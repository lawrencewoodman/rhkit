package rule

import (
	"github.com/lawrencewoodman/dlit"
	"reflect"
	"sort"
	"testing"
)

func TestNewOr(t *testing.T) {
	cases := []struct {
		ruleA Rule
		ruleB Rule
	}{
		{ruleA: NewLEFV("flow", dlit.MustNew(1.05)),
			ruleB: NewGEFV("flow", dlit.MustNew(1.06))},
		{ruleA: NewLEFV("flow", dlit.MustNew(1.05)),
			ruleB: NewEQFF("flow", "rate")},
		{ruleA: NewEQFF("flow", "rate"),
			ruleB: NewLEFV("flow", dlit.MustNew(1.05))},
		{ruleA: NewLEFV("rate", dlit.MustNew(1.05)),
			ruleB: NewLEFV("flow", dlit.MustNew(2.70))},
		{ruleA: NewLEFV("rate", dlit.MustNew(1.05)),
			ruleB: NewGEFV("flow", dlit.MustNew(1.05))},
		{ruleA: NewLEFV("flow", dlit.MustNew(1.05)),
			ruleB: NewGEFV("flow", dlit.MustNew(2.1))},
		{ruleA: NewLEFV("flow", dlit.MustNew(1.05)),
			ruleB: NewGEFV("rate", dlit.MustNew(2.1))},
		{ruleA: NewGEFV("flow", dlit.MustNew(2.1)),
			ruleB: NewLEFV("flow", dlit.MustNew(1.05))},
		{ruleA: NewGEFV("flow", dlit.MustNew(1.05)),
			ruleB: NewGEFV("rate", dlit.MustNew(2.07))},
		{ruleA: NewGEFV("flow", dlit.MustNew(1.05)),
			ruleB: MustNewBetweenFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1))},
		{ruleA: MustNewBetweenFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1)),
			ruleB: NewGEFV("flow", dlit.MustNew(1.05))},
		{ruleA: NewGEFV("rate", dlit.MustNew(1)),
			ruleB: MustNewBetweenFV("rate", dlit.MustNew(0), dlit.MustNew(2))},
		{ruleA: MustNewBetweenFV("rate", dlit.MustNew(0), dlit.MustNew(2)),
			ruleB: NewGEFV("rate", dlit.MustNew(1))},
		{ruleA: MustNewBetweenFV("rate", dlit.MustNew(1), dlit.MustNew(5)),
			ruleB: MustNewBetweenFV("rate", dlit.MustNew(7), dlit.MustNew(8)),
		},
		{ruleA: MustNewBetweenFV("flow", dlit.MustNew(1.7), dlit.MustNew(5.2)),
			ruleB: MustNewBetweenFV("flow", dlit.MustNew(7.3), dlit.MustNew(8.9)),
		},

		{ruleA: NewInFV("group",
			[]*dlit.Literal{
				dlit.NewString("collingwood"), dlit.NewString("drake"),
			},
		),
			ruleB: NewInFV("group",
				[]*dlit.Literal{
					dlit.NewString("nelson"), dlit.NewString("mountbatten"),
				},
			),
		},
		{ruleA: NewGEFV("flowA", dlit.MustNew(1.05)),
			ruleB: MustNewOutsideFV("flowB", dlit.MustNew(0.7), dlit.MustNew(2.1)),
		},
		{ruleA: MustNewOutsideFV("flowA", dlit.MustNew(0.7), dlit.MustNew(2.1)),
			ruleB: NewGEFV("flowB", dlit.MustNew(1.05)),
		},
		{ruleA: NewGEFV("rateA", dlit.MustNew(1)),
			ruleB: MustNewOutsideFV("rateB", dlit.MustNew(0), dlit.MustNew(2)),
		},
		{ruleA: MustNewOutsideFV("rateA", dlit.MustNew(0), dlit.MustNew(2)),
			ruleB: NewGEFV("rateB", dlit.MustNew(1)),
		},
		{ruleA: MustNewOutsideFV("rateA", dlit.MustNew(1), dlit.MustNew(5)),
			ruleB: MustNewOutsideFV("rateB", dlit.MustNew(7), dlit.MustNew(8)),
		},
		{ruleA: MustNewOutsideFV("flowA", dlit.MustNew(1.7), dlit.MustNew(5.2)),
			ruleB: MustNewOutsideFV("flowB", dlit.MustNew(7.3), dlit.MustNew(8.9)),
		},
		{ruleA: MustNewBetweenFV("flow", dlit.MustNew(0.7), dlit.MustNew(22.1)),
			ruleB: MustNewOutsideFV("flow", dlit.MustNew(1.05), dlit.MustNew(17.5)),
		},
		{ruleA: MustNewOutsideFV("flow", dlit.MustNew(1.05), dlit.MustNew(17.5)),
			ruleB: MustNewBetweenFV("flow", dlit.MustNew(0.7), dlit.MustNew(22.1)),
		},
	}
	for _, c := range cases {
		r, err := NewOr(c.ruleA, c.ruleB)
		if r == nil {
			t.Errorf("NewOr(%s, %s) rule got: nil, want: !nil", c.ruleA, c.ruleB)
		}
		if err != nil {
			t.Errorf("NewOr(%s, %s) got err: %s", c.ruleA, c.ruleB, err)
		}
	}
}

