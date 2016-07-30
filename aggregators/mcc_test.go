package aggregators

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/goal"
	"testing"
)

func TestMCCGetResult(t *testing.T) {
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
		map[string]*dlit.Literal{
			"income": dlit.MustNew(9),
			"cost":   dlit.MustNew(2),
			"band":   dlit.MustNew(7),
		},
	}
	goals := []*goal.Goal{}
	cases := []struct {
		records   []map[string]*dlit.Literal
		ruleExpr  *dexpr.Expr
		checkExpr *dexpr.Expr
	}{
		{records: records,
			ruleExpr:  dexpr.MustNew("cost > 2"),
			checkExpr: dexpr.MustNew("got == 1.0"),
		},
		{records: records,
			ruleExpr:  dexpr.MustNew("band == 9"),
			checkExpr: dexpr.MustNew("got >= 0 && got <= 1.0"),
		},
		{records: records,
			ruleExpr:  dexpr.MustNew("band != 9"),
			checkExpr: dexpr.MustNew("got >= -1.0 && got <= 0"),
		},
		{records: records,
			ruleExpr:  dexpr.MustNew("cost <= 2"),
			checkExpr: dexpr.MustNew("got == -1.0"),
		},
		{records: []map[string]*dlit.Literal{},
			ruleExpr:  dexpr.MustNew("1==1"),
			checkExpr: dexpr.MustNew("got == 0"),
		},
	}
	callFuncs := map[string]dexpr.CallFun{}
	for _, c := range cases {
		mccCostGt2, err := New("mccCostGt2", "mcc", "cost > 2")
		if err != nil {
			t.Fatalf("New(\"mccCostGt2\", \"mcc\", \"cost > 2\") err: %s", err)
		}
		aggregators := []Aggregator{mccCostGt2}

		for _, record := range c.records {
			isTrue, err := c.ruleExpr.EvalBool(record, callFuncs)
			if err != nil {
				t.Fatalf("EvalBool(%v, callFuncs) err: %v", record, err)
			}
			mccCostGt2.NextRecord(record, isTrue)
		}
		numRecords := int64(len(c.records))
		got := mccCostGt2.GetResult(aggregators, goals, numRecords)
		vars := map[string]*dlit.Literal{"got": got}
		isCorrect, err := c.checkExpr.EvalBool(vars, callFuncs)
		if err != nil {
			t.Fatalf("EvalBool(%v, callFuncs) err: %v", vars, err)
		}
		if !isCorrect {
			t.Errorf("EvalBool(%v, callFuncs) got: %v, want: %v",
				vars, got, c.checkExpr)
		}
	}
}

/*
func TestMCCCloneNew(t *testing.T) {
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(3),
		"band":   dlit.MustNew(4),
		"cost":   dlit.MustNew(4),
	}
	goals := []*goal.Goal{}
	numRecords := int64(1)
	mccCostGt2, err := New("mccCostGt2", "mcc", "cost > 2")
	if err != nil {
		t.Fatalf("New(\"mccCostGt2\", \"mcc\", \"cost > 2\") err: %s", err)
	}
	mccCostGt2_2 := mccCostGt2.CloneNew()
	aggregators := []Aggregator{}
	want := int64(100)
	mccCostGt2.NextRecord(record, true)
	got1 := mccCostGt2.GetResult(aggregators, goals, numRecords)
	got2 := mccCostGt2_2.GetResult(aggregators, goals, numRecords)

	gotInt1, gotIsInt1 := got1.Int()
	if !gotIsInt1 || gotInt1 != want {
		t.Errorf("GetResult() got: %d, want: %d", gotInt1, want)
	}
	gotInt2, gotIsInt2 := got2.Int()
	if !gotIsInt2 || gotInt2 != 0 {
		t.Errorf("GetResult() got: %d, want: %d", gotInt1, 0)
	}
}
*/
