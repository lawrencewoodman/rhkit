// Copyright (C) 2016-2017 vLife Systems Ltd <http://vlifesystems.com>
// Licensed under an MIT licence.  Please see LICENSE.md for details.

// Package assessment assesses rules to meet user defined goals
package assessment

import (
	"errors"
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/aggregator"
	"github.com/vlifesystems/rhkit/rule"
	"sort"
)

type Assessment struct {
	NumRecords      int64
	RuleAssessments []*RuleAssessment
	flags           map[string]bool
}

type RuleAssessment struct {
	Rule        rule.Rule
	Aggregators map[string]*dlit.Literal
	Goals       []*GoalAssessment
}

type GoalAssessment struct {
	Expr   string
	Passed bool
}

func newAssessment(numRecords int64) *Assessment {
	return &Assessment{
		NumRecords: numRecords,
		flags: map[string]bool{
			"sorted":  false,
			"refined": false,
		},
	}
}

func (a *Assessment) AddRuleAssessors(ruleAssessors []*ruleAssessor) error {
	ruleAssessments := make([]*RuleAssessment, len(ruleAssessors))
	for i, ruleAssessment := range ruleAssessors {
		rule := ruleAssessment.Rule
		aggregatorInstancesMap, err :=
			aggregator.InstancesToMap(
				ruleAssessment.Aggregators,
				ruleAssessment.Goals,
				a.NumRecords,
			)
		if err != nil {
			return err
		}
		goalAssessments := make([]*GoalAssessment, len(ruleAssessment.Goals))
		for j, goal := range ruleAssessment.Goals {
			passed, err := goal.Assess(aggregatorInstancesMap)
			if err != nil {
				return err
			}
			goalAssessments[j] = &GoalAssessment{goal.String(), passed}
		}
		delete(aggregatorInstancesMap, "numRecords")
		ruleAssessments[i] = &RuleAssessment{
			Rule:        rule,
			Aggregators: aggregatorInstancesMap,
			Goals:       goalAssessments,
		}
	}
	a.RuleAssessments = ruleAssessments
	return nil
}

func (a *Assessment) Sort(s []SortOrder) {
	sort.Sort(by{a.RuleAssessments, s})
	a.flags["sorted"] = true
}

func (a *Assessment) IsSorted() bool {
	return a.flags["sorted"]
}

func (a *Assessment) IsRefined() bool {
	return a.flags["refined"]
}

// TODO: Test this
func (a *Assessment) IsEqual(o *Assessment) bool {
	if a.NumRecords != o.NumRecords {
		return false
	}

	if len(a.RuleAssessments) != len(o.RuleAssessments) {
		return false
	}
	for i, ruleAssessment := range a.RuleAssessments {
		if !ruleAssessment.IsEqual(o.RuleAssessments[i]) {
			return false
		}
	}

	if len(a.flags) != len(o.flags) {
		return false
	}
	for k, v := range a.flags {
		if v != o.flags[k] {
			return false
		}
	}

	return true
}

// Tidy up rule assessments by removing poor and poorer similar rules
// For example this removes all rules poorer than the True rule
func (sortedAssessment *Assessment) Refine() {
	if !sortedAssessment.IsSorted() {
		panic("Assessment isn't sorted")
	}
	sortedAssessment.excludePoorRules()
	sortedAssessment.excludeSameRecordsRules()
	sortedAssessment.excludePoorerOverlappingRules()
	sortedAssessment.flags["refined"] = true
}

func (a *Assessment) Merge(o *Assessment) (*Assessment, error) {
	if a.NumRecords != o.NumRecords {
		// TODO: Create error type
		err := errors.New("Can't merge assessments: Number of records don't match")
		return nil, err
	}
	newRuleAssessments := append(a.RuleAssessments, o.RuleAssessments...)
	flags := map[string]bool{
		"sorted":  false,
		"refined": false,
	}
	return &Assessment{a.NumRecords, newRuleAssessments, flags}, nil
}

