package rule

import (
	"errors"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
	"reflect"
	"strings"
	"testing"
)

func TestInFVNew_panics(t *testing.T) {
	values := []*dlit.Literal{}
	wantPanic := "NewInFV: Must contain at least one value"
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("New() didn't panic")
		} else if r.(string) != wantPanic {
			t.Errorf("New() - got panic: %s, wanted: %s",
				r, wantPanic)
		}
	}()
	field := "station"
	NewInFV(field, values)
}

func TestInFVString(t *testing.T) {
	cases := []struct {
		values []*dlit.Literal
		want   string
	}{
		{values: []*dlit.Literal{
			dlit.MustNew(7.892),
			dlit.MustNew("harry"),
			dlit.MustNew(""),
			dlit.MustNew(" harry "),
			dlit.MustNew("fred and win"),
		},
			want: "in(station,\"7.892\",\"harry\",\"\",\" harry \",\"fred and win\")",
		},
	}

	field := "station"
	for _, c := range cases {
		r := NewInFV(field, c.values)
		got := r.String()
		if got != c.want {
			t.Errorf("String() got: %s, want: %s", got, c.want)
		}
	}
}

func TestInFVFields(t *testing.T) {
	field := "station"
	values := []*dlit.Literal{dlit.MustNew(7.892)}
	r := NewInFV(field, values)
	wantFields := []string{field}
	gotFields := r.Fields()
	if !reflect.DeepEqual(gotFields, wantFields) {
		t.Errorf("Fields() got: %v, want: %v", gotFields, wantFields)
	}
}

