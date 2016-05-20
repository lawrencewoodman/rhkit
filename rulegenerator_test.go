package rulehunter

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/description"
	"github.com/vlifesystems/rulehunter/rule"
	"regexp"
	"sort"
	"testing"
)

func TestGenerateRules_1(t *testing.T) {
	testPurpose := "Ensure generates correct rules for each field"
	inputDescription := &description.Description{
		map[string]*description.Field{
			"team": &description.Field{
				Kind: description.STRING,
				Values: []*dlit.Literal{
					dlit.MustNew("a"), dlit.MustNew("b"), dlit.MustNew("c"),
				},
			},
			"teamOut": &description.Field{
				Kind: description.STRING,
				Values: []*dlit.Literal{
					dlit.MustNew("a"), dlit.MustNew("c"), dlit.MustNew("d"),
					dlit.MustNew("e"), dlit.MustNew("f"),
				},
			},
			"teamBob": &description.Field{
				Kind: description.STRING,
				Values: []*dlit.Literal{
					dlit.MustNew("a"), dlit.MustNew("b"), dlit.MustNew("c"),
				},
			},
			"camp": &description.Field{
				Kind: description.STRING,
				Values: []*dlit.Literal{
					dlit.MustNew("arthur"), dlit.MustNew("offa"),
					dlit.MustNew("richard"), dlit.MustNew("owen"),
				},
			},
			"level": &description.Field{
				Kind:  description.INT,
				Min:   dlit.MustNew(0),
				Max:   dlit.MustNew(5),
				MaxDP: 0,
				Values: []*dlit.Literal{
					dlit.MustNew(0), dlit.MustNew(1), dlit.MustNew(2),
					dlit.MustNew(3), dlit.MustNew(4), dlit.MustNew(5),
				},
			},
			"levelBob": &description.Field{
				Kind:  description.INT,
				Min:   dlit.MustNew(0),
				Max:   dlit.MustNew(5),
				MaxDP: 0,
				Values: []*dlit.Literal{
					dlit.MustNew(0), dlit.MustNew(1), dlit.MustNew(2),
					dlit.MustNew(3), dlit.MustNew(4), dlit.MustNew(5),
				},
			},
			"flow": &description.Field{
				Kind:  description.FLOAT,
				Min:   dlit.MustNew(0),
				Max:   dlit.MustNew(10.5),
				MaxDP: 2,
				Values: []*dlit.Literal{
					dlit.MustNew(0.0), dlit.MustNew(2.34), dlit.MustNew(10.5),
				},
			},
			"position": &description.Field{
				Kind:  description.INT,
				Min:   dlit.MustNew(1),
				Max:   dlit.MustNew(13),
				MaxDP: 0,
				Values: []*dlit.Literal{
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
		wantRules []*rule.Rule
	}{
		{"team", []*rule.Rule{
			rule.MustNew("team == \"a\""),
			rule.MustNew("team == \"b\""), rule.MustNew("team == \"c\""),
			rule.MustNew("team != \"a\""),
			rule.MustNew("team != \"b\""),
			rule.MustNew("team != \"c\""),
			rule.MustNew("team == teamOut"),
			rule.MustNew("team != teamOut"),
		}},
		{"teamOut", []*rule.Rule{
			rule.MustNew("teamOut == \"a\""),
			rule.MustNew("teamOut == \"c\""),
			rule.MustNew("teamOut == \"d\""),
			rule.MustNew("teamOut == \"e\""),
			rule.MustNew("teamOut == \"f\""),
			rule.MustNew("teamOut != \"a\""),
			rule.MustNew("teamOut != \"c\""),
			rule.MustNew("teamOut != \"d\""),
			rule.MustNew("teamOut != \"e\""),
			rule.MustNew("teamOut != \"f\""),
			rule.MustNew("in(teamOut,\"a\",\"c\")"),
			rule.MustNew("in(teamOut,\"a\",\"d\")"),
			rule.MustNew("in(teamOut,\"a\",\"e\")"),
			rule.MustNew("in(teamOut,\"a\",\"f\")"),
			rule.MustNew("in(teamOut,\"c\",\"d\")"),
			rule.MustNew("in(teamOut,\"c\",\"e\")"),
			rule.MustNew("in(teamOut,\"c\",\"f\")"),
			rule.MustNew("in(teamOut,\"d\",\"e\")"),
			rule.MustNew("in(teamOut,\"d\",\"f\")"),
			rule.MustNew("in(teamOut,\"e\",\"f\")"),
			rule.MustNew("in(teamOut,\"a\",\"c\",\"d\")"),
			rule.MustNew("in(teamOut,\"a\",\"c\",\"e\")"),
			rule.MustNew("in(teamOut,\"a\",\"c\",\"f\")"),
			rule.MustNew("in(teamOut,\"a\",\"d\",\"e\")"),
			rule.MustNew("in(teamOut,\"a\",\"d\",\"f\")"),
			rule.MustNew("in(teamOut,\"a\",\"e\",\"f\")"),
			rule.MustNew("in(teamOut,\"c\",\"d\",\"e\")"),
			rule.MustNew("in(teamOut,\"c\",\"d\",\"f\")"),
			rule.MustNew("in(teamOut,\"c\",\"e\",\"f\")"),
			rule.MustNew("in(teamOut,\"d\",\"e\",\"f\")"),
			rule.MustNew("ni(teamOut,\"a\",\"c\")"),
			rule.MustNew("ni(teamOut,\"a\",\"d\")"),
			rule.MustNew("ni(teamOut,\"a\",\"e\")"),
			rule.MustNew("ni(teamOut,\"a\",\"f\")"),
			rule.MustNew("ni(teamOut,\"c\",\"d\")"),
			rule.MustNew("ni(teamOut,\"c\",\"e\")"),
			rule.MustNew("ni(teamOut,\"c\",\"f\")"),
			rule.MustNew("ni(teamOut,\"d\",\"e\")"),
			rule.MustNew("ni(teamOut,\"d\",\"f\")"),
			rule.MustNew("ni(teamOut,\"e\",\"f\")"),
			rule.MustNew("ni(teamOut,\"a\",\"c\",\"d\")"),
			rule.MustNew("ni(teamOut,\"a\",\"c\",\"e\")"),
			rule.MustNew("ni(teamOut,\"a\",\"c\",\"f\")"),
			rule.MustNew("ni(teamOut,\"a\",\"d\",\"e\")"),
			rule.MustNew("ni(teamOut,\"a\",\"d\",\"f\")"),
			rule.MustNew("ni(teamOut,\"a\",\"e\",\"f\")"),
			rule.MustNew("ni(teamOut,\"c\",\"d\",\"e\")"),
			rule.MustNew("ni(teamOut,\"c\",\"d\",\"f\")"),
			rule.MustNew("ni(teamOut,\"c\",\"e\",\"f\")"),
			rule.MustNew("ni(teamOut,\"d\",\"e\",\"f\")"),
		}},
		{"level", []*rule.Rule{
			rule.MustNew("level == 0"),
			rule.MustNew("level == 1"),
			rule.MustNew("level == 2"),
			rule.MustNew("level == 3"),
			rule.MustNew("level == 4"),
			rule.MustNew("level == 5"),
			rule.MustNew("level != 0"),
			rule.MustNew("level != 1"),
			rule.MustNew("level != 2"),
			rule.MustNew("level != 3"),
			rule.MustNew("level != 4"),
			rule.MustNew("level != 5"),
			rule.MustNew("level < position"),
			rule.MustNew("level <= position"),
			rule.MustNew("level != position"),
			rule.MustNew("level >= position"),
			rule.MustNew("level > position"),
			rule.MustNew("level == position"),
			rule.MustNew("level >= 0"),
			rule.MustNew("level >= 1"),
			rule.MustNew("level >= 2"),
			rule.MustNew("level >= 3"),
			rule.MustNew("level >= 4"),
			rule.MustNew("level <= 1"),
			rule.MustNew("level <= 2"),
			rule.MustNew("level <= 3"),
			rule.MustNew("level <= 4"),
			rule.MustNew("level <= 5"),
			rule.MustNew("in(level,\"0\",\"1\")"),
			rule.MustNew("in(level,\"0\",\"2\")"),
			rule.MustNew("in(level,\"0\",\"3\")"),
			rule.MustNew("in(level,\"0\",\"4\")"),
			rule.MustNew("in(level,\"0\",\"5\")"),
			rule.MustNew("in(level,\"1\",\"2\")"),
			rule.MustNew("in(level,\"1\",\"3\")"),
			rule.MustNew("in(level,\"1\",\"4\")"),
			rule.MustNew("in(level,\"1\",\"5\")"),
			rule.MustNew("in(level,\"2\",\"3\")"),
			rule.MustNew("in(level,\"2\",\"4\")"),
			rule.MustNew("in(level,\"2\",\"5\")"),
			rule.MustNew("in(level,\"3\",\"4\")"),
			rule.MustNew("in(level,\"3\",\"5\")"),
			rule.MustNew("in(level,\"4\",\"5\")"),
			rule.MustNew("in(level,\"0\",\"1\",\"2\")"),
			rule.MustNew("in(level,\"0\",\"1\",\"3\")"),
			rule.MustNew("in(level,\"0\",\"1\",\"4\")"),
			rule.MustNew("in(level,\"0\",\"1\",\"5\")"),
			rule.MustNew("in(level,\"0\",\"2\",\"3\")"),
			rule.MustNew("in(level,\"0\",\"2\",\"4\")"),
			rule.MustNew("in(level,\"0\",\"2\",\"5\")"),
			rule.MustNew("in(level,\"0\",\"3\",\"4\")"),
			rule.MustNew("in(level,\"0\",\"3\",\"5\")"),
			rule.MustNew("in(level,\"0\",\"4\",\"5\")"),
			rule.MustNew("in(level,\"1\",\"2\",\"3\")"),
			rule.MustNew("in(level,\"1\",\"2\",\"4\")"),
			rule.MustNew("in(level,\"1\",\"2\",\"5\")"),
			rule.MustNew("in(level,\"1\",\"3\",\"4\")"),
			rule.MustNew("in(level,\"1\",\"3\",\"5\")"),
			rule.MustNew("in(level,\"1\",\"4\",\"5\")"),
			rule.MustNew("in(level,\"2\",\"3\",\"4\")"),
			rule.MustNew("in(level,\"2\",\"3\",\"5\")"),
			rule.MustNew("in(level,\"2\",\"4\",\"5\")"),
			rule.MustNew("in(level,\"3\",\"4\",\"5\")"),
			rule.MustNew("in(level,\"0\",\"1\",\"2\",\"3\")"),
			rule.MustNew("in(level,\"0\",\"1\",\"2\",\"4\")"),
			rule.MustNew("in(level,\"0\",\"1\",\"2\",\"5\")"),
			rule.MustNew("in(level,\"0\",\"1\",\"3\",\"4\")"),
			rule.MustNew("in(level,\"0\",\"1\",\"3\",\"5\")"),
			rule.MustNew("in(level,\"0\",\"1\",\"4\",\"5\")"),
			rule.MustNew("in(level,\"0\",\"2\",\"3\",\"4\")"),
			rule.MustNew("in(level,\"0\",\"2\",\"3\",\"5\")"),
			rule.MustNew("in(level,\"0\",\"2\",\"4\",\"5\")"),
			rule.MustNew("in(level,\"0\",\"3\",\"4\",\"5\")"),
			rule.MustNew("in(level,\"1\",\"2\",\"3\",\"4\")"),
			rule.MustNew("in(level,\"1\",\"2\",\"3\",\"5\")"),
			rule.MustNew("in(level,\"1\",\"2\",\"4\",\"5\")"),
			rule.MustNew("in(level,\"1\",\"3\",\"4\",\"5\")"),
			rule.MustNew("in(level,\"2\",\"3\",\"4\",\"5\")"),
			rule.MustNew("ni(level,\"0\",\"1\")"),
			rule.MustNew("ni(level,\"0\",\"2\")"),
			rule.MustNew("ni(level,\"0\",\"3\")"),
			rule.MustNew("ni(level,\"0\",\"4\")"),
			rule.MustNew("ni(level,\"0\",\"5\")"),
			rule.MustNew("ni(level,\"1\",\"2\")"),
			rule.MustNew("ni(level,\"1\",\"3\")"),
			rule.MustNew("ni(level,\"1\",\"4\")"),
			rule.MustNew("ni(level,\"1\",\"5\")"),
			rule.MustNew("ni(level,\"2\",\"3\")"),
			rule.MustNew("ni(level,\"2\",\"4\")"),
			rule.MustNew("ni(level,\"2\",\"5\")"),
			rule.MustNew("ni(level,\"3\",\"4\")"),
			rule.MustNew("ni(level,\"3\",\"5\")"),
			rule.MustNew("ni(level,\"4\",\"5\")"),
			rule.MustNew("ni(level,\"0\",\"1\",\"2\")"),
			rule.MustNew("ni(level,\"0\",\"1\",\"3\")"),
			rule.MustNew("ni(level,\"0\",\"1\",\"4\")"),
			rule.MustNew("ni(level,\"0\",\"1\",\"5\")"),
			rule.MustNew("ni(level,\"0\",\"2\",\"3\")"),
			rule.MustNew("ni(level,\"0\",\"2\",\"4\")"),
			rule.MustNew("ni(level,\"0\",\"2\",\"5\")"),
			rule.MustNew("ni(level,\"0\",\"3\",\"4\")"),
			rule.MustNew("ni(level,\"0\",\"3\",\"5\")"),
			rule.MustNew("ni(level,\"0\",\"4\",\"5\")"),
			rule.MustNew("ni(level,\"1\",\"2\",\"3\")"),
			rule.MustNew("ni(level,\"1\",\"2\",\"4\")"),
			rule.MustNew("ni(level,\"1\",\"2\",\"5\")"),
			rule.MustNew("ni(level,\"1\",\"3\",\"4\")"),
			rule.MustNew("ni(level,\"1\",\"3\",\"5\")"),
			rule.MustNew("ni(level,\"1\",\"4\",\"5\")"),
			rule.MustNew("ni(level,\"2\",\"3\",\"4\")"),
			rule.MustNew("ni(level,\"2\",\"3\",\"5\")"),
			rule.MustNew("ni(level,\"2\",\"4\",\"5\")"),
			rule.MustNew("ni(level,\"3\",\"4\",\"5\")"),
			rule.MustNew("ni(level,\"0\",\"1\",\"2\",\"3\")"),
			rule.MustNew("ni(level,\"0\",\"1\",\"2\",\"4\")"),
			rule.MustNew("ni(level,\"0\",\"1\",\"2\",\"5\")"),
			rule.MustNew("ni(level,\"0\",\"1\",\"3\",\"4\")"),
			rule.MustNew("ni(level,\"0\",\"1\",\"3\",\"5\")"),
			rule.MustNew("ni(level,\"0\",\"1\",\"4\",\"5\")"),
			rule.MustNew("ni(level,\"0\",\"2\",\"3\",\"4\")"),
			rule.MustNew("ni(level,\"0\",\"2\",\"3\",\"5\")"),
			rule.MustNew("ni(level,\"0\",\"2\",\"4\",\"5\")"),
			rule.MustNew("ni(level,\"0\",\"3\",\"4\",\"5\")"),
			rule.MustNew("ni(level,\"1\",\"2\",\"3\",\"4\")"),
			rule.MustNew("ni(level,\"1\",\"2\",\"3\",\"5\")"),
			rule.MustNew("ni(level,\"1\",\"2\",\"4\",\"5\")"),
			rule.MustNew("ni(level,\"1\",\"3\",\"4\",\"5\")"),
			rule.MustNew("ni(level,\"2\",\"3\",\"4\",\"5\")"),
		}},
		{"flow", []*rule.Rule{
			rule.MustNew("flow == 0"),
			rule.MustNew("flow == 2.34"),
			rule.MustNew("flow == 10.5"),
			rule.MustNew("flow != 0"),
			rule.MustNew("flow != 2.34"),
			rule.MustNew("flow != 10.5"),
			rule.MustNew("flow < level"),
			rule.MustNew("flow <= level"),
			rule.MustNew("flow != level"),
			rule.MustNew("flow >= level"),
			rule.MustNew("flow > level"),
			rule.MustNew("flow == level"),
			rule.MustNew("flow < position"),
			rule.MustNew("flow <= position"),
			rule.MustNew("flow != position"),
			rule.MustNew("flow >= position"),
			rule.MustNew("flow > position"),
			rule.MustNew("flow == position"),
			rule.MustNew("flow >= 0"),
			rule.MustNew("flow >= 1.05"),
			rule.MustNew("flow >= 2.1"),
			rule.MustNew("flow >= 3.15"),
			rule.MustNew("flow >= 4.2"),
			rule.MustNew("flow >= 5.25"),
			rule.MustNew("flow >= 6.3"),
			rule.MustNew("flow >= 7.35"),
			rule.MustNew("flow >= 8.4"),
			rule.MustNew("flow >= 9.45"),
			rule.MustNew("flow <= 1.05"),
			rule.MustNew("flow <= 2.1"),
			rule.MustNew("flow <= 3.15"),
			rule.MustNew("flow <= 4.2"),
			rule.MustNew("flow <= 5.25"),
			rule.MustNew("flow <= 6.3"),
			rule.MustNew("flow <= 7.35"),
			rule.MustNew("flow <= 8.4"),
			rule.MustNew("flow <= 9.45"),
		}},
		{"position", []*rule.Rule{
			rule.MustNew("position == 1"),
			rule.MustNew("position == 2"),
			rule.MustNew("position == 3"),
			rule.MustNew("position == 4"),
			rule.MustNew("position == 5"),
			rule.MustNew("position == 6"),
			rule.MustNew("position == 7"),
			rule.MustNew("position == 8"),
			rule.MustNew("position == 9"),
			rule.MustNew("position == 10"),
			rule.MustNew("position == 11"),
			rule.MustNew("position == 12"),
			rule.MustNew("position == 13"),
			rule.MustNew("position != 1"),
			rule.MustNew("position != 2"),
			rule.MustNew("position != 3"),
			rule.MustNew("position != 4"),
			rule.MustNew("position != 5"),
			rule.MustNew("position != 6"),
			rule.MustNew("position != 7"),
			rule.MustNew("position != 8"),
			rule.MustNew("position != 9"),
			rule.MustNew("position != 10"),
			rule.MustNew("position != 11"),
			rule.MustNew("position != 12"),
			rule.MustNew("position != 13"),
			rule.MustNew("position >= 1"),
			rule.MustNew("position >= 2"),
			rule.MustNew("position >= 3"),
			rule.MustNew("position >= 4"),
			rule.MustNew("position >= 5"),
			rule.MustNew("position >= 6"),
			rule.MustNew("position >= 7"),
			rule.MustNew("position >= 8"),
			rule.MustNew("position >= 9"),
			rule.MustNew("position >= 10"),
			rule.MustNew("position >= 11"),
			rule.MustNew("position >= 12"),
			rule.MustNew("position <= 2"),
			rule.MustNew("position <= 3"),
			rule.MustNew("position <= 4"),
			rule.MustNew("position <= 5"),
			rule.MustNew("position <= 6"),
			rule.MustNew("position <= 7"),
			rule.MustNew("position <= 8"),
			rule.MustNew("position <= 9"),
			rule.MustNew("position <= 10"),
			rule.MustNew("position <= 11"),
			rule.MustNew("position <= 12"),
			rule.MustNew("position <= 13"),
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
		rulesMatch, msg := matchRules(gotFieldRules, c.wantRules)
		if !rulesMatch {
			gotFieldRuleStrs := rulesToSortedStrings(gotFieldRules)
			wantRuleStrs := rulesToSortedStrings(c.wantRules)
			t.Errorf("Test: %s\n", testPurpose)
			t.Errorf("matchRules() rules don't match for field: %s - %s\ngot: %s\nwant: %s\n",
				c.field, msg, gotFieldRuleStrs, wantRuleStrs)
		}
	}
}

func TestGenerateRules_2(t *testing.T) {
	testPurpose := "Ensure generates a 'true()' rule"
	inputDescription := &description.Description{
		map[string]*description.Field{
			"team": &description.Field{
				Kind: description.STRING,
				Values: []*dlit.Literal{
					dlit.MustNew("a"), dlit.MustNew("b"), dlit.MustNew("c"),
				},
			},
			"teamOut": &description.Field{
				Kind: description.STRING,
				Values: []*dlit.Literal{
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
		inRules   []*rule.Rule
		wantRules []*rule.Rule
	}{
		{[]*rule.Rule{
			rule.MustNew("team == \"a\""),
			rule.MustNew("band > 4"),
			rule.MustNew("in(team,\"red\",\"green\",\"blue\")"),
		},
			[]*rule.Rule{
				rule.MustNew("team == \"a\" && band > 4"),
				rule.MustNew("team == \"a\" || band > 4"),
				rule.MustNew("team == \"a\" && in(team,\"red\",\"green\",\"blue\")"),
				rule.MustNew("team == \"a\" || in(team,\"red\",\"green\",\"blue\")"),
				rule.MustNew("band > 4 && in(team,\"red\",\"green\",\"blue\")"),
				rule.MustNew("band > 4 || in(team,\"red\",\"green\",\"blue\")"),
			}},
		{[]*rule.Rule{rule.MustNew("team == \"a\"")}, []*rule.Rule{}},
		{[]*rule.Rule{}, []*rule.Rule{}},
	}

	for _, c := range cases {
		gotRules := CombineRules(c.inRules)
		rulesMatch, msg := matchRules(gotRules, c.wantRules)
		if !rulesMatch {
			gotRuleStrs := rulesToSortedStrings(gotRules)
			wantRuleStrs := rulesToSortedStrings(c.wantRules)
			t.Errorf("matchRules() rules don't match: %s\ngot: %s\nwant: %s\n",
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
	rules []*rule.Rule,
) []*rule.Rule {
	fieldRules := make([]*rule.Rule, 0)
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

func rulesToSortedStrings(rules []*rule.Rule) []string {
	r := make([]string, len(rules))
	for i, rule := range rules {
		r[i] = rule.String()
	}
	sort.Strings(r)
	return r
}
