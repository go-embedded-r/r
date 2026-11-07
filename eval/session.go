package eval

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/go-embedded-r/r/syntax"
)

// Unbounded, given as MaxSteps, MaxAlloc or Timeout, removes that limit. It
// is the only way to do so: a zero selects the default, and any other
// negative value panics in NewSession. Use it only for trusted code.
const Unbounded = -1

// The defaults of Options, chosen by measurement (TestMeasureLimits): the
// golden snippets and the heaviest cells a notebook runs on purpose (a loop
// of 10^6 iterations: 3 million steps; a list grown by 20 000 assignments:
// 1.6 GB of copies, 0.63 s) each get about ten times the largest legitimate
// use.
const (
	DefaultMaxDepth = 5000
	DefaultMaxSteps = 30_000_000
	DefaultMaxAlloc = 16 << 30
	DefaultTimeout  = 10 * time.Second
)

// Options bound a session. A zero field selects its default; Unbounded
// (MaxSteps, MaxAlloc, Timeout only) removes the limit explicitly.
type Options struct {
	// MaxDepth bounds the nesting of evaluations: R's option "expressions"
	// (default 5000). Exceeding it is R's error "evaluation nested too deeply".
	MaxDepth int
	// MaxVector bounds the length of any vector the runtime allocates.
	MaxVector int
	// MaxAlloc bounds the bytes of vector storage one Eval allocates, in
	// total over the evaluation, copies included (default DefaultMaxAlloc).
	// Exceeding it ends the evaluation with an error.
	MaxAlloc int
	// MaxSteps bounds the number of calls and loop iterations evaluated by
	// one Eval (default DefaultMaxSteps).
	MaxSteps int
	// Timeout bounds the time one Eval may take (default DefaultTimeout); a
	// context passed to EvalContext bounds it too.
	Timeout time.Duration
	// Parse bounds the source text of an Eval.
	Parse syntax.Limits
}

func (o Options) withDefaults() Options {
	if o.MaxDepth <= 0 {
		o.MaxDepth = DefaultMaxDepth
	}
	if o.MaxVector <= 0 {
		o.MaxVector = 1 << 28
	}
	// inside the session, 0 means no limit
	norm := func(v, def int64) int64 {
		switch {
		case v == 0:
			return def
		case v == Unbounded:
			return 0
		case v < 0:
			panic("eval: a negative limit other than eval.Unbounded")
		}
		return v
	}
	o.MaxSteps = int(min(norm(int64(o.MaxSteps), DefaultMaxSteps), math.MaxInt))
	o.MaxAlloc = int(min(norm(int64(o.MaxAlloc), DefaultMaxAlloc), math.MaxInt))
	o.Timeout = time.Duration(norm(int64(o.Timeout), int64(DefaultTimeout)))
	return o
}

// Warning is a warning raised while evaluating.
type Warning struct {
	Message string
	Call    string
}

// Result is what one Eval returns.
type Result struct {
	// Value is the value of the last expression.
	Value Value
	// Visible says whether R would print it.
	Visible bool
	// Output is what the cell printed on standard output: what print() and
	// cat() wrote, and the auto-printed value of every visible top-level
	// expression, as an R console shows them.
	Output string
	// Stderr is what the cell wrote on standard error: message() output.
	Stderr string
	// Warnings are the warnings raised and not muffled.
	Warnings []Warning
}

// Error is an R error that reached the top level. Message is R's
// conditionMessage and Call the deparsed conditionCall ("" for none).
type Error struct {
	Message string
	Call    string
	// Condition is the condition object that was signalled (nil for errors
	// that no handler could see: limits and interrupts).
	Condition Value
	// Output and Stderr are what the cell had written before it failed.
	Output, Stderr string
	// Warnings are the warnings raised before the error.
	Warnings []Warning
	err      error
}

func (e *Error) Error() string {
	if e.Call == "" {
		return "Error: " + e.Message
	}
	return "Error in " + e.Call + " : " + e.Message
}

// Unwrap is context.DeadlineExceeded or context.Canceled for an evaluation
// that was interrupted, and nil otherwise.
func (e *Error) Unwrap() error { return e.err }

// Session is one interpreter: a global environment that persists from one
// Eval to the next. It is not safe for concurrent use.
type Session struct {
	Base, Global, Empty *Env
	opts                Options

	frames   []*frame
	warnings []Warning
	visible  bool
	steps    int
	depth    int
	alloc    int64
	tick     int
	ctx      context.Context
	deadline time.Time

	handlers []handler
	restarts []restartEntry
	restartID int
	lastErr   string // geterrmessage()
	nextID   int
	out      strings.Builder
	stderr   strings.Builder

	noQuote    bool // print(quote = FALSE) is running
	noDispatch bool // the next primitive runs its default, not a method

	// flow is set by break, next and return in statement position, and read
	// by the construct that owns them; see eval.go.
	flow    flowCode
	flowVal Value
	flowEnv *Env
}

