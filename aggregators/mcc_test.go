package aggregators

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/goal"
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
		{records: records,
			ruleExpr:  dexpr.MustNew("1 != 1"),
			checkExpr: dexpr.MustNew("got == 0"),
		},
		{records: []map[string]*dlit.Literal{},
			ruleExpr:  dexpr.MustNew("1 == 1"),
			checkExpr: dexpr.MustNew("got == 0"),
		},
	}
	callFuncs := map[string]dexpr.CallFun{}
	for _, c := range cases {
		mccCostGt2Desc, err := New("mccCostGt2", "mcc", "cost > 2")
		if err != nil {
			t.Fatalf("New(\"mccCostGt2\", \"mcc\", \"cost > 2\") err: %s", err)
		}
		mccCostGt2 := mccCostGt2Desc.New()
		instances := []AggregatorInstance{mccCostGt2}

		for _, record := range c.records {
			isTrue, err := c.ruleExpr.EvalBool(record, callFuncs)
			if err != nil {
				t.Fatalf("EvalBool(%v, callFuncs) err: %v", record, err)
			}
			mccCostGt2.NextRecord(record, isTrue)
		}
		numRecords := int64(len(c.records))
		got := mccCostGt2.GetResult(instances, goals, numRecords)
		vars := map[string]*dlit.Literal{"got": got}
		isCorrect, err := c.checkExpr.EvalBool(vars, callFuncs)
		if err != nil {
			t.Fatalf("EvalBool(%v, callFuncs) err: %v", vars, err)
		}
		if !isCorrect {
			t.Errorf("GetResult() (c.ruleExpr: %s) got: %v, want: %v",
				c.ruleExpr, got, c.checkExpr)
		}
	}
}
