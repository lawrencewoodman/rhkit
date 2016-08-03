package aggregators

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/goal"
	"testing"
)

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

func (li *LitInstance) GetName() string {
	return li.name
}

func (li *LitInstance) NextRecord(
	record map[string]*dlit.Literal,
	isRuleTrue bool,
) error {
	return nil
}

func (li *LitInstance) GetResult(
	aggregatorInstances []AggregatorInstance,
	goals []*goal.Goal,
	numRecords int64,
) *dlit.Literal {
	return dlit.MustNew(li.result)
}
