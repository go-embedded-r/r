package eval

import (
	"math"
	"strings"
)

// setAttr returns a copy of vector v with attribute name set to val (NULL
// removes it), keeping the order of the others.
func setAttr(v Value, name string, val Value) Value {
	old := attrsOf(v)
	out := make([]Attr, 0, len(old)+1)
	done := false
	for _, a := range old {
		if a.Name == name {
			if _, isNull := val.(nullT); !isNull {
				out = append(out, Attr{Name: name, Value: val})
			}
			done = true
			continue
		}
		out = append(out, a)
	}
	if !done {
		if _, isNull := val.(nullT); !isNull {
			out = append(out, Attr{Name: name, Value: val})
		}
	}
	if len(out) == 0 {
		out = nil
	}
	return withAttr(v, out)
}

func hasAttrs(v Value) bool { return len(attrsOf(v)) > 0 }

func isVector(v Value) bool {
	switch v.(type) {
	case *Logical, *Integer, *Double, *String, *List:
		return true
	}
	return false
}

// classAttr is the class attribute as strings, nil when there is none.
func classAttr(v Value) []string {
	c, ok := getAttr(v, "class").(*String)
	if !ok {
		return nil
	}
	out := make([]string, len(c.V))
	for i, e := range c.V {
		out[i] = e.S
	}
	return out
}

func dimsOf(v Value) []int {
	switch d := getAttr(v, "dim").(type) {
	case *Integer:
		out := make([]int, len(d.V))
		for i, e := range d.V {
			out[i] = int(e)
		}
		return out
	case *Double:
		out := make([]int, len(d.V))
		for i, e := range d.V {
			out[i] = int(e)
		}
		return out
	}
	return nil
}

// implicitClass is class(x) for an object without a class attribute.
func implicitClass(v Value) []string {
	if d := dimsOf(v); len(d) == 2 {
		return []string{"matrix", "array"}
	} else if len(d) > 0 {
		return []string{"array"}
	}
	switch x := v.(type) {
	case nullT:
		return []string{"NULL"}
	case *Logical:
		return []string{"logical"}
	case *Integer:
		return []string{"integer"}
	case *Double:
		return []string{"numeric"}
	case *String:
		return []string{"character"}
	case *List:
		return []string{"list"}
	case *Closure, *Builtin:
		return []string{"function"}
	case *Symbol:
		return []string{"name"}
	case *Env:
		return []string{"environment"}
	case *Lang:
		if s, ok := x.Fun.(*Symbol); ok {
			switch s.Name {
			case "if", "while", "for", "=", "<-", "(", "{":
				return []string{s.Name}
			}
		}
		return []string{"call"}
	}
	return []string{"unknown"}
}

// classOf is class(x).
func classOf(v Value) []string {
	if c := classAttr(v); c != nil {
		return c
	}
	return implicitClass(v)
}

// dispatchClass is the class vector S3 dispatch walks: the class attribute,
// or the implicit class with the type as R's dispatch sees it.
func dispatchClass(v Value) []string {
	if c := classAttr(v); c != nil {
		return c
	}
	out := implicitClass(v)
	switch v.(type) {
	case *Integer:
		out = append(out, "numeric")
	case *Double:
		out = []string{}
		if d := dimsOf(v); len(d) == 2 {
			out = append(out, "matrix", "array")
		} else if len(d) > 0 {
			out = append(out, "array")
		}
		out = append(out, "double", "numeric")
	}
	if _, isInt := v.(*Integer); isInt {
		if d := dimsOf(v); len(d) > 0 {
			out = append([]string{}, implicitClass(v)...)
			out = append(out, "integer", "numeric")
		}
	}
	return out
}

func strVec(ss []string) *String {
	out := &String{V: make([]Str, len(ss))}
	for i, s := range ss {
		out.V[i] = Str{S: s}
	}
	return out
}

func (s *Session) needVector(c *callInfo, x Value, what string) {
	if _, isNull := x.(nullT); isNull {
		s.errorf(c.call, "attempt to set an attribute on NULL")
	}
	if _, isEnv := x.(*Env); isEnv && what == "an attribute" {
		return
	}
	if !isVector(x) {
		s.errorf(c.call, "cannot set %s on this object", what)
	}
}

// setNames implements names(x) <- value.
func (s *Session) setNames(c *callInfo, x, value Value) Value {
	if _, xNull := x.(nullT); xNull {
		if _, vNull := value.(nullT); vNull {
			return Null
		}
	}
	s.needVector(c, x, "names")
	if _, isNull := value.(nullT); isNull {
		return setAttr(x, "names", Null)
	}
	var nm []Str
	switch kindOf(value) {
	case kLogical, kInteger, kDouble, kString:
		nm = append([]Str(nil), asStrings(value)...)
	default:
		if sym, ok := value.(*Symbol); ok {
			nm = []Str{{S: sym.Name}}
		} else {
			s.errorf(c.call, "invalid type (%s) for 'names': must be vector or NULL", typeName(value))
		}
	}
	n := Length(x)
	if len(nm) > n {
		s.errorf(c.call, "'names' attribute [%d] must be the same length as the vector [%d]", len(nm), n)
	}
	for len(nm) < n {
		nm = append(nm, Str{NA: true})
	}
	return setAttr(x, "names", &String{V: nm})
}

