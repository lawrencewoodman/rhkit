/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package reduceinput

import (
	"github.com/lawrencewoodman/dlit_go"
	"github.com/lawrencewoodman/rulehunter/input"
	"io"
)

type ReduceInput struct {
	input      input.Input
	recordNum  int
	numRecords int
	err        error
}

func New(input input.Input, numRecords int) (input.Input, error) {
	return &ReduceInput{
		input:      input,
		recordNum:  -1,
		numRecords: numRecords,
		err:        nil,
	}, nil
}

func (r *ReduceInput) Clone() (input.Input, error) {
	i, err := r.input.Clone()
	return i, err
}

func (r *ReduceInput) Next() bool {
	if r.Err() != nil {
		return false
	}
	if r.recordNum < r.numRecords {
		r.recordNum++
		return r.input.Next()
	}
	r.err = io.EOF
	return false
}

func (r *ReduceInput) Err() error {
	if r.err == io.EOF {
		return nil
	}
	return r.input.Err()
}

func (r *ReduceInput) Read() (map[string]*dlit.Literal, error) {
	record, err := r.input.Read()
	return record, err
}

func (r *ReduceInput) Rewind() error {
	if r.Err() != nil {
		return r.Err()
	}
	r.recordNum = -1
	r.err = r.input.Rewind()
	return r.err
}

// This should only be called by Experiment.Close() ordinarily
func (r *ReduceInput) Close() error {
	return r.input.Close()
}