func TestNewOr_errors(t *testing.T) {
	cases := []struct {
		ruleA      Rule
		ruleB      Rule
		wantErrStr string
	}{
		{ruleA: NewLEFV("flow", dlit.MustNew(1.05)),
			ruleB:      NewTrue(),
			wantErrStr: "can't Or rule: flow <= 1.05, with: true()",
		},
		{ruleA: NewTrue(),
			ruleB:      NewLEFV("flow", dlit.MustNew(1.05)),
			wantErrStr: "can't Or rule: true(), with: flow <= 1.05",
		},
		{ruleA: NewLEFV("flow", dlit.MustNew(1.05)),
			ruleB:      NewLEFV("flow", dlit.MustNew(2.07)),
			wantErrStr: "can't Or rule: flow <= 1.05, with: flow <= 2.07",
		},
		{ruleA: NewGEFV("flow", dlit.MustNew(1.05)),
			ruleB:      NewGEFV("flow", dlit.MustNew(2.07)),
			wantErrStr: "can't Or rule: flow >= 1.05, with: flow >= 2.07",
		},
		{ruleA: NewGEFV("flow", dlit.MustNew(1.05)),
			ruleB:      MustNewOutsideFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1)),
			wantErrStr: "can't Or rule: flow >= 1.05, with: flow <= 0.7 || flow >= 2.1",
		},
		{ruleA: MustNewOutsideFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1)),
			ruleB:      NewGEFV("flow", dlit.MustNew(1.05)),
			wantErrStr: "can't Or rule: flow <= 0.7 || flow >= 2.1, with: flow >= 1.05",
		},
		{ruleA: MustNewOutsideFV("flow", dlit.MustNew(1.7), dlit.MustNew(5.2)),
			ruleB:      MustNewOutsideFV("flow", dlit.MustNew(7.3), dlit.MustNew(8.9)),
			wantErrStr: "can't Or rule: flow <= 1.7 || flow >= 5.2, with: flow <= 7.3 || flow >= 8.9",
		},
	}
	for _, c := range cases {
		r, err := NewOr(c.ruleA, c.ruleB)
		if r != nil {
			t.Errorf("NewOr(%s, %s) rule got: %s, want: nil", c.ruleA, c.ruleB, r)
		}
		if err == nil {
			t.Errorf("NewOr(%s, %s) got err: nil, want: %s",
				c.ruleA, c.ruleB, c.wantErrStr)
		} else if err.Error() != c.wantErrStr {
			t.Errorf("NewOr(%s, %s) got err: %s, want: %s",
				c.ruleA, c.ruleB, err, c.wantErrStr)
		}
	}
}

func TestMustNewOr(t *testing.T) {
	ruleA := NewEQFF("flow", "rate")
	ruleB := NewEQFF("income", "cost")
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustNewOr(%s, %s) panic: %s", ruleA, ruleB, r)
		}
	}()
	MustNewOr(ruleA, ruleB)
}

func TestMustNewOr_panic(t *testing.T) {
	ruleA := NewTrue()
	ruleB := NewEQFF("income", "cost")
	paniced := false
	wantPanic := "can't Or rule: true(), with: income == cost"
	defer func() {
		if r := recover(); r != nil {
			if r.(error).Error() == wantPanic {
				paniced = true
			} else {
				t.Errorf("MustNewOr(%s, %s) - got panic: %s, want: %s",
					ruleA, ruleB, r, wantPanic)
			}
		}
	}()
	MustNewOr(ruleA, ruleB)
	if !paniced {
		t.Errorf("MustNewOr(%s, %s) - failed to panic with: %s",
			ruleA, ruleB, wantPanic)
	}
}

