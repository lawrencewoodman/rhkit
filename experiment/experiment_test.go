package experiment

import (
	"errors"
	"github.com/lawrencewoodman/ddataset"
	"github.com/lawrencewoodman/ddataset/dcsv"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/aggregators"
	"github.com/vlifesystems/rhkit/goal"
	"github.com/vlifesystems/rhkit/rule"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	// Field: 'p_1234567890outcome' is there to check allowed characters
	fieldNames := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "p_1234567890outcome", "y"}
	expectedExperiments := []*Experiment{
		{
			Dataset: dcsv.New(
				filepath.Join("..", "fixtures", "bank.csv"),
				true,
				rune(';'),
				fieldNames,
			),
			RuleFields: []string{"age", "job", "marital", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "p_1234567890outcome", "y",
			},
			RuleComplexity: rule.Complexity{Arithmetic: true},
			Aggregators: []aggregators.AggregatorSpec{
				aggregators.MustNew("numMatches", "count", "true()"),
				aggregators.MustNew("percentMatches", "calc",
					"roundto(100.0 * numMatches / numRecords, 2)"),
				// num_married to check for allowed characters
				aggregators.MustNew("num_married", "count", "marital == \"married\""),
				aggregators.MustNew("numSignedUp", "count", "y == \"yes\""),
				aggregators.MustNew("cost", "calc", "numMatches * 4.5"),
				aggregators.MustNew("income", "calc", "numSignedUp * 24"),
				aggregators.MustNew("profit", "calc", "income - cost"),
				aggregators.MustNew("goalsScore", "goalsscore"),
			},
			Goals: []*goal.Goal{goal.MustNew("profit > 0")},
			SortOrder: []SortField{
				{"profit", DESCENDING},
				{"numSignedUp", DESCENDING},
				{"cost", ASCENDING},
				{"numMatches", DESCENDING},
				{"percentMatches", DESCENDING},
				{"goalsScore", DESCENDING},
			},
		},
		{
			Dataset: dcsv.New(
				filepath.Join("..", "fixtures", "bank.csv"),
				true,
				rune(';'),
				fieldNames,
			),
			RuleFields: []string{"age", "job", "marital", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "p_1234567890outcome", "y",
			},
			RuleComplexity: rule.Complexity{Arithmetic: false},
			Aggregators: []aggregators.AggregatorSpec{
				aggregators.MustNew("numMatches", "count", "true()"),
				aggregators.MustNew("percentMatches", "calc",
					"roundto(100.0 * numMatches / numRecords, 2)"),
				// num_married to check for allowed characters
				aggregators.MustNew("num_married", "count", "marital == \"married\""),
				aggregators.MustNew("numSignedUp", "count", "y == \"yes\""),
				aggregators.MustNew("cost", "calc", "numMatches * 4.5"),
				aggregators.MustNew("income", "calc", "numSignedUp * 24"),
				aggregators.MustNew("profit", "calc", "income - cost"),
				aggregators.MustNew("goalsScore", "goalsscore"),
			},
			Goals: []*goal.Goal{goal.MustNew("profit > 0")},
			SortOrder: []SortField{
				{"profit", DESCENDING},
				{"numSignedUp", DESCENDING},
				{"cost", ASCENDING},
				{"numMatches", DESCENDING},
				{"percentMatches", DESCENDING},
				{"goalsScore", DESCENDING},
			},
		},
		{
			Dataset: dcsv.New(
				filepath.Join("..", "fixtures", "bank.csv"),
				true,
				rune(';'),
				fieldNames,
			),
			RuleFields: []string{"age", "job", "marital", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "p_1234567890outcome", "y",
			},
			RuleComplexity: rule.Complexity{Arithmetic: false},
			Aggregators: []aggregators.AggregatorSpec{
				aggregators.MustNew("numMatches", "count", "true()"),
				aggregators.MustNew("percentMatches", "calc",
					"roundto(100.0 * numMatches / numRecords, 2)"),
				// num_married to check for allowed characters
				aggregators.MustNew("num_married", "count", "marital == \"married\""),
				aggregators.MustNew("numSignedUp", "count", "y == \"yes\""),
				aggregators.MustNew("cost", "calc", "numMatches * 4.5"),
				aggregators.MustNew("income", "calc", "numSignedUp * 24"),
				aggregators.MustNew("profit", "calc", "income - cost"),
				aggregators.MustNew("goalsScore", "goalsscore"),
			},
			Goals: []*goal.Goal{goal.MustNew("profit > 0")},
			SortOrder: []SortField{
				{"profit", DESCENDING},
				{"numSignedUp", DESCENDING},
				{"cost", ASCENDING},
				{"numMatches", DESCENDING},
				{"percentMatches", DESCENDING},
				{"goalsScore", DESCENDING},
			},
			Rules: []rule.Rule{
				rule.NewEQFV("job", dlit.MustNew("manager")),
				rule.NewGEFV("age", dlit.MustNew(27)),
				rule.NewLEFV("balance", dlit.MustNew(1500)),
			},
		},
	}
	cases := []struct {
		experimentDesc *ExperimentDesc
		want           *Experiment
	}{
		{&ExperimentDesc{
			Dataset: dcsv.New(
				filepath.Join("..", "fixtures", "bank.csv"),
				true,
				rune(';'),
				fieldNames,
			),
			RuleFields: []string{"age", "job", "marital", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "p_1234567890outcome", "y",
			},
			RuleComplexity: rule.Complexity{Arithmetic: true},
			Aggregators: []*AggregatorDesc{
				{"num_married", "count", "marital == \"married\""},
				{"numSignedUp", "count", "y == \"yes\""},
				{"cost", "calc", "numMatches * 4.5"},
				{"income", "calc", "numSignedUp * 24"},
				{"profit", "calc", "income - cost"},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				{"profit", "descending"},
				{"numSignedUp", "descending"},
				{"cost", "ascending"},
				{"numMatches", "descending"},
				{"percentMatches", "descending"},
				{"goalsScore", "descending"},
			}},
			expectedExperiments[0],
		},
		{&ExperimentDesc{
			Dataset: dcsv.New(
				filepath.Join("..", "fixtures", "bank.csv"),
				true,
				rune(';'),
				fieldNames,
			),
			RuleFields: []string{"age", "job", "marital", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "p_1234567890outcome", "y",
			},
			RuleComplexity: rule.Complexity{Arithmetic: false},
			Aggregators: []*AggregatorDesc{
				{"num_married", "count", "marital == \"married\""},
				{"numSignedUp", "count", "y == \"yes\""},
				{"cost", "calc", "numMatches * 4.5"},
				{"income", "calc", "numSignedUp * 24"},
				{"profit", "calc", "income - cost"},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				{"profit", "descending"},
				{"numSignedUp", "descending"},
				{"cost", "ascending"},
				{"numMatches", "descending"},
				{"percentMatches", "descending"},
				{"goalsScore", "descending"},
			}},
			expectedExperiments[1],
		},
		{&ExperimentDesc{
			Dataset: dcsv.New(
				filepath.Join("..", "fixtures", "bank.csv"),
				true,
				rune(';'),
				fieldNames,
			),
			RuleFields: []string{"age", "job", "marital", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "p_1234567890outcome", "y",
			},
			RuleComplexity: rule.Complexity{Arithmetic: false},
			Aggregators: []*AggregatorDesc{
				{"num_married", "count", "marital == \"married\""},
				{"numSignedUp", "count", "y == \"yes\""},
				{"cost", "calc", "numMatches * 4.5"},
				{"income", "calc", "numSignedUp * 24"},
				{"profit", "calc", "income - cost"},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				{"profit", "descending"},
				{"numSignedUp", "descending"},
				{"cost", "ascending"},
				{"numMatches", "descending"},
				{"percentMatches", "descending"},
				{"goalsScore", "descending"},
			},
			Rules: []string{
				"job == \"manager\"",
				"age >= 27",
				"balance <= 1500",
			}},
			expectedExperiments[2],
		},
	}
	for i, c := range cases {
		got, err := New(c.experimentDesc)
		if err != nil {
			t.Errorf("(%d)  -New(%v) err: %s", i, c.experimentDesc, err)
		}
		if err := checkExperimentsMatch(got, c.want); err != nil {
			t.Errorf(
				"(%d) - New(%v)\n experiments don't match: %s\n got: %v\n want: %v",
				i, c.experimentDesc, err, got, c.want,
			)
		}
	}
}

