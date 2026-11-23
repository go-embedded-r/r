package eval_test

import (
	"testing"
	"time"

	"github.com/go-embedded-r/r/eval"
)

// Bombs of slice 26: extreme parameters and quantile arguments of the d/p/q
// functions (the deep-quantile bisection is bounded in steps), huge recycled
// results (refused under the vector limit), and the session stays usable.
func TestSlice26Bombs(t *testing.T) {
	underMemoryLimit(t)
	small := func() *eval.Session {
		return eval.NewSession(eval.Options{MaxVector: 1 << 17, MaxAlloc: 1 << 28, Timeout: 30 * time.Second, MaxDepth: 200})
	}
	t.Run("extreme parameters are bounded", func(t *testing.T) {
		s := small()
		start := time.Now()
		for _, src := range []string{
			`qgamma(-1e300, 1e300, log.p = TRUE)`, `qgamma(-1000, 1e-300, log.p = TRUE)`, `qbeta(-1e10, 1e-300, 1e300, log.p = TRUE)`,
			`qt(-1e300, 1e-300, log.p = TRUE)`, `qt(-1e5, 1e300, log.p = TRUE)`, `qf(-1e300, 1e300, 1e-300, log.p = TRUE)`,
			`qchisq(-1e300, 1e300, log.p = TRUE)`, `qweibull(-1e300, 1e-300, 1e300, log.p = TRUE)`, `qlnorm(-1e300, 0, 1e300, log.p = TRUE)`,
			`qnorm(-1e300, 0, 1e300, log.p = TRUE)`, `qcauchy(-1e300, log.p = TRUE)`, `qlogis(-1e300, log.p = TRUE)`, `qexp(-1e300, 1e-300, log.p = TRUE)`,
			`pgamma(1e300, 1e300)`, `pbeta(0.5, 1e300, 1e300)`, `pt(1e300, 1e300)`, `pf(1e300, 1e300, 1e300)`,
			`pchisq(1e300, 1e-300)`, `dgamma(1e-300, 1e-300, 1e-300)`, `dbeta(1e-300, 1e-300, 1e-300)`,
			`qgamma(0.5, 1e300)`, `qbeta(0.5, 1e300, 1e-300)`, `qt(0.5, 1e-300)`, `qf(0.5, 1e-300, 1e300)`,
			`qunif(0.5, -1e308, 1e308)`, `qweibull(0.5, 1e300, 1e-300)`,
		} {
			if _, err := s.Eval(`suppressWarnings(` + src + `)`); err != nil {
				t.Errorf("%s: %v", src, err)
			}
		}
		if d := time.Since(start); d > 20*time.Second {
			t.Errorf("extreme parameters took %v", d)
		}
	})
	t.Run("results that do not fit", func(t *testing.T) {
		s := small()
		const refusal = "cannot allocate vector"
		refused(t, s, `pnorm(1:131073)`, refusal, "too large")
		refused(t, s, `qgamma(0.5, 1:131073)`, refusal, "too large")
		refused(t, s, `dbeta(0.5, 2, rep(3, 131073))`, refusal, "too large")
		refused(t, s, `pt(numeric(2^31), 3)`, refusal, "too large")
		refused(t, s, `qnorm(0.5, numeric(2^31))`, refusal, "too large")
	})
	t.Run("what is allowed still works", func(t *testing.T) {
		s := small()
		res, err := s.Eval(`length(pnorm(1:1e5)); signif(qnorm(0.975), 6)`)
		if err != nil || res.Output != "[1] 100000\n[1] 1.95996\n" {
			t.Fatalf("%v %q", err, res.Output)
		}
	})
}
