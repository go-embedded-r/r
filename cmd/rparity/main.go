// Command rparity compares the Go parser with the reference R 4.6.2 and
// prints the parity tables of the README.
//
//	rparity probe OP FILE       # print R's answer to OP for each snippet of FILE
//	rparity snippets            # parity on testdata/snippets.R, errors.R, edge.R
//	rparity corpus DIR          # parity on every .R file under DIR (never committed)
//	rparity bench DIR           # parse speed of the Go parser and of R's parse()
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var osExit = os.Exit

func main() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "rparity:", err)
		osExit(1)
	}
}

// RootCmd builds the command tree.
func RootCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "rparity", Short: "Compare the Go parser with R " + "4.6.2", SilenceUsage: true}
	cmd.AddCommand(probeCmd(), snippetsCmd(), corpusCmd(), numbersCmd(), mutateCmd(), benchCmd(), checkCmd(), primcountCmd())
	return cmd
}
