package eval

import (
	"math"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-ndarray/ndarray/random"
	"github.com/go-statistics/stats"
)

// Slice 24 (M4): random streams. Written from the help pages (?Random,
// ?RNGkind, ?set.seed, ?sample, ?Normal, ?jitter), from Matsumoto and
// Nishimura (1998) for the generator, and from experiments with the oracle
// for everything else (the layout of .Random.seed, the seeding, how many
// 32-bit words each function draws). The generator itself is go-ndarray's
// random.MT19937.

const (
	mtWords   = 624
	rngI2_32m = 2.328306437080797e-10 // 2^-32 (R's i2_32m1)
)

var rngKindNames = []string{"Wichmann-Hill", "Marsaglia-Multicarry", "Super-Duper", "Mersenne-Twister",
	"Knuth-TAOCP", "user-supplied", "Knuth-TAOCP-2002", "L'Ecuyer-CMRG"}
var rngNormalNames = []string{"Buggy Kinderman-Ramage", "Ahrens-Dieter", "Box-Muller", "user-supplied", "Inversion", "Kinderman-Ramage"}
var rngSampleNames = []string{"Rounding", "Rejection"}

var stdNormal = stats.Normal{Loc: 0, Scale: 1}

var rngClock atomic.Uint64

// rngKinds returns the Normal and sample kinds in force (default Inversion,
// Rejection). rngCode is 0 until a state was loaded or written.
func (s *Session) rngKinds() (norm, samp int) {
	if s.rngCode == 0 {
		return 4, 1
	}
	c := s.rngCode - 1
	return (c % 10000) / 100, c / 10000
}

func (s *Session) setRngKinds(norm, samp int) { s.rngCode = 1 + 100*norm + 10000*samp }

// mtScramble is the seeding of set.seed: the seed goes through fifty steps of
// the congruential generator 69069*s+1, then one more step per element of
// .Random.seed; the first element is the position, the rest are the words.
func mtScramble(seed uint32) (key [mtWords]uint32) {
	for j := 0; j < 50; j++ {
		seed = 69069*seed + 1
	}
	for j := 0; j <= mtWords; j++ {
		seed = 69069*seed + 1
		if j > 0 {
			key[j-1] = seed
		}
	}
	return key
}

// mtSgenrand is the default initialisation of Matsumoto and Nishimura's
// reference code (sgenrand), used when the position says "not initialised".
func mtSgenrand(seed uint32) (key [mtWords]uint32) {
	for i := 0; i < mtWords; i++ {
		key[i] = seed & 0xffff0000
		seed = 69069*seed + 1
		key[i] |= (seed & 0xffff0000) >> 16
		seed = 69069*seed + 1
	}
	return key
}

func rngSeed(g *random.MT19937, seed uint32) {
	g.SetState(random.MT19937State{Key: mtScramble(seed), Pos: mtWords})
}

func rngClockSeed() uint32 {
	t := uint64(time.Now().UnixNano())
	t ^= rngClock.Add(0x9e3779b97f4a7c15) * 0xff51afd7ed558ccd
	return uint32(t ^ t>>32)
}

var randomSeedSym = Intern(".Random.seed")

