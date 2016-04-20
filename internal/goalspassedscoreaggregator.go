/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package internal

import (
	"github.com/lawrencewoodman/dlit_go"
)

type GoalsPassedScoreAggregator struct {
	name string
}

func NewGoalsPassedScoreAggregator(
	name string,
) (*GoalsPassedScoreAggregator, error) {
	a := &GoalsPassedScoreAggregator{name: name}
	return a, nil
}

func (a *GoalsPassedScoreAggregator) CloneNew() Aggregator {
	return &GoalsPassedScoreAggregator{name: a.name}
}

func (a *GoalsPassedScoreAggregator) GetName() string {
	return a.name
}

func (a *GoalsPassedScoreAggregator) GetArg() string {
	return ""
}

func (a *GoalsPassedScoreAggregator) NextRecord(
	record map[string]*dlit.Literal,
	isRuleTrue bool,
) error {
	return nil
}

func (a *GoalsPassedScoreAggregator) GetResult(
	aggregators []Aggregator,
	goals []*Goal,
	numRecords int64,
) *dlit.Literal {
	aggregatorsMap, err :=
		AggregatorsToMap(aggregators, goals, numRecords, a.name)
	if err != nil {
		return dlit.MustNew(err)
	}
	numGoalsPassed := 0.0
	increment := 1.0
	for _, goal := range goals {
		hasPassed, err := goal.Assess(aggregatorsMap)
		if err != nil {
			return dlit.MustNew(err)
		}

		if hasPassed {
			numGoalsPassed += increment
		} else {
			increment = 0.001
		}
	}
	return dlit.MustNew(numGoalsPassed)
}

func (a *GoalsPassedScoreAggregator) IsEqual(o Aggregator) bool {
	if _, ok := o.(*GoalsPassedScoreAggregator); !ok {
		return false
	}
	return a.name == o.GetName()
}
