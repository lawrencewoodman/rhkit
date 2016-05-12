/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package internal

import (
	"fmt"
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
)

type SumAggregator struct {
	name string
	sum  *dlit.Literal
	expr *dexpr.Expr
}

var sumExpr = dexpr.MustNew("sum+value")

func NewSumAggregator(name string, exprStr string) (*SumAggregator, error) {
	expr, err := dexpr.New(exprStr)
	if err != nil {
		return nil, err
	}
	ca := &SumAggregator{
		name: name,
		sum:  dlit.MustNew(0),
		expr: expr,
	}
	return ca, nil
}

func (a *SumAggregator) CloneNew() Aggregator {
	return &SumAggregator{
		name: a.name,
		sum:  dlit.MustNew(0),
		expr: a.expr,
	}
}

func (a *SumAggregator) GetName() string {
	return a.name
}

func (a *SumAggregator) GetArg() string {
	return a.expr.String()
}

func (a *SumAggregator) NextRecord(
	record map[string]*dlit.Literal,
	isRuleTrue bool,
) error {
	if isRuleTrue {
		exprValue := a.expr.Eval(record, CallFuncs)
		_, valueIsFloat := exprValue.Float()
		if !valueIsFloat {
			return fmt.Errorf("Value isn't a float: %s", exprValue)
		}

		vars := map[string]*dlit.Literal{
			"sum":   a.sum,
			"value": exprValue,
		}
		a.sum = sumExpr.Eval(vars, CallFuncs)
	}
	return nil
}

func (a *SumAggregator) GetResult(
	aggregators []Aggregator,
	goals []*Goal,
	numRecords int64,
) *dlit.Literal {
	return a.sum
}

func (a *SumAggregator) IsEqual(o Aggregator) bool {
	if _, ok := o.(*SumAggregator); !ok {
		return false
	}
	return a.name == o.GetName() && a.GetArg() == o.GetArg()
}
