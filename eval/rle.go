package eval

import (
	"math"
	"strings"
)

// Run length encoding, written from rle.Rd and experiments with the R 4.6.2
// oracle. It is a Go native, not a baselib function: a blind R draft of rle
// converged on R's own 51-token run (see README, "Clean-room process").
func init() {
	regBuiltin("rle", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x")
		x := a[0]
		if x == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		if !isVector(x) || !onlyNames(x) {
			s.errorf(c.call, "'x' must be a vector of an atomic type")
		}
		n := Length(x)
		if _, ok := x.(*List); ok && n > 1 {
			s.errorf(c.call, "comparison of these types is not implemented")
		}
		nm := namesOf(x)
		var ends []int
		for i := 0; i < n; i++ {
			if i == n-1 || differs(x, i, i+1) {
				ends = append(ends, i)
			}
		}
		lens := make([]int32, len(ends))
		var lnames []Str
		prev := -1
		for k, e := range ends {
			lens[k] = int32(e - prev)
			prev = e
			if nm != nil {
				if e+1 < n && k < len(ends)-1 {
					lnames = append(lnames, nm.V[e+1])
				} else {
					lnames = append(lnames, Str{})
				}
			}
		}
		lv := &Integer{V: lens}
		if lens == nil {
			lv.V = []int32{}
		}
		if nm != nil && len(ends) > 0 {
			lv.Attr = namesAttr(lnames)
		}
		vals := s.subsetPositions(x, ends)
		if len(ends) == 0 {
			vals = withAttr(x, nil)
		}
		return &List{V: []Value{lv, vals}, Attr: []Attr{
			{Name: "names", Value: &String{V: []Str{{S: "lengths"}, {S: "values"}}}},
			{Name: "class", Value: mkString("rle")},
		}}
	})

	regBuiltin("inverse.rle", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "...")
		x := a[0]
		if x == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		l, ok := x.(*List)
		if !ok {
			if _, isNull := x.(nullT); isNull {
				s.errorf(c.call, "invalid 'rle' structure")
			}
			s.errorf(c.call, "$ operator is invalid for atomic vectors")
		}
		lens, vals := listField(l, "lengths"), listField(l, "values")
		if lens == nil || vals == nil || Length(lens) != Length(vals) {
			s.errorf(c.call, "invalid 'rle' structure")
		}
		var d []float64
		switch lens.(type) {
		case *Integer, *Double, *Logical:
			d = asDoubles(coerceTo(lens, kDouble))
		default:
			s.errorf(callOf("rep.int(v, le)"), "invalid 'times' value")
		}
		var pos []int
		total := 0
		for i, v := range d {
			if math.IsNaN(v) || v < 0 {
				s.errorf(callOf("rep.int(v, le)"), "invalid 'times' value")
			}
			if v > float64(s.opts.MaxVector) {
				s.checkLen(c.call, s.opts.MaxVector+1)
			}
			total += int(v)
			s.checkLen(c.call, total)
			for j := 0; j < int(v); j++ {
				pos = append(pos, i)
			}
		}
		out := s.subsetPositions(vals, pos)
		if at := attrsOf(out); len(at) > 0 {
			var keep []Attr
			for _, p := range at {
				if p.Name == "class" {
					keep = append(keep, p)
				}
			}
			for _, p := range at {
				if p.Name != "names" && p.Name != "class" {
					keep = append(keep, p)
				}
			}
			out = withAttr(out, keep)
		}
		return out
	})
}

// onlyNames reports whether names are the only attribute of v.
func onlyNames(v Value) bool {
	for _, a := range attrsOf(v) {
		if a.Name != "names" {
			return false
		}
	}
	return true
}

// differs reports whether elements i and j of an atomic vector are unequal,
// a missing value counting as unequal to everything.
func differs(x Value, i, j int) bool {
	switch v := x.(type) {
	case *Logical:
		return v.V[i] == naInt || v.V[j] == naInt || v.V[i] != v.V[j]
	case *Integer:
		return v.V[i] == naInt || v.V[j] == naInt || v.V[i] != v.V[j]
	case *Double:
		return math.IsNaN(v.V[i]) || math.IsNaN(v.V[j]) || v.V[i] != v.V[j]
	case *String:
		return v.V[i].NA || v.V[j].NA || v.V[i].S != v.V[j].S
	}
	return true
}

// listField is x$name: an exact match, else a unique partial one.
func listField(l *List, name string) Value {
	nm := namesOf(l)
	if nm == nil {
		return nil
	}
	hit := -1
	for i, n := range nm.V {
		if !n.NA && n.S == name {
			return l.V[i]
		}
	}
	for i, n := range nm.V {
		if !n.NA && strings.HasPrefix(n.S, name) {
			if hit >= 0 {
				return nil
			}
			hit = i
		}
	}
	if hit < 0 {
		return nil
	}
	return l.V[hit]
}
