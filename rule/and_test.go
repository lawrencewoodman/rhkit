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
		{ruleA: NewLEFVF("flow", 1.05), ruleB: NewGEFVF("flow", 0.0)},
		{ruleA: NewLEFVF("flow", 1.05), ruleB: NewEQFF("flow", "rate")},
		{ruleA: NewEQFF("flow", "rate"), ruleB: NewLEFVF("flow", 1.05)},
		{ruleA: NewLEFVF("rate", 1.05), ruleB: NewLEFVF("flow", 2.70)},
		{ruleA: NewLEFVF("rate", 1.05), ruleB: NewGEFVF("flow", 1.05)},
		{ruleA: NewGEFVF("flow", 1.05), ruleB: NewLEFVF("flow", 2.1)},
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
		{ruleA: NewLEFVF("flow", 1.05),
			ruleB:      NewLEFVF("flow", 2.70),
			wantErrStr: "can't And rule: flow <= 1.05, with: flow <= 2.7",
		},
		{ruleA: NewLEFVF("flow", 1.05),
			ruleB:      NewLEFVF("flow", 1.05),
			wantErrStr: "can't And rule: flow <= 1.05, with: flow <= 1.05",
		},
		{ruleA: NewGEFVF("flow", 1.05),
			ruleB:      NewGEFVF("flow", 1.05),
			wantErrStr: "can't And rule: flow >= 1.05, with: flow >= 1.05",
		},
		{ruleA: NewLEFVF("flow", 1.05),
			ruleB:      NewGEFVF("flow", 1.05),
			wantErrStr: "can't And rule: flow <= 1.05, with: flow >= 1.05",
		},
		{ruleA: NewLEFVF("flow", 1.05),
			ruleB:      NewGEFVF("flow", 2.1),
			wantErrStr: "can't And rule: flow <= 1.05, with: flow >= 2.1",
		},
		{ruleA: NewGEFVF("flow", 2.1),
			ruleB:      NewLEFVF("flow", 1.05),
			wantErrStr: "can't And rule: flow >= 2.1, with: flow <= 1.05",
		},
		{ruleA: NewGEFVF("flow", 2.1),
			ruleB:      NewLEFVF("flow", 2.1),
			wantErrStr: "can't And rule: flow >= 2.1, with: flow <= 2.1",
		},
		{ruleA: NewInFV("group", []*dlit.Literal{
			dlit.NewString("bob"),
			dlit.NewString("fred"),
			dlit.NewString("albert"),
		}),
			ruleB:      NewEQFVS("group", "norris"),
			wantErrStr: "can't And rule: in(group,\"bob\",\"fred\",\"albert\"), with: group == \"norris\"",
		},
		{ruleA: NewEQFVS("group", "norris"),
			ruleB: NewInFV("group", []*dlit.Literal{
				dlit.NewString("bob"),
				dlit.NewString("fred"),
				dlit.NewString("albert"),
			}),
			wantErrStr: "can't And rule: group == \"norris\", with: in(group,\"bob\",\"fred\",\"albert\")",
		},
		{ruleA: NewTrue(),
			ruleB:      NewEQFF("flow", "rate"),
			wantErrStr: "can't And rule: true(), with: flow == rate",
		},
		{ruleA: NewEQFF("flow", "rate"),
			ruleB:      NewTrue(),
			wantErrStr: "can't And rule: flow == rate, with: true()",
		},
	}
	for _, c := range cases {
		r, err := NewAnd(c.ruleA, c.ruleB)
		if r != nil {
			t.Errorf("NewAnd(%s, %s) rule got: %s, want: nil", c.ruleA, c.ruleB, r)
		}
		if err == nil {
			t.Errorf("NewAnd(%s, %s) got err: nil, want: %s",
				c.ruleA, c.ruleB, c.wantErrStr)
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
		{ruleA: MustNewAnd(NewEQFVI("income", 5), NewEQFF("flow", "flow")),
			ruleB: MustNewAnd(NewEQFVI("cost", 6), NewEQFF("flow", "flow")),
			want:  "(income == 5 && flow == flow) && (cost == 6 && flow == flow)",
		},
		{ruleA: MustNewOr(NewEQFVI("income", 5), NewEQFF("flow", "flow")),
			ruleB: MustNewOr(NewEQFVI("cost", 6), NewEQFF("flow", "flow")),
			want:  "(income == 5 || flow == flow) && (cost == 6 || flow == flow)",
		},
		{ruleA: MustNewBetweenFVI("flow", 5, 10),
			ruleB: NewEQFF("income", "cost"),
			want:  "(flow >= 5 && flow <= 10) && income == cost",
		},
		{ruleA: NewEQFF("income", "cost"),
			ruleB: MustNewBetweenFVI("flow", 5, 10),
			want:  "income == cost && (flow >= 5 && flow <= 10)",
		},
		{ruleA: MustNewBetweenFVF("flow", 5.24, 10.89),
			ruleB: NewEQFF("income", "cost"),
			want:  "(flow >= 5.24 && flow <= 10.89) && income == cost",
		},
		{ruleA: NewEQFF("income", "cost"),
			ruleB: MustNewBetweenFVF("flow", 5.24, 10.89),
			want:  "income == cost && (flow >= 5.24 && flow <= 10.89)",
		},
		{ruleA: MustNewOutsideFVI("flow", 5, 10),
			ruleB: NewEQFF("income", "cost"),
			want:  "(flow <= 5 || flow >= 10) && income == cost",
		},
		{ruleA: NewEQFF("income", "cost"),
			ruleB: MustNewOutsideFVI("flow", 5, 10),
			want:  "income == cost && (flow <= 5 || flow >= 10)",
		},
		{ruleA: MustNewOutsideFVF("flow", 5.24, 10.89),
			ruleB: NewEQFF("income", "cost"),
			want:  "(flow <= 5.24 || flow >= 10.89) && income == cost",
		},
		{ruleA: NewEQFF("income", "cost"),
			ruleB: MustNewOutsideFVF("flow", 5.24, 10.89),
			want:  "income == cost && (flow <= 5.24 || flow >= 10.89)",
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
