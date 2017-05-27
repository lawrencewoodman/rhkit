/*
	Copyright (C) 2016-2017 vLife Systems Ltd <http://vlifesystems.com>
	This file is part of rhkit.

	rhkit is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	rhkit is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with rhkit; see the file COPYING.  If not, see
	<http://www.gnu.org/licenses/>.
*/

package aggregators

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/goal"
	"github.com/vlifesystems/rhkit/internal/dexprfuncs"
)

// mcc represents a Matthews correlation coefficient aggregator
// see: https://en.wikipedia.org/wiki/Matthews_correlation_coefficient
type mccAggregator struct{}

type mccSpec struct {
	name string
	expr *dexpr.Expr
}

type mccInstance struct {
	spec              *mccSpec
	numTruePositives  int64
	numTrueNegatives  int64
	numFalsePositives int64
	numFalseNegatives int64
}

// This is calculated using a dexpr because this easily handles errors and
// overflow/underflow errors
var mccExpr = dexpr.MustNew(
	"((tp*tn)-(fp*fn))/sqrt((tp+fp)*(tp+fn)*(tn+fp)*(tn+fn))",
	dexprfuncs.CallFuncs,
)

var radicandIsZeroExpr = dexpr.MustNew(
	"((tp+fp) * (tp+fn) * (tn+fp) * (tn+fn)) == 0",
	dexprfuncs.CallFuncs,
)

func init() {
	Register("mcc", &mccAggregator{})
}

func (a *mccAggregator) MakeSpec(
	name string,
	expr string,
) (AggregatorSpec, error) {
	dexpr, err := dexpr.New(expr, dexprfuncs.CallFuncs)
	if err != nil {
		return nil, err
	}
	d := &mccSpec{
		name: name,
		expr: dexpr,
	}
	return d, nil
}

func (ad *mccSpec) New() AggregatorInstance {
	return &mccInstance{
		spec:              ad,
		numTruePositives:  0,
		numTrueNegatives:  0,
		numFalsePositives: 0,
		numFalseNegatives: 0,
	}
}

func (ad *mccSpec) Name() string {
	return ad.name
}

func (ad *mccSpec) Kind() string {
	return "mcc"
}

func (ad *mccSpec) Arg() string {
	return ad.expr.String()
}

func (ai *mccInstance) Name() string {
	return ai.spec.name
}

func (ai *mccInstance) NextRecord(
	record map[string]*dlit.Literal,
	isRuleTrue bool,
) error {
	matchExprIsTrue, err := ai.spec.expr.EvalBool(record)
	if err != nil {
		return err
	}
	if matchExprIsTrue {
		if isRuleTrue {
			ai.numTruePositives++
		} else {
			ai.numFalseNegatives++
		}
	} else {
		if isRuleTrue {
			ai.numFalsePositives++
		} else {
			ai.numTrueNegatives++
		}
	}

	return nil
}

func (ai *mccInstance) Result(
	aggregatorInstances []AggregatorInstance,
	goals []*goal.Goal,
	numRecords int64,
) *dlit.Literal {
	if numRecords == 0 {
		return dlit.MustNew(0)
	}

	vars := map[string]*dlit.Literal{
		"tp": dlit.MustNew(ai.numTruePositives),
		"tn": dlit.MustNew(ai.numTrueNegatives),
		"fp": dlit.MustNew(ai.numFalsePositives),
		"fn": dlit.MustNew(ai.numFalseNegatives),
	}
	radIsZero, err := radicandIsZeroExpr.EvalBool(vars)
	if err != nil {
		return dlit.MustNew(err)
	}
	if radIsZero {
		return dlit.MustNew(0)
	}
	return mccExpr.Eval(vars)
}
