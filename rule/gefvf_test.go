package rule

import (
	"github.com/lawrencewoodman/dlit"
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
		min   *dlit.Literal
		max   *dlit.Literal
		maxDP int
		stage int
		want  []Rule
	}{
		{min: dlit.MustNew(500),
			max:   dlit.MustNew(1000),
			maxDP: 2,
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
		{min: dlit.MustNew(790),
			max:   dlit.MustNew(1000),
			maxDP: 2,
			stage: 1,
			want: []Rule{
				NewGEFVF(field, float64(791.6)),
				NewGEFVF(field, float64(793.7)),
				NewGEFVF(field, float64(795.8)),
				NewGEFVF(field, float64(797.9)),
				NewGEFVF(field, float64(802.1)),
				NewGEFVF(field, float64(804.2)),
				NewGEFVF(field, float64(806.3)),
				NewGEFVF(field, float64(808.4)),
				NewGEFVF(field, float64(810.5)),
				NewGEFVF(field, float64(812.6)),
				NewGEFVF(field, float64(814.7)),
				NewGEFVF(field, float64(816.8)),
				NewGEFVF(field, float64(818.9)),
			},
		},
		{min: dlit.MustNew(500),
			max:   dlit.MustNew(810),
			maxDP: 2,
			stage: 1,
			want: []Rule{
				NewGEFVF(field, float64(772.1)),
				NewGEFVF(field, float64(775.2)),
				NewGEFVF(field, float64(778.3)),
				NewGEFVF(field, float64(781.4)),
				NewGEFVF(field, float64(784.5)),
				NewGEFVF(field, float64(787.6)),
				NewGEFVF(field, float64(790.7)),
				NewGEFVF(field, float64(793.8)),
				NewGEFVF(field, float64(796.9)),
				NewGEFVF(field, float64(803.1)),
				NewGEFVF(field, float64(806.2)),
				NewGEFVF(field, float64(809.3)),
			},
		},
		{min: dlit.MustNew(799),
			max:   dlit.MustNew(801),
			maxDP: 0,
			stage: 1,
			want:  []Rule{},
		},
		{min: dlit.MustNew(500),
			max:   dlit.MustNew(1000),
			maxDP: 2,
			stage: 2,
			want: []Rule{
				NewGEFVF(field, float64(777.5)),
				NewGEFVF(field, float64(780)),
				NewGEFVF(field, float64(782.5)),
				NewGEFVF(field, float64(785)),
				NewGEFVF(field, float64(787.5)),
				NewGEFVF(field, float64(790)),
				NewGEFVF(field, float64(792.5)),
				NewGEFVF(field, float64(795)),
				NewGEFVF(field, float64(797.5)),
				NewGEFVF(field, float64(802.5)),
				NewGEFVF(field, float64(805)),
				NewGEFVF(field, float64(807.5)),
				NewGEFVF(field, float64(810)),
				NewGEFVF(field, float64(812.5)),
				NewGEFVF(field, float64(815)),
				NewGEFVF(field, float64(817.5)),
				NewGEFVF(field, float64(820)),
				NewGEFVF(field, float64(822.5)),
			},
		},
	}
	for _, c := range cases {
		got := rule.Tweak(c.min, c.max, c.maxDP, c.stage)
		if err := checkRulesMatch(got, c.want); err != nil {
			t.Errorf("Tweak: %s, got: %s, want: %s", err, got, c.want)
		}
	}
}
