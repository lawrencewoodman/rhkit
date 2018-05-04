// Copyright (C) 2018 vLife Systems Ltd <http://vlifesystems.com>
// Licensed under an MIT licence.  Please see LICENSE.md for details.

package testhelpers

type GenerationDesc struct {
	DFields     []string
	DArithmetic bool
	DDeny       map[string][]string
}

func (gd GenerationDesc) Fields() []string {
	return gd.DFields
}

func (gd GenerationDesc) Arithmetic() bool {
	return gd.DArithmetic
}

func (gd GenerationDesc) Deny(generatorName string, field string) bool {
	fields := gd.DDeny[generatorName]
	for _, f := range fields {
		if f == field {
			return true
		}
	}
	return false
}
