package eval_test

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/go-embedded-r/r/eval"
	"github.com/go-embedded-r/r/internal/gen"
	"github.com/go-embedded-r/r/rds"
)

func warnTexts(w []eval.Warning) []string {
	var out []string
	for _, x := range w {
		out = append(out, x.Message)
	}
	return out
}

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
	bad := 0
	for _, rec := range recs {
		s := eval.NewSession(eval.Options{})
		res, err := s.EvalIn(rec.Src, eval.NewEnv(s.Global))
		var why string
		switch {
		case rec.OK && err != nil:
			why = "R succeeds, Go fails: " + err.Error()
		case !rec.OK && err == nil:
			why = "R fails (" + rec.Error + "), Go succeeds"
		case !rec.OK:
			e, ok := err.(*eval.Error)
			switch {
			case !ok:
				why = "Go: " + err.Error()
			case e.Message != rec.Error || e.Call != rec.Call:
				why = "error\n  R:  " + rec.Error + " [" + rec.Call + "]\n  Go: " + e.Message + " [" + e.Call + "]"
			}
			if why == "" {
				why = cmpWarnings(rec.Warnings, nil)
			}
		default:
			if rec.Typeof != "closure" && rec.Typeof != "builtin" && rec.Typeof != "special" {
				b, eerr := rds.Encode(eval.ToRDS(res.Value))
				if eerr != nil {
					why = eerr.Error()
				} else if hex.EncodeToString(b) != rec.Hex {
					want, _, _ := rds.Decode(mustHex(rec.Hex), rds.Limits{})
					got, _, _ := rds.Decode(b, rds.Limits{})
					why = "value\n  R:  " + describe(want) + "\n  Go: " + describe(got)
				}
			}
			if why == "" && res.Visible != rec.Visible {
				why = "visibility differs"
			}
			if why == "" {
				why = cmpWarnings(rec.Warnings, warnTexts(res.Warnings))
			}
		}
		if why != "" {
			bad++
			t.Errorf("%q: %s", rec.Src, why)
		}
	}
	t.Logf("%d of %d snippets agree with R", len(recs)-bad, len(recs))
}

func cmpWarnings(want, got []string) string {
	if strings.Join(want, "\x00") != strings.Join(got, "\x00") {
		return "warnings\n  R:  " + strings.Join(want, " | ") + "\n  Go: " + strings.Join(got, " | ")
	}
	return ""
}

func mustHex(s string) []byte { b, _ := hex.DecodeString(s); return b }

func describe(v rds.Value) string {
	if v == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%+v", v)
}
