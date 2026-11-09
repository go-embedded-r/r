package eval

import (
	"math"
	"strconv"
	"strings"

	"github.com/go-embedded-r/r/syntax"
)

const printWidth = 80
const printDigits = 7

// realFormat is the common format of a double vector: the field width, the
// digits after the point, and whether scientific notation is used (R's
// formatReal with 7 significant digits).
type realFormat struct {
	w, d int
	e    bool
}

func formatRealVec(x []float64, digits int) realFormat { return formatRealVecSci(x, digits, 0) }

// formatRealVecSci is formatRealVec with R's scipen penalty on scientific notation.
func formatRealVecSci(x []float64, digits, scipen int) realFormat {
	naflag, nanflag, posinf, neginf := false, false, false, false
	neg := 0
	mxl, rgt, mxsl, mxns, mnl := math.MinInt32, math.MinInt32, math.MinInt32, math.MinInt32, math.MaxInt32
	for _, f := range x {
		switch {
		case math.IsNaN(f):
			if isNAReal(f) {
				naflag = true
			} else {
				nanflag = true
			}
			continue
		case math.IsInf(f, 1):
			posinf = true
			continue
		case math.IsInf(f, -1):
			neginf = true
			continue
		}
		negI, kpower, nsig, rw := syntax.Scientific(f, digits)
		left := kpower + 1
		if rw {
			left--
		}
		sleft := negI
		if left <= 0 {
			sleft++
		} else {
			sleft += left
		}
		right := nsig - left
		if negI == 1 {
			neg = 1
		}
		if right > rgt {
			rgt = right
		}
		if left > mxl {
			mxl = left
		}
		if left < mnl {
			mnl = left
		}
		if sleft > mxsl {
			mxsl = sleft
		}
		if nsig > mxns {
			mxns = nsig
		}
	}
	var f realFormat
	if mxns != math.MinInt32 {
		if mxl < 0 {
			mxsl = 1 + neg
		}
		if rgt < 0 {
			rgt = 0
		}
		wF := mxsl + rgt
		if rgt != 0 {
			wF++
		}
		e := 1
		if mxl > 100 || mnl <= -99 {
			e = 2
		}
		d := mxns - 1
		w := neg + d + 4 + e
		if d > 0 {
			w++
		}
		if wF <= w+scipen {
			f = realFormat{w: wF, d: rgt}
		} else {
			f = realFormat{w: w, d: d, e: true}
		}
	}
	if naflag && f.w < 2 {
		f.w = 2
	}
	if nanflag && f.w < 3 {
		f.w = 3
	}
	if posinf && f.w < 3 {
		f.w = 3
	}
	if neginf && f.w < 4 {
		f.w = 4
	}
	return f
}

func (f realFormat) encode(x float64) string {
	var s string
	if x == 0 {
		x = 0 // -0 is written 0
	}
	switch {
	case isNAReal(x):
		s = "NA"
	case math.IsNaN(x):
		s = "NaN"
	case math.IsInf(x, 1):
		s = "Inf"
	case math.IsInf(x, -1):
		s = "-Inf"
	case f.e:
		s = strconv.FormatFloat(x, 'e', f.d, 64)
	default:
		s = strconv.FormatFloat(x, 'f', f.d, 64)
	}
	return padLeft(s, f.w)
}

func padLeft(s string, w int) string {
	if n := displayWidth(s); n < w {
		return strings.Repeat(" ", w-n) + s
	}
	return s
}

func padRight(s string, w int) string {
	if n := displayWidth(s); n < w {
		return s + strings.Repeat(" ", w-n)
	}
	return s
}

