package eval_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-embedded-r/r/eval"
)

// FuzzReadTable: the readers (scan, read.table and its relatives, count.fields,
// type.convert) and the writers, given text of any shape and any combination of
// the options, end with a value or an R error, never a panic or a hang, and
// leave the session usable.
func FuzzReadTable(f *testing.F) {
	f.Add([]byte("a,b\n1,x\n2,y\n"), uint8(1), uint8(0), uint8(0))
	f.Add([]byte("a b\nr1 1 2\nr2 3 4"), uint8(0), uint8(1), uint8(1))
	f.Add([]byte("x\n\"p\nq\"\n'a''b'\n"), uint8(2), uint8(2), uint8(3))
	f.Add([]byte("1 2 3\n4 5\n6\n"), uint8(0), uint8(3), uint8(4))
	f.Add([]byte("\"never closed\n1,2\n"), uint8(1), uint8(4), uint8(5))
	f.Add([]byte("a\tb\n\\t\\n\\x41\\101\n"), uint8(3), uint8(5), uint8(6))
	f.Add([]byte("# c\n1;2,5\n3;4,5 # d\n"), uint8(4), uint8(6), uint8(7))
	f.Add([]byte(",,,\n,\n\n \n"), uint8(1), uint8(7), uint8(8))
	f.Add([]byte("1e5 0x1F Inf -inf NaN NA TRUE T 1+2i 2147483648\n"), uint8(0), uint8(0), uint8(9))
	f.Fuzz(func(t *testing.T, data []byte, mode, flags, opt uint8) {
		if len(data) > 512 {
			data = data[:512]
		}
		txt := strings.ToValidUTF8(strings.ReplaceAll(string(data), "\x00", " "), " ")
		lit := fmt.Sprintf("%q", txt)
		seps := []string{"", ",", ";", "\t", "|", " ", "\n"}
		sep := seps[int(mode)%len(seps)]
		quotes := []string{"\"'", "\"", "'", ""}
		quote := quotes[int(flags)%len(quotes)]
		decs := []string{".", ","}
		dec := decs[int(opt)%2]
		b := func(bit uint, v *uint8) string {
			if *v&(1<<bit) != 0 {
				return "TRUE"
			}
			return "FALSE"
		}
		comment := []string{"#", "", "%"}[int(opt>>1)%3]
		common := fmt.Sprintf("sep = %q, quote = %q, dec = %q, fill = %s, strip.white = %s, comment.char = %q, blank.lines.skip = %s",
			sep, quote, dec, b(3, &flags), b(4, &flags), comment, b(5, &flags))
		s := eval.NewSession(eval.Options{MaxVector: 1 << 16, MaxAlloc: 1 << 26, MaxFileBytes: 1 << 20, Timeout: 2 * time.Second})
		progs := []string{
			fmt.Sprintf("scan(text = %s, what = %s, %s, quiet = TRUE, nlines = %d)", lit, []string{"0", `""`, "0L", "TRUE", "list(0, \"\")", "list(a = \"\", b = 0L, c = NULL)", "complex()"}[int(mode>>3)%7], common, int(opt>>4)),
			fmt.Sprintf("read.table(text = %s, header = %s, %s, nrows = %d, skip = %d, flush = %s, allowEscapes = %s)", lit, []string{"TRUE", "FALSE"}[int(flags>>5)%2], common, int(opt>>2)%5-1, int(opt>>6), b(6, &flags), b(7, &flags)),
			fmt.Sprintf("read.csv(text = %s, row.names = %d)", lit, int(opt>>3)%3),
			fmt.Sprintf("read.delim(text = %s, colClasses = %q, na.strings = %q)", lit, []string{"character", "numeric", "NULL", "factor", "integer"}[int(opt)%5], []string{"NA", "", "x"}[int(flags)%3]),
			fmt.Sprintf("tf <- tempfile(); writeLines(%s, tf); count.fields(tf, sep = %q, quote = %q, comment.char = %q)", lit, sep, quote, comment),
			fmt.Sprintf("type.convert(strsplit(%s, \"[ ,\\n]\")[[1]], as.is = %s, dec = %q, numerals = %q)", lit, b(0, &opt), dec, []string{"allow.loss", "warn.loss", "no.loss"}[int(flags)%3]),
			fmt.Sprintf("d <- tryCatch(read.csv(text = %s), error = function(e) data.frame(x = 1)); write.csv(d); write.table(d, quote = %s, sep = %q, na = \"-\", dec = %q, row.names = %s)", lit, b(1, &opt), sep+";", dec, b(2, &opt)),
		}
		for _, p := range progs {
			s.Eval(p)
		}
		if _, err := s.Eval("NULL"); err != nil {
			t.Fatalf("session broken after %q: %v", progs, err)
		}
	})
}
