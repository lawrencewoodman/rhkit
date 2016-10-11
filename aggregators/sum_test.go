package aggregators

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/goal"
	"testing"
)

func TestSumGetResult(t *testing.T) {
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
	}
	goals := []*goal.Goal{}
	profitDesc, err := New("profit", "sum", "income-cost")
	if err != nil {
		t.Errorf("New(\"profit\", \"sum\", \"income-cost\") err == %s", err)
	}
	profit := profitDesc.New()
	instances := []AggregatorInstance{profit}

	for i, record := range records {
		profit.NextRecord(record, i != 2)
	}
	want := 5.3
	numRecords := int64(len(records))
	got := profit.GetResult(instances, goals, numRecords)
	gotFloat, gotIsFloat := got.Float()
	if !gotIsFloat || gotFloat != want {
		t.Errorf("GetResult() got: %f, want: %f", got, want)
	}
}
