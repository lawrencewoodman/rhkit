package rule

import (
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
	"reflect"
	"testing"
)

func TestGEFVIString(t *testing.T) {
	field := "income"
	value := int64(893)
	want := "income >= 893"
	r := NewGEFVI(field, value)
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestGEFVIValue(t *testing.T) {
	field := "income"
	value := int64(893)
	r := NewGEFVI(field, value)
	got := r.Value()
	if got.String() != "893" {
		t.Errorf("Value() got: %s, want: %f", got, value)
	}
}

func TestGEFVIIsTrue(t *testing.T) {
	cases := []struct {
		field string
		value int64
		want  bool
	}{
		{"income", 19, true},
		{"income", 18, true},
		{"income", 20, false},
		{"income", -20, true},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"cost":   dlit.MustNew(20),
		"flow":   dlit.MustNew(124.564),
	}
	for _, c := range cases {
		r := NewGEFVI(c.field, c.value)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) (rule: %s) err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestGEFVIIsTrue_errors(t *testing.T) {
	cases := []struct {
		field   string
		value   int64
		wantErr error
	}{
		{field: "fred",
			value:   7,
			wantErr: InvalidRuleError{Rule: NewGEFVI("fred", 7)},
		},
		{field: "band",
			value:   7894,
			wantErr: IncompatibleTypesRuleError{Rule: NewGEFVI("band", 7894)},
		},
		{field: "flow",
			value:   7894,
			wantErr: IncompatibleTypesRuleError{Rule: NewGEFVI("flow", 7894)},
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"flow":   dlit.MustNew(124.564),
		"band":   dlit.NewString("alpha"),
	}
	for _, c := range cases {
		r := NewGEFVI(c.field, c.value)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}

func TestGEFVITweak(t *testing.T) {
	field := "income"
	value := int64(800)
	rule := NewGEFVI(field, value)
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
					Kind: fieldtype.Int,
					Min:  dlit.MustNew(500),
					Max:  dlit.MustNew(1000),
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
					Kind: fieldtype.Int,
					Min:  dlit.MustNew(790),
					Max:  dlit.MustNew(1000),
				},
			},
		},
			stage:       1,
			minNumRules: 18,
			maxNumRules: 20,
			min:         dlit.MustNew(792),
			max:         dlit.MustNew(819),
			mid:         dlit.MustNew(805),
			maxDP:       0,
		},
		{description: &description.Description{
			map[string]*description.Field{
				"income": &description.Field{
					Kind: fieldtype.Int,
					Min:  dlit.MustNew(500),
					Max:  dlit.MustNew(810),
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
				"income": &description.Field{
					Kind: fieldtype.Int,
					Min:  dlit.MustNew(798),
					Max:  dlit.MustNew(805),
				},
			},
		},
			stage:       1,
			minNumRules: 0,
			maxNumRules: 0,
			min:         dlit.MustNew(798),
			max:         dlit.MustNew(800),
			mid:         dlit.MustNew(805),
			maxDP:       0,
		},
		{description: &description.Description{
			map[string]*description.Field{
				"income": &description.Field{
					Kind: fieldtype.Int,
					Min:  dlit.MustNew(500),
					Max:  dlit.MustNew(1000),
				},
			},
		},
			stage:       2,
			minNumRules: 18,
			maxNumRules: 20,
			min:         dlit.MustNew(778),
			max:         dlit.MustNew(823),
			mid:         dlit.MustNew(798),
			maxDP:       0,
		},
	}
	complyFunc := func(r Rule) error {
		x, ok := r.(*GEFVI)
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

func TestGEFVIGetFields(t *testing.T) {
	r := NewGEFVI("income", 5)
	want := []string{"income"}
	got := r.GetFields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetFields() got: %s, want: %s", got, want)
	}
}

func TestGEFVIOverlaps(t *testing.T) {
	cases := []struct {
		ruleA *GEFVI
		ruleB Rule
		want  bool
	}{
		{ruleA: NewGEFVI("band", 7),
			ruleB: NewGEFVI("band", 6),
			want:  true,
		},
		{ruleA: NewGEFVI("band", 7),
			ruleB: NewGEFVI("rate", 6),
			want:  false,
		},
		{ruleA: NewGEFVI("band", 7),
			ruleB: NewLEFVI("band", 6),
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
