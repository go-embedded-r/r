package eval

import (
	"math"
	"strings"

	"github.com/go-embedded-r/r/syntax"
)

// matchFormals matches the arguments of a builtin to formals: exact names
// first, then the remaining positions in order. Absent formals are nil.
func (s *Session) matchFormals(c *callInfo, formals ...string) []Value {
	out := make([]Value, len(formals))
	used := make([]bool, len(c.args))
	var unused []string
	show := func(i int) string {
		d := strings.Join(syntax.Deparse(toSyntax(c.args[i]), nil), " ")
		if c.names[i] != "" {
			return c.names[i] + " = " + d
		}
		return d
	}
	for i, n := range c.names {
		if n == "" {
			continue
		}
		for j, f := range formals {
			if f == n && out[j] == nil {
				out[j], used[i] = c.args[i], true
				break
			}
		}
		if !used[i] {
			unused = append(unused, show(i))
			used[i] = true
		}
	}
	j := 0
	for i := range c.args {
		if used[i] {
			continue
		}
		for j < len(formals) && out[j] != nil {
			j++
		}
		if j >= len(formals) {
			unused = append(unused, show(i))
			continue
		}
		out[j] = c.args[i]
	}
	switch len(unused) {
	case 0:
	case 1:
		s.errorf(c.call, "unused argument (%s)", unused[0])
	default:
		s.errorf(c.call, "unused arguments (%s)", strings.Join(unused, ", "))
	}
	return out
}

func (s *Session) countArg(c *callInfo, v Value, what string) int {
	k := kindOf(v)
	if k != kInteger && k != kDouble && k != kLogical {
		s.errorf(c.call, "non-numeric matrix extent")
	}
	d := asDoubles(v)
	if len(d) == 0 || math.IsNaN(d[0]) || d[0] > math.MaxInt32 {
		s.errorf(c.call, "invalid '%s' value (too large or NA)", what)
	}
	if d[0] < 0 {
		s.errorf(c.call, "invalid '%s' value (< 0)", what)
	}
	return int(d[0])
}

// recycleData is data recycled to n elements (NA of its type when empty).
func recycleData(data Value, n int) Value {
	l := Length(data)
	pos := make([]int, n)
	if l == 0 {
		k := kindOf(data)
		if k == kRaw {
			return &Raw{V: make([]byte, n)}
		}
		if k != kInteger && k != kDouble && k != kString && k != kList && k != kComplex {
			k = kLogical
		}
		na := coerceTo(&Logical{V: []int32{naInt}}, k)
		return takeElems(na, pos)
	}
	for i := range pos {
		pos[i] = i % l
	}
	return takeElems(data, pos)
}

func ceilDiv(a, b int) int {
	if b == 0 {
		return 0
	}
	return (a + b - 1) / b
}

