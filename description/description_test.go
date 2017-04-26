package description

import (
	"github.com/lawrencewoodman/dlit"
	"github.com/vlifesystems/rhkit/internal/fieldtype"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

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
				fieldtype.Float,
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
				fieldtype.Float,
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
				fieldtype.Int,
				dlit.MustNew(21),
				dlit.MustNew(87),
				0,
				map[string]Value{}, -1},
			"score": &Field{
				fieldtype.Int,
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
	got, err := LoadDescriptionJSON(filename)
	if err != nil {
		t.Fatalf("LoadDescriptionJSON: %s", err)
	}
	if err := got.CheckEqual(description); err != nil {
		t.Errorf("LoadDescriptionJSON got not expected: %s", err)
	}
}
