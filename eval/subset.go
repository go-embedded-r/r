package eval

import (
	"math"
	"strings"
)

// kindName is the type R names in messages.
func kindName(k int) string {
	return [...]string{"NULL", "logical", "integer", "double", "character", "list", "other"}[k]
}

// naPos marks an NA subscript.
const naPos = -1

// resolveIndex turns a subscript into 0-based positions for a vector of
// length n with the given names. A position at or beyond n is out of
// bounds (NA for extraction, an extension for assignment); naPos is an NA
// subscript. For a character subscript that matches no name, newName holds
// the name the new element would get.
func (s *Session) resolveIndex(c *callInfo, idx Value, n int, names *String) (pos []int, newName map[int]string) {
	switch x := idx.(type) {
	case missingT:
		pos = make([]int, n)
		for i := range pos {
			pos[i] = i
		}
		return
	case *Logical:
		m := len(x.V)
		total := n
		if m > total {
			total = m
		}
		if m == 0 {
			return nil, nil
		}
		for i := 0; i < total; i++ {
			switch x.V[i%m] {
			case naInt:
				pos = append(pos, naPos)
			case 0:
			default:
				pos = append(pos, i)
			}
		}
		return
	case *Integer, *Double:
		var vals []float64
		switch v := x.(type) {
		case *Integer:
			vals = asDoubles(v)
		case *Double:
			vals = v.V
		}
		neg, posCount := 0, 0
		for _, f := range vals {
			switch {
			case math.IsNaN(f):
				posCount++
			case f < 0 && f > -1:
			case f < 0:
				neg++
			case f >= 1:
				posCount++
			}
		}
		if neg > 0 && posCount > 0 {
			s.errorf(c.call, "only 0's may be mixed with negative subscripts")
		}
		if neg > 0 {
			drop := make([]bool, n)
			for _, f := range vals {
				if f <= -1 {
					if k := int(-f) - 1; k >= 0 && k < n {
						drop[k] = true
					}
				}
			}
			for i := 0; i < n; i++ {
				if !drop[i] {
					pos = append(pos, i)
				}
			}
			return
		}
		for _, f := range vals {
			switch {
			case math.IsNaN(f):
				pos = append(pos, naPos)
			case f >= 1:
				if f > float64(s.opts.MaxVector) {
					s.errorf(c.call, "subscript too large")
				}
				pos = append(pos, int(f)-1)
			}
		}
		return
	case nullT:
		return nil, nil
	case *String:
		newName = map[int]string{}
		next := n
		seen := map[string]int{}
		for _, e := range x.V {
			if e.NA {
				pos = append(pos, naPos)
				continue
			}
			found := -1
			if names != nil {
				for j, nm := range names.V {
					if !nm.NA && nm.S == e.S && e.S != "" {
						found = j
						break
					}
				}
			}
			if found < 0 {
				if j, ok := seen[e.S]; ok {
					found = j
				} else {
					found = next
					seen[e.S] = next
					newName[next] = e.S
					next++
				}
			}
			pos = append(pos, found)
		}
		return
	}
	s.errorf(c.call, "invalid subscript type '%s'", typeName(idx))
	return
}

// takeElems builds a vector of the type of x from the elements at pos; a
// position outside x, or naPos, is NA (NULL in a list).
func takeElems(x Value, pos []int) Value {
	switch v := x.(type) {
	case *Logical:
		out := make([]int32, len(pos))
		for i, p := range pos {
			if p < 0 || p >= len(v.V) {
				out[i] = naInt
			} else {
				out[i] = v.V[p]
			}
		}
		return &Logical{V: out}
	case *Integer:
		out := make([]int32, len(pos))
		for i, p := range pos {
			if p < 0 || p >= len(v.V) {
				out[i] = naInt
			} else {
				out[i] = v.V[p]
			}
		}
		return &Integer{V: out}
	case *Double:
		out := make([]float64, len(pos))
		for i, p := range pos {
			if p < 0 || p >= len(v.V) {
				out[i] = naReal
			} else {
				out[i] = v.V[p]
			}
		}
		return &Double{V: out}
	case *String:
		out := make([]Str, len(pos))
		for i, p := range pos {
			if p < 0 || p >= len(v.V) {
				out[i] = Str{NA: true}
			} else {
				out[i] = v.V[p]
			}
		}
		return &String{V: out}
	case *List:
		out := make([]Value, len(pos))
		for i, p := range pos {
			if p < 0 || p >= len(v.V) {
				out[i] = Null
			} else {
				out[i] = v.V[p]
			}
		}
		return &List{V: out}
	}
	return Null
}

