/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	experimentFilename := os.Args[1]
	experiment, err := LoadExperiment(experimentFilename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Experiment: ", experiment)
}
