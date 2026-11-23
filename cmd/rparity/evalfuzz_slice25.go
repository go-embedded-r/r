package main

import (
	"fmt"
	"math/rand"
)

// The slice25 mode: random programs over the special functions of slice 25:
// gamma, lgamma, beta, lbeta, digamma, trigamma, psigamma, besselI/J/K/Y and
// the exp and ^ edges. The first line is "# family". A value ends in "# libm"
// when it is compared within --ulp units in the last place after being rounded
// to twelve digits; the edge families (zeros, poles, infinities, NA, warnings)
// are compared exactly.

func fuzzSlice25(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	// an argument vector: ordinary values with the occasional edge
	arg := func(pos bool) string {
		n := 1 + r.Intn(6)
		s := "c("
		for i := 0; i < n; i++ {
			if i > 0 {
				s += ", "
			}
			switch r.Intn(12) {
			case 0:
				s += pick([]string{"0", "NA", "Inf", "1e5"})
			case 1:
				if !pos {
					s += pick([]string{"-1", "-0.5", "-2.5", "-Inf"})
					continue
				}
				fallthrough
			default:
				v := float64(r.Intn(4000)) / 64
				if !pos && r.Intn(4) == 0 {
					v = -v
				}
				s += fmt.Sprintf("%g", v)
			}
		}
		return s + ")"
	}
	nu := func() string {
		return pick([]string{"0", "1", "2", "0.5", "1.5", "-0.5", "-1", "2.75", "c(0, 1)", "c(0.5, 1.5, 2)", "10", "25.5", "-2.25"})
	}
	var family, body string
	switch r.Intn(8) {
	case 0:
		family = "gamma"
		body = "signif(" + pick([]string{"gamma", "lgamma", "digamma", "trigamma"}) + "(" + arg(false) + "), 8) # libm"
	case 1:
		family = "psigamma"
		body = fmt.Sprintf("signif(psigamma(%s, %s), 8) # libm", arg(true), pick([]string{"0", "1", "2", "3", "4", "5", "c(0, 1, 2)", "2.4", "2.6"}))
	case 2:
		family = "beta"
		body = "signif(" + pick([]string{"beta", "lbeta"}) + "(" + arg(true) + ", " + arg(true) + "), 8) # libm"
	case 3:
		family = "besselJ"
		body = fmt.Sprintf("signif(besselJ(%s, %s), 8) # libm", arg(true), nu())
	case 4:
		family = "besselY"
		body = fmt.Sprintf("signif(besselY(%s, %s), 8) # libm", arg(true), nu())
	case 5:
		family = "besselIK"
		fn := pick([]string{"besselI", "besselK"})
		body = fmt.Sprintf("signif(%s(%s, %s, %s), 8) # libm", fn, arg(true), nu(), pick([]string{"TRUE", "FALSE"}))
	case 6:
		family = "edges"
		body = pick([]string{
			"gamma(c(0, -1, Inf, -Inf, NA, NaN))", "lgamma(c(0, -1, Inf, -Inf, NA))", "digamma(c(0, -1, Inf, -Inf, NA))",
			"trigamma(c(0, -1, Inf, -Inf, NA))", "beta(c(-1, 0, 1, Inf, NA), c(2, 2, 0, 2, 1))", "lbeta(c(-1, 0, 1, Inf, NA), c(2, 2, 0, 2, 1))",
			"suppressWarnings(besselJ(c(-1, 0, Inf, NA), 0))", "suppressWarnings(besselK(c(-1, 0, Inf, NA), 1))",
			"suppressWarnings(besselY(c(-1, 0, Inf, NA), 0))", "suppressWarnings(besselI(c(-1, 0, Inf, NA), 0))",
			"psigamma(c(0, -1, Inf, -Inf), 3)", "psigamma(1, c(-1, NA, NaN))",
			"tryCatch(besselJ(1), error = conditionMessage)", "tryCatch(besselJ('a', 1), error = conditionMessage)",
			"tryCatch(beta(1), error = conditionMessage)", "tryCatch(gamma('a'), error = conditionMessage)",
			"names(gamma(c(a = 1, b = 2)))", "dim(besselJ(matrix(1:4, 2), 1))", "length(besselK(1:3, c(0, 1)))",

			"withCallingHandlers(besselY(1, 2e7), warning = function(w) {print(conditionMessage(w)); invokeRestart('muffleWarning')})",
			"withCallingHandlers(besselJ(-1, 1), warning = function(w) {print(conditionMessage(w)); invokeRestart('muffleWarning')})",
			"withCallingHandlers(psigamma(1, 101), warning = function(w) {print(conditionMessage(w)); invokeRestart('muffleWarning')})",
		})
	default:
		family = "exppow"
		body = pick([]string{
			fmt.Sprintf("signif(exp(c(%d, %g, %g, -%g)), 14)", r.Intn(700)-350, float64(r.Intn(5000))/64, float64(r.Intn(900))/64, float64(r.Intn(11000))/16) + " # libm",
			fmt.Sprintf("signif(%g^c(%g, %g, %d), 14)", 0.25+float64(r.Intn(400))/16, float64(r.Intn(900))/64, -float64(r.Intn(900))/64, r.Intn(40)-20) + " # libm",
			"exp(c(709.79, -745.2, -1e5, Inf, -Inf, NA))", "(-8)^c(1/3, 3, 0.5, 1e15, -3)", "2^c(1024, -1074, -1075)", "0^c(-1, 0, 1, NA)", "c(-Inf, Inf)^c(-1, 0.5, 2, 3)",
		})
	}
	return "# " + family + "\n" + body
}
