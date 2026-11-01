package rds

import (
	"errors"
	"math"
	"strconv"
)

// ErrUnsupportedAltrep is wrapped by the error [Altrep.Materialize] returns
// for a class it cannot expand.
var ErrUnsupportedAltrep = errors.New("rds: unsupported ALTREP class")

// Materialize expands an ALTREP object into the ordinary vector it stands
// for. The classes of R's base package that serialize() writes for ordinary
// data are handled: compact_intseq and compact_realseq (the result of `a:b`
// and seq_len), the wrap_* classes (a vector with cached metadata), and
// deferred_string over an integer vector. maxLen bounds the length of the
// vector produced (a compact sequence is a few bytes however long it is).
func (a *Altrep) Materialize(maxLen int) (Value, error) {
	switch a.Class.Name {
	case "compact_intseq", "compact_realseq":
		st, ok := a.State.(*Real)
		if !ok || len(st.V) != 3 {
			return nil, errors.New("rds: malformed compact sequence state")
		}
		n, start, step := st.V[0], st.V[1], st.V[2]
		if !(n >= 0) || n > float64(maxLen) {
			return nil, errors.New("rds: compact sequence length out of range")
		}
		count := int(n)
		if a.Class.Name == "compact_realseq" {
			out := make([]float64, count)
			for i := range out {
				out[i] = start + float64(i)*step
			}
			return &Real{Meta: a.Meta, V: out}, nil
		}
		out := make([]int32, count)
		last := start + float64(count-1)*step
		if count > 0 && (start < math.MinInt32+1 || start > math.MaxInt32 || last < math.MinInt32+1 || last > math.MaxInt32) {
			return nil, errors.New("rds: compact integer sequence overflows")
		}
		for i := range out {
			out[i] = int32(start + float64(i)*step)
		}
		return &Integer{Meta: a.Meta, V: out}, nil
	case "wrap_real", "wrap_integer", "wrap_logical", "wrap_string", "wrap_complex", "wrap_raw", "wrap_list":
		st, ok := a.State.(*List)
		if !ok || len(st.V) != 2 {
			return nil, errors.New("rds: malformed wrapper state")
		}
		return withMeta(st.V[0], a.Meta)
	case "deferred_string":
		st, ok := a.State.(*Pairlist)
		if !ok || len(st.Items) != 1 {
			return nil, errors.New("rds: malformed deferred_string state")
		}
		arg := st.Items[0].Value
		if inner, ok := arg.(*Altrep); ok {
			v, err := inner.Materialize(maxLen)
			if err != nil {
				return nil, err
			}
			arg = v
		}
		ints, ok := arg.(*Integer)
		if !ok {
			return nil, ErrUnsupportedAltrep
		}
		if len(ints.V) > maxLen {
			return nil, errors.New("rds: vector longer than the limit")
		}
		out := make([]Str, len(ints.V))
		for i, x := range ints.V {
			if x == NAInteger {
				out[i] = Str{NA: true}
			} else {
				out[i] = Str{S: strconv.Itoa(int(x)), Gp: 64}
			}
		}
		return &Character{Meta: a.Meta, V: out}, nil
	}
	return nil, ErrUnsupportedAltrep
}

// withMeta returns a shallow copy of the vector v whose attributes are m's.
func withMeta(v Value, m Meta) (Value, error) {
	switch x := v.(type) {
	case *Logical:
		return &Logical{Meta: m, V: x.V}, nil
	case *Integer:
		return &Integer{Meta: m, V: x.V}, nil
	case *Real:
		return &Real{Meta: m, V: x.V}, nil
	case *Complex:
		return &Complex{Meta: m, V: x.V}, nil
	case *Raw:
		return &Raw{Meta: m, V: x.V}, nil
	case *Character:
		return &Character{Meta: m, V: x.V}, nil
	case *List:
		return &List{Meta: m, Expr: x.Expr, V: x.V}, nil
	}
	return nil, errors.New("rds: wrapper around a non-vector")
}
