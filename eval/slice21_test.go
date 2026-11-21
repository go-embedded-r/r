package eval

import (
	"math"
	"math/big"
	"math/rand"
	"testing"
)

// The accuracy of the empirical distribution function and of poly, asserted
// against math/big, and the guards of the natives of slice 21.
//
//	ecdf   each y = (count at or below the knot) / n is one correctly rounded
//	       division of exact integers: |y - c/n| <= u |c/n|, u = 2^-53
//	poly   (orthogonal, degree d <= 6, n points) ||Z'Z - I||max <= 64 n u and
//	       ||Z'1||max <= 64 n u, the products formed exactly; the raw columns
//	       x^k have a relative error of at most 8 u
//	predict.poly at the points that built the basis gives the basis back to 4 u

func evalTrue(t *testing.T, s *Session, src string) {
	t.Helper()
	res, err := s.Eval(src)
	if err != nil {
		t.Fatalf("%s: %v", src, err)
	}
	if res.Output != "[1] TRUE\n" {
		t.Fatalf("%s: got %q, want TRUE", src, res.Output)
	}
}

func TestEcdfAgainstBig(t *testing.T) {
	r := rand.New(rand.NewSource(121))
	worst := 0.0
	for trial := 0; trial < 40; trial++ {
		n := 1 + r.Intn(300)
		x := make([]float64, n)
		for i := range x {
			x[i] = float64(r.Intn(50)) / 7 // ties
		}
		s := NewSession(Options{})
		s.Set("x", &Double{V: x})
		linEval(t, s, "k <- ..ecdfknots(x); Fn <- ecdf(x); v <- Fn(k$x)")
		v := linGet(t, s, "v")
		kv, _ := s.Get("k")
		kl := kv.(*List)
		knots := asDoubles(kl.V[0])
		ys := asDoubles(kl.V[1])
		if len(knots) != len(ys) || len(v) != len(ys) {
			t.Fatalf("lengths %d %d %d", len(knots), len(ys), len(v))
		}
		for i, kx := range knots {
			count := 0
			for _, e := range x {
				if e <= kx {
					count++
				}
			}
			want, _ := new(big.Rat).SetFrac64(int64(count), int64(n)).Float64()
			exact := new(big.Float).SetPrec(300).Quo(bigf(float64(count)), bigf(float64(n)))
			diff := new(big.Float).SetPrec(300).Sub(bigf(ys[i]), exact)
			df, _ := diff.Float64()
			rel := math.Abs(df) / (float64(count) / float64(n))
			worst = math.Max(worst, rel/unitRoundoff)
			if ys[i] != want || v[i] != want || rel > unitRoundoff {
				t.Fatalf("n=%d knot %g: y=%g step=%g want %g (rel %g)", n, kx, ys[i], v[i], want, rel)
			}
		}
	}
	t.Logf("ecdf: worst relative error / u = %.3f", worst)
}

func TestEcdfRefusals(t *testing.T) {
	s := NewSession(Options{})
	for _, src := range []string{`..ecdfknots("a")`, `..ecdfknots(numeric())`, `..ecdfknots(c(NA, NaN))`, `..ecdfknots()`, `ecdf(NA_real_)`} {
		if _, err := s.Eval(src); err == nil {
			t.Fatalf("%s was not refused", src)
		}
	}
	evalTrue(t, s, `identical(..ecdfknots(c(2, NA, 1, 2))$n, 3L)`)
}

// An empty numeric argument where a scalar is read is an error, not an index
// panic (found by FuzzEval: qr(,)).
func TestEmptyScalarArgumentsRefused(t *testing.T) {
	s := NewSession(Options{})
	for _, src := range []string{`qr(,)`, `qr(diag(2), tol = numeric())`, `qr.qy(qr(diag(2)), diag(2), numeric())`} {
		if _, err := s.Eval(src); err == nil {
			t.Fatalf("%s was not refused", src)
		}
	}
}

