package reduceinput

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rulehunter/csvinput"
	"github.com/vlifesystems/rulehunter/input"
	"path/filepath"
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
		cinput := mustNewCsvInput(c.fieldNames, c.filename, ';', false)
		_, err := New(cinput, c.numRecords)
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
			errors.New("wrong number of field names for input")},
		{filepath.Join("..", "fixtures", "bank.csv"), ';',
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome", "y"},
			20, nil},
	}
	for _, c := range cases {
		input := mustNewCsvInput(c.fieldNames, c.filename, c.separator, false)
		records, err := New(input, c.numRecords)
		if err != nil {
			t.Errorf("Read() - New() - filename: %q err: %q", c.filename, err)
		}
		for records.Next() {
			records.Read()
		}
		if c.wantErr == nil {
			if records.Err() != nil {
				t.Errorf("Read() - filename: %q wantErr: %s, got error: %s",
					c.filename, c.wantErr, records.Err())
			}
		} else {
			if records.Err() == nil ||
				records.Err().Error() != c.wantErr.Error() {
				t.Errorf("Read() - filename: %q wantErr: %s, got error: %s",
					c.filename, c.wantErr, records.Err())
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
		input := mustNewCsvInput(c.fieldNames, c.filename, c.separator, false)
		records, err := New(input, c.numRecords)
		if err != nil {
			t.Errorf("New() - filename: %q err: %q", c.filename, err)
		}
		recordNum := -1
		for records.Next() {
			recordNum++
		}
		if records.Next() {
			t.Errorf("records.Next() - Return true, despite having finished")
		}
		if recordNum != c.numRecords {
			t.Errorf("records.Next() - recordNum: %d, numRecords: %d",
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
			errors.New("input has been closed")},
	}
	for _, c := range cases {
		input := mustNewCsvInput(c.fieldNames, c.filename, c.separator, false)
		records, err := New(input, c.numRecords)
		if err != nil {
			t.Errorf("New() - filename: %q err: %q", c.filename, err)
		}
		i := 0
		for records.Next() {
			if i == c.stopRow {
				if err := records.Close(); err != nil {
					t.Errorf("records.Close() - Err: %d", err)
				}
				break
			}
			i++
		}
		if i != c.stopRow {
			t.Errorf("records.Next() - Not stopped at row: %d", c.stopRow)
		}
		if records.Next() {
			t.Errorf("records.Next() - Return true, despite records being closed")
		}
		if records.Err() == nil || records.Err().Error() != c.wantErr.Error() {
			t.Errorf("records.Err() - err: %s, want err: %s", records.Err(), c.wantErr)
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
		input := mustNewCsvInput(c.fieldNames, c.filename, ';', false)
		records, err := New(input, c.numRecords)
		if err != nil {
			t.Errorf("New() - filename: %q err: %q", c.filename, err)
		}
		for i := 0; i < 5; i++ {
			gotNumRows := 0
			for records.Next() {
				record, err := records.Read()
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
				if err := records.Err(); err != nil {
					t.Errorf("Err() - filename: %s err: %s", c.filename, err)
				}
				gotNumRows++
			}
			if gotNumRows != c.wantNumRows {
				t.Errorf("Read() - filename: %q gotNumRows:: %d, want: %d",
					c.filename, gotNumRows, c.wantNumRows)
			}
			if err := records.Rewind(); err != nil {
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
		input := mustNewCsvInput(c.fieldNames, c.filename, c.separator, false)
		records, err := New(input, c.numRecords)
		if err != nil {
			t.Errorf("New() - filename: %q err: %q", c.filename, err)
		}
		for records.Next() {
			records.Read()
		}
		err = records.Rewind()
		if err.Error() != c.wantErr.Error() {
			t.Errorf("Rewind() - err: %s, wantErr: %s", err, c.wantErr)
		}
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

func mustNewCsvInput(
	fieldNames []string,
	filename string,
	separator rune,
	skipFirstLine bool,
) input.Input {
	input, err := csvinput.New(fieldNames, filename, separator, skipFirstLine)
	if err != nil {
		panic(fmt.Sprintf("Can't create new csvinput for filename: %s", filename))
	}
	return input
}
