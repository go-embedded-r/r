package eval

import (
	"math"
	"strings"
)

// The raw type (RAWSXP) and the bit functions, written from ?raw, ?rawConnection,
// ?bitwAnd, ?packBits and the oracle's behaviour.

// rawOperand turns a raw operand of a comparison or of match() into what it
// compares as: its hex text against a character vector, its bytes otherwise.
func rawOperand(r *Raw, other Value) Value {
	if _, ok := other.(*String); ok {
		return &String{V: asStrings(r), Attr: r.Attr}
	}
	if _, ok := other.(*Raw); ok {
		return &Integer{V: asInts(r), Attr: r.Attr}
	}
	return &Integer{V: asInts(r), Attr: r.Attr}
}

func (s *Session) rawLogic(op string, c *callInfo, x, y *Raw) Value {
	nx, ny := len(x.V), len(y.V)
	n := recycleLen(nx, ny)
	s.checkLen(c.call, n)
	if n > 0 && (n%nx != 0 || n%ny != 0) {
		s.warningf(c.call, "longer object length is not a multiple of shorter object length")
	}
	out := make([]byte, n)
	for i := range out {
		if op == "&" {
			out[i] = x.V[i%nx] & y.V[i%ny]
		} else {
			out[i] = x.V[i%nx] | y.V[i%ny]
		}
	}
	return &Raw{V: out, Attr: copyAttrs(x, y, n)}
}

// complexReals is the real parts of z, and whether any imaginary part that is
// not zero was dropped.
func complexReals(z *Complex) ([]float64, bool) {
	out := make([]float64, len(z.V))
	warn := false
	for i, e := range z.V {
		if math.IsNaN(real(e)) {
			out[i] = naReal
			continue
		}
		out[i] = real(e)
		if imag(e) != 0 {
			warn = true
		}
	}
	return out, warn
}

func (s *Session) rawArg(c *callInfo, v Value, name string) *Raw {
	r, ok := v.(*Raw)
	if !ok {
		s.errorf(c.call, "argument '%s' must be a raw vector", name)
	}
	return r
}

func (s *Session) asRaw(c *callInfo) Value {
	x := firstArg(c)
	switch v := x.(type) {
	case nullT:
		return &Raw{V: []byte{}}
	case *Raw:
		return &Raw{V: v.V}
	case *List:
		for _, e := range v.V {
			if Length(e) != 1 || !isVector(e) || kindOf(e) == kList {
				s.errorf(nil, "'list' object cannot be coerced to type 'raw'")
			}
		}
		comb := s.combine(&callInfo{call: c.call, env: c.env, args: v.V, names: make([]string, len(v.V))})
		return s.asRaw(&callInfo{call: c.call, env: c.env, args: []Value{withAttr(comb, nil)}, names: []string{""}})
	case *Logical, *Integer, *Double, *Complex, *String:
		var d []float64
		warnCoerce := false
		switch t := v.(type) {
		case *Double:
			d = t.V
		case *Complex:
			var w bool
			d, w = complexReals(t)
			if w {
				s.warningf(c.call, "imaginary parts discarded in coercion")
			}
		case *String:
			d = make([]float64, len(t.V))
			for i, e := range t.V {
				f, ok := parseNumber(e)
				if !ok {
					f, warnCoerce = naReal, true
				}
				d[i] = f
			}
		default:
			d = asDoubles(x)
		}
		if warnCoerce {
			s.warningf(c.call, "NAs introduced by coercion")
		}
		out := make([]byte, len(d))
		bad := false
		for i, f := range d {
			if math.IsNaN(f) || f < 0 || f >= 256 {
				bad = true
				continue
			}
			out[i] = byte(f)
		}
		if bad {
			s.warningf(c.call, "out-of-range values treated as 0 in coercion to raw")
		}
		return &Raw{V: out}
	}
	s.errorf(c.call, "cannot coerce type '%s' to vector of type 'raw'", typeName(x))
	return nil
}

func nulEscape(b []byte) string {
	var sb strings.Builder
	for _, e := range b {
		if e == 0 {
			sb.WriteString("\\0")
		} else {
			sb.WriteByte(e)
		}
	}
	return sb.String()
}

// intBits32 is the 32 bits of each element, least significant first.
func bitsOf(vals []uint64, width int) []byte {
	out := make([]byte, 0, len(vals)*width)
	for _, v := range vals {
		for b := 0; b < width; b++ {
			out = append(out, byte(v>>uint(b))&1)
		}
	}
	return out
}

