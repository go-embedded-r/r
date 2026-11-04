package eval

import (
	"strings"
)

// These are functions of R's base package that are closures there (R code
// over .Internal calls); they are provided natively until base R is
// rewritten (M3), with the behaviour the oracle shows.

func (s *Session) modeKind(c *callInfo, mode string) int {
	switch mode {
	case "logical":
		return kLogical
	case "integer":
		return kInteger
	case "numeric", "double":
		return kDouble
	case "character":
		return kString
	case "list":
		return kList
	}
	s.errorf(c.call, "vector: cannot make a vector of mode '%s'.", mode)
	return 0
}

func (s *Session) makeVector(c *callInfo, kind int, n Value) Value {
	var length int
	switch kindOf(n) {
	case kLogical, kInteger, kDouble:
		d := asDoubles(n)
		if len(d) != 1 || d[0] != d[0] || d[0] < 0 {
			s.errorf(c.call, "invalid 'length' argument")
		}
		length = int(min(d[0], float64(s.opts.MaxVector)+1))
	default:
		s.errorf(c.call, "invalid 'length' argument")
	}
	s.checkLen(c.call, length)
	switch kind {
	case kLogical:
		return &Logical{V: make([]int32, length)}
	case kInteger:
		return &Integer{V: make([]int32, length)}
	case kDouble:
		return &Double{V: make([]float64, length)}
	case kString:
		return &String{V: make([]Str, length)}
	}
	out := &List{V: make([]Value, length)}
	for i := range out.V {
		out.V[i] = Null
	}
	return out
}

func identical(a, b Value) bool {
	if kindOf(a) != kindOf(b) || Length(a) != Length(b) {
		return false
	}
	if !identicalAttrs(attrsOf(a), attrsOf(b)) {
		return false
	}
	switch x := a.(type) {
	case nullT:
		return true
	case *Logical:
		return equalInt32(x.V, b.(*Logical).V)
	case *Integer:
		return equalInt32(x.V, b.(*Integer).V)
	case *Double:
		y := b.(*Double).V
		for i, e := range x.V {
			switch {
			case isNAReal(e) != isNAReal(y[i]):
				return false
			case e != e && y[i] != y[i]:
			case e != y[i]:
				return false
			}
		}
		return true
	case *String:
		y := b.(*String).V
		for i, e := range x.V {
			if e != y[i] {
				return false
			}
		}
		return true
	case *List:
		y := b.(*List).V
		for i, e := range x.V {
			if !identical(e, y[i]) {
				return false
			}
		}
		return true
	case *Lang:
		return callText(x) == callText(b.(*Lang)) || x == b
	}
	return a == b
}

