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

package rule

import (
	"github.com/lawrencewoodman/ddataset"
	"strconv"
)

// EQFVF represents a rule determening if fieldA == floatValue
type EQFVF struct {
	field string
	value float64
}

func NewEQFVF(field string, value float64) Rule {
	return &EQFVF{field: field, value: value}
}

func (r *EQFVF) String() string {
	return r.field + " == " + strconv.FormatFloat(r.value, 'f', -1, 64)
}

func (r *EQFVF) GetInNiParts() (bool, string, string) {
	return false, "", ""
}

func (r *EQFVF) IsTrue(record ddataset.Record) (bool, error) {
	lh, ok := record[r.field]
	if !ok {
		return false, InvalidRuleError(r.String())
	}

	lhFloat, lhIsFloat := lh.Float()
	if lhIsFloat {
		return lhFloat == r.value, nil
	}

	// TODO: Return error saying incompatible types
	return false, InvalidRuleError(r.String())
}
