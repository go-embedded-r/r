package eval

import "sort"

// Locked and active bindings, delayed assignment, and the checks every
// assignment makes: written from the documentation of lockBinding,
// makeActiveBinding and delayedAssign and from what the oracle answers.

// frame is the environment that owns the bindings of e: e itself, or the
// environment e is another name for.
func (e *Env) frame() *Env {
	if e.alias != nil {
		return e.alias
	}
	return e
}

// bval is the value a binding currently has: for an active binding, what its
// function returns. Promises are left alone.
func (s *Session) bval(b *binding) Value {
	if b.active != nil {
		return s.callFunction(b.active, nil)
	}
	return b.val
}

// bforce is bval with a promise forced.
func (s *Session) bforce(b *binding) Value {
	v := s.bval(b)
	if p, ok := v.(*Promise); ok {
		return s.force(p)
	}
	return v
}

// assignTo binds sym to v in e as an assignment written by the user does:
// locked bindings and environments refuse it, an active binding runs its
// function with v, and the empty environment holds nothing.
func (s *Session) assignTo(e *Env, sym *Symbol, v Value, call *Lang) {
	if !e.special && e.alias == nil && e != s.Empty {
		e.set(sym, v)
		return
	}
	if e == s.Empty {
		s.errorf(call, "cannot assign values in the empty environment")
	}
	f := e.frame()
	if b := f.lookup(sym); b != nil {
		if b.locked || f.allLocked {
			s.errorf(call, "cannot change value of locked binding for '%s'", sym.Name)
		}
		if b.active != nil {
			s.callFunction(b.active, []Value{v})
			return
		}
	} else if f.locked {
		if f == s.Base {
			s.errorf(call, "cannot add binding of '%s' to the base environment", sym.Name)
		}
		s.errorf(call, "cannot add bindings to a locked environment")
	}
	f.set(sym, v)
}

// lockAllBindings locks every binding of the frame, present and future ones
// excluded.
func (e *Env) lockAllBindings() {
	e.special = true
	if e.big != nil {
		for _, b := range e.big {
			b.locked = true
		}
		return
	}
	for i := range e.small {
		e.small[i].locked = true
	}
}

// bindingLocked reports whether the binding of f is locked.
func (e *Env) bindingLocked(b *binding) bool { return b.locked || e.allLocked }

// materializeLocks turns the frame-wide lock of the base frame into one
// lock per binding, so that a single one can be undone.
func (e *Env) materializeLocks() {
	if e.allLocked {
		e.allLocked = false
		e.lockAllBindings()
	}
}

func (s *Session) symArg(c *callInfo, v Value) *Symbol {
	switch x := v.(type) {
	case *Symbol:
		return x
	case *String:
		if len(x.V) >= 1 && !x.V[0].NA {
			return Intern(x.V[0].S)
		}
	}
	s.errorf(c.call, "not a symbol")
	return nil
}

func (s *Session) envOnly(c *callInfo, v Value) *Env {
	e, ok := v.(*Env)
	if !ok {
		s.errorf(c.call, "not an environment")
	}
	return e
}

// boundIn is the binding of sym in the frame of e, or an R error.
func (s *Session) boundIn(c *callInfo, e *Env, sym *Symbol) *binding {
	b := e.lookup(sym)
	if b == nil {
		s.errorf(c.call, "no binding for \"%s\"", sym.Name)
	}
	return b
}

