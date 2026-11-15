package eval_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/go-embedded-r/r/eval"
)

// Issue 42: deep data and deep call nesting ended the instance on js/wasm with
// an uncaught JS RangeError. Every test here asserts an R error, by its class
// and its message, and then that the same session still evaluates 1 + 1. They
// run under the platform's default bound (5000 natively, 1000 on js and
// wasip1), so on WebAssembly they prove that the bound is below the engine's.

const (
	deepData = "x <- list(); for (i in 1:50000) x <- list(x); "
	cstack   = "C stack usage"
	exprDeep = "evaluation nested too deeply: infinite recursion / options(expressions=)?"
)

func evalString(t *testing.T, s *eval.Session, src string) string {
	t.Helper()
	r, err := s.Eval(src)
	if err != nil {
		t.Fatalf("%s: %v", src, err)
	}
	st, ok := r.Value.(*eval.String)
	if !ok || len(st.V) != 1 {
		t.Fatalf("%s: %#v", src, r.Value)
	}
	return st.V[0].S
}

// assertStillAlive is the point of the issue: the session survives.
func assertStillAlive(t *testing.T, s *eval.Session) {
	t.Helper()
	r, err := s.Eval("1 + 1")
	if err != nil {
		t.Fatalf("session after the error: %v", err)
	}
	if d, ok := r.Value.(*eval.Double); !ok || len(d.V) != 1 || d.V[0] != 2 {
		t.Fatalf("1 + 1 after the error: %#v", r.Value)
	}
}

func TestDepthBoundOnData(t *testing.T) {
	ops := []struct{ name, op string }{
		{"identical", "identical(x, x)"},
		{"print", "print(x)"},
		{"format", "format(x)"},
		{"deparse", "deparse(x)"},
		{"unlist", "unlist(x)"},
		{"anyNA", "anyNA(x, recursive = TRUE)"},
		{"saveRDS", "saveRDS(x, '/deep.rds')"},
		{"unique", "unique(list(x, x))"},
	}
	s := eval.NewSession(eval.Options{})
	for _, o := range ops {
		t.Run(o.name, func(t *testing.T) {
			if _, err := s.Eval(deepData); err != nil {
				t.Fatalf("building the data: %v", err)
			}
			// as a top-level error
			_, err := s.Eval(o.op)
			var e *eval.Error
			if !errors.As(err, &e) {
				t.Fatalf("%s: want an R error, got %v", o.op, err)
			}
			if !strings.HasPrefix(e.Message, cstack) || !strings.HasSuffix(e.Message, "is too close to the limit") {
				t.Errorf("%s: message %q", o.op, e.Message)
			}
			// as a condition that tryCatch sees, with R's classes
			cls := evalString(t, s, "paste(class(tryCatch("+o.op+", error = function(e) e)), collapse = ' ')")
			if cls != "CStackOverflowError stackOverflowError error condition" {
				t.Errorf("%s: class %q", o.op, cls)
			}
			assertStillAlive(t, s)
		})
	}
}

func TestDepthBoundOnCalls(t *testing.T) {
	srcs := map[string]string{
		"plain":     "f <- function(n) f(n + 1); f(1)",
		"sum":       "f <- function(n) 1 + f(n + 1); f(0)",
		"sapply":    "f <- function() sapply(1, function(i) f()); f()",
		"tryCatch":  "f <- function() tryCatch(f(), finally = NULL); f()",
		"do.call":   "f <- function() do.call(f, list()); f()",
		"Reduce":    "f <- function() Reduce(function(a, b) f(), 1:2); f()",
		"S3":        "g <- function(x) UseMethod('g'); g.default <- function(x) g(x); g(1)",
		"callingH":  "f <- function() withCallingHandlers(f(), warning = function(w) 0); f()",
		"recursion": "f <- function(n) if (n == 0) 0 else 1 + f(n - 1); f(1e6)",
		// all.equal is R code of ours: it recurses through the evaluator
		"all.equal": deepData + "all.equal(x, x)",
	}
	s := eval.NewSession(eval.Options{})
	for name, src := range srcs {
		t.Run(name, func(t *testing.T) {
			_, err := s.Eval(src)
			var e *eval.Error
			if !errors.As(err, &e) || e.Message != exprDeep {
				t.Fatalf("want %q, got %v", exprDeep, err)
			}
			wrapped := "paste(class(tryCatch({" + src + "}, error = function(e) e)), collapse = ' ')"
			if cls := evalString(t, s, wrapped); cls != "expressionStackOverflowError stackOverflowError error condition" {
				t.Errorf("class %q", cls)
			}
			assertStillAlive(t, s)
		})
	}
}

