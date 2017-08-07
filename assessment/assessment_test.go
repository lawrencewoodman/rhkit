package assessment

import (
	"errors"
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/aggregator"
	"github.com/vlifesystems/rhkit/goal"
	"github.com/vlifesystems/rhkit/internal/testhelpers"
	"github.com/vlifesystems/rhkit/rule"
	"reflect"
	"testing"
)

func TestAddRuleAssessors_error(t *testing.T) {
	numRecords := int64(100)
	as := aggregator.MustNew("a", "calc", "3+4")
	ai := as.New()
	ruleAssessors := []*ruleAssessor{
		{
			Rule:        rule.NewEQFV("month", dlit.NewString("May")),
			Aggregators: []aggregator.Instance{ai},
			Goals:       []*goal.Goal{goal.MustNew("cost > 3")},
		},
	}
	wantErr := dexpr.InvalidExprError{
		Expr: "cost > 3",
		Err:  dexpr.VarNotExistError("cost"),
	}
	assessment := newAssessment(numRecords)
	err := assessment.AddRuleAssessors(ruleAssessors)
	if err == nil || err.Error() != wantErr.Error() {
		t.Errorf("AddRuleAssessors: err: %s, wantErr: %s", err, wantErr)
	}
}

func TestRules(t *testing.T) {
	var gotRules []rule.Rule
	assessment := Assessment{NumRecords: 8,
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(9)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", true},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(456)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(3)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("76.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", true},
				},
			},
			{
				Rule: rule.NewGEFV("cost", dlit.MustNew(1.2)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
		},
	}
	cases := []struct {
		numRules     int
		passNumRules bool
		wantRules    []rule.Rule
	}{
		{0, true, []rule.Rule{}},
		{1, true, []rule.Rule{rule.NewGEFV("band", dlit.MustNew(9))}},
		{2, true, []rule.Rule{
			rule.NewGEFV("band", dlit.MustNew(9)),
			rule.NewGEFV("band", dlit.MustNew(456)),
		}},
		{4, true, []rule.Rule{
			rule.NewGEFV("band", dlit.MustNew(9)),
			rule.NewGEFV("band", dlit.MustNew(456)),
			rule.NewGEFV("band", dlit.MustNew(3)),
			rule.NewGEFV("cost", dlit.MustNew(1.2)),
		},
		},
		{5, true, []rule.Rule{
			rule.NewGEFV("band", dlit.MustNew(9)),
			rule.NewGEFV("band", dlit.MustNew(456)),
			rule.NewGEFV("band", dlit.MustNew(3)),
			rule.NewGEFV("cost", dlit.MustNew(1.2)),
		},
		},
		{0, false, []rule.Rule{
			rule.NewGEFV("band", dlit.MustNew(9)),
			rule.NewGEFV("band", dlit.MustNew(456)),
			rule.NewGEFV("band", dlit.MustNew(3)),
			rule.NewGEFV("cost", dlit.MustNew(1.2)),
		},
		},
	}
	for _, c := range cases {
		if c.passNumRules {
			gotRules = assessment.Rules(c.numRules)
		} else {
			gotRules = assessment.Rules()
		}
		if !reflect.DeepEqual(gotRules, c.wantRules) {
			t.Errorf("Rules() passNumRules: %t, numRules: %d rules don't match\ngot: %s\nwant: %s\n",
				c.passNumRules, c.numRules, gotRules, c.wantRules)
		}
	}
}

