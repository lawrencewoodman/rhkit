package rulehunter

import (
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/rule"
	"regexp"
	"sort"
	"testing"
)

func TestGenerateRules_1(t *testing.T) {
	testPurpose := "Ensure generates correct rules for each field"
	inputDescription := &Description{
		map[string]*fieldDescription{
			"team": &fieldDescription{
				kind: ftString,
				values: []*dlit.Literal{
					dlit.MustNew("a"), dlit.MustNew("b"), dlit.MustNew("c"),
				},
			},
			"teamOut": &fieldDescription{
				kind: ftString,
				values: []*dlit.Literal{
					dlit.MustNew("a"), dlit.MustNew("c"), dlit.MustNew("d"),
					dlit.MustNew("e"), dlit.MustNew("f"),
				},
			},
			"teamBob": &fieldDescription{
				kind: ftString,
				values: []*dlit.Literal{
					dlit.MustNew("a"), dlit.MustNew("b"), dlit.MustNew("c"),
				},
			},
			"camp": &fieldDescription{
				kind: ftString,
				values: []*dlit.Literal{
					dlit.MustNew("arthur"), dlit.MustNew("offa"),
					dlit.MustNew("richard"), dlit.MustNew("owen"),
				},
			},
			"level": &fieldDescription{
				kind:  ftInt,
				min:   dlit.MustNew(0),
				max:   dlit.MustNew(5),
				maxDP: 0,
				values: []*dlit.Literal{
					dlit.MustNew(0), dlit.MustNew(1), dlit.MustNew(2),
					dlit.MustNew(3), dlit.MustNew(4), dlit.MustNew(5),
				},
			},
			"levelBob": &fieldDescription{
				kind:  ftInt,
				min:   dlit.MustNew(0),
				max:   dlit.MustNew(5),
				maxDP: 0,
				values: []*dlit.Literal{
					dlit.MustNew(0), dlit.MustNew(1), dlit.MustNew(2),
					dlit.MustNew(3), dlit.MustNew(4), dlit.MustNew(5),
				},
			},
			"flow": &fieldDescription{
				kind:  ftFloat,
				min:   dlit.MustNew(0),
				max:   dlit.MustNew(10.5),
				maxDP: 2,
				values: []*dlit.Literal{
					dlit.MustNew(0.0), dlit.MustNew(2.34), dlit.MustNew(10.5),
				},
			},
			"position": &fieldDescription{
				kind:  ftInt,
				min:   dlit.MustNew(1),
				max:   dlit.MustNew(13),
				maxDP: 0,
				values: []*dlit.Literal{
					dlit.MustNew(1), dlit.MustNew(2), dlit.MustNew(3),
					dlit.MustNew(4), dlit.MustNew(5), dlit.MustNew(6),
					dlit.MustNew(7), dlit.MustNew(8), dlit.MustNew(9),
					dlit.MustNew(10), dlit.MustNew(11), dlit.MustNew(12),
					dlit.MustNew(13),
				},
			},
		}}
	excludeFields := []string{"teamBob", "levelBob"}
	cases := []struct {
		field     string
		wantRules []rule.Rule
	}{
		{"team", []rule.Rule{
			rule.MustNewDRule("team == \"a\""),
			rule.MustNewDRule("team == \"b\""), rule.MustNewDRule("team == \"c\""),
			rule.MustNewDRule("team != \"a\""),
			rule.MustNewDRule("team != \"b\""),
			rule.MustNewDRule("team != \"c\""),
			rule.MustNewDRule("team == teamOut"),
			rule.MustNewDRule("team != teamOut"),
		}},
		{"teamOut", []rule.Rule{
			rule.MustNewDRule("teamOut == \"a\""),
			rule.MustNewDRule("teamOut == \"c\""),
			rule.MustNewDRule("teamOut == \"d\""),
			rule.MustNewDRule("teamOut == \"e\""),
			rule.MustNewDRule("teamOut == \"f\""),
			rule.MustNewDRule("teamOut != \"a\""),
			rule.MustNewDRule("teamOut != \"c\""),
			rule.MustNewDRule("teamOut != \"d\""),
			rule.MustNewDRule("teamOut != \"e\""),
			rule.MustNewDRule("teamOut != \"f\""),
			rule.MustNewDRule("in(teamOut,\"a\",\"c\")"),
			rule.MustNewDRule("in(teamOut,\"a\",\"d\")"),
			rule.MustNewDRule("in(teamOut,\"a\",\"e\")"),
			rule.MustNewDRule("in(teamOut,\"a\",\"f\")"),
			rule.MustNewDRule("in(teamOut,\"c\",\"d\")"),
			rule.MustNewDRule("in(teamOut,\"c\",\"e\")"),
			rule.MustNewDRule("in(teamOut,\"c\",\"f\")"),
			rule.MustNewDRule("in(teamOut,\"d\",\"e\")"),
			rule.MustNewDRule("in(teamOut,\"d\",\"f\")"),
			rule.MustNewDRule("in(teamOut,\"e\",\"f\")"),
			rule.MustNewDRule("in(teamOut,\"a\",\"c\",\"d\")"),
			rule.MustNewDRule("in(teamOut,\"a\",\"c\",\"e\")"),
			rule.MustNewDRule("in(teamOut,\"a\",\"c\",\"f\")"),
			rule.MustNewDRule("in(teamOut,\"a\",\"d\",\"e\")"),
			rule.MustNewDRule("in(teamOut,\"a\",\"d\",\"f\")"),
			rule.MustNewDRule("in(teamOut,\"a\",\"e\",\"f\")"),
			rule.MustNewDRule("in(teamOut,\"c\",\"d\",\"e\")"),
			rule.MustNewDRule("in(teamOut,\"c\",\"d\",\"f\")"),
			rule.MustNewDRule("in(teamOut,\"c\",\"e\",\"f\")"),
			rule.MustNewDRule("in(teamOut,\"d\",\"e\",\"f\")"),
			rule.MustNewDRule("ni(teamOut,\"a\",\"c\")"),
			rule.MustNewDRule("ni(teamOut,\"a\",\"d\")"),
			rule.MustNewDRule("ni(teamOut,\"a\",\"e\")"),
			rule.MustNewDRule("ni(teamOut,\"a\",\"f\")"),
			rule.MustNewDRule("ni(teamOut,\"c\",\"d\")"),
			rule.MustNewDRule("ni(teamOut,\"c\",\"e\")"),
			rule.MustNewDRule("ni(teamOut,\"c\",\"f\")"),
			rule.MustNewDRule("ni(teamOut,\"d\",\"e\")"),
			rule.MustNewDRule("ni(teamOut,\"d\",\"f\")"),
			rule.MustNewDRule("ni(teamOut,\"e\",\"f\")"),
			rule.MustNewDRule("ni(teamOut,\"a\",\"c\",\"d\")"),
			rule.MustNewDRule("ni(teamOut,\"a\",\"c\",\"e\")"),
			rule.MustNewDRule("ni(teamOut,\"a\",\"c\",\"f\")"),
			rule.MustNewDRule("ni(teamOut,\"a\",\"d\",\"e\")"),
			rule.MustNewDRule("ni(teamOut,\"a\",\"d\",\"f\")"),
			rule.MustNewDRule("ni(teamOut,\"a\",\"e\",\"f\")"),
			rule.MustNewDRule("ni(teamOut,\"c\",\"d\",\"e\")"),
			rule.MustNewDRule("ni(teamOut,\"c\",\"d\",\"f\")"),
			rule.MustNewDRule("ni(teamOut,\"c\",\"e\",\"f\")"),
			rule.MustNewDRule("ni(teamOut,\"d\",\"e\",\"f\")"),
		}},
		{"level", []rule.Rule{
			rule.MustNewDRule("level == 0"),
			rule.MustNewDRule("level == 1"),
			rule.MustNewDRule("level == 2"),
			rule.MustNewDRule("level == 3"),
			rule.MustNewDRule("level == 4"),
			rule.MustNewDRule("level == 5"),
			rule.MustNewDRule("level != 0"),
			rule.MustNewDRule("level != 1"),
			rule.MustNewDRule("level != 2"),
			rule.MustNewDRule("level != 3"),
			rule.MustNewDRule("level != 4"),
			rule.MustNewDRule("level != 5"),
			rule.MustNewDRule("level < position"),
			rule.MustNewDRule("level <= position"),
			rule.MustNewDRule("level != position"),
			rule.MustNewDRule("level >= position"),
			rule.MustNewDRule("level > position"),
			rule.MustNewDRule("level == position"),
			rule.MustNewDRule("level >= 0"),
			rule.MustNewDRule("level >= 1"),
			rule.MustNewDRule("level >= 2"),
			rule.MustNewDRule("level >= 3"),
			rule.MustNewDRule("level >= 4"),
			rule.MustNewDRule("level <= 1"),
			rule.MustNewDRule("level <= 2"),
			rule.MustNewDRule("level <= 3"),
			rule.MustNewDRule("level <= 4"),
			rule.MustNewDRule("level <= 5"),
			rule.MustNewDRule("in(level,\"0\",\"1\")"),
			rule.MustNewDRule("in(level,\"0\",\"2\")"),
			rule.MustNewDRule("in(level,\"0\",\"3\")"),
			rule.MustNewDRule("in(level,\"0\",\"4\")"),
			rule.MustNewDRule("in(level,\"0\",\"5\")"),
			rule.MustNewDRule("in(level,\"1\",\"2\")"),
			rule.MustNewDRule("in(level,\"1\",\"3\")"),
			rule.MustNewDRule("in(level,\"1\",\"4\")"),
			rule.MustNewDRule("in(level,\"1\",\"5\")"),
			rule.MustNewDRule("in(level,\"2\",\"3\")"),
			rule.MustNewDRule("in(level,\"2\",\"4\")"),
			rule.MustNewDRule("in(level,\"2\",\"5\")"),
			rule.MustNewDRule("in(level,\"3\",\"4\")"),
			rule.MustNewDRule("in(level,\"3\",\"5\")"),
			rule.MustNewDRule("in(level,\"4\",\"5\")"),
			rule.MustNewDRule("in(level,\"0\",\"1\",\"2\")"),
			rule.MustNewDRule("in(level,\"0\",\"1\",\"3\")"),
			rule.MustNewDRule("in(level,\"0\",\"1\",\"4\")"),
			rule.MustNewDRule("in(level,\"0\",\"1\",\"5\")"),
			rule.MustNewDRule("in(level,\"0\",\"2\",\"3\")"),
			rule.MustNewDRule("in(level,\"0\",\"2\",\"4\")"),
			rule.MustNewDRule("in(level,\"0\",\"2\",\"5\")"),
			rule.MustNewDRule("in(level,\"0\",\"3\",\"4\")"),
			rule.MustNewDRule("in(level,\"0\",\"3\",\"5\")"),
			rule.MustNewDRule("in(level,\"0\",\"4\",\"5\")"),
			rule.MustNewDRule("in(level,\"1\",\"2\",\"3\")"),
			rule.MustNewDRule("in(level,\"1\",\"2\",\"4\")"),
			rule.MustNewDRule("in(level,\"1\",\"2\",\"5\")"),
			rule.MustNewDRule("in(level,\"1\",\"3\",\"4\")"),
			rule.MustNewDRule("in(level,\"1\",\"3\",\"5\")"),
			rule.MustNewDRule("in(level,\"1\",\"4\",\"5\")"),
			rule.MustNewDRule("in(level,\"2\",\"3\",\"4\")"),
			rule.MustNewDRule("in(level,\"2\",\"3\",\"5\")"),
			rule.MustNewDRule("in(level,\"2\",\"4\",\"5\")"),
			rule.MustNewDRule("in(level,\"3\",\"4\",\"5\")"),
			rule.MustNewDRule("in(level,\"0\",\"1\",\"2\",\"3\")"),
			rule.MustNewDRule("in(level,\"0\",\"1\",\"2\",\"4\")"),
			rule.MustNewDRule("in(level,\"0\",\"1\",\"2\",\"5\")"),
			rule.MustNewDRule("in(level,\"0\",\"1\",\"3\",\"4\")"),
			rule.MustNewDRule("in(level,\"0\",\"1\",\"3\",\"5\")"),
			rule.MustNewDRule("in(level,\"0\",\"1\",\"4\",\"5\")"),
			rule.MustNewDRule("in(level,\"0\",\"2\",\"3\",\"4\")"),
			rule.MustNewDRule("in(level,\"0\",\"2\",\"3\",\"5\")"),
			rule.MustNewDRule("in(level,\"0\",\"2\",\"4\",\"5\")"),
			rule.MustNewDRule("in(level,\"0\",\"3\",\"4\",\"5\")"),
			rule.MustNewDRule("in(level,\"1\",\"2\",\"3\",\"4\")"),
			rule.MustNewDRule("in(level,\"1\",\"2\",\"3\",\"5\")"),
			rule.MustNewDRule("in(level,\"1\",\"2\",\"4\",\"5\")"),
			rule.MustNewDRule("in(level,\"1\",\"3\",\"4\",\"5\")"),
			rule.MustNewDRule("in(level,\"2\",\"3\",\"4\",\"5\")"),
			rule.MustNewDRule("ni(level,\"0\",\"1\")"),
			rule.MustNewDRule("ni(level,\"0\",\"2\")"),
			rule.MustNewDRule("ni(level,\"0\",\"3\")"),
			rule.MustNewDRule("ni(level,\"0\",\"4\")"),
			rule.MustNewDRule("ni(level,\"0\",\"5\")"),
			rule.MustNewDRule("ni(level,\"1\",\"2\")"),
			rule.MustNewDRule("ni(level,\"1\",\"3\")"),
			rule.MustNewDRule("ni(level,\"1\",\"4\")"),
			rule.MustNewDRule("ni(level,\"1\",\"5\")"),
			rule.MustNewDRule("ni(level,\"2\",\"3\")"),
			rule.MustNewDRule("ni(level,\"2\",\"4\")"),
			rule.MustNewDRule("ni(level,\"2\",\"5\")"),
			rule.MustNewDRule("ni(level,\"3\",\"4\")"),
			rule.MustNewDRule("ni(level,\"3\",\"5\")"),
			rule.MustNewDRule("ni(level,\"4\",\"5\")"),
			rule.MustNewDRule("ni(level,\"0\",\"1\",\"2\")"),
			rule.MustNewDRule("ni(level,\"0\",\"1\",\"3\")"),
			rule.MustNewDRule("ni(level,\"0\",\"1\",\"4\")"),
			rule.MustNewDRule("ni(level,\"0\",\"1\",\"5\")"),
			rule.MustNewDRule("ni(level,\"0\",\"2\",\"3\")"),
			rule.MustNewDRule("ni(level,\"0\",\"2\",\"4\")"),
			rule.MustNewDRule("ni(level,\"0\",\"2\",\"5\")"),
			rule.MustNewDRule("ni(level,\"0\",\"3\",\"4\")"),
			rule.MustNewDRule("ni(level,\"0\",\"3\",\"5\")"),
			rule.MustNewDRule("ni(level,\"0\",\"4\",\"5\")"),
			rule.MustNewDRule("ni(level,\"1\",\"2\",\"3\")"),
			rule.MustNewDRule("ni(level,\"1\",\"2\",\"4\")"),
			rule.MustNewDRule("ni(level,\"1\",\"2\",\"5\")"),
			rule.MustNewDRule("ni(level,\"1\",\"3\",\"4\")"),
			rule.MustNewDRule("ni(level,\"1\",\"3\",\"5\")"),
			rule.MustNewDRule("ni(level,\"1\",\"4\",\"5\")"),
			rule.MustNewDRule("ni(level,\"2\",\"3\",\"4\")"),
			rule.MustNewDRule("ni(level,\"2\",\"3\",\"5\")"),
			rule.MustNewDRule("ni(level,\"2\",\"4\",\"5\")"),
			rule.MustNewDRule("ni(level,\"3\",\"4\",\"5\")"),
			rule.MustNewDRule("ni(level,\"0\",\"1\",\"2\",\"3\")"),
			rule.MustNewDRule("ni(level,\"0\",\"1\",\"2\",\"4\")"),
			rule.MustNewDRule("ni(level,\"0\",\"1\",\"2\",\"5\")"),
			rule.MustNewDRule("ni(level,\"0\",\"1\",\"3\",\"4\")"),
			rule.MustNewDRule("ni(level,\"0\",\"1\",\"3\",\"5\")"),
			rule.MustNewDRule("ni(level,\"0\",\"1\",\"4\",\"5\")"),
			rule.MustNewDRule("ni(level,\"0\",\"2\",\"3\",\"4\")"),
			rule.MustNewDRule("ni(level,\"0\",\"2\",\"3\",\"5\")"),
			rule.MustNewDRule("ni(level,\"0\",\"2\",\"4\",\"5\")"),
			rule.MustNewDRule("ni(level,\"0\",\"3\",\"4\",\"5\")"),
			rule.MustNewDRule("ni(level,\"1\",\"2\",\"3\",\"4\")"),
			rule.MustNewDRule("ni(level,\"1\",\"2\",\"3\",\"5\")"),
			rule.MustNewDRule("ni(level,\"1\",\"2\",\"4\",\"5\")"),
			rule.MustNewDRule("ni(level,\"1\",\"3\",\"4\",\"5\")"),
			rule.MustNewDRule("ni(level,\"2\",\"3\",\"4\",\"5\")"),
		}},
		{"flow", []rule.Rule{
			rule.MustNewDRule("flow == 0"),
			rule.MustNewDRule("flow == 2.34"),
			rule.MustNewDRule("flow == 10.5"),
			rule.MustNewDRule("flow != 0"),
			rule.MustNewDRule("flow != 2.34"),
			rule.MustNewDRule("flow != 10.5"),
			rule.MustNewDRule("flow < level"),
			rule.MustNewDRule("flow <= level"),
			rule.MustNewDRule("flow != level"),
			rule.MustNewDRule("flow >= level"),
			rule.MustNewDRule("flow > level"),
			rule.MustNewDRule("flow == level"),
			rule.MustNewDRule("flow < position"),
			rule.MustNewDRule("flow <= position"),
			rule.MustNewDRule("flow != position"),
			rule.MustNewDRule("flow >= position"),
			rule.MustNewDRule("flow > position"),
			rule.MustNewDRule("flow == position"),
			rule.MustNewDRule("flow >= 0"),
			rule.MustNewDRule("flow >= 1.05"),
			rule.MustNewDRule("flow >= 2.1"),
			rule.MustNewDRule("flow >= 3.15"),
			rule.MustNewDRule("flow >= 4.2"),
			rule.MustNewDRule("flow >= 5.25"),
			rule.MustNewDRule("flow >= 6.3"),
			rule.MustNewDRule("flow >= 7.35"),
			rule.MustNewDRule("flow >= 8.4"),
			rule.MustNewDRule("flow >= 9.45"),
			rule.MustNewDRule("flow <= 1.05"),
			rule.MustNewDRule("flow <= 2.1"),
			rule.MustNewDRule("flow <= 3.15"),
			rule.MustNewDRule("flow <= 4.2"),
			rule.MustNewDRule("flow <= 5.25"),
			rule.MustNewDRule("flow <= 6.3"),
			rule.MustNewDRule("flow <= 7.35"),
			rule.MustNewDRule("flow <= 8.4"),
			rule.MustNewDRule("flow <= 9.45"),
		}},
		{"position", []rule.Rule{
			rule.MustNewDRule("position == 1"),
			rule.MustNewDRule("position == 2"),
			rule.MustNewDRule("position == 3"),
			rule.MustNewDRule("position == 4"),
			rule.MustNewDRule("position == 5"),
			rule.MustNewDRule("position == 6"),
			rule.MustNewDRule("position == 7"),
			rule.MustNewDRule("position == 8"),
			rule.MustNewDRule("position == 9"),
			rule.MustNewDRule("position == 10"),
			rule.MustNewDRule("position == 11"),
			rule.MustNewDRule("position == 12"),
			rule.MustNewDRule("position == 13"),
			rule.MustNewDRule("position != 1"),
			rule.MustNewDRule("position != 2"),
			rule.MustNewDRule("position != 3"),
			rule.MustNewDRule("position != 4"),
			rule.MustNewDRule("position != 5"),
			rule.MustNewDRule("position != 6"),
			rule.MustNewDRule("position != 7"),
			rule.MustNewDRule("position != 8"),
			rule.MustNewDRule("position != 9"),
			rule.MustNewDRule("position != 10"),
			rule.MustNewDRule("position != 11"),
			rule.MustNewDRule("position != 12"),
			rule.MustNewDRule("position != 13"),
			rule.MustNewDRule("position >= 1"),
			rule.MustNewDRule("position >= 2"),
			rule.MustNewDRule("position >= 3"),
			rule.MustNewDRule("position >= 4"),
			rule.MustNewDRule("position >= 5"),
			rule.MustNewDRule("position >= 6"),
			rule.MustNewDRule("position >= 7"),
			rule.MustNewDRule("position >= 8"),
			rule.MustNewDRule("position >= 9"),
			rule.MustNewDRule("position >= 10"),
			rule.MustNewDRule("position >= 11"),
			rule.MustNewDRule("position >= 12"),
			rule.MustNewDRule("position <= 2"),
			rule.MustNewDRule("position <= 3"),
			rule.MustNewDRule("position <= 4"),
			rule.MustNewDRule("position <= 5"),
			rule.MustNewDRule("position <= 6"),
			rule.MustNewDRule("position <= 7"),
			rule.MustNewDRule("position <= 8"),
			rule.MustNewDRule("position <= 9"),
			rule.MustNewDRule("position <= 10"),
			rule.MustNewDRule("position <= 11"),
			rule.MustNewDRule("position <= 12"),
			rule.MustNewDRule("position <= 13"),
		}},
	}

	rules, err := GenerateRules(inputDescription, excludeFields)
	if err != nil {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("GenerateRules(%q, %q) err: %q",
			inputDescription, excludeFields, err)
	}

	for _, c := range cases {
		gotFieldRules := getFieldRules(c.field, rules)
		rulesMatch, msg := matchRulesUnordered(gotFieldRules, c.wantRules)
		if !rulesMatch {
			gotFieldRuleStrs := rulesToSortedStrings(gotFieldRules)
			wantRuleStrs := rulesToSortedStrings(c.wantRules)
			t.Errorf("Test: %s\n", testPurpose)
			t.Errorf("matchRulesUnordered() rules don't match for field: %s - %s\ngot: %s\nwant: %s\n",
				c.field, msg, gotFieldRuleStrs, wantRuleStrs)
		}
	}
}

