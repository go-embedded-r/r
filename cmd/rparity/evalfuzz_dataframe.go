package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// The dataframe mode: random programs over data frames and tables: building,
// printing, indexing, assignment, binding, merge, subset, transform, order,
// split, table, cut, expand.grid and interaction. The first line is "# family".
// Everything is compared exactly, printed text included. There is no exclusion class: a
// case that differs is a defect or a divergence listed in the README.

func dfCol(r *rand.Rand, n int, kind int) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	elems := make([]string, n)
	for i := range elems {
		switch kind {
		case 0:
			elems[i] = pick([]string{"1L", "2L", "3L", "10L", "-4L", "NA_integer_", "0L"})
		case 1:
			elems[i] = pick([]string{"1.5", "2", "-0.25", "1e5", "123.456", "NA", "0", "1e-3"})
		case 2:
			elems[i] = pick([]string{"\"a\"", "\"bb\"", "\"c\"", "\"a\"", "NA", "\"dd d\"", "\"\""})
		default:
			elems[i] = pick([]string{"TRUE", "FALSE", "NA"})
		}
	}
	return "c(" + strings.Join(elems, ", ") + ")"
}

func dfSrc(r *rand.Rand, n int, withFactor bool) string {
	cols := []string{"a = " + dfCol(r, n, 0), "b = " + dfCol(r, n, 1), "c = " + dfCol(r, n, 2), "d = " + dfCol(r, n, 3)}
	r.Shuffle(len(cols), func(i, j int) { cols[i], cols[j] = cols[j], cols[i] })
	cols = cols[:2+r.Intn(3)]
	if withFactor && r.Intn(2) == 0 {
		cols = append(cols, "f = factor("+dfCol(r, n, 2)+")")
	}
	extra := ""
	if r.Intn(4) == 0 {
		extra = ", stringsAsFactors = TRUE"
	}
	if r.Intn(5) == 0 {
		names := make([]string, n)
		for i := range names {
			names[i] = fmt.Sprintf("\"r%d\"", i+1)
		}
		extra += ", row.names = c(" + strings.Join(names, ", ") + ")"
	}
	return "data.frame(" + strings.Join(cols, ", ") + extra + ")"
}

