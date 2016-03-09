package main

import (
	"errors"
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/dlit_go"
	"sort"
	"testing"
)

func TestAssessRules(t *testing.T) {
	rules := []*dexpr.Expr{
		mustNewDExpr("band > 4"),
		mustNewDExpr("band > 3"),
		mustNewDExpr("cost > 1.2"),
	}
	aggregators := []Aggregator{
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
			"income": mustNewLit(3),
			"cost":   mustNewLit(4.5),
			"band":   mustNewLit(4),
		},
		map[string]*dlit.Literal{
			"income": mustNewLit(3),
			"cost":   mustNewLit(3.2),
			"band":   mustNewLit(7),
		},
		map[string]*dlit.Literal{
			"income": mustNewLit(2),
			"cost":   mustNewLit(1.2),
			"band":   mustNewLit(4),
		},
		map[string]*dlit.Literal{
			"income": mustNewLit(0),
			"cost":   mustNewLit(0),
			"band":   mustNewLit(9),
		},
	}
	wantReport := Report{NumRecords: int64(len(records)),
		RuleReports: []*RuleReport{
			&RuleReport{
				Rule: "band > 4",
				Aggregators: map[string]string{
					"numMatches":     "2",
					"percentMatches": "50",
					"numIncomeGt2":   "1",
					"numBandGt4":     "2",
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
			&RuleReport{
				Rule: "band > 3",
				Aggregators: map[string]string{
					"numMatches":     "4",
					"percentMatches": "100",
					"numIncomeGt2":   "2",
					"numBandGt4":     "2",
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
			&RuleReport{
				Rule: "cost > 1.2",
				Aggregators: map[string]string{
					"numMatches":     "2",
					"percentMatches": "50",
					"numIncomeGt2":   "2",
					"numBandGt4":     "1",
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
	gotReport, err := AssessRules(rules, aggregators, goals, input)
	if err != nil {
		t.Errorf("AssessRules(%q, %q, %q, input) - err: %q",
			rules, aggregators, goals, err)
	}
	if !gotReport.IsEqual(&wantReport) {
		t.Errorf("AssessRules(%q, %q, %q, input)\ngot: %q\nwant: %q\n",
			rules, aggregators, goals, gotReport, wantReport)
	}
}

func TestAssessRules_errors(t *testing.T) {
	cases := []struct {
		rules       []*dexpr.Expr
		aggregators []Aggregator
		goals       []*dexpr.Expr
		wantErr     error
	}{
		{[]*dexpr.Expr{mustNewDExpr("band ^^ 3")},
			[]Aggregator{mustNewCountAggregator("numIncomeGt2", "income > 2")},
			[]*dexpr.Expr{mustNewDExpr("numIncomeGt2 == 1")},
			errors.New("Invalid operator: \"^\"")},
		{[]*dexpr.Expr{mustNewDExpr("hand > 3")},
			[]Aggregator{mustNewCountAggregator("numIncomeGt2", "income > 2")},
			[]*dexpr.Expr{mustNewDExpr("numIncomeGt2 == 1")},
			errors.New("Variable doesn't exist: hand")},
		{[]*dexpr.Expr{mustNewDExpr("band > 3")},
			[]Aggregator{mustNewCountAggregator("numIncomeGt2", "bincome > 2")},
			[]*dexpr.Expr{mustNewDExpr("numIncomeGt2 == 1")},
			errors.New("Variable doesn't exist: bincome")},
		{[]*dexpr.Expr{mustNewDExpr("band > 3")},
			[]Aggregator{mustNewCountAggregator("numIncomeGt2", "income > 2")},
			[]*dexpr.Expr{mustNewDExpr("numIncomeGt == 1")},
			errors.New("Variable doesn't exist: numIncomeGt")},
	}
	records := []map[string]*dlit.Literal{
		map[string]*dlit.Literal{
			"income": mustNewLit(3),
			"cost":   mustNewLit(4.5),
			"band":   mustNewLit(4),
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
