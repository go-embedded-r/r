// Package oracle drives the reference R as a differential oracle.
//
// One R process is started per session and kept alive. A request is a length
// prefix and an R serialization (written by package rds) of a list holding the
// operation and its arguments; a reply is one line of hex holding the
// serialization of a list R built (see oracle.R for the framing). Everything
// the test suite and the golden files know about "what R does" comes through
// here, from the release pinned in [Version].
//
// R's sources are GPL-2+ and are never copied into this repository: R is
// only ever run, as a black box, and what it prints is what the goldens hold.
package oracle

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/go-embedded-r/r/rds"
)

// Version is the pinned reference release (R.version$major + "." + minor).
const Version = "4.6.2"

// SVNRevision is the Subversion revision of the pinned release,
// R.version$`svn rev`.
const SVNRevision = "90784"

//go:embed oracle.R
var script string

// bootstrap is the -e program: Rscript passes -e through a fixed-size buffer
// that the oracle's source outgrew, so the source is sent on stdin instead.
const bootstrap = `.rgo_stdin <- file("stdin", "rb"); eval(parse(text = rawToChar(readBin(.rgo_stdin, "raw", readBin(.rgo_stdin, "integer", 1L, size = 4L, endian = "big")))))`

// Script returns the R source of the oracle.
func Script() string { return script }

// ErrNoR is returned by [Command] when no R can be found.
var ErrNoR = errors.New("oracle: no Rscript found (set RGO_RSCRIPT, or put the pinned release on PATH)")

// ErrVersion is wrapped by the error Start returns when the R it found is not
// the pinned release (CI runners ship their own R).
var ErrVersion = errors.New("oracle: wrong R release")

// Unavailable reports whether err means "the pinned R is not here": no R, or
// another release. Tests that need the oracle skip on it (and fail instead
// when RGO_REQUIRE_ORACLE is set).
func Unavailable(err error) bool { return errors.Is(err, ErrNoR) || errors.Is(err, ErrVersion) }

// ErrR is wrapped by every error that R itself reported in answer to a
// request (as opposed to a failure of the process or the protocol).
var ErrR = errors.New("R error")

// Command returns the program and arguments that run the reference R.
// $RGO_RSCRIPT, when set, is the Rscript executable; otherwise Rscript is
// looked up on PATH.
func Command() ([]string, error) {
	if r := os.Getenv("RGO_RSCRIPT"); r != "" {
		return []string{r}, nil
	}
	p, err := exec.LookPath("Rscript")
	if err != nil {
		return nil, ErrNoR
	}
	// an R of another release on PATH (a CI image may carry one) is no oracle
	if out, err := exec.Command(p, "--version").CombinedOutput(); err != nil || !strings.Contains(string(out), "version "+Version) {
		return nil, ErrNoR
	}
	return []string{p}, nil
}

// Options configure a session. The zero value is usable.
type Options struct {
	// Timeout bounds every Call (default 60 s). A call that exceeds it kills
	// the process; the session is then over.
	Timeout time.Duration
	// MaxVSize is R_MAX_VSIZE (default "2G"): the limit on R's vector heap,
	// so that a runaway allocation fails inside R.
	MaxVSize string
}

// Process is a running oracle.
type Process struct {
	mu      sync.Mutex
	cmd     *exec.Cmd
	in      io.WriteCloser
	out     *bufio.Reader
	stderr  *cappedBuffer
	timeout time.Duration
	greet   Greeting
	// broken is the failure that ended the session: once a write or a read
	// has failed the process is gone (or the stream out of step), and every
	// later Call reports it rather than racing the dead pipe.
	broken  error
	once    sync.Once
	waitErr error
}

// Greeting is what R reports about itself before the first request.
type Greeting struct {
	Version  string // "4.6.2"
	SVN      string // "90784"
	Platform string // R.version$platform, e.g. "aarch64-apple-darwin20"
	Banner   string // R.version.string
}

