package rule

import (
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
	"reflect"
	"testing"
)

func TestMulGEFString(t *testing.T) {
	fieldA := "income"
	fieldB := "balance"
	value := dlit.MustNew(8.93)
	want := "income * balance >= 8.93"
	r := NewMulGEF(fieldA, fieldB, value)
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestMulGEFValue(t *testing.T) {
	fieldA := "income"
	fieldB := "balance"
	value := dlit.MustNew(8.93)
	r := NewMulGEF(fieldA, fieldB, value)
	got := r.Value()
	if got.String() != value.String() {
		t.Errorf("String() got: %s, want: %s", got, value)
	}
}

func TestMulGEFIsTrue(t *testing.T) {
	cases := []struct {
		fieldA string
		fieldB string
		value  *dlit.Literal
		want   bool
	}{
		{"income", "balance", dlit.MustNew(61), false},
		{"income", "balance", dlit.MustNew(60.12), false},
		{"income", "balance", dlit.MustNew(60), true},
		{"income", "balance", dlit.MustNew(-60), true},
		{"income", "balance", dlit.MustNew(59.89), true},
		{"flow", "cost", dlit.MustNew(2491.2), true},
		{"flow", "cost", dlit.MustNew(2491.21), false},
		{"flow", "cost", dlit.MustNew(2491.19), true},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(4),
		"balance": dlit.MustNew(15),
		"cost":    dlit.MustNew(20),
		"flow":    dlit.MustNew(124.56),
	}
	for _, c := range cases {
		r := NewMulGEF(c.fieldA, c.fieldB, c.value)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) (rule: %s) err: %v", r, err)
		} else if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestMulGEFIsTrue_errors(t *testing.T) {
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
				Rule: NewMulGEF("fred", "flow", dlit.MustNew(7.894)),
			},
		},
		{fieldA: "flow",
			fieldB: "fred",
			value:  dlit.MustNew(7.894),
			wantErr: InvalidRuleError{
				Rule: NewMulGEF("flow", "fred", dlit.MustNew(7.894)),
			},
		},
		{fieldA: "band",
			fieldB: "flow",
			value:  dlit.MustNew(7.894),
			wantErr: IncompatibleTypesRuleError{
				Rule: NewMulGEF("band", "flow", dlit.MustNew(7.894)),
			},
		},
		{fieldA: "flow",
			fieldB: "band",
			value:  dlit.MustNew(7.894),
			wantErr: IncompatibleTypesRuleError{
				Rule: NewMulGEF("flow", "band", dlit.MustNew(7.894)),
			},
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"flow":   dlit.MustNew(124.564),
		"band":   dlit.NewString("alpha"),
	}
	for _, c := range cases {
		r := NewMulGEF(c.fieldA, c.fieldB, c.value)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}

func TestMulGEFGetFields(t *testing.T) {
	r := NewMulGEF("income", "cost", dlit.MustNew(5.5))
	want := []string{"income", "cost"}
	got := r.GetFields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetFields() got: %s, want: %s", got, want)
	}
}