// rngLoad reads .Random.seed as GetRNGstate does. It returns false when the
// value is ignored (a warning was raised): the caller then randomizes.
func (s *Session) rngLoad(g *random.MT19937, v Value, call *Lang) bool {
	iv, ok := v.(*Integer)
	if !ok {
		s.warningf(call, "'.Random.seed' is not an integer vector but of type '%s', so ignored", typeName(v))
		return false
	}
	if len(iv.V) == 0 {
		s.warningf(call, "'.Random.seed[1]' is not a valid integer, so ignored")
		return false
	}
	tmp := iv.V[0]
	if tmp == naInt || tmp < 0 || tmp > 11000 {
		s.warningf(call, "'.Random.seed[1]' is not a valid integer, so ignored")
		return false
	}
	if tmp%100 > 7 {
		s.warningf(call, "'.Random.seed[1]' is not a valid RNG kind so ignored")
		return false
	}
	if (tmp%10000)/100 > 5 {
		s.warningf(call, "'.Random.seed[1]' is not a valid Normal type, so ignored")
		return false
	}
	if tmp%100 != 3 {
		s.errorf(call, "RNG kind '%s' is not implemented: only Mersenne-Twister is available", rngKindNames[tmp%100])
	}
	if len(iv.V) < mtWords+2 {
		s.errorf(call, "'.Random.seed' has wrong length")
	}
	pos := iv.V[1]
	if pos == naInt || pos < 0 {
		// R reads outside its state here; we refuse before the generator sees it
		s.errorf(call, "'.Random.seed[2]' is not a valid position")
	}
	var st random.MT19937State
	for i := range st.Key {
		st.Key[i] = uint32(iv.V[2+i])
	}
	switch {
	case pos >= 1 && pos <= mtWords-1:
		st.Pos = int(pos)
	case pos == mtWords+1:
		st.Key = mtSgenrand(4357)
		st.Pos = mtWords
	default: // 0, 624 and anything beyond 625 regenerate
		st.Pos = mtWords
	}
	g.SetState(st)
	s.setRngKinds(int((tmp%10000)/100), int(tmp/10000))
	return true
}

// rngBegin is GetRNGstate: it loads .Random.seed from the global environment,
// or randomizes when there is none or it is ignored.
func (s *Session) rngBegin() *random.MT19937 {
	call := s.currentCall()
	g := new(random.MT19937)
	if v, found := s.getVar(".Random.seed", s.Global, false, false); found {
		if s.rngLoad(g, v, call) {
			return g
		}
	}
	rngSeed(g, rngClockSeed())
	return g
}

// rngEnd is PutRNGstate: a fresh plain integer vector of 626 in the global
// environment.
func (s *Session) rngEnd(g *random.MT19937) {
	st := g.State()
	out := make([]int32, mtWords+2)
	norm, samp := s.rngKinds()
	out[0] = int32(3 + 100*norm + 10000*samp)
	out[1] = int32(st.Pos)
	for i, w := range st.Key {
		out[2+i] = int32(w)
	}
	s.assignTo(s.Global, randomSeedSym, &Integer{V: out}, s.currentCall())
}

// rngPeek reads the kinds from a valid .Random.seed without warnings.
func (s *Session) rngPeek() {
	v, _ := s.getVar(".Random.seed", s.Global, false, false) // absent: v is nil, not an *Integer
	iv, ok := v.(*Integer)
	if !ok || len(iv.V) == 0 {
		return
	}
	t := iv.V[0]
	if t == naInt || t < 0 || t > 11000 || t%100 > 7 || (t%10000)/100 > 5 {
		return
	}
	s.setRngKinds(int((t%10000)/100), int(t/10000))
}

func unifRand(g *random.MT19937) float64 {
	v := float64(g.Uint32()) * 2.3283064365386963e-10
	if v <= 0 {
		return 0.5 * rngI2_32m
	}
	// the largest word gives 1 - 2^-32 here, so there is no upper fixup
	return v
}

// normRand is the Inversion method: two uniforms make a 27-bit-finer one.
func normRand(g *random.MT19937) float64 {
	const big = 134217728.0
	u := unifRand(g)
	u = math.Floor(big*u) + unifRand(g)
	return stdNormal.PPF(u / big)
}

// unifIndex draws an integer in [0, dn): by rejection on masked bits, or the
// old rounding of dn*u.
func unifIndex(g *random.MT19937, dn float64, rounding bool) float64 {
	if rounding {
		return math.Floor(dn * unifRand(g))
	}
	// the callers draw from a population of at least one
	bits := int(math.Ceil(math.Log2(dn)))
	// each attempt succeeds with probability above one half, so the bound is
	// never reached by chance (2^-4096); it only keeps a broken test finite
	for range 4096 {
		var v int64
		for n := 0; n <= bits; n += 16 {
			v = 65536*v + int64(math.Floor(unifRand(g)*65536))
		}
		if bits < 64 {
			v &= int64(1)<<uint(bits) - 1
		}
		if dv := float64(v); dn > dv {
			return dv
		}
	}
	panic("unifIndex: rejection sampling did not terminate")
}

