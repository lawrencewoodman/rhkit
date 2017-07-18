package rule

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
	"github.com/vlifesystems/rhkit/internal/testhelpers"
	"testing"
)

func TestTweak_1(t *testing.T) {
	testPurposes := []string{"Ensure that results are only from tweakable rules"}
	inputDescription := &description.Description{
		map[string]*description.Field{
			"band": {
				fieldtype.Number, dlit.MustNew(3), dlit.MustNew(40), 0,
				map[string]description.Value{}, 0},
			"age": {
				fieldtype.Number, dlit.MustNew(4), dlit.MustNew(90), 0,
				map[string]description.Value{}, 0},
			"flow": {
				fieldtype.Number, dlit.MustNew(50), dlit.MustNew(400), 2,
				map[string]description.Value{}, 0},
		}}
	rulesIn := []Rule{
		NewGEFV("band", dlit.MustNew(4)),
		NewGEFV("band", dlit.MustNew(20)),
		NewGTFF("band", "team"),
		NewGEFV("age", dlit.MustNew(7)),
		NewGEFV("age", dlit.MustNew(8)),
		MustNewBetweenFV("age", dlit.MustNew(21), dlit.MustNew(39)),
		NewGEFV("flow", dlit.MustNew(60.7)),
		NewGEFV("flow", dlit.MustNew(70.20)),
		NewGEFV("flow", dlit.MustNew(100.5)),
		NewGTFF("age", "band"),
		NewInFV(
			"stage",
			testhelpers.MakeStringsDlitSlice("20", "21", "22"),
		),
	}
	gotRules := Tweak(1, rulesIn, inputDescription)

	numAgeGERules := 0
	numAgeBetweenRules := 0
	numBandGERules := 0
	numFlowGERules := 0
	numOtherRules := 0
	for _, gotRule := range gotRules {
		switch x := gotRule.(type) {
		case True:
			continue
		case *GEFV:
			field := x.Fields()[0]
			if field == "band" {
				numBandGERules++
			} else if field == "age" {
				numAgeGERules++
			}
			if x.Fields()[0] == "flow" {
				numFlowGERules++
			}
		case *BetweenFV:
			if x.Fields()[0] == "age" {
				numAgeBetweenRules++
			}
		case Tweaker:
			numOtherRules++
		default:
			printTestPurposes(t, testPurposes)
			t.Fatalf("Tweak(%s) rule isn't tweakable: %s", rulesIn, gotRule)
		}
	}

	if numBandGERules < 7 {
		printTestPurposes(t, testPurposes)
		t.Errorf("Tweak(%v) wrong number(%d) of rules: band >= ? - got: %v",
			rulesIn, numBandGERules, gotRules)
	}
	if numAgeGERules < 4 {
		printTestPurposes(t, testPurposes)
		t.Errorf("Tweak(%v) wrong number(%d) of rules: age >= ? - got: %v",
			rulesIn, numAgeGERules, gotRules)
	}
	if numFlowGERules < 9 {
		printTestPurposes(t, testPurposes)
		t.Errorf("Tweak(%v) wrong number(%d) of rules: flow >= ? - got: %v",
			rulesIn, numFlowGERules, gotRules)
	}
	if numAgeBetweenRules < 6 {
		printTestPurposes(t, testPurposes)
		t.Errorf("Tweak(%v) wrong number(%d) of rules: age >= ? && age <= ?- got: %v",
			rulesIn, numAgeBetweenRules, gotRules)
	}
	if numOtherRules != 0 {
		printTestPurposes(t, testPurposes)
		t.Errorf("Tweak(%v) wrong number(%d) of other rules - got: %v",
			rulesIn, numOtherRules, gotRules)
	}
}

