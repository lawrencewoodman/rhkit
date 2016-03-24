/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */

package rulehunter

import (
	"errors"
	"fmt"
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/dlit_go"
	"github.com/lawrencewoodman/rulehunter/internal"
	"regexp"
)

type Rule struct {
	expr *dexpr.Expr
}

type ErrInvalidRule string

func (e ErrInvalidRule) Error() string {
	return string(e)
}

func (r *Rule) String() string {
	return r.expr.String()
}

func newRule(exprStr string) (*Rule, error) {
	expr, err := dexpr.New(exprStr)
	if err != nil {
		return nil, ErrInvalidRule(fmt.Sprintf("Invalid rule: %s", exprStr))
	}
	return &Rule{expr}, nil
}

func mustNewRule(exprStr string) *Rule {
	rule, err := newRule(exprStr)
	if err != nil {
		panic(err)
	}
	return rule
}

func (r *Rule) getTweakableParts() (bool, string, string, string) {
	ruleStr := r.String()
	isTweakable := isTweakableRegexp.MatchString(ruleStr)
	if !isTweakable {
		return false, "", "", ""
	}
	fieldName := matchTweakablePartsRegexp.ReplaceAllString(ruleStr, "$1")
	operator := matchTweakablePartsRegexp.ReplaceAllString(ruleStr, "$2")
	value := matchTweakablePartsRegexp.ReplaceAllString(ruleStr, "$3")
	return isTweakable, fieldName, operator, value
}

func (r *Rule) getInNiParts() (bool, string, string) {
	ruleStr := r.String()
	isInNi := isInNiRegexp.MatchString(ruleStr)
	if !isInNi {
		return false, "", ""
	}
	operator := matchInNiPartsRegexp.ReplaceAllString(ruleStr, "$1")
	fieldName := matchInNiPartsRegexp.ReplaceAllString(ruleStr, "$3")
	return isInNi, operator, fieldName
}

func (r *Rule) isTrue(record map[string]*dlit.Literal) (bool, error) {
	isTrue, err := r.expr.EvalBool(record, internal.CallFuncs)
	// TODO: Create an error type for rule rather than coopting the dexpr one
	return isTrue, err
}

func (r *Rule) cloneWithValue(newValue string) (*Rule, error) {
	isTweakable, fieldName, operator, _ := r.getTweakableParts()
	if !isTweakable {
		return nil, errors.New(fmt.Sprintf("Can't clone non-tweakable rule: %s", r))
	}
	newRule, err :=
		newRule(fmt.Sprintf("%s %s %s", fieldName, operator, newValue))
	return newRule, err
}

var isTweakableRegexp = regexp.MustCompile("^[^( ]* (<|<=|>=|>) \\d+\\.?\\d*$")
var matchTweakablePartsRegexp = regexp.MustCompile("^([^( ]*) (<|<=|>=|>) (\\d+\\.?\\d*)$")
var isInNiRegexp = regexp.MustCompile("^(in|ni)(\\()([^ ,]+)(.*\\))$")
var matchInNiPartsRegexp = regexp.MustCompile("^(in|ni)(\\()([^ ,]+)(.*\\))$")
