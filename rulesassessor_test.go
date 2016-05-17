package rulehunter

import (
	"errors"
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/experiment"
	"github.com/vlifesystems/rulehunter/internal"
	"reflect"
	"testing"
)

func TestAssessRules(t *testing.T) {
	fields := []string{"income", "band", "cost"}
	rules := []*Rule{
		mustNewRule("band > 4"),
		mustNewRule("band > 3"),
		mustNewRule("cost > 1.2"),
	}
	aggregators := []*experiment.AggregatorDesc{
		&experiment.AggregatorDesc{"numIncomeGt2", "count", "income > 2"},
		&experiment.AggregatorDesc{"numBandGt4", "count", "band > 4"},
	}
	goals := []string{
		"numIncomeGt2 == 1",
		"numIncomeGt2 == 2",
		"numIncomeGt2 == 3",
		"numIncomeGt2 == 4",
		"numBandGt4 == 1",
		"numBandGt4 == 2",
		"numBandGt4 == 3",
		"numBandGt4 == 4",
	}
	records := []map[string]*dlit.Literal{
		map[string]*dlit.Literal{
			"income": dlit.MustNew(3),
			"cost":   dlit.MustNew(4.5),
			"band":   dlit.MustNew(4),
		},
		map[string]*dlit.Literal{
			"income": dlit.MustNew(3),
			"cost":   dlit.MustNew(3.2),
			"band":   dlit.MustNew(7),
		},
		map[string]*dlit.Literal{
			"income": dlit.MustNew(2),
			"cost":   dlit.MustNew(1.2),
			"band":   dlit.MustNew(4),
		},
		map[string]*dlit.Literal{
			"income": dlit.MustNew(0),
			"cost":   dlit.MustNew(0),
			"band":   dlit.MustNew(9),
		},
	}
	input := NewLiteralInput(records)
	experimentDesc := &experiment.ExperimentDesc{
		Title:         "",
		Input:         input,
		Fields:        fields,
		ExcludeFields: []string{},
		Aggregators:   aggregators,
		Goals:         goals,
		SortOrder:     []*experiment.SortDesc{},
	}
	experiment := mustNewExperiment(experimentDesc)
	wantAssessment := Assessment{
		NumRecords: int64(len(records)),
		Flags: map[string]bool{
			"sorted":  false,
			"refined": false,
		},
		RuleAssessments: []*RuleAssessment{
			&RuleAssessment{
				Rule: mustNewRule("band > 4"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"numBandGt4":     dlit.MustNew("2"),
					"numGoalsPassed": dlit.MustNew(1.001),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
					&GoalAssessment{"numIncomeGt2 == 3", false},
					&GoalAssessment{"numIncomeGt2 == 4", false},
					&GoalAssessment{"numBandGt4 == 1", false},
					&GoalAssessment{"numBandGt4 == 2", true},
					&GoalAssessment{"numBandGt4 == 3", false},
					&GoalAssessment{"numBandGt4 == 4", false},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("band > 3"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"numBandGt4":     dlit.MustNew("2"),
					"numGoalsPassed": dlit.MustNew(0.002),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
					&GoalAssessment{"numIncomeGt2 == 3", false},
					&GoalAssessment{"numIncomeGt2 == 4", false},
					&GoalAssessment{"numBandGt4 == 1", false},
					&GoalAssessment{"numBandGt4 == 2", true},
					&GoalAssessment{"numBandGt4 == 3", false},
					&GoalAssessment{"numBandGt4 == 4", false},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("cost > 1.2"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"numBandGt4":     dlit.MustNew("1"),
					"numGoalsPassed": dlit.MustNew(0.002),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
					&GoalAssessment{"numIncomeGt2 == 3", false},
					&GoalAssessment{"numIncomeGt2 == 4", false},
					&GoalAssessment{"numBandGt4 == 1", true},
					&GoalAssessment{"numBandGt4 == 2", false},
					&GoalAssessment{"numBandGt4 == 3", false},
					&GoalAssessment{"numBandGt4 == 4", false},
				},
			},
		},
	}
	gotAssessment, err := AssessRules(rules, experiment)
	if err != nil {
		t.Errorf("AssessRules(%q, %q, %q, input) - err: %q",
			rules, aggregators, goals, err)
	}

	assessmentsEqual, msg := matchAssessments(gotAssessment, &wantAssessment)
	if !assessmentsEqual {
		t.Errorf("AssessRules(%q, %q, %q, input)\nassessments don't match: %s\n",
			rules, aggregators, goals, msg)
	}
}

func TestAssessRules_errors(t *testing.T) {
	fields := []string{"income", "cost"}
	cases := []struct {
		rules       []*Rule
		aggregators []*experiment.AggregatorDesc
		goals       []string
		wantErr     error
	}{
		{[]*Rule{mustNewRule("band ^^ 3")},
			[]*experiment.AggregatorDesc{
				&experiment.AggregatorDesc{"numIncomeGt2", "count", "income > 2"},
			},
			[]string{"numIncomeGt2 == 1"},
			errors.New("Invalid operator: \"^\"")},
		{[]*Rule{mustNewRule("hand > 3")},
			[]*experiment.AggregatorDesc{
				&experiment.AggregatorDesc{"numIncomeGt2", "count", "income > 2"},
			},
			[]string{"numIncomeGt2 == 1"},
			errors.New("Variable doesn't exist: hand")},
		{[]*Rule{mustNewRule("band > 3")},
			[]*experiment.AggregatorDesc{
				&experiment.AggregatorDesc{"numIncomeGt2", "count", "bincome > 2"},
			},
			[]string{"numIncomeGt2 == 1"},
			errors.New("Variable doesn't exist: bincome")},
		{[]*Rule{mustNewRule("band > 3")},
			[]*experiment.AggregatorDesc{
				&experiment.AggregatorDesc{"numIncomeGt2", "count", "income > 2"},
			},
			[]string{"numIncomeGt == 1"},
			errors.New("Variable doesn't exist: numIncomeGt")},
	}
	records := []map[string]*dlit.Literal{
		map[string]*dlit.Literal{
			"income": dlit.MustNew(3),
			"cost":   dlit.MustNew(4.5),
			"band":   dlit.MustNew(4),
		},
	}
	input := NewLiteralInput(records)
	for _, c := range cases {
		experimentDesc := &experiment.ExperimentDesc{
			Title:         "",
			Input:         input,
			Fields:        fields,
			ExcludeFields: []string{},
			Aggregators:   c.aggregators,
			Goals:         c.goals,
			SortOrder:     []*experiment.SortDesc{},
		}
		experiment := mustNewExperiment(experimentDesc)
		_, err := AssessRules(c.rules, experiment)
		if err == nil || err.Error() != c.wantErr.Error() {
			t.Errorf("AssessRules(%q, %q) - err: %s, wantErr: %s",
				c.rules, experiment, err, c.wantErr)
		}
	}
}

func TestAssessRulesMP(t *testing.T) {
	fields := []string{"income", "band"}
	aggregators := []*experiment.AggregatorDesc{
		&experiment.AggregatorDesc{"numIncomeGt2", "count", "income > 2"},
		&experiment.AggregatorDesc{"numBandGt4", "count", "band > 4"},
	}
	goals := []string{
		"numIncomeGt2 == 1",
		"numIncomeGt2 == 2",
		"numIncomeGt2 == 3",
		"numIncomeGt2 == 4",
		"numBandGt4 == 1",
		"numBandGt4 == 2",
		"numBandGt4 == 3",
		"numBandGt4 == 4",
	}
	records := []map[string]*dlit.Literal{
		map[string]*dlit.Literal{
			"income": dlit.MustNew(3),
			"cost":   dlit.MustNew(4.5),
			"band":   dlit.MustNew(4),
		},
		map[string]*dlit.Literal{
			"income": dlit.MustNew(3),
			"cost":   dlit.MustNew(3.2),
			"band":   dlit.MustNew(7),
		},
		map[string]*dlit.Literal{
			"income": dlit.MustNew(2),
			"cost":   dlit.MustNew(1.2),
			"band":   dlit.MustNew(4),
		},
		map[string]*dlit.Literal{
			"income": dlit.MustNew(0),
			"cost":   dlit.MustNew(0),
			"band":   dlit.MustNew(9),
		},
	}
	cases := []struct {
		rules []*Rule
	}{
		{[]*Rule{
			mustNewRule("band > 4"),
			mustNewRule("band > 3"),
			mustNewRule("cost > 1.2"),
		}},
		{[]*Rule{
			mustNewRule("band > 4"),
			mustNewRule("cost > 1.2"),
		}},
		{[]*Rule{}},
	}

	input := NewLiteralInput(records)
	experimentDesc := &experiment.ExperimentDesc{
		Title:         "",
		Input:         input,
		Fields:        fields,
		ExcludeFields: []string{},
		Aggregators:   aggregators,
		Goals:         goals,
		SortOrder:     []*experiment.SortDesc{},
	}
	experiment := mustNewExperiment(experimentDesc)
	maxProcesses := 4
	for _, cs := range cases {
		wantAssessment, err :=
			AssessRules(cs.rules, experiment)
		if err != nil {
			t.Errorf("AssessRules(%q, %q) - err: %q",
				cs.rules, experiment, err)
		}
		c := make(chan *AssessRulesMPOutcome)
		progress := 0.0
		var gotAssessment *Assessment
		go AssessRulesMP(cs.rules, experiment, maxProcesses, c)

		numRuns := 0
		lastProgress := -1.0
		for o := range c {
			numRuns++
			progress = o.Progress
			if o.Err != nil {
				t.Errorf("AssessRulesMP(%q, %q, ...) - err: %q",
					cs.rules, experiment, o.Err)
			}
			if progress <= lastProgress {
				t.Errorf("AssessRulesMP(%q, %q, ...) - progress not increasing in order: this: %f, last: %f",
					cs.rules, experiment, progress, lastProgress)
			}
			if o.Finished {
				gotAssessment = o.Assessment
			}
		}
		if progress != 1.0 {
			t.Errorf("AssessRulesMP(%q, %q, ...) - progress didn't finish at 100, but: %d",
				cs.rules, experiment, progress)
		}
		if numRuns < len(cs.rules) {
			t.Errorf("AssessRulesMP(%q, %q, ...) - only made %d runs",
				cs.rules, experiment, numRuns)
		}
		assessmentsEqual, msg := matchAssessments(gotAssessment, wantAssessment)
		if !assessmentsEqual {
			t.Errorf("AssessRulesMP(%q, %q, ...)\nassessments don't match: %s",
				cs.rules, experiment, msg)
		}
	}
}

func TestSort(t *testing.T) {
	assessment := Assessment{
		NumRecords: 8,
		Flags: map[string]bool{
			"sorted": false,
		},
		RuleAssessments: []*RuleAssessment{
			&RuleAssessment{
				Rule: mustNewRule("band > 9"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
					"numGoalsPassed": dlit.MustNew(0.003),
					"numIncomeGt2":   dlit.MustNew("3"),
					"numBandGt4":     dlit.MustNew("2"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
					&GoalAssessment{"numIncomeGt2 == 3", false},
					&GoalAssessment{"numIncomeGt2 == 4", false},
					&GoalAssessment{"numBandGt4 == 1", false},
					&GoalAssessment{"numBandGt4 == 2", true},
					&GoalAssessment{"numBandGt4 == 3", false},
					&GoalAssessment{"numBandGt4 == 4", true},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("band > 456"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
					"numGoalsPassed": dlit.MustNew(1.001),
					"numIncomeGt2":   dlit.MustNew("1"),
					"numBandGt4":     dlit.MustNew("2"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
					&GoalAssessment{"numIncomeGt2 == 3", false},
					&GoalAssessment{"numIncomeGt2 == 4", false},
					&GoalAssessment{"numBandGt4 == 1", false},
					&GoalAssessment{"numBandGt4 == 2", true},
					&GoalAssessment{"numBandGt4 == 3", false},
					&GoalAssessment{"numBandGt4 == 4", false},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("band > 3"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("76.3"),
					"numGoalsPassed": dlit.MustNew(0.002),
					"numIncomeGt2":   dlit.MustNew("2"),
					"numBandGt4":     dlit.MustNew("2"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
					&GoalAssessment{"numIncomeGt2 == 3", false},
					&GoalAssessment{"numIncomeGt2 == 4", false},
					&GoalAssessment{"numBandGt4 == 1", false},
					&GoalAssessment{"numBandGt4 == 2", true},
					&GoalAssessment{"numBandGt4 == 3", false},
					&GoalAssessment{"numBandGt4 == 4", false},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("cost > 1.2"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
					"numGoalsPassed": dlit.MustNew(0.002),
					"numIncomeGt2":   dlit.MustNew("2"),
					"numBandGt4":     dlit.MustNew("1"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
					&GoalAssessment{"numIncomeGt2 == 3", false},
					&GoalAssessment{"numIncomeGt2 == 4", false},
					&GoalAssessment{"numBandGt4 == 1", true},
					&GoalAssessment{"numBandGt4 == 2", false},
					&GoalAssessment{"numBandGt4 == 3", false},
					&GoalAssessment{"numBandGt4 == 4", false},
				},
			},
		},
	}
	cases := []struct {
		sortOrder []experiment.SortField
		wantRules []string
	}{
		{[]experiment.SortField{
			experiment.SortField{"numGoalsPassed", experiment.ASCENDING},
		},
			[]string{"band > 456", "band > 9", "band > 3", "cost > 1.2"}},
		{[]experiment.SortField{
			experiment.SortField{"percentMatches", experiment.DESCENDING},
		},
			[]string{"band > 3", "band > 9", "cost > 1.2", "band > 456"}},
		{[]experiment.SortField{
			experiment.SortField{"percentMatches", experiment.ASCENDING},
		},
			[]string{"cost > 1.2", "band > 456", "band > 9", "band > 3"}},
		{[]experiment.SortField{
			experiment.SortField{"percentMatches", experiment.ASCENDING},
			experiment.SortField{"numIncomeGt2", experiment.ASCENDING},
		},
			[]string{"band > 456", "cost > 1.2", "band > 9", "band > 3"}},
		{[]experiment.SortField{
			experiment.SortField{"percentMatches", experiment.DESCENDING},
			experiment.SortField{"numIncomeGt2", experiment.ASCENDING},
		},
			[]string{"band > 9", "band > 3", "cost > 1.2", "band > 456"}},
	}
	for _, c := range cases {
		assessment.Sort(c.sortOrder)
		if !assessment.Flags["sorted"] {
			t.Errorf("Sort(%s) 'sorted' flag not set", c.sortOrder)
		}
		gotRules := getAssessmentRules(&assessment)
		rulesMatch, msg := matchRules(gotRules, c.wantRules)
		if !rulesMatch {
			t.Errorf("matchRules() rules don't match: %s\ngot: %s\nwant: %s\n",
				msg, gotRules, c.wantRules)
		}
	}
}

func TestGetRules(t *testing.T) {
	var gotRules []*Rule
	assessment := Assessment{NumRecords: 8,
		RuleAssessments: []*RuleAssessment{
			&RuleAssessment{
				Rule: mustNewRule("band > 9"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", true},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("band > 456"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("band > 3"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("76.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", true},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("cost > 1.2"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
		},
	}
	cases := []struct {
		numRules     int
		passNumRules bool
		wantRules    []*Rule
	}{
		{0, true, []*Rule{}},
		{1, true, []*Rule{mustNewRule("band > 9")}},
		{2, true, []*Rule{mustNewRule("band > 9"), mustNewRule("band > 456")}},
		{4, true, []*Rule{
			mustNewRule("band > 9"),
			mustNewRule("band > 456"),
			mustNewRule("band > 3"),
			mustNewRule("cost > 1.2"),
		},
		},
		{5, true, []*Rule{
			mustNewRule("band > 9"),
			mustNewRule("band > 456"),
			mustNewRule("band > 3"),
			mustNewRule("cost > 1.2"),
		},
		},
		{0, false, []*Rule{
			mustNewRule("band > 9"),
			mustNewRule("band > 456"),
			mustNewRule("band > 3"),
			mustNewRule("cost > 1.2"),
		},
		},
	}
	for _, c := range cases {
		if c.passNumRules {
			gotRules = assessment.GetRules(c.numRules)
		} else {
			gotRules = assessment.GetRules()
		}
		if !reflect.DeepEqual(gotRules, c.wantRules) {
			t.Errorf("GetRules() passNumRules: %t, numRules: %d rules don't match\ngot: %s\nwant: %s\n",
				c.passNumRules, c.numRules, gotRules, c.wantRules)
		}
	}
}

func TestMerge(t *testing.T) {
	assessment1 := &Assessment{
		NumRecords: 8,
		Flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			&RuleAssessment{
				Rule: mustNewRule("band > 9"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", true},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("band > 456"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("band > 3"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("76.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", true},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("cost > 1.2"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
		},
	}
	assessment2 := &Assessment{
		NumRecords: 8,
		Flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			&RuleAssessment{
				Rule: mustNewRule("band > 16"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("8"),
					"percentMatches": dlit.MustNew("5.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", true},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("team == \"Pi\""),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("19"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("band > 36"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("6.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("cost > 1.27"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("3.5"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
		},
	}

	wantAssessment := &Assessment{
		NumRecords: 8,
		Flags: map[string]bool{
			"sorted":  false,
			"refined": false,
		},
		RuleAssessments: []*RuleAssessment{
			&RuleAssessment{
				Rule: mustNewRule("band > 9"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", true},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("band > 456"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("band > 3"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("76.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", true},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("cost > 1.2"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("band > 16"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("8"),
					"percentMatches": dlit.MustNew("5.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", true},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("team == \"Pi\""),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("19"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("band > 36"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("6.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("cost > 1.27"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("3.5"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
		},
	}
	gotAssessment, err := assessment1.Merge(assessment2)
	if err != nil {
		t.Errorf("Merge() error: %s", err)
		return
	}

	assessmentsEqual, msg := matchAssessments(gotAssessment, wantAssessment)
	if !assessmentsEqual {
		t.Errorf("Merge() assessments don't match: %s\n", msg)
	}
}

func TestMerge_errors(t *testing.T) {
	assessment1 := &Assessment{
		NumRecords: 8,
		Flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			&RuleAssessment{
				Rule: mustNewRule("band > 9"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", true},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("band > 456"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
		},
	}
	assessment2 := &Assessment{NumRecords: 2,
		RuleAssessments: []*RuleAssessment{
			&RuleAssessment{
				Rule: mustNewRule("band > 16"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("8"),
					"percentMatches": dlit.MustNew("5.3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", true},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("team == \"Pi\""),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("19"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numMatches > 3 ", false},
				},
			},
		},
	}
	wantError :=
		errors.New("Can't merge assessments: Number of records don't match")
	_, err := assessment1.Merge(assessment2)
	if err == nil {
		t.Errorf("Merge() not error, expected: %s", wantError)
		return
	}
	if err.Error() != wantError.Error() {
		t.Errorf("Merge() got error: %s, want: %s", err, wantError)
	}
}

func TestRefine(t *testing.T) {
	sortedAssessment := &Assessment{
		NumRecords: 20,
		Flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			&RuleAssessment{
				Rule: mustNewRule("band > 4"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("5"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", true},
					&GoalAssessment{"numIncomeGt2 == 2", false},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("in(band,\"4\",\"3\",\"2\")"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
					"numIncomeGt2":   dlit.MustNew("4"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("in(team,\"a\",\"b\")"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
					"numIncomeGt2":   dlit.MustNew("4"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("in(band,\"99\",\"23\")"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
					"numIncomeGt2":   dlit.MustNew("4"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("band > 3"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
					"numIncomeGt2":   dlit.MustNew("4"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("band > 9"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
					"numIncomeGt2":   dlit.MustNew("4"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("in(band,\"9\",\"2\")"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
					"numIncomeGt2":   dlit.MustNew("4"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("band == 7"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("0"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("true()"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("cost > 1.2"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("3"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
					&GoalAssessment{"numIncomeGt2 == 2", true},
				},
			},
		},
	}
	wantRules := []string{
		"band > 4",
		"in(band,\"4\",\"3\",\"2\")",
		"in(team,\"a\",\"b\")",
		"in(band,\"99\",\"23\")",
		"band > 3",
		"true()",
	}
	numSimilarRules := 2
	sortedAssessment.Refine(numSimilarRules)
	gotRules := getAssessmentRules(sortedAssessment)
	rulesMatch, msg := matchRules(gotRules, wantRules)
	if !rulesMatch {
		t.Errorf("matchRules() rules don't match: %s\ngot: %s\nwant: %s\n",
			msg, gotRules, wantRules)
	}
}

func TestRefine_panic_1(t *testing.T) {
	testPurpose := "Ensure panics if assessment not sorted"
	unsortedAssessment := &Assessment{
		NumRecords: 20,
		RuleAssessments: []*RuleAssessment{
			&RuleAssessment{
				Rule: mustNewRule("band > 4"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{},
			},
			&RuleAssessment{
				Rule: mustNewRule("true()"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
				},
				Goals: []*GoalAssessment{},
			},
		},
	}
	paniced := false
	wantPanic := "Assessment isn't sorted"
	defer func() {
		if r := recover(); r != nil {
			if r.(string) == wantPanic {
				paniced = true
			} else {
				t.Errorf("Test: %s\n", testPurpose)
				t.Errorf("Refine() - got panic: %s, wanted: %s",
					r, wantPanic)
			}
		}
	}()
	numSimilarRules := 1
	unsortedAssessment.Refine(numSimilarRules)
	if !paniced {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("Refine() - failed to panic with: %s", wantPanic)
	}
}

func TestRefine_panic_2(t *testing.T) {
	testPurpose := "Ensure panics if 'true()' rule missing"
	sortedAssessment := &Assessment{
		NumRecords: 20,
		Flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			&RuleAssessment{
				Rule: mustNewRule("band > 4"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{},
			},
			&RuleAssessment{
				Rule: mustNewRule("team > 7"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
				},
				Goals: []*GoalAssessment{},
			},
		},
	}
	paniced := false
	wantPanic := "No 'true()' rule found"
	defer func() {
		if r := recover(); r != nil {
			if r.(string) == wantPanic {
				paniced = true
			} else {
				t.Errorf("Test: %s\n", testPurpose)
				t.Errorf("Refine() - got panic: %s, wanted: %s",
					r, wantPanic)
			}
		}
	}()
	numSimilarRules := 1
	sortedAssessment.Refine(numSimilarRules)
	if !paniced {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("Refine() - failed to panic with: %s", wantPanic)
	}
}

func TestLimitRuleAssessments(t *testing.T) {
	refinedAssessment := &Assessment{
		NumRecords: 20,
		Flags: map[string]bool{
			"sorted":  true,
			"refined": true,
		},
		RuleAssessments: []*RuleAssessment{
			&RuleAssessment{
				Rule: mustNewRule("band > 4"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches": dlit.MustNew("2"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", true},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("in(band,\"4\",\"3\",\"2\")"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches": dlit.MustNew("4"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("in(team,\"a\",\"b\")"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches": dlit.MustNew("4"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("in(band,\"99\",\"23\")"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches": dlit.MustNew("4"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
				},
			},
			&RuleAssessment{
				Rule: mustNewRule("true()"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches": dlit.MustNew("2"),
				},
				Goals: []*GoalAssessment{
					&GoalAssessment{"numIncomeGt2 == 1", false},
				},
			},
		},
	}
	cases := []struct {
		numRules  int
		wantRules []string
	}{
		{3,
			[]string{
				"band > 4",
				"in(band,\"4\",\"3\",\"2\")",
				"in(team,\"a\",\"b\")",
				"true()",
			},
		},
		{4,
			[]string{
				"band > 4",
				"in(band,\"4\",\"3\",\"2\")",
				"in(team,\"a\",\"b\")",
				"in(band,\"99\",\"23\")",
				"true()",
			},
		},
		{5,
			[]string{
				"band > 4",
				"in(band,\"4\",\"3\",\"2\")",
				"in(team,\"a\",\"b\")",
				"in(band,\"99\",\"23\")",
				"true()",
			},
		},
	}
	for _, c := range cases {
		limitedAssessment := refinedAssessment.LimitRuleAssessments(c.numRules)
		gotRules := getAssessmentRules(limitedAssessment)
		rulesMatch, msg := matchRules(gotRules, c.wantRules)
		if !rulesMatch {
			t.Errorf("matchRules() rules don't match: %s\nnumRules: %d\ngot: %s\nwant: %s\n",
				msg, c.numRules, gotRules, c.wantRules)
		}
	}
}

func TestLimitRuleAssessment_panic_1(t *testing.T) {
	testPurpose := "Ensure panics if assessment not sorted"
	unsortedAssessment := &Assessment{
		NumRecords: 20,
		RuleAssessments: []*RuleAssessment{
			&RuleAssessment{
				Rule: mustNewRule("band > 4"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{},
			},
			&RuleAssessment{
				Rule: mustNewRule("true()"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
				},
				Goals: []*GoalAssessment{},
			},
		},
	}
	paniced := false
	wantPanic := "Assessment isn't sorted"
	defer func() {
		if r := recover(); r != nil {
			if r.(string) == wantPanic {
				paniced = true
			} else {
				t.Errorf("Test: %s\n", testPurpose)
				t.Errorf("LimitRuleAssessments() - got panic: %s, wanted: %s",
					r, wantPanic)
			}
		}
	}()
	numRules := 1
	unsortedAssessment.LimitRuleAssessments(numRules)
	if !paniced {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("Refine() - failed to panic with: %s", wantPanic)
	}
}

func TestLimitRuleAssessment_panic_2(t *testing.T) {
	testPurpose := "Ensure panics if assessment not refined"
	unsortedAssessment := &Assessment{
		NumRecords: 20,
		Flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			&RuleAssessment{
				Rule: mustNewRule("band > 4"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{},
			},
			&RuleAssessment{
				Rule: mustNewRule("true()"),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
				},
				Goals: []*GoalAssessment{},
			},
		},
	}
	paniced := false
	wantPanic := "Assessment isn't refined"
	defer func() {
		if r := recover(); r != nil {
			if r.(string) == wantPanic {
				paniced = true
			} else {
				t.Errorf("Test: %s\n", testPurpose)
				t.Errorf("LimitRuleAssessments() - got panic: %s, wanted: %s",
					r, wantPanic)
			}
		}
	}()
	numRules := 1
	unsortedAssessment.LimitRuleAssessments(numRules)
	if !paniced {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("Refine() - failed to panic with: %s", wantPanic)
	}
}

/******************************
 *  Helper functions
 ******************************/

func makeGoalAssessment(expr string, passed bool) *internal.Goal {
	g, err := internal.NewGoal(expr)
	if err != nil {
		panic(fmt.Sprintf("Couldn't create goal: %s", expr))
	}
	g.SetPassed(passed)
	return g
}

func getAssessmentRules(assessment *Assessment) []string {
	rules := make([]string, len(assessment.RuleAssessments))
	for i, ruleAssessment := range assessment.RuleAssessments {
		rules[i] = ruleAssessment.Rule.String()
	}
	return rules
}

func matchAssessments(assessment1, assessment2 *Assessment) (bool, string) {
	if assessment1.NumRecords != assessment2.NumRecords {
		return false, "Number of records don't match"
	}
	if !reflect.DeepEqual(assessment1.Flags, assessment2.Flags) {
		return false,
			fmt.Sprintf("Flags don't match: %s, %s",
				assessment1.Flags,
				assessment2.Flags)
	}
	if len(assessment1.RuleAssessments) != len(assessment2.RuleAssessments) {
		return false, "Number of rule assessments don't match"
	}

	for _, ruleAssessment1 := range assessment1.RuleAssessments {
		ruleFound := false
		for _, ruleAssessment2 := range assessment2.RuleAssessments {
			if ruleAssessment1.Rule.String() == ruleAssessment2.Rule.String() {
				ruleFound = true
				if !ruleAssessment1.isEqual(ruleAssessment2) {
					return false,
						fmt.Sprintf("RuleAssessments don't match:\n %s\n %s",
							ruleAssessment1,
							ruleAssessment2)
				}
			}
		}
		if !ruleFound {
			return false, fmt.Sprintf("Rule doesn't exist: %s", ruleAssessment1.Rule)
		}
	}
	return true, ""
}

func mustNewExperiment(ed *experiment.ExperimentDesc) *experiment.Experiment {
	e, err := experiment.New(ed)
	if err != nil {
		panic(fmt.Sprintf("Can't create Experiment: %s", err))
	}
	return e
}
