package eval

import (
	"math"
	"strconv"
	"strings"
)

// str (slice 17): the compact display of a structure, written from the help
// page of utils::str and from what the oracle prints for vectors, factors,
// dates, matrices and attributes. The oracle refuses to print what str shows of
// a list, a data frame or a named vector (its guard sees a base closure), so
// those layouts follow the help page and are checked by eval/str17_test.go,
// not by goldens.

type strOpts struct {
	vecLen   float64 // vec.len, 4 by default
	digits   int     // digits.d, 3 by default
	giveAttr bool
	giveHead bool
	maxLevel int // -1: no limit
	indent   string
	width    int
}

func (s *Session) strFormatNum(c *callInfo, v Value, o strOpts) []string {
	f := func(drop bool) []string {
		args := []Value{v, mkInt(int32(o.digits)), mkLogical(true), mkLogical(drop)}
		r := s.callNamed(c, "format", args, []string{"", "digits", "trim", "drop0trailing"})
		out := make([]string, 0, Length(r))
		for _, e := range asStrings(r) {
			out = append(out, e.S)
		}
		return out
	}
	items := f(true)
	for _, it := range items {
		if strings.Contains(it, "e") {
			return f(false)
		}
	}
	return items
}

func strQuote(s string) string {
	return "\"" + strings.NewReplacer("\\", "\\\\", "\"", "\\\"").Replace(s) + "\""
}

// strVector writes the header and the first elements of an atomic vector.
func (s *Session) strVector(c *callInfo, v Value, o strOpts) string {
	n := Length(v)
	d := dimsOf(v)
	var typ string
	var shown int
	vl := o.vecLen
	switch kindOf(v) {
	case kInteger:
		typ, shown = "int", int(math.RoundToEven(2.5*vl))
	case kDouble:
		typ, shown = "num", int(math.RoundToEven(2.5*vl))
	case kLogical:
		typ, shown = "logi", int(math.RoundToEven(1.5*vl))
	case kString:
		typ = "chr"
	case kComplex:
		typ, shown = "cplx", int(math.RoundToEven(1.25*vl))
	case kRaw:
		typ, shown = "raw", int(math.RoundToEven(2.5*vl))
	default:
		typ, shown = typeName(v), 5
	}
	var head string
	switch {
	case n == 0 && len(d) == 0:
		return " " + typ + "(0) "
	case len(d) > 0:
		parts := make([]string, len(d))
		for i, e := range d {
			parts[i] = "1:" + strconv.Itoa(e)
			if e == 1 {
				parts[i] = "1" // an extent of one is shown as 1
			}
		}
		head = " " + typ + " [" + strings.Join(parts, ", ") + "]"
	case n == 1:
		head = " " + typ
	default:
		head = " " + typ + " [1:" + strconv.Itoa(n) + "]"
	}
	if named := namesOf(v); named != nil && len(d) == 0 {
		head = " Named" + head
	}
	if cl := classAttr(v); cl != nil && !inheritsClass(v, "factor") {
		quoted := make([]string, len(cl))
		for i, e := range cl {
			quoted[i] = "'" + e + "'"
		}
		head = " " + strings.Join(quoted, ", ") + head
	}
	var items []string
	switch kindOf(v) {
	case kString:
		sv := asStrings(v)
		width, k := 0, 0
		for k < len(sv) {
			it := "NA"
			if !sv[k].NA {
				it = strQuote(sv[k].S)
			}
			add := len(it)
			if k > 0 {
				add++
			}
			if k > 0 && width+add > o.width-16 {
				break
			}
			width += add
			items = append(items, it)
			k++
		}
		shown = k
		if o.vecLen != 4 {
			shown = int(math.RoundToEven(1.25 * vl))
			if shown < len(items) {
				items = items[:shown]
			}
		}
	case kRaw:
		for _, e := range asStrings(v) {
			items = append(items, e.S)
		}
	case kDouble:
		cut := min(n, shown)
		items = s.strFormatNum(c, &Double{V: v.(*Double).V[:cut]}, o)
	case kComplex:
		cut := min(n, shown)
		for _, e := range asStrings(&Complex{V: v.(*Complex).V[:cut]}) {
			items = append(items, e.S)
		}
	default:
		for _, e := range asStrings(v) {
			if e.NA {
				items = append(items, "NA")
			} else {
				items = append(items, e.S)
			}
		}
	}
	if kindOf(v) != kString && len(items) > shown {
		items = items[:shown]
	}
	more := ""
	if n > len(items) {
		more = " ..."
	}
	if !o.giveHead {
		return strings.Join(items, " ") + more
	}
	return head + " " + strings.Join(items, " ") + more
}

