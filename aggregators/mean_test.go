package aggregators

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/goal"
	"testing"
)

func TestNewMean_error(t *testing.T) {
	_, err := New("a", "mean", "3+4+{")
	wantErr := "can't make aggregator: a, error: " +
		dexpr.InvalidExprError{
			Expr: "3+4+{",
			Err:  dexpr.ErrSyntax,
		}.Error()
	if err.Error() != wantErr {
		t.Errorf("New: gotErr: %s, wantErr: %s", err, wantErr)
	}
}

func TestMeanResult(t *testing.T) {
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
			"income": dlit.MustNew(2),
			"cost":   dlit.MustNew(1.2),
			"band":   dlit.MustNew(4),
		},
		{
			"income": dlit.MustNew(9),
			"cost":   dlit.MustNew(2),
			"band":   dlit.MustNew(9),
		},
		{
			"income": dlit.MustNew(3.98),
			"cost":   dlit.MustNew(1.2),
			"band":   dlit.MustNew(9),
		},
	}
	goals := []*goal.Goal{}
	cases := []struct {
		records []map[string]*dlit.Literal
		rule    func(int) bool
		want    float64
	}{
		{records, func(i int) bool { return i != 2 }, 2.02},
		{records, func(i int) bool { return false }, 0},
		{[]map[string]*dlit.Literal{}, func(i int) bool { return true }, 0},
	}
	for _, c := range cases {
		meanProfitDesc := MustNew("meanProfit", "mean", "income-cost")
		meanProfit := meanProfitDesc.New()
		instances := []AggregatorInstance{meanProfit}

		for i, record := range c.records {
			meanProfit.NextRecord(record, c.rule(i))
		}
		numRecords := int64(len(records))
		got := meanProfit.Result(instances, goals, numRecords)
		gotFloat, gotIsFloat := got.Float()
		if !gotIsFloat || gotFloat != c.want {
			t.Errorf("Result() got: %v, want: %f", got, c.want)
		}
	}
}

func TestMeanNextRecord_errors(t *testing.T) {
	as := MustNew("a", "mean", "cost + 2")
	ai := as.New()
	cases := []struct {
		record map[string]*dlit.Literal
		want   error
	}{
		{record: map[string]*dlit.Literal{},
			want: dexpr.InvalidExprError{
				Expr: "cost + 2",
				Err:  dexpr.VarNotExistError("cost"),
			},
		},
		{record: map[string]*dlit.Literal{"cost": dlit.NewString("hello")},
			want: dexpr.InvalidExprError{
				Expr: "cost + 2",
				Err:  dexpr.ErrIncompatibleTypes,
			},
		},
	}
	for _, c := range cases {
		got := ai.NextRecord(c.record, true)
		if got == nil || got.Error() != c.want.Error() {
			t.Errorf("NextRecord: got: %s, want: %s", got, c.want)
		}
	}
}

func TestMeanSpecName(t *testing.T) {
	name := "a"
	as := MustNew(name, "mean", "income - cost")
	got := as.Name()
	if got != name {
		t.Errorf("Name - got: %s, want: %s", got, name)
	}
}

func TestMeanSpecKind(t *testing.T) {
	kind := "mean"
	as := MustNew("a", kind, "income - cost")
	got := as.Kind()
	if got != kind {
		t.Errorf("Kind - got: %s, want: %s", got, kind)
	}
}

func TestMeanSpecArg(t *testing.T) {
	arg := "income - cost"
	as := MustNew("a", "mean", arg)
	got := as.Arg()
	if got != arg {
		t.Errorf("Arg - got: %s, want: %s", got, arg)
	}
}

func TestMeanInstanceName(t *testing.T) {
	as := MustNew("abc", "mean", "cost + 2")
	ai := as.New()
	got := ai.Name()
	want := "abc"
	if got != want {
		t.Errorf("Name: got: %s, want: %s", got, want)
	}
}
