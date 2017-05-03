package rule

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
	"reflect"
	"testing"
)

func TestAddLEFString(t *testing.T) {
	fieldA := "income"
	fieldB := "balance"
	value := dlit.MustNew(8.93)
	want := "income + balance <= 8.93"
	r := NewAddLEF(fieldA, fieldB, value)
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestAddLEFIsTrue(t *testing.T) {
	cases := []struct {
		fieldA string
		fieldB string
		value  *dlit.Literal
		want   bool
	}{
		{"income", "balance", dlit.MustNew(19), true},
		{"income", "balance", dlit.MustNew(19.12), true},
		{"income", "balance", dlit.MustNew(20), true},
		{"income", "balance", dlit.MustNew(-20), false},
		{"income", "balance", dlit.MustNew(18.34), false},
		{"flow", "cost", dlit.MustNew(144.564), true},
		{"flow", "cost", dlit.MustNew(144.565), true},
		{"flow", "cost", dlit.MustNew(144.563), false},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(4),
		"balance": dlit.MustNew(15),
		"cost":    dlit.MustNew(20),
		"flow":    dlit.MustNew(124.564),
	}
	for _, c := range cases {
		r := NewAddLEF(c.fieldA, c.fieldB, c.value)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) (rule: %s) err: %v", r, err)
		} else if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestAddLEFIsTrue_errors(t *testing.T) {
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
				Rule: NewAddLEF("fred", "flow", dlit.MustNew(7.894)),
			},
		},
		{fieldA: "flow",
			fieldB: "fred",
			value:  dlit.MustNew(7.894),
			wantErr: InvalidRuleError{
				Rule: NewAddLEF("flow", "fred", dlit.MustNew(7.894)),
			},
		},
		{fieldA: "band",
			fieldB: "flow",
			value:  dlit.MustNew(7.894),
			wantErr: IncompatibleTypesRuleError{
				Rule: NewAddLEF("band", "flow", dlit.MustNew(7.894)),
			},
		},
		{fieldA: "flow",
			fieldB: "band",
			value:  dlit.MustNew(7.894),
			wantErr: IncompatibleTypesRuleError{
				Rule: NewAddLEF("flow", "band", dlit.MustNew(7.894)),
			},
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"flow":   dlit.MustNew(124.564),
		"band":   dlit.NewString("alpha"),
	}
	for _, c := range cases {
		r := NewAddLEF(c.fieldA, c.fieldB, c.value)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}

func TestAddLEFGetFields(t *testing.T) {
	r := NewAddLEF("income", "cost", dlit.MustNew(5.5))
	want := []string{"income", "cost"}
	got := r.GetFields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetFields() got: %s, want: %s", got, want)
	}
}

func TestAddLEFOverlaps(t *testing.T) {
	cases := []struct {
		ruleA *AddLEF
		ruleB Rule
		want  bool
	}{
		{ruleA: NewAddLEF("band", "cost", dlit.MustNew(7.3)),
			ruleB: NewAddLEF("band", "cost", dlit.MustNew(6.5)),
			want:  true,
		},
		{ruleA: NewAddLEF("band", "balance", dlit.MustNew(7.3)),
			ruleB: NewAddLEF("rate", "balance", dlit.MustNew(6.5)),
			want:  false,
		},
		{ruleA: NewAddLEF("band", "balance", dlit.MustNew(7.3)),
			ruleB: NewAddLEF("band", "rate", dlit.MustNew(6.5)),
			want:  false,
		},
		{ruleA: NewAddLEF("band", "cost", dlit.MustNew(7.3)),
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

func TestAddLEFTweak(t *testing.T) {
	fieldA := "income"
	fieldB := "balance"
	value := int64(800)
	rule := NewAddLEF(fieldA, fieldB, dlit.MustNew(value))
	cases := []struct {
		description *description.Description
		stage       int
		want        []Rule
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
			stage: 1,
			want: []Rule{
				NewAddLEF(fieldA, fieldB, dlit.MustNew(755)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(760)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(765)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(770)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(775)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(780)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(785)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(790)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(795)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(805)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(810)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(815)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(820)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(825)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(830)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(835)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(840)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(845)),
			},
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
			stage: 1,
			want: []Rule{
				NewAddLEF(fieldA, fieldB, dlit.MustNew(792)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(794)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(796)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(798)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(802)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(804)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(806)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(808)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(810)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(812)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(814)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(816)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(818)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(820)),
			},
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
			stage: 1,
			want: []Rule{
				NewAddLEF(fieldA, fieldB, dlit.MustNew(771)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(773)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(775)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(777)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(779)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(781)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(783)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(785)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(787)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(789)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(791)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(793)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(795)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(797)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(801)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(803)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(805)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(807)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(809)),
			},
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
			stage: 1,
			want: []Rule{
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(801)),
			},
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
			stage: 1,
			want: []Rule{
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.36)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.363)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.4)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.43)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.434)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.5)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.505)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.51)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.576)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.58)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.6)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.647)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.65)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.7)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.718)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.72)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.789)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.79)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.8)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.86)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.9)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.93)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799.931)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(800.002)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(800.07)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(800.073)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(800.1)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(800.14)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(800.144)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(800.2)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(800.215)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(800.22)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(800.286)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(800.29)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(800.3)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(800.357)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(800.36)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(800.4)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(800.428)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(800.43)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(800.499)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(800.5)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(800.57)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(800.6)),
			},
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
			stage: 2,
			want: []Rule{
				NewAddLEF(fieldA, fieldB, dlit.MustNew(778)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(781)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(784)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(787)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(790)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(793)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(796)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(799)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(802)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(805)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(808)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(811)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(814)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(817)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(820)),
				NewAddLEF(fieldA, fieldB, dlit.MustNew(823)),
			},
		},
	}
	for i, c := range cases {
		got := rule.Tweak(c.description, c.stage)
		if err := checkRulesMatch(got, c.want); err != nil {
			t.Errorf("Tweak(%d): %s, got: %s", i, err, got)
		}
	}
}

/**************************
 *  Benchmarks
 **************************/

func BenchmarkAddLEFIsTrue(b *testing.B) {
	record := map[string]*dlit.Literal{
		"band":   dlit.MustNew(23),
		"income": dlit.MustNew(1024),
		"cost":   dlit.MustNew(890),
	}
	r := NewAddLEF("cost", "income", dlit.MustNew(900.23))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = r.IsTrue(record)
	}
}
