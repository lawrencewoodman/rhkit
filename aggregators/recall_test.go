package aggregators

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/goal"
	"testing"
)

func TestNewRecall_error(t *testing.T) {
	_, err := New("a", "recall", "3>4{")
	wantErr := "can't make aggregator: a, error: " +
		dexpr.InvalidExprError{
			Expr: "3>4{",
			Err:  dexpr.ErrSyntax,
		}.Error()
	if err.Error() != wantErr {
		t.Errorf("New: gotErr: %s, wantErr: %s", err, wantErr)
	}
}

func TestRecallResult(t *testing.T) {
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
		want    float64
	}{
		{records, 0.75},
		{records[3:4], 0},
		{[]map[string]*dlit.Literal{}, 0},
	}
	for _, c := range cases {
		recallCostGt2Desc := MustNew("recallCostGt2", "recall", "cost > 2")
		for i := 0; i < 5; i++ {
			recallCostGt2 := recallCostGt2Desc.New()
			instances := []AggregatorInstance{recallCostGt2}

			for i, record := range c.records {
				recallCostGt2.NextRecord(record, i != 1 && i != 2)
			}
			numRecords := int64(len(c.records))
			got := recallCostGt2.Result(instances, goals, numRecords)
			gotFloat, gotIsFloat := got.Float()
			if !gotIsFloat || gotFloat != c.want {
				t.Errorf("Result() got: %v, want: %v", got, c.want)
			}
		}
	}
}

func TestRecallNextRecord_error(t *testing.T) {
	as := MustNew("a", "recall", "cost > 2")
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

func TestRecallSpecName(t *testing.T) {
	name := "a"
	as := MustNew(name, "recall", "cost > 2")
	got := as.Name()
	if got != name {
		t.Errorf("Name - got: %s, want: %s", got, name)
	}
}

func TestRecallSpecKind(t *testing.T) {
	kind := "recall"
	as := MustNew("a", kind, "cost > 2")
	got := as.Kind()
	if got != kind {
		t.Errorf("Kind - got: %s, want: %s", got, kind)
	}
}

func TestRecallSpecArg(t *testing.T) {
	arg := "cost > 2"
	as := MustNew("a", "recall", arg)
	got := as.Arg()
	if got != arg {
		t.Errorf("Arg - got: %s, want: %s", got, arg)
	}
}

func TestRecallInstanceName(t *testing.T) {
	as := MustNew("abc", "recall", "cost + 2")
	ai := as.New()
	got := ai.Name()
	want := "abc"
	if got != want {
		t.Errorf("Name: got: %s, want: %s", got, want)
	}
}
