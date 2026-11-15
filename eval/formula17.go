package eval

// Formulas (slice 17): the `~` primitive. Written from the help pages
// ?tilde and ?formula and from the oracle's behaviour.

func init() {
	regSpecial("~", func(s *Session, c *callInfo) Value {
		s.visible = true
		if getAttr(c.call, "class") != nil {
			return c.call // a formula that is evaluated again stays as it is
		}
		out := &Lang{Fun: c.call.Fun, Args: c.call.Args}
		out.Attr = []Attr{{Name: "class", Value: mkString("formula")}, {Name: ".Environment", Value: c.env}}
		return out
	})
}

// langAsList is the call as the list of its parts, the function first.
func langAsList(l *Lang) *List {
	out := &List{V: make([]Value, 0, len(l.Args)+1)}
	out.V = append(out.V, l.Fun)
	named := false
	for _, a := range l.Args {
		out.V = append(out.V, a.Value)
		named = named || a.HasName
	}
	if named {
		nm := make([]Str, len(out.V))
		for i, a := range l.Args {
			nm[i+1] = Str{S: a.Name}
		}
		out.Attr = namesAttr(nm)
	}
	return out
}

// listAsLang is as.call for a list, keeping the attributes given.
func listAsLang(v *List, attr []Attr) Value {
	if len(v.V) == 0 {
		return Null
	}
	nm := namesOf(v)
	l := &Lang{Fun: v.V[0], Attr: attr}
	for i := 1; i < len(v.V); i++ {
		a := Arg{Value: v.V[i]}
		if nm != nil && nm.V[i].S != "" && !nm.V[i].NA {
			a.Name, a.HasName = nm.V[i].S, true
		}
		l.Args = append(l.Args, a)
	}
	return l
}

// subassignLang runs a replacement on a call by way of its list form.
func (s *Session) subassignLang(c *callInfo, l *Lang, run func(*callInfo) Value) Value {
	c2 := *c
	c2.args = append([]Value{langAsList(l)}, c.args[1:]...)
	res := run(&c2)
	if lv, ok := res.(*List); ok {
		return listAsLang(lv, l.Attr)
	}
	return res
}

func init() {
	// registerS3method(genname, class, method, envir): the method is found by
	// dispatch, without a visible name of its own, as the methods of a package
	// that its namespace does not export.
	regBuiltin("registerS3method", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "genname", "class", "method", "envir")
		for i, name := range []string{"genname", "class", "method"} {
			if a[i] == nil {
				s.errorf(c.call, "argument \"%s\" is missing, with no default", name)
			}
		}
		gen, cls := s.stringArg(c, a[0], "genname"), s.stringArg(c, a[1], "class")
		m := a[2]
		if st, ok := m.(*String); ok && len(st.V) == 1 {
			env := c.env
			if e, ok := a[3].(*Env); ok {
				env = e
			}
			m = s.findMethod(st.V[0].S, env)
		}
		switch m.(type) {
		case *Closure, *Builtin:
		default:
			s.errorf(c.call, "bad method")
		}
		if s.s3reg == nil {
			s.s3reg = map[string]Value{}
		}
		s.s3reg[gen+"."+cls] = m
		s.visible = false
		return Null
	})
}
