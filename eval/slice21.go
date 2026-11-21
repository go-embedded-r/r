package eval

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode/utf8"
)

// regexec: the match and its parenthesised groups, one integer vector of
// start positions per text, with the lengths in the match.length attribute.
func init() {
	regBuiltin("regexec", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "pattern", "text", "ignore.case", "perl", "fixed", "useBytes")
		pat, okp := s.patternOf(c, a[0])
		x := s.textOf(c, a[1])
		fl := s.reFlags(c, map[string]Value{"ignore.case": a[2], "perl": a[3], "fixed": a[4]})
		var m *matcher
		if okp {
			m = s.compileRegexGrep(c, pat, x, fl.fixed, fl.icase, fl.perl)
		}
		ascii := okp && allASCII(pat, x)
		out := make([]Value, len(x))
		for i, e := range x {
			if !okp || e.NA {
				out[i] = &Integer{V: []int32{naInt}, Attr: []Attr{{Name: "match.length", Value: &Integer{V: []int32{naInt}}}}}
				continue
			}
			loc := m.find(e.S, 0)
			if loc == nil {
				out[i] = &Integer{V: []int32{-1}, Attr: matchAttrs([]int32{-1}, ascii)}
				continue
			}
			n := len(loc) / 2
			pos := make([]int32, n)
			lens := make([]int32, n)
			for k := 0; k < n; k++ {
				if loc[2*k] < 0 {
					pos[k], lens[k] = 0, 0
					continue
				}
				pos[k] = int32(charIndex(e.S, loc[2*k]) + 1)
				lens[k] = int32(utf8.RuneCountInString(e.S[loc[2*k]:loc[2*k+1]]))
			}
			out[i] = &Integer{V: pos, Attr: matchAttrs(lens, ascii)}
		}
		return &List{V: out}
	})
}

// .approxeval(x, y, v, method, yleft, yright, f): the step or linear
// interpolation of the sorted knots (x, y) at v (method 1 linear, 2 constant).
// A NA or NaN stays itself; outside the knots the answer is yleft or yright.
func init() {
	regBuiltin("..approxeval", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "y", "v", "method", "yleft", "yright", "f")
		for _, v := range a {
			if v == nil {
				s.errorf(c.call, "7 arguments required")
			}
		}
		x, y, v := asDoubles(a[0]), asDoubles(a[1]), asDoubles(a[2])
		one := func(v Value) float64 {
			d := asDoubles(v)
			if len(d) == 0 {
				return naReal
			}
			return d[0]
		}
		method, yl, yr, f := int(one(a[3])), one(a[4]), one(a[5]), one(a[6])
		n := len(x)
		out := make([]float64, len(v))
		for k, t := range v {
			switch {
			case math.IsNaN(t) || n == 0 || len(y) != n:
				out[k] = t
				if n == 0 || len(y) != n {
					out[k] = naReal
				}
			case t < x[0]:
				out[k] = yl
			case t > x[n-1]:
				out[k] = yr
			default:
				// the last knot not above t
				lo, hi := 0, n-1
				for hi-lo > 1 {
					mid := (lo + hi) / 2
					if x[mid] <= t {
						lo = mid
					} else {
						hi = mid
					}
				}
				if x[hi] == t {
					lo = hi
				}
				switch {
				case t == x[lo]:
					out[k] = y[lo]
				case method == 2:
					switch f {
					case 0:
						out[k] = y[lo]
					case 1:
						out[k] = y[hi]
					default:
						out[k] = y[lo]*(1-f) + y[hi]*f
					}
				default:
					out[k] = y[lo] + (y[hi]-y[lo])*((t-x[lo])/(x[hi]-x[lo]))
				}
			}
		}
		return &Double{V: out}
	})
}

