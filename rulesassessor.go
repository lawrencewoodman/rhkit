/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package rulehunter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/dlit_go"
	"github.com/lawrencewoodman/rulehunter/internal"
	"io"
	"os"
	"sort"
)

type Assessment struct {
	NumRecords      int64
	RuleAssessments []*RuleFinalAssessment
	Flags           map[string]bool
}

type RuleFinalAssessment struct {
	Rule        *internal.Rule
	Aggregators map[string]*dlit.Literal
	Goals       map[string]bool
}

// by implements sort.Interface for []*RuleFinalAssessments based
// on the sortFields
type by struct {
	ruleAssessments []*RuleFinalAssessment
	sortFields      []SortField
}

func (b by) Len() int { return len(b.ruleAssessments) }
func (b by) Swap(i, j int) {
	b.ruleAssessments[i], b.ruleAssessments[j] =
		b.ruleAssessments[j], b.ruleAssessments[i]
}
func (b by) Less(i, j int) bool {
	var vI *dlit.Literal
	var vJ *dlit.Literal
	for _, sortField := range b.sortFields {
		field := sortField.Field
		direction := sortField.Direction
		// TODO: Perhaps ignore case
		if field == "numGoalsPassed" {
			// TODO: Work out if this should be calculated here, or elsewhere?
			vI = calcNumGoalsPassedScore(b.ruleAssessments[i])
			vJ = calcNumGoalsPassedScore(b.ruleAssessments[j])
		} else {
			vI = b.ruleAssessments[i].Aggregators[field]
			vJ = b.ruleAssessments[j].Aggregators[field]
		}
		c := compareDlitNums(vI, vJ)

		if direction == DESCENDING {
			c *= -1
		}
		if c < 0 {
			return true
		} else if c > 0 {
			return false
		}
	}

	ruleLenI := len(b.ruleAssessments[i].Rule.String())
	ruleLenJ := len(b.ruleAssessments[j].Rule.String())
	return ruleLenI < ruleLenJ
}

func compareDlitNums(l1 *dlit.Literal, l2 *dlit.Literal) int {
	i1, l1IsInt := l1.Int()
	i2, l2IsInt := l2.Int()
	if l1IsInt && l2IsInt {
		if i1 < i2 {
			return -1
		}
		if i1 > i2 {
			return 1
		}
		return 0
	}

	f1, l1IsFloat := l1.Float()
	f2, l2IsFloat := l2.Float()

	if l1IsFloat && l2IsFloat {
		if f1 < f2 {
			return -1
		}
		if f1 > f2 {
			return 1
		}
		return 0
	}
	panic(fmt.Sprintf("Can't compare numbers: %s, %s", l1, l2))
}

func (r *Assessment) Sort(s []SortField) {
	sort.Sort(by{r.RuleAssessments, s})
	r.Flags["sorted"] = true
}

// TODO: Test this
func (r *Assessment) IsEqual(o *Assessment) bool {
	if r.NumRecords != o.NumRecords {
		return false
	}
	for i, ruleAssessment := range r.RuleAssessments {
		if !ruleAssessment.isEqual(o.RuleAssessments[i]) {
			return false
		}
	}
	if len(r.Flags) != len(o.Flags) {
		return false
	}
	for k, v := range r.Flags {
		if v != o.Flags[k] {
			return false
		}
	}
	return true
}

func (r *RuleFinalAssessment) isEqual(o *RuleFinalAssessment) bool {
	if r.Rule.String() != o.Rule.String() {
		return false
	}
	if len(r.Aggregators) != len(o.Aggregators) {
		return false
	}
	for aName, value := range r.Aggregators {
		if o.Aggregators[aName].String() != value.String() {
			return false
		}
	}
	if len(r.Goals) != len(o.Goals) {
		return false
	}
	for gName, passed := range r.Goals {
		if o.Goals[gName] != passed {
			return false
		}
	}
	return true
}

func (r *RuleFinalAssessment) String() string {
	return fmt.Sprintf("Rule: %s, Aggregators: %s, Goals: %s",
		r.Rule, r.Aggregators, r.Goals)
}

type JReport struct {
	NumRecords      int64
	RuleAssessments []*JRuleReport
}

type JRuleReport struct {
	Rule        string
	Aggregators map[string]string
	Goals       map[string]bool
}

func (a *Assessment) ToJSON() (string, error) {
	jRuleAssessments := make([]*JRuleReport, len(a.RuleAssessments))
	for i, ruleAssessment := range a.RuleAssessments {
		jRuleAssessments[i] = makeJRuleReport(ruleAssessment)
	}
	jReport := &JReport{a.NumRecords, jRuleAssessments}
	b, err := json.MarshalIndent(jReport, "", "  ")
	if err != nil {
		os.Stdout.Write(b)
	}
	return string(b[:]), err
}

