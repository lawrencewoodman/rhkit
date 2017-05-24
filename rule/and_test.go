package rule

import (
	"github.com/lawrencewoodman/dlit"
	"reflect"
	"sort"
	"testing"
)

func TestNewAnd(t *testing.T) {
	cases := []struct {
		ruleA Rule
		ruleB Rule
	}{
		{ruleA: NewLEFV("flow", dlit.MustNew(1.05)),
			ruleB: NewGEFV("flow", dlit.MustNew(0.0))},
		{ruleA: NewLEFV("flow", dlit.MustNew(1.05)),
			ruleB: NewEQFF("flow", "rate")},
		{ruleA: NewEQFF("flow", "rate"),
			ruleB: NewLEFV("flow", dlit.MustNew(1.05))},
		{ruleA: NewLEFV("rate", dlit.MustNew(1.05)),
			ruleB: NewLEFV("flow", dlit.MustNew(2.70))},
		{ruleA: NewLEFV("rate", dlit.MustNew(1.05)),
			ruleB: NewGEFV("flow", dlit.MustNew(1.05))},
		{ruleA: NewGEFV("flow", dlit.MustNew(1.05)),
			ruleB: NewLEFV("flow", dlit.MustNew(2.1))},
		{ruleA: NewGEFV("rate", dlit.MustNew(1.05)),
			ruleB: MustNewBetweenFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1))},
		{ruleA: MustNewBetweenFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1)),
			ruleB: NewGEFV("rate", dlit.MustNew(1.05))},
		{ruleA: MustNewOutsideFV("flow", dlit.MustNew(1), dlit.MustNew(4)),
			ruleB: NewGEFV("flow", dlit.MustNew(0))},
		{ruleA: NewGEFV("flow", dlit.MustNew(0)),
			ruleB: MustNewOutsideFV("flow", dlit.MustNew(1), dlit.MustNew(4))},
		{ruleA: MustNewOutsideFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1)),
			ruleB: NewGEFV("flow", dlit.MustNew(0.2))},
		{ruleA: NewGEFV("flow", dlit.MustNew(0.2)),
			ruleB: MustNewOutsideFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1))},
		{ruleA: MustNewOutsideFV("rate", dlit.MustNew(0), dlit.MustNew(2)),
			ruleB: NewGEFV("rate", dlit.MustNew(-5)),
		},
		{ruleA: NewGEFV("rate", dlit.MustNew(-5)),
			ruleB: MustNewOutsideFV("rate", dlit.MustNew(0), dlit.MustNew(2)),
		},
		{ruleA: MustNewOutsideFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1)),
			ruleB: NewGEFV("flow", dlit.MustNew(0.5)),
		},
		{ruleA: NewGEFV("flow", dlit.MustNew(0.5)),
			ruleB: MustNewOutsideFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1)),
		},
		{ruleA: MustNewOutsideFV("rate", dlit.MustNew(0), dlit.MustNew(2)),
			ruleB: NewLEFV("rate", dlit.MustNew(4)),
		},
		{ruleA: NewLEFV("rate", dlit.MustNew(4)),
			ruleB: MustNewOutsideFV("rate", dlit.MustNew(0), dlit.MustNew(2)),
		},
		{ruleA: MustNewOutsideFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1)),
			ruleB: NewLEFV("flow", dlit.MustNew(5.5)),
		},
		{ruleA: NewLEFV("flow", dlit.MustNew(5.7)),
			ruleB: MustNewOutsideFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1)),
		},
		{ruleA: MustNewBetweenFV("flow", dlit.MustNew(0.7), dlit.MustNew(22.1)),
			ruleB: MustNewOutsideFV("flow", dlit.MustNew(1.05), dlit.MustNew(17.5)),
		},
		{ruleA: MustNewOutsideFV("flow", dlit.MustNew(1.05), dlit.MustNew(17.5)),
			ruleB: MustNewBetweenFV("flow", dlit.MustNew(0.7), dlit.MustNew(22.1)),
		},
		{ruleA: MustNewBetweenFV("rate", dlit.MustNew(0), dlit.MustNew(22)),
			ruleB: MustNewOutsideFV("rate", dlit.MustNew(1), dlit.MustNew(17)),
		},
		{ruleA: MustNewOutsideFV("rate", dlit.MustNew(1), dlit.MustNew(17)),
			ruleB: MustNewBetweenFV("rate", dlit.MustNew(0), dlit.MustNew(22)),
		},
		{ruleA: NewInFV("group", []*dlit.Literal{
			dlit.NewString("bob"),
			dlit.NewString("fred"),
			dlit.NewString("albert"),
		}),
			ruleB: NewEQFF("team", "group"),
		},
		{ruleA: NewEQFF("team", "group"),
			ruleB: NewInFV("group", []*dlit.Literal{
				dlit.NewString("bob"),
				dlit.NewString("fred"),
				dlit.NewString("albert"),
			}),
		},
		{ruleA: NewInFV("group", []*dlit.Literal{
			dlit.NewString("bob"),
			dlit.NewString("fred"),
			dlit.NewString("albert"),
		}),
			ruleB: NewEQFF("group", "team"),
		},
		{ruleA: NewEQFF("group", "team"),
			ruleB: NewInFV("group", []*dlit.Literal{
				dlit.NewString("bob"),
				dlit.NewString("fred"),
				dlit.NewString("albert"),
			}),
		},
		{ruleA: MustNewAnd(NewEQFF("team", "group"), NewEQFF("flow", "rate")),
			ruleB: NewInFV("group", []*dlit.Literal{
				dlit.NewString("bob"),
				dlit.NewString("fred"),
				dlit.NewString("albert"),
			}),
		},
		{ruleA: NewInFV("group", []*dlit.Literal{
			dlit.NewString("bob"),
			dlit.NewString("fred"),
			dlit.NewString("albert"),
		}),
			ruleB: MustNewAnd(NewEQFF("team", "group"), NewEQFF("flow", "rate")),
		},
		{ruleA: NewEQFV("group", dlit.MustNew("bob")),
			ruleB: NewEQFV("team", dlit.MustNew("ruth"))},
	}
	for _, c := range cases {
		r, err := NewAnd(c.ruleA, c.ruleB)
		if err != nil {
			t.Errorf("NewAnd(%s, %s) got err: %s", c.ruleA, c.ruleB, err)
			continue
		}
		if r == nil {
			t.Errorf("NewAnd(%s, %s) rule got: nil, want: !nil", c.ruleA, c.ruleB)
		}
	}
}