func (s *Session) setDim(c *callInfo, x, value Value) Value {
	s.needVector(c, x, "dim")
	if _, isNull := value.(nullT); isNull {
		y := setAttr(x, "dim", Null)
		return setAttr(y, "dimnames", Null)
	}
	var dims []int32
	switch kindOf(value) {
	case kInteger, kLogical:
		dims = asInts(value)
	case kDouble:
		for _, f := range asDoubles(value) {
			if math.IsNaN(f) {
				dims = append(dims, naInt)
			} else {
				dims = append(dims, int32(f))
			}
		}
	default:
		s.errorf(c.call, "invalid second argument, must be vector or NULL")
	}
	if len(dims) == 0 {
		s.errorf(c.call, "invalid second argument, must be vector or NULL")
	}
	prod := 1
	for _, d := range dims {
		if d == naInt {
			s.errorf(c.call, "the dims contain missing values")
		}
		if d < 0 {
			s.errorf(c.call, "the dims contain negative values")
		}
		prod *= int(d)
	}
	if prod != Length(x) {
		s.errorf(c.call, "dims [product %d] do not match the length of object [%d]", prod, Length(x))
	}
	y := setAttr(x, "names", Null)
	y = setAttr(y, "dimnames", Null)
	return setAttr(y, "dim", &Integer{V: dims})
}

func (s *Session) setClass(c *callInfo, x, value Value) Value {
	if _, isNull := x.(nullT); isNull {
		s.errorf(c.call, "attempt to set an attribute on NULL")
	}
	if _, isEnv := x.(*Env); !isVector(x) && !isEnv {
		s.errorf(c.call, "cannot set class on this object")
	}
	switch v := value.(type) {
	case nullT:
		return setAttr(x, "class", Null)
	case *Logical, *Integer, *Double:
		return s.setClass(c, x, &String{V: asStrings(v)})
	case *String:
		if len(v.V) == 0 {
			return setAttr(x, "class", Null)
		}
		return setAttr(x, "class", &String{V: append([]Str(nil), v.V...)})
	}
	s.errorf(c.call, "attempt to set invalid 'class' attribute")
	return nil
}

func (s *Session) setAttrByName(c *callInfo, x Value, which string, value Value) Value {
	switch which {
	case "names":
		return s.setNames(c, x, value)
	case "dim":
		return s.setDim(c, x, value)
	case "class":
		return s.setClass(c, x, value)
	}
	s.needVector(c, x, "an attribute")
	return setAttr(x, which, value)
}

func (s *Session) whichArg(c *callInfo, v Value) string {
	st, ok := v.(*String)
	if !ok || len(st.V) != 1 || st.V[0].NA {
		s.errorf(c.call, "'name' must be non-null character string")
	}
	return st.V[0].S
}

