/*************************
 *  Test helper functions
 *************************/
package rule

import (
	"errors"
	"fmt"
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/internal"
	"github.com/vlifesystems/rhkit/internal/dexprfuncs"
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

func isRuleInRange(rule Rule, min *dlit.Literal, max *dlit.Literal) bool {
	switch x := rule.(type) {
	case *BetweenFV:
		return isInRange(x.Min(), min, max) && isInRange(x.Max(), min, max)
	case *OutsideFV:
		return isInRange(x.Low(), min, max) && isInRange(x.High(), min, max)
	case Valuer:
		return isInRange(x.Value(), min, max)
	}
	return true
}

func isInRange(v, min, max *dlit.Literal) bool {
	vars := map[string]*dlit.Literal{
		"min": min,
		"max": max,
		"v":   v,
	}
	inRange, err :=
		dexpr.EvalBool("v >= min && v <= max", dexprfuncs.CallFuncs, vars)
	return inRange && err == nil
}

func checkMaxDP(rules []Rule, maxDP int) error {
	if len(rules) == 0 {
		return nil
	}
	numMaxDP := 0
	for _, r := range rules {
		switch x := r.(type) {
		case *BetweenFV:
			numMinDP := internal.NumDecPlaces(x.Min().String())
			if numMinDP == maxDP {
				numMaxDP++
			} else if numMinDP > maxDP {
				return fmt.Errorf("rule has too many d.p., got:%d, want %d <= %d",
					numMinDP, numMinDP, maxDP)
			}
			numMaxDP := internal.NumDecPlaces(x.Max().String())
			if numMaxDP == maxDP {
				numMaxDP++
			} else if numMaxDP > maxDP {
				return fmt.Errorf("rule has too many d.p., got:%d, want %d <= %d",
					numMaxDP, numMaxDP, maxDP)
			}
		case *OutsideFV:
			numLowDP := internal.NumDecPlaces(x.Low().String())
			if numLowDP == maxDP {
				numMaxDP++
			} else if numLowDP > maxDP {
				return fmt.Errorf("rule has too many d.p., got:%d, want %d <= %d",
					numLowDP, numLowDP, maxDP)
			}
			numHighDP := internal.NumDecPlaces(x.High().String())
			if numHighDP == maxDP {
				numHighDP++
			} else if numHighDP > maxDP {
				return fmt.Errorf("rule has too many d.p., got:%d, want %d <= %d",
					numHighDP, numHighDP, maxDP)
			}
		case Valuer:
			numDP := internal.NumDecPlaces(x.Value().String())
			if numDP == maxDP {
				numMaxDP++
			} else if numDP > maxDP {
				return fmt.Errorf("rule has too many d.p., got:%d, want %d <= %d",
					numDP, numDP, maxDP)
			}
		}
	}
	if numMaxDP == 0 {
		return fmt.Errorf("there were no rules that had maxDP: %d", maxDP)
	}
	return nil
}

func valuesAreSpread(rules []Rule, mid *dlit.Literal) bool {
	if len(rules) < 3 {
		return true
	}
	values := []*dlit.Literal{}
	for _, r := range rules {
		switch x := r.(type) {
		case *BetweenFV:
			values = append(values, x.Min())
			values = append(values, x.Max())
		case *OutsideFV:
			values = append(values, x.Low())
			values = append(values, x.High())
		case Valuer:
			values = append(values, x.Value())
		}
	}

	isBelowMidExpr := dexpr.MustNew("v < mid", dexprfuncs.CallFuncs)
	isAboveMidExpr := dexpr.MustNew("v > mid", dexprfuncs.CallFuncs)
	numBelowMid := 0
	numAboveMid := 0
	vars := map[string]*dlit.Literal{"mid": mid}
	for _, v := range values {
		vars["v"] = v
		isBelowMid, err := isBelowMidExpr.EvalBool(vars)
		if err != nil {
			panic(err)
		}
		isAboveMid, err := isAboveMidExpr.EvalBool(vars)
		if err != nil {
			panic(err)
		}
		if isBelowMid {
			numBelowMid++
		} else if isAboveMid {
			numAboveMid++
		}
	}
	if numBelowMid == 0 || numAboveMid == 0 {
		return false
	}
	return true
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
	for _, r := range rules {
		if !isRuleInRange(r, min, max) {
			return fmt.Errorf("rule value isn't in range: %s", r)
		}
		if err := complyFunc(r); err != nil {
			return err
		}
	}
	if !valuesAreSpread(rules, mid) {
		return fmt.Errorf("rules are not evenly spread: %s", rules)
	}
	if err := checkMaxDP(rules, maxDP); err != nil {
		return err
	}
	return nil
}

func matchRulesUnordered(rules1 []Rule, rules2 []Rule) error {
	if len(rules1) != len(rules2) {
		return errors.New("rules different lengths")
	}
	return rulesContain(rules1, rules2)
}

func rulesContain(gotRules []Rule, wantRules []Rule) error {
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
