// Command roracle regenerates everything this repository derives from the
// reference R, and checks which R is installed.
//
//	roracle version   # which R answers, and whether it is the pinned release
//	roracle eval SRC  # evaluate SRC in R, print its captured output
//	roracle deparse   # testdata/golden/deparse.jsonl  (from testdata/snippets.R)
//	roracle rds       # testdata/golden/rds.jsonl      (from testdata/values.R)
//	roracle regen     # both goldens
//
// The reference is the release pinned in internal/oracle, run through the
// Rscript named by $RGO_RSCRIPT or found on PATH.
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/go-embedded-r/r/internal/gen"
	"github.com/go-embedded-r/r/internal/oracle"
	"github.com/go-embedded-r/r/rds"
)

var osExit = os.Exit

// startOracle is replaced by tests.
var startOracle = func() (oracleSession, error) {
	return oracle.Start(oracle.Options{})
}

// oracleSession is what the commands need from a running R.
type oracleSession interface {
	gen.Oracle
	Eval(src string) (*oracle.Reply, error)
	Greeting() oracle.Greeting
	Close() error
}

func main() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "roracle:", err)
		osExit(1)
	}
}

// RootCmd builds the command tree.
func RootCmd() *cobra.Command {
	var root string
	cmd := &cobra.Command{
		Use:          "roracle",
		Short:        "Regenerate the files derived from the reference R " + oracle.Version,
		SilenceUsage: true,
	}
	cmd.PersistentFlags().StringVarP(&root, "root", "C", ".", "repository root")

	steps := map[string]func(o oracleSession, root string, out io.Writer) error{
		"deparse": writeDeparse,
		"rds":     writeRDS,
	}
	run := func(names ...string) func(*cobra.Command, []string) error {
		return func(c *cobra.Command, _ []string) error {
			o, err := startOracle()
			if err != nil {
				return err
			}
			defer o.Close()
			for _, n := range names {
				if err := steps[n](o, root, c.OutOrStdout()); err != nil {
					return err
				}
			}
			return nil
		}
	}
	cmd.AddCommand(&cobra.Command{Use: "deparse", Short: "write testdata/golden/deparse.jsonl", Args: cobra.NoArgs, RunE: run("deparse")})
	cmd.AddCommand(&cobra.Command{Use: "rds", Short: "write testdata/golden/rds.jsonl", Args: cobra.NoArgs, RunE: run("rds")})
	cmd.AddCommand(&cobra.Command{Use: "regen", Short: "regenerate every golden file", Args: cobra.NoArgs, RunE: run("deparse", "rds")})
	cmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "print the identity of the R that answers",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			o, err := startOracle()
			if err != nil {
				return err
			}
			defer o.Close()
			g := o.Greeting()
			fmt.Fprintf(c.OutOrStdout(), "%s\nR %s (r%s) on %s\n", g.Banner, g.Version, g.SVN, g.Platform)
			return nil
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "eval SRC",
		Short: "evaluate SRC in R and print its captured output and warnings",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			o, err := startOracle()
			if err != nil {
				return err
			}
			defer o.Close()
			r, err := o.Eval(args[0])
			if r != nil {
				for _, l := range r.Output {
					fmt.Fprintln(c.OutOrStdout(), l)
				}
				for _, w := range r.Warnings {
					fmt.Fprintln(c.ErrOrStderr(), "warning:", w)
				}
			}
			return err
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "call OP ARG...",
		Short: "send one request to the oracle and print a character reply, one string per line",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			o, err := startOracle()
			if err != nil {
				return err
			}
			defer o.Close()
			vals := make([]any, len(args)-1)
			for i, a := range args[1:] {
				vals[i] = a
			}
			r, err := o.Call(args[0], vals...)
			if r != nil {
				if ch, ok := r.Value.(*rds.Character); ok {
					for _, s := range ch.V {
						fmt.Fprintln(c.OutOrStdout(), s.S)
					}
				}
			}
			return err
		},
	})
	return cmd
}

func write(root, rel string, data []byte, out io.Writer) error {
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return err
	}
	fmt.Fprintf(out, "wrote %s (%d bytes)\n", rel, len(data))
	return nil
}

func writeDeparse(o oracleSession, root string, out io.Writer) error {
	snippets, err := gen.ReadSnippets(root, "testdata/snippets.R")
	if err != nil {
		return err
	}
	data, err := gen.DeparseGolden(o, snippets)
	if err != nil {
		return err
	}
	return write(root, "testdata/golden/deparse.jsonl", data, out)
}

func writeRDS(o oracleSession, root string, out io.Writer) error {
	exprs, err := gen.ReadSnippets(root, "testdata/values.R")
	if err != nil {
		return err
	}
	data, err := gen.RDSGolden(o, exprs)
	if err != nil {
		return err
	}
	return write(root, "testdata/golden/rds.jsonl", data, out)
}
