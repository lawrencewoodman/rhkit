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
	"strings"
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
				Kind:  fieldtype.Int,
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
				Kind:  fieldtype.Float,
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
				Kind:  fieldtype.Int,
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
		rule.NewEQFVS("team", "a"),
		rule.NewNEFVS("team", "a"),
		rule.NewEQFF("team", "teamOut"),
		rule.NewNEFF("team", "teamOut"),
		rule.NewInFV("teamOut", makeStringsDlitSlice("a", "d")),
		rule.NewEQFVI("level", 0),
		rule.NewEQFVI("level", 1),
		rule.NewNEFVI("level", 0),
		rule.NewNEFVI("level", 1),
		rule.NewLTFF("level", "position"),
		rule.NewLEFF("level", "position"),
		rule.NewNEFF("level", "position"),
		rule.NewGEFF("level", "position"),
		rule.NewGTFF("level", "position"),
		rule.NewEQFF("level", "position"),
		rule.NewGEFVI("level", 0),
		rule.NewGEFVI("level", 1),
		rule.NewLEFVI("level", 4),
		rule.NewLEFVI("level", 5),
		rule.NewInFV("level", makeStringsDlitSlice("0", "1")),
		rule.NewInFV("level", makeStringsDlitSlice("0", "3")),
		rule.NewGEFVF("flow", 2.1),
		rule.NewGEFVF("flow", 3.15),
		rule.NewLEFVF("flow", 4.2),
		rule.NewLEFVF("flow", 5.25),
		rule.NewAddLEF("level", "position", dlit.MustNew(12)),
		rule.NewAddGEF("level", "position", dlit.MustNew(12)),
	}

	got := GenerateRules(inputDescription, ruleFields)
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
	rules := GenerateRules(inputDescription, ruleFields)

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
		rule.NewEQFVS("directionIn", "de"),
		rule.NewEQFVS("directionIn", "gogledd"),
		rule.NewEQFVS("directionOut", "dwyrain"),
		rule.NewEQFVS("directionOut", "gorllewin"),
		rule.MustNewAnd(
			rule.NewEQFVS("directionIn", "de"),
			rule.NewEQFVS("directionOut", "dwyrain"),
		),
		rule.MustNewAnd(
			rule.NewEQFVS("directionIn", "de"),
			rule.NewEQFVS("directionOut", "gorllewin"),
		),
		rule.MustNewAnd(
			rule.NewEQFVS("directionIn", "gogledd"),
			rule.NewEQFVS("directionOut", "dwyrain"),
		),
		rule.MustNewAnd(
			rule.NewEQFVS("directionIn", "gogledd"),
			rule.NewEQFVS("directionOut", "gorllewin"),
		),
		rule.MustNewOr(
			rule.NewEQFVS("directionIn", "de"),
			rule.NewEQFVS("directionIn", "gogledd"),
		),
		rule.MustNewOr(
			rule.NewEQFVS("directionIn", "de"),
			rule.NewEQFVS("directionOut", "dwyrain"),
		),
		rule.MustNewOr(
			rule.NewEQFVS("directionIn", "de"),
			rule.NewEQFVS("directionOut", "gorllewin"),
		),
		rule.MustNewOr(
			rule.NewEQFVS("directionIn", "gogledd"),
			rule.NewEQFVS("directionOut", "dwyrain"),
		),
		rule.MustNewOr(
			rule.NewEQFVS("directionIn", "gogledd"),
			rule.NewEQFVS("directionOut", "gorllewin"),
		),
		rule.MustNewOr(
			rule.NewEQFVS("directionOut", "dwyrain"),
			rule.NewEQFVS("directionOut", "gorllewin"),
		),
		rule.NewTrue(),
	}

	got := GenerateRules(inputDescription, ruleFields)
	rule.Sort(got)
	rule.Sort(want)
	if err := matchRulesUnordered(got, want); err != nil {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("matchRulesUnordered: %s\n got: %s\nwant: %s\n",
			err, got, want)
	}
}

func TestGenerateIntRules(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"flow": &description.Field{
				Kind:   fieldtype.Int,
				Min:    dlit.MustNew(1),
				Max:    dlit.MustNew(20),
				Values: map[string]description.Value{},
			},
		},
	}
	ruleFields := []string{"flow"}
	wantRules := []rule.Rule{
		rule.NewLEFVI("flow", 2),
		rule.MustNewAnd(rule.NewLEFVI("flow", 2), rule.NewGEFVI("flow", 1)),
		rule.MustNewOr(rule.NewLEFVI("flow", 2), rule.NewGEFVI("flow", 3)),
		rule.MustNewOr(rule.NewLEFVI("flow", 2), rule.NewGEFVI("flow", 3)),
		rule.MustNewOr(rule.NewLEFVI("flow", 2), rule.NewGEFVI("flow", 4)),
		rule.MustNewOr(rule.NewLEFVI("flow", 2), rule.NewGEFVI("flow", 19)),
		rule.NewLEFVI("flow", 3),
		rule.MustNewAnd(rule.NewLEFVI("flow", 3), rule.NewGEFVI("flow", 1)),
		rule.MustNewAnd(rule.NewLEFVI("flow", 3), rule.NewGEFVI("flow", 2)),
		rule.MustNewOr(rule.NewLEFVI("flow", 3), rule.NewGEFVI("flow", 4)),
		rule.MustNewOr(rule.NewLEFVI("flow", 3), rule.NewGEFVI("flow", 19)),
		rule.NewLEFVI("flow", 4),
		rule.MustNewAnd(rule.NewLEFVI("flow", 4), rule.NewGEFVI("flow", 1)),
		rule.MustNewAnd(rule.NewLEFVI("flow", 4), rule.NewGEFVI("flow", 2)),
		rule.MustNewAnd(rule.NewLEFVI("flow", 4), rule.NewGEFVI("flow", 3)),
		rule.MustNewOr(rule.NewLEFVI("flow", 4), rule.NewGEFVI("flow", 5)),
		rule.MustNewOr(rule.NewLEFVI("flow", 4), rule.NewGEFVI("flow", 19)),
		rule.NewLEFVI("flow", 20),

		rule.NewGEFVI("flow", 1), /* TODO: This is a pointless rule */

		rule.NewGEFVI("flow", 2),
		rule.NewGEFVI("flow", 3),
		rule.NewGEFVI("flow", 19),
	}

	got := generateIntRules(inputDescription, ruleFields, "flow")
	if err := rulesContain(got, wantRules); err != nil {
		t.Errorf("GenerateIntRules: %s, got: %s", err, got)
	}
	if len(got) <= len(wantRules) {
		t.Errorf("GenerateIntRules: There should be more rules generated: %d",
			len(got))
	}
}

