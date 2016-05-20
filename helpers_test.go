/*************************
 *  Test helper functions
 *************************/
package rulehunter

import (
	"fmt"
	"github.com/lawrencewoodman/dexpr"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/input"
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
	records    [][]string
	fieldNames []string
	position   int
	isClosed   bool
}

func NewLiteralInput(fieldNames []string, records [][]string) input.Input {
	return &LiteralInput{records: records, fieldNames: fieldNames, position: -1}
}

func (l *LiteralInput) Clone() (input.Input, error) {
	return NewLiteralInput(l.fieldNames, l.records), nil
}

func (l *LiteralInput) Close() error {
	return nil
}

func (l *LiteralInput) Next() bool {
	if !l.isClosed && (l.position+1) < len(l.records) {
		l.position++
		return true
	}
	return false
}

func (l *LiteralInput) Read() (map[string]*dlit.Literal, error) {
	line := l.records[l.position]
	record := make(map[string]*dlit.Literal, len(l.fieldNames))
	for i, v := range line {
		record[l.fieldNames[i]] = dlit.MustNew(v)
	}
	return record, nil
}

func (l *LiteralInput) Err() error {
	return nil
}

func (l *LiteralInput) Rewind() error {
	l.position = -1
	return nil
}

func (l *LiteralInput) GetFieldNames() []string {
	return l.fieldNames
}
