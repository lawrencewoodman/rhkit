package rulehunter

import (
	"errors"
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/dlit_go"
	"github.com/lawrencewoodman/rulehunter/internal"
	"reflect"
	"testing"
)

func TestAssessRules(t *testing.T) {
	rules := []*Rule{
		mustNewRule("band > 4"),
		mustNewRule("band > 3"),
		mustNewRule("cost > 1.2"),
	}
	inAggregators := []internal.Aggregator{
		mustNewCountAggregator("numIncomeGt2", "income > 2"),
		mustNewCountAggregator("numBandGt4", "band > 4"),
	}
	goals := []*dexpr.Expr{
		mustNewDExpr("numIncomeGt2 == 1"),
		mustNewDExpr("numIncomeGt2 == 2"),
		mustNewDExpr("numIncomeGt2 == 3"),
		mustNewDExpr("numIncomeGt2 == 4"),
		mustNewDExpr("numBandGt4 == 1"),
		mustNewDExpr("numBandGt4 == 2"),
		mustNewDExpr("numBandGt4 == 3"),
		mustNewDExpr("numBandGt4 == 4"),
	}
	records := []map[string]*dlit.Literal{
		map[string]*dlit.Literal{
			"income": dlit.MustNew(3),
			"cost":   dlit.MustNew(4.5),
			"band":   dlit.MustNew(4),
		},
		map[string]*dlit.Literal{
			"income": dlit.MustNew(3),
			"cost":   dlit.MustNew(3.2),
			"band":   dlit.MustNew(7),
		},
		map[string]*dlit.Literal{
			"income": dlit.MustNew(2),
			"cost":   dlit.MustNew(1.2),
			"band":   dlit.MustNew(4),
		},
		map[string]*dlit.Literal{
			"income": dlit.MustNew(0),
			"cost":   dlit.MustNew(0),
			"band":   dlit.MustNew(9),
		},
	}
	wantAssessment := Assessment{
		NumRecords: int64(len(records)),
		Flags: map[string]bool{
			"sorted": false,
		},
		RuleAssessments: []*RuleFinalAssessment{
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 4"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"numBandGt4":     dlit.MustNew("2"),
				},
				Goals: map[string]bool{
					"numIncomeGt2 == 1": true,
					"numIncomeGt2 == 2": false,
					"numIncomeGt2 == 3": false,
					"numIncomeGt2 == 4": false,
					"numBandGt4 == 1":   false,
					"numBandGt4 == 2":   true,
					"numBandGt4 == 3":   false,
					"numBandGt4 == 4":   false,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 3"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"numBandGt4":     dlit.MustNew("2"),
				},
				Goals: map[string]bool{
					"numIncomeGt2 == 1": false,
					"numIncomeGt2 == 2": true,
					"numIncomeGt2 == 3": false,
					"numIncomeGt2 == 4": false,
					"numBandGt4 == 1":   false,
					"numBandGt4 == 2":   true,
					"numBandGt4 == 3":   false,
					"numBandGt4 == 4":   false,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("cost > 1.2"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"numBandGt4":     dlit.MustNew("1"),
				},
				Goals: map[string]bool{
					"numIncomeGt2 == 1": false,
					"numIncomeGt2 == 2": true,
					"numIncomeGt2 == 3": false,
					"numIncomeGt2 == 4": false,
					"numBandGt4 == 1":   true,
					"numBandGt4 == 2":   false,
					"numBandGt4 == 3":   false,
					"numBandGt4 == 4":   false,
				},
			},
		},
	}
	input := NewLiteralInput(records)
	gotAssessment, err := AssessRules(rules, inAggregators, goals, input)
	if err != nil {
		t.Errorf("AssessRules(%q, %q, %q, input) - err: %q",
			rules, inAggregators, goals, err)
	}
	if !gotAssessment.IsEqual(&wantAssessment) {
		t.Errorf("AssessRules(%q, %q, %q, input)\ngot: %q\nwant: %q\n",
			rules, inAggregators, goals, gotAssessment, wantAssessment)
	}
}

