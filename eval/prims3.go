package eval

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/go-embedded-r/r/syntax"
)

// typeof, deparse, as.name, unique, match, factor and friends: the base
// functions every other base function leans on.

func init() {
	regBuiltin("typeof", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "typeof")
		if b, ok := x.(*Builtin); ok {
			if b.Special {
				return mkString("special")
			}
			return mkString("builtin")
		}
		return mkString(typeName(x))
	})
	regBuiltin("mode", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "mode")
		switch v := x.(type) {
		case *Integer, *Double:
			return mkString("numeric")
		case *Closure, *Builtin:
			return mkString("function")
		case *Symbol:
			return mkString("name")
		case *Lang:
			if f, ok := v.Fun.(*Symbol); ok && f.Name == "(" {
				return mkString("(")
			}
			return mkString("call")
		}
		return mkString(typeName(x))
	})
	regBuiltin("storage.mode", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "storage.mode")
		switch x.(type) {
		case *Closure, *Builtin:
			return mkString("function")
		case *Symbol:
			return mkString("symbol")
		}
		return mkString(typeName(x))
	})
	asName := func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "as.name")
		switch v := x.(type) {
		case *Symbol:
			return v
		case *String:
			if len(v.V) >= 1 && !v.V[0].NA && v.V[0].S != "" {
				return Intern(v.V[0].S)
			}
			if len(v.V) == 0 {
				s.errorf(c.call, "invalid data of mode 'character' (too short)")
			}
			s.errorf(c.call, "attempt to use zero-length variable name")
		case *Logical, *Integer, *Double:
			if Length(v) >= 1 {
				return Intern(asStrings(v)[0].S)
			}
			s.errorf(c.call, "invalid type/length (symbol/0) in vector allocation")
		}
		s.errorf(c.call, "cannot coerce type '%s' to vector of type 'symbol'", typeName(x))
		return nil
	}
	regBuiltin("as.name", asName)
	regBuiltin("as.symbol", asName)
	regBuiltin("deparse", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "expr", "width.cutoff", "backtick", "control", "nlines")
		if a[0] == nil {
			s.errorf(c.call, "argument \"expr\" is missing, with no default")
		}
		var node syntax.Node
		if cl, ok := a[0].(*Closure); ok {
			fn, ok := toSyntax(&Lang{Fun: Intern("function"), Args: []Arg{{Value: &Formals{F: cl.Formals}}, {Value: cl.Body}}}).(*syntax.Function)
			if ok {
				return strVec(syntax.DeparseClosure(fn.Params, fn.Body))
			}
		}
		node = toSyntax(a[0])
		return strVec(syntax.Deparse(node, nil))
	})

	regBuiltin("unique", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "incomparables", "fromLast")
		x := a[0]
		if x == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		if _, isNull := x.(nullT); isNull {
			return Null
		}
		if !isVector(x) {
			s.errorf(c.call, "unique() applies only to vectors")
		}
		dup := s.duplicatedFlags(c, x, a[2] != nil && s.truthy(c, a[2], "fromLast"))
		var pos []int
		for i, d := range dup {
			if !d {
				pos = append(pos, i)
			}
		}
		if pos == nil {
			pos = []int{}
		}
		res := s.subsetPositions(x, pos)
		if nm := namesOf(res); nm != nil {
			var keep []Attr
			for _, at := range attrsOf(res) {
				if at.Name != "names" {
					keep = append(keep, at)
				}
			}
			return withAttr(res, keep)
		}
		return res
	})
	regBuiltin("duplicated", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "incomparables", "fromLast")
		x := a[0]
		if x == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		if _, isNull := x.(nullT); isNull {
			return &Logical{V: []int32{}}
		}
		dup := s.duplicatedFlags(c, x, a[2] != nil && s.truthy(c, a[2], "fromLast"))
		out := make([]int32, len(dup))
		for i, d := range dup {
			out[i] = b2i(d)
		}
		return &Logical{V: out}
	})
	regBuiltin("anyDuplicated", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "incomparables", "fromLast")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		if _, isNull := a[0].(nullT); isNull {
			return mkInt(0)
		}
		for i, d := range s.duplicatedFlags(c, a[0], a[2] != nil && s.truthy(c, a[2], "fromLast")) {
			if d {
				return mkInt(int32(i + 1))
			}
		}
		return mkInt(0)
	})
	regBuiltin("match", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "table", "nomatch", "incomparables")
		if a[0] == nil || a[1] == nil {
			s.errorf(c.call, "argument \"%s\" is missing, with no default", map[bool]string{true: "x", false: "table"}[a[0] == nil])
		}
		nomatch := int32(naInt)
		if a[2] != nil {
			nm := asDoublesOrNil(a[2])
			if len(nm) < 1 {
				s.errorf(c.call, "invalid 'nomatch' argument")
			}
			nomatch = naInt
			if !math.IsNaN(nm[0]) {
				nomatch = int32(nm[0])
			}
		}
		res := s.matchValues(c, a[0], a[1], nomatch)
		return &Integer{V: res}
	})
	regBuiltin("%in%", func(s *Session, c *callInfo) Value {
		if len(c.args) != 2 {
			s.errorf(c.call, "operator needs two arguments")
		}
		res := s.matchValues(c, c.args[0], c.args[1], 0)
		out := make([]int32, len(res))
		for i, r := range res {
			out[i] = b2i(r > 0)
		}
		return &Logical{V: out}
	})

	regBuiltin("levels", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "levels")
		if l := getAttr(x, "levels"); l != nil {
			return l
		}
		return Null
	})
	regBuiltin("nlevels", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "nlevels")
		return mkInt(int32(Length(orNull(getAttr(x, "levels")))))
	})
	regBuiltin("factor", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "levels", "labels", "exclude", "ordered", "nmax")
		var x Value = &String{V: []Str{}}
		if a[0] != nil {
			x = a[0]
		}
		if _, isNull := x.(nullT); isNull {
			x = &String{V: []Str{}}
		}
		if !isVector(x) {
			s.errorf(c.call, "'x' must be atomic")
		}
		if _, isList := x.(*List); isList {
			s.errorf(c.call, "'x' must be atomic")
		}
		ordered := a[4] != nil && s.truthy(c, a[4], "ordered")
		if a[4] == nil && inheritsClass(x, "ordered") {
			ordered = true
		}
		xs := s.factorStrings(x)
		var levels []Str
		if a[1] != nil {
			levels = s.factorStrings(a[1])
		} else {
			// sort(unique(x)) with NA removed
			seen := map[string]bool{}
			var u []Str
			var pos []int
			for i, e := range xs {
				if e.NA || seen[e.S] {
					continue
				}
				seen[e.S] = true
				u = append(u, e)
				pos = append(pos, i)
			}
			if _, isStr := x.(*String); isStr || (isObject(x) && inheritsClass(x, "factor")) {
				sort.SliceStable(u, func(i, j int) bool { return collate(u[i].S, u[j].S) < 0 })
			} else {
				d := asDoubles(x)
				idx := make([]int, len(u))
				for i := range idx {
					idx[i] = i
				}
				sort.SliceStable(idx, func(i, j int) bool { return d[pos[idx[i]]] < d[pos[idx[j]]] })
				u2 := make([]Str, len(u))
				for i, k := range idx {
					u2[i] = u[k]
				}
				u = u2
			}
			levels = u
			if _, keep := a[3].(nullT); keep {
				for _, e := range xs {
					if e.NA {
						levels = append(levels, Str{NA: true})
						break
					}
				}
			}
		}
		// exclude: NA by default
		if a[3] != nil {
			if ex, ok := a[3].(*String); ok {
				drop := map[string]bool{}
				for _, e := range ex.V {
					drop[e.S] = true
				}
				var keep []Str
				for _, l := range levels {
					if !drop[l.S] {
						keep = append(keep, l)
					}
				}
				levels = keep
			}
		}
		var cleaned []Str
		_, keepNA := a[3].(nullT)
		for _, l := range levels {
			if !l.NA || keepNA {
				cleaned = append(cleaned, l)
			}
		}
		levels = cleaned
		index := map[string]int32{}
		naIdx := int32(0)
		for i, l := range levels {
			if l.NA {
				if naIdx == 0 {
					naIdx = int32(i + 1)
				}
				continue
			}
			if _, dup := index[l.S]; !dup {
				index[l.S] = int32(i + 1)
			}
		}
		codes := make([]int32, len(xs))
		for i, e := range xs {
			if e.NA {
				codes[i] = naInt
				if naIdx > 0 {
					codes[i] = naIdx
				}
			} else if k, ok := index[e.S]; ok {
				codes[i] = k
			} else {
				codes[i] = naInt
			}
		}
		labels := levels
		if a[2] != nil {
			ls := s.factorStrings(a[2])
			switch {
			case len(ls) == len(levels):
				labels = ls
			case len(ls) == 1:
				labels = make([]Str, len(levels))
				for i := range labels {
					labels[i] = Str{S: ls[0].S + itoa(i+1)}
				}
			default:
				s.errorf(c.call, "invalid 'labels'; length %d should be 1 or %d", len(ls), len(levels))
			}
			// duplicated labels merge their levels
			first := map[string]int32{}
			remap := make([]int32, len(labels))
			var uniq []Str
			for i, l := range labels {
				if k, ok := first[l.S]; ok {
					remap[i] = k
				} else {
					uniq = append(uniq, l)
					first[l.S] = int32(len(uniq))
					remap[i] = int32(len(uniq))
				}
			}
			if len(uniq) != len(labels) {
				for i, cd := range codes {
					if cd != naInt {
						codes[i] = remap[cd-1]
					}
				}
				labels = uniq
			}
		}
		res := &Integer{V: codes}
		var attr []Attr
		if nm := namesOf(x); nm != nil {
			attr = append(attr, Attr{Name: "names", Value: nm})
		}
		attr = append(attr, Attr{Name: "levels", Value: &String{V: append([]Str(nil), labels...)}})
		cls := []string{"factor"}
		if ordered {
			cls = []string{"ordered", "factor"}
		}
		attr = append(attr, Attr{Name: "class", Value: strVec(cls)})
		res.Attr = attr
		return res
	})
	regBuiltin("as.factor", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "as.factor")
		if inheritsClass(x, "factor") {
			return x
		}
		return builtins["factor"].fn(s, &callInfo{call: c.call, env: c.env, args: []Value{x}, names: []string{""}})
	})
}

