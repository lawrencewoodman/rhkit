package rule

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/testhelpers"
)

func TestCountEQVFNew_panics(t *testing.T) {
	value := dlit.NewString("station")
	fields := []string{"place"}
	num := int64(3)
	wantPanic := "NewCountEQVF: Must contain at least two fields"
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("New() didn't panic")
		} else if r.(string) != wantPanic {
			t.Errorf("New() - got panic: %s, wanted: %s",
				r, wantPanic)
		}
	}()
	NewCountEQVF(value, fields, num)
}

func TestCountEQVFString(t *testing.T) {
	cases := []struct {
		rule Rule
		want string
	}{
		{rule: NewCountEQVF(dlit.MustNew(7.892), []string{"flowIn", "flowOut"}, 2),
			want: "count(\"7.892\", flowIn, flowOut) == 2",
		},
	}

	for _, c := range cases {
		got := c.rule.String()
		if got != c.want {
			t.Errorf("String() got: %s, want: %s", got, c.want)
		}
	}
}

func TestCountEQVFFields(t *testing.T) {
	fields := []string{"station", "destination"}
	value := dlit.MustNew(7.892)
	num := int64(2)
	r := NewCountEQVF(value, fields, num)
	wantFields := fields
	gotFields := r.Fields()
	if !reflect.DeepEqual(gotFields, wantFields) {
		t.Errorf("Fields() got: %v, want: %v", gotFields, wantFields)
	}
}

func TestCountEQVFIsTrue(t *testing.T) {
	cases := []struct {
		rule Rule
		want bool
	}{
		{rule: NewCountEQVF(
			dlit.MustNew("harry"),
			[]string{"station1", "station2", "station3"},
			int64(0),
		),
			want: false,
		},
		{rule: NewCountEQVF(
			dlit.MustNew("harry"),
			[]string{"station1", "station2", "station3"},
			int64(1),
		),
			want: true,
		},
		{rule: NewCountEQVF(
			dlit.MustNew("true"),
			[]string{"success1", "success2", "success3"},
			int64(0),
		),
			want: false,
		},
		{rule: NewCountEQVF(
			dlit.MustNew("true"),
			[]string{"success1", "success2", "success3"},
			int64(1),
		),
			want: true,
		},
		{rule: NewCountEQVF(
			dlit.MustNew("true"),
			[]string{"success1", "success2", "success3"},
			int64(2),
		),
			want: false,
		},
		{rule: NewCountEQVF(
			dlit.MustNew("true"),
			[]string{"success1", "success2", "success3"},
			int64(3),
		),
			want: false,
		},
		{rule: NewCountEQVF(
			dlit.MustNew("TRUE"),
			[]string{"success1", "success2", "success3"},
			int64(2),
		),
			want: true,
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
		"success3": dlit.MustNew("TRUE"),
		"band":     dlit.MustNew("alpha"),
	}
	for _, c := range cases {
		got, err := c.rule.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", c.rule, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t",
				c.rule, got, c.want)
		}
	}
}

