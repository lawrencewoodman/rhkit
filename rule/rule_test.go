package rule

import "testing"

func TestInvalidRuleErrorError(t *testing.T) {
	r := NewTrue()
	err := InvalidRuleError{r}
	want := "invalid rule: true()"
	got := err.Error()
	if got != want {
		t.Errorf("Error() got: %s, want: %s", got, want)
	}
}

func TestIncompatibleTypesRuleErrorError(t *testing.T) {
	r := NewTrue()
	err := IncompatibleTypesRuleError{r}
	want := "incompatible types in rule: true()"
	got := err.Error()
	if got != want {
		t.Errorf("Error() got: %s, want: %s", got, want)
	}
}
