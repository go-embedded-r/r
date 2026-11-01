package rds

import (
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

// Limits bound what the decoder accepts. The zero value of a field selects
// its default.
type Limits struct {
	// MaxBytes bounds the input after decompression (default 256 MiB). A
	// compressed input that expands beyond it is refused, not truncated.
	MaxBytes int
	// MaxDepth bounds the nesting of the object tree (default 2000). The
	// cells of a pairlist do not count: they are read iteratively.
	MaxDepth int
}

const (
	defaultMaxBytes = 256 << 20
	defaultMaxDepth = 2000
)

func (l Limits) withDefaults() Limits {
	if l.MaxBytes <= 0 {
		l.MaxBytes = defaultMaxBytes
	}
	if l.MaxDepth <= 0 {
		l.MaxDepth = defaultMaxDepth
	}
	return l
}

// ErrTooLarge is wrapped by the error returned when the input, after
// decompression, exceeds [Limits.MaxBytes].
var ErrTooLarge = errors.New("rds: input exceeds the size limit")

// Header is the preamble of a serialized stream.
type Header struct {
	Version   int32  // 2 or 3
	Writer    uint32 // R version that wrote it, packed as 65536*major+256*minor+patch
	MinReader uint32 // the oldest R version able to read it
	Encoding  string // native encoding at write time (version 3 only)
}

// DefaultHeader is what [Encode] writes: version 3 as written by R 4.6.2.
var DefaultHeader = Header{Version: 3, Writer: 0x00040602, MinReader: 0x00030500, Encoding: "UTF-8"}

// Decode reads one serialized object from b, which may be gzip- or
// bzip2-compressed (detected by magic number, as readRDS does for files).
// The XDR binary format, versions 2 and 3, is supported; the ASCII and
// native-binary variants, xz compression, byte-code and promises are refused
// with an *[Error].
func Decode(b []byte, lim Limits) (Value, Header, error) {
	lim = lim.withDefaults()
	b, err := decompress(b, lim.MaxBytes)
	if err != nil {
		return nil, Header{}, err
	}
	d := &decoder{b: b, lim: lim, syms: map[string]*Sym{}}
	h, err := d.header()
	if err != nil {
		return nil, Header{}, err
	}
	v, err := d.item(0)
	if err != nil {
		return nil, Header{}, err
	}
	if d.pos != len(d.b) {
		return nil, Header{}, d.fail("%d bytes of trailing data", len(d.b)-d.pos)
	}
	return v, h, nil
}

// DecodeReader is Decode over a reader, reading at most Limits.MaxBytes
// (plus one byte to detect an overflow) from it.
func DecodeReader(r io.Reader, lim Limits) (Value, Header, error) {
	lim = lim.withDefaults()
	b, err := io.ReadAll(io.LimitReader(r, int64(lim.MaxBytes)+1))
	if err != nil {
		return nil, Header{}, err
	}
	return Decode(b, lim)
}

func decompress(b []byte, max int) ([]byte, error) {
	var r io.Reader
	switch {
	case len(b) >= 2 && b[0] == 0x1f && b[1] == 0x8b:
		z, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return nil, &Error{Offset: -1, Msg: "gzip: " + err.Error()}
		}
		r = z
	case len(b) >= 3 && string(b[:3]) == "BZh":
		r = bzip2.NewReader(bytes.NewReader(b))
	case len(b) >= 6 && string(b[:6]) == "\xfd7zXZ\x00":
		return nil, &Error{Offset: -1, Msg: "xz compression is not supported"}
	default:
		if len(b) > max {
			return nil, tooLarge()
		}
		return b, nil
	}
	out, err := io.ReadAll(io.LimitReader(r, int64(max)+1))
	if err != nil {
		return nil, &Error{Offset: -1, Msg: "decompression: " + err.Error()}
	}
	if len(out) > max {
		return nil, tooLarge()
	}
	return out, nil
}

type sizeError struct{ e *Error }

func (s sizeError) Error() string { return s.e.Error() }
func (sizeError) Unwrap() error   { return ErrTooLarge }

