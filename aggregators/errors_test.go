package aggregators

import "testing"

func TestInvalidNameErrorError(t *testing.T) {
	e := InvalidNameError("bob")
	want := "invalid aggregator name: bob"
	got := e.Error()
	if got != want {
		t.Errorf("Error() got: %s, want: %s", got, want)
	}
}

func TestNameClashErrorError(t *testing.T) {
	e := NameClashError("bob")
	want := "aggregator name clashes with field name: bob"
	got := e.Error()
	if got != want {
		t.Errorf("Error() got: %s, want: %s", got, want)
	}
}

func TestNameReservedErrorError(t *testing.T) {
	e := NameReservedError("bob")
	want := "aggregator name reserved: bob"
	got := e.Error()
	if got != want {
		t.Errorf("Error() got: %s, want: %s", got, want)
	}
}
