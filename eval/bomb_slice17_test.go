package eval_test

import (
	"strings"
	"testing"

	"github.com/go-embedded-r/r/eval"
)

// Bombs of slice 17: a formula that expands to 2^20 terms, one nested 100 000
// deep, a table of groups, a structure for str, a model frame: each is refused
// before it is built, or through the shared depth counter.
func TestSlice17Bombs(t *testing.T) {
	underMemoryLimit(t)
	s := eval.NewSession(eval.Options{MaxAlloc: 1 << 30})
	vars := `paste0("v", 1:20)`
	for _, src := range []string{
		"terms(as.formula(paste('y ~ (', paste(" + vars + ", collapse = '+'), ')^20')))",
		"terms(as.formula(paste('y ~', paste(" + vars + ", collapse = '*'))))",
		"terms(y ~ a^1e9)",
		"terms(reformulate(paste0('v', 1:100000)))",
		"terms(reformulate(paste0('v', 1:100000)), keep.order = TRUE)",
		"f <- y ~ x; for (i in 1:100000) f[[3]] <- call('+', f[[3]], quote(z)); terms(f)",
		"f <- y ~ x; for (i in 1:100000) f[[3]] <- call('+', f[[3]], quote(z)); all.vars(f)",
		"f <- y ~ x; for (i in 1:100000) f[[3]] <- call('+', f[[3]], quote(z)); update(f, . ~ .)",
		"f <- y ~ x; for (i in 1:100000) f[[3]] <- call('+', f[[3]], quote(z)); deparse(f)",
		"model.frame(y ~ x, data.frame(y = numeric(2^40), x = 1))",
		"model.frame(y ~ x, list(y = numeric(2^40), x = 1))",
		"aggregate(numeric(10), rep(list(1:10), 8), sum, drop = FALSE)",
		"aggregate(numeric(2^40), list(1), sum)",
		"aggregate(y ~ g, data.frame(y = numeric(2^40), g = 1), sum)",
		"stack(list(numeric(2^40)))",
		"unstack(data.frame(values = numeric(2^40), ind = 1))",
		"rowsum(numeric(2^40), 1)",
		"rowsum(matrix(0, 2^20, 2^20), 1:2^20)",
		"array2DF(array(0, c(2^20, 2^20, 2^20)))",
		"x <- 0; for (i in 1:100000) x <- list(x); str(x)",
		"x <- 0; for (i in 1:100000) x <- list(x); str(list(x), max.level = 2^30)",
		"Vectorize(function(x) x)(numeric(2^40))",
		"x <- as.POSIXlt(Sys.time()); x[rep(1, 2^40)]",
		"x <- as.POSIXlt(Sys.time()); x$sec <- numeric(2^40)",
		"as.data.frame(array(0, c(2^20, 2^20, 2^20)))",
		"sort_by(data.frame(a = numeric(2^40)), ~ a)",
		"get_all_vars(y ~ x, list(y = numeric(2^40), x = 1))",
		"na.omit(numeric(2^40))",
		"complete.cases(numeric(2^40))",
		"chkDots(numeric(2^40))",
	} {
		t.Run(src, func(t *testing.T) {
			refused(t, s, src, "allocate", "too large", "limit", "invalid", "memory", "too big", "cannot", "not allowed", "must", "overflow", "negative", "NA", "dims", "nested too deeply", "too many", "expands", "more than", "power")
		})
	}
}

// str of a list, a data frame and a named vector: the oracle cannot print
// them (its guard refuses), so the layouts of the help page are checked here.
func TestStrStructures(t *testing.T) {
	cases := []struct{ src, want string }{
		{`str(list(1, "a"))`, "List of 2\n $ : num 1\n $ : chr \"a\"\n"},
		{`str(list(a = 1:3, b = list(c = "x")))`, "List of 2\n $ a: int [1:3] 1 2 3\n $ b:List of 1\n  ..$ c: chr \"x\"\n"},
		{`str(data.frame(a = 1:2, b = c("x", "y")))`, "'data.frame':\t2 obs. of  2 variables:\n $ a: int  1 2\n $ b: chr  \"x\" \"y\"\n"},
		{`str(data.frame(x = 1))`, "'data.frame':\t1 obs. of  1 variable:\n $ x: num 1\n"},
		{`str(c(a = 1, b = 2))`, " Named num [1:2] 1 2\n - attr(*, \"names\")= chr [1:2] \"a\" \"b\"\n"},
		{`str(list(a = 1, bb = list(x = 1, y = list(z = 2))))`, "List of 2\n $ a : num 1\n $ bb:List of 2\n  ..$ x: num 1\n  ..$ y:List of 1\n  .. ..$ z: num 2\n"},
		{`str(list(a = NULL))`, "List of 1\n $ a: NULL\n"},
	}
	for _, c := range cases {
		s := eval.NewSession(eval.Options{})
		r, err := s.Eval(c.src)
		if err != nil || r.Output != c.want {
			t.Errorf("%s:\n got %q (%v)\nwant %q", c.src, r.Output, err, c.want)
		}
	}
}

// Vectorize: the oracle refuses to print what it returns (a closure of base),
// so its behaviour is checked here against the help page.
func TestVectorize(t *testing.T) {
	cases := []struct{ src, want string }{
		{`Vectorize(function(a, b) a + b)(1:3, 2)`, "[1] 3 4 5\n"},
		{`f <- Vectorize(function(n, m) paste(n, m), "n"); f(1:2, "z")`, "[1] \"1 z\" \"2 z\"\n"},
		{`Vectorize(rep.int)(1:3, 3:1)`, "[[1]]\n[1] 1 1 1\n\n[[2]]\n[1] 2 2\n\n[[3]]\n[1] 3\n\n"},
		{`Vectorize(function(x) x)(c(a = 1))`, "a \n1 \n"},
		{`Vectorize(function(x) x, USE.NAMES = FALSE)(c(a = 1))`, "[1] 1\n"},
		{`Vectorize(function(x, y) x + y, vectorize.args = "x")(1:2, 10)`, "[1] 11 12\n"},
		{`vf <- Vectorize(function(x, y = 2) x * y, SIMPLIFY = FALSE); vf(1:2)`, "[[1]]\n[1] 2\n\n[[2]]\n[1] 4\n\n"},
	}
	for _, c := range cases {
		s := eval.NewSession(eval.Options{})
		r, err := s.Eval(c.src)
		if err != nil || r.Output != c.want {
			t.Errorf("%s:\n got %q (%v)\nwant %q", c.src, r.Output, err, c.want)
		}
	}
	s := eval.NewSession(eval.Options{})
	if _, err := s.Eval(`Vectorize(function(x) x, "zz")`); err == nil || !strings.Contains(err.Error(), "must specify names of formal arguments") {
		t.Errorf("unknown argument: %v", err)
	}
}
