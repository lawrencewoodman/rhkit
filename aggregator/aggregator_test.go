package aggregator

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/goal"
	"reflect"
	"testing"
)

func TestRegister_panic1(t *testing.T) {
	paniced := false
	wantPanic := "aggregator.Register aggregator is nil"
	defer func() {
		if r := recover(); r != nil {
			if r.(string) == wantPanic {
				paniced = true
			} else {
				t.Errorf("Register: got panic: %s, wanted: %s", r, wantPanic)
			}
		}
	}()
	Register("dummy", nil)
	if !paniced {
		t.Errorf("Register: failed to panic with: %s", wantPanic)
	}
}

func TestRegister_panic2(t *testing.T) {
	paniced := false
	wantPanic := "aggregator.Register called twice for aggregator: goalsscore"
	defer func() {
		if r := recover(); r != nil {
			if r.(string) == wantPanic {
				paniced = true
			} else {
				t.Errorf("Register: got panic: %s, wanted: %s", r, wantPanic)
			}
		}
	}()
	Register("goalsscore", &goalsScoreAggregator{})
	if !paniced {
		t.Errorf("Register: failed to panic with: %s", wantPanic)
	}
}

func TestNew_error(t *testing.T) {
	cases := []struct {
		name    string
		kind    string
		args    []string
		wantErr error
	}{
		{name: "a",
			kind: "goalsscore",
			args: []string{"5+6"},
			wantErr: DescError{
				Name: "a",
				Kind: "goalsscore",
				Err:  ErrInvalidNumArgs,
			},
		},
		{name: "a",
			kind: "calc",
			args: []string{},
			wantErr: DescError{
				Name: "a",
				Kind: "calc",
				Err:  ErrInvalidNumArgs,
			},
		},
		{name: "a",
			kind: "calc",
			args: []string{"3+4", "5+6"},
			wantErr: DescError{
				Name: "a",
				Kind: "calc",
				Err:  ErrInvalidNumArgs,
			},
		},
		{name: "a",
			kind: "invalid",
			args: []string{"3+4"},
			wantErr: DescError{
				Name: "a",
				Kind: "invalid",
				Err:  ErrUnregisteredKind,
			},
		},
		{name: "2nums",
			kind: "calc",
			args: []string{"3+4"},
			wantErr: DescError{
				Name: "2nums",
				Kind: "calc",
				Err:  ErrInvalidName,
			},
		},
	}
	for i, c := range cases {
		_, err := New(c.name, c.kind, c.args...)
		if err == nil || err.Error() != c.wantErr.Error() {
			t.Errorf("(%d) New: gotErr: %s, wantErr: %s", i, err, c.wantErr)
		}
	}
}

func TestMustNew_panic(t *testing.T) {
	paniced := false
	wantPanic := DescError{
		Name: "a",
		Kind: "goalsscore",
		Err:  ErrInvalidNumArgs,
	}.Error()
	defer func() {
		if r := recover(); r != nil {
			if r.(error).Error() == wantPanic {
				paniced = true
			} else {
				t.Errorf("MustNew: got panic: %s, wanted: %s", r, wantPanic)
			}
		}
	}()
	got := MustNew("a", "goalsscore", "5")
	if !paniced {
		t.Errorf("MustNew: got: %s, failed to panic with: %s", got, wantPanic)
	}
}

func TestInstancesToMap(t *testing.T) {
	instances := []Instance{
		MustNewLitInstance("numSignedUp", "54"),
		MustNewLitInstance("profit", "54.25"),
		MustNewLitInstance("cost", "203"),
		MustNewLitInstance("income", "257.25"),
	}
	goals := []*goal.Goal{}
	numRecords := int64(100)
	cases := []struct {
		thisName string
		want     map[string]*dlit.Literal
	}{
		{"numSignedUp",
			map[string]*dlit.Literal{
				"numRecords": dlit.MustNew(100),
			},
		},
		{"profit",
			map[string]*dlit.Literal{
				"numRecords":  dlit.MustNew(100),
				"numSignedUp": dlit.MustNew(54),
			},
		},
		{"cost",
			map[string]*dlit.Literal{
				"numRecords":  dlit.MustNew(100),
				"numSignedUp": dlit.MustNew(54),
				"profit":      dlit.MustNew(54.25),
			},
		},
		{"income",
			map[string]*dlit.Literal{
				"numRecords":  dlit.MustNew(100),
				"numSignedUp": dlit.MustNew(54),
				"profit":      dlit.MustNew(54.25),
				"cost":        dlit.MustNew(203),
			},
		},
		{"",
			map[string]*dlit.Literal{
				"numRecords":  dlit.MustNew(100),
				"numSignedUp": dlit.MustNew(54),
				"profit":      dlit.MustNew(54.25),
				"cost":        dlit.MustNew(203),
				"income":      dlit.MustNew(257.25),
			},
		},
	}

	var gotResults map[string]*dlit.Literal
	var err error
	for _, c := range cases {
		if c.thisName == "" {
			gotResults, err =
				InstancesToMap(instances, goals, numRecords)
		} else {
			gotResults, err =
				InstancesToMap(instances, goals, numRecords, c.thisName)
		}
		if err != nil {
			t.Errorf("InstancesToMap(..., %s) err: %s", c.thisName, err)
		}

		if !doAggregatorMapsMatch(gotResults, c.want) {
			t.Errorf("InstancesToMap(..., %s) got: %s, want: %s",
				c.thisName, gotResults, c.want)
		}
	}
}

