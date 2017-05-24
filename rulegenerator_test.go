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

func TestGenerateNumRules(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"flow": &description.Field{
				Kind:   fieldtype.Number,
				Min:    dlit.MustNew(0),
				Max:    dlit.MustNew(10.5),
				MaxDP:  2,
				Values: map[string]description.Value{},
			},
			"flowConstant": &description.Field{
				Kind:   fieldtype.Number,
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
				rule.NewLEFV("flow", dlit.MustNew(1.05)),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(1.05)),
					rule.NewGEFV("flow", dlit.MustNew(2.1))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(1.05)),
					rule.NewGEFV("flow", dlit.MustNew(3.68))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(1.05)),
					rule.NewGEFV("flow", dlit.MustNew(4.2))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(1.05)),
					rule.NewGEFV("flow", dlit.MustNew(5.25))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(1.05)),
					rule.NewGEFV("flow", dlit.MustNew(6.83))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(1.05)),
					rule.NewGEFV("flow", dlit.MustNew(7.88))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(1.05)),
					rule.NewGEFV("flow", dlit.MustNew(8.93))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(1.05)),
					rule.NewGEFV("flow", dlit.MustNew(9.45))),
				rule.NewLEFV("flow", dlit.MustNew(2.1)),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(2.1)),
					rule.NewGEFV("flow", dlit.MustNew(1.05))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(2.1)),
					rule.NewGEFV("flow", dlit.MustNew(3.68))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(2.1)),
					rule.NewGEFV("flow", dlit.MustNew(4.2))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(2.1)),
					rule.NewGEFV("flow", dlit.MustNew(5.25))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(2.1)),
					rule.NewGEFV("flow", dlit.MustNew(6.83))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(2.1)),
					rule.NewGEFV("flow", dlit.MustNew(7.88))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(2.1)),
					rule.NewGEFV("flow", dlit.MustNew(8.93))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(2.1)),
					rule.NewGEFV("flow", dlit.MustNew(9.45))),
				rule.NewLEFV("flow", dlit.MustNew(3.68)),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(3.68)),
					rule.NewGEFV("flow", dlit.MustNew(1.05))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(3.68)),
					rule.NewGEFV("flow", dlit.MustNew(2.1))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(3.68)),
					rule.NewGEFV("flow", dlit.MustNew(4.2))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(3.68)),
					rule.NewGEFV("flow", dlit.MustNew(5.25))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(3.68)),
					rule.NewGEFV("flow", dlit.MustNew(6.83))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(3.68)),
					rule.NewGEFV("flow", dlit.MustNew(7.88))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(3.68)),
					rule.NewGEFV("flow", dlit.MustNew(8.93))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(3.68)),
					rule.NewGEFV("flow", dlit.MustNew(9.45))),
				rule.NewLEFV("flow", dlit.MustNew(4.2)),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(4.2)),
					rule.NewGEFV("flow", dlit.MustNew(1.05))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(4.2)),
					rule.NewGEFV("flow", dlit.MustNew(2.1))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(4.2)),
					rule.NewGEFV("flow", dlit.MustNew(3.68))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(4.2)),
					rule.NewGEFV("flow", dlit.MustNew(5.25))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(4.2)),
					rule.NewGEFV("flow", dlit.MustNew(6.83))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(4.2)),
					rule.NewGEFV("flow", dlit.MustNew(7.88))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(4.2)),
					rule.NewGEFV("flow", dlit.MustNew(8.93))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(4.2)),
					rule.NewGEFV("flow", dlit.MustNew(9.45))),
				rule.NewLEFV("flow", dlit.MustNew(5.25)),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(5.25)),
					rule.NewGEFV("flow", dlit.MustNew(1.05))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(5.25)),
					rule.NewGEFV("flow", dlit.MustNew(2.1))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(5.25)),
					rule.NewGEFV("flow", dlit.MustNew(3.68))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(5.25)),
					rule.NewGEFV("flow", dlit.MustNew(4.2))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(5.25)),
					rule.NewGEFV("flow", dlit.MustNew(6.83))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(5.25)),
					rule.NewGEFV("flow", dlit.MustNew(7.88))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(5.25)),
					rule.NewGEFV("flow", dlit.MustNew(8.93))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(5.25)),
					rule.NewGEFV("flow", dlit.MustNew(9.45))),
				rule.NewLEFV("flow", dlit.MustNew(6.83)),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(6.83)),
					rule.NewGEFV("flow", dlit.MustNew(1.05))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(6.83)),
					rule.NewGEFV("flow", dlit.MustNew(2.1))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(6.83)),
					rule.NewGEFV("flow", dlit.MustNew(3.68))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(6.83)),
					rule.NewGEFV("flow", dlit.MustNew(4.2))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(6.83)),
					rule.NewGEFV("flow", dlit.MustNew(5.25))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(6.83)),
					rule.NewGEFV("flow", dlit.MustNew(7.88))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(6.83)),
					rule.NewGEFV("flow", dlit.MustNew(8.93))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(6.83)),
					rule.NewGEFV("flow", dlit.MustNew(9.45))),
				rule.NewLEFV("flow", dlit.MustNew(7.88)),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(7.88)),
					rule.NewGEFV("flow", dlit.MustNew(1.05))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(7.88)),
					rule.NewGEFV("flow", dlit.MustNew(2.1))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(7.88)),
					rule.NewGEFV("flow", dlit.MustNew(3.68))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(7.88)),
					rule.NewGEFV("flow", dlit.MustNew(4.2))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(7.88)),
					rule.NewGEFV("flow", dlit.MustNew(5.25))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(7.88)),
					rule.NewGEFV("flow", dlit.MustNew(6.83))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(7.88)),
					rule.NewGEFV("flow", dlit.MustNew(8.93))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(7.88)),
					rule.NewGEFV("flow", dlit.MustNew(9.45))),
				rule.NewLEFV("flow", dlit.MustNew(8.93)),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(8.93)),
					rule.NewGEFV("flow", dlit.MustNew(1.05))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(8.93)),
					rule.NewGEFV("flow", dlit.MustNew(2.1))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(8.93)),
					rule.NewGEFV("flow", dlit.MustNew(3.68))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(8.93)),
					rule.NewGEFV("flow", dlit.MustNew(4.2))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(8.93)),
					rule.NewGEFV("flow", dlit.MustNew(5.25))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(8.93)),
					rule.NewGEFV("flow", dlit.MustNew(6.83))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(8.93)),
					rule.NewGEFV("flow", dlit.MustNew(7.88))),
				rule.MustNewOr(rule.NewLEFV("flow", dlit.MustNew(8.93)),
					rule.NewGEFV("flow", dlit.MustNew(9.45))),
				rule.NewLEFV("flow", dlit.MustNew(9.45)),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(9.45)),
					rule.NewGEFV("flow", dlit.MustNew(1.05))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(9.45)),
					rule.NewGEFV("flow", dlit.MustNew(2.1))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(9.45)),
					rule.NewGEFV("flow", dlit.MustNew(3.68))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(9.45)),
					rule.NewGEFV("flow", dlit.MustNew(4.2))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(9.45)),
					rule.NewGEFV("flow", dlit.MustNew(5.25))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(9.45)),
					rule.NewGEFV("flow", dlit.MustNew(6.83))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(9.45)),
					rule.NewGEFV("flow", dlit.MustNew(7.88))),
				rule.MustNewAnd(rule.NewLEFV("flow", dlit.MustNew(9.45)),
					rule.NewGEFV("flow", dlit.MustNew(8.93))),
				rule.NewGEFV("flow", dlit.MustNew(1.05)),
				rule.NewGEFV("flow", dlit.MustNew(2.1)),
				rule.NewGEFV("flow", dlit.MustNew(3.68)),
				rule.NewGEFV("flow", dlit.MustNew(4.2)),
				rule.NewGEFV("flow", dlit.MustNew(5.25)),
				rule.NewGEFV("flow", dlit.MustNew(6.83)),
				rule.NewGEFV("flow", dlit.MustNew(7.88)),
				rule.NewGEFV("flow", dlit.MustNew(8.93)),
				rule.NewGEFV("flow", dlit.MustNew(9.45)),
			},
		},
		{field: "flowConstant", want: []rule.Rule{}},
	}

	for i, c := range cases {
		for complexity := 1; complexity < 10; complexity++ {
			got :=
				generateNumRules(inputDescription, ruleFields, complexity, c.field)
			if err := rulesContain(got, c.want); err != nil {
				t.Errorf("(%d) GenerateNumRules: %s - got: %v", i, err, got)
			}
			if len(got) < len(c.want) {
				t.Errorf(
					"GenerateFloatRules: There should be more rules generated for field: %s, got: %d",
					c.field, len(got),
				)
			}
		}
	}
}

