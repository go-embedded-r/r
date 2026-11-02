package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/spf13/cobra"

	"github.com/go-embedded-r/r/internal/oracle"
	"github.com/go-embedded-r/r/rds"
	"github.com/go-embedded-r/r/syntax"
)

// benchCmd times the Go parser and R's parse() on the same texts, in
// interleaved rounds, one core each.
func benchCmd() *cobra.Command {
	var rounds int
	cmd := &cobra.Command{
		Use:   "bench DIR",
		Short: "parse speed: Go against R's parse(text=, keep.source=FALSE) on the .R files under DIR",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			files, err := corpusFiles(args[0])
			if err != nil {
				return err
			}
			var texts []string
			total := 0
			for _, f := range files {
				b, err := os.ReadFile(f)
				if err != nil {
					return err
				}
				s := string(b)
				if !utf8.ValidString(s) {
					continue
				}
				if _, err := syntax.ParseAll(s); err != nil {
					continue // only inputs both parse: the comparison is of equal work
				}
				texts = append(texts, strings.TrimSuffix(s, "\n"))
				total += len(s)
			}
			p, err := oracle.Start(oracle.Options{Timeout: 10 * time.Minute})
			if err != nil {
				return err
			}
			defer p.Close()
			vec := rds.Strings(texts...)
			var goT, rT []float64
			for i := 0; i < rounds+1; i++ {
				t0 := time.Now()
				for _, s := range texts {
					syntax.ParseAll(s)
				}
				g := time.Since(t0).Seconds()
				r, err := p.Call("benchparse", vec)
				if err != nil {
					return err
				}
				rv := r.Value.(*rds.Real).V[0]
				if i == 0 {
					continue // warm-up round
				}
				goT, rT = append(goT, g), append(rT, rv)
			}
			sort.Float64s(goT)
			sort.Float64s(rT)
			g, r := goT[len(goT)/2], rT[len(rT)/2]
			fmt.Fprintf(c.OutOrStdout(), "%d files, %.1f MB, %d rounds (median)\n  Go  %.3f s\n  R   %.3f s\n  Go/R = %.3f (%.1fx faster)\n", len(texts), float64(total)/1e6, rounds, g, r, g/r, r/g)
			return nil
		},
	}
	cmd.Flags().IntVar(&rounds, "rounds", 7, "timed rounds, after one warm-up round")
	return cmd
}
