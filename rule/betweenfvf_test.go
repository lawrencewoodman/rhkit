package rule

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
	"reflect"
	"testing"
)

func TestNewBetweenFVF(t *testing.T) {
	min := float64(5.89)
	max := float64(6.72)
	r, err := NewBetweenFVF("flow", min, max)
	if err != nil {
		t.Errorf("NewBetweenFVF(%s, %f, %f) got err: %s", "flow", min, max, err)
	}
	if r == nil {
		t.Errorf("NewBetweenFVF(%s, %f, %f) got r: nil", "flow", min, max)
	}
}

func TestNewBetweenFVF_errors(t *testing.T) {
	cases := []struct {
		min        float64
		max        float64
		wantErrStr string
	}{
		{min: 5,
			max:        5,
			wantErrStr: "can't create Between rule where max: 5 <= min: 5",
		},
		{min: 5.65,
			max:        5.65,
			wantErrStr: "can't create Between rule where max: 5.65 <= min: 5.65",
		},
		{min: 6.72,
			max:        5.89,
			wantErrStr: "can't create Between rule where max: 5.89 <= min: 6.72",
		},
	}
	field := "flow"
	for _, c := range cases {
		r, err := NewBetweenFVF(field, c.min, c.max)
		if r != nil {
			t.Errorf("NewBetweenFVF(%s, %f, %f) rule got: %s, want: nil",
				field, c.min, c.max, r)
		}
		if err == nil {
			t.Errorf("NewBetweenFVF(%s, %f, %f) got err: nil, want: %s",
				field, c.min, c.max, c.wantErrStr)
		} else if err.Error() != c.wantErrStr {
			t.Errorf("NewBetweenFVF(%s, %f, %f) got err: %s, want: %s",
				field, c.min, c.max, err, c.wantErrStr)
		}
	}
}

