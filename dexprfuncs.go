/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package main

import (
	"errors"
	"fmt"
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"math"
)

var callFuncs = map[string]dexpr.CallFun{
	"roundto": roundTo,
}

// This uses round half-up to tie-break
func roundTo(args []*dlit.Literal) (*dlit.Literal, error) {
	if len(args) > 2 {
		err := errors.New("Too many arguments")
		r, _ := dlit.New(err)
		return r, err
	}
	x, isFloat := args[0].Float()
	if !isFloat {
		if args[0].IsError() {
			// TODO: Create a function in dlit to access error directly
			return args[0], errors.New(args[0].String())
		}
		err := errors.New(fmt.Sprintf("Can't convert to float: %s", args[0]))
		r, _ := dlit.New(err)
		return r, err
	}
	p, isInt := args[1].Int()
	if !isInt {
		if args[1].IsError() {
			// TODO: Create a function in dlit to access error directly
			return args[1], errors.New(args[1].String())
		}
		err := errors.New(fmt.Sprintf("Can't convert to int: %s", args[0]))
		r, _ := dlit.New(err)
		return r, err
	}
	shift := math.Pow(10, float64(p))
	r, err := dlit.New(math.Floor(.5+x*shift) / shift)
	return r, err
}