// The bound is Options.MaxDepth, per session: a larger one lets data through
// that the default refuses, and the default of the next session is unchanged.
func TestDepthBoundIsAnOption(t *testing.T) {
	small := eval.NewSession(eval.Options{MaxDepth: 50})
	if _, err := small.Eval("x <- list(); for (i in 1:60) x <- list(x); identical(x, x)"); err == nil || !strings.Contains(err.Error(), cstack) {
		t.Errorf("60 levels under MaxDepth 50: %v", err)
	}
	assertStillAlive(t, small)
	if r, err := small.Eval("x <- list(); for (i in 1:30) x <- list(x); identical(x, x)"); err != nil || r.Value.(*eval.Logical).V[0] != 1 {
		t.Errorf("30 levels under MaxDepth 50: %v", err)
	}
	// the bound is shared with the evaluator: the frames in use count
	if _, err := small.Eval("f <- function(n) if (n == 0) identical(x, x) else f(n - 1); f(15)"); err == nil || !strings.Contains(err.Error(), cstack) {
		t.Errorf("30 levels of data below 15 calls: %v", err)
	}
	assertStillAlive(t, small)
	// a counter that a refusal leaves unbalanced would, after many of them,
	// refuse everything
	s := eval.NewSession(eval.Options{MaxDepth: 20})
	for i := 0; i < 200; i++ {
		if _, err := s.Eval("x <- list(); for (i in 1:30) x <- list(x); tryCatch(identical(x, x), error = function(e) 0)"); err != nil {
			t.Fatalf("round %d: %v", i, err)
		}
	}
	assertStillAlive(t, s)
}

// readRDS: a stream that its own session could write (a larger bound) is
// refused by a session with the default one.
func TestDepthBoundOnReadRDS(t *testing.T) {
	deeper := eval.DefaultMaxDepth + 500
	w := eval.NewSession(eval.Options{MaxDepth: 4 * deeper})
	src := "x <- list(); for (i in 1:" + itoa(deeper) + ") x <- list(x); saveRDS(x, '/d.rds', compress = FALSE); readBin('/d.rds', 'raw', 1e8)"
	r, err := w.Eval(src)
	if err != nil {
		t.Fatalf("writing: %v", err)
	}
	s := eval.NewSession(eval.Options{})
	s.Set("bytes", r.Value)
	if _, err := s.Eval("writeBin(bytes, '/d.rds')"); err != nil {
		t.Fatalf("copying: %v", err)
	}
	_, err = s.Eval("readRDS('/d.rds')")
	var e *eval.Error
	if !errors.As(err, &e) || !strings.HasPrefix(e.Message, cstack) {
		t.Fatalf("readRDS: want %q, got %v", cstack, err)
	}
	if cls := evalString(t, s, "paste(class(tryCatch(readRDS('/d.rds'), error = function(e) e)), collapse = ' ')"); cls != "CStackOverflowError stackOverflowError error condition" {
		t.Errorf("class %q", cls)
	}
	assertStillAlive(t, s)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for ; n > 0; n /= 10 {
		b = append([]byte{byte('0' + n%10)}, b...)
	}
	return string(b)
}

// Data below the bound is left alone.
func TestDepthBoundLeavesNormalDataAlone(t *testing.T) {
	s := eval.NewSession(eval.Options{})
	src := "x <- list(); for (i in 1:100) x <- list(x); " +
		"c(identical(x, x), length(unlist(x)) == 0, length(deparse(x)) > 0, length(capture.output(print(x))) > 0, !anyNA(x, recursive = TRUE))"
	r, err := s.Eval(src)
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range r.Value.(*eval.Logical).V {
		if v != 1 {
			t.Errorf("check %d is %d", i, v)
		}
	}
}
