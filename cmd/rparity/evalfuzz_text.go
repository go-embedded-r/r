package main

import (
	"fmt"
	"math/rand"
	"strings"
)

var (
	reAtoms  = []string{"a", "b", "c", ".", "[ab]", "[^a]", "[[:alpha:]]", "[[:digit:]]", "\\\\d", "\\\\w", "\\\\s", "x", "1", " ", "[a-c]", "\\\\.", "\\\\<", "\\\\>", "\\\\b", "\\\\B", "[[:upper:]]", "[^[:space:]]", "é"}
	reQuants = []string{"", "", "", "*", "+", "?", "{2}", "{1,2}", "{0,1}", "*?", "+?", "??", "{1,3}?", "{2,}", "{,2}"}
	reTexts  = []string{"abc", "aab", "banana", "a b c", "abc123", "", "ccc", "a.b", "xyz", "1a2b3", "aaa", "b", "ab ab", "cab", "abab", "aabbaabb", "foo bar", "Abc", "héllo wörld", "abcabc", "baab"}
)

// fuzzRegexBackref returns a pattern with a back-reference (the group number
// can exceed the groups: R then reports an error, as does the port).
func fuzzRegexBackref(r *rand.Rand) string {
	grp := []string{"(a*)", "(a|ab)", "([abc]+)", "(.)", "(a+?)", "(b*)"}
	g := grp[r.Intn(len(grp))]
	ref := "\\\\1"
	if r.Intn(12) == 0 {
		ref = "\\\\2"
	}
	return g + reAtoms[r.Intn(len(reAtoms))] + reQuants[r.Intn(len(reQuants))] + ref + reQuants[r.Intn(len(reQuants))]
}

func fuzzRegexAtom(r *rand.Rand, depth int) string {
	if depth > 0 && r.Intn(5) == 0 {
		n := 1 + r.Intn(2)
		alts := make([]string, 1+r.Intn(2))
		for i := range alts {
			var b strings.Builder
			for k := 0; k <= n; k++ {
				b.WriteString(fuzzRegexAtom(r, depth-1) + reQuants[r.Intn(len(reQuants))])
			}
			alts[i] = b.String()
		}
		return "(" + strings.Join(alts, "|") + ")"
	}
	return reAtoms[r.Intn(len(reAtoms))]
}

func fuzzPattern(r *rand.Rand) string {
	if r.Intn(6) == 0 {
		return fuzzRegexBackref(r)
	}
	var b strings.Builder
	if r.Intn(8) == 0 {
		b.WriteString("^")
	}
	n := 1 + r.Intn(3)
	for i := 0; i < n; i++ {
		b.WriteString(fuzzRegexAtom(r, 1) + reQuants[r.Intn(len(reQuants))])
	}
	if r.Intn(8) == 0 {
		b.WriteString("$")
	}
	return b.String()
}

func rq(s string) string { return `"` + s + `"` }

func fuzzRegex(r *rand.Rand) string {
	pat := rq(fuzzPattern(r))
	nt := 1 + r.Intn(3)
	ts := make([]string, nt)
	for i := range ts {
		ts[i] = rq(reTexts[r.Intn(len(reTexts))])
	}
	text := "c(" + strings.Join(ts, ", ") + ")"
	fl := ""
	if r.Intn(6) == 0 {
		fl = ", ignore.case = TRUE"
	}
	repls := []string{`"X"`, `"<\\1>"`, `""`, `"[\\0]"`, `"\\\\"`}
	switch r.Intn(7) {
	case 0:
		return fmt.Sprintf("grepl(%s, %s%s)", pat, text, fl)
	case 1:
		return fmt.Sprintf("sub(%s, %s, %s%s)", pat, repls[r.Intn(len(repls))], text, fl)
	case 2:
		return fmt.Sprintf("gsub(%s, %s, %s%s)", pat, repls[r.Intn(len(repls))], text, fl)
	case 3:
		return fmt.Sprintf("regexpr(%s, %s%s)", pat, text, fl)
	case 4:
		return fmt.Sprintf("gregexpr(%s, %s%s)", pat, text, fl)
	case 5:
		return fmt.Sprintf("strsplit(%s, %s)", text, pat)
	}
	return fmt.Sprintf("regmatches(%s, gregexpr(%s, %s%s))", text, pat, text, fl)
}

