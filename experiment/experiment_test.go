package experiment

import (
	"errors"
	"fmt"
	"github.com/vlifesystems/rulehunter/aggregators"
	"github.com/vlifesystems/rulehunter/csvinput"
	"github.com/vlifesystems/rulehunter/goal"
	"github.com/vlifesystems/rulehunter/input"
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
			Input: mustNewCsvInput(
				fieldNames,
				filepath.Join("..", "fixtures", "bank.csv"),
				rune(';'),
				true,
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
			Input: mustNewCsvInput(
				fieldNames,
				filepath.Join("..", "fixtures", "bank.csv"),
				rune(';'),
				true,
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
		experimentsMatch, reason := experimentMatch(got, c.want)
		if !experimentsMatch {
			t.Errorf("New(%q)\n Reason: %s\n got: %q\n want: %q",
				c.experimentDesc, reason, got, c.want)
		}
	}
}

func TestNew_errors(t *testing.T) {
	fieldNames := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y",
	}
	input := mustNewCsvInput(
		fieldNames,
		filepath.Join("..", "fixtures", "bank.csv"),
		rune(';'),
		true,
	)

	cases := []struct {
		experimentDesc *ExperimentDesc
		wantErr        error
	}{
		{&ExperimentDesc{
			Title:         "This is a nice title",
			Input:         input,
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
			Input:         input,
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
			Input:         input,
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
			Input:         input,
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
			Input:         input,
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
			Input:         input,
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
			Input:         input,
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
			Input:         input,
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
			Input:         input,
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
			Input:         input,
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

func TestClose(t *testing.T) {
	fieldNames := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y",
	}
	input := mustNewCsvInput(
		fieldNames,
		filepath.Join("..", "fixtures", "bank.csv"),
		rune(';'),
		true,
	)

	experimentDesc := &ExperimentDesc{
		Title:         "This is a nice title",
		Input:         input,
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
			&SortDesc{"cost", "ascending"},
		},
	}
	experiment, err := New(experimentDesc)
	if err != nil {
		t.Errorf("New(%q) err: %q", experimentDesc, err)
	}
	if !experiment.Input.Next() {
		t.Errorf("Next() return false on first call")
	}
	if err := experiment.Close(); err != nil {
		t.Errorf("Close() err: %s", err)
	}
	if experiment.Input.Next() {
		t.Errorf("Next() return true on second call")
	}
}

/***********************
   Helper functions
************************/

func experimentMatch(e1 *Experiment, e2 *Experiment) (bool, string) {
	if e1.Title != e2.Title {
		return false, fmt.Sprintf("Titles don't match",
			e1.Title, e2.Title)
	}
	if !areStringArraysEqual(e1.ExcludeFieldNames, e2.ExcludeFieldNames) {
		return false, "ExcludeFieldNames don't match"
	}
	if !areGoalExpressionsEqual(e1.Goals, e2.Goals) {
		return false, "Goals don't match"
	}
	if !areAggregatorsEqual(e1.Aggregators, e2.Aggregators) {
		return false, "Aggregators don't match"
	}
	if !areSortOrdersEqual(e1.SortOrder, e2.SortOrder) {
		return false, "Sort Orders don't match"
	}
	inputsEqual, msg := areInputsEqual(e1.Input, e2.Input)
	return inputsEqual, msg
}

func areInputsEqual(i1, i2 input.Input) (bool, string) {
	for {
		i1Next := i1.Next()
		i2Next := i2.Next()
		if i1Next != i2Next {
			return false, "Inputs don't finish at same point"
		}
		if !i1Next {
			break
		}

		i1Record, i1Err := i1.Read()
		i2Record, i2Err := i2.Read()
		if i1Err != i2Err {
			return false, "Inputs don't error at same point"
		} else if i1Err == nil && i2Err == nil {
			if !reflect.DeepEqual(i1Record, i2Record) {
				return false, "Inputs don't match"
			}
		}
	}
	if i1.Err() != i2.Err() {
		return false, "Inputs final error doesn't match"
	}
	return true, ""
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

func mustNewCsvInput(
	fieldNames []string,
	filename string,
	separator rune,
	skipFirstLine bool,
) input.Input {
	input, err :=
		csvinput.New(fieldNames, filename, separator, skipFirstLine)
	if err != nil {
		panic(fmt.Sprintf("Couldn't create Csv Input: %s", err))
	}
	return input
}
