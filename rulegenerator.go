/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package main

import (
	"errors"
	"fmt"
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit_go"
	"regexp"
	"sort"
	"strings"
)

type ruleGeneratorFunc func(map[string]*FieldDescription,
	[]string, string) ([]*dexpr.Expr, error)

func GenerateRules(
	fieldDescriptions map[string]*FieldDescription,
	excludeFields []string) ([]*dexpr.Expr, error) {
	rules := make([]*dexpr.Expr, 0)
	ruleGenerators := []ruleGeneratorFunc{
		generateIntRules, generateFloatRules, generateStringRules,
		generateCompareNumericRules, generateCompareStringRules,
		generateInNiRules,
	}
	for field, _ := range fieldDescriptions {
		if !stringInSlice(field, excludeFields) {
			for _, ruleGenerator := range ruleGenerators {
				newRules, err := ruleGenerator(fieldDescriptions, excludeFields, field)
				if err != nil {
					return nil, err
				}
				rules = append(rules, newRules...)
			}
		}
	}
	return rules, nil
}

func stringInSlice(s string, strings []string) bool {
	for _, x := range strings {
		if x == s {
			return true
		}
	}
	return false
}

func generateIntRules(
	fieldDescriptions map[string]*FieldDescription,
	excludeFields []string, field string) ([]*dexpr.Expr, error) {
	fd := fieldDescriptions[field]
	if fd.Kind != INT {
		return []*dexpr.Expr{}, nil
	}
	rulesMap := make(map[string]*dexpr.Expr)
	min, _ := fd.Min.Int()
	max, _ := fd.Max.Int()
	values := fd.Values
	diff := max - min
	step := diff / 10
	if step == 0 {
		step = 1
	}
	// i set to 0 to make more tweakable
	for i := int64(0); i < diff; i += step {
		n := min + i
		exprStr := fmt.Sprintf("%s >= %d", field, n)
		rule, err := dexpr.New(exprStr)
		if err != nil {
			return nil, err
		}
		rulesMap[exprStr] = rule
	}

	for i := step; i <= diff; i += step {
		n := min + i
		exprStr := fmt.Sprintf("%s <= %d", field, n)
		rule, err := dexpr.New(exprStr)
		if err != nil {
			return nil, err
		}
		rulesMap[exprStr] = rule
	}

	if len(values) >= 2 {
		for _, v := range values {
			n, isInt := v.Int()
			if !isInt {
				return nil, errors.New(fmt.Sprintf("value isn't int: %s", v))
			}
			exprStr := fmt.Sprintf("%s == %d", field, n)
			rule, err := dexpr.New(exprStr)
			if err != nil {
				return nil, err
			}
			rulesMap[exprStr] = rule

			exprStr = fmt.Sprintf("%s != %d", field, n)
			rule, err = dexpr.New(exprStr)
			if err != nil {
				return nil, err
			}
			rulesMap[exprStr] = rule
		}
	}
	rules := rulesMapToArray(rulesMap)
	return rules, nil
}

var matchTrailingZerosRegexp = regexp.MustCompile("^(\\d+.[^ 0])(0+)$")
var matchAllZerosRegexp = regexp.MustCompile("^(0)(.0+)$")

func floatToString(f float64, maxDP int) string {
	fStr := fmt.Sprintf("%.*f", maxDP, f)
	fStr = matchTrailingZerosRegexp.ReplaceAllString(fStr, "$1")
	return matchAllZerosRegexp.ReplaceAllString(fStr, "0")
}

