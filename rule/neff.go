// Copyright (C) 2016-2017 vLife Systems Ltd <http://vlifesystems.com>
// Licensed under an MIT licence.  Please see LICENSE.md for details.

package rule

import (
	"github.com/lawrencewoodman/ddataset"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal"
)

// NEFF represents a rule determining if fieldA != fieldB
type NEFF struct {
	fieldA string
	fieldB string
}

func init() {
	registerGenerator("NEFF", generateNEFF)
}

func NewNEFF(fieldA, fieldB string) Rule {
	return &NEFF{fieldA: fieldA, fieldB: fieldB}
}

func (r *NEFF) String() string {
	return r.fieldA + " != " + r.fieldB
}

func (r *NEFF) IsTrue(record ddataset.Record) (bool, error) {
	lh, ok := record[r.fieldA]
	if !ok {
		return false, InvalidRuleError{Rule: r}
	}
	rh, ok := record[r.fieldB]
	if !ok {
		return false, InvalidRuleError{Rule: r}
	}

	lhInt, lhIsInt := lh.Int()
	rhInt, rhIsInt := rh.Int()
	if lhIsInt && rhIsInt {
		return lhInt != rhInt, nil
	}

	rhFloat, rhIsFloat := rh.Float()
	lhFloat, lhIsFloat := lh.Float()
	if lhIsFloat && rhIsFloat {
		return lhFloat != rhFloat, nil
	}

	// Don't compare bools as otherwise with the way that floats or ints
	// are cast to bools you would find that "True" == 1.0 because they would
	// both convert to true bools
	lhErr := lh.Err()
	rhErr := rh.Err()
	if lhErr != nil || rhErr != nil {
		return false, IncompatibleTypesRuleError{Rule: r}
	}

	return lh.String() != rh.String(), nil
}

func (r *NEFF) Fields() []string {
	return []string{r.fieldA, r.fieldB}
}

func generateNEFF(
	inputDescription *description.Description,
	generationDesc GenerationDescriber,
	field string,
) []Rule {
	fd := inputDescription.Fields[field]
	if fd.Kind != description.String && fd.Kind != description.Number {
		return []Rule{}
	}
	fieldNum := description.CalcFieldNum(inputDescription.Fields, field)
	rules := make([]Rule, 0)
	for oField, oFd := range inputDescription.Fields {
		if oFd.Kind == fd.Kind {
			oFieldNum := description.CalcFieldNum(inputDescription.Fields, oField)
			numSharedValues := calcNumSharedValues(fd, oFd)
			if fieldNum < oFieldNum &&
				numSharedValues >= 2 &&
				internal.IsStringInSlice(oField, generationDesc.Fields()) {
				r := NewNEFF(field, oField)
				rules = append(rules, r)
			}
		}
	}
	return rules
}
