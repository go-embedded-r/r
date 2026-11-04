package eval

import (
	"strings"

	"github.com/go-embedded-r/r/syntax"
)

// A handler of the condition system: tryCatch's are exiting, those of
// withCallingHandlers are calling.
type handler struct {
	class   string
	fn      Value
	calling bool
	id      int
}

// unwindSignal is the panic that unwinds to the tryCatch with this id.
type unwindSignal struct {
	id      int
	cond    Value
	handler Value
}

// restartSignal is the panic of invokeRestart.
type restartSignal struct{ name string }

func callValue(c *Lang) Value {
	if c == nil {
		return Null
	}
	return c
}

func makeCondition(msg string, call Value, classes ...string) Value {
	l := &List{V: []Value{mkString(msg), call}}
	l.Attr = []Attr{{Name: "names", Value: strVec([]string{"message", "call"})}, {Name: "class", Value: strVec(classes)}}
	return l
}

func (s *Session) isCondition(v Value) bool {
	for _, c := range classAttr(v) {
		if c == "condition" {
			return true
		}
	}
	return false
}

func condField(v Value, name string) Value {
	if l, ok := v.(*List); ok {
		if nm := namesOf(l); nm != nil {
			for i, e := range nm.V {
				if e.S == name {
					return l.V[i]
				}
			}
		}
	}
	return Null
}

func condMessage(v Value) string {
	if st, ok := condField(v, "message").(*String); ok && len(st.V) > 0 {
		return st.V[0].S
	}
	return ""
}

func condCall(v Value) *Lang {
	l, _ := condField(v, "call").(*Lang)
	return l
}

// callFunction applies a function value to evaluated arguments.
func (s *Session) callFunction(fn Value, args []Value) Value {
	call := &Lang{Fun: Intern("h")}
	for _, a := range args {
		call.Args = append(call.Args, Arg{Value: a})
	}
	env := s.Global
	switch f := fn.(type) {
	case *Closure:
		sup := make([]Arg, len(args))
		for i, a := range args {
			sup[i] = Arg{Value: a}
		}
		return s.applyClosure(f, call, sup, env, nil)
	case *Builtin:
		ci := &callInfo{call: call, env: env, args: args, names: make([]string, len(args))}
		return f.fn(s, ci)
	}
	s.errorf(call, "attempt to apply non-function")
	return nil
}

// signal runs the handlers of the condition, from the innermost. A calling
// handler runs with the handlers below it; an exiting one unwinds to its
// tryCatch.
func (s *Session) signal(cond Value) {
	cls := classAttr(cond)
	for i := len(s.handlers) - 1; i >= 0; i-- {
		h := s.handlers[i]
		match := false
		for _, c := range cls {
			if c == h.class {
				match = true
				break
			}
		}
		if !match {
			continue
		}
		if !h.calling {
			panic(&unwindSignal{id: h.id, cond: cond, handler: h.fn})
		}
		saved := s.handlers
		s.handlers = append([]handler(nil), saved[:i]...)
		s.callFunction(h.fn, []Value{cond})
		s.handlers = saved
	}
}

// raiseError signals a simpleError and, when nothing handles it, ends the
// evaluation with it.
func (s *Session) raiseError(msg string, call *Lang) {
	cond := makeCondition(msg, callValue(call), "simpleError", "error", "condition")
	s.signal(cond)
	panic(&rError{msg: msg, call: call})
}

// raiseCondition is stop(cond).
func (s *Session) raiseCondition(cond Value) {
	s.signal(cond)
	panic(&rError{msg: condMessage(cond), call: condCall(cond)})
}

func (s *Session) withRestart(name string, f func()) (invoked bool) {
	base := s.handlers
	s.restarts = append(s.restarts, name)
	n := len(s.restarts)
	defer func() {
		s.restarts = s.restarts[:n-1]
		if r := recover(); r != nil {
			if rs, ok := r.(*restartSignal); ok && rs.name == name {
				s.handlers = base
				invoked = true
				return
			}
			panic(r)
		}
	}()
	f()
	return false
}

