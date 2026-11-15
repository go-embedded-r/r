package eval

import "math"

// summary.table and print.summary.table (slice 18, the test in slice 21),
// written from the help page of summary.table and from the oracle. They carry
// the number of factors and of cases and, for a table of several factors, the
// chi-squared test of independence of all of them: the statistic against the
// product of the margins, and its p-value from pchisq (go-statistics/stats).

// independenceTest answers the statistic, its degrees of freedom, whether every
// expected count is at least 5, and the p-value (the upper tail of the
// chi-squared distribution) of the table x of several factors.
func independenceTest(s *Session, c *callInfo, x Value, total Value) (stat, df float64, approxOK bool, pv Value) {
	dim := getAttr(x, "dim").(*Integer).V
	cells := asDoubles(s.callNamed(c, "unclass", []Value{x}, []string{""}))
	n := asDoubles(total)[0]
	// the relative frequencies of each factor: the margins of the table over n
	margins := make([][]float64, len(dim))
	for k := range dim {
		margins[k] = make([]float64, dim[k])
	}
	stride := 1
	for k := range dim {
		// plain sums in order, as R's: they are the definition of the test
		// statistic, and the rounding is the one R has on this platform
		for i, v := range cells {
			margins[k][(i/stride)%int(dim[k])] += v / n
		}
		stride *= int(dim[k])
	}
	approxOK = true
	for i, v := range cells {
		e, st := 1.0, 1
		for k := range dim {
			e *= margins[k][(i/st)%int(dim[k])]
			st *= int(dim[k])
		}
		e *= n
		d := v - e
		stat += d * d / e
		if !(e >= 5) {
			approxOK = false
		}
	}
	df = 1
	for k := range dim {
		df *= float64(dim[k])
	}
	df--
	for k := range dim {
		df -= float64(dim[k] - 1)
	}
	pv = mkDouble(chisqValue(chisqP, stat, df, 0, false, false, false))
	return
}

func init() {
	regBuiltin("summary.table", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "object", "...")
		if a[0] == nil {
			s.errorf(c.call, "argument \"object\" is missing, with no default")
		}
		x := a[0]
		if !inheritsFrom(x, "table") {
			s.errorf(c.call, "'object' must inherit from class “table”")
		}
		nvars := int32(1)
		if d, ok := getAttr(x, "dim").(*Integer); ok {
			nvars = int32(len(d.V))
		}
		total := s.callNamed(c, "sum", []Value{s.callNamed(c, "unclass", []Value{x}, []string{""})}, []string{""})
		out := &List{V: []Value{mkInt(nvars), total}}
		names := []string{"n.vars", "n.cases"}
		if nvars > 1 {
			stat, df, ok, pv := independenceTest(s, c, x, total)
			out.V = append(out.V, mkDouble(stat), mkDouble(df), &Logical{V: []int32{b2i(ok)}}, pv)
			names = append(names, "statistic", "parameter", "approx.ok", "p.value")
			// the call is an element even when the table has none (NULL)
			call := getAttr(x, "call")
			if call == nil {
				call = Null
			}
			out.V = append(out.V, call)
			names = append(names, "call")
		}
		out.Attr = namesAttr(strsOf(names))
		out.Attr = append(out.Attr, Attr{Name: "class", Value: mkString("summary.table")})
		return out
	})
	regBuiltin("print.summary.table", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "digits", "...")
		x, ok := a[0].(*List)
		if !ok {
			s.errorf(c.call, "invalid 'x' argument")
		}
		get := func(name string) Value {
			if nm := namesOf(x); nm != nil {
				for i, e := range nm.V {
					if e.S == name {
						return x.V[i]
					}
				}
			}
			return nil
		}
		cat := func(args ...Value) {
			names := make([]string, len(args)+1)
			s.callNamed(c, "cat", append(args, mkString("")), append(names[:len(args)], "sep"))
		}
		if call := get("call"); call != nil && call != Null {
			s.callNamed(c, "cat", []Value{mkString("Call: ")}, []string{""})
			s.callNamed(c, "print", []Value{call}, []string{""})
		}
		if n := get("n.cases"); n != nil {
			s.callNamed(c, "cat", []Value{mkString("Number of cases in table:"), n, mkString("\n")}, []string{"", "", ""})
		}
		nv := get("n.vars")
		if nv != nil {
			s.callNamed(c, "cat", []Value{mkString("Number of factors:"), nv, mkString("\n")}, []string{"", "", ""})
		}
		if vars, isInt := nv.(*Integer); isInt && len(vars.V) == 1 && vars.V[0] > 1 {
			digits := a[1]
			if digits == nil {
				d := asDoubles(s.callNamed(c, "getOption", []Value{mkString("digits")}, []string{""}))
				digits = mkDouble(math.Max(1, d[0]-3))
			}
			ch := get("statistic")
			dg := asDoubles(digits)[0]
			chv := asDoubles(ch)[0]
			rounded := s.callNamed(c, "round", []Value{ch, mkDouble(math.Max(0, dg-math.Log10(chv)))}, []string{"", ""})
			pval := s.callNamed(c, "format.pval", []Value{get("p.value"), digits, mkDouble(0)}, []string{"", "", "eps"})
			cat(mkString("Test for independence of all factors:\n"))
			cat(mkString("\tChisq = "), s.callNamed(c, "format", []Value{rounded}, []string{""}),
				mkString(", df = "), get("parameter"), mkString(", p-value = "), pval, mkString("\n"))
			if ok, isLgl := get("approx.ok").(*Logical); isLgl && len(ok.V) == 1 && ok.V[0] == 0 {
				cat(mkString("\tChi-squared approximation may be incorrect\n"))
			}
		}
		return s.invisible(a[0])
	})
}
