package eval

import (
	"math"
	"math/big"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/go-embedded-r/r/syntax"
)

// Slice 22: the small natives that need no host and no shared engine. Written
// from the help pages of R 4.6.2 and the oracle's answers (docs/m3-plan.md).

// ..charmatch(x, table, nomatch): the position of the unique exact match, else
// of the unique partial (prefix) match; 0 when several match, nomatch when none.
func init() {
	regBuiltin("..charmatch", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "table", "nomatch")
		x, ok1 := a[0].(*String)
		t, ok2 := a[1].(*String)
		if !ok1 || !ok2 || a[2] == nil {
			s.errorf(c.call, "invalid arguments")
		}
		nm := asInts(coerceTo(a[2], kInteger))
		nomatch := int(naInt)
		if len(nm) > 0 {
			nomatch = int(nm[0])
		}
		text := func(e Str) string {
			if e.NA {
				return "NA"
			}
			return e.S
		}
		out := make([]int32, len(x.V))
		for i, xe := range x.V {
			xs := text(xe)
			exact, partial := 0, 0
			ie, ip := 0, 0
			for j, te := range t.V {
				ts := text(te)
				if !strings.HasPrefix(ts, xs) {
					continue
				}
				if len(ts) == len(xs) {
					exact++
					if exact == 1 {
						ie = j + 1
					}
				} else {
					partial++
					if partial == 1 {
						ip = j + 1
					}
				}
			}
			switch {
			case exact == 1:
				out[i] = int32(ie)
			case exact > 1:
				out[i] = 0
			case partial == 1:
				out[i] = int32(ip)
			case partial > 1:
				out[i] = 0
			default:
				out[i] = int32(nomatch)
			}
		}
		return &Integer{V: out}
	})
}

// ..findInterval(x, vec, rightmost.closed, all.inside, left.open): the number
// of knots of vec at or below x (below x when left.open), as findInterval
// documents. vec must be sorted and free of NA.
func init() {
	regBuiltin("..findInterval", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "vec", "rightmost.closed", "all.inside", "left.open")
		for _, v := range a {
			if v == nil {
				s.errorf(c.call, "5 arguments required")
			}
		}
		x, vec := asDoubles(a[0]), asDoubles(a[1])
		rc, ai, lo := s.truthy(c, a[2], "rightmost.closed"), s.truthy(c, a[3], "all.inside"), s.truthy(c, a[4], "left.open")
		n := len(vec)
		for i, v := range vec {
			if math.IsNaN(v) || (i > 0 && v < vec[i-1]) {
				s.errorf(c.call, "'vec' must be sorted non-decreasingly and not contain NAs")
			}
		}
		out := make([]int32, len(x))
		for i, t := range x {
			if math.IsNaN(t) {
				out[i] = naInt
				continue
			}
			// count the knots at or below t (below, when left.open)
			lo2, hi := 0, n
			for lo2 < hi {
				mid := (lo2 + hi) / 2
				if (lo && vec[mid] < t) || (!lo && vec[mid] <= t) {
					lo2 = mid + 1
				} else {
					hi = mid
				}
			}
			k := lo2
			if rc && n > 0 {
				if !lo && t == vec[n-1] {
					k = n - 1
				}
				if lo && t == vec[0] {
					k = 1
				}
			}
			if ai {
				// R tests the two ends in turn, never both: with one knot
				// the answer is 1 below it and 0 at or above it
				if k < 1 {
					k = 1
				} else if k > n-1 {
					k = n - 1
				}
			}
			out[i] = int32(k)
		}
		return &Integer{V: out}
	})
}

// chooseExact is choose(n, k) for 0 <= k <= 1000, correctly rounded: an exact
// integer when n is an integer, an exact rational otherwise.
func chooseExact(n float64, k int) float64 {
	if n == math.Trunc(n) {
		// symmetric: use the smaller of k and n-k (n >= 0 here)
		if n >= 0 {
			if float64(k) > n {
				return 0
			}
			if float64(k) > n-float64(k) {
				k = int(n) - k
			}
		}
		r := big.NewInt(1)
		bn := new(big.Float).SetFloat64(n)
		ni, _ := bn.Int(nil)
		tmp := new(big.Int)
		for j := 0; j < k; j++ {
			tmp.Sub(ni, big.NewInt(int64(j)))
			r.Mul(r, tmp)
			r.Quo(r, big.NewInt(int64(j+1)))
		}
		f, _ := new(big.Float).SetInt(r).Float64()
		return f
	}
	num := big.NewRat(1, 1)
	bn := new(big.Rat)
	bn.SetFloat64(n)
	for j := 0; j < k; j++ {
		num.Mul(num, new(big.Rat).Sub(bn, big.NewRat(int64(j), 1)))
	}
	fact := big.NewInt(1)
	for j := 2; j <= k; j++ {
		fact.Mul(fact, big.NewInt(int64(j)))
	}
	num.Quo(num, new(big.Rat).SetInt(fact))
	f, _ := num.Float64()
	return f
}

