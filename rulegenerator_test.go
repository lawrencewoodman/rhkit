package rulehunter

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/rule"
	"regexp"
	"sort"
	"testing"
)

func TestGenerateRules_1(t *testing.T) {
	testPurpose := "Ensure generates correct rules for each field"
	fieldDescriptions := map[string]*FieldDescription{
		"team": &FieldDescription{
			Kind: STRING,
			Values: []*dlit.Literal{
				dlit.MustNew("a"), dlit.MustNew("b"), dlit.MustNew("c"),
			},
		},
		"teamOut": &FieldDescription{
			Kind: STRING,
			Values: []*dlit.Literal{
				dlit.MustNew("a"), dlit.MustNew("c"), dlit.MustNew("d"),
				dlit.MustNew("e"), dlit.MustNew("f"),
			},
		},
		"teamBob": &FieldDescription{
			Kind: STRING,
			Values: []*dlit.Literal{
				dlit.MustNew("a"), dlit.MustNew("b"), dlit.MustNew("c"),
			},
		},
		"camp": &FieldDescription{
			Kind: STRING,
			Values: []*dlit.Literal{
				dlit.MustNew("arthur"), dlit.MustNew("offa"),
				dlit.MustNew("richard"), dlit.MustNew("owen"),
			},
		},
		"level": &FieldDescription{
			Kind:  INT,
			Min:   dlit.MustNew(0),
			Max:   dlit.MustNew(5),
			MaxDP: 0,
			Values: []*dlit.Literal{
				dlit.MustNew(0), dlit.MustNew(1), dlit.MustNew(2),
				dlit.MustNew(3), dlit.MustNew(4), dlit.MustNew(5),
			},
		},
		"levelBob": &FieldDescription{
			Kind:  INT,
			Min:   dlit.MustNew(0),
			Max:   dlit.MustNew(5),
			MaxDP: 0,
			Values: []*dlit.Literal{
				dlit.MustNew(0), dlit.MustNew(1), dlit.MustNew(2),
				dlit.MustNew(3), dlit.MustNew(4), dlit.MustNew(5),
			},
		},
		"flow": &FieldDescription{
			Kind:  FLOAT,
			Min:   dlit.MustNew(0),
			Max:   dlit.MustNew(10.5),
			MaxDP: 2,
			Values: []*dlit.Literal{
				dlit.MustNew(0.0), dlit.MustNew(2.34), dlit.MustNew(10.5),
			},
		},
		"position": &FieldDescription{
			Kind:  INT,
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
	}
	excludeFields := []string{"teamBob", "levelBob"}
	cases := []struct {
		field     string
		wantRules []string
	}{
		{"team", []string{
			"team == \"a\"", "team == \"b\"", "team == \"c\"",
			"team != \"a\"", "team != \"b\"", "team != \"c\"",
			"team == teamOut", "team != teamOut",
		}},
		{"teamOut", []string{
			"teamOut == \"a\"", "teamOut == \"c\"",
			"teamOut == \"d\"", "teamOut == \"e\"", "teamOut == \"f\"",
			"teamOut != \"a\"", "teamOut != \"c\"",
			"teamOut != \"d\"", "teamOut != \"e\"", "teamOut != \"f\"",
			"in(teamOut,\"a\",\"c\")", "in(teamOut,\"a\",\"d\")",
			"in(teamOut,\"a\",\"e\")", "in(teamOut,\"a\",\"f\")",
			"in(teamOut,\"c\",\"d\")", "in(teamOut,\"c\",\"e\")",
			"in(teamOut,\"c\",\"f\")",
			"in(teamOut,\"d\",\"e\")", "in(teamOut,\"d\",\"f\")",
			"in(teamOut,\"e\",\"f\")",
			"in(teamOut,\"a\",\"c\",\"d\")", "in(teamOut,\"a\",\"c\",\"e\")",
			"in(teamOut,\"a\",\"c\",\"f\")",
			"in(teamOut,\"a\",\"d\",\"e\")", "in(teamOut,\"a\",\"d\",\"f\")",
			"in(teamOut,\"a\",\"e\",\"f\")",
			"in(teamOut,\"c\",\"d\",\"e\")", "in(teamOut,\"c\",\"d\",\"f\")",
			"in(teamOut,\"c\",\"e\",\"f\")",
			"in(teamOut,\"d\",\"e\",\"f\")",
			"ni(teamOut,\"a\",\"c\")", "ni(teamOut,\"a\",\"d\")",
			"ni(teamOut,\"a\",\"e\")", "ni(teamOut,\"a\",\"f\")",
			"ni(teamOut,\"c\",\"d\")", "ni(teamOut,\"c\",\"e\")",
			"ni(teamOut,\"c\",\"f\")",
			"ni(teamOut,\"d\",\"e\")", "ni(teamOut,\"d\",\"f\")",
			"ni(teamOut,\"e\",\"f\")",
			"ni(teamOut,\"a\",\"c\",\"d\")", "ni(teamOut,\"a\",\"c\",\"e\")",
			"ni(teamOut,\"a\",\"c\",\"f\")",
			"ni(teamOut,\"a\",\"d\",\"e\")", "ni(teamOut,\"a\",\"d\",\"f\")",
			"ni(teamOut,\"a\",\"e\",\"f\")",
			"ni(teamOut,\"c\",\"d\",\"e\")", "ni(teamOut,\"c\",\"d\",\"f\")",
			"ni(teamOut,\"c\",\"e\",\"f\")",
			"ni(teamOut,\"d\",\"e\",\"f\")",
		}},
		{"level", []string{
			"level == 0", "level == 1", "level == 2",
			"level == 3", "level == 4", "level == 5",
			"level != 0", "level != 1", "level != 2",
			"level != 3", "level != 4", "level != 5",
			"level < position", "level <= position", "level != position",
			"level >= position", "level > position", "level == position",
			"level >= 0", "level >= 1", "level >= 2",
			"level >= 3", "level >= 4",
			"level <= 1", "level <= 2", "level <= 3",
			"level <= 4", "level <= 5",
			"in(level,\"0\",\"1\")", "in(level,\"0\",\"2\")", "in(level,\"0\",\"3\")",
			"in(level,\"0\",\"4\")", "in(level,\"0\",\"5\")",
			"in(level,\"1\",\"2\")", "in(level,\"1\",\"3\")", "in(level,\"1\",\"4\")",
			"in(level,\"1\",\"5\")",
			"in(level,\"2\",\"3\")", "in(level,\"2\",\"4\")", "in(level,\"2\",\"5\")",
			"in(level,\"3\",\"4\")", "in(level,\"3\",\"5\")",
			"in(level,\"4\",\"5\")",
			"in(level,\"0\",\"1\",\"2\")", "in(level,\"0\",\"1\",\"3\")",
			"in(level,\"0\",\"1\",\"4\")", "in(level,\"0\",\"1\",\"5\")",
			"in(level,\"0\",\"2\",\"3\")", "in(level,\"0\",\"2\",\"4\")",
			"in(level,\"0\",\"2\",\"5\")", "in(level,\"0\",\"3\",\"4\")",
			"in(level,\"0\",\"3\",\"5\")", "in(level,\"0\",\"4\",\"5\")",
			"in(level,\"1\",\"2\",\"3\")", "in(level,\"1\",\"2\",\"4\")",
			"in(level,\"1\",\"2\",\"5\")", "in(level,\"1\",\"3\",\"4\")",
			"in(level,\"1\",\"3\",\"5\")", "in(level,\"1\",\"4\",\"5\")",
			"in(level,\"2\",\"3\",\"4\")", "in(level,\"2\",\"3\",\"5\")",
			"in(level,\"2\",\"4\",\"5\")",
			"in(level,\"3\",\"4\",\"5\")",
			"in(level,\"0\",\"1\",\"2\",\"3\")", "in(level,\"0\",\"1\",\"2\",\"4\")",
			"in(level,\"0\",\"1\",\"2\",\"5\")", "in(level,\"0\",\"1\",\"3\",\"4\")",
			"in(level,\"0\",\"1\",\"3\",\"5\")", "in(level,\"0\",\"1\",\"4\",\"5\")",
			"in(level,\"0\",\"2\",\"3\",\"4\")", "in(level,\"0\",\"2\",\"3\",\"5\")",
			"in(level,\"0\",\"2\",\"4\",\"5\")", "in(level,\"0\",\"3\",\"4\",\"5\")",
			"in(level,\"1\",\"2\",\"3\",\"4\")", "in(level,\"1\",\"2\",\"3\",\"5\")",
			"in(level,\"1\",\"2\",\"4\",\"5\")", "in(level,\"1\",\"3\",\"4\",\"5\")",
			"in(level,\"2\",\"3\",\"4\",\"5\")",
			"ni(level,\"0\",\"1\")", "ni(level,\"0\",\"2\")", "ni(level,\"0\",\"3\")",
			"ni(level,\"0\",\"4\")", "ni(level,\"0\",\"5\")",
			"ni(level,\"1\",\"2\")", "ni(level,\"1\",\"3\")", "ni(level,\"1\",\"4\")",
			"ni(level,\"1\",\"5\")",
			"ni(level,\"2\",\"3\")", "ni(level,\"2\",\"4\")", "ni(level,\"2\",\"5\")",
			"ni(level,\"3\",\"4\")", "ni(level,\"3\",\"5\")",
			"ni(level,\"4\",\"5\")",
			"ni(level,\"0\",\"1\",\"2\")", "ni(level,\"0\",\"1\",\"3\")",
			"ni(level,\"0\",\"1\",\"4\")", "ni(level,\"0\",\"1\",\"5\")",
			"ni(level,\"0\",\"2\",\"3\")", "ni(level,\"0\",\"2\",\"4\")",
			"ni(level,\"0\",\"2\",\"5\")", "ni(level,\"0\",\"3\",\"4\")",
			"ni(level,\"0\",\"3\",\"5\")", "ni(level,\"0\",\"4\",\"5\")",
			"ni(level,\"1\",\"2\",\"3\")", "ni(level,\"1\",\"2\",\"4\")",
			"ni(level,\"1\",\"2\",\"5\")", "ni(level,\"1\",\"3\",\"4\")",
			"ni(level,\"1\",\"3\",\"5\")", "ni(level,\"1\",\"4\",\"5\")",
			"ni(level,\"2\",\"3\",\"4\")", "ni(level,\"2\",\"3\",\"5\")",
			"ni(level,\"2\",\"4\",\"5\")", "ni(level,\"3\",\"4\",\"5\")",
			"ni(level,\"0\",\"1\",\"2\",\"3\")", "ni(level,\"0\",\"1\",\"2\",\"4\")",
			"ni(level,\"0\",\"1\",\"2\",\"5\")", "ni(level,\"0\",\"1\",\"3\",\"4\")",
			"ni(level,\"0\",\"1\",\"3\",\"5\")", "ni(level,\"0\",\"1\",\"4\",\"5\")",
			"ni(level,\"0\",\"2\",\"3\",\"4\")", "ni(level,\"0\",\"2\",\"3\",\"5\")",
			"ni(level,\"0\",\"2\",\"4\",\"5\")", "ni(level,\"0\",\"3\",\"4\",\"5\")",
			"ni(level,\"1\",\"2\",\"3\",\"4\")", "ni(level,\"1\",\"2\",\"3\",\"5\")",
			"ni(level,\"1\",\"2\",\"4\",\"5\")", "ni(level,\"1\",\"3\",\"4\",\"5\")",
			"ni(level,\"2\",\"3\",\"4\",\"5\")",
		}},
		{"flow", []string{
			"flow == 0", "flow == 2.34", "flow == 10.5",
			"flow != 0", "flow != 2.34", "flow != 10.5",
			"flow < level", "flow <= level", "flow != level",
			"flow >= level", "flow > level", "flow == level",
			"flow < position", "flow <= position", "flow != position",
			"flow >= position", "flow > position", "flow == position",
			"flow >= 0",
			"flow >= 1.05",
			"flow >= 2.1",
			"flow >= 3.15",
			"flow >= 4.2",
			"flow >= 5.25",
			"flow >= 6.3",
			"flow >= 7.35",
			"flow >= 8.4",
			"flow >= 9.45",
			"flow <= 1.05",
			"flow <= 2.1",
			"flow <= 3.15",
			"flow <= 4.2",
			"flow <= 5.25",
			"flow <= 6.3",
			"flow <= 7.35",
			"flow <= 8.4",
			"flow <= 9.45",
		}},
		{"position", []string{
			"position == 1", "position == 2", "position == 3",
			"position == 4", "position == 5", "position == 6",
			"position == 7", "position == 8", "position == 9",
			"position == 10", "position == 11", "position == 12",
			"position == 13",
			"position != 1", "position != 2", "position != 3",
			"position != 4", "position != 5", "position != 6",
			"position != 7", "position != 8", "position != 9",
			"position != 10", "position != 11", "position != 12",
			"position != 13",
			"position >= 1", "position >= 2", "position >= 3",
			"position >= 4", "position >= 5", "position >= 6",
			"position >= 7", "position >= 8", "position >= 9",
			"position >= 10", "position >= 11", "position >= 12",
			"position <= 2", "position <= 3", "position <= 4",
			"position <= 5", "position <= 6", "position <= 7",
			"position <= 8", "position <= 9", "position <= 10",
			"position <= 11", "position <= 12", "position <= 13",
		}},
	}

	rules, err := GenerateRules(fieldDescriptions, excludeFields)
	if err != nil {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("GenerateRules(%q, %q) err: %q",
			fieldDescriptions, excludeFields, err)
	}

	for _, c := range cases {
		gotFieldRules := getFieldRules(c.field, rules)
		rulesMatch, msg := matchRules(gotFieldRules, c.wantRules)
		if !rulesMatch {
			sort.Strings(gotFieldRules)
			sort.Strings(c.wantRules)
			t.Errorf("Test: %s\n", testPurpose)
			t.Errorf("matchRules() rules don't match for field: %s - %s\ngot: %s\nwant: %s\n",
				c.field, msg, gotFieldRules, c.wantRules)
		}
	}
}

func TestGenerateRules_2(t *testing.T) {
	testPurpose := "Ensure generates a 'true()' rule"
	fieldDescriptions := map[string]*FieldDescription{
		"team": &FieldDescription{
			Kind: STRING,
			Values: []*dlit.Literal{
				dlit.MustNew("a"), dlit.MustNew("b"), dlit.MustNew("c"),
			},
		},
		"teamOut": &FieldDescription{
			Kind: STRING,
			Values: []*dlit.Literal{
				dlit.MustNew("a"), dlit.MustNew("c"), dlit.MustNew("d"),
				dlit.MustNew("e"), dlit.MustNew("f"),
			},
		},
	}
	excludeFields := []string{}

	rules, err := GenerateRules(fieldDescriptions, excludeFields)
	if err != nil {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("GenerateRules(%q, %q) err: %q",
			fieldDescriptions, excludeFields, err)
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
			fieldDescriptions, excludeFields)
	}
}

