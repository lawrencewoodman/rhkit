/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package main

import "github.com/lawrencewoodman/dlit"

type Aggregator interface {
	CloneNew() Aggregator
	GetName() string
	NextRecord(map[string]*dlit.Literal, bool) error
	GetResult([]Aggregator, int64) *dlit.Literal
}

// TODO: Make the thisName optional
// TODO: Test this
func AggregatorsToMap(
	aggregators []Aggregator,
	numRecords int64,
	thisName string) map[string]*dlit.Literal {
	r := make(map[string]*dlit.Literal, len(aggregators))
	numRecordsL, _ := dlit.New(numRecords)
	r["numRecords"] = numRecordsL
	for _, aggregator := range aggregators {
		if thisName == aggregator.GetName() {
			break
		}
		r[aggregator.GetName()] = aggregator.GetResult(aggregators, numRecords)
	}
	return r
}
