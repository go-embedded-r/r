package eval

import (
	"os"
	"strings"

	"github.com/go-embedded-r/r/syntax"
)

// Connections: file connections over the session's file system, text
// connections, and the three terminal ones, with readLines, writeLines, cat
// to a file or connection, sink and capture. Written from the help pages of
// connections, readLines, writeLines, textConnection and sink, and from what
// the oracle does; nothing of R's sources is read.

const (
	kStd = iota
	kFile
	kText
)

// conn is one entry of the session's table of connections.
type conn struct {
	id       int
	kind     int
	class    string
	desc     string
	mode     string // the mode it was opened with; "" while it never was
	open     bool
	canRead  bool
	canWrite bool
	binary   bool
	path     string // file connections

	data []byte // what an open-for-reading connection has left to read
	pos  int

	// output text connections
	lines   []string
	partial []byte
	varEnv  *Env // nil for an anonymous one
	varName string

	sinked int // diversions this connection is the target of
}

type sinkEntry struct {
	c     *conn
	ours  bool // opened by sink(), closed with it
	split bool
}

type ioState struct {
	fs    fileSystem
	root  *os.Root
	cwd   string
	tmp   string
	conns map[int]*conn
	next  int
	sinks [2][]sinkEntry // output, message
	mark  [2]int
	env   map[string]string
}

func (s *Session) initIO() {
	io := &s.io
	if s.opts.FileRoot != "" {
		r, err := os.OpenRoot(s.opts.FileRoot)
		if err != nil {
			panic("eval: Options.FileRoot: " + err.Error())
		}
		io.root = r
		io.fs = &rootFS{root: r, max: int64(s.opts.MaxFileBytes)}
	} else {
		io.fs = newMemFS(int64(s.opts.MaxFileBytes))
	}
	io.cwd = "/"
	io.env = map[string]string{}
	for k, v := range s.opts.Env {
		io.env[k] = v
	}
	io.conns = map[int]*conn{
		0: {id: 0, kind: kStd, class: "terminal", desc: "stdin", mode: "r", open: true, canRead: true},
		1: {id: 1, kind: kStd, class: "terminal", desc: "stdout", mode: "w", open: true, canWrite: true},
		2: {id: 2, kind: kStd, class: "terminal", desc: "stderr", mode: "w", open: true, canWrite: true},
	}
	io.next = 3
}

// Close releases the host directory a session was given with
// Options.FileRoot; a session without one has nothing to release.
func (s *Session) Close() error {
	if s.io.root != nil {
		r := s.io.root
		s.io.root = nil
		return r.Close()
	}
	return nil
}

func connValue(c *conn) Value {
	return &Integer{V: []int32{int32(c.id)}, Attr: []Attr{{Name: "class", Value: strVec([]string{c.class, "connection"})}}}
}

func isConnValue(v Value) bool {
	if _, ok := v.(*Integer); !ok {
		return false
	}
	for _, c := range classAttr(v) {
		if c == "connection" {
			return true
		}
	}
	return false
}

// getConn is the connection a value stands for: R's "invalid connection" for
// one that was closed, "'con' is not a connection" for a plain value.
func (s *Session) getConn(call *Lang, v Value, what string) *conn {
	if !isConnValue(v) || Length(v) != 1 {
		s.errorf(call, "'%s' is not a connection", what)
	}
	cn := s.io.conns[int(v.(*Integer).V[0])]
	if cn == nil {
		s.errorf(call, "invalid connection")
	}
	return cn
}

func (s *Session) newConn(cn *conn) *conn {
	if len(s.io.conns) >= 128 {
		s.errorf(s.currentCall(), "all connections are in use")
	}
	cn.id = s.io.next
	s.io.next++
	s.io.conns[cn.id] = cn
	return cn
}

// ---- opening, closing -----------------------------------------------------

func (s *Session) abspath(p string) string {
	if !strings.HasPrefix(p, "/") {
		p = s.io.cwd + "/" + p
	}
	return cleanPath(p)
}

