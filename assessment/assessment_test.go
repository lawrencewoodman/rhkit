package assessment

import (
	"path/filepath"
	"reflect"
	"sync"
	"testing"

	"github.com/lawrencewoodman/ddataset/dcsv"
	"github.com/lawrencewoodman/ddataset/dtruncate"
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/aggregator"
	"github.com/vlifesystems/rhkit/goal"
	"github.com/vlifesystems/rhkit/internal/testhelpers"
	"github.com/vlifesystems/rhkit/rule"
)

func TestAssessRules(t *testing.T) {
	rules := []rule.Rule{
		rule.NewGEFV("band", dlit.MustNew(5)),
		rule.NewGEFV("band", dlit.MustNew(4)),
		rule.NewGEFV("cost", dlit.MustNew(1.3)),
	}
	aggregatorDescs := []*aggregator.Desc{
		{"numIncomeGt2", "count", "income > 2"},
		{"numBandGt4", "count", "band > 4"},
	}
	goalExprs := []string{
		"numIncomeGt2 == 1",
		"numIncomeGt2 == 2",
		"numIncomeGt2 == 3",
		"numIncomeGt2 == 4",
		"numBandGt4 == 1",
		"numBandGt4 == 2",
		"numBandGt4 == 3",
		"numBandGt4 == 4",
	}
	fields := []string{"income", "cost", "band"}
	records := [][]string{
		{"3", "4.5", "4"},
		{"3", "3.2", "7"},
		{"2", "1.2", "4"},
		{"0", "0", "9"},
	}
	dataset := testhelpers.NewLiteralDataset(fields, records)
	aggregatorSpecs, err := aggregator.MakeSpecs(fields, aggregatorDescs)
	if err != nil {
		t.Fatalf("MakeSpecs: %s", err)
	}
	goals, err := goal.MakeGoals(goalExprs)
	if err != nil {
		t.Fatalf("MakeGoals: %s", err)
	}
	wantIsSorted := false
	wantIsRefined := false
	wantNumRecords := int64(len(records))
	wantRuleAssessments := []*RuleAssessment{
		{
			Rule: rule.NewGEFV("band", dlit.MustNew(5)),
			Aggregators: map[string]*dlit.Literal{
				"numMatches":     dlit.MustNew("2"),
				"percentMatches": dlit.MustNew("50"),
				"numIncomeGt2":   dlit.MustNew("1"),
				"numBandGt4":     dlit.MustNew("2"),
				"goalsScore":     dlit.MustNew(1.001),
			},
			Goals: []*GoalAssessment{
				{"numIncomeGt2 == 1", true},
				{"numIncomeGt2 == 2", false},
				{"numIncomeGt2 == 3", false},
				{"numIncomeGt2 == 4", false},
				{"numBandGt4 == 1", false},
				{"numBandGt4 == 2", true},
				{"numBandGt4 == 3", false},
				{"numBandGt4 == 4", false},
			},
		},
		{
			Rule: rule.NewGEFV("band", dlit.MustNew(4)),
			Aggregators: map[string]*dlit.Literal{
				"numMatches":     dlit.MustNew("4"),
				"percentMatches": dlit.MustNew("100"),
				"numIncomeGt2":   dlit.MustNew("2"),
				"numBandGt4":     dlit.MustNew("2"),
				"goalsScore":     dlit.MustNew(0.002),
			},
			Goals: []*GoalAssessment{
				{"numIncomeGt2 == 1", false},
				{"numIncomeGt2 == 2", true},
				{"numIncomeGt2 == 3", false},
				{"numIncomeGt2 == 4", false},
				{"numBandGt4 == 1", false},
				{"numBandGt4 == 2", true},
				{"numBandGt4 == 3", false},
				{"numBandGt4 == 4", false},
			},
		},
		{
			Rule: rule.NewGEFV("cost", dlit.MustNew(1.3)),
			Aggregators: map[string]*dlit.Literal{
				"numMatches":     dlit.MustNew("2"),
				"percentMatches": dlit.MustNew("50"),
				"numIncomeGt2":   dlit.MustNew("2"),
				"numBandGt4":     dlit.MustNew("1"),
				"goalsScore":     dlit.MustNew(0.002),
			},
			Goals: []*GoalAssessment{
				{"numIncomeGt2 == 1", false},
				{"numIncomeGt2 == 2", true},
				{"numIncomeGt2 == 3", false},
				{"numIncomeGt2 == 4", false},
				{"numBandGt4 == 1", true},
				{"numBandGt4 == 2", false},
				{"numBandGt4 == 3", false},
				{"numBandGt4 == 4", false},
			},
		},
	}
	gotAssessment := New(aggregatorSpecs, goals)
	err = gotAssessment.AssessRules(dataset, rules)
	if err != nil {
		t.Errorf("AssessRules: %v", err)
	}

	assessmentsMatch := areAssessmentsEqv(
		gotAssessment,
		wantNumRecords,
		wantIsSorted,
		wantIsRefined,
		wantRuleAssessments,
	)
	if !assessmentsMatch {
		t.Errorf("AssessRules: assessments don't match")
		t.Errorf("got: %v", gotAssessment)
		t.Errorf("wantRuleAssessments: %v, wantNumRecords: %d, wantIsSorted: %t, wantIsRefined: %t",
			wantRuleAssessments,
			wantNumRecords, wantIsSorted, wantIsRefined)
	}
}

