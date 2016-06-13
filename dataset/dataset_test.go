package dataset

import (
	"errors"
	"testing"
)

func TestCheckFieldNamesValid(t *testing.T) {
	cases := []struct {
		fieldNames []string
		wantErr    error
	}{
		{[]string{"name"}, errors.New("must specify at least two field names")},
		{[]string{"name", "de^pt"}, errors.New("invalid field name: de^pt")},
		{[]string{"name", "dept"}, nil},
	}
	for _, c := range cases {
		err := CheckFieldNamesValid(c.fieldNames)
		if !errorMatch(err, c.wantErr) {
			t.Errorf("CheckFieldNames(%s) err: %s, wantErr: %s",
				c.fieldNames, err, c.wantErr)
		}
	}
}

/*********************
 *  Helper functions
 *********************/
func errorMatch(e1 error, e2 error) bool {
	if e1 == nil && e2 == nil {
		return true
	}
	if e1 == nil || e2 == nil {
		return false
	}
	if e1.Error() == e2.Error() {
		return true
	}
	return false
}