// formatElements gives the printed text of each element of an atomic vector
// (before padding) and whether the elements are left-justified.
func formatElements(v Value, quote bool) (items []string, left bool) {
	switch x := v.(type) {
	case *Logical:
		for _, e := range x.V {
			switch e {
			case naInt:
				items = append(items, "NA")
			case 0:
				items = append(items, "FALSE")
			default:
				items = append(items, "TRUE")
			}
		}
	case *Integer:
		for _, e := range x.V {
			if e == naInt {
				items = append(items, "NA")
			} else {
				items = append(items, strconv.Itoa(int(e)))
			}
		}
	case *Double:
		f := formatRealVec(x.V, printDigits)
		for _, e := range x.V {
			items = append(items, f.encode(e))
		}
	case *String:
		left = true
		for _, e := range x.V {
			switch {
			case e.NA && quote:
				items = append(items, "NA")
			case e.NA:
				items = append(items, "<NA>")
			case quote:
				items = append(items, syntax.EncodeString(e.S, '"'))
			default:
				items = append(items, e.S)
			}
		}
	}
	return
}

func maxWidth(items []string) int {
	w := 0
	for _, s := range items {
		if n := displayWidth(s); n > w {
			w = n
		}
	}
	return w
}

func (s *Session) printVector(v Value, quote bool) {
	n := Length(v)
	if n == 0 {
		s.out.WriteString(emptyVectorText(v))
		s.out.WriteString("\n")
		return
	}
	items, left := formatElements(v, quote)
	if nm := namesOf(v); nm != nil {
		s.printNamedVector(v, items, nm, quote)
		return
	}
	w := maxWidth(items)
	labw := len(strconv.Itoa(n)) + 2
	width := labw
	label := func(i int) string {
		l := "[" + strconv.Itoa(i) + "]"
		return padLeft(l, labw)
	}
	s.out.WriteString(label(1))
	for i, it := range items {
		if i > 0 && width+w+1 > printWidth {
			s.out.WriteString("\n")
			s.out.WriteString(label(i + 1))
			width = labw
		}
		s.out.WriteString(" ")
		if left {
			s.out.WriteString(padRight(it, w))
		} else {
			s.out.WriteString(padLeft(it, w))
		}
		width += w + 1
	}
	s.out.WriteString("\n")
}

func emptyVectorText(v Value) string {
	named := ""
	if namesOf(v) != nil {
		named = "named "
	}
	switch v.(type) {
	case *Logical:
		return named + "logical(0)"
	case *Integer:
		return named + "integer(0)"
	case *Double:
		return named + "numeric(0)"
	case *String:
		return named + "character(0)"
	}
	return named + "list()"
}

func (s *Session) printNamedVector(v Value, items []string, nm *String, quote bool) {
	names := make([]string, len(items))
	for i := range items {
		switch {
		case i >= len(nm.V) || nm.V[i].NA:
			names[i] = "<NA>"
		default:
			names[i] = nm.V[i].S
		}
	}
	w := maxWidth(items)
	if nw := maxWidth(names); nw > w {
		w = nw
	}
	perLine := printWidth / (w + 1)
	if perLine < 1 {
		perLine = 1
	}
	for start := 0; start < len(items); start += perLine {
		end := start + perLine
		if end > len(items) {
			end = len(items)
		}
		for k := start; k < end; k++ {
			s.out.WriteString(padLeft(names[k], w) + " ")
		}
		s.out.WriteString("\n")
		for k := start; k < end; k++ {
			s.out.WriteString(padLeft(items[k], w) + " ")
		}
		s.out.WriteString("\n")
	}
}

func syntacticName(name string) bool {
	return name != "" && isSyntactic(name)
}

func isSyntactic(name string) bool {
	d := syntax.Deparse(&syntax.Sym{Name: name}, nil)
	return len(d) == 1 && d[0] == name
}

