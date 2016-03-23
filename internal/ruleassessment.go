/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package internal

import (
	"errors"
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/dlit_go"
)

type RuleAssessment struct {
	Rule        *Rule
	Aggregators []Aggregator
	Goals       []*dexpr.Expr
}

// Note: This clones the aggregators to ensure the results are specific
//       to the rule.
func NewRuleAssessment(
	rule *Rule,
	aggregators []Aggregator,
	goals []*dexpr.Expr,
) *RuleAssessment {
	cloneAggregators := make([]Aggregator, len(aggregators))
	for i, a := range aggregators {
		cloneAggregators[i] = a.CloneNew()
	}
	return &RuleAssessment{Rule: rule, Aggregators: cloneAggregators,
		Goals: goals}
}

func (ra *RuleAssessment) NextRecord(record map[string]*dlit.Literal) error {
	var ruleIsTrue bool
	var err error
	for _, aggregator := range ra.Aggregators {
		ruleIsTrue, err = ra.Rule.IsTrue(record)
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

func (ra *RuleAssessment) GetAggregatorValue(
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
