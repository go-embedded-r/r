package eval

import (
	"math"
	"sort"
	"strings"
)

// Slice 5: vectors. as.vector, order, pmax/pmin, which.max/which.min,
// tabulate, row/column sums and means, t, diag, mean.

func init() {
	regBuiltin("as.vector", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "mode")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		mode := "any"
		if a[1] != nil {
			mode = s.stringArg(c, a[1], "mode")
		}
		x := a[0]
		if inheritsClass(x, "factor") {
			x = &String{V: s.factorStrings(x)}
		}
		switch mode {
		case "any":
			switch v := x.(type) {
			case *List:
				return withAttr(&List{V: v.V, Expr: v.Expr}, keepNamesOnly(v))
			case *Logical, *Integer, *Double, *String, *Complex, *Raw:
				return withAttr(x, nil)
			case *Closure, *Builtin:
				s.errorf(&Lang{Fun: Intern("as.vector"), Args: []Arg{{Value: Intern("x")}, {Value: Intern("mode")}}}, "cannot coerce type '%s' to vector of type 'any'", typeName(x))
			}
			return x
		case "list":
			if l, ok := x.(*List); ok {
				if l.Expr {
					return &List{V: l.V, Attr: l.Attr}
				}
				return l
			}
			return &List{V: s.asListElems(c, x), Attr: keepNamesOnly(x)}
		case "complex":
			cc := &callInfo{call: c.call, env: c.env, args: []Value{withAttr(x, nil)}, names: []string{""}}
			return s.asComplexFn(cc)
		case "raw":
			cc := &callInfo{call: c.call, env: c.env, args: []Value{withAttr(x, nil)}, names: []string{""}}
			return s.asRaw(cc)
		case "expression":
			return s.asExpression(c, x)
		case "logical", "integer", "numeric", "double", "character":
			if l, ok := x.(*List); ok {
				for _, e := range l.V {
					switch e.(type) {
					case *Logical, *Integer, *Double, *String:
						if Length(e) == 1 {
							continue
						}
					}
					s.errorf(c.call, "'list' object cannot be coerced to type '%s'", map[string]string{"logical": "logical", "integer": "integer", "numeric": "double", "double": "double", "character": "character"}[mode])
				}
				x = s.combine(&callInfo{call: c.call, env: c.env, args: l.V, names: make([]string, len(l.V))})
			}
			var k int
			switch mode {
			case "logical":
				k = kLogical
			case "integer":
				k = kInteger
			case "character":
				k = kString
			default:
				k = kDouble
			}
			if !isVector(x) {
				if _, isNull := x.(nullT); isNull {
					return emptyOf(k)
				}
				s.errorf(c.call, "cannot coerce type '%s' to vector of type '%s'", typeName(x), mode)
			}
			if _, isList := x.(*List); isList {
				s.errorf(c.call, "'list' object cannot be coerced to type '%s'", map[string]string{"logical": "logical", "integer": "integer", "numeric": "double", "double": "double", "character": "character"}[mode])
			}
			cc := &callInfo{call: c.call, env: c.env, args: []Value{withAttr(x, nil)}, names: []string{""}}
			switch k {
			case kLogical:
				return s.asLogical(cc)
			case kInteger:
				return s.asInteger(cc)
			case kString:
				return s.asCharacter(cc)
			}
			return s.asDouble(cc, "as.vector")
		}
		s.errorf(c.call, "invalid 'mode' argument")
		return nil
	})

	regBuiltin("order", func(s *Session, c *callInfo) Value {
		var keys []Value
		decreasing := []bool{false}
		naLast := int32(1)
		for i, a := range c.args {
			switch c.names[i] {
			case "decreasing":
				l := logicalOf(s, c, a, "'decreasing' must be TRUE or FALSE")
				decreasing = nil
				for _, v := range l {
					decreasing = append(decreasing, v == 1)
				}
				if len(decreasing) == 0 {
					s.errorf(c.call, "'decreasing' must be TRUE or FALSE")
				}
			case "na.last":
				// asLogical(CAR(args)) in R: the first element, NA when empty
				naLast = naInt
				if l := logicalOf(s, c, a, "invalid 'na.last' argument"); len(l) > 0 {
					naLast = l[0]
				}
			case "method":
			default:
				keys = append(keys, a)
			}
		}
		if len(keys) == 0 {
			return &Integer{V: []int32{}}
		}
		n := Length(keys[0])
		for ki, k := range keys {
			if Length(k) != n {
				s.errorf(c.call, "argument lengths differ")
			}
			switch k.(type) {
			case *Logical, *Integer, *Double, *String, *Complex:
			case nullT:
				s.errorf(c.call, "argument %d is not a vector", ki+1)
			default:
				s.errorf(c.call, "unimplemented type '%s' in 'orderVector1'\n", typeName(k))
			}
		}
		idx := make([]int, 0, n)
		for i := 0; i < n; i++ {
			idx = append(idx, i)
		}
		keyNA := func(k Value, i int) bool {
			switch v := k.(type) {
			case *Logical:
				return v.V[i] == naInt
			case *Integer:
				return v.V[i] == naInt
			case *Double:
				return math.IsNaN(v.V[i])
			case *String:
				return v.V[i].NA
			case *Complex:
				return math.IsNaN(real(v.V[i])) || math.IsNaN(imag(v.V[i]))
			}
			return false
		}
		cmp := func(k Value, i, j int) int {
			switch v := k.(type) {
			case *Logical:
				return cmpInt(int(v.V[i]), int(v.V[j]))
			case *Integer:
				return cmpInt(int(v.V[i]), int(v.V[j]))
			case *Double:
				switch {
				case v.V[i] < v.V[j]:
					return -1
				case v.V[i] > v.V[j]:
					return 1
				}
				return 0
			case *String:
				return collate(v.V[i].S, v.V[j].S)
			case *Complex:
				switch {
				case cless(v.V[i], v.V[j]):
					return -1
				case cless(v.V[j], v.V[i]):
					return 1
				}
				return 0
			}
			return 0
		}
		if naLast == naInt {
			var keep []int
			for _, i := range idx {
				ok := true
				for _, k := range keys {
					if keyNA(k, i) {
						ok = false
					}
				}
				if ok {
					keep = append(keep, i)
				}
			}
			idx = keep
		}
		sort.SliceStable(idx, func(p, q int) bool {
			i, j := idx[p], idx[q]
			for ki, k := range keys {
				dec := decreasing[ki%len(decreasing)]
				ni, nj := keyNA(k, i), keyNA(k, j)
				switch {
				case ni && nj:
					continue
				case ni:
					return naLast == 0
				case nj:
					return naLast != 0
				}
				r := cmp(k, i, j)
				if r == 0 {
					continue
				}
				if dec {
					return r > 0
				}
				return r < 0
			}
			return false
		})
		out := make([]int32, len(idx))
		for i, v := range idx {
			out[i] = int32(v + 1)
		}
		return &Integer{V: out}
	})

	pm := func(name string, wantMax bool) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			narm := false
			var args []Value
			for i, a := range c.args {
				if c.names[i] == "na.rm" {
					narm = s.truthy(c, a, "na.rm")
					continue
				}
				args = append(args, a)
			}
			if len(args) == 0 {
				s.errorf(c.call, "no arguments")
			}
			k := kLogical
			n := 0
			zero := false
			for _, a := range args {
				switch a.(type) {
				case *Logical, *Integer, *Double, *String:
				case nullT:
				default:
					s.errorf(c.call, "invalid input type")
				}
				if kindOf(a) > k {
					k = kindOf(a)
				}
				if Length(a) == 0 {
					zero = true
				}
				if Length(a) > n {
					n = Length(a)
				}
			}
			if zero {
				if k == kLogical {
					k = kInteger
				}
				return emptyOf(k)
			}
			if k == kLogical {
				k = kInteger
			}
			s.checkLen(c.call, n)
			conv := make([]Value, len(args))
			for i, a := range args {
				conv[i] = coerceTo(a, k)
				if Length(a)%1 != 0 || (n%Length(a) != 0) {
					s.warningf(c.call, "an argument will be fractionally recycled")
				}
			}
			var res Value
			switch k {
			case kInteger:
				out := make([]int32, n)
				for i := range out {
					var best int32
					have, na := false, false
					for _, a := range conv {
						e := a.(*Integer).V[i%Length(a)]
						if e == naInt {
							if !narm {
								na = true
							}
							continue
						}
						if !have || (wantMax && e > best) || (!wantMax && e < best) {
							best, have = e, true
						}
					}
					if na || !have {
						out[i] = naInt
					} else {
						out[i] = best
					}
				}
				res = &Integer{V: out}
			case kDouble:
				out := make([]float64, n)
				for i := range out {
					var best float64
					have, na := false, false
					var nan float64
					for _, a := range conv {
						e := a.(*Double).V[i%Length(a)]
						if math.IsNaN(e) {
							if !narm {
								na = true
								if nan == 0 || (isNAReal(e) && !isNAReal(nan)) {
									nan = e
								}
							}
							continue
						}
						if !have || (wantMax && e > best) || (!wantMax && e < best) {
							best, have = e, true
						}
					}
					switch {
					case na:
						out[i] = nan
					case !have:
						out[i] = naReal
					default:
						out[i] = best
					}
				}
				res = &Double{V: out}
			case kString:
				out := make([]Str, n)
				for i := range out {
					var best Str
					have, na := false, false
					for _, a := range conv {
						e := a.(*String).V[i%Length(a)]
						if e.NA {
							if !narm {
								na = true
							}
							continue
						}
						if !have || (wantMax && collate(e.S, best.S) > 0) || (!wantMax && collate(e.S, best.S) < 0) {
							best, have = e, true
						}
					}
					if na || !have {
						out[i] = Str{NA: true}
					} else {
						out[i] = best
					}
				}
				res = &String{V: out}
			}
			// the attributes of the first argument that has the result's length
			if Length(args[0]) == n {
				res = withAttr(res, attrsOf(args[0]))
			}
			return res
		})
	}
	pm("pmax", true)
	pm("pmin", false)

	which := func(name string, wantMax bool) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			x := s.nargs1(c, name)
			var d []float64
			switch x.(type) {
			case *Logical, *Integer, *Double:
				d = asDoubles(x)
			case *String:
				d = asDoubles(s.asDouble(&callInfo{call: c.call, env: c.env, args: []Value{x}, names: []string{""}}, name))
			case nullT:
			default:
				s.errorf(c.call, "cannot coerce type '%s' to vector of type 'double'", typeName(x))
			}
			best := -1
			for i, e := range d {
				if math.IsNaN(e) {
					continue
				}
				if best < 0 || (wantMax && e > d[best]) || (!wantMax && e < d[best]) {
					best = i
				}
			}
			if best < 0 {
				return &Integer{V: []int32{}}
			}
			res := &Integer{V: []int32{int32(best + 1)}}
			if nm := namesOf(x); nm != nil {
				res.Attr = namesAttr([]Str{nm.V[best]})
			}
			return res
		})
	}
	which("which.max", true)
	which("which.min", false)

	regBuiltin("tabulate", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "bin", "nbins")
		if a[0] == nil {
			s.errorf(c.call, "argument \"bin\" is missing, with no default")
		}
		if inheritsClass(a[0], "factor") {
			a[0] = withAttr(a[0], nil)
		}
		switch a[0].(type) {
		case *Integer, *Logical, *Double:
		default:
			s.errorf(c.call, "'bin' must be numeric or a factor")
		}
		bins := asInts(coerceTo(a[0], kInteger))
		nb := 1
		if a[1] != nil {
			d := asDoublesOrNil(a[1])
			if len(d) != 1 || math.IsNaN(d[0]) || d[0] < 0 {
				s.errorf(c.call, "invalid 'nbin' argument")
			}
			if d[0] > float64(s.opts.MaxVector) {
				s.errorf(c.call, "vector size too large")
			}
			nb = int(d[0])
		} else {
			for _, b := range bins {
				if b != naInt && int(b) > nb {
					nb = int(b)
				}
			}
		}
		s.checkLen(c.call, nb)
		out := make([]int32, nb)
		for _, b := range bins {
			if b != naInt && b >= 1 && int(b) <= nb {
				out[b-1]++
			}
		}
		return &Integer{V: out}
	})

	margin := func(name string, rows, mean bool) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			a := s.matchFormals(c, "x", "na.rm", "dims")
			x := a[0]
			if x == nil {
				s.errorf(c.call, "argument \"x\" is missing, with no default")
			}
			if isObject(x) && inheritsClass(x, "data.frame") {
				s.errorf(c.call, "data frames are not supported")
			}
			dims := dimsOf(x)
			if len(dims) < 2 {
				s.errorf(c.call, "'x' must be an array of at least two dimensions")
			}
			switch x.(type) {
			case *Logical, *Integer, *Double:
			default:
				s.errorf(c.call, "'x' must be numeric or complex")
			}
			narm := a[1] != nil && s.truthy(c, a[1], "na.rm")
			nr, nc := dims[0], 1
			for _, d := range dims[1:] {
				nc *= d
			}
			v := asDoubles(x)
			var out []float64
			if rows {
				out = make([]float64, nr)
				cnt := make([]int, nr)
				for j := 0; j < nc; j++ {
					for i := 0; i < nr; i++ {
						e := v[i+j*nr]
						if narm && math.IsNaN(e) {
							continue
						}
						out[i] += e
						cnt[i]++
					}
				}
				if mean {
					for i := range out {
						if cnt[i] == 0 {
							out[i] = math.NaN()
						} else {
							out[i] /= float64(cnt[i])
						}
					}
				}
			} else {
				out = make([]float64, nc)
				for j := 0; j < nc; j++ {
					var sum float64
					cnt := 0
					for i := 0; i < nr; i++ {
						e := v[i+j*nr]
						if narm && math.IsNaN(e) {
							continue
						}
						sum += e
						cnt++
					}
					if mean {
						if cnt == 0 {
							sum = math.NaN()
						} else {
							sum /= float64(cnt)
						}
					}
					out[j] = sum
				}
			}
			res := &Double{V: out}
			if len(dims) == 2 {
				d := 0
				if !rows {
					d = 1
				}
				if nm := dimNamesOf(x, d); nm != nil {
					res.Attr = namesAttr(nm.V)
				}
			}
			return res
		})
	}
	margin("rowSums", true, false)
	margin("colSums", false, false)
	margin("rowMeans", true, true)
	margin("colMeans", false, true)

	regBuiltin("t", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "t")
		if !isVector(x) {
			s.errorf(c.call, "argument is not a matrix")
		}
		dims := dimsOf(x)
		switch len(dims) {
		case 0:
			n := Length(x)
			res := withAttr(x, nil)
			attr := []Attr{{Name: "dim", Value: &Integer{V: []int32{1, int32(n)}}}}
			if nm := namesOf(x); nm != nil {
				attr = append(attr, Attr{Name: "dimnames", Value: &List{V: []Value{Null, nm}}})
			}
			return withAttr(res, attr)
		case 2:
			nr, nc := dims[0], dims[1]
			pos := make([]int, nr*nc)
			for i := 0; i < nr; i++ {
				for j := 0; j < nc; j++ {
					pos[j+i*nc] = i + j*nr
				}
			}
			res := takeElems(x, pos)
			attr := []Attr{{Name: "dim", Value: &Integer{V: []int32{int32(nc), int32(nr)}}}}
			if dn, ok := getAttr(x, "dimnames").(*List); ok && len(dn.V) == 2 {
				l := &List{V: []Value{dn.V[1], dn.V[0]}}
				if nn := namesOf(dn); nn != nil && len(nn.V) == 2 {
					l.Attr = namesAttr([]Str{nn.V[1], nn.V[0]})
				}
				attr = append(attr, Attr{Name: "dimnames", Value: l})
			}
			for _, at := range attrsOf(x) {
				if at.Name != "dim" && at.Name != "dimnames" && at.Name != "names" {
					attr = append(attr, at)
				}
			}
			return withAttr(res, attr)
		}
		s.errorf(&Lang{Fun: Intern("t.default"), Args: c.call.Args}, "argument is not a matrix")
		return nil
	})

	regBuiltin("diag", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "nrow", "ncol", "names")
		x := a[0]
		if x == nil {
			x = mkDouble(1)
		}
		switch x.(type) {
		case *Logical, *Integer, *Double:
		default:
			s.errorf(c.call, "invalid 'nrow' value (too large or NA)")
		}
		if d := dimsOf(x); len(d) == 2 {
			m := min(d[0], d[1])
			pos := make([]int, m)
			for i := range pos {
				pos[i] = i + i*d[0]
			}
			return takeElems(withAttr(x, nil), pos)
		}
		var n, nc int
		vals := asDoubles(x)
		switch {
		case a[1] == nil && a[2] == nil && Length(x) == 1 && !isObject(x):
			n = int(vals[0])
			if vals[0] < 0 || math.IsNaN(vals[0]) {
				s.errorf(c.call, "invalid 'nrow' value (< 0)")
			}
			vals = []float64{1}
			nc = n
		default:
			n = Length(x)
			if a[1] != nil {
				n = s.countArg(c, a[1], "nrow")
			}
			nc = n
			if a[2] != nil {
				nc = s.countArg(c, a[2], "ncol")
			} else if a[1] == nil {
				nc = n
			}
		}
		if n > 0 && nc > s.opts.MaxVector/n {
			s.errorf(c.call, "vector size too large")
		}
		s.checkLen(c.call, n*nc)
		out := make([]float64, n*nc)
		if len(vals) > 0 {
			for i := 0; i < min(n, nc); i++ {
				out[i+i*n] = vals[i%len(vals)]
			}
		}
		_, isInt := x.(*Integer)
		var res Value = &Double{V: out}
		if isInt && !(a[1] == nil && a[2] == nil && Length(x) == 1) {
			res = coerceTo(res, kInteger)
		}
		return withAttr(res, []Attr{{Name: "dim", Value: &Integer{V: []int32{int32(n), int32(nc)}}}})
	})

	regBuiltin("mean", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "trim", "na.rm")
		x := a[0]
		if x == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		narm := a[2] != nil && s.truthy(c, a[2], "na.rm")
		switch xv := x.(type) {
		case *Logical, *Integer, *Double:
		case *Complex:
			var zs []complex128
			for _, z := range xv.V {
				if narm && (math.IsNaN(real(z)) || math.IsNaN(imag(z))) {
					continue
				}
				zs = append(zs, z)
			}
			var sum complex128
			for _, z := range zs {
				sum += z
			}
			n := complex(float64(len(zs)), 0)
			m := complex(real(sum)/real(n), imag(sum)/real(n))
			if len(zs) > 0 && !math.IsInf(real(m), 0) && !math.IsNaN(real(m)) && !math.IsInf(imag(m), 0) && !math.IsNaN(imag(m)) {
				var t complex128
				for _, z := range zs {
					t += z - m
				}
				m += complex(real(t)/real(n), imag(t)/real(n))
			}
			return &Complex{V: []complex128{m}}
		default:
			s.warningf(c.call, "argument is not numeric or logical: returning NA")
			return &Double{V: []float64{naReal}}
		}
		d := asDoubles(x)
		if _, isI := x.(*Integer); isI || kindOf(x) == kLogical {
			var sum float64
			n := 0
			for _, e := range asInts(x) {
				if e == naInt {
					if narm {
						continue
					}
					return &Double{V: []float64{naReal}}
				}
				sum += float64(e)
				n++
			}
			if n == 0 {
				return &Double{V: []float64{math.NaN()}}
			}
			return mkDouble(sum / float64(n))
		}
		if narm {
			var keep []float64
			for _, e := range d {
				if !math.IsNaN(e) {
					keep = append(keep, e)
				}
			}
			d = keep
		}
		if len(d) == 0 {
			return &Double{V: []float64{math.NaN()}}
		}
		var sum float64
		for _, e := range d {
			sum += e
		}
		m := sum / float64(len(d))
		if math.IsInf(sum, 0) {
			finite := true
			for _, e := range d {
				if math.IsInf(e, 0) {
					finite = false
				}
			}
			if finite {
				m = 0
				for _, e := range d {
					m += e / float64(len(d))
				}
			}
		}
		if !math.IsInf(m, 0) && !math.IsNaN(m) {
			var t float64
			for _, e := range d {
				t += e - m
			}
			m += t / float64(len(d))
		}
		return mkDouble(m)
	})
}

