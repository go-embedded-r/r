package main

import (
	"fmt"
	"math/rand"
)

// The slice24 mode: random programs over the random streams of slice 24:
// set.seed, runif, rnorm (Inversion), sample and sample.int, jitter, RNGkind,
// RNGversion and .Random.seed. The first line is "# family". The oracle's R
// session is shared across programs, so every program starts by setting the
// kinds it relies on and a seed, and one that changes a kind sets the default
// back before it ends. A value ends in "# libm" when it is compared within
// --ulp units in the last place after being rounded to seven digits (rnorm is
// held to a stated bound at the uniform it draws, not to R's libm).

func fuzzSlice24(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	const reset = `RNGkind("Mersenne-Twister", "Inversion", "Rejection"); `
	seed := fmt.Sprintf("set.seed(%d); ", r.Intn(100000)-1000)
	var family, body string
	switch r.Intn(7) {
	case 0:
		family = "runif"
		body = reset + seed + pick([]string{
			fmt.Sprintf("runif(%d)", r.Intn(12)), fmt.Sprintf("runif(%d, %d, %d)", 1+r.Intn(8), r.Intn(5)-2, 3+r.Intn(9)),
			fmt.Sprintf("runif(%d, c(0, 10), c(1, 20))", 1+r.Intn(8)), "runif(c(5, 5, 5))", "suppressWarnings(runif(3, 2, 1))",
			fmt.Sprintf("x <- runif(%d); c(min(x) >= 0, max(x) < 1)", 1+r.Intn(50)),
		})
	case 1:
		family = "rnorm"
		body = reset + seed + "signif(" + pick([]string{
			fmt.Sprintf("rnorm(%d)", 1+r.Intn(12)), fmt.Sprintf("rnorm(%d, %d, %d)", 1+r.Intn(8), r.Intn(7)-3, r.Intn(5)),
			fmt.Sprintf("rnorm(%d, c(0, 100), c(1, 0.5))", 1+r.Intn(8)), "rnorm(c(5, 5, 5))",
			fmt.Sprintf("mean(rnorm(%d))", 10+r.Intn(500)),
		}) + ", 7) # libm"
	case 2:
		family = "sample"
		body = reset + seed + pick([]string{
			fmt.Sprintf("sample(%d)", 1+r.Intn(30)), fmt.Sprintf("sample(%d, %d)", 5+r.Intn(30), 1+r.Intn(5)),
			fmt.Sprintf("sample(%d, %d, replace = TRUE)", 1+r.Intn(30), 1+r.Intn(20)), "sample(c('a', 'b', 'c', 'd'))",
			fmt.Sprintf("sample(%d, %d, TRUE, prob = c(0.5, 0.25, 0.125, 0.125, 0, 1)[1:%d])", 3+r.Intn(3), 1+r.Intn(12), 3+r.Intn(3)),
			fmt.Sprintf("sample(%d, %d, prob = (1:%d)/%d)", 4+r.Intn(8), 1+r.Intn(3), 4+r.Intn(8), 4+r.Intn(8)),
			fmt.Sprintf("sample.int(%d, %d)", 10+r.Intn(1000), 1+r.Intn(5)), "sample.int(2^31 + 10, 3)", "sample(character(0))",
			fmt.Sprintf("sample(%d, 3, useHash = TRUE)", 100+r.Intn(1000)),
		})
	case 3:
		family = "jitter"
		body = reset + seed + pick([]string{
			fmt.Sprintf("signif(jitter(1:%d), 7)", 2+r.Intn(8)), "signif(jitter(c(1, 1, 2), factor = 2), 7)",
			"signif(jitter(rep(3, 4)), 7)", "signif(jitter(c(1.5, 2.5, 9), amount = 0.5), 7)", "signif(jitter(numeric(0)), 7)",
		}) + " # libm"
	case 4:
		family = "kinds"
		body = pick([]string{
			reset + `x <- RNGkind(sample.kind = "Rounding"); ` + seed + "y <- sample(10); RNGkind(sample.kind = 'Rejection'); list(x, y)",
			reset + `RNGkind(normal.kind = "Box-Muller"); z <- RNGkind(); RNGkind(normal.kind = "Inversion"); z`,
			reset + `RNGkind("Mer", "Inv", "Rej")`,
			reset + `RNGversion("3.5.0"); a <- RNGkind(); RNGversion("4.0.0"); a`,
			reset + `tryCatch(RNGkind("foo"), error = conditionMessage)`,
			reset + `tryCatch(RNGkind(normal.kind = "foo"), error = conditionMessage)`,
		})
	case 5:
		family = "state"
		body = reset + seed + pick([]string{
			".Random.seed[1:5]", "invisible(runif(3)); .Random.seed[2]", "a <- .Random.seed; x <- runif(4); .Random.seed <- a; identical(x, runif(4))",
			"length(.Random.seed)", "invisible(sample(10)); .Random.seed[1:3]",
			"s <- .Random.seed; s[2] <- 700L; .Random.seed <- s; a <- runif(2); " + seed + "identical(a, runif(2))",
			"s <- .Random.seed; s[2] <- 0L; .Random.seed <- s; a <- runif(2); " + seed + "identical(a, runif(2))",
			"s <- .Random.seed; s[2] <- 623L; .Random.seed <- s; invisible(runif(1)); .Random.seed[2]",
		})
	default:
		family = "seed"
		body = reset + pick([]string{
			"set.seed(1); runif(2)", "set.seed(0); runif(2)", "set.seed(-5); runif(2)", "set.seed(2147483647); runif(2)",
			"set.seed(3.9); runif(2)", "set.seed('7'); runif(2)", "set.seed(NULL); invisible(runif(1)); length(.Random.seed)",
			"tryCatch(set.seed(NA), error = conditionMessage)", "set.seed(42, kind = 'Mersenne-Twister'); runif(1)",
		})
	}
	return "# " + family + "\n" + body
}
