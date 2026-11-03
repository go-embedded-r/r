package eval

// coerceTo promotes a vector to the given kind (logical < integer < double
// < character < list), keeping its attributes.
func coerceTo(x Value, kind int) Value {
	if kindOf(x) == kind {
		return x
	}
	attr := attrsOf(x)
	switch kind {
	case kInteger:
		return &Integer{V: append([]int32(nil), asInts(x)...), Attr: attr}
	case kDouble:
		return &Double{V: append([]float64(nil), asDoubles(x)...), Attr: attr}
	case kString:
		return &String{V: append([]Str(nil), asStrings(x)...), Attr: attr}
	case kList:
		n := Length(x)
		out := make([]Value, n)
		for i := range out {
			out[i] = element(x, i)
		}
		return &List{V: out, Attr: attr}
	}
	return x
}

// emptyOf is a zero-length vector of the kind.
func emptyOf(kind int) Value {
	switch kind {
	case kLogical:
		return &Logical{V: []int32{}}
	case kInteger:
		return &Integer{V: []int32{}}
	case kDouble:
		return &Double{V: []float64{}}
	case kString:
		return &String{V: []Str{}}
	}
	return &List{V: []Value{}}
}

// assignAt returns a copy of x (same kind as value) with the elements at pos
// replaced by those of value, recycled; the result is extended to n elements.
func assignAt(x Value, pos []int, value Value, n int) Value {
	lv := Length(value)
	switch v := x.(type) {
	case *Logical:
		out := make([]int32, n)
		copy(out, v.V)
		for i := len(v.V); i < n; i++ {
			out[i] = naInt
		}
		vv := value.(*Logical).V
		for i, p := range pos {
			if p >= 0 {
				out[p] = vv[i%lv]
			}
		}
		return &Logical{V: out, Attr: v.Attr}
	case *Integer:
		out := make([]int32, n)
		copy(out, v.V)
		for i := len(v.V); i < n; i++ {
			out[i] = naInt
		}
		vv := value.(*Integer).V
		for i, p := range pos {
			if p >= 0 {
				out[p] = vv[i%lv]
			}
		}
		return &Integer{V: out, Attr: v.Attr}
	case *Double:
		out := make([]float64, n)
		copy(out, v.V)
		for i := len(v.V); i < n; i++ {
			out[i] = naReal
		}
		vv := value.(*Double).V
		for i, p := range pos {
			if p >= 0 {
				out[p] = vv[i%lv]
			}
		}
		return &Double{V: out, Attr: v.Attr}
	case *String:
		out := make([]Str, n)
		copy(out, v.V)
		for i := len(v.V); i < n; i++ {
			out[i] = Str{NA: true}
		}
		vv := value.(*String).V
		for i, p := range pos {
			if p >= 0 {
				out[p] = vv[i%lv]
			}
		}
		return &String{V: out, Attr: v.Attr}
	case *List:
		out := make([]Value, n)
		copy(out, v.V)
		for i := len(v.V); i < n; i++ {
			out[i] = Null
		}
		vv := value.(*List).V
		for i, p := range pos {
			if p >= 0 {
				out[p] = vv[i%lv]
			}
		}
		return &List{V: out, Attr: v.Attr}
	}
	return x
}

// extendNames makes the names attribute of the extended vector.
func extendNames(x Value, old, n int, newName map[int]string) []Attr {
	attr := attrsOf(x)
	nm := namesOf(x)
	if nm == nil && len(newName) == 0 {
		return attr
	}
	names := make([]Str, n)
	for i := range names {
		if nm != nil && i < len(nm.V) {
			names[i] = nm.V[i]
		}
	}
	for p, name := range newName {
		if p < n {
			names[p] = Str{S: name}
		}
	}
	out := make([]Attr, 0, len(attr)+1)
	set := false
	for _, a := range attr {
		if a.Name == "names" {
			out = append(out, Attr{Name: "names", Value: &String{V: names}})
			set = true
		} else {
			out = append(out, a)
		}
	}
	if !set {
		out = append(out, Attr{Name: "names", Value: &String{V: names}})
	}
	return out
}

