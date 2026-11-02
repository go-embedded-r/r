package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/spf13/cobra"

	"github.com/go-embedded-r/r/internal/oracle"
	"github.com/go-embedded-r/r/internal/parity"
	"github.com/go-embedded-r/r/syntax"
)

// numbersCmd measures how the Go parser's numeric literals compare with the
// reference R's, on random decimal literals: exact matches, a few ulps off,
// and worse.
func numbersCmd() *cobra.Command {
	var n int
	var seed int64
	cmd := &cobra.Command{
		Use:   "numbers",
		Short: "compare decimal literal values with R's on random literals",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			p, err := oracle.Start(oracle.Options{})
			if err != nil {
				return err
			}
			defer p.Close()
			rng := rand.New(rand.NewSource(seed))
			var lits []string
			for len(lits) < n {
				lits = append(lits, randomLiteral(rng))
			}
			exact, close, far := 0, 0, 0
			var bad []string
			const batch = 500
			for i := 0; i < len(lits); i += batch {
				j := min(i+batch, len(lits))
				src := strings.Join(lits[i:j], "\n")
				rl, rerr, err := parity.TreeLines(p, src)
				if err != nil {
					return err
				}
				if rerr != "" {
					return fmt.Errorf("R: %s", rerr)
				}
				f, gerr := syntax.ParseAll(src)
				if gerr != nil {
					return gerr
				}
				gl := syntax.Dump(f.Exprs)
				for k := range lits[i:j] {
					r, g := rl[k+1], gl[k+1]
					switch {
					case r == g:
						exact++
					default:
						if _, ok := parity.NumClose(r, g); ok {
							close++
						} else {
							far++
						}
						if len(bad) < 15 {
							bad = append(bad, fmt.Sprintf("%s  R=%s Go=%s", lits[i+k], strings.TrimPrefix(r, "."), strings.TrimPrefix(g, ".")))
						}
					}
				}
			}
			fmt.Fprintf(c.OutOrStdout(), "%d literals: %d identical, %d within 16 ulp, %d further\n", len(lits), exact, close, far)
			for _, b := range bad {
				fmt.Fprintln(c.OutOrStdout(), "  ", b)
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&n, "n", 20000, "literals")
	cmd.Flags().Int64Var(&seed, "seed", 1, "random seed")
	return cmd
}

func randomLiteral(rng *rand.Rand) string {
	var b strings.Builder
	nd := 1 + rng.Intn(20)
	if rng.Intn(4) == 0 {
		nd = 1 + rng.Intn(5)
	}
	for i := 0; i < nd; i++ {
		b.WriteByte(byte('0' + rng.Intn(10)))
	}
	s := strings.TrimLeft(b.String(), "0")
	if s == "" {
		s = "1"
	}
	if rng.Intn(2) == 0 && len(s) > 1 {
		k := 1 + rng.Intn(len(s)-1)
		s = s[:k] + "." + s[k:]
	}
	switch rng.Intn(3) {
	case 0:
	case 1:
		s += fmt.Sprintf("e%d", rng.Intn(60)-30)
	default:
		s += fmt.Sprintf("e%d", rng.Intn(660)-330)
	}
	return s
}
