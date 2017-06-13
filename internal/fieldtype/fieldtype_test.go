package fieldtype

import (
	"fmt"
	"testing"
)

func TestFieldTypeNew(t *testing.T) {
	cases := []struct {
		in   string
		want FieldType
	}{
		{"Unknown", Unknown},
		{"Ignore", Ignore},
		{"Number", Number},
		{"String", String},
	}

	for _, c := range cases {
		got := New(c.in)
		if got != c.want {
			t.Errorf("New: got: %s, want: %s", got, c.want)
		}
	}
}

func TestFieldTypeNew_panic(t *testing.T) {
	kind := "invalid"
	paniced := false
	wantPanic := fmt.Sprintf("unsupported type: %s", kind)
	defer func() {
		if r := recover(); r != nil {
			if r.(string) == wantPanic {
				paniced = true
			} else {
				t.Errorf("New: got panic: %s, wanted: %s", r, wantPanic)
			}
		}
	}()
	got := New(kind)
	if !paniced {
		t.Errorf("New: got: %s, failed to panic with: %s", got, wantPanic)
	}
}

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
			t.Errorf("String: c.in:%d got: %s, want: %s", c.in, got, c.want)
		}
	}
}

func TestFieldTypeString_panic(t *testing.T) {
	kind := FieldType(99)
	paniced := false
	wantPanic := fmt.Sprintf("unsupported type: %d", kind)
	defer func() {
		if r := recover(); r != nil {
			if r.(string) == wantPanic {
				paniced = true
			} else {
				t.Errorf("String: got panic: %s, wanted: %s", r, wantPanic)
			}
		}
	}()
	got := kind.String()
	if !paniced {
		t.Errorf("String: got: %s, failed to panic with: %s", got, wantPanic)
	}
}