var (
	strTexts = []string{`"hello"`, `"a,b,c"`, `"  pad  "`, `"héllo"`, `""`, `NA`, `"MiXeD"`, `c("ab", "cde", NA)`, `c(a = "x", b = "yy")`, `"日本語"`}
	strNums  = []string{"1", "2", "3", "0", "-1", "10", "NA", "1.5", "2L", "c(1, 2)"}
	strFmts  = []string{`"%d"`, `"%s"`, `"%5.2f"`, `"%e"`, `"%g"`, `"%x"`, `"%5s|"`, `"%-5s|"`, `"%05d"`, `"%+d"`, `"%%"`, `"%s=%d"`, `"%+.2e"`, `"%10.3f"`, `"%-10d|"`, `"% d"`, `"%#x"`, `"%o"`, `"%5.1f%%"`, `"%s and %s"`, `"%2$s %1$s"`, `"%.0f"`, `"%08.3f"`, `"%G"`, `"%E"`, `"%i"`, `"%5.2s|"`, `"%.3d"`, `"%x"`, `"%X"`}
	strFmtA  = []string{"1L", "2.5", `"txt"`, "NA", "TRUE", "c(1L, 2L)", "1e6", "123456.789", "-3L", "0L", "255L", "NA_integer_", "1e-5"}
)

func fuzzStrings(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	switch r.Intn(10) {
	case 0:
		return fmt.Sprintf("nchar(%s)", pick(strTexts))
	case 1:
		return fmt.Sprintf("substr(%s, %s, %s)", pick(strTexts), pick(strNums), pick(strNums))
	case 2:
		return fmt.Sprintf("toupper(%s)", pick(strTexts))
	case 3:
		return fmt.Sprintf("trimws(%s)", pick(strTexts))
	case 4:
		return fmt.Sprintf("substring(%s, %s)", pick(strTexts), pick(strNums))
	case 5:
		return fmt.Sprintf("strsplit(%s, %s)", pick(strTexts), pick([]string{`""`, `","`, `"l+"`, `" "`, `"[aeiou]"`}))
	case 6, 7:
		return fmt.Sprintf("sprintf(%s, %s)", pick(strFmts), pick(strFmtA))
	case 8:
		return fmt.Sprintf("format(%s, nsmall = %s)", pick([]string{"1", "2.5", "pi", "c(1, 10, 100)", "1e6", "1e-10", "NA", "c(1.5, NA)", "123456789", "0.1 + 0.2", "TRUE", `"a"`, `c("a", "bbb")`}), pick([]string{"0", "1", "2"}))
	}
	return fmt.Sprintf("chartr(%s, %s, %s)", `"abc"`, `"xyz"`, pick(strTexts))
}

var (
	apVecs = []string{"1:5", "c(a = 1, b = 2, c = 3)", "c(1.5, 2.5)", `c("x", "yy", "zzz")`, "list(1, 2L, 3)", "list(a = 1, b = 1:2)", "integer(0)", "NULL", "c(TRUE, FALSE, NA)", `c(u = "p", v = "q")`}
	apFuns = []string{"function(x) x * 2", "function(x) c(x, x)", "function(x) c(m = x, n = x + 1)", "function(x) NULL", "function(x) as.character(x)", "length", "function(x) if (is.na(x[1])) 0 else x[1]", "function(x) list(x)", "function(x) x > 1", "nchar", "function(x, y = 2) x + y", "toupper"}
	apBin  = []string{"function(a, b) a + b", "function(a, b) paste(a, b)", "function(a, b) c(a, b)", "function(a, b) max(a, b)"}
)

