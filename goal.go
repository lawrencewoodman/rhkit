/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package main

import (
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/dlit_go"
	"github.com/lawrencewoodman/rulehunter/internal/aggregators"
	"github.com/lawrencewoodman/rulehunter/internal/dexprfuncs"
)

// TODO: Create a Goal type

// TODO: See if used anywhere
func HasGoalPassed(
	goal *dexpr.Expr,
	_aggregators []aggregators.Aggregator,
	numRecords int64,
) (bool, error) {
	results := aggregators.AggregatorsToMap(_aggregators, numRecords, "")

	isTrue, err := goal.EvalBool(results, dexprfuncs.CallFuncs)
	if err != nil {
		return false, err
	}
	return isTrue, nil
}

// TODO: test this
func GoalsToMap(
	goals []*dexpr.Expr,
	aggregators map[string]*dlit.Literal,
) (map[string]bool, error) {
	var err error
	r := make(map[string]bool, len(goals))

	for _, goal := range goals {
		r[goal.String()], err = goal.EvalBool(aggregators, dexprfuncs.CallFuncs)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}
