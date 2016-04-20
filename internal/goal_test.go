package internal

import (
	"github.com/lawrencewoodman/dlit_go"
	"testing"
)

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
		goal, err := NewGoal(c.goalStr)
		if err != nil {
			t.Errorf("NewGoal(%s) err: %s", c.goalStr, err)
		}

		got, err := goal.Assess(aggregators)
		if err != nil {
			t.Errorf("Assess(%q) goal: %s, err: %s", aggregators, goal, err)
		}
		if got != c.want {
			t.Errorf("Assess(%q) goal: %s, got: %s, want: %s",
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
		wantErr string
	}{
		{"bob > 4999", "Variable doesn't exist: bob"},
		{"roundbob(percentMatches,2) == 5.23", "Function doesn't exist: roundbob"},
	}
	for _, c := range cases {
		goal, err := NewGoal(c.goalStr)
		if err != nil {
			t.Errorf("NewGoal(%s) err: %s", c.goalStr, err)
		}

		_, err = goal.Assess(aggregators)
		if err == nil || err.Error() != c.wantErr {
			t.Errorf("Assess(%s) goal: %s, wantErr, gotErr: %s",
				aggregators, goal, c.wantErr, err)
		}
	}
}