type frame struct {
	call     *Lang
	fn       *Closure
	env      *Env
	supplied []Arg         // the arguments as supplied (promises), for UseMethod and NextMethod
	di       *dispatchInfo // set when the function runs as an S3 method
	caller   *Env
	onexit   []Value
}

// rError is the panic value of an R error that no handler took.
type rError struct {
	msg   string
	call  *Lang
	cond  Value
	abort bool // a limit or an interrupt: not an R condition
	err   error
}

// NewSession makes a session whose global environment is empty, enclosed by
// the base environment holding the primitives.
func NewSession(opts Options) *Session {
	installS3()
	s := &Session{opts: opts.withDefaults()}
	s.Empty = &Env{name: "R_EmptyEnv"}
	s.Base = &Env{name: "base", parent: s.Empty}
	for name, b := range builtins {
		s.Base.Set(name, b)
	}
	lower, upper := make([]string, 26), make([]string, 26)
	for i := range lower {
		lower[i], upper[i] = string(rune('a'+i)), string(rune('A'+i))
	}
	s.Base.Set("pi", mkDouble(math.Pi))
	s.Base.Set("T", mkLogical(true))
	s.Base.Set("F", mkLogical(false))
	s.Base.Set("letters", strVec(lower))
	s.Base.Set("LETTERS", strVec(upper))
	s.Global = &Env{parent: s.Base, name: "R_GlobalEnv"}
	s.installBaselib()
	s.resetCounters()
	return s
}

// Eval parses src and evaluates its expressions in order in the global
// environment, as one notebook cell. The result is that of the last one.
// Errors are *Error (a runtime error: R's message and call, with the output
// and warnings so far), *syntax.Error or *syntax.LimitError.
func (s *Session) Eval(src string) (*Result, error) {
	return s.EvalContext(context.Background(), src)
}

// EvalContext is Eval that ends, with an *Error wrapping the context's
// error, when ctx is done; Options.Timeout bounds it as well.
func (s *Session) EvalContext(ctx context.Context, src string) (*Result, error) {
	return s.EvalIn(ctx, src, s.Global)
}

// EvalIn is EvalContext in the given environment (a child of the global one,
// for a cell that must not leave names behind).
func (s *Session) EvalIn(ctx context.Context, src string, env *Env) (res *Result, err error) {
	f, perr := syntax.Parse(src, syntax.Options{Limits: s.opts.Parse})
	if perr != nil {
		return nil, perr
	}
	s.warnings = nil
	s.steps, s.tick, s.alloc, s.depth = 0, 0, 0, 0
	s.frames = s.frames[:0]
	s.handlers, s.restarts = nil, nil
	s.flow = flowNone
	s.noDispatch = false
	s.out.Reset()
	s.stderr.Reset()
	if ctx == nil {
		ctx = context.Background()
	}
	s.ctx, s.deadline = ctx, time.Time{}
	if s.opts.Timeout > 0 {
		s.deadline = time.Now().Add(s.opts.Timeout)
	}
	res = &Result{Value: Null}
	defer func() {
		s.ctx = nil
		if r := recover(); r != nil {
			e := &Error{Output: s.out.String(), Stderr: s.stderr.String(), Warnings: s.warnings}
			switch x := r.(type) {
			case *rError:
				e.Message, e.Call, e.Condition, e.err = x.msg, callText(x.call), x.cond, x.err
			case *flowSignal:
				e.Message = noLoopMessage(x.f)
			case *unwindSignal:
				e.Message = "unhandled condition"
			case *restartSignal:
				e.Message = "no 'restart' '" + x.name + "' found"
			default:
				panic(r)
			}
			res, err = nil, e
		}
	}()
	for _, e := range f.Exprs {
		s.visible = true
		v := s.eval(FromSyntax(e), env)
		res.Value, res.Visible = v, s.visible
		if s.visible {
			s.printValue(v, env)
		}
	}
	res.Output, res.Stderr, res.Warnings = s.out.String(), s.stderr.String(), s.warnings
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

// errorf raises an R error with the given call: a condition that the
// handlers see before the evaluation ends.
func (s *Session) errorf(call *Lang, format string, a ...any) {
	s.raiseError(fmt.Sprintf(format, a...), call)
}

// warningf raises a warning condition.
func (s *Session) warningf(call *Lang, format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	s.doWarning(makeCondition(msg, callValue(call), "simpleWarning", "warning", "condition"), msg, call)
}

// currentCall is the call of the innermost function being evaluated, the
// call R attaches to errors raised by the evaluator itself.
func (s *Session) currentCall() *Lang {
	if n := len(s.frames); n > 0 {
		return s.frames[n-1].call
	}
	return nil
}

// Builtins returns the names of the primitives and native functions this
// evaluator provides.
func Builtins() []string {
	loadBaselib()
	out := make([]string, 0, len(builtins)+len(baselibNames))
	for n := range builtins {
		out = append(out, n)
	}
	out = append(out, baselibNames...)
	sort.Strings(out)
	return out
}

// resetCounters clears the budgets that installing the library spent.
func (s *Session) resetCounters() {
	s.steps, s.tick, s.alloc, s.depth = 0, 0, 0, 0
	s.visible = false
}
