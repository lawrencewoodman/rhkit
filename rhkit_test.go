/*
 *   Integration test for package
 */
package rhkit_test

import (
	"fmt"
	"github.com/lawrencewoodman/ddataset/dcsv"
	"github.com/vlifesystems/rhkit"
	"github.com/vlifesystems/rhkit/experiment"
	"github.com/vlifesystems/rhkit/rule"
	"path/filepath"
	"testing"
)

func TestAll(t *testing.T) {
	fieldNames := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y"}
	experimentDesc := &experiment.ExperimentDesc{
		Title: "This is a jolly nice title",
		Dataset: dcsv.New(
			filepath.Join("fixtures", "bank.csv"),
			true,
			rune(';'),
			fieldNames,
		),
		RuleFields: []string{"age", "job", "marital", "default",
			"balance", "housing", "loan", "contact", "day", "month", "duration",
			"campaign", "pdays", "previous", "poutcome", "y",
		},
		Aggregators: []*experiment.AggregatorDesc{
			&experiment.AggregatorDesc{"numSignedUp", "count", "y == \"yes\""},
			&experiment.AggregatorDesc{"cost", "calc", "numMatches * 4.5"},
			&experiment.AggregatorDesc{"income", "calc", "numSignedUp * 24"},
			&experiment.AggregatorDesc{"profit", "calc", "income - cost"},
			&experiment.AggregatorDesc{"oddFigure", "sum", "balance - age"},
			&experiment.AggregatorDesc{
				"percentMarried",
				"precision",
				"marital == \"married\"",
			},
		},
		Goals: []string{"profit > 0"},
		SortOrder: []*experiment.SortDesc{
			&experiment.SortDesc{"profit", "descending"},
			&experiment.SortDesc{"numSignedUp", "descending"},
			&experiment.SortDesc{"goalsScore", "descending"},
		},
	}
	experiment, err := experiment.New(experimentDesc)
	if err != nil {
		t.Errorf("experiment.New(%s) - err: %s", experimentDesc, err)
		return
	}
	if err = processDataset(experiment); err != nil {
		t.Errorf("processDataset() - err: %s", err)
		return
	}
}

/******************************
 *  Helper functions
 ******************************/
func processDataset(experiment *experiment.Experiment) error {
	var assessment *rhkit.Assessment
	var newAssessment *rhkit.Assessment
	var err error
	fieldDescriptions, err := rhkit.DescribeDataset(experiment.Dataset)
	if err != nil {
		return fmt.Errorf("describer.DescribeDataset(experiment.dataset) - err: %s",
			err)
	}
	complexity := 5
	rules := rule.Generate(
		fieldDescriptions,
		experiment.RuleFieldNames,
		complexity,
	)
	if len(rules) < 2 {
		return fmt.Errorf(
			"rhkit.GenerateRules(%v, %v) - not enough rules generated",
			fieldDescriptions,
			experiment.RuleFieldNames,
		)
	}

	assessment, err = rhkit.AssessRules(rules, experiment)
	if err != nil {
		return fmt.Errorf("rhkit.AssessRules(rules, %v) - err: %s",
			experiment, err)
	}

	assessment.Sort(experiment.SortOrder)
	assessment.Refine()
	rules = assessment.Rules()

	tweakableRules := rule.Tweak(
		1,
		rules,
		fieldDescriptions,
	)
	if len(tweakableRules) < 2 {
		return fmt.Errorf("rhkit.TweakRules(sortedRules, %v) - not enough rules generated",

			fieldDescriptions)
	}

	newAssessment, err = rhkit.AssessRules(tweakableRules, experiment)
	if err != nil {
		return fmt.Errorf("rhkit.assessRules(tweakableRules, %v) - err: %s",
			experiment, err)
	}

	assessment, err = assessment.Merge(newAssessment)
	if err != nil {
		return fmt.Errorf("assessment.Merge(assessment2) - err: %s", err)
	}
	assessment.Sort(experiment.SortOrder)
	assessment.Refine()

	rules = assessment.Rules()
	reducedDPRules := rule.ReduceDP(rules)
	if len(reducedDPRules) < 2 {
		return fmt.Errorf("rule.ReduceDP: not enough rules generated")
	}

	newAssessment, err = rhkit.AssessRules(reducedDPRules, experiment)
	if err != nil {
		return fmt.Errorf("rhkit.assessRules(reducedDPRules, %v) - err: %s",
			experiment, err)
	}

	assessment, err = assessment.Merge(newAssessment)
	if err != nil {
		return fmt.Errorf("assessment.Merge: %s", err)
	}
	assessment.Sort(experiment.SortOrder)
	assessment.Refine()

	numRulesToCombine := 50
	bestNonCombinedRules := assessment.Rules(numRulesToCombine)
	combinedRules :=
		rule.Combine(bestNonCombinedRules)
	if len(combinedRules) < 2 {
		return fmt.Errorf("rhkit.CombineRules(bestNonCombinedRules) - not enough rules generated")
	}

	newAssessment, err = rhkit.AssessRules(combinedRules, experiment)
	if err != nil {
		return fmt.Errorf("rhkit.assessRules(combinedRules, %v) - err: %s",
			experiment, err)
	}

	assessment, err = assessment.Merge(newAssessment)
	if err != nil {
		return fmt.Errorf("assessment3.Merge: %s", err)
	}
	assessment.Sort(experiment.SortOrder)
	assessment.Refine()

	finalNumRuleAssessments := 100
	assessment.TruncateRuleAssessments(finalNumRuleAssessments)
	return nil
}
