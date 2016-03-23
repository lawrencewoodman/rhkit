package internal

import (
	"fmt"
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/dlit_go"
	"testing"
)

func TestCalcGetResult(t *testing.T) {
	aggregators := []Aggregator{
		mustNewCalcAggregator("a", "3 + 4"),
		mustNewCalcAggregator("b", "5 + 6"),
		mustNewCalcAggregator("c", "a + b"),
		mustNewCalcAggregator("2NumRecords", "numRecords * 2"),
		mustNewCalcAggregator("d", "a + e"),
	}
	want := []*dlit.Literal{
		dlit.MustNew(7),
		dlit.MustNew(11),
		dlit.MustNew(18),
		dlit.MustNew(24),
		dlit.MustNew(dexpr.ErrInvalidExpr("Variable doesn't exist: e")),
	}
	numRecords := int64(12)
	for i, aggregator := range aggregators {
		got := aggregator.GetResult(aggregators, numRecords)
		if got.String() != want[i].String() {
			t.Errorf("GetResult() i: %d got: %s, want: %s", i, got, want[i])
		}
	}
}

func TestCalcCloneNew(t *testing.T) {
	aggregators := []Aggregator{
		mustNewCalcAggregator("a", "3 + 4"),
		mustNewCalcAggregator("b", "5 + 6"),
		mustNewCalcAggregator("c", "a + b"),
	}
	numRecords := int64(12)
	aggregatorD := aggregators[2].CloneNew()
	gotC := aggregators[2].GetResult(aggregators, numRecords)
	gotD := aggregatorD.GetResult(aggregators, numRecords)

	if gotC.String() != gotD.String() && gotC.String() != "18" {
		t.Errorf("CloneNew() gotC: %s, gotD: %s", gotC, gotD)
	}
}

/************************
 *   Helper functions
 ************************/
func mustNewCalcAggregator(name string, expr string) *CalcAggregator {
	c, err := NewCalcAggregator(name, expr)
	if err != nil {
		panic(fmt.Sprintf("Can't create CalcAggregator: %s", err))
	}
	return c
}
