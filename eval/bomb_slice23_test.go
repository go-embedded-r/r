package eval_test

import (
	"testing"
	"time"

	"github.com/go-embedded-r/r/eval"
)

// Bombs of slice 23: lchoose on a huge k or huge vectors, abbreviate of a huge
// string or vector, agrep with a long pattern and enormous caps, and many
// connections for showConnections. Each is refused or answered within bounds,
// under a memory limit, and the session stays usable.
func TestSlice23Bombs(t *testing.T) {
	underMemoryLimit(t)
	small := func() *eval.Session {
		return eval.NewSession(eval.Options{MaxVector: 1 << 17, MaxAlloc: 1 << 28, MaxFileBytes: 1 << 20, Timeout: 30 * time.Second, MaxDepth: 200})
	}
	const refusal = "cannot allocate vector"
	t.Run("lchoose", func(t *testing.T) {
		s := small()
		res, err := s.Eval(`length(lchoose(1e300, 2^40))`)
		if err != nil {
			t.Fatal(err)
		}
		if res.Output != "[1] 1\n" {
			t.Fatalf("got %q", res.Output)
		}
		refused(t, s, `lchoose(numeric(2^40), 2)`, refusal, "too large")
		refused(t, s, `lchoose(1:3, rep(1, 2^40))`, refusal, "too large")
	})
	t.Run("abbreviate", func(t *testing.T) {
		s := small()
		refused(t, s, `abbreviate(rep("abc", 2^40))`, refusal, "too large")
		res, err := s.Eval(`nchar(abbreviate(strrep("abcdefg ", 5000), 10))`)
		if err != nil {
			t.Fatal(err)
		}
		if res.Output == "" {
			t.Fatal("no answer")
		}
	})
	t.Run("agrep", func(t *testing.T) {
		s := small()
		// a long pattern against a long text with caps beyond both lengths is
		// a bounded dynamic programme, answered or refused, never a hang
		res, err := s.Eval(`agrepl(strrep("ab", 200), strrep("ba", 200), max.distance = 1e9)`)
		if err != nil {
			t.Fatal(err)
		}
		if res.Output != "[1] TRUE\n" {
			t.Fatalf("got %q", res.Output)
		}
		refused(t, s, `agrepl("a", rep("b", 2^40))`, refusal, "too large")
	})
	t.Run("showConnections", func(t *testing.T) {
		s := small()
		res, err := s.Eval(`for (i in 1:50) assign(paste0("c", i), textConnection("a")); n <- nrow(showConnections()); closeAllConnections(); c(n, nrow(showConnections()))`)
		if err != nil {
			t.Fatal(err)
		}
		_ = res
		res, err = s.Eval(`nrow(showConnections())`)
		if err != nil {
			t.Fatal(err)
		}
		if res.Output != "[1] 0\n" {
			t.Fatalf("got %q", res.Output)
		}
	})
}
