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
				return pointMass{p[0]}, true
			}
			return mustNew(stats.NewNormal(p[0], p[1]))
		}},
	{name: "unif", ncp: -1,
		params: []famParam{{"min", 0, false}, {"max", 1, false}},
		build: func(p []float64) (dist, bool) {
			if !finite(p[0]) || !finite(p[1]) || p[1] < p[0] {
				return nil, false
			}
			if p[1] == p[0] {
				return pointMass{p[0]}, true
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
			if math.IsInf(sc, 1) {
				return nil, false
			}
			if sc == 0 {
				return pointMass{0}, true
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
				return pointMass{0}, true
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
				return betaMass{}, true
			case a == 0 || a/b == 0:
				return pointMass{0}, true
			case b == 0 || b/a == 0:
				return pointMass{1}, true
			case math.IsInf(a, 1) && math.IsInf(b, 1):
				return pointMass{0.5}, true
			case math.IsInf(a, 1):
				return pointMass{1}, true
			case math.IsInf(b, 1):
				return pointMass{0}, true
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
				return pointMass{1}, true
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
				return pointMass{0}, true
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
				return pointMass{math.Exp(p[0])}, true
			}
			return mustNew(stats.NewLogNormal(p[1], 0, math.Exp(p[0])))
		}},
	{name: "logis", ncp: -1,
		params: []famParam{{"location", 0, false}, {"scale", 1, false}},
		build: func(p []float64) (dist, bool) {
			if p[1] < 0 {
				return nil, false
			}
			if p[1] == 0 {
				return pointMass{p[0]}, true
			}
			return mustNew(stats.NewLogistic(p[0], p[1]))
		}},
	{name: "cauchy", ncp: -1,
		params: []famParam{{"location", 0, false}, {"scale", 1, false}},
		build: func(p []float64) (dist, bool) {
			if p[1] < 0 {
				return nil, false
			}
			if p[1] == 0 {
				return pointMass{p[0]}, true
			}
			return mustNew(stats.NewCauchy(p[0], p[1]))
		}},
	{name: "weibull", ncp: -1,
		params: []famParam{{"shape", 0, true}, {"scale", 1, false}},
		build: func(p []float64) (dist, bool) {
			if p[0] <= 0 || p[1] <= 0 {
				return nil, false
			}
			return mustNew(stats.NewWeibull(p[0], 0, p[1]))
		}},
}

// betaMass is the beta(0, 0) law: half the mass at 0, half at 1.
type betaMass struct{}

func (betaMass) PDF(x float64) float64 {
	if x == 0 || x == 1 {
		return math.Inf(1)
	}
	return 0
}
func (b betaMass) LogPDF(x float64) float64 { return math.Log(b.PDF(x)) }
func (betaMass) CDF(x float64) float64 {
	switch {
	case x < 0:
		return 0
	case x < 1:
		return 0.5
	}
	return 1
}
func (b betaMass) LogCDF(x float64) float64 { return math.Log(b.CDF(x)) }
func (b betaMass) SF(x float64) float64     { return 1 - b.CDF(x) }
func (b betaMass) LogSF(x float64) float64  { return math.Log(b.SF(x)) }
func (betaMass) PPF(p float64) float64 {
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
		logp = s.boolArg(c, a[len(a)-1], false, "log")
		lower = true
	} else {
		lower = s.boolArg(c, a[len(a)-2], true, "lower.tail")
		logp = s.boolArg(c, a[len(a)-1], false, "log.p")
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
		return math.NaN()
	}
	switch kind {
	case kindD:
		if logp {
			return d.LogPDF(x)
		}
		return d.PDF(x)
	case kindP:
		switch {
		case lower && logp:
			return d.LogCDF(x)
		case lower:
			return d.CDF(x)
		case logp:
			return d.LogSF(x)
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
	return deepQuantile(d, p, lower)
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
			a = b - 1
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
			b = a + 1
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
	ia, ib := orderedBits(a), orderedBits(b)
	for ib-ia > 1 {
		m := ia + (ib-ia)/2
		if g(fromOrderedBits(m)) < 0 {
			ia = m
		} else {
			ib = m
		}
	}
	return fromOrderedBits(ib)
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