func (s *Session) strFactor(v Value, o strOpts) string {
	lev, _ := getAttr(v, "levels").(*String)
	nl := 0
	if lev != nil {
		nl = len(lev.V)
	}
	ordered := inheritsClass(v, "ordered")
	cum, show := 0, nl
	for i := 0; i < nl; i++ {
		cum += 3 + len(strQuote(lev.V[i].S)) - 2
		if nl > 1 && cum > 13 {
			show = i + 1
			break
		}
	}
	sep := ","
	if ordered {
		sep = "<"
	}
	var ls []string
	for i := 0; i < show; i++ {
		ls = append(ls, strQuote(lev.V[i].S))
	}
	levels := strings.Join(ls, sep)
	if show < nl {
		levels += sep + ".."
	}
	kind := "Factor"
	if ordered {
		kind = "Ord.factor"
	}
	word := "levels"
	if nl == 1 {
		word = "level"
	}
	codes := asInts(v)
	shownCodes := int(math.RoundToEven(2.5 * o.vecLen))
	var cs []string
	for i, code := range codes {
		if i >= shownCodes {
			break
		}
		if code == naInt {
			cs = append(cs, "NA")
		} else {
			cs = append(cs, strconv.Itoa(int(code)))
		}
	}
	more := ""
	if len(codes) > shownCodes {
		more = " ..."
	}
	return " " + kind + " w/ " + strconv.Itoa(nl) + " " + word + " " + levels + ": " + strings.Join(cs, " ") + more
}

func (s *Session) strDateTime(c *callInfo, v Value, kind string, o strOpts) string {
	n := Length(v)
	shown := min(n, 40)
	head := kind + "[1:" + strconv.Itoa(n) + "], format: "
	var items []string
	if shown > 0 {
		cut := s.callNamed(c, "[", []Value{v, &Integer{V: seqInts(shown)}}, []string{"", ""})
		width := 0
		for k, e := range asStrings(s.callNamed(c, "format", []Value{cut}, []string{""})) {
			it := strQuote(e.S)
			add := len(it)
			if k > 0 {
				add++
			}
			if k > 0 && width+add > o.width-16 {
				break
			}
			width += add
			items = append(items, it)
		}
		shown = len(items)
	}
	more := ""
	if n > shown {
		more = " ..."
	}
	return " " + head + strings.Join(items, " ") + more
}

func seqInts(n int) []int32 {
	out := make([]int32, n)
	for i := range out {
		out[i] = int32(i + 1)
	}
	return out
}

// strAttrs writes the attributes str shows below an object.
func (s *Session) strAttrs(c *callInfo, v Value, o strOpts, out *strings.Builder, skip map[string]bool) {
	if !o.giveAttr {
		return
	}
	for _, a := range attrsOf(v) {
		if skip[a.Name] {
			continue
		}
		out.WriteString(o.indent + "- attr(*, \"" + a.Name + "\")=")
		s.strValue(c, a.Value, strOpts{vecLen: o.vecLen, digits: o.digits, giveAttr: o.giveAttr, giveHead: true, maxLevel: o.maxLevel, indent: o.indent + " ..", width: o.width}, out, 1)
	}
}