// Assessment must be sorted and refined first
func (a *Assessment) TruncateRuleAssessments(
	numRuleAssessments int,
) *Assessment {
	if !a.IsSorted() {
		panic("Assessment isn't sorted")
	}
	if !a.IsRefined() {
		panic("Assessment isn't refined")
	}
	if len(a.RuleAssessments) < numRuleAssessments {
		numRuleAssessments = len(a.RuleAssessments)
	}
	numNonTrueRuleAssessments := numRuleAssessments - 1

	ruleAssessments := make([]*RuleAssessment, numRuleAssessments)
	for i := 0; i < numNonTrueRuleAssessments; i++ {
		ruleAssessments[i] = a.RuleAssessments[i].clone()
	}

	if numRuleAssessments > 0 {
		trueRuleAssessment := a.RuleAssessments[len(a.RuleAssessments)-1]
		if _, isTrueRule := trueRuleAssessment.Rule.(rule.True); !isTrueRule {
			panic("Assessment doesn't have True rule last")
		}

		ruleAssessments[numNonTrueRuleAssessments] = trueRuleAssessment
	}

	flags := map[string]bool{
		"sorted":  true,
		"refined": true,
	}
	return &Assessment{a.NumRecords, ruleAssessments, flags}
}

// Can optionally pass maximum number of rules to return
func (a *Assessment) Rules(args ...int) []rule.Rule {
	var numRules int
	switch len(args) {
	case 0:
		numRules = len(a.RuleAssessments)
	case 1:
		numRules = args[0]
		if len(a.RuleAssessments) < numRules {
			numRules = len(a.RuleAssessments)
		}
	default:
		panic(fmt.Sprintf("incorrect number of arguments passed: %d", len(args)))
	}

	r := make([]rule.Rule, numRules)
	for i, ruleAssessment := range a.RuleAssessments {
		if i >= numRules {
			break
		}
		r[i] = ruleAssessment.Rule
	}
	return r
}

func (sortedAssessment *Assessment) excludeSameRecordsRules() {
	if len(sortedAssessment.RuleAssessments) < 2 {
		return
	}
	lastAggregators := sortedAssessment.RuleAssessments[0].Aggregators
	if len(lastAggregators) <= 3 {
		return
	}

	goodRuleAssessments := make([]*RuleAssessment, 1)
	goodRuleAssessments[0] = sortedAssessment.RuleAssessments[0]
	for _, a := range sortedAssessment.RuleAssessments[1:] {
		aggregatorsMatch := true
		for k, v := range lastAggregators {
			if a.Aggregators[k].String() != v.String() {
				aggregatorsMatch = false
			}
		}
		switch a.Rule.(type) {
		case rule.True:
			if aggregatorsMatch {
				goodRuleAssessments[len(goodRuleAssessments)-1] = a
			} else {
				goodRuleAssessments = append(goodRuleAssessments, a)
			}
			break
		default:
			if !aggregatorsMatch {
				goodRuleAssessments = append(goodRuleAssessments, a)
			}
		}
		lastAggregators = a.Aggregators
	}
	sortedAssessment.RuleAssessments = goodRuleAssessments
}

func (sortedAssessment *Assessment) excludePoorRules() {
	trueFound := false
	goodRuleAssessments := make([]*RuleAssessment, 0)
	for _, a := range sortedAssessment.RuleAssessments {
		numMatches, numMatchesIsInt := a.Aggregators["numMatches"].Int()
		if !numMatchesIsInt {
			panic("numMatches aggregator isn't an int")
		}
		if numMatches > 1 {
			goodRuleAssessments = append(goodRuleAssessments, a)
		}
		if _, isTrueRule := a.Rule.(rule.True); isTrueRule {
			trueFound = true
			break
		}
	}
	if !trueFound {
		panic("No True rule found")
	}
	sortedAssessment.RuleAssessments = goodRuleAssessments
}

func (sortedAssessment *Assessment) excludePoorerOverlappingRules() {
	goodRuleAssessments := make([]*RuleAssessment, 0)
	for i, aI := range sortedAssessment.RuleAssessments {
		switch xI := aI.Rule.(type) {
		case rule.Overlapper:
			overlaps := false
			for j, aJ := range sortedAssessment.RuleAssessments {
				if j >= i {
					break
				}
				if xI.Overlaps(aJ.Rule) {
					overlaps = true
				}
			}
			if !overlaps {
				goodRuleAssessments = append(goodRuleAssessments, aI)
			}
		default:
			goodRuleAssessments = append(goodRuleAssessments, aI)
		}
	}
	sortedAssessment.RuleAssessments = goodRuleAssessments
}

func (r *RuleAssessment) IsEqual(o *RuleAssessment) bool {
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
	for i, goal := range r.Goals {
		if !goal.IsEqual(o.Goals[i]) {
			return false
		}
	}
	return true
}

func (r *RuleAssessment) clone() *RuleAssessment {
	return &RuleAssessment{r.Rule, r.Aggregators, r.Goals}
}

func (g *GoalAssessment) IsEqual(o *GoalAssessment) bool {
	return g.Expr == o.Expr && g.Passed == o.Passed
}
