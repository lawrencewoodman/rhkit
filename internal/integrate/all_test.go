/*
 *   Integration test for package
 */
package rulehunter

import (
	"fmt"
	"github.com/lawrencewoodman/rulehunter"
	"github.com/lawrencewoodman/rulehunter/csvinput"
	"github.com/lawrencewoodman/rulehunter/input"
	"github.com/lawrencewoodman/rulehunter/reduceinput"
	"path/filepath"
	"runtime"
	"testing"
)

func TestAll_full(t *testing.T) {
	fieldNames := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y"}
	input, err := csvinput.New(
		fieldNames,
		filepath.Join("..", "..", "fixtures", "bank.csv"),
		rune(';'),
		true,
	)
	if err != nil {
		t.Errorf("csvInput.New() - err: %s", err)
		return
	}
	if err = processInput(input, fieldNames); err != nil {
		t.Errorf("processInput() - err: %s", err)
		return
	}
}

func TestAll_reduced(t *testing.T) {
	fieldNames := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y"}
	numRecords := 5

	input, err := csvinput.New(
		fieldNames,
		filepath.Join("..", "..", "fixtures", "bank.csv"),
		rune(';'),
		true,
	)
	if err != nil {
		t.Errorf("csvInput.New() - err: %s", err)
		return
	}

	records, err := reduceinput.New(input, numRecords)
	if err != nil {
		t.Errorf("reduceInput.New() - err: %s", err)
		return
	}

	if err = processInput(records, fieldNames); err != nil {
		t.Errorf("processInput() - err: %s", err)
		return
	}
}

/****************************
 *  Helper functions
 ****************************/
func assessRules(
	rules []*rulehunter.Rule,
	experiment *rulehunter.Experiment,
) (*rulehunter.Assessment, error) {
	var assessment *rulehunter.Assessment
	maxProcesses := runtime.NumCPU()
	c := make(chan *rulehunter.AssessRulesMPOutcome)

	go rulehunter.AssessRulesMP(
		rules,
		experiment.Aggregators,
		experiment.Goals,
		experiment.Input,
		maxProcesses,
		c,
	)
	for o := range c {
		if o.Err != nil {
			return nil, o.Err
		}
		assessment = o.Assessment
	}
	return assessment, nil
}

func processInput(input input.Input, fieldNames []string) error {
	experimentDesc := &rulehunter.ExperimentDesc{
		Title:         "This is a jolly nice title",
		Input:         input,
		Fields:        fieldNames,
		ExcludeFields: []string{"education"},
		Aggregators: []*rulehunter.AggregatorDesc{
			&rulehunter.AggregatorDesc{"numSignedUp", "count", "y == \"yes\""},
			&rulehunter.AggregatorDesc{"cost", "calc", "numMatches * 4.5"},
			&rulehunter.AggregatorDesc{"income", "calc", "numSignedUp * 24"},
			&rulehunter.AggregatorDesc{"profit", "calc", "income - cost"},
			&rulehunter.AggregatorDesc{"oddFigure", "sum", "balance - age"},
			&rulehunter.AggregatorDesc{
				"percentMarried",
				"percent",
				"marital == \"married\"",
			},
		},
		Goals: []string{"profit > 0"},
		SortOrder: []*rulehunter.SortDesc{
			&rulehunter.SortDesc{"profit", "descending"},
			&rulehunter.SortDesc{"numSignedUp", "descending"},
			&rulehunter.SortDesc{"numGoalsPassed", "descending"},
		},
	}
	experiment, err := rulehunter.MakeExperiment(experimentDesc)
	if err != nil {
		return fmt.Errorf("rulehunter.MakeExperiment(%s) - err: %s",
			experimentDesc, err)
	}
	defer experiment.Close()

	fieldDescriptions, err := rulehunter.DescribeInput(experiment.Input)
	if err != nil {
		return fmt.Errorf("rulehunter.DescribeInput(experiment.input) - err: %s",
			err)
	}
	rules, err :=
		rulehunter.GenerateRules(fieldDescriptions, experiment.ExcludeFieldNames)
	if err != nil {
		return fmt.Errorf("rulehunter.GenerateRules(%q, %q) - err: %s",
			fieldDescriptions, experiment.ExcludeFieldNames, err)
	}
	if len(rules) < 2 {
		return fmt.Errorf("rulehunter.GenerateRules(%q, %q) - not enough rules generated",

			fieldDescriptions, experiment.ExcludeFieldNames)
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

	bestNonCombinedRules := assessment3.GetRules()
	numRulesToCombine := 50
	combinedRules :=
		rulehunter.CombineRules(bestNonCombinedRules[:numRulesToCombine])
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
	assessment5.LimitRuleAssessments(finalNumRuleAssessments)
	return nil
}
