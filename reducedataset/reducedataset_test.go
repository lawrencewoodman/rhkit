package reducedataset

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/vlifesystems/rulehunter/csvdataset"
	"github.com/vlifesystems/rulehunter/dataset"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestOpen(t *testing.T) {
	cases := []struct {
		filename   string
		fieldNames []string
		numRecords int
	}{
		{filepath.Join("..", "fixtures", "bank.csv"),
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome", "y"},
			10},
	}
	for _, c := range cases {
		ds := mustNewCsvDataset(c.fieldNames, c.filename, ';', false)
		rds := New(ds, c.numRecords)
		if _, err := rds.Open(); err != nil {
			t.Errorf("Open() err: %s", err)
		}
	}
}

func TestOpen_errors(t *testing.T) {
	cases := []struct {
		filename   string
		fieldNames []string
		numRecords int
		wantErr    error
	}{
		{"missing.csv",
			[]string{"age", "occupation"},
			10,
			&os.PathError{"open", "missing.csv",
				errors.New("no such file or directory")}},
	}
	for _, c := range cases {
		ds := mustNewCsvDataset(c.fieldNames, c.filename, ';', false)
		rds := New(ds, c.numRecords)
		_, err := rds.Open()
		if err.Error() != c.wantErr.Error() {
			t.Errorf("Open() - filename: %s, err: %s, wantErr: %s",
				c.filename, err, c.wantErr)
		}
	}
}

func TestGetFieldNames(t *testing.T) {
	filename := filepath.Join("..", "fixtures", "bank.csv")
	fieldNames := []string{
		"age", "job", "marital", "education", "default", "balance",
		"housing", "loan", "contact", "day", "month", "duration", "campaign",
		"pdays", "previous", "poutcome", "y",
	}
	numRecords := 3
	ds := mustNewCsvDataset(fieldNames, filename, ';', false)
	rds := New(ds, numRecords)

	got := rds.GetFieldNames()
	if !reflect.DeepEqual(got, fieldNames) {
		t.Errorf("GetFieldNames() - got: %s, want: %s", got, fieldNames)
	}
}

func TestErr(t *testing.T) {
	cases := []struct {
		filename   string
		separator  rune
		fieldNames []string
		numRecords int
		wantErr    error
	}{
		{filepath.Join("..", "fixtures", "invalid_numfields_at_102.csv"), ',',
			[]string{"band", "score", "team", "points", "rating"},
			105,
			&csv.ParseError{102, 0, errors.New("wrong number of fields in line")}},
		{filepath.Join("..", "fixtures", "bank.csv"), ';',
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome"},
			4,
			errors.New("wrong number of field names for dataset")},
		{filepath.Join("..", "fixtures", "bank.csv"), ';',
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome", "y"},
			20, nil},
	}
	for _, c := range cases {
		ds := mustNewCsvDataset(c.fieldNames, c.filename, c.separator, false)
		rds := New(ds, c.numRecords)
		conn, err := rds.Open()
		if err != nil {
			t.Errorf("Open() - filename: %s, err: %s", c.filename, err)
		}
		for conn.Next() {
			conn.Read()
		}
		if c.wantErr == nil {
			if conn.Err() != nil {
				t.Errorf("Read() - filename: %s, wantErr: %s, got error: %s",
					c.filename, c.wantErr, conn.Err())
			}
		} else {
			if conn.Err() == nil || conn.Err().Error() != c.wantErr.Error() {
				t.Errorf("Read() - filename: %s, wantErr: %s, got error: %s",
					c.filename, c.wantErr, conn.Err())
			}
		}
	}
}

func TestNext(t *testing.T) {
	cases := []struct {
		filename   string
		separator  rune
		fieldNames []string
		numRecords int
	}{
		{filepath.Join("..", "fixtures", "bank.csv"), ';',
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome"}, 4},
		{filepath.Join("..", "fixtures", "invalid_numfields_at_102.csv"), ',',
			[]string{"band", "score", "team", "points", "rating"}, 50},
	}
	for _, c := range cases {
		ds := mustNewCsvDataset(c.fieldNames, c.filename, c.separator, false)
		rds := New(ds, c.numRecords)
		conn, err := rds.Open()
		if err != nil {
			t.Errorf("Open() - filename: %s, err: %s", c.filename, err)
		}
		recordNum := -1
		for conn.Next() {
			recordNum++
		}
		if conn.Next() {
			t.Errorf("conn.Next() - Return true, despite having finished")
		}
		if recordNum != c.numRecords {
			t.Errorf("conn.Next() - recordNum: %d, numRecords: %d",
				recordNum, c.numRecords)
		}
	}
}

func TestNext_errors(t *testing.T) {
	cases := []struct {
		filename   string
		separator  rune
		fieldNames []string
		stopRow    int
		numRecords int
		wantErr    error
	}{
		{filepath.Join("..", "fixtures", "bank.csv"), ';',
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome"}, 2, 4,
			errors.New("connection has been closed")},
	}
	for _, c := range cases {
		ds := mustNewCsvDataset(c.fieldNames, c.filename, c.separator, false)
		rds := New(ds, c.numRecords)
		conn, err := rds.Open()
		if err != nil {
			t.Errorf("Open() - filename: %s, err: %s", c.filename, err)
		}
		i := 0
		for conn.Next() {
			if i == c.stopRow {
				if err := conn.Close(); err != nil {
					t.Errorf("conn.Close() - Err: %d", err)
				}
				break
			}
			i++
		}
		if i != c.stopRow {
			t.Errorf("conn.Next() - Not stopped at row: %d", c.stopRow)
		}
		if conn.Next() {
			t.Errorf("conn.Next() - Return true, despite reducedDataset being closed")
		}
		if conn.Err() == nil || conn.Err().Error() != c.wantErr.Error() {
			t.Errorf("conn.Err() - err: %s, want err: %s", conn.Err(), c.wantErr)
		}
	}
}

/*************************
 *   Helper functions
 *************************/

func checkDatasetsEqual(i1, i2 dataset.Conn) error {
	for {
		i1Next := i1.Next()
		i2Next := i2.Next()
		if i1Next != i2Next {
			return errors.New("Datasets don't finish at same point")
		}
		if !i1Next {
			break
		}

		i1Record, i1Err := i1.Read()
		i2Record, i2Err := i2.Read()
		if i1Err != i2Err {
			return errors.New("Datasets don't error at same point")
		} else if i1Err == nil && i2Err == nil {
			if !matchRecords(i1Record, i2Record) {
				return errors.New("Datasets don't match")
			}
		}
	}
	if i1.Err() != i2.Err() {
		return errors.New("Datasets final error doesn't match")
	}
	return nil
}

func matchRecords(r1 dataset.Record, r2 dataset.Record) bool {
	if len(r1) != len(r2) {
		return false
	}
	for fieldName, value := range r1 {
		if value.String() != r2[fieldName].String() {
			return false
		}
	}
	return true
}

func mustNewCsvDataset(
	fieldNames []string,
	filename string,
	separator rune,
	skipFirstLine bool,
) dataset.Dataset {
	dataset, err := csvdataset.New(fieldNames, filename, separator, skipFirstLine)
	if err != nil {
		panic(fmt.Sprintf("Can't create new csvdataset for filename: %s", filename))
	}
	return dataset
}
