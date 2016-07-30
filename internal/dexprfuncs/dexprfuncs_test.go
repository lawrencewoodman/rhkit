package dexprfuncs

import (
	"errors"
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestSqrt(t *testing.T) {
	cases := []struct {
		in   *dlit.Literal
		want *dlit.Literal
	}{
		{in: dlit.MustNew(16), want: dlit.MustNew(4)},
		{in: dlit.MustNew(2.25 * 2.25), want: dlit.MustNew(2.25)},
	}

	for _, c := range cases {
		got, err := sqrt([]*dlit.Literal{c.in})
		if err != nil {
			t.Errorf("sqrt(%v) err: %v", []*dlit.Literal{c.in}, err)
		}
		if got.String() != c.want.String() {
			t.Errorf("sqrt(%v) got: %s, want: %s",
				[]*dlit.Literal{c.in}, got, c.want)
		}
	}
}

var errThisIsAnError = errors.New("this is an error")

func TestSqrt_errors(t *testing.T) {
	cases := []struct {
		in   []*dlit.Literal
		want *dlit.Literal
		err  error
	}{
		{in: []*dlit.Literal{},
			want: dlit.MustNew(WrongNumOfArgsError{Got: 0, Want: 1}),
			err:  WrongNumOfArgsError{Got: 0, Want: 1},
		},
		{in: []*dlit.Literal{dlit.MustNew(23), dlit.MustNew(4)},
			want: dlit.MustNew(WrongNumOfArgsError{Got: 2, Want: 1}),
			err:  WrongNumOfArgsError{Got: 2, Want: 1},
		},
		{in: []*dlit.Literal{dlit.MustNew(errThisIsAnError)},
			want: dlit.MustNew(errThisIsAnError),
			err:  errThisIsAnError,
		},
		{in: []*dlit.Literal{dlit.NewString("hello")},
			want: dlit.MustNew(
				CantConvertToTypeError{Kind: "float", Value: dlit.NewString("hello")},
			),
			err: CantConvertToTypeError{
				Kind:  "float",
				Value: dlit.NewString("hello"),
			},
		},
	}

	for _, c := range cases {
		got, err := sqrt(c.in)
		checkErrorMatch(t, fmt.Sprintf("sqrt(%v)", c.in), err, c.err)
		if got.String() != c.want.String() {
			t.Errorf("sqrt(%v) got: %s, want: %s", c.in, got, c.want)
		}
	}
}

/*************************************
 *  Helper functions
 *************************************/

func checkErrorMatch(t *testing.T, context string, got, want error) {
	if got == nil && want == nil {
		return
	}
	if got == nil || want == nil {
		t.Errorf("%s got err: %s, want : %s", context, got, want)
	}
	switch x := want.(type) {
	case CantConvertToTypeError:
		checkCantConvertToTypeError(t, context, got, x)
		return
	}
	if got != want {
		t.Errorf("%s got err: %s, want : %s", context, got, want)
	}
}

func checkCantConvertToTypeError(
	t *testing.T,
	context string,
	got, want error,
) {
	gerr, ok := got.(CantConvertToTypeError)
	if !ok {
		t.Errorf(
			"%s got err type: %T, want error type: CantConvertToTypeError",
			context,
			got,
		)
	}
	werr, ok := want.(CantConvertToTypeError)
	if !ok {
		panic("want isn't type CantConvertToTypeError")
	}
	if gerr.Kind != werr.Kind {
		t.Errorf("%s got: %s, want: %s", context, got, want)
	}
	if gerr.Value.String() != werr.Value.String() {
		t.Errorf("%s got: %s, want: %s", context, got, want)
	}
}
