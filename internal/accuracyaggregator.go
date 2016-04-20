/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package internal

import (
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/dlit_go"
)

type AccuracyAggregator struct {
	name       string
	numMatches int64
	expr       *dexpr.Expr
}

var accuracyExpr = dexpr.MustNew("roundto(100*numMatches/numRecords,2)")

func NewAccuracyAggregator(name string, expr string) (*AccuracyAggregator, error) {
	dexpr, err := dexpr.New(expr)
	if err != nil {
		return nil, err
	}
	ca :=
		&AccuracyAggregator{name: name, numMatches: 0, expr: dexpr}
	return ca, nil
}

func (a *AccuracyAggregator) CloneNew() Aggregator {
	return &AccuracyAggregator{
		name:       a.name,
		numMatches: 0,
		expr:       a.expr,
	}
}

func (a *AccuracyAggregator) GetName() string {
	return a.name
}

func (a *AccuracyAggregator) GetArg() string {
	return a.expr.String()
}

func (a *AccuracyAggregator) NextRecord(record map[string]*dlit.Literal,
	isRuleTrue bool) error {
	matchExprIsTrue, err := a.expr.EvalBool(record, CallFuncs)
	if err != nil {
		return err
	}
	if isRuleTrue == matchExprIsTrue {
		a.numMatches++
	}
	return nil
}

func (a *AccuracyAggregator) GetResult(
	aggregators []Aggregator,
	goals []*Goal,
	numRecords int64,
) *dlit.Literal {
	if numRecords == 0 {
		return dlit.MustNew(0)
	}

	vars := map[string]*dlit.Literal{
		"numRecords": dlit.MustNew(numRecords),
		"numMatches": dlit.MustNew(a.numMatches),
	}
	return accuracyExpr.Eval(vars, CallFuncs)
}

func (a *AccuracyAggregator) IsEqual(o Aggregator) bool {
	if _, ok := o.(*AccuracyAggregator); !ok {
		return false
	}
	return a.name == o.GetName() && a.GetArg() == o.GetArg()
}
