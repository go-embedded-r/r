package syntax_test

import (
	"errors"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/go-embedded-r/r/internal/gen"
	"github.com/go-embedded-r/r/syntax"
)

func seedCorpus(f *testing.F) {
	recs, err := gen.ReadJSONL[gen.ParseRecord]("../testdata/golden/parse.jsonl")
	if err != nil {
		f.Fatal(err)
	}
	for _, r := range recs {
		f.Add(r.Src)
	}
}

// FuzzParse: any input parses or fails with a typed error, never panics or
// hangs; what parses has tokens inside the source, a dump that is stable, and
// names that can be computed.
func FuzzParse(f *testing.F) {
	seedCorpus(f)
	lim := syntax.Limits{MaxBytes: 1 << 20, MaxDepth: 300, MaxNodes: 200000}
	f.Fuzz(func(t *testing.T, src string) {
		file, err := syntax.ParseAllLimits(src, lim)
		if err != nil {
			var se *syntax.Error
			if !errors.As(err, &se) && !errors.Is(err, syntax.ErrLimit) {
				t.Fatalf("untyped error %T: %v", err, err)
			}
			if se != nil && (se.Line < 0 || se.Col < 0) {
				t.Fatalf("bad position %+v", se)
			}
		}
		if file == nil {
			return
		}
		for _, tk := range file.Tokens {
			if tk.StartByte < 0 || tk.EndByte > len(file.Source) || tk.StartByte > tk.EndByte {
				t.Fatalf("token %+v outside the source", tk)
			}
		}
		if err == nil {
			a := strings.Join(syntax.Dump(file.Exprs), "\n")
			b := strings.Join(syntax.Dump(file.Exprs), "\n")
			if a != b {
				t.Fatal("dump not stable")
			}
			syntax.NamesOf(file.Exprs)
		}
	})
}

// FuzzLexNumbers targets the numeric literal rules.
func FuzzLexNumbers(f *testing.F) {
	for _, s := range []string{"1", "0x1p3", "1e5L", ".5", "1.5.5", "0x", "1e+", "9999999999999999999999e999", "0x.8p-1074"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, src string) {
		if len(src) > 200 {
			return
		}
		syntax.ParseAll(src)
	})
}

// FuzzDeparse: whatever parses can be deparsed, to at least one line of valid
// UTF-8 (deparse() is not a round trip: doubles keep 15 digits, 1i is written
// 0+1i, so only the shape of the output is checked here; its content is
// checked against R by cmd/rparity).
func FuzzDeparse(f *testing.F) {
	seedCorpus(f)
	lim := syntax.Limits{MaxBytes: 1 << 16, MaxDepth: 100, MaxNodes: 20000}
	f.Fuzz(func(t *testing.T, src string) {
		file, err := syntax.ParseAllLimits(src, lim)
		if err != nil {
			return
		}
		for _, e := range file.Exprs {
			lines := syntax.Deparse(e, nil)
			if len(lines) == 0 {
				t.Fatal("no line")
			}
			if !utf8.ValidString(strings.Join(lines, "\n")) && utf8.ValidString(src) {
				t.Fatalf("invalid UTF-8 from valid input: %q", lines)
			}
		}
	})
}
