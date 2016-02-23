/*************************
 *  Test helper functions
 *************************/
package main

import (
	"fmt"
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
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

func mustNewLit(v interface{}) *dlit.Literal {
	l, err := dlit.New(v)
	if err != nil {
		panic(fmt.Sprintf("Can't create dlit.Literal: %q", err))
	}
	return l
}

func mustNewDExpr(expr string) *dexpr.Expr {
	dexpr, err := dexpr.New(expr)
	if err != nil {
		panic(fmt.Sprintf("Can't create dexpr.Expr: %q", err))
	}
	return dexpr
}

func mustNewCountAggregator(name string, expr string) *CountAggregator {
	c, err := NewCountAggregator(name, expr)
	if err != nil {
		panic(fmt.Sprintf("Can't create CountAggregator: %s", err))
	}
	return c
}

func mustNewCalcAggregator(name string, expr string) *CalcAggregator {
	c, err := NewCalcAggregator(name, expr)
	if err != nil {
		panic(fmt.Sprintf("Can't create CalcAggregator: %s", err))
	}
	return c
}
