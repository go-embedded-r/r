package eval

import "sort"

// A minimal namespace model. The evaluator has no package loading: the one
// namespace is that of base, an environment that is another name for the base
// frame (as in R, where it is a distinct environment whose enclosure is the
// global one). `::` and `:::` resolve in it; the packages R installs with
// base are known by name and resolve through our base environment, which
// holds what we implement of them; any other package does not exist.

// installedPackages are the packages the oracle's R ships with.
var installedPackages = map[string]bool{
	"compiler": true, "datasets": true, "graphics": true, "grDevices": true,
	"grid": true, "methods": true, "parallel": true, "splines": true,
	"stats": true, "stats4": true, "tcltk": true, "tools": true, "utils": true,
}

// packageNotFound raises R's packageNotFoundError text, with the call R shows.
func (s *Session) packageNotFound(pkg, argName string) {
	call := &Lang{Fun: Intern("loadNamespace"), Args: []Arg{{Value: Intern(argName)}}}
	s.errorf(call, "there is no package called ‘%s’", pkg)
}

// namespaceNamed is the namespace called name. Packages other than base are
// not available: unknown ones are an error as in R, the installed ones an
// error that says why.
func (s *Session) namespaceNamed(c *callInfo, name, argName string) *Env {
	if name == "base" {
		return s.BaseNS
	}
	if installedPackages[name] {
		s.errorf(c.call, "namespace '%s' is not supported by this evaluator (no package loading)", name)
	}
	s.packageNotFound(name, argName)
	return nil
}

// nsName reads a namespace or package name given as a string or a symbol.
func nsName(v Value) (string, bool) {
	switch x := v.(type) {
	case *Symbol:
		return x.Name, true
	case *String:
		if len(x.V) >= 1 && !x.V[0].NA {
			return x.V[0].S, true
		}
	}
	return "", false
}

// asNamespace is asNamespace(ns) for the call c: an environment must be a
// namespace, a name must be that of an available one.
func (s *Session) asNamespace(c *callInfo, v Value) *Env {
	switch x := v.(type) {
	case *Env:
		if x.ns {
			return x
		}
	case *String, *Symbol:
		if name, ok := nsName(x); ok {
			return s.namespaceNamed(c, name, "name")
		}
	}
	s.errorf(c.call, "not a namespace")
	return nil
}

// topEnv walks from e outwards to the first environment that is a top level
// one: the global environment, the base environment or a namespace.
func (s *Session) topEnv(e *Env) *Env {
	for ; e != nil; e = e.parent {
		if e == s.Global || e == s.Base || e.ns {
			return e
		}
	}
	return s.Global
}

// modeArg reads the mode argument of get, exists and mget: "any" when absent.
func (s *Session) modeArg(c *callInfo, v Value) string {
	if v == nil {
		return "any"
	}
	st, ok := v.(*String)
	if !ok || len(st.V) < 1 || st.V[0].NA {
		s.errorf(c.call, "invalid 'mode' argument")
	}
	return st.V[0].S // the first of several, as R reads it
}

// nameArg reads the name argument of get, exists, get0 and assign.
func (s *Session) nameArg(c *callInfo, v Value) string {
	st, ok := v.(*String)
	if !ok || len(st.V) == 0 {
		s.errorf(c.call, "invalid first argument")
	}
	if len(st.V) > 1 {
		s.errorf(c.call, "first argument has length > 1")
	}
	return st.V[0].S
}

// modeOf is mode(v) for the values a binding can hold.
func modeOf(v Value) string {
	switch v.(type) {
	case *Integer, *Double:
		return "numeric"
	case *String:
		return "character"
	case *Logical:
		return "logical"
	case *List:
		return "list"
	case *Closure, *Builtin:
		return "function"
	case *Env:
		return "environment"
	case *Symbol:
		return "name"
	case *Lang:
		return "call"
	case nullT:
		return "NULL"
	}
	return "any"
}

