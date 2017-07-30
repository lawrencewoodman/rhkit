// Copyright (C) 2016-2017 vLife Systems Ltd <http://vlifesystems.com>
// Licensed under an MIT licence.  Please see LICENSE.md for details.

// Package experiment handles initialization and validation of experiment
package experiment

import (
	"github.com/lawrencewoodman/ddataset"
	"github.com/vlifesystems/rhkit/aggregators"
	"github.com/vlifesystems/rhkit/assessment"
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
	SortOrder      []assessment.SortDesc
	Rules          []string
}

type Experiment struct {
	Dataset        ddataset.Dataset
	RuleFields     []string
	RuleComplexity rule.Complexity
	Aggregators    []aggregators.Spec
	Goals          []*goal.Goal
	SortOrder      []assessment.SortOrder
	Rules          []rule.Rule
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

	sortOrder, err := assessment.MakeSortOrders(aggregators, e.SortOrder)
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
	if err := checkRuleFieldsValid(e); err != nil {
		return err
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
