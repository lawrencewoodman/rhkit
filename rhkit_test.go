package rhkit

import (
	"fmt"
	"github.com/lawrencewoodman/ddataset/dcsv"
	"github.com/vlifesystems/rhkit/aggregators"
	"github.com/vlifesystems/rhkit/experiment"
	"path/filepath"
	"testing"
)

func TestProcess(t *testing.T) {
	fieldNames := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y"}
	experimentDesc := &experiment.ExperimentDesc{
		Dataset: dcsv.New(
			filepath.Join("fixtures", "bank.csv"),
			true,
			rune(';'),
			fieldNames,
		),
		RuleFields: []string{"age", "job", "marital", "default",
			"balance", "housing", "loan", "contact", "day", "month", "duration",
			"campaign", "pdays", "previous", "poutcome", "y",
		},
		Aggregators: []*aggregators.Desc{
			{"numSignedUp", "count", "y == \"yes\""},
			{"cost", "calc", "numMatches * 4.5"},
			{"income", "calc", "numSignedUp * 24"},
			{"profit", "calc", "income - cost"},
			{"oddFigure", "sum", "balance - age"},
			{
				"percentMarried",
				"precision",
				"marital == \"married\"",
			},
		},
		Goals: []string{"profit > 0"},
		SortOrder: []*experiment.SortDesc{
			{"profit", "descending"},
			{"numSignedUp", "descending"},
			{"goalsScore", "descending"},
		},
	}
	experiment, err := experiment.New(experimentDesc)
	if err != nil {
		t.Fatalf("experiment.New(%s) - err: %s", experimentDesc, err)
	}
	for maxNumRules := 0; maxNumRules < 1500; maxNumRules += 100 {
		maxNumRules := maxNumRules
		t.Run(fmt.Sprintf("maxNumRules %d", maxNumRules), func(t *testing.T) {
			t.Parallel()
			ass, err := Process(experiment, maxNumRules)
			if err != nil {
				t.Errorf("Process: %s", err)
			}
			numRules := len(ass.Rules())
			if numRules > maxNumRules || (maxNumRules > 0 && numRules < 1) {
				t.Errorf("Process - numRules: %d, maxNumRules: %d",
					numRules, maxNumRules)
			}
		})
	}
}

func TestProcess_user_rules(t *testing.T) {
	fieldNames := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y"}
	experimentDesc := &experiment.ExperimentDesc{
		Dataset: dcsv.New(
			filepath.Join("fixtures", "bank.csv"),
			true,
			rune(';'),
			fieldNames,
		),
		RuleFields: []string{"age", "job", "marital", "default",
			"balance", "housing", "loan", "contact", "day", "month", "duration",
			"campaign", "pdays", "previous", "poutcome", "y",
		},
		Aggregators: []*aggregators.Desc{
			{"numSignedUp", "count", "y == \"yes\""},
			{"cost", "calc", "numMatches * 4.5"},
			{"income", "calc", "numSignedUp * 24"},
			{"profit", "calc", "income - cost"},
			{"oddFigure", "sum", "balance - age"},
			{
				"percentMarried",
				"precision",
				"marital == \"married\"",
			},
		},
		Goals: []string{"profit > 0"},
		SortOrder: []*experiment.SortDesc{
			{"profit", "descending"},
			{"numSignedUp", "descending"},
			{"goalsScore", "descending"},
		},
		Rules: []string{
			"age > 30",
			"age > 30 && duration > 79",
			"age > 30 && pdays > 5",
			"age <= 19 || age >= 37",
			"month == \"may\"",
			"month == \"unknown\"",
		},
	}
	wantRules := []string{
		"age > 30",
		"age > 30 && duration > 79",
		"age <= 19 || age >= 37",
		"month == \"may\"",
	}
	experiment, err := experiment.New(experimentDesc)
	if err != nil {
		t.Fatalf("experiment.New(%s) - err: %s", experimentDesc, err)
	}
	maxNumRules := 50
	ass, err := Process(experiment, maxNumRules)
	if err != nil {
		t.Errorf("Process: %s", err)
	}

	rules := ass.Rules()
	numRules := len(rules)
	if numRules > maxNumRules || (maxNumRules > 0 && numRules < 1) {
		t.Errorf("Process - numRules: %d, maxNumRules: %d",
			numRules, maxNumRules)
	}

	for _, wantRule := range wantRules {
		foundRule := false
		for _, r := range rules {
			if wantRule == r.String() {
				foundRule = true
			}
		}
		if !foundRule {
			t.Errorf("Process: couldn't find rule: %s", wantRule)
		}
	}
}
