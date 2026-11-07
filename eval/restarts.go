package eval

// Restarts: withRestarts, computeRestarts, invokeRestart, tryInvokeRestart.
// A restart is established with a name and, for withRestarts, a function;
// invoking it unwinds to its withRestarts and returns what the function
// returns. The muffling restarts of warning() and message() have no function.

type restartEntry struct {
	name string
	id   int
	fn   Value // nil for the built-in muffling restarts
	cond Value
}

func restartObject(e restartEntry) Value {
	l := &List{V: []Value{mkString(e.name), Null}}
	l.Attr = []Attr{{Name: "names", Value: strVec([]string{"name", "exit"})}, {Name: "class", Value: mkString("restart")}}
	return l
}

// findRestart looks a restart up by name (innermost first) or by object.
func (s *Session) findRestart(v Value) (restartEntry, bool) {
	name := ""
	switch x := v.(type) {
	case *String:
		if len(x.V) == 1 {
			name = x.V[0].S
		}
	case *List:
		if inheritsClass(x, "restart") && len(x.V) > 0 {
			if st, ok := x.V[0].(*String); ok && len(st.V) == 1 {
				name = st.V[0].S
			}
		}
	}
	for i := len(s.restarts) - 1; i >= 0; i-- {
		if s.restarts[i].name == name {
			return s.restarts[i], true
		}
	}
	return restartEntry{}, false
}

func init() {
	invoke := func(try bool) func(s *Session, c *callInfo) Value {
		return func(s *Session, c *callInfo) Value {
			if len(c.args) == 0 {
				s.errorf(c.call, "argument \"r\" is missing, with no default")
			}
			e, ok := s.findRestart(c.args[0])
			if !ok {
				if try {
					s.visible = false
					return Null
				}
				if st, isStr := c.args[0].(*String); isStr && len(st.V) == 1 {
					s.errorf(c.call, "no 'restart' '%s' found", st.V[0].S)
				}
				s.errorf(c.call, "bad restart specification")
			}
			panic(&restartSignal{name: e.name, id: e.id, args: append([]Value(nil), c.args[1:]...)})
		}
	}
	regBuiltin("invokeRestart", invoke(false))
	regBuiltin("tryInvokeRestart", invoke(true))
	regBuiltin("computeRestarts", func(s *Session, c *callInfo) Value {
		out := make([]Value, 0, len(s.restarts))
		for i := len(s.restarts) - 1; i >= 0; i-- {
			out = append(out, restartObject(s.restarts[i]))
		}
		return &List{V: out}
	})
	regBuiltin("findRestart", func(s *Session, c *callInfo) Value {
		if len(c.args) == 0 {
			s.errorf(c.call, "argument \"name\" is missing, with no default")
		}
		if e, ok := s.findRestart(c.args[0]); ok {
			return restartObject(e)
		}
		return Null
	})
	regSpecial("withRestarts", func(s *Session, c *callInfo) Value {
		var expr Value
		haveExpr := false
		type spec struct {
			name string
			fn   Value
		}
		var specs []spec
		for _, a := range c.raw {
			switch {
			case a.HasName && a.Name == "expr":
				expr, haveExpr = a.Value, true
			case !a.HasName && !haveExpr:
				expr, haveExpr = a.Value, true
			case a.HasName:
				specs = append(specs, spec{a.Name, s.eval(a.Value, c.env)})
			}
		}
		if !haveExpr {
			s.errorf(c.call, "argument \"expr\" is missing, with no default")
		}
		base := len(s.restarts)
		baseHandlers := s.handlers
		depth, nframes := s.depth, len(s.frames)
		var ids []int
		for i := len(specs) - 1; i >= 0; i-- {
			s.restartID++
			ids = append(ids, s.restartID)
			s.restarts = append(s.restarts, restartEntry{name: specs[i].name, id: s.restartID, fn: specs[i].fn})
		}
		var result Value
		func() {
			defer func() {
				s.restarts = s.restarts[:base]
				if r := recover(); r != nil {
					if rs, ok := r.(*restartSignal); ok {
						for k, id := range ids {
							if rs.id == id {
								fn := specs[len(specs)-1-k].fn
								s.handlers = baseHandlers
								s.depth, s.frames = depth, s.frames[:nframes]
								result = s.callFunction(fn, rs.args)
								return
							}
						}
					}
					panic(r)
				}
			}()
			result = s.eval(expr, c.env)
		}()
		return result
	})
}
