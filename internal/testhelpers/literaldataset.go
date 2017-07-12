/*************************
 *  Test helper functions
 *************************/
package testhelpers

import (
	"github.com/lawrencewoodman/ddataset"
	"github.com/lawrencewoodman/dlit"
)

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
) ddataset.Dataset {
	return &LiteralDataset{
		records:    records,
		fieldNames: fieldNames,
	}
}

func (l *LiteralDataset) Open() (ddataset.Conn, error) {
	return &LiteralDatasetConn{
		dataset:  l,
		position: -1,
		isClosed: false,
	}, nil
}

func (l *LiteralDataset) Fields() []string {
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

func (lc *LiteralDatasetConn) Read() ddataset.Record {
	line := lc.dataset.records[lc.position]
	record := make(ddataset.Record, len(lc.dataset.fieldNames))
	for i, v := range line {
		record[lc.dataset.fieldNames[i]] = dlit.MustNew(v)
	}
	return record
}

func (lc *LiteralDatasetConn) Err() error {
	return nil
}