// rngInvalid raises the error of the random generators for a bad argument.
func (s *Session) rngInvalid(call *Lang) { s.errorf(call, "invalid arguments") }

// rngCount reads the n of runif and rnorm: a vector of length > 1 stands for
// its length.
func (s *Session) rngCount(v Value, call *Lang) int {
	var d float64
	switch x := v.(type) {
	case *Logical, *Integer, *Double:
		l := Length(x)
		if l != 1 {
			return l
		}
		d = asDoubles(x)[0]
	case *String:
		if len(x.V) != 1 {
			return len(x.V)
		}
		f, ok := parseNumber(x.V[0])
		if !ok {
			s.warningf(call, "NAs introduced by coercion")
			f = rNaN
		}
		d = f
	default:
		s.rngInvalid(call)
	}
	if math.IsNaN(d) || d < 0 || d > 4.5e15 {
		s.rngInvalid(call)
	}
	n := int(d)
	s.checkLen(call, n)
	return n
}

func (s *Session) rngNumArg(v Value, call *Lang) []float64 {
	switch v.(type) {
	case *Logical, *Integer, *Double:
		return asDoubles(v)
	}
	s.rngInvalid(call)
	return nil
}

func init() {
	regBuiltin("..runif", func(s *Session, c *callInfo) Value {
		call := s.currentCall()
		n := s.rngCount(c.args[0], call)
		mn, mx := s.rngNumArg(c.args[1], call), s.rngNumArg(c.args[2], call)
		out := make([]float64, n)
		if n == 0 {
			return &Double{V: out}
		}
		if len(mn) == 0 || len(mx) == 0 {
			for i := range out {
				out[i] = naReal
			}
			s.warningf(call, "NAs produced")
			return &Double{V: out}
		}
		g := s.rngBegin()
		bad := false
		for i := range out {
			if i&255 == 0 {
				s.checkInterrupt()
			}
			lo, hi := mn[i%len(mn)], mx[i%len(mx)]
			switch {
			case math.IsNaN(lo) || math.IsNaN(hi) || math.IsInf(lo, 0) || math.IsInf(hi, 0) || hi < lo:
				out[i], bad = rNaN, true
			case lo == hi:
				out[i] = lo
			default:
				out[i] = lo + (hi-lo)*unifRand(g)
			}
		}
		s.rngEnd(g)
		if bad {
			s.warningf(call, "NAs produced")
		}
		return &Double{V: out}
	})
	regBuiltin("..rnorm", func(s *Session, c *callInfo) Value {
		call := s.currentCall()
		n := s.rngCount(c.args[0], call)
		mu, sd := s.rngNumArg(c.args[1], call), s.rngNumArg(c.args[2], call)
		out := make([]float64, n)
		if n == 0 {
			return &Double{V: out}
		}
		if len(mu) == 0 || len(sd) == 0 {
			for i := range out {
				out[i] = naReal
			}
			s.warningf(call, "NAs produced")
			return &Double{V: out}
		}
		g := s.rngBegin()
		bad := false
		for i := range out {
			if i&255 == 0 {
				s.checkInterrupt()
			}
			m, sg := mu[i%len(mu)], sd[i%len(sd)]
			switch {
			case math.IsNaN(m) || math.IsInf(sg, 0) || math.IsNaN(sg) || sg < 0:
				out[i], bad = rNaN, true
			case sg == 0 || math.IsInf(m, 0):
				out[i] = m
			default:
				if nk, _ := s.rngKinds(); nk != 4 {
					s.errorf(call, "normal.kind '%s' is not implemented: only Inversion is available", rngNormalNames[nk])
				}
				out[i] = m + sg*normRand(g)
			}
		}
		s.rngEnd(g)
		if bad {
			s.warningf(call, "NAs produced")
		}
		return &Double{V: out}
	})

	regBuiltin("..set.seed", func(s *Session, c *callInfo) Value {
		call := s.currentCall()
		kind, nk, sk := s.rngKindArgs(call, c.args[1], c.args[2], c.args[3])
		randomize := false
		var seed int32
		if _, isNull := c.args[0].(nullT); isNull {
			randomize = true
		} else {
			seed = s.rngSeedArg(call, c.args[0])
		}
		s.rngApplyKinds(call, kind, nk, sk)
		g := new(random.MT19937)
		if randomize {
			rngSeed(g, rngClockSeed())
		} else {
			rngSeed(g, uint32(seed))
		}
		s.rngEnd(g)
		return Null
	})
	regBuiltin("..RNGkind", func(s *Session, c *callInfo) Value {
		call := s.currentCall()
		kind, nk, sk := s.rngKindArgs(call, c.args[0], c.args[1], c.args[2])
		s.rngPeek()
		norm, samp := s.rngKinds()
		old := &String{V: []Str{{S: "Mersenne-Twister"}, {S: rngNormalNames[norm]}, {S: rngSampleNames[samp]}}}
		if kind < 0 && nk < 0 && sk < 0 {
			return old
		}
		g := s.rngBegin()
		if kind >= 0 {
			// a new kind reseeds the generator from the old stream
			u := unifRand(g)
			rngSeed(g, uint32(u*4294967295.0))
		}
		s.rngApplyKinds(call, kind, nk, sk)
		s.rngEnd(g)
		return old
	})
}

