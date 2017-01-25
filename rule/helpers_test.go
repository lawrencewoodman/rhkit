/*************************
 *  Test helper functions
 *************************/
package rule

import (
	"fmt"
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
