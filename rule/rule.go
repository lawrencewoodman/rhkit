/*
	Copyright (C) 2016-2017 vLife Systems Ltd <http://vlifesystems.com>
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

// Package rule implements rules to be tested against a dataset
package rule

import (
	"fmt"
	"github.com/lawrencewoodman/ddataset"
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/dexprfuncs"
	"sort"
	"strconv"
	"strings"
)

type Rule interface {
	fmt.Stringer
	IsTrue(record ddataset.Record) (bool, error)
	GetFields() []string
}

type Tweaker interface {
	Tweak(*description.Description, int) []Rule
}

type Overlapper interface {
	Overlaps(o Rule) bool
}

type InvalidRuleError struct {
	Rule Rule
}

type IncompatibleTypesRuleError struct {
	Rule Rule
}

func (e InvalidRuleError) Error() string {
	return "invalid rule: " + e.Rule.String()
}

func (e IncompatibleTypesRuleError) Error() string {
	return "incompatible types in rule: " + e.Rule.String()
}

// Sort sorts the rules in place using their .String() method
func Sort(rules []Rule) {
	sort.Sort(byString(rules))
}

// Uniq returns the slices of Rules with duplicates removed
func Uniq(rules []Rule) []Rule {
	results := []Rule{}
	mResults := map[string]interface{}{}
	for _, r := range rules {
		if _, ok := mResults[r.String()]; !ok {
			mResults[r.String()] = nil
			results = append(results, r)
		}
	}
	return results
}

func commaJoinValues(values []*dlit.Literal) string {
	str := fmt.Sprintf("\"%s\"", values[0].String())
	for _, v := range values[1:] {
		str += fmt.Sprintf(",\"%s\"", v)
	}
	return str
}

// byString implements sort.Interface for []Rule
type byString []Rule

func (rs byString) Len() int { return len(rs) }
func (rs byString) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}
func (rs byString) Less(i, j int) bool {
	return strings.Compare(rs[i].String(), rs[j].String()) == -1
}

// byNumber implements sort.Interface for []*dlit.Literal
type byNumber []*dlit.Literal

func (bn byNumber) Len() int { return len(bn) }
func (bn byNumber) Swap(i, j int) {
	bn[i], bn[j] = bn[j], bn[i]
}

var compareLitExpr = dexpr.MustNew("a < b", dexprfuncs.CallFuncs)

func (bn byNumber) Less(i, j int) bool {
	vars := map[string]*dlit.Literal{
		"a": bn[i],
		"b": bn[j],
	}
	if r, err := compareLitExpr.EvalBool(vars); err != nil {
		panic(err)
	} else {
		return r
	}
}

func truncateFloat(f float64, maxDP int) float64 {
	v := fmt.Sprintf("%.*f", maxDP, f)
	nf, _ := strconv.ParseFloat(v, 64)
	return nf
}

func generatePoints(
	value, min, max *dlit.Literal,
	maxDP int,
	stage int,
) []*dlit.Literal {
	points := make(map[string]*dlit.Literal)
	vars := map[string]*dlit.Literal{
		"min":   min,
		"max":   max,
		"maxDP": dlit.MustNew(maxDP),
		"stage": dlit.MustNew(stage),
		"value": value,
	}
	vars["step"] =
		dexpr.Eval("(max - min) / (10 * stage)", dexprfuncs.CallFuncs, vars)
	vars["low"] = dexpr.Eval("value - step", dexprfuncs.CallFuncs, vars)
	vars["high"] = dexpr.Eval("value + step", dexprfuncs.CallFuncs, vars)
	vars["interStep"] = dexpr.Eval("step/10", dexprfuncs.CallFuncs, vars)

	nextNExpr := dexpr.MustNew("n + interStep", dexprfuncs.CallFuncs)
	stopExpr := dexpr.MustNew("interStep <= 0 || n > high", dexprfuncs.CallFuncs)
	roundExpr := dexpr.MustNew("roundto(n, maxDP)", dexprfuncs.CallFuncs)
	isValidExpr := dexpr.MustNew(
		"newValue != value && newValue != low && newValue != high && "+
			"newValue > min && newValue < max",
		dexprfuncs.CallFuncs,
	)
	vars["n"] = dexpr.Eval("value - step", dexprfuncs.CallFuncs, vars)
	for {
		if stop, err := stopExpr.EvalBool(vars); stop || err != nil {
			break
		}
		v := roundExpr.Eval(vars)
		vars["newValue"] = v
		if ok, err := isValidExpr.EvalBool(vars); ok && err == nil {
			points[v.String()] = v
		}
		vars["n"] = nextNExpr.Eval(vars)
	}

	r := make([]*dlit.Literal, len(points))
	i := 0
	for _, p := range points {
		r[i] = p
		i++
	}

	sort.Sort(byNumber(r))
	return r
}
