package reducedataset

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/csvdataset"
	"github.com/vlifesystems/rulehunter/dataset"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
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
		dataset := mustNewCsvDataset(c.fieldNames, c.filename, ';', false)
		_, err := New(dataset, c.numRecords)
		if err != nil {
			t.Errorf("New(filename: %q) err: %q", c.filename, err)
		}
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
		dataset := mustNewCsvDataset(c.fieldNames, c.filename, c.separator, false)
		reducedDataset, err := New(dataset, c.numRecords)
		if err != nil {
			t.Errorf("Read() - New() - filename: %q err: %q", c.filename, err)
		}
		for reducedDataset.Next() {
			reducedDataset.Read()
		}
		if c.wantErr == nil {
			if reducedDataset.Err() != nil {
				t.Errorf("Read() - filename: %q wantErr: %s, got error: %s",
					c.filename, c.wantErr, reducedDataset.Err())
			}
		} else {
			if reducedDataset.Err() == nil ||
				reducedDataset.Err().Error() != c.wantErr.Error() {
				t.Errorf("Read() - filename: %q wantErr: %s, got error: %s",
					c.filename, c.wantErr, reducedDataset.Err())
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
		dataset := mustNewCsvDataset(c.fieldNames, c.filename, c.separator, false)
		reducedDataset, err := New(dataset, c.numRecords)
		if err != nil {
			t.Errorf("New() - filename: %q err: %q", c.filename, err)
		}
		recordNum := -1
		for reducedDataset.Next() {
			recordNum++
		}
		if reducedDataset.Next() {
			t.Errorf("reducedDataset.Next() - Return true, despite having finished")
		}
		if recordNum != c.numRecords {
			t.Errorf("reducedDataset.Next() - recordNum: %d, numRecords: %d",
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
			errors.New("dataset has been closed")},
	}
	for _, c := range cases {
		dataset := mustNewCsvDataset(c.fieldNames, c.filename, c.separator, false)
		reducedDataset, err := New(dataset, c.numRecords)
		if err != nil {
			t.Errorf("New() - filename: %q err: %q", c.filename, err)
		}
		i := 0
		for reducedDataset.Next() {
			if i == c.stopRow {
				if err := reducedDataset.Close(); err != nil {
					t.Errorf("reducedDataset.Close() - Err: %d", err)
				}
				break
			}
			i++
		}
		if i != c.stopRow {
			t.Errorf("reducedDataset.Next() - Not stopped at row: %d", c.stopRow)
		}
		if reducedDataset.Next() {
			t.Errorf("reducedDataset.Next() - Return true, despite reducedDataset being closed")
		}
		if reducedDataset.Err() == nil ||
			reducedDataset.Err().Error() != c.wantErr.Error() {
			t.Errorf("reducedDataset.Err() - err: %s, want err: %s",
				reducedDataset.Err(), c.wantErr)
		}
	}
}

func TestRewind(t *testing.T) {
	cases := []struct {
		filename        string
		fieldNames      []string
		numRecords      int
		wantNumColumns  int
		wantNumRows     int
		wantThirdRecord map[string]*dlit.Literal
	}{
		{filepath.Join("..", "fixtures", "bank.csv"),
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome", "y"},
			20, 17, 10,
			map[string]*dlit.Literal{
				"age":       dlit.MustNew(32),
				"job":       dlit.MustNew("entrepreneur"),
				"marital":   dlit.MustNew("married"),
				"education": dlit.MustNew("secondary"),
				"default":   dlit.MustNew("no"),
				"balance":   dlit.MustNew(2),
				"housing":   dlit.MustNew("yes"),
				"loan":      dlit.MustNew("yes"),
				"contact":   dlit.MustNew("unknown"),
				"day":       dlit.MustNew(5),
				"month":     dlit.MustNew("may"),
				"duration":  dlit.MustNew(76),
				"campaign":  dlit.MustNew(1),
				"pdays":     dlit.MustNew(-1),
				"previous":  dlit.MustNew(0),
				"poutcome":  dlit.MustNew("unknown"),
				"y":         dlit.MustNew("no")}},
	}
	for _, c := range cases {
		dataset := mustNewCsvDataset(c.fieldNames, c.filename, ';', false)
		reducedDataset, err := New(dataset, c.numRecords)
		if err != nil {
			t.Errorf("New() - filename: %q err: %q", c.filename, err)
		}
		for i := 0; i < 5; i++ {
			gotNumRows := 0
			for reducedDataset.Next() {
				record, err := reducedDataset.Read()
				if err != nil {
					t.Errorf("Read() - filename: %q err: %q", c.filename, err)
				}

				gotNumColumns := len(record)
				if gotNumColumns != c.wantNumColumns {
					t.Errorf("Read() - filename: %q gotNumColumns: %d, want: %d",
						c.filename, gotNumColumns, c.wantNumColumns)
				}
				if gotNumRows == 2 && !matchRecords(record, c.wantThirdRecord) {
					t.Errorf("Read() - filename: %q got: %q, want: %q",
						c.filename, record, c.wantThirdRecord)
				}
				if err := reducedDataset.Err(); err != nil {
					t.Errorf("Err() - filename: %s err: %s", c.filename, err)
				}
				gotNumRows++
			}
			if gotNumRows != c.wantNumRows {
				t.Errorf("Read() - filename: %q gotNumRows:: %d, want: %d",
					c.filename, gotNumRows, c.wantNumRows)
			}
			if err := reducedDataset.Rewind(); err != nil {
				t.Errorf("Rewind() - filename: %s err: %s", c.filename, err)
			}
		}
	}
}

func TestRewind_errors(t *testing.T) {
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
	}
	for _, c := range cases {
		dataset := mustNewCsvDataset(c.fieldNames, c.filename, c.separator, false)
		reducedDataset, err := New(dataset, c.numRecords)
		if err != nil {
			t.Errorf("New() - filename: %q err: %q", c.filename, err)
		}
		for reducedDataset.Next() {
			reducedDataset.Read()
		}
		err = reducedDataset.Rewind()
		if err.Error() != c.wantErr.Error() {
			t.Errorf("Rewind() - err: %s, wantErr: %s", err, c.wantErr)
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
	dataset := mustNewCsvDataset(fieldNames, filename, ';', false)
	reducedDataset, err := New(dataset, numRecords)
	if err != nil {
		t.Errorf("New() - filename: %s err: %s", filename, err)
	}

	got := reducedDataset.GetFieldNames()
	if !reflect.DeepEqual(got, fieldNames) {
		t.Errorf("GetFieldNames() - got: %q want: %q", got, fieldNames)
	}
}

/*************************
 *   Helper functions
 *************************/

func matchRecords(r1 map[string]*dlit.Literal,
	r2 map[string]*dlit.Literal) bool {
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