func fuzzDataFrame(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	n := 1 + r.Intn(5)
	df := dfSrc(r, n, true)
	var family, body string
	idx := func() string {
		return pick([]string{"1", "2:3", "-1", "c(TRUE, FALSE)", "c(1, 1)", "n + 1", "c(1, NA)", "integer(0)", "NULL", "\"r1\"", "which(df$a > 0)", "df$a > 0"})
	}
	col := func() string {
		return pick([]string{"1", "2", "\"a\"", "\"b\"", "c(1, 2)", "-1", "\"zz\"", "c(\"a\", \"b\")", "c(TRUE, FALSE)"})
	}
	switch r.Intn(16) {
	case 0:
		family, body = "df-build", df
	case 1:
		family, body = "df-print", "print("+df+pick([]string{"", ", digits = 3", ", row.names = FALSE", ", right = FALSE"})+")"
	case 2:
		family, body = "df-index", "df <- "+df+"; n <- nrow(df); df["+idx()+", "+col()+pick([]string{"", ", drop = FALSE", ", drop = TRUE"})+"]"
	case 3:
		family, body = "df-index", "df <- "+df+"; n <- nrow(df); df["+col()+"]"
	case 4:
		family, body = "df-index2", "df <- "+df+"; list(df[[1]], df$b, df[[\"a\"]], names(df), dim(df), dimnames(df))"
	case 5:
		family, body = "df-assign", "df <- "+df+"; df$"+pick([]string{"a", "z", "b"})+" <- "+pick([]string{"NULL", "1", "df$a", "1:2", "\"q\"", "n:1"})+"; df"
	case 6:
		family, body = "df-assign", "df <- "+df+"; n <- nrow(df); df["+idx()+", "+col()+"] <- "+pick([]string{"0L", "NA", "\"x\"", "list(1L, 2L)", "99", "TRUE"})+"; df"
	case 7:
		family, body = "df-assign", "df <- "+df+"; df[[\""+pick([]string{"a", "new"})+"\"]] <- "+pick([]string{"1", "NULL", "df$a", "letters[1:n]"})+"; df"
	case 8:
		family, body = "df-bind", "df <- "+df+"; "+pick([]string{"rbind(df, df)", "cbind(df, z = 1)", "rbind(df, df[1, ])", "cbind(df, df)", "rbind(df[0, ], df)", "do.call(rbind, list(df, df, df))"})
	case 9:
		other := dfSrc(r, 1+r.Intn(4), false)
		family, body = "df-merge", "x <- "+df+"; y <- "+other+"; merge(x, y"+pick([]string{"", ", all = TRUE", ", all.x = TRUE", ", by = \"a\"", ", by = character(0)", ", sort = FALSE"})+")"
	case 10:
		family, body = "df-subset", "df <- "+df+"; "+pick([]string{"subset(df, a > 0)", "subset(df, !is.na(b), select = a)", "transform(df, z = a * 2)", "transform(df, a = NULL)", "df[order(df$a), ]", "df[order(df$b, decreasing = TRUE), ]", "unique(df)", "duplicated(df)", "as.matrix(df)", "t(df)", "is.na(df)", "format(df)", "split(df, df$a > 0)", "lapply(df, class)", "sapply(df, is.numeric)", "with(df, a)"})
	case 11:
		vals := dfCol(r, 3+r.Intn(8), 2)
		family, body = "table1", "x <- "+vals+"; table(x"+pick([]string{"", ", useNA = \"ifany\"", ", useNA = \"always\"", ", exclude = \"a\"", ", dnn = \"v\"", ", exclude = NULL"})+")"
	case 12:
		m := 3 + r.Intn(6)
		family, body = "table2", "x <- "+dfCol(r, m, 2)+"; y <- "+dfCol(r, m, 3)+"; t2 <- table(x, y"+pick([]string{"", ", useNA = \"ifany\""})+"); "+pick([]string{"t2", "as.data.frame(t2)", "dim(t2)", "t2[1, ]", "prop.table(t2, 1)", "margin.table(t2, 2)", "summary(c(t2))"})
	case 13:
		family, body = "cut", "cut(c("+pick([]string{"1, 2, 3, 4, 5", "0.5, 7.25, 3", "-3, 0, 3, NA", "10, 20, 30, 40"})+"), "+pick([]string{"3", "c(0, 2, 5, 50)", "c(-5, 0, 5)", "2, dig.lab = 2", "c(0, 4, 10), right = FALSE", "c(0, 4, 40), include.lowest = TRUE, labels = c(\"lo\", \"hi\")", "c(0, 4, 40), labels = FALSE"})+")"
	case 14:
		family, body = "expand", pick([]string{"expand.grid(a = 1:2, b = c(\"x\", \"y\", \"z\"))", "expand.grid(1:2, 3:4, KEEP.OUT.ATTRS = FALSE)", "interaction(c(\"a\", \"b\", \"a\"), c(\"x\", \"x\", \"y\"))", "interaction(c(\"a\", \"b\"), c(\"x\", \"y\"), lex.order = TRUE, drop = TRUE)", "nrow(expand.grid(1:3, 1:4, 1:2))"})
	default:
		family, body = "factor-ops", "f <- factor("+dfCol(r, 2+r.Intn(5), 2)+"); "+pick([]string{"f[2] <- \"a\"; f", "c(f, factor(\"zz\"))", "droplevels(f[1])", "addNA(f)", "levels(f)[1] <- \"q\"; f", "f == \"a\"", "table(f, useNA = \"always\")", "factor(f, exclude = NULL)"})
	}
	return "# " + family + "\n" + body
}
