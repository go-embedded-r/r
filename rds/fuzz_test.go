package rds_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/go-embedded-r/r/internal/gen"
	"github.com/go-embedded-r/r/rds"
)

// FuzzDecode feeds arbitrary bytes to the decoder, with tight limits so that
// a hostile input is refused rather than slow. What decodes must encode, and
// the bytes it encodes to must be a fixed point of decode-then-encode.
func FuzzDecode(f *testing.F) {
	recs, err := gen.ReadJSONL[gen.RDSRecord]("../testdata/golden/rds.jsonl")
	if err != nil {
		f.Fatal(err)
	}
	for _, r := range recs {
		b, _ := hex.DecodeString(r.Hex)
		f.Add(b)
	}
	f.Add([]byte("X\n"))
	f.Add([]byte("\x1f\x8b\x08\x00"))
	lim := rds.Limits{MaxBytes: 1 << 20, MaxDepth: 200}
	f.Fuzz(func(t *testing.T, in []byte) {
		v, h, err := rds.Decode(in, lim)
		if err != nil {
			return
		}
		out, err := rds.EncodeWith(v, h, 200)
		if err != nil {
			t.Fatalf("decoded but cannot encode: %v", err)
		}
		v2, h2, err := rds.Decode(out, lim)
		if err != nil {
			t.Fatalf("own output does not decode: %v", err)
		}
		out2, err := rds.EncodeWith(v2, h2, 200)
		if err != nil || !bytes.Equal(out, out2) {
			t.Fatalf("encode is not a fixed point: %v\n%x\n%x", err, out, out2)
		}
	})
}
