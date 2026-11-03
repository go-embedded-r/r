// Package parity compares the Go parser with the reference R 4.6.2 on
// arbitrary source text, at three levels: the syntax tree (against the dump
// of R's own language objects), the terminal tokens (against
// utils::getParseData) and the error message (against parse()'s).
package parity

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/go-embedded-r/r/internal/oracle"
	"github.com/go-embedded-r/r/rds"
	"github.com/go-embedded-r/r/syntax"
)

// Oracle is the part of the reference R the comparison uses.
type Oracle interface {
	Call(op string, args ...any) (*oracle.Reply, error)
}

// Outcome is the result of comparing one source text.
type Outcome struct {
	// Parsed is set when R parsed the text.
	Parsed bool
	// Tree, Tokens and Error are the three verdicts. For text R parses, Tree
	// and Tokens apply (Error is true when Go agrees it parses); for text R
	// rejects, Error applies (Tree and Tokens are true when Go agrees it is
	// rejected).
	Tree, Tokens, Error bool
	// Deparse is set when Go's deparse of every top-level expression equals R's, in
	// both ways R writes one: deparse() and print().
	Deparse bool
	// ULP counts the numeric literals that differ from R's by a few units in
	// the last place (see syntax's package documentation).
	ULP int
	// Detail describes the first difference.
	Detail string
}

// Extras are inputs that cannot be written in a text file of snippets:
// carriage returns, form feeds, odd spaces and control characters.
func Extras() []string {
	return []string{
		"x <- 1\r\ny <- 2",
		"x <- 1 \r\n",
		"a +\r\nb",
		"\tx <- 1",
		"a\t\t+ b",
		"x\t<-\t'é'\t# é",
		"\f x",
		"x \f+ y",
		"\x01",
		"x \x7f y",
		" x",
		"x + y",
		"x   y",
		" x",
		"é\tz",
		"\"é\tx\"; y z",
		"'a\nb' z",
		"x <- 'unterminated\ny <- 2",
		"x <- r\"(unterminated",
		"\u3000x", "\u2007x", "\u202fx", "\u1680x", "\u200bx", "x\u2009+ y",
		"x # comment é\n\ty",
		"1 +\n\n\n\n2",
		"if (a)\n\n\n b else\n\n\n c",
		"f(\n\n)",
		"{\n\n\n}",
		"x <- 'a' 'b'",
		"`a\nb`",
		"function(x)\n\n x",
		"x[\n1\n]",
		"c(1,\n\n2,\n\n)",
		"a ||\n\n b",
		"a\n|| b",
		"a\n+ b",
		"-\n1",
		"!\n\nx",
		"x$\n y",
		"x |>\n f()",
		"x\n|> f()",
		"f(length=>)", "a => b", "a => b => c", "x |> f() => y",
		"x <- '\xff'", "\xff", "x \xe9 y", "'\xe9'", "`a\xc3`", "r\"(\xc3)\"", "# \xff\nx", "x <- '\xe2\x82'", "\xe2\x82",
		"x <- '\u202e'", "'a\u2066b'", "r'(\u202a)'", "\u202e",
		"x <- 1;;", "if (a) b\n\n\nelse c", "{if (a) b\n\n\nelse c}", "f(if (a) b\n else c)",
	}
}

// GoDeparse deparses e the way the oracle's deparse2 operation does in the
// given mode ("default": deparse(); "print": print() of a call or name).
func GoDeparse(e syntax.Node, mode string) []string {
	var opt *syntax.DeparseOptions
	if _, isSym := e.(*syntax.Sym); isSym && mode == "default" {
		opt = &syntax.DeparseOptions{NoBackquote: true}
	}
	return syntax.Deparse(e, opt)
}

