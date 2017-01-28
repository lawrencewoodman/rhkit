package rhkit

import (
	"errors"
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/rule"
	"regexp"
	"sort"
	"strings"
	"testing"
)

func TestGenerateRules_1(t *testing.T) {
	testPurpose := "Ensure generates correct rules for each field"
	inputDescription := &Description{
		map[string]*fieldDescription{
			"team": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"a": valueDescription{dlit.NewString("a"), 3},
					"b": valueDescription{dlit.NewString("b"), 3},
					"c": valueDescription{dlit.NewString("c"), 3},
				},
			},
			"teamOut": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"a": valueDescription{dlit.NewString("a"), 3},
					"c": valueDescription{dlit.NewString("c"), 1},
					"d": valueDescription{dlit.NewString("d"), 3},
					"e": valueDescription{dlit.NewString("e"), 3},
					"f": valueDescription{dlit.NewString("f"), 3},
				},
			},
			"level": &fieldDescription{
				kind:  ftInt,
				min:   dlit.MustNew(0),
				max:   dlit.MustNew(5),
				maxDP: 0,
				values: map[string]valueDescription{
					"0": valueDescription{dlit.NewString("0"), 3},
					"1": valueDescription{dlit.NewString("1"), 3},
					"2": valueDescription{dlit.NewString("2"), 1},
					"3": valueDescription{dlit.NewString("3"), 3},
					"4": valueDescription{dlit.NewString("4"), 3},
					"5": valueDescription{dlit.NewString("5"), 3},
				},
			},
			"flow": &fieldDescription{
				kind:  ftFloat,
				min:   dlit.MustNew(0),
				max:   dlit.MustNew(10.5),
				maxDP: 2,
				values: map[string]valueDescription{
					"0.0":  valueDescription{dlit.NewString("0.0"), 3},
					"2.34": valueDescription{dlit.NewString("2.34"), 3},
					"10.5": valueDescription{dlit.NewString("10.5"), 3},
				},
			},
			"position": &fieldDescription{
				kind:  ftInt,
				min:   dlit.MustNew(1),
				max:   dlit.MustNew(13),
				maxDP: 0,
				values: map[string]valueDescription{
					"1":  valueDescription{dlit.NewString("1"), 3},
					"2":  valueDescription{dlit.NewString("2"), 3},
					"3":  valueDescription{dlit.NewString("3"), 3},
					"4":  valueDescription{dlit.NewString("4"), 3},
					"5":  valueDescription{dlit.NewString("5"), 3},
					"6":  valueDescription{dlit.NewString("6"), 3},
					"7":  valueDescription{dlit.NewString("7"), 3},
					"8":  valueDescription{dlit.NewString("8"), 3},
					"9":  valueDescription{dlit.NewString("9"), 3},
					"10": valueDescription{dlit.NewString("10"), 3},
					"11": valueDescription{dlit.NewString("11"), 3},
					"12": valueDescription{dlit.NewString("12"), 3},
					"13": valueDescription{dlit.NewString("13"), 3},
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
	}

	got := GenerateRules(inputDescription, ruleFields)
	if err := rulesContain(got, wantRules); err != nil {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("GenerateRules: %s", err)
	}
}

