package experiment

import (
	"errors"
	"github.com/lawrencewoodman/ddataset"
	"github.com/lawrencewoodman/ddataset/dcsv"
	"github.com/vlifesystems/rhkit/aggregators"
	"github.com/vlifesystems/rhkit/goal"
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
		&Experiment{},
		&Experiment{
			Title: "This is a jolly nice title",
			Dataset: dcsv.New(
				filepath.Join("..", "fixtures", "bank.csv"),
				true,
				rune(';'),
				fieldNames,
			),
			RuleFieldNames: []string{"age", "job", "marital", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "p_1234567890outcome", "y",
			},
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
				SortField{"profit", DESCENDING},
				SortField{"numSignedUp", DESCENDING},
				SortField{"cost", ASCENDING},
				SortField{"numMatches", DESCENDING},
				SortField{"percentMatches", DESCENDING},
				SortField{"goalsScore", DESCENDING},
			},
		},
	}
	cases := []struct {
		experimentDesc *ExperimentDesc
		want           *Experiment
	}{
		{&ExperimentDesc{
			Title: "This is a jolly nice title",
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
			Aggregators: []*AggregatorDesc{
				&AggregatorDesc{"num_married", "count", "marital == \"married\""},
				&AggregatorDesc{"numSignedUp", "count", "y == \"yes\""},
				&AggregatorDesc{"cost", "calc", "numMatches * 4.5"},
				&AggregatorDesc{"income", "calc", "numSignedUp * 24"},
				&AggregatorDesc{"profit", "calc", "income - cost"},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"profit", "descending"},
				&SortDesc{"numSignedUp", "descending"},
				&SortDesc{"cost", "ascending"},
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
				&SortDesc{"goalsScore", "descending"},
			}},
			expectedExperiments[1],
		},
	}
	for _, c := range cases {
		got, err := New(c.experimentDesc)
		if err != nil {
			t.Errorf("New(%v) err: %s", c.experimentDesc, err)
		}
		if err := checkExperimentsMatch(got, c.want); err != nil {
			t.Errorf("New(%v)\n experiments don't match: %s\n got: %v\n want: %v",
				c.experimentDesc, err, got, c.want)
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
			Title:   "This is a nice title",
			Dataset: dataset,
			RuleFields: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{
				&AggregatorDesc{"numSignedUp", "count", "y == \"yes\""},
				&AggregatorDesc{"cost", "calc", "numMatches * 4.5"},
				&AggregatorDesc{"income", "calc", "numSignedUp * 24"},
				&AggregatorDesc{"profit", "calc", "income - cost"},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"profit", "descending"},
				&SortDesc{"numSignedUp", "descending"},
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
				&SortDesc{"age", "ascending"},
			}},
			InvalidSortFieldError("age"),
		},
		{&ExperimentDesc{
			Title:   "This is a nice title",
			Dataset: dataset,
			RuleFields: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{},
			Goals:       []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "Descending"},
			}},
			&InvalidSortDirectionError{"numMatches", "Descending"},
		},
		{&ExperimentDesc{
			Title:   "This is a nice title",
			Dataset: dataset,
			RuleFields: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{},
			Goals:       []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"percentMatches", "Ascending"},
			}},
			&InvalidSortDirectionError{"percentMatches", "Ascending"},
		},
		{&ExperimentDesc{
			Title:       "This is a nice title",
			Dataset:     dataset,
			RuleFields:  []string{},
			Aggregators: []*AggregatorDesc{},
			Goals:       []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
			}},
			ErrNoRuleFieldsSpecified,
		},
		{&ExperimentDesc{
			Title:   "This is a nice title",
			Dataset: dataset,
			RuleFields: []string{"age", "bob", "job", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{},
			Goals:       []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
			}},
			InvalidRuleFieldError("bob"),
		},
		{&ExperimentDesc{
			Title:   "This is a nice title",
			Dataset: dataset,
			RuleFields: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{
				&AggregatorDesc{"pdays", "count", "day > 2"},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
			}},
			AggregatorNameClashError("pdays"),
		},
		{&ExperimentDesc{
			Title:   "This is a nice title",
			Dataset: dataset,
			RuleFields: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{
				&AggregatorDesc{"numMatches", "count", "y == \"yes\""},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
			}},
			AggregatorNameReservedError("numMatches"),
		},
		{&ExperimentDesc{
			Title:   "This is a nice title",
			Dataset: dataset,
			RuleFields: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{
				&AggregatorDesc{"percentMatches", "percent", "y == \"yes\""},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
			}},
			AggregatorNameReservedError("percentMatches"),
		},
		{&ExperimentDesc{
			Title:   "This is a nice title",
			Dataset: dataset,
			RuleFields: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{
				&AggregatorDesc{"goalsScore", "count", "y == \"yes\""},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
			}},
			AggregatorNameReservedError("goalsScore"),
		},
		{&ExperimentDesc{
			Title:   "This is a nice title",
			Dataset: dataset,
			RuleFields: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{
				&AggregatorDesc{"3numSignedUp", "count", "y == \"yes\""},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
			}},
			InvalidAggregatorNameError("3numSignedUp"),
		},
		{&ExperimentDesc{
			Title:   "This is a nice title",
			Dataset: dataset,
			RuleFields: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{
				&AggregatorDesc{"num-signed-up", "count", "y == \"yes\""},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
			}},
			InvalidAggregatorNameError("num-signed-up"),
		},
		{&ExperimentDesc{
			Title:   "This is a nice title",
			Dataset: dataset,
			RuleFields: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y",
			},
			Aggregators: []*AggregatorDesc{
				&AggregatorDesc{"numSignedUp", "count", "y == \"yes\""},
			},
			Goals: []string{"profit > > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
			}},
			goal.InvalidGoalError("profit > > 0"),
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
	if e1.Title != e2.Title {
		return errors.New("Titles don't match")
	}
	if !areStringArraysEqual(e1.RuleFieldNames, e2.RuleFieldNames) {
		return errors.New("RuleFieldNames don't match")
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
