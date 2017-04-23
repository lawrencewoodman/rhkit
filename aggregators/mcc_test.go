package aggregators

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/goal"
	"github.com/vlifesystems/rhkit/internal/dexprfuncs"
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
			ruleExpr:  dexpr.MustNew("cost > 2", dexprfuncs.CallFuncs),
			checkExpr: dexpr.MustNew("got == 1.0", dexprfuncs.CallFuncs),
		},
		{records: records,
			ruleExpr:  dexpr.MustNew("band == 9", dexprfuncs.CallFuncs),
			checkExpr: dexpr.MustNew("got >= 0 && got <= 1.0", dexprfuncs.CallFuncs),
		},
		{records: records,
			ruleExpr:  dexpr.MustNew("band != 9", dexprfuncs.CallFuncs),
			checkExpr: dexpr.MustNew("got >= -1.0 && got <= 0", dexprfuncs.CallFuncs),
		},
		{records: records,
			ruleExpr:  dexpr.MustNew("cost <= 2", dexprfuncs.CallFuncs),
			checkExpr: dexpr.MustNew("got == -1.0", dexprfuncs.CallFuncs),
		},
		{records: records,
			ruleExpr:  dexpr.MustNew("1 != 1", dexprfuncs.CallFuncs),
			checkExpr: dexpr.MustNew("got == 0", dexprfuncs.CallFuncs),
		},
		{records: []map[string]*dlit.Literal{},
			ruleExpr:  dexpr.MustNew("1 == 1", dexprfuncs.CallFuncs),
			checkExpr: dexpr.MustNew("got == 0", dexprfuncs.CallFuncs),
		},
	}
	for _, c := range cases {
		mccCostGt2Desc := MustNew("mccCostGt2", "mcc", "cost > 2")
		mccCostGt2 := mccCostGt2Desc.New()
		instances := []AggregatorInstance{mccCostGt2}

		for _, record := range c.records {
			isTrue, err := c.ruleExpr.EvalBool(record)
			if err != nil {
				t.Fatalf("EvalBool(%v, callFuncs) err: %v", record, err)
			}
			mccCostGt2.NextRecord(record, isTrue)
		}
		numRecords := int64(len(c.records))
		got := mccCostGt2.GetResult(instances, goals, numRecords)
		vars := map[string]*dlit.Literal{"got": got}
		isCorrect, err := c.checkExpr.EvalBool(vars)
		if err != nil {
			t.Fatalf("EvalBool(%v, callFuncs) err: %v", vars, err)
		}
		if !isCorrect {
			t.Errorf("GetResult() (c.ruleExpr: %s) got: %v, want: %v",
				c.ruleExpr, got, c.checkExpr)
		}
	}
}

func TestMCCSpecGetName(t *testing.T) {
	name := "a"
	as := MustNew(name, "mcc", "band > 4")
	got := as.GetName()
	if got != name {
		t.Errorf("GetName - got: %s, want: %s", got, name)
	}
}

func TestMCCSpecGetKind(t *testing.T) {
	kind := "mcc"
	as := MustNew("a", kind, "band > 4")
	got := as.GetKind()
	if got != kind {
		t.Errorf("GetKind - got: %s, want: %s", got, kind)
	}
}

func TestMCCSpecGetArg(t *testing.T) {
	arg := "band > 4"
	as := MustNew("a", "mcc", arg)
	got := as.GetArg()
	if got != arg {
		t.Errorf("GetArg - got: %s, want: %s", got, arg)
	}
}
