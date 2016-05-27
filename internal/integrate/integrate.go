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

package integrate

import (
	"fmt"
	"github.com/vlifesystems/rulehunter"
	"github.com/vlifesystems/rulehunter/assessment"
	"github.com/vlifesystems/rulehunter/experiment"
	"github.com/vlifesystems/rulehunter/rule"
	"runtime"
)

func ProcessDataset(experiment *experiment.Experiment) error {
	fieldDescriptions, err := rulehunter.DescribeDataset(experiment.Dataset)
	if err != nil {
		return fmt.Errorf("describer.DescribeDataset(experiment.dataset) - err: %s",
			err)
	}
	rules, err :=
		rulehunter.GenerateRules(fieldDescriptions, experiment.ExcludeFieldNames)
	if err != nil {
		return fmt.Errorf("rulehunter.GenerateRules(%q, %q) - err: %s",
			fieldDescriptions, experiment.ExcludeFieldNames, err)
	}
	if len(rules) < 2 {
		return fmt.Errorf(
			"rulehunter.GenerateRules(%q, %q) - not enough rules generated",
			fieldDescriptions,
			experiment.ExcludeFieldNames,
		)
	}

	assessment, err := assessRules(rules, experiment)
	if err != nil {
		return fmt.Errorf("rulehunter.assessRules(rules, %q) - err: %s",
			experiment, err)
	}

	assessment.Sort(experiment.SortOrder)
	assessment.Refine(3)
	sortedRules := assessment.GetRules()

	tweakableRules := rulehunter.TweakRules(sortedRules, fieldDescriptions)
	if len(tweakableRules) < 2 {
		return fmt.Errorf("rulehunter.TweakRules(sortedRules, %q) - not enough rules generated",

			fieldDescriptions)
	}

	assessment2, err := assessRules(tweakableRules, experiment)
	if err != nil {
		return fmt.Errorf("rulehunter.assessRules(tweakableRules, %q) - err: %s",
			experiment, err)
	}

	assessment3, err := assessment.Merge(assessment2)
	if err != nil {
		return fmt.Errorf("assessment.Merge(assessment2) - err: %s", err)
	}
	assessment3.Sort(experiment.SortOrder)
	assessment3.Refine(1)

	numRulesToCombine := 50
	bestNonCombinedRules := assessment3.GetRules(numRulesToCombine)
	combinedRules :=
		rulehunter.CombineRules(bestNonCombinedRules)
	if len(combinedRules) < 2 {
		return fmt.Errorf("rulehunter.CombineRules(bestNonCombinedRules) - not enough rules generated")
	}

	assessment4, err := assessRules(combinedRules, experiment)
	if err != nil {
		return fmt.Errorf("rulehunter.assessRules(combinedRules, %q) - err: %s",
			experiment, err)
	}

	assessment5, err := assessment3.Merge(assessment4)
	if err != nil {
		return fmt.Errorf("assessment3.Merge(assessment4) - err: %s", err)
	}
	assessment5.Sort(experiment.SortOrder)
	assessment5.Refine(1)

	finalNumRuleAssessments := 100
	assessment5.TruncateRuleAssessments(finalNumRuleAssessments)
	return nil
}

func assessRules(
	rules []*rule.Rule,
	experiment *experiment.Experiment,
) (*assessment.Assessment, error) {
	var assessment *assessment.Assessment
	maxProcesses := runtime.NumCPU()
	c := make(chan *rulehunter.AssessRulesMPOutcome)

	go rulehunter.AssessRulesMP(rules, experiment, maxProcesses, c)
	for o := range c {
		if o.Err != nil {
			return nil, o.Err
		}
		assessment = o.Assessment
	}
	return assessment, nil
}
