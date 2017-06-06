package rule

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal"
	"github.com/vlifesystems/rhkit/internal/dexprfuncs"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
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

func TestGenerate(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"team": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"a": description.Value{dlit.NewString("a"), 3},
					"b": description.Value{dlit.NewString("b"), 3},
					"c": description.Value{dlit.NewString("c"), 3},
				},
			},
			"teamOut": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"a": description.Value{dlit.NewString("a"), 3},
					"c": description.Value{dlit.NewString("c"), 1},
					"d": description.Value{dlit.NewString("d"), 3},
					"e": description.Value{dlit.NewString("e"), 3},
					"f": description.Value{dlit.NewString("f"), 3},
				},
			},
			"level": &description.Field{
				Kind:  fieldtype.Number,
				Min:   dlit.MustNew(0),
				Max:   dlit.MustNew(5),
				MaxDP: 0,
				Values: map[string]description.Value{
					"0": description.Value{dlit.NewString("0"), 3},
					"1": description.Value{dlit.NewString("1"), 3},
					"2": description.Value{dlit.NewString("2"), 1},
					"3": description.Value{dlit.NewString("3"), 3},
					"4": description.Value{dlit.NewString("4"), 3},
					"5": description.Value{dlit.NewString("5"), 3},
				},
			},
			"flow": &description.Field{
				Kind:  fieldtype.Number,
				Min:   dlit.MustNew(0),
				Max:   dlit.MustNew(10.5),
				MaxDP: 2,
				Values: map[string]description.Value{
					"0.0":  description.Value{dlit.NewString("0.0"), 3},
					"2.34": description.Value{dlit.NewString("2.34"), 3},
					"10.5": description.Value{dlit.NewString("10.5"), 3},
				},
			},
			"position": &description.Field{
				Kind:  fieldtype.Number,
				Min:   dlit.MustNew(1),
				Max:   dlit.MustNew(13),
				MaxDP: 0,
				Values: map[string]description.Value{
					"1":  description.Value{dlit.NewString("1"), 3},
					"2":  description.Value{dlit.NewString("2"), 3},
					"3":  description.Value{dlit.NewString("3"), 3},
					"4":  description.Value{dlit.NewString("4"), 3},
					"5":  description.Value{dlit.NewString("5"), 3},
					"6":  description.Value{dlit.NewString("6"), 3},
					"7":  description.Value{dlit.NewString("7"), 3},
					"8":  description.Value{dlit.NewString("8"), 3},
					"9":  description.Value{dlit.NewString("9"), 3},
					"10": description.Value{dlit.NewString("10"), 3},
					"11": description.Value{dlit.NewString("11"), 3},
					"12": description.Value{dlit.NewString("12"), 3},
					"13": description.Value{dlit.NewString("13"), 3},
				},
			},
		}}
	ruleFields :=
		[]string{"team", "teamOut", "level", "flow", "position"}
	wantRules := []Rule{
		NewTrue(),
		NewEQFV("team", dlit.MustNew("a")),
		NewNEFV("team", dlit.MustNew("a")),
		NewEQFF("team", "teamOut"),
		NewNEFF("team", "teamOut"),
		NewInFV("teamOut", makeStringsDlitSlice("a", "d")),
		NewEQFV("level", dlit.MustNew(0)),
		NewEQFV("level", dlit.MustNew(1)),
		NewNEFV("level", dlit.MustNew(0)),
		NewNEFV("level", dlit.MustNew(1)),
		NewLTFF("level", "position"),
		NewLEFF("level", "position"),
		NewNEFF("level", "position"),
		NewGEFF("level", "position"),
		NewGTFF("level", "position"),
		NewEQFF("level", "position"),
		NewGEFV("level", dlit.MustNew(1)),
		NewLEFV("level", dlit.MustNew(4)),
		NewInFV("level", makeStringsDlitSlice("0", "1")),
		NewInFV("level", makeStringsDlitSlice("0", "3")),
		NewGEFV("flow", dlit.MustNew(2.1)),
		NewGEFV("flow", dlit.MustNew(3.68)),
		NewLEFV("flow", dlit.MustNew(4.2)),
		NewLEFV("flow", dlit.MustNew(5.25)),
		NewAddLEF("level", "position", dlit.MustNew(12)),
		NewAddGEF("level", "position", dlit.MustNew(12)),
		NewMulLEF("flow", "level", dlit.MustNew(26.25)),
		NewMulGEF("flow", "level", dlit.MustNew(23.63)),
		MustNewBetweenFV("position", dlit.MustNew(9), dlit.MustNew(12)),
		MustNewOutsideFV("position", dlit.MustNew(9), dlit.MustNew(12)),
	}
	complexity := 10
	got := Generate(inputDescription, ruleFields, complexity)
	if err := rulesContain(got, wantRules); err != nil {
		t.Errorf("Generate: %s", err)
	}
}