func orNull(v Value) Value {
	if v == nil {
		return Null
	}
	return v
}

func inheritsClass(v Value, cl string) bool {
	return hasString(classAttr(v), cl)
}

// factorStrings is as.character for the arguments of factor(): a factor
// gives its labels.
func (s *Session) factorStrings(x Value) []Str {
	if inheritsClass(x, "factor") {
		lev, _ := getAttr(x, "levels").(*String)
		codes := asInts(x)
		out := make([]Str, len(codes))
		for i, cd := range codes {
			if cd == naInt || lev == nil || int(cd) > len(lev.V) || cd < 1 {
				out[i] = Str{NA: true}
			} else {
				out[i] = lev.V[cd-1]
			}
		}
		return out
	}
	return asStrings(x)
}

func uniqueKeepAttrs(x, res Value) []Attr {
	attr := attrsOf(res)
	if inheritsClass(x, "factor") {
		return attr
	}
	return attr
}

// keyOfElem is a comparable key for element i of an atomic vector.
func keyOfElem(v Value, i int) any {
	switch x := v.(type) {
	case *Logical:
		return int64(x.V[i])
	case *Integer:
		return int64(x.V[i])
	case *Double:
		f := x.V[i]
		switch {
		case isNAReal(f):
			return "NA_real_"
		case math.IsNaN(f):
			return "NaN"
		case f == 0:
			return float64(0)
		}
		return f
	case *String:
		if x.V[i].NA {
			return nil
		}
		return x.V[i].S
	case *Complex:
		z := x.V[i]
		if isNAComplex(z) {
			return "NA_complex_"
		}
		re, im := real(z), imag(z)
		if re == 0 {
			re = 0
		}
		if im == 0 {
			im = 0
		}
		if math.IsNaN(re) || math.IsNaN(im) {
			return "NaN_complex_"
		}
		return complex(re, im)
	case *Raw:
		return int64(x.V[i])
	}
	return nil
}

