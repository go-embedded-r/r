package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// The slice23 mode: random programs over the additions of slice 23: lchoose,
// abbreviate, agrep and agrepl, the locale functions, Sys.which, summary of
// warnings, the small dot helpers and c(recursive =). The first line is
// "# family". A value ends in "# libm" when it is compared within --ulp units in
// the last place after being rounded to seven digits (lchoose is held to a
// stated bound, not to R's lgamma).

func fuzzSlice23(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	words := []string{"Alabama", "Alaska", "Arizona", "Illinois", "Indiana", "New Hampshire", "North Carolina", "abc def", "International", "x", "lazy", "laser", "lasy"}
	strs := func(n int) string {
		v := make([]string, n)
		for i := range v {
			v[i] = fmt.Sprintf("%q", words[r.Intn(len(words))])
		}
		return "c(" + strings.Join(v, ", ") + ")"
	}
	var family, body string
	switch r.Intn(7) {
	case 0:
		family = "lchoose"
		body = "signif(" + pick([]string{
			fmt.Sprintf("lchoose(%d, %d)", 2+r.Intn(80), r.Intn(12)), fmt.Sprintf("lchoose(-%d, %d)", 1+r.Intn(20), r.Intn(9)),
			fmt.Sprintf("lchoose(%d.5, %d)", r.Intn(20), r.Intn(8)), fmt.Sprintf("lchoose(%d, 0:%d)", 5+r.Intn(30), 3+r.Intn(6)),
			fmt.Sprintf("lchoose(%d, %d)", 100+r.Intn(900), 1+r.Intn(500)),
		}) + ", 7) # libm"
	case 1:
		family = "abbreviate"
		body = "x <- " + strs(1+r.Intn(5)) + "; " + pick([]string{
			fmt.Sprintf("abbreviate(x, %d)", 2+r.Intn(6)), fmt.Sprintf("abbreviate(x, %d, strict = TRUE)", 2+r.Intn(6)),
			fmt.Sprintf("abbreviate(x, %d, use.classes = FALSE)", 2+r.Intn(6)), fmt.Sprintf("abbreviate(x, %d, named = FALSE)", 2+r.Intn(6)),
			fmt.Sprintf("abbreviate(x, %d, dot = TRUE)", 2+r.Intn(6)), fmt.Sprintf("abbreviate(x, %d, method = 'left.kept')", 2+r.Intn(6)),
		})
	case 2:
		family = "agrep"
		body = "x <- " + strs(2+r.Intn(5)) + "; " + pick([]string{
			"agrep('lasy', x)", "agrepl('lasy', x)", "agrep('lasy', x, max.distance = 2, value = TRUE)", "agrepl('LASY', x, ignore.case = TRUE)",
			"agrepl('lasy', x, max.distance = list(sub = 0, all = 1))", "agrepl('abc', x, max.distance = 0)", "agrep('ala', x, max.distance = 1)",
			"agrepl('lasy', x, max.distance = list(ins = 0, del = 0, sub = 1))",
		})
	case 3:
		family = "locale"
		body = pick([]string{
			"Sys.getlocale('LC_COLLATE')", "Sys.getlocale('LC_CTYPE')", "tryCatch(Sys.getlocale('xx'), error = conditionMessage)",
			"suppressWarnings(Sys.setlocale('LC_TIME', 'xx_YY'))", "length(Sys.localeconv())", "Sys.which('zz-none')", "Sys.which(c('zz-a', NA))",
			"tryCatch(Sys.setlocale('zz', 'C'), error = conditionMessage)",
		})
	case 4:
		family = "warnings"
		body = "w <- structure(list(quote(f(x)), quote(f(x)), NULL, quote(g(y))), names = c('a', 'a', 'b', " + pick([]string{"'a'", "'c'"}) + "), class = 'warnings'); " + pick([]string{
			"print(summary(w))", "attr(summary(w), 'counts')", "length(summary(w))", "print(summary(w[0]))",
		})
	case 5:
		family = "helpers"
		body = pick([]string{
			".gt(c(3, 1, 2), 1, 2)", ".gtn(c(1, 2, 2), TRUE)", ".gtn(c(1, 2, 3), FALSE)", ".format.zeros(c('0.0', '1.5'), '-')",
			".rangeNum(c(3, 1, 2), na.rm = FALSE, finite = FALSE, isNumeric = is.numeric)",
			".rangeNum(c(3, NA, Inf), na.rm = TRUE, finite = FALSE, isNumeric = is.numeric)", ".tryResumeInterrupt()",
			"m <- taskCallbackManager(); m$add(function(...) TRUE, name = 'q'); names(m$callbacks())",
		})
	default:
		family = "combine"
		body = pick([]string{
			"c(1, list(2, list(3, 'a')), recursive = TRUE)", "c(list(a = 1), b = 2, recursive = TRUE)", "c(a = 1, b = 2, use.names = FALSE)",
			"c(list(a = 1, b = list(c = 2)), recursive = TRUE, use.names = FALSE)", "c(list(1, 2), recursive = FALSE)", "c(NULL, recursive = TRUE)",
		})
	}
	return "# " + family + "\n" + body
}
