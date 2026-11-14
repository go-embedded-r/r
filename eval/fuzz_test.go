package eval_test

import (
	"testing"
	"time"

	"github.com/go-embedded-r/r/eval"
	"github.com/go-embedded-r/r/internal/gen"
	"github.com/go-embedded-r/r/syntax"
)

// FuzzEval: whatever source is given, evaluation ends with a value or an error
// and never panics, hangs or allocates without bound.
func FuzzEval(f *testing.F) {
	recs, err := gen.ReadEvalGoldens("..")
	if err != nil {
		f.Fatal(err)
	}
	for _, r := range recs {
		f.Add(r.Src)
	}
	f.Fuzz(func(t *testing.T, src string) {
		s := eval.NewSession(eval.Options{
			MaxDepth: 200, MaxVector: 1 << 16, MaxSteps: 20000, MaxAlloc: 1 << 28, MaxFileBytes: 1 << 22, Timeout: 2 * time.Second,
			Parse: syntax.Limits{MaxBytes: 1 << 14, MaxDepth: 100, MaxNodes: 5000},
		})
		s.Eval(src)
		// the session is usable afterwards (the source may have redefined anything)
		if _, err := s.Eval("NULL"); err != nil {
			t.Fatalf("session broken after %q: %v", src, err)
		}
	})
}