func TestInFVIsTrue(t *testing.T) {
	cases := []struct {
		field  string
		values []*dlit.Literal
		want   bool
	}{
		{field: "station1",
			values: []*dlit.Literal{
				dlit.MustNew(7.892),
				dlit.MustNew("harry"),
				dlit.MustNew(""),
				dlit.MustNew(" harry "),
				dlit.MustNew("fred and win"),
				dlit.MustNew("true"),
			},
			want: true,
		},
		{field: "station2",
			values: []*dlit.Literal{
				dlit.MustNew(7.892),
				dlit.MustNew("harry"),
				dlit.MustNew(""),
				dlit.MustNew(" harry "),
				dlit.MustNew("fred and win"),
				dlit.MustNew("true"),
			},
			want: true,
		},
		{field: "station3",
			values: []*dlit.Literal{
				dlit.MustNew(7.892),
				dlit.MustNew("harry"),
				dlit.MustNew(""),
				dlit.MustNew(" harry "),
				dlit.MustNew("fred and win"),
				dlit.MustNew("true"),
			},
			want: false,
		},
		{field: "flow1",
			values: []*dlit.Literal{
				dlit.MustNew(7.892),
				dlit.MustNew("harry"),
				dlit.MustNew(""),
				dlit.MustNew(" harry "),
				dlit.MustNew("fred and win"),
				dlit.MustNew("true"),
			},
			want: true,
		},
		{field: "flow2",
			values: []*dlit.Literal{
				dlit.MustNew(7.892),
				dlit.MustNew("harry"),
				dlit.MustNew(""),
				dlit.MustNew(" harry "),
				dlit.MustNew("fred and win"),
				dlit.MustNew("true"),
			},
			want: false,
		},
		{field: "success1",
			values: []*dlit.Literal{
				dlit.MustNew(7.892),
				dlit.MustNew("harry"),
				dlit.MustNew(""),
				dlit.MustNew(" harry "),
				dlit.MustNew("fred and win"),
				dlit.MustNew("true"),
			},
			want: true,
		},
		{field: "success2",
			values: []*dlit.Literal{
				dlit.MustNew(7.892),
				dlit.MustNew("harry"),
				dlit.MustNew(""),
				dlit.MustNew(" harry "),
				dlit.MustNew("fred and win"),
				dlit.MustNew("true"),
			},
			want: false,
		},
	}
	record := map[string]*dlit.Literal{
		"station1": dlit.MustNew("harry"),
		"station2": dlit.MustNew(" harry "),
		"station3": dlit.MustNew("  harry  "),
		"flow1":    dlit.MustNew(7.892),
		"flow2":    dlit.MustNew(7.893),
		"success1": dlit.MustNew("true"),
		"success2": dlit.MustNew("TRUE"),
		"band":     dlit.MustNew("alpha"),
	}
	for _, c := range cases {
		r := NewInFV(c.field, c.values)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestInFVIsTrue_errors(t *testing.T) {
	cases := []struct {
		field   string
		values  []*dlit.Literal
		wantErr error
	}{
		{field: "fred",
			values: []*dlit.Literal{dlit.NewString("hello")},
			wantErr: InvalidRuleError{
				Rule: NewInFV("fred", []*dlit.Literal{dlit.NewString("hello")}),
			},
		},
		{field: "problem",
			values: []*dlit.Literal{dlit.NewString("hello")},
			wantErr: IncompatibleTypesRuleError{
				Rule: NewInFV("problem", []*dlit.Literal{dlit.NewString("hello")}),
			},
		},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"flow":    dlit.MustNew(124.564),
		"band":    dlit.NewString("alpha"),
		"problem": dlit.MustNew(errors.New("this is an error")),
	}
	for _, c := range cases {
		r := NewInFV(c.field, c.values)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}

func TestInFVOverlaps(t *testing.T) {
	cases := []struct {
		ruleA *InFV
		ruleB Rule
		want  bool
	}{
		{ruleA: NewInFV("band", []*dlit.Literal{
			dlit.NewString("4"), dlit.NewString("3"), dlit.NewString("2")},
		),
			ruleB: NewInFV("band", []*dlit.Literal{
				dlit.NewString("9"), dlit.NewString("2")},
			),
			want: true,
		},
		{ruleA: NewInFV("band", []*dlit.Literal{
			dlit.NewString("9"), dlit.NewString("2")},
		),
			ruleB: NewInFV("band", []*dlit.Literal{
				dlit.NewString("4"), dlit.NewString("3"), dlit.NewString("2")},
			),
			want: true,
		},
		{ruleA: NewInFV("rate", []*dlit.Literal{
			dlit.NewString("4"), dlit.NewString("3"), dlit.NewString("2")},
		),
			ruleB: NewInFV("band", []*dlit.Literal{
				dlit.NewString("9"), dlit.NewString("2")},
			),
			want: false,
		},
		{ruleA: NewInFV("band", []*dlit.Literal{
			dlit.NewString("4"), dlit.NewString("3"), dlit.NewString("2")},
		),
			ruleB: NewLEFV("band", dlit.MustNew(6)),
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

func TestGenerateInFV(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"band": &description.Field{
				Kind:   fieldtype.Number,
				Min:    dlit.MustNew(1),
				Max:    dlit.MustNew(3),
				Values: map[string]description.Value{},
			},
			"flow": &description.Field{
				Kind:   fieldtype.Number,
				Min:    dlit.MustNew(1),
				Max:    dlit.MustNew(3),
				MaxDP:  2,
				Values: map[string]description.Value{},
			},
			"groupA": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Fred":    description.Value{dlit.NewString("Fred"), 3},
					"Mary":    description.Value{dlit.NewString("Mary"), 4},
					"Rebecca": description.Value{dlit.NewString("Rebecca"), 2},
				},
			},

			"groupB": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Fred":    description.Value{dlit.NewString("Fred"), 3},
					"Mary":    description.Value{dlit.NewString("Mary"), 4},
					"Rebecca": description.Value{dlit.NewString("Rebecca"), 2},
					"Harry":   description.Value{dlit.NewString("Harry"), 2},
					"Dinah":   description.Value{dlit.NewString("Dinah"), 2},
					"Israel":  description.Value{dlit.NewString("Israel"), 2},
					"Sarah":   description.Value{dlit.NewString("Sarah"), 2},
					"Ishmael": description.Value{dlit.NewString("Ishmael"), 2},
					"Caen":    description.Value{dlit.NewString("Caen"), 2},
					"Abel":    description.Value{dlit.NewString("Abel"), 2},
					"Noah":    description.Value{dlit.NewString("Noah"), 2},
					"Isaac":   description.Value{dlit.NewString("Isaac"), 2},
					"Moses":   description.Value{dlit.NewString("Moses"), 2},
				},
			},
			"groupC": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Fred":    description.Value{dlit.NewString("Fred"), 3},
					"Mary":    description.Value{dlit.NewString("Mary"), 4},
					"Rebecca": description.Value{dlit.NewString("Rebecca"), 2},
					"Harry":   description.Value{dlit.NewString("Harry"), 2},
				},
			},
			"groupD": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Fred":    description.Value{dlit.NewString("Fred"), 3},
					"Mary":    description.Value{dlit.NewString("Mary"), 4},
					"Rebecca": description.Value{dlit.NewString("Rebecca"), 1},
					"Harry":   description.Value{dlit.NewString("Harry"), 2},
				},
			},
			"groupE": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Fred":    description.Value{dlit.NewString("Fred"), 3},
					"Mary":    description.Value{dlit.NewString("Mary"), 4},
					"Rebecca": description.Value{dlit.NewString("Rebecca"), 2},
					"Harry":   description.Value{dlit.NewString("Harry"), 2},
					"Juliet":  description.Value{dlit.NewString("Juliet"), 2},
				},
			},
		},
	}
	cases := []struct {
		field string
		want  []Rule
	}{
		{field: "band",
			want: []Rule{},
		},
		{field: "flow",
			want: []Rule{},
		},
		{field: "groupA",
			want: []Rule{},
		},
		{field: "groupB",
			want: []Rule{},
		},
		{field: "groupC",
			want: []Rule{
				NewInFV("groupC", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Harry"),
				}),
				NewInFV("groupC", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Mary"),
				}),
				NewInFV("groupC", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Rebecca"),
				}),
				NewInFV("groupC", []*dlit.Literal{
					dlit.NewString("Harry"),
					dlit.NewString("Mary"),
				}),
				NewInFV("groupC", []*dlit.Literal{
					dlit.NewString("Harry"),
					dlit.NewString("Rebecca"),
				}),
				NewInFV("groupC", []*dlit.Literal{
					dlit.NewString("Mary"),
					dlit.NewString("Rebecca"),
				}),
			},
		},
		{field: "groupD",
			want: []Rule{
				NewInFV("groupD", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Harry"),
				}),
				NewInFV("groupD", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Mary"),
				}),
				NewInFV("groupD", []*dlit.Literal{
					dlit.NewString("Harry"),
					dlit.NewString("Mary"),
				}),
			},
		},
		{field: "groupE",
			want: []Rule{
				NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Harry"),
				}),
				NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Harry"),
					dlit.NewString("Juliet"),
				}),
				NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Harry"),
					dlit.NewString("Mary"),
				}),
				NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Harry"),
					dlit.NewString("Rebecca"),
				}),
				NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Juliet"),
					dlit.NewString("Rebecca"),
				}),
				NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Juliet"),
					dlit.NewString("Mary"),
				}),
				NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Mary"),
				}),
				NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Mary"),
					dlit.NewString("Rebecca"),
				}),
				NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Rebecca"),
				}),
				NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Fred"),
					dlit.NewString("Juliet"),
				}),
				NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Harry"),
					dlit.NewString("Mary"),
				}),
				NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Harry"),
					dlit.NewString("Rebecca"),
				}),
				NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Harry"),
					dlit.NewString("Juliet"),
				}),
				NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Harry"),
					dlit.NewString("Juliet"),
					dlit.NewString("Mary"),
				}),
				NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Harry"),
					dlit.NewString("Juliet"),
					dlit.NewString("Rebecca"),
				}),
				NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Harry"),
					dlit.NewString("Mary"),
					dlit.NewString("Rebecca"),
				}),
				NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Juliet"),
					dlit.NewString("Mary"),
				}),
				NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Juliet"),
					dlit.NewString("Mary"),
					dlit.NewString("Rebecca"),
				}),
				NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Mary"),
					dlit.NewString("Rebecca"),
				}),
				NewInFV("groupE", []*dlit.Literal{
					dlit.NewString("Juliet"),
					dlit.NewString("Rebecca"),
				}),
			},
		},
	}
	ruleFields :=
		[]string{"band", "flow", "groupA", "groupB", "groupC", "groupD", "groupE"}
	complexity := Complexity{}
	for i, c := range cases {
		got := generateInFV(inputDescription, ruleFields, complexity, c.field)
		if err := matchRulesUnordered(got, c.want); err != nil {
			t.Errorf("(%d) matchRulesUnordered() rules don't match: %s\ngot: %s\nwant: %s\n",
				i, err, got, c.want)
		}
	}
}

