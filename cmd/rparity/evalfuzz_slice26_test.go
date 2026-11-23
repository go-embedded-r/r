package main

import (
	"math/rand"
	"regexp"
	"strings"
	"testing"
)

// TestFuzzSlice26Shape holds the generator's shape: every snippet names a
// family on its first line, ordinary calls end in "# libm" after nine digits and
// edge calls after six, and a first argument is sometimes one bare value and
// sometimes a single edge value.
func TestFuzzSlice26Shape(t *testing.T) {
	r := rand.New(rand.NewSource(26))
	known := map[string]bool{}
	for _, f := range fuzzFams26 {
		known[f.name] = true
	}
	var edgeSix, libm, bareFirst, edgeFirst int
	bareRe := regexp.MustCompile(`^signif\([dpq][a-z]+\(([^c(][^,)]*)[,)]`)
	edgeRe := regexp.MustCompile(`^signif\([dpq][a-z]+\((NA|NaN|Inf|-Inf|1e30|-1e30)[,)]`)
	for i := 0; i < 6000; i++ {
		s := fuzzSlice26(r)
		head, call, ok := strings.Cut(s, "\n")
		if !ok || !strings.HasPrefix(head, "# ") || !known[strings.TrimPrefix(head, "# ")] {
			t.Fatalf("snippet %d: bad first line in %q", i, s)
		}
		switch {
		case strings.HasSuffix(call, ", 9) # libm"):
			libm++
		case strings.HasSuffix(call, ", 6)"):
			edgeSix++
		default:
			t.Fatalf("snippet %d: neither ordinary nor edge ending: %q", i, s)
		}
		if bareRe.MatchString(call) {
			bareFirst++
		}
		if edgeRe.MatchString(call) {
			edgeFirst++
		}
	}
	for name, n := range map[string]int{"ordinary": libm, "edge": edgeSix, "bare first argument": bareFirst, "edge first argument": edgeFirst} {
		if n == 0 {
			t.Errorf("no %s snippet in 6000", name)
		}
	}
	// the edge share is about a quarter
	if edgeSix < 1000 || edgeSix > 2000 {
		t.Errorf("edge snippets: %d of 6000, want about 1500", edgeSix)
	}
}
