/*************************
 *  Test helper functions
 *************************/
package ruleassessment

import (
	"fmt"
	"github.com/vlifesystems/rulehunter/internal"
)

func errorMatch(e1 error, e2 error) bool {
	if e1 == nil && e2 == nil {
		return true
	}
	if e1 == nil || e2 == nil {
		return false
	}
	if e1.Error() == e2.Error() {
		return true
	}
	return false
}

func mustNewGoalsPassedScoreAggregator(
	name string,
) *internal.GoalsPassedScoreAggregator {
	a, err := internal.NewGoalsPassedScoreAggregator(name)
	if err != nil {
		panic(fmt.Sprintf("Can't create GoalsPassedScoreAggregator: %s", err))
	}
	return a
}
