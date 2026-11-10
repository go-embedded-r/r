package eval

import (
	"math"
	"strings"
)

// The primitives of slice 4: the rest of the math group, cumulative
// functions, the type predicates, the language primitives (switch,
// substitute, call, ...) and the matrix products. Written from R's
// documentation and the oracle's behaviour.

func init() {
	for name, f := range map[string]func(float64) float64{
		"sin": math.Sin, "cos": math.Cos, "tan": math.Tan,
		"asin": math.Asin, "acos": math.Acos, "atan": math.Atan,
		"sinh": math.Sinh, "cosh": math.Cosh, "tanh": math.Tanh,
		"asinh": math.Asinh, "acosh": math.Acosh, "atanh": math.Atanh,
		"log10": log10Fn, "log2": log2Fn, "log1p": math.Log1p, "expm1": math.Expm1,
		"gamma": gammaFn, "lgamma": lgammaFn,
		"sinpi": sinpi, "cospi": cospi, "tanpi": tanpi,
		"trunc": math.Trunc,
	} {
		regBuiltin(name, (*Session)(nil).math1(name, f, false))
	}
	regBuiltin("sign", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "sign")
		switch v := x.(type) {
		case *Integer, *Logical:
			in := asInts(v)
			out := make([]float64, len(in))
			for i, e := range in {
				switch {
				case e == naInt:
					out[i] = naReal
				case e > 0:
					out[i] = 1
				case e < 0:
					out[i] = -1
				}
			}
			return &Double{V: out, Attr: attrsOf(x)}
		case *Double:
			out := make([]float64, len(v.V))
			for i, e := range v.V {
				switch {
				case math.IsNaN(e):
					out[i] = e
				case e > 0:
					out[i] = 1
				case e < 0:
					out[i] = -1
				}
			}
			return &Double{V: out, Attr: v.Attr}
		}
		s.errorf(c.call, "non-numeric argument to mathematical function")
		return nil
	})
	regBuiltin("log", func(s *Session, c *callInfo) Value {
		if len(c.args) == 0 {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		a := s.matchFormals(c, "x", "base")
		x := a[0]
		switch xv := x.(type) {
		case *Logical, *Integer, *Double:
		case *Complex:
			var base *complex128
			if a[1] != nil {
				switch a[1].(type) {
				case *Logical, *Integer, *Double, *Complex:
				default:
					s.errorf(c.call, "non-numeric argument to mathematical function")
				}
				if bz := asComplex(a[1]); len(bz) > 0 {
					base = &bz[0]
				}
			}
			out := make([]complex128, len(xv.V))
			for i, z := range xv.V {
				out[i] = complexLog(z, base)
			}
			return &Complex{V: out, Attr: xv.Attr}
		default:
			s.errorf(c.call, "non-numeric argument to mathematical function")
		}
		f := math.Log
		if a[1] != nil {
			bd := asDoublesOrNil(a[1])
			if len(bd) < 1 || kindOf(a[1]) == kString || kindOf(a[1]) == kOther {
				s.errorf(c.call, "non-numeric argument to mathematical function")
			}
			base := bd[0]
			switch {
			case base == 10:
				f = log10Fn
			case base == 2:
				f = log2Fn
			default:
				lb := math.Log(base)
				f = func(v float64) float64 { return math.Log(v) / lb }
			}
		}
		d := asDoubles(x)
		out := make([]float64, len(d))
		nan := false
		for i, e := range d {
			if math.IsNaN(e) {
				out[i] = e
				continue
			}
			out[i] = f(e)
			if math.IsNaN(out[i]) {
				nan = true
			}
		}
		if nan {
			s.warningf(c.call, "NaNs produced")
		}
		return &Double{V: out, Attr: attrsOf(x)}
	})
	regBuiltin("atan2", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "y", "x")
		if a[0] == nil || a[1] == nil {
			s.errorf(c.call, "argument \"%s\" is missing, with no default", map[bool]string{true: "y", false: "x"}[a[0] == nil])
		}
		y, x := asDoubles(a[0]), asDoubles(a[1])
		n := max(len(y), len(x))
		if len(y) == 0 || len(x) == 0 {
			return &Double{V: []float64{}}
		}
		out := make([]float64, n)
		for i := range out {
			out[i] = math.Atan2(y[i%len(y)], x[i%len(x)])
		}
		return &Double{V: out}
	})
	regBuiltin("round", func(s *Session, c *callInfo) Value { return s.roundLike(c, false) })
	regBuiltin("signif", func(s *Session, c *callInfo) Value { return s.roundLike(c, true) })

	cum := func(name string) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			x := s.nargs1(c, name)
			if _, isNull := x.(nullT); isNull {
				return &Double{V: []float64{}}
			}
			switch xv := x.(type) {
			case *Logical, *Integer, *Double:
			case *Complex:
				if name == "cumsum" || name == "cumprod" {
					out := make([]complex128, len(xv.V))
					var acc complex128
					if name == "cumprod" {
						acc = 1
					}
					for i, z := range xv.V {
						if name == "cumsum" {
							acc += z
						} else {
							acc = cmul(acc, z)
						}
						out[i] = acc
					}
					var nmAttr []Attr
					if n := namesOf(x); n != nil {
						nmAttr = namesAttr(n.V)
					}
					return &Complex{V: out, Attr: nmAttr}
				}
				s.errorf(c.call, "'%s' not defined for complex numbers", name)
			case *String:
				x = s.asDouble(&callInfo{call: c.call, env: c.env, args: []Value{x}, names: []string{""}}, name)
			default:
				s.errorf(c.call, "'%s' not defined for this type", name)
			}
			var nm []Attr
			if n := namesOf(x); n != nil {
				nm = namesAttr(n.V)
			}
			_, isDbl := x.(*Double)
			if !isDbl && name != "cumprod" {
				in := asInts(x)
				out := make([]int32, len(in))
				var acc int64
				na := false
				for i, e := range in {
					if na || e == naInt {
						na = true
						out[i] = naInt
						continue
					}
					switch name {
					case "cumsum":
						acc += int64(e)
						if acc > math.MaxInt32 || acc < -math.MaxInt32 {
							s.warningf(c.call, "integer overflow in 'cumsum'; use 'cumsum(as.numeric(.))'")
							na = true
							out[i] = naInt
							continue
						}
					case "cummax":
						if i == 0 || int64(e) > acc {
							acc = int64(e)
						}
					case "cummin":
						if i == 0 || int64(e) < acc {
							acc = int64(e)
						}
					}
					out[i] = int32(acc)
				}
				return &Integer{V: out, Attr: nm}
			}
			d := asDoubles(x)
			out := make([]float64, len(d))
			var acc float64
			switch name {
			case "cumprod":
				acc = 1
			}
			for i, e := range d {
				switch name {
				case "cumsum":
					acc += e
				case "cumprod":
					acc *= e
				case "cummax":
					if i == 0 || math.IsNaN(acc) {
						if i == 0 {
							acc = e
						}
					} else if math.IsNaN(e) || e > acc {
						acc = e
					}
				case "cummin":
					if i == 0 || math.IsNaN(acc) {
						if i == 0 {
							acc = e
						}
					} else if math.IsNaN(e) || e < acc {
						acc = e
					}
				}
				out[i] = acc
			}
			return &Double{V: out, Attr: nm}
		})
	}
	for _, n := range []string{"cumsum", "cumprod", "cummax", "cummin"} {
		cum(n)
	}

	regBuiltin("range", func(s *Session, c *callInfo) Value {
		finite := false
		var args []Value
		var names []string
		for i, a := range c.args {
			if c.names[i] == "finite" {
				finite = s.truthy(c, a, "finite")
				continue
			}
			args = append(args, a)
			names = append(names, c.names[i])
		}
		hasRm := false
		for _, n := range names {
			if n == "na.rm" {
				hasRm = true
			}
		}
		if finite {
			for i, a := range args {
				if names[i] == "na.rm" {
					continue
				}
				if d, ok := a.(*Double); ok {
					var keep []float64
					for _, e := range d.V {
						if !math.IsNaN(e) && !math.IsInf(e, 0) {
							keep = append(keep, e)
						}
					}
					args[i] = &Double{V: keep}
				}
			}
			if !hasRm {
				args = append(args, mkLogical(true))
				names = append(names, "na.rm")
			}
		}
		sub := &callInfo{call: c.call, env: c.env, args: args, names: names}
		lo := builtins["min"].fn(s, sub)
		hi := builtins["max"].fn(s, sub)
		return s.combine(&callInfo{call: c.call, env: c.env, args: []Value{lo, hi}, names: []string{"", ""}})
	})

	predicate := func(name string, f func(float64) bool, strOK bool) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			x := s.nargs1(c, name)
			var out []int32
			switch v := x.(type) {
			case *Double:
				out = make([]int32, len(v.V))
				for i, e := range v.V {
					if f(e) {
						out[i] = 1
					}
				}
			case *Complex:
				out = make([]int32, len(v.V))
				for i, z := range v.V {
					re, im := real(z), imag(z)
					switch name {
					case "is.finite":
						out[i] = b2i(!math.IsNaN(re) && !math.IsInf(re, 0) && !math.IsNaN(im) && !math.IsInf(im, 0))
					case "is.infinite":
						out[i] = b2i(math.IsInf(re, 0) || math.IsInf(im, 0))
					default:
						out[i] = b2i((math.IsNaN(re) && !isNAReal(re)) || (math.IsNaN(im) && !isNAReal(im)))
					}
				}
			case *Raw:
				out = make([]int32, len(v.V))
				if name == "is.finite" {
					for i := range out {
						out[i] = 1
					}
				}
			case *Integer, *Logical:
				in := asInts(v)
				out = make([]int32, len(in))
				for i, e := range in {
					out[i] = b2i(name == "is.finite" && e != naInt)
				}
			case *String:
				if !strOK {
					s.errorf(c.call, "default method not implemented for type 'character'")
				}
				out = make([]int32, len(v.V))
			case nullT:
				return &Logical{V: []int32{}}
			default:
				s.errorf(c.call, "default method not implemented for type '%s'", typeName(x))
			}
			return withAttr(&Logical{V: out}, keepShapeAttrs(x))
		})
	}
	predicate("is.finite", func(x float64) bool { return !math.IsNaN(x) && !math.IsInf(x, 0) }, true)
	predicate("is.infinite", func(x float64) bool { return math.IsInf(x, 0) }, true)
	predicate("is.nan", func(x float64) bool { return math.IsNaN(x) && !isNAReal(x) }, true)

	regBuiltin("anyNA", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "recursive")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		rec := a[1] != nil && s.truthy(c, a[1], "recursive")
		return mkLogical(anyNAValue(a[0], rec))
	})
	regBuiltin("nzchar", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "keepNA")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		keep := int32(0)
		if a[1] != nil {
			keep = logicalOf(s, c, a[1], "invalid 'keepNA' value")[0]
		}
		var x []Str
		switch v := a[0].(type) {
		case *String:
			x = v.V
		case *Logical, *Integer, *Double:
			x = asStrings(v)
		case nullT:
			return &Logical{V: []int32{}}
		default:
			s.errorf(c.call, "'nzchar()' requires a character vector")
		}
		out := make([]int32, len(x))
		for i, e := range x {
			switch {
			case e.NA && keep == 1:
				out[i] = naInt
			case e.NA && keep == naInt:
				out[i] = 1
			case e.NA:
				out[i] = 1
			default:
				out[i] = b2i(e.S != "")
			}
		}
		return &Logical{V: out}
	})

	typeTest := func(name string, f func(v Value) bool) {
		regBuiltin(name, func(s *Session, c *callInfo) Value { return mkLogical(f(s.nargs1(c, name))) })
	}
	typeTest("is.matrix", func(v Value) bool { return len(dimsOf(v)) == 2 })
	typeTest("is.array", func(v Value) bool { return len(dimsOf(v)) >= 1 })
	typeTest("is.object", isObject)
	typeTest("is.recursive", func(v Value) bool {
		switch v.(type) {
		case *List, *Closure, *Builtin, *Env, *Lang:
			return true
		}
		return false
	})
	typeTest("is.language", func(v Value) bool {
		switch x := v.(type) {
		case *Lang, *Symbol:
			return true
		case *List:
			return x.Expr
		}
		return false
	})
	typeTest("is.pairlist", func(v Value) bool { _, null := v.(nullT); return null })
	typeTest("is.expression", func(v Value) bool { l, ok := v.(*List); return ok && l.Expr })
	typeTest("is.raw", func(v Value) bool { _, ok := v.(*Raw); return ok })
	typeTest("isS4", func(Value) bool { return false })
	regBuiltin("interactive", func(s *Session, c *callInfo) Value { return mkLogical(false) })
	regBuiltin("gc.time", func(s *Session, c *callInfo) Value { return &Double{V: make([]float64, 5)} })

	regBuiltin("length<-", func(s *Session, c *callInfo) Value {
		if len(c.args) != 2 {
			s.errorf(c.call, "2 arguments are required")
		}
		x := c.args[0]
		d := asDoublesOrNil(c.args[1])
		if len(d) != 1 || d[0] != d[0] || d[0] < 0 || kindOf(c.args[1]) == kString {
			s.errorf(c.call, "invalid value")
		}
		if d[0] > float64(s.opts.MaxVector) {
			s.errorf(c.call, "vector size too large")
		}
		n := int(d[0])
		s.checkLen(c.call, n)
		if _, isNull := x.(nullT); isNull {
			if n == 0 {
				return Null
			}
			s.errorf(c.call, "invalid argument")
		}
		if !isVector(x) {
			s.errorf(c.call, "invalid argument")
		}
		pos := make([]int, n)
		for i := range pos {
			pos[i] = i
		}
		if l, ok := x.(*List); ok {
			out := make([]Value, n)
			for i := range out {
				if i < len(l.V) {
					out[i] = l.V[i]
				} else {
					out[i] = Null
				}
			}
			res := &List{V: out, Expr: l.Expr}
			if nm := namesOf(l); nm != nil {
				res.Attr = namesAttr(fitNamesEmpty(nm.V, n))
			}
			return res
		}
		res := takeElems(x, pos)
		if nm := namesOf(x); nm != nil {
			res = withAttr(res, namesAttr(fitNamesEmpty(nm.V, n)))
		}
		return res
	})
	regBuiltin("storage.mode<-", func(s *Session, c *callInfo) Value {
		if len(c.args) != 2 {
			s.errorf(c.call, "2 arguments are required")
		}
		mode := s.stringArg(c, c.args[1], "value")
		var k int
		switch mode {
		case "logical":
			k = kLogical
		case "integer":
			k = kInteger
		case "double", "numeric":
			k = kDouble
		case "character":
			k = kString
		case "list":
			k = kList
		default:
			s.errorf(c.call, "invalid value")
		}
		if !isVector(c.args[0]) {
			s.errorf(c.call, "invalid value")
		}
		if kindOf(c.args[0]) == k {
			return c.args[0]
		}
		return coerceTo(c.args[0], k)
	})
	regBuiltin("xtfrm", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "xtfrm")
		switch x.(type) {
		case *Integer, *Double:
			if isObject(x) {
				return setAttr(x, "class", Null)
			}
			return x
		}
		s.errorf(c.call, "xtfrm is only implemented for numeric vectors")
		return nil
	})

	// language primitives
	regSpecial("switch", s3Switch)
	regSpecial("substitute", func(s *Session, c *callInfo) Value {
		if len(c.raw) == 0 {
			s.errorf(c.call, "argument \"expr\" is missing, with no default")
		}
		env := c.env
		if len(c.raw) > 1 {
			ev := s.eval(c.raw[1].Value, c.env)
			switch e := ev.(type) {
			case *Env:
				env = e
			case *List:
				env = s.evalEnvArg(c, e, s.Empty)
			default:
				s.errorf(c.call, "invalid environment specified")
			}
		}
		if env == s.Global {
			return c.raw[0].Value
		}
		return s.substitute(c.raw[0].Value, env)
	})
	regBuiltin("nargs", func(s *Session, c *callInfo) Value {
		_, fr := s.frameIndex(c.env)
		if fr == nil {
			return mkInt(0)
		}
		n := 0
		for _, a := range fr.supplied {
			if !isMissingValue(a.Value) {
				n++
			}
		}
		return mkInt(int32(n))
	})
	regBuiltin("...length", func(s *Session, c *callInfo) Value {
		d := s.dotsOf(c.env)
		if d == nil {
			s.errorf(nil, "incorrect context: the current call has no '...' to look in")
		}
		return mkInt(int32(len(d.Args)))
	})
	regBuiltin("...elt", func(s *Session, c *callInfo) Value {
		d := s.dotsOf(c.env)
		if d == nil {
			s.errorf(c.call, "incorrect context: the current call has no '...' to look in")
		}
		n := s.intArg(c, s.nargs1(c, "...elt"), "n")
		if n < 1 || n > len(d.Args) {
			s.errorf(s.currentCall(), "the ... list contains fewer than %d elements", n)
		}
		return s.argValue(d.Args[n-1].Value)
	})
	regBuiltin("...names", func(s *Session, c *callInfo) Value {
		d := s.dotsOf(c.env)
		if d == nil {
			s.errorf(c.call, "incorrect context: the current call has no '...' to look in")
		}
		any := false
		out := make([]Str, len(d.Args))
		for i, a := range d.Args {
			out[i] = Str{S: a.Name}
			if a.Name != "" {
				any = true
			}
		}
		if !any {
			return Null
		}
		return &String{V: out}
	})
	regBuiltin("call", func(s *Session, c *callInfo) Value {
		if len(c.args) == 0 {
			s.errorf(c.call, "first argument must be a character string")
		}
		st, ok := c.args[0].(*String)
		if !ok || len(st.V) != 1 || st.V[0].NA {
			s.errorf(c.call, "first argument must be a character string")
		}
		l := &Lang{Fun: Intern(st.V[0].S)}
		for i := 1; i < len(c.args); i++ {
			l.Args = append(l.Args, Arg{Name: c.names[i], HasName: c.names[i] != "", Value: c.args[i]})
		}
		return l
	})
	regBuiltin("as.call", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "as.call")
		switch v := x.(type) {
		case *Lang:
			return v
		case *List:
			if len(v.V) == 0 {
				s.errorf(c.call, "invalid length 0 argument")
			}
			nm := namesOf(v)
			l := &Lang{Fun: v.V[0]}
			for i := 1; i < len(v.V); i++ {
				a := Arg{Value: v.V[i]}
				if nm != nil && nm.V[i].S != "" && !nm.V[i].NA {
					a.Name, a.HasName = nm.V[i].S, true
				}
				l.Args = append(l.Args, a)
			}
			return l
		}
		s.errorf(c.call, "invalid argument list")
		return nil
	})
	regSpecial("forceAndCall", func(s *Session, c *callInfo) Value {
		if len(c.raw) < 2 {
			s.errorf(c.call, "argument \"FUN\" is missing, with no default")
		}
		call := &Lang{Fun: c.raw[1].Value, Args: c.raw[2:]}
		return s.eval(call, c.env)
	})
	regBuiltin(".Primitive", func(s *Session, c *callInfo) Value {
		name := s.stringArg(c, s.nargs1(c, ".Primitive"), "name")
		if b, ok := builtins[name]; ok {
			return b
		}
		s.errorf(c.call, "no such primitive function")
		return nil
	})
	subset := func(name, target string) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			saved := s.noDispatch
			s.noDispatch = true
			defer func() { s.noDispatch = saved }()
			ci := &callInfo{call: c.call, env: c.env, args: c.args, names: c.names}
			return builtins[target].fn(s, ci)
		})
	}
	subset(".subset", "[")
	subset(".subset2", "[[")
	// matrix products
	regBuiltin("%*%", func(s *Session, c *callInfo) Value { return s.matprod(c, false, false) })
	regBuiltin("crossprod", func(s *Session, c *callInfo) Value { return s.matprod(c, true, false) })
	regBuiltin("tcrossprod", func(s *Session, c *callInfo) Value { return s.matprod(c, false, true) })

	regBuiltin("proc.time", func(s *Session, c *callInfo) Value {
		res := &Double{V: []float64{0, 0, 0, 0, 0}}
		res.Attr = []Attr{{Name: "names", Value: strVec([]string{"user.self", "sys.self", "elapsed", "user.child", "sys.child"})}, {Name: "class", Value: mkString("proc_time")}}
		return res
	})

	// what is not available: refused with an error that says so, never faked
	for _, n := range unsupportedPrims {
		n := n
		regBuiltin(n, func(s *Session, c *callInfo) Value {
			s.errorf(c.call, "'%s' is not supported by this evaluator (%s)", n, unsupportedWhy(n))
			return nil
		})
	}
	for _, n := range []string{"~", "@", "@<-", "Exec", "Tailcall", "declare", "browser", "standardGeneric", ".Internal", ".isMethodsDispatchOn", "pos.to.env", "tracemem", "untracemem", "retracemem"} {
		n := n
		regSpecial(n, func(s *Session, c *callInfo) Value {
			s.errorf(c.call, "'%s' is not supported by this evaluator (%s)", n, unsupportedWhy(n))
			return nil
		})
	}
}

