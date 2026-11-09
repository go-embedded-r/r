package evalcheck

import (
	"encoding/hex"

	"github.com/go-embedded-r/r/internal/gen"
	"github.com/go-embedded-r/r/rds"
)

// encIgnore, when set, makes the value comparison blind to the encoding mark
// of strings (the CHARSXP flag that says UTF-8 or native): R leaves what a
// connection reads unmarked, the evaluator does not track marks.
var encIgnore bool

// CheckIgnoringEncodingMarks is Check, except that the strings of the value
// are compared without their encoding marks (everything else exactly).
func CheckIgnoringEncodingMarks(rec gen.EvalRecord) string {
	encIgnore = true
	defer func() { encIgnore = false }()
	return check(rec, 0, nil)
}

// unmarked re-encodes a serialized value with every string mark cleared.
func unmarked(b []byte) []byte {
	v, _, err := rds.Decode(b, rds.Limits{})
	if err != nil {
		return b
	}
	clearMarks(v)
	out, err := rds.Encode(v)
	if err != nil {
		return b
	}
	return out
}

func clearMarks(v rds.Value) {
	switch x := v.(type) {
	case *rds.Character:
		for i := range x.V {
			x.V[i].Gp = 0
		}
		clearMeta(x.Meta)
	case *rds.List:
		for _, e := range x.V {
			clearMarks(e)
		}
		clearMeta(x.Meta)
	case *rds.Real:
		clearMeta(x.Meta)
	case *rds.Integer:
		clearMeta(x.Meta)
	case *rds.Logical:
		clearMeta(x.Meta)
	}
}

func clearMeta(m rds.Meta) {
	for _, a := range m.Attrs {
		clearMarks(a.Value)
	}
}

func hexUnmarked(h string) string {
	b, err := hex.DecodeString(h)
	if err != nil {
		return h
	}
	return hex.EncodeToString(unmarked(b))
}
