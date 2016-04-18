package internal

import (
	"github.com/lawrencewoodman/dlit_go"
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
	percentCostGt2, err := NewPercentAggregator("percentCostGt2", "cost > 2")
	if err != nil {
		t.Errorf("NewPercentAggregator(\"percentCostGt2\", \"cost > 2\") err == %s",
			err)
	}
	aggregators := []Aggregator{percentCostGt2}

	for i, record := range records {
		percentCostGt2.NextRecord(record, i != 1)
	}
	want := 33.33
	numRecords := int64(len(records))
	got := percentCostGt2.GetResult(aggregators, numRecords)
	gotFloat, gotIsFloat := got.Float()
	if !gotIsFloat || gotFloat != want {
		t.Errorf("GetResult() got: %s, want: %s", got, want)
	}
}

func TestPercentCloneNew(t *testing.T) {
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(3),
		"band":   dlit.MustNew(4),
		"cost":   dlit.MustNew(4),
	}
	numRecords := int64(1)
	percentCostGt2, err := NewPercentAggregator("percentCostGt2", "cost > 2")
	if err != nil {
		t.Errorf("NewPercentAggregator(\"percentCostGt2\", \"cost > 2\") err == %s",
			err)
	}
	percentCostGt2_2 := percentCostGt2.CloneNew()
	aggregators := []Aggregator{}
	want := int64(100)
	percentCostGt2.NextRecord(record, true)
	got1 := percentCostGt2.GetResult(aggregators, numRecords)
	got2 := percentCostGt2_2.GetResult(aggregators, numRecords)

	gotInt1, gotIsInt1 := got1.Int()
	if !gotIsInt1 || gotInt1 != want {
		t.Errorf("GetResult() got: %s, want: %s", gotInt1, want)
	}
	gotInt2, gotIsInt2 := got2.Int()
	if !gotIsInt2 || gotInt2 != 0 {
		t.Errorf("GetResult() got: %s, want: %s", gotInt1, 0)
	}
}