// unsupportedPrims are the primitives that need a type or a runtime this
// evaluator does not have yet; Unsupported lists them.
var unsupportedPrims = []string{
	"is.single",
	".C", ".Call", ".Call.graphics", ".External", ".External.graphics", ".External2", ".Fortran",
	"lazyLoadDBfetch", ".cache_class", ".primTrace", ".primUntrace", "unCfillPOSIXlt",
}

var unsupportedSpecials = []string{"~", "@", "@<-", "Exec", "Tailcall", "declare", "browser", "standardGeneric", ".Internal", ".isMethodsDispatchOn", "pos.to.env", "tracemem", "untracemem", "retracemem"}

// Unsupported lists the primitives that exist as names but refuse to run:
// they are not counted as implemented.
func Unsupported() []string {
	return append(append([]string(nil), unsupportedPrims...), unsupportedSpecials...)
}

func unsupportedWhy(n string) string {
	switch n {
	case "is.single":
		return "there is no single type"
	case ".C", ".Call", ".Call.graphics", ".External", ".External.graphics", ".External2", ".Fortran", "lazyLoadDBfetch":
		return "foreign code is not loaded; routines are ported to Go instead"
	case "~":
		return "formulas need attributes on calls"
	case "@", "@<-", "standardGeneric", "isS4", ".isMethodsDispatchOn":
		return "S4 is not implemented"
	case ".Internal":
		return "the .Internal functions are rewritten as Go functions, see docs/m3-plan.md"
	}
	return "not implemented yet"
}

