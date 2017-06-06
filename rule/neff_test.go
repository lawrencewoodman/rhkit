package rule

import (
	"errors"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
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
			"bandA": &description.Field{
				Kind:   fieldtype.Number,
				Min:    dlit.MustNew(1),
				Max:    dlit.MustNew(3),
				Values: map[string]description.Value{},
			},
			"groupA": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Nelson":      description.Value{dlit.NewString("Nelson"), 3},
					"Collingwood": description.Value{dlit.NewString("Collingwood"), 1},
					"Mountbatten": description.Value{dlit.NewString("Mountbatten"), 1},
					"Drake":       description.Value{dlit.NewString("Drake"), 2},
				},
			},
			"groupB": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Nelson":      description.Value{dlit.NewString("Nelson"), 3},
					"Mountbatten": description.Value{dlit.NewString("Mountbatten"), 1},
					"Drake":       description.Value{dlit.NewString("Drake"), 2},
				},
			},
			"groupC": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Nelson": description.Value{dlit.NewString("Nelson"), 3},
					"Drake":  description.Value{dlit.NewString("Drake"), 2},
				},
			},
			"groupD": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Drake": description.Value{dlit.NewString("Drake"), 2},
				},
			},
			"groupE": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Drake":       description.Value{dlit.NewString("Drake"), 2},
					"Chaucer":     description.Value{dlit.NewString("Chaucer"), 2},
					"Shakespeare": description.Value{dlit.NewString("Shakespeare"), 2},
					"Marlowe":     description.Value{dlit.NewString("Marlowe"), 2},
				},
			},
			"groupF": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Nelson":      description.Value{dlit.NewString("Nelson"), 3},
					"Drake":       description.Value{dlit.NewString("Drake"), 2},
					"Chaucer":     description.Value{dlit.NewString("Chaucer"), 2},
					"Shakespeare": description.Value{dlit.NewString("Shakespeare"), 2},
					"Marlowe":     description.Value{dlit.NewString("Marlowe"), 2},
				},
			},
			"bandB": &description.Field{
				Kind: fieldtype.Number,
				Min:  dlit.MustNew(1),
				Max:  dlit.MustNew(3),
				Values: map[string]description.Value{
					"1": description.Value{dlit.NewString("1"), 3},
					"2": description.Value{dlit.NewString("2"), 2},
					"3": description.Value{dlit.NewString("3"), 1},
				},
			},
			"bandC": &description.Field{
				Kind: fieldtype.Number,
				Min:  dlit.MustNew(2),
				Max:  dlit.MustNew(7),
				Values: map[string]description.Value{
					"7": description.Value{dlit.NewString("7"), 3},
					"2": description.Value{dlit.NewString("2"), 2},
					"6": description.Value{dlit.NewString("6"), 1},
				},
			},
			"bandD": &description.Field{
				Kind: fieldtype.Number,
				Min:  dlit.MustNew(2),
				Max:  dlit.MustNew(8),
				Values: map[string]description.Value{
					"3": description.Value{dlit.NewString("3"), 3},
					"2": description.Value{dlit.NewString("2"), 2},
					"8": description.Value{dlit.NewString("8"), 1},
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
	ruleFields :=
		[]string{
			"bandA", "groupA", "groupB", "groupC", "groupD",
			"groupE", "groupF", "bandB", "bandC", "bandD",
		}
	complexity := 10
	for _, c := range cases {
		got := generateNEFF(
			inputDescription,
			ruleFields,
			complexity,
			c.field,
		)
		if err := matchRulesUnordered(got, c.want); err != nil {
			t.Errorf("matchRulesUnordered() rules don't match: %s\ngot: %s\nwant: %s\n",
				err, got, c.want)
		}
	}
}