func (s *Session) subsettable(c *callInfo, x Value) {
	if !isVector(x) {
		s.errorf(c.call, "object of type '%s' is not subsettable", typeName(x))
	}
}

// subsetArgs separates the drop and exact arguments from the subscripts.
func (s *Session) subsetArgs(c *callInfo) (x Value, idx []Value, drop, exact bool, hasExact bool) {
	drop, exact = true, true
	x = c.args[0]
	for i := 1; i < len(c.args); i++ {
		switch c.names[i] {
		case "drop":
			if l := logicalOf(s, c, c.args[i], "invalid 'drop' argument"); len(l) > 0 {
				drop = l[0] != 0
			}
		case "exact":
			hasExact = true
			if l := logicalOf(s, c, c.args[i], "invalid 'exact' argument"); len(l) > 0 {
				exact = l[0] == 1
			}
		default:
			idx = append(idx, c.args[i])
		}
	}
	return
}

func (s *Session) subset1(c *callInfo) Value {
	if len(c.args) == 0 {
		s.errorf(c.call, "invalid subsetting")
	}
	x, idx, drop, _, _ := s.subsetArgs(c)
	if _, isNull := x.(nullT); isNull {
		return Null
	}
	s.subsettable(c, x)
	dims := dimsOf(x)
	if len(idx) == 0 {
		return x
	}
	if len(idx) >= 2 {
		if len(idx) != len(dims) {
			s.errorf(c.call, "incorrect number of dimensions")
		}
		return s.subsetArray(c, x, idx, dims, drop)
	}
	n := Length(x)
	nm := namesOf(x)
	if _, miss := idx[0].(missingT); miss {
		return x
	}
	pos, _ := s.resolveIndex(c, idx[0], n, nm)
	res := takeElems(x, pos)
	var attr []Attr
	if nm != nil {
		names := make([]Str, len(pos))
		for i, p := range pos {
			if p < 0 || p >= n {
				names[i] = Str{NA: true}
			} else {
				names[i] = nm.V[p]
			}
		}
		attr = namesAttr(names)
	}
	// a factor keeps its levels and class; other objects lose their attributes
	if cl := classAttr(x); cl != nil && len(cl) == 1 && cl[0] == "factor" {
		attr = append(attr, Attr{Name: "levels", Value: getAttr(x, "levels")}, Attr{Name: "class", Value: getAttr(x, "class")})
	}
	if attr != nil {
		return withAttr(res, attr)
	}
	return res
}

// subsetArray is x[i, j, ...] on an array.
func (s *Session) subsetArray(c *callInfo, x Value, idx []Value, dims []int, drop bool) Value {
	k := len(dims)
	var dn *List
	if l, ok := getAttr(x, "dimnames").(*List); ok {
		dn = l
	}
	sel := make([][]int, k)
	for d := 0; d < k; d++ {
		var names *String
		if dn != nil && d < len(dn.V) {
			names, _ = dn.V[d].(*String)
		}
		if _, isStr := idx[d].(*String); isStr && names == nil {
			s.errorf(c.call, "no 'dimnames' attribute for array")
		}
		pos, _ := s.resolveIndex(c, idx[d], dims[d], names)
		for _, p := range pos {
			if p >= dims[d] {
				s.errorf(c.call, "subscript out of bounds")
			}
		}
		sel[d] = pos
	}
	total := 1
	for _, p := range sel {
		total *= len(p)
	}
	s.checkLen(c.call, total)
	flat := make([]int, total)
	counter := make([]int, k)
	for i := 0; i < total; i++ {
		off, stride := 0, 1
		na := false
		for d := 0; d < k; d++ {
			p := sel[d][counter[d]]
			if p == naPos {
				na = true
			}
			off += p * stride
			stride *= dims[d]
		}
		if na {
			flat[i] = naPos
		} else {
			flat[i] = off
		}
		for d := 0; d < k; d++ {
			counter[d]++
			if counter[d] < len(sel[d]) {
				break
			}
			counter[d] = 0
		}
	}
	res := takeElems(x, flat)
	newDims := make([]int, k)
	for d := range newDims {
		newDims[d] = len(sel[d])
	}
	subNames := func(d int) Value {
		if dn == nil || d >= len(dn.V) {
			return Null
		}
		nm, ok := dn.V[d].(*String)
		if !ok {
			return Null
		}
		if len(sel[d]) == 0 {
			return Null
		}
		out := make([]Str, len(sel[d]))
		for i, p := range sel[d] {
			if p == naPos {
				out[i] = Str{NA: true}
			} else {
				out[i] = nm.V[p]
			}
		}
		return &String{V: out}
	}
	if drop {
		var keep []int
		for d := 0; d < k; d++ {
			if newDims[d] != 1 {
				keep = append(keep, d)
			}
		}
		if len(keep) <= 1 {
			if len(keep) == 1 {
				if nm := subNames(keep[0]); nm != Null {
					return withAttr(res, []Attr{{Name: "names", Value: nm}})
				}
			}
			return res
		}
		if len(keep) < k {
			nd := make([]int32, len(keep))
			var ndn []Value
			any := false
			for i, d := range keep {
				nd[i] = int32(newDims[d])
				v := subNames(d)
				ndn = append(ndn, v)
				if v != Null {
					any = true
				}
			}
			attr := []Attr{{Name: "dim", Value: &Integer{V: nd}}}
			if any {
				attr = append(attr, Attr{Name: "dimnames", Value: &List{V: ndn}})
			}
			return withAttr(res, attr)
		}
	}
	nd := make([]int32, k)
	for d := range nd {
		nd[d] = int32(newDims[d])
	}
	attr := []Attr{{Name: "dim", Value: &Integer{V: nd}}}
	if dn != nil {
		ndn := make([]Value, k)
		any := false
		for d := 0; d < k; d++ {
			ndn[d] = subNames(d)
			if ndn[d] != Null {
				any = true
			}
		}
		if total == 0 {
			any = true
		}
		if any {
			var l = &List{V: ndn}
			if dnn := namesOf(dn); dnn != nil {
				l.Attr = namesAttr(dnn.V)
			}
			attr = append(attr, Attr{Name: "dimnames", Value: l})
		}
	}
	return withAttr(res, attr)
}

