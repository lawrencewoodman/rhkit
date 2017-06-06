package rhkit

import (
	"errors"
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
	"github.com/vlifesystems/rhkit/rule"
	"regexp"
	"sort"
	"testing"
)

func TestGenerateRules_1(t *testing.T) {
	testPurpose := "Ensure generates correct rules for each field"
	inputDescription := &description.Description{
		map[string]*description.Field{
			"team": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"a": description.Value{dlit.NewString("a"), 3},
					"b": description.Value{dlit.NewString("b"), 3},
					"c": description.Value{dlit.NewString("c"), 3},
				},
			},
			"teamOut": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"a": description.Value{dlit.NewString("a"), 3},
					"c": description.Value{dlit.NewString("c"), 1},
					"d": description.Value{dlit.NewString("d"), 3},
					"e": description.Value{dlit.NewString("e"), 3},
					"f": description.Value{dlit.NewString("f"), 3},
				},
			},
			"level": &description.Field{
				Kind:  fieldtype.Number,
				Min:   dlit.MustNew(0),
				Max:   dlit.MustNew(5),
				MaxDP: 0,
				Values: map[string]description.Value{
					"0": description.Value{dlit.NewString("0"), 3},
					"1": description.Value{dlit.NewString("1"), 3},
					"2": description.Value{dlit.NewString("2"), 1},
					"3": description.Value{dlit.NewString("3"), 3},
					"4": description.Value{dlit.NewString("4"), 3},
					"5": description.Value{dlit.NewString("5"), 3},
				},
			},
			"flow": &description.Field{
				Kind:  fieldtype.Number,
				Min:   dlit.MustNew(0),
				Max:   dlit.MustNew(10.5),
				MaxDP: 2,
				Values: map[string]description.Value{
					"0.0":  description.Value{dlit.NewString("0.0"), 3},
					"2.34": description.Value{dlit.NewString("2.34"), 3},
					"10.5": description.Value{dlit.NewString("10.5"), 3},
				},
			},
			"position": &description.Field{
				Kind:  fieldtype.Number,
				Min:   dlit.MustNew(1),
				Max:   dlit.MustNew(13),
				MaxDP: 0,
				Values: map[string]description.Value{
					"1":  description.Value{dlit.NewString("1"), 3},
					"2":  description.Value{dlit.NewString("2"), 3},
					"3":  description.Value{dlit.NewString("3"), 3},
					"4":  description.Value{dlit.NewString("4"), 3},
					"5":  description.Value{dlit.NewString("5"), 3},
					"6":  description.Value{dlit.NewString("6"), 3},
					"7":  description.Value{dlit.NewString("7"), 3},
					"8":  description.Value{dlit.NewString("8"), 3},
					"9":  description.Value{dlit.NewString("9"), 3},
					"10": description.Value{dlit.NewString("10"), 3},
					"11": description.Value{dlit.NewString("11"), 3},
					"12": description.Value{dlit.NewString("12"), 3},
					"13": description.Value{dlit.NewString("13"), 3},
				},
			},
		}}
	ruleFields :=
		[]string{"team", "teamOut", "level", "flow", "position"}
	wantRules := []rule.Rule{
		rule.NewEQFV("team", dlit.MustNew("a")),
		rule.NewNEFV("team", dlit.MustNew("a")),
		rule.NewEQFF("team", "teamOut"),
		rule.NewNEFF("team", "teamOut"),
		rule.NewInFV("teamOut", makeStringsDlitSlice("a", "d")),
		rule.NewEQFV("level", dlit.MustNew(0)),
		rule.NewEQFV("level", dlit.MustNew(1)),
		rule.NewNEFV("level", dlit.MustNew(0)),
		rule.NewNEFV("level", dlit.MustNew(1)),
		rule.NewLTFF("level", "position"),
		rule.NewLEFF("level", "position"),
		rule.NewNEFF("level", "position"),
		rule.NewGEFF("level", "position"),
		rule.NewGTFF("level", "position"),
		rule.NewEQFF("level", "position"),
		rule.NewGEFV("level", dlit.MustNew(1)),
		rule.NewLEFV("level", dlit.MustNew(4)),
		rule.NewInFV("level", makeStringsDlitSlice("0", "1")),
		rule.NewInFV("level", makeStringsDlitSlice("0", "3")),
		rule.NewGEFV("flow", dlit.MustNew(2.1)),
		rule.NewGEFV("flow", dlit.MustNew(3.68)),
		rule.NewLEFV("flow", dlit.MustNew(4.2)),
		rule.NewLEFV("flow", dlit.MustNew(5.25)),
		rule.NewAddLEF("level", "position", dlit.MustNew(12)),
		rule.NewAddGEF("level", "position", dlit.MustNew(12)),
		rule.NewMulLEF("flow", "level", dlit.MustNew(26.25)),
		rule.NewMulGEF("flow", "level", dlit.MustNew(23.63)),
	}
	complexity := 10
	got := GenerateRules(inputDescription, ruleFields, complexity)
	if err := rulesContain(got, wantRules); err != nil {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("GenerateRules: %s", err)
	}
}

