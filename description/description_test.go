package description

import (
	"errors"
	"fmt"
	"github.com/lawrencewoodman/ddataset"
	"github.com/lawrencewoodman/ddataset/dcsv"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
	"github.com/vlifesystems/rhkit/internal/testhelpers"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

var flowRecords = [][]string{
	[]string{"a", "7", "2.6", "9.9", "22", "1", "a"},
	[]string{"c", "7.3", "2.8789", "9.97", "21", "4", "b"},
	[]string{"b", "9", "3", "10", "23", "2", "c"},
	[]string{"f", "14", "5", "10.94", "25", "3", "d"},
	[]string{"b", "15.1", "2", "9.9a", "27", "5", "e"},

	[]string{"g", "7", "2.6", "9.9", "32", "5", "f"},
	[]string{"i", "7.3", "2.8", "9.97", "31", "4", "g"},
	[]string{"k", "9", "3", "10", "33", "1", "h"},
	[]string{"l", "14", "5", "10.94", "35", "3", "i"},
	[]string{"m", "15.1", "2", "9.9a", "37", "2", "j"},

	[]string{"z", "7", "2.6", "9.9", "42", "5", "k"},
	[]string{"u", "7.3", "2.8", "9.97", "41", "5", "l"},
	[]string{"b", "9", "3", "10", "43", "2", "m"},
	[]string{"a", "14", "5", "10.94", "45", "1", "n"},
	[]string{"n", "15.1", "2", "9.9a", "47", "4", "o"},

	[]string{"t", "7", "2.6", "9.9", "22", "3", "p"},
	[]string{"s", "7.3", "2.8", "9.97", "21", "5", "q"},
	[]string{"x", "9", "3", "10", "53", "2", "r"},
	[]string{"y", "14", "5", "10.94", "55", "3", "s"},
	[]string{"v", "15.1", "2", "9.9a", "57", "4", "t"},

	[]string{"h", "7", "2.6", "9.9", "62", "1", "u"},
	[]string{"j", "7.3", "2.8", "9.97", "61", "5", "v"},
	[]string{"o", "9", "3", "10", "63", "4", "w"},
	[]string{"p", "14", "5", "10.94", "65", "3", "x"},
	[]string{"q", "15.1", "2", "9.9a", "27", "2", "y"},

	[]string{"9", "7", "2.6", "9.9", "72", "4", "z"},
	[]string{"7", "7.3", "2.8", "9.97", "71", "5", "aa"},
	[]string{"6", "9", "3", "10", "73", "4", "ab"},
	[]string{"5", "14", "5", "10.94", "75", "2", "ac"},
	[]string{"4", "15.1", "2", "9.9a", "77", "1", "ad"},

	[]string{"8", "7", "2.6", "9.9", "82", "5", "ae"},
	[]string{"1", "7.3", "2.8", "9.97", "81", "4", "af"},
	[]string{"2", "9", "3", "10", "83", "3", "a"},
	[]string{"3", "14", "5", "10.94", "85", "2", "b"},
	[]string{"8", "15.1", "2", "9.9b", "87", "1", "c"},
}

func TestDescribeDataset(t *testing.T) {
	fieldNames :=
		[]string{"band", "inputA", "inputB", "version", "flow", "score", "method"}
	expected := &Description{
		map[string]*Field{
			"band": &Field{fieldtype.String, nil, nil, 0,
				map[string]Value{
					"a": Value{dlit.MustNew("a"), 2},
					"b": Value{dlit.MustNew("b"), 3},
					"c": Value{dlit.MustNew("c"), 1},
					"f": Value{dlit.MustNew("f"), 1},
					"g": Value{dlit.MustNew("g"), 1},
					"h": Value{dlit.MustNew("h"), 1},
					"i": Value{dlit.MustNew("i"), 1},
					"j": Value{dlit.MustNew("j"), 1},
					"k": Value{dlit.MustNew("k"), 1},
					"l": Value{dlit.MustNew("l"), 1},
					"m": Value{dlit.MustNew("m"), 1},
					"n": Value{dlit.MustNew("n"), 1},
					"o": Value{dlit.MustNew("o"), 1},
					"p": Value{dlit.MustNew("p"), 1},
					"q": Value{dlit.MustNew("q"), 1},
					"s": Value{dlit.MustNew("s"), 1},
					"t": Value{dlit.MustNew("t"), 1},
					"u": Value{dlit.MustNew("u"), 1},
					"v": Value{dlit.MustNew("v"), 1},
					"x": Value{dlit.MustNew("x"), 1},
					"y": Value{dlit.MustNew("y"), 1},
					"z": Value{dlit.MustNew("z"), 1},
					"1": Value{dlit.MustNew("1"), 1},
					"2": Value{dlit.MustNew("2"), 1},
					"3": Value{dlit.MustNew("3"), 1},
					"4": Value{dlit.MustNew("4"), 1},
					"5": Value{dlit.MustNew("5"), 1},
					"6": Value{dlit.MustNew("6"), 1},
					"7": Value{dlit.MustNew("7"), 1},
					"8": Value{dlit.MustNew("8"), 2},
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
		}}
	dataset := testhelpers.NewLiteralDataset(fieldNames, flowRecords)
	d, err := DescribeDataset(dataset)
	if err != nil {
		t.Errorf("DescribeDataset(dataset) err: %s", err)
	}
	if err := d.CheckEqual(expected); err != nil {
		t.Errorf("DescibeDataset(dataset) got not expected: %s", err)
	}
}

func TestDescribeDataset_errors(t *testing.T) {
	fieldNames :=
		[]string{"band", "inputA", "inputB", "version", "flow", "score", "method"}
	dataset := testhelpers.NewLiteralDataset(fieldNames, flowRecords)
	cases := []struct {
		stage int
		err   error
	}{
		{stage: 0, err: errors.New("can't open database")},
		{stage: 1, err: errors.New("read error")},
	}
	for _, c := range cases {
		fdataset := NewFailingDataset(dataset, c.stage, c.err)
		_, err := DescribeDataset(fdataset)
		if err == nil || err.Error() != c.err.Error() {
			t.Errorf("DescribeDataset(dataset) err: %s, want: %s", err, c.err)
		}
	}
}
func TestDescriptionNew(t *testing.T) {
	got := New()
	if len(got.Fields) != 0 {
		t.Errorf("New got len(got.Fields): %d, want: 0", len(got.Fields))
	}
}

func TestDescriptionNextRecord(t *testing.T) {
	wantDescription := &Description{
		map[string]*Field{
			"platform": &Field{fieldtype.Ignore, nil, nil, 0,
				map[string]Value{},
				-1,
			},
			"month": &Field{fieldtype.String, nil, nil, 0,
				map[string]Value{
					"jan":   Value{dlit.MustNew("jan"), 3},
					"feb":   Value{dlit.MustNew("feb"), 3},
					"march": Value{dlit.MustNew("march"), 3},
					"april": Value{dlit.MustNew("april"), 3},
					"may":   Value{dlit.MustNew("may"), 12},
					"june":  Value{dlit.MustNew("june"), 3},
					"july":  Value{dlit.MustNew("july"), 9},
					"aug":   Value{dlit.MustNew("aug"), 3},
					"sept":  Value{dlit.MustNew("sept"), 3},
				},
				9,
			},
			"grade": &Field{fieldtype.String, nil, nil, 0,
				map[string]Value{
					"0":     Value{dlit.MustNew("0"), 3},
					"5":     Value{dlit.MustNew("5"), 3},
					"707":   Value{dlit.MustNew("707"), 3},
					"22":    Value{dlit.MustNew("22"), 5},
					"alpha": Value{dlit.MustNew("alpha"), 3},
					"15":    Value{dlit.MustNew("15"), 6},
					"beta":  Value{dlit.MustNew("beta"), 3},
					"gamma": Value{dlit.MustNew("gamma"), 3},
					"23":    Value{dlit.MustNew("23"), 3},
					"-22":   Value{dlit.MustNew("-22"), 1},
					"delta": Value{dlit.MustNew("delta"), 3},
					"6":     Value{dlit.MustNew("6"), 3},
					"98":    Value{dlit.MustNew("98"), 3},
				},
				13,
			},
			"rate": &Field{
				fieldtype.Number,
				dlit.MustNew(-9.3456),
				dlit.MustNew(282),
				4,
				map[string]Value{},
				-1,
			},
			"success": &Field{fieldtype.String, nil, nil, 0,
				map[string]Value{
					"true":  Value{dlit.MustNew("true"), 24},
					"false": Value{dlit.MustNew("false"), 18},
				},
				2,
			},
		},
	}
	fields := []string{"platform", "month", "grade", "rate", "success"}
	dataset := dcsv.New(
		filepath.Join("fixtures", "launch.csv"),
		true,
		rune(','),
		fields,
	)
	description := New()
	conn, err := dataset.Open()
	if err != nil {
		t.Fatalf("dataset.Open: %s", err)
	}
	defer conn.Close()

	for conn.Next() {
		record := conn.Read()
		description.NextRecord(record)
	}
	if err := conn.Err(); err != nil {
		t.Fatalf("NextRecord Err: %s", err)
	}
	if err := description.CheckEqual(wantDescription); err != nil {
		t.Errorf("NextRecord description not expected: %s", err)
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
		testhelpers.CheckErrorMatch(
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
		testhelpers.CheckErrorMatch(
			t,
			fmt.Sprintf("(%d) CheckEqual: ", i),
			got,
			c.want,
		)
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
		testhelpers.CheckErrorMatch(
			t,
			fmt.Sprintf("(%d) CheckEqual: ", i),
			got,
			c.want,
		)
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

type FailingDataset struct {
	dataset  ddataset.Dataset
	stage    int
	errStage int
	err      error
}

type FailingDatasetConn struct {
	dataset *FailingDataset
	conn    ddataset.Conn
}

func NewFailingDataset(
	dataset ddataset.Dataset,
	stage int,
	err error,
) ddataset.Dataset {
	return &FailingDataset{
		dataset:  dataset,
		stage:    0,
		errStage: stage,
		err:      err,
	}
}

func (l *FailingDataset) Open() (ddataset.Conn, error) {
	if l.errStage == 0 {
		return nil, l.err
	}
	l.stage++
	conn, err := l.dataset.Open()
	if err != nil {
		return nil, err
	}
	return &FailingDatasetConn{
		dataset: l,
		conn:    conn,
	}, nil
}

func (l *FailingDataset) Fields() []string {
	return l.dataset.Fields()
}

func (lc *FailingDatasetConn) Close() error {
	if lc.dataset.stage == lc.dataset.errStage {
		return lc.dataset.err
	}
	lc.dataset.stage++
	return lc.conn.Close()
}

func (lc *FailingDatasetConn) Next() bool {
	if lc.dataset.stage == lc.dataset.errStage {
		return false
	}
	return lc.conn.Next()
}

func (lc *FailingDatasetConn) Read() ddataset.Record {
	return lc.conn.Read()
}

func (lc *FailingDatasetConn) Err() error {
	if lc.dataset.stage == lc.dataset.errStage {
		return lc.dataset.err
	}
	lc.dataset.stage++
	return lc.conn.Err()
}
