package eval

import (
	"embed"
	"sort"
	"sync"

	"github.com/go-embedded-r/r/syntax"
)

// The base library written in R by us: the closures of base that are glue
// over the primitives. Each file under baselib/ is original work (BSD),
// written from R's documentation and checked against the oracle; see
// docs/m3-plan.md for the clean-room rule. The files are parsed once and
// evaluated into the base environment of every new session.

//go:embed baselib/*.R
var baselibFS embed.FS

var (
	baselibOnce  sync.Once
	baselibExprs []Value
	baselibNames []string
)

func loadBaselib() {
	baselibOnce.Do(func() {
		entries, err := baselibFS.ReadDir("baselib")
		if err != nil {
			panic(err)
		}
		for _, e := range entries {
			src, err := baselibFS.ReadFile("baselib/" + e.Name())
			if err != nil {
				panic(err)
			}
			f, err := syntax.ParseAll(string(src))
			if err != nil {
				panic("baselib/" + e.Name() + ": " + err.Error())
			}
			for _, n := range f.Exprs {
				baselibExprs = append(baselibExprs, FromSyntax(n))
				if call, ok := n.(*syntax.Call); ok {
					if fs, ok := call.Fun.(*syntax.Sym); ok && fs.Name == "<-" && len(call.Args) == 2 {
						if sym, ok := call.Args[0].Value.(*syntax.Sym); ok {
							baselibNames = append(baselibNames, sym.Name)
						}
					}
				}
			}
		}
		sort.Strings(baselibNames)
	})
}

// installBaselib evaluates the library into the session's base environment.
func (s *Session) installBaselib() {
	loadBaselib()
	for _, e := range baselibExprs {
		s.eval(e, s.Base)
	}
}
