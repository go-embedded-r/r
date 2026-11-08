package evalcheck

import (
	"fmt"
	"strings"

	"github.com/go-embedded-r/r/rds"
)

// dump prints a serialized value compactly, attributes included, for the
// message of a difference.
func dump(v rds.Value) string {
	var b strings.Builder
	dumpTo(&b, v, 0)
	return b.String()
}

func dumpMeta(b *strings.Builder, m rds.Meta, depth int) {
	for _, a := range m.Attrs {
		fmt.Fprintf(b, " @%s=", a.Name)
		dumpTo(b, a.Value, depth+1)
	}
}

func dumpTo(b *strings.Builder, v rds.Value, depth int) {
	if depth > 6 {
		b.WriteString("...")
		return
	}
	switch x := v.(type) {
	case *rds.Logical:
		fmt.Fprintf(b, "lgl%v", x.V)
		dumpMeta(b, x.Meta, depth)
	case *rds.Integer:
		fmt.Fprintf(b, "int%v", x.V)
		dumpMeta(b, x.Meta, depth)
	case *rds.Real:
		fmt.Fprintf(b, "dbl%v", x.V)
		dumpMeta(b, x.Meta, depth)
	case *rds.Character:
		b.WriteString("chr[")
		for i, s := range x.V {
			if i > 0 {
				b.WriteByte(' ')
			}
			if s.NA {
				b.WriteString("NA")
			} else {
				fmt.Fprintf(b, "%q", s.S)
			}
		}
		b.WriteString("]")
		dumpMeta(b, x.Meta, depth)
	case *rds.List:
		b.WriteString("list(")
		for i, e := range x.V {
			if i > 0 {
				b.WriteString(", ")
			}
			dumpTo(b, e, depth+1)
		}
		b.WriteString(")")
		dumpMeta(b, x.Meta, depth)
	case *rds.Sym:
		fmt.Fprintf(b, "sym(%s)", x.Name)
	default:
		fmt.Fprintf(b, "%T", v)
	}
}
