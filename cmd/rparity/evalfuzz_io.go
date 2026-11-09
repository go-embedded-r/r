package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// The io mode: random programs over text connections, files, sinks,
// capture.output, the file-system functions, readRDS/saveRDS, source and the
// environment variables. Every program works in its own temporary directory
// and closes what it opens; paths never reach the output (a warning that names
// one goes through sub()). The first line is "# family". The mode has no
// counted exclusion: everything is compared exactly.

var ioLines = []string{
	"", "a", "hello world", "  lead", "trail  ", "tab\\there", "x,y", "\u00e9\u4e2d", "NA", "a\\rb", "q\\\"uote", "1 2 3", "# c", "\\\\back", "z",
}

func fuzzIO(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	q := func(s string) string { return "\"" + s + "\"" }
	lines := func(n int) string {
		parts := make([]string, n)
		for i := range parts {
			parts[i] = q(pick(ioLines))
		}
		return "c(" + strings.Join(parts, ", ") + ")"
	}
	norm := `sub(f, "FILE", conditionMessage(w), fixed = TRUE)`
	var family, body string
	switch r.Intn(12) {
	case 0:
		family = "text-read"
		var steps []string
		for i := 0; i < 1+r.Intn(4); i++ {
			steps = append(steps, fmt.Sprintf("list(%s, readLines(tc, %s))", fmt.Sprint(i), pick([]string{"1", "2", "0", "-1", "n = 3", "5"})))
		}
		src := lines(1 + r.Intn(5))
		if r.Intn(3) == 0 {
			src = q(strings.Join([]string{pick(ioLines), pick(ioLines)}, pick([]string{"\\n", "\\r\\n", "\\r"})))
		}
		body = "tc <- textConnection(" + src + "); r <- list(" + strings.Join(steps, ", ") + ", isOpen(tc)); close(tc); r"
	case 1:
		family = "text-write"
		var ops []string
		for i := 0; i < 1+r.Intn(5); i++ {
			switch r.Intn(4) {
			case 0:
				ops = append(ops, "writeLines("+lines(1+r.Intn(3))+", tc"+pick([]string{"", ", sep = \"|\"", ", sep = \"\""})+")")
			case 1:
				ops = append(ops, "cat("+q(pick(ioLines))+", "+q(pick(ioLines))+", file = tc"+pick([]string{"", ", sep = \"\\n\"", ", sep = \"-\""})+")")
			case 2:
				ops = append(ops, "cat("+pick([]string{"1:3", "TRUE", "2.5", "\"s\\n\""})+", file = tc)")
			default:
				ops = append(ops, "writeLines("+q(pick(ioLines))+", tc)")
			}
		}
		mode := pick([]string{"\"w\"", "\"a\""})
		body = "tc <- textConnection(\"zz\", " + mode + ", local = TRUE); " + strings.Join(ops, "; ") +
			"; mid <- list(zz, isIncomplete(tc), textConnectionValue(tc)); close(tc); list(mid, zz)"
	case 3:
		family = "file-roundtrip"
		app := pick([]string{"", ", append = TRUE"})
		sep := pick([]string{"", ", sep = \"\\n\"", ", sep = \", \""})
		body = "f <- tempfile(); cat(" + q(pick(ioLines)) + ", " + q(pick(ioLines)) + ", file = f" + sep + "); cat(" + q(pick(ioLines)) + ", file = f" + app + sep + ");" +
			" r <- tryCatch(readLines(f), warning = function(w) " + norm + "); sz <- file.size(f); unlink(f); list(r, sz)"
	case 4:
		family = "file-lines"
		n := pick([]string{"-1L", "1L", "2L", "0L", "10L"})
		body = "f <- tempfile(); writeLines(" + lines(1+r.Intn(5)) + ", f" + pick([]string{"", ", sep = \"\\r\\n\"", ", sep = \"|\""}) + "); r <- list(suppressWarnings(readLines(f, n = " + n + ")), file.size(f), file.exists(f)); unlink(f); r"
	case 5:
		family = "file-conn"
		body = "f <- tempfile(); con <- file(f, " + pick([]string{"\"w\"", "\"wt\"", "\"a\""}) + "); writeLines(" + lines(2) + ", con); cat(" + q(pick(ioLines)) + ", file = con); r1 <- isOpen(con); close(con);" +
			" con <- file(f, \"r\"); a <- suppressWarnings(readLines(con, 1)); b <- suppressWarnings(readLines(con, 1)); d <- suppressWarnings(readLines(con)); close(con); unlink(f); list(r1, a, b, d)"
	case 6:
		family = "file-ops"
		names := []string{"a.txt", "b.csv", "C.txt", ".h", "d_1.R"}
		var ops []string
		for i := 0; i < 2+r.Intn(5); i++ {
			n := pick(names)
			switch r.Intn(8) {
			case 0:
				ops = append(ops, fmt.Sprintf("file.create(file.path(d, %s))", q(n)))
			case 1:
				ops = append(ops, fmt.Sprintf("writeLines(%s, file.path(d, %s))", lines(2), q(n)))
			case 2:
				ops = append(ops, fmt.Sprintf("suppressWarnings(file.copy(file.path(d, %s), file.path(d, %s)))", q(n), q(pick(names))))
			case 3:
				ops = append(ops, fmt.Sprintf("suppressWarnings(file.rename(file.path(d, %s), file.path(d, %s)))", q(n), q(pick(names))))
			case 4:
				ops = append(ops, fmt.Sprintf("suppressWarnings(file.remove(file.path(d, %s)))", q(n)))
			case 5:
				ops = append(ops, fmt.Sprintf("suppressWarnings(dir.create(file.path(d, %s)))", q("sub"+fmt.Sprint(r.Intn(2)))))
			case 6:
				ops = append(ops, fmt.Sprintf("file.exists(file.path(d, %s))", q(n)))
			default:
				ops = append(ops, fmt.Sprintf("suppressWarnings(file.size(file.path(d, %s)))", q(n)))
			}
		}
		lf := pick([]string{"list.files(d)", "list.files(d, all.files = TRUE)", "list.files(d, pattern = \"\\\\.txt$\")", "list.files(d, recursive = TRUE)", "list.files(d, all.files = TRUE, recursive = TRUE, include.dirs = TRUE)", "list.dirs(d, full.names = FALSE)", "list.files(d, pattern = \"^[a-c]\", ignore.case = TRUE)"})
		body = "d <- tempfile(); dir.create(d); res <- list(" + strings.Join(ops, ", ") + "); out <- " + lf + "; unlink(d, recursive = TRUE); list(res, out, dir.exists(d))"
	case 7:
		family = "sink-capture"
		items := []string{"print(1:3)", "cat(\"a\\nb\")", "print(\"s\")", "x <- 1", "invisible(2)", "message(\"m\")", "writeLines(\"w\")", "print(NULL)", "(function() 5)()", "c(a = 1)"}
		switch r.Intn(3) {
		case 0:
			body = "capture.output(" + pick(items) + ", " + pick(items) + ")"
		case 1:
			body = "tc <- textConnection(\"zz\", \"w\"); sink(tc); " + pick(items) + "; " + pick(items) + "; sink(); close(tc); zz"
		default:
			body = "t1 <- textConnection(\"z1\", \"w\"); t2 <- textConnection(\"z2\", \"w\"); sink(t1); print(1); sink(t2); " + pick(items) + "; sink(); " + pick(items) + "; sink(); close(t1); close(t2); list(z1, z2)"
		}
	case 8:
		family = "rds"
		vals := []string{"1:5", "c(a = 1.5, b = NA, c = -Inf)", "c(\"x\", NA, \"\u00e9\")", "list(1, \"a\", NULL, list(TRUE, 2L))", "factor(c(\"u\", \"v\", \"u\"))", "matrix(1:6, 2, dimnames = list(NULL, c(\"p\", \"q\", \"r\")))", "NULL", "as.Date(\"2024-02-29\")", "TRUE", "quote(f(x, y = 2))", "structure(list(a = 1), class = \"k\")", "c(1.5, NaN, NA, 1e300)", "integer(0)", "character(0)"}
		v := pick(vals)
		opt := pick([]string{"", ", compress = FALSE", ", version = 2", ", compress = TRUE"})
		body = "x <- " + v + "; f <- tempfile(); saveRDS(x, f" + opt + "); y <- readRDS(f); unlink(f); list(identical(x, y), y)"
	case 9:
		family = "env"
		k := pick([]string{"RGO_A", "RGO_B", "RGO_C"})
		body = "Sys.setenv(" + k + " = " + q(pick(ioLines)) + ", RGO_Z = \"9\"); r1 <- Sys.getenv(" + pick([]string{"\"" + k + "\"", "c(\"" + k + "\", \"RGO_NONE\")", "c(\"RGO_Z\", \"" + k + "\")"}) + pick([]string{"", ", names = FALSE", ", names = TRUE", ", unset = NA", ", unset = \"d\""}) + "); r2 <- Sys.unsetenv(\"" + k + "\"); list(r1, r2, Sys.getenv(\"" + k + "\", unset = NA))"
	case 10:
		family = "source"
		code := pick([]string{"x1 <- 1; x1 + 1", "f1 <- function(a) a * 2\\nf1(4)", "invisible(7)", "y1 <- c(1, 2)\\ny1[2]", "if (TRUE) 3 else 4", "NULL", "z1 <- 'text'; nchar(z1)", ""})
		how := pick([]string{"", ", local = TRUE", ", print.eval = TRUE", ", local = e"})
		body = "e <- new.env(); tc <- textConnection(" + q(code) + "); r <- source(tc" + how + "); close(tc); list(r, ls(e))"
	default:
		family = "cat-sep"
		var args []string
		for i := 0; i < 1+r.Intn(4); i++ {
			args = append(args, pick([]string{"\"a\"", "1L", "2.5", "TRUE", "\"b c\"", "NULL", "\"\"", "character(0)"}))
		}
		seps := []string{"", ", sep = \"\"", ", sep = \"\\n\"", ", sep = \", \"", ", sep = c(\",\", \"\\n\")", ", sep = \" \\n\"", ", sep = \"\\n\\n\""}
		body = "cat(" + strings.Join(args, ", ") + pick(seps) + ")"
	}
	for i := 0; i < len(body); i++ {
		if body[i] >= 0x80 {
			return "# excl:" + exclEncoding + "\n" + body
		}
	}
	return "# " + family + "\n" + body
}

// exclEncoding is the one counted exclusion of the io mode.
const exclEncoding = "encoding"
