package eval

import (
	"fmt"
	"strings"

	"github.com/go-embedded-r/r/syntax"
)

// Options bound a session. A zero field selects its default.
type Options struct {
	// MaxDepth bounds the nesting of evaluations: R's option "expressions"
	// (default 5000). Exceeding it is R's error "evaluation nested too deeply".
	MaxDepth int
	// MaxVector bounds the length of any vector the runtime allocates.
	MaxVector int
	// MaxSteps bounds the number of calls evaluated by one Eval (0: no bound).
	MaxSteps int
	// Parse bounds the source text of an Eval.
	Parse syntax.Limits
}

func (o Options) withDefaults() Options {
	if o.MaxDepth <= 0 {
		o.MaxDepth = 5000
	}
	if o.MaxVector <= 0 {
		o.MaxVector = 1 << 28
	}
	return o
}

// Warning is a warning raised while evaluating.
type Warning struct {
	Message string
	Call    string
}

// Result is what one Eval returns.
type Result struct {
	Value    Value
	Visible  bool
	Warnings []Warning
}

// Error is an R error that reached the top level. Message is R's
// conditionMessage and Call the deparsed conditionCall ("" for none).
type Error struct {
	Message string
	Call    string
}

func (e *Error) Error() string {
	if e.Call == "" {
		return "Error: " + e.Message
	}
	return "Error in " + e.Call + " : " + e.Message
}

// Session is one interpreter: a global environment that persists from one
// Eval to the next. It is not safe for concurrent use.
type Session struct {
	Base, Global *Env
	opts         Options

	frames   []*frame
	warnings []Warning
	visible  bool
	steps    int
	depth    int

	// flow is set by break, next and return in statement position, and read
	// by the construct that owns them; see eval.go.
	flow    flowCode
	flowVal Value
	flowEnv *Env
}

type frame struct {
	call *Lang
	fn   *Closure
	env  *Env
}

// rError is the panic value of an R error.
type rError struct {
	msg  string
	call *Lang
}

// NewSession makes a session whose global environment is empty, enclosed by
// the base environment holding the primitives.
func NewSession(opts Options) *Session {
	s := &Session{opts: opts.withDefaults()}
	s.Base = &Env{name: "base"}
	for name, b := range builtins {
		s.Base.Set(name, b)
	}
	s.Global = &Env{parent: s.Base, name: "R_GlobalEnv"}
	return s
}

// Eval parses src and evaluates its expressions in order in the global
// environment. The result is that of the last one. Errors are *Error (a
// runtime error), *syntax.Error or *syntax.LimitError.
func (s *Session) Eval(src string) (*Result, error) { return s.EvalIn(src, s.Global) }

// EvalIn is Eval in the given environment (a child of the global one, for
// a cell that must not leave names behind).
func (s *Session) EvalIn(src string, env *Env) (res *Result, err error) {
	f, perr := syntax.Parse(src, syntax.Options{Limits: s.opts.Parse})
	if perr != nil {
		return nil, perr
	}
	s.warnings = nil
	s.steps = 0
	s.frames = s.frames[:0]
	s.flow = flowNone
	res = &Result{Value: Null}
	defer func() {
		if r := recover(); r != nil {
			switch e := r.(type) {
			case *rError:
				err = &Error{Message: e.msg, Call: callText(e.call)}
			case *flowSignal:
				err = &Error{Message: noLoopMessage(e.f)}
			default:
				panic(r)
			}
			res = nil
		}
	}()
	for _, e := range f.Exprs {
		s.visible = true
		v := s.eval(FromSyntax(e), env)
		res.Value, res.Visible = v, s.visible
	}
	res.Warnings = s.warnings
	return res, nil
}

func noLoopMessage(f flowCode) string {
	if f == flowReturn {
		return "no function to return from, jumping to top level"
	}
	return "no loop for break/next, jumping to top level"
}

// Get and Set read and write a variable of the global environment.
func (s *Session) Get(name string) (Value, bool) { return s.Global.Get(name) }
func (s *Session) Set(name string, v Value)      { s.Global.Set(name, v) }

func callText(c *Lang) string {
	if c == nil {
		return ""
	}
	return strings.Join(syntax.Deparse(toSyntax(c), nil), "\n")
}

// errorf raises an R error with the given call.
func (s *Session) errorf(call *Lang, format string, a ...any) {
	panic(&rError{msg: fmt.Sprintf(format, a...), call: call})
}

// warningf records a warning.
func (s *Session) warningf(call *Lang, format string, a ...any) {
	s.warnings = append(s.warnings, Warning{Message: fmt.Sprintf(format, a...), Call: callText(call)})
}

// currentCall is the call of the innermost function being evaluated, the
// call R attaches to errors raised by the evaluator itself.
func (s *Session) currentCall() *Lang {
	if n := len(s.frames); n > 0 {
		return s.frames[n-1].call
	}
	return nil
}
