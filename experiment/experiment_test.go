package experiment

import (
	"errors"
	"github.com/lawrencewoodman/ddataset"
	"github.com/lawrencewoodman/ddataset/dcsv"
	"github.com/vlifesystems/rulehunter/aggregators"
	"github.com/vlifesystems/rulehunter/goal"
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
			ExcludeFieldNames: []string{"education"},
			Aggregators: []aggregators.Aggregator{
				// num_married to check for allowed characters
				aggregators.MustNew("num_married", "count", "marital == \"married\""),
				aggregators.MustNew("numSignedUp", "count", "y == \"yes\""),
				aggregators.MustNew("cost", "calc", "numMatches * 4.5"),
				aggregators.MustNew("income", "calc", "numSignedUp * 24"),
				aggregators.MustNew("profit", "calc", "income - cost")},
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
			ExcludeFields: []string{"education"},
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
			t.Errorf("New(%q) err: %s", c.experimentDesc, err)
		}
		if err := checkExperimentsMatch(got, c.want); err != nil {
			t.Errorf("New(%q)\n experiments don't match: %s\n got: %q\n want: %q",
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
			Title:         "This is a nice title",
			Dataset:       dataset,
			ExcludeFields: []string{},
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
			errors.New("Invalid sort field: age"),
		},
		{&ExperimentDesc{
			Title:         "This is a nice title",
			Dataset:       dataset,
			ExcludeFields: []string{},
			Aggregators:   []*AggregatorDesc{},
			Goals:         []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "Descending"},
			}},
			errors.New("Invalid sort direction: Descending, for field: numMatches"),
		},
		{&ExperimentDesc{
			Title:         "This is a nice title",
			Dataset:       dataset,
			ExcludeFields: []string{},
			Aggregators:   []*AggregatorDesc{},
			Goals:         []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"percentMatches", "Ascending"},
			}},
			errors.New("Invalid sort direction: Ascending, for field: percentMatches"),
		},
		{&ExperimentDesc{
			Title:         "This is a nice title",
			Dataset:       dataset,
			ExcludeFields: []string{"bob"},
			Aggregators:   []*AggregatorDesc{},
			Goals:         []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
			}},
			errors.New("Invalid exclude field: bob"),
		},
		{&ExperimentDesc{
			Title:         "This is a nice title",
			Dataset:       dataset,
			ExcludeFields: []string{},
			Aggregators: []*AggregatorDesc{
				&AggregatorDesc{"pdays", "count", "day > 2"},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
			}},
			errors.New("Aggregator name clashes with field name: pdays"),
		},
		{&ExperimentDesc{
			Title:         "This is a nice title",
			Dataset:       dataset,
			ExcludeFields: []string{},
			Aggregators: []*AggregatorDesc{
				&AggregatorDesc{"numMatches", "count", "y == \"yes\""},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
			}},
			errors.New("Aggregator name reserved: numMatches"),
		},
		{&ExperimentDesc{
			Title:         "This is a nice title",
			Dataset:       dataset,
			ExcludeFields: []string{},
			Aggregators: []*AggregatorDesc{
				&AggregatorDesc{"percentMatches", "percent", "y == \"yes\""},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
			}},
			errors.New("Aggregator name reserved: percentMatches"),
		},
		{&ExperimentDesc{
			Title:         "This is a nice title",
			Dataset:       dataset,
			ExcludeFields: []string{},
			Aggregators: []*AggregatorDesc{
				&AggregatorDesc{"goalsScore", "count", "y == \"yes\""},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
			}},
			errors.New("Aggregator name reserved: goalsScore"),
		},
		{&ExperimentDesc{
			Title:         "This is a nice title",
			Dataset:       dataset,
			ExcludeFields: []string{},
			Aggregators: []*AggregatorDesc{
				&AggregatorDesc{"3numSignedUp", "count", "y == \"yes\""},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
			}},
			errors.New("Invalid aggregator name: 3numSignedUp"),
		},
		{&ExperimentDesc{
			Title:         "This is a nice title",
			Dataset:       dataset,
			ExcludeFields: []string{},
			Aggregators: []*AggregatorDesc{
				&AggregatorDesc{"num-signed-up", "count", "y == \"yes\""},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
			}},
			errors.New("Invalid aggregator name: num-signed-up"),
		},
	}
	for _, c := range cases {
		_, err := New(c.experimentDesc)
		if err == nil || c.wantErr.Error() != err.Error() {
			t.Errorf("New(%q) err: %q, wantErr: %q",
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
	if !areStringArraysEqual(e1.ExcludeFieldNames, e2.ExcludeFieldNames) {
		return errors.New("ExcludeFieldNames don't match")
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
	a1 []aggregators.Aggregator,
	a2 []aggregators.Aggregator,
) bool {
	if len(a1) != len(a2) {
		return false
	}
	for i, e := range a1 {
		if !e.IsEqual(a2[i]) {
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
