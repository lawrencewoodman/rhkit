package rhkit

import (
	"fmt"
	"github.com/lawrencewoodman/ddataset/dcsv"
	"github.com/vlifesystems/rhkit/aggregator"
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
	aggregatorDescs := []*aggregator.Desc{
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
	aggregators, err := aggregator.MakeSpecs(dataset.Fields(), aggregatorDescs)
	if err != nil {
		t.Fatalf("MakeSpecs: %s", err)
	}
	sortOrder, err := assessment.MakeSortOrders(aggregators, sortOrderDescs)
	if err != nil {
		t.Fatalf("MakeSortOrders: %s", err)
	}
	rules := []rule.Rule{}

	cases := []struct {
		opts            Options
		wantMinNumRules int
		wantMaxNumRules int
	}{
		{opts: Options{MaxNumRules: 0, RuleFields: []string{}},
			wantMinNumRules: 1,
			wantMaxNumRules: 1,
		},
		{opts: Options{MaxNumRules: 1, RuleFields: []string{}},
			wantMinNumRules: 1,
			wantMaxNumRules: 1,
		},
		{opts: Options{MaxNumRules: 1500, RuleFields: []string{}},
			wantMinNumRules: 1,
			wantMaxNumRules: 1,
		},
		{opts: Options{MaxNumRules: 0, RuleFields: ruleFields},
			wantMinNumRules: 1,
			wantMaxNumRules: 1,
		},
		{opts: Options{MaxNumRules: 1, RuleFields: ruleFields},
			wantMinNumRules: 1,
			wantMaxNumRules: 1,
		},
		{opts: Options{MaxNumRules: 100, RuleFields: ruleFields},
			wantMinNumRules: 100,
			wantMaxNumRules: 100,
		},
		{opts: Options{MaxNumRules: 500, RuleFields: ruleFields},
			wantMinNumRules: 500,
			wantMaxNumRules: 500,
		},
		{opts: Options{MaxNumRules: 3000, RuleFields: ruleFields},
			wantMinNumRules: 1400,
			wantMaxNumRules: 1600,
		},
		{opts: Options{
			MaxNumRules:             3000,
			RuleFields:              ruleFields,
			GenerateArithmeticRules: true,
		},
			wantMinNumRules: 1100,
			wantMaxNumRules: 1400,
		},
	}
	for i, c := range cases {
		opts := c.opts
		i := i
		wantMinNumRules := c.wantMinNumRules
		wantMaxNumRules := c.wantMaxNumRules
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()
			ass, err := Process(dataset, aggregators, goals, sortOrder, rules, opts)
			if err != nil {
				t.Fatalf("(%d) Process: %s", i, err)
			}
			numRules := len(ass.Rules())
			if numRules > wantMaxNumRules || numRules < wantMinNumRules {
				t.Errorf("(%d) Process - opts: %v - gotNumRules: %d, wantMinNumRules: %d wantMaxNumRules: %d",
					i, opts, numRules, wantMinNumRules, wantMaxNumRules)
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
	aggregatorDescs := []*aggregator.Desc{
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
	maxNumRules := 50
	wantRules := []string{
		"age > 30",
		"age > 30 && duration > 79",
		"age <= 19 || age >= 37",
	}

	goals, err := goal.MakeGoals(goalExprs)
	if err != nil {
		t.Fatalf("MakeGoals: %s", err)
	}
	aggregators, err := aggregator.MakeSpecs(dataset.Fields(), aggregatorDescs)
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
	opts := Options{
		MaxNumRules:             maxNumRules,
		RuleFields:              []string{},
		GenerateArithmeticRules: true,
	}
	ass, err := Process(dataset, aggregators, goals, sortOrder, rules, opts)
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
