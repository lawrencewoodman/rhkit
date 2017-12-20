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

func TestCountLTVFNew_panics(t *testing.T) {
	value := dlit.NewString("station")
	fields := []string{"place"}
	num := int64(3)
	wantPanic := "NewCountLTVF: Must contain at least two fields"
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("New() didn't panic")
		} else if r.(string) != wantPanic {
			t.Errorf("New() - got panic: %s, wanted: %s",
				r, wantPanic)
		}
	}()
	NewCountLTVF(value, fields, num)
}

func TestCountLTVFString(t *testing.T) {
	cases := []struct {
		rule Rule
		want string
	}{
		{rule: NewCountLTVF(dlit.MustNew(7.892), []string{"flowIn", "flowOut"}, 2),
			want: "count(\"7.892\", flowIn, flowOut) < 2",
		},
	}

	for _, c := range cases {
		got := c.rule.String()
		if got != c.want {
			t.Errorf("String() got: %s, want: %s", got, c.want)
		}
	}
}

func TestCountLTVFFields(t *testing.T) {
	fields := []string{"station", "destination"}
	value := dlit.MustNew(7.892)
	num := int64(2)
	r := NewCountLTVF(value, fields, num)
	wantFields := fields
	gotFields := r.Fields()
	if !reflect.DeepEqual(gotFields, wantFields) {
		t.Errorf("Fields() got: %v, want: %v", gotFields, wantFields)
	}
}

