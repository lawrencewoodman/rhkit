/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package aggregators

import (
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/dlit_go"
	"github.com/lawrencewoodman/rulehunter/internal/dexprfuncs"
)

type Count struct {
	name       string
	numMatches int64
	expr       *dexpr.Expr
}

func NewCount(name string, expr string) (*Count, error) {
	dexpr, err := dexpr.New(expr)
	if err != nil {
		return nil, err
	}
	ca := &Count{name: name, numMatches: 0, expr: dexpr}
	return ca, nil
}

func (a *Count) CloneNew() Aggregator {
	return &Count{name: a.name, numMatches: 0, expr: a.expr}
}

func (a *Count) GetName() string {
	return a.name
}

func (a *Count) GetArg() string {
	return a.expr.String()
}

func (a *Count) NextRecord(record map[string]*dlit.Literal,
	isRuleTrue bool) error {
	countExprIsTrue, err := a.expr.EvalBool(record, dexprfuncs.CallFuncs)
	if err != nil {
		return err
	}
	if isRuleTrue && countExprIsTrue {
		a.numMatches++
	}
	return nil
}

func (a *Count) GetResult(
	aggregators []Aggregator, numRecords int64) *dlit.Literal {
	l := dlit.MustNew(a.numMatches)
	return l
}

func (a *Count) IsEqual(o Aggregator) bool {
	if _, ok := o.(*Count); !ok {
		return false
	}
	return a.name == o.GetName() && a.GetArg() == o.GetArg()
}
