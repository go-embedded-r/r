package eval

import (
	"math"

	"github.com/go-statistics/stats"
)

// Slice 26 (M4): d, p and q of the continuous families norm, unif, exp, gamma,
// beta, t, f, chisq, lnorm, logis, cauchy and weibull, with log, lower.tail and
// log.p. The numerics are go-statistics/stats (each method has a stated ulp
// budget, listed in the README); this file is the R glue: argument matching
// with the defaults of each help page, recycling, NA, the attributes of the
// first full-length argument, the domain rules (NaN with "NaNs produced"),
// the point masses of the degenerate parameters, and the quantile at a log
// probability.
//
// Limits, listed in the README: ncp for t, beta and F is refused (go-statistics
// has no non-central t, beta or F); a non-central chi-squared with df = 0 is
// NaN (go-statistics/stats#26).

// dist is what every family of stats provides, plus the support.
type dist interface {
	PDF(float64) float64
	LogPDF(float64) float64
	CDF(float64) float64
	LogCDF(float64) float64
	SF(float64) float64
	LogSF(float64) float64
	PPF(float64) float64
	ISF(float64) float64
	Support() (float64, float64)
}

// pointMass is the distribution concentrated on m (a zero scale, a gamma shape
// of 0, ...): d is Inf at m, p a step, q is m.
type pointMass struct{ m float64 }

func (p pointMass) PDF(x float64) float64 {
	if x == p.m {
		return math.Inf(1)
	}
	return 0
}
func (p pointMass) LogPDF(x float64) float64 { return math.Log(p.PDF(x)) }
func (p pointMass) CDF(x float64) float64 {
	if x >= p.m {
		return 1
	}
	return 0
}
func (p pointMass) LogCDF(x float64) float64 { return math.Log(p.CDF(x)) }
func (p pointMass) SF(x float64) float64     { return 1 - p.CDF(x) }
func (p pointMass) LogSF(x float64) float64  { return math.Log(p.SF(x)) }
func (p pointMass) PPF(float64) float64      { return p.m }
func (p pointMass) ISF(float64) float64      { return p.m }
func (p pointMass) Support() (float64, float64) {
	return p.m, p.m
}

// famParam is one parameter of a family: its R name and its default (required
// when req).
type famParam struct {
	name string
	def  float64
	req  bool
}

// family describes one distribution: the parameters in the order of the R
// formals, and build, which turns one row of parameter values into a
// distribution (ok false: the parameters are invalid, NaN).
type family struct {
	name   string
	params []famParam
	// ncp is the index in params of the ncp formal, or -1.
	ncp   int
	build func(p []float64) (dist, bool)
	// fix, when not nil, rewrites the argument values (gamma: rate and scale).
	fix func(s *Session, c *callInfo, a []Value, cols [][]float64) [][]float64
}

// qEnds holds, for the families whose quantile R answers at the probabilities 0
// and 1 before it looks at the parameters, the two ends: qgamma(0, -1) is 0.
var qEnds = map[string][2]float64{
	"norm":  {math.Inf(-1), math.Inf(1)},
	"logis": {math.Inf(-1), math.Inf(1)},
	"t":     {math.Inf(-1), math.Inf(1)},
	"gamma": {0, math.Inf(1)},
	"chisq": {0, math.Inf(1)},
	"lnorm": {0, math.Inf(1)},
}

// endsLaw answers only the ends of the support, for the quantile of invalid
// parameters at the probabilities 0 and 1.
type endsLaw struct{ lo, hi float64 }

func (endsLaw) PDF(float64) float64    { return math.NaN() }
func (endsLaw) LogPDF(float64) float64 { return math.NaN() }
func (endsLaw) CDF(float64) float64    { return math.NaN() }
func (endsLaw) LogCDF(float64) float64 { return math.NaN() }
func (endsLaw) SF(float64) float64     { return math.NaN() }
func (endsLaw) LogSF(float64) float64  { return math.NaN() }
func (endsLaw) PPF(float64) float64    { return math.NaN() }
func (endsLaw) ISF(float64) float64    { return math.NaN() }
func (e endsLaw) Support() (float64, float64) {
	return e.lo, e.hi
}

func mustNew[T dist](d T, err error) (dist, bool) {
	if err != nil {
		return nil, false
	}
	return d, true
}

