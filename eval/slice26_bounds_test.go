package eval

import (
	"math"
	"math/big"
	"testing"
)

// The accuracy of d, p and q of the twelve families of slice 26, asserted
// against oracles built here from series and closed forms at 320 bits, on the
// four scales (lower and upper tail, plain and log). The oracles share nothing
// with go-statistics: the normal law is the Taylor series of erf, the gamma and
// beta laws the series of the incomplete gamma and beta functions with the
// log-gamma of the double-double package of the slice 25 oracles, t and F go
// through the beta law, the logistic and Cauchy laws are closed forms.
//
// The stated bound of each family is in the README; a kernel that moves past
// one fails here. The control at the end proves the oracles can report an error.

func nb() *big.Float { return new(big.Float).SetPrec(expPrec) }

func bAdd(a, b *big.Float) *big.Float { return nb().Add(a, b) }
func bSub(a, b *big.Float) *big.Float { return nb().Sub(a, b) }
func bMul(a, b *big.Float) *big.Float { return nb().Mul(a, b) }
func bQuo(a, b *big.Float) *big.Float { return nb().Quo(a, b) }
func bNeg(a *big.Float) *big.Float    { return nb().Neg(a) }
func bInt(n int) *big.Float           { return nb().SetInt64(int64(n)) }
func bFlt(x float64) *big.Float       { return nb().SetFloat64(x) }
func bOne() *big.Float                { return bInt(1) }

// bLn is the natural log of a positive big float, by Halley's iteration on bigExp
// from the float64 log.
func bLn(x *big.Float) *big.Float {
	f, _ := x.Float64()
	l := bFlt(math.Log(f))
	two := bInt(2)
	for i := 0; i < 4; i++ {
		e := bigExp(l)
		l = bAdd(l, bMul(two, bQuo(bSub(x, e), bAdd(x, e))))
	}
	return l
}

// bPow is x^y for x > 0.
func bPow(x, y *big.Float) *big.Float { return bigExp(bMul(y, bLn(x))) }

func bPi() *big.Float { return nb().Set(bigPi) }

// bPhi is the standard normal distribution function, |x| < 8.
func bPhi(x *big.Float) *big.Float {
	z := bQuo(x, nb().Sqrt(bInt(2)))
	z2 := bMul(z, z)
	term := nb().Set(z)
	sum := nb().Set(z)
	for n := 1; n < 200; n++ {
		term = bQuo(bMul(term, z2), bInt(-n))
		sum = bAdd(sum, bQuo(term, bInt(2*n+1)))
	}
	erf := bMul(bInt(2), bQuo(sum, nb().Sqrt(bPi())))
	return bQuo(bAdd(bOne(), erf), bInt(2))
}

// bAtan is atan by the halving identity and the Taylor series.
func bAtan(x *big.Float) *big.Float {
	if x.Sign() < 0 {
		return bNeg(bAtan(bNeg(x)))
	}
	h := 0
	for x.Cmp(bFlt(0.25)) > 0 {
		x = bQuo(x, bAdd(bOne(), nb().Sqrt(bAdd(bOne(), bMul(x, x)))))
		h++
	}
	x2 := bMul(x, x)
	term := nb().Set(x)
	sum := nb().Set(x)
	for n := 1; n < 200; n++ {
		term = bNeg(bMul(term, x2))
		sum = bAdd(sum, bQuo(term, bInt(2*n+1)))
	}
	for ; h > 0; h-- {
		sum = bMul(sum, bInt(2))
	}
	return sum
}

func ddBig(a float64, b float64) *big.Float { return bAdd(bFlt(a), bFlt(b)) }