// TreeLines returns the dump of R's parse of src, or R's error message.
func TreeLines(o Oracle, src string) (lines []string, rErr string, err error) {
	r, err := o.Call("tree", src)
	if r != nil && !r.OK {
		return nil, r.Message, nil
	}
	if err != nil {
		return nil, "", err
	}
	return strs(r.Value), "", nil
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

// TokenLines formats the Go tokens like the oracle's "tokens" operation.
func TokenLines(toks []syntax.Token) []string {
	out := make([]string, 0, len(toks))
	sorted := append([]syntax.Token(nil), toks...)
	// R orders by (line1, col1)
	sort.SliceStable(sorted, func(i, j int) bool { return less(sorted[i], sorted[j]) })
	for _, t := range sorted {
		out = append(out, fmt.Sprintf("%d:%d-%d:%d %s %s", t.Start.Line, t.Start.Col, t.End.Line, t.End.Col, t.Kind, hex.EncodeToString([]byte(t.Text))))
	}
	return out
}

func less(a, b syntax.Token) bool {
	if a.Start.Line != b.Start.Line {
		return a.Start.Line < b.Start.Line
	}
	return a.Start.Col < b.Start.Col
}

// Compare parses src with both and compares them.
func Compare(o Oracle, src string) (Outcome, error) {
	var out Outcome
	rlines, rErr, err := TreeLines(o, src)
	if err != nil {
		return out, err
	}
	f, gerr := syntax.ParseAll(src)
	out.Parsed = rErr == ""
	if !out.Parsed {
		out.Tree, out.Tokens, out.Deparse = gerr != nil, gerr != nil, gerr != nil
		switch {
		case gerr == nil:
			out.Detail = "R rejects: " + firstLine(rErr) + "; Go accepts"
		case gerr.Error() == rErr:
			out.Error = true
		default:
			out.Detail = "error text\n  R:  " + rErr + "\n  Go: " + gerr.Error()
		}
		return out, nil
	}
	if gerr != nil {
		out.Detail = "R parses; Go says " + gerr.Error()
		return out, nil
	}
	out.Error = true
	out.Deparse = true
	for _, mode := range []string{"default", "print"} {
		r, err := o.Call("deparse2", src, mode)
		if err != nil {
			return out, err
		}
		rl, _ := r.Value.(*rds.List)
		if rl == nil || len(rl.V) != len(f.Exprs) {
			out.Deparse = false
			if out.Detail == "" {
				out.Detail = "deparse: R returned a different number of expressions"
			}
			continue
		}
		for i, e := range f.Exprs {
			want := strs(rl.V[i])
			got := GoDeparse(e, mode)
			if strings.Join(got, "\n") != strings.Join(want, "\n") {
				out.Deparse = false
				if out.Detail == "" {
					out.Detail = firstLineDiff(mode, i+1, want, got)
				}
			}
		}
	}
	glines := syntax.Dump(f.Exprs)
	var treeDetail string
	out.Tree, out.ULP, treeDetail = compareTrees(rlines, glines)
	if out.Detail == "" {
		out.Detail = treeDetail
	}

	r, err := o.Call("tokens", src)
	if errors.Is(err, oracle.ErrR) {
		// R cannot build parse data for this text (it refuses invalid
		// UTF-8 when it splits the source into lines): nothing to compare
		out.Tokens = true
		return out, nil
	}
	if err != nil {
		return out, err
	}
	rt := strs(r.Value)
	gt := TokenLines(f.Tokens)
	out.Tokens = true
	if len(rt) != len(gt) {
		out.Tokens = false
		if out.Detail == "" {
			out.Detail = fmt.Sprintf("tokens: R has %d, Go %d", len(rt), len(gt))
			out.Detail += firstTokenDiff(rt, gt)
		}
	} else {
		for i := range rt {
			if rt[i] != gt[i] {
				out.Tokens = false
				if out.Detail == "" {
					out.Detail = "token\n  R:  " + decodeTok(rt[i]) + "\n  Go: " + decodeTok(gt[i])
				}
				break
			}
		}
	}
	return out, nil
}

func firstTokenDiff(rt, gt []string) string {
	for i := 0; i < len(rt) && i < len(gt); i++ {
		if rt[i] != gt[i] {
			return "; first difference\n  R:  " + decodeTok(rt[i]) + "\n  Go: " + decodeTok(gt[i])
		}
	}
	if len(rt) > len(gt) {
		return "; Go lacks " + decodeTok(rt[len(gt)])
	}
	return "; Go has extra " + decodeTok(gt[len(rt)])
}

func decodeTok(s string) string {
	i := strings.LastIndexByte(s, ' ')
	if i < 0 {
		return s
	}
	b, err := hex.DecodeString(s[i+1:])
	if err != nil {
		return s
	}
	return s[:i] + " " + strconv.Quote(string(b))
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

// compareTrees compares two dumps. Doubles that differ by at most 16 units in
// the last place are counted, not failed.
func compareTrees(r, g []string) (ok bool, ulp int, detail string) {
	n := len(r)
	if len(g) < n {
		n = len(g)
	}
	for i := 0; i < n; i++ {
		if r[i] == g[i] {
			continue
		}
		if d, same := NumClose(r[i], g[i]); same {
			ulp += d
			continue
		}
		return false, ulp, fmt.Sprintf("tree line %d\n  R:  %s\n  Go: %s", i+1, decodeLine(r[i]), decodeLine(g[i]))
	}
	if len(r) != len(g) {
		return false, ulp, fmt.Sprintf("tree: R has %d lines, Go %d", len(r), len(g))
	}
	return true, ulp, ""
}

func decodeLine(s string) string {
	// the node kind is the first word after the dots; hex arguments follow
	t := strings.TrimLeft(s, ".")
	prefix := s[:len(s)-len(t)]
	f := strings.Fields(t)
	if len(f) == 2 && (f[0] == "sym" || f[0] == "str" || f[0] == "name" || f[0] == "formal") && f[1] != "NA" {
		if b, err := hex.DecodeString(f[1]); err == nil {
			return prefix + f[0] + " " + strconv.Quote(string(b))
		}
	}
	return s
}

// NumClose reports whether two dump lines are doubles (or complex numbers)
// within 16 ulps of each other, and how many of their parts differ.
func NumClose(a, b string) (int, bool) {
	fa, fb := strings.Fields(a), strings.Fields(b)
	if len(fa) != len(fb) || len(fa) < 2 || strings.TrimLeft(fa[0], ".") != strings.TrimLeft(fb[0], ".") || len(fa[0]) != len(fb[0]) {
		return 0, false
	}
	kind := strings.TrimLeft(fa[0], ".")
	if kind != "dbl" && kind != "cplx" {
		return 0, false
	}
	n := 0
	for i := 1; i < len(fa); i++ {
		if fa[i] == fb[i] {
			continue
		}
		x, err1 := strconv.ParseUint(fa[i], 16, 64)
		y, err2 := strconv.ParseUint(fb[i], 16, 64)
		if err1 != nil || err2 != nil {
			return 0, false
		}
		fx, fy := math.Float64frombits(x), math.Float64frombits(y)
		if math.IsNaN(fx) || math.IsNaN(fy) || math.IsInf(fx, 0) != math.IsInf(fy, 0) {
			return 0, false
		}
		d := int64(x) - int64(y)
		if d < 0 {
			d = -d
		}
		if d > 16 {
			return 0, false
		}
		n++
	}
	return n, true
}

// ErrUnavailable is returned when the pinned R is not available.
var ErrUnavailable = errors.New("parity: the reference R is not available")

func firstLineDiff(mode string, expr int, want, got []string) string {
	for i := 0; i < len(want) || i < len(got); i++ {
		var w, g string
		if i < len(want) {
			w = want[i]
		}
		if i < len(got) {
			g = got[i]
		}
		if w != g || i >= len(want) || i >= len(got) {
			return fmt.Sprintf("deparse (%s) of expression %d, line %d\n  R:  %q\n  Go: %q", mode, expr, i+1, w, g)
		}
	}
	return "deparse"
}
