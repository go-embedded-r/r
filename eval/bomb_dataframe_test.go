package eval_test

import (
	"testing"

	"github.com/go-embedded-r/r/eval"
)

// Bombs of data frames and tables: every size is checked before the
// allocation, whether it comes from a recycled column, a cross product or the
// number of cells of a table.
func TestDataFrameBombs(t *testing.T) {
	underMemoryLimit(t)
	s := eval.NewSession(eval.Options{MaxAlloc: 1 << 30})
	for _, src := range []string{
		"data.frame(x = 1:1e12)",
		"data.frame(x = 1:2^31)",
		"data.frame(a = numeric(2^28), b = 1:3)",
		"as.data.frame(matrix(0, 1e6, 1e6))",
		"expand.grid(1:1e6, 1:1e6)",
		"expand.grid(a = 1:2^20, b = 1:2^20, c = 1:2^20)",
		"merge(data.frame(a = 1:1e5), data.frame(b = 1:1e5))",
		"merge(data.frame(k = rep(1, 1e5), a = 1), data.frame(k = rep(1, 1e5), b = 1))",
		"table(1:2^28)",
		"table(1:5000, 1:5000, 1:5000)",
		"df <- data.frame(a = 1); df[1e12, ] <- 1",
		"df <- data.frame(a = 1); df[2^31, \"a\"] <- 1",
		"cut(1:3, 1e12)",
		"cbind(numeric(2^28), numeric(2^28), numeric(2^28), numeric(2^28))",
		"rbind(numeric(2^28), numeric(2^28), numeric(2^28), numeric(2^28))",
		"interaction(1:1e5, 1:1e5)",
	} {
		refused(t, s, src, "allocate", "too large", "limit", "invalid", "memory", "too big", "cross product")
	}
}
