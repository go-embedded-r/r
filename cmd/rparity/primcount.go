package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"

	"github.com/go-embedded-r/r/syntax"
)

// primcountCmd counts, over the R files under DIR, how often each primitive
// of R's base environment is called (a call site whose function is the bare
// name) and referenced as a value, and in how many files. It reads R's
// sources to count; nothing of them is kept.
func primcountCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "primcount DIR",
		Short: "count calls of each primitive in the R files under DIR, as a Markdown table",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			files, err := corpusFiles(args[0])
			if err != nil {
				return err
			}
			type count struct{ calls, refs, files int }
			counts := map[string]*count{}
			internals := map[string]int{}
			prims := syntax.PrimitiveNames()
			parsed := 0
			for _, f := range files {
				b, err := os.ReadFile(f)
				if err != nil {
					return err
				}
				if !utf8.Valid(b) {
					continue
				}
				file, err := syntax.ParseAll(string(b))
				if err != nil {
					continue
				}
				parsed++
				seen := map[string]bool{}
				var walk func(n syntax.Node, asCallee bool)
				note := func(name string, call bool) {
					if !prims[name] {
						return
					}
					cn := counts[name]
					if cn == nil {
						cn = &count{}
						counts[name] = cn
					}
					if call {
						cn.calls++
					} else {
						cn.refs++
					}
					if !seen[name] {
						seen[name] = true
						cn.files++
					}
				}
				walk = func(n syntax.Node, asCallee bool) {
					switch x := n.(type) {
					case *syntax.Sym:
						if !asCallee {
							note(x.Name, false)
						}
					case *syntax.Call:
						if s, ok := x.Fun.(*syntax.Sym); ok {
							if s.Name == ".Internal" && len(x.Args) == 1 {
								if in, ok := x.Args[0].Value.(*syntax.Call); ok && in.FunName() != "" {
									internals[in.FunName()]++
								}
							}
							note(s.Name, true)
						} else {
							walk(x.Fun, false)
						}
						for _, a := range x.Args {
							if a.Value != nil {
								walk(a.Value, false)
							}
						}
					case *syntax.Function:
						for _, p := range x.Params {
							if p.Default != nil {
								walk(p.Default, false)
							}
						}
						walk(x.Body, false)
					}
				}
				for _, e := range file.Exprs {
					walk(e, false)
				}
			}
			type row struct {
				name string
				c    count
			}
			var rows []row
			for n := range prims {
				cn := counts[n]
				if cn == nil {
					cn = &count{}
				}
				rows = append(rows, row{n, *cn})
			}
			sort.Slice(rows, func(i, j int) bool {
				if rows[i].c.calls != rows[j].c.calls {
					return rows[i].c.calls > rows[j].c.calls
				}
				return rows[i].name < rows[j].name
			})
			out := c.OutOrStdout()
			fmt.Fprintf(out, "%d files parsed, %d primitives\n\n| primitive | calls | as value | files |\n|---|---|---|---|\n", parsed, len(rows))
			for _, r := range rows {
				fmt.Fprintf(out, "| `%s` | %d | %d | %d |\n", strings.ReplaceAll(r.name, "|", `\|`), r.c.calls, r.c.refs, r.c.files)
			}
			type irow struct {
				name string
				n    int
			}
			var irows []irow
			for n, k := range internals {
				irows = append(irows, irow{n, k})
			}
			sort.Slice(irows, func(i, j int) bool {
				if irows[i].n != irows[j].n {
					return irows[i].n > irows[j].n
				}
				return irows[i].name < irows[j].name
			})
			fmt.Fprintf(out, "\n%d distinct `.Internal` names called\n\n| .Internal | call sites |\n|---|---|\n", len(irows))
			for _, r := range irows {
				fmt.Fprintf(out, "| `%s` | %d |\n", r.name, r.n)
			}
			return nil
		},
	}
}
