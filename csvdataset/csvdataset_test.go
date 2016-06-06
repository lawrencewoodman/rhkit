package csvdataset

import (
	"encoding/csv"
	"errors"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/dataset"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	cases := []struct {
		filename   string
		fieldNames []string
	}{
		{filepath.Join("..", "fixtures", "bank.csv"),
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome", "y"}},
	}
	for _, c := range cases {
		_, err := New(c.fieldNames, c.filename, ';', false)
		if err != nil {
			t.Errorf("New(filename: %q) err: %q", c.filename, err)
		}
	}
}

func TestNew_errors(t *testing.T) {
	cases := []struct {
		filename   string
		fieldNames []string
		wantErr    error
	}{
		{"missing.csv",
			[]string{"age", "occupation"},
			&os.PathError{"open", "missing.csv",
				errors.New("no such file or directory")}},
		{filepath.Join("..", "fixtures", "bank.csv"),
			[]string{"age"},
			errors.New("Must specify at least two field names")},
		{filepath.Join("..", "fixtures", "bank.csv"),
			[]string{"age", "job", "marital", "education", "_default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome", "y"},
			errors.New("Invalid field name: _default")},
	}
	for _, c := range cases {
		_, err := New(c.fieldNames, c.filename, ';', false)
		if err.Error() != c.wantErr.Error() {
			t.Errorf("New(filename: %q) err: %q, wantErr: %q",
				c.filename, err, c.wantErr)
		}
	}
}

func TestClone(t *testing.T) {
	filename := filepath.Join("..", "fixtures", "bank.csv")
	fieldNames := []string{
		"age", "job", "marital", "education", "default", "balance",
		"housing", "loan", "contact", "day", "month", "duration", "campaign",
		"pdays", "previous", "poutcome", "y",
	}
	dataset, err := New(fieldNames, filename, ';', true)
	cDataset, err := dataset.Clone()
	if err != nil {
		t.Errorf("Clone() err: %q", err)
	}
	if err := checkDatasetsEqual(cDataset, dataset); err != nil {
		t.Errorf("Clone() datasets are not equal: %s", err)
	}
}

