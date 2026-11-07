package eval

import (
	"strings"
	"sync"
)

// dispatchInfo is what a method running under S3 dispatch knows.
type dispatchInfo struct {
	generic string
	classes []string // the dispatch classes, then "default"
	idx     int      // the class whose method is running
	object  Value
}

// frameOf is the frame of the function whose frame environment is env.
func (s *Session) frameOf(env *Env) *frame {
	for i := len(s.frames) - 1; i >= 0; i-- {
		if s.frames[i].env == env {
			return s.frames[i]
		}
	}
	return nil
}

// findMethod looks for the function named generic.class from env outwards.
func (s *Session) findMethod(name string, env *Env) Value {
	sym := Intern(name)
	for e := env; e != nil; e = e.parent {
		if b := e.lookup(sym); b != nil {
			v := b.val
			if p, ok := v.(*Promise); ok {
				v = s.force(p)
			}
			switch v.(type) {
			case *Closure, *Builtin:
				return v
			}
		}
	}
	return nil
}

func isObject(v Value) bool { return getAttr(v, "class") != nil }

// callMethod runs a method (a closure or a builtin) with supplied arguments.
func (s *Session) callMethod(m Value, name string, call *Lang, supplied []Arg, env *Env, di *dispatchInfo) Value {
	mc := &Lang{Fun: Intern(name), Args: call.Args}
	switch f := m.(type) {
	case *Closure:
		return s.applyClosure(f, mc, supplied, env, di)
	case *Builtin:
		ci := &callInfo{call: mc, env: env}
		for _, a := range supplied {
			if f.Special {
				ci.raw = append(ci.raw, a)
				continue
			}
			ci.args = append(ci.args, s.argValue(a.Value))
			ci.names = append(ci.names, a.Name)
		}
		return f.fn(s, ci)
	}
	return Null
}

func argsOf(args []Value, names []string) []Arg {
	out := make([]Arg, len(args))
	for i, a := range args {
		out[i] = Arg{Value: a}
		if names[i] != "" {
			out[i].Name, out[i].HasName = names[i], true
		}
	}
	return out
}

// dispatchInternal is S3 dispatch for a primitive generic ([, [[, length,
// ...): when the first argument has a class attribute and a method exists,
// the method runs.
func (s *Session) dispatchInternal(name string, c *callInfo, args []Value, names []string) (Value, bool) {
	if s.noDispatch {
		s.noDispatch = false
		return nil, false
	}
	if len(args) == 0 || !isObject(args[0]) {
		return nil, false
	}
	classes := classAttr(args[0])
	for i, cl := range classes {
		if m := s.findMethod(name+"."+cl, c.env); m != nil {
			di := &dispatchInfo{generic: name, classes: append(append([]string{}, classes...), "default"), idx: i, object: args[0]}
			return s.callMethod(m, name+"."+cl, c.call, argsOf(args, names), c.env, di), true
		}
	}
	return nil, false
}

// group generics
var (
	opsGroup  = map[string]bool{"+": true, "-": true, "*": true, "/": true, "^": true, "%%": true, "%/%": true, "==": true, "!=": true, "<": true, "<=": true, ">": true, ">=": true, "&": true, "|": true, "!": true}
	mathGroup = map[string]bool{"abs": true, "sqrt": true, "exp": true, "floor": true, "ceiling": true,
		"sin": true, "cos": true, "tan": true, "asin": true, "acos": true, "atan": true, "sinh": true, "cosh": true, "tanh": true,
		"asinh": true, "acosh": true, "atanh": true, "log": true, "log10": true, "log2": true, "log1p": true, "expm1": true,
		"gamma": true, "lgamma": true, "sinpi": true, "cospi": true, "tanpi": true, "trunc": true, "sign": true, "round": true, "signif": true,
		"cumsum": true, "cumprod": true, "cummax": true, "cummin": true}
	summaryGroup = map[string]bool{"sum": true, "prod": true, "min": true, "max": true, "any": true, "all": true, "range": true}
	internalGen  = []string{"length", "names", "names<-", "as.character", "as.integer", "as.double", "as.numeric", "as.logical", "c", "is.na", "dim", "dim<-", "unlist", "rep", "seq_along", "mean", "t", "as.vector"}
)

