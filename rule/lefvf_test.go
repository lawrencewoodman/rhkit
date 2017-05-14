package rule

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
	"reflect"
	"testing"
)

func TestLEFVFString(t *testing.T) {
	field := "income"
	value := 8.93
	want := "income <= 8.93"
	r := NewLEFVF(field, value)
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestLEFVFIsTrue(t *testing.T) {
	cases := []struct {
		field string
		value float64
		want  bool
	}{
		{"income", 19, true},
		{"income", 19.12, true},
		{"income", 20, true},
		{"income", -20, false},
		{"income", 18.34, false},
		{"flow", 124.564, true},
		{"flow", 124.565, true},
		{"flow", 124.563, false},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"cost":   dlit.MustNew(20),
		"flow":   dlit.MustNew(124.564),
	}
	for _, c := range cases {
		r := NewLEFVF(c.field, c.value)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) (rule: %s) err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestLEFVFIsTrue_errors(t *testing.T) {
	cases := []struct {
		field   string
		value   float64
		wantErr error
	}{
		{field: "fred",
			value:   7.894,
			wantErr: InvalidRuleError{Rule: NewLEFVF("fred", 7.894)},
		},
		{field: "band",
			value:   7.894,
			wantErr: IncompatibleTypesRuleError{Rule: NewLEFVF("band", 7.894)},
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"flow":   dlit.MustNew(124.564),
		"band":   dlit.NewString("alpha"),
	}
	for _, c := range cases {
		r := NewLEFVF(c.field, c.value)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}

func TestLEFVFGetFields(t *testing.T) {
	r := NewLEFVF("income", 5.5)
	want := []string{"income"}
	got := r.GetFields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetFields() got: %s, want: %s", got, want)
	}
}

func TestLEFVFTweak(t *testing.T) {
	field := "income"
	value := float64(800)
	rule := NewLEFVF(field, value)
	cases := []struct {
		description *description.Description
		stage       int
		want        []Rule
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
			stage: 1,
			want: []Rule{
				NewLEFVF(field, float64(755)),
				NewLEFVF(field, float64(760)),
				NewLEFVF(field, float64(765)),
				NewLEFVF(field, float64(770)),
				NewLEFVF(field, float64(775)),
				NewLEFVF(field, float64(780)),
				NewLEFVF(field, float64(785)),
				NewLEFVF(field, float64(790)),
				NewLEFVF(field, float64(795)),
				NewLEFVF(field, float64(805)),
				NewLEFVF(field, float64(810)),
				NewLEFVF(field, float64(815)),
				NewLEFVF(field, float64(820)),
				NewLEFVF(field, float64(825)),
				NewLEFVF(field, float64(830)),
				NewLEFVF(field, float64(835)),
				NewLEFVF(field, float64(840)),
				NewLEFVF(field, float64(845)),
			},
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
			stage: 1,
			want: []Rule{
				NewLEFVF(field, float64(791.55)),
				NewLEFVF(field, float64(793.1)),
				NewLEFVF(field, float64(794.65)),
				NewLEFVF(field, float64(796.2)),
				NewLEFVF(field, float64(797.75)),
				NewLEFVF(field, float64(799.3)),
				NewLEFVF(field, float64(800.85)),
				NewLEFVF(field, float64(802.4)),
				NewLEFVF(field, float64(803.95)),
				NewLEFVF(field, float64(805.5)),
				NewLEFVF(field, float64(807.05)),
				NewLEFVF(field, float64(808.6)),
				NewLEFVF(field, float64(810.15)),
				NewLEFVF(field, float64(811.7)),
				NewLEFVF(field, float64(813.25)),
				NewLEFVF(field, float64(814.8)),
				NewLEFVF(field, float64(816.35)),
				NewLEFVF(field, float64(817.9)),
				NewLEFVF(field, float64(819.45)),
			},
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
			stage: 1,
			want: []Rule{
				NewLEFVF(field, float64(771.05)),
				NewLEFVF(field, float64(773.1)),
				NewLEFVF(field, float64(775.15)),
				NewLEFVF(field, float64(777.2)),
				NewLEFVF(field, float64(779.25)),
				NewLEFVF(field, float64(781.3)),
				NewLEFVF(field, float64(783.35)),
				NewLEFVF(field, float64(785.4)),
				NewLEFVF(field, float64(787.45)),
				NewLEFVF(field, float64(789.5)),
				NewLEFVF(field, float64(791.55)),
				NewLEFVF(field, float64(793.6)),
				NewLEFVF(field, float64(795.65)),
				NewLEFVF(field, float64(797.7)),
				NewLEFVF(field, float64(799.75)),
				NewLEFVF(field, float64(801.8)),
				NewLEFVF(field, float64(803.85)),
				NewLEFVF(field, float64(805.9)),
				NewLEFVF(field, float64(807.95)),
			},
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
			stage: 1,
			want:  []Rule{},
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
			stage: 2,
			want: []Rule{
				NewLEFVF(field, float64(777.5)),
				NewLEFVF(field, float64(780)),
				NewLEFVF(field, float64(782.5)),
				NewLEFVF(field, float64(785)),
				NewLEFVF(field, float64(787.5)),
				NewLEFVF(field, float64(790)),
				NewLEFVF(field, float64(792.5)),
				NewLEFVF(field, float64(795)),
				NewLEFVF(field, float64(797.5)),
				NewLEFVF(field, float64(802.5)),
				NewLEFVF(field, float64(805)),
				NewLEFVF(field, float64(807.5)),
				NewLEFVF(field, float64(810)),
				NewLEFVF(field, float64(812.5)),
				NewLEFVF(field, float64(815)),
				NewLEFVF(field, float64(817.5)),
				NewLEFVF(field, float64(820)),
				NewLEFVF(field, float64(822.5)),
			},
		},
	}
	for i, c := range cases {
		got := rule.Tweak(c.description, c.stage)
		if err := checkRulesMatch(got, c.want); err != nil {
			t.Errorf("(%d) Tweak: %s, got: %s, want: %s", i, err, got, c.want)
		}
	}
}

func TestLEFVFOverlaps(t *testing.T) {
	cases := []struct {
		ruleA *LEFVF
		ruleB Rule
		want  bool
	}{
		{ruleA: NewLEFVF("band", 7.3),
			ruleB: NewLEFVF("band", 6.5),
			want:  true,
		},
		{ruleA: NewLEFVF("band", 7.3),
			ruleB: NewLEFVF("rate", 6.5),
			want:  false,
		},
		{ruleA: NewLEFVF("band", 7.3),
			ruleB: NewGEFVF("band", 6.5),
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
