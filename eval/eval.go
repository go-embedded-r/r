package eval

import (
	"strconv"
	"strings"

	"github.com/go-embedded-r/r/syntax"
)

type flowCode int

const (
	flowNone flowCode = iota
	flowBreak
	flowNext
	flowReturn
)

// flowSignal is the panic value of a break, next or return reached outside
// statement position (through a promise, or an argument of a builtin).
type flowSignal struct {
	f   flowCode
	val Value
	env *Env
}

// callInfo is what a primitive receives.
type callInfo struct {
	call  *Lang
	env   *Env
	raw   []Arg   // the arguments as written (specials)
	args  []Value // the evaluated arguments (builtins)
	names []string
}

var builtins = map[string]*Builtin{}

func regBuiltin(name string, fn func(s *Session, c *callInfo) Value) {
	builtins[name] = &Builtin{Name: name, fn: fn}
}

func regSpecial(name string, fn func(s *Session, c *callInfo) Value) {
	builtins[name] = &Builtin{Name: name, Special: true, fn: fn}
}

// eval evaluates v in env. A break, next or return signalled by a construct
// in statement position is turned into a panic here: only evalStmt lets it
// travel as a flag.
func (s *Session) eval(v Value, env *Env) Value {
	switch x := v.(type) {
	case *Symbol:
		s.visible = true
		return s.evalSym(x, env)
	case *Lang:
		r := s.evalLang(x, env)
		if s.flow != flowNone {
			sig := &flowSignal{f: s.flow, val: s.flowVal, env: s.flowEnv}
			s.flow, s.flowVal, s.flowEnv = flowNone, nil, nil
			panic(sig)
		}
		return r
	case *Promise:
		s.visible = true
		return s.force(x)
	}
	s.visible = true
	return v
}

// evalStmt is eval for a statement: control flow stays in s.flow.
func (s *Session) evalStmt(v Value, env *Env) Value {
	if l, ok := v.(*Lang); ok {
		return s.evalLang(l, env)
	}
	return s.eval(v, env)
}

func (s *Session) evalSym(sym *Symbol, env *Env) Value {
	b, _ := env.find(sym)
	if b == nil {
		s.errorf(s.currentCall(), "object '%s' not found", sym.Name)
	}
	switch x := b.val.(type) {
	case *Promise:
		return s.force(x)
	case missingT:
		s.errorf(s.currentCall(), "argument \"%s\" is missing, with no default", sym.Name)
	}
	return b.val
}

func (s *Session) force(p *Promise) Value {
	if p.forced {
		return p.val
	}
	if p.forcing {
		s.errorf(s.currentCall(), "promise already under evaluation: recursive default argument reference or earlier problems?")
	}
	p.forcing = true
	v := s.eval(p.Expr, p.Env)
	p.val, p.forced, p.forcing, p.Env = v, true, false, nil
	return v
}

func (s *Session) findFun(sym *Symbol, env *Env, call *Lang) Value {
	for e := env; e != nil; e = e.parent {
		b := e.lookup(sym)
		if b == nil {
			continue
		}
		v := b.val
		if p, ok := v.(*Promise); ok {
			v = s.force(p)
		}
		switch v.(type) {
		case *Closure, *Builtin:
			return v
		}
	}
	s.errorf(call, "could not find function \"%s\"", sym.Name)
	return nil
}

func (s *Session) evalLang(c *Lang, env *Env) Value {
	s.depth++
	defer func() { s.depth-- }()
	if s.depth > s.opts.MaxDepth {
		s.errorf(nil, "evaluation nested too deeply: infinite recursion / options(expressions=)?")
	}
	if s.opts.MaxSteps > 0 {
		s.steps++
		if s.steps > s.opts.MaxSteps {
			s.errorf(nil, "evaluation step limit reached")
		}
	}
	var fn Value
	switch f := c.Fun.(type) {
	case *Symbol:
		fn = s.findFun(f, env, c)
	default:
		fn = s.eval(c.Fun, env)
	}
	s.visible = true
	switch f := fn.(type) {
	case *Builtin:
		ci := &callInfo{call: c, env: env, raw: c.Args}
		if !f.Special {
			s.evalArgs(ci)
		}
		return f.fn(s, ci)
	case *Closure:
		return s.applyClosure(f, c, s.promiseArgs(c.Args, env), env)
	}
	s.errorf(s.currentCall(), "attempt to apply non-function")
	return nil
}

