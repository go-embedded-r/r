package eval

import (
	"math"

	"github.com/go-statistics/stats"
)

// pchisq and qchisq (slice 21), bound to the chi-squared distributions of
// go-statistics/stats (central and non-central). Only the R glue is here:
// recycling over the arguments, NA, the attributes of the longest argument, the
// domain errors (NaN with the warning "NaNs produced"), lower.tail and log.p.
//
// Limits, listed in the README: with log.p = TRUE qchisq goes through exp(p),
// so a probability below about e^-745 is 0; and a non-central distribution with
// df = 0 is not in go-statistics, so it answers NaN.

type chisqKind int

const (
	chisqP chisqKind = iota
	chisqQ
)

func chisqFn(name string, kind chisqKind) {
	regBuiltin(name, func(s *Session, c *callInfo) Value {
		first := "q"
		if kind == chisqQ {
			first = "p"
		}
		a := s.matchFormals(c, first, "df", "ncp", "lower.tail", "log.p")
		for i, nm := range []string{first, "df"} {
			if a[i] == nil {
				s.errorf(c.call, "argument \"%s\" is missing, with no default", nm)
			}
		}
		lower := s.boolArg(c, a[3], true, "lower.tail")
		logp := s.boolArg(c, a[4], false, "log.p")
		args := []Value{a[0], a[1]}
		hasNcp := a[2] != nil
		if hasNcp {
			args = append(args, a[2])
		}
		cols := make([][]float64, len(args))
		n := 0
		for i, v := range args {
			switch v.(type) {
			case *Double, *Integer, *Logical:
			default:
				s.errorf(c.call, "Non-numeric argument to mathematical function")
			}
			cols[i] = asDoubles(v)
			if len(cols[i]) == 0 {
				return &Double{V: []float64{}}
			}
			n = max(n, len(cols[i]))
		}
		out := make([]float64, n)
		nan := false
		for i := range out {
			x, df := cols[0][i%len(cols[0])], cols[1][i%len(cols[1])]
			ncp := 0.0
			if hasNcp {
				ncp = cols[2][i%len(cols[2])]
			}
			if isNAReal(x) || isNAReal(df) || (hasNcp && isNAReal(ncp)) {
				out[i] = naReal
				continue
			}
			if math.IsNaN(x) || math.IsNaN(df) || math.IsNaN(ncp) {
				out[i] = math.NaN()
				continue
			}
			out[i] = chisqValue(kind, x, df, ncp, hasNcp, lower, logp)
			if math.IsNaN(out[i]) {
				nan = true
			}
		}
		if nan {
			s.warningf(c.call, "NaNs produced")
		}
		var attrs []Attr
		for i, v := range args {
			if len(cols[i]) == n {
				attrs = attrsOf(v)
				break
			}
		}
		return &Double{V: out, Attr: attrs}
	})
}

// chisqValue is one value of pchisq (kind chisqP, x a quantile) or qchisq
// (chisqQ, x a probability), none of the arguments NaN.
func chisqValue(kind chisqKind, x, df, ncp float64, hasNcp, lower, logp bool) float64 {
	if df < 0 || ncp < 0 || math.IsInf(df, 0) || math.IsInf(ncp, 0) {
		return math.NaN()
	}
	var dist interface {
		CDF(float64) float64
		SF(float64) float64
		LogCDF(float64) float64
		LogSF(float64) float64
		PPF(float64) float64
		ISF(float64) float64
	}
	switch {
	case hasNcp && ncp > 0:
		d, err := stats.NewNoncentralChiSquared(df, ncp, 0, 1)
		if err != nil {
			return math.NaN()
		}
		dist = d
	case df == 0:
		// a point mass at 0, as pgamma with shape 0
		if kind == chisqQ {
			if !validProb(x, logp) {
				return math.NaN()
			}
			return 0
		}
		cdf := 0.0
		if x > 0 {
			cdf = 1
		}
		if !lower {
			cdf = 1 - cdf
		}
		if logp {
			return math.Log(cdf)
		}
		return cdf
	default:
		d, err := stats.NewChiSquared(df, 0, 1)
		if err != nil {
			return math.NaN()
		}
		dist = d
	}
	if kind == chisqP {
		switch {
		case lower && logp:
			return dist.LogCDF(x)
		case lower:
			return dist.CDF(x)
		case logp:
			return dist.LogSF(x)
		}
		return dist.SF(x)
	}
	if !validProb(x, logp) {
		return math.NaN()
	}
	p := x
	if logp {
		p = math.Exp(x)
	}
	if lower {
		return dist.PPF(p)
	}
	return dist.ISF(p)
}

func validProb(p float64, logp bool) bool {
	if logp {
		return p <= 0
	}
	return p >= 0 && p <= 1
}

func init() {
	chisqFn("pchisq", chisqP)
	chisqFn("qchisq", chisqQ)
}