var families = []family{
	{name: "norm", ncp: -1,
		params: []famParam{{"mean", 0, false}, {"sd", 1, false}},
		build: func(p []float64) (dist, bool) {
			if p[1] < 0 {
				return nil, false
			}
			if p[1] == 0 {
				return massIn{pointMass{p[0]}, math.Inf(-1), math.Inf(1)}, true
			}
			return lsOf(stats.NewNormal(0, 1))(p[0], p[1])
		}},
	{name: "unif", ncp: -1,
		params: []famParam{{"min", 0, false}, {"max", 1, false}},
		build: func(p []float64) (dist, bool) {
			if p[1] < p[0] {
				return nil, false
			}
			if !finite(p[0]) || !finite(p[1]) {
				if p[1] == p[0] {
					return nil, false
				}
				return densZero{}, true
			}
			if p[1] == p[0] {
				return pointNaN{pointMass{p[0]}, false}, true
			}
			return mustNew(stats.NewUniform(p[0], p[1]-p[0]))
		}},
	{name: "exp", ncp: -1,
		params: []famParam{{"rate", 1, false}},
		build: func(p []float64) (dist, bool) {
			if p[0] < 0 {
				return nil, false
			}
			sc := 1 / p[0]
			if sc == 0 {
				return pointNaN{zeroPlus{pointMass{0}}, true}, true
			}
			if math.IsInf(sc, 1) {
				return infScale{nanInf: true}, true
			}
			return mustNew(stats.NewExponential(0, sc))
		}},
	{name: "gamma", ncp: -1,
		params: []famParam{{"shape", 0, true}, {"scale", 1, false}},
		fix:    fixGamma,
		build: func(p []float64) (dist, bool) {
			if p[0] < 0 || p[1] <= 0 {
				return nil, false
			}
			if p[0] == 0 {
				return zeroPlus{pointMass{0}}, true
			}
			if math.IsInf(p[1], 1) {
				return infScale{atZero: p[0] < 1}, true
			}
			if math.IsInf(p[0], 1) {
				return negMass{pointMass{math.Inf(1)}}, true
			}
			return mustNew(stats.NewGamma(p[0], 0, p[1]))
		}},
	{name: "beta", ncp: 2,
		params: []famParam{{"shape1", 0, true}, {"shape2", 0, true}, {"ncp", 0, false}},
		build: func(p []float64) (dist, bool) {
			a, b := p[0], p[1]
			if a < 0 || b < 0 {
				return nil, false
			}
			switch {
			case a == 0 && b == 0:
				return betaMass{both: true}, true
			case a == 0 || a/b == 0:
				return betaMass{m: 0}, true
			case b == 0 || b/a == 0:
				return betaMass{m: 1}, true
			case math.IsInf(a, 1) && math.IsInf(b, 1):
				return betaMass{m: 0.5}, true
			}
			return mustNew(stats.NewBeta(a, b, 0, 1))
		}},
	{name: "t", ncp: 1,
		params: []famParam{{"df", 0, true}, {"ncp", 0, false}},
		build: func(p []float64) (dist, bool) {
			if p[0] <= 0 {
				return nil, false
			}
			if math.IsInf(p[0], 1) {
				return mustNew(stats.NewNormal(0, 1))
			}
			return mustNew(stats.NewStudentT(p[0], 0, 1))
		}},
	{name: "f", ncp: 2,
		params: []famParam{{"df1", 0, true}, {"df2", 0, true}, {"ncp", 0, false}},
		build: func(p []float64) (dist, bool) {
			m, n := p[0], p[1]
			if m <= 0 || n <= 0 {
				return nil, false
			}
			switch {
			case math.IsInf(m, 1) && math.IsInf(n, 1):
				return fMass{pointMass{1}}, true
			case math.IsInf(n, 1):
				return mustNew(stats.NewChiSquared(m, 0, 1/m))
			case math.IsInf(m, 1):
				c, err := stats.NewChiSquared(n, 0, 1)
				if err != nil {
					return nil, false
				}
				return invChi{c, n}, true
			}
			return mustNew(stats.NewF(m, n, 0, 1))
		}},
	{name: "chisq", ncp: 1,
		params: []famParam{{"df", 0, true}, {"ncp", 0, false}},
		build: func(p []float64) (dist, bool) {
			df, ncp := p[0], p[1]
			if df < 0 || ncp < 0 || math.IsInf(df, 0) || math.IsInf(ncp, 0) {
				return nil, false
			}
			if ncp > 0 {
				return mustNew(stats.NewNoncentralChiSquared(df, ncp, 0, 1))
			}
			if df == 0 {
				return zeroPlus{pointMass{0}}, true
			}
			return mustNew(stats.NewChiSquared(df, 0, 1))
		}},
	{name: "lnorm", ncp: -1,
		params: []famParam{{"meanlog", 0, false}, {"sdlog", 1, false}},
		build: func(p []float64) (dist, bool) {
			if p[1] < 0 {
				return nil, false
			}
			if p[1] == 0 {
				return massIn{pointMass{math.Exp(p[0])}, 0, math.Inf(1)}, true
			}
			// exp of the normal law, as R does: accurate in the far tails of the
			// log scale, where go-statistics' log-normal underflows
			n, ok := lsOf(stats.NewNormal(0, 1))(p[0], p[1])
			if !ok {
				return nil, false
			}
			return expOf{n}, true
		}},
	{name: "logis", ncp: -1,
		params: []famParam{{"location", 0, false}, {"scale", 1, false}},
		build: func(p []float64) (dist, bool) {
			if p[1] < 0 {
				return nil, false
			}
			if p[1] == 0 {
				return nanLaw{p[0], math.Inf(-1), math.Inf(1)}, true
			}
			return lsOf(stats.NewLogistic(0, 1))(p[0], p[1])
		}},
	{name: "cauchy", ncp: -1,
		params: []famParam{{"location", 0, false}, {"scale", 1, false}},
		build: func(p []float64) (dist, bool) {
			if p[1] < 0 {
				return nil, false
			}
			if p[1] == 0 {
				return nanLaw{p[0], p[0], p[0]}, true
			}
			return lsOf(stats.NewCauchy(0, 1))(p[0], p[1])
		}},
	{name: "weibull", ncp: -1,
		params: []famParam{{"shape", 0, true}, {"scale", 1, false}},
		build: func(p []float64) (dist, bool) {
			if p[0] <= 0 || p[1] <= 0 {
				return nil, false
			}
			if math.IsInf(p[0], 1) {
				return weibMass{p[1]}, true
			}
			if math.IsInf(p[1], 1) {
				return infScale{atZero: p[0] < 1, nanPos: p[0] < 1}, true
			}
			return mustNew(stats.NewWeibull(p[0], 0, p[1]))
		}},
}

