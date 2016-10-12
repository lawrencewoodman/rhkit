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

package integrate

import (
	"fmt"
	"github.com/vlifesystems/rhkit"
	"github.com/vlifesystems/rhkit/experiment"
)

func ProcessDataset(experiment *experiment.Experiment) error {
	fieldDescriptions, err := rhkit.DescribeDataset(experiment.Dataset)
	if err != nil {
		return fmt.Errorf("describer.DescribeDataset(experiment.dataset) - err: %s",
			err)
	}
	rules, err :=
		rhkit.GenerateRules(fieldDescriptions, experiment.ExcludeFieldNames)
	if err != nil {
		return fmt.Errorf("rhkit.GenerateRules(%q, %q) - err: %s",
			fieldDescriptions, experiment.ExcludeFieldNames, err)
	}
	if len(rules) < 2 {
		return fmt.Errorf(
			"rhkit.GenerateRules(%q, %q) - not enough rules generated",
			fieldDescriptions,
			experiment.ExcludeFieldNames,
		)
	}

	assessment, err := rhkit.AssessRules(rules, experiment)
	if err != nil {
		return fmt.Errorf("rhkit.AssessRules(rules, %v) - err: %s",
			experiment, err)
	}

	assessment.Sort(experiment.SortOrder)
	assessment.Refine(3)
	sortedRules := assessment.GetRules()

	tweakableRules := rhkit.TweakRules(sortedRules, fieldDescriptions)
	if len(tweakableRules) < 2 {
		return fmt.Errorf("rhkit.TweakRules(sortedRules, %v) - not enough rules generated",

			fieldDescriptions)
	}

	assessment2, err := rhkit.AssessRules(tweakableRules, experiment)
	if err != nil {
		return fmt.Errorf("rhkit.assessRules(tweakableRules, %v) - err: %s",
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
		rhkit.CombineRules(bestNonCombinedRules)
	if len(combinedRules) < 2 {
		return fmt.Errorf("rhkit.CombineRules(bestNonCombinedRules) - not enough rules generated")
	}

	assessment4, err := rhkit.AssessRules(combinedRules, experiment)
	if err != nil {
		return fmt.Errorf("rhkit.assessRules(combinedRules, %v) - err: %s",
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