func TestGenerateAddRules(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"flowIn": &description.Field{
				Kind:   fieldtype.Number,
				Min:    dlit.MustNew(0),
				Max:    dlit.MustNew(10.5),
				MaxDP:  2,
				Values: map[string]description.Value{},
			},
			"flowOut": &description.Field{
				Kind:   fieldtype.Number,
				Min:    dlit.MustNew(1.2),
				Max:    dlit.MustNew(5.2),
				MaxDP:  2,
				Values: map[string]description.Value{},
			},
			"flowSideA": &description.Field{
				Kind:   fieldtype.Number,
				Min:    dlit.MustNew(2.2),
				Max:    dlit.MustNew(2.2),
				MaxDP:  1,
				Values: map[string]description.Value{},
			},
			"flowSideB": &description.Field{
				Kind:   fieldtype.Number,
				Min:    dlit.MustNew(2.2),
				Max:    dlit.MustNew(2.2),
				MaxDP:  1,
				Values: map[string]description.Value{},
			},
		},
	}
	cases := []struct {
		field        string
		complexities []int
		want         []rule.Rule
	}{
		{field: "flowIn",
			complexities: []int{1, 2},
			want:         []rule.Rule{},
		},
		{field: "flowIn",
			complexities: []int{3, 4, 5, 6, 7, 8, 9, 10},
			want: []rule.Rule{
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(1.92)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(10.62)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(11.35)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(12.07)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(12.8)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(13.52)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(14.25)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(14.97)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(2.65)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(3.38)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(4.1)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(4.82)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(5.55)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(6.27)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(7)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(7.72)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(8.45)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(9.17)),
				rule.NewAddLEF("flowIn", "flowOut", dlit.MustNew(9.9)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(1.92)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(10.62)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(11.35)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(12.07)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(12.8)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(13.52)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(14.25)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(14.97)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(2.65)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(3.38)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(4.1)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(4.82)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(5.55)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(6.27)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(7)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(7.72)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(8.45)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(9.17)),
				rule.NewAddGEF("flowIn", "flowOut", dlit.MustNew(9.9)),

				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(10.08)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(10.6)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(11.13)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(11.65)),
				rule.NewAddLEF("flowIn", "flowSideA", dlit.MustNew(12.18)),
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
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(10.08)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(10.6)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(11.13)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(11.65)),
				rule.NewAddGEF("flowIn", "flowSideA", dlit.MustNew(12.18)),
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
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(10.08)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(10.6)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(11.13)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(11.65)),
				rule.NewAddLEF("flowIn", "flowSideB", dlit.MustNew(12.18)),
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
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(10.08)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(10.6)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(11.13)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(11.65)),
				rule.NewAddGEF("flowIn", "flowSideB", dlit.MustNew(12.18)),
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
			},
		},
		{field: "flowOut",
			complexities: []int{3, 4, 5, 6, 7, 8, 9, 10},
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
				rule.NewAddGEF("flowOut", "flowSideB", dlit.MustNew(7.)),
			},
		},
		{field: "flowSideA",
			complexities: []int{10},
			want:         []rule.Rule{},
		},
		{field: "flowSideB",
			complexities: []int{10},
			want:         []rule.Rule{},
		},
	}
	ruleFields := []string{"flowIn", "flowOut", "flowSideA", "flowSideB"}

	for i, c := range cases {
		for _, complexity := range c.complexities {
			got := generateAddRules(inputDescription, ruleFields, complexity, c.field)
			if err := matchRulesUnordered(got, c.want); err != nil {
				gotRuleStrs := rulesToSortedStrings(got)
				wantRuleStrs := rulesToSortedStrings(c.want)
				t.Errorf(
					"(%d) matchRulesUnordered() rules don't match: %s\ngot: %s\nwant: %s\n",
					i, err, gotRuleStrs, wantRuleStrs)
			}
		}
	}
}

