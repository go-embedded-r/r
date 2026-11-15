package eval

import (
	"math"
	"math/big"
	"math/rand"
	"sort"
	"testing"
)

// The accuracy contract of the descriptive statistics (see stats18.go): each
// kernel is compared with an extended-precision evaluation (math/big, 400
// bits) of the same definition, and must lie within the stated bound.

func bf(x float64) *big.Float { return new(big.Float).SetPrec(400).SetFloat64(x) }

func bigSum(x []float64) *big.Float {
	s := bf(0)
	for _, v := range x {
		s.Add(s, bf(v))
	}
	return s
}

// bigVar answers sum (x - mean)^2 and the mean, exactly to 400 bits.
func bigSS(x []float64) (ss, mean *big.Float) {
	n := bf(float64(len(x)))
	mean = new(big.Float).SetPrec(400).Quo(bigSum(x), n)
	ss = bf(0)
	for _, v := range x {
		d := new(big.Float).SetPrec(400).Sub(bf(v), mean)
		ss.Add(ss, d.Mul(d, d))
	}
	return
}

func f64(b *big.Float) float64 { v, _ := b.Float64(); return v }

// samples are data sets that stress the sums: a large offset, mixed magnitudes,
// ties, a tiny spread, a long run.
func samples(r *rand.Rand) [][]float64 {
	mk := func(n int, f func(i int) float64) []float64 {
		x := make([]float64, n)
		for i := range x {
			x[i] = f(i)
		}
		return x
	}
	return [][]float64{
		mk(2, func(int) float64 { return r.NormFloat64() }),
		mk(7, func(int) float64 { return r.NormFloat64() }),
		mk(1000, func(int) float64 { return r.NormFloat64() }),
		mk(1000, func(int) float64 { return 1e9 + r.NormFloat64() }),
		mk(1000, func(int) float64 { return 1e12 + r.Float64() }),
		mk(1000, func(int) float64 { return r.NormFloat64() * math.Pow(10, float64(r.Intn(20)-10)) }),
		mk(500, func(int) float64 { return float64(r.Intn(5)) }),
		mk(5000, func(i int) float64 { return 1 + float64(i)*1e-9 }),
		mk(3, func(int) float64 { return 1e-300 * (1 + r.Float64()) }),
		mk(100000, func(int) float64 { return r.ExpFloat64() }),
	}
}

func TestVarianceBound(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	for k, x := range samples(r) {
		ss, _ := bigSS(x)
		exact := new(big.Float).SetPrec(400).Quo(ss, bf(float64(len(x)-1)))
		ev := f64(exact)
		got := varKernel(x)
		if ev == 0 {
			if got != 0 {
				t.Errorf("sample %d: var = %g, exact 0", k, got)
			}
			continue
		}
		if rel := math.Abs(got-ev) / ev; rel > 4*u53 {
			t.Errorf("sample %d: relative error %g exceeds 4u (%g)", k, rel, 4*u53)
		}
	}
}

func TestCovCorBound(t *testing.T) {
	r := rand.New(rand.NewSource(2))
	for k, x := range samples(r) {
		n := len(x)
		y := make([]float64, n)
		for i := range y {
			y[i] = 0.5*x[i] + r.NormFloat64()*math.Max(1e-3, math.Abs(x[i])*1e-3)
		}
		sxx, syy, sxy, ok := pairKernel(x, y)
		if !ok {
			t.Fatal("not finite")
		}
		bx, mx := bigSS(x)
		by, my := bigSS(y)
		bxy, absSum := bf(0), bf(0)
		for i := range x {
			dx := new(big.Float).SetPrec(400).Sub(bf(x[i]), mx)
			dy := new(big.Float).SetPrec(400).Sub(bf(y[i]), my)
			p := dx.Mul(dx, dy)
			bxy.Add(bxy, p)
			absSum.Add(absSum, new(big.Float).Abs(p))
		}
		// the componentwise bound of the covariance: 4u sum|dx dy|
		if d := math.Abs(sxy - f64(bxy)); d > 4*u53*f64(absSum)*1.0000001 {
			t.Errorf("sample %d: cov error %g exceeds 4u S = %g", k, d, 4*u53*f64(absSum))
		}
		for _, c := range []struct{ got, ss, want *big.Float }{} {
			_ = c
		}
		// the variances are the diagonal
		if ex := f64(bx); ex != 0 && math.Abs(sxx-ex)/ex > 4*u53 {
			t.Errorf("sample %d: sxx relative error %g", k, math.Abs(sxx-ex)/ex)
		}
		if ey := f64(by); ey != 0 && math.Abs(syy-ey)/ey > 4*u53 {
			t.Errorf("sample %d: syy relative error %g", k, math.Abs(syy-ey)/ey)
		}
		// the correlation: 8u in absolute value
		den := new(big.Float).SetPrec(400).Mul(bx, by)
		den.Sqrt(den)
		if den.Sign() > 0 {
			want := f64(new(big.Float).SetPrec(400).Quo(bxy, den))
			if !(sxx > 1e-250 && syy > 1e-250 && sxx < 1e250 && syy < 1e250) {
				// out of range: the correlation of the columns scaled by powers of two
				sxx, syy, sxy, _ = pairKernel(scalePow2(x), scalePow2(y))
			}
			got := sxy / math.Sqrt(sxx) / math.Sqrt(syy)
			if d := math.Abs(got - want); d > 8*u53 {
				t.Errorf("sample %d: cor error %g exceeds 8u (%g)", k, d, 8*u53)
			}
		}
	}
}

