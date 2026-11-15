package eval_test

import (
	"testing"
	"time"

	"github.com/go-embedded-r/r/eval"
)

// Bombs of slice 20: the linear algebra builds matrices from the sizes of its
// arguments, and every result is refused against MaxVector before anything is
// allocated (so 2^40 elements, or MaxVector + 1, never reach the allocator),
// under a memory limit, leaving the session usable. The shared depth counter
// bounds the recursion of R code around the natives.
func TestSlice20Bombs(t *testing.T) {
	underMemoryLimit(t)
	small := func() *eval.Session {
		return eval.NewSession(eval.Options{MaxVector: 1 << 17, MaxAlloc: 1 << 28, MaxFileBytes: 1 << 20, Timeout: 30 * time.Second, MaxDepth: 200})
	}
	const refusal = "cannot allocate vector"
	t.Run("a complete Q of a tall matrix", func(t *testing.T) {
		s := small()
		// 100 000 x 100 000 elements
		refused(t, s, `qr.Q(qr(matrix(1, 100000, 1)), complete = TRUE)`, refusal, "too large")
		refused(t, s, `qr.X(qr(matrix(1, 100000, 1)), complete = TRUE, ncol = 100000)`, refusal, "too large")
	})
	t.Run("ncol of 2^40 and MaxVector", func(t *testing.T) {
		s := small()
		refused(t, s, `qr.X(qr(matrix(1:4, 2)), ncol = 2^40)`, refusal, "too large")
		refused(t, s, `qr.X(qr(matrix(1:4, 2)), ncol = 131073)`, refusal, "too large")
		refused(t, s, `qr.X(qr(matrix(1:4, 2)), ncol = -1)`, "invalid 'ncol'")
		refused(t, s, `chol2inv(matrix(1, 2, 2), size = 2^40)`, "size")
		refused(t, s, `backsolve(matrix(1, 2, 2), 1:2, k = 2^40)`, "invalid 'k'")
	})
	t.Run("a square result from a vector", func(t *testing.T) {
		s := small()
		refused(t, s, `svd(matrix(1, 100000, 1), nu = 100000)`, refusal, "too large")
		refused(t, s, `svd(matrix(1, 1, 100000), nv = 100000)`, refusal, "too large")
		refused(t, s, `La.svd(matrix(1, 100000, 1), nu = 100000)`, refusal, "too large")
		refused(t, s, `solve(matrix(1, 1, 1), 1:2)`, "must be compatible")
	})
	t.Run("a huge number of levels", func(t *testing.T) {
		s := small()
		for _, src := range []string{
			`contr.treatment(2^40)`,
			`contr.treatment(100000)`,
			`contr.sum(100000)`,
			`contr.helmert(100000)`,
			`contr.poly(100000)`,
			`contr.poly(2^40)`,
			`contr.SAS(100000)`,
		} {
			refused(t, s, src, refusal, "too large", "limit", "memory", "invalid")
		}
	})
	t.Run("model.matrix with many levels", func(t *testing.T) {
		s := small()
		if _, err := s.Eval(`d <- data.frame(f = factor(seq_len(100000)))`); err != nil {
			t.Fatal(err)
		}
		refused(t, s, `model.matrix(~ f, d)`, refusal, "too large", "limit", "memory")
		refused(t, s, `model.matrix(~ 0 + f, d)`, refusal, "too large", "limit", "memory")
	})
	t.Run("the depth counter", func(t *testing.T) {
		s := small()
		refused(t, s, `f <- function(n) if (n == 0) solve(matrix(1)) else f(n - 1); f(1e5)`, "nested too deeply", "C stack", "evaluation")
		refused(t, s, `g <- function(n) if (n == 0) det(matrix(1)) else 1 + g(n - 1); g(1e6)`, "nested too deeply", "C stack", "evaluation")
	})
	t.Run("what is allowed still works", func(t *testing.T) {
		s := small()
		res, err := s.Eval(`solve(matrix(c(2, 0, 0, 2), 2))[1, 1]`)
		if err != nil {
			t.Fatal(err)
		}
		if got := res.Output; got == "" {
			t.Fatal("no output")
		}
	})
}
