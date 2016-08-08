package rule

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestNewDRule_errors(t *testing.T) {
	cases := []struct {
		ruleString string
		wantError  error
	}{
		{"7 {} 3", InvalidRuleError("7 {} 3")},
	}
	for _, c := range cases {
		_, err := NewDRule(c.ruleString)
		if err == nil {
			t.Errorf("newDRule(%s) no error, expected: %s", c.ruleString, c.wantError)
			return
		}
		if err.Error() != c.wantError.Error() {
			t.Errorf("newDRule(%s) got error: %s, want error: %s",
				c.ruleString, err, c.wantError)
		}
	}
}

func TestIsTrue(t *testing.T) {
	cases := []struct {
		rule       Rule
		wantIsTrue bool
	}{
		{MustNewDRule("band > 3"), true},
		{MustNewDRule("band == 2"), false},
	}
	record := map[string]*dlit.Literal{
		"cost": dlit.MustNew(4.5),
		"band": dlit.MustNew(4),
	}
	for _, c := range cases {
		gotIsTrue, err := c.rule.IsTrue(record)
		if err != nil {
			t.Errorf("isTrue(%s) rule: %s err: %s", record, c.rule, err)
		}
		if gotIsTrue != c.wantIsTrue {
			t.Errorf("isTrue(%s) got: %t want: %t", record, gotIsTrue, c.wantIsTrue)
		}
	}
}

func TestIsTrue_errors(t *testing.T) {
	cases := []struct {
		rule      Rule
		wantError error
	}{
		{MustNewDRule("band > 3"),
			dexpr.ErrInvalidExpr{
				Expr: "band > 3",
				Err:  dexpr.ErrVarNotExist("band"),
			},
		},
	}
	record := map[string]*dlit.Literal{
		"cost":   dlit.MustNew(4.5),
		"length": dlit.MustNew(4),
	}
	for _, c := range cases {
		_, err := c.rule.IsTrue(record)
		if err == nil {
			t.Errorf("isTrue(%s) no error, expected: %s", record, c.wantError)
		}
		if err.Error() != c.wantError.Error() {
			t.Errorf("isTrue(%s) got error: %s, want error: %s", record,
				err, c.wantError)
		}
	}
}

func TestString(t *testing.T) {
	cases := []struct {
		rule Rule
		want string
	}{
		{MustNewDRule("band > 3"), "band > 3"},
		{MustNewDRule("in(Band, \"a\", \"bb\")"), "in(Band, \"a\", \"bb\")"},
	}
	for _, c := range cases {
		got := c.rule.String()
		if got != c.want {
			t.Errorf("c.rule.String() got: %s want: %s", got, c.want)
		}
	}
}

func TestGetInNiParts(t *testing.T) {
	cases := []struct {
		rule          Rule
		wantIsInNi    bool
		wantOperator  string
		wantFieldName string
	}{
		{MustNewDRule("band > 3"), false, "", ""},
		{MustNewDRule("band == 2"), false, "", ""},
		{MustNewDRule("in(band, \"a\", \"b\")"), true, "in", "band"},
		{MustNewDRule("in(flow, \"4\", \"6\")"), true, "in", "flow"},
		{MustNewDRule("ni(band, \"a\", \"b\")"), true, "ni", "band"},
		{MustNewDRule("ni(flow, \"4\", \"6\")"), true, "ni", "flow"},
	}
	for _, c := range cases {
		gotIsInNi, gotOperator, gotFieldName := c.rule.GetInNiParts()
		if gotIsInNi != c.wantIsInNi {
			t.Errorf("getInNIParts() rule: %s, got isInNi: %t want: %t",
				c.rule, gotIsInNi, c.wantIsInNi)
		}
		if gotOperator != c.wantOperator {
			t.Errorf("getInNIParts() rule: %s, got operator: %s want: %s",
				c.rule, gotOperator, c.wantOperator)
		}
		if gotFieldName != c.wantFieldName {
			t.Errorf("getInNIParts() rule: %s, got fieldName: %s want: %s",
				c.rule, gotFieldName, c.wantFieldName)
		}
	}
}

/**************************
 *  Benchmarks
 **************************/

func BenchmarkDRuleIsTrue(b *testing.B) {
	record := map[string]*dlit.Literal{
		"band":   dlit.MustNew(23),
		"income": dlit.MustNew(1024),
		"cost":   dlit.MustNew(890),
	}
	r := MustNewDRule("cost < income")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = r.IsTrue(record)
	}
}
