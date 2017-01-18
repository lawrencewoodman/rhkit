/*
	Copyright (C) 2016 vLife Systems Ltd <http://vlifesystems.com>
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

package rule

import (
	"fmt"
	"github.com/lawrencewoodman/ddataset"
	"github.com/lawrencewoodman/dlit"
)

// Or represents a rule determening if ruleA OR ruleB
type Or struct {
	ruleA Rule
	ruleB Rule
}

func NewOr(ruleA Rule, ruleB Rule) (Rule, error) {
	_, ruleAIsTrue := ruleA.(True)
	_, ruleBIsTrue := ruleB.(True)
	if ruleAIsTrue || ruleBIsTrue {
		return nil, fmt.Errorf("can't Or rule: %s, with: %s", ruleA, ruleB)
	}
	inRuleA, ruleAIsIn := ruleA.(*InFV)
	inRuleB, ruleBIsIn := ruleB.(*InFV)
	if ruleAIsIn && ruleBIsIn {
		return handleInRules(inRuleA, inRuleB), nil
	}

	return &Or{ruleA: ruleA, ruleB: ruleB}, nil
}

func handleInRules(ruleA, ruleB *InFV) Rule {
	if ruleA.GetFields()[0] != ruleB.GetFields()[0] {
		return &Or{ruleA: ruleA, ruleB: ruleB}
	}
	newValues := []*dlit.Literal{}
	mNewValues := map[string]interface{}{}
	for _, v := range ruleA.GetValues() {
		if _, ok := mNewValues[v.String()]; !ok {
			mNewValues[v.String()] = nil
			newValues = append(newValues, v)
		}
	}
	for _, v := range ruleB.GetValues() {
		if _, ok := mNewValues[v.String()]; !ok {
			mNewValues[v.String()] = nil
			newValues = append(newValues, v)
		}
	}
	return NewInFV(ruleA.GetFields()[0], newValues)
}

func MustNewOr(ruleA Rule, ruleB Rule) Rule {
	r, err := NewOr(ruleA, ruleB)
	if err != nil {
		panic(err)
	}
	return r
}

func (r *Or) String() string {
	// TODO: Consider making this OR rather than ||
	aStr := r.ruleA.String()
	bStr := r.ruleB.String()
	switch r.ruleA.(type) {
	case *And:
		aStr = "(" + aStr + ")"
	case *Or:
		aStr = "(" + aStr + ")"
	}
	switch r.ruleB.(type) {
	case *And:
		bStr = "(" + bStr + ")"
	case *Or:
		bStr = "(" + bStr + ")"
	}
	return fmt.Sprintf("%s || %s", aStr, bStr)
}

func (r *Or) IsTrue(record ddataset.Record) (bool, error) {
	lh, err := r.ruleA.IsTrue(record)
	if err != nil {
		return false, InvalidRuleError{Rule: r}
	}
	rh, err := r.ruleB.IsTrue(record)
	if err != nil {
		return false, InvalidRuleError{Rule: r}
	}
	return lh || rh, nil
}

func (r *Or) GetFields() []string {
	results := []string{}
	mResults := map[string]interface{}{}
	for _, f := range r.ruleA.GetFields() {
		if _, ok := mResults[f]; !ok {
			mResults[f] = nil
			results = append(results, f)
		}
	}
	for _, f := range r.ruleB.GetFields() {
		if _, ok := mResults[f]; !ok {
			mResults[f] = nil
			results = append(results, f)
		}
	}
	return results
}
