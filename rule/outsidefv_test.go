package rule

import (
	"fmt"
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/dexprfuncs"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
	"reflect"
	"testing"
)

func TestNewOutsideFV(t *testing.T) {
	low := dlit.MustNew(5.78)
	high := dlit.MustNew(6.44)
	r, err := NewOutsideFV("flow", low, high)
	if err != nil {
		t.Errorf("NewOutsideFV(%s, %s, %s) got err: %s", "flow", low, high, err)
	}
	if r == nil {
		t.Errorf("NewOutsideFV(%s, %s, %s) got r: nil", "flow", low, high)
	}
}

func TestNewOutsideFV_errors(t *testing.T) {
	cases := []struct {
		low        *dlit.Literal
		high       *dlit.Literal
		wantErrStr string
	}{
		{low: dlit.MustNew(5.78),
			high:       dlit.MustNew(5.78),
			wantErrStr: "can't create Outside rule where high: 5.78 <= low: 5.78",
		},
		{low: dlit.MustNew(6.23),
			high:       dlit.MustNew(5.35),
			wantErrStr: "can't create Outside rule where high: 5.35 <= low: 6.23",
		},
	}
	field := "flow"
	for _, c := range cases {
		r, err := NewOutsideFV(field, c.low, c.high)
		if r != nil {
			t.Errorf("NewOutsideFV(%s, %s, %s) rule got: %s, want: nil",
				field, c.low, c.high, r)
		}
		if err == nil {
			t.Errorf("NewOutsideFV(%s, %s, %s) got err: nil, want: %s",
				field, c.low, c.high, c.wantErrStr)
		} else if err.Error() != c.wantErrStr {
			t.Errorf("NewOutsideFV(%s, %s, %s) got err: %s, want: %s",
				field, c.low, c.high, err, c.wantErrStr)
		}
	}
}

func TestOutsideFVString(t *testing.T) {
	field := "flow"
	low := dlit.MustNew(183.78)
	high := dlit.MustNew(287.28)
	want := "flow <= 183.78 || flow >= 287.28"
	r, err := NewOutsideFV(field, low, high)
	if err != nil {
		t.Fatalf("NewOutsideFV: %s", err)
	}
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestOutsideFVIsTrue(t *testing.T) {
	cases := []struct {
		field string
		low   *dlit.Literal
		high  *dlit.Literal
		want  bool
	}{
		{field: "income", low: dlit.MustNew(20.23),
			high: dlit.MustNew(21.45), want: true},
		{field: "income", low: dlit.MustNew(19.63),
			high: dlit.MustNew(21.92), want: true},
		{field: "income", low: dlit.MustNew(30.28),
			high: dlit.MustNew(50.28), want: true},
		{field: "income", low: dlit.MustNew(10.24),
			high: dlit.MustNew(12.78), want: true},
		{field: "income", low: dlit.MustNew(10.78),
			high: dlit.MustNew(19.63), want: true},
		{field: "income", low: dlit.MustNew(18.82),
			high: dlit.MustNew(21.23), want: false},
		{field: "income", low: dlit.MustNew(10.23),
			high: dlit.MustNew(20.48), want: false},
		{field: "cost", low: dlit.MustNew(25.89),
			high: dlit.MustNew(30.28), want: true},
		{field: "cost", low: dlit.MustNew(15.24),
			high: dlit.MustNew(25.89), want: true},
		{field: "cost", low: dlit.MustNew(25.88),
			high: dlit.MustNew(72.4), want: false},
		{field: "cost", low: dlit.MustNew(24.),
			high: dlit.MustNew(25.90), want: false},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19.63),
		"cost":   dlit.MustNew(25.89),
	}
	for _, c := range cases {
		r, err := NewOutsideFV(c.field, c.low, c.high)
		if err != nil {
			t.Fatalf("NewOutsideFV: %s", err)
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

func TestOutsideFVIsTrue_errors(t *testing.T) {
	low := dlit.MustNew(18.47)
	high := dlit.MustNew(20.23)
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19),
		"cost":   dlit.MustNew(18),
		"band":   dlit.NewString("alpha"),
	}
	cases := []struct {
		field   string
		wantErr error
	}{{field: "rate",
		wantErr: InvalidRuleError{Rule: MustNewOutsideFV("rate", low, high)}},
		{field: "band",
			wantErr: IncompatibleTypesRuleError{
				Rule: MustNewOutsideFV("band", low, high),
			}},
	}
	for _, c := range cases {
		r, err := NewOutsideFV(c.field, low, high)
		if err != nil {
			t.Fatalf("NewOutsideFV: %s", err)
		}
		_, err = r.IsTrue(record)
		if err == nil || err.Error() != c.wantErr.Error() {
			t.Errorf("IsTrue(record) rule: %s, err: %v, want: %v", r, err, c.wantErr)
		}
	}
}

func TestOutsideFVFields(t *testing.T) {
	field := "rate"
	low := dlit.MustNew(18.54)
	high := dlit.MustNew(20.302)
	want := []string{"rate"}
	r, err := NewOutsideFV(field, low, high)
	if err != nil {
		t.Fatalf("NewOutsideFV: %s", err)
	}
	got := r.Fields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Fields() got: %s, want: %s", got, want)
	}
}

