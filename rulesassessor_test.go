package main

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"io"
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
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     mustNewLit(2),
					"percentMatches": mustNewLit(50),
					"numIncomeGt2":   mustNewLit(1),
					"numBandGt4":     mustNewLit(2),
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
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     mustNewLit(4),
					"percentMatches": mustNewLit(100),
					"numIncomeGt2":   mustNewLit(2),
					"numBandGt4":     mustNewLit(2),
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
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     mustNewLit(2),
					"percentMatches": mustNewLit(50),
					"numIncomeGt2":   mustNewLit(2),
					"numBandGt4":     mustNewLit(1),
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

/*****************************
 *    Helper functions
 *****************************/
type LiteralInput struct {
	records  []map[string]*dlit.Literal
	position int
}

func NewLiteralInput(records []map[string]*dlit.Literal) Input {
	return &LiteralInput{records: records, position: 0}
}

func (l *LiteralInput) Read() (map[string]*dlit.Literal, error) {
	if l.position < len(l.records) {
		record := l.records[l.position]
		l.position++
		return record, nil
	}
	return map[string]*dlit.Literal{}, io.EOF
}