func TestAssessRules_errors(t *testing.T) {
	cases := []struct {
		rules       []*Rule
		aggregators []internal.Aggregator
		goals       []*dexpr.Expr
		wantErr     error
	}{
		{[]*Rule{mustNewRule("band ^^ 3")},
			[]internal.Aggregator{
				mustNewCountAggregator("numIncomeGt2", "income > 2")},
			[]*dexpr.Expr{mustNewDExpr("numIncomeGt2 == 1")},
			errors.New("Invalid operator: \"^\"")},
		{[]*Rule{mustNewRule("hand > 3")},
			[]internal.Aggregator{
				mustNewCountAggregator("numIncomeGt2", "income > 2")},
			[]*dexpr.Expr{mustNewDExpr("numIncomeGt2 == 1")},
			errors.New("Variable doesn't exist: hand")},
		{[]*Rule{mustNewRule("band > 3")},
			[]internal.Aggregator{
				mustNewCountAggregator("numIncomeGt2", "bincome > 2")},
			[]*dexpr.Expr{mustNewDExpr("numIncomeGt2 == 1")},
			errors.New("Variable doesn't exist: bincome")},
		{[]*Rule{mustNewRule("band > 3")},
			[]internal.Aggregator{
				mustNewCountAggregator("numIncomeGt2", "income > 2")},
			[]*dexpr.Expr{mustNewDExpr("numIncomeGt == 1")},
			errors.New("Variable doesn't exist: numIncomeGt")},
	}
	records := []map[string]*dlit.Literal{
		map[string]*dlit.Literal{
			"income": dlit.MustNew(3),
			"cost":   dlit.MustNew(4.5),
			"band":   dlit.MustNew(4),
		},
	}
	input := NewLiteralInput(records)
	for _, c := range cases {
		_, err := AssessRules(c.rules, c.aggregators, c.goals, input)
		if err.Error() != c.wantErr.Error() {
			t.Errorf("AssessRules(%q, %q, %q, input) - err: %s, wantErr: %s",
				c.rules, c.aggregators, c.goals, err, c.wantErr)
		}
	}
}

func TestSort(t *testing.T) {
	assessment := Assessment{
		NumRecords: 8,
		Flags: map[string]bool{
			"sorted": false,
		},
		RuleAssessments: []*RuleFinalAssessment{
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 9"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
					"numIncomeGt2":   dlit.MustNew("3"),
					"numBandGt4":     dlit.MustNew("2"),
				},
				Goals: map[string]bool{
					"numIncomeGt2 == 1": false,
					"numIncomeGt2 == 2": true,
					"numIncomeGt2 == 3": false,
					"numIncomeGt2 == 4": false,
					"numBandGt4 == 1":   false,
					"numBandGt4 == 2":   true,
					"numBandGt4 == 3":   false,
					"numBandGt4 == 4":   true,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 456"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"numBandGt4":     dlit.MustNew("2"),
				},
				Goals: map[string]bool{
					"numIncomeGt2 == 1": true,
					"numIncomeGt2 == 2": false,
					"numIncomeGt2 == 3": false,
					"numIncomeGt2 == 4": false,
					"numBandGt4 == 1":   false,
					"numBandGt4 == 2":   true,
					"numBandGt4 == 3":   false,
					"numBandGt4 == 4":   false,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 3"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("76.3"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"numBandGt4":     dlit.MustNew("2"),
				},
				Goals: map[string]bool{
					"numIncomeGt2 == 1": false,
					"numIncomeGt2 == 2": true,
					"numIncomeGt2 == 3": false,
					"numIncomeGt2 == 4": false,
					"numBandGt4 == 1":   false,
					"numBandGt4 == 2":   true,
					"numBandGt4 == 3":   false,
					"numBandGt4 == 4":   false,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("cost > 1.2"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"numBandGt4":     dlit.MustNew("1"),
				},
				Goals: map[string]bool{
					"numIncomeGt2 == 1": false,
					"numIncomeGt2 == 2": true,
					"numIncomeGt2 == 3": false,
					"numIncomeGt2 == 4": false,
					"numBandGt4 == 1":   true,
					"numBandGt4 == 2":   false,
					"numBandGt4 == 3":   false,
					"numBandGt4 == 4":   false,
				},
			},
		},
	}
	cases := []struct {
		sortOrder []SortField
		wantRules []string
	}{
		{[]SortField{SortField{"numGoalsPassed", ASCENDING}},
			[]string{"band > 456", "band > 9", "band > 3", "cost > 1.2"}},
		{[]SortField{SortField{"percentMatches", DESCENDING}},
			[]string{"band > 3", "band > 9", "cost > 1.2", "band > 456"}},
		{[]SortField{SortField{"percentMatches", ASCENDING}},
			[]string{"cost > 1.2", "band > 456", "band > 9", "band > 3"}},
		{[]SortField{SortField{"percentMatches", ASCENDING},
			SortField{"numIncomeGt2", ASCENDING}},
			[]string{"band > 456", "cost > 1.2", "band > 9", "band > 3"}},
		{[]SortField{SortField{"percentMatches", DESCENDING},
			SortField{"numIncomeGt2", ASCENDING}},
			[]string{"band > 9", "band > 3", "cost > 1.2", "band > 456"}},
	}
	for _, c := range cases {
		assessment.Sort(c.sortOrder)
		if !assessment.Flags["sorted"] {
			t.Errorf("Sort(%s) 'sorted' flag not set", c.sortOrder)
		}
		gotRules := getAssessmentRules(&assessment)
		rulesMatch, msg := matchRules(gotRules, c.wantRules)
		if !rulesMatch {
			t.Errorf("matchRules() rules don't match: %s\ngot: %s\nwant: %s\n",
				msg, gotRules, c.wantRules)
		}
	}
}

