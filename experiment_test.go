package rulehunter

import (
	"errors"
	"fmt"
	"github.com/lawrencewoodman/rulehunter/csvinput"
	"github.com/lawrencewoodman/rulehunter/input"
	"github.com/lawrencewoodman/rulehunter/internal"
	"io"
	"path/filepath"
	"reflect"
	"testing"
)

func TestMakeExperiment(t *testing.T) {
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
				filepath.Join("fixtures", "bank.csv"),
				rune(';'),
				true,
			),
			FieldNames:        fieldNames,
			ExcludeFieldNames: []string{"education"},
			Aggregators: []internal.Aggregator{
				// num_married to check for allowed characters
				mustNewCountAggregator("num_married", "marital == \"married\""),
				mustNewCountAggregator("numSignedUp", "y == \"yes\""),
				mustNewCalcAggregator("cost", "numMatches * 4.5"),
				mustNewCalcAggregator("income", "numSignedUp * 24"),
				mustNewCalcAggregator("profit", "income - cost")},
			Goals: []*internal.Goal{mustNewGoal("profit > 0")},
			SortOrder: []SortField{
				SortField{"profit", DESCENDING},
				SortField{"numSignedUp", DESCENDING},
				SortField{"cost", ASCENDING},
				SortField{"numMatches", DESCENDING},
				SortField{"percentMatches", DESCENDING},
				SortField{"numGoalsPassed", DESCENDING},
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
				filepath.Join("fixtures", "bank.csv"),
				rune(';'),
				true,
			),
			Fields:        fieldNames,
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
				&SortDesc{"numGoalsPassed", "descending"},
			}},
			expectedExperiments[1],
		},
	}
	for _, c := range cases {
		got, err := MakeExperiment(c.experimentDesc)
		if err != nil {
			t.Errorf("MakeExperiment(%q) err: %s", c.experimentDesc, err)
		}
		experimentsMatch, reason := experimentMatch(got, c.want)
		if !experimentsMatch {
			t.Errorf("MakeExperiment(%q)\n Reason: %s\n got: %q\n want: %q",
				c.experimentDesc, reason, got, c.want)
		}
	}
}

func TestMakeExperiment_errors(t *testing.T) {
	fieldNames := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "poutcome", "y",
	}
	input := mustNewCsvInput(
		fieldNames,
		filepath.Join("fixtures", "bank.csv"),
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
			Fields:        []string{"age"},
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
			}},
			errors.New("Must specify at least two field names"),
		},
		{&ExperimentDesc{
			Title:         "This is a nice title",
			Input:         input,
			Fields:        fieldNames,
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
			Fields:        fieldNames,
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
			Fields:        fieldNames,
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
			Fields:        fieldNames,
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
			Fields:        []string{"month", "y", "numSignedUp"},
			ExcludeFields: []string{},
			Aggregators: []*AggregatorDesc{
				&AggregatorDesc{"numSignedUp", "count", "y == \"yes\""},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
			}},
			errors.New("Aggregator name clashes with field name: numSignedUp"),
		},
		{&ExperimentDesc{
			Title:         "This is a nice title",
			Input:         input,
			Fields:        []string{"month", "y", "numSignedUp"},
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
			Fields:        []string{"month", "y", "numSignedUp"},
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
			Fields:        []string{"month", "y", "numSignedUp"},
			ExcludeFields: []string{},
			Aggregators: []*AggregatorDesc{
				&AggregatorDesc{"numGoalsPassed", "count", "y == \"yes\""},
			},
			Goals: []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
			}},
			errors.New("Aggregator name reserved: numGoalsPassed"),
		},
		{&ExperimentDesc{
			Title:         "This is a nice title",
			Input:         input,
			Fields:        []string{"month", "y", "numSignedUp"},
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
			Fields:        []string{"month", "y", "numSignedUp"},
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
		{&ExperimentDesc{
			Title:         "This is a nice title",
			Input:         input,
			Fields:        []string{"month", "profit£", "numSignedUp"},
			ExcludeFields: []string{},
			Aggregators:   []*AggregatorDesc{},
			Goals:         []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
			}},
			errors.New("Invalid field name: profit£"),
		},
		{&ExperimentDesc{
			Title:         "This is a nice title",
			Input:         input,
			Fields:        []string{"month", "profit", "num-said-yes"},
			ExcludeFields: []string{},
			Aggregators:   []*AggregatorDesc{},
			Goals:         []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
			}},
			errors.New("Invalid field name: num-said-yes"),
		},
		{&ExperimentDesc{
			Title:         "This is a nice title",
			Input:         input,
			Fields:        []string{"month", "profit", "2pay"},
			ExcludeFields: []string{},
			Aggregators:   []*AggregatorDesc{},
			Goals:         []string{"profit > 0"},
			SortOrder: []*SortDesc{
				&SortDesc{"numMatches", "descending"},
				&SortDesc{"percentMatches", "descending"},
			}},
			errors.New("Invalid field name: 2pay"),
		},
	}
	for _, c := range cases {
		_, err := MakeExperiment(c.experimentDesc)
		if err == nil || c.wantErr.Error() != err.Error() {
			t.Errorf("MakeExperiment(%q) err: %q, wantErr: %q",
				c.experimentDesc, err, c.wantErr)
		}
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
	if !areStringArraysEqual(e1.FieldNames, e2.FieldNames) {
		return false, "FieldNames don't match"
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
	for {
		e1Record, e1Err := e1.Input.Read()
		e2Record, e2Err := e2.Input.Read()
		if e1Err != e2Err {
			return false, "Inputs don't error at same point"
		} else if e1Err == nil && e2Err == nil {
			if !reflect.DeepEqual(e1Record, e2Record) {
				return false, "Inputs don't match"
			}
		} else if e1Err == io.EOF || e2Err == io.EOF {
			break
		}
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

func areGoalExpressionsEqual(g1 []*internal.Goal, g2 []*internal.Goal) bool {
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
	a1 []internal.Aggregator,
	a2 []internal.Aggregator,
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
