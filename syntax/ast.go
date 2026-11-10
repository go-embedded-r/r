package syntax

// Node is an R language object: the result of parse() is a list of Nodes.
// The tree has the shape R gives it: every operator, control-flow construct,
// subscript and brace is a call whose function is a symbol (`if`, `+`, `[`,
// `{`, `(`, ...), exactly as quote(x + 1)[[1]] is `+`. Only function
// definitions have a node of their own, because R stores their formals and
// body outside the argument list.
type Node interface {
	// Span is the source range of the node.
	Span() Span
	isNode()
}

// ConstKind is the type of a constant.
type ConstKind int

// The constant types: NULL, logical, integer, double, complex, character.
const (
	Null ConstKind = iota
	Logical
	Integer
	Double
	Complex
	String
	// Raw is never produced by the parser: it is a byte (Int) that deparse
	// writes as 0x01, inside the as.raw(...) call an evaluator raw vector
	// deparses to.
	Raw
)

// Sym is a symbol (a name, possibly backquoted in the source).
type Sym struct {
	Pos  Span
	Name string
}

// Const is a constant of length one (or NULL).
type Const struct {
	Pos  Span
	Kind ConstKind
	// Logical: 0 or 1, NA when NA is set. Integer: Int, NA when NA is set.
	// Double: Double (NA_real_ has the bits of rds.NAReal; use NA to tell).
	// Complex: Re and Im. String: Str, NA when NA is set.
	Int      int32
	Double   float64
	Im       float64
	Str      string
	NA       bool
	LogicalV bool
	// ph marks R's pipe placeholder `_`, which is a string constant until
	// the pipe rewriting replaces it.
	ph bool
}

// Arg is one argument of a call.
type Arg struct {
	// Name is the argument name; HasName tells `f(a = 1)` from `f(1)` (and
	// `f("" = 1)` from both).
	Name    string
	HasName bool
	// Value is nil for an empty argument: f(,), x[1, ].
	Value Node
	Pos   Span
}

// Call is a function call, and every operator.
type Call struct {
	Pos  Span
	Fun  Node
	Args []Arg
	// Elems marks the call an evaluator expression vector deparses to,
	// expression(...): its arguments are elements, so an `=` call among them
	// is not parenthesized as an argument would be.
	Elems bool
	h     int // height of the subtree, checked against Limits.MaxDepth
}

// Param is one formal argument of a function definition.
type Param struct {
	Name string
	// Default is nil when there is none.
	Default Node
	Pos     Span
}

// Function is `function(params) body` or `\(params) body`.
type Function struct {
	Pos    Span
	Params []Param
	Body   Node
	// Lambda is set for the backslash form; both are the same call to
	// `function` in R.
	Lambda bool
	h      int
}

func (n *Sym) Span() Span      { return n.Pos }
func (n *Const) Span() Span    { return n.Pos }
func (n *Call) Span() Span     { return n.Pos }
func (n *Function) Span() Span { return n.Pos }
func (*Sym) isNode()           {}
func (*Const) isNode()         {}
func (*Call) isNode()          {}
func (*Function) isNode()      {}

// FunName returns the name of the symbol a call's function is, or "".
func (c *Call) FunName() string {
	if s, ok := c.Fun.(*Sym); ok {
		return s.Name
	}
	return ""
}
