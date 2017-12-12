package description

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/lawrencewoodman/ddataset"
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/internal/testhelpers"
)

var flowRecords = [][]string{
	{"a", "7", "2.6", "9.9", "22", "1", "a"},
	{"c", "7.3", "2.8789", "9.97", "21", "4", "b"},
	{"b", "9", "3", "10", "23", "2", "c"},
	{"f", "14", "5", "10.94", "25", "3", "d"},
	{"b", "15.1", "2", "9.9a", "27", "5", "e"},

	{"g", "7", "2.6", "9.9", "32", "5", "f"},
	{"i", "7.3", "2.8", "9.97", "31", "4", "g"},
	{"k", "9", "3", "10", "33", "1", "h"},
	{"l", "14", "5", "10.94", "35", "3", "i"},
	{"m", "15.1", "2", "9.9a", "37", "2", "j"},

	{"z", "7", "2.6", "9.9", "42", "5", "k"},
	{"u", "7.3", "2.8", "9.97", "41", "5", "l"},
	{"b", "9", "3", "10", "43", "2", "m"},
	{"a", "14", "5", "10.94", "45", "1", "n"},
	{"n", "15.1", "2", "9.9a", "47", "4", "o"},

	{"t", "7", "2.6", "9.9", "22", "3", "p"},
	{"s", "7.3", "2.8", "9.97", "21", "5", "q"},
	{"x", "9", "3", "10", "53", "2", "r"},
	{"y", "14", "5", "10.94", "55", "3", "s"},
	{"v", "15.1", "2", "9.9a", "57", "4", "t"},

	{"h", "7", "2.6", "9.9", "62", "1", "u"},
	{"j", "7.3", "2.8", "9.97", "61", "5", "v"},
	{"o", "9", "3", "10", "63", "4", "w"},
	{"p", "14", "5", "10.94", "65", "3", "x"},
	{"q", "15.1", "2", "9.9a", "27", "2", "y"},

	{"9", "7", "2.6", "9.9", "72", "4", "z"},
	{"7", "7.3", "2.8", "9.97", "71", "5", "aa"},
	{"6", "9", "3", "10", "73", "4", "ab"},
	{"5", "14", "5", "10.94", "75", "2", "ac"},
	{"4", "15.1", "2", "9.9a", "77", "1", "ad"},

	{"8", "7", "2.6", "9.9", "82", "5", "ae"},
	{"1", "7.3", "2.8", "9.97", "81", "4", "af"},
	{"2", "9", "3", "10", "83", "3", "a"},
	{"3", "14", "5", "10.94", "85", "2", "b"},
	{"8", "15.1", "2", "9.9b", "87", "1", "c"},
}

