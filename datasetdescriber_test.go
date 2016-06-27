package rulehunter

import (
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
				[]*dlit.Literal{dlit.MustNew("a"), dlit.MustNew("b"),
					dlit.MustNew("c"), dlit.MustNew("f"), dlit.MustNew("g"),
					dlit.MustNew("i"), dlit.MustNew("j"), dlit.MustNew("k"),
					dlit.MustNew("l"), dlit.MustNew("m"), dlit.MustNew("n"),
					dlit.MustNew("o"), dlit.MustNew("p"), dlit.MustNew("q"),
					dlit.MustNew("s"), dlit.MustNew("t"), dlit.MustNew("u"),
					dlit.MustNew("v"), dlit.MustNew("x"), dlit.MustNew("y"),
					dlit.MustNew("z"), dlit.MustNew("1"), dlit.MustNew("2"),
					dlit.MustNew("3"), dlit.MustNew("4"), dlit.MustNew("5"),
					dlit.MustNew("6"), dlit.MustNew("7"), dlit.MustNew("8"),
					dlit.MustNew("9"), dlit.MustNew("h")}, 0},
			"inputA": &fieldDescription{
				ftFloat,
				dlit.MustNew(7),
				dlit.MustNew(15.1),
				1,
				[]*dlit.Literal{dlit.MustNew(7), dlit.MustNew(7.3),
					dlit.MustNew(9), dlit.MustNew(14), dlit.MustNew(15.1)}, 0},
			"inputB": &fieldDescription{
				ftFloat,
				dlit.MustNew(2),
				dlit.MustNew(5),
				4,
				[]*dlit.Literal{dlit.MustNew(2.6), dlit.MustNew(2.8789),
					dlit.MustNew(3), dlit.MustNew(5), dlit.MustNew(2),
					dlit.MustNew(2.8)}, 0},
			"version": &fieldDescription{ftString, nil, nil, 0,
				[]*dlit.Literal{dlit.MustNew("9.9"), dlit.MustNew("9.97"),
					dlit.MustNew("10"), dlit.MustNew("10.94"), dlit.MustNew("9.9a"),
					dlit.MustNew("9.9b")}, 0},
			"flow": &fieldDescription{
				ftInt,
				dlit.MustNew(21),
				dlit.MustNew(87),
				0,
				[]*dlit.Literal{}, 0},
			"score": &fieldDescription{
				ftInt,
				dlit.MustNew(1),
				dlit.MustNew(5),
				0,
				[]*dlit.Literal{dlit.MustNew(1), dlit.MustNew(2), dlit.MustNew(3),
					dlit.MustNew(4), dlit.MustNew(5)}, 0},
			"method": &fieldDescription{ftIgnore, nil, nil, 0,
				[]*dlit.Literal{}, 0},
		}}
	dataset := NewLiteralDataset(fieldNames, flowRecords)
	d, err := DescribeDataset(dataset)
	if err != nil {
		t.Errorf("DescribeDataset(dataset) err: %s", err)
	}
	if !descriptionsEqual(d, expected) {
		t.Errorf("DescibeDataset(dataset) got: %s, want: %s", d, expected)
	}
}

/*************************
 *   Helper functions
 *************************/
func descriptionsEqual(d1 *Description, d2 *Description) bool {
	return fieldDescriptionsEqual(d1.fields, d2.fields)
}

func fieldDescriptionsEqual(
	fds1 map[string]*fieldDescription,
	fds2 map[string]*fieldDescription,
) bool {
	for field, fd1 := range fds1 {
		fd2, ok := fds2[field]
		if ok && !fieldDescriptionEqual(fd1, fd2) {
			return false
		}
	}
	return true
}

func fieldDescriptionEqual(
	fd1 *fieldDescription,
	fd2 *fieldDescription,
) bool {
	if fd1.kind != fd2.kind || len(fd1.values) != len(fd2.values) {
		return false
	}
	if fd1.kind == ftInt || fd1.kind == ftFloat {
		if fd1.min.String() != fd2.min.String() ||
			fd1.max.String() != fd2.max.String() {
			return false
		}
	}
	if fd1.kind == ftFloat {
		if fd1.maxDP != fd2.maxDP {
			return false
		}
	}
	for _, fd1V := range fd1.values {
		found := false
		for _, fd2V := range fd2.values {
			if fd1V.String() == fd2V.String() {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}
