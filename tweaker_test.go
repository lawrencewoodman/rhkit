package main

import (
	"github.com/lawrencewoodman/dlit_go"
	"github.com/lawrencewoodman/rulehunter/internal"
	"testing"
)

func TestTweakRules_1(t *testing.T) {
	testPurposes := []string{"Ensure that results are only from tweakable rules"}
	fieldDescriptions := map[string]*FieldDescription{
		"band": &FieldDescription{INT, dlit.MustNew(3), dlit.MustNew(40), 0,
			[]*dlit.Literal{}, 0},
		"age": &FieldDescription{INT, dlit.MustNew(4), dlit.MustNew(30), 0,
			[]*dlit.Literal{}, 0},
		"flow": &FieldDescription{FLOAT, dlit.MustNew(50), dlit.MustNew(400), 2,
			[]*dlit.Literal{}, 0},
	}
	rulesIn := []*internal.Rule{
		internal.MustNewRule("band > 4"),
		internal.MustNewRule("band > 20"),
		internal.MustNewRule("band > team"),
		internal.MustNewRule("age > 7"),
		internal.MustNewRule("age >= 8"),
		internal.MustNewRule("flow >= 60.7"),
		internal.MustNewRule("flow >= 70.20"),
		internal.MustNewRule("flow > 100.5"),
		internal.MustNewRule("age > band"),
		internal.MustNewRule("in(stage,\"20\",\"21\",\"22\")"),
	}
	gotRules := TweakRules(rulesIn, fieldDescriptions)

	numBandGtRules := 0
	numFlowGeqRules := 0
	numOtherRules := 0
	for _, rule := range gotRules {
		if rule.String() == "true()" {
			continue
		}
		isTweakable, field, operator, _ := rule.GetTweakableParts()
		if !isTweakable {
			printTestPurposes(t, testPurposes)
			t.Errorf("TweakRules(%s) rule isn't tweakable: %s", rulesIn, rule)
		}

		if field == "band" && operator == ">" {
			numBandGtRules++
		} else if field == "flow" && operator == ">=" {
			numFlowGeqRules++
		} else {
			numOtherRules++
		}
	}

	if numBandGtRules < 10 {
		printTestPurposes(t, testPurposes)
		t.Errorf("TweakRules(%q) wrong number(%d) of rules: band > ? - got: %q",
			rulesIn, numBandGtRules, gotRules)
	}
	if numFlowGeqRules < 9 {
		printTestPurposes(t, testPurposes)
		t.Errorf("TweakRules(%q) wrong number(%d) of rules: flow >= ? - got: %q",
			rulesIn, numFlowGeqRules, gotRules)
	}
	if numOtherRules != 0 {
		printTestPurposes(t, testPurposes)
		t.Errorf("TweakRules(%q) wrong number(%d) of other rules - got: %q",
			numOtherRules, gotRules)
	}
}

func TestTweakRules_2(t *testing.T) {
	testPurposes := []string{
		"Ensure that generates a range of int numbers between current ones",
		"Ensure only operates on first 3 in group",
	}
	fieldDescriptions := map[string]*FieldDescription{
		"age": &FieldDescription{INT, dlit.MustNew(20), dlit.MustNew(40), 0,
			[]*dlit.Literal{}, 0},
	}
	rulesIn := []*internal.Rule{
		internal.MustNewRule("age <= 40"),
		internal.MustNewRule("age <= 20"),
		internal.MustNewRule("age <= 50"),
		internal.MustNewRule("age <= 60"),
	}
	gotRules := TweakRules(rulesIn, fieldDescriptions)

	num20To40 := 0
	num40To50 := 0
	numOther := 0
	for _, rule := range gotRules {
		if rule.String() == "true()" {
			continue
		}
		isTweakable, field, operator, value := rule.GetTweakableParts()
		if !isTweakable && field != "age" && operator != "<=" {
			printTestPurposes(t, testPurposes)
			t.Errorf("TweakRules(%s) invalid rule(%s): isTweakable: %s, field: %s, operator: %s",
				rulesIn, rule, isTweakable, field, operator)
		}
		l := dlit.MustNew(value)
		n, nIsInt := l.Int()
		if !nIsInt {
			printTestPurposes(t, testPurposes)
			t.Errorf("TweakRules(%s) invalid rule(%s): value isn't int", rulesIn, rule)
		} else if n >= 20 && n < 40 {
			num20To40++
		} else if n >= 40 && n < 50 {
			num40To50++
		} else {
			numOther++
		}
	}

	if num20To40 < 9 {
		printTestPurposes(t, testPurposes)
		t.Errorf("TweakRules(%q) wrong number(%d) of rules 20 to 40, got: %q",
			rulesIn, num20To40, gotRules)
	}
	if num40To50 < 9 {
		printTestPurposes(t, testPurposes)
		t.Errorf("TweakRules(%q) wrong number(%d) of rules 40 to 50, got: %q",
			rulesIn, num40To50, gotRules)
	}
	if numOther != 0 {
		printTestPurposes(t, testPurposes)
		t.Errorf("TweakRules(%q) wrong number(%d) of other rules - got: %q",
			numOther, gotRules)
	}
}