func tooLarge() error {
	return sizeError{&Error{Offset: -1, Msg: "input exceeds the size limit"}}
}

type decoder struct {
	b    []byte
	pos  int
	lim  Limits
	refs []Value
	syms map[string]*Sym
}

func (d *decoder) fail(format string, a ...any) error {
	return &Error{Offset: d.pos, Msg: fmt.Sprintf(format, a...)}
}

func (d *decoder) remaining() int { return len(d.b) - d.pos }

func (d *decoder) int32() (int32, error) {
	if d.remaining() < 4 {
		return 0, d.fail("unexpected end of input")
	}
	v := int32(binary.BigEndian.Uint32(d.b[d.pos:]))
	d.pos += 4
	return v, nil
}

func (d *decoder) header() (Header, error) {
	var h Header
	if d.remaining() < 2 {
		return h, d.fail("unexpected end of input")
	}
	switch string(d.b[:2]) {
	case "X\n":
	case "A\n", "B\n":
		return h, d.fail("the %q serialization format is not supported (only XDR)", d.b[0])
	default:
		return h, d.fail("not a serialized R object (no XDR header)")
	}
	d.pos = 2
	v, err := d.int32()
	if err != nil {
		return h, err
	}
	if v != 2 && v != 3 {
		return h, d.fail("unsupported serialization version %d", v)
	}
	h.Version = v
	w, err := d.int32()
	if err != nil {
		return h, err
	}
	m, err := d.int32()
	if err != nil {
		return h, err
	}
	h.Writer, h.MinReader = uint32(w), uint32(m)
	if v == 3 {
		n, err := d.int32()
		if err != nil {
			return h, err
		}
		if n < 0 || int(n) > d.remaining() {
			return h, d.fail("bad encoding-name length %d", n)
		}
		h.Encoding = string(d.b[d.pos : d.pos+int(n)])
		d.pos += int(n)
	}
	return h, nil
}

// flags is the first word of an item.
type flags uint32

func (f flags) typ() Type     { return Type(f & 0xff) }
func (f flags) obj() bool     { return f&0x100 != 0 }
func (f flags) hasAttr() bool { return f&0x200 != 0 }
func (f flags) hasTag() bool  { return f&0x400 != 0 }
func (f flags) gp() uint16    { return uint16(f >> 12) }

func (d *decoder) item(depth int) (Value, error) {
	if depth > d.lim.MaxDepth {
		return nil, d.fail("nesting deeper than %d", d.lim.MaxDepth)
	}
	w, err := d.int32()
	if err != nil {
		return nil, err
	}
	return d.dispatch(flags(w), depth)
}