// doWarning signals a warning condition; when it is not muffled it is kept.
func (s *Session) doWarning(cond Value, msg string, call *Lang) {
	if !s.withRestart("muffleWarning", func() { s.signal(cond) }) {
		s.warnings = append(s.warnings, Warning{Message: msg, Call: callText(call)})
	}
}

func (s *Session) pasteArgs(c *callInfo, skip map[string]bool) string {
	var b strings.Builder
	for i, a := range c.args {
		if skip[c.names[i]] {
			continue
		}
		for _, e := range asStrings(a) {
			if e.NA {
				b.WriteString("NA")
			} else {
				b.WriteString(e.S)
			}
		}
		if sym, ok := a.(*Symbol); ok {
			b.WriteString(sym.Name)
		}
	}
	return b.String()
}

func (s *Session) flagArg(c *callInfo, name string, def bool) bool {
	for i, n := range c.names {
		if n == name {
			if l := logicalOf(s, c, c.args[i], "invalid argument"); len(l) > 0 && l[0] != naInt {
				return l[0] != 0
			}
		}
	}
	return def
}

func (s *Session) runOnExit(fr *frame) {
	if len(fr.onexit) == 0 {
		return
	}
	f, fv, fe := s.flow, s.flowVal, s.flowEnv
	vis := s.visible
	s.flow, s.flowVal, s.flowEnv = flowNone, nil, nil
	exprs := fr.onexit
	fr.onexit = nil
	for _, e := range exprs {
		s.eval(e, fr.env)
	}
	s.flow, s.flowVal, s.flowEnv = f, fv, fe
	s.visible = vis
}

func (s *Session) abort(msg string, err error) {
	panic(&rError{msg: msg, abort: true, err: err})
}

