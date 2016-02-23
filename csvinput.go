/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package main

import (
	"encoding/csv"
	"github.com/lawrencewoodman/dlit"
	"os"
)

type CsvInput struct {
	reader     *csv.Reader
	fieldNames []string
	filename   string
	separator  rune
}

func NewCsvInput(fieldNames []string, filename string,
	separator rune) (*CsvInput, error) {
	f, err := os.Open(filename)
	if err != nil {
		return &CsvInput{}, err
	}

	r := csv.NewReader(f)
	r.Comma = separator
	return &CsvInput{reader: r, fieldNames: fieldNames, filename: filename,
		separator: separator}, nil
}

func (c *CsvInput) Read() (map[string]*dlit.Literal, error) {
	recordLits := make(map[string]*dlit.Literal)
	record, err := c.reader.Read()
	if err != nil {
		return recordLits, err
	}
	for i, field := range record {
		l, err := dlit.New(field)
		if err != nil {
			return recordLits, err
		}
		recordLits[c.fieldNames[i]] = l
	}
	return recordLits, nil
}
