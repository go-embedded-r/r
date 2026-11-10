package eval

import "math"

// coerceTo promotes a vector to the given kind (logical < integer < double
// < character < list), keeping its attributes.
func coerceTo(x Value, kind int) Value {
	if kindOf(x) == kind {
		return x
	}
	attr := attrsOf(x)
	switch kind {
	case kLogical:
		// a downward conversion, for the functions that ask for logical
		d := asDoubles(x)
		out := make([]int32, len(d))
		for i, e := range d {
			switch {
			case math.IsNaN(e):
				out[i] = naInt
			case e != 0:
				out[i] = 1
			}
		}
		return &Logical{V: out, Attr: attr}
	case kInteger:
		if dv, ok := x.(*Double); ok {
			out := make([]int32, len(dv.V))
			for i, e := range dv.V {
				if math.IsNaN(e) || e >= 2147483648 || e <= -2147483649 {
					out[i] = naInt
				} else {
					out[i] = int32(e)
				}
			}
			return &Integer{V: out, Attr: attr}
		}
		return &Integer{V: append([]int32(nil), asInts(x)...), Attr: attr}
	case kDouble:
		return &Double{V: append([]float64(nil), asDoubles(x)...), Attr: attr}
	case kComplex:
		return &Complex{V: append([]complex128(nil), asComplex(x)...), Attr: attr}
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
	case kComplex:
		return &Complex{V: []complex128{}}
	case kRaw:
		return &Raw{V: []byte{}}
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
		return &List{V: out, Attr: v.Attr, Expr: v.Expr}
	case *Complex:
		out := make([]complex128, n)
		copy(out, v.V)
		for i := len(v.V); i < n; i++ {
			out[i] = naComplex
		}
		vv := value.(*Complex).V
		for i, p := range pos {
			if p >= 0 {
				out[p] = vv[i%lv]
			}
		}
		return &Complex{V: out, Attr: v.Attr}
	case *Raw:
		out := make([]byte, n)
		copy(out, v.V)
		vv := value.(*Raw).V
		for i, p := range pos {
			if p >= 0 {
				out[p] = vv[i%lv]
			}
		}
		return &Raw{V: out, Attr: v.Attr}
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
	if kx == kRaw || kv == kRaw {
		// raw mixes with nothing but itself and lists
		switch {
		case kx == kRaw && kv == kRaw:
			k = kRaw
		case kx == kList || kv == kList:
			k = kList
		case kx == kNull || kv == kNull:
			k = kx + kv - kNull
		default:
			s.errorf(c.call, "incompatible types (from %s to %s) in subassignment type fix", typeName(x), typeName(value))
		}
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

// commonKindFlat is commonKind for a subscript that treats x as a plain
// vector: coercing an atomic array to a list keeps its names only.
func (s *Session) commonKindFlat(c *callInfo, x, value Value) (Value, Value) {
	kx := kindOf(x)
	cx, cv := s.commonKind(c, x, value)
	if _, isList := cx.(*List); isList && kx != kList {
		var keep []Attr
		for _, a := range attrsOf(cx) {
			if a.Name == "names" {
				keep = append(keep, a)
			}
		}
		cx = withAttr(cx, keep)
	}
	return cx, cv
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
		if l, isL := value.(*List); isL && len(l.V) == 0 {
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
			if len(idx) == 2 {
				s.errorf(c.call, "incorrect number of subscripts on matrix")
			}
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
	_, xIsL := x.(*List)
	if _, vnull := value.(nullT); vnull && (xIsL || len(pos) > 0) {
		l, isList := x.(*List)
		if !isList {
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
		res := &List{V: out, Attr: attrsOf(l), Expr: l.Expr}
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
	// a logical subscript longer than x extends x, even when it selects nothing
	ext := 0
	if len(idx) > 0 {
		if lg, ok := idx[0].(*Logical); ok {
			ext = len(lg.V)
		}
	}
	if len(pos) == 0 {
		cx := x
		if _, vnull := value.(nullT); !vnull {
			cx, _ = s.commonKindFlat(c, x, value)
		}
		if ext > n {
			s.checkLen(c.call, ext)
			res := assignAt(cx, nil, emptyOf(kindOf(cx)), ext)
			return withAttr(res, dropDims(extendNames(cx, n, ext, nil)))
		}
		return cx
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
	if ext > newN {
		newN = ext
	}
	s.checkLen(c.call, newN)
	cx, cv := s.commonKindFlat(c, x, value)
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
		if _, isStr := idx[d].(*String); isStr && names == nil {
			call := c.call
			if c.tmpCall != nil {
				call = c.tmpCall
			}
			s.errorf(call, "no 'dimnames' attribute for array")
		}
		pos, _ := s.resolveIndex(c, idx[d], dims[d], names)
		var kept []int
		for _, p := range pos {
			if p == naPos {
				if _, vnull := value.(nullT); vnull || Length(value) > 1 {
					s.errorf(c.call, "NAs are not allowed in subscripted assignments")
				}
				continue
			}
			if p >= dims[d] {
				call := c.call
				if c.tmpCall != nil {
					call = c.tmpCall
				}
				s.errorf(call, "subscript out of bounds")
			}
			kept = append(kept, p)
		}
		sel[d] = kept
		total *= len(kept)
	}
	lv := Length(value)
	if total == 0 {
		if _, vnull := value.(nullT); vnull {
			return x
		}
		cx, _ := s.commonKindFlat(c, x, value)
		return cx
	}
	if lv == 0 {
		if _, vnull := value.(nullT); vnull {
			s.errorf(c.call, "number of items to replace is not a multiple of replacement length")
		}
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
	cx, cv := s.commonKindFlat(c, x, value)
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
			s.assignTo(env, Intern(st.V[0].S), value, c.call)
			return env
		}
		s.errorf(c.call, "wrong args for environment subassignment")
	}
	wasNull := false
	if _, isNull := x.(nullT); isNull {
		if _, vnull := value.(nullT); vnull {
			return Null
		}
		x = emptyOf(kList)
		wasNull = true
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
	// the first subscript decides several of R's error messages
	firstSub := func() {
		if k := kindOf(i); (k == kLogical || k == kInteger || k == kDouble) && Length(i) > 0 {
			nx := Length(x)
			f := asDoubles(i)[0]
			single := Length(i) == 1
			less := func(kind string) {
				s.errorf(c.call, "attempt to select less than one element in %s", kind)
			}
			switch {
			case k == kLogical && single && f == 0:
				less("integerOneIndex")
			case single && math.IsNaN(f):
				if nx == 0 {
					less("integerOneIndex")
				}
				s.errorf(c.call, "attempt to select more than one element in integerOneIndex")
			case k != kLogical && f > -1 && f < 1 && !math.IsNaN(f):
				less("OneIndex" + idxKind(i))
			case k != kLogical && f <= -1 && nx <= 1:
				less("OneIndex" + idxKind(i))
			case k != kLogical && f <= -1 && nx >= 3:
				s.errorf(c.call, "attempt to select more than one element in OneIndex%s", idxKind(i))
			}
			if k == kLogical && single && f == 1 {
				i = mkInt(1)
			}
		}
	}
	if Length(i) > 1 {
		firstSub()
	}
	// recursive: x[[c(1, 2)]] <- v is x[[1]][[2]] <- v
	if Length(i) > 1 {
		l, isList := x.(*List)
		if wasNull {
			call := c.call
			if c.tmpCall != nil {
				call = c.tmpCall
			}
			s.errorf(call, "no such index at level 1\n")
		}
		if !isList {
			call := c.call
			if c.tmpCall != nil {
				call = c.tmpCall
			}
			s.errorf(call, "attempt to select more than one element in vectorIndex")
		}
		first, rest := element(i, 0), dropFirst(i)
		sub := s.getElement(c, l, first, true)
		inner := s.subassign2(&callInfo{call: c.call, env: c.env, args: []Value{sub, rest, value}, names: []string{"", "", ""}})
		return s.subassign2(&callInfo{call: c.call, env: c.env, args: []Value{x, first, inner}, names: []string{"", "", ""}})
	}
	n := Length(x)
	nm := namesOf(x)
	_, xIsList := x.(*List)
	flat := wasNull || !xIsList
	if !xIsList && Length(value) == 0 && kindOf(value) != kOther {
		s.errorf(c.call, "replacement has length zero")
	}
	if !xIsList && Length(value) > 1 && kindOf(value) != kOther {
		s.errorf(c.call, "more elements supplied than there are to replace")
	}
	if Length(i) == 0 {
		s.errorf(c.call, "attempt to select less than one element in OneIndex")
	}
	firstSub()
	pos, newName := s.resolveIndex(c, i, n, nm)
	if lg, isLgl := i.(*Logical); isLgl && len(pos) == 0 && len(lg.V) == 1 {
		s.errorf(c.call, "attempt to select less than one element in integerOneIndex")
	}
	if flat && len(pos) == 1 && pos[0] == naPos {
		if n == 0 {
			s.errorf(c.call, "attempt to select less than one element in integerOneIndex")
		}
		s.errorf(c.call, "attempt to select more than one element in integerOneIndex")
	}
	if flat && len(pos) == 0 {
		s.errorf(c.call, "attempt to select less than one element in OneIndex <real>")
	}
	if len(pos) != 1 || pos[0] == naPos {
		s.errorf(c.call, "[[ ]] subscript out of bounds")
	}
	p := pos[0]
	if f := asDoublesOrNil(i); f != nil && f[0] <= -1 && n != 2 {
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
			res := &List{V: out, Attr: attrsOf(l), Expr: l.Expr}
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
	if kindOf(value) != kOther && lv != 1 {
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
	cx, cv := s.commonKindFlat(c, x, value)
	if lst, ok := value.(*List); ok {
		cv = &List{V: []Value{lst}}
	}
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
		s.assignTo(v, Intern(name), value, c.call)
		return v
	case nullT:
		if _, vnull := value.(nullT); vnull {
			return Null
		}
		x = &List{V: []Value{}}
	case *Logical, *Integer, *Double, *String, *Complex, *Raw:
		s.warningf(c.call, "Coercing LHS to a list")
		var keep []Attr
		for _, a := range attrsOf(x) {
			if a.Name == "names" {
				keep = append(keep, a)
			}
		}
		x = withAttr(coerceTo(x, kList), keep)
	case *List:
		if v.Expr {
			x = &List{V: v.V, Attr: v.Attr}
		}
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
