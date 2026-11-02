package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"

	"github.com/go-embedded-r/r/internal/oracle"
	"github.com/go-embedded-r/r/internal/parity"
)

const noise = "()[]{}'\"`%,;=<>-+*/^$@:|&!?\\#_ \n\t.0aE1Lx_~"

// mutate returns src with k random edits: a deleted, duplicated, replaced,
// inserted or swapped character.
func mutate(rng *rand.Rand, src string, k int) string {
	r := []rune(src)
	for ; k > 0 && len(r) > 0; k-- {
		i := rng.Intn(len(r))
		switch rng.Intn(5) {
		case 0:
			r = append(r[:i], r[i+1:]...)
		case 1:
			r = append(r[:i+1], r[i:]...)
		case 2:
			r[i] = rune(noise[rng.Intn(len(noise))])
		case 3:
			r = append(r[:i], append([]rune{rune(noise[rng.Intn(len(noise))])}, r[i:]...)...)
		default:
			if i+1 < len(r) {
				r[i], r[i+1] = r[i+1], r[i]
			}
		}
	}
	return string(r)
}

// mutateCmd is a differential fuzz: random edits of real R code, parsed by
// both R and Go, verdicts compared (tree, tokens, error text).
func mutateCmd() *cobra.Command {
	var n int
	var seed int64
	var maxFail int
	cmd := &cobra.Command{
		Use:   "mutate DIR",
		Short: "differential fuzz: random edits of the R files under DIR",
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
			rng := rand.New(rand.NewSource(seed))
			t := &tally{name: fmt.Sprintf("mutants of %s (seed %d)", args[0], seed)}
			for t.n < n {
				b, err := os.ReadFile(files[rng.Intn(len(files))])
				if err != nil || !utf8.Valid(b) {
					continue
				}
				lines := strings.Split(string(b), "\n")
				if len(lines) < 3 {
					continue
				}
				// a window of a few lines
				start := rng.Intn(len(lines))
				end := min(len(lines), start+1+rng.Intn(6))
				src := mutate(rng, strings.Join(lines[start:end], "\n"), 1+rng.Intn(3))
				o, err := parity.Compare(p, src)
				if err != nil {
					return fmt.Errorf("%q: %w", src, err)
				}
				t.add(fmt.Sprintf("%q", src), o)
			}
			t.print(c.OutOrStdout(), maxFail)
			return nil
		},
	}
	cmd.Flags().IntVar(&n, "n", 20000, "mutants")
	cmd.Flags().Int64Var(&seed, "seed", 1, "random seed")
	cmd.Flags().IntVar(&maxFail, "max-diff", 15, "differences to print")
	return cmd
}