// lsOf returns the constructor of a location-scale law from the standard one:
// z = (x - loc) / scale, with the rules R applies to an infinite location or
// scale (a zero density, a step) that stats, which wants both finite, refuses.
func lsOf[T dist](std T, err error) func(loc, scale float64) (dist, bool) {
	return func(loc, scale float64) (dist, bool) {
		if err != nil || scale < 0 {
			return nil, false
		}
		return locScale{std, loc, scale}, true
	}
}

// locScale is the standard law std moved to loc and scaled by scale (> 0, maybe
// infinite).
type locScale struct {
	std        dist
	loc, scale float64
}

func (l locScale) z(x float64) float64 { return (x - l.loc) / l.scale }

func (l locScale) PDF(x float64) float64 {
	if math.IsInf(l.scale, 1) {
		return 0
	}
	if math.IsInf(x, 0) && x == l.loc {
		return math.NaN()
	}
	z := l.z(x)
	if !finite(z) {
		return 0
	}
	return l.std.PDF(z) / l.scale
}

func (l locScale) LogPDF(x float64) float64 {
	if math.IsInf(l.scale, 1) {
		return math.Inf(-1)
	}
	if math.IsInf(x, 0) && x == l.loc {
		return math.NaN()
	}
	z := l.z(x)
	if !finite(z) {
		return math.Inf(-1)
	}
	return l.std.LogPDF(z) - math.Log(l.scale)
}

// tail answers a cdf-like call: the standard law's f at z, or the value it
// saturates to (below or above the location) when z is not finite.
func (l locScale) tail(x float64, f func(float64) float64, below, above float64) float64 {
	if math.IsInf(x, 0) && x == l.loc {
		return math.NaN()
	}
	z := l.z(x)
	// R: a z that is not finite (also Inf / Inf) is below or above the location
	if !finite(z) {
		if x < l.loc {
			return below
		}
		return above
	}
	return f(z)
}

func (l locScale) CDF(x float64) float64 { return l.tail(x, l.std.CDF, 0, 1) }
func (l locScale) SF(x float64) float64  { return l.tail(x, l.std.SF, 1, 0) }
func (l locScale) LogCDF(x float64) float64 {
	return l.tail(x, l.std.LogCDF, math.Inf(-1), 0)
}
func (l locScale) LogSF(x float64) float64 {
	return l.tail(x, l.std.LogSF, 0, math.Inf(-1))
}
func (l locScale) PPF(p float64) float64 { return l.loc + l.scale*l.std.PPF(p) }
func (l locScale) ISF(q float64) float64 { return l.loc + l.scale*l.std.ISF(q) }
func (l locScale) Support() (float64, float64) {
	return math.Inf(-1), math.Inf(1)
}

// expOf is exp of a law: the log-normal with infinite parameters.
type expOf struct{ n dist }

// infSd is true for a log-normal whose sdlog is infinite: R answers NaN for its
// density at x = Inf.
func (e expOf) infSd(x float64) bool {
	l, ok := e.n.(locScale)
	return ok && math.IsInf(l.scale, 1) && math.IsInf(x, 1)
}

func (e expOf) PDF(x float64) float64 {
	if e.infSd(x) {
		return math.NaN()
	}
	if x <= 0 {
		return 0
	}
	return e.n.PDF(math.Log(x)) / x
}
func (e expOf) LogPDF(x float64) float64 {
	if e.infSd(x) {
		return math.NaN()
	}
	if x <= 0 {
		return math.Inf(-1)
	}
	return e.n.LogPDF(math.Log(x)) - math.Log(x)
}
func (e expOf) CDF(x float64) float64 {
	if x <= 0 {
		return 0
	}
	return e.n.CDF(math.Log(x))
}
func (e expOf) SF(x float64) float64 {
	if x <= 0 {
		return 1
	}
	return e.n.SF(math.Log(x))
}
func (e expOf) LogCDF(x float64) float64 {
	if x <= 0 {
		return math.Inf(-1)
	}
	return e.n.LogCDF(math.Log(x))
}
func (e expOf) LogSF(x float64) float64 {
	if x <= 0 {
		return 0
	}
	return e.n.LogSF(math.Log(x))
}
func (e expOf) PPF(p float64) float64 { return math.Exp(e.n.PPF(p)) }
func (e expOf) ISF(q float64) float64 { return math.Exp(e.n.ISF(q)) }
func (expOf) Support() (float64, float64) {
	return 0, math.Inf(1)
}

// pointNaN is a point mass whose density R answers NaN for (the exponential
// with rate Inf, the uniform with min = max).
type pointNaN struct {
	dist
	neg bool // a log tail of 0 is -0 (exp), not +0 (unif)
}

func (p pointNaN) interior(x float64) bool { return p.neg && x > 0 }