func TestAssessRules_goroutines(t *testing.T) {
	numGoRoutines := 10
	numRulesEach := 300
	if testing.Short() {
		numRulesEach = 20
	}
	allRules := make([][]rule.Rule, numGoRoutines)
	flatRules := make([]rule.Rule, numGoRoutines*numRulesEach)
	k := 0
	for i := 0; i < numGoRoutines; i++ {
		rules := make([]rule.Rule, numRulesEach)
		for j := 0; j < numRulesEach; j++ {
			rules[j] = rule.NewGEFV(
				"band",
				dlit.MustNew(0.001*(float64(i)*1100+float64(j))),
			)
			flatRules[k] = rules[j]
			k++
		}
		allRules[i] = rules
	}
	aggregatorDescs := []*aggregator.Desc{
		{"numIncomeGt2", "count", "income > 2"},
		{"numBandGt4", "count", "band > 4"},
	}
	goalExprs := []string{
		"numIncomeGt2 == 1",
		"numIncomeGt2 == 2",
		"numIncomeGt2 == 3",
		"numIncomeGt2 == 4",
		"numBandGt4 == 1",
		"numBandGt4 == 2",
		"numBandGt4 == 3",
		"numBandGt4 == 4",
	}
	sortOrder := []SortOrder{}
	fields := []string{"income", "cost", "band"}
	numRecords := 10000
	records := make([][]string, numRecords)
	for i := 0; i < numRecords; i++ {
		switch i % 4 {
		case 0:
			records[i] = []string{"3", "4.5", "4"}
		case 1:
			records[i] = []string{"3", "3.2", "7"}
		case 2:
			records[i] = []string{"2", "1.2", "4"}
		case 3:
			records[i] = []string{"0", "0", "9"}
		default:
			t.Fatalf("io dear")
		}
	}
	dataset := testhelpers.NewLiteralDataset(fields, records)
	aggregatorSpecs, err := aggregator.MakeSpecs(fields, aggregatorDescs)
	if err != nil {
		t.Fatalf("MakeSpecs: %s", err)
	}
	goals, err := goal.MakeGoals(goalExprs)
	if err != nil {
		t.Fatalf("MakeGoals: %s", err)
	}
	wantAssessment := New(aggregatorSpecs, goals)
	err = wantAssessment.AssessRules(dataset, flatRules)
	if err != nil {
		t.Fatalf("AssessRules: %v", err)
	}
	gotAssessment := New(aggregatorSpecs, goals)

	var wg sync.WaitGroup
	wg.Add(numGoRoutines)
	for i := 0; i < numGoRoutines; i++ {
		rules := allRules[i]
		go func() {
			err := gotAssessment.AssessRules(dataset, rules)
			defer wg.Done()
			if err != nil {
				t.Fatalf("AssessRules: %s", err)
			}
		}()
	}
	wg.Wait()
	gotAssessment.Sort(sortOrder)
	wantAssessment.Sort(sortOrder)

	if !wantAssessment.IsEqual(gotAssessment) {
		t.Errorf("AssessRules assessments don't match\n - got num ruleAssessments: %d, want: %d\n",
			len(gotAssessment.RuleAssessments), len(wantAssessment.RuleAssessments))
	}
}

func TestAssessRules_errors(t *testing.T) {
	cases := []struct {
		rules           []rule.Rule
		aggregatorDescs []*aggregator.Desc
		goalExprs       []string
		wantErr         error
	}{
		{[]rule.Rule{rule.NewGEFV("hand", dlit.MustNew(3))},
			[]*aggregator.Desc{
				{"numIncomeGt2", "count", "income > 2"},
			},
			[]string{"numIncomeGt2 == 1"},
			rule.InvalidRuleError{Rule: rule.NewGEFV("hand", dlit.MustNew(3))},
		},
		{[]rule.Rule{rule.NewGEFV("band", dlit.MustNew(3))},
			[]*aggregator.Desc{
				{"numIncomeGt2", "count", "bincome > 2"},
			},
			[]string{"numIncomeGt2 == 1"},
			dexpr.InvalidExprError{
				Expr: "bincome > 2",
				Err:  dexpr.VarNotExistError("bincome"),
			},
		},
		{[]rule.Rule{rule.NewGEFV("band", dlit.MustNew(3))},
			[]*aggregator.Desc{
				{"numIncomeGt2", "count", "income > 2"},
			},
			[]string{"numIncomeGt == 1"},
			dexpr.InvalidExprError{
				Expr: "numIncomeGt == 1",
				Err:  dexpr.VarNotExistError("numIncomeGt"),
			},
		},
	}
	fields := []string{"income", "cost", "band"}
	records := [][]string{
		{"3", "4.5", "4"},
	}
	dataset := testhelpers.NewLiteralDataset(fields, records)
	for i, c := range cases {
		aggregatorSpecs, err := aggregator.MakeSpecs(fields, c.aggregatorDescs)
		if err != nil {
			t.Fatalf("(%d) MakeSpecs: %s", i, err)
		}
		goals, err := goal.MakeGoals(c.goalExprs)
		if err != nil {
			t.Fatalf("(%d) MakeGoals: %s", i, err)
		}
		gotAssessment := New(aggregatorSpecs, goals)
		err = gotAssessment.AssessRules(dataset, c.rules)
		if err == nil || err.Error() != c.wantErr.Error() {
			t.Errorf("(%d) AssessRules - err: %s, wantErr: %s", i, err, c.wantErr)
		}
	}
}