// evalArgs evaluates the arguments of a builtin call, expanding `...`.
func (s *Session) evalArgs(ci *callInfo) {
	for _, a := range ci.raw {
		if sym, ok := a.Value.(*Symbol); ok && sym.Name == "..." {
			if d := s.dotsOf(ci.env); d != nil {
				for _, da := range d.Args {
					ci.args = append(ci.args, s.argValue(da.Value))
					ci.names = append(ci.names, da.Name)
				}
			}
			continue
		}
		ci.args = append(ci.args, s.argValue(a.Value, ci.env))
		ci.names = append(ci.names, a.Name)
	}
}

// argValue evaluates one argument: a promise is forced, a constant stands
// for itself, anything else is evaluated in env.
func (s *Session) argValue(v Value, env ...*Env) Value {
	switch x := v.(type) {
	case missingT:
		return v
	case *Promise:
		return s.force(x)
	}
	if len(env) == 0 {
		return v
	}
	return s.eval(v, env[0])
}

func (s *Session) dotsOf(env *Env) *Dots {
	if b, _ := env.find(Intern("...")); b != nil {
		d, _ := b.val.(*Dots)
		return d
	}
	return nil
}

// promiseArgs turns the arguments of a closure call into promises.
func (s *Session) promiseArgs(raw []Arg, env *Env) []Arg {
	out := make([]Arg, 0, len(raw))
	for _, a := range raw {
		switch x := a.Value.(type) {
		case *Symbol:
			if x.Name == "..." {
				if d := s.dotsOf(env); d != nil {
					out = append(out, d.Args...)
				}
				continue
			}
			out = append(out, Arg{Name: a.Name, HasName: a.HasName, Value: &Promise{Expr: x, Env: env}})
		case *Lang:
			out = append(out, Arg{Name: a.Name, HasName: a.HasName, Value: &Promise{Expr: x, Env: env}})
		default:
			out = append(out, a) // constants and the empty argument
		}
	}
	return out
}

// applyClosure calls f with supplied (promises, constants, Missing).
func (s *Session) applyClosure(f *Closure, call *Lang, supplied []Arg, caller *Env) (result Value) {
	fenv := &Env{parent: f.Env, fn: f}
	s.bindArgs(f, call, supplied, fenv)
	s.frames = append(s.frames, &frame{call: call, fn: f, env: fenv})
	defer func() {
		s.frames = s.frames[:len(s.frames)-1]
		if r := recover(); r != nil {
			if sig, ok := r.(*flowSignal); ok && sig.f == flowReturn && sig.env == fenv {
				result = sig.val
				return
			}
			panic(r)
		}
	}()
	v := s.evalStmt(f.Body, fenv)
	if s.flow != flowNone {
		switch {
		case s.flow == flowReturn && s.flowEnv == fenv:
			v = s.flowVal
			s.flow, s.flowVal, s.flowEnv = flowNone, nil, nil
		default:
			sig := &flowSignal{f: s.flow, val: s.flowVal, env: s.flowEnv}
			s.flow, s.flowVal, s.flowEnv = flowNone, nil, nil
			if sig.f != flowReturn {
				// a break or next that no loop of this function owns
				s.errorf(s.currentCall(), "no loop for break/next, jumping to top level")
			}
			panic(sig) // a return for another function's frame
		}
	}
	return v
}

