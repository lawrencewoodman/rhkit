/*************************
 *  Test helper functions
 *************************/
package main

import (
	"fmt"
	"github.com/lawrencewoodman/dexpr_go"
	"github.com/lawrencewoodman/dlit_go"
	"github.com/lawrencewoodman/rulehunter/internal"
	"io"
)

func errorMatch(e1 error, e2 error) bool {
	if e1 == nil && e2 == nil {
		return true
	}
	if e1 == nil || e2 == nil {
		return false
	}
	if e1.Error() == e2.Error() {
		return true
	}
	return false
}

func mustNewDExpr(expr string) *dexpr.Expr {
	dexpr, err := dexpr.New(expr)
	if err != nil {
		panic(fmt.Sprintf("Can't create dexpr.Expr: %q", err))
	}
	return dexpr
}

func mustNewCountAggregator(
	name string,
	expr string,
) *internal.CountAggregator {
	c, err := internal.NewCountAggregator(name, expr)
	if err != nil {
		panic(fmt.Sprintf("Can't create CountAggregator: %s", err))
	}
	return c
}

func mustNewCalcAggregator(
	name string,
	expr string,
) *internal.CalcAggregator {
	c, err := internal.NewCalcAggregator(name, expr)
	if err != nil {
		panic(fmt.Sprintf("Can't create CalcAggregator: %s", err))
	}
	return c
}

func matchRules(rules1 []string, rules2 []string) (bool, string) {
	if len(rules1) != len(rules2) {
		return false, "rules different lengths"
	}
	for _, rule1 := range rules1 {
		found := false
		for _, rule2 := range rules2 {
			if rule1 == rule2 {
				found = true
				break
			}
		}
		if !found {
			return false, fmt.Sprintf("rule doesn't exist: %s", rule1)
		}
	}
	return true, ""
}

type LiteralInput struct {
	records  []map[string]*dlit.Literal
	position int
}

func NewLiteralInput(records []map[string]*dlit.Literal) internal.Input {
	return &LiteralInput{records: records, position: 0}
}

func (l *LiteralInput) Read() (map[string]*dlit.Literal, error) {
	if l.position < len(l.records) {
		record := l.records[l.position]
		l.position++
		return record, nil
	}
	return map[string]*dlit.Literal{}, io.EOF
}

func (l *LiteralInput) Rewind() error {
	l.position = 0
	return nil
}
