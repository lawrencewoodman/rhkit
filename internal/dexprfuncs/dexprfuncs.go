/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package dexprfuncs

import (
	"errors"
	"fmt"
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/dlit_go"
	"math"
)

var CallFuncs = map[string]dexpr.CallFun{
	"roundto": roundTo,
	"in":      in,
	"ni":      ni,
	"true":    alwaysTrue,
}

// This uses round half-up to tie-break
func roundTo(args []*dlit.Literal) (*dlit.Literal, error) {
	if len(args) > 2 {
		err := errors.New("Too many arguments")
		r := dlit.MustNew(err)
		return r, err
	}
	x, isFloat := args[0].Float()
	if !isFloat {
		if args[0].IsError() {
			// TODO: Create a function in dlit to access error directly
			return args[0], errors.New(args[0].String())
		}
		err := errors.New(fmt.Sprintf("Can't convert to float: %s", args[0]))
		r := dlit.MustNew(err)
		return r, err
	}
	p, isInt := args[1].Int()
	if !isInt {
		if args[1].IsError() {
			// TODO: Create a function in dlit to access error directly
			return args[1], errors.New(args[1].String())
		}
		err := errors.New(fmt.Sprintf("Can't convert to int: %s", args[0]))
		r := dlit.MustNew(err)
		return r, err
	}
	shift := math.Pow(10, float64(p))
	r, err := dlit.New(math.Floor(.5+x*shift) / shift)
	return r, err
}

// Is a string IN a list of strings
func in(args []*dlit.Literal) (*dlit.Literal, error) {
	if len(args) < 2 {
		err := errors.New("Too few arguments")
		r := dlit.MustNew(err)
		return r, err
	}
	needle := args[0]
	haystack := args[1:]
	for _, v := range haystack {
		if needle.String() == v.String() {
			r, err := dlit.New(true)
			return r, err
		}
	}
	r, err := dlit.New(false)
	return r, err
}

// Is a string NI a list of strings
func ni(args []*dlit.Literal) (*dlit.Literal, error) {
	if len(args) < 2 {
		err := errors.New("Too few arguments")
		r := dlit.MustNew(err)
		return r, err
	}
	needle := args[0]
	haystack := args[1:]
	for _, v := range haystack {
		if needle.String() == v.String() {
			r, err := dlit.New(false)
			return r, err
		}
	}
	r, err := dlit.New(true)
	return r, err
}

// Returns true
func alwaysTrue(args []*dlit.Literal) (*dlit.Literal, error) {
	return dlit.MustNew(true), nil
}