func TestNewAnd_errors(t *testing.T) {
	cases := []struct {
		ruleA      Rule
		ruleB      Rule
		wantErrStr string
	}{
		{ruleA: NewLEFV("flow", dlit.MustNew(1.05)),
			ruleB:      NewLEFV("flow", dlit.MustNew(2.70)),
			wantErrStr: "can't And rule: flow <= 1.05, with: flow <= 2.7",
		},
		{ruleA: NewLEFV("flow", dlit.MustNew(1.05)),
			ruleB:      NewLEFV("flow", dlit.MustNew(1.05)),
			wantErrStr: "can't And rule: flow <= 1.05, with: flow <= 1.05",
		},
		{ruleA: NewGEFV("flow", dlit.MustNew(1.05)),
			ruleB:      NewGEFV("flow", dlit.MustNew(1.05)),
			wantErrStr: "can't And rule: flow >= 1.05, with: flow >= 1.05",
		},
		{ruleA: NewLEFV("flow", dlit.MustNew(1.05)),
			ruleB:      NewGEFV("flow", dlit.MustNew(1.05)),
			wantErrStr: "can't And rule: flow <= 1.05, with: flow >= 1.05",
		},
		{ruleA: NewLEFV("flow", dlit.MustNew(1.05)),
			ruleB:      NewGEFV("flow", dlit.MustNew(2.1)),
			wantErrStr: "can't And rule: flow <= 1.05, with: flow >= 2.1",
		},
		{ruleA: NewGEFV("flow", dlit.MustNew(2.1)),
			ruleB:      NewLEFV("flow", dlit.MustNew(1.05)),
			wantErrStr: "can't And rule: flow >= 2.1, with: flow <= 1.05",
		},
		{ruleA: NewGEFV("flow", dlit.MustNew(2.1)),
			ruleB:      NewLEFV("flow", dlit.MustNew(2.1)),
			wantErrStr: "can't And rule: flow >= 2.1, with: flow <= 2.1",
		},
		{ruleA: NewInFV("group", []*dlit.Literal{
			dlit.NewString("bob"),
			dlit.NewString("fred"),
			dlit.NewString("albert"),
		}),
			ruleB:      NewEQFV("group", dlit.MustNew("norris")),
			wantErrStr: "can't And rule: in(group,\"bob\",\"fred\",\"albert\"), with: group == \"norris\"",
		},
		{ruleA: NewEQFV("group", dlit.MustNew("norris")),
			ruleB: NewInFV("group", []*dlit.Literal{
				dlit.NewString("bob"),
				dlit.NewString("fred"),
				dlit.NewString("albert"),
			}),
			wantErrStr: "can't And rule: group == \"norris\", with: in(group,\"bob\",\"fred\",\"albert\")",
		},
		{ruleA: NewInFV("group", []*dlit.Literal{
			dlit.NewString("bob"),
			dlit.NewString("fred"),
			dlit.NewString("albert"),
		}),
			ruleB: NewInFV("group", []*dlit.Literal{
				dlit.NewString("harry"),
				dlit.NewString("fred"),
				dlit.NewString("albert"),
			}),
			wantErrStr: "can't And rule: in(group,\"bob\",\"fred\",\"albert\"), with: in(group,\"harry\",\"fred\",\"albert\")",
		},
		{ruleA: NewTrue(),
			ruleB:      NewEQFF("flow", "rate"),
			wantErrStr: "can't And rule: true(), with: flow == rate",
		},
		{ruleA: NewEQFF("flow", "rate"),
			ruleB:      NewTrue(),
			wantErrStr: "can't And rule: flow == rate, with: true()",
		},
		{ruleA: NewEQFV("flow", dlit.MustNew(1.05)),
			ruleB:      NewEQFV("flow", dlit.MustNew(2.1)),
			wantErrStr: "can't And rule: flow == 1.05, with: flow == 2.1",
		},
		{ruleA: NewEQFV("flow", dlit.MustNew(1.05)),
			ruleB:      NewNEFV("flow", dlit.MustNew(2.1)),
			wantErrStr: "can't And rule: flow == 1.05, with: flow != 2.1",
		},
		{ruleA: NewEQFV("group", dlit.MustNew(1)),
			ruleB:      NewEQFV("group", dlit.MustNew(2)),
			wantErrStr: "can't And rule: group == 1, with: group == 2",
		},
		{ruleA: NewEQFV("team", dlit.MustNew("oren")),
			ruleB:      NewEQFV("team", dlit.MustNew("melyn")),
			wantErrStr: "can't And rule: team == \"oren\", with: team == \"melyn\"",
		},
		{ruleA: NewEQFV("team", dlit.MustNew("oren")),
			ruleB:      NewNEFV("team", dlit.MustNew("melyn")),
			wantErrStr: "can't And rule: team == \"oren\", with: team != \"melyn\"",
		},
		{ruleA: NewNEFV("group", dlit.MustNew(1)),
			ruleB:      NewNEFV("group", dlit.MustNew(2)),
			wantErrStr: "can't And rule: group != 1, with: group != 2",
		},
		{ruleA: NewNEFV("group", dlit.MustNew(1)),
			ruleB:      NewEQFV("group", dlit.MustNew(2)),
			wantErrStr: "can't And rule: group != 1, with: group == 2",
		},
		{ruleA: NewNEFV("team", dlit.MustNew("oren")),
			ruleB:      NewNEFV("team", dlit.MustNew("melyn")),
			wantErrStr: "can't And rule: team != \"oren\", with: team != \"melyn\"",
		},
		{ruleA: NewNEFV("team", dlit.MustNew("oren")),
			ruleB:      NewEQFV("team", dlit.MustNew("melyn")),
			wantErrStr: "can't And rule: team != \"oren\", with: team == \"melyn\"",
		},
		{ruleA: MustNewBetweenFV("rate", dlit.MustNew(1), dlit.MustNew(5)),
			ruleB:      MustNewBetweenFV("rate", dlit.MustNew(7), dlit.MustNew(8)),
			wantErrStr: "can't And rule: rate >= 1 && rate <= 5, with: rate >= 7 && rate <= 8",
		},
		{ruleA: MustNewOutsideFV("rate", dlit.MustNew(1), dlit.MustNew(5)),
			ruleB:      MustNewOutsideFV("rate", dlit.MustNew(7), dlit.MustNew(8)),
			wantErrStr: "can't And rule: rate <= 1 || rate >= 5, with: rate <= 7 || rate >= 8",
		},
		{ruleA: MustNewOutsideFV("flow", dlit.MustNew(1.7), dlit.MustNew(5.2)),
			ruleB:      MustNewOutsideFV("flow", dlit.MustNew(7.3), dlit.MustNew(8.9)),
			wantErrStr: "can't And rule: flow <= 1.7 || flow >= 5.2, with: flow <= 7.3 || flow >= 8.9",
		},
		{ruleA: NewGEFV("flow", dlit.MustNew(1.05)),
			ruleB:      MustNewBetweenFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1)),
			wantErrStr: "can't And rule: flow >= 1.05, with: flow >= 0.7 && flow <= 2.1",
		},
		{ruleA: MustNewBetweenFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1)),
			ruleB:      NewGEFV("flow", dlit.MustNew(1.05)),
			wantErrStr: "can't And rule: flow >= 0.7 && flow <= 2.1, with: flow >= 1.05",
		},
		{ruleA: MustNewOutsideFV("rate", dlit.MustNew(0), dlit.MustNew(2)),
			ruleB:      NewGEFV("rate", dlit.MustNew(1)),
			wantErrStr: "can't And rule: rate <= 0 || rate >= 2, with: rate >= 1",
		},
		{ruleA: MustNewOutsideFV("rate", dlit.MustNew(0), dlit.MustNew(2)),
			ruleB:      NewGEFV("rate", dlit.MustNew(3)),
			wantErrStr: "can't And rule: rate <= 0 || rate >= 2, with: rate >= 3",
		},
		{ruleA: NewGEFV("rate", dlit.MustNew(1)),
			ruleB:      MustNewOutsideFV("rate", dlit.MustNew(0), dlit.MustNew(2)),
			wantErrStr: "can't And rule: rate >= 1, with: rate <= 0 || rate >= 2",
		},
		{ruleA: NewGEFV("rate", dlit.MustNew(3)),
			ruleB:      MustNewOutsideFV("rate", dlit.MustNew(0), dlit.MustNew(2)),
			wantErrStr: "can't And rule: rate >= 3, with: rate <= 0 || rate >= 2",
		},
		{ruleA: MustNewOutsideFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1)),
			ruleB:      NewGEFV("flow", dlit.MustNew(1.05)),
			wantErrStr: "can't And rule: flow <= 0.7 || flow >= 2.1, with: flow >= 1.05",
		},
		{ruleA: MustNewOutsideFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1)),
			ruleB:      NewGEFV("flow", dlit.MustNew(2.05)),
			wantErrStr: "can't And rule: flow <= 0.7 || flow >= 2.1, with: flow >= 2.05",
		},
		{ruleA: NewGEFV("flow", dlit.MustNew(1.05)),
			ruleB:      MustNewOutsideFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1)),
			wantErrStr: "can't And rule: flow >= 1.05, with: flow <= 0.7 || flow >= 2.1",
		},
		{ruleA: NewGEFV("flow", dlit.MustNew(2.05)),
			ruleB:      MustNewOutsideFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1)),
			wantErrStr: "can't And rule: flow >= 2.05, with: flow <= 0.7 || flow >= 2.1",
		},
		{ruleA: MustNewOutsideFV("rate", dlit.MustNew(0), dlit.MustNew(2)),
			ruleB:      NewLEFV("rate", dlit.MustNew(1)),
			wantErrStr: "can't And rule: rate <= 0 || rate >= 2, with: rate <= 1",
		},
		{ruleA: MustNewOutsideFV("rate", dlit.MustNew(0), dlit.MustNew(2)),
			ruleB:      NewLEFV("rate", dlit.MustNew(-1)),
			wantErrStr: "can't And rule: rate <= 0 || rate >= 2, with: rate <= -1",
		},
		{ruleA: NewLEFV("rate", dlit.MustNew(1)),
			ruleB:      MustNewOutsideFV("rate", dlit.MustNew(0), dlit.MustNew(2)),
			wantErrStr: "can't And rule: rate <= 1, with: rate <= 0 || rate >= 2",
		},
		{ruleA: NewLEFV("rate", dlit.MustNew(-1)),
			ruleB:      MustNewOutsideFV("rate", dlit.MustNew(0), dlit.MustNew(2)),
			wantErrStr: "can't And rule: rate <= -1, with: rate <= 0 || rate >= 2",
		},
		{ruleA: MustNewOutsideFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1)),
			ruleB:      NewLEFV("flow", dlit.MustNew(1.05)),
			wantErrStr: "can't And rule: flow <= 0.7 || flow >= 2.1, with: flow <= 1.05",
		},
		{ruleA: MustNewOutsideFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1)),
			ruleB:      NewLEFV("flow", dlit.MustNew(0.5)),
			wantErrStr: "can't And rule: flow <= 0.7 || flow >= 2.1, with: flow <= 0.5",
		},
		{ruleA: NewLEFV("flow", dlit.MustNew(1.05)),
			ruleB:      MustNewOutsideFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1)),
			wantErrStr: "can't And rule: flow <= 1.05, with: flow <= 0.7 || flow >= 2.1",
		},
		{ruleA: NewLEFV("flow", dlit.MustNew(0.5)),
			ruleB:      MustNewOutsideFV("flow", dlit.MustNew(0.7), dlit.MustNew(2.1)),
			wantErrStr: "can't And rule: flow <= 0.5, with: flow <= 0.7 || flow >= 2.1",
		},
	}
	for _, c := range cases {
		r, err := NewAnd(c.ruleA, c.ruleB)
		if r != nil || err == nil {
			t.Errorf("NewAnd(%s, %s) got rule: %s, want: nil, err: %v, want: %s",
				c.ruleA, c.ruleB, r, err, c.wantErrStr)
		} else if err.Error() != c.wantErrStr {
			t.Errorf("NewAnd(%s, %s) got err: %s, want: %s",
				c.ruleA, c.ruleB, err, c.wantErrStr)
		}
	}
}

