// Package cleanroom measures how close our R files (eval/baselib/*.R) are to
// the R code of R's own library, so that a rewrite that "came out close to
// R's text" is caught by a gate instead of by memory.
//
// Method: both files are cut into tokens by this repository's own lexer
// (package syntax; comments dropped), the tokens are cut into k-grams,
// the k-gram hashes are thinned by winnowing (Schleimer, Wilkerson, Aiken,
// "Winnowing: Local Algorithms for Document Fingerprinting", SIGMOD 2003,
// section 3, Definition 1: in each window of w consecutive hashes keep the
// minimum, the rightmost on a tie; any shared token run of at least
// t = w+k-1 tokens shares a fingerprint, and no run shorter than k does).
//
// Two token streams are compared: Exact (identifiers kept) and Normalised
// (every identifier replaced by one placeholder, which catches a copy whose
// variables were renamed). For a pair (ours, theirs) the score is the
// containment |F(ours) ∩ F(theirs)| / |F(ours)|: the fraction of our
// fingerprints that also occur in theirs. The longest common token run is
// reported with line numbers on both sides.
//
// The R sources are GPL: this package reads them at check time only and
// reports file:line references, never their text.
package cleanroom

import (
	"fmt"
	"hash/fnv"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-embedded-r/r/syntax"
)

// Params are the winnowing parameters. K is the noise threshold (shorter
// matches are never detected), W the window size; the guarantee threshold is
// K+W-1 tokens.
type Params struct{ K, W int }

// Default parameters, chosen by measurement (README, "Clean-room process").
var Default = Params{K: 8, W: 4}

// Stream is one tokenised file in both views.
type Stream struct {
	Path  string
	Lines []int    // source line of every token
	Data  []bool   // token is a literal or a comma
	Exact []string // token text (strings and numbers included)
	Norm  []string // identifiers replaced by "ID"
}

