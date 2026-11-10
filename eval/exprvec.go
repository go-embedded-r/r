package eval

import (
	"strings"

	"github.com/go-embedded-r/r/syntax"
)

// Expression vectors (EXPRSXP): written from ?expression, ?parse and the
// oracle's behaviour. An expression vector is a *List with Expr set, so every
// list operation applies to it; the ones that must keep the flag (subsetting,
// c, rev, rep, assignment) do.

// parseText parses text as R's parse(text =) does: the lines are joined, n
// bounds the number of expressions (n < 0: all), and a failure is the error R
// signals.
func (s *Session) parseText(c *callInfo, text string, n int, desc string) *List {
	f, err := syntax.Parse(text, syntax.Options{Limits: s.opts.Parse})
	if err != nil {
		msg := err.Error()
		if strings.HasPrefix(msg, "<text>") {
			msg = desc + msg[len("<text>"):]
		}
		s.errorf(c.call, "%s", msg)
	}
	out := &List{V: []Value{}, Expr: true}
	for _, e := range f.Exprs {
		if n >= 0 && len(out.V) >= n {
			break
		}
		out.V = append(out.V, FromSyntax(e))
	}
	return out
}

// textLines is the character data of parse(text =): a character vector, a
// symbol, or anything as.character accepts; NULL is no lines.
func (s *Session) textLines(c *callInfo, v Value) []Str {
	switch x := v.(type) {
	case nullT:
		return nil
	case *Symbol:
		return []Str{{S: x.Name}}
	case *String:
		return x.V
	case *Logical, *Integer, *Double, *Complex, *Raw:
		return asStrings(x)
	}
	s.errorf(c.call, "argument is not a character vector")
	return nil
}

func joinLines(lines []Str) string {
	parts := make([]string, len(lines))
	for i, l := range lines {
		if l.NA {
			parts[i] = "NA"
		} else {
			parts[i] = l.S
		}
	}
	return strings.Join(parts, "\n")
}

func (s *Session) asExpression(c *callInfo, x Value) Value {
	switch v := x.(type) {
	case *List:
		if v.Expr {
			return v
		}
		return &List{V: v.V, Attr: keepNamesOnly(v), Expr: true}
	case nullT:
		return &List{V: []Value{}, Expr: true}
	case *Lang, *Symbol:
		return &List{V: []Value{x}, Expr: true}
	case *Logical, *Integer, *Double, *String, *Complex, *Raw:
		n := Length(v)
		s.checkLen(c.call, n)
		out := &List{V: make([]Value, n), Expr: true, Attr: keepNamesOnly(v)}
		for i := range out.V {
			out.V[i] = element(v, i)
		}
		return out
	}
	s.errorf(c.call, "cannot coerce type '%s' to vector of type 'expression'", typeName(x))
	return nil
}

func init() {
	regSpecial("expression", func(s *Session, c *callInfo) Value {
		out := &List{V: make([]Value, len(c.raw)), Expr: true}
		named := false
		names := make([]Str, len(c.raw))
		for i, a := range c.raw {
			out.V[i] = a.Value
			if a.HasName {
				names[i] = Str{S: a.Name}
				named = true
			}
		}
		if named {
			out.Attr = namesAttr(names)
		}
		return out
	})
	regBuiltin("as.expression", func(s *Session, c *callInfo) Value {
		return s.asExpression(c, firstArg(c))
	})
	regBuiltin("parse", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "file", "n", "text", "prompt", "keep.source", "srcfile", "encoding")
		n := -1
		if a[1] != nil {
			if _, isNull := a[1].(nullT); !isNull {
				d := asDoublesOrNil(a[1])
				if len(d) != 1 || d[0] != d[0] {
					s.errorf(c.call, "invalid '%s' argument", "n")
				}
				if d[0] >= 0 {
					n = int(d[0])
				}
			}
		}
		if a[2] != nil {
			if _, isNull := a[2].(nullT); !isNull || a[0] == nil {
				lines := s.textLines(c, a[2])
				s.chargeBytes(len(joinLines(lines)))
				return s.parseText(c, joinLines(lines), n, "<text>")
			}
		}
		if a[0] == nil {
			s.errorf(c.call, "argument \"file\" is missing, with no default")
		}
		call := &Lang{Fun: Intern("file"), Args: []Arg{{Value: Intern("file")}, {Value: mkString("r")}}}
		data, desc := s.readAll(c, a[0], "file", call, "cannot open file '%s': %s")
		return s.parseText(c, strings.TrimSuffix(string(data), "\n"), n, desc)
	})
	regBuiltin("str2expression", func(s *Session, c *callInfo) Value {
		st, ok := s.nargs1(c, "str2expression").(*String)
		if !ok {
			s.errorf(c.call, "argument must be a character vector")
		}
		return s.parseText(c, joinLines(st.V), -1, "<text>")
	})
	regBuiltin("str2lang", func(s *Session, c *callInfo) Value {
		st, ok := s.nargs1(c, "str2lang").(*String)
		if !ok || len(st.V) != 1 {
			s.errorf(c.call, "argument must be a character string")
		}
		r := s.parseText(c, joinLines(st.V), -1, "<text>")
		if len(r.V) != 1 {
			s.errorf(c.call, "parsing result not of length one, but %d", len(r.V))
		}
		return r.V[0]
	})
}
