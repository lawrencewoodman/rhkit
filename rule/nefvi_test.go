package rule

import (
	"errors"
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestNEFVIString(t *testing.T) {
	cases := []struct {
		value int64
		want  string
	}{
		{value: 789, want: "income != 789"},
		{value: -789, want: "income != -789"},
		{value: 0, want: "income != 0"},
	}
	field := "income"
	for _, c := range cases {
		r := NewNEFVI(field, c.value)
		got := r.String()
		if got != c.want {
			t.Errorf("String() got: %s, want: %s", got, c.want)
		}
	}
}

func TestNEFVIGetInNiParts(t *testing.T) {
	field := "income"
	value := int64(29)
	r := NewNEFVI(field, value)
	a, b, c := r.GetInNiParts()
	if a || b != "" || c != "" {
		t.Errorf("GetInNiParts() got: %t,\"%s\",\"%s\" - want: %t,\"\",\"\"",
			a, b, c, false)
	}
}

func TestNEFVIIsTrue(t *testing.T) {
	cases := []struct {
		field string
		value int64
		want  bool
	}{
		{"income", 19.0, false},
		{"income", -19.0, true},
		{"income", 20.0, true},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"flow":   dlit.MustNew(124.564),
		"band":   dlit.MustNew("alpha"),
	}
	for _, c := range cases {
		r := NewNEFVI(c.field, c.value)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestNEFVIIsTrue_errors(t *testing.T) {
	cases := []struct {
		field   string
		value   int64
		wantErr error
	}{
		{field: "fred",
			value:   8,
			wantErr: InvalidRuleError{Rule: NewNEFVI("fred", 8)},
		},
		{field: "band",
			value:   8,
			wantErr: IncompatibleTypesRuleError{Rule: NewNEFVI("band", 8)},
		},
		{field: "flow",
			value:   8,
			wantErr: IncompatibleTypesRuleError{Rule: NewNEFVI("flow", 8)},
		},
		{field: "problem",
			value:   8,
			wantErr: IncompatibleTypesRuleError{Rule: NewNEFVI("problem", 8)},
		},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"flow":    dlit.MustNew(124.564),
		"band":    dlit.NewString("alpha"),
		"problem": dlit.MustNew(errors.New("this is an error")),
	}
	for _, c := range cases {
		r := NewNEFVI(c.field, c.value)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}
