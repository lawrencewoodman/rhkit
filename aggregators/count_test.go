package aggregators

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/goal"
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
	goals := []*goal.Goal{}
	numBandGt4Desc, err := New("numBandGt4", "count", "band > 4")
	if err != nil {
		t.Errorf("New(\"numBandGt4\", \"count\", \"band > 4\") err: %v", err)
	}
	numBandGt4 := numBandGt4Desc.New()
	instances := []AggregatorInstance{numBandGt4}

	for i, record := range records {
		numBandGt4.NextRecord(record, i != 3)
	}
	numRecords := int64(len(records))
	want := int64(2)
	got := numBandGt4.GetResult(instances, goals, numRecords)
	gotInt, gotIsInt := got.Int()
	if !gotIsInt || gotInt != want {
		t.Errorf("New(\"numBandGt4\", \"count\", \"band > 4\") got: %v, want: %v",
			got, want)
	}
}