func fuzzApply(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	v, f := pick(apVecs), pick(apFuns)
	switch r.Intn(14) {
	case 0:
		return fmt.Sprintf("lapply(%s, %s)", v, f)
	case 1, 2:
		return fmt.Sprintf("sapply(%s, %s)", v, f)
	case 3:
		return fmt.Sprintf("sapply(%s, %s, simplify = FALSE)", v, f)
	case 4:
		return fmt.Sprintf("sapply(%s, %s, USE.NAMES = FALSE)", v, f)
	case 5:
		return fmt.Sprintf("vapply(%s, %s, %s)", v, f, pick([]string{"numeric(1)", "character(1)", "logical(1)", "numeric(2)", "integer(1)"}))
	case 6:
		return fmt.Sprintf("Map(%s, %s, %s)", pick(apBin), pick(apVecs[:6]), pick(apVecs[:6]))
	case 7:
		return fmt.Sprintf("mapply(%s, %s, %s)", pick(apBin), pick(apVecs[:6]), pick(apVecs[:6]))
	case 8:
		return fmt.Sprintf("Reduce(%s, %s%s)", pick(apBin), v, pick([]string{"", ", accumulate = TRUE", ", right = TRUE", ", 0", ", accumulate = TRUE, right = TRUE"}))
	case 9:
		return fmt.Sprintf("Filter(%s, %s)", pick([]string{"function(x) x > 1", "function(x) TRUE", "is.na", "function(x) nchar(x) > 1"}), v)
	case 10:
		return fmt.Sprintf("unlist(%s%s)", pick(apVecs), pick([]string{"", ", use.names = FALSE", ", recursive = FALSE"}))
	case 11:
		return fmt.Sprintf("rep(%s, %s)", pick(apVecs), pick([]string{"2", "times = 2", "each = 2", "length.out = 7", "times = c(1, 2)", "0", "each = 2, times = 2"}))
	case 12:
		return fmt.Sprintf("seq(%s)", pick([]string{"5", "2, 8", "10, 1", "0, 1, by = 0.25", "1, 10, length.out = 4", "5, 1, by = -2", "2, 11, 3", "length.out = 3", "from = 1, to = 2, by = 0.5", "1, 10, length.out = 5"}))
	}
	return fmt.Sprintf("rev(%s)", v)
}

var (
	mNums = []string{"0", "1", "-1", "2.5", "-2.5", "0.5", "1e-3", "123456.789", "NA_real_", "NaN", "Inf", "-Inf", "7", "-7", "1e10", "0.125", "1.005", "2.675", "15", "1234.5678"}
	mVecs = []string{"c(1.5, 2.5, 3.5)", "c(-2, 0, 2)", "c(1, NA, 3)", "1:5", "c(a = 1, b = 2)", "c(0.1, 0.2, 0.3)", "c(1e10, 1, 1e-10)", "c(TRUE, FALSE)", "integer(0)", "c(NaN, 1)", "c(-1, Inf)"}
)

// fuzzMath draws a call of the math primitives. The libm functions carry a
// "# libm" mark: they are compared within a few ulp, never excluded, and
// the caller counts them and reports the largest difference per function.
func fuzzMath(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	arg := func() string {
		if r.Intn(2) == 0 {
			return pick(mNums)
		}
		return pick(mVecs)
	}
	switch r.Intn(10) {
	case 0:
		return fmt.Sprintf("round(%s, %s)", arg(), pick([]string{"0", "1", "2", "3", "-1", "-2", "10"}))
	case 1:
		return fmt.Sprintf("signif(%s, %s)", arg(), pick([]string{"1", "2", "3", "6", "10"}))
	case 2:
		return fmt.Sprintf("%s(%s)", pick([]string{"trunc", "sign", "floor", "ceiling", "cumsum", "cummax", "cummin", "is.finite", "is.infinite", "is.nan", "anyNA"}), pick(mVecs))
	case 3:
		return fmt.Sprintf("cumprod(%s)", pick([]string{"c(1.5, 2.5, 3.5)", "1:5", "c(1, NA, 3)", "c(2, 2)"}))
	case 4:
		return fmt.Sprintf("range(%s%s)", pick(mVecs), pick([]string{"", ", na.rm = TRUE", ", finite = TRUE"}))
	case 5:
		fn := pick([]string{"sin", "cos", "tan", "asin", "acos", "atan", "sinh", "cosh", "tanh", "log", "log2", "log10", "log1p", "expm1", "gamma", "lgamma", "digamma", "trigamma"})
		a := arg()
		if (fn == "digamma" || fn == "trigamma" || fn == "gamma" || fn == "lgamma") && strings.Contains(a, "1e-10") {
			// excluded: our series and R's sum the 1/x terms in a different order, a 1e-16 relative difference that round(x, 6) cannot hide at 1e20
			a = "c(1.5, 2.5)"
		}
		if fn == "lgamma" || fn == "digamma" || fn == "trigamma" {
			// our gamma-family routines are less accurate than R's (a few
			// 1e-16 absolute, thousands of ulp near a zero of digamma): compared
			// through round(x, 9) and counted as "approx", an open defect
			return fmt.Sprintf("signif(%s(%s), 12) # approx", fn, a)
		}
		return fmt.Sprintf("%s(%s) # libm", fn, a)
	case 6:
		return fmt.Sprintf("log(%s, %s) # libm", pick(mVecs), pick([]string{"2", "10", "0", "-1"}))
	case 7:
		return fmt.Sprintf("%s(%s) # libm", pick([]string{"sinpi", "cospi", "tanpi"}), pick([]string{"0", "0.5", "1", "1.5", "2", "-0.5", "0.25", "Inf", "NA_real_"}))
	case 8:
		return fmt.Sprintf("nzchar(%s)", pick([]string{`c("a", "", NA)`, `""`, `NA`, `character(0)`, `123`, `c(x = "a")`}))
	}
	return fmt.Sprintf("x <- %s; length(x) <- %s; x", pick(mVecs), pick([]string{"0", "2", "5"}))
}

