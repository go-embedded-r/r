package eval_test

import (
	"testing"

	"github.com/go-embedded-r/r/eval"
)

// BenchmarkGsubGroup is gsub with a capture group over 20 000 lines, the
// case where the loop around the engine, not the engine, dominated.
func BenchmarkGsubGroup(b *testing.B) {
	s := eval.NewSession(eval.Options{MaxSteps: eval.Unbounded, MaxAlloc: eval.Unbounded, Timeout: eval.Unbounded})
	if _, err := s.Eval(`x <- rep(c("the quick brown fox jumps over the lazy dog", "error: file not found at line 42", "2026-11-07 12:00:01 INFO started job 17"), length.out = 20000)`); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := s.Eval(`y <- gsub("([a-z]+) ([a-z]+)", "\\2 \\1", x)`); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGrepl(b *testing.B) {
	s := eval.NewSession(eval.Options{MaxSteps: eval.Unbounded, MaxAlloc: eval.Unbounded, Timeout: eval.Unbounded})
	if _, err := s.Eval(`x <- rep(c("the quick brown fox jumps over the lazy dog", "error: file not found at line 42", "2026-11-07 12:00:01 INFO started job 17"), length.out = 20000)`); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := s.Eval(`y <- grepl("error", x)`); err != nil {
			b.Fatal(err)
		}
	}
}