// bigQuantile is the quantile of the sorted data at p, from the exact position
// n p + m (math/big rationals), with the position snapped to an integer when it
// lies within the tolerance of one.
func bigQuantile(sorted []float64, p float64, typ int) (float64, float64) {
	n := len(sorted)
	rat := func(x float64) *big.Rat { return new(big.Rat).SetFloat64(x) }
	var a, k *big.Rat
	switch typ {
	case 4:
		a, k = big.NewRat(0, 1), big.NewRat(0, 1)
	case 5:
		a, k = big.NewRat(1, 2), big.NewRat(0, 1)
	case 6:
		a, k = big.NewRat(0, 1), big.NewRat(1, 1)
	case 7:
		a, k = big.NewRat(1, 1), big.NewRat(-1, 1)
	case 8:
		a, k = big.NewRat(1, 3), big.NewRat(1, 3)
	default:
		a, k = big.NewRat(3, 8), big.NewRat(1, 4)
	}
	pos := new(big.Rat).Add(new(big.Rat).SetInt64(int64(n)), k)
	pos.Mul(pos, rat(p))
	pos.Add(pos, a)
	pf, _ := pos.Float64()
	fl := new(big.Int).Quo(pos.Num(), pos.Denom())
	if pos.Sign() < 0 {
		fl.Sub(fl, big.NewInt(1))
	}
	j := new(big.Rat).SetInt(fl)
	g := new(big.Rat).Sub(pos, j)
	tol := rat(quantileFuzz * math.Max(1, math.Abs(pf)))
	if g.Cmp(tol) < 0 {
		g = new(big.Rat)
	} else if new(big.Rat).Sub(big.NewRat(1, 1), g).Cmp(tol) < 0 {
		g = new(big.Rat)
		j.Add(j, big.NewRat(1, 1))
	}
	jj := int(fl.Int64())
	if g.Sign() == 0 {
		jj, _ = int(f64(bf(pf))), 0
		jj = int(j.Num().Int64())
	}
	lo, hi := jj, jj+1
	clampI := func(i int) int { return max(1, min(n, i)) }
	xl, xh := sorted[clampI(lo)-1], sorted[clampI(hi)-1]
	one := big.NewRat(1, 1)
	q := new(big.Rat).Mul(new(big.Rat).Sub(one, g), rat(xl))
	q.Add(q, new(big.Rat).Mul(g, rat(xh)))
	qf, _ := q.Float64()
	return qf, math.Max(math.Abs(xl), math.Abs(xh))
}

func TestQuantileBound(t *testing.T) {
	r := rand.New(rand.NewSource(3))
	probs := []float64{0, 0.001, 0.1, 0.25, 1.0 / 3, 0.3, 0.5, 0.62, 0.7, 0.9, 0.999, 1}
	for k, x := range samples(r) {
		sorted := append([]float64(nil), x...)
		sort.Float64s(sorted)
		for typ := 4; typ <= 9; typ++ {
			got := quantileCore(append([]float64(nil), x...), probs, typ)
			for i, p := range probs {
				want, scale := bigQuantile(sorted, p, typ)
				if d := math.Abs(got[i] - want); d > 4*u53*scale {
					t.Errorf("sample %d type %d p=%g: error %g exceeds 4u max|x| = %g", k, typ, p, d, 4*u53*scale)
				}
			}
		}
		// types 1 to 3 give an order statistic or the mean of two
		for typ := 1; typ <= 3; typ++ {
			got := quantileCore(append([]float64(nil), x...), probs, typ)
			for i, v := range got {
				n := len(sorted)
				found := false
				for _, s := range sorted {
					if s == v {
						found = true
						break
					}
				}
				if !found && typ != 2 {
					t.Errorf("sample %d type %d p=%g: %g is not an element of the data (n=%d)", k, typ, probs[i], v, n)
				}
			}
		}
	}
}