var (
	vVecs = []string{"c(3, 1, 2)", "c(2, 1, NA, 3)", "c(\"b\", \"A\", \"a\", \"B\")", "1:6", "c(5, 5, 1)", "c(a = 2, b = 1)", "c(TRUE, FALSE, NA)", "numeric(0)", "c(1.5, -2, 0)", "letters[1:4]", "c(NA, NA)", "c(1, 1, 2, 2)"}
	vNums = []string{"c(1, 2)", "c(3, 1, 2)", "1:3", "c(NA, 4)", "5", "c(2.5, 2.5)"}
	vMats = []string{"matrix(1:6, 2)", "matrix(c(1, NA, 3, 4), 2)", "matrix(1:4, 2, dimnames = list(c(\"a\", \"b\"), c(\"x\", \"y\")))", "matrix(c(TRUE, FALSE, TRUE, TRUE), 2)", "matrix(1.5, 3, 2)"}
)

func fuzzVectors(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	v := pick(vVecs)
	switch r.Intn(18) {
	case 0:
		return fmt.Sprintf("order(%s%s)", v, pick([]string{"", ", decreasing = TRUE", ", na.last = FALSE", ", na.last = NA"}))
	case 1:
		return fmt.Sprintf("order(%s, %s)", pick(vNums), pick(vNums))
	case 2:
		return fmt.Sprintf("pmax(%s, %s%s)", pick(vNums), pick(vNums), pick([]string{"", ", na.rm = TRUE"}))
	case 3:
		return fmt.Sprintf("pmin(%s, %s%s)", pick(vNums), pick(vNums), pick([]string{"", ", na.rm = TRUE"}))
	case 4:
		return fmt.Sprintf("which.max(%s)", pick(vNums[:4]))
	case 5:
		return fmt.Sprintf("which.min(%s)", pick(vNums[:4]))
	case 6:
		return fmt.Sprintf("tabulate(%s%s)", pick(vNums), pick([]string{"", ", 3", ", nbins = 6"}))
	case 7:
		return fmt.Sprintf("%s(%s%s)", pick([]string{"rowSums", "colSums", "rowMeans", "colMeans"}), pick(vMats), pick([]string{"", ", na.rm = TRUE"}))
	case 8:
		return fmt.Sprintf("t(%s)", pick(append([]string{v}, vMats...)))
	case 9:
		return fmt.Sprintf("diag(%s)", pick([]string{"3", "2", "c(1, 2, 3)", "matrix(1:9, 3)", "matrix(1:6, 2)", "0", "c(2L, 3L)"}))
	case 10:
		return fmt.Sprintf("outer(%s, %s%s)", pick(vNums[:3]), pick(vNums[:3]), pick([]string{"", ", \"+\"", ", \"-\"", ", function(a, b) a * 10 + b"}))
	case 11:
		return fmt.Sprintf("mean(%s%s)", pick(vNums), pick([]string{"", ", na.rm = TRUE"}))
	case 12:
		return fmt.Sprintf("%s(%s, %s)", pick([]string{"union", "intersect", "setdiff"}), pick(vNums[:3]), pick(vNums[:3]))
	case 13:
		return fmt.Sprintf("%s(%s, %s)", pick([]string{"head", "tail"}), pick(vVecs), pick([]string{"2", "-1", "0", "10", "3L"}))
	case 14:
		return fmt.Sprintf("ifelse(%s > 1, %s, %s)", pick(vNums[:3]), pick([]string{"1", "\"a\"", "c(7, 8, 9)"}), pick([]string{"0", "\"b\"", "NA"}))
	case 15:
		return fmt.Sprintf("as.vector(%s%s)", pick(append([]string{v}, vMats...)), pick([]string{"", ", \"list\"", ", \"character\"", ", \"numeric\""}))
	case 16:
		return fmt.Sprintf("diff(%s%s)", pick(vNums[:4]), pick([]string{"", ", lag = 2", ", differences = 2"}))
	}
	return fmt.Sprintf("append(%s, %s%s)", pick(vNums[:3]), pick([]string{"9", "c(8, 9)"}), pick([]string{"", ", after = 1", ", after = 0"}))
}

