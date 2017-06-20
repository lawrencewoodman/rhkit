package rule

import (
	"errors"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
	"math"
	"reflect"
	"testing"
)

func TestNEFVString(t *testing.T) {
	cases := []struct {
		field string
		value *dlit.Literal
		want  string
	}{
		{field: "income", value: dlit.MustNew(7.8903), want: "income != 7.8903"},
		{field: "income", value: dlit.MustNew(7.890300), want: "income != 7.8903"},
		{field: "income", value: dlit.MustNew(7.), want: "income != 7"},
		{field: "income", value: dlit.MustNew(7.00), want: "income != 7"},
		{field: "income", value: dlit.MustNew(7), want: "income != 7"},
		{field: "income", value: dlit.MustNew(0.34), want: "income != 0.34"},
		{field: "income", value: dlit.MustNew(0.3400), want: "income != 0.34"},
		{field: "income", value: dlit.MustNew(0.0), want: "income != 0"},
		{field: "income", value: dlit.MustNew(-53.9), want: "income != -53.9"},
		{field: "name", value: dlit.MustNew("borris"),
			want: "name != \"borris\""},
		{field: "name", value: dlit.MustNew("bo   rris"),
			want: "name != \"bo   rris\""},
		{field: "name", value: dlit.MustNew("  borris  "),
			want: "name != \"  borris  \""},
		{field: "name", value: dlit.MustNew(""), want: "name != \"\""},
	}
	for _, c := range cases {
		r := NewNEFV(c.field, c.value)
		got := r.String()
		if got != c.want {
			t.Errorf("String() got: %s, want: %s", got, c.want)
		}
	}
}

func TestNEFVIsTrue(t *testing.T) {
	cases := []struct {
		field string
		value *dlit.Literal
		want  bool
	}{
		{"income", dlit.MustNew(19.0), false},
		{"income", dlit.MustNew(math.MaxInt64), true},
		{"income", dlit.MustNew(-19.0), true},
		{"income", dlit.MustNew(20.0), true},
		{"flow", dlit.MustNew(124.564), false},
		{"flow", dlit.MustNew(-124.564), true},
		{"flow", dlit.MustNew(20.0), true},
		{"flow", dlit.MustNew(124.5645), true},
		{"flow", dlit.MustNew(124.565), true},
		{"band", dlit.MustNew("hello"), true},
		{"band", dlit.MustNew("alpha"), false},
		{"band", dlit.MustNew("ALPHA"), true},
		{"band", dlit.MustNew(8.9), true},
		{"success", dlit.MustNew("TRUE"), false},
		{"success", dlit.MustNew("true"), true},
		{"success", dlit.MustNew("1"), true},
		{"bigNums", dlit.MustNew(math.MaxInt64), false},
		{"bigNums", dlit.MustNew("1"), true},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"flow":    dlit.MustNew(124.564),
		"band":    dlit.MustNew("alpha"),
		"success": dlit.MustNew("TRUE"),
		"bigNums": dlit.MustNew(math.MaxInt64),
	}
	for _, c := range cases {
		r := NewNEFV(c.field, c.value)
		got, err := r.IsTrue(record)
		if err != nil {
			t.Errorf("IsTrue(record) rule: %s, err: %v", r, err)
		}
		if got != c.want {
			t.Errorf("IsTrue(record) (rule: %s) got: %t, want: %t", r, got, c.want)
		}
	}
}

func TestNEFVIsTrue_errors(t *testing.T) {
	cases := []struct {
		field   string
		value   *dlit.Literal
		wantErr error
	}{
		{"fred", dlit.MustNew(8.9),
			InvalidRuleError{Rule: NewNEFV("fred", dlit.MustNew(8.9))}},
		{"problem", dlit.MustNew(8.9),
			IncompatibleTypesRuleError{Rule: NewNEFV("problem", dlit.MustNew(8.9))}},
	}
	record := map[string]*dlit.Literal{
		"income":  dlit.MustNew(19),
		"flow":    dlit.MustNew(124.564),
		"band":    dlit.NewString("alpha"),
		"problem": dlit.MustNew(errors.New("this is an error")),
	}
	for _, c := range cases {
		r := NewNEFV(c.field, c.value)
		_, gotErr := r.IsTrue(record)
		if err := checkErrorMatch(gotErr, c.wantErr); err != nil {
			t.Errorf("IsTrue(record) rule: %s - %s", r, err)
		}
	}
}

func TestNEFVFields(t *testing.T) {
	r := NewNEFV("income", dlit.MustNew(5.5))
	want := []string{"income"}
	got := r.Fields()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Fields() got: %s, want: %s", got, want)
	}
}

func TestGenerateNEFV(t *testing.T) {
	inputDescription := &description.Description{
		map[string]*description.Field{
			"band": &description.Field{
				Kind: fieldtype.Number,
				Min:  dlit.MustNew(1),
				Max:  dlit.MustNew(4),
				Values: map[string]description.Value{
					"1": description.Value{dlit.NewString("1"), 3},
					"2": description.Value{dlit.NewString("2"), 1},
					"3": description.Value{dlit.NewString("3"), 2},
					"4": description.Value{dlit.NewString("4"), 5},
				},
			},
			"flow": &description.Field{
				Kind:  fieldtype.Number,
				Min:   dlit.MustNew(1),
				Max:   dlit.MustNew(4),
				MaxDP: 2,
				Values: map[string]description.Value{
					"1":    description.Value{dlit.NewString("1"), 3},
					"2":    description.Value{dlit.NewString("2"), 1},
					"2.90": description.Value{dlit.NewString("2.90"), 1},
					"3.37": description.Value{dlit.NewString("3.37"), 2},
					"3.3":  description.Value{dlit.NewString("3.3"), 2},
					"4":    description.Value{dlit.NewString("4"), 5},
				},
			},
			"group": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"Nelson":      description.Value{dlit.NewString("Nelson"), 3},
					"Collingwood": description.Value{dlit.NewString("Collingwood"), 1},
					"Mountbatten": description.Value{dlit.NewString("Mountbatten"), 1},
					"Drake":       description.Value{dlit.NewString("Drake"), 2},
				},
			},
			"month": &description.Field{
				Kind: fieldtype.String,
				Values: map[string]description.Value{
					"May":  description.Value{dlit.NewString("May"), 3},
					"June": description.Value{dlit.NewString("June"), 2},
				},
			},
		},
	}
	cases := []struct {
		field string
		want  []Rule
	}{
		{field: "band",
			want: []Rule{
				NewNEFV("band", dlit.MustNew(1)),
				NewNEFV("band", dlit.MustNew(3)),
				NewNEFV("band", dlit.MustNew(4)),
			},
		},
		{field: "flow",
			want: []Rule{
				NewNEFV("flow", dlit.MustNew(1)),
				NewNEFV("flow", dlit.MustNew(3.37)),
				NewNEFV("flow", dlit.MustNew(3.3)),
				NewNEFV("flow", dlit.MustNew(4)),
			},
		},
		{field: "group",
			want: []Rule{
				NewNEFV("group", dlit.MustNew("Nelson")),
				NewNEFV("group", dlit.MustNew("Drake")),
			},
		},
		{field: "month",
			want: []Rule{},
		},
	}
	ruleFields := []string{"band", "flow", "group", "month"}
	complexity := Complexity{}
	for _, c := range cases {
		got := generateNEFV(inputDescription, ruleFields, complexity, c.field)
		if err := matchRulesUnordered(got, c.want); err != nil {
			t.Errorf("matchRulesUnordered() rules don't match: %s\ngot: %s\nwant: %s\n",
				err, got, c.want)
		}
	}
}
