package eval

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

// The matcher type, its compilation and the search itself are in
// regex_tre.go; this file holds the R-level functions built on them.

// nextChar is the byte length of the character at str[i:].
func nextChar(str string, i int) int {
	if i >= len(str) {
		return 1
	}
	_, n := utf8.DecodeRuneInString(str[i:])
	return n
}

// gregexprMatches lists the matches of gregexpr: after an empty match the
// search resumes one character later, and it stops at the end of the text.
func (m *matcher) gregexprMatches(str string) [][]int {
	var out [][]int
	pos := 0
	for {
		loc := m.find(str, pos)
		if loc == nil {
			break
		}
		out = append(out, loc)
		if loc[1] == loc[0] {
			pos = loc[1] + nextChar(str, loc[1])
		} else {
			pos = loc[1]
		}
		if pos >= len(str) {
			break
		}
	}
	return out
}

// allMatches lists the matches of sub / gsub: after an empty
// match the search resumes one character later, and it stops at the end.
func (m *matcher) allMatches(str string) [][]int {
	var out [][]int
	pos, prevEnd := 0, -1
	for pos <= len(str) {
		loc := m.find(str, pos)
		if loc == nil {
			break
		}
		if loc[1] == loc[0] && loc[0] == prevEnd {
			// an empty match right after the previous match does not count
			if loc[0] >= len(str) {
				break
			}
			pos = loc[0] + nextChar(str, loc[0])
			continue
		}
		out = append(out, loc)
		prevEnd = loc[1]
		if loc[1] == loc[0] {
			if loc[1] >= len(str) {
				break
			}
			pos = loc[1] + nextChar(str, loc[1])
		} else {
			pos = loc[1]
		}
	}
	return out
}

func charIndex(str string, byteOff int) int { return utf8.RuneCountInString(str[:byteOff]) }

// expandReplacement applies \1..\9 (and \0) of a sub/gsub replacement.
func expandReplacement(repl, str string, loc []int) string {
	if !strings.Contains(repl, "\\") {
		return repl
	}
	var b strings.Builder
	for i := 0; i < len(repl); i++ {
		ch := repl[i]
		if ch == '\\' && i+1 < len(repl) {
			i++
			n := repl[i]
			if n >= '1' && n <= '9' {
				g := int(n - '0')
				if 2*g+1 < len(loc) && loc[2*g] >= 0 {
					b.WriteString(str[loc[2*g]:loc[2*g+1]])
				}
				continue
			}
			b.WriteByte(n)
			continue
		}
		b.WriteByte(ch)
	}
	return b.String()
}

type reArgs struct {
	icase, perl, fixed, invert, value bool
}

func (s *Session) reFlags(c *callInfo, a map[string]Value) reArgs {
	var r reArgs
	get := func(n string) bool {
		if v, ok := a[n]; ok && v != nil {
			return s.truthy(c, v, n)
		}
		return false
	}
	r.icase, r.perl, r.fixed, r.invert, r.value = get("ignore.case"), get("perl"), get("fixed"), get("invert"), get("value")
	return r
}

func (s *Session) patternOf(c *callInfo, v Value) (string, bool) {
	st, ok := v.(*String)
	if !ok {
		if v != nil && isVector(v) {
			vs := asStrings(v)
			if len(vs) > 0 {
				if len(vs) > 1 {
					s.warningf(c.call, "argument 'pattern' has length > 1 and only the first element will be used")
				}
				return vs[0].S, !vs[0].NA
			}
		}
		s.errorf(c.call, "invalid 'pattern' argument")
	}
	if len(st.V) < 1 {
		s.errorf(c.call, "invalid 'pattern' argument")
	}
	if len(st.V) > 1 {
		s.warningf(c.call, "argument 'pattern' has length > 1 and only the first element will be used")
	}
	return st.V[0].S, !st.V[0].NA
}

func (s *Session) textOf(c *callInfo, v Value) []Str {
	if v == nil {
		s.errorf(c.call, "argument \"x\" is missing, with no default")
	}
	if _, ok := v.(*List); ok {
		return asStringsList(v)
	}
	return asStrings(v)
}

func asStringsList(v Value) []Str {
	l := v.(*List)
	out := make([]Str, len(l.V))
	for i, e := range l.V {
		if st, ok := e.(*String); ok && len(st.V) == 1 {
			out[i] = st.V[0]
		} else {
			out[i] = Str{S: strings.Join(deparseValue(e), "")}
		}
	}
	return out
}

