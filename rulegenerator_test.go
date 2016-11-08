package rhkit

import (
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/rule"
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
	ruleFields := []string{"team", "teamOut", "camp", "level", "flow", "position"}
	cases := []struct {
		field     string
		wantRules []rule.Rule
	}{
		{"team", []rule.Rule{
			rule.NewEQFVS("team", "a"),
			rule.NewEQFVS("team", "b"), rule.NewEQFVS("team", "c"),
			rule.NewNEFVS("team", "a"),
			rule.NewNEFVS("team", "b"),
			rule.NewNEFVS("team", "c"),
			rule.NewEQFF("team", "teamOut"),
			rule.NewNEFF("team", "teamOut"),
		}},
		{"teamOut", []rule.Rule{
			rule.NewEQFVS("teamOut", "a"),
			rule.NewEQFVS("teamOut", "c"),
			rule.NewEQFVS("teamOut", "d"),
			rule.NewEQFVS("teamOut", "e"),
			rule.NewEQFVS("teamOut", "f"),
			rule.NewNEFVS("teamOut", "a"),
			rule.NewNEFVS("teamOut", "c"),
			rule.NewNEFVS("teamOut", "d"),
			rule.NewNEFVS("teamOut", "e"),
			rule.NewNEFVS("teamOut", "f"),
			rule.NewInFV("teamOut", makeStringsDlitSlice("a", "c")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("a", "d")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("a", "e")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("a", "f")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("c", "d")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("c", "e")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("c", "f")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("d", "e")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("d", "f")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("e", "f")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("a", "c", "d")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("a", "c", "e")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("a", "c", "f")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("a", "d", "e")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("a", "d", "f")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("a", "e", "f")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("c", "d", "e")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("c", "d", "f")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("c", "e", "f")),
			rule.NewInFV("teamOut", makeStringsDlitSlice("d", "e", "f")),
		}},
		{"level", []rule.Rule{
			rule.NewEQFVI("level", 0),
			rule.NewEQFVI("level", 1),
			rule.NewEQFVI("level", 2),
			rule.NewEQFVI("level", 3),
			rule.NewEQFVI("level", 4),
			rule.NewEQFVI("level", 5),
			rule.NewNEFVI("level", 0),
			rule.NewNEFVI("level", 1),
			rule.NewNEFVI("level", 2),
			rule.NewNEFVI("level", 3),
			rule.NewNEFVI("level", 4),
			rule.NewNEFVI("level", 5),
			rule.NewLTFF("level", "position"),
			rule.NewLEFF("level", "position"),
			rule.NewNEFF("level", "position"),
			rule.NewGEFF("level", "position"),
			rule.NewGTFF("level", "position"),
			rule.NewEQFF("level", "position"),
			rule.NewGEFVI("level", 0),
			rule.NewGEFVI("level", 1),
			rule.NewGEFVI("level", 2),
			rule.NewGEFVI("level", 3),
			rule.NewGEFVI("level", 4),
			rule.NewLEFVI("level", 1),
			rule.NewLEFVI("level", 2),
			rule.NewLEFVI("level", 3),
			rule.NewLEFVI("level", 4),
			rule.NewLEFVI("level", 5),
			rule.NewInFV("level", makeStringsDlitSlice("0", "1")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "2")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "3")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "4")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("1", "2")),
			rule.NewInFV("level", makeStringsDlitSlice("1", "3")),
			rule.NewInFV("level", makeStringsDlitSlice("1", "4")),
			rule.NewInFV("level", makeStringsDlitSlice("1", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("2", "3")),
			rule.NewInFV("level", makeStringsDlitSlice("2", "4")),
			rule.NewInFV("level", makeStringsDlitSlice("2", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("3", "4")),
			rule.NewInFV("level", makeStringsDlitSlice("3", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("4", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "1", "2")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "1", "3")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "1", "4")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "1", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "2", "3")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "2", "4")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "2", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "3", "4")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "3", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "4", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("1", "2", "3")),
			rule.NewInFV("level", makeStringsDlitSlice("1", "2", "4")),
			rule.NewInFV("level", makeStringsDlitSlice("1", "2", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("1", "3", "4")),
			rule.NewInFV("level", makeStringsDlitSlice("1", "3", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("1", "4", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("2", "3", "4")),
			rule.NewInFV("level", makeStringsDlitSlice("2", "3", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("2", "4", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("3", "4", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "1", "2", "3")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "1", "2", "4")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "1", "2", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "1", "3", "4")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "1", "3", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "1", "4", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "2", "3", "4")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "2", "3", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "2", "4", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("0", "3", "4", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("1", "2", "3", "4")),
			rule.NewInFV("level", makeStringsDlitSlice("1", "2", "3", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("1", "2", "4", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("1", "3", "4", "5")),
			rule.NewInFV("level", makeStringsDlitSlice("2", "3", "4", "5")),
		}},
		{"flow", []rule.Rule{
			rule.NewEQFVI("flow", 0),
			rule.NewEQFVF("flow", 2.34),
			rule.NewEQFVF("flow", 10.5),
			rule.NewNEFVI("flow", 0),
			rule.NewNEFVF("flow", 2.34),
			rule.NewNEFVF("flow", 10.5),
			rule.NewLTFF("flow", "level"),
			rule.NewLEFF("flow", "level"),
			rule.NewNEFF("flow", "level"),
			rule.NewGEFF("flow", "level"),
			rule.NewGTFF("flow", "level"),
			rule.NewEQFF("flow", "level"),
			rule.NewLTFF("flow", "position"),
			rule.NewLEFF("flow", "position"),
			rule.NewNEFF("flow", "position"),
			rule.NewGEFF("flow", "position"),
			rule.NewGTFF("flow", "position"),
			rule.NewEQFF("flow", "position"),
			rule.NewGEFVF("flow", 0.0),
			rule.NewGEFVF("flow", 1.05),
			rule.NewGEFVF("flow", 2.1),
			rule.NewGEFVF("flow", 3.15),
			rule.NewGEFVF("flow", 4.2),
			rule.NewGEFVF("flow", 5.25),
			rule.NewGEFVF("flow", 6.3),
			rule.NewGEFVF("flow", 7.35),
			rule.NewGEFVF("flow", 8.4),
			rule.NewGEFVF("flow", 9.45),
			rule.NewLEFVF("flow", 1.05),
			rule.NewLEFVF("flow", 2.1),
			rule.NewLEFVF("flow", 3.15),
			rule.NewLEFVF("flow", 4.2),
			rule.NewLEFVF("flow", 5.25),
			rule.NewLEFVF("flow", 6.3),
			rule.NewLEFVF("flow", 7.35),
			rule.NewLEFVF("flow", 8.4),
			rule.NewLEFVF("flow", 9.45),
		}},
		{"position", []rule.Rule{
			rule.NewEQFVI("position", 1),
			rule.NewEQFVI("position", 2),
			rule.NewEQFVI("position", 3),
			rule.NewEQFVI("position", 4),
			rule.NewEQFVI("position", 5),
			rule.NewEQFVI("position", 6),
			rule.NewEQFVI("position", 7),
			rule.NewEQFVI("position", 8),
			rule.NewEQFVI("position", 9),
			rule.NewEQFVI("position", 10),
			rule.NewEQFVI("position", 11),
			rule.NewEQFVI("position", 12),
			rule.NewEQFVI("position", 13),
			rule.NewNEFVI("position", 1),
			rule.NewNEFVI("position", 2),
			rule.NewNEFVI("position", 3),
			rule.NewNEFVI("position", 4),
			rule.NewNEFVI("position", 5),
			rule.NewNEFVI("position", 6),
			rule.NewNEFVI("position", 7),
			rule.NewNEFVI("position", 8),
			rule.NewNEFVI("position", 9),
			rule.NewNEFVI("position", 10),
			rule.NewNEFVI("position", 11),
			rule.NewNEFVI("position", 12),
			rule.NewNEFVI("position", 13),
			rule.NewGEFVI("position", 1),
			rule.NewGEFVI("position", 2),
			rule.NewGEFVI("position", 3),
			rule.NewGEFVI("position", 4),
			rule.NewGEFVI("position", 5),
			rule.NewGEFVI("position", 6),
			rule.NewGEFVI("position", 7),
			rule.NewGEFVI("position", 8),
			rule.NewGEFVI("position", 9),
			rule.NewGEFVI("position", 10),
			rule.NewGEFVI("position", 11),
			rule.NewGEFVI("position", 12),
			rule.NewLEFVI("position", 2),
			rule.NewLEFVI("position", 3),
			rule.NewLEFVI("position", 4),
			rule.NewLEFVI("position", 5),
			rule.NewLEFVI("position", 6),
			rule.NewLEFVI("position", 7),
			rule.NewLEFVI("position", 8),
			rule.NewLEFVI("position", 9),
			rule.NewLEFVI("position", 10),
			rule.NewLEFVI("position", 11),
			rule.NewLEFVI("position", 12),
			rule.NewLEFVI("position", 13),
		}},
	}

	rules, err := GenerateRules(inputDescription, ruleFields)
	if err != nil {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("GenerateRules(%v, %v) err: %v", inputDescription, ruleFields, err)
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
	ruleFields := []string{"team", "teamOut"}

	rules, err := GenerateRules(inputDescription, ruleFields)
	if err != nil {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("GenerateRules(%v, %v) err: %v", inputDescription, ruleFields, err)
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
		t.Errorf("GenerateRules(%v, %v)  - True rule missing",
			inputDescription, ruleFields)
	}
}

func TestGenerateRules_3(t *testing.T) {
	testPurpose := "Ensure generates correct combination rules"
	inputDescription := &Description{
		map[string]*fieldDescription{
			"team": &fieldDescription{
				kind: ftString,
				values: []*dlit.Literal{
					dlit.MustNew("a"), dlit.MustNew("b"), dlit.MustNew("c"),
				},
			},
			"teamAlt": &fieldDescription{
				kind: ftString,
				values: []*dlit.Literal{
					dlit.MustNew("a"),
					dlit.MustNew("b"),
					dlit.MustNew("c"),
					dlit.MustNew("d"),
				},
			},
			"monthIn": &fieldDescription{
				kind: ftInt,
				min:  dlit.MustNew(1),
				max:  dlit.MustNew(3),
				values: []*dlit.Literal{
					dlit.MustNew(1), dlit.MustNew(2), dlit.MustNew(3),
				},
			},
			"monthOut": &fieldDescription{
				kind: ftInt,
				min:  dlit.MustNew(1),
				max:  dlit.MustNew(3),
				values: []*dlit.Literal{
					dlit.MustNew(1), dlit.MustNew(2), dlit.MustNew(3), dlit.MustNew(4),
				},
			},
			"win": &fieldDescription{
				kind:   ftString,
				values: []*dlit.Literal{dlit.MustNew("t"), dlit.MustNew("f")},
			},
		}}
	cases := []struct {
		field      string
		ruleFields []string
		wantRules  []rule.Rule
	}{
		{field: "team",
			ruleFields: []string{"team", "monthIn"},
			wantRules: []rule.Rule{
				rule.NewEQFVS("team", "a"),
				rule.NewEQFVS("team", "b"),
				rule.NewEQFVS("team", "c"),
				rule.NewNEFVS("team", "a"),
				rule.NewNEFVS("team", "b"),
				rule.NewNEFVS("team", "c"),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewNEFVI("monthIn", 3)),
			}},
		{field: "team",
			ruleFields: []string{"team", "monthOut"},
			wantRules: []rule.Rule{
				rule.NewEQFVS("team", "a"),
				rule.NewEQFVS("team", "b"),
				rule.NewEQFVS("team", "c"),
				rule.NewNEFVS("team", "a"),
				rule.NewNEFVS("team", "b"),
				rule.NewNEFVS("team", "c"),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewEQFVI("monthOut", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewEQFVI("monthOut", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewEQFVI("monthOut", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewEQFVI("monthOut", 4)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewNEFVI("monthOut", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewNEFVI("monthOut", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewNEFVI("monthOut", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewNEFVI("monthOut", 4)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewEQFVI("monthOut", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewEQFVI("monthOut", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewEQFVI("monthOut", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewEQFVI("monthOut", 4)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewNEFVI("monthOut", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewNEFVI("monthOut", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewNEFVI("monthOut", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewNEFVI("monthOut", 4)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewEQFVI("monthOut", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewEQFVI("monthOut", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewEQFVI("monthOut", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewEQFVI("monthOut", 4)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewNEFVI("monthOut", 1)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewNEFVI("monthOut", 2)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewNEFVI("monthOut", 3)),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewNEFVI("monthOut", 4)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewEQFVI("monthOut", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewEQFVI("monthOut", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewEQFVI("monthOut", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewEQFVI("monthOut", 4)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewNEFVI("monthOut", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewNEFVI("monthOut", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewNEFVI("monthOut", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewNEFVI("monthOut", 4)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewEQFVI("monthOut", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewEQFVI("monthOut", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewEQFVI("monthOut", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewEQFVI("monthOut", 4)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewNEFVI("monthOut", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewNEFVI("monthOut", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewNEFVI("monthOut", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewNEFVI("monthOut", 4)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewEQFVI("monthOut", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewEQFVI("monthOut", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewEQFVI("monthOut", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewEQFVI("monthOut", 4)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewNEFVI("monthOut", 1)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewNEFVI("monthOut", 2)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewNEFVI("monthOut", 3)),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewNEFVI("monthOut", 4)),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "2"))),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "3"))),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "4"))),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "3"))),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "4"))),
				rule.NewAnd(rule.NewEQFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("3", "4"))),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "2"))),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "3"))),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "4"))),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "3"))),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "4"))),
				rule.NewAnd(rule.NewEQFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("3", "4"))),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "2"))),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "3"))),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "4"))),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "3"))),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "4"))),
				rule.NewAnd(rule.NewEQFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("3", "4"))),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "2"))),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "3"))),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "4"))),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "3"))),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "4"))),
				rule.NewAnd(rule.NewNEFVS("team", "a"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("3", "4"))),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "2"))),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "3"))),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "4"))),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "3"))),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "4"))),
				rule.NewAnd(rule.NewNEFVS("team", "b"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("3", "4"))),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "2"))),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "3"))),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("1", "4"))),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "3"))),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("2", "4"))),
				rule.NewAnd(rule.NewNEFVS("team", "c"),
					rule.NewInFV("monthOut", makeStringsDlitSlice("3", "4"))),
			}},
		{field: "teamAlt",
			ruleFields: []string{"teamAlt", "monthIn"},
			wantRules: []rule.Rule{
				rule.NewEQFVS("teamAlt", "a"),
				rule.NewEQFVS("teamAlt", "b"),
				rule.NewEQFVS("teamAlt", "c"),
				rule.NewEQFVS("teamAlt", "d"),
				rule.NewNEFVS("teamAlt", "a"),
				rule.NewNEFVS("teamAlt", "b"),
				rule.NewNEFVS("teamAlt", "c"),
				rule.NewNEFVS("teamAlt", "d"),
				rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "b")),
				rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "c")),
				rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "d")),
				rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "c")),
				rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "d")),
				rule.NewInFV("teamAlt", makeStringsDlitSlice("c", "d")),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "a"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "a"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "a"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "a"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "a"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "a"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "b"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "b"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "b"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "b"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "b"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "b"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "c"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "c"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "c"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "c"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "c"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "c"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "d"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "d"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "d"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "d"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "d"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewEQFVS("teamAlt", "d"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "a"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "a"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "a"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "a"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "a"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "a"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "b"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "b"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "b"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "b"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "b"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "b"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "c"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "c"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "c"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "c"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "c"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "c"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "d"),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "d"),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "d"),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "d"),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "d"),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewNEFVS("teamAlt", "d"),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "b")),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "b")),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "b")),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "b")),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "b")),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "b")),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "c")),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "c")),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "c")),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "c")),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "c")),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "c")),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "d")),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "d")),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "d")),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "d")),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "d")),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("a", "d")),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "c")),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "c")),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "c")),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "c")),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "c")),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "c")),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "d")),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "d")),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "d")),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "d")),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "d")),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("b", "d")),
					rule.NewNEFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("c", "d")),
					rule.NewEQFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("c", "d")),
					rule.NewEQFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("c", "d")),
					rule.NewEQFVI("monthIn", 3)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("c", "d")),
					rule.NewNEFVI("monthIn", 1)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("c", "d")),
					rule.NewNEFVI("monthIn", 2)),
				rule.NewAnd(rule.NewInFV("teamAlt", makeStringsDlitSlice("c", "d")),
					rule.NewNEFVI("monthIn", 3)),
			}},
	}

	for _, c := range cases {
		rules, err := GenerateRules(inputDescription, c.ruleFields)
		if err != nil {
			t.Errorf("Test: %s\n", testPurpose)
			t.Fatalf("GenerateRules(%v, %v) err: %v",
				inputDescription, c.ruleFields, err)
		}

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

func TestCombinedRules(t *testing.T) {
	cases := []struct {
		inRules   []rule.Rule
		wantRules []rule.Rule
	}{
		{[]rule.Rule{
			rule.NewEQFVS("team", "a"),
			rule.NewGEFVI("band", 4),
			rule.NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
		},
			[]rule.Rule{
				rule.NewAnd(rule.NewEQFVS("team", "a"), rule.NewGEFVI("band", 4)),
				rule.NewOr(rule.NewEQFVS("team", "a"), rule.NewGEFVI("band", 4)),
				rule.NewAnd(
					rule.NewEQFVS("team", "a"),
					rule.NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
				),
				rule.NewOr(
					rule.NewEQFVS("team", "a"),
					rule.NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
				),
				rule.NewAnd(
					rule.NewGEFVI("band", 4),
					rule.NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
				),
				rule.NewOr(
					rule.NewGEFVI("band", 4),
					rule.NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
				),
			}},
		{[]rule.Rule{rule.NewEQFVS("team", "a")}, []rule.Rule{}},
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
var matchFieldInRegexp = regexp.MustCompile("^((in\\()+)([^ ,]+)(.*)$")
var matchFieldMatchRegexp = regexp.MustCompile("^([^ (]+)( .*)$")

func getFieldRules(
	field string,
	rules []rule.Rule,
) []rule.Rule {
	fieldRules := make([]rule.Rule, 0)
	for _, rule := range rules {
		ruleStr := rule.String()
		ruleField := matchFieldMatchRegexp.ReplaceAllString(ruleStr, "$1")
		ruleField = matchFieldInRegexp.ReplaceAllString(ruleField, "$3")
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
