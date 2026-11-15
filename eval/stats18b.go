package eval

import (
	"math"
	"sort"
	"strconv"
	"strings"
)

// quantile, fivenum, IQR, mad and weighted.mean (slice 18), over the kernels of
// stats18.go.

// fgFormat is formatC(x, format = "fg", width = 1, digits = d): d significant
// digits in fixed notation, the trailing zeros dropped.
func fgFormat(x float64, d int) string {
	if x == 0 {
		return "0"
	}
	if math.IsInf(x, 0) {
		if x > 0 {
			return "Inf"
		}
		return "-Inf"
	}
	e := strconv.FormatFloat(math.Abs(x), 'e', d-1, 64)
	mant, exps, _ := strings.Cut(e, "e")
	exp, _ := strconv.Atoi(exps)
	digits := strings.Replace(mant, ".", "", 1)
	var out string
	switch {
	case exp >= len(digits)-1:
		out = digits + strings.Repeat("0", exp-len(digits)+1)
	case exp >= 0:
		out = digits[:exp+1] + "." + digits[exp+1:]
	default:
		out = "0." + strings.Repeat("0", -exp-1) + digits
	}
	if strings.Contains(out, ".") {
		out = strings.TrimRight(out, "0")
		out = strings.TrimSuffix(out, ".")
	}
	if x < 0 {
		out = "-" + out
	}
	return out
}

func (s *Session) statIntArg(c *callInfo, v Value, def int, what string) int {
	if v == nil {
		return def
	}
	d := asDoublesOrNil(v)
	if len(d) != 1 || d[0] != d[0] {
		s.errorf(c.call, "invalid '%s' argument", what)
	}
	return int(d[0])
}

// quantileCore answers the quantiles of the plain numbers x (no missing
// values; modified) at probs.
func quantileCore(x []float64, probs []float64, typ int) []float64 {
	out := make([]float64, len(probs))
	n := len(x)
	if n == 0 {
		for i := range out {
			out[i] = naReal
		}
		return out
	}
	selectK(x, quantileIndices(n, probs, typ))
	for i, p := range probs {
		if p != p {
			out[i] = naReal
			continue
		}
		out[i] = quantileAt(x, p, typ)
	}
	return out
}

