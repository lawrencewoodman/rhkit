package fieldtype

import (
	"fmt"
	"testing"
)

func TestFieldTypeString(t *testing.T) {
	cases := []struct {
		in   FieldType
		want string
	}{
		{Unknown, "Unknown"},
		{Ignore, "Ignore"},
		{Number, "Number"},
		{String, "String"},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("String() c.in:%d got: %s, want: %s", c.in, got, c.want)
		}
	}
}

func TestFieldTypeString_panic(t *testing.T) {
	kind := FieldType(99)
	paniced := false
	wantPanic := fmt.Sprintf("Unsupported type: %d", kind)
	defer func() {
		if r := recover(); r != nil {
			if r.(string) == wantPanic {
				paniced = true
			} else {
				t.Errorf("String() - got panic: %s, wanted: %s", r, wantPanic)
			}
		}
	}()
	got := kind.String()
	if !paniced {
		t.Errorf("String() - got: %s, failed to panic with: %s", got, wantPanic)
	}
}
