// Package rds reads and writes R's serialization format (version 2 and 3,
// XDR byte order): the format of serialize(), saveRDS() and readRDS().
//
// The format is described in "R Internals", section "Serialization
// Formats". The numeric codes of the pseudo-types it lists but does not
// number (references, the special environments, ALTREP) were observed on the
// reference R 4.6.2 (see internal/oracle) and are checked against its output
// by the golden tests. Nothing here is translated from R's sources.
//
// A decoded document is a tree of [Value]s that keeps everything needed to
// write the same bytes again: attributes in order, the object bit, the
// general-purpose bits, the encoding flags of every string, NaN payloads,
// reference sharing of environments and symbols, and ALTREP objects in their
// serialized (not expanded) form; [Altrep.Materialize] expands them.
//
// The decoder is built for untrusted input: it refuses inputs larger than
// [Limits.MaxBytes] (after decompression), nesting deeper than
// [Limits.MaxDepth], and any length that the remaining bytes cannot hold,
// before allocating.
package rds

import "fmt"

// Type is an R SEXPTYPE (or, above 200, a serialization pseudo-type).
type Type uint8

// SEXPTYPE codes (R Internals, "SEXPs"; R Language Definition, typeof()).
const (
	NILSXP     Type = 0
	SYMSXP     Type = 1
	LISTSXP    Type = 2
	CLOSXP     Type = 3
	ENVSXP     Type = 4
	PROMSXP    Type = 5
	LANGSXP    Type = 6
	SPECIALSXP Type = 7
	BUILTINSXP Type = 8
	CHARSXP    Type = 9
	LGLSXP     Type = 10
	INTSXP     Type = 13
	REALSXP    Type = 14
	CPLXSXP    Type = 15
	STRSXP     Type = 16
	DOTSXP     Type = 17
	VECSXP     Type = 19
	EXPRSXP    Type = 20
	BCODESXP   Type = 21
	EXTPTRSXP  Type = 22
	WEAKREFSXP Type = 23
	RAWSXP     Type = 24
	OBJSXP     Type = 25 // S4SXP before R 4.4
)

// Serialization pseudo-types, observed on R 4.6.2.
const (
	altrepSXP        Type = 238
	attrListSXP      Type = 239
	attrLangSXP      Type = 240
	baseEnvSXP       Type = 241
	emptyEnvSXP      Type = 242
	persistSXP       Type = 247
	packageSXP       Type = 248
	namespaceSXP     Type = 249
	baseNamespaceSXP Type = 250
	missingArgSXP    Type = 251
	unboundValueSXP  Type = 252
	globalEnvSXP     Type = 253
	nilValueSXP      Type = 254
	refSXP           Type = 255
)

// NA values of the atomic types, as R stores them.
const (
	// NAInteger is NA_integer_ and NA_logical_: INT_MIN.
	NAInteger int32 = -1 << 31
	// NAReal is the bit pattern of NA_real_: a NaN whose low word is 1954.
	NAReal uint64 = 0x7FF00000000007A2
)

// Value is a node of a decoded document.
type Value interface {
	// Type is the SEXPTYPE the node is written as.
	Type() Type
}

// Attr is one attribute, in the order R stores them.
type Attr struct {
	Name  string
	Value Value
}

// Meta is what every non-trivial node carries besides its data.
type Meta struct {
	// Attrs are the attributes in storage order.
	Attrs []Attr
	// Obj is the object bit (set when the node has a class attribute, and
	// for S4 instances).
	Obj bool
	// Gp is the 16-bit general-purpose field ("levels" in the flag word).
	Gp uint16
}

// Attr returns the value of the named attribute, or nil.
func (m *Meta) Attr(name string) Value {
	for _, a := range m.Attrs {
		if a.Name == name {
			return a.Value
		}
	}
	return nil
}

func (m *Meta) meta() *Meta { return m }

// WithMeta is implemented by every node that carries a [Meta].
type WithMeta interface {
	Value
	meta() *Meta
}

// Null is R's NULL.
var Null Value = nullValue{}

type nullValue struct{}

func (nullValue) Type() Type { return NILSXP }

// Special is one of the singleton objects the format writes as a bare
// pseudo-type: the global, empty and base environments, the base namespace,
// the missing-argument marker and the unbound-value marker.
type Special struct {
	name string
	code Type
}

// Type implements [Value]; the environments report ENVSXP, the two markers
// SYMSXP, as typeof() does.
func (s *Special) Type() Type {
	switch s.code {
	case missingArgSXP, unboundValueSXP:
		return SYMSXP
	}
	return ENVSXP
}

func (s *Special) String() string { return s.name }

// The singletons.
var (
	GlobalEnv     = &Special{"R_GlobalEnv", globalEnvSXP}
	EmptyEnv      = &Special{"R_EmptyEnv", emptyEnvSXP}
	BaseEnv       = &Special{"R_BaseEnv", baseEnvSXP}
	BaseNamespace = &Special{"R_BaseNamespace", baseNamespaceSXP}
	MissingArg    = &Special{"R_MissingArg", missingArgSXP}
	UnboundValue  = &Special{"R_UnboundValue", unboundValueSXP}
)

// Sym is a symbol (SYMSXP). Within one document, equal names share one *Sym.
type Sym struct {
	Name string
	// Gp is the flag-word levels of the symbol's print name (its encoding
	// bits), kept so that the bytes round-trip.
	Gp uint16
}

