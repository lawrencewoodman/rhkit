package dexprfuncs

import (
	"errors"
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"testing"
)

var errThisIsAnError = errors.New("this is an error")

func TestIfFunc(t *testing.T) {
	cases := []struct {
		in   []*dlit.Literal
		want *dlit.Literal
	}{
		{in: []*dlit.Literal{
			dlit.MustNew(1),
			dlit.NewString("fred"),
			dlit.NewString("martha"),
		},
			want: dlit.NewString("fred"),
		},
		{in: []*dlit.Literal{
			dlit.MustNew(0),
			dlit.NewString("fred"),
			dlit.NewString("martha"),
		},
			want: dlit.NewString("martha"),
		},
	}

	for i, c := range cases {
		got, err := ifFunc(c.in)
		if err != nil {
			t.Errorf("[%d] ifFunc: %s", i, err)
		}
		if got.String() != c.want.String() {
			t.Errorf("[%d] ifFunc got: %s, want: %s", i, got, c.want)
		}
	}
}

func TestIfFunc_errors(t *testing.T) {
	cases := []struct {
		in   []*dlit.Literal
		want *dlit.Literal
		err  error
	}{
		{in: []*dlit.Literal{
			dlit.NewString("bob"),
			dlit.NewString("martha"),
		},
			want: dlit.MustNew(WrongNumOfArgsError{Got: 2, Want: 3}),
			err:  WrongNumOfArgsError{Got: 2, Want: 3},
		},
		{in: []*dlit.Literal{
			dlit.NewString("bob"),
			dlit.NewString("fred"),
			dlit.NewString("martha"),
		},
			want: dlit.MustNew(
				CantConvertToTypeError{Kind: "bool", Value: dlit.NewString("bob")},
			),
			err: CantConvertToTypeError{Kind: "bool", Value: dlit.NewString("bob")},
		},
		{in: []*dlit.Literal{
			dlit.MustNew(errThisIsAnError),
			dlit.NewString("fred"),
			dlit.NewString("martha"),
		},
			want: dlit.MustNew(errThisIsAnError),
			err:  errThisIsAnError,
		},
	}

	for i, c := range cases {
		got, err := ifFunc(c.in)
		checkErrorMatch(t, fmt.Sprintf("[%d] ifFunc", i), err, c.err)
		if got.String() != c.want.String() {
			t.Errorf("[%d] ifFunc got: %s, want: %s", i, got, c.want)
		}
	}
}

func TestIfErrFunc(t *testing.T) {
	cases := []struct {
		in   []*dlit.Literal
		want *dlit.Literal
	}{
		{in: []*dlit.Literal{
			dlit.NewString("fred"),
			dlit.NewString("martha"),
		},
			want: dlit.NewString("fred"),
		},
		{in: []*dlit.Literal{
			dlit.MustNew(errThisIsAnError),
			dlit.NewString("martha"),
		},
			want: dlit.NewString("martha"),
		},
	}

	for i, c := range cases {
		got, err := ifErrFunc(c.in)
		if err != nil {
			t.Errorf("[%d] ifErrFunc: %s", i, err)
		}
		if got.String() != c.want.String() {
			t.Errorf("[%d] ifErrFunc got: %s, want: %s", i, got, c.want)
		}
	}
}

func TestIfErrFunc_errors(t *testing.T) {
	cases := []struct {
		in   []*dlit.Literal
		want *dlit.Literal
		err  error
	}{
		{in: []*dlit.Literal{
			dlit.NewString("martha"),
		},
			want: dlit.MustNew(WrongNumOfArgsError{Got: 1, Want: 2}),
			err:  WrongNumOfArgsError{Got: 1, Want: 2},
		},
		{in: []*dlit.Literal{
			dlit.NewString("martha"),
			dlit.NewString("fred"),
			dlit.NewString("rebecca"),
		},
			want: dlit.MustNew(WrongNumOfArgsError{Got: 3, Want: 2}),
			err:  WrongNumOfArgsError{Got: 3, Want: 2},
		},
	}

	for i, c := range cases {
		got, err := ifErrFunc(c.in)
		checkErrorMatch(t, fmt.Sprintf("[%d] ifErrFunc", i), err, c.err)
		if got.String() != c.want.String() {
			t.Errorf("[%d] ifErrFunc got: %s, want: %s", i, got, c.want)
		}
	}
}

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
			dlit.MustNew(6),
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

