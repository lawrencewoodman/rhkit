package rhkit

import (
	"fmt"
	"github.com/lawrencewoodman/dlit"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestDescriptionWriteLoadJSON(t *testing.T) {
	description := &Description{
		map[string]*fieldDescription{
			"band": &fieldDescription{ftString, nil, nil, 0,
				map[string]valueDescription{
					"a": valueDescription{dlit.MustNew("a"), 2},
					"b": valueDescription{dlit.MustNew("b"), 3},
					"c": valueDescription{dlit.MustNew("c"), 70},
					"f": valueDescription{dlit.MustNew("f"), 22},
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
	if err := checkDescriptionsEqual(got, description); err != nil {
		t.Errorf("LoadDescriptionJSON got not expected: %s", err)
	}
}

/*************************
 *   Helper functions
 *************************/
func checkDescriptionsEqual(dGot *Description, dWant *Description) error {
	return fieldDescriptionsEqual(dGot.Fields, dWant.Fields)
}

func fieldDescriptionsEqual(
	fdsGot map[string]*fieldDescription,
	fdsWant map[string]*fieldDescription,
) error {
	if len(fdsGot) != len(fdsWant) {
		return fmt.Errorf(
			"Number of FieldDescriptions doesn't match. got: %d, want: %d\n",
			len(fdsGot), len(fdsWant),
		)
	}
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
	if fdGot.Kind != fdWant.Kind {
		return fmt.Errorf("got field kind: %s, want: %s", fdGot.Kind, fdWant.Kind)
	}
	if len(fdGot.Values) != len(fdWant.Values) {
		return fmt.Errorf("got %d values, want: %d",
			len(fdGot.Values), len(fdWant.Values))
	}
	if fdGot.Kind == ftInt || fdGot.Kind == ftFloat {
		if fdGot.Min.String() != fdWant.Min.String() ||
			fdGot.Max.String() != fdWant.Max.String() {
			return fmt.Errorf("got min: %s and max: %s, want min: %s and max: %s",
				fdGot.Min, fdGot.Max, fdWant.Min, fdWant.Max)
		}
	}
	if fdGot.Kind == ftFloat {
		if fdGot.MaxDP != fdWant.MaxDP {
			return fmt.Errorf("got maxDP: %d, want: %d", fdGot.MaxDP, fdWant.MaxDP)
		}
	}

	if fdGot.NumValues != fdWant.NumValues {
		return fmt.Errorf("got numValues: %d, numValues: %d",
			fdGot.NumValues, fdWant.NumValues)
	}

	return fieldValuesEqual(fdGot.Values, fdWant.Values)
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
		if vdG.Num != vdW.Num || vdG.Value.String() != vdW.Value.String() {
			return fmt.Errorf("got valueDescription: %s, want: %s", vdG, vdW)
		}
	}
	return nil
}
