package eval

import (
	"math"
	"strconv"
	"strings"
)

// The apply family, do.call, and the vector helpers they need (unlist, rep,
// rev, which, seq).

// callWith applies fn to already evaluated arguments (named or not).
func (s *Session) callWith(fn Value, call *Lang, args []Arg, env *Env) Value {
	switch f := fn.(type) {
	case *Closure:
		return s.applyClosure(f, call, args, env, nil)
	case *Builtin:
		ci := &callInfo{call: call, env: env}
		if f.Special {
			ci.raw = args
		} else {
			for _, a := range args {
				ci.args = append(ci.args, a.Value)
				ci.names = append(ci.names, a.Name)
			}
		}
		return f.fn(s, ci)
	}
	s.errorf(call, "attempt to apply non-function")
	return nil
}

// matchFun is match.fun: a function, or the name of one.
func (s *Session) matchFun(c *callInfo, v Value, env *Env) Value {
	switch f := v.(type) {
	case *Closure, *Builtin:
		return v
	case *String:
		if len(f.V) == 1 && !f.V[0].NA {
			return s.findFun(Intern(f.V[0].S), env, c.call)
		}
	case *Symbol:
		return s.findFun(f, env, c.call)
	}
	s.errorf(c.call, "'%s' is not a function, character or symbol", strings.Join(syntaxDeparse(v), ""))
	return nil
}

// asListElems splits X into its elements, as as.list does for lapply.
func (s *Session) asListElems(c *callInfo, x Value) []Value {
	switch v := x.(type) {
	case *List:
		return v.V
	case nullT:
		return nil
	case *Env:
		var out []Value
		for _, n := range v.Names() {
			val, _ := v.Get(n)
			out = append(out, val)
		}
		return out
	case *Logical, *Integer, *Double, *String:
		n := Length(v)
		out := make([]Value, n)
		for i := range out {
			e := takeElems(v, []int{i})
			if f := classAttr(v); len(f) == 1 && f[0] == "factor" {
				e = withAttr(e, nil)
			}
			out[i] = e
		}
		return out
	}
	s.errorf(c.call, "argument is not a vector")
	return nil
}

var (
	funXi = &Lang{Fun: Intern("FUN"), Args: []Arg{
		{Value: &Lang{Fun: Intern("[["), Args: []Arg{{Value: Intern("X")}, {Value: Intern("i")}}}},
		{Value: Intern("...")},
	}}
)

// extraArgs are the arguments after the named formals of the apply functions.
func (s *Session) splitArgs(c *callInfo, formals []string) (matched []Value, extra []Arg) {
	matched = make([]Value, len(formals))
	used := make([]bool, len(c.args))
	for i, n := range c.names {
		if n == "" {
			continue
		}
		for j, f := range formals {
			if f == n && matched[j] == nil {
				matched[j], used[i] = c.args[i], true
				break
			}
		}
	}
	j := 0
	for i := range c.args {
		if used[i] || c.names[i] != "" {
			continue
		}
		for j < len(formals) && matched[j] != nil {
			j++
		}
		if j >= len(formals) {
			break
		}
		matched[j], used[i] = c.args[i], true
	}
	for i := range c.args {
		if !used[i] {
			extra = append(extra, Arg{Name: c.names[i], HasName: c.names[i] != "", Value: c.args[i]})
		}
	}
	return
}

func (s *Session) lapplyCore(c *callInfo, x Value, fn Value, extra []Arg) ([]Value, *String) {
	elems := s.asListElems(c, x)
	out := make([]Value, len(elems))
	for i, e := range elems {
		args := append([]Arg{{Value: e}}, extra...)
		out[i] = s.callWith(fn, funXi, args, c.env)
	}
	return out, namesOf(x)
}

func listOf(vals []Value, names *String) *List {
	l := &List{V: vals}
	if l.V == nil {
		l.V = []Value{}
	}
	if names != nil {
		if len(vals) == 0 {
			l.Attr = namesAttr(names.V)
		} else {
			l.Attr = namesAttr(fitNames(names.V, len(vals)))
		}
	}
	return l
}