// printValueTagged prints v with the tag prefix of the enclosing lists.
func (s *Session) printValueTagged(v Value, prefix string, env *Env) {
	switch x := v.(type) {
	case nullT:
		s.out.WriteString("NULL\n")
	case *Logical, *Integer, *Double, *String:
		if d := dimsOf(v); len(d) >= 2 {
			s.printArray(v, d, !s.noQuote)
		} else {
			s.printVector(v, !s.noQuote)
		}
		s.printAttributes(v, prefix, env)
	case *List:
		if len(x.V) == 0 {
			if namesOf(x) != nil {
				s.out.WriteString("named list()\n")
			} else {
				s.out.WriteString("list()\n")
			}
			s.printAttributes(v, prefix, env)
			return
		}
		nm := namesOf(x)
		for i, e := range x.V {
			tag := prefix + "[[" + strconv.Itoa(i+1) + "]]"
			if nm != nil && nm.V[i].NA {
				tag = prefix + "$<NA>"
			} else if nm != nil && nm.V[i].S != "" {
				if syntacticName(nm.V[i].S) {
					tag = prefix + "$" + nm.V[i].S
				} else {
					tag = prefix + "$`" + nm.V[i].S + "`"
				}
			}
			s.out.WriteString(tag + "\n")
			// an element that is an object is printed through print()
			if isObject(e) {
				s.printObject(e, tag, s.Global)
			} else {
				s.printValueTagged(e, tag, env)
			}
			s.out.WriteString("\n")
		}
		s.printAttributes(v, prefix, env)
	case *Symbol:
		s.out.WriteString(strings.Join(syntax.Deparse(&syntax.Sym{Name: x.Name}, nil), "\n") + "\n")
	case *Lang:
		s.out.WriteString(strings.Join(syntax.Deparse(toSyntax(x), nil), "\n") + "\n")
	case *Closure:
		fn, ok := toSyntax(&Lang{Fun: Intern("function"), Args: []Arg{{Value: &Formals{F: x.Formals}}, {Value: x.Body}}}).(*syntax.Function)
		if ok {
			s.out.WriteString(strings.Join(syntax.DeparseClosure(fn.Params, fn.Body), "\n") + "\n")
		}
		if x.Env != s.Global && x.Env != s.Base {
			s.out.WriteString("<environment>\n")
		}
	case *Builtin:
		s.out.WriteString("function (...) .Primitive(\"" + x.Name + "\")\n")
	case *Env:
		switch {
		case x == s.BaseNS:
			s.out.WriteString("<environment: namespace:base>\n")
		case envName(s, x) != "":
			s.out.WriteString("<environment: " + envName(s, x) + ">\n")
		default:
			s.out.WriteString("<environment>\n")
		}
	case missingT:
		s.out.WriteString("\n") // the empty symbol prints as an empty line
	default:
		s.out.WriteString("<object>\n")
	}
}

func hasString(l []string, w string) bool {
	for _, e := range l {
		if e == w {
			return true
		}
	}
	return false
}

func (s *Session) printAttributes(v Value, prefix string, env *Env) {
	for _, a := range attrsOf(v) {
		if (a.Name == "names" && len(dimsOf(v)) < 2) || ((a.Name == "dim" || a.Name == "dimnames") && len(dimsOf(v)) >= 1) {
			continue
		}
		s.out.WriteString("attr(,\"" + a.Name + "\")\n")
		if isObject(a.Value) {
			s.printObject(a.Value, "attr(,\""+a.Name+"\")", s.Global)
		} else {
			s.printValueTagged(a.Value, "attr(,\""+a.Name+"\")", env)
		}
	}
}

// printObject is print() for a value with a class attribute: a user method,
// the condition printer, or the default.
func (s *Session) printObject(v Value, prefix string, env *Env, extra ...Arg) {
	classes := classAttr(v)
	for i, cl := range classes {
		if m := s.findMethod("print."+cl, env); m != nil {
			di := &dispatchInfo{generic: "print", classes: append(append([]string{}, classes...), "default"), idx: i, object: v}
			call := &Lang{Fun: Intern("print"), Args: append([]Arg{{Value: Intern("x")}}, extra...)}
			s.callMethod(m, "print."+cl, call, append([]Arg{{Value: v}}, extra...), env, di)
			return
		}
		if cl == "condition" {
			s.printCondition(v)
			return
		}
		if (cl == "factor" || cl == "ordered") && kindOf(v) == kInteger {
			s.printFactor(v)
			return
		}
	}
	s.printValueTagged(v, prefix, env)
}