// getElement is x[[i]] for one subscript on a vector.
func (s *Session) getElement(c *callInfo, x Value, idx Value, exact bool) Value {
	n := Length(x)
	nm := namesOf(x)
	var p int
	switch iv := idx.(type) {
	case *String:
		if len(iv.V) > 1 {
			if _, isList := x.(*List); !isList {
				s.errorf(c.call, "attempt to select more than one element in vectorIndex")
			}
		}
		if len(iv.V) != 1 {
			s.errorf(c.call, "subscript out of bounds")
		}
		p = -2
		if !iv.V[0].NA && nm != nil {
			for j, e := range nm.V {
				if !e.NA && e.S == iv.V[0].S {
					p = j
					break
				}
			}
			if p < 0 && !exact {
				hits := 0
				for j, e := range nm.V {
					if !e.NA && strings.HasPrefix(e.S, iv.V[0].S) {
						p, hits = j, hits+1
					}
				}
				if hits != 1 {
					p = -2
				}
			}
		}
		if p == -2 {
			if _, isList := x.(*List); isList {
				return Null
			}
			s.errorf(c.call, "subscript out of bounds")
		}
	case *Logical, *Integer, *Double:
		d := asDoubles(iv)
		if len(d) != 1 {
			if len(d) == 0 {
				s.errorf(c.call, "attempt to select less than one element in get1index")
			}
			s.errorf(c.call, "attempt to select more than one element in vectorIndex")
		}
		f := d[0]
		switch {
		case math.IsNaN(f):
			if _, isList := x.(*List); isList {
				return Null
			}
			s.errorf(c.call, "subscript out of bounds")
		case f >= 1:
			p = int(f) - 1
		case f <= -1:
			if n == 2 && math.Trunc(f) == -1 {
				p = 1
			} else if n == 2 && math.Trunc(f) == -2 {
				p = 0
			} else {
				s.errorf(c.call, "invalid negative subscript in get1index%s", idxKind(idx))
			}
		default:
			if _, isLgl := iv.(*Logical); isLgl {
				s.errorf(c.call, "attempt to select less than one element in integerOneIndex")
			}
			s.errorf(c.call, "attempt to select less than one element in get1index%s", idxKind(idx))
		}
		if p >= n {
			s.errorf(c.call, "subscript out of bounds")
		}
	case nullT:
		s.errorf(c.call, "attempt to select less than one element in get1index")
	default:
		s.errorf(c.call, "invalid subscript type '%s'", typeName(idx))
	}
	if l, ok := x.(*List); ok {
		return l.V[p]
	}
	return takeElems(x, []int{p})
}

