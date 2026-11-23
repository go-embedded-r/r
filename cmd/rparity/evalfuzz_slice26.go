package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// The slice26 mode: random calls of d, p and q of the twelve continuous
// families of slice 26, over ordinary values, the parameter edges (0, negative,
// Inf, NA, NaN, huge) and the flag combinations log, lower.tail and log.p. The
// first line is "# family". Ordinary values end in "# libm": they are compared
// within --ulp units in the last place after being rounded to nine digits; the
// edge calls are compared after rounding to six digits.

type fuzzFam struct {
	name string
	// ordinary and edge argument generators for each parameter
	params [][2][]string
	// the support-ish values of x
	xs []string
	// call template pieces: the named parameters
	names []string
}

var fuzzFams26 = []fuzzFam{
	{"norm", [][2][]string{
		{{"0", "1", "-2.5", "10", "0.3"}, {"-1", "Inf", "-Inf", "NA"}},
		{{"1", "0.5", "2", "7.25", "0.01"}, {"0", "-1", "Inf", "NaN"}}},
		[]string{"-8", "-3", "-0.5", "0", "0.7", "2", "5", "30", "-40"}, []string{"mean", "sd"}},
	{"unif", [][2][]string{
		{{"0", "-1", "2.5", "-10"}, {"Inf", "NA"}},
		{{"1", "3", "2.75", "11"}, {"-1", "0", "Inf", "NaN"}}},
		[]string{"-3", "0", "0.25", "0.5", "1", "2.6", "4", "20"}, []string{"min", "max"}},
	{"exp", [][2][]string{
		{{"1", "0.5", "3", "0.01", "20"}, {"0", "-1", "Inf", "NA"}}},
		[]string{"-1", "0", "0.1", "1", "5", "30", "120", "1e3"}, []string{"rate"}},
	{"gamma", [][2][]string{
		{{"1", "0.5", "2", "7.5", "50", "0.01", "1000"}, {"0", "-1", "Inf", "NA"}},
		{{"1", "0.5", "3", "2.5"}, {"0", "-1", "Inf", "NaN"}}},
		[]string{"-1", "0", "0.1", "1", "3", "10", "60", "1000", "1e5"}, []string{"shape", "scale"}},
	{"beta", [][2][]string{
		{{"1", "0.5", "2", "7.5", "50", "0.01", "300"}, {"0", "-1", "Inf", "NA"}},
		{{"1", "0.5", "3", "2.5", "100", "0.05"}, {"0", "-1", "Inf", "NaN"}}},
		[]string{"-0.5", "0", "0.001", "0.1", "0.5", "0.9", "0.999", "1", "1.5"}, []string{"shape1", "shape2"}},
	{"t", [][2][]string{
		{{"1", "2", "5", "30", "0.5", "1000", "1e8"}, {"0", "-1", "Inf", "NA"}}},
		[]string{"-1e3", "-30", "-3", "-0.5", "0", "0.7", "2", "10", "200"}, []string{"df"}},
	{"f", [][2][]string{
		{{"1", "2", "5", "30", "0.5", "200"}, {"0", "-1", "Inf", "NA"}},
		{{"1", "3", "10", "2.5", "100"}, {"0", "-1", "Inf", "NaN"}}},
		[]string{"-1", "0", "0.01", "0.5", "1", "3", "20", "1e3"}, []string{"df1", "df2"}},
	{"chisq", [][2][]string{
		{{"1", "2", "5", "30", "0.5", "200", "1e4"}, {"0", "-1", "Inf", "NA"}}},
		[]string{"-1", "0", "0.01", "0.5", "1", "3", "20", "300", "1e4"}, []string{"df"}},
	{"lnorm", [][2][]string{
		{{"0", "1", "-2", "5"}, {"-1", "Inf", "NA"}},
		{{"1", "0.5", "2", "0.1"}, {"0", "-1", "Inf", "NaN"}}},
		[]string{"-1", "0", "0.01", "0.5", "1", "3", "20", "1e5"}, []string{"meanlog", "sdlog"}},
	{"logis", [][2][]string{
		{{"0", "1", "-2", "10"}, {"Inf", "NA", "-Inf"}},
		{{"1", "0.5", "2", "7"}, {"0", "-1", "Inf", "NaN"}}},
		[]string{"-800", "-30", "-3", "0", "0.7", "2", "30", "800"}, []string{"location", "scale"}},
	{"cauchy", [][2][]string{
		{{"0", "1", "-2", "10"}, {"Inf", "NA", "-Inf"}},
		{{"1", "0.5", "2", "7"}, {"0", "-1", "Inf", "NaN"}}},
		[]string{"-1e10", "-30", "-3", "0", "0.7", "2", "30", "1e10"}, []string{"location", "scale"}},
	{"weibull", [][2][]string{
		{{"1", "0.5", "2", "7.5", "50"}, {"0", "-1", "Inf", "NA"}},
		{{"1", "0.5", "3", "2.5"}, {"0", "-1", "Inf", "NaN"}}},
		[]string{"-1", "0", "0.01", "0.5", "1", "3", "20", "100"}, []string{"shape", "scale"}},
}