// bGammaP is the regularised lower incomplete gamma function P(a, x), x > 0.
func bGammaP(a, x *big.Float) *big.Float {
	af, _ := a.Float64()
	lg := ddLogGamma(af + 1)
	pre := bigExp(bSub(bSub(bMul(a, bLn(x)), x), ddBig(lg.Hi, lg.Lo)))
	term := bOne()
	sum := bOne()
	for n := 1; n < 20000; n++ {
		term = bQuo(bMul(term, x), bAdd(a, bInt(n)))
		sum = bAdd(sum, term)
		if term.MantExp(nil) < sum.MantExp(nil)-expPrec-4 {
			break
		}
	}
	return bMul(pre, sum)
}

// bBetaI is the regularised incomplete beta function I_x(a, b) for 0 < x < 1,
// by the power series, on the side of the symmetry where it converges quickly;
// it returns the lower and the upper tail, each computed directly where it is
// the small one.
func bBetaI(x, a, b *big.Float) (lower, upper *big.Float) {
	series := func(x, a, b *big.Float) *big.Float {
		af, _ := a.Float64()
		bf, _ := b.Float64()
		lb := ddLogBetaDD(af, bf)
		l1x := bLn(bSub(bOne(), x))
		pre := bigExp(bSub(bAdd(bMul(a, bLn(x)), bMul(b, l1x)), bAdd(ddBig(lb.Hi, lb.Lo), bLn(a))))
		term := bOne()
		sum := bOne()
		ab := bAdd(a, b)
		for n := 0; n < 40000; n++ {
			term = bMul(term, bMul(bQuo(bAdd(ab, bInt(n)), bAdd(a, bInt(n+1))), x))
			sum = bAdd(sum, term)
			if term.MantExp(nil) < sum.MantExp(nil)-expPrec-4 {
				break
			}
		}
		return bMul(pre, sum)
	}
	xf, _ := x.Float64()
	af, _ := a.Float64()
	bf, _ := b.Float64()
	if xf < (af+1)/(af+bf+2) {
		lower = series(x, a, b)
		return lower, bSub(bOne(), lower)
	}
	upper = series(bSub(bOne(), x), b, a)
	return bSub(bOne(), upper), upper
}

// an oracle gives the lower and the upper tail at x (big), and the density.
type oracle struct {
	name  string
	tails func(x float64) (lower, upper *big.Float)
}

func oracleOf(name string, p []float64) oracle {
	b := bFlt
	switch name {
	case "norm":
		return oracle{name, func(x float64) (*big.Float, *big.Float) {
			z := bQuo(bSub(b(x), b(p[0])), b(p[1]))
			return bPhi(z), bPhi(bNeg(z))
		}}
	case "lnorm":
		return oracle{name, func(x float64) (*big.Float, *big.Float) {
			z := bQuo(bSub(bLn(b(x)), b(p[0])), b(p[1]))
			return bPhi(z), bPhi(bNeg(z))
		}}
	case "exp":
		return oracle{name, func(x float64) (*big.Float, *big.Float) {
			sf := bigExp(bNeg(bMul(b(p[0]), b(x))))
			return bSub(bOne(), sf), sf
		}}
	case "weibull":
		return oracle{name, func(x float64) (*big.Float, *big.Float) {
			t := bPow(bQuo(b(x), b(p[1])), b(p[0]))
			sf := bigExp(bNeg(t))
			return bSub(bOne(), sf), sf
		}}
	case "logis":
		return oracle{name, func(x float64) (*big.Float, *big.Float) {
			z := bQuo(bSub(b(x), b(p[0])), b(p[1]))
			return bQuo(bOne(), bAdd(bOne(), bigExp(bNeg(z)))), bQuo(bOne(), bAdd(bOne(), bigExp(z)))
		}}
	case "cauchy":
		return oracle{name, func(x float64) (*big.Float, *big.Float) {
			z := bQuo(bSub(b(x), b(p[0])), b(p[1]))
			a := bQuo(bAtan(z), bPi())
			h := b(0.5)
			return bAdd(h, a), bSub(h, a)
		}}
	case "gamma":
		return oracle{name, func(x float64) (*big.Float, *big.Float) {
			l := bGammaP(b(p[0]), bQuo(b(x), b(p[1])))
			return l, bSub(bOne(), l)
		}}
	case "chisq":
		return oracle{name, func(x float64) (*big.Float, *big.Float) {
			l := bGammaP(b(p[0]/2), b(x/2))
			return l, bSub(bOne(), l)
		}}
	case "beta":
		return oracle{name, func(x float64) (*big.Float, *big.Float) {
			return bBetaI(b(x), b(p[0]), b(p[1]))
		}}
	case "t":
		return oracle{name, func(x float64) (*big.Float, *big.Float) {
			n := b(p[0])
			ax := b(math.Abs(x))
			u := bQuo(n, bAdd(n, bMul(ax, ax)))
			l, _ := bBetaI(u, bQuo(n, bInt(2)), b(0.5))
			half := bQuo(l, bInt(2)) // the tail beyond |x|
			if x >= 0 {
				return bSub(bOne(), half), half
			}
			return half, bSub(bOne(), half)
		}}
	case "f":
		return oracle{name, func(x float64) (*big.Float, *big.Float) {
			d1, d2 := b(p[0]), b(p[1])
			t := bMul(d1, b(x))
			w := bQuo(t, bAdd(t, d2))
			return bBetaI(w, bQuo(d1, bInt(2)), bQuo(d2, bInt(2)))
		}}
	}
	panic("no oracle for " + name)
}

