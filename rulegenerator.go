/*
	Copyright (C) 2016-2017 vLife Systems Ltd <http://vlifesystems.com>
	This file is part of rhkit.

	rhkit is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	rhkit is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with rhkit; see the file COPYING.  If not, see
	<http://www.gnu.org/licenses/>.
*/

package rhkit

import (
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/rule"
)

// GenerateRules generates rules for the ruleFields.
// complexity is used to indicate how complex and in turn how many rules
// to produce it takes a number 1 to 10.
func GenerateRules(
	inputDescription *description.Description,
	ruleFields []string,
	complexity int,
) []rule.Rule {
	if complexity < 1 || complexity > 10 {
		panic("complexity must be in range 1..10")
	}
	rules := make([]rule.Rule, 1)
	rules[0] = rule.NewTrue()
	newRules := rule.Generate(inputDescription, ruleFields, complexity)
	rules = append(rules, newRules...)

	if len(ruleFields) == 2 {
		cRules := CombineRules(rules)
		rules = append(rules, cRules...)
	}
	rule.Sort(rules)
	return rules
}

func CombineRules(rules []rule.Rule) []rule.Rule {
	rule.Sort(rules)
	combinedRules := make([]rule.Rule, 0)
	numRules := len(rules)
	for i := 0; i < numRules-1; i++ {
		for j := i + 1; j < numRules; j++ {
			if andRule, err := rule.NewAnd(rules[i], rules[j]); err == nil {
				combinedRules = append(combinedRules, andRule)
			}
			if orRule, err := rule.NewOr(rules[i], rules[j]); err == nil {
				combinedRules = append(combinedRules, orRule)
			}
		}
	}
	return rule.Uniq(combinedRules)
}