func cleanPath(p string) string {
	var parts []string
	for _, e := range strings.Split(p, "/") {
		switch e {
		case "", ".":
		case "..":
			if len(parts) > 0 {
				parts = parts[:len(parts)-1]
			}
		default:
			parts = append(parts, e)
		}
	}
	return "/" + strings.Join(parts, "/")
}

// openConn opens cn in the given mode; a failure is R's warning and error,
// raised with call.
func (s *Session) openConn(call *Lang, cn *conn, mode string) {
	if mode == "" {
		mode = "r"
	}
	if cn.kind == kStd {
		return
	}
	if strings.Trim(mode[1:], "+bt") != "" || !strings.ContainsRune("rwa", rune(mode[0])) {
		s.openFailed(call, cn, errInvalid)
	}
	plus := strings.Contains(mode, "+")
	cn.binary = strings.Contains(mode, "b")
	switch mode[0] {
	case 'r':
		var data []byte
		if cn.path != "/dev/null" && cn.desc != "stdin" {
			var err error
			data, err = s.io.fs.Read(cn.path)
			if err != nil {
				s.openFailed(call, cn, err)
			}
		}
		cn.data, cn.pos = data, 0
		cn.canRead, cn.canWrite = true, plus
	case 'w', 'a':
		if cn.path != "/dev/null" {
			if err := s.io.fs.Write(cn.path, nil, mode[0] == 'a'); err != nil {
				s.openFailed(call, cn, err)
			}
		}
		cn.canWrite, cn.canRead = true, false
	default:
		s.errorf(call, "invalid '%s' argument", "open")
	}
	cn.mode, cn.open = mode, true
}

func (s *Session) openFailed(call *Lang, cn *conn, err error) {
	s.warningf(call, "cannot open file '%s': %s", cn.desc, errText(err))
	s.errorf(call, "cannot open the connection")
}

// closeConn closes cn; destroy also removes it from the table.
func (s *Session) closeConn(cn *conn, destroy bool) {
	if cn.kind == kText && cn.canWrite && len(cn.partial) > 0 {
		cn.lines = append(cn.lines, string(cn.partial))
		cn.partial = nil
		s.syncTextVar(cn)
	}
	cn.open, cn.data, cn.pos = false, nil, 0
	if cn.kind == kFile {
		cn.canRead, cn.canWrite = false, false
	}
	if destroy {
		delete(s.io.conns, cn.id)
	}
}

// syncTextVar stores the completed lines of an output text connection in its
// variable.
func (s *Session) syncTextVar(cn *conn) {
	if cn.varEnv == nil {
		return
	}
	out := make([]Str, len(cn.lines))
	for i, l := range cn.lines {
		out[i] = Str{S: l}
	}
	cn.varEnv.Set(cn.varName, &String{V: out})
}

// ---- writing --------------------------------------------------------------

// connWrite writes text to an open connection; a failure is an R error.
func (s *Session) connWrite(call *Lang, cn *conn, text string) {
	if err := s.rawWrite(cn, text); err != nil {
		s.errorf(call, "cannot write to this connection")
	}
}

type notWritable struct{}

func (notWritable) Error() string { return "not writable" }

func (s *Session) rawWrite(cn *conn, text string) error {
	if !cn.canWrite {
		return notWritable{}
	}
	switch cn.kind {
	case kStd:
		if cn.id == 2 {
			s.stderr.WriteString(text)
		} else {
			s.out.WriteString(text)
		}
	case kText:
		cn.partial = append(cn.partial, text...)
		changed := false
		for {
			i := strings.IndexByte(string(cn.partial), '\n')
			if i < 0 {
				break
			}
			cn.lines = append(cn.lines, string(cn.partial[:i]))
			cn.partial = cn.partial[i+1:]
			changed = true
		}
		if changed {
			s.syncTextVar(cn)
		}
	case kFile:
		if cn.path == "/dev/null" {
			return nil
		}
		return s.io.fs.Write(cn.path, []byte(text), true)
	}
	return nil
}

