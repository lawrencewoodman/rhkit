/*************************
 *  Test helper functions
 *************************/
package rule

import (
	"fmt"
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/internal"
	"github.com/vlifesystems/rhkit/internal/dexprfuncs"
	"math"
)

func checkErrorMatch(got, want error) error {
	switch x := want.(type) {
	case InvalidRuleError:
		gerr, ok := got.(InvalidRuleError)
		if !ok {
			return fmt.Errorf("got err type: %T, want error type: InvalidRuleError",
				got)
		}
		if x.Rule.String() != gerr.Rule.String() {
			return fmt.Errorf("got Rule: %s, Rule: %s", gerr.Rule, x.Rule)
		}
		return nil
	case IncompatibleTypesRuleError:
		gerr, ok := got.(IncompatibleTypesRuleError)
		if !ok {
			return fmt.Errorf("got err type: %T, want error type: IncompatibleTypesRuleError",
				got)
		}
		if x.Rule.String() != gerr.Rule.String() {
			return fmt.Errorf("got Rule: %s, Rule: %s", gerr.Rule, x.Rule)
		}
		return nil
	}
	if got.Error() != want.Error() {
		return fmt.Errorf("got err: %v, want err: %v", got, want)
	}
	return nil
}

func checkRulesMatch(got, want []Rule) error {
	if len(got) != len(want) {
		return fmt.Errorf("len(got): %d != len(want): %d", len(got), len(want))
	}
	for i, r := range want {
		if got[i].String() != r.String() {
			return fmt.Errorf("got != want, got[%d]: %s, want[%d]: %s", i, got[i], i, r)
		}
	}
	return nil
}

func checkRulesComply(
	rules []Rule,
	minNumRules int,
	maxNumRules int,
	min *dlit.Literal,
	max *dlit.Literal,
	mid *dlit.Literal,
	maxDP int,
	complyFunc func(Rule) error,
) error {
	if len(rules) == 0 && maxNumRules == 0 {
		return nil
	}
	if len(rules) < minNumRules || len(rules) > maxNumRules {
		return fmt.Errorf("len(rules): %d, want %d >= %d && %d <= %d",
			len(rules), len(rules), minNumRules, len(rules), maxNumRules)
	}
	uRules := Uniq(rules)
	if len(uRules) != len(rules) {
		return fmt.Errorf("len(rules): %d, %d aren't unique", len(rules),
			len(rules)-len(uRules))
	}
	numMaxDP := 0
	numBelowMid := 0
	numAboveMid := 0
	for _, r := range rules {
		if x, ok := r.(Valuer); ok {
			v := x.Value()
			vars := map[string]*dlit.Literal{
				"min": min,
				"max": max,
				"mid": mid,
				"v":   v,
			}
			inRange, err :=
				dexpr.EvalBool("v >= min && v <= max", dexprfuncs.CallFuncs, vars)
			if !inRange || err != nil {
				return fmt.Errorf(
					"rule value isn't in range, got: %s, want: %s >= %s && %s <= %s (%s)",
					v,
					v,
					min,
					v,
					max,
					r,
				)
			}

			isBelowMid, err := dexpr.EvalBool("v < mid", dexprfuncs.CallFuncs, vars)
			if err != nil {
				panic(err)
			}
			isAboveMid, err := dexpr.EvalBool("v > mid", dexprfuncs.CallFuncs, vars)
			if err != nil {
				panic(err)
			}
			if isBelowMid {
				numBelowMid++
			} else if isAboveMid {
				numAboveMid++
			}

			if err := complyFunc(r); err != nil {
				return err
			}
		}
		numDP := internal.NumDecPlaces(r.String())
		if numDP == maxDP {
			numMaxDP++
		} else if numDP > maxDP {
			return fmt.Errorf("rule has too many d.p., got:%d, want %d <= %d",
				numDP, numDP, maxDP)
		}
	}

	if math.Abs(float64(numAboveMid-numBelowMid)) > 3 {
		return fmt.Errorf("rules are not evenly spread: %s", rules)
	}

	if numMaxDP == 0 {
		return fmt.Errorf("there were no rules that had maxDP: %d", maxDP)
	}
	return nil
}
