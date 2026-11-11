package eval_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-embedded-r/r/eval"
)

// FuzzReadBin: the decoder of readBin and readChar, given bytes of any shape
// and any size, byte order and count, ends with a value or an R error, never a
// panic, and leaves the session usable.
func FuzzReadBin(f *testing.F) {
	f.Add([]byte{1, 2, 3, 4, 5, 6, 7, 8}, uint8(0), uint8(4), uint8(0), int16(2))
	f.Add([]byte{0xff, 0xfe}, uint8(0), uint8(1), uint8(1), int16(5))
	f.Add([]byte("ab\x00cd\x00e"), uint8(5), uint8(0), uint8(0), int16(3))
	f.Add([]byte{0, 0, 0xf0, 0x3f, 0, 0, 0, 0}, uint8(1), uint8(8), uint8(2), int16(1))
	f.Add([]byte{}, uint8(4), uint8(0), uint8(0), int16(-1))
	f.Fuzz(func(t *testing.T, data []byte, what, size, flags uint8, n int16) {
		if len(data) > 256 {
			data = data[:256]
		}
		parts := make([]string, len(data))
		for i, b := range data {
			parts[i] = fmt.Sprint(b)
		}
		raw := "as.raw(c(" + strings.Join(parts, ",") + "))"
		if len(data) == 0 {
			raw = "raw(0)"
		}
		whats := []string{"integer", "double", "logical", "complex", "raw", "character", "numeric", "foo"}
		sz := ""
		if size%9 != 0 {
			sz = fmt.Sprintf(", size = %d", size%17)
		}
		endian := []string{"little", "big", "swap"}[flags%3]
		signed := "TRUE"
		if flags&4 != 0 {
			signed = "FALSE"
		}
		s := eval.NewSession(eval.Options{MaxVector: 1 << 16, MaxAlloc: 1 << 26, MaxFileBytes: 1 << 20, Timeout: 2 * time.Second})
		src := fmt.Sprintf("readBin(%s, %q, n = %d%s, signed = %s, endian = %q)", raw, whats[int(what)%len(whats)], n, sz, signed, endian)
		s.Eval(src)
		s.Eval(fmt.Sprintf("readChar(%s, %d, useBytes = %v)", raw, n, flags&8 != 0))
		s.Eval(fmt.Sprintf("con <- rawConnection(%s); readBin(con, %q, %d%s); readBin(con, \"raw\", 3); close(con)", raw, whats[int(what)%len(whats)], n, sz))
		if _, err := s.Eval("NULL"); err != nil {
			t.Fatalf("session broken after %q: %v", src, err)
		}
	})
}