func init() {
	regBuiltin("as.raw", func(s *Session, c *callInfo) Value { return s.asRaw(c) })
	regBuiltin("charToRaw", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "charToRaw")
		st, ok := x.(*String)
		if !ok || len(st.V) == 0 {
			s.errorf(c.call, "argument must be a character vector of length 1")
		}
		if len(st.V) > 1 {
			s.warningf(c.call, "argument should be a character vector of length 1\nall but the first element will be ignored")
		}
		e := st.V[0]
		if e.NA {
			return &Raw{V: []byte("NA")}
		}
		return &Raw{V: []byte(e.S)}
	})
	regBuiltin("rawToChar", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "multiple")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		r := s.rawArg(c, a[0], "x")
		multiple := false
		if a[1] != nil {
			l := logicalOf(s, c, a[1], "argument 'multiple' must be TRUE or FALSE")
			if len(l) != 1 || l[0] == naInt {
				s.errorf(c.call, "argument 'multiple' must be TRUE or FALSE")
			}
			multiple = l[0] == 1
		}
		s.checkLen(c.call, len(r.V))
		if multiple {
			out := make([]Str, len(r.V))
			for i, b := range r.V {
				if b != 0 {
					out[i] = Str{S: string([]byte{b})}
				}
			}
			return &String{V: out}
		}
		b := r.V
		for len(b) > 0 && b[len(b)-1] == 0 {
			b = b[:len(b)-1]
		}
		for _, e := range b {
			if e == 0 {
				s.errorf(c.call, "embedded nul in string: '%s'", nulEscape(b))
			}
		}
		return mkString(string(b))
	})
	regBuiltin("rawToBits", func(s *Session, c *callInfo) Value {
		r := s.rawArg(c, s.nargs1(c, "rawToBits"), "x")
		s.checkLen(c.call, 8*len(r.V))
		vals := make([]uint64, len(r.V))
		for i, b := range r.V {
			vals[i] = uint64(b)
		}
		return &Raw{V: bitsOf(vals, 8)}
	})
	regBuiltin("intToBits", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "intToBits")
		var iv []int32
		switch x.(type) {
		case *Logical, *Integer:
			iv = asInts(x)
		case *Double:
			iv = s.asInteger(&callInfo{call: c.call, env: c.env, args: []Value{x}, names: []string{""}}).(*Integer).V
		case *String:
			iv = s.asInteger(&callInfo{call: c.call, env: c.env, args: []Value{x}, names: []string{""}}).(*Integer).V
		default:
			s.errorf(c.call, "argument 'x' must be an integer vector")
		}
		s.checkLen(c.call, 32*len(iv))
		vals := make([]uint64, len(iv))
		for i, v := range iv {
			vals[i] = uint64(uint32(v))
		}
		return &Raw{V: bitsOf(vals, 32)}
	})
	regBuiltin("packBits", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "type")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		typ := "raw"
		if a[1] != nil {
			st, ok := a[1].(*String)
			if !ok || len(st.V) < 1 {
				s.errorf(c.call, "'arg' must be NULL or a character vector")
			}
			t := st.V[0].S
			matched := ""
			for _, cand := range []string{"raw", "integer", "double"} {
				if strings.HasPrefix(cand, t) && t != "" {
					matched = cand
				}
			}
			if matched == "" {
				s.errorf(c.call, "'arg' should be one of %s", "‘raw’, ‘integer’, ‘double’")
			}
			typ = matched
		}
		var bits []byte
		switch x := a[0].(type) {
		case *Raw:
			bits = make([]byte, len(x.V))
			for i, b := range x.V {
				bits[i] = b & 1
			}
		case *Logical, *Integer:
			iv := asInts(x)
			bits = make([]byte, len(iv))
			for i, v := range iv {
				if v == naInt {
					s.errorf(c.call, "argument 'x' must not contain NAs")
				}
				bits[i] = byte(v & 1)
			}
		default:
			s.errorf(c.call, "argument 'x' must be raw, integer or logical")
		}
		width := map[string]int{"raw": 8, "integer": 32, "double": 64}[typ]
		if len(bits)%width != 0 {
			s.errorf(c.call, "argument 'x' must be a multiple of %d long", width)
		}
		n := len(bits) / width
		word := func(i int) uint64 {
			var v uint64
			for b := 0; b < width; b++ {
				v |= uint64(bits[i*width+b]) << uint(b)
			}
			return v
		}
		switch typ {
		case "raw":
			out := make([]byte, n)
			for i := range out {
				out[i] = byte(word(i))
			}
			return &Raw{V: out}
		case "integer":
			out := make([]int32, n)
			for i := range out {
				out[i] = int32(uint32(word(i)))
			}
			return &Integer{V: out}
		}
		out := make([]float64, n)
		for i := range out {
			out[i] = math.Float64frombits(word(i))
		}
		return &Double{V: out}
	})
	regBuiltin("rawShift", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "n")
		if a[0] == nil || a[1] == nil {
			s.errorf(c.call, "argument \"%s\" is missing, with no default", map[bool]string{true: "x", false: "n"}[a[0] == nil])
		}
		r := s.rawArg(c, a[0], "x")
		var n int
		switch a[1].(type) {
		case *Logical, *Integer, *Double:
			d := asDoubles(a[1])
			if len(d) != 1 || math.IsNaN(d[0]) || d[0] < -8 || d[0] > 8 {
				s.errorf(c.call, "argument 'n' must be a small integer")
			}
			n = int(d[0])
		default:
			s.errorf(c.call, "argument 'n' must be a small integer")
		}
		out := make([]byte, len(r.V))
		for i, b := range r.V {
			if n >= 0 {
				out[i] = b << uint(n)
			} else {
				out[i] = b >> uint(-n)
			}
		}
		return &Raw{V: out}
	})

	bitw := func(name string, f func(a, b uint32) uint32) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			a := s.matchFormals(c, "a", "b")
			if a[0] == nil || a[1] == nil {
				s.errorf(c.call, "argument \"%s\" is missing, with no default", map[bool]string{true: "a", false: "b"}[a[0] == nil])
			}
			x, y := a[0], a[1]
			if kindOf(x) != kindOf(y) {
				s.errorf(c.call, "'a' and 'b' must have the same type")
			}
			xi, yi := s.bitwInts(c, name, x), s.bitwInts(c, name, y)
			nx, ny := len(xi), len(yi)
			n := recycleLen(nx, ny)
			s.checkLen(c.call, n)
			out := make([]int32, n)
			for i := range out {
				p, q := xi[i%nx], yi[i%ny]
				if p == naInt || q == naInt {
					out[i] = naInt
					continue
				}
				out[i] = int32(f(uint32(p), uint32(q)))
			}
			return &Integer{V: out}
		})
	}
	bitw("bitwAnd", func(a, b uint32) uint32 { return a & b })
	bitw("bitwOr", func(a, b uint32) uint32 { return a | b })
	bitw("bitwXor", func(a, b uint32) uint32 { return a ^ b })
	regBuiltin("bitwNot", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "bitwNot")
		xi := s.bitwInts(c, "bitwNot", x)
		out := make([]int32, len(xi))
		for i, v := range xi {
			if v == naInt {
				out[i] = naInt
			} else {
				out[i] = ^v
			}
		}
		return &Integer{V: out}
	})
	shift := func(name string, left bool) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			a := s.matchFormals(c, "a", "n")
			if a[0] == nil || a[1] == nil {
				s.errorf(c.call, "argument \"%s\" is missing, with no default", map[bool]string{true: "a", false: "n"}[a[0] == nil])
			}
			xi := s.bitwInts(c, name, a[0])
			var ni []int32
			switch a[1].(type) {
			case *Logical, *Integer:
				ni = asInts(a[1])
			case *Double:
				ni = s.asInteger(&callInfo{call: c.call, env: c.env, args: []Value{a[1]}, names: []string{""}}).(*Integer).V
			default:
				s.errorf(c.call, "unimplemented type '%s' in '%s'", typeName(a[1]), name)
			}
			nx, nn := len(xi), len(ni)
			n := recycleLen(nx, nn)
			s.checkLen(c.call, n)
			out := make([]int32, n)
			for i := range out {
				p, q := xi[i%nx], ni[i%nn]
				if p == naInt || q == naInt || q < 0 || q > 31 {
					out[i] = naInt
					continue
				}
				var r uint32
				if left {
					r = uint32(p) << uint(q)
				} else {
					r = uint32(p) >> uint(q)
				}
				out[i] = int32(r)
			}
			return &Integer{V: out}
		})
	}
	shift("bitwShiftL", true)
	shift("bitwShiftR", false)
}

// bitwInts is the integer elements of an operand of the bitw functions: an
// integer vector, or a double one truncated.
func (s *Session) bitwInts(c *callInfo, name string, x Value) []int32 {
	switch x.(type) {
	case *Integer:
		return asInts(x)
	case *Double:
		return s.asInteger(&callInfo{call: c.call, env: c.env, args: []Value{x}, names: []string{""}}).(*Integer).V
	}
	s.errorf(c.call, "unimplemented type '%s' in '%s'\n", typeName(x), name)
	return nil
}
