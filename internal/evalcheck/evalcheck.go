// Package evalcheck runs a snippet through the Go evaluator and compares the
// outcome with what R 4.6.2 did (an EvalRecord): value bytes, visibility,
// warnings, error message and call.
package evalcheck

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"strings"

	"github.com/go-embedded-r/r/eval"
	"github.com/go-embedded-r/r/internal/gen"
	"github.com/go-embedded-r/r/rds"
)

// Check returns "" when the Go evaluator agrees with the record, else a
// description of the first difference.
func Check(rec gen.EvalRecord) string { return check(rec, 0, nil) }

// CheckNormalized is Check, except that the printed output of both sides goes
// through norm before it is compared (for text R makes unpredictable, such as
// the address of an environment). Value, visibility, errors and warnings are
// still compared exactly.
func CheckNormalized(rec gen.EvalRecord, norm func(string) string) string {
	outNorm = norm
	defer func() { outNorm = nil }()
	return check(rec, 0, nil)
}

// outNorm, when set, is applied to the outputs compared by check.
var outNorm func(string) string

func normOut(s string) string {
	if outNorm != nil {
		return outNorm(s)
	}
	return s
}

// CheckULP is Check, except that double values may differ by up to maxULP
// units in the last place; the largest difference seen is recorded in worst.
// It is for results that depend on the C library (libm), nothing else.
func CheckULP(rec gen.EvalRecord, maxULP int64, worst *int64) string {
	return check(rec, maxULP, worst)
}

func check(rec gen.EvalRecord, maxULP int64, worst *int64) string {
	s := eval.NewSession(eval.Options{})
	res, err := s.EvalIn(context.Background(), rec.Src, eval.NewEnv(s.Global))
	if e, ok := err.(*eval.Error); ok && !rec.OK {
		if normOut(e.Output) != normOut(rec.Output) {
			return fmt.Sprintf("output before the error\n  R:  %q\n  Go: %q", rec.Output, e.Output)
		}
		if e.Stderr != rec.Stderr {
			return fmt.Sprintf("stderr before the error\n  R:  %q\n  Go: %q", rec.Stderr, e.Stderr)
		}
	}
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
		b, eerr := rds.Encode(gen.NormalizeNaN(eval.ToRDS(res.Value)))
		if eerr != nil {
			return eerr.Error()
		}
		got, want := hex.EncodeToString(b), rec.Hex
		if encIgnore {
			got, want = hexUnmarked(got), hexUnmarked(want)
		}
		if got != want {
			want, _, _ := rds.Decode(mustHex(rec.Hex), rds.Limits{})
			got, _, _ := rds.Decode(b, rds.Limits{})
			if maxULP > 0 && ulpEqual(want, got, maxULP, worst) {
				goto rest
			}
			return fmt.Sprintf("value\n  R:  %s\n  Go: %s\n  R bytes  %s\n  Go bytes %s", dump(want), dump(got), rec.Hex, hex.EncodeToString(b))
		}
	}
rest:
	if res.Visible != rec.Visible {
		return "visibility differs"
	}
	if normOut(res.Output) != normOut(rec.Output) {
		return fmt.Sprintf("output\n  R:  %q\n  Go: %q", rec.Output, res.Output)
	}
	if res.Stderr != rec.Stderr {
		return fmt.Sprintf("stderr\n  R:  %q\n  Go: %q", rec.Stderr, res.Stderr)
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

// ulpDistance is the number of representable doubles between a and b.
func ulpDistance(a, b float64) int64 {
	if a == b {
		return 0
	}
	ia, ib := int64(math.Float64bits(a)), int64(math.Float64bits(b))
	if ia < 0 {
		ia = math.MinInt64 - ia
	}
	if ib < 0 {
		ib = math.MinInt64 - ib
	}
	d := ia - ib
	if d < 0 {
		d = -d
	}
	return d
}

func ulpEqual(a, b rds.Value, maxULP int64, worst *int64) bool {
	switch x := a.(type) {
	case *rds.Real:
		y, ok := b.(*rds.Real)
		if !ok || len(x.V) != len(y.V) || !ulpMeta(x.Meta, y.Meta, maxULP, worst) {
			return false
		}
		for i := range x.V {
			p, q := x.V[i], y.V[i]
			if math.IsNaN(p) || math.IsNaN(q) {
				if math.IsNaN(p) != math.IsNaN(q) {
					return false
				}
				continue
			}
			if math.IsInf(p, 0) || math.IsInf(q, 0) {
				if p != q {
					return false
				}
				continue
			}
			d := ulpDistance(p, q)
			if d > maxULP {
				return false
			}
			if worst != nil && d > *worst {
				*worst = d
			}
		}
		return true
	case *rds.List:
		y, ok := b.(*rds.List)
		if !ok || x.Expr != y.Expr || len(x.V) != len(y.V) || !ulpMeta(x.Meta, y.Meta, maxULP, worst) {
			return false
		}
		for i := range x.V {
			if !ulpEqual(x.V[i], y.V[i], maxULP, worst) {
				return false
			}
		}
		return true
	}
	ea, e1 := rds.Encode(a)
	eb, e2 := rds.Encode(b)
	return e1 == nil && e2 == nil && string(ea) == string(eb)
}

func ulpMeta(a, b rds.Meta, maxULP int64, worst *int64) bool {
	if a.Obj != b.Obj || a.Gp != b.Gp || len(a.Attrs) != len(b.Attrs) {
		return false
	}
	for i := range a.Attrs {
		if a.Attrs[i].Name != b.Attrs[i].Name || !ulpEqual(a.Attrs[i].Value, b.Attrs[i].Value, maxULP, worst) {
			return false
		}
	}
	return true
}
