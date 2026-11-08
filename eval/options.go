package eval

import "sort"

// options(): a per-session store of named values. Only the options the
// evaluator itself reads have an effect (digits.secs for date-times); the
// others are kept and returned so that code which saves and restores them
// works.

func defaultOptions() map[string]Value {
	return map[string]Value{
		"digits":                 mkInt(7),
		"OutDec":                 mkString("."),
		"width":                  mkInt(80),
		"max.print":              mkInt(99999),
		"scipen":                 mkInt(0),
		"warn":                   mkInt(0),
		"expressions":            mkInt(5000),
		"warning.length":         mkInt(1000),
		"nwarnings":              mkInt(50),
		"useFancyQuotes":         mkLogical(true),
		"stringsAsFactors":       mkLogical(false),
		"keep.source":            mkLogical(false),
		"encoding":               mkString("native.enc"),
		"digits.secs":            nil, // unset
		"warnPartialMatchDollar": mkLogical(false),
	}
}

func init() {
	regBuiltin("options", func(s *Session, c *callInfo) Value {
		if s.options == nil {
			s.options = defaultOptions()
		}
		type kv struct {
			name string
			val  Value
		}
		var sets []kv
		var gets []string
		for i, a := range c.args {
			switch {
			case c.names[i] != "":
				sets = append(sets, kv{c.names[i], a})
			default:
				switch x := a.(type) {
				case *List:
					nm := namesOf(x)
					if nm == nil && len(x.V) > 0 {
						s.errorf(c.call, "list argument has no valid names")
					}
					for j, v := range x.V {
						sets = append(sets, kv{nm.V[j].S, v})
					}
				case *String:
					for _, e := range x.V {
						gets = append(gets, e.S)
					}
				case nullT:
				default:
					s.errorf(c.call, "invalid argument")
				}
			}
		}
		listOf := func(names []string) *List {
			out := &List{V: make([]Value, len(names))}
			for i, n := range names {
				if v := s.options[n]; v != nil {
					out.V[i] = v
				} else {
					out.V[i] = Null
				}
			}
			out.Attr = namesAttr(strsOf(names))
			return out
		}
		if len(sets) > 0 {
			old := make([]string, len(sets))
			for i, e := range sets {
				old[i] = e.name
			}
			res := listOf(old)
			for _, e := range sets {
				if _, isNull := e.val.(nullT); isNull {
					delete(s.options, e.name)
				} else {
					s.options[e.name] = e.val
				}
			}
			return s.invisible(res)
		}
		if len(gets) > 0 {
			return listOf(gets)
		}
		names := make([]string, 0, len(s.options))
		for n, v := range s.options {
			if v != nil {
				names = append(names, n)
			}
		}
		sort.Slice(names, func(i, j int) bool { return collate(names[i], names[j]) < 0 })
		return listOf(names)
	})
}