func (pointNaN) PDF(float64) float64    { return math.NaN() }
func (pointNaN) LogPDF(float64) float64 { return math.NaN() }

// Support: the exponential with rate Inf has no upper end (R: qexp(1, Inf) is
// NaN, as Inf * 0 is).
func (p pointNaN) Support() (float64, float64) {
	lo, hi := p.dist.Support()
	if p.neg {
		hi = math.NaN()
	}
	return lo, hi
}

// zeroPlus is the gamma law with shape 0: R puts the mass just above 0, so the
// cdf is 0 at 0 and 1 above.
type zeroPlus struct{ pointMass }

func (zeroPlus) CDF(x float64) float64 {
	if x > 0 {
		return 1
	}
	return 0
}
func (z zeroPlus) LogCDF(x float64) float64 { return math.Log(z.CDF(x)) }
func (z zeroPlus) SF(x float64) float64     { return 1 - z.CDF(x) }
func (z zeroPlus) LogSF(x float64) float64  { return math.Log(z.SF(x)) }
func (zeroPlus) Support() (float64, float64) {
	return 0, math.Inf(1)
}

// negMass is a point mass whose log tails of 0 are -0 in R (gamma, shape Inf).
type negMass struct{ pointMass }

func (negMass) interior(x float64) bool { return x > 0 }

// R's density of the gamma law with an infinite shape is 0 everywhere, x = Inf
// included.
func (negMass) PDF(float64) float64    { return 0 }
func (negMass) LogPDF(float64) float64 { return math.Inf(-1) }
func (negMass) Support() (float64, float64) {
	return 0, math.Inf(1)
}

// massIn is a point mass (a scale 0) whose quantiles at the probabilities 0 and
// 1 are the ends lo and hi of the support of the family: R answers those
// before it looks at the scale.
type massIn struct {
	pointMass
	lo, hi float64
}

func (m massIn) Support() (float64, float64) { return m.lo, m.hi }

// densZero is the uniform law with an infinite end: R answers a zero density
// and NaN for the distribution and the quantiles.
type densZero struct{}

func (densZero) PDF(float64) float64    { return 0 }
func (densZero) LogPDF(float64) float64 { return math.Inf(-1) }
func (densZero) CDF(float64) float64    { return math.NaN() }
func (densZero) LogCDF(float64) float64 { return math.NaN() }
func (densZero) SF(float64) float64     { return math.NaN() }
func (densZero) LogSF(float64) float64  { return math.NaN() }
func (densZero) PPF(float64) float64    { return math.NaN() }
func (densZero) ISF(float64) float64    { return math.NaN() }
func (densZero) Support() (float64, float64) {
	return math.NaN(), math.NaN()
}

func (pointNaN) exact() {}
func (zeroPlus) exact() {}
func (densZero) exact() {}

// nanLaw is the logistic or Cauchy law with scale 0: R answers NaN for the
// density and the distribution, and the location for a quantile (lo and hi
// are the quantiles of probability 0 and 1).
type nanLaw struct{ loc, lo, hi float64 }

func (nanLaw) PDF(float64) float64    { return math.NaN() }
func (nanLaw) LogPDF(float64) float64 { return math.NaN() }
func (nanLaw) CDF(float64) float64    { return math.NaN() }
func (nanLaw) LogCDF(float64) float64 { return math.NaN() }
func (nanLaw) SF(float64) float64     { return math.NaN() }
func (nanLaw) LogSF(float64) float64  { return math.NaN() }
func (n nanLaw) PPF(float64) float64  { return n.loc }
func (n nanLaw) ISF(float64) float64  { return n.loc }
func (n nanLaw) Support() (float64, float64) {
	return n.lo, n.hi
}
func (nanLaw) exact() {}

// infScale is the exponential law with rate 0 (scale Inf): no density, no mass
// anywhere finite.
type infScale struct {
	nanInf bool // the density at x = Inf is NaN (exp), not 0 (gamma)
	atZero bool // the density at x = 0 is Inf (a shape below 1), not 0
	nanPos bool // the density at a finite x > 0 is NaN (Weibull, shape below 1)
}

func (i infScale) PDF(x float64) float64 {
	switch {
	case i.nanInf && math.IsInf(x, 1):
		return math.NaN()
	case i.atZero && x == 0:
		return math.Inf(1)
	case i.nanPos && x > 0 && finite(x):
		return math.NaN()
	}
	return 0
}
func (i infScale) LogPDF(x float64) float64 {
	switch {
	case i.nanInf && math.IsInf(x, 1):
		return math.NaN()
	case i.atZero && x == 0:
		return math.Inf(1)
	case i.nanPos && x > 0 && finite(x):
		return math.NaN()
	}
	return math.Inf(-1)
}

// the tails: R answers NaN at x = Inf (Inf / Inf), a constant elsewhere
func infNaN(x, v float64) float64 {
	if math.IsInf(x, 1) {
		return math.NaN()
	}
	return v
}
func (infScale) CDF(x float64) float64    { return infNaN(x, 0) }
func (infScale) LogCDF(x float64) float64 { return infNaN(x, math.Inf(-1)) }
func (infScale) SF(x float64) float64     { return infNaN(x, 1) }
func (infScale) LogSF(x float64) float64  { return infNaN(x, 0) }
func (infScale) PPF(float64) float64      { return math.Inf(1) }
func (infScale) ISF(float64) float64      { return math.Inf(1) }
func (infScale) Support() (float64, float64) {
	return 0, math.Inf(1)
}

