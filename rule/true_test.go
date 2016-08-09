package rule

import (
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestTrueString(t *testing.T) {
	want := "true()"
	r := NewTrue()
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestTrueGetInNiParts(t *testing.T) {
	r := NewTrue()
	a, b, c := r.GetInNiParts()
	if a || b != "" || c != "" {
		t.Errorf("GetInNiParts() got: %t,\"%s\",\"%s\" - want: %t,\"\",\"\"",
			a, b, c, false)
	}
}

func TestTrueIsTrue(t *testing.T) {
	record := map[string]*dlit.Literal{
		"station": dlit.MustNew("harry"),
		"flow":    dlit.MustNew(7.892),
	}
	r := NewTrue()
	want := true
	got, err := r.IsTrue(record)
	if err != nil {
		t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
	}
	if got != want {
		t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, want)
	}
}