// Tokenise lexes src with the repository's lexer. A syntax error keeps the
// tokens before it. Formal argument lists are dropped (see below).
func Tokenise(path string, src []byte) *Stream {
	f, _ := syntax.ParseAll(string(src))
	s := &Stream{Path: path}
	if f == nil {
		return s
	}
	// A formal argument list is fixed by the documented signature of the
	// function, so independent code has to repeat it: its tokens (the
	// parentheses and everything between them, defaults included) are left
	// out of both views, and the function header is `function` alone.
	afterFunction, depth := false, 0
	for _, t := range f.Tokens {
		if t.Kind == "COMMENT" {
			continue
		}
		if depth > 0 {
			switch t.Text {
			case "(":
				depth++
			case ")":
				depth--
			}
			continue
		}
		if afterFunction && t.Text == "(" {
			afterFunction, depth = false, 1
			continue
		}
		afterFunction = t.Text == "function" || t.Text == `\`
		s.Lines = append(s.Lines, t.Start.Line)
		s.Data = append(s.Data, t.Kind == "STR_CONST" || t.Kind == "NUM_CONST" || t.Kind == "NULL_CONST" || t.Text == ",")
		s.Exact = append(s.Exact, t.Text)
		if strings.HasPrefix(t.Kind, "SYMBOL") || t.Kind == "SLOT" {
			s.Norm = append(s.Norm, "ID")
		} else {
			s.Norm = append(s.Norm, t.Text)
		}
	}
	return s
}

func kgramHashes(toks []string, k int) []uint64 {
	if k < 1 || len(toks) < k {
		return nil
	}
	h := make([]uint64, len(toks)-k+1)
	for i := range h {
		f := fnv.New64a()
		for _, t := range toks[i : i+k] {
			f.Write([]byte(t))
			f.Write([]byte{0})
		}
		h[i] = f.Sum64()
	}
	return h
}

// Fingerprints winnows the k-gram hashes of toks: the set of selected hashes.
func Fingerprints(toks []string, p Params) map[uint64]struct{} {
	h := kgramHashes(toks, p.K)
	out := map[uint64]struct{}{}
	for _, i := range Winnow(h, p.W) {
		out[h[i]] = struct{}{}
	}
	return out
}

// Winnow returns the positions selected by the winnowing of h with window w
// (Definition 1 of the paper): in each window of w consecutive hashes, the
// minimum, the rightmost one on a tie; each position once, in order. A
// sequence shorter than w is one window.
func Winnow(h []uint64, w int) []int {
	if w < 1 {
		w = 1
	}
	if len(h) == 0 {
		return nil
	}
	if len(h) < w {
		w = len(h)
	}
	var sel []int
	last := -1
	for i := 0; i+w <= len(h); i++ {
		m := i
		for j := i + 1; j < i+w; j++ {
			if h[j] <= h[m] {
				m = j
			}
		}
		if m != last {
			sel = append(sel, m)
			last = m
		}
	}
	return sel
}

// Containment is |a ∩ b| / |a| (0 for an empty a).
func Containment(a, b map[uint64]struct{}) float64 {
	if len(a) == 0 {
		return 0
	}
	n := 0
	for h := range a {
		if _, ok := b[h]; ok {
			n++
		}
	}
	return float64(n) / float64(len(a))
}

// Hits is |a ∩ b|: the number of fingerprints two files share.
func Hits(a, b map[uint64]struct{}) int {
	n := 0
	for h := range a {
		if _, ok := b[h]; ok {
			n++
		}
	}
	return n
}

// Jaccard is |a ∩ b| / |a ∪ b|.
func Jaccard(a, b map[uint64]struct{}) float64 {
	if len(a)+len(b) == 0 {
		return 0
	}
	n := 0
	for h := range a {
		if _, ok := b[h]; ok {
			n++
		}
	}
	return float64(n) / float64(len(a)+len(b)-n)
}

// Run is a common token run. Len counts the code tokens of the run (not
// string or number literals, nor commas: a table of constants such as the
// month names is data, and independent code has to contain it too); Raw
// counts every token. The lines are those of the first and the last token.
type Run struct {
	Len, Raw             int
	OurStart, OurEnd     int // lines
	TheirStart, TheirEnd int
}

// LongestRun finds the longest common contiguous token run of a and b.
// Seeds are the k-gram matches; runs are extended both ways.
func LongestRun(a, b *Stream, norm bool, k int) Run {
	var best Run
	for _, r := range Runs(a, b, norm, k, 1) {
		if r.Len > best.Len {
			best = r
		}
	}
	return best
}

// Runs lists the maximal common runs of a and b that hold at least min code
// tokens, in the order of our file.
func Runs(a, b *Stream, norm bool, k, min int) []Run {
	ta, tb := a.Exact, b.Exact
	if norm {
		ta, tb = a.Norm, b.Norm
	}
	ha, hb := kgramHashes(ta, k), kgramHashes(tb, k)
	if len(ha) == 0 || len(hb) == 0 {
		return nil
	}
	pos := map[uint64][]int{}
	for j, h := range hb {
		if len(pos[h]) < 64 { // bound the work on a repeated k-gram
			pos[h] = append(pos[h], j)
		}
	}
	var out []Run
	for i, h := range ha {
		for _, j := range pos[h] {
			if i > 0 && j > 0 && ta[i-1] == tb[j-1] {
				continue // not the start of the run
			}
			n, code := 0, 0
			for i+n < len(ta) && j+n < len(tb) && ta[i+n] == tb[j+n] {
				if !a.Data[i+n] {
					code++
				}
				n++
			}
			if code >= min && n > 0 {
				out = append(out, Run{Len: code, Raw: n,
					OurStart: a.Lines[i], OurEnd: a.Lines[i+n-1],
					TheirStart: b.Lines[j], TheirEnd: b.Lines[j+n-1]})
			}
		}
	}
	return out
}

// Ref is one reference file, tokenised and fingerprinted once.
type Ref struct {
	*Stream
	FExact, FNorm map[uint64]struct{}
}

// NewRef tokenises and fingerprints a file.
func NewRef(path string, src []byte, p Params) *Ref {
	s := Tokenise(path, src)
	return &Ref{Stream: s, FExact: Fingerprints(s.Exact, p), FNorm: Fingerprints(s.Norm, p)}
}

// LoadRefs reads every .R file under root/src/library/*/R (recursively); the
// paths of the result are relative to root. A missing tree is an error.
func LoadRefs(root string, p Params) ([]*Ref, error) {
	libs, err := filepath.Glob(filepath.Join(root, "src", "library", "*", "R"))
	if err != nil || len(libs) == 0 {
		return nil, fmt.Errorf("no src/library/*/R under %s", root)
	}
	var refs []*Ref
	for _, dir := range libs {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(path, ".R") {
				return err
			}
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			rel, _ := filepath.Rel(root, path)
			refs = append(refs, NewRef(filepath.ToSlash(rel), b, p))
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	sort.Slice(refs, func(i, j int) bool { return refs[i].Path < refs[j].Path })
	return refs, nil
}

// Result is the verdict for one of our files: the best-matching R file by
// each measure, with both scores and the longest runs against that file.
type Result struct {
	File               string
	Tokens             int
	Best               string  // by Normalised containment
	Exact, Norm        float64 // containment vs Best
	ExactJac, NormJac  float64
	RunExact, RunNorm  Run // against Best
	BestExactFile      string
	BestExact          float64 // highest Exact containment over all refs
	BestExactHits      int     // most exact fingerprints shared with one ref
	BestHitsFile       string
	LongestAny         Run // longest exact run over all refs
	LongestAnyFile     string
	LongestAnyNorm     Run // longest normalised run over all refs
	LongestAnyNormFile string
}

// Check scores one of our files against all refs.
func Check(path string, src []byte, refs []*Ref, p Params) Result {
	ours := NewRef(path, src, p)
	r := Result{File: path, Tokens: len(ours.Exact)}
	var bestN float64 = -1
	for _, ref := range refs {
		if ref.Path == path {
			continue
		}
		cn := Containment(ours.FNorm, ref.FNorm)
		ce := Containment(ours.FExact, ref.FExact)
		if cn > bestN {
			bestN = cn
			r.Best, r.Norm, r.Exact = ref.Path, cn, ce
			r.NormJac, r.ExactJac = Jaccard(ours.FNorm, ref.FNorm), Jaccard(ours.FExact, ref.FExact)
		}
		if h := Hits(ours.FExact, ref.FExact); h > r.BestExactHits {
			r.BestExactHits, r.BestHitsFile = h, ref.Path
		}
		if ce > r.BestExact {
			r.BestExact, r.BestExactFile = ce, ref.Path
		}
	}
	// runs: seeded by every ref (cheap: hash lookups)
	for _, ref := range refs {
		if ref.Path == path {
			continue
		}
		if x := LongestRun(ours.Stream, ref.Stream, false, p.K); x.Len > r.LongestAny.Len {
			r.LongestAny, r.LongestAnyFile = x, ref.Path
		}
		if x := LongestRun(ours.Stream, ref.Stream, true, p.K); x.Len > r.LongestAnyNorm.Len {
			r.LongestAnyNorm, r.LongestAnyNormFile = x, ref.Path
		}
		if ref.Path == r.Best {
			r.RunExact = LongestRun(ours.Stream, ref.Stream, false, p.K)
			r.RunNorm = LongestRun(ours.Stream, ref.Stream, true, p.K)
		}
	}
	return r
}

// Limits is the gate: a file fails when its Normalised containment against
// one R file, or its longest exact or normalised run against any, reaches it.
type Limits struct {
	Containment      float64 // normalised, against one R file
	ExactContainment float64 // identifiers kept, against one R file
	SmallHits        int     // files of 30..MinTokens-1 tokens: shared exact fingerprints
	ExactRun         int     // code tokens
	NormRun          int
}

// ListRuns reports, one line each, every exact run of at least min code
// tokens that src shares with a reference file: the places to rewrite.
func ListRuns(path string, src []byte, refs []*Ref, p Params, min int) string {
	ours := Tokenise(path, src)
	var sb strings.Builder
	for _, ref := range refs {
		if ref.Path == path {
			continue
		}
		for _, r := range Runs(ours, ref.Stream, false, p.K, min) {
			fmt.Fprintf(&sb, "  run %d tok: %s:%d-%d vs %s:%d-%d\n", r.Len, path, r.OurStart, r.OurEnd, ref.Path, r.TheirStart, r.TheirEnd)
		}
	}
	return sb.String()
}

// Percentile is the nearest-rank p-th percentile (0..100) of xs (0 if empty).
func Percentile(xs []float64, p float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	s := append([]float64(nil), xs...)
	sort.Float64s(s)
	i := int(math.Ceil(p/100*float64(len(s)))) - 1
	if i < 0 {
		i = 0
	}
	return s[i]
}

// Summary describes the distribution of the scores of results with at least
// minTokens tokens: the noise floor when the inputs are independent code.
func Summary(res []Result, minTokens int) string {
	var ne, nn, re, rn []float64
	for _, r := range res {
		if r.Tokens < minTokens {
			continue
		}
		ne = append(ne, r.BestExact)
		nn = append(nn, r.Norm)
		re = append(re, float64(r.LongestAny.Len))
		rn = append(rn, float64(r.LongestAnyNorm.Len))
	}
	line := func(name string, xs []float64) string {
		return fmt.Sprintf("  %-26s p50 %.2f  p90 %.2f  p95 %.2f  p99 %.2f  max %.2f\n", name, Percentile(xs, 50), Percentile(xs, 90), Percentile(xs, 95), Percentile(xs, 99), Percentile(xs, 100))
	}
	return fmt.Sprintf("summary: %d files with >= %d tokens\n", len(nn), minTokens) +
		line("exact containment", ne) + line("normalised containment", nn) +
		line("exact run (code tokens)", re) + line("normalised run (code tok)", rn)
}

// SmallSummary describes the exact fingerprints a file of 30 to MinTokens-1
// tokens shares with one R file, over the results given: the floor of the
// small-file rule.
func SmallSummary(res []Result) string {
	var hits []float64
	for _, r := range res {
		if r.Tokens >= 30 && r.Tokens < MinTokens {
			hits = append(hits, float64(r.BestExactHits))
		}
	}
	return fmt.Sprintf("small files (30 <= tokens < %d): %d\n  shared exact fingerprints  p50 %.0f  p90 %.0f  p95 %.0f  p99 %.0f  max %.0f\n",
		MinTokens, len(hits), Percentile(hits, 50), Percentile(hits, 90), Percentile(hits, 95), Percentile(hits, 99), Percentile(hits, 100))
}

// MinTokens is the size from which the containment limits apply: the noise
// floor was measured on files of at least this many tokens, and a file of a
// hundred tokens has a dozen fingerprints, so one match moves it by 0.08. The
// run limits apply to every file.
const MinTokens = 200

// Gate is the limit the repository must stay under (README, "Clean-room
// process"). Each is above what independent code reached when measured:
// exact containment 0.20 (R's tests and demos: p99 0.19; our test inputs:
// max 0.04), normalised containment 0.60 (max 0.52), exact run 40 code
// tokens (p99 38; our test inputs max 14), normalised run 60 (p99 54), and
// for files of 30 to 199 tokens 4 shared exact fingerprints (independent
// small files: max 2 of 32). Formal argument lists are not fingerprinted, so
// the signatures that help pages fix do not count.
var Gate = Limits{Containment: 0.60, ExactContainment: 0.20, SmallHits: 4, ExactRun: 40, NormRun: 60}

// Violations lists why r fails the gate, as references (no R text).
func (l Limits) Violations(r Result) []string {
	var v []string
	big := r.Tokens >= MinTokens // a handful of fingerprints is not a proportion
	if big && r.Norm >= l.Containment {
		v = append(v, fmt.Sprintf("normalised containment %.2f >= %.2f against %s", r.Norm, l.Containment, r.Best))
	}
	if big && r.BestExact >= l.ExactContainment {
		v = append(v, fmt.Sprintf("exact containment %.2f >= %.2f against %s", r.BestExact, l.ExactContainment, r.BestExactFile))
	}
	if !big && r.Tokens >= 30 && r.BestExactHits >= l.SmallHits {
		v = append(v, fmt.Sprintf("%d exact fingerprints shared (small file) >= %d with %s", r.BestExactHits, l.SmallHits, r.BestHitsFile))
	}
	if r.LongestAny.Len >= l.ExactRun {
		v = append(v, fmt.Sprintf("exact run of %d tokens >= %d: ours %s:%d-%d, theirs %s:%d-%d", r.LongestAny.Len, l.ExactRun,
			r.File, r.LongestAny.OurStart, r.LongestAny.OurEnd, r.LongestAnyFile, r.LongestAny.TheirStart, r.LongestAny.TheirEnd))
	}
	if r.LongestAnyNorm.Len >= l.NormRun {
		v = append(v, fmt.Sprintf("normalised run of %d tokens >= %d: ours %s:%d-%d, theirs %s:%d-%d", r.LongestAnyNorm.Len, l.NormRun,
			r.File, r.LongestAnyNorm.OurStart, r.LongestAnyNorm.OurEnd, r.LongestAnyNormFile, r.LongestAnyNorm.TheirStart, r.LongestAnyNorm.TheirEnd))
	}
	return v
}

// RSourceDir finds the R source tree: $RGO_RSOURCE, else the developer's
// durable checkout. ok is false when there is none.
func RSourceDir() (string, bool) {
	cands := []string{os.Getenv("RGO_RSOURCE")}
	if home, err := os.UserHomeDir(); err == nil {
		cands = append(cands, filepath.Join(home, "Documents", "VCS", "GIT", "localhost", "r-source-R-4-6-2", "R-4.6.2"))
	}
	for _, c := range cands {
		if c == "" {
			continue
		}
		if st, err := os.Stat(filepath.Join(c, "src", "library", "base", "R")); err == nil && st.IsDir() {
			return c, true
		}
	}
	return "", false
}
