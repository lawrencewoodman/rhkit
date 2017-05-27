package aggregators

import (
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