// methodFor finds the method of generic (or of its group) for one class vector.
func (s *Session) methodFor(op, group string, classes []string, env *Env) (Value, string, int) {
	for i, cl := range classes {
		if m := s.findMethod(op+"."+cl, env); m != nil {
			return m, op + "." + cl, i
		}
		if m := s.findMethod(group+"."+cl, env); m != nil {
			return m, group + "." + cl, i
		}
	}
	return nil, "", 0
}

func (s *Session) dispatchGroup(op, group string, c *callInfo) (Value, bool) {
	if s.noDispatch {
		s.noDispatch = false
		return nil, false
	}
	if len(c.args) == 0 {
		return nil, false
	}
	var cands []Value
	if group == "Ops" {
		for _, a := range c.args {
			if isObject(a) {
				cands = append(cands, a)
			}
		}
	} else if isObject(c.args[0]) {
		cands = []Value{c.args[0]}
	}
	if len(cands) == 0 {
		return nil, false
	}
	var found Value
	var fname string
	var di *dispatchInfo
	for k, a := range cands {
		classes := classAttr(a)
		m, name, i := s.methodFor(op, group, classes, c.env)
		if m == nil {
			continue
		}
		if found == nil {
			found, fname = m, name
			di = &dispatchInfo{generic: op, classes: append(append([]string{}, classes...), "default"), idx: i, object: a}
		} else if name != fname && k > 0 {
			s.warningf(c.call, "Incompatible methods (\"%s\", \"%s\") for \"%s\"", fname, name, op)
			return nil, false
		}
	}
	if found == nil {
		return nil, false
	}
	return s.callMethod(found, fname, c.call, argsOf(c.args, c.names), c.env, di), true
}

var s3once sync.Once

// installS3 wraps the primitives that dispatch.
func installS3() {
	s3once.Do(func() {
		for name, b := range builtins {
			name, b := name, b
			orig := b.fn
			switch {
			case opsGroup[name]:
				b.fn = func(s *Session, c *callInfo) Value {
					if r, ok := s.dispatchGroup(name, "Ops", c); ok {
						return r
					}
					return orig(s, c)
				}
			case mathGroup[name]:
				b.fn = func(s *Session, c *callInfo) Value {
					if r, ok := s.dispatchGroup(name, "Math", c); ok {
						return r
					}
					return orig(s, c)
				}
			case summaryGroup[name]:
				b.fn = func(s *Session, c *callInfo) Value {
					if r, ok := s.dispatchGroup(name, "Summary", c); ok {
						return r
					}
					return orig(s, c)
				}
			}
		}
		for _, name := range internalGen {
			name := name
			b, ok := builtins[name]
			if !ok {
				continue
			}
			orig := b.fn
			b.fn = func(s *Session, c *callInfo) Value {
				if r, ok := s.dispatchInternal(name, c, c.args, c.names); ok {
					return r
				}
				return orig(s, c)
			}
		}
	})
}

func (s *Session) implicitDispatchClasses(v Value) []string {
	return append(dispatchClass(v), "default")
}

