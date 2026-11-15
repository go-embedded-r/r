package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

// The slice16 mode: random programs over the closures of slice 16: pretty,
// strwrap, version numbers, the rounding, truncating, cutting and sequencing of
// dates and times, the summaries of dates, the marks of format, formatC and
// prettyNum, the options of deparse, by, unique and duplicated of matrices and
// kronecker. The first line is "# family". Everything is compared exactly,
// printed text included; text that R builds with typographic quotes is printed
// with cat, because R leaves such strings unmarked and ours marks them.

// fuzzNum is a number of a random magnitude, written as R source.
func fuzzNum(r *rand.Rand) string {
	switch r.Intn(5) {
	case 0:
		return fmt.Sprint(r.Intn(200) - 100)
	case 1:
		return fmt.Sprintf("%.3g", (r.Float64()-0.5)*math.Pow(10, float64(r.Intn(9)-4)))
	case 2:
		return fmt.Sprintf("%.2g", r.Float64()*math.Pow(10, float64(r.Intn(13)-6)))
	case 3:
		return fmt.Sprintf("%d", r.Intn(100000)*(1+r.Intn(1000)))
	default:
		return fmt.Sprintf("%.6g", r.NormFloat64()*math.Pow(10, float64(r.Intn(7)-3)))
	}
}