func generateFloatRules(
	fieldDescriptions map[string]*FieldDescription,
	excludeFields []string, field string) ([]*dexpr.Expr, error) {
	fd := fieldDescriptions[field]
	if fd.Kind != FLOAT {
		return []*dexpr.Expr{}, nil
	}
	rulesMap := make(map[string]*dexpr.Expr)
	min, _ := fd.Min.Float()
	max, _ := fd.Max.Float()
	maxDP := fd.MaxDP
	values := fd.Values
	diff := max - min
	step := diff / 10.0

	// i set to 0 to make more tweakable
	for i := float64(0); i < diff; i += step {
		n := min + i
		nStr := floatToString(n, maxDP)
		exprStr := fmt.Sprintf("%s >= %s", field, nStr)
		rule, err := dexpr.New(exprStr)
		if err != nil {
			return nil, err
		}
		rulesMap[exprStr] = rule
	}

	for i := step; i <= diff; i += step {
		n := min + i
		nStr := floatToString(n, maxDP)
		exprStr := fmt.Sprintf("%s <= %s", field, nStr)
		rule, err := dexpr.New(exprStr)
		if err != nil {
			return nil, err
		}
		rulesMap[exprStr] = rule
	}

	if len(values) >= 2 {
		for _, v := range values {
			n, isFloat := v.Float()
			if !isFloat {
				return nil, errors.New(fmt.Sprintf("value isn't float: %s", v))
			}
			nStr := floatToString(n, maxDP)
			exprStr := fmt.Sprintf("%s == %s", field, nStr)
			rule, err := dexpr.New(exprStr)
			if err != nil {
				return nil, err
			}
			rulesMap[exprStr] = rule

			exprStr = fmt.Sprintf("%s != %s", field, nStr)
			rule, err = dexpr.New(exprStr)
			if err != nil {
				return nil, err
			}
			rulesMap[exprStr] = rule
		}
	}
	rules := rulesMapToArray(rulesMap)
	return rules, nil
}

func generateCompareNumericRules(
	fieldDescriptions map[string]*FieldDescription,
	excludeFields []string, field string) ([]*dexpr.Expr, error) {
	fd := fieldDescriptions[field]
	if fd.Kind != INT && fd.Kind != FLOAT {
		return []*dexpr.Expr{}, nil
	}
	fieldNum := calcFieldNum(fieldDescriptions, field)
	rulesMap := make(map[string]*dexpr.Expr)
	exprFmts := []string{
		"%s < %s", "%s <= %s", "%s == %s", "%s != %s", "%s >= %s", "%s > %s",
	}

	for oField, oFd := range fieldDescriptions {
		oFieldNum := calcFieldNum(fieldDescriptions, oField)
		isComparable := hasComparableNumberRange(fd, oFd)
		if fieldNum < oFieldNum &&
			isComparable &&
			!stringInSlice(oField, excludeFields) {
			for _, exprFmt := range exprFmts {
				exprStr := fmt.Sprintf(exprFmt, field, oField)
				rule, err := dexpr.New(exprStr)
				if err != nil {
					return nil, err
				}
				rulesMap[exprStr] = rule
			}
		}
	}
	rules := rulesMapToArray(rulesMap)
	return rules, nil
}

func generateCompareStringRules(
	fieldDescriptions map[string]*FieldDescription,
	excludeFields []string, field string) ([]*dexpr.Expr, error) {
	fd := fieldDescriptions[field]
	if fd.Kind != STRING {
		return []*dexpr.Expr{}, nil
	}
	fieldNum := calcFieldNum(fieldDescriptions, field)
	rulesMap := make(map[string]*dexpr.Expr)
	exprFmts := []string{
		"%s == %s", "%s != %s",
	}
	for oField, oFd := range fieldDescriptions {
		if oFd.Kind == STRING {
			oFieldNum := calcFieldNum(fieldDescriptions, oField)
			numSharedValues := calcNumSharedValues(fd, oFd)
			if fieldNum < oFieldNum &&
				numSharedValues >= 2 &&
				!stringInSlice(oField, excludeFields) {
				for _, exprFmt := range exprFmts {
					exprStr := fmt.Sprintf(exprFmt, field, oField)
					rule, err := dexpr.New(exprStr)
					if err != nil {
						return nil, err
					}
					rulesMap[exprStr] = rule
				}
			}
		}
	}
	rules := rulesMapToArray(rulesMap)
	return rules, nil
}

func calcNumSharedValues(fd1 *FieldDescription, fd2 *FieldDescription) int {
	numShared := 0
	for _, v1 := range fd1.Values {
		for _, v2 := range fd2.Values {
			if v1.String() == v2.String() {
				numShared++
			}
		}
	}
	return numShared
}

