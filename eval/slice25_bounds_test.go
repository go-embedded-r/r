package eval

import (
	"math"
	"math/big"
	"math/rand"
	"testing"

	"github.com/go-numerics/special"
)

// The accuracy of the kernels slice 25 binds, measured against the oracles of
// slice25_oracle_test.go at sampled points. The bounds are the ones this slice
// states in the README; a kernel that moves past one fails here.

type boundCase struct {
	name  string
	bound float64 // in ulps of the exact value
	worst float64
	at    string
}

func (b *boundCase) see(t *testing.T, got, want float64, at string) {
	if e := ulpErr(got, want); e > b.worst {
		b.worst, b.at = e, at
	}
}

func (b *boundCase) done(t *testing.T) {
	t.Helper()
	t.Logf("%s: worst %.2f ulp at %s (bound %.0f)", b.name, b.worst, b.at, b.bound)
	if b.worst > b.bound {
		t.Errorf("%s: %.2f ulp at %s exceeds the bound of %.0f", b.name, b.worst, b.at, b.bound)
	}
}

func grid(rng *rand.Rand, lo, hi float64, step float64, n int) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = lo + step*math.Floor(rng.Float64()*(hi-lo)/step)
	}
	return out
}

func fmtAt(xs ...float64) string {
	s := ""
	for _, x := range xs {
		s += " " + formatG(x, 17, false, false)
	}
	return s
}

func TestSpecialBoundsGamma(t *testing.T) {
	rng := rand.New(rand.NewSource(25))
	g := &boundCase{name: "gamma", bound: 2}
	lg := &boundCase{name: "lgamma", bound: 2}
	dg := &boundCase{name: "digamma", bound: 2}
	tg := &boundCase{name: "trigamma", bound: 2}
	pg := &boundCase{name: "psigamma(x, 2..5)", bound: 2}
	for _, x := range grid(rng, 1.0/64, 170, 1.0/64, 600) {
		g.see(t, gammaFn(x), ddGamma(x), fmtAt(x))
		if math.Abs(x-1) > 0.05 && math.Abs(x-2) > 0.05 {
			lg.see(t, lgammaFn(x), ddLogGamma(x).Hi, fmtAt(x))
		}
	}
	for _, x := range grid(rng, 1.0/64, 200, 1.0/64, 600) {
		if math.Abs(x-1.4616321449683623) > 0.05 {
			dg.see(t, digammaFn(x), ddDigamma(x), fmtAt(x))
		}
		tg.see(t, trigammaFn(x), ddPolygamma(1, x), fmtAt(x))
		for m := 2; m <= 5; m++ {
			pg.see(t, psigammaN(x, m), ddPolygamma(m, x), fmtAt(x, float64(m)))
		}
	}
	for _, c := range []*boundCase{g, lg, dg, tg, pg} {
		c.done(t)
	}
}

// TestOracleControl: the oracles are not vacuous (the bound tests report 0 ulp
// when the kernel is right, so a control proves they can report more).
func TestOracleControl(t *testing.T) {
	if e := ulpErr(math.Nextafter(1.5, 2), 1.5); e != 1 {
		t.Errorf("ulpErr of one step = %g", e)
	}
	for _, c := range []struct {
		name      string
		got, want float64
	}{
		{"gamma(1/2)", ddGamma(0.5), math.Sqrt(math.Pi)},
		{"gamma(5)", ddGamma(5), 24},
		{"digamma(1)", ddDigamma(1), -0.5772156649015329},
		{"trigamma(1)", ddPolygamma(1, 1), math.Pi * math.Pi / 6},
		{"psigamma(1, 2)", ddPolygamma(2, 1), -2.4041138063191885},
		{"lgamma(10)", ddLogGamma(10).Hi, 12.801827480081469},
		{"beta(2,3)", ddBeta(2, 3), 1.0 / 12},
		{"J0(1)", bigSeries(0, 1, -1), 0.7651976865579666},
		{"I1(1)", bigSeries(1, 1, +1), 0.5651591039924851},
		{"J_1/2(1)", bigSeries(0.5, 1, -1), math.Sqrt(2/math.Pi) * math.Sin(1)},
	} {
		if e := ulpErr(c.got, c.want); e > 2 {
			t.Errorf("oracle %s = %.17g, want %.17g (%g ulp)", c.name, c.got, c.want, e)
		}
	}
	// the code under test must register when it is wrong
	b := &boundCase{name: "control"}
	b.see(t, math.Nextafter(gammaFn(3.5), 10), ddGamma(3.5), "3.5")
	if b.worst < 1 {
		t.Error("a one-ulp error was not seen")
	}
}

// The outliers of the accuracy table: the region of the root of digamma
// (1.4616...), where the value is small against the terms that make it, and
// lbeta near its zero at a = b = 1.
func TestSpecialBoundsOutliers(t *testing.T) {
	dg := &boundCase{name: "digamma within 0.05 of its root", bound: 64}
	for i := -51; i <= 51; i++ {
		x := math.Floor(1.4616321449683623*1024)/1024 + float64(i)/1024
		dg.see(t, digammaFn(x), ddDigamma(x), fmtAt(x))
	}
	dg.done(t)
	dt := &boundCase{name: "digamma, tiny and huge arguments", bound: 64}
	for k := 5; k <= 40; k++ {
		x := math.Ldexp(1, -k)
		dt.see(t, digammaFn(x), ddDigamma(x), fmtAt(x))
	}
	for _, x := range []float64{1 << 10, 1 << 20, 1 << 30, 1 << 40, 1 << 50, 1.5 * (1 << 20), 12345.5} {
		dt.see(t, digammaFn(x), ddDigamma(x), fmtAt(x))
	}
	dt.done(t)
	lb := &boundCase{name: "lbeta with |value| < 0.5", bound: 64}
	rng := rand.New(rand.NewSource(28))
	for _, a := range grid(rng, 0.5, 3, 1.0/64, 300) {
		for _, b := range grid(rng, 0.5, 3, 1.0/64, 3) {
			if w := ddLogBeta(a, b); w != 0 && math.Abs(w) < 0.5 {
				lb.see(t, lbetaFn(a, b), w, fmtAt(a, b))
			}
		}
	}
	lb.done(t)
	lg := &boundCase{name: "lgamma within 0.05 of its roots 1 and 2", bound: 64}
	for i := -51; i <= 51; i++ {
		for _, r := range []float64{1, 2} {
			if x := r + float64(i)/1024; x != r {
				lg.see(t, lgammaFn(x), ddLogGamma(x).Hi, fmtAt(x))
			}
		}
	}
	lg.done(t)
}