// fMass is the F law with both degrees of freedom infinite: a mass at 1 whose
// distribution function R answers 1/2 at 1 itself.
type fMass struct{ pointMass }

func (fMass) CDF(x float64) float64 {
	switch {
	case x < 1:
		return 0
	case x == 1:
		return 0.5
	}
	return 1
}
func (f fMass) SF(x float64) float64     { return 1 - f.CDF(x) }
func (f fMass) LogCDF(x float64) float64 { return math.Log(f.CDF(x)) }
func (f fMass) LogSF(x float64) float64  { return math.Log(f.SF(x)) }
func (fMass) Support() (float64, float64) {
	return 0, math.Inf(1)
}

// weibMass is the Weibull law with an infinite shape: R evaluates the formulas
// at the limit, so the density is Inf at x = scale, NaN inside (0, Inf) and 0
// below 0 and at Inf, the distribution function steps from 0 to 1 - 1/e at the
// scale and to 1 above it, and the quantile of every probability inside (0, 1)
// is the scale.
type weibMass struct{ scale float64 }

func (w weibMass) PDF(x float64) float64 {
	switch {
	case x < 0 || math.IsInf(x, 1):
		return 0
	case x == w.scale:
		return math.Inf(1)
	}
	return math.NaN()
}
func (w weibMass) LogPDF(x float64) float64 {
	if p := w.PDF(x); p == 0 {
		return math.Inf(-1)
	} else if p > 0 {
		return p
	}
	return math.NaN()
}
func (w weibMass) CDF(x float64) float64 {
	switch {
	case x < w.scale:
		return 0
	case x == w.scale:
		return -math.Expm1(-1)
	}
	return 1
}
func (w weibMass) SF(x float64) float64 {
	switch {
	case x < w.scale:
		return 1
	case x == w.scale:
		return math.Exp(-1)
	}
	return 0
}
func (w weibMass) LogCDF(x float64) float64 {
	switch {
	case x < w.scale:
		return math.Inf(-1)
	case x == w.scale:
		return math.Log1p(-math.Exp(-1))
	}
	return 0
}
func (w weibMass) LogSF(x float64) float64 {
	switch {
	case x < w.scale:
		return 0
	case x == w.scale:
		return -1
	}
	return math.Inf(-1)
}
func (w weibMass) PPF(float64) float64 { return w.scale }
func (w weibMass) ISF(float64) float64 { return w.scale }
func (weibMass) Support() (float64, float64) {
	return 0, math.Inf(1)
}
func (weibMass) exact() {}

// interior: R's log tail is -0 where it is computed from the formula (0 < x),
// +0 where the answer is returned early (x <= 0).
func (weibMass) interior(x float64) bool { return x > 0 }

// betaMass is a degenerate beta law (a shape 0 or infinite): its cdf is the one
// R answers, which for a mass at 0 or 1 does not look at the range of x.
// both is beta(0, 0), half the mass at each end.
type betaMass struct {
	m    float64
	both bool
}

func (b betaMass) PDF(x float64) float64 {
	if x == b.m || b.both && (x == 0 || x == 1) {
		return math.Inf(1)
	}
	return 0
}
func (b betaMass) LogPDF(x float64) float64 { return math.Log(b.PDF(x)) }
func (b betaMass) CDF(x float64) float64 {
	switch {
	case x >= 1:
		return 1
	case b.both:
		return 0.5
	case b.m == 0:
		return 1
	case b.m == 1:
		return 0
	case x < 0.5:
		return 0
	}
	return 1
}
func (b betaMass) LogCDF(x float64) float64 { return math.Log(b.CDF(x)) }
func (b betaMass) SF(x float64) float64     { return 1 - b.CDF(x) }
func (b betaMass) LogSF(x float64) float64  { return math.Log(b.SF(x)) }
func (b betaMass) PPF(p float64) float64 {
	if !b.both {
		return b.m
	}
	if p < 0.5 {
		return 0
	}
	if p > 0.5 {
		return 1
	}
	return 0.5
}
func (b betaMass) ISF(q float64) float64 { return b.PPF(1 - q) }
func (betaMass) Support() (float64, float64) {
	return 0, 1
}

// invChi is the F law with infinite numerator df: df2 / chisq(df2), the
// reciprocal law. Only the cdf side is needed for p, d and q of R.
type invChi struct {
	c stats.ChiSquared
	n float64
}

func (v invChi) PDF(x float64) float64 {
	if x <= 0 {
		return 0
	}
	return v.c.PDF(v.n/x) * v.n / (x * x)
}
func (v invChi) LogPDF(x float64) float64 {
	if x <= 0 {
		return math.Inf(-1)
	}
	return v.c.LogPDF(v.n/x) + math.Log(v.n) - 2*math.Log(x)
}
func (v invChi) CDF(x float64) float64 {
	if x <= 0 {
		return 0
	}
	return v.c.SF(v.n / x)
}
func (v invChi) LogCDF(x float64) float64 {
	if x <= 0 {
		return math.Inf(-1)
	}
	return v.c.LogSF(v.n / x)
}
func (v invChi) SF(x float64) float64 {
	if x <= 0 {
		return 1
	}
	return v.c.CDF(v.n / x)
}
func (v invChi) LogSF(x float64) float64 {
	if x <= 0 {
		return 0
	}
	return v.c.LogCDF(v.n / x)
}
func (v invChi) PPF(p float64) float64 { return v.n / v.c.ISF(p) }
func (v invChi) ISF(q float64) float64 { return v.n / v.c.PPF(q) }
func (invChi) Support() (float64, float64) {
	return 0, math.Inf(1)
}