// lookupMode finds name from env (through the enclosures when inherits) among
// the bindings whose value has the given mode. A promise is forced only when
// force is set or the mode needs the value.
func (s *Session) lookupMode(name string, env *Env, inherits bool, mode string, force bool) (Value, bool) {
	sym := Intern(name)
	for e := env; e != nil; e = e.parent {
		if b := e.lookup(sym); b != nil {
			if mode == "any" && !force && b.active == nil {
				if _, miss := b.val.(missingT); !miss {
					return b.val, true
				}
			} else {
				v := s.bforce(b)
				if _, miss := v.(missingT); !miss && (mode == "any" || modeOf(v) == mode) {
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

func init() {
	regBuiltin("environmentName", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "env")
		e, ok := a[0].(*Env)
		if !ok {
			return mkString("")
		}
		if n := envName(s, e); n != "" {
			return mkString(n)
		}
		if nm := getAttr(e, "name"); nm != nil {
			return nm // whatever the attribute holds, as R answers
		}
		return mkString("")
	})
	regBuiltin("topenv", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "envir", "matchThisEnv")
		start := c.env
		if a[0] != nil {
			e, ok := a[0].(*Env)
			if !ok {
				return s.Global
			}
			start = e
		}
		if m, ok := a[1].(*Env); ok {
			for e := start; e != nil; e = e.parent {
				if e == m {
					return e
				}
			}
		}
		return s.topEnv(start)
	})
	regBuiltin("asNamespace", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "ns", "base.OK")
		if a[0] == nil {
			s.errorf(c.call, "argument \"ns\" is missing, with no default")
		}
		return s.asNamespace(c, a[0])
	})
	regBuiltin("getNamespace", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "name")
		switch x := a[0].(type) {
		case *Env:
			s.errorf(nil, "environments cannot be coerced to other types")
		case *String, *Symbol:
			if n, ok := nsName(x); ok {
				return s.namespaceNamed(c, n, "name")
			}
		case *Integer, *Double, *Logical:
			if st := asStrings(x); len(st) > 0 {
				return s.namespaceNamed(c, st[0].S, "name")
			}
		}
		s.errorf(c.call, "invalid 'name' argument")
		return nil
	})
	regBuiltin("isNamespace", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "ns")
		e, ok := a[0].(*Env)
		return mkLogical(ok && e.ns)
	})
	regBuiltin("isBaseNamespace", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "ns")
		return mkLogical(a[0] == Value(s.BaseNS))
	})
	regBuiltin("getNamespaceName", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "ns")
		var ns *Env
		switch x := a[0].(type) {
		case *String, *Symbol:
			if n, ok := nsName(x); ok {
				ns = s.namespaceNamed(c, n, "name")
			}
		case *Env:
			if x.ns {
				ns = x
			}
		}
		if ns == nil {
			s.errorf(&Lang{Fun: Intern("asNamespace"), Args: []Arg{{Value: Intern("ns")}}}, "not a namespace")
		}
		return mkString(ns.nsName())
	})
	regBuiltin("getNamespaceExports", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "ns")
		ns := s.asNamespace(c, a[0])
		names := ns.Names()
		out := make([]string, len(names))
		copy(out, names)
		sort.Strings(out)
		return strVec(out)
	})
	regBuiltin("isNamespaceLoaded", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "name")
		n, ok := nsName(a[0])
		if !ok {
			s.errorf(c.call, "bad namespace name")
		}
		return mkLogical(n == "base")
	})
	regBuiltin("requireNamespace", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "package", "quietly")
		n, ok := nsName(a[0])
		if !ok {
			s.errorf(c.call, "invalid 'package' argument")
		}
		if n == "base" {
			return s.invisible(mkLogical(true))
		}
		quiet := a[1] != nil && logicalFlag(a[1])
		if !quiet {
			s.stderr.WriteString("Loading required namespace: " + n + "\n")
		}
		if installedPackages[n] {
			if !quiet {
				s.stderr.WriteString("Failed with error:  ‘namespace '" + n + "' is not supported by this evaluator (no package loading)’\n")
			}
		} else if !quiet {
			s.stderr.WriteString("Failed with error:  ‘there is no package called ‘" + n + "’’\n")
		}
		return s.invisible(mkLogical(false))
	})
	for _, op := range []string{"::", ":::"} {
		exported := op == "::"
		regSpecial(op, func(s *Session, c *callInfo) Value {
			if len(c.raw) != 2 {
				s.errorf(c.call, "bad arguments")
			}
			pkg, ok1 := nsName(c.raw[0].Value)
			name, ok2 := nsName(c.raw[1].Value)
			if !ok1 || !ok2 {
				s.errorf(c.call, "bad arguments")
			}
			if pkg != "base" && !installedPackages[pkg] {
				s.packageNotFound(pkg, "x")
			}
			// the installed packages resolve through what base holds of them
			if b := s.Base.lookup(Intern(name)); b != nil {
				return s.bforce(b)
			}
			switch {
			case pkg == "base":
				s.errorf(s.currentCall(), "object '%s' not found", name)
			case exported:
				s.errorf(c.call, "'%s' is not an exported object from 'namespace:%s'", name, pkg)
			}
			s.errorf(s.currentCall(), "object '%s' not found", name)
			return nil
		})
	}
}

// nsName is the name of a namespace environment.
func (e *Env) nsName() string { return "base" }
