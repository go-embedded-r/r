package eval

// cbind and rbind, written from cbind.Rd and experiments with the R 4.6.2
// oracle: vectors and matrices bound along a dimension, the names that come
// from the arguments, and the dispatch to a method of a class (the data frame
// methods are in baselib/dataframe.R).

type bindArg struct {
	idx    int // 1-based position among the arguments
	data   Value
	mat    bool
	nr, nc int
	n      int
	label  string // the name that heads the vector, or ""
	dn     []Value
}

func concatKind(vals []Value, kind int) Value {
	switch kind {
	case kLogical, kInteger:
		var out []int32
		for _, v := range vals {
			out = append(out, asInts(v)...)
		}
		if kind == kLogical {
			return &Logical{V: out}
		}
		return &Integer{V: out}
	case kDouble:
		var out []float64
		for _, v := range vals {
			out = append(out, v.(*Double).V...)
		}
		return &Double{V: out}
	case kString:
		var out []Str
		for _, v := range vals {
			out = append(out, v.(*String).V...)
		}
		return &String{V: out}
	case kComplex:
		var out []complex128
		for _, v := range vals {
			out = append(out, v.(*Complex).V...)
		}
		return &Complex{V: out}
	case kRaw:
		var out []byte
		for _, v := range vals {
			out = append(out, v.(*Raw).V...)
		}
		return &Raw{V: out}
	}
	var out []Value
	for _, v := range vals {
		out = append(out, v.(*List).V...)
	}
	return &List{V: out}
}

