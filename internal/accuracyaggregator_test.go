package internal

import (
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestAccuracyGetResult(t *testing.T) {
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
	goals := []*Goal{}
	cases := []struct {
		records []map[string]*dlit.Literal
		want    float64
	}{
		{records, 44.44},
		{[]map[string]*dlit.Literal{}, 0},
	}
	for _, c := range cases {
		accuracyCostGt2, err := NewAccuracyAggregator("accuracyCostGt2", "cost > 2")
		if err != nil {
			t.Errorf("NewAccuracyAggregator(\"accuracyCostGt2\", \"cost > 2\") err == %s",
				err)
		}
		aggregators := []Aggregator{accuracyCostGt2}

		for i, record := range c.records {
			accuracyCostGt2.NextRecord(record, i != 1 && i != 2)
		}
		numRecords := int64(len(c.records))
		got := accuracyCostGt2.GetResult(aggregators, goals, numRecords)
		gotFloat, gotIsFloat := got.Float()
		if !gotIsFloat || gotFloat != c.want {
			t.Errorf("GetResult() got: %s, want: %s", got, c.want)
		}
	}
}

func TestAccuracyCloneNew(t *testing.T) {
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(3),
		"band":   dlit.MustNew(4),
		"cost":   dlit.MustNew(4),
	}
	goals := []*Goal{}
	numRecords := int64(1)
	accuracyCostGt2, err := NewAccuracyAggregator("accuracyCostGt2", "cost > 2")
	if err != nil {
		t.Errorf("NewAccuracyAggregator(\"accuracyCostGt2\", \"cost > 2\") err == %s",
			err)
	}
	accuracyCostGt2_2 := accuracyCostGt2.CloneNew()
	aggregators := []Aggregator{}
	want := int64(100)
	accuracyCostGt2.NextRecord(record, true)
	got1 := accuracyCostGt2.GetResult(aggregators, goals, numRecords)
	got2 := accuracyCostGt2_2.GetResult(aggregators, goals, numRecords)

	gotInt1, gotIsInt1 := got1.Int()
	if !gotIsInt1 || gotInt1 != want {
		t.Errorf("GetResult() got: %s, want: %s", gotInt1, want)
	}
	gotInt2, gotIsInt2 := got2.Int()
	if !gotIsInt2 || gotInt2 != 0 {
		t.Errorf("GetResult() got: %s, want: %s", gotInt1, 0)
	}
}
