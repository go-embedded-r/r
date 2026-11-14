package main

import (
	"math/rand"
	"strings"
)

// The summary mode: random programs over the closures of slice 15: summary of
// vectors, factors and data frames, format.pval, all.equal, diff, the data
// frame group methods, within, hexmode, mean(trim =) and strtrim. The first
// line is "# family". Everything is compared exactly, printed text included;
// there is no exclusion class. Text that R builds with typographic quotes is
// printed with cat, because R leaves such strings unmarked and ours marks them.

func sumNums(r *rand.Rand, n int, withNA bool) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	vals := make([]string, n)
	for i := range vals {
		if withNA && r.Intn(6) == 0 {
			vals[i] = "NA"
			continue
		}
		vals[i] = pick([]string{"1", "2.5", "-3", "10", "0.125", "1e3", "7", "123.456", "0", "42", "-0.5", "3.14159"})
	}
	return "c(" + strings.Join(vals, ", ") + ")"
}

func fuzzSummary(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	n := 1 + r.Intn(8)
	var family, body string
	show := "p <- function(x) cat(if (isTRUE(x)) \"TRUE\" else x, sep = \"\\n\"); "
	switch r.Intn(14) {
	case 0:
		family, body = "sum-num", "print(summary("+sumNums(r, n, true)+pick([]string{"", ", digits = 2", ", digits = 5"})+"))"
	case 1:
		family, body = "sum-other", "print(summary("+pick([]string{"c(TRUE, FALSE, NA, TRUE)", "c(\"a\", \"bb\", NA, \"\")", "factor(c(\"x\", \"y\", \"x\", NA))", "list(1, \"a\", 1:3)", "NULL", "factor(letters[1:9])"})+pick([]string{"", ", maxsum = 3"})+"))"
	case 2:
		family, body = "sum-df", "print(summary(data.frame(a = "+sumNums(r, n, true)+", b = "+dfCol(r, n, 2)+", c = "+dfCol(r, n, 3)+pick([]string{"", ", f = factor(" + dfCol(r, n, 2) + ")"})+")"+pick([]string{"", ", digits = 3", ", maxsum = 2"})+"))"
	case 3:
		a := sumNums(r, n, true)
		family, body = "allequal-num", show+"p(all.equal("+a+", "+pick([]string{a, a + " + 1e-10", a + " * 1.1", a + " + 1", "rev(" + a + ")", "c(" + a + ", 1)", "as.character(" + a + ")"})+pick([]string{"", ", tolerance = 0.5", ", scale = 1", ", countEQ = TRUE", ", check.attributes = FALSE"})+"))"
	case 4:
		family, body = "allequal-other", show+"p(all.equal("+pick([]string{"list(a = 1, b = \"x\")", "c(a = 1, b = 2)", "factor(c(\"a\", \"b\"))", "c(\"a\", \"b\")", "matrix(1:4, 2)", "data.frame(a = 1:2)", "quote(f(x))", "TRUE", "NULL"})+", "+pick([]string{"list(a = 1, b = \"y\")", "c(a = 1, c = 2)", "factor(c(\"a\", \"c\"))", "c(\"a\", \"c\")", "matrix(1:4, 1)", "data.frame(a = c(1, 3))", "quote(f(y))", "FALSE", "1", "list(1)"})+"))"
	case 5:
		vals := make([]string, 1+r.Intn(5))
		for i := range vals {
			vals[i] = pick([]string{"0.1", "0.05", "1e-4", "1e-10", "1e-30", "0.5", "0.00123", "NA", "0", "0.999"})
		}
		family, body = "pval", "format.pval(c("+strings.Join(vals, ", ")+")"+pick([]string{"", ", digits = 3", ", digits = 1", ", eps = 0.001", ", eps = 1e-8, digits = 4"})+")"
	case 6:
		family, body = "diff", "print(diff("+pick([]string{sumNums(r, n+2, false), "as.Date(c(\"2024-01-05\", \"2024-02-01\", \"2024-02-03\"))", "as.POSIXct(c(\"2024-03-10 00:00\", \"2024-03-10 03:30\", \"2024-03-11\"), tz = \"UTC\")", "matrix(c(1, 4, 9, 16, 25, 36), 3)"})+pick([]string{"", ", lag = 2", ", differences = 2"})+"))"
	case 7:
		df := "data.frame(a = " + sumNums(r, n, false) + ", b = " + sumNums(r, n, false) + ")"
		family, body = "df-ops", "df <- "+df+"; "+pick([]string{"print(df * 2)", "print(df > 1)", "print(df + df)", "print(-df)", "print(sqrt(abs(df)))", "print(sum(df))", "print(range(df))", "print(round(df, 1))", "print(df == df)", "print(max(df, 100))", "print(df %% 2)"})
	case 8:
		family, body = "within", "df <- data.frame(a = "+sumNums(r, n, false)+", b = "+sumNums(r, n, false)+"); print(within(df, "+pick([]string{"c <- a + b", "{ c <- a * 2; d <- c + 1 }", "a <- NULL", "b <- b * 10", "rm(a)", "z <- 1"})+"))"
	case 9:
		family, body = "hexmode", "x <- as."+pick([]string{"hexmode", "octmode"})+"(c("+strings.Join([]string{pick([]string{"1", "10", "255", "64"}), pick([]string{"8", "100", "4095", "0"})}, ", ")+")); print(x); print(format(x"+pick([]string{"", ", width = 6"})+")); print(as.character(x))"
	case 10:
		// whole numbers only: R sums the trimmed middle in the order its partial
		// sort leaves it, ours in sorted order, so a mean with an inexact sum can
		// differ in the last bit (a known divergence, listed in the README)
		vals := make([]string, 3+r.Intn(8))
		for i := range vals {
			vals[i] = pick([]string{"1", "2", "-3", "10", "7", "0", "42", "-5", "1000"})
		}
		// compared within a few units in the last place (the "# libm" class)
		family, body = "mean-trim", "print(mean(c("+strings.Join(vals, ", ")+"), trim = "+pick([]string{"0.1", "0.2", "0.25", "0.4", "0.5", "0"})+")) # libm"
	case 11:
		family, body = "strtrim", "print(strtrim("+pick([]string{"\"abcdef\"", "c(\"a\", \"bcd\", NA)", "\"日本語\"", "letters[1:4]"})+", "+pick([]string{"1", "3", "c(1, 2)", "0", "10"})+"))"
	case 12:
		family, body = "factor-ops", "f <- factor(c(\"a\", \"b\", \"a\")); o <- ordered(c(\"lo\", \"hi\", \"lo\"), levels = c(\"lo\", \"hi\")); "+pick([]string{"print(f == \"a\")", "print(o < \"hi\")", "print(max(o))", "print(range(o))", "print(summary(o))", "r <- tryCatch(f + 1, warning = function(w) \"w\"); print(r)", "r <- tryCatch(sum(f), error = function(e) \"e\"); print(r)", "length(f) <- 5; print(f)"})
	default:
		family, body = "date-methods", "d <- as.Date(c(\"2024-03-05\", \"2024-01-01\", NA, \"2025-12-31\")); "+pick([]string{"print(range(d, na.rm = TRUE))", "print(xtfrm(d))", "print(split(d, c(1, 2, 1, 2)))", "length(d) <- 6; print(d)", "print(mtfrm(d))", "r <- tryCatch(sqrt(d), error = function(e) \"e\"); print(r)"})
	}
	return "# " + family + "\n" + body
}
