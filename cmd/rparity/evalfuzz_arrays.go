package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// The arrays mode: random programs over matrices and arrays (apply, sweep,
// scale, aperm, the margin functions, the dimnames accessors and
// replacements, lower.tri, diag<-), grouping (split, unsplit, tapply, rle,
// sequence) and the small helpers of slice 11. The first line is "# family".
// Nothing is excluded: everything is compared exactly.
func fuzzArrays(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	nums := func(n int, withNA bool) string {
		parts := make([]string, n)
		for i := range parts {
			switch {
			case withNA && r.Intn(8) == 0:
				parts[i] = "NA"
			case r.Intn(3) == 0:
				parts[i] = fmt.Sprintf("%d.5", r.Intn(9)-4)
			default:
				parts[i] = fmt.Sprint(r.Intn(11) - 5)
			}
		}
		return "c(" + strings.Join(parts, ", ") + ")"
	}
	names := func(prefix string, n int) string {
		parts := make([]string, n)
		for i := range parts {
			parts[i] = fmt.Sprintf("\"%s%d\"", prefix, i+1)
		}
		return "c(" + strings.Join(parts, ", ") + ")"
	}
	dims := func() (int, int) { return 1 + r.Intn(4), 1 + r.Intn(4) }
	mat := func(withNA bool) (string, int, int) {
		nr, nc := dims()
		m := fmt.Sprintf("matrix(%s, %d, %d", nums(nr*nc, withNA), nr, nc)
		switch r.Intn(4) {
		case 1:
			m += ", dimnames = list(" + names("r", nr) + ", " + names("c", nc) + ")"
		case 2:
			m += ", dimnames = list(rows = " + names("r", nr) + ", cols = " + names("c", nc) + ")"
		}
		return m + ")", nr, nc
	}
	funs := []string{"sum", "range", "rev", "length", "function(v) v * 2", "function(v) v[v > 0]", "max", "function(v) c(lo = min(v), hi = max(v))", "function(v) paste(v, collapse = \"-\")", "function(v) if (sum(v) > 0) 1 else c(1, 2)"}
	var family, body string
	switch r.Intn(14) {
	case 0:
		m, _, _ := mat(true)
		family = "apply"
		body = "apply(" + m + ", " + pick([]string{"1", "2", "c(1, 2)", "c(2, 1)"}) + ", " + pick(funs) + pick([]string{"", ", na.rm = TRUE"}) + ")"
		if strings.Contains(body, "na.rm = TRUE") && !strings.Contains(body, ", sum") && !strings.Contains(body, ", max") && !strings.Contains(body, ", range") {
			body = strings.Replace(body, ", na.rm = TRUE", "", 1)
		}
	case 1:
		m, _, nc := mat(false)
		family = "apply-array"
		a := fmt.Sprintf("array(%s, c(2, 3, %d))", nums(6*nc, false), nc)
		body = "apply(" + a + ", " + pick([]string{"1", "2", "3", "c(1, 3)", "c(2, 3)", "c(1, 2)"}) + ", " + pick([]string{"sum", "range", "function(v) v[1]", "length"}) + ")"
		_ = m
	case 2:
		m, nr, nc := mat(false)
		family = "sweep"
		if r.Intn(2) == 0 {
			body = "sweep(" + m + ", 1, " + nums(nr, false) + pick([]string{"", ", \"+\"", ", \"*\"", ", \"/\""}) + ")"
		} else {
			body = "sweep(" + m + ", 2, " + nums(nc, false) + pick([]string{"", ", \"+\"", ", \"*\""}) + ")"
		}
	case 3:
		m, _, _ := mat(r.Intn(3) == 0)
		family = "scale"
		body = "scale(" + m + pick([]string{"", ", center = FALSE", ", scale = FALSE", ", center = FALSE, scale = FALSE"}) + ")"
	case 4:
		_, _, _ = mat(false)
		family = "aperm"
		a := fmt.Sprintf("array(1:24, 2:4%s)", pick([]string{"", ", list(c(\"a\", \"b\"), c(\"p\", \"q\", \"r\"), NULL)", ", list(A = c(\"a\", \"b\"), B = c(\"p\", \"q\", \"r\"), C = c(\"u\", \"v\", \"w\", \"x\"))"}))
		body = "aperm(" + a + ", " + pick([]string{"c(1, 2, 3)", "c(2, 1, 3)", "c(3, 2, 1)", "c(2, 3, 1)", "c(3, 1, 2)", "c(1, 3, 2)"}) + ")"
	case 5:
		m, _, _ := mat(false)
		family = "margins"
		f := pick([]string{"prop.table", "proportions", "marginSums", "margin.table"})
		body = f + "(abs(" + m + ")" + pick([]string{"", ", 1", ", 2"}) + ")"
	case 6:
		m, nr, nc := mat(false)
		family = "dimnames"
		switch r.Intn(5) {
		case 0:
			body = "list(rownames(" + m + "), colnames(" + m + "))"
		case 1:
			body = "list(rownames(" + m + ", do.NULL = FALSE), colnames(" + m + ", do.NULL = FALSE, prefix = \"V\"))"
		case 2:
			body = "x <- " + m + "; rownames(x) <- " + pick([]string{names("a", nr), "NULL"}) + "; colnames(x) <- " + pick([]string{names("z", nc), "NULL"}) + "; x"
		case 3:
			body = "lower.tri(" + m + pick([]string{"", ", diag = TRUE"}) + ")"
		default:
			body = "x <- " + m + "; diag(x) <- " + pick([]string{"0", "9", nums(min(nr, nc), false)}) + "; x"
		}
	case 7:
		family = "split"
		n := 1 + r.Intn(8)
		fp := make([]string, n)
		for i := range fp {
			fp[i] = fmt.Sprintf("\"%s\"", pick([]string{"a", "b", "c", "d"}))
		}
		f := "c(" + strings.Join(fp, ", ") + ")"
		if r.Intn(3) == 0 {
			// a grouping factor shorter than the data is recycled (with a warning when it does not divide)
			f = "c(" + strings.Join(fp[:1+r.Intn(n)], ", ") + ")"
		}
		switch r.Intn(4) {
		case 0:
			body = "split(" + nums(n, false) + ", " + f + ")"
		case 1:
			body = "split(1:" + fmt.Sprint(n) + ", factor(" + f + ", levels = c(\"d\", \"c\", \"b\", \"a\"))" + pick([]string{"", ", drop = TRUE"}) + ")"
		case 2:
			body = "x <- " + nums(n, false) + "; f <- " + f + "; unsplit(split(x, f), f)"
		default:
			body = "x <- 1:" + fmt.Sprint(n) + "; split(x, " + f + ") <- list(-1L, -2L, -3L, -4L); x"
		}
	case 8:
		family = "tapply"
		n := 2 + r.Intn(8)
		f := func() string {
			p := make([]string, n)
			for i := range p {
				p[i] = fmt.Sprintf("\"%s\"", pick([]string{"a", "b", "c"}))
			}
			return "c(" + strings.Join(p, ", ") + ")"
		}
		switch r.Intn(3) {
		case 0:
			body = "tapply(" + nums(n, false) + ", " + f() + ", " + pick([]string{"sum", "mean", "length", "range", "max", "function(v) v"}) + ")"
		case 1:
			body = "tapply(" + nums(n, false) + ", list(" + f() + ", " + f() + "), " + pick([]string{"sum", "length", "max"}) + pick([]string{"", ", default = 0"}) + ")"
		default:
			body = "tapply(" + nums(n, false) + ", factor(" + f() + ", levels = c(\"a\", \"b\", \"c\", \"z\")), sum" + pick([]string{"", ", default = -1"}) + ")"
		}
	case 9:
		family = "rle"
		n := r.Intn(9)
		p := make([]string, n)
		for i := range p {
			p[i] = pick([]string{"1", "1", "2", "NA", "3"})
		}
		x := "c(" + strings.Join(p, ", ") + ")"
		if n == 0 {
			x = "numeric(0)"
		}
		body = pick([]string{"unclass(rle(" + x + "))", "inverse.rle(rle(" + x + "))", "identical(inverse.rle(rle(" + x + ")), " + x + ")"})
	case 10:
		family = "sequence"
		n := r.Intn(4)
		p := make([]string, n)
		for i := range p {
			p[i] = fmt.Sprint(r.Intn(5))
		}
		nv := "c(" + strings.Join(p, ", ") + ")"
		if n == 0 {
			nv = "integer(0)"
		}
		body = "sequence(" + nv + pick([]string{"", ", from = 2L", ", by = 2L", ", from = 5L, by = -1L"}) + ")"
	case 11:
		family = "helpers"
		switch r.Intn(6) {
		case 0:
			body = "replace(" + nums(5, false) + ", " + pick([]string{"2", "c(1, 5)", "7", "c(TRUE, FALSE)"}) + ", " + pick([]string{"0", "NA", "c(9, 8)"}) + ")"
		case 1:
			body = "setequal(" + nums(1+r.Intn(3), true) + ", " + nums(1+r.Intn(3), true) + ")"
		case 2:
			body = "zapsmall(c(" + pick([]string{"1", "100", "0.5"}) + ", " + pick([]string{"1e-10", "1e-3", "0", "NA"}) + ")" + pick([]string{"", ", digits = 3", ", digits = 8"}) + ")"
		case 3:
			body = "factorial(" + pick([]string{"0", "5", "c(1, 2, 3, 4)", "10", "NA"}) + ")"
		case 4:
			body = "c(" + pick([]string{"NULL", "1", "NA"}) + " %||% 5, " + pick([]string{"1", "3", "NA"}) + " %notin% c(1, 2))"
		default:
			body = "grepv(" + pick([]string{"\"a\"", "\"^b\"", "\"[0-9]\""}) + ", c(\"abc\", \"bcd\", \"x1\", \"a2\")" + pick([]string{"", ", invert = TRUE"}) + ")"
		}
	case 12:
		family = "replicate"
		body = pick([]string{
			"replicate(" + fmt.Sprint(r.Intn(4)) + ", 1:2)",
			"replicate(3, c(a = 1, b = 2))",
			"replicate(2, matrix(1:4, 2))",
			"replicate(2, \"s\", simplify = FALSE)",
			"sapply(1:3, function(i) matrix(i, 2, 2), simplify = \"array\")",
			"sapply(c(u = 1, v = 2), function(i) c(i, i + 1), simplify = \"array\")",
		})
	default:
		family = "bquote"
		body = pick([]string{
			"bquote(x + .(" + fmt.Sprint(r.Intn(9)) + "))",
			"y <- " + fmt.Sprint(r.Intn(9)) + "; bquote(f(.(y), z, k = .(y + 1)))",
			"bquote(g(a, .(c(1, 2))))",
			"bquote(h(..(c(1, 2)), 3), splice = TRUE)",
			"q <- quote(f(a, b = 2, c)); list(q[[1]], q[[3]], q[-1], length(q))",
		})
	}
	return "# " + family + "\n" + body
}