func dropDims(attr []Attr) []Attr {
	var out []Attr
	for _, a := range attr {
		if a.Name != "dim" && a.Name != "dimnames" {
			out = append(out, a)
		}
	}
	return out
}

// valueFor prepares the common kind of x and value for an assignment.
func (s *Session) commonKind(c *callInfo, x, value Value) (Value, Value) {
	kx, kv := kindOf(x), kindOf(value)
	if kv == kOther {
		kv = kList
	}
	k := kx
	if kv > k {
		k = kv
	}
	if kx == kOther {
		s.errorf(c.call, "object of type '%s' is not subsettable", typeName(x))
	}
	x = coerceTo(x, k)
	switch {
	case kindOf(value) == kOther:
		value = &List{V: []Value{value}}
	default:
		value = coerceTo(value, k)
	}
	return x, value
}

func (s *Session) subassign1(c *callInfo) Value {
	if len(c.args) < 3 {
		s.errorf(c.call, "invalid subscript in assignment")
	}
	x := c.args[0]
	value := c.args[len(c.args)-1]
	idx := c.args[1 : len(c.args)-1]
	if _, isNull := x.(nullT); isNull {
		if _, vnull := value.(nullT); vnull {
			return Null
		}
		k := kindOf(value)
		if k == kOther {
			k = kList
		}
		x = emptyOf(k)
	}
	if !isVector(x) {
		s.errorf(c.call, "object of type '%s' is not subsettable", typeName(x))
	}
	dims := dimsOf(x)
	if len(idx) >= 2 {
		if len(idx) != len(dims) {
			s.errorf(c.call, "incorrect number of subscripts")
		}
		return s.subassignArray(c, x, idx, dims, value)
	}
	n := Length(x)
	nm := namesOf(x)
	var pos []int
	var newName map[int]string
	if len(idx) == 0 {
		pos, _ = s.resolveIndex(c, missingT{}, n, nm)
	} else {
		pos, newName = s.resolveIndex(c, idx[0], n, nm)
	}
	// x[i] <- NULL deletes from a list
	if _, vnull := value.(nullT); vnull {
		l, isList := x.(*List)
		if !isList {
			if len(pos) == 0 {
				return x
			}
			s.errorf(c.call, "replacement has length zero")
		}
		drop := make([]bool, n)
		for _, p := range pos {
			if p >= 0 && p < n {
				drop[p] = true
			}
		}
		var out []Value
		var names []Str
		for i, e := range l.V {
			if !drop[i] {
				out = append(out, e)
				if nm != nil {
					names = append(names, nm.V[i])
				}
			}
		}
		res := &List{V: out, Attr: attrsOf(l)}
		if out == nil {
			res.V = []Value{}
		}
		if nm != nil {
			res.Attr = extendNamesTo(attrsOf(l), names)
		}
		return res
	}
	lv := Length(value)
	if kindOf(value) == kOther {
		lv = 1
	}
	if len(pos) == 0 {
		return x
	}
	if lv == 0 {
		s.errorf(c.call, "replacement has length zero")
	}
	hasNA := false
	maxp := n - 1
	for _, p := range pos {
		if p == naPos {
			hasNA = true
		}
		if p > maxp {
			maxp = p
		}
	}
	if hasNA && lv > 1 {
		s.errorf(c.call, "NAs are not allowed in subscripted assignments")
	}
	if len(pos)%lv != 0 {
		s.warningf(c.call, "number of items to replace is not a multiple of replacement length")
	}
	newN := maxp + 1
	s.checkLen(c.call, newN)
	cx, cv := s.commonKind(c, x, value)
	res := assignAt(cx, pos, cv, newN)
	if newN > n {
		attr := extendNames(cx, n, newN, newName)
		res = withAttr(res, dropDims(attr))
	}
	return res
}

func extendNamesTo(attr []Attr, names []Str) []Attr {
	out := make([]Attr, 0, len(attr))
	for _, a := range attr {
		if a.Name == "names" {
			out = append(out, Attr{Name: "names", Value: &String{V: names}})
		} else {
			out = append(out, a)
		}
	}
	return out
}