func TestNew_errors(t *testing.T) {
	fieldNames := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y",
	}
	dataset := dcsv.New(
		filepath.Join("..", "fixtures", "bank.csv"),
		true,
		rune(';'),
		fieldNames,
	)

	cases := []struct {
		experimentDesc *ExperimentDesc
		wantErr        error
	}{
		{&ExperimentDesc{
			Dataset: dataset,
			RuleFields: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{
				{"numSignedUp", "count", "y == \"yes\""},
				{"cost", "calc", "numMatches * 4.5"},
				{"income", "calc", "numSignedUp * 24"},
				{"profit", "calc", "income - cost"},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				{"profit", "descending"},
				{"numSignedUp", "descending"},
				{"numMatches", "descending"},
				{"percentMatches", "descending"},
				{"age", "ascending"},
			}},
			InvalidSortFieldError("age"),
		},
		{&ExperimentDesc{
			Dataset: dataset,
			RuleFields: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{},
			Goals:       []string{"profit > 0"},
			SortOrder: []*SortDesc{
				{"numMatches", "Descending"},
			}},
			&InvalidSortDirectionError{"numMatches", "Descending"},
		},
		{&ExperimentDesc{
			Dataset: dataset,
			RuleFields: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{},
			Goals:       []string{"profit > 0"},
			SortOrder: []*SortDesc{
				{"percentMatches", "Ascending"},
			}},
			&InvalidSortDirectionError{"percentMatches", "Ascending"},
		},
		{&ExperimentDesc{
			Dataset:     dataset,
			RuleFields:  []string{},
			Aggregators: []*AggregatorDesc{},
			Goals:       []string{"profit > 0"},
			SortOrder: []*SortDesc{
				{"numMatches", "descending"},
				{"percentMatches", "descending"},
			}},
			ErrNoRuleFieldsSpecified,
		},
		{&ExperimentDesc{
			Dataset: dataset,
			RuleFields: []string{"age", "bob", "job", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{},
			Goals:       []string{"profit > 0"},
			SortOrder: []*SortDesc{
				{"numMatches", "descending"},
				{"percentMatches", "descending"},
			}},
			InvalidRuleFieldError("bob"),
		},
		{&ExperimentDesc{
			Dataset: dataset,
			RuleFields: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{
				{"pdays", "count", "day > 2"},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				{"numMatches", "descending"},
				{"percentMatches", "descending"},
			}},
			AggregatorNameClashError("pdays"),
		},
		{&ExperimentDesc{
			Dataset: dataset,
			RuleFields: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{
				{"numMatches", "count", "y == \"yes\""},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				{"numMatches", "descending"},
				{"percentMatches", "descending"},
			}},
			AggregatorNameReservedError("numMatches"),
		},
		{&ExperimentDesc{
			Dataset: dataset,
			RuleFields: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{
				{"percentMatches", "percent", "y == \"yes\""},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				{"numMatches", "descending"},
				{"percentMatches", "descending"},
			}},
			AggregatorNameReservedError("percentMatches"),
		},
		{&ExperimentDesc{
			Dataset: dataset,
			RuleFields: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{
				{"goalsScore", "count", "y == \"yes\""},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				{"numMatches", "descending"},
				{"percentMatches", "descending"},
			}},
			AggregatorNameReservedError("goalsScore"),
		},
		{&ExperimentDesc{
			Dataset: dataset,
			RuleFields: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{
				{"3numSignedUp", "count", "y == \"yes\""},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				{"numMatches", "descending"},
				{"percentMatches", "descending"},
			}},
			InvalidAggregatorNameError("3numSignedUp"),
		},
		{&ExperimentDesc{
			Dataset: dataset,
			RuleFields: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{
				{"num-signed-up", "count", "y == \"yes\""},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				{"numMatches", "descending"},
				{"percentMatches", "descending"},
			}},
			InvalidAggregatorNameError("num-signed-up"),
		},
		{&ExperimentDesc{
			Dataset: dataset,
			RuleFields: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{
				{"numSignedUp", "count", "y == \"yes\""},
			},
			Goals: []string{"profit > > 0"},
			SortOrder: []*SortDesc{
				{"numMatches", "descending"},
				{"percentMatches", "descending"},
			}},
			goal.InvalidGoalError("profit > > 0"),
		},
		{&ExperimentDesc{
			Dataset: dataset,
			RuleFields: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{
				{"numSignedUp", "count", "y == \"yes\""},
			},
			Goals: []string{},
			SortOrder: []*SortDesc{
				{"numMatches", "descending"},
				{"percentMatches", "descending"},
			},
			Rules: []string{
				"balance >= 27",
				"day > > 0",
				"month < 27",
			},
		},
			rule.InvalidExprError{Expr: "day > > 0"},
		},
	}
	for _, c := range cases {
		_, err := New(c.experimentDesc)
		if err == nil || c.wantErr.Error() != err.Error() {
			t.Errorf("New(%v) err: %v, wantErr: %v",
				c.experimentDesc, err, c.wantErr)
		}
	}
}