func TestRoundTo(t *testing.T) {
	cases := []struct {
		in   []*dlit.Literal
		want *dlit.Literal
	}{
		{in: []*dlit.Literal{dlit.MustNew(0), dlit.MustNew(5)},
			want: dlit.MustNew(0)},
		{in: []*dlit.Literal{dlit.MustNew(3), dlit.MustNew(5)},
			want: dlit.MustNew(3)},
		{in: []*dlit.Literal{dlit.MustNew(3), dlit.MustNew(1)},
			want: dlit.MustNew(3)},
		{in: []*dlit.Literal{dlit.MustNew(3), dlit.MustNew(0)},
			want: dlit.MustNew(3)},
		{in: []*dlit.Literal{dlit.MustNew(5), dlit.MustNew(183)},
			want: dlit.MustNew(5)},
		{in: []*dlit.Literal{dlit.MustNew(2.445), dlit.MustNew(17)},
			want: dlit.MustNew(2.445)},
		{in: []*dlit.Literal{dlit.MustNew(2.445), dlit.MustNew(3)},
			want: dlit.MustNew(2.445)},
		{in: []*dlit.Literal{dlit.MustNew(2.445), dlit.MustNew(2)},
			want: dlit.MustNew(2.44)},
		{in: []*dlit.Literal{dlit.MustNew(2.445), dlit.MustNew(1)},
			want: dlit.MustNew(2.4)},
		{in: []*dlit.Literal{dlit.MustNew(2.445), dlit.MustNew(0)},
			want: dlit.MustNew(2)},
		{in: []*dlit.Literal{dlit.MustNew(2.512), dlit.MustNew(4)},
			want: dlit.MustNew(2.512)},
		{in: []*dlit.Literal{dlit.MustNew(2.512), dlit.MustNew(3)},
			want: dlit.MustNew(2.512)},
		{in: []*dlit.Literal{dlit.MustNew(2.512), dlit.MustNew(2)},
			want: dlit.MustNew(2.51)},
		{in: []*dlit.Literal{dlit.MustNew(2.512), dlit.MustNew(1)},
			want: dlit.MustNew(2.5)},
		{in: []*dlit.Literal{dlit.MustNew(2.5), dlit.MustNew(0)},
			want: dlit.MustNew(3)},
		{in: []*dlit.Literal{dlit.MustNew(-23.5), dlit.MustNew(1)},
			want: dlit.MustNew(-23.5)},
		{in: []*dlit.Literal{dlit.MustNew(-23.5), dlit.MustNew(0)},
			want: dlit.MustNew(-23)},
		{in: []*dlit.Literal{dlit.MustNew(0.268), dlit.MustNew(3)},
			want: dlit.MustNew(0.268)},
		{in: []*dlit.Literal{dlit.MustNew(0.268), dlit.MustNew(2)},
			want: dlit.MustNew(0.27)},
		{in: []*dlit.Literal{dlit.MustNew(0.268), dlit.MustNew(1)},
			want: dlit.MustNew(0.3)},
		{in: []*dlit.Literal{dlit.MustNew(0.268), dlit.MustNew(0)},
			want: dlit.MustNew(0)},
		{in: []*dlit.Literal{dlit.MustNew(0.258), dlit.MustNew(2)},
			want: dlit.MustNew(0.26)},
	}

	for _, c := range cases {
		got, err := roundTo(c.in)
		if err != nil {
			t.Errorf("roundTo(%v) err: %v", c.in, err)
		}
		if got.String() != c.want.String() {
			t.Errorf("roundTo(%v) got: %s, want: %s", c.in, got, c.want)
		}
	}
}

// Tests that an int is always returned no matter the number of dp
// This is important because rounding errors can be created because
// of the use of a float
func TestRoundTo_int(t *testing.T) {
	want := dlit.MustNew(5)
	for dp := 0; dp < 1000; dp++ {
		in := []*dlit.Literal{dlit.MustNew(5), dlit.MustNew(dp)}
		got, err := roundTo(in)
		if err != nil {
			t.Errorf("roundTo(%v) err: %v", in, err)
		}
		if got.String() != want.String() {
			t.Errorf("roundTo(%v) got: %s, want: %s", in, got, want)
		}
	}
}

// Tests that a float is always returned with a dp less then or equal
// the original number.
// This is important because rounding errors can be created because
// of the use of a float
func TestRoundTo_float(t *testing.T) {
	want := dlit.MustNew(5.55)
	for dp := 2; dp < 1000; dp++ {
		in := []*dlit.Literal{dlit.MustNew(5.55), dlit.MustNew(dp)}
		got, err := roundTo(in)
		if err != nil {
			t.Errorf("roundTo(%v) err: %v", in, err)
		}
		if got.String() != want.String() {
			t.Errorf("roundTo(%v) got: %s, want: %s", in, got, want)
		}
	}
}

