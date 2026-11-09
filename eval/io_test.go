package eval_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-embedded-r/r/eval"
	"github.com/go-embedded-r/r/syntax"
)

// The files a session sees are private by default and confined when they are
// the host's: these tests check what the goldens cannot, the boundary.

func evalOut(t *testing.T, s *eval.Session, src string) string {
	t.Helper()
	res, err := s.Eval(src)
	if err != nil {
		t.Fatalf("%q: %v", src, err)
	}
	return res.Output
}

func TestDefaultFileSystemIsPrivate(t *testing.T) {
	skipOnWasm(t)
	host := filepath.Join(t.TempDir(), "private", "f.txt")
	s := eval.NewSession(eval.Options{})
	src := `dir.create("` + filepath.ToSlash(filepath.Dir(host)) + `", recursive = TRUE); writeLines("secret", "` + filepath.ToSlash(host) + `"); print(readLines("` + filepath.ToSlash(host) + `")); print(file.exists("` + filepath.ToSlash(host) + `"))`
	if out := evalOut(t, s, src); !strings.Contains(out, `"secret"`) || !strings.Contains(out, "TRUE") {
		t.Fatalf("the session cannot use its own file system: %q", out)
	}
	if _, err := os.Stat(host); err == nil {
		t.Fatalf("a default session wrote %s on the host", host)
	}
	// another session does not see it
	if out := evalOut(t, eval.NewSession(eval.Options{}), `print(file.exists("`+filepath.ToSlash(host)+`"))`); !strings.Contains(out, "FALSE") {
		t.Fatalf("sessions share files: %q", out)
	}
	// no host file is readable either
	if _, err := eval.NewSession(eval.Options{}).Eval(`readLines("/etc/hosts")`); err == nil {
		t.Fatal("a default session read /etc/hosts")
	}
}