func (s *Session) duplicatedFlags(c *callInfo, x Value, fromLast bool) []bool {
	n := Length(x)
	dup := make([]bool, n)
	if l, ok := x.(*List); ok {
		order := make([]int, n)
		for i := range order {
			if fromLast {
				order[i] = n - 1 - i
			} else {
				order[i] = i
			}
		}
		// elements are told apart by identical(), within buckets of a cheap
		// key, so that a long list is not compared pairwise
		seen := map[string][]Value{}
		for _, i := range order {
			k := bucketKey(l.V[i])
			found := false
			for _, e := range seen[k] {
				if identical(e, l.V[i]) {
					found = true
					break
				}
			}
			if found {
				dup[i] = true
			} else {
				seen[k] = append(seen[k], l.V[i])
			}
		}
		return dup
	}
	seen := map[any]bool{}
	for k := 0; k < n; k++ {
		i := k
		if fromLast {
			i = n - 1 - k
		}
		key := keyOfElem(x, i)
		if seen[key] {
			dup[i] = true
		}
		seen[key] = true
	}
	return dup
}

// matchValues is match(x, table, nomatch): the common type of the two
// arguments decides how elements are compared.
func (s *Session) matchValues(c *callInfo, x, table Value, nomatch int32) []int32 {
	for _, v := range []Value{x, table} {
		switch v.(type) {
		case nullT, *Logical, *Integer, *Double, *String, *Complex, *Raw:
		case *List:
		default:
			s.errorf(c.call, "'match' requires vector arguments")
		}
	}
	if inheritsClass(x, "factor") {
		x = &String{V: s.factorStrings(x)}
	}
	if inheritsClass(table, "factor") {
		table = &String{V: s.factorStrings(table)}
	}
	nx := Length(x)
	out := make([]int32, nx)
	_, xl := x.(*List)
	_, tl := table.(*List)
	if xl || tl {
		xe, te := s.asListElems(c, x), s.asListElems(c, table)
		for i := range out {
			out[i] = nomatch
			for j, t := range te {
				if identical(xe[i], t) {
					out[i] = int32(j + 1)
					break
				}
			}
		}
		return out
	}
	if r, ok := x.(*Raw); ok {
		x = rawOperand(r, table)
	}
	if r, ok := table.(*Raw); ok {
		table = rawOperand(r, x)
	}
	kx, kt := kindOf(x), kindOf(table)
	k := max(kx, kt)
	if kx == kNull || kt == kNull {
		for i := range out {
			out[i] = nomatch
		}
		return out
	}
	cx, ct := coerceTo(x, k), coerceTo(table, k)
	idx := map[any]int32{}
	for j := Length(ct) - 1; j >= 0; j-- {
		idx[keyOfElem(ct, j)] = int32(j + 1)
	}
	for i := range out {
		if j, ok := idx[keyOfElem(cx, i)]; ok {
			out[i] = j
		} else {
			out[i] = nomatch
		}
	}
	return out
}

var _ = strings.Compare

// bucketKey is a key that identical values share: the type, the length, and
// the first few elements of an atomic vector. Values with the same key are
// then compared with identical().
func bucketKey(v Value) string {
	n := Length(v)
	key := typeName(v) + ":" + strconv.Itoa(n)
	switch x := v.(type) {
	case *Logical, *Integer, *Double, *String, *Complex, *Raw:
		for i := 0; i < n && i < 4; i++ {
			key += fmt.Sprintf("|%v", keyOfElem(x, i))
		}
	}
	return key
}
