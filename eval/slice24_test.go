package eval_test

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-embedded-r/r/eval"
)

// Slice 24: random streams. The goldens (testdata/eval/slice24*.R) pin the
// stream to R; these tests pin what the goldens cannot: the bound on rnorm
// against an extended-precision oracle, the refusals that are ours (R reads
// out of its state there) and the state transitions.

const bigPrec = 400

// piDigits is pi to 100 digits.
const piDigits = "3.1415926535897932384626433832795028841971693993751058209749445923078164062862089986280348253421170679"

// bigPhi is the standard normal distribution function at x, by the Taylor
// series of erf at bigPrec bits (alternating, so no cancellation beyond what
// 400 bits absorb for |x| < 8).
func bigPhi(x *big.Float) *big.Float {
	pi, _, _ := big.ParseFloat(piDigits, 10, bigPrec, big.ToNearestEven)
	two := big.NewFloat(2).SetPrec(bigPrec)
	z := new(big.Float).SetPrec(bigPrec).Quo(x, new(big.Float).SetPrec(bigPrec).Sqrt(two))
	z2 := new(big.Float).SetPrec(bigPrec).Mul(z, z)
	term := new(big.Float).SetPrec(bigPrec).Set(z) // z^(2n+1)/n! with sign
	sum := new(big.Float).SetPrec(bigPrec).Set(z)
	for n := 1; n < 160; n++ {
		term.Mul(term, z2)
		term.Quo(term, big.NewFloat(float64(-n)).SetPrec(bigPrec))
		t := new(big.Float).SetPrec(bigPrec).Quo(term, big.NewFloat(float64(2*n+1)).SetPrec(bigPrec))
		sum.Add(sum, t)
	}
	erf := new(big.Float).SetPrec(bigPrec).Quo(sum, new(big.Float).SetPrec(bigPrec).Sqrt(pi))
	erf.Mul(erf, two)
	r := new(big.Float).SetPrec(bigPrec).Add(big.NewFloat(1).SetPrec(bigPrec), erf)
	return r.Quo(r, two)
}

// bigPPF inverts bigPhi at p by bisection around the float64 guess g.
func bigPPF(p, g float64) float64 {
	bp := new(big.Float).SetPrec(bigPrec).SetFloat64(p)
	w := 1e-5 * math.Max(1, math.Abs(g))
	lo := new(big.Float).SetPrec(bigPrec).SetFloat64(g - w)
	hi := new(big.Float).SetPrec(bigPrec).SetFloat64(g + w)
	half := big.NewFloat(0.5).SetPrec(bigPrec)
	for i := 0; i < 90; i++ {
		mid := new(big.Float).SetPrec(bigPrec).Add(lo, hi)
		mid.Mul(mid, half)
		if bigPhi(mid).Cmp(bp) < 0 {
			lo = mid
		} else {
			hi = mid
		}
	}
	mid := new(big.Float).SetPrec(bigPrec).Add(lo, hi)
	mid.Mul(mid, half)
	f, _ := mid.Float64()
	return f
}

func numbers(t *testing.T, s *eval.Session, src string) []float64 {
	t.Helper()
	res, err := s.Eval(src)
	if err != nil {
		t.Fatalf("%q: %v", src, err)
	}
	var out []float64
	for _, f := range strings.Fields(res.Output) {
		v, err := strconv.ParseFloat(f, 64)
		if err != nil {
			t.Fatalf("%q: %q is not a number", src, f)
		}
		out = append(out, v)
	}
	return out
}

func newRandSession() *eval.Session {
	return eval.NewSession(eval.Options{MaxVector: 1 << 22, MaxAlloc: 1 << 28, Timeout: 60 * time.Second})
}

