package aggregators

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/goal"
	"testing"
)

func TestNewCalc_error(t *testing.T) {
	_, err := New("a", "calc", "3+4+{")
	wantErr := DescError{
		Name: "a",
		Kind: "calc",
		Err: dexpr.InvalidExprError{
			Expr: "3+4+{",
			Err:  dexpr.ErrSyntax,
		},
	}.Error()
	if err.Error() != wantErr {
		t.Errorf("New: gotErr: %s, wantErr: %s", err, wantErr)
	}
}

func TestCalcNextRecord(t *testing.T) {
	as := MustNew("a", "calc", "3+4")
	ai := as.New()
	record := map[string]*dlit.Literal{}
	got := ai.NextRecord(record, true)
	if got != nil {
		t.Errorf("NextRecord: got: %s, want: nil", got)
	}
}

func TestCalcResult(t *testing.T) {
	aggregatorSpecs := []AggregatorSpec{
		MustNew("a", "calc", "3 + 4"),
		MustNew("b", "calc", "5 + 6"),
		MustNew("c", "calc", "a + b"),
		MustNew("2NumRecords", "calc", "numRecords * 2"),
		MustNew("d", "calc", "a + e"),
		MustNew("f", "calc", "a + d"),
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
		got := instance.Result(instances, goals, numRecords)
		if got.String() != want[i].String() {
			t.Errorf("(%d) Result: got: %s, want: %s", i, got, want[i])
		}
	}
}

func TestCalcSpecName(t *testing.T) {
	name := "a"
	as := MustNew(name, "calc", "3+4")
	got := as.Name()
	if got != name {
		t.Errorf("Name - got: %s, want: %s", got, name)
	}
}

func TestCalcSpecKind(t *testing.T) {
	kind := "calc"
	as := MustNew("a", kind, "3+4")
	got := as.Kind()
	if got != kind {
		t.Errorf("Kind - got: %s, want: %s", got, kind)
	}
}

func TestCalcSpecArg(t *testing.T) {
	arg := "3+4"
	as := MustNew("a", "calc", arg)
	got := as.Arg()
	if got != arg {
		t.Errorf("Arg - got: %s, want: %s", got, arg)
	}
}
