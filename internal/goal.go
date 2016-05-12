/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package internal

import (
	"fmt"
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
)

type Goal struct {
	expr       *dexpr.Expr
	passed     bool
	calculated bool
}

type ErrInvalidGoal string

func (e ErrInvalidGoal) Error() string {
	return string(e)
}

func NewGoal(exprStr string) (*Goal, error) {
	expr, err := dexpr.New(exprStr)
	if err != nil {
		return nil, ErrInvalidGoal(fmt.Sprintf("Invalid goal: %s", exprStr))
	}
	return &Goal{expr, false, false}, nil
}

func (g *Goal) String() string {
	return g.expr.String()
}

func (g *Goal) Clone() *Goal {
	return &Goal{g.expr, false, false}
}

// This is to make other functions easier to test
func (g *Goal) SetPassed(b bool) {
	g.passed = b
	g.calculated = true
}

func (g *Goal) Assess(aggregators map[string]*dlit.Literal) (bool, error) {
	if g.calculated {
		return g.passed, nil
	}
	passed, err := g.expr.EvalBool(aggregators, CallFuncs)
	if err != nil {
		g.passed = passed
		g.calculated = true
	}
	return passed, err
}

func (g *Goal) HasPassed() bool {
	if !g.calculated {
		panic(fmt.Sprintf("Goal hasn't been assessed: %s", g))
	}
	return g.passed
}