// chooseOne is choose(n, k) for a rounded k.
func chooseOne(n, k float64) float64 {
	switch {
	case k < 0:
		return 0
	case k == 0:
		return 1
	case math.IsInf(n, 0):
		if n > 0 || math.Mod(k, 2) == 0 {
			return math.Inf(1)
		}
		return math.Inf(-1)
	}
	isInt := math.Abs(n-math.RoundToEven(n)) <= 1e-7*math.Max(1, math.Abs(n))
	if isInt {
		n = math.RoundToEven(n)
		if n < 0 {
			r := chooseOne(-n+k-1, k)
			if math.Mod(k, 2) != 0 {
				r = -r
			}
			return r
		}
		if k > n {
			return 0
		}
		if k > n-k {
			k = n - k
		}
		if k <= 1000 {
			return chooseExact(n, int(k))
		}
	} else if k <= 1000 {
		return chooseExact(n, int(k))
	}
	// a long product: the gamma function
	l1, s1 := math.Lgamma(n + 1)
	l2, s2 := math.Lgamma(k + 1)
	l3, s3 := math.Lgamma(n - k + 1)
	r := math.Exp(l1 - l2 - l3)
	if isInt {
		r = math.RoundToEven(r)
	}
	return float64(s1*s2*s3) * r
}

func init() {
	regBuiltin("choose", func(s *Session, c *callInfo) Value {
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
		m := max(len(n), len(k))
		out := make([]float64, m)
		warned := false
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
				s.checkInterrupt() // each element can cost a thousand big-number steps
			}
			out[i] = chooseOne(ni, kr)
		}
		src := a[0]
		if len(k) > len(n) {
			src = a[1]
		}
		return &Double{V: out, Attr: attrsOf(src)}
	})
}

// numToBits and numToInts: the 64 bits of each double, the lowest first as 0
// and 1 raw bytes, and its two 32-bit halves, the lower word first.
func init() {
	regBuiltin("numToBits", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "numToBits")
		d := asDoubles(x)
		if len(d) > s.opts.MaxVector/64 {
			s.checkLen(c.call, s.opts.MaxVector+1) // 64 bytes per number: refused before anything is built
		}
		out := make([]byte, 0, 64*len(d))
		for _, f := range d {
			u := math.Float64bits(f)
			for i := 0; i < 64; i++ {
				out = append(out, byte(u>>uint(i))&1)
			}
		}
		return &Raw{V: out}
	})
	regBuiltin("numToInts", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "numToInts")
		d := asDoubles(x)
		if len(d) > s.opts.MaxVector/2 {
			s.checkLen(c.call, s.opts.MaxVector+1)
		}
		out := make([]int32, 0, 2*len(d))
		for _, f := range d {
			u := math.Float64bits(f)
			out = append(out, int32(uint32(u)), int32(uint32(u>>32)))
		}
		return &Integer{V: out}
	})
}

// Sys.getpid: the process number. validEnc: whether each string is valid in
// its encoding (UTF-8 here).
func init() {
	regBuiltin("Sys.getpid", func(s *Session, c *callInfo) Value {
		if len(c.args) != 0 {
			s.errorf(c.call, "unused argument")
		}
		return mkInt(int32(os.Getpid()))
	})
	regBuiltin("validEnc", func(s *Session, c *callInfo) Value {
		x, ok := s.nargs1(c, "validEnc").(*String)
		if !ok {
			s.errorf(c.call, "invalid 'x' argument")
		}
		out := make([]int32, len(x.V))
		for i, e := range x.V {
			if e.NA || utf8.ValidString(e.S) {
				out[i] = 1
			}
		}
		return &Logical{V: out}
	})
}

// The browser is never open here: its queries fail as R's do outside it.
func init() {
	for _, n := range []string{"browserText", "browserCondition", "browserSetDebug"} {
		regBuiltin(n, func(s *Session, c *callInfo) Value {
			a := s.matchFormals(c, "n")
			if a[0] != nil {
				if v := asInts(coerceTo(a[0], kInteger)); len(v) != 1 || v[0] == naInt {
					s.errorf(c.call, "invalid 'n' argument")
				}
			}
			s.errorf(c.call, "no browser context to query")
			return Null
		})
	}
}