func init() {
	regBuiltin("attr", func(s *Session, c *callInfo) Value {
		if len(c.args) < 2 || len(c.args) > 3 {
			s.errorf(c.call, "either 2 or 3 arguments are required")
		}
		which := s.whichArg(c, c.args[1])
		exact := false
		if len(c.args) == 3 {
			exact = len(logicalOf(s, c, c.args[2], "invalid 'exact' value")) > 0 && logicalOf(s, c, c.args[2], "")[0] == 1
		}
		x := c.args[0]
		if which == "names" {
			if n := namesOf(x); n != nil {
				return n
			}
			return Null
		}
		var found Value
		hits := 0
		for _, a := range attrsOf(x) {
			if a.Name == which {
				return a.Value
			}
			if !exact && strings.HasPrefix(a.Name, which) {
				found = a.Value
				hits++
			}
		}
		if hits == 1 {
			return found
		}
		return Null
	})
	regBuiltin("attr<-", func(s *Session, c *callInfo) Value {
		if len(c.args) != 3 {
			s.errorf(c.call, "3 arguments are required")
		}
		return s.setAttrByName(c, c.args[0], s.whichArg(c, c.args[1]), c.args[2])
	})
	regBuiltin("attributes", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "attributes")
		at := attrsOf(x)
		if len(at) == 0 {
			return Null
		}
		out := &List{}
		var names []Str
		// names come first
		for _, a := range at {
			if a.Name == "names" {
				out.V = append(out.V, a.Value)
				names = append(names, Str{S: a.Name})
			}
		}
		for _, a := range at {
			if a.Name != "names" {
				out.V = append(out.V, a.Value)
				names = append(names, Str{S: a.Name})
			}
		}
		out.Attr = namesAttr(names)
		return out
	})
	regBuiltin("attributes<-", func(s *Session, c *callInfo) Value {
		if len(c.args) != 2 {
			s.errorf(c.call, "2 arguments are required")
		}
		x := c.args[0]
		if _, isNull := x.(nullT); isNull {
			s.errorf(c.call, "attempt to set an attribute on NULL")
		}
		y := withAttr(x, nil)
		switch v := c.args[1].(type) {
		case nullT:
			return y
		case *List:
			nm := namesOf(v)
			if nm == nil && len(v.V) > 0 {
				s.errorf(c.call, "attributes must be named")
			}
			// dim first, as R does, so that dimnames can be checked
			for i, e := range v.V {
				if nm.V[i].S == "dim" {
					y = s.setAttrByName(c, y, "dim", e)
				}
			}
			for i, e := range v.V {
				if nm.V[i].S != "dim" {
					y = s.setAttrByName(c, y, nm.V[i].S, e)
				}
			}
			return y
		}
		s.errorf(c.call, "attributes must be a list or NULL")
		return nil
	})
	regBuiltin("names<-", func(s *Session, c *callInfo) Value {
		if len(c.args) != 2 {
			s.errorf(c.call, "2 arguments are required")
		}
		return s.setNames(c, c.args[0], c.args[1])
	})
	regBuiltin("class", func(s *Session, c *callInfo) Value { return strVec(classOf(s.nargs1(c, "class"))) })
	regBuiltin("oldClass", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "oldClass")
		if a := getAttr(x, "class"); a != nil {
			return a
		}
		return Null
	})
	regBuiltin("class<-", func(s *Session, c *callInfo) Value {
		if len(c.args) != 2 {
			s.errorf(c.call, "2 arguments are required")
		}
		return s.setClass(c, c.args[0], c.args[1])
	})
	regBuiltin("oldClass<-", func(s *Session, c *callInfo) Value {
		if len(c.args) != 2 {
			s.errorf(c.call, "2 arguments are required")
		}
		return s.setClass(c, c.args[0], c.args[1])
	})
	regBuiltin("unclass", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "unclass")
		if getAttr(x, "class") == nil {
			return x
		}
		return setAttr(x, "class", Null)
	})
	regBuiltin("dim", func(s *Session, c *callInfo) Value {
		if d := getAttr(s.nargs1(c, "dim"), "dim"); d != nil {
			return d
		}
		return Null
	})
	regBuiltin("dim<-", func(s *Session, c *callInfo) Value {
		if len(c.args) != 2 {
			s.errorf(c.call, "2 arguments are required")
		}
		return s.setDim(c, c.args[0], c.args[1])
	})
	regBuiltin("levels<-", func(s *Session, c *callInfo) Value {
		if len(c.args) != 2 {
			s.errorf(c.call, "2 arguments are required")
		}
		s.needVector(c, c.args[0], "levels")
		return setAttr(c.args[0], "levels", c.args[1])
	})
	regBuiltin("dimnames<-", func(s *Session, c *callInfo) Value {
		if len(c.args) != 2 {
			s.errorf(c.call, "2 arguments are required")
		}
		x := c.args[0]
		dims := dimsOf(x)
		switch v := c.args[1].(type) {
		case nullT:
			return setAttr(x, "dimnames", Null)
		case *List:
			if len(dims) == 0 {
				s.errorf(c.call, "'dimnames' applied to non-array")
			}
			if len(v.V) != len(dims) {
				s.errorf(c.call, "length of 'dimnames' [%d] must match that of 'dims' [%d]", len(v.V), len(dims))
			}
			out := &List{V: make([]Value, len(v.V)), Attr: attrsOf(v)}
			for i, e := range v.V {
				switch ev := e.(type) {
				case nullT:
					out.V[i] = Null
				case *Logical, *Integer, *Double, *String:
					if Length(ev) == 0 {
						out.V[i] = Null
						continue
					}
					if Length(ev) != dims[i] {
						s.errorf(c.call, "length of 'dimnames' [%d] not equal to array extent", i+1)
					}
					out.V[i] = &String{V: asStrings(ev)}
				default:
					s.errorf(c.call, "invalid type (%s) for 'dimnames' (must be a vector)", typeName(e))
				}
			}
			return setAttr(x, "dimnames", out)
		}
		s.errorf(c.call, "'dimnames' must be a list")
		return nil
	})
	regBuiltin("dimnames", func(s *Session, c *callInfo) Value {
		if d := getAttr(s.nargs1(c, "dimnames"), "dimnames"); d != nil {
			return d
		}
		return Null
	})
}
