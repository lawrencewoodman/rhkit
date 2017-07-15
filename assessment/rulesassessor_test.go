package assessment

import (
	"fmt"
	"github.com/lawrencewoodman/ddataset/dcsv"
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/experiment"
	"github.com/vlifesystems/rhkit/internal/testhelpers"
	"github.com/vlifesystems/rhkit/rule"
	"path/filepath"
	"testing"
)

func TestAssessRules(t *testing.T) {
	rules := []rule.Rule{
		rule.NewGEFV("band", dlit.MustNew(5)),
		rule.NewGEFV("band", dlit.MustNew(4)),
		rule.NewGEFV("cost", dlit.MustNew(1.3)),
	}
	aggregators := []*experiment.AggregatorDesc{
		&experiment.AggregatorDesc{"numIncomeGt2", "count", "income > 2"},
		&experiment.AggregatorDesc{"numBandGt4", "count", "band > 4"},
	}
	goals := []string{
		"numIncomeGt2 == 1",
		"numIncomeGt2 == 2",
		"numIncomeGt2 == 3",
		"numIncomeGt2 == 4",
		"numBandGt4 == 1",
		"numBandGt4 == 2",
		"numBandGt4 == 3",
		"numBandGt4 == 4",
	}
	fieldNames := []string{"income", "cost", "band"}
	records := [][]string{
		[]string{"3", "4.5", "4"},
		[]string{"3", "3.2", "7"},
		[]string{"2", "1.2", "4"},
		[]string{"0", "0", "9"},
	}
	dataset := testhelpers.NewLiteralDataset(fieldNames, records)
	experimentDesc := &experiment.ExperimentDesc{
		Title:       "",
		Dataset:     dataset,
		RuleFields:  []string{"income", "cost", "band"},
		Aggregators: aggregators,
		Goals:       goals,
		SortOrder:   []*experiment.SortDesc{},
	}
	experiment := mustNewExperiment(experimentDesc)
	wantIsSorted := false
	wantIsRefined := false
	wantNumRecords := int64(len(records))
	wantRuleAssessments := []*RuleAssessment{
		&RuleAssessment{
			Rule: rule.NewGEFV("band", dlit.MustNew(5)),
			Aggregators: map[string]*dlit.Literal{
				"numMatches":     dlit.MustNew("2"),
				"percentMatches": dlit.MustNew("50"),
				"numIncomeGt2":   dlit.MustNew("1"),
				"numBandGt4":     dlit.MustNew("2"),
				"goalsScore":     dlit.MustNew(1.001),
			},
			Goals: []*GoalAssessment{
				&GoalAssessment{"numIncomeGt2 == 1", true},
				&GoalAssessment{"numIncomeGt2 == 2", false},
				&GoalAssessment{"numIncomeGt2 == 3", false},
				&GoalAssessment{"numIncomeGt2 == 4", false},
				&GoalAssessment{"numBandGt4 == 1", false},
				&GoalAssessment{"numBandGt4 == 2", true},
				&GoalAssessment{"numBandGt4 == 3", false},
				&GoalAssessment{"numBandGt4 == 4", false},
			},
		},
		&RuleAssessment{
			Rule: rule.NewGEFV("band", dlit.MustNew(4)),
			Aggregators: map[string]*dlit.Literal{
				"numMatches":     dlit.MustNew("4"),
				"percentMatches": dlit.MustNew("100"),
				"numIncomeGt2":   dlit.MustNew("2"),
				"numBandGt4":     dlit.MustNew("2"),
				"goalsScore":     dlit.MustNew(0.002),
			},
			Goals: []*GoalAssessment{
				&GoalAssessment{"numIncomeGt2 == 1", false},
				&GoalAssessment{"numIncomeGt2 == 2", true},
				&GoalAssessment{"numIncomeGt2 == 3", false},
				&GoalAssessment{"numIncomeGt2 == 4", false},
				&GoalAssessment{"numBandGt4 == 1", false},
				&GoalAssessment{"numBandGt4 == 2", true},
				&GoalAssessment{"numBandGt4 == 3", false},
				&GoalAssessment{"numBandGt4 == 4", false},
			},
		},
		&RuleAssessment{
			Rule: rule.NewGEFV("cost", dlit.MustNew(1.3)),
			Aggregators: map[string]*dlit.Literal{
				"numMatches":     dlit.MustNew("2"),
				"percentMatches": dlit.MustNew("50"),
				"numIncomeGt2":   dlit.MustNew("2"),
				"numBandGt4":     dlit.MustNew("1"),
				"goalsScore":     dlit.MustNew(0.002),
			},
			Goals: []*GoalAssessment{
				&GoalAssessment{"numIncomeGt2 == 1", false},
				&GoalAssessment{"numIncomeGt2 == 2", true},
				&GoalAssessment{"numIncomeGt2 == 3", false},
				&GoalAssessment{"numIncomeGt2 == 4", false},
				&GoalAssessment{"numBandGt4 == 1", true},
				&GoalAssessment{"numBandGt4 == 2", false},
				&GoalAssessment{"numBandGt4 == 3", false},
				&GoalAssessment{"numBandGt4 == 4", false},
			},
		},
	}
	gotAssessment, err := AssessRules(rules, experiment)
	if err != nil {
		t.Errorf("AssessRules(%v, %v) - err: %v", rules, experiment, err)
	}

	assessmentsMatch := areAssessmentsEqv(
		gotAssessment,
		wantNumRecords,
		wantIsSorted,
		wantIsRefined,
		wantRuleAssessments,
	)
	if !assessmentsMatch {
		t.Errorf("AssessRules(%v, %v)\nassessments don't match\n - got: %v\n - wantRuleAssessments: %v, wantNumRecords: %d, wantIsSorted: %t, wantIsRefined: %t\n",
			rules, experiment, gotAssessment, wantRuleAssessments,
			wantNumRecords, wantIsSorted, wantIsRefined)
	}
}

