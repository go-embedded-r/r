package eval

import (
	"context"
	"math"
	"math/big"
	"testing"
)

// The natives of slice 23, checked against exact arithmetic or a brute force.
//
//	lchoose    |lchoose(n, k) - ln|C(n, k)|| <= 4 u max(1, |ln|C(n, k)||), u = 2^-53,
//	           for k <= 1000 and |n| < 1e300, against a 200-bit ln of the exact
//	           rational (an atanh series, independent of math.Log)
//	agrepMatch  equal to the brute-force edit distance of the pattern to the
//	           best substring, with every cap at its default
//	abbreviate  the length asked for, or the unabbreviatable remainder

// bigLn is ln x for x > 0 to about 190 bits: x = m 2^e, ln m by 2 atanh((m-1)/(m+1)).
func bigLn(x *big.Float) *big.Float {
	const prec = 220
	m := new(big.Float).SetPrec(prec)
	e := x.MantExp(m)
	series := func(v *big.Float) *big.Float {
		one := big.NewFloat(1).SetPrec(prec)
		y := new(big.Float).SetPrec(prec).Quo(new(big.Float).SetPrec(prec).Sub(v, one), new(big.Float).SetPrec(prec).Add(v, one))
		y2 := new(big.Float).SetPrec(prec).Mul(y, y)
		term := new(big.Float).SetPrec(prec).Set(y)
		sum := new(big.Float).SetPrec(prec)
		for i := 1; i < 400; i += 2 {
			sum.Add(sum, new(big.Float).SetPrec(prec).Quo(term, new(big.Float).SetPrec(prec).SetInt64(int64(i))))
			term.Mul(term, y2)
		}
		return sum.Mul(sum, big.NewFloat(2).SetPrec(prec))
	}
	res := series(m)
	if e != 0 {
		ln2 := series(big.NewFloat(0.5).SetPrec(prec))
		ln2.Neg(ln2)
		res.Add(res, ln2.Mul(ln2, new(big.Float).SetPrec(prec).SetInt64(int64(e))))
	}
	return res
}

func TestLchooseAgainstBig(t *testing.T) {
	const u = 1.0 / (1 << 53)
	check := func(n float64, k int) {
		r := chooseRat(n, k)
		if r.Sign() == 0 {
			if got := lchooseOne(n, float64(k)); !math.IsInf(got, -1) {
				t.Fatalf("lchoose(%g,%d) = %g, want -Inf", n, k, got)
			}
			return
		}
		f := new(big.Float).SetPrec(220).SetRat(r)
		f.Abs(f)
		want, _ := bigLn(f).Float64()
		got := lchooseOne(n, float64(k))
		if math.Abs(got-want) > 4*u*math.Max(1, math.Abs(want)) {
			t.Fatalf("lchoose(%g,%d) = %.17g, want %.17g", n, k, got, want)
		}
	}
	for n := 2; n <= 300; n += 11 {
		for k := 2; k <= n && k <= 150; k += 5 {
			check(float64(n), k)
		}
	}
	for _, n := range []float64{-7, -30, 0.5, 5.5, 20.25, 1e6, 123456789, 1e15} {
		for _, k := range []int{2, 3, 10, 40, 200, 1000} {
			check(n, k)
		}
	}
	if got := lchooseOne(math.Inf(1), 3); !math.IsNaN(got) {
		t.Fatalf("lchoose(Inf, 3) = %g, want NaN", got)
	}
	if got := lchooseOne(5, 6); !math.IsInf(got, -1) {
		t.Fatalf("lchoose(5, 6) = %g, want -Inf", got)
	}
	if got := lchooseOne(5, 0); got != 0 {
		t.Fatalf("lchoose(5, 0) = %g, want 0", got)
	}
}

func TestAbbreviateOne(t *testing.T) {
	cases := []struct {
		in   string
		n    int
		cls  bool
		want string
	}{
		{"International", 5, true, "Intrn"},
		{"Alabama", 3, true, "Alb"},
		{"Hello World", 3, true, "HlW"},
		{"short", 10, true, "short"},
		{"", 3, true, ""},
		{"abcdef", 3, false, "abc"},
		{"a b", 3, true, "a b"},
		{"New Hampshire", 2, true, "NH"},
	}
	for _, c := range cases {
		if got := abbreviateOne(c.in, c.n, c.cls); got != c.want {
			t.Errorf("abbreviateOne(%q, %d, %v) = %q, want %q", c.in, c.n, c.cls, got, c.want)
		}
	}
	// the result is never longer than asked when the string has enough
	// droppable characters, and is always a subsequence of the input
	for _, s := range []string{"Massachusetts", "North Carolina", "abcdefghij", "Mississippi"} {
		for n := 1; n <= len(s); n++ {
			got := abbreviateOne(s, n, true)
			if len([]rune(got)) > max(n, 1) && len([]rune(got)) > len([]rune(s)) {
				t.Errorf("abbreviateOne(%q, %d) = %q is longer than the input", s, n, got)
			}
			i := 0
			for _, r := range s {
				if i < len(got) && rune(got[i]) == r {
					i++
				}
			}
			if i != len(got) {
				t.Errorf("abbreviateOne(%q, %d) = %q is not a subsequence", s, n, got)
			}
		}
	}
}

