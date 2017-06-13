package aggregators

import (
	"errors"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/goal"
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
		aggType string
		args    []string
		wantErr error
	}{
		{aggType: "goalsscore",
			args:    []string{"5+6"},
			wantErr: errors.New("invalid number of arguments for aggregator: goalsscore"),
		},
		{aggType: "calc",
			args:    []string{},
			wantErr: errors.New("invalid number of arguments for aggregator: calc"),
		},
		{aggType: "calc",
			args:    []string{"3+4", "5+6"},
			wantErr: errors.New("invalid number of arguments for aggregator: calc"),
		},
	}
	for i, c := range cases {
		_, err := New("a", c.aggType, c.args...)
		if err == nil || err.Error() != c.wantErr.Error() {
			t.Errorf("(%d) New: gotErr: %s, wantErr: %s", i, err, c.wantErr)
		}
	}
}

func TestMustNew_panic(t *testing.T) {
	paniced := false
	wantPanic := "invalid number of arguments for aggregator: goalsscore"
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
	instances := []AggregatorInstance{
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
	aggregatorInstances []AggregatorInstance,
	goals []*goal.Goal,
	numRecords int64,
) *dlit.Literal {
	return dlit.MustNew(li.result)
}
