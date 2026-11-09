package eval

import (
	"strings"
)

func (s *Session) boolArg(c *callInfo, v Value, def bool, what string) bool {
	if v == nil {
		return def
	}
	switch x := v.(type) {
	case *Logical:
		if len(x.V) >= 1 && x.V[0] != naInt {
			return x.V[0] != 0
		}
	case *Integer, *Double:
		if d := asDoublesOrNil(v); len(d) >= 1 && d[0] == d[0] {
			return d[0] != 0
		}
	}
	s.errorf(c.call, "invalid '%s' argument", what)
	return def
}

// optString is a length-one character argument that may be absent.
func (s *Session) optString(c *callInfo, v Value, def, what string) string {
	if v == nil {
		return def
	}
	return s.stringArg(c, v, what)
}

// noMethod is R's error for a generic that finds no method.
func (s *Session) noMethod(generic string, x Value) {
	cl := dispatchClass(x)
	desc := `"` + cl[0] + `"`
	if len(cl) > 1 {
		desc = `"c('` + strings.Join(cl, "', '") + `')"`
	}
	call := &Lang{Fun: Intern("UseMethod"), Args: []Arg{{Value: mkString(generic)}}}
	s.errorf(call, "no applicable method for '%s' applied to an object of class %s", generic, desc)
}

// connGeneric registers a generic over connections: the builtin dispatches on
// the class of its first argument (user methods, then the .connection one).
func connGeneric(name string, method func(s *Session, c *callInfo) Value) {
	regBuiltin(name+".connection", method)
	regBuiltin(name, func(s *Session, c *callInfo) Value {
		if len(c.args) > 0 {
			if v, ok := s.dispatchInternal(name, c, c.args, c.names); ok {
				return v
			}
			s.noMethod(name, c.args[0])
		}
		s.errorf(c.call, "argument \"con\" is missing, with no default")
		return nil
	})
}

// descriptionOf is deparse1(substitute(object)) for the first argument of a
// builtin: the text R gives a text connection as its name.
func descriptionOf(c *callInfo, formal string, v Value) string {
	if c.call != nil {
		for _, a := range c.call.Args {
			if a.HasName && a.Name == formal {
				return strings.Join(deparseValue(a.Value), " ")
			}
		}
		for _, a := range c.call.Args {
			if !a.HasName {
				return strings.Join(deparseValue(a.Value), " ")
			}
		}
	}
	return strings.Join(deparseValue(v), " ")
}

// splitLines reads up to n lines (all when n < 0) from the unread data of cn;
// the last one is reported incomplete when it has no end of line.
func splitLines(cn *conn, n int) (lines []string, incomplete bool) {
	d := cn.data
	for (n < 0 || len(lines) < n) && cn.pos < len(d) {
		i := cn.pos
		for i < len(d) && d[i] != '\n' && d[i] != '\r' {
			i++
		}
		lines = append(lines, string(d[cn.pos:i]))
		if i >= len(d) {
			cn.pos = len(d)
			incomplete = true
			break
		}
		if d[i] == '\r' && i+1 < len(d) && d[i+1] == '\n' {
			i++
		}
		cn.pos = i + 1
	}
	return lines, incomplete
}

func (s *Session) modeName(cn *conn) (mode, text, opened, canRead, canWrite string) {
	yes := func(b bool) string {
		if b {
			return "yes"
		}
		return "no"
	}
	mode = cn.mode
	if mode == "" {
		mode = "r"
	}
	text = "text"
	if cn.binary {
		text = "binary"
	}
	opened = "closed"
	if cn.open {
		opened = "opened"
	}
	if cn.kind == kFile && !cn.open {
		return mode, text, opened, "yes", "yes"
	}
	return mode, text, opened, yes(cn.canRead), yes(cn.canWrite)
}

