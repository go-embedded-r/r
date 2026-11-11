package eval_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-embedded-r/r/eval"
	"github.com/go-embedded-r/r/syntax"
)

// Bombs of the raw, complex and expression types, and of the binary
// connection functions: every allocation is bounded before it is made.

func TestRawComplexBombs(t *testing.T) {
	underMemoryLimit(t)
	s := eval.NewSession(eval.Options{MaxAlloc: 1 << 30})
	for _, src := range []string{
		"raw(1e12)", "raw(2^28 + 1)", "complex(1e12)", "complex(length.out = 1e12)", "vector(\"raw\", 1e12)", "vector(\"complex\", 1e12)",
		"rep(as.raw(1), 1e12)", "rep(1i, 1e12)", "as.raw(1:2)[1e12]", "x <- raw(2^28); c(x, x)", "rawToChar(raw(2^28 + 1))",
		"charToRaw(strrep(\"a\", 1e12))", "intToBits(1:2^24)", "rawToBits(raw(2^26))", "numeric(2^28) + 1i",
		"matrix(raw(1), 1e6, 1e6)", "readBin(as.raw(1:8), \"integer\", n = 1e12)", "readBin(as.raw(1:8), \"character\", n = 1e12)",
		"readChar(as.raw(1:8), 1e12)", "writeChar(\"a\", raw(), nchars = 1e12)", "expression(1)[1e12]", "length(x <- expression(a)) <- 1e12",
	} {
		refused(t, s, src, "allocate", "too large", "limit", "invalid", "nchars", "memory")
	}
}

// writeBin and the raw connection are bounded by the allocation budget and by
// MaxFileBytes, and a huge exponent of a complex number costs no loop.
func TestBinaryOutputBombs(t *testing.T) {
	underMemoryLimit(t)
	small := eval.NewSession(eval.Options{MaxAlloc: 64 << 20})
	refused(t, small, "writeBin(numeric(5e6), raw())", "memory limit")
	refused(t, small, "rawToChar(rep(as.raw(65), 9e6))", "memory limit")
	files := eval.NewSession(eval.Options{MaxFileBytes: 1 << 20})
	refused(t, files, "con <- rawConnection(raw(0), \"w\"); writeBin(raw(2e6), con)", "file size limit")
	refused(t, files, "con <- rawConnection(raw(0), \"w\"); for (i in 1:1e6) writeChar(strrep(\"a\", 100), con)", "file size limit", "step limit")
	refused(t, files, "f <- tempfile(); writeBin(raw(2e6), f)", "file size limit")
	for _, src := range []string{"(1+1i)^1e9", "(1+1i)^65536L", "(1+1i)^-1e300", "2^(1e300+1e300i)"} {
		done := make(chan error, 1)
		go func() { _, err := files.Eval(src); done <- err }()
		select {
		case err := <-done:
			if err != nil {
				t.Errorf("%s: %v", src, err)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("%s did not end", src)
		}
	}
}

// parse(text=) obeys the parse limits of the session.
func TestParseTextBombs(t *testing.T) {
	underMemoryLimit(t)
	s := eval.NewSession(eval.Options{Parse: syntax.Limits{MaxBytes: 1 << 14, MaxDepth: 100, MaxNodes: 5000}})
	refused(t, s, "parse(text = strrep(\"(\", 1e5))", "limit", "exceeded")
	refused(t, s, "parse(text = strrep(\"a;\", 1e5))", "limit", "exceeded")
	refused(t, s, "str2lang(strrep(\"(\", 1e4))", "limit", "exceeded", "contextstack overflow")
}

// What a decoder is given is untrusted: readBin and readChar over every size
// and byte order on bytes of any shape ends with a value or an R error.
func TestReadBinNeverPanics(t *testing.T) {
	underMemoryLimit(t)
	var raws []string
	for _, n := range []int{0, 1, 2, 3, 4, 7, 8, 9, 16, 17, 33} {
		parts := make([]string, n)
		for i := range parts {
			parts[i] = fmt.Sprint((i*37 + n*11) % 256)
		}
		raws = append(raws, "as.raw(c("+strings.Join(parts, ",")+"))")
	}
	raws = append(raws, "as.raw(rep(255L, 16))", "as.raw(rep(0L, 16))", "as.raw(c(97, 0, 98, 0, 0))", "as.raw(c(255, 254, 253))")
	whats := []string{"\"integer\"", "\"double\"", "\"numeric\"", "\"logical\"", "\"complex\"", "\"character\"", "\"raw\"", "integer()", "TRUE", "1i", "\"foo\""}
	s := eval.NewSession(eval.Options{MaxVector: 1 << 16, MaxAlloc: 1 << 28, MaxFileBytes: 1 << 20, Timeout: 2 * time.Second})
	for _, r := range raws {
		for _, w := range whats {
			for _, extra := range []string{"", ", n = 3", ", n = 0", ", size = 1", ", size = 2", ", size = 4", ", size = 8", ", size = 3", ", size = 16", ", signed = FALSE, size = 1", ", signed = FALSE, size = 2", ", endian = \"big\"", ", endian = \"swap\", n = 100", ", n = 1e12", ", n = -1", ", size = NA, n = 2"} {
				src := "readBin(" + r + ", " + w + extra + ")"
				func() {
					defer func() {
						if p := recover(); p != nil {
							t.Errorf("%s: panic %v", src, p)
						}
					}()
					_, _ = s.Eval(src)
				}()
			}
		}
		for _, n := range []string{"1", "c(1, 2, 3)", "0", "100", "-1", "NA", "1e12"} {
			for _, ub := range []string{"", ", useBytes = TRUE"} {
				src := "readChar(" + r + ", " + n + ub + ")"
				func() {
					defer func() {
						if p := recover(); p != nil {
							t.Errorf("%s: panic %v", src, p)
						}
					}()
					_, _ = s.Eval(src)
				}()
			}
		}
	}
	if _, err := s.Eval("NULL"); err != nil {
		t.Fatalf("session broken: %v", err)
	}
}