func TestOutsideFVTweak(t *testing.T) {
	field := "income"
	cases := []struct {
		fdMin   *dlit.Literal
		fdMax   *dlit.Literal
		fdMaxDP int
		low     *dlit.Literal
		high    *dlit.Literal
		wantNum int
	}{
		{fdMin: dlit.MustNew(500),
			fdMax:   dlit.MustNew(2000),
			fdMaxDP: 2,
			low:     dlit.MustNew(800),
			high:    dlit.MustNew(1000),
			wantNum: 150,
		},
		{fdMin: dlit.MustNew(18),
			fdMax:   dlit.MustNew(95),
			fdMaxDP: 0,
			low:     dlit.MustNew(22),
			high:    dlit.MustNew(26),
			wantNum: 80,
		},
	}
	inRangeExpr := dexpr.MustNew(
		"lowV > fdMin && highV < fdMax && lowV != low && highV != high",
		dexprfuncs.CallFuncs,
	)
	for _, c := range cases {
		description := &description.Description{
			map[string]*description.Field{
				"income": &description.Field{
					Kind:  fieldtype.Number,
					Min:   c.fdMin,
					Max:   c.fdMax,
					MaxDP: c.fdMaxDP,
				},
			},
		}
		rule := MustNewOutsideFV(field, c.low, c.high)
		got := rule.Tweak(description, 1)
		numGot := len(got)
		if numGot < c.wantNum {
			t.Errorf("Tweak - got too few rules returned: %d, got: %v", numGot, got)
		}
		uniqueRules := Uniq(got)
		if len(uniqueRules) != numGot {
			t.Errorf("Tweak - num uniqueRules: %d != num got: %d",
				len(uniqueRules), numGot)
		}
		vars := map[string]*dlit.Literal{
			"fdMin": c.fdMin,
			"fdMax": c.fdMax,
			"low":   c.low,
			"high":  c.high,
		}
		for _, r := range got {
			switch x := r.(type) {
			case *OutsideFV:
				vars["lowV"] = x.Low()
				vars["highV"] = x.High()
				if ok, err := inRangeExpr.EvalBool(vars); !ok || err != nil {
					t.Errorf("Tweak - invalid rule: %s", r)
				}
			default:
				t.Errorf("Tweak - invalid rule: %s", r)
			}
		}
	}
}

