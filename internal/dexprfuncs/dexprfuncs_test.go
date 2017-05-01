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

func TestPow(t *testing.T) {
	cases := []struct {
		in   []*dlit.Literal
		want *dlit.Literal
	}{
		{in: []*dlit.Literal{dlit.MustNew(0), dlit.MustNew(3)},
			want: dlit.MustNew(0),
		},
		{in: []*dlit.Literal{dlit.MustNew(0), dlit.MustNew(1.23)},
			want: dlit.MustNew(0),
		},
		{in: []*dlit.Literal{dlit.MustNew(0), dlit.MustNew(0)},
			want: dlit.MustNew(1),
		},
		{in: []*dlit.Literal{dlit.MustNew(2), dlit.MustNew(3)},
			want: dlit.MustNew(8),
		},
		{in: []*dlit.Literal{dlit.MustNew(4), dlit.MustNew(4.5)},
			want: dlit.MustNew(512),
		},
		{in: []*dlit.Literal{dlit.MustNew(2.5), dlit.MustNew(3)},
			want: dlit.MustNew(15.625),
		},
	}

	for _, c := range cases {
		got, err := pow(c.in)
		if err != nil {
			t.Errorf("pow(%v) err: %v", c.in, err)
		}
		if got.String() != c.want.String() {
			t.Errorf("pow(%v) got: %s, want: %s", c.in, got, c.want)
		}
	}
}

