package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/go-embedded-r/r/internal/gen"
	"github.com/go-embedded-r/r/internal/oracle"
	"github.com/go-embedded-r/r/rds"
)

func probeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "probe OP FILE",
		Short: "print R's answer to OP (tree, tokens, parseerror) for each snippet of FILE",
		Args:  cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			snippets, err := gen.ReadSnippets(".", args[1])
			if err != nil {
				return err
			}
			p, err := oracle.Start(oracle.Options{})
			if err != nil {
				return err
			}
			defer p.Close()
			out := c.OutOrStdout()
			for _, s := range snippets {
				fmt.Fprintf(out, "=== %q\n", s)
				r, err := p.Call(args[0], s)
				if r == nil {
					return err
				}
				if !r.OK {
					fmt.Fprintf(out, "ERROR %s\n", r.Message)
					continue
				}
				if ch, ok := r.Value.(*rds.Character); ok {
					for _, l := range ch.V {
						fmt.Fprintln(out, l.S)
					}
				} else if r.Value == rds.Null {
					fmt.Fprintln(out, "NULL")
				}
			}
			return nil
		},
	}
}