func TestGenerateFloatRules(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"flow": &description.Field{
				Kind:   fieldtype.Float,
				Min:    dlit.MustNew(0),
				Max:    dlit.MustNew(10.5),
				MaxDP:  2,
				Values: map[string]description.Value{},
			},
			"flowConstant": &description.Field{
				Kind:   fieldtype.Float,
				Min:    dlit.MustNew(10.5),
				Max:    dlit.MustNew(10.5),
				MaxDP:  1,
				Values: map[string]description.Value{},
			},
		},
	}
	ruleFields := []string{"flow", "flowConstant"}
	cases := []struct {
		field string
		want  []rule.Rule
	}{
		{field: "flow",
			want: []rule.Rule{
				rule.NewLEFVF("flow", 1.05),
				rule.MustNewAnd(rule.NewLEFVF("flow", 1.05), rule.NewGEFVF("flow", 0.0)),
				rule.MustNewOr(rule.NewLEFVF("flow", 1.05), rule.NewGEFVF("flow", 2.1)),
				rule.MustNewOr(rule.NewLEFVF("flow", 1.05), rule.NewGEFVF("flow", 3.15)),
				rule.MustNewOr(rule.NewLEFVF("flow", 1.05), rule.NewGEFVF("flow", 4.2)),
				rule.MustNewOr(rule.NewLEFVF("flow", 1.05), rule.NewGEFVF("flow", 5.25)),
				rule.MustNewOr(rule.NewLEFVF("flow", 1.05), rule.NewGEFVF("flow", 6.3)),
				rule.MustNewOr(rule.NewLEFVF("flow", 1.05), rule.NewGEFVF("flow", 7.35)),
				rule.MustNewOr(rule.NewLEFVF("flow", 1.05), rule.NewGEFVF("flow", 8.4)),
				rule.MustNewOr(rule.NewLEFVF("flow", 1.05), rule.NewGEFVF("flow", 9.45)),
				rule.NewLEFVF("flow", 2.1),
				rule.MustNewAnd(rule.NewLEFVF("flow", 2.1), rule.NewGEFVF("flow", 0.0)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 2.1), rule.NewGEFVF("flow", 1.05)),
				rule.MustNewOr(rule.NewLEFVF("flow", 2.1), rule.NewGEFVF("flow", 3.15)),
				rule.MustNewOr(rule.NewLEFVF("flow", 2.1), rule.NewGEFVF("flow", 4.2)),
				rule.MustNewOr(rule.NewLEFVF("flow", 2.1), rule.NewGEFVF("flow", 5.25)),
				rule.MustNewOr(rule.NewLEFVF("flow", 2.1), rule.NewGEFVF("flow", 6.3)),
				rule.MustNewOr(rule.NewLEFVF("flow", 2.1), rule.NewGEFVF("flow", 7.35)),
				rule.MustNewOr(rule.NewLEFVF("flow", 2.1), rule.NewGEFVF("flow", 8.4)),
				rule.MustNewOr(rule.NewLEFVF("flow", 2.1), rule.NewGEFVF("flow", 9.45)),
				rule.NewLEFVF("flow", 3.15),
				rule.MustNewAnd(rule.NewLEFVF("flow", 3.15), rule.NewGEFVF("flow", 0.0)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 3.15), rule.NewGEFVF("flow", 1.05)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 3.15), rule.NewGEFVF("flow", 2.1)),
				rule.MustNewOr(rule.NewLEFVF("flow", 3.15), rule.NewGEFVF("flow", 4.2)),
				rule.MustNewOr(rule.NewLEFVF("flow", 3.15), rule.NewGEFVF("flow", 5.25)),
				rule.MustNewOr(rule.NewLEFVF("flow", 3.15), rule.NewGEFVF("flow", 6.3)),
				rule.MustNewOr(rule.NewLEFVF("flow", 3.15), rule.NewGEFVF("flow", 7.35)),
				rule.MustNewOr(rule.NewLEFVF("flow", 3.15), rule.NewGEFVF("flow", 8.4)),
				rule.MustNewOr(rule.NewLEFVF("flow", 3.15), rule.NewGEFVF("flow", 9.45)),
				rule.NewLEFVF("flow", 4.2),
				rule.MustNewAnd(rule.NewLEFVF("flow", 4.2), rule.NewGEFVF("flow", 0.0)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 4.2), rule.NewGEFVF("flow", 1.05)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 4.2), rule.NewGEFVF("flow", 2.1)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 4.2), rule.NewGEFVF("flow", 3.15)),
				rule.MustNewOr(rule.NewLEFVF("flow", 4.2), rule.NewGEFVF("flow", 5.25)),
				rule.MustNewOr(rule.NewLEFVF("flow", 4.2), rule.NewGEFVF("flow", 6.3)),
				rule.MustNewOr(rule.NewLEFVF("flow", 4.2), rule.NewGEFVF("flow", 7.35)),
				rule.MustNewOr(rule.NewLEFVF("flow", 4.2), rule.NewGEFVF("flow", 8.4)),
				rule.MustNewOr(rule.NewLEFVF("flow", 4.2), rule.NewGEFVF("flow", 9.45)),
				rule.NewLEFVF("flow", 5.25),
				rule.MustNewAnd(rule.NewLEFVF("flow", 5.25), rule.NewGEFVF("flow", 0.0)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 5.25), rule.NewGEFVF("flow", 1.05)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 5.25), rule.NewGEFVF("flow", 2.1)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 5.25), rule.NewGEFVF("flow", 3.15)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 5.25), rule.NewGEFVF("flow", 4.2)),
				rule.MustNewOr(rule.NewLEFVF("flow", 5.25), rule.NewGEFVF("flow", 6.3)),
				rule.MustNewOr(rule.NewLEFVF("flow", 5.25), rule.NewGEFVF("flow", 7.35)),
				rule.MustNewOr(rule.NewLEFVF("flow", 5.25), rule.NewGEFVF("flow", 8.4)),
				rule.MustNewOr(rule.NewLEFVF("flow", 5.25), rule.NewGEFVF("flow", 9.45)),
				rule.NewLEFVF("flow", 6.3),
				rule.MustNewAnd(rule.NewLEFVF("flow", 6.3), rule.NewGEFVF("flow", 0.0)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 6.3), rule.NewGEFVF("flow", 1.05)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 6.3), rule.NewGEFVF("flow", 2.1)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 6.3), rule.NewGEFVF("flow", 3.15)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 6.3), rule.NewGEFVF("flow", 4.2)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 6.3), rule.NewGEFVF("flow", 5.25)),
				rule.MustNewOr(rule.NewLEFVF("flow", 6.3), rule.NewGEFVF("flow", 7.35)),
				rule.MustNewOr(rule.NewLEFVF("flow", 6.3), rule.NewGEFVF("flow", 8.4)),
				rule.MustNewOr(rule.NewLEFVF("flow", 6.3), rule.NewGEFVF("flow", 9.45)),
				rule.NewLEFVF("flow", 7.35),
				rule.MustNewAnd(rule.NewLEFVF("flow", 7.35), rule.NewGEFVF("flow", 0.0)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 7.35), rule.NewGEFVF("flow", 1.05)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 7.35), rule.NewGEFVF("flow", 2.1)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 7.35), rule.NewGEFVF("flow", 3.15)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 7.35), rule.NewGEFVF("flow", 4.2)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 7.35), rule.NewGEFVF("flow", 5.25)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 7.35), rule.NewGEFVF("flow", 6.3)),
				rule.MustNewOr(rule.NewLEFVF("flow", 7.35), rule.NewGEFVF("flow", 8.4)),
				rule.MustNewOr(rule.NewLEFVF("flow", 7.35), rule.NewGEFVF("flow", 9.45)),
				rule.NewLEFVF("flow", 8.4),
				rule.MustNewAnd(rule.NewLEFVF("flow", 8.4), rule.NewGEFVF("flow", 0.0)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 8.4), rule.NewGEFVF("flow", 1.05)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 8.4), rule.NewGEFVF("flow", 2.1)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 8.4), rule.NewGEFVF("flow", 3.15)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 8.4), rule.NewGEFVF("flow", 4.2)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 8.4), rule.NewGEFVF("flow", 5.25)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 8.4), rule.NewGEFVF("flow", 6.3)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 8.4), rule.NewGEFVF("flow", 7.35)),
				rule.MustNewOr(rule.NewLEFVF("flow", 8.4), rule.NewGEFVF("flow", 9.45)),
				rule.NewLEFVF("flow", 9.45),
				rule.MustNewAnd(rule.NewLEFVF("flow", 9.45), rule.NewGEFVF("flow", 0.0)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 9.45), rule.NewGEFVF("flow", 1.05)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 9.45), rule.NewGEFVF("flow", 2.1)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 9.45), rule.NewGEFVF("flow", 3.15)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 9.45), rule.NewGEFVF("flow", 4.2)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 9.45), rule.NewGEFVF("flow", 5.25)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 9.45), rule.NewGEFVF("flow", 6.3)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 9.45), rule.NewGEFVF("flow", 7.35)),
				rule.MustNewAnd(rule.NewLEFVF("flow", 9.45), rule.NewGEFVF("flow", 8.4)),
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
			},
		},
		{field: "flowConstant", want: []rule.Rule{}},
	}

	for _, c := range cases {
		got := generateFloatRules(inputDescription, ruleFields, c.field)
		if err := rulesContain(got, c.want); err != nil {
			t.Errorf("GenerateFloatRules: %s", err)
		}
		if len(got) < len(c.want) {
			t.Errorf(
				"GenerateFloatRules: There should be more rules generated for field: %s, got: %d",
				c.field, len(got),
			)
		}
	}
}

