package eval

import "fmt"

// One counter, Session.depth, bounds both the nesting of evaluations and the
// recursion of the natives that walk the structure of a value (identical,
// printing, deparse, unlist, anyNA, serialization, ...). Options.MaxDepth is
// its limit. A walker calls descend before it recurses and defers ascend:
//
//	s.descend()
//	defer s.ascend()
//
// descend checks BEFORE it counts, so a refusal leaves the counter as it was.
// Past the limit the walker raises R's own error, which the session survives:
// the Go stack, and on js/wasm the engine's stack, never gets there.

// stackBytesPerLevel turns a depth into the figure of R's "C stack usage"
// message. It is an estimate of what one level costs R, for the wording only.
const stackBytesPerLevel = 1600

// walkSession is the session of the exported functions that walk a value
// without one (ToRDS, FromSyntax): a bare counter with the default bound. A
// walk too deep there panics with the R error, as the session's own would.
func walkSession() *Session {
	return &Session{opts: Options{MaxDepth: DefaultMaxDepth}}
}

// descend enters one level of a recursive walk over data.
func (s *Session) descend() {
	if s.depth >= s.opts.MaxDepth {
		s.dataTooDeep()
	}
	s.depth++
}

// ascend leaves the level descend entered.
func (s *Session) ascend() { s.depth-- }

// dataTooDeep is R's error for a recursion on data that went too deep.
func (s *Session) dataTooDeep() {
	s.raiseStack(fmt.Sprintf("C stack usage  %d is too close to the limit", (s.depth+1)*stackBytesPerLevel), "CStackOverflowError")
}

// evalTooDeep is R's error for evaluations nested past options(expressions=).
func (s *Session) evalTooDeep() {
	s.raiseStack("evaluation nested too deeply: infinite recursion / options(expressions=)?", "expressionStackOverflowError")
}

// raiseStack signals the condition R signals for a stack overflow: its class
// is the specific one, then stackOverflowError, error, condition.
func (s *Session) raiseStack(msg, class string) {
	cond := makeCondition(msg, Null, class, "stackOverflowError", "error", "condition")
	s.signal(cond)
	panic(&rError{msg: msg})
}
