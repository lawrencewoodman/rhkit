// Copyright (C) 2016-2017 vLife Systems Ltd <http://vlifesystems.com>
// Licensed under an MIT licence.  Please see LICENCE.md for details.

// Package goal handles goal expressions and their results
package goal

import (
	"fmt"
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/internal/dexprfuncs"
)

type Goal struct {
	expr *dexpr.Expr
}

type InvalidGoalError string

func (e InvalidGoalError) Error() string {
	return "invalid goal: " + string(e)
}

func New(exprStr string) (*Goal, error) {
	expr, err := dexpr.New(exprStr, dexprfuncs.CallFuncs)
	if err != nil {
		return nil, InvalidGoalError(exprStr)
	}
	return &Goal{expr}, nil
}

// This should only be used for testing
func MustNew(expr string) *Goal {
	g, err := New(expr)
	if err != nil {
		panic(fmt.Sprintf("can't create goal: %s", err))
	}
	return g
}

func (g *Goal) String() string {
	return g.expr.String()
}

func (g *Goal) Assess(aggregators map[string]*dlit.Literal) (bool, error) {
	passed, err := g.expr.EvalBool(aggregators)
	return passed, err
}
