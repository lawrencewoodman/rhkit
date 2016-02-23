package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// Ensure that correct number is returned
func TestLoadExperiment(t *testing.T) {
	expectedExperiments := []Experiment{
		Experiment{},
		Experiment{
			FileFormatVersion: "0.1",
			Title:             "This is a jolly nice title",
			InputFilename:     "/tmp/bank.csv",
			FieldNames: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y"},
			ExcludeFieldNames:     []string{"education"},
			IsFirstLineFieldNames: true,
			Separator:             ";",
			Aggregators: []ExperimentAggregator{
				ExperimentAggregator{"numSignedUp", "count", "y eq \"yes\""},
				ExperimentAggregator{"cost", "calc", "numMatches * 4.5"},
				ExperimentAggregator{"income", "calc", "numSignedUp * 24"},
				ExperimentAggregator{"profit", "calc", "income - cost"}},
			Goals: []string{"profit > 0"},
			SortOrder: []SortField{
				SortField{"profit", "descending"},
				SortField{"numSignedUp", "descending"}},
		},
	}
	cases := []struct {
		filename string
		want     Experiment
		wantErr  error
	}{
		{"missingfile.json", expectedExperiments[0],
			&os.PathError{"open", "missingfile.json",
				errors.New("no such file or directory")}},
		{filepath.Join("fixtures", "missingFileFormatVersion.json"),
			expectedExperiments[0],
			&ErrInvalidField{"fileFormatVersion", "",
				errors.New("Must have a valid version number")}},
		{filepath.Join("fixtures", "singleFieldName.json"),
			expectedExperiments[0],
			&ErrInvalidField{"fieldNames", "[\"age\"]",
				errors.New("Must specify at least two field names")}},
		{filepath.Join("fixtures", "bank.json"),
			expectedExperiments[1], nil},
	}
	for _, c := range cases {
		got, err := LoadExperiment(c.filename)
		if !errorMatch(c.wantErr, err) ||
			(err == nil && !experimentMatch(got, c.want)) {
			t.Errorf("LoadExperiment(%q) err: %q, wantErr: %q - got: %q, want: %q",
				c.filename, err, c.wantErr, got, c.want)
		}
	}
}

/***********************
   Helper functions
************************/

func experimentMatch(e1 Experiment, e2 Experiment) bool {
	if e1.FileFormatVersion != e2.FileFormatVersion ||
		e1.Title != e2.Title ||
		e1.InputFilename != e2.InputFilename ||
		e1.IsFirstLineFieldNames != e2.IsFirstLineFieldNames ||
		e1.Separator != e2.Separator ||
		!areStringArraysEqual(e1.FieldNames, e2.FieldNames) ||
		!areStringArraysEqual(e1.ExcludeFieldNames, e2.ExcludeFieldNames) ||
		!areStringArraysEqual(e1.Goals, e2.Goals) ||
		!areAggregatorsEqual(e1.Aggregators, e2.Aggregators) ||
		!areSortOrdersEqual(e1.SortOrder, e2.SortOrder) {
		return false
	}
	return true
}

func areStringArraysEqual(a1 []string, a2 []string) bool {
	if len(a1) != len(a2) {
		return false
	}
	for i, e := range a1 {
		if e != a2[i] {
			return false
		}
	}
	return true
}

func areAggregatorsEqual(a1 []ExperimentAggregator,
	a2 []ExperimentAggregator) bool {
	if len(a1) != len(a2) {
		return false
	}
	for i, e := range a1 {
		if e != a2[i] {
			return false
		}
	}
	return true
}

func areSortOrdersEqual(so1 []SortField, so2 []SortField) bool {
	if len(so1) != len(so2) {
		return false
	}
	for i, e := range so1 {
		if e != so2[i] {
			return false
		}
	}
	return true
}
