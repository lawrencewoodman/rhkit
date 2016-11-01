package rule

import (
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestOrString(t *testing.T) {
	ruleA := NewTrue()
	ruleB := NewEQFF("income", "cost")
	want := "true() || income == cost"
	r := NewOr(ruleA, ruleB)
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestOrIsTrue(t *testing.T) {
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
			want:  true,
		},
		{ruleA: NewTrue(),
			ruleB: NewNEFF("income", "income"),
			want:  true,
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
		r := NewOr(c.ruleA, c.ruleB)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestOrIsTrue_errors(t *testing.T) {
	cases := []struct {
		ruleA   Rule
		ruleB   Rule
		wantErr error
	}{
		{ruleA: NewTrue(),
			ruleB: NewEQFF("fred", "income"),
			wantErr: InvalidRuleError{
				Rule: NewOr(NewTrue(), NewEQFF("fred", "income")),
			},
		},
		{ruleA: NewEQFF("fred", "income"),
			ruleB: NewTrue(),
			wantErr: InvalidRuleError{
				Rule: NewOr(NewEQFF("fred", "income"), NewTrue()),
			},
		},
		{ruleA: NewEQFF("fred", "income"),
			ruleB: NewEQFF("bob", "cost"),
			wantErr: InvalidRuleError{
				Rule: NewOr(NewEQFF("fred", "income"), NewEQFF("bob", "cost")),
			},
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"cost":   dlit.MustNew(18),
		"band":   dlit.NewString("alpha"),
	}
	for _, c := range cases {
		r := NewOr(c.ruleA, c.ruleB)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}
