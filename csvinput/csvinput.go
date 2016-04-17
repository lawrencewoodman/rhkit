/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package csvinput

import (
	"encoding/csv"
	"errors"
	"github.com/lawrencewoodman/dlit_go"
	"github.com/lawrencewoodman/rulehunter/input"
	"os"
)

type CsvInput struct {
	file          *os.File
	reader        *csv.Reader
	fieldNames    []string
	filename      string
	separator     rune
	skipFirstLine bool
}

func New(fieldNames []string, filename string,
	separator rune, skipFirstLine bool) (input.Input, error) {
	f, r, err := makeCsvReader(filename, separator, skipFirstLine)
	if err != nil {
		return nil, err
	}
	r.Comma = separator
	return &CsvInput{file: f, reader: r, fieldNames: fieldNames,
		filename: filename, separator: separator,
		skipFirstLine: skipFirstLine}, nil
}

func (c *CsvInput) Clone() (input.Input, error) {
	newC, err :=
		New(c.fieldNames, c.filename, c.separator, c.skipFirstLine)
	return newC, err
}

func (c *CsvInput) Read() (map[string]*dlit.Literal, error) {
	recordLits := make(map[string]*dlit.Literal)
	record, err := c.reader.Read()
	if err != nil {
		return recordLits, err
	}
	if len(record) != len(c.fieldNames) {
		// TODO: Create specific error type for this
		return recordLits, errors.New("wrong number of field names for input")
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

func (c *CsvInput) Rewind() error {
	var err error
	if err = c.file.Close(); err != nil {
		return err
	}
	c.file, c.reader, err =
		makeCsvReader(c.filename, c.separator, c.skipFirstLine)
	if err != nil {
		return err
	}
	return nil
}

func makeCsvReader(filename string, separator rune,
	skipFirstLine bool) (*os.File, *csv.Reader, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	r := csv.NewReader(f)
	r.Comma = separator
	if skipFirstLine {
		_, err := r.Read()
		if err != nil {
			return nil, nil, err
		}
	}
	return f, r, err
}