func TestBetweenFVFString(t *testing.T) {
	field := "flow"
	min := float64(183.92837)
	max := float64(287.87442)
	want := "flow >= 183.92837 && flow <= 287.87442"
	r, err := NewBetweenFVF(field, min, max)
	if err != nil {
		t.Fatalf("NewBetweenFVF: %s", err)
	}
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestBetweenFVFIsTrue(t *testing.T) {
	cases := []struct {
		field string
		min   float64
		max   float64
		want  bool
	}{
		{field: "income", min: 18.3, max: 20.45, want: true},
		{field: "income", min: 19.19, max: 20.78, want: true},
		{field: "income", min: 19.81, max: 19.83, want: true},
		{field: "income", min: 18.78, max: 19.92, want: true},
		{field: "income", min: 10.12, max: 25.986, want: true},
		{field: "income", min: 10.34, max: 19.81, want: false},
		{field: "income", min: 19.83, max: 30.5, want: false},
		{field: "cost", min: 20.67, max: 30.89, want: true},
		{field: "cost", min: 20.23, max: 25.98, want: true},
		{field: "cost", min: 25.98, max: 30.7, want: true},
		{field: "cost", min: 20.2, max: 25.97, want: false},
		{field: "cost", min: 25.99, max: 30.7, want: false},
	}
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19.82),
		"cost":   dlit.MustNew(25.98),
	}
	for _, c := range cases {
		r, err := NewBetweenFVF(c.field, c.min, c.max)
		if err != nil {
			t.Fatalf("NewBetweenFVF: %s", err)
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

func TestBetweenFVFIsTrue_errors(t *testing.T) {
	field := "rate"
	min := float64(18.72)
	max := float64(20.64)
	record := map[string]*dlit.Literal{
		"income": dlit.MustNew(19.27),
		"cost":   dlit.MustNew(18.34),
		"band":   dlit.NewString("alpha"),
	}
	r, err := NewBetweenFVF(field, min, max)
	if err != nil {
		t.Fatalf("NewBetweenFVF: %s", err)
	}
	wantErr := InvalidRuleError{Rule: r}
	_, err = r.IsTrue(record)
	if err != wantErr {
		t.Errorf("IsTrue(record) rule: %s, err: %v, want: %v", r, err, wantErr)
	}
}

func TestBetweenFVFGetFields(t *testing.T) {
	field := "rate"
	min := float64(18.72)
	max := float64(20.72)
	want := []string{"rate"}
	r, err := NewBetweenFVF(field, min, max)
	if err != nil {
		t.Fatalf("NewBetweenFVF: %s", err)
	}
	got := r.GetFields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetFields() got: %s, want: %s", got, want)
	}
}

func TestBetweenFVFTweak(t *testing.T) {
	field := "income"
	cases := []struct {
		fdMin   float64
		fdMax   float64
		min     float64
		max     float64
		wantNum int
	}{
		{fdMin: float64(500),
			fdMax:   float64(2000),
			min:     float64(800),
			max:     float64(1000),
			wantNum: 150,
		},
		{fdMin: float64(18),
			fdMax:   float64(95),
			min:     float64(22),
			max:     float64(26),
			wantNum: 150,
		},
	}
	complexity := 10
	for _, c := range cases {
		description := &description.Description{
			map[string]*description.Field{
				"income": &description.Field{
					Kind:  fieldtype.Float,
					Min:   dlit.MustNew(c.fdMin),
					Max:   dlit.MustNew(c.fdMax),
					MaxDP: 2,
				},
			},
		}
		rule := MustNewBetweenFVF(field, c.min, c.max)
		got := rule.Tweak(description, complexity, 1)
		numGot := len(got)
		if numGot < c.wantNum {
			t.Errorf("Tweak - got too few rules returned: %d", numGot)
		}
		uniqueRules := Uniq(got)
		if len(uniqueRules) != numGot {
			t.Errorf("Tweak - num uniqueRules: %d != num got: %d",
				len(uniqueRules), numGot)
		}
		for _, r := range got {
			switch x := r.(type) {
			case *BetweenFVF:
				minV := x.GetMin()
				maxV := x.GetMax()
				if minV <= c.fdMin ||
					maxV >= c.fdMax ||
					minV == c.min ||
					maxV == c.max {
					t.Errorf("Tweak - invalid rule: %s", r)
				}
			default:
				t.Errorf("Tweak - invalid rule: %s", r)
			}
		}
	}
}

func TestBetweenFVFOverlaps(t *testing.T) {
	cases := []struct {
		ruleA *BetweenFVF
		ruleB Rule
		want  bool
	}{
		{ruleA: MustNewBetweenFVF("rate", 5.3, 10.4),
			ruleB: MustNewBetweenFVF("rate", 1.3, 5.4),
			want:  true,
		},
		{ruleA: MustNewBetweenFVF("rate", 5.3, 10.4),
			ruleB: MustNewBetweenFVF("rate", 10.3, 15.4),
			want:  true,
		},
		{ruleA: MustNewBetweenFVF("rate", 5.3, 10.4),
			ruleB: MustNewBetweenFVF("rate", 6.3, 20.4),
			want:  true,
		},
		{ruleA: MustNewBetweenFVF("rate", 6.3, 20.4),
			ruleB: MustNewBetweenFVF("rate", 5.3, 10.4),
			want:  true,
		},
		{ruleA: MustNewBetweenFVF("rate", 5.3, 10.4),
			ruleB: MustNewBetweenFVF("rate", 1, 5.2),
			want:  false,
		},
		{ruleA: MustNewBetweenFVF("rate", 5.3, 10.4),
			ruleB: MustNewBetweenFVF("rate", 10.5, 20),
			want:  false,
		},
		{ruleA: MustNewBetweenFVF("rate", 5.3, 10.4),
			ruleB: MustNewBetweenFVF("flow", 1.3, 5.4),
			want:  false,
		},
		{ruleA: MustNewBetweenFVF("rate", 5.3, 10.4),
			ruleB: NewLEFVF("flow", 6.3),
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
