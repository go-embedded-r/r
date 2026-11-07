package eval

// The call stack: sys.nframe, sys.calls, sys.frames, sys.frame, sys.on.exit,
// Recall, match.call, geterrmessage.

func (s *Session) stackIndex(c *callInfo, which Value, def int) int {
	cur, _ := s.frameIndex(c.env)
	w := def
	if which != nil {
		w = s.intArg(c, which, "which")
	}
	switch {
	case w > 0:
		return w - 1
	case w < 0:
		return cur + w
	}
	return -1
}

func init() {
	regBuiltin("sys.nframe", func(s *Session, c *callInfo) Value {
		cur, fr := s.frameIndex(c.env)
		if fr == nil {
			return mkInt(0)
		}
		return mkInt(int32(cur + 1))
	})
	regBuiltin("sys.calls", func(s *Session, c *callInfo) Value {
		cur, fr := s.frameIndex(c.env)
		if fr == nil {
			return Null
		}
		out := make([]Value, 0, cur+1)
		for i := 0; i <= cur; i++ {
			out = append(out, s.frames[i].call)
		}
		return &List{V: out}
	})
	regBuiltin("sys.frames", func(s *Session, c *callInfo) Value {
		cur, fr := s.frameIndex(c.env)
		if fr == nil {
			return Null
		}
		out := make([]Value, 0, cur+1)
		for i := 0; i <= cur; i++ {
			out = append(out, s.frames[i].env)
		}
		return &List{V: out}
	})
	regBuiltin("sys.frame", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "which")
		idx := s.stackIndex(c, a[0], 0)
		if idx == -1 && (a[0] == nil || s.intArg(c, a[0], "which") == 0) {
			return s.Global
		}
		if idx < 0 || idx >= len(s.frames) {
			s.errorf(c.call, "not that many frames on the stack")
		}
		return s.frames[idx].env
	})
	regBuiltin("sys.on.exit", func(s *Session, c *callInfo) Value {
		_, fr := s.frameIndex(c.env)
		if fr == nil || len(fr.onexit) == 0 {
			return Null
		}
		if len(fr.onexit) == 1 {
			return fr.onexit[0]
		}
		l := &Lang{Fun: Intern("{")}
		for _, e := range fr.onexit {
			l.Args = append(l.Args, Arg{Value: e})
		}
		return l
	})
	regBuiltin("Recall", func(s *Session, c *callInfo) Value {
		_, fr := s.frameIndex(c.env)
		if fr == nil {
			s.errorf(c.call, "Recall called from outside a closure")
		}
		args := make([]Arg, len(c.args))
		for i, v := range c.args {
			args[i] = Arg{Name: c.names[i], HasName: c.names[i] != "", Value: v}
		}
		return s.applyClosure(fr.fn, fr.call, args, c.env, nil)
	})
	regBuiltin("geterrmessage", func(s *Session, c *callInfo) Value { return mkString(s.lastErr) })
	regBuiltin("match.call", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "definition", "call", "expand.dots", "envir")
		_, fr := s.frameIndex(c.env)
		var fn *Closure
		var call *Lang
		var callerEnv *Env
		if fr != nil {
			fn, call, callerEnv = fr.fn, fr.call, fr.caller
		}
		if a[0] != nil {
			if f, ok := a[0].(*Closure); ok {
				fn = f
			} else if _, isNull := a[0].(nullT); !isNull {
				s.errorf(c.call, "invalid 'definition' argument")
			}
		}
		if a[1] != nil {
			l, ok := a[1].(*Lang)
			if !ok {
				s.errorf(c.call, "invalid 'call' argument")
			}
			call = l
			callerEnv = nil
		}
		if fn == nil || call == nil {
			s.errorf(c.call, "match.call() was called from outside a function")
		}
		expand := a[2] == nil || s.truthy(c, a[2], "expand.dots")
		// the supplied arguments as written, `...` replaced by what it holds
		var supplied []Arg
		for _, ar := range call.Args {
			if sym, ok := ar.Value.(*Symbol); ok && sym.Name == "..." {
				if callerEnv != nil {
					if d := s.dotsOf(callerEnv); d != nil {
						for _, da := range d.Args {
							v := da.Value
							if p, ok := v.(*Promise); ok {
								v = p.Expr
							}
							supplied = append(supplied, Arg{Name: da.Name, HasName: da.Name != "", Value: v})
						}
					}
				}
				continue
			}
			supplied = append(supplied, ar)
		}
		tmp := &Env{}
		s.bindArgs(fn, call, supplied, tmp)
		out := &Lang{Fun: call.Fun}
		for _, fm := range fn.Formals {
			if fm.Name.Name == "..." {
				b := tmp.lookup(fm.Name)
				d, _ := b.val.(*Dots)
				if d == nil || len(d.Args) == 0 {
					continue
				}
				if expand {
					for _, da := range d.Args {
						out.Args = append(out.Args, Arg{Name: da.Name, HasName: da.Name != "", Value: da.Value})
					}
				} else {
					l := &List{V: make([]Value, len(d.Args))}
					anyName := false
					names := make([]Str, len(d.Args))
					for i, da := range d.Args {
						l.V[i] = da.Value
						names[i] = Str{S: da.Name}
						if da.Name != "" {
							anyName = true
						}
					}
					if anyName {
						l.Attr = namesAttr(names)
					}
					out.Args = append(out.Args, Arg{Name: "...", HasName: true, Value: l})
				}
				continue
			}
			if tmp.missing[fm.Name] {
				continue
			}
			b := tmp.lookup(fm.Name)
			if b == nil {
				continue
			}
			out.Args = append(out.Args, Arg{Name: fm.Name.Name, HasName: true, Value: b.val})
		}
		return out
	})
}
