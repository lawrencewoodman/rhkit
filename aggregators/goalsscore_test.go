package aggregators

import (
	"testing"
)

func TestGoalsScoreSpecGetName(t *testing.T) {
	name := "a"
	as := MustNew(name, "goalsscore")
	got := as.GetName()
	if got != name {
		t.Errorf("GetName - got: %s, want: %s", got, name)
	}
}

func TestGoalsScoreSpecGetKind(t *testing.T) {
	kind := "goalsscore"
	as := MustNew("a", kind)
	got := as.GetKind()
	if got != kind {
		t.Errorf("GetKind - got: %s, want: %s", got, kind)
	}
}

func TestGoalsScoreSpecGetArg(t *testing.T) {
	arg := ""
	as := MustNew("a", "goalsscore")
	got := as.GetArg()
	if got != arg {
		t.Errorf("GetArg - got: %s, want: %s", got, arg)
	}
}