func (s *Session) printCondition(v Value) {
	cls := classAttr(v)[0]
	msg := condMessage(v)
	if cl := condCall(v); cl != nil {
		s.out.WriteString("<" + cls + " in " + strings.Join(syntax.Deparse(toSyntax(cl), nil), "") + ": " + msg + ">\n")
	} else {
		s.out.WriteString("<" + cls + ": " + msg + ">\n")
	}
}

// printValue is the auto-print of a visible top-level value.
func (s *Session) printValue(v Value, env *Env) {
	if isObject(v) {
		// the methods are looked for from the global environment, where the
		// console prints from
		s.printObject(v, "", s.Global)
		return
	}
	s.printValueTagged(v, "", env)
}

// catString is the text cat() writes for one element.
func catStrings(v Value) []string {
	switch x := v.(type) {
	case nullT:
		return nil
	case *Logical, *Integer, *String:
		items, _ := formatElements(x, false)
		if st, ok := x.(*String); ok {
			for i := range items {
				if st.V[i].NA {
					items[i] = "NA"
				}
			}
		}
		return items
	case *Double:
		out := make([]string, len(x.V))
		for i, e := range x.V {
			f := formatRealVec([]float64{e}, printDigits)
			out[i] = strings.TrimLeft(f.encode(e), " ")
		}
		return out
	case *Symbol:
		return []string{x.Name}
	case *List:
		var out []string
		for _, e := range x.V {
			switch e.(type) {
			case *Logical, *Integer, *Double, *String:
				if Length(e) == 1 {
					out = append(out, catStrings(e)...)
					continue
				}
			case nullT:
				continue
			}
			out = append(out, "<list>")
		}
		return out
	}
	return nil
}

func init() {
	regBuiltin("print", func(s *Session, c *callInfo) Value {
		x := s.nargs1First(c)
		for i, n := range c.names {
			if n == "quote" && s.truthy(c, c.args[i], "quote") == false {
				s.noQuote = true
				defer func() { s.noQuote = false }()
			}
		}
		if isObject(x) {
			var extra []Arg
			for i := 1; i < len(c.args); i++ {
				extra = append(extra, Arg{Name: c.names[i], HasName: c.names[i] != "", Value: c.args[i]})
			}
			s.printObject(x, "", c.env, extra...)
		} else {
			s.printValueTagged(x, "", c.env)
		}
		s.visible = false
		return x
	})
	regBuiltin("cat", func(s *Session, c *callInfo) Value {
		sep := []string{" "}
		var items []string
		var file Value
		appendMode := false
		for i, a := range c.args {
			switch c.names[i] {
			case "sep":
				if st, ok := a.(*String); ok && len(st.V) > 0 {
					sep = nil
					for _, e := range st.V {
						sep = append(sep, e.S)
					}
				}
			case "file":
				file = a
			case "append":
				appendMode = s.boolArg(c, a, false, "append")
			case "fill", "labels":
			default:
				switch a.(type) {
				case *Closure, *Builtin, *Env:
					s.errorf(c.call, "argument %d (type '%s') cannot be handled by 'cat'", i+1, typeName(a))
				}
				items = append(items, catStrings(a)...)
			}
		}
		var sb strings.Builder
		for i, it := range items {
			sb.WriteString(it)
			if i < len(items)-1 {
				sb.WriteString(sep[i%len(sep)])
			}
		}
		// a separator with a newline in it makes the output end with one
		for _, sp := range sep {
			if strings.Contains(sp, "\n") {
				sb.WriteByte('\n')
				break
			}
		}
		if file == nil {
			s.out.WriteString(sb.String())
		} else {
			s.writeTo(c, file, "file", appendMode, sb.String(), catOpenCall)
		}
		s.visible = false
		return Null
	})
}

// catOpenCall is the call R shows when cat cannot open the file it was given.
var catOpenCall = callOf(`file(file, ifelse(append, "a", "w"))`)

func (s *Session) nargs1First(c *callInfo) Value {
	if len(c.args) == 0 {
		s.errorf(c.call, "argument \"x\" is missing, with no default")
	}
	return c.args[0]
}