func TestGenerateAddRules(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"flowIn": &description.Field{
				Kind:   fieldtype.Float,
				Min:    dlit.MustNew(0),
				Max:    dlit.MustNew(10.5),
				MaxDP:  2,
				Values: map[string]description.Value{},
			},
			"flowOut": &description.Field{
				Kind:   fieldtype.Float,
				Min:    dlit.MustNew(1.2),
				Max:    dlit.MustNew(5.2),
				MaxDP:  2,
				Values: map[string]description.Value{},
			},
			"flowSideA": &description.Field{
				Kind:   fieldtype.Float,
				Min:    dlit.MustNew(2.2),
				Max:    dlit.MustNew(2.2),
				MaxDP:  1,
				Values: map[string]description.Value{},
			},
			"flowSideB": &description.Field{
				Kind:   fieldtype.Float,
				Min:    dlit.MustNew(2.2),
				Max:    dlit.MustNew(2.2),
				MaxDP:  1,
				Values: map[string]description.Value{},
			},
		},
	}
	cases := []struct {
		field string
		want  []rule.Rule
	}{
		{field: "flowIn",
			want: []rule.Rule{
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(1.92)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(2.65)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(3.38)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(4.1)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(4.83)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(5.55)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(6.27)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(7)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(7.72)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(8.45)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(9.17)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(9.9)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(10.62)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(11.35)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(12.07)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(12.8)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(13.52)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(14.25)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(14.97)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(15.7)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(1.2)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(1.92)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(2.65)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(3.38)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(4.1)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(4.83)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(5.55)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(6.27)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(7)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(7.72)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(8.45)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(9.17)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(9.9)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(10.62)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(11.35)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(12.07)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(12.8)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(13.52)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(14.25)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(14.97)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(15.7)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(2.73)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(3.25)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(3.78)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(4.3)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(4.83)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(5.35)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(5.88)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(6.4)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(6.93)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(7.45)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(7.98)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(8.5)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(9.03)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(9.55)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(10.08)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(10.6)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(11.13)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(11.65)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(12.18)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(2.2)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(2.73)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(3.25)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(3.78)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(4.3)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(4.83)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(5.35)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(5.88)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(6.4)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(6.93)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(7.45)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(7.98)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(8.5)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(9.03)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(9.55)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(10.08)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(10.6)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(11.13)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(11.65)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(12.18)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(2.73)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(3.25)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(3.78)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(4.3)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(4.83)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(5.35)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(5.88)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(6.4)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(6.93)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(7.45)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(7.98)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(8.5)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(9.03)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(9.55)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(10.08)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(10.6)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(11.13)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(11.65)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(12.18)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(2.2)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(2.73)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(3.25)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(3.78)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(4.3)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(4.83)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(5.35)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(5.88)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(6.4)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(6.93)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(7.45)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(7.98)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(8.5)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(9.03)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(9.55)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(10.08)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(10.6)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(11.13)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(11.65)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(12.18)),
			},
		},
		{field: "flowOut",
			want: []rule.Rule{
				rule.NewAddLEF("flowOut", "flowSideA", dlit.MustNew(3.6)),
				rule.NewAddLEF("flowOut", "flowSideA", dlit.MustNew(3.8)),
				rule.NewAddLEF("flowOut", "flowSideA", dlit.MustNew(4)),
				rule.NewAddLEF("flowOut", "flowSideA", dlit.MustNew(4.2)),
				rule.NewAddLEF("flowOut", "flowSideA", dlit.MustNew(4.4)),
				rule.NewAddLEF("flowOut", "flowSideA", dlit.MustNew(4.6)),
				rule.NewAddLEF("flowOut", "flowSideA", dlit.MustNew(4.8)),
				rule.NewAddLEF("flowOut", "flowSideA", dlit.MustNew(5)),
				rule.NewAddLEF("flowOut", "flowSideA", dlit.MustNew(5.2)),
				rule.NewAddLEF("flowOut", "flowSideA", dlit.MustNew(5.4)),
				rule.NewAddLEF("flowOut", "flowSideA", dlit.MustNew(5.6)),
				rule.NewAddLEF("flowOut", "flowSideA", dlit.MustNew(5.8)),
				rule.NewAddLEF("flowOut", "flowSideA", dlit.MustNew(6)),
				rule.NewAddLEF("flowOut", "flowSideA", dlit.MustNew(6.2)),
				rule.NewAddLEF("flowOut", "flowSideA", dlit.MustNew(6.4)),
				rule.NewAddLEF("flowOut", "flowSideA", dlit.MustNew(6.6)),
				rule.NewAddLEF("flowOut", "flowSideA", dlit.MustNew(6.8)),
				rule.NewAddLEF("flowOut", "flowSideA", dlit.MustNew(7)),
				rule.NewAddLEF("flowOut", "flowSideA", dlit.MustNew(7.2)),
				rule.NewAddGEF("flowOut", "flowSideA", dlit.MustNew(3.4)),
				rule.NewAddGEF("flowOut", "flowSideA", dlit.MustNew(3.6)),
				rule.NewAddGEF("flowOut", "flowSideA", dlit.MustNew(3.8)),
				rule.NewAddGEF("flowOut", "flowSideA", dlit.MustNew(4)),
				rule.NewAddGEF("flowOut", "flowSideA", dlit.MustNew(4.2)),
				rule.NewAddGEF("flowOut", "flowSideA", dlit.MustNew(4.4)),
				rule.NewAddGEF("flowOut", "flowSideA", dlit.MustNew(4.6)),
				rule.NewAddGEF("flowOut", "flowSideA", dlit.MustNew(4.8)),
				rule.NewAddGEF("flowOut", "flowSideA", dlit.MustNew(5)),
				rule.NewAddGEF("flowOut", "flowSideA", dlit.MustNew(5.2)),
				rule.NewAddGEF("flowOut", "flowSideA", dlit.MustNew(5.4)),
				rule.NewAddGEF("flowOut", "flowSideA", dlit.MustNew(5.6)),
				rule.NewAddGEF("flowOut", "flowSideA", dlit.MustNew(5.8)),
				rule.NewAddGEF("flowOut", "flowSideA", dlit.MustNew(6)),
				rule.NewAddGEF("flowOut", "flowSideA", dlit.MustNew(6.2)),
				rule.NewAddGEF("flowOut", "flowSideA", dlit.MustNew(6.4)),
				rule.NewAddGEF("flowOut", "flowSideA", dlit.MustNew(6.6)),
				rule.NewAddGEF("flowOut", "flowSideA", dlit.MustNew(6.8)),
				rule.NewAddGEF("flowOut", "flowSideA", dlit.MustNew(7)),
				rule.NewAddGEF("flowOut", "flowSideA", dlit.MustNew(7.2)),
				rule.NewAddLEF("flowOut", "flowSideB", dlit.MustNew(3.6)),
				rule.NewAddLEF("flowOut", "flowSideB", dlit.MustNew(3.8)),
				rule.NewAddLEF("flowOut", "flowSideB", dlit.MustNew(4)),
				rule.NewAddLEF("flowOut", "flowSideB", dlit.MustNew(4.2)),
				rule.NewAddLEF("flowOut", "flowSideB", dlit.MustNew(4.4)),
				rule.NewAddLEF("flowOut", "flowSideB", dlit.MustNew(4.6)),
				rule.NewAddLEF("flowOut", "flowSideB", dlit.MustNew(4.8)),
				rule.NewAddLEF("flowOut", "flowSideB", dlit.MustNew(5)),
				rule.NewAddLEF("flowOut", "flowSideB", dlit.MustNew(5.2)),
				rule.NewAddLEF("flowOut", "flowSideB", dlit.MustNew(5.4)),
				rule.NewAddLEF("flowOut", "flowSideB", dlit.MustNew(5.6)),
				rule.NewAddLEF("flowOut", "flowSideB", dlit.MustNew(5.8)),
				rule.NewAddLEF("flowOut", "flowSideB", dlit.MustNew(6)),
				rule.NewAddLEF("flowOut", "flowSideB", dlit.MustNew(6.2)),
				rule.NewAddLEF("flowOut", "flowSideB", dlit.MustNew(6.4)),
				rule.NewAddLEF("flowOut", "flowSideB", dlit.MustNew(6.6)),
				rule.NewAddLEF("flowOut", "flowSideB", dlit.MustNew(6.8)),
				rule.NewAddLEF("flowOut", "flowSideB", dlit.MustNew(7)),
				rule.NewAddLEF("flowOut", "flowSideB", dlit.MustNew(7.2)),
				rule.NewAddGEF("flowOut", "flowSideB", dlit.MustNew(3.4)),
				rule.NewAddGEF("flowOut", "flowSideB", dlit.MustNew(3.6)),
				rule.NewAddGEF("flowOut", "flowSideB", dlit.MustNew(3.8)),
				rule.NewAddGEF("flowOut", "flowSideB", dlit.MustNew(4)),
				rule.NewAddGEF("flowOut", "flowSideB", dlit.MustNew(4.2)),
				rule.NewAddGEF("flowOut", "flowSideB", dlit.MustNew(4.4)),
				rule.NewAddGEF("flowOut", "flowSideB", dlit.MustNew(4.6)),
				rule.NewAddGEF("flowOut", "flowSideB", dlit.MustNew(4.8)),
				rule.NewAddGEF("flowOut", "flowSideB", dlit.MustNew(5)),
				rule.NewAddGEF("flowOut", "flowSideB", dlit.MustNew(5.2)),
				rule.NewAddGEF("flowOut", "flowSideB", dlit.MustNew(5.4)),
				rule.NewAddGEF("flowOut", "flowSideB", dlit.MustNew(5.6)),
				rule.NewAddGEF("flowOut", "flowSideB", dlit.MustNew(5.8)),
				rule.NewAddGEF("flowOut", "flowSideB", dlit.MustNew(6)),
				rule.NewAddGEF("flowOut", "flowSideB", dlit.MustNew(6.2)),
				rule.NewAddGEF("flowOut", "flowSideB", dlit.MustNew(6.4)),
				rule.NewAddGEF("flowOut", "flowSideB", dlit.MustNew(6.6)),
				rule.NewAddGEF("flowOut", "flowSideB", dlit.MustNew(6.8)),
				rule.NewAddGEF("flowOut", "flowSideB", dlit.MustNew(7)),
				rule.NewAddGEF("flowOut", "flowSideB", dlit.MustNew(7.2)),
			},
		},
		{field: "flowSideA",
			want: []rule.Rule{},
		},
		{field: "flowSideB",
			want: []rule.Rule{},
		},
	}
	ruleFields := []string{"flowIn", "flowOut", "flowSideA", "flowSideB"}

	for _, c := range cases {
		got := generateAddRules(inputDescription, ruleFields, c.field)
		if err := matchRulesUnordered(got, c.want); err != nil {
			gotRuleStrs := rulesToSortedStrings(got)
			wantRuleStrs := rulesToSortedStrings(c.want)
			t.Errorf("matchRulesUnordered() rules don't match: %s\ngot: %s\nwant: %s\n",
				err, gotRuleStrs, wantRuleStrs)
		}
	}
}

