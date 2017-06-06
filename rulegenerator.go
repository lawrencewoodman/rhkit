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
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal"
	"github.com/vlifesystems/rhkit/internal/dexprfuncs"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
	"github.com/vlifesystems/rhkit/rule"
)

type ruleGeneratorFunc func(
	*description.Description,
	[]string,
	int,
	string,
) []rule.Rule

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
	ruleGenerators := []ruleGeneratorFunc{
		generateCompareNumericRules, generateCompareStringRules,
	}
	rules[0] = rule.NewTrue()
	for field := range inputDescription.Fields {
		if internal.StringInSlice(field, ruleFields) {
			for _, ruleGenerator := range ruleGenerators {
				newRules :=
					ruleGenerator(inputDescription, ruleFields, complexity, field)
				rules = append(rules, newRules...)
			}
		}
	}
	extraRules := rule.Generate(inputDescription, ruleFields, complexity)
	rules = append(rules, extraRules...)

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

func generateCompareNumericRules(
	inputDescription *description.Description,
	ruleFields []string,
	complexity int,
	field string,
) []rule.Rule {
	fd := inputDescription.Fields[field]
	if fd.Kind != fieldtype.Number {
		return []rule.Rule{}
	}
	fieldNum := description.CalcFieldNum(inputDescription.Fields, field)
	rulesMap := make(map[string]rule.Rule)
	ruleNewFuncs := []func(string, string) rule.Rule{
		rule.NewLTFF,
		rule.NewLEFF,
		rule.NewEQFF,
		rule.NewNEFF,
		rule.NewGEFF,
		rule.NewGTFF,
	}

	for oField, oFd := range inputDescription.Fields {
		oFieldNum := description.CalcFieldNum(inputDescription.Fields, oField)
		isComparable := hasComparableNumberRange(fd, oFd)
		if fieldNum < oFieldNum && isComparable &&
			internal.StringInSlice(oField, ruleFields) {
			for _, ruleNewFunc := range ruleNewFuncs {
				r := ruleNewFunc(field, oField)
				rulesMap[r.String()] = r
			}
		}
	}
	rules := rulesMapToArray(rulesMap)
	return rules
}

func generateCompareStringRules(
	inputDescription *description.Description,
	ruleFields []string,
	complexity int,
	field string,
) []rule.Rule {
	fd := inputDescription.Fields[field]
	if fd.Kind != fieldtype.String {
		return []rule.Rule{}
	}
	fieldNum := description.CalcFieldNum(inputDescription.Fields, field)
	rulesMap := make(map[string]rule.Rule)
	ruleNewFuncs := []func(string, string) rule.Rule{
		rule.NewEQFF,
		rule.NewNEFF,
	}
	for oField, oFd := range inputDescription.Fields {
		if oFd.Kind == fieldtype.String {
			oFieldNum := description.CalcFieldNum(inputDescription.Fields, oField)
			numSharedValues := calcNumSharedValues(fd, oFd)
			if fieldNum < oFieldNum && numSharedValues >= 2 &&
				internal.StringInSlice(oField, ruleFields) {
				for _, ruleNewFunc := range ruleNewFuncs {
					r := ruleNewFunc(field, oField)
					rulesMap[r.String()] = r
				}
			}
		}
	}
	rules := rulesMapToArray(rulesMap)
	return rules
}

func calcNumSharedValues(
	fd1 *description.Field,
	fd2 *description.Field,
) int {
	numShared := 0
	for _, vd1 := range fd1.Values {
		if _, ok := fd2.Values[vd1.Value.String()]; ok {
			numShared++
		}
	}
	return numShared
}

func isNumberField(fd *description.Field) bool {
	return fd.Kind == fieldtype.Number
}

var compareExpr *dexpr.Expr = dexpr.MustNew(
	"min1 < max2 && max1 > min2",
	dexprfuncs.CallFuncs,
)

func hasComparableNumberRange(
	fd1 *description.Field,
	fd2 *description.Field,
) bool {
	if !isNumberField(fd1) || !isNumberField(fd2) {
		return false
	}
	var isComparable bool
	vars := map[string]*dlit.Literal{
		"min1": fd1.Min,
		"max1": fd1.Max,
		"min2": fd2.Min,
		"max2": fd2.Max,
	}
	isComparable, err := compareExpr.EvalBool(vars)
	return err == nil && isComparable
}

func rulesMapToArray(rulesMap map[string]rule.Rule) []rule.Rule {
	rules := make([]rule.Rule, len(rulesMap))
	i := 0
	for _, expr := range rulesMap {
		rules[i] = expr
		i++
	}
	return rules
}
