package eval

// Slice 25: the special functions. gamma, lgamma, digamma, trigamma, beta,
// lbeta, psigamma and the four Bessel functions run on go-numerics/special; the
// edge cases and the warnings are those the oracle shows. exp and the finite
// core of ^ are rounded through double-double (go-numerics/special/dd), so
// their last bit does not depend on the platform's libm (go-embedded-r/r#14).

import (
	"math"
	"math/big"

	"github.com/go-numerics/special"
	ddp "github.com/go-numerics/special/dd"
)

// cG is C's "%g" for a double, as R's warning texts print it.
func cG(v float64) string {
	switch {
	case math.IsInf(v, 1):
		return "inf"
	case math.IsInf(v, -1):
		return "-inf"
	case math.IsNaN(v):
		return "nan"
	}
	return formatG(v, 6, false, false)
}

// isNonPosInt reports x in {0, -1, -2, ..., -Inf}.
func isNonPosInt(x float64) bool { return x <= 0 && x == math.Floor(x) }

func gammaFn(x float64) float64 {
	switch {
	case math.IsNaN(x):
		return x
	case x == 0 || (x < 0 && x == math.Floor(x)):
		return math.NaN() // the poles, and -Inf
	case math.IsInf(x, 1):
		return x
	}
	return special.Gamma(x)
}

func lgammaFn(x float64) float64 {
	switch {
	case math.IsNaN(x):
		return x
	case isNonPosInt(x), math.IsInf(x, 1):
		return math.Inf(1)
	}
	return special.LogGamma(x)
}

func digammaFn(x float64) float64 {
	switch {
	case math.IsNaN(x):
		return x
	case isNonPosInt(x):
		return math.NaN()
	case math.IsInf(x, 1):
		return x
	}
	return special.Digamma(x)
}

func trigammaFn(x float64) float64 {
	switch {
	case math.IsNaN(x):
		return x
	case isNonPosInt(x):
		return math.Inf(1)
	case math.IsInf(x, 1):
		return 0
	}
	return overflowNaN(special.Polygamma(1, x))
}

// overflowNaN: R reports a polygamma that overflows (x finite, nonzero and
// tiny) as NaN, not as the infinity it is.
func overflowNaN(r float64) float64 {
	if math.IsInf(r, 0) {
		return math.NaN()
	}
	return r
}

// psigammaN is the n-th derivative of digamma, 0 <= n <= 100.
func psigammaN(x float64, n int) float64 {
	switch n {
	case 0:
		return digammaFn(x)
	case 1:
		return trigammaFn(x)
	}
	switch {
	case math.IsNaN(x):
		return x
	case isNonPosInt(x):
		if n&1 == 1 {
			return math.Inf(1)
		}
		return math.NaN()
	case math.IsInf(x, 1):
		return 0
	}
	return overflowNaN(special.Polygamma(n, x))
}

// betaFn and lbetaFn: Beta and its logarithm, with R's edges: a negative
// argument is NaN, a zero one Inf, an infinite one 0 (-Inf for the log).
func betaFn(a, b float64) float64 {
	switch {
	case math.IsNaN(a) || math.IsNaN(b):
		return a + b
	case a < 0 || b < 0:
		return math.NaN()
	case a == 0 || b == 0:
		return math.Inf(1)
	case math.IsInf(a, 0) || math.IsInf(b, 0):
		return 0
	}
	return special.Beta(a, b)
}

func lbetaFn(a, b float64) float64 {
	switch {
	case math.IsNaN(a) || math.IsNaN(b):
		return a + b
	case a < 0 || b < 0:
		return math.NaN()
	case a == 0 || b == 0:
		return math.Inf(1)
	case math.IsInf(a, 0) || math.IsInf(b, 0):
		return math.Inf(-1)
	}
	return special.LogBeta(a, b)
}

// Bessel functions. The limits are the ones the oracle shows: nu beyond them,
// or an argument beyond the range of the algorithm, gives the warning R gives.

const (
	besselMaxNuJ = 1e9 // max_alpha_j
	besselMaxNuY = 1e7 // max_alpha_y
	besselMaxXJ  = 1e24
	besselMaxXY  = 1e8
	besselMaxXI  = 1e5 // beyond it the scaled I is 0

	besselHankelMin = 1e8 // beyond it J comes from jvHankel, not special.Jv (see its comment)
)

type warnf func(format string, a ...any)