func TestRead(t *testing.T) {
	cases := []struct {
		filename        string
		skipFirstLine   bool
		fieldNames      []string
		wantNumColumns  int
		wantNumRows     int
		wantThirdRecord map[string]*dlit.Literal
	}{
		{filepath.Join("..", "fixtures", "bank.csv"), false,
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome", "y"},
			17, 10,
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
		{filepath.Join("..", "fixtures", "bank.csv"), true,
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome", "y"},
			17, 9,
			map[string]*dlit.Literal{
				"age":       dlit.MustNew(74),
				"job":       dlit.MustNew("blue-collar"),
				"marital":   dlit.MustNew("married"),
				"education": dlit.MustNew("unknown"),
				"default":   dlit.MustNew("no"),
				"balance":   dlit.MustNew(1506),
				"housing":   dlit.MustNew("yes"),
				"loan":      dlit.MustNew("no"),
				"contact":   dlit.MustNew("unknown"),
				"day":       dlit.MustNew(5),
				"month":     dlit.MustNew("may"),
				"duration":  dlit.MustNew(92),
				"campaign":  dlit.MustNew(1),
				"pdays":     dlit.MustNew(-1),
				"previous":  dlit.MustNew(0),
				"poutcome":  dlit.MustNew("unknown"),
				"y":         dlit.MustNew("no")}},
	}
	for _, c := range cases {
		dataset, err := New(c.fieldNames, c.filename, ';', c.skipFirstLine)
		if err != nil {
			t.Errorf("Read() - New() - filename: %q err: %q", c.filename, err)
		}
		gotNumRows := 0
		for dataset.Next() {
			record, err := dataset.Read()
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
			gotNumRows++
		}
		if err := dataset.Err(); err != nil {
			t.Errorf("Read() - filename: %q err: %s", c.filename, err)
		}
		if gotNumRows != c.wantNumRows {
			t.Errorf("Read() - filename: %q gotNumRows:: %d, want: %d",
				c.filename, gotNumRows, c.wantNumRows)
		}
	}
}

func TestRead_errors(t *testing.T) {
	cases := []struct {
		filename   string
		separator  rune
		fieldNames []string
		errRow     int
		wantErr    error
	}{
		{filepath.Join("..", "fixtures", "invalid_numfields_at_102.csv"), ',',
			[]string{"band", "score", "team", "points", "rating"}, 101,
			&csv.ParseError{102, 0, errors.New("wrong number of fields in line")}},
		{filepath.Join("..", "fixtures", "bank.csv"), ';',
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome"}, 0,
			errors.New("wrong number of field names for dataset")},
	}
	for _, c := range cases {
		dataset, err := New(c.fieldNames, c.filename, c.separator, false)
		if err != nil {
			t.Errorf("Read() - New() - filename: %q err: %q", c.filename, err)
		}
		row := 0
		for dataset.Next() {
			_, err := dataset.Read()
			if row == c.errRow {
				if err == nil {
					t.Errorf("Read() - filename: %q Failed to raise error", c.filename)
					return
				}
			}
			row++
		}
		if dataset.Err().Error() != c.wantErr.Error() {
			t.Errorf("Read() - filename: %q Failed to raise error", c.filename)
		}
	}
}

func TestRead_errors2(t *testing.T) {
	cases := []struct {
		filename   string
		separator  rune
		fieldNames []string
		wantErr    error
	}{
		{filepath.Join("..", "fixtures", "bank.csv"), ';',
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome", "y"},
			errors.New("wrong number of field names for dataset")},
	}
	for _, c := range cases {
		dataset, err := New(c.fieldNames, c.filename, c.separator, false)
		if err != nil {
			t.Errorf("Read() - New() - filename: %q err: %q", c.filename, err)
		}
		_, err = dataset.Read()
		if err.Error() != c.wantErr.Error() {
			t.Errorf("Read() - filename: %q got error: %s, want error: %s",
				c.filename, err, c.wantErr)
			return
		}
		if dataset.Err().Error() != c.wantErr.Error() {
			t.Errorf("Read() - filename: %q got error: %s, want error: %s",
				c.filename, dataset.Err().Error(), c.wantErr)
		}
	}
}

func TestErr(t *testing.T) {
	cases := []struct {
		filename   string
		separator  rune
		fieldNames []string
		wantErr    error
	}{
		{filepath.Join("..", "fixtures", "invalid_numfields_at_102.csv"), ',',
			[]string{"band", "score", "team", "points", "rating"},
			&csv.ParseError{102, 0, errors.New("wrong number of fields in line")}},
		{filepath.Join("..", "fixtures", "bank.csv"), ';',
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome"},
			errors.New("wrong number of field names for dataset")},
		{filepath.Join("..", "fixtures", "bank.csv"), ';',
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome", "y"}, nil},
	}
	for _, c := range cases {
		dataset, err := New(c.fieldNames, c.filename, c.separator, false)
		if err != nil {
			t.Errorf("Read() - New() - filename: %q err: %q", c.filename, err)
		}
		for dataset.Next() {
			dataset.Read()
		}
		if c.wantErr == nil {
			if dataset.Err() != nil {
				t.Errorf("Read() - filename: %q wantErr: %s, got error: %s",
					c.filename, c.wantErr, dataset.Err())
			}
		} else {
			if dataset.Err() == nil ||
				dataset.Err().Error() != c.wantErr.Error() {
				t.Errorf("Read() - filename: %q wantErr: %s, got error: %s",
					c.filename, c.wantErr, dataset.Err())
			}
		}
	}
}

func TestNext(t *testing.T) {
	cases := []struct {
		filename   string
		separator  rune
		fieldNames []string
	}{
		{filepath.Join("..", "fixtures", "bank.csv"), ';',
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome"}},
		{filepath.Join("..", "fixtures", "invalid_numfields_at_102.csv"), ',',
			[]string{"band", "score", "team", "points", "rating"}},
	}
	for _, c := range cases {
		dataset, err := New(c.fieldNames, c.filename, c.separator, false)
		if err != nil {
			t.Errorf("New() - filename: %q err: %q", c.filename, err)
		}
		for dataset.Next() {
		}
		if dataset.Next() {
			t.Errorf("dataset.Next() - Return true, despite having finished")
		}
	}
}

func TestNext_errors(t *testing.T) {
	cases := []struct {
		filename   string
		separator  rune
		fieldNames []string
		stopRow    int
		wantErr    error
	}{
		{filepath.Join("..", "fixtures", "bank.csv"), ';',
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome"}, 2,
			errors.New("dataset has been closed")},
	}
	for _, c := range cases {
		dataset, err := New(c.fieldNames, c.filename, c.separator, false)
		if err != nil {
			t.Errorf("New() - filename: %q err: %q", c.filename, err)
		}
		i := 0
		for dataset.Next() {
			if i == c.stopRow {
				if err := dataset.Close(); err != nil {
					t.Errorf("dataset.Close() - Err: %d", err)
				}
				break
			}
			i++
		}
		if i != c.stopRow {
			t.Errorf("dataset.Next() - Not stopped at row: %d", c.stopRow)
		}
		if dataset.Next() {
			t.Errorf("dataset.Next() - Return true, despite dataset being closed")
		}
		if dataset.Err() == nil || dataset.Err().Error() != c.wantErr.Error() {
			t.Errorf("dataset.Err() - err: %s, want err: %s",
				dataset.Err(), c.wantErr)
		}
	}
}

func TestRewind(t *testing.T) {
	cases := []struct {
		filename        string
		fieldNames      []string
		wantNumColumns  int
		wantNumRows     int
		wantThirdRecord map[string]*dlit.Literal
	}{
		{filepath.Join("..", "fixtures", "bank.csv"),
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome", "y"},
			17, 10,
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
		dataset, err := New(c.fieldNames, c.filename, ';', false)
		if err != nil {
			t.Errorf("Read() - New() - filename: %q err: %q", c.filename, err)
		}
		for i := 0; i < 5; i++ {
			gotNumRows := 0
			for dataset.Next() {
				record, err := dataset.Read()
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
				if err := dataset.Err(); err != nil {
					t.Errorf("Err() - filename: %s err: %s", c.filename, err)
				}
				gotNumRows++
			}
			if gotNumRows != c.wantNumRows {
				t.Errorf("Read() - filename: %q gotNumRows:: %d, want: %d",
					c.filename, gotNumRows, c.wantNumRows)
			}
			if err := dataset.Rewind(); err != nil {
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
		wantErr    error
	}{
		{filepath.Join("..", "fixtures", "invalid_numfields_at_102.csv"), ',',
			[]string{"band", "score", "team", "points", "rating"},
			&csv.ParseError{102, 0, errors.New("wrong number of fields in line")}},
	}
	for _, c := range cases {
		dataset, err := New(c.fieldNames, c.filename, c.separator, false)
		if err != nil {
			t.Errorf("New() - filename: %q err: %q", c.filename, err)
		}
		for dataset.Next() {
			dataset.Read()
		}
		err = dataset.Rewind()
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
	dataset, err := New(fieldNames, filename, ';', false)
	if err != nil {
		t.Errorf("New(%s, %s, ...) err: %q", fieldNames, filename, err)
	}

	got := dataset.GetFieldNames()
	if !reflect.DeepEqual(got, fieldNames) {
		t.Errorf("GetFieldNames() - got: %q want: %q", got, fieldNames)
	}
}

/*************************
 *   Helper functions
 *************************/

func checkDatasetsEqual(i1, i2 dataset.Dataset) error {
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

func matchRecords(
	r1 map[string]*dlit.Literal,
	r2 map[string]*dlit.Literal,
) bool {
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
