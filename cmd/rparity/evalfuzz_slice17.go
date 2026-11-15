package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// The slice17 mode: random programs over the closures of slice 17: the terms
// of a formula (labels, order, factors, intercept), update and reformulate,
// all.vars, model frames with a subset and a way of dealing with NA, aggregate
// by one or two groups and with a formula, stack and unstack, rowsum,
// sort_by of a formula, array2DF, the subsetting and replacement of POSIXlt,
// and str of vectors. The first line is "# family". A formula is printed
// through deparse() and the labels of its terms, never as an object: R prints
// the environment of one made in a function.

func fuzzData(r *rand.Rand) string {
	n := 3 + r.Intn(8)
	gs := []string{"a", "b", "c", "d"}
	g, h, v, w := make([]string, n), make([]string, n), make([]string, n), make([]string, n)
	for i := 0; i < n; i++ {
		g[i] = "\"" + gs[r.Intn(len(gs))] + "\""
		h[i] = "\"" + []string{"x", "y"}[r.Intn(2)] + "\""
		v[i] = fmt.Sprint(r.Intn(20))
		if r.Intn(6) == 0 {
			w[i] = "NA"
		} else {
			w[i] = fmt.Sprintf("%.1f", float64(r.Intn(100))/4)
		}
	}
	return "d <- data.frame(g = c(" + strings.Join(g, ", ") + "), h = c(" + strings.Join(h, ", ") + "), v = c(" + strings.Join(v, ", ") + "), w = c(" + strings.Join(w, ", ") + "), stringsAsFactors = FALSE); "
}

