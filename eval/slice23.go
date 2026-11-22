package eval

import (
	_ "embed"
	"math"
	"math/big"
	"sort"
	"strings"
	"sync/atomic"
	"unicode"
)

// Slice 23: lchoose, abbreviate, approximate matching (agrep), OlsonNames and
// the connection table (showConnections, closeAllConnections).

//go:embed olson.txt
var olsonText string

// olsonVersion is the tz database release the list was taken from.
const olsonVersion = "2026e"

// chooseRat is choose(n, k) exactly, for 0 <= k <= 1000 and finite n: the
// product n (n-1) ... (n-k+1) / k! over the rationals.
func chooseRat(n float64, k int) *big.Rat {
	bn := new(big.Rat)
	bn.SetFloat64(n)
	num := big.NewRat(1, 1)
	for j := 0; j < k; j++ {
		num.Mul(num, new(big.Rat).Sub(bn, big.NewRat(int64(j), 1)))
	}
	fact := big.NewInt(1)
	for j := 2; j <= k; j++ {
		fact.Mul(fact, big.NewInt(int64(j)))
	}
	return num.Quo(num, new(big.Rat).SetInt(fact))
}

// logRat is log|r| for a nonzero rational: the 53-bit mantissa m in [0.5, 1)
// and the exponent e of r = m 2^e give log(m) + e log 2. The absolute error is
// at most 4 u max(1, |result|) with u = 2^-53 (checked by TestLchooseAgainstBig).
func logRat(r *big.Rat) float64 {
	f := new(big.Float).SetPrec(200).SetRat(r)
	f.Abs(f)
	m := new(big.Float)
	e := f.MantExp(m)
	mf, _ := m.Float64()
	return math.Log(mf) + float64(e)*math.Ln2
}

// lchooseOne is log|choose(n, k)| for a rounded k.
func lchooseOne(n, k float64) float64 {
	switch {
	case k < 0:
		return math.Inf(-1)
	case k == 0:
		return 0
	case math.IsInf(n, 0):
		return math.NaN()
	case k == 1:
		return math.Log(math.Abs(n))
	}
	isInt := math.Abs(n-math.RoundToEven(n)) <= 1e-7*math.Max(1, math.Abs(n))
	if isInt {
		n = math.RoundToEven(n)
		if n < 0 {
			return lchooseOne(-n+k-1, k)
		}
		if k > n {
			return math.Inf(-1)
		}
		if k > n-k {
			k = n - k
		}
	}
	if k <= 1000 && math.Abs(n) < 1e300 {
		return logRat(chooseRat(n, int(k)))
	}
	l1, _ := math.Lgamma(n + 1)
	l2, _ := math.Lgamma(k + 1)
	l3, _ := math.Lgamma(n - k + 1)
	return l1 - l2 - l3
}

// abbrevVowel: the lowercase vowels abbreviate drops first.
func abbrevVowel(r rune) bool {
	switch r {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	}
	switch {
	case r >= 224 && r <= 229, r >= 232 && r <= 239, r >= 242 && r <= 246, r >= 248 && r <= 253:
		return true
	}
	switch r {
	case 257, 259, 261, 275, 277, 279, 281, 283, 297, 299, 301, 303, 305, 333, 335, 337, 361, 363, 365, 367, 369, 371:
		return true
	}
	return false
}

