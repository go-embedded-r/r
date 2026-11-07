package main

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/go-embedded-r/r/eval"
	"github.com/go-embedded-r/r/syntax"
)

// missingCmd lists the primitives of R that eval does not implement yet.
func missingCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "missing",
		Short: "list the primitives of R that package eval does not implement",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			have := map[string]bool{}
			for _, n := range eval.Builtins() {
				have[n] = true
			}
			for _, n := range eval.Unsupported() {
				have[n] = false // present as a name, refuses to run
			}
			var miss []string
			total := 0
			for n := range syntax.PrimitiveNames() {
				total++
				if !have[n] {
					miss = append(miss, n)
				}
			}
			sort.Strings(miss)
			fmt.Fprintf(c.OutOrStdout(), "%d of %d primitives missing\n", len(miss), total)
			for _, n := range miss {
				fmt.Fprintln(c.OutOrStdout(), n)
			}
			return nil
		},
	}
}
