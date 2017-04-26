package rule

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
	"reflect"
	"testing"
)

func TestNewOutsideFVI(t *testing.T) {
	low := int64(5)
	high := int64(6)
	r, err := NewOutsideFVI("flow", low, high)
	if err != nil {
		t.Errorf("NewOutsideFVI(%s, %d, %d) got err: %s", "flow", low, high, err)
	}
	if r == nil {
		t.Errorf("NewOutsideFVI(%s, %d, %d) got r: nil", "flow", low, high)
	}
}

func TestNewOutsideFVI_errors(t *testing.T) {
	cases := []struct {
		low        int64
		high       int64
		wantErrStr string
	}{
		{low: 5,
			high:       5,
			wantErrStr: "can't create Outside rule where high: 5 <= low: 5",
		},
		{low: 6,
			high:       5,
			wantErrStr: "can't create Outside rule where high: 5 <= low: 6",
		},
	}
	field := "flow"
	for _, c := range cases {
		r, err := NewOutsideFVI(field, c.low, c.high)
		if r != nil {
			t.Errorf("NewOutsideFVI(%s, %d, %d) rule got: %s, want: nil",
				field, c.low, c.high, r)
		}
		if err == nil {
			t.Errorf("NewOutsideFVI(%s, %d, %d) got err: nil, want: %s",
				field, c.low, c.high, c.wantErrStr)
		} else if err.Error() != c.wantErrStr {
			t.Errorf("NewOutsideFVI(%s, %d, %d) got err: %s, want: %s",
				field, c.low, c.high, err, c.wantErrStr)
		}
	}
}

func TestOutsideFVIString(t *testing.T) {
	field := "flow"
	low := int64(183)
	high := int64(287)
	want := "flow <= 183 || flow >= 287"
	r, err := NewOutsideFVI(field, low, high)
	if err != nil {
		t.Fatalf("NewOutsideFVI: %s", err)
	}
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestOutsideFVIIsTrue(t *testing.T) {
	cases := []struct {
		field string
		low   int64
		high  int64
		want  bool
	}{
		{field: "income", low: 20, high: 21, want: true},
		{field: "income", low: 19, high: 21, want: true},
		{field: "income", low: 30, high: 50, want: true},
		{field: "income", low: 10, high: 12, want: true},
		{field: "income", low: 10, high: 19, want: true},
		{field: "income", low: 18, high: 21, want: false},
		{field: "income", low: 10, high: 20, want: false},
		{field: "cost", low: 25, high: 30, want: true},
		{field: "cost", low: 15, high: 25, want: true},
		{field: "cost", low: 24, high: 30, want: false},
		{field: "cost", low: 24, high: 26, want: false},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"cost":   dlit.MustNew(25),
	}
	for _, c := range cases {
		r, err := NewOutsideFVI(c.field, c.low, c.high)
		if err != nil {
			t.Fatalf("NewOutsideFVI: %s", err)
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

func TestOutsideFVIIsTrue_errors(t *testing.T) {
	field := "rate"
	low := int64(18)
	high := int64(20)
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"cost":   dlit.MustNew(18),
		"band":   dlit.NewString("alpha"),
	}
	r, err := NewOutsideFVI(field, low, high)
	if err != nil {
		t.Fatalf("NewOutsideFVI: %s", err)
	}
	wantErr := InvalidRuleError{Rule: r}
	_, err = r.IsTrue(record)
	if err != wantErr {
		t.Errorf("IsTrue(record) rule: %s, err: %v, want: %v", r, err, wantErr)
	}
}

func TestOutsideFVIGetFields(t *testing.T) {
	field := "rate"
	low := int64(18)
	high := int64(20)
	want := []string{"rate"}
	r, err := NewOutsideFVI(field, low, high)
	if err != nil {
		t.Fatalf("NewOutsideFVI: %s", err)
	}
	got := r.GetFields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetFields() got: %s, want: %s", got, want)
	}
}

func TestOutsideFVITweak(t *testing.T) {
	field := "income"
	low := int64(800)
	high := int64(1000)
	rule := MustNewOutsideFVI(field, low, high)
	fdMin := int64(500)
	fdMax := int64(2000)
	description := &description.Description{
		map[string]*description.Field{
			"income": &description.Field{
				Kind: fieldtype.Int,
				Min:  dlit.MustNew(fdMin),
				Max:  dlit.MustNew(fdMax),
			},
		},
	}
	got := rule.Tweak(description, 1)
	numGot := len(got)
	if numGot < 300 {
		t.Errorf("Tweak - got too few rules returned: %d", numGot)
	}
	uniqueRules := Uniq(got)
	if len(uniqueRules) != numGot {
		t.Errorf("Tweak - num uniqueRules: %d != num got: %d",
			len(uniqueRules), numGot)
	}
	for _, r := range got {
		switch x := r.(type) {
		case *OutsideFVI:
			lowV := x.GetLow()
			highV := x.GetHigh()
			if lowV <= fdMin || highV >= fdMax || (lowV == low && highV == high) {
				t.Errorf("Tweak - invalid rule: %s", r)
			}
		default:
			t.Errorf("Tweak - invalid rule: %s", r)
		}
	}
}

func TestOutsideFVIOverlaps(t *testing.T) {
	cases := []struct {
		ruleA *OutsideFVI
		ruleB Rule
		want  bool
	}{
		{ruleA: MustNewOutsideFVI("band", 7, 120),
			ruleB: MustNewOutsideFVI("band", 6, 50),
			want:  true,
		},
		{ruleA: MustNewOutsideFVI("band", 7, 50),
			ruleB: MustNewOutsideFVI("rate", 6, 90),
			want:  false,
		},
		{ruleA: MustNewOutsideFVI("band", 7, 40),
			ruleB: NewGEFVI("band", 6),
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
