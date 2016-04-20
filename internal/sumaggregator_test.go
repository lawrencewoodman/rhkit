package internal

import (
	"github.com/lawrencewoodman/dlit_go"
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
	goals := []*Goal{}
	profit, err := NewSumAggregator("profit", "income-cost")
	if err != nil {
		t.Errorf("NewSumAggregator(\"profit\", \"income-cost\") err == %s", err)
	}
	aggregators := []Aggregator{profit}

	for i, record := range records {
		profit.NextRecord(record, i != 2)
	}
	want := 5.3
	numRecords := int64(len(records))
	got := profit.GetResult(aggregators, goals, numRecords)
	gotFloat, gotIsFloat := got.Float()
	if !gotIsFloat || gotFloat != want {
		t.Errorf("GetResult() got: %s, want: %s", got, want)
	}
}

func TestSumCloneNew(t *testing.T) {
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(3),
		"band":   dlit.MustNew(4),
	}
	goals := []*Goal{}
	numRecords := int64(1)
	totalIncome, err := NewSumAggregator("totalIncome", "income")
	if err != nil {
		t.Errorf("NewSumAggregator(\"totalIncome\", \"income\") err: %s", err)
	}
	totalIncome_2 := totalIncome.CloneNew()
	aggregators := []Aggregator{}
	want := int64(3)
	totalIncome.NextRecord(record, true)
	got1 := totalIncome.GetResult(aggregators, goals, numRecords)
	got2 := totalIncome_2.GetResult(aggregators, goals, numRecords)

	gotInt1, gotIsInt1 := got1.Int()
	if !gotIsInt1 || gotInt1 != want {
		t.Errorf("GetResult() got: %s, want: %s", gotInt1, want)
	}
	gotInt2, gotIsInt2 := got2.Int()
	if !gotIsInt2 || gotInt2 != 0 {
		t.Errorf("GetResult() got: %s, want: %s", gotInt1, 0)
	}
}
