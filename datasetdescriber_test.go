package rhkit

import (
	"fmt"
	"github.com/lawrencewoodman/dlit"
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
		map[string]*fieldDescription{
			"band": &fieldDescription{ftString, nil, nil, 0,
				map[string]valueDescription{
					"a": valueDescription{dlit.MustNew("a"), 2},
					"b": valueDescription{dlit.MustNew("b"), 3},
					"c": valueDescription{dlit.MustNew("c"), 1},
					"f": valueDescription{dlit.MustNew("f"), 1},
					"g": valueDescription{dlit.MustNew("g"), 1},
					"h": valueDescription{dlit.MustNew("h"), 1},
					"i": valueDescription{dlit.MustNew("i"), 1},
					"j": valueDescription{dlit.MustNew("j"), 1},
					"k": valueDescription{dlit.MustNew("k"), 1},
					"l": valueDescription{dlit.MustNew("l"), 1},
					"m": valueDescription{dlit.MustNew("m"), 1},
					"n": valueDescription{dlit.MustNew("n"), 1},
					"o": valueDescription{dlit.MustNew("o"), 1},
					"p": valueDescription{dlit.MustNew("p"), 1},
					"q": valueDescription{dlit.MustNew("q"), 1},
					"s": valueDescription{dlit.MustNew("s"), 1},
					"t": valueDescription{dlit.MustNew("t"), 1},
					"u": valueDescription{dlit.MustNew("u"), 1},
					"v": valueDescription{dlit.MustNew("v"), 1},
					"x": valueDescription{dlit.MustNew("x"), 1},
					"y": valueDescription{dlit.MustNew("y"), 1},
					"z": valueDescription{dlit.MustNew("z"), 1},
					"1": valueDescription{dlit.MustNew("1"), 1},
					"2": valueDescription{dlit.MustNew("2"), 1},
					"3": valueDescription{dlit.MustNew("3"), 1},
					"4": valueDescription{dlit.MustNew("4"), 1},
					"5": valueDescription{dlit.MustNew("5"), 1},
					"6": valueDescription{dlit.MustNew("6"), 1},
					"7": valueDescription{dlit.MustNew("7"), 1},
					"8": valueDescription{dlit.MustNew("8"), 2},
					"9": valueDescription{dlit.MustNew("9"), 1},
				},
				31,
			},
			"inputA": &fieldDescription{
				ftFloat,
				dlit.MustNew(7),
				dlit.MustNew(15.1),
				1,
				map[string]valueDescription{
					"7":    valueDescription{dlit.MustNew(7), 7},
					"7.3":  valueDescription{dlit.MustNew(7.3), 7},
					"9":    valueDescription{dlit.MustNew(9), 7},
					"14":   valueDescription{dlit.MustNew(14), 7},
					"15.1": valueDescription{dlit.MustNew(15.1), 7},
				},
				5,
			},
			"inputB": &fieldDescription{
				ftFloat,
				dlit.MustNew(2),
				dlit.MustNew(5),
				4,
				map[string]valueDescription{
					"2.6":    valueDescription{dlit.MustNew(2.6), 7},
					"2.8789": valueDescription{dlit.MustNew(2.8789), 1},
					"3":      valueDescription{dlit.MustNew(3), 7},
					"5":      valueDescription{dlit.MustNew(5), 7},
					"2":      valueDescription{dlit.MustNew(2), 7},
					"2.8":    valueDescription{dlit.MustNew(2.8), 6},
				},
				6,
			},
			"version": &fieldDescription{ftString, nil, nil, 0,
				map[string]valueDescription{
					"9.9":   valueDescription{dlit.MustNew("9.9"), 7},
					"9.97":  valueDescription{dlit.MustNew("9.97"), 7},
					"10":    valueDescription{dlit.MustNew("10"), 7},
					"10.94": valueDescription{dlit.MustNew("10.94"), 7},
					"9.9a":  valueDescription{dlit.MustNew("9.9a"), 6},
					"9.9b":  valueDescription{dlit.MustNew("9.9b"), 1},
				},
				6,
			},
			"flow": &fieldDescription{
				ftInt,
				dlit.MustNew(21),
				dlit.MustNew(87),
				0,
				map[string]valueDescription{}, -1},
			"score": &fieldDescription{
				ftInt,
				dlit.MustNew(1),
				dlit.MustNew(5),
				0,
				map[string]valueDescription{
					"1": valueDescription{dlit.MustNew(1), 6},
					"2": valueDescription{dlit.MustNew(2), 7},
					"3": valueDescription{dlit.MustNew(3), 6},
					"4": valueDescription{dlit.MustNew(4), 8},
					"5": valueDescription{dlit.MustNew(5), 8},
				}, 5,
			},
			"method": &fieldDescription{ftIgnore, nil, nil, 0,
				map[string]valueDescription{}, -1},
		}}
	dataset := NewLiteralDataset(fieldNames, flowRecords)
	d, err := DescribeDataset(dataset)
	if err != nil {
		t.Errorf("DescribeDataset(dataset) err: %s", err)
	}
	if err := checkDescriptionsEqual(d, expected); err != nil {
		t.Errorf("DescibeDataset(dataset) got not expected: %s", err)
	}
}