// fitNames pads names with NA, or cuts them, to n, as names<- does.
func fitNames(nm []Str, n int) []Str {
	out := make([]Str, n)
	for i := range out {
		if i < len(nm) {
			out[i] = nm[i]
		} else {
			out[i] = Str{NA: true}
		}
	}
	return out
}

// simplify2array for the cases sapply and mapply produce: all results of one
// length become a vector (1) or a matrix (>1).
func (s *Session) simplify(c *callInfo, vals []Value, names *String) Value {
	if len(vals) == 0 {
		return &List{V: []Value{}}
	}
	n := Length(vals[0])
	for _, v := range vals {
		if Length(v) != n || !isVector(v) {
			return listOf(vals, names)
		}
		if _, isNull := v.(nullT); isNull {
			return listOf(vals, names)
		}
	}
	if n == 1 {
		for _, v := range vals {
			if _, isList := v.(*List); isList {
				// unlist(recursive = FALSE) of length-one lists
				out := make([]Value, len(vals))
				for i, v := range vals {
					out[i] = v.(*List).V[0]
				}
				return listOf(out, names)
			}
		}
		cc := &callInfo{call: c.call, env: c.env, args: vals, names: make([]string, len(vals))}
		inner := false
		for _, v := range vals {
			if namesOf(v) != nil {
				inner = true
			}
		}
		if names != nil && !inner {
			for i := range vals {
				vals[i] = withAttr(vals[i], nil)
			}
			res := s.combine(cc)
			return withAttr(res, namesAttr(fitNames(names.V, len(vals))))
		}
		if names != nil {
			fn := fitNames(names.V, len(vals))
			for i := range cc.names {
				if fn[i].NA {
					cc.names[i] = "NA"
				} else {
					cc.names[i] = fn[i].S
				}
			}
		}
		return s.combine(cc)
	}
	if n == 0 {
		return listOf(vals, names)
	}
	cc := &callInfo{call: c.call, env: c.env, args: make([]Value, len(vals)), names: make([]string, len(vals))}
	for i, v := range vals {
		cc.args[i] = withAttr(v, nil)
	}
	res := s.combine(cc)
	res = withAttr(res, nil)
	res = setAttr(res, "dim", &Integer{V: []int32{int32(n), int32(len(vals))}})
	var rown Value = Null
	if rn := namesOf(vals[0]); rn != nil {
		rown = rn
	}
	var coln Value = Null
	if names != nil {
		coln = &String{V: fitNames(names.V, len(vals))}
	}
	if rown != Null || coln != Null {
		res = setAttr(res, "dimnames", &List{V: []Value{rown, coln}})
	}
	return res
}

func (s *Session) useNames(x Value, names *String, use bool) *String {
	if names != nil {
		return names
	}
	if !use {
		return nil
	}
	if st, ok := x.(*String); ok {
		return &String{V: append([]Str(nil), st.V...)}
	}
	return nil
}

