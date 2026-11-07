package eval

import (
	"math"
)

// numeric kinds, ordered as R promotes them
const (
	kNull = iota
	kLogical
	kInteger
	kDouble
	kString
	kList
	kOther
)

func kindOf(v Value) int {
	switch v.(type) {
	case nullT:
		return kNull
	case *Logical:
		return kLogical
	case *Integer:
		return kInteger
	case *Double:
		return kDouble
	case *String:
		return kString
	case *List:
		return kList
	}
	return kOther
}

// asInts returns the elements of a logical or integer vector (NULL: none).
func asInts(v Value) []int32 {
	switch x := v.(type) {
	case *Logical:
		return x.V
	case *Integer:
		return x.V
	}
	return nil
}

// asDoubles converts a logical, integer or double vector to doubles, NA to NA.
func asDoubles(v Value) []float64 {
	switch x := v.(type) {
	case *Double:
		return x.V
	case *Logical, *Integer:
		iv := asInts(v)
		out := make([]float64, len(iv))
		for i, e := range iv {
			if e == naInt {
				out[i] = naReal
			} else {
				out[i] = float64(e)
			}
		}
		return out
	}
	return nil
}

// copyAttrs returns the attributes of the result of a binary operation on
// operands of lengths nx, ny (result length n): those of x, then those of y
// not already present, for each operand whose length is n.
func copyAttrs(x, y Value, n int) []Attr {
	var out []Attr
	add := func(v Value) {
		if Length(v) != n {
			return
		}
		for _, a := range attrsOf(v) {
			dup := false
			for _, o := range out {
				if o.Name == a.Name {
					dup = true
				}
			}
			if !dup {
				out = append(out, a)
			}
		}
	}
	add(x)
	add(y)
	return out
}

func (s *Session) checkLen(call *Lang, n int) {
	if n > s.opts.MaxVector {
		s.errorf(call, "cannot allocate vector of length %d", n)
	}
	if s.opts.MaxAlloc > 0 {
		s.alloc += 8 * int64(n)
		if s.alloc > int64(s.opts.MaxAlloc) {
			s.abort("memory limit exceeded: more than "+itoa(s.opts.MaxAlloc)+" bytes of vectors allocated", nil)
		}
	}
}

func recycleLen(nx, ny int) int {
	if nx == 0 || ny == 0 {
		return 0
	}
	if nx > ny {
		return nx
	}
	return ny
}

// conform checks the dimensions of the operands of a binary operation: two
// arrays must have the same dims, and an array must not be shorter than the
// other operand.
func (s *Session) conform(c *callInfo, x, y Value) {
	dx, dy := dimsOf(x), dimsOf(y)
	switch {
	case len(dx) > 0 && len(dy) > 0:
		if len(dx) != len(dy) {
			s.errorf(c.call, "non-conformable arrays")
		}
		for i := range dx {
			if dx[i] != dy[i] {
				s.errorf(c.call, "non-conformable arrays")
			}
		}
	case len(dx) > 0 && Length(x) > 0 && Length(y) > Length(x):
		s.errorf(nil, "dims [product %d] do not match the length of object [%d]", Length(x), Length(y))
	case len(dy) > 0 && Length(y) > 0 && Length(x) > Length(y):
		s.errorf(nil, "dims [product %d] do not match the length of object [%d]", Length(y), Length(x))
	}
}

func (s *Session) arith(op string, c *callInfo) Value {
	switch len(c.args) {
	case 1:
		return s.unary(op, c)
	case 2:
	default:
		s.errorf(c.call, "operator needs one or two arguments")
	}
	x, y := c.args[0], c.args[1]
	kx, ky := kindOf(x), kindOf(y)
	if kx > kDouble || ky > kDouble {
		s.errorf(c.call, "non-numeric argument to binary operator")
	}
	s.conform(c, x, y)
	nx, ny := Length(x), Length(y)
	n := recycleLen(nx, ny)
	s.checkLen(c.call, n)
	if n > 0 && (n%nx != 0 || n%ny != 0) {
		s.warningf(c.call, "longer object length is not a multiple of shorter object length")
	}
	attrs := copyAttrs(x, y, n)
	toDouble := kx == kDouble || ky == kDouble || op == "/" || op == "^"
	if toDouble {
		a, b := asDoubles(x), asDoubles(y)
		out := make([]float64, n)
		switch op {
		case "+":
			for i := range out {
				out[i] = a[i%nx] + b[i%ny]
			}
		case "-":
			for i := range out {
				out[i] = a[i%nx] - b[i%ny]
			}
		case "*":
			for i := range out {
				out[i] = a[i%nx] * b[i%ny]
			}
		case "/":
			for i := range out {
				out[i] = a[i%nx] / b[i%ny]
			}
		case "^":
			for i := range out {
				out[i] = rpow(a[i%nx], b[i%ny])
			}
		case "%%":
			for i := range out {
				var w bool
				out[i], w = rmod(a[i%nx], b[i%ny])
				if w {
					s.warningf(c.call, "probable complete loss of accuracy in modulus")
				}
			}
		case "%/%":
			for i := range out {
				out[i] = rintdiv(a[i%nx], b[i%ny])
			}
		}
		if kx <= kInteger && ky <= kInteger && (op == "/" || op == "^") {
			// two integer-like operands: R stores NA_real_ for an NA element
			// without computing, so its bits stay those of NA_real_; for the
			// power, 1^NA and NA^0 are 1
			ia, ib := asInts(x), asInts(y)
			for i := range out {
				p, q := ia[i%nx], ib[i%ny]
				switch {
				case op == "^" && (p == 1 || q == 0):
					out[i] = 1
				case p == naInt || q == naInt:
					out[i] = naReal
				}
			}
		}
		return &Double{V: out, Attr: attrs}
	}
	a, b := asInts(x), asInts(y)
	out := make([]int32, n)
	overflow := false
	for i := range out {
		p, q := a[i%nx], b[i%ny]
		if p == naInt || q == naInt {
			out[i] = naInt
			continue
		}
		var r int64
		switch op {
		case "+":
			r = int64(p) + int64(q)
		case "-":
			r = int64(p) - int64(q)
		case "*":
			r = int64(p) * int64(q)
		case "%%":
			if q == 0 {
				out[i] = naInt
				continue
			}
			if p >= 0 && q > 0 {
				r = int64(p % q)
			} else {
				m, _ := rmod(float64(p), float64(q))
				r = int64(m)
			}
		case "%/%":
			if q == 0 {
				out[i] = naInt
				continue
			}
			r = int64(math.Floor(float64(p) / float64(q)))
		}
		if r > math.MaxInt32 || r <= math.MinInt32 {
			out[i] = naInt
			overflow = true
		} else {
			out[i] = int32(r)
		}
	}
	if overflow {
		s.warningf(c.call, "NAs produced by integer overflow")
	}
	return &Integer{V: out, Attr: attrs}
}

