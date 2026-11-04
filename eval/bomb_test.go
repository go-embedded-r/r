package eval_test

import (
	"context"
	"errors"
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

func TestDefaultsAreNonZero(t *testing.T) {
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
}

func TestAllocLimit(t *testing.T) {
	underMemoryLimit(t)
	s := eval.NewSession(eval.Options{MaxAlloc: 64 << 20})
	refused(t, s, "x <- numeric(1e6); repeat x <- x + 1", "memory limit")
}

func TestTimeoutAndCancel(t *testing.T) {
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