func init() {
	regBuiltin("lapply", func(s *Session, c *callInfo) Value {
		m, extra := s.splitArgs(c, []string{"X", "FUN"})
		if m[0] == nil || m[1] == nil {
			s.errorf(c.call, "argument \"%s\" is missing, with no default", map[bool]string{true: "X", false: "FUN"}[m[0] == nil])
		}
		fn := s.matchFun(c, m[1], c.env)
		vals, names := s.lapplyCore(c, m[0], fn, extra)
		return listOf(vals, names)
	})
	regBuiltin("sapply", func(s *Session, c *callInfo) Value {
		m, extra := s.splitArgs(c, []string{"X", "FUN"})
		if m[0] == nil || m[1] == nil {
			s.errorf(c.call, "argument \"%s\" is missing, with no default", map[bool]string{true: "X", false: "FUN"}[m[0] == nil])
		}
		simp, useN := true, true
		var rest []Arg
		for _, a := range extra {
			switch a.Name {
			case "simplify":
				simp = s.truthy(c, a.Value, "simplify")
			case "USE.NAMES":
				useN = s.truthy(c, a.Value, "USE.NAMES")
			default:
				rest = append(rest, a)
			}
		}
		fn := s.matchFun(c, m[1], c.env)
		vals, names := s.lapplyCore(c, m[0], fn, rest)
		names = s.useNames(m[0], names, useN)
		if !simp {
			return listOf(vals, names)
		}
		return s.simplify(c, vals, names)
	})
	regBuiltin("vapply", func(s *Session, c *callInfo) Value {
		m, extra := s.splitArgs(c, []string{"X", "FUN", "FUN.VALUE"})
		for i, n := range []string{"X", "FUN", "FUN.VALUE"} {
			if m[i] == nil {
				s.errorf(c.call, "argument \"%s\" is missing, with no default", n)
			}
		}
		useN := true
		var rest []Arg
		for _, a := range extra {
			if a.Name == "USE.NAMES" {
				useN = s.truthy(c, a.Value, "USE.NAMES")
			} else {
				rest = append(rest, a)
			}
		}
		fn := s.matchFun(c, m[1], c.env)
		fv := m[2]
		if !isVector(fv) {
			s.errorf(c.call, "'FUN.VALUE' must be a vector")
		}
		vals, names := s.lapplyCore(c, m[0], fn, rest)
		names = s.useNames(m[0], names, useN)
		n := Length(fv)
		for i, v := range vals {
			if Length(v) != n {
				s.errorf(c.call, "values must be length %d,\n but FUN(X[[%d]]) result is length %d", n, i+1, Length(v))
			}
			kv, kf := kindOf(v), kindOf(fv)
			ok := kv == kf || (kf == kDouble && (kv == kInteger || kv == kLogical)) || (kf == kInteger && kv == kLogical)
			if !ok {
				s.errorf(c.call, "values must be type '%s',\n but FUN(X[[%d]]) result is type '%s'", typeName(fv), i+1, typeName(v))
			}
		}
		if len(vals) == 0 {
			if n != 1 {
				return setAttr(emptyOf(kindOf(fv)), "dim", &Integer{V: []int32{int32(n), 0}})
			}
			return emptyOf(kindOf(fv))
		}
		conv := make([]Value, len(vals))
		for i, v := range vals {
			conv[i] = coerceTo(v, kindOf(fv))
		}
		if n == 1 {
			return s.simplify(c, conv, names)
		}
		return s.simplify(c, conv, names)
	})
	mapCore := func(s *Session, c *callInfo, fnV Value, dots []Arg, more []Arg, useNames bool) ([]Value, *String) {
		fn := s.matchFun(c, fnV, c.env)
		lens := make([]int, len(dots))
		elems := make([][]Value, len(dots))
		maxLen, anyZero := 0, false
		for i, d := range dots {
			elems[i] = s.asListElems(c, d.Value)
			lens[i] = len(elems[i])
			if lens[i] == 0 {
				anyZero = true
			}
			if lens[i] > maxLen {
				maxLen = lens[i]
			}
		}
		if anyZero {
			maxLen = 0
			if useNames && len(dots) > 0 {
				if nm := namesOf(dots[0].Value); nm != nil && len(nm.V) > 0 {
					s.errorf(&Lang{Fun: Intern("<-"), Args: []Arg{{Value: &Lang{Fun: Intern("names"), Args: []Arg{{Value: Intern("answer")}}}}, {Value: Intern("names1")}}},
						"'names' attribute [%d] must be the same length as the vector [0]", len(nm.V))
				}
			}
		}
		for _, l := range lens {
			if maxLen > 0 && l > 0 && maxLen%l != 0 {
				s.warningf(c.call, "longer argument not a multiple of length of shorter")
				break
			}
		}
		out := make([]Value, maxLen)
		for k := 0; k < maxLen; k++ {
			args := make([]Arg, 0, len(dots)+len(more))
			for i, d := range dots {
				args = append(args, Arg{Name: d.Name, HasName: d.HasName, Value: elems[i][k%lens[i]]})
			}
			args = append(args, more...)
			out[k] = s.callWith(fn, c.call, args, c.env)
		}
		var names *String
		if useNames && len(dots) > 0 {
			names = s.useNames(dots[0].Value, namesOf(dots[0].Value), true)
		}
		return out, names
	}
	regBuiltin("Map", func(s *Session, c *callInfo) Value {
		if len(c.args) == 0 {
			s.errorf(c.call, "argument \"f\" is missing, with no default")
		}
		fi := -1
		for i, n := range c.names {
			if n == "f" {
				fi = i
			}
		}
		if fi < 0 {
			for i, n := range c.names {
				if n == "" {
					fi = i
					break
				}
			}
		}
		if fi < 0 {
			s.errorf(c.call, "argument \"f\" is missing, with no default")
		}
		var dots []Arg
		for i := range c.args {
			if i != fi {
				dots = append(dots, Arg{Name: c.names[i], HasName: c.names[i] != "", Value: c.args[i]})
			}
		}
		vals, names := mapCore(s, c, c.args[fi], dots, nil, true)
		return listOf(vals, names)
	})
	regBuiltin("mapply", func(s *Session, c *callInfo) Value {
		fi := -1
		for i, n := range c.names {
			if n == "FUN" {
				fi = i
			}
		}
		if fi < 0 {
			for i, n := range c.names {
				if n == "" {
					fi = i
					break
				}
			}
		}
		if fi < 0 {
			s.errorf(c.call, "argument \"FUN\" is missing, with no default")
		}
		simp, useN := true, true
		var more []Arg
		var dots []Arg
		for i := range c.args {
			switch {
			case i == fi:
			case c.names[i] == "SIMPLIFY":
				simp = s.truthy(c, c.args[i], "SIMPLIFY")
			case c.names[i] == "USE.NAMES":
				useN = s.truthy(c, c.args[i], "USE.NAMES")
			case c.names[i] == "MoreArgs":
				if l, ok := c.args[i].(*List); ok {
					nm := namesOf(l)
					for j, v := range l.V {
						a := Arg{Value: v}
						if nm != nil && nm.V[j].S != "" {
							a.Name, a.HasName = nm.V[j].S, true
						}
						more = append(more, a)
					}
				}
			default:
				dots = append(dots, Arg{Name: c.names[i], HasName: c.names[i] != "", Value: c.args[i]})
			}
		}
		vals, names := mapCore(s, c, c.args[fi], dots, more, useN)
		if !simp {
			return listOf(vals, names)
		}
		return s.simplify(c, vals, names)
	})
	regBuiltin("Reduce", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "f", "x", "init", "right", "accumulate", "simplify")
		if a[0] == nil || a[1] == nil {
			s.errorf(c.call, "argument \"%s\" is missing, with no default", map[bool]string{true: "f", false: "x"}[a[0] == nil])
		}
		fn := s.matchFun(c, a[0], c.env)
		right := a[3] != nil && s.truthy(c, a[3], "right")
		acc := a[4] != nil && s.truthy(c, a[4], "accumulate")
		simp := a[5] == nil || s.truthy(c, a[5], "simplify")
		elems := append([]Value(nil), s.asListElems(c, a[1])...)
		if right {
			for i, j := 0, len(elems)-1; i < j; i, j = i+1, j-1 {
				elems[i], elems[j] = elems[j], elems[i]
			}
		}
		var cur Value
		start := 0
		if a[2] != nil {
			cur = a[2]
		} else {
			if len(elems) == 0 {
				return Null
			}
			cur = elems[0]
			start = 1
		}
		call := &Lang{Fun: Intern("f"), Args: []Arg{{Value: Intern("init")}, {Value: &Lang{Fun: Intern("[["), Args: []Arg{{Value: Intern("x")}, {Value: Intern("i")}}}}}}
		hist := []Value{cur}
		for _, e := range elems[start:] {
			var args []Arg
			if right {
				args = []Arg{{Value: e}, {Value: cur}}
			} else {
				args = []Arg{{Value: cur}, {Value: e}}
			}
			cur = s.callWith(fn, call, args, c.env)
			hist = append(hist, cur)
		}
		if !acc {
			return cur
		}
		if right {
			for i, j := 0, len(hist)-1; i < j; i, j = i+1, j-1 {
				hist[i], hist[j] = hist[j], hist[i]
			}
		}
		if simp {
			allOne := true
			for _, h := range hist {
				if Length(h) != 1 || !isVector(h) {
					allOne = false
				}
			}
			if allOne {
				return s.simplify(c, hist, nil)
			}
		}
		return listOf(hist, nil)
	})
	regBuiltin("Filter", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "f", "x")
		if a[0] == nil || a[1] == nil {
			s.errorf(c.call, "argument \"%s\" is missing, with no default", map[bool]string{true: "f", false: "x"}[a[0] == nil])
		}
		fn := s.matchFun(c, a[0], c.env)
		elems := s.asListElems(c, a[1])
		var flat []int32
		for _, e := range elems {
			r := s.callWith(fn, c.call, []Arg{{Value: e}}, c.env)
			flat = append(flat, logicalOf(s, c, r, "invalid argument")...)
		}
		var idx []int32
		for i, v := range flat {
			if v != 0 && v != naInt {
				idx = append(idx, int32(i+1))
			}
		}
		if idx == nil {
			idx = []int32{}
		}
		return s.callWith(builtins["["], c.call, []Arg{{Value: a[1]}, {Value: &Integer{V: idx}}}, c.env)
	})
	findPos := func(name string) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			a := s.matchFormals(c, "f", "x", "right", "nomatch")
			if a[0] == nil || a[1] == nil {
				s.errorf(c.call, "argument \"%s\" is missing, with no default", map[bool]string{true: "f", false: "x"}[a[0] == nil])
			}
			fn := s.matchFun(c, a[0], c.env)
			elems := s.asListElems(c, a[1])
			right := a[2] != nil && s.truthy(c, a[2], "right")
			idx := make([]int, len(elems))
			for i := range idx {
				if right {
					idx[i] = len(elems) - 1 - i
				} else {
					idx[i] = i
				}
			}
			for _, i := range idx {
				r := s.callWith(fn, c.call, []Arg{{Value: elems[i]}}, c.env)
				l := logicalOf(s, c, r, "invalid argument")
				if len(l) > 0 && l[0] != naInt && l[0] != 0 {
					if name == "Find" {
						return elems[i]
					}
					return mkInt(int32(i + 1))
				}
			}
			if name == "Find" {
				return Null
			}
			if a[3] != nil {
				return a[3]
			}
			return &Integer{V: []int32{naInt}}
		})
	}
	findPos("Find")
	findPos("Position")
	regBuiltin("do.call", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "what", "args", "quote", "envir")
		if a[0] == nil {
			s.errorf(c.call, "argument \"what\" is missing, with no default")
		}
		env := c.env
		if a[3] != nil {
			env = s.envArg(c, a[3], "envir")
		}
		var fun Value
		switch w := a[0].(type) {
		case *Closure, *Builtin:
			fun = w
		case *String:
			if len(w.V) != 1 || w.V[0].NA {
				s.errorf(c.call, "'what' must be a function or character string")
			}
			fun = Intern(w.V[0].S)
		default:
			s.errorf(c.call, "'what' must be a function or character string")
		}
		var args []Arg
		switch l := a[1].(type) {
		case nil:
		case nullT:
		case *List:
			nm := namesOf(l)
			for i, v := range l.V {
				arg := Arg{Value: v}
				if nm != nil && nm.V[i].S != "" && !nm.V[i].NA {
					arg.Name, arg.HasName = nm.V[i].S, true
				}
				switch v.(type) {
				case *Lang, *Symbol:
					if a[2] != nil && s.truthy(c, a[2], "quote") {
						arg.Value = &Lang{Fun: Intern("quote"), Args: []Arg{{Value: v}}}
					}
				}
				args = append(args, arg)
			}
		default:
			s.errorf(c.call, "second argument must be a list")
		}
		call := &Lang{Fun: fun, Args: args}
		return s.eval(call, env)
	})
	regBuiltin("match.fun", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "FUN", "descend")
		if a[0] == nil {
			s.errorf(c.call, "argument \"FUN\" is missing, with no default")
		}
		return s.matchFun(c, a[0], c.env)
	})
	regBuiltin("identity", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		return a[0]
	})

	regBuiltin("unlist", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "recursive", "use.names")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		rec := a[1] == nil || s.truthy(c, a[1], "recursive")
		use := a[2] == nil || s.truthy(c, a[2], "use.names")
		l, ok := a[0].(*List)
		if !ok {
			return a[0]
		}
		var leaves []Value
		var names []string
		s.flatten(l, "", rec, use, &leaves, &names)
		if len(leaves) == 0 {
			return Null
		}
		if !rec {
			anyList := false
			for _, e := range l.V {
				if _, ok := e.(*List); ok {
					anyList = true
				}
			}
			if anyList {
				var elems []Value
				for _, e := range l.V {
					switch v := e.(type) {
					case *List:
						elems = append(elems, v.V...)
					case nullT:
					default:
						for i := 0; i < Length(v); i++ {
							elems = append(elems, takeElems(withAttr(v, nil), []int{i}))
						}
					}
				}
				res := &List{V: elems}
				if use {
					any := false
					for _, n := range names {
						if n != "" {
							any = true
						}
					}
					if any && len(names) == len(elems) {
						res.Attr = namesAttr(strSlice(names))
					}
				}
				return res
			}
		}
		cc := &callInfo{call: c.call, env: c.env, args: leaves, names: make([]string, len(leaves))}
		res := s.combine(cc)
		if use {
			any := false
			for _, n := range names {
				if n != "" {
					any = true
				}
			}
			res = withAttr(res, nil)
			if any {
				res = withAttr(res, namesAttr(strSlice(names)))
			}
		} else {
			res = withAttr(res, nil)
		}
		return res
	})
	regBuiltin("rev", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x")
		x := a[0]
		if x == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		if _, isNull := x.(nullT); isNull {
			return Null
		}
		n := Length(x)
		pos := make([]int, n)
		for i := range pos {
			pos[i] = n - 1 - i
		}
		return s.subsetPositions(x, pos)
	})
	regBuiltin("which", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "arr.ind", "useNames")
		l, ok := a[0].(*Logical)
		if !ok {
			s.errorf(c.call, "argument to 'which' is not logical")
		}
		var out []int32
		var nm []Str
		names := namesOf(l)
		for i, v := range l.V {
			if v != 0 && v != naInt {
				out = append(out, int32(i+1))
				if names != nil {
					nm = append(nm, names.V[i])
				}
			}
		}
		res := &Integer{V: out}
		if res.V == nil {
			res.V = []int32{}
		}
		if names != nil {
			res.Attr = namesAttr(nm)
		}
		return res
	})
	regBuiltin("rep", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "times", "each", "length.out")
		x := a[0]
		if x == nil {
			s.errorf(c.call, "attempt to replicate an object of type 'NULL'")
		}
		if _, isNull := x.(nullT); isNull {
			if a[3] != nil {
				s.warningf(c.call, "'x' is NULL so the result will be NULL")
			}
			return Null
		}
		if !isVector(x) {
			s.errorf(c.call, "attempt to replicate an object of type '%s'", typeName(x))
		}
		n := Length(x)
		each := 1
		if a[2] != nil {
			each = s.countOf(c, a[2], "each")
		}
		base := make([]int, 0, n*each)
		for i := 0; i < n; i++ {
			for j := 0; j < each; j++ {
				base = append(base, i)
			}
		}
		var pos []int
		switch {
		case a[3] != nil:
			lo := s.countOf(c, a[3], "length.out")
			s.checkLen(c.call, lo)
			if len(base) == 0 && lo > 0 {
				pos = make([]int, lo)
				for i := range pos {
					pos[i] = -1
				}
			} else if len(base) > 0 {
				for i := 0; i < lo; i++ {
					pos = append(pos, base[i%len(base)])
				}
			}
		case a[1] != nil && len(base) == 0:
			pos = nil
		case a[1] != nil:
			t := asDoublesOrNil(a[1])
			if len(t) == len(base) && len(t) > 1 {
				total := 0
				for _, v := range t {
					if v != v || v < 0 {
						s.errorf(c.call, "invalid 'times' argument")
					}
					total += int(v)
					s.checkLen(c.call, total)
				}
				for i, v := range t {
					for j := 0; j < int(v); j++ {
						pos = append(pos, base[i])
					}
				}
			} else {
				if len(t) != 1 || t[0] != t[0] || t[0] < 0 {
					s.errorf(c.call, "invalid 'times' argument")
				}
				if t[0] > float64(s.opts.MaxVector) {
					s.errorf(c.call, "vector size too large")
				}
				times := int(t[0])
				if len(base) > 0 {
					s.checkLen(c.call, len(base)*times)
				}
				for k := 0; k < times; k++ {
					pos = append(pos, base...)
				}
			}
		default:
			pos = base
		}
		out := s.subsetPositions(x, pos)
		if nm := namesOf(x); nm == nil {
			return withAttr(out, nil)
		}
		return out
	})
	regBuiltin("seq", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "from", "to", "by", "length.out", "along.with")
		num := func(v Value, what string) float64 {
			d := asDoublesOrNil(v)
			if len(d) != 1 || d[0] != d[0] {
				s.errorf(c.call, "'%s' must be of length 1", what)
			}
			if math.IsInf(d[0], 0) {
				s.errorf(&Lang{Fun: Intern("seq.default"), Args: c.call.Args}, "'%s' must be a finite number", what)
			}
			return d[0]
		}
		isInt := func(v Value) bool { _, ok := v.(*Integer); return ok }
		if a[4] != nil {
			return s.seqLen(Length(a[4]))
		}
		if a[0] != nil && a[1] == nil && a[2] == nil && a[3] == nil {
			if Length(a[0]) == 1 {
				f := num(a[0], "from")
				if f >= 1 {
					return s.seqLen(int(f))
				}
				return s.colonValues(1, f)
			}
			return s.seqLen(Length(a[0]))
		}
		from, to := 1.0, 1.0
		if a[0] != nil {
			from = num(a[0], "from")
		}
		if a[1] != nil {
			to = num(a[1], "to")
		}
		allInt := (a[0] == nil || isInt(a[0]) || from == float64(int64(from))) && a[2] == nil
		switch {
		case a[3] != nil:
			lo := int(math.Min(num(a[3], "length.out"), float64(s.opts.MaxVector)+1))
			if lo < 0 {
				s.errorf(c.call, "'length.out' must be a non-negative number")
			}
			s.checkLen(c.call, lo)
			by := 1.0
			switch {
			case a[2] != nil:
				by = num(a[2], "by")
				if a[0] == nil && a[1] != nil {
					from = to - by*float64(lo-1)
				}
			case a[1] != nil && a[0] != nil:
				if lo > 1 {
					by = (to - from) / float64(lo-1)
				}
			case a[1] != nil:
				from = to - float64(lo-1)
			}
			out := make([]float64, lo)
			for i := range out {
				out[i] = from + float64(i)*by
			}
			if (a[0] == nil && a[1] == nil && a[2] == nil) || (a[2] != nil && isInt(a[2]) && a[1] == nil && (a[0] == nil || isInt(a[0]))) {
				return toIntVec(out)
			}
			return &Double{V: out}
		case a[2] != nil:
			by := num(a[2], "by")
			if by == 0 && from == to {
				return a[0]
			}
			if (to-from)/by < 0 || by == 0 {
				s.errorf(&Lang{Fun: Intern("seq.default"), Args: c.call.Args}, "wrong sign in 'by' argument")
			}
			nf := (to-from)/by + 1e-10
			if nf > float64(s.opts.MaxVector) {
				s.errorf(c.call, "vector size too large")
			}
			n := int(nf)
			s.checkLen(c.call, n+1)
			out := make([]float64, n+1)
			for i := range out {
				out[i] = from + float64(i)*by
			}
			if (a[0] == nil || isInt(a[0])) && isInt(a[2]) {
				return toIntVec(out)
			}
			return &Double{V: out}
		}
		_ = allInt
		return s.colonValues(from, to)
	})
}

