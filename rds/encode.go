package rds

import (
	"encoding/binary"
	"math"
)

// maxPackedIndex is the largest reference index that fits in the flag word.
const maxPackedIndex = math.MaxInt32 >> 8

// Encode writes v as an XDR version-3 stream with [DefaultHeader], the bytes
// R 4.6.2's serialize(v, NULL) produces for the same object.
func Encode(v Value) ([]byte, error) { return EncodeWith(v, DefaultHeader, 0) }

// EncodeWith writes v with the given header. maxDepth bounds the nesting
// (0 selects the decoder's default). Version 2 cannot hold ALTREP objects:
// [Altrep.Materialize] them first.
func EncodeWith(v Value, h Header, maxDepth int) ([]byte, error) {
	if h.Version != 2 && h.Version != 3 {
		return nil, &Error{Offset: -1, Msg: "unsupported serialization version"}
	}
	if maxDepth <= 0 {
		maxDepth = defaultMaxDepth
	}
	e := &encoder{ver: h.Version, max: maxDepth, envs: map[Value]int{}, syms: map[string]int{}}
	e.buf = append(e.buf, 'X', '\n')
	e.int32(h.Version)
	e.int32(int32(h.Writer))
	e.int32(int32(h.MinReader))
	if h.Version == 3 {
		e.int32(int32(len(h.Encoding)))
		e.buf = append(e.buf, h.Encoding...)
	}
	if err := e.item(v, 0); err != nil {
		return nil, err
	}
	return e.buf, nil
}

type encoder struct {
	buf  []byte
	ver  int32
	max  int
	nref int
	envs map[Value]int // reference objects by identity
	syms map[string]int
}

func (e *encoder) int32(v int32) { e.buf = binary.BigEndian.AppendUint32(e.buf, uint32(v)) }

func (e *encoder) float(v float64) { e.buf = binary.BigEndian.AppendUint64(e.buf, math.Float64bits(v)) }

func (e *encoder) fail(msg string) error { return &Error{Offset: -1, Msg: msg} }

func mkFlags(t Type, m *Meta, tag bool) int32 {
	f := uint32(t) | uint32(m.Gp)<<12
	if m.Obj {
		f |= 0x100
	}
	if len(m.Attrs) > 0 {
		f |= 0x200
	}
	if tag {
		f |= 0x400
	}
	return int32(f)
}

func (e *encoder) length(n int) {
	if int64(n) <= math.MaxInt32 {
		e.int32(int32(n))
		return
	}
	e.int32(-1)
	e.int32(int32(uint32(uint64(n) >> 32)))
	e.int32(int32(uint32(uint64(n))))
}

// ref writes a back-reference to the idx-th reference object (1-based).
func (e *encoder) ref(idx int) {
	if idx <= maxPackedIndex {
		e.int32(int32(idx<<8 | int(refSXP)))
		return
	}
	e.int32(int32(refSXP))
	e.int32(int32(idx))
}

func (e *encoder) charsxp(s Str) {
	if s.NA {
		e.int32(int32(uint32(CHARSXP) | uint32(s.Gp)<<12))
		e.int32(-1)
		return
	}
	e.int32(int32(uint32(CHARSXP) | uint32(s.Gp)<<12))
	e.int32(int32(len(s.S)))
	e.buf = append(e.buf, s.S...)
}

func (e *encoder) attrs(m *Meta, depth int) error {
	if len(m.Attrs) == 0 {
		return nil
	}
	return e.attrList(m.Attrs, depth)
}

func (e *encoder) attrList(attrs []Attr, depth int) error {
	for _, a := range attrs {
		e.int32(int32(uint32(LISTSXP) | 0x400))
		s := NewStr(a.Name)
		if err := e.item(&Sym{Name: a.Name, Gp: s.Gp}, depth+1); err != nil {
			return err
		}
		if err := e.item(a.Value, depth+1); err != nil {
			return err
		}
	}
	e.int32(int32(nilValueSXP))
	return nil
}

func (e *encoder) stringVec(t Type, info []string, v Value) {
	e.nref++
	e.envs[v] = e.nref
	e.int32(int32(t))
	e.int32(0)
	e.int32(int32(len(info)))
	for _, s := range info {
		e.charsxp(NewStr(s))
	}
}

