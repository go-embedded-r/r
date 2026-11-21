package eval_test

import (
	"testing"
	"time"

	"github.com/go-embedded-r/r/eval"
)

// Bombs of slice 22: numToBits and numToInts asked for 64 or 2 elements per
// number past the vector limit, encodeString padded to a huge width, choose on
// huge arguments, findInterval and charmatch on huge inputs, and the nesting of
// the R code around the new closures. Each is refused (or answered) without
// building the huge object, under a memory limit, and the session stays usable.
func TestSlice22Bombs(t *testing.T) {
	underMemoryLimit(t)
	small := func() *eval.Session {
		return eval.NewSession(eval.Options{MaxVector: 1 << 17, MaxAlloc: 1 << 28, MaxFileBytes: 1 << 20, Timeout: 30 * time.Second, MaxDepth: 200})
	}
	const refusal = "cannot allocate vector"
	t.Run("numToBits and numToInts", func(t *testing.T) {
		s := small()
		refused(t, s, `numToBits(numeric(2^40))`, refusal, "too large")
		refused(t, s, `numToBits(numeric(2049))`, refusal, "too large")
		refused(t, s, `numToInts(numeric(65537))`, refusal, "too large")
		res, err := s.Eval(`length(numToBits(numeric(2048)))`)
		if err != nil {
			t.Fatal(err)
		}
		if res.Output != "[1] 131072\n" {
			t.Fatalf("got %q", res.Output)
		}
	})
	t.Run("encodeString", func(t *testing.T) {
		s := small()
		refused(t, s, `encodeString("a", width = 2^31 - 1)`, refusal, "too large")
		refused(t, s, `encodeString(rep("a", 1000), width = 1000)`, refusal, "too large")
		refused(t, s, `encodeString(rep("a", 2^40))`, refusal, "too large")
	})
	t.Run("choose", func(t *testing.T) {
		s := small()
		// a huge k is a product of at most a thousand terms or one gamma ratio
		res, err := s.Eval(`length(choose(1e300, 2^40))`)
		if err != nil {
			t.Fatal(err)
		}
		if res.Output != "[1] 1\n" {
			t.Fatalf("got %q", res.Output)
		}
		refused(t, s, `choose(numeric(2^40), 2)`, refusal, "too large")
		refused(t, s, `choose(1:3, rep(1, 2^40))`, refusal, "too large")
	})
	t.Run("findInterval and charmatch", func(t *testing.T) {
		s := small()
		refused(t, s, `findInterval(numeric(2^40), 1:3)`, refusal, "too large")
		refused(t, s, `findInterval(1, numeric(2^40))`, refusal, "too large")
		refused(t, s, `charmatch(rep("a", 2^40), "a")`, refusal, "too large")
		refused(t, s, `char.expand(rep("a", 2^40), "a")`, refusal, "too large")
	})
	t.Run("the depth counter", func(t *testing.T) {
		s := small()
		refused(t, s, `f <- function(n) if (n == 0) 1 else callCC(function(k) f(n - 1)); f(1e5)`, "nested too deeply", "C stack", "evaluation")
		refused(t, s, `g <- function(n) if (n == 0) 1 else rapply(list(g(n - 1)), identity); g(1e5)`, "nested too deeply", "C stack", "evaluation")
		refused(t, s, `h <- function(n) if (n == 0) 1 else .mapply(function(x) h(n - 1), list(1), NULL); h(1e5)`, "nested too deeply", "C stack", "evaluation")
	})
	t.Run("what is allowed still works", func(t *testing.T) {
		s := small()
		res, err := s.Eval(`findInterval(2.5, 1:3)`)
		if err != nil {
			t.Fatal(err)
		}
		if res.Output != "[1] 2\n" {
			t.Fatalf("got %q", res.Output)
		}
	})
}