func TestFileRootIsConfined(t *testing.T) {
	skipOnWasm(t)
	dir := t.TempDir()
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "secret.txt"), []byte("outside\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	symlinked := os.Symlink(outside, filepath.Join(dir, "escape")) == nil
	s := eval.NewSession(eval.Options{FileRoot: dir})
	defer s.Close()
	evalOut(t, s, `writeLines(c("a", "b"), "/inside.txt"); dir.create("/sub"); cat("x", file = "/sub/f.txt")`)
	if b, err := os.ReadFile(filepath.Join(dir, "inside.txt")); err != nil || string(b) != "a\nb\n" {
		t.Fatalf("the file is not under FileRoot: %q %v", b, err)
	}
	if b, err := os.ReadFile(filepath.Join(dir, "sub", "f.txt")); err != nil || string(b) != "x" {
		t.Fatalf("sub/f.txt: %q %v", b, err)
	}
	// ".." is clamped at the root: this names <root>/etc/passwd
	for _, p := range []string{"/../../../../etc/passwd", "../../../../etc/passwd", "/sub/../../../etc/hosts", "/./../etc/hosts"} {
		if _, err := s.Eval(`readLines("` + p + `")`); err == nil {
			t.Errorf("%s was read", p)
		}
	}
	if symlinked {
		for _, src := range []string{`readLines("/escape/secret.txt")`, `writeLines("pwn", "/escape/new.txt")`, `file.exists("/escape/secret.txt")`} {
			res, err := s.Eval(src)
			if src[0] == 'f' {
				if err != nil || strings.Contains(res.Output, "TRUE") || (res.Value != nil && res.Visible && strings.Contains(strings.ToUpper(res.Output), "TRUE")) {
					t.Errorf("%s: %v %v", src, res, err)
				}
				continue
			}
			if err == nil {
				t.Errorf("%s was not refused", src)
			}
		}
		if _, err := os.Stat(filepath.Join(outside, "new.txt")); err == nil {
			t.Fatal("a file was written outside FileRoot through a symbolic link")
		}
	}
	// the host's own tree is not reachable by another name either
	if _, err := s.Eval(`unlink("/", recursive = TRUE)`); err != nil {
		t.Fatalf("unlink of the root: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("the root itself was removed: %v", err)
	}
}

func TestEnvironmentIsNotTheHosts(t *testing.T) {
	t.Setenv("RGO_HOST_SECRET", "hunter2")
	s := eval.NewSession(eval.Options{Env: map[string]string{"GIVEN": "yes"}})
	if out := evalOut(t, s, `print(Sys.getenv("RGO_HOST_SECRET")); print(Sys.getenv("GIVEN")); print(Sys.getenv("HOME"))`); strings.Contains(out, "hunter2") || !strings.Contains(out, `"yes"`) {
		t.Fatalf("session environment: %q", out)
	}
}

func TestFileBombs(t *testing.T) {
	underMemoryLimit(t)
	s := eval.NewSession(eval.Options{MaxFileBytes: 1 << 20, MaxSteps: eval.Unbounded, Timeout: 3 * time.Second})
	refused(t, s, `f <- tempfile(); x <- strrep("x", 1e5); repeat cat(x, file = f, append = TRUE)`, "file size limit", "cannot write", "File too large")
	refused(t, s, `f <- tempfile(); x <- strrep("x", 2e6); writeLines(x, f)`, "file size limit", "cannot write", "File too large")
	refused(t, s, `tc <- textConnection("out", "w"); x <- strrep("x", 1e5); repeat writeLines(x, tc)`, "file size limit")
	refused(t, s, `x <- strrep("x", 1e5); sink(tempfile()); repeat print(x)`, "file size limit", "memory limit", "time limit")
	refused(t, s, `for (i in 1:200) textConnection("a")`, "all connections are in use")
	refused(t, s, `for (i in 1:30) sink(textConnection(NULL, "w"))`, "sink stack is full", "all connections are in use")
	s2 := eval.NewSession(eval.Options{MaxFileBytes: 1 << 20})
	// the size limit binds what is read as well
	if _, err := s2.Eval(`f <- tempfile(); writeLines(strrep("y", 9e5), f); file.copy(f, g <- tempfile()); file.copy(f, h <- tempfile())`); err != nil {
		t.Fatalf("under the limit: %v", err)
	}
	refused(t, s2, `f <- tempfile(); writeLines(strrep("y", 9e5), f); for (i in 1:5) file.copy(f, tempfile())`, "cannot", "limit")
}

// A compressed document that expands past the limit is refused, not expanded.
func TestReadRDSBomb(t *testing.T) {
	underMemoryLimit(t)
	dir := t.TempDir()
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	zeros := make([]byte, 1<<20)
	for i := 0; i < 64; i++ {
		zw.Write(zeros)
	}
	zw.Close()
	if err := os.WriteFile(filepath.Join(dir, "bomb.rds"), buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	s := eval.NewSession(eval.Options{FileRoot: dir, MaxFileBytes: 1 << 20})
	defer s.Close()
	refused(t, s, `readRDS("/bomb.rds")`, "unknown input format")
}

// Many small writes to a text connection and a file stay linear.
func TestWritesAreLinear(t *testing.T) {
	skipOnWasm(t)
	s := eval.NewSession(eval.Options{MaxSteps: eval.Unbounded})
	src := `tc <- textConnection("out", "w"); for (i in 1:20000) writeLines("line", tc); close(tc); length(out)`
	res, err := s.Eval(src)
	if err != nil || !strings.Contains(res.Output, "20000") {
		t.Fatalf("%v %v", res, err)
	}
	src = `f <- tempfile(); con <- file(f, "w"); for (i in 1:20000) writeLines("line", con); close(con); length(readLines(f))`
	if res, err = s.Eval(src); err != nil || !strings.Contains(res.Output, "20000") {
		t.Fatalf("%v %v", res, err)
	}
}

// A sink left open by a cell keeps diverting in the next one, and what it
// diverted is not in any Result.
func TestSinkSpansCells(t *testing.T) {
	s := eval.NewSession(eval.Options{})
	res, err := s.Eval(`tc <- textConnection("kept", "w"); sink(tc); print(1)`)
	if err != nil || res.Output != "" {
		t.Fatalf("output leaked: %q %v", res.Output, err)
	}
	res, err = s.Eval(`print(2); sink(); close(tc); kept`)
	if err != nil || strings.Contains(res.Output, "[1] 2") && !strings.Contains(res.Output, `"[1] 1" "[1] 2"`) {
		t.Fatalf("second cell: %q %v", res.Output, err)
	}
	if !strings.Contains(res.Output, `"[1] 1" "[1] 2"`) {
		t.Fatalf("kept: %q", res.Output)
	}
}

// Windows and other exotic path forms are only file names inside the session's
// tree: none reaches outside, whatever the host.
func TestExoticPathsStayInside(t *testing.T) {
	skipOnWasm(t)
	dir := t.TempDir()
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "secret.txt"), []byte("outside\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, s := range []*eval.Session{eval.NewSession(eval.Options{}), eval.NewSession(eval.Options{FileRoot: dir})} {
		for _, p := range []string{`C:\\Windows\\win.ini`, `C:/Windows/win.ini`, `C:win.ini`, `\\\\server\\share\\x`, `\\\\?\\C:\\x`, `//server/share/x`, `/\\..\\..\\etc\\hosts`, `..\\..\\x`, filepath.ToSlash(filepath.Join(outside, "secret.txt"))} {
			if _, err := s.Eval(`readLines("` + p + `")`); err == nil {
				t.Errorf("%s was read", p)
			}
		}
		s.Close()
	}
}

// A symbolic link swapped between inside and outside while the session opens
// through it never lets a read out: os.Root resolves and opens in one step.
func TestSymlinkSwapRace(t *testing.T) {
	skipOnWasm(t)
	dir := t.TempDir()
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "f.txt"), []byte("OUTSIDE\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "in"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "in", "f.txt"), []byte("inside\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "swap")
	if os.Symlink(filepath.Join(dir, "in"), link) != nil {
		t.Skip("no symbolic links here")
	}
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; ; i++ {
			select {
			case <-stop:
				return
			default:
			}
			target := filepath.Join(dir, "in")
			if i%2 == 1 {
				target = outside
			}
			tmp := link + ".tmp"
			os.Remove(tmp)
			if os.Symlink(target, tmp) == nil {
				os.Rename(tmp, link)
			}
		}
	}()
	s := eval.NewSession(eval.Options{FileRoot: dir, MaxSteps: eval.Unbounded})
	defer s.Close()
	for i := 0; i < 3000; i++ {
		res, err := s.Eval(`tryCatch(readLines("/swap/f.txt", warn = FALSE), error = function(e) "err", warning = function(w) "warn")`)
		if err == nil && strings.Contains(res.Output, "OUTSIDE") {
			t.Fatal("read through a swapped link left the root")
		}
		if res != nil && res.Value != nil {
			if strings.Contains(res.Output, "OUTSIDE") {
				t.Fatal("left the root")
			}
		}
	}
	close(stop)
	<-done
}