func TestProcessRecord(t *testing.T) {
	rules := []rule.Rule{
		rule.NewGEFV("band", dlit.MustNew(5)),
		rule.NewGEFV("band", dlit.MustNew(4)),
		rule.NewGEFV("cost", dlit.MustNew(1.3)),
	}
	aggregatorDescs := []*aggregator.Desc{
		{"numIncomeGt2", "count", "income > 2"},
		{"numBandGt4", "count", "band > 4"},
	}
	goalExprs := []string{
		"numIncomeGt2 == 1",
		"numIncomeGt2 == 2",
		"numIncomeGt2 == 3",
		"numIncomeGt2 == 4",
		"numBandGt4 == 1",
		"numBandGt4 == 2",
		"numBandGt4 == 3",
		"numBandGt4 == 4",
	}
	fields := []string{"income", "cost", "band"}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(3),
		"cost":   dlit.MustNew(4.5),
		"band":   dlit.NewString("4"),
	}
	aggregatorSpecs, err := aggregator.MakeSpecs(fields, aggregatorDescs)
	if err != nil {
		t.Fatalf("MakeSpecs: %s", err)
	}
	goals, err := goal.MakeGoals(goalExprs)
	if err != nil {
		t.Fatalf("MakeGoals: %s", err)
	}
	wantIsSorted := false
	wantIsRefined := false
	wantNumRecords := int64(1)
	wantRuleAssessments := []*RuleAssessment{
		{
			Rule: rule.NewGEFV("band", dlit.MustNew(5)),
			Aggregators: map[string]*dlit.Literal{
				"numMatches":     dlit.MustNew("0"),
				"percentMatches": dlit.MustNew("0"),
				"numIncomeGt2":   dlit.MustNew("0"),
				"numBandGt4":     dlit.MustNew("0"),
				"goalsScore":     dlit.MustNew(0),
			},
			Goals: []*GoalAssessment{
				{"numIncomeGt2 == 1", false},
				{"numIncomeGt2 == 2", false},
				{"numIncomeGt2 == 3", false},
				{"numIncomeGt2 == 4", false},
				{"numBandGt4 == 1", false},
				{"numBandGt4 == 2", false},
				{"numBandGt4 == 3", false},
				{"numBandGt4 == 4", false},
			},
		},
		{
			Rule: rule.NewGEFV("band", dlit.MustNew(4)),
			Aggregators: map[string]*dlit.Literal{
				"numMatches":     dlit.MustNew("1"),
				"percentMatches": dlit.MustNew("100"),
				"numIncomeGt2":   dlit.MustNew("1"),
				"numBandGt4":     dlit.MustNew("0"),
				"goalsScore":     dlit.MustNew(1),
			},
			Goals: []*GoalAssessment{
				{"numIncomeGt2 == 1", true},
				{"numIncomeGt2 == 2", false},
				{"numIncomeGt2 == 3", false},
				{"numIncomeGt2 == 4", false},
				{"numBandGt4 == 1", false},
				{"numBandGt4 == 2", false},
				{"numBandGt4 == 3", false},
				{"numBandGt4 == 4", false},
			},
		},
		{
			Rule: rule.NewGEFV("cost", dlit.MustNew(1.3)),
			Aggregators: map[string]*dlit.Literal{
				"numMatches":     dlit.MustNew("1"),
				"percentMatches": dlit.MustNew("100"),
				"numIncomeGt2":   dlit.MustNew("1"),
				"numBandGt4":     dlit.MustNew("0"),
				"goalsScore":     dlit.MustNew(1),
			},
			Goals: []*GoalAssessment{
				{"numIncomeGt2 == 1", true},
				{"numIncomeGt2 == 2", false},
				{"numIncomeGt2 == 3", false},
				{"numIncomeGt2 == 4", false},
				{"numBandGt4 == 1", false},
				{"numBandGt4 == 2", false},
				{"numBandGt4 == 3", false},
				{"numBandGt4 == 4", false},
			},
		},
	}
	gotAssessment := New(aggregatorSpecs, goals)
	gotAssessment.AddRules(rules)
	if err := gotAssessment.ProcessRecord(record); err != nil {
		t.Fatalf("ProcessRecord: %s", err)
	}
	if err := gotAssessment.Update(); err != nil {
		t.Fatalf("Update: %s", err)
	}

	assessmentsMatch := areAssessmentsEqv(
		gotAssessment,
		wantNumRecords,
		wantIsSorted,
		wantIsRefined,
		wantRuleAssessments,
	)
	if !assessmentsMatch {
		t.Errorf("ProcessRecord: assessments don't match")
		t.Errorf("got: %v", gotAssessment)
		t.Errorf("wantRuleAssessments: %v, wantNumRecords: %d, wantIsSorted: %t, wantIsRefined: %t",
			wantRuleAssessments,
			wantNumRecords, wantIsSorted, wantIsRefined)
	}
}

func TestAssessmentAssessRules_dataset_changed_error(t *testing.T) {
	rules := []rule.Rule{
		rule.NewGEFV("band", dlit.MustNew(5)),
		rule.NewGEFV("band", dlit.MustNew(4)),
		rule.NewGEFV("cost", dlit.MustNew(1.3)),
	}
	aggregatorDescs := []*aggregator.Desc{
		{"numIncomeGt2", "count", "income > 2"},
		{"numBandGt4", "count", "band > 4"},
	}
	goalExprs := []string{
		"numIncomeGt2 == 1",
		"numBandGt4 == 4",
	}
	fields := []string{"income", "cost", "band"}
	records := [][]string{
		{"3", "4.5", "4"},
		{"3", "3.2", "7"},
		{"2", "1.2", "4"},
		{"0", "0", "9"},
	}
	dataset := testhelpers.NewLiteralDataset(fields, records)
	aggregatorSpecs, err := aggregator.MakeSpecs(fields, aggregatorDescs)
	if err != nil {
		t.Fatalf("MakeSpecs: %s", err)
	}
	goals, err := goal.MakeGoals(goalExprs)
	if err != nil {
		t.Fatalf("MakeGoals: %s", err)
	}
	trueRules := []rule.Rule{rule.NewTrue()}
	wantErr := ErrNumRecordsChanged
	assessment := New(aggregatorSpecs, goals)
	err = assessment.AssessRules(dataset, trueRules)
	if err != nil {
		t.Fatalf("AssessRules: %v", err)
	}
	tDataset := dtruncate.New(dataset, int64(len(records)-1))
	err = assessment.AssessRules(tDataset, rules)
	if err == nil || err.Error() != wantErr.Error() {
		t.Errorf("AssessRules - err: %s, wantErr: %s", err, wantErr)
	}
}

