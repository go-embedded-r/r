package eval_test

import (
	"testing"

	"github.com/go-embedded-r/r/eval"
)

// The base library is evaluated into every new session: its cost is the cost
// of a session.
func BenchmarkNewSession(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		eval.NewSession(eval.Options{})
	}
}
