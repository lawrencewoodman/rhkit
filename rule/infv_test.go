package rule

import (
	"errors"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/testhelpers"
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
		{ruleA: NewInFV("band", []*dlit.Literal{
			dlit.NewString("4"), dlit.NewString("3"), dlit.NewString("2")},
		),
			ruleB: NewInFV("band", []*dlit.Literal{
				dlit.NewString("9"), dlit.NewString("7")},
			),
			want: false,
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
			"band": {
				Kind:   description.Number,
				Min:    dlit.MustNew(1),
				Max:    dlit.MustNew(3),
				Values: map[string]description.Value{},
			},
			"flow": {
				Kind:   description.Number,
				Min:    dlit.MustNew(1),
				Max:    dlit.MustNew(3),
				MaxDP:  2,
				Values: map[string]description.Value{},
			},
			"groupA": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 2},
				},
			},

			"groupB": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 2},
					"Harry":   {dlit.NewString("Harry"), 2},
					"Dinah":   {dlit.NewString("Dinah"), 2},
					"Israel":  {dlit.NewString("Israel"), 2},
					"Sarah":   {dlit.NewString("Sarah"), 2},
					"Ishmael": {dlit.NewString("Ishmael"), 2},
					"Caen":    {dlit.NewString("Caen"), 2},
					"Abel":    {dlit.NewString("Abel"), 2},
					"Noah":    {dlit.NewString("Noah"), 2},
					"Isaac":   {dlit.NewString("Isaac"), 2},
					"Moses":   {dlit.NewString("Moses"), 2},
				},
			},
			"groupC": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 2},
					"Harry":   {dlit.NewString("Harry"), 2},
				},
			},
			"groupD": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 1},
					"Harry":   {dlit.NewString("Harry"), 2},
				},
			},
			"groupE": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 2},
					"Harry":   {dlit.NewString("Harry"), 2},
					"Juliet":  {dlit.NewString("Juliet"), 2},
				},
			},
		},
	}
	want := []Rule{
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
	}
	generationDesc := testhelpers.GenerationDesc{
		DFields: []string{"band", "flow", "groupA", "groupB",
			"groupC", "groupD", "groupE"},
		DArithmetic: false,
	}
	got := generateInFV(inputDescription, generationDesc)
	if err := matchRulesUnordered(got, want); err != nil {
		t.Errorf("matchRulesUnordered() rules don't match: %s\ngot: %s\nwant: %s\n",
			err, got, want)
	}
}

