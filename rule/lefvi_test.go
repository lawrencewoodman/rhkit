package rule

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
	"reflect"
	"testing"
)

func TestLEFVIString(t *testing.T) {
	field := "income"
	value := int64(893)
	want := "income <= 893"
	r := NewLEFVI(field, value)
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestLEFVIIsTrue(t *testing.T) {
	cases := []struct {
		field string
		value int64
		want  bool
	}{
		{"income", 19, true},
		{"income", 20, true},
		{"income", -20, false},
		{"income", 18, false},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"cost":   dlit.MustNew(20),
		"flow":   dlit.MustNew(124.564),
	}
	for _, c := range cases {
		r := NewLEFVI(c.field, c.value)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) (rule: %s) err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestLEFVIIsTrue_errors(t *testing.T) {
	cases := []struct {
		field   string
		value   int64
		wantErr error
	}{
		{field: "fred",
			value:   7894,
			wantErr: InvalidRuleError{Rule: NewLEFVI("fred", 7894)},
		},
		{field: "band",
			value:   7894,
			wantErr: IncompatibleTypesRuleError{Rule: NewLEFVI("band", 7894)},
		},
		{field: "flow",
			value:   7894,
			wantErr: IncompatibleTypesRuleError{Rule: NewLEFVI("flow", 7894)},
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"flow":   dlit.MustNew(124.564),
		"band":   dlit.NewString("alpha"),
	}
	for _, c := range cases {
		r := NewLEFVI(c.field, c.value)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}

func TestLEFVITweak(t *testing.T) {
	field := "income"
	value := int64(800)
	rule := NewLEFVI(field, value)
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
				NewLEFVI(field, int64(755)),
				NewLEFVI(field, int64(760)),
				NewLEFVI(field, int64(765)),
				NewLEFVI(field, int64(770)),
				NewLEFVI(field, int64(775)),
				NewLEFVI(field, int64(780)),
				NewLEFVI(field, int64(785)),
				NewLEFVI(field, int64(790)),
				NewLEFVI(field, int64(795)),
				NewLEFVI(field, int64(805)),
				NewLEFVI(field, int64(810)),
				NewLEFVI(field, int64(815)),
				NewLEFVI(field, int64(820)),
				NewLEFVI(field, int64(825)),
				NewLEFVI(field, int64(830)),
				NewLEFVI(field, int64(835)),
				NewLEFVI(field, int64(840)),
				NewLEFVI(field, int64(845)),
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
				NewLEFVI(field, int64(792)),
				NewLEFVI(field, int64(794)),
				NewLEFVI(field, int64(796)),
				NewLEFVI(field, int64(798)),
				NewLEFVI(field, int64(802)),
				NewLEFVI(field, int64(804)),
				NewLEFVI(field, int64(806)),
				NewLEFVI(field, int64(808)),
				NewLEFVI(field, int64(810)),
				NewLEFVI(field, int64(812)),
				NewLEFVI(field, int64(814)),
				NewLEFVI(field, int64(816)),
				NewLEFVI(field, int64(818)),
				NewLEFVI(field, int64(820)),
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
				NewLEFVI(field, int64(771)),
				NewLEFVI(field, int64(773)),
				NewLEFVI(field, int64(775)),
				NewLEFVI(field, int64(777)),
				NewLEFVI(field, int64(779)),
				NewLEFVI(field, int64(781)),
				NewLEFVI(field, int64(783)),
				NewLEFVI(field, int64(785)),
				NewLEFVI(field, int64(787)),
				NewLEFVI(field, int64(789)),
				NewLEFVI(field, int64(791)),
				NewLEFVI(field, int64(793)),
				NewLEFVI(field, int64(795)),
				NewLEFVI(field, int64(797)),
				NewLEFVI(field, int64(799)),
				NewLEFVI(field, int64(801)),
				NewLEFVI(field, int64(803)),
				NewLEFVI(field, int64(805)),
				NewLEFVI(field, int64(807)),
				NewLEFVI(field, int64(809)),
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
				NewLEFVI(field, int64(778)),
				NewLEFVI(field, int64(781)),
				NewLEFVI(field, int64(784)),
				NewLEFVI(field, int64(787)),
				NewLEFVI(field, int64(790)),
				NewLEFVI(field, int64(793)),
				NewLEFVI(field, int64(796)),
				NewLEFVI(field, int64(799)),
				NewLEFVI(field, int64(802)),
				NewLEFVI(field, int64(805)),
				NewLEFVI(field, int64(808)),
				NewLEFVI(field, int64(811)),
				NewLEFVI(field, int64(814)),
				NewLEFVI(field, int64(817)),
				NewLEFVI(field, int64(820)),
				NewLEFVI(field, int64(823)),
			},
		},
	}
	complexity := 10
	for _, c := range cases {
		got := rule.Tweak(c.description, complexity, c.stage)
		if err := checkRulesMatch(got, c.want); err != nil {
			t.Errorf("Tweak: %s, got: %s", err, got)
		}
	}
}

func TestLEFVIGetFields(t *testing.T) {
	r := NewLEFVI("income", 5)
	want := []string{"income"}
	got := r.GetFields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetFields() got: %s, want: %s", got, want)
	}
}

func TestLEFVIOverlaps(t *testing.T) {
	cases := []struct {
		ruleA *LEFVI
		ruleB Rule
		want  bool
	}{
		{ruleA: NewLEFVI("band", 7),
			ruleB: NewLEFVI("band", 6),
			want:  true,
		},
		{ruleA: NewLEFVI("band", 7),
			ruleB: NewLEFVI("rate", 6),
			want:  false,
		},
		{ruleA: NewLEFVI("band", 7),
			ruleB: NewGEFVI("band", 6),
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
