// Package repohygiene holds a test that refuses stray tracked files.
package repohygiene

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
)

const maxBlob = 1 << 20

// stray reports why a tracked file must not be committed, or "".
func stray(name string, size int64) string {
	base := path.Base(name)
	for _, pat := range []string{"*.test", ".*.atomicwrite-*", "*.orig", "*.rej"} {
		if ok, _ := path.Match(pat, base); ok {
			return "matches " + pat
		}
	}
	inTestdata := false
	for _, seg := range strings.Split(name, "/") {
		if seg == "testdata" {
			inTestdata = true
		}
	}
	if size > maxBlob && !inTestdata {
		return "over 1 MiB outside testdata"
	}
	return ""
}

func TestStrayRules(t *testing.T) {
	cases := []struct {
		name string
		size int64
		bad  bool
	}{
		{"eval.test", 1, true},
		{"eval/.slice24.go.atomicwrite-123", 1, true},
		{"a/b.orig", 1, true},
		{"a/b.rej", 1, true},
		{"a/big.bin", maxBlob + 1, true},
		{"testdata/big.bin", maxBlob + 1, false},
		{"eval/x.go", 10, false},
	}
	for _, c := range cases {
		if got := stray(c.name, c.size) != ""; got != c.bad {
			t.Errorf("stray(%q, %d) = %v, want %v", c.name, c.size, got, c.bad)
		}
	}
}

func TestNoStrayTrackedFiles(t *testing.T) {
	top, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		t.Skip("not a git checkout")
	}
	root := strings.TrimSpace(string(top))
	cmd := exec.Command("git", "ls-files", "-z")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		t.Skipf("git ls-files failed: %v", err)
	}
	for _, f := range bytes.Split(out, []byte{0}) {
		name := string(f)
		if name == "" {
			continue
		}
		fi, err := os.Lstat(root + "/" + name)
		if err != nil {
			continue // tracked but deleted in the working tree
		}
		if why := stray(name, fi.Size()); why != "" {
			t.Errorf("tracked file %s %s; git rm it", name, why)
		}
	}
}