func TestGenerate_combinations(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"directionIn": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"gogledd": description.Value{dlit.MustNew("gogledd"), 3},
					"de":      description.Value{dlit.MustNew("de"), 3},
				},
			},
			"directionOut": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"dwyrain":   description.Value{dlit.MustNew("dwyrain"), 3},
					"gorllewin": description.Value{dlit.MustNew("gorllewin"), 3},
				},
			},
		}}
	ruleFields := []string{"directionIn", "directionOut"}
	want := []Rule{
		NewEQFV("directionIn", dlit.MustNew("de")),
		NewEQFV("directionIn", dlit.MustNew("gogledd")),
		NewEQFV("directionOut", dlit.MustNew("dwyrain")),
		NewEQFV("directionOut", dlit.MustNew("gorllewin")),
		MustNewAnd(
			NewEQFV("directionIn", dlit.MustNew("de")),
			NewEQFV("directionOut", dlit.MustNew("dwyrain")),
		),
		MustNewAnd(
			NewEQFV("directionIn", dlit.MustNew("de")),
			NewEQFV("directionOut", dlit.MustNew("gorllewin")),
		),
		MustNewAnd(
			NewEQFV("directionIn", dlit.MustNew("gogledd")),
			NewEQFV("directionOut", dlit.MustNew("dwyrain")),
		),
		MustNewAnd(
			NewEQFV("directionIn", dlit.MustNew("gogledd")),
			NewEQFV("directionOut", dlit.MustNew("gorllewin")),
		),
		MustNewOr(
			NewEQFV("directionIn", dlit.MustNew("de")),
			NewEQFV("directionIn", dlit.MustNew("gogledd")),
		),
		MustNewOr(
			NewEQFV("directionIn", dlit.MustNew("de")),
			NewEQFV("directionOut", dlit.MustNew("dwyrain")),
		),
		MustNewOr(
			NewEQFV("directionIn", dlit.MustNew("de")),
			NewEQFV("directionOut", dlit.MustNew("gorllewin")),
		),
		MustNewOr(
			NewEQFV("directionIn", dlit.MustNew("gogledd")),
			NewEQFV("directionOut", dlit.MustNew("dwyrain")),
		),
		MustNewOr(
			NewEQFV("directionIn", dlit.MustNew("gogledd")),
			NewEQFV("directionOut", dlit.MustNew("gorllewin")),
		),
		MustNewOr(
			NewEQFV("directionOut", dlit.MustNew("dwyrain")),
			NewEQFV("directionOut", dlit.MustNew("gorllewin")),
		),
		NewTrue(),
	}

	complexity := 10
	got := Generate(inputDescription, ruleFields, complexity)
	Sort(got)
	Sort(want)
	if err := matchRulesUnordered(got, want); err != nil {
		t.Errorf("matchRulesUnordered: %s\n got: %s\nwant: %s\n",
			err, got, want)
	}
}

