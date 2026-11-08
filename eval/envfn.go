package eval

import (
	"sort"

	"github.com/go-embedded-r/r/syntax"
)

// Environments as values, and the functions that look at the call stack.

// frameOf is the call frame whose environment is env, or nil.
func (s *Session) frameIndex(env *Env) (int, *frame) {
	for i := len(s.frames) - 1; i >= 0; i-- {
		if s.frames[i].env == env {
			return i, s.frames[i]
		}
	}
	return -1, nil
}

// parentFrame is parent.frame(n) seen from env.
func (s *Session) parentFrame(env *Env, n int) *Env {
	e := env
	for i := 0; i < n; i++ {
		_, fr := s.frameIndex(e)
		if fr == nil {
			return s.Global
		}
		e = fr.caller
		if e == nil {
			return s.Global
		}
	}
	return e
}

func (s *Session) envArg(c *callInfo, v Value, what string) *Env {
	switch e := v.(type) {
	case *Env:
		return e
	case nullT:
		s.errorf(c.call, "use of NULL environment is defunct")
	}
	s.errorf(c.call, "invalid '%s' argument", what)
	return nil
}

func (s *Session) stringArg(c *callInfo, v Value, what string) string {
	st, ok := v.(*String)
	if !ok || len(st.V) < 1 || st.V[0].NA {
		s.errorf(c.call, "invalid '%s' argument", what)
	}
	return st.V[0].S
}

func (s *Session) intArg(c *callInfo, v Value, what string) int {
	d := asDoublesOrNil(v)
	if len(d) < 1 || d[0] != d[0] {
		s.errorf(c.call, "invalid '%s' argument", what)
	}
	return int(d[0])
}

func envName(s *Session, e *Env) string {
	switch e {
	case s.Global:
		return "R_GlobalEnv"
	case s.Base, s.BaseNS:
		return "base"
	case s.Empty:
		return "R_EmptyEnv"
	}
	return ""
}

func (s *Session) invisible(v Value) Value {
	s.visible = false
	return v
}

// getVar looks name up from env; mode "function" skips non-functions.
func (s *Session) getVar(name string, env *Env, inherits bool, fnOnly bool) (Value, bool) {
	sym := Intern(name)
	for e := env; e != nil; e = e.parent {
		if b := e.lookup(sym); b != nil {
			v := s.bval(b)
			if p, ok := v.(*Promise); ok {
				v = s.force(p)
			}
			if _, miss := v.(missingT); !miss {
				if !fnOnly {
					return v, true
				}
				switch v.(type) {
				case *Closure, *Builtin:
					return v, true
				}
			}
		}
		if !inherits {
			break
		}
	}
	return nil, false
}

func (s *Session) modeIsFunction(c *callInfo, mode Value) bool {
	if mode == nil {
		return false
	}
	return s.stringArg(c, mode, "mode") == "function"
}