func TestGenerateCompareNumericRules(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"band": &description.Field{
				Kind:   fieldtype.Int,
				Min:    dlit.MustNew(1),
				Max:    dlit.MustNew(3),
				Values: map[string]description.Value{},
			},
			"flowIn": &description.Field{
				Kind:   fieldtype.Float,
				Min:    dlit.MustNew(1),
				Max:    dlit.MustNew(4),
				MaxDP:  2,
				Values: map[string]description.Value{},
			},
			"flowOut": &description.Field{
				Kind:   fieldtype.Float,
				Min:    dlit.MustNew(0.95),
				Max:    dlit.MustNew(4.1),
				MaxDP:  2,
				Values: map[string]description.Value{},
			},
			"rateIn": &description.Field{
				Kind:   fieldtype.Float,
				Min:    dlit.MustNew(4.2),
				Max:    dlit.MustNew(8.9),
				MaxDP:  2,
				Values: map[string]description.Value{},
			},
			"rateOut": &description.Field{
				Kind:   fieldtype.Float,
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
				rule.NewNEFF("band", "flowIn"),
				rule.NewNEFF("band", "flowOut"),
				rule.NewLTFF("band", "flowIn"),
				rule.NewLTFF("band", "flowOut"),
				rule.NewLEFF("band", "flowIn"),
				rule.NewLEFF("band", "flowOut"),
				rule.NewEQFF("band", "flowIn"),
				rule.NewEQFF("band", "flowOut"),
				rule.NewGTFF("band", "flowIn"),
				rule.NewGTFF("band", "flowOut"),
				rule.NewGEFF("band", "flowIn"),
				rule.NewGEFF("band", "flowOut"),
			},
		},
		{field: "flowIn",
			want: []rule.Rule{
				rule.NewNEFF("flowIn", "flowOut"),
				rule.NewLTFF("flowIn", "flowOut"),
				rule.NewLEFF("flowIn", "flowOut"),
				rule.NewEQFF("flowIn", "flowOut"),
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
	for _, c := range cases {
		got := generateCompareNumericRules(inputDescription, ruleFields, c.field)
		if err := matchRulesUnordered(got, c.want); err != nil {
			gotRuleStrs := rulesToSortedStrings(got)
			wantRuleStrs := rulesToSortedStrings(c.want)
			t.Errorf("matchRulesUnordered() rules don't match: %s\ngot: %s\nwant: %s\n",
				err, gotRuleStrs, wantRuleStrs)
		}
	}
}

func TestGenerateCompareStringRules(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"band": &description.Field{
				Kind:   fieldtype.Int,
				Min:    dlit.MustNew(1),
				Max:    dlit.MustNew(3),
				Values: map[string]description.Value{},
			},
			"groupA": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Nelson":      description.Value{dlit.NewString("Nelson"), 3},
					"Collingwood": description.Value{dlit.NewString("Collingwood"), 1},
					"Mountbatten": description.Value{dlit.NewString("Mountbatten"), 1},
					"Drake":       description.Value{dlit.NewString("Drake"), 2},
				},
			},
			"groupB": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Nelson":      description.Value{dlit.NewString("Nelson"), 3},
					"Mountbatten": description.Value{dlit.NewString("Mountbatten"), 1},
					"Drake":       description.Value{dlit.NewString("Drake"), 2},
				},
			},
			"groupC": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Nelson": description.Value{dlit.NewString("Nelson"), 3},
					"Drake":  description.Value{dlit.NewString("Drake"), 2},
				},
			},
			"groupD": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Drake": description.Value{dlit.NewString("Drake"), 2},
				},
			},
			"groupE": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Drake":       description.Value{dlit.NewString("Drake"), 2},
					"Chaucer":     description.Value{dlit.NewString("Chaucer"), 2},
					"Shakespeare": description.Value{dlit.NewString("Shakespeare"), 2},
					"Marlowe":     description.Value{dlit.NewString("Marlowe"), 2},
				},
			},
			"groupF": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Nelson":      description.Value{dlit.NewString("Nelson"), 3},
					"Drake":       description.Value{dlit.NewString("Drake"), 2},
					"Chaucer":     description.Value{dlit.NewString("Chaucer"), 2},
					"Shakespeare": description.Value{dlit.NewString("Shakespeare"), 2},
					"Marlowe":     description.Value{dlit.NewString("Marlowe"), 2},
				},
			},
		},
	}
	cases := []struct {
		field string
		want  []rule.Rule
	}{
		{field: "band",
			want: []rule.Rule{},
		},
		{field: "groupA",
			want: []rule.Rule{
				rule.NewEQFF("groupA", "groupB"),
				rule.NewNEFF("groupA", "groupB"),
				rule.NewEQFF("groupA", "groupC"),
				rule.NewNEFF("groupA", "groupC"),
				rule.NewEQFF("groupA", "groupF"),
				rule.NewNEFF("groupA", "groupF"),
			},
		},
		{field: "groupB",
			want: []rule.Rule{
				rule.NewEQFF("groupB", "groupC"),
				rule.NewNEFF("groupB", "groupC"),
				rule.NewEQFF("groupB", "groupF"),
				rule.NewNEFF("groupB", "groupF"),
			},
		},
		{field: "groupC",
			want: []rule.Rule{
				rule.NewEQFF("groupC", "groupF"),
				rule.NewNEFF("groupC", "groupF"),
			},
		},
		{field: "groupD",
			want: []rule.Rule{},
		},
		{field: "groupE",
			want: []rule.Rule{
				rule.NewEQFF("groupE", "groupF"),
				rule.NewNEFF("groupE", "groupF"),
			},
		},
	}
	ruleFields :=
		[]string{"band", "groupA", "groupB", "groupC", "groupD", "groupE", "groupF"}
	for _, c := range cases {
		got := generateCompareStringRules(inputDescription, ruleFields, c.field)
		if err := matchRulesUnordered(got, c.want); err != nil {
			gotRuleStrs := rulesToSortedStrings(got)
			wantRuleStrs := rulesToSortedStrings(c.want)
			t.Errorf("matchRulesUnordered() rules don't match: %s\ngot: %s\nwant: %s\n",
				err, gotRuleStrs, wantRuleStrs)
		}
	}
}

