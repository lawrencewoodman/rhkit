package rule

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
	"reflect"
	"testing"
)

func TestMulGEFString(t *testing.T) {
	fieldA := "income"
	fieldB := "balance"
	value := dlit.MustNew(8.93)
	want := "income * balance >= 8.93"
	r := NewMulGEF(fieldA, fieldB, value)
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestMulGEFIsTrue(t *testing.T) {
	cases := []struct {
		fieldA string
		fieldB string
		value  *dlit.Literal
		want   bool
	}{
		{"income", "balance", dlit.MustNew(61), false},
		{"income", "balance", dlit.MustNew(60.12), false},
		{"income", "balance", dlit.MustNew(60), true},
		{"income", "balance", dlit.MustNew(-60), true},
		{"income", "balance", dlit.MustNew(59.89), true},
		{"flow", "cost", dlit.MustNew(2491.2), true},
		{"flow", "cost", dlit.MustNew(2491.21), false},
		{"flow", "cost", dlit.MustNew(2491.19), true},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(4),
		"balance": dlit.MustNew(15),
		"cost":    dlit.MustNew(20),
		"flow":    dlit.MustNew(124.56),
	}
	for _, c := range cases {
		r := NewMulGEF(c.fieldA, c.fieldB, c.value)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) (rule: %s) err: %v", r, err)
		} else if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestMulGEFIsTrue_errors(t *testing.T) {
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
				Rule: NewMulGEF("fred", "flow", dlit.MustNew(7.894)),
			},
		},
		{fieldA: "flow",
			fieldB: "fred",
			value:  dlit.MustNew(7.894),
			wantErr: InvalidRuleError{
				Rule: NewMulGEF("flow", "fred", dlit.MustNew(7.894)),
			},
		},
		{fieldA: "band",
			fieldB: "flow",
			value:  dlit.MustNew(7.894),
			wantErr: IncompatibleTypesRuleError{
				Rule: NewMulGEF("band", "flow", dlit.MustNew(7.894)),
			},
		},
		{fieldA: "flow",
			fieldB: "band",
			value:  dlit.MustNew(7.894),
			wantErr: IncompatibleTypesRuleError{
				Rule: NewMulGEF("flow", "band", dlit.MustNew(7.894)),
			},
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"flow":   dlit.MustNew(124.564),
		"band":   dlit.NewString("alpha"),
	}
	for _, c := range cases {
		r := NewMulGEF(c.fieldA, c.fieldB, c.value)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}

func TestMulGEFGetFields(t *testing.T) {
	r := NewMulGEF("income", "cost", dlit.MustNew(5.5))
	want := []string{"income", "cost"}
	got := r.GetFields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetFields() got: %s, want: %s", got, want)
	}
}

func TestMulGEFOverlaps(t *testing.T) {
	cases := []struct {
		ruleA *MulGEF
		ruleB Rule
		want  bool
	}{
		{ruleA: NewMulGEF("band", "cost", dlit.MustNew(7.3)),
			ruleB: NewMulGEF("band", "cost", dlit.MustNew(6.5)),
			want:  true,
		},
		{ruleA: NewMulGEF("band", "balance", dlit.MustNew(7.3)),
			ruleB: NewMulGEF("rate", "balance", dlit.MustNew(6.5)),
			want:  false,
		},
		{ruleA: NewMulGEF("band", "balance", dlit.MustNew(7.3)),
			ruleB: NewMulGEF("band", "rate", dlit.MustNew(6.5)),
			want:  false,
		},
		{ruleA: NewMulGEF("band", "cost", dlit.MustNew(7.3)),
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

func TestMulGEFTweak(t *testing.T) {
	fieldA := "income"
	fieldB := "balance"
	value := int64(156250)
	rule := NewMulGEF(fieldA, fieldB, dlit.MustNew(value))
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
				NewMulGEF(fieldA, fieldB, dlit.MustNew(139375)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(141250)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(143125)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(145000)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(146875)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(148750)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(150625)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(152500)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(154375)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(158125)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(160000)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(161875)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(163750)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(165625)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(167500)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(169375)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(171250)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(173125)),
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
				NewMulGEF(fieldA, fieldB, dlit.MustNew(149500)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(150250)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(151000)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(151750)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(152500)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(153250)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(154000)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(154750)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(155500)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(157000)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(157750)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(158500)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(159250)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(160000)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(160750)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(161500)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(162250)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(163000)),
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
				NewMulGEF(fieldA, fieldB, dlit.MustNew(147253)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(147555)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(147858)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(148160)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(148463)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(148765)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(149068)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(149370)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(149673)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(149975)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(150278)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(150580)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(150883)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(151185)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(151488)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(151790)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(152093)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(152395)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(152698)),
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
					Min:   dlit.MustNew(300.924),
					Max:   dlit.MustNew(505),
					MaxDP: 3,
				},
			},
		},
			stage: 1,
			want: []Rule{
				NewMulGEF(fieldA, fieldB, dlit.MustNew(147337.556)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(147556.632)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(147775.708)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(147994.784)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(148213.86)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(148432.936)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(148652.012)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(148871.088)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(149090.164)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(149309.24)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(149528.316)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(149747.392)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(149966.468)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(150185.544)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(150404.62)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(150623.696)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(150842.772)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(151061.848)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(151280.924)),
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
				NewMulGEF(fieldA, fieldB, dlit.MustNew(149500)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(150250)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(151000)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(151750)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(152500)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(153250)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(154000)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(154750)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(155500)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(157000)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(157750)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(158500)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(159250)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(160000)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(160750)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(161500)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(162250)),
				NewMulGEF(fieldA, fieldB, dlit.MustNew(163000)),
			},
		},
	}
	for i, c := range cases {
		got := rule.Tweak(c.description, c.stage)
		if err := checkRulesMatch(got, c.want); err != nil {
			t.Errorf("(%d) Tweak: %s, got: %s", i, err, got)
		}
	}
}

/**************************
 *  Benchmarks
 **************************/

func BenchmarkMulGEFIsTrue(b *testing.B) {
	record := map[string]*dlit.Literal{
		"band":   dlit.MustNew(23),
		"income": dlit.MustNew(1024),
		"cost":   dlit.MustNew(890),
	}
	r := NewMulGEF("cost", "income", dlit.MustNew(900.23))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = r.IsTrue(record)
	}
}
