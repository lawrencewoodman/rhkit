/*************************
 *  Test helper functions
 *************************/
package rhkit

import (
	"github.com/lawrencewoodman/dlit"
)

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

func makeStringsDlitSlice(strings ...string) []*dlit.Literal {
	r := make([]*dlit.Literal, len(strings))
	for i, s := range strings {
		r[i] = dlit.NewString(s)
	}
	return r
}