func TestGenerateInRules_1(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"band": &description.Field{
				Kind:   fieldtype.Int,
				Min:    dlit.MustNew(1),
				Max:    dlit.MustNew(3),
				Values: map[string]description.Value{},
			},
			"flow": &description.Field{
				Kind:   fieldtype.Float,
				Min:    dlit.MustNew(1),
				Max:    dlit.MustNew(3),
				MaxDP:  2,
				Values: map[string]description.Value{},
			},
			"groupA": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Fred":    description.Value{dlit.NewString("Fred"), 3},
					"Mary":    description.Value{dlit.NewString("Mary"), 4},
					"Rebecca": description.Value{dlit.NewString("Rebecca"), 2},
				},
			},

			"groupB": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Fred":    description.Value{dlit.NewString("Fred"), 3},
					"Mary":    description.Value{dlit.NewString("Mary"), 4},
					"Rebecca": description.Value{dlit.NewString("Rebecca"), 2},
					"Harry":   description.Value{dlit.NewString("Harry"), 2},
					"Dinah":   description.Value{dlit.NewString("Dinah"), 2},
					"Israel":  description.Value{dlit.NewString("Israel"), 2},
					"Sarah":   description.Value{dlit.NewString("Sarah"), 2},
					"Ishmael": description.Value{dlit.NewString("Ishmael"), 2},
					"Caen":    description.Value{dlit.NewString("Caen"), 2},
					"Abel":    description.Value{dlit.NewString("Abel"), 2},
					"Noah":    description.Value{dlit.NewString("Noah"), 2},
					"Isaac":   description.Value{dlit.NewString("Isaac"), 2},
					"Moses":   description.Value{dlit.NewString("Moses"), 2},
				},
			},
			"groupC": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Fred":    description.Value{dlit.NewString("Fred"), 3},
					"Mary":    description.Value{dlit.NewString("Mary"), 4},
					"Rebecca": description.Value{dlit.NewString("Rebecca"), 2},
					"Harry":   description.Value{dlit.NewString("Harry"), 2},
				},
			},
			"groupD": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Fred":    description.Value{dlit.NewString("Fred"), 3},
					"Mary":    description.Value{dlit.NewString("Mary"), 4},
					"Rebecca": description.Value{dlit.NewString("Rebecca"), 1},
					"Harry":   description.Value{dlit.NewString("Harry"), 2},
				},
			},
			"groupE": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Fred":    description.Value{dlit.NewString("Fred"), 3},
					"Mary":    description.Value{dlit.NewString("Mary"), 4},
					"Rebecca": description.Value{dlit.NewString("Rebecca"), 2},
					"Harry":   description.Value{dlit.NewString("Harry"), 2},
					"Juliet":  description.Value{dlit.NewString("Juliet"), 2},
				},
			},
		},
	}
	cases := []struct {
		field string
		want  []rule.Rule
	}{
		{field: "band",
			want: []rule.Rule{},
		},
		{field: "flow",
			want: []rule.Rule{},
		},
		{field: "groupA",
			want: []rule.Rule{},
		},
		{field: "groupB",
			want: []rule.Rule{},
		},
		{field: "groupC",
			want: []rule.Rule{
				rule.NewInFV("groupC", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Harry"),
				}),
				rule.NewInFV("groupC", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Mary"),
				}),
				rule.NewInFV("groupC", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Rebecca"),
				}),
				rule.NewInFV("groupC", []*dlit.Literal{
					dlit.NewString("Harry"),
					dlit.NewString("Mary"),
				}),
				rule.NewInFV("groupC", []*dlit.Literal{
					dlit.NewString("Harry"),
					dlit.NewString("Rebecca"),
				}),
				rule.NewInFV("groupC", []*dlit.Literal{
					dlit.NewString("Mary"),
					dlit.NewString("Rebecca"),
				}),
			},
		},
		{field: "groupD",
			want: []rule.Rule{
				rule.NewInFV("groupD", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Harry"),
				}),
				rule.NewInFV("groupD", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Mary"),
				}),
				rule.NewInFV("groupD", []*dlit.Literal{
					dlit.NewString("Harry"),
					dlit.NewString("Mary"),
				}),
			},
		},
		{field: "groupE",
			want: []rule.Rule{
				rule.NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Harry"),
				}),
				rule.NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Harry"),
					dlit.NewString("Juliet"),
				}),
				rule.NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Harry"),
					dlit.NewString("Mary"),
				}),
				rule.NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Harry"),
					dlit.NewString("Rebecca"),
				}),
				rule.NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Juliet"),
					dlit.NewString("Rebecca"),
				}),
				rule.NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Juliet"),
					dlit.NewString("Mary"),
				}),
				rule.NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Mary"),
				}),
				rule.NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Mary"),
					dlit.NewString("Rebecca"),
				}),
				rule.NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Rebecca"),
				}),
				rule.NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Juliet"),
				}),
				rule.NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Harry"),
					dlit.NewString("Mary"),
				}),
				rule.NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Harry"),
					dlit.NewString("Rebecca"),
				}),
				rule.NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Harry"),
					dlit.NewString("Juliet"),
				}),
				rule.NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Harry"),
					dlit.NewString("Juliet"),
					dlit.NewString("Mary"),
				}),
				rule.NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Harry"),
					dlit.NewString("Juliet"),
					dlit.NewString("Rebecca"),
				}),
				rule.NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Harry"),
					dlit.NewString("Mary"),
					dlit.NewString("Rebecca"),
				}),
				rule.NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Juliet"),
					dlit.NewString("Mary"),
				}),
				rule.NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Juliet"),
					dlit.NewString("Mary"),
					dlit.NewString("Rebecca"),
				}),
				rule.NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Mary"),
					dlit.NewString("Rebecca"),
				}),
				rule.NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Juliet"),
					dlit.NewString("Rebecca"),
				}),
			},
		},
	}
	ruleFields :=
		[]string{"band", "flow", "groupA", "groupB", "groupC", "groupD", "groupE"}
	for _, c := range cases {
		got := generateInRules(inputDescription, ruleFields, c.field)
		if err := matchRulesUnordered(got, c.want); err != nil {
			gotRuleStrs := rulesToSortedStrings(got)
			wantRuleStrs := rulesToSortedStrings(c.want)
			t.Errorf("matchRulesUnordered() rules don't match: %s\ngot: %s\nwant: %s\n",
				err, gotRuleStrs, wantRuleStrs)
		}
	}
}

