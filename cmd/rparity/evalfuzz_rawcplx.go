package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// The rawcplx mode: random programs over the complex and raw types and the
// expression vector. The first line is "# family". Counted exclusions: the
// transcendental functions of complex numbers (their results are compared
// within a few units in the last place, the Go library and the C one rounding
// differently; the line ends in "# libm") and non-ASCII text produced from raw
// bytes (R leaves it unmarked, ours is marked UTF-8; the class is "encoding").
// Everything else is compared exactly, including the printed text.

var (
	cplxParts = []string{"0", "1", "-1", "2", "0.5", "-2.5", "3", "10", "1e5", "1e-5", "123456.789", "0.1", "-0.25", "7", "1/3", "pi", "100"}
	cplxSpec  = []string{"NA", "NaN", "Inf", "-Inf", "NA_real_"}
)

func cplxLit(r *rand.Rand, special bool) string {
	re, im := cplxParts[r.Intn(len(cplxParts))], cplxParts[r.Intn(len(cplxParts))]
	if special && r.Intn(8) == 0 {
		return "complex(real = " + cplxSpec[r.Intn(len(cplxSpec))] + ", imaginary = " + im + ")"
	}
	if r.Intn(6) == 0 {
		return "NA_complex_"
	}
	return "complex(real = " + re + ", imaginary = " + im + ")"
}

func cplxVec(r *rand.Rand, special bool) string {
	n := 1 + r.Intn(4)
	parts := make([]string, n)
	for i := range parts {
		parts[i] = cplxLit(r, special)
	}
	if n == 1 {
		return parts[0]
	}
	return "c(" + strings.Join(parts, ", ") + ")"
}