func besselJ(x, nu float64, w warnf) float64 {
	switch {
	case math.IsNaN(x) || math.IsNaN(nu):
		return x + nu
	case x < 0:
		w("value out of range in 'bessel_j'")
		return math.NaN()
	}
	if nu < 0 {
		a := -nu
		var r float64
		if a-math.Floor(a) != 0.5 {
			r = besselJ(x, a, w) * cospi(a)
		}
		if a != math.Floor(a) {
			r -= besselY(x, a, w) * sinpi(a)
		}
		return r
	}
	switch {
	case nu > besselMaxNuJ:
		w("besselJ(x, nu): nu=%s > max_alpha_j (= 1e+09): too large for bessel_j() algorithm", cG(nu))
		return math.NaN()
	case x > besselMaxXJ:
		w("value out of range in 'J_bessel'")
		return 0
	case x == 0:
		if nu == 0 {
			return 1
		}
		return 0
	}
	var r float64
	hankel := false
	if x > besselHankelMin {
		r, hankel = jvHankel(nu, x)
	}
	if !hankel {
		r = special.Jv(nu, x)
	}
	if r == 0 && nu > 0 {
		w("bessel_j(%s,nu=%s): precision lost in result", cG(x), cG(nu))
	}
	return r
}

func besselY(x, nu float64, w warnf) float64 {
	switch {
	case math.IsNaN(x) || math.IsNaN(nu):
		return x + nu
	case x < 0:
		w("value out of range in 'bessel_y'")
		return math.NaN()
	}
	if nu < 0 {
		a := -nu
		if a > besselMaxNuY {
			w("besselY(x, -nu): nu = %s > 1e+07 =: max_alpha_y, too large for bessel_y() algorithm", cG(a))
			return math.NaN()
		}
		var r float64
		if a-math.Floor(a) != 0.5 {
			r = besselY(x, a, w) * cospi(a)
		}
		if a != math.Floor(a) {
			r += besselJ(x, a, w) * sinpi(a)
		}
		return r
	}
	switch {
	case nu > besselMaxNuY:
		w("besselY(x, nu): nu=%s > 1e+07 =: max_alpha_y, too large for bessel_y() algorithm", cG(nu))
		return math.NaN()
	case x > besselMaxXY:
		return 0
	case x == 0:
		return math.Inf(-1)
	}
	r := special.Yv(nu, x)
	if math.IsInf(r, 0) && nu > 0 {
		w("bessel_y(%s,nu=%s): precision lost in result", cG(x), cG(nu))
	}
	return r
}

func besselI(x, nu float64, scaled bool, w warnf) float64 {
	switch {
	case math.IsNaN(x) || math.IsNaN(nu):
		return x + nu
	case x < 0:
		w("value out of range in 'bessel_i'")
		return math.NaN()
	}
	if nu < 0 {
		a := -nu
		r := besselI(x, a, scaled, w)
		if a != math.Floor(a) {
			k := besselK(x, a, scaled, w)
			if scaled {
				k *= 2 * rexp(-2*x)
			} else {
				k *= 2
			}
			r += k / math.Pi * sinpi(a)
		}
		return r
	}
	switch {
	case x == 0:
		if nu == 0 {
			return 1
		}
		return 0
	case math.IsInf(x, 1):
		if scaled {
			return 0
		}
		return x
	case scaled && x > besselMaxXI:
		return 0
	}
	var r float64
	if scaled {
		r = special.Ive(nu, x)
	} else {
		r = special.Iv(nu, x)
	}
	if r == 0 && nu > 0 {
		w("bessel_i(%s,nu=%s): precision lost in result", cG(x), cG(nu))
	}
	return r
}

func besselK(x, nu float64, scaled bool, w warnf) float64 {
	switch {
	case math.IsNaN(x) || math.IsNaN(nu):
		return x + nu
	case x < 0:
		w("value out of range in 'bessel_k'")
		return math.NaN()
	case x == 0:
		return math.Inf(1)
	case math.IsInf(x, 1):
		return 0
	}
	nu = math.Abs(nu)
	if scaled {
		return special.Kve(nu, x)
	}
	return special.Kv(nu, x)
}

