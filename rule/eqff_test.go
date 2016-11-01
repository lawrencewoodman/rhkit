package rule

import (
	"errors"
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestEQFFString(t *testing.T) {
	fieldA := "income"
	fieldB := "cost"
	want := "income == cost"
	r := NewEQFF(fieldA, fieldB)
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestEQFFIsTrue(t *testing.T) {
	cases := []struct {
		fieldA string
		fieldB string
		want   bool
	}{
		{"income", "cost", false},
		{"cost", "income", false},
		{"cost", "band", false},
		{"income", "income", true},
		{"flowIn", "flowOut", false},
		{"flowOut", "flowIn", false},
		{"flowIn", "flowIn", true},
		{"flowIn", "band", false},
		{"income", "flowIn", false},
		{"flowIn", "income", false},
		{"band", "band", true},
		{"band", "trueA", false},
		{"trueA", "trueB", false},
		{"trueA", "trueA", true},
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
		r := NewEQFF(c.fieldA, c.fieldB)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestEQFFIsTrue_errors(t *testing.T) {
	cases := []struct {
		fieldA  string
		fieldB  string
		wantErr error
	}{
		{fieldA: "fred",
			fieldB:  "income",
			wantErr: InvalidRuleError{NewEQFF("fred", "income")},
		},
		{fieldA: "income",
			fieldB:  "fred",
			wantErr: InvalidRuleError{NewEQFF("income", "fred")},
		},
		{fieldA: "income",
			fieldB:  "problem",
			wantErr: IncompatibleTypesRuleError{NewEQFF("income", "problem")},
		},
		{fieldA: "problem",
			fieldB:  "income",
			wantErr: IncompatibleTypesRuleError{NewEQFF("problem", "income")},
		},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"flow":    dlit.MustNew(124.564),
		"band":    dlit.NewString("alpha"),
		"problem": dlit.MustNew(errors.New("this is an error")),
	}
	for _, c := range cases {
		r := NewEQFF(c.fieldA, c.fieldB)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}
