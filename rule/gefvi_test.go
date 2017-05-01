package rule

import (
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
		want        []Rule
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
			stage: 1,
			want: []Rule{
				NewGEFVI(field, int64(755)),
				NewGEFVI(field, int64(760)),
				NewGEFVI(field, int64(765)),
				NewGEFVI(field, int64(770)),
				NewGEFVI(field, int64(775)),
				NewGEFVI(field, int64(780)),
				NewGEFVI(field, int64(785)),
				NewGEFVI(field, int64(790)),
				NewGEFVI(field, int64(795)),
				NewGEFVI(field, int64(805)),
				NewGEFVI(field, int64(810)),
				NewGEFVI(field, int64(815)),
				NewGEFVI(field, int64(820)),
				NewGEFVI(field, int64(825)),
				NewGEFVI(field, int64(830)),
				NewGEFVI(field, int64(835)),
				NewGEFVI(field, int64(840)),
				NewGEFVI(field, int64(845)),
			},
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
			stage: 1,
			want: []Rule{
				NewGEFVI(field, int64(792)),
				NewGEFVI(field, int64(794)),
				NewGEFVI(field, int64(796)),
				NewGEFVI(field, int64(798)),
				NewGEFVI(field, int64(802)),
				NewGEFVI(field, int64(804)),
				NewGEFVI(field, int64(806)),
				NewGEFVI(field, int64(808)),
				NewGEFVI(field, int64(810)),
				NewGEFVI(field, int64(812)),
				NewGEFVI(field, int64(814)),
				NewGEFVI(field, int64(816)),
				NewGEFVI(field, int64(818)),
				NewGEFVI(field, int64(820)),
			},
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
			stage: 1,
			want: []Rule{
				NewGEFVI(field, int64(771)),
				NewGEFVI(field, int64(773)),
				NewGEFVI(field, int64(775)),
				NewGEFVI(field, int64(777)),
				NewGEFVI(field, int64(779)),
				NewGEFVI(field, int64(781)),
				NewGEFVI(field, int64(783)),
				NewGEFVI(field, int64(785)),
				NewGEFVI(field, int64(787)),
				NewGEFVI(field, int64(789)),
				NewGEFVI(field, int64(791)),
				NewGEFVI(field, int64(793)),
				NewGEFVI(field, int64(795)),
				NewGEFVI(field, int64(797)),
				NewGEFVI(field, int64(799)),
				NewGEFVI(field, int64(801)),
				NewGEFVI(field, int64(803)),
				NewGEFVI(field, int64(805)),
				NewGEFVI(field, int64(807)),
				NewGEFVI(field, int64(809)),
			},
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
			stage: 1,
			want:  []Rule{},
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
			stage: 2,
			want: []Rule{
				NewGEFVI(field, int64(778)),
				NewGEFVI(field, int64(781)),
				NewGEFVI(field, int64(784)),
				NewGEFVI(field, int64(787)),
				NewGEFVI(field, int64(790)),
				NewGEFVI(field, int64(793)),
				NewGEFVI(field, int64(796)),
				NewGEFVI(field, int64(799)),
				NewGEFVI(field, int64(802)),
				NewGEFVI(field, int64(805)),
				NewGEFVI(field, int64(808)),
				NewGEFVI(field, int64(811)),
				NewGEFVI(field, int64(814)),
				NewGEFVI(field, int64(817)),
				NewGEFVI(field, int64(820)),
				NewGEFVI(field, int64(823)),
			},
		},
	}
	for _, c := range cases {
		got := rule.Tweak(c.description, c.stage)
		if err := checkRulesMatch(got, c.want); err != nil {
			t.Errorf("Tweak: %s, got: %s", err, got)
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