func TestCountEQVFIsTrue_errors(t *testing.T) {
	cases := []struct {
		rule    Rule
		wantErr error
	}{
		{rule: NewCountEQVF(
			dlit.NewString("alpha"),
			[]string{"nonexistant", "income"},
			int64(1),
		),
			wantErr: InvalidRuleError{
				Rule: NewCountEQVF(
					dlit.NewString("alpha"),
					[]string{"nonexistant", "income"},
					int64(1),
				),
			},
		},
		{rule: NewCountEQVF(
			dlit.NewString("alpha"),
			[]string{"problem", "income"},
			int64(1),
		),
			wantErr: IncompatibleTypesRuleError{
				Rule: NewCountEQVF(
					dlit.NewString("alpha"),
					[]string{"problem", "income"},
					int64(1),
				),
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
		_, gotErr := c.rule.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", c.rule, err)
		}
	}
}

var wantGenerateCountEQVFRules = []Rule{
	NewCountEQVF(
		dlit.NewString("Fred"),
		[]string{"groupC", "groupD"},
		int64(0),
	),
	NewCountEQVF(
		dlit.NewString("Fred"),
		[]string{"groupC", "groupD"},
		int64(1),
	),
	NewCountEQVF(
		dlit.NewString("Fred"),
		[]string{"groupC", "groupD"},
		int64(2),
	),
	NewCountEQVF(
		dlit.NewString("Fred"),
		[]string{"groupA", "groupD"},
		int64(0),
	),
	NewCountEQVF(
		dlit.NewString("Fred"),
		[]string{"groupA", "groupD"},
		int64(1),
	),
	NewCountEQVF(
		dlit.NewString("Fred"),
		[]string{"groupA", "groupD"},
		int64(2),
	),
	NewCountEQVF(
		dlit.NewString("Fred"),
		[]string{"groupA", "groupC"},
		int64(0),
	),
	NewCountEQVF(
		dlit.NewString("Fred"),
		[]string{"groupA", "groupC"},
		int64(1),
	),
	NewCountEQVF(
		dlit.NewString("Fred"),
		[]string{"groupA", "groupC"},
		int64(2),
	),
	NewCountEQVF(
		dlit.NewString("Fred"),
		[]string{"groupA", "groupC", "groupD"},
		int64(0),
	),
	NewCountEQVF(
		dlit.NewString("Fred"),
		[]string{"groupA", "groupC", "groupD"},
		int64(1),
	),
	NewCountEQVF(
		dlit.NewString("Fred"),
		[]string{"groupA", "groupC", "groupD"},
		int64(2),
	),
	NewCountEQVF(
		dlit.NewString("Fred"),
		[]string{"groupA", "groupC", "groupD"},
		int64(3),
	),
	NewCountEQVF(
		dlit.NewString("Mary"),
		[]string{"groupC", "groupD"},
		int64(0),
	),
	NewCountEQVF(
		dlit.NewString("Mary"),
		[]string{"groupC", "groupD"},
		int64(1),
	),
	NewCountEQVF(
		dlit.NewString("Mary"),
		[]string{"groupC", "groupD"},
		int64(2),
	),
	NewCountEQVF(
		dlit.NewString("Mary"),
		[]string{"groupA", "groupD"},
		int64(0),
	),
	NewCountEQVF(
		dlit.NewString("Mary"),
		[]string{"groupA", "groupD"},
		int64(1),
	),
	NewCountEQVF(
		dlit.NewString("Mary"),
		[]string{"groupA", "groupD"},
		int64(2),
	),
	NewCountEQVF(
		dlit.NewString("Mary"),
		[]string{"groupA", "groupC"},
		int64(0),
	),
	NewCountEQVF(
		dlit.NewString("Mary"),
		[]string{"groupA", "groupC"},
		int64(1),
	),
	NewCountEQVF(
		dlit.NewString("Mary"),
		[]string{"groupA", "groupC"},
		int64(2),
	),
	NewCountEQVF(
		dlit.NewString("Mary"),
		[]string{"groupA", "groupC", "groupD"},
		int64(0),
	),
	NewCountEQVF(
		dlit.NewString("Mary"),
		[]string{"groupA", "groupC", "groupD"},
		int64(1),
	),
	NewCountEQVF(
		dlit.NewString("Mary"),
		[]string{"groupA", "groupC", "groupD"},
		int64(2),
	),
	NewCountEQVF(
		dlit.NewString("Mary"),
		[]string{"groupA", "groupC", "groupD"},
		int64(3),
	),
	NewCountEQVF(
		dlit.NewString("Rebecca"),
		[]string{"groupC", "groupD"},
		int64(0),
	),
	NewCountEQVF(
		dlit.NewString("Rebecca"),
		[]string{"groupC", "groupD"},
		int64(1),
	),
	NewCountEQVF(
		dlit.NewString("Rebecca"),
		[]string{"groupC", "groupD"},
		int64(2),
	),
	NewCountEQVF(
		dlit.NewString("Rebecca"),
		[]string{"groupA", "groupD"},
		int64(0),
	),
	NewCountEQVF(
		dlit.NewString("Rebecca"),
		[]string{"groupA", "groupD"},
		int64(1),
	),
	NewCountEQVF(
		dlit.NewString("Rebecca"),
		[]string{"groupA", "groupD"},
		int64(2),
	),
	NewCountEQVF(
		dlit.NewString("Rebecca"),
		[]string{"groupA", "groupC"},
		int64(0),
	),
	NewCountEQVF(
		dlit.NewString("Rebecca"),
		[]string{"groupA", "groupC"},
		int64(1),
	),
	NewCountEQVF(
		dlit.NewString("Rebecca"),
		[]string{"groupA", "groupC"},
		int64(2),
	),
	NewCountEQVF(
		dlit.NewString("Rebecca"),
		[]string{"groupA", "groupC", "groupD"},
		int64(0),
	),
	NewCountEQVF(
		dlit.NewString("Rebecca"),
		[]string{"groupA", "groupC", "groupD"},
		int64(1),
	),
	NewCountEQVF(
		dlit.NewString("Rebecca"),
		[]string{"groupA", "groupC", "groupD"},
		int64(2),
	),
	NewCountEQVF(
		dlit.NewString("Rebecca"),
		[]string{"groupA", "groupC", "groupD"},
		int64(3),
	),
	NewCountEQVF(
		dlit.NewString("Harry"),
		[]string{"groupC", "groupD"},
		int64(0),
	),
	NewCountEQVF(
		dlit.NewString("Harry"),
		[]string{"groupC", "groupD"},
		int64(1),
	),
	NewCountEQVF(
		dlit.NewString("Harry"),
		[]string{"groupC", "groupD"},
		int64(2),
	),
}

func TestGenerateCountEQVF(t *testing.T) {
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
				NumValues: 3,
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
				NumValues: 13,
			},
			"groupC": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 2},
					"Harry":   {dlit.NewString("Harry"), 2},
				},
				NumValues: 4,
			},
			"groupCa": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 2},
					"Harry":   {dlit.NewString("Harry"), 2},
				},
				NumValues: 4,
			},
			"groupD": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 1},
					"Harry":   {dlit.NewString("Harry"), 2},
				},
				NumValues: 4,
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
				NumValues: 5,
			},
		},
	}
	generationDesc := testhelpers.GenerationDesc{
		DFields: []string{"band", "flow", "groupA", "groupB",
			"groupC", "groupD", "groupE"},
		DArithmetic: false,
	}
	got := generateCountEQVF(inputDescription, generationDesc)
	if err := matchRulesUnordered(got, wantGenerateCountEQVFRules); err != nil {
		t.Errorf("matchRulesUnordered() rules don't match: %s\ngot: %s\nwant: %s\n",
			err, got, wantGenerateCountEQVFRules)
	}
}

