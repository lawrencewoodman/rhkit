package rule

import (
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestAndString(t *testing.T) {
	ruleA := NewTrue()
	ruleB := NewEQFF("income", "cost")
	want := "true() && income == cost"
	r := NewAnd(ruleA, ruleB)
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestAndGetInNiParts(t *testing.T) {
	ruleA := NewTrue()
	ruleB := NewEQFF("income", "cost")
	r := NewAnd(ruleA, ruleB)
	a, b, c := r.GetInNiParts()
	if a || b != "" || c != "" {
		t.Errorf("GetInNiParts() got: %t,\"%s\",\"%s\" - want: %t,\"\",\"\"",
			a, b, c, false)
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
		ruleA   Rule
		ruleB   Rule
		wantErr error
	}{
		{ruleA: NewTrue(),
			ruleB:   NewEQFF("fred", "income"),
			wantErr: InvalidRuleError("true() && fred == income"),
		},
		{ruleA: NewEQFF("fred", "income"),
			ruleB:   NewTrue(),
			wantErr: InvalidRuleError("fred == income && true()"),
		},
		{ruleA: NewEQFF("fred", "income"),
			ruleB:   NewEQFF("bob", "cost"),
			wantErr: InvalidRuleError("fred == income && bob == cost"),
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"cost":   dlit.MustNew(18),
		"band":   dlit.NewString("alpha"),
	}
	for _, c := range cases {
		r := NewAnd(c.ruleA, c.ruleB)
		_, err := r.IsTrue(record)
		if err != c.wantErr {
			t.Errorf("IsTrue(record) rule: %s, err: %v, want: %v", r, err, c.wantErr)
		}
	}
}