// rnorm by Inversion is within 1 ulp of the exact inverse of the normal
// distribution function at the uniform it actually drew (two words of the
// stream make one uniform).
func TestRnormAgainstBig(t *testing.T) {
	const n = 30
	s := newRandSession()
	for _, seed := range []int{1, 2} {
		xs := numbers(t, s, fmt.Sprintf(`set.seed(%d); cat(sprintf("%%.17g", rnorm(%d)), sep = "\n")`, seed, n))
		us := numbers(t, s, fmt.Sprintf(`set.seed(%d); cat(sprintf("%%.17g", runif(%d)), sep = "\n")`, seed, 2*n))
		if len(xs) != n || len(us) != 2*n {
			t.Fatalf("got %d and %d numbers", len(xs), len(us))
		}
		for i, x := range xs {
			const b = 134217728.0
			u := (math.Floor(b*us[2*i]) + us[2*i+1]) / b
			want := bigPPF(u, x)
			ulp := math.Nextafter(math.Abs(want), math.Inf(1)) - math.Abs(want)
			if d := math.Abs(x - want); d > ulp {
				t.Errorf("seed %d draw %d: rnorm %.17g, exact %.17g at u=%.17g (%.2g ulp)", seed, i, x, want, u, d/ulp)
			}
		}
	}
}

func TestRandomSeedStateTransitions(t *testing.T) {
	s := newRandSession()
	run := func(src string) string {
		t.Helper()
		res, err := s.Eval(src)
		if err != nil {
			t.Fatalf("%q: %v", src, err)
		}
		return res.Output
	}
	// no state until a draw needs one
	run(`if (exists(".Random.seed", envir = globalenv())) rm(list = ".Random.seed", envir = globalenv())`)
	if got := run(`exists(".Random.seed")`); got != "[1] FALSE\n" {
		t.Fatalf("got %q", got)
	}
	if got := run(`invisible(runif(0)); exists(".Random.seed")`); got != "[1] FALSE\n" {
		t.Fatalf("runif(0) created a state: %q", got)
	}
	if got := run(`invisible(rnorm(0)); exists(".Random.seed")`); got != "[1] FALSE\n" {
		t.Fatalf("rnorm(0) created a state: %q", got)
	}
	if got := run(`invisible(RNGkind()); exists(".Random.seed")`); got != "[1] FALSE\n" {
		t.Fatalf("RNGkind() created a state: %q", got)
	}
	// a draw makes a plain integer vector of 626: kind, position, 624 words
	if got := run(`invisible(runif(1)); c(length(.Random.seed), typeof(.Random.seed) == "integer", .Random.seed[1])`); got != "[1]   626     1 10403\n" {
		t.Fatalf("got %q", got)
	}
	if got := run(`set.seed(1); .Random.seed[1:2]`); got != "[1] 10403   624\n" {
		t.Fatalf("got %q", got)
	}
	// drawing advances the position; restoring the vector replays the stream
	if got := run(`set.seed(9); a <- .Random.seed; x <- runif(5); .Random.seed <- a; identical(x, runif(5))`); got != "[1] TRUE\n" {
		t.Fatalf("got %q", got)
	}
	if got := run(`set.seed(9); invisible(runif(3)); .Random.seed[2]`); got != "[1] 3\n" {
		t.Fatalf("got %q", got)
	}
	// a position of 0, 624 and beyond 625 regenerate, as 624 does
	for _, p := range []string{"0L", "624L", "700L", "1000000L", "2147483647L"} {
		got := run(fmt.Sprintf(`set.seed(1); s <- .Random.seed; s[2] <- %s; .Random.seed <- s; a <- runif(2); set.seed(1); identical(a, runif(2))`, p))
		if got != "[1] TRUE\n" {
			t.Errorf("position %s: %q", p, got)
		}
	}
	// the end of the table: position 623 hands out one word, then regenerates
	if got := run(`set.seed(1); s <- .Random.seed; s[2] <- 623L; .Random.seed <- s; invisible(runif(1)); .Random.seed[2]`); got != "[1] 624\n" {
		t.Fatalf("got %q", got)
	}
}

