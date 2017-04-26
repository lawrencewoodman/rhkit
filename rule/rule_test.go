package rule

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/internal"
	"testing"
)

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

// TODO: Expand this test
func TestGeneratePoints(t *testing.T) {
	cases := []struct {
		value   *dlit.Literal
		min     *dlit.Literal
		max     *dlit.Literal
		maxDP   int
		stage   int
		wantNum int
	}{
		{value: dlit.MustNew(5),
			min:     dlit.MustNew(10),
			max:     dlit.MustNew(10),
			maxDP:   0,
			stage:   1,
			wantNum: 0,
		},
		{value: dlit.MustNew(5),
			min:     dlit.MustNew(10),
			max:     dlit.MustNew(10),
			maxDP:   50,
			stage:   1,
			wantNum: 0,
		},
		{value: dlit.MustNew(800),
			min:     dlit.MustNew(500),
			max:     dlit.MustNew(1000),
			maxDP:   0,
			stage:   1,
			wantNum: 18,
		},
		{value: dlit.MustNew(5),
			min:     dlit.MustNew(1),
			max:     dlit.MustNew(10),
			maxDP:   0,
			stage:   1,
			wantNum: 2,
		},
		{value: dlit.MustNew(5),
			min:     dlit.MustNew(1),
			max:     dlit.MustNew(10),
			maxDP:   3,
			stage:   1,
			wantNum: 18,
		},
	}
	for _, c := range cases {
		got := generatePoints(c.value, c.min, c.max, c.maxDP, c.stage)
		if len(got) != c.wantNum {
			t.Errorf("generatePoints(%s, %s, %s, %d, %d) got: %s, len(want): %d",
				c.value, c.min, c.max, c.maxDP, c.stage, got, c.wantNum)
		}
		for _, v := range got {
			// TODO: Extend this test of validity
			if v.String() == c.value.String() ||
				v.String() == c.min.String() ||
				v.String() == c.max.String() ||
				internal.NumDecPlaces(v.String()) > c.maxDP {
				t.Errorf("generatePoints(%s, %s, %s, %d, %d) invalid point: %s",
					c.value, c.min, c.max, c.maxDP, c.stage, v)
			}
		}
	}
}