func TestSpecialBoundsBeta(t *testing.T) {
	rng := rand.New(rand.NewSource(26))
	b := &boundCase{name: "beta", bound: 2}
	lb := &boundCase{name: "lbeta", bound: 2}
	as, bs := grid(rng, 1.0/64, 80, 1.0/64, 600), grid(rng, 1.0/64, 80, 1.0/64, 600)
	for i := range as {
		a, c := as[i], bs[i]
		want := ddLogBeta(a, c)
		if a+c < 100 {
			b.see(t, betaFn(a, c), ddBeta(a, c), fmtAt(a, c))
		}
		if math.Abs(want) > 0.5 {
			lb.see(t, lbetaFn(a, c), want, fmtAt(a, c))
		}
	}
	b.done(t)
	lb.done(t)
}

func TestSpecialBoundsBessel(t *testing.T) {
	rng := rand.New(rand.NewSource(27))
	ji := &boundCase{name: "besselI", bound: 2}
	jj := &boundCase{name: "besselJ, integer order (ulps of the amplitude)", bound: 2}
	jh := &boundCase{name: "besselJ, half-integer order (ulps of the amplitude)", bound: 64}
	for _, nu := range []float64{0, 1, 2, 5, 10, 0.5, 1.5, 2.5} {
		for _, x := range grid(rng, 0.05, 60, 1.0/16, 150) {
			ji.see(t, special.Iv(nu, x), bigSeries(nu, x, +1), fmtAt(nu, x))
			// J near a zero is judged against the amplitude of the function there
			want := bigSeries(nu, x, -1)
			amp := math.Min(1, math.Sqrt(2/(math.Pi*x)))
			if nu > x {
				amp = math.Abs(want)
			}
			got := special.Jv(nu, x)
			u := math.Nextafter(amp, 2) - amp
			c := jj
			if nu != math.Floor(nu) {
				c = jh
			}
			if e := math.Abs(got-want) / u; e > c.worst {
				c.worst, c.at = e, fmtAt(nu, x)
			}
		}
	}
	ji.done(t)
	jj.done(t)
	jh.done(t)
}

// TestJvHankel: the large-x kernel agrees with special.Jv where that one is
// still sound (x up to 1e4), to 1e-13 of the amplitude, and its phase is exact
// where Jv's is not (checked against R in testdata/eval/slice25b.R).
func TestJvHankel(t *testing.T) {
	rng := rand.New(rand.NewSource(29))
	worst := 0.0
	for _, nu := range []float64{0, 1, 2, 0.5, 3.25, 10} {
		for _, x := range grid(rng, 2000, 10000, 1.0/16, 100) {
			h, ok := jvHankel(nu, x)
			if !ok {
				t.Fatalf("jvHankel(%g, %g) did not converge", nu, x)
			}
			amp := math.Sqrt(2 / (math.Pi * x))
			if e := math.Abs(h-special.Jv(nu, x)) / amp; e > worst {
				worst = e
			}
		}
	}
	t.Logf("worst %.3g of the amplitude", worst)
	if worst > 1e-13 {
		t.Errorf("jvHankel differs from Jv by %.3g of the amplitude", worst)
	}
	if _, ok := jvHankel(1e6, 1e8); ok {
		t.Error("jvHankel claimed convergence for nu^2 >> x")
	}
}

// TestJvHankelPhase: at x = 1e10, 1e15 and 1e20 the phase x - pi/4 is reduced a
// second way (the fractional part of chi/2pi, in 300-bit arithmetic) and J0
// rebuilt from it. R's own besselJ is not the oracle there: it disagrees with
// the exact phase beyond about 1e8 (README, known divergences).
func TestJvHankelPhase(t *testing.T) {
	pi := bigPi320()
	for _, x := range []float64{1e10, 1e15, 1e20} {
		chi := new(big.Float).SetPrec(300).Quo(pi, big.NewFloat(4))
		chi.Sub(new(big.Float).SetPrec(300).SetFloat64(x), chi)
		frac := new(big.Float).SetPrec(300).Quo(chi, new(big.Float).SetPrec(300).Mul(pi, big.NewFloat(2)))
		n, _ := frac.Int(nil)
		frac.Sub(frac, new(big.Float).SetPrec(300).SetInt(n))
		f, _ := frac.Float64()
		amp := math.Sqrt(2 / (math.Pi * x))
		want := amp * math.Cos(2*math.Pi*f) // leading term; P-1 and Q are below 1e-10
		got, ok := jvHankel(0, x)
		if !ok {
			t.Fatalf("no convergence at %g", x)
		}
		if e := math.Abs(got-want) / amp; e > 1e-9 {
			t.Errorf("J0(%g) = %.17g, rebuilt from the exact phase %.17g (%.3g of the amplitude)", x, got, want, e)
		}
	}
}
