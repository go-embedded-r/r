package eval_test

import (
	"runtime"
	"testing"

	"github.com/go-embedded-r/r/internal/evalcheck"
	"github.com/go-embedded-r/r/internal/gen"
)

// Every snippet of testdata/eval, evaluated by the Go evaluator, gives what R
// 4.6.2 gave when the golden was written: the value (serialized, byte for
// byte), its visibility, the warnings, or the error message and call.
func TestGoldens(t *testing.T) {
	recs, err := gen.ReadJSONL[gen.EvalRecord]("../testdata/golden/eval.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	if len(recs) < 250 {
		t.Fatalf("only %d records", len(recs))
	}
	bad, skipped := 0, 0
	for _, rec := range recs {
		if rec.Platform != "" && rec.Platform != runtime.GOOS {
			skipped++ // recorded for another platform, where R's C library answers differently
			continue
		}
		if why := evalcheck.Check(rec); why != "" {
			bad++
			t.Errorf("%q: %s", rec.Src, why)
		}
	}
	t.Logf("%d of %d snippets agree with R (%d recorded for another platform, not compared here)", len(recs)-bad-skipped, len(recs)-skipped, skipped)
}
