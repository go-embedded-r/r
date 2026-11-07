package eval

import (
	"math"
	"path"
	"strconv"
	"strings"
)

// Slice 6: more of the string functions of base.

// insertBigMark puts mark between groups of three digits in the integer
// part of a formatted number.
func insertBigMark(s, mark string) string {
	if mark == "" {
		return s
	}
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '-' || s[start] == '+') {
		start++
	}
	end := start
	for end < len(s) && s[end] >= '0' && s[end] <= '9' {
		end++
	}
	digits := s[start:end]
	if len(digits) <= 3 {
		return s
	}
	var b strings.Builder
	for i, ch := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			b.WriteString(mark)
		}
		b.WriteRune(ch)
	}
	return s[:start] + b.String() + s[end:]
}

func init() {
	regBuiltin("strrep", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "times")
		if a[0] == nil || a[1] == nil {
			s.errorf(c.call, "argument \"%s\" is missing, with no default", map[bool]string{true: "x", false: "times"}[a[0] == nil])
		}
		x := s.charVec(c, a[0], "strrep")
		t := asDoublesOrNil(a[1])
		if len(x) == 0 || len(t) == 0 {
			return &String{V: []Str{}}
		}
		n := max(len(x), len(t))
		out := make([]Str, n)
		for i := range out {
			e, k := x[i%len(x)], t[i%len(t)]
			switch {
			case e.NA || math.IsNaN(k):
				out[i] = Str{NA: true}
			case k < 0:
				s.errorf(c.call, "invalid 'times' value")
			default:
				if float64(len(e.S))*k > float64(s.opts.MaxVector) {
					s.errorf(c.call, "vector size too large")
				}
				out[i] = Str{S: strings.Repeat(e.S, int(k))}
			}
		}
		res := &String{V: out}
		if n == len(x) {
			res.Attr = keepNamesOnly(a[0])
		}
		return res
	})
	regBuiltin("substr<-", func(s *Session, c *callInfo) Value {
		if len(c.args) != 4 {
			s.errorf(c.call, "4 arguments are required")
		}
		x, ok := c.args[0].(*String)
		if !ok {
			call := c.tmpCall
			if call == nil {
				call = c.call
			}
			s.errorf(call, "replacing substrings in a non-character object")
		}
		st, sp := s.intsOf(c, c.args[1], "start"), s.intsOf(c, c.args[2], "stop")
		val := s.charVec(c, c.args[3], "substr<-")
		if len(x.V) > 0 && (len(st) == 0 || len(sp) == 0 || len(val) == 0) {
			s.errorf(c.call, "invalid substring arguments")
		}
		out := make([]Str, len(x.V))
		for i, e := range x.V {
			b, f, v := st[i%len(st)], sp[i%len(sp)], val[i%len(val)]
			if e.NA || v.NA || b == naInt || f == naInt {
				out[i] = Str{NA: true}
				continue
			}
			rs, vs := []rune(e.S), []rune(v.S)
			start, stop := int(b), int(f)
			if start < 1 {
				start = 1
			}
			if stop > len(rs) {
				stop = len(rs)
			}
			if start > stop {
				out[i] = e
				continue
			}
			n := min(stop-start+1, len(vs))
			copy(rs[start-1:start-1+n], vs[:n])
			out[i] = Str{S: string(rs)}
		}
		return &String{V: out, Attr: x.Attr}
	})
	regBuiltin("strtoi", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "base")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		base := 10
		if a[1] != nil {
			base = s.intArg(c, a[1], "base")
			if base != 0 && (base < 2 || base > 36) {
				s.errorf(c.call, "invalid '%s' argument", "base")
			}
		}
		x := s.charVec(c, a[0], "strtoi")
		out := make([]int32, len(x))
		for i, e := range x {
			out[i] = naInt
			if e.NA {
				continue
			}
			str := strings.TrimLeft(e.S, " \t\n")
			v, err := strconv.ParseInt(str, base, 64)
			if err != nil || v > math.MaxInt32 || v < -math.MaxInt32 {
				if base == 16 || base == 0 {
					low := strings.ToLower(str)
					if strings.HasPrefix(low, "0x") {
						if v2, err2 := strconv.ParseInt(low[2:], 16, 64); err2 == nil && v2 <= math.MaxInt32 {
							out[i] = int32(v2)
						}
					}
				}
				continue
			}
			out[i] = int32(v)
		}
		return withAttr(&Integer{V: out}, keepNamesOnly(a[0]))
	})
	regBuiltin("make.unique", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "names", "sep")
		names, ok := a[0].(*String)
		if !ok {
			s.errorf(c.call, "'names' must be a character vector")
		}
		sep := "."
		if a[1] != nil {
			sep = s.stringArg(c, a[1], "sep")
		}
		out := make([]Str, len(names.V))
		copy(out, names.V)
		seen := map[string]bool{}
		for _, e := range out {
			seen[e.S] = true
		}
		first := map[string]bool{}
		next := map[string]int{}
		for i, e := range names.V {
			if !first[e.S] {
				first[e.S] = true
				continue
			}
			k := next[e.S]
			for {
				k++
				cand := e.S + sep + strconv.Itoa(k)
				if !seen[cand] {
					out[i] = Str{S: cand}
					seen[cand] = true
					break
				}
			}
			next[e.S] = k
		}
		return &String{V: out}
	})
	regBuiltin("basename", func(s *Session, c *callInfo) Value {
		x, ok := s.nargs1(c, "basename").(*String)
		if !ok {
			s.errorf(c.call, "a character vector argument expected")
		}
		out := make([]Str, len(x.V))
		for i, e := range x.V {
			if e.NA {
				out[i] = e
				continue
			}
			t := strings.TrimRight(e.S, "/")
			if k := strings.LastIndexByte(t, '/'); k >= 0 {
				t = t[k+1:]
			}
			out[i] = Str{S: t}
		}
		return &String{V: out}
	})
	regBuiltin("dirname", func(s *Session, c *callInfo) Value {
		x, ok := s.nargs1(c, "dirname").(*String)
		if !ok {
			s.errorf(c.call, "a character vector argument expected")
		}
		out := make([]Str, len(x.V))
		for i, e := range x.V {
			if e.NA {
				out[i] = e
				continue
			}
			t := strings.TrimRight(e.S, "/")
			switch {
			case t == "" && e.S != "":
				out[i] = Str{S: "/"}
			case t == "":
				out[i] = Str{S: ""}
			default:
				k := strings.LastIndexByte(t, '/')
				switch {
				case k < 0:
					out[i] = Str{S: "."}
				default:
					d := strings.TrimRight(t[:k], "/")
					if d == "" {
						d = "/"
					}
					out[i] = Str{S: d}
				}
			}
		}
		return &String{V: out}
	})
	regBuiltin("formatC", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "width", "digits", "format", "flag", "mode", "big.mark")
		x := a[0]
		if x == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		switch x.(type) {
		case *Integer, *Double, *String:
		default:
			s.errorf(c.call, "unsupported type ")
		}
		width := 0
		if a[1] != nil {
			width = s.intArg(c, a[1], "width")
		}
		digits := -1
		digitsGiven := false
		if a[2] != nil {
			if _, isNull := a[2].(nullT); !isNull {
				digits, digitsGiven = s.intArg(c, a[2], "digits"), true
			}
		}
		flag := ""
		if a[4] != nil {
			flag = s.stringArg(c, a[4], "flag")
		}
		if width < 0 {
			width = -width
			flag += "-"
		}
		mark := ""
		if a[6] != nil {
			mark = s.stringArg(c, a[6], "big.mark")
		}
		format := ""
		if a[3] != nil {
			format = s.stringArg(c, a[3], "format")
		} else {
			switch x.(type) {
			case *String:
				format = "s"
			case *Double:
				format = "g"
			default:
				format = "d"
			}
		}
		if width > maxFieldWidth {
			s.errorf(c.call, "invalid 'width' argument")
		}
		n := Length(x)
		out := make([]Str, n)
		for i := 0; i < n; i++ {
			var body string
			switch format {
			case "s":
				e := asStrings(x)[i]
				body = e.S
				if e.NA {
					body = "NA"
				}
			case "d":
				d := asDoubles(x)[i]
				switch {
				case math.IsNaN(d):
					body = "NA"
				default:
					body = strconv.FormatFloat(math.Round(d), 'f', 0, 64)
				}
			case "f", "e", "E", "g", "G":
				d := asDoubles(x)[i]
				switch {
				case isNAReal(d) || (kindOf(x) != kDouble && math.IsNaN(d)):
					body = "NA"
				case math.IsNaN(d):
					body = "NaN"
				case math.IsInf(d, 1):
					body = "Inf"
				case math.IsInf(d, -1):
					body = "-Inf"
				default:
					p := digits
					if !digitsGiven {
						p = 4 // formatC's default for real numbers
					}
					verb := format
					switch {
					case p < 0 && (verb == "f" || verb == "e" || verb == "E"):
						p = 6
					case p < 0:
						p = -1
					}
					if verb == "g" || verb == "G" {
						body = formatG(d, p, verb == "G", strings.Contains(flag, "#"))
					} else {
						body = strconv.FormatFloat(d, verb[0], p, 64)
						if verb == "e" || verb == "E" {
							body = fixExp(body)
						}
					}
				}
			default:
				s.errorf(c.call, "'format' must be one of {\"f\",\"e\",\"E\",\"g\",\"G\", \"fg\", \"s\"}")
			}
			if strings.Contains(flag, "+") && len(body) > 0 && body[0] != '-' && body != "NA" {
				body = "+" + body
			} else if strings.Contains(flag, " ") && len(body) > 0 && body[0] != '-' && body != "NA" {
				body = " " + body
			}
			body = insertBigMark(body, mark)
			if w := displayWidth(body); w < width {
				switch {
				case strings.Contains(flag, "-"):
					body += strings.Repeat(" ", width-w)
				case strings.Contains(flag, "0") && format != "s" && body != "NA":
					sign := ""
					if len(body) > 0 && (body[0] == '-' || body[0] == '+' || body[0] == ' ') {
						sign, body = body[:1], body[1:]
					}
					body = sign + strings.Repeat("0", width-w) + body
				default:
					body = strings.Repeat(" ", width-w) + body
				}
			}
			out[i] = Str{S: body}
		}
		return withAttr(&String{V: out}, keepShapeAttrs(x))
	})
	regBuiltin("prettyNum", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "big.mark")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		mark := ""
		if a[1] != nil {
			mark = s.stringArg(c, a[1], "big.mark")
		}
		var x []Str
		if st, ok := a[0].(*String); ok {
			x = st.V
		} else {
			for _, e := range asStrings(a[0]) {
				if e.NA {
					e = Str{S: "NA"}
				}
				x = append(x, e)
			}
		}
		out := make([]Str, len(x))
		for i, e := range x {
			if e.NA {
				out[i] = e
				continue
			}
			out[i] = Str{S: insertBigMark(e.S, mark)}
		}
		return withAttr(&String{V: out}, keepShapeAttrs(a[0]))
	})
}

