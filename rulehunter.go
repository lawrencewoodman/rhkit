/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package main

import (
	"encoding/json"
	"fmt"
	"github.com/lawrencewoodman/dexpr"
	"log"
	"os"
)

func main() {
	var experiment *Experiment
	var input Input
	var err error
	experimentFilename := os.Args[1]
	experiment, err = LoadExperiment(experimentFilename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Experiment: %s\n\n", experiment)
	// TODO: Create input from LoadExperiment
	input, err = NewCsvInput(experiment.FieldNames, experiment.InputFilename,
		experiment.Separator, experiment.IsFirstLineFieldNames)
	if err != nil {
		log.Fatal(err)
	}

	fieldDescriptions, err := DescribeInput(input)
	if err != nil {
		panic(fmt.Sprintf("Couldn't describe input: %s", err))
	}

	fmt.Printf("Input Description: %q\n\n", fieldDescriptions)

	rules := []*dexpr.Expr{
		makeRule("age > 20"),
		makeRule("age > 30"),
		makeRule("age > 40"),
		makeRule("marital == \"married\""),
	}
	report, err :=
		AssessRules(rules, experiment.Aggregators, experiment.Goals, input)
	if err != nil {
		fmt.Printf("Couldn't make report: %s\n", err)
	} else {
		b, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			fmt.Printf("Couldn't make report json: %s\n", err)
		} else {
			os.Stdout.Write(b)
		}
	}
}

func makeRule(expr string) *dexpr.Expr {
	r, err := dexpr.New(expr)
	if err != nil {
		panic("Can't make rule")
	}
	return r
}