func TestPathBombs(t *testing.T) {
	underMemoryLimit(t)
	s := eval.NewSession(eval.Options{})
	refused(t, s, `writeLines("x", strrep("a", 5000))`, "file name too long")
	refused(t, s, `dir.create(strrep("a/", 400), recursive = TRUE)`, "file name too long")
	// the number of files is bounded
	refused(t, eval.NewSession(eval.Options{MaxSteps: eval.Unbounded, Timeout: 20 * time.Second}),
		`d <- tempfile(); dir.create(d); i <- 0; repeat { i <- i + 1; if (!file.create(file.path(d, i), showWarnings = FALSE)) stop("files are bounded") }`, "files are bounded")
}

// Sys.sleep ends at the time limit and at a cancelled context.
func TestSleepHonoursLimits(t *testing.T) {
	skipOnWasm(t)
	s := eval.NewSession(eval.Options{Timeout: 200 * time.Millisecond})
	start := time.Now()
	refused(t, s, `Sys.sleep(1e6)`, "time limit")
	if time.Since(start) > 3*time.Second {
		t.Fatalf("slept past the limit: %v", time.Since(start))
	}
	s = eval.NewSession(eval.Options{Timeout: eval.Unbounded})
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	start = time.Now()
	if _, err := s.EvalContext(ctx, `Sys.sleep(1e6)`); !errors.Is(err, context.DeadlineExceeded) || time.Since(start) > 3*time.Second {
		t.Fatalf("context ignored by Sys.sleep: %v after %v", err, time.Since(start))
	}
}

// Untrusted input to readRDS and source is bounded: the document by the
// decoder's limits, the code by the session's parser limits.
func TestUntrustedDocuments(t *testing.T) {
	underMemoryLimit(t)
	dir := t.TempDir()
	deep := strings.Repeat("(", 20000) + "1" + strings.Repeat(")", 20000)
	if err := os.WriteFile(filepath.Join(dir, "deep.R"), []byte(deep), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "junk.rds"), []byte("\x1f\x8b not gzip at all"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := eval.NewSession(eval.Options{FileRoot: dir, Parse: syntax.Limits{MaxBytes: 1 << 14, MaxDepth: 100, MaxNodes: 5000}})
	defer s.Close()
	refused(t, s, `source("/deep.R")`, "limit", "deep", "nest")
	refused(t, s, `readRDS("/junk.rds")`, "unknown input format")
}
