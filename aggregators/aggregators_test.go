package aggregators

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/goal"
	"testing"
)

func TestAggregatorsToMap(t *testing.T) {
	aggregators := []Aggregator{
		MustNewLitAggregator("numSignedUp", "54"),
		MustNewLitAggregator("profit", "54.25"),
		MustNewLitAggregator("cost", "203"),
		MustNewLitAggregator("income", "257.25"),
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
				AggregatorsToMap(aggregators, goals, numRecords)
		} else {
			gotResults, err =
				AggregatorsToMap(aggregators, goals, numRecords, c.thisName)
		}
		if err != nil {
			t.Errorf("AggregatorsToMap(..., %s) err: %s", c.thisName, err)
		}

		if !doAggregatorMapsMatch(gotResults, c.want) {
			t.Errorf("AggregatorsToMap(..., %s) got: %s, want: %s",
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

type LitAggregator struct {
	name   string
	result string
}

func MustNewLitAggregator(name, result string) *LitAggregator {
	return &LitAggregator{name: name, result: result}
}

func (l *LitAggregator) CloneNew() Aggregator {
	MustNewLitAggregator(l.name, l.result)
	return nil
}

func (l *LitAggregator) GetName() string {
	return l.name
}

func (l *LitAggregator) GetArg() string {
	return ""
}

func (l *LitAggregator) NextRecord(
	record map[string]*dlit.Literal,
	isRuleTrue bool,
) error {
	return nil
}

func (l *LitAggregator) GetResult(
	aggregators []Aggregator,
	goals []*goal.Goal,
	numRecords int64,
) *dlit.Literal {
	return dlit.MustNew(l.result)
}

func (l *LitAggregator) IsEqual(o Aggregator) bool {
	if _, ok := o.(*LitAggregator); !ok {
		return false
	}
	return l.name == o.GetName() && l.GetArg() == o.GetArg()
}