func allWhole(x []float64) bool {
	for _, v := range x {
		if v != float64(int32(v)) {
			return false
		}
	}
	return true
}

func toIntVec(x []float64) Value {
	out := make([]int32, len(x))
	for i, v := range x {
		out[i] = int32(v)
	}
	return &Integer{V: out}
}

func strSlice(ss []string) []Str {
	out := make([]Str, len(ss))
	for i, s := range ss {
		out[i] = Str{S: s}
	}
	return out
}

func (s *Session) seqLen(n int) Value {
	s.checkLen(nil, n)
	out := make([]int32, n)
	for i := range out {
		out[i] = int32(i + 1)
	}
	return &Integer{V: out}
}

// colonValues is from:to for doubles, as the colon operator does.
func (s *Session) colonValues(from, to float64) Value {
	return s.colon(&callInfo{args: []Value{mkDouble(from), mkDouble(to)}, names: []string{"", ""}})
}

func (s *Session) countOf(c *callInfo, v Value, what string) int {
	d := asDoublesOrNil(v)
	if len(d) != 1 || d[0] != d[0] || d[0] < 0 {
		s.errorf(c.call, "invalid '%s' argument", what)
	}
	if d[0] > float64(s.opts.MaxVector) {
		s.errorf(c.call, "vector size too large")
	}
	return int(d[0])
}