func TestAssessRules_errors(t *testing.T) {
	cases := []struct {
		rules       []rule.Rule
		aggregators []*experiment.AggregatorDesc
		goals       []string
		wantErr     error
	}{
		{[]rule.Rule{rule.NewGEFV("hand", dlit.MustNew(3))},
			[]*experiment.AggregatorDesc{
				&experiment.AggregatorDesc{"numIncomeGt2", "count", "income > 2"},
			},
			[]string{"numIncomeGt2 == 1"},
			rule.InvalidRuleError{Rule: rule.NewGEFV("hand", dlit.MustNew(3))},
		},
		{[]rule.Rule{rule.NewGEFV("band", dlit.MustNew(3))},
			[]*experiment.AggregatorDesc{
				&experiment.AggregatorDesc{"numIncomeGt2", "count", "bincome > 2"},
			},
			[]string{"numIncomeGt2 == 1"},
			dexpr.InvalidExprError{
				Expr: "bincome > 2",
				Err:  dexpr.VarNotExistError("bincome"),
			},
		},
		{[]rule.Rule{rule.NewGEFV("band", dlit.MustNew(3))},
			[]*experiment.AggregatorDesc{
				&experiment.AggregatorDesc{"numIncomeGt2", "count", "income > 2"},
			},
			[]string{"numIncomeGt == 1"},
			dexpr.InvalidExprError{
				Expr: "numIncomeGt == 1",
				Err:  dexpr.VarNotExistError("numIncomeGt"),
			},
		},
	}
	fieldNames := []string{"income", "cost", "band"}
	records := [][]string{
		[]string{"3", "4.5", "4"},
	}
	dataset := testhelpers.NewLiteralDataset(fieldNames, records)
	for _, c := range cases {
		experimentDesc := &experiment.ExperimentDesc{
			Title:       "",
			Dataset:     dataset,
			RuleFields:  []string{"income", "cost", "band"},
			Aggregators: c.aggregators,
			Goals:       c.goals,
			SortOrder:   []*experiment.SortDesc{},
		}
		experiment := mustNewExperiment(experimentDesc)
		_, err := AssessRules(c.rules, experiment)
		if err == nil || err.Error() != c.wantErr.Error() {
			t.Errorf("AssessRules(%v, %v) - err: %s, wantErr: %s",
				c.rules, experiment, err, c.wantErr)
		}
	}
}

