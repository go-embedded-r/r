package cleanroom

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/go-embedded-r/r/syntax"
)

// The sequence of hashes of Figure 2(d) of Schleimer, Wilkerson, Aiken (2003)
// with window 4: the paper's Figure 2(g) selects the positions below (the
// paper numbers from 0).
func TestWinnowPaperFigure2(t *testing.T) {
	h := []uint64{77, 74, 42, 17, 98, 50, 17, 98, 8, 88, 67, 39, 77, 74, 42, 17, 98}
	got := Winnow(h, 4)
	want := []int{3, 6, 8, 11, 15}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Winnow = %v, want %v", got, want)
	}
}

func TestWinnowEdges(t *testing.T) {
	if Winnow(nil, 4) != nil {
		t.Fatal("empty input selects nothing")
	}
	if got := Winnow([]uint64{5, 3}, 4); !reflect.DeepEqual(got, []int{1}) {
		t.Fatalf("short input is one window: %v", got)
	}
	if got := Winnow([]uint64{1, 1, 1}, 2); !reflect.DeepEqual(got, []int{1, 2}) {
		t.Fatalf("ties select the rightmost: %v", got)
	}
	if got := Winnow([]uint64{9, 8, 7}, 0); len(got) != 3 {
		t.Fatalf("w < 1 means 1: %v", got)
	}
}

const sample = `
total <- function(values, weights = NULL, na.rm = FALSE) {
  if (is.null(weights)) weights <- rep(1, length(values))
  keep <- if (na.rm) !is.na(values) else rep(TRUE, length(values))
  sum(values[keep] * weights[keep]) / sum(weights[keep])
}
`

// renamed is sample with every variable renamed and the comment changed.
const renamed = `
tot2 <- function(xs, ws = NULL, drop = FALSE) {
  # comment
  if (is.null(ws)) ws <- rep(1, length(xs))
  ok <- if (drop) !is.na(xs) else rep(TRUE, length(xs))
  sum(xs[ok] * ws[ok]) / sum(ws[ok])
}
`

func TestRenamedCopyIsCaughtByTheNormalisedView(t *testing.T) {
	p := Params{K: 5, W: 3}
	a, b := NewRef("a.R", []byte(sample), p), NewRef("b.R", []byte(renamed), p)
	if c := Containment(a.FNorm, b.FNorm); c != 1 {
		t.Fatalf("normalised containment of a renamed copy = %v, want 1", c)
	}
	if c := Containment(a.FExact, b.FExact); c >= 0.9 {
		t.Fatalf("exact containment of a renamed copy = %v, want well below 1", c)
	}
	if r := LongestRun(a.Stream, b.Stream, true, p.K); r.Len < 30 || r.OurStart != 2 {
		t.Fatalf("normalised run = %+v", r)
	}
	if r := LongestRun(a.Stream, b.Stream, false, p.K); r.Len > 12 {
		t.Fatalf("exact run of a renamed copy = %+v", r)
	}
	if j := Jaccard(a.FNorm, a.FNorm); j != 1 {
		t.Fatalf("self Jaccard = %v", j)
	}
}

func TestDataTokensDoNotCountInARun(t *testing.T) {
	data := `x <- c("a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l")` + "\n"
	p := Params{K: 5, W: 3}
	a := NewRef("a.R", []byte(data), p)
	r := LongestRun(a.Stream, a.Stream, false, p.K)
	if r.Raw != len(a.Exact) || r.Len > 6 {
		t.Fatalf("run = %+v of %d tokens", r, len(a.Exact))
	}
}

func TestPercentile(t *testing.T) {
	xs := []float64{5, 1, 4, 2, 3}
	if Percentile(xs, 50) != 3 || Percentile(xs, 100) != 5 || Percentile(xs, 0) != 1 || Percentile(nil, 50) != 0 {
		t.Fatal("percentile")
	}
	if xs[0] != 5 {
		t.Fatal("input reordered")
	}
	if s := Summary([]Result{{Tokens: 10}, {Tokens: 300, Norm: 0.5}}, 200); !strings.Contains(s, "1 files") {
		t.Fatal(s)
	}
}