// fixGamma resolves rate and scale of the gamma functions: cols holds shape,
// rate, scale in the order of the R formals.
func fixGamma(s *Session, c *callInfo, a []Value, cols [][]float64) [][]float64 {
	rate, scale := a[1], a[2]
	if rate != nil && scale != nil {
		r, sc := cols[1], cols[2]
		n := max(len(r), len(sc))
		if len(r) > 0 && len(sc) > 0 {
			for i := 0; i < n; i++ {
				if math.Abs(r[i%len(r)]*sc[i%len(sc)]-1) >= 1e-15 {
					s.errorf(c.call, "specify 'rate' or 'scale' but not both")
				}
			}
			s.warningf(c.call, "specify 'rate' or 'scale' but not both")
		}
	}
	var sc []float64
	switch {
	case scale != nil:
		sc = cols[2]
	case rate != nil:
		sc = make([]float64, len(cols[1]))
		for i, r := range cols[1] {
			sc[i] = 1 / r
		}
	default:
		sc = []float64{1}
	}
	return [][]float64{cols[0], sc}
}

type distKind int

const (
	kindD distKind = iota
	kindP
	kindQ
)

func init() {
	for i := range families {
		f := &families[i]
		regDist("d"+f.name, f, kindD)
		regDist("p"+f.name, f, kindP)
		regDist("q"+f.name, f, kindQ)
	}
}

func regDist(name string, f *family, kind distKind) {
	regBuiltin(name, func(s *Session, c *callInfo) Value { return s.distCall(c, f, kind) })
}

// formals of one function: the first argument, the parameters (gamma lists
// rate and scale for its single scale parameter), the flags.
func (f *family) formals(kind distKind) []string {
	first := map[distKind]string{kindD: "x", kindP: "q", kindQ: "p"}[kind]
	out := []string{first}
	if f.name == "gamma" {
		out = append(out, "shape", "rate", "scale")
	} else if f.name == "exp" {
		out = append(out, "rate")
	} else {
		for _, p := range f.params {
			out = append(out, p.name)
		}
	}
	if kind == kindD {
		return append(out, "log")
	}
	return append(out, "lower.tail", "log.p")
}

func (s *Session) distCall(c *callInfo, f *family, kind distKind) Value {
	formals := f.formals(kind)
	a := s.matchFormals(c, formals...)
	if a[0] == nil {
		s.errorf(c.call, "argument \"%s\" is missing, with no default", formals[0])
	}
	nflag := 1
	if kind != kindD {
		nflag = 2
	}
	vals := a[1 : len(a)-nflag]
	var lower, logp bool
	if kind == kindD {
		logp = s.distFlag(c, a[len(a)-1], false, "log")
		lower = true
	} else {
		lower = s.distFlag(c, a[len(a)-2], true, "lower.tail")
		logp = s.distFlag(c, a[len(a)-1], false, "log.p")
	}
	// the parameter columns, defaults filled in
	np := len(f.params)
	if f.name == "gamma" {
		np = 3
	}
	pnames := make([]famParam, 0, np)
	if f.name == "gamma" {
		pnames = append(pnames, f.params[0], famParam{"rate", 1, false}, famParam{"scale", 1, false})
	} else {
		pnames = append(pnames, f.params...)
	}
	for i, p := range pnames {
		if p.req && vals[i] == nil {
			s.errorf(c.call, "argument \"%s\" is missing, with no default", p.name)
		}
	}
	all := append([]Value{a[0]}, vals...)
	cols := make([][]float64, len(all))
	for i, v := range all {
		if v == nil {
			cols[i] = []float64{pnames[i-1].def}
			continue
		}
		switch v.(type) {
		case *Double, *Integer, *Logical:
		default:
			s.errorf(c.call, "Non-numeric argument to mathematical function")
		}
		cols[i] = asDoubles(v)
	}
	for i, v := range all {
		if v != nil && len(cols[i]) == 0 {
			return &Double{V: []float64{}}
		}
	}
	pcols := cols[1:]
	if f.fix != nil {
		pcols = f.fix(s, c, vals, pcols)
	}
	if f.ncp >= 0 && vals[f.ncp] != nil && f.name != "chisq" {
		for _, v := range cols[1+f.ncp] {
			if v != 0 && !math.IsNaN(v) {
				s.errorf(c.call, "%s: a non-central distribution (ncp != 0) is not implemented: go-statistics/stats has no non-central %s", name3(kind)+f.name, f.name)
			}
		}
	}
	n := len(cols[0])
	for _, pc := range pcols {
		n = max(n, len(pc))
	}
	out := make([]float64, n)
	row := make([]float64, len(pcols))
	nan := false
	for i := range out {
		if i&1023 == 0 {
			s.checkInterrupt()
		}
		x := cols[0][i%len(cols[0])]
		na, nn := isNAReal(x), math.IsNaN(x)
		for j, pc := range pcols {
			row[j] = pc[i%len(pc)]
			na = na || isNAReal(row[j])
			nn = nn || math.IsNaN(row[j])
		}
		if na {
			out[i] = naReal
			continue
		}
		if nn {
			out[i] = math.NaN()
			continue
		}
		out[i] = f.value(kind, x, row, lower, logp)
		if math.IsNaN(out[i]) {
			nan = true
		}
	}
	if nan {
		s.warningf(c.call, "NaNs produced")
	}
	var attrs []Attr
	if len(cols[0]) == n && a[0] != nil {
		attrs = attrsOf(a[0])
	} else {
		for j, pc := range pcols {
			if len(pc) == n {
				if v := paramValue(f, vals, j); v != nil {
					attrs = attrsOf(v)
					break
				}
			}
		}
	}
	return &Double{V: out, Attr: attrs}
}