func (d *decoder) dispatch(f flags, depth int) (Value, error) {
	t := f.typ()
	switch t {
	case nilValueSXP:
		return Null, nil
	case emptyEnvSXP:
		return EmptyEnv, nil
	case baseEnvSXP:
		return BaseEnv, nil
	case globalEnvSXP:
		return GlobalEnv, nil
	case unboundValueSXP:
		return UnboundValue, nil
	case missingArgSXP:
		return MissingArg, nil
	case baseNamespaceSXP:
		return BaseNamespace, nil
	case refSXP:
		return d.ref(f)
	case persistSXP:
		return nil, d.fail("persistent references are not supported")
	case packageSXP, namespaceSXP:
		info, err := d.stringVec()
		if err != nil {
			return nil, err
		}
		var v Value
		if t == packageSXP {
			v = &PackageEnv{Info: info}
		} else {
			v = &NamespaceEnv{Info: info}
		}
		d.refs = append(d.refs, v)
		return v, nil
	case SYMSXP:
		s, err := d.charsxp()
		if err != nil {
			return nil, err
		}
		sym, ok := d.syms[s.S]
		if !ok {
			sym = &Sym{Name: s.S, Gp: s.Gp}
			d.syms[s.S] = sym
		}
		d.refs = append(d.refs, sym)
		return sym, nil
	case ENVSXP:
		return d.env(depth)
	case LISTSXP, LANGSXP, DOTSXP:
		return d.pairlist(f, depth)
	case CLOSXP:
		return d.closure(f, depth)
	case SPECIALSXP, BUILTINSXP:
		n, err := d.int32()
		if err != nil {
			return nil, err
		}
		if n < 0 || int(n) > d.remaining() {
			return nil, d.fail("bad builtin-name length %d", n)
		}
		name := string(d.b[d.pos : d.pos+int(n)])
		d.pos += int(n)
		return &Builtin{Kind: t, Name: name}, nil
	case LGLSXP, INTSXP:
		n, err := d.length(4)
		if err != nil {
			return nil, err
		}
		v := make([]int32, n)
		for i := range v {
			v[i] = int32(binary.BigEndian.Uint32(d.b[d.pos+4*i:]))
		}
		d.pos += 4 * n
		m, err := d.meta(f, depth)
		if err != nil {
			return nil, err
		}
		if t == LGLSXP {
			return &Logical{Meta: m, V: v}, nil
		}
		return &Integer{Meta: m, V: v}, nil
	case REALSXP:
		n, err := d.length(8)
		if err != nil {
			return nil, err
		}
		v := make([]float64, n)
		for i := range v {
			v[i] = math.Float64frombits(binary.BigEndian.Uint64(d.b[d.pos+8*i:]))
		}
		d.pos += 8 * n
		m, err := d.meta(f, depth)
		if err != nil {
			return nil, err
		}
		return &Real{Meta: m, V: v}, nil
	case CPLXSXP:
		n, err := d.length(16)
		if err != nil {
			return nil, err
		}
		v := make([]complex128, n)
		for i := range v {
			re := math.Float64frombits(binary.BigEndian.Uint64(d.b[d.pos+16*i:]))
			im := math.Float64frombits(binary.BigEndian.Uint64(d.b[d.pos+16*i+8:]))
			v[i] = complex(re, im)
		}
		d.pos += 16 * n
		m, err := d.meta(f, depth)
		if err != nil {
			return nil, err
		}
		return &Complex{Meta: m, V: v}, nil
	case RAWSXP:
		n, err := d.length(1)
		if err != nil {
			return nil, err
		}
		v := append([]byte(nil), d.b[d.pos:d.pos+n]...)
		d.pos += n
		m, err := d.meta(f, depth)
		if err != nil {
			return nil, err
		}
		return &Raw{Meta: m, V: v}, nil
	case STRSXP:
		// every element is at least a flag word and a length word
		n, err := d.length(8)
		if err != nil {
			return nil, err
		}
		v := make([]Str, n)
		for i := range v {
			s, err := d.charsxp()
			if err != nil {
				return nil, err
			}
			v[i] = s
		}
		m, err := d.meta(f, depth)
		if err != nil {
			return nil, err
		}
		return &Character{Meta: m, V: v}, nil
	case VECSXP, EXPRSXP:
		n, err := d.length(4)
		if err != nil {
			return nil, err
		}
		v := make([]Value, n)
		for i := range v {
			if v[i], err = d.item(depth + 1); err != nil {
				return nil, err
			}
		}
		m, err := d.meta(f, depth)
		if err != nil {
			return nil, err
		}
		return &List{Meta: m, Expr: t == EXPRSXP, V: v}, nil
	case EXTPTRSXP:
		p := &ExtPtr{}
		d.refs = append(d.refs, p)
		var err error
		if p.Prot, err = d.item(depth + 1); err != nil {
			return nil, err
		}
		if p.Tag, err = d.item(depth + 1); err != nil {
			return nil, err
		}
		if p.Meta, err = d.meta(f, depth); err != nil {
			return nil, err
		}
		return p, nil
	case WEAKREFSXP:
		w := &WeakRef{}
		d.refs = append(d.refs, w)
		var err error
		if w.Meta, err = d.meta(f, depth); err != nil {
			return nil, err
		}
		return w, nil
	case OBJSXP:
		m, err := d.meta(f, depth)
		if err != nil {
			return nil, err
		}
		return &S4{Meta: m}, nil
	case altrepSXP:
		return d.altrep(f, depth)
	case BCODESXP:
		return nil, d.fail("byte-code is not supported")
	case PROMSXP:
		return nil, d.fail("promises are not supported")
	}
	return nil, d.fail("unknown type code %d", t)
}