func init() {
	regBuiltin("textConnection", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "object", "open", "local", "name", "encoding")
		open := s.optString(c, a[1], "r", "open")
		local := s.boolArg(c, a[2], false, "local")
		if a[0] == nil {
			s.errorf(c.call, "argument \"object\" is missing, with no default")
		}
		name := descriptionOf(c, "object", a[0])
		if a[3] != nil {
			name = s.stringArg(c, a[3], "name")
		}
		cn := &conn{kind: kText, class: "textConnection", desc: name, mode: open, open: true}
		switch {
		case open == "r" || open == "rt":
			st, ok := a[0].(*String)
			if !ok {
				s.errorf(c.call, "invalid 'text' argument")
			}
			var sb strings.Builder
			for _, e := range st.V {
				if e.NA {
					sb.WriteString("NA")
				} else {
					sb.WriteString(e.S)
				}
				sb.WriteByte('\n')
			}
			cn.data, cn.canRead = []byte(sb.String()), true
		case open == "w" || open == "a" || open == "wt" || open == "at":
			cn.canWrite = true
			switch x := a[0].(type) {
			case nullT:
			case *String:
				if len(x.V) != 1 || x.V[0].NA {
					s.errorf(c.call, "invalid 'object' argument")
				}
				cn.varName = x.V[0].S
				cn.varEnv = s.Global
				if local {
					cn.varEnv = c.env
				}
				if open[0] == 'a' {
					if old, ok := cn.varEnv.Get(cn.varName); ok {
						if os, ok := old.(*String); ok {
							cn.lines = append(cn.lines, os.V...)
						}
					} else {
						s.warningf(c.call, "text connection: appending to a non-existent char vector")
					}
				}
				s.syncTextVar(cn)
			default:
				s.errorf(c.call, "invalid 'object' argument")
			}
		default:
			s.errorf(c.call, "unsupported mode")
		}
		return connValue(s.newConn(cn))
	})
	regBuiltin("textConnectionValue", func(s *Session, c *callInfo) Value {
		cn := s.getConn(c.call, s.nargs1(c, "textConnectionValue"), "con")
		if cn.kind != kText || !cn.canWrite {
			s.errorf(c.call, "'con' is not an output textConnection")
		}
		s.flushSinks()
		return &String{V: cn.lines[:len(cn.lines):len(cn.lines)]}
	})
	regBuiltin("file",func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "description", "open", "blocking", "encoding", "raw", "method")
		desc := s.optString(c, a[0], "", "description")
		open := s.optString(c, a[1], "", "open")
		if desc == "" {
			s.errorf(c.call, "an anonymous temporary file (description = \"\") is not supported")
		}
		cn := &conn{kind: kFile, class: "file", desc: desc, path: s.abspath(desc)}
		if open != "" {
			s.openConn(c.call, cn, open)
		}
		return connValue(s.newConn(cn))
	})
	connGeneric("open", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "con", "open", "blocking")
		cn := s.getConn(c.call, a[0], "con")
		if cn.open {
			s.warningf(c.call, "connection is already open")
			return s.invisible(Null)
		}
		s.openConn(c.call, cn, s.optString(c, a[1], "r", "open"))
		return s.invisible(Null)
	})
	connGeneric("close", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "con", "type")
		cn := s.getConn(c.call, a[0], "con")
		if cn.kind == kStd {
			s.errorf(c.call, "cannot close standard connections")
		}
		if cn.sinked > 0 {
			s.errorf(c.call, "cannot close 'output' sink connection")
		}
		s.flushSinks()
		s.closeConn(cn, true)
		return s.invisible(Null)
	})
	connGeneric("flush", func(s *Session, c *callInfo) Value {
		s.getConn(c.call, s.nargs1(c, "flush"), "con")
		s.flushSinks()
		return s.invisible(Null)
	})
	regBuiltin("isOpen", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "con", "rw")
		cn := s.getConn(c.call, a[0], "con")
		rw := s.optString(c, a[1], "", "rw")
		switch {
		case rw == "":
			return mkLogical(cn.open)
		case strings.HasPrefix("read", rw):
			return mkLogical(cn.open && cn.canRead)
		case strings.HasPrefix("write", rw):
			return mkLogical(cn.open && cn.canWrite)
		}
		s.errorf(c.call, "unknown 'rw' value")
		return nil
	})
	regBuiltin("isIncomplete", func(s *Session, c *callInfo) Value {
		cn := s.getConn(c.call, s.nargs1(c, "isIncomplete"), "con")
		s.flushSinks()
		return mkLogical(cn.kind == kText && cn.canWrite && len(cn.partial) > 0)
	})
	std := func(name string, id int) {
		regBuiltin(name, func(s *Session, c *callInfo) Value { return connValue(s.io.conns[id]) })
	}
	std("stdin", 0)
	std("stdout", 1)
	std("stderr", 2)
	regBuiltin("nullfile", func(s *Session, c *callInfo) Value { return mkString("/dev/null") })
	regBuiltin("isatty", func(s *Session, c *callInfo) Value {
		s.getConn(c.call, s.nargs1(c, "isatty"), "con")
		return mkLogical(false)
	})
	regBuiltin("summary.connection", func(s *Session, c *callInfo) Value {
		cn := s.getConn(c.call, s.nargs1First(c), "object")
		mode, text, opened, canRead, canWrite := s.modeName(cn)
		vals := []string{cn.desc, cn.class, mode, text, opened, canRead, canWrite}
		out := &List{V: make([]Value, len(vals))}
		for i, v := range vals {
			out.V[i] = mkString(v)
		}
		out.Attr = namesAttr(strsOf([]string{"description", "class", "mode", "text", "opened", "can read", "can write"}))
		return out
	})
	regBuiltin("readLines", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "con", "n", "ok", "warn", "encoding", "skipNul")
		n := -1
		if a[1] != nil {
			d := asDoublesOrNil(a[1])
			if len(d) != 1 || d[0] != d[0] {
				s.errorf(c.call, "invalid '%s' argument", "n")
			}
			n = int(max(min(d[0], 1<<31-1), -1))
		}
		ok := s.boolArg(c, a[2], true, "ok")
		warn := s.boolArg(c, a[3], true, "warn")
		s.flushSinks()
		var cn *conn
		temp := false
		switch x := a[0].(type) {
		case nil:
			cn = s.io.conns[0]
		case *String:
			desc := s.stringArg(c, x, "con")
			cn = &conn{kind: kFile, class: "file", desc: desc, path: s.abspath(desc)}
			s.openConn(fileCall("con", "r"), cn, "r")
			temp = true
		default:
			cn = s.getConn(c.call, a[0], "con")
		}
		if !cn.open {
			s.openConn(c.call, cn, "rt")
			temp = true
		}
		defer func() {
			if temp {
				s.closeConn(cn, false)
			}
		}()
		if !cn.canRead {
			s.errorf(c.call, "cannot read from this connection")
		}
		lines, incomplete := splitLines(cn, n)
		if incomplete && warn {
			s.warningf(c.call, "incomplete final line found on '%s'", cn.desc)
		}
		if n > 0 && len(lines) < n && !ok {
			s.errorf(c.call, "too few lines read in readLines")
		}
		out := make([]Str, len(lines))
		for i, l := range lines {
			out[i] = Str{S: l}
		}
		return &String{V: out}
	})
	regBuiltin("writeLines", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "text", "con", "sep", "useBytes")
		text, ok := a[0].(*String)
		if !ok {
			s.errorf(c.call, "can only write character objects")
		}
		sep := s.optString(c, a[2], "\n", "sep")
		var sb strings.Builder
		for _, e := range text.V {
			if e.NA {
				sb.WriteString("NA")
			} else {
				sb.WriteString(e.S)
			}
			sb.WriteString(sep)
		}
		switch x := a[1].(type) {
		case nil:
			s.writeTo(c, nil, "con", false, sb.String(), nil)
		case *String:
			s.writeTo(c, x, "con", false, sb.String(), nil)
		default:
			if !isConnValue(x) {
				s.errorf(c.call, "'con' is not a connection")
			}
			s.writeTo(c, x, "con", false, sb.String(), nil)
		}
		return s.invisible(Null)
	})
	regBuiltin("readline", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "prompt")
		if a[0] != nil {
			s.out.WriteString(s.optString(c, a[0], "", "prompt"))
		}
		s.out.WriteString("\n")
		return mkString("")
	})
	regBuiltin("sink", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "file", "append", "type", "split")
		k := s.sinkType(c, a[2])
		appendMode := s.boolArg(c, a[1], false, "append")
		split := s.boolArg(c, a[3], false, "split")
		s.flushSinks()
		switch x := a[0].(type) {
		case nil, nullT:
			if !s.popSink(k) {
				s.warningf(c.call, "no sink to remove")
			}
		case *String:
			if len(x.V) != 1 || x.V[0].NA {
				s.errorf(c.call, "invalid connection")
			}
			m := "w"
			if appendMode {
				m = "a"
			}
			cn := &conn{kind: kFile, class: "file", desc: x.V[0].S, path: s.abspath(x.V[0].S)}
			s.openConn(callOf(`file(file, if (append) "a" else "w")`), cn, m)
			s.newConn(cn)
			s.pushSink(k, sinkEntry{c: cn, ours: true, split: split})
		default:
			if !isConnValue(x) {
				s.errorf(c.call, "'file' must be NULL, a connection or a character string")
			}
			cn := s.getConn(c.call, x, "file")
			ours := false
			if !cn.open {
				s.openConn(c.call, cn, "wt")
				ours = true
			}
			if !cn.canWrite {
				s.errorf(c.call, "cannot write to this connection")
			}
			if len(s.io.sinks[k]) >= 20 {
				s.errorf(c.call, "sink stack is full")
			}
			s.pushSink(k, sinkEntry{c: cn, ours: ours, split: split})
		}
		return s.invisible(Null)
	})
	regBuiltin("sink.number", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "type")
		k := s.sinkType(c, a[0])
		st := s.io.sinks[k]
		if k == 1 {
			if len(st) == 0 {
				return mkInt(2)
			}
			return mkInt(int32(st[len(st)-1].c.id))
		}
		return mkInt(int32(len(st)))
	})
}

// sinkType reads type = c("output", "message"), partially matched.
func (s *Session) sinkType(c *callInfo, v Value) int {
	if v == nil {
		return 0
	}
	st, ok := v.(*String)
	if !ok || len(st.V) < 1 || st.V[0].NA {
		s.errorf(c.call, "invalid '%s' argument", "type")
	}
	t := st.V[0].S
	switch {
	case t != "" && strings.HasPrefix("output", t):
		return 0
	case t != "" && strings.HasPrefix("message", t):
		return 1
	}
	s.errorf(c.call, "'arg' should be one of %s, %s", "‘output’", "‘message’")
	return 0
}
