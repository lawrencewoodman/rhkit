package aggregators

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/goal"
	"testing"
)

func TestCalcGetResult(t *testing.T) {
	aggregators := []Aggregator{
		MustNew("a", "calc", "3 + 4"),
		MustNew("b", "calc", "5 + 6"),
		MustNew("c", "calc", "a + b"),
		MustNew("2NumRecords", "calc", "numRecords * 2"),
		MustNew("d", "calc", "a + e"),
	}
	goals := []*goal.Goal{}
	want := []*dlit.Literal{
		dlit.MustNew(7),
		dlit.MustNew(11),
		dlit.MustNew(18),
		dlit.MustNew(24),
		dlit.MustNew(dexpr.ErrInvalidExpr{
			Expr: "a + e",
			Err:  dexpr.ErrVarNotExist("e"),
		}),
	}
	numRecords := int64(12)
	for i, aggregator := range aggregators {
		got := aggregator.GetResult(aggregators, goals, numRecords)
		if got.String() != want[i].String() {
			t.Errorf("GetResult() i: %d got: %s, want: %s", i, got, want[i])
		}
	}
}

func TestCalcCloneNew(t *testing.T) {
	aggregators := []Aggregator{
		MustNew("a", "calc", "3 + 4"),
		MustNew("b", "calc", "5 + 6"),
		MustNew("c", "calc", "a + b"),
	}
	goals := []*goal.Goal{}
	numRecords := int64(12)
	aggregatorD := aggregators[2].CloneNew()
	gotC := aggregators[2].GetResult(aggregators, goals, numRecords)
	gotD := aggregatorD.GetResult(aggregators, goals, numRecords)

	if gotC.String() != gotD.String() && gotC.String() != "18" {
		t.Errorf("CloneNew() gotC: %s, gotD: %s", gotC, gotD)
	}
}
