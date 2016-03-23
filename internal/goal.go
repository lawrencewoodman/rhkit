/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package internal

import (
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/dlit_go"
)

// TODO: Create a Goal type

func GoalsToMap(
	goals []*dexpr.Expr,
	aggregators map[string]*dlit.Literal,
) (map[string]bool, error) {
	var err error
	r := make(map[string]bool, len(goals))

	for _, goal := range goals {
		r[goal.String()], err = goal.EvalBool(aggregators, CallFuncs)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}
