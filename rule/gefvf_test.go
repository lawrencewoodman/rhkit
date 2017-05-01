package rule

import (
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
				NewGEFVF(field, float64(755)),
				NewGEFVF(field, float64(760)),
				NewGEFVF(field, float64(765)),
				NewGEFVF(field, float64(770)),
				NewGEFVF(field, float64(775)),
				NewGEFVF(field, float64(780)),
				NewGEFVF(field, float64(785)),
				NewGEFVF(field, float64(790)),
				NewGEFVF(field, float64(795)),
				NewGEFVF(field, float64(805)),
				NewGEFVF(field, float64(810)),
				NewGEFVF(field, float64(815)),
				NewGEFVF(field, float64(820)),
				NewGEFVF(field, float64(825)),
				NewGEFVF(field, float64(830)),
				NewGEFVF(field, float64(835)),
				NewGEFVF(field, float64(840)),
				NewGEFVF(field, float64(845)),
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
				NewGEFVF(field, float64(791.55)),
				NewGEFVF(field, float64(791.6)),
				NewGEFVF(field, float64(792)),
				NewGEFVF(field, float64(793)),
				NewGEFVF(field, float64(793.1)),
				NewGEFVF(field, float64(794.6)),
				NewGEFVF(field, float64(794.65)),
				NewGEFVF(field, float64(795)),
				NewGEFVF(field, float64(796)),
				NewGEFVF(field, float64(796.2)),
				NewGEFVF(field, float64(797.7)),
				NewGEFVF(field, float64(797.75)),
				NewGEFVF(field, float64(798)),
				NewGEFVF(field, float64(799)),
				NewGEFVF(field, float64(799.3)),
				NewGEFVF(field, float64(800.8)),
				NewGEFVF(field, float64(800.85)),
				NewGEFVF(field, float64(801)),
				NewGEFVF(field, float64(802)),
				NewGEFVF(field, float64(802.4)),
				NewGEFVF(field, float64(803.9)),
				NewGEFVF(field, float64(803.95)),
				NewGEFVF(field, float64(804)),
				NewGEFVF(field, float64(805)),
				NewGEFVF(field, float64(805.5)),
				NewGEFVF(field, float64(807)),
				NewGEFVF(field, float64(807.05)),
				NewGEFVF(field, float64(808.6)),
				NewGEFVF(field, float64(809)),
				NewGEFVF(field, float64(810)),
				NewGEFVF(field, float64(810.1)),
				NewGEFVF(field, float64(810.15)),
				NewGEFVF(field, float64(811.7)),
				NewGEFVF(field, float64(812)),
				NewGEFVF(field, float64(813)),
				NewGEFVF(field, float64(813.2)),
				NewGEFVF(field, float64(813.25)),
				NewGEFVF(field, float64(814.8)),
				NewGEFVF(field, float64(815)),
				NewGEFVF(field, float64(816)),
				NewGEFVF(field, float64(816.3)),
				NewGEFVF(field, float64(816.35)),
				NewGEFVF(field, float64(817.9)),
				NewGEFVF(field, float64(818)),
				NewGEFVF(field, float64(819)),
				NewGEFVF(field, float64(819.4)),
				NewGEFVF(field, float64(819.45)),
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
				NewGEFVF(field, float64(771)),
				NewGEFVF(field, float64(771.05)),
				NewGEFVF(field, float64(771.1)),
				NewGEFVF(field, float64(773)),
				NewGEFVF(field, float64(773.1)),
				NewGEFVF(field, float64(775)),
				NewGEFVF(field, float64(775.1)),
				NewGEFVF(field, float64(775.15)),
				NewGEFVF(field, float64(777)),
				NewGEFVF(field, float64(777.2)),
				NewGEFVF(field, float64(779)),
				NewGEFVF(field, float64(779.2)),
				NewGEFVF(field, float64(779.25)),
				NewGEFVF(field, float64(781)),
				NewGEFVF(field, float64(781.3)),
				NewGEFVF(field, float64(783)),
				NewGEFVF(field, float64(783.3)),
				NewGEFVF(field, float64(783.35)),
				NewGEFVF(field, float64(785)),
				NewGEFVF(field, float64(785.4)),
				NewGEFVF(field, float64(787)),
				NewGEFVF(field, float64(787.4)),
				NewGEFVF(field, float64(787.45)),
				NewGEFVF(field, float64(789)),
				NewGEFVF(field, float64(789.5)),
				NewGEFVF(field, float64(791.5)),
				NewGEFVF(field, float64(791.55)),
				NewGEFVF(field, float64(792)),
				NewGEFVF(field, float64(793.6)),
				NewGEFVF(field, float64(794)),
				NewGEFVF(field, float64(795.6)),
				NewGEFVF(field, float64(795.65)),
				NewGEFVF(field, float64(796)),
				NewGEFVF(field, float64(797.7)),
				NewGEFVF(field, float64(798)),
				NewGEFVF(field, float64(799.7)),
				NewGEFVF(field, float64(799.75)),
				NewGEFVF(field, float64(801.8)),
				NewGEFVF(field, float64(802)),
				NewGEFVF(field, float64(803.8)),
				NewGEFVF(field, float64(803.85)),
				NewGEFVF(field, float64(804)),
				NewGEFVF(field, float64(805.9)),
				NewGEFVF(field, float64(806)),
				NewGEFVF(field, float64(807.9)),
				NewGEFVF(field, float64(807.95)),
				NewGEFVF(field, float64(808)),
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
				NewGEFVF(field, float64(777.5)),
				NewGEFVF(field, float64(778)),
				NewGEFVF(field, float64(780)),
				NewGEFVF(field, float64(782.5)),
				NewGEFVF(field, float64(783)),
				NewGEFVF(field, float64(785)),
				NewGEFVF(field, float64(787.5)),
				NewGEFVF(field, float64(788)),
				NewGEFVF(field, float64(790)),
				NewGEFVF(field, float64(792.5)),
				NewGEFVF(field, float64(793)),
				NewGEFVF(field, float64(795)),
				NewGEFVF(field, float64(797.5)),
				NewGEFVF(field, float64(798)),
				NewGEFVF(field, float64(802.5)),
				NewGEFVF(field, float64(803)),
				NewGEFVF(field, float64(805)),
				NewGEFVF(field, float64(807.5)),
				NewGEFVF(field, float64(808)),
				NewGEFVF(field, float64(810)),
				NewGEFVF(field, float64(812.5)),
				NewGEFVF(field, float64(813)),
				NewGEFVF(field, float64(815)),
				NewGEFVF(field, float64(817.5)),
				NewGEFVF(field, float64(818)),
				NewGEFVF(field, float64(820)),
				NewGEFVF(field, float64(822.5)),
				NewGEFVF(field, float64(823)),
			},
		},
	}
	for _, c := range cases {
		got := rule.Tweak(c.description, c.stage)
		if err := checkRulesMatch(got, c.want); err != nil {
			t.Errorf("Tweak: %s, got: %s, want: %s", err, got, c.want)
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