func TestGenerateRules_2(t *testing.T) {
	testPurpose := "Ensure generates a True rule"
	inputDescription := &Description{
		map[string]*fieldDescription{
			"team": &fieldDescription{
				kind: ftString,
				values: []*dlit.Literal{
					dlit.MustNew("a"), dlit.MustNew("b"), dlit.MustNew("c"),
				},
			},
			"teamOut": &fieldDescription{
				kind: ftString,
				values: []*dlit.Literal{
					dlit.MustNew("a"), dlit.MustNew("c"), dlit.MustNew("d"),
					dlit.MustNew("e"), dlit.MustNew("f"),
				},
			},
		}}
	excludeFields := []string{}

	rules, err := GenerateRules(inputDescription, excludeFields)
	if err != nil {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("GenerateRules(%q, %q) err: %q",
			inputDescription, excludeFields, err)
	}

	trueRuleFound := false
	for _, r := range rules {
		if _, isTrueRule := r.(rule.True); isTrueRule {
			trueRuleFound = true
			break
		}
	}
	if !trueRuleFound {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("GenerateRules(%q, %q)  - True rule missing",
			inputDescription, excludeFields)
	}
}

func TestCombinedRules(t *testing.T) {
	cases := []struct {
		inRules   []rule.Rule
		wantRules []rule.Rule
	}{
		{[]rule.Rule{
			rule.MustNewDRule("team == \"a\""),
			rule.MustNewDRule("band > 4"),
			rule.MustNewDRule("in(team,\"red\",\"green\",\"blue\")"),
		},
			[]rule.Rule{
				rule.MustNewDRule("team == \"a\" && band > 4"),
				rule.MustNewDRule("team == \"a\" || band > 4"),
				rule.MustNewDRule("team == \"a\" && in(team,\"red\",\"green\",\"blue\")"),
				rule.MustNewDRule("team == \"a\" || in(team,\"red\",\"green\",\"blue\")"),
				rule.MustNewDRule("band > 4 && in(team,\"red\",\"green\",\"blue\")"),
				rule.MustNewDRule("band > 4 || in(team,\"red\",\"green\",\"blue\")"),
			}},
		{[]rule.Rule{rule.MustNewDRule("team == \"a\"")}, []rule.Rule{}},
		{[]rule.Rule{}, []rule.Rule{}},
	}

	for _, c := range cases {
		gotRules := CombineRules(c.inRules)
		rulesMatch, msg := matchRulesUnordered(gotRules, c.wantRules)
		if !rulesMatch {
			gotRuleStrs := rulesToSortedStrings(gotRules)
			wantRuleStrs := rulesToSortedStrings(c.wantRules)
			t.Errorf("matchRulesUnordered() rules don't match: %s\ngot: %s\nwant: %s\n",
				msg, gotRuleStrs, wantRuleStrs)
		}
	}
}