// distFlag reads log, lower.tail or log.p: R takes the integer value of the
// first element, and NA is not zero, so it is TRUE.
func (s *Session) distFlag(c *callInfo, v Value, def bool, what string) bool {
	switch x := v.(type) {
	case *Logical:
		if len(x.V) >= 1 && x.V[0] == naInt {
			return true
		}
	case *Integer:
		if len(x.V) >= 1 && x.V[0] == naInt {
			return true
		}
	case *Double:
		if len(x.V) >= 1 && math.IsNaN(x.V[0]) {
			return true
		}
	}
	return s.boolArg(c, v, def, what)
}

func name3(kind distKind) string {
	return map[distKind]string{kindD: "d", kindP: "p", kindQ: "q"}[kind]
}

// paramValue is the argument Value that supplied parameter column j.
func paramValue(f *family, vals []Value, j int) Value {
	if f.name == "gamma" && j == 1 {
		if vals[2] != nil {
			return vals[2]
		}
		return vals[1]
	}
	if j < len(vals) {
		return vals[j]
	}
	return nil
}

// value is one d, p or q (no NaN argument).
func (f *family) value(kind distKind, x float64, row []float64, lower, logp bool) float64 {
	d, ok := f.build(row)
	if !ok {
		if e, has := qEnds[f.name]; has && kind == kindQ {
			return quantile(endsLaw{e[0], e[1]}, x, lower, logp)
		}
		return math.NaN()
	}
	lo, hi := d.Support()
	// go-statistics' log densities and log tails are the log of the plain
	// value in places: -Inf where the plain value underflows although x is
	// inside the support and the true log is finite. R answers a number
	// there; NaN with the warning is the honest answer (README, limits).
	_, isExact := d.(exact)
	logOK := func(v float64) float64 {
		if !isExact && math.IsInf(v, -1) && x > lo && x < hi {
			return math.NaN()
		}
		return v
	}
	inside := func(lower bool) bool {
		if zs, ok := d.(zeroSigner); ok {
			return zs.interior(x)
		}
		return !isExact && (x > lo && x < hi || lower && x == hi && math.IsInf(hi, 1) && (f.name == "exp" || f.name == "weibull"))
	}
	switch kind {
	case kindD:
		if logp {
			v := logOK(d.LogPDF(x))
			if v == 0 && (f.name == "unif" || f.name == "gamma") {
				return math.Copysign(0, -1) // R's log density is -log(scale) - ...: -0
			}
			return v
		}
		return d.PDF(x)
	case kindP:
		switch {
		case lower && logp:
			return negZero(logOK(d.LogCDF(x)), inside(true))
		case lower:
			return d.CDF(x)
		case logp:
			return negZero(logOK(d.LogSF(x)), inside(false))
		}
		return d.SF(x)
	}
	return quantile(d, x, lower, logp)
}

func validProb(p float64, logp bool) bool {
	if logp {
		return p <= 0
	}
	return p >= 0 && p <= 1
}

// quantile is q of R for the probability p (a log probability when logp).
func quantile(d dist, p float64, lower, logp bool) float64 {
	if !validProb(p, logp) {
		return math.NaN()
	}
	lo, hi := d.Support()
	// the probability of the end of the support the answer is next to
	var zero, one bool
	if logp {
		zero, one = math.IsInf(p, -1), p == 0
	} else {
		zero, one = p == 0, p == 1
	}
	if lower && zero || !lower && one {
		return lo
	}
	if lower && one || !lower && zero {
		return hi
	}
	if !logp {
		if lower {
			return d.PPF(p)
		}
		return d.ISF(p)
	}
	// log scale: the tail the probability is small in is the one asked for
	// directly when p < -ln 2, the other tail goes through -expm1(p).
	small := p < -math.Ln2
	switch {
	case small && math.Exp(p) >= 0x1p-1000:
		if lower {
			return d.PPF(math.Exp(p))
		}
		return d.ISF(math.Exp(p))
	case !small:
		if lower {
			return d.ISF(-math.Expm1(p))
		}
		return d.PPF(-math.Expm1(p))
	}
	// the laws that are a transform of another: the deep quantile of the law
	// underneath, transformed (the overflow of the scaling is then the answer)
	switch v := d.(type) {
	case pointMass, massIn, nanLaw, zeroPlus, negMass, pointNaN, endsLaw, fMass, weibMass:
		return v.PPF(0.5)
	case locScale:
		if finite(v.scale) {
			return v.loc + v.scale*quantile(v.std, p, lower, logp)
		}
	case expOf:
		return math.Exp(quantile(v.n, p, lower, logp))
	}
	return deepQuantile(d, p, lower)
}

