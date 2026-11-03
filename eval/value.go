// Package eval is the evaluator of the Go R interpreter: R values, environments,
// promises, closures and the primitives, written from the R Language
// Definition and checked against R 4.6.2 (see docs/m2-design.md). Nothing is
// translated from R's sources.
package eval

import (
	"math"
	"runtime"
	"sync"

	"github.com/go-embedded-r/r/rds"
	"github.com/go-embedded-r/r/syntax"
)

// Type is a SEXPTYPE code.
type Type = rds.Type

// Value is an R object.
type Value interface{ Type() Type }

// NA values, as R stores them.
const (
	naInt = math.MinInt32
)

var naReal = math.Float64frombits(rds.NAReal)

// rNaN is R_NaN: the NaN the hardware makes of 0/0 (the sign bit is set on
// x86, clear on arm64), which is what R serializes for it.
var rNaN = func() float64 {
	if runtime.GOARCH == "amd64" || runtime.GOARCH == "386" {
		return math.Float64frombits(0xFFF8000000000000)
	}
	return math.Float64frombits(0x7FF8000000000000)
}()

// isNAReal: R's ISNA, the NaN whose low word is 1954.
func isNAReal(x float64) bool {
	return math.IsNaN(x) && math.Float64bits(x)&0xFFFFFFFF == 1954
}

// Attr is one attribute; the slice keeps R's order.
type Attr struct {
	Name  string
	Value Value
}

// Str is an element of a character vector.
type Str struct {
	S  string
	NA bool
}

type nullT struct{}

func (nullT) Type() Type { return rds.NILSXP }

// Null is R's NULL.
var Null Value = nullT{}

type missingT struct{}

func (missingT) Type() Type { return rds.SYMSXP }

// Missing is the empty argument (R_MissingArg).
var Missing Value = missingT{}

// Symbol is an interned name.
type Symbol struct{ Name string }

func (*Symbol) Type() Type { return rds.SYMSXP }

var symtab sync.Map

// Intern returns the one Symbol with this name.
func Intern(name string) *Symbol {
	if s, ok := symtab.Load(name); ok {
		return s.(*Symbol)
	}
	s, _ := symtab.LoadOrStore(name, &Symbol{Name: name})
	return s.(*Symbol)
}

// The vectors. Their slices are never modified in place once a value is
// visible to R code (copy on modify; see docs/m2-design.md).
type (
	Logical struct {
		V    []int32
		Attr []Attr
	}
	Integer struct {
		V    []int32
		Attr []Attr
	}
	Double struct {
		V    []float64
		Attr []Attr
	}
	String struct {
		V    []Str
		Attr []Attr
	}
	List struct {
		V    []Value
		Attr []Attr
	}
)

func (*Logical) Type() Type { return rds.LGLSXP }
func (*Integer) Type() Type { return rds.INTSXP }
func (*Double) Type() Type  { return rds.REALSXP }
func (*String) Type() Type  { return rds.STRSXP }
func (*List) Type() Type    { return rds.VECSXP }

// Arg is one argument of a call.
type Arg struct {
	Name    string
	HasName bool
	Value   Value // Missing for an empty argument
}

// Lang is a call.
type Lang struct {
	Fun  Value
	Args []Arg
}

func (*Lang) Type() Type { return rds.LANGSXP }

// Formal is one formal argument of a closure.
type Formal struct {
	Name    *Symbol
	Default Value // nil when there is none
}

// Formals is the formal-argument list of a `function` expression, which sits
// in the second place of the call to `function`.
type Formals struct{ F []Formal }

func (*Formals) Type() Type { return rds.LISTSXP }

// Closure is an R function.
type Closure struct {
	Formals []Formal
	Body    Value
	Env     *Env
}

func (*Closure) Type() Type { return rds.CLOSXP }

// Builtin is a primitive: Special ones get their arguments unevaluated.
type Builtin struct {
	Name    string
	Special bool
	fn      func(s *Session, c *callInfo) Value
}

func (b *Builtin) Type() Type {
	if b.Special {
		return rds.SPECIALSXP
	}
	return rds.BUILTINSXP
}

// Promise is a delayed argument.
type Promise struct {
	Expr    Value
	Env     *Env
	val     Value
	forced  bool
	forcing bool
}

func (*Promise) Type() Type { return rds.PROMSXP }

