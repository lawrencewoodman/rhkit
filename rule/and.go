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

	tRuleA, ruleAIsTweakable := ruleA.(TweakableRule)
	tRuleB, ruleBIsTweakable := ruleB.(TweakableRule)
	if !ruleAIsTweakable || !ruleBIsTweakable {
		_, ruleAIsIn := ruleA.(*InFV)
		_, ruleBIsIn := ruleB.(*InFV)
		ruleAFields := ruleA.GetFields()
		ruleBFields := ruleB.GetFields()
		if (ruleAIsIn &&
			len(ruleBFields) == 1 &&
			ruleAFields[0] == ruleBFields[0]) ||
			(ruleBIsIn &&
				len(ruleAFields) == 1 &&
				ruleAFields[0] == ruleBFields[0]) {
			return nil, fmt.Errorf("can't And rule: %s, with: %s", ruleA, ruleB)
		} else {
			return &And{ruleA: ruleA, ruleB: ruleB}, nil
		}
	}
	r, err := joinTweakableRulesBetween(tRuleA, tRuleB)
	return r, err
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
	}
	switch r.ruleB.(type) {
	case *And:
		bStr = "(" + bStr + ")"
	case *Or:
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

func joinTweakableRulesBetween(
	ruleA TweakableRule,
	ruleB TweakableRule,
) (Rule, error) {
	var r Rule
	var err error
	fieldA := ruleA.GetFields()[0]
	fieldB := ruleB.GetFields()[0]

	if fieldA == fieldB {
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
			return nil, fmt.Errorf("can't And rule: %s, with: %s", ruleA, ruleB)
		}
		return r, err
	}
	return &And{ruleA: ruleA, ruleB: ruleB}, nil
}
