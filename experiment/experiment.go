// Copyright (C) 2016-2017 vLife Systems Ltd <http://vlifesystems.com>
// Licensed under an MIT licence.  Please see LICENSE.md for details.

// Package experiment handles initialization and validation of experiment
package experiment

import (
	"github.com/lawrencewoodman/ddataset"
	"github.com/vlifesystems/rhkit/aggregators"
	"github.com/vlifesystems/rhkit/goal"
	"github.com/vlifesystems/rhkit/internal"
	"github.com/vlifesystems/rhkit/rule"
)

type ExperimentDesc struct {
	Dataset        ddataset.Dataset
	RuleFields     []string
	RuleComplexity rule.Complexity
	Aggregators    []*aggregators.Desc
	Goals          []string
	SortOrder      []*SortDesc
	Rules          []string
}

type SortDesc struct {
	AggregatorName string
	Direction      string
}

type Experiment struct {
	Dataset        ddataset.Dataset
	RuleFields     []string
	RuleComplexity rule.Complexity
	Aggregators    []aggregators.Spec
	Goals          []*goal.Goal
	SortOrder      []SortField
	Rules          []rule.Rule
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

// Create a new Experiment from the description
func New(e *ExperimentDesc) (*Experiment, error) {
	if err := checkExperimentDescValid(e); err != nil {
		return nil, err
	}
	goals, err := goal.MakeGoals(e.Goals)
	if err != nil {
		return nil, err
	}
	aggregators, err :=
		aggregators.MakeSpecs(e.Dataset.Fields(), e.Aggregators)
	if err != nil {
		return nil, err
	}

	sortOrder, err := makeSortOrder(e.SortOrder)
	if err != nil {
		return nil, err
	}

	rules, err := rule.MakeDynamicRules(e.Rules)
	if err != nil {
		return nil, err
	}

	return &Experiment{
		Dataset:        e.Dataset,
		RuleFields:     e.RuleFields,
		RuleComplexity: e.RuleComplexity,
		Aggregators:    aggregators,
		Goals:          goals,
		SortOrder:      sortOrder,
		Rules:          rules,
	}, nil
}

func checkExperimentDescValid(e *ExperimentDesc) error {
	if err := checkSortDescsValid(e); err != nil {
		return err
	}
	if err := checkRuleFieldsValid(e); err != nil {
		return err
	}
	return nil
}

func checkSortDescsValid(e *ExperimentDesc) error {
	for _, sortDesc := range e.SortOrder {
		if sortDesc.Direction != "ascending" && sortDesc.Direction != "descending" {
			return &InvalidSortDirectionError{
				sortDesc.AggregatorName,
				sortDesc.Direction,
			}
		}
		sortName := sortDesc.AggregatorName
		nameFound := false
		for _, aggregator := range e.Aggregators {
			if aggregator.Name == sortName {
				nameFound = true
				break
			}
		}
		if !nameFound &&
			sortName != "percentMatches" &&
			sortName != "numMatches" &&
			sortName != "goalsScore" {
			return InvalidSortFieldError(sortName)
		}
	}
	return nil
}

func checkRuleFieldsValid(e *ExperimentDesc) error {
	if len(e.RuleFields) == 0 {
		return ErrNoRuleFieldsSpecified
	}
	fieldNames := e.Dataset.Fields()
	for _, ruleField := range e.RuleFields {
		if !internal.IsIdentifierValid(ruleField) {
			return InvalidRuleFieldError(ruleField)
		}
		if !internal.IsStringInSlice(ruleField, fieldNames) {
			return InvalidRuleFieldError(ruleField)
		}
	}
	return nil
}

func makeSortOrder(eSortOrder []*SortDesc) ([]SortField, error) {
	r := make([]SortField, len(eSortOrder))
	for i, eSortField := range eSortOrder {
		field := eSortField.AggregatorName
		direction := eSortField.Direction
		// TODO: Make case insensitive
		if direction == "ascending" {
			r[i] = SortField{field, ASCENDING}
		} else {
			r[i] = SortField{field, DESCENDING}
		}
	}
	return r, nil
}
