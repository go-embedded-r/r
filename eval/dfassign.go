package eval

import "strconv"

// `[<-.data.frame` and the matching step of transform: a Go native, because
// a first draft of each in R (baselib/dfmethods.R) converged on R's own text
// (exact runs of 67 and 86 tokens against dataframe.R and frametools.R; see
// README, "Clean-room process"). Written from Extract.data.frame.Rd and
// experiments with the R 4.6.2 oracle. The columns themselves are assigned
// by the R generic `[<-` (so that factors and dates keep their methods),
// through the small helper ..dfsetcol of baselib/dataframe.R.

// dfIndex resolves one subscript against n existing positions. names, when
// not nil, is what a character subscript is matched against; the positions
// of names that are not there are numbered after n and returned in added.
func (s *Session) dfIndex(c *callInfo, sub Value, n int, names []Str) (pos []int, added []Str) {
	switch v := sub.(type) {
	case *String:
		for _, e := range v.V {
			found := -1
			for k, nm := range names {
				if !nm.NA && !e.NA && nm.S == e.S {
					found = k
					break
				}
			}
			if found < 0 {
				for k, a := range added {
					if a.S == e.S {
						found = n + k
					}
				}
			}
			if found < 0 {
				added = append(added, e)
				found = n + len(added) - 1
			}
			pos = append(pos, found)
		}
	case nullT:
		return nil, nil
	case *Logical:
		if len(v.V) == 0 {
			return nil, nil
		}
		for i := 0; i < max(n, len(v.V)); i++ {
			if v.V[i%len(v.V)] == 1 {
				pos = append(pos, i)
			}
		}
	case *Integer, *Double:
		d := asDoubles(coerceTo(v, kDouble))
		neg := false
		for _, e := range d {
			neg = neg || e < 0
		}
		if neg {
			drop := map[int]bool{}
			for _, e := range d {
				if e > 0 {
					s.errorf(c.call, "can't mix positive and negative subscripts")
				}
				drop[int(-e)-1] = true
			}
			for i := 0; i < n; i++ {
				if !drop[i] {
					pos = append(pos, i)
				}
			}
			return pos, nil
		}
		for _, e := range d {
			if e >= 1 {
				pos = append(pos, int(e)-1)
			}
		}
	default:
		s.errorf(c.call, "invalid subscript type '%s'", typeName(sub))
	}
	return pos, added
}

