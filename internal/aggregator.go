/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package internal

import "github.com/lawrencewoodman/dlit_go"

type Aggregator interface {
	CloneNew() Aggregator
	GetName() string
	GetArg() string
	GetResult([]Aggregator, int64) *dlit.Literal
	NextRecord(map[string]*dlit.Literal, bool) error
	IsEqual(Aggregator) bool
}

// TODO: Make the thisName optional
// TODO: Test this
func AggregatorsToMap(
	aggregators []Aggregator,
	numRecords int64,
	thisName string) map[string]*dlit.Literal {
	r := make(map[string]*dlit.Literal, len(aggregators))
	numRecordsL := dlit.MustNew(numRecords)
	r["numRecords"] = numRecordsL
	for _, aggregator := range aggregators {
		if thisName == aggregator.GetName() {
			break
		}
		r[aggregator.GetName()] = aggregator.GetResult(aggregators, numRecords)
	}
	return r
}
