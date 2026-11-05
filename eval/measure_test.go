package eval

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"
	"time"
)

// TestMeasureLimits reports the largest steps, allocation and time of the
// golden snippets, to choose the defaults of Options. Run it with
// RGO_MEASURE=1; it is not a gate.
func TestMeasureLimits(t *testing.T) {
	if os.Getenv("RGO_MEASURE") == "" {
		t.Skip("set RGO_MEASURE=1")
	}
	f, err := os.Open("../testdata/golden/eval.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<26)
	var maxSteps int
	var maxAlloc int64
	var maxTime time.Duration
	var sSrc, aSrc string
	n := 0
	srcs := workloads
	for sc.Scan() {
		var rec struct{ Src string }
		if json.Unmarshal(sc.Bytes(), &rec) != nil {
			continue
		}
		srcs = append(srcs, rec.Src)
	}
	for _, src := range srcs {
		rec := struct{ Src string }{src}
		s := NewSession(Options{MaxSteps: 1 << 30, MaxAlloc: 1<<31 - 1, Timeout: Unbounded})
		t0 := time.Now()
		s.Eval(rec.Src)
		d := time.Since(t0)
		if n < len(workloads) {
			t.Logf("workload %q: steps %d alloc %d time %v", src, s.steps, s.alloc, d)
		}
		n++
		if s.steps > maxSteps {
			maxSteps, sSrc = s.steps, rec.Src
		}
		if s.alloc > maxAlloc {
			maxAlloc, aSrc = s.alloc, rec.Src
		}
		if d > maxTime {
			maxTime = d
		}
	}
	t.Logf("%d snippets: max steps %d (%q), max alloc %d bytes (%q), max time %v", n, maxSteps, sSrc, maxAlloc, aSrc, maxTime)
}

// workloads are cells a notebook user writes on purpose: loops of a million
// iterations, naive recursion, vectors of ten million doubles.
var workloads = []string{
	"s <- 0; for (i in 1:1000000) s <- s + i; s",
	"f <- function(n) if (n < 2) n else f(n - 1) + f(n - 2); f(22)",
	"x <- numeric(1e7); x[1] <- 1; sum(x)",
	"x <- 1:1e7; sum(as.numeric(x))",
	"f <- function(n) if (n == 0) 0 else 1 + f(n - 1); f(1000)",
	"x <- numeric(10000); for (i in 1:10000) x[i] <- i^2; sum(x)",
	"l <- list(); for (i in 1:20000) l[[i]] <- i; length(l)",
	"i <- 0; while (i < 500000) i <- i + 1; i",
}
