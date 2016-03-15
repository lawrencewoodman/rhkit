package main

import (
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/dlit_go"
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
		_, err := NewRule(c.ruleString)
		if err == nil {
			t.Errorf("NewRule(%s) no error, expected: %s", c.ruleString, c.wantError)
			return
		}
		if err.Error() != c.wantError.Error() {
			t.Errorf("NewRule(%s) got error: %s, want error: %s",
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
			t.Errorf("IsTrue(%s) got: %s want: %s", record, gotIsTrue, c.wantIsTrue)
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
			t.Errorf("IsTrue(%s) got: %s want: %s", got, c.want)
		}
	}
}