func TestCombine(t *testing.T) {
	cases := []struct {
		in   []Rule
		want []Rule
	}{
		{in: []Rule{
			NewEQFV("group", dlit.MustNew("a")),
			NewGEFV("band", dlit.MustNew(4)),
			NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
		},
			want: []Rule{
				MustNewAnd(
					NewGEFV("band", dlit.MustNew(4)),
					NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
				),
				MustNewAnd(
					NewGEFV("band", dlit.MustNew(4)),
					NewEQFV("group", dlit.MustNew("a")),
				),
				MustNewOr(
					NewGEFV("band", dlit.MustNew(4)),
					NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
				),
				MustNewOr(
					NewGEFV("band", dlit.MustNew(4)),
					NewEQFV("group", dlit.MustNew("a")),
				),
				MustNewAnd(
					NewEQFV("group", dlit.MustNew("a")),
					NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
				),
				MustNewOr(
					NewEQFV("group", dlit.MustNew("a")),
					NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
				),
			},
		},
		{in: []Rule{
			NewEQFV("team", dlit.MustNew("a")),
			NewGEFV("band", dlit.MustNew(4)),
			NewGEFV("band", dlit.MustNew(2)),
			NewGEFV("flow", dlit.MustNew(6)),
		},
			want: []Rule{
				MustNewAnd(NewGEFV("band", dlit.MustNew(2)),
					NewGEFV("flow", dlit.MustNew(6))),
				MustNewAnd(NewGEFV("band", dlit.MustNew(2)),
					NewEQFV("team", dlit.MustNew("a"))),
				MustNewOr(NewGEFV("band", dlit.MustNew(2)),
					NewGEFV("flow", dlit.MustNew(6))),
				MustNewOr(NewGEFV("band", dlit.MustNew(2)),
					NewEQFV("team", dlit.MustNew("a"))),
				MustNewAnd(NewGEFV("band", dlit.MustNew(4)),
					NewGEFV("flow", dlit.MustNew(6))),
				MustNewAnd(NewGEFV("band", dlit.MustNew(4)),
					NewEQFV("team", dlit.MustNew("a"))),
				MustNewOr(NewGEFV("band", dlit.MustNew(4)),
					NewGEFV("flow", dlit.MustNew(6))),
				MustNewOr(NewGEFV("band", dlit.MustNew(4)),
					NewEQFV("team", dlit.MustNew("a"))),
				MustNewAnd(NewGEFV("flow", dlit.MustNew(6)),
					NewEQFV("team", dlit.MustNew("a"))),
				MustNewOr(NewGEFV("flow", dlit.MustNew(6)),
					NewEQFV("team", dlit.MustNew("a"))),
			},
		},
		{in: []Rule{
			NewInFV("team", makeStringsDlitSlice("pink", "yellow", "blue")),
			NewInFV("team", makeStringsDlitSlice("red", "green", "blue")),
		},
			want: []Rule{
				NewInFV("team",
					makeStringsDlitSlice("pink", "yellow", "blue", "red", "green")),
			},
		},
		{in: []Rule{
			NewInFV("team", makeStringsDlitSlice("pink", "yellow", "blue")),
			NewInFV("group", makeStringsDlitSlice("red", "green", "blue")),
		},
			want: []Rule{
				MustNewAnd(
					NewInFV("group", makeStringsDlitSlice("red", "green", "blue")),
					NewInFV("team", makeStringsDlitSlice("pink", "yellow", "blue")),
				),
				MustNewOr(
					NewInFV("group", makeStringsDlitSlice("red", "green", "blue")),
					NewInFV("team", makeStringsDlitSlice("pink", "yellow", "blue")),
				),
			},
		},
		{in: []Rule{NewEQFV("team", dlit.MustNew("a"))},
			want: []Rule{}},
		{in: []Rule{}, want: []Rule{}},
	}

	for _, c := range cases {
		gotRules := Combine(c.in)
		if err := matchRulesUnordered(gotRules, c.want); err != nil {
			gotRuleStrs := rulesToSortedStrings(gotRules)
			wantRuleStrs := rulesToSortedStrings(c.want)
			t.Errorf("matchRulesUnordered() rules don't match: %s\n got: %s\n want: %s\n",
				err, gotRuleStrs, wantRuleStrs)
		}
	}
}