// Positions outside the generator's table would index out of range in the
// generator; R reads outside its state there. We refuse before the generator
// sees them, and a refused load leaves the session usable.
func TestRandomSeedRefusals(t *testing.T) {
	s := newRandSession()
	for _, p := range []string{"-1L", "-2147483647L", "NA_integer_"} {
		refused(t, s, fmt.Sprintf(`set.seed(1); s <- .Random.seed; s[2] <- %s; .Random.seed <- s; runif(1)`, p), "not a valid position")
		refused(t, s, fmt.Sprintf(`set.seed(1); s <- .Random.seed; s[2] <- %s; .Random.seed <- s; sample(5)`, p), "not a valid position")
		refused(t, s, fmt.Sprintf(`set.seed(1); s <- .Random.seed; s[2] <- %s; .Random.seed <- s; rnorm(1)`, p), "not a valid position")
	}
	refused(t, s, `set.seed(1); .Random.seed <- .Random.seed[1:625]; runif(1)`, "wrong length")
	refused(t, s, `set.seed(1); .Random.seed <- .Random.seed[1:2]; runif(1)`, "wrong length")
	refused(t, s, `set.seed(1); s <- .Random.seed; s[1] <- 10401L; .Random.seed <- s; runif(1)`, "not implemented")
	refused(t, s, `set.seed(1); s <- .Random.seed; s[1] <- 10407L; .Random.seed <- s; set.seed(1, kind = "Wichmann-Hill")`, "not implemented")
	refused(t, s, `RNGkind("Knuth-TAOCP")`, "not implemented")
	refused(t, s, `RNGkind(normal.kind = "Box-Muller"); rnorm(1)`, "not implemented")
	refused(t, s, `RNGkind(normal.kind = "Ahrens-Dieter"); rnorm(1)`, "not implemented")
	refused(t, s, `RNGversion("1.6.0")`, "not implemented")
	refused(t, s, `RNGkind(normal.kind = "foo")`, "not a valid choice")
	refused(t, s, `RNGkind(sample.kind = "foo")`, "not a valid choice")
	refused(t, s, `RNGkind("foo")`, "not a valid abbreviation")
	refused(t, s, `RNGkind(1)`, "must be a character string")
	refused(t, s, `RNGkind(normal.kind = 1)`, "must be a character string")
	refused(t, s, `RNGkind(sample.kind = 1)`, "must be a character string")
	refused(t, s, `RNGkind(normal.kind = c("Inversion", "Inversion"))`, "must be a character string")
	refused(t, s, `set.seed(NA)`, "supplied seed is not a valid integer")
	refused(t, s, `set.seed(numeric(0))`, "supplied seed is not a valid integer")
	refused(t, s, `set.seed(list(1))`, "supplied seed is not a valid integer")
	refused(t, s, `set.seed(NA_character_)`, "supplied seed is not a valid integer")
	refused(t, s, `set.seed(character(0))`, "supplied seed is not a valid integer")
	refused(t, s, `set.seed(NA_integer_)`, "supplied seed is not a valid integer")
	refused(t, s, `set.seed(NA)`, "supplied seed is not a valid integer")
	refused(t, s, `runif(-1)`, "invalid arguments")
	refused(t, s, `runif(NA)`, "invalid arguments")
	refused(t, s, `runif(list(1))`, "invalid arguments")
	refused(t, s, `runif(1, "a")`, "invalid arguments")
	refused(t, s, `rnorm(1, list(1))`, "invalid arguments")
	// after every refusal the session still draws
	res, err := s.Eval(`set.seed(1); length(runif(3))`)
	if err != nil || res.Output != "[1] 3\n" {
		t.Fatalf("%v %q", err, res.Output)
	}
}

