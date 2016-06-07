/*
	Copyright (C) 2016 vLife Systems Ltd <http://vlifesystems.com>
	This file is part of Rulehunter.

	Rulehunter is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	Rulehunter is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with Rulehunter; see the file COPYING.  If not, see
	<http://www.gnu.org/licenses/>.
*/
package reducedataset

import (
	"github.com/vlifesystems/rulehunter/dataset"
	"io"
)

type ReduceDataset struct {
	dataset    dataset.Dataset
	numRecords int
}

type ReduceDatasetConn struct {
	dataset   *ReduceDataset
	conn      dataset.Conn
	recordNum int
	err       error
}

func New(dataset dataset.Dataset, numRecords int) dataset.Dataset {
	return &ReduceDataset{
		dataset:    dataset,
		numRecords: numRecords,
	}
}

func (r *ReduceDataset) Open() (dataset.Conn, error) {
	conn, err := r.dataset.Open()
	if err != nil {
		return nil, err
	}
	return &ReduceDatasetConn{
		dataset:   r,
		conn:      conn,
		recordNum: -1,
		err:       nil,
	}, nil
}

func (r *ReduceDataset) GetFieldNames() []string {
	return r.dataset.GetFieldNames()
}

func (rc *ReduceDatasetConn) Next() bool {
	if rc.conn.Err() != nil {
		return false
	}
	if rc.recordNum < rc.dataset.numRecords {
		rc.recordNum++
		return rc.conn.Next()
	}
	rc.err = io.EOF
	return false
}

func (rc *ReduceDatasetConn) Err() error {
	if rc.err == io.EOF {
		return nil
	}
	return rc.conn.Err()
}

func (rc *ReduceDatasetConn) Read() (dataset.Record, error) {
	record, err := rc.conn.Read()
	return record, err
}

func (rc *ReduceDatasetConn) Close() error {
	return rc.conn.Close()
}