func TestEcdfLine(t *testing.T) {
	s := NewSession(Options{})
	for _, c := range []struct{ src, want string }{
		// the three first against R 4.6.2's print.ecdf
		{`..ecdfline(c(1, 2), 5)`, " x[1:2] =      1,      2\n"},
		{`..ecdfline(c(1.5, 2.5, 3.5), 3)`, " x[1:3] =  1.5,  2.5,  3.5\n"},
		{`..ecdfline(c(-3, 0.5, 2, 2.25, 7, 8, 9, 123456, 1e10), 5)`, " x[1:9] =     -3,    0.5,      2,  ..., 1.2346e+05,  1e+10\n"},
		{`..ecdfline(1, 5)`, " x[1:1] =      1\n"},
		{`..ecdfline(c(1, 2, 3, 4), 1)`, " x[1:4] =  1,  2,  3,  4\n"},
		{`..ecdfline(c(1, 2, 3, 4, 5, 6), 1)`, " x[1:6] =  1,  2,  3,  ...,  5,  6\n"},
		{`..ecdfline(c(1, 2), 99)`, " x[1:2] =      1,      2\n"},
	} {
		res, err := s.Eval(c.src)
		if err != nil {
			t.Fatal(err)
		}
		got := asStrings(res.Value)
		if len(got) != 1 || got[0].S != c.want {
			t.Fatalf("%s: %v, want %q", c.src, got, c.want)
		}
	}
}

func TestPolyBound(t *testing.T) {
	r := rand.New(rand.NewSource(122))
	worst := 0.0
	for trial := 0; trial < 40; trial++ {
		n := 8 + r.Intn(120)
		d := 1 + r.Intn(6)
		x := make([]float64, n)
		for i := range x {
			x[i] = r.NormFloat64()*3 + 10
		}
		s := NewSession(Options{})
		s.Set("x", &Double{V: x})
		s.Set("d", &Double{V: []float64{float64(d)}})
		linEval(t, s, "Z <- unclass(poly(x, d)); dim(Z) <- c(length(x), d); p <- predict(poly(x, d), x); dim(p) <- dim(Z); R <- unclass(poly(x, d, raw = TRUE)); dim(R) <- dim(Z)")
		z := linGet(t, s, "Z")
		zt := transposeData(z, n, d)
		gram := make([]float64, d*d)
		for i := 0; i < d; i++ {
			gram[i+i*d] = 1
		}
		res := exactResidual(zt, d, n, z, d, gram)
		bound := 64 * float64(n) * unitRoundoff
		worst = math.Max(worst, res/bound)
		if res > bound {
			t.Fatalf("n=%d d=%d: ||Z'Z-I|| = %g exceeds %g", n, d, res, bound)
		}
		ones := make([]float64, n)
		for i := range ones {
			ones[i] = 1
		}
		zero := make([]float64, d)
		if c := exactResidual(zt, d, n, ones, 1, zero); c > bound {
			t.Fatalf("n=%d d=%d: ||Z'1|| = %g exceeds %g", n, d, c, bound)
		}
		p := linGet(t, s, "p")
		for i := range z {
			if math.Abs(p[i]-z[i]) > 4*unitRoundoff*math.Max(1, math.Abs(z[i])) {
				t.Fatalf("predict differs from the basis at %d: %g vs %g", i, p[i], z[i])
			}
		}
		raw := linGet(t, s, "R")
		for k := 1; k <= d; k++ {
			for i := 0; i < n; i++ {
				ex := new(big.Float).SetPrec(300).SetFloat64(1)
				for j := 0; j < k; j++ {
					ex.Mul(ex, bigf(x[i]))
				}
				df, _ := new(big.Float).SetPrec(300).Sub(bigf(raw[i+(k-1)*n]), ex).Float64()
				exf, _ := ex.Float64()
				if math.Abs(df) > 8*unitRoundoff*math.Abs(exf) {
					t.Fatalf("raw x^%d at %g: %g, exact %g", k, x[i], raw[i+(k-1)*n], exf)
				}
			}
		}
	}
	t.Logf("poly: worst ||Z'Z - I|| / (64 n u) = %.4f", worst)
}

