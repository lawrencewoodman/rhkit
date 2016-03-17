/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */

package main

import (
	"errors"
	"fmt"
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/dlit_go"
	"regexp"
)

type Rule struct {
	expr *dexpr.Expr
}

type ErrInvalidRule string

func (e ErrInvalidRule) Error() string {
	return string(e)
}

func NewRule(exprStr string) (*Rule, error) {
	expr, err := dexpr.New(exprStr)
	if err != nil {
		return nil, ErrInvalidRule(fmt.Sprintf("Invalid rule: %s", exprStr))
	}
	return &Rule{expr}, nil
}

func (r *Rule) GetTweakableParts() (bool, string, string, string) {
	ruleStr := r.String()
	isTweakable := isTweakableRegexp.MatchString(ruleStr)
	if !isTweakable {
		return false, "", "", ""
	}
	fieldName := matchPartsRegexp.ReplaceAllString(ruleStr, "$1")
	operator := matchPartsRegexp.ReplaceAllString(ruleStr, "$2")
	value := matchPartsRegexp.ReplaceAllString(ruleStr, "$3")
	return isTweakable, fieldName, operator, value
}

func (r *Rule) IsTrue(record map[string]*dlit.Literal) (bool, error) {
	isTrue, err := r.expr.EvalBool(record, callFuncs)
	// TODO: Create an error type for rule rather than coopting the dexpr one
	return isTrue, err
}

func (r *Rule) String() string {
	return r.expr.String()
}

func (r *Rule) CloneWithValue(newValue string) (*Rule, error) {
	isTweakable, fieldName, operator, _ := r.GetTweakableParts()
	if !isTweakable {
		return nil, errors.New(fmt.Sprintf("Can't clone non-tweakable rule: %s", r))
	}
	newRule, err :=
		NewRule(fmt.Sprintf("%s %s %s", fieldName, operator, newValue))
	return newRule, err
}

var isTweakableRegexp = regexp.MustCompile("^[^( ]* (<|<=|>=|>) \\d+\\.?\\d*$")
var matchPartsRegexp = regexp.MustCompile("^([^( ]*) (<|<=|>=|>) (\\d+\\.?\\d*)$")
