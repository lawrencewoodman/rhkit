/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package main

import (
	"fmt"
	"github.com/lawrencewoodman/dexpr_go"
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
	prettyPrintFieldDescriptions(fieldDescriptions)

	fmt.Printf("Generating rules...")
	rules, err := GenerateRules(fieldDescriptions, experiment.ExcludeFieldNames)
	if err != nil {
		panic(fmt.Sprintf("Couldn't make report: %s\n", err))
	}
	fmt.Printf("ok (%d generated)\n", len(rules))

	fmt.Printf("Assessing rules...")
	assessment, err :=
		AssessRules(rules, experiment.Aggregators, experiment.Goals, input)
	if err != nil {
		panic(fmt.Sprintf("Couldn't make report: %s\n", err))
	}
	fmt.Printf("ok\n")

	assessment.Sort(experiment.SortOrder)
	sortedRules := assessment.GetRules()

	fmt.Printf("Tweaking rules...")
	tweakableRules := TweakRules(sortedRules, fieldDescriptions)
	fmt.Printf("ok (%d generated)\n", len(tweakableRules))

	fmt.Printf("Assessing rules...")
	assessment2, err :=
		AssessRules(tweakableRules, experiment.Aggregators, experiment.Goals, input)
	if err != nil {
		panic(fmt.Sprintf("Couldn't make report: %s\n", err))
	}
	fmt.Printf("ok\n")

	assessment3, err := assessment.Merge(assessment2)
	if err != nil {
		panic(err)
	}
	assessment3.Sort(experiment.SortOrder)
	s, err := assessment3.ToJSON()
	if err != nil {
		panic(fmt.Sprintf("Couldn't make report: %s\n", err))
	} else {
		fmt.Println(s)
	}
}

func makeRule(expr string) *dexpr.Expr {
	r, err := dexpr.New(expr)
	if err != nil {
		panic("Can't make rule")
	}
	return r
}

func prettyPrintFieldDescriptions(fds map[string]*FieldDescription) {
	fmt.Println("Input Description\n")
	for field, fd := range fds {
		fmt.Println("--------------------------")
		fmt.Printf("%s\n--------------------------\n", field)
		prettyPrintFieldDescription(fd)
	}
	fmt.Println("\n")
}

func prettyPrintFieldDescription(fd *FieldDescription) {
	fmt.Printf("Kind: %s\n", fd.Kind)
	fmt.Printf("Min: %s\n", fd.Min)
	fmt.Printf("Max: %s\n", fd.Max)
	fmt.Printf("MaxDP: %d\n", fd.MaxDP)
	fmt.Printf("Values: %s\n", fd.Values)
}
