/*
	Copyright (C) 2016 vLife Systems Ltd <http://vlifesystems.com>
	This file is part of Rulehunter.

	Rulehunter is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	Rulehunter is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with Rulehunter; see the file COPYING.  If not, see
	<http://www.gnu.org/licenses/>.
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

// This should only be used for testing
func MustNewGoal(expr string) *Goal {
	g, err := NewGoal(expr)
	if err != nil {
		panic(fmt.Sprintf("Can't create goal: %s", err))
	}
	return g
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