// Test that will generate if has 12 values and ensures that has correct
// number of values in In rule
func TestGenerateInRules_2(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"group": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Fred":    description.Value{dlit.NewString("Fred"), 3},
					"Mary":    description.Value{dlit.NewString("Mary"), 4},
					"Rebecca": description.Value{dlit.NewString("Rebecca"), 2},
					"Harry":   description.Value{dlit.NewString("Harry"), 2},
					"Dinah":   description.Value{dlit.NewString("Dinah"), 2},
					"Israel":  description.Value{dlit.NewString("Israel"), 2},
					"Sarah":   description.Value{dlit.NewString("Sarah"), 2},
					"Ishmael": description.Value{dlit.NewString("Ishmael"), 2},
					"Caen":    description.Value{dlit.NewString("Caen"), 2},
					"Abel":    description.Value{dlit.NewString("Abel"), 2},
					"Noah":    description.Value{dlit.NewString("Noah"), 2},
					"Isaac":   description.Value{dlit.NewString("Isaac"), 2},
				},
			},
		},
	}
	ruleFields := []string{"group"}
	got := generateInRules(inputDescription, ruleFields, "group")
	if len(got) < 1000 {
		t.Errorf("generateInRules: got too few rules: %d", len(got))
	}
	for _, r := range got {
		numValues := strings.Count(r.String(), ",")
		if numValues < 2 || numValues > 5 {
			t.Errorf("generateInRules: wrong number of values in rule: %s", r)
		}
	}
}

