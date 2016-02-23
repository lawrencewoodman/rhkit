/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package main

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
)

type CalcAggregator struct {
	name string
	expr *dexpr.Expr
}

func NewCalcAggregator(name string, expr string) (*CalcAggregator, error) {
	dexpr, err := dexpr.New(expr)
	if err != nil {
		return nil, err
	}
	ca := &CalcAggregator{name: name, expr: dexpr}
	return ca, nil
}

func (a *CalcAggregator) CloneNew() Aggregator {
	return &CalcAggregator{name: a.name, expr: a.expr}
}

func (a *CalcAggregator) GetName() string {
	return a.name
}

func (a *CalcAggregator) NextRecord(
	record map[string]*dlit.Literal, isRuleTrue bool) error {
	return nil
}

func (a *CalcAggregator) GetResult(
	aggregators []Aggregator, numRecords int64) *dlit.Literal {
	return a.expr.Eval(AggregatorsToMap(aggregators, numRecords, a.name),
		callFuncs)
}
