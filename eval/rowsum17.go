package eval

// rowsum.default and rowsum.data.frame (slice 17) are Go natives, written from
// the help page and the oracle. Their first draft in R matched R's rowsum.R in
// a run of 53 exact tokens (limit 40) and was not reworked.

// rowsumGroups answers the distinct groups, in the order they come out
// (sorted, with a missing one last, unless reorder is FALSE), and for each row
// the index of its group.
func (s *Session) rowsumGroups(c *callInfo, group Value, reorder bool, rows int) ([]Value, []int, Value) {
	if Length(group) != rows {
		s.errorf(c.call, "incorrect length for 'group'")
	}
	miss := s.callNamed(c, "anyNA", []Value{group}, []string{""})
	if l := asInts(miss); len(l) == 1 && l[0] == 1 {
		s.warningf(c.call, "missing values for 'group'")
	}
	u := s.callNamed(c, "unique", []Value{group}, []string{""})
	if reorder {
		sorted := s.callNamed(c, "sort", []Value{u}, []string{""})
		isna := s.callNamed(c, "is.na", []Value{u}, []string{""})
		nas := s.callNamed(c, "[", []Value{u, isna}, []string{"", ""})
		u = s.callNamed(c, "c", []Value{sorted, nas}, []string{"", ""})
	}
	hit := asInts(s.callNamed(c, "match", []Value{group, u}, []string{"", ""}))
	idx := make([]int, len(hit))
	for i, h := range hit {
		idx[i] = int(h) - 1
	}
	n := Length(u)
	groups := make([]Value, n)
	for i := range groups {
		groups[i] = element(u, i)
	}
	return groups, idx, s.callNamed(c, "as.character", []Value{u}, []string{""})
}

// rowsumColumn adds the elements of one numeric column by group.
func rowsumColumn(col Value, off, rows int, idx []int, ng int, naRm bool) Value {
	switch v := col.(type) {
	case *Integer:
		out := make([]int32, ng)
		na := make([]bool, ng)
		for r := 0; r < rows; r++ {
			x := v.V[off+r]
			if x == naInt {
				if !naRm {
					na[idx[r]] = true
				}
				continue
			}
			out[idx[r]] += x
		}
		for g := range out {
			if na[g] {
				out[g] = naInt
			}
		}
		return &Integer{V: out}
	case *Double:
		out := make([]float64, ng)
		for r := 0; r < rows; r++ {
			x := v.V[off+r]
			if naRm && x != x {
				continue
			}
			out[idx[r]] += x
		}
		return &Double{V: out}
	}
	return col
}

func init() {
	regBuiltin("rowsum.default", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "group", "reorder", "na.rm")
		if a[0] == nil || a[1] == nil {
			s.errorf(c.call, "argument \"%s\" is missing, with no default", map[bool]string{true: "x", false: "group"}[a[0] == nil])
		}
		x := a[0]
		if k := kindOf(x); k != kInteger && k != kDouble {
			s.errorf(c.call, "'x' must be numeric")
		}
		reorder := a[2] == nil || s.truthy(c, a[2], "reorder")
		naRm := a[3] != nil && s.truthy(c, a[3], "na.rm")
		rows, cols := Length(x), 1
		if d := dimsOf(x); len(d) == 2 {
			rows, cols = d[0], d[1]
		}
		groups, idx, labels := s.rowsumGroups(c, a[1], reorder, rows)
		ng := len(groups)
		var parts []Value
		for k := 0; k < cols; k++ {
			parts = append(parts, rowsumColumn(x, k*rows, rows, idx, ng, naRm))
		}
		var out Value
		if kindOf(x) == kInteger {
			data := make([]int32, 0, ng*cols)
			for _, p := range parts {
				data = append(data, p.(*Integer).V...)
			}
			out = &Integer{V: data}
		} else {
			data := make([]float64, 0, ng*cols)
			for _, p := range parts {
				data = append(data, p.(*Double).V...)
			}
			out = &Double{V: data}
		}
		var cn Value = Null
		if len(dimsOf(x)) == 2 {
			if nm := dimnamesOf(x, 1); nm != nil {
				cn = nm
			}
		}
		out = withAttr(out, []Attr{
			{Name: "dim", Value: &Integer{V: []int32{int32(ng), int32(cols)}}},
			{Name: "dimnames", Value: &List{V: []Value{labels, cn}}},
		})
		return out
	})
	regBuiltin("rowsum.data.frame", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "group", "reorder", "na.rm")
		if a[0] == nil || a[1] == nil {
			s.errorf(c.call, "argument \"%s\" is missing, with no default", map[bool]string{true: "x", false: "group"}[a[0] == nil])
		}
		df, ok := a[0].(*List)
		if !ok || !inheritsClass(df, "data.frame") {
			s.errorf(c.call, "not a data frame")
		}
		reorder := a[2] == nil || s.truthy(c, a[2], "reorder")
		naRm := a[3] != nil && s.truthy(c, a[3], "na.rm")
		rows := int(rowNamesCount(getAttr(df, "row.names")))
		groups, idx, labels := s.rowsumGroups(c, a[1], reorder, rows)
		out := &List{V: make([]Value, len(df.V))}
		for i, col := range df.V {
			if k := kindOf(col); (k != kInteger && k != kDouble) || isObject(col) {
				s.errorf(c.call, "non-numeric data frame in rowsum")
			}
			out.V[i] = rowsumColumn(col, 0, rows, idx, len(groups), naRm)
		}
		out.Attr = []Attr{
			{Name: "names", Value: getAttr(df, "names")},
			{Name: "row.names", Value: labels},
			{Name: "class", Value: mkString("data.frame")},
		}
		return out
	})
}
