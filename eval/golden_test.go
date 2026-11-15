package eval_test

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/go-embedded-r/r/internal/evalcheck"
	"github.com/go-embedded-r/r/internal/gen"
)

// ulpMark ends the snippets whose double values are compared within
// evalcheck.GoldenULP units in the last place (see README, slice 16).
const ulpMark = "# ulp"

// No golden shard may reach MaxEvalShard: the pre-push gate refuses blobs of
// 1 MiB, and the fix for a shard that grows is to split its snippet file.
func TestGoldenShardsStaySmall(t *testing.T) {
	entries, err := os.ReadDir("../" + gen.EvalShardDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) < 20 {
		t.Fatalf("only %d shards", len(entries))
	}
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			t.Fatal(err)
		}
		if info.Size() >= gen.MaxEvalShard {
			t.Errorf("%s is %d bytes, at or over %d: split testdata/eval/%s.R", e.Name(), info.Size(), gen.MaxEvalShard, strings.TrimSuffix(e.Name(), ".jsonl"))
		}
	}
}

// Every snippet of testdata/eval, evaluated by the Go evaluator, gives what R
// 4.6.2 gave when the golden was written: the value (serialized, byte for
// byte), its visibility, the warnings, or the error message and call.
func TestGoldens(t *testing.T) {
	recs, err := gen.ReadEvalGoldens("..")
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
		check := evalcheck.Check
		if strings.HasSuffix(rec.Src, ulpMark) {
			// a value whose last bits depend on whether the platform fuses a
			// multiply and an add: compared within GoldenULP units in the last
			// place, and its printed digits are not compared
			check = func(r gen.EvalRecord) string {
				var worst int64
				return evalcheck.CheckULPNormalized(r, evalcheck.GoldenULP, &worst, func(string) string { return "" })
			}
		}
		if why := check(rec); why != "" {
			bad++
			t.Errorf("%q: %s", rec.Src, why)
		}
	}
	t.Logf("%d of %d snippets agree with R (%d recorded for another platform, not compared here)", len(recs)-bad-skipped, len(recs)-skipped, skipped)
}