func TestMakeSpecs(t *testing.T) {
	fields := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "y"}
	desc := []*Desc{
		{"num_married", "count", "marital == \"married\""},
		{"numSignedUp", "count", "y == \"yes\""},
		{"cost", "calc", "numMatches * 4.5"},
		{"income", "calc", "numSignedUp * 24"},
		{"profit", "calc", "income - cost"},
	}
	want := []Spec{
		MustNew("numMatches", "count", "true()"),
		MustNew("percentMatches", "calc",
			"iferr(roundto(100.0 * numMatches / numRecords, 2), 0)"),
		// num_married to check for allowed characters
		MustNew("num_married", "count", "marital == \"married\""),
		MustNew("numSignedUp", "count", "y == \"yes\""),
		MustNew("cost", "calc", "numMatches * 4.5"),
		MustNew("income", "calc", "numSignedUp * 24"),
		MustNew("profit", "calc", "income - cost"),
		MustNew("goalsScore", "goalsscore"),
	}
	got, err := MakeSpecs(fields, desc)
	if err != nil {
		t.Errorf("MakeSpecs(%v): %s", desc, err)
	}
	if !areAggregatorsEqual(got, want) {
		t.Errorf("MakeSpecs(%v) got: %v, want: %v",
			desc, got, want)
	}
}

func TestMakeSpecs_errors(t *testing.T) {
	fields := []string{"age", "job", "marital", "education", "default",
		"balance", "housing", "loan", "contact", "day", "month", "duration",
		"campaign", "pdays", "previous", "y"}
	cases := []struct {
		desc    []*Desc
		wantErr error
	}{
		{desc: []*Desc{
			{"pdays", "count", "day > 2"},
		},
			wantErr: DescError{Name: "pdays", Kind: "count", Err: ErrNameClash},
		},
		{desc: []*Desc{
			{"numMatches", "count", "y == \"yes\""},
		},
			wantErr: DescError{
				Name: "numMatches",
				Kind: "count",
				Err:  ErrNameReserved,
			},
		},
		{desc: []*Desc{
			{"percentMatches", "percent", "y == \"yes\""},
		},
			wantErr: DescError{
				Name: "percentMatches",
				Kind: "percent",
				Err:  ErrNameReserved,
			},
		},
		{desc: []*Desc{
			{"goalsScore", "count", "y == \"yes\""},
		},
			wantErr: DescError{
				Name: "goalsScore",
				Kind: "count",
				Err:  ErrNameReserved,
			},
		},
		{desc: []*Desc{
			{"3numSignedUp", "count", "y == \"yes\""},
		},
			wantErr: DescError{
				Name: "3numSignedUp",
				Kind: "count",
				Err:  ErrInvalidName,
			},
		},
		{desc: []*Desc{
			{"num-signed-up", "count", "y == \"yes\""},
		},
			wantErr: DescError{
				Name: "num-signed-up",
				Kind: "count",
				Err:  ErrInvalidName,
			},
		},
		{desc: []*Desc{
			{"something", "nothing", "y == \"yes\""},
		},
			wantErr: DescError{
				Name: "something",
				Kind: "nothing",
				Err:  ErrUnregisteredKind,
			},
		},
	}
	for i, c := range cases {
		_, err := MakeSpecs(fields, c.desc)
		if err == nil || err.Error() != c.wantErr.Error() {
			t.Errorf("(%d) MakeSpecs: err: %s, wantErr: %s",
				i, err, c.wantErr)
		}
	}
}

/**********************
 *  Helper functions
 **********************/

func doAggregatorMapsMatch(am1, am2 map[string]*dlit.Literal) bool {
	if len(am1) != len(am2) {
		return false
	}
	for k, v := range am1 {
		if am2Val, ok := am2[k]; !ok || am2Val.String() != v.String() {
			return false
		}
	}
	return true
}

type LitInstance struct {
	name   string
	result string
}

func MustNewLitInstance(name, result string) *LitInstance {
	return &LitInstance{name: name, result: result}
}

func (li *LitInstance) Name() string {
	return li.name
}

func (li *LitInstance) NextRecord(
	record map[string]*dlit.Literal,
	isRuleTrue bool,
) error {
	return nil
}

func (li *LitInstance) Result(
	aggregatorInstances []Instance,
	goals []*goal.Goal,
	numRecords int64,
) *dlit.Literal {
	return dlit.MustNew(li.result)
}

func areAggregatorsEqual(a1 []Spec, a2 []Spec) bool {
	if len(a1) != len(a2) {
		return false
	}
	for i, a := range a1 {
		if reflect.TypeOf(a) != reflect.TypeOf(a2[i]) ||
			a.Name() != a2[i].Name() ||
			a.Arg() != a2[i].Arg() {
			return false
		}
	}
	return true
}