func cmpInt(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func keepNamesOnly(x Value) []Attr {
	if nm := namesOf(x); nm != nil {
		return namesAttr(nm.V)
	}
	return nil
}

func init() {
	regSpecial("stopifnot", func(s *Session, c *callInfo) Value {
		for _, a := range c.raw {
			if a.HasName && (a.Name == "exprs" || a.Name == "exprObject" || a.Name == "local") {
				continue
			}
			v := s.eval(a.Value, c.env)
			ok := false
			if l, isL := v.(*Logical); isL {
				ok = true
				for _, e := range l.V {
					if e != 1 {
						ok = false
						break
					}
				}
			}
			if ok {
				continue
			}
			msg := a.Name
			if !a.HasName || msg == "" {
				d := ""
				for i, line := range syntaxDeparse(a.Value) {
					if i > 0 {
						d += " "
					}
					d += line
				}
				switch l, isL := v.(*Logical); {
				case isL && len(l.V) > 1:
					msg = d + " are not all TRUE"
				default:
					msg = d + " is not TRUE"
				}
			}
			var call *Lang
			if _, fr := s.frameIndex(c.env); fr != nil {
				call = fr.call
			}
			s.errorf(call, "%s", msg)
		}
		s.visible = false
		return Null
	})

	regBuiltin("rank", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "na.last", "ties.method")
		x := a[0]
		if x == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		ties := "average"
		if a[2] != nil {
			ties = s.stringArg(c, a[2], "ties.method")
		}
		switch ties {
		case "average", "first", "min", "max":
		default:
			s.errorf(c.call, "'arg' should be one of %s", "‘average’, ‘first’, ‘last’, ‘random’, ‘max’, ‘min’")
		}
		switch x.(type) {
		case *Logical, *Integer, *Double, *String, *Complex:
		default:
			s.errorf(c.call, "argument is not a vector")
		}
		n := Length(x)
		oc := &callInfo{call: c.call, env: c.env, args: []Value{x}, names: []string{""}}
		ord := builtins["order"].fn(s, oc).(*Integer).V
		keyEq := func(i, j int) bool { return keyOfElem(x, i) == keyOfElem(x, j) }
		ranks := make([]float64, n)
		for i := 0; i < n; {
			j := i
			for j+1 < n && keyEq(int(ord[i])-1, int(ord[j+1])-1) {
				j++
			}
			for k := i; k <= j; k++ {
				var r float64
				switch ties {
				case "average":
					r = float64(i+j)/2 + 1
				case "min":
					r = float64(i + 1)
				case "max":
					r = float64(j + 1)
				case "first":
					r = float64(k + 1)
				}
				ranks[int(ord[k])-1] = r
			}
			i = j + 1
		}
		var res Value = &Double{V: ranks}
		if ties != "average" {
			res = coerceTo(res, kInteger)
		}
		if nm := namesOf(x); nm != nil {
			res = withAttr(res, namesAttr(nm.V))
		}
		return res
	})
}

