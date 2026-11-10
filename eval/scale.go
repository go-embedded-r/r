package eval

import "math"

// scale.default's arithmetic, written from scale.Rd and experiments with the
// R 4.6.2 oracle. It is a Go native because the first blind R draft scored
// 0.26 exact containment against R's scale.R (gate: 0.20); see the README,
// "Clean-room process". The R side (baselib/scale_rle.R) only converts to a
// matrix and calls it.
func init() {
	regBuiltin("..scale_matrix", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "center", "scale")
		m := a[0]
		dim, _ := getAttr(m, "dim").(*Integer)
		if dim == nil || len(dim.V) != 2 {
			s.errorf(c.call, "'x' must be a matrix")
		}
		nr, nc := int(dim.V[0]), int(dim.V[1])
		// a dim attribute that disagrees with the data (R refuses to set one, but
		// a value can reach us that way) is not a matrix
		if nr < 0 || nc < 0 || (nr > 0 && nc > Length(m)/nr+1) || nr*nc != Length(m) {
			s.errorf(c.call, "'x' must be a matrix")
		}
		var colnames Value = Null
		if dn, ok := getAttr(m, "dimnames").(*List); ok && len(dn.V) == 2 {
			colnames = dn.V[1]
		}
		switch m.(type) {
		case *Logical, *Integer, *Double:
		default:
			isFlag := func(v Value, want bool) bool {
				l, ok := v.(*Logical)
				return ok && len(l.V) == 1 && l.V[0] != naInt && (l.V[0] != 0) == want
			}
			switch {
			case isFlag(a[1], true):
				s.errorf(callOf("colMeans(x, na.rm = TRUE)"), "'x' must be numeric or complex")
			case isFlag(a[1], false) && isFlag(a[2], false):
				return m
			}
			s.errorf(callOf("FUN(x, aperm(array(STATS, dims[perm]), order(perm)), ...)"), "non-numeric argument to binary operator")
		}
		x := append([]float64(nil), asDoubles(coerceTo(m, kDouble))...)
		s.checkLen(c.call, len(x))

		// vecArg resolves center or scale: a flag for TRUE/FALSE, else one
		// value per column.
		vecArg := func(v Value, what string) (vals []float64, given Value, on bool) {
			if l, ok := v.(*Logical); ok && len(l.V) == 1 {
				if l.V[0] == naInt {
					s.errorf(c.call, "missing value where TRUE/FALSE needed")
				}
				return nil, nil, l.V[0] != 0
			}
			nv := v
			switch v.(type) {
			case *Integer, *Double:
			default:
				nv = coerceTo(v, kDouble)
			}
			if Length(nv) != nc {
				s.errorf(c.call, "length of '%s' must equal the number of columns of 'x'", what)
			}
			return asDoubles(coerceTo(nv, kDouble)), nv, true
		}
		named := func(vals []float64, from Value) Value {
			out := &Double{V: vals}
			if from != nil {
				if nm := namesOf(from); nm != nil {
					out.Attr = namesAttr(nm.V)
				}
			} else if cn, ok := colnames.(*String); ok {
				out.Attr = namesAttr(cn.V)
			}
			return out
		}

		var ctr, scl Value
		cv, cgiven, con := vecArg(a[1], "center")
		if con {
			if cgiven == nil {
				cv = make([]float64, nc)
				for j := 0; j < nc; j++ {
					sum, k := 0.0, 0
					for i := 0; i < nr; i++ {
						if v := x[j*nr+i]; !math.IsNaN(v) {
							sum += v
							k++
						}
					}
					cv[j] = sum / float64(k)
				}
			}
			for j := 0; j < nc; j++ {
				for i := 0; i < nr; i++ {
					x[j*nr+i] -= cv[j]
				}
			}
			ctr = named(append([]float64(nil), cv...), cgiven)
		}
		sv, sgiven, son := vecArg(a[2], "scale")
		if son {
			if sgiven == nil {
				sv = make([]float64, nc)
				for j := 0; j < nc; j++ {
					ss, k := 0.0, 0
					for i := 0; i < nr; i++ {
						if v := x[j*nr+i]; !math.IsNaN(v) {
							ss += v * v
							k++
						}
					}
					sv[j] = math.Sqrt(ss / math.Max(1, float64(k-1)))
				}
			}
			for j := 0; j < nc; j++ {
				for i := 0; i < nr; i++ {
					x[j*nr+i] /= sv[j]
				}
			}
			scl = named(append([]float64(nil), sv...), sgiven)
		}
		if ctr == nil && scl == nil {
			return m
		}
		attr := []Attr{{Name: "dim", Value: dim}}
		if dn := getAttr(m, "dimnames"); dn != nil {
			attr = append(attr, Attr{Name: "dimnames", Value: dn})
		}
		if ctr != nil {
			attr = append(attr, Attr{Name: "scaled:center", Value: ctr})
		}
		if scl != nil {
			attr = append(attr, Attr{Name: "scaled:scale", Value: scl})
		}
		return &Double{V: x, Attr: attr}
	})
}