// jvHankel is J_nu(x) for large x by Hankel's asymptotic expansion (Abramowitz
// and Stegun 9.2.5, 9.2.9, 9.2.10): sqrt(2/(pi x)) (P cos chi - Q sin chi) with
// chi = x - (nu/2 + 1/4) pi. chi is reduced modulo 2 pi in 320-bit arithmetic,
// so the phase is exact for any finite x; special.Jv loses it beyond about
// 1e8. ok is false when the series does not converge to 1e-17 (nu^2 not small
// against x), and the caller falls back to special.Jv.
func jvHankel(nu, x float64) (r float64, ok bool) {
	mu := 4 * nu * nu
	t, p, q := 1.0, 1.0, 0.0
	for k := 1; ; k++ {
		next := t * (mu - float64((2*k-1)*(2*k-1))) / (float64(k) * 8 * x)
		if math.Abs(next) >= math.Abs(t) && k > 1 || k > 60 {
			return 0, false
		}
		t = next
		switch k & 3 {
		case 1:
			q += t
		case 2:
			p -= t
		case 3:
			q -= t
		case 0:
			p += t
		}
		if math.Abs(t) < 1e-17 {
			break
		}
	}
	const prec = 320
	pi := bigPi320()
	twoPi := new(big.Float).SetPrec(prec).Mul(pi, big.NewFloat(2))
	chi := new(big.Float).SetPrec(prec).Mul(pi, new(big.Float).SetPrec(prec).Quo(big.NewFloat(2*nu+1), big.NewFloat(4)))
	chi.Sub(new(big.Float).SetPrec(prec).SetFloat64(x), chi)
	q0 := new(big.Float).SetPrec(prec).Quo(chi, twoPi)
	n, _ := q0.Int(nil) // truncated toward zero
	chi.Sub(chi, new(big.Float).SetPrec(prec).Mul(new(big.Float).SetPrec(prec).SetInt(n), twoPi))
	c, _ := chi.Float64() // in (-2 pi, 2 pi): the cosine and sine of it are the phase's
	return math.Sqrt(2/(math.Pi*x)) * (p*math.Cos(c) - q*math.Sin(c)), true
}

var bigPi320Once = func() *big.Float {
	f, _, _ := big.ParseFloat("3.14159265358979323846264338327950288419716939937510582097494459230781640628620899862803482534211706798214808651", 10, 320, big.ToNearestEven)
	return f
}()

func bigPi320() *big.Float { return bigPi320Once }

// besselArgs checks the numeric arguments of a binary or ternary math
// function and recycles them: the doubles of each, and the common length.
func (s *Session) mathArgs(c *callInfo, a []Value) (d [][]float64, m int) {
	d = make([][]float64, len(a))
	for i, v := range a {
		switch v.(type) {
		case *Logical, *Integer, *Double:
		default:
			s.errorf(c.call, "non-numeric argument to mathematical function")
		}
		d[i] = asDoubles(v)
		if len(d[i]) == 0 {
			return d, 0
		}
		m = max(m, len(d[i]))
	}
	return d, m
}

// mathAttrs: the attributes come from the first argument of full length.
func mathAttrs(a []Value, d [][]float64, m int) []Attr {
	for i := range a {
		if len(d[i]) == m {
			return attrsOf(a[i])
		}
	}
	return nil
}

// needArgs reports the first missing required argument as R does.
func (s *Session) needArgs(c *callInfo, a []Value, names ...string) {
	for i, n := range names {
		if a[i] == nil {
			s.errorf(c.call, "argument \"%s\" is missing, with no default", n)
		}
	}
}

