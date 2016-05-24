/*
 *   Integration test for package
 */
package integrate

import (
	"github.com/vlifesystems/rulehunter/csvinput"
	"github.com/vlifesystems/rulehunter/experiment"
	"github.com/vlifesystems/rulehunter/reduceinput"
	"path/filepath"
	"testing"
)

func TestAll_full(t *testing.T) {
	fieldNames := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y"}
	records, err := csvinput.New(
		fieldNames,
		filepath.Join("..", "..", "fixtures", "bank.csv"),
		rune(';'),
		true,
	)
	if err != nil {
		t.Errorf("csvInput.New() - err: %s", err)
		return
	}
	experimentDesc := &experiment.ExperimentDesc{
		Title:         "This is a jolly nice title",
		Input:         records,
		ExcludeFields: []string{"education"},
		Aggregators: []*experiment.AggregatorDesc{
			&experiment.AggregatorDesc{"numSignedUp", "count", "y == \"yes\""},
			&experiment.AggregatorDesc{"cost", "calc", "numMatches * 4.5"},
			&experiment.AggregatorDesc{"income", "calc", "numSignedUp * 24"},
			&experiment.AggregatorDesc{"profit", "calc", "income - cost"},
			&experiment.AggregatorDesc{"oddFigure", "sum", "balance - age"},
			&experiment.AggregatorDesc{
				"percentMarried",
				"percent",
				"marital == \"married\"",
			},
		},
		Goals: []string{"profit > 0"},
		SortOrder: []*experiment.SortDesc{
			&experiment.SortDesc{"profit", "descending"},
			&experiment.SortDesc{"numSignedUp", "descending"},
			&experiment.SortDesc{"goalsScore", "descending"},
		},
	}
	experiment, err := experiment.New(experimentDesc)
	if err != nil {
		t.Errorf("experiment.New(%s) - err: %s", experimentDesc, err)
		return
	}
	defer experiment.Close()
	if err = ProcessInput(experiment); err != nil {
		t.Errorf("ProcessInput() - err: %s", err)
		return
	}
}

func TestAll_reduced(t *testing.T) {
	fieldNames := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y"}
	numRecords := 5

	input, err := csvinput.New(
		fieldNames,
		filepath.Join("..", "..", "fixtures", "bank.csv"),
		rune(';'),
		true,
	)
	if err != nil {
		t.Errorf("csvInput.New() - err: %s", err)
		return
	}

	records, err := reduceinput.New(input, numRecords)
	if err != nil {
		t.Errorf("reduceInput.New() - err: %s", err)
		return
	}

	experimentDesc := &experiment.ExperimentDesc{
		Title:         "This is a jolly nice title",
		Input:         records,
		ExcludeFields: []string{"education"},
		Aggregators: []*experiment.AggregatorDesc{
			&experiment.AggregatorDesc{"numSignedUp", "count", "y == \"yes\""},
			&experiment.AggregatorDesc{"cost", "calc", "numMatches * 4.5"},
			&experiment.AggregatorDesc{"income", "calc", "numSignedUp * 24"},
			&experiment.AggregatorDesc{"profit", "calc", "income - cost"},
			&experiment.AggregatorDesc{"oddFigure", "sum", "balance - age"},
			&experiment.AggregatorDesc{
				"percentMarried",
				"percent",
				"marital == \"married\"",
			},
		},
		Goals: []string{"profit > 0"},
		SortOrder: []*experiment.SortDesc{
			&experiment.SortDesc{"profit", "descending"},
			&experiment.SortDesc{"numSignedUp", "descending"},
			&experiment.SortDesc{"goalsScore", "descending"},
		},
	}
	experiment, err := experiment.New(experimentDesc)
	if err != nil {
		t.Errorf("experiment.New(%s) - err: %s", experimentDesc, err)
		return
	}

	defer experiment.Close()
	if err = ProcessInput(experiment); err != nil {
		t.Errorf("ProcessInput() - err: %s", err)
		return
	}
}