// Dots is the value of `...`: the promises of the extra arguments.
type Dots struct {
	Args []Arg // Value is a *Promise, a constant, or Missing
}

func (*Dots) Type() Type { return rds.DOTSXP }

// ---- constructors and accessors -------------------------------------------

func mkLogical(b bool) *Logical {
	if b {
		return &Logical{V: []int32{1}}
	}
	return &Logical{V: []int32{0}}
}
func mkInt(i int32) *Integer     { return &Integer{V: []int32{i}} }
func mkDouble(f float64) *Double { return &Double{V: []float64{f}} }
func mkString(s string) *String  { return &String{V: []Str{{S: s}}} }

// Length is length(x) for the types of this slice.
func Length(v Value) int {
	switch x := v.(type) {
	case *Logical:
		return len(x.V)
	case *Integer:
		return len(x.V)
	case *Double:
		return len(x.V)
	case *String:
		return len(x.V)
	case *List:
		return len(x.V)
	case *Lang:
		return len(x.Args) + 1
	case nullT:
		return 0
	}
	return 1
}

func attrsOf(v Value) []Attr {
	switch x := v.(type) {
	case *Logical:
		return x.Attr
	case *Integer:
		return x.Attr
	case *Double:
		return x.Attr
	case *String:
		return x.Attr
	case *List:
		return x.Attr
	}
	return nil
}

func getAttr(v Value, name string) Value {
	for _, a := range attrsOf(v) {
		if a.Name == name {
			return a.Value
		}
	}
	return nil
}

func namesOf(v Value) *String {
	if n, ok := getAttr(v, "names").(*String); ok {
		return n
	}
	return nil
}

// withAttr returns a shallow copy of vector v with attributes attr.
func withAttr(v Value, attr []Attr) Value {
	switch x := v.(type) {
	case *Logical:
		return &Logical{V: x.V, Attr: attr}
	case *Integer:
		return &Integer{V: x.V, Attr: attr}
	case *Double:
		return &Double{V: x.V, Attr: attr}
	case *String:
		return &String{V: x.V, Attr: attr}
	case *List:
		return &List{V: x.V, Attr: attr}
	}
	return v
}

func namesAttr(n []Str) []Attr {
	return []Attr{{Name: "names", Value: &String{V: n}}}
}

// ---- conversions to and from the other packages ---------------------------

// FromSyntax converts a parsed expression to a language object.
func FromSyntax(n syntax.Node) Value {
	switch x := n.(type) {
	case nil:
		return Missing
	case *syntax.Sym:
		return Intern(x.Name)
	case *syntax.Const:
		switch x.Kind {
		case syntax.Null:
			return Null
		case syntax.Logical:
			switch {
			case x.NA:
				return &Logical{V: []int32{naInt}}
			case x.LogicalV:
				return mkLogical(true)
			}
			return mkLogical(false)
		case syntax.Integer:
			if x.NA {
				return &Integer{V: []int32{naInt}}
			}
			return mkInt(x.Int)
		case syntax.Double:
			if x.NA {
				return mkDouble(naReal)
			}
			return mkDouble(x.Double)
		case syntax.String:
			if x.NA {
				return &String{V: []Str{{NA: true}}}
			}
			return mkString(x.Str)
		}
		return Null // complex constants are not in this slice
	case *syntax.Function:
		f := &Formals{}
		for _, p := range x.Params {
			fm := Formal{Name: Intern(p.Name)}
			if p.Default != nil {
				fm.Default = FromSyntax(p.Default)
			}
			f.F = append(f.F, fm)
		}
		return &Lang{Fun: Intern("function"), Args: []Arg{{Value: f}, {Value: FromSyntax(x.Body)}}}
	case *syntax.Call:
		l := &Lang{Fun: FromSyntax(x.Fun)}
		for _, a := range x.Args {
			l.Args = append(l.Args, Arg{Name: a.Name, HasName: a.HasName, Value: FromSyntax(a.Value)})
		}
		return l
	}
	return Null
}