func init() {
	regSpecial("match.arg", func(s *Session, c *callInfo) Value {
		if len(c.raw) == 0 {
			s.errorf(c.call, "argument \"arg\" is missing, with no default")
		}
		var argE, chE Value
		several := false
		pos := 0
		for _, a := range c.raw {
			switch {
			case a.HasName && a.Name == "arg":
				argE = a.Value
			case a.HasName && a.Name == "choices":
				chE = a.Value
			case a.HasName && a.Name == "several.ok":
				several = s.lglFlag(c, a.Value, "several.ok")
			case !a.HasName:
				switch pos {
				case 0:
					argE = a.Value
				case 1:
					chE = a.Value
				case 2:
					several = s.lglFlag(c, a.Value, "several.ok")
				}
				pos++
			}
		}
		if argE == nil {
			s.errorf(c.call, "argument \"arg\" is missing, with no default")
		}
		var choices Value
		identicalToDefault := false
		if chE != nil {
			choices = s.eval(chE, c.env)
		} else {
			sym, ok := argE.(*Symbol)
			if !ok {
				s.errorf(c.call, "'arg' must be a symbol when 'choices' is missing")
			}
			_, fr := s.frameIndex(c.env)
			if fr == nil {
				s.errorf(c.call, "match.arg() must be called from a function")
			}
			var def Value
			for _, f := range fr.fn.Formals {
				if f.Name == sym {
					def = f.Default
				}
			}
			if def == nil {
				s.errorf(c.call, "no default choices for '%s'", sym.Name)
			}
			choices = s.eval(def, c.env)
			if s.isMissingArg(sym, c.env) {
				identicalToDefault = true
			}
		}
		arg := s.eval(argE, c.env)
		cs, ok := choices.(*String)
		if !ok {
			s.errorf(c.call, "'choices' must be a character vector")
		}
		if _, isNull := arg.(nullT); isNull {
			return &String{V: []Str{cs.V[0]}}
		}
		as, ok := arg.(*String)
		if !ok {
			s.errorf(c.call, "'arg' must be NULL or a character vector")
		}
		if !several {
			if identicalToDefault || identical(arg, choices) {
				return &String{V: []Str{cs.V[0]}}
			}
			if len(as.V) != 1 {
				s.errorf(c.call, "'arg' must be of length 1")
			}
		} else if len(as.V) == 0 {
			s.errorf(c.call, "'arg' must be of length >= 1")
		}
		var out []Str
		for _, a := range as.V {
			hit := -1
			exact := false
			n := 0
			for j, ch := range cs.V {
				if ch.S == a.S {
					hit, exact = j, true
					break
				}
				if strings.HasPrefix(ch.S, a.S) && a.S != "" {
					hit = j
					n++
				}
			}
			if !exact && n != 1 {
				hit = -1
			}
			if hit < 0 {
				if several {
					continue
				}
				var q []string
				for _, ch := range cs.V {
					q = append(q, "“"+ch.S+"”")
				}
				s.errorf(c.call, "'arg' should be one of %s", strings.Join(q, ", "))
			}
			out = append(out, cs.V[hit])
		}
		if len(out) == 0 {
			s.errorf(c.call, "'arg' should be one of ...")
		}
		return &String{V: out}
	})
}
