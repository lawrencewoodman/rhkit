package rulehunter

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/rule"
	"testing"
)

func TestTweakRules_1(t *testing.T) {
	testPurposes := []string{"Ensure that results are only from tweakable rules"}
	inputDescription := &Description{
		map[string]*fieldDescription{
			"band": &fieldDescription{
				ftInt, dlit.MustNew(3), dlit.MustNew(40), 0,
				[]*dlit.Literal{}, 0},
			"age": &fieldDescription{
				ftInt, dlit.MustNew(4), dlit.MustNew(30), 0,
				[]*dlit.Literal{}, 0},
			"flow": &fieldDescription{
				ftFloat, dlit.MustNew(50), dlit.MustNew(400), 2,
				[]*dlit.Literal{}, 0},
		}}
	rulesIn := []rule.Rule{
		rule.NewGEFVI("band", 4),
		rule.NewGEFVI("band", 20),
		rule.NewGTFF("band", "team"),
		rule.NewGEFVI("age", 7),
		rule.NewGEFVI("age", 8),
		rule.NewGEFVF("flow", 60.7),
		rule.NewGEFVF("flow", 70.20),
		rule.NewGEFVF("flow", 100.5),
		rule.NewGTFF("age", "band"),
		rule.MustNewDRule("in(stage,\"20\",\"21\",\"22\")"),
	}
	gotRules := TweakRules(rulesIn, inputDescription)

	numBandGtRules := 0
	numFlowGeqRules := 0
	numOtherRules := 0
	for _, gotRule := range gotRules {
		if _, ruleIsTrue := gotRule.(rule.True); ruleIsTrue {
			continue
		}
		tRule, isTweakable := gotRule.(rule.TweakableRule)
		if !isTweakable {
			printTestPurposes(t, testPurposes)
			t.Fatalf("TweakRules(%s) rule isn't tweakable: %s", rulesIn, gotRule)
			continue
		}

		field, operator, _ := tRule.GetTweakableParts()

		if field == "band" && operator == ">=" {
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
			rulesIn, numOtherRules, gotRules)
	}
}

func TestTweakRules_2(t *testing.T) {
	testPurposes := []string{
		"Ensure that generates a range of int numbers between current ones",
		"Ensure only operates on first 3 in group",
	}
	inputDescription := &Description{
		map[string]*fieldDescription{
			"age": &fieldDescription{
				ftInt, dlit.MustNew(20), dlit.MustNew(40), 0,
				[]*dlit.Literal{}, 0,
			},
		}}
	rulesIn := []rule.Rule{
		rule.NewLEFVI("age", 40),
		rule.NewLEFVI("age", 20),
		rule.NewLEFVI("age", 50),
		rule.NewLEFVI("age", 60),
	}
	gotRules := TweakRules(rulesIn, inputDescription)

	num20To40 := 0
	num40To50 := 0
	numOther := 0
	for _, gotRule := range gotRules {
		if _, isRuleTrue := gotRule.(rule.True); isRuleTrue {
			continue
		}
		tRule, isTweakable := gotRule.(rule.TweakableRule)
		if !isTweakable {
			printTestPurposes(t, testPurposes)
			t.Fatalf("TweakRules(%s) invalid rule(%s): isTweakable: %t",
				rulesIn, gotRule, isTweakable)
		}
		field, operator, value := tRule.GetTweakableParts()
		if field != "age" && operator != "<=" {
			printTestPurposes(t, testPurposes)
			t.Fatalf("TweakRules(%s) invalid rule(%s): field: %s, operator: %s",
				rulesIn, gotRule, field, operator)
		}
		l := dlit.MustNew(value)
		n, nIsInt := l.Int()
		if !nIsInt {
			printTestPurposes(t, testPurposes)
			t.Errorf("TweakRules(%s) invalid rule(%s): value isn't int",
				rulesIn, gotRule)
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
			rulesIn, numOther, gotRules)
	}
}

func TestTweakRules_3(t *testing.T) {
	testPurposes := []string{
		"Ensure that generates a range of float numbers between current ones",
		"Ensure only operates on first 3 in group",
		"Ensure that decimal places are no greater than maxDP for field",
	}
	inputDescription := &Description{
		map[string]*fieldDescription{
			"flow": &fieldDescription{
				ftFloat, dlit.MustNew(4), dlit.MustNew(30), 6,
				[]*dlit.Literal{}, 0,
			},
		}}
	rulesIn := []rule.Rule{
		rule.NewLEFVF("flow", 40.78234),
		rule.NewLEFVF("flow", 24.89),
		rule.NewLEFVF("flow", 52.604956),
		rule.NewLEFVF("flow", 65.80),
	}
	wantMaxDP := inputDescription.fields["flow"].maxDP
	wantMinDP := 0
	gotRules := TweakRules(rulesIn, inputDescription)

	num24To41 := 0
	num41To53 := 0
	numOther := 0
	gotMaxDP := 0
	gotMinDP := 100
	for _, gotRule := range gotRules {
		if _, ruleIsTrue := gotRule.(rule.True); ruleIsTrue {
			continue
		}
		tRule, isTweakable := gotRule.(rule.TweakableRule)
		if !isTweakable {
			printTestPurposes(t, testPurposes)
			t.Fatalf("TweakRules(%s) invalid rule(%s): isTweakable: %t",
				rulesIn, gotRule, isTweakable)
		}
		field, operator, value := tRule.GetTweakableParts()
		if field != "flow" && operator != "<=" {
			printTestPurposes(t, testPurposes)
			t.Fatalf("TweakRules(%s) invalid rule(%s): field: %s, operator: %s",
				rulesIn, gotRule, field, operator)
		}
		l := dlit.MustNew(value)
		n, nIsFloat := l.Float()
		if !nIsFloat {
			printTestPurposes(t, testPurposes)
			t.Errorf("TweakRules(%s) invalid rule(%s): value isn't float",
				rulesIn, gotRule)
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
		"Ensure that generates a True rule",
	}
	inputDescription := &Description{
		map[string]*fieldDescription{
			"flow": &fieldDescription{
				ftFloat, dlit.MustNew(4), dlit.MustNew(30), 6,
				[]*dlit.Literal{}, 0,
			},
		}}
	rulesIn := []rule.Rule{
		rule.MustNewDRule("flow <= 40.78234"),
		rule.MustNewDRule("flow <= 24.89"),
		rule.MustNewDRule("flow <= 52.604956"),
		rule.NewTrue(),
	}

	gotRules := TweakRules(rulesIn, inputDescription)
	trueRuleFound := false
	for _, r := range gotRules {
		if _, ruleIsTrue := r.(rule.True); ruleIsTrue {
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