// Start launches the oracle and checks that it is the pinned release.
func Start(opt Options) (*Process, error) {
	argv, err := Command()
	if err != nil {
		return nil, err
	}
	if opt.Timeout <= 0 {
		opt.Timeout = 60 * time.Second
	}
	if opt.MaxVSize == "" {
		opt.MaxVSize = "2G"
	}
	argv = append(argv, "--vanilla", "-e", bootstrap)
	cmd := exec.Command(argv[0], argv[1:]...)
	// A fixed locale and zone: goldens must not depend on the host.
	// R 4.6 ignores the --max-vsize option ("no longer supported", printed on
	// stdout); the limit is the R_MAX_VSIZE environment variable.
	cmd.Env = append(os.Environ(), "LC_ALL=C.UTF-8", "LANG=C.UTF-8", "TZ=UTC", "R_MAX_VSIZE="+opt.MaxVSize)
	stderr := &cappedBuffer{max: 8 << 10}
	cmd.Stderr = stderr
	// StdinPipe and StdoutPipe fail only when the stream is already set or
	// the process has started, neither of which is the case here.
	in, _ := cmd.StdinPipe()
	out, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("oracle: starting %q: %w", argv[0], err)
	}
	frame := binary.BigEndian.AppendUint32(nil, uint32(len(script)))
	frame = append(frame, script...)
	// A failed write means the process is already gone (it ends before it
	// reads its script when, for instance, R cannot start): the greeting read
	// below then fails with what the process said on its standard error.
	in.Write(frame)
	p := &Process{cmd: cmd, in: in, out: bufio.NewReaderSize(out, 1<<20), stderr: stderr, timeout: opt.Timeout}
	r, err := p.readReply()
	if err != nil {
		p.Close()
		return nil, fmt.Errorf("oracle: no greeting: %w", err)
	}
	p.greet = Greeting{
		Version:  r.String("version"),
		SVN:      r.String("svn"),
		Platform: r.String("platform"),
		Banner:   r.String("banner"),
	}
	if p.greet.Version != Version || p.greet.SVN != SVNRevision {
		g := p.greet
		p.Close()
		return nil, fmt.Errorf("%w: R %s (r%s) is running, the pinned release is %s (r%s)", ErrVersion, g.Version, g.SVN, Version, SVNRevision)
	}
	return p, nil
}

// Greeting is the identity R reported at start-up.
func (p *Process) Greeting() Greeting { return p.greet }

// Reply is a decoded answer. Fields not sent by R are zero.
type Reply struct {
	OK       bool
	Value    rds.Value // the value of the request, as R serialized it
	Output   []string  // captured standard output, one string per line
	Warnings []string
	Class    []string // the condition classes when !OK
	Message  string   // the condition message when !OK
	fields   map[string]rds.Value
}

// String returns a character field of the raw reply (first element).
func (r *Reply) String(name string) string {
	if c, ok := r.fields[name].(*rds.Character); ok && len(c.V) > 0 {
		return c.V[0].S
	}
	return ""
}

func strs(v rds.Value) []string {
	c, ok := v.(*rds.Character)
	if !ok {
		return nil
	}
	out := make([]string, len(c.V))
	for i, s := range c.V {
		out[i] = s.S
	}
	return out
}

func parseReply(raw []byte) (*Reply, error) {
	v, _, err := rds.Decode(raw, rds.Limits{})
	if err != nil {
		return nil, err
	}
	l, ok := v.(*rds.List)
	if !ok {
		return nil, errors.New("oracle: reply is not a list")
	}
	names, _ := l.Attr("names").(*rds.Character)
	if names == nil || len(names.V) != len(l.V) {
		return nil, errors.New("oracle: reply has no names")
	}
	r := &Reply{fields: map[string]rds.Value{}}
	for i, n := range names.V {
		r.fields[n.S] = l.V[i]
	}
	if ok, _ := r.fields["ok"].(*rds.Logical); ok != nil && len(ok.V) == 1 {
		r.OK = ok.V[0] == 1
	}
	r.Value = r.fields["value"]
	r.Output = strs(r.fields["output"])
	r.Warnings = strs(r.fields["warnings"])
	r.Class = strs(r.fields["class"])
	r.Message = r.String("message")
	return r, nil
}

