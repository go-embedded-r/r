package oracle

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-embedded-r/r/rds"
)

// The test binary doubles as a fake Rscript (RGO_FAKE selects its behaviour),
// so the protocol and its failure modes are tested without R.
func TestMain(m *testing.M) {
	if mode := os.Getenv("RGO_FAKE"); mode != "" {
		fakeR(mode)
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func named(kv ...any) *rds.List {
	l := &rds.List{}
	names := &rds.Character{}
	for i := 0; i < len(kv); i += 2 {
		names.V = append(names.V, rds.NewStr(kv[i].(string)))
		l.V = append(l.V, kv[i+1].(rds.Value))
	}
	l.Attrs = []rds.Attr{{Name: "names", Value: names}}
	return l
}

func lgl(b bool) rds.Value {
	if b {
		return &rds.Logical{V: []int32{1}}
	}
	return &rds.Logical{V: []int32{0}}
}

func emit(v rds.Value) {
	b, _ := rds.Encode(v)
	fmt.Println(hex.EncodeToString(b))
}

func fakeR(mode string) {
	version := Version
	if mode == "oldr" {
		version = "4.5.3"
	}
	if mode == "nogreet" {
		return
	}
	if mode == "x" {
		return
	}
	if mode == "badhexgreet" {
		fmt.Println("not hex")
		return
	}
	if mode == "notlistgreet" {
		emit(&rds.Integer{V: []int32{1}})
		return
	}
	emit(named("ok", lgl(true), "ready", lgl(true), "version", rds.Strings(version),
		"svn", rds.Strings(SVNRevision), "platform", rds.Strings("fake"), "banner", rds.Strings("fake R")))
	in := bufio.NewReader(os.Stdin)
	// the bootstrap frame: the script, which the fake ignores
	var sn [4]byte
	if _, err := io.ReadFull(in, sn[:]); err != nil {
		return
	}
	io.CopyN(io.Discard, in, int64(binary.BigEndian.Uint32(sn[:])))
	for {
		var n [4]byte
		if _, err := io.ReadFull(in, n[:]); err != nil {
			return
		}
		body := make([]byte, binary.BigEndian.Uint32(n[:]))
		if _, err := io.ReadFull(in, body); err != nil {
			return
		}
		req, _, err := rds.Decode(body, rds.Limits{})
		if err != nil {
			fmt.Fprintln(os.Stderr, "fake: bad request:", err)
			os.Exit(3)
		}
		op := req.(*rds.List).V[0].(*rds.Character).V[0].S
		switch {
		case mode == "die":
			fmt.Fprintln(os.Stderr, "fake: dying")
			os.Exit(2)
		case mode == "garbage":
			fmt.Println("zz")
		case mode == "notlist":
			emit(&rds.Integer{V: []int32{1}})
		case mode == "nonames":
			emit(&rds.List{})
		case mode == "hang":
			time.Sleep(time.Minute)
		case op == "fail":
			emit(named("ok", lgl(false), "class", rds.Strings("simpleError", "error"), "message", rds.Strings("boom")))
		case op == "echo":
			emit(named("ok", lgl(true), "value", req.(*rds.List).V[1], "output", rds.Strings("line"), "warnings", rds.Strings("w1", "w2")))
		case op == "reserialize":
			emit(named("ok", lgl(true), "value", &rds.Integer{V: []int32{1}}))
		default:
			emit(named("ok", lgl(true)))
		}
	}
}

func fake(t *testing.T, mode string, opt Options) (*Process, error) {
	t.Helper()
	t.Setenv("RGO_FAKE", mode)
	t.Setenv("RGO_RSCRIPT", os.Args[0])
	return Start(opt)
}

func TestFakeProtocol(t *testing.T) {
	p, err := fake(t, "good", Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if g := p.Greeting(); g.Version != Version || g.Platform != "fake" || g.Banner != "fake R" {
		t.Errorf("greeting %+v", g)
	}
	r, err := p.Call("echo", "text", []byte{1, 2}, rds.Null)
	if err != nil {
		t.Fatal(err)
	}
	if got := r.Value.(*rds.Character).V[0].S; got != "text" {
		t.Errorf("value %q", got)
	}
	if strings.Join(r.Output, "|") != "line" || strings.Join(r.Warnings, "|") != "w1|w2" {
		t.Errorf("%+v", r)
	}
	r, err = p.Call("fail")
	if !errors.Is(err, ErrR) || r == nil || r.Message != "boom" || strings.Join(r.Class, "/") != "simpleError/error" {
		t.Errorf("failure: %v %+v", err, r)
	}
	if _, err := p.Call("echo", 3.5); err == nil || !strings.Contains(err.Error(), "unsupported argument") {
		t.Errorf("bad argument: %v", err)
	}
	if _, err := p.Call("ping"); err != nil {
		t.Error(err)
	}
	if r, err := p.Eval("1"); err != nil || r == nil {
		t.Error(err)
	}
	if _, err := p.Deparse("1"); err != nil {
		t.Error(err)
	}
	if _, err := p.ParseData("1"); err != nil {
		t.Error(err)
	}
	if b, err := p.Reserialize([]byte{1}); err == nil {
		t.Errorf("reserialize accepted a non-raw value: %v", b)
	}
}

func TestFakeStartFailures(t *testing.T) {
	for mode, want := range map[string]string{
		"oldr":         "pinned release",
		"nogreet":      "no greeting",
		"badhexgreet":  "malformed reply",
		"notlistgreet": "not a list",
	} {
		_, err := fake(t, mode, Options{})
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: %v, want %q", mode, err, want)
		}
		if got, want := Unavailable(err), mode == "oldr"; got != want {
			t.Errorf("%s: Unavailable = %v", mode, got)
		}
	}
	t.Setenv("RGO_RSCRIPT", "/nonexistent/Rscript")
	if _, err := Start(Options{}); err == nil || !strings.Contains(err.Error(), "starting") {
		t.Errorf("missing binary: %v", err)
	}
}

func TestFakeBrokenSessions(t *testing.T) {
	for mode, want := range map[string]string{
		"garbage": "malformed reply",
		"notlist": "not a list",
		"nonames": "no names",
		"die":     "fake: dying",
	} {
		p, err := fake(t, mode, Options{})
		if err != nil {
			t.Fatal(err)
		}
		_, err = p.Call("ping")
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: %v, want %q", mode, err, want)
		}
		if _, err := p.Call("ping"); err == nil || !strings.Contains(err.Error(), "gone") {
			t.Errorf("%s: second call: %v", mode, err)
		}
		p.Close()
	}
}

func TestFakeTimeout(t *testing.T) {
	p, err := fake(t, "hang", Options{Timeout: 300 * time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	start := time.Now()
	if _, err := p.Call("ping"); err == nil {
		t.Fatal("a hung R answered")
	}
	if d := time.Since(start); d > 10*time.Second {
		t.Errorf("timeout took %v", d)
	}
}

func TestCommand(t *testing.T) {
	t.Setenv("RGO_RSCRIPT", "/x/Rscript")
	if c, err := Command(); err != nil || c[0] != "/x/Rscript" {
		t.Errorf("%v %v", c, err)
	}
	t.Setenv("RGO_RSCRIPT", "")
	t.Setenv("PATH", t.TempDir())
	if _, err := Command(); !errors.Is(err, ErrNoR) {
		t.Errorf("empty PATH: %v", err)
	}
	if _, err := Start(Options{}); !errors.Is(err, ErrNoR) {
		t.Errorf("Start without R: %v", err)
	}
	if !strings.Contains(Script(), "serialize") {
		t.Error("embedded script is empty")
	}
}

func TestCappedBuffer(t *testing.T) {
	c := &cappedBuffer{max: 4}
	c.Write([]byte("abc"))
	c.Write([]byte("defg"))
	if n, _ := c.Write([]byte("more")); n != 4 || c.String() != "abcd" {
		t.Errorf("%q", c.String())
	}
}

func TestParseReplyErrors(t *testing.T) {
	if _, err := parseReply([]byte("junk")); err == nil {
		t.Error("junk accepted")
	}
	b, _ := rds.Encode(rds.Strings("a"))
	if _, err := parseReply(b); err == nil || !strings.Contains(err.Error(), "not a list") {
		t.Errorf("%v", err)
	}
	if got := strs(rds.Null); got != nil {
		t.Errorf("strs(NULL) = %v", got)
	}
	r := &Reply{}
	if r.String("x") != "" {
		t.Error("missing field")
	}
}

// ---- live tests: they need the pinned R and skip without it --------------

func live(t *testing.T) *Process {
	t.Helper()
	p, err := Start(Options{Timeout: 120 * time.Second})
	if Unavailable(err) {
		if os.Getenv("RGO_REQUIRE_ORACLE") != "" {
			t.Fatal(err)
		}
		t.Skip(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { p.Close() })
	return p
}

func TestLiveGreeting(t *testing.T) {
	p := live(t)
	g := p.Greeting()
	if g.Version != Version || g.SVN != SVNRevision || !strings.Contains(g.Banner, Version) {
		t.Errorf("%+v", g)
	}
}

func TestLiveEval(t *testing.T) {
	p := live(t)
	r, err := p.Eval("cat('hi\\n'); warning('careful'); c(1L, NA)")
	if err != nil {
		t.Fatal(err)
	}
	if got := r.Value.(*rds.Integer).V; len(got) != 2 || got[0] != 1 || got[1] != rds.NAInteger {
		t.Errorf("value %v", got)
	}
	if len(r.Output) != 1 || r.Output[0] != "hi" || len(r.Warnings) != 1 || r.Warnings[0] != "careful" {
		t.Errorf("output %q warnings %q", r.Output, r.Warnings)
	}

	r, err = p.Eval("stop('nope')")
	if !errors.Is(err, ErrR) || r.Message != "nope" || r.Class[0] != "simpleError" {
		t.Errorf("error reply: %v %+v", err, r)
	}
	if _, err = p.Eval("1 +"); !errors.Is(err, ErrR) {
		t.Errorf("parse error: %v", err)
	}
	// the session survives errors, and each request has its own environment
	if _, err := p.Eval("leaked <- 1"); err != nil {
		t.Fatal(err)
	}
	if _, err := p.Eval("leaked"); !errors.Is(err, ErrR) {
		t.Errorf("a variable leaked between requests: %v", err)
	}
	if _, err := p.Call("no-such-op"); !errors.Is(err, ErrR) {
		t.Errorf("unknown op: %v", err)
	}
}

func TestLiveParseTools(t *testing.T) {
	p := live(t)
	r, err := p.Deparse("a+b; f(x,y=1)")
	if err != nil {
		t.Fatal(err)
	}
	l := r.Value.(*rds.List)
	if len(l.V) != 2 || l.V[0].(*rds.Character).V[0].S != "a + b" {
		t.Errorf("deparse: %+v", l)
	}
	r, err = p.ParseData("x <- 1")
	if err != nil {
		t.Fatal(err)
	}
	df := r.Value.(*rds.List)
	names := df.Attr("names").(*rds.Character)
	var tokens []string
	for i, n := range names.V {
		if n.S == "token" {
			for _, s := range df.V[i].(*rds.Character).V {
				tokens = append(tokens, s.S)
			}
		}
	}
	if got := strings.Join(tokens, " "); !strings.Contains(got, "LEFT_ASSIGN") || !strings.Contains(got, "NUM_CONST") {
		t.Errorf("tokens %q", got)
	}
}

// R reads our serialization and writes back what Go writes for the same
// object, for every value of the golden set.
func TestLiveReserialize(t *testing.T) {
	p := live(t)
	recs, err := readGolden("../../testdata/golden/rds.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	for _, rec := range recs {
		raw, _ := hex.DecodeString(rec.Hex)
		got, err := p.Reserialize(raw)
		if err != nil {
			t.Errorf("%s: %v", rec.Src, err)
			continue
		}
		if !bytes.Equal(got, raw) {
			t.Errorf("%s: R rewrote the bytes\n got %x\nwant %x", rec.Src, got, raw)
		}
		v, _, err := rds.Decode(raw, rds.Limits{})
		if err != nil {
			t.Fatal(err)
		}
		goBytes, err := rds.Encode(v)
		if err != nil {
			t.Fatal(err)
		}
		if got, err := p.Reserialize(goBytes); err != nil || !bytes.Equal(got, goBytes) {
			t.Errorf("%s: R does not accept Go's bytes unchanged: %v", rec.Src, err)
		}
	}
}

func TestLiveTimeoutKillsR(t *testing.T) {
	p, err := Start(Options{Timeout: 700 * time.Millisecond})
	if Unavailable(err) {
		t.Skip(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	start := time.Now()
	if _, err := p.Call("sleep", "30"); err == nil {
		t.Fatal("sleep(30) returned within the timeout")
	}
	if d := time.Since(start); d > 20*time.Second {
		t.Errorf("took %v", d)
	}
	if _, err := p.Eval("1"); err == nil {
		t.Error("session usable after a timeout")
	}
}

func TestLiveMemoryLimit(t *testing.T) {
	p, err := Start(Options{MaxVSize: "100M", Timeout: 60 * time.Second})
	if Unavailable(err) {
		t.Skip(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	r, err := p.Eval("numeric(5e8)")
	if !errors.Is(err, ErrR) || !strings.Contains(r.Message, "vector memory limit") {
		t.Errorf("allocation not refused: %v %+v", err, r)
	}
	if _, err := p.Eval("1"); err != nil {
		t.Errorf("session lost after the refusal: %v", err)
	}
}

// ---- the clean-room guard: the oracle never shows R's R source --------------

const refusal = "oracle: printing base closures is disallowed by the clean-room rule"

// attempts that must be refused, from the plain to the roundabout
var guardAttempts = []string{
	"errorCondition", // the value itself
	"print(errorCondition)",
	"deparse(errorCondition)",
	"body(errorCondition)",
	"dput(errorCondition)",
	"as.list(errorCondition)",
	"format(errorCondition)",
	"toString(errorCondition)",
	"utils::capture.output(print(errorCondition))",
	"utils::capture.output(errorCondition)",
	"base::print.default(errorCondition)",
	"print.function(errorCondition)",
	"serialize(errorCondition, NULL)",
	"stats::sd",                                                   // another package
	"stats:::print.density",                                       // :::
	"list(1, list(2, errorCondition))",                            // inside a list
	"structure(list(), f = errorCondition)",                       // inside an attribute
	"f <- errorCondition; environment(f) <- globalenv(); body(f)", // a copy outside the namespace
	"f <- errorCondition; environment(f) <- globalenv(); f",
	"utils::getAnywhere('errorCondition')",
	"get('errorCondition', baseenv())",
	"match.fun('errorCondition')",
	"Negate(is.na)", // a closure built by base
	"methods::show(errorCondition)",
	"utils::str(errorCondition)",
	".Internal(print.function(errorCondition, FALSE, NULL))",
	"do.call('deparse', list(errorCondition))",
	"Reduce(function(a, b) b, list(1, errorCondition))",
}

func TestLiveBaseClosuresAreRefused(t *testing.T) {
	p := live(t)
	for _, src := range guardAttempts {
		r, err := p.Eval(src)
		if err == nil {
			t.Errorf("%q was answered: %+v", src, r)
			continue
		}
		if r == nil || r.Message != refusal {
			t.Errorf("%q: %v", src, err)
		}
	}
	// the console's auto-print path (what the goldens use)
	for _, src := range []string{"errorCondition", "stats::sd", "x <- errorCondition; x", "print(errorCondition)"} {
		r, err := p.Call("evalcase", src)
		if err != nil {
			t.Fatal(err)
		}
		l := r.Value.(*rds.List)
		names := l.Attr("names").(*rds.Character)
		got := map[string]rds.Value{}
		for i, n := range names.V {
			got[n.S] = l.V[i]
		}
		if m, _ := got["message"].(*rds.Character); m == nil || len(m.V) != 1 || m.V[0].S != refusal {
			t.Errorf("evalcase %q: %v", src, got)
		}
		if o, _ := got["output"].(*rds.Character); o != nil && len(o.V) > 0 && strings.Contains(o.V[0].S, "function") {
			t.Errorf("evalcase %q printed %q", src, o.V[0].S)
		}
	}
}

func TestLiveOwnClosuresStillPrint(t *testing.T) {
	p := live(t)
	r, err := p.Eval("f <- function(x, y = 2) x + y; utils::capture.output(print(f)); deparse(f); body(f); format(f); as.list(f); identity")
	_ = r
	if err == nil || r.Message != refusal {
		t.Fatalf("identity is a base closure and must be refused: %v", err)
	}
	r, err = p.Eval("f <- function(x, y = 2) x + y; g <- f; paste(utils::capture.output(print(f)), deparse(g), format(body(f)), length(as.list(f)), sep = '|')")
	if err != nil {
		t.Fatal(err)
	}
	if s := r.Value.(*rds.Character).V[0].S; !strings.Contains(s, "x + y") {
		t.Errorf("a user closure must still print: %q", s)
	}
	// non-closures through the same functions are untouched
	r, err = p.Eval("paste(deparse(quote(a + b)), format(1:3, width = 3), toString(1:3), body(function() 1), sep = ';')")
	if err != nil {
		t.Fatal(err)
	}
}
