package rule

import (
	"errors"
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestEQFVSString(t *testing.T) {
	cases := []struct {
		value string
		want  string
	}{
		{value: "borris", want: "name == \"borris\""},
		{value: "bo   rris", want: "name == \"bo   rris\""},
		{value: "  borris  ", want: "name == \"  borris  \""},
		{value: "", want: "name == \"\""},
		{value: "-232.4", want: "name == \"-232.4\""},
	}
	field := "name"
	for _, c := range cases {
		r := NewEQFVS(field, c.value)
		got := r.String()
		if got != c.want {
			t.Errorf("String() got: %s, want: %s", got, c.want)
		}
	}
}

func TestEQFVSGetInNiParts(t *testing.T) {
	field := "name"
	value := "borris"
	r := NewEQFVS(field, value)
	a, b, c := r.GetInNiParts()
	if a || b != "" || c != "" {
		t.Errorf("GetInNiParts() got: %t,\"%s\",\"%s\" - want: %t,\"\",\"\"",
			a, b, c, false)
	}
}

func TestEQFVSIsTrue(t *testing.T) {
	cases := []struct {
		field string
		value string
		want  bool
	}{
		{"band", "alpha", true},
		{"band", " alpha", false},
		{"band", "alpha ", false},
		{"band", "ALPHA", false},
		{"band", "Alpha", false},
		{"success", "TRUE", true},
		{"success", "true", false},
		{"success", "1", false},
		{"income", "19", true},
		{"flow", "7.893", true},
		{"flow", "7.894", false},
		{"flow", "7.8930", false},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"flow":    dlit.MustNew(7.893),
		"success": dlit.MustNew("TRUE"),
		"band":    dlit.MustNew("alpha"),
	}
	for _, c := range cases {
		r := NewEQFVS(c.field, c.value)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestEQFVSIsTrue_errors(t *testing.T) {
	cases := []struct {
		field   string
		value   string
		wantErr error
	}{
		{"fred", "hello", InvalidRuleError{Rule: NewEQFVS("fred", "hello")}},
		{"problem", "hi", IncompatibleTypesRuleError{
			Rule: NewEQFVS("problem", "hi"),
		}},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"flow":    dlit.MustNew(124.564),
		"band":    dlit.NewString("alpha"),
		"problem": dlit.MustNew(errors.New("this is an error")),
	}
	for _, c := range cases {
		r := NewEQFVS(c.field, c.value)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}
