package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/go-embedded-r/r/eval"
)

// The corpus is built by R code both sides evaluate, outside the timed part.
const benchSetup = `i <- 1:%d; x <- paste0("id", i, " user", i %%%% 97, "@host", i %%%% 13, ".example.org ", i * 7, " abc", i %%%% 1000, "-xyz")`

var benchExprs = []struct{ name, code string }{
	{"grepl", `sum(grepl("[0-9]+@host[0-9]+\\.example", x))`},
	{"gsub-class", `sum(nchar(gsub("[aeiou]", "_", x)))`},
	{"gsub-group", `sum(nchar(gsub("([a-z]+)([0-9]+)", "\\2\\1", x)))`},
	{"regexpr", `sum(regexpr("abc[0-9]+", x))`},
	{"backref", `sum(grepl("([0-9])\\1", x))`},
}

// regexbenchCmd times the port's regex functions against R's own, in
// interleaved rounds, one core each. R is timed inside its process (startup
// excluded). It refuses to measure when the 1-minute load is above 8.
func regexbenchCmd() *cobra.Command {
	var rounds, n int
	var maxLoad float64
	cmd := &cobra.Command{
		Use:   "regexbench",
		Short: "regex speed: this port (TRE port) against R's grepl/gsub/regexpr on one corpus",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			out, err := exec.Command("uptime").Output()
			if err != nil {
				return err
			}
			f := strings.Fields(strings.ReplaceAll(string(out), ",", " "))
			var load float64
			for i, w := range f {
				if w == "averages:" || w == "average:" {
					load, _ = strconv.ParseFloat(f[i+1], 64)
				}
			}
			if load > maxLoad {
				return fmt.Errorf("1-minute load %.1f is above %.1f: measurement deferred", load, maxLoad)
			}
			rscript := os.Getenv("RGO_RSCRIPT")
			if rscript == "" {
				return fmt.Errorf("RGO_RSCRIPT is not set")
			}
			setup := fmt.Sprintf(benchSetup, n)
			tmp, err := os.MkdirTemp("", "regexbench")
			if err != nil {
				return err
			}
			defer os.RemoveAll(tmp)
			type res struct{ g, r []float64 }
			results := map[string]*res{}
			for round := 0; round < rounds+1; round++ {
				for _, e := range benchExprs {
					s := eval.NewSession(eval.Options{MaxSteps: eval.Unbounded, Timeout: eval.Unbounded, MaxAlloc: eval.Unbounded})
					if _, err := s.Eval(setup); err != nil {
						return err
					}
					t0 := time.Now()
					if _, err := s.Eval(e.code); err != nil {
						return fmt.Errorf("%s: %v", e.name, err)
					}
					g := time.Since(t0).Seconds()
					script := setup + "\nt0 <- proc.time()[[3]]; invisible(" + e.code + "); cat(proc.time()[[3]] - t0)\n"
					sf := filepath.Join(tmp, "b.R")
					if err := os.WriteFile(sf, []byte(script), 0o600); err != nil {
						return err
					}
					rout, err := exec.Command(rscript, sf).CombinedOutput()
					if err != nil {
						return fmt.Errorf("Rscript: %v: %s", err, rout)
					}
					r, _ := strconv.ParseFloat(strings.TrimSpace(string(rout)), 64)
					if round == 0 {
						continue // warm-up round
					}
					if results[e.name] == nil {
						results[e.name] = &res{}
					}
					results[e.name].g = append(results[e.name].g, g)
					results[e.name].r = append(results[e.name].r, r)
				}
			}
			med := func(v []float64) float64 { sort.Float64s(v); return v[len(v)/2] }
			fmt.Fprintf(c.OutOrStdout(), "regexbench: %d lines, %d rounds, load %.1f\n", n, rounds, load)
			for _, e := range benchExprs {
				r := results[e.name]
				g, rr := med(r.g), med(r.r)
				fmt.Fprintf(c.OutOrStdout(), "  %-11s go %8.1f ms   R %8.1f ms   go/R %5.2f\n", e.name, g*1000, rr*1000, g/rr)
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&rounds, "rounds", 5, "interleaved rounds")
	cmd.Flags().IntVar(&n, "n", 20000, "lines")
	cmd.Flags().Float64Var(&maxLoad, "max-load", 8, "refuse to measure above this 1-minute load")
	return cmd
}
