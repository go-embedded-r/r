package eval

import (
	"math"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Support for data frames, written from the help pages of data.frame,
// row.names and .row_names_info and from experiments with the R 4.6.2 oracle:
// the compact form of the row.names attribute (integer c(NA, -n)), which is
// what attr(x, "row.names") hides from the user.

// compactRowNames is what R stores when row.names is set to the integers 1..n.
func compactRowNames(v Value) Value {
	iv, ok := v.(*Integer)
	if !ok || len(iv.V) == 0 {
		return v
	}
	if len(iv.V) == 2 && iv.V[0] == naInt && iv.V[1] != naInt {
		return v // already compact
	}
	for i, e := range iv.V {
		if e != int32(i+1) {
			return v
		}
	}
	return &Integer{V: []int32{naInt, -int32(len(iv.V))}}
}

// rowNamesCount is the number of rows a row.names attribute describes.
func rowNamesCount(v Value) int {
	if iv, ok := v.(*Integer); ok && len(iv.V) == 2 && iv.V[0] == naInt && iv.V[1] != naInt {
		if iv.V[1] < 0 {
			return int(-iv.V[1])
		}
		return int(iv.V[1])
	}
	return Length(v)
}

// expandRowNames is attr(x, "row.names"): the compact form is 1:n.
func expandRowNames(v Value) Value {
	if iv, ok := v.(*Integer); ok && len(iv.V) == 2 && iv.V[0] == naInt && iv.V[1] != naInt {
		n := rowNamesCount(v)
		out := make([]int32, n)
		for i := range out {
			out[i] = int32(i + 1)
		}
		return &Integer{V: out}
	}
	return v
}

var reservedWords = map[string]bool{"if": true, "else": true, "repeat": true, "while": true, "function": true,
	"for": true, "next": true, "break": true, "TRUE": true, "FALSE": true, "NULL": true, "Inf": true, "NaN": true,
	"NA": true, "NA_integer_": true, "NA_real_": true, "NA_character_": true, "NA_complex_": true, "in": true}

// validName is one element of make.names(allow_ = TRUE).
func validName(x Str) string {
	if x.NA {
		return "NA"
	}
	var b strings.Builder
	for _, r := range x.S {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_' {
			b.WriteRune(r)
		} else {
			b.WriteByte('.')
		}
	}
	n := b.String()
	if n == "" {
		return "X"
	}
	first, _ := utf8.DecodeRuneInString(n)
	switch {
	case unicode.IsLetter(first):
	case first == '.':
		if len(n) > 1 && n[1] >= '0' && n[1] <= '9' {
			n = "X" + n
		}
	default:
		n = "X" + n
	}
	if reservedWords[n] {
		n += "."
	}
	return n
}

// factorAssign is `[<-.factor` and `[[<-.factor`: the value is matched to the
// levels (a value that is none of them is NA, with a warning) and stored as
// codes. A Go native: a first draft in R converged on R's own text.
func factorAssign(name string) func(s *Session, c *callInfo) Value {
	return func(s *Session, c *callInfo) Value {
		n := len(c.args)
		if n < 2 || c.names[n-1] != "value" {
			s.errorf(c.call, "argument \"value\" is missing, with no default")
		}
		x, value := c.args[0], c.args[n-1]
		lev, _ := getAttr(x, "levels").(*String)
		if lev == nil {
			lev = &String{}
		}
		vs := s.factorStrings(value)
		codes := make([]int32, len(vs))
		bad := false
		for i, e := range vs {
			codes[i] = naInt
			if e.NA {
				continue
			}
			found := false
			for k, l := range lev.V {
				if !l.NA && l.S == e.S {
					codes[i], found = int32(k+1), true
					break
				}
			}
			bad = bad || !found
		}
		if bad {
			s.warningf(c.call, "invalid factor level, NA generated")
		}
		var keep []Attr
		var cls, lv Value
		for _, a := range attrsOf(x) {
			switch a.Name {
			case "class":
				cls = a.Value
			case "levels":
				lv = a.Value
			default:
				keep = append(keep, a)
			}
		}
		args := append([]Value{withAttr(x, keep)}, c.args[1:n-1]...)
		args = append(args, &Integer{V: codes})
		res := builtins[name].fn(s, &callInfo{call: c.call, env: c.env, args: args, names: c.names})
		return withAttr(res, append(append([]Attr(nil), attrsOf(res)...), Attr{Name: "levels", Value: lv}, Attr{Name: "class", Value: cls}))
	}
}

