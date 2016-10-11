package aggregators

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/goal"
	"testing"
)

func TestRecallGetResult(t *testing.T) {
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
		want    float64
	}{
		{records, 0.75},
		{records[3:4], 0},
		{[]map[string]*dlit.Literal{}, 0},
	}
	for _, c := range cases {
		recallCostGt2Desc, err := New("recallCostGt2", "recall", "cost > 2")
		if err != nil {
			t.Errorf("New(\"recallCostGt2\", \"recall\", \"cost > 2\") err == %s",
				err)
		}
		for i := 0; i < 5; i++ {
			recallCostGt2 := recallCostGt2Desc.New()
			instances := []AggregatorInstance{recallCostGt2}

			for i, record := range c.records {
				recallCostGt2.NextRecord(record, i != 1 && i != 2)
			}
			numRecords := int64(len(c.records))
			got := recallCostGt2.GetResult(instances, goals, numRecords)
			gotFloat, gotIsFloat := got.Float()
			if !gotIsFloat || gotFloat != c.want {
				t.Errorf("GetResult() got: %v, want: %v", got, c.want)
			}
		}
	}
}
