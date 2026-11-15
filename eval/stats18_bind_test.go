package eval

import (
	"math"
	"math/big"
	"math/rand"
	"sort"
	"testing"
)

// The kernels that slice 21 binds to go-statistics/stats (ranks, Kendall's
// tau-b), held to the bounds stated in stats18.go against math/big.

// naiveRanks is the definition: 1 + the number of smaller values + half the
// number of other equal values.
func naiveRanks(x []float64) []float64 {
	r := make([]float64, len(x))
	for i, v := range x {
		less, eq := 0, 0
		for _, w := range x {
			if w < v {
				less++
			} else if w == v {
				eq++
			}
		}
		r[i] = float64(less) + float64(eq+1)/2
	}
	return r
}

func TestRanksAreExact(t *testing.T) {
	r := rand.New(rand.NewSource(11))
	for trial := 0; trial < 300; trial++ {
		n := 1 + r.Intn(60)
		x := make([]float64, n)
		for i := range x {
			switch r.Intn(6) {
			case 0:
				x[i] = math.Inf(1 - 2*r.Intn(2))
			case 1:
				x[i] = math.Copysign(0, float64(1-2*r.Intn(2)))
			default:
				x[i] = float64(r.Intn(8)) / 2
			}
		}
		got, want := ranks(x), naiveRanks(x)
		for i := range got {
			if got[i] != want[i] {
				t.Fatalf("trial %d: rank %d of %v is %v, want %v", trial, i, x, got[i], want[i])
			}
		}
	}
}

// TestKendallTauBound: tau-b within 8u (the stated bound of cor) of the exact
// ratio of the signed pair counts, which are integers.
func TestKendallTauBound(t *testing.T) {
	r := rand.New(rand.NewSource(12))
	worst := 0.0
	for trial := 0; trial < 200; trial++ {
		n := 2 + r.Intn(300)
		k := 1 + r.Intn(40)
		x, y := make([]float64, n), make([]float64, n)
		for i := range x {
			x[i] = float64(r.Intn(k)) + r.Float64()*float64(r.Intn(2))
			y[i] = x[i]*float64(r.Intn(3)-1) + float64(r.Intn(k))
		}
		sxy, sxx, syy := kendallSums(x, y)
		got, ok := kendallTau(x, y)
		if sxx == 0 || syy == 0 {
			if ok {
				t.Fatalf("trial %d: constant column but tau = %v", trial, got)
			}
			continue
		}
		if !ok {
			t.Fatalf("trial %d: tau refused", trial)
		}
		den := new(big.Float).SetPrec(400).Mul(bf(sxx), bf(syy))
		den.Sqrt(den)
		want := f64(new(big.Float).SetPrec(400).Quo(bf(sxy), den))
		d := math.Abs(got-want) / (8 * u53)
		worst = math.Max(worst, d)
		if d > 1 {
			t.Errorf("trial %d: tau error %g u exceeds 8u", trial, d*8)
		}
	}
	t.Logf("worst tau error: %.3g of the 8u bound", worst)
}

func TestKendallTauLarge(t *testing.T) {
	// the O(n log n) path on a size the O(n^2) kernel would take minutes for
	r := rand.New(rand.NewSource(13))
	n := 200000
	x, y := make([]float64, n), make([]float64, n)
	for i := range x {
		x[i] = r.NormFloat64()
		y[i] = x[i] + r.NormFloat64()
	}
	tau, ok := kendallTau(x, y)
	// tau = (2/pi) asin(rho), rho = 1/sqrt 2: 0.5
	if !ok || math.Abs(tau-0.5) > 0.01 {
		t.Errorf("tau = %v, %v; want about 0.5", tau, ok)
	}
	_ = sort.Float64s
}
