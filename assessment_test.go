package rhkit

import (
	"errors"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/experiment"
	"github.com/vlifesystems/rhkit/rule"
	"reflect"
	"testing"
)

func TestGetRules(t *testing.T) {
	var gotRules []rule.Rule
	assessment := Assessment{NumRecords: 8,
		RuleAssessments: []*RuleAssessment{
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(9)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", true},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(456)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(3)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("76.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", true},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("cost", dlit.MustNew(1.2)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
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
			gotRules = assessment.GetRules(c.numRules)
		} else {
			gotRules = assessment.GetRules()
		}
		if !reflect.DeepEqual(gotRules, c.wantRules) {
			t.Errorf("GetRules() passNumRules: %t, numRules: %d rules don't match\ngot: %s\nwant: %s\n",
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
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(9)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", true},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(456)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(3)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("76.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", true},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("cost", dlit.MustNew(1.2)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
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
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(16)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("8"),
					"percentMatches": dlit.MustNew("5.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", true},
				},
			},
			&RuleAssessment{
				Rule: rule.NewEQFV("team", dlit.MustNew("Pi")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("19"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(36)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("6.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("cost", dlit.MustNew(1.27)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("3.5"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
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
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(9)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", true},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(456)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(3)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("76.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", true},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("cost", dlit.MustNew(1.2)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(16)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("8"),
					"percentMatches": dlit.MustNew("5.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", true},
				},
			},
			&RuleAssessment{
				Rule: rule.NewEQFV("team", dlit.MustNew("Pi")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("19"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(36)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("6.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("cost", dlit.MustNew(1.27)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("3.5"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
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
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(9)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", true},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(456)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
		},
	}
	assessment2 := &Assessment{NumRecords: 2,
		RuleAssessments: []*RuleAssessment{
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(16)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("8"),
					"percentMatches": dlit.MustNew("5.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", true},
				},
			},
			&RuleAssessment{
				Rule: rule.NewEQFV("team", dlit.MustNew("Pi")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("19"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
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
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("150"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(5)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("150"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
				Rule: rule.NewInFV("band", makeStringsDlitSlice("4", "3", "2")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("149"),
					"percentMatches": dlit.MustNew("49"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
				},
			},
			&RuleAssessment{
				Rule: rule.NewInFV("team", makeStringsDlitSlice("a", "b")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("148"),
					"percentMatches": dlit.MustNew("48"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
				},
			},
			&RuleAssessment{
				Rule: rule.NewInFV("band", makeStringsDlitSlice("99", "23")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("147"),
					"percentMatches": dlit.MustNew("47"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(3)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("146"),
					"percentMatches": dlit.MustNew("46"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(9)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("145"),
					"percentMatches": dlit.MustNew("45"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
				},
			},
			&RuleAssessment{
				Rule: rule.NewInFV("band", makeStringsDlitSlice("9", "2")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("144"),
					"percentMatches": dlit.MustNew("44"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
				},
			},
			&RuleAssessment{
				Rule: rule.NewEQFV("band", dlit.MustNew(7)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("0"),
					"percentMatches": dlit.MustNew("43"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
				},
			},
			&RuleAssessment{
				Rule: rule.NewEQFV("band", dlit.MustNew(8)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("142"),
					"percentMatches": dlit.MustNew("42"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
				},
			},
			&RuleAssessment{
				Rule: rule.NewTrue(),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("142"),
					"percentMatches": dlit.MustNew("42"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("cost", dlit.MustNew(1.2)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("141"),
					"percentMatches": dlit.MustNew("41"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
				},
			},
		},
	}
	wantRules := []rule.Rule{
		rule.NewGEFV("band", dlit.MustNew(4)),
		rule.NewInFV("band", makeStringsDlitSlice("4", "3", "2")),
		rule.NewInFV("team", makeStringsDlitSlice("a", "b")),
		rule.NewInFV("band", makeStringsDlitSlice("99", "23")),
		rule.NewTrue(),
	}
	sortedAssessment.Refine()
	gotRules := sortedAssessment.GetRules()

	if !matchRules(gotRules, wantRules) {
		t.Errorf("matchRules() rules don't match:\ngot: %s\nwant: %s\n",
			gotRules, wantRules)
	}
}

func TestRefine_between(t *testing.T) {
	sortedAssessment := &Assessment{
		NumRecords: 20,
		flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			&RuleAssessment{
				Rule: rule.MustNewBetweenFV("band", dlit.MustNew(5), dlit.MustNew(7)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("150"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("149"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
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
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
				Rule: rule.MustNewBetweenFV("band", dlit.MustNew(5), dlit.MustNew(6)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("147"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
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
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(5)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("145"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
				Rule: rule.MustNewBetweenFV("band", dlit.MustNew(7), dlit.MustNew(16)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("144"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
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
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
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
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
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
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
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
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
				Rule: rule.NewLEFV("band", dlit.MustNew(5)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("139"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(6)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("138"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
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
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
				Rule: rule.NewTrue(),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("136"),
					"percentMatches": dlit.MustNew("42"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
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
	gotRules := sortedAssessment.GetRules()

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
			&RuleAssessment{
				Rule: rule.MustNewOutsideFV("band", dlit.MustNew(5), dlit.MustNew(7)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("150"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("149"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
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
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
				Rule: rule.MustNewOutsideFV("band", dlit.MustNew(5), dlit.MustNew(6)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("147"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
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
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(5)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("145"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
				Rule: rule.MustNewOutsideFV("band", dlit.MustNew(7), dlit.MustNew(16)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("144"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
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
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
				Rule: rule.NewTrue(),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("142"),
					"percentMatches": dlit.MustNew("42"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
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
	gotRules := sortedAssessment.GetRules()

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
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{},
			},
			&RuleAssessment{
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
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{},
			},
			&RuleAssessment{
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
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches": dlit.MustNew("2"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", true},
				},
			},
			&RuleAssessment{
				Rule: rule.NewInFV("band", makeStringsDlitSlice("4", "3", "2")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches": dlit.MustNew("4"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
				},
			},
			&RuleAssessment{
				Rule: rule.NewInFV("team", makeStringsDlitSlice("a", "b")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches": dlit.MustNew("4"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
				},
			},
			&RuleAssessment{
				Rule: rule.NewInFV("band", makeStringsDlitSlice("99", "23")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches": dlit.MustNew("4"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
				},
			},
			&RuleAssessment{
				Rule: rule.NewTrue(),
				Aggregators: map[string]*dlit.Literal{
					"numMatches": dlit.MustNew("2"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
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
				rule.NewInFV("band", makeStringsDlitSlice("4", "3", "2")),
				rule.NewTrue(),
			},
		},
		{4,
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(4)),
				rule.NewInFV("band", makeStringsDlitSlice("4", "3", "2")),
				rule.NewInFV("team", makeStringsDlitSlice("a", "b")),
				rule.NewTrue(),
			},
		},
		{5,
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(4)),
				rule.NewInFV("band", makeStringsDlitSlice("4", "3", "2")),
				rule.NewInFV("team", makeStringsDlitSlice("a", "b")),
				rule.NewInFV("band", makeStringsDlitSlice("99", "23")),
				rule.NewTrue(),
			},
		},
		{6,
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(4)),
				rule.NewInFV("band", makeStringsDlitSlice("4", "3", "2")),
				rule.NewInFV("team", makeStringsDlitSlice("a", "b")),
				rule.NewInFV("band", makeStringsDlitSlice("99", "23")),
				rule.NewTrue(),
			},
		},
	}
	for _, c := range cases {
		limitedAssessment := refinedAssessment.TruncateRuleAssessments(c.numRules)
		gotRules := limitedAssessment.GetRules()
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
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{},
			},
			&RuleAssessment{
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
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{},
			},
			&RuleAssessment{
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
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(9)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
					"goalsScore":     dlit.MustNew(0.003),
					"numIncomeGt2":   dlit.MustNew("3"),
					"numBandGt4":     dlit.MustNew("2"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
					&GoalAssessment{"numIncomeGt2 == 3", false},
					&GoalAssessment{"numIncomeGt2 == 4", false},
					&GoalAssessment{"numBandGt4 == 1", false},
					&GoalAssessment{"numBandGt4 == 2", true},
					&GoalAssessment{"numBandGt4 == 3", false},
					&GoalAssessment{"numBandGt4 == 4", true},
				},
			},
			&RuleAssessment{
				Rule: rule.NewGEFV("band", dlit.MustNew(456)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
					"goalsScore":     dlit.MustNew(1.001),
					"numIncomeGt2":   dlit.MustNew("2"),
					"numBandGt4":     dlit.MustNew("2"),
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
				Rule: rule.NewGEFV("band", dlit.MustNew(3)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("76.3"),
					"goalsScore":     dlit.MustNew(0.002),
					"numIncomeGt2":   dlit.MustNew("2"),
					"numBandGt4":     dlit.MustNew("2"),
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
				Rule: rule.NewGEFV("cost", dlit.MustNew(1.2)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
					"goalsScore":     dlit.MustNew(0.002),
					"numIncomeGt2":   dlit.MustNew("1"),
					"numBandGt4":     dlit.MustNew("1"),
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
		},
	}
	cases := []struct {
		sortOrder []experiment.SortField
		wantRules []rule.Rule
	}{
		{[]experiment.SortField{
			experiment.SortField{"goalsScore", experiment.ASCENDING},
		},
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(3)),
				rule.NewGEFV("cost", dlit.MustNew(1.2)),
				rule.NewGEFV("band", dlit.MustNew(9)),
				rule.NewGEFV("band", dlit.MustNew(456)),
			}},
		{[]experiment.SortField{
			experiment.SortField{"percentMatches", experiment.DESCENDING},
		},
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(3)),
				rule.NewGEFV("band", dlit.MustNew(9)),
				rule.NewGEFV("band", dlit.MustNew(456)),
				rule.NewGEFV("cost", dlit.MustNew(1.2)),
			}},
		{[]experiment.SortField{
			experiment.SortField{"percentMatches", experiment.ASCENDING},
		},
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(456)),
				rule.NewGEFV("cost", dlit.MustNew(1.2)),
				rule.NewGEFV("band", dlit.MustNew(9)),
				rule.NewGEFV("band", dlit.MustNew(3)),
			}},
		{[]experiment.SortField{
			experiment.SortField{"percentMatches", experiment.ASCENDING},
			experiment.SortField{"numIncomeGt2", experiment.ASCENDING},
		},
			[]rule.Rule{
				rule.NewGEFV("cost", dlit.MustNew(1.2)),
				rule.NewGEFV("band", dlit.MustNew(456)),
				rule.NewGEFV("band", dlit.MustNew(9)),
				rule.NewGEFV("band", dlit.MustNew(3)),
			}},
		{[]experiment.SortField{
			experiment.SortField{"percentMatches", experiment.DESCENDING},
			experiment.SortField{"numIncomeGt2", experiment.ASCENDING},
		},
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(3)),
				rule.NewGEFV("band", dlit.MustNew(9)),
				rule.NewGEFV("cost", dlit.MustNew(1.2)),
				rule.NewGEFV("band", dlit.MustNew(456)),
			}},
		{[]experiment.SortField{},
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
		gotRules := assessment.GetRules()
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
