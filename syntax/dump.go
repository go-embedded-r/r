package syntax

import (
	"encoding/hex"
	"math"
	"strconv"
	"strings"
)

// Dump renders the expressions in the canonical line format the oracle's
// "tree" operation produces from R's own language objects (see
// internal/oracle/oracle.R): "top N", then one node per line, children
// indented by one "." per level; names and strings are hex (UTF-8 bytes),
// doubles are their IEEE bits. Two parses are the same tree exactly when
// their dumps are equal.
func Dump(exprs []Node) []string {
	d := &dumper{}
	d.emit(0, "top "+strconv.Itoa(len(exprs)))
	for _, e := range exprs {
		d.node(e, 1)
	}
	return d.lines
}

type dumper struct{ lines []string }

func (d *dumper) emit(depth int, s string) {
	d.lines = append(d.lines, strings.Repeat(".", depth)+s)
}

func hx(s string) string { return hex.EncodeToString([]byte(s)) }

func bits(f float64) string {
	b := math.Float64bits(f)
	var raw [8]byte
	for i := 0; i < 8; i++ {
		raw[i] = byte(b >> (56 - 8*i))
	}
	return hex.EncodeToString(raw[:])
}

func (d *dumper) node(n Node, depth int) {
	switch x := n.(type) {
	case nil:
		d.emit(depth, "missing")
	case *Sym:
		d.emit(depth, "sym "+hx(x.Name))
	case *Const:
		switch x.Kind {
		case Null:
			d.emit(depth, "null")
		case Logical:
			v := "F"
			if x.NA {
				v = "NA"
			} else if x.LogicalV {
				v = "T"
			}
			d.emit(depth, "lgl "+v)
		case Integer:
			if x.NA {
				d.emit(depth, "int NA")
			} else {
				d.emit(depth, "int "+strconv.Itoa(int(x.Int)))
			}
		case Double:
			d.emit(depth, "dbl "+bits(x.Double))
		case Complex:
			d.emit(depth, "cplx "+bits(x.Double)+" "+bits(x.Im))
		case String:
			if x.NA {
				d.emit(depth, "str NA")
			} else {
				d.emit(depth, "str "+hx(x.Str))
			}
		}
	case *Function:
		d.emit(depth, "fn")
		for _, p := range x.Params {
			d.emit(depth+1, "formal "+hx(p.Name))
			if p.Default == nil {
				d.emit(depth+2, "nodefault")
			} else {
				d.node(p.Default, depth+2)
			}
		}
		d.node(x.Body, depth+1)
	case *Call:
		d.emit(depth, "call")
		d.node(x.Fun, depth+1)
		for _, a := range x.Args {
			if a.HasName {
				d.emit(depth+1, "name "+hx(a.Name))
			}
			d.node(a.Value, depth+1)
		}
	}
}