func init() {
	regBuiltin("[<-.data.frame", func(s *Session, c *callInfo) Value {
		n := len(c.args)
		if n < 2 || n > 4 {
			s.errorf(c.call, "need 0, 1, or 2 subscripts")
		}
		x, value := c.args[0], c.args[n-1]
		sub := c.args[1 : n-1]
		var iSub, jSub Value
		switch len(sub) {
		case 2:
			iSub, jSub = sub[0], sub[1]
		case 1:
			jSub = sub[0]
		}
		hasI := iSub != nil && !isMissingValue(iSub)
		hasJ := jSub != nil && !isMissingValue(jSub)
		cols, _ := x.(*List)
		if cols == nil {
			s.errorf(c.call, "invalid data frame")
		}
		rn := getAttr(x, "row.names")
		nrows := rowNamesCount(rn)
		autoRN := false
		if iv, ok := rn.(*Integer); ok && len(iv.V) == 2 && iv.V[0] == naInt && iv.V[1] < 0 {
			autoRN = true
		}
		var cn []Str
		if nm := namesOf(x); nm != nil {
			cn = append(cn, nm.V...)
		}
		ncol0 := len(cn)
		vals := append([]Value(nil), cols.V...)
		var jj []int
		if hasJ {
			var added []Str
			jj, added = s.dfIndex(c, jSub, ncol0, cn)
			cn = append(cn, added...)
			for _, p := range jj {
				for len(cn) <= p {
					cn = append(cn, Str{S: "V" + strconv.Itoa(len(cn)+1)})
				}
			}
		} else {
			for k := range vals {
				jj = append(jj, k)
			}
		}
		attrs := func(nms []Str, rnv Value, rnFirst bool) []Attr {
			names := Attr{Name: "names", Value: &String{V: nms}}
			row := Attr{Name: "row.names", Value: rnv}
			class := Attr{Name: "class", Value: getAttr(x, "class")}
			if rnFirst {
				return []Attr{row, names, class}
			}
			return []Attr{names, row, class}
		}
		if _, isNull := value.(nullT); isNull && hasJ && !hasI {
			drop := map[int]bool{}
			for _, p := range jj {
				drop[p] = true
			}
			var keepV []Value
			var keepN []Str
			for k := range vals {
				if !drop[k] {
					keepV = append(keepV, vals[k])
					keepN = append(keepN, cn[k])
				}
			}
			return &List{V: keepV, Attr: attrs(keepN, rn, false)}
		}
		var ii []int
		var newRN []Str
		if hasI {
			var rnames []Str
			if st, ok := expandRowNames(rn).(*String); ok {
				rnames = st.V
			} else {
				for k := 0; k < nrows; k++ {
					rnames = append(rnames, Str{S: strconv.Itoa(k + 1)})
				}
				if iv, ok := expandRowNames(rn).(*Integer); ok {
					for k, e := range iv.V {
						rnames[k] = Str{S: strconv.Itoa(int(e))}
					}
				}
			}
			switch sv := iSub.(type) {
			case *Logical:
				for _, e := range sv.V {
					if e == naInt {
						s.errorf(c.call, "missing values are not allowed in subscripted assignments of data frames")
					}
				}
			case *Integer, *Double:
				for _, e := range asDoubles(coerceTo(sv, kDouble)) {
					if e != e {
						s.errorf(c.call, "missing values are not allowed in subscripted assignments of data frames")
					}
				}
			}
			ii, newRN = s.dfIndex(c, iSub, nrows, rnames)
		} else {
			for k := 0; k < nrows; k++ {
				ii = append(ii, k)
			}
		}
		newn := nrows
		for _, p := range ii {
			newn = max(newn, p+1)
		}
		s.checkLen(c.call, newn)
		nrep := len(ii)
		// the values, one per column
		var per []Value
		if l, ok := value.(*List); ok {
			per = l.V
		} else if m := Length(value); m > nrep && nrep > 0 && len(jj) > 1 {
			for k := range jj {
				pos := make([]int, nrep)
				for r := range pos {
					pos[r] = (k*nrep + r) % m
				}
				per = append(per, takeElems(withAttr(value, nil), pos))
			}
		} else {
			per = []Value{value}
		}
		if len(per) == 0 {
			s.errorf(c.call, "replacement has length zero")
		}
		if l, ok := value.(*List); ok && len(l.V) > len(jj) && len(jj) > 0 && !(hasI && len(ii) == 0 && allBelow(jj, ncol0)) {
			s.warningf(c.call, "provided %d variables to replace %d variables", len(l.V), len(jj))
		}
		setcol := s.findMethod("..dfsetcol", s.Base)
		touched := map[int]bool{}
		for k, p := range jj {
			if hasI && len(ii) == 0 && allBelow(jj, ncol0) {
				continue // no row to assign, no new column: the columns are left as they are
			}
			var col Value
			if p < len(vals) {
				col = vals[p]
			} else {
				col = &Logical{V: constNA(nrows)}
			}
			for len(vals) <= p {
				vals = append(vals, Null)
			}
			idx := make([]int32, len(ii))
			for r, e := range ii {
				idx[r] = int32(e + 1)
			}
			vals[p] = s.callFunction(setcol, []Value{col, &Integer{V: idx}, per[k%len(per)], mkLogical(hasI), mkInt(int32(newn))})
			touched[p] = true
		}
		if newn > nrows {
			pad := s.findMethod("..dfsetcol", s.Base)
			for k := range vals {
				if !touched[k] {
					vals[k] = s.callFunction(pad, []Value{vals[k], &Integer{V: []int32{}}, Null, mkLogical(false), mkInt(int32(newn))})
				}
			}
			if len(newRN) == 0 && autoRN {
				seq := make([]int32, newn)
				for k := range seq {
					seq[k] = int32(k + 1)
				}
				rn = compactRowNames(&Integer{V: seq})
			} else {
				var all []Str
				if st, ok := expandRowNames(rn).(*String); ok {
					all = append(all, st.V...)
				} else if iv, ok := expandRowNames(rn).(*Integer); ok {
					for _, e := range iv.V {
						all = append(all, Str{S: strconv.Itoa(int(e))})
					}
				}
				if len(newRN) > 0 {
					all = append(all, newRN...)
				}
				for len(all) < newn {
					all = append(all, Str{S: strconv.Itoa(len(all) + 1)})
				}
				rn = &String{V: all}
			}
		}
		return &List{V: vals, Attr: attrs(cn, rn, len(cn) > ncol0)}
	})

	// ..subsetDF(x, rows, cols, drop): the selection step of subset() for a
	// data frame: rows is NULL or a logical vector (NA counts as FALSE), cols
	// NULL or a subscript already evaluated.
	regBuiltin("..subsetDF", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "rows", "cols", "drop")
		var rows, cols Value = missingT{}, missingT{}
		if _, isNull := a[1].(nullT); !isNull {
			lg, ok := a[1].(*Logical)
			if !ok {
				s.errorf(callOf("subset.data.frame(x, subset, select)"), "'subset' must be logical")
			}
			v := make([]int32, len(lg.V))
			for i, e := range lg.V {
				if e == 1 {
					v[i] = 1
				}
			}
			rows = &Logical{V: v}
		}
		if _, isNull := a[2].(nullT); !isNull {
			cols = a[2]
		}
		args := []Value{a[0], rows, cols, a[3]}
		return builtins["["].fn(s, &callInfo{call: c.call, env: c.env, args: args, names: []string{"", "", "", "drop"}})
	})

	// ..transformDF(data, e): the columns of data named in the list e are
	// replaced, the others are returned in the second element.
	regBuiltin("..transformDF", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "data", "e")
		data, ok1 := a[0].(*List)
		e, ok2 := a[1].(*List)
		if !ok1 || !ok2 {
			s.errorf(c.call, "invalid arguments")
		}
		nrows := rowNamesCount(getAttr(data, "row.names"))
		dn := namesOf(data)
		en := namesOf(e)
		vals := append([]Value(nil), data.V...)
		var extra []Value
		var extraN []Str
		gone := map[int]bool{}
		for k, v := range e.V {
			name := Str{}
			if en != nil {
				name = en.V[k]
			}
			hit := -1
			if dn != nil {
				for j, nm := range dn.V {
					if nm.S == name.S && !nm.NA {
						hit = j
						break
					}
				}
			}
			_, isNull := v.(nullT)
			if hit < 0 {
				extra = append(extra, v)
				extraN = append(extraN, name)
				continue
			}
			if isNull {
				gone[hit] = true
				continue
			}
			if m := Length(v); m != nrows && m > 0 && nrows > 0 {
				if m > nrows || nrows%m != 0 {
					s.errorf(callOf("`[<-.data.frame`(`*tmp*`, hit, value = e[matched])"), "replacement has %d row%s, data has %d", m, map[bool]string{true: "", false: "s"}[m == 1], nrows)
				}
				pos := make([]int, nrows)
				for r := range pos {
					pos[r] = r % m
				}
				v = takeElems(v, pos)
			}
			vals[hit] = v
		}
		out := &List{V: vals, Attr: attrsOf(data)}
		if len(gone) > 0 {
			var keepV []Value
			var keepN []Str
			for j := range vals {
				if !gone[j] {
					keepV = append(keepV, vals[j])
					keepN = append(keepN, dn.V[j])
				}
			}
			out = &List{V: keepV}
			for _, at := range attrsOf(data) {
				if at.Name == "names" {
					at.Value = &String{V: keepN}
				}
				out.Attr = append(out.Attr, at)
			}
		}
		ex := &List{V: extra}
		if len(extra) > 0 {
			ex.Attr = namesAttr(extraN)
		}
		return &List{V: []Value{out, ex}}
	})
}

func allBelow(pos []int, n int) bool {
	for _, p := range pos {
		if p >= n {
			return false
		}
	}
	return true
}

func constNA(n int) []int32 {
	out := make([]int32, n)
	for i := range out {
		out[i] = naInt
	}
	return out
}
