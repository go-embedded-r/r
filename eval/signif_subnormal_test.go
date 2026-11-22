package eval

import "testing"

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
