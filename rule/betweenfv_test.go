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

func TestNewBetweenFV(t *testing.T) {
	min := dlit.MustNew(5.89)
	max := dlit.MustNew(6.72)
	r, err := NewBetweenFV("flow", min, max)
	if err != nil {
		t.Errorf("NewBetweenFV(%s, %s, %s) got err: %s", "flow", min, max, err)
	}
	if r == nil {
		t.Errorf("NewBetweenFV(%s, %s, %s) got r: nil", "flow", min, max)
	}
}

func TestNewBetweenFV_errors(t *testing.T) {
	cases := []struct {
		min        *dlit.Literal
		max        *dlit.Literal
		wantErrStr string
	}{
		{min: dlit.MustNew(5),
			max:        dlit.MustNew(5),
			wantErrStr: "can't create Between rule where max: 5 <= min: 5",
		},
		{min: dlit.MustNew(5.65),
			max:        dlit.MustNew(5.65),
			wantErrStr: "can't create Between rule where max: 5.65 <= min: 5.65",
		},
		{min: dlit.MustNew(6.72),
			max:        dlit.MustNew(5.89),
			wantErrStr: "can't create Between rule where max: 5.89 <= min: 6.72",
		},
	}
	field := "flow"
	for _, c := range cases {
		r, err := NewBetweenFV(field, c.min, c.max)
		if r != nil {
			t.Errorf("NewBetweenFV(%s, %s, %s) rule got: %s, want: nil",
				field, c.min, c.max, r)
		}
		if err == nil {
			t.Errorf("NewBetweenFV(%s, %s, %s) got err: nil, want: %s",
				field, c.min, c.max, c.wantErrStr)
		} else if err.Error() != c.wantErrStr {
			t.Errorf("NewBetweenFV(%s, %s, %s) got err: %s, want: %s",
				field, c.min, c.max, err, c.wantErrStr)
		}
	}
}

func TestBetweenFVString(t *testing.T) {
	field := "flow"
	min := dlit.MustNew(183.92837)
	max := dlit.MustNew(287.87442)
	want := "flow >= 183.92837 && flow <= 287.87442"
	r, err := NewBetweenFV(field, min, max)
	if err != nil {
		t.Fatalf("NewBetweenFV: %s", err)
	}
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestBetweenFVIsTrue(t *testing.T) {
	cases := []struct {
		field string
		min   *dlit.Literal
		max   *dlit.Literal
		want  bool
	}{
		{field: "income", min: dlit.MustNew(18.3),
			max: dlit.MustNew(20.45), want: true},
		{field: "income", min: dlit.MustNew(19.19),
			max: dlit.MustNew(20.78), want: true},
		{field: "income", min: dlit.MustNew(19.81),
			max: dlit.MustNew(19.83), want: true},
		{field: "income", min: dlit.MustNew(18.78),
			max: dlit.MustNew(19.92), want: true},
		{field: "income", min: dlit.MustNew(10.12),
			max: dlit.MustNew(25.986), want: true},
		{field: "income", min: dlit.MustNew(10.34),
			max: dlit.MustNew(19.81), want: false},
		{field: "income", min: dlit.MustNew(19.83),
			max: dlit.MustNew(30.5), want: false},
		{field: "cost", min: dlit.MustNew(20.67),
			max: dlit.MustNew(30.89), want: true},
		{field: "cost", min: dlit.MustNew(20.23),
			max: dlit.MustNew(25.98), want: true},
		{field: "cost", min: dlit.MustNew(25.98),
			max: dlit.MustNew(30.7), want: true},
		{field: "cost", min: dlit.MustNew(20.2),
			max: dlit.MustNew(25.97), want: false},
		{field: "cost", min: dlit.MustNew(25.99),
			max: dlit.MustNew(30.7), want: false},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19.82),
		"cost":   dlit.MustNew(25.98),
	}
	for _, c := range cases {
		r, err := NewBetweenFV(c.field, c.min, c.max)
		if err != nil {
			t.Fatalf("NewBetweenFV: %s", err)
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

func TestBetweenFVIsTrue_errors(t *testing.T) {
	min := dlit.MustNew(18.72)
	max := dlit.MustNew(20.64)
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19.27),
		"cost":   dlit.MustNew(18.34),
		"band":   dlit.NewString("alpha"),
	}
	cases := []struct {
		field   string
		wantErr error
	}{{field: "rate",
		wantErr: InvalidRuleError{Rule: MustNewBetweenFV("rate", min, max)}},
		{field: "band",
			wantErr: IncompatibleTypesRuleError{
				Rule: MustNewBetweenFV("band", min, max),
			}},
	}
	for _, c := range cases {
		r, err := NewBetweenFV(c.field, min, max)
		if err != nil {
			t.Fatalf("NewBetweenFV: %s", err)
		}
		_, err = r.IsTrue(record)
		if err == nil || err.Error() != c.wantErr.Error() {
			t.Errorf("IsTrue(record) rule: %s, err: %v, want: %v", r, err, c.wantErr)
		}
	}
}

func TestBetweenFVFields(t *testing.T) {
	field := "rate"
	min := dlit.MustNew(18.72)
	max := dlit.MustNew(20.72)
	want := []string{"rate"}
	r, err := NewBetweenFV(field, min, max)
	if err != nil {
		t.Fatalf("NewBetweenFV: %s", err)
	}
	got := r.Fields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Fields() got: %s, want: %s", got, want)
	}
}

