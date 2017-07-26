package assessment

import (
	"fmt"
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/aggregators"
	"github.com/vlifesystems/rhkit/goal"
	"github.com/vlifesystems/rhkit/internal/testhelpers"
	"github.com/vlifesystems/rhkit/rule"
	"testing"
)

func TestNextRecord(t *testing.T) {
	// It is important for this test to reuse the goals
	// to ensure that they are cloned properly.
	inAggregators := []aggregators.Spec{
		aggregators.MustNew("numIncomeGt2", "count", "income > 2"),
		aggregators.MustNew("numBandGt4", "count", "band > 4"),
		aggregators.MustNew("goalsScore", "goalsscore"),
	}
	records := [4]map[string]*dlit.Literal{
		{
			"income": dlit.MustNew(3),
			"cost":   dlit.MustNew(4.5),
			"band":   dlit.MustNew(4),
		},
		{
			"income": dlit.MustNew(3),
			"cost":   dlit.MustNew(3.2),
			"band":   dlit.MustNew(7),
		},
		{
			"income": dlit.MustNew(2),
			"cost":   dlit.MustNew(1.2),
			"band":   dlit.MustNew(4),
		},
		{
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
		rule             rule.Rule
		wantNumIncomeGt2 int64
		wantNumBandGt4   int64
		wantGoalsScore   float64
	}{
		{rule.NewGEFV("band", dlit.MustNew(5)), 1, 2, 2.0},
		{rule.NewGEFV("band", dlit.MustNew(3)), 2, 2, 0.001},
		{rule.NewGEFV("cost", dlit.MustNew(1.3)), 2, 1, 0},
	}
	for _, c := range cases {
		ra := newRuleAssessor(c.rule, inAggregators, goals)
		for _, record := range records {
			err := ra.NextRecord(record)
			if err != nil {
				t.Errorf("nextRecord(%v) rule: %s, aggregators: %v, goals: %v - err: %v",
					record, c.rule, inAggregators, goals, err)
			}
		}
		gotNumIncomeGt2, gt2Exists :=
			ra.AggregatorValue("numIncomeGt2", numRecords)
		if !gt2Exists {
			t.Errorf("numIncomeGt2 aggregator doesn't exist")
		}
		gotNumIncomeGt2Int, gt2IsInt := gotNumIncomeGt2.Int()
		if !gt2IsInt {
			t.Errorf("numIncomeGt2 aggregator can't be int")
		}
		if gotNumIncomeGt2Int != c.wantNumIncomeGt2 {
			t.Errorf("nextRecord() rule: %s, aggregators: %v, goals: %v - wantNumIncomeGt2: %d, got: %d",
				c.rule, inAggregators, goals, c.wantNumIncomeGt2, gotNumIncomeGt2Int)
		}
		gotNumBandGt4, gt4Exists :=
			ra.AggregatorValue("numBandGt4", numRecords)
		if !gt4Exists {
			t.Errorf("numBandGt4 aggregator doesn't exist")
		}
		gotNumBandGt4Int, gt4IsInt := gotNumBandGt4.Int()
		if !gt4IsInt {
			t.Errorf("numBandGt4 aggregator can't be int")
		}
		if gotNumBandGt4Int != c.wantNumBandGt4 {
			t.Errorf("nextRecord() rule: %s, aggregators: %v, goals: %v - wantNumBandGt4: %d, got: %d",
				c.rule, inAggregators, goals, c.wantNumBandGt4, gotNumBandGt4Int)
		}
		gotGoalsScore, goalsScoreExists :=
			ra.AggregatorValue("goalsScore", numRecords)
		if !goalsScoreExists {
			t.Errorf("goalsScore aggregator doesn't exist")
		}
		gotGoalsScoreFloat, goalsScoreIsFloat := gotGoalsScore.Float()
		if !goalsScoreIsFloat {
			t.Errorf("goalsScore aggregator can't be float")
		}
		if gotGoalsScoreFloat != c.wantGoalsScore {
			t.Errorf("nextRecord() rule: %s, aggregators: %v, goals: %v - wantGoalsScore: %f, got: %f",
				c.rule, inAggregators, goals, c.wantGoalsScore, gotGoalsScore)
		}
	}
}

func TestNextRecord_Errors(t *testing.T) {
	records := [4]map[string]*dlit.Literal{
		{"income": dlit.MustNew(3), "band": dlit.MustNew(4)},
		{"income": dlit.MustNew(3), "band": dlit.MustNew(7)},
		{"income": dlit.MustNew(2), "band": dlit.MustNew(4)},
		{"income": dlit.MustNew(0), "band": dlit.MustNew(9)},
	}
	goals := []*goal.Goal{goal.MustNew("numIncomeGt2 == 1")}
	cases := []struct {
		rule        rule.Rule
		aggregators []aggregators.Spec
		wantErr     error
	}{
		{rule.NewGEFV("band", dlit.MustNew(4)),
			[]aggregators.Spec{
				aggregators.MustNew("numIncomeGt2", "count", "fred > 2")},
			dexpr.InvalidExprError{
				Expr: "fred > 2",
				Err:  dexpr.VarNotExistError("fred"),
			},
		},
		{rule.NewGEFV("band", dlit.MustNew(4)),
			[]aggregators.Spec{
				aggregators.MustNew("numIncomeGt2", "count", "income > 2")}, nil},
		{rule.NewGEFV("hand", dlit.MustNew(4)),
			[]aggregators.Spec{
				aggregators.MustNew("numIncomeGt2", "count", "income > 2")},
			rule.InvalidRuleError{Rule: rule.NewGEFV("hand", dlit.MustNew(4))},
		},
	}
	for _, c := range cases {
		ra := newRuleAssessor(c.rule, c.aggregators, goals)
		for _, record := range records {
			err := ra.NextRecord(record)
			context := fmt.Sprintf(
				"NextRecord(%v) rule: %v, aggregators: %v, goals: %v err: %v, wantErr: %v",
				record, c.rule, c.aggregators, goals, err, c.wantErr,
			)
			testhelpers.CheckErrorMatch(t, context, err, c.wantErr)
		}
	}
}
