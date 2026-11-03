package eval

import (
	"context"
	"time"
)

func init() {
	regSpecial("quote", func(s *Session, c *callInfo) Value {
		if len(c.raw) != 1 {
			s.errorf(c.call, "%d arguments passed to 'quote' which requires 1", len(c.raw))
		}
		return c.raw[0].Value
	})
	regSpecial("function", func(s *Session, c *callInfo) Value {
		f := c.raw[0].Value.(*Formals)
		return &Closure{Formals: f.F, Body: c.raw[1].Value, Env: c.env}
	})
	regSpecial("{", func(s *Session, c *callInfo) Value {
		var v Value = Null
		for _, a := range c.raw {
			v = s.evalStmt(a.Value, c.env)
			if s.flow != flowNone {
				return v
			}
		}
		return v
	})
	regBuiltin("(", func(s *Session, c *callInfo) Value {
		s.visible = true
		if len(c.args) != 1 {
			s.errorf(c.call, "%d arguments passed to '(' which requires 1", len(c.args))
		}
		return c.args[0]
	})
	regBuiltin("invisible", func(s *Session, c *callInfo) Value {
		s.visible = false
		switch len(c.args) {
		case 0:
			return Null
		case 1:
			return c.args[0]
		}
		s.errorf(c.call, "%d arguments passed to 'invisible' which requires 1", len(c.args))
		return nil
	})
	regSpecial("if", func(s *Session, c *callInfo) Value {
		if len(c.raw) < 2 {
			s.errorf(c.call, "argument \"%s\" is missing, with no default", "yes")
		}
		cond := s.eval(c.raw[0].Value, c.env)
		if s.asCondition(cond, c.call) {
			return s.evalStmt(c.raw[1].Value, c.env)
		}
		if len(c.raw) > 2 {
			return s.evalStmt(c.raw[2].Value, c.env)
		}
		s.visible = false
		return Null
	})
	regSpecial("while", func(s *Session, c *callInfo) Value {
		for {
			cond := s.eval(c.raw[0].Value, c.env)
			if !s.asCondition(cond, c.call) {
				break
			}
			if s.loopBody(c.raw[1].Value, c.env) || s.flow != flowNone {
				break
			}
		}
		s.endLoop()
		return Null
	})
	regSpecial("repeat", func(s *Session, c *callInfo) Value {
		for {
			if s.loopBody(c.raw[0].Value, c.env) || s.flow != flowNone {
				break
			}
		}
		s.endLoop()
		return Null
	})
	regSpecial("for", func(s *Session, c *callInfo) Value {
		sym, ok := c.raw[0].Value.(*Symbol)
		if !ok {
			s.errorf(c.call, "non-symbol loop variable")
		}
		seq := s.eval(c.raw[1].Value, c.env)
		n := Length(seq)
		switch seq.(type) {
		case nullT, *Logical, *Integer, *Double, *String, *List:
		default:
			s.errorf(c.call, "invalid for() loop sequence")
		}
		for i := 0; i < n; i++ {
			c.env.set(sym, element(seq, i))
			if s.loopBody(c.raw[2].Value, c.env) || s.flow != flowNone {
				break
			}
		}
		s.endLoop()
		return Null
	})
	regSpecial("break", func(s *Session, c *callInfo) Value {
		s.flow = flowBreak
		return Null
	})
	regSpecial("next", func(s *Session, c *callInfo) Value {
		s.flow = flowNext
		return Null
	})
	regSpecial("return", func(s *Session, c *callInfo) Value {
		var v Value = Null
		switch len(c.raw) {
		case 0:
		case 1:
			v = s.eval(c.raw[0].Value, c.env)
		default:
			s.errorf(c.call, "multi-argument returns are not permitted")
		}
		if c.env.fn == nil {
			s.errorf(c.call, "no function to return from, jumping to top level")
		}
		s.flow, s.flowVal, s.flowEnv = flowReturn, v, c.env
		return v
	})
	assign := func(super bool) func(s *Session, c *callInfo) Value {
		return func(s *Session, c *callInfo) Value {
			if len(c.raw) != 2 {
				s.errorf(c.call, "invalid assignment")
			}
			var sym *Symbol
			switch t := c.raw[0].Value.(type) {
			case *Symbol:
				sym = t
			case *String:
				if len(t.V) == 1 && !t.V[0].NA {
					sym = Intern(t.V[0].S)
				}
			case *Lang:
				v := s.eval(c.raw[1].Value, c.env)
				s.complexAssign(c, t, v, super)
				s.visible = false
				return v
			}
			if sym == nil {
				s.errorf(c.call, "invalid (do_set) left-hand side to assignment")
			}
			v := s.eval(c.raw[1].Value, c.env)
			if super {
				target := s.Global
				for e := c.env.parent; e != nil; e = e.parent {
					if e.lookup(sym) != nil {
						target = e
						break
					}
				}
				target.set(sym, v)
			} else {
				c.env.set(sym, v)
			}
			s.visible = false
			return v
		}
	}
	regSpecial("<-", assign(false))
	regSpecial("=", assign(false))
	regSpecial("<<-", assign(true))
	regSpecial("missing", func(s *Session, c *callInfo) Value {
		if len(c.raw) != 1 {
			s.errorf(c.call, "%d arguments passed to 'missing' which requires 1", len(c.raw))
		}
		var sym *Symbol
		switch t := c.raw[0].Value.(type) {
		case *Symbol:
			sym = t
		case *String:
			if len(t.V) == 1 {
				sym = Intern(t.V[0].S)
			}
		}
		if sym == nil {
			s.errorf(c.call, "invalid use of 'missing'")
		}
		if c.env.fn == nil {
			s.errorf(c.call, "'missing' can only be used for arguments")
		}
		if c.env.lookup(sym) == nil {
			s.errorf(c.call, "'missing' can only be used for arguments")
		}
		return mkLogical(s.isMissingArg(sym, c.env))
	})
}

