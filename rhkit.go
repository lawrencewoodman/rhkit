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

// package rhkit is used to find rules in a Dataset to satisfy user defined
// goals
package rhkit

import (
	"errors"
	"github.com/vlifesystems/rhkit/assessment"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/experiment"
	"github.com/vlifesystems/rhkit/rule"
)

var ErrNoRulesGenerated = errors.New("no rules generated")

type DescribingError struct {
	Err error
}

func (e DescribingError) Error() string {
	return "problem describing dataset: " + e.Err.Error()
}

type AssessError struct {
	Err error
}

func (e AssessError) Error() string {
	return "problem assessing rules: " + e.Err.Error()
}

type MergeError struct {
	Err error
}

func (e MergeError) Error() string {
	return "problem merging rules: " + e.Err.Error()
}

// Process processes an Experiment and returns an assessment
func Process(
	experiment *experiment.Experiment,
	maxNumRules int,
) (*assessment.Assessment, error) {
	var ass *assessment.Assessment
	var newAss *assessment.Assessment
	var err error
	fieldDescriptions, err := description.DescribeDataset(experiment.Dataset)
	if err != nil {
		return nil, DescribingError{Err: err}
	}
	rules := rule.Generate(
		fieldDescriptions,
		experiment.RuleFields,
		experiment.RuleComplexity,
	)
	if len(rules) < 2 {
		return nil, ErrNoRulesGenerated
	}

	ass, err = assessment.AssessRules(rules, experiment)
	if err != nil {
		return nil, AssessError{Err: err}
	}

	ass.Sort(experiment.SortOrder)
	ass.Refine()
	rules = ass.Rules()

	tweakableRules := rule.Tweak(
		1,
		rules,
		fieldDescriptions,
	)

	newAss, err = assessment.AssessRules(tweakableRules, experiment)
	if err != nil {
		return nil, AssessError{Err: err}
	}

	ass, err = ass.Merge(newAss)
	if err != nil {
		return nil, MergeError{Err: err}
	}
	ass.Sort(experiment.SortOrder)
	ass.Refine()

	rules = ass.Rules()
	reducedDPRules := rule.ReduceDP(rules)

	newAss, err = assessment.AssessRules(reducedDPRules, experiment)
	if err != nil {
		return nil, AssessError{Err: err}
	}

	ass, err = ass.Merge(newAss)
	if err != nil {
		return nil, MergeError{Err: err}
	}
	ass.Sort(experiment.SortOrder)
	ass.Refine()

	numRulesToCombine := 50
	bestNonCombinedRules := ass.Rules(numRulesToCombine)
	combinedRules := rule.Combine(bestNonCombinedRules)

	newAss, err = assessment.AssessRules(combinedRules, experiment)
	if err != nil {
		return nil, AssessError{Err: err}
	}

	ass, err = ass.Merge(newAss)
	if err != nil {
		return nil, MergeError{Err: err}
	}
	ass.Sort(experiment.SortOrder)
	ass.Refine()

	return ass.TruncateRuleAssessments(maxNumRules), nil
}
