// Copyright (C) 2016-2017 vLife Systems Ltd <http://vlifesystems.com>
// Licensed under an MIT licence.  Please see LICENSE.md for details.

// Package testhelpers is used to help test rhkit
package testhelpers

import (
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"os"
	"path/filepath"
	"testing"
)

func MakeStringsDlitSlice(strings ...string) []*dlit.Literal {
	r := make([]*dlit.Literal, len(strings))
	for i, s := range strings {
		r[i] = dlit.NewString(s)
	}
	return r
}

func CheckErrorMatch(t *testing.T, context string, got, want error) {
	if got == nil && want == nil {
		return
	}
	if got == nil || want == nil {
		t.Errorf("%s got err: %s, want : %s", context, got, want)
	}
	if perr, ok := want.(*os.PathError); ok {
		if err := checkPathErrorMatch(got, perr); err != nil {
			t.Errorf("%s %s", context, err)
		}
	}
	if got.Error() != want.Error() {
		t.Errorf("%s got err: %s, want : %s", context, got, want)
	}
}

func checkPathErrorMatch(checkErr error, wantErr *os.PathError) error {
	perr, ok := checkErr.(*os.PathError)
	if !ok {
		return fmt.Errorf("got err type: %T, want error type: os.PathError",
			checkErr)
	}
	if perr.Op != wantErr.Op {
		return fmt.Errorf("got perr.Op: %s, want: %s", perr.Op, wantErr.Op)
	}
	if filepath.Clean(perr.Path) != filepath.Clean(wantErr.Path) {
		return fmt.Errorf("got perr.Path: %s, want: %s", perr.Path, wantErr.Path)
	}
	if perr.Err != wantErr.Err {
		return fmt.Errorf("got perr.Err: %s, want: %s", perr.Err, wantErr.Err)
	}
	return nil
}