func TestBetweenFVTweak(t *testing.T) {
	field := "income"
	cases := []struct {
		fdMin   *dlit.Literal
		fdMax   *dlit.Literal
		min     *dlit.Literal
		max     *dlit.Literal
		wantNum int
	}{
		{fdMin: dlit.MustNew(500),
			fdMax:   dlit.MustNew(2000),
			min:     dlit.MustNew(800),
			max:     dlit.MustNew(1000),
			wantNum: 150,
		},
		{fdMin: dlit.MustNew(18),
			fdMax:   dlit.MustNew(95),
			min:     dlit.MustNew(22),
			max:     dlit.MustNew(26),
			wantNum: 150,
		},
	}
	inRangeExpr := dexpr.MustNew(
		"minV >= fdMin && maxV <= fdMax && minV != min && maxV != max",
		dexprfuncs.CallFuncs,
	)
	for _, c := range cases {
		description := &description.Description{
			map[string]*description.Field{
				"income": {
					Kind:  fieldtype.Number,
					Min:   c.fdMin,
					Max:   c.fdMax,
					MaxDP: 2,
				},
			},
		}
		rule := MustNewBetweenFV(field, c.min, c.max)
		got := rule.Tweak(description, 1)
		numGot := len(got)
		if numGot < c.wantNum {
			t.Errorf("Tweak - got too few rules returned: %d", numGot)
		}
		uniqueRules := Uniq(got)
		if len(uniqueRules) != numGot {
			t.Errorf("Tweak - num uniqueRules: %d != num got: %d",
				len(uniqueRules), numGot)
		}
		vars := map[string]*dlit.Literal{
			"fdMin": c.fdMin,
			"fdMax": c.fdMax,
			"min":   c.min,
			"max":   c.max,
		}
		for _, r := range got {
			switch x := r.(type) {
			case *BetweenFV:
				vars["minV"] = x.Min()
				vars["maxV"] = x.Max()
				if ok, err := inRangeExpr.EvalBool(vars); !ok || err != nil {
					t.Errorf("Tweak - invalid rule: %s", r)
				}
			default:
				t.Errorf("Tweak - invalid rule: %s", r)
			}
		}
	}
}

func TestBetweenFVOverlaps(t *testing.T) {
	cases := []struct {
		ruleA *BetweenFV
		ruleB Rule
		want  bool
	}{
		{ruleA: MustNewBetweenFV("rate", dlit.MustNew(5.3), dlit.MustNew(10.4)),
			ruleB: MustNewBetweenFV("rate", dlit.MustNew(1.3), dlit.MustNew(5.4)),
			want:  true,
		},
		{ruleA: MustNewBetweenFV("rate", dlit.MustNew(5.3), dlit.MustNew(10.4)),
			ruleB: MustNewBetweenFV("rate", dlit.MustNew(10.3), dlit.MustNew(15.4)),
			want:  true,
		},
		{ruleA: MustNewBetweenFV("rate", dlit.MustNew(5.3), dlit.MustNew(10.4)),
			ruleB: MustNewBetweenFV("rate", dlit.MustNew(6.3), dlit.MustNew(20.4)),
			want:  true,
		},
		{ruleA: MustNewBetweenFV("rate", dlit.MustNew(6.3), dlit.MustNew(20.4)),
			ruleB: MustNewBetweenFV("rate", dlit.MustNew(5.3), dlit.MustNew(10.4)),
			want:  true,
		},
		{ruleA: MustNewBetweenFV("rate", dlit.MustNew(5.3), dlit.MustNew(10.4)),
			ruleB: MustNewBetweenFV("rate", dlit.MustNew(1), dlit.MustNew(5.2)),
			want:  false,
		},
		{ruleA: MustNewBetweenFV("rate", dlit.MustNew(5.3), dlit.MustNew(10.4)),
			ruleB: MustNewBetweenFV("rate", dlit.MustNew(10.5), dlit.MustNew(20)),
			want:  false,
		},
		{ruleA: MustNewBetweenFV("rate", dlit.MustNew(5.3), dlit.MustNew(10.4)),
			ruleB: MustNewBetweenFV("flow", dlit.MustNew(1.3), dlit.MustNew(5.4)),
			want:  false,
		},
		{ruleA: MustNewBetweenFV("rate", dlit.MustNew(5.3), dlit.MustNew(10.4)),
			ruleB: NewLEFV("flow", dlit.MustNew(6.3)),
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

func TestGenerateBetweenFV(t *testing.T) {
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
				"income": {
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
				"income": {
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
				"income": {
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
				"income": {
					Kind:  fieldtype.Number,
					Min:   dlit.MustNew(700),
					Max:   dlit.MustNew(800),
					MaxDP: 0,
				},
				"month": {
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
			x, ok := r.(*BetweenFV)
			if !ok {
				return fmt.Errorf("wrong type: %T (%s)", r, r)
			}
			if x.field != c.field {
				return fmt.Errorf("field isn't correct for rule: %s", r)
			}
			return nil
		}
		got := generateBetweenFV(c.description, ruleFields, complexity, c.field)
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
			t.Errorf("(%d) generateBetweenFV: %s", i, err)
		}
	}
}
