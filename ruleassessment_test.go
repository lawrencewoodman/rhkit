package main

import (
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/dlit_go"
	"testing"
)

func TestNextRecord(t *testing.T) {
	// It is important for this test to reuse the aggregators to ensure that
	// they are cloned properly.
	aggregators := []Aggregator{
		mustNewCountAggregator("numIncomeGt2", "income > 2"),
		mustNewCountAggregator("numBandGt4", "band > 4"),
	}
	records := [4]map[string]*dlit.Literal{
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
	numRecords := int64(len(records))
	cases := []struct {
		rule             *dexpr.Expr
		goals            []*dexpr.Expr
		wantNumIncomeGt2 int64
		wantNumBandGt4   int64
	}{
		{mustNewDExpr("band > 4"),
			[]*dexpr.Expr{
				mustNewDExpr("numIncomeGt2 == 1"),
				mustNewDExpr("numBandGt4 == 2"),
			}, 1, 2},
		{mustNewDExpr("band > 3"),
			[]*dexpr.Expr{
				mustNewDExpr("numIncomeGt2 == 2"),
				mustNewDExpr("numBandGt4 == 2"),
			}, 2, 2},
		{mustNewDExpr("cost > 1.2"),
			[]*dexpr.Expr{
				mustNewDExpr("numIncomeGt2 == 2"),
				mustNewDExpr("numBandGt4 == 1"),
			}, 2, 1},
	}
	for _, c := range cases {
		ra := NewRuleAssessment(c.rule, aggregators, c.goals)
		for _, record := range records {
			err := ra.NextRecord(record)
			if err != nil {
				t.Errorf("NextRecord(%q) rule: %s, aggregators: %q, goals: %q - err: %q",
					record, c.rule, aggregators, c.goals, err)
			}
		}
		gotNumIncomeGt2, gt2Exists :=
			ra.GetAggregatorValue("numIncomeGt2", numRecords)
		if !gt2Exists {
			t.Errorf("numIncomeGt2 aggregator doesn't exist")
		}
		gotNumIncomeGt2Int, gt2IsInt := gotNumIncomeGt2.Int()
		if !gt2IsInt {
			t.Errorf("numIncomeGt2 aggregator can't be int")
		}
		gotNumBandGt4, gt4Exists :=
			ra.GetAggregatorValue("numBandGt4", numRecords)
		if !gt4Exists {
			t.Errorf("numBandGt4 aggregator doesn't exist")
		}
		gotNumBandGt4Int, gt4IsInt := gotNumBandGt4.Int()
		if !gt4IsInt {
			t.Errorf("numBandGt4 aggregator can't be int")
		}
		if gotNumIncomeGt2Int != c.wantNumIncomeGt2 {
			t.Errorf("NextRecord() rule: %s, aggregators: %q, goals: %q - wantNumIncomeGt2: %d, got: %d",
				c.rule, aggregators, c.goals, c.wantNumIncomeGt2, gotNumIncomeGt2Int)
		}
		if gotNumBandGt4Int != c.wantNumBandGt4 {
			t.Errorf("NextRecord() rule: %s, aggregators: %q, goals: %q - wantNumBandGt4: %d, got: %d",
				c.rule, aggregators, c.goals, c.wantNumBandGt4, gotNumBandGt4Int)
		}
	}
}

func TestNextRecord_Errors(t *testing.T) {
	records := [4]map[string]*dlit.Literal{
		map[string]*dlit.Literal{"income": dlit.MustNew(3), "band": dlit.MustNew(4)},
		map[string]*dlit.Literal{"income": dlit.MustNew(3), "band": dlit.MustNew(7)},
		map[string]*dlit.Literal{"income": dlit.MustNew(2), "band": dlit.MustNew(4)},
		map[string]*dlit.Literal{"income": dlit.MustNew(0), "band": dlit.MustNew(9)},
	}
	cases := []struct {
		rule        *dexpr.Expr
		aggregators []Aggregator
		goals       []*dexpr.Expr
		wantErr     error
	}{
		{mustNewDExpr("band > 4"),
			[]Aggregator{mustNewCountAggregator("numIncomeGt2", "income > 2")},
			[]*dexpr.Expr{mustNewDExpr("numIncomeGt2 == 1")}, nil},
		{mustNewDExpr("band > 4"),
			[]Aggregator{mustNewCountAggregator("numIncomeGt2", "fred > 2")},
			[]*dexpr.Expr{mustNewDExpr("numIncomeGt2 == 1")},
			dexpr.ErrInvalidExpr("Variable doesn't exist: fred")},
		{mustNewDExpr("band > 4"),
			[]Aggregator{mustNewCountAggregator("numIncomeGt2", "income > 2")},
			[]*dexpr.Expr{mustNewDExpr("numIncomeGt == 1")}, nil},
		{mustNewDExpr("hand > 4"),
			[]Aggregator{mustNewCountAggregator("numIncomeGt2", "income > 2")},
			[]*dexpr.Expr{mustNewDExpr("numIncomeGt == 1")},
			dexpr.ErrInvalidExpr("Variable doesn't exist: hand")},
		{mustNewDExpr("band ^^ 4"),
			[]Aggregator{mustNewCountAggregator("numIncomeGt2", "income > 2")},
			[]*dexpr.Expr{mustNewDExpr("numIncomeGt == 1")},
			dexpr.ErrInvalidExpr("Invalid operator: \"^\"")},
	}
	for _, c := range cases {
		ra := NewRuleAssessment(c.rule, c.aggregators, c.goals)
		for _, record := range records {
			err := ra.NextRecord(record)
			if !errorMatch(c.wantErr, err) {
				t.Errorf("NextRecord(%q) rule: %q, aggregators: %q, goals: %q err: %q, wantErr: %q",
					record, c.rule, c.aggregators, c.goals, err, c.wantErr)
			}
		}
	}
}