func TestRoundTo_errors(t *testing.T) {
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
			dlit.MustNew(6),
		},
			want: dlit.MustNew(errThisIsAnError),
			err:  errThisIsAnError,
		},
		{in: []*dlit.Literal{
			dlit.MustNew(6.7),
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
		{in: []*dlit.Literal{
			dlit.MustNew(4.3),
			dlit.MustNew(6.7),
		},
			want: dlit.MustNew(
				CantConvertToTypeError{Kind: "int", Value: dlit.MustNew(6.7)},
			),
			err: CantConvertToTypeError{
				Kind:  "int",
				Value: dlit.MustNew(6.7),
			},
		},
	}

	for i, c := range cases {
		got, err := roundTo(c.in)
		checkErrorMatch(t, fmt.Sprintf("(%d) roundTo(%v)", i, c.in), err, c.err)
		if got.String() != c.want.String() {
			t.Errorf("(%d) roundTo(%v) got: %s, want: %s", i, c.in, got, c.want)
		}
	}
}

func TestIn(t *testing.T) {
	cases := []struct {
		in   []*dlit.Literal
		want *dlit.Literal
	}{
		{in: []*dlit.Literal{
			dlit.MustNew(3),
			dlit.MustNew(5),
			dlit.MustNew(7),
			dlit.MustNew("fred"),
			dlit.MustNew(9),
			dlit.MustNew(7),
		},
			want: falseLiteral},
		{in: []*dlit.Literal{
			dlit.MustNew(3),
			dlit.MustNew(5),
			dlit.MustNew(7),
			dlit.MustNew("fred"),
			dlit.MustNew(3),
			dlit.MustNew(7),
		},
			want: trueLiteral},
	}

	for _, c := range cases {
		got, err := in(c.in)
		if err != nil {
			t.Errorf("in(%v) err: %v", c.in, err)
		}
		if got.String() != c.want.String() {
			t.Errorf("in(%v) got: %s, want: %s", c.in, got, c.want)
		}
	}
}

func TestIn_errors(t *testing.T) {
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
			dlit.MustNew(6),
		},
			want: dlit.MustNew(errThisIsAnError),
			err:  errThisIsAnError,
		},
		{in: []*dlit.Literal{
			dlit.MustNew(6.7),
			dlit.MustNew(errThisIsAnError),
		},
			want: dlit.MustNew(errThisIsAnError),
			err:  errThisIsAnError,
		},
	}

	for i, c := range cases {
		got, err := in(c.in)
		checkErrorMatch(t, fmt.Sprintf("(%d) in(%v)", i, c.in), err, c.err)
		if got.String() != c.want.String() {
			t.Errorf("(%d) in(%v) got: %s, want: %s", i, c.in, got, c.want)
		}
	}
}

func TestNi(t *testing.T) {
	cases := []struct {
		in   []*dlit.Literal
		want *dlit.Literal
	}{
		{in: []*dlit.Literal{
			dlit.MustNew(3),
			dlit.MustNew(5),
			dlit.MustNew(7),
			dlit.MustNew("fred"),
			dlit.MustNew(9),
			dlit.MustNew(7),
		},
			want: trueLiteral},
		{in: []*dlit.Literal{
			dlit.MustNew(3),
			dlit.MustNew(5),
			dlit.MustNew(7),
			dlit.MustNew("fred"),
			dlit.MustNew(3),
			dlit.MustNew(7),
		},
			want: falseLiteral},
	}

	for _, c := range cases {
		got, err := ni(c.in)
		if err != nil {
			t.Errorf("ni(%v) err: %v", c.in, err)
		}
		if got.String() != c.want.String() {
			t.Errorf("ni(%v) got: %s, want: %s", c.in, got, c.want)
		}
	}
}

func TestNi_errors(t *testing.T) {
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
			dlit.MustNew(6),
		},
			want: dlit.MustNew(errThisIsAnError),
			err:  errThisIsAnError,
		},
		{in: []*dlit.Literal{
			dlit.MustNew(6.7),
			dlit.MustNew(errThisIsAnError),
		},
			want: dlit.MustNew(errThisIsAnError),
			err:  errThisIsAnError,
		},
	}

	for i, c := range cases {
		got, err := ni(c.in)
		checkErrorMatch(t, fmt.Sprintf("(%d) ni(%v)", i, c.in), err, c.err)
		if got.String() != c.want.String() {
			t.Errorf("(%d) ni(%v) got: %s, want: %s", i, c.in, got, c.want)
		}
	}
}

func TestAlwaysTrue(t *testing.T) {
	want := trueLiteral
	got, err := alwaysTrue([]*dlit.Literal{})
	if err != nil {
		t.Errorf("alwaysTrue err: %s", err)
	}
	if got.String() != want.String() {
		t.Errorf("alwaysTrue: got: %s, want: %s", got, want)
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
		return
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