func init() {
	regBuiltin("matrix", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "data", "nrow", "ncol", "byrow", "dimnames")
		data := a[0]
		if data == nil {
			data = &Logical{V: []int32{naInt}}
		}
		if _, isNull := data.(nullT); isNull {
			data = &Logical{}
		}
		if !isVector(data) {
			s.errorf(c.call, "'data' must be of a vector type, was '%s'", typeName(data))
		}
		l := Length(data)
		byrow := false
		if a[3] != nil {
			lg := logicalOf(s, c, a[3], "invalid 'byrow' argument")
			if len(lg) != 1 || lg[0] == naInt {
				s.errorf(c.call, "invalid 'byrow' argument")
			}
			byrow = lg[0] != 0
		}
		var nr, nc int
		switch {
		case a[1] == nil && a[2] == nil:
			nr, nc = l, 1
		case a[2] == nil:
			nr = s.countArg(c, a[1], "nrow")
			nc = ceilDiv(l, nr)
		case a[1] == nil:
			nc = s.countArg(c, a[2], "ncol")
			nr = ceilDiv(l, nc)
		default:
			nr, nc = s.countArg(c, a[1], "nrow"), s.countArg(c, a[2], "ncol")
		}
		if nr > 0 && nc > s.opts.MaxVector/nr {
			s.errorf(c.call, "vector size too large")
		}
		n := nr * nc
		s.checkLen(c.call, n)
		if l > 1 && n%l != 0 {
			switch {
			case (l > nr && l/nr*nr != l) || (l < nr && nr/l*l != nr):
				s.warningf(c.call, "data length [%d] is not a sub-multiple or multiple of the number of rows [%d]", l, nr)
			case (l > nc && l/nc*nc != l) || (l < nc && nc/l*l != nc):
				s.warningf(c.call, "data length [%d] is not a sub-multiple or multiple of the number of columns [%d]", l, nc)
			default:
				s.warningf(c.call, "data length differs from size of matrix: [%d != %d x %d]", l, nr, nc)
			}
		} else if l > 1 && n == 0 {
			s.warningf(c.call, "data length exceeds size of matrix")
		}
		var res Value
		if byrow && l > 0 {
			pos := make([]int, n)
			for j := 0; j < nc; j++ {
				for i := 0; i < nr; i++ {
					pos[j*nr+i] = (i*nc + j) % l
				}
			}
			res = takeElems(data, pos)
		} else {
			res = recycleData(data, n)
		}
		res = setAttr(withAttr(res, nil), "dim", &Integer{V: []int32{int32(nr), int32(nc)}})
		if a[4] != nil {
			if _, isNull := a[4].(nullT); !isNull {
				res = setAttr(res, "dimnames", s.checkDimnames(c, a[4], []int32{int32(nr), int32(nc)}))
			}
		}
		return res
	})
	regBuiltin("array", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "data", "dim", "dimnames")
		data := a[0]
		if data == nil {
			data = &Logical{V: []int32{naInt}}
		}
		if _, isNull := data.(nullT); isNull {
			s.errorf(c.call, "'data' must be of a vector type, was 'NULL'")
		}
		var dims []int32
		if a[1] == nil {
			dims = []int32{int32(Length(data))}
		} else {
			if k := kindOf(a[1]); k != kInteger && k != kDouble && k != kLogical {
				s.errorf(c.call, "'dims' cannot be of length 0")
			}
			for _, f := range asDoubles(a[1]) {
				if math.IsNaN(f) || f < 0 || f > math.MaxInt32 {
					s.errorf(c.call, "negative length vectors are not allowed")
				}
				dims = append(dims, int32(f))
			}
		}
		if len(dims) == 0 {
			s.errorf(c.call, "'dims' cannot be of length 0")
		}
		n := 1
		for _, d := range dims {
			if d > 0 && n > s.opts.MaxVector/int(d) {
				s.errorf(c.call, "'dim' specifies too large an array")
			}
			n *= int(d)
		}
		s.checkLen(c.call, n)
		res := setAttr(withAttr(recycleData(data, n), nil), "dim", &Integer{V: dims})
		if a[2] != nil {
			if _, isNull := a[2].(nullT); !isNull {
				res = setAttr(res, "dimnames", s.checkDimnames(c, a[2], dims))
			}
		}
		return res
	})
}

// checkDimnames validates the dimnames argument of matrix() and array().
func (s *Session) checkDimnames(c *callInfo, dn Value, dims []int32) Value {
	l, ok := dn.(*List)
	if !ok {
		s.errorf(c.call, "'dimnames' must be a list")
	}
	if len(l.V) > len(dims) {
		s.errorf(c.call, "length of 'dimnames' [%d] must match that of 'dims' [%d]", len(l.V), len(dims))
	}
	for i, e := range l.V {
		if _, isNull := e.(nullT); isNull {
			continue
		}
		if k := kindOf(e); k != kString && k != kInteger && k != kDouble && k != kLogical {
			s.errorf(c.call, "invalid type (%s) for 'dimnames' (must be a vector)", typeName(e))
		}
		if n := Length(e); n != 0 && n != int(dims[i]) {
			s.errorf(c.call, "length of 'dimnames' [%d] not equal to array extent", i+1)
		}
	}
	if len(l.V) < len(dims) {
		padded := &List{V: append([]Value(nil), l.V...), Attr: l.Attr}
		for len(padded.V) < len(dims) {
			padded.V = append(padded.V, Null)
		}
		if nm := namesOf(l); nm != nil {
			padded.Attr = namesAttr(fitNamesEmpty(nm.V, len(dims)))
		}
		return padded
	}
	return l
}