// The warnings of a state that is ignored: the draw then randomizes.
func TestRandomSeedIgnored(t *testing.T) {
	s := newRandSession()
	for _, c := range []struct{ assign, warn string }{
		{`.Random.seed <- as.numeric(.Random.seed)`, "is not an integer vector"},
		{`.Random.seed <- integer(0)`, "is not a valid integer"},
		{`s <- .Random.seed; s[1] <- NA; .Random.seed <- s`, "is not a valid integer"},
		{`s <- .Random.seed; s[1] <- -3L; .Random.seed <- s`, "is not a valid integer"},
		{`s <- .Random.seed; s[1] <- 11001L; .Random.seed <- s`, "is not a valid integer"},
		{`s <- .Random.seed; s[1] <- 10408L; .Random.seed <- s`, "is not a valid RNG kind"},
		{`s <- .Random.seed; s[1] <- 10603L; .Random.seed <- s`, "is not a valid Normal type"},
	} {
		src := fmt.Sprintf(`set.seed(1); %s; tryCatch(runif(1), warning = conditionMessage)`, c.assign)
		res, err := s.Eval(src)
		if err != nil {
			t.Fatalf("%q: %v", src, err)
		}
		if !strings.Contains(res.Output, c.warn) {
			t.Errorf("%q: %q does not hold %q", src, res.Output, c.warn)
		}
	}
	// a value that is not a vector at all
	res, err := s.Eval(`set.seed(1); .Random.seed <- "a"; tryCatch(runif(1), warning = conditionMessage)`)
	if err != nil || !strings.Contains(res.Output, "is not an integer vector but of type 'character'") {
		t.Fatalf("%v %q", err, res.Output)
	}
}

func TestRandomKinds(t *testing.T) {
	s := newRandSession()
	run := func(src string) string {
		t.Helper()
		res, err := s.Eval(src)
		if err != nil {
			t.Fatalf("%q: %v", src, err)
		}
		return res.Output
	}
	if got := run(`set.seed(1); RNGkind()`); got != `[1] "Mersenne-Twister" "Inversion"        "Rejection"       `+"\n" {
		t.Fatalf("got %q", got)
	}
	if got := run(`suppressWarnings(RNGkind(sample.kind = "Rounding")); .Random.seed[1]`); got != "[1] 403\n" {
		t.Fatalf("got %q", got)
	}
	if got := run(`RNGkind(sample.kind = "Rejection"); .Random.seed[1]`); got != "[1] 10403\n" {
		t.Fatalf("got %q", got)
	}
	if got := run(`RNGkind(normal.kind = "Box-Muller"); .Random.seed[1]; RNGkind(normal.kind = "default"); .Random.seed[1]`); got != "[1] 10203\n[1] 10403\n" {
		t.Fatalf("got %q", got)
	}
	// abbreviations
	if got := run(`RNGkind("Mers", "Inv", "Rej"); .Random.seed[1]`); got != "[1] 10403\n" {
		t.Fatalf("got %q", got)
	}
	if got := run(`x <- RNGkind(character(0)); length(x)`); got != "[1] 3\n" {
		t.Fatalf("got %q", got)
	}
	// RNGversion maps versions to kinds
	if got := run(`suppressWarnings(RNGversion("3.5.0")); .Random.seed[1]; RNGversion("3.6.0"); .Random.seed[1]`); got != "[1] 403\n[1] 10403\n" {
		t.Fatalf("got %q", got)
	}
	if got := run(`RNGversion("4.0.0"); .Random.seed[1]`); got != "[1] 10403\n" {
		t.Fatalf("got %q", got)
	}
	// a new kind reseeds from the old stream (the words are in the golden)
	if got := run(`set.seed(7); invisible(RNGkind("Mersenne-Twister")); .Random.seed[1:5]`); got != "[1]       10403         624  -922521600 -1804554239  1194196430\n" {
		t.Fatalf("got %q", got)
	}
}

