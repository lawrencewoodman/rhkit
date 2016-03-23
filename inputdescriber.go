/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package main

import (
	"fmt"
	"github.com/lawrencewoodman/dlit_go"
	"github.com/lawrencewoodman/rulehunter/internal"
	"io"
	"math"
	"strings"
)

type FieldDescription struct {
	Kind      kind
	Min       *dlit.Literal
	Max       *dlit.Literal
	MaxDP     int
	Values    []*dlit.Literal
	numValues int
}

type kind int

const (
	UNKNOWN kind = iota
	IGNORE
	INT
	FLOAT
	STRING
)

func (fd *FieldDescription) String() string {
	return fmt.Sprintf("Kind: %s, Min: %s, Max: %s, MaxDP: %d, Values: %s",
		fd.Kind, fd.Min, fd.Max, fd.MaxDP, fd.Values)
}

func (k kind) String() string {
	switch k {
	case UNKNOWN:
		return "Unknown"
	case IGNORE:
		return "Ignore"
	case INT:
		return "Int"
	case FLOAT:
		return "Float"
	case STRING:
		return "String"
	}
	panic(fmt.Sprintf("Unsupported kind: %d", k))
}

func DescribeInput(input internal.Input) (map[string]*FieldDescription, error) {
	input.Rewind()
	fd := make(map[string]*FieldDescription)
	firstRecord := true
	for {
		record, err := input.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fd, err
		}
		if firstRecord {
			initFieldDescriptions(record, fd)
			firstRecord = false
		}
		updateFieldDescriptions(record, fd)
	}
	return fd, nil
}

func initFieldDescriptions(
	record map[string]*dlit.Literal, fd map[string]*FieldDescription) {
	for field, value := range record {
		fd[field] = &FieldDescription{Kind: UNKNOWN, Min: value, Max: value}
	}
}

func updateFieldDescriptions(
	record map[string]*dlit.Literal, fd map[string]*FieldDescription) {
	for field, value := range record {
		updateFieldDescription(value, fd[field])
	}
}

func updateFieldDescription(value *dlit.Literal, fd *FieldDescription) {
	updateFieldKind(value, fd)
	updateFieldValues(value, fd)
	updateFieldNumBoundaries(value, fd)
}

func updateFieldKind(value *dlit.Literal, fd *FieldDescription) {
	fdKind := fd.Kind
	switch fdKind {
	case UNKNOWN:
		fallthrough
	case INT:
		if _, isInt := value.Int(); isInt {
			fd.Kind = INT
			break
		}
		fallthrough
	case FLOAT:
		if _, isFloat := value.Float(); isFloat {
			fd.Kind = FLOAT
			break
		}
		fd.Kind = STRING
	}
}

func updateFieldValues(value *dlit.Literal, fd *FieldDescription) {
	// Chose 31 so could hold each day in month
	maxNumValues := 31
	if fd.Kind == IGNORE || fd.Kind == UNKNOWN || fd.numValues == -1 {
		return
	}
	for _, v := range fd.Values {
		if v.String() == value.String() {
			return
		}
	}
	if fd.numValues >= maxNumValues {
		if fd.Kind == STRING {
			fd.Kind = IGNORE
		}
		fd.Values = []*dlit.Literal{}
		fd.numValues = -1
		return
	}
	fd.Values = append(fd.Values, value)
	fd.numValues++
}

func updateFieldNumBoundaries(value *dlit.Literal, fd *FieldDescription) {
	if fd.Kind == INT {
		valueInt, valueIsInt := value.Int()
		minInt, minIsInt := fd.Min.Int()
		maxInt, maxIsInt := fd.Max.Int()
		if !valueIsInt || !minIsInt || !maxIsInt {
			panic("Type mismatch")
		}
		fd.Min = dlit.MustNew(minI(minInt, valueInt))
		fd.Max = dlit.MustNew(maxI(maxInt, valueInt))
	} else if fd.Kind == FLOAT {
		valueFloat, valueIsFloat := value.Float()
		minFloat, minIsFloat := fd.Min.Float()
		maxFloat, maxIsFloat := fd.Max.Float()
		if !valueIsFloat || !minIsFloat || !maxIsFloat {
			panic("Type mismatch")
		}
		fd.Min = dlit.MustNew(math.Min(minFloat, valueFloat))
		fd.Max = dlit.MustNew(math.Max(maxFloat, valueFloat))
		fd.MaxDP = int(maxI(int64(fd.MaxDP), int64(numDecPlaces(value.String()))))
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

func numDecPlaces(s string) int {
	i := strings.IndexByte(s, '.')
	if i > -1 {
		return len(s) - i - 1
	}
	return 0
}