// rngSeedArg is the seed of set.seed as an integer, or an error.
func (s *Session) rngSeedArg(call *Lang, v Value) int32 {
	bad := func() { s.errorf(call, "supplied seed is not a valid integer") }
	switch x := v.(type) {
	case *Logical:
		if len(x.V) == 0 || x.V[0] == naInt {
			bad()
		}
		return x.V[0]
	case *Integer:
		if len(x.V) == 0 || x.V[0] == naInt {
			bad()
		}
		return x.V[0]
	case *Double:
		if len(x.V) == 0 || math.IsNaN(x.V[0]) {
			bad()
		}
		t := math.Trunc(x.V[0])
		if math.Abs(t) >= 2147483648 {
			s.warningf(call, "NAs introduced by coercion to integer range")
			bad()
		}
		return int32(t)
	case *String:
		if len(x.V) == 0 {
			bad()
		}
		f, ok := parseNumber(x.V[0])
		if !ok {
			s.warningf(call, "NAs introduced by coercion")
			bad()
		}
		if math.IsNaN(f) {
			bad()
		}
		t := math.Trunc(f)
		if math.Abs(t) >= 2147483648 {
			s.warningf(call, "NAs introduced by coercion to integer range")
			bad()
		}
		return int32(t)
	}
	bad()
	return 0
}

// rngPmatch is pmatch for one name: exact, else a unique prefix; the last
// index (len(names)) is "default".
func rngPmatch(x string, names []string) int {
	all := append(append([]string(nil), names...), "default")
	for i, n := range all {
		if n == x {
			return i
		}
	}
	hit := -1
	for i, n := range all {
		if x != "" && strings.HasPrefix(n, x) {
			if hit >= 0 {
				return -1
			}
			hit = i
		}
	}
	return hit
}

// rngKindArgs resolves the kind, normal.kind and sample.kind arguments; -1
// means NULL (unchanged), "default" is resolved to the default kind.
func (s *Session) rngKindArgs(call *Lang, kind, norm, samp Value) (int, int, int) {
	str := func(v Value, what string) (string, bool) {
		if _, isNull := v.(nullT); isNull {
			return "", false
		}
		fail := func() {
			if what == "kind" {
				s.errorf(call, "'kind' must be a character string (RNG to be used).")
			}
			s.errorf(call, "'%s' must be a character string", what)
		}
		x, ok := v.(*String)
		if !ok {
			fail()
		}
		if len(x.V) == 0 && what == "kind" {
			return "", false
		}
		if len(x.V) != 1 {
			fail()
		}
		if x.V[0].NA {
			return "NA", true
		}
		return x.V[0].S, true
	}
	k, nk, sk := -1, -1, -1
	if name, ok := str(kind, "kind"); ok {
		i := rngPmatch(name, rngKindNames)
		if i < 0 {
			s.errorf(call, "'%s' is not a valid abbreviation of an RNG", name)
		}
		if i == len(rngKindNames) {
			i = 3
		}
		k = i
	}
	if name, ok := str(norm, "normal.kind"); ok {
		i := rngPmatch(name, rngNormalNames)
		if i < 0 {
			s.errorf(call, "'NA' is not a valid choice")
		}
		if i == len(rngNormalNames) {
			i = 4
		}
		nk = i
	}
	if name, ok := str(samp, "sample.kind"); ok {
		i := rngPmatch(name, rngSampleNames)
		if i < 0 {
			s.errorf(call, "'NA' is not a valid choice")
		}
		if i == len(rngSampleNames) {
			i = 1
		}
		sk = i
	}
	return k, nk, sk
}

