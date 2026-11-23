package eval

import (
	"math"
	"math/big"
)

// signifSubnormal is signif(x, d) for a positive subnormal x and 1 <= d <= 22,
// with this contract: round the exact binary value of x to d significant
// decimal digits (ties to the even digit; none can occur here, since the exact
// expansion of a subnormal has hundreds of digits), then take the nearest
// representable double of that decimal. It is computed with integers, so the
// answer does not depend on the rounding of any intermediate product.
func signifSubnormal(x float64, d int) float64 {
	// x = k * 2^-1074 exactly.
	k := new(big.Int).SetUint64(math.Float64bits(x))
	den := new(big.Int).Lsh(big.NewInt(1), 1074)

	// e = floor(log10(x)) of the exact value: the shortest form is a start (it
	// can sit on the other side of a power of ten), corrected against 10^e.
	e := decimalExponent(x)
	for !geqPow10(k, den, e) {
		e--
	}
	for geqPow10(k, den, e+1) {
		e++
	}

	// q = round(x * 10^s), s = d - 1 - e > 0.
	s := d - 1 - e
	p := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(s)), nil)
	n := new(big.Int).Mul(k, p)
	q, r := new(big.Int).QuoRem(n, den, new(big.Int))
	r.Lsh(r, 1)
	if c := r.Cmp(den); c > 0 || (c == 0 && q.Bit(0) == 1) {
		q.Add(q, big.NewInt(1))
	}
	f, _ := new(big.Rat).SetFrac(q, p).Float64()
	return f
}

// geqPow10 reports whether k/den >= 10^e, for e < 0 (a subnormal is below
// 2^-1022, so its decimal exponent and the one above it are both negative).
func geqPow10(k, den *big.Int, e int) bool {
	p := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-e)), nil)
	return new(big.Int).Mul(k, p).Cmp(den) >= 0
}
