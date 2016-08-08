package rule

import (
	"github.com/lawrencewoodman/dlit"
	"testing"
)

/**************************
 *  Benchmarks
 **************************/

func BenchmarkLTFFIsTrue(b *testing.B) {
	record := map[string]*dlit.Literal{
		"band":   dlit.MustNew(23),
		"income": dlit.MustNew(1024),
		"cost":   dlit.MustNew(890),
	}
	r := NewLTFF("cost", "income")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = r.IsTrue(record)
	}
}