// meta reads the attributes that follow the data of a non-pairlist item.
func (d *decoder) meta(f flags, depth int) (Meta, error) {
	m := Meta{Obj: f.obj(), Gp: f.gp()}
	if f.hasAttr() {
		a, err := d.attrs(depth)
		if err != nil {
			return m, err
		}
		m.Attrs = a
	}
	return m, nil
}

func (d *decoder) attrs(depth int) ([]Attr, error) {
	v, err := d.item(depth + 1)
	if err != nil {
		return nil, err
	}
	return d.attrsOf(v)
}

func (d *decoder) attrsOf(v Value) ([]Attr, error) {
	if v == Null {
		return nil, nil
	}
	p, ok := v.(*Pairlist)
	if !ok || p.Kind != LISTSXP || p.Tail != nil || len(p.Meta.Attrs) != 0 {
		return nil, d.fail("attributes are not a plain tagged pairlist")
	}
	out := make([]Attr, len(p.Items))
	for i, it := range p.Items {
		if it.Tag == nil {
			return nil, d.fail("untagged attribute")
		}
		out[i] = Attr{Name: it.Tag.Name, Value: it.Value}
	}
	return out, nil
}

// length reads a vector length and checks that elem*length bytes remain.
func (d *decoder) length(elem int) (int, error) {
	n, err := d.int32()
	if err != nil {
		return 0, err
	}
	var n64 int64
	switch {
	case n >= 0:
		n64 = int64(n)
	case n == -1:
		hi, err := d.int32()
		if err != nil {
			return 0, err
		}
		lo, err := d.int32()
		if err != nil {
			return 0, err
		}
		n64 = int64(uint32(hi))<<32 | int64(uint32(lo))
		if n64 < 0 {
			return 0, d.fail("vector length out of range")
		}
	default:
		return 0, d.fail("negative vector length %d", n)
	}
	if n64 > int64(d.remaining()/elem) {
		return 0, d.fail("vector of %d elements does not fit in the %d bytes left", n64, d.remaining())
	}
	return int(n64), nil
}

// charsxp reads a CHARSXP item (flag word, length, bytes).
func (d *decoder) charsxp() (Str, error) {
	w, err := d.int32()
	if err != nil {
		return Str{}, err
	}
	f := flags(w)
	if f.typ() != CHARSXP {
		return Str{}, d.fail("expected a CHARSXP, found type %d", f.typ())
	}
	n, err := d.int32()
	if err != nil {
		return Str{}, err
	}
	if n == -1 {
		return Str{NA: true, Gp: f.gp()}, nil
	}
	if n < 0 || int(n) > d.remaining() {
		return Str{}, d.fail("bad string length %d", n)
	}
	s := string(d.b[d.pos : d.pos+int(n)])
	d.pos += int(n)
	return Str{S: s, Gp: f.gp()}, nil
}

// stringVec reads the body of a namespace, package or persistent reference:
// a zero, a count, then that many CHARSXPs.
func (d *decoder) stringVec() ([]string, error) {
	z, err := d.int32()
	if err != nil {
		return nil, err
	}
	if z != 0 {
		return nil, d.fail("unexpected %d where a zero is required", z)
	}
	n, err := d.length(8)
	if err != nil {
		return nil, err
	}
	out := make([]string, n)
	for i := range out {
		s, err := d.charsxp()
		if err != nil {
			return nil, err
		}
		out[i] = s.S
	}
	return out, nil
}

