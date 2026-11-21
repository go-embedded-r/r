package eval

import (
	"math"
	"math/big"
	"math/rand"
	"testing"
)

// The natives of slice 22, checked against exact arithmetic or a brute force.
//
//	choose      |choose(n, k) - C(n, k)| <= u |C(n, k)|, u = 2^-53, for integer n
//	            (|n| <= 400) and half-integer n with k <= 40, against math/big;
//	            for k > 1000 the result goes through lgamma and is held to 1e-9
//	findInterval  equal to a linear count of the knots at or below x
//	numToBits / numToInts  invert to the original bits
//	charmatch   the exact match wins, then the unique prefix, else 0

func bigBinomial(n, k int64) *big.Int { return new(big.Int).Binomial(n, k) }

func TestChooseAgainstBig(t *testing.T) {
	const u = 1.0 / (1 << 53)
	worst := 0.0
	for n := int64(0); n <= 400; n += 7 {
		for k := int64(0); k <= n && k <= 200; k += 3 {
			exact := new(big.Float).SetInt(bigBinomial(n, k))
			want, _ := exact.Float64()
			got := chooseOne(float64(n), float64(k))
			if math.IsInf(want, 0) {
				if !math.IsInf(got, 1) {
					t.Fatalf("choose(%d,%d) = %g, want +Inf", n, k, got)
				}
				continue
			}
			if rel := math.Abs(got-want) / want; rel > worst {
				worst = rel
			}
			if math.Abs(got-want) > u*want {
				t.Fatalf("choose(%d,%d) = %g, want %g", n, k, got, want)
			}
		}
	}
	// half-integer n: the product n (n-1) ... (n-k+1) / k! as an exact rational
	for _, n2 := range []int64{1, 3, 5, 9, 21, 41, 101} {
		for k := int64(0); k <= 40; k += 4 {
			num := big.NewRat(1, 1)
			for j := int64(0); j < k; j++ {
				num.Mul(num, big.NewRat(n2-2*j, 2))
			}
			num.Quo(num, new(big.Rat).SetInt(new(big.Int).MulRange(1, max(k, 1))))
			want, _ := num.Float64()
			got := chooseOne(float64(n2)/2, float64(k))
			if want == 0 {
				continue
			}
			if math.Abs(got-want) > u*math.Abs(want) {
				t.Fatalf("choose(%d/2,%d) = %g, want %g", n2, k, got, want)
			}
		}
	}
	// the reflection for a negative integer n: choose(-n, k) = (-1)^k choose(n+k-1, k)
	for n := int64(1); n < 30; n++ {
		for k := int64(0); k < 12; k++ {
			want := float64(bigBinomialFloat(n+k-1, k))
			if k%2 == 1 {
				want = -want
			}
			if got := chooseOne(float64(-n), float64(k)); got != want {
				t.Fatalf("choose(%d,%d) = %g, want %g", -n, k, got, want)
			}
		}
	}
	// beyond 1000 terms the gamma function is used
	exact, _ := new(big.Float).SetInt(bigBinomial(3000, 1500)).Float64()
	got := chooseOne(3000, 1500)
	if math.IsInf(exact, 0) {
		if !math.IsInf(got, 1) {
			t.Fatalf("choose(3000,1500) = %g, want +Inf", got)
		}
	} else if math.Abs(got-exact) > 1e-9*exact {
		t.Fatalf("choose(3000,1500) = %g, want %g", got, exact)
	}
	exact, _ = new(big.Float).SetInt(bigBinomial(1500, 700)).Float64()
	if got := chooseOne(1500, 700); math.Abs(got-exact) > 1e-9*exact && !math.IsInf(exact, 0) {
		t.Fatalf("choose(1500,700) = %g, want %g", got, exact)
	}
	t.Logf("worst relative error of choose over the grid: %.3g (bound 2^-53 = %.3g)", worst, u)
}

func bigBinomialFloat(n, k int64) float64 {
	f, _ := new(big.Float).SetInt(bigBinomial(n, k)).Float64()
	return f
}