func fuzzSlice16(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	var family, body string
	switch r.Intn(11) {
	case 0:
		family = "pretty"
		args := ""
		if r.Intn(2) == 0 {
			args += fmt.Sprintf(", n = %d", r.Intn(15))
		}
		if r.Intn(4) == 0 {
			args += fmt.Sprintf(", min.n = %d", r.Intn(5))
		}
		if r.Intn(5) == 0 {
			args += ", high.u.bias = " + pick([]string{"0.5", "1", "3", "5", "10"})
		}
		if r.Intn(5) == 0 {
			args += ", u5.bias = " + pick([]string{"0", "0.5", "2", "6"})
		}
		if r.Intn(6) == 0 {
			args += ", shrink.sml = " + pick([]string{"0.1", "0.5", "1"})
		}
		a, b := fuzzNum(r), fuzzNum(r)
		if r.Intn(6) == 0 {
			b = a
		}
		body = "print(pretty(c(" + a + ", " + b + ")" + args + "))"
	case 1:
		family = "strwrap"
		words := []string{"the", "quick", "brown", "fox", "jumps.", "over", "lazy", "dog!", "It", "was", "a", "sunny", "day?", "and", "fine", "x", "ab", "ünï", "日本語", "(see.)", "done.\\\""}
		seps := []string{" ", " ", " ", "  ", "   ", "\\t", "\\n", "\\n\\n"}
		var sb strings.Builder
		for i, n := 0, 1+r.Intn(25); i < n; i++ {
			if i > 0 {
				sb.WriteString(pick(seps))
			}
			sb.WriteString(pick(words))
		}
		args := fmt.Sprintf(", width = %d", 1+r.Intn(40))
		if r.Intn(3) == 0 {
			args += fmt.Sprintf(", indent = %d", r.Intn(4))
		}
		if r.Intn(3) == 0 {
			args += fmt.Sprintf(", exdent = %d", r.Intn(4))
		}
		if r.Intn(3) == 0 {
			args += ", prefix = " + pick([]string{"\"> \"", "\"#\"", "\"\""})
		}
		if r.Intn(4) == 0 {
			args += ", initial = \"* \""
		}
		body = "print(strwrap(\"" + sb.String() + "\"" + args + "))"
	case 2:
		family = "version"
		vs := func() string {
			parts := make([]string, 1+r.Intn(3))
			for i := range parts {
				parts[i] = fmt.Sprint(r.Intn(12))
			}
			return strings.Join(parts, ".")
		}
		v := make([]string, 2+r.Intn(4))
		for i := range v {
			v[i] = "\"" + vs() + "\""
		}
		x := "numeric_version(c(" + strings.Join(v, ", ") + "))"
		body = "x <- " + x + "; " + pick([]string{"print(x)", "print(sort(x))", "print(max(x))", "print(range(x))", "print(x > \"3.4\")", "print(x == x[1])", "print(unique(x))", "print(order(x))", "print(rev(x))", "print(format(x))", "print(x[[1]])", "print(as.list(x)[1])", "print(compareVersion(" + v[0] + ", " + v[1] + "))", "print(c(x, x)[2])"})
	case 3:
		family = "dt-trunc"
		t := fmt.Sprintf("as.POSIXct(\"%04d-%02d-%02d %02d:%02d:%02d\", tz = \"%s\")", 1990+r.Intn(60), 1+r.Intn(12), 1+r.Intn(28), r.Intn(24), r.Intn(60), r.Intn(60), pick([]string{"UTC", "Europe/Paris", "America/New_York"}))
		u := pick([]string{"secs", "mins", "hours", "days", "months", "years"})
		body = pick([]string{"print(trunc(" + t + ", \"" + u + "\"))", "print(round(" + t + ", \"" + u + "\"))", "print(format(round(" + t + ", \"" + u + "\"), \"%Y-%m-%d %H:%M:%S\"))"})
	case 4:
		family = "dt-date"
		d := fmt.Sprintf("as.Date(\"%04d-%02d-%02d\")", 1950+r.Intn(100), 1+r.Intn(12), 1+r.Intn(28))
		body = pick([]string{"print(trunc(" + d + ", \"months\"))", "print(trunc(" + d + ", \"years\"))", "print(trunc(" + d + "))", "print(round(" + d + " + 0.5))"})
	case 5:
		family = "dt-cut"
		n := 2 + r.Intn(8)
		ds := make([]string, n)
		for i := range ds {
			ds[i] = fmt.Sprintf("\"2024-%02d-%02d\"", 1+r.Intn(12), 1+r.Intn(28))
		}
		d := "as.Date(c(" + strings.Join(ds, ", ") + "))"
		by := pick([]string{"\"day\"", "\"week\"", "\"month\"", "\"2 months\"", "\"quarter\"", "\"year\"", "3"})
		body = "d <- " + d + "; print(cut(d, " + by + pick([]string{"", ", right = TRUE", ", start.on.monday = FALSE", ", labels = FALSE"}) + "))"
	case 6:
		family = "dt-seq"
		t := fmt.Sprintf("as.POSIXct(\"2024-%02d-%02d %02d:%02d:00\", tz = \"%s\")", 1+r.Intn(12), 1+r.Intn(28), r.Intn(24), r.Intn(60), pick([]string{"UTC", "Europe/Paris"}))
		by := pick([]string{"\"day\"", "\"month\"", "\"2 months\"", "\"year\"", "\"week\"", "\"hour\"", "\"30 min\"", "\"quarter\"", "\"DSTday\"", "86400", "\"-1 day\"", "\"-1 month\""})
		body = "print(seq(" + t + ", by = " + by + ", length.out = " + fmt.Sprint(1+r.Intn(8)) + "))"
	case 7:
		family = "dt-summary"
		n := 2 + r.Intn(8)
		ds := make([]string, n)
		for i := range ds {
			ds[i] = fmt.Sprintf("\"20%02d-%02d-%02d\"", 10+r.Intn(15), 1+r.Intn(12), 1+r.Intn(28))
		}
		d := "as.Date(c(" + strings.Join(ds, ", ") + pick([]string{"", ", NA"}) + "))"
		body = "d <- " + d + "; " + pick([]string{"print(summary(d))", "print(summary(as.POSIXct(d)))", "print(summary(d, digits = 4))", "print(class(summary(d)))", "dput(summary(d))"})
	case 8:
		family = "marks"
		x := "c(" + fuzzNum(r) + ", " + fuzzNum(r) + ", " + pick([]string{"0", "1234567.891", "0.5"}) + ")"
		opts := pick([]string{"big.mark = \",\"", "big.mark = \" \", big.interval = 2", "decimal.mark = \",\"", "drop0trailing = TRUE", "zero.print = \"-\"", "small.mark = \" \", small.interval = 3", "big.mark = \".\", decimal.mark = \",\""})
		if strings.Contains(opts, "small.mark") {
			// R groups the digits of a mantissa and of its exponent as if they
			// were decimals; only numbers that print in fixed notation are compared
			x = "c(1234.56789, 0.123456789, 7)"
		}
		switch r.Intn(3) {
		case 0:
			body = "print(format(" + x + ", " + opts + pick([]string{"", ", nsmall = 2", ", digits = 4"}) + "))"
		case 1:
			body = "print(formatC(" + x + ", format = " + pick([]string{"\"f\"", "\"e\"", "\"g\""}) + ", digits = " + fmt.Sprint(r.Intn(6)) + ", " + opts + "))"
		default:
			body = "print(prettyNum(" + x + ", " + opts + "))"
		}
	case 9:
		family = "deparse"
		x := pick([]string{"c(a = 1L, b = NA)", "c(1.5, NA)", "list(p = 1L, q = \"s\")", "factor(c(\"a\", \"b\"))", "1:10", "c(0.1, 1/3)", "NA_character_", "quote(x + y)", "matrix(1:4, 2)", "as.Date(\"2024-01-02\")", "c(a = 1, b = 2)", "letters[1:14]"})
		ctl := pick([]string{"NULL", "\"keepInteger\"", "c(\"keepNA\", \"niceNames\")", "c(\"showAttributes\", \"niceNames\")", "\"showAttributes\"", "\"niceNames\"", "\"digits17\"", "\"all\"", "\"quoteExpressions\"", "c(\"keepNA\", \"keepInteger\", \"niceNames\", \"showAttributes\")"})
		body = pick([]string{"print(deparse(" + x + ", control = " + ctl + "))", "print(deparse(" + x + ", width.cutoff = " + fmt.Sprint(20+r.Intn(60)) + "))", "dput(" + x + ")", "dput(" + x + ", control = " + ctl + ")"})
	default:
		family = "by-unique"
		body = pick([]string{
			"df <- data.frame(g = c(\"a\", \"b\", \"a\", \"c\"), x = 1:4); print(by(df, df$g, function(d) sum(d$x)))",
			"df <- data.frame(g = c(\"a\", \"b\", \"a\", \"c\"), x = 1:4); print(by(df$x, df$g, mean))",
			"m <- matrix(c(1, 2, 1, 2, 3, 4), 3, byrow = TRUE); print(unique(m)); print(duplicated(m)); print(anyDuplicated(m))",
			"m <- matrix(c(1, 1, 2, 2, 1, 1), 2); print(unique(m, MARGIN = 2)); print(duplicated(m, MARGIN = 2))",
			"A <- matrix(1:4, 2); B <- matrix(c(1, 0, 0, 1), 2); print(A %x% B)",
			"print(1:" + fmt.Sprint(2+r.Intn(3)) + " %x% 1:" + fmt.Sprint(2+r.Intn(3)) + ")",
			"print(as.table(matrix(1:6, 2)))",
			"print(data.matrix(data.frame(a = 1:3, b = c(\"x\", \"y\", \"x\"))))",
			"print(provideDimnames(matrix(1:6, 2)))",
		})
	}
	return "# " + family + "\n" + body
}