// subsetPositions is x[pos+1] for plain positions, keeping names.
func (s *Session) subsetPositions(x Value, pos []int) Value {
	if l, ok := x.(*List); ok {
		out := make([]Value, len(pos))
		for i, p := range pos {
			if p < 0 || p >= len(l.V) {
				out[i] = Null
			} else {
				out[i] = l.V[p]
			}
		}
		res := &List{V: out}
		if nm := namesOf(l); nm != nil {
			names := make([]Str, len(pos))
			for i, p := range pos {
				names[i] = nameAt(nm, p)
			}
			res.Attr = namesAttr(names)
		}
		return res
	}
	res := takeElems(x, pos)
	if nm := namesOf(x); nm != nil {
		names := make([]Str, len(pos))
		for i, p := range pos {
			names[i] = nameAt(nm, p)
		}
		res = withAttr(res, namesAttr(names))
	}
	if cl := classAttr(x); len(cl) == 1 && cl[0] == "factor" {
		res = withAttr(res, append(attrsOf(res), Attr{Name: "levels", Value: getAttr(x, "levels")}, Attr{Name: "class", Value: getAttr(x, "class")}))
	}
	return res
}

// flatten collects the leaves of nested lists with R's name rules: a name
// given to a vector of n > 1 unnamed leaves becomes name1 ... namen.
func (s *Session) flatten(l *List, base string, rec, use bool, leaves *[]Value, names *[]string) {
	nm := namesOf(l)
	for i, e := range l.V {
		tag := ""
		if nm != nil && !nm.V[i].NA {
			tag = nm.V[i].S
		}
		name := joinName(base, tag)
		switch v := e.(type) {
		case *List:
			if rec {
				s.flatten(v, name, rec, use, leaves, names)
				continue
			}
			s.leafParts(e, name, use, leaves, names)
		case nullT:
		default:
			s.leafParts(e, name, use, leaves, names)
		}
	}
}

