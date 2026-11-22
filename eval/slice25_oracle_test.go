package eval

// The oracles of slice 25: double-double (go-numerics/special/dd, ~2^-100)
// for gamma, lgamma, digamma, polygamma and beta, and math/big power series
// (600 bits) for the Bessel functions I and J. Each oracle is independent of
// the code under test.

import (
	"math"
	"math/big"

	ddp "github.com/go-numerics/special/dd"
)

func ddFrac(num, den float64) ddp.DD { return ddp.Div(ddp.F(num), ddp.F(den)) }

var ddLn2Pi = func() ddp.DD {
	v, err := ddp.Parse("1.83787706640934548356065947281123527972279494727556682563430308096553139185452079538948659727190839524401129324926867")
	if err != nil {
		panic(err)
	}
	return v
}()

// stirlingShift returns z = x + n >= 40 (exact when x has few bits) and n.
func stirlingShift(x float64) (z float64, n int) {
	for z = x; z < 40; z++ {
		n++
	}
	return z, n
}

// B2k/(2k(2k-1)), k = 1..10: the coefficients of the Stirling series of ln Gamma.
var stirlingLG = []ddp.DD{
	ddFrac(1, 12), ddFrac(-1, 360), ddFrac(1, 1260), ddFrac(-1, 1680), ddFrac(1, 1188),
	ddFrac(-691, 360360), ddFrac(1, 156), ddFrac(-3617, 122400), ddFrac(43867, 244188), ddFrac(-174611, 125400),
}

// B2k/(2k): the coefficients of the series of digamma.
var stirlingPsi = []ddp.DD{
	ddFrac(1, 12), ddFrac(-1, 120), ddFrac(1, 252), ddFrac(-1, 240), ddFrac(1, 132),
	ddFrac(-691, 32760), ddFrac(1, 12), ddFrac(-3617, 8160), ddFrac(43867, 14364), ddFrac(-174611, 6600),
}

// B2k.
var bernoulli2k = []ddp.DD{
	ddFrac(1, 6), ddFrac(-1, 30), ddFrac(1, 42), ddFrac(-1, 30), ddFrac(5, 66),
	ddFrac(-691, 2730), ddFrac(7, 6), ddFrac(-3617, 510), ddFrac(43867, 798), ddFrac(-174611, 330),
}

func ddPowInt(z float64, k int) ddp.DD {
	r := ddp.F(1)
	for i := 0; i < k; i++ {
		r = ddp.MulF(r, z)
	}
	return r
}

// ddLogGamma is ln Gamma(x), x > 0 with few bits.
func ddLogGamma(x float64) ddp.DD {
	z, n := stirlingShift(x)
	r := ddp.Sub(ddp.MulF(ddp.Log(z), z-0.5), ddp.F(z))
	r = ddp.Add(r, ddp.MulF(ddLn2Pi, 0.5))
	zp := ddp.F(z)
	zz := ddp.Mul(ddp.F(z), ddp.F(z))
	inv := ddp.Div(ddp.F(1), zp)
	for _, c := range stirlingLG {
		r = ddp.Add(r, ddp.Mul(c, inv))
		inv = ddp.Div(inv, zz)
	}
	for i := 0; i < n; i++ {
		r = ddp.Sub(r, ddp.Log(x+float64(i)))
	}
	return r
}

func ddGamma(x float64) float64 { return ddp.Exp(ddLogGamma(x)) }

func ddDigamma(x float64) float64 {
	z, n := stirlingShift(x)
	r := ddp.Sub(ddp.Log(z), ddp.Div(ddp.F(0.5), ddp.F(z)))
	zz := ddp.Mul(ddp.F(z), ddp.F(z))
	inv := ddp.Div(ddp.F(1), zz)
	for _, c := range stirlingPsi {
		r = ddp.Sub(r, ddp.Mul(c, inv))
		inv = ddp.Div(inv, zz)
	}
	for i := 0; i < n; i++ {
		r = ddp.Sub(r, ddp.Div(ddp.F(1), ddp.F(x+float64(i))))
	}
	return r.Hi + r.Lo
}