func init() {
	regSpecial("UseMethod", func(s *Session, c *callInfo) Value {
		if len(c.raw) == 0 || len(c.raw) > 2 {
			s.errorf(c.call, "there must be a 'generic' argument")
		}
		g := s.eval(c.raw[0].Value, c.env)
		gs, ok := g.(*String)
		if !ok || len(gs.V) != 1 {
			s.errorf(c.call, "'generic' argument must be a character string")
		}
		generic := gs.V[0].S
		fr := s.frameOf(c.env)
		if fr == nil {
			s.errorf(c.call, "UseMethod called from outside a function")
		}
		var obj Value = Null
		if len(c.raw) == 2 {
			obj = s.eval(c.raw[1].Value, c.env)
		} else if len(fr.fn.Formals) > 0 {
			f0 := fr.fn.Formals[0].Name
			if f0.Name == "..." {
				if d, _ := c.env.lookup(f0).val.(*Dots); d != nil && len(d.Args) > 0 {
					obj = s.argValue(d.Args[0].Value)
				}
			} else if b := c.env.lookup(f0); b != nil {
				switch bv := b.val.(type) {
				case missingT:
				case *Promise:
					obj = s.force(bv)
				default:
					obj = bv
				}
			}
		}
		classes := s.implicitDispatchClasses(obj)
		callEnv := fr.caller
		if callEnv == nil {
			callEnv = c.env
		}
		for i, cl := range classes {
			name := generic + "." + cl
			m := s.findMethod(name, callEnv)
			if m == nil {
				m = s.findMethod(name, fr.fn.Env)
			}
			if m == nil {
				continue
			}
			di := &dispatchInfo{generic: generic, classes: classes, idx: i, object: obj}
			v := s.callMethod(m, name, fr.call, fr.supplied, callEnv, di)
			s.flow, s.flowVal, s.flowEnv = flowReturn, v, c.env
			return v
		}
		cls := classOf(obj)
		desc := cls[0]
		if len(cls) > 1 {
			parts := make([]string, len(cls))
			for i, x := range cls {
				parts[i] = "'" + x + "'"
			}
			desc = "c(" + strings.Join(parts, ", ") + ")"
		}
		if classAttr(obj) == nil {
			dc := dispatchClass(obj)
			if len(dc) > 1 {
				parts := make([]string, len(dc))
				for i, x := range dc {
					parts[i] = "'" + x + "'"
				}
				desc = "c(" + strings.Join(parts, ", ") + ")"
			}
		}
		s.errorf(c.call, "no applicable method for '%s' applied to an object of class \"%s\"", generic, desc)
		return nil
	})
	regBuiltin("NextMethod", func(s *Session, c *callInfo) Value {
		fr := s.frameOf(c.env)
		if fr == nil {
			s.errorf(c.call, "NextMethod called from outside a method dispatch")
		}
		di := fr.di
		if di == nil {
			s.errorf(c.call, "generic function not specified")
		}
		// the arguments: the formals of the running method, as promises to
		// their current values, then those given to NextMethod
		var supplied []Arg
		matched := map[string]bool{}
		named := map[string]bool{}
		for _, a := range fr.supplied {
			if a.HasName && a.Name != "" {
				named[a.Name] = true
			}
		}
		for _, fm := range fr.fn.Formals {
			if fm.Name.Name == "..." {
				if d, ok := c.env.lookup(fm.Name).val.(*Dots); ok {
					supplied = append(supplied, d.Args...)
				}
				continue
			}
			if c.env.missing[fm.Name] {
				continue
			}
			arg := Arg{Value: &Promise{Expr: fm.Name, Env: c.env}}
			if named[fm.Name.Name] {
				arg.Name, arg.HasName = fm.Name.Name, true
			}
			supplied = append(supplied, arg)
			matched[fm.Name.Name] = true
		}
		// extra arguments of NextMethod(generic, object, ...)
		for i, a := range c.args {
			if i < 2 && c.names[i] == "" {
				continue
			}
			if c.names[i] == "generic" || c.names[i] == "object" {
				continue
			}
			na := Arg{Value: a}
			if c.names[i] != "" {
				na.Name, na.HasName = c.names[i], true
				replaced := false
				for j := range supplied {
					if supplied[j].HasName && supplied[j].Name == na.Name {
						supplied[j] = na
						replaced = true
					}
				}
				if replaced {
					continue
				}
			}
			supplied = append(supplied, na)
		}
		for j := di.idx + 1; j < len(di.classes); j++ {
			name := di.generic + "." + di.classes[j]
			if m := s.findMethod(name, c.env); m != nil {
				ndi := &dispatchInfo{generic: di.generic, classes: di.classes, idx: j, object: di.object}
				return s.callMethod(m, name, fr.call, supplied, c.env, ndi)
			}
		}
		// no method left: the internal default
		if b, ok := builtins[di.generic]; ok {
			s.noDispatch = true
			ci := &callInfo{call: fr.call, env: c.env}
			for _, a := range supplied {
				ci.args = append(ci.args, s.argValue(a.Value))
				ci.names = append(ci.names, a.Name)
			}
			if b.Special {
				ci.raw = supplied
			}
			v := b.fn(s, ci)
			s.noDispatch = false
			return v
		}
		s.errorf(c.call, "no more methods for '%s'", di.generic)
		return nil
	})
}
