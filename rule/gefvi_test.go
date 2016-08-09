package rule

import (
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestGEFVIString(t *testing.T) {
	field := "income"
	value := int64(893)
	want := "income >= 893"
	r := NewGEFVI(field, value)
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestGEFVIGetInNiParts(t *testing.T) {
	field := "income"
	value := int64(893)
	r := NewGEFVI(field, value)
	a, b, c := r.GetInNiParts()
	if a || b != "" || c != "" {
		t.Errorf("GetInNiParts() got: %t,\"%s\",\"%s\" - want: %t,\"\",\"\"",
			a, b, c, false)
	}
}

func TestGEFVIGetTweakableParts(t *testing.T) {
	field := "income"
	value := int64(893)
	r := NewGEFVI(field, value)
	a, b, c := r.GetTweakableParts()
	if a != field || b != ">=" || c != "893" {
		t.Errorf("GetInNiParts() got: \"%s\",\"%s\",\"%s\" - want: \"%s\",\">=\",\"8.93\"",
			a, b, c, field)
	}
}

func TestGEFVIIsTrue(t *testing.T) {
	cases := []struct {
		field string
		value int64
		want  bool
	}{
		{"income", 19, true},
		{"income", 18, true},
		{"income", 20, false},
		{"income", -20, true},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"cost":   dlit.MustNew(20),
		"flow":   dlit.MustNew(124.564),
	}
	for _, c := range cases {
		r := NewGEFVI(c.field, c.value)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) (rule: %s) err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestGEFVIIsTrue_errors(t *testing.T) {
	cases := []struct {
		field   string
		value   int64
		wantErr error
	}{
		{"fred", 7, InvalidRuleError("fred >= 7")},
		{"band", 7894, InvalidRuleError("band >= 7894")},
		{"flow", 7894, InvalidRuleError("flow >= 7894")},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"flow":   dlit.MustNew(124.564),
		"band":   dlit.NewString("alpha"),
	}
	for _, c := range cases {
		r := NewGEFVI(c.field, c.value)
		_, err := r.IsTrue(record)
		if err != c.wantErr {
			t.Errorf("IsTrue(record) rule: %s, err: %v, want: %v", r, err, c.wantErr)
		}
	}
}

func TestGEFVICloneWithValue(t *testing.T) {
	field := "income"
	value := int64(893)
	r := NewGEFVI(field, value)
	want := "income >= -27489"
	cr := r.CloneWithValue(int64(-27489))
	got := cr.String()
	if got != want {
		t.Errorf("CloseWithValue() got: %s, want: %s", got, want)
	}
}

func TestGEFVICloneWithValue_panics(t *testing.T) {
	wantPanic := "can't clone with newValue: fred of type string, need type int64"
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("CloneWithValue() didn't panic")
		} else if r.(string) != wantPanic {
			t.Errorf("CloseWithValue() - got panic: %s, wanted: %s",
				r, wantPanic)
		}
	}()
	field := "income"
	value := int64(893)
	r := NewGEFVI(field, value)
	r.CloneWithValue("fred")
}
