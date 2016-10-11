package rhkit

import (
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
		got := numDecPlaces(c.in)
		if got != c.want {
			t.Errorf("numDecPlaces(%s) got: %d, want: %d", c.in, got, c.want)
		}
	}
}
