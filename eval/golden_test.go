package eval_test

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/go-embedded-r/r/internal/evalcheck"
	"github.com/go-embedded-r/r/internal/gen"
)

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
		if why := evalcheck.Check(rec); why != "" {
			bad++
			t.Errorf("%q: %s", rec.Src, why)
		}
	}
	t.Logf("%d of %d snippets agree with R (%d recorded for another platform, not compared here)", len(recs)-bad-skipped, len(recs)-skipped, skipped)
}