func familyOf(name string) *family {
	for i := range families {
		if families[i].name == name {
			return &families[i]
		}
	}
	panic("no family " + name)
}

func toF(x *big.Float) float64 { f, _ := x.Float64(); return f }

// boundPoint is one x with the parameters of a family.
type boundPoint struct {
	p  []float64
	xs []float64
}

// bigQuantile inverts a tail by bisection around the answer g: the p of the
// lower tail (or the q of the upper one), to 100 halvings of a relative window.
func s26Quantile(o oracle, lowerTail bool, p, g float64) float64 {
	pb := bFlt(p)
	w := 1e-6 * math.Max(math.Abs(g), 1e-300)
	lo, hi := bFlt(g-w), bFlt(g+w)
	half := bFlt(0.5)
	for i := 0; i < 100; i++ {
		mid := bMul(bAdd(lo, hi), half)
		l, u := o.tailsBig(mid)
		var cmp int
		if lowerTail {
			cmp = l.Cmp(pb)
		} else {
			cmp = pb.Cmp(u) // the upper tail falls as x grows
		}
		if cmp < 0 {
			lo = mid
		} else {
			hi = mid
		}
	}
	return toF(bMul(bAdd(lo, hi), half))
}

// tailsBig evaluates the tails at a big x (the grid of the oracles is float64,
// so the mid point is rounded to float64 first: the bisection stops improving at
// one ulp of x, which is the resolution of the answer).
func (o oracle) tailsBig(x *big.Float) (*big.Float, *big.Float) {
	f, _ := x.Float64()
	return o.tails(f)
}

