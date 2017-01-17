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

func TestLEFVIGetTweakableParts(t *testing.T) {
	field := "income"
	value := int64(893)
	r := NewLEFVI(field, value)
	a, b, c := r.GetTweakableParts()
	if a != field || b != "<=" || c != "893" {
		t.Errorf("GetTweakableParts() got: \"%s\",\"%s\",\"%s\" - want: \"%s\",\"<=\",\"893\"",
			a, b, c, field)
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

func TestLEFVICloneWithValue(t *testing.T) {
	field := "income"
	value := int64(893)
	r := NewLEFVI(field, value)
	want := "income <= -27489"
	cr := r.CloneWithValue(int64(-27489))
	got := cr.String()
	if got != want {
		t.Errorf("CloseWithValue() got: %s, want: %s", got, want)
	}
}

func TestLEFVICloneWithValue_panics(t *testing.T) {
	wantPanic := "can't clone with newValue: fred of type string, need type int64"
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("CloneWithValue() didn't panic")
		} else if r.(string) != wantPanic {
			t.Errorf("CloseWithValue() - got panic: %s, wanted: %s",
				r, wantPanic)
		}
	}()
	field := "income"
	value := int64(893)
	r := NewLEFVI(field, value)
	r.CloneWithValue("fred")
}

func TestLEFVIGetFields(t *testing.T) {
	r := NewLEFVI("income", 5)
	want := []string{"income"}
	got := r.GetFields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetFields() got: %s, want: %s", got, want)
	}
}
