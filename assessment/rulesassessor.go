// Copyright (C) 2016-2017 vLife Systems Ltd <http://vlifesystems.com>
// Licensed under an MIT licence.  Please see LICENSE.md for details.

package assessment

import (
	"github.com/lawrencewoodman/ddataset"
	"github.com/vlifesystems/rhkit/aggregator"
	"github.com/vlifesystems/rhkit/goal"
	"github.com/vlifesystems/rhkit/rule"
)

// AssessRules tests the rules for a Dataset against the aggregators and goals
// supplied and returns an Assessment along with any errors
func AssessRules(
	dataset ddataset.Dataset,
	rules []rule.Rule,
	aggregatorSpecs []aggregator.Spec,
	goals []*goal.Goal,
) (*Assessment, error) {
	ruleAssessments := make([]*RuleAssessment, len(rules))
	for i, rule := range rules {
		ruleAssessments[i] = newRuleAssessment(rule, aggregatorSpecs, goals)
	}
	numRecords, err := processDataset(dataset, ruleAssessments)
	if err != nil {
		return &Assessment{}, err
	}
	assessment := New(numRecords)
	err = assessment.addRuleAssessments(ruleAssessments)
	return assessment, err
}

func processDataset(
	dataset ddataset.Dataset,
	ruleAssessments []*RuleAssessment,
) (int64, error) {
	numRecords := int64(0)
	conn, err := dataset.Open()
	if err != nil {
		return numRecords, err
	}
	defer conn.Close()

	for conn.Next() {
		record := conn.Read()
		numRecords++
		for _, ruleAssessment := range ruleAssessments {
			err := ruleAssessment.NextRecord(record)
			if err != nil {
				return numRecords, err
			}
		}
	}

	return numRecords, conn.Err()
}
