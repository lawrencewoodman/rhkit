package rule

import (
	"errors"
	"github.com/lawrencewoodman/dlit"
	"reflect"
	"testing"
)

func TestNEFFString(t *testing.T) {
	fieldA := "income"
	fieldB := "cost"
	want := "income != cost"
	r := NewNEFF(fieldA, fieldB)
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestNEFFIsTrue(t *testing.T) {
	cases := []struct {
		fieldA string
		fieldB string
		want   bool
	}{
		{"income", "cost", true},
		{"cost", "income", true},
		{"cost", "band", true},
		{"income", "income", false},
		{"flowIn", "flowOut", true},
		{"flowOut", "flowIn", true},
		{"flowIn", "flowIn", false},
		{"flowIn", "band", true},
		{"income", "flowIn", true},
		{"flowIn", "income", true},
		{"band", "band", false},
		{"band", "trueA", true},
		{"trueA", "trueB", true},
		{"trueA", "trueA", false},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"cost":    dlit.MustNew(20),
		"flowIn":  dlit.MustNew(124.564),
		"flowOut": dlit.MustNew(124.565),
		"band":    dlit.MustNew("alpha"),
		"trueA":   dlit.MustNew("true"),
		"trueB":   dlit.MustNew("TRUE"),
	}
	for _, c := range cases {
		r := NewNEFF(c.fieldA, c.fieldB)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestNEFFIsTrue_errors(t *testing.T) {
	cases := []struct {
		fieldA  string
		fieldB  string
		wantErr error
	}{
		{fieldA: "fred",
			fieldB:  "income",
			wantErr: InvalidRuleError{NewNEFF("fred", "income")},
		},
		{fieldA: "income",
			fieldB:  "fred",
			wantErr: InvalidRuleError{NewNEFF("income", "fred")},
		},
		{fieldA: "income",
			fieldB:  "problem",
			wantErr: IncompatibleTypesRuleError{NewNEFF("income", "problem")},
		},
		{fieldA: "problem",
			fieldB:  "income",
			wantErr: IncompatibleTypesRuleError{NewNEFF("problem", "income")},
		},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"flow":    dlit.MustNew(124.564),
		"band":    dlit.NewString("alpha"),
		"problem": dlit.MustNew(errors.New("this is an error")),
	}
	for _, c := range cases {
		r := NewNEFF(c.fieldA, c.fieldB)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}

func TestNEFFFields(t *testing.T) {
	r := NewNEFF("income", "cost")
	want := []string{"income", "cost"}
	got := r.Fields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Fields() got: %s, want: %s", got, want)
	}
}
