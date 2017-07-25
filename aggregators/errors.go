// Copyright (C) 2016-2017 vLife Systems Ltd <http://vlifesystems.com>
// Licensed under an MIT licence.  Please see LICENSE.md for details.

package aggregators

// NameReservedError indicates that an aggregator name can't be
// used because it is a reserved name
type NameReservedError string

func (e NameReservedError) Error() string {
	return "aggregator name reserved: " + string(e)
}

// NameClashError indicates that an aggregator name can't be used
// because it clashes with a field name
type NameClashError string

func (e NameClashError) Error() string {
	return "aggregator name clashes with field name: " + string(e)
}

// InvalidNameError indicates that the aggregator name is invalid
type InvalidNameError string

func (e InvalidNameError) Error() string {
	return "invalid aggregator name: " + string(e)
}
