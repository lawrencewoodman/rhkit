package rule

import (
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestGEFVFString(t *testing.T) {
	field := "income"
	value := 8.93
	want := "income >= 8.93"
	r := NewGEFVF(field, value)
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestGEFVFGetInNiParts(t *testing.T) {
	field := "income"
	value := 8.93
	r := NewGEFVF(field, value)
	a, b, c := r.GetInNiParts()
	if a || b != "" || c != "" {
		t.Errorf("GetInNiParts() got: %t,\"%s\",\"%s\" - want: %t,\"\",\"\"",
			a, b, c, false)
	}
}

func TestGEFVFGetTweakableParts(t *testing.T) {
	field := "income"
	value := 8.93
	r := NewGEFVF(field, value)
	a, b, c := r.GetTweakableParts()
	if a != field || b != ">=" || c != "8.93" {
		t.Errorf("GetInNiParts() got: \"%s\",\"%s\",\"%s\" - want: \"%s\",\">=\",\"8.93\"",
			a, b, c, field)
	}
}

func TestGEFVFIsTrue(t *testing.T) {
	cases := []struct {
		field string
		value float64
		want  bool
	}{
		{"income", 19, true},
		{"income", 19.12, false},
		{"income", 20, false},
		{"income", -20, true},
		{"income", 18.34, true},
		{"flow", 124.564, true},
		{"flow", 124.565, false},
		{"flow", 124.563, true},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"cost":   dlit.MustNew(20),
		"flow":   dlit.MustNew(124.564),
	}
	for _, c := range cases {
		r := NewGEFVF(c.field, c.value)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) (rule: %s) err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestGEFVFIsTrue_errors(t *testing.T) {
	cases := []struct {
		field   string
		value   float64
		wantErr error
	}{
		{field: "fred",
			value:   7.894,
			wantErr: InvalidRuleError{Rule: NewGEFVF("fred", 7.894)},
		},
		{field: "band",
			value:   7.894,
			wantErr: IncompatibleTypesRuleError{Rule: NewGEFVF("band", 7.894)},
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"flow":   dlit.MustNew(124.564),
		"band":   dlit.NewString("alpha"),
	}
	for _, c := range cases {
		r := NewGEFVF(c.field, c.value)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}

func TestGEFVFCloneWithValue(t *testing.T) {
	field := "income"
	value := 8.93
	r := NewGEFVF(field, value)
	want := "income >= -27.489"
	cr := r.CloneWithValue(-27.489)
	got := cr.String()
	if got != want {
		t.Errorf("CloseWithValue() got: %s, want: %s", got, want)
	}
}

func TestGEFVFCloneWithValue_panics(t *testing.T) {
	wantPanic := "can't clone with newValue: fred of type string, need type float64"
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("CloneWithValue() didn't panic")
		} else if r.(string) != wantPanic {
			t.Errorf("CloseWithValue() - got panic: %s, wanted: %s",
				r, wantPanic)
		}
	}()
	field := "income"
	value := 8.93
	r := NewGEFVF(field, value)
	r.CloneWithValue("fred")
}