func (e *encoder) item(v Value, depth int) error {
	if depth > e.max {
		return e.fail("nesting too deep")
	}
	switch x := v.(type) {
	case nil:
		return e.fail("nil Value")
	case nullValue:
		e.int32(int32(nilValueSXP))
	case *Special:
		e.int32(int32(x.code))
	case *Sym:
		if idx, ok := e.syms[x.Name]; ok {
			e.ref(idx)
			return nil
		}
		e.nref++
		e.syms[x.Name] = e.nref
		e.int32(int32(SYMSXP))
		e.charsxp(Str{S: x.Name, Gp: x.Gp})
	case *PackageEnv:
		if idx, ok := e.envs[v]; ok {
			e.ref(idx)
			return nil
		}
		e.stringVec(packageSXP, x.Info, v)
	case *NamespaceEnv:
		if idx, ok := e.envs[v]; ok {
			e.ref(idx)
			return nil
		}
		e.stringVec(namespaceSXP, x.Info, v)
	case *Env:
		if idx, ok := e.envs[v]; ok {
			e.ref(idx)
			return nil
		}
		e.nref++
		e.envs[v] = e.nref
		e.int32(int32(ENVSXP))
		e.int32(x.Locked)
		for _, c := range []Value{x.Enclos, x.Frame, x.HashTab} {
			if err := e.item(c, depth+1); err != nil {
				return err
			}
		}
		if len(x.Attrs) == 0 {
			e.int32(int32(nilValueSXP))
			return nil
		}
		return e.attrList(x.Attrs, depth)
	case *Pairlist:
		return e.pairlist(x, depth)
	case *Closure:
		hasEnv := x.Env != Null && x.Env != nil
		e.int32(mkFlags(CLOSXP, &x.Meta, hasEnv))
		if err := e.attrs(&x.Meta, depth); err != nil {
			return err
		}
		if hasEnv {
			if err := e.item(x.Env, depth+1); err != nil {
				return err
			}
		}
		if err := e.item(x.Formals, depth+1); err != nil {
			return err
		}
		return e.item(x.Body, depth+1)
	case *Builtin:
		e.int32(int32(x.Kind))
		e.int32(int32(len(x.Name)))
		e.buf = append(e.buf, x.Name...)
	case *Logical:
		return e.ints(LGLSXP, &x.Meta, x.V, depth)
	case *Integer:
		return e.ints(INTSXP, &x.Meta, x.V, depth)
	case *Real:
		e.int32(mkFlags(REALSXP, &x.Meta, false))
		e.length(len(x.V))
		for _, f := range x.V {
			e.float(f)
		}
		return e.attrs(&x.Meta, depth)
	case *Complex:
		e.int32(mkFlags(CPLXSXP, &x.Meta, false))
		e.length(len(x.V))
		for _, c := range x.V {
			e.float(real(c))
			e.float(imag(c))
		}
		return e.attrs(&x.Meta, depth)
	case *Raw:
		e.int32(mkFlags(RAWSXP, &x.Meta, false))
		e.length(len(x.V))
		e.buf = append(e.buf, x.V...)
		return e.attrs(&x.Meta, depth)
	case *Character:
		e.int32(mkFlags(STRSXP, &x.Meta, false))
		e.length(len(x.V))
		for _, s := range x.V {
			e.charsxp(s)
		}
		return e.attrs(&x.Meta, depth)
	case *List:
		e.int32(mkFlags(x.Type(), &x.Meta, false))
		e.length(len(x.V))
		for _, c := range x.V {
			if err := e.item(c, depth+1); err != nil {
				return err
			}
		}
		return e.attrs(&x.Meta, depth)
	case *ExtPtr:
		if idx, ok := e.envs[v]; ok {
			e.ref(idx)
			return nil
		}
		e.nref++
		e.envs[v] = e.nref
		e.int32(mkFlags(EXTPTRSXP, &x.Meta, false))
		if err := e.item(x.Prot, depth+1); err != nil {
			return err
		}
		if err := e.item(x.Tag, depth+1); err != nil {
			return err
		}
		return e.attrs(&x.Meta, depth)
	case *WeakRef:
		if idx, ok := e.envs[v]; ok {
			e.ref(idx)
			return nil
		}
		e.nref++
		e.envs[v] = e.nref
		e.int32(mkFlags(WEAKREFSXP, &x.Meta, false))
		return e.attrs(&x.Meta, depth)
	case *S4:
		e.int32(mkFlags(OBJSXP, &x.Meta, false))
		return e.attrs(&x.Meta, depth)
	case *Altrep:
		if e.ver < 3 {
			return e.fail("version 2 cannot hold ALTREP objects")
		}
		e.int32(int32(uint32(altrepSXP)|uint32(x.Gp)<<12) | objBit(x.Obj))
		info := &Pairlist{Kind: LISTSXP, Items: []Item{
			{Value: x.Class}, {Value: x.Pkg}, {Value: &Integer{V: []int32{int32(x.Of)}}},
		}}
		if err := e.item(info, depth+1); err != nil {
			return err
		}
		if err := e.item(x.State, depth+1); err != nil {
			return err
		}
		if len(x.Attrs) == 0 {
			e.int32(int32(nilValueSXP))
			return nil
		}
		return e.attrList(x.Attrs, depth)
	default:
		return e.fail("unsupported Go value in the tree")
	}
	return nil
}

func objBit(o bool) int32 {
	if o {
		return 0x100
	}
	return 0
}

func (e *encoder) ints(t Type, m *Meta, v []int32, depth int) error {
	e.int32(mkFlags(t, m, false))
	e.length(len(v))
	for _, x := range v {
		e.int32(x)
	}
	return e.attrs(m, depth)
}

func (e *encoder) pairlist(p *Pairlist, depth int) error {
	if len(p.Items) == 0 {
		return e.fail("empty pairlist (use Null)")
	}
	if p.Kind != LISTSXP && p.Kind != LANGSXP && p.Kind != DOTSXP {
		return e.fail("bad pairlist kind")
	}
	for i, it := range p.Items {
		var f int32
		if i == 0 {
			f = mkFlags(p.Kind, &p.Meta, it.Tag != nil)
		} else {
			f = int32(uint32(LISTSXP) | uint32(it.Gp)<<12)
			if it.Tag != nil {
				f |= 0x400
			}
		}
		e.int32(f)
		if i == 0 {
			if err := e.attrs(&p.Meta, depth); err != nil {
				return err
			}
		}
		if it.Tag != nil {
			if err := e.item(it.Tag, depth+1); err != nil {
				return err
			}
		}
		if err := e.item(it.Value, depth+1); err != nil {
			return err
		}
	}
	if p.Tail == nil || p.Tail == Null {
		e.int32(int32(nilValueSXP))
		return nil
	}
	return e.item(p.Tail, depth+1)
}
