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
	"fmt"
	"github.com/lawrencewoodman/ddataset"
)

// OutsideFVI represents a rule determining if:
// field <= intValue || field >= intValue
type OutsideFVI struct {
	field string
	low   int64
	high  int64
}

func NewOutsideFVI(field string, low int64, high int64) (Rule, error) {
	if high <= low {
		return nil,
			fmt.Errorf("can't create Outside rule where high: %d <= low: %d",
				high, low)
	}
	return &OutsideFVI{field: field, low: low, high: high}, nil
}

func MustNewOutsideFVI(field string, low int64, high int64) Rule {
	r, err := NewOutsideFVI(field, low, high)
	if err != nil {
		panic(err)
	}
	return r
}

func (r *OutsideFVI) String() string {
	return fmt.Sprintf("%s <= %d || %s >= %d", r.field, r.low, r.field, r.high)
}

func (r *OutsideFVI) IsTrue(record ddataset.Record) (bool, error) {
	value, ok := record[r.field]
	if !ok {
		return false, InvalidRuleError{Rule: r}
	}

	valueInt, valueIsInt := value.Int()
	if valueIsInt {
		return valueInt <= r.low || valueInt >= r.high, nil
	}

	return false, IncompatibleTypesRuleError{Rule: r}
}

func (r *OutsideFVI) GetFields() []string {
	return []string{r.field}
}
