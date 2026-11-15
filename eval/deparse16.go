package eval

import (
	"strings"

	"github.com/go-embedded-r/r/syntax"
)

// The deparse options (slice 16): width.cutoff, backtick, control and nlines.
// Written from the help page of deparse and from the oracle's answers.

// deparseOptNames are the names deparse() takes for its control argument, in
// the order .deparseOpts lists them.
var deparseOptNames = []string{"all", "keepInteger", "quoteExpressions", "showAttributes", "useSource",
	"warnIncomplete", "delayPromises", "keepNA", "S_compatible", "hexNumeric", "digits17", "niceNames", "exact"}

// deparseControl reads the names of a control argument into what each asks
// for; an unknown name is an error.
func (s *Session) deparseControl(c *callInfo, v Value) deparseCtl {
	have := map[string]bool{}
	switch x := v.(type) {
	case nullT:
	case *String:
		for _, e := range x.V {
			name, ok := "", false
			if !e.NA {
				var hits []string
				for _, o := range deparseOptNames {
					if o == e.S {
						hits = []string{o}
						break
					}
					if strings.HasPrefix(o, e.S) && e.S != "" {
						hits = append(hits, o)
					}
				}
				if len(hits) == 1 {
					name, ok = hits[0], true
				}
			}
			if !ok {
				s.errorf(c.call, "deparse option ‘%s’ is not recognized", e.S)
			}
			switch name {
			case "all", "exact":
				for _, o := range []string{"keepInteger", "quoteExpressions", "showAttributes", "useSource", "warnIncomplete", "keepNA", "niceNames"} {
					have[o] = true
				}
				if name == "all" {
					have["digits17"] = true
				} else {
					have["hexNumeric"] = true
				}
			default:
				have[name] = true
			}
		}
	default:
		s.errorf(c.call, "invalid 'control' argument")
	}
	return deparseCtl{noAttr: !have["showAttributes"], noNice: !have["niceNames"], noKeepInt: !have["keepInteger"],
		noKeepNA: !have["keepNA"], quote: have["quoteExpressions"], digits17: have["digits17"], hex: have["hexNumeric"]}
}

// deparseArgs turns the arguments of deparse() after the expression into the
// options of the syntax writer, what toSyntax is to write, and the number of
// lines to keep (0 for all).
func (s *Session) deparseArgs(c *callInfo, cutoff, backtick, control, nlines, expr Value) (*syntax.DeparseOptions, deparseCtl, int) {
	opt := &syntax.DeparseOptions{}
	if cutoff != nil {
		w := asDoubles(cutoff)
		if len(w) == 0 || w[0] != w[0] || w[0] < 20 || w[0] > 500 {
			s.warningf(c.call, "invalid 'cutoff' value for 'deparse', using default")
		} else {
			opt.Cutoff = int(w[0])
		}
	}
	if backtick != nil {
		opt.NoBackquote = !s.truthy(c, backtick, "backtick")
	} else {
		// backquotes by default only for a call, an expression or a function
		switch x := expr.(type) {
		case *Lang, *Closure:
		case *List:
			opt.NoBackquote = !x.Expr
		default:
			opt.NoBackquote = true
		}
	}
	ctl := deparseCtl{}
	if control != nil {
		ctl = s.deparseControl(c, control)
	}
	opt.NoKeepInteger, opt.NoKeepNA = ctl.noKeepInt, ctl.noKeepNA
	opt.Digits17, opt.HexNumeric = ctl.digits17, ctl.hex
	n := 0
	if nlines != nil {
		if v := asDoubles(nlines); len(v) > 0 && v[0] == v[0] && v[0] > 0 && v[0] < 1<<30 {
			n = int(v[0])
		}
	}
	return opt, ctl, n
}
