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
)

// And represents a rule determening if ruleA AND ruleB
type And struct {
	ruleA Rule
	ruleB Rule
}

func NewAnd(ruleA Rule, ruleB Rule) (Rule, error) {
	_, ruleAIsTrue := ruleA.(True)
	_, ruleBIsTrue := ruleB.(True)
	if ruleAIsTrue || ruleBIsTrue {
		return nil, fmt.Errorf("can't And rule: %s, with: %s", ruleA, ruleB)
	}

	skip, r := tryJoinTweakableRulesBetween(ruleA, ruleB)
	if !skip {
		if r != nil {
			return r, nil
		}
		return nil, fmt.Errorf("can't And rule: %s, with: %s", ruleA, ruleB)
	}
	skip, r = tryInRule(ruleA, ruleB)
	if !skip {
		if r != nil {
			return r, nil
		}
		return nil, fmt.Errorf("can't And rule: %s, with: %s", ruleA, ruleB)
	}
	skip, r = tryEqNeRule(ruleA, ruleB)
	if !skip {
		if r != nil {
			return r, nil
		}
		return nil, fmt.Errorf("can't And rule: %s, with: %s", ruleA, ruleB)
	}
	return &And{ruleA: ruleA, ruleB: ruleB}, nil
}

func tryInRule(ruleA, ruleB Rule) (skip bool, newRule Rule) {
	_, ruleAIsIn := ruleA.(*InFV)
	_, ruleBIsIn := ruleB.(*InFV)
	if !ruleAIsIn && !ruleBIsIn {
		return true, nil
	}
	fieldsA := ruleA.GetFields()
	fieldsB := ruleB.GetFields()
	if len(fieldsA) == 1 && len(fieldsB) == 1 && fieldsA[0] == fieldsB[0] {
		return false, nil
	}
	return false, &And{ruleA: ruleA, ruleB: ruleB}
}

func tryEqNeRule(ruleA, ruleB Rule) (skip bool, newRule Rule) {
	ruleAFields := ruleA.GetFields()
	ruleBFields := ruleB.GetFields()
	if len(ruleAFields) != 1 && len(ruleBFields) != 1 {
		return true, nil
	}
	fieldA := ruleA.GetFields()[0]
	fieldB := ruleB.GetFields()[0]

	if fieldA != fieldB {
		return true, nil
	}
	switch ruleA.(type) {
	case *EQFVI:
	case *EQFVF:
	case *EQFVS:
	case *NEFVI:
	case *NEFVF:
	case *NEFVS:
	default:
		return true, nil
	}
	switch ruleB.(type) {
	case *EQFVI:
	case *EQFVF:
	case *EQFVS:
	case *NEFVI:
	case *NEFVF:
	case *NEFVS:
	default:
		return false, &And{ruleA: ruleA, ruleB: ruleB}
	}
	return false, nil
}