func init() {
	regBuiltin("addNA", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "ifany")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		fac := func(args []Value, names []string) Value {
			return builtins["factor"].fn(s, &callInfo{call: c.call, env: c.env, args: args, names: names})
		}
		x := a[0]
		if !inheritsClass(x, "factor") {
			x = fac([]Value{x}, []string{""})
		}
		lev, _ := getAttr(x, "levels").(*String)
		if lev == nil {
			lev = &String{}
		}
		hasNA := false
		for _, l := range lev.V {
			hasNA = hasNA || l.NA
		}
		if a[1] != nil && s.truthy(c, a[1], "ifany") {
			anyNA := false
			for _, v := range x.(*Integer).V {
				anyNA = anyNA || v == naInt
			}
			if !anyNA {
				return x
			}
		}
		nl := append([]Str(nil), lev.V...)
		if !hasNA {
			nl = append(nl, Str{NA: true})
		}
		return fac([]Value{x, &String{V: nl}, Null}, []string{"", "levels", "exclude"})
	})
	regBuiltin("[<-.factor", factorAssign("[<-"))
	regBuiltin("[[<-.factor", factorAssign("[[<-"))
	regBuiltin("make.names", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "names", "unique", "allow_")
		if a[0] == nil {
			s.errorf(c.call, "argument \"names\" is missing, with no default")
		}
		x := a[0]
		if _, ok := x.(*String); !ok {
			x = coerceTo(x, kString)
		}
		xs := x.(*String)
		out := make([]Str, len(xs.V))
		allowUnderscore := a[2] == nil || s.truthy(c, a[2], "allow_")
		for i, e := range xs.V {
			n := validName(e)
			if !allowUnderscore {
				n = strings.ReplaceAll(n, "_", ".")
			}
			out[i] = Str{S: n}
		}
		var res Value = &String{V: out}
		if a[1] != nil && s.truthy(c, a[1], "unique") {
			res = builtins["make.unique"].fn(s, &callInfo{call: c.call, env: c.env, args: []Value{res}, names: []string{""}})
		}
		return res
	})
	regBuiltin(".bincode", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "breaks", "right", "include.lowest")
		x, br := asDoubles(coerceTo(a[0], kDouble)), asDoubles(coerceTo(a[1], kDouble))
		right := a[2] == nil || s.truthy(c, a[2], "right")
		lowest := a[3] != nil && s.truthy(c, a[3], "include.lowest")
		nb := len(br)
		for i := 1; i < nb; i++ {
			if br[i-1] > br[i] {
				s.errorf(c.call, "'breaks' is not sorted")
			}
		}
		out := make([]int32, len(x))
		for i, v := range x {
			out[i] = naInt
			if math.IsNaN(v) || nb < 2 {
				continue
			}
			lo, hi := 0, nb-1
			if v < br[lo] || br[hi] < v || (v == br[map[bool]int{true: lo, false: hi}[right]] && !lowest) {
				continue
			}
			for hi-lo >= 2 {
				mid := (hi + lo) / 2
				if v > br[mid] || (!right && v == br[mid]) {
					lo = mid
				} else {
					hi = mid
				}
			}
			out[i] = int32(lo + 1)
		}
		return &Integer{V: out}
	})
	regBuiltin(".set_row_names", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "n")
		if a[0] == nil {
			s.errorf(c.call, "argument \"n\" is missing, with no default")
		}
		d := asDoubles(coerceTo(a[0], kDouble))
		if len(d) != 1 || d[0] != d[0] || d[0] <= 0 {
			return &Integer{V: []int32{}}
		}
		if d[0] > 2147483647 {
			s.errorf(c.call, "invalid 'n' argument")
		}
		return &Integer{V: []int32{naInt, -int32(d[0])}}
	})
	regBuiltin(".row_names_info", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "type")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		typ := int32(1)
		if a[1] != nil {
			if d := asDoubles(coerceTo(a[1], kDouble)); len(d) == 1 {
				typ = int32(d[0])
			}
		}
		rn := getAttr(a[0], "row.names")
		if rn == nil {
			if typ == 0 {
				return Null
			}
			return mkInt(0)
		}
		if typ == 0 {
			return rn
		}
		n := int32(rowNamesCount(rn))
		if typ == 2 {
			return mkInt(n)
		}
		if iv, ok := rn.(*Integer); ok && len(iv.V) == 2 && iv.V[0] == naInt && iv.V[1] != naInt {
			return mkInt(iv.V[1]) // negative: automatic row names
		}
		return mkInt(n)
	})
}
