package eval_test

import (
	"testing"

	"github.com/go-embedded-r/r/eval"
)

// Bombs of slice 15: every size that comes from an argument (a count, a
// repeat, the dimensions of an array, a width) is checked before the
// allocation, whatever the closure that receives it.
func TestSlice15Bombs(t *testing.T) {
	underMemoryLimit(t)
	s := eval.NewSession(eval.Options{MaxAlloc: 1 << 30})
	for _, src := range []string{
		"gl(2, 1e12)",
		"gl(1e6, 1e6)",
		"gl(2^31, 2)",
		"slice.index(array(0, c(2^20, 2^20, 2^20)), 1)",
		"slice.index(1:2, 2^31)",
		"arrayInd(seq_len(2^28), c(2^30, 2))",
		"gl(2^20, 2^12)",
		"gl(3, 2^31)",
		"rep(1:2^20, each = 2^20)",
		"x <- list(); for (i in 1:100000) x <- list(x); all.equal(x, x)",
		"strtrim(strrep(\"a\", 2^30), 2^31 - 1)",
		"as.function(alist(a = , a))(numeric(2^40))",
		"formals(function(a) a) <- rep(alist(x = ), 2^28)",
		"diff(numeric(2^28), lag = 1, differences = 2^30)",
		"format.pval(numeric(2^30))",
		"summary(numeric(2^40))",
		"summary(data.frame(a = 1:2^40))",
		"mean(numeric(2^40), trim = 0.1)",
		"all.equal(numeric(2^40), 1)",
		"all.equal(1:2^31, numeric(2^31))",
		"df <- data.frame(a = 1:3); within(df, b <- numeric(2^40))",
		"df <- data.frame(a = 1:3); df * numeric(2^40)",
		"toString(numeric(2^40))",
		"as.hexmode(numeric(2^40))",
		"format(as.hexmode(1L), width = 2^31 - 1)",
		"format(as.octmode(1L), width = 2^40)",
		"write(numeric(2^40), tempfile())",
		"length(x <- structure(1:3, class = \"factor\")) <- 2^40",
		"x <- as.Date(\"2024-01-01\"); length(x) <- 2^40",
		"x <- factor(\"a\"); length(x) <- 2^40",
		"labels(array(0, c(2^20, 2^20, 2^20)))",
		"labels(numeric(2^40))",
		"isSymmetric(matrix(0, 2^20, 2^20))",
		"as.array(numeric(2^40))",
		"alist(a = , b = )[rep(1, 2^40)]",
		"pmax.int(numeric(2^40), 1)",
	} {
		t.Run(src, func(t *testing.T) {
			refused(t, s, src, "allocate", "too large", "limit", "invalid", "memory", "too big", "cannot", "not allowed", "must", "overflow", "negative", "NA", "dims", "nested too deeply")
		})
	}
}
