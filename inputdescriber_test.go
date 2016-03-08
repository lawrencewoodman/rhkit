package main

import (
	"github.com/lawrencewoodman/dlit_go"
	"path/filepath"
	"testing"
)

func TestDescribeInput(t *testing.T) {
	filename := filepath.Join("fixtures", "flow.csv")
	skipFirstLine := true
	fieldNames :=
		[]string{"band", "inputA", "inputB", "version", "flow", "score", "method"}
	expected := map[string]*FieldDescription{
		"band": &FieldDescription{STRING, nil, nil, 0,
			[]*dlit.Literal{mustNewLit("a"), mustNewLit("b"),
				mustNewLit("c"), mustNewLit("f"), mustNewLit("g"),
				mustNewLit("i"), mustNewLit("j"), mustNewLit("k"),
				mustNewLit("l"), mustNewLit("m"), mustNewLit("n"),
				mustNewLit("o"), mustNewLit("p"), mustNewLit("q"),
				mustNewLit("s"), mustNewLit("t"), mustNewLit("u"),
				mustNewLit("v"), mustNewLit("x"), mustNewLit("y"),
				mustNewLit("z"), mustNewLit("1"), mustNewLit("2"),
				mustNewLit("3"), mustNewLit("4"), mustNewLit("5"),
				mustNewLit("6"), mustNewLit("7"), mustNewLit("8"),
				mustNewLit("9"), mustNewLit("h")}, 0},
		"inputA": &FieldDescription{FLOAT, mustNewLit(7), mustNewLit(15.1), 1,
			[]*dlit.Literal{mustNewLit(7), mustNewLit(7.3),
				mustNewLit(9), mustNewLit(14), mustNewLit(15.1)}, 0},
		"inputB": &FieldDescription{FLOAT, mustNewLit(2), mustNewLit(5), 4,
			[]*dlit.Literal{mustNewLit(2.6), mustNewLit(2.8789),
				mustNewLit(3), mustNewLit(5), mustNewLit(2),
				mustNewLit(2.8)}, 0},
		"version": &FieldDescription{STRING, nil, nil, 0,
			[]*dlit.Literal{mustNewLit("9.9"), mustNewLit("9.97"),
				mustNewLit("10"), mustNewLit("10.94"), mustNewLit("9.9a"),
				mustNewLit("9.9b")}, 0},
		"flow": &FieldDescription{INT, mustNewLit(21), mustNewLit(87), 0,
			[]*dlit.Literal{}, 0},
		"score": &FieldDescription{INT, mustNewLit(1), mustNewLit(5), 0,
			[]*dlit.Literal{mustNewLit(1), mustNewLit(2), mustNewLit(3),
				mustNewLit(4), mustNewLit(5)}, 0},
		"method": &FieldDescription{IGNORE, nil, nil, 0,
			[]*dlit.Literal{}, 0},
	}
	input, err := NewCsvInput(fieldNames, filename, ',', skipFirstLine)
	if err != nil {
		t.Errorf("NewCsvInput() - err: %q", filename, err)
	}
	fd, err := DescribeInput(input)
	if err != nil {
		t.Errorf("DescribeInput(input) err: %s", err)
	}
	if !fieldDescriptionsEqual(fd, expected) {
		t.Errorf("fieldDescriptionsEqual(%q, %q) not equal", fd, expected)
	}
}

/*************************
 *   Helper functions
 *************************/
func fieldDescriptionsEqual(
	fds1 map[string]*FieldDescription, fds2 map[string]*FieldDescription) bool {
	for field, fd1 := range fds1 {
		fd2, ok := fds2[field]
		if ok && !fieldDescriptionEqual(fd1, fd2) {
			return false
		}
	}
	return true
}
func fieldDescriptionEqual(fd1 *FieldDescription, fd2 *FieldDescription) bool {
	if fd1.Kind != fd2.Kind || len(fd1.Values) != len(fd2.Values) {
		return false
	}
	if fd1.Kind == INT || fd1.Kind == FLOAT {
		if fd1.Min.String() != fd2.Min.String() ||
			fd1.Max.String() != fd2.Max.String() {
			return false
		}
	}
	if fd1.Kind == FLOAT {
		if fd1.MaxDP != fd2.MaxDP {
			return false
		}
	}
	for _, fd1V := range fd1.Values {
		found := false
		for _, fd2V := range fd2.Values {
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