func TestMustNewAnd(t *testing.T) {
	ruleA := NewEQFF("flow", "rate")
	ruleB := NewEQFF("income", "cost")
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustNewAnd(%s, %s) panic: %s", ruleA, ruleB, r)
		}
	}()
	MustNewAnd(ruleA, ruleB)
}

func TestMustNewAnd_panic(t *testing.T) {
	ruleA := NewTrue()
	ruleB := NewEQFF("income", "cost")
	paniced := false
	wantPanic := "can't And rule: true(), with: income == cost"
	defer func() {
		if r := recover(); r != nil {
			if r.(error).Error() == wantPanic {
				paniced = true
			} else {
				t.Errorf("MustNewAnd(%s, %s) - got panic: %s, want: %s",
					ruleA, ruleB, r, wantPanic)
			}
		}
	}()
	MustNewAnd(ruleA, ruleB)
	if !paniced {
		t.Errorf("MustNewAnd(%s, %s) - failed to panic with: %s",
			ruleA, ruleB, wantPanic)
	}
}

func TestAndString(t *testing.T) {
	cases := []struct {
		ruleA Rule
		ruleB Rule
		want  string
	}{
		{ruleA: NewEQFF("flow", "flow"),
			ruleB: NewEQFF("income", "cost"),
			want:  "flow == flow && income == cost",
		},
		{ruleA: MustNewAnd(NewEQFF("flow", "flow"), NewEQFF("flow", "flow")),
			ruleB: NewEQFF("income", "cost"),
			want:  "(flow == flow && flow == flow) && income == cost",
		},
		{ruleA: MustNewOr(NewEQFF("flow", "flow"), NewEQFF("flow", "flow")),
			ruleB: NewEQFF("income", "cost"),
			want:  "(flow == flow || flow == flow) && income == cost",
		},
		{ruleA: NewEQFF("income", "cost"),
			ruleB: MustNewAnd(NewEQFF("flow", "flow"), NewEQFF("flow", "flow")),
			want:  "income == cost && (flow == flow && flow == flow)",
		},
		{ruleA: NewEQFF("income", "cost"),
			ruleB: MustNewOr(NewEQFF("flow", "flow"), NewEQFF("flow", "flow")),
			want:  "income == cost && (flow == flow || flow == flow)",
		},
		{ruleA: MustNewAnd(NewEQFV("income", dlit.MustNew(5)),
			NewEQFF("flow", "flow")),
			ruleB: MustNewAnd(NewEQFV("cost", dlit.MustNew(6)),
				NewEQFF("flow", "flow")),
			want: "(income == 5 && flow == flow) && (cost == 6 && flow == flow)",
		},
		{ruleA: MustNewOr(NewEQFV("income", dlit.MustNew(5)),
			NewEQFF("flow", "flow")),
			ruleB: MustNewOr(NewEQFV("cost", dlit.MustNew(6)), NewEQFF("flow", "flow")),
			want:  "(income == 5 || flow == flow) && (cost == 6 || flow == flow)",
		},
		{ruleA: MustNewBetweenFV("flow", dlit.MustNew(5), dlit.MustNew(10)),
			ruleB: NewEQFF("income", "cost"),
			want:  "(flow >= 5 && flow <= 10) && income == cost",
		},
		{ruleA: NewEQFF("income", "cost"),
			ruleB: MustNewBetweenFV("flow", dlit.MustNew(5), dlit.MustNew(10)),
			want:  "income == cost && (flow >= 5 && flow <= 10)",
		},
		{ruleA: MustNewBetweenFV("flow", dlit.MustNew(5.24), dlit.MustNew(10.89)),
			ruleB: NewEQFF("income", "cost"),
			want:  "(flow >= 5.24 && flow <= 10.89) && income == cost",
		},
		{ruleA: NewEQFF("income", "cost"),
			ruleB: MustNewBetweenFV("flow", dlit.MustNew(5.24), dlit.MustNew(10.89)),
			want:  "income == cost && (flow >= 5.24 && flow <= 10.89)",
		},
		{ruleA: MustNewOutsideFV("flow", dlit.MustNew(5), dlit.MustNew(10)),
			ruleB: NewEQFF("income", "cost"),
			want:  "(flow <= 5 || flow >= 10) && income == cost",
		},
		{ruleA: NewEQFF("income", "cost"),
			ruleB: MustNewOutsideFV("flow", dlit.MustNew(5), dlit.MustNew(10)),
			want:  "income == cost && (flow <= 5 || flow >= 10)",
		},
		{ruleA: MustNewOutsideFV("flow", dlit.MustNew(5.24), dlit.MustNew(10.89)),
			ruleB: NewEQFF("income", "cost"),
			want:  "(flow <= 5.24 || flow >= 10.89) && income == cost",
		},
		{ruleA: NewEQFF("income", "cost"),
			ruleB: MustNewOutsideFV("flow", dlit.MustNew(5.24), dlit.MustNew(10.89)),
			want:  "income == cost && (flow <= 5.24 || flow >= 10.89)",
		},
		{ruleA: NewGEFV("income", dlit.MustNew(100)),
			ruleB: NewLEFV("income", dlit.MustNew(500)),
			want:  "income >= 100 && income <= 500",
		},
		{ruleA: NewLEFV("income", dlit.MustNew(500)),
			ruleB: NewGEFV("income", dlit.MustNew(100)),
			want:  "income >= 100 && income <= 500",
		},
		{ruleA: NewGEFV("rate", dlit.MustNew(100.78)),
			ruleB: NewLEFV("rate", dlit.MustNew(500.24)),
			want:  "rate >= 100.78 && rate <= 500.24",
		},
		{ruleA: NewLEFV("rate", dlit.MustNew(500.24)),
			ruleB: NewGEFV("rate", dlit.MustNew(100.78)),
			want:  "rate >= 100.78 && rate <= 500.24",
		},
	}
	for _, c := range cases {
		r := MustNewAnd(c.ruleA, c.ruleB)
		got := r.String()
		if got != c.want {
			t.Errorf("String() got: %s, want: %s", got, c.want)
		}
	}
}