func TestViolationsReportReferencesOnly(t *testing.T) {
	r := Result{File: "f.R", Tokens: 500, Norm: 0.9, Best: "x.R", BestExact: 0.5, BestExactFile: "y.R",
		LongestAny: Run{Len: 50, OurStart: 1, OurEnd: 3, TheirStart: 7, TheirEnd: 9}, LongestAnyFile: "y.R",
		LongestAnyNorm: Run{Len: 70}, LongestAnyNormFile: "z.R"}
	if v := Gate.Violations(r); len(v) != 4 {
		t.Fatalf("violations = %v", v)
	}
	small := Result{File: "s.R", Tokens: 100, BestExactHits: 8, BestHitsFile: "x.R"}
	if v := Gate.Violations(small); len(v) != 1 {
		t.Fatalf("small file sharing 8 fingerprints: %v", v)
	}
	small.BestExactHits = 7
	if v := Gate.Violations(small); len(v) != 0 {
		t.Fatalf("small file sharing 7 fingerprints: %v", v)
	}
	if v := Gate.Violations(Result{}); len(v) != 0 {
		t.Fatalf("clean result flagged: %v", v)
	}
}

func rTree(t *testing.T) string {
	t.Helper()
	dir, ok := RSourceDir()
	if !ok {
		if os.Getenv("RGO_REQUIRE_RSOURCE") != "" {
			t.Fatal("RGO_REQUIRE_RSOURCE is set but no R source tree (set RGO_RSOURCE)")
		}
		t.Skip("no R source tree (set RGO_RSOURCE to R-4.6.2/)")
	}
	return dir
}

// renameAll rewrites src with every identifier renamed: what a careless
// "restructuring" of copied code amounts to.
func renameAll(src []byte) []byte {
	f, _ := syntax.ParseAll(string(src))
	var out []byte
	last := 0
	for _, tk := range f.Tokens {
		if strings.HasPrefix(tk.Kind, "SYMBOL") && tk.EndByte <= len(src) && tk.StartByte >= last {
			out = append(out, src[last:tk.StartByte]...)
			out = append(out, "zz_"...)
			out = append(out, tk.Text...)
			last = tk.EndByte
		}
	}
	return append(out, src[last:]...)
}

// The planted positive: a file of R's own with its identifiers renamed must
// fail the gate (otherwise the gate below passes for the wrong reason).
func TestGateCatchesARenamedCopyOfRSource(t *testing.T) {
	root := rTree(t)
	refs, err := LoadRefs(root, Default)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"src/library/base/R/ifelse.R", "src/library/base/R/sets.R", "src/library/base/R/merge.R", "src/library/base/R/strwrap.R"} {
		b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
		if err != nil {
			t.Fatal(err)
		}
		r := Check("control.R", renameAll(b), refs, Default)
		if len(Gate.Violations(r)) == 0 {
			t.Errorf("%s with renamed identifiers passed the gate: %+v", name, r)
		}
	}
}

// The gate: no file of eval/baselib is closer to R's library code than the
// noise floor measured in README ("Clean-room process") allows.
func TestBaselibIsCleanRoom(t *testing.T) {
	root := rTree(t)
	refs, err := LoadRefs(root, Default)
	if err != nil {
		t.Fatal(err)
	}
	files, _ := filepath.Glob(filepath.Join("..", "..", "eval", "baselib", "*.R"))
	if len(files) == 0 {
		t.Fatal("no eval/baselib/*.R")
	}
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		r := Check(filepath.ToSlash(filepath.Join("eval", "baselib", filepath.Base(f))), b, refs, Default)
		for _, v := range Gate.Violations(r) {
			t.Errorf("%s: %s", r.File, v)
		}
	}
}

func FuzzFingerprint(f *testing.F) {
	f.Add(sample, 5, 3)
	f.Add("", 1, 1)
	f.Add("x <- \"unterminated", 8, 4)
	f.Fuzz(func(t *testing.T, src string, k, w int) {
		if k > 64 || w > 64 {
			return
		}
		p := Params{K: k, W: w}
		a := NewRef("a.R", []byte(src), p)
		if len(a.FNorm) > 0 && Containment(a.FNorm, a.FNorm) != 1 {
			t.Fatal("self containment")
		}
		r := LongestRun(a.Stream, a.Stream, true, k)
		if r.Raw > len(a.Norm) || r.Len > r.Raw {
			t.Fatalf("run %+v of %d", r, len(a.Norm))
		}
	})
}