func TestGetRules(t *testing.T) {
	assessment := Assessment{NumRecords: 8,
		RuleAssessments: []*RuleFinalAssessment{
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 9"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": true,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 456"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": false,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 3"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("76.3"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": true,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("cost > 1.2"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": false,
				},
			},
		},
	}
	wantRules := []*Rule{
		mustNewRule("band > 9"),
		mustNewRule("band > 456"),
		mustNewRule("band > 3"),
		mustNewRule("cost > 1.2"),
	}
	gotRules := assessment.GetRules()
	if !reflect.DeepEqual(gotRules, wantRules) {
		t.Errorf("GetRuleString() rules don't match\ngot: %s\nwant: %s\n",
			gotRules, wantRules)
	}
}

func TestMerge(t *testing.T) {
	assessment1 := &Assessment{
		NumRecords: 8,
		Flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleFinalAssessment{
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 9"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": true,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 456"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": false,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 3"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("76.3"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": true,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("cost > 1.2"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": false,
				},
			},
		},
	}
	assessment2 := &Assessment{
		NumRecords: 8,
		Flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleFinalAssessment{
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 16"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("8"),
					"percentMatches": dlit.MustNew("5.3"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": true,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("team == \"Pi\""),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("19"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": false,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 36"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("6.3"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": false,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("cost > 1.27"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("3.5"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": false,
				},
			},
		},
	}

	wantAssessment := &Assessment{
		NumRecords: 8,
		Flags: map[string]bool{
			"sorted": false,
		},
		RuleAssessments: []*RuleFinalAssessment{
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 9"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": true,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 456"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": false,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 3"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("76.3"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": true,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("cost > 1.2"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": false,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 16"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("8"),
					"percentMatches": dlit.MustNew("5.3"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": true,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("team == \"Pi\""),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("19"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": false,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 36"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("6.3"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": false,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("cost > 1.27"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("3.5"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": false,
				},
			},
		},
	}
	gotAssessment, err := assessment1.Merge(assessment2)
	if err != nil {
		t.Errorf("Merge() error: %s", err)
		return
	}
	if !reflect.DeepEqual(gotAssessment, wantAssessment) {
		t.Errorf("Merge() got assessment: %s\nwant: %s\n",
			gotAssessment, wantAssessment)
	}
}

