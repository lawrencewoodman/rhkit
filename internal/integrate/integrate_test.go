/*
 *   Integration test for package
 */
package integrate

import (
	"github.com/lawrencewoodman/ddataset/dcsv"
	"github.com/vlifesystems/rulehunter/experiment"
	"path/filepath"
	"testing"
)

func TestAll(t *testing.T) {
	fieldNames := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y"}
	experimentDesc := &experiment.ExperimentDesc{
		Title: "This is a jolly nice title",
		Dataset: dcsv.New(
			filepath.Join("..", "..", "fixtures", "bank.csv"),
			true,
			rune(';'),
			fieldNames,
		),
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
