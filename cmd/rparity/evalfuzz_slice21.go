package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// The slice21 mode: random programs over the additions of slice 21: ftable,
// ecdf, poly and predict.poly, read.fwf, the .rowSums family, gregexec,
// regmatches<-, formatDL, pairlist, the restart accessors and the methods of
// factor, noquote and table. The first line is "# family". A value ends in
// "# libm" when it is compared within --ulp units in the last place after being
// rounded to seven digits (poly is accurate to a stated bound, not bit
// identical to R).

func fuzzSlice21(r *rand.Rand) string {
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
	r7 := func(e string) string { return "round(" + e + ", 7) + 0" }
	var family, body string
	switch r.Intn(9) {
	case 0:
		family = "ftable"
		n := 4 + r.Intn(8)
		body = "d <- data.frame(a = " + strs(n, []string{"x", "y"}) + ", b = " + strs(n, []string{"p", "q", "s"}) + ", c = " + strs(n, []string{"u", "v"}) + "); t3 <- table(d); " + pick([]string{
			"ftable(t3)", "ftable(t3, row.vars = 1)", "ftable(t3, col.vars = c(1, 3))", "ftable(t3, row.vars = c('c', 'a'))", "as.matrix(ftable(t3))",
			"as.data.frame(ftable(t3))", "format(ftable(t3), quote = FALSE)", "ftable(d$a, d$b)", "ftable(a ~ b, d)", "ftable(a + b ~ c, d)", "dim(ftable(t3, row.vars = 2))",
			"attributes(ftable(t3))$row.vars", "as.table(ftable(t3))", "write.ftable(ftable(t3), quote = FALSE)",
		})
	case 1:
		family = "ecdf"
		n := 2 + r.Intn(12)
		body = "x <- " + nums(n, 9) + "; Fn <- ecdf(x); " + pick([]string{
			"Fn(c(-1, 0, 2.5, 4, 100))", "knots(Fn)", "summary(Fn)", "round(quantile(Fn, c(0.1, 0.5, 0.9)), 7)", "class(Fn)", "Fn(NA)", "environment(Fn)$nobs",
			"unclass(summary(Fn))[1:3]", "round(quantile(Fn, 0.25, type = 1), 7)", "Fn(sort(x))", "length(knots(Fn))",
		})
	case 2:
		family = "poly"
		n := 6 + r.Intn(8)
		deg := 1 + r.Intn(4)
		d := fmt.Sprint(deg)
		body = "x <- " + nums(n, 40) + " + (1:" + fmt.Sprint(n) + ") / 13; p <- poly(x, " + d + "); " + pick([]string{
			r7("unclass(p)[, 1]"), r7("crossprod(p)"), "attr(p, 'degree')", "dim(p)", "class(p)", r7("predict(p, c(0.5, 7.25))"), r7("poly(x, " + d + ", raw = TRUE)[, " + d + "]"),
			r7("attr(p, 'coefs')$alpha"), "signif(attr(p, 'coefs')$norm2, 8)", r7("poly(x, degree = " + d + ", simple = TRUE)[1:3, ]"),
			"dim(model.matrix(~ poly(x, " + d + ")))", "colnames(poly(x, " + d + "))", "tryCatch(poly(x, 40), error = function(e) conditionMessage(e))",
		}) + " # libm"
	case 3:
		family = "readfwf"
		rows := make([]string, 1+r.Intn(5))
		for i := range rows {
			rows[i] = fmt.Sprintf("%q", fmt.Sprint(10000+r.Intn(89999))+pick([]string{"ab", "xy", "q", ""})+fmt.Sprint(r.Intn(90)))
		}
		w := pick([]string{"c(2, 3, 2)", "c(1, -2, 3)", "c(5, 2, 3)", "c(2, 2, 2, 2)", "c(10, 1)", "c(3, -1, 4)"})
		body = "tf <- tempfile(); writeLines(c(" + strings.Join(rows, ", ") + "), tf); " + pick([]string{
			"read.fwf(tf, " + w + ")", "read.fwf(tf, " + w + ", header = TRUE)", "read.fwf(tf, " + w + ", skip = 1)", "read.fwf(tf, " + w + ", n = 2)",
			"read.fwf(tf, " + w + ", col.names = c('u', 'v', 'w', 'z')[seq_len(sum(" + w + " > 0))])", "sapply(read.fwf(tf, " + w + "), class)",
			"read.fwf(tf, " + w + ", colClasses = 'character')", "dim(read.fwf(tf, " + w + ", sep = ','))",
		})
	case 4:
		family = "rowcol"
		nr, nc := 2+r.Intn(4), 2+r.Intn(4)
		body = "m <- matrix(" + nums(nr*nc, 9) + ", " + fmt.Sprint(nr) + "); " + pick([]string{
			fmt.Sprintf(".rowSums(m, %d, %d)", nr, nc), fmt.Sprintf(".colSums(m, %d, %d)", nr, nc), fmt.Sprintf(".rowMeans(m, %d, %d)", nr, nc),
			fmt.Sprintf(".colMeans(m, %d, %d)", nr, nc), fmt.Sprintf(".row(c(%d, %d))", nr, nc), fmt.Sprintf(".col(c(%d, %d))", nr, nc),
			fmt.Sprintf(".rowSums(m, %d, %d, na.rm = TRUE)", nr, nc), "identical(.row(dim(m)), row(m))", "identical(.col(dim(m)), col(m))",
		})
	case 5:
		family = "regex"
		body = "x <- " + strs(2+r.Intn(3), []string{"abc123", "x9y8", "no digits", "7", "a1b22c333"}) + "; " + pick([]string{
			"regmatches(x, gregexec('([a-z])([0-9]+)', x))", "m <- regexpr('[0-9]+', x); regmatches(x, m) <- 'N'; x", "m <- gregexpr('[0-9]+', x); regmatches(x, m) <- list('#'); x",
			"m <- gregexpr('[0-9]', x); regmatches(x, m, invert = TRUE)", "gregexec('[0-9]+', x)", "regexec('([a-z]+)([0-9]*)', x)",
			"m <- regexpr('[a-z]', x); regmatches(x, m, invert = NA)",
		})
	case 6:
		family = "formatDL"
		n := 1 + r.Intn(4)
		x := strs(n, []string{"alpha", "b", "gamma ray", "dd"})
		y := strs(n, []string{"first description here", "second", "a longer line of text that wraps around the margin of the page"})
		body = pick([]string{
			"formatDL(" + x + ", " + y + ")", "formatDL(" + x + ", " + y + ", style = 'list')", "formatDL(" + x + ", " + y + ", width = 30)",
			"formatDL(" + x + ", " + y + ", indent = 5)", "formatDL(structure(" + y + ", names = " + x + "))",
		})
	case 7:
		family = "listtypes"
		body = pick([]string{
			"typeof(pairlist(a = 1, b = 2))", "length(pairlist(1, 2, 3))", "is.null(pairlist())", "names(as.pairlist(list(p = 1, q = 'a')))", "is.pairlist(pairlist(a = 1))",
			"class(noquote(letters[1:3])[2])", "class(c(noquote('a'), 'b'))", "unclass(factor(c('a', 'b', 'a'))[[2]])", "as.logical(factor(c('TRUE', 'FALSE', 'TRUE')))",
			"as.vector(factor(c('u', 'v')))", "x <- structure(list(a = 1, b = 2, c = 3), class = 'listof'); class(x[2:3])",
			"as.vector(table(c('a', 'b', 'a'))[2])", "x <- structure(list(1, 2, 3), class = 'simple.list'); length(x[2:3])",
			"all.equal(y ~ x, y ~ x)", "all.equal(y ~ x, y ~ z)", "d <- data.frame(a = c(1, NA, 3), b = 1:3); anyNA(d)", "d <- data.frame(a = c(1, 1, 3), b = c(2, 2, 3)); anyDuplicated(d)",
			"as.vector(data.frame(a = 1:2, b = 3:4), 'list')", "merge(data.frame(k = 1:3, a = 4:6), data.frame(k = c(2, 3, 4), b = 7:9))",
			"d <- data.frame(g = c('a', 'b', 'a'), v = 1:3); split(d$v, d$g) <- list(c(10L, 30L), 20L); d$v",
			"chooseOpsMethod(1, 2, NULL, NULL, 'x', FALSE)", "is.null(nameOfClass(1))",
		})
	default:
		family = "restarts"
		body = pick([]string{
			"withRestarts(computeRestarts()[[1]]$name, myR = function() 1)", "withRestarts(isRestart(computeRestarts()[[1]]), myR = function() 1)",
			"withRestarts(restartDescription(computeRestarts()[[1]]), myR = list(handler = function() 1, description = 'does it'))",
			"withRestarts(invokeRestart('myR', 5), myR = function(v) v * 2)", "withRestarts(is.null(restartFormals(computeRestarts()[[1]])), myR = function() 1)",
			"tryCatch(invokeRestart('nope'), error = function(e) conditionMessage(e))", "class(messageCondition('hello'))", "conditionMessage(messageCondition('hi', call = NULL))",
			"class(errorCondition('m', class = 'my'))", "class(warningCondition('m', class = 'my'))", "withAutoprint({ x <- 1:3; x * 2 })",
			"is.function(print.function)", "length(sys.status())",
			"tryCatch(stop('boom'), error = function(e) class(e))",
			"suppressPackageStartupMessages(packageStartupMessage('hi'))", "withCallingHandlers(packageStartupMessage('hi'), packageStartupMessage = function(m) invokeRestart('muffleMessage'))",
		})
	}
	return "# " + family + "\n" + body
}
