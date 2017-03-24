/*
	Copyright (C) 2016-2017 vLife Systems Ltd <http://vlifesystems.com>
	This file is part of rhkit.

	rhkit is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	rhkit is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with rhkit; see the file COPYING.  If not, see
	<http://www.gnu.org/licenses/>.
*/

// Package experiment handles initialization and validation of experiment
package experiment

import (
	"github.com/lawrencewoodman/ddataset"
	"github.com/vlifesystems/rhkit/aggregators"
	"github.com/vlifesystems/rhkit/goal"
	"github.com/vlifesystems/rhkit/internal"
)

type ExperimentDesc struct {
	Title       string
	Dataset     ddataset.Dataset
	RuleFields  []string
	Aggregators []*AggregatorDesc
	Goals       []string
	SortOrder   []*SortDesc
}

type AggregatorDesc struct {
	Name     string
	Function string
	Arg      string
}

type SortDesc struct {
	AggregatorName string
	Direction      string
}

type Experiment struct {
	Title          string
	Dataset        ddataset.Dataset
	RuleFieldNames []string
	Aggregators    []aggregators.AggregatorSpec
	Goals          []*goal.Goal
	SortOrder      []SortField
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
	goals, err := makeGoals(e.Goals)
	if err != nil {
		return nil, err
	}
	aggregators, err := makeAggregators(e.Aggregators)
	if err != nil {
		return nil, err
	}

	sortOrder, err := makeSortOrder(e.SortOrder)
	if err != nil {
		return nil, err
	}

	return &Experiment{
		Title:          e.Title,
		Dataset:        e.Dataset,
		RuleFieldNames: e.RuleFields,
		Aggregators:    aggregators,
		Goals:          goals,
		SortOrder:      sortOrder,
	}, nil
}

func checkExperimentDescValid(e *ExperimentDesc) error {
	if err := checkSortDescsValid(e); err != nil {
		return err
	}

	if err := checkRuleFieldsValid(e); err != nil {
		return err
	}

	if err := checkAggregatorsValid(e); err != nil {
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
	fieldNames := e.Dataset.GetFieldNames()
	for _, ruleField := range e.RuleFields {
		if !internal.IsIdentifierValid(ruleField) {
			return InvalidRuleFieldError(ruleField)
		}
		if !isStringInSlice(ruleField, fieldNames) {
			return InvalidRuleFieldError(ruleField)
		}
	}
	return nil
}

func isStringInSlice(needle string, haystack []string) bool {
	for _, s := range haystack {
		if needle == s {
			return true
		}
	}
	return false
}

func checkAggregatorsValid(e *ExperimentDesc) error {
	fieldNames := e.Dataset.GetFieldNames()
	for _, aggregator := range e.Aggregators {
		if !internal.IsIdentifierValid(aggregator.Name) {
			return InvalidAggregatorNameError(aggregator.Name)
		}
		if isStringInSlice(aggregator.Name, fieldNames) {
			return AggregatorNameClashError(aggregator.Name)
		}
		if aggregator.Name == "percentMatches" ||
			aggregator.Name == "numMatches" ||
			aggregator.Name == "goalsScore" {
			return AggregatorNameReservedError(aggregator.Name)
		}
	}
	return nil
}

func makeGoals(exprs []string) ([]*goal.Goal, error) {
	var err error
	r := make([]*goal.Goal, len(exprs))
	for i, expr := range exprs {
		r[i], err = goal.New(expr)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

func makeAggregators(
	eAggregators []*AggregatorDesc,
) ([]aggregators.AggregatorSpec, error) {
	var err error
	r := make([]aggregators.AggregatorSpec, len(eAggregators))
	for i, ea := range eAggregators {
		r[i], err = aggregators.New(ea.Name, ea.Function, ea.Arg)
		if err != nil {
			return r, err
		}
	}
	return addDefaultAggregators(r), nil
}

func addDefaultAggregators(
	aggregatorSpecs []aggregators.AggregatorSpec,
) []aggregators.AggregatorSpec {
	newAggregatorSpecs := make([]aggregators.AggregatorSpec, 2)
	newAggregatorSpecs[0] = aggregators.MustNew("numMatches", "count", "true()")
	newAggregatorSpecs[1] = aggregators.MustNew(
		"percentMatches",
		"calc",
		"roundto(100.0 * numMatches / numRecords, 2)",
	)
	goalsScoreAggregatorSpec := aggregators.MustNew("goalsScore", "goalsscore")
	newAggregatorSpecs = append(newAggregatorSpecs, aggregatorSpecs...)
	newAggregatorSpecs = append(newAggregatorSpecs, goalsScoreAggregatorSpec)
	return newAggregatorSpecs
}

func makeSortOrder(eSortOrder []*SortDesc) ([]SortField, error) {
	r := make([]SortField, len(eSortOrder))
	for i, eSortField := range eSortOrder {
		field := eSortField.AggregatorName
		direction := eSortField.Direction
		// TODO: Make case insensitive
		if direction == "ascending" {
			r[i] = SortField{field, ASCENDING}
		} else if direction == "descending" {
			r[i] = SortField{field, DESCENDING}
		} else {
			err := &InvalidSortDirectionError{field, direction}
			return r, err
		}
	}
	return r, nil
}
