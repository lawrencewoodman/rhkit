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
	"github.com/lawrencewoodman/ddataset"
)

type ReduceDataset struct {
	dataset    ddataset.Dataset
	numRecords int
}

type ReduceDatasetConn struct {
	dataset   *ReduceDataset
	conn      ddataset.Conn
	recordNum int
	err       error
}

func New(dataset ddataset.Dataset, numRecords int) ddataset.Dataset {
	return &ReduceDataset{
		dataset:    dataset,
		numRecords: numRecords,
	}
}

func (r *ReduceDataset) Open() (ddataset.Conn, error) {
	conn, err := r.dataset.Open()
	if err != nil {
		return nil, err
	}
	return &ReduceDatasetConn{
		dataset:   r,
		conn:      conn,
		recordNum: 0,
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
	return false
}

func (rc *ReduceDatasetConn) Err() error {
	return rc.conn.Err()
}

func (rc *ReduceDatasetConn) Read() ddataset.Record {
	return rc.conn.Read()
}

func (rc *ReduceDatasetConn) Close() error {
	return rc.conn.Close()
}
