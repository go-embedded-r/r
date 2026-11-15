package eval_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-embedded-r/r/eval"
)

// Bombs of slice 18: the readers parse untrusted text and are bounded by
// MaxFileBytes (what a file can hold) and MaxVector (what they build). Each
// attack below is refused with an R error before anything large is allocated,
// under a memory limit, and leaves the session usable. The files are written in
// chunks, since no single vector may pass MaxVector either.
func TestSlice18Bombs(t *testing.T) {
	underMemoryLimit(t)
	small := func() *eval.Session {
		return eval.NewSession(eval.Options{MaxVector: 1 << 17, MaxAlloc: 1 << 28, MaxFileBytes: 4 << 20, Timeout: 30 * time.Second})
	}
	t.Run("a huge line", func(t *testing.T) {
		s := small()
		// one line of 240 000 fields: more than MaxVector items or columns
		if _, err := s.Eval(`con <- file("big.csv", "w"); for (i in 1:4) cat(strrep("a,", 60000), file = con); cat("\n", file = con); close(con)`); err != nil {
			t.Fatal(err)
		}
		for _, src := range []string{
			`scan("big.csv", what = "", sep = ",", quiet = TRUE)`,
			`read.csv("big.csv", header = FALSE)`,
			`read.table("big.csv", sep = ",")`,
			`read.csv("big.csv")`,
		} {
			refused(t, s, src, "cannot allocate vector", "too large", "limit", "memory")
		}
	})
	t.Run("a huge field count", func(t *testing.T) {
		s := small()
		if _, err := s.Eval(`con <- file("wide.txt", "w"); for (i in 1:4) cat(strrep("a ", 60000), file = con); cat("\n", file = con); close(con)`); err != nil {
			t.Fatal(err)
		}
		for _, src := range []string{
			`read.table("wide.txt", header = TRUE)`,
			`read.table("wide.txt", header = FALSE, colClasses = "character")`,
			`read.table("wide.txt", col.names = paste0("v", 1:130000))`,
		} {
			refused(t, s, src, "cannot allocate vector", "too large", "limit", "memory", "too many")
		}
	})
	t.Run("a quote that never closes", func(t *testing.T) {
		s := small()
		// 3 MB of text after one quote: count.fields refuses it, scan and read.table
		// stop at the end of the data with R's warning
		if _, err := s.Eval(`con <- file("open.txt", "w"); cat("a b\n\"", file = con); for (i in 1:30) cat(strrep("x", 1e5), file = con); cat("\n", file = con); close(con)`); err != nil {
			t.Fatal(err)
		}
		refused(t, s, `count.fields("open.txt")`, "quoted string on line 2 terminated by EOF")
		res, err := s.Eval(`x <- scan("open.txt", what = "", quiet = TRUE); length(x)`)
		if err != nil {
			t.Fatalf("scan of an open quote: %v", err)
		}
		if len(res.Warnings) == 0 || !strings.Contains(fmt.Sprint(res.Warnings), "EOF within quoted string") {
			t.Fatalf("no warning for an open quote: %v", res.Warnings)
		}
		// and a text that would pass the limit of a file cannot be written
		refused(t, s, `con <- file("huge.txt", "w"); for (i in 1:50) cat(strrep("x", 1e5), file = con)`, "limit", "cannot", "file too", "disk", "write", "space")
	})
	t.Run("2^31 rows", func(t *testing.T) {
		s := small()
		for _, src := range []string{
			`scan(text = rep("1", 2^31), quiet = TRUE)`,
			`read.table(text = rep("1 2", 2^31), nrows = 2^31)`,
			`read.csv(text = rep("1,2", 2^31))`,
			`type.convert(rep("1", 2^31), as.is = TRUE)`,
			`count.fields(textConnection(rep("1 2", 2^31)))`,
			`write.csv(data.frame(x = numeric(2^31)))`,
			`write.table(matrix(0, 2^20, 2^10))`,
			`scan(text = "1 2", what = numeric(2^31))`,
		} {
			refused(t, s, src, "cannot allocate vector", "too large", "limit", "memory", "invalid", "must", "negative", "too big")
		}
		// a bound on the number of rows is a bound, not an allocation
		if _, err := s.Eval(`stopifnot(nrow(read.table(text = "a\n1\n2", header = TRUE, nrows = 2^31)) == 2L, length(scan(text = "1 2 3", n = 2^31, quiet = TRUE)) == 3L)`); err != nil {
			t.Fatalf("nrows = 2^31 allocated or failed: %v", err)
		}
	})
	t.Run("more rows than a vector", func(t *testing.T) {
		s := small()
		if _, err := s.Eval(`con <- file("rows.txt", "w"); for (i in 1:3) writeLines(rep("1", 60000), con); close(con)`); err != nil {
			t.Fatal(err)
		}
		refused(t, s, `scan("rows.txt", quiet = TRUE)`, "cannot allocate vector")
		refused(t, s, `read.table("rows.txt")`, "cannot allocate vector")
		refused(t, s, `count.fields("rows.txt")`, "cannot allocate vector")
	})
	t.Run("a file too large to write", func(t *testing.T) {
		s := eval.NewSession(eval.Options{MaxVector: 1 << 24, MaxAlloc: 1 << 28, MaxFileBytes: 1 << 16, Timeout: 20 * time.Second})
		refused(t, s, `write.csv(data.frame(x = 1:1e5, y = "some text"), "f.csv")`, "limit", "cannot write", "file", "too")
	})
}

// A file system root, as a notebook host sets one, is read through the same
// readers and cannot be left.
func TestSlice18FileRoot(t *testing.T) {
	dir := t.TempDir()
	s := eval.NewSession(eval.Options{FileRoot: dir})
	defer s.Close()
	if _, err := s.Eval(`write.csv(data.frame(a = 1:3, b = c("x", "y", "z")), "t.csv", row.names = FALSE); d <- read.csv("t.csv"); stopifnot(identical(d, data.frame(a = 1:3, b = c("x", "y", "z"))))`); err != nil {
		t.Fatal(err)
	}
	for _, src := range []string{`read.csv("/etc/hosts")`, `scan("/etc/passwd", what = "")`, `read.table("../../etc/hosts")`, `count.fields("/etc/hosts")`} {
		if _, err := s.Eval(src); err == nil {
			t.Fatalf("%q read a file outside the root", src)
		}
	}
}
