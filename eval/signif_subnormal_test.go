package eval

import (
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"testing"
)

// TestSignifSubnormal holds the exponent of a subnormal number: math.Log10 of
// it is wrong on amd64, so signif used to return 0.
func TestSignifSubnormal(t *testing.T) {
	for _, c := range []struct{ x, want float64 }{
		{1.837924e-321, 1.837924e-321},
		{-6.265e-321, -6.265e-321},
		{4.9406564584124654e-324, 4.9406564584124654e-324},
	} {
		if got := fsignif(c.x, 12); got != c.want {
			t.Errorf("fsignif(%g, 12) = %g, want %g", c.x, got, c.want)
		}
	}
}

// signifContract is the contract of signif(x, d) for a subnormal x, written
// without the code under test: the exact binary value is rounded to d
// significant decimal digits by math/big (to nearest, ties to even), and the
// decimal is read back as the nearest double.
func signifContract(t *testing.T, x float64, d int) float64 {
	t.Helper()
	exact := new(big.Float).SetPrec(2000).SetFloat64(x)
	f, err := strconv.ParseFloat(exact.Text('e', d-1), 64)
	if err != nil && !math.IsInf(f, 0) {
		t.Fatalf("ParseFloat: %v", err)
	}
	return f
}

// The subnormals: every one of the first 3000, every power-of-ten boundary
// (the largest subnormal under 10^e and the smallest over it, where a shortest
// decimal form can sit on the other side), and random ones across the range;
// each with d from 1 to 22 and both signs.
func TestSignifSubnormalContract(t *testing.T) {
	var xs []float64
	for b := uint64(1); b <= 3000; b++ {
		xs = append(xs, math.Float64frombits(b))
	}
	for e := -323; e <= -308; e++ {
		p, _ := strconv.ParseFloat("1e"+strconv.Itoa(e), 64)
		b := math.Float64bits(p)
		for _, c := range []uint64{b - 2, b - 1, b, b + 1, b + 2} {
			if c >= 1 && c < 1<<52 {
				xs = append(xs, math.Float64frombits(c))
			}
		}
	}
	xs = append(xs, math.Float64frombits(1<<52-1), math.Float64frombits(1<<52-2))
	rng := rand.New(rand.NewSource(61))
	for i := 0; i < 3000; i++ {
		xs = append(xs, math.Float64frombits(1+uint64(rng.Int63n(1<<52-1))))
	}
	for _, x := range xs {
		for d := 1; d <= 22; d++ {
			want := signifContract(t, x, d)
			if got := fsignif(x, float64(d)); got != want {
				t.Fatalf("fsignif(%g [bits %#x], %d) = %g, contract %g", x, math.Float64bits(x), d, got, want)
			}
			if got := fsignif(-x, float64(d)); got != -want {
				t.Fatalf("fsignif(-%g, %d) = %g, contract %g", x, d, got, -want)
			}
		}
	}
}

// The recorded divergence: R 4.6.2 (linux/amd64) answers these a few units of
// 2^-1074 away from the rounding of the exact value (README, "signif of a
// subnormal"). The test holds both numbers, so a change on either side shows.
func TestSignifSubnormalDivergesFromR(t *testing.T) {
	for _, c := range []struct {
		bits  uint64
		d     int
		r, we float64
	}{
		{1<<52 - 1, 1, 1.9999999999999993e-308, 1.9999999999999998e-308},
		{1 << 51, 1, 9.9999999999999942e-309, 9.9999999999999991e-309},
		{1629120193449413, 1, 7.9999999999999953e-309, 8.0000000000000003e-309},
	} {
		x := math.Float64frombits(c.bits)
		if got := fsignif(x, float64(c.d)); got != c.we {
			t.Errorf("fsignif(%g, %d) = %.17g, want %.17g", x, c.d, got, c.we)
		}
		if signifContract(t, x, c.d) != c.we || c.r == c.we {
			t.Errorf("case %#x: the contract value or R's value is wrong", c.bits)
		}
	}
}

// d that is not a whole number is truncated (as fsignif always did), and a d
// below 1 is 1.
func TestSignifSubnormalDigits(t *testing.T) {
	x := 1.234567e-315
	want := signifContract(t, x, 3)
	if got := fsignif(x, 3.9); got != want {
		t.Errorf("fsignif(x, 3.9) = %g, want %g", got, want)
	}
	if got, w1 := fsignif(x, 0), signifContract(t, x, 1); got != w1 {
		t.Errorf("fsignif(x, 0) = %g, want %g", got, w1)
	}
}