func init() {
	for _, e := range []struct {
		name string
		f    func(a, b float64) float64
	}{{"beta", betaFn}, {"lbeta", lbetaFn}} {
		regBuiltin(e.name, func(s *Session, c *callInfo) Value {
			a := s.matchFormals(c, "a", "b")
			s.needArgs(c, a, "a", "b")
			d, m := s.mathArgs(c, a)
			if m == 0 && e.name == "lbeta" { // lbeta is a closure in R: no attributes
				return &Double{V: []float64{}}
			}
			out := make([]float64, m)
			nan := false
			for i := range out {
				x, y := d[0][i%len(d[0])], d[1][i%len(d[1])]
				out[i] = e.f(x, y)
				if math.IsNaN(out[i]) && !math.IsNaN(x) && !math.IsNaN(y) {
					nan = true
				}
			}
			if nan {
				s.warningf(c.call, "NaNs produced")
			}
			return &Double{V: out, Attr: mathAttrs(a, d, m)}
		})
	}

	regBuiltin("psigamma", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "deriv")
		s.needArgs(c, a, "x")
		if a[1] == nil {
			a[1] = &Double{V: []float64{0}}
		}
		d, m := s.mathArgs(c, a)
		out := make([]float64, m)
		nan := false
		for i := range out {
			x, dv := d[0][i%len(d[0])], d[1][i%len(d[1])]
			if math.IsNaN(x) || math.IsNaN(dv) {
				out[i] = x + dv
				continue
			}
			n := math.RoundToEven(dv)
			switch {
			case n > 100:
				k := math.Min(n, math.MaxInt32) // C's int conversion saturates in the text
				s.warningf(c.call, "deriv = %d > %d (= n_max)", int64(k), 100)
				out[i] = math.NaN()
			case n < 0:
				out[i] = math.NaN()
			default:
				out[i] = psigammaN(x, int(n))
			}
			if math.IsNaN(out[i]) {
				nan = true
			}
		}
		if nan {
			s.warningf(c.call, "NaNs produced")
		}
		return &Double{V: out, Attr: mathAttrs(a, d, m)}
	})

	for _, name := range []string{"besselJ", "besselY", "besselI", "besselK"} {
		scalable := name == "besselI" || name == "besselK"
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			var a []Value
			if scalable {
				a = s.matchFormals(c, "x", "nu", "expon.scaled")
			} else {
				a = s.matchFormals(c, "x", "nu")
			}
			s.needArgs(c, a, "x", "nu")
			if scalable && a[2] == nil {
				a[2] = &Logical{V: []int32{0}}
			}
			if sv, ok := a[len(a)-1].(*String); ok && scalable {
				lv := make([]int32, len(sv.V)) // as.logical("T") is TRUE, other text NA
				for i, e := range sv.V {
					switch {
					case e.NA:
						lv[i] = naInt
					case e.S == "T" || e.S == "TRUE" || e.S == "true" || e.S == "True":
						lv[i] = 1
					case e.S == "F" || e.S == "FALSE" || e.S == "false" || e.S == "False":
					default:
						lv[i] = naInt
					}
				}
				a[2] = &Logical{V: lv}
			}
			d, m := s.mathArgs(c, a)
			w := func(format string, args ...any) { s.warningf(c.call, format, args...) }
			out := make([]float64, m)
			nan := false
			for i := range out {
				if i&255 == 0 {
					s.checkInterrupt()
				}
				x, nu := d[0][i%len(d[0])], d[1][i%len(d[1])]
				scaled, expNaN := false, false
				if scalable {
					e := d[2][i%len(d[2])]
					scaled, expNaN = e != 0, math.IsNaN(e)
				}
				switch {
				case expNaN:
					out[i] = naReal
				case name == "besselJ":
					out[i] = besselJ(x, nu, w)
				case name == "besselY":
					out[i] = besselY(x, nu, w)
				case name == "besselI":
					out[i] = besselI(x, nu, scaled, w)
				default:
					out[i] = besselK(x, nu, scaled, w)
				}
				if math.IsNaN(out[i]) && !math.IsNaN(x) && !math.IsNaN(nu) && !expNaN {
					nan = true
				}
			}
			if nan {
				s.warningf(c.call, "NaNs produced")
			}
			return &Double{V: out, Attr: mathAttrs(a, d, m)}
		})
	}

	regBuiltin("exp", (*Session)(nil).math1("exp", rexp, false))
}

// rexp is exp rounded through double-double: correctly rounded except where
// the exact value lies within 2^-76 (relative) of a rounding boundary.
func rexp(x float64) float64 {
	switch {
	case math.IsNaN(x):
		return x
	case math.IsInf(x, 1):
		return x
	case math.IsInf(x, -1):
		return 0
	}
	return ddp.Exp(ddp.F(x))
}

// rpowFinite is x^y for finite x != 0 and finite y (callers handle the rest):
// a negative base needs an integer exponent. exp(y log x) is formed in
// double-double, so the result is correctly rounded but for the same
// boundary cases as rexp.
func rpowFinite(x, y float64) float64 {
	neg := false
	if x < 0 {
		if y != math.Floor(y) {
			return rNaN
		}
		x = -x
		neg = math.Mod(y, 2) != 0
	}
	var r float64
	if est := math.Pow(x, y); y == math.Floor(y) && math.Abs(y) <= 1024 && est > 1e-290 && est < 1e290 {
		// An integer power can be an exact tie between two doubles (41^10 / 2^40):
		// binary powering in double-double keeps it exact where it fits, and the
		// rounding of Hi + Lo then breaks the tie to even, as a correctly rounded
		// pow does.
		p, base := ddp.F(1), ddp.F(x)
		for n := int(math.Abs(y)); n > 0; n >>= 1 {
			if n&1 == 1 {
				p = ddp.Mul(p, base)
			}
			base = ddp.Sqr(base)
		}
		if y < 0 {
			p = ddp.Div(ddp.F(1), p)
		}
		r = p.Hi + p.Lo
	} else {
		r = ddp.Exp(ddp.MulF(ddp.Log(x), y))
	}
	if neg {
		return -r
	}
	return r
}
