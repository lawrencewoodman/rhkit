/*
 *   Integration test for package
 */
package rulehunter

import (
	"github.com/lawrencewoodman/rulehunter"
	"path/filepath"
	"runtime"
	"testing"
)

func TestAll(t *testing.T) {
	var experiment *rulehunter.Experiment
	var err error
	experimentFilename := filepath.Join("fixtures", "bank.json")
	experiment, err = rulehunter.LoadExperiment(experimentFilename)
	if err != nil {
		t.Errorf("rulehunter.LoadExperiment(%s) - err: %s",
			experimentFilename, err)
		return
	}

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
	_, jsonErr := assessment5.ToJSON()
	if jsonErr != nil {
		t.Errorf("assessment5.ToJSON() - err: %s", jsonErr)
		return
	}
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
