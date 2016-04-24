/*
 *   Integration test for package
 */
package rulehunter

import (
	"github.com/lawrencewoodman/rulehunter"
	"github.com/lawrencewoodman/rulehunter/csvinput"
	"path/filepath"
	"runtime"
	"testing"
)

func TestAll(t *testing.T) {
	fieldNames := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y"}
	input, err := csvinput.New(
		fieldNames,
		filepath.Join("fixtures", "bank.csv"),
		rune(';'),
		true,
	)
	if err != nil {
		t.Errorf("rulehunter.NewCsvInput() - err: %s", err)
		return
	}
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
		t.Errorf("rulehunter.MakeExperiment(%s) - err: %s", experimentDesc, err)
		return
	}
	defer experiment.Close()

	fieldDescriptions, err := rulehunter.DescribeInput(experiment.Input)
	if err != nil {
		t.Errorf("rulehunter.DescribeInput(experiment.input) - err: %s", err)
		return
	}
	rules, err :=
		rulehunter.GenerateRules(fieldDescriptions, experiment.ExcludeFieldNames)
	if err != nil {
		t.Errorf("rulehunter.GenerateRules(%q, %q) - err: %s",
			fieldDescriptions, experiment.ExcludeFieldNames, err)
		return
	}
	if len(rules) < 2 {
		t.Errorf("rulehunter.GenerateRules(%q, %q) - not enough rules generated",
			fieldDescriptions, experiment.ExcludeFieldNames)
		return
	}

	assessment, err := assessRules(rules, experiment)
	if err != nil {
		t.Errorf("rulehunter.assessRules(rules, %q) - err: %s",
			experiment, err)
		return
	}

	assessment.Sort(experiment.SortOrder)
	assessment.Refine(3)
	sortedRules := assessment.GetRules()

	tweakableRules := rulehunter.TweakRules(sortedRules, fieldDescriptions)
	if len(tweakableRules) < 2 {
		t.Errorf("rulehunter.TweakRules(sortedRules, %q) - not enough rules generated",
			fieldDescriptions)
		return
	}

	assessment2, err := assessRules(tweakableRules, experiment)
	if err != nil {
		t.Errorf("rulehunter.assessRules(tweakableRules, %q) - err: %s",
			experiment, err)
		return
	}

	assessment3, err := assessment.Merge(assessment2)
	if err != nil {
		t.Errorf("assessment.Merge(assessment2) - err: %s", err)
		return
	}
	assessment3.Sort(experiment.SortOrder)
	assessment3.Refine(1)

	bestNonCombinedRules := assessment3.GetRules()
	numRulesToCombine := 50
	combinedRules :=
		rulehunter.CombineRules(bestNonCombinedRules[:numRulesToCombine])
	if len(combinedRules) < 2 {
		t.Errorf("rulehunter.CombineRules(bestNonCombinedRules) - not enough rules generated")
		return
	}

	assessment4, err := assessRules(combinedRules, experiment)
	if err != nil {
		t.Errorf("rulehunter.assessRules(combinedRules, %q) - err: %s",
			experiment, err)
		return
	}

	assessment5, err := assessment3.Merge(assessment4)
	if err != nil {
		t.Errorf("assessment3.Merge(assessment4) - err: %s", err)
		return
	}
	assessment5.Sort(experiment.SortOrder)
	assessment5.Refine(1)
}

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
