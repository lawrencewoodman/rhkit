package rule

import (
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
	"reflect"
	"testing"
)

func TestGEFVFString(t *testing.T) {
	field := "income"
	value := 8.93
	want := "income >= 8.93"
	r := NewGEFVF(field, value)
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestGEFVFValue(t *testing.T) {
	field := "income"
	value := 8.93
	r := NewGEFVF(field, value)
	got := r.Value()
	if got.String() != "8.93" {
		t.Errorf("Value() got: %s, want: %f", got, value)
	}
}

func TestGEFVFIsTrue(t *testing.T) {
	cases := []struct {
		field string
		value float64
		want  bool
	}{
		{"income", 19, true},
		{"income", 19.12, false},
		{"income", 20, false},
		{"income", -20, true},
		{"income", 18.34, true},
		{"flow", 124.564, true},
		{"flow", 124.565, false},
		{"flow", 124.563, true},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"cost":   dlit.MustNew(20),
		"flow":   dlit.MustNew(124.564),
	}
	for _, c := range cases {
		r := NewGEFVF(c.field, c.value)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) (rule: %s) err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestGEFVFIsTrue_errors(t *testing.T) {
	cases := []struct {
		field   string
		value   float64
		wantErr error
	}{
		{field: "fred",
			value:   7.894,
			wantErr: InvalidRuleError{Rule: NewGEFVF("fred", 7.894)},
		},
		{field: "band",
			value:   7.894,
			wantErr: IncompatibleTypesRuleError{Rule: NewGEFVF("band", 7.894)},
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"flow":   dlit.MustNew(124.564),
		"band":   dlit.NewString("alpha"),
	}
	for _, c := range cases {
		r := NewGEFVF(c.field, c.value)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}

func TestGEFVFGetFields(t *testing.T) {
	r := NewGEFVF("income", 5.5)
	want := []string{"income"}
	got := r.GetFields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetFields() got: %s, want: %s", got, want)
	}
}

func TestGEFVFTweak(t *testing.T) {
	field := "income"
	value := float64(800)
	rule := NewGEFVF(field, value)
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
				"income": &description.Field{
					Kind:  fieldtype.Float,
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
				"income": &description.Field{
					Kind:  fieldtype.Float,
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
				"income": &description.Field{
					Kind:  fieldtype.Float,
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
				"income": &description.Field{
					Kind:  fieldtype.Float,
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
				"income": &description.Field{
					Kind:  fieldtype.Float,
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
		x, ok := r.(*GEFVF)
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

func TestGEFVFOverlaps(t *testing.T) {
	cases := []struct {
		ruleA *GEFVF
		ruleB Rule
		want  bool
	}{
		{ruleA: NewGEFVF("band", 7.3),
			ruleB: NewGEFVF("band", 6.5),
			want:  true,
		},
		{ruleA: NewGEFVF("band", 7.3),
			ruleB: NewGEFVF("rate", 6.5),
			want:  false,
		},
		{ruleA: NewGEFVF("band", 7.3),
			ruleB: NewLEFVF("band", 6.5),
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
