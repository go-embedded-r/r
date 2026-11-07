package eval

import (
	"math"
	"sort"
)

func init() {
	regBuiltin("sort", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "decreasing", "na.last")
		x := a[0]
		if x == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		dec := a[1] != nil && s.truthy(c, a[1], "decreasing")
		switch x.(type) {
		case *Logical, *Integer, *Double, *String:
		case nullT:
			return Null
		default:
			s.errorf(c.call, "'x' must be atomic")
		}
		n := Length(x)
		var idx []int
		isNA := func(i int) bool {
			switch v := x.(type) {
			case *Logical:
				return v.V[i] == naInt
			case *Integer:
				return v.V[i] == naInt
			case *Double:
				return math.IsNaN(v.V[i])
			case *String:
				return v.V[i].NA
			}
			return false
		}
		for i := 0; i < n; i++ {
			if !isNA(i) {
				idx = append(idx, i)
			}
		}
		less := func(i, j int) bool {
			switch v := x.(type) {
			case *Logical:
				return v.V[i] < v.V[j]
			case *Integer:
				return v.V[i] < v.V[j]
			case *Double:
				return v.V[i] < v.V[j]
			case *String:
				return collate(v.V[i].S, v.V[j].S) < 0
			}
			return false
		}
		sort.SliceStable(idx, func(p, q int) bool {
			if dec {
				return less(idx[q], idx[p])
			}
			return less(idx[p], idx[q])
		})
		if idx == nil {
			idx = []int{}
		}
		return s.subsetPositions(x, idx)
	})
}