// zeroSigner is implemented by the laws that know, for a log tail that is
// exactly 0, whether R's answer is -0 (x inside the range where the tail is
// computed) or +0 (x saturates the tail).
type zeroSigner interface{ interior(x float64) bool }

func (l locScale) interior(x float64) bool { return finite(x) && finite(l.z(x)) }
func (e expOf) interior(x float64) bool {
	z, ok := e.n.(zeroSigner)
	return ok && x > 0 && finite(x) && z.interior(math.Log(x))
}

// deepQuantile solves LogCDF(x) = p (lower) or LogSF(x) = p when exp(p) is
// too small for the plain quantile: a bisection over the float64 order of x,
// on the log probability, which is monotone in x.
func deepQuantile(d dist, p float64, lower bool) float64 {
	lo, hi := d.Support()
	// g is increasing in x and crosses zero at the answer
	g := func(x float64) float64 {
		if lower {
			return d.LogCDF(x) - p
		}
		return p - d.LogSF(x)
	}
	// start from the quantile at the smallest normal probability and walk
	// away from the centre until the sign changes
	var a, b float64
	if lower {
		b = d.PPF(0x1p-1000)
		if !(b > lo) {
			b = lo + (hi-lo)/2
			if math.IsInf(lo, -1) || math.IsInf(hi, 1) {
				b = 0
			}
		}
		a = lo
		if math.IsInf(lo, -1) {
			a = b - math.Max(1, math.Abs(b))
			for k := 0; g(a) > 0 && k < 1100; k++ {
				a = b - 2*(b-a)
			}
		}
	} else {
		a = d.ISF(0x1p-1000)
		if !(a < hi) {
			a = 0
		}
		b = hi
		if math.IsInf(hi, 1) {
			b = a + math.Max(1, math.Abs(a))
			for k := 0; g(b) < 0 && k < 1100; k++ {
				b = a + 2*(b-a)
			}
		}
	}
	if math.IsInf(a, 0) {
		a = math.Copysign(math.MaxFloat64, a)
	}
	if math.IsInf(b, 0) {
		b = math.Copysign(math.MaxFloat64, b)
	}
	// the root is beyond the largest finite value: it is infinite
	if lower && math.IsInf(lo, -1) && g(a) > 0 {
		return lo
	}
	if !lower && math.IsInf(hi, 1) && g(b) < 0 {
		return hi
	}
	ia, ib := orderedBits(a), orderedBits(b)
	for ib-ia > 1 {
		m := ia + (ib-ia)/2
		if g(fromOrderedBits(m)) < 0 {
			ia = m
		} else {
			ib = m
		}
	}
	// the root lies between the end of the support and the next float: the
	// end is the nearest answer
	if lower && !math.IsInf(lo, 0) && fromOrderedBits(ia) == lo {
		return lo
	}
	if !lower && !math.IsInf(hi, 0) && fromOrderedBits(ib) == hi {
		return hi
	}
	res := fromOrderedBits(ib)
	// a log probability that stats cannot follow (it underflows to -Inf before
	// the true value does) leaves the bisection on the edge of the finite
	// values, not on a root: refuse with NaN rather than answer from it
	if r := g(res); math.IsNaN(r) || math.Abs(r) > 1e-9*math.Abs(p) {
		return math.NaN()
	}
	return res
}

// orderedBits maps float64 to int64 monotonically (so that bisection halves the
// number of representable values, not the interval).
func orderedBits(x float64) int64 {
	b := int64(math.Float64bits(x))
	if b < 0 {
		return math.MinInt64 - b
	}
	return b
}

func fromOrderedBits(i int64) float64 {
	if i < 0 {
		i = math.MinInt64 - i
	}
	return math.Float64frombits(uint64(i))
}

// chisqKind and chisqValue keep the one-value entry the chi-squared test and
// its tests use: pchisq (chisqP, x a quantile) or qchisq (chisqQ).
type chisqKind int

const (
	chisqP chisqKind = iota
	chisqQ
)

func chisqValue(kind chisqKind, x, df, ncp float64, hasNcp, lower, logp bool) float64 {
	if !hasNcp {
		ncp = 0
	}
	k := kindP
	if kind == chisqQ {
		k = kindQ
	}
	for i := range families {
		if families[i].name == "chisq" {
			return families[i].value(k, x, []float64{df, ncp}, lower, logp)
		}
	}
	return math.NaN()
}

// exact marks the laws of this file whose log values are computed here, never
// the underflowing log of go-statistics' plain value.
type exact interface{ exact() }

func (pointMass) exact() {}
func (betaMass) exact()  {}
func (infScale) exact()  {}
func (locScale) exact()  {}
func (expOf) exact()     {}

// negZero turns a log probability of exactly 0 into -0, the sign R's log1p
// based formulas leave, when the point is inside the support (inside); at or
// beyond an end of the support, and for the degenerate laws, R answers +0.
func negZero(v float64, inside bool) float64 {
	if v == 0 && inside {
		return math.Copysign(0, -1)
	}
	if v == 0 {
		return 0
	}
	return v
}
