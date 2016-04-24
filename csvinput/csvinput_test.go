package csvinput

import (
	"encoding/csv"
	"errors"
	"github.com/lawrencewoodman/dlit_go"
	"os"
	"path/filepath"
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
		{"missing.csv", []string{},
			&os.PathError{"open", "missing.csv",
				errors.New("no such file or directory")}},
	}
	for _, c := range cases {
		_, err := New(c.fieldNames, c.filename, ';', false)
		if err.Error() != c.wantErr.Error() {
			t.Errorf("New(filename: %q) err: %q, wantErr: %q",
				c.filename, err, c.wantErr)
		}
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
		records, err := New(c.fieldNames, c.filename, ';', c.skipFirstLine)
		if err != nil {
			t.Errorf("Read() - New() - filename: %q err: %q", c.filename, err)
		}
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
			gotNumRows++
		}
		if err := records.Err(); err != nil {
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
			errors.New("wrong number of field names for input")},
	}
	for _, c := range cases {
		records, err := New(c.fieldNames, c.filename, c.separator, false)
		if err != nil {
			t.Errorf("Read() - New() - filename: %q err: %q", c.filename, err)
		}
		row := 0
		for records.Next() {
			_, err := records.Read()
			if row == c.errRow {
				if err == nil {
					t.Errorf("Read() - filename: %q Failed to raise error", c.filename)
					return
				}
			}
			row++
		}
		if records.Err().Error() != c.wantErr.Error() {
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
			errors.New("wrong number of field names for input")},
	}
	for _, c := range cases {
		records, err := New(c.fieldNames, c.filename, c.separator, false)
		if err != nil {
			t.Errorf("Read() - New() - filename: %q err: %q", c.filename, err)
		}
		_, err = records.Read()
		if err.Error() != c.wantErr.Error() {
			t.Errorf("Read() - filename: %q got error: %s, want error: %s",
				c.filename, err, c.wantErr)
			return
		}
		if records.Err().Error() != c.wantErr.Error() {
			t.Errorf("Read() - filename: %q got error: %s, want error: %s",
				c.filename, records.Err().Error(), c.wantErr)
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
			errors.New("wrong number of field names for input")},
		{filepath.Join("..", "fixtures", "bank.csv"), ';',
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome", "y"}, nil},
	}
	for _, c := range cases {
		records, err := New(c.fieldNames, c.filename, c.separator, false)
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
	}{
		{filepath.Join("..", "fixtures", "bank.csv"), ';',
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome"}},
		{filepath.Join("..", "fixtures", "invalid_numfields_at_102.csv"), ',',
			[]string{"band", "score", "team", "points", "rating"}},
	}
	for _, c := range cases {
		records, err := New(c.fieldNames, c.filename, c.separator, false)
		if err != nil {
			t.Errorf("New() - filename: %q err: %q", c.filename, err)
		}
		for records.Next() {
		}
		if records.Next() {
			t.Errorf("records.Next() - Return true, despite having finished")
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
			errors.New("input has been closed")},
	}
	for _, c := range cases {
		records, err := New(c.fieldNames, c.filename, c.separator, false)
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
		input, err := New(c.fieldNames, c.filename, ';', false)
		if err != nil {
			t.Errorf("Read() - New() - filename: %q err: %q", c.filename, err)
		}
		for i := 0; i < 5; i++ {
			gotNumRows := 0
			for input.Next() {
				record, err := input.Read()
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
				if err := input.Err(); err != nil {
					t.Errorf("Err() - filename: %s err: %s", c.filename, err)
				}
				gotNumRows++
			}
			if gotNumRows != c.wantNumRows {
				t.Errorf("Read() - filename: %q gotNumRows:: %d, want: %d",
					c.filename, gotNumRows, c.wantNumRows)
			}
			if err := input.Rewind(); err != nil {
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
			errors.New("wrong number of field names for input")},
	}
	for _, c := range cases {
		input, err := New(c.fieldNames, c.filename, ';', false)
		if err != nil {
			t.Errorf("New() - filename: %q err: %q", c.filename, err)
		}
		for input.Next() {
			input.Read()
		}
		err = input.Rewind()
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
