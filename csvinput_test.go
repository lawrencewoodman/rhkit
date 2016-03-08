package main

import (
	"encoding/csv"
	"errors"
	"github.com/lawrencewoodman/dlit_go"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestNewCsvInput(t *testing.T) {
	cases := []struct {
		filename   string
		fieldNames []string
		wantErr    error
	}{
		{"missing.csv", []string{},
			&os.PathError{"open", "missing.csv",
				errors.New("no such file or directory")}},
		{filepath.Join("fixtures", "bank.csv"),
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome", "y"}, nil},
	}
	for _, c := range cases {
		_, err := NewCsvInput(c.fieldNames, c.filename, ';', false)
		if !errorMatch(c.wantErr, err) {
			t.Errorf("NewCsvInput(filename: %q) err: %q, wantErr: %q",
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
		{filepath.Join("fixtures", "bank.csv"), false,
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome", "y"},
			17, 10,
			map[string]*dlit.Literal{
				"age":       mustNewLit(32),
				"job":       mustNewLit("entrepreneur"),
				"marital":   mustNewLit("married"),
				"education": mustNewLit("secondary"),
				"default":   mustNewLit("no"),
				"balance":   mustNewLit(2),
				"housing":   mustNewLit("yes"),
				"loan":      mustNewLit("yes"),
				"contact":   mustNewLit("unknown"),
				"day":       mustNewLit(5),
				"month":     mustNewLit("may"),
				"duration":  mustNewLit(76),
				"campaign":  mustNewLit(1),
				"pdays":     mustNewLit(-1),
				"previous":  mustNewLit(0),
				"poutcome":  mustNewLit("unknown"),
				"y":         mustNewLit("no")}},
		{filepath.Join("fixtures", "bank.csv"), true,
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome", "y"},
			17, 9,
			map[string]*dlit.Literal{
				"age":       mustNewLit(74),
				"job":       mustNewLit("blue-collar"),
				"marital":   mustNewLit("married"),
				"education": mustNewLit("unknown"),
				"default":   mustNewLit("no"),
				"balance":   mustNewLit(1506),
				"housing":   mustNewLit("yes"),
				"loan":      mustNewLit("no"),
				"contact":   mustNewLit("unknown"),
				"day":       mustNewLit(5),
				"month":     mustNewLit("may"),
				"duration":  mustNewLit(92),
				"campaign":  mustNewLit(1),
				"pdays":     mustNewLit(-1),
				"previous":  mustNewLit(0),
				"poutcome":  mustNewLit("unknown"),
				"y":         mustNewLit("no")}},
	}
	for _, c := range cases {
		i, err := NewCsvInput(c.fieldNames, c.filename, ';', c.skipFirstLine)
		if err != nil {
			t.Errorf("Read() - NewCsvInput() - filename: %q err: %q", c.filename, err)
		}
		gotNumRows := 0
		for {
			record, err := i.Read()
			if err == io.EOF {
				break
			} else if err != nil {
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
		{filepath.Join("fixtures", "invalid_numfields_at_102.csv"), ',',
			[]string{"band", "score", "team", "points", "rating"}, 101,
			&csv.ParseError{102, 0, errors.New("wrong number of fields in line")}},
		{filepath.Join("fixtures", "bank.csv"), ';',
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome"}, -1,
			errors.New("wrong number of field names for input")},
	}
	for _, c := range cases {
		i, err := NewCsvInput(c.fieldNames, c.filename, c.separator, false)
		if err != nil {
			t.Errorf("Read() - NewCsvInput() - filename: %q err: %q", c.filename, err)
		}
		row := 0
		raisedCorrectErr := false
		for {
			_, err := i.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				if row == c.errRow || c.errRow == -1 {
					if err.Error() != c.wantErr.Error() {
						t.Errorf("Read() - filename: %q err: %q, wantErr: %q",
							c.filename, err, c.wantErr)
					} else {
						raisedCorrectErr = true
					}
				} else {
					t.Errorf("Read() - filename: %q err: %q", c.filename, err)
				}
			}
			row++
		}
		if !raisedCorrectErr {
			t.Errorf("Read() - filename: %q failed to raise error: %q",
				c.filename, c.wantErr)
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
		{filepath.Join("fixtures", "bank.csv"),
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome", "y"},
			17, 10,
			map[string]*dlit.Literal{
				"age":       mustNewLit(32),
				"job":       mustNewLit("entrepreneur"),
				"marital":   mustNewLit("married"),
				"education": mustNewLit("secondary"),
				"default":   mustNewLit("no"),
				"balance":   mustNewLit(2),
				"housing":   mustNewLit("yes"),
				"loan":      mustNewLit("yes"),
				"contact":   mustNewLit("unknown"),
				"day":       mustNewLit(5),
				"month":     mustNewLit("may"),
				"duration":  mustNewLit(76),
				"campaign":  mustNewLit(1),
				"pdays":     mustNewLit(-1),
				"previous":  mustNewLit(0),
				"poutcome":  mustNewLit("unknown"),
				"y":         mustNewLit("no")}},
	}
	for _, c := range cases {
		input, err := NewCsvInput(c.fieldNames, c.filename, ';', false)
		if err != nil {
			t.Errorf("Read() - NewCsvInput() - filename: %q err: %q", c.filename, err)
		}
		for i := 0; i < 5; i++ {
			gotNumRows := 0
			for {
				record, err := input.Read()
				if err == io.EOF {
					break
				} else if err != nil {
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
			if gotNumRows != c.wantNumRows {
				t.Errorf("Read() - filename: %q gotNumRows:: %d, want: %d",
					c.filename, gotNumRows, c.wantNumRows)
			}
			if err := input.Rewind(); err != nil {
				t.Errorf("Rewind() - filename: %q err: %s", c.filename, err)
			}
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