// bindArgs matches the supplied arguments to the formals, in R's three
// passes (exact name, partial name, position), and binds them in fenv.
func (s *Session) bindArgs(f *Closure, call *Lang, supplied []Arg, fenv *Env) {
	formals := f.Formals
	nf := len(formals)
	bound := make([]Value, nf)
	have := make([]bool, nf)
	used := make([]bool, len(supplied))
	exact := make([]bool, len(supplied))
	dotsAt := -1
	for i, fm := range formals {
		if fm.Name.Name == "..." {
			dotsAt = i
		}
	}
	// 1. exact
	for i, fm := range formals {
		if i == dotsAt {
			continue
		}
		for j, a := range supplied {
			if used[j] || !a.HasName || a.Name == "" || a.Name != fm.Name.Name {
				continue
			}
			if have[i] {
				s.errorf(call, "formal argument \"%s\" matched by multiple actual arguments", fm.Name.Name)
			}
			bound[i], have[i], used[j], exact[j] = a.Value, true, true, true
		}
	}
	// 2. partial (formals before ... only)
	partialOf := make([]int, len(supplied))
	for j := range partialOf {
		partialOf[j] = -1
	}
	for i, fm := range formals {
		if have[i] || i == dotsAt || (dotsAt >= 0 && i > dotsAt) {
			continue
		}
		for j, a := range supplied {
			if exact[j] || !a.HasName || a.Name == "" || !strings.HasPrefix(fm.Name.Name, a.Name) {
				continue
			}
			if have[i] {
				s.errorf(call, "formal argument \"%s\" matched by multiple actual arguments", fm.Name.Name)
			}
			if partialOf[j] >= 0 {
				s.errorf(call, "argument %d matches multiple formal arguments", j+1)
			}
			bound[i], have[i], used[j] = a.Value, true, true
			partialOf[j] = i
		}
	}
	// 3. positional
	fi := 0
	for j, a := range supplied {
		if used[j] || (a.HasName && a.Name != "") {
			continue
		}
		for fi < nf && (have[fi] || fi == dotsAt) {
			if fi == dotsAt {
				break
			}
			fi++
		}
		if fi >= nf || fi == dotsAt {
			break
		}
		bound[fi], have[fi], used[j] = a.Value, true, true
		fi++
	}
	// 4. the rest: dots, or an error
	var rest []Arg
	for j, a := range supplied {
		if !used[j] {
			rest = append(rest, a)
		}
	}
	if len(rest) > 0 && dotsAt < 0 {
		parts := make([]string, len(rest))
		for i, a := range rest {
			txt := strings.Join(syntax.Deparse(toSyntax(a.Value), nil), " ")
			if a.HasName && a.Name != "" {
				txt = a.Name + " = " + txt
			}
			parts[i] = txt
		}
		word := "argument"
		if len(rest) > 1 {
			word = "arguments"
		}
		s.errorf(call, "unused %s (%s)", word, strings.Join(parts, ", "))
	}
	for i, fm := range formals {
		switch {
		case i == dotsAt:
			fenv.set(fm.Name, &Dots{Args: rest})
		case have[i] && !isMissingValue(bound[i]):
			fenv.set(fm.Name, bound[i])
		case fm.Default != nil:
			fenv.set(fm.Name, &Promise{Expr: fm.Default, Env: fenv})
			markMissing(fenv, fm.Name)
		default:
			fenv.set(fm.Name, Missing)
			markMissing(fenv, fm.Name)
		}
	}
}

func isMissingValue(v Value) bool { _, ok := v.(missingT); return ok }

func markMissing(e *Env, sym *Symbol) {
	if e.missing == nil {
		e.missing = map[*Symbol]bool{}
	}
	e.missing[sym] = true
}

// isMissingArg is missing(sym) in a function frame: the formal was not
// supplied, or was supplied as a symbol that is itself missing in the caller.
func (s *Session) isMissingArg(sym *Symbol, env *Env) bool {
	if env.missing[sym] {
		return true
	}
	b := env.lookup(sym)
	if b == nil {
		return false
	}
	switch x := b.val.(type) {
	case missingT:
		return true
	case *Promise:
		if inner, ok := x.Expr.(*Symbol); ok && x.Env != nil && x.Env.fn != nil {
			return s.isMissingArg(inner, x.Env)
		}
	}
	return false
}

// ---- conditions of `if`, `while` ------------------------------------------

// asCondition coerces the value of an if/while condition as R does.
func (s *Session) asCondition(v Value, call *Lang) bool {
	if Length(v) > 1 {
		s.errorf(call, "the condition has length > 1")
	}
	if Length(v) == 0 {
		s.errorf(call, "argument is of length zero")
	}
	var r int32
	switch x := v.(type) {
	case *Logical:
		r = x.V[0]
	case *Integer:
		r = x.V[0]
		if r != naInt && r != 0 {
			r = 1
		}
	case *Double:
		switch {
		case x.V[0] != x.V[0]:
			r = naInt
		case x.V[0] != 0:
			r = 1
		}
	case *String:
		if x.V[0].NA {
			r = naInt
		} else {
			switch x.V[0].S {
			case "T", "TRUE", "true", "True":
				r = 1
			case "F", "FALSE", "false", "False":
				r = 0
			default:
				s.errorf(call, "argument is not interpretable as logical")
			}
		}
	default:
		s.errorf(call, "argument is not interpretable as logical")
	}
	if r == naInt {
		if _, isLgl := v.(*Logical); isLgl {
			s.errorf(call, "missing value where TRUE/FALSE needed")
		}
		s.errorf(call, "argument is not interpretable as logical")
	}
	return r != 0
}

func itoa(i int) string { return strconv.Itoa(i) }
