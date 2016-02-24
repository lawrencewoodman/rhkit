/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package main

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
)

type CountAggregator struct {
	name       string
	numMatches int64
	expr       *dexpr.Expr
}

func NewCountAggregator(name string, expr string) (*CountAggregator, error) {
	dexpr, err := dexpr.New(expr)
	if err != nil {
		return nil, err
	}
	ca := &CountAggregator{name: name, numMatches: 0, expr: dexpr}
	return ca, nil
}

func (a *CountAggregator) CloneNew() Aggregator {
	return &CountAggregator{name: a.name, numMatches: 0, expr: a.expr}
}

func (a *CountAggregator) GetName() string {
	return a.name
}

func (a *CountAggregator) GetArg() string {
	return a.expr.String()
}

func (a *CountAggregator) NextRecord(record map[string]*dlit.Literal,
	isRuleTrue bool) error {
	countExprIsTrue, err := a.expr.EvalBool(record, callFuncs)
	if err != nil {
		return err
	}
	if isRuleTrue && countExprIsTrue {
		a.numMatches++
	}
	return nil
}

func (a *CountAggregator) GetResult(
	aggregators []Aggregator, numRecords int64) *dlit.Literal {
	l, _ := dlit.New(a.numMatches)
	return l
}

func (a *CountAggregator) IsEqual(o Aggregator) bool {
	if _, ok := o.(*CountAggregator); !ok {
		return false
	}
	return a.name == o.GetName() && a.GetArg() == o.GetArg()
}