func TestPow_errors(t *testing.T) {
	cases := []struct {
		in   []*dlit.Literal
		want *dlit.Literal
		err  error
	}{
		{in: []*dlit.Literal{},
			want: dlit.MustNew(WrongNumOfArgsError{Got: 0, Want: 2}),
			err:  WrongNumOfArgsError{Got: 0, Want: 2},
		},
		{in: []*dlit.Literal{dlit.MustNew(23)},
			want: dlit.MustNew(WrongNumOfArgsError{Got: 1, Want: 2}),
			err:  WrongNumOfArgsError{Got: 1, Want: 2},
		},
		{in: []*dlit.Literal{dlit.MustNew(23), dlit.MustNew(4), dlit.MustNew(5)},
			want: dlit.MustNew(WrongNumOfArgsError{Got: 3, Want: 2}),
			err:  WrongNumOfArgsError{Got: 3, Want: 2},
		},
		{in: []*dlit.Literal{
			dlit.MustNew(errThisIsAnError),
			dlit.MustNew(errThisIsAnError),
		},
			want: dlit.MustNew(errThisIsAnError),
			err:  errThisIsAnError,
		},
		{in: []*dlit.Literal{dlit.NewString("hello"), dlit.MustNew(4)},
			want: dlit.MustNew(
				CantConvertToTypeError{Kind: "float", Value: dlit.NewString("hello")},
			),
			err: CantConvertToTypeError{
				Kind:  "float",
				Value: dlit.NewString("hello"),
			},
		},
		{in: []*dlit.Literal{dlit.MustNew(4), dlit.NewString("hello")},
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
		got, err := pow(c.in)
		checkErrorMatch(t, fmt.Sprintf("pow(%v)", c.in), err, c.err)
		if got.String() != c.want.String() {
			t.Errorf("pow(%v) got: %s, want: %s", c.in, got, c.want)
		}
	}
}

func TestMin(t *testing.T) {
	cases := []struct {
		in   []*dlit.Literal
		want *dlit.Literal
	}{
		{in: []*dlit.Literal{dlit.MustNew(16), dlit.MustNew(27)},
			want: dlit.MustNew(16),
		},
		{in: []*dlit.Literal{dlit.MustNew(27), dlit.MustNew(16)},
			want: dlit.MustNew(16),
		},
		{in: []*dlit.Literal{dlit.MustNew(27), dlit.MustNew(8), dlit.MustNew(17)},
			want: dlit.MustNew(8),
		},
		{in: []*dlit.Literal{dlit.MustNew(27), dlit.MustNew(18), dlit.MustNew(17)},
			want: dlit.MustNew(17),
		},
	}

	for _, c := range cases {
		got, err := min(c.in)
		if err != nil {
			t.Errorf("min(%v) err: %v", c.in, err)
		}
		if got.String() != c.want.String() {
			t.Errorf("min(%v) got: %s, want: %s",
				c.in, got, c.want)
		}
	}
}

func TestMin_errors(t *testing.T) {
	cases := []struct {
		in   []*dlit.Literal
		want *dlit.Literal
		err  error
	}{
		{in: []*dlit.Literal{},
			want: dlit.MustNew(ErrTooFewArguments),
			err:  ErrTooFewArguments,
		},
		{in: []*dlit.Literal{dlit.MustNew(23)},
			want: dlit.MustNew(ErrTooFewArguments),
			err:  ErrTooFewArguments,
		},
		{in: []*dlit.Literal{
			dlit.MustNew(errThisIsAnError),
			dlit.MustNew(errThisIsAnError),
		},
			want: dlit.MustNew(errThisIsAnError),
			err:  errThisIsAnError,
		},
		{in: []*dlit.Literal{
			dlit.MustNew(5),
			dlit.MustNew(errThisIsAnError),
		},
			want: dlit.MustNew(errThisIsAnError),
			err:  errThisIsAnError,
		},
		{in: []*dlit.Literal{
			dlit.MustNew(errThisIsAnError),
			dlit.MustNew(5),
		},
			want: dlit.MustNew(errThisIsAnError),
			err:  errThisIsAnError,
		},
		{in: []*dlit.Literal{dlit.NewString("hello"), dlit.MustNew(4)},
			want: dlit.MustNew(ErrIncompatibleTypes),
			err:  ErrIncompatibleTypes,
		},
		{in: []*dlit.Literal{dlit.MustNew(4), dlit.NewString("hello")},
			want: dlit.MustNew(ErrIncompatibleTypes),
			err:  ErrIncompatibleTypes,
		},
	}

	for _, c := range cases {
		got, err := min(c.in)
		checkErrorMatch(t, fmt.Sprintf("min(%v)", c.in), err, c.err)
		if got.String() != c.want.String() {
			t.Errorf("min(%v) got: %s, want: %s", c.in, got, c.want)
		}
	}
}

func TestMax(t *testing.T) {
	cases := []struct {
		in   []*dlit.Literal
		want *dlit.Literal
	}{
		{in: []*dlit.Literal{dlit.MustNew(16), dlit.MustNew(27)},
			want: dlit.MustNew(27),
		},
		{in: []*dlit.Literal{dlit.MustNew(27), dlit.MustNew(16)},
			want: dlit.MustNew(27),
		},
		{in: []*dlit.Literal{dlit.MustNew(17), dlit.MustNew(8), dlit.MustNew(27)},
			want: dlit.MustNew(27),
		},
		{in: []*dlit.Literal{dlit.MustNew(27), dlit.MustNew(18), dlit.MustNew(17)},
			want: dlit.MustNew(27),
		},
		{in: []*dlit.Literal{dlit.MustNew(18), dlit.MustNew(27), dlit.MustNew(17)},
			want: dlit.MustNew(27),
		},
	}

	for _, c := range cases {
		got, err := max(c.in)
		if err != nil {
			t.Errorf("max(%v) err: %v", c.in, err)
		}
		if got.String() != c.want.String() {
			t.Errorf("max(%v) got: %s, want: %s",
				c.in, got, c.want)
		}
	}
}

func TestMax_errors(t *testing.T) {
	cases := []struct {
		in   []*dlit.Literal
		want *dlit.Literal
		err  error
	}{
		{in: []*dlit.Literal{},
			want: dlit.MustNew(ErrTooFewArguments),
			err:  ErrTooFewArguments,
		},
		{in: []*dlit.Literal{dlit.MustNew(23)},
			want: dlit.MustNew(ErrTooFewArguments),
			err:  ErrTooFewArguments,
		},
		{in: []*dlit.Literal{
			dlit.MustNew(errThisIsAnError),
			dlit.MustNew(errThisIsAnError),
		},
			want: dlit.MustNew(errThisIsAnError),
			err:  errThisIsAnError,
		},
		{in: []*dlit.Literal{
			dlit.MustNew(5),
			dlit.MustNew(errThisIsAnError),
		},
			want: dlit.MustNew(errThisIsAnError),
			err:  errThisIsAnError,
		},
		{in: []*dlit.Literal{
			dlit.MustNew(errThisIsAnError),
			dlit.MustNew(5),
		},
			want: dlit.MustNew(errThisIsAnError),
			err:  errThisIsAnError,
		},
		{in: []*dlit.Literal{dlit.NewString("hello"), dlit.MustNew(4)},
			want: dlit.MustNew(ErrIncompatibleTypes),
			err:  ErrIncompatibleTypes,
		},
		{in: []*dlit.Literal{dlit.MustNew(4), dlit.NewString("hello")},
			want: dlit.MustNew(ErrIncompatibleTypes),
			err:  ErrIncompatibleTypes,
		},
	}

	for i, c := range cases {
		got, err := max(c.in)
		checkErrorMatch(t, fmt.Sprintf("max(%v)", c.in), err, c.err)
		if got.String() != c.want.String() {
			t.Errorf("i: %d, max(%v) got: %s, want: %s", i, c.in, got, c.want)
		}
	}
}

/*************************
 *       Benchmarks
 *************************/
func BenchmarkAlwaysTrue(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StartTimer()
		l, _ := alwaysTrue([]*dlit.Literal{})
		b.StopTimer()
		if v, ok := l.Bool(); !ok || !v {
			b.Errorf("alwaysTrue - got: %v, want: %t ", l, true)
		}
	}
}

