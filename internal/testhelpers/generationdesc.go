// Copyright (C) 2017 vLife Systems Ltd <http://vlifesystems.com>
// Licensed under an MIT licence.  Please see LICENSE.md for details.

package testhelpers

type GenerationDesc struct {
	DFields     []string
	DArithmetic bool
}

func (gd GenerationDesc) Fields() []string {
	return gd.DFields
}

func (gd GenerationDesc) Arithmetic() bool {
	return gd.DArithmetic
}
