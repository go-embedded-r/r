package eval

import "testing"

// A zero-length argument gives a zero-length result that keeps the attributes
// of the first argument (R's math2 and the Bessel functions), except lbeta,
// which is a closure in R and drops them (oracle R 4.6.2).
func TestSpecialEmptyAttrs(t *testing.T) {
	for _, c := range []struct{ src, want string }{
		{`names(attributes(beta(structure(numeric(0), foo = "a"), 1)))`, "[1] \"foo\"\n"},
		{`names(attributes(psigamma(structure(numeric(0), foo = "a"), 1)))`, "[1] \"foo\"\n"},
		{`names(attributes(besselJ(structure(numeric(0), foo = "a"), 1)))`, "[1] \"foo\"\n"},
		{`names(attributes(besselY(structure(numeric(0), foo = "a"), 1)))`, "[1] \"foo\"\n"},
		{`names(attributes(besselI(structure(numeric(0), foo = "a"), 1)))`, "[1] \"foo\"\n"},
		{`names(attributes(besselK(structure(numeric(0), foo = "a"), 1)))`, "[1] \"foo\"\n"},
		{`is.null(attributes(lbeta(1, structure(numeric(0), foo = "a"))))`, "[1] TRUE\n"},
		{`length(lbeta(numeric(0), 1)) + length(beta(1, numeric(0)))`, "[1] 0\n"},
	} {
		res, err := NewSession(Options{}).Eval(c.src)
		if err != nil {
			t.Fatalf("%s: %v", c.src, err)
		}
		if res.Output != c.want {
			t.Errorf("%s printed %q, want %q", c.src, res.Output, c.want)
		}
	}
}