// The guards of restartSpec: only a list with a handler is split.
func TestRestartSpec(t *testing.T) {
	fn := &Closure{}
	if f, d := restartSpec(fn); f != Value(fn) || d != nil {
		t.Fatalf("a function: got %v %v", f, d)
	}
	if f, d := restartSpec(mkString("x")); f == nil || d != nil {
		t.Fatalf("a string: got %v %v", f, d)
	}
	named := func(names []string, vs ...Value) *List {
		return &List{V: vs, Attr: []Attr{{Name: "names", Value: strVec(names)}}}
	}
	noHandler := named([]string{"description"}, mkString("d"))
	if f, d := restartSpec(noHandler); f != Value(noHandler) || d != nil {
		t.Fatalf("no handler: got %v %v", f, d)
	}
	unnamed := &List{V: []Value{fn}}
	if f, d := restartSpec(unnamed); f != Value(unnamed) || d != nil {
		t.Fatalf("unnamed: got %v %v", f, d)
	}
	both := named([]string{"handler", "description"}, fn, mkString("d"))
	if f, d := restartSpec(both); f != Value(fn) || d == nil {
		t.Fatalf("handler and description: got %v %v", f, d)
	}
	only := named([]string{"handler"}, fn)
	if f, d := restartSpec(only); f != Value(fn) || d != nil {
		t.Fatalf("handler: got %v %v", f, d)
	}
}

func TestSlice21Guards(t *testing.T) {
	s := NewSession(Options{})
	// ..approxeval: an argument of length 0 reads as NA, one of length 1 as itself
	evalTrue(t, s, `is.na(..approxeval(1:3, c(2, 4, 8), 0, 1, numeric(0), 9, 0))`)
	evalTrue(t, s, `is.na(..approxeval(1:3, c(2, 4, 8), 5, 1, 0, numeric(0), 0))`)
	evalTrue(t, s, `..approxeval(1:3, c(2, 4, 8), 0, 1, 7, 9, 0) == 7`)
	evalTrue(t, s, `..approxeval(1:3, c(2, 4, 8), 5, 1, 7, 9, 0) == 9`)
	evalTrue(t, s, `..approxeval(1:3, c(2, 4, 8), 1.5, 1, 7, 9, 0) == 3`)
	evalTrue(t, s, `..approxeval(1:3, c(2, 4, 8), 2.5, 2, 7, 9, 0.5) == 6`)
	evalTrue(t, s, `is.na(..approxeval(numeric(), numeric(), 1, 1, 0, 1, 0))`)
	// sys.parents: empty at the top level, the parent numbers inside calls
	evalTrue(t, s, `length(sys.parents()) == 0L`)
	evalTrue(t, s, `f <- function() sys.parents(); identical(f(), 0L)`)
	evalTrue(t, s, `g <- function() f(); identical(g(), c(0L, 1L))`)
	// pairlist: no argument is NULL, otherwise a pairlist
	evalTrue(t, s, `is.null(pairlist())`)
	evalTrue(t, s, `identical(typeof(pairlist(a = 1)), "pairlist")`)
	evalTrue(t, s, `is.null(..mkpair(1))`)
	evalTrue(t, s, `is.null(..mkpair(list()))`)
	evalTrue(t, s, `identical(typeof(..mkpair(list(a = 1))), "pairlist")`)
	evalTrue(t, s, `identical(names(..mkpair(list(a = 1))), "a")`)
	if _, err := s.Eval(`..mkpair()`); err == nil {
		t.Fatal("..mkpair() was not refused")
	}
	if _, err := s.Eval(`..approxeval(1)`); err == nil {
		t.Fatal("..approxeval(1) was not refused")
	}
}

// regmatches(invert = NA) interleaves the gaps and the matches, and is.pairlist
// is true for a pairlist (values checked against R 4.6.2).
func TestRegmatchesInvertNAAndPairlist(t *testing.T) {
	s := NewSession(Options{})
	for _, src := range []string{
		`identical(regmatches(c("a1b22", "zz"), gregexpr("[0-9]+", c("a1b22", "zz")), invert = NA), list(c("a", "1", "b", "22", ""), "zz"))`,
		`identical(regmatches(c("a1b", "zz", NA), regexpr("[0-9]", c("a1b", "zz", NA)), invert = NA)[1:2], list(c("a", "1", "b"), "zz"))`,
		`is.na(regmatches(NA_character_, regexpr("a", NA_character_), invert = NA)[[1]])`,
		`is.pairlist(pairlist(a = 1)) && is.pairlist(NULL) && !is.pairlist(list(1))`,
		`identical(formatDL(structure(c("first description", "second"), names = c("b", "gamma"))), formatDL(c("b", "gamma"), c("first description", "second")))`,
		`is.integer(environment(ecdf(c(3, 1, 2)))$nobs)`,
	} {
		evalTrue(t, s, src)
	}
	if _, err := s.Eval(`read.fwf(textConnection("abc"), c(1, 1), skip = 1)`); err == nil {
		t.Fatal("read.fwf with no data line should fail")
	}
}
