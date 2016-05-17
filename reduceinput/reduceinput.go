/*
	Copyright (C) 2016 vLife Systems Ltd <http://vlifesystems.com>
	This file is part of Rulehunter.

	Rulehunter is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	Rulehunter is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with Rulehunter; see the file COPYING.  If not, see
	<http://www.gnu.org/licenses/>.
*/
package reduceinput

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/input"
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

func (r *ReduceInput) GetFieldNames() []string {
	return r.input.GetFieldNames()
}

func (r *ReduceInput) Close() error {
	return r.input.Close()
}