/*************************************
 *    Helper Functions
 *************************************/
var matchFieldInNiRegexp = regexp.MustCompile("^((in\\(|ni\\()+)([^ ,]+)(.*)$")
var matchFieldMatchRegexp = regexp.MustCompile("^([^ (]+)( .*)$")

func getFieldRules(
	field string,
	rules []rule.Rule,
) []rule.Rule {
	fieldRules := make([]rule.Rule, 0)
	for _, rule := range rules {
		ruleStr := rule.String()
		ruleField := matchFieldMatchRegexp.ReplaceAllString(ruleStr, "$1")
		ruleField = matchFieldInNiRegexp.ReplaceAllString(ruleField, "$3")
		if field == ruleField {
			fieldRules = append(fieldRules, rule)
		}
	}
	return fieldRules
}

func rulesToSortedStrings(rules []rule.Rule) []string {
	r := make([]string, len(rules))
	for i, rule := range rules {
		r[i] = rule.String()
	}
	sort.Strings(r)
	return r
}

func matchRulesUnordered(
	rules1 []rule.Rule,
	rules2 []rule.Rule,
) (bool, string) {
	if len(rules1) != len(rules2) {
		return false, "rules different lengths"
	}
	for _, rule1 := range rules1 {
		found := false
		for _, rule2 := range rules2 {
			if rule1.String() == rule2.String() {
				found = true
				break
			}
		}
		if !found {
			return false, fmt.Sprintf("rule doesn't exist: %s", rule1)
		}
	}
	return true, ""
}
