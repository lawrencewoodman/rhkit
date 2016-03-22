/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package aggregators

import (
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/dlit_go"
	"github.com/lawrencewoodman/rulehunter/internal/dexprfuncs"
)

type Calc struct {
	name string
	expr *dexpr.Expr
}

func NewCalc(name string, expr string) (*Calc, error) {
	dexpr, err := dexpr.New(expr)
	if err != nil {
		return nil, err
	}
	ca := &Calc{name: name, expr: dexpr}
	return ca, nil
}

func (a *Calc) CloneNew() Aggregator {
	return &Calc{name: a.name, expr: a.expr}
}

func (a *Calc) GetName() string {
	return a.name
}

func (a *Calc) GetArg() string {
	return a.expr.String()
}

func (a *Calc) NextRecord(
	record map[string]*dlit.Literal, isRuleTrue bool) error {
	return nil
}

func (a *Calc) GetResult(
	aggregators []Aggregator, numRecords int64) *dlit.Literal {
	return a.expr.Eval(AggregatorsToMap(aggregators, numRecords, a.name),
		dexprfuncs.CallFuncs)
}

func (a *Calc) IsEqual(o Aggregator) bool {
	if _, ok := o.(*Calc); !ok {
		return false
	}
	return a.name == o.GetName() && a.GetArg() == o.GetArg()
}