/***********************
   Helper functions
************************/

func checkExperimentsMatch(e1 *Experiment, e2 *Experiment) error {
	if !areStringArraysEqual(e1.RuleFields, e2.RuleFields) {
		return errors.New("RuleFields don't match")
	}
	if !areRuleComplexitiesEqual(e1.RuleComplexity, e2.RuleComplexity) {
		return errors.New("RuleComplexities don't match")
	}
	if !areGoalExpressionsEqual(e1.Goals, e2.Goals) {
		return errors.New("Goals don't match")
	}
	if !areAggregatorsEqual(e1.Aggregators, e2.Aggregators) {
		return errors.New("Aggregators don't match")
	}
	if !areSortOrdersEqual(e1.SortOrder, e2.SortOrder) {
		return errors.New("Sort Orders don't match")
	}
	if !areRulesEqual(e1.Rules, e2.Rules) {
		return errors.New("Rules don't match")
	}
	return checkDatasetsEqual(e1.Dataset, e2.Dataset)
}

func checkDatasetsEqual(ds1, ds2 ddataset.Dataset) error {
	conn1, err := ds1.Open()
	if err != nil {
		return err
	}
	conn2, err := ds2.Open()
	if err != nil {
		return err
	}
	for {
		conn1Next := conn1.Next()
		conn2Next := conn2.Next()
		if conn1Next != conn2Next {
			return errors.New("Datasets don't finish at same point")
		}
		if !conn1Next {
			break
		}

		conn1Record := conn1.Read()
		conn2Record := conn2.Read()
		if !reflect.DeepEqual(conn1Record, conn2Record) {
			return errors.New("Datasets don't match")
		}
	}
	if conn1.Err() != conn2.Err() {
		return errors.New("Datasets final error doesn't match")
	}
	return nil
}