func (s *Session) bind(c *callInfo, byCol bool) Value {
	name := "rbind"
	if byCol {
		name = "cbind"
	}
	// arguments, minus deparse.level and the arguments for a method
	var vals []Value
	var tags []string
	var exprs []Value
	sameLen := len(c.call.Args) == len(c.args)
	for i, a := range c.args {
		if c.names[i] == "deparse.level" {
			continue
		}
		vals = append(vals, a)
		tags = append(tags, c.names[i])
		if sameLen {
			exprs = append(exprs, c.call.Args[i].Value)
		} else {
			exprs = append(exprs, nil)
		}
	}
	// a class with a method decides
	for _, a := range vals {
		if !isObject(a) {
			continue
		}
		for i, cl := range classAttr(a) {
			if m := s.findMethod(name+"."+cl, c.env); m != nil {
				di := &dispatchInfo{generic: name, classes: append(append([]string{}, classAttr(a)...), "default"), idx: i, object: a}
				return s.callMethod(m, name+"."+cl, c.call, argsOf(c.args, c.names), c.env, di)
			}
		}
	}
	var args []bindArg
	kind := kNull
	allRaw := true
	for i, a := range vals {
		if _, isNull := a.(nullT); isNull {
			continue
		}
		switch a.(type) {
		case *Logical, *Integer, *Double, *String, *Complex, *Raw, *List:
		default:
			s.errorf(c.call, "cannot create a matrix from type '%s'", typeName(a))
		}
		ba := bindArg{idx: i + 1, data: withAttr(a, nil), n: Length(a)}
		if d := dimsOf(a); len(d) == 2 {
			ba.mat, ba.nr, ba.nc = true, d[0], d[1]
			if dn, ok := getAttr(a, "dimnames").(*List); ok {
				ba.dn = dn.V
			}
		} else {
			if nm := namesOf(a); nm != nil {
				ba.dn = []Value{nm}
			}
			ba.label = tags[i]
			if ba.label == "" {
				if sym, ok := exprs[i].(*Symbol); ok {
					ba.label = sym.Name
				}
			}
		}
		if _, isRaw := a.(*Raw); !isRaw {
			allRaw = false
		}
		if k := kindOf(a); k != kRaw && k > kind {
			kind = k
		}
		args = append(args, ba)
	}
	if len(args) == 0 {
		return Null
	}
	if allRaw {
		kind = kRaw
	} else if kind == kNull {
		kind = kLogical
	}
	// the size of the other dimension
	m, seenMat := 0, false
	for _, a := range args {
		if a.mat {
			o := a.nr
			if !byCol {
				o = a.nc
			}
			if seenMat && o != m {
				if byCol {
					s.errorf(c.call, "number of rows of matrices must match (see arg %d)", a.idx)
				}
				s.errorf(c.call, "number of columns of matrices must match (see arg %d)", a.idx)
			}
			m, seenMat = o, true
		}
	}
	if !seenMat {
		for _, a := range args {
			m = max(m, a.n)
		}
	}
	ext := 0
	for _, a := range args {
		switch {
		case a.mat && byCol:
			ext += a.nc
		case a.mat:
			ext += a.nr
		case a.n == 0 && m > 0:
		default:
			ext++
		}
	}
	if float64(m)*float64(ext) > float64(s.opts.MaxVector) {
		s.errorf(c.call, "cannot allocate vector of length %.0f", float64(m)*float64(ext))
	}
	s.checkLen(c.call, m*ext)
	var sources []Value
	for _, a := range args {
		sources = append(sources, coerceTo(a.data, kind))
	}
	src := concatKind(sources, kind)
	pos := make([]int, m*ext)
	var bindNames []Str
	anyBind := false
	var otherNames Value
	off, base := 0, 0
	for _, a := range args {
		switch {
		case a.mat:
			k := a.nr
			if byCol {
				k = a.nc
			}
			for j := 0; j < k; j++ {
				for r := 0; r < m; r++ {
					var p int
					if byCol {
						p = base + r + j*a.nr
					} else {
						p = base + j + r*a.nr
					}
					if byCol {
						pos[r+(off+j)*m] = p
					} else {
						pos[(off+j)+r*ext] = p
					}
				}
				label := Str{S: ""}
				bi := 1
				if !byCol {
					bi = 0
				}
				if len(a.dn) == 2 {
					if st, ok := a.dn[bi].(*String); ok && j < len(st.V) {
						label = st.V[j]
						anyBind = true
					}
				}
				bindNames = append(bindNames, label)
			}
			if oi := map[bool]int{true: 0, false: 1}[byCol]; otherNames == nil && len(a.dn) == 2 {
				if st, ok := a.dn[oi].(*String); ok {
					otherNames = st
				}
			}
			off += k
		case a.n == 0 && m > 0:
		default:
			if a.n > 0 && m%a.n != 0 {
				if byCol {
					s.warningf(c.call, "number of rows of result is not a multiple of vector length (arg %d)", a.idx)
				} else {
					s.warningf(c.call, "number of columns of result is not a multiple of vector length (arg %d)", a.idx)
				}
			}
			for r := 0; r < m; r++ {
				p := base
				if a.n > 0 {
					p = base + r%a.n
				}
				if byCol {
					pos[r+off*m] = p
				} else {
					pos[off+r*ext] = p
				}
			}
			if a.label != "" {
				anyBind = true
			}
			bindNames = append(bindNames, Str{S: a.label})
			if otherNames == nil && len(a.dn) == 1 && a.n == m {
				otherNames = a.dn[0]
			}
			off++
		}
		base += a.n
	}
	out := takeElems(src, pos)
	if kind == kList {
		out = &List{V: out.(*List).V}
	}
	dims := []int32{int32(m), int32(ext)}
	if !byCol {
		dims = []int32{int32(ext), int32(m)}
	}
	out = setAttr(out, "dim", &Integer{V: dims})
	if anyBind || otherNames != nil {
		var bn Value = Null
		if anyBind {
			bn = &String{V: bindNames}
		}
		on := otherNames
		if on == nil {
			on = Null
		}
		dn := &List{V: []Value{on, bn}}
		if !byCol {
			dn.V[0], dn.V[1] = bn, on
		}
		out = setAttr(out, "dimnames", dn)
	}
	return out
}

func init() {
	regBuiltin("cbind", func(s *Session, c *callInfo) Value { return s.bind(c, true) })
	regBuiltin("rbind", func(s *Session, c *callInfo) Value { return s.bind(c, false) })
}