func (s *Session) subassignArray(c *callInfo, x Value, idx []Value, dims []int, value Value) Value {
	k := len(dims)
	var dn *List
	if l, ok := getAttr(x, "dimnames").(*List); ok {
		dn = l
	}
	sel := make([][]int, k)
	total := 1
	for d := 0; d < k; d++ {
		var names *String
		if dn != nil && d < len(dn.V) {
			names, _ = dn.V[d].(*String)
		}
		pos, _ := s.resolveIndex(c, idx[d], dims[d], names)
		for _, p := range pos {
			if p >= dims[d] || p == naPos {
				s.errorf(c.call, "subscript out of bounds")
			}
		}
		sel[d] = pos
		total *= len(pos)
	}
	lv := Length(value)
	if total == 0 {
		return x
	}
	if lv == 0 {
		s.errorf(c.call, "replacement has length zero")
	}
	if total%lv != 0 {
		s.errorf(c.call, "number of items to replace is not a multiple of replacement length")
	}
	flat := make([]int, total)
	counter := make([]int, k)
	for i := 0; i < total; i++ {
		off, stride := 0, 1
		for d := 0; d < k; d++ {
			off += sel[d][counter[d]] * stride
			stride *= dims[d]
		}
		flat[i] = off
		for d := 0; d < k; d++ {
			counter[d]++
			if counter[d] < len(sel[d]) {
				break
			}
			counter[d] = 0
		}
	}
	cx, cv := s.commonKind(c, x, value)
	return assignAt(cx, flat, cv, Length(x))
}

func (s *Session) subassign2(c *callInfo) Value {
	if len(c.args) < 3 {
		s.errorf(c.call, "[[ ]] with missing subscript")
	}
	x := c.args[0]
	value := c.args[len(c.args)-1]
	idx := c.args[1 : len(c.args)-1]
	if env, ok := x.(*Env); ok {
		if st, ok := idx[0].(*String); ok && len(idx) == 1 && len(st.V) == 1 {
			env.set(Intern(st.V[0].S), value)
			return env
		}
		s.errorf(c.call, "wrong args for environment subassignment")
	}
	if _, isNull := x.(nullT); isNull {
		if _, vnull := value.(nullT); vnull {
			return Null
		}
		k := kindOf(value)
		if k == kOther || Length(value) != 1 {
			k = kList
		}
		x = emptyOf(k)
	}
	if !isVector(x) {
		s.errorf(c.call, "object of type '%s' is not subsettable", typeName(x))
	}
	if len(idx) >= 2 {
		dims := dimsOf(x)
		if len(idx) != len(dims) {
			s.errorf(c.call, "[[ ]] improper number of subscripts")
		}
		return s.subassignArray(c, x, idx, dims, value)
	}
	if len(idx) == 0 {
		s.errorf(c.call, "[[ ]] with missing subscript")
	}
	i := idx[0]
	if _, miss := i.(missingT); miss {
		s.errorf(c.call, "[[ ]] with missing subscript")
	}
	// recursive: x[[c(1, 2)]] <- v is x[[1]][[2]] <- v
	if Length(i) > 1 {
		l, isList := x.(*List)
		if !isList {
			s.errorf(c.call, "[[ ]] subscript out of bounds")
		}
		first, rest := element(i, 0), dropFirst(i)
		sub := s.getElement(c, l, first, true)
		inner := s.subassign2(&callInfo{call: c.call, env: c.env, args: []Value{sub, rest, value}, names: []string{"", "", ""}})
		return s.subassign2(&callInfo{call: c.call, env: c.env, args: []Value{x, first, inner}, names: []string{"", "", ""}})
	}
	if Length(i) == 0 {
		s.errorf(c.call, "[[ ]] with missing subscript")
	}
	n := Length(x)
	nm := namesOf(x)
	pos, newName := s.resolveIndex(c, i, n, nm)
	if len(pos) != 1 || pos[0] == naPos {
		s.errorf(c.call, "[[ ]] subscript out of bounds")
	}
	p := pos[0]
	if f := asDoublesOrNil(i); f != nil && f[0] < 0 {
		s.errorf(c.call, "invalid negative subscript in get1index <real>")
	}
	if l, isList := x.(*List); isList {
		if _, vnull := value.(nullT); vnull {
			if p >= n {
				return x
			}
			var out []Value
			var names []Str
			for j, e := range l.V {
				if j != p {
					out = append(out, e)
					if nm != nil {
						names = append(names, nm.V[j])
					}
				}
			}
			res := &List{V: out, Attr: attrsOf(l)}
			if out == nil {
				res.V = []Value{}
			}
			if nm != nil {
				res.Attr = extendNamesTo(attrsOf(l), names)
			}
			return res
		}
		newN := n
		if p >= n {
			newN = p + 1
		}
		s.checkLen(c.call, newN)
		res := assignAt(l, []int{p}, &List{V: []Value{value}}, newN)
		if newN > n {
			res = withAttr(res, dropDims(extendNames(l, n, newN, newName)))
		}
		return res
	}
	// atomic
	lv := Length(value)
	if kindOf(value) == kOther {
		lv = 1
	}
	if _, vnull := value.(nullT); vnull {
		s.errorf(c.call, "replacement has length zero")
	}
	if kindOf(value) != kList && kindOf(value) != kOther && lv != 1 {
		if lv == 0 {
			s.errorf(c.call, "replacement has length zero")
		}
		s.errorf(c.call, "more elements supplied than there are to replace")
	}
	newN := n
	if p >= n {
		newN = p + 1
	}
	s.checkLen(c.call, newN)
	cx, cv := s.commonKind(c, x, value)
	res := assignAt(cx, []int{p}, cv, newN)
	if newN > n {
		res = withAttr(res, dropDims(extendNames(cx, n, newN, newName)))
	}
	return res
}

