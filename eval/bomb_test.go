package eval_test

import (
	"context"
	"errors"
	"runtime"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	"github.com/go-embedded-r/r/eval"
)

// The bomb tests assert a REFUSAL: an error that names the limit, and a
// session that is still usable afterwards. They run under a Go memory limit
// (and, in CI, under ulimit -v), never against the whole memory of the host.

func underMemoryLimit(t *testing.T) {
	t.Helper()
	skipOnWasm(t)
	old := debug.SetMemoryLimit(768 << 20)
	t.Cleanup(func() { debug.SetMemoryLimit(old) })
}

func refused(t *testing.T, s *eval.Session, src string, want ...string) {
	t.Helper()
	_, err := s.Eval(src)
	if err == nil {
		t.Fatalf("%q was not refused", src)
	}
	msg := err.Error()
	ok := false
	for _, w := range want {
		if strings.Contains(msg, w) {
			ok = true
		}
	}
	if !ok {
		t.Fatalf("%q: refused with %q, want one of %q", src, msg, want)
	}
	if res, err := s.Eval("1 + 1"); err != nil || res.Value == nil {
		t.Fatalf("session unusable after %q: %v", src, err)
	}
}

func skipOnWasm(t *testing.T) {
	t.Helper()
	if runtime.GOARCH == "wasm" {
		t.Skip("the bombs run natively; wasm hosts have their own memory limits")
	}
}

func TestDefaultsAreNonZero(t *testing.T) {
	skipOnWasm(t)
	// the zero Options is bounded
	s := eval.NewSession(eval.Options{})
	refused(t, s, "while (TRUE) 0", "step limit", "time limit")
}

func TestBombsAreRefusedByDefault(t *testing.T) {
	underMemoryLimit(t)
	s := eval.NewSession(eval.Options{})
	refused(t, s, "numeric(1e12)", "vector", "too large", "limit")
	refused(t, s, "1:1e12", "vector", "too large", "limit")
	refused(t, s, "f <- function(n) f(n + 1); f(0)", "nested too deeply")
	refused(t, s, "f <- function(n) 1 + f(n + 1); f(0)", "nested too deeply")
	refused(t, s, "x <- 0; repeat x <- x + 1", "step limit", "time limit")
	refused(t, s, "rep(1, 1e12)", "vector size", "too large", "limit")
	refused(t, s, "rep(\"a\", times = 1e10)", "vector size", "too large", "limit")
	refused(t, s, "sprintf(\"%1000000000d\", 1L)", "maximal")
	refused(t, s, "format(1, width = 1e9)", "width")
	refused(t, s, "seq(1, 2, by = 1e-300)", "vector size", "too large", "limit")
	refused(t, s, "seq(1, 2, by = 1e-300)", "vector size", "too large", "limit")
	refused(t, s, "matrix(0, 1e6, 1e6)", "vector size", "too large", "limit")
	refused(t, s, "array(0, c(1e5, 1e5, 10))", "too large", "vector size", "limit")
}

// RE2 has no catastrophic backtracking: a classic bomb pattern is answered
// at once, and a huge repetition is refused at compile time.
func TestRegexBombs(t *testing.T) {
	skipOnWasm(t)
	s := eval.NewSession(eval.Options{})
	res, err := s.Eval("x <- paste(rep(\"a\", 20000), collapse = \"\"); grepl(\"(a*)*b\", x)")
	if err != nil || res.Value == nil {
		t.Fatalf("the evil pattern should be answered: %v", err)
	}
	refused(t, s, "grepl(\"(((a{1000}){1000}){1000})\", \"a\")", "invalid regular expression")
	refused(t, s, "grepl(\"a\", \"a\", perl = TRUE)", "perl = TRUE is not supported")
	refused(t, s, "grepl(\"(a)\\\\1\", \"aa\")", "back-references")
}

func TestAllocLimit(t *testing.T) {
	underMemoryLimit(t)
	s := eval.NewSession(eval.Options{MaxAlloc: 64 << 20})
	refused(t, s, "x <- numeric(1e6); repeat x <- x + 1", "memory limit")
}

func TestTimeoutAndCancel(t *testing.T) {
	skipOnWasm(t)
	s := eval.NewSession(eval.Options{MaxSteps: eval.Unbounded, Timeout: 150 * time.Millisecond})
	_, err := s.Eval("while (TRUE) 0")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("want a deadline error, got %v", err)
	}
	s = eval.NewSession(eval.Options{MaxSteps: eval.Unbounded, Timeout: eval.Unbounded})
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	_, err = s.EvalContext(ctx, "while (TRUE) 0")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("a context deadline must stop an unbounded session, got %v", err)
	}
}

func TestUnboundedIsStamped(t *testing.T) {
	skipOnWasm(t)
	// the explicit opt-out works...
	s := eval.NewSession(eval.Options{MaxSteps: eval.Unbounded, MaxAlloc: eval.Unbounded, Timeout: eval.Unbounded})
	res, err := s.Eval("s <- 0; for (i in 1:100000) s <- s + i; s")
	if err != nil || res.Value == nil {
		t.Fatalf("unbounded session failed: %v", err)
	}
	// ...and no other negative value is taken for it
	defer func() {
		if recover() == nil {
			t.Fatal("a negative limit other than Unbounded must panic")
		}
	}()
	eval.NewSession(eval.Options{MaxSteps: -5})
}