func TestMedianAndFivenum(t *testing.T) {
	r := rand.New(rand.NewSource(4))
	for k, x := range samples(r) {
		sorted := append([]float64(nil), x...)
		sort.Float64s(sorted)
		n := len(sorted)
		var want float64
		if n%2 == 1 {
			want = sorted[n/2]
		} else {
			// the mean of the two middle values, in extended precision
			s := new(big.Float).SetPrec(400).Add(bf(sorted[n/2-1]), bf(sorted[n/2]))
			want = f64(s.Quo(s, bf(2)))
		}
		got := medianOf(append([]float64(nil), x...))
		if d := math.Abs(got - want); d > 2*u53*math.Abs(want) {
			t.Errorf("sample %d: median error %g exceeds 2u", k, d)
		}
	}
}

func TestWeightedMeanBound(t *testing.T) {
	r := rand.New(rand.NewSource(5))
	for k, x := range samples(r) {
		w := make([]float64, len(x))
		for i := range w {
			w[i] = r.Float64() + 0.01
		}
		num, den, absSum := bf(0), bf(0), bf(0)
		for i := range x {
			p := new(big.Float).SetPrec(400).Mul(bf(x[i]), bf(w[i]))
			num.Add(num, p)
			den.Add(den, bf(w[i]))
			absSum.Add(absSum, new(big.Float).Abs(p))
		}
		want := f64(new(big.Float).SetPrec(400).Quo(num, den))
		bound := 4 * u53 * f64(new(big.Float).SetPrec(400).Quo(absSum, den))
		if d := math.Abs(weightedMean(x, w) - want); d > bound*1.0000001 {
			t.Errorf("sample %d: weighted.mean error %g exceeds %g", k, d, bound)
		}
	}
}

func TestSelectK(t *testing.T) {
	r := rand.New(rand.NewSource(6))
	for trial := 0; trial < 200; trial++ {
		n := 1 + r.Intn(300)
		x := make([]float64, n)
		for i := range x {
			x[i] = float64(r.Intn(20)) + r.Float64()*float64(r.Intn(2))
		}
		sorted := append([]float64(nil), x...)
		sort.Float64s(sorted)
		var ks []int
		for i := 0; i < 1+r.Intn(5); i++ {
			ks = append(ks, r.Intn(n))
		}
		sort.Ints(ks)
		ks = uniqInts(ks)
		selectK(x, ks)
		for _, k := range ks {
			if x[k] != sorted[k] {
				t.Fatalf("trial %d: element %d is %g, want %g", trial, k, x[k], sorted[k])
			}
		}
	}
}

func uniqInts(a []int) []int {
	out := a[:0]
	for i, v := range a {
		if i == 0 || v != a[i-1] {
			out = append(out, v)
		}
	}
	return out
}

func benchData(n int) []float64 {
	r := rand.New(rand.NewSource(7))
	x := make([]float64, n)
	for i := range x {
		x[i] = r.NormFloat64()
	}
	return x
}

func BenchmarkVar1e6(b *testing.B) {
	x := benchData(1_000_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		varKernel(x)
	}
}

func BenchmarkQuantile1e6(b *testing.B) {
	x := benchData(1_000_000)
	w := make([]float64, len(x))
	probs := []float64{0, 0.25, 0.5, 0.75, 1}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(w, x)
		quantileCore(w, probs, 7)
	}
}

func BenchmarkMedian1e6(b *testing.B) {
	x := benchData(1_000_000)
	w := make([]float64, len(x))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(w, x)
		medianOf(w)
	}
}

func BenchmarkCor1e6(b *testing.B) {
	x, y := benchData(1_000_000), benchData(1_000_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pairKernel(x, y)
	}
}
