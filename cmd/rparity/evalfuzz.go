package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/go-embedded-r/r/internal/evalcheck"
	"github.com/go-embedded-r/r/internal/gen"
	"github.com/go-embedded-r/r/internal/oracle"
)

var (
	fuzzNums = []string{"0", "1", "-1", "2", "0.5", "-2.5", "3", "10", "1e308", "-1e308", "Inf", "-Inf", "NaN", "NA_real_", "1e-300", "7", "100"}
	fuzzInts = []string{"0L", "1L", "-1L", "2L", "3L", "7L", "-7L", "2147483647L", "-2147483647L", "NA_integer_", "46341L", "100000L"}
	fuzzLgls = []string{"TRUE", "FALSE", "NA"}
	fuzzBin  = []string{"+", "-", "*", "/", "^", "%%", "%/%", "==", "!=", "<", "<=", ">", ">=", "&", "|", "&&", "||"}
	fuzzFns  = []string{"sum", "min", "max", "prod", "any", "all", "is.na", "abs", "length", "as.integer", "as.numeric", "as.logical", "sqrt", "floor", "ceiling", "seq_along", "c", "is.null", "is.numeric"}
)

func fuzzLeaf(r *rand.Rand) string {
	var pool []string
	switch r.Intn(4) {
	case 0:
		pool = fuzzNums
	case 1:
		pool = fuzzInts
	case 2:
		pool = fuzzLgls
	default:
		pool = fuzzNums
	}
	if r.Intn(5) == 0 {
		n := 2 + r.Intn(3)
		parts := make([]string, n)
		for i := range parts {
			parts[i] = pool[r.Intn(len(pool))]
		}
		return "c(" + strings.Join(parts, ", ") + ")"
	}
	return pool[r.Intn(len(pool))]
}

func fuzzExpr(r *rand.Rand, depth int) string {
	if depth == 0 || r.Intn(4) == 0 {
		return fuzzLeaf(r)
	}
	switch r.Intn(6) {
	case 0:
		return "-(" + fuzzExpr(r, depth-1) + ")"
	case 1:
		return "!(" + fuzzExpr(r, depth-1) + ")"
	case 2:
		f := fuzzFns[r.Intn(len(fuzzFns))]
		if r.Intn(4) == 0 {
			return f + "(" + fuzzExpr(r, depth-1) + ", " + fuzzExpr(r, depth-1) + ")"
		}
		return f + "(" + fuzzExpr(r, depth-1) + ")"
	case 3:
		return "if (" + fuzzExpr(r, depth-1) + ") " + fuzzExpr(r, depth-1) + " else " + fuzzExpr(r, depth-1)
	}
	op := fuzzBin[r.Intn(len(fuzzBin))]
	if op == "^" {
		// libm's pow and Go's differ in the last bit for most fractional
		// exponents; the exponents here are integers and 0.5
		exps := []string{"0", "1", "2", "3", "-1", "-2", "0.5", "NA", "NA_real_", "Inf", "-Inf", "2L", "NA_integer_"}
		return "(" + fuzzExpr(r, depth-1) + ") ^ (" + exps[r.Intn(len(exps))] + ")"
	}
	return "(" + fuzzExpr(r, depth-1) + ") " + op + " (" + fuzzExpr(r, depth-1) + ")"
}

// evalfuzzCmd evaluates random expressions with R and with Go and compares.
func evalfuzzCmd() *cobra.Command {
	var n int
	var seed int64
	var maxFail int
	var mode string
	cmd := &cobra.Command{
		Use:   "evalfuzz",
		Short: "differential fuzz of the evaluator: random arithmetic expressions, R against Go",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			p, err := oracle.Start(oracle.Options{})
			if err != nil {
				return err
			}
			defer p.Close()
			r := rand.New(rand.NewSource(seed))
			ok, bad, libm := 0, 0, 0
			perFn := map[string][2]int{}
			var failures []string
			const batch = 200
			for done := 0; done < n; done += batch {
				srcs := make([]string, batch)
				for i := range srcs {
					if mode == "math" {
						srcs[i] = fuzzMath(r)
						if strings.HasSuffix(srcs[i], "# libm") {
							libm++
						}
					} else if mode == "apply" {
						srcs[i] = fuzzApply(r)
					} else if mode == "regex" {
						srcs[i] = fuzzRegex(r)
					} else if mode == "strings" {
						srcs[i] = fuzzStrings(r)
					} else if mode == "subset" {
						srcs[i] = fuzzSubset(r)
					} else {
						srcs[i] = fuzzExpr(r, 3)
					}
				}
				data, err := gen.EvalGolden(p, srcs)
				if err != nil {
					return err
				}
				path := c.OutOrStdout()
				_ = path
				recs, err := gen.ParseJSONL[gen.EvalRecord](data)
				if err != nil {
					return err
				}
				for _, rec := range recs {
					fn := rec.Src
					if i := strings.IndexByte(fn, '('); i > 0 {
						fn = fn[:i]
					}
					t := perFn[fn]
					if why := evalcheck.Check(rec); why != "" {
						bad++
						t[1]++
						if len(failures) < maxFail {
							failures = append(failures, fmt.Sprintf("%s: %s", rec.Src, why))
						}
					} else {
						ok++
						t[0]++
					}
					perFn[fn] = t
				}
			}
			fmt.Fprintf(c.OutOrStdout(), "evalfuzz (seed %d): %d expressions, %d agree with R, %d differ\n", seed, ok+bad, ok, bad)
			if libm > 0 {
				fmt.Fprintf(c.OutOrStdout(), "  of which %d libm functions compared through round(x, 6): last-bit equality with the C library excluded\n", libm)
			}
			names := make([]string, 0, len(perFn))
			for k := range perFn {
				names = append(names, k)
			}
			sort.Strings(names)
			for _, k := range names {
				fmt.Fprintf(c.OutOrStdout(), "  %-10s %6d agree, %d differ\n", k, perFn[k][0], perFn[k][1])
			}
			for _, f := range failures {
				fmt.Fprintln(c.OutOrStdout(), "  DIFF", f)
			}
			if bad > 0 {
				return fmt.Errorf("%d expressions differ from R", bad)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&mode, "mode", "arith", "arith, subset, regex, strings, apply or math")
	cmd.Flags().IntVar(&n, "n", 2000, "expressions")
	cmd.Flags().Int64Var(&seed, "seed", 1, "random seed")
	cmd.Flags().IntVar(&maxFail, "max-diff", 15, "differences to print")
	return cmd
}