func (s *Session) subset2(c *callInfo) Value {
	if len(c.args) == 0 {
		s.errorf(c.call, "invalid subsetting")
	}
	x, idx, _, exact, _ := s.subsetArgs(c)
	if _, isNull := x.(nullT); isNull {
		return Null
	}
	if env, ok := x.(*Env); ok {
		if len(idx) == 1 {
			if st, ok := idx[0].(*String); ok && len(st.V) == 1 {
				if b := env.lookup(Intern(st.V[0].S)); b != nil {
					if p, isP := b.val.(*Promise); isP {
						return s.force(p)
					}
					return b.val
				}
				return Null
			}
		}
		s.errorf(c.call, "wrong arguments for subsetting an environment")
	}
	s.subsettable(c, x)
	if len(idx) == 0 {
		s.errorf(c.call, "invalid subscript")
	}
	if len(idx) >= 2 {
		dims := dimsOf(x)
		if len(idx) != len(dims) {
			s.errorf(c.call, "incorrect number of subscripts")
		}
		r := s.subsetArray(c, x, idx, dims, true)
		if Length(r) != 1 {
			s.errorf(c.call, "subscript out of bounds")
		}
		if l, ok := r.(*List); ok {
			return l.V[0]
		}
		return withAttr(r, nil)
	}
	i := idx[0]
	if _, miss := i.(missingT); miss {
		s.errorf(c.call, "invalid subscript type 'symbol'")
	}
	// recursive indexing of lists: x[[c(1, 2)]] is x[[1]][[2]]
	if Length(i) > 1 {
		if _, isList := x.(*List); isList {
			cur := x
			for k := 0; k < Length(i); k++ {
				if !isVector(cur) {
					s.errorf(c.call, "subscript out of bounds")
				}
				if k > 0 && Length(cur) == 0 {
					s.errorf(c.call, "subscript out of bounds")
				}
				cur = s.getElement(c, cur, element(i, k), exact)
			}
			return cur
		}
	}
	return s.getElement(c, x, i, exact)
}

func (s *Session) dollar(c *callInfo) Value {
	if len(c.raw) != 2 {
		s.errorf(c.call, "invalid subsetting")
	}
	x := s.eval(c.raw[0].Value, c.env)
	var name string
	switch n := c.raw[1].Value.(type) {
	case *Symbol:
		name = n.Name
	case *String:
		if len(n.V) != 1 {
			s.errorf(c.call, "invalid subscript type 'character'")
		}
		name = n.V[0].S
	default:
		s.errorf(c.call, "invalid subscript type '%s'", typeName(c.raw[1].Value))
	}
	if r, ok := s.dispatchInternal("$", c, []Value{x, mkString(name)}, []string{"", ""}); ok {
		return r
	}
	return s.dollarGet(c, x, name)
}

func (s *Session) dollarGet(c *callInfo, x Value, name string) Value {
	switch v := x.(type) {
	case nullT:
		return Null
	case *Lang:
		// a call is indexed by its argument names, as a pairlist is
		for _, a := range v.Args {
			if a.HasName && a.Name == name {
				return a.Value
			}
		}
		hit, hits := -1, 0
		for j, a := range v.Args {
			if a.HasName && strings.HasPrefix(a.Name, name) {
				hit, hits = j, hits+1
			}
		}
		if hits == 1 {
			return v.Args[hit].Value
		}
		return Null
	case *Env:
		if b := v.lookup(Intern(name)); b != nil {
			if p, isP := b.val.(*Promise); isP {
				return s.force(p)
			}
			return b.val
		}
		return Null
	case *List:
		nm := namesOf(v)
		if nm == nil {
			return Null
		}
		for j, e := range nm.V {
			if !e.NA && e.S == name {
				return v.V[j]
			}
		}
		hit, hits := -1, 0
		for j, e := range nm.V {
			if !e.NA && strings.HasPrefix(e.S, name) {
				hit, hits = j, hits+1
			}
		}
		if hits == 1 {
			return v.V[hit]
		}
		return Null
	case *Logical, *Integer, *Double, *String:
		s.errorf(c.call, "$ operator is invalid for atomic vectors")
	}
	s.errorf(c.call, "object of type '%s' is not subsettable", typeName(x))
	return nil
}

func init() {
	regBuiltin("[", func(s *Session, c *callInfo) Value {
		if len(c.args) > 0 {
			if r, ok := s.dispatchInternal("[", c, c.args, c.names); ok {
				return r
			}
		}
		return s.subset1(c)
	})
	regBuiltin("[[", func(s *Session, c *callInfo) Value {
		if len(c.args) > 0 {
			if r, ok := s.dispatchInternal("[[", c, c.args, c.names); ok {
				return r
			}
		}
		return s.subset2(c)
	})
	regSpecial("$", func(s *Session, c *callInfo) Value { return s.dollar(c) })
}

func idxKind(idx Value) string {
	if _, ok := idx.(*Integer); ok {
		return " <integer>"
	}
	return " <real>"
}
