// Copyright (C) 2016-2017 vLife Systems Ltd <http://vlifesystems.com>
// Licensed under an MIT licence.  Please see LICENSE.md for details.

// package rhkit is used to find rules in a Dataset to satisfy user defined
// goals
package rhkit

import (
	"errors"
	"github.com/lawrencewoodman/ddataset"
	"github.com/vlifesystems/rhkit/aggregator"
	"github.com/vlifesystems/rhkit/assessment"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/goal"
	"github.com/vlifesystems/rhkit/rule"
)

// ErrNoRulesGenerated indicates that no rules were generated
var ErrNoRulesGenerated = errors.New("no rules generated")

// DescribeError indicates an error describing a Dataset
type DescribeError struct {
	Err error
}

func (e DescribeError) Error() string {
	return "problem describing dataset: " + e.Err.Error()
}

// GenerateRulesError indicates an error generating rules
type GenerateRulesError struct {
	Err error
}

func (e GenerateRulesError) Error() string {
	return "problem generating rules: " + e.Err.Error()
}

// AssessError indicates an error assessing rules
type AssessError struct {
	Err error
}

func (e AssessError) Error() string {
	return "problem assessing rules: " + e.Err.Error()
}

// MergeError indicates an error Merging assessments
type MergeError struct {
	Err error
}

func (e MergeError) Error() string {
	return "problem merging assessments: " + e.Err.Error()
}

type Options struct {
	MaxNumRules    int
	GenerateRules  bool
	RuleComplexity rule.Complexity
}

// Process processes a Dataset to find Rules to meet the supplied requirements
func Process(
	dataset ddataset.Dataset,
	ruleFields []string,
	aggregators []aggregator.Spec,
	goals []*goal.Goal,
	sortOrder []assessment.SortOrder,
	rules []rule.Rule,
	opts Options,
) (*assessment.Assessment, error) {
	fieldDescriptions, err := description.DescribeDataset(dataset)
	if err != nil {
		return nil, DescribeError{Err: err}
	}

	if !opts.GenerateRules {
		rules = append(rules, rule.NewTrue())
	}
	ass, err := assessment.AssessRules(
		dataset,
		rules,
		aggregators,
		goals,
	)
	if err != nil {
		return nil, AssessError{Err: err}
	}

	if opts.GenerateRules {
		generatedAss, err := processGenerate(
			dataset,
			ruleFields,
			fieldDescriptions,
			aggregators,
			goals,
			sortOrder,
			len(rules),
			opts,
		)
		if err != nil {
			return nil, err
		}

		ass, err = ass.Merge(generatedAss)
		if err != nil {
			return nil, MergeError{Err: err}
		}
	}

	ass.Sort(sortOrder)
	return ass, nil
}

func processGenerate(
	dataset ddataset.Dataset,
	ruleFields []string,
	fieldDescriptions *description.Description,
	aggregators []aggregator.Spec,
	goals []*goal.Goal,
	sortOrder []assessment.SortOrder,
	numUserRules int,
	opts Options,
) (*assessment.Assessment, error) {
	var ass *assessment.Assessment
	var newAss *assessment.Assessment
	var err error

	generatedRules, err := rule.Generate(
		fieldDescriptions,
		ruleFields,
		opts.RuleComplexity,
	)
	if err != nil {
		return nil, GenerateRulesError{Err: err}
	}
	if len(generatedRules) < 2 {
		return nil, ErrNoRulesGenerated
	}

	ass, err = assessment.AssessRules(
		dataset,
		generatedRules,
		aggregators,
		goals,
	)
	if err != nil {
		return nil, AssessError{Err: err}
	}

	ass.Sort(sortOrder)
	ass.Refine()
	bestRules := ass.Rules()

	tweakableRules := rule.Tweak(1, bestRules, fieldDescriptions)
	newAss, err = assessment.AssessRules(
		dataset,
		tweakableRules,
		aggregators,
		goals,
	)
	if err != nil {
		return nil, AssessError{Err: err}
	}

	ass, err = ass.Merge(newAss)
	if err != nil {
		return nil, MergeError{Err: err}
	}
	ass.Sort(sortOrder)
	ass.Refine()

	bestRules = ass.Rules()
	reducedDPRules := rule.ReduceDP(bestRules)

	newAss, err = assessment.AssessRules(
		dataset,
		reducedDPRules,
		aggregators,
		goals,
	)
	if err != nil {
		return nil, AssessError{Err: err}
	}

	ass, err = ass.Merge(newAss)
	if err != nil {
		return nil, MergeError{Err: err}
	}
	ass.Sort(sortOrder)
	ass.Refine()

	numRulesToCombine := 50
	bestNonCombinedRules := ass.Rules(numRulesToCombine)
	combinedRules := rule.Combine(bestNonCombinedRules)

	newAss, err = assessment.AssessRules(
		dataset,
		combinedRules,
		aggregators,
		goals,
	)
	if err != nil {
		return nil, AssessError{Err: err}
	}

	ass, err = ass.Merge(newAss)
	if err != nil {
		return nil, MergeError{Err: err}
	}
	ass.Sort(sortOrder)
	ass.Refine()

	if opts.MaxNumRules-numUserRules < 1 {
		return ass.TruncateRuleAssessments(1), nil
	}

	return ass.TruncateRuleAssessments(opts.MaxNumRules - numUserRules), nil
}
