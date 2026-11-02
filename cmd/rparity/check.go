package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// checkCmd reads a report written by the other commands and fails unless
// every "tree", "tokens" and "errors" line is n / n. (A verdict is data, so
// it is computed here, not by a shell pipeline.)
func checkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check REPORT",
		Short: "fail unless every verdict line of a report is n / n",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			f, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer f.Close()
			sc := bufio.NewScanner(f)
			sc.Buffer(make([]byte, 1<<20), 1<<20)
			lines, bad := 0, 0
			for sc.Scan() {
				fs := strings.Fields(sc.Text())
				if len(fs) >= 4 && fs[2] == "/" && (fs[0] == "tree" || fs[0] == "tokens" || fs[0] == "errors") {
					lines++
					if fs[1] != fs[3] {
						bad++
						fmt.Fprintln(c.ErrOrStderr(), "not at parity:", sc.Text())
					}
				}
			}
			if err := sc.Err(); err != nil {
				return err
			}
			if lines == 0 {
				return fmt.Errorf("%s holds no verdict line", args[0])
			}
			if bad > 0 {
				return fmt.Errorf("%d of %d verdicts differ from R", bad, lines)
			}
			fmt.Fprintf(c.OutOrStdout(), "%d verdicts, all at parity\n", lines)
			return nil
		},
	}
}