func fuzzSlice17(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	var family, body string
	switch r.Intn(9) {
	case 0:
		family = "terms"
		vars := []string{"a", "b", "c", "d"}
		term := func() string {
			x, y := pick(vars), pick(vars)
			switch r.Intn(9) {
			case 0:
				return x
			case 1:
				return x + ":" + y
			case 2:
				return x + "*" + y
			case 3:
				return "(" + x + "+" + y + ")^2"
			case 4:
				return x + "/" + y
			case 5:
				return "I(" + x + "^2)"
			case 6:
				return "log(" + x + ")"
			case 7:
				return x + "%in%" + y
			}
			return "(" + x + "+" + y + ")"
		}
		f := "y ~ " + term()
		for i, k := 0, r.Intn(3); i < k; i++ {
			f += pick([]string{" + ", " - ", " + "}) + term()
		}
		if r.Intn(4) == 0 {
			f += pick([]string{" - 1", " + 0", " - 0", " + 1"})
		}
		opt := pick([]string{"", ", keep.order = TRUE", ", simplify = TRUE", ", specials = \"log\""})
		body = "t <- terms(" + f + opt + "); print(attr(t, \"term.labels\")); print(attr(t, \"order\")); print(attr(t, \"intercept\")); print(attr(t, \"factors\")); print(deparse(t))"
	case 1:
		family = "update"
		old := pick([]string{"y ~ x + z", "y ~ x * z", "log(y) ~ x", "~ x + z", "y ~ a + b + c"})
		upd := pick([]string{". ~ . + w", ". ~ . - x", "log(.) ~ .", ". ~ .^2", "~ . + k", ". ~ . * q", "z ~ ."})
		body = "print(deparse(update(" + old + ", " + upd + ")))"
	case 2:
		family = "vars"
		e := pick([]string{"y ~ x + log(z)", "a + b * c", "f(x, y)[z]", "if (p) q else r", "x$y + w", "function(a) a + b"})
		body = "print(all.vars(quote(" + e + "))); print(all.names(quote(" + e + "))); print(all.vars(quote(" + e + "), unique = FALSE))"
	case 3:
		family = "aggregate"
		fun := pick([]string{"sum", "mean", "length", "max", "min"})
		by := pick([]string{"list(d$g)", "list(g = d$g)", "list(g = d$g, h = d$h)", "list(h = d$h, g = d$g)"})
		f := pick([]string{"v ~ g", "v ~ g + h", "cbind(v, w) ~ g", "w ~ h", "v ~ ."})
		body = fuzzData(r) + pick([]string{
			"print(aggregate(d[c(\"v\", \"w\")], " + by + ", " + fun + "))",
			"print(aggregate(d$v, " + by + ", " + fun + "))",
			"print(aggregate(d[c(\"v\", \"w\")], " + by + ", " + fun + ", na.rm = TRUE))",
			"print(aggregate(" + f + ", d[c(\"g\", \"h\", \"v\", \"w\")], " + fun + "))",
			"print(aggregate(" + f + ", d[c(\"g\", \"h\", \"v\", \"w\")], " + fun + ", na.action = na.pass))",
			"print(aggregate(d[\"v\"], " + by + ", " + fun + ", drop = FALSE))",
		})
	case 4:
		family = "stack"
		body = pick([]string{
			"print(stack(list(a = 1:" + fmt.Sprint(1+r.Intn(4)) + ", b = " + fmt.Sprint(r.Intn(9)) + ")))",
			"s <- stack(data.frame(x = 1:" + fmt.Sprint(2+r.Intn(3)) + ", y = 5:" + fmt.Sprint(6+r.Intn(3)) + ")); print(s); print(unstack(s))",
			"print(unstack(data.frame(v = 1:6, f = factor(rep(c(\"p\", \"q\", \"r\"), 2))), v ~ f))",
			"print(stack(list(u = 1, w = 2), drop = TRUE))",
			"print(levels(stack(list(zz = 1, aa = 2))$ind))",
		})
	case 5:
		family = "rowsum"
		n := 3 + r.Intn(6)
		xs, gs := make([]string, n), make([]string, n)
		for i := range xs {
			xs[i] = fmt.Sprint(r.Intn(30))
			gs[i] = "\"" + []string{"u", "v", "w"}[r.Intn(3)] + "\""
		}
		body = "x <- c(" + strings.Join(xs, ", ") + "); g <- c(" + strings.Join(gs, ", ") + "); " + pick([]string{
			"print(rowsum(x, g))",
			"print(rowsum(x, g, reorder = FALSE))",
			"print(rowsum(matrix(x, ncol = 1), g))",
			"print(rowsum(data.frame(p = x, q = x * 2), g))",
			"print(rowsum(as.numeric(x), g))",
		})
	case 6:
		family = "frame"
		body = "P <- function(x) cat(capture.output(print(x)), sep = \"\\n\"); " + fuzzData(r) + pick([]string{
			"mf <- model.frame(v ~ w + g, d, na.action = na.omit); P(mf); print(attr(attr(mf, \"terms\"), \"dataClasses\"))",
			"P(model.frame(v ~ g, d, subset = v > 5, na.action = na.pass))",
			"print(model.response(model.frame(v ~ w, d, na.action = na.exclude)))",
			"print(get_all_vars(v ~ w, d))",
			"print(sort_by(d, ~ g + v))",
			"print(sort_by(d, ~ v, decreasing = TRUE))",
			"print(complete.cases(d))",
			"print(na.omit(d))",
		})
	case 7:
		family = "posixlt"
		x := "x <- as.POSIXlt(c(\"2024-01-15 10:30:00\", \"2024-02-20 11:45:30\", \"2024-03-25 23:59:59\", \"2023-12-31 00:00:00\"), tz = \"UTC\"); "
		body = x + pick([]string{
			"print(x[" + fmt.Sprint(1+r.Intn(4)) + "])",
			"print(x[c(" + fmt.Sprint(1+r.Intn(4)) + ", " + fmt.Sprint(1+r.Intn(4)) + ")])",
			"print(x[-" + fmt.Sprint(1+r.Intn(4)) + "])",
			"print(x[[" + fmt.Sprint(1+r.Intn(4)) + "]])",
			"y <- x; y[" + fmt.Sprint(1+r.Intn(4)) + "] <- x[" + fmt.Sprint(1+r.Intn(4)) + "]; print(y)",
			"y <- x; y$hour <- " + fmt.Sprint(r.Intn(24)) + "L; print(format(y))",
			"y <- x; y$year <- y$year + " + fmt.Sprint(r.Intn(5)) + "L; print(format(y))",
			"print(length(x[" + fmt.Sprint(r.Intn(8)) + "]))",
			"print(rev(x))",
			"print(c(x[1], x[2]))",
		})
	default:
		family = "array2DF"
		body = pick([]string{
			"print(array2DF(array(1:" + fmt.Sprint(4) + ", c(2, 2), dimnames = list(a = c(\"x\", \"y\"), b = c(\"p\", \"q\")))))",
			"print(array2DF(matrix(1:6, 2)))",
			"print(array2DF(table(c(\"a\", \"b\", \"a\", \"c\"))))",
			"print(as.data.frame(matrix(1:6, 3)))",
			"print(as.data.frame(array(1:" + fmt.Sprint(8) + ", c(2, 2, 2))))",
			"str(" + pick([]string{"1:" + fmt.Sprint(1+r.Intn(30)), "c(1.5, 22.25, 333.125)", "letters[1:" + fmt.Sprint(1+r.Intn(26)) + "]", "factor(letters[1:" + fmt.Sprint(1+r.Intn(8)) + "])", "c(TRUE, FALSE, NA)", "as.Date(\"2024-01-01\") + 0:" + fmt.Sprint(r.Intn(10)), "matrix(1:" + fmt.Sprint(2*(1+r.Intn(6))) + ", 2)"}) + ")",
		})
	}
	return "# " + family + "\n" + body
}
