package eval

// summary.table and print.summary.table (slice 18), written from the help
// page of summary.table and from the oracle. They carry the number of factors
// and of cases; the chi-squared test of independence is not computed (it
// needs pchisq, which is not implemented): the README records the divergence.

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
		out := &List{V: []Value{mkInt(nvars), total}, Attr: namesAttr(strsOf([]string{"n.vars", "n.cases"}))}
		if nvars > 1 {
			if call := getAttr(x, "call"); call != nil {
				out.V = append(out.V, call)
				out.Attr = namesAttr(strsOf([]string{"n.vars", "n.cases", "call"}))
			}
		}
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
		if call := get("call"); call != nil {
			s.callNamed(c, "cat", []Value{mkString("Call: ")}, []string{""})
			s.callNamed(c, "print", []Value{call}, []string{""})
		}
		if n := get("n.cases"); n != nil {
			s.callNamed(c, "cat", []Value{mkString("Number of cases in table:"), n, mkString("\n")}, []string{"", "", ""})
		}
		if n := get("n.vars"); n != nil {
			s.callNamed(c, "cat", []Value{mkString("Number of factors:"), n, mkString("\n")}, []string{"", "", ""})
		}
		return s.invisible(a[0])
	})
}