func TestAssessmentAddRuleAssessments_error(t *testing.T) {
	aggregatorSpecs := []aggregator.Spec{aggregator.MustNew("a", "calc", "3+4")}
	goals := []*goal.Goal{goal.MustNew("cost > 3")}
	ruleAssessments := []*RuleAssessment{
		newRuleAssessment(
			rule.NewEQFV("month", dlit.NewString("May")),
			aggregatorSpecs,
			goals,
		),
	}
	wantErr := dexpr.InvalidExprError{
		Expr: "cost > 3",
		Err:  dexpr.VarNotExistError("cost"),
	}
	assessment := New(aggregatorSpecs, goals)
	err := assessment.addRuleAssessments(ruleAssessments)
	if err == nil || err.Error() != wantErr.Error() {
		t.Errorf("AddRuleAssessors: err: %s, wantErr: %s", err, wantErr)
	}
}

func TestRules(t *testing.T) {
	var gotRules []rule.Rule
	assessment := Assessment{NumRecords: 8,
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(9)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", true},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(456)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(3)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("76.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", true},
				},
			},
			{
				Rule: rule.NewGEFV("cost", dlit.MustNew(1.2)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
		},
	}
	cases := []struct {
		numRules     int
		passNumRules bool
		wantRules    []rule.Rule
	}{
		{0, true, []rule.Rule{}},
		{1, true, []rule.Rule{rule.NewGEFV("band", dlit.MustNew(9))}},
		{2, true, []rule.Rule{
			rule.NewGEFV("band", dlit.MustNew(9)),
			rule.NewGEFV("band", dlit.MustNew(456)),
		}},
		{4, true, []rule.Rule{
			rule.NewGEFV("band", dlit.MustNew(9)),
			rule.NewGEFV("band", dlit.MustNew(456)),
			rule.NewGEFV("band", dlit.MustNew(3)),
			rule.NewGEFV("cost", dlit.MustNew(1.2)),
		},
		},
		{5, true, []rule.Rule{
			rule.NewGEFV("band", dlit.MustNew(9)),
			rule.NewGEFV("band", dlit.MustNew(456)),
			rule.NewGEFV("band", dlit.MustNew(3)),
			rule.NewGEFV("cost", dlit.MustNew(1.2)),
		},
		},
		{0, false, []rule.Rule{
			rule.NewGEFV("band", dlit.MustNew(9)),
			rule.NewGEFV("band", dlit.MustNew(456)),
			rule.NewGEFV("band", dlit.MustNew(3)),
			rule.NewGEFV("cost", dlit.MustNew(1.2)),
		},
		},
	}
	for _, c := range cases {
		if c.passNumRules {
			gotRules = assessment.Rules(c.numRules)
		} else {
			gotRules = assessment.Rules()
		}
		if !reflect.DeepEqual(gotRules, c.wantRules) {
			t.Errorf("Rules() passNumRules: %t, numRules: %d rules don't match\ngot: %s\nwant: %s\n",
				c.passNumRules, c.numRules, gotRules, c.wantRules)
		}
	}
}