// Test that will generate correct number of values in In based on number
// of fields
func TestGenerateInFV_num_fields(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"group": &description.Field{
				Kind: fieldtype.String,
			},
			"flowA": &description.Field{
				Kind: fieldtype.Number,
			},
			"flowB": &description.Field{
				Kind: fieldtype.Number,
			},
		},
	}
	cases := []struct {
		ruleFields       []string
		groupValues      map[string]description.Value
		wantMinNumRules  int
		wantMaxNumRules  int
		wantMaxNumValues int
	}{
		{ruleFields: []string{"group", "flowA", "flowB"},
			groupValues: map[string]description.Value{
				"Fred":    description.Value{dlit.NewString("Fred"), 3},
				"Mary":    description.Value{dlit.NewString("Mary"), 4},
				"Rebecca": description.Value{dlit.NewString("Rebecca"), 2},
				"Harry":   description.Value{dlit.NewString("Harry"), 2},
				"Dinah":   description.Value{dlit.NewString("Dinah"), 2},
				"Israel":  description.Value{dlit.NewString("Israel"), 2},
				"Sarah":   description.Value{dlit.NewString("Sarah"), 2},
				"Ishmael": description.Value{dlit.NewString("Ishmael"), 2},
				"Caen":    description.Value{dlit.NewString("Caen"), 2},
				"Abel":    description.Value{dlit.NewString("Abel"), 2},
				"Noah":    description.Value{dlit.NewString("Noah"), 2},
				"Isaac":   description.Value{dlit.NewString("Isaac"), 2},
			},
			wantMinNumRules:  1000,
			wantMaxNumRules:  2000,
			wantMaxNumValues: 5,
		},
		{ruleFields: []string{"group", "flowA"},
			groupValues: map[string]description.Value{
				"Fred":    description.Value{dlit.NewString("Fred"), 3},
				"Mary":    description.Value{dlit.NewString("Mary"), 4},
				"Rebecca": description.Value{dlit.NewString("Rebecca"), 2},
				"Harry":   description.Value{dlit.NewString("Harry"), 2},
				"Dinah":   description.Value{dlit.NewString("Dinah"), 2},
				"Israel":  description.Value{dlit.NewString("Israel"), 2},
				"Sarah":   description.Value{dlit.NewString("Sarah"), 2},
				"Ishmael": description.Value{dlit.NewString("Ishmael"), 2},
				"Caen":    description.Value{dlit.NewString("Caen"), 2},
				"Abel":    description.Value{dlit.NewString("Abel"), 2},
				"Noah":    description.Value{dlit.NewString("Noah"), 2},
				"Isaac":   description.Value{dlit.NewString("Isaac"), 2},
			},
			wantMinNumRules:  2000,
			wantMaxNumRules:  4000,
			wantMaxNumValues: 8,
		},
	}
	complexity := Complexity{}
	for i, c := range cases {
		inputDescription.Fields["group"].Values = c.groupValues
		got := generateInFV(inputDescription, c.ruleFields, complexity, "group")
		if len(got) < c.wantMinNumRules || len(got) > c.wantMaxNumRules {
			t.Errorf("(%d) generateInFV: got wrong number of rules: %d", i, len(got))
		}
		for _, r := range got {
			numValues := strings.Count(r.String(), ",")
			if numValues < 2 || numValues > c.wantMaxNumValues {
				t.Errorf("(%d) generateInFV: wrong number of values in rule: %s", i, r)
			}
		}
	}
}