func generateStringRules(
	fieldDescriptions map[string]*FieldDescription,
	excludeFields []string, field string) ([]*dexpr.Expr, error) {
	fd := fieldDescriptions[field]
	if fd.Kind != STRING {
		return []*dexpr.Expr{}, nil
	}
	rulesMap := make(map[string]*dexpr.Expr)

	for _, v := range fd.Values {
		s := v.String()
		exprStr := fmt.Sprintf("%s == \"%s\"", field, s)
		rule, err := dexpr.New(exprStr)
		if err != nil {
			return nil, err
		}
		rulesMap[exprStr] = rule
		if len(fd.Values) > 2 {
			exprStr := fmt.Sprintf("%s != \"%s\"", field, s)
			rule, err := dexpr.New(exprStr)
			if err != nil {
				return nil, err
			}
			rulesMap[exprStr] = rule
		}
	}
	rules := rulesMapToArray(rulesMap)
	return rules, nil
}

func isNumberField(fd *FieldDescription) bool {
	return fd.Kind == INT || fd.Kind == FLOAT
}

func hasComparableNumberRange(
	fd1 *FieldDescription, fd2 *FieldDescription) bool {
	if !isNumberField(fd1) || !isNumberField(fd2) {
		return false
	}
	var isComparable bool
	vars := map[string]*dlit.Literal{
		"min1": fd1.Min,
		"max1": fd1.Max,
		"min2": fd2.Min,
		"max2": fd2.Max,
	}
	funcs := map[string]dexpr.CallFun{}
	expr, err := dexpr.New("min1 < max2 && max1 > min2")
	if err != nil {
		return false
	}
	isComparable, err = expr.EvalBool(vars, funcs)
	if err != nil {
		return false
	}
	return isComparable
}

func rulesMapToArray(rulesMap map[string]*dexpr.Expr) []*dexpr.Expr {
	rules := make([]*dexpr.Expr, len(rulesMap))
	i := 0
	for _, expr := range rulesMap {
		rules[i] = expr
		i++
	}
	return rules
}

func generateInNiRules(
	fieldDescriptions map[string]*FieldDescription,
	excludeFields []string, field string) ([]*dexpr.Expr, error) {
	fd := fieldDescriptions[field]
	numValues := len(fd.Values)
	if fd.Kind != STRING && fd.Kind != FLOAT && fd.Kind != INT ||
		numValues <= 3 || numValues > 12 {
		return []*dexpr.Expr{}, nil
	}
	rulesMap := make(map[string]*dexpr.Expr)
	exprFmts := []string{
		"in(%s,%s)", "ni(%s,%s)",
	}
	for i := 3; ; i++ {
		numOnBits := calcNumOnBits(i)
		if numOnBits >= numValues {
			break
		}
		if numOnBits >= 2 && numOnBits <= 5 && numOnBits < (numValues-1) {
			compareValuesStr := makeCompareValuesStr(fd.Values, i)
			for _, exprFmt := range exprFmts {
				exprStr := fmt.Sprintf(exprFmt, field, compareValuesStr)
				rule, err := dexpr.New(exprStr)
				if err != nil {
					return nil, err
				}
				rulesMap[exprStr] = rule
			}

		}
	}
	rules := rulesMapToArray(rulesMap)
	return rules, nil
}

func makeCompareValues(values []*dlit.Literal, i int) []*dlit.Literal {
	bStr := fmt.Sprintf("%b", i)
	numOnBits := calcNumOnBits(i)
	numValues := len(values)
	j := numValues - 1
	compareValues := make([]*dlit.Literal, numOnBits)
	k := 0
	for _, b := range reverseString(bStr) {
		if b == '1' {
			compareValues[k] = values[numValues-1-j]
			k++
		}
		j -= 1
	}
	return compareValues
}

func reverseString(s string) (r string) {
	for _, v := range s {
		r = string(v) + r
	}
	return
}

func makeCompareValuesStr(values []*dlit.Literal, i int) string {
	compareValues := makeCompareValues(values, i)
	str := fmt.Sprintf("\"%s\"", compareValues[0].String())
	for _, v := range compareValues[1:] {
		str += fmt.Sprintf(",\"%s\"", v)
	}
	return str
}

func calcNumOnBits(i int) int {
	bStr := fmt.Sprintf("%b", i)
	return strings.Count(bStr, "1")
}

func calcFieldNum(
	fieldDescriptions map[string]*FieldDescription, fieldN string) int {
	fields := make([]string, len(fieldDescriptions))
	i := 0
	for field, _ := range fieldDescriptions {
		fields[i] = field
		i++
	}
	sort.Strings(fields)
	j := 0
	for _, field := range fields {
		if field == fieldN {
			return j
		}
		j++
	}
	panic("Can't field field in fieldDescriptions")
}