// rngApplyKinds installs the resolved kinds: only Mersenne-Twister exists
// here, the Normal kind is stored (rnorm refuses all but Inversion).
func (s *Session) rngApplyKinds(call *Lang, kind, nk, sk int) {
	if kind >= 0 && kind != 3 {
		s.errorf(call, "RNG kind '%s' is not implemented: only Mersenne-Twister is available", rngKindNames[kind])
	}
	norm, samp := s.rngKinds()
	if nk >= 0 {
		norm = nk
	}
	if sk >= 0 {
		samp = sk
		if sk == 0 {
			s.warningf(call, "non-uniform 'Rounding' sampler used")
		}
	}
	s.setRngKinds(norm, samp)
}

// ---- sample ----------------------------------------------------------------

// revsort sorts a into descending order by heapsort, permuting ib with it
// (the order of ties is the heapsort's, which sample's output depends on).
func revsort(a []float64, ib []int) {
	n := len(a)
	if n <= 1 {
		return
	}
	// 1-based, as in Numerical Recipes
	l := n>>1 + 1
	ir := n
	at := func(i int) *float64 { return &a[i-1] }
	bt := func(i int) *int { return &ib[i-1] }
	for {
		var ra float64
		var ii int
		if l > 1 {
			l--
			ra, ii = *at(l), *bt(l)
		} else {
			ra, ii = *at(ir), *bt(ir)
			*at(ir), *bt(ir) = *at(1), *bt(1)
			ir--
			if ir == 1 {
				*at(1), *bt(1) = ra, ii
				return
			}
		}
		i, j := l, l<<1
		for j <= ir {
			if j < ir && *at(j) > *at(j + 1) {
				j++
			}
			if ra > *at(j) {
				*at(i), *bt(i) = *at(j), *bt(j)
				i = j
				j += i
			} else {
				j = ir + 1
			}
		}
		*at(i), *bt(i) = ra, ii
	}
}

func (s *Session) sampleSize(v Value, call *Lang) int {
	bad := func() { s.errorf(call, "invalid 'size' argument") }
	var d float64
	switch x := v.(type) {
	case *Logical, *Integer, *Double:
		if Length(x) != 1 {
			bad()
		}
		d = asDoubles(x)[0]
	case *String:
		if len(x.V) != 1 {
			bad()
		}
		f, ok := parseNumber(x.V[0])
		if !ok {
			bad()
		}
		d = f
	default:
		bad()
	}
	if math.IsNaN(d) || d < 0 || d > 4.5e15 {
		bad()
	}
	return int(d)
}

func (s *Session) sampleN(v Value, call *Lang) float64 {
	switch v.(type) {
	case *Logical, *Integer, *Double:
		if d := asDoubles(v); len(d) >= 1 {
			return d[0]
		}
	}
	s.errorf(call, "invalid first argument")
	return 0
}

func (s *Session) sampleReplace(v Value, call *Lang) bool {
	bad := func() { s.errorf(call, "invalid 'replace' argument") }
	switch x := v.(type) {
	case *Logical:
		if len(x.V) == 1 && x.V[0] != naInt {
			return x.V[0] != 0
		}
	case *Integer:
		if len(x.V) == 1 && x.V[0] != naInt {
			return x.V[0] != 0
		}
	case *Double:
		if len(x.V) == 1 && !math.IsNaN(x.V[0]) {
			return x.V[0] != 0
		}
	case *String:
		if len(x.V) == 1 && !x.V[0].NA {
			switch x.V[0].S {
			case "TRUE", "true", "T", "True":
				return true
			case "FALSE", "false", "F", "False":
				return false
			}
		}
	}
	bad()
	return false
}

