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
	"github.com/lawrencewoodman/dexpr"
	"github.com/vlifesystems/rulehunter/internal/dexprfuncs"
	"regexp"
)

type DRule struct {
	expr *dexpr.Expr
}

func NewDRule(exprStr string) (Rule, error) {
	expr, err := dexpr.New(exprStr)
	if err != nil {
		return nil, InvalidRuleError(exprStr)
	}
	return &DRule{expr: expr}, nil
}

func MustNewDRule(exprStr string) Rule {
	rule, err := NewDRule(exprStr)
	if err != nil {
		panic(err)
	}
	return rule
}

func (r *DRule) String() string {
	return r.expr.String()
}

func (r *DRule) GetInNiParts() (bool, string, string) {
	ruleStr := r.String()
	isInNi := isInNiRegexp.MatchString(ruleStr)
	if !isInNi {
		return false, "", ""
	}
	operator := matchInNiPartsRegexp.ReplaceAllString(ruleStr, "$1")
	fieldName := matchInNiPartsRegexp.ReplaceAllString(ruleStr, "$3")
	return isInNi, operator, fieldName
}

func (r *DRule) IsTrue(record ddataset.Record) (bool, error) {
	isTrue, err := r.expr.EvalBool(record, dexprfuncs.CallFuncs)
	// TODO: Create an error type for rule rather than coopting the dexpr one
	return isTrue, err
}

var isInNiRegexp = regexp.MustCompile("^(in|ni)(\\()([^ ,]+)(.*\\))$")
var matchInNiPartsRegexp = regexp.MustCompile("^(in|ni)(\\()([^ ,]+)(.*\\))$")
