package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-embedded-r/r/internal/oracle"
	"github.com/go-embedded-r/r/rds"
)

// fakeSession answers every request with canned data.
type fakeSession struct {
	evalErr error
	closed  *bool
}

func (f fakeSession) Deparse(string) (*oracle.Reply, error) {
	return &oracle.Reply{OK: true, Value: &rds.List{V: []rds.Value{rds.Strings("x")}}}, nil
}

func (f fakeSession) Call(string, ...any) (*oracle.Reply, error) {
	return &oracle.Reply{OK: true, Value: &rds.List{V: []rds.Value{
		&rds.Raw{V: []byte{1}}, rds.Strings("raw"), &rds.Integer{V: []int32{1}}, rds.Strings("raw"),
	}}}, nil
}

func (f fakeSession) Eval(string) (*oracle.Reply, error) {
	return &oracle.Reply{OK: true, Output: []string{"out1", "out2"}, Warnings: []string{"warn"}}, f.evalErr
}

func (f fakeSession) Greeting() oracle.Greeting {
	return oracle.Greeting{Banner: "fake banner", Version: "4.6.2", SVN: "1", Platform: "plat"}
}

func (f fakeSession) Close() error {
	if f.closed != nil {
		*f.closed = true
	}
	return nil
}

func withFake(t *testing.T, s oracleSession, err error) {
	t.Helper()
	old := startOracle
	startOracle = func() (oracleSession, error) { return s, err }
	t.Cleanup(func() { startOracle = old })
}

func run(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	cmd := RootCmd()
	var out, errb bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errb)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return out.String(), errb.String(), err
}

func repoCopy(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, f := range []string{"testdata/snippets.R", "testdata/values.R", "testdata/errors.R", "testdata/edge.R"} {
		b, err := os.ReadFile(filepath.Join("../..", f))
		if err != nil {
			t.Fatal(err)
		}
		os.MkdirAll(filepath.Dir(filepath.Join(dir, f)), 0o755)
		os.WriteFile(filepath.Join(dir, f), b, 0o644)
	}
	return dir
}

func TestVersionAndEval(t *testing.T) {
	closed := false
	withFake(t, fakeSession{closed: &closed}, nil)
	out, _, err := run(t, "version")
	if err != nil || !strings.Contains(out, "fake banner") || !strings.Contains(out, "R 4.6.2 (r1) on plat") || !closed {
		t.Errorf("%q %v closed=%v", out, err, closed)
	}
	out, errOut, err := run(t, "eval", "1")
	if err != nil || out != "out1\nout2\n" || !strings.Contains(errOut, "warning: warn") {
		t.Errorf("%q %q %v", out, errOut, err)
	}
	withFake(t, fakeSession{evalErr: errors.New("R said no")}, nil)
	if _, _, err := run(t, "eval", "1"); err == nil {
		t.Error("evaluation error swallowed")
	}
	if _, _, err := run(t, "eval"); err == nil {
		t.Error("eval without an argument accepted")
	}
}

func TestStartFailure(t *testing.T) {
	withFake(t, nil, oracle.ErrNoR)
	for _, c := range []string{"version", "regen", "deparse", "rds"} {
		if _, _, err := run(t, c, "-C", t.TempDir()); !errors.Is(err, oracle.ErrNoR) {
			t.Errorf("%s: %v", c, err)
		}
	}
	if _, _, err := run(t, "eval", "1"); !errors.Is(err, oracle.ErrNoR) {
		t.Errorf("eval: %v", err)
	}
}

func TestRegenWritesGoldens(t *testing.T) {
	withFake(t, fakeSession{}, nil)
	dir := repoCopy(t)
	out, _, err := run(t, "regen", "-C", dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{"testdata/golden/deparse.jsonl", "testdata/golden/rds.jsonl", "testdata/golden/parse.jsonl"} {
		if !strings.Contains(out, "wrote "+f) {
			t.Errorf("no report for %s in %q", f, out)
		}
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Error(err)
		}
	}
}

func TestRegenFailures(t *testing.T) {
	withFake(t, fakeSession{}, nil)
	empty := t.TempDir() // no inputs
	for _, c := range []string{"deparse", "rds", "parse"} {
		if _, _, err := run(t, c, "-C", empty); err == nil {
			t.Errorf("%s: missing inputs accepted", c)
		}
	}
	// an unwritable target
	dir := repoCopy(t)
	os.WriteFile(filepath.Join(dir, "testdata/golden"), []byte("a file, not a directory"), 0o644)
	for _, c := range []string{"deparse", "rds"} {
		if _, _, err := run(t, c, "-C", dir); err == nil {
			t.Errorf("%s: unwritable target accepted", c)
		}
	}
	// the oracle returns something the generators cannot use
	withFake(t, badSession{}, nil)
	dir = repoCopy(t)
	for _, c := range []string{"deparse", "rds"} {
		if _, _, err := run(t, c, "-C", dir); err == nil {
			t.Errorf("%s: unusable answer accepted", c)
		}
	}
}

type badSession struct{ fakeSession }

func (badSession) Deparse(string) (*oracle.Reply, error) {
	return &oracle.Reply{OK: true, Value: rds.Null}, nil
}

func (badSession) Call(string, ...any) (*oracle.Reply, error) {
	return &oracle.Reply{OK: true, Value: rds.Null}, nil
}

func TestMainExitsNonZero(t *testing.T) {
	withFake(t, nil, oracle.ErrNoR)
	oldExit, oldArgs := osExit, os.Args
	defer func() { osExit, os.Args = oldExit, oldArgs }()
	code := -1
	osExit = func(c int) { code = c }
	os.Args = []string{"roracle", "version"}
	main()
	if code != 1 {
		t.Errorf("exit code %d", code)
	}
}

// The real thing, when R is available: regenerating into a scratch copy of
// the inputs gives the committed goldens.
func TestLiveRegen(t *testing.T) {
	old := startOracle
	defer func() { startOracle = old }()
	p, err := oracle.Start(oracle.Options{})
	if oracle.Unavailable(err) {
		if os.Getenv("RGO_REQUIRE_ORACLE") != "" {
			t.Fatal(err)
		}
		t.Skip(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	p.Close()
	startOracle = func() (oracleSession, error) { return oracle.Start(oracle.Options{}) }
	dir := repoCopy(t)
	if _, _, err := run(t, "regen", "-C", dir); err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{"testdata/golden/deparse.jsonl", "testdata/golden/rds.jsonl", "testdata/golden/parse.jsonl"} {
		got, _ := os.ReadFile(filepath.Join(dir, f))
		want, _ := os.ReadFile(filepath.Join("../..", f))
		if !bytes.Equal(got, want) {
			t.Errorf("%s differs from the committed golden", f)
		}
	}
	out, _, err := run(t, "version")
	if err != nil || !strings.Contains(out, "R 4.6.2 (r90784)") {
		t.Errorf("%q %v", out, err)
	}
}
