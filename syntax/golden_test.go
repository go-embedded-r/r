package syntax_test

import (
	"strings"
	"testing"

	"github.com/go-embedded-r/r/internal/gen"
	"github.com/go-embedded-r/r/internal/parity"
	"github.com/go-embedded-r/r/syntax"
)

// Every golden record: R's tree, tokens and error message, which the Go
// parser must reproduce exactly. No R is needed.
func TestGoldenParity(t *testing.T) {
	recs, err := gen.ReadJSONL[gen.ParseRecord]("../testdata/golden/parse.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	if len(recs) < 250 {
		t.Fatalf("only %d records", len(recs))
	}
	for _, rec := range recs {
		f, err := syntax.ParseAll(rec.Src)
		if rec.Error != "" {
			if err == nil {
				t.Errorf("%q: R rejects (%s), Go accepts", rec.Src, rec.Error)
			} else if err.Error() != rec.Error {
				t.Errorf("%q:\n R:  %s\n Go: %s", rec.Src, rec.Error, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("%q: R parses, Go says %v", rec.Src, err)
			continue
		}
		if got := syntax.Dump(f.Exprs); strings.Join(got, "\n") != strings.Join(rec.Tree, "\n") {
			t.Errorf("%q: tree differs\n R:  %v\n Go: %v", rec.Src, rec.Tree, got)
		}
		for mode, want := range map[string][][]string{"default": rec.Deparse, "print": rec.Print} {
			if len(want) != len(f.Exprs) {
				t.Errorf("%q: %d deparsed expressions in the golden, %d parsed", rec.Src, len(want), len(f.Exprs))
				continue
			}
			for i, e := range f.Exprs {
				if got := parity.GoDeparse(e, mode); strings.Join(got, "\n") != strings.Join(want[i], "\n") {
					t.Errorf("%q: deparse (%s) of expression %d\n R:  %q\n Go: %q", rec.Src, mode, i+1, want[i], got)
				}
			}
		}
		if got := parity.TokenLines(f.Tokens); strings.Join(got, "\n") != strings.Join(rec.Tokens, "\n") {
			t.Errorf("%q: tokens differ\n R:  %v\n Go: %v", rec.Src, rec.Tokens, got)
		}
	}
}
