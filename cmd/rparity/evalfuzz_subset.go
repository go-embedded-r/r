package main

import (
	"math/rand"
	"strings"
)

var (
	subBases = []string{
		"c(10, 20, 30, 40)", "1:5", "c(a = 1, b = 2, c = 3)", "c(TRUE, FALSE, NA)",
		"c(\"p\", \"q\", \"r\")", "c(x = \"p\", y = \"q\")", "list(1, \"a\", TRUE)",
		"list(a = 1, b = \"s\", c = 1:3)", "numeric(0)", "NULL", "c(a = 1L, b = NA)",
		"matrix(1:6, 2)", "matrix(1:6, 2, dimnames = list(c(\"r1\", \"r2\"), c(\"A\", \"B\", \"C\")))",
		"array(1:24, c(2, 3, 4))",
	}
	subIdx = []string{
		"1", "2", "5", "0", "-1", "-2", "NA", "NA_integer_", "TRUE", "FALSE", "c(TRUE, FALSE)",
		"c(1, 3)", "c(-1, -2)", "c(0, 1)", "c(1, NA)", "1:2", "2:1", "c(1, 1)", "\"a\"", "\"b\"",
		"c(\"a\", \"c\")", "\"zz\"", "1.9", "-1.9", "integer(0)", "c(TRUE, NA)", "c(1, -1)", "4", "c(2, 4)",
		"\"x\"", "\"p\"", "NULL",
	}
	subVal = []string{
		"99", "c(7, 8)", "NA", "\"s\"", "TRUE", "NULL", "list(1)", "1:3", "c(a = 5)", "0L", "2.5",
	}
	subDim = []string{"1", "2", "3", "c(1, 2)", "TRUE", "-1", "0", "NA", "\"A\"", "\"r2\"", "\"z\"", "5", "integer(0)", "1:2"}
)

// fuzzSubset builds a random subsetting or sub-assignment program.
func fuzzSubset(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	x := pick(subBases)
	var body string
	switch r.Intn(9) {
	case 0:
		body = "x[" + pick(subIdx) + "]"
	case 1:
		body = "x[[" + pick(subIdx) + "]]"
	case 2:
		body = "x[" + pick(subIdx) + "] <- " + pick(subVal) + "; x"
	case 3:
		body = "x[[" + pick(subIdx) + "]] <- " + pick(subVal) + "; x"
	case 4:
		body = "x$" + pick([]string{"a", "b", "zz", "c"})
	case 5:
		body = "x$" + pick([]string{"a", "b", "new"}) + " <- " + pick(subVal) + "; x"
	case 6:
		body = "x[" + pick(subDim) + ", " + pick(subDim) + "]"
		if r.Intn(2) == 0 {
			body = "x[" + pick(subDim) + ", " + pick(subDim) + ", drop = FALSE]"
		}
	case 7:
		body = "x[" + pick(subDim) + ", " + pick(subDim) + "] <- " + pick(subVal) + "; x"
	default:
		body = "names(x) <- " + pick([]string{"c(\"u\", \"v\")", "NULL", "\"n\"", "c(\"a\", NA)", "letters"}) + "; x"
	}
	return strings.Join([]string{"x <- " + x, body}, "; ")
}
