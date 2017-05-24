package rule

import (
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
	"reflect"
	"testing"
)

func TestAddGEFString(t *testing.T) {
	fieldA := "income"
	fieldB := "balance"
	value := dlit.MustNew(8.93)
	want := "income + balance >= 8.93"
	r := NewAddGEF(fieldA, fieldB, value)
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestAddGEFIsTrue(t *testing.T) {
	cases := []struct {
		fieldA string
		fieldB string
		value  *dlit.Literal
		want   bool
	}{
		{"income", "balance", dlit.MustNew(19), true},
		{"income", "balance", dlit.MustNew(19.12), false},
		{"income", "balance", dlit.MustNew(20), false},
		{"income", "balance", dlit.MustNew(-20), true},
		{"income", "balance", dlit.MustNew(18.34), true},
		{"flow", "cost", dlit.MustNew(144.564), true},
		{"flow", "cost", dlit.MustNew(144.565), false},
		{"flow", "cost", dlit.MustNew(144.563), true},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(4),
		"balance": dlit.MustNew(15),
		"cost":    dlit.MustNew(20),
		"flow":    dlit.MustNew(124.564),
	}
	for _, c := range cases {
		r := NewAddGEF(c.fieldA, c.fieldB, c.value)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) (rule: %s) err: %v", r, err)
		} else if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestAddGEFIsTrue_errors(t *testing.T) {
	cases := []struct {
		fieldA  string
		fieldB  string
		value   *dlit.Literal
		wantErr error
	}{
		{fieldA: "fred",
			fieldB: "flow",
			value:  dlit.MustNew(7.894),
			wantErr: InvalidRuleError{
				Rule: NewAddGEF("fred", "flow", dlit.MustNew(7.894)),
			},
		},
		{fieldA: "flow",
			fieldB: "fred",
			value:  dlit.MustNew(7.894),
			wantErr: InvalidRuleError{
				Rule: NewAddGEF("flow", "fred", dlit.MustNew(7.894)),
			},
		},
		{fieldA: "band",
			fieldB: "flow",
			value:  dlit.MustNew(7.894),
			wantErr: IncompatibleTypesRuleError{
				Rule: NewAddGEF("band", "flow", dlit.MustNew(7.894)),
			},
		},
		{fieldA: "flow",
			fieldB: "band",
			value:  dlit.MustNew(7.894),
			wantErr: IncompatibleTypesRuleError{
				Rule: NewAddGEF("flow", "band", dlit.MustNew(7.894)),
			},
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"flow":   dlit.MustNew(124.564),
		"band":   dlit.NewString("alpha"),
	}
	for _, c := range cases {
		r := NewAddGEF(c.fieldA, c.fieldB, c.value)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}

func TestAddGEFGetFields(t *testing.T) {
	r := NewAddGEF("income", "cost", dlit.MustNew(5.5))
	want := []string{"income", "cost"}
	got := r.GetFields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetFields() got: %s, want: %s", got, want)
	}
}

func TestAddGEFOverlaps(t *testing.T) {
	cases := []struct {
		ruleA *AddGEF
		ruleB Rule
		want  bool
	}{
		{ruleA: NewAddGEF("band", "cost", dlit.MustNew(7.3)),
			ruleB: NewAddGEF("band", "cost", dlit.MustNew(6.5)),
			want:  true,
		},
		{ruleA: NewAddGEF("band", "balance", dlit.MustNew(7.3)),
			ruleB: NewAddGEF("rate", "balance", dlit.MustNew(6.5)),
			want:  false,
		},
		{ruleA: NewAddGEF("band", "balance", dlit.MustNew(7.3)),
			ruleB: NewAddGEF("band", "rate", dlit.MustNew(6.5)),
			want:  false,
		},
		{ruleA: NewAddGEF("band", "cost", dlit.MustNew(7.3)),
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

func TestAddGEFTweak(t *testing.T) {
	fieldA := "income"
	fieldB := "balance"
	value := int64(800)
	rule := NewAddGEF(fieldA, fieldB, dlit.MustNew(value))
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
				"balance": &description.Field{
					Kind: fieldtype.Int,
					Min:  dlit.MustNew(250),
					Max:  dlit.MustNew(500),
				},
				"income": &description.Field{
					Kind: fieldtype.Int,
					Min:  dlit.MustNew(250),
					Max:  dlit.MustNew(500),
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
				"balance": &description.Field{
					Kind: fieldtype.Int,
					Min:  dlit.MustNew(250),
					Max:  dlit.MustNew(300),
				},
				"income": &description.Field{
					Kind: fieldtype.Int,
					Min:  dlit.MustNew(540),
					Max:  dlit.MustNew(700),
				},
			},
		},
			stage:       1,
			minNumRules: 18,
			maxNumRules: 20,
			min:         dlit.MustNew(792),
			max:         dlit.MustNew(819),
			mid:         dlit.MustNew(804),
			maxDP:       0,
		},
		{description: &description.Description{
			map[string]*description.Field{
				"balance": &description.Field{
					Kind: fieldtype.Int,
					Min:  dlit.MustNew(200),
					Max:  dlit.MustNew(300),
				},
				"income": &description.Field{
					Kind: fieldtype.Int,
					Min:  dlit.MustNew(300),
					Max:  dlit.MustNew(510),
				},
			},
		},
			stage:       1,
			minNumRules: 18,
			maxNumRules: 20,
			min:         dlit.MustNew(771),
			max:         dlit.MustNew(808),
			mid:         dlit.MustNew(787),
			maxDP:       0,
		},
		{description: &description.Description{
			map[string]*description.Field{
				"balance": &description.Field{
					Kind: fieldtype.Int,
					Min:  dlit.MustNew(200),
					Max:  dlit.MustNew(300),
				},
				"income": &description.Field{
					Kind: fieldtype.Int,
					Min:  dlit.MustNew(590),
					Max:  dlit.MustNew(510),
				},
			},
		},
			stage:       1,
			minNumRules: 2,
			maxNumRules: 2,
			min:         dlit.MustNew(799),
			max:         dlit.MustNew(801),
			mid:         dlit.MustNew(800),
			maxDP:       0,
		},
		{description: &description.Description{
			map[string]*description.Field{
				"balance": &description.Field{
					Kind:  fieldtype.Float,
					Min:   dlit.MustNew(200),
					Max:   dlit.MustNew(300),
					MaxDP: 0,
				},
				"income": &description.Field{
					Kind:  fieldtype.Float,
					Min:   dlit.MustNew(597.924),
					Max:   dlit.MustNew(505),
					MaxDP: 3,
				},
			},
		},
			stage:       1,
			minNumRules: 18,
			maxNumRules: 20,
			min:         dlit.MustNew(799),
			max:         dlit.MustNew(801),
			mid:         dlit.MustNew(800),
			maxDP:       3,
		},
		{description: &description.Description{
			map[string]*description.Field{
				"balance": &description.Field{
					Kind: fieldtype.Int,
					Min:  dlit.MustNew(200),
					Max:  dlit.MustNew(300),
				},
				"income": &description.Field{
					Kind: fieldtype.Int,
					Min:  dlit.MustNew(300),
					Max:  dlit.MustNew(700),
				},
			},
		},
			stage:       2,
			minNumRules: 18,
			maxNumRules: 20,
			min:         dlit.MustNew(778),
			max:         dlit.MustNew(824),
			mid:         dlit.MustNew(800),
			maxDP:       0,
		},
	}
	complyFunc := func(r Rule) error {
		x, ok := r.(*AddGEF)
		if !ok {
			return fmt.Errorf("wrong type: %T (%s)", r, r)
		}
		if x.fieldA != "balance" || x.fieldB != "income" {
			return fmt.Errorf("fields aren't correct for rule: %s", r)
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

/**************************
 *  Benchmarks
 **************************/

func BenchmarkAddGEFIsTrue(b *testing.B) {
	record := map[string]*dlit.Literal{
		"band":   dlit.MustNew(23),
		"income": dlit.MustNew(1024),
		"cost":   dlit.MustNew(890),
	}
	r := NewAddGEF("cost", "income", dlit.MustNew(900.23))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = r.IsTrue(record)
	}
}