func TestGenerateRules_2(t *testing.T) {
	testPurpose := "Ensure generates a True rule"
	inputDescription := &description.Description{
		map[string]*description.Field{
			"team": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"a": description.Value{dlit.MustNew("a"), 3},
					"b": description.Value{dlit.MustNew("b"), 3},
					"c": description.Value{dlit.MustNew("c"), 3},
				},
			},
			"teamOut": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"a": description.Value{dlit.MustNew("a"), 3},
					"c": description.Value{dlit.MustNew("c"), 3},
					"d": description.Value{dlit.MustNew("d"), 3},
					"e": description.Value{dlit.MustNew("e"), 3},
					"f": description.Value{dlit.MustNew("f"), 3},
				},
			},
		}}
	ruleFields := []string{"team", "teamOut"}
	complexity := 10
	got := GenerateRules(inputDescription, ruleFields, complexity)

	trueRuleFound := false
	for _, r := range got {
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
	inputDescription := &description.Description{
		map[string]*description.Field{
			"directionIn": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"gogledd": description.Value{dlit.MustNew("gogledd"), 3},
					"de":      description.Value{dlit.MustNew("de"), 3},
				},
			},
			"directionOut": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"dwyrain":   description.Value{dlit.MustNew("dwyrain"), 3},
					"gorllewin": description.Value{dlit.MustNew("gorllewin"), 3},
				},
			},
		}}
	ruleFields := []string{"directionIn", "directionOut"}
	want := []rule.Rule{
		rule.NewEQFV("directionIn", dlit.MustNew("de")),
		rule.NewEQFV("directionIn", dlit.MustNew("gogledd")),
		rule.NewEQFV("directionOut", dlit.MustNew("dwyrain")),
		rule.NewEQFV("directionOut", dlit.MustNew("gorllewin")),
		rule.MustNewAnd(
			rule.NewEQFV("directionIn", dlit.MustNew("de")),
			rule.NewEQFV("directionOut", dlit.MustNew("dwyrain")),
		),
		rule.MustNewAnd(
			rule.NewEQFV("directionIn", dlit.MustNew("de")),
			rule.NewEQFV("directionOut", dlit.MustNew("gorllewin")),
		),
		rule.MustNewAnd(
			rule.NewEQFV("directionIn", dlit.MustNew("gogledd")),
			rule.NewEQFV("directionOut", dlit.MustNew("dwyrain")),
		),
		rule.MustNewAnd(
			rule.NewEQFV("directionIn", dlit.MustNew("gogledd")),
			rule.NewEQFV("directionOut", dlit.MustNew("gorllewin")),
		),
		rule.MustNewOr(
			rule.NewEQFV("directionIn", dlit.MustNew("de")),
			rule.NewEQFV("directionIn", dlit.MustNew("gogledd")),
		),
		rule.MustNewOr(
			rule.NewEQFV("directionIn", dlit.MustNew("de")),
			rule.NewEQFV("directionOut", dlit.MustNew("dwyrain")),
		),
		rule.MustNewOr(
			rule.NewEQFV("directionIn", dlit.MustNew("de")),
			rule.NewEQFV("directionOut", dlit.MustNew("gorllewin")),
		),
		rule.MustNewOr(
			rule.NewEQFV("directionIn", dlit.MustNew("gogledd")),
			rule.NewEQFV("directionOut", dlit.MustNew("dwyrain")),
		),
		rule.MustNewOr(
			rule.NewEQFV("directionIn", dlit.MustNew("gogledd")),
			rule.NewEQFV("directionOut", dlit.MustNew("gorllewin")),
		),
		rule.MustNewOr(
			rule.NewEQFV("directionOut", dlit.MustNew("dwyrain")),
			rule.NewEQFV("directionOut", dlit.MustNew("gorllewin")),
		),
		rule.NewTrue(),
	}

	complexity := 10
	got := GenerateRules(inputDescription, ruleFields, complexity)
	rule.Sort(got)
	rule.Sort(want)
	if err := matchRulesUnordered(got, want); err != nil {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("matchRulesUnordered: %s\n got: %s\nwant: %s\n",
			err, got, want)
		for i, r := range got {
			t.Errorf("rule(%d): %s", i, r)
		}
	}
}

