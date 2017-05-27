package rule

import (
	"github.com/lawrencewoodman/dlit"
	"reflect"
	"testing"
)

func TestGEFFString(t *testing.T) {
	fieldA := "income"
	fieldB := "cost"
	want := "income >= cost"
	r := NewGEFF(fieldA, fieldB)
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestGEFFIsTrue(t *testing.T) {
	cases := []struct {
		fieldA string
		fieldB string
		want   bool
	}{
		{"income", "cost", false},
		{"cost", "income", true},
		{"income", "income", true},
		{"flowIn", "flowOut", false},
		{"flowOut", "flowIn", true},
		{"flowIn", "flowIn", true},
		{"income", "flowIn", false},
		{"flowIn", "income", true},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"cost":    dlit.MustNew(20),
		"flowIn":  dlit.MustNew(124.564),
		"flowOut": dlit.MustNew(124.565),
	}
	for _, c := range cases {
		r := NewGEFF(c.fieldA, c.fieldB)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestGEFFIsTrue_errors(t *testing.T) {
	cases := []struct {
		fieldA  string
		fieldB  string
		wantErr error
	}{
		{fieldA: "income",
			fieldB:  "band",
			wantErr: IncompatibleTypesRuleError{Rule: NewGEFF("income", "band")},
		},
		{fieldA: "band",
			fieldB:  "income",
			wantErr: IncompatibleTypesRuleError{Rule: NewGEFF("band", "income")},
		},
		{fieldA: "flow",
			fieldB:  "band",
			wantErr: IncompatibleTypesRuleError{Rule: NewGEFF("flow", "band")},
		},
		{fieldA: "band",
			fieldB:  "flow",
			wantErr: IncompatibleTypesRuleError{Rule: NewGEFF("band", "flow")},
		},
		{fieldA: "fred",
			fieldB:  "income",
			wantErr: InvalidRuleError{Rule: NewGEFF("fred", "income")},
		},
		{fieldA: "income",
			fieldB:  "fred",
			wantErr: InvalidRuleError{Rule: NewGEFF("income", "fred")},
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"flow":   dlit.MustNew(124.564),
		"band":   dlit.NewString("alpha"),
	}
	for _, c := range cases {
		r := NewGEFF(c.fieldA, c.fieldB)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}

func TestGEFFFields(t *testing.T) {
	r := NewGEFF("income", "cost")
	want := []string{"income", "cost"}
	got := r.Fields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Fields() got: %s, want: %s", got, want)
	}
}
