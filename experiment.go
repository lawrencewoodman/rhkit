/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// TODO: Consider making an internal structure and external structure
// This would allow Separator to be a rune for example
type Experiment struct {
	FileFormatVersion     string
	Title                 string
	InputFilename         string
	FieldNames            []string
	ExcludeFieldNames     []string
	IsFirstLineFieldNames bool
	Separator             string
	Aggregators           []ExperimentAggregator
	Goals                 []string
	SortOrder             []SortField
}

type ExperimentAggregator struct {
	Name     string
	Function string
	Arg      string
}

type SortField struct {
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

func LoadExperiment(filename string) (Experiment, error) {
	var f *os.File
	var e Experiment
	var err error

	f, err = os.Open(filename)
	if err != nil {
		return e, err
	}

	dec := json.NewDecoder(f)
	if err = dec.Decode(&e); err != nil {
		return e, err
	}
	err = checkExperimentValid(e)
	return e, err
}

func checkExperimentValid(e Experiment) error {
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
	return nil
}