// sampleProb reads and normalises prob (length n): it must be finite and
// non-negative with enough positive entries.
func (s *Session) sampleProb(v Value, n int, k int, replace bool, call *Lang) []float64 {
	var p []float64
	switch x := v.(type) {
	case *Logical, *Integer, *Double:
		p = append([]float64(nil), asDoubles(x)...)
	case *String:
		p = make([]float64, len(x.V))
		warned := false
		for i, e := range x.V {
			f, ok := parseNumber(e)
			if !ok {
				if !warned {
					s.warningf(call, "NAs introduced by coercion")
					warned = true
				}
				f = naReal
			}
			p[i] = f
		}
	case *List:
		p = make([]float64, len(x.V))
		for i, e := range x.V {
			d := asDoubles(e)
			if len(d) != 1 {
				s.errorf(call, "'list' object cannot be coerced to type 'double'")
			}
			p[i] = d[0]
		}
	default:
		s.errorf(call, "invalid 'prob' argument")
	}
	if len(p) != n {
		s.errorf(call, "incorrect number of probabilities")
	}
	npos := 0
	sum := 0.0
	for _, e := range p {
		if math.IsNaN(e) || math.IsInf(e, 0) {
			s.errorf(call, "NA in probability vector")
		}
		if e < 0 {
			s.errorf(call, "negative probability")
		}
		if e > 0 {
			npos++
			sum += e
		}
	}
	if npos == 0 || (!replace && k > npos) {
		s.errorf(call, "too few positive probabilities")
	}
	for i := range p {
		p[i] /= sum
	}
	return p
}

func sampleResult(vals []float64, dn float64) Value {
	if dn <= math.MaxInt32 {
		out := make([]int32, len(vals))
		for i, v := range vals {
			out[i] = int32(v)
		}
		return &Integer{V: out}
	}
	return &Double{V: vals}
}

func init() {
	regBuiltin("..sample", func(s *Session, c *callInfo) Value {
		call := s.currentCall()
		dn := s.sampleN(c.args[0], call)
		if math.IsNaN(dn) || dn < 0 || dn > 4.5e15 {
			s.errorf(call, "invalid first argument")
		}
		k := s.sampleSize(c.args[1], call)
		replace := s.sampleReplace(c.args[2], call)
		if math.IsNaN(dn) || dn < 0 || dn > 4.5e15 || (k > 0 && dn == 0) {
			s.errorf(call, "invalid first argument")
		}
		s.checkLen(call, k)
		_, isNull := c.args[3].(nullT)
		var p []float64
		n := int(dn)
		if !isNull {
			s.checkLen(call, n)
			p = s.sampleProb(c.args[3], n, k, replace, call)
		} else if !replace && k > n {
			s.errorf(call, "cannot take a sample larger than the population when 'replace = FALSE'")
		}
		if isNull && !replace && k >= 2 {
			s.checkLen(call, n)
		}
		g := s.rngBegin()
		_, samp := s.rngKinds()
		rounding := samp == 0
		out := make([]float64, k)
		switch {
		case !isNull && replace:
			nc := 0
			for _, e := range p {
				if float64(n)*e > 0.1 {
					nc++
				}
			}
			if nc > 200 {
				s.walkerSample(g, p, out, rounding)
			} else {
				s.probSampleReplace(g, p, out)
			}
		case !isNull:
			s.probSampleNoReplace(g, p, out)
		case replace || k < 2:
			for i := range out {
				if i&255 == 0 {
					s.checkInterrupt()
				}
				out[i] = unifIndex(g, dn, rounding) + 1
			}
		default:
			x := make([]int32, n)
			for i := range x {
				x[i] = int32(i)
			}
			m := n
			for i := range out {
				if i&255 == 0 {
					s.checkInterrupt()
				}
				j := int(unifIndex(g, float64(m), rounding))
				out[i] = float64(x[j]) + 1
				m--
				x[j] = x[m]
			}
		}
		s.rngEnd(g)
		if !isNull {
			dn = float64(n)
		}
		return sampleResult(out, dn)
	})
	regBuiltin("..sample2", func(s *Session, c *callInfo) Value {
		call := s.currentCall()
		dn := s.sampleN(c.args[0], call)
		if math.IsNaN(dn) || dn < 0 || dn > 4.5e15 {
			s.errorf(call, "invalid first argument")
		}
		k := s.sampleSize(c.args[1], call)
		if math.IsNaN(dn) || dn < 0 || dn > 4.5e15 || (k > 0 && dn == 0) {
			s.errorf(call, "invalid first argument")
		}
		if float64(k) > dn/2 {
			s.errorf(call, "This algorithm is for size <= n/2")
		}
		s.checkLen(call, k)
		g := s.rngBegin()
		_, samp := s.rngKinds()
		out := make([]float64, k)
		seen := make(map[float64]struct{}, k)
		for i := range out {
			if i&255 == 0 {
				s.checkInterrupt()
			}
			for {
				v := unifIndex(g, dn, samp == 0) + 1
				if _, dup := seen[v]; dup {
					continue
				}
				seen[v] = struct{}{}
				out[i] = v
				break
			}
		}
		s.rngEnd(g)
		return sampleResult(out, dn)
	})
}

