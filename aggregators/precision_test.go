package aggregators

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/goal"
	"testing"
)

func TestPrecisionGetResult(t *testing.T) {
	records := []map[string]*dlit.Literal{
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
			"income": dlit.MustNew(3),
			"cost":   dlit.MustNew(1.2),
			"band":   dlit.MustNew(7),
		},
		map[string]*dlit.Literal{
			"income": dlit.MustNew(2),
			"cost":   dlit.MustNew(1.2),
			"band":   dlit.MustNew(4),
		},
		map[string]*dlit.Literal{
			"income": dlit.MustNew(2),
			"cost":   dlit.MustNew(5.6),
			"band":   dlit.MustNew(4),
		},
		map[string]*dlit.Literal{
			"income": dlit.MustNew(2),
			"cost":   dlit.MustNew(0.6),
			"band":   dlit.MustNew(4),
		},
		map[string]*dlit.Literal{
			"income": dlit.MustNew(2),
			"cost":   dlit.MustNew(0.8),
			"band":   dlit.MustNew(4),
		},
		map[string]*dlit.Literal{
			"income": dlit.MustNew(9),
			"cost":   dlit.MustNew(2),
			"band":   dlit.MustNew(9),
		},
		map[string]*dlit.Literal{
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
			got := precisionCostGt2.GetResult(instances, goals, numRecords)
			gotFloat, gotIsFloat := got.Float()
			if !gotIsFloat || gotFloat != c.want {
				t.Errorf("GetResult() got: %v, want: %v", got, c.want)
			}
		}
	}
}

func TestPrecisionSpecGetName(t *testing.T) {
	name := "a"
	as := MustNew(name, "precision", "cost > 2")
	got := as.GetName()
	if got != name {
		t.Errorf("GetName - got: %s, want: %s", got, name)
	}
}

func TestPrecisionSpecGetKind(t *testing.T) {
	kind := "precision"
	as := MustNew("a", kind, "cost > 2")
	got := as.GetKind()
	if got != kind {
		t.Errorf("GetKind - got: %s, want: %s", got, kind)
	}
}

func TestPrecisionSpecGetArg(t *testing.T) {
	arg := "cost > 2"
	as := MustNew("a", "precision", arg)
	got := as.GetArg()
	if got != arg {
		t.Errorf("GetArg - got: %s, want: %s", got, arg)
	}
}
