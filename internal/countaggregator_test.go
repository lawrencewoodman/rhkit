package internal

import (
	"github.com/lawrencewoodman/dlit_go"
	"testing"
)

func TestCountGetResult(t *testing.T) {
	records := [4]map[string]*dlit.Literal{
		map[string]*dlit.Literal{"income": dlit.MustNew(3), "band": dlit.MustNew(4)},
		map[string]*dlit.Literal{"income": dlit.MustNew(3), "band": dlit.MustNew(7)},
		map[string]*dlit.Literal{"income": dlit.MustNew(2), "band": dlit.MustNew(4)},
		map[string]*dlit.Literal{"income": dlit.MustNew(0), "band": dlit.MustNew(9)},
	}
	numBandGt4, err := NewCountAggregator("numBandGt4", "band > 4")
	if err != nil {
		t.Errorf("NewCount(\"numBandGt4\", \"band > 4\") err == %q",
			err)
	}
	aggregators := []Aggregator{numBandGt4}

	for _, record := range records {
		numBandGt4.NextRecord(record, true)
	}
	numRecords := int64(len(records))
	got := numBandGt4.GetResult(aggregators, numRecords)
	gotInt, gotIsInt := got.Int()
	if !gotIsInt || gotInt != 2 {
		t.Errorf("NewCount(\"numBandGt4\", \"band > 4\") == %q, want: %q",
			got, 2)
	}
}

func TestCountCloneNew(t *testing.T) {
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(3),
		"band":   dlit.MustNew(4),
	}
	numRecords := int64(1)
	numBandGt4, err := NewCountAggregator("numBandGt4", "band > 3")
	if err != nil {
		t.Errorf("NewCount(\"numBandGt4\", \"band > 3\") err == %q",
			err)
	}
	numBandGt4_2 := numBandGt4.CloneNew()
	aggregators := []Aggregator{}

	numBandGt4.NextRecord(record, true)
	got1 := numBandGt4.GetResult(aggregators, numRecords)
	got2 := numBandGt4_2.GetResult(aggregators, numRecords)

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

/*
TODO:
test compile-count-2 {Ensure that ignores false isRuleTrue rules} -setup {
  set prospects {{3 4} {3 7} {2 4} {0 9}}
  set fields {income band}
  set numBandGt4NewCmd [aggregator compile $fields {count {$band > 4}}]
  set numBandGt4 [{*}$numBandGt4NewCmd]
  foreach prospect $prospects {
    $numBandGt4 nextProspect $prospect true
  }
  $numBandGt4 nextProspect {50 100} false
} -body {
  $numBandGt4 getResult {}
} -result {2}
*/
