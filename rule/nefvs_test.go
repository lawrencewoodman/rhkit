package rule

import (
	"errors"
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestNEFVSString(t *testing.T) {
	cases := []struct {
		value string
		want  string
	}{
		{value: "borris", want: "name != \"borris\""},
		{value: "bo   rris", want: "name != \"bo   rris\""},
		{value: "  borris  ", want: "name != \"  borris  \""},
		{value: "", want: "name != \"\""},
		{value: "-232.4", want: "name != \"-232.4\""},
	}
	field := "name"
	for _, c := range cases {
		r := NewNEFVS(field, c.value)
		got := r.String()
		if got != c.want {
			t.Errorf("String() got: %s, want: %s", got, c.want)
		}
	}
}

func TestNEFVSGetInNiParts(t *testing.T) {
	field := "name"
	value := "borris"
	r := NewNEFVS(field, value)
	a, b, c := r.GetInNiParts()
	if a || b != "" || c != "" {
		t.Errorf("GetInNiParts() got: %t,\"%s\",\"%s\" - want: %t,\"\",\"\"",
			a, b, c, false)
	}
}

func TestNEFVSIsTrue(t *testing.T) {
	cases := []struct {
		field string
		value string
		want  bool
	}{
		{"band", "alpha", false},
		{"band", " alpha", true},
		{"band", "alpha ", true},
		{"band", "ALPHA", true},
		{"band", "Alpha", true},
		{"success", "TRUE", false},
		{"success", "true", true},
		{"success", "1", true},
		{"income", "19", false},
		{"flow", "7.893", false},
		{"flow", "7.894", true},
		{"flow", "7.8930", true},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"flow":    dlit.MustNew(7.893),
		"success": dlit.MustNew("TRUE"),
		"band":    dlit.MustNew("alpha"),
	}
	for _, c := range cases {
		r := NewNEFVS(c.field, c.value)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestNEFVSIsTrue_errors(t *testing.T) {
	cases := []struct {
		field   string
		value   string
		wantErr error
	}{
		{field: "fred",
			value:   "hello",
			wantErr: InvalidRuleError{Rule: NewNEFVS("fred", "hello")},
		},
		{field: "problem",
			value:   "hi",
			wantErr: IncompatibleTypesRuleError{Rule: NewNEFVS("problem", "hi")},
		},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"flow":    dlit.MustNew(124.564),
		"band":    dlit.NewString("alpha"),
		"problem": dlit.MustNew(errors.New("this is an error")),
	}
	for _, c := range cases {
		r := NewNEFVS(c.field, c.value)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}
