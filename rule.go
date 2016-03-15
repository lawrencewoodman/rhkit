/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */

package main

import (
	"fmt"
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/dlit_go"
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

func (r *Rule) IsTrue(record map[string]*dlit.Literal) (bool, error) {
	isTrue, err := r.expr.EvalBool(record, callFuncs)
	// TODO: Create an error type for rule rather than coopting the dexpr one
	return isTrue, err
}

func (r *Rule) String() string {
	return r.expr.String()
}