func TestOrString(t *testing.T) {
	cases := []struct {
		ruleA Rule
		ruleB Rule
		want  string
	}{
		{ruleA: NewEQFF("flow", "flow"),
			ruleB: NewEQFF("income", "cost"),
			want:  "flow == flow || income == cost",
		},
		{ruleA: MustNewAnd(NewEQFF("flow", "flow"), NewEQFF("flow", "flow")),
			ruleB: NewEQFF("income", "cost"),
			want:  "(flow == flow && flow == flow) || income == cost",
		},
		{ruleA: MustNewOr(NewEQFF("flow", "flow"), NewEQFF("flow", "flow")),
			ruleB: NewEQFF("income", "cost"),
			want:  "(flow == flow || flow == flow) || income == cost",
		},
		{ruleA: NewEQFF("income", "cost"),
			ruleB: MustNewAnd(NewEQFF("flow", "flow"), NewEQFF("flow", "flow")),
			want:  "income == cost || (flow == flow && flow == flow)",
		},
		{ruleA: NewEQFF("income", "cost"),
			ruleB: MustNewOr(NewEQFF("flow", "flow"), NewEQFF("flow", "flow")),
			want:  "income == cost || (flow == flow || flow == flow)",
		},
		{ruleA: MustNewAnd(NewEQFV("income", dlit.MustNew(5)),
			NewEQFF("flow", "flow")),
			ruleB: MustNewAnd(NewEQFV("cost", dlit.MustNew(6)),
				NewEQFF("flow", "flow")),
			want: "(income == 5 && flow == flow) || (cost == 6 && flow == flow)",
		},
		{ruleA: MustNewOr(NewEQFV("income", dlit.MustNew(5)),
			NewEQFF("flow", "flow")),
			ruleB: MustNewOr(NewEQFV("cost", dlit.MustNew(6)),
				NewEQFF("flow", "flow")),
			want: "(income == 5 || flow == flow) || (cost == 6 || flow == flow)",
		},
		{ruleA: MustNewBetweenFV("flow", dlit.MustNew(5.24), dlit.MustNew(10.89)),
			ruleB: NewEQFF("income", "cost"),
			want:  "(flow >= 5.24 && flow <= 10.89) || income == cost",
		},
		{ruleA: NewEQFF("income", "cost"),
			ruleB: MustNewBetweenFV("flow", dlit.MustNew(5.24), dlit.MustNew(10.89)),
			want:  "income == cost || (flow >= 5.24 && flow <= 10.89)",
		},
		{ruleA: MustNewOutsideFV("flow", dlit.MustNew(5), dlit.MustNew(10)),
			ruleB: NewEQFF("income", "cost"),
			want:  "(flow <= 5 || flow >= 10) || income == cost",
		},
		{ruleA: NewEQFF("income", "cost"),
			ruleB: MustNewOutsideFV("flow", dlit.MustNew(5), dlit.MustNew(10)),
			want:  "income == cost || (flow <= 5 || flow >= 10)",
		},
		{ruleA: NewInFV("group",
			[]*dlit.Literal{
				dlit.NewString("collingwood"), dlit.NewString("drake"),
				dlit.NewString("nelson"),
			},
		),
			ruleB: NewInFV("group",
				[]*dlit.Literal{
					dlit.NewString("nelson"), dlit.NewString("mountbatten"),
				},
			),
			want: "in(group,\"collingwood\",\"drake\",\"nelson\",\"mountbatten\")",
		},
		{ruleA: NewInFV("team",
			[]*dlit.Literal{
				dlit.NewString("collingwood"), dlit.NewString("drake"),
				dlit.NewString("nelson"),
			},
		),
			ruleB: NewInFV("group",
				[]*dlit.Literal{
					dlit.NewString("nelson"), dlit.NewString("mountbatten"),
				},
			),
			want: "in(team,\"collingwood\",\"drake\",\"nelson\") || in(group,\"nelson\",\"mountbatten\")",
		},
	}
	for _, c := range cases {
		r := MustNewOr(c.ruleA, c.ruleB)
		got := r.String()
		if got != c.want {
			t.Errorf("String() got: %s, want: %s", got, c.want)
		}
	}
}