func equalInt32(a, b []int32) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func identicalAttrs(a, b []Attr) bool {
	if len(a) != len(b) {
		return false
	}
	for _, x := range a {
		found := false
		for _, y := range b {
			if x.Name == y.Name {
				found = identical(x.Value, y.Value)
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (s *Session) paste(c *callInfo, defSep string) Value {
	sep := defSep
	var collapse *string
	var cols [][]Str
	for i, a := range c.args {
		switch c.names[i] {
		case "sep":
			if st, ok := a.(*String); ok && len(st.V) == 1 {
				sep = st.V[0].S
			} else {
				s.errorf(c.call, "invalid separator")
			}
			continue
		case "collapse":
			if st, ok := a.(*String); ok && len(st.V) == 1 {
				v := st.V[0].S
				collapse = &v
			} else if _, isNull := a.(nullT); !isNull {
				s.errorf(c.call, "invalid 'collapse' argument")
			}
			continue
		}
		var strs []Str
		switch x := a.(type) {
		case *Symbol:
			strs = []Str{{S: x.Name}}
		case *List:
			for _, e := range x.V {
				if Length(e) == 1 && isVector(e) && kindOf(e) != kList {
					strs = append(strs, asStrings(e)...)
				} else {
					strs = append(strs, Str{S: callText2(e)})
				}
			}
		default:
			strs = asStrings(a)
		}
		cols = append(cols, strs)
	}
	n := 0
	for _, col := range cols {
		if len(col) > n {
			n = len(col)
		}
	}
	out := make([]Str, n)
	for i := 0; i < n; i++ {
		var parts []string
		for _, col := range cols {
			if len(col) == 0 {
				parts = append(parts, "") // a zero-length argument is ""
				continue
			}
			e := col[i%len(col)]
			if e.NA {
				parts = append(parts, "NA")
			} else {
				parts = append(parts, e.S)
			}
		}
		out[i] = Str{S: strings.Join(parts, sep)}
	}
	if collapse != nil {
		parts := make([]string, len(out))
		for i, e := range out {
			parts[i] = e.S
		}
		return mkString(strings.Join(parts, *collapse))
	}
	return &String{V: out}
}

func callText2(v Value) string {
	if l, ok := v.(*Lang); ok {
		return callText(l)
	}
	return typeName(v)
}

func init() {
	regBuiltin("structure", func(s *Session, c *callInfo) Value {
		if len(c.args) == 0 {
			s.errorf(c.call, "argument \".Data\" is missing, with no default")
		}
		di := 0
		for i := range c.args {
			if c.names[i] == ".Data" || c.names[i] == "" {
				di = i
				break
			}
		}
		x := c.args[di]
		if _, isNull := x.(nullT); isNull && len(c.args) > 1 {
			s.errorf(c.call, "attempt to set an attribute on NULL")
		}
		for i := 0; i < len(c.args); i++ {
			if i == di {
				continue
			}
			name := c.names[i]
			if name == "" {
				s.errorf(c.call, "attributes must be named")
			}
			if name == ".Names" {
				name = "names"
			}
			x = s.setAttrByName(c, x, name, c.args[i])
		}
		return x
	})
	regBuiltin("inherits", func(s *Session, c *callInfo) Value {
		if len(c.args) < 2 {
			s.errorf(c.call, "argument \"what\" is missing, with no default")
		}
		what, ok := c.args[1].(*String)
		if !ok {
			s.errorf(c.call, "'what' must be a character vector or an object with a nameOfClass() method")
		}
		which := s.flagArg(c, "which", false)
		cls := classOf(c.args[0])
		if which {
			out := make([]int32, len(what.V))
			for i, w := range what.V {
				for j, cl := range cls {
					if cl == w.S {
						out[i] = int32(j + 1)
						break
					}
				}
			}
			return &Integer{V: out}
		}
		for _, w := range what.V {
			for _, cl := range cls {
				if cl == w.S {
					return mkLogical(true)
				}
			}
		}
		return mkLogical(false)
	})
	for name, kind := range map[string]int{"numeric": kDouble, "double": kDouble, "integer": kInteger, "logical": kLogical, "character": kString} {
		kind := kind
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			if len(c.args) == 0 {
				return s.makeVector(c, kind, mkInt(0))
			}
			return s.makeVector(c, kind, c.args[0])
		})
	}
	regBuiltin("vector", func(s *Session, c *callInfo) Value {
		mode, n := "logical", Value(mkInt(0))
		pos := 0
		for i, a := range c.args {
			switch {
			case c.names[i] == "mode" || (c.names[i] == "" && pos == 0):
				if st, ok := a.(*String); ok && len(st.V) == 1 {
					mode = st.V[0].S
				}
				if c.names[i] == "" {
					pos++
				}
			case c.names[i] == "length" || (c.names[i] == "" && pos == 1):
				n = a
				if c.names[i] == "" {
					pos++
				}
			}
		}
		return s.makeVector(c, s.modeKind(c, mode), n)
	})
	regBuiltin("identical", func(s *Session, c *callInfo) Value {
		if len(c.args) < 2 {
			s.errorf(c.call, "argument \"y\" is missing, with no default")
		}
		return mkLogical(identical(c.args[0], c.args[1]))
	})
	regBuiltin("paste", func(s *Session, c *callInfo) Value { return s.paste(c, " ") })
	regBuiltin("paste0", func(s *Session, c *callInfo) Value { return s.paste(c, "") })
}
