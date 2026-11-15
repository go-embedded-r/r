package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// The slice18 mode: random programs over the readers and writers of slice 18
// (scan, read.table and its relatives, count.fields, type.convert,
// write.table, write.csv), xtabs and summary.table, and the descriptive
// statistics (var, sd, cov, cor in its three methods, median, quantile of the
// nine types, fivenum, IQR, mad, weighted.mean). The first line is "# family".
// A statistic ends in "# libm": it is compared within --ulp units in the last
// place (the kernels are accurate to a stated bound, GNU R's are not exact),
// and counted; the readers and writers are compared exactly.

func fuzzCell(r *rand.Rand, sep string) string {
	pool := []string{"1", "2", "3", "10", "-4", "1.5", "2.25", "0.5", "NA", "", "TRUE", "FALSE", "abc", "x y", "'q r'", "\"p q\"", "\"a" + sep + "b\"", "1e3", "007", "Inf", "T", "F", "NaN", "0x10", " 5 "}
	return pool[r.Intn(len(pool))]
}

func fuzzTable(r *rand.Rand, sep string, header bool) string {
	nc := 1 + r.Intn(4)
	nr := 1 + r.Intn(5)
	var lines []string
	if header {
		h := make([]string, nc)
		for i := range h {
			h[i] = []string{"a", "b", "c", "d", "a b", "1x", "x.y"}[r.Intn(7)]
		}
		lines = append(lines, strings.Join(h, sep))
	}
	for i := 0; i < nr; i++ {
		row := make([]string, nc)
		for j := range row {
			row[j] = fuzzCell(r, sep)
			if sep == "" || sep == " " {
				row[j] = strings.ReplaceAll(strings.ReplaceAll(row[j], " ", ""), "x y", "xy")
				if row[j] == "" {
					row[j] = "7"
				}
			}
		}
		if r.Intn(8) == 0 && nc > 1 {
			row = row[:nc-1]
		}
		lines = append(lines, strings.Join(row, sep))
		if r.Intn(10) == 0 {
			lines = append(lines, "")
		}
		if r.Intn(12) == 0 {
			lines = append(lines, "# note")
		}
	}
	return strings.Join(lines, "\n")
}

func fuzzVec18(r *rand.Rand, withNA bool) string {
	n := 2 + r.Intn(9)
	v := make([]string, n)
	for i := range v {
		switch {
		case withNA && r.Intn(7) == 0:
			v[i] = "NA"
		case r.Intn(3) == 0:
			v[i] = fmt.Sprintf("%.2f", float64(r.Intn(400))/8-20)
		default:
			v[i] = fmt.Sprint(r.Intn(25) - 5)
		}
	}
	return "c(" + strings.Join(v, ", ") + ")"
}

