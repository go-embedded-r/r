package eval_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/go-embedded-r/r/eval"
	"github.com/go-embedded-r/r/internal/gen"
	"github.com/go-embedded-r/r/syntax"
)

func TestSessionPersists(t *testing.T) {
	s := eval.NewSession(eval.Options{})
	if _, err := s.Eval("x <- 40; f <- function(y) y + x"); err != nil {
		t.Fatal(err)
	}
	r, err := s.Eval("f(2)")
	if err != nil {
		t.Fatal(err)
	}
	if d := r.Value.(*eval.Double); d.V[0] != 42 || !r.Visible {
		t.Errorf("%+v visible=%v", d, r.Visible)
	}
	if v, ok := s.Get("x"); !ok || v.(*eval.Double).V[0] != 40 {
		t.Errorf("Get: %v %v", v, ok)
	}
	s.Set("z", eval.Null)
	if _, err := s.Eval("is.null(z)"); err != nil {
		t.Error(err)
	}
	r, _ = s.Eval("y <- 1")
	if r.Visible {
		t.Error("an assignment is invisible")
	}
	if names := s.Global.Names(); strings.Join(names, ",") != "x,f,z,y" {
		t.Errorf("names %v", names)
	}
}

// R's do_radixsort reads na.last with asLogical, which is NA for an empty
// vector; the evaluator used to index element 0 of it (FuzzEval crasher from
// the main CI lane).
func TestOrderEmptyNaLast(t *testing.T) {
	s := eval.NewSession(eval.Options{})
	r, err := s.Eval("order(c(2, NA, 1), na.last = logical(0))")
	if err != nil {
		t.Fatal(err)
	}
	in, ok := r.Value.(*eval.Integer)
	if !ok || len(in.V) != 2 || in.V[0] != 3 || in.V[1] != 1 {
		t.Errorf("got %#v, want c(3L, 1L) as for na.last = NA", r.Value)
	}
}

func TestSessionErrors(t *testing.T) {
	s := eval.NewSession(eval.Options{})
	_, err := s.Eval("f <- function(a) a + undefined_name; f(1)")
	var e *eval.Error
	if !errors.As(err, &e) || e.Message != "object 'undefined_name' not found" || e.Call != "f(1)" {
		t.Fatalf("%#v", err)
	}
	if !strings.HasPrefix(err.Error(), "Error in f(1) : object") {
		t.Error(err)
	}
	// the session survives an error, and keeps what ran before it
	if _, err := s.Eval("g <- 2; stop_here()"); err == nil {
		t.Fatal("no error")
	}
	if _, ok := s.Get("g"); !ok {
		t.Error("g lost")
	}
	if _, err := s.Eval("g * 2"); err != nil {
		t.Error(err)
	}
	_, err = s.Eval("1 +")
	var pe *syntax.Error
	if !errors.As(err, &pe) {
		t.Errorf("a parse error is a *syntax.Error: %T", err)
	}
	_, err = s.Eval("1 + 1")
	if err != nil {
		t.Error(err)
	}
	var top *eval.Error
	_, err = s.Eval("stop_at_top")
	if !errors.As(err, &top) || top.Call != "" || err.Error() != "Error: object 'stop_at_top' not found" {
		t.Errorf("%v", err)
	}
}

func TestSessionWarnings(t *testing.T) {
	s := eval.NewSession(eval.Options{})
	r, err := s.Eval("1:3 + 1:2; sqrt(-1)")
	if err != nil || len(r.Warnings) != 2 || r.Warnings[1].Message != "NaNs produced" || r.Warnings[1].Call != "sqrt(-1)" {
		t.Fatalf("%+v %v", r, err)
	}
	r, _ = s.Eval("1")
	if len(r.Warnings) != 0 {
		t.Error("warnings are per Eval")
	}
}

