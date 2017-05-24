package rhkit

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/description"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
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
	expected := &description.Description{
		map[string]*description.Field{
			"band": &description.Field{fieldtype.String, nil, nil, 0,
				map[string]description.Value{
					"a": description.Value{dlit.MustNew("a"), 2},
					"b": description.Value{dlit.MustNew("b"), 3},
					"c": description.Value{dlit.MustNew("c"), 1},
					"f": description.Value{dlit.MustNew("f"), 1},
					"g": description.Value{dlit.MustNew("g"), 1},
					"h": description.Value{dlit.MustNew("h"), 1},
					"i": description.Value{dlit.MustNew("i"), 1},
					"j": description.Value{dlit.MustNew("j"), 1},
					"k": description.Value{dlit.MustNew("k"), 1},
					"l": description.Value{dlit.MustNew("l"), 1},
					"m": description.Value{dlit.MustNew("m"), 1},
					"n": description.Value{dlit.MustNew("n"), 1},
					"o": description.Value{dlit.MustNew("o"), 1},
					"p": description.Value{dlit.MustNew("p"), 1},
					"q": description.Value{dlit.MustNew("q"), 1},
					"s": description.Value{dlit.MustNew("s"), 1},
					"t": description.Value{dlit.MustNew("t"), 1},
					"u": description.Value{dlit.MustNew("u"), 1},
					"v": description.Value{dlit.MustNew("v"), 1},
					"x": description.Value{dlit.MustNew("x"), 1},
					"y": description.Value{dlit.MustNew("y"), 1},
					"z": description.Value{dlit.MustNew("z"), 1},
					"1": description.Value{dlit.MustNew("1"), 1},
					"2": description.Value{dlit.MustNew("2"), 1},
					"3": description.Value{dlit.MustNew("3"), 1},
					"4": description.Value{dlit.MustNew("4"), 1},
					"5": description.Value{dlit.MustNew("5"), 1},
					"6": description.Value{dlit.MustNew("6"), 1},
					"7": description.Value{dlit.MustNew("7"), 1},
					"8": description.Value{dlit.MustNew("8"), 2},
					"9": description.Value{dlit.MustNew("9"), 1},
				},
				31,
			},
			"inputA": &description.Field{
				fieldtype.Number,
				dlit.MustNew(7),
				dlit.MustNew(15.1),
				1,
				map[string]description.Value{
					"7":    description.Value{dlit.MustNew(7), 7},
					"7.3":  description.Value{dlit.MustNew(7.3), 7},
					"9":    description.Value{dlit.MustNew(9), 7},
					"14":   description.Value{dlit.MustNew(14), 7},
					"15.1": description.Value{dlit.MustNew(15.1), 7},
				},
				5,
			},
			"inputB": &description.Field{
				fieldtype.Number,
				dlit.MustNew(2),
				dlit.MustNew(5),
				4,
				map[string]description.Value{
					"2.6":    description.Value{dlit.MustNew(2.6), 7},
					"2.8789": description.Value{dlit.MustNew(2.8789), 1},
					"3":      description.Value{dlit.MustNew(3), 7},
					"5":      description.Value{dlit.MustNew(5), 7},
					"2":      description.Value{dlit.MustNew(2), 7},
					"2.8":    description.Value{dlit.MustNew(2.8), 6},
				},
				6,
			},
			"version": &description.Field{fieldtype.String, nil, nil, 0,
				map[string]description.Value{
					"9.9":   description.Value{dlit.MustNew("9.9"), 7},
					"9.97":  description.Value{dlit.MustNew("9.97"), 7},
					"10":    description.Value{dlit.MustNew("10"), 7},
					"10.94": description.Value{dlit.MustNew("10.94"), 7},
					"9.9a":  description.Value{dlit.MustNew("9.9a"), 6},
					"9.9b":  description.Value{dlit.MustNew("9.9b"), 1},
				},
				6,
			},
			"flow": &description.Field{
				fieldtype.Number,
				dlit.MustNew(21),
				dlit.MustNew(87),
				0,
				map[string]description.Value{}, -1},
			"score": &description.Field{
				fieldtype.Number,
				dlit.MustNew(1),
				dlit.MustNew(5),
				0,
				map[string]description.Value{
					"1": description.Value{dlit.MustNew(1), 6},
					"2": description.Value{dlit.MustNew(2), 7},
					"3": description.Value{dlit.MustNew(3), 6},
					"4": description.Value{dlit.MustNew(4), 8},
					"5": description.Value{dlit.MustNew(5), 8},
				}, 5,
			},
			"method": &description.Field{fieldtype.Ignore, nil, nil, 0,
				map[string]description.Value{}, -1},
		}}
	dataset := NewLiteralDataset(fieldNames, flowRecords)
	d, err := DescribeDataset(dataset)
	if err != nil {
		t.Errorf("DescribeDataset(dataset) err: %s", err)
	}
	if err := d.CheckEqual(expected); err != nil {
		t.Errorf("DescibeDataset(dataset) got not expected: %s", err)
	}
}