// withWritable runs f on a connection that is open for writing: one that is
// not open is opened in mode "wt" (or "at" when appendMode) for the call and
// closed again.
func (s *Session) withWritable(call, openCall *Lang, cn *conn, appendMode bool, f func()) {
	s.flushSinks()
	if cn.open {
		if !cn.canWrite {
			s.errorf(call, "cannot write to this connection")
		}
		f()
		return
	}
	m := "wt"
	if appendMode {
		m = "at"
	}
	s.openConn(openCall, cn, m)
	defer s.closeConn(cn, false)
	f()
}

// fileCall is the call R shows for a failure to open the file a function was
// given by name: file(con, "r"), or file(file, "w") for cat.
func fileCall(arg, mode string) *Lang {
	return &Lang{Fun: Intern("file"), Args: []Arg{{Value: Intern(arg)}, {Value: mkString(mode)}}}
}

// callOf is the call a source text spells, for the calls R's own closures
// make inside a function (so that an error shows them).
func callOf(src string) *Lang {
	f, err := syntax.ParseAll(src)
	if err != nil || len(f.Exprs) != 1 {
		panic("eval: callOf: " + src)
	}
	return FromSyntax(f.Exprs[0]).(*Lang)
}

// writeTo is where cat and writeLines send text: a path (opened for the
// call), a connection, or the standard output when nil. openCall is the call
// R shows if a path cannot be opened (nil: file(<argName>, "w")).
func (s *Session) writeTo(c *callInfo, target Value, argName string, appendMode bool, text string, openCall *Lang) {
	switch t := target.(type) {
	case nil:
		s.out.WriteString(text)
	case *String:
		if len(t.V) != 1 || t.V[0].NA {
			s.errorf(c.call, "invalid '%s' argument", argName)
		}
		if t.V[0].S == "" {
			s.out.WriteString(text)
			return
		}
		m := "w"
		if appendMode {
			m = "a"
		}
		cn := &conn{kind: kFile, class: "file", desc: t.V[0].S, path: s.abspath(t.V[0].S)}
		if openCall == nil {
			openCall = fileCall(argName, m)
		}
		s.openConn(openCall, cn, m)
		s.connWrite(c.call, cn, text)
		s.closeConn(cn, false)
	default:
		cn := s.getConn(c.call, target, argName)
		s.withWritable(c.call, c.call, cn, appendMode, func() { s.connWrite(c.call, cn, text) })
	}
}

// ---- sinks ----------------------------------------------------------------

func (s *Session) streamBuf(k int) *strings.Builder {
	if k == 0 {
		return &s.out
	}
	return &s.stderr
}

// flushSinks moves what the evaluator wrote since the last flush into the
// connection a sink diverts it to. The streams are written as before and
// diverted here, so that no writer has to know about diversions.
func (s *Session) flushSinks() {
	for k := 0; k < 2; k++ {
		st := s.io.sinks[k]
		if len(st) == 0 {
			continue
		}
		b := s.streamBuf(k)
		str := b.String()
		m := min(s.io.mark[k], len(str))
		seg := str[m:]
		if seg == "" {
			continue
		}
		top := st[len(st)-1]
		if top.split {
			s.io.mark[k] = len(str)
		} else {
			b.Reset()
			b.WriteString(str[:m])
		}
		_ = s.rawWrite(top.c, seg)
	}
}

func (s *Session) pushSink(k int, e sinkEntry) {
	s.flushSinks()
	s.io.sinks[k] = append(s.io.sinks[k], e)
	s.io.mark[k] = s.streamBuf(k).Len()
	e.c.sinked++
}

func (s *Session) popSink(k int) bool {
	s.flushSinks()
	st := s.io.sinks[k]
	if len(st) == 0 {
		return false
	}
	top := st[len(st)-1]
	s.io.sinks[k] = st[:len(st)-1]
	s.io.mark[k] = s.streamBuf(k).Len()
	top.c.sinked--
	if top.ours {
		s.closeConn(top.c, true)
	}
	return true
}
