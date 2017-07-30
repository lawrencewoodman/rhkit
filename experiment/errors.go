// Copyright (C) 2016-2017 vLife Systems Ltd <http://vlifesystems.com>
// Licensed under an MIT licence.  Please see LICENSE.md for details.

package experiment

import "errors"

var ErrNoRuleFieldsSpecified = errors.New("no rule fields specified")

type InvalidRuleFieldError string

func (e InvalidRuleFieldError) Error() string {
	return "invalid rule field: " + string(e)
}
