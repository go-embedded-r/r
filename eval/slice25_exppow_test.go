package eval

import (
	"math"
	"math/big"
	"math/rand"
	"testing"
)

// The stated bound of exp and ^ (slice 25): the result is the correctly rounded
// value, except where the exact value lies within 2^-76 (relative) of a rounding
// boundary. The oracle is a 320-bit exp and log built here from a Taylor series
// and Newton's iteration, independent of the dd package the kernels use.

const expPrec = 320

func bigExp(x *big.Float) *big.Float {
	const k = 24
	r := new(big.Float).SetPrec(expPrec).SetMantExp(x, -k)
	sum := new(big.Float).SetPrec(expPrec).SetInt64(1)
	term := new(big.Float).SetPrec(expPrec).SetInt64(1)
	for i := 1; i < 200; i++ {
		term.Mul(term, r)
		term.Quo(term, new(big.Float).SetPrec(expPrec).SetInt64(int64(i)))
		sum.Add(sum, term)
		if term.Sign() == 0 || term.MantExp(nil) < -expPrec-8 {
			break
		}
	}
	for i := 0; i < k; i++ {
		sum.Mul(sum, sum)
	}
	return sum
}

func bigLog(x float64) *big.Float {
	bx := new(big.Float).SetPrec(expPrec).SetFloat64(x)
	l := new(big.Float).SetPrec(expPrec).SetFloat64(math.Log(x))
	two := new(big.Float).SetPrec(expPrec).SetInt64(2)
	for i := 0; i < 4; i++ { // quadratic... cubic convergence of Halley's step
		e := bigExp(l)
		num := new(big.Float).SetPrec(expPrec).Sub(bx, e)
		den := new(big.Float).SetPrec(expPrec).Add(bx, e)
		l.Add(l, new(big.Float).SetPrec(expPrec).Mul(two, num.Quo(num, den)))
	}
	return l
}

func TestExpPowBounds(t *testing.T) {
	rng := rand.New(rand.NewSource(30))
	e := &boundCase{name: "exp", bound: 0}
	p := &boundCase{name: "^ (non-integer and integer exponents)", bound: 0}
	for i := 0; i < 400; i++ {
		x := (rng.Float64()*2 - 1) * 700
		want, _ := bigExp(new(big.Float).SetPrec(expPrec).SetFloat64(x)).Float64()
		e.see(t, rexp(x), want, fmtAt(x))
	}
	for i := 0; i < 400; i++ {
		x := math.Floor(rng.Float64()*4096) / 64 // few-bit bases, exact in big
		y := (rng.Float64()*2 - 1) * 40
		if i%3 == 0 {
			y = math.Floor(y)
		}
		if x == 0 {
			continue
		}
		prod := new(big.Float).SetPrec(expPrec).Mul(bigLog(x), new(big.Float).SetPrec(expPrec).SetFloat64(y))
		want, _ := bigExp(prod).Float64()
		p.see(t, rpowFinite(x, y), want, fmtAt(x, y))
	}
	// an integer power of a few-bit base is exact in big arithmetic
	for _, c := range [][2]float64{{2.5625, 10}, {41.0 / 16, 9}, {1.1, 7}, {3, 40}, {5, 30}, {7, 20}} {
		w := new(big.Float).SetPrec(2000).SetFloat64(c[0])
		want := new(big.Float).SetPrec(2000).SetInt64(1)
		for i := 0; i < int(c[1]); i++ {
			want.Mul(want, w)
		}
		wf, _ := want.Float64()
		p.see(t, rpowFinite(c[0], c[1]), wf, fmtAt(c[0], c[1]))
	}
	e.done(t)
	p.done(t)
}
