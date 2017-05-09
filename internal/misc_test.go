package internal

import (
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestNumDecPlaces(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"5", 0},
		{"25", 0},
		{"235", 0},
		{"235.0", 0},
		{"235.00", 0},
		{".5", 1},
		{".50", 1},
		{"0.5", 1},
		{"00.5", 1},
		{"00.50", 1},
		{"1.50", 1},
		{"123.50", 1},
		{".23", 2},
		{".230", 2},
		{"0.230", 2},
		{"00.230", 2},
		{"5.230", 2},
		{"25.230", 2},
		{"325.230", 2},
		{".234", 3},
		{".2340", 3},
		{"0.2340", 3},
		{"00.2340", 3},
		{"5.2340", 3},
		{"25.2340", 3},
		{"325.2340", 3},
	}

	for _, c := range cases {
		got := NumDecPlaces(c.in)
		if got != c.want {
			t.Errorf("NumDecPlaces(%s) got: %d, want: %d", c.in, got, c.want)
		}
	}
}

func TestGeneratePoints(t *testing.T) {
	cases := []struct {
		min        *dlit.Literal
		max        *dlit.Literal
		maxShortDP int
		maxDP      int
		want       []*dlit.Literal
	}{
		{min: dlit.MustNew(10), max: dlit.MustNew(15), maxShortDP: 0, maxDP: 0,
			want: []*dlit.Literal{
				dlit.MustNew(11), dlit.MustNew(12), dlit.MustNew(13), dlit.MustNew(14),
			},
		},
		{min: dlit.MustNew(10), max: dlit.MustNew(50), maxShortDP: 0, maxDP: 0,
			want: []*dlit.Literal{
				dlit.MustNew(12), dlit.MustNew(14), dlit.MustNew(16), dlit.MustNew(18),
				dlit.MustNew(20), dlit.MustNew(22), dlit.MustNew(24), dlit.MustNew(26),
				dlit.MustNew(28), dlit.MustNew(30), dlit.MustNew(32), dlit.MustNew(34),
				dlit.MustNew(36), dlit.MustNew(38), dlit.MustNew(40), dlit.MustNew(42),
				dlit.MustNew(44), dlit.MustNew(46), dlit.MustNew(48),
			},
		},
		{min: dlit.MustNew(10), max: dlit.MustNew(15), maxShortDP: 0, maxDP: 1,
			want: []*dlit.Literal{
				dlit.MustNew(10), dlit.MustNew(10.3), dlit.MustNew(10.6),
				dlit.MustNew(10.9), dlit.MustNew(11), dlit.MustNew(11.2),
				dlit.MustNew(11.5), dlit.MustNew(11.8), dlit.MustNew(12),
				dlit.MustNew(12.1), dlit.MustNew(12.4), dlit.MustNew(12.7),
				dlit.MustNew(13), dlit.MustNew(13.3), dlit.MustNew(13.6),
				dlit.MustNew(13.9), dlit.MustNew(14), dlit.MustNew(14.2),
				dlit.MustNew(14.5),
			},
		},
		{min: dlit.MustNew(10.2678945), max: dlit.MustNew(15.2382745),
			maxShortDP: 0, maxDP: 6,
			want: []*dlit.Literal{
				dlit.MustNew(10.516414), dlit.MustNew(10.764933), dlit.MustNew(11),
				dlit.MustNew(11.013452), dlit.MustNew(11.261971),
				dlit.MustNew(11.51049), dlit.MustNew(11.759009),
				dlit.MustNew(12), dlit.MustNew(12.007528), dlit.MustNew(12.256047),
				dlit.MustNew(12.504566), dlit.MustNew(12.753085), dlit.MustNew(13),
				dlit.MustNew(13.001604), dlit.MustNew(13.250123),
				dlit.MustNew(13.498642), dlit.MustNew(13.747161),
				dlit.MustNew(13.99568), dlit.MustNew(14), dlit.MustNew(14.244199),
				dlit.MustNew(14.492718), dlit.MustNew(14.741237),
				dlit.MustNew(14.989756), dlit.MustNew(15),
			},
		},
		{min: dlit.MustNew(10.2678945), max: dlit.MustNew(15.2382745),
			maxShortDP: 1, maxDP: 6,
			want: []*dlit.Literal{
				dlit.MustNew(10.5), dlit.MustNew(10.516414), dlit.MustNew(10.764933),
				dlit.MustNew(10.8), dlit.MustNew(11), dlit.MustNew(11.013452),
				dlit.MustNew(11.261971), dlit.MustNew(11.3), dlit.MustNew(11.5),
				dlit.MustNew(11.51049), dlit.MustNew(11.759009), dlit.MustNew(11.8),
				dlit.MustNew(12), dlit.MustNew(12.007528), dlit.MustNew(12.256047),
				dlit.MustNew(12.3), dlit.MustNew(12.5), dlit.MustNew(12.504566),
				dlit.MustNew(12.753085), dlit.MustNew(12.8), dlit.MustNew(13),
				dlit.MustNew(13.001604), dlit.MustNew(13.250123),
				dlit.MustNew(13.3), dlit.MustNew(13.498642), dlit.MustNew(13.5),
				dlit.MustNew(13.7), dlit.MustNew(13.747161), dlit.MustNew(13.99568),
				dlit.MustNew(14), dlit.MustNew(14.2), dlit.MustNew(14.244199),
				dlit.MustNew(14.492718), dlit.MustNew(14.5), dlit.MustNew(14.7),
				dlit.MustNew(14.741237), dlit.MustNew(14.989756), dlit.MustNew(15),
			},
		},
		{min: dlit.MustNew(10.2678945), max: dlit.MustNew(15.2382745),
			maxShortDP: 2, maxDP: 6,
			want: []*dlit.Literal{
				dlit.MustNew(10.5), dlit.MustNew(10.516414), dlit.MustNew(10.52),
				dlit.MustNew(10.76), dlit.MustNew(10.764933), dlit.MustNew(10.8),
				dlit.MustNew(11), dlit.MustNew(11.01), dlit.MustNew(11.013452),
				dlit.MustNew(11.26), dlit.MustNew(11.261971), dlit.MustNew(11.3),
				dlit.MustNew(11.5), dlit.MustNew(11.51), dlit.MustNew(11.51049),
				dlit.MustNew(11.759009), dlit.MustNew(11.76), dlit.MustNew(11.8),
				dlit.MustNew(12), dlit.MustNew(12.007528), dlit.MustNew(12.01),
				dlit.MustNew(12.256047), dlit.MustNew(12.26), dlit.MustNew(12.3),
				dlit.MustNew(12.5), dlit.MustNew(12.504566), dlit.MustNew(12.75),
				dlit.MustNew(12.753085), dlit.MustNew(12.8), dlit.MustNew(13),
				dlit.MustNew(13.001604), dlit.MustNew(13.25), dlit.MustNew(13.250123),
				dlit.MustNew(13.3), dlit.MustNew(13.498642), dlit.MustNew(13.5),
				dlit.MustNew(13.7), dlit.MustNew(13.747161), dlit.MustNew(13.75),
				dlit.MustNew(13.99568), dlit.MustNew(14), dlit.MustNew(14.2),
				dlit.MustNew(14.24), dlit.MustNew(14.244199), dlit.MustNew(14.49),
				dlit.MustNew(14.492718), dlit.MustNew(14.5), dlit.MustNew(14.7),
				dlit.MustNew(14.74), dlit.MustNew(14.741237),
				dlit.MustNew(14.989756), dlit.MustNew(14.99), dlit.MustNew(15),
			},
		},
	}
	for i, c := range cases {
		got := GeneratePoints(c.min, c.max, c.maxShortDP, c.maxDP)
		if err := checkPoints(got, c.want); err != nil {
			t.Errorf("(%d)GeneratePoints: %s", i, err)
		}
	}
}

func checkPoints(got, want []*dlit.Literal) error {
	if len(got) != len(want) {
		return fmt.Errorf("len(got): %d != len(want): %d, got: %v",
			len(got), len(want), got)
	}
	for i, p := range got {
		if want[i].String() != p.String() {
			return fmt.Errorf("got: %v, want: %v", got, want)
		}
	}
	return nil
}
