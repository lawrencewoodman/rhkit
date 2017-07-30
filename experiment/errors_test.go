package experiment

import "testing"

func TestErrNoRuleFieldsSpecifiedError(t *testing.T) {
	want := "no rule fields specified"
	got := ErrNoRuleFieldsSpecified.Error()
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
