package aggregators

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/goal"
	"testing"
)

func TestPercentGetResult(t *testing.T) {
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
	}
	goals := []*goal.Goal{}
	cases := []struct {
		records []map[string]*dlit.Literal
		rule    func(int) bool
		want    float64
	}{
		{records, func(i int) bool { return i != 1 }, 33.33},
		{records, func(i int) bool { return false }, 0},
		{[]map[string]*dlit.Literal{}, func(i int) bool { return true }, 0},
	}
	for _, c := range cases {
		percentCostGt2Desc, err := New("percentCostGt2", "percent", "cost > 2")
		if err != nil {
			t.Errorf("New(\"percentCostGt2\", \"percent\", \"cost > 2\") err == %s",
				err)
		}
		percentCostGt2 := percentCostGt2Desc.New()
		instances := []AggregatorInstance{percentCostGt2}

		for i, record := range c.records {
			percentCostGt2.NextRecord(record, c.rule(i))
		}
		numRecords := int64(len(c.records))
		got := percentCostGt2.GetResult(instances, goals, numRecords)
		gotFloat, gotIsFloat := got.Float()
		if !gotIsFloat || gotFloat != c.want {
			t.Errorf("GetResult() got: %v, want: %f", got, c.want)
		}
	}
}
