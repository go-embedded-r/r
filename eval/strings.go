package eval

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/go-embedded-r/r/syntax"
)

func syntaxDeparse(v Value) []string { return syntax.Deparse(toSyntax(v), nil) }

// charVec is as.character for the character functions: strings stay, other
// atomic vectors are converted, anything else is an error with R's wording.
func (s *Session) charVec(c *callInfo, v Value, fname string) []Str {
	switch x := v.(type) {
	case *String:
		return x.V
	case *Logical, *Integer, *Double:
		if cl := classAttr(v); len(cl) == 1 && cl[0] == "factor" {
			s.errorf(c.call, "'%s' requires a character vector", fname)
		}
		return asStrings(v)
	case nullT:
		return nil
	case *List:
		return asStringsList(v)
	}
	s.errorf(c.call, "non-character argument")
	return nil
}

// keepShape returns res with the names, dim and dimnames of x.
func keepShape(res Value, x Value) Value {
	var attr []Attr
	for _, a := range attrsOf(x) {
		if a.Name == "names" || a.Name == "dim" || a.Name == "dimnames" {
			attr = append(attr, a)
		}
	}
	if attr == nil {
		return res
	}
	return withAttr(res, attr)
}

// runeWidth is the display width of a character: 2 for East Asian wide ones.
func runeWidth(r rune) int {
	switch {
	case r == 0:
		return 0
	case unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r) || unicode.Is(unicode.Cf, r):
		return 0
	case r >= 0x1100 && (r <= 0x115f || r == 0x2329 || r == 0x232a ||
		(r >= 0x2e80 && r <= 0xa4cf && r != 0x303f) ||
		(r >= 0xac00 && r <= 0xd7a3) || (r >= 0xf900 && r <= 0xfaff) ||
		(r >= 0xfe30 && r <= 0xfe6f) || (r >= 0xff00 && r <= 0xff60) ||
		(r >= 0xffe0 && r <= 0xffe6) || (r >= 0x1f300 && r <= 0x1f64f) ||
		(r >= 0x1f900 && r <= 0x1f9ff) || (r >= 0x20000 && r <= 0x3fffd)):
		return 2
	}
	return 1
}

func displayWidth(str string) int {
	w := 0
	for _, r := range str {
		w += runeWidth(r)
	}
	return w
}

func substrRunes(str string, start, stop int) string {
	if start < 1 {
		start = 1
	}
	if stop < start {
		return ""
	}
	rs := []rune(str)
	if start > len(rs) {
		return ""
	}
	if stop > len(rs) {
		stop = len(rs)
	}
	return string(rs[start-1 : stop])
}

func (s *Session) intsOf(c *callInfo, v Value, what string) []int32 {
	switch v.(type) {
	case *Integer, *Double, *Logical:
		d := asDoubles(v)
		out := make([]int32, len(d))
		for i, f := range d {
			switch {
			case math.IsNaN(f):
				out[i] = naInt
			case f > math.MaxInt32:
				out[i] = math.MaxInt32
			case f < -math.MaxInt32:
				out[i] = -math.MaxInt32
			default:
				out[i] = int32(f)
			}
		}
		return out
	}
	s.errorf(c.call, "invalid substring argument")
	return nil
}

