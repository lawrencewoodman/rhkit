/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package main

import (
	"errors"
	"fmt"
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/dlit_go"
	"io"
	"sort"
	"strconv"
)

type Report struct {
	NumRecords  int64
	RuleReports []*RuleReport
}

type RuleReport struct {
	Rule        string
	Aggregators map[string]string
	Goals       map[string]bool
}

type SortField struct {
	Field     string
	Direction direction
}

type direction int

const (
	ASCENDING direction = iota
	DESCENDING
)

func (d direction) String() string {
	if d == ASCENDING {
		return "ascending"
	}
	return "descending"
}

// by implements sort.Interface for []*RuleReports based on the sortFields
type by struct {
	ruleReports []*RuleReport
	sortFields  []SortField
}

func (b by) Len() int { return len(b.ruleReports) }
func (b by) Swap(i, j int) {
	b.ruleReports[i], b.ruleReports[j] = b.ruleReports[j], b.ruleReports[i]
}
func (b by) Less(i, j int) bool {
	var vI string
	var vJ string
	for _, sortField := range b.sortFields {
		field := sortField.Field
		direction := sortField.Direction
		// TODO: Perhaps ignore case
		if field == "numGoalsPassed" {
			// TODO: Work out if this should be calculated here, or elsewhere?
			vI = calcNumGoalsPassedScore(b.ruleReports[i])
			vJ = calcNumGoalsPassedScore(b.ruleReports[j])
		} else {
			vI = b.ruleReports[i].Aggregators[field]
			vJ = b.ruleReports[j].Aggregators[field]
		}
		c := compareStrNums(vI, vJ)

		if direction == DESCENDING {
			c *= -1
		}
		if c < 0 {
			return true
		} else if c > 0 {
			return false
		}
	}

	ruleLenI := len(b.ruleReports[i].Rule)
	ruleLenJ := len(b.ruleReports[j].Rule)
	return ruleLenI < ruleLenJ
}

func compareStrNums(nStr1 string, nStr2 string) int {
	i1, errI1 := strconv.ParseInt(nStr1, 10, 64)
	i2, errI2 := strconv.ParseInt(nStr2, 10, 64)
	if errI1 == nil && errI2 == nil {
		if i1 < i2 {
			return -1
		}
		if i1 > i2 {
			return 1
		}
		return 0
	}
	f1, errF1 := strconv.ParseFloat(nStr1, 64)
	f2, errF2 := strconv.ParseFloat(nStr2, 64)
	if errF1 == nil && errF2 == nil {
		if f1 < f2 {
			return -1
		}
		if f1 > f2 {
			return 1
		}
		return 0
	}
	panic(fmt.Sprintf("Can't compare strings as numbers: %s, %s", nStr1, nStr2))
}

func (r *Report) Sort(s []SortField) {
	sort.Sort(by{r.RuleReports, s})
}

// TODO: Test this
func (r *Report) IsEqual(o *Report) bool {
	if r.NumRecords != o.NumRecords {
		return false
	}
	for i, ruleReport := range r.RuleReports {
		if !ruleReport.isEqual(o.RuleReports[i]) {
			return false
		}
	}
	return true
}

func (r *RuleReport) isEqual(o *RuleReport) bool {
	if r.Rule != o.Rule {
		return false
	}
	if len(r.Aggregators) != len(o.Aggregators) {
		return false
	}
	for aName, value := range r.Aggregators {
		if o.Aggregators[aName] != value {
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

func (r *RuleReport) String() string {
	return fmt.Sprintf("Rule: %s, Aggregators: %s, Goals: %s",
		r.Rule, r.Aggregators, r.Goals)
}

type ErrNameConflict string

func (e ErrNameConflict) Error() string {
	return string(e)
}

// need a progress callback and a specifier for how often to report
func AssessRules(rules []*dexpr.Expr, aggregators []Aggregator,
	goals []*dexpr.Expr, input Input) (*Report, error) {
	var allAggregators []Aggregator
	var numRecords int64
	var err error

	allAggregators, err = prependDefaultAggregators(aggregators)
	if err != nil {
		return &Report{}, err
	}
	/*
		TODO: Put this test somewhere else
		err := checkForNameConflicts(fieldNames, aggregators)
		if err != nil {
			return &[]RuleAssessment{}, err
		}
	*/

	ruleAssessments := make([]*RuleAssessment, len(rules))
	for i, rule := range rules {
		ruleAssessments[i] = NewRuleAssessment(rule, allAggregators, goals)
	}
	numRecords, err = processInput(input, ruleAssessments)
	if err != nil {
		return &Report{}, err
	}
	goodRuleAssessments, err := filterGoodReports(ruleAssessments, numRecords)
	if err != nil {
		return &Report{}, err
	}

	report, err := makeReport(numRecords, goodRuleAssessments, goals)
	return report, err
}

func makeReport(numRecords int64, goodRuleAssessments []*RuleAssessment,
	goals []*dexpr.Expr) (*Report, error) {
	ruleReports := make([]*RuleReport, len(goodRuleAssessments))
	for i, ruleAssessment := range goodRuleAssessments {
		rule := ruleAssessment.rule.String()
		aggregators :=
			AggregatorsToMap(ruleAssessment.aggregators, numRecords, "")
		goals, err := GoalsToMap(ruleAssessment.goals, aggregators)
		if err != nil {
			return &Report{}, err
		}
		delete(aggregators, "numRecords")
		ruleReports[i] = &RuleReport{Rule: rule,
			Aggregators: makeRuleReportAggregators(aggregators), Goals: goals}
	}
	return &Report{NumRecords: numRecords, RuleReports: ruleReports}, nil
}

func makeRuleReportAggregators(
	aMap map[string]*dlit.Literal) map[string]string {
	r := make(map[string]string, len(aMap))
	for n, v := range aMap {
		r[n] = v.String()
	}
	return r
}

func filterGoodReports(
	ruleAssessments []*RuleAssessment,
	numRecords int64) ([]*RuleAssessment, error) {
	goodRuleAssessments := make([]*RuleAssessment, 0)
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

func processInput(input Input,
	ruleAssessments []*RuleAssessment) (int64, error) {
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

func prependDefaultAggregators(aggregators []Aggregator) ([]Aggregator, error) {
	newAggregators := make([]Aggregator, 2)
	numMatchesAggregator, err := NewCountAggregator("numMatches", "1==1")
	if err != nil {
		return newAggregators, err
	}
	percentMatchesAggregator, err :=
		NewCalcAggregator("percentMatches",
			"roundto(100.0 * numMatches / numRecords, 2)")
	if err != nil {
		return newAggregators, err
	}
	newAggregators[0] = numMatchesAggregator
	newAggregators[1] = percentMatchesAggregator
	newAggregators = append(newAggregators, aggregators...)
	return newAggregators, nil
}

func calcNumGoalsPassedScore(r *RuleReport) string {
	numGoalsPassed := 0.0
	increment := 1.0
	for _, goalPassed := range r.Goals {
		if goalPassed {
			numGoalsPassed += increment
		} else {
			increment = 0.001
		}
	}
	return fmt.Sprintf("%f", numGoalsPassed)
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
