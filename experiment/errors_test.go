package experiment

import (
	"testing"
)

func TestErrNoRuleFieldsSpecifiedError(t *testing.T) {
	want := "no rule fields specified"
	got := ErrNoRuleFieldsSpecified.Error()
	if got != want {
		t.Errorf("Error() got: %s, want: %s", got, want)
	}
}

func TestInvalidSortDirectionErrorError(t *testing.T) {
	e := &InvalidSortDirectionError{"numMatches", "Ascending"}
	want := "invalid sort direction: Ascending, for field: numMatches"
	got := e.Error()
	if got != want {
		t.Errorf("Error() got: %s, want: %s", got, want)
	}
}

func TestInvalidSortFieldErrorError(t *testing.T) {
	e := InvalidSortFieldError("bob")
	want := "invalid sort field: bob"
	got := e.Error()
	if got != want {
		t.Errorf("Error() got: %s, want: %s", got, want)
	}
}

func TestInvalidRuleFieldErrorError(t *testing.T) {
	e := InvalidRuleFieldError("bob")
	want := "invalid rule field: bob"
	got := e.Error()
	if got != want {
		t.Errorf("Error() got: %s, want: %s", got, want)
	}
}

func TestInvalidAggregatorNameErrorError(t *testing.T) {
	e := InvalidAggregatorNameError("bob")
	want := "invalid aggregator name: bob"
	got := e.Error()
	if got != want {
		t.Errorf("Error() got: %s, want: %s", got, want)
	}
}

func TestAggregatorNameClashErrorError(t *testing.T) {
	e := AggregatorNameClashError("bob")
	want := "aggregator name clashes with field name: bob"
	got := e.Error()
	if got != want {
		t.Errorf("Error() got: %s, want: %s", got, want)
	}
}

func TestAggregatorNameReservedErrorError(t *testing.T) {
	e := AggregatorNameReservedError("bob")
	want := "aggregator name reserved: bob"
	got := e.Error()
	if got != want {
		t.Errorf("Error() got: %s, want: %s", got, want)
	}
}