func TestMulGEFOverlaps(t *testing.T) {
	cases := []struct {
		ruleA *MulGEF
		ruleB Rule
		want  bool
	}{
		{ruleA: NewMulGEF("band", "cost", dlit.MustNew(7.3)),
			ruleB: NewMulGEF("band", "cost", dlit.MustNew(6.5)),
			want:  true,
		},
		{ruleA: NewMulGEF("band", "balance", dlit.MustNew(7.3)),
			ruleB: NewMulGEF("rate", "balance", dlit.MustNew(6.5)),
			want:  false,
		},
		{ruleA: NewMulGEF("band", "balance", dlit.MustNew(7.3)),
			ruleB: NewMulGEF("band", "rate", dlit.MustNew(6.5)),
			want:  false,
		},
		{ruleA: NewMulGEF("band", "cost", dlit.MustNew(7.3)),
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

func TestMulGEFTweak(t *testing.T) {
	fieldA := "income"
	fieldB := "balance"
	value := int64(156250)
	rule := NewMulGEF(fieldA, fieldB, dlit.MustNew(value))
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
					Kind: fieldtype.Number,
					Min:  dlit.MustNew(250),
					Max:  dlit.MustNew(500),
				},
				"income": &description.Field{
					Kind: fieldtype.Number,
					Min:  dlit.MustNew(250),
					Max:  dlit.MustNew(500),
				},
			},
		},
			stage:       1,
			minNumRules: 18,
			maxNumRules: 20,
			min:         dlit.MustNew(139375),
			max:         dlit.MustNew(173125),
			mid:         dlit.MustNew(152500),
			maxDP:       0,
		},
		{description: &description.Description{
			map[string]*description.Field{
				"balance": &description.Field{
					Kind: fieldtype.Number,
					Min:  dlit.MustNew(250),
					Max:  dlit.MustNew(300),
				},
				"income": &description.Field{
					Kind: fieldtype.Number,
					Min:  dlit.MustNew(540),
					Max:  dlit.MustNew(700),
				},
			},
		},
			stage:       1,
			minNumRules: 18,
			maxNumRules: 20,
			min:         dlit.MustNew(149500),
			max:         dlit.MustNew(163000),
			mid:         dlit.MustNew(154750),
			maxDP:       0,
		},
		{description: &description.Description{
			map[string]*description.Field{
				"balance": &description.Field{
					Kind: fieldtype.Number,
					Min:  dlit.MustNew(200),
					Max:  dlit.MustNew(300),
				},
				"income": &description.Field{
					Kind: fieldtype.Number,
					Min:  dlit.MustNew(300),
					Max:  dlit.MustNew(510),
				},
			},
		},
			stage:       1,
			minNumRules: 18,
			maxNumRules: 20,
			min:         dlit.MustNew(147253),
			max:         dlit.MustNew(152698),
			mid:         dlit.MustNew(149673),
			maxDP:       0,
		},
		{description: &description.Description{
			map[string]*description.Field{
				"balance": &description.Field{
					Kind:  fieldtype.Number,
					Min:   dlit.MustNew(200),
					Max:   dlit.MustNew(300),
					MaxDP: 0,
				},
				"income": &description.Field{
					Kind:  fieldtype.Number,
					Min:   dlit.MustNew(300.924),
					Max:   dlit.MustNew(505),
					MaxDP: 3,
				},
			},
		},
			stage:       1,
			minNumRules: 18,
			maxNumRules: 20,
			min:         dlit.MustNew(147337),
			max:         dlit.MustNew(151281),
			mid:         dlit.MustNew(149090),
			maxDP:       3,
		},
		{description: &description.Description{
			map[string]*description.Field{
				"balance": &description.Field{
					Kind: fieldtype.Number,
					Min:  dlit.MustNew(200),
					Max:  dlit.MustNew(300),
				},
				"income": &description.Field{
					Kind: fieldtype.Number,
					Min:  dlit.MustNew(300),
					Max:  dlit.MustNew(700),
				},
			},
		},
			stage:       2,
			minNumRules: 18,
			maxNumRules: 20,
			min:         dlit.MustNew(149500),
			max:         dlit.MustNew(163000),
			mid:         dlit.MustNew(154750),
			maxDP:       0,
		},
	}
	complyFunc := func(r Rule) error {
		x, ok := r.(*MulGEF)
		if !ok {
			return fmt.Errorf("wrong type: %T (%s)", r, r)
		}
		if x.fieldA != "income" || x.fieldB != "balance" {
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

func BenchmarkMulGEFIsTrue(b *testing.B) {
	record := map[string]*dlit.Literal{
		"band":   dlit.MustNew(23),
		"income": dlit.MustNew(1024),
		"cost":   dlit.MustNew(890),
	}
	r := NewMulGEF("cost", "income", dlit.MustNew(900.23))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = r.IsTrue(record)
	}
}