func TestCountLTVFIsTrue(t *testing.T) {
	cases := []struct {
		rule Rule
		want bool
	}{
		{rule: NewCountLTVF(
			dlit.MustNew("harry"),
			[]string{"station1", "station2", "station3"},
			int64(0),
		),
			want: false,
		},
		{rule: NewCountLTVF(
			dlit.MustNew("harry"),
			[]string{"station1", "station2", "station3"},
			int64(1),
		),
			want: false,
		},
		{rule: NewCountLTVF(
			dlit.MustNew("harry"),
			[]string{"station1", "station2", "station3"},
			int64(2),
		),
			want: true,
		},
		{rule: NewCountLTVF(
			dlit.MustNew("harry"),
			[]string{"station1", "station2", "station3"},
			int64(3),
		),
			want: true,
		},
		{rule: NewCountLTVF(
			dlit.MustNew("true"),
			[]string{"success1", "success2", "success3"},
			int64(0),
		),
			want: false,
		},
		{rule: NewCountLTVF(
			dlit.MustNew("true"),
			[]string{"success1", "success2", "success3"},
			int64(1),
		),
			want: false,
		},
		{rule: NewCountLTVF(
			dlit.MustNew("true"),
			[]string{"success1", "success2", "success3"},
			int64(2),
		),
			want: true,
		},
		{rule: NewCountLTVF(
			dlit.MustNew("true"),
			[]string{"success1", "success2", "success3"},
			int64(3),
		),
			want: true,
		},
		{rule: NewCountLTVF(
			dlit.MustNew("TRUE"),
			[]string{"success1", "success2", "success3"},
			int64(2),
		),
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

func TestCountLTVFIsTrue_errors(t *testing.T) {
	cases := []struct {
		rule    Rule
		wantErr error
	}{
		{rule: NewCountLTVF(
			dlit.NewString("alpha"),
			[]string{"nonexistant", "income"},
			int64(1),
		),
			wantErr: InvalidRuleError{
				Rule: NewCountLTVF(
					dlit.NewString("alpha"),
					[]string{"nonexistant", "income"},
					int64(1),
				),
			},
		},
		{rule: NewCountLTVF(
			dlit.NewString("alpha"),
			[]string{"problem", "income"},
			int64(1),
		),
			wantErr: IncompatibleTypesRuleError{
				Rule: NewCountLTVF(
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

var wantGenerateCountLTVFRules = []Rule{
	NewCountLTVF(
		dlit.NewString("Fred"),
		[]string{"groupC", "groupD"},
		int64(2),
	),
	NewCountLTVF(
		dlit.NewString("Fred"),
		[]string{"groupA", "groupD"},
		int64(2),
	),
	NewCountLTVF(
		dlit.NewString("Fred"),
		[]string{"groupA", "groupC"},
		int64(2),
	),
	NewCountLTVF(
		dlit.NewString("Fred"),
		[]string{"groupA", "groupC", "groupD"},
		int64(2),
	),
	NewCountLTVF(
		dlit.NewString("Fred"),
		[]string{"groupA", "groupC", "groupD"},
		int64(3),
	),
	NewCountLTVF(
		dlit.NewString("Mary"),
		[]string{"groupC", "groupD"},
		int64(2),
	),
	NewCountLTVF(
		dlit.NewString("Mary"),
		[]string{"groupA", "groupD"},
		int64(2),
	),
	NewCountLTVF(
		dlit.NewString("Mary"),
		[]string{"groupA", "groupC"},
		int64(2),
	),
	NewCountLTVF(
		dlit.NewString("Mary"),
		[]string{"groupA", "groupC", "groupD"},
		int64(2),
	),
	NewCountLTVF(
		dlit.NewString("Mary"),
		[]string{"groupA", "groupC", "groupD"},
		int64(3),
	),
	NewCountLTVF(
		dlit.NewString("Rebecca"),
		[]string{"groupC", "groupD"},
		int64(2),
	),
	NewCountLTVF(
		dlit.NewString("Rebecca"),
		[]string{"groupA", "groupD"},
		int64(2),
	),
	NewCountLTVF(
		dlit.NewString("Rebecca"),
		[]string{"groupA", "groupC"},
		int64(2),
	),
	NewCountLTVF(
		dlit.NewString("Rebecca"),
		[]string{"groupA", "groupC", "groupD"},
		int64(2),
	),
	NewCountLTVF(
		dlit.NewString("Rebecca"),
		[]string{"groupA", "groupC", "groupD"},
		int64(3),
	),
	NewCountLTVF(
		dlit.NewString("Harry"),
		[]string{"groupC", "groupD"},
		int64(2),
	),
}

func TestGenerateCountLTVF(t *testing.T) {
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
		},
	}
	generationDesc := testhelpers.GenerationDesc{
		DFields: []string{"band", "flow", "groupA", "groupB",
			"groupC", "groupD", "groupE"},
		DArithmetic: false,
	}
	got := generateCountLTVF(inputDescription, generationDesc)
	if err := matchRulesUnordered(got, wantGenerateCountLTVFRules); err != nil {
		t.Errorf("matchRulesUnordered() rules don't match: %s\ngot: %s\nwant: %s\n",
			err, got, wantGenerateCountLTVFRules)
	}
}

// Test that will generate correct number of fields
func TestGenerateCountLTVF_num_fields(t *testing.T) {
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
			wantMinNumRules:  3,
			wantMaxNumRules:  3,
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
			DFields: []string{"groupA", "groupC", "groupG", "groupH",
				"groupI", "groupJ"},
			DArithmetic: false,
		},
			wantMinNumRules:  400,
			wantMaxNumRules:  500,
			wantMaxNumFields: 6,
		},
		{generationDesc: testhelpers.GenerationDesc{
			DFields: []string{"groupA", "groupC", "groupG", "groupH",
				"groupI", "groupJ", "groupK", "groupL"},
			DArithmetic: false,
		},
			wantMinNumRules:  2000,
			wantMaxNumRules:  2050,
			wantMaxNumFields: 8,
		},
		{generationDesc: testhelpers.GenerationDesc{
			DFields: []string{"groupA", "groupC", "groupG", "groupH",
				"groupI", "groupJ", "groupK", "groupL", "groupM"},
			DArithmetic: false,
		},
			wantMinNumRules:  2050,
			wantMaxNumRules:  2100,
			wantMaxNumFields: 9,
		},
		{generationDesc: testhelpers.GenerationDesc{
			DFields: []string{"groupA", "groupC", "groupG", "groupH",
				"groupI", "groupJ", "groupK", "groupL", "groupM", "groupN"},
			DArithmetic: false,
		},
			wantMinNumRules:  3000,
			wantMaxNumRules:  4000,
			wantMaxNumFields: 10,
		},
		{generationDesc: testhelpers.GenerationDesc{
			DFields: []string{"groupA", "groupC", "groupG", "groupH",
				"groupI", "groupJ", "groupK", "groupL", "groupM", "groupN", "groupO"},
			DArithmetic: false,
		},
			wantMinNumRules:  1400,
			wantMaxNumRules:  1500,
			wantMaxNumFields: 11,
		},
		{generationDesc: testhelpers.GenerationDesc{
			DFields: []string{"groupC", "groupD", "groupG", "groupH", "groupI",
				"groupJ", "groupK", "groupL", "groupM", "groupN", "groupO"},
			DArithmetic: false,
		},
			wantMinNumRules:  1500,
			wantMaxNumRules:  1600,
			wantMaxNumFields: 12,
		},
		{generationDesc: testhelpers.GenerationDesc{
			DFields: []string{"groupC", "groupD", "groupG", "groupH", "groupI",
				"groupJ", "groupK", "groupL", "groupM", "groupN", "groupO",
				"groupP"},
			DArithmetic: false,
		},
			wantMinNumRules:  2000,
			wantMaxNumRules:  2050,
			wantMaxNumFields: 12,
		},
		{generationDesc: testhelpers.GenerationDesc{
			DFields: []string{"groupC", "groupD", "groupG", "groupH", "groupI",
				"groupJ", "groupK", "groupL", "groupM", "groupN", "groupO",
				"groupP", "groupQ"},
			DArithmetic: false,
		},
			wantMinNumRules:  2050,
			wantMaxNumRules:  3000,
			wantMaxNumFields: 12,
		},
	}
	for i, c := range cases {
		got := generateCountLTVF(inputDescription, c.generationDesc)
		if len(got) < c.wantMinNumRules || len(got) > c.wantMaxNumRules {
			t.Errorf("(%d) generateCountLTVF: got wrong number of rules: %d",
				i, len(got))
		}
		for _, r := range got {
			numFields := strings.Count(r.String(), ",")
			if numFields < 2 || numFields > c.wantMaxNumFields {
				t.Errorf("(%d) generateCountLTVF: wrong number of values in rule: %s",
					i, r)
			}
		}
	}
}
