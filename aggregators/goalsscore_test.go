package aggregators

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/goal"
	"testing"
)

func TestGoalsScoreSpecName(t *testing.T) {
	name := "a"
	as := MustNew(name, "goalsscore")
	got := as.Name()
	if got != name {
		t.Errorf("Name - got: %s, want: %s", got, name)
	}
}

func TestGoalsScoreSpecKind(t *testing.T) {
	kind := "goalsscore"
	as := MustNew("a", kind)
	got := as.Kind()
	if got != kind {
		t.Errorf("Kind - got: %s, want: %s", got, kind)
	}
}

func TestGoalsScoreSpecArg(t *testing.T) {
	arg := ""
	as := MustNew("a", "goalsscore")
	got := as.Arg()
	if got != arg {
		t.Errorf("Arg - got: %s, want: %s", got, arg)
	}
}

func TestGoalsScoreNextRecord(t *testing.T) {
	as := MustNew("a", "goalsscore")
	ai := as.New()
	record := map[string]*dlit.Literal{}
	got := ai.NextRecord(record, true)
	if got != nil {
		t.Errorf("NextRecord: got: %s, want: %s", got, nil)
	}
}

func TestGoalsScoreResult(t *testing.T) {
	aggregatorSpecs := []AggregatorSpec{
		MustNew("income", "calc", "3 + 4"),
		MustNew("costs", "calc", "5 + 6"),
		MustNew("profit", "calc", "costs - income"),
		MustNew("highest", "calc", "7+10"),
		MustNew("lowest", "calc", "5"),
		MustNew("mid", "calc", "3"),
		MustNew("goalsScore", "goalsscore"),
	}
	cases := []struct {
		goals []*goal.Goal
		want  *dlit.Literal
	}{
		{goals: []*goal.Goal{
			goal.MustNew("income > 10"),
			goal.MustNew("costs < 2"),
			goal.MustNew("profit >= 6"),
			goal.MustNew("highest >= 18"),
			goal.MustNew("lowest <= 4"),
		},
			want: dlit.MustNew(0),
		},
		{goals: []*goal.Goal{
			goal.MustNew("income > 10"),
			goal.MustNew("costs < 2"),
			goal.MustNew("profit >= 2"),
			goal.MustNew("highest >= 15"),
			goal.MustNew("lowest <= 4"),
		},
			want: dlit.MustNew(0.002),
		},
		{goals: []*goal.Goal{
			goal.MustNew("income > 6"),
			goal.MustNew("costs < 2"),
			goal.MustNew("profit >= 2"),
		},
			want: dlit.MustNew(1.001),
		},
		{goals: []*goal.Goal{
			goal.MustNew("income > 6"),
			goal.MustNew("costs < 20"),
			goal.MustNew("profit >= 10"),
			goal.MustNew("highest >= 15"),
			goal.MustNew("lowest <= 7"),
		},
			want: dlit.MustNew(2.002),
		},
		{goals: []*goal.Goal{
			goal.MustNew("income > 6"),
			goal.MustNew("costs < 20"),
			goal.MustNew("profit >= nothing"),
			goal.MustNew("highest >= 15"),
			goal.MustNew("lowest <= 7"),
		},
			want: dlit.MustNew(dexpr.InvalidExprError{
				Expr: "profit >= nothing",
				Err:  dexpr.VarNotExistError("nothing"),
			}),
		},
	}
	numRecords := int64(12)
	instances := make([]AggregatorInstance, len(aggregatorSpecs))
	for i, aggregatorSpec := range aggregatorSpecs {
		instances[i] = aggregatorSpec.New()
	}
	goalsScoreInstance := instances[len(instances)-1]
	for i, c := range cases {
		got := goalsScoreInstance.Result(instances, c.goals, numRecords)
		if got.String() != c.want.String() {
			t.Errorf("(%d) Result: got: %s, want: %s", i, got, c.want)
		}
	}
}

func TestGoalsScoreResult_aggregator_error(t *testing.T) {
	aggregatorSpecs := []AggregatorSpec{
		MustNew("mid", "calc", "a+e"),
		MustNew("goalsScore", "goalsscore"),
	}
	goals := []*goal.Goal{
		goal.MustNew("mid > 5"),
	}
	want := dlit.MustNew(dexpr.InvalidExprError{
		Expr: "a+e",
		Err:  dexpr.VarNotExistError("a"),
	})
	numRecords := int64(12)
	instances := make([]AggregatorInstance, len(aggregatorSpecs))
	for i, aggregatorSpec := range aggregatorSpecs {
		instances[i] = aggregatorSpec.New()
	}
	goalsScoreInstance := instances[len(instances)-1]
	got := goalsScoreInstance.Result(instances, goals, numRecords)
	if got.String() != want.String() {
		t.Errorf("Result: got: %s, want: %s", got, want)
	}
}
