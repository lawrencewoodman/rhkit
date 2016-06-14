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
	"syscall"
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
	filename := "missing.csv"
	fieldNames := []string{"age", "occupation"}
	numRecords := 10
	wantErr := &os.PathError{"open", "missing.csv", syscall.ENOENT}
	ds := mustNewCsvDataset(fieldNames, filename, ';', false)
	rds := New(ds, numRecords)
	_, err := rds.Open()
	if err := checkPathErrorMatch(err, wantErr); err != nil {
		t.Errorf("Open() - filename: %s - problem with error: %s",
			filename, err)
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
		filename       string
		separator      rune
		hasHeader      bool
		fieldNames     []string
		wantNumRecords int
	}{
		{filepath.Join("..", "fixtures", "bank.csv"), ';', true,
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome", "y"}, 4},
		{filepath.Join("..", "fixtures", "invalid_numfields_at_102.csv"),
			',', false,
			[]string{"band", "score", "team", "points", "rating"}, 50},
	}
	for _, c := range cases {
		ds := mustNewCsvDataset(c.fieldNames, c.filename, c.separator, c.hasHeader)
		rds := New(ds, c.wantNumRecords)
		conn, err := rds.Open()
		if err != nil {
			t.Errorf("Open() - filename: %s, err: %s", c.filename, err)
		}
		numRecords := 0
		for conn.Next() {
			numRecords++
		}
		if conn.Next() {
			t.Errorf("conn.Next() - Return true, despite having finished")
		}
		if numRecords != c.wantNumRecords {
			t.Errorf("conn.Next() - filename: %s, wantNumRecords: %d, gotNumRecords: %d",
				c.filename, c.wantNumRecords, numRecords)
		}
	}
}

func TestNext_errors(t *testing.T) {
	cases := []struct {
		filename   string
		separator  rune
		hasHeader  bool
		fieldNames []string
		stopRow    int
		numRecords int
		wantErr    error
	}{
		{filename: filepath.Join("..", "fixtures", "bank.csv"),
			separator: ';',
			hasHeader: true,
			fieldNames: []string{"age", "job", "marital", "education", "default",
				"balance", "housing", "loan", "contact", "day", "month", "duration",
				"campaign", "pdays", "previous", "poutcome", "y"},
			stopRow:    2,
			numRecords: 4,
			wantErr:    errors.New("connection has been closed")},
	}
	for _, c := range cases {
		ds := mustNewCsvDataset(c.fieldNames, c.filename, c.separator, c.hasHeader)
		rds := New(ds, c.numRecords)
		conn, err := rds.Open()
		if err != nil {
			t.Errorf("Open() - filename: %s, err: %s", c.filename, err)
		}
		recordNum := 0
		for conn.Next() {
			if recordNum == c.stopRow {
				if err := conn.Close(); err != nil {
					t.Errorf("conn.Close() - Err: %d", err)
				}
				break
			}
			recordNum++
		}
		if recordNum != c.stopRow {
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

		i1Record := i1.Read()
		i2Record := i2.Read()
		if !matchRecords(i1Record, i2Record) {
			return errors.New("Datasets don't match")
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

func checkPathErrorMatch(
	checkErr error,
	wantErr *os.PathError,
) error {
	perr, ok := checkErr.(*os.PathError)
	if !ok {
		return errors.New("error isn't a os.PathError")
	}
	if perr.Op != wantErr.Op {
		return fmt.Errorf("wanted perr.Op: %s, got: %s", perr.Op, wantErr.Op)
	}
	if filepath.Clean(perr.Path) != filepath.Clean(wantErr.Path) {
		return fmt.Errorf("wanted perr.Path: %s, got: %s", perr.Path, wantErr.Path)
	}
	if perr.Err != wantErr.Err {
		return fmt.Errorf("wanted perr.Err: %s, got: %s", perr.Err, wantErr.Err)
	}
	return nil
}
