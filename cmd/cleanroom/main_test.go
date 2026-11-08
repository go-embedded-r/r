package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-embedded-r/r/internal/cleanroom"
)

func run(t *testing.T, args ...string) (string, error) {
	t.Helper()
	c := RootCmd()
	var out bytes.Buffer
	c.SetOut(&out)
	c.SetErr(&out)
	c.SetArgs(args)
	err := c.Execute()
	return out.String(), err
}

func TestCheckAndFloor(t *testing.T) {
	root, ok := cleanroom.RSourceDir()
	if !ok {
		if os.Getenv("RGO_REQUIRE_RSOURCE") != "" {
			t.Fatal("no R source tree")
		}
		t.Skip("no R source tree")
	}
	dir := t.TempDir()
	f := filepath.Join(dir, "own.R")
	body := "f <- function(a, b) {\n  s <- 0\n  for (i in seq_len(a)) s <- s + i * b\n  s\n}\n"
	if err := os.WriteFile(f, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := run(t, "check", "--rsrc", root, f)
	if err != nil || !strings.Contains(out, "own.R") {
		t.Fatalf("check: %v\n%s", err, out)
	}
	out, err = run(t, "floor", "--rsrc", root, f)
	if err != nil || !strings.Contains(out, "summary:") {
		t.Fatalf("floor: %v\n%s", err, out)
	}
	// a limit of 0 makes any file fail the run gate
	if _, err = run(t, "check", "--rsrc", root, "--max-exact-run", "1", f); err == nil {
		t.Fatal("a gate of 1 token must fail")
	}
}

func TestErrors(t *testing.T) {
	if _, err := run(t, "check", "--rsrc", t.TempDir()); err == nil {
		t.Fatal("an empty R tree must be an error")
	}
	if _, err := run(t, "floor"); err == nil {
		t.Fatal("floor needs files")
	}
}

func TestMainExit(t *testing.T) {
	old, oldArgs := osExit, os.Args
	defer func() { osExit, os.Args = old, oldArgs }()
	code := 0
	osExit = func(c int) { code = c }
	os.Args = []string{"cleanroom", "floor"}
	main()
	if code != 1 {
		t.Fatalf("exit code %d", code)
	}
}
