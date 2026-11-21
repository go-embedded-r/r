package eval_test

import (
	"testing"
	"time"

	"github.com/go-embedded-r/r/eval"
)

// Bombs of slice 21: read.fwf with huge widths and many fields, ecdf, poly and
// ftable asked for 2^40 elements or MaxVector + 1, and nesting of the R code
// around the new natives. Each is refused before anything is allocated, under
// a memory limit, and the session stays usable.
func TestSlice21Bombs(t *testing.T) {
	underMemoryLimit(t)
	small := func() *eval.Session {
		return eval.NewSession(eval.Options{MaxVector: 1 << 17, MaxAlloc: 1 << 28, MaxFileBytes: 1 << 20, Timeout: 30 * time.Second, MaxDepth: 200})
	}
	const refusal = "cannot allocate vector"
	t.Run("read.fwf with huge widths", func(t *testing.T) {
		s := small()
		refused(t, s, `read.fwf(textConnection("abc"), widths = rep(1, 2^40))`, refusal, "too large", "limit", "invalid")
		refused(t, s, `read.fwf(textConnection("abc"), widths = rep(1, 131073))`, refusal, "too large", "limit", "invalid")
		refused(t, s, `read.fwf(textConnection("abc"), widths = numeric(2^40))`, refusal, "too large", "limit", "invalid")
		refused(t, s, `read.fwf(textConnection("abc"), widths = list(1, 2))`, "not supported")
		refused(t, s, `read.fwf(textConnection("abc"), widths = integer())`, "invalid", "widths", "argument", "no ")
		// a width beyond any line reads NA, without building a field of that width
		res, err := s.Eval(`dim(read.fwf(textConnection(c("abc", "de")), widths = c(2, 2^40, 1)))`)
		if err != nil {
			t.Fatal(err)
		}
		if res.Output != "[1] 2 3\n" {
			t.Fatalf("got %q", res.Output)
		}
	})
	t.Run("ecdf", func(t *testing.T) {
		s := small()
		refused(t, s, `ecdf(numeric(2^40))`, refusal, "too large")
		refused(t, s, `ecdf(rep(1, 131073))`, refusal, "too large")
		refused(t, s, `ecdf(numeric())`, "1 or more")
		refused(t, s, `knots(ecdf(1))(2^40)`, "attempt to apply", "invalid", "could not find", "not a function", "argument")
		refused(t, s, `ecdf(1:3)(numeric(2^40))`, refusal, "too large")
	})
	t.Run("poly", func(t *testing.T) {
		s := small()
		refused(t, s, `poly(1:10, degree = 2^40)`, "must be less than", refusal, "too large")
		refused(t, s, `poly(1:10, degree = 2^40, raw = TRUE)`, refusal, "too large")
		refused(t, s, `poly(numeric(2^40), 2)`, refusal, "too large")
		refused(t, s, `poly(1:5, 1:5, degree = 2^40)`, refusal, "too large", "must be less than", "invalid")
		refused(t, s, `poly(1:20, 1:20, 1:20, 1:20, 1:20, 1:20, degree = 18)`, refusal, "too large", "limit", "memory")
		refused(t, s, `poly(1:10, degree = 0)`, "must be at least 1")
		refused(t, s, `poly(c(1, NA, 3), 2)`, "missing values")
		refused(t, s, `x <- poly(1:10, 2); predict(x, numeric(2^40))`, refusal, "too large")
	})
	t.Run("ftable", func(t *testing.T) {
		s := small()
		refused(t, s, `ftable(array(0, c(2^20, 2^20)))`, refusal, "too large")
		refused(t, s, `ftable(table(1:100000, 1:100000))`, refusal, "too large", "limit", "memory")
		refused(t, s, `ftable(array(0, c(512, 512, 2)), row.vars = 2^40)`, "invalid", "row.vars", "subscript", "out of bounds", "too large")
		refused(t, s, `ftable(1:3, 1:3, row.vars = 5)`, "invalid", "row.vars", "subscript", "out of bounds", "must", "incorrect")
	})
	t.Run("the other natives", func(t *testing.T) {
		s := small()
		refused(t, s, `gregexec("a", rep("a", 2^40))`, refusal, "too large")
		refused(t, s, `formatDL(rep("a", 2^40), "b")`, refusal, "too large", "incorrect", "same")
		refused(t, s, `.rowSums(1, 2^40, 2)`, "too short", "invalid", "too large", refusal, "length")
		refused(t, s, `.colMeans(1, 2, 2^40)`, "too short", "invalid", "too large", refusal, "length")
		refused(t, s, `.row(c(2^40, 2^40))`, refusal, "too large", "invalid")
		refused(t, s, `.col(c(2^20, 2^20))`, refusal, "too large", "invalid")
	})
	t.Run("the depth counter", func(t *testing.T) {
		s := small()
		refused(t, s, `f <- function(n) if (n == 0) ecdf(1) else f(n - 1); f(1e5)`, "nested too deeply", "C stack", "evaluation")
		refused(t, s, `g <- function(n) if (n == 0) 1 else withRestarts(g(n - 1), r = function() 1); g(1e5)`, "nested too deeply", "C stack", "evaluation")
		refused(t, s, `h <- function(n) if (n == 0) 1 else withAutoprint(h(n - 1)); h(1e5)`, "nested too deeply", "C stack", "evaluation")
		refused(t, s, `k <- function(n) if (n == 0) pairlist(1) else pairlist(k(n - 1)); k(1e5)`, "nested too deeply", "C stack", "evaluation")
	})
	t.Run("what is allowed still works", func(t *testing.T) {
		s := small()
		res, err := s.Eval(`ecdf(c(1, 2, 2, 3))(2)`)
		if err != nil {
			t.Fatal(err)
		}
		if res.Output != "[1] 0.75\n" {
			t.Fatalf("got %q", res.Output)
		}
	})
}