func makeJRuleReport(r *RuleFinalAssessment) *JRuleReport {
	aggregators := make(map[string]string, len(r.Aggregators))
	for n, l := range r.Aggregators {
		aggregators[n] = l.String()
	}
	return &JRuleReport{r.Rule.String(), aggregators, r.Goals}
}

func (a *Assessment) GetRules() []*internal.Rule {
	r := make([]*internal.Rule, len(a.RuleAssessments))
	for i, ruleAssessment := range a.RuleAssessments {
		r[i] = ruleAssessment.Rule
	}
	return r
}

// Tidy up rule assessments by removing poor and poorer similar rules
// For example this removes all rules poorer than the 'true()' rule
func (sortedAssessment *Assessment) Refine(numSimilarRules int) {
	if !sortedAssessment.Flags["sorted"] {
		panic("Assessment isn't sorted")
	}
	sortedAssessment.excludePoorRules()
	sortedAssessment.excludePoorerInNiRules(numSimilarRules)
	sortedAssessment.excludePoorerTweakableRules(numSimilarRules)
}

func (sortedAssessment *Assessment) excludePoorRules() {
	trueFound := false
	goodRuleAssessments := make([]*RuleFinalAssessment, 0)
	for _, a := range sortedAssessment.RuleAssessments {
		numMatches, numMatchesIsInt := a.Aggregators["numMatches"].Int()
		if !numMatchesIsInt {
			panic("numMatches aggregator isn't an int")
		}
		if numMatches > 1 {
			goodRuleAssessments = append(goodRuleAssessments, a)
		}
		if a.Rule.String() == "true()" {
			trueFound = true
			break
		}
	}
	if !trueFound {
		panic("No 'true()' rule found")
	}
	sortedAssessment.RuleAssessments = goodRuleAssessments
}

func (sortedAssessment *Assessment) excludePoorerInNiRules(
	numSimilarRules int,
) {
	goodRuleAssessments := make([]*RuleFinalAssessment, 0)
	inFields := make(map[string]int)
	niFields := make(map[string]int)
	for _, a := range sortedAssessment.RuleAssessments {
		rule := a.Rule
		isInNiRule, operator, field := rule.GetInNiParts()
		if !isInNiRule {
			goodRuleAssessments = append(goodRuleAssessments, a)
		} else if operator == "in" {
			n, ok := inFields[field]
			if !ok {
				goodRuleAssessments = append(goodRuleAssessments, a)
				inFields[field] = 1
			} else if n < numSimilarRules {
				goodRuleAssessments = append(goodRuleAssessments, a)
				inFields[field]++
			}
		} else if operator == "ni" {
			n, ok := niFields[field]
			if !ok {
				goodRuleAssessments = append(goodRuleAssessments, a)
				niFields[field] = 1
			} else if n < numSimilarRules {
				goodRuleAssessments = append(goodRuleAssessments, a)
				niFields[field]++
			}
		}
	}
	sortedAssessment.RuleAssessments = goodRuleAssessments
}

func (sortedAssessment *Assessment) excludePoorerTweakableRules(
	numSimilarRules int,
) {
	goodRuleAssessments := make([]*RuleFinalAssessment, 0)
	fieldOperatorIDs := make(map[string]int)
	for _, a := range sortedAssessment.RuleAssessments {
		rule := a.Rule
		isTweakable, field, operator, _ := rule.GetTweakableParts()
		if !isTweakable {
			goodRuleAssessments = append(goodRuleAssessments, a)
		} else {
			fieldOperatorID := fmt.Sprintf("%s^%s", field, operator)
			n, ok := fieldOperatorIDs[fieldOperatorID]
			if !ok {
				goodRuleAssessments = append(goodRuleAssessments, a)
				fieldOperatorIDs[fieldOperatorID] = 1
			} else if n < numSimilarRules {
				goodRuleAssessments = append(goodRuleAssessments, a)
				fieldOperatorIDs[fieldOperatorID]++
			}
		}
	}
	sortedAssessment.RuleAssessments = goodRuleAssessments
}

type ErrNameConflict string

func (e ErrNameConflict) Error() string {
	return string(e)
}

func (a *Assessment) Merge(o *Assessment) (*Assessment, error) {
	if a.NumRecords != o.NumRecords {
		// TODO: Create error type
		err := errors.New("Can't merge assessments: Number of records don't match")
		return nil, err
	}
	newRuleAssessments := append(a.RuleAssessments, o.RuleAssessments...)
	flags := map[string]bool{
		"sorted": false,
	}
	return &Assessment{a.NumRecords, newRuleAssessments, flags}, nil
}

