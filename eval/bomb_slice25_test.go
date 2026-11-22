package eval_test

import (
	"testing"
	"time"

	"github.com/go-embedded-r/r/eval"
)

// Bombs of slice 25: a huge order nu or argument x for the Bessel functions
// (R sizes a work array by nu; here nothing is allocated by it), results that
// would not fit, and a psigamma order of 2^31 or more. Each is refused or
// bounded in time under a memory limit, and the session stays usable.
func TestSlice25Bombs(t *testing.T) {
	underMemoryLimit(t)
	small := func() *eval.Session {
		return eval.NewSession(eval.Options{MaxVector: 1 << 17, MaxAlloc: 1 << 28, Timeout: 30 * time.Second, MaxDepth: 200})
	}
	t.Run("huge orders and arguments are bounded", func(t *testing.T) {
		s := small()
		start := time.Now()
		for _, src := range []string{
			`suppressWarnings(besselJ(1, 1e300))`, `suppressWarnings(besselY(1, 1e300))`, `suppressWarnings(besselI(1, 1e300))`,
			`suppressWarnings(besselK(1, 1e300))`, `suppressWarnings(besselJ(1, -1e300))`, `suppressWarnings(besselY(1, -1e300))`,
			`suppressWarnings(besselI(1, -1e300))`, `suppressWarnings(besselK(1, -1e300))`,
			`suppressWarnings(besselJ(1e300, 1))`, `suppressWarnings(besselY(1e300, 1))`, `suppressWarnings(besselI(1e300, 1))`,
			`suppressWarnings(besselK(1e300, 1))`, `suppressWarnings(besselI(1e300, 1, TRUE))`, `suppressWarnings(besselK(1e300, 1, TRUE))`,
			`suppressWarnings(besselJ(1e9, 1e9))`, `suppressWarnings(besselJ(1e8, 1e4))`, `suppressWarnings(besselJ(1e17, 3))`,
			`suppressWarnings(besselI(1e5, 1e8, TRUE))`, `suppressWarnings(besselK(1e5, 1e8, TRUE))`, `suppressWarnings(besselY(1e8, 1e6))`,
			`suppressWarnings(besselJ(1e8, 2^53))`, `suppressWarnings(besselJ(.Machine$double.xmax, .Machine$double.xmax))`,
		} {
			if _, err := s.Eval(src); err != nil {
				t.Errorf("%s: %v", src, err)
			}
		}
		if d := time.Since(start); d > 20*time.Second {
			t.Errorf("huge orders took %v", d)
		}
	})
	t.Run("results that do not fit", func(t *testing.T) {
		s := small()
		const refusal = "cannot allocate vector"
		refused(t, s, `besselJ(1:131073, 1)`, refusal, "too large")
		refused(t, s, `besselK(1, 1:131073)`, refusal, "too large")
		refused(t, s, `beta(1:131073, 2)`, refusal, "too large")
		refused(t, s, `psigamma(1, rep(2, 131073))`, refusal, "too large")
		refused(t, s, `gamma(numeric(2^31))`, refusal, "too large")
		refused(t, s, `besselI(numeric(2^31), 1)`, refusal, "too large")
	})
	t.Run("psigamma order", func(t *testing.T) {
		s := small()
		for _, src := range []string{`psigamma(1, 2^31)`, `psigamma(1, 2^53)`, `psigamma(1, 1e300)`, `psigamma(1, -1e300)`, `psigamma(1, .Machine$double.xmax)`} {
			if _, err := s.Eval(`suppressWarnings(` + src + `)`); err != nil {
				t.Errorf("%s: %v", src, err)
			}
		}
	})
	t.Run("what is allowed still works", func(t *testing.T) {
		s := small()
		res, err := s.Eval(`length(besselJ(1:1e5, 1)); signif(gamma(5), 6)`)
		if err != nil || res.Output != "[1] 100000\n[1] 24\n" {
			t.Fatalf("%v %q", err, res.Output)
		}
	})
}