func init() {
	regBuiltin("nchar", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "type", "allowNA", "keepNA")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		typ := "chars"
		if a[1] != nil {
			typ = s.stringArg(c, a[1], "type")
			ok := false
			for _, t := range []string{"bytes", "chars", "width"} {
				if strings.HasPrefix(t, typ) && typ != "" {
					typ, ok = t, true
					break
				}
			}
			if !ok {
				s.errorf(c.call, "invalid 'type' argument")
			}
		}
		keepNA := int32(naInt)
		if a[3] != nil {
			keepNA = logicalOf(s, c, a[3], "invalid 'keepNA' argument")[0]
		}
		x := s.charVec(c, a[0], "nchar()")
		out := make([]int32, len(x))
		for i, e := range x {
			if e.NA {
				na := keepNA == 1 || (keepNA == naInt && typ != "width")
				if na {
					out[i] = naInt
				} else {
					out[i] = 2
				}
				continue
			}
			switch typ {
			case "bytes":
				out[i] = int32(len(e.S))
			case "width":
				out[i] = int32(displayWidth(e.S))
			default:
				out[i] = int32(utf8.RuneCountInString(e.S))
			}
		}
		return keepShape(&Integer{V: out}, a[0])
	})
	regBuiltin("substr", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "start", "stop")
		if a[0] == nil || a[1] == nil || a[2] == nil {
			s.errorf(c.call, "argument \"%s\" is missing, with no default", []string{"x", "start", "stop"}[map[bool]int{true: 0, false: map[bool]int{true: 1, false: 2}[a[1] == nil]}[a[0] == nil]])
		}
		x := s.charVec(c, a[0], "substr")
		st, sp := s.intsOf(c, a[1], "start"), s.intsOf(c, a[2], "stop")
		if len(x) > 0 && (len(st) == 0 || len(sp) == 0) {
			s.errorf(c.call, "invalid substring arguments")
		}
		out := make([]Str, len(x))
		for i, e := range x {
			b, f := st[i%len(st)], sp[i%len(sp)]
			if e.NA || b == naInt || f == naInt {
				out[i] = Str{NA: true}
				continue
			}
			out[i] = Str{S: substrRunes(e.S, int(b), int(f))}
		}
		return withAttr(&String{V: out}, attrsOf(a[0]))
	})
	regBuiltin("substring", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "text", "first", "last")
		if a[0] == nil || a[1] == nil {
			s.errorf(c.call, "argument \"%s\" is missing, with no default", map[bool]string{true: "text", false: "first"}[a[0] == nil])
		}
		x := s.charVec(c, a[0], "substring")
		st := s.intsOf(c, a[1], "first")
		sp := []int32{1000000}
		if a[2] != nil {
			sp = s.intsOf(c, a[2], "last")
		}
		if len(x) == 0 || len(st) == 0 || len(sp) == 0 {
			return &String{V: []Str{}}
		}
		n := max(len(x), len(st), len(sp))
		out := make([]Str, n)
		for i := range out {
			e, b, f := x[i%len(x)], st[i%len(st)], sp[i%len(sp)]
			if e.NA || b == naInt || f == naInt {
				out[i] = Str{NA: true}
				continue
			}
			out[i] = Str{S: substrRunes(e.S, int(b), int(f))}
		}
		res := &String{V: out}
		if n == len(x) {
			if nm := namesOf(a[0]); nm != nil {
				res.Attr = namesAttr(nm.V)
			}
		}
		return res
	})
	caseFn := func(name string, f func(string) string) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			a := s.matchFormals(c, "x")
			if a[0] == nil {
				s.errorf(c.call, "argument \"x\" is missing, with no default")
			}
			x := s.charVec(c, a[0], name)
			out := make([]Str, len(x))
			for i, e := range x {
				if e.NA {
					out[i] = e
				} else {
					out[i] = Str{S: f(e.S)}
				}
			}
			if _, isStr := a[0].(*String); isStr {
				return withAttr(&String{V: out}, attrsOf(a[0]))
			}
			return keepShape(&String{V: out}, a[0])
		})
	}
	caseFn("toupper", strings.ToUpper)
	caseFn("tolower", strings.ToLower)
	regBuiltin("casefold", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "upper")
		up := a[1] != nil && s.truthy(c, a[1], "upper")
		name := "tolower"
		if up {
			name = "toupper"
		}
		return s.callWith(builtins[name], c.call, []Arg{{Value: a[0]}}, c.env)
	})
	regBuiltin("chartr", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "old", "new", "x")
		if a[0] == nil || a[1] == nil || a[2] == nil {
			s.errorf(c.call, "argument is missing, with no default")
		}
		old, nw := s.stringArg(c, a[0], "old"), s.stringArg(c, a[1], "new")
		expand := func(p string) []rune {
			rs := []rune(p)
			var out []rune
			for i := 0; i < len(rs); i++ {
				if i+2 < len(rs) && rs[i+1] == '-' {
					if rs[i] > rs[i+2] {
						s.errorf(c.call, "decreasing range specification ('%c-%c')", rs[i], rs[i+2])
					}
					for r := rs[i]; r <= rs[i+2]; r++ {
						out = append(out, r)
					}
					i += 2
					continue
				}
				out = append(out, rs[i])
			}
			return out
		}
		o, n := expand(old), expand(nw)
		if len(o) > len(n) {
			s.errorf(c.call, "'old' is longer than 'new'")
		}
		tr := map[rune]rune{}
		for i, r := range o {
			tr[r] = n[i]
		}
		x := s.charVec(c, a[2], "chartr")
		out := make([]Str, len(x))
		for i, e := range x {
			if e.NA {
				out[i] = e
				continue
			}
			out[i] = Str{S: strings.Map(func(r rune) rune {
				if t, ok := tr[r]; ok {
					return t
				}
				return r
			}, e.S)}
		}
		return withAttr(&String{V: out}, attrsOf(a[2]))
	})
	regBuiltin("trimws", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "which", "whitespace")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		which := "both"
		if a[1] != nil {
			which = s.stringArg(c, a[1], "which")
			ok := false
			for _, w := range []string{"both", "left", "right"} {
				if strings.HasPrefix(w, which) && which != "" {
					which, ok = w, true
					break
				}
			}
			if !ok {
				s.errorf(c.call, "'arg' should be one of %s", "‘both’, ‘left’, ‘right’")
			}
		}
		ws := "[ \t\r\n]"
		if a[2] != nil {
			ws = s.stringArg(c, a[2], "whitespace")
		}
		var left, right *matcher
		if which != "right" {
			left = s.compileRegex(c, "^("+ws+")+", false, false, false)
		}
		if which != "left" {
			right = s.compileRegex(c, "("+ws+")+$", false, false, false)
		}
		x := s.charVec(c, a[0], "trimws")
		out := make([]Str, len(x))
		for i, e := range x {
			if e.NA {
				out[i] = e
				continue
			}
			str := e.S
			if left != nil {
				if loc := left.find(str, 0); loc != nil {
					str = str[loc[1]:]
				}
			}
			if right != nil {
				if loc := right.find(str, 0); loc != nil {
					str = str[:loc[0]]
				}
			}
			out[i] = Str{S: str}
		}
		return withAttr(&String{V: out}, attrsOf(a[0]))
	})
	regBuiltin("strsplit", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "split", "fixed", "perl", "useBytes")
		x, ok := a[0].(*String)
		if !ok {
			s.errorf(c.call, "non-character argument")
		}
		fixed := a[2] != nil && s.truthy(c, a[2], "fixed")
		perl := a[3] != nil && s.truthy(c, a[3], "perl")
		var splits []Str
		if a[1] != nil {
			splits = s.charVec(c, a[1], "strsplit")
		}
		if len(splits) == 0 {
			splits = []Str{{S: ""}}
		}
		cache := map[string]*matcher{}
		out := make([]Value, len(x.V))
		for i, e := range x.V {
			if e.NA {
				out[i] = &String{V: []Str{{NA: true}}}
				continue
			}
			sp := splits[i%len(splits)]
			if sp.S == "" || sp.NA {
				var parts []Str
				for _, r := range e.S {
					parts = append(parts, Str{S: string(r)})
				}
				if parts == nil {
					parts = []Str{}
				}
				out[i] = &String{V: parts}
				continue
			}
			m := cache[sp.S]
			if m == nil {
				m = s.compileRegex(c, sp.S, fixed, false, perl)
				cache[sp.S] = m
			}
			var parts []Str
			rest := e.S
			for len(rest) > 0 {
				loc := m.find(rest, 0)
				if loc == nil {
					parts = append(parts, Str{S: rest})
					break
				}
				if loc[1] == 0 {
					// the match ends at the start: one character is split off
					n := nextChar(rest, 0)
					parts = append(parts, Str{S: rest[:n]})
					rest = rest[n:]
					continue
				}
				parts = append(parts, Str{S: rest[:loc[0]]})
				rest = rest[loc[1]:]
			}
			if parts == nil {
				parts = []Str{}
			}
			out[i] = &String{V: parts}
		}
		return listOfNamed(out, namesOf(x))
	})
	regBuiltin("Encoding", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x")
		x, ok := a[0].(*String)
		if !ok {
			s.errorf(c.call, "a character vector argument expected")
		}
		out := make([]Str, len(x.V))
		for i, e := range x.V {
			out[i] = Str{S: "unknown"}
			for j := 0; j < len(e.S); j++ {
				if e.S[j] >= 0x80 {
					out[i] = Str{S: "UTF-8"}
					break
				}
			}
		}
		return &String{V: out}
	})
	for _, n := range []string{"enc2utf8", "enc2native"} {
		n := n
		regBuiltin(n, func(s *Session, c *callInfo) Value {
			a := s.matchFormals(c, "x")
			if _, ok := a[0].(*String); !ok {
				s.errorf(c.call, "argument is not a character vector")
			}
			return a[0]
		})
	}
	regBuiltin("validUTF8", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x")
		x, ok := a[0].(*String)
		if !ok {
			s.errorf(c.call, "invalid 'x' argument")
		}
		out := make([]int32, len(x.V))
		for i, e := range x.V {
			if utf8.ValidString(e.S) {
				out[i] = 1
			}
		}
		return &Logical{V: out}
	})
	regBuiltin("utf8ToInt", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x")
		x, ok := a[0].(*String)
		if !ok || len(x.V) != 1 {
			s.errorf(c.call, "argument must be a character vector of length 1")
		}
		if x.V[0].NA {
			return &Integer{V: []int32{naInt}}
		}
		if !utf8.ValidString(x.V[0].S) {
			return &Integer{V: []int32{naInt}}
		}
		out := []int32{}
		for _, r := range x.V[0].S {
			out = append(out, int32(r))
		}
		return &Integer{V: out}
	})
	regBuiltin("intToUtf8", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "multiple", "allow_surrogate_pairs")
		d := asDoublesOrNil(a[0])
		multiple := a[1] != nil && s.truthy(c, a[1], "multiple")
		if multiple {
			out := make([]Str, len(d))
			for i, f := range d {
				if f != f {
					out[i] = Str{NA: true}
				} else if f == 0 {
					out[i] = Str{S: ""}
				} else {
					out[i] = Str{S: string(rune(int(f)))}
				}
			}
			return &String{V: out}
		}
		var b strings.Builder
		for _, f := range d {
			if f != f {
				return &String{V: []Str{{NA: true}}}
			}
			if f != 0 {
				b.WriteRune(rune(int(f)))
			}
		}
		return mkString(b.String())
	})
	regBuiltin("sprintf", func(s *Session, c *callInfo) Value { return s.sprintf(c) })
	regBuiltin("format", func(s *Session, c *callInfo) Value { return s.format(c) })
}