func TestMerge(t *testing.T) {
	assessment1 := &Assessment{
		NumRecords: 8,
		flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(9)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", true},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(456)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(3)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("76.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", true},
				},
			},
			{
				Rule: rule.NewGEFV("cost", dlit.MustNew(1.2)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
		},
	}
	assessment2 := &Assessment{
		NumRecords: 8,
		flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(16)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("8"),
					"percentMatches": dlit.MustNew("5.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", true},
				},
			},
			{
				Rule: rule.NewEQFV("team", dlit.MustNew("Pi")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("19"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(36)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("6.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
			{
				Rule: rule.NewGEFV("cost", dlit.MustNew(1.27)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("3.5"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
		},
	}

	wantAssessment := &Assessment{
		NumRecords: 8,
		flags: map[string]bool{
			"sorted":  false,
			"refined": false,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(9)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", true},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(456)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(3)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("76.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", true},
				},
			},
			{
				Rule: rule.NewGEFV("cost", dlit.MustNew(1.2)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(16)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("8"),
					"percentMatches": dlit.MustNew("5.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", true},
				},
			},
			{
				Rule: rule.NewEQFV("team", dlit.MustNew("Pi")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("19"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(36)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("6.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
			{
				Rule: rule.NewGEFV("cost", dlit.MustNew(1.27)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("3.5"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
		},
	}
	gotAssessment, err := assessment1.Merge(assessment2)
	if err != nil {
		t.Errorf("Merge() error: %s", err)
		return
	}

	if !gotAssessment.IsEqual(wantAssessment) {
		t.Errorf("Merge() assessments don't match\n - got: %v\n - want: %v\n",
			gotAssessment, wantAssessment)
	}
}

func TestMerge_errors(t *testing.T) {
	assessment1 := &Assessment{
		NumRecords: 8,
		flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(9)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", true},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(456)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
		},
	}
	assessment2 := &Assessment{NumRecords: 2,
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(16)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("8"),
					"percentMatches": dlit.MustNew("5.3"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", true},
				},
			},
			{
				Rule: rule.NewEQFV("team", dlit.MustNew("Pi")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("3"),
					"percentMatches": dlit.MustNew("19"),
				},
				Goals: []*GoalAssessment{
					{"numMatches > 3 ", false},
				},
			},
		},
	}
	wantErr := ErrNumRecordsChanged
	_, err := assessment1.Merge(assessment2)
	if err == nil || err != wantErr {
		t.Errorf("Merge - err: %s, wantErr: %s", err, wantErr)
	}
}

func TestRefine(t *testing.T) {
	sortedAssessment := &Assessment{
		NumRecords: 20,
		flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("0.005"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(5)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("150"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.NewInFV(
					"band",
					testhelpers.MakeStringsDlitSlice("4", "3", "2"),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("149"),
					"percentMatches": dlit.MustNew("49"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
			{
				Rule: rule.NewInFV("team", testhelpers.MakeStringsDlitSlice("a", "b")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("148"),
					"percentMatches": dlit.MustNew("48"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
			{
				Rule: rule.NewInFV(
					"band",
					testhelpers.MakeStringsDlitSlice("99", "23"),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("147"),
					"percentMatches": dlit.MustNew("47"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(3)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("146"),
					"percentMatches": dlit.MustNew("46"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(9)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("145"),
					"percentMatches": dlit.MustNew("45"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
			{
				Rule: rule.NewInFV("band", testhelpers.MakeStringsDlitSlice("9", "2")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("144"),
					"percentMatches": dlit.MustNew("44"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
			{
				Rule: rule.NewEQFV("band", dlit.MustNew(7)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("1"),
					"percentMatches": dlit.MustNew("0.0045"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
			{
				Rule: rule.NewEQFV("band", dlit.MustNew(8)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("142"),
					"percentMatches": dlit.MustNew("42"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
			{
				Rule: rule.NewTrue(),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("142"),
					"percentMatches": dlit.MustNew("42"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
			{
				Rule: rule.NewGEFV("cost", dlit.MustNew(1.2)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("141"),
					"percentMatches": dlit.MustNew("41"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
		},
	}
	wantRules := []rule.Rule{
		rule.NewGEFV("band", dlit.MustNew(4)),
		rule.NewInFV("band", testhelpers.MakeStringsDlitSlice("4", "3", "2")),
		rule.NewInFV("team", testhelpers.MakeStringsDlitSlice("a", "b")),
		rule.NewInFV("band", testhelpers.MakeStringsDlitSlice("99", "23")),
		rule.NewTrue(),
	}
	sortedAssessment.Refine()
	gotRules := sortedAssessment.Rules()

	if !matchRules(gotRules, wantRules) {
		t.Errorf("matchRules() rules don't match:\ngot: %s\nwant: %s\n",
			gotRules, wantRules)
	}
}

func TestRefine_few_ruleassessments(t *testing.T) {
	cases := []struct {
		in        *Assessment
		wantRules []rule.Rule
	}{
		{in: &Assessment{
			NumRecords: 20,
			flags: map[string]bool{
				"sorted": true,
			},
			RuleAssessments: []*RuleAssessment{
				{
					Rule: rule.NewTrue(),
					Aggregators: map[string]*dlit.Literal{
						"numMatches":     dlit.MustNew("142"),
						"percentMatches": dlit.MustNew("42"),
						"numIncomeGt2":   dlit.MustNew("2"),
						"goalsScore":     dlit.MustNew(0.1),
					},
					Goals: []*GoalAssessment{
						{"numIncomeGt2 == 1", false},
						{"numIncomeGt2 == 2", true},
					},
				},
			},
		},
			wantRules: []rule.Rule{rule.NewTrue()},
		},
		{in: &Assessment{
			NumRecords: 20,
			flags: map[string]bool{
				"sorted": true,
			},
			RuleAssessments: []*RuleAssessment{
				{
					Rule: rule.NewEQFV("month", dlit.NewString("april")),
					Aggregators: map[string]*dlit.Literal{
						"numMatches":     dlit.MustNew("142"),
						"percentMatches": dlit.MustNew("42"),
						"goalsScore":     dlit.MustNew(0.1),
					},
					Goals: []*GoalAssessment{
						{"numIncomeGt2 == 1", false},
						{"numIncomeGt2 == 2", true},
					},
				},
				{
					Rule: rule.NewTrue(),
					Aggregators: map[string]*dlit.Literal{
						"numMatches":     dlit.MustNew("142"),
						"percentMatches": dlit.MustNew("42"),
						"numIncomeGt2":   dlit.MustNew("2"),
						"goalsScore":     dlit.MustNew(0.1),
					},
					Goals: []*GoalAssessment{
						{"numIncomeGt2 == 1", false},
						{"numIncomeGt2 == 2", true},
					},
				},
			},
		},
			wantRules: []rule.Rule{
				rule.NewEQFV("month", dlit.NewString("april")),
				rule.NewTrue(),
			},
		},
	}
	for _, c := range cases {
		c.in.Refine()
		gotRules := c.in.Rules()

		if !matchRules(gotRules, c.wantRules) {
			t.Errorf("matchRules() rules don't match:\ngot: %s\nwant: %s\n",
				gotRules, c.wantRules)
		}
	}
}

func TestRefine_between(t *testing.T) {
	sortedAssessment := &Assessment{
		NumRecords: 20,
		flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.MustNewBetweenFV("band", dlit.MustNew(5), dlit.MustNew(7)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("150"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("149"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewBetweenFV(
					"rate",
					dlit.MustNew(16.2),
					dlit.MustNew(17.93),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("148"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewBetweenFV("band", dlit.MustNew(5), dlit.MustNew(6)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("147"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewBetweenFV(
					"rate",
					dlit.MustNew(16.2),
					dlit.MustNew(17.89),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("146"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(5)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("145"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewBetweenFV("band", dlit.MustNew(7), dlit.MustNew(16)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("144"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewBetweenFV(
					"rate",
					dlit.MustNew(15.2),
					dlit.MustNew(17.89),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("143"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewBetweenFV(
					"rate",
					dlit.MustNew(10.2),
					dlit.MustNew(16.3),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("142"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewBetweenFV(
					"rate",
					dlit.MustNew(50.1),
					dlit.MustNew(60.3),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("141"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewBetweenFV(
					"rate",
					dlit.MustNew(16.1),
					dlit.MustNew(20.3),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("140"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.NewLEFV("band", dlit.MustNew(5)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("139"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(6)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("138"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewBetweenFV(
					"rate",
					dlit.MustNew(1.2),
					dlit.MustNew(6.3),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("137"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.NewTrue(),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("136"),
					"percentMatches": dlit.MustNew("42"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
		},
	}
	wantRules := []rule.Rule{
		rule.MustNewBetweenFV("band", dlit.MustNew(5), dlit.MustNew(7)),
		rule.NewGEFV("band", dlit.MustNew(4)),
		rule.MustNewBetweenFV(
			"rate",
			dlit.MustNew(16.2),
			dlit.MustNew(17.93),
		),
		rule.MustNewBetweenFV(
			"rate",
			dlit.MustNew(50.1),
			dlit.MustNew(60.3),
		),
		rule.NewLEFV("band", dlit.MustNew(5)),
		rule.MustNewBetweenFV(
			"rate",
			dlit.MustNew(1.2),
			dlit.MustNew(6.3),
		),
		rule.NewTrue(),
	}
	sortedAssessment.Refine()
	gotRules := sortedAssessment.Rules()

	if !matchRules(gotRules, wantRules) {
		t.Errorf("matchRules() rules don't match:\ngot: %s\nwant: %s\n",
			gotRules, wantRules)
	}
}

func TestRefine_outside(t *testing.T) {
	sortedAssessment := &Assessment{
		NumRecords: 20,
		flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.MustNewOutsideFV("band", dlit.MustNew(5), dlit.MustNew(7)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("150"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("149"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewOutsideFV(
					"rate",
					dlit.MustNew(16.2),
					dlit.MustNew(17.93),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("148"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewOutsideFV("band", dlit.MustNew(5), dlit.MustNew(6)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("147"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewOutsideFV(
					"rate",
					dlit.MustNew(16.2),
					dlit.MustNew(17.89),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("146"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(5)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("145"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewOutsideFV("band", dlit.MustNew(7), dlit.MustNew(16)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("144"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.MustNewOutsideFV(
					"rate",
					dlit.MustNew(15.2),
					dlit.MustNew(17.89),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("143"),
					"percentMatches": dlit.MustNew("50"),
					"numIncomeGt2":   dlit.MustNew("1"),
					"goalsScore":     dlit.MustNew(1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
				},
			},
			{
				Rule: rule.NewTrue(),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("142"),
					"percentMatches": dlit.MustNew("42"),
					"numIncomeGt2":   dlit.MustNew("2"),
					"goalsScore":     dlit.MustNew(0.1),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
				},
			},
		},
	}
	wantRules := []rule.Rule{
		rule.MustNewOutsideFV("band", dlit.MustNew(5), dlit.MustNew(7)),
		rule.NewGEFV("band", dlit.MustNew(4)),
		rule.MustNewOutsideFV("rate", dlit.MustNew(16.2), dlit.MustNew(17.93)),
		rule.NewTrue(),
	}
	sortedAssessment.Refine()
	gotRules := sortedAssessment.Rules()

	if !matchRules(gotRules, wantRules) {
		t.Errorf("matchRules() rules don't match:\ngot: %s\nwant: %s\n",
			gotRules, wantRules)
	}
}

func TestRefine_panic_1(t *testing.T) {
	testPurpose := "Ensure panics if assessment not sorted"
	unsortedAssessment := &Assessment{
		NumRecords: 20,
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{},
			},
			{
				Rule: rule.NewTrue(),
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
	unsortedAssessment.Refine()
	if !paniced {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("Refine() - failed to panic with: %s", wantPanic)
	}
}

func TestRefine_panic_2(t *testing.T) {
	testPurpose := "Ensure panics if True rule missing"
	sortedAssessment := &Assessment{
		NumRecords: 20,
		flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{},
			},
			{
				Rule: rule.NewGEFV("team", dlit.MustNew(7)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("100"),
				},
				Goals: []*GoalAssessment{},
			},
		},
	}
	paniced := false
	wantPanic := "No True rule found"
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
	sortedAssessment.Refine()
	if !paniced {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("Refine() - failed to panic with: %s", wantPanic)
	}
}

func TestTruncateRuleAssessments(t *testing.T) {
	refinedAssessment := &Assessment{
		NumRecords: 20,
		flags: map[string]bool{
			"sorted":  true,
			"refined": true,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches": dlit.MustNew("2"),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
				},
			},
			{
				Rule: rule.NewInFV(
					"band",
					testhelpers.MakeStringsDlitSlice("4", "3", "2"),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches": dlit.MustNew("4"),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
				},
			},
			{
				Rule: rule.NewInFV("team", testhelpers.MakeStringsDlitSlice("a", "b")),
				Aggregators: map[string]*dlit.Literal{
					"numMatches": dlit.MustNew("4"),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
				},
			},
			{
				Rule: rule.NewInFV(
					"band",
					testhelpers.MakeStringsDlitSlice("99", "23"),
				),
				Aggregators: map[string]*dlit.Literal{
					"numMatches": dlit.MustNew("4"),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
				},
			},
			{
				Rule: rule.NewTrue(),
				Aggregators: map[string]*dlit.Literal{
					"numMatches": dlit.MustNew("2"),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
				},
			},
		},
	}
	// The increasing order of the numRules is important as this also checks that
	// the ruleAssessments are cloned properly
	cases := []struct {
		numRules  int
		wantRules []rule.Rule
	}{
		{0,
			[]rule.Rule{},
		},
		{1,
			[]rule.Rule{
				rule.NewTrue(),
			},
		},
		{2,
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(4)),
				rule.NewTrue(),
			},
		},
		{3,
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(4)),
				rule.NewInFV("band", testhelpers.MakeStringsDlitSlice("4", "3", "2")),
				rule.NewTrue(),
			},
		},
		{4,
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(4)),
				rule.NewInFV("band", testhelpers.MakeStringsDlitSlice("4", "3", "2")),
				rule.NewInFV("team", testhelpers.MakeStringsDlitSlice("a", "b")),
				rule.NewTrue(),
			},
		},
		{5,
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(4)),
				rule.NewInFV("band", testhelpers.MakeStringsDlitSlice("4", "3", "2")),
				rule.NewInFV("team", testhelpers.MakeStringsDlitSlice("a", "b")),
				rule.NewInFV("band", testhelpers.MakeStringsDlitSlice("99", "23")),
				rule.NewTrue(),
			},
		},
		{6,
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(4)),
				rule.NewInFV("band", testhelpers.MakeStringsDlitSlice("4", "3", "2")),
				rule.NewInFV("team", testhelpers.MakeStringsDlitSlice("a", "b")),
				rule.NewInFV("band", testhelpers.MakeStringsDlitSlice("99", "23")),
				rule.NewTrue(),
			},
		},
	}
	for _, c := range cases {
		limitedAssessment := refinedAssessment.TruncateRuleAssessments(c.numRules)
		gotRules := limitedAssessment.Rules()
		if !matchRules(gotRules, c.wantRules) {
			t.Errorf("matchRules() rules don't match:\nnumRules: %d\ngot: %s\nwant: %s\n",
				c.numRules, gotRules, c.wantRules)
		}
	}
}

func TestTruncateRuleAssessment_panic_1(t *testing.T) {
	testPurpose := "Ensure panics if assessment not sorted"
	unsortedAssessment := &Assessment{
		NumRecords: 20,
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{},
			},
			{
				Rule: rule.NewTrue(),
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
				t.Errorf("TruncateRuleAssessments() - got panic: %s, wanted: %s",
					r, wantPanic)
			}
		}
	}()
	numRules := 1
	unsortedAssessment.TruncateRuleAssessments(numRules)
	if !paniced {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("Refine() - failed to panic with: %s", wantPanic)
	}
}

func TestTruncateRuleAssessment_panic_2(t *testing.T) {
	testPurpose := "Ensure panics if assessment not refined"
	unsortedAssessment := &Assessment{
		NumRecords: 20,
		flags: map[string]bool{
			"sorted": true,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(4)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
				},
				Goals: []*GoalAssessment{},
			},
			{
				Rule: rule.NewTrue(),
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
				t.Errorf("TruncateRuleAssessments() - got panic: %s, wanted: %s",
					r, wantPanic)
			}
		}
	}()
	numRules := 1
	unsortedAssessment.TruncateRuleAssessments(numRules)
	if !paniced {
		t.Errorf("Test: %s\n", testPurpose)
		t.Errorf("Refine() - failed to panic with: %s", wantPanic)
	}
}

func TestSort(t *testing.T) {
	assessment := Assessment{
		NumRecords: 8,
		flags: map[string]bool{
			"sorted": false,
		},
		RuleAssessments: []*RuleAssessment{
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(9)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("5"),
					"percentMatches": dlit.MustNew("65.3"),
					"goalsScore":     dlit.MustNew(0.003),
					"numIncomeGt2":   dlit.MustNew("3"),
					"numBandGt4":     dlit.MustNew("2"),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
					{"numIncomeGt2 == 3", false},
					{"numIncomeGt2 == 4", false},
					{"numBandGt4 == 1", false},
					{"numBandGt4 == 2", true},
					{"numBandGt4 == 3", false},
					{"numBandGt4 == 4", true},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(456)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
					"goalsScore":     dlit.MustNew(1.001),
					"numIncomeGt2":   dlit.MustNew("2"),
					"numBandGt4":     dlit.MustNew("2"),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", true},
					{"numIncomeGt2 == 2", false},
					{"numIncomeGt2 == 3", false},
					{"numIncomeGt2 == 4", false},
					{"numBandGt4 == 1", false},
					{"numBandGt4 == 2", true},
					{"numBandGt4 == 3", false},
					{"numBandGt4 == 4", false},
				},
			},
			{
				Rule: rule.NewGEFV("band", dlit.MustNew(3)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("4"),
					"percentMatches": dlit.MustNew("76.3"),
					"goalsScore":     dlit.MustNew(0.002),
					"numIncomeGt2":   dlit.MustNew("2"),
					"numBandGt4":     dlit.MustNew("2"),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
					{"numIncomeGt2 == 3", false},
					{"numIncomeGt2 == 4", false},
					{"numBandGt4 == 1", false},
					{"numBandGt4 == 2", true},
					{"numBandGt4 == 3", false},
					{"numBandGt4 == 4", false},
				},
			},
			{
				Rule: rule.NewGEFV("cost", dlit.MustNew(1.2)),
				Aggregators: map[string]*dlit.Literal{
					"numMatches":     dlit.MustNew("2"),
					"percentMatches": dlit.MustNew("50"),
					"goalsScore":     dlit.MustNew(0.002),
					"numIncomeGt2":   dlit.MustNew("1"),
					"numBandGt4":     dlit.MustNew("1"),
				},
				Goals: []*GoalAssessment{
					{"numIncomeGt2 == 1", false},
					{"numIncomeGt2 == 2", true},
					{"numIncomeGt2 == 3", false},
					{"numIncomeGt2 == 4", false},
					{"numBandGt4 == 1", true},
					{"numBandGt4 == 2", false},
					{"numBandGt4 == 3", false},
					{"numBandGt4 == 4", false},
				},
			},
		},
	}
	cases := []struct {
		sortOrder []SortOrder
		wantRules []rule.Rule
	}{
		{[]SortOrder{
			{"goalsScore", ASCENDING},
		},
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(3)),
				rule.NewGEFV("cost", dlit.MustNew(1.2)),
				rule.NewGEFV("band", dlit.MustNew(9)),
				rule.NewGEFV("band", dlit.MustNew(456)),
			}},
		{[]SortOrder{
			{"percentMatches", DESCENDING},
		},
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(3)),
				rule.NewGEFV("band", dlit.MustNew(9)),
				rule.NewGEFV("band", dlit.MustNew(456)),
				rule.NewGEFV("cost", dlit.MustNew(1.2)),
			}},
		{[]SortOrder{
			{"percentMatches", ASCENDING},
		},
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(456)),
				rule.NewGEFV("cost", dlit.MustNew(1.2)),
				rule.NewGEFV("band", dlit.MustNew(9)),
				rule.NewGEFV("band", dlit.MustNew(3)),
			}},
		{[]SortOrder{
			{"percentMatches", ASCENDING},
			{"numIncomeGt2", ASCENDING},
		},
			[]rule.Rule{
				rule.NewGEFV("cost", dlit.MustNew(1.2)),
				rule.NewGEFV("band", dlit.MustNew(456)),
				rule.NewGEFV("band", dlit.MustNew(9)),
				rule.NewGEFV("band", dlit.MustNew(3)),
			}},
		{[]SortOrder{
			{"percentMatches", DESCENDING},
			{"numIncomeGt2", ASCENDING},
		},
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(3)),
				rule.NewGEFV("band", dlit.MustNew(9)),
				rule.NewGEFV("cost", dlit.MustNew(1.2)),
				rule.NewGEFV("band", dlit.MustNew(456)),
			}},
		{[]SortOrder{},
			[]rule.Rule{
				rule.NewGEFV("band", dlit.MustNew(3)),
				rule.NewGEFV("band", dlit.MustNew(9)),
				rule.NewGEFV("band", dlit.MustNew(456)),
				rule.NewGEFV("cost", dlit.MustNew(1.2)),
			}},
	}
	for _, c := range cases {
		assessment.Sort(c.sortOrder)
		if !assessment.IsSorted() {
			t.Errorf("Sort(%s) 'sorted' flag not set", c.sortOrder)
		}
		gotRules := assessment.Rules()
		if !matchRules(gotRules, c.wantRules) {
			t.Errorf("matchRules() rules don't match:\n - sortOrder: %s\n - got: %s\n - want: %s\n",
				c.sortOrder, gotRules, c.wantRules)
		}
	}
}

/*************************
       Benchmarks
*************************/
func BenchmarkAssessRules(b *testing.B) {
	b.StopTimer()
	var numRules int64 = 3000
	rules := make([]rule.Rule, numRules)
	for i := int64(0); i < numRules; i++ {
		if i%2 == 0 {
			rules[i] = rule.NewGEFV("age", dlit.MustNew(i%50))
		} else {
			rules[i] = rule.NewGEFV("day", dlit.MustNew(i%20))
		}
	}

	fields := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y"}
	dataset := dcsv.New(
		filepath.Join("fixtures", "bank_big.csv"),
		true,
		rune(';'),
		fields,
	)
	aggregatorDescs := []*aggregator.Desc{
		{"numMarried", "count", "marital == \"married\""},
		{"numSignedUp", "count", "y == \"yes\""},
		{"cost", "calc", "numMatches * 4.5"},
		{"income", "calc", "numSignedUp * 24"},
		{"profit", "calc", "income - cost"},
	}
	goalExprs := []string{
		"profit > 0",
		"numSignedUp > 3",
		"numMarried > 2",
	}
	aggregatorSpecs, err := aggregator.MakeSpecs(fields, aggregatorDescs)
	if err != nil {
		b.Fatalf("MakeSpecs: %s", err)
	}
	goals, err := goal.MakeGoals(goalExprs)
	if err != nil {
		b.Fatalf("MakeGoals: %s", err)
	}
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		assessment := New(aggregatorSpecs, goals)
		err := assessment.AssessRules(dataset, rules)
		if err != nil {
			b.Errorf("AssessRules: %s", err)
		}
	}
}

/******************************
 *  Helper functions
 ******************************/

// Match the rules including their order
func matchRules(rules1 []rule.Rule, rules2 []rule.Rule) bool {
	if len(rules1) != len(rules2) {
		return false
	}
	for i, rule1 := range rules1 {
		if rule1.String() != rules2[i].String() {
			return false
		}
	}
	return true
}

// Are the assessments equivalent.  The ruleAssessments must match
// but don't have to be in the same order if both assessments are
// unsorted. If both are unsorted then this will sort the assessments
func areAssessmentsEqv(
	got *Assessment,
	wantNumRecords int64,
	wantIsSorted bool,
	wantIsRefined bool,
	wantRuleAssessments []*RuleAssessment,
) bool {
	if got.NumRecords != wantNumRecords {
		return false
	}
	if got.IsSorted() != wantIsSorted {
		return false
	}
	if got.IsRefined() != wantIsRefined {
		return false
	}
	if len(got.RuleAssessments) != len(wantRuleAssessments) {
		return false
	}
	for _, gotRuleAssesment := range got.RuleAssessments {
		found := false
		for _, wantRuleAssessment := range wantRuleAssessments {
			if gotRuleAssesment.IsEqual(wantRuleAssessment) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