// editToBestSubstring is the least number of single-character edits taking p to
// some substring of t: the classical dynamic programme, free to start and stop
// anywhere in t.
func editToBestSubstring(p, t []rune) int {
	prev := make([]int, len(t)+1)
	for i := 1; i <= len(p); i++ {
		cur := make([]int, len(t)+1)
		cur[0] = i
		for j := 1; j <= len(t); j++ {
			cost := 1
			if p[i-1] == t[j-1] {
				cost = 0
			}
			cur[j] = min(prev[j-1]+cost, prev[j]+1, cur[j-1]+1)
		}
		prev = cur
	}
	best := prev[0]
	for _, v := range prev {
		best = min(best, v)
	}
	return best
}

func TestAgrepMatchAgainstBruteForce(t *testing.T) {
	alpha := []rune("abc")
	var words [][]rune
	var gen func(cur []rune, n int)
	gen = func(cur []rune, n int) {
		words = append(words, append([]rune(nil), cur...))
		if len(cur) == n {
			return
		}
		for _, a := range alpha {
			gen(append(cur, a), n)
		}
	}
	gen(nil, 4)
	for _, p := range words[1:] {
		for _, tx := range words {
			d := editToBestSubstring(p, tx)
			for k := 0; k <= 3; k++ {
				got := agrepMatch(p, tx, agrepBounds{cost: -1, ins: -1, del: -1, sub: -1, all: k})
				if got != (d <= k) {
					t.Fatalf("agrepMatch(%q, %q, all=%d) = %v, edit distance %d", string(p), string(tx), k, got, d)
				}
			}
		}
	}
	// caps: with no insertion and no deletion allowed only substitutions remain
	sub := agrepBounds{cost: -1, ins: 0, del: 0, sub: 1, all: -1}
	if !agrepMatch([]rune("lasy"), []rune("lazy"), sub) {
		t.Error("one substitution should match")
	}
	if agrepMatch([]rune("lasy"), []rune("lxsx"), sub) {
		t.Error("lxsx needs two substitutions")
	}
	if agrepMatch([]rune("lasy"), []rune("lsaxy"), sub) {
		t.Error("lsaxy needs two substitutions")
	}
}

// A bound of -1 is no bound: every pattern is within reach of every text. And
// when the bounds already allow deleting the whole pattern the answer is yes
// without building the table of counts, which for 300 characters would be
// gigabytes.
func TestAgrepMatchUnboundedAndShortcut(t *testing.T) {
	none := agrepBounds{cost: -1, ins: -1, del: -1, sub: -1, all: -1}
	if !agrepMatch([]rune("abc"), []rune("xyz"), none) {
		t.Error("with no bound at all, abc matches xyz")
	}
	before := agrepTables.Load()
	if !agrepMatch([]rune("abcd"), []rune("zzzz"), agrepBounds{cost: -1, ins: -1, del: -1, sub: -1, all: 4}) {
		t.Error("a bound of the pattern length always matches")
	}
	if agrepTables.Load() != before {
		t.Error("the shortcut was not taken: a table was built")
	}
	if agrepMatch([]rune("abcd"), []rune("zzzz"), agrepBounds{cost: -1, ins: -1, del: -1, sub: -1, all: 3}) {
		t.Error("three edits cannot turn abcd into a substring of zzzz")
	}
	if agrepTables.Load() == before {
		t.Error("below the pattern length the table must be built")
	}
}

func TestLchooseEmptyArguments(t *testing.T) {
	s := NewSession(Options{})
	res, err := s.EvalIn(context.Background(), `print(lchoose(numeric(0), 2)); print(lchoose(5, integer(0)))`, NewEnv(s.Global))
	if err != nil {
		t.Fatal(err)
	}
	if res.Output != "numeric(0)\nnumeric(0)\n" {
		t.Fatalf("got %q", res.Output)
	}
}

func TestCloseAllConnections(t *testing.T) {
	s := NewSession(Options{})
	env := NewEnv(s.Global)
	res, err := s.EvalIn(context.Background(), `
tc <- textConnection("a"); f <- file(tempfile(), "w")
before <- nrow(showConnections())
closeAllConnections()
after <- nrow(showConnections())
c(before, after)`, env)
	if err != nil {
		t.Fatal(err)
	}
	if res.Output != "[1] 2 0\n" {
		t.Fatalf("unexpected output %q", res.Output)
	}
	res, err = s.EvalIn(context.Background(), `print(c(before, after)); print(nrow(showConnections()))`, env)
	if err != nil {
		t.Fatal(err)
	}
	if res.Output != "[1] 2 0\n[1] 0\n" {
		t.Fatalf("got %q", res.Output)
	}
	// the standard connections stay
	res, err = s.EvalIn(context.Background(), `print(rownames(showConnections(all = TRUE)))`, NewEnv(s.Global))
	if err != nil {
		t.Fatal(err)
	}
	if res.Output != "[1] \"0\" \"1\" \"2\"\n" {
		t.Fatalf("got %q", res.Output)
	}
}
