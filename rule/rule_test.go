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

func TestSort(t *testing.T) {
	in := []Rule{
		NewEQFVS("band", "b"),
		NewGEFVI("flow", 3),
		NewEQFVS("band", "a"),
		NewGEFVI("flow", 2),
	}
	want := []Rule{
		NewEQFVS("band", "a"),
		NewEQFVS("band", "b"),
		NewGEFVI("flow", 2),
		NewGEFVI("flow", 3),
	}
	Sort(in)
	if len(in) != len(want) {
		t.Fatalf("Sort - len(in) != len(want)")
	}
	for i, r := range want {
		if in[i].String() != r.String() {
			t.Fatalf("Sort - got: %v, want: %v", in, want)
		}
	}
}

func TestUniq(t *testing.T) {
	in := []Rule{
		NewEQFVS("band", "b"),
		NewEQFVS("band", "a"),
		NewGEFVI("flow", 3),
		NewEQFVS("band", "a"),
		NewGEFVI("flow", 2),
	}
	want := []Rule{
		NewEQFVS("band", "b"),
		NewEQFVS("band", "a"),
		NewGEFVI("flow", 3),
		NewGEFVI("flow", 2),
	}
	got := Uniq(in)
	if len(got) != len(want) {
		t.Fatalf("Sort - len(got) != len(want)")
	}
	for i, r := range want {
		if got[i].String() != r.String() {
			t.Fatalf("Sort - got: %v, want: %v", got, want)
		}
	}
}