func TestDescribeDataset(t *testing.T) {
	fieldNames :=
		[]string{"band", "inputA", "inputB", "version", "flow", "score", "method"}
	expected := &Description{
		map[string]*Field{
			"band": {String, nil, nil, 0,
				map[string]Value{
					"a": {dlit.MustNew("a"), 2},
					"b": {dlit.MustNew("b"), 3},
					"c": {dlit.MustNew("c"), 1},
					"f": {dlit.MustNew("f"), 1},
					"g": {dlit.MustNew("g"), 1},
					"h": {dlit.MustNew("h"), 1},
					"i": {dlit.MustNew("i"), 1},
					"j": {dlit.MustNew("j"), 1},
					"k": {dlit.MustNew("k"), 1},
					"l": {dlit.MustNew("l"), 1},
					"m": {dlit.MustNew("m"), 1},
					"n": {dlit.MustNew("n"), 1},
					"o": {dlit.MustNew("o"), 1},
					"p": {dlit.MustNew("p"), 1},
					"q": {dlit.MustNew("q"), 1},
					"s": {dlit.MustNew("s"), 1},
					"t": {dlit.MustNew("t"), 1},
					"u": {dlit.MustNew("u"), 1},
					"v": {dlit.MustNew("v"), 1},
					"x": {dlit.MustNew("x"), 1},
					"y": {dlit.MustNew("y"), 1},
					"z": {dlit.MustNew("z"), 1},
					"1": {dlit.MustNew("1"), 1},
					"2": {dlit.MustNew("2"), 1},
					"3": {dlit.MustNew("3"), 1},
					"4": {dlit.MustNew("4"), 1},
					"5": {dlit.MustNew("5"), 1},
					"6": {dlit.MustNew("6"), 1},
					"7": {dlit.MustNew("7"), 1},
					"8": {dlit.MustNew("8"), 2},
					"9": {dlit.MustNew("9"), 1},
				},
				31,
			},
			"inputA": {
				Number,
				dlit.MustNew(7),
				dlit.MustNew(15.1),
				1,
				map[string]Value{
					"7":    {dlit.MustNew(7), 7},
					"7.3":  {dlit.MustNew(7.3), 7},
					"9":    {dlit.MustNew(9), 7},
					"14":   {dlit.MustNew(14), 7},
					"15.1": {dlit.MustNew(15.1), 7},
				},
				5,
			},
			"inputB": {
				Number,
				dlit.MustNew(2),
				dlit.MustNew(5),
				4,
				map[string]Value{
					"2.6":    {dlit.MustNew(2.6), 7},
					"2.8789": {dlit.MustNew(2.8789), 1},
					"3":      {dlit.MustNew(3), 7},
					"5":      {dlit.MustNew(5), 7},
					"2":      {dlit.MustNew(2), 7},
					"2.8":    {dlit.MustNew(2.8), 6},
				},
				6,
			},
			"version": {String, nil, nil, 0,
				map[string]Value{
					"9.9":   {dlit.MustNew("9.9"), 7},
					"9.97":  {dlit.MustNew("9.97"), 7},
					"10":    {dlit.MustNew("10"), 7},
					"10.94": {dlit.MustNew("10.94"), 7},
					"9.9a":  {dlit.MustNew("9.9a"), 6},
					"9.9b":  {dlit.MustNew("9.9b"), 1},
				},
				6,
			},
			"flow": {
				Number,
				dlit.MustNew(21),
				dlit.MustNew(87),
				0,
				map[string]Value{}, -1},
			"score": {
				Number,
				dlit.MustNew(1),
				dlit.MustNew(5),
				0,
				map[string]Value{
					"1": {dlit.MustNew(1), 6},
					"2": {dlit.MustNew(2), 7},
					"3": {dlit.MustNew(3), 6},
					"4": {dlit.MustNew(4), 8},
					"5": {dlit.MustNew(5), 8},
				}, 5,
			},
			"method": {Ignore, nil, nil, 0,
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

func TestDescribeDataset_dataset_errors(t *testing.T) {
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

func TestDescribeDataset_field_errors(t *testing.T) {
	fieldNames :=
		[]string{"band", "5hello", "inputB", "version", "flow", "score", "method"}
	dataset := testhelpers.NewLiteralDataset(fieldNames, flowRecords)
	wantErr := InvalidFieldError("5hello")
	_, err := DescribeDataset(dataset)
	if err == nil || err.Error() != wantErr.Error() {
		t.Errorf("DescribeDataset - err: %s, wantErr: %s", err, wantErr)
	}
}

func TestDescriptionMarshalUnmarshalJSON(t *testing.T) {
	description := &Description{
		map[string]*Field{
			"band": {String, nil, nil, 0,
				map[string]Value{
					"a": {dlit.MustNew("a"), 2},
					"b": {dlit.MustNew("b"), 3},
					"c": {dlit.MustNew("c"), 70},
					"f": {dlit.MustNew("f"), 22},
					"9": {dlit.MustNew("9"), 1},
				},
				31,
			},
			"inputA": {
				Number,
				dlit.MustNew(7),
				dlit.MustNew(15.1),
				1,
				map[string]Value{
					"7":    {dlit.MustNew(7), 7},
					"7.3":  {dlit.MustNew(7.3), 7},
					"9":    {dlit.MustNew(9), 7},
					"14":   {dlit.MustNew(14), 7},
					"15.1": {dlit.MustNew(15.1), 7},
				},
				5,
			},
			"inputB": {
				Number,
				dlit.MustNew(2),
				dlit.MustNew(5),
				4,
				map[string]Value{
					"2.6":    {dlit.MustNew(2.6), 7},
					"2.8789": {dlit.MustNew(2.8789), 1},
					"3":      {dlit.MustNew(3), 7},
					"5":      {dlit.MustNew(5), 7},
					"2":      {dlit.MustNew(2), 7},
					"2.8":    {dlit.MustNew(2.8), 6},
				},
				6,
			},
			"version": {String, nil, nil, 0,
				map[string]Value{
					"9.9":   {dlit.MustNew("9.9"), 7},
					"9.97":  {dlit.MustNew("9.97"), 7},
					"10":    {dlit.MustNew("10"), 7},
					"10.94": {dlit.MustNew("10.94"), 7},
					"9.9a":  {dlit.MustNew("9.9a"), 6},
					"9.9b":  {dlit.MustNew("9.9b"), 1},
				},
				6,
			},
			"flow": {
				Number,
				dlit.MustNew(21),
				dlit.MustNew(87),
				0,
				map[string]Value{}, -1},
			"score": {
				Number,
				dlit.MustNew(1),
				dlit.MustNew(5),
				0,
				map[string]Value{
					"1": {dlit.MustNew(1), 6},
					"2": {dlit.MustNew(2), 7},
					"3": {dlit.MustNew(3), 6},
					"4": {dlit.MustNew(4), 8},
					"5": {dlit.MustNew(5), 8},
				}, 5,
			},
			"method": {Ignore, nil, nil, 0,
				map[string]Value{}, -1},
		},
	}
	b, err := json.Marshal(description)
	if err != nil {
		t.Fatalf("Marshal: %s", err)
	}
	var got Description
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("Unmarshal: %s", err)
	}
	if err := got.CheckEqual(description); err != nil {
		t.Errorf("Unmarshal got not expected: %s", err)
	}
}

func TestDescriptionCheckEqual(t *testing.T) {
	descriptions := []*Description{
		{
			map[string]*Field{
				"band": {String, nil, nil, 0,
					map[string]Value{
						"a": {dlit.MustNew("a"), 2},
						"b": {dlit.MustNew("b"), 3},
						"c": {dlit.MustNew("c"), 70},
						"f": {dlit.MustNew("f"), 22},
						"9": {dlit.MustNew("9"), 1},
					},
					31,
				},
				"inputA": {
					Number,
					dlit.MustNew(7),
					dlit.MustNew(15.1),
					1,
					map[string]Value{
						"7":    {dlit.MustNew(7), 7},
						"7.3":  {dlit.MustNew(7.3), 7},
						"9":    {dlit.MustNew(9), 7},
						"14":   {dlit.MustNew(14), 7},
						"15.1": {dlit.MustNew(15.1), 7},
					},
					5,
				},
				"inputB": {
					Number,
					dlit.MustNew(2),
					dlit.MustNew(5),
					4,
					map[string]Value{
						"2.6":    {dlit.MustNew(2.6), 7},
						"2.8789": {dlit.MustNew(2.8789), 1},
						"3":      {dlit.MustNew(3), 7},
						"5":      {dlit.MustNew(5), 7},
						"2":      {dlit.MustNew(2), 7},
						"2.8":    {dlit.MustNew(2.8), 6},
					},
					6,
				},
			},
		},
		{
			map[string]*Field{
				"band": {String, nil, nil, 0,
					map[string]Value{
						"a": {dlit.MustNew("a"), 2},
						"b": {dlit.MustNew("b"), 3},
						"c": {dlit.MustNew("c"), 70},
						"f": {dlit.MustNew("f"), 22},
						"9": {dlit.MustNew("9"), 1},
					},
					31,
				},
				"inputB": {
					Number,
					dlit.MustNew(2),
					dlit.MustNew(5),
					4,
					map[string]Value{
						"2.6":    {dlit.MustNew(2.6), 7},
						"2.8789": {dlit.MustNew(2.8789), 1},
						"3":      {dlit.MustNew(3), 7},
						"5":      {dlit.MustNew(5), 7},
						"2":      {dlit.MustNew(2), 7},
						"2.8":    {dlit.MustNew(2.8), 6},
					},
					6,
				},
			},
		},
		{
			map[string]*Field{
				"strata": {String, nil, nil, 0,
					map[string]Value{
						"a": {dlit.MustNew("a"), 2},
						"b": {dlit.MustNew("b"), 3},
						"c": {dlit.MustNew("c"), 70},
						"f": {dlit.MustNew("f"), 22},
						"9": {dlit.MustNew("9"), 1},
					},
					31,
				},
				"inputA": {
					Number,
					dlit.MustNew(7),
					dlit.MustNew(15.1),
					1,
					map[string]Value{
						"7":    {dlit.MustNew(7), 7},
						"7.3":  {dlit.MustNew(7.3), 7},
						"9":    {dlit.MustNew(9), 7},
						"14":   {dlit.MustNew(14), 7},
						"15.1": {dlit.MustNew(15.1), 7},
					},
					5,
				},
				"inputB": {
					Number,
					dlit.MustNew(2),
					dlit.MustNew(5),
					4,
					map[string]Value{
						"2.6":    {dlit.MustNew(2.6), 7},
						"2.8789": {dlit.MustNew(2.8789), 1},
						"3":      {dlit.MustNew(3), 7},
						"5":      {dlit.MustNew(5), 7},
						"2":      {dlit.MustNew(2), 7},
						"2.8":    {dlit.MustNew(2.8), 6},
					},
					6,
				},
			},
		},
		{
			map[string]*Field{
				"band": {String, nil, nil, 0,
					map[string]Value{
						"a": {dlit.MustNew("a"), 2},
						"b": {dlit.MustNew("b"), 3},
						"c": {dlit.MustNew("c"), 70},
						"f": {dlit.MustNew("f"), 22},
						"9": {dlit.MustNew("9"), 1},
					},
					31,
				},
				"inputA": {
					Number,
					dlit.MustNew(6),
					dlit.MustNew(15.1),
					1,
					map[string]Value{
						"7":    {dlit.MustNew(7), 7},
						"7.3":  {dlit.MustNew(7.3), 7},
						"9":    {dlit.MustNew(9), 7},
						"14":   {dlit.MustNew(14), 7},
						"15.1": {dlit.MustNew(15.1), 7},
					},
					5,
				},
				"inputB": {
					Number,
					dlit.MustNew(2),
					dlit.MustNew(5),
					4,
					map[string]Value{
						"2.6":    {dlit.MustNew(2.6), 7},
						"2.8789": {dlit.MustNew(2.8789), 1},
						"3":      {dlit.MustNew(3), 7},
						"5":      {dlit.MustNew(5), 7},
						"2":      {dlit.MustNew(2), 7},
						"2.8":    {dlit.MustNew(2.8), 6},
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

func TestDescriptionFieldNames(t *testing.T) {
	description := &Description{
		Fields: map[string]*Field{
			"band": {String, nil, nil, 0,
				map[string]Value{
					"a": {dlit.MustNew("a"), 2},
					"b": {dlit.MustNew("b"), 3},
					"c": {dlit.MustNew("c"), 70},
					"f": {dlit.MustNew("f"), 22},
					"9": {dlit.MustNew("9"), 1},
				},
				31,
			},
			"inputA": {
				Number,
				dlit.MustNew(7),
				dlit.MustNew(15.1),
				1,
				map[string]Value{
					"7":    {dlit.MustNew(7), 7},
					"7.3":  {dlit.MustNew(7.3), 7},
					"9":    {dlit.MustNew(9), 7},
					"14":   {dlit.MustNew(14), 7},
					"15.1": {dlit.MustNew(15.1), 7},
				},
				5,
			},
			"inputB": {
				Number,
				dlit.MustNew(2),
				dlit.MustNew(5),
				4,
				map[string]Value{
					"2.6":    {dlit.MustNew(2.6), 7},
					"2.8789": {dlit.MustNew(2.8789), 1},
					"3":      {dlit.MustNew(3), 7},
					"5":      {dlit.MustNew(5), 7},
					"2":      {dlit.MustNew(2), 7},
					"2.8":    {dlit.MustNew(2.8), 6},
				},
				6,
			},
		},
	}
	want := []string{"band", "inputA", "inputB"}
	got := description.FieldNames()
	sort.Strings(got)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FieldNames - got: %v, want: %v", got, want)
	}
}

func TestFieldCheckEqual(t *testing.T) {
	fields := []*Field{
		{String, nil, nil, 0,
			map[string]Value{
				"a": {dlit.MustNew("a"), 2},
				"b": {dlit.MustNew("b"), 3},
				"c": {dlit.MustNew("c"), 70},
				"f": {dlit.MustNew("f"), 22},
				"9": {dlit.MustNew("9"), 1},
			},
			31,
		},
		{String, nil, nil, 0,
			map[string]Value{
				"a": {dlit.MustNew("a"), 2},
				"b": {dlit.MustNew("b"), 3},
				"c": {dlit.MustNew("c"), 70},
				"f": {dlit.MustNew("f"), 22},
				"9": {dlit.MustNew("9"), 1},
			},
			18,
		},
		{
			Number,
			dlit.MustNew(2),
			dlit.MustNew(5),
			4,
			map[string]Value{
				"2.6":    {dlit.MustNew(2.6), 7},
				"2.8789": {dlit.MustNew(2.8789), 1},
				"3":      {dlit.MustNew(3), 7},
				"5":      {dlit.MustNew(5), 7},
				"2":      {dlit.MustNew(2), 7},
				"2.8":    {dlit.MustNew(2.8), 6},
			},
			6,
		},
		{
			Number,
			dlit.MustNew(7),
			dlit.MustNew(5),
			4,
			map[string]Value{
				"2.6":    {dlit.MustNew(2.6), 7},
				"2.8789": {dlit.MustNew(2.8789), 1},
				"3":      {dlit.MustNew(3), 7},
				"5":      {dlit.MustNew(5), 7},
				"2":      {dlit.MustNew(2), 7},
				"2.8":    {dlit.MustNew(2.8), 6},
			},
			6,
		},
		{
			Number,
			dlit.MustNew(2),
			dlit.MustNew(4),
			4,
			map[string]Value{
				"2.6":    {dlit.MustNew(2.6), 7},
				"2.8789": {dlit.MustNew(2.8789), 1},
				"3":      {dlit.MustNew(3), 7},
				"5":      {dlit.MustNew(5), 7},
				"2":      {dlit.MustNew(2), 7},
				"2.8":    {dlit.MustNew(2.8), 6},
			},
			6,
		},
		{
			Number,
			dlit.MustNew(2),
			dlit.MustNew(5),
			2,
			map[string]Value{
				"2.6":    {dlit.MustNew(2.6), 7},
				"2.8789": {dlit.MustNew(2.8789), 1},
				"3":      {dlit.MustNew(3), 7},
				"5":      {dlit.MustNew(5), 7},
				"2":      {dlit.MustNew(2), 7},
				"2.8":    {dlit.MustNew(2.8), 6},
			},
			6,
		},
		{
			Number,
			dlit.MustNew(7),
			dlit.MustNew(5),
			4,
			map[string]Value{
				"2.6":    {dlit.MustNew(2.6), 7},
				"2.8789": {dlit.MustNew(2.8789), 1},
				"3.3":    {dlit.MustNew(3.3), 6},
				"5":      {dlit.MustNew(5), 7},
				"2":      {dlit.MustNew(2), 7},
				"2.8":    {dlit.MustNew(2.8), 6},
			},
			6,
		},
		{
			Number,
			dlit.MustNew(7),
			dlit.MustNew(5),
			4,
			map[string]Value{
				"2.6":    {dlit.MustNew(2.6), 7},
				"2.8789": {dlit.MustNew(2.8789), 1},
				"3.3":    {dlit.MustNew(3.3), 6},
				"5":      {dlit.MustNew(5), 7},
				"2":      {dlit.MustNew(2), 7},
				"2.8":    {dlit.MustNew(2.8), 6},
				"8.8":    {dlit.MustNew(8.8), 6},
			},
			6,
		},
		{
			Number,
			dlit.MustNew(7),
			dlit.MustNew(5),
			4,
			map[string]Value{
				"2.6":    {dlit.MustNew(2.6), 7},
				"2.8789": {dlit.MustNew(2.8789), 1},
				"3.3":    {dlit.MustNew(3.3), 3},
				"5":      {dlit.MustNew(5), 7},
				"2":      {dlit.MustNew(2), 7},
				"2.8":    {dlit.MustNew(2.8), 6},
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
			"inputA": {
				Number,
				dlit.MustNew(7),
				dlit.MustNew(15.1),
				1,
				map[string]Value{
					"7":    {dlit.MustNew(7), 7},
					"7.3":  {dlit.MustNew(7.3), 7},
					"9":    {dlit.MustNew(9), 7},
					"14":   {dlit.MustNew(14), 7},
					"15.1": {dlit.MustNew(15.1), 7},
				},
				5,
			},
			"band": {String, nil, nil, 0,
				map[string]Value{
					"a": {dlit.MustNew("a"), 2},
					"b": {dlit.MustNew("b"), 3},
					"c": {dlit.MustNew("c"), 70},
					"f": {dlit.MustNew("f"), 22},
					"9": {dlit.MustNew("9"), 1},
				},
				31,
			},
			"inputB": {
				Number,
				dlit.MustNew(2),
				dlit.MustNew(5),
				4,
				map[string]Value{
					"2.6":    {dlit.MustNew(2.6), 7},
					"2.8789": {dlit.MustNew(2.8789), 1},
					"3":      {dlit.MustNew(3), 7},
					"5":      {dlit.MustNew(5), 7},
					"2":      {dlit.MustNew(2), 7},
					"2.8":    {dlit.MustNew(2.8), 6},
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
			"inputA": {
				Number,
				dlit.MustNew(7),
				dlit.MustNew(15.1),
				1,
				map[string]Value{
					"7":    {dlit.MustNew(7), 7},
					"7.3":  {dlit.MustNew(7.3), 7},
					"9":    {dlit.MustNew(9), 7},
					"14":   {dlit.MustNew(14), 7},
					"15.1": {dlit.MustNew(15.1), 7},
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
		t.Errorf("CalcFieldNum: got: %d, failed to panic with: %s", got, wantPanic)
	}
}

func TestInvalidFieldErrorError(t *testing.T) {
	err := InvalidFieldError("bob")
	want := "invalid field: bob"
	got := err.Error()
	if got != want {
		t.Errorf("Error() got: %s, want: %s", got, want)
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

func (l *FailingDataset) Release() error {
	return nil
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
