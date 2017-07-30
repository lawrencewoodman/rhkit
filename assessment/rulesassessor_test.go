package assessment

import (
	"github.com/lawrencewoodman/ddataset/dcsv"
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/aggregators"
	"github.com/vlifesystems/rhkit/goal"
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
	aggregatorDescs := []*aggregators.Desc{
		{"numIncomeGt2", "count", "income > 2"},
		{"numBandGt4", "count", "band > 4"},
	}
	goalExprs := []string{
		"numIncomeGt2 == 1",
		"numIncomeGt2 == 2",
		"numIncomeGt2 == 3",
		"numIncomeGt2 == 4",
		"numBandGt4 == 1",
		"numBandGt4 == 2",
		"numBandGt4 == 3",
		"numBandGt4 == 4",
	}
	fields := []string{"income", "cost", "band"}
	records := [][]string{
		{"3", "4.5", "4"},
		{"3", "3.2", "7"},
		{"2", "1.2", "4"},
		{"0", "0", "9"},
	}
	dataset := testhelpers.NewLiteralDataset(fields, records)
	aggregatorSpecs, err := aggregators.MakeSpecs(fields, aggregatorDescs)
	if err != nil {
		t.Fatalf("MakeSpecs: %s", err)
	}
	goals, err := goal.MakeGoals(goalExprs)
	if err != nil {
		t.Fatalf("MakeGoals: %s", err)
	}
	wantIsSorted := false
	wantIsRefined := false
	wantNumRecords := int64(len(records))
	wantRuleAssessments := []*RuleAssessment{
		{
			Rule: rule.NewGEFV("band", dlit.MustNew(5)),
			Aggregators: map[string]*dlit.Literal{
				"numMatches":     dlit.MustNew("2"),
				"percentMatches": dlit.MustNew("50"),
				"numIncomeGt2":   dlit.MustNew("1"),
				"numBandGt4":     dlit.MustNew("2"),
				"goalsScore":     dlit.MustNew(1.001),
			},
			Goals: []*GoalAssessment{
				{"numIncomeGt2 == 1", true},
				{"numIncomeGt2 == 2", false},
				{"numIncomeGt2 == 3", false},
				{"numIncomeGt2 == 4", false},
				{"numBandGt4 == 1", false},
				{"numBandGt4 == 2", true},
				{"numBandGt4 == 3", false},
				{"numBandGt4 == 4", false},
			},
		},
		{
			Rule: rule.NewGEFV("band", dlit.MustNew(4)),
			Aggregators: map[string]*dlit.Literal{
				"numMatches":     dlit.MustNew("4"),
				"percentMatches": dlit.MustNew("100"),
				"numIncomeGt2":   dlit.MustNew("2"),
				"numBandGt4":     dlit.MustNew("2"),
				"goalsScore":     dlit.MustNew(0.002),
			},
			Goals: []*GoalAssessment{
				{"numIncomeGt2 == 1", false},
				{"numIncomeGt2 == 2", true},
				{"numIncomeGt2 == 3", false},
				{"numIncomeGt2 == 4", false},
				{"numBandGt4 == 1", false},
				{"numBandGt4 == 2", true},
				{"numBandGt4 == 3", false},
				{"numBandGt4 == 4", false},
			},
		},
		{
			Rule: rule.NewGEFV("cost", dlit.MustNew(1.3)),
			Aggregators: map[string]*dlit.Literal{
				"numMatches":     dlit.MustNew("2"),
				"percentMatches": dlit.MustNew("50"),
				"numIncomeGt2":   dlit.MustNew("2"),
				"numBandGt4":     dlit.MustNew("1"),
				"goalsScore":     dlit.MustNew(0.002),
			},
			Goals: []*GoalAssessment{
				{"numIncomeGt2 == 1", false},
				{"numIncomeGt2 == 2", true},
				{"numIncomeGt2 == 3", false},
				{"numIncomeGt2 == 4", false},
				{"numBandGt4 == 1", true},
				{"numBandGt4 == 2", false},
				{"numBandGt4 == 3", false},
				{"numBandGt4 == 4", false},
			},
		},
	}
	gotAssessment, err := AssessRules(dataset, rules, aggregatorSpecs, goals)
	if err != nil {
		t.Errorf("AssessRules: %v", err)
	}

	assessmentsMatch := areAssessmentsEqv(
		gotAssessment,
		wantNumRecords,
		wantIsSorted,
		wantIsRefined,
		wantRuleAssessments,
	)
	if !assessmentsMatch {
		t.Errorf("AssessRules: assessments don't match")
		t.Errorf("got: %v", gotAssessment)
		t.Errorf("wantRuleAssessments: %v, wantNumRecords: %d, wantIsSorted: %t, wantIsRefined: %t",
			wantRuleAssessments,
			wantNumRecords, wantIsSorted, wantIsRefined)
	}
}

