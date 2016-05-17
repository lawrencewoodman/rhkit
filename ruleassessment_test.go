package rulehunter

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/internal"
	"testing"
)

func TestNextRecord(t *testing.T) {
	// It is important for this test to reuse the aggregators and goals
	// to ensure that they are cloned properly.
	inAggregators := []internal.Aggregator{
		internal.MustNewCountAggregator("numIncomeGt2", "income > 2"),
		internal.MustNewCountAggregator("numBandGt4", "band > 4"),
		mustNewGoalsPassedScoreAggregator("numGoalsPassed"),
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
	goals := []*internal.Goal{
		internal.MustNewGoal("numIncomeGt2 == 1"),
		internal.MustNewGoal("numBandGt4 == 2"),
	}
	numRecords := int64(len(records))
	cases := []struct {
		rule               *Rule
		wantNumIncomeGt2   int64
		wantNumBandGt4     int64
		wantNumGoalsPassed float64
	}{
		{mustNewRule("band > 4"), 1, 2, 2.0},
		{mustNewRule("band > 3"), 2, 2, 0.001},
		{mustNewRule("cost > 1.2"), 2, 1, 0},
	}
	for _, c := range cases {
		ra := newRuleAssessment(c.rule, inAggregators, goals)
		for _, record := range records {
			err := ra.nextRecord(record)
			if err != nil {
				t.Errorf("nextRecord(%q) rule: %s, aggregators: %q, goals: %q - err: %q",
					record, c.rule, inAggregators, goals, err)
			}
		}
		gotNumIncomeGt2, gt2Exists :=
			ra.getAggregatorValue("numIncomeGt2", numRecords)
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
			ra.getAggregatorValue("numBandGt4", numRecords)
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
			ra.getAggregatorValue("numGoalsPassed", numRecords)
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
	goals := []*internal.Goal{internal.MustNewGoal("numIncomeGt2 == 1")}
	cases := []struct {
		rule        *Rule
		aggregators []internal.Aggregator
		wantErr     error
	}{
		{mustNewRule("band > 4"),
			[]internal.Aggregator{
				internal.MustNewCountAggregator("numIncomeGt2", "fred > 2")},
			dexpr.ErrInvalidExpr("Variable doesn't exist: fred")},
		{mustNewRule("band > 4"),
			[]internal.Aggregator{
				internal.MustNewCountAggregator("numIncomeGt2", "income > 2")}, nil},
		{mustNewRule("hand > 4"),
			[]internal.Aggregator{
				internal.MustNewCountAggregator("numIncomeGt2", "income > 2")},
			dexpr.ErrInvalidExpr("Variable doesn't exist: hand")},
		{mustNewRule("band ^^ 4"),
			[]internal.Aggregator{
				internal.MustNewCountAggregator("numIncomeGt2", "income > 2")},
			dexpr.ErrInvalidExpr("Invalid operator: \"^\"")},
	}
	for _, c := range cases {
		ra := newRuleAssessment(c.rule, c.aggregators, goals)
		for _, record := range records {
			err := ra.nextRecord(record)
			if !errorMatch(c.wantErr, err) {
				t.Errorf("NextRecord(%q) rule: %q, aggregators: %q, goals: %q err: %q, wantErr: %q",
					record, c.rule, c.aggregators, goals, err, c.wantErr)
				return
			}
		}
	}
}
