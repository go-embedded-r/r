package eval

import (
	"bytes"
	"compress/gzip"
	"strings"

	"github.com/go-embedded-r/r/rds"
	"github.com/go-embedded-r/r/syntax"
)

// readRDS, saveRDS, source and sys.source: R's serialization format through
// package rds, and R code read from a file or connection and evaluated.

// fromRDS converts a decoded document to a value. Closures, environments
// other than the four standard ones, external pointers and S4 instances are
// not values of this evaluator yet.
func (s *Session) fromRDS(c *callInfo, v rds.Value, depth int) Value {
	if depth > 1000 {
		s.errorf(c.call, "serialized object is nested too deeply")
	}
	attrs := func(m rds.Meta) []Attr {
		var out []Attr
		for _, a := range m.Attrs {
			out = append(out, Attr{Name: a.Name, Value: s.fromRDS(c, a.Value, depth+1)})
		}
		return out
	}
	switch x := v.(type) {
	case *rds.Altrep:
		m, err := x.Materialize(s.opts.MaxVector)
		if err != nil {
			s.errorf(c.call, "unknown input format")
		}
		return s.fromRDS(c, m, depth+1)
	case *rds.Logical:
		return &Logical{V: append([]int32(nil), x.V...), Attr: attrs(x.Meta)}
	case *rds.Integer:
		return &Integer{V: append([]int32(nil), x.V...), Attr: attrs(x.Meta)}
	case *rds.Real:
		return &Double{V: append([]float64(nil), x.V...), Attr: attrs(x.Meta)}
	case *rds.Character:
		out := &String{V: make([]Str, len(x.V)), Attr: attrs(x.Meta)}
		for i, e := range x.V {
			out.V[i] = Str{S: e.S, NA: e.NA}
		}
		return out
	case *rds.List:
		if x.Expr {
			break
		}
		out := &List{V: make([]Value, len(x.V)), Attr: attrs(x.Meta)}
		for i, e := range x.V {
			out.V[i] = s.fromRDS(c, e, depth+1)
		}
		return out
	case *rds.Special:
		switch x {
		case rds.GlobalEnv:
			return s.Global
		case rds.BaseEnv:
			return s.Base
		case rds.EmptyEnv:
			return s.Empty
		case rds.BaseNamespace:
			return s.BaseNS
		case rds.MissingArg:
			return Missing
		}
	case *rds.Sym:
		return Intern(x.Name)
	case *rds.Pairlist:
		switch x.Kind {
		case rds.LANGSXP:
			if len(x.Items) == 0 {
				break
			}
			l := &Lang{Fun: s.fromRDS(c, x.Items[0].Value, depth+1)}
			for _, it := range x.Items[1:] {
				arg := Arg{Value: s.fromRDS(c, it.Value, depth+1)}
				if it.Tag != nil {
					arg.Name, arg.HasName = it.Tag.Name, true
				}
				l.Args = append(l.Args, arg)
			}
			return l
		case rds.LISTSXP:
			// a pairlist is read as the list of its values
			out := &List{V: make([]Value, len(x.Items))}
			names := make([]Str, len(x.Items))
			named := false
			for i, it := range x.Items {
				out.V[i] = s.fromRDS(c, it.Value, depth+1)
				if it.Tag != nil {
					names[i], named = Str{S: it.Tag.Name}, true
				}
			}
			if named {
				out.Attr = namesAttr(names)
			}
			return out
		}
	}
	if v == rds.Null {
		return Null
	}
	s.errorf(c.call, "this evaluator cannot restore an object of type '%d' from a serialization", v.Type())
	return nil
}

// rdsCheck refuses what ToRDS cannot write.
func (s *Session) rdsCheck(c *callInfo, v Value, depth int) {
	if depth > 1000 {
		s.errorf(c.call, "object is nested too deeply to serialize")
	}
	for _, a := range attrsOf(v) {
		s.rdsCheck(c, a.Value, depth+1)
	}
	switch x := v.(type) {
	case *Closure, *Builtin, *Promise, *Dots:
		s.errorf(c.call, "serializing a %s is not supported by this evaluator", typeName(v))
	case *Env:
		if x != s.Global && x != s.Base && x != s.Empty && x != s.BaseNS {
			s.errorf(c.call, "serializing an environment is not supported by this evaluator")
		}
	case *List:
		for _, e := range x.V {
			s.rdsCheck(c, e, depth+1)
		}
	case *Lang:
		s.rdsCheck(c, x.Fun, depth+1)
		for _, a := range x.Args {
			s.rdsCheck(c, a.Value, depth+1)
		}
	}
}