func TestTweakRules_3(t *testing.T) {
	testPurposes := []string{
		"Ensure that generates a range of float numbers between current ones",
		"Ensure only operates on first 3 in group",
		"Ensure that decimal places are no greater than maxDP for field",
	}
	fieldDescriptions := map[string]*FieldDescription{
		"flow": &FieldDescription{FLOAT, dlit.MustNew(4), dlit.MustNew(30), 6,
			[]*dlit.Literal{}, 0},
	}
	rulesIn := []*internal.Rule{
		internal.MustNewRule("flow <= 40.78234"),
		internal.MustNewRule("flow <= 24.89"),
		internal.MustNewRule("flow <= 52.604956"),
		internal.MustNewRule("flow <= 65.80"),
	}
	wantMaxDP := fieldDescriptions["flow"].MaxDP
	wantMinDP := 0
	gotRules := TweakRules(rulesIn, fieldDescriptions)

	num24To41 := 0
	num41To53 := 0
	numOther := 0
	gotMaxDP := 0
	gotMinDP := 100
	for _, rule := range gotRules {
		if rule.String() == "true()" {
			continue
		}
		isTweakable, field, operator, value := rule.GetTweakableParts()
		if !isTweakable && field != "flow" && operator != "<=" {
			printTestPurposes(t, testPurposes)
			t.Errorf("TweakRules(%s) invalid rule(%s): isTweakable: %s, field: %s, operator: %s",
				rulesIn, rule, isTweakable, field, operator)
		}
		l := dlit.MustNew(value)
		n, nIsFloat := l.Float()
		if !nIsFloat {
			printTestPurposes(t, testPurposes)
			t.Errorf("TweakRules(%s) invalid rule(%s): value isn't float", rulesIn, rule)
		} else if n >= 24 && n < 41 {
			num24To41++
		} else if n >= 41 && n < 53 {
			num41To53++
		} else {
			numOther++
		}
		valueDP := numDecPlaces(value)
		if valueDP > gotMaxDP {
			gotMaxDP = valueDP
		}
		if valueDP < gotMinDP {
			gotMinDP = valueDP
		}
	}

	if num24To41 < 9 {
		printTestPurposes(t, testPurposes)
		t.Errorf("TweakRules(%q) wrong number(%d) of rules 24 to 41, got: %q",
			rulesIn, num24To41, gotRules)
	}
	if num41To53 < 9 {
		printTestPurposes(t, testPurposes)
		t.Errorf("TweakRules(%q) wrong number(%d) of rules 41 to 53, got: %q",

			rulesIn, num41To53, gotRules)
	}

	if numOther != 0 {
		printTestPurposes(t, testPurposes)
		t.Errorf("TweakRules(%q) wrong number(%d) of other rules - got: %q",
			rulesIn, numOther, gotRules)
	}

	if gotMinDP != wantMinDP {
		printTestPurposes(t, testPurposes)
		t.Errorf("TweakRules(%q) maxDP for rules to big got, %d, want: %d, rules: %q",
			rulesIn, gotMinDP, wantMinDP, gotRules)
	}
	if gotMaxDP != wantMaxDP {
		printTestPurposes(t, testPurposes)
		t.Errorf("TweakRules(%q) maxDP for rules to big got, %d, want: %d, rules: %q",
			rulesIn, gotMaxDP, wantMaxDP, gotRules)
	}
}

func TestTweakRules_4(t *testing.T) {
	testPurposes := []string{
		"Ensure that generates a 'true()' rule",
	}
	fieldDescriptions := map[string]*FieldDescription{
		"flow": &FieldDescription{FLOAT, dlit.MustNew(4), dlit.MustNew(30), 6,
			[]*dlit.Literal{}, 0},
	}
	rulesIn := []*internal.Rule{
		internal.MustNewRule("flow <= 40.78234"),
		internal.MustNewRule("flow <= 24.89"),
		internal.MustNewRule("flow <= 52.604956"),
		internal.MustNewRule("true()"),
	}

	gotRules := TweakRules(rulesIn, fieldDescriptions)
	trueRuleFound := false
	for _, rule := range gotRules {
		if rule.String() == "true()" {
			trueRuleFound = true
			break
		}
	}
	if !trueRuleFound {
		printTestPurposes(t, testPurposes)
		t.Errorf("TweakRules(%s)  - No 'true' rule found", rulesIn)
	}
}

/**************************************
 *    Helper functions
 **************************************/
func printTestPurposes(t *testing.T, testPurposes []string) {
	t.Errorf("Test: %s\n", testPurposes[0])
	for _, p := range testPurposes[1:] {
		t.Errorf("      %s\n", p)
	}
}
