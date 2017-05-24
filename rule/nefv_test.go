package rule

import (
	"errors"
	"github.com/lawrencewoodman/dlit"
	"reflect"
	"testing"
)

func TestNEFVString(t *testing.T) {
	cases := []struct {
		field string
		value *dlit.Literal
		want  string
	}{
		{field: "income", value: dlit.MustNew(7.8903), want: "income != 7.8903"},
		{field: "income", value: dlit.MustNew(7.890300), want: "income != 7.8903"},
		{field: "income", value: dlit.MustNew(7.), want: "income != 7"},
		{field: "income", value: dlit.MustNew(7.00), want: "income != 7"},
		{field: "income", value: dlit.MustNew(7), want: "income != 7"},
		{field: "income", value: dlit.MustNew(0.34), want: "income != 0.34"},
		{field: "income", value: dlit.MustNew(0.3400), want: "income != 0.34"},
		{field: "income", value: dlit.MustNew(0.0), want: "income != 0"},
		{field: "income", value: dlit.MustNew(-53.9), want: "income != -53.9"},
		{field: "name", value: dlit.MustNew("borris"),
			want: "name != \"borris\""},
		{field: "name", value: dlit.MustNew("bo   rris"),
			want: "name != \"bo   rris\""},
		{field: "name", value: dlit.MustNew("  borris  "),
			want: "name != \"  borris  \""},
		{field: "name", value: dlit.MustNew(""), want: "name != \"\""},
	}
	for _, c := range cases {
		r := NewNEFV(c.field, c.value)
		got := r.String()
		if got != c.want {
			t.Errorf("String() got: %s, want: %s", got, c.want)
		}
	}
}

func TestNEFVIsTrue(t *testing.T) {
	cases := []struct {
		field string
		value *dlit.Literal
		want  bool
	}{
		{"income", dlit.MustNew(19.0), false},
		{"income", dlit.MustNew(-19.0), true},
		{"income", dlit.MustNew(20.0), true},
		{"flow", dlit.MustNew(124.564), false},
		{"flow", dlit.MustNew(-124.564), true},
		{"flow", dlit.MustNew(20.0), true},
		{"flow", dlit.MustNew(124.5645), true},
		{"flow", dlit.MustNew(124.565), true},
		{"band", dlit.MustNew("hello"), true},
		{"band", dlit.MustNew("alpha"), false},
		{"band", dlit.MustNew("ALPHA"), true},
		{"band", dlit.MustNew(8.9), true},
		{"success", dlit.MustNew("TRUE"), false},
		{"success", dlit.MustNew("true"), true},
		{"success", dlit.MustNew("1"), true},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"flow":    dlit.MustNew(124.564),
		"band":    dlit.MustNew("alpha"),
		"success": dlit.MustNew("TRUE"),
	}
	for _, c := range cases {
		r := NewNEFV(c.field, c.value)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestNEFVIsTrue_errors(t *testing.T) {
	cases := []struct {
		field   string
		value   *dlit.Literal
		wantErr error
	}{
		{"fred", dlit.MustNew(8.9),
			InvalidRuleError{Rule: NewNEFV("fred", dlit.MustNew(8.9))}},
		{"problem", dlit.MustNew(8.9),
			IncompatibleTypesRuleError{Rule: NewNEFV("problem", dlit.MustNew(8.9))}},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"flow":    dlit.MustNew(124.564),
		"band":    dlit.NewString("alpha"),
		"problem": dlit.MustNew(errors.New("this is an error")),
	}
	for _, c := range cases {
		r := NewNEFV(c.field, c.value)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}

func TestNEFVGetFields(t *testing.T) {
	r := NewNEFV("income", dlit.MustNew(5.5))
	want := []string{"income"}
	got := r.GetFields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetFields() got: %s, want: %s", got, want)
	}
}