func TestGenerateCompareNumericRules(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"band": &description.Field{
				Kind:   fieldtype.Number,
				Min:    dlit.MustNew(1),
				Max:    dlit.MustNew(3),
				Values: map[string]description.Value{},
			},
			"flowIn": &description.Field{
				Kind:   fieldtype.Number,
				Min:    dlit.MustNew(1),
				Max:    dlit.MustNew(4),
				MaxDP:  2,
				Values: map[string]description.Value{},
			},
			"flowOut": &description.Field{
				Kind:   fieldtype.Number,
				Min:    dlit.MustNew(0.95),
				Max:    dlit.MustNew(4.1),
				MaxDP:  2,
				Values: map[string]description.Value{},
			},
			"rateIn": &description.Field{
				Kind:   fieldtype.Number,
				Min:    dlit.MustNew(4.2),
				Max:    dlit.MustNew(8.9),
				MaxDP:  2,
				Values: map[string]description.Value{},
			},
			"rateOut": &description.Field{
				Kind:   fieldtype.Number,
				Min:    dlit.MustNew(0.1),
				Max:    dlit.MustNew(0.9),
				MaxDP:  2,
				Values: map[string]description.Value{},
			},
			"group": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Nelson":      description.Value{dlit.NewString("Nelson"), 3},
					"Collingwood": description.Value{dlit.NewString("Collingwood"), 1},
					"Mountbatten": description.Value{dlit.NewString("Mountbatten"), 1},
					"Drake":       description.Value{dlit.NewString("Drake"), 2},
				},
			},
		},
	}
	cases := []struct {
		field string
		want  []rule.Rule
	}{
		{field: "band",
			want: []rule.Rule{
				rule.NewLTFF("band", "flowIn"),
				rule.NewLTFF("band", "flowOut"),
				rule.NewLEFF("band", "flowIn"),
				rule.NewLEFF("band", "flowOut"),
				rule.NewGTFF("band", "flowIn"),
				rule.NewGTFF("band", "flowOut"),
				rule.NewGEFF("band", "flowIn"),
				rule.NewGEFF("band", "flowOut"),
			},
		},
		{field: "flowIn",
			want: []rule.Rule{
				rule.NewLTFF("flowIn", "flowOut"),
				rule.NewLEFF("flowIn", "flowOut"),
				rule.NewGTFF("flowIn", "flowOut"),
				rule.NewGEFF("flowIn", "flowOut"),
			},
		},
		{field: "flowOut",
			want: []rule.Rule{},
		},
		{field: "rateIn",
			want: []rule.Rule{},
		},
		{field: "rateOut",
			want: []rule.Rule{},
		},
		{field: "group",
			want: []rule.Rule{},
		},
	}
	ruleFields :=
		[]string{"band", "flowIn", "flowOut", "rateIn", "rateOut", "group"}
	complexity := 10
	for _, c := range cases {
		got := generateCompareNumericRules(
			inputDescription,
			ruleFields,
			complexity,
			c.field,
		)
		if err := matchRulesUnordered(got, c.want); err != nil {
			gotRuleStrs := rulesToSortedStrings(got)
			wantRuleStrs := rulesToSortedStrings(c.want)
			t.Errorf("matchRulesUnordered() rules don't match: %s\ngot: %s\nwant: %s\n",
				err, gotRuleStrs, wantRuleStrs)
		}
	}
}