func TestOutsideFVOverlaps(t *testing.T) {
	cases := []struct {
		ruleA *OutsideFV
		ruleB Rule
		want  bool
	}{
		{ruleA: MustNewOutsideFV("band", dlit.MustNew(7.9), dlit.MustNew(120.9)),
			ruleB: MustNewOutsideFV("band", dlit.MustNew(6.3), dlit.MustNew(50.3)),
			want:  true,
		},
		{ruleA: MustNewOutsideFV("band", dlit.MustNew(7.9), dlit.MustNew(50.9)),
			ruleB: MustNewOutsideFV("rate", dlit.MustNew(6.3), dlit.MustNew(90.3)),
			want:  false,
		},
		{ruleA: MustNewOutsideFV("band", dlit.MustNew(7.9), dlit.MustNew(40.9)),
			ruleB: NewGEFV("band", dlit.MustNew(6.3)),
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

func TestGenerateOutsideFV(t *testing.T) {
	cases := []struct {
		description *description.Description
		field       string
		minNumRules int
		maxNumRules int
		min         *dlit.Literal
		max         *dlit.Literal
		mid         *dlit.Literal
		maxDP       int
	}{
		{description: &description.Description{
			map[string]*description.Field{
				"income": &description.Field{
					Kind:  fieldtype.Number,
					Min:   dlit.MustNew(500),
					Max:   dlit.MustNew(1000),
					MaxDP: 2,
				},
			},
		},
			field:       "income",
			minNumRules: 18 * 19 / 2,
			maxNumRules: 20 * 21 / 2,
			min:         dlit.MustNew(500),
			max:         dlit.MustNew(1000),
			mid:         dlit.MustNew(750),
			maxDP:       0,
		},
		{description: &description.Description{
			map[string]*description.Field{
				"income": &description.Field{
					Kind:  fieldtype.Number,
					Min:   dlit.MustNew(790.73),
					Max:   dlit.MustNew(1000),
					MaxDP: 2,
				},
			},
		},
			field:       "income",
			minNumRules: 18 * 19 / 2,
			maxNumRules: 20 * 21 / 2,
			min:         dlit.MustNew(790),
			max:         dlit.MustNew(1000),
			mid:         dlit.MustNew(903),
			maxDP:       2,
		},
		{description: &description.Description{
			map[string]*description.Field{
				"income": &description.Field{
					Kind:  fieldtype.Number,
					Min:   dlit.MustNew(799),
					Max:   dlit.MustNew(801),
					MaxDP: 0,
				},
			},
		},
			field:       "income",
			minNumRules: 0,
			maxNumRules: 0,
			min:         dlit.MustNew(799),
			max:         dlit.MustNew(801),
			mid:         dlit.MustNew(800),
			maxDP:       0,
		},
		{description: &description.Description{
			map[string]*description.Field{
				"income": &description.Field{
					Kind:  fieldtype.Number,
					Min:   dlit.MustNew(700),
					Max:   dlit.MustNew(800),
					MaxDP: 0,
				},
				"month": &description.Field{
					Kind: fieldtype.String,
				},
			},
		},
			field:       "month",
			minNumRules: 0,
			maxNumRules: 0,
			min:         dlit.MustNew(0),
			max:         dlit.MustNew(0),
			mid:         dlit.MustNew(0),
			maxDP:       0,
		},
	}
	complexity := Complexity{}
	ruleFields := []string{"income"}
	for i, c := range cases {
		complyFunc := func(r Rule) error {
			x, ok := r.(*OutsideFV)
			if !ok {
				return fmt.Errorf("wrong type: %T (%s)", r, r)
			}
			if x.field != c.field {
				return fmt.Errorf("field isn't correct for rule: %s", r)
			}
			return nil
		}
		got := generateOutsideFV(c.description, ruleFields, complexity, c.field)
		err := checkRulesComply(
			got,
			c.minNumRules,
			c.maxNumRules,
			c.min,
			c.max,
			c.mid,
			c.maxDP,
			complyFunc,
		)
		if err != nil {
			t.Errorf("(%d) generateOutsideFV: %s", i, err)
		}
	}
}