func fuzzSlice18(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	q := func(s string) string { return fmt.Sprintf("%q", s) }
	var family, body string
	switch r.Intn(9) {
	case 0:
		family = "scan"
		sep := pick([]string{"", ",", ";", "\t", " "})
		txt := fuzzTable(r, sep+"", false)
		if sep == "" {
			txt = strings.ReplaceAll(txt, "\"a"+sep+"b\"", "ab")
		}
		what := pick([]string{"0", "\"\"", "list(0, \"\")", "list(\"\", \"\", \"\")", "0L", "TRUE", "list(a = \"\", b = \"\")"})
		opts := pick([]string{"", ", fill = TRUE", ", flush = TRUE", ", strip.white = TRUE", ", blank.lines.skip = FALSE", ", comment.char = \"#\"", ", na.strings = c(\"NA\", \"abc\")", ", nlines = 2", ", skip = 1", ", multi.line = FALSE", ", dec = \",\""})
		sp := ""
		if sep != "" {
			sp = ", sep = " + q(sep)
		}
		body = "r <- tryCatch(scan(text = " + q(txt) + ", what = " + what + sp + opts + ", quiet = TRUE), error = function(e) conditionMessage(e)); r"
	case 1:
		family = "readtable"
		sep := pick([]string{"", ",", ";", "\t"})
		header := r.Intn(2) == 0
		txt := fuzzTable(r, sep+"", header)
		fn := "read.table"
		args := ", header = " + map[bool]string{true: "TRUE", false: "FALSE"}[header]
		if sep != "" {
			args += ", sep = " + q(sep)
		}
		switch r.Intn(4) {
		case 0:
			fn, args = "read.csv", ""
			txt = fuzzTable(r, ",", true)
		case 1:
			args += ", fill = TRUE"
		case 2:
			args += pick([]string{", stringsAsFactors = TRUE", ", na.strings = c(\"NA\", \"abc\", \"\")", ", colClasses = \"character\"", ", row.names = 1", ", nrows = 2", ", skip = 1", ", check.names = FALSE", ", strip.white = TRUE", ", comment.char = \"\"", ", tryLogical = FALSE"})
		}
		body = "d <- tryCatch(" + fn + "(text = " + q(txt) + args + "), error = function(e) conditionMessage(e), warning = function(w) conditionMessage(w)); d"
	case 2:
		family = "countfields"
		sep := pick([]string{"", ",", ";"})
		txt := fuzzTable(r, sep+"", false)
		sp := ""
		if sep != "" {
			sp = ", sep = " + q(sep)
		}
		body = "tf <- tempfile(); writeLines(" + q(txt) + ", tf); r <- tryCatch(count.fields(tf" + sp + pick([]string{"", ", blank.lines.skip = FALSE", ", comment.char = \"\"", ", quote = \"\"", ", skip = 1"}) + "), error = function(e) conditionMessage(e)); r"
	case 3:
		family = "typeconvert"
		n := 1 + r.Intn(6)
		cells := make([]string, n)
		for i := range cells {
			cells[i] = q(fuzzCell(r, ","))
		}
		body = "type.convert(c(" + strings.Join(cells, ", ") + "), as.is = " + pick([]string{"TRUE", "FALSE"}) + pick([]string{"", ", na.strings = c(\"NA\", \"\")", ", tryLogical = FALSE", ", dec = \",\"", ", numerals = \"no.loss\""}) + ")"
	case 4:
		family = "write"
		d := "d <- data.frame(n = " + fuzzVec18(r, true) + ", s = c(\"a\", \"b\\\"c\", NA, \"x y\", \"\", \"é\", \"p,q\", \"r\", \"s\", \"t\", \"u\", \"v\")[1:2], stringsAsFactors = FALSE); "
		_ = d
		n := 2 + r.Intn(4)
		d = fmt.Sprintf("n <- %d; d <- data.frame(n = %s[1:n], s = c(\"a\", \"b\\\"c\", NA, \"x y\", \"\", \"p,q\")[1:n], l = c(TRUE, NA, FALSE, TRUE, NA, FALSE)[1:n], f = factor(c(\"u\", \"v\", \"u\", NA, \"v\", \"u\")[1:n])); ", n, "c(1.5, NA, -3, 100000, 0.1, 1e-5, 7)")
		body = d + pick([]string{
			"write.table(d)", "write.csv(d)", "write.csv2(d, row.names = FALSE)", "write.table(d, quote = FALSE, sep = \"\\t\", na = \"-\")",
			"write.table(d, col.names = NA, qmethod = \"double\")", "write.table(d, row.names = FALSE, col.names = FALSE, sep = \";\", dec = \",\")",
			"write.csv(d, quote = 2:3, na = \"\")", "write.table(as.matrix(d[, c(1, 3)]))",
			"f <- tempfile(); write.csv(d, f, row.names = FALSE); r <- read.csv(f); print(r); print(sapply(r, class))",
		})
	case 5:
		family = "stats"
		x, y := fuzzVec18(r, true), fuzzVec18(r, true)
		body = "x <- " + x + "; y <- " + y + "; n <- min(length(x), length(y)); x <- x[1:n]; y <- y[1:n]; round(" + pick([]string{
			"var(x, na.rm = TRUE)", "sd(x, na.rm = TRUE)", "var(x)", "median(x)", "median(x, na.rm = TRUE)", "fivenum(x)", "IQR(x, na.rm = TRUE)",
			"mad(x, na.rm = TRUE)", "weighted.mean(x, abs(y) + 1, na.rm = TRUE)",
			"weighted.mean(x, abs(y) + 1)", "cov(x, y, use = \"complete.obs\")", "cor(x, y, use = \"pairwise.complete.obs\")",
			"cor(x, y, method = \"spearman\", use = \"complete.obs\")", "cor(x, y, method = \"kendall\", use = \"complete.obs\")",
			"cov(cbind(x, y), use = \"pairwise.complete.obs\")", "cor(cbind(a = x, b = y, c = x + y), use = \"complete.obs\")", "var(cbind(x, y), na.rm = TRUE)",
		}) + ", 7) # libm"
	case 6:
		family = "quantile"
		x := fuzzVec18(r, true)
		probs := pick([]string{"c(0, .25, .5, .75, 1)", "c(.1, .33, .9)", "seq(0, 1, 0.1)", "c(0.025, 0.975)", "1/3", "c(0.3, 0.7)"})
		body = "x <- " + x + "; round(quantile(x, " + probs + ", type = " + fmt.Sprint(1+r.Intn(9)) + ", na.rm = TRUE" + pick([]string{"", ", names = FALSE"}) + "), 9) # libm"
	case 7:
		family = "xtabs"
		n := 3 + r.Intn(8)
		g, h, v := make([]string, n), make([]string, n), make([]string, n)
		for i := 0; i < n; i++ {
			g[i], h[i] = q([]string{"a", "b", "c"}[r.Intn(3)]), q([]string{"x", "y"}[r.Intn(2)])
			v[i] = fmt.Sprint(r.Intn(9))
		}
		body = "d <- data.frame(g = c(" + strings.Join(g, ", ") + "), h = c(" + strings.Join(h, ", ") + "), v = c(" + strings.Join(v, ", ") + ")); " + pick([]string{
			"print(xtabs(~ g, d))", "print(xtabs(~ g + h, d))", "print(xtabs(v ~ g + h, d))", "x <- xtabs(v ~ g, d); print(x); print(summary(x))", "print(xtabs(~ g + h, d, subset = v > 2))",
			"print(summary(xtabs(~ g + h, d))$n.cases)",
		})
	default:
		family = "typeconvert2"
		body = "d <- data.frame(a = c(\"1\", \"2\", \"x\"), b = c(\"1.5\", \"2\", \"3\"), c = c(\"T\", \"F\", \"T\"), stringsAsFactors = FALSE); print(type.convert(d, as.is = TRUE)); print(sapply(type.convert(d, as.is = FALSE), class))"
	}
	return "# " + family + "\n" + body
}
