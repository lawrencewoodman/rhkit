package rule

import (
	"errors"
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestEQFVIString(t *testing.T) {
	cases := []struct {
		value int64
		want  string
	}{
		{value: 789, want: "income == 789"},
		{value: -789, want: "income == -789"},
		{value: 0, want: "income == 0"},
	}
	field := "income"
	for _, c := range cases {
		r := NewEQFVI(field, c.value)
		got := r.String()
		if got != c.want {
			t.Errorf("String() got: %s, want: %s", got, c.want)
		}
	}
}

func TestEQFVIIsTrue(t *testing.T) {
	cases := []struct {
		field string
		value int64
		want  bool
	}{
		{"income", 19.0, true},
		{"income", -19.0, false},
		{"income", 20.0, false},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"flow":   dlit.MustNew(124.564),
		"band":   dlit.MustNew("alpha"),
	}
	for _, c := range cases {
		r := NewEQFVI(c.field, c.value)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestEQFVIIsTrue_errors(t *testing.T) {
	cases := []struct {
		field   string
		value   int64
		wantErr error
	}{
		{"fred", 8, InvalidRuleError{Rule: NewEQFVI("fred", 8)}},
		{"band", 8, IncompatibleTypesRuleError{Rule: NewEQFVI("band", 8)}},
		{"flow", 8, IncompatibleTypesRuleError{Rule: NewEQFVI("flow", 8)}},
		{"problem", 8, IncompatibleTypesRuleError{Rule: NewEQFVI("problem", 8)}},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"flow":    dlit.MustNew(124.564),
		"band":    dlit.NewString("alpha"),
		"problem": dlit.MustNew(errors.New("this is an error")),
	}
	for _, c := range cases {
		r := NewEQFVI(c.field, c.value)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}
