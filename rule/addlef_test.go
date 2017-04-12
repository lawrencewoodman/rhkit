package rule

import (
	"github.com/lawrencewoodman/dlit"
	"reflect"
	"testing"
)

func TestAddLEFString(t *testing.T) {
	fieldA := "income"
	fieldB := "balance"
	value := 8.93
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
		value  float64
		want   bool
	}{
		{"income", "balance", 19, true},
		{"income", "balance", 19.12, true},
		{"income", "balance", 20, true},
		{"income", "balance", -20, false},
		{"income", "balance", 18.34, false},
		{"flow", "cost", 144.564, true},
		{"flow", "cost", 144.565, true},
		{"flow", "cost", 144.563, false},
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
		value   float64
		wantErr error
	}{
		{fieldA: "fred",
			fieldB:  "flow",
			value:   7.894,
			wantErr: InvalidRuleError{Rule: NewAddLEF("fred", "flow", 7.894)},
		},
		{fieldA: "flow",
			fieldB:  "fred",
			value:   7.894,
			wantErr: InvalidRuleError{Rule: NewAddLEF("flow", "fred", 7.894)},
		},
		{fieldA: "band",
			fieldB: "flow",
			value:  7.894,
			wantErr: IncompatibleTypesRuleError{
				Rule: NewAddLEF("band", "flow", 7.894),
			},
		},
		{fieldA: "flow",
			fieldB: "band",
			value:  7.894,
			wantErr: IncompatibleTypesRuleError{
				Rule: NewAddLEF("flow", "band", 7.894),
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
	r := NewAddLEF("income", "cost", 5.5)
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
		{ruleA: NewAddLEF("band", "cost", 7.3),
			ruleB: NewAddLEF("band", "cost", 6.5),
			want:  true,
		},
		{ruleA: NewAddLEF("band", "balance", 7.3),
			ruleB: NewAddLEF("rate", "balance", 6.5),
			want:  false,
		},
		{ruleA: NewAddLEF("band", "balance", 7.3),
			ruleB: NewAddLEF("band", "rate", 6.5),
			want:  false,
		},
		{ruleA: NewAddLEF("band", "cost", 7.3),
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
