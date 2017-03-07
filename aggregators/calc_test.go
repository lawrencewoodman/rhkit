package aggregators

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/goal"
	"testing"
)

func TestCalcGetResult(t *testing.T) {
	aggregatorSpecs := []AggregatorSpec{
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
		dlit.MustNew(dexpr.InvalidExprError{
			Expr: "a + e",
			Err:  dexpr.VarNotExistError("e"),
		}),
	}
	numRecords := int64(12)
	instances := make([]AggregatorInstance, len(aggregatorSpecs))
	for i, aggregatorSpec := range aggregatorSpecs {
		instances[i] = aggregatorSpec.New()
	}
	for i, instance := range instances {
		got := instance.GetResult(instances, goals, numRecords)
		if got.String() != want[i].String() {
			t.Errorf("GetResult() i: %d got: %s, want: %s", i, got, want[i])
		}
	}
}

func TestCalcSpecGetName(t *testing.T) {
	name := "a"
	as := MustNew(name, "calc", "3+4")
	got := as.GetName()
	if got != name {
		t.Errorf("GetName - got: %s, want: %s", got, name)
	}
}

func TestCalcSpecGetKind(t *testing.T) {
	kind := "calc"
	as := MustNew("a", kind, "3+4")
	got := as.GetKind()
	if got != kind {
		t.Errorf("GetKind - got: %s, want: %s", got, kind)
	}
}

func TestCalcSpecGetArg(t *testing.T) {
	arg := "3+4"
	as := MustNew("a", "calc", arg)
	got := as.GetArg()
	if got != arg {
		t.Errorf("GetArg - got: %s, want: %s", got, arg)
	}
}