func BenchmarkIn(b *testing.B) {
	b.StopTimer()
	haystack := []*dlit.Literal{
		dlit.MustNew(7),
		dlit.MustNew(21),
		dlit.MustNew(789),
		dlit.MustNew(56),
	}
	cases := []struct {
		needle *dlit.Literal
		want   bool
	}{
		{needle: dlit.MustNew(7), want: true},
		{needle: dlit.MustNew(789), want: true},
		{needle: dlit.MustNew(56), want: true},
		{needle: dlit.MustNew(89), want: false},
		{needle: dlit.MustNew(102), want: false},
		{needle: dlit.MustNew(78), want: false},
	}
	for n := 0; n < b.N; n++ {
		for _, c := range cases {
			b.StartTimer()
			l, _ := in(append([]*dlit.Literal{c.needle}, haystack...))
			b.StopTimer()
			if v, ok := l.Bool(); !ok || v != c.want {
				b.Errorf("alwaysTrue - got: %v, want: %t ", l, c.want)
			}
		}
	}
}

func BenchmarkNi(b *testing.B) {
	b.StopTimer()
	haystack := []*dlit.Literal{
		dlit.MustNew(7),
		dlit.MustNew(21),
		dlit.MustNew(789),
		dlit.MustNew(56),
	}
	cases := []struct {
		needle *dlit.Literal
		want   bool
	}{
		{needle: dlit.MustNew(7), want: false},
		{needle: dlit.MustNew(789), want: false},
		{needle: dlit.MustNew(56), want: false},
		{needle: dlit.MustNew(89), want: true},
		{needle: dlit.MustNew(102), want: true},
		{needle: dlit.MustNew(78), want: true},
	}
	for n := 0; n < b.N; n++ {
		for _, c := range cases {
			b.StartTimer()
			l, _ := ni(append([]*dlit.Literal{c.needle}, haystack...))
			b.StopTimer()
			if v, ok := l.Bool(); !ok || v != c.want {
				b.Errorf("alwaysTrue - got: %v, want: %t ", l, c.want)
			}
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