func TestCombinedRules(t *testing.T) {
	cases := []struct {
		inRules   []*rule.Rule
		wantRules []string
	}{
		{[]*rule.Rule{
			rule.MustNew("team == \"a\""),
			rule.MustNew("band > 4"),
			rule.MustNew("in(team,\"red\",\"green\",\"blue\")"),
		},
			[]string{
				"team == \"a\" && band > 4",
				"team == \"a\" || band > 4",
				"team == \"a\" && in(team,\"red\",\"green\",\"blue\")",
				"team == \"a\" || in(team,\"red\",\"green\",\"blue\")",
				"band > 4 && in(team,\"red\",\"green\",\"blue\")",
				"band > 4 || in(team,\"red\",\"green\",\"blue\")",
			}},
		{[]*rule.Rule{rule.MustNew("team == \"a\"")}, []string{}},
		{[]*rule.Rule{}, []string{}},
	}

	for _, c := range cases {
		gotRules := CombineRules(c.inRules)
		gotRuleStrs := rulesToStrings(gotRules)
		rulesMatch, msg := matchRules(gotRuleStrs, c.wantRules)
		if !rulesMatch {
			sort.Strings(gotRuleStrs)
			sort.Strings(c.wantRules)
			t.Errorf("matchRules() rules don't match: %s\ngot: %s\nwant: %s\n",
				msg, gotRuleStrs, c.wantRules)
		}
	}
}

/*************************************
 *    Helper Functions
 *************************************/
var matchFieldInNiRegexp = regexp.MustCompile("^((in\\(|ni\\()+)([^ ,]+)(.*)$")
var matchFieldMatchRegexp = regexp.MustCompile("^([^ (]+)( .*)$")

func getFieldRules(
	field string, rules []*rule.Rule) []string {
	fieldRules := make([]string, 0)
	for _, rule := range rules {
		ruleStr := rule.String()
		ruleField := matchFieldMatchRegexp.ReplaceAllString(ruleStr, "$1")
		ruleField = matchFieldInNiRegexp.ReplaceAllString(ruleField, "$3")
		if field == ruleField {
			fieldRules = append(fieldRules, ruleStr)
		}
	}
	return fieldRules
}

func rulesToStrings(rules []*rule.Rule) []string {
	r := make([]string, len(rules))
	for i, rule := range rules {
		r[i] = rule.String()
	}
	return r
}