func (d *decoder) ref(f flags) (Value, error) {
	idx := int64(uint32(f) >> 8)
	if idx == 0 {
		n, err := d.int32()
		if err != nil {
			return nil, err
		}
		idx = int64(n)
	}
	if idx < 1 || idx > int64(len(d.refs)) {
		return nil, d.fail("reference %d out of range (%d defined)", idx, len(d.refs))
	}
	return d.refs[idx-1], nil
}

func (d *decoder) env(depth int) (Value, error) {
	locked, err := d.int32()
	if err != nil {
		return nil, err
	}
	e := &Env{Locked: locked}
	d.refs = append(d.refs, e)
	if e.Enclos, err = d.item(depth + 1); err != nil {
		return nil, err
	}
	if e.Frame, err = d.item(depth + 1); err != nil {
		return nil, err
	}
	if e.HashTab, err = d.item(depth + 1); err != nil {
		return nil, err
	}
	a, err := d.attrs(depth)
	if err != nil {
		return nil, err
	}
	e.Attrs = a
	return e, nil
}

func (d *decoder) pairlist(f flags, depth int) (Value, error) {
	p := &Pairlist{Kind: f.typ(), Meta: Meta{Obj: f.obj(), Gp: f.gp()}}
	if f.hasAttr() {
		a, err := d.attrs(depth)
		if err != nil {
			return nil, err
		}
		p.Attrs = a
	}
	cur := f
	for {
		it := Item{Gp: cur.gp()}
		if cur.hasTag() {
			t, err := d.item(depth + 1)
			if err != nil {
				return nil, err
			}
			s, ok := t.(*Sym)
			if !ok {
				return nil, d.fail("pairlist tag is not a symbol")
			}
			it.Tag = s
		}
		var err error
		if it.Value, err = d.item(depth + 1); err != nil {
			return nil, err
		}
		p.Items = append(p.Items, it)
		w, err := d.int32()
		if err != nil {
			return nil, err
		}
		next := flags(w)
		switch next.typ() {
		case nilValueSXP:
			return p, nil
		case LISTSXP:
			if next.hasAttr() || next.obj() {
				return nil, d.fail("attributes on an interior pairlist cell are not supported")
			}
			cur = next
		default:
			tail, err := d.dispatch(next, depth+1)
			if err != nil {
				return nil, err
			}
			p.Tail = tail
			return p, nil
		}
	}
}

func (d *decoder) closure(f flags, depth int) (Value, error) {
	c := &Closure{Meta: Meta{Obj: f.obj(), Gp: f.gp()}, Env: Null}
	var err error
	if f.hasAttr() {
		if c.Attrs, err = d.attrs(depth); err != nil {
			return nil, err
		}
	}
	if f.hasTag() {
		if c.Env, err = d.item(depth + 1); err != nil {
			return nil, err
		}
	}
	if c.Formals, err = d.item(depth + 1); err != nil {
		return nil, err
	}
	if c.Body, err = d.item(depth + 1); err != nil {
		return nil, err
	}
	return c, nil
}

func (d *decoder) altrep(f flags, depth int) (Value, error) {
	info, err := d.item(depth + 1)
	if err != nil {
		return nil, err
	}
	p, ok := info.(*Pairlist)
	if !ok || p.Kind != LISTSXP || len(p.Items) != 3 || p.Tail != nil {
		return nil, d.fail("malformed ALTREP class info")
	}
	cls, ok1 := p.Items[0].Value.(*Sym)
	pkg, ok2 := p.Items[1].Value.(*Sym)
	of, ok3 := p.Items[2].Value.(*Integer)
	if !ok1 || !ok2 || !ok3 || len(of.V) != 1 || of.V[0] < 0 || of.V[0] > 255 {
		return nil, d.fail("malformed ALTREP class info")
	}
	state, err := d.item(depth + 1)
	if err != nil {
		return nil, err
	}
	a := &Altrep{Meta: Meta{Obj: f.obj(), Gp: f.gp()}, Class: cls, Pkg: pkg, Of: Type(of.V[0]), State: state}
	if a.Attrs, err = d.attrs(depth); err != nil {
		return nil, err
	}
	return a, nil
}