// Type implements [Value].
func (*Sym) Type() Type { return SYMSXP }

// Item is one cell of a pairlist, call or dots list.
type Item struct {
	Tag   *Sym // nil when untagged
	Value Value
	Gp    uint16
}

// Pairlist is a LISTSXP, LANGSXP or DOTSXP: a chain of cells. Meta belongs to
// the first cell. Tail is the final cdr when it is not NULL (an improper
// list, which the format allows).
type Pairlist struct {
	Meta
	Kind  Type
	Items []Item
	Tail  Value // nil means NULL
}

// Type implements [Value].
func (p *Pairlist) Type() Type { return p.Kind }

// Closure is a CLOSXP.
type Closure struct {
	Meta
	Env     Value // the enclosing environment
	Formals Value // a *Pairlist or Null
	Body    Value
}

// Type implements [Value].
func (*Closure) Type() Type { return CLOSXP }

// Env is an ordinary environment (ENVSXP). Frame and HashTab are the two
// representations R uses; exactly one is non-NULL for a non-empty
// environment.
type Env struct {
	Meta
	Locked  int32
	Enclos  Value
	Frame   Value
	HashTab Value
}

// Type implements [Value].
func (*Env) Type() Type { return ENVSXP }

// NamespaceEnv is a namespace environment, written by name and version.
type NamespaceEnv struct{ Info []string }

// Type implements [Value].
func (*NamespaceEnv) Type() Type { return ENVSXP }

// PackageEnv is a package environment ("package:stats"), written by name.
type PackageEnv struct{ Info []string }

// Type implements [Value].
func (*PackageEnv) Type() Type { return ENVSXP }

// Builtin is a SPECIALSXP or BUILTINSXP, written by name.
type Builtin struct {
	Kind Type
	Name string
}

// Type implements [Value].
func (b *Builtin) Type() Type { return b.Kind }

// ExtPtr is an external pointer: only its protected value and tag survive.
type ExtPtr struct {
	Meta
	Prot, Tag Value
}

// Type implements [Value].
func (*ExtPtr) Type() Type { return EXTPTRSXP }

// WeakRef is a weak reference (nothing but attributes is written).
type WeakRef struct{ Meta }

// Type implements [Value].
func (*WeakRef) Type() Type { return WEAKREFSXP }

// S4 is an S4 instance without a vector payload (OBJSXP).
type S4 struct{ Meta }

// Type implements [Value].
func (*S4) Type() Type { return OBJSXP }

// Logical is a LGLSXP; TRUE is 1, FALSE 0, NA [NAInteger].
type Logical struct {
	Meta
	V []int32
}

// Type implements [Value].
func (*Logical) Type() Type { return LGLSXP }

// Integer is an INTSXP; NA is [NAInteger].
type Integer struct {
	Meta
	V []int32
}

// Type implements [Value].
func (*Integer) Type() Type { return INTSXP }

// Real is a REALSXP. NA_real_ is the NaN with bits [NAReal]; payloads are
// preserved.
type Real struct {
	Meta
	V []float64
}

// Type implements [Value].
func (*Real) Type() Type { return REALSXP }

// Complex is a CPLXSXP.
type Complex struct {
	Meta
	V []complex128
}

// Type implements [Value].
func (*Complex) Type() Type { return CPLXSXP }

// Raw is a RAWSXP.
type Raw struct {
	Meta
	V []byte
}

// Type implements [Value].
func (*Raw) Type() Type { return RAWSXP }

// Str is one element of a character vector (a CHARSXP).
type Str struct {
	S  string
	NA bool // NA_character_; S is empty
	// Gp is the CHARSXP levels: 64 ASCII, 8 UTF-8, 4 Latin-1, 2 bytes.
	Gp uint16
}

// Character is a STRSXP.
type Character struct {
	Meta
	V []Str
}

// Type implements [Value].
func (*Character) Type() Type { return STRSXP }

// List is a VECSXP, or an EXPRSXP when Expr is set.
type List struct {
	Meta
	Expr bool
	V    []Value
}

// Type implements [Value].
func (l *List) Type() Type {
	if l.Expr {
		return EXPRSXP
	}
	return VECSXP
}

// Altrep is an ALTREP object in its serialized form: the class (and the
// package defining it), the vector type it stands for, and its state.
type Altrep struct {
	Meta
	Class, Pkg *Sym
	Of         Type
	State      Value
}

// Type implements [Value]: the type of the vector it stands for.
func (a *Altrep) Type() Type { return a.Of }

// Strings is a convenience constructor for a character vector.
func Strings(s ...string) *Character {
	c := &Character{V: make([]Str, len(s))}
	for i, x := range s {
		c.V[i] = NewStr(x)
	}
	return c
}

// NewStr builds a non-NA string with the flag R would set for it in a UTF-8
// session: 64 for pure ASCII, 8 otherwise.
func NewStr(s string) Str {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return Str{S: s, Gp: 8}
		}
	}
	return Str{S: s, Gp: 64}
}

// Error is the type of every failure of the decoder and encoder.
type Error struct {
	Offset int // byte offset in the (decompressed) input, -1 when not applicable
	Msg    string
}

func (e *Error) Error() string {
	if e.Offset < 0 {
		return "rds: " + e.Msg
	}
	return fmt.Sprintf("rds: offset %d: %s", e.Offset, e.Msg)
}