// toSyntax converts a language object back, for deparse.
func toSyntax(v Value) syntax.Node {
	switch x := v.(type) {
	case missingT:
		return nil
	case *Symbol:
		return &syntax.Sym{Name: x.Name}
	case nullT:
		return &syntax.Const{Kind: syntax.Null}
	case *Logical:
		if len(x.V) == 1 {
			c := &syntax.Const{Kind: syntax.Logical}
			if x.V[0] == naInt {
				c.NA = true
			} else {
				c.LogicalV = x.V[0] != 0
			}
			return c
		}
	case *Integer:
		if len(x.V) == 1 {
			c := &syntax.Const{Kind: syntax.Integer, Int: x.V[0]}
			c.NA = x.V[0] == naInt
			return c
		}
	case *Double:
		if len(x.V) == 1 {
			c := &syntax.Const{Kind: syntax.Double, Double: x.V[0]}
			c.NA = isNAReal(x.V[0])
			return c
		}
	case *String:
		if len(x.V) == 1 {
			return &syntax.Const{Kind: syntax.String, Str: x.V[0].S, NA: x.V[0].NA}
		}
	case *Promise:
		return toSyntax(x.Expr)
	case *Lang:
		if s, ok := x.Fun.(*Symbol); ok && s.Name == "function" && len(x.Args) == 2 {
			if f, ok := x.Args[0].Value.(*Formals); ok {
				fn := &syntax.Function{Body: toSyntax(x.Args[1].Value)}
				for _, p := range f.F {
					prm := syntax.Param{Name: p.Name.Name}
					if p.Default != nil {
						prm.Default = toSyntax(p.Default)
					}
					fn.Params = append(fn.Params, prm)
				}
				return fn
			}
		}
		c := &syntax.Call{Fun: toSyntax(x.Fun)}
		for _, a := range x.Args {
			c.Args = append(c.Args, syntax.Arg{Name: a.Name, HasName: a.HasName, Value: toSyntax(a.Value)})
		}
		return c
	}
	return &syntax.Sym{Name: "<object>"}
}

// ToRDS converts a value to the serialization tree, so that it can be
// compared with, or written like, what R serializes.
func ToRDS(v Value) rds.Value {
	meta := func(attr []Attr) rds.Meta {
		var m rds.Meta
		for _, a := range attr {
			m.Attrs = append(m.Attrs, rds.Attr{Name: a.Name, Value: ToRDS(a.Value)})
			if a.Name == "class" {
				m.Obj = true
			}
		}
		return m
	}
	switch x := v.(type) {
	case nullT:
		return rds.Null
	case *Symbol:
		return &rds.Sym{Name: x.Name, Gp: rds.NewStr(x.Name).Gp}
	case missingT:
		return rds.MissingArg
	case *Logical:
		return &rds.Logical{Meta: meta(x.Attr), V: append([]int32(nil), x.V...)}
	case *Integer:
		return &rds.Integer{Meta: meta(x.Attr), V: append([]int32(nil), x.V...)}
	case *Double:
		return &rds.Real{Meta: meta(x.Attr), V: append([]float64(nil), x.V...)}
	case *String:
		out := &rds.Character{Meta: meta(x.Attr), V: make([]rds.Str, len(x.V))}
		for i, s := range x.V {
			if s.NA {
				out.V[i] = rds.Str{NA: true}
			} else {
				out.V[i] = rds.NewStr(s.S)
			}
		}
		return out
	case *List:
		out := &rds.List{Meta: meta(x.Attr), V: make([]rds.Value, len(x.V))}
		for i, e := range x.V {
			out.V[i] = ToRDS(e)
		}
		return out
	case *Formals:
		if len(x.F) == 0 {
			return rds.Null
		}
		p := &rds.Pairlist{Kind: rds.LISTSXP}
		for _, f := range x.F {
			var dv rds.Value = rds.MissingArg
			if f.Default != nil {
				dv = ToRDS(f.Default)
			}
			p.Items = append(p.Items, rds.Item{Tag: &rds.Sym{Name: f.Name.Name, Gp: rds.NewStr(f.Name.Name).Gp}, Value: dv})
		}
		return p
	case *Lang:
		p := &rds.Pairlist{Kind: rds.LANGSXP}
		p.Items = append(p.Items, rds.Item{Value: ToRDS(x.Fun)})
		for _, a := range x.Args {
			it := rds.Item{Value: ToRDS(a.Value)}
			if a.HasName {
				it.Tag = &rds.Sym{Name: a.Name, Gp: rds.NewStr(a.Name).Gp}
			}
			p.Items = append(p.Items, it)
		}
		if s, ok := x.Fun.(*Symbol); ok && s.Name == "function" && len(x.Args) == 2 {
			p.Items = append(p.Items, rds.Item{Value: rds.Null}) // the absent srcref
		}
		return p
	}
	return rds.Null
}