// probSampleReplace draws with replacement by inversion on the cumulative
// probabilities sorted in descending order.
func (s *Session) probSampleReplace(g *random.MT19937, p []float64, out []float64) {
	n := len(p)
	perm := make([]int, n)
	for i := range perm {
		perm[i] = i + 1
	}
	revsort(p, perm)
	for i := 1; i < n; i++ {
		p[i] += p[i-1]
	}
	for i := range out {
		if i&255 == 0 {
			s.checkInterrupt()
		}
		ru := unifRand(g)
		j := 0
		for ; j < n-1; j++ {
			if ru <= p[j] {
				break
			}
		}
		out[i] = float64(perm[j])
	}
}

// probSampleNoReplace draws one at a time from the remaining mass, the
// probabilities sorted in descending order.
func (s *Session) probSampleNoReplace(g *random.MT19937, p []float64, out []float64) {
	n := len(p)
	perm := make([]int, n)
	for i := range perm {
		perm[i] = i + 1
	}
	revsort(p, perm)
	total := 1.0
	n1 := n - 1
	for i := range out {
		if i&255 == 0 {
			s.checkInterrupt()
		}
		rt := total * unifRand(g)
		mass := 0.0
		j := 0
		for ; j < n1; j++ {
			mass += p[j]
			if rt <= mass {
				break
			}
		}
		out[i] = float64(perm[j])
		total -= p[j]
		for m := j; m < n1; m++ {
			p[m] = p[m+1]
			perm[m] = perm[m+1]
		}
		n1--
	}
}

// walkerSample is Walker's alias method (Walker 1977), used with replacement
// when more than 200 probabilities are not tiny.
func (s *Session) walkerSample(g *random.MT19937, p []float64, out []float64, rounding bool) {
	n := len(p)
	q := make([]float64, n)
	a := make([]int, n)
	hl := make([]int, n)
	h, l := -1, n
	for i := 0; i < n; i++ {
		q[i] = p[i] * float64(n)
		if q[i] < 1 {
			h++
			hl[h] = i
		} else {
			l--
			hl[l] = i
		}
	}
	if h >= 0 && l < n {
		for k := 0; k < n-1; k++ {
			i, j := hl[k], hl[l]
			a[i] = j
			q[j] += q[i] - 1
			if q[j] < 1 {
				l++
			}
			if l >= n {
				break
			}
		}
	}
	for i := range out {
		if i&255 == 0 {
			s.checkInterrupt()
		}
		k := int(unifIndex(g, float64(n), rounding))
		if unifRand(g) < q[k] {
			out[i] = float64(k + 1)
		} else {
			out[i] = float64(a[k] + 1)
		}
	}
}
