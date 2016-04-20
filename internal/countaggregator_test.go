package internal

import (
	"github.com/lawrencewoodman/dlit_go"
	"testing"
)

func TestCountGetResult(t *testing.T) {
	records := []map[string]*dlit.Literal{
		map[string]*dlit.Literal{"income": dlit.MustNew(3), "band": dlit.MustNew(4)},
		map[string]*dlit.Literal{"income": dlit.MustNew(3), "band": dlit.MustNew(7)},
		map[string]*dlit.Literal{"income": dlit.MustNew(2), "band": dlit.MustNew(4)},
		map[string]*dlit.Literal{"income": dlit.MustNew(2), "band": dlit.MustNew(6)},
		map[string]*dlit.Literal{"income": dlit.MustNew(0), "band": dlit.MustNew(9)},
	}
	goals := []*Goal{}
	numBandGt4, err := NewCountAggregator("numBandGt4", "band > 4")
	if err != nil {
		t.Errorf("NewCount(\"numBandGt4\", \"band > 4\") err == %q",
			err)
	}
	aggregators := []Aggregator{numBandGt4}

	for i, record := range records {
		numBandGt4.NextRecord(record, i != 3)
	}
	numRecords := int64(len(records))
	want := int64(2)
	got := numBandGt4.GetResult(aggregators, goals, numRecords)
	gotInt, gotIsInt := got.Int()
	if !gotIsInt || gotInt != want {
		t.Errorf("NewCount(\"numBandGt4\", \"band > 4\") == %q, want: %q",
			got, want)
	}
}

func TestCountCloneNew(t *testing.T) {
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(3),
		"band":   dlit.MustNew(4),
	}
	goals := []*Goal{}
	numRecords := int64(1)
	numBandGt4, err := NewCountAggregator("numBandGt4", "band > 3")
	if err != nil {
		t.Errorf("NewCount(\"numBandGt4\", \"band > 3\") err == %q",
			err)
	}
	numBandGt4_2 := numBandGt4.CloneNew()
	aggregators := []Aggregator{}

	numBandGt4.NextRecord(record, true)
	got1 := numBandGt4.GetResult(aggregators, goals, numRecords)
	got2 := numBandGt4_2.GetResult(aggregators, goals, numRecords)

	gotInt1, gotIsInt1 := got1.Int()
	if !gotIsInt1 || gotInt1 != 1 {
		t.Errorf("NewCount(\"numBandGt4\", \"band > 4\") == %q, want: %q",
			gotInt1, 1)
	}
	gotInt2, gotIsInt2 := got2.Int()
	if !gotIsInt2 || gotInt2 != 0 {
		t.Errorf("NewCount(\"numBandGt4\", \"band > 4\") == %q, want: %q",
			gotInt2, 0)
	}
}
