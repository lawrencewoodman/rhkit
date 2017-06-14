package description

import (
	"errors"
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestDescriptionNew(t *testing.T) {
	got := New()
	if len(got.Fields) != 0 {
		t.Errorf("New got len(got.Fields): %d, want: 0", len(got.Fields))
	}
}

func TestDescriptionWriteLoadJSON(t *testing.T) {
	description := &Description{
		map[string]*Field{
			"band": &Field{fieldtype.String, nil, nil, 0,
				map[string]Value{
					"a": Value{dlit.MustNew("a"), 2},
					"b": Value{dlit.MustNew("b"), 3},
					"c": Value{dlit.MustNew("c"), 70},
					"f": Value{dlit.MustNew("f"), 22},
					"9": Value{dlit.MustNew("9"), 1},
				},
				31,
			},
			"inputA": &Field{
				fieldtype.Number,
				dlit.MustNew(7),
				dlit.MustNew(15.1),
				1,
				map[string]Value{
					"7":    Value{dlit.MustNew(7), 7},
					"7.3":  Value{dlit.MustNew(7.3), 7},
					"9":    Value{dlit.MustNew(9), 7},
					"14":   Value{dlit.MustNew(14), 7},
					"15.1": Value{dlit.MustNew(15.1), 7},
				},
				5,
			},
			"inputB": &Field{
				fieldtype.Number,
				dlit.MustNew(2),
				dlit.MustNew(5),
				4,
				map[string]Value{
					"2.6":    Value{dlit.MustNew(2.6), 7},
					"2.8789": Value{dlit.MustNew(2.8789), 1},
					"3":      Value{dlit.MustNew(3), 7},
					"5":      Value{dlit.MustNew(5), 7},
					"2":      Value{dlit.MustNew(2), 7},
					"2.8":    Value{dlit.MustNew(2.8), 6},
				},
				6,
			},
			"version": &Field{fieldtype.String, nil, nil, 0,
				map[string]Value{
					"9.9":   Value{dlit.MustNew("9.9"), 7},
					"9.97":  Value{dlit.MustNew("9.97"), 7},
					"10":    Value{dlit.MustNew("10"), 7},
					"10.94": Value{dlit.MustNew("10.94"), 7},
					"9.9a":  Value{dlit.MustNew("9.9a"), 6},
					"9.9b":  Value{dlit.MustNew("9.9b"), 1},
				},
				6,
			},
			"flow": &Field{
				fieldtype.Number,
				dlit.MustNew(21),
				dlit.MustNew(87),
				0,
				map[string]Value{}, -1},
			"score": &Field{
				fieldtype.Number,
				dlit.MustNew(1),
				dlit.MustNew(5),
				0,
				map[string]Value{
					"1": Value{dlit.MustNew(1), 6},
					"2": Value{dlit.MustNew(2), 7},
					"3": Value{dlit.MustNew(3), 6},
					"4": Value{dlit.MustNew(4), 8},
					"5": Value{dlit.MustNew(5), 8},
				}, 5,
			},
			"method": &Field{fieldtype.Ignore, nil, nil, 0,
				map[string]Value{}, -1},
		},
	}
	tempDir, err := ioutil.TempDir("", "rulehunter_test")
	if err != nil {
		t.Fatalf("TempDir() err: %s", err)
	}
	defer os.RemoveAll(tempDir)
	filename := filepath.Join(tempDir, "fd.json")
	if err := description.WriteJSON(filename); err != nil {
		t.Fatalf("WriteJSON: %s", err)
	}
	got, err := LoadJSON(filename)
	if err != nil {
		t.Fatalf("LoadJSON: %s", err)
	}
	if err := got.CheckEqual(description); err != nil {
		t.Errorf("LoadJSON got not expected: %s", err)
	}
}

func TestDescriptionLoadJSON_errors(t *testing.T) {
	cases := []struct {
		filename string
		wantErr  error
	}{
		{filename: filepath.Join("fixtures", "nonexistant.json"),
			wantErr: &os.PathError{
				"open",
				filepath.Join("fixtures", "nonexistant.json"),
				syscall.ENOENT,
			},
		},
		{filename: filepath.Join("fixtures", "broken.json"),
			wantErr: errors.New("unexpected EOF"),
		},
	}
	for i, c := range cases {
		_, err := LoadJSON(c.filename)
		checkErrorMatch(
			t,
			fmt.Sprintf("(%d) LoadJSON:", i),
			err,
			c.wantErr,
		)
	}
}

func TestDescriptionCheckEqual(t *testing.T) {
	descriptions := []*Description{
		&Description{
			map[string]*Field{
				"band": &Field{fieldtype.String, nil, nil, 0,
					map[string]Value{
						"a": Value{dlit.MustNew("a"), 2},
						"b": Value{dlit.MustNew("b"), 3},
						"c": Value{dlit.MustNew("c"), 70},
						"f": Value{dlit.MustNew("f"), 22},
						"9": Value{dlit.MustNew("9"), 1},
					},
					31,
				},
				"inputA": &Field{
					fieldtype.Number,
					dlit.MustNew(7),
					dlit.MustNew(15.1),
					1,
					map[string]Value{
						"7":    Value{dlit.MustNew(7), 7},
						"7.3":  Value{dlit.MustNew(7.3), 7},
						"9":    Value{dlit.MustNew(9), 7},
						"14":   Value{dlit.MustNew(14), 7},
						"15.1": Value{dlit.MustNew(15.1), 7},
					},
					5,
				},
				"inputB": &Field{
					fieldtype.Number,
					dlit.MustNew(2),
					dlit.MustNew(5),
					4,
					map[string]Value{
						"2.6":    Value{dlit.MustNew(2.6), 7},
						"2.8789": Value{dlit.MustNew(2.8789), 1},
						"3":      Value{dlit.MustNew(3), 7},
						"5":      Value{dlit.MustNew(5), 7},
						"2":      Value{dlit.MustNew(2), 7},
						"2.8":    Value{dlit.MustNew(2.8), 6},
					},
					6,
				},
			},
		},
		&Description{
			map[string]*Field{
				"band": &Field{fieldtype.String, nil, nil, 0,
					map[string]Value{
						"a": Value{dlit.MustNew("a"), 2},
						"b": Value{dlit.MustNew("b"), 3},
						"c": Value{dlit.MustNew("c"), 70},
						"f": Value{dlit.MustNew("f"), 22},
						"9": Value{dlit.MustNew("9"), 1},
					},
					31,
				},
				"inputB": &Field{
					fieldtype.Number,
					dlit.MustNew(2),
					dlit.MustNew(5),
					4,
					map[string]Value{
						"2.6":    Value{dlit.MustNew(2.6), 7},
						"2.8789": Value{dlit.MustNew(2.8789), 1},
						"3":      Value{dlit.MustNew(3), 7},
						"5":      Value{dlit.MustNew(5), 7},
						"2":      Value{dlit.MustNew(2), 7},
						"2.8":    Value{dlit.MustNew(2.8), 6},
					},
					6,
				},
			},
		},
		&Description{
			map[string]*Field{
				"strata": &Field{fieldtype.String, nil, nil, 0,
					map[string]Value{
						"a": Value{dlit.MustNew("a"), 2},
						"b": Value{dlit.MustNew("b"), 3},
						"c": Value{dlit.MustNew("c"), 70},
						"f": Value{dlit.MustNew("f"), 22},
						"9": Value{dlit.MustNew("9"), 1},
					},
					31,
				},
				"inputA": &Field{
					fieldtype.Number,
					dlit.MustNew(7),
					dlit.MustNew(15.1),
					1,
					map[string]Value{
						"7":    Value{dlit.MustNew(7), 7},
						"7.3":  Value{dlit.MustNew(7.3), 7},
						"9":    Value{dlit.MustNew(9), 7},
						"14":   Value{dlit.MustNew(14), 7},
						"15.1": Value{dlit.MustNew(15.1), 7},
					},
					5,
				},
				"inputB": &Field{
					fieldtype.Number,
					dlit.MustNew(2),
					dlit.MustNew(5),
					4,
					map[string]Value{
						"2.6":    Value{dlit.MustNew(2.6), 7},
						"2.8789": Value{dlit.MustNew(2.8789), 1},
						"3":      Value{dlit.MustNew(3), 7},
						"5":      Value{dlit.MustNew(5), 7},
						"2":      Value{dlit.MustNew(2), 7},
						"2.8":    Value{dlit.MustNew(2.8), 6},
					},
					6,
				},
			},
		},
		&Description{
			map[string]*Field{
				"band": &Field{fieldtype.String, nil, nil, 0,
					map[string]Value{
						"a": Value{dlit.MustNew("a"), 2},
						"b": Value{dlit.MustNew("b"), 3},
						"c": Value{dlit.MustNew("c"), 70},
						"f": Value{dlit.MustNew("f"), 22},
						"9": Value{dlit.MustNew("9"), 1},
					},
					31,
				},
				"inputA": &Field{
					fieldtype.Number,
					dlit.MustNew(6),
					dlit.MustNew(15.1),
					1,
					map[string]Value{
						"7":    Value{dlit.MustNew(7), 7},
						"7.3":  Value{dlit.MustNew(7.3), 7},
						"9":    Value{dlit.MustNew(9), 7},
						"14":   Value{dlit.MustNew(14), 7},
						"15.1": Value{dlit.MustNew(15.1), 7},
					},
					5,
				},
				"inputB": &Field{
					fieldtype.Number,
					dlit.MustNew(2),
					dlit.MustNew(5),
					4,
					map[string]Value{
						"2.6":    Value{dlit.MustNew(2.6), 7},
						"2.8789": Value{dlit.MustNew(2.8789), 1},
						"3":      Value{dlit.MustNew(3), 7},
						"5":      Value{dlit.MustNew(5), 7},
						"2":      Value{dlit.MustNew(2), 7},
						"2.8":    Value{dlit.MustNew(2.8), 6},
					},
					6,
				},
			},
		},
	}
	cases := []struct {
		ndxA int
		ndxB int
		want error
	}{
		{0, 0, nil},
		{0, 1, errors.New("number of Fields doesn't match: 3 != 2")},
		{0, 2, errors.New("missing field: band")},
		{0, 3, errors.New("description for field: inputA, Min not equal: 7 != 6")},
	}
	for i, c := range cases {
		got := descriptions[c.ndxA].CheckEqual(descriptions[c.ndxB])
		checkErrorMatch(t, fmt.Sprintf("(%d) CheckEqual: ", i), got, c.want)
	}
}

func TestFieldCheckEqual(t *testing.T) {
	fields := []*Field{
		&Field{fieldtype.String, nil, nil, 0,
			map[string]Value{
				"a": Value{dlit.MustNew("a"), 2},
				"b": Value{dlit.MustNew("b"), 3},
				"c": Value{dlit.MustNew("c"), 70},
				"f": Value{dlit.MustNew("f"), 22},
				"9": Value{dlit.MustNew("9"), 1},
			},
			31,
		},
		&Field{fieldtype.String, nil, nil, 0,
			map[string]Value{
				"a": Value{dlit.MustNew("a"), 2},
				"b": Value{dlit.MustNew("b"), 3},
				"c": Value{dlit.MustNew("c"), 70},
				"f": Value{dlit.MustNew("f"), 22},
				"9": Value{dlit.MustNew("9"), 1},
			},
			18,
		},
		&Field{
			fieldtype.Number,
			dlit.MustNew(2),
			dlit.MustNew(5),
			4,
			map[string]Value{
				"2.6":    Value{dlit.MustNew(2.6), 7},
				"2.8789": Value{dlit.MustNew(2.8789), 1},
				"3":      Value{dlit.MustNew(3), 7},
				"5":      Value{dlit.MustNew(5), 7},
				"2":      Value{dlit.MustNew(2), 7},
				"2.8":    Value{dlit.MustNew(2.8), 6},
			},
			6,
		},
		&Field{
			fieldtype.Number,
			dlit.MustNew(7),
			dlit.MustNew(5),
			4,
			map[string]Value{
				"2.6":    Value{dlit.MustNew(2.6), 7},
				"2.8789": Value{dlit.MustNew(2.8789), 1},
				"3":      Value{dlit.MustNew(3), 7},
				"5":      Value{dlit.MustNew(5), 7},
				"2":      Value{dlit.MustNew(2), 7},
				"2.8":    Value{dlit.MustNew(2.8), 6},
			},
			6,
		},
		&Field{
			fieldtype.Number,
			dlit.MustNew(2),
			dlit.MustNew(4),
			4,
			map[string]Value{
				"2.6":    Value{dlit.MustNew(2.6), 7},
				"2.8789": Value{dlit.MustNew(2.8789), 1},
				"3":      Value{dlit.MustNew(3), 7},
				"5":      Value{dlit.MustNew(5), 7},
				"2":      Value{dlit.MustNew(2), 7},
				"2.8":    Value{dlit.MustNew(2.8), 6},
			},
			6,
		},
		&Field{
			fieldtype.Number,
			dlit.MustNew(2),
			dlit.MustNew(5),
			2,
			map[string]Value{
				"2.6":    Value{dlit.MustNew(2.6), 7},
				"2.8789": Value{dlit.MustNew(2.8789), 1},
				"3":      Value{dlit.MustNew(3), 7},
				"5":      Value{dlit.MustNew(5), 7},
				"2":      Value{dlit.MustNew(2), 7},
				"2.8":    Value{dlit.MustNew(2.8), 6},
			},
			6,
		},
		&Field{
			fieldtype.Number,
			dlit.MustNew(7),
			dlit.MustNew(5),
			4,
			map[string]Value{
				"2.6":    Value{dlit.MustNew(2.6), 7},
				"2.8789": Value{dlit.MustNew(2.8789), 1},
				"3.3":    Value{dlit.MustNew(3.3), 6},
				"5":      Value{dlit.MustNew(5), 7},
				"2":      Value{dlit.MustNew(2), 7},
				"2.8":    Value{dlit.MustNew(2.8), 6},
			},
			6,
		},
		&Field{
			fieldtype.Number,
			dlit.MustNew(7),
			dlit.MustNew(5),
			4,
			map[string]Value{
				"2.6":    Value{dlit.MustNew(2.6), 7},
				"2.8789": Value{dlit.MustNew(2.8789), 1},
				"3.3":    Value{dlit.MustNew(3.3), 6},
				"5":      Value{dlit.MustNew(5), 7},
				"2":      Value{dlit.MustNew(2), 7},
				"2.8":    Value{dlit.MustNew(2.8), 6},
				"8.8":    Value{dlit.MustNew(8.8), 6},
			},
			6,
		},
		&Field{
			fieldtype.Number,
			dlit.MustNew(7),
			dlit.MustNew(5),
			4,
			map[string]Value{
				"2.6":    Value{dlit.MustNew(2.6), 7},
				"2.8789": Value{dlit.MustNew(2.8789), 1},
				"3.3":    Value{dlit.MustNew(3.3), 3},
				"5":      Value{dlit.MustNew(5), 7},
				"2":      Value{dlit.MustNew(2), 7},
				"2.8":    Value{dlit.MustNew(2.8), 6},
			},
			6,
		},
	}
	cases := []struct {
		ndxA int
		ndxB int
		want error
	}{
		{0, 0, nil},
		{0, 1, errors.New("NumValues not equal: 31 != 18")},
		{0, 2, errors.New("Kind not equal: String != Number")},
		{2, 3, errors.New("Min not equal: 2 != 7")},
		{2, 4, errors.New("Max not equal: 5 != 4")},
		{2, 5, errors.New("MaxDP not equal: 4 != 2")},
		{6, 7, errors.New("number of Values not equal: 6 != 7")},
		{3, 6, errors.New("Value missing: 3")},
		{6, 8, errors.New("Value not equal for: 3.3, {3.3 6} != {3.3 3}")},
	}
	for i, c := range cases {
		got := fields[c.ndxA].checkEqual(fields[c.ndxB])
		checkErrorMatch(t, fmt.Sprintf("(%d) CheckEqual: ", i), got, c.want)
	}
}

func TestDescriptionCalcFieldNum(t *testing.T) {
	description := &Description{
		map[string]*Field{
			"inputA": &Field{
				fieldtype.Number,
				dlit.MustNew(7),
				dlit.MustNew(15.1),
				1,
				map[string]Value{
					"7":    Value{dlit.MustNew(7), 7},
					"7.3":  Value{dlit.MustNew(7.3), 7},
					"9":    Value{dlit.MustNew(9), 7},
					"14":   Value{dlit.MustNew(14), 7},
					"15.1": Value{dlit.MustNew(15.1), 7},
				},
				5,
			},
			"band": &Field{fieldtype.String, nil, nil, 0,
				map[string]Value{
					"a": Value{dlit.MustNew("a"), 2},
					"b": Value{dlit.MustNew("b"), 3},
					"c": Value{dlit.MustNew("c"), 70},
					"f": Value{dlit.MustNew("f"), 22},
					"9": Value{dlit.MustNew("9"), 1},
				},
				31,
			},
			"inputB": &Field{
				fieldtype.Number,
				dlit.MustNew(2),
				dlit.MustNew(5),
				4,
				map[string]Value{
					"2.6":    Value{dlit.MustNew(2.6), 7},
					"2.8789": Value{dlit.MustNew(2.8789), 1},
					"3":      Value{dlit.MustNew(3), 7},
					"5":      Value{dlit.MustNew(5), 7},
					"2":      Value{dlit.MustNew(2), 7},
					"2.8":    Value{dlit.MustNew(2.8), 6},
				},
				6,
			},
		},
	}
	cases := []struct {
		field string
		want  int
	}{
		{"band", 0},
		{"inputA", 1},
		{"inputB", 2},
	}
	for i, c := range cases {
		got := CalcFieldNum(description.Fields, c.field)
		if got != c.want {
			t.Errorf("(%d) CalcFieldNum: got: %d, want: %d", i, got, c.want)
		}
	}
}

func TestDescriptionCalcFieldNum_panic(t *testing.T) {
	description := &Description{
		map[string]*Field{
			"inputA": &Field{
				fieldtype.Number,
				dlit.MustNew(7),
				dlit.MustNew(15.1),
				1,
				map[string]Value{
					"7":    Value{dlit.MustNew(7), 7},
					"7.3":  Value{dlit.MustNew(7.3), 7},
					"9":    Value{dlit.MustNew(9), 7},
					"14":   Value{dlit.MustNew(14), 7},
					"15.1": Value{dlit.MustNew(15.1), 7},
				},
				5,
			},
		},
	}
	paniced := false
	field := "borris"
	wantPanic := "can't find field in Field descriptions: " + field
	defer func() {
		if r := recover(); r != nil {
			if r.(string) == wantPanic {
				paniced = true
			} else {
				t.Errorf("CalcFieldNum: got panic: %s, want: %s", r, wantPanic)
			}
		}
	}()
	got := CalcFieldNum(description.Fields, field)
	if !paniced {
		t.Errorf("CalcFieldNum: got: %s, failed to panic with: %s", got, wantPanic)
	}
}

/*************************************
 *  Helper functions
 *************************************/
func checkErrorMatch(t *testing.T, context string, got, want error) {
	if got == nil && want == nil {
		return
	}
	if got == nil || want == nil {
		t.Errorf("%s got err: %s, want : %s", context, got, want)
	}
	if perr, ok := want.(*os.PathError); ok {
		if err := checkPathErrorMatch(got, perr); err != nil {
			t.Errorf("%s %s", context, err)
		}
	}
	if got.Error() != want.Error() {
		t.Errorf("%s got err: %s, want : %s", context, got, want)
	}
}

func checkPathErrorMatch(checkErr error, wantErr *os.PathError) error {
	perr, ok := checkErr.(*os.PathError)
	if !ok {
		return fmt.Errorf("got err type: %T, want error type: os.PathError",
			checkErr)
	}
	if perr.Op != wantErr.Op {
		return fmt.Errorf("got perr.Op: %s, want: %s", perr.Op, wantErr.Op)
	}
	if filepath.Clean(perr.Path) != filepath.Clean(wantErr.Path) {
		return fmt.Errorf("got perr.Path: %s, want: %s", perr.Path, wantErr.Path)
	}
	if perr.Err != wantErr.Err {
		return fmt.Errorf("got perr.Err: %s, want: %s", perr.Err, wantErr.Err)
	}
	return nil
}