func init() {
	regBuiltin("quantile.default", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "probs", "na.rm", "names", "type", "digits", "...")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		x := a[0]
		naRm := s.boolArg(c, a[2], false, "na.rm")
		names := s.boolArg(c, a[3], true, "names")
		typ := s.statIntArg(c, a[4], 7, "type")
		digits := s.statIntArg(c, a[5], 7, "digits")
		probs := []float64{0, 0.25, 0.5, 0.75, 1}
		if a[1] != nil {
			p, ok := numericVector(a[1])
			if !ok {
				s.errorf(c.call, "non-numeric argument to binary operator")
			}
			probs = p
		}
		ordered := inheritsFrom(x, "ordered")
		if inheritsFrom(x, "factor") {
			if !ordered {
				s.errorf(c.call, "(unordered) factors are not allowed")
			}
			if typ != 1 && typ != 3 {
				s.errorf(c.call, "'type' must be 1 or 3 for ordered factors")
			}
		}
		if typ < 1 || typ > 9 {
			s.errorf(c.call, "'type' must be an integer in 1..9")
		}
		xs, ok := numericVector(x)
		if !ok {
			s.errorf(c.call, "non-numeric argument to binary operator")
		}
		if naRm {
			xs = dropNaN(xs)
		} else if hasNaN(xs) {
			s.errorf(c.call, "missing values and NaN's not allowed if 'na.rm' is FALSE")
		}
		eps := 100 * 0x1p-52
		clipped := make([]float64, len(probs))
		for i, p := range probs {
			if p == p && (p < -eps || p > 1+eps) {
				s.errorf(c.call, "'probs' outside [0,1]")
			}
			clipped[i] = math.Max(0, math.Min(1, p))
			if p != p {
				clipped[i] = p
			}
		}
		work := append([]float64(nil), xs...)
		vals := quantileCore(work, clipped, typ)
		var res Value = &Double{V: vals}
		if classAttr(x) != nil && typ <= 3 {
			// the order statistics keep the class of the data (ordered factors,
			// dates)
			var attrs []Attr
			for _, at := range attrsOf(x) {
				if at.Name != "names" && at.Name != "dim" && at.Name != "dimnames" {
					attrs = append(attrs, at)
				}
			}
			if ordered {
				iv := make([]int32, len(vals))
				for i, v := range vals {
					if v != v {
						iv[i] = naInt
					} else {
						iv[i] = int32(v)
					}
				}
				res = &Integer{V: iv}
			}
			res = withAttr(res, attrs)
		}
		if names {
			nm := make([]Str, len(probs))
			d := max(2, digits)
			for i, p := range probs {
				if p == p {
					nm[i] = Str{S: fgFormat(100*p, d) + "%"}
				}
			}
			res = withAttr(res, append(append([]Attr(nil), attrsOf(res)...), Attr{Name: "names", Value: &String{V: nm}}))
		}
		return res
	})

	regBuiltin("fivenum", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "na.rm")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		naRm := s.boolArg(c, a[1], true, "na.rm")
		xs, ok := numericVector(a[0])
		if !ok {
			s.errorf(c.call, "non-numeric argument to binary operator")
		}
		if hasNaN(xs) {
			if !naRm {
				return naLogical5()
			}
			xs = dropNaN(xs)
		}
		n := len(xs)
		if n == 0 {
			return naLogical5()
		}
		w := append([]float64(nil), xs...)
		sort.Float64s(w)
		n4 := math.Floor(float64(n+3)/2) / 2
		d := []float64{1, n4, float64(n+1) / 2, float64(n) + 1 - n4, float64(n)}
		out := make([]float64, 5)
		for i, e := range d {
			out[i] = 0.5 * (w[int(math.Floor(e))-1] + w[int(math.Ceil(e))-1])
		}
		return &Double{V: out}
	})

	regBuiltin("IQR", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "na.rm", "type")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		naRm := s.boolArg(c, a[1], false, "na.rm")
		typ := s.statIntArg(c, a[2], 7, "type")
		xs, ok := numericVector(a[0])
		if !ok {
			s.errorf(c.call, "non-numeric argument to binary operator")
		}
		if naRm {
			xs = dropNaN(xs)
		} else if hasNaN(xs) {
			s.errorf(callOf("quantile.default(as.numeric(x), c(0.25, 0.75), na.rm = na.rm, names = FALSE, type = type)"),
				"missing values and NaN's not allowed if 'na.rm' is FALSE")
		}
		if typ < 1 || typ > 9 {
			s.errorf(callOf("quantile.default(as.numeric(x), c(0.25, 0.75), na.rm = na.rm, names = FALSE, type = type)"),
				"'type' must be an integer in 1..9")
		}
		q := quantileCore(append([]float64(nil), xs...), []float64{0.25, 0.75}, typ)
		return mkDouble(q[1] - q[0])
	})

	regBuiltin("mad", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "center", "constant", "na.rm", "low", "high")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		naRm := s.boolArg(c, a[3], false, "na.rm")
		low := s.boolArg(c, a[4], false, "low")
		high := s.boolArg(c, a[5], false, "high")
		constant := 1.4826
		if a[2] != nil {
			d := asDoublesOrNil(a[2])
			if len(d) != 1 {
				s.errorf(c.call, "invalid '%s' argument", "constant")
			}
			constant = d[0]
		}
		xs, ok := numericVector(a[0])
		if !ok {
			s.errorf(c.call, "non-numeric argument to binary operator")
		}
		if naRm {
			xs = dropNaN(xs)
		}
		if low && high {
			s.errorf(c.call, "'low' and 'high' cannot be both TRUE")
		}
		var center float64
		if a[1] != nil {
			d := asDoublesOrNil(a[1])
			if len(d) != 1 {
				s.errorf(c.call, "invalid '%s' argument", "center")
			}
			center = d[0]
		} else if hasNaN(xs) {
			center = naReal
		} else {
			center = medianOf(append([]float64(nil), xs...))
		}
		n := len(xs)
		if n == 0 {
			return mkDouble(naReal)
		}
		dev := make([]float64, n)
		for i, v := range xs {
			dev[i] = math.Abs(v - center)
		}
		if hasNaN(dev) {
			return mkDouble(naReal)
		}
		var m float64
		switch {
		case n%2 == 1:
			m = medianOf(dev)
		case low:
			selectK(dev, []int{n/2 - 1})
			m = dev[n/2-1]
		case high:
			selectK(dev, []int{n / 2})
			m = dev[n/2]
		default:
			m = medianOf(dev)
		}
		return mkDouble(constant * m)
	})

	regBuiltin("weighted.mean.default", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "w", "na.rm", "...")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		naRm := s.boolArg(c, a[2], false, "na.rm")
		if a[1] == nil {
			return s.callNamed(c, "mean", []Value{a[0], mkLogical(naRm)}, []string{"", "na.rm"})
		}
		xs, ok1 := numericVector(a[0])
		ws, ok2 := numericVector(a[1])
		if !ok1 || !ok2 {
			s.errorf(callOf("x * w"), "non-numeric argument to binary operator")
		}
		if len(xs) != len(ws) {
			s.errorf(c.call, "'x' and 'w' must have the same length")
		}
		if naRm {
			var xk, wk []float64
			for i, v := range xs {
				if v == v {
					xk, wk = append(xk, v), append(wk, ws[i])
				}
			}
			xs, ws = xk, wk
		}
		return mkDouble(weightedMean(xs, ws))
	})
}

func naLogical5() Value {
	return &Logical{V: []int32{naInt, naInt, naInt, naInt, naInt}}
}

// weightedMean is sum(x w) / sum(w) with the products added exactly (fused
// products, compensated sums) when everything is finite.
func weightedMean(x, w []float64) float64 {
	allFinite := true
	for i := range x {
		if !finite(x[i]) || !finite(w[i]) {
			allFinite = false
			break
		}
	}
	if !allFinite {
		var num, den float64
		for i := range x {
			num += x[i] * w[i]
			den += w[i]
		}
		return num / den
	}
	var s, c, ws, wc float64
	add := func(a, comp *float64, v float64) {
		t := *a + v
		if math.Abs(*a) >= math.Abs(v) {
			*comp += (*a - t) + v
		} else {
			*comp += (v - t) + *a
		}
		*a = t
	}
	for i := range x {
		p := float64(x[i] * w[i])
		c += math.FMA(x[i], w[i], -p)
		add(&s, &c, p)
		add(&ws, &wc, w[i])
	}
	return (s + c) / (ws + wc)
}
