package rule

import (
	"github.com/lawrencewoodman/dlit"
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
		min   *dlit.Literal
		max   *dlit.Literal
		stage int
		want  []Rule
	}{
		{min: dlit.MustNew(500),
			max:   dlit.MustNew(1000),
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
		{min: dlit.MustNew(790),
			max:   dlit.MustNew(1000),
			stage: 1,
			want: []Rule{
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
				NewLEFVI(field, int64(811)),
				NewLEFVI(field, int64(813)),
				NewLEFVI(field, int64(815)),
				NewLEFVI(field, int64(817)),
				NewLEFVI(field, int64(819)),
			},
		},
		{min: dlit.MustNew(500),
			max:   dlit.MustNew(810),
			stage: 1,
			want: []Rule{
				NewLEFVI(field, int64(772)),
				NewLEFVI(field, int64(775)),
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
			},
		},
		{min: dlit.MustNew(798),
			max:   dlit.MustNew(805),
			stage: 1,
			want:  []Rule{},
		},
		{min: dlit.MustNew(500),
			max:   dlit.MustNew(1000),
			stage: 2,
			want: []Rule{
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
				NewLEFVI(field, int64(811)),
				NewLEFVI(field, int64(813)),
				NewLEFVI(field, int64(815)),
				NewLEFVI(field, int64(817)),
				NewLEFVI(field, int64(819)),
				NewLEFVI(field, int64(821)),
				NewLEFVI(field, int64(823)),
			},
		},
	}
	for _, c := range cases {
		got := rule.Tweak(c.min, c.max, 0, c.stage)
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