func TestGenerateMulRules(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"flowIn": &description.Field{
				Kind:   fieldtype.Number,
				Min:    dlit.MustNew(0),
				Max:    dlit.MustNew(10.5),
				MaxDP:  2,
				Values: map[string]description.Value{},
			},
			"flowOut": &description.Field{
				Kind:   fieldtype.Number,
				Min:    dlit.MustNew(1.2),
				Max:    dlit.MustNew(5.2),
				MaxDP:  2,
				Values: map[string]description.Value{},
			},
			"flowSideA": &description.Field{
				Kind:   fieldtype.Number,
				Min:    dlit.MustNew(2.2),
				Max:    dlit.MustNew(2.2),
				MaxDP:  1,
				Values: map[string]description.Value{},
			},
			"flowSideB": &description.Field{
				Kind:   fieldtype.Number,
				Min:    dlit.MustNew(2.2),
				Max:    dlit.MustNew(2.2),
				MaxDP:  1,
				Values: map[string]description.Value{},
			},
		},
	}
	cases := []struct {
		field        string
		complexities []int
		want         []rule.Rule
	}{
		{field: "flowIn",
			complexities: []int{1, 2},
			want:         []rule.Rule{},
		},
		{field: "flowIn",
			complexities: []int{3, 4, 5, 6, 7, 8, 9, 10},
			want: []rule.Rule{
				rule.NewMulLEF("flowIn", "flowOut", dlit.MustNew(10.92)),
				rule.NewMulLEF("flowIn", "flowOut", dlit.MustNew(13.65)),
				rule.NewMulLEF("flowIn", "flowOut", dlit.MustNew(16.38)),
				rule.NewMulLEF("flowIn", "flowOut", dlit.MustNew(19.11)),
				rule.NewMulLEF("flowIn", "flowOut", dlit.MustNew(2.73)),
				rule.NewMulLEF("flowIn", "flowOut", dlit.MustNew(21.84)),
				rule.NewMulLEF("flowIn", "flowOut", dlit.MustNew(24.57)),
				rule.NewMulLEF("flowIn", "flowOut", dlit.MustNew(27.3)),
				rule.NewMulLEF("flowIn", "flowOut", dlit.MustNew(30.03)),
				rule.NewMulLEF("flowIn", "flowOut", dlit.MustNew(32.76)),
				rule.NewMulLEF("flowIn", "flowOut", dlit.MustNew(35.49)),
				rule.NewMulLEF("flowIn", "flowOut", dlit.MustNew(38.22)),
				rule.NewMulLEF("flowIn", "flowOut", dlit.MustNew(40.95)),
				rule.NewMulLEF("flowIn", "flowOut", dlit.MustNew(43.68)),
				rule.NewMulLEF("flowIn", "flowOut", dlit.MustNew(46.41)),
				rule.NewMulLEF("flowIn", "flowOut", dlit.MustNew(49.14)),
				rule.NewMulLEF("flowIn", "flowOut", dlit.MustNew(5.46)),
				rule.NewMulLEF("flowIn", "flowOut", dlit.MustNew(51.87)),
				rule.NewMulLEF("flowIn", "flowOut", dlit.MustNew(8.19)),

				rule.NewMulGEF("flowIn", "flowOut", dlit.MustNew(10.92)),
				rule.NewMulGEF("flowIn", "flowOut", dlit.MustNew(13.65)),
				rule.NewMulGEF("flowIn", "flowOut", dlit.MustNew(16.38)),
				rule.NewMulGEF("flowIn", "flowOut", dlit.MustNew(19.11)),
				rule.NewMulGEF("flowIn", "flowOut", dlit.MustNew(2.73)),
				rule.NewMulGEF("flowIn", "flowOut", dlit.MustNew(21.84)),
				rule.NewMulGEF("flowIn", "flowOut", dlit.MustNew(24.57)),
				rule.NewMulGEF("flowIn", "flowOut", dlit.MustNew(27.3)),
				rule.NewMulGEF("flowIn", "flowOut", dlit.MustNew(30.03)),
				rule.NewMulGEF("flowIn", "flowOut", dlit.MustNew(32.76)),
				rule.NewMulGEF("flowIn", "flowOut", dlit.MustNew(35.49)),
				rule.NewMulGEF("flowIn", "flowOut", dlit.MustNew(38.22)),
				rule.NewMulGEF("flowIn", "flowOut", dlit.MustNew(40.95)),
				rule.NewMulGEF("flowIn", "flowOut", dlit.MustNew(43.68)),
				rule.NewMulGEF("flowIn", "flowOut", dlit.MustNew(46.41)),
				rule.NewMulGEF("flowIn", "flowOut", dlit.MustNew(49.14)),
				rule.NewMulGEF("flowIn", "flowOut", dlit.MustNew(5.46)),
				rule.NewMulGEF("flowIn", "flowOut", dlit.MustNew(51.87)),
				rule.NewMulGEF("flowIn", "flowOut", dlit.MustNew(8.19)),

				rule.NewMulLEF("flowIn", "flowSideA", dlit.MustNew(1.16)),
				rule.NewMulLEF("flowIn", "flowSideA", dlit.MustNew(10.4)),
				rule.NewMulLEF("flowIn", "flowSideA", dlit.MustNew(11.55)),
				rule.NewMulLEF("flowIn", "flowSideA", dlit.MustNew(12.7)),
				rule.NewMulLEF("flowIn", "flowSideA", dlit.MustNew(13.86)),
				rule.NewMulLEF("flowIn", "flowSideA", dlit.MustNew(15.01)),
				rule.NewMulLEF("flowIn", "flowSideA", dlit.MustNew(16.17)),
				rule.NewMulLEF("flowIn", "flowSideA", dlit.MustNew(17.33)),
				rule.NewMulLEF("flowIn", "flowSideA", dlit.MustNew(18.48)),
				rule.NewMulLEF("flowIn", "flowSideA", dlit.MustNew(19.64)),
				rule.NewMulLEF("flowIn", "flowSideA", dlit.MustNew(2.31)),
				rule.NewMulLEF("flowIn", "flowSideA", dlit.MustNew(20.79)),
				rule.NewMulLEF("flowIn", "flowSideA", dlit.MustNew(21.95)),
				rule.NewMulLEF("flowIn", "flowSideA", dlit.MustNew(3.47)),
				rule.NewMulLEF("flowIn", "flowSideA", dlit.MustNew(4.62)),
				rule.NewMulLEF("flowIn", "flowSideA", dlit.MustNew(5.78)),
				rule.NewMulLEF("flowIn", "flowSideA", dlit.MustNew(6.93)),
				rule.NewMulLEF("flowIn", "flowSideA", dlit.MustNew(8.09)),
				rule.NewMulLEF("flowIn", "flowSideA", dlit.MustNew(9.24)),

				rule.NewMulGEF("flowIn", "flowSideA", dlit.MustNew(1.16)),
				rule.NewMulGEF("flowIn", "flowSideA", dlit.MustNew(10.4)),
				rule.NewMulGEF("flowIn", "flowSideA", dlit.MustNew(11.55)),
				rule.NewMulGEF("flowIn", "flowSideA", dlit.MustNew(12.7)),
				rule.NewMulGEF("flowIn", "flowSideA", dlit.MustNew(13.86)),
				rule.NewMulGEF("flowIn", "flowSideA", dlit.MustNew(15.01)),
				rule.NewMulGEF("flowIn", "flowSideA", dlit.MustNew(16.17)),
				rule.NewMulGEF("flowIn", "flowSideA", dlit.MustNew(17.33)),
				rule.NewMulGEF("flowIn", "flowSideA", dlit.MustNew(18.48)),
				rule.NewMulGEF("flowIn", "flowSideA", dlit.MustNew(19.64)),
				rule.NewMulGEF("flowIn", "flowSideA", dlit.MustNew(2.31)),
				rule.NewMulGEF("flowIn", "flowSideA", dlit.MustNew(20.79)),
				rule.NewMulGEF("flowIn", "flowSideA", dlit.MustNew(21.95)),
				rule.NewMulGEF("flowIn", "flowSideA", dlit.MustNew(3.47)),
				rule.NewMulGEF("flowIn", "flowSideA", dlit.MustNew(4.62)),
				rule.NewMulGEF("flowIn", "flowSideA", dlit.MustNew(5.78)),
				rule.NewMulGEF("flowIn", "flowSideA", dlit.MustNew(6.93)),
				rule.NewMulGEF("flowIn", "flowSideA", dlit.MustNew(8.09)),
				rule.NewMulGEF("flowIn", "flowSideA", dlit.MustNew(9.24)),

				rule.NewMulLEF("flowIn", "flowSideB", dlit.MustNew(1.16)),
				rule.NewMulLEF("flowIn", "flowSideB", dlit.MustNew(10.4)),
				rule.NewMulLEF("flowIn", "flowSideB", dlit.MustNew(11.55)),
				rule.NewMulLEF("flowIn", "flowSideB", dlit.MustNew(12.7)),
				rule.NewMulLEF("flowIn", "flowSideB", dlit.MustNew(13.86)),
				rule.NewMulLEF("flowIn", "flowSideB", dlit.MustNew(15.01)),
				rule.NewMulLEF("flowIn", "flowSideB", dlit.MustNew(16.17)),
				rule.NewMulLEF("flowIn", "flowSideB", dlit.MustNew(17.33)),
				rule.NewMulLEF("flowIn", "flowSideB", dlit.MustNew(18.48)),
				rule.NewMulLEF("flowIn", "flowSideB", dlit.MustNew(19.64)),
				rule.NewMulLEF("flowIn", "flowSideB", dlit.MustNew(2.31)),
				rule.NewMulLEF("flowIn", "flowSideB", dlit.MustNew(20.79)),
				rule.NewMulLEF("flowIn", "flowSideB", dlit.MustNew(21.95)),
				rule.NewMulLEF("flowIn", "flowSideB", dlit.MustNew(3.47)),
				rule.NewMulLEF("flowIn", "flowSideB", dlit.MustNew(4.62)),
				rule.NewMulLEF("flowIn", "flowSideB", dlit.MustNew(5.78)),
				rule.NewMulLEF("flowIn", "flowSideB", dlit.MustNew(6.93)),
				rule.NewMulLEF("flowIn", "flowSideB", dlit.MustNew(8.09)),
				rule.NewMulLEF("flowIn", "flowSideB", dlit.MustNew(9.24)),

				rule.NewMulGEF("flowIn", "flowSideB", dlit.MustNew(1.16)),
				rule.NewMulGEF("flowIn", "flowSideB", dlit.MustNew(10.4)),
				rule.NewMulGEF("flowIn", "flowSideB", dlit.MustNew(11.55)),
				rule.NewMulGEF("flowIn", "flowSideB", dlit.MustNew(12.7)),
				rule.NewMulGEF("flowIn", "flowSideB", dlit.MustNew(13.86)),
				rule.NewMulGEF("flowIn", "flowSideB", dlit.MustNew(15.01)),
				rule.NewMulGEF("flowIn", "flowSideB", dlit.MustNew(16.17)),
				rule.NewMulGEF("flowIn", "flowSideB", dlit.MustNew(17.33)),
				rule.NewMulGEF("flowIn", "flowSideB", dlit.MustNew(18.48)),
				rule.NewMulGEF("flowIn", "flowSideB", dlit.MustNew(19.64)),
				rule.NewMulGEF("flowIn", "flowSideB", dlit.MustNew(2.31)),
				rule.NewMulGEF("flowIn", "flowSideB", dlit.MustNew(20.79)),
				rule.NewMulGEF("flowIn", "flowSideB", dlit.MustNew(21.95)),
				rule.NewMulGEF("flowIn", "flowSideB", dlit.MustNew(3.47)),
				rule.NewMulGEF("flowIn", "flowSideB", dlit.MustNew(4.62)),
				rule.NewMulGEF("flowIn", "flowSideB", dlit.MustNew(5.78)),
				rule.NewMulGEF("flowIn", "flowSideB", dlit.MustNew(6.93)),
				rule.NewMulGEF("flowIn", "flowSideB", dlit.MustNew(8.09)),
				rule.NewMulGEF("flowIn", "flowSideB", dlit.MustNew(9.24)),
			},
		},
		{field: "flowOut",
			complexities: []int{3, 4, 5, 6, 7, 8, 9, 10},
			want: []rule.Rule{
				rule.NewMulLEF("flowOut", "flowSideA", dlit.MustNew(10.12)),
				rule.NewMulLEF("flowOut", "flowSideA", dlit.MustNew(10.56)),
				rule.NewMulLEF("flowOut", "flowSideA", dlit.MustNew(11)),
				rule.NewMulLEF("flowOut", "flowSideA", dlit.MustNew(11.44)),
				rule.NewMulLEF("flowOut", "flowSideA", dlit.MustNew(3.08)),
				rule.NewMulLEF("flowOut", "flowSideA", dlit.MustNew(3.52)),
				rule.NewMulLEF("flowOut", "flowSideA", dlit.MustNew(3.96)),
				rule.NewMulLEF("flowOut", "flowSideA", dlit.MustNew(4.4)),
				rule.NewMulLEF("flowOut", "flowSideA", dlit.MustNew(4.84)),
				rule.NewMulLEF("flowOut", "flowSideA", dlit.MustNew(5.28)),
				rule.NewMulLEF("flowOut", "flowSideA", dlit.MustNew(5.72)),
				rule.NewMulLEF("flowOut", "flowSideA", dlit.MustNew(6.16)),
				rule.NewMulLEF("flowOut", "flowSideA", dlit.MustNew(6.6)),
				rule.NewMulLEF("flowOut", "flowSideA", dlit.MustNew(7.04)),
				rule.NewMulLEF("flowOut", "flowSideA", dlit.MustNew(7.48)),
				rule.NewMulLEF("flowOut", "flowSideA", dlit.MustNew(7.92)),
				rule.NewMulLEF("flowOut", "flowSideA", dlit.MustNew(8.36)),
				rule.NewMulLEF("flowOut", "flowSideA", dlit.MustNew(8.8)),
				rule.NewMulLEF("flowOut", "flowSideA", dlit.MustNew(9.24)),
				rule.NewMulLEF("flowOut", "flowSideA", dlit.MustNew(9.68)),

				rule.NewMulGEF("flowOut", "flowSideA", dlit.MustNew(10.12)),
				rule.NewMulGEF("flowOut", "flowSideA", dlit.MustNew(10.56)),
				rule.NewMulGEF("flowOut", "flowSideA", dlit.MustNew(11)),
				rule.NewMulGEF("flowOut", "flowSideA", dlit.MustNew(11.44)),
				rule.NewMulGEF("flowOut", "flowSideA", dlit.MustNew(3.08)),
				rule.NewMulGEF("flowOut", "flowSideA", dlit.MustNew(3.52)),
				rule.NewMulGEF("flowOut", "flowSideA", dlit.MustNew(3.96)),
				rule.NewMulGEF("flowOut", "flowSideA", dlit.MustNew(4.4)),
				rule.NewMulGEF("flowOut", "flowSideA", dlit.MustNew(4.84)),
				rule.NewMulGEF("flowOut", "flowSideA", dlit.MustNew(5.28)),
				rule.NewMulGEF("flowOut", "flowSideA", dlit.MustNew(5.72)),
				rule.NewMulGEF("flowOut", "flowSideA", dlit.MustNew(6.16)),
				rule.NewMulGEF("flowOut", "flowSideA", dlit.MustNew(6.6)),
				rule.NewMulGEF("flowOut", "flowSideA", dlit.MustNew(7.04)),
				rule.NewMulGEF("flowOut", "flowSideA", dlit.MustNew(7.48)),
				rule.NewMulGEF("flowOut", "flowSideA", dlit.MustNew(7.92)),
				rule.NewMulGEF("flowOut", "flowSideA", dlit.MustNew(8.36)),
				rule.NewMulGEF("flowOut", "flowSideA", dlit.MustNew(8.8)),
				rule.NewMulGEF("flowOut", "flowSideA", dlit.MustNew(9.24)),
				rule.NewMulGEF("flowOut", "flowSideA", dlit.MustNew(9.68)),

				rule.NewMulLEF("flowOut", "flowSideB", dlit.MustNew(10.12)),
				rule.NewMulLEF("flowOut", "flowSideB", dlit.MustNew(10.56)),
				rule.NewMulLEF("flowOut", "flowSideB", dlit.MustNew(11)),
				rule.NewMulLEF("flowOut", "flowSideB", dlit.MustNew(11.44)),
				rule.NewMulLEF("flowOut", "flowSideB", dlit.MustNew(3.08)),
				rule.NewMulLEF("flowOut", "flowSideB", dlit.MustNew(3.52)),
				rule.NewMulLEF("flowOut", "flowSideB", dlit.MustNew(3.96)),
				rule.NewMulLEF("flowOut", "flowSideB", dlit.MustNew(4.4)),
				rule.NewMulLEF("flowOut", "flowSideB", dlit.MustNew(4.84)),
				rule.NewMulLEF("flowOut", "flowSideB", dlit.MustNew(5.28)),
				rule.NewMulLEF("flowOut", "flowSideB", dlit.MustNew(5.72)),
				rule.NewMulLEF("flowOut", "flowSideB", dlit.MustNew(6.16)),
				rule.NewMulLEF("flowOut", "flowSideB", dlit.MustNew(6.6)),
				rule.NewMulLEF("flowOut", "flowSideB", dlit.MustNew(7.04)),
				rule.NewMulLEF("flowOut", "flowSideB", dlit.MustNew(7.48)),
				rule.NewMulLEF("flowOut", "flowSideB", dlit.MustNew(7.92)),
				rule.NewMulLEF("flowOut", "flowSideB", dlit.MustNew(8.36)),
				rule.NewMulLEF("flowOut", "flowSideB", dlit.MustNew(8.8)),
				rule.NewMulLEF("flowOut", "flowSideB", dlit.MustNew(9.24)),
				rule.NewMulLEF("flowOut", "flowSideB", dlit.MustNew(9.68)),

				rule.NewMulGEF("flowOut", "flowSideB", dlit.MustNew(10.12)),
				rule.NewMulGEF("flowOut", "flowSideB", dlit.MustNew(10.56)),
				rule.NewMulGEF("flowOut", "flowSideB", dlit.MustNew(11)),
				rule.NewMulGEF("flowOut", "flowSideB", dlit.MustNew(11.44)),
				rule.NewMulGEF("flowOut", "flowSideB", dlit.MustNew(3.08)),
				rule.NewMulGEF("flowOut", "flowSideB", dlit.MustNew(3.52)),
				rule.NewMulGEF("flowOut", "flowSideB", dlit.MustNew(3.96)),
				rule.NewMulGEF("flowOut", "flowSideB", dlit.MustNew(4.4)),
				rule.NewMulGEF("flowOut", "flowSideB", dlit.MustNew(4.84)),
				rule.NewMulGEF("flowOut", "flowSideB", dlit.MustNew(5.28)),
				rule.NewMulGEF("flowOut", "flowSideB", dlit.MustNew(5.72)),
				rule.NewMulGEF("flowOut", "flowSideB", dlit.MustNew(6.16)),
				rule.NewMulGEF("flowOut", "flowSideB", dlit.MustNew(6.6)),
				rule.NewMulGEF("flowOut", "flowSideB", dlit.MustNew(7.04)),
				rule.NewMulGEF("flowOut", "flowSideB", dlit.MustNew(7.48)),
				rule.NewMulGEF("flowOut", "flowSideB", dlit.MustNew(7.92)),
				rule.NewMulGEF("flowOut", "flowSideB", dlit.MustNew(8.36)),
				rule.NewMulGEF("flowOut", "flowSideB", dlit.MustNew(8.8)),
				rule.NewMulGEF("flowOut", "flowSideB", dlit.MustNew(9.24)),
				rule.NewMulGEF("flowOut", "flowSideB", dlit.MustNew(9.68)),
			},
		},
		{field: "flowSideA",
			complexities: []int{10},
			want:         []rule.Rule{},
		},
		{field: "flowSideB",
			complexities: []int{10},
			want:         []rule.Rule{},
		},
	}
	ruleFields := []string{"flowIn", "flowOut", "flowSideA", "flowSideB"}

	for i, c := range cases {
		for _, complexity := range c.complexities {
			got := generateMulRules(inputDescription, ruleFields, complexity, c.field)
			if err := matchRulesUnordered(got, c.want); err != nil {
				gotRuleStrs := rulesToSortedStrings(got)
				wantRuleStrs := rulesToSortedStrings(c.want)
				t.Errorf(
					"(%d) matchRulesUnordered() rules don't match: %s\ngot: %s\nwant: %s\n",
					i, err, gotRuleStrs, wantRuleStrs)
			}
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

func TestGenerateCompareStringRules(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"band": &description.Field{
				Kind:   fieldtype.Number,
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
	complexity := 10
	for _, c := range cases {
		got := generateCompareStringRules(
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

func TestGenerateInRules_1(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"band": &description.Field{
				Kind:   fieldtype.Number,
				Min:    dlit.MustNew(1),
				Max:    dlit.MustNew(3),
				Values: map[string]description.Value{},
			},
			"flow": &description.Field{
				Kind:   fieldtype.Number,
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
		field        string
		complexities []int
		want         []rule.Rule
	}{
		{field: "band",
			complexities: []int{10},
			want:         []rule.Rule{},
		},
		{field: "flow",
			complexities: []int{10},
			want:         []rule.Rule{},
		},
		{field: "groupA",
			complexities: []int{10},
			want:         []rule.Rule{},
		},
		{field: "groupB",
			complexities: []int{1, 2, 3, 4, 5, 6},
			want:         []rule.Rule{},
		},
		{field: "groupC",
			complexities: []int{1, 2, 3, 4},
			want:         []rule.Rule{},
		},
		{field: "groupC",
			complexities: []int{5, 6},
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
			complexities: []int{5, 6},
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
			complexities: []int{5, 6},
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
		for _, complexity := range c.complexities {
			got := generateInRules(inputDescription, ruleFields, complexity, c.field)
			if err := matchRulesUnordered(got, c.want); err != nil {
				gotRuleStrs := rulesToSortedStrings(got)
				wantRuleStrs := rulesToSortedStrings(c.want)
				t.Errorf("matchRulesUnordered() rules don't match: %s\ngot: %s\nwant: %s\n",
					err, gotRuleStrs, wantRuleStrs)
			}
		}
	}
}

// Test that will generate correct number of values in In rule for complexity
func TestGenerateInRules_2(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"group": &description.Field{
				Kind: fieldtype.String,
			},
		},
	}
	cases := []struct {
		groupValues      map[string]description.Value
		complexities     []int
		wantMinNumRules  int
		wantMaxNumRules  int
		wantMaxNumValues int
	}{
		{groupValues: map[string]description.Value{
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
			complexities:     []int{1, 2, 3, 4},
			wantMinNumRules:  0,
			wantMaxNumRules:  0,
			wantMaxNumValues: 5,
		},
		{groupValues: map[string]description.Value{
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
			complexities:     []int{5, 6},
			wantMinNumRules:  1000,
			wantMaxNumRules:  2000,
			wantMaxNumValues: 5,
		},
	}
	ruleFields := []string{"group"}
	for _, c := range cases {
		for _, complexity := range c.complexities {
			inputDescription.Fields["group"].Values = c.groupValues
			got := generateInRules(inputDescription, ruleFields, complexity, "group")
			if len(got) < c.wantMinNumRules || len(got) > c.wantMaxNumRules {
				t.Errorf("generateInRules: got wrong number of rules: %d", len(got))
			}
			for _, r := range got {
				numValues := strings.Count(r.String(), ",")
				if numValues < 2 || numValues > c.wantMaxNumValues {
					t.Errorf("generateInRules: wrong number of values in rule: %s", r)
				}
			}
		}
	}
}

func TestGenerateValueRules(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"band": &description.Field{
				Kind: fieldtype.Number,
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
				Kind:  fieldtype.Number,
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
				rule.NewEQFV("band", dlit.MustNew(1)),
				rule.NewNEFV("band", dlit.MustNew(1)),
				rule.NewEQFV("band", dlit.MustNew(3)),
				rule.NewNEFV("band", dlit.MustNew(3)),
				rule.NewEQFV("band", dlit.MustNew(4)),
				rule.NewNEFV("band", dlit.MustNew(4)),
			},
		},
		{field: "flow",
			want: []rule.Rule{
				rule.NewEQFV("flow", dlit.MustNew(1)),
				rule.NewNEFV("flow", dlit.MustNew(1)),
				rule.NewEQFV("flow", dlit.MustNew(3.37)),
				rule.NewNEFV("flow", dlit.MustNew(3.37)),
				rule.NewEQFV("flow", dlit.MustNew(3.3)),
				rule.NewNEFV("flow", dlit.MustNew(3.3)),
				rule.NewEQFV("flow", dlit.MustNew(4)),
				rule.NewNEFV("flow", dlit.MustNew(4)),
			},
		},
		{field: "group",
			want: []rule.Rule{
				rule.NewEQFV("group", dlit.MustNew("Nelson")),
				rule.NewNEFV("group", dlit.MustNew("Nelson")),
				rule.NewEQFV("group", dlit.MustNew("Drake")),
				rule.NewNEFV("group", dlit.MustNew("Drake")),
			},
		},
	}
	ruleFields := []string{"band", "flow", "group"}
	complexity := 5
	for _, c := range cases {
		got := generateValueRules(inputDescription, ruleFields, complexity, c.field)
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