// ---- sprintf -----------------------------------------------------------

// maxFieldWidth bounds one formatted field: a width of 10^9 would be a
// memory bomb.
const maxFieldWidth = 1 << 24

func (s *Session) sprintf(c *callInfo) Value {
	if len(c.args) == 0 {
		s.errorf(c.call, "'fmt' is not a character vector")
	}
	fi := 0
	for i, n := range c.names {
		if n == "fmt" {
			fi = i
		}
	}
	fmtV, ok := c.args[fi].(*String)
	if !ok {
		s.errorf(c.call, "'fmt' is not a character vector")
	}
	var args []Value
	for i, a := range c.args {
		if i != fi {
			args = append(args, a)
		}
	}
	n := len(fmtV.V)
	for _, a := range args {
		if Length(a) == 0 {
			return &String{V: []Str{}}
		}
		if Length(a) > n {
			n = Length(a)
		}
	}
	if len(fmtV.V) == 0 {
		return &String{V: []Str{}}
	}
	used := make([]bool, len(args))
	out := make([]Str, n)
	for i := 0; i < n; i++ {
		f := fmtV.V[i%len(fmtV.V)]
		if f.NA {
			out[i] = Str{NA: true}
			continue
		}
		str, na := s.sprintf1(c, f.S, args, i, used)
		if na {
			out[i] = Str{NA: true}
		} else {
			out[i] = Str{S: str}
		}
	}
	unused := 0
	for _, u := range used {
		if !u {
			unused++
		}
	}
	if unused == 1 {
		s.warningf(c.call, "one argument not used by format '%s'", fmtV.V[0].S)
	} else if unused > 1 {
		s.warningf(c.call, "%d arguments not used by format '%s'", unused, fmtV.V[0].S)
	}
	return &String{V: out}
}