// encodeString(x, width, quote, na.encode, justify): the strings as print
// writes them, padded to width.
func init() {
	regBuiltin("encodeString", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "width", "quote", "na.encode", "justify")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		x, ok := a[0].(*String)
		if !ok {
			switch a[0].(type) {
			case *Logical, *Integer, *Double:
				x, _ = coerceTo(a[0], kString).(*String)
			}
			if x == nil {
				s.errorf(c.call, "a character vector argument expected")
			}
		}
		width := 0
		wNA := false
		if a[1] != nil {
			if _, isNull := a[1].(nullT); !isNull {
				w := asInts(coerceTo(a[1], kInteger))
				if len(w) != 1 {
					s.errorf(c.call, "invalid 'width' value")
				}
				if w[0] == naInt {
					wNA = true
				} else if w[0] < 0 {
					s.errorf(c.call, "invalid 'width' value")
				} else {
					width = int(w[0])
				}
			}
		}
		quote := byte(0)
		if a[2] != nil {
			q, ok := a[2].(*String)
			if !ok || len(q.V) != 1 || q.V[0].NA {
				s.errorf(c.call, "invalid 'quote' value")
			}
			if q.V[0].S != "" {
				quote = q.V[0].S[0]
			}
		}
		naEncode := true
		if a[3] != nil {
			naEncode = s.truthy(c, a[3], "na.encode")
		}
		just := 0
		if a[4] != nil {
			j, ok := a[4].(*String)
			if !ok || len(j.V) == 0 {
				s.errorf(c.call, "'arg' must be of length 1")
			}
			opts := []string{"left", "right", "centre", "none"}
			if len(j.V) == len(opts) {
				equal := true
				for i, e := range j.V {
					if e.S != opts[i] {
						equal = false
					}
				}
				if equal {
					j = &String{V: j.V[:1]}
				}
			}
			if len(j.V) != 1 {
				s.errorf(c.call, "'arg' must be of length 1")
			}
			found := -1
			for i, o := range opts {
				if strings.HasPrefix(o, j.V[0].S) && j.V[0].S != "" {
					if found >= 0 {
						found = -2
						break
					}
					found = i
				}
			}
			if found < 0 {
				s.errorf(c.call, "'arg' should be one of %s", `“left”, “right”, “centre”, “none”`)
			}
			just = found
		}
		enc := make([]string, len(x.V))
		for i, e := range x.V {
			switch {
			case e.NA && quote == 0:
				enc[i] = "<NA>"
			case e.NA:
				enc[i] = "NA"
			case quote == 0:
				enc[i] = encodeUnquoted(e.S)
			default:
				enc[i] = syntax.EncodeString(e.S, quote)
			}
		}
		if width > 0 && len(x.V) > 0 && float64(width)*float64(len(x.V)) > float64(s.opts.MaxVector) {
			s.checkLen(c.call, s.opts.MaxVector+1) // the padding would pass the vector limit
		}
		if wNA {
			width = 0
			for i, e := range enc {
				if x.V[i].NA && !naEncode {
					continue
				}
				n := displayWidth(e)
				if x.V[i].NA {
					n = 2 // an NA counts as "NA", even when it prints as <NA>
				}
				if n > width {
					width = n
				}
			}
		}
		out := make([]Str, len(x.V))
		for i, e := range enc {
			if x.V[i].NA && !naEncode {
				out[i] = Str{NA: true}
				continue
			}
			n := displayWidth(e)
			if n < width && just != 3 {
				switch just {
				case 0:
					e = e + strings.Repeat(" ", width-n)
				case 1:
					e = strings.Repeat(" ", width-n) + e
				case 2:
					l := (width - n) / 2
					e = strings.Repeat(" ", l) + e + strings.Repeat(" ", width-n-l)
				}
			}
			out[i] = Str{S: e}
		}
		return &String{V: out, Attr: attrsOf(x)}
	})
}

// format.info(x, digits, nsmall): the width, the digits after the point and
// the exponent digits that format() would use for the whole of x.
func init() {
	regBuiltin("format.info", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "digits", "nsmall")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		digits := printDigits
		if a[1] != nil {
			if _, isNull := a[1].(nullT); !isNull {
				d := asInts(coerceTo(a[1], kInteger))
				if len(d) != 1 || d[0] == naInt || d[0] < 1 || d[0] > 22 {
					s.errorf(c.call, "invalid 'digits' argument")
				}
				digits = int(d[0])
			}
		}
		nsmall := 0
		if a[2] != nil {
			n := asInts(coerceTo(a[2], kInteger))
			if len(n) != 1 || n[0] == naInt || n[0] < 0 || n[0] > 20 {
				s.errorf(c.call, "invalid 'nsmall' argument")
			}
			nsmall = int(n[0])
		}
		switch x := a[0].(type) {
		case *Logical:
			w := 0
			for _, e := range x.V {
				n := 5
				switch e {
				case naInt:
					n = 2
				case 0:
					n = 5
				default:
					n = 4
				}
				if n > w {
					w = n
				}
			}
			return mkInt(int32(w))
		case *Integer:
			items, _ := formatElements(x, false)
			return mkInt(int32(maxWidth(items)))
		case *String:
			w := 0
			for _, e := range x.V {
				n := 2
				if !e.NA {
					n = displayWidth(e.S)
				}
				if n > w {
					w = n
				}
			}
			return mkInt(int32(w))
		case *Double:
			f := formatRealVec(x.V, digits)
			if !f.e && f.d < nsmall {
				if f.d == 0 {
					f.w++
				}
				f.w += nsmall - f.d
				f.d = nsmall
			}
			e := 0
			if f.e {
				e = 1
				for _, v := range x.V {
					if math.IsNaN(v) || math.IsInf(v, 0) || v == 0 {
						continue
					}
					if m := math.Abs(math.Log10(math.Abs(v))); m >= 99.5 {
						e = 2
					}
				}
			}
			return &Integer{V: []int32{int32(f.w), int32(f.d), int32(e)}}
		}
		s.errorf(c.call, "atomic vector arguments only")
		return Null
	})
}
