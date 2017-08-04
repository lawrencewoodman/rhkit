package rhkit

import (
	"fmt"
	"github.com/lawrencewoodman/ddataset/dcsv"
	"github.com/vlifesystems/rhkit/aggregators"
	"github.com/vlifesystems/rhkit/assessment"
	"github.com/vlifesystems/rhkit/goal"
	"github.com/vlifesystems/rhkit/rule"
	"path/filepath"
	"testing"
)

func TestProcess(t *testing.T) {
	fields := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y"}
	dataset := dcsv.New(
		filepath.Join("fixtures", "bank.csv"),
		true,
		rune(';'),
		fields,
	)
	ruleFields := []string{"age", "job", "marital", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y",
	}
	aggregatorDescs := []*aggregators.Desc{
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
	}
	goalExprs := []string{"profit > 0"}
	sortOrderDescs := []assessment.SortDesc{
		{"profit", "descending"},
		{"numSignedUp", "descending"},
		{"goalsScore", "descending"},
	}
	goals, err := goal.MakeGoals(goalExprs)
	if err != nil {
		t.Fatalf("MakeGoals: %s", err)
	}
	aggregators, err := aggregators.MakeSpecs(dataset.Fields(), aggregatorDescs)
	if err != nil {
		t.Fatalf("MakeSpecs: %s", err)
	}
	sortOrder, err := assessment.MakeSortOrders(aggregators, sortOrderDescs)
	if err != nil {
		t.Fatalf("MakeSortOrders: %s", err)
	}
	ruleComplexity := rule.Complexity{Arithmetic: true}
	rules := []rule.Rule{}

	for maxNumRules := 0; maxNumRules < 1500; maxNumRules += 100 {
		maxNumRules := maxNumRules
		t.Run(fmt.Sprintf("maxNumRules %d", maxNumRules), func(t *testing.T) {
			t.Parallel()
			ass, err := Process(
				dataset,
				ruleFields,
				ruleComplexity,
				aggregators,
				goals,
				sortOrder,
				rules,
				maxNumRules,
			)
			if err != nil {
				t.Errorf("Process: %s", err)
				return
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
	fields := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y"}
	dataset := dcsv.New(
		filepath.Join("fixtures", "bank.csv"),
		true,
		rune(';'),
		fields,
	)
	ruleFields := []string{"age", "job", "marital", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y",
	}
	aggregatorDescs := []*aggregators.Desc{
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
	}
	goalExprs := []string{"profit > 0"}
	sortOrderDescs := []assessment.SortDesc{
		{"profit", "descending"},
		{"numSignedUp", "descending"},
		{"goalsScore", "descending"},
	}
	ruleExprs := []string{
		"age > 30",
		"age > 30 && duration > 79",
		"age > 30 && pdays > 5",
		"age <= 19 || age >= 37",
		"month == \"may\"",
		"month == \"unknown\"",
	}
	ruleComplexity := rule.Complexity{Arithmetic: true}
	maxNumRules := 50
	wantRules := []string{
		"age > 30",
		"age > 30 && duration > 79",
		"age <= 19 || age >= 37",
		"month == \"may\"",
	}

	goals, err := goal.MakeGoals(goalExprs)
	if err != nil {
		t.Fatalf("MakeGoals: %s", err)
	}
	aggregators, err := aggregators.MakeSpecs(dataset.Fields(), aggregatorDescs)
	if err != nil {
		t.Fatalf("MakeSpecs: %s", err)
	}
	sortOrder, err := assessment.MakeSortOrders(aggregators, sortOrderDescs)
	if err != nil {
		t.Fatalf("MakeSortOrders: %s", err)
	}
	rules, err := rule.MakeDynamicRules(ruleExprs)
	if err != nil {
		t.Fatalf("MakeDynamicRules: %s", err)
	}
	ass, err := Process(
		dataset,
		ruleFields,
		ruleComplexity,
		aggregators,
		goals,
		sortOrder,
		rules,
		maxNumRules,
	)
	if err != nil {
		t.Errorf("Process: %s", err)
	}

	gotRules := ass.Rules()
	numRules := len(gotRules)
	if numRules > maxNumRules || (maxNumRules > 0 && numRules < 1) {
		t.Errorf("Process - numRules: %d, maxNumRules: %d",
			numRules, maxNumRules)
	}

	for _, wantRule := range wantRules {
		foundRule := false
		for _, r := range gotRules {
			if wantRule == r.String() {
				foundRule = true
			}
		}
		if !foundRule {
			t.Errorf("Process: couldn't find rule: %s", wantRule)
		}
	}
}