// abbreviateOne shortens one string to minlength characters by dropping, from
// the right, lowercase vowels, then the other lowercase letters, then
// whatever else is not the first character of a word.
func abbreviateOne(x string, minlength int, useClasses bool) string {
	rs := []rune(strings.Trim(x, " "))
	if len(rs) <= minlength {
		return string(rs)
	}
	n := 0
	first := make([]bool, len(rs))
	for i, r := range rs {
		first[i] = i == 0 || rs[i-1] == ' '
		if r != ' ' {
			n++
		}
	}
	drop := make([]bool, len(rs))
	phases := []func(rune) bool{
		func(r rune) bool { return abbrevVowel(r) },
		func(r rune) bool { return unicode.IsLower(r) },
		func(r rune) bool { return true },
	}
	if !useClasses {
		phases = phases[2:]
	}
	for _, ok := range phases {
		for i := len(rs) - 1; i >= 0 && n > minlength; i-- {
			if drop[i] || first[i] || rs[i] == ' ' || !ok(rs[i]) {
				continue
			}
			drop[i] = true
			n--
		}
	}
	var b strings.Builder
	for i, r := range rs {
		if !drop[i] && r != ' ' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// agrepTables counts the tables of counts agrepMatch has built.
var agrepTables atomic.Int64

// agrepBounds are the resolved limits of an approximate match; a negative
// field is no limit.
type agrepBounds struct{ cost, ins, del, sub, all int }

// agrepMaxCells bounds the table of counts agrepMatch builds for one text.
const agrepMaxCells = 1 << 24

// agrepClamp resolves the bounds: no limit is the sum of the lengths, and
// there can be no more deletions or substitutions than the pattern has
// characters, no more insertions than the text has; a larger bound only
// grows the table.
func agrepClamp(p, t []rune, b agrepBounds) (total, mi, md, ms int) {
	limit := func(v int) int {
		if v < 0 {
			return len(p) + len(t)
		}
		return v
	}
	total = min(limit(b.cost), limit(b.all))
	return total, min(limit(b.ins), total, len(t)), min(limit(b.del), total, len(p)), min(limit(b.sub), total, len(p))
}

// agrepCells is the size of the table agrepMatch would build for p, t and b
// (0 when the answer is known without one).
func agrepCells(p, t []rune, b agrepBounds) float64 {
	_, mi, md, ms := agrepClamp(p, t, b)
	if md >= len(p) {
		return 0
	}
	return float64(mi+1) * float64(md+1) * float64(ms+1) * float64(len(t)+1)
}

// agrepMatch reports whether pattern p occurs in text t within the bounds:
// some substring of t is reached from p by at most that many insertions
// (characters of t that p lacks), deletions (characters of p that t lacks)
// and substitutions, each at unit cost.
func agrepMatch(p, t []rune, b agrepBounds) bool {
	total, mi, md, ms := agrepClamp(p, t, b)
	if md >= len(p) {
		// deleting the whole pattern is within reach: it matches the empty substring
		return true
	}
	// the state of a cell is the set of (ins, del, sub) counts that reach it
	agrepTables.Add(1)
	di, dd, ds := mi+1, md+1, ms+1
	if float64(di)*float64(dd)*float64(ds)*float64(len(t)+1) > agrepMaxCells {
		// callers refuse such a table before asking; reaching here is a bug
		panic("agrepMatch: table too large")
	}
	cube := di * dd * ds
	idx := func(a, d, s int) int { return (a*dd+d)*ds + s }
	rows := [2][][]bool{}
	for r := range rows {
		rows[r] = make([][]bool, len(t)+1)
		for j := range rows[r] {
			rows[r][j] = make([]bool, cube)
		}
	}
	// row i = 0: matching the empty prefix of p anywhere costs nothing
	for j := range rows[0] {
		rows[0][j][idx(0, 0, 0)] = true
	}
	cur := 0
	for i := 1; i <= len(p); i++ {
		prev, next := rows[cur], rows[1-cur]
		for j := range next {
			clear(next[j])
		}
		for j := 0; j <= len(t); j++ {
			dst := next[j]
			for a := 0; a < di; a++ {
				for d := 0; d < dd; d++ {
					for s := 0; s < ds; s++ {
						if a+d+s > total {
							continue
						}
						// deletion: p[i-1] is missing from t
						if d > 0 && prev[j][idx(a, d-1, s)] {
							dst[idx(a, d, s)] = true
						}
						if j == 0 {
							continue
						}
						// insertion: t[j-1] has no counterpart in p
						if a > 0 && next[j-1][idx(a-1, d, s)] {
							dst[idx(a, d, s)] = true
						}
						// match or substitution
						if p[i-1] == t[j-1] && prev[j-1][idx(a, d, s)] {
							dst[idx(a, d, s)] = true
						}
						if s > 0 && p[i-1] != t[j-1] && prev[j-1][idx(a, d, s-1)] {
							dst[idx(a, d, s)] = true
						}
					}
				}
			}
		}
		cur = 1 - cur
	}
	for _, cell := range rows[cur] {
		for a := 0; a < di; a++ {
			for d := 0; d < dd; d++ {
				for s := 0; s < ds; s++ {
					if cell[idx(a, d, s)] && a+d+s <= total {
						return true
					}
				}
			}
		}
	}
	return false
}

func init() {
	// .rangeNum(..., na.rm, finite, isNumeric): the range of the flattened
	// arguments; a numeric vector drops the non-finite (finite) or the missing
	// (na.rm) values first, any other type is ordered by min and max, with
	// finite implying na.rm. The three controls are taken by exact name.
	regBuiltin(".rangeNum", func(s *Session, c *callInfo) Value {
		var naRm, finite, isNum Value
		data := &callInfo{call: c.call, env: c.env}
		for i, a := range c.args {
			switch c.names[i] {
			case "na.rm":
				naRm = a
			case "finite":
				finite = a
			case "isNumeric":
				isNum = a
			default:
				data.args = append(data.args, a)
				data.names = append(data.names, c.names[i])
			}
		}
		if naRm == nil || finite == nil || isNum == nil {
			s.errorf(c.call, "argument \"%s\" is missing, with no default", map[bool]string{true: "na.rm", false: "finite"}[naRm == nil])
		}
		call := func(name string, names []string, args ...Value) Value {
			return builtins[name].fn(s, &callInfo{call: c.call, env: c.env, args: args, names: names})
		}
		data.args = append(data.args, &Logical{V: []int32{1}})
		data.names = append(data.names, "recursive")
		x := builtins["c"].fn(s, data)
		fin := s.truthy(c, finite, "finite")
		rm := s.truthy(c, naRm, "na.rm")
		keep := func(pred string, negate bool) {
			p := call(pred, []string{""}, x)
			if negate {
				p = call("!", []string{""}, p)
			}
			x = call("[", []string{"", ""}, x, p)
		}
		if s.truthy(c, s.callFunction(isNum, []Value{x}), "isNumeric") {
			switch {
			case fin:
				keep("is.finite", false)
			case rm:
				keep("is.na", true)
			}
			return call("c", []string{"", ""}, call("min", []string{""}, x), call("max", []string{""}, x))
		}
		narm := &Logical{V: []int32{0}}
		if fin || rm {
			narm.V[0] = 1
		}
		nm := []string{"", "na.rm"}
		return call("c", []string{"", ""}, call("min", nm, x, narm), call("max", nm, x, narm))
	})
	regBuiltin("lchoose", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "n", "k")
		if a[0] == nil || a[1] == nil {
			s.errorf(c.call, "argument \"%s\" is missing, with no default", map[bool]string{true: "n", false: "k"}[a[0] == nil])
		}
		for _, v := range a {
			switch v.(type) {
			case *Logical, *Integer, *Double:
			default:
				s.errorf(c.call, "non-numeric argument to mathematical function")
			}
		}
		n, k := asDoubles(a[0]), asDoubles(a[1])
		if len(n) == 0 || len(k) == 0 {
			return &Double{V: []float64{}}
		}
		out := make([]float64, max(len(n), len(k)))
		warned, nanWarned := false, false
		for i := range out {
			ni, ki := n[i%len(n)], k[i%len(k)]
			if math.IsNaN(ni) || math.IsNaN(ki) {
				out[i] = ni + ki
				continue
			}
			kr := math.RoundToEven(ki)
			if math.Abs(ki-kr) > 1e-7 && !warned {
				warned = true
				s.warningf(c.call, "'k' (%.2f) must be integer, rounded to %.0f", ki, kr)
			}
			if i&255 == 0 {
				s.checkInterrupt()
			}
			out[i] = lchooseOne(ni, kr)
			if math.IsNaN(out[i]) && !nanWarned {
				nanWarned = true
				s.warningf(c.call, "NaNs produced")
			}
		}
		src := a[0]
		if len(k) > len(n) {
			src = a[1]
		}
		return &Double{V: out, Attr: attrsOf(src)}
	})

	// ..abbrev1(x, minlength, use.classes): the stripping step of abbreviate
	regBuiltin("..abbrev1", func(s *Session, c *callInfo) Value {
		x, _ := c.args[0].(*String)
		ml := int(asDoubles(c.args[1])[0])
		uc := asDoubles(c.args[2])[0] != 0
		out := make([]Str, len(x.V))
		for i, e := range x.V {
			if e.NA {
				out[i] = Str{NA: true}
				continue
			}
			if i&255 == 0 {
				s.checkInterrupt()
			}
			out[i] = Str{S: abbreviateOne(e.S, ml, uc)}
		}
		return &String{V: out}
	})
	regBuiltin("..nonascii", func(s *Session, c *callInfo) Value {
		x, _ := c.args[0].(*String)
		for _, e := range x.V {
			for i := 0; i < len(e.S); i++ {
				if e.S[i] >= 0x80 {
					return mkLogical(true)
				}
			}
		}
		return mkLogical(false)
	})

	// ..agrepGo(pattern, x, ignore.case, bounds): bounds is c(cost, ins, del,
	// sub, all) with NA for no limit, each already an integer count
	regBuiltin("..agrepGo", func(s *Session, c *callInfo) Value {
		p, _ := c.args[0].(*String)
		x, _ := c.args[1].(*String)
		ic := asDoubles(c.args[2])[0] != 0
		bd := asDoubles(c.args[3])
		get := func(i int) int {
			if math.IsNaN(bd[i]) {
				return -1
			}
			return int(min(bd[i], 1<<20))
		}
		// R's "insertions" are the pattern characters the text lacks and its
		// "deletions" the text characters the pattern lacks: the other way round
		// from agrepMatch's reading
		b := agrepBounds{get(0), get(2), get(1), get(3), get(4)}
		fold := func(str string) []rune {
			if ic {
				str = strings.ToLower(str)
			}
			return []rune(str)
		}
		pr := fold(p.V[0].S)
		out := make([]int32, len(x.V))
		for i, e := range x.V {
			if e.NA {
				continue
			}
			if i&15 == 0 {
				s.checkInterrupt()
			}
			tx := fold(e.S)
			if agrepCells(pr, tx, b) > agrepMaxCells {
				s.errorf(c.call, "approximate matching of a pattern of %d characters within these distances is too large", len(pr))
			}
			if agrepMatch(pr, tx, b) {
				out[i] = 1
			}
		}
		return &Logical{V: out}
	})

	regBuiltin("OlsonNames", func(s *Session, c *callInfo) Value {
		s.matchFormals(c, "tzdir")
		names := strings.Fields(olsonText)
		sort.Strings(names)
		r := strVec(names)
		r.Attr = []Attr{{Name: "Version", Value: mkString(olsonVersion)}}
		return r
	})

	regBuiltin("showConnections", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "all")
		all := s.boolArg(c, a[0], false, "all")
		var ids []int
		for id := range s.io.conns {
			if all || id >= 3 {
				ids = append(ids, id)
			}
		}
		sort.Ints(ids)
		cols := []string{"description", "class", "mode", "text", "isopen", "can read", "can write"}
		out := make([]Str, 0, len(ids)*len(cols))
		rn := make([]string, len(ids))
		cells := make([][]string, len(cols))
		for r, id := range ids {
			cn := s.io.conns[id]
			mode, text, opened, canRead, canWrite := s.modeName(cn)
			rn[r] = itoa(id)
			for j, v := range []string{cn.desc, cn.class, mode, text, opened, canRead, canWrite} {
				cells[j] = append(cells[j], v)
			}
		}
		for j := range cols {
			for _, v := range cells[j] {
				out = append(out, Str{S: v})
			}
		}
		var rnv Value = Null
		if len(rn) > 0 {
			rnv = strVec(rn)
		}
		return &String{V: out, Attr: []Attr{
			{Name: "dim", Value: &Integer{V: []int32{int32(len(ids)), int32(len(cols))}}},
			{Name: "dimnames", Value: &List{V: []Value{rnv, strVec(cols)}}},
		}}
	})

	regBuiltin("closeAllConnections", func(s *Session, c *callInfo) Value {
		s.flushSinks()
		s.io.sinks = [2][]sinkEntry{}
		var ids []int
		for id := range s.io.conns {
			if id >= 3 {
				ids = append(ids, id)
			}
		}
		sort.Ints(ids)
		for _, id := range ids {
			if cn, ok := s.io.conns[id]; ok {
				cn.sinked = 0
				s.closeConn(cn, true)
			}
		}
		return s.invisible(Null)
	})
}
