package eval

// Support for data frames, written from the help pages of data.frame,
// row.names and .row_names_info and from experiments with the R 4.6.2 oracle:
// the compact form of the row.names attribute (integer c(NA, -n)), which is
// what attr(x, "row.names") hides from the user.

// compactRowNames is what R stores when row.names is set to the integers 1..n.
func compactRowNames(v Value) Value {
	iv, ok := v.(*Integer)
	if !ok || len(iv.V) == 0 {
		return v
	}
	if len(iv.V) == 2 && iv.V[0] == naInt && iv.V[1] != naInt {
		return v // already compact
	}
	for i, e := range iv.V {
		if e != int32(i+1) {
			return v
		}
	}
	return &Integer{V: []int32{naInt, -int32(len(iv.V))}}
}

// rowNamesCount is the number of rows a row.names attribute describes.
func rowNamesCount(v Value) int {
	if iv, ok := v.(*Integer); ok && len(iv.V) == 2 && iv.V[0] == naInt && iv.V[1] != naInt {
		if iv.V[1] < 0 {
			return int(-iv.V[1])
		}
		return int(iv.V[1])
	}
	return Length(v)
}

// expandRowNames is attr(x, "row.names"): the compact form is 1:n.
func expandRowNames(v Value) Value {
	if iv, ok := v.(*Integer); ok && len(iv.V) == 2 && iv.V[0] == naInt && iv.V[1] != naInt {
		n := rowNamesCount(v)
		out := make([]int32, n)
		for i := range out {
			out[i] = int32(i + 1)
		}
		return &Integer{V: out}
	}
	return v
}

func init() {
	regBuiltin(".set_row_names", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "n")
		if a[0] == nil {
			s.errorf(c.call, "argument \"n\" is missing, with no default")
		}
		d := asDoubles(coerceTo(a[0], kDouble))
		if len(d) != 1 || d[0] != d[0] || d[0] <= 0 {
			return &Integer{V: []int32{}}
		}
		if d[0] > 2147483647 {
			s.errorf(c.call, "invalid 'n' argument")
		}
		return &Integer{V: []int32{naInt, -int32(d[0])}}
	})
	regBuiltin(".row_names_info", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "type")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		typ := int32(1)
		if a[1] != nil {
			if d := asDoubles(coerceTo(a[1], kDouble)); len(d) == 1 {
				typ = int32(d[0])
			}
		}
		rn := getAttr(a[0], "row.names")
		if rn == nil {
			if typ == 0 {
				return Null
			}
			return mkInt(0)
		}
		if typ == 0 {
			return rn
		}
		n := int32(rowNamesCount(rn))
		if typ == 2 {
			return mkInt(n)
		}
		if iv, ok := rn.(*Integer); ok && len(iv.V) == 2 && iv.V[0] == naInt && iv.V[1] != naInt {
			return mkInt(iv.V[1]) // negative: automatic row names
		}
		return mkInt(n)
	})
}
