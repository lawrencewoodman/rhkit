package rule

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
	"reflect"
	"testing"
)

func TestNewOutsideFVF(t *testing.T) {
	low := float64(5.78)
	high := float64(6.44)
	r, err := NewOutsideFVF("flow", low, high)
	if err != nil {
		t.Errorf("NewOutsideFVF(%s, %f, %f) got err: %s", "flow", low, high, err)
	}
	if r == nil {
		t.Errorf("NewOutsideFVF(%s, %f, %f) got r: nil", "flow", low, high)
	}
}

func TestNewOutsideFVF_errors(t *testing.T) {
	cases := []struct {
		low        float64
		high       float64
		wantErrStr string
	}{
		{low: 5.78,
			high:       5.78,
			wantErrStr: "can't create Outside rule where high: 5.78 <= low: 5.78",
		},
		{low: 6.23,
			high:       5.35,
			wantErrStr: "can't create Outside rule where high: 5.35 <= low: 6.23",
		},
	}
	field := "flow"
	for _, c := range cases {
		r, err := NewOutsideFVF(field, c.low, c.high)
		if r != nil {
			t.Errorf("NewOutsideFVF(%s, %f, %f) rule got: %s, want: nil",
				field, c.low, c.high, r)
		}
		if err == nil {
			t.Errorf("NewOutsideFVF(%s, %f, %f) got err: nil, want: %s",
				field, c.low, c.high, c.wantErrStr)
		} else if err.Error() != c.wantErrStr {
			t.Errorf("NewOutsideFVF(%s, %f, %f) got err: %s, want: %s",
				field, c.low, c.high, err, c.wantErrStr)
		}
	}
}

func TestOutsideFVFString(t *testing.T) {
	field := "flow"
	low := float64(183.78)
	high := float64(287.28)
	want := "flow <= 183.78 || flow >= 287.28"
	r, err := NewOutsideFVF(field, low, high)
	if err != nil {
		t.Fatalf("NewOutsideFVF: %s", err)
	}
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestOutsideFVFIsTrue(t *testing.T) {
	cases := []struct {
		field string
		low   float64
		high  float64
		want  bool
	}{
		{field: "income", low: 20.23, high: 21.45, want: true},
		{field: "income", low: 19.63, high: 21.92, want: true},
		{field: "income", low: 30.28, high: 50.28, want: true},
		{field: "income", low: 10.24, high: 12.78, want: true},
		{field: "income", low: 10.78, high: 19.63, want: true},
		{field: "income", low: 18.82, high: 21.23, want: false},
		{field: "income", low: 10.23, high: 20.48, want: false},
		{field: "cost", low: 25.89, high: 30.28, want: true},
		{field: "cost", low: 15.24, high: 25.89, want: true},
		{field: "cost", low: 25.88, high: 72.4, want: false},
		{field: "cost", low: 24., high: 25.90, want: false},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19.63),
		"cost":   dlit.MustNew(25.89),
	}
	for _, c := range cases {
		r, err := NewOutsideFVF(c.field, c.low, c.high)
		if err != nil {
			t.Fatalf("NewOutsideFVF: %s", err)
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

func TestOutsideFVFIsTrue_errors(t *testing.T) {
	field := "rate"
	low := float64(18.47)
	high := float64(20.23)
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"cost":   dlit.MustNew(18),
		"band":   dlit.NewString("alpha"),
	}
	r, err := NewOutsideFVF(field, low, high)
	if err != nil {
		t.Fatalf("NewOutsideFVF: %s", err)
	}
	wantErr := InvalidRuleError{Rule: r}
	_, err = r.IsTrue(record)
	if err != wantErr {
		t.Errorf("IsTrue(record) rule: %s, err: %v, want: %v", r, err, wantErr)
	}
}

func TestOutsideFVFGetFields(t *testing.T) {
	field := "rate"
	low := float64(18.54)
	high := float64(20.302)
	want := []string{"rate"}
	r, err := NewOutsideFVF(field, low, high)
	if err != nil {
		t.Fatalf("NewOutsideFVF: %s", err)
	}
	got := r.GetFields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetFields() got: %s, want: %s", got, want)
	}
}

func TestOutsideFVFTweak(t *testing.T) {
	field := "income"
	low := float64(800)
	high := float64(1000)
	rule := MustNewOutsideFVF(field, low, high)
	fdMin := float64(500)
	fdMax := float64(2000)
	description := &description.Description{
		map[string]*description.Field{
			"income": &description.Field{
				Kind: fieldtype.Float,
				Min:  dlit.MustNew(fdMin),
				Max:  dlit.MustNew(fdMax),
			},
		},
	}
	got := rule.Tweak(description, 1)
	numGot := len(got)
	if numGot < 150 {
		t.Errorf("Tweak - got too few rules returned: %d", numGot)
	}
	uniqueRules := Uniq(got)
	if len(uniqueRules) != numGot {
		t.Errorf("Tweak - num uniqueRules: %d != num got: %d",
			len(uniqueRules), numGot)
	}
	for _, r := range got {
		switch x := r.(type) {
		case *OutsideFVF:
			lowV := x.GetLow()
			highV := x.GetHigh()
			if lowV <= fdMin || highV >= fdMax || lowV == low || highV == high {
				t.Errorf("Tweak - invalid rule: %s", r)
			}
		default:
			t.Errorf("Tweak - invalid rule: %s", r)
		}
	}
}

func TestOutsideFVFOverlaps(t *testing.T) {
	cases := []struct {
		ruleA *OutsideFVF
		ruleB Rule
		want  bool
	}{
		{ruleA: MustNewOutsideFVF("band", 7.9, 120.9),
			ruleB: MustNewOutsideFVF("band", 6.3, 50.3),
			want:  true,
		},
		{ruleA: MustNewOutsideFVF("band", 7.9, 50.9),
			ruleB: MustNewOutsideFVF("rate", 6.3, 90.3),
			want:  false,
		},
		{ruleA: MustNewOutsideFVF("band", 7.9, 40.9),
			ruleB: NewGEFVF("band", 6.3),
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