func TestMakeCompareValues(t *testing.T) {
	values1 := map[string]description.Value{
		"a": description.Value{dlit.MustNew("a"), 2},
		"c": description.Value{dlit.MustNew("c"), 2},
		"d": description.Value{dlit.MustNew("d"), 2},
		"e": description.Value{dlit.MustNew("e"), 2},
		"f": description.Value{dlit.MustNew("f"), 2},
	}
	values2 := map[string]description.Value{
		"a": description.Value{dlit.MustNew("a"), 2},
		"c": description.Value{dlit.MustNew("c"), 1},
		"d": description.Value{dlit.MustNew("d"), 2},
		"e": description.Value{dlit.MustNew("e"), 2},
		"f": description.Value{dlit.MustNew("f"), 2},
	}
	cases := []struct {
		values map[string]description.Value
		i      int
		want   []*dlit.Literal
	}{
		{
			values: values1,
			i:      2,
			want:   []*dlit.Literal{dlit.NewString("c")},
		},
		{
			values: values2,
			i:      2,
			want:   []*dlit.Literal{},
		},
		{
			values: values1,
			i:      3,
			want:   []*dlit.Literal{dlit.NewString("a"), dlit.NewString("c")},
		},
		{
			values: values2,
			i:      3,
			want:   []*dlit.Literal{},
		},
		{
			values: values1,
			i:      4,
			want:   []*dlit.Literal{dlit.NewString("d")},
		},
		{
			values: values2,
			i:      4,
			want:   []*dlit.Literal{dlit.NewString("d")},
		},
		{
			values: values1,
			i:      5,
			want:   []*dlit.Literal{dlit.NewString("a"), dlit.NewString("d")},
		},
		{
			values: values2,
			i:      5,
			want:   []*dlit.Literal{dlit.NewString("a"), dlit.NewString("d")},
		},
		{
			values: values1,
			i:      6,
			want:   []*dlit.Literal{dlit.NewString("c"), dlit.NewString("d")},
		},
		{
			values: values2,
			i:      6,
			want:   []*dlit.Literal{},
		},
		{
			values: values1,
			i:      7,
			want: []*dlit.Literal{
				dlit.NewString("a"),
				dlit.NewString("c"),
				dlit.NewString("d"),
			},
		},
		{
			values: values2,
			i:      7,
			want:   []*dlit.Literal{},
		},
		{
			values: values1,
			i:      8,
			want:   []*dlit.Literal{dlit.NewString("e")},
		},
		{
			values: values2,
			i:      8,
			want:   []*dlit.Literal{dlit.NewString("e")},
		},
		{
			values: values1,
			i:      9,
			want:   []*dlit.Literal{dlit.NewString("a"), dlit.NewString("e")},
		},
		{
			values: values2,
			i:      9,
			want:   []*dlit.Literal{dlit.NewString("a"), dlit.NewString("e")},
		},
		{
			values: values1,
			i:      10,
			want:   []*dlit.Literal{dlit.NewString("c"), dlit.NewString("e")},
		},
		{
			values: values2,
			i:      10,
			want:   []*dlit.Literal{},
		},
		{
			values: values1,
			i:      11,
			want: []*dlit.Literal{
				dlit.NewString("a"),
				dlit.NewString("c"),
				dlit.NewString("e"),
			},
		},
		{
			values: values2,
			i:      11,
			want:   []*dlit.Literal{},
		},
		{
			values: values1,
			i:      12,
			want:   []*dlit.Literal{dlit.NewString("d"), dlit.NewString("e")},
		},
		{
			values: values2,
			i:      12,
			want:   []*dlit.Literal{dlit.NewString("d"), dlit.NewString("e")},
		},
		{
			values: values1,
			i:      13,
			want: []*dlit.Literal{
				dlit.NewString("a"),
				dlit.NewString("d"),
				dlit.NewString("e"),
			},
		},
		{
			values: values1,
			i:      14,
			want: []*dlit.Literal{
				dlit.NewString("c"),
				dlit.NewString("d"),
				dlit.NewString("e"),
			},
		},
		{
			values: values1,
			i:      15,
			want: []*dlit.Literal{
				dlit.NewString("a"),
				dlit.NewString("c"),
				dlit.NewString("d"),
				dlit.NewString("e"),
			},
		},
		{
			values: values1,
			i:      16,
			want:   []*dlit.Literal{dlit.NewString("f")},
		},
		{
			values: values2,
			i:      16,
			want:   []*dlit.Literal{dlit.NewString("f")},
		},
		{
			values: values1,
			i:      17,
			want:   []*dlit.Literal{dlit.NewString("a"), dlit.NewString("f")},
		},
		{
			values: values2,
			i:      17,
			want:   []*dlit.Literal{dlit.NewString("a"), dlit.NewString("f")},
		},
	}
	for _, c := range cases {
		got := makeCompareValues(c.values, c.i)
		if len(got) != len(c.want) {
			t.Errorf("makeCompareValues(%v, %d) got: %v, want: %v",
				c.values, c.i, got, c.want)
		}
		for j, v := range got {
			o := c.want[j]
			if o.String() != v.String() {
				t.Errorf("makeCompareValues(%v, %d) got: %v, want: %v",
					c.values, c.i, got, c.want)
			}
		}
	}
}
