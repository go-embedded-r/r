// Package evalcheck runs a snippet through the Go evaluator and compares the
// outcome with what R 4.6.2 did (an EvalRecord): value bytes, visibility,
// warnings, error message and call.
package evalcheck

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/go-embedded-r/r/eval"
	"github.com/go-embedded-r/r/internal/gen"
	"github.com/go-embedded-r/r/rds"
)

// Check returns "" when the Go evaluator agrees with the record, else a
// description of the first difference.
func Check(rec gen.EvalRecord) string {
	s := eval.NewSession(eval.Options{})
	res, err := s.EvalIn(rec.Src, eval.NewEnv(s.Global))
	switch {
	case rec.OK && err != nil:
		return "R succeeds, Go fails: " + err.Error()
	case !rec.OK && err == nil:
		return "R fails (" + rec.Error + "), Go succeeds"
	case !rec.OK:
		e, ok := err.(*eval.Error)
		switch {
		case !ok:
			return "Go: " + err.Error()
		case e.Message != rec.Error || e.Call != rec.Call:
			return "error\n  R:  " + rec.Error + " [" + rec.Call + "]\n  Go: " + e.Message + " [" + e.Call + "]"
		}
		return ""
	}
	if rec.Typeof != "closure" && rec.Typeof != "builtin" && rec.Typeof != "special" {
		b, eerr := rds.Encode(eval.ToRDS(res.Value))
		if eerr != nil {
			return eerr.Error()
		}
		if hex.EncodeToString(b) != rec.Hex {
			want, _, _ := rds.Decode(mustHex(rec.Hex), rds.Limits{})
			got, _, _ := rds.Decode(b, rds.Limits{})
			return fmt.Sprintf("value\n  R:  %+v\n  Go: %+v\n  R bytes  %s\n  Go bytes %s", want, got, rec.Hex, hex.EncodeToString(b))
		}
	}
	if res.Visible != rec.Visible {
		return "visibility differs"
	}
	var got []string
	for _, w := range res.Warnings {
		got = append(got, w.Message)
	}
	if strings.Join(rec.Warnings, "\x00") != strings.Join(got, "\x00") {
		return "warnings\n  R:  " + strings.Join(rec.Warnings, " | ") + "\n  Go: " + strings.Join(got, " | ")
	}
	return ""
}

func mustHex(s string) []byte { b, _ := hex.DecodeString(s); return b }
