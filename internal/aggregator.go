/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package internal

import (
	"errors"
	"github.com/lawrencewoodman/dlit"
)

type Aggregator interface {
	CloneNew() Aggregator
	GetName() string
	GetArg() string
	GetResult([]Aggregator, []*Goal, int64) *dlit.Literal
	NextRecord(map[string]*dlit.Literal, bool) error
	IsEqual(Aggregator) bool
}

// TODO: Make the thisName optional
// TODO: Test this
func AggregatorsToMap(
	aggregators []Aggregator,
	goals []*Goal,
	numRecords int64,
	thisName string,
) (map[string]*dlit.Literal, error) {
	r := make(map[string]*dlit.Literal, len(aggregators))
	numRecordsL := dlit.MustNew(numRecords)
	r["numRecords"] = numRecordsL
	for _, aggregator := range aggregators {
		if thisName == aggregator.GetName() {
			break
		}
		l := aggregator.GetResult(aggregators, goals, numRecords)
		if l.IsError() {
			return r, errors.New(l.String())
		}
		r[aggregator.GetName()] = l
	}
	return r, nil
}
