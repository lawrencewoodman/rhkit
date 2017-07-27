package goal

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestNew_errors(t *testing.T) {
	exprStr := "543d)"
	wantErr := InvalidGoalError(exprStr)
	_, err := New(exprStr)
	if err == nil || err != wantErr {
		t.Errorf("New(%v) err: %v, want: %v", exprStr, err, wantErr)
	}
}

func TestMustNew_panic(t *testing.T) {
	paniced := false
	wantPanic := "can't create goal: invalid goal: 7+8{"
	defer func() {
		if r := recover(); r != nil {
			if r.(string) == wantPanic {
				paniced = true
			} else {
				t.Errorf("MustNew: got panic: %s, wanted: %s", r, wantPanic)
			}
		}
	}()
	MustNew("7+8{")
	if !paniced {
		t.Errorf("MustNew: failed to panic with: %s", wantPanic)
	}
}

func TestString(t *testing.T) {
	exprStr := "profit > 55"
	goal, err := New(exprStr)
	if err != nil {
		t.Errorf("New(%v) err: %v", exprStr, err)
	}
	got := goal.String()
	if got != exprStr {
		t.Errorf("String() got: %s, want: %s", got, exprStr)
	}
}

func TestInvalidGoalErrorError(t *testing.T) {
	err := InvalidGoalError("b43d)")
	want := "invalid goal: b43d)"
	got := err.Error()
	if got != want {
		t.Errorf("Error() got: %v, want: %v", got, want)
	}
}

func TestAssess(t *testing.T) {
	aggregators := map[string]*dlit.Literal{
		"totalIncome":    dlit.MustNew(5000),
		"totalCost":      dlit.MustNew(4001),
		"percentMatches": dlit.MustNew(5.235),
	}
	cases := []struct {
		goalStr string
		want    bool
	}{
		{"totalIncome > 4999", true},
		{"totalIncome > 5000", false},
		{"totalIncome + totalCost > 9000", true},
		{"totalIncome + totalCost > 9001", false},
		{"roundto(percentMatches,2) == 5.24", true},
		{"roundto(percentMatches,2) == 5.23", false},
	}
	for _, c := range cases {
		goal, err := New(c.goalStr)
		if err != nil {
			t.Errorf("New(%s) err: %s", c.goalStr, err)
		}

		got, err := goal.Assess(aggregators)
		if err != nil {
			t.Errorf("Assess(%v) goal: %s, err: %s", aggregators, goal, err)
		}
		if got != c.want {
			t.Errorf("Assess(%v) goal: %s, got: %t, want: %t",
				aggregators, goal, got, c.want)
		}
	}
}

func TestAssess_errors(t *testing.T) {
	aggregators := map[string]*dlit.Literal{
		"totalIncome":    dlit.MustNew(5000),
		"totalCost":      dlit.MustNew(4001),
		"percentMatches": dlit.MustNew(5.235),
	}
	cases := []struct {
		goalStr string
		wantErr error
	}{
		{"bob > 4999",
			dexpr.InvalidExprError{
				Expr: "bob > 4999",
				Err:  dexpr.VarNotExistError("bob"),
			},
		},
		{"roundbob(percentMatches,2) == 5.23",
			dexpr.InvalidExprError{
				Expr: "roundbob(percentMatches,2) == 5.23",
				Err:  dexpr.FunctionNotExistError("roundbob"),
			},
		},
	}
	for _, c := range cases {
		goal, err := New(c.goalStr)
		if err != nil {
			t.Errorf("New(%s) err: %s", c.goalStr, err)
		}

		_, err = goal.Assess(aggregators)
		if err != c.wantErr {
			t.Errorf("Assess(%v) goal: %s, wantErr: %s, gotErr: %s",
				aggregators, goal, c.wantErr, err)
		}
	}
}

func TestMakeGoals(t *testing.T) {
	cases := []struct {
		exprs []string
		want  []*Goal
	}{
		{exprs: []string{
			"profit > 27",
			"cost <= 37",
			"numMatches >= 1500",
		},
			want: []*Goal{
				MustNew("profit > 27"),
				MustNew("cost <= 37"),
				MustNew("numMatches >= 1500"),
			},
		},
		{exprs: []string{}},
	}
	for _, c := range cases {
		got, err := MakeGoals(c.exprs)
		if err != nil {
			t.Fatalf("MakeGoals: %s", err)
		}
		if len(got) != len(c.exprs) {
			t.Fatalf("MakeGoals got: %s, want: %s", got, c.want)
		}
		for i, g := range got {
			if g.String() != c.want[i].String() {
				t.Fatalf("MakeGoals got: %s, want: %s", got, c.want)
			}
		}
	}
}

func TestMakeGoals_errors(t *testing.T) {
	exprs := []string{
		"job == \"manager\"",
		"age > > 27",
		"balance <= 1500",
	}
	wantErr := InvalidGoalError("age > > 27")
	_, err := MakeGoals(exprs)
	if err == nil || err.Error() != wantErr.Error() {
		t.Fatalf("MakeGoals err: %s, wantErr: %s", err, wantErr)
	}
}