func init() {
	regBuiltin("stop", func(s *Session, c *callInfo) Value {
		callFlag := s.flagArg(c, "call.", true)
		if len(c.args) > 0 && c.names[0] == "" && s.isCondition(c.args[0]) {
			s.raiseCondition(c.args[0])
		}
		var call *Lang
		if callFlag {
			call = s.currentCall()
		}
		s.raiseError(s.pasteArgs(c, map[string]bool{"call.": true, "domain": true}), call)
		return nil
	})
	regBuiltin("warning", func(s *Session, c *callInfo) Value {
		callFlag := s.flagArg(c, "call.", true)
		var cond Value
		var msg string
		if len(c.args) > 0 && c.names[0] == "" && s.isCondition(c.args[0]) {
			cond, msg = c.args[0], condMessage(c.args[0])
		} else {
			msg = s.pasteArgs(c, map[string]bool{"call.": true, "domain": true, "immediate.": true, "noBreaks.": true})
			var call *Lang
			if callFlag {
				call = s.currentCall()
			}
			cond = makeCondition(msg, callValue(call), "simpleWarning", "warning", "condition")
		}
		s.doWarning(cond, msg, condCall(cond))
		s.visible = false
		return mkString(msg)
	})
	regBuiltin("message", func(s *Session, c *callInfo) Value {
		var cond Value
		if len(c.args) > 0 && c.names[0] == "" && s.isCondition(c.args[0]) {
			cond = c.args[0]
		} else {
			msg := s.pasteArgs(c, map[string]bool{"domain": true, "appendLF": true})
			if s.flagArg(c, "appendLF", true) {
				msg += "\n"
			}
			cond = makeCondition(msg, callValue(s.currentCall()), "simpleMessage", "message", "condition")
		}
		if !s.withRestart("muffleMessage", func() { s.signal(cond) }) {
			s.stderr.WriteString(condMessage(cond))
		}
		s.visible = false
		return Null
	})
	regBuiltin("signalCondition", func(s *Session, c *callInfo) Value {
		if len(c.args) == 0 {
			s.errorf(c.call, "argument \"cond\" is missing, with no default")
		}
		s.signal(c.args[0])
		return Null
	})
	mk := func(name string, classes ...string) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			msg := ""
			var call Value = Null
			for i, a := range c.args {
				switch {
				case c.names[i] == "call" || (c.names[i] == "" && i == 1):
					call = a
				case c.names[i] == "message" || (c.names[i] == "" && i == 0):
					if st, ok := a.(*String); ok && len(st.V) > 0 {
						msg = st.V[0].S
					}
				}
			}
			return makeCondition(msg, call, classes...)
		})
	}
	mk("simpleError", "simpleError", "error", "condition")
	mk("simpleWarning", "simpleWarning", "warning", "condition")
	mk("simpleCondition", "simpleCondition", "condition")
	mk("simpleMessage", "simpleMessage", "message", "condition")
	regBuiltin("conditionMessage", func(s *Session, c *callInfo) Value {
		return condField(s.nargs1(c, "conditionMessage"), "message")
	})
	regBuiltin("conditionCall", func(s *Session, c *callInfo) Value {
		return condField(s.nargs1(c, "conditionCall"), "call")
	})
	regBuiltin("invokeRestart", func(s *Session, c *callInfo) Value {
		if len(c.args) == 0 {
			s.errorf(c.call, "argument \"r\" is missing, with no default")
		}
		st, ok := c.args[0].(*String)
		if !ok || len(st.V) != 1 {
			s.errorf(c.call, "bad restart specification")
		}
		for i := len(s.restarts) - 1; i >= 0; i-- {
			if s.restarts[i] == st.V[0].S {
				panic(&restartSignal{name: st.V[0].S})
			}
		}
		s.errorf(c.call, "no 'restart' '%s' found", st.V[0].S)
		return nil
	})
	regSpecial("tryCatch", func(s *Session, c *callInfo) Value {
		var expr, finally Value
		haveExpr := false
		type spec struct {
			class string
			fn    Value
		}
		var specs []spec
		for _, a := range c.raw {
			switch {
			case a.HasName && a.Name == "finally":
				finally = a.Value
			case a.HasName && a.Name == "expr":
				expr, haveExpr = a.Value, true
			case !a.HasName && !haveExpr:
				expr, haveExpr = a.Value, true
			case a.HasName:
				specs = append(specs, spec{a.Name, s.eval(a.Value, c.env)})
			}
		}
		id := s.nextID
		s.nextID++
		base := s.handlers
		if finally != nil {
			defer func() {
				f, fv, fe := s.flow, s.flowVal, s.flowEnv
				s.flow, s.flowVal, s.flowEnv = flowNone, nil, nil
				saved := s.handlers
				s.handlers = base
				s.eval(finally, c.env)
				s.handlers = saved
				s.flow, s.flowVal, s.flowEnv = f, fv, fe
			}()
		}
		h := append([]handler(nil), base...)
		for i := len(specs) - 1; i >= 0; i-- {
			h = append(h, handler{class: specs[i].class, fn: specs[i].fn, id: id})
		}
		s.handlers = h
		return func() (res Value) {
			depth, nframes := s.depth, len(s.frames)
			defer func() {
				s.handlers = base
				if r := recover(); r != nil {
					if us, ok := r.(*unwindSignal); ok && us.id == id {
						s.depth, s.frames = depth, s.frames[:nframes]
						res = s.callFunction(us.handler, []Value{us.cond})
						return
					}
					panic(r)
				}
			}()
			if !haveExpr {
				s.errorf(c.call, "argument \"expr\" is missing, with no default")
			}
			return s.eval(expr, c.env)
		}()
	})
	regSpecial("withCallingHandlers", func(s *Session, c *callInfo) Value {
		var expr Value
		haveExpr := false
		var hs []handler
		for _, a := range c.raw {
			switch {
			case a.HasName && a.Name == "expr":
				expr, haveExpr = a.Value, true
			case !a.HasName && !haveExpr:
				expr, haveExpr = a.Value, true
			case a.HasName:
				hs = append(hs, handler{class: a.Name, fn: s.eval(a.Value, c.env), calling: true})
			}
		}
		base := s.handlers
		h := append([]handler(nil), base...)
		for i := len(hs) - 1; i >= 0; i-- {
			h = append(h, hs[i])
		}
		s.handlers = h
		defer func() { s.handlers = base }()
		if !haveExpr {
			s.errorf(c.call, "argument \"expr\" is missing, with no default")
		}
		return s.eval(expr, c.env)
	})
	muffler := func(restart string) *Builtin {
		return &Builtin{Name: "muffle", fn: func(s *Session, c *callInfo) Value { panic(&restartSignal{name: restart}) }}
	}
	suppress := func(name, class, restart string) {
		regSpecial(name, func(s *Session, c *callInfo) Value {
			if len(c.raw) == 0 {
				return Null
			}
			base := s.handlers
			s.handlers = append(append([]handler(nil), base...), handler{class: class, fn: muffler(restart), calling: true})
			defer func() { s.handlers = base }()
			return s.eval(c.raw[0].Value, c.env)
		})
	}
	suppress("suppressWarnings", "warning", "muffleWarning")
	suppress("suppressMessages", "message", "muffleMessage")
	regSpecial("on.exit", func(s *Session, c *callInfo) Value {
		fr := s.frameOf(c.env)
		var expr Value
		add, after := false, true
		pos := 0
		for _, a := range c.raw {
			switch {
			case a.HasName && a.Name == "expr":
				expr = a.Value
			case a.HasName && a.Name == "add":
				add = len(logicalOf(s, nil2(c), s.eval(a.Value, c.env), "invalid 'add' argument")) > 0 && logicalOf(s, nil2(c), s.eval(a.Value, c.env), "")[0] == 1
			case a.HasName && a.Name == "after":
				after = logicalOf(s, nil2(c), s.eval(a.Value, c.env), "invalid 'after' argument")[0] == 1
			default:
				switch pos {
				case 0:
					expr = a.Value
				case 1:
					add = logicalOf(s, nil2(c), s.eval(a.Value, c.env), "invalid 'add' argument")[0] == 1
				case 2:
					after = logicalOf(s, nil2(c), s.eval(a.Value, c.env), "invalid 'after' argument")[0] == 1
				}
				pos++
			}
		}
		s.visible = false
		if fr == nil {
			return Null
		}
		switch {
		case expr == nil:
			fr.onexit = nil
		case !add:
			fr.onexit = []Value{expr}
		case after:
			fr.onexit = append(fr.onexit, expr)
		default:
			fr.onexit = append([]Value{expr}, fr.onexit...)
		}
		return Null
	})
	regSpecial("try", func(s *Session, c *callInfo) Value {
		var expr Value
		silent := false
		for i, a := range c.raw {
			switch {
			case a.HasName && a.Name == "silent":
				silent = silentArg(s, c, a.Value)
			case a.HasName && a.Name == "expr", !a.HasName && i == 0:
				expr = a.Value
			case !a.HasName && i == 1:
				silent = silentArg(s, c, a.Value)
			}
		}
		if expr == nil {
			return Null
		}
		id := s.nextID
		s.nextID++
		base := s.handlers
		catcher := &Builtin{Name: "try", fn: func(s *Session, ci *callInfo) Value { return Null }}
		s.handlers = append(append([]handler(nil), base...), handler{class: "error", fn: catcher, id: id})
		depth, nframes := s.depth, len(s.frames)
		return func() (res Value) {
			defer func() {
				s.handlers = base
				if r := recover(); r != nil {
					if us, ok := r.(*unwindSignal); ok && us.id == id {
						s.depth, s.frames = depth, s.frames[:nframes]
						msg := condMessage(us.cond)
						text := "Error : " + msg + "\n"
						if cl := condCall(us.cond); cl != nil {
							text = "Error in " + callText(cl) + " : " + msg + "\n"
						}
						if !silent {
							s.stderr.WriteString(text)
						}
						out := &String{V: []Str{{S: text}}}
						out.Attr = []Attr{{Name: "class", Value: mkString("try-error")}, {Name: "condition", Value: us.cond}}
						s.visible = false
						res = out
						return
					}
					panic(r)
				}
			}()
			return s.eval(expr, c.env)
		}()
	})
}

func nil2(c *callInfo) *callInfo { return c }

var _ = syntax.Deparse

func silentArg(s *Session, c *callInfo, e Value) bool {
	l := logicalOf(s, nil2(c), s.eval(e, c.env), "invalid 'silent' argument")
	if len(l) == 0 {
		s.errorf(c.call, "invalid 'silent' argument")
	}
	return l[0] == 1
}
