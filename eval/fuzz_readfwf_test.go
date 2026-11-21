package eval_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-embedded-r/r/eval"
)

// FuzzReadFwf: read.fwf, given text of any shape, widths of either sign and
// any size, and any combination of header, skip, n, sep and col.names, ends
// with a data frame or an R error, never a panic or a hang, and leaves the
// session usable. The corpus is committed in testdata/fuzz/FuzzReadFwf.
func FuzzReadFwf(f *testing.F) {
	f.Add([]byte("12345\n67890\n"), uint8(2), uint8(3), uint8(0), uint8(0))
	f.Add([]byte("h1h2\nabcd\nefgh\n"), uint8(2), uint8(2), uint8(0), uint8(1))
	f.Add([]byte("ab\ncdefg\n\nx"), uint8(1), uint8(250), uint8(7), uint8(2))
	f.Add([]byte("1 2 3\n4 5 6\n"), uint8(1), uint8(129), uint8(1), uint8(4))
	f.Add([]byte("\"q\"x\n"), uint8(0), uint8(255), uint8(255), uint8(8))
	f.Add([]byte(""), uint8(3), uint8(3), uint8(3), uint8(16))
	f.Add([]byte("日本語テキスト\nあいうえお\n"), uint8(2), uint8(2), uint8(2), uint8(32))
	f.Fuzz(func(t *testing.T, data []byte, w1, w2, w3, flags uint8) {
		if len(data) > 512 {
			data = data[:512]
		}
		txt := strings.ToValidUTF8(strings.ReplaceAll(string(data), "\x00", " "), " ")
		width := func(b uint8) string { return fmt.Sprint(int(int8(b))) } // -128 to 127: negative skips a field
		widths := fmt.Sprintf("c(%s, %s, %s)", width(w1), width(w2), width(w3))
		if flags&64 != 0 {
			widths = fmt.Sprintf("rep(%s, %d)", width(w1), int(w2)) // many fields
		}
		if flags&128 != 0 {
			widths = fmt.Sprintf("c(%s, 2^40, %s)", width(w1), width(w3)) // a width beyond any line
		}
		seps := []string{"\t", ",", " ", ""}
		opts := fmt.Sprintf("header = %s, skip = %d, n = %d, sep = %q",
			[]string{"FALSE", "TRUE"}[flags&1], int(flags>>1)&3, int(flags>>3)&3-1, seps[int(flags>>2)&3])
		if flags&32 != 0 {
			opts += `, col.names = c("a", "b", "c")`
		}
		s := eval.NewSession(eval.Options{MaxVector: 1 << 16, MaxAlloc: 1 << 26, MaxFileBytes: 1 << 20, Timeout: 2 * time.Second, MaxDepth: 200})
		prog := fmt.Sprintf("read.fwf(textConnection(%q), widths = %s, %s)", txt, widths, opts)
		s.Eval(prog)
		if _, err := s.Eval("NULL"); err != nil {
			t.Fatalf("session broken after %q: %v", prog, err)
		}
	})
}