func init() {
	regBuiltin("lockBinding", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "sym", "env")
		sym, e := s.symArg(c, a[0]), s.envOnly(c, a[1])
		b := s.boundIn(c, e, sym)
		b.locked = true
		e.frame().special = true
		return s.invisible(Null)
	})
	regBuiltin("unlockBinding", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "sym", "env")
		sym, e := s.symArg(c, a[0]), s.envOnly(c, a[1])
		e.frame().materializeLocks()
		s.boundIn(c, e, sym).locked = false
		return s.invisible(Null)
	})
	regBuiltin("bindingIsLocked", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "sym", "env")
		sym, e := s.symArg(c, a[0]), s.envOnly(c, a[1])
		return mkLogical(e.frame().bindingLocked(s.boundIn(c, e, sym)))
	})
	regBuiltin("lockEnvironment", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "env", "bindings")
		e := s.envOnly(c, a[0])
		bindings := a[1] != nil && s.truthy(c, a[1], "bindings")
		f := e.frame()
		f.locked, f.special = true, true
		if bindings {
			f.lockAllBindings()
		}
		return s.invisible(Null)
	})
	regBuiltin("environmentIsLocked", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "env")
		return mkLogical(s.envOnly(c, a[0]).frame().locked)
	})
	regBuiltin("makeActiveBinding", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "sym", "fun", "env")
		sym := s.symArg(c, a[0])
		switch a[1].(type) {
		case *Closure, *Builtin:
		default:
			s.errorf(c.call, "not a function")
		}
		e := s.envOnly(c, a[2])
		f := e.frame()
		if b := f.lookup(sym); b != nil {
			if b.active == nil {
				s.errorf(c.call, "symbol already has a regular binding")
			}
			if f.bindingLocked(b) {
				s.errorf(c.call, "cannot change active binding if binding is locked")
			}
			b.active = a[1]
			return s.invisible(Null)
		}
		if f.locked {
			s.errorf(c.call, "cannot add bindings to a locked environment")
		}
		f.set(sym, Null)
		f.lookup(sym).active = a[1]
		f.special = true
		return s.invisible(Null)
	})
	regBuiltin("bindingIsActive", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "sym", "env")
		sym, e := s.symArg(c, a[0]), s.envOnly(c, a[1])
		return mkLogical(s.boundIn(c, e, sym).active != nil)
	})
	regBuiltin("activeBindingFunction", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "sym", "env")
		sym, e := s.symArg(c, a[0]), s.envOnly(c, a[1])
		b := s.boundIn(c, e, sym)
		if b.active == nil {
			s.errorf(c.call, "no active binding for \"%s\"", sym.Name)
		}
		return b.active
	})
	regSpecial("delayedAssign", func(s *Session, c *callInfo) Value {
		formals := []string{"x", "value", "eval.env", "assign.env"}
		got := make([]Value, len(formals))
		have := make([]bool, len(formals))
		used := make([]bool, len(c.raw))
		for i, ra := range c.raw {
			if !ra.HasName {
				continue
			}
			for j, f := range formals {
				if f == ra.Name && !have[j] {
					got[j], have[j], used[i] = ra.Value, true, true
					break
				}
			}
		}
		j := 0
		for i, ra := range c.raw {
			if used[i] {
				continue
			}
			for j < len(formals) && have[j] {
				j++
			}
			if j >= len(formals) {
				s.errorf(c.call, "unused argument")
			}
			got[j], have[j] = ra.Value, true
		}
		if !have[0] {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		if !have[1] {
			s.errorf(c.call, "argument \"value\" is missing, with no default")
		}
		name, ok := s.eval(got[0], c.env).(*String)
		if !ok || len(name.V) == 0 || name.V[0].NA {
			s.errorf(c.call, "invalid first argument")
		}
		evalEnv, assignEnv := c.env, c.env
		if have[2] {
			evalEnv = s.envArg(c, s.eval(got[2], c.env), "eval.env")
		}
		if have[3] {
			assignEnv = s.envArg(c, s.eval(got[3], c.env), "assign.env")
		}
		s.assignTo(assignEnv, Intern(name.V[0].S), &Promise{Expr: got[1], Env: evalEnv}, c.call)
		return s.invisible(Null)
	})

	// as.list.environment, with the values forced; extra arguments are the
	// `...` of the method and ignored.
	regBuiltin("as.list.environment", func(s *Session, c *callInfo) Value {
		a := s.matchLoose(c, "x", "all.names", "sorted")
		e, ok := a[0].(*Env)
		if !ok {
			s.errorf(c.call, "argument must be an environment")
		}
		all := a[1] != nil && logicalFlag(a[1])
		sorted := a[2] != nil && logicalFlag(a[2])
		var names []string
		for _, n := range e.Names() {
			if !all && len(n) > 0 && n[0] == '.' {
				continue
			}
			names = append(names, n)
		}
		if sorted {
			sort.SliceStable(names, func(i, j int) bool { return collate(names[i], names[j]) < 0 })
		}
		out := &List{V: make([]Value, len(names))}
		nm := make([]Str, len(names))
		for i, n := range names {
			out.V[i] = s.bforce(e.lookup(Intern(n)))
			nm[i] = Str{S: n}
		}
		if len(names) > 0 {
			out.Attr = namesAttr(nm)
		}
		return out
	})
}

