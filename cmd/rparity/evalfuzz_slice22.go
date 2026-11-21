package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// The slice22 mode: random programs over the additions of slice 22: charmatch
// and char.expand, findInterval, choose, numToBits and numToInts, encodeString,
// format.info, sort.int and sort.list, rapply, asplit, drop, max.col, callCC and
// the assignment methods of POSIXct and difftime. The first line is "# family".
// A value ends in "# libm" when it is compared within --ulp units in the last
// place after being rounded to seven digits (choose is held to a stated bound,
// not to R's product loop).

func fuzzSlice22(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	nums := func(n, hi int) string {
		v := make([]string, n)
		for i := range v {
			v[i] = fmt.Sprint(r.Intn(hi))
		}
		return "c(" + strings.Join(v, ", ") + ")"
	}
	strs := func(n int, alpha []string) string {
		v := make([]string, n)
		for i := range v {
			v[i] = fmt.Sprintf("%q", alpha[r.Intn(len(alpha))])
		}
		return "c(" + strings.Join(v, ", ") + ")"
	}
	var family, body string
	switch r.Intn(9) {
	case 0:
		family = "charmatch"
		tbl := strs(2+r.Intn(4), []string{"mean", "median", "mode", "me", "max", "min", ""})
		body = "tb <- " + tbl + "; " + pick([]string{
			"charmatch('me', tb)", "charmatch(c('m', 'me', 'x', ''), tb)", "charmatch(NA, tb)", "charmatch('mean', tb, nomatch = 0L)",
			"char.expand('me', tb, nomatch = NULL)", "tryCatch(char.expand('zz', tb), error = conditionMessage)", "charmatch(character(), tb)",
			"charmatch('m', character())", "charmatch(factor('ma'), tb)",
		})
	case 1:
		family = "findInterval"
		n := 1 + r.Intn(6)
		body = "v <- sort(" + nums(n, 10) + "); x <- c(" + fmt.Sprint(r.Intn(12)-1) + ", " + fmt.Sprint(r.Intn(12)-1) + ".5, NA, " + fmt.Sprint(r.Intn(10)) + "); " + pick([]string{
			"findInterval(x, v)", "findInterval(x, v, rightmost.closed = TRUE)", "findInterval(x, v, all.inside = TRUE)", "findInterval(x, v, left.open = TRUE)",
			"findInterval(x, v, left.open = TRUE, rightmost.closed = TRUE)", "findInterval(x, v, left.open = TRUE, all.inside = TRUE)",
			"tryCatch(findInterval(1, rev(v + 0:" + fmt.Sprint(n-1) + ")), error = conditionMessage)", "length(findInterval(numeric(), v))",
		})
	case 2:
		family = "choose"
		body = pick([]string{
			fmt.Sprintf("choose(%d, %d)", r.Intn(60), r.Intn(12)), fmt.Sprintf("choose(-%d, %d)", 1+r.Intn(20), r.Intn(9)), fmt.Sprintf("choose(%d.5, %d)", r.Intn(20), r.Intn(8)),
			fmt.Sprintf("choose(%d, 0:%d)", 5+r.Intn(30), 3+r.Intn(6)), fmt.Sprintf("choose(c(a = %d, b = %d), 2)", r.Intn(30), r.Intn(30)),
			fmt.Sprintf("choose(%d, %d.4)", 5+r.Intn(30), r.Intn(5)), fmt.Sprintf("choose(%d, -1)", r.Intn(9)), "choose(NA, 2)",
		}) + " # libm"
	case 3:
		family = "numbits"
		x := nums(1+r.Intn(3), 50) + " / " + fmt.Sprint(1+r.Intn(7))
		body = "x <- " + x + "; " + pick([]string{
			"numToBits(x)[1:12]", "numToInts(x)", "length(numToBits(x))", "sum(as.integer(numToBits(x)))", "as.integer(numToBits(x))[53:64]",
			"numToInts(-x)", "numToBits(c(NA, Inf, -Inf, 0))[1:4]",
		})
	case 4:
		family = "encodeString"
		x := strs(1+r.Intn(4), []string{"a", "bb", "a\nb", "q\"x", "tab\t", "", "\\"})
		body = "x <- " + x + "; " + pick([]string{
			"encodeString(x)", "encodeString(x, quote = '\"')", "encodeString(x, quote = \"'\")", "encodeString(x, width = NA)", "encodeString(x, width = 6, justify = 'right')",
			"encodeString(x, width = 7, justify = 'centre')", "encodeString(c(x, NA))", "encodeString(c(x, NA), na.encode = FALSE)", "encodeString(c(a = x[1]))",
			"encodeString(x, width = 5, justify = 'none')",
		})
	case 5:
		family = "formatinfo"
		body = pick([]string{
			"format.info(c(1, 10, 100))", "format.info(c(1.5, NA, 22.25))", "format.info(1e10)", "format.info(1e-5)", "format.info(123456789)", "format.info(pi, digits = 3)",
			"format.info(1, nsmall = 3)", "format.info(1:" + fmt.Sprint(1+r.Intn(1000)) + ")", "format.info(c(TRUE, NA))", "format.info(c('ab', NA, 'cde'))",
			"format.info(c(-1, 1))", "format.info(1e300)", "format.info(c(0.1, 0.12, 0.123))",
		})
	case 6:
		family = "sorting"
		x := nums(3+r.Intn(8), 9)
		body = "x <- " + x + "; x[2] <- NA; " + pick([]string{
			"sort.int(x)", "sort.int(x, na.last = TRUE)", "sort.int(x, decreasing = TRUE)", "sort.int(x, index.return = TRUE)", "sort.list(x)", "sort.list(x, decreasing = TRUE)",
			"sort.list(x, na.last = NA)", "sort.list(x, na.last = FALSE)", "sort.int(c(b = 2, a = 1))", "sort.int(letters[c(3, 1, 2)])",
		})
	case 7:
		family = "structure"
		body = pick([]string{
			"rapply(list(1, 'a', list(2, 3)), function(z) z * 10, classes = 'numeric')", "rapply(list(a = 1, b = list(c = 2, d = 'x')), function(z) z + 1, classes = 'numeric', how = 'replace')",
			"rapply(list(1:3, 4:6), sum)", "rapply(list(a = 1, b = 'x'), toupper, classes = 'character', deflt = NA, how = 'unlist')",
			"asplit(matrix(1:6, 2), 1)", "asplit(matrix(1:6, 2), 2)", "asplit(array(1:8, c(2, 2, 2)), c(1, 3))", "drop(matrix(1:3, 1))", "drop(array(1:6, c(1, 2, 3)))",
			"drop(array(1:4, c(2, 1, 2), dimnames = list(c('a', 'b'), 'k', c('u', 'v'))))", "max.col(matrix(c(1, 5, 3, 9, 2, 2), 2), 'first')", "max.col(matrix(c(1, 5, 3, 9, 2, 2), 2), 'last')",
			"t.default(matrix(1:6, 2))", "x <- 1:2; comment(x) <- 'hi'; comment(x)",
		})
	default:
		family = "callcc"
		body = pick([]string{
			"callCC(function(k) { k(" + fmt.Sprint(r.Intn(9)) + "); 99 })", "callCC(function(k) " + fmt.Sprint(r.Intn(9)) + ")",
			"p <- as.POSIXct('2020-01-01', tz = 'UTC') + 0:3; p[2] <- as.POSIXct('2021-06-01', tz = 'UTC'); p",
			"p <- as.POSIXct('2020-01-01', tz = 'UTC') + 0:3; p[c(1, 3)] <- NA; attributes(p)",
			"d <- as.difftime(c(1, 2, 3), units = 'mins'); d[2] <- as.difftime(1, units = 'hours'); d",
			"d <- as.difftime(c(1, 2, 3), units = 'secs'); d[3] <- 9; unclass(d)",
			".mapply(function(a, b) a * b, list(1:3, 4:6), NULL)", ".makeMessage('a', 1:2, appendLF = TRUE)", "as.character(factor(c('b', 'a', 'b')))",
			"f <- factor(c('a', 'b')); levels(f) <- c('x', 'x'); .valid.factor(f)", "dontCheck(" + fmt.Sprint(r.Intn(9)) + ")",
		})
	}
	return "# " + family + "\n" + body
}