func (s *Session) unary(op string, c *callInfo) Value {
	x := c.args[0]
	if op != "+" && op != "-" {
		s.errorf(c.call, "invalid unary operator")
	}
	switch v := x.(type) {
	case *Logical, *Integer:
		iv := asInts(v)
		out := make([]int32, len(iv))
		for i, e := range iv {
			switch {
			case e == naInt:
				out[i] = naInt
			case op == "-":
				out[i] = -e
			default:
				out[i] = e
			}
		}
		return &Integer{V: out, Attr: attrsOf(x)}
	case *Double:
		if op == "+" {
			return v
		}
		out := make([]float64, len(v.V))
		for i, e := range v.V {
			out[i] = -e
		}
		return &Double{V: out, Attr: v.Attr}
	}
	s.errorf(c.call, "invalid argument to unary operator")
	return nil
}

// rpow is R's R_pow: x^y with its special cases.
func rpow(x, y float64) float64 {
	if x == 1 || y == 0 {
		return 1
	}
	if x == 0 {
		switch {
		case y > 0:
			return 0
		case y < 0:
			return math.Inf(1)
		}
		return y // NA or NaN
	}
	if !math.IsInf(x, 0) && !math.IsNaN(x) && !math.IsInf(y, 0) && !math.IsNaN(y) {
		if y == 2 {
			return x * x
		}
		r := math.Pow(x, y)
		if math.IsNaN(r) {
			return rNaN // C's pow gives the hardware NaN, Go's a different payload
		}
		return r
	}
	if math.IsNaN(x) || math.IsNaN(y) {
		return x + y
	}
	if math.IsInf(x, 0) {
		if x > 0 {
			if y < 0 {
				return 0
			}
			return math.Inf(1)
		}
		if !math.IsInf(y, 0) && y == math.Floor(y) {
			if y < 0 {
				return 0
			}
			if m, _ := rmod(y, 2); m != 0 {
				return x
			}
			return -x
		}
	}
	if math.IsInf(y, 0) {
		if x >= 0 {
			if y > 0 {
				if x >= 1 {
					return math.Inf(1)
				}
				return 0
			}
			if x < 1 {
				return math.Inf(1)
			}
			return 0
		}
	}
	return rNaN
}

const dblEpsilon = 2.220446049250313e-16

// rmod is x %% y for doubles; the bool is true when R warns about a loss of
// accuracy.
func rmod(x1, x2 float64) (float64, bool) {
	if x2 == 0 {
		return rNaN, false
	}
	if math.Abs(x2)*dblEpsilon > 1 && !math.IsInf(x1, 0) && !math.IsNaN(x1) && math.Abs(x1) <= math.Abs(x2) {
		switch {
		case math.Abs(x1) == math.Abs(x2):
			return 0, false
		case (x1 < 0 && x2 > 0) || (x2 < 0 && x1 > 0):
			return x1 + x2, false
		}
		return x1, false
	}
	q := x1 / x2
	warn := false
	if !math.IsInf(q, 0) && !math.IsNaN(q) && math.Abs(q)*dblEpsilon > 1 {
		warn = true
	}
	// R compiled for arm64 fuses these multiply-subtracts (a NaN negated by the
	// fusion keeps the sign bit); written as FMA, the result is the same on every
	// platform
	// (the multiplicands are passed in the order that gives the negated
	// floor(q) priority when both are NaN, as in R's compiled code)
	tmp := math.FMA(x2, -math.Floor(q), x1)
	q = math.Floor(tmp / x2)
	return math.FMA(x2, -q, tmp), warn
}

// rintdiv is x %/% y for doubles.
func rintdiv(x1, x2 float64) float64 {
	q := x1 / x2
	if x2 == 0 || math.Abs(q)*dblEpsilon > 1 || math.IsInf(q, 0) || math.IsNaN(q) {
		return q
	}
	if math.Abs(q) < 1 {
		switch {
		case q < 0:
			return -1
		case (x1 < 0 && x2 > 0) || (x1 > 0 && x2 < 0):
			return -1
		}
		return 0
	}
	tmp := math.FMA(x2, -math.Floor(q), x1)
	return math.Floor(q) + math.Floor(tmp/x2)
}

func init() {
	for _, op := range []string{"+", "-", "*", "/", "^", "%%", "%/%"} {
		op := op
		regBuiltin(op, func(s *Session, c *callInfo) Value { return s.arith(op, c) })
	}
}
