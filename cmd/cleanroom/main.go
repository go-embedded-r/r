// Command cleanroom scores our R files against R's own library code, to
// enforce the clean-room rule (NOTICE): token k-gram winnowing, exact and with
// identifiers normalised, plus the longest common token run. It reads R's GPL
// sources at run time and prints only file:line references.
//
//	cleanroom check [files...]   # default: eval/baselib/*.R; exit 1 over the gate
//	cleanroom floor files...     # same scores for independent code (noise floor)
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/go-embedded-r/r/internal/cleanroom"
	"github.com/spf13/cobra"
)

var osExit = os.Exit

func main() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		osExit(1)
	}
}

// RootCmd builds the command tree.
func RootCmd() *cobra.Command {
	var rsrc string
	var k, w int
	var lim = cleanroom.Gate
	root := &cobra.Command{Use: "cleanroom", Short: "similarity of our R files to R's library code", SilenceUsage: true}
	root.PersistentFlags().StringVar(&rsrc, "rsrc", "", "R 4.6.2 source tree (default $RGO_RSOURCE, else the durable checkout)")
	root.PersistentFlags().IntVarP(&k, "k", "k", cleanroom.Default.K, "k-gram length in tokens (noise threshold)")
	root.PersistentFlags().IntVarP(&w, "window", "w", cleanroom.Default.W, "winnowing window")
	root.PersistentFlags().Float64Var(&lim.Containment, "max-containment", lim.Containment, "gate: normalised containment limit")
	root.PersistentFlags().Float64Var(&lim.ExactContainment, "max-exact-containment", lim.ExactContainment, "gate: exact containment limit")
	root.PersistentFlags().IntVar(&lim.ExactRun, "max-exact-run", lim.ExactRun, "gate: longest exact token run")
	root.PersistentFlags().IntVar(&lim.NormRun, "max-norm-run", lim.NormRun, "gate: longest normalised token run")

	var list int
	root.PersistentFlags().IntVar(&list, "list", 0, "also list every common exact run of at least this many code tokens (file:line references only)")
	run := func(gate bool) func(*cobra.Command, []string) error {
		return func(cmd *cobra.Command, args []string) error {
			p := cleanroom.Params{K: k, W: w}
			if rsrc == "" {
				d, ok := cleanroom.RSourceDir()
				if !ok {
					return fmt.Errorf("no R source tree: pass --rsrc or set RGO_RSOURCE")
				}
				rsrc = d
			}
			refs, err := cleanroom.LoadRefs(rsrc, p)
			if err != nil {
				return err
			}
			if len(args) == 0 {
				args, _ = filepath.Glob("eval/baselib/*.R")
			}
			if len(args) == 0 {
				return fmt.Errorf("no input files")
			}
			var res []cleanroom.Result
			for _, f := range args {
				b, err := os.ReadFile(f)
				if err != nil {
					return err
				}
				res = append(res, cleanroom.Check(f, b, refs, p))
				if list > 0 {
					fmt.Fprint(cmd.OutOrStdout(), cleanroom.ListRuns(f, b, refs, p, list))
				}
			}
			sort.SliceStable(res, func(i, j int) bool { return res[i].Norm > res[j].Norm })
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "k=%d w=%d (guarantee t=%d tokens) refs=%d\n", p.K, p.W, p.K+p.W-1, len(refs))
			fmt.Fprintf(out, "%-32s %6s %-34s %6s %6s %6s %6s  %s | %s\n", "file", "tokens", "best R file (normalised)", "exact", "norm", "jacE", "jacN", "longest exact run (any R file)", "longest normalised run")
			bad := 0
			for _, r := range res {
				fmt.Fprintf(out, "%-32s %6d %-34s %6.2f %6.2f %6.2f %6.2f  %d tok %s:%d-%d vs %s:%d-%d | %d tok %s:%d-%d vs %s:%d-%d\n",
					r.File, r.Tokens, r.Best, r.Exact, r.Norm, r.ExactJac, r.NormJac,
					r.LongestAny.Len, r.File, r.LongestAny.OurStart, r.LongestAny.OurEnd, r.LongestAnyFile, r.LongestAny.TheirStart, r.LongestAny.TheirEnd,
					r.LongestAnyNorm.Len, r.File, r.LongestAnyNorm.OurStart, r.LongestAnyNorm.OurEnd, r.LongestAnyNormFile, r.LongestAnyNorm.TheirStart, r.LongestAnyNorm.TheirEnd)
				if v := lim.Violations(r); len(v) > 0 && gate {
					bad++
					for _, s := range v {
						fmt.Fprintf(out, "  FAIL %s: %s\n", r.File, s)
					}
				}
			}
			if !gate {
				fmt.Fprint(out, cleanroom.Summary(res, 200))
			}
			if bad > 0 {
				return fmt.Errorf("%d file(s) over the clean-room gate", bad)
			}
			return nil
		}
	}
	root.AddCommand(&cobra.Command{Use: "check [files...]", Short: "score files (default eval/baselib/*.R) and fail over the gate", RunE: run(true)})
	root.AddCommand(&cobra.Command{Use: "floor files...", Short: "score independent code: the noise floor", Args: cobra.MinimumNArgs(1), RunE: run(false)})
	return root
}
