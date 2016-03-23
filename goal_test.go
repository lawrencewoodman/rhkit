package main

import (
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/dlit_go"
	"reflect"
	"testing"
)

func TestGoalsToMap(t *testing.T) {
	allAggregators := map[string]*dlit.Literal{
		"totalIncome":    dlit.MustNew(5000),
		"totalCost":      dlit.MustNew(4001),
		"percentMatches": dlit.MustNew(5.235),
	}
	goals := []*dexpr.Expr{
		mustNewDExpr("totalIncome > 4999"),
		mustNewDExpr("totalIncome > 5000"),
		mustNewDExpr("totalIncome + totalCost > 9000"),
		mustNewDExpr("totalIncome + totalCost > 9001"),
		mustNewDExpr("roundto(percentMatches,2) == 5.24"),
		mustNewDExpr("roundto(percentMatches,2) == 5.23"),
	}
	wantPasses := map[string]bool{
		"totalIncome > 4999":                true,
		"totalIncome > 5000":                false,
		"totalIncome + totalCost > 9000":    true,
		"totalIncome + totalCost > 9001":    false,
		"roundto(percentMatches,2) == 5.24": true,
		"roundto(percentMatches,2) == 5.23": false,
	}
	gotPasses, err := GoalsToMap(goals, allAggregators)
	if err != nil {
		t.Errorf("GoalsToMap(%q, %q) err: %s", goals, allAggregators, err)
	}
	if !reflect.DeepEqual(gotPasses, wantPasses) {
		t.Errorf("GoalsToMap(%q, %q) got: %s, want: %s",
			goals, allAggregators, gotPasses, wantPasses)
	}
}

func TestGoalsToMap_errors(t *testing.T) {
	allAggregators := map[string]*dlit.Literal{
		"totalIncome":    dlit.MustNew(5000),
		"totalCost":      dlit.MustNew(4001),
		"percentMatches": dlit.MustNew(5.235),
	}
	goals := []*dexpr.Expr{
		mustNewDExpr("totalIncome > 4999"),
		mustNewDExpr("totalIncome > 5000"),
		mustNewDExpr("totalIncome + totalCost > 9000"),
		mustNewDExpr("totalIncome + totalCost > 9001"),
		mustNewDExpr("roundbob(percentMatches,2) == 5.24"),
		mustNewDExpr("roundto(percentMatches,2) == 5.23"),
	}
	wantError := "Function doesn't exist: roundbob"
	_, err := GoalsToMap(goals, allAggregators)
	if err.Error() != wantError {
		t.Errorf("GoalsToMap(%q, %q) err: %s, wantError: %s",
			goals, allAggregators, err, wantError)
	}
}
