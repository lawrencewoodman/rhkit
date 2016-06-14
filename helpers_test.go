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

type LiteralDatasetConn struct {
	dataset  *LiteralDataset
	position int
	isClosed bool
}

func NewLiteralDataset(
	fieldNames []string,
	records [][]string,
) dataset.Dataset {
	return &LiteralDataset{
		records:    records,
		fieldNames: fieldNames,
	}
}

func (l *LiteralDataset) Open() (dataset.Conn, error) {
	return &LiteralDatasetConn{
		dataset:  l,
		position: -1,
		isClosed: false,
	}, nil
}

func (l *LiteralDataset) GetFieldNames() []string {
	return l.fieldNames
}

func (lc *LiteralDatasetConn) Close() error {
	return nil
}

func (lc *LiteralDatasetConn) Next() bool {
	if !lc.isClosed && (lc.position+1) < len(lc.dataset.records) {
		lc.position++
		return true
	}
	return false
}

func (lc *LiteralDatasetConn) Read() dataset.Record {
	line := lc.dataset.records[lc.position]
	record := make(dataset.Record, len(lc.dataset.fieldNames))
	for i, v := range line {
		record[lc.dataset.fieldNames[i]] = dlit.MustNew(v)
	}
	return record
}

func (lc *LiteralDatasetConn) Err() error {
	return nil
}