func TestSampleSemantics(t *testing.T) {
	s := newRandSession()
	run := func(src string) string {
		t.Helper()
		res, err := s.Eval(src)
		if err != nil {
			t.Fatalf("%q: %v", src, err)
		}
		return res.Output
	}
	if got := run(`set.seed(1); sample(10)`); got != " [1]  9  4  7  1  2  5  3 10  6  8\n" {
		t.Fatalf("got %q", got)
	}
	// sample(n) for a scalar n >= 1, sample(x) otherwise
	if got := run(`set.seed(1); sample(c(10))`); got != " [1]  9  4  7  1  2  5  3 10  6  8\n" {
		t.Fatalf("got %q", got)
	}
	// a set of positive probabilities picks only from the support
	if got := run(`set.seed(1); all(sample(5, 50, TRUE, prob = c(0, 1, 0, 0, 1)) %in% c(2, 5))`); got != "[1] TRUE\n" {
		t.Fatalf("got %q", got)
	}
	if got := run(`set.seed(1); sort(sample(6)); `); got != "[1] 1 2 3 4 5 6\n" {
		t.Fatalf("got %q", got)
	}
	// the result of a population above 2^31-1 is double
	if got := run(`set.seed(1); typeof(sample.int(2^31, 2)); typeof(sample.int(2^31 - 1, 2)); typeof(sample.int(10, 2))`); got != "[1] \"double\"\n[1] \"integer\"\n[1] \"integer\"\n" {
		t.Fatalf("got %q", got)
	}
	refused(t, s, `sample.int(5, 3, prob = c(1, 1, 1))`, "incorrect number of probabilities")
	refused(t, s, `sample.int(5, 3, prob = c(1, 1, 1, 1, NA))`, "NA in probability vector")
	refused(t, s, `sample.int(5, 3, prob = c(-1, 1, 1, 1, 1))`, "negative probability")
	refused(t, s, `sample.int(5, 3, prob = rep(0, 5))`, "too few positive probabilities")
	refused(t, s, `sample.int(5, 3, prob = c(1, 0, 0, 0, 0))`, "too few positive probabilities")
	refused(t, s, `sample.int(10, 20)`, "cannot take a sample larger")
	refused(t, s, `sample.int(10, -1)`, "invalid 'size' argument")
	refused(t, s, `sample.int(10, NA)`, "invalid 'size' argument")
	refused(t, s, `sample.int(10, "a")`, "invalid 'size' argument")
	refused(t, s, `sample.int(10, c(1, 2))`, "invalid 'size' argument")
	refused(t, s, `sample.int(10, list(1))`, "invalid 'size' argument")
	refused(t, s, `sample.int(10, 3, replace = NA)`, "invalid 'replace' argument")
	refused(t, s, `sample.int(10, 3, replace = "a")`, "invalid 'replace' argument")
	refused(t, s, `sample.int(-1)`, "invalid first argument")
	refused(t, s, `sample.int(5e15, 2)`, "invalid first argument")
	refused(t, s, `sample.int(0, 1)`, "invalid first argument")
	refused(t, s, `sample.int(2^40, 2^39 + 1, useHash = TRUE)`, "size <= n/2")
	refused(t, s, `sample.int(10, 3, TRUE, useHash = TRUE)`, "replace")
}

func FuzzRandomSeed(f *testing.F) {
	f.Add(int32(10403), int32(624), uint8(2), int32(1))
	f.Add(int32(10403), int32(-1), uint8(2), int32(1))
	f.Add(int32(10403), int32(625), uint8(2), int32(0))
	f.Add(int32(10403), int32(2147483647), uint8(2), int32(-1))
	f.Add(int32(10403), int32(-2147483648), uint8(2), int32(7))
	f.Add(int32(403), int32(1), uint8(2), int32(99))
	f.Add(int32(10303), int32(623), uint8(1), int32(5))
	f.Add(int32(-1), int32(0), uint8(0), int32(5))
	f.Add(int32(11001), int32(0), uint8(3), int32(5))
	f.Add(int32(10407), int32(300), uint8(2), int32(5))
	f.Fuzz(func(t *testing.T, head, pos int32, extra uint8, fill int32) {
		n := 622 + int(extra%4) // 622..625 words after the position, so 624 is the valid length
		s := eval.NewSession(eval.Options{MaxVector: 1 << 16, MaxAlloc: 1 << 26, Timeout: 2 * time.Second, MaxDepth: 200})
		prog := fmt.Sprintf(`.Random.seed <- c(%dL, %dL, rep(%dL, %d)); r <- c(runif(2), rnorm(2), sample(5), sample.int(100, 3, TRUE))`, head, pos, fill, n)
		s.Eval(prog)
		if _, err := s.Eval("NULL"); err != nil {
			t.Fatalf("session broken after %q: %v", prog, err)
		}
		// a state written by a draw is always a valid one
		res, err := s.Eval(`set.seed(1); s <- .Random.seed; .Random.seed <- s; length(runif(1))`)
		if err != nil || res.Output != "[1] 1\n" {
			t.Fatalf("%v %q", err, res.Output)
		}
	})
}