func TestOrIsTrue(t *testing.T) {
	cases := []struct {
		ruleA Rule
		ruleB Rule
		want  bool
	}{
		{ruleA: NewEQFF("income", "income"),
			ruleB: NewEQFF("cost", "cost"),
			want:  true,
		},
		{ruleA: NewNEFF("income", "income"),
			ruleB: NewEQFF("cost", "cost"),
			want:  true,
		},
		{ruleA: NewEQFF("cost", "cost"),
			ruleB: NewNEFF("income", "income"),
			want:  true,
		},
		{ruleA: NewNEFF("income", "income"),
			ruleB: NewNEFF("income", "income"),
			want:  false,
		},
		{ruleA: NewInFV("group",
			[]*dlit.Literal{
				dlit.NewString("collingwood"), dlit.NewString("drake"),
				dlit.NewString("nelson"),
			},
		),
			ruleB: NewInFV("group",
				[]*dlit.Literal{
					dlit.NewString("nelson"), dlit.NewString("mountbatten"),
				},
			),
			want: true,
		},
		{ruleA: NewInFV("group",
			[]*dlit.Literal{
				dlit.NewString("mountbatten"), dlit.NewString("drake"),
				dlit.NewString("nelson"),
			},
		),
			ruleB: NewInFV("group",
				[]*dlit.Literal{
					dlit.NewString("nelson"), dlit.NewString("collingwood"),
				},
			),
			want: true,
		},
		{ruleA: NewInFV("team",
			[]*dlit.Literal{
				dlit.NewString("first"), dlit.NewString("second"),
				dlit.NewString("third"),
			},
		),
			ruleB: NewInFV("group",
				[]*dlit.Literal{
					dlit.NewString("nelson"), dlit.NewString("mountbatten"),
				},
			),
			want: true,
		},
		{ruleA: NewInFV("team",
			[]*dlit.Literal{
				dlit.NewString("first"), dlit.NewString("second"),
				dlit.NewString("third"),
			},
		),
			ruleB: NewInFV("group",
				[]*dlit.Literal{
					dlit.NewString("collingwood"), dlit.NewString("nelson"),
				},
			),
			want: false,
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"cost":   dlit.MustNew(20),
		"group":  dlit.NewString("mountbatten"),
		"team":   dlit.NewString("ace"),
	}
	for _, c := range cases {
		r := MustNewOr(c.ruleA, c.ruleB)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestOrIsTrue_errors(t *testing.T) {
	cases := []struct {
		ruleA   Rule
		ruleB   Rule
		wantErr error
	}{
		{ruleA: NewEQFF("flow", "flow"),
			ruleB: NewEQFF("fred", "income"),
			wantErr: InvalidRuleError{
				Rule: MustNewOr(NewEQFF("flow", "flow"), NewEQFF("fred", "income")),
			},
		},
		{ruleA: NewEQFF("fred", "income"),
			ruleB: NewEQFF("flow", "flow"),
			wantErr: InvalidRuleError{
				Rule: MustNewOr(NewEQFF("fred", "income"), NewEQFF("flow", "flow")),
			},
		},
		{ruleA: NewEQFF("fred", "income"),
			ruleB: NewEQFF("bob", "cost"),
			wantErr: InvalidRuleError{
				Rule: MustNewOr(NewEQFF("fred", "income"), NewEQFF("bob", "cost")),
			},
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"cost":   dlit.MustNew(18),
		"band":   dlit.NewString("alpha"),
	}
	for _, c := range cases {
		r := MustNewOr(c.ruleA, c.ruleB)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}

func TestOrGetFields(t *testing.T) {
	cases := []struct {
		ruleA Rule
		ruleB Rule
		want  []string
	}{
		{ruleA: NewEQFF("flow", "flow"),
			ruleB: NewEQFF("income", "cost"),
			want:  []string{"flow", "income", "cost"},
		},
		{ruleA: MustNewAnd(NewEQFF("flowIn", "flowOut"), NewEQFF("rate", "flow")),
			ruleB: NewEQFF("income", "cost"),
			want:  []string{"flow", "flowIn", "flowOut", "rate", "income", "cost"},
		},
		{
			ruleA: NewEQFF("income", "cost"),
			ruleB: MustNewAnd(NewEQFF("flowIn", "flowOut"), NewEQFF("rate", "flow")),
			want:  []string{"flow", "flowIn", "flowOut", "rate", "income", "cost"},
		},
		{ruleA: NewEQFF("income", "cost"),
			ruleB: MustNewOr(NewEQFF("flowIn", "flowOut"), NewEQFF("rate", "flow")),
			want:  []string{"flow", "flowIn", "flowOut", "rate", "income", "cost"},
		},
	}
	for _, c := range cases {
		r := MustNewOr(c.ruleA, c.ruleB)
		got := r.GetFields()
		sort.Strings(got)
		sort.Strings(c.want)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("GetFields() got: %s, want: %s", got, c.want)
		}
	}
}