// need a progress callback and a specifier for how often to report
func AssessRules(rules []*internal.Rule, aggregators []internal.Aggregator,
	goals []*dexpr.Expr, input internal.Input) (*Assessment, error) {
	var allAggregators []internal.Aggregator
	var numRecords int64
	var err error

	allAggregators, err = prependDefaultAggregators(aggregators)
	if err != nil {
		return &Assessment{}, err
	}
	/*
		TODO: Put this test somewhere else
		err := checkForNameConflicts(fieldNames, aggregators)
		if err != nil {
			return &[]RuleAssessment{}, err
		}
	*/

	ruleAssessments := make([]*internal.RuleAssessment, len(rules))
	for i, rule := range rules {
		ruleAssessments[i] = internal.NewRuleAssessment(rule, allAggregators, goals)
	}
	numRecords, err = processInput(input, ruleAssessments)
	if err != nil {
		return &Assessment{}, err
	}
	goodRuleAssessments, err := filterGoodReports(ruleAssessments, numRecords)
	if err != nil {
		return &Assessment{}, err
	}

	assessment, err := makeAssessment(numRecords, goodRuleAssessments, goals)
	return assessment, err
}

func makeAssessment(
	numRecords int64,
	goodRuleAssessments []*internal.RuleAssessment,
	goals []*dexpr.Expr,
) (*Assessment, error) {
	ruleAssessments := make([]*RuleFinalAssessment, len(goodRuleAssessments))
	for i, ruleAssessment := range goodRuleAssessments {
		rule := ruleAssessment.Rule
		aggregatorsMap :=
			internal.AggregatorsToMap(ruleAssessment.Aggregators, numRecords, "")
		goals, err := internal.GoalsToMap(ruleAssessment.Goals, aggregatorsMap)
		if err != nil {
			return &Assessment{}, err
		}
		delete(aggregatorsMap, "numRecords")
		ruleAssessments[i] = &RuleFinalAssessment{
			Rule:        rule,
			Aggregators: aggregatorsMap,
			Goals:       goals,
		}
	}
	flags := map[string]bool{
		"sorted": false,
	}
	assessment := &Assessment{
		NumRecords:      numRecords,
		RuleAssessments: ruleAssessments,
		Flags:           flags,
	}
	return assessment, nil
}

func filterGoodReports(
	ruleAssessments []*internal.RuleAssessment,
	numRecords int64) ([]*internal.RuleAssessment, error) {
	goodRuleAssessments := make([]*internal.RuleAssessment, 0)
	for _, ruleAssessment := range ruleAssessments {
		numMatches, exists :=
			ruleAssessment.GetAggregatorValue("numMatches", numRecords)
		if !exists {
			// TODO: Create a proper error for this?
			err := errors.New("numMatches doesn't exist in aggregators")
			return goodRuleAssessments, err
		}
		numMatchesInt, isInt := numMatches.Int()
		if !isInt {
			// TODO: Create a proper error for this?
			err := errors.New(fmt.Sprintf("Can't cast to Int: %q", numMatches))
			return goodRuleAssessments, err
		}
		if numMatchesInt > 0 {
			goodRuleAssessments = append(goodRuleAssessments, ruleAssessment)
		}
	}
	return goodRuleAssessments, nil
}

func processInput(input internal.Input,
	ruleAssessments []*internal.RuleAssessment) (int64, error) {
	numRecords := int64(0)
	// TODO: test this rewinds properly
	if err := input.Rewind(); err != nil {
		return numRecords, err
	}

	for {
		record, err := input.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return numRecords, err
		}
		numRecords++
		for _, ruleAssessment := range ruleAssessments {
			err := ruleAssessment.NextRecord(record)
			if err != nil {
				return numRecords, err
			}
		}
	}
	return numRecords, nil
}

func prependDefaultAggregators(
	aggregators []internal.Aggregator,
) ([]internal.Aggregator, error) {
	newAggregators := make([]internal.Aggregator, 2)
	numMatchesAggregator, err := internal.NewCountAggregator("numMatches", "1==1")
	if err != nil {
		return newAggregators, err
	}
	percentMatchesAggregator, err :=
		internal.NewCalcAggregator("percentMatches",
			"roundto(100.0 * numMatches / numRecords, 2)")
	if err != nil {
		return newAggregators, err
	}
	newAggregators[0] = numMatchesAggregator
	newAggregators[1] = percentMatchesAggregator
	newAggregators = append(newAggregators, aggregators...)
	return newAggregators, nil
}

func calcNumGoalsPassedScore(r *RuleFinalAssessment) *dlit.Literal {
	numGoalsPassed := 0.0
	increment := 1.0
	for _, goalPassed := range r.Goals {
		if goalPassed {
			numGoalsPassed += increment
		} else {
			increment = 0.001
		}
	}
	return dlit.MustNew(numGoalsPassed)
}

/* TODO: Put this somewhere else
func checkForNameConflicts(fields []string, aggregators []Aggregator) error {
	for _, aggregator := range aggregators {
		for _, fieldName := range fields {
			if aggregator.GetName() == fieldName {
				return ErrNameConflict(
					fmt.Sprintf("Aggregator name and field name conflict: %s",
						fieldName))
			}
		}
	}
	return nil
}
*/