// Test that will generate correct number of fields
func TestGenerateCountEQVF_num_fields(t *testing.T) {
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
				NumValues: 3,
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
				NumValues: 13,
			},
			"groupC": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 2},
					"Harry":   {dlit.NewString("Harry"), 2},
				},
				NumValues: 4,
			},
			"groupD": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 1},
					"Harry":   {dlit.NewString("Harry"), 2},
				},
				NumValues: 4,
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
				NumValues: 5,
			},
			"groupF": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 2},
					"Harry":   {dlit.NewString("Harry"), 2},
					"Juliet":  {dlit.NewString("Juliet"), 2},
				},
				NumValues: 5,
			},
			"groupG": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 1},
					"Harry":   {dlit.NewString("Harry"), 2},
				},
				NumValues: 4,
			},
			"groupH": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 1},
					"Harry":   {dlit.NewString("Harry"), 2},
				},
				NumValues: 4,
			},
			"groupI": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 1},
					"Harry":   {dlit.NewString("Harry"), 2},
				},
				NumValues: 4,
			},
			"groupJ": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 1},
					"Harry":   {dlit.NewString("Harry"), 2},
				},
				NumValues: 4,
			},
			"groupK": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 1},
					"Harry":   {dlit.NewString("Harry"), 2},
				},
				NumValues: 4,
			},
			"groupL": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 1},
					"Harry":   {dlit.NewString("Harry"), 2},
				},
				NumValues: 4,
			},
			"groupM": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 1},
					"Harry":   {dlit.NewString("Harry"), 2},
				},
				NumValues: 4,
			},
			"groupN": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 1},
					"Harry":   {dlit.NewString("Harry"), 2},
				},
				NumValues: 4,
			},
			"groupO": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 1},
					"Harry":   {dlit.NewString("Harry"), 2},
				},
				NumValues: 4,
			},
			"groupP": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 1},
					"Harry":   {dlit.NewString("Harry"), 2},
				},
				NumValues: 4,
			},
			"groupQ": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Fred":    {dlit.NewString("Fred"), 3},
					"Mary":    {dlit.NewString("Mary"), 4},
					"Rebecca": {dlit.NewString("Rebecca"), 1},
					"Harry":   {dlit.NewString("Harry"), 2},
				},
				NumValues: 4,
			},
		},
	}
	cases := []struct {
		generationDesc   GenerationDescriber
		wantMinNumRules  int
		wantMaxNumRules  int
		wantMaxNumFields int
	}{
		{generationDesc: testhelpers.GenerationDesc{
			DFields:     []string{"groupB", "groupC"},
			DArithmetic: false,
		},
			wantMinNumRules:  0,
			wantMaxNumRules:  0,
			wantMaxNumFields: 0,
		},
		{generationDesc: testhelpers.GenerationDesc{
			DFields:     []string{"groupE", "groupF"},
			DArithmetic: false,
		},
			wantMinNumRules:  0,
			wantMaxNumRules:  0,
			wantMaxNumFields: 0,
		},
		{generationDesc: testhelpers.GenerationDesc{
			DFields:     []string{"groupA", "groupC"},
			DArithmetic: false,
		},
			wantMinNumRules:  9,
			wantMaxNumRules:  9,
			wantMaxNumFields: 2,
		},
		{generationDesc: testhelpers.GenerationDesc{
			DFields:     []string{"groupA", "groupC", "groupG"},
			DArithmetic: false,
		},
			wantMinNumRules:  15,
			wantMaxNumRules:  50,
			wantMaxNumFields: 3,
		},
		{generationDesc: testhelpers.GenerationDesc{
			DFields:     []string{"groupA", "groupC", "groupG"},
			DArithmetic: false,
			DDeny:       map[string][]string{"CountEQFV": []string{"groupG"}},
		},
			wantMinNumRules:  9,
			wantMaxNumRules:  9,
			wantMaxNumFields: 2,
		},
		{generationDesc: testhelpers.GenerationDesc{
			DFields: []string{"groupA", "groupC", "groupG", "groupH",
				"groupI", "groupJ"},
			DArithmetic: false,
		},
			wantMinNumRules:  600,
			wantMaxNumRules:  900,
			wantMaxNumFields: 6,
		},
		{generationDesc: testhelpers.GenerationDesc{
			DFields: []string{"groupA", "groupC", "groupG", "groupH",
				"groupI", "groupJ", "groupK", "groupL"},
			DArithmetic: false,
		},
			wantMinNumRules:  3000,
			wantMaxNumRules:  4000,
			wantMaxNumFields: 8,
		},
		{generationDesc: testhelpers.GenerationDesc{
			DFields: []string{"groupA", "groupC", "groupG", "groupH",
				"groupI", "groupJ", "groupK", "groupL", "groupM"},
			DArithmetic: false,
		},
			wantMinNumRules:  3500,
			wantMaxNumRules:  4000,
			wantMaxNumFields: 9,
		},
		{generationDesc: testhelpers.GenerationDesc{
			DFields: []string{"groupA", "groupC", "groupG", "groupH",
				"groupI", "groupJ", "groupK", "groupL", "groupM", "groupN"},
			DArithmetic: false,
		},
			wantMinNumRules:  6000,
			wantMaxNumRules:  6100,
			wantMaxNumFields: 10,
		},
		{generationDesc: testhelpers.GenerationDesc{
			DFields: []string{"groupA", "groupC", "groupG", "groupH",
				"groupI", "groupJ", "groupK", "groupL", "groupM", "groupN", "groupO"},
			DArithmetic: false,
		},
			wantMinNumRules:  3000,
			wantMaxNumRules:  3100,
			wantMaxNumFields: 11,
		},
		{generationDesc: testhelpers.GenerationDesc{
			DFields: []string{"groupC", "groupD", "groupG", "groupH", "groupI",
				"groupJ", "groupK", "groupL", "groupM", "groupN", "groupO"},
			DArithmetic: false,
		},
			wantMinNumRules:  3300,
			wantMaxNumRules:  4000,
			wantMaxNumFields: 12,
		},
		{generationDesc: testhelpers.GenerationDesc{
			DFields: []string{"groupC", "groupD", "groupG", "groupH", "groupI",
				"groupJ", "groupK", "groupL", "groupM", "groupN", "groupO",
				"groupP"},
			DArithmetic: false,
		},
			wantMinNumRules:  4000,
			wantMaxNumRules:  4500,
			wantMaxNumFields: 12,
		},
		{generationDesc: testhelpers.GenerationDesc{
			DFields: []string{"groupC", "groupD", "groupG", "groupH", "groupI",
				"groupJ", "groupK", "groupL", "groupM", "groupN", "groupO",
				"groupP", "groupQ"},
			DArithmetic: false,
		},
			wantMinNumRules:  5000,
			wantMaxNumRules:  6000,
			wantMaxNumFields: 12,
		},
	}
	for i, c := range cases {
		got := generateCountEQVF(inputDescription, c.generationDesc)
		if len(got) < c.wantMinNumRules || len(got) > c.wantMaxNumRules {
			t.Errorf("(%d) generateCountEQVF: got wrong number of rules: %d",
				i, len(got))
		}
		for _, r := range got {
			numFields := strings.Count(r.String(), ",")
			if numFields < 2 || numFields > c.wantMaxNumFields {
				t.Errorf("(%d) generateCountEQVF: wrong number of values in rule: %s",
					i, r)
			}
		}
	}
}