func TestAndIsTrue(t *testing.T) {
	cases := []struct {
		ruleA Rule
		ruleB Rule
		want  bool
	}{
		{ruleA: NewEQFF("cost", "cost"),
			ruleB: NewEQFF("income", "income"),
			want:  true,
		},
		{ruleA: NewNEFF("income", "income"),
			ruleB: NewEQFF("cost", "cost"),
			want:  false,
		},
		{ruleA: NewEQFF("cost", "cost"),
			ruleB: NewNEFF("income", "income"),
			want:  false,
		},
		{ruleA: NewNEFF("cost", "cost"),
			ruleB: NewNEFF("income", "income"),
			want:  false,
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"cost":   dlit.MustNew(20),
	}
	for _, c := range cases {
		r := MustNewAnd(c.ruleA, c.ruleB)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestAndIsTrue_errors(t *testing.T) {
	cases := []struct {
		ruleA Rule
		ruleB Rule
	}{
		{ruleA: NewEQFF("cost", "cost"),
			ruleB: NewEQFF("fred", "income"),
		},
		{ruleA: NewEQFF("fred", "income"),
			ruleB: NewEQFF("cost", "cost"),
		},
		{ruleA: NewEQFF("fred", "income"),
			ruleB: NewEQFF("bob", "cost"),
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"cost":   dlit.MustNew(18),
		"band":   dlit.NewString("alpha"),
	}
	for _, c := range cases {
		r := MustNewAnd(c.ruleA, c.ruleB)
		wantErr := InvalidRuleError{Rule: r}
		_, err := r.IsTrue(record)
		if err != wantErr {
			t.Errorf("IsTrue(record) rule: %s, err: %v, want: %v", r, err, wantErr)
		}
	}
}

func TestAndGetFields(t *testing.T) {
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
		r := MustNewAnd(c.ruleA, c.ruleB)
		got := r.GetFields()
		sort.Strings(got)
		sort.Strings(c.want)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("GetFields() got: %s, want: %s", got, c.want)
		}
	}
}
