package main

import (
	"fmt"
	"math/rand"
	"strings"
)

var (
	reAtoms  = []string{"a", "b", "c", ".", "[ab]", "[^a]", "[[:alpha:]]", "[[:digit:]]", "\\\\d", "\\\\w", "\\\\s", "x", "1", " ", "[a-c]", "\\\\."}
	reQuants = []string{"", "", "", "*", "+", "?", "{2}", "{1,2}", "{0,1}"}
	reTexts  = []string{"abc", "aab", "banana", "a b c", "abc123", "", "ccc", "a.b", "xyz", "1a2b3", "aaa", "b", "ab ab", "cab"}
)

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
	strFmts  = []string{`"%d"`, `"%s"`, `"%5.2f"`, `"%e"`, `"%g"`, `"%x"`, `"%5s|"`, `"%-5s|"`, `"%05d"`, `"%+d"`, `"%%"`, `"%s=%d"`}
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