// ddPolygamma is psi^(m)(x), x > 0, 1 <= m <= 6.
func ddPolygamma(m int, x float64) float64 {
	z, n := stirlingShift(x)
	fact := func(k int) float64 {
		r := 1.0
		for i := 2; i <= k; i++ {
			r *= float64(i)
		}
		return r
	}
	sign := 1.0
	if m%2 == 0 {
		sign = -1
	}
	// psi^(m)(z) = (-1)^(m+1) [ (m-1)!/z^m + m!/(2 z^(m+1)) + sum B2k (2k+m-1)!/((2k)! z^(2k+m)) ]
	t := ddp.Div(ddp.F(fact(m-1)), ddPowInt(z, m))
	t = ddp.Add(t, ddp.Div(ddp.F(fact(m)/2), ddPowInt(z, m+1)))
	for k := 1; k <= 10; k++ {
		ratio := 1.0
		for j := 2*k + 1; j <= 2*k+m-1; j++ {
			ratio *= float64(j)
		}
		t = ddp.Add(t, ddp.Div(ddp.MulF(bernoulli2k[k-1], ratio), ddPowInt(z, 2*k+m)))
	}
	r := ddp.MulF(t, sign)
	for i := 0; i < n; i++ {
		// psi^(m)(x) = psi^(m)(x+1) - (-1)^m m!/x^(m+1)
		r = ddp.Add(r, ddp.MulF(ddp.Div(ddp.F(fact(m)), ddPowInt(x+float64(i), m+1)), sign))
	}
	return r.Hi + r.Lo
}

func ddLogBetaDD(a, b float64) ddp.DD {
	return ddp.Sub(ddp.Add(ddLogGamma(a), ddLogGamma(b)), ddLogGamma(a+b))
}

func ddLogBeta(a, b float64) float64 { r := ddLogBetaDD(a, b); return r.Hi + r.Lo }

func ddBeta(a, b float64) float64 { return ddp.Exp(ddLogBetaDD(a, b)) }

// ulpErr is |got-want| in units of the last place of want.
func ulpErr(got, want float64) float64 {
	if got == want {
		return 0
	}
	if math.IsNaN(got) || math.IsNaN(want) || math.IsInf(got, 0) || math.IsInf(want, 0) {
		return math.Inf(1)
	}
	w := math.Abs(want)
	u := math.Nextafter(w, math.Inf(1)) - w
	if u == 0 {
		u = 5e-324
	}
	return math.Abs(got-want) / u
}

const bigPrec = 600

func bigF(x float64) *big.Float { return new(big.Float).SetPrec(bigPrec).SetFloat64(x) }

// bigPi: pi to 600 bits, from its decimal expansion.
var bigPi = func() *big.Float {
	f, _, err := big.ParseFloat("3.14159265358979323846264338327950288419716939937510582097494459230781640628620899862803482534211706798214808651328230664709384460955058223172535940812848111745028410270193852110555964462294895493038196", 10, bigPrec, big.ToNearestEven)
	if err != nil {
		panic(err)
	}
	return f
}()

// bigSeries sums (sign)^k (x/2)^(2k+nu) / (k! Gamma(k+nu+1)) for the
// Bessel functions J (sign -1) and I (sign +1), nu a non-negative integer or
// half-integer, x <= 60. Gamma(nu+1) is exact: n!, or the half-integer product
// times sqrt(pi).
func bigSeries(nu, x float64, sign int) float64 {
	h := new(big.Float).SetPrec(bigPrec).Quo(bigF(x), bigF(2)) // x/2
	h2 := new(big.Float).SetPrec(bigPrec).Mul(h, h)
	// term0 = (x/2)^nu / Gamma(nu+1)
	integer := nu == math.Floor(nu)
	term := bigF(1)
	var g *big.Float // Gamma(nu+1)
	if integer {
		g = bigF(1)
		for i := 2; i <= int(nu); i++ {
			g.Mul(g, bigF(float64(i)))
		}
		for i := 0; i < int(nu); i++ {
			term.Mul(term, h)
		}
	} else {
		sq := new(big.Float).SetPrec(bigPrec).Sqrt(bigPi)
		g = new(big.Float).SetPrec(bigPrec).Set(sq) // Gamma(1/2)
		for k := 0.5; k < nu; k++ {
			g.Mul(g, bigF(k)) // Gamma(nu+1) = Gamma(1/2) * (1/2)(3/2)...(nu)
		}
		g.Mul(g, bigF(nu)) // the last factor: Gamma(nu+1) = nu Gamma(nu)
		// the loop above multiplied by 1/2, 3/2, ..., nu-1: the product of Gamma(nu) is complete
		hs := new(big.Float).SetPrec(bigPrec).Sqrt(h)
		for i := 0; i < int(nu-0.5); i++ {
			term.Mul(term, h)
		}
		term.Mul(term, hs) // h^(nu): h^(integer part) * sqrt(h)
		_ = hs
	}
	term.Quo(term, g)
	sum := new(big.Float).SetPrec(bigPrec).Set(term)
	for k := 1; k < 400; k++ {
		term.Mul(term, h2)
		term.Quo(term, bigF(float64(k)))
		term.Quo(term, bigF(float64(k)+nu))
		if sign < 0 {
			term.Neg(term)
		}
		sum.Add(sum, term)
	}
	f, _ := sum.Float64()
	return f
}
