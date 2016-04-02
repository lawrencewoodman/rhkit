package rulehunter

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
		"inputA": &FieldDescription{FLOAT, dlit.MustNew(7), dlit.MustNew(15.1), 1,
			[]*dlit.Literal{dlit.MustNew(7), dlit.MustNew(7.3),
				dlit.MustNew(9), dlit.MustNew(14), dlit.MustNew(15.1)}, 0},
		"inputB": &FieldDescription{FLOAT, dlit.MustNew(2), dlit.MustNew(5), 4,
			[]*dlit.Literal{dlit.MustNew(2.6), dlit.MustNew(2.8789),
				dlit.MustNew(3), dlit.MustNew(5), dlit.MustNew(2),
				dlit.MustNew(2.8)}, 0},
		"version": &FieldDescription{STRING, nil, nil, 0,
			[]*dlit.Literal{dlit.MustNew("9.9"), dlit.MustNew("9.97"),
				dlit.MustNew("10"), dlit.MustNew("10.94"), dlit.MustNew("9.9a"),
				dlit.MustNew("9.9b")}, 0},
		"flow": &FieldDescription{INT, dlit.MustNew(21), dlit.MustNew(87), 0,
			[]*dlit.Literal{}, 0},
		"score": &FieldDescription{INT, dlit.MustNew(1), dlit.MustNew(5), 0,
			[]*dlit.Literal{dlit.MustNew(1), dlit.MustNew(2), dlit.MustNew(3),
				dlit.MustNew(4), dlit.MustNew(5)}, 0},
		"method": &FieldDescription{IGNORE, nil, nil, 0,
			[]*dlit.Literal{}, 0},
	}
	input, err := newCsvInput(fieldNames, filename, ',', skipFirstLine)
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
