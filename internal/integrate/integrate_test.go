/*
 *   Integration test for package
 */
package integrate

import (
	"github.com/vlifesystems/rulehunter/csvdataset"
	"github.com/vlifesystems/rulehunter/experiment"
	"github.com/vlifesystems/rulehunter/reducedataset"
	"path/filepath"
	"testing"
)

func TestAll_full(t *testing.T) {
	fieldNames := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y"}
	dataset, err := csvdataset.New(
		fieldNames,
		filepath.Join("..", "..", "fixtures", "bank.csv"),
		rune(';'),
		true,
	)
	if err != nil {
		t.Errorf("csvDataset.New() - err: %s", err)
		return
	}
	experimentDesc := &experiment.ExperimentDesc{
		Title:         "This is a jolly nice title",
		Dataset:       dataset,
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
	if err = ProcessDataset(experiment); err != nil {
		t.Errorf("ProcessDataset() - err: %s", err)
		return
	}
}

func TestAll_reduced(t *testing.T) {
	fieldNames := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y"}
	numRecords := 5

	dataset, err := csvdataset.New(
		fieldNames,
		filepath.Join("..", "..", "fixtures", "bank.csv"),
		rune(';'),
		true,
	)
	if err != nil {
		t.Errorf("csvDataset.New() - err: %s", err)
		return
	}

	reducedDataset := reducedataset.New(dataset, numRecords)

	experimentDesc := &experiment.ExperimentDesc{
		Title:         "This is a jolly nice title",
		Dataset:       reducedDataset,
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

	if err = ProcessDataset(experiment); err != nil {
		t.Errorf("ProcessDataset() - err: %s", err)
		return
	}
}