func fuzzSlice26(r *rand.Rand) string {
	f := fuzzFams26[r.Intn(len(fuzzFams26))]
	kind := []string{"d", "p", "q"}[r.Intn(3)]
	logq := r.Intn(3) == 0
	lower := r.Intn(2) == 0
	dlog := kind == "d" && r.Intn(2) == 0
	// values left out of a call, with the reason (README, limits):
	//  - a log density or log tail of gamma, beta, t, F, chi-squared, exp,
	//    weibull or unif far in a tail: go-statistics underflows to -Inf there
	//    and R answers a number (upstream gap);
	//  - a deep log-probability quantile of those families (same gap), and of
	//    F, where R's own answer is off by up to 1 % (pf of R's answer is not p);
	//  - an infinite or zero degrees of freedom of chisq (R answers NaN with a
	//    warning for pchisq), an infinite location of lnorm with sdlog 0, an
	//    infinite location or scale of logis and cauchy (R's answers are
	//    inconsistent from one function to the next).
	loc := f.name == "norm" || f.name == "lnorm" || f.name == "logis" || f.name == "cauchy"
	banned := map[string]bool{}
	if !loc && (dlog || logq) {
		for _, s := range []string{"1e5", "1e3", "1000", "1e30", "-1e30", "-1e3", "1e4", "200", "300", "50", "120", "60"} {
			banned[s] = true
		}
	}
	switch f.name {
	case "chisq", "lnorm", "logis", "cauchy":
		banned["Inf"] = true
		banned["-Inf"] = true
	}
	// R's own quantile is off by 1e-6 to 1e-2 in relative terms where pt or pf
	// of its answer is not the probability (the Go answer is the closer one),
	// and R answers its smallest double where the true quantile underflows
	switch {
	case f.name == "f" && kind == "q":
		banned["0.5"], banned["1e-10"] = true, true
		if logq {
			banned["1"] = true
		}
	case f.name == "t" && kind == "q":
		banned["0.5"] = true
	case (f.name == "beta" || f.name == "gamma") && kind == "q":
		banned["0.01"] = true // a quantile of 1e-300: the ninth digit is the last bit
	case f.name == "weibull" && kind == "d":
		banned["Inf"] = true
	}
	// R answers NaN for the density of F at x = Inf with some degrees of
	// freedom; the gamma tail at x = Inf with an infinite shape and scale
	banInfX := f.name == "f" && kind == "d"
	keep := func(vs []string, isX bool) []string {
		var out []string
		for _, s := range vs {
			if banned[s] && !(isX && (s == "Inf" || s == "-Inf")) || isX && banInfX && s == "Inf" {
				continue
			}
			out = append(out, s)
		}
		return out
	}
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	edge := r.Intn(4) == 0
	// one argument vector: one to four values
	vec := func(vals []string, ev []string, isX bool) string {
		vals, ev = keep(vals, isX), keep(ev, isX)
		n := 1 + r.Intn(4)
		if n == 1 {
			if edge && r.Intn(2) == 0 {
				return pick(ev)
			}
			return pick(vals)
		}
		parts := make([]string, n)
		for i := range parts {
			if edge && r.Intn(3) == 0 {
				parts[i] = pick(ev)
			} else {
				parts[i] = pick(vals)
			}
		}
		return "c(" + strings.Join(parts, ", ") + ")"
	}
	var first string
	switch kind {
	case "q":
		first = vec([]string{"0.001", "0.1", "0.25", "0.5", "0.9", "0.999", "1e-10"}, []string{"0", "1", "-0.5", "1.5", "NA", "NaN"}, true)
	default:
		first = vec(f.xs, []string{"Inf", "-Inf", "NA", "NaN", "1e30", "-1e30"}, true)
	}
	if kind == "q" && logq {
		// below about -745 the log-cdf of go-statistics for gamma, beta, t, f and
		// chi-squared underflows (upstream gap, README): R is compared only above
		deep := []string{"-1e-5", "-0.1", "-0.7", "-2"}
		if loc {
			deep = append(deep, "-30", "-700", "-1e5", "-1e10")
		}
		first = vec(deep, []string{"0", "-Inf", "0.5", "NA", "NaN"}, true)
	}
	params := append([][2][]string(nil), f.params...)
	if f.name == "gamma" && kind == "p" {
		params[1] = [2][]string{params[1][0], {"0", "-1", "NaN"}}
	}
	args := first
	// a parameter with a default may be left out; once one is, the others are
	// named so that the call stays well defined
	hasDef := map[string]bool{"mean": true, "sd": true, "min": true, "max": true, "rate": true, "scale": true, "meanlog": true, "sdlog": true, "location": true}
	named := false
	for i, nm := range f.names {
		v := vec(params[i][0], params[i][1], false)
		if hasDef[nm] && r.Intn(5) == 0 {
			named = true
			continue
		}
		switch {
		case f.name == "gamma" && nm == "scale" && r.Intn(2) == 0:
			args += ", rate = " + vec([]string{"1", "0.5", "3", "2"}, []string{"0", "-1", "Inf", "NaN"}, false)
		case named || r.Intn(2) == 0:
			args += ", " + nm + " = " + v
		default:
			args += ", " + v
		}
	}
	switch kind {
	case "d":
		if dlog {
			args += ", log = TRUE"
		}
	default:
		if !lower {
			args += ", lower.tail = FALSE"
		}
		if logq && kind == "p" {
			args += ", log.p = TRUE"
		} else if kind == "q" && logq {
			args += ", log.p = TRUE"
		}
	}
	call := fmt.Sprintf("%s%s(%s)", kind, f.name, args)
	if edge {
		return "# " + f.name + "\n" + "signif(" + call + ", 6)"
	}
	return "# " + f.name + "\n" + "signif(" + call + ", 9) # libm"
}

// nameAll is a no-op hook kept for clarity: every parameter was either named
// or positional in order, and a skipped parameter is only ever one that has a
// default, so later parameters are always named below.
func nameAll(args string, f fuzzFam, kind string) string { return args }
