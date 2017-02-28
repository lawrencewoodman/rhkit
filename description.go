/*
	Copyright (C) 2016 vLife Systems Ltd <http://vlifesystems.com>
	This file is part of rhkit.

	rhkit is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	rhkit is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with rhkit; see the file COPYING.  If not, see
	<http://www.gnu.org/licenses/>.
*/

package rhkit

import (
	"encoding/json"
	"fmt"
	"github.com/lawrencewoodman/ddataset"
	"github.com/lawrencewoodman/dlit"
	"io/ioutil"
	"math"
	"os"
)

type Description struct {
	fields map[string]*fieldDescription
}

type fieldDescription struct {
	kind      fieldType
	min       *dlit.Literal
	max       *dlit.Literal
	maxDP     int
	values    map[string]valueDescription
	numValues int
}

type valueDescription struct {
	value *dlit.Literal
	num   int
}

func LoadDescriptionJSON(filename string) (*Description, error) {
	var dj descriptionJ

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	if err := dec.Decode(&dj); err != nil {
		return nil, err
	}

	fields := make(map[string]*fieldDescription, len(dj.Fields))
	for field, fd := range dj.Fields {
		values := make(map[string]valueDescription, len(fd.Values))
		for v, vd := range fd.Values {
			values[v] = valueDescription{
				value: dlit.NewString(vd.Value),
				num:   vd.Num,
			}
		}
		fields[field] = &fieldDescription{
			kind:      newFieldType(fd.Kind),
			min:       dlit.NewString(fd.Min),
			max:       dlit.NewString(fd.Max),
			maxDP:     fd.MaxDP,
			values:    values,
			numValues: fd.NumValues,
		}
	}
	d := &Description{fields: fields}
	return d, nil
}

func (d *Description) WriteJSON(filename string) error {
	fields := make(map[string]*fieldDescriptionJ, len(d.fields))
	for field, fd := range d.fields {
		fields[field] = newFieldDescriptionJ(fd)
	}
	dj := descriptionJ{Fields: fields}
	json, err := json.Marshal(dj)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, json, 0640)
}

// Create a New Description.
func newDescription() *Description {
	fd := map[string]*fieldDescription{}
	return &Description{fd}
}

type descriptionJ struct {
	Fields map[string]*fieldDescriptionJ
}

type fieldDescriptionJ struct {
	Kind      string
	Min       string
	Max       string
	MaxDP     int
	Values    map[string]valueDescriptionJ
	NumValues int
}

type valueDescriptionJ struct {
	Value string
	Num   int
}

func newFieldDescriptionJ(fd *fieldDescription) *fieldDescriptionJ {
	values := make(map[string]valueDescriptionJ, len(fd.values))
	for v, vd := range fd.values {
		values[v] = valueDescriptionJ{
			Value: vd.value.String(),
			Num:   vd.num,
		}
	}
	min := ""
	max := ""
	if fd.min != nil {
		min = fd.min.String()
	}
	if fd.max != nil {
		max = fd.max.String()
	}
	return &fieldDescriptionJ{
		Kind:      fd.kind.String(),
		Min:       min,
		Max:       max,
		MaxDP:     fd.maxDP,
		Values:    values,
		NumValues: fd.numValues,
	}
}

func (fd *fieldDescription) String() string {
	return fmt.Sprintf("Kind: %s, Min: %s, Max: %s, MaxDP: %d, Values: %s",
		fd.kind, fd.min, fd.max, fd.maxDP, fd.values)
}

// Analyse this record
func (d *Description) NextRecord(record ddataset.Record) {
	if len(d.fields) == 0 {
		for field, value := range record {
			d.fields[field] = &fieldDescription{
				kind:   ftUnknown,
				min:    value,
				max:    value,
				values: map[string]valueDescription{},
			}
		}
	}

	for field, value := range record {
		d.fields[field].processValue(value)
	}
}

func (f *fieldDescription) processValue(value *dlit.Literal) {
	f.updateKind(value)
	f.updateValues(value)
	f.updateNumBoundaries(value)
}

func (f *fieldDescription) updateKind(value *dlit.Literal) {
	switch f.kind {
	case ftUnknown:
		fallthrough
	case ftInt:
		if _, isInt := value.Int(); isInt {
			f.kind = ftInt
			break
		}
		fallthrough
	case ftFloat:
		if _, isFloat := value.Float(); isFloat {
			f.kind = ftFloat
			break
		}
		f.kind = ftString
	}
}

func (f *fieldDescription) updateValues(value *dlit.Literal) {
	// Chose 31 so could hold each day in month
	const maxNumValues = 31
	if f.kind == ftIgnore ||
		f.kind == ftUnknown ||
		f.numValues == -1 {
		return
	}
	if vd, ok := f.values[value.String()]; ok {
		f.values[value.String()] = valueDescription{vd.value, vd.num + 1}
		return
	}
	if f.numValues >= maxNumValues {
		if f.kind == ftString {
			f.kind = ftIgnore
		}
		f.values = map[string]valueDescription{}
		f.numValues = -1
		return
	}
	f.numValues++
	f.values[value.String()] = valueDescription{value, 1}
}

func (f *fieldDescription) updateNumBoundaries(value *dlit.Literal) {
	if f.kind == ftInt {
		valueInt, valueIsInt := value.Int()
		minInt, minIsInt := f.min.Int()
		maxInt, maxIsInt := f.max.Int()
		if !valueIsInt || !minIsInt || !maxIsInt {
			panic("Type mismatch")
		}
		f.min = dlit.MustNew(minI(minInt, valueInt))
		f.max = dlit.MustNew(maxI(maxInt, valueInt))
	} else if f.kind == ftFloat {
		valueFloat, valueIsFloat := value.Float()
		minFloat, minIsFloat := f.min.Float()
		maxFloat, maxIsFloat := f.max.Float()
		if !valueIsFloat || !minIsFloat || !maxIsFloat {
			panic("Type mismatch")
		}
		f.min = dlit.MustNew(math.Min(minFloat, valueFloat))
		f.max = dlit.MustNew(math.Max(maxFloat, valueFloat))
		f.maxDP =
			int(maxI(int64(f.maxDP), int64(numDecPlaces(value.String()))))
	}
}

func minI(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func maxI(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}
