package rule

import (
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestDynamicNew_errors(t *testing.T) {
	expr := "profit > > 0"
	wantErr := InvalidExprError{Expr: "profit > > 0"}
	_, err := NewDynamic(expr)
	if err == nil || err.Error() != wantErr.Error() {
		t.Fatalf("NewDynamic err: %s, wantErr: %s", err, wantErr)
	}
}

func TestDynamicString(t *testing.T) {
	expr := "income <= cost"
	r, err := NewDynamic(expr)
	if err != nil {
		t.Fatalf("NewDynamic: %s", err)
	}
	got := r.String()
	if got != expr {
		t.Errorf("String() got: %s, want: %s", got, expr)
	}
}

func TestDynamicIsTrue(t *testing.T) {
	cases := []struct {
		expr string
		want bool
	}{
		{"income <= cost", true},
		{"cost <= income", false},
		{"income <= income", true},
		{"flowIn <= flowOut", true},
		{"flowOut <= flowIn", false},
		{"flowIn <= flowIn", true},
		{"income <= flowIn", true},
		{"flowIn <= income", false},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"cost":    dlit.MustNew(20),
		"flowIn":  dlit.MustNew(124.564),
		"flowOut": dlit.MustNew(124.565),
	}
	for _, c := range cases {
		r, err := NewDynamic(c.expr)
		if err != nil {
			t.Fatalf("NewDynamic: %s", err)
		}
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestDynamicIsTrue_errors(t *testing.T) {
	cases := []struct {
		expr    string
		wantErr error
	}{
		{expr: "income <= band",
			wantErr: IncompatibleTypesRuleError{Rule: NewLEFF("income", "band")},
		},
		{expr: "band <= income",
			wantErr: IncompatibleTypesRuleError{Rule: NewLEFF("band", "income")},
		},
		{expr: "flow <= band",
			wantErr: IncompatibleTypesRuleError{Rule: NewLEFF("flow", "band")},
		},
		{expr: "band <= flow",
			wantErr: IncompatibleTypesRuleError{Rule: NewLEFF("band", "flow")},
		},
		{expr: "fred <= income",
			wantErr: InvalidRuleError{Rule: NewLEFF("fred", "income")},
		},
		{expr: "income <= fred",
			wantErr: InvalidRuleError{Rule: NewLEFF("income", "fred")},
		},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"flow":   dlit.MustNew(124.564),
		"band":   dlit.NewString("alpha"),
	}

	for _, c := range cases {
		r, err := NewDynamic(c.expr)
		if err != nil {
			t.Fatalf("NewDynamic: %s", err)
		}
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}

func TestDynamicFields(t *testing.T) {
	expr := "income <= cost"
	r, err := NewDynamic(expr)
	if err != nil {
		t.Fatalf("NewDynamic: %s", err)
	}
	want := []string{}
	got := r.Fields()
	if len(got) != 0 {
		t.Errorf("Fields() got: %s, want: %s", got, want)
	}
}