func TestGenerateValueRules(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"band": &description.Field{
				Kind: fieldtype.Int,
				Min:  dlit.MustNew(1),
				Max:  dlit.MustNew(4),
				Values: map[string]description.Value{
					"1": description.Value{dlit.NewString("1"), 3},
					"2": description.Value{dlit.NewString("2"), 1},
					"3": description.Value{dlit.NewString("3"), 2},
					"4": description.Value{dlit.NewString("4"), 5},
				},
			},
			"flow": &description.Field{
				Kind:  fieldtype.Float,
				Min:   dlit.MustNew(1),
				Max:   dlit.MustNew(4),
				MaxDP: 2,
				Values: map[string]description.Value{
					"1":    description.Value{dlit.NewString("1"), 3},
					"2":    description.Value{dlit.NewString("2"), 1},
					"2.90": description.Value{dlit.NewString("2.90"), 1},
					"3.37": description.Value{dlit.NewString("3.37"), 2},
					"3.3":  description.Value{dlit.NewString("3.3"), 2},
					"4":    description.Value{dlit.NewString("4"), 5},
				},
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
				rule.NewEQFVI("band", 1),
				rule.NewNEFVI("band", 1),
				rule.NewEQFVI("band", 3),
				rule.NewNEFVI("band", 3),
				rule.NewEQFVI("band", 4),
				rule.NewNEFVI("band", 4),
			},
		},
		{field: "flow",
			want: []rule.Rule{
				rule.NewEQFVF("flow", 1),
				rule.NewNEFVF("flow", 1),
				rule.NewEQFVF("flow", 3.37),
				rule.NewNEFVF("flow", 3.37),
				rule.NewEQFVF("flow", 3.3),
				rule.NewNEFVF("flow", 3.3),
				rule.NewEQFVF("flow", 4),
				rule.NewNEFVF("flow", 4),
			},
		},
		{field: "group",
			want: []rule.Rule{
				rule.NewEQFVS("group", "Nelson"),
				rule.NewNEFVS("group", "Nelson"),
				rule.NewEQFVS("group", "Drake"),
				rule.NewNEFVS("group", "Drake"),
			},
		},
	}
	ruleFields := []string{"band", "flow", "group"}
	for _, c := range cases {
		got := generateValueRules(inputDescription, ruleFields, c.field)
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
			rule.NewEQFVS("group", "a"),
			rule.NewGEFVI("band", 4),
			rule.NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
		},
			want: []rule.Rule{
				rule.MustNewAnd(
					rule.NewGEFVI("band", 4),
					rule.NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
				),
				rule.MustNewAnd(
					rule.NewGEFVI("band", 4),
					rule.NewEQFVS("group", "a"),
				),
				rule.MustNewOr(
					rule.NewGEFVI("band", 4),
					rule.NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
				),
				rule.MustNewOr(
					rule.NewGEFVI("band", 4),
					rule.NewEQFVS("group", "a"),
				),
				rule.MustNewAnd(
					rule.NewEQFVS("group", "a"),
					rule.NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
				),
				rule.MustNewOr(
					rule.NewEQFVS("group", "a"),
					rule.NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
				),
			},
		},
		{in: []rule.Rule{
			rule.NewEQFVS("team", "a"),
			rule.NewGEFVI("band", 4),
			rule.NewGEFVI("band", 2),
			rule.NewGEFVI("flow", 6),
		},
			want: []rule.Rule{
				rule.MustNewAnd(rule.NewGEFVI("band", 2), rule.NewGEFVI("flow", 6)),
				rule.MustNewAnd(rule.NewGEFVI("band", 2), rule.NewEQFVS("team", "a")),
				rule.MustNewOr(rule.NewGEFVI("band", 2), rule.NewGEFVI("flow", 6)),
				rule.MustNewOr(rule.NewGEFVI("band", 2), rule.NewEQFVS("team", "a")),
				rule.MustNewAnd(rule.NewGEFVI("band", 4), rule.NewGEFVI("flow", 6)),
				rule.MustNewAnd(rule.NewGEFVI("band", 4), rule.NewEQFVS("team", "a")),
				rule.MustNewOr(rule.NewGEFVI("band", 4), rule.NewGEFVI("flow", 6)),
				rule.MustNewOr(rule.NewGEFVI("band", 4), rule.NewEQFVS("team", "a")),
				rule.MustNewAnd(rule.NewGEFVI("flow", 6), rule.NewEQFVS("team", "a")),
				rule.MustNewOr(rule.NewGEFVI("flow", 6), rule.NewEQFVS("team", "a")),
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
		{in: []rule.Rule{rule.NewEQFVS("team", "a")}, want: []rule.Rule{}},
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

func TestMakeCompareValues(t *testing.T) {
	values1 := map[string]description.Value{
		"a": description.Value{dlit.MustNew("a"), 2},
		"c": description.Value{dlit.MustNew("c"), 2},
		"d": description.Value{dlit.MustNew("d"), 2},
		"e": description.Value{dlit.MustNew("e"), 2},
		"f": description.Value{dlit.MustNew("f"), 2},
	}
	values2 := map[string]description.Value{
		"a": description.Value{dlit.MustNew("a"), 2},
		"c": description.Value{dlit.MustNew("c"), 1},
		"d": description.Value{dlit.MustNew("d"), 2},
		"e": description.Value{dlit.MustNew("e"), 2},
		"f": description.Value{dlit.MustNew("f"), 2},
	}
	cases := []struct {
		values map[string]description.Value
		i      int
		want   []*dlit.Literal
	}{
		{
			values: values1,
			i:      2,
			want:   []*dlit.Literal{dlit.NewString("c")},
		},
		{
			values: values2,
			i:      2,
			want:   []*dlit.Literal{},
		},
		{
			values: values1,
			i:      3,
			want:   []*dlit.Literal{dlit.NewString("a"), dlit.NewString("c")},
		},
		{
			values: values2,
			i:      3,
			want:   []*dlit.Literal{},
		},
		{
			values: values1,
			i:      4,
			want:   []*dlit.Literal{dlit.NewString("d")},
		},
		{
			values: values2,
			i:      4,
			want:   []*dlit.Literal{dlit.NewString("d")},
		},
		{
			values: values1,
			i:      5,
			want:   []*dlit.Literal{dlit.NewString("a"), dlit.NewString("d")},
		},
		{
			values: values2,
			i:      5,
			want:   []*dlit.Literal{dlit.NewString("a"), dlit.NewString("d")},
		},
		{
			values: values1,
			i:      6,
			want:   []*dlit.Literal{dlit.NewString("c"), dlit.NewString("d")},
		},
		{
			values: values2,
			i:      6,
			want:   []*dlit.Literal{},
		},
		{
			values: values1,
			i:      7,
			want: []*dlit.Literal{
				dlit.NewString("a"),
				dlit.NewString("c"),
				dlit.NewString("d"),
			},
		},
		{
			values: values2,
			i:      7,
			want:   []*dlit.Literal{},
		},
		{
			values: values1,
			i:      8,
			want:   []*dlit.Literal{dlit.NewString("e")},
		},
		{
			values: values2,
			i:      8,
			want:   []*dlit.Literal{dlit.NewString("e")},
		},
		{
			values: values1,
			i:      9,
			want:   []*dlit.Literal{dlit.NewString("a"), dlit.NewString("e")},
		},
		{
			values: values2,
			i:      9,
			want:   []*dlit.Literal{dlit.NewString("a"), dlit.NewString("e")},
		},
		{
			values: values1,
			i:      10,
			want:   []*dlit.Literal{dlit.NewString("c"), dlit.NewString("e")},
		},
		{
			values: values2,
			i:      10,
			want:   []*dlit.Literal{},
		},
		{
			values: values1,
			i:      11,
			want: []*dlit.Literal{
				dlit.NewString("a"),
				dlit.NewString("c"),
				dlit.NewString("e"),
			},
		},
		{
			values: values2,
			i:      11,
			want:   []*dlit.Literal{},
		},
		{
			values: values1,
			i:      12,
			want:   []*dlit.Literal{dlit.NewString("d"), dlit.NewString("e")},
		},
		{
			values: values2,
			i:      12,
			want:   []*dlit.Literal{dlit.NewString("d"), dlit.NewString("e")},
		},
		{
			values: values1,
			i:      13,
			want: []*dlit.Literal{
				dlit.NewString("a"),
				dlit.NewString("d"),
				dlit.NewString("e"),
			},
		},
		{
			values: values1,
			i:      14,
			want: []*dlit.Literal{
				dlit.NewString("c"),
				dlit.NewString("d"),
				dlit.NewString("e"),
			},
		},
		{
			values: values1,
			i:      15,
			want: []*dlit.Literal{
				dlit.NewString("a"),
				dlit.NewString("c"),
				dlit.NewString("d"),
				dlit.NewString("e"),
			},
		},
		{
			values: values1,
			i:      16,
			want:   []*dlit.Literal{dlit.NewString("f")},
		},
		{
			values: values2,
			i:      16,
			want:   []*dlit.Literal{dlit.NewString("f")},
		},
		{
			values: values1,
			i:      17,
			want:   []*dlit.Literal{dlit.NewString("a"), dlit.NewString("f")},
		},
		{
			values: values2,
			i:      17,
			want:   []*dlit.Literal{dlit.NewString("a"), dlit.NewString("f")},
		},
	}
	for _, c := range cases {
		got := makeCompareValues(c.values, c.i)
		if len(got) != len(c.want) {
			t.Errorf("makeCompareValues(%v, %d) got: %v, want: %v",
				c.values, c.i, got, c.want)
		}
		for j, v := range got {
			o := c.want[j]
			if o.String() != v.String() {
				t.Errorf("makeCompareValues(%v, %d) got: %v, want: %v",
					c.values, c.i, got, c.want)
			}
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
