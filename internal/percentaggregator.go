/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package internal

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
)

type PercentAggregator struct {
	name       string
	numRecords int64
	numMatches int64
	expr       *dexpr.Expr
}

var percentExpr = dexpr.MustNew("roundto(100*numMatches/numRecords,2)")

func NewPercentAggregator(name string, expr string) (*PercentAggregator, error) {
	dexpr, err := dexpr.New(expr)
	if err != nil {
		return nil, err
	}
	ca :=
		&PercentAggregator{name: name, numMatches: 0, numRecords: 0, expr: dexpr}
	return ca, nil
}

func (a *PercentAggregator) CloneNew() Aggregator {
	return &PercentAggregator{
		name:       a.name,
		numMatches: 0,
		numRecords: 0,
		expr:       a.expr,
	}
}

func (a *PercentAggregator) GetName() string {
	return a.name
}

func (a *PercentAggregator) GetArg() string {
	return a.expr.String()
}

func (a *PercentAggregator) NextRecord(record map[string]*dlit.Literal,
	isRuleTrue bool) error {
	countExprIsTrue, err := a.expr.EvalBool(record, CallFuncs)
	if err != nil {
		return err
	}
	if isRuleTrue {
		a.numRecords++
		if countExprIsTrue {
			a.numMatches++
		}
	}
	return nil
}

func (a *PercentAggregator) GetResult(
	aggregators []Aggregator,
	goals []*Goal,
	numRecords int64,
) *dlit.Literal {
	if a.numRecords == 0 {
		return dlit.MustNew(0)
	}

	vars := map[string]*dlit.Literal{
		"numRecords": dlit.MustNew(a.numRecords),
		"numMatches": dlit.MustNew(a.numMatches),
	}
	return percentExpr.Eval(vars, CallFuncs)
}

func (a *PercentAggregator) IsEqual(o Aggregator) bool {
	if _, ok := o.(*PercentAggregator); !ok {
		return false
	}
	return a.name == o.GetName() && a.GetArg() == o.GetArg()
}
