package rule

import (
	"errors"
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestEQFVFString(t *testing.T) {
	cases := []struct {
		value float64
		want  string
	}{
		{value: 7.8903, want: "income == 7.8903"},
		{value: 7.890300, want: "income == 7.8903"},
		{value: 7., want: "income == 7"},
		{value: 7.00, want: "income == 7"},
		{value: 7, want: "income == 7"},
		{value: 0.34, want: "income == 0.34"},
		{value: 0.3400, want: "income == 0.34"},
		{value: 0.0, want: "income == 0"},
		{value: -53.9, want: "income == -53.9"},
	}
	field := "income"
	for _, c := range cases {
		r := NewEQFVF(field, c.value)
		got := r.String()
		if got != c.want {
			t.Errorf("String() got: %s, want: %s", got, c.want)
		}
	}
}

func TestEQFVFGetInNiParts(t *testing.T) {
	field := "income"
	value := 29.3
	r := NewEQFVF(field, value)
	a, b, c := r.GetInNiParts()
	if a || b != "" || c != "" {
		t.Errorf("GetInNiParts() got: %t,\"%s\",\"%s\" - want: %t,\"\",\"\"",
			a, b, c, false)
	}
}

func TestEQFVFIsTrue(t *testing.T) {
	cases := []struct {
		field string
		value float64
		want  bool
	}{
		{"income", 19.0, true},
		{"income", -19.0, false},
		{"income", 20.0, false},
		{"flow", 124.564, true},
		{"flow", -124.564, false},
		{"flow", 20.0, false},
		{"flow", 124.5645, false},
		{"flow", 124.565, false},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"flow":   dlit.MustNew(124.564),
		"band":   dlit.MustNew("alpha"),
	}
	for _, c := range cases {
		r := NewEQFVF(c.field, c.value)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestEQFVFIsTrue_errors(t *testing.T) {
	cases := []struct {
		field   string
		value   float64
		wantErr error
	}{
		{"fred", 8.9, InvalidRuleError("fred == 8.9")},
		{"band", 8.9, InvalidRuleError("band == 8.9")},
		{"problem", 8.9, InvalidRuleError("problem == 8.9")},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"flow":    dlit.MustNew(124.564),
		"band":    dlit.NewString("alpha"),
		"problem": dlit.MustNew(errors.New("this is an error")),
	}
	for _, c := range cases {
		r := NewEQFVF(c.field, c.value)
		_, err := r.IsTrue(record)
		if err != c.wantErr {
			t.Errorf("IsTrue(record) rule: %s, err: %v, want: %v", r, err, c.wantErr)
		}
	}
}