func areStringArraysEqual(a1 []string, a2 []string) bool {
	if len(a1) != len(a2) {
		return false
	}
	for i, e := range a1 {
		if e != a2[i] {
			return false
		}
	}
	return true
}

func areRuleComplexitiesEqual(c1, c2 rule.Complexity) bool {
	return c1.Arithmetic == c2.Arithmetic
}

func areGoalExpressionsEqual(g1 []*goal.Goal, g2 []*goal.Goal) bool {
	if len(g1) != len(g2) {
		return false
	}
	for i, g := range g1 {
		if g.String() != g2[i].String() {
			return false
		}
	}
	return true
}

func areAggregatorsEqual(
	a1 []aggregators.AggregatorSpec,
	a2 []aggregators.AggregatorSpec,
) bool {
	if len(a1) != len(a2) {
		return false
	}
	for i, a := range a1 {
		if reflect.TypeOf(a) != reflect.TypeOf(a2[i]) ||
			a.Name() != a2[i].Name() ||
			a.Arg() != a2[i].Arg() {
			return false
		}
	}
	return true
}

func areSortOrdersEqual(so1 []SortField, so2 []SortField) bool {
	if len(so1) != len(so2) {
		return false
	}
	for i, sf1 := range so1 {
		sf2 := so2[i]
		if sf1.Field != sf2.Field || sf1.Direction != sf2.Direction {
			return false
		}
	}
	return true
}

func areRulesEqual(r1 []rule.Rule, r2 []rule.Rule) bool {
	if len(r1) != len(r2) {
		return false
	}
	for i, r := range r1 {
		if r.String() != r2[i].String() {
			return false
		}
	}
	return true

}