func tryJoinTweakableRulesBetween(
	ruleA Rule,
	ruleB Rule,
) (skip bool, newRule Rule) {
	_, ruleAIsTweakable := ruleA.(TweakableRule)
	_, ruleBIsTweakable := ruleB.(TweakableRule)
	if !ruleAIsTweakable || !ruleBIsTweakable {
		return true, nil
	}

	var r Rule
	var err error
	fieldA := ruleA.GetFields()[0]
	fieldB := ruleB.GetFields()[0]

	if fieldA == fieldB {
		_, ruleAIsBetweenFVI := ruleA.(*BetweenFVI)
		_, ruleAIsBetweenFVF := ruleA.(*BetweenFVF)
		_, ruleBIsBetweenFVI := ruleB.(*BetweenFVI)
		_, ruleBIsBetweenFVF := ruleB.(*BetweenFVF)
		_, ruleAIsOutsideFVI := ruleA.(*OutsideFVI)
		_, ruleAIsOutsideFVF := ruleA.(*OutsideFVF)
		_, ruleBIsOutsideFVI := ruleB.(*OutsideFVI)
		_, ruleBIsOutsideFVF := ruleB.(*OutsideFVF)
		if (ruleAIsBetweenFVI && !ruleBIsBetweenFVI) ||
			(!ruleAIsBetweenFVI && ruleBIsBetweenFVI) ||
			(ruleAIsBetweenFVF && !ruleBIsBetweenFVF) ||
			(!ruleAIsBetweenFVF && ruleBIsBetweenFVF) ||
			(ruleAIsOutsideFVI && !ruleBIsOutsideFVI) ||
			(!ruleAIsOutsideFVI && ruleBIsOutsideFVI) ||
			(ruleAIsOutsideFVF && !ruleBIsOutsideFVF) ||
			(!ruleAIsOutsideFVF && ruleBIsOutsideFVF) {
			return true, nil
		}

		GEFVIRuleA, ruleAIsGEFVI := ruleA.(*GEFVI)
		LEFVIRuleA, ruleAIsLEFVI := ruleA.(*LEFVI)
		GEFVIRuleB, ruleBIsGEFVI := ruleB.(*GEFVI)
		LEFVIRuleB, ruleBIsLEFVI := ruleB.(*LEFVI)
		GEFVFRuleA, ruleAIsGEFVF := ruleA.(*GEFVF)
		LEFVFRuleA, ruleAIsLEFVF := ruleA.(*LEFVF)
		GEFVFRuleB, ruleBIsGEFVF := ruleB.(*GEFVF)
		LEFVFRuleB, ruleBIsLEFVF := ruleB.(*LEFVF)
		if ruleAIsGEFVI && ruleBIsLEFVI {
			r, err = NewBetweenFVI(
				fieldA,
				GEFVIRuleA.GetValue(),
				LEFVIRuleB.GetValue(),
			)
		} else if ruleAIsLEFVI && ruleBIsGEFVI {
			r, err = NewBetweenFVI(
				fieldA,
				GEFVIRuleB.GetValue(),
				LEFVIRuleA.GetValue(),
			)
		} else if ruleAIsGEFVF && ruleBIsLEFVF {
			r, err = NewBetweenFVF(
				fieldA,
				GEFVFRuleA.GetValue(),
				LEFVFRuleB.GetValue(),
			)
		} else if ruleAIsLEFVF && ruleBIsGEFVF {
			r, err = NewBetweenFVF(
				fieldA,
				GEFVFRuleB.GetValue(),
				LEFVFRuleA.GetValue(),
			)
		} else {
			err = fmt.Errorf("can't And rule: %s, with: %s", ruleA, ruleB)
		}
		if err != nil {
			return false, nil
		}
		return false, r
	}
	return true, nil
}

func MustNewAnd(ruleA Rule, ruleB Rule) Rule {
	r, err := NewAnd(ruleA, ruleB)
	if err != nil {
		panic(err)
	}
	return r
}

func (r *And) String() string {
	// TODO: Consider making this AND rather than &&
	aStr := r.ruleA.String()
	bStr := r.ruleB.String()
	switch r.ruleA.(type) {
	case *And:
		aStr = "(" + aStr + ")"
	case *Or:
		aStr = "(" + aStr + ")"
	case *BetweenFVI:
		aStr = "(" + aStr + ")"
	case *BetweenFVF:
		aStr = "(" + aStr + ")"
	case *OutsideFVI:
		aStr = "(" + aStr + ")"
	case *OutsideFVF:
		aStr = "(" + aStr + ")"
	}
	switch r.ruleB.(type) {
	case *And:
		bStr = "(" + bStr + ")"
	case *Or:
		bStr = "(" + bStr + ")"
	case *BetweenFVI:
		bStr = "(" + bStr + ")"
	case *BetweenFVF:
		bStr = "(" + bStr + ")"
	case *OutsideFVI:
		bStr = "(" + bStr + ")"
	case *OutsideFVF:
		bStr = "(" + bStr + ")"
	}
	return fmt.Sprintf("%s && %s", aStr, bStr)
}

func (r *And) IsTrue(record ddataset.Record) (bool, error) {
	lh, err := r.ruleA.IsTrue(record)
	if err != nil {
		return false, InvalidRuleError{Rule: r}
	}
	rh, err := r.ruleB.IsTrue(record)
	if err != nil {
		return false, InvalidRuleError{Rule: r}
	}
	return lh && rh, nil
}

func (r *And) GetFields() []string {
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