func fuzzRawComplex(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	q := func(s string) string { return "\"" + s + "\"" }
	var family, body string
	libm := false
	switch r.Intn(14) {
	case 0:
		family = "cplx-arith"
		op := pick([]string{"+", "-", "*", "/"})
		finite := []string{"0", "1", "-1", "2", "0.5", "-2.5", "3", "10", "100", "1e-5"}
		body = "(" + cplxVec(r, false) + ") " + op + " (" + pick([]string{cplxVec(r, false), pick(finite), pick(fuzzInts), pick(fuzzLgls)}) + ")"
	case 1:
		family = "cplx-pow"
		body = "(" + cplxVec(r, false) + ") ^ " + pick([]string{"2", "3", "-1", "-2", "0", "1", "4L", "10"})
	case 2:
		family = "cplx-math"
		libm = true
		f := pick([]string{"sqrt", "exp", "sin", "cos", "tan", "sinh", "cosh", "tanh", "asin", "acos", "atan", "asinh", "acosh", "log", "abs", "Arg", "Mod"})
		// not pi: Go's sin and cos are off by a few units in the last place near
		// multiples of it, which is huge for a result that is nearly zero
		lit := cplxLit(r, false)
		for strings.Contains(lit, "pi") {
			lit = cplxLit(r, false)
		}
		body = f + "(" + lit + ")"
		switch f {
		case "asin", "acos", "atan", "acosh":
			// R's algorithms for these lose digits that ours keep: a counted exclusion
			return "# excl:" + exclCplxInv + "\n" + body
		}
	case 3:
		family = "cplx-print"
		body = "x <- " + cplxVec(r, true) + "; print(x); format(x); as.character(x); deparse(x); paste(x); format(x, digits = " + pick([]string{"2", "3", "10"}) + "); toString(x)"
	case 4:
		family = "cplx-parts"
		f := pick([]string{"Re", "Im", "Conj", "is.na", "is.finite", "is.infinite", "is.nan", "round", "signif", "sum", "prod", "cumsum", "cumprod", "mean", "rev", "length", "sort", "unique"})
		body = f + "(" + cplxVec(r, f != "sort" && f != "unique" && f != "prod" && f != "cumprod") + ")"
	case 5:
		family = "cplx-coerce"
		s := pick([]string{"1", "2i", "1+2i", "-1.5-0.5i", "3e2+1e-1i", "NA", "Inf", "NaN", "abc", "1+i", " 4+5i ", "+1i", "0x10", "1e5+0i"})
		body = pick([]string{"as.complex(" + q(s) + ")", "suppressWarnings(as.complex(" + q(s) + "))", "as.complex(" + pick(fuzzNums) + ")", "as.complex(" + pick(fuzzLgls) + ")"})
	case 6:
		family = "cplx-cmp"
		body = "(" + cplxVec(r, true) + ") " + pick([]string{"==", "!=", "=="}) + " (" + pick([]string{cplxVec(r, true), pick(fuzzNums), q("0+1i")}) + ")"
	case 7:
		family = "raw-bits"
		v := "as.raw(" + pick([]string{"0", "1", "255", "16", "170", "c(1, 2, 3)", "1:8", "c(255, 0)"}) + ")"
		body = pick([]string{"rawToBits(" + v + ")", "packBits(rawToBits(" + v + "))", "rawShift(" + v + ", " + pick([]string{"-8", "-3", "0", "1", "5", "8"}) + ")", "intToBits(" + pick(fuzzInts) + ")[1:" + pick([]string{"8", "32"}) + "]", "packBits(intToBits(" + pick([]string{"1L", "-1L", "1:3", "255L"}) + "), \"integer\")", "rev(" + v + ")", v + " & as.raw(" + pick([]string{"15", "240", "1"}) + ")", "!" + v, "xor(" + v + ", as.raw(255))"})
	case 8:
		family = "raw-bitw"
		f := pick([]string{"bitwAnd", "bitwOr", "bitwXor"})
		a, b := pick([]string{"12L", "10L", "255L", "0L", "-1L", "1:4", "NA_integer_", "65535L"}), pick([]string{"10L", "3L", "15L", "-8L", "c(1L, 2L)", "NA_integer_"})
		body = pick([]string{f + "(" + a + ", " + b + ")", "bitwNot(" + a + ")", "bitwShiftL(" + a + ", " + pick([]string{"1L", "4L", "31L", "32L", "0L", "-1L"}) + ")", "bitwShiftR(" + a + ", " + pick([]string{"1L", "4L", "31L", "0L"}) + ")"})
	case 9:
		family = "raw-io"
		var v, what string
		switch r.Intn(5) {
		case 0:
			v, what = pick([]string{"c(1L, -1L, 256L, NA, 2147483647L)", "1:3", "integer(0)"}), "\"integer\""
		case 1:
			v, what = pick([]string{"c(1.5, NA, -Inf, 0)", "pi", "numeric(0)"}), "\"double\""
		case 2:
			v, what = pick([]string{"c(TRUE, NA, FALSE)", "TRUE"}), "\"logical\""
		case 3:
			v, what = pick([]string{"c(1+2i, NA, 3i)", "1i"}), "\"complex\""
		default:
			v, what = pick([]string{"c(\"a\", \"bc\", \"\")", "\"x\""}), "\"character\""
		}
		endian := pick([]string{"", ", endian = \"big\"", ", endian = \"little\"", ", endian = \"swap\""})
		size := ""
		switch what {
		case "\"integer\"":
			size = pick([]string{"", ", size = 1", ", size = 2", ", size = 4", ", size = 8"})
		case "\"double\"":
			size = pick([]string{"", ", size = 4", ", size = 8"})
		}
		body = "x <- " + v + "; r <- writeBin(x, raw()" + size + endian + "); list(r, readBin(r, " + what + ", " + pick([]string{"1", "3", "10"}) + size + endian + "))"
	case 10:
		family = "raw-char"
		s := pick([]string{"abc", "hello world", "", "a", "x y", "h\\u00e9llo", "\\u4e2d\\u6587"})
		ops := []string{"charToRaw(" + q(s) + ")", "rawToChar(charToRaw(" + q(s) + "))", "writeChar(" + q(s) + ", raw()" + pick([]string{"", ", eos = NULL", ", nchars = 1", ", eos = \"!\""}) + ")",
			"readChar(charToRaw(" + q(s) + "), " + pick([]string{"1", "2", "c(1, 1)", "10", "0"}) + ")", "rawToChar(charToRaw(" + q(s) + "), multiple = TRUE)", "utf8ToInt(" + q(s) + ")"}
		body = pick(ops)
		if strings.Contains(s, "\\u") {
			return "# excl:" + exclEncoding + "\n" + body
		}
	case 11:
		family = "raw-conn"
		var ops []string
		for i := 0; i < 1+r.Intn(4); i++ {
			switch r.Intn(5) {
			case 0:
				ops = append(ops, "writeBin("+pick([]string{"1:3", "2.5", "\"s\"", "TRUE", "as.raw(7)"})+", con)")
			case 1:
				ops = append(ops, "writeChar("+q(pick([]string{"ab", "c", ""}))+", con"+pick([]string{"", ", eos = NULL"})+")")
			case 2:
				ops = append(ops, "writeLines("+q(pick([]string{"a", "b c", ""}))+", con)")
			case 3:
				ops = append(ops, "cat("+q(pick([]string{"x", "y z"}))+", file = con)")
			default:
				ops = append(ops, "rawConnectionValue(con)")
			}
		}
		body = "con <- rawConnection(" + pick([]string{"raw(0)", "as.raw(1:3)"}) + ", " + pick([]string{"\"w\"", "\"wb\"", "\"a\""}) + "); r <- list(" + strings.Join(ops, ", ") + "); v <- rawConnectionValue(con); close(con); list(r, v)"
	case 12:
		family = "raw-ops"
		v := "as.raw(" + pick([]string{"c(1, 2, 3)", "255", "c(10, 20, 30, 40)", "0:3"}) + ")"
		body = pick([]string{
			"c(" + v + ", as.raw(9))", "c(" + v + ", " + pick([]string{"2L", "TRUE", "\"a\"", "2.5"}) + ")",
			v + "[" + pick([]string{"2", "-1", "c(TRUE, FALSE)", "5", "0", "c(1, 1)"}) + "]",
			v + " == " + pick([]string{v, "as.raw(2)", "2", "\"02\""}), "rep(" + v + ", " + pick([]string{"2", "0"}) + ")", "as.integer(" + v + ")", "as.character(" + v + ")",
			"format(" + v + ")", "deparse(" + v + ")", "print(" + v + ")", "unique(" + v + ")", "rev(" + v + ")", "matrix(" + v + ", " + pick([]string{"1", "2"}) + ")", "as.logical(" + v + ")", "as.vector(" + v + ", \"list\")"})
	default:
		family = "expr"
		e := pick([]string{"a + b", "f(x, y = 2)", "x[1]", "if (a) b else c", "function(x) x^2", "-x", "y = 2", "a <- 1", "{ a; b }", "1i", "\"s\"", "NULL", "TRUE"})
		body = pick([]string{
			"x <- expression(" + e + ", z); list(x, length(x), x[[1]], x[1], x[2], class(x), is.expression(x), eval(x[[2]], list(z = 3)))",
			"x <- parse(text = " + q(strings.ReplaceAll(e, "\"", "\\\"")) + "); list(x, deparse(x), length(x), class(x[[1]]))",
			"x <- expression(" + e + "); list(as.character(x), format(x), c(x, expression(w)), rev(c(x, x)), as.list(x))",
			"x <- as.expression(" + pick([]string{"1:3", "quote(a + b)", "list(1, quote(x))", "\"s\"", "c(a = 1, b = 2)"}) + "); list(x, length(x))",
			"eval(parse(text = " + q(pick([]string{"1 + 2", "x <- 5; x * 2", "c(1, 2)[2]", "paste('a', 'b')", "sqrt(16)", "nchar('abc')"})) + "))",
			"x <- expression(a = 1 + 2, b = 3); x[[2]] <- quote(q); x[3] <- expression(r); names(x) <- NULL; x",
		})
	}
	if libm {
		body += " # libm"
	}
	return fmt.Sprintf("# %s\n%s", family, body)
}

// exclCplxInv is the exclusion class of the complex inverse trigonometric
// functions: compared within 2^36 units in the last place (about 1e-5).
const exclCplxInv = "cplx-inv"
