package rule

import (
	"errors"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/testhelpers"
	"reflect"
	"testing"
)

func TestNEFFString(t *testing.T) {
	fieldA := "income"
	fieldB := "cost"
	want := "income != cost"
	r := NewNEFF(fieldA, fieldB)
	got := r.String()
	if got != want {
		t.Errorf("String() got: %s, want: %s", got, want)
	}
}

func TestNEFFIsTrue(t *testing.T) {
	cases := []struct {
		fieldA string
		fieldB string
		want   bool
	}{
		{"income", "cost", true},
		{"cost", "income", true},
		{"cost", "band", true},
		{"income", "income", false},
		{"flowIn", "flowOut", true},
		{"flowOut", "flowIn", true},
		{"flowIn", "flowIn", false},
		{"flowIn", "band", true},
		{"income", "flowIn", true},
		{"flowIn", "income", true},
		{"band", "band", false},
		{"band", "trueA", true},
		{"trueA", "trueB", true},
		{"trueA", "trueA", false},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"cost":    dlit.MustNew(20),
		"flowIn":  dlit.MustNew(124.564),
		"flowOut": dlit.MustNew(124.565),
		"band":    dlit.MustNew("alpha"),
		"trueA":   dlit.MustNew("true"),
		"trueB":   dlit.MustNew("TRUE"),
	}
	for _, c := range cases {
		r := NewNEFF(c.fieldA, c.fieldB)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestNEFFIsTrue_errors(t *testing.T) {
	cases := []struct {
		fieldA  string
		fieldB  string
		wantErr error
	}{
		{fieldA: "fred",
			fieldB:  "income",
			wantErr: InvalidRuleError{NewNEFF("fred", "income")},
		},
		{fieldA: "income",
			fieldB:  "fred",
			wantErr: InvalidRuleError{NewNEFF("income", "fred")},
		},
		{fieldA: "income",
			fieldB:  "problem",
			wantErr: IncompatibleTypesRuleError{NewNEFF("income", "problem")},
		},
		{fieldA: "problem",
			fieldB:  "income",
			wantErr: IncompatibleTypesRuleError{NewNEFF("problem", "income")},
		},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"flow":    dlit.MustNew(124.564),
		"band":    dlit.NewString("alpha"),
		"problem": dlit.MustNew(errors.New("this is an error")),
	}
	for _, c := range cases {
		r := NewNEFF(c.fieldA, c.fieldB)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}

func TestNEFFFields(t *testing.T) {
	r := NewNEFF("income", "cost")
	want := []string{"income", "cost"}
	got := r.Fields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Fields() got: %s, want: %s", got, want)
	}
}

func TestGenerateNEFF(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"bandA": {
				Kind:   description.Number,
				Min:    dlit.MustNew(1),
				Max:    dlit.MustNew(3),
				Values: map[string]description.Value{},
			},
			"groupA": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Nelson":      {dlit.NewString("Nelson"), 3},
					"Collingwood": {dlit.NewString("Collingwood"), 1},
					"Mountbatten": {dlit.NewString("Mountbatten"), 1},
					"Drake":       {dlit.NewString("Drake"), 2},
				},
			},
			"groupB": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Nelson":      {dlit.NewString("Nelson"), 3},
					"Mountbatten": {dlit.NewString("Mountbatten"), 1},
					"Drake":       {dlit.NewString("Drake"), 2},
				},
			},
			"groupC": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Nelson": {dlit.NewString("Nelson"), 3},
					"Drake":  {dlit.NewString("Drake"), 2},
				},
			},
			"groupD": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Drake": {dlit.NewString("Drake"), 2},
				},
			},
			"groupE": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Drake":       {dlit.NewString("Drake"), 2},
					"Chaucer":     {dlit.NewString("Chaucer"), 2},
					"Shakespeare": {dlit.NewString("Shakespeare"), 2},
					"Marlowe":     {dlit.NewString("Marlowe"), 2},
				},
			},
			"groupF": {
				Kind: description.String,
				Values: map[string]description.Value{
					"Nelson":      {dlit.NewString("Nelson"), 3},
					"Drake":       {dlit.NewString("Drake"), 2},
					"Chaucer":     {dlit.NewString("Chaucer"), 2},
					"Shakespeare": {dlit.NewString("Shakespeare"), 2},
					"Marlowe":     {dlit.NewString("Marlowe"), 2},
				},
			},
			"bandB": {
				Kind: description.Number,
				Min:  dlit.MustNew(1),
				Max:  dlit.MustNew(3),
				Values: map[string]description.Value{
					"1": {dlit.NewString("1"), 3},
					"2": {dlit.NewString("2"), 2},
					"3": {dlit.NewString("3"), 1},
				},
			},
			"bandC": {
				Kind: description.Number,
				Min:  dlit.MustNew(2),
				Max:  dlit.MustNew(7),
				Values: map[string]description.Value{
					"7": {dlit.NewString("7"), 3},
					"2": {dlit.NewString("2"), 2},
					"6": {dlit.NewString("6"), 1},
				},
			},
			"bandD": {
				Kind: description.Number,
				Min:  dlit.MustNew(2),
				Max:  dlit.MustNew(8),
				Values: map[string]description.Value{
					"3": {dlit.NewString("3"), 3},
					"2": {dlit.NewString("2"), 2},
					"8": {dlit.NewString("8"), 1},
				},
			},
		},
	}
	cases := []struct {
		field string
		want  []Rule
	}{
		{field: "bandA",
			want: []Rule{},
		},
		{field: "groupA",
			want: []Rule{
				NewNEFF("groupA", "groupB"),
				NewNEFF("groupA", "groupC"),
				NewNEFF("groupA", "groupF"),
			},
		},
		{field: "groupB",
			want: []Rule{
				NewNEFF("groupB", "groupC"),
				NewNEFF("groupB", "groupF"),
			},
		},
		{field: "groupC",
			want: []Rule{
				NewNEFF("groupC", "groupF"),
			},
		},
		{field: "groupD",
			want: []Rule{},
		},
		{field: "groupE",
			want: []Rule{
				NewNEFF("groupE", "groupF"),
			},
		},
		{field: "bandB",
			want: []Rule{
				NewNEFF("bandB", "bandD"),
			},
		},
	}
	generationDesc := testhelpers.GenerationDesc{
		DFields: []string{
			"bandA", "groupA", "groupB", "groupC", "groupD",
			"groupE", "groupF", "bandB", "bandC", "bandD",
		},
		DArithmetic: false,
	}
	for _, c := range cases {
		got := generateNEFF(inputDescription, generationDesc, c.field)
		if err := matchRulesUnordered(got, c.want); err != nil {
			t.Errorf("matchRulesUnordered() rules don't match: %s\ngot: %s\nwant: %s\n",
				err, got, c.want)
		}
	}
}