func joinName(base, tag string) string {
	switch {
	case base != "" && tag != "":
		return base + "." + tag
	case base != "":
		return base
	}
	return tag
}

// leafParts adds an atomic vector (or, for recursive = FALSE, a list's
// elements) as leaves, numbering the names of a named vector of length > 1.
func (s *Session) leafParts(e Value, name string, use bool, leaves *[]Value, names *[]string) {
	n := Length(e)
	inner := namesOf(e)
	if l, ok := e.(*List); ok {
		for i, v := range l.V {
			*leaves = append(*leaves, v)
			tag := ""
			if inner != nil {
				tag = inner.V[i].S
			}
			*names = append(*names, partName(name, tag, i, n))
		}
		return
	}
	*leaves = append(*leaves, withAttr(e, nil))
	for i := 0; i < n; i++ {
		tag := ""
		if inner != nil && !inner.V[i].NA {
			tag = inner.V[i].S
		}
		*names = append(*names, partName(name, tag, i, n))
	}
	// the leaf above stands for n elements: keep one name per element
}

func partName(base, tag string, i, n int) string {
	switch {
	case base != "" && tag != "":
		return base + "." + tag
	case base != "" && n > 1:
		return base + strconv.Itoa(i+1)
	case base != "":
		return base
	}
	return tag
}

func nameAt(nm *String, p int) Str {
	if p < 0 || p >= len(nm.V) {
		return Str{NA: true}
	}
	return nm.V[p]
}