func deparseValue(v Value) []string { return syntaxDeparse(v) }

func intAttrs(extra ...Attr) []Attr { return extra }

func matchAttrs(lens []int32, useBytes bool) []Attr {
	if !useBytes {
		// as R, matching characters of non-ASCII text reports only match.length
		return []Attr{{Name: "match.length", Value: &Integer{V: lens}}}
	}
	return []Attr{
		{Name: "match.length", Value: &Integer{V: lens}},
		{Name: "index.type", Value: mkString("chars")},
		{Name: "useBytes", Value: mkLogical(useBytes)},
	}
}

// allASCII reports whether the pattern and every text are ASCII: R then
// matches bytes and says so in the useBytes attribute.
func allASCII(pat string, x []Str) bool {
	ascii := func(s string) bool {
		for i := 0; i < len(s); i++ {
			if s[i] >= 0x80 {
				return false
			}
		}
		return true
	}
	if !ascii(pat) {
		return false
	}
	for _, e := range x {
		if !e.NA && !ascii(e.S) {
			return false
		}
	}
	return true
}

func init() {
	regBuiltin("grepl", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "pattern", "x", "ignore.case", "perl", "fixed", "useBytes")
		pat, okp := s.patternOf(c, a[0])
		x := s.textOf(c, a[1])
		fl := s.reFlags(c, map[string]Value{"ignore.case": a[2], "perl": a[3], "fixed": a[4]})
		out := make([]int32, len(x))
		if !okp {
			for i := range out {
				out[i] = naInt
			}
			return &Logical{V: out}
		}
		m := s.compileRegexGrep(c, pat, x, fl.fixed, fl.icase, fl.perl)
		for i, e := range x {
			if !e.NA && m.matches(e.S) {
				out[i] = 1
			}
		}
		return &Logical{V: out}
	})
	regBuiltin("grep", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "pattern", "x", "ignore.case", "perl", "value", "fixed", "useBytes", "invert")
		pat, okp := s.patternOf(c, a[0])
		x := s.textOf(c, a[1])
		fl := s.reFlags(c, map[string]Value{"ignore.case": a[2], "perl": a[3], "value": a[4], "fixed": a[5], "invert": a[7]})
		var pos []int
		if okp {
			m := s.compileRegexGrep(c, pat, x, fl.fixed, fl.icase, fl.perl)
			for i, e := range x {
				hit := !e.NA && m.matches(e.S)
				if e.NA {
					continue
				}
				if hit != fl.invert {
					pos = append(pos, i)
				}
			}
		}
		if fl.value {
			out := make([]Str, len(pos))
			for i, p := range pos {
				out[i] = x[p]
			}
			res := &String{V: out}
			if nm := namesOf(a[1]); nm != nil {
				names := make([]Str, len(pos))
				for i, p := range pos {
					names[i] = nm.V[p]
				}
				res.Attr = namesAttr(names)
			}
			return res
		}
		out := make([]int32, len(pos))
		for i, p := range pos {
			out[i] = int32(p + 1)
		}
		res := &Integer{V: out}
		if nm := namesOf(a[1]); nm != nil {
			names := make([]Str, len(pos))
			for i, p := range pos {
				names[i] = nm.V[p]
			}
			res.Attr = namesAttr(names)
		}
		return res
	})
	subFn := func(global bool) func(s *Session, c *callInfo) Value {
		return func(s *Session, c *callInfo) Value {
			a := s.matchFormals(c, "pattern", "replacement", "x", "ignore.case", "perl", "fixed", "useBytes")
			if a[1] == nil {
				s.errorf(c.call, "argument \"replacement\" is missing, with no default")
			}
			pat, okp := s.patternOf(c, a[0])
			rs, ok := a[1].(*String)
			if !ok || len(rs.V) < 1 {
				rv := asStrings(a[1])
				if len(rv) < 1 {
					s.errorf(c.call, "invalid 'replacement' argument")
				}
				rs = &String{V: rv}
			}
			if len(rs.V) > 1 {
				s.warningf(c.call, "argument 'replacement' has length > 1 and only the first element will be used")
			}
			x := s.textOf(c, a[2])
			fl := s.reFlags(c, map[string]Value{"ignore.case": a[3], "perl": a[4], "fixed": a[5]})
			out := make([]Str, len(x))
			if !okp {
				for i := range out {
					out[i] = Str{NA: true}
				}
				return &String{V: out}
			}
			m := s.compileRegex(c, pat, x, fl.fixed, fl.icase, fl.perl)
			repl, replNA := rs.V[0].S, rs.V[0].NA
			for i, e := range x {
				if e.NA {
					out[i] = Str{NA: true}
					continue
				}
				str := e.S
				matches := m.allMatches(str)
				if !global && len(matches) > 1 {
					matches = matches[:1]
				}
				if len(matches) == 0 {
					out[i] = e
					continue
				}
				if replNA {
					out[i] = Str{NA: true}
					continue
				}
				var b strings.Builder
				last := 0
				for _, loc := range matches {
					b.WriteString(str[last:loc[0]])
					if fl.fixed {
						b.WriteString(repl)
					} else {
						b.WriteString(expandReplacement(repl, str, loc))
					}
					last = loc[1]
				}
				b.WriteString(str[last:])
				out[i] = Str{S: b.String()}
			}
			res := &String{V: out}
			if _, isStr := a[2].(*String); isStr {
				res.Attr = attrsOf(a[2])
			} else if nm := namesOf(a[2]); nm != nil {
				res.Attr = namesAttr(nm.V)
			}
			return res
		}
	}
	regBuiltin("sub", subFn(false))
	regBuiltin("gsub", subFn(true))
	regBuiltin("regexpr", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "pattern", "text", "ignore.case", "perl", "fixed", "useBytes")
		pat, okp := s.patternOf(c, a[0])
		x := s.textOf(c, a[1])
		fl := s.reFlags(c, map[string]Value{"ignore.case": a[2], "perl": a[3], "fixed": a[4]})
		pos := make([]int32, len(x))
		lens := make([]int32, len(x))
		var m *matcher
		if okp {
			m = s.compileRegexGrep(c, pat, x, fl.fixed, fl.icase, fl.perl)
		}
		for i, e := range x {
			if !okp || e.NA {
				pos[i], lens[i] = naInt, naInt
				continue
			}
			loc := m.find(e.S, 0)
			if loc == nil {
				pos[i], lens[i] = -1, -1
				continue
			}
			pos[i] = int32(charIndex(e.S, loc[0]) + 1)
			lens[i] = int32(utf8.RuneCountInString(e.S[loc[0]:loc[1]]))
		}
		if !allASCII(pat, x) {
			return &Integer{V: pos, Attr: []Attr{{Name: "match.length", Value: &Integer{V: lens}}}}
		}
		return &Integer{V: pos, Attr: matchAttrs(lens, true)}
	})
	regBuiltin("gregexpr", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "pattern", "text", "ignore.case", "perl", "fixed", "useBytes")
		pat, okp := s.patternOf(c, a[0])
		x := s.textOf(c, a[1])
		fl := s.reFlags(c, map[string]Value{"ignore.case": a[2], "perl": a[3], "fixed": a[4]})
		var m *matcher
		if okp {
			m = s.compileRegexGrep(c, pat, x, fl.fixed, fl.icase, fl.perl)
		}
		out := make([]Value, len(x))
		ascii := okp && allASCII(pat, x)
		for i, e := range x {
			if !okp || e.NA {
				out[i] = &Integer{V: []int32{naInt}, Attr: []Attr{{Name: "match.length", Value: &Integer{V: []int32{naInt}}}}}
				continue
			}
			var ms [][]int
			if e.S != "" { // R finds no match, empty or not, in an empty text
				ms = m.gregexprMatches(e.S)
			}
			if len(ms) == 0 {
				out[i] = &Integer{V: []int32{-1}, Attr: matchAttrs([]int32{-1}, ascii)}
				continue
			}
			pos := make([]int32, len(ms))
			lens := make([]int32, len(ms))
			for k, loc := range ms {
				pos[k] = int32(charIndex(e.S, loc[0]) + 1)
				lens[k] = int32(utf8.RuneCountInString(e.S[loc[0]:loc[1]]))
			}
			out[i] = &Integer{V: pos, Attr: matchAttrs(lens, ascii)}
		}
		return &List{V: out}
	})
	regBuiltin("regmatches", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "m", "invert")
		x := s.textOf(c, a[0])
		invert := a[2] != nil && s.truthy(c, a[2], "invert")
		runes := func(str string, start, length int) string {
			return substrRunes(str, start, start+length-1)
		}
		switch m := a[1].(type) {
		case *Integer:
			ml, _ := getAttr(m, "match.length").(*Integer)
			if ml == nil || len(m.V) != len(x) {
				s.errorf(c.call, "'x' and 'm' must have the same length")
			}
			if !invert {
				var out []Str
				var names []Str
				nm := namesOf(a[0])
				for i, e := range x {
					if m.V[i] == naInt || m.V[i] < 0 || e.NA {
						continue
					}
					out = append(out, Str{S: runes(e.S, int(m.V[i]), int(ml.V[i]))})
					if nm != nil {
						names = append(names, nm.V[i])
					}
				}
				res := &String{V: out}
				if res.V == nil {
					res.V = []Str{}
				}
				if nm != nil {
					res.Attr = namesAttr(names)
				}
				return res
			}
			out := make([]Value, len(x))
			for i, e := range x {
				if e.NA {
					out[i] = &String{V: []Str{{NA: true}}}
					continue
				}
				if m.V[i] < 0 {
					out[i] = &String{V: []Str{e}}
					continue
				}
				st, ln := int(m.V[i]), int(ml.V[i])
				out[i] = &String{V: []Str{{S: substrRunes(e.S, 1, st-1)}, {S: substrRunes(e.S, st+ln, 1<<30)}}}
			}
			return listOfNamed(out, namesOf(a[0]))
		case *List:
			if len(m.V) != len(x) {
				s.errorf(c.call, "'x' and 'm' must have the same length")
			}
			out := make([]Value, len(x))
			for i, e := range x {
				mi, _ := m.V[i].(*Integer)
				if mi == nil {
					s.errorf(c.call, "invalid 'm' argument")
				}
				ml, _ := getAttr(mi, "match.length").(*Integer)
				if e.NA {
					out[i] = &String{V: []Str{}}
					if invert {
						out[i] = &String{V: []Str{{NA: true}}}
					}
					continue
				}
				if !invert {
					var parts []Str
					for k, p := range mi.V {
						if p < 0 || p == naInt {
							continue
						}
						parts = append(parts, Str{S: runes(e.S, int(p), int(ml.V[k]))})
					}
					if parts == nil {
						parts = []Str{}
					}
					out[i] = &String{V: parts}
					continue
				}
				var parts []Str
				prev := 1
				for k, p := range mi.V {
					if p < 0 || p == naInt {
						continue
					}
					parts = append(parts, Str{S: substrRunes(e.S, prev, int(p)-1)})
					prev = int(p) + int(ml.V[k])
				}
				parts = append(parts, Str{S: substrRunes(e.S, prev, 1<<30)})
				out[i] = &String{V: parts}
			}
			return listOfNamed(out, namesOf(a[0]))
		}
		s.errorf(c.call, "invalid 'm' argument")
		return nil
	})
	regBuiltin("startsWith", func(s *Session, c *callInfo) Value { return s.affix(c, true) })
	regBuiltin("endsWith", func(s *Session, c *callInfo) Value { return s.affix(c, false) })
}

func (s *Session) affix(c *callInfo, start bool) Value {
	a := s.matchFormals(c, "x", "prefix")
	x, ok1 := a[0].(*String)
	p, ok2 := a[1].(*String)
	if !ok1 || !ok2 {
		s.errorf(c.call, "non-character object(s)")
	}
	n := len(x.V)
	if len(p.V) > n && n > 0 {
		n = len(p.V)
	}
	if len(x.V) == 0 || len(p.V) == 0 {
		return &Logical{V: []int32{}}
	}
	out := make([]int32, n)
	for i := range out {
		xe, pe := x.V[i%len(x.V)], p.V[i%len(p.V)]
		switch {
		case xe.NA || pe.NA:
			out[i] = naInt
		case start && strings.HasPrefix(xe.S, pe.S), !start && strings.HasSuffix(xe.S, pe.S):
			out[i] = 1
		}
	}
	return &Logical{V: out}
}

func listOfNamed(vals []Value, names *String) *List {
	l := &List{V: vals}
	if l.V == nil {
		l.V = []Value{}
	}
	if names != nil {
		l.Attr = namesAttr(names.V)
	}
	return l
}

var _ = strconv.Itoa