var boundPoints = map[string][]boundPoint{
	"norm":    {{[]float64{0, 1}, []float64{-7.5, -5, -2.5, -1, -0.25, 0.5, 1.5, 3, 6, 7.9}}, {[]float64{3, 0.5}, []float64{1, 2.5, 4, 5.5}}},
	"lnorm":   {{[]float64{0, 1}, []float64{1e-3, 0.05, 0.5, 1, 3, 20, 500}}, {[]float64{-2, 0.5}, []float64{0.01, 0.2, 1}}},
	"exp":     {{[]float64{1}, []float64{1e-5, 0.1, 1, 5, 30, 120}}, {[]float64{0.01}, []float64{3, 70, 900}}},
	"weibull": {{[]float64{2, 3}, []float64{1e-3, 0.1, 1, 3, 7, 20}}, {[]float64{0.5, 1}, []float64{1e-6, 0.01, 1, 20, 100}}, {[]float64{7.5, 2}, []float64{0.5, 1.5, 2, 2.5}}},
	"logis":   {{[]float64{0, 1}, []float64{-40, -10, -3, -0.1, 0.5, 2, 15, 40}}, {[]float64{1, 0.3}, []float64{-2, 0.5, 1.2, 3}}},
	"cauchy":  {{[]float64{0, 1}, []float64{-1e6, -50, -3, -0.2, 0.4, 2, 30, 1e6}}, {[]float64{2, 0.1}, []float64{1.5, 1.95, 2.05, 3}}},
	"gamma":   {{[]float64{2.5, 1}, []float64{0.05, 0.5, 2, 5, 12, 30}}, {[]float64{0.5, 2}, []float64{1e-4, 0.1, 1, 8, 30}}, {[]float64{30, 0.5}, []float64{8, 14, 15, 20, 30}}},
	"chisq":   {{[]float64{3}, []float64{0.05, 0.5, 2, 6, 15, 40}}, {[]float64{40}, []float64{15, 30, 40, 60, 90}}, {[]float64{0.5}, []float64{1e-4, 0.1, 3, 20}}},
	"beta":    {{[]float64{2, 3}, []float64{1e-3, 0.1, 0.4, 0.7, 0.95, 0.999}}, {[]float64{0.5, 0.5}, []float64{0.01, 0.3, 0.5, 0.9}}, {[]float64{20, 7}, []float64{0.4, 0.7, 0.9}}},
	"t":       {{[]float64{5}, []float64{-30, -4, -1, -0.1, 0.5, 2, 8, 60}}, {[]float64{1.5}, []float64{-100, -2, 0.7, 12}}, {[]float64{40}, []float64{-6, -1, 2, 5}}},
	"f":       {{[]float64{3, 10}, []float64{0.01, 0.4, 1, 3, 12, 60}}, {[]float64{20, 2.5}, []float64{0.2, 1, 5, 40}}, {[]float64{1, 1}, []float64{0.01, 1, 90}}},
}

// stated bounds in ulps of the exact value, per scale (README, slice 26)
var boundTable = map[string]struct{ plain, logp, q float64 }{
	"norm":    {4, 8, 8},
	"lnorm":   {64, 64, 16}, // conditioning: the rounding of log(x) is amplified by z^2
	"exp":     {4, 4, 8},
	"weibull": {64, 8, 16}, // conditioning: exp(-t) amplifies the rounding of t

	"logis":  {4, 8, 8},
	"cauchy": {8, 8, 16},
	"gamma":  {1024, 1024, 64}, // go-statistics' upper tail: 565 ulp at pgamma(30, 2.5, lower.tail = FALSE)
	"chisq":  {16, 16, 16},
	"beta":   {16, 16, 16},
	"t":      {16, 16, 16},
	"f":      {16, 16, 16},
}

