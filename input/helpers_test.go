/*************************
 *  Test helper functions
 *************************/
package input

import (
	"github.com/lawrencewoodman/dlit"
)

type LiteralInput struct {
	records    [][]string
	fieldNames []string
	position   int
	isClosed   bool
}

func NewLiteralInput(fieldNames []string, records [][]string) Input {
	return &LiteralInput{records: records, fieldNames: fieldNames, position: -1}
}

func (l *LiteralInput) Clone() (Input, error) {
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