func symOrString(v Value) string {
	switch x := v.(type) {
	case *Symbol:
		return x.Name
	case *String:
		if len(x.V) == 1 {
			return x.V[0].S
		}
	}
	return ""
}

func b2i(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func keepShapeAttrs(x Value) []Attr {
	var out []Attr
	for _, a := range attrsOf(x) {
		if a.Name == "names" || a.Name == "dim" || a.Name == "dimnames" {
			out = append(out, a)
		}
	}
	return out
}

func fitNamesEmpty(nm []Str, n int) []Str {
	out := make([]Str, n)
	for i := range out {
		if i < len(nm) {
			out[i] = nm[i]
		}
	}
	return out
}

func anyNAValue(v Value, rec bool) bool {
	switch x := v.(type) {
	case *Logical:
		for _, e := range x.V {
			if e == naInt {
				return true
			}
		}
	case *Integer:
		for _, e := range x.V {
			if e == naInt {
				return true
			}
		}
	case *Double:
		for _, e := range x.V {
			if math.IsNaN(e) {
				return true
			}
		}
	case *String:
		for _, e := range x.V {
			if e.NA {
				return true
			}
		}
	case *Complex:
		for _, e := range x.V {
			if math.IsNaN(real(e)) || math.IsNaN(imag(e)) {
				return true
			}
		}
	case *List:
		for _, e := range x.V {
			if rec {
				if anyNAValue(e, rec) {
					return true
				}
			} else if Length(e) == 1 && isVector(e) {
				if _, isList := e.(*List); !isList && anyNAValue(e, false) {
					return true
				}
			}
		}
	}
	return false
}

// ---- round, signif -----------------------------------------------------------

func (s *Session) roundLike(c *callInfo, sig bool) Value {
	if len(c.args) == 0 {
		s.errorf(c.call, "argument \"x\" is missing, with no default")
	}
	a := s.matchFormals(c, "x", "digits")
	x := a[0]
	switch x.(type) {
	case *Logical, *Integer, *Double, *Complex:
	default:
		s.errorf(c.call, "non-numeric argument to mathematical function")
	}
	dig := 0.0
	if sig {
		dig = 6
	}
	if a[1] != nil {
		switch a[1].(type) {
		case *Logical, *Integer, *Double:
		default:
			s.errorf(c.call, "non-numeric argument to mathematical function")
		}
		dd := asDoubles(a[1])
		if len(dd) == 0 {
			s.errorf(c.call, "invalid second argument of length 0")
		}
		dig = dd[0]
	}
	if cx, ok := x.(*Complex); ok {
		out := make([]complex128, len(cx.V))
		digits := int(dig)
		if math.IsNaN(dig) {
			digits = 0
		}
		for i, z := range cx.V {
			out[i] = roundComplex(z, digits, sig, func(v float64, dg int) float64 {
				if math.IsNaN(v) || math.IsInf(v, 0) {
					return v
				}
				return fround(v, float64(dg))
			})
		}
		return &Complex{V: out, Attr: cx.Attr}
	}
	d := asDoubles(x)
	out := make([]float64, len(d))
	for i, e := range d {
		switch {
		case math.IsNaN(e) || math.IsInf(e, 0):
			out[i] = e
		case math.IsNaN(dig):
			out[i] = math.NaN()
		case sig:
			out[i] = fsignif(e, dig)
		default:
			out[i] = fround(e, dig)
		}
	}
	return &Double{V: out, Attr: attrsOf(x)}
}

func privateRint(x float64) float64 { return math.RoundToEven(x) }

// fround is round(x, dig) with the algorithm of R >= 4.0.0: the candidates
// are the two decimal numbers around x, the nearer one wins, ties go to the
// even digit.
func fround(x, dig float64) float64 {
	const dblDig = 15
	if dig > 308 {
		return x
	}
	if dig < -308 {
		return 0
	}
	if dig == 0 {
		return privateRint(x)
	}
	neg := x < 0
	if neg {
		x = -x
	}
	var r float64
	if dig > 0 {
		l10x := math.Log10(2) * (0.5 + float64(math.Ilogb(x)))
		if l10x+dig > dblDig {
			r = x
		} else {
			p10 := math.Pow(10, dig)
			x10 := p10 * x
			i10 := math.Floor(x10)
			xd := i10 / p10
			xu := math.Ceil(x10) / p10
			du, dd := xu-x, x-xd
			if dd < du || (dd == du && math.Mod(i10, 2) == 0) {
				r = xd
			} else {
				r = xu
			}
		}
	} else {
		p10 := math.Pow(10, -dig)
		r = privateRint(x/p10) * p10
	}
	if neg {
		return -r
	}
	return r
}

func fsignif(x, dig float64) float64 {
	if x == 0 {
		return x
	}
	d := int(dig)
	if d > 22 {
		return x
	}
	if d < 1 {
		d = 1
	}
	neg := x < 0
	if neg {
		x = -x
	}
	e10 := d - 1 - int(math.Floor(math.Log10(x)))
	var r float64
	switch {
	case e10 > 306:
		p := math.Pow(10, float64(e10-306))
		r = privateRint(x*1e306*p) / 1e306 / p
	case e10 >= 0:
		p := math.Pow(10, float64(e10))
		r = privateRint(x*p) / p
	default:
		p := math.Pow(10, float64(-e10))
		r = privateRint(x/p) * p
	}
	if neg {
		return -r
	}
	return r
}

func gammaFn(x float64) float64 {
	if x == 0 || (x < 0 && x == math.Floor(x)) {
		return math.NaN()
	}
	return math.Gamma(x)
}

func lgammaFn(x float64) float64 {
	if x <= 0 && x == math.Floor(x) {
		return math.Inf(1)
	}
	v, _ := math.Lgamma(x)
	return v
}

func sinpi(x float64) float64 {
	if math.IsInf(x, 0) {
		return math.NaN()
	}
	x = math.Mod(x, 2)
	if x <= -1 {
		x += 2
	} else if x > 1 {
		x -= 2
	}
	switch x {
	case 0, 1:
		return 0
	case 0.5:
		return 1
	case -0.5:
		return -1
	}
	return math.Sin(math.Pi * x)
}

func cospi(x float64) float64 {
	if math.IsInf(x, 0) {
		return math.NaN()
	}
	x = math.Mod(math.Abs(x), 2)
	switch x {
	case 0.5, 1.5:
		return 0
	case 1:
		return -1
	case 0:
		return 1
	}
	return math.Cos(math.Pi * x)
}

func tanpi(x float64) float64 {
	if math.IsInf(x, 0) {
		return math.NaN()
	}
	x = math.Mod(x, 1)
	if x <= -0.5 {
		x++
	} else if x > 0.5 {
		x--
	}
	switch x {
	case 0:
		return 0
	case 0.5:
		return math.NaN()
	case 0.25:
		return 1
	case -0.25:
		return -1
	}
	return math.Tan(math.Pi * x)
}

// ---- switch, substitute --------------------------------------------------------

func s3Switch(s *Session, c *callInfo) Value {
	if len(c.raw) == 0 {
		s.errorf(c.call, "'EXPR' is missing")
	}
	sel := s.eval(c.raw[0].Value, c.env)
	alts := c.raw[1:]
	if k := kindOf(sel); !isVector(sel) || Length(sel) != 1 || k == kList || k == kOther {
		s.errorf(c.call, "EXPR must be a length 1 vector")
	}
	if _, isNull := sel.(nullT); isNull {
		s.errorf(c.call, "EXPR must be a length 1 vector")
	}
	if isObject(sel) {
		if cl := classAttr(sel); len(cl) == 1 && cl[0] == "factor" {
			s.warningf(c.call, "EXPR is a \"factor\", treated as integer.\n Consider using '%s' instead.", "switch(as.character( * ), ...)")
		}
	}
	if len(alts) == 0 {
		s.warningf(c.call, "'switch' with no alternatives")
	}
	if st, ok := sel.(*String); ok {
		s.checkSwitchDefaults(c, alts)
		if st.V[0].NA {
			return s.switchDefault(c, alts, "NA")
		}
		key := st.V[0].S
		for i, a := range alts {
			if a.HasName && a.Name == key {
				// fall through empty alternatives
				for j := i; j < len(alts); j++ {
					if !isMissingValue(alts[j].Value) {
						return s.eval(alts[j].Value, c.env)
					}
				}
				s.visible = false
				return Null
			}
		}
		return s.switchDefault(c, alts, key)
	}
	d := asDoubles(sel)
	if math.IsNaN(d[0]) {
		s.visible = false
		return Null
	}
	n := int(d[0])
	if n >= 1 && n <= len(alts) {
		if isMissingValue(alts[n-1].Value) {
			s.errorf(nil, "empty alternative in numeric switch")
		}
		return s.eval(alts[n-1].Value, c.env)
	}
	s.visible = false
	return Null
}

func (s *Session) checkSwitchDefaults(c *callInfo, alts []Arg) {
	var shown []string
	for _, a := range alts {
		if !a.HasName || a.Name == "" {
			shown = append(shown, strings.Join(syntaxDeparse(a.Value), " "))
		}
	}
	if len(shown) > 1 {
		s.errorf(nil, "duplicate 'switch' defaults: '%s' and '%s'", shown[0], shown[1])
	}
}

func (s *Session) switchDefault(c *callInfo, alts []Arg, key string) Value {
	var def *Arg
	nd := 0
	for i := range alts {
		if !alts[i].HasName || alts[i].Name == "" {
			if nd == 0 {
				def = &alts[i]
			}
			nd++
		}
	}
	if nd > 1 {
		var shown []string
		for _, a := range alts {
			if !a.HasName || a.Name == "" {
				shown = append(shown, strings.Join(syntaxDeparse(a.Value), " "))
			}
		}
		s.errorf(c.call, "duplicate 'switch' defaults: '%s' and '%s'", shown[0], shown[1])
	}
	if def != nil && !isMissingValue(def.Value) {
		return s.eval(def.Value, c.env)
	}
	s.visible = false
	return Null
}

// substitute replaces the symbols bound in env by their promise expressions
// (or their values, for ordinary variables), expanding `...`.
func (s *Session) substitute(v Value, env *Env) Value {
	switch x := v.(type) {
	case *Symbol:
		b := env.lookup(x)
		if b == nil {
			return x
		}
		switch bv := b.val.(type) {
		case *Promise:
			return bv.Expr
		case missingT:
			return x
		case *Dots:
			return x
		case *Closure, *Builtin, *Env:
			return x
		default:
			return bv
		}
	case *Lang:
		out := &Lang{Fun: s.substitute(x.Fun, env)}
		for _, a := range x.Args {
			if sym, ok := a.Value.(*Symbol); ok && sym.Name == "..." {
				if b := env.lookup(sym); b != nil {
					if d, ok := b.val.(*Dots); ok {
						for _, da := range d.Args {
							na := Arg{Name: da.Name, HasName: da.Name != ""}
							if p, ok := da.Value.(*Promise); ok {
								na.Value = p.Expr
							} else {
								na.Value = da.Value
							}
							out.Args = append(out.Args, na)
						}
						continue
					}
				}
			}
			if _, isForm := a.Value.(*Formals); isForm {
				out.Args = append(out.Args, a)
				continue
			}
			out.Args = append(out.Args, Arg{Name: a.Name, HasName: a.HasName, Value: s.substitute(a.Value, env)})
		}
		return out
	}
	return v
}

// ---- matrix products -----------------------------------------------------------

func (s *Session) matprod(c *callInfo, cross, tcross bool) Value {
	if len(c.args) < 1 {
		s.errorf(c.call, "argument \"x\" is missing, with no default")
	}
	x := c.args[0]
	var y Value = x
	hasY := len(c.args) > 1
	if hasY {
		y = c.args[1]
		if _, isNull := y.(nullT); isNull {
			y, hasY = x, false
		}
	}
	for _, v := range []Value{x, y} {
		switch v.(type) {
		case *Logical, *Integer, *Double:
		default:
			s.errorf(c.call, "requires numeric/complex matrix/vector arguments")
		}
	}
	xd, yd := dimsOf(x), dimsOf(y)
	// shapes: nrx x ncx and nry x ncy, with vectors promoted as R does
	var nrx, ncx, nry, ncy int
	lx, ly := Length(x), Length(y)
	switch {
	case len(xd) == 2:
		nrx, ncx = xd[0], xd[1]
	case len(xd) > 2:
		s.errorf(c.call, "non-conformable arguments")
	}
	switch {
	case len(yd) == 2:
		nry, ncy = yd[0], yd[1]
	case len(yd) > 2:
		s.errorf(c.call, "non-conformable arguments")
	}
	switch {
	case !cross && !tcross:
		switch {
		case len(xd) != 2 && len(yd) != 2:
			if lx == ly {
				nrx, ncx, nry, ncy = 1, lx, ly, 1
			} else if ly == 1 {
				nrx, ncx, nry, ncy = lx, 1, 1, 1
			} else {
				nrx, ncx, nry, ncy = 1, lx, ly, 1
			}
		case len(xd) != 2:
			if lx == nry {
				nrx, ncx = 1, lx
			} else if nry == 1 {
				nrx, ncx = lx, 1
			} else {
				nrx, ncx = 1, lx
			}
		case len(yd) != 2:
			if ly == ncx {
				nry, ncy = ly, 1
			} else if ncx == 1 {
				nry, ncy = 1, ly
			} else {
				nry, ncy = ly, 1
			}
		}
		if ncx != nry {
			s.errorf(c.call, "non-conformable arguments")
		}
		return s.gemm(c, x, y, nrx, ncx, nry, ncy, false, false, xd, yd)
	case cross:
		if len(xd) != 2 {
			nrx, ncx = lx, 1
		}
		if len(yd) != 2 {
			nry, ncy = ly, 1
		}
		if !hasY {
			nry, ncy = nrx, ncx
		}
		if nrx != nry {
			s.errorf(nil, "non-conformable arguments")
		}
		return s.gemm(c, x, y, nrx, ncx, nry, ncy, true, false, xd, yd)
	}
	if len(xd) != 2 {
		nrx, ncx = lx, 1
	}
	if len(yd) != 2 {
		nry, ncy = ly, 1
	}
	if !hasY {
		nry, ncy = nrx, ncx
	}
	if ncx != ncy {
		s.errorf(c.call, "non-conformable arguments")
	}
	return s.gemm(c, x, y, nrx, ncx, nry, ncy, false, true, xd, yd)
}

// gemm multiplies op(x) by op(y), op being a transpose when asked.
func (s *Session) gemm(c *callInfo, x, y Value, nrx, ncx, nry, ncy int, tx, ty bool, xd, yd []int) Value {
	xv, yv := asDoubles(x), asDoubles(y)
	m, k, n := nrx, ncx, ncy
	if tx {
		m, k = ncx, nrx
	}
	if ty {
		n = nry
	}
	if m > 0 && n > s.opts.MaxVector/m {
		s.errorf(c.call, "vector size too large")
	}
	s.checkLen(c.call, m*n)
	at := func(i, p int) float64 {
		if tx {
			return xv[p+i*nrx]
		}
		return xv[i+p*nrx]
	}
	bt := func(p, j int) float64 {
		if ty {
			return yv[j+p*nry]
		}
		return yv[p+j*nry]
	}
	out := make([]float64, m*n)
	for j := 0; j < n; j++ {
		for i := 0; i < m; i++ {
			var sum float64
			for p := 0; p < k; p++ {
				sum += at(i, p) * bt(p, j)
			}
			out[i+j*m] = sum
		}
	}
	res := &Double{V: out}
	attr := []Attr{{Name: "dim", Value: &Integer{V: []int32{int32(m), int32(n)}}}}
	// dimnames follow the operands' own, for the dimensions that survive
	rowN := dimNamesOf(x, boolInt(tx))
	colN := dimNamesOf(y, boolInt(!ty))
	if rowN != nil || colN != nil {
		var r, cc Value = Null, Null
		if rowN != nil {
			r = rowN
		}
		if colN != nil {
			cc = colN
		}
		attr = append(attr, Attr{Name: "dimnames", Value: &List{V: []Value{r, cc}}})
	}
	res.Attr = attr
	return res
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func dimNamesOf(v Value, d int) *String {
	if dn, ok := getAttr(v, "dimnames").(*List); ok && d < len(dn.V) {
		if nm, ok := dn.V[d].(*String); ok {
			return nm
		}
	}
	return nil
}

func digammaFn(x float64) float64 {
	if math.IsNaN(x) {
		return x
	}
	if x <= 0 && x == math.Floor(x) {
		return math.NaN()
	}
	if x < 0 {
		return digammaFn(1-x) - math.Pi/math.Tan(math.Pi*x)
	}
	r := 0.0
	for x < 10 {
		r -= 1 / x
		x++
	}
	f := 1 / (x * x)
	t := f * (-1.0/12 + f*(1.0/120+f*(-1.0/252+f*(1.0/240+f*(-1.0/132+f*(691.0/32760+f*(-1.0/12)))))))
	return r + math.Log(x) - 0.5/x + t
}

func trigammaFn(x float64) float64 {
	if math.IsNaN(x) {
		return x
	}
	if x <= 0 && x == math.Floor(x) {
		return math.Inf(1)
	}
	if x < 0 {
		sn := math.Sin(math.Pi * x)
		return math.Pi*math.Pi/(sn*sn) - trigammaFn(1-x)
	}
	r := 0.0
	for x < 10 {
		r += 1 / (x * x)
		x++
	}
	f := 1 / (x * x)
	t := 1/x + f/2 + (1/x)*f*(1.0/6+f*(-1.0/30+f*(1.0/42+f*(-1.0/30+f*(5.0/66+f*(-691.0/2730+f*(7.0/6)))))))
	return r + t
}

func init() {
	regBuiltin("digamma", (*Session)(nil).math1("digamma", digammaFn, false))
	regBuiltin("trigamma", (*Session)(nil).math1("trigamma", trigammaFn, false))
	regBuiltin(".class2", func(s *Session, c *callInfo) Value { return strVec(dispatchClass(s.nargs1(c, ".class2"))) })
	regBuiltin("seq.int", func(s *Session, c *callInfo) Value { return builtins["seq"].fn(s, c) })
}

// seqErrCall is the call R names in seq's errors: seq.default, or seq.int
// when the primitive was called.
func seqErrCall(c *callInfo) *Lang {
	if f, ok := c.call.Fun.(*Symbol); ok && f.Name == "seq.int" {
		return c.call
	}
	return &Lang{Fun: Intern("seq.default"), Args: c.call.Args}
}

// log2Fn and log10Fn: Go's math.Log2 and math.Log10 go through Frexp and
// lose several bits near 1 (log2(1.005) is 8 ulp off); the logarithm divided
// by the logarithm of the base is within an ulp, and exact powers are exact.
func log2Fn(x float64) float64 {
	if fr, e := math.Frexp(x); fr == 0.5 {
		return float64(e - 1)
	}
	return math.Log(x) / math.Ln2
}

func log10Fn(x float64) float64 {
	r := math.Log(x) / math.Ln10
	if k := math.Round(r); k == r || (math.Abs(r-k) < 1e-9 && math.Abs(k) <= 22 && math.Pow(10, k) == x) {
		return k
	}
	return r
}