func TestLimits(t *testing.T) {
	s := eval.NewSession(eval.Options{MaxDepth: 100})
	_, err := s.Eval("f <- function(n) f(n + 1); f(1)")
	var e *eval.Error
	if !errors.As(err, &e) || !strings.Contains(e.Message, "nested too deeply") {
		t.Errorf("deep recursion: %v", err)
	}
	// the default limit is R's 5000, and recursion that stays under it works
	s = eval.NewSession(eval.Options{})
	if r, err := s.Eval("f <- function(n) if (n == 0) 0 else 1 + f(n - 1); f(500)"); err != nil || r.Value.(*eval.Double).V[0] != 500 {
		t.Errorf("recursion 500: %v", err)
	}
	if _, err := s.Eval("f <- function(n) f(n + 1); f(1)"); err == nil {
		t.Error("infinite recursion accepted")
	}
	s = eval.NewSession(eval.Options{MaxVector: 1000})
	if _, err := s.Eval("1:100000"); err == nil || !strings.Contains(err.Error(), "cannot allocate") {
		t.Errorf("vector limit: %v", err)
	}
	s = eval.NewSession(eval.Options{MaxSteps: 1000})
	if _, err := s.Eval("repeat 1 + 1"); err == nil || !strings.Contains(err.Error(), "step limit") {
		t.Errorf("step limit: %v", err)
	}
	s = eval.NewSession(eval.Options{Parse: syntax.Limits{MaxBytes: 10}})
	if _, err := s.Eval("1 + 1 + 1 + 1 + 1 + 1"); !errors.Is(err, syntax.ErrLimit) {
		t.Errorf("parse limit: %v", err)
	}
}

func TestLoopsAreIterative(t *testing.T) {
	s := eval.NewSession(eval.Options{MaxDepth: 50})
	r, err := s.Eval("x <- 0; for (i in 1:100000) x <- x + i; x")
	if err != nil || r.Value.(*eval.Double).V[0] != 5000050000 {
		t.Errorf("%v %v", r, err)
	}
}

// Every primitive and native function the evaluator registers is exercised by
// the goldens, so each is checked against the oracle.
func TestPrimitiveCoverage(t *testing.T) {
	prims := syntax.PrimitiveNames()
	snippets, err := gen.EvalInputs("..")
	if err != nil {
		t.Fatal(err)
	}
	used := map[string]bool{}
	for _, src := range snippets {
		f, err := syntax.ParseAll(src)
		if err != nil {
			continue
		}
		var walk func(syntax.Node)
		walk = func(n syntax.Node) {
			switch x := n.(type) {
			case *syntax.Call:
				if s, ok := x.Fun.(*syntax.Sym); ok {
					used[s.Name] = true
					if (s.Name == "<-" || s.Name == "=" || s.Name == "<<-") && len(x.Args) == 2 {
						// f(g(x), i) <- v calls `f<-`(...) and then `g<-`(...)
						for t := x.Args[0].Value; ; {
							tc, ok := t.(*syntax.Call)
							if !ok {
								break
							}
							if fs, ok := tc.Fun.(*syntax.Sym); ok {
								used[fs.Name+"<-"] = true
							}
							if len(tc.Args) == 0 {
								break
							}
							t = tc.Args[0].Value
						}
					}
				} else {
					walk(x.Fun)
				}
				for _, a := range x.Args {
					walk(a.Value)
				}
			case *syntax.Function:
				used["function"] = true
				for _, p := range x.Params {
					walk(p.Default)
				}
				walk(x.Body)
			}
		}
		for _, e := range f.Exprs {
			walk(e)
		}
	}
	have := eval.Builtins()
	var missing []string
	natives := 0
	refuse := map[string]bool{}
	for _, n := range eval.Unsupported() {
		refuse[n] = true
	}
	implemented := 0
	for _, n := range have {
		if refuse[n] {
			continue
		}
		if strings.HasPrefix(n, "..") && !strings.HasPrefix(n, "...") {
			// a helper of ours, not a function of R: its callers are what the goldens exercise
			continue
		}
		implemented++
		if !prims[n] {
			natives++
		}
		if !used[n] {
			missing = append(missing, n)
		}
	}
	if len(missing) > 0 {
		t.Errorf("registered but not exercised by testdata/eval: %v", missing)
	}
	t.Logf("%d of the %d primitives of R are implemented, and %d base functions that are closures in R; each is covered by the oracle goldens", implemented-natives, len(prims), natives)
}

// The primitives that need a type or runtime we do not have are refused with
// an error that says so; they are not counted as implemented.
func TestUnsupportedPrimitivesRefuse(t *testing.T) {
	for _, n := range eval.Unsupported() {
		s := eval.NewSession(eval.Options{})
		_, err := s.Eval("f <- `" + n + "`; f()")
		if err == nil || !strings.Contains(err.Error(), "not supported") && !strings.Contains(err.Error(), "argument") {
			t.Errorf("%s: want a refusal, got %v", n, err)
		}
	}
}
