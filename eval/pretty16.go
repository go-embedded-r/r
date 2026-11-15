package eval

import (
	"math"
)

// pretty() (slice 16): about n+1 equally spaced round values covering a
// range, written from the help page of pretty (the unit is 1, 2, 5 or 10
// times a power of ten, chosen by the two bias arguments) and checked against
// the oracle on random ranges (rparity evalfuzz --mode pretty).

// prettyCore returns the first and last multiple of the unit, the number of
// intervals and the unit for the range [lo, hi].
func prettyCore(lo, hi float64, ndiv, minN int, shrink, h, h5 float64, epsCorrect int) (ns, nu float64, nd int, unit, lo2, hi2 float64) {
	const roundingEps = 1e-10
	dx := hi - lo
	var cell float64
	small := false
	if dx == 0 && hi == 0 {
		cell, small = 1, true
	} else {
		cell = math.Max(math.Abs(lo), math.Abs(hi))
		u := 1.5 / (1 + h5)
		if h5 >= 1.5*h+0.5 {
			u = 1 / (1 + h)
		}
		u = (1 + u) * math.Max(1, float64(ndiv)) * 2.220446049250313e-16
		small = dx < cell*u*3
	}
	if small {
		if cell > 10 {
			cell = 9 + cell/10
		}
		cell *= shrink
		if minN > 1 {
			cell /= float64(minN)
		}
	} else {
		cell = dx
		if ndiv > 1 {
			cell /= float64(ndiv)
		}
	}
	if cell < 20*2.2250738585072014e-308 { // 20 * DBL_MIN
		cell = 20 * 2.2250738585072014e-308
	} else if cell > math.MaxFloat64/1.25 {
		cell = math.MaxFloat64 / 1.25
	}
	base := math.Pow10(int(math.Floor(math.Log10(cell))))
	unit = base
	if 2*base-cell < h*(cell-unit) {
		unit = 2 * base
		if 5*base-cell < h5*(cell-unit) {
			unit = 5 * base
			if 10*base-cell < h*(cell-unit) {
				unit = 10 * base
			}
		}
	}
	ns = math.Floor(lo/unit + roundingEps)
	nu = math.Ceil(hi/unit - roundingEps)
	if epsCorrect != 0 && (epsCorrect > 1 || !small) {
		if lo != 0 {
			lo *= 1 - 2.220446049250313e-16
		} else {
			lo = -2.2250738585072014e-308
		}
		if hi != 0 {
			hi *= 1 + 2.220446049250313e-16
		} else {
			hi = 2.2250738585072014e-308
		}
	}
	for ns*unit > lo+roundingEps*unit {
		ns--
	}
	for nu*unit < hi-roundingEps*unit {
		nu++
	}
	k := int(0.5 + nu - ns)
	if k < minN {
		k = minN - k
		if ns >= 0 {
			nu += float64(k / 2)
			ns -= float64(k/2 + k%2)
		} else {
			ns -= float64(k / 2)
			nu += float64(k/2 + k%2)
		}
		nd = minN
	} else {
		nd = k
	}
	return ns, nu, nd, unit, lo, hi
}

func init() {
	// ..prettyCore(lo, hi, n, min.n, shrink.sml, high.u.bias, u5.bias, eps.correct)
	// is c(first, last, intervals, unit); the R code builds the sequence.
	regBuiltin("..prettyCore", func(s *Session, c *callInfo) Value {
		if len(c.args) != 8 {
			s.errorf(c.call, "%d arguments passed to '..prettyCore' which requires 8", len(c.args))
		}
		num := func(i int) float64 {
			v := asDoubles(c.args[i])
			if len(v) != 1 {
				s.errorf(c.call, "invalid argument %d", i+1)
			}
			return v[0]
		}
		lo, hi := num(0), num(1)
		n, minN := num(2), num(3)
		if math.IsNaN(n) || n < 0 || n > 1e9 {
			s.errorf(c.call, "invalid 'n' argument")
		}
		if math.IsNaN(minN) || minN < 0 || minN > n {
			s.errorf(c.call, "invalid 'min.n' argument")
		}
		if math.IsNaN(lo) || math.IsNaN(hi) || math.IsInf(lo, 0) || math.IsInf(hi, 0) {
			s.errorf(c.call, "invalid 'x' argument")
		}
		if math.Abs(lo) > math.MaxFloat64/4 || math.Abs(hi) > math.MaxFloat64/4 {
			s.errorf(c.call, "R_pretty(): very large range, which may give inaccurate results")
		}
		ns, nu, nd, unit, lo2, hi2 := prettyCore(lo, hi, int(n), int(minN), num(4), num(5), num(6), int(num(7)))
		if nu-ns > 1e9 {
			s.errorf(c.call, "'n' is too large for a range of this size")
		}
		return &Double{V: []float64{ns, nu, float64(nd), unit, lo2, hi2}}
	})
}
