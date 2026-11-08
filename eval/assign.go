package eval

// setVar binds sym in env, or for <<- in the nearest enclosing environment
// that has it (the global one when none does).
func (s *Session) setVar(env *Env, sym *Symbol, v Value, super bool) {
	if !super {
		s.assignTo(env, sym, v, s.currentCall())
		return
	}
	target := s.Global
	for e := env.parent; e != nil; e = e.parent {
		if e.lookup(sym) != nil {
			target = e
			break
		}
	}
	s.assignTo(target, sym, v, s.currentCall())
}

// targetValue is the current value of the variable (or expression) that a
// replacement modifies.
func (s *Session) targetValue(c *callInfo, inner Value, super bool) Value {
	var sym *Symbol
	switch t := inner.(type) {
	case *Symbol:
		sym = t
	case *String:
		if len(t.V) == 1 {
			sym = Intern(t.V[0].S)
		}
	case *Lang:
		return s.eval(t, c.env)
	}
	if sym == nil {
		s.errorf(c.call, "target of assignment expands to non-language object")
	}
	start := c.env
	if super {
		start = c.env.parent
	}
	var b *binding
	if start != nil {
		b, _ = start.find(sym)
	}
	if b == nil {
		s.errorf(nil, "object '%s' not found", sym.Name)
	}
	if b.active != nil {
		return s.callFunction(b.active, nil)
	}
	switch v := b.val.(type) {
	case *Promise:
		return s.force(v)
	case missingT:
		s.errorf(c.call, "argument \"%s\" is missing, with no default", sym.Name)
	}
	return b.val
}

// complexAssign performs f(g(x), a) <- value as R does: x <- `g<-`(x,
// `f<-`(g(x), a, value = value)), one replacement at a time.
func (s *Session) complexAssign(c *callInfo, target Value, value Value, super bool) {
	switch t := target.(type) {
	case *Symbol:
		s.setVar(c.env, t, value, super)
	case *String:
		if len(t.V) != 1 {
			s.errorf(c.call, "invalid (do_set) left-hand side to assignment")
		}
		s.setVar(c.env, Intern(t.V[0].S), value, super)
	case *Lang:
		f, ok := t.Fun.(*Symbol)
		if !ok {
			s.errorf(c.call, "invalid function in complex assignment")
		}
		if len(t.Args) == 0 {
			s.errorf(c.call, "invalid assignment target")
		}
		inner := t.Args[0].Value
		if _, isNull := inner.(nullT); isNull {
			s.errorf(c.call, "invalid (NULL) left side of assignment")
		}
		cur := s.targetValue(c, inner, super)
		nv := s.callReplacement(c, f.Name, cur, t.Args[1:], value)
		s.complexAssign(c, inner, nv, super)
	default:
		s.errorf(c.call, "target of assignment expands to non-language object")
	}
}

func (s *Session) callReplacement(c *callInfo, fname string, cur Value, extra []Arg, value Value) Value {
	rname := Intern(fname + "<-")
	fn := s.findFun(rname, c.env, c.call)
	switch f := fn.(type) {
	case *Builtin:
		tmp := &Lang{Fun: rname, Args: append([]Arg{{Value: Intern("*tmp*")}}, extra...)}
		ci := &callInfo{call: c.call, tmpCall: tmp, env: c.env, args: []Value{cur}, names: []string{""}}
		for i, a := range extra {
			if sym, isSym := a.Value.(*Symbol); isSym && (fname == "$" || fname == "@") && i == 0 {
				ci.args = append(ci.args, mkString(sym.Name))
				ci.names = append(ci.names, "")
				continue
			}
			if sym, isSym := a.Value.(*Symbol); isSym && sym.Name == "..." {
				if d := s.dotsOf(c.env); d != nil {
					for _, da := range d.Args {
						ci.args = append(ci.args, s.argValue(da.Value))
						ci.names = append(ci.names, da.Name)
					}
				}
				continue
			}
			ci.args = append(ci.args, s.argValue(a.Value, c.env))
			ci.names = append(ci.names, a.Name)
		}
		ci.args = append(ci.args, value)
		ci.names = append(ci.names, "value")
		tmp.Args = append(tmp.Args, Arg{Name: "value", HasName: true, Value: value})
		return f.fn(s, ci)
	case *Closure:
		supplied := []Arg{{Value: cur}}
		supplied = append(supplied, s.promiseArgs(extra, c.env)...)
		supplied = append(supplied, Arg{Name: "value", HasName: true, Value: value})
		call := &Lang{Fun: rname, Args: append(append([]Arg{{Value: Intern("*tmp*")}}, extra...), Arg{Name: "value", HasName: true, Value: value})}
		return s.applyClosure(f, call, supplied, c.env, nil)
	}
	return Null
}
