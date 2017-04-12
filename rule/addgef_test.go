package rule

import (
	"github.com/lawrencewoodman/dlit"
	"reflect"
	"testing"
)

func TestAddGEFString(t *testing.T) {
	fieldA := "income"
	fieldB := "balance"
	value := 8.93
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
		value  float64
		want   bool
	}{
		{"income", "balance", 19, true},
		{"income", "balance", 19.12, false},
		{"income", "balance", 20, false},
		{"income", "balance", -20, true},
		{"income", "balance", 18.34, true},
		{"flow", "cost", 144.564, true},
		{"flow", "cost", 144.565, false},
		{"flow", "cost", 144.563, true},
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
		value   float64
		wantErr error
	}{
		{fieldA: "fred",
			fieldB:  "flow",
			value:   7.894,
			wantErr: InvalidRuleError{Rule: NewAddGEF("fred", "flow", 7.894)},
		},
		{fieldA: "flow",
			fieldB:  "fred",
			value:   7.894,
			wantErr: InvalidRuleError{Rule: NewAddGEF("flow", "fred", 7.894)},
		},
		{fieldA: "band",
			fieldB: "flow",
			value:  7.894,
			wantErr: IncompatibleTypesRuleError{
				Rule: NewAddGEF("band", "flow", 7.894),
			},
		},
		{fieldA: "flow",
			fieldB: "band",
			value:  7.894,
			wantErr: IncompatibleTypesRuleError{
				Rule: NewAddGEF("flow", "band", 7.894),
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
	r := NewAddGEF("income", "cost", 5.5)
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
		{ruleA: NewAddGEF("band", "cost", 7.3),
			ruleB: NewAddGEF("band", "cost", 6.5),
			want:  true,
		},
		{ruleA: NewAddGEF("band", "balance", 7.3),
			ruleB: NewAddGEF("rate", "balance", 6.5),
			want:  false,
		},
		{ruleA: NewAddGEF("band", "balance", 7.3),
			ruleB: NewAddGEF("band", "rate", 6.5),
			want:  false,
		},
		{ruleA: NewAddGEF("band", "cost", 7.3),
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