func TestMerge(t *testing.T) {
	assessment1 := &Assessment{
		NumRecords: 8,
		flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(9)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", true},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(456)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(3)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("76.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", true},
				},
			},
			{
				Rule: rule.NewGEFV("cost", dlit.MustNew(1.2)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
		},
	}
	assessment2 := &Assessment{
		NumRecords: 8,
		flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(16)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("8"),
					"percentMatches": dlit.MustNew("5.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", true},
				},
			},
			{
				Rule: rule.NewEQFV("team", dlit.MustNew("Pi")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("19"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(36)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("6.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
			{
				Rule: rule.NewGEFV("cost", dlit.MustNew(1.27)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("3.5"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
		},
	}

	wantAssessment := &Assessment{
		NumRecords: 8,
		flags: map[string]bool{
			"sorted":  false,
			"refined": false,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(9)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", true},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(456)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(3)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("76.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", true},
				},
			},
			{
				Rule: rule.NewGEFV("cost", dlit.MustNew(1.2)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(16)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("8"),
					"percentMatches": dlit.MustNew("5.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", true},
				},
			},
			{
				Rule: rule.NewEQFV("team", dlit.MustNew("Pi")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("19"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(36)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("6.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
			{
				Rule: rule.NewGEFV("cost", dlit.MustNew(1.27)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("3.5"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
		},
	}
	gotAssessment, err := assessment1.Merge(assessment2)
	if err != nil {
		t.Errorf("Merge() error: %s", err)
		return
	}

	if !gotAssessment.IsEqual(wantAssessment) {
		t.Errorf("Merge() assessments don't match\n - got: %v\n - want: %v\n",
			gotAssessment, wantAssessment)
	}
}

func TestMerge_errors(t *testing.T) {
	assessment1 := &Assessment{
		NumRecords: 8,
		flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(9)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", true},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(456)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
		},
	}
	assessment2 := &Assessment{NumRecords: 2,
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(16)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("8"),
					"percentMatches": dlit.MustNew("5.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", true},
				},
			},
			{
				Rule: rule.NewEQFV("team", dlit.MustNew("Pi")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("19"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
		},
	}
	wantError :=
		errors.New("Can't merge assessments: Number of records don't match")
	_, err := assessment1.Merge(assessment2)
	if err == nil {
		t.Errorf("Merge() not error, expected: %s", wantError)
		return
	}
	if err.Error() != wantError.Error() {
		t.Errorf("Merge() got error: %s, want: %s", err, wantError)
	}
}

func TestRefine(t *testing.T) {
	sortedAssessment := &Assessment{
		NumRecords: 20,
		flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("150"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(5)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("150"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.NewInFV(
					"band",
					testhelpers.MakeStringsDlitSlice("4", "3", "2"),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("149"),
					"percentMatches": dlit.MustNew("49"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
			{
				Rule: rule.NewInFV("team", testhelpers.MakeStringsDlitSlice("a", "b")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("148"),
					"percentMatches": dlit.MustNew("48"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
			{
				Rule: rule.NewInFV(
					"band",
					testhelpers.MakeStringsDlitSlice("99", "23"),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("147"),
					"percentMatches": dlit.MustNew("47"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(3)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("146"),
					"percentMatches": dlit.MustNew("46"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(9)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("145"),
					"percentMatches": dlit.MustNew("45"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
			{
				Rule: rule.NewInFV("band", testhelpers.MakeStringsDlitSlice("9", "2")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("144"),
					"percentMatches": dlit.MustNew("44"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
			{
				Rule: rule.NewEQFV("band", dlit.MustNew(7)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("0"),
					"percentMatches": dlit.MustNew("43"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
			{
				Rule: rule.NewEQFV("band", dlit.MustNew(8)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("142"),
					"percentMatches": dlit.MustNew("42"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
			{
				Rule: rule.NewTrue(),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("142"),
					"percentMatches": dlit.MustNew("42"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
			{
				Rule: rule.NewGEFV("cost", dlit.MustNew(1.2)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("141"),
					"percentMatches": dlit.MustNew("41"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
		},
	}
	wantRules := []rule.Rule{
		rule.NewGEFV("band", dlit.MustNew(4)),
		rule.NewInFV("band", testhelpers.MakeStringsDlitSlice("4", "3", "2")),
		rule.NewInFV("team", testhelpers.MakeStringsDlitSlice("a", "b")),
		rule.NewInFV("band", testhelpers.MakeStringsDlitSlice("99", "23")),
		rule.NewTrue(),
	}
	sortedAssessment.Refine()
	gotRules := sortedAssessment.Rules()

	if !matchRules(gotRules, wantRules) {
		t.Errorf("matchRules() rules don't match:\ngot: %s\nwant: %s\n",
			gotRules, wantRules)
	}
}

func TestRefine_few_ruleassessments(t *testing.T) {
	cases := []struct {
		in        *Assessment
		wantRules []rule.Rule
	}{
		{in: &Assessment{
			NumRecords: 20,
			flags: map[string]bool{
				"sorted": true,
			},
			RuleAssessments: []*RuleAssessment{
				{
					Rule: rule.NewTrue(),
					Aggregators: map[string]*dlit.Literal{
						"numMatches":     dlit.MustNew("142"),
						"percentMatches": dlit.MustNew("42"),
						"numIncomeGt2":   dlit.MustNew("2"),
						"goalsScore":     dlit.MustNew(0.1),
					},
					Goals: []*GoalAssessment{
						{"numIncomeGt2 == 1", false},
						{"numIncomeGt2 == 2", true},
					},
				},
			},
		},
			wantRules: []rule.Rule{rule.NewTrue()},
		},
		{in: &Assessment{
			NumRecords: 20,
			flags: map[string]bool{
				"sorted": true,
			},
			RuleAssessments: []*RuleAssessment{
				{
					Rule: rule.NewEQFV("month", dlit.NewString("april")),
					Aggregators: map[string]*dlit.Literal{
						"numMatches":     dlit.MustNew("142"),
						"percentMatches": dlit.MustNew("42"),
						"goalsScore":     dlit.MustNew(0.1),
					},
					Goals: []*GoalAssessment{
						{"numIncomeGt2 == 1", false},
						{"numIncomeGt2 == 2", true},
					},
				},
				{
					Rule: rule.NewTrue(),
					Aggregators: map[string]*dlit.Literal{
						"numMatches":     dlit.MustNew("142"),
						"percentMatches": dlit.MustNew("42"),
						"numIncomeGt2":   dlit.MustNew("2"),
						"goalsScore":     dlit.MustNew(0.1),
					},
					Goals: []*GoalAssessment{
						{"numIncomeGt2 == 1", false},
						{"numIncomeGt2 == 2", true},
					},
				},
			},
		},
			wantRules: []rule.Rule{
				rule.NewEQFV("month", dlit.NewString("april")),
				rule.NewTrue(),
			},
		},
	}
	for _, c := range cases {
		c.in.Refine()
		gotRules := c.in.Rules()

		if !matchRules(gotRules, c.wantRules) {
			t.Errorf("matchRules() rules don't match:\ngot: %s\nwant: %s\n",
				gotRules, c.wantRules)
		}
	}
}

func TestRefine_between(t *testing.T) {
	sortedAssessment := &Assessment{
		NumRecords: 20,
		flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.MustNewBetweenFV("band", dlit.MustNew(5), dlit.MustNew(7)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("150"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("149"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewBetweenFV(
					"rate",
					dlit.MustNew(16.2),
					dlit.MustNew(17.93),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("148"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewBetweenFV("band", dlit.MustNew(5), dlit.MustNew(6)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("147"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewBetweenFV(
					"rate",
					dlit.MustNew(16.2),
					dlit.MustNew(17.89),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("146"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(5)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("145"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewBetweenFV("band", dlit.MustNew(7), dlit.MustNew(16)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("144"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewBetweenFV(
					"rate",
					dlit.MustNew(15.2),
					dlit.MustNew(17.89),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("143"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewBetweenFV(
					"rate",
					dlit.MustNew(10.2),
					dlit.MustNew(16.3),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("142"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewBetweenFV(
					"rate",
					dlit.MustNew(50.1),
					dlit.MustNew(60.3),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("141"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewBetweenFV(
					"rate",
					dlit.MustNew(16.1),
					dlit.MustNew(20.3),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("140"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.NewLEFV("band", dlit.MustNew(5)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("139"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(6)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("138"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewBetweenFV(
					"rate",
					dlit.MustNew(1.2),
					dlit.MustNew(6.3),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("137"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.NewTrue(),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("136"),
					"percentMatches": dlit.MustNew("42"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
		},
	}
	wantRules := []rule.Rule{
		rule.MustNewBetweenFV("band", dlit.MustNew(5), dlit.MustNew(7)),
		rule.NewGEFV("band", dlit.MustNew(4)),
		rule.MustNewBetweenFV(
			"rate",
			dlit.MustNew(16.2),
			dlit.MustNew(17.93),
		),
		rule.MustNewBetweenFV(
			"rate",
			dlit.MustNew(50.1),
			dlit.MustNew(60.3),
		),
		rule.NewLEFV("band", dlit.MustNew(5)),
		rule.MustNewBetweenFV(
			"rate",
			dlit.MustNew(1.2),
			dlit.MustNew(6.3),
		),
		rule.NewTrue(),
	}
	sortedAssessment.Refine()
	gotRules := sortedAssessment.Rules()

	if !matchRules(gotRules, wantRules) {
		t.Errorf("matchRules() rules don't match:\ngot: %s\nwant: %s\n",
			gotRules, wantRules)
	}
}

func TestRefine_outside(t *testing.T) {
	sortedAssessment := &Assessment{
		NumRecords: 20,
		flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.MustNewOutsideFV("band", dlit.MustNew(5), dlit.MustNew(7)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("150"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("149"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewOutsideFV(
					"rate",
					dlit.MustNew(16.2),
					dlit.MustNew(17.93),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("148"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewOutsideFV("band", dlit.MustNew(5), dlit.MustNew(6)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("147"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewOutsideFV(
					"rate",
					dlit.MustNew(16.2),
					dlit.MustNew(17.89),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("146"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(5)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("145"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewOutsideFV("band", dlit.MustNew(7), dlit.MustNew(16)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("144"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewOutsideFV(
					"rate",
					dlit.MustNew(15.2),
					dlit.MustNew(17.89),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("143"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.NewTrue(),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("142"),
					"percentMatches": dlit.MustNew("42"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
		},
	}
	wantRules := []rule.Rule{
		rule.MustNewOutsideFV("band", dlit.MustNew(5), dlit.MustNew(7)),
		rule.NewGEFV("band", dlit.MustNew(4)),
		rule.MustNewOutsideFV("rate", dlit.MustNew(16.2), dlit.MustNew(17.93)),
		rule.NewTrue(),
	}
	sortedAssessment.Refine()
	gotRules := sortedAssessment.Rules()

	if !matchRules(gotRules, wantRules) {
		t.Errorf("matchRules() rules don't match:\ngot: %s\nwant: %s\n",
			gotRules, wantRules)
	}
}

func TestRefine_panic_1(t *testing.T) {
	testPurpose := "Ensure panics if assessment not sorted"
	unsortedAssessment := &Assessment{
		NumRecords: 20,
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{},
			},
			{
				Rule: rule.NewTrue(),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
				},
				Goals: []*GoalAssessment{},
			},
		},
	}
	paniced := false
	wantPanic := "Assessment isn't sorted"
	defer func() {
		if r := recover(); r != nil {
			if r.(string) == wantPanic {
				paniced = true
			} else {
				t.Errorf("Test: %s\n", testPurpose)
				t.Errorf("Refine() - got panic: %s, wanted: %s",
					r, wantPanic)
			}
		}
	}()
	unsortedAssessment.Refine()
	if !paniced {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("Refine() - failed to panic with: %s", wantPanic)
	}
}

func TestRefine_panic_2(t *testing.T) {
	testPurpose := "Ensure panics if True rule missing"
	sortedAssessment := &Assessment{
		NumRecords: 20,
		flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{},
			},
			{
				Rule: rule.NewGEFV("team", dlit.MustNew(7)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
				},
				Goals: []*GoalAssessment{},
			},
		},
	}
	paniced := false
	wantPanic := "No True rule found"
	defer func() {
		if r := recover(); r != nil {
			if r.(string) == wantPanic {
				paniced = true
			} else {
				t.Errorf("Test: %s\n", testPurpose)
				t.Errorf("Refine() - got panic: %s, wanted: %s",
					r, wantPanic)
			}
		}
	}()
	sortedAssessment.Refine()
	if !paniced {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("Refine() - failed to panic with: %s", wantPanic)
	}
}

func TestTruncateRuleAssessments(t *testing.T) {
	refinedAssessment := &Assessment{
		NumRecords: 20,
		flags: map[string]bool{
			"sorted":  true,
			"refined": true,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches": dlit.MustNew("2"),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
				},
			},
			{
				Rule: rule.NewInFV(
					"band",
					testhelpers.MakeStringsDlitSlice("4", "3", "2"),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches": dlit.MustNew("4"),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
				},
			},
			{
				Rule: rule.NewInFV("team", testhelpers.MakeStringsDlitSlice("a", "b")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches": dlit.MustNew("4"),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
				},
			},
			{
				Rule: rule.NewInFV(
					"band",
					testhelpers.MakeStringsDlitSlice("99", "23"),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches": dlit.MustNew("4"),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
				},
			},
			{
				Rule: rule.NewTrue(),
				Aggregators: map[string]*dlit.Literal{
					"numMatches": dlit.MustNew("2"),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
				},
			},
		},
	}
	// The increasing order of the numRules is important as this also checks that
	// the ruleAssessments are cloned properly
	cases := []struct {
		numRules  int
		wantRules []rule.Rule
	}{
		{0,
			[]rule.Rule{},
		},
		{1,
			[]rule.Rule{
				rule.NewTrue(),
			},
		},
		{2,
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(4)),
				rule.NewTrue(),
			},
		},
		{3,
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(4)),
				rule.NewInFV("band", testhelpers.MakeStringsDlitSlice("4", "3", "2")),
				rule.NewTrue(),
			},
		},
		{4,
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(4)),
				rule.NewInFV("band", testhelpers.MakeStringsDlitSlice("4", "3", "2")),
				rule.NewInFV("team", testhelpers.MakeStringsDlitSlice("a", "b")),
				rule.NewTrue(),
			},
		},
		{5,
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(4)),
				rule.NewInFV("band", testhelpers.MakeStringsDlitSlice("4", "3", "2")),
				rule.NewInFV("team", testhelpers.MakeStringsDlitSlice("a", "b")),
				rule.NewInFV("band", testhelpers.MakeStringsDlitSlice("99", "23")),
				rule.NewTrue(),
			},
		},
		{6,
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(4)),
				rule.NewInFV("band", testhelpers.MakeStringsDlitSlice("4", "3", "2")),
				rule.NewInFV("team", testhelpers.MakeStringsDlitSlice("a", "b")),
				rule.NewInFV("band", testhelpers.MakeStringsDlitSlice("99", "23")),
				rule.NewTrue(),
			},
		},
	}
	for _, c := range cases {
		limitedAssessment := refinedAssessment.TruncateRuleAssessments(c.numRules)
		gotRules := limitedAssessment.Rules()
		if !matchRules(gotRules, c.wantRules) {
			t.Errorf("matchRules() rules don't match:\nnumRules: %d\ngot: %s\nwant: %s\n",
				c.numRules, gotRules, c.wantRules)
		}
	}
}

func TestTruncateRuleAssessment_panic_1(t *testing.T) {
	testPurpose := "Ensure panics if assessment not sorted"
	unsortedAssessment := &Assessment{
		NumRecords: 20,
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{},
			},
			{
				Rule: rule.NewTrue(),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
				},
				Goals: []*GoalAssessment{},
			},
		},
	}
	paniced := false
	wantPanic := "Assessment isn't sorted"
	defer func() {
		if r := recover(); r != nil {
			if r.(string) == wantPanic {
				paniced = true
			} else {
				t.Errorf("Test: %s\n", testPurpose)
				t.Errorf("TruncateRuleAssessments() - got panic: %s, wanted: %s",
					r, wantPanic)
			}
		}
	}()
	numRules := 1
	unsortedAssessment.TruncateRuleAssessments(numRules)
	if !paniced {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("Refine() - failed to panic with: %s", wantPanic)
	}
}

func TestTruncateRuleAssessment_panic_2(t *testing.T) {
	testPurpose := "Ensure panics if assessment not refined"
	unsortedAssessment := &Assessment{
		NumRecords: 20,
		flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{},
			},
			{
				Rule: rule.NewTrue(),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
				},
				Goals: []*GoalAssessment{},
			},
		},
	}
	paniced := false
	wantPanic := "Assessment isn't refined"
	defer func() {
		if r := recover(); r != nil {
			if r.(string) == wantPanic {
				paniced = true
			} else {
				t.Errorf("Test: %s\n", testPurpose)
				t.Errorf("TruncateRuleAssessments() - got panic: %s, wanted: %s",
					r, wantPanic)
			}
		}
	}()
	numRules := 1
	unsortedAssessment.TruncateRuleAssessments(numRules)
	if !paniced {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("Refine() - failed to panic with: %s", wantPanic)
	}
}

func TestSort(t *testing.T) {
	assessment := Assessment{
		NumRecords: 8,
		flags: map[string]bool{
			"sorted": false,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(9)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
					"goalsScore":     dlit.MustNew(0.003),
					"numIncomeGt2":   dlit.MustNew("3"),
					"numBandGt4":     dlit.MustNew("2"),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
					{"numIncomeGt2 == 3", false},
					{"numIncomeGt2 == 4", false},
					{"numBandGt4 == 1", false},
					{"numBandGt4 == 2", true},
					{"numBandGt4 == 3", false},
					{"numBandGt4 == 4", true},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(456)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
					"goalsScore":     dlit.MustNew(1.001),
					"numIncomeGt2":   dlit.MustNew("2"),
					"numBandGt4":     dlit.MustNew("2"),
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
				Rule: rule.NewGEFV("band", dlit.MustNew(3)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("76.3"),
					"goalsScore":     dlit.MustNew(0.002),
					"numIncomeGt2":   dlit.MustNew("2"),
					"numBandGt4":     dlit.MustNew("2"),
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
				Rule: rule.NewGEFV("cost", dlit.MustNew(1.2)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
					"goalsScore":     dlit.MustNew(0.002),
					"numIncomeGt2":   dlit.MustNew("1"),
					"numBandGt4":     dlit.MustNew("1"),
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
		},
	}
	cases := []struct {
		sortOrder []SortOrder
		wantRules []rule.Rule
	}{
		{[]SortOrder{
			{"goalsScore", ASCENDING},
		},
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(3)),
				rule.NewGEFV("cost", dlit.MustNew(1.2)),
				rule.NewGEFV("band", dlit.MustNew(9)),
				rule.NewGEFV("band", dlit.MustNew(456)),
			}},
		{[]SortOrder{
			{"percentMatches", DESCENDING},
		},
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(3)),
				rule.NewGEFV("band", dlit.MustNew(9)),
				rule.NewGEFV("band", dlit.MustNew(456)),
				rule.NewGEFV("cost", dlit.MustNew(1.2)),
			}},
		{[]SortOrder{
			{"percentMatches", ASCENDING},
		},
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(456)),
				rule.NewGEFV("cost", dlit.MustNew(1.2)),
				rule.NewGEFV("band", dlit.MustNew(9)),
				rule.NewGEFV("band", dlit.MustNew(3)),
			}},
		{[]SortOrder{
			{"percentMatches", ASCENDING},
			{"numIncomeGt2", ASCENDING},
		},
			[]rule.Rule{
				rule.NewGEFV("cost", dlit.MustNew(1.2)),
				rule.NewGEFV("band", dlit.MustNew(456)),
				rule.NewGEFV("band", dlit.MustNew(9)),
				rule.NewGEFV("band", dlit.MustNew(3)),
			}},
		{[]SortOrder{
			{"percentMatches", DESCENDING},
			{"numIncomeGt2", ASCENDING},
		},
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(3)),
				rule.NewGEFV("band", dlit.MustNew(9)),
				rule.NewGEFV("cost", dlit.MustNew(1.2)),
				rule.NewGEFV("band", dlit.MustNew(456)),
			}},
		{[]SortOrder{},
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(3)),
				rule.NewGEFV("band", dlit.MustNew(9)),
				rule.NewGEFV("band", dlit.MustNew(456)),
				rule.NewGEFV("cost", dlit.MustNew(1.2)),
			}},
	}
	for _, c := range cases {
		assessment.Sort(c.sortOrder)
		if !assessment.IsSorted() {
			t.Errorf("Sort(%s) 'sorted' flag not set", c.sortOrder)
		}
		gotRules := assessment.Rules()
		if !matchRules(gotRules, c.wantRules) {
			t.Errorf("matchRules() rules don't match:\n - sortOrder: %s\n - got: %s\n - want: %s\n",
				c.sortOrder, gotRules, c.wantRules)
		}
	}
}

/******************************
 *  Helper functions
 ******************************/

// Match the rules including their order
func matchRules(rules1 []rule.Rule, rules2 []rule.Rule) bool {
	if len(rules1) != len(rules2) {
		return false
	}
	for i, rule1 := range rules1 {
		if rule1.String() != rules2[i].String() {
			return false
		}
	}
	return true
}
