/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package internal

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

func (a *CalcAggregator) GetArg() string {
	return a.expr.String()
}

func (a *CalcAggregator) NextRecord(
	record map[string]*dlit.Literal, isRuleTrue bool) error {
	return nil
}

func (a *CalcAggregator) GetResult(
	aggregators []Aggregator,
	goals []*Goal,
	numRecords int64,
) *dlit.Literal {
	aggregatorsMap, err :=
		AggregatorsToMap(aggregators, goals, numRecords, a.name)
	if err != nil {
		return dlit.MustNew(err)
	}
	return a.expr.Eval(aggregatorsMap, CallFuncs)
}

func (a *CalcAggregator) IsEqual(o Aggregator) bool {
	if _, ok := o.(*CalcAggregator); !ok {
		return false
	}
	return a.name == o.GetName() && a.GetArg() == o.GetArg()
}
