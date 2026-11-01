package gen

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/go-embedded-r/r/internal/oracle"
	"github.com/go-embedded-r/r/rds"
)

const root = "../.."

func TestSplit(t *testing.T) {
	in := "## header\n#---\na\nb\n#---\n\n#---\nc\r\n#---\n## not a snippet\n#---\nd\n"
	want := []string{"a\nb", "c", "d"}
	if got := Split(in); !reflect.DeepEqual(got, want) {
		t.Errorf("%q, want %q", got, want)
	}
}

func TestReadSnippets(t *testing.T) {
	if _, err := ReadSnippets(root, "nope.R"); err == nil {
		t.Error("missing file accepted")
	}
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "empty.R"), []byte("## only a header\n"), 0o644)
	if _, err := ReadSnippets(dir, "empty.R"); err == nil || !strings.Contains(err.Error(), "no snippet") {
		t.Errorf("empty file: %v", err)
	}
}

// The committed goldens are in step with the committed inputs: the same
// snippets, in the same order. No R is needed.
func TestGoldensMatchInputs(t *testing.T) {
	snippets, err := ReadSnippets(root, "testdata/snippets.R")
	if err != nil {
		t.Fatal(err)
	}
	dep, err := ReadJSONL[DeparseRecord](filepath.Join(root, "testdata/golden/deparse.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	if len(dep) != len(snippets) {
		t.Fatalf("%d deparse records for %d snippets", len(dep), len(snippets))
	}
	for i, r := range dep {
		if r.Src != snippets[i] {
			t.Errorf("deparse record %d is for %q, not %q", i, r.Src, snippets[i])
		}
		if (r.Error == "") == (len(r.Deparse) == 0) {
			t.Errorf("record %d has both or neither of deparse and error", i)
		}
	}
	values, err := ReadSnippets(root, "testdata/values.R")
	if err != nil {
		t.Fatal(err)
	}
	rdsRecs, err := ReadJSONL[RDSRecord](filepath.Join(root, "testdata/golden/rds.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	if len(rdsRecs) != len(values) {
		t.Fatalf("%d rds records for %d expressions", len(rdsRecs), len(values))
	}
	for i, r := range rdsRecs {
		if r.Src != values[i] {
			t.Errorf("rds record %d is for %q, not %q", i, r.Src, values[i])
		}
	}
}

func TestReadJSONL(t *testing.T) {
	if _, err := ReadJSONL[DeparseRecord]("nope"); err == nil {
		t.Error("missing file accepted")
	}
	p := filepath.Join(t.TempDir(), "x.jsonl")
	os.WriteFile(p, []byte(`{"src":"a","surprise":1}`+"\n"), 0o644)
	if _, err := ReadJSONL[DeparseRecord](p); err == nil {
		t.Error("unknown field accepted")
	}
}

// fakeOracle answers from canned replies.
type fakeOracle struct {
	deparse func(string) (*oracle.Reply, error)
	call    func(string, ...any) (*oracle.Reply, error)
}

func (f fakeOracle) Deparse(s string) (*oracle.Reply, error) { return f.deparse(s) }
func (f fakeOracle) Call(op string, a ...any) (*oracle.Reply, error) {
	return f.call(op, a...)
}

func TestGeneratorFailures(t *testing.T) {
	boom := errors.New("pipe closed")
	_, err := DeparseGolden(fakeOracle{deparse: func(string) (*oracle.Reply, error) { return nil, boom }}, []string{"x"})
	if !errors.Is(err, boom) {
		t.Errorf("transport error: %v", err)
	}
	for name, v := range map[string]rds.Value{
		"not a list":  rds.Null,
		"not strings": &rds.List{V: []rds.Value{rds.Null}},
	} {
		_, err := DeparseGolden(fakeOracle{deparse: func(string) (*oracle.Reply, error) { return &oracle.Reply{OK: true, Value: v}, nil }}, []string{"x"})
		if err == nil {
			t.Errorf("deparse %s accepted", name)
		}
	}
	_, err = RDSGolden(fakeOracle{call: func(string, ...any) (*oracle.Reply, error) { return nil, boom }}, []string{"x"})
	if !errors.Is(err, boom) {
		t.Errorf("rds transport error: %v", err)
	}
	for name, v := range map[string]rds.Value{
		"not a list":  rds.Null,
		"wrong shape": &rds.List{V: []rds.Value{rds.Null, rds.Null, rds.Null, rds.Null}},
	} {
		_, err := RDSGolden(fakeOracle{call: func(string, ...any) (*oracle.Reply, error) { return &oracle.Reply{OK: true, Value: v}, nil }}, []string{"x"})
		if err == nil {
			t.Errorf("rds %s accepted", name)
		}
	}
	// a parse error is data, not a failure
	out, err := DeparseGolden(fakeOracle{deparse: func(string) (*oracle.Reply, error) {
		return &oracle.Reply{Message: "unexpected end of input"}, oracle.ErrR
	}}, []string{"1 +"})
	if err != nil || !strings.Contains(string(out), `"error":"unexpected end of input"`) {
		t.Errorf("%s %v", out, err)
	}
}

// With the pinned R, regenerating both goldens reproduces the committed
// files byte for byte (CI runs this on Linux, the files were written on
// macOS).
func TestLiveRegeneration(t *testing.T) {
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
	defer p.Close()

	snippets, _ := ReadSnippets(root, "testdata/snippets.R")
	got, err := DeparseGolden(p, snippets)
	if err != nil {
		t.Fatal(err)
	}
	want, _ := os.ReadFile(filepath.Join(root, "testdata/golden/deparse.jsonl"))
	compareLines(t, "deparse.jsonl", got, want)

	values, _ := ReadSnippets(root, "testdata/values.R")
	got, err = RDSGolden(p, values)
	if err != nil {
		t.Fatal(err)
	}
	want, _ = os.ReadFile(filepath.Join(root, "testdata/golden/rds.jsonl"))
	compareLines(t, "rds.jsonl", got, want)
}

// compareLines reports which records of a golden differ, so that a platform
// difference names the input it concerns.
func compareLines(t *testing.T, name string, got, want []byte) {
	t.Helper()
	if bytes.Equal(got, want) {
		return
	}
	g := strings.Split(string(got), "\n")
	w := strings.Split(string(want), "\n")
	t.Errorf("%s is stale or platform-dependent (%d records regenerated, %d committed); regenerate with go run ./cmd/roracle regen", name, len(g), len(w))
	shown := 0
	for i := 0; i < len(g) && i < len(w) && shown < 5; i++ {
		if g[i] != w[i] {
			t.Errorf("record %d differs:\n regenerated %.300s\n committed   %.300s", i+1, g[i], w[i])
			shown++
		}
	}
}