// fixExp makes Go's exponent ("e+06") what C prints ("e+06"): two digits at least.
func fixExp(s string) string {
	i := strings.LastIndexAny(s, "eE")
	if i < 0 || i+2 >= len(s) {
		return s
	}
	sign, digits := s[i+1], s[i+2:]
	if len(digits) < 2 {
		digits = "0" + digits
	}
	return s[:i+1] + string(sign) + digits
}

// formatG is C's %g with precision p (default 6; -1 asks for it).
func formatG(d float64, p int, upper, alt bool) string {
	if p < 0 {
		p = 6
	}
	if p == 0 {
		p = 1
	}
	if d == 0 {
		return "0"
	}
	exp := int(math.Floor(math.Log10(math.Abs(d))))
	// rounding may carry into the next power of ten
	if r, _ := strconv.ParseFloat(strconv.FormatFloat(d, 'e', p-1, 64), 64); r != 0 {
		exp = int(math.Floor(math.Log10(math.Abs(r))))
	}
	var s string
	if exp < -4 || exp >= p {
		s = strconv.FormatFloat(d, 'e', p-1, 64)
		if !alt {
			m, e, _ := strings.Cut(s, "e")
			if strings.Contains(m, ".") {
				m = strings.TrimRight(strings.TrimRight(m, "0"), ".")
			}
			s = m + "e" + e
		}
		s = fixExp(s)
	} else {
		s = strconv.FormatFloat(d, 'f', max(p-1-exp, 0), 64)
		if !alt && strings.Contains(s, ".") {
			s = strings.TrimRight(strings.TrimRight(s, "0"), ".")
		}
	}
	if upper {
		s = strings.ToUpper(s)
	}
	return s
}

var _ = path.Base
