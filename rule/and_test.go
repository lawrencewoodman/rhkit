package rule

import (
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestAndString(t *testing.T) {
	cases := []struct {
		ruleA Rule
		ruleB Rule
		want  string
	}{
		{ruleA: NewTrue(),
			ruleB: NewEQFF("income", "cost"),
			want:  "true() && income == cost",
		},
		{ruleA: NewAnd(NewTrue(), NewTrue()),
			ruleB: NewEQFF("income", "cost"),
			want:  "(true() && true()) && income == cost",
		},
		{ruleA: NewOr(NewTrue(), NewTrue()),
			ruleB: NewEQFF("income", "cost"),
			want:  "(true() || true()) && income == cost",
		},
		{ruleA: NewEQFF("income", "cost"),
			ruleB: NewAnd(NewTrue(), NewTrue()),
			want:  "income == cost && (true() && true())",
		},
		{ruleA: NewEQFF("income", "cost"),
			ruleB: NewOr(NewTrue(), NewTrue()),
			want:  "income == cost && (true() || true())",
		},
		{ruleA: NewAnd(NewEQFVI("income", 5), NewTrue()),
			ruleB: NewAnd(NewEQFVI("cost", 6), NewTrue()),
			want:  "(income == 5 && true()) && (cost == 6 && true())",
		},
		{ruleA: NewOr(NewEQFVI("income", 5), NewTrue()),
			ruleB: NewOr(NewEQFVI("cost", 6), NewTrue()),
			want:  "(income == 5 || true()) && (cost == 6 || true())",
		},
	}
	for _, c := range cases {
		r := NewAnd(c.ruleA, c.ruleB)
		got := r.String()
		if got != c.want {
			t.Errorf("String() got: %s, want: %s", got, c.want)
		}
	}
}

func TestAndIsTrue(t *testing.T) {
	cases := []struct {
		ruleA Rule
		ruleB Rule
		want  bool
	}{
		{ruleA: NewTrue(),
			ruleB: NewTrue(),
			want:  true,
		},
		{ruleA: NewNEFF("income", "income"),
			ruleB: NewTrue(),
			want:  false,
		},
		{ruleA: NewTrue(),
			ruleB: NewNEFF("income", "income"),
			want:  false,
		},
		{ruleA: NewNEFF("income", "income"),
			ruleB: NewNEFF("income", "income"),
			want:  false,
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"cost":   dlit.MustNew(20),
	}
	for _, c := range cases {
		r := NewAnd(c.ruleA, c.ruleB)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestAndIsTrue_errors(t *testing.T) {
	cases := []struct {
		ruleA Rule
		ruleB Rule
	}{
		{ruleA: NewTrue(),
			ruleB: NewEQFF("fred", "income"),
		},
		{ruleA: NewEQFF("fred", "income"),
			ruleB: NewTrue(),
		},
		{ruleA: NewEQFF("fred", "income"),
			ruleB: NewEQFF("bob", "cost"),
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"cost":   dlit.MustNew(18),
		"band":   dlit.NewString("alpha"),
	}
	for _, c := range cases {
		r := NewAnd(c.ruleA, c.ruleB)
		wantErr := InvalidRuleError{Rule: r}
		_, err := r.IsTrue(record)
		if err != wantErr {
			t.Errorf("IsTrue(record) rule: %s, err: %v, want: %v", r, err, wantErr)
		}
	}
}