func init() {
	// withVisible(x): the value of x and whether R would have printed it.
	regSpecial("withVisible", func(s *Session, c *callInfo) Value {
		if len(c.raw) != 1 {
			s.errorf(c.call, "%d arguments passed to 'withVisible' which requires 1", len(c.raw))
		}
		s.visible = true
		v := s.eval(c.raw[0].Value, c.env)
		vis := s.visible
		s.visible = true
		return &List{V: []Value{v, mkLogical(vis)}, Attr: namesAttr([]Str{{S: "value"}, {S: "visible"}})}
	})
	regBuiltin("as.list.default", func(s *Session, c *callInfo) Value {
		a := s.matchLoose(c, "x")
		switch x := a[0].(type) {
		case *List:
			return x
		case *Lang:
			return callAsList(x)
		}
		asVec := &Lang{Fun: Intern("as.vector"), Args: []Arg{{Value: Intern("x")}, {Value: mkString("list")}}}
		switch a[0].(type) {
		case *Env, *Closure, *Builtin:
			s.errorf(asVec, "cannot coerce type '%s' to vector of type 'list'", typeName(a[0]))
		}
		ci := &callInfo{call: asVec, env: c.env, args: []Value{a[0], mkString("list")}, names: []string{"", ""}}
		return builtins["as.vector"].fn(s, ci)
	})
	// as.list of a closure: its formals, then its body.
	regBuiltin("as.list.function", func(s *Session, c *callInfo) Value {
		a := s.matchLoose(c, "x")
		f, ok := a[0].(*Closure)
		if !ok {
			s.errorf(c.call, "argument is not a function")
		}
		out := &List{V: make([]Value, 0, len(f.Formals)+1)}
		names := make([]Str, 0, len(f.Formals)+1)
		for _, fm := range f.Formals {
			d := fm.Default
			if d == nil {
				d = Missing
			}
			out.V = append(out.V, d)
			names = append(names, Str{S: fm.Name.Name})
		}
		out.V = append(out.V, f.Body)
		names = append(names, Str{S: ""})
		out.Attr = namesAttr(names)
		return out
	})
}

// matchLoose is matchFormals for a method whose signature ends in `...`:
// arguments that match no formal are ignored.
func (s *Session) matchLoose(c *callInfo, formals ...string) []Value {
	out := make([]Value, len(formals))
	used := make([]bool, len(c.args))
	for i, n := range c.names {
		for j, f := range formals {
			if n != "" && f == n && out[j] == nil {
				out[j], used[i] = c.args[i], true
				break
			}
		}
	}
	j := 0
	for i := range c.args {
		if used[i] || c.names[i] != "" {
			continue
		}
		for j < len(formals) && out[j] != nil {
			j++
		}
		if j < len(formals) {
			out[j] = c.args[i]
		}
	}
	return out
}

// callAsList is as.list of a call: the function and the arguments, tagged as
// written.
func callAsList(l *Lang) Value {
	out := &List{V: make([]Value, 0, len(l.Args)+1)}
	names := []Str{{S: ""}}
	out.V = append(out.V, l.Fun)
	tagged := false
	for _, arg := range l.Args {
		out.V = append(out.V, arg.Value)
		names = append(names, Str{S: arg.Name})
		tagged = tagged || (arg.HasName && arg.Name != "")
	}
	if tagged {
		out.Attr = namesAttr(names)
	}
	return out
}

// logicalFlag reads a TRUE/FALSE argument the way R's asLogical does: NA and
// anything that is not a logical-like scalar count as FALSE.
func logicalFlag(v Value) bool {
	switch x := v.(type) {
	case *Logical:
		return len(x.V) > 0 && x.V[0] != naInt && x.V[0] != 0
	case *Integer:
		return len(x.V) > 0 && x.V[0] != naInt && x.V[0] != 0
	case *Double:
		return len(x.V) > 0 && x.V[0] == x.V[0] && x.V[0] != 0
	}
	return false
}
