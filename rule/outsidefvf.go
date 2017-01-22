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

package rule

import (
	"errors"
	"github.com/lawrencewoodman/ddataset"
	"strconv"
)

// OutsideFVF represents a rule determining if:
// field <= intValue || field >= intValue
type OutsideFVF struct {
	field string
	low   float64
	high  float64
}

func NewOutsideFVF(field string, low float64, high float64) (Rule, error) {
	if high <= low {
		msg := "can't create Outside rule where high: " +
			strconv.FormatFloat(high, 'f', -1, 64) + " <= low: " +
			strconv.FormatFloat(low, 'f', -1, 64)
		return nil, errors.New(msg)
	}
	return &OutsideFVF{field: field, low: low, high: high}, nil
}

func (r *OutsideFVF) String() string {
	return r.field + " <= " + strconv.FormatFloat(r.low, 'f', -1, 64) +
		" || " + r.field + " >= " + strconv.FormatFloat(r.high, 'f', -1, 64)
}

func (r *OutsideFVF) IsTrue(record ddataset.Record) (bool, error) {
	value, ok := record[r.field]
	if !ok {
		return false, InvalidRuleError{Rule: r}
	}

	valueFloat, valueIsFloat := value.Float()
	if valueIsFloat {
		return valueFloat <= r.low || valueFloat >= r.high, nil
	}

	return false, IncompatibleTypesRuleError{Rule: r}
}

func (r *OutsideFVF) GetFields() []string {
	return []string{r.field}
}