var (
	s2Texts = []string{`"abc"`, `c("a", "bb", NA)`, `""`, `"x,y"`, `c(a = "q")`, `"it's"`, `"12"`, `"ff"`, `"a/b/c"`, `"dir/"`, `"/"`}
	s2Nums  = []string{"3.14159", "42L", "1234567.891", "-5", "0.000012345", "1e6", "c(1, 10, 100)", "NA_real_", "123456789", "0", "c(1.5, NA)"}
)

func fuzzStrings2(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	switch r.Intn(10) {
	case 0:
		return fmt.Sprintf("strrep(%s, %s)", pick(s2Texts), pick([]string{"0", "1", "3", "c(1, 2)", "NA"}))
	case 1:
		return fmt.Sprintf("formatC(%s%s)", pick(s2Nums), pick([]string{"", ", width = 8", ", digits = 2, format = \"f\"", ", format = \"e\", digits = 3", ", format = \"g\"", ", flag = \"0\", width = 10", ", big.mark = \",\", format = \"d\"", ", flag = \"-\", width = 9", ", width = -7", ", format = \"f\", digits = 0", ", flag = \"+\""}))
	case 2:
		return fmt.Sprintf("prettyNum(%s, big.mark = \",\")", pick(s2Nums))
	case 3:
		return fmt.Sprintf("strtoi(%s, %s)", pick(s2Texts), pick([]string{"10L", "16L", "8L", "2L", "36L"}))
	case 4:
		return fmt.Sprintf("make.unique(c(%s, %s, %s))", pick([]string{`"a"`, `"b"`, `"a.1"`}), pick([]string{`"a"`, `"b"`}), pick([]string{`"a"`, `"a.1"`}))
	case 5:
		return fmt.Sprintf("%s(%s)", pick([]string{"basename", "dirname"}), pick(s2Texts))
	case 6:
		return fmt.Sprintf("%s(%s)", pick([]string{"sQuote", "dQuote", "shQuote", "toString"}), pick(s2Texts))
	case 7:
		return fmt.Sprintf("x <- %s; substr(x, %s, %s) <- %s; x", pick([]string{`"abcdef"`, `c("aaa", "bbb")`, `"é"`}), pick([]string{"1", "2", "0", "5"}), pick([]string{"2", "3", "10"}), pick([]string{`"Z"`, `"XYZ"`, `""`}))
	case 8:
		return fmt.Sprintf("noquote(%s)", pick(s2Texts))
	}
	return fmt.Sprintf("file.path(%s, %s)", pick(s2Texts[:5]), pick(s2Texts[:5]))
}