/*************************
       Benchmarks
*************************/
func BenchmarkAssessRules(b *testing.B) {
	b.StopTimer()
	var numRules int64 = 3000
	rules := make([]rule.Rule, numRules)
	for i := int64(0); i < numRules; i++ {
		if i%2 == 0 {
			rules[i] = rule.NewGEFV("age", dlit.MustNew(i%50))
		} else {
			rules[i] = rule.NewGEFV("day", dlit.MustNew(i%20))
		}
	}

	fieldNames := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y"}
	experimentDesc := &experiment.ExperimentDesc{
		Title: "",
		Dataset: dcsv.New(
			filepath.Join("fixtures", "bank_big.csv"),
			true,
			rune(';'),
			fieldNames,
		),
		RuleFields: []string{"age", "job", "default",
			"balance", "housing", "loan", "contact", "day", "month", "duration",
			"campaign", "pdays", "previous", "poutcome",
		},
		Aggregators: []*experiment.AggregatorDesc{
			&experiment.AggregatorDesc{"numMarried", "count", "marital == \"married\""},
			&experiment.AggregatorDesc{"numSignedUp", "count", "y == \"yes\""},
			&experiment.AggregatorDesc{"cost", "calc", "numMatches * 4.5"},
			&experiment.AggregatorDesc{"income", "calc", "numSignedUp * 24"},
			&experiment.AggregatorDesc{"profit", "calc", "income - cost"},
		},
		Goals: []string{
			"profit > 0",
			"numSignedUp > 3",
			"numMarried > 2",
		},
		SortOrder: []*experiment.SortDesc{
			&experiment.SortDesc{"profit", "descending"},
			&experiment.SortDesc{"numSignedUp", "descending"},
			&experiment.SortDesc{"cost", "ascending"},
			&experiment.SortDesc{"numMatches", "descending"},
			&experiment.SortDesc{"percentMatches", "descending"},
			&experiment.SortDesc{"goalsScore", "descending"},
		},
	}
	experiment := mustNewExperiment(experimentDesc)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		_, err := AssessRules(rules, experiment)
		if err != nil {
			b.Errorf("AssessRules(%v, %v) - err: %v", rules, experiment, err)
		}
	}
}

/******************************
 *  Helper functions
 ******************************/

func mustNewExperiment(ed *experiment.ExperimentDesc) *experiment.Experiment {
	e, err := experiment.New(ed)
	if err != nil {
		panic(fmt.Sprintf("Can't create Experiment: %s", err))
	}
	return e
}

// Are the assessments equivalent.  The ruleAssessments must match
// but don't have to be in the same order if both assessments are
// unsorted. If both are unsorted then this will sort the assessments
func areAssessmentsEqv(
	got *Assessment,
	wantNumRecords int64,
	wantIsSorted bool,
	wantIsRefined bool,
	wantRuleAssessments []*RuleAssessment,
) bool {
	if got.NumRecords != wantNumRecords {
		return false
	}
	if got.IsSorted() != wantIsSorted {
		return false
	}
	if got.IsRefined() != wantIsRefined {
		return false
	}
	for _, gotRuleAssesment := range got.RuleAssessments {
		found := false
		for _, wantRuleAssessment := range wantRuleAssessments {
			if gotRuleAssesment.IsEqual(wantRuleAssessment) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