func TestCombineRules(t *testing.T) {
	cases := []struct {
		in   []rule.Rule
		want []rule.Rule
	}{
		{in: []rule.Rule{
			rule.NewEQFV("group", dlit.MustNew("a")),
			rule.NewGEFV("band", dlit.MustNew(4)),
			rule.NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
		},
			want: []rule.Rule{
				rule.MustNewAnd(
					rule.NewGEFV("band", dlit.MustNew(4)),
					rule.NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
				),
				rule.MustNewAnd(
					rule.NewGEFV("band", dlit.MustNew(4)),
					rule.NewEQFV("group", dlit.MustNew("a")),
				),
				rule.MustNewOr(
					rule.NewGEFV("band", dlit.MustNew(4)),
					rule.NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
				),
				rule.MustNewOr(
					rule.NewGEFV("band", dlit.MustNew(4)),
					rule.NewEQFV("group", dlit.MustNew("a")),
				),
				rule.MustNewAnd(
					rule.NewEQFV("group", dlit.MustNew("a")),
					rule.NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
				),
				rule.MustNewOr(
					rule.NewEQFV("group", dlit.MustNew("a")),
					rule.NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
				),
			},
		},
		{in: []rule.Rule{
			rule.NewEQFV("team", dlit.MustNew("a")),
			rule.NewGEFV("band", dlit.MustNew(4)),
			rule.NewGEFV("band", dlit.MustNew(2)),
			rule.NewGEFV("flow", dlit.MustNew(6)),
		},
			want: []rule.Rule{
				rule.MustNewAnd(rule.NewGEFV("band", dlit.MustNew(2)),
					rule.NewGEFV("flow", dlit.MustNew(6))),
				rule.MustNewAnd(rule.NewGEFV("band", dlit.MustNew(2)),
					rule.NewEQFV("team", dlit.MustNew("a"))),
				rule.MustNewOr(rule.NewGEFV("band", dlit.MustNew(2)),
					rule.NewGEFV("flow", dlit.MustNew(6))),
				rule.MustNewOr(rule.NewGEFV("band", dlit.MustNew(2)),
					rule.NewEQFV("team", dlit.MustNew("a"))),
				rule.MustNewAnd(rule.NewGEFV("band", dlit.MustNew(4)),
					rule.NewGEFV("flow", dlit.MustNew(6))),
				rule.MustNewAnd(rule.NewGEFV("band", dlit.MustNew(4)),
					rule.NewEQFV("team", dlit.MustNew("a"))),
				rule.MustNewOr(rule.NewGEFV("band", dlit.MustNew(4)),
					rule.NewGEFV("flow", dlit.MustNew(6))),
				rule.MustNewOr(rule.NewGEFV("band", dlit.MustNew(4)),
					rule.NewEQFV("team", dlit.MustNew("a"))),
				rule.MustNewAnd(rule.NewGEFV("flow", dlit.MustNew(6)),
					rule.NewEQFV("team", dlit.MustNew("a"))),
				rule.MustNewOr(rule.NewGEFV("flow", dlit.MustNew(6)),
					rule.NewEQFV("team", dlit.MustNew("a"))),
			},
		},
		{in: []rule.Rule{
			rule.NewInFV("team", makeStringsDlitSlice("pink", "yellow", "blue")),
			rule.NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
		},
			want: []rule.Rule{
				rule.NewInFV("team",
					makeStringsDlitSlice("pink", "yellow", "blue", "red", "green")),
			},
		},
		{in: []rule.Rule{
			rule.NewInFV("team", makeStringsDlitSlice("pink", "yellow", "blue")),
			rule.NewInFV("group", makeStringsDlitSlice("red", "green", "blue")),
		},
			want: []rule.Rule{
				rule.MustNewAnd(
					rule.NewInFV("group", makeStringsDlitSlice("red", "green", "blue")),
					rule.NewInFV("team", makeStringsDlitSlice("pink", "yellow", "blue")),
				),
				rule.MustNewOr(
					rule.NewInFV("group", makeStringsDlitSlice("red", "green", "blue")),
					rule.NewInFV("team", makeStringsDlitSlice("pink", "yellow", "blue")),
				),
			},
		},
		{in: []rule.Rule{rule.NewEQFV("team", dlit.MustNew("a"))},
			want: []rule.Rule{}},
		{in: []rule.Rule{}, want: []rule.Rule{}},
	}

	for _, c := range cases {
		gotRules := CombineRules(c.in)
		if err := matchRulesUnordered(gotRules, c.want); err != nil {
			gotRuleStrs := rulesToSortedStrings(gotRules)
			wantRuleStrs := rulesToSortedStrings(c.want)
			t.Errorf("matchRulesUnordered() rules don't match: %s\n got: %s\n want: %s\n",
				err, gotRuleStrs, wantRuleStrs)
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
) error {
	if len(rules1) != len(rules2) {
		return errors.New("rules different lengths")
	}
	return rulesContain(rules1, rules2)
}

func rulesContain(gotRules []rule.Rule, wantRules []rule.Rule) error {
	for _, wRule := range wantRules {
		found := false
		for _, gRule := range gotRules {
			if gRule.String() == wRule.String() {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("rule doesn't exist: %s", wRule)
		}
	}
	return nil
}