func TestFindIntervalAgainstCount(t *testing.T) {
	r := rand.New(rand.NewSource(222))
	s := NewSession(Options{})
	for trial := 0; trial < 200; trial++ {
		nv := r.Intn(8)
		vec := make([]float64, nv)
		for i := range vec {
			vec[i] = float64(r.Intn(10))
		}
		for i := 1; i < nv; i++ {
			if vec[i] < vec[i-1] {
				vec[i] = vec[i-1]
			}
		}
		x := make([]float64, 12)
		for i := range x {
			x[i] = float64(r.Intn(24))/2 - 1
		}
		for _, lo := range []bool{false, true} {
			s.Set("x", &Double{V: x})
			s.Set("v", &Double{V: vec})
			linEval(t, s, "r <- findInterval(x, v, left.open = "+map[bool]string{true: "TRUE", false: "FALSE"}[lo]+")")
			got := linGet(t, s, "r")
			for i, xi := range x {
				want := 0
				for _, kv := range vec {
					if kv < xi || (!lo && kv == xi) {
						want++
					}
				}
				if int(got[i]) != want {
					t.Fatalf("findInterval(%g, %v, left.open=%v) = %v, want %d", xi, vec, lo, got[i], want)
				}
			}
		}
	}
}

func TestNumToBitsRoundTrip(t *testing.T) {
	s := NewSession(Options{})
	xs := []float64{0, 1, -2.5, math.Pi, math.Inf(-1), 5e-324, math.MaxFloat64, math.Copysign(0, -1)}
	s.Set("x", &Double{V: xs})
	linEval(t, s, "b <- as.integer(numToBits(x)); w <- numToInts(x)")
	b := linGet(t, s, "b")
	w := linGet(t, s, "w")
	for i, x := range xs {
		var u uint64
		for j := 63; j >= 0; j-- {
			u = u<<1 | uint64(b[64*i+j])
		}
		if u != math.Float64bits(x) {
			t.Fatalf("numToBits(%g) rebuilds %#x, want %#x", x, u, math.Float64bits(x))
		}
		if math.Float64bits(x)>>32 == 0x80000000 {
			continue // the high word is INT_MIN, which R prints as NA: not a number to rebuild
		}
		lo, hi := uint32(int32(w[2*i])), uint32(int32(w[2*i+1]))
		if uint64(hi)<<32|uint64(lo) != math.Float64bits(x) {
			t.Fatalf("numToInts(%g) rebuilds %#x", x, uint64(hi)<<32|uint64(lo))
		}
	}
}

func TestCharmatchRules(t *testing.T) {
	s := NewSession(Options{})
	cases := []struct{ src, want string }{
		{`charmatch("m", c("mean", "median"))`, "[1] 0\n"},
		{`charmatch("med", c("mean", "median"))`, "[1] 2\n"},
		{`charmatch("mean", c("mean", "meanx"))`, "[1] 1\n"},
		{`charmatch("a", c("a", "a"))`, "[1] 0\n"},
		{`charmatch("z", "a", nomatch = -1L)`, "[1] -1\n"},
		{`charmatch("", "")`, "[1] 1\n"},
	}
	for _, c := range cases {
		res, err := s.Eval(c.src)
		if err != nil {
			t.Fatalf("%s: %v", c.src, err)
		}
		if res.Output != c.want {
			t.Fatalf("%s: got %q, want %q", c.src, res.Output, c.want)
		}
	}
}

func TestEncodeStringGuards(t *testing.T) {
	s := NewSession(Options{})
	for _, src := range []string{
		`encodeString("a", width = -1)`,
		`encodeString("a", width = c(1, 2))`,
		`encodeString("a", quote = c("'", "\""))`,
		`encodeString("a", justify = "diagonal")`,
		`encodeString(list(1))`,
	} {
		if _, err := s.Eval(src); err == nil {
			t.Fatalf("%s: no error", src)
		}
	}
}

func TestChooseExactEdges(t *testing.T) {
	if got := chooseExact(3, 5); got != 0 {
		t.Fatalf("chooseExact(3, 5) = %g, want 0", got)
	}
	if got := chooseExact(7, 0); got != 1 {
		t.Fatalf("chooseExact(7, 0) = %g, want 1", got)
	}
	if got := chooseExact(1e300, 1000); !math.IsInf(got, 1) {
		t.Fatalf("chooseExact(1e300, 1000) = %g, want +Inf", got)
	}
}

func TestChooseEmptyArguments(t *testing.T) {
	s := NewSession(Options{})
	for src, want := range map[string]string{
		"choose(numeric(), 2)": "numeric(0)\n",
		"choose(3, integer())": "numeric(0)\n",
		"choose(4, 4)":         "[1] 1\n",
		"choose(5, 5)":         "[1] 1\n",
	} {
		res, err := s.Eval(src)
		if err != nil {
			t.Fatalf("%s: %v", src, err)
		}
		if res.Output != want {
			t.Fatalf("%s: got %q, want %q", src, res.Output, want)
		}
	}
}