func (s *Session) sprintf1(c *callInfo, f string, args []Value, row int, used []bool) (string, bool) {
	var b strings.Builder
	next := 0
	for i := 0; i < len(f); i++ {
		ch := f[i]
		if ch != '%' {
			b.WriteByte(ch)
			continue
		}
		i++
		if i >= len(f) {
			s.errorf(c.call, "unrecognised format specification '%s'", f[i-1:])
		}
		if f[i] == '%' {
			b.WriteByte('%')
			continue
		}
		start := i - 1
		// n$
		argIdx := -1
		j := i
		for j < len(f) && f[j] >= '0' && f[j] <= '9' {
			j++
		}
		if j < len(f) && j > i && f[j] == '$' {
			argIdx, _ = strconv.Atoi(f[i:j])
			argIdx--
			i = j + 1
		}
		flags := ""
		for i < len(f) && strings.IndexByte("-+ 0#", f[i]) >= 0 {
			flags += string(f[i])
			i++
		}
		width := ""
		if i < len(f) && f[i] == '*' {
			if next >= len(args) {
				s.errorf(c.call, "too few arguments")
			}
			w := asDoubles(args[next])
			used[next] = true
			next++
			if len(w) > 0 {
				wv := int(w[row%len(w)])
				if wv < 0 {
					flags += "-"
					wv = -wv
				}
				width = strconv.Itoa(wv)
			}
			i++
		} else {
			for i < len(f) && f[i] >= '0' && f[i] <= '9' {
				width += string(f[i])
				i++
			}
		}
		prec := ""
		hasPrec := false
		if i < len(f) && f[i] == '.' {
			hasPrec = true
			i++
			for i < len(f) && f[i] >= '0' && f[i] <= '9' {
				prec += string(f[i])
				i++
			}
			if prec == "" {
				prec = "0"
			}
		}
		for i < len(f) && (f[i] == 'l' || f[i] == 'h') {
			i++
		}
		for _, num := range []string{width, prec} {
			if v, err := strconv.Atoi(num); num != "" && (err != nil || v > maxFieldWidth) {
				s.errorf(c.call, "required resulting string length is greater than maximal %d", maxFieldWidth)
			}
		}
		if i >= len(f) {
			s.errorf(c.call, "unrecognised format specification '%s'", f[start:])
		}
		conv := f[i]
		spec := f[start : i+1]
		var arg Value
		if argIdx >= 0 {
			if argIdx >= len(args) {
				s.errorf(c.call, "reference to non-existent argument %d", argIdx+1)
			}
			arg = args[argIdx]
			used[argIdx] = true
		} else {
			if next >= len(args) {
				s.errorf(c.call, "too few arguments")
			}
			arg = args[next]
			used[next] = true
			next++
		}
		idx := row % Length(arg)
		goFlags := "%" + flags + width
		if hasPrec {
			goFlags += "." + prec
		}
		isNAv := func() bool {
			switch v := arg.(type) {
			case *Integer:
				return v.V[idx] == naInt
			case *Logical:
				return v.V[idx] == naInt
			case *Double:
				return math.IsNaN(v.V[idx]) && isNAReal(v.V[idx])
			case *String:
				return v.V[idx].NA
			}
			return false
		}
		padStr := func(str string) string {
			w, _ := strconv.Atoi(width)
			n := utf8.RuneCountInString(str)
			if n >= w {
				return str
			}
			if strings.Contains(flags, "-") {
				return str + strings.Repeat(" ", w-n)
			}
			if strings.Contains(flags, "0") && conv != 's' {
				return strings.Repeat("0", w-n) + str
			}
			return strings.Repeat(" ", w-n) + str
		}
		switch conv {
		case 'd', 'i':
			switch v := arg.(type) {
			case *Integer, *Logical:
				if isNAv() {
					b.WriteString(padStr("NA"))
					continue
				}
				b.WriteString(fmt.Sprintf(goFlags+"d", asInts(v)[idx]))
			case *Double:
				x := v.V[idx]
				if isNAv() {
					b.WriteString(padStr("NA"))
					continue
				}
				if math.IsNaN(x) || math.IsInf(x, 0) || x != math.Trunc(x) {
					s.errorf(c.call, "invalid format '%s'; use format %%f, %%e, %%g or %%a for numeric objects", spec)
				}
				b.WriteString(fmt.Sprintf(goFlags+"d", int64(x)))
			default:
				s.errorf(c.call, "invalid format '%s'; use format %%s for character objects", spec)
			}
		case 'f', 'e', 'E', 'g', 'G', 'a', 'A':
			var x float64
			switch v := arg.(type) {
			case *Integer, *Logical:
				if isNAv() {
					b.WriteString(padStr("NA"))
					continue
				}
				x = float64(asInts(v)[idx])
			case *Double:
				x = v.V[idx]
				if isNAv() {
					b.WriteString(padStr("NA"))
					continue
				}
			default:
				s.errorf(c.call, "invalid format '%s'; use format %%s for character objects", spec)
			}
			switch {
			case math.IsNaN(x):
				b.WriteString(padStr("NaN"))
			case math.IsInf(x, 1):
				if strings.Contains(flags, "+") {
					b.WriteString(padStr("+Inf"))
				} else if strings.Contains(flags, " ") {
					b.WriteString(padStr(" Inf"))
				} else {
					b.WriteString(padStr("Inf"))
				}
			case math.IsInf(x, -1):
				b.WriteString(padStr("-Inf"))
			default:
				p := goFlags
				if !hasPrec && (conv == 'g' || conv == 'G') {
					p += ".6"
				} else if !hasPrec && (conv == 'f' || conv == 'e' || conv == 'E') {
					p += ".6"
				}
				if conv == 'a' || conv == 'A' {
					conv = 'x'
				}
				b.WriteString(fmt.Sprintf(p+string(conv), x))
			}
		case 'x', 'X', 'o':
			var x int64
			switch v := arg.(type) {
			case *Logical:
				s.errorf(c.call, "invalid format '%s'; use format %%d or %%i for logical objects", spec)
			case *Integer:
				if isNAv() {
					b.WriteString(padStr("NA"))
					continue
				}
				x = int64(uint32(asInts(v)[idx]))
			case *Double:
				if isNAv() {
					b.WriteString(padStr("NA"))
					continue
				}
				if v.V[idx] != math.Trunc(v.V[idx]) {
					s.errorf(c.call, "invalid format '%s'; use format %%f, %%e, %%g or %%a for numeric objects", spec)
				}
				x = int64(v.V[idx])
			default:
				s.errorf(c.call, "invalid format '%s'; use format %%s for character objects", spec)
			}
			b.WriteString(fmt.Sprintf(goFlags+string(conv), x))
		case 's':
			var str string
			switch v := arg.(type) {
			case *String:
				if v.V[idx].NA {
					str = "NA"
				} else {
					str = v.V[idx].S
				}
			case *Logical, *Integer, *Double:
				e := asStrings(v)[idx]
				if e.NA {
					str = "NA"
				} else {
					str = e.S
				}
			case *List:
				e := asStringsList(v)[idx]
				str = e.S
			default:
				s.errorf(c.call, "unsupported type")
			}
			if hasPrec {
				p, _ := strconv.Atoi(prec)
				if rs := []rune(str); len(rs) > p {
					str = string(rs[:p])
				}
			}
			b.WriteString(padStr(str))
		default:
			s.errorf(c.call, "unrecognised format specification '%s'", spec)
		}
	}
	return b.String(), false
}

