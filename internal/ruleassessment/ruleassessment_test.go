package ruleassessment

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/aggregators"
	"github.com/vlifesystems/rulehunter/goal"
	"github.com/vlifesystems/rulehunter/rule"
	"testing"
)

func TestNextRecord(t *testing.T) {
	// It is important for this test to reuse the aggregators and goals
	// to ensure that they are cloned properly.
	inAggregators := []aggregators.Aggregator{
		aggregators.MustNew("numIncomeGt2", "count", "income > 2"),
		aggregators.MustNew("numBandGt4", "count", "band > 4"),
		aggregators.MustNew("numGoalsPassed", "goalspassedscore"),
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
	goals := []*goal.Goal{
		goal.MustNew("numIncomeGt2 == 1"),
		goal.MustNew("numBandGt4 == 2"),
	}
	numRecords := int64(len(records))
	cases := []struct {
		rule               *rule.Rule
		wantNumIncomeGt2   int64
		wantNumBandGt4     int64
		wantNumGoalsPassed float64
	}{
		{rule.MustNew("band > 4"), 1, 2, 2.0},
		{rule.MustNew("band > 3"), 2, 2, 0.001},
		{rule.MustNew("cost > 1.2"), 2, 1, 0},
	}
	for _, c := range cases {
		ra := New(c.rule, inAggregators, goals)
		for _, record := range records {
			err := ra.NextRecord(record)
			if err != nil {
				t.Errorf("nextRecord(%q) rule: %s, aggregators: %q, goals: %q - err: %q",
					record, c.rule, inAggregators, goals, err)
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
		if gotNumIncomeGt2Int != c.wantNumIncomeGt2 {
			t.Errorf("nextRecord() rule: %s, aggregators: %q, goals: %q - wantNumIncomeGt2: %d, got: %d",
				c.rule, inAggregators, goals, c.wantNumIncomeGt2, gotNumIncomeGt2Int)
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
		if gotNumBandGt4Int != c.wantNumBandGt4 {
			t.Errorf("nextRecord() rule: %s, aggregators: %q, goals: %q - wantNumBandGt4: %d, got: %d",
				c.rule, inAggregators, goals, c.wantNumBandGt4, gotNumBandGt4Int)
		}
		gotNumGoalsPassed, goalsPassedExists :=
			ra.GetAggregatorValue("numGoalsPassed", numRecords)
		if !goalsPassedExists {
			t.Errorf("numGoalsPassed aggregator doesn't exist")
		}
		gotNumGoalsPassedFloat, goalsPassedIsFloat := gotNumGoalsPassed.Float()
		if !goalsPassedIsFloat {
			t.Errorf("numGoalsPassed aggregator can't be float")
		}
		if gotNumGoalsPassedFloat != c.wantNumGoalsPassed {
			t.Errorf("nextRecord() rule: %s, aggregators: %q, goals: %q - wantNumGoalsPassed: %d, got: %d",
				c.rule, inAggregators, goals, c.wantNumGoalsPassed, gotNumGoalsPassed)
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
	goals := []*goal.Goal{goal.MustNew("numIncomeGt2 == 1")}
	cases := []struct {
		rule        *rule.Rule
		aggregators []aggregators.Aggregator
		wantErr     error
	}{
		{rule.MustNew("band > 4"),
			[]aggregators.Aggregator{
				aggregators.MustNew("numIncomeGt2", "count", "fred > 2")},
			dexpr.ErrInvalidExpr("Variable doesn't exist: fred")},
		{rule.MustNew("band > 4"),
			[]aggregators.Aggregator{
				aggregators.MustNew("numIncomeGt2", "count", "income > 2")}, nil},
		{rule.MustNew("hand > 4"),
			[]aggregators.Aggregator{
				aggregators.MustNew("numIncomeGt2", "count", "income > 2")},
			dexpr.ErrInvalidExpr("Variable doesn't exist: hand")},
		{rule.MustNew("band ^^ 4"),
			[]aggregators.Aggregator{
				aggregators.MustNew("numIncomeGt2", "count", "income > 2")},
			dexpr.ErrInvalidExpr("Invalid operator: \"^\"")},
	}
	for _, c := range cases {
		ra := New(c.rule, c.aggregators, goals)
		for _, record := range records {
			err := ra.NextRecord(record)
			if !errorMatch(c.wantErr, err) {
				t.Errorf("NextRecord(%q) rule: %q, aggregators: %q, goals: %q err: %q, wantErr: %q",
					record, c.rule, c.aggregators, goals, err, c.wantErr)
				return
			}
		}
	}
}