func TestAssessRules_errors(t *testing.T) {
	cases := []struct {
		rules           []rule.Rule
		aggregatorDescs []*aggregators.Desc
		goalExprs       []string
		wantErr         error
	}{
		{[]rule.Rule{rule.NewGEFV("hand", dlit.MustNew(3))},
			[]*aggregators.Desc{
				{"numIncomeGt2", "count", "income > 2"},
			},
			[]string{"numIncomeGt2 == 1"},
			rule.InvalidRuleError{Rule: rule.NewGEFV("hand", dlit.MustNew(3))},
		},
		{[]rule.Rule{rule.NewGEFV("band", dlit.MustNew(3))},
			[]*aggregators.Desc{
				{"numIncomeGt2", "count", "bincome > 2"},
			},
			[]string{"numIncomeGt2 == 1"},
			dexpr.InvalidExprError{
				Expr: "bincome > 2",
				Err:  dexpr.VarNotExistError("bincome"),
			},
		},
		{[]rule.Rule{rule.NewGEFV("band", dlit.MustNew(3))},
			[]*aggregators.Desc{
				{"numIncomeGt2", "count", "income > 2"},
			},
			[]string{"numIncomeGt == 1"},
			dexpr.InvalidExprError{
				Expr: "numIncomeGt == 1",
				Err:  dexpr.VarNotExistError("numIncomeGt"),
			},
		},
	}
	fields := []string{"income", "cost", "band"}
	records := [][]string{
		{"3", "4.5", "4"},
	}
	dataset := testhelpers.NewLiteralDataset(fields, records)
	for i, c := range cases {
		aggregatorSpecs, err := aggregators.MakeSpecs(fields, c.aggregatorDescs)
		if err != nil {
			t.Fatalf("(%d) MakeSpecs: %s", i, err)
		}
		goals, err := goal.MakeGoals(c.goalExprs)
		if err != nil {
			t.Fatalf("(%d) MakeGoals: %s", i, err)
		}
		_, err = AssessRules(dataset, c.rules, aggregatorSpecs, goals)
		if err == nil || err.Error() != c.wantErr.Error() {
			t.Errorf("(%d) AssessRules - err: %s, wantErr: %s", i, err, c.wantErr)
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

	fields := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y"}
	dataset := dcsv.New(
		filepath.Join("fixtures", "bank_big.csv"),
		true,
		rune(';'),
		fields,
	)
	aggregatorDescs := []*aggregators.Desc{
		{"numMarried", "count", "marital == \"married\""},
		{"numSignedUp", "count", "y == \"yes\""},
		{"cost", "calc", "numMatches * 4.5"},
		{"income", "calc", "numSignedUp * 24"},
		{"profit", "calc", "income - cost"},
	}
	goalExprs := []string{
		"profit > 0",
		"numSignedUp > 3",
		"numMarried > 2",
	}
	aggregatorSpecs, err := aggregators.MakeSpecs(fields, aggregatorDescs)
	if err != nil {
		b.Fatalf("MakeSpecs: %s", err)
	}
	goals, err := goal.MakeGoals(goalExprs)
	if err != nil {
		b.Fatalf("MakeGoals: %s", err)
	}
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		_, err := AssessRules(dataset, rules, aggregatorSpecs, goals)
		if err != nil {
			b.Errorf("AssessRules: %s", err)
		}
	}
}

/******************************
 *  Helper functions
 ******************************/

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
