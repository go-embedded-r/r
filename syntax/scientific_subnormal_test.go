package syntax

import (
	"math"
	"testing"
)

// TestScientificSubnormal holds the decimal exponent of a subnormal number:
// math.Log10 of it is wrong on amd64 (go-embedded-r/r#58), so print() wrote
// another exponent than R.
func TestScientificSubnormal(t *testing.T) {
	for _, c := range []struct {
		x            float64
		digits       int
		kpower, nsig int
		neg          int
	}{
		{1.837924e-321, 7, -321, 7, 0},
		{-6.265e-321, 7, -321, 7, 1}, // a subnormal is a multiple of 4.94e-324: 7 digits are not 6.265
		{4.9406564584124654e-324, 7, -324, 7, 0},
		{2.5e-320, 7, -320, 7, 0},
		{2.2250738585072014e-308, 7, -308, 7, 0},
		{2.225073858507201e-308, 7, -308, 7, 0}, // the largest subnormal
		{1, 7, 0, 1, 0},
		{999.9999999999999, 7, 3, 1, 0},
		{1e-300, 7, -300, 1, 0},
		{math.MaxFloat64, 7, 308, 7, 0},
	} {
		neg, kpower, nsig, _ := Scientific(c.x, c.digits)
		if neg != c.neg || kpower != c.kpower || nsig != c.nsig {
			t.Errorf("Scientific(%g, %d) = (%d, %d, %d), want (%d, %d, %d)", c.x, c.digits, neg, kpower, nsig, c.neg, c.kpower, c.nsig)
		}
	}
}

func TestDecimalExponent(t *testing.T) {
	for _, c := range []struct {
		x    float64
		want int
	}{
		{1, 0}, {9.999999999999998, 0}, {10, 1}, {0.1, -1}, {0.09999999999999999, -2},
		{1e22, 22}, {1e-22, -22}, {1.8e-321, -321}, {5e-324, -324}, {math.MaxFloat64, 308},
		{2.2250738585072014e-308, -308}, {1e100, 100}, {1e-100, -100},
	} {
		if got := decimalExponent(c.x); got != c.want {
			t.Errorf("decimalExponent(%g) = %d, want %d", c.x, got, c.want)
		}
	}
}
