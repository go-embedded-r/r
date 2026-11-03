package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"

	"github.com/go-embedded-r/r/internal/gen"
	"github.com/go-embedded-r/r/internal/oracle"
	"github.com/go-embedded-r/r/internal/parity"
)

// tally accumulates verdicts.
type tally struct {
	name                     string
	n, parsed, rejected      int
	tree, tokens, errs       int // matches
	treeOf, tokensOf, errsOf int // denominators
	dep                      int // deparse matches (denominator: parsed)
	ulp                      int
	failures                 []string
	skipped                  int
	// noDeparse: the file holds platform-dependent numeric literals (R formats
	// them with the C library and, on x86-64 Linux, in long double), so the
	// deparse verdict is reported but not part of the parity gate.
	noDeparse bool
}

func (t *tally) add(label string, o parity.Outcome) {
	t.n++
	if o.Parsed {
		t.parsed++
		t.treeOf++
		t.tokensOf++
		if o.Tree {
			t.tree++
		}
		if o.Tokens {
			t.tokens++
		}
		if o.Deparse {
			t.dep++
		}
	} else {
		t.rejected++
		t.errsOf++
		if o.Error {
			t.errs++
		}
	}
	t.ulp += o.ULP
	ok := o.Tree && o.Tokens && o.Error && (o.Deparse || t.noDeparse)
	if !ok {
		t.failures = append(t.failures, label+": "+o.Detail)
	}
}

func (t *tally) print(w io.Writer, maxFail int) {
	fmt.Fprintf(w, "%-28s %6d inputs (%d parse, %d rejected)\n", t.name, t.n, t.parsed, t.rejected)
	fmt.Fprintf(w, "  tree    %5d / %-5d\n", t.tree, t.treeOf)
	fmt.Fprintf(w, "  tokens  %5d / %-5d\n", t.tokens, t.tokensOf)
	if t.noDeparse {
		fmt.Fprintf(w, "  deparse of platform-dependent literals, not gated: %d of %d equal to this platform's R\n", t.dep, t.parsed)
	} else {
		fmt.Fprintf(w, "  deparse %5d / %-5d byte for byte (deparse() and print())\n", t.dep, t.parsed)
	}
	fmt.Fprintf(w, "  errors  %5d / %-5d byte for byte\n", t.errs, t.errsOf)
	if t.ulp > 0 {
		fmt.Fprintf(w, "  numeric literals within 16 ulp of R's (counted as equal): %d\n", t.ulp)
	}
	if t.skipped > 0 {
		fmt.Fprintf(w, "  skipped (not valid UTF-8): %d\n", t.skipped)
	}
	for i, f := range t.failures {
		if i >= maxFail {
			fmt.Fprintf(w, "  ... %d more differences\n", len(t.failures)-maxFail)
			break
		}
		fmt.Fprintf(w, "  DIFF %s\n", f)
	}
}

func snippetsCmd() *cobra.Command {
	var maxFail int
	cmd := &cobra.Command{
		Use:   "snippets",
		Short: "parity on the committed inputs (testdata/*.R)",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			p, err := oracle.Start(oracle.Options{})
			if err != nil {
				return err
			}
			defer p.Close()
			for _, f := range []string{"testdata/snippets.R", "testdata/errors.R", "testdata/edge.R", "testdata/deparse.R", "testdata/numbers.R", "testdata/values.R"} {
				snippets, err := gen.ReadSnippets(".", f)
				if err != nil {
					return err
				}
				if f == "testdata/values.R" {
					snippets = append(snippets, parity.Extras()...)
				}
				t := &tally{name: f, noDeparse: f == "testdata/numbers.R"}
				for i, s := range snippets {
					o, err := parity.Compare(p, s)
					if err != nil {
						return fmt.Errorf("%s #%d: %w", f, i, err)
					}
					t.add(fmt.Sprintf("#%d %q", i, s), o)
				}
				t.print(c.OutOrStdout(), maxFail)
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&maxFail, "max-diff", 10, "differences to print per file")
	return cmd
}

func corpusFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && (strings.HasSuffix(path, ".R") || strings.HasSuffix(path, ".r")) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func corpusCmd() *cobra.Command {
	var maxFail int
	cmd := &cobra.Command{
		Use:   "corpus DIR",
		Short: "parity on every .R file under DIR (read locally, never committed)",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			files, err := corpusFiles(args[0])
			if err != nil {
				return err
			}
			p, err := oracle.Start(oracle.Options{})
			if err != nil {
				return err
			}
			defer p.Close()
			t := &tally{name: "corpus " + args[0]}
			bytesTotal := 0
			for _, path := range files {
				b, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				if !utf8.Valid(b) {
					t.skipped++
					continue
				}
				src := strings.TrimSuffix(string(b), "\n")
				bytesTotal += len(b)
				o, err := parity.Compare(p, src)
				if err != nil {
					return fmt.Errorf("%s: %w", path, err)
				}
				t.add(path, o)
			}
			fmt.Fprintf(c.OutOrStdout(), "%d files, %.1f MB\n", len(files), float64(bytesTotal)/1e6)
			t.print(c.OutOrStdout(), maxFail)
			return nil
		},
	}
	cmd.Flags().IntVar(&maxFail, "max-diff", 10, "differences to print")
	return cmd
}
