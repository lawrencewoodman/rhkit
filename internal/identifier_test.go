package internal

import (
	"testing"
)

func TestIsIdentifierValid(t *testing.T) {
	cases := []struct {
		id   string
		want bool
	}{
		{"5hello", false},
		{"5Hello", false},
		{"h5ello", true},
		{"H5ello", true},
		{"H5Ello", true},
		{"hello", true},
		{"Hello", true},
		{"HELLO", true},
		{"h5_ello", true},
		{"h_ello", true},
		{"hello_everyone", true},
		{"HELLO_EVERYONE", true},
		{"_ello", false},
		{"e llo", false},
		{"h-ello", false},
	}

	for _, c := range cases {
		got := IsIdentifierValid(c.id)
		if got != c.want {
			t.Errorf("IsIdentifierValid(%s) got: %s, want: %s", c.id, got, c.want)
		}
	}
}
