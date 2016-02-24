package main

import (
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"testing"
)

func TestAssess(t *testing.T) {
	cases := []struct {
		goal        *dexpr.Expr
		aggregators []Aggregator
		wantPassed  bool
		wantErr     error
	}{
		{mustNewDExpr("totalIncome > 5000"),
			[]Aggregator{&DummyAggregator{"totalIncome", mustNewLit(5000)}},
			false, nil},
		{mustNewDExpr("totalIncome > 5000"),
			[]Aggregator{&DummyAggregator{"totalIncome", mustNewLit(5001)}},
			true, nil},
		{mustNewDExpr("totalCosts < 5000"),
			[]Aggregator{&DummyAggregator{"totalIncome", mustNewLit(9000)}},
			false,
			dexpr.ErrInvalidExpr("Variable doesn't exist: totalCosts")},
	}
	numRecords := int64(12)
	for _, c := range cases {
		gotPassed, err := HasGoalPassed(c.goal, c.aggregators, numRecords)
		if !errorMatch(c.wantErr, err) {
			t.Errorf("HasGoalPassed(%q, %q) expr: %s - err: %q, wantErr: %q",
				c.goal, c.aggregators, err, c.wantErr)
		}
		if gotPassed != c.wantPassed {
			t.Errorf("HasGoalPassed(%q, %q) want: %s, got: %s",
				c.goal, c.aggregators, c.wantPassed, gotPassed)
		}
	}
}

type DummyAggregator struct {
	name   string
	result *dlit.Literal
}

func (d *DummyAggregator) CloneNew() Aggregator {
	return &DummyAggregator{name: d.name, result: d.result}
}

func (d *DummyAggregator) GetName() string {
	return d.name
}

func (d *DummyAggregator) GetArg() string {
	return ""
}

func (d *DummyAggregator) NextRecord(record map[string]*dlit.Literal,
	isRuleTrue bool) error {
	return nil
}

func (d *DummyAggregator) GetResult(aggregators []Aggregator,
	numRecords int64) *dlit.Literal {
	return d.result
}

func (a *DummyAggregator) IsEqual(o Aggregator) bool {
	if _, ok := o.(*DummyAggregator); !ok {
		return false
	}
	return a.name == o.GetName()
}
