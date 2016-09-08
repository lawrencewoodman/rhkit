package aggregators

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/goal"
	"testing"
)

func TestMeanGetResult(t *testing.T) {
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
			"income": dlit.MustNew(9),
			"cost":   dlit.MustNew(2),
			"band":   dlit.MustNew(9),
		},
		map[string]*dlit.Literal{
			"income": dlit.MustNew(3.98),
			"cost":   dlit.MustNew(1.2),
			"band":   dlit.MustNew(9),
		},
	}
	goals := []*goal.Goal{}
	cases := []struct {
		records []map[string]*dlit.Literal
		want    float64
	}{
		{records, 2.02},
		{[]map[string]*dlit.Literal{}, 0},
	}
	for _, c := range cases {
		meanProfitDesc, err := New("meanProfit", "mean", "income-cost")
		if err != nil {
			t.Errorf("New(\"meanProfit\", \"mean\", \"income-cost\") err == %s", err)
		}
		meanProfit := meanProfitDesc.New()
		instances := []AggregatorInstance{meanProfit}

		for i, record := range c.records {
			meanProfit.NextRecord(record, i != 2)
		}
		numRecords := int64(len(records))
		got := meanProfit.GetResult(instances, goals, numRecords)
		gotFloat, gotIsFloat := got.Float()
		if !gotIsFloat || gotFloat != c.want {
			t.Errorf("GetResult() got: %v, want: %f", got, c.want)
		}
	}
}