/*************************
 *   Helper functions
 *************************/
func checkDescriptionsEqual(dGot *Description, dWant *Description) error {
	return fieldDescriptionsEqual(dGot.fields, dWant.fields)
}

func fieldDescriptionsEqual(
	fdsGot map[string]*fieldDescription,
	fdsWant map[string]*fieldDescription,
) error {
	for field, fdG := range fdsGot {
		fdW, ok := fdsWant[field]
		if !ok {
			return fmt.Errorf("Field Description missing for field: %s", field)
		}
		if err := fieldDescriptionEqual(fdG, fdW); err != nil {
			return fmt.Errorf("Field Description for field: %s, %s", field, err)
		}
	}
	return nil
}

func fieldDescriptionEqual(
	fdGot *fieldDescription,
	fdWant *fieldDescription,
) error {
	if fdGot.kind != fdWant.kind {
		return fmt.Errorf("got field kind: %s, want: %s", fdGot.kind, fdWant.kind)
	}
	if len(fdGot.values) != len(fdWant.values) {
		return fmt.Errorf("got %d values, want: %d",
			len(fdGot.values), len(fdWant.values))
	}
	if fdGot.kind == ftInt || fdGot.kind == ftFloat {
		if fdGot.min.String() != fdWant.min.String() ||
			fdGot.max.String() != fdWant.max.String() {
			return fmt.Errorf("got min: %s and max: %s, want min: %s and max: %s",
				fdGot.min, fdGot.max, fdWant.min, fdWant.max)
		}
	}
	if fdGot.kind == ftFloat {
		if fdGot.maxDP != fdWant.maxDP {
			return fmt.Errorf("got maxDP: %d, want: %d", fdGot.maxDP, fdWant.maxDP)
		}
	}
	if fdGot.kind == ftFloat {
		if fdGot.maxDP != fdWant.maxDP {
			return fmt.Errorf("got maxDP: %d, want: %d", fdGot.maxDP, fdWant.maxDP)
		}
	}

	if fdGot.numValues != fdWant.numValues {
		return fmt.Errorf("got numValues: %d, numValues: %d",
			fdGot.numValues, fdWant.numValues)
	}

	return fieldValuesEqual(fdGot.values, fdWant.values)
}

func fieldValuesEqual(
	vdsGot map[string]valueDescription,
	vdsWant map[string]valueDescription,
) error {
	if len(vdsGot) != len(vdsWant) {
		return fmt.Errorf("got %d valueDescriptions, want: %d",
			len(vdsGot), len(vdsWant))
	}
	for k, vdW := range vdsWant {
		vdG, ok := vdsGot[k]
		if !ok {
			return fmt.Errorf("valueDescription missing value: %s", k)
		}
		if vdG.num != vdW.num || vdG.value.String() != vdW.value.String() {
			return fmt.Errorf("got valueDescription: %s, want: %s", vdG, vdW)
		}
	}
	return nil
}
