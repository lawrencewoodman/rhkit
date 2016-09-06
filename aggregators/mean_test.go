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
	meanProfitDesc, err := New("meanProfit", "mean", "income-cost")
	if err != nil {
		t.Errorf("New(\"meanProfit\", \"mean\", \"income-cost\") err == %s", err)
	}
	meanProfit := meanProfitDesc.New()
	instances := []AggregatorInstance{meanProfit}

	for i, record := range records {
		meanProfit.NextRecord(record, i != 2)
	}
	want := 2.02
	numRecords := int64(len(records))
	got := meanProfit.GetResult(instances, goals, numRecords)
	gotFloat, gotIsFloat := got.Float()
	if !gotIsFloat || gotFloat != want {
		t.Errorf("GetResult() got: %v, want: %f", got, want)
	}
}

func TestMeanGetResult_dbz(t *testing.T) {
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
	}
	goals := []*goal.Goal{}
	meanProfitDesc, err := New("meanProfit", "mean", "income-cost")
	if err != nil {
		t.Errorf("New(\"meanProfit\", \"mean\", \"income-cost\") err == %s", err)
	}
	meanProfit := meanProfitDesc.New()
	instances := []AggregatorInstance{meanProfit}

	for i, record := range records {
		meanProfit.NextRecord(record, i == 4)
	}
	want := float64(0)
	numRecords := int64(len(records))
	got := meanProfit.GetResult(instances, goals, numRecords)
	gotFloat, gotIsFloat := got.Float()
	if !gotIsFloat || gotFloat != want {
		t.Errorf("GetResult() got: %v, want: %f", got, want)
	}
}