func TestTweak_2(t *testing.T) {
	testPurposes := []string{
		"Ensure that generates a range of int numbers between current ones",
	}
	inputDescription := &description.Description{
		map[string]*description.Field{
			"age": {
				fieldtype.Number, dlit.MustNew(10), dlit.MustNew(80), 0,
				map[string]description.Value{}, 0,
			},
		}}
	rulesIn := []Rule{
		NewLEFV("age", dlit.MustNew(40)),
		NewLEFV("age", dlit.MustNew(20)),
		NewLEFV("age", dlit.MustNew(50)),
		NewLEFV("age", dlit.MustNew(60)),
	}
	gotRules := Tweak(1, rulesIn, inputDescription)

	num10To20 := 0
	num20To40 := 0
	num40To50 := 0
	num50To80 := 0
	numOther := 0
	for _, gotRule := range gotRules {
		switch x := gotRule.(type) {
		case True:
			continue
		case *LEFV:
			if x.Fields()[0] != "age" {
				printTestPurposes(t, testPurposes)
				t.Fatalf("Tweak(%s) invalid rule(%s): ", rulesIn, gotRule)
			}
			n := x.Value()
			nInt, _ := n.Int()
			if nInt >= 10 && nInt < 20 {
				num10To20++
			} else if nInt >= 20 && nInt < 40 {
				num20To40++
			} else if nInt >= 40 && nInt < 50 {
				num40To50++
			} else if nInt >= 50 && nInt < 80 {
				num50To80++
			} else {
				numOther++
			}
		case Tweaker:
			continue
		default:
			printTestPurposes(t, testPurposes)
			t.Fatalf("Tweak(%s) invalid rule(%s)", rulesIn, gotRule)
		}
	}

	if num10To20 < 6 {
		printTestPurposes(t, testPurposes)
		t.Errorf("Tweak(%v) wrong number(%d) of rules 10 to 20, got: %v",
			rulesIn, num10To20, gotRules)
	}
	if num20To40 < 9 {
		printTestPurposes(t, testPurposes)
		t.Errorf("Tweak(%v) wrong number(%d) of rules 20 to 40, got: %v",
			rulesIn, num20To40, gotRules)
	}
	if num40To50 < 9 {
		printTestPurposes(t, testPurposes)
		t.Errorf("Tweak(%v) wrong number(%d) of rules 40 to 50, got: %v",
			rulesIn, num40To50, gotRules)
	}
	if num50To80 < 9 {
		printTestPurposes(t, testPurposes)
		t.Errorf("Tweak(%v) wrong number(%d) of rules 50 to 80, got: %v",
			rulesIn, num50To80, gotRules)
	}
	if numOther != 0 {
		printTestPurposes(t, testPurposes)
		t.Errorf("Tweak(%v) wrong number(%d) of other rules - got: %v",
			rulesIn, numOther, gotRules)
	}
}

func TestTweak_3(t *testing.T) {
	testPurposes := []string{
		"Ensure that generates a range of float numbers between current ones",
		"Ensure that decimal places are no greater than maxDP for field",
	}
	inputDescription := &description.Description{
		map[string]*description.Field{
			"flow": {
				Kind:      fieldtype.Number,
				Min:       dlit.MustNew(10),
				Max:       dlit.MustNew(80),
				MaxDP:     6,
				Values:    map[string]description.Value{},
				NumValues: 0,
			},
		}}
	rulesIn := []Rule{
		NewLEFV("flow", dlit.MustNew(40.78234)),
		NewLEFV("flow", dlit.MustNew(24.89)),
		NewLEFV("flow", dlit.MustNew(52.604956)),
		NewLEFV("flow", dlit.MustNew(65.80)),
	}
	wantMaxDP := inputDescription.Fields["flow"].MaxDP
	wantMinDP := 0
	gotRules := Tweak(1, rulesIn, inputDescription)

	num10To24 := 0
	num24To41 := 0
	num41To53 := 0
	num53To80 := 0
	numOther := 0
	gotMaxDP := 0
	gotMinDP := 100
	for _, gotRule := range gotRules {
		switch x := gotRule.(type) {
		case True:
			continue
		case *LEFV:
			if x.Fields()[0] != "flow" {
				printTestPurposes(t, testPurposes)
				t.Fatalf("Tweak(%s) invalid rule(%s)", rulesIn, gotRule)
			}
			n := x.Value()
			nFloat, _ := n.Float()
			if nFloat >= 10 && nFloat < 24 {
				num10To24++
			} else if nFloat >= 24 && nFloat < 41 {
				num24To41++
			} else if nFloat >= 41 && nFloat < 53 {
				num41To53++
			} else if nFloat >= 53 && nFloat < 80 {
				num53To80++
			} else {
				numOther++
			}
			valueDP := internal.NumDecPlaces(x.Value().String())
			if valueDP > gotMaxDP {
				gotMaxDP = valueDP
			}
			if valueDP < gotMinDP {
				gotMinDP = valueDP
			}
		default:
			printTestPurposes(t, testPurposes)
			t.Fatalf("Tweak(%s) invalid rule(%s)",
				rulesIn, gotRule)
		}
	}

	if num10To24 < 8 {
		printTestPurposes(t, testPurposes)
		t.Errorf("Tweak(%v) wrong number(%d) of rules 10 to 24, got: %v",
			rulesIn, num10To24, gotRules)
	}
	if num24To41 < 9 {
		printTestPurposes(t, testPurposes)
		t.Errorf("Tweak(%v) wrong number(%d) of rules 24 to 41, got: %v",
			rulesIn, num24To41, gotRules)
	}
	if num41To53 < 9 {
		printTestPurposes(t, testPurposes)
		t.Errorf("Tweak(%v) wrong number(%d) of rules 41 to 53, got: %v",

			rulesIn, num41To53, gotRules)
	}
	if num53To80 < 9 {
		printTestPurposes(t, testPurposes)
		t.Errorf("Tweak(%v) wrong number(%d) of rules 53 to 80, got: %v",

			rulesIn, num53To80, gotRules)
	}

	if numOther != 0 {
		printTestPurposes(t, testPurposes)
		t.Errorf("Tweak(%v) wrong number(%d) of other rules - got: %v",
			rulesIn, numOther, gotRules)
	}

	if gotMinDP != wantMinDP {
		printTestPurposes(t, testPurposes)
		t.Errorf("Tweak(%v) maxDP for rules to big got, %d, want: %d, rules: %v",
			rulesIn, gotMinDP, wantMinDP, gotRules)
	}
	if gotMaxDP != wantMaxDP {
		printTestPurposes(t, testPurposes)
		t.Errorf("Tweak(%v) maxDP for rules to big got, %d, want: %d, rules: %v",
			rulesIn, gotMaxDP, wantMaxDP, gotRules)
	}
}