// ..ecdfknots(x): the knots of the empirical distribution function of x, the
// sorted distinct values without NA, with the proportion of observations at or
// below each, as list(x, y, n). Each proportion is one correctly rounded
// division of an exact count by n, so its error is at most half a unit in the
// last place (the test asserts it against math/big).
func init() {
	regBuiltin("..ecdfknots", func(s *Session, c *callInfo) Value {
		if len(c.args) != 1 {
			s.errorf(c.call, "1 argument required")
		}
		switch c.args[0].(type) {
		case *Double, *Integer, *Logical:
		default:
			s.errorf(c.call, "'x' must be numeric")
		}
		in := asDoubles(c.args[0])
		xs := make([]float64, 0, len(in))
		for _, v := range in {
			if !math.IsNaN(v) {
				xs = append(xs, v)
			}
		}
		if len(xs) < 1 {
			s.errorf(c.call, "'x' must have 1 or more non-missing values")
		}
		sort.Float64s(xs)
		n := float64(len(xs))
		var vals, cum []float64
		for i, v := range xs {
			if i+1 < len(xs) && xs[i+1] == v {
				continue
			}
			vals = append(vals, v)
			cum = append(cum, float64(i+1)/n)
		}
		l := &List{V: []Value{&Double{V: vals}, &Double{V: cum}, &Integer{V: []int32{int32(len(xs))}}}}
		l.Attr = []Attr{{Name: "names", Value: strVec([]string{"x", "y", "n"})}}
		return l
	})
}

// ..ecdfline(x, digits): the line print.ecdf gives under its call, the first
// three knots, the last two and the count, each formatted like formatC with
// the digits significant digits (a Go native: the draft in R converged on R's
// own text).
func init() {
	regBuiltin("..ecdfline", func(s *Session, c *callInfo) Value {
		if len(c.args) != 2 {
			s.errorf(c.call, "2 arguments required")
		}
		xs := asDoubles(c.args[0])
		d := asDoubles(c.args[1])
		dig := 5
		if len(d) == 1 && !math.IsNaN(d[0]) && d[0] >= 1 && d[0] <= 22 {
			dig = int(d[0])
		}
		n := len(xs)
		form := func(lo, hi int) string {
			var b strings.Builder
			for i := lo; i < hi; i++ {
				if i > lo {
					b.WriteString(", ")
				}
				f := cFormatG(xs[i], dig)
				// formatC pads each knot to digits + 1 characters
				for k := len(f); k < dig+1; k++ {
					b.WriteByte(' ')
				}
				b.WriteString(f)
			}
			return b.String()
		}
		var b strings.Builder
		fmt.Fprintf(&b, " x[1:%d] = %s", n, form(0, min(3, n)))
		if n > 3 {
			b.WriteString(", ")
		}
		if n > 5 {
			b.WriteString(" ..., ")
		}
		if n >= 4 {
			b.WriteString(form(max(4, n-1)-1, n))
		}
		b.WriteString("\n")
		return mkString(b.String())
	})
}

// Natives of slice 21: pairlist (the list of a call's arguments, typeof
// "pairlist") and its helper for as.pairlist.

func init() {
	// sys.parents: for each frame up to the caller's, the number of its parent frame
	regBuiltin("sys.parents", func(s *Session, c *callInfo) Value {
		cur, fr := s.frameIndex(c.env)
		if fr == nil {
			return &Integer{}
		}
		out := make([]int32, 0, cur+1)
		for i := 0; i <= cur; i++ {
			p := int32(0)
			if ce := s.frames[i].caller; ce != nil && ce != s.Global {
				if idx, f := s.frameIndex(ce); f != nil {
					p = int32(idx + 1)
				}
			}
			out = append(out, p)
		}
		return &Integer{V: out}
	})
	regBuiltin("pairlist", func(s *Session, c *callInfo) Value {
		v := s.list(c)
		l, ok := v.(*List)
		if !ok || len(l.V) == 0 {
			return Null
		}
		l.Pair = true
		return l
	})
	regBuiltin("..mkpair", func(s *Session, c *callInfo) Value {
		if len(c.args) != 1 {
			s.errorf(c.call, "1 argument required")
		}
		l, ok := c.args[0].(*List)
		if !ok || len(l.V) == 0 {
			return Null
		}
		return &List{V: l.V, Attr: l.Attr, Pair: true}
	})
}
