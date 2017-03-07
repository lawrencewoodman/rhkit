package aggregators

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/goal"
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
	numBandGt4Desc := MustNew("numBandGt4", "count", "band > 4")
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

func TestCountSpecGetName(t *testing.T) {
	name := "a"
	as := MustNew(name, "count", "band > 4")
	got := as.GetName()
	if got != name {
		t.Errorf("GetName - got: %s, want: %s", got, name)
	}
}

func TestCountSpecGetKind(t *testing.T) {
	kind := "count"
	as := MustNew("a", kind, "band > 4")
	got := as.GetKind()
	if got != kind {
		t.Errorf("GetKind - got: %s, want: %s", got, kind)
	}
}

func TestCountSpecGetArg(t *testing.T) {
	arg := "band > 4"
	as := MustNew("a", "count", arg)
	got := as.GetArg()
	if got != arg {
		t.Errorf("GetArg - got: %s, want: %s", got, arg)
	}
}
