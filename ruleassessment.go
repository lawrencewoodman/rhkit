/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package rulehunter

import (
	"errors"
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/dlit_go"
	"github.com/lawrencewoodman/rulehunter/internal"
)

type ruleAssessment struct {
	Rule        *Rule
	Aggregators []internal.Aggregator
	Goals       []*dexpr.Expr
}

// Note: This clones the aggregators to ensure the results are specific
//       to the rule.
func newRuleAssessment(
	rule *Rule,
	aggregators []internal.Aggregator,
	goals []*dexpr.Expr,
) *ruleAssessment {
	cloneAggregators := make([]internal.Aggregator, len(aggregators))
	for i, a := range aggregators {
		cloneAggregators[i] = a.CloneNew()
	}
	return &ruleAssessment{Rule: rule, Aggregators: cloneAggregators,
		Goals: goals}
}

func (ra *ruleAssessment) nextRecord(record map[string]*dlit.Literal) error {
	var ruleIsTrue bool
	var err error
	for _, aggregator := range ra.Aggregators {
		ruleIsTrue, err = ra.Rule.isTrue(record)
		if err != nil {
			return err
		}
		err = aggregator.NextRecord(record, ruleIsTrue)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ra *ruleAssessment) getAggregatorValue(
	name string, numRecords int64) (*dlit.Literal, bool) {
	for _, aggregator := range ra.Aggregators {
		if aggregator.GetName() == name {
			return aggregator.GetResult(ra.Aggregators, numRecords), true
		}
	}
	// TODO: Test and create specific error type
	err := errors.New("Aggregator doesn't exist")
	l := dlit.MustNew(err)
	return l, false
}
