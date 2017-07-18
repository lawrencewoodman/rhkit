package aggregators

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/goal"
	"testing"
)

func TestNewPrecision_error(t *testing.T) {
	_, err := New("a", "precision", "3>4{")
	wantErr := "can't make aggregator: a, error: " +
		dexpr.InvalidExprError{
			Expr: "3>4{",
			Err:  dexpr.ErrSyntax,
		}.Error()
	if err.Error() != wantErr {
		t.Errorf("New: gotErr: %s, wantErr: %s", err, wantErr)
	}
}

func TestPrecisionResult(t *testing.T) {
	records := []map[string]*dlit.Literal{
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
			"income": dlit.MustNew(3),
			"cost":   dlit.MustNew(1.2),
			"band":   dlit.MustNew(7),
		},
		{
			"income": dlit.MustNew(2),
			"cost":   dlit.MustNew(1.2),
			"band":   dlit.MustNew(4),
		},
		{
			"income": dlit.MustNew(2),
			"cost":   dlit.MustNew(5.6),
			"band":   dlit.MustNew(4),
		},
		{
			"income": dlit.MustNew(2),
			"cost":   dlit.MustNew(0.6),
			"band":   dlit.MustNew(4),
		},
		{
			"income": dlit.MustNew(2),
			"cost":   dlit.MustNew(0.8),
			"band":   dlit.MustNew(4),
		},
		{
			"income": dlit.MustNew(9),
			"cost":   dlit.MustNew(2),
			"band":   dlit.MustNew(9),
		},
		{
			"income": dlit.MustNew(9),
			"cost":   dlit.MustNew(3),
			"band":   dlit.MustNew(9),
		},
	}
	goals := []*goal.Goal{}
	cases := []struct {
		records []map[string]*dlit.Literal
		rule    func(int) bool
		want    float64
	}{
		{records, func(i int) bool { return i != 1 && i != 2 }, 0.4286},
		{records, func(i int) bool { return false }, 0},
		{[]map[string]*dlit.Literal{}, func(i int) bool { return true }, 0},
	}
	for _, c := range cases {
		precisionCostGt2Desc := MustNew("precisionCostGt2", "precision", "cost > 2")
		for i := 0; i < 5; i++ {
			precisionCostGt2 := precisionCostGt2Desc.New()
			instances := []AggregatorInstance{precisionCostGt2}

			for i, record := range c.records {
				precisionCostGt2.NextRecord(record, c.rule(i))
			}
			numRecords := int64(len(c.records))
			got := precisionCostGt2.Result(instances, goals, numRecords)
			gotFloat, gotIsFloat := got.Float()
			if !gotIsFloat || gotFloat != c.want {
				t.Errorf("Result() got: %v, want: %v", got, c.want)
			}
		}
	}
}

func TestPrecisionNextRecord_error(t *testing.T) {
	as := MustNew("a", "precision", "cost > 2")
	ai := as.New()
	record := map[string]*dlit.Literal{}
	got := ai.NextRecord(record, true)
	want := dexpr.InvalidExprError{
		Expr: "cost > 2",
		Err:  dexpr.VarNotExistError("cost"),
	}
	if got == nil || got.Error() != want.Error() {
		t.Errorf("NextRecord: got: %s, want: %s", got, want)
	}
}

func TestPrecisionSpecName(t *testing.T) {
	name := "a"
	as := MustNew(name, "precision", "cost > 2")
	got := as.Name()
	if got != name {
		t.Errorf("Name - got: %s, want: %s", got, name)
	}
}

func TestPrecisionSpecKind(t *testing.T) {
	kind := "precision"
	as := MustNew("a", kind, "cost > 2")
	got := as.Kind()
	if got != kind {
		t.Errorf("Kind - got: %s, want: %s", got, kind)
	}
}

func TestPrecisionSpecArg(t *testing.T) {
	arg := "cost > 2"
	as := MustNew("a", "precision", arg)
	got := as.Arg()
	if got != arg {
		t.Errorf("Arg - got: %s, want: %s", got, arg)
	}
}

func TestPrecisionInstanceName(t *testing.T) {
	as := MustNew("abc", "precision", "cost + 2")
	ai := as.New()
	got := ai.Name()
	want := "abc"
	if got != want {
		t.Errorf("Name: got: %s, want: %s", got, want)
	}
}
