package rule

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/internal"
	"github.com/vlifesystems/rhkit/internal/dexprfuncs"
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
		NewEQFV("band", dlit.MustNew("b")),
		NewGEFV("flow", dlit.MustNew(3)),
		NewEQFV("band", dlit.MustNew("a")),
		NewGEFV("flow", dlit.MustNew(2)),
	}
	want := []Rule{
		NewEQFV("band", dlit.MustNew("a")),
		NewEQFV("band", dlit.MustNew("b")),
		NewGEFV("flow", dlit.MustNew(2)),
		NewGEFV("flow", dlit.MustNew(3)),
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
		NewEQFV("band", dlit.MustNew("b")),
		NewEQFV("band", dlit.MustNew("a")),
		NewGEFV("flow", dlit.MustNew(3)),
		NewEQFV("band", dlit.MustNew("a")),
		NewGEFV("flow", dlit.MustNew(2)),
	}
	want := []Rule{
		NewEQFV("band", dlit.MustNew("b")),
		NewEQFV("band", dlit.MustNew("a")),
		NewGEFV("flow", dlit.MustNew(3)),
		NewGEFV("flow", dlit.MustNew(2)),
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
func TestGenerateTweakPoints(t *testing.T) {
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
			wantNum: 0,
		},
		{value: dlit.MustNew(10),
			min:     dlit.MustNew(1),
			max:     dlit.MustNew(20),
			maxDP:   0,
			stage:   1,
			wantNum: 2,
		},
		{value: dlit.MustNew(5),
			min:     dlit.MustNew(1),
			max:     dlit.MustNew(10),
			maxDP:   3,
			stage:   1,
			wantNum: 16,
		},
		{value: dlit.MustNew(5),
			min:     dlit.MustNew(1),
			max:     dlit.MustNew(10),
			maxDP:   3,
			stage:   1,
			wantNum: 18,
		},
		{value: dlit.MustNew(5),
			min:     dlit.MustNew(1.278123),
			max:     dlit.MustNew(10.47529284),
			maxDP:   6,
			stage:   1,
			wantNum: 18,
		},
		{value: dlit.MustNew(5),
			min:     dlit.MustNew(1.278123),
			max:     dlit.MustNew(10.47529284),
			maxDP:   6,
			stage:   1,
			wantNum: 18,
		},
		{value: dlit.MustNew(5),
			min:     dlit.MustNew(1.278123),
			max:     dlit.MustNew(10.47529284),
			maxDP:   6,
			stage:   1,
			wantNum: 18,
		},
		{value: dlit.MustNew(800),
			min:     dlit.MustNew(790),
			max:     dlit.MustNew(1000),
			maxDP:   3,
			stage:   1,
			wantNum: 18,
		},
		{value: dlit.MustNew(990),
			min:     dlit.MustNew(790),
			max:     dlit.MustNew(1000),
			maxDP:   3,
			stage:   1,
			wantNum: 19,
		},
	}
	isValidExpr := dexpr.MustNew(
		"v != value && v > min && v < max && vNumDecPlaces <= maxDP",
		dexprfuncs.CallFuncs,
	)
	for i, c := range cases {
		vars := map[string]*dlit.Literal{
			"value": c.value,
			"min":   c.min,
			"max":   c.max,
			"maxDP": dlit.MustNew(c.maxDP),
		}
		got := generateTweakPoints(
			c.value,
			c.min,
			c.max,
			c.maxDP,
			c.stage,
		)
		if len(got) < c.wantNum || len(got) > (c.wantNum+2) {
			t.Errorf("(%d) generateTweakPoints(%s, %s, %s, %d, %d) got: %s, len(want): %d",
				i, c.value, c.min, c.max, c.maxDP, c.stage, got, c.wantNum)
		}
		for _, v := range got {
			vars["v"] = v
			vars["vNumDecPlaces"] = dlit.MustNew(internal.NumDecPlaces(v.String()))
			// TODO: Extend this test of validity
			if isValid, err := isValidExpr.EvalBool(vars); !isValid || err != nil {
				t.Errorf("(%d) generateTweakPoints(%s, %s, %s, %d, %d) invalid point: %s",
					i, c.value, c.min, c.max, c.maxDP, c.stage, v)
			}
		}
	}
}

func TestRoundRules(t *testing.T) {
	field := "income"
	cases := []struct {
		in   *dlit.Literal
		want []Rule
	}{
		{in: dlit.MustNew(5), want: []Rule{
			NewLEFV(field, dlit.MustNew(5)),
		}},
		{in: dlit.MustNew(2.5), want: []Rule{
			NewLEFV(field, dlit.MustNew(2.5)),
			NewLEFV(field, dlit.MustNew(3)),
		}},
		{in: dlit.MustNew(2.25), want: []Rule{
			NewLEFV(field, dlit.MustNew(2.25)),
			NewLEFV(field, dlit.MustNew(2.3)),
			NewLEFV(field, dlit.MustNew(2)),
		}},
	}

	makeRule := func(p *dlit.Literal) Rule {
		return NewLEFV(field, p)
	}
	for _, c := range cases {
		got := roundRules(c.in, makeRule)
		if len(got) != len(c.want) {
			t.Errorf("roundRules got: %s, want: %s", got, c.want)
			continue
		}
		for i, n := range got {
			if n.String() != c.want[i].String() {
				t.Errorf("roundRules got: %s, want: %s", got, c.want)
			}
		}
	}
}

func TestReduceDP(t *testing.T) {
	in := []Rule{
		NewLEFV("income", dlit.MustNew(7.772)),
		NewGEFV("flow", dlit.MustNew(7.9265)),
		NewGEFF("flow", "income"),
		NewAddLEF("balance", "income", dlit.MustNew(1024.23)),
		NewAddGEF("balance", "income", dlit.MustNew(1024.23)),
	}
	want := []Rule{
		NewLEFV("income", dlit.MustNew(7.772)),
		NewLEFV("income", dlit.MustNew(7.77)),
		NewLEFV("income", dlit.MustNew(7.8)),
		NewLEFV("income", dlit.MustNew(8)),
		NewGEFV("flow", dlit.MustNew(7.9265)),
		NewGEFV("flow", dlit.MustNew(7.927)),
		NewGEFV("flow", dlit.MustNew(7.93)),
		NewGEFV("flow", dlit.MustNew(7.9)),
		NewGEFV("flow", dlit.MustNew(8)),
		NewAddLEF("balance", "income", dlit.MustNew(1024.23)),
		NewAddLEF("balance", "income", dlit.MustNew(1024.2)),
		NewAddLEF("balance", "income", dlit.MustNew(1024)),
		NewAddGEF("balance", "income", dlit.MustNew(1024.23)),
		NewAddGEF("balance", "income", dlit.MustNew(1024.2)),
		NewAddGEF("balance", "income", dlit.MustNew(1024)),
		NewTrue(),
	}
	got := ReduceDP(in)
	if len(got) != len(want) {
		t.Errorf("ReduceDP got: %s, want: %s", got, want)
		return
	}
	for i, n := range got {
		if n.String() != want[i].String() {
			t.Errorf("ReduceDP got: %s, want: %s", got, want)
		}
	}
}