func TestMerge_errors(t *testing.T) {
	assessment1 := &Assessment{
		NumRecords: 8,
		Flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleFinalAssessment{
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 9"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": true,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 456"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": false,
				},
			},
		},
	}
	assessment2 := &Assessment{NumRecords: 2,
		RuleAssessments: []*RuleFinalAssessment{
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 16"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("8"),
					"percentMatches": dlit.MustNew("5.3"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": true,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("team == \"Pi\""),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("19"),
				},
				Goals: map[string]bool{
					"numMatches > 3 ": false,
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
		Flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleFinalAssessment{
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 4"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("5"),
				},
				Goals: map[string]bool{
					"numIncomeGt2 == 1": true,
					"numIncomeGt2 == 2": false,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("in(band,\"4\",\"3\",\"2\")"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
					"numIncomeGt2":   dlit.MustNew("4"),
				},
				Goals: map[string]bool{
					"numIncomeGt2 == 1": false,
					"numIncomeGt2 == 2": true,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("in(team,\"a\",\"b\")"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
					"numIncomeGt2":   dlit.MustNew("4"),
				},
				Goals: map[string]bool{
					"numIncomeGt2 == 1": false,
					"numIncomeGt2 == 2": true,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("in(band,\"99\",\"23\")"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
					"numIncomeGt2":   dlit.MustNew("4"),
				},
				Goals: map[string]bool{
					"numIncomeGt2 == 1": false,
					"numIncomeGt2 == 2": true,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 3"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
					"numIncomeGt2":   dlit.MustNew("4"),
				},
				Goals: map[string]bool{
					"numIncomeGt2 == 1": false,
					"numIncomeGt2 == 2": true,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 9"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
					"numIncomeGt2":   dlit.MustNew("4"),
				},
				Goals: map[string]bool{
					"numIncomeGt2 == 1": false,
					"numIncomeGt2 == 2": true,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("in(band,\"9\",\"2\")"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
					"numIncomeGt2":   dlit.MustNew("4"),
				},
				Goals: map[string]bool{
					"numIncomeGt2 == 1": false,
					"numIncomeGt2 == 2": true,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("band == 7"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("0"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("3"),
				},
				Goals: map[string]bool{
					"numIncomeGt2 == 1": false,
					"numIncomeGt2 == 2": true,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("true()"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("3"),
				},
				Goals: map[string]bool{
					"numIncomeGt2 == 1": false,
					"numIncomeGt2 == 2": true,
				},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("cost > 1.2"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("3"),
				},
				Goals: map[string]bool{
					"numIncomeGt2 == 1": false,
					"numIncomeGt2 == 2": true,
				},
			},
		},
	}
	wantRules := []string{
		"band > 4",
		"in(band,\"4\",\"3\",\"2\")",
		"in(team,\"a\",\"b\")",
		"in(band,\"99\",\"23\")",
		"band > 3",
		"true()",
	}
	numSimilarRules := 2
	sortedAssessment.Refine(numSimilarRules)
	gotRules := getAssessmentRules(sortedAssessment)
	rulesMatch, msg := matchRules(gotRules, wantRules)
	if !rulesMatch {
		t.Errorf("matchRules() rules don't match: %s\ngot: %s\nwant: %s\n",
			msg, gotRules, wantRules)
	}
}

func TestRefine_panic_1(t *testing.T) {
	testPurpose := "Ensure panics if assessment not sorted"
	unsortedAssessment := &Assessment{
		NumRecords: 20,
		RuleAssessments: []*RuleFinalAssessment{
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 4"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: map[string]bool{},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("true()"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
				},
				Goals: map[string]bool{},
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
	numSimilarRules := 1
	unsortedAssessment.Refine(numSimilarRules)
	if !paniced {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("Refine() - failed to panic with: %s", wantPanic)
	}
}

func TestRefine_panic_2(t *testing.T) {
	testPurpose := "Ensure panics if 'true()' rule missing"
	sortedAssessment := &Assessment{
		NumRecords: 20,
		Flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleFinalAssessment{
			&RuleFinalAssessment{
				Rule: mustNewRule("band > 4"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: map[string]bool{},
			},
			&RuleFinalAssessment{
				Rule: mustNewRule("team > 7"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
				},
				Goals: map[string]bool{},
			},
		},
	}
	paniced := false
	wantPanic := "No 'true()' rule found"
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
	numSimilarRules := 1
	sortedAssessment.Refine(numSimilarRules)
	if !paniced {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("Refine() - failed to panic with: %s", wantPanic)
	}
}

/******************************
 *  Helper functions
 ******************************/

func getAssessmentRules(assessment *Assessment) []string {
	rules := make([]string, len(assessment.RuleAssessments))
	for i, ruleAssessment := range assessment.RuleAssessments {
		rules[i] = ruleAssessment.Rule.String()
	}
	return rules
}
