package eval

import (
	"strings"

	"github.com/go-embedded-r/r/syntax"
)

// is.primitive and strtrim, written from their help pages and the oracle.
func init() {
	prim := syntax.PrimitiveNames()
	regBuiltin("is.primitive", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "is.primitive")
		b, ok := x.(*Builtin)
		if ok && prim[b.Name] {
			return &Logical{V: []int32{1}}
		}
		return &Logical{V: []int32{0}}
	})
	regBuiltin("strtrim", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "width")
		x, ok := a[0].(*String)
		if !ok {
			if a[0] == nil {
				s.errorf(c.call, "argument \"x\" is missing, with no default")
			}
			if !isVector(a[0]) {
				s.errorf(c.call, "strtrim() requires a character vector")
			}
			x = &String{V: asStrings(coerceTo(a[0], kString)), Attr: attrsOf(a[0])}
		}
		if a[1] == nil {
			s.errorf(c.call, "argument \"width\" is missing, with no default")
		}
		switch a[1].(type) {
		case *Integer, *Double:
		default:
			s.errorf(c.call, "invalid 'width' argument")
		}
		w := asInts(coerceTo(a[1], kInteger))
		// a shorter width is recycled only when it divides the length of x
		if len(w) == 0 || (len(w) < len(x.V) && len(x.V)%len(w) != 0) {
			s.errorf(c.call, "invalid 'width' argument")
		}
		for _, e := range w {
			if e == naInt || e < 0 {
				s.errorf(c.call, "invalid 'width' argument")
			}
		}
		out := make([]Str, len(x.V))
		for i, e := range x.V {
			if e.NA {
				out[i] = e
				continue
			}
			lim, acc := int(w[i%len(w)]), 0
			var b strings.Builder
			for _, r := range e.S {
				rw := runeWidth(r)
				if acc+rw > lim {
					break
				}
				acc += rw
				b.WriteRune(r)
			}
			out[i] = Str{S: b.String()}
		}
		return withAttr(&String{V: out}, attrsOf(x))
	})
}
