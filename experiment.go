/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/rulehunter/internal"
	"os"
)

type Experiment struct {
	FileFormatVersion string
	Title             string
	Input             internal.Input
	FieldNames        []string
	ExcludeFieldNames []string
	Aggregators       []internal.Aggregator
	Goals             []*dexpr.Expr
	SortOrder         []SortField
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

type experimentFile struct {
	FileFormatVersion     string
	Title                 string
	InputFilename         string
	FieldNames            []string
	ExcludeFieldNames     []string
	IsFirstLineFieldNames bool
	Separator             string
	Aggregators           []experimentAggregator
	Goals                 []string
	SortOrder             []experimentSortField
}

type experimentAggregator struct {
	Name     string
	Function string
	Arg      string
}

type experimentSortField struct {
	AggregatorName string
	Direction      string
}

type ErrInvalidField struct {
	FieldName string
	Value     string
	Err       error
}

func (e *ErrInvalidField) Error() string {
	return fmt.Sprintf("Field: %q has Value: %q - %s", e.FieldName, e.Value, e.Err)
}

func LoadExperiment(filename string) (*Experiment, error) {
	var f *os.File
	var e experimentFile
	var experiment *Experiment
	var err error

	f, err = os.Open(filename)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(f)
	if err = dec.Decode(&e); err != nil {
		return nil, err
	}
	err = checkExperimentValid(e)
	if err != nil {
		return nil, err
	}
	experiment, err = makeExperiment(e)
	return experiment, err
}

func checkExperimentValid(e experimentFile) error {
	if e.FileFormatVersion == "" {
		return &ErrInvalidField{"fileFormatVersion", e.FileFormatVersion,
			errors.New("Must have a valid version number")}
	}
	// TODO: Test this more fully
	if len(e.FieldNames) < 2 {
		return &ErrInvalidField{"fieldNames",
			fmt.Sprintf("%q", e.FieldNames),
			errors.New("Must specify at least two field names")}
	}

	if len(e.Separator) != 1 {
		return &ErrInvalidField{"separator",
			fmt.Sprintf("%q", e.Separator),
			errors.New("Must contain one character only")}
	}
	return nil
}

func makeExperiment(e experimentFile) (*Experiment, error) {
	var goals []*dexpr.Expr
	var aggregators []internal.Aggregator
	var sortOrder []SortField
	var input internal.Input
	var err error
	goals, err = makeGoals(e.Goals)
	if err != nil {
		return nil, err
	}
	aggregators, err = makeAggregators(e.Aggregators)
	if err != nil {
		return nil, err
	}

	sortOrder, err = makeSortOrder(e.SortOrder)
	if err != nil {
		return nil, err
	}

	input, err = internal.NewCsvInput(e.FieldNames, e.InputFilename,
		rune(e.Separator[0]), e.IsFirstLineFieldNames)
	if err != nil {
		return nil, err
	}

	return &Experiment{
		FileFormatVersion: e.FileFormatVersion,
		Title:             e.Title,
		Input:             input,
		FieldNames:        e.FieldNames,
		ExcludeFieldNames: e.ExcludeFieldNames,
		Aggregators:       aggregators,
		Goals:             goals,
		SortOrder:         sortOrder,
	}, nil
}

func makeGoal(expr string) (*dexpr.Expr, error) {
	r, err := dexpr.New(expr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Can't make goal: %s", err))
	}
	return r, nil
}

func makeGoals(exprs []string) ([]*dexpr.Expr, error) {
	var err error
	r := make([]*dexpr.Expr, len(exprs))
	for i, s := range exprs {
		r[i], err = makeGoal(s)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

func makeAggregator(name, aggType, arg string) (internal.Aggregator, error) {
	var r internal.Aggregator
	var err error
	switch aggType {
	case "calc":
		r, err = internal.NewCalcAggregator(name, arg)
		return r, err
	case "count":
		r, err = internal.NewCountAggregator(name, arg)
		return r, err
	default:
		err = errors.New("Unrecognized aggregator")
	}
	if err != nil {
		// TODO: Make custome error type
		err = errors.New(fmt.Sprintf("Can't make aggregator: %s", err))
	}
	return r, err
}

func makeAggregators(
	eAggregators []experimentAggregator) ([]internal.Aggregator, error) {
	var err error
	r := make([]internal.Aggregator, len(eAggregators))
	for i, ea := range eAggregators {
		r[i], err = makeAggregator(ea.Name, ea.Function, ea.Arg)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

func makeSortOrder(eSortOrder []experimentSortField) ([]SortField, error) {
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
			err := errors.New(fmt.Sprintf("Invalid sort direction: %s, for field: %s",
				direction, field))
			return r, err
		}
	}
	return r, nil
}
