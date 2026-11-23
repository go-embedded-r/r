package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// budget is testdata/sizebudget.json: the size of this program under each wasm
// target when the budget was last set, and how far above it a build may go.
// Growing past it is a decision (go-embedded-r/r#62): raise the number in the
// same pull request and say what the bytes buy.
type budget struct {
	TolerancePercent float64          `json:"tolerance_percent"`
	Bytes            map[string]int64 `json:"bytes"`
}

func TestWasmSizeBudget(t *testing.T) {
	if testing.Short() {
		t.Skip("builds two wasm binaries")
	}
	if runtime.GOOS == "js" || runtime.GOOS == "wasip1" {
		t.Skip("needs the go tool")
	}
	raw, err := os.ReadFile(filepath.Join("..", "..", "testdata", "sizebudget.json"))
	if err != nil {
		t.Fatal(err)
	}
	var b budget
	if err := json.Unmarshal(raw, &b); err != nil {
		t.Fatal(err)
	}
	if b.TolerancePercent <= 0 || len(b.Bytes) == 0 {
		t.Fatalf("sizebudget.json is incomplete: %+v", b)
	}
	for goos, base := range b.Bytes {
		t.Run(goos, func(t *testing.T) {
			out := filepath.Join(t.TempDir(), "rsize.wasm")
			cmd := exec.Command("go", "build", "-o", out, ".")
			cmd.Env = append(os.Environ(), "GOOS="+goos, "GOARCH=wasm")
			if msg, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("go build: %v\n%s", err, msg)
			}
			fi, err := os.Stat(out)
			if err != nil {
				t.Fatal(err)
			}
			limit := int64(float64(base) * (1 + b.TolerancePercent/100))
			t.Logf("%s: %d bytes (budget %d, limit %d)", goos, fi.Size(), base, limit)
			if fi.Size() > limit {
				t.Errorf("%s wasm is %d bytes, over the budget %d by more than %.0f%%", goos, fi.Size(), base, b.TolerancePercent)
			}
			if fi.Size() < base*9/10 {
				t.Errorf("%s wasm is %d bytes, more than 10%% under the budget %d: lower the budget in testdata/sizebudget.json", goos, fi.Size(), base)
			}
		})
	}
}
