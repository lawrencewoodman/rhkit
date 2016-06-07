package rulehunter

import (
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/internal"
	"regexp"
	"sort"
	"testing"
)

func TestGenerateRules_1(t *testing.T) {
	testPurpose := "Ensure generates correct rules for each field"
	inputDescription := &Description{
		map[string]*fieldDescription{
			"team": &fieldDescription{
				kind: internal.STRING,
				values: []*dlit.Literal{
					dlit.MustNew("a"), dlit.MustNew("b"), dlit.MustNew("c"),
				},
			},
			"teamOut": &fieldDescription{
				kind: internal.STRING,
				values: []*dlit.Literal{
					dlit.MustNew("a"), dlit.MustNew("c"), dlit.MustNew("d"),
					dlit.MustNew("e"), dlit.MustNew("f"),
				},
			},
			"teamBob": &fieldDescription{
				kind: internal.STRING,
				values: []*dlit.Literal{
					dlit.MustNew("a"), dlit.MustNew("b"), dlit.MustNew("c"),
				},
			},
			"camp": &fieldDescription{
				kind: internal.STRING,
				values: []*dlit.Literal{
					dlit.MustNew("arthur"), dlit.MustNew("offa"),
					dlit.MustNew("richard"), dlit.MustNew("owen"),
				},
			},
			"level": &fieldDescription{
				kind:  internal.INT,
				min:   dlit.MustNew(0),
				max:   dlit.MustNew(5),
				maxDP: 0,
				values: []*dlit.Literal{
					dlit.MustNew(0), dlit.MustNew(1), dlit.MustNew(2),
					dlit.MustNew(3), dlit.MustNew(4), dlit.MustNew(5),
				},
			},
			"levelBob": &fieldDescription{
				kind:  internal.INT,
				min:   dlit.MustNew(0),
				max:   dlit.MustNew(5),
				maxDP: 0,
				values: []*dlit.Literal{
					dlit.MustNew(0), dlit.MustNew(1), dlit.MustNew(2),
					dlit.MustNew(3), dlit.MustNew(4), dlit.MustNew(5),
				},
			},
			"flow": &fieldDescription{
				kind:  internal.FLOAT,
				min:   dlit.MustNew(0),
				max:   dlit.MustNew(10.5),
				maxDP: 2,
				values: []*dlit.Literal{
					dlit.MustNew(0.0), dlit.MustNew(2.34), dlit.MustNew(10.5),
				},
			},
			"position": &fieldDescription{
				kind:  internal.INT,
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
		wantRules []*Rule
	}{
		{"team", []*Rule{
			mustNewRule("team == \"a\""),
			mustNewRule("team == \"b\""), mustNewRule("team == \"c\""),
			mustNewRule("team != \"a\""),
			mustNewRule("team != \"b\""),
			mustNewRule("team != \"c\""),
			mustNewRule("team == teamOut"),
			mustNewRule("team != teamOut"),
		}},
		{"teamOut", []*Rule{
			mustNewRule("teamOut == \"a\""),
			mustNewRule("teamOut == \"c\""),
			mustNewRule("teamOut == \"d\""),
			mustNewRule("teamOut == \"e\""),
			mustNewRule("teamOut == \"f\""),
			mustNewRule("teamOut != \"a\""),
			mustNewRule("teamOut != \"c\""),
			mustNewRule("teamOut != \"d\""),
			mustNewRule("teamOut != \"e\""),
			mustNewRule("teamOut != \"f\""),
			mustNewRule("in(teamOut,\"a\",\"c\")"),
			mustNewRule("in(teamOut,\"a\",\"d\")"),
			mustNewRule("in(teamOut,\"a\",\"e\")"),
			mustNewRule("in(teamOut,\"a\",\"f\")"),
			mustNewRule("in(teamOut,\"c\",\"d\")"),
			mustNewRule("in(teamOut,\"c\",\"e\")"),
			mustNewRule("in(teamOut,\"c\",\"f\")"),
			mustNewRule("in(teamOut,\"d\",\"e\")"),
			mustNewRule("in(teamOut,\"d\",\"f\")"),
			mustNewRule("in(teamOut,\"e\",\"f\")"),
			mustNewRule("in(teamOut,\"a\",\"c\",\"d\")"),
			mustNewRule("in(teamOut,\"a\",\"c\",\"e\")"),
			mustNewRule("in(teamOut,\"a\",\"c\",\"f\")"),
			mustNewRule("in(teamOut,\"a\",\"d\",\"e\")"),
			mustNewRule("in(teamOut,\"a\",\"d\",\"f\")"),
			mustNewRule("in(teamOut,\"a\",\"e\",\"f\")"),
			mustNewRule("in(teamOut,\"c\",\"d\",\"e\")"),
			mustNewRule("in(teamOut,\"c\",\"d\",\"f\")"),
			mustNewRule("in(teamOut,\"c\",\"e\",\"f\")"),
			mustNewRule("in(teamOut,\"d\",\"e\",\"f\")"),
			mustNewRule("ni(teamOut,\"a\",\"c\")"),
			mustNewRule("ni(teamOut,\"a\",\"d\")"),
			mustNewRule("ni(teamOut,\"a\",\"e\")"),
			mustNewRule("ni(teamOut,\"a\",\"f\")"),
			mustNewRule("ni(teamOut,\"c\",\"d\")"),
			mustNewRule("ni(teamOut,\"c\",\"e\")"),
			mustNewRule("ni(teamOut,\"c\",\"f\")"),
			mustNewRule("ni(teamOut,\"d\",\"e\")"),
			mustNewRule("ni(teamOut,\"d\",\"f\")"),
			mustNewRule("ni(teamOut,\"e\",\"f\")"),
			mustNewRule("ni(teamOut,\"a\",\"c\",\"d\")"),
			mustNewRule("ni(teamOut,\"a\",\"c\",\"e\")"),
			mustNewRule("ni(teamOut,\"a\",\"c\",\"f\")"),
			mustNewRule("ni(teamOut,\"a\",\"d\",\"e\")"),
			mustNewRule("ni(teamOut,\"a\",\"d\",\"f\")"),
			mustNewRule("ni(teamOut,\"a\",\"e\",\"f\")"),
			mustNewRule("ni(teamOut,\"c\",\"d\",\"e\")"),
			mustNewRule("ni(teamOut,\"c\",\"d\",\"f\")"),
			mustNewRule("ni(teamOut,\"c\",\"e\",\"f\")"),
			mustNewRule("ni(teamOut,\"d\",\"e\",\"f\")"),
		}},
		{"level", []*Rule{
			mustNewRule("level == 0"),
			mustNewRule("level == 1"),
			mustNewRule("level == 2"),
			mustNewRule("level == 3"),
			mustNewRule("level == 4"),
			mustNewRule("level == 5"),
			mustNewRule("level != 0"),
			mustNewRule("level != 1"),
			mustNewRule("level != 2"),
			mustNewRule("level != 3"),
			mustNewRule("level != 4"),
			mustNewRule("level != 5"),
			mustNewRule("level < position"),
			mustNewRule("level <= position"),
			mustNewRule("level != position"),
			mustNewRule("level >= position"),
			mustNewRule("level > position"),
			mustNewRule("level == position"),
			mustNewRule("level >= 0"),
			mustNewRule("level >= 1"),
			mustNewRule("level >= 2"),
			mustNewRule("level >= 3"),
			mustNewRule("level >= 4"),
			mustNewRule("level <= 1"),
			mustNewRule("level <= 2"),
			mustNewRule("level <= 3"),
			mustNewRule("level <= 4"),
			mustNewRule("level <= 5"),
			mustNewRule("in(level,\"0\",\"1\")"),
			mustNewRule("in(level,\"0\",\"2\")"),
			mustNewRule("in(level,\"0\",\"3\")"),
			mustNewRule("in(level,\"0\",\"4\")"),
			mustNewRule("in(level,\"0\",\"5\")"),
			mustNewRule("in(level,\"1\",\"2\")"),
			mustNewRule("in(level,\"1\",\"3\")"),
			mustNewRule("in(level,\"1\",\"4\")"),
			mustNewRule("in(level,\"1\",\"5\")"),
			mustNewRule("in(level,\"2\",\"3\")"),
			mustNewRule("in(level,\"2\",\"4\")"),
			mustNewRule("in(level,\"2\",\"5\")"),
			mustNewRule("in(level,\"3\",\"4\")"),
			mustNewRule("in(level,\"3\",\"5\")"),
			mustNewRule("in(level,\"4\",\"5\")"),
			mustNewRule("in(level,\"0\",\"1\",\"2\")"),
			mustNewRule("in(level,\"0\",\"1\",\"3\")"),
			mustNewRule("in(level,\"0\",\"1\",\"4\")"),
			mustNewRule("in(level,\"0\",\"1\",\"5\")"),
			mustNewRule("in(level,\"0\",\"2\",\"3\")"),
			mustNewRule("in(level,\"0\",\"2\",\"4\")"),
			mustNewRule("in(level,\"0\",\"2\",\"5\")"),
			mustNewRule("in(level,\"0\",\"3\",\"4\")"),
			mustNewRule("in(level,\"0\",\"3\",\"5\")"),
			mustNewRule("in(level,\"0\",\"4\",\"5\")"),
			mustNewRule("in(level,\"1\",\"2\",\"3\")"),
			mustNewRule("in(level,\"1\",\"2\",\"4\")"),
			mustNewRule("in(level,\"1\",\"2\",\"5\")"),
			mustNewRule("in(level,\"1\",\"3\",\"4\")"),
			mustNewRule("in(level,\"1\",\"3\",\"5\")"),
			mustNewRule("in(level,\"1\",\"4\",\"5\")"),
			mustNewRule("in(level,\"2\",\"3\",\"4\")"),
			mustNewRule("in(level,\"2\",\"3\",\"5\")"),
			mustNewRule("in(level,\"2\",\"4\",\"5\")"),
			mustNewRule("in(level,\"3\",\"4\",\"5\")"),
			mustNewRule("in(level,\"0\",\"1\",\"2\",\"3\")"),
			mustNewRule("in(level,\"0\",\"1\",\"2\",\"4\")"),
			mustNewRule("in(level,\"0\",\"1\",\"2\",\"5\")"),
			mustNewRule("in(level,\"0\",\"1\",\"3\",\"4\")"),
			mustNewRule("in(level,\"0\",\"1\",\"3\",\"5\")"),
			mustNewRule("in(level,\"0\",\"1\",\"4\",\"5\")"),
			mustNewRule("in(level,\"0\",\"2\",\"3\",\"4\")"),
			mustNewRule("in(level,\"0\",\"2\",\"3\",\"5\")"),
			mustNewRule("in(level,\"0\",\"2\",\"4\",\"5\")"),
			mustNewRule("in(level,\"0\",\"3\",\"4\",\"5\")"),
			mustNewRule("in(level,\"1\",\"2\",\"3\",\"4\")"),
			mustNewRule("in(level,\"1\",\"2\",\"3\",\"5\")"),
			mustNewRule("in(level,\"1\",\"2\",\"4\",\"5\")"),
			mustNewRule("in(level,\"1\",\"3\",\"4\",\"5\")"),
			mustNewRule("in(level,\"2\",\"3\",\"4\",\"5\")"),
			mustNewRule("ni(level,\"0\",\"1\")"),
			mustNewRule("ni(level,\"0\",\"2\")"),
			mustNewRule("ni(level,\"0\",\"3\")"),
			mustNewRule("ni(level,\"0\",\"4\")"),
			mustNewRule("ni(level,\"0\",\"5\")"),
			mustNewRule("ni(level,\"1\",\"2\")"),
			mustNewRule("ni(level,\"1\",\"3\")"),
			mustNewRule("ni(level,\"1\",\"4\")"),
			mustNewRule("ni(level,\"1\",\"5\")"),
			mustNewRule("ni(level,\"2\",\"3\")"),
			mustNewRule("ni(level,\"2\",\"4\")"),
			mustNewRule("ni(level,\"2\",\"5\")"),
			mustNewRule("ni(level,\"3\",\"4\")"),
			mustNewRule("ni(level,\"3\",\"5\")"),
			mustNewRule("ni(level,\"4\",\"5\")"),
			mustNewRule("ni(level,\"0\",\"1\",\"2\")"),
			mustNewRule("ni(level,\"0\",\"1\",\"3\")"),
			mustNewRule("ni(level,\"0\",\"1\",\"4\")"),
			mustNewRule("ni(level,\"0\",\"1\",\"5\")"),
			mustNewRule("ni(level,\"0\",\"2\",\"3\")"),
			mustNewRule("ni(level,\"0\",\"2\",\"4\")"),
			mustNewRule("ni(level,\"0\",\"2\",\"5\")"),
			mustNewRule("ni(level,\"0\",\"3\",\"4\")"),
			mustNewRule("ni(level,\"0\",\"3\",\"5\")"),
			mustNewRule("ni(level,\"0\",\"4\",\"5\")"),
			mustNewRule("ni(level,\"1\",\"2\",\"3\")"),
			mustNewRule("ni(level,\"1\",\"2\",\"4\")"),
			mustNewRule("ni(level,\"1\",\"2\",\"5\")"),
			mustNewRule("ni(level,\"1\",\"3\",\"4\")"),
			mustNewRule("ni(level,\"1\",\"3\",\"5\")"),
			mustNewRule("ni(level,\"1\",\"4\",\"5\")"),
			mustNewRule("ni(level,\"2\",\"3\",\"4\")"),
			mustNewRule("ni(level,\"2\",\"3\",\"5\")"),
			mustNewRule("ni(level,\"2\",\"4\",\"5\")"),
			mustNewRule("ni(level,\"3\",\"4\",\"5\")"),
			mustNewRule("ni(level,\"0\",\"1\",\"2\",\"3\")"),
			mustNewRule("ni(level,\"0\",\"1\",\"2\",\"4\")"),
			mustNewRule("ni(level,\"0\",\"1\",\"2\",\"5\")"),
			mustNewRule("ni(level,\"0\",\"1\",\"3\",\"4\")"),
			mustNewRule("ni(level,\"0\",\"1\",\"3\",\"5\")"),
			mustNewRule("ni(level,\"0\",\"1\",\"4\",\"5\")"),
			mustNewRule("ni(level,\"0\",\"2\",\"3\",\"4\")"),
			mustNewRule("ni(level,\"0\",\"2\",\"3\",\"5\")"),
			mustNewRule("ni(level,\"0\",\"2\",\"4\",\"5\")"),
			mustNewRule("ni(level,\"0\",\"3\",\"4\",\"5\")"),
			mustNewRule("ni(level,\"1\",\"2\",\"3\",\"4\")"),
			mustNewRule("ni(level,\"1\",\"2\",\"3\",\"5\")"),
			mustNewRule("ni(level,\"1\",\"2\",\"4\",\"5\")"),
			mustNewRule("ni(level,\"1\",\"3\",\"4\",\"5\")"),
			mustNewRule("ni(level,\"2\",\"3\",\"4\",\"5\")"),
		}},
		{"flow", []*Rule{
			mustNewRule("flow == 0"),
			mustNewRule("flow == 2.34"),
			mustNewRule("flow == 10.5"),
			mustNewRule("flow != 0"),
			mustNewRule("flow != 2.34"),
			mustNewRule("flow != 10.5"),
			mustNewRule("flow < level"),
			mustNewRule("flow <= level"),
			mustNewRule("flow != level"),
			mustNewRule("flow >= level"),
			mustNewRule("flow > level"),
			mustNewRule("flow == level"),
			mustNewRule("flow < position"),
			mustNewRule("flow <= position"),
			mustNewRule("flow != position"),
			mustNewRule("flow >= position"),
			mustNewRule("flow > position"),
			mustNewRule("flow == position"),
			mustNewRule("flow >= 0"),
			mustNewRule("flow >= 1.05"),
			mustNewRule("flow >= 2.1"),
			mustNewRule("flow >= 3.15"),
			mustNewRule("flow >= 4.2"),
			mustNewRule("flow >= 5.25"),
			mustNewRule("flow >= 6.3"),
			mustNewRule("flow >= 7.35"),
			mustNewRule("flow >= 8.4"),
			mustNewRule("flow >= 9.45"),
			mustNewRule("flow <= 1.05"),
			mustNewRule("flow <= 2.1"),
			mustNewRule("flow <= 3.15"),
			mustNewRule("flow <= 4.2"),
			mustNewRule("flow <= 5.25"),
			mustNewRule("flow <= 6.3"),
			mustNewRule("flow <= 7.35"),
			mustNewRule("flow <= 8.4"),
			mustNewRule("flow <= 9.45"),
		}},
		{"position", []*Rule{
			mustNewRule("position == 1"),
			mustNewRule("position == 2"),
			mustNewRule("position == 3"),
			mustNewRule("position == 4"),
			mustNewRule("position == 5"),
			mustNewRule("position == 6"),
			mustNewRule("position == 7"),
			mustNewRule("position == 8"),
			mustNewRule("position == 9"),
			mustNewRule("position == 10"),
			mustNewRule("position == 11"),
			mustNewRule("position == 12"),
			mustNewRule("position == 13"),
			mustNewRule("position != 1"),
			mustNewRule("position != 2"),
			mustNewRule("position != 3"),
			mustNewRule("position != 4"),
			mustNewRule("position != 5"),
			mustNewRule("position != 6"),
			mustNewRule("position != 7"),
			mustNewRule("position != 8"),
			mustNewRule("position != 9"),
			mustNewRule("position != 10"),
			mustNewRule("position != 11"),
			mustNewRule("position != 12"),
			mustNewRule("position != 13"),
			mustNewRule("position >= 1"),
			mustNewRule("position >= 2"),
			mustNewRule("position >= 3"),
			mustNewRule("position >= 4"),
			mustNewRule("position >= 5"),
			mustNewRule("position >= 6"),
			mustNewRule("position >= 7"),
			mustNewRule("position >= 8"),
			mustNewRule("position >= 9"),
			mustNewRule("position >= 10"),
			mustNewRule("position >= 11"),
			mustNewRule("position >= 12"),
			mustNewRule("position <= 2"),
			mustNewRule("position <= 3"),
			mustNewRule("position <= 4"),
			mustNewRule("position <= 5"),
			mustNewRule("position <= 6"),
			mustNewRule("position <= 7"),
			mustNewRule("position <= 8"),
			mustNewRule("position <= 9"),
			mustNewRule("position <= 10"),
			mustNewRule("position <= 11"),
			mustNewRule("position <= 12"),
			mustNewRule("position <= 13"),
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
	testPurpose := "Ensure generates a 'true()' rule"
	inputDescription := &Description{
		map[string]*fieldDescription{
			"team": &fieldDescription{
				kind: internal.STRING,
				values: []*dlit.Literal{
					dlit.MustNew("a"), dlit.MustNew("b"), dlit.MustNew("c"),
				},
			},
			"teamOut": &fieldDescription{
				kind: internal.STRING,
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
	for _, rule := range rules {
		if rule.String() == "true()" {
			trueRuleFound = true
			break
		}
	}
	if !trueRuleFound {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("GenerateRules(%q, %q)  - 'true()' rule missing",
			inputDescription, excludeFields)
	}
}

func TestCombinedRules(t *testing.T) {
	cases := []struct {
		inRules   []*Rule
		wantRules []*Rule
	}{
		{[]*Rule{
			mustNewRule("team == \"a\""),
			mustNewRule("band > 4"),
			mustNewRule("in(team,\"red\",\"green\",\"blue\")"),
		},
			[]*Rule{
				mustNewRule("team == \"a\" && band > 4"),
				mustNewRule("team == \"a\" || band > 4"),
				mustNewRule("team == \"a\" && in(team,\"red\",\"green\",\"blue\")"),
				mustNewRule("team == \"a\" || in(team,\"red\",\"green\",\"blue\")"),
				mustNewRule("band > 4 && in(team,\"red\",\"green\",\"blue\")"),
				mustNewRule("band > 4 || in(team,\"red\",\"green\",\"blue\")"),
			}},
		{[]*Rule{mustNewRule("team == \"a\"")}, []*Rule{}},
		{[]*Rule{}, []*Rule{}},
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
	rules []*Rule,
) []*Rule {
	fieldRules := make([]*Rule, 0)
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

func rulesToSortedStrings(rules []*Rule) []string {
	r := make([]string, len(rules))
	for i, rule := range rules {
		r[i] = rule.String()
	}
	sort.Strings(r)
	return r
}

func matchRulesUnordered(
	rules1 []*Rule,
	rules2 []*Rule,
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
