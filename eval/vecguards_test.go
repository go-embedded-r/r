package eval_test

import "testing"

// The guards of ':', 'c' and the type names, with values on both sides.
// Expected values were settled with the oracle (R 4.6.2).
func TestVecGuards(t *testing.T) {
	s := newRandSession()
	run := func(src string) string {
		t.Helper()
		res, err := s.Eval(src)
		if err != nil {
			t.Fatalf("%q: %v", src, err)
		}
		return res.Output
	}
	// Inf:Inf subtracts to NaN: numeric(0), not a vector of garbage length
	if got := run(`length(Inf:Inf); length(-Inf:-Inf); length(1:3)`); got != "[1] 0\n[1] 0\n[1] 3\n" {
		t.Errorf("Inf:Inf: %q", got)
	}
	// a string that is no number is NA: an error, not a zero
	refused(t, s, `suppressWarnings("a":3)`, "NA/NaN argument")
	refused(t, s, `3:suppressWarnings("a")`, "NA/NaN argument")
	if got := run(`length("2":4)`); got != "[1] 3\n" {
		t.Errorf("string range: %q", got)
	}
	// the type of a special and of a builtin
	if got := run(`typeof(` + "`if`" + `); typeof(sum); tryCatch(as.numeric(` + "`if`" + `), error = conditionMessage); tryCatch(as.numeric(sum), error = conditionMessage)`); got != "[1] \"special\"\n[1] \"builtin\"\n[1] \"cannot coerce type 'special' to vector of type 'double'\"\n[1] \"cannot coerce type 'builtin' to vector of type 'double'\"\n" {
		t.Errorf("types: %q", got)
	}
	// names of c(): none without a tag, "" for an untagged element among tagged ones
	if got := run(`names(c(1, 2)); names(c(a = 1, 2)); names(c(1, b = 2))`); got != "NULL\n[1] \"a\" \"\" \n[1] \"\"  \"b\"\n" {
		t.Errorf("names: %q", got)
	}
}
