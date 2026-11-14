package eval

import (
	"fmt"
	"math"
	"sort"
)

// cut.default, written from cut.Rd and experiments with the R 4.6.2 oracle.
// A Go native: a draft in R converged on R's own text.

func cFormatG(v float64, dig int) string {
	switch {
	case math.IsInf(v, 1):
		return " Inf" // formatC pads an infinite value to the width of -Inf
	case math.IsInf(v, -1):
		return "-Inf"
	case math.IsNaN(v):
		return "NaN"
	}
	return fmt.Sprintf("%.*g", dig, v+0)
}

func init() {
	regBuiltin("cut.default", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "breaks", "labels", "include.lowest", "right", "dig.lab", "ordered_result", "...")
		if a[0] == nil || a[1] == nil {
			s.errorf(c.call, "argument \"%s\" is missing, with no default", map[bool]string{true: "x", false: "breaks"}[a[0] == nil])
		}
		switch a[0].(type) {
		case *Integer, *Double, *Logical:
		default:
			s.errorf(c.call, "'x' must be numeric")
		}
		if inheritsClass(a[0], "factor") {
			s.errorf(c.call, "'x' must be numeric")
		}
		x := asDoubles(coerceTo(a[0], kDouble))
		br := asDoubles(coerceTo(a[1], kDouble))
		right := a[4] == nil || s.truthy(c, a[4], "right")
		lowest := a[3] != nil && s.truthy(c, a[3], "include.lowest")
		ordered := a[6] != nil && s.truthy(c, a[6], "ordered_result")
		dig := 3
		if a[5] != nil {
			if d := asDoubles(coerceTo(a[5], kDouble)); len(d) == 1 {
				dig = int(d[0])
			}
		}
		if len(br) == 1 {
			if math.IsNaN(br[0]) || br[0] < 2 {
				s.errorf(c.call, "invalid number of intervals")
			}
			if br[0] > float64(s.opts.MaxVector) {
				s.errorf(c.call, "cannot allocate vector of length %.0f", br[0]+1)
			}
			nb := int(br[0]) + 1
			s.checkLen(c.call, nb)
			lo, hi := math.Inf(1), math.Inf(-1)
			for _, v := range x {
				if !math.IsNaN(v) {
					lo, hi = math.Min(lo, v), math.Max(hi, v)
				}
			}
			dx := hi - lo
			from, to := lo, hi
			if dx == 0 {
				dx = math.Abs(lo)
				if lo == 0 {
					dx = 1
				}
				from, to = lo-dx/1000, hi+dx/1000
			}
			br = make([]float64, nb)
			by := (to - from) / float64(nb-1)
			for i := range br {
				br[i] = from + float64(i)*by
			}
			br[nb-1] = to
			if hi-lo != 0 {
				br[0], br[nb-1] = lo-dx/1000, hi+dx/1000
			}
		} else {
			br = append([]float64(nil), br...)
			sort.Float64s(br)
		}
		nb := len(br)
		for i := 1; i < nb; i++ {
			if br[i] == br[i-1] {
				s.errorf(c.call, "'breaks' are not unique")
			}
		}
		codesOnly := false
		var labels []Str
		switch l := a[2].(type) {
		case nil, nullT:
			ok := false
			var ch []string
			for d := dig; d <= max(12, dig); d++ {
				ch = make([]string, nb)
				for i, v := range br {
					ch[i] = cFormatG(v, d)
				}
				ok = true
				for i := 1; i < nb; i++ {
					ok = ok && ch[i] != ch[i-1]
				}
				if ok {
					break
				}
			}
			for i := 0; i+1 < nb; i++ {
				if !ok {
					labels = append(labels, Str{S: fmt.Sprintf("Range_%d", i+1)})
					continue
				}
				open, shut := "[", ")"
				if right {
					open, shut = "(", "]"
				}
				labels = append(labels, Str{S: open + ch[i] + "," + ch[i+1] + shut})
			}
			if ok && lowest && nb > 1 {
				if right {
					labels[0].S = "[" + labels[0].S[1:]
				} else {
					k := len(labels) - 1
					labels[k].S = labels[k].S[:len(labels[k].S)-1] + "]"
				}
			}
		case *Logical:
			if len(l.V) == 1 && l.V[0] == 0 {
				codesOnly = true
			} else {
				labels = s.factorStrings(l)
			}
		default:
			labels = s.factorStrings(l)
		}
		if !codesOnly && len(labels) != nb-1 {
			s.errorf(c.call, "number of intervals and length of 'labels' differ")
		}
		codes := builtins[".bincode"].fn(s, &callInfo{call: c.call, env: c.env,
			args: []Value{&Double{V: x}, &Double{V: br}, mkLogical(right), mkLogical(lowest)}, names: []string{"", "", "", ""}}).(*Integer)
		if codesOnly {
			return codes
		}
		cls := []string{"factor"}
		if ordered {
			cls = []string{"ordered", "factor"}
		}
		codes.Attr = []Attr{{Name: "levels", Value: &String{V: labels}}, {Name: "class", Value: strVec(cls)}}
		return codes
	})
}
