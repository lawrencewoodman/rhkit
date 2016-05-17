/*
	Copyright (C) 2016 vLife Systems Ltd <http://vlifesystems.com>
	This file is part of Rulehunter.

	Rulehunter is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	Rulehunter is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with Rulehunter; see the file COPYING.  If not, see
	<http://www.gnu.org/licenses/>.
*/

package ruleassessment

import (
	"errors"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/internal"
	"github.com/vlifesystems/rulehunter/rule"
)

type RuleAssessment struct {
	Rule        *rule.Rule
	Aggregators []internal.Aggregator
	Goals       []*internal.Goal
}

func New(
	rule *rule.Rule,
	aggregators []internal.Aggregator,
	goals []*internal.Goal,
) *RuleAssessment {
	// Clone the aggregators and goals to ensure the results are
	// specific to this rule
	cloneAggregators := make([]internal.Aggregator, len(aggregators))
	for i, a := range aggregators {
		cloneAggregators[i] = a.CloneNew()
	}
	cloneGoals := make([]*internal.Goal, len(goals))
	for i, g := range goals {
		cloneGoals[i] = g.Clone()
	}
	return &RuleAssessment{Rule: rule, Aggregators: cloneAggregators,
		Goals: cloneGoals}
}

func (ra *RuleAssessment) NextRecord(record map[string]*dlit.Literal) error {
	var ruleIsTrue bool
	var err error
	for _, aggregator := range ra.Aggregators {
		ruleIsTrue, err = ra.Rule.IsTrue(record)
		if err != nil {
			return err
		}
		err = aggregator.NextRecord(record, ruleIsTrue)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ra *RuleAssessment) GetAggregatorValue(
	name string,
	numRecords int64,
) (*dlit.Literal, bool) {
	for _, aggregator := range ra.Aggregators {
		if aggregator.GetName() == name {
			return aggregator.GetResult(ra.Aggregators, ra.Goals, numRecords), true
		}
	}
	// TODO: Test and create specific error type
	err := errors.New("Aggregator doesn't exist")
	l := dlit.MustNew(err)
	return l, false
}
