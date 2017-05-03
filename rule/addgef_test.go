package rule

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
	"reflect"
	"testing"
)

func TestAddGEFString(t *testing.T) {
	fieldA := "income"
	fieldB := "balance"
	value := dlit.MustNew(8.93)
	want := "income + balance >= 8.93"
	r := NewAddGEF(fieldA, fieldB, value)
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestAddGEFIsTrue(t *testing.T) {
	cases := []struct {
		fieldA string
		fieldB string
		value  *dlit.Literal
		want   bool
	}{
		{"income", "balance", dlit.MustNew(19), true},
		{"income", "balance", dlit.MustNew(19.12), false},
		{"income", "balance", dlit.MustNew(20), false},
		{"income", "balance", dlit.MustNew(-20), true},
		{"income", "balance", dlit.MustNew(18.34), true},
		{"flow", "cost", dlit.MustNew(144.564), true},
		{"flow", "cost", dlit.MustNew(144.565), false},
		{"flow", "cost", dlit.MustNew(144.563), true},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(4),
		"balance": dlit.MustNew(15),
		"cost":    dlit.MustNew(20),
		"flow":    dlit.MustNew(124.564),
	}
	for _, c := range cases {
		r := NewAddGEF(c.fieldA, c.fieldB, c.value)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) (rule: %s) err: %v", r, err)
		} else if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestAddGEFIsTrue_errors(t *testing.T) {
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
				Rule: NewAddGEF("fred", "flow", dlit.MustNew(7.894)),
			},
		},
		{fieldA: "flow",
			fieldB: "fred",
			value:  dlit.MustNew(7.894),
			wantErr: InvalidRuleError{
				Rule: NewAddGEF("flow", "fred", dlit.MustNew(7.894)),
			},
		},
		{fieldA: "band",
			fieldB: "flow",
			value:  dlit.MustNew(7.894),
			wantErr: IncompatibleTypesRuleError{
				Rule: NewAddGEF("band", "flow", dlit.MustNew(7.894)),
			},
		},
		{fieldA: "flow",
			fieldB: "band",
			value:  dlit.MustNew(7.894),
			wantErr: IncompatibleTypesRuleError{
				Rule: NewAddGEF("flow", "band", dlit.MustNew(7.894)),
			},
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"flow":   dlit.MustNew(124.564),
		"band":   dlit.NewString("alpha"),
	}
	for _, c := range cases {
		r := NewAddGEF(c.fieldA, c.fieldB, c.value)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}

func TestAddGEFGetFields(t *testing.T) {
	r := NewAddGEF("income", "cost", dlit.MustNew(5.5))
	want := []string{"income", "cost"}
	got := r.GetFields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetFields() got: %s, want: %s", got, want)
	}
}

func TestAddGEFOverlaps(t *testing.T) {
	cases := []struct {
		ruleA *AddGEF
		ruleB Rule
		want  bool
	}{
		{ruleA: NewAddGEF("band", "cost", dlit.MustNew(7.3)),
			ruleB: NewAddGEF("band", "cost", dlit.MustNew(6.5)),
			want:  true,
		},
		{ruleA: NewAddGEF("band", "balance", dlit.MustNew(7.3)),
			ruleB: NewAddGEF("rate", "balance", dlit.MustNew(6.5)),
			want:  false,
		},
		{ruleA: NewAddGEF("band", "balance", dlit.MustNew(7.3)),
			ruleB: NewAddGEF("band", "rate", dlit.MustNew(6.5)),
			want:  false,
		},
		{ruleA: NewAddGEF("band", "cost", dlit.MustNew(7.3)),
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

func TestAddGEFTweak(t *testing.T) {
	fieldA := "income"
	fieldB := "balance"
	value := int64(800)
	rule := NewAddGEF(fieldA, fieldB, dlit.MustNew(value))
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
				NewAddGEF(fieldA, fieldB, dlit.MustNew(755)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(760)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(765)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(770)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(775)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(780)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(785)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(790)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(795)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(805)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(810)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(815)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(820)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(825)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(830)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(835)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(840)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(845)),
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
				NewAddGEF(fieldA, fieldB, dlit.MustNew(792)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(794)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(796)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(798)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(802)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(804)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(806)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(808)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(810)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(812)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(814)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(816)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(818)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(820)),
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
				NewAddGEF(fieldA, fieldB, dlit.MustNew(771)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(773)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(775)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(777)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(779)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(781)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(783)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(785)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(787)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(789)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(791)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(793)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(795)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(797)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(801)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(803)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(805)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(807)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(809)),
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
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(801)),
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
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.36)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.363)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.4)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.43)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.434)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.5)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.505)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.51)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.576)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.58)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.6)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.647)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.65)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.7)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.718)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.72)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.789)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.79)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.8)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.86)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.9)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.93)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799.931)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(800.002)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(800.07)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(800.073)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(800.1)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(800.14)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(800.144)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(800.2)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(800.215)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(800.22)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(800.286)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(800.29)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(800.3)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(800.357)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(800.36)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(800.4)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(800.428)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(800.43)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(800.499)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(800.5)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(800.57)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(800.6)),
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
				NewAddGEF(fieldA, fieldB, dlit.MustNew(778)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(781)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(784)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(787)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(790)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(793)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(796)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(799)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(802)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(805)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(808)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(811)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(814)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(817)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(820)),
				NewAddGEF(fieldA, fieldB, dlit.MustNew(823)),
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

func BenchmarkAddGEFIsTrue(b *testing.B) {
	record := map[string]*dlit.Literal{
		"band":   dlit.MustNew(23),
		"income": dlit.MustNew(1024),
		"cost":   dlit.MustNew(890),
	}
	r := NewAddGEF("cost", "income", dlit.MustNew(900.23))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = r.IsTrue(record)
	}
}
