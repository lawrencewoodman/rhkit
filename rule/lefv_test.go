package rule

import (
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/testhelpers"
	"reflect"
	"testing"
)

func TestLEFVString(t *testing.T) {
	field := "income"
	value := dlit.MustNew(8.93)
	want := "income <= 8.93"
	r := NewLEFV(field, value)
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestLEFVValue(t *testing.T) {
	field := "income"
	value := dlit.MustNew(8.93)
	r := NewLEFV(field, value)
	got := r.Value()
	if got.String() != "8.93" {
		t.Errorf("Value() got: %s, want: %s", got, value)
	}
}

func TestLEFVIsTrue(t *testing.T) {
	cases := []struct {
		field string
		value *dlit.Literal
		want  bool
	}{
		{"income", dlit.MustNew(19), true},
		{"income", dlit.MustNew(19.12), true},
		{"income", dlit.MustNew(20), true},
		{"income", dlit.MustNew(-20), false},
		{"income", dlit.MustNew(18.34), false},
		{"flow", dlit.MustNew(124.564), true},
		{"flow", dlit.MustNew(124.565), true},
		{"flow", dlit.MustNew(124.563), false},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"cost":   dlit.MustNew(20),
		"flow":   dlit.MustNew(124.564),
	}
	for _, c := range cases {
		r := NewLEFV(c.field, c.value)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) (rule: %s) err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestLEFVIsTrue_errors(t *testing.T) {
	cases := []struct {
		field   string
		value   *dlit.Literal
		wantErr error
	}{
		{field: "fred",
			value:   dlit.MustNew(7.894),
			wantErr: InvalidRuleError{Rule: NewLEFV("fred", dlit.MustNew(7.894))},
		},
		{field: "band",
			value: dlit.MustNew(7.894),
			wantErr: IncompatibleTypesRuleError{
				Rule: NewLEFV("band", dlit.MustNew(7.894)),
			},
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"flow":   dlit.MustNew(124.564),
		"band":   dlit.NewString("alpha"),
	}
	for _, c := range cases {
		r := NewLEFV(c.field, c.value)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}

func TestLEFVFields(t *testing.T) {
	r := NewLEFV("income", dlit.MustNew(5.5))
	want := []string{"income"}
	got := r.Fields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Fields() got: %s, want: %s", got, want)
	}
}

func TestLEFVTweak(t *testing.T) {
	field := "income"
	value := dlit.MustNew(800)
	rule := NewLEFV(field, value)
	cases := []struct {
		description *description.Description
		stage       int
		minNumRules int
		maxNumRules int
		min         *dlit.Literal
		max         *dlit.Literal
		mid         *dlit.Literal
		maxDP       int
	}{
		{description: &description.Description{
			map[string]*description.Field{
				"income": {
					Kind:  description.Number,
					Min:   dlit.MustNew(500),
					Max:   dlit.MustNew(1000),
					MaxDP: 2,
				},
			},
		},
			stage:       1,
			minNumRules: 18,
			maxNumRules: 20,
			min:         dlit.MustNew(755),
			max:         dlit.MustNew(845),
			mid:         dlit.MustNew(800),
			maxDP:       0,
		},
		{description: &description.Description{
			map[string]*description.Field{
				"income": {
					Kind:  description.Number,
					Min:   dlit.MustNew(790),
					Max:   dlit.MustNew(1000),
					MaxDP: 2,
				},
			},
		},
			stage:       1,
			minNumRules: 18,
			maxNumRules: 20,
			min:         dlit.MustNew(790),
			max:         dlit.MustNew(820),
			mid:         dlit.MustNew(803),
			maxDP:       2,
		},
		{description: &description.Description{
			map[string]*description.Field{
				"income": {
					Kind:  description.Number,
					Min:   dlit.MustNew(500),
					Max:   dlit.MustNew(810),
					MaxDP: 2,
				},
			},
		},
			stage:       1,
			minNumRules: 18,
			maxNumRules: 20,
			min:         dlit.MustNew(770),
			max:         dlit.MustNew(808),
			mid:         dlit.MustNew(787),
			maxDP:       2,
		},
		{description: &description.Description{
			map[string]*description.Field{
				"income": {
					Kind:  description.Number,
					Min:   dlit.MustNew(799),
					Max:   dlit.MustNew(801),
					MaxDP: 0,
				},
			},
		},
			stage:       1,
			minNumRules: 0,
			maxNumRules: 0,
			min:         dlit.MustNew(770),
			max:         dlit.MustNew(787),
			mid:         dlit.MustNew(808),
			maxDP:       0,
		},

		{description: &description.Description{
			map[string]*description.Field{
				"income": {
					Kind:  description.Number,
					Min:   dlit.MustNew(500),
					Max:   dlit.MustNew(1000),
					MaxDP: 2,
				},
			},
		},
			stage:       2,
			minNumRules: 18,
			maxNumRules: 20,
			min:         dlit.MustNew(777),
			max:         dlit.MustNew(823),
			mid:         dlit.MustNew(797),
			maxDP:       1,
		},
	}
	complyFunc := func(r Rule) error {
		x, ok := r.(*LEFV)
		if !ok {
			return fmt.Errorf("wrong type: %T (%s)", r, r)
		}
		if x.field != "income" {
			return fmt.Errorf("field isn't correct for rule: %s", r)
		}
		return nil
	}
	for i, c := range cases {
		got := rule.Tweak(c.description, c.stage)
		err := checkRulesComply(
			got,
			c.minNumRules,
			c.maxNumRules,
			c.min,
			c.max,
			c.mid,
			c.maxDP,
			complyFunc,
		)
		if err != nil {
			t.Errorf("(%d) Tweak: %s", i, err)
		}
	}
}

func TestGenerateLEFV(t *testing.T) {
	cases := []struct {
		description *description.Description
		field       string
		minNumRules int
		maxNumRules int
		min         *dlit.Literal
		max         *dlit.Literal
		mid         *dlit.Literal
		maxDP       int
	}{
		{description: &description.Description{
			map[string]*description.Field{
				"income": {
					Kind:  description.Number,
					Min:   dlit.MustNew(500),
					Max:   dlit.MustNew(1000),
					MaxDP: 2,
				},
			},
		},
			field:       "income",
			minNumRules: 18,
			maxNumRules: 20,
			min:         dlit.MustNew(500),
			max:         dlit.MustNew(1000),
			mid:         dlit.MustNew(750),
			maxDP:       0,
		},
		{description: &description.Description{
			map[string]*description.Field{
				"income": {
					Kind:  description.Number,
					Min:   dlit.MustNew(790.73),
					Max:   dlit.MustNew(1000),
					MaxDP: 2,
				},
			},
		},
			field:       "income",
			minNumRules: 18,
			maxNumRules: 20,
			min:         dlit.MustNew(790),
			max:         dlit.MustNew(1000),
			mid:         dlit.MustNew(903),
			maxDP:       2,
		},
		{description: &description.Description{
			map[string]*description.Field{
				"income": {
					Kind:  description.Number,
					Min:   dlit.MustNew(799),
					Max:   dlit.MustNew(801),
					MaxDP: 0,
				},
			},
		},
			field:       "income",
			minNumRules: 1,
			maxNumRules: 1,
			min:         dlit.MustNew(799),
			max:         dlit.MustNew(801),
			mid:         dlit.MustNew(800),
			maxDP:       0,
		},
		{description: &description.Description{
			map[string]*description.Field{
				"income": {
					Kind:  description.Number,
					Min:   dlit.MustNew(799),
					Max:   dlit.MustNew(801),
					MaxDP: 0,
				},
				"month": {
					Kind: description.String,
				},
			},
		},
			field:       "month",
			minNumRules: 0,
			maxNumRules: 0,
			min:         dlit.MustNew(0),
			max:         dlit.MustNew(0),
			mid:         dlit.MustNew(0),
			maxDP:       0,
		},
	}
	complyFunc := func(r Rule) error {
		x, ok := r.(*LEFV)
		if !ok {
			return fmt.Errorf("wrong type: %T (%s)", r, r)
		}
		if x.field != "income" {
			return fmt.Errorf("field isn't correct for rule: %s", r)
		}
		return nil
	}
	generationDesc := testhelpers.GenerationDesc{
		DFields:     []string{"income"},
		DArithmetic: false,
	}
	for i, c := range cases {
		got := generateLEFV(c.description, generationDesc, c.field)
		err := checkRulesComply(
			got,
			c.minNumRules,
			c.maxNumRules,
			c.min,
			c.max,
			c.mid,
			c.maxDP,
			complyFunc,
		)
		if err != nil {
			t.Errorf("(%d) generateLEFV: %s", i, err)
		}
	}
}

func TestLEFVOverlaps(t *testing.T) {
	cases := []struct {
		ruleA *LEFV
		ruleB Rule
		want  bool
	}{
		{ruleA: NewLEFV("band", dlit.MustNew(7.3)),
			ruleB: NewLEFV("band", dlit.MustNew(6.5)),
			want:  true,
		},
		{ruleA: NewLEFV("band", dlit.MustNew(7.3)),
			ruleB: NewLEFV("rate", dlit.MustNew(6.5)),
			want:  false,
		},
		{ruleA: NewLEFV("band", dlit.MustNew(7.3)),
			ruleB: NewGEFV("band", dlit.MustNew(6.5)),
			want:  false,
		},
	}
	for _, c := range cases {
		got := c.ruleA.Overlaps(c.ruleB)
		if got != c.want {
			t.Errorf("Overlaps - ruleA: %s, ruleB: %s - got: %t, want: %t",
				c.ruleA, c.ruleB, got, c.want)
		}
	}
}

/**************************
 *  Benchmarks
 **************************/

func BenchmarkLEFVIsTrue(b *testing.B) {
	record := map[string]*dlit.Literal{
		"band":   dlit.MustNew(23),
		"income": dlit.MustNew(1024),
		"cost":   dlit.MustNew(890.32),
	}
	r := NewLEFV("cost", dlit.MustNew(900.47))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = r.IsTrue(record)
	}
}
