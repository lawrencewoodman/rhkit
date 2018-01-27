package aggregator

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/goal"
	"testing"
)

func TestNewMean_error(t *testing.T) {
	_, err := New("a", "mean", "3+4+{")
	wantErr := DescError{
		Name: "a",
		Kind: "mean",
		Err: dexpr.InvalidExprError{
			Expr: "3+4+{",
			Err:  dexpr.ErrSyntax,
		},
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
		{
			"income": dlit.MustNew(9),
			"cost":   dlit.MustNew(9),
			"band":   dlit.MustNew(9),
		},
		{
			"income": dlit.MustNew(11),
			"cost":   dlit.MustNew(11),
			"band":   dlit.MustNew(9),
		},
		{
			"income": dlit.MustNew(9),
			"cost":   dlit.MustNew(8),
			"band":   dlit.MustNew(9),
		},
	}
	goals := []*goal.Goal{}
	cases := []struct {
		records []map[string]*dlit.Literal
		rule    func(int) bool
		want    float64
	}{
		{records, func(i int) bool { return i != 2 && i <= 4 }, 2.02},
		{records, func(i int) bool { return i > 4 }, 0.33},
		{records, func(i int) bool { return false }, 0},
		{[]map[string]*dlit.Literal{}, func(i int) bool { return true }, 0},
	}
	for ci, c := range cases {
		meanProfitDesc := MustNew("meanProfit", "mean", "income-cost")
		meanProfit := meanProfitDesc.New()
		instances := []Instance{meanProfit}

		for i, record := range c.records {
			meanProfit.NextRecord(record, c.rule(i))
		}
		numRecords := int64(len(records))
		got := meanProfit.Result(instances, goals, numRecords)
		gotFloat, gotIsFloat := got.Float()
		if !gotIsFloat || gotFloat != c.want {
			t.Errorf("(%d) - Result, got: %v, want: %f", ci, got, c.want)
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

/*************************
 *       Benchmarks
 *************************/

func BenchmarkMeanNextRecord(b *testing.B) {
	as := MustNew("a", "mean", "cost + 2")
	ai := as.New()
	record := map[string]*dlit.Literal{"cost": dlit.NewString("17.89245")}
	for n := 0; n < b.N; n++ {
		b.StartTimer()
		got := ai.NextRecord(record, true)
		b.StopTimer()
		if got != nil {
			b.Errorf("NextRecord: %s", got)
		}
	}
}
