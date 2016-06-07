package rulehunter

import (
	"errors"
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestNewRule_errors(t *testing.T) {
	cases := []struct {
		ruleString string
		wantError  error
	}{
		{"7 {} 3", ErrInvalidRule("Invalid rule: 7 {} 3")},
	}
	for _, c := range cases {
		_, err := newRule(c.ruleString)
		if err == nil {
			t.Errorf("newRule(%s) no error, expected: %s", c.ruleString, c.wantError)
			return
		}
		if err.Error() != c.wantError.Error() {
			t.Errorf("newRule(%s) got error: %s, want error: %s",
				c.ruleString, err, c.wantError)
		}
	}
}

func TestIsTrue(t *testing.T) {
	cases := []struct {
		rule       *Rule
		wantIsTrue bool
	}{
		{mustNewRule("band > 3"), true},
		{mustNewRule("band == 2"), false},
	}
	record := map[string]*dlit.Literal{
		"cost": dlit.MustNew(4.5),
		"band": dlit.MustNew(4),
	}
	for _, c := range cases {
		gotIsTrue, err := c.rule.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(%s) rule: %s err: %s", record, c.rule, err)
		}
		if gotIsTrue != c.wantIsTrue {
			t.Errorf("IsTrue(%s) got: %t want: %t", record, gotIsTrue, c.wantIsTrue)
		}
	}
}

func TestIsTrue_errors(t *testing.T) {
	cases := []struct {
		rule      *Rule
		wantError error
	}{
		{mustNewRule("band > 3"),
			dexpr.ErrInvalidExpr("Variable doesn't exist: band")},
	}
	record := map[string]*dlit.Literal{
		"cost":   dlit.MustNew(4.5),
		"length": dlit.MustNew(4),
	}
	for _, c := range cases {
		_, err := c.rule.IsTrue(record)
		if err == nil {
			t.Errorf("IsTrue(%s) no error, expected: %s", record, c.wantError)
		}
		if err.Error() != c.wantError.Error() {
			t.Errorf("IsTrue(%s) got error: %s, want error: %s", record,
				err, c.wantError)
		}
	}
}

func TestString(t *testing.T) {
	cases := []struct {
		rule *Rule
		want string
	}{
		{mustNewRule("band > 3"), "band > 3"},
		{mustNewRule("in(Band, \"a\", \"bb\")"), "in(Band, \"a\", \"bb\")"},
	}
	for _, c := range cases {
		got := c.rule.String()
		if got != c.want {
			t.Errorf("c.rule.String() got: %s want: %s", got, c.want)
		}
	}
}

func TestGetTweakableParts(t *testing.T) {
	cases := []struct {
		rule            *Rule
		wantIsTweakable bool
		wantFieldName   string
		wantOperator    string
		wantValue       string
	}{
		{mustNewRule("band > 3"), true, "band", ">", "3"},
		{mustNewRule("band == 2"), false, "", "", ""},
		{mustNewRule("in(band, \"a\", \"b\")"), false, "", "", ""},
	}
	for _, c := range cases {
		gotIsTweakable, gotFieldName, gotOperator, gotValue :=
			c.rule.GetTweakableParts()
		if gotIsTweakable != c.wantIsTweakable {
			t.Errorf("GetTweakableParts() rule: %s, got isTweakable: %t want: %t",
				c.rule, gotIsTweakable, c.wantIsTweakable)
		}
		if gotFieldName != c.wantFieldName {
			t.Errorf("GetTweakableParts() rule: %s, got fieldName: %s want: %s",
				c.rule, gotFieldName, c.wantFieldName)
		}
		if gotOperator != c.wantOperator {
			t.Errorf("GetTweakableParts() rule: %s, got operator: %s want: %s",
				c.rule, gotOperator, c.wantOperator)
		}
		if gotValue != c.wantValue {
			t.Errorf("GetTweakableParts() rule: %s, got value: %s want: %s",
				c.rule, gotValue, c.wantValue)
		}
	}
}

func TestGetInNiParts(t *testing.T) {
	cases := []struct {
		rule          *Rule
		wantIsInNi    bool
		wantOperator  string
		wantFieldName string
	}{
		{mustNewRule("band > 3"), false, "", ""},
		{mustNewRule("band == 2"), false, "", ""},
		{mustNewRule("in(band, \"a\", \"b\")"), true, "in", "band"},
		{mustNewRule("in(flow, \"4\", \"6\")"), true, "in", "flow"},
		{mustNewRule("ni(band, \"a\", \"b\")"), true, "ni", "band"},
		{mustNewRule("ni(flow, \"4\", \"6\")"), true, "ni", "flow"},
	}
	for _, c := range cases {
		gotIsInNi, gotOperator, gotFieldName := c.rule.GetInNiParts()
		if gotIsInNi != c.wantIsInNi {
			t.Errorf("GetInNIParts() rule: %s, got isInNi: %t want: %t",
				c.rule, gotIsInNi, c.wantIsInNi)
		}
		if gotOperator != c.wantOperator {
			t.Errorf("GetInNIParts() rule: %s, got operator: %s want: %s",
				c.rule, gotOperator, c.wantOperator)
		}
		if gotFieldName != c.wantFieldName {
			t.Errorf("GetInNIParts() rule: %s, got fieldName: %s want: %s",
				c.rule, gotFieldName, c.wantFieldName)
		}
	}
}

func TestCloneWithValue(t *testing.T) {
	cases := []struct {
		rule     *Rule
		newValue string
		wantRule *Rule
	}{
		{mustNewRule("band > 3"), "20", mustNewRule("band > 20")},
	}
	for _, c := range cases {
		gotRule, err := c.rule.CloneWithValue(c.newValue)
		if err != nil {
			t.Errorf("CloneWithValue(%s) rule: %s, err: %s", c.newValue, c.rule, err)
		}
		if gotRule.String() != c.wantRule.String() {
			t.Errorf("CloneWithValue(%s) rule: %s, got: %s, want: %s",
				c.newValue, c.rule, gotRule, c.wantRule)
		}
	}
}

func TestCloneWithValue_errors(t *testing.T) {
	cases := []struct {
		rule      *Rule
		newValue  string
		wantError error
	}{
		{mustNewRule("band > 3 && band < 9"), "20",
			errors.New("Can't clone non-tweakable rule: band > 3 && band < 9")},
	}
	for _, c := range cases {
		_, err := c.rule.CloneWithValue(c.newValue)
		if err == nil {
			t.Errorf("CloneWithValue(%s) rule: %s, no error, expected: %s",
				c.newValue, c.rule, c.wantError)
		}
		if err.Error() != c.wantError.Error() {
			t.Errorf("CloneWithValue(%s) rule: %s, got error: %s, want: %s",
				c.newValue, c.rule, err, c.wantError)
		}
	}
}
