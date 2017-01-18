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
		{ruleA: NewLEFVF("flow", 1.05), ruleB: NewGEFVF("flow", 0.0)},
		{ruleA: NewLEFVF("flow", 1.05), ruleB: NewEQFF("flow", "rate")},
		{ruleA: NewEQFF("flow", "rate"), ruleB: NewLEFVF("flow", 1.05)},
		{ruleA: NewLEFVF("flow", 1.05), ruleB: NewLEFVF("flow", 2.70)},
		{ruleA: NewLEFVF("rate", 1.05), ruleB: NewLEFVF("flow", 2.70)},
		{ruleA: NewLEFVF("flow", 1.05), ruleB: NewGEFVF("flow", 1.05)},
		{ruleA: NewLEFVF("rate", 1.05), ruleB: NewGEFVF("flow", 1.05)},
		{ruleA: NewLEFVF("flow", 1.05), ruleB: NewGEFVF("flow", 2.1)},
		{ruleA: NewLEFVF("flow", 1.05), ruleB: NewGEFVF("rate", 2.1)},
		{ruleA: NewLEFVF("flow", 1.05), ruleB: NewGEFVF("flow", 1.05)},
		{ruleA: NewGEFVF("flow", 2.1), ruleB: NewLEFVF("flow", 1.05)},
		{ruleA: NewGEFVF("flow", 2.1), ruleB: NewLEFVF("flow", 2.1)},
		{ruleA: NewGEFVF("flow", 1.05), ruleB: NewLEFVF("flow", 2.1)},
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
		{ruleA: NewLEFVF("flow", 1.05),
			ruleB:      NewTrue(),
			wantErrStr: "can't Or rule: flow <= 1.05, with: true()",
		},
		{ruleA: NewTrue(),
			ruleB:      NewLEFVF("flow", 1.05),
			wantErrStr: "can't Or rule: true(), with: flow <= 1.05",
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
		{ruleA: MustNewAnd(NewEQFVI("income", 5), NewEQFF("flow", "flow")),
			ruleB: MustNewAnd(NewEQFVI("cost", 6), NewEQFF("flow", "flow")),
			want:  "(income == 5 && flow == flow) || (cost == 6 && flow == flow)",
		},
		{ruleA: MustNewOr(NewEQFVI("income", 5), NewEQFF("flow", "flow")),
			ruleB: MustNewOr(NewEQFVI("cost", 6), NewEQFF("flow", "flow")),
			want:  "(income == 5 || flow == flow) || (cost == 6 || flow == flow)",
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
