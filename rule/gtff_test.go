package rule

import (
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestGTFFString(t *testing.T) {
	fieldA := "income"
	fieldB := "cost"
	want := "income > cost"
	r := NewGTFF(fieldA, fieldB)
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestGTFFGetInNiParts(t *testing.T) {
	fieldA := "income"
	fieldB := "cost"
	r := NewGTFF(fieldA, fieldB)
	a, b, c := r.GetInNiParts()
	if a || b != "" || c != "" {
		t.Errorf("GetInNiParts() got: %t,\"%s\",\"%s\" - want: %t,\"\",\"\"",
			a, b, c, false)
	}
}

func TestGTFFIsTrue(t *testing.T) {
	cases := []struct {
		fieldA string
		fieldB string
		want   bool
	}{
		{"income", "cost", false},
		{"cost", "income", true},
		{"income", "income", false},
		{"flowIn", "flowOut", false},
		{"flowOut", "flowIn", true},
		{"flowIn", "flowIn", false},
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
		r := NewGTFF(c.fieldA, c.fieldB)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestGTFFIsTrue_errors(t *testing.T) {
	cases := []struct {
		fieldA  string
		fieldB  string
		wantErr error
	}{
		{"income", "band", InvalidRuleError("income > band")},
		{"band", "income", InvalidRuleError("band > income")},
		{"flow", "band", InvalidRuleError("flow > band")},
		{"band", "flow", InvalidRuleError("band > flow")},
		{"fred", "income", InvalidRuleError("fred > income")},
		{"income", "fred", InvalidRuleError("income > fred")},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"flow":   dlit.MustNew(124.564),
		"band":   dlit.NewString("alpha"),
	}
	for _, c := range cases {
		r := NewGTFF(c.fieldA, c.fieldB)
		_, err := r.IsTrue(record)
		if err != c.wantErr {
			t.Errorf("IsTrue(record) rule: %s, err: %v, want: %v", r, err, c.wantErr)
		}
	}
}