func init() {
	regBuiltin("environment", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "fun")
		if a[0] == nil {
			return c.env
		}
		switch f := a[0].(type) {
		case nullT:
			return c.env
		case *Closure:
			return f.Env
		}
		return Null
	})
	regBuiltin("environment<-", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "fun", "value")
		f, ok := a[0].(*Closure)
		if !ok {
			s.errorf(c.call, "replacement object is not an environment")
		}
		e, ok := a[1].(*Env)
		if !ok {
			s.errorf(c.call, "replacement object is not an environment")
		}
		return &Closure{Formals: f.Formals, Body: f.Body, Env: e}
	})
	regBuiltin("new.env", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "hash", "parent", "size")
		parent := c.env
		if a[1] != nil {
			parent = s.envArg(c, a[1], "enclos")
		}
		return NewEnv(parent)
	})
	regBuiltin("globalenv", func(s *Session, c *callInfo) Value { return s.Global })
	regBuiltin("emptyenv", func(s *Session, c *callInfo) Value { return s.Empty })
	regBuiltin("baseenv", func(s *Session, c *callInfo) Value { return s.Base })
	regBuiltin("parent.env", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "env")
		e, ok := a[0].(*Env)
		if !ok {
			s.errorf(c.call, "argument is not an environment")
		}
		if e == s.Empty {
			s.errorf(c.call, "the empty environment has no parent")
		}
		if e.parent == nil {
			return s.Empty
		}
		return e.parent
	})
	regBuiltin("parent.env<-", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "env", "value")
		e, ok := a[0].(*Env)
		p, ok2 := a[1].(*Env)
		if !ok || !ok2 {
			s.errorf(c.call, "argument is not an environment")
		}
		e.parent = p
		return e
	})
	regBuiltin("parent.frame", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "n")
		n := 1
		if a[0] != nil {
			n = s.intArg(c, a[0], "n")
			if n < 1 {
				s.errorf(c.call, "invalid 'n' value")
			}
		}
		return s.parentFrame(c.env, n)
	})
	regBuiltin("sys.call", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "which")
		which := 0
		if a[0] != nil {
			which = s.intArg(c, a[0], "which")
		}
		cur, _ := s.frameIndex(c.env)
		idx := cur
		switch {
		case which > 0:
			idx = which - 1
		case which < 0:
			idx = cur + which
		}
		if idx < 0 && which == 0 {
			return Null
		}
		if idx < 0 || idx >= len(s.frames) {
			s.errorf(c.call, "not that many frames on the stack")
		}
		return s.frames[idx].call
	})
	regBuiltin("sys.function", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "which")
		which := 0
		if a[0] != nil {
			which = s.intArg(c, a[0], "which")
		}
		cur, _ := s.frameIndex(c.env)
		idx := cur
		switch {
		case which > 0:
			idx = which - 1
		case which < 0:
			idx = cur + which
		}
		if idx < 0 || idx >= len(s.frames) {
			s.errorf(c.call, "not that many frames on the stack")
		}
		return s.frames[idx].fn
	})
	regBuiltin("assign", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "value", "pos", "envir", "inherits", "immediate")
		if a[0] == nil || a[1] == nil {
			s.errorf(c.call, "argument \"%s\" is missing, with no default", map[bool]string{true: "x", false: "value"}[a[0] == nil])
		}
		name, ok := a[0].(*String)
		if !ok || len(name.V) == 0 {
			s.errorf(c.call, "invalid first argument")
		}
		env := c.env
		if a[3] != nil {
			env = s.envArg(c, a[3], "envir")
		} else if a[2] != nil {
			if e, ok := a[2].(*Env); ok {
				env = e
			}
		}
		inherits := a[4] != nil && s.truthy(c, a[4], "inherits")
		sym := Intern(name.V[0].S)
		target := env
		if inherits {
			if _, e := env.find(sym); e != nil {
				target = e
			}
		}
		s.assignTo(target, sym, a[1], c.call)
		return s.invisible(a[1])
	})
	regBuiltin("get", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "pos", "envir", "mode", "inherits")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		name := s.nameArg(c, a[0])
		env := c.env
		if a[2] != nil {
			env = s.envArg(c, a[2], "envir")
		} else if e, ok := a[1].(*Env); ok {
			env = e
		}
		inherits := a[4] == nil || s.truthy(c, a[4], "inherits")
		mode := s.modeArg(c, a[3])
		v, ok := s.lookupMode(name, env, inherits, mode, true)
		if !ok {
			if mode != "any" {
				s.errorf(c.call, "object '%s' of mode '%s' was not found", name, mode)
			}
			s.errorf(c.call, "object '%s' not found", name)
		}
		return v
	})
	regBuiltin("get0", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "envir", "mode", "inherits", "ifnotfound")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		name := s.nameArg(c, a[0])
		env := c.env
		if a[1] != nil {
			env = s.envArg(c, a[1], "envir")
		}
		inherits := a[3] == nil || s.truthy(c, a[3], "inherits")
		if v, ok := s.lookupMode(name, env, inherits, s.modeArg(c, a[2]), true); ok {
			return v
		}
		if a[4] != nil {
			return a[4]
		}
		return Null
	})
	regBuiltin("exists", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "where", "envir", "frame", "mode", "inherits")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		name := s.nameArg(c, a[0])
		env := c.env
		if a[2] != nil {
			env = s.envArg(c, a[2], "envir")
		} else if e, ok := a[1].(*Env); ok {
			env = e
		}
		inherits := a[5] == nil || s.truthy(c, a[5], "inherits")
		_, ok := s.lookupMode(name, env, inherits, s.modeArg(c, a[4]), false)
		return mkLogical(ok)
	})
	regBuiltin("ls", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "name", "pos", "envir", "all.names", "pattern", "sorted")
		env := c.env
		if a[2] != nil {
			env = s.envArg(c, a[2], "envir")
		} else if e, ok := a[0].(*Env); ok {
			env = e
		}
		all := a[3] != nil && s.truthy(c, a[3], "all.names")
		var names []string
		for _, n := range env.Names() {
			if _, ok := env.Get(n); !ok {
				continue
			}
			if !all && len(n) > 0 && n[0] == '.' {
				continue
			}
			names = append(names, n)
		}
		sort.SliceStable(names, func(i, j int) bool { return collate(names[i], names[j]) < 0 })
		return strVec(names)
	})
	regSpecial("rm", func(s *Session, c *callInfo) Value {
		env := c.env
		var names []string
		inherits := false
		for _, ra := range c.raw {
			switch {
			case ra.HasName && ra.Name == "inherits":
				inherits = s.truthy(c, s.eval(ra.Value, c.env), "inherits")
			case ra.HasName && ra.Name == "list":
				v := s.eval(ra.Value, c.env)
				if st, ok := v.(*String); ok {
					for _, e := range st.V {
						names = append(names, e.S)
					}
				}
			case ra.HasName && ra.Name == "envir":
				env = s.envArg(c, s.eval(ra.Value, c.env), "envir")
			case ra.HasName:
			default:
				switch x := ra.Value.(type) {
				case *Symbol:
					names = append(names, x.Name)
				case *String:
					for _, e := range x.V {
						names = append(names, e.S)
					}
				default:
					s.errorf(c.call, "... must contain names or character strings")
				}
			}
		}
		for _, n := range names {
			sym := Intern(n)
			target := env
			if inherits {
				_, e := env.find(sym)
				if e == nil {
					// the search ends at the empty environment, which holds no bindings
					s.errorf(c.call, "cannot remove bindings from a locked environment")
				}
				target = e
			}
			switch env { // the base environments named by the call itself
			case s.BaseNS:
				s.errorf(c.call, "cannot remove variables from base namespace")
			case s.Base:
				s.errorf(c.call, "cannot remove variables from the base environment")
			}
			if target.frame().locked {
				s.errorf(c.call, "cannot remove bindings from a locked environment")
			}
			if !target.remove(sym) {
				s.warningf(c.call, "object '%s' not found", n)
			}
		}
		return s.invisible(Null)
	})
	regSpecial("local", func(s *Session, c *callInfo) Value {
		// local(expr, envir = new.env()): R runs eval(quote(expr), envir)
		var expr, envirExpr Value
		haveExpr, haveEnvir := false, false
		for _, ra := range c.raw {
			switch {
			case ra.HasName && ra.Name == "envir":
				envirExpr, haveEnvir = ra.Value, true
			case ra.HasName && ra.Name == "expr":
				expr, haveExpr = ra.Value, true
			}
		}
		for _, ra := range c.raw {
			if ra.HasName {
				continue
			}
			if !haveExpr {
				expr, haveExpr = ra.Value, true
			} else if !haveEnvir {
				envirExpr, haveEnvir = ra.Value, true
			}
		}
		quoted := &Lang{Fun: Intern("quote")}
		if haveExpr {
			quoted.Args = []Arg{{Value: expr}}
		}
		if !haveEnvir {
			envirExpr = &Lang{Fun: Intern("new.env")}
		}
		evalCall := &Lang{Fun: Intern("eval"), Args: []Arg{{Value: quoted}, {Value: envirExpr}}}
		if !haveExpr {
			s.errorf(evalCall, "argument is missing, with no default")
		}
		env := NewEnv(c.env)
		if haveEnvir {
			env = s.envArg(c, s.eval(envirExpr, c.env), "envir")
		}
		return s.underEval(evalCall, func() Value { return s.eval(expr, env) })
	})
	regSpecial("evalq", func(s *Session, c *callInfo) Value {
		if len(c.raw) == 0 {
			s.errorf(c.call, "argument \"expr\" is missing, with no default")
		}
		env := c.env
		if len(c.raw) > 1 {
			env = s.evalEnvArg(c, s.eval(c.raw[1].Value, c.env), nil)
		}
		return s.underEval(c.call, func() Value { return s.eval(c.raw[0].Value, env) })
	})
	regBuiltin("eval", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "expr", "envir", "enclos")
		if a[0] == nil {
			s.errorf(c.call, "argument \"expr\" is missing, with no default")
		}
		env := c.env
		if a[1] != nil {
			var enclos *Env
			if a[2] != nil {
				enclos = s.envArg(c, a[2], "enclos")
			}
			env = s.evalEnvArg(c, a[1], enclos)
		}
		switch x := a[0].(type) {
		case *Lang, *Symbol:
			return s.underEval(c.call, func() Value { return s.eval(x, env) })
		}
		s.visible = true
		return a[0]
	})
	regBuiltin("as.environment", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x")
		switch x := a[0].(type) {
		case *Env:
			return x
		case nullT:
			s.errorf(c.call, "using 'as.environment(NULL)' is defunct")
		case *List:
			e := NewEnv(s.Empty)
			s.fillEnv(c, e, x)
			return e
		case *Double, *Integer:
			if d := asDoubles(x); len(d) == 1 && d[0] == 1 {
				return s.Global
			}
		case *String:
			if len(x.V) == 1 {
				switch x.V[0].S {
				case ".GlobalEnv", "R_GlobalEnv":
					return s.Global
				case "package:base":
					return s.Base
				}
				s.errorf(c.call, "no item called \"%s\" on the search list", x.V[0].S)
			}
		}
		s.errorf(c.call, "invalid object for 'as.environment'")
		return nil
	})
	regBuiltin("list2env", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "envir", "parent")
		l, ok := a[0].(*List)
		if !ok {
			s.errorf(c.call, "first argument must be a named list")
		}
		var env *Env
		if a[1] != nil {
			env = s.envArg(c, a[1], "envir")
		} else {
			parent := c.env
			if a[2] != nil {
				parent = s.envArg(c, a[2], "parent")
			}
			env = NewEnv(parent)
		}
		s.fillEnv(c, env, l)
		return env
	})
	regBuiltin("mget", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "envir", "mode", "ifnotfound", "inherits")
		env := c.env
		if a[1] != nil {
			env = s.envArg(c, a[1], "envir")
		}
		st, ok := a[0].(*String)
		if !ok {
			s.errorf(c.call, "invalid first argument")
		}
		n := len(st.V)
		modes := []string{"any"}
		if a[2] != nil {
			ms, ok := a[2].(*String)
			if !ok || len(ms.V) == 0 {
				s.errorf(c.call, "invalid 'mode' argument")
			}
			modes = modes[:0]
			for _, m := range ms.V {
				modes = append(modes, m.S)
			}
		}
		var nf []Value
		if a[3] != nil {
			if l, isList := a[3].(*List); isList {
				nf = l.V
			} else {
				nf = s.asListElems(c, a[3])
			}
			if len(nf) != 1 && len(nf) != n {
				s.errorf(c.call, "wrong length for 'ifnotfound' argument")
			}
		}
		inherits := a[4] != nil && s.truthy(c, a[4], "inherits")
		out := &List{V: make([]Value, n)}
		if n > 0 {
			out.Attr = namesAttr(st.V)
		} else {
			out.Attr = []Attr{{Name: "names", Value: &String{V: []Str{}}}}
		}
		for i, name := range st.V {
			v, found := s.lookupMode(name.S, env, inherits, modes[i%len(modes)], true)
			if !found {
				if nf == nil {
					s.errorf(nil, "value for ‘%s’ not found", name.S)
				}
				v = nf[i%len(nf)]
				switch v.(type) {
				case *Closure, *Builtin:
					v = s.callFunction(v, []Value{mkString(name.S)})
				}
			}
			out.V[i] = v
		}
		return out
	})
	regBuiltin("body", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "fun")
		if f, ok := a[0].(*Closure); ok {
			return f.Body
		}
		return Null
	})
}