func asDoublesOrNil(v Value) []float64 {
	switch kindOf(v) {
	case kInteger, kDouble:
		return asDoubles(v)
	}
	return nil
}

func dropFirst(v Value) Value {
	n := Length(v)
	pos := make([]int, 0, n-1)
	for i := 1; i < n; i++ {
		pos = append(pos, i)
	}
	return takeElems(v, pos)
}

func (s *Session) subassignDollar(c *callInfo) Value {
	if len(c.args) != 3 {
		s.errorf(c.call, "invalid subscript")
	}
	x, nameV, value := c.args[0], c.args[1], c.args[2]
	st, ok := nameV.(*String)
	if !ok || len(st.V) != 1 {
		s.errorf(c.call, "invalid subscript type '%s'", typeName(nameV))
	}
	name := st.V[0].S
	switch v := x.(type) {
	case *Env:
		v.set(Intern(name), value)
		return v
	case nullT:
		if _, vnull := value.(nullT); vnull {
			return &List{V: []Value{}}
		}
		x = &List{V: []Value{}}
	case *Logical, *Integer, *Double, *String:
		s.warningf(c.call, "Coercing LHS to a list")
		x = coerceTo(x, kList)
	case *List:
	default:
		s.errorf(c.call, "object of type '%s' is not subsettable", typeName(x))
	}
	return s.subassign2(&callInfo{call: c.call, env: c.env, args: []Value{x, nameV, value}, names: []string{"", "", ""}})
}

func init() {
	regBuiltin("[<-", func(s *Session, c *callInfo) Value {
		if len(c.args) > 0 {
			if r, ok := s.dispatchInternal("[<-", c, c.args, c.names); ok {
				return r
			}
		}
		return s.subassign1(c)
	})
	regBuiltin("[[<-", func(s *Session, c *callInfo) Value {
		if len(c.args) > 0 {
			if r, ok := s.dispatchInternal("[[<-", c, c.args, c.names); ok {
				return r
			}
		}
		return s.subassign2(c)
	})
	regBuiltin("$<-", func(s *Session, c *callInfo) Value {
		if len(c.args) > 0 {
			if r, ok := s.dispatchInternal("$<-", c, c.args, c.names); ok {
				return r
			}
		}
		return s.subassignDollar(c)
	})
}