// ---- format ------------------------------------------------------------

func (s *Session) format(c *callInfo) Value {
	a := s.matchFormals(c, "x", "trim", "digits", "nsmall", "justify", "width", "na.encode", "scientific")
	x := a[0]
	if x == nil {
		s.errorf(c.call, "argument \"x\" is missing, with no default")
	}
	trim := a[1] != nil && s.truthy(c, a[1], "trim")
	digits := printDigits
	if a[2] != nil {
		if _, isNull := a[2].(nullT); !isNull {
			digits = s.intArg(c, a[2], "digits")
			if digits < 0 || digits > 22 {
				s.errorf(c.call, "invalid value %d for 'digits' argument", digits)
			}
		}
	}
	nsmall := 0
	if a[3] != nil {
		nsmall = s.intArg(c, a[3], "nsmall")
		if nsmall < 0 || nsmall > 20 {
			s.errorf(c.call, "invalid 'nsmall' argument")
		}
	}
	justify := "left"
	if a[4] != nil {
		j := s.stringArg(c, a[4], "justify")
		ok := false
		for _, w := range []string{"left", "right", "centre", "none"} {
			if strings.HasPrefix(w, j) && j != "" {
				justify, ok = w, true
				break
			}
		}
		if !ok {
			s.errorf(c.call, "'arg' should be one of %s", "‘left’, ‘right’, ‘centre’, ‘none’")
		}
	}
	width := 0
	if a[5] != nil {
		if _, isNull := a[5].(nullT); !isNull {
			width = s.intArg(c, a[5], "width")
			if width > maxFieldWidth {
				s.errorf(c.call, "invalid 'width' argument")
			}
		}
	}
	scipen := 0
	if a[7] != nil {
		if l, ok := a[7].(*Logical); ok && len(l.V) == 1 && l.V[0] != naInt {
			if l.V[0] != 0 {
				scipen = -100
			} else {
				scipen = 100
			}
		}
	}
	var items []string
	left := false
	switch v := x.(type) {
	case nullT:
		return mkString("NULL")
	case *List:
		out := make([]Str, len(v.V))
		for i, e := range v.V {
			if _, isNull := e.(nullT); isNull {
				out[i] = Str{S: "NULL"}
				continue
			}
			if !isVector(e) {
				s.errorf(c.call, "unsupported object for format()")
			}
			parts := s.format(&callInfo{call: c.call, env: c.env, args: []Value{e, mkLogical(true)}, names: []string{"", "trim"}})
			ps := parts.(*String).V
			var sb strings.Builder
			for k, p := range ps {
				if k > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(p.S)
			}
			out[i] = Str{S: sb.String()}
		}
		return keepShape(&String{V: out}, x)
	case *Logical, *Integer:
		if cl := classAttr(x); len(cl) == 1 && cl[0] == "factor" {
			lev, _ := getAttr(x, "levels").(*String)
			for _, code := range asInts(x) {
				if code == naInt || lev == nil || int(code) > len(lev.V) {
					items = append(items, "NA")
				} else {
					items = append(items, lev.V[code-1].S)
				}
			}
			left = true
			break
		}
		items, _ = formatElements(v, false)
	case *Double:
		f := formatRealVecSci(v.V, digits, scipen)
		if nsmall > 0 && !f.e && f.d < nsmall {
			f.w += nsmall - f.d
			if f.d == 0 {
				f.w++
			}
			f.d = nsmall
		}
		for _, e := range v.V {
			items = append(items, strings.TrimLeft(f.encode(e), " "))
		}
	case *String:
		left = true
		for _, e := range v.V {
			if e.NA {
				items = append(items, "NA")
			} else {
				items = append(items, e.S)
			}
		}
	default:
		s.errorf(c.call, "unsupported object for format()")
	}
	// common width
	w := width
	if !trim || left {
		for _, it := range items {
			if n := displayWidth(it); n > w && (!trim || left) {
				w = n
			}
		}
	}
	if trim && !left {
		w = width
	}
	out := make([]Str, len(items))
	for i, it := range items {
		n := displayWidth(it)
		pad := w - n
		if pad < 0 {
			pad = 0
		}
		switch {
		case left && justify == "left":
			it += strings.Repeat(" ", pad)
		case left && justify == "right":
			it = strings.Repeat(" ", pad) + it
		case left && justify == "centre":
			l := pad / 2
			it = strings.Repeat(" ", l) + it + strings.Repeat(" ", pad-l)
		case left:
		default:
			it = strings.Repeat(" ", pad) + it
		}
		out[i] = Str{S: it}
	}
	return keepShape(&String{V: out}, x)
}
