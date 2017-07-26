package aggregators

import "testing"

func TestDescErrorError(t *testing.T) {
	e := DescError{Name: "myName", Kind: "myKind", Err: ErrUnregisteredKind}
	want := "problem with aggregator description - name: myName, kind: myKind (unregistered kind)"
	got := e.Error()
	if got != want {
		t.Errorf("Error() got: %s, want: %s", got, want)
	}
}
