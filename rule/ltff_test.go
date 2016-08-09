package rule

import (
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestLTFFString(t *testing.T) {
	fieldA := "income"
	fieldB := "cost"
	want := "income < cost"
	r := NewLTFF(fieldA, fieldB)
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestLTFFGetInNiParts(t *testing.T) {
	fieldA := "income"
	fieldB := "cost"
	r := NewLTFF(fieldA, fieldB)
	a, b, c := r.GetInNiParts()
	if a || b != "" || c != "" {
		t.Errorf("GetInNiParts() got: %t,\"%s\",\"%s\" - want: %t,\"\",\"\"",
			a, b, c, false)
	}
}

func TestLTFFIsTrue(t *testing.T) {
	cases := []struct {
		fieldA string
		fieldB string
		want   bool
	}{
		{"income", "cost", true},
		{"cost", "income", false},
		{"income", "income", false},
		{"flowIn", "flowOut", true},
		{"flowOut", "flowIn", false},
		{"flowIn", "flowIn", false},
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
		r := NewLTFF(c.fieldA, c.fieldB)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestLTFFIsTrue_errors(t *testing.T) {
	cases := []struct {
		fieldA  string
		fieldB  string
		wantErr error
	}{
		{"income", "band", InvalidRuleError("income < band")},
		{"band", "income", InvalidRuleError("band < income")},
		{"flow", "band", InvalidRuleError("flow < band")},
		{"band", "flow", InvalidRuleError("band < flow")},
		{"fred", "income", InvalidRuleError("fred < income")},
		{"income", "fred", InvalidRuleError("income < fred")},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"flow":   dlit.MustNew(124.564),
		"band":   dlit.NewString("alpha"),
	}
	for _, c := range cases {
		r := NewLTFF(c.fieldA, c.fieldB)
		_, err := r.IsTrue(record)
		if err != c.wantErr {
			t.Errorf("IsTrue(record) rule: %s, err: %v, want: %v", r, err, c.wantErr)
		}
	}
}

/**************************
 *  Benchmarks
 **************************/

func BenchmarkLTFFIsTrue(b *testing.B) {
	record := map[string]*dlit.Literal{
		"band":   dlit.MustNew(23),
		"income": dlit.MustNew(1024),
		"cost":   dlit.MustNew(890),
	}
	r := NewLTFF("cost", "income")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = r.IsTrue(record)
	}
}