func TestGenerateRules_2(t *testing.T) {
	testPurpose := "Ensure generates a True rule"
	inputDescription := &Description{
		map[string]*fieldDescription{
			"team": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"a": valueDescription{dlit.MustNew("a"), 3},
					"b": valueDescription{dlit.MustNew("b"), 3},
					"c": valueDescription{dlit.MustNew("c"), 3},
				},
			},
			"teamOut": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"a": valueDescription{dlit.MustNew("a"), 3},
					"c": valueDescription{dlit.MustNew("c"), 3},
					"d": valueDescription{dlit.MustNew("d"), 3},
					"e": valueDescription{dlit.MustNew("e"), 3},
					"f": valueDescription{dlit.MustNew("f"), 3},
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
	inputDescription := &Description{
		map[string]*fieldDescription{
			"directionIn": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"gogledd": valueDescription{dlit.MustNew("gogledd"), 3},
					"de":      valueDescription{dlit.MustNew("de"), 3},
				},
			},
			"directionOut": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"dwyrain":   valueDescription{dlit.MustNew("dwyrain"), 3},
					"gorllewin": valueDescription{dlit.MustNew("gorllewin"), 3},
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
	inputDescription := &Description{
		map[string]*fieldDescription{
			"flow": &fieldDescription{
				kind:   ftInt,
				min:    dlit.MustNew(1),
				max:    dlit.MustNew(20),
				values: map[string]valueDescription{},
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
	inputDescription := &Description{
		map[string]*fieldDescription{
			"flow": &fieldDescription{
				kind:   ftFloat,
				min:    dlit.MustNew(0),
				max:    dlit.MustNew(10.5),
				maxDP:  2,
				values: map[string]valueDescription{},
			},
		},
	}
	ruleFields := []string{"flow"}
	wantRules := []rule.Rule{
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
	}

	got := generateFloatRules(inputDescription, ruleFields, "flow")
	if err := rulesContain(got, wantRules); err != nil {
		t.Errorf("GenerateFloatRules: %s", err)
	}
	if len(got) <= len(wantRules) {
		t.Errorf("GenerateFloatRules: There should be more rules generated: %d",
			len(got))
	}
}

func TestGenerateCompareNumericRules(t *testing.T) {
	inputDescription := &Description{
		map[string]*fieldDescription{
			"band": &fieldDescription{
				kind:   ftInt,
				min:    dlit.MustNew(1),
				max:    dlit.MustNew(3),
				values: map[string]valueDescription{},
			},
			"flowIn": &fieldDescription{
				kind:   ftFloat,
				min:    dlit.MustNew(1),
				max:    dlit.MustNew(4),
				maxDP:  2,
				values: map[string]valueDescription{},
			},
			"flowOut": &fieldDescription{
				kind:   ftFloat,
				min:    dlit.MustNew(0.95),
				max:    dlit.MustNew(4.1),
				maxDP:  2,
				values: map[string]valueDescription{},
			},
			"rateIn": &fieldDescription{
				kind:   ftFloat,
				min:    dlit.MustNew(4.2),
				max:    dlit.MustNew(8.9),
				maxDP:  2,
				values: map[string]valueDescription{},
			},
			"rateOut": &fieldDescription{
				kind:   ftFloat,
				min:    dlit.MustNew(0.1),
				max:    dlit.MustNew(0.9),
				maxDP:  2,
				values: map[string]valueDescription{},
			},
			"group": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"Nelson":      valueDescription{dlit.NewString("Nelson"), 3},
					"Collingwood": valueDescription{dlit.NewString("Collingwood"), 1},
					"Mountbatten": valueDescription{dlit.NewString("Mountbatten"), 1},
					"Drake":       valueDescription{dlit.NewString("Drake"), 2},
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
	inputDescription := &Description{
		map[string]*fieldDescription{
			"band": &fieldDescription{
				kind:   ftInt,
				min:    dlit.MustNew(1),
				max:    dlit.MustNew(3),
				values: map[string]valueDescription{},
			},
			"groupA": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"Nelson":      valueDescription{dlit.NewString("Nelson"), 3},
					"Collingwood": valueDescription{dlit.NewString("Collingwood"), 1},
					"Mountbatten": valueDescription{dlit.NewString("Mountbatten"), 1},
					"Drake":       valueDescription{dlit.NewString("Drake"), 2},
				},
			},
			"groupB": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"Nelson":      valueDescription{dlit.NewString("Nelson"), 3},
					"Mountbatten": valueDescription{dlit.NewString("Mountbatten"), 1},
					"Drake":       valueDescription{dlit.NewString("Drake"), 2},
				},
			},
			"groupC": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"Nelson": valueDescription{dlit.NewString("Nelson"), 3},
					"Drake":  valueDescription{dlit.NewString("Drake"), 2},
				},
			},
			"groupD": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"Drake": valueDescription{dlit.NewString("Drake"), 2},
				},
			},
			"groupE": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"Drake":       valueDescription{dlit.NewString("Drake"), 2},
					"Chaucer":     valueDescription{dlit.NewString("Chaucer"), 2},
					"Shakespeare": valueDescription{dlit.NewString("Shakespeare"), 2},
					"Marlowe":     valueDescription{dlit.NewString("Marlowe"), 2},
				},
			},
			"groupF": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"Nelson":      valueDescription{dlit.NewString("Nelson"), 3},
					"Drake":       valueDescription{dlit.NewString("Drake"), 2},
					"Chaucer":     valueDescription{dlit.NewString("Chaucer"), 2},
					"Shakespeare": valueDescription{dlit.NewString("Shakespeare"), 2},
					"Marlowe":     valueDescription{dlit.NewString("Marlowe"), 2},
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
	inputDescription := &Description{
		map[string]*fieldDescription{
			"band": &fieldDescription{
				kind:   ftInt,
				min:    dlit.MustNew(1),
				max:    dlit.MustNew(3),
				values: map[string]valueDescription{},
			},
			"flow": &fieldDescription{
				kind:   ftFloat,
				min:    dlit.MustNew(1),
				max:    dlit.MustNew(3),
				maxDP:  2,
				values: map[string]valueDescription{},
			},
			"groupA": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"Fred":    valueDescription{dlit.NewString("Fred"), 3},
					"Mary":    valueDescription{dlit.NewString("Mary"), 4},
					"Rebecca": valueDescription{dlit.NewString("Rebecca"), 2},
				},
			},

			"groupB": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"Fred":    valueDescription{dlit.NewString("Fred"), 3},
					"Mary":    valueDescription{dlit.NewString("Mary"), 4},
					"Rebecca": valueDescription{dlit.NewString("Rebecca"), 2},
					"Harry":   valueDescription{dlit.NewString("Harry"), 2},
					"Dinah":   valueDescription{dlit.NewString("Dinah"), 2},
					"Israel":  valueDescription{dlit.NewString("Israel"), 2},
					"Sarah":   valueDescription{dlit.NewString("Sarah"), 2},
					"Ishmael": valueDescription{dlit.NewString("Ishmael"), 2},
					"Caen":    valueDescription{dlit.NewString("Caen"), 2},
					"Abel":    valueDescription{dlit.NewString("Abel"), 2},
					"Noah":    valueDescription{dlit.NewString("Noah"), 2},
					"Isaac":   valueDescription{dlit.NewString("Isaac"), 2},
					"Moses":   valueDescription{dlit.NewString("Moses"), 2},
				},
			},
			"groupC": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"Fred":    valueDescription{dlit.NewString("Fred"), 3},
					"Mary":    valueDescription{dlit.NewString("Mary"), 4},
					"Rebecca": valueDescription{dlit.NewString("Rebecca"), 2},
					"Harry":   valueDescription{dlit.NewString("Harry"), 2},
				},
			},
			"groupD": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"Fred":    valueDescription{dlit.NewString("Fred"), 3},
					"Mary":    valueDescription{dlit.NewString("Mary"), 4},
					"Rebecca": valueDescription{dlit.NewString("Rebecca"), 1},
					"Harry":   valueDescription{dlit.NewString("Harry"), 2},
				},
			},
			"groupE": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"Fred":    valueDescription{dlit.NewString("Fred"), 3},
					"Mary":    valueDescription{dlit.NewString("Mary"), 4},
					"Rebecca": valueDescription{dlit.NewString("Rebecca"), 2},
					"Harry":   valueDescription{dlit.NewString("Harry"), 2},
					"Juliet":  valueDescription{dlit.NewString("Juliet"), 2},
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
	inputDescription := &Description{
		map[string]*fieldDescription{
			"group": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"Fred":    valueDescription{dlit.NewString("Fred"), 3},
					"Mary":    valueDescription{dlit.NewString("Mary"), 4},
					"Rebecca": valueDescription{dlit.NewString("Rebecca"), 2},
					"Harry":   valueDescription{dlit.NewString("Harry"), 2},
					"Dinah":   valueDescription{dlit.NewString("Dinah"), 2},
					"Israel":  valueDescription{dlit.NewString("Israel"), 2},
					"Sarah":   valueDescription{dlit.NewString("Sarah"), 2},
					"Ishmael": valueDescription{dlit.NewString("Ishmael"), 2},
					"Caen":    valueDescription{dlit.NewString("Caen"), 2},
					"Abel":    valueDescription{dlit.NewString("Abel"), 2},
					"Noah":    valueDescription{dlit.NewString("Noah"), 2},
					"Isaac":   valueDescription{dlit.NewString("Isaac"), 2},
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
	inputDescription := &Description{
		map[string]*fieldDescription{
			"band": &fieldDescription{
				kind: ftInt,
				min:  dlit.MustNew(1),
				max:  dlit.MustNew(4),
				values: map[string]valueDescription{
					"1": valueDescription{dlit.NewString("1"), 3},
					"2": valueDescription{dlit.NewString("2"), 1},
					"3": valueDescription{dlit.NewString("3"), 2},
					"4": valueDescription{dlit.NewString("4"), 5},
				},
			},
			"flow": &fieldDescription{
				kind:  ftFloat,
				min:   dlit.MustNew(1),
				max:   dlit.MustNew(4),
				maxDP: 2,
				values: map[string]valueDescription{
					"1":    valueDescription{dlit.NewString("1"), 3},
					"2":    valueDescription{dlit.NewString("2"), 1},
					"2.90": valueDescription{dlit.NewString("2.90"), 1},
					"3.37": valueDescription{dlit.NewString("3.37"), 2},
					"3.3":  valueDescription{dlit.NewString("3.3"), 2},
					"4":    valueDescription{dlit.NewString("4"), 5},
				},
			},
			"group": &fieldDescription{
				kind: ftString,
				values: map[string]valueDescription{
					"Nelson":      valueDescription{dlit.NewString("Nelson"), 3},
					"Collingwood": valueDescription{dlit.NewString("Collingwood"), 1},
					"Mountbatten": valueDescription{dlit.NewString("Mountbatten"), 1},
					"Drake":       valueDescription{dlit.NewString("Drake"), 2},
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
	values1 := map[string]valueDescription{
		"a": valueDescription{dlit.MustNew("a"), 2},
		"c": valueDescription{dlit.MustNew("c"), 2},
		"d": valueDescription{dlit.MustNew("d"), 2},
		"e": valueDescription{dlit.MustNew("e"), 2},
		"f": valueDescription{dlit.MustNew("f"), 2},
	}
	values2 := map[string]valueDescription{
		"a": valueDescription{dlit.MustNew("a"), 2},
		"c": valueDescription{dlit.MustNew("c"), 1},
		"d": valueDescription{dlit.MustNew("d"), 2},
		"e": valueDescription{dlit.MustNew("e"), 2},
		"f": valueDescription{dlit.MustNew("f"), 2},
	}
	cases := []struct {
		values map[string]valueDescription
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
			t.Errorf("makeCompareValues(%s, %d) got: %v, want: %v",
				c.values, c.i, got, c.want)
		}
		for j, v := range got {
			o := c.want[j]
			if o.String() != v.String() {
				t.Errorf("makeCompareValues(%s, %d) got: %v, want: %v",
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