func (s *Session) fillEnv(c *callInfo, e *Env, l *List) {
	nm := namesOf(l)
	if len(l.V) > 0 && nm == nil {
		s.errorf(c.call, "names(x) must be a character vector of the same length as x")
	}
	for i, v := range l.V {
		if nm.V[i].S == "" || nm.V[i].NA {
			s.errorf(c.call, "attempt to use zero-length variable name")
		}
		s.assignTo(e, Intern(nm.V[i].S), v, c.call)
	}
}

// evalEnvArg turns the envir argument of eval into an environment.
func (s *Session) evalEnvArg(c *callInfo, v Value, enclos *Env) *Env {
	switch e := v.(type) {
	case *Env:
		return e
	case *List:
		if enclos == nil {
			enclos = c.env
		}
		env := NewEnv(enclos)
		if nm := namesOf(e); nm != nil {
			for i, x := range e.V {
				if nm.V[i].S != "" && !nm.V[i].NA {
					env.set(Intern(nm.V[i].S), x)
				}
			}
		}
		return env
	case nullT:
		if enclos == nil {
			enclos = c.env
		}
		return NewEnv(enclos)
	case *Double, *Integer:
		if Length(v) != 1 {
			s.errorf(c.call, "numeric 'envir' arg not of length one")
		}
	}
	s.errorf(c.call, "invalid 'envir' argument of type '%s'", typeName(v))
	return nil
}

func (s *Session) truthy(c *callInfo, v Value, what string) bool {
	l := logicalOf(s, c, v, "invalid '"+what+"' argument")
	if len(l) != 1 || l[0] == naInt {
		s.errorf(c.call, "invalid '%s' argument", what)
	}
	return l[0] != 0
}

// remove deletes sym from the frame; it reports whether it was bound.
func (e *Env) remove(sym *Symbol) bool {
	if e.alias != nil {
		return e.alias.remove(sym)
	}
	if e.lookup(sym) == nil {
		return false
	}
	if e.big != nil {
		delete(e.big, sym)
	} else {
		for i := range e.small {
			if e.small[i].sym == sym {
				e.small = append(e.small[:i], e.small[i+1:]...)
				break
			}
		}
	}
	for i, o := range e.order {
		if o == sym {
			e.order = append(e.order[:i], e.order[i+1:]...)
			break
		}
	}
	return true
}

var _ = syntax.Deparse
