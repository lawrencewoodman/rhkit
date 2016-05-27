/*************************
 *  Test helper functions
 *************************/
package rulehunter

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/dataset"
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

type LiteralDataset struct {
	records    [][]string
	fieldNames []string
	position   int
	isClosed   bool
}

func NewLiteralDataset(
	fieldNames []string,
	records [][]string,
) dataset.Dataset {
	return &LiteralDataset{records: records, fieldNames: fieldNames, position: -1}
}

func (l *LiteralDataset) Clone() (dataset.Dataset, error) {
	return NewLiteralDataset(l.fieldNames, l.records), nil
}

func (l *LiteralDataset) Close() error {
	return nil
}

func (l *LiteralDataset) Next() bool {
	if !l.isClosed && (l.position+1) < len(l.records) {
		l.position++
		return true
	}
	return false
}

func (l *LiteralDataset) Read() (map[string]*dlit.Literal, error) {
	line := l.records[l.position]
	record := make(map[string]*dlit.Literal, len(l.fieldNames))
	for i, v := range line {
		record[l.fieldNames[i]] = dlit.MustNew(v)
	}
	return record, nil
}

func (l *LiteralDataset) Err() error {
	return nil
}

func (l *LiteralDataset) Rewind() error {
	l.position = -1
	return nil
}

func (l *LiteralDataset) GetFieldNames() []string {
	return l.fieldNames
}