func TestTweak_4(t *testing.T) {
	testPurposes := []string{
		"Ensure that generates a True rule",
	}
	inputDescription := &description.Description{
		map[string]*description.Field{
			"flow": {
				fieldtype.Number, dlit.MustNew(4), dlit.MustNew(30), 6,
				map[string]description.Value{}, 0,
			},
		}}
	rulesIn := []Rule{
		NewLEFV("flow", dlit.MustNew(40.78234)),
		NewLEFV("flow", dlit.MustNew(24.89)),
		NewLEFV("flow", dlit.MustNew(52.604956)),
		NewTrue(),
	}

	gotRules := Tweak(1, rulesIn, inputDescription)
	trueRuleFound := false
	for _, r := range gotRules {
		if _, ruleIsTrue := r.(True); ruleIsTrue {
			trueRuleFound = true
			break
		}
	}
	if !trueRuleFound {
		printTestPurposes(t, testPurposes)
		t.Errorf("Tweak(%s)  - No 'true' rule found", rulesIn)
	}
}

func TestTweak_5(t *testing.T) {
	testPurposes := []string{"Ensure that are rules are unique"}
	inputDescription := &description.Description{
		map[string]*description.Field{
			"band": {
				fieldtype.Number, dlit.MustNew(3), dlit.MustNew(40), 0,
				map[string]description.Value{}, 0},
			"age": {
				fieldtype.Number, dlit.MustNew(4), dlit.MustNew(30), 0,
				map[string]description.Value{}, 0},
			"flow": {
				fieldtype.Number, dlit.MustNew(50), dlit.MustNew(400), 2,
				map[string]description.Value{}, 0},
		}}
	rulesIn := []Rule{
		NewGEFV("band", dlit.MustNew(4)),
		NewGEFV("band", dlit.MustNew(5)),
		NewGEFV("band", dlit.MustNew(6)),
		NewGEFV("band", dlit.MustNew(20)),
		NewGTFF("band", "team"),
		NewGEFV("age", dlit.MustNew(7)),
		NewGEFV("age", dlit.MustNew(8)),
		NewGEFV("flow", dlit.MustNew(60.7)),
		NewGEFV("flow", dlit.MustNew(70.20)),
		NewGEFV("flow", dlit.MustNew(100.5)),
		NewGTFF("age", "band"),
		NewInFV(
			"stage",
			testhelpers.MakeStringsDlitSlice("20", "21", "22"),
		),
	}
	gotRules := Tweak(1, rulesIn, inputDescription)

	for _, gotRule := range gotRules {
		count := 0
		for _, r := range gotRules {
			if gotRule.String() == r.String() {
				count++
				if count > 1 {
					printTestPurposes(t, testPurposes)
					t.Fatalf("Tweak - rule isn't unique: %s", gotRule)
				}
			}
		}
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