func TestSlice26Bounds(t *testing.T) {
	for _, name := range []string{"norm", "lnorm", "exp", "weibull", "logis", "cauchy", "gamma", "chisq", "beta", "t", "f"} {
		name := name
		t.Run(name, func(t *testing.T) {
			if testing.Short() {
				t.Skip("series oracles at 320 bits")
			}
			f := familyOf(name)
			bt := boundTable[name]
			pl := &boundCase{name: name + " p/d plain", bound: bt.plain}
			lg := &boundCase{name: name + " log scale", bound: bt.logp}
			qq := &boundCase{name: name + " quantile", bound: bt.q}
			for _, pt := range boundPoints[name] {
				o := oracleOf(name, pt.p)
				row := append([]float64(nil), pt.p...)
				for len(row) < len(f.params) {
					row = append(row, 0)
				}
				if name == "gamma" {
					row = []float64{pt.p[0], pt.p[1]}
				}
				for _, x := range pt.xs {
					lo, up := o.tails(x)
					wl, wu := toF(lo), toF(up)
					at := fmtAt(append([]float64{x}, pt.p...)...)
					pl.see(t, f.value(kindP, x, row, true, false), wl, at)
					pl.see(t, f.value(kindP, x, row, false, false), wu, at)
					if wl > 0 && wl < 1 {
						lg.see(t, f.value(kindP, x, row, true, true), toF(bLn(lo)), at)
					}
					if wu > 0 && wu < 1 {
						lg.see(t, f.value(kindP, x, row, false, true), toF(bLn(up)), at)
					}
					// the quantile is the inverse: bisect the oracle around our answer
					for _, lower := range []bool{true, false} {
						pr := wl
						if !lower {
							pr = wu
						}
						if pr <= 1e-12 || pr >= 1-1e-9 {
							continue
						}
						g := f.value(kindQ, pr, row, lower, false)
						want := s26Quantile(o, lower, pr, g)
						// the conditioning of the inverse: the answer cannot be more
						// accurate than the probability's own rounding allows
						qq.see(t, g, want, at+" p="+fmtAt(pr))
					}
				}
			}
			for _, c := range []*boundCase{pl, lg, qq} {
				c.done(t)
			}
		})
	}
}

// TestSlice26OracleControl: the oracles are not vacuous. Known values, and an
// error that must be seen.
func TestSlice26OracleControl(t *testing.T) {
	for _, c := range []struct {
		name      string
		got, want float64
	}{
		{"Phi(1)", toF(bPhi(bFlt(1))), 0.8413447460685429},
		{"Phi(-3)", toF(bPhi(bFlt(-3))), 0.0013498980316300946},
		{"atan(1)", toF(bAtan(bFlt(1))), math.Pi / 4},
		{"atan(-1e6)", toF(bAtan(bFlt(-1e6))), -1.5707953267948966},
		{"P(1, 1)", toF(bGammaP(bFlt(1), bFlt(1))), 1 - math.Exp(-1)},
		{"P(0.5, 2)", toF(bGammaP(bFlt(0.5), bFlt(2))), math.Erf(math.Sqrt(2))},
		{"I(1/2; 2, 3)", toF(bBetaI2(bFlt(0.5), bFlt(2), bFlt(3))), 0.6875},
		{"ln 10", toF(bLn(bFlt(10))), math.Ln10},
		{"2^0.5", toF(bPow(bFlt(2), bFlt(0.5))), math.Sqrt2},
	} {
		if e := ulpErr(c.got, c.want); e > 2 {
			t.Errorf("oracle %s = %.17g, want %.17g (%g ulp)", c.name, c.got, c.want, e)
		}
	}
	// the other tail of the beta oracle: 1 - I_0.9(2, 3) = I_0.1(3, 2) = 0.0037
	if _, u := bBetaI(bFlt(0.9), bFlt(2), bFlt(3)); math.Abs(toF(u)/0.0037-1) > 1e-13 { // 0.9 is not exactly 9/10
		t.Errorf("upper tail of I_0.9(2, 3) = %.17g", toF(u))
	}
	// the bound check registers an error
	b := &boundCase{name: "control"}
	b.see(t, math.Nextafter(0.25, 1), 0.25, "0.25")
	if b.worst < 1 {
		t.Error("a one-ulp error was not seen")
	}
	// and the bisection finds a known quantile: the median of the logistic law
	o := oracleOf("logis", []float64{0, 1})
	if q := s26Quantile(o, true, 0.25, -1.098612); ulpErr(q, -math.Log(3)) > 4 {
		t.Errorf("lower quartile of the logistic law = %.17g, want -ln 3", q)
	}
	if q := s26Quantile(o, false, 0.25, 1.098612); ulpErr(q, math.Log(3)) > 4 {
		t.Errorf("upper quartile of the logistic law = %.17g, want ln 3", q)
	}
}

func bBetaI2(x, a, b *big.Float) *big.Float { l, _ := bBetaI(x, a, b); return l }