// Test that will generate correct number of values in In based on number
// of fields
func TestGenerateInFV_num_fields(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"group": {
				Kind: description.String,
			},
			"flowA": {
				Kind: description.Number,
			},
			"flowB": {
				Kind: description.Number,
			},
			"flowC": {
				Kind: description.Number,
			},
		},
	}
	cases := []struct {
		generationDesc   GenerationDescriber
		groupValues      map[string]description.Value
		wantMinNumRules  int
		wantMaxNumRules  int
		wantMaxNumValues int
	}{
		{
			generationDesc: testhelpers.GenerationDesc{
				DFields:     []string{"group", "flowA", "flowB"},
				DArithmetic: false,
			},
			groupValues: map[string]description.Value{
				"Fred":    {dlit.NewString("Fred"), 3},
				"Mary":    {dlit.NewString("Mary"), 4},
				"Rebecca": {dlit.NewString("Rebecca"), 2},
				"Harry":   {dlit.NewString("Harry"), 2},
				"Dinah":   {dlit.NewString("Dinah"), 2},
				"Israel":  {dlit.NewString("Israel"), 2},
				"Sarah":   {dlit.NewString("Sarah"), 2},
				"Ishmael": {dlit.NewString("Ishmael"), 2},
				"Caen":    {dlit.NewString("Caen"), 2},
				"Abel":    {dlit.NewString("Abel"), 2},
				"Noah":    {dlit.NewString("Noah"), 2},
				"Isaac":   {dlit.NewString("Isaac"), 2},
			},
			wantMinNumRules:  1000,
			wantMaxNumRules:  2000,
			wantMaxNumValues: 5,
		},
		{
			generationDesc: testhelpers.GenerationDesc{
				DFields:     []string{"group", "flowA"},
				DArithmetic: false,
			},
			groupValues: map[string]description.Value{
				"Fred":    {dlit.NewString("Fred"), 3},
				"Mary":    {dlit.NewString("Mary"), 4},
				"Rebecca": {dlit.NewString("Rebecca"), 2},
				"Harry":   {dlit.NewString("Harry"), 2},
				"Dinah":   {dlit.NewString("Dinah"), 2},
				"Israel":  {dlit.NewString("Israel"), 2},
				"Sarah":   {dlit.NewString("Sarah"), 2},
				"Ishmael": {dlit.NewString("Ishmael"), 2},
				"Caen":    {dlit.NewString("Caen"), 2},
				"Abel":    {dlit.NewString("Abel"), 2},
				"Noah":    {dlit.NewString("Noah"), 2},
				"Isaac":   {dlit.NewString("Isaac"), 2},
			},
			wantMinNumRules:  2000,
			wantMaxNumRules:  4000,
			wantMaxNumValues: 8,
		},
		{
			generationDesc: testhelpers.GenerationDesc{
				DFields:     []string{"group", "flowA", "flowB"},
				DArithmetic: false,
				DDeny:       map[string][]string{"INFV": []string{"group"}},
			},
			groupValues:      map[string]description.Value{},
			wantMinNumRules:  0,
			wantMaxNumRules:  0,
			wantMaxNumValues: 0,
		},
	}
	for i, c := range cases {
		inputDescription.Fields["group"].Values = c.groupValues
		got := generateInFV(inputDescription, c.generationDesc)
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

func TestLiteralCombinations(t *testing.T) {
	cases := []struct {
		values []*dlit.Literal
		min    int
		max    int
		want   [][]*dlit.Literal
	}{
		{values: []*dlit.Literal{
			dlit.NewString("a"),
		},
			min:  2,
			max:  3,
			want: [][]*dlit.Literal{},
		},
		{values: []*dlit.Literal{
			dlit.NewString("a"),
			dlit.NewString("c"),
		},
			min: 2,
			max: 3,
			want: [][]*dlit.Literal{
				[]*dlit.Literal{dlit.NewString("a"), dlit.NewString("c")},
			},
		},
		{values: []*dlit.Literal{
			dlit.NewString("a"),
			dlit.NewString("b"),
			dlit.NewString("c"),
		},
			min: 2,
			max: 2,
			want: [][]*dlit.Literal{
				[]*dlit.Literal{dlit.NewString("b"), dlit.NewString("c")},
				[]*dlit.Literal{dlit.NewString("a"), dlit.NewString("c")},
				[]*dlit.Literal{dlit.NewString("a"), dlit.NewString("b")},
			},
		},
		{values: []*dlit.Literal{
			dlit.NewString("a"),
			dlit.NewString("b"),
			dlit.NewString("c"),
		},
			min: 2,
			max: 3,
			want: [][]*dlit.Literal{
				[]*dlit.Literal{dlit.NewString("b"), dlit.NewString("c")},
				[]*dlit.Literal{dlit.NewString("a"), dlit.NewString("c")},
				[]*dlit.Literal{dlit.NewString("a"), dlit.NewString("b")},
				[]*dlit.Literal{
					dlit.NewString("a"),
					dlit.NewString("b"),
					dlit.NewString("c"),
				},
			},
		},
		{values: []*dlit.Literal{
			dlit.NewString("a"),
			dlit.NewString("b"),
			dlit.NewString("c"),
		},
			min: 2,
			max: 4,
			want: [][]*dlit.Literal{
				[]*dlit.Literal{dlit.NewString("b"), dlit.NewString("c")},
				[]*dlit.Literal{dlit.NewString("a"), dlit.NewString("c")},
				[]*dlit.Literal{dlit.NewString("a"), dlit.NewString("b")},
				[]*dlit.Literal{
					dlit.NewString("a"),
					dlit.NewString("b"),
					dlit.NewString("c"),
				},
			},
		},
	}
	for _, c := range cases {
		got := literalCombinations(c.values, c.min, c.max)
		if len(got) != len(c.want) {
			t.Errorf("literalCombinations(%v, %d, %d) got: %v, want: %v",
				c.values, c.min, c.max, got, c.want)
			continue
		}
		for i, subset := range got {
			if len(subset) != len(c.want[i]) {
				t.Errorf("literalCombinations(%v, %d, %d) got: %v, want: %v",
					c.values, c.min, c.max, got, c.want)
				continue
			}
			for j, v := range subset {
				if v.String() != c.want[i][j].String() {
					t.Errorf("literalCombinations(%v, %d, %d) got: %v, want: %v",
						c.values, c.min, c.max, got, c.want)
				}
			}
		}
	}
}
