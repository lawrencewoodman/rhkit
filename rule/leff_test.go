package rule

import (
	"github.com/lawrencewoodman/dlit"
	"reflect"
	"testing"
)

func TestLEFFString(t *testing.T) {
	fieldA := "income"
	fieldB := "cost"
	want := "income <= cost"
	r := NewLEFF(fieldA, fieldB)
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestLEFFIsTrue(t *testing.T) {
	cases := []struct {
		fieldA string
		fieldB string
		want   bool
	}{
		{"income", "cost", true},
		{"cost", "income", false},
		{"income", "income", true},
		{"flowIn", "flowOut", true},
		{"flowOut", "flowIn", false},
		{"flowIn", "flowIn", true},
		{"income", "flowIn", true},
		{"flowIn", "income", false},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"cost":    dlit.MustNew(20),
		"flowIn":  dlit.MustNew(124.564),
		"flowOut": dlit.MustNew(124.565),
	}
	for _, c := range cases {
		r := NewLEFF(c.fieldA, c.fieldB)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestLEFFIsTrue_errors(t *testing.T) {
	cases := []struct {
		fieldA  string
		fieldB  string
		wantErr error
	}{
		{fieldA: "income",
			fieldB:  "band",
			wantErr: IncompatibleTypesRuleError{Rule: NewLEFF("income", "band")},
		},
		{fieldA: "band",
			fieldB:  "income",
			wantErr: IncompatibleTypesRuleError{Rule: NewLEFF("band", "income")},
		},
		{fieldA: "flow",
			fieldB:  "band",
			wantErr: IncompatibleTypesRuleError{Rule: NewLEFF("flow", "band")},
		},
		{fieldA: "band",
			fieldB:  "flow",
			wantErr: IncompatibleTypesRuleError{Rule: NewLEFF("band", "flow")},
		},
		{fieldA: "fred",
			fieldB:  "income",
			wantErr: InvalidRuleError{Rule: NewLEFF("fred", "income")},
		},
		{fieldA: "income",
			fieldB:  "fred",
			wantErr: InvalidRuleError{Rule: NewLEFF("income", "fred")},
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"flow":   dlit.MustNew(124.564),
		"band":   dlit.NewString("alpha"),
	}
	for _, c := range cases {
		r := NewLEFF(c.fieldA, c.fieldB)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}

func TestLEFFFields(t *testing.T) {
	r := NewLEFF("income", "cost")
	want := []string{"income", "cost"}
	got := r.Fields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Fields() got: %s, want: %s", got, want)
	}
}