// strValue writes str of one object; level is the nesting level.
func (s *Session) strValue(c *callInfo, v Value, o strOpts, out *strings.Builder, level int) {
	s.descend()
	defer s.ascend()
	skip := map[string]bool{"dim": true, "class": true, "levels": true}
	switch x := v.(type) {
	case nullT:
		out.WriteString(" NULL\n")
	case *Symbol:
		out.WriteString(" symbol " + x.Name + "\n")
	case *Lang:
		out.WriteString(" language " + strings.Join(s.syntaxDeparse(x), " ") + "\n")
	case *Closure, *Builtin:
		out.WriteString("function ()  \n")
	case *Env:
		out.WriteString("<environment>\n")
	case *List:
		if x.Expr {
			out.WriteString("  expression(" + strconv.Itoa(len(x.V)) + ")\n")
			return
		}
		if inheritsClass(x, "data.frame") {
			nrow := rowNamesCount(getAttr(x, "row.names"))
			word := "variables"
			if len(x.V) == 1 {
				word = "variable"
			}
			out.WriteString("'data.frame':\t" + strconv.Itoa(nrow) + " obs. of  " + strconv.Itoa(len(x.V)) + " " + word + ":\n")
			s.strListBody(c, x, o, out, level, true)
			return
		}
		if len(x.V) == 0 {
			out.WriteString(" list()\n")
			return
		}
		out.WriteString("List of " + strconv.Itoa(len(x.V)) + "\n")
		s.strListBody(c, x, o, out, level, false)
		skip["names"] = true
		s.strAttrs(c, v, o, out, skip)
	default:
		switch {
		case inheritsClass(v, "factor") && kindOf(v) == kInteger:
			out.WriteString(s.strFactor(v, o) + "\n")
		case inheritsClass(v, "Date") && kindOf(v) == kDouble:
			out.WriteString(s.strDateTime(c, v, "Date", o) + "\n")
		case inheritsClass(v, "POSIXct") && kindOf(v) == kDouble:
			out.WriteString(s.strDateTime(c, v, "POSIXct", o) + "\n")
		default:
			out.WriteString(s.strVector(c, v, o) + "\n")
			skip["names"] = namesOf(v) == nil
			s.strAttrs(c, v, o, out, skip)
		}
	}
}

func (s *Session) strListBody(c *callInfo, x *List, o strOpts, out *strings.Builder, level int, frame bool) {
	if o.maxLevel >= 0 && level > o.maxLevel {
		return
	}
	nm := namesOf(x)
	width := 0
	names := make([]string, len(x.V))
	for i := range x.V {
		if nm != nil {
			names[i] = nm.V[i].S
		}
		if len(names[i]) > width {
			width = len(names[i])
		}
	}
	for i, e := range x.V {
		pad := names[i] + strings.Repeat(" ", width-len(names[i]))
		out.WriteString(o.indent + "$ " + pad + ":")
		sub := o
		sub.indent = o.indent + " .."
		if frame {
			// a column of a data frame is written without its length
			var b strings.Builder
			s.strValue(c, e, sub, &b, level+1)
			text := b.String()
			if n := Length(e); n > 1 && len(dimsOf(e)) == 0 && !inheritsClass(e, "factor") && !inheritsClass(e, "Date") {
				text = strings.Replace(text, " [1:"+strconv.Itoa(n)+"]", " ", 1)
			}
			out.WriteString(text)
			continue
		}
		s.strValue(c, e, sub, out, level+1)
	}
}

func init() {
	regBuiltin("..strDefault", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "object", "max.level", "vec.len", "digits.d", "nchar.max", "give.attr", "give.head", "indent.str")
		if a[0] == nil {
			s.errorf(c.call, "argument \"object\" is missing, with no default")
		}
		o := strOpts{vecLen: 4, digits: 3, giveAttr: true, giveHead: true, maxLevel: -1, width: 80, indent: " "}
		if w, ok := s.options["width"]; ok {
			if d := asDoubles(w); len(d) == 1 && d[0] > 20 {
				o.width = int(d[0])
			}
		}
		if a[1] != nil {
			if d := asDoubles(a[1]); len(d) == 1 && d[0] == d[0] {
				o.maxLevel = int(d[0])
			}
		}
		if a[2] != nil {
			if d := asDoubles(a[2]); len(d) == 1 && d[0] >= 0 {
				o.vecLen = d[0]
			}
		}
		if a[3] != nil {
			o.digits = s.intArg(c, a[3], "digits.d")
		}
		if a[5] != nil {
			o.giveAttr = s.truthy(c, a[5], "give.attr")
		}
		if a[6] != nil {
			o.giveHead = s.truthy(c, a[6], "give.head")
		}
		if a[7] != nil {
			o.indent = s.stringArg(c, a[7], "indent.str")
		}
		var b strings.Builder
		s.strValue(c, a[0], o, &b, 1)
		text := b.String()
		if !o.giveHead {
			text = strings.TrimPrefix(text, " ")
		}
		s.out.WriteString(text)
		s.visible = false
		return Null
	})
}