func (p *Process) readReply() (*Reply, error) {
	line, err := p.out.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("read: %w%s", err, p.stderrTail())
	}
	raw, err := hex.DecodeString(strings.TrimRight(line, "\r\n"))
	if err != nil {
		return nil, fmt.Errorf("malformed reply: %w", err)
	}
	return parseReply(raw)
}

// stderrTail ends the process (it is broken when this is asked) and waits
// for it, which also waits for the copy of its standard error to finish, so
// that what R said before it died is in the message.
func (p *Process) stderrTail() string {
	p.cmd.Process.Kill()
	p.wait()
	if s := strings.TrimSpace(p.stderr.String()); s != "" {
		return "; R said: " + s
	}
	return ""
}

// Call sends one request. Each argument is a string (a character vector of
// length one), a []byte (a raw vector) or an [rds.Value]. When R reports
// failure the reply is returned together with an error wrapping [ErrR].
func (p *Process) Call(op string, args ...any) (*Reply, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.broken != nil {
		return nil, fmt.Errorf("oracle: process is gone: %w", p.broken)
	}
	req := &rds.List{V: []rds.Value{rds.Strings(op)}}
	for _, a := range args {
		switch x := a.(type) {
		case string:
			req.V = append(req.V, rds.Strings(x))
		case []byte:
			req.V = append(req.V, &rds.Raw{V: x})
		case rds.Value:
			req.V = append(req.V, x)
		default:
			return nil, fmt.Errorf("oracle: unsupported argument type %T", a)
		}
	}
	body, err := rds.Encode(req)
	if err != nil {
		return nil, err
	}
	frame := binary.BigEndian.AppendUint32(nil, uint32(len(body)))
	frame = append(frame, body...)
	timer := time.AfterFunc(p.timeout, func() { p.cmd.Process.Kill() })
	defer timer.Stop()
	if _, err := p.in.Write(frame); err != nil {
		p.broken = fmt.Errorf("oracle: write: %w%s", err, p.stderrTail())
		return nil, p.broken
	}
	r, err := p.readReply()
	if err != nil {
		p.broken = fmt.Errorf("oracle: %w", err)
		return nil, p.broken
	}
	if !r.OK {
		return r, fmt.Errorf("%w: %s: %s", ErrR, strings.Join(r.Class, "/"), r.Message)
	}
	return r, nil
}

// Close stops the oracle: R exits when its input ends.
func (p *Process) Close() error {
	p.in.Close()
	return p.wait()
}

func (p *Process) wait() error {
	p.once.Do(func() { p.waitErr = p.cmd.Wait() })
	return p.waitErr
}

// Eval evaluates the R source text (every expression, in a fresh child of
// the global environment) and returns the reply for the last value.
func (p *Process) Eval(src string) (*Reply, error) { return p.Call("eval", src) }

// Deparse parses the text and deparses each top-level expression.
func (p *Process) Deparse(src string) (*Reply, error) { return p.Call("deparse", src) }

// ParseData returns getParseData() of the parse of the text, as a data frame.
func (p *Process) ParseData(src string) (*Reply, error) { return p.Call("parsedata", src) }

// Reserialize has R unserialize the bytes and serialize the object again,
// and returns the bytes R wrote.
func (p *Process) Reserialize(b []byte) ([]byte, error) {
	r, err := p.Call("reserialize", b)
	if err != nil {
		return nil, err
	}
	raw, ok := r.Value.(*rds.Raw)
	if !ok {
		return nil, errors.New("oracle: reserialize did not return a raw vector")
	}
	return raw.V, nil
}

// cappedBuffer keeps the first max bytes written to it; it is the stderr of R.
type cappedBuffer struct {
	mu  sync.Mutex
	max int
	buf bytes.Buffer
}

func (c *cappedBuffer) Write(b []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if room := c.max - c.buf.Len(); room > 0 {
		c.buf.Write(b[:min(room, len(b))])
	}
	return len(b), nil
}

func (c *cappedBuffer) String() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.buf.String()
}
