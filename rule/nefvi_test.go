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
		{"fred", 8, InvalidRuleError("fred != 8")},
		{"band", 8, InvalidRuleError("band != 8")},
		{"flow", 8, InvalidRuleError("flow != 8")},
		{"problem", 8, InvalidRuleError("problem != 8")},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"flow":    dlit.MustNew(124.564),
		"band":    dlit.NewString("alpha"),
		"problem": dlit.MustNew(errors.New("this is an error")),
	}
	for _, c := range cases {
		r := NewNEFVI(c.field, c.value)
		_, err := r.IsTrue(record)
		if err != c.wantErr {
			t.Errorf("IsTrue(record) rule: %s, err: %v, want: %v", r, err, c.wantErr)
		}
	}
}
