package eval_test

import (
	"testing"

	"github.com/go-embedded-r/r/eval"
)

// Bombs of slice 16: a count, a width or a range that comes from an argument
// is checked before anything is allocated, and a recursion over a value is
// bounded by the shared depth counter.
func TestSlice16Bombs(t *testing.T) {
	underMemoryLimit(t)
	s := eval.NewSession(eval.Options{MaxAlloc: 1 << 30})
	for _, src := range []string{
		"seq(as.POSIXct(0, tz = \"UTC\"), by = \"sec\", length.out = 2^40)",
		"seq(as.POSIXct(0, tz = \"UTC\"), by = \"month\", length.out = 2^40)",
		"seq(as.POSIXct(0, tz = \"UTC\"), as.POSIXct(2^40, tz = \"UTC\"), by = \"sec\")",
		"seq(as.POSIXct(0, tz = \"UTC\"), as.POSIXct(2^40, tz = \"UTC\"), by = \"DSTday\")",
		"cut(as.POSIXct(c(0, 2^40), tz = \"UTC\"), \"secs\")",
		"cut(as.Date(c(0, 2^40)), \"day\", labels = FALSE)",
		"cut(as.Date(c(\"2024-01-01\", \"2024-02-01\")), 2^40)",
		"strwrap(strrep(\"a \", 2^40), 10)",
		"strwrap(\"a b\", width = 2^40, indent = 2^40)",
		"kronecker(numeric(2^20), numeric(2^20))",
		"matrix(0, 2^20, 2^10) %x% matrix(0, 2^10, 2^10)",
		"pretty(c(0, 1), n = 2^40)",
		"pretty(c(0, 1e300), n = 1e9)",
		"format(numeric(2^40), big.mark = \",\")",
		"formatC(numeric(2^40), big.mark = \",\")",
		"prettyNum(numeric(2^40), drop0trailing = TRUE)",
		"x <- 0; for (i in 1:100000) x <- list(x); deparse(x)",
		"x <- 0; for (i in 1:100000) x <- list(x); dput(x)",
		"x <- 0; for (i in 1:100000) x <- list(x); deparse(x, control = \"all\")",
		"x <- 0; for (i in 1:100000) x <- structure(list(x), class = \"a\"); deparse(x)",
		"deparse(1:2^30, nlines = 1)",
		"deparse(numeric(2^40))",
		"dump(\"x\", file = tempfile(), envir = list2env(list(x = numeric(2^40))))",
		"by(data.frame(a = 1:2), 1:2^40, nrow)",
		"unique(matrix(0, 2^20, 2^20))",
		"duplicated(array(0, c(2^20, 2^20, 2^20)), MARGIN = 3)",
		"as.table(array(0, c(2^20, 2^20, 2^20)))",
		"provideDimnames(array(0, c(2^20, 2^20, 2^20)))",
		"data.matrix(data.frame(a = 1:2^40))",
		"list2DF(list(a = 1:2^31), nrow = 2^31)",
		"single(2^40)",
		"numeric_version(strrep(\"1.\", 2^30))",
		"v <- numeric_version(rep(\"1.2\", 2^28)); v[2^30]",
		"sort(numeric_version(rep(\"1.2\", 2^28)))",
		"rep(numeric_version(\"1.2\"), 2^40)",
		"file.info(rep(\"a\", 2^28))",
		"file.info(strrep(\"a\", 2^30))",
		"print.rle(rle(numeric(2^40)))",
		"summary(as.Date(numeric(2^40)))",
		"trunc(as.POSIXct(numeric(2^40), tz = \"UTC\"), \"months\")",
		"round(as.POSIXct(numeric(2^40), tz = \"UTC\"), \"years\")",
		"mode(x) <- \"list\"; x <- numeric(2^40)",
		"simplify2array(rep(list(1:2^20), 2^20))",
		"sort_by(numeric(2^40), 1)",
	} {
		t.Run(src, func(t *testing.T) {
			refused(t, s, src, "allocate", "too large", "limit", "invalid", "memory", "too big", "cannot", "not allowed", "must", "overflow", "negative", "NA", "dims", "nested too deeply", "too many")
		})
	}
}