// readAll is the whole content of a file named by a path or read through a
// file connection; what a failure is, the caller says.
func (s *Session) readAll(c *callInfo, v Value, what string, failCall *Lang, failWarn string) ([]byte, string) {
	switch x := v.(type) {
	case *String:
		p := s.stringArg(c, x, what)
		data, err := s.io.fs.Read(s.abspath(p))
		if err != nil {
			s.warningf(failCall, failWarn, p, errText(err))
			s.errorf(failCall, "cannot open the connection")
		}
		return data, p
	}
	cn := s.getConn(c.call, v, what)
	if cn.kind != kFile && cn.kind != kText {
		s.errorf(c.call, "cannot read from this connection")
	}
	if cn.open {
		if !cn.canRead {
			s.errorf(c.call, "cannot read from this connection")
		}
		data := cn.data[cn.pos:]
		cn.pos = len(cn.data)
		return data, cn.desc
	}
	if cn.kind == kFile {
		data, err := s.io.fs.Read(cn.path)
		if err != nil {
			s.openFailed(c.call, cn, err)
		}
		return data, cn.desc
	}
	s.errorf(c.call, "cannot read from this connection")
	return nil, ""
}

func init() {
	regBuiltin("readRDS", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "file", "refhook")
		if a[0] == nil {
			s.errorf(c.call, "argument \"file\" is missing, with no default")
		}
		call := &Lang{Fun: Intern("gzfile"), Args: []Arg{{Value: Intern("file")}, {Value: mkString("rb")}}}
		data, _ := s.readAll(c, a[0], "file", call, "cannot open compressed file '%s', probable reason '%s'")
		v, _, err := rds.Decode(data, rds.Limits{MaxBytes: s.opts.MaxFileBytes})
		if err != nil {
			s.errorf(c.call, "unknown input format")
		}
		return s.fromRDS(c, v, 0)
	})
	regBuiltin("saveRDS", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "object", "file", "ascii", "version", "compress", "refhook")
		if a[0] == nil {
			s.errorf(c.call, "argument \"object\" is missing, with no default")
		}
		if a[1] == nil {
			s.errorf(c.call, "'file' must be a character string or connection")
		}
		if s.boolArg(c, a[2], false, "ascii") {
			s.errorf(c.call, "'ascii = TRUE' is not supported by this evaluator")
		}
		hdr := rds.DefaultHeader
		if a[3] != nil {
			if _, null := a[3].(nullT); !null {
				switch s.intArg(c, a[3], "version") {
				case 3:
				case 2:
					hdr = rds.Header{Version: 2, Writer: 0x00040602, MinReader: 0x00020300}
				default:
					s.errorf(c.call, "version %d is not supported", s.intArg(c, a[3], "version"))
				}
			}
		}
		compress := "gzip"
		switch x := a[4].(type) {
		case nil:
		case *Logical:
			if len(x.V) == 1 && x.V[0] == 0 {
				compress = ""
			}
		case *String:
			compress = s.stringArg(c, x, "compress")
			if compress != "gzip" {
				s.errorf(c.call, "compression '%s' is not supported by this evaluator (only gzip)", compress)
			}
		}
		s.rdsCheck(c, a[0], 0)
		b, err := rds.EncodeWith(ToRDS(a[0]), hdr, 0)
		if err != nil {
			s.errorf(c.call, "%s", err.Error())
		}
		if compress == "gzip" {
			var buf bytes.Buffer
			zw := gzip.NewWriter(&buf)
			_, _ = zw.Write(b)
			_ = zw.Close()
			b = buf.Bytes()
		}
		s.flushSinks()
		switch x := a[1].(type) {
		case *String:
			p := s.stringArg(c, x, "file")
			if err := s.io.fs.Write(s.abspath(p), b, false); err != nil {
				call := callOf("gzfile(file, mode)")
				s.warningf(call, "cannot open compressed file '%s', probable reason '%s'", p, errText(err))
				s.errorf(call, "cannot open the connection")
			}
		default:
			cn := s.getConn(c.call, x, "file")
			if cn.kind != kFile {
				s.errorf(c.call, "cannot write to this connection")
			}
			if err := s.io.fs.Write(cn.path, b, cn.open); err != nil {
				s.errorf(c.call, "cannot write to this connection")
			}
		}
		return s.invisible(Null)
	})
	regBuiltin("source", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "file", "local", "echo", "print.eval", "exprs", "spaced", "verbose",
			"prompt.echo", "max.deparse.length", "width.cutoff", "deparseCtrl", "chdir", "encoding",
			"continue.echo", "skip.echo", "keep.source")
		if a[0] == nil {
			s.errorf(c.call, "argument \"file\" is missing, with no default")
		}
		if a[4] != nil {
			s.errorf(c.call, "'exprs' is not supported by this evaluator (it needs expression vectors)")
		}
		echo := s.boolArg(c, a[2], false, "echo")
		if echo || s.boolArg(c, a[6], false, "verbose") {
			s.errorf(c.call, "'echo' and 'verbose' are not supported by this evaluator")
		}
		printEval := s.boolArg(c, a[3], echo, "print.eval")
		env := s.Global
		switch l := a[1].(type) {
		case nil:
		case *Env:
			env = l
		default:
			if s.boolArg(c, l, false, "local") {
				env = c.env
			}
		}
		call := &Lang{Fun: Intern("file"), Args: []Arg{{Value: Intern("filename")}, {Value: mkString("r")}, {Name: "encoding", HasName: true, Value: Intern("encoding")}}}
		data, desc := s.readAll(c, a[0], "file", call, "cannot open file '%s': %s")
		if _, isPath := a[0].(*String); !isPath {
			desc = descriptionOf(c, "file", a[0])
		}
		value, visible := s.runSource(c, strings.TrimSuffix(string(data), "\n"), desc, env, printEval)
		if value == nil {
			return s.invisible(Null)
		}
		out := &List{V: []Value{value, mkLogical(visible)}, Attr: namesAttr(strsOf([]string{"value", "visible"}))}
		return s.invisible(out)
	})
	regBuiltin("sys.source", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "file", "envir", "chdir", "keep.source", "keep.parse.data", "toplevel.env")
		if a[0] == nil {
			s.errorf(c.call, "argument \"file\" is missing, with no default")
		}
		env := s.Base
		if a[1] != nil {
			env = s.envArg(c, a[1], "envir")
		}
		call := &Lang{Fun: Intern("file"), Args: []Arg{{Value: Intern("file")}, {Value: mkString("r")}}}
		data, desc := s.readAll(c, a[0], "file", call, "cannot open file '%s': %s")
		s.runSource(c, strings.TrimSuffix(string(data), "\n"), desc, env, false)
		return s.invisible(Null)
	})
}

// sourceEvalCall is the call R names for an error raised directly by code
// that source() evaluates.
var sourceEvalCall = callOf("eval(ei, envir)")

// runSource parses text and evaluates its expressions in env; the value and
// visibility of the last one are returned.
func (s *Session) runSource(c *callInfo, text, desc string, env *Env, printEval bool) (Value, bool) {
	f, err := syntax.Parse(text, syntax.Options{Limits: s.opts.Parse})
	if err != nil {
		msg := err.Error()
		if strings.HasPrefix(msg, "<text>") {
			msg = desc + msg[len("<text>"):]
		}
		s.errorf(c.call, "%s", msg)
	}
	var value Value // nil while nothing was evaluated
	visible := false
	for _, e := range f.Exprs {
		s.visible = true
		value = s.underEval(sourceEvalCall, func() Value { return s.eval(FromSyntax(e), env) })
		visible = s.visible
		if visible && printEval {
			s.printValue(value, env)
		}
	}
	return value, visible
}
