package main

import (
	"encoding/csv"
	"errors"
	"github.com/lawrencewoodman/dlit"
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
		_, err := NewCsvInput(c.fieldNames, c.filename, ';')
		if !errorMatch(c.wantErr, err) {
			t.Errorf("NewCsvInput(filename: %q) err: %q, wantErr: %q",
				c.filename, err, c.wantErr)
		}
	}
}

func TestRead(t *testing.T) {
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
		i, err := NewCsvInput(c.fieldNames, c.filename, ';')
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
		fieldNames []string
		errRow     int
		wantErr    error
	}{
		{filepath.Join("fixtures", "invalid_numfields_at_102.csv"),
			[]string{"age", "job", "marital", "education", "default", "balance",
				"housing", "loan", "contact", "day", "month", "duration", "campaign",
				"pdays", "previous", "poutcome", "y"},
			101,
			&csv.ParseError{102, 0, errors.New("wrong number of fields in line")}},
	}
	for _, c := range cases {
		i, err := NewCsvInput(c.fieldNames, c.filename, ',')
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
				if row == c.errRow {
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
