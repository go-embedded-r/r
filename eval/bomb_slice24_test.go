package eval_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/go-embedded-r/r/eval"
)

// Bombs of slice 24: draws and samples asked for 2^31 or 2^40 elements, and a
// shuffle of a population that does not fit. Each is refused before anything
// is allocated, under a memory limit, and the session stays usable.
func TestSlice24Bombs(t *testing.T) {
	underMemoryLimit(t)
	small := func() *eval.Session {
		return eval.NewSession(eval.Options{MaxVector: 1 << 17, MaxAlloc: 1 << 28, Timeout: 30 * time.Second, MaxDepth: 200})
	}
	const refusal = "cannot allocate vector"
	t.Run("runif and rnorm", func(t *testing.T) {
		s := small()
		refused(t, s, `runif(2^31)`, refusal, "too large")
		refused(t, s, `runif(2^40)`, refusal, "too large")
		refused(t, s, `runif(131073)`, refusal, "too large")
		refused(t, s, `rnorm(2^31)`, refusal, "too large")
		refused(t, s, `rnorm(4.5e15)`, refusal, "too large")
		refused(t, s, `rnorm(4.6e15)`, "invalid arguments")
		refused(t, s, `runif(Inf)`, "invalid arguments")
		refused(t, s, `jitter(numeric(2^31))`, refusal, "too large")
	})
	t.Run("sample", func(t *testing.T) {
		s := small()
		refused(t, s, `sample.int(2^31, useHash = FALSE)`, refusal, "too large")
		refused(t, s, `sample.int(2^40, useHash = FALSE)`, refusal, "too large")
		refused(t, s, `sample.int(4e15, 2, FALSE, useHash = FALSE)`, refusal, "too large")
		refused(t, s, `sample.int(131073)`, refusal, "too large")
		refused(t, s, `sample.int(10, 2^31, TRUE)`, refusal, "too large")
		refused(t, s, `sample.int(10, 2^40, TRUE)`, refusal, "too large")
		refused(t, s, `sample.int(2^40, 2^30, useHash = TRUE)`, refusal, "too large")
		refused(t, s, `sample.int(2^31, 2^30)`, refusal, "too large")
		refused(t, s, `sample.int(2^31, 3, prob = rep(1, 2^31))`, refusal, "too large", "incorrect number")
		refused(t, s, `sample.int(2^40, 3, TRUE, prob = 1:5)`, refusal, "too large", "incorrect number")
		refused(t, s, `sample(numeric(2^31))`, refusal, "too large")
		refused(t, s, `sample(2e5, 10, TRUE, prob = rep(1, 2e5))`, refusal, "too large")
	})
	t.Run("what is allowed still works", func(t *testing.T) {
		s := small()
		if strconv.IntSize == 64 { // a population past 2^31 - 1 needs a 64-bit int
			res, err := s.Eval(`set.seed(1); x <- sample.int(2^40, 3); length(x)`)
			if err != nil {
				t.Fatal(err)
			}
			if res.Output != "[1] 3\n" {
				t.Fatalf("got %q", res.Output)
			}
		}
		res, err := s.Eval(`set.seed(1); length(runif(1e5))`)
		if err != nil || res.Output != "[1] 100000\n" {
			t.Fatalf("%v %q", err, res.Output)
		}
	})
}