// checkInterrupt ends the evaluation when its context is done or its time
// is up.
func (s *Session) checkInterrupt() {
	if s.ctx != nil {
		if err := s.ctx.Err(); err != nil {
			s.abort("evaluation interrupted: "+err.Error(), err)
		}
	}
	if !s.deadline.IsZero() && time.Now().After(s.deadline) {
		s.abort("evaluation interrupted: time limit of "+s.opts.Timeout.String()+" exceeded", context.DeadlineExceeded)
	}
}

// loopBody evaluates the body of a loop once; it reports whether the loop
// must end because of a break. A next ends the iteration only.
func (s *Session) loopBody(body Value, env *Env) (brk bool) {
	defer func() {
		if r := recover(); r != nil {
			if sig, ok := r.(*flowSignal); ok && (sig.f == flowBreak || sig.f == flowNext) {
				brk = sig.f == flowBreak
				return
			}
			panic(r)
		}
	}()
	s.tick++
	if s.tick&127 == 0 {
		s.checkInterrupt()
	}
	s.evalStmt(body, env)
	switch s.flow {
	case flowBreak:
		s.flow = flowNone
		return true
	case flowNext:
		s.flow = flowNone
	}
	return false
}

// element is x[[i]] stripped of attributes, the loop variable's value.
func element(x Value, i int) Value {
	switch v := x.(type) {
	case *Logical:
		return &Logical{V: []int32{v.V[i]}}
	case *Integer:
		return &Integer{V: []int32{v.V[i]}}
	case *Double:
		return &Double{V: []float64{v.V[i]}}
	case *String:
		return &String{V: []Str{v.V[i]}}
	case *List:
		return v.V[i]
	}
	return Null
}

// endLoop makes the value of a loop invisible, unless it ended by a return,
// whose visibility is that of the returned value.
func (s *Session) endLoop() {
	if s.flow == flowNone {
		s.visible = false
	}
}
