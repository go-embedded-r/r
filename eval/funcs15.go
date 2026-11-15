package eval

import (
	"math"
	"math/cmplx"
	"strconv"
	"strings"

	"github.com/go-embedded-r/r/syntax"
)

// Slice 15: function objects (formals, body<-, formals<-, as.function) and the
// default methods of the generics that are Go natives with internal dispatch.
// Written from the help pages (formals.Rd, body.Rd, as.function.Rd) and
// experiments with the R 4.6.2 oracle.

func closureFormalsList(f *Closure) Value {
	if len(f.Formals) == 0 {
		return Null
	}
	l := &List{V: make([]Value, len(f.Formals))}
	nm := make([]Str, len(f.Formals))
	for i, fm := range f.Formals {
		nm[i] = Str{S: fm.Name.Name}
		if fm.Default == nil {
			l.V[i] = Missing
		} else {
			l.V[i] = fm.Default
		}
	}
	l.Attr = namesAttr(nm)
	return l
}

// formalsFromList turns a named list (or NULL) into formal arguments; an empty
// symbol is an argument without a default.
func (s *Session) formalsFromList(c *callInfo, v Value) []Formal {
	switch l := v.(type) {
	case nullT:
		return nil
	case *List:
		nm := namesOf(l)
		out := make([]Formal, len(l.V))
		for i, x := range l.V {
			if nm == nil || nm.V[i].S == "" {
				s.errorf(c.call, "invalid formal argument list for \"function\"")
			}
			out[i].Name = Intern(nm.V[i].S)
			if _, miss := x.(missingT); !miss {
				out[i].Default = x
			}
		}
		return out
	}
	s.errorf(c.call, "invalid formal argument list for \"function\"")
	return nil
}

func init() {
	regBuiltin("formals", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "fun", "envir")
		fn := a[0]
		if st, ok := fn.(*String); ok && len(st.V) == 1 {
			if v, found := s.lookupMode(st.V[0].S, c.env, true, "function", true); found {
				fn = v
			}
		}
		if f, ok := fn.(*Closure); ok {
			return closureFormalsList(f)
		}
		return Null
	})
	regBuiltin("formals<-", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "fun", "envir", "value")
		f, ok := a[0].(*Closure)
		if !ok {
			s.errorf(c.call, "invalid 'fun' argument")
		}
		nf := &Closure{Formals: s.formalsFromList(c, a[2]), Body: f.Body, Env: f.Env}
		if a[1] != nil {
			if e, ok := a[1].(*Env); ok {
				nf.Env = e
			}
		}
		if at := attrsOf(f); len(at) > 0 {
			return withAttr(nf, at)
		}
		return nf
	})
	regBuiltin("body<-", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "fun", "envir", "value")
		f, ok := a[0].(*Closure)
		if !ok {
			s.errorf(c.call, "invalid 'fun' argument")
		}
		nf := &Closure{Formals: f.Formals, Body: a[2], Env: f.Env}
		if a[1] != nil {
			if e, ok := a[1].(*Env); ok {
				nf.Env = e
			}
		}
		if a[2] == nil {
			nf.Body = Null
		}
		if at := attrsOf(f); len(at) > 0 {
			return withAttr(nf, at)
		}
		return nf
	})
	regBuiltin("as.function.default", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "envir")
		l, ok := a[0].(*List)
		if !ok {
			s.errorf(c.call, "list argument expected")
		}
		if len(l.V) == 0 {
			s.errorf(c.call, "argument must have length at least 1")
		}
		env := c.env
		if a[1] != nil {
			if e, ok := a[1].(*Env); ok {
				env = e
			}
		}
		n := len(l.V) - 1
		nm := namesOf(l)
		out := make([]Formal, n)
		for i := 0; i < n; i++ {
			if nm == nil || nm.V[i].S == "" {
				s.errorf(c.call, "invalid formal argument list for \"function\"")
			}
			out[i].Name = Intern(nm.V[i].S)
			if _, miss := l.V[i].(missingT); !miss {
				out[i].Default = l.V[i]
			}
		}
		return &Closure{Formals: out, Body: l.V[n], Env: env}
	})

	regBuiltin("sys.parent", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "n")
		n := 1
		if a[0] != nil {
			n = s.intArg(c, a[0], "n")
			if n < 1 {
				s.errorf(c.call, "invalid 'n' value")
			}
		}
		e := s.parentFrame(c.env, n)
		if e == nil || e == s.Global {
			return mkInt(0)
		}
		idx, fr := s.frameIndex(e)
		if fr == nil {
			return mkInt(0)
		}
		return mkInt(int32(idx + 1))
	})

	// the default method of a native that dispatches inside: the same code
	// with the dispatch switched off for the one call
	def := func(name, target string) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			saved := s.noDispatch
			s.noDispatch = true
			defer func() { s.noDispatch = saved }()
			ci := &callInfo{call: c.call, env: c.env, args: c.args, names: c.names, raw: c.raw}
			return builtins[target].fn(s, ci)
		})
	}
	for _, g := range []string{"unique", "duplicated", "anyDuplicated", "rev", "sort", "xtfrm", "range", "mean", "seq", "format"} {
		def(g+".default", g)
	}
}

// as.hexmode and as.octmode are Go natives: a blind R draft of each converged
// on R's own 77-token run of octhex.R (README, "Clean-room process").
func init() {
	conv := func(name, class string, base int) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			a := s.matchFormals(c, "x")
			x := a[0]
			if x == nil {
				s.errorf(c.call, "argument \"x\" is missing, with no default")
			}
			for _, cl := range classAttr(x) {
				if cl == class {
					return x
				}
			}
			fail := func() {
				s.errorf(c.call, "'x' cannot be coerced to class \"%s\"", class)
			}
			switch v := x.(type) {
			case *Double:
				ints := make([]int32, len(v.V))
				for i, f := range v.V {
					switch {
					case math.IsNaN(f):
						ints[i] = naInt
					case f == math.Trunc(f) && math.Abs(f) <= math.MaxInt32:
						ints[i] = int32(f)
					default:
						fail()
					}
				}
				return setAttr(withAttr(&Integer{V: ints}, attrsOf(v)), "class", mkString(class))
			case *Integer:
				return setAttr(x, "class", mkString(class))
			case *String:
				ints := make([]int32, len(v.V))
				for i, e := range v.V {
					if e.NA {
						ints[i] = naInt
						continue
					}
					t := strings.TrimSpace(e.S)
					if base == 16 {
						t = strings.TrimPrefix(strings.TrimPrefix(t, "0x"), "0X")
					}
					n, err := strconv.ParseInt(t, base, 32)
					if err != nil {
						fail()
					}
					ints[i] = int32(n)
				}
				return setAttr(withAttr(&Integer{V: ints}, attrsOf(v)), "class", mkString(class))
			}
			fail()
			return nil
		})
	}
	conv("as.hexmode", "hexmode", 16)
	conv("as.octmode", "octmode", 8)
}

// diff.POSIXt: the lagged differences of times as a difftime whose units are
// chosen from the first differences, as difftime() chooses them.
func init() {
	regBuiltin("diff.POSIXt", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "lag", "differences", "...")
		x := a[0]
		if x == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		lag, diffs := 1, 1
		bad := func() {
			s.errorf(c.call, "'lag' and 'differences' must be integers >= 1")
		}
		for i, p := range []*int{&lag, &diffs} {
			if _, blank := a[i+1].(missingT); a[i+1] == nil || blank {
				continue
			}
			d := asDoubles(coerceTo(a[i+1], kDouble))
			if Length(a[i+1]) != 1 || len(d) != 1 || math.IsNaN(d[0]) || d[0] > math.MaxInt32 {
				bad()
			}
			*p = int(d[0])
		}
		if lag < 1 || diffs < 1 {
			bad()
		}
		isCT := false
		for _, cl := range classAttr(x) {
			if cl == "POSIXct" {
				isCT = true
			}
		}
		if !isCT {
			f := s.findMethod("as.POSIXct", c.env)
			if f == nil {
				s.errorf(c.call, "could not find function \"as.POSIXct\"")
			}
			x = s.callFunction(f, []Value{x})
		}
		d := append([]float64(nil), asDoubles(x)...)
		mk := func(v []float64, units string) Value {
			return &Double{V: v, Attr: []Attr{{Name: "class", Value: mkString("difftime")}, {Name: "units", Value: mkString(units)}}}
		}
		if lag*diffs >= len(d) {
			return mk([]float64{}, "secs")
		}
		step := func(v []float64) []float64 {
			out := make([]float64, len(v)-lag)
			for i := range out {
				out[i] = v[i+lag] - v[i]
			}
			return out
		}
		d = step(d)
		min := math.Inf(1)
		for _, e := range d {
			if !math.IsNaN(e) && math.Abs(e) < min {
				min = math.Abs(e)
			}
		}
		units, div := "secs", 1.0
		switch {
		case math.IsInf(min, 1) || min < 60:
		case min < 3600:
			units, div = "mins", 60
		case min < 86400:
			units, div = "hours", 3600
		default:
			units, div = "days", 86400
		}
		for i := range d {
			d[i] /= div
		}
		for k := 1; k < diffs; k++ {
			d = step(d)
		}
		return mk(d, units)
	})
}

// labels.default and mostattributes<- are Go natives: blind R drafts of both
// converged on R's own runs (62 tokens of labels.R, 66 normalised of attr.R).
func init() {
	regBuiltin("labels.default", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "object", "...")
		x := a[0]
		if x == nil {
			s.errorf(c.call, "argument \"object\" is missing, with no default")
		}
		seq := func(n int) Value {
			out := make([]Str, n)
			for i := range out {
				out[i] = Str{S: strconv.Itoa(i + 1)}
			}
			return &String{V: out}
		}
		if d := dimsOf(x); len(d) > 0 {
			var dn *List
			if l, ok := getAttr(x, "dimnames").(*List); ok {
				dn = l
			}
			out := &List{V: make([]Value, len(d))}
			for i := range d {
				if dn != nil && i < len(dn.V) && Length(dn.V[i]) > 0 {
					out.V[i] = dn.V[i]
				} else {
					out.V[i] = seq(d[i])
				}
			}
			if dn != nil {
				if nm := namesOf(dn); nm != nil {
					out.Attr = namesAttr(nm.V)
				}
			}
			return out
		}
		if nm := namesOf(x); nm != nil && len(nm.V) > 0 {
			return nm
		}
		return seq(Length(x))
	})

	regBuiltin("mostattributes<-", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "obj", "value")
		obj, val := a[0], a[1]
		if obj == nil || val == nil {
			s.errorf(c.call, "argument \"value\" is missing, with no default")
		}
		if Length(val) == 0 {
			return obj
		}
		l, ok := val.(*List)
		if !ok {
			call := c.call
			if c.tmpCall != nil {
				call = c.tmpCall
			}
			s.errorf(call, "'value' must be a list")
		}
		nm := namesOf(l)
		get := func(name string) Value {
			if nm == nil {
				return nil
			}
			for i, n := range nm.V {
				if n.S == name {
					return l.V[i]
				}
			}
			return nil
		}
		names, dim, dimnames := get("names"), get("dim"), get("dimnames")
		obj = withAttr(obj, nil)
		if nm != nil {
			for i, n := range nm.V {
				if n.S == "names" || n.S == "dim" || n.S == "dimnames" || n.S == "" {
					continue
				}
				obj = s.setAttrByName(c, obj, n.S, l.V[i])
			}
		}
		n := Length(obj)
		var dm []int
		if dim != nil {
			prod := 1
			for _, d := range asInts(coerceTo(dim, kInteger)) {
				prod *= int(d)
			}
			if prod == n {
				obj = s.setDim(c, obj, dim)
				dm = dimsOf(obj)
			}
		}
		if dimnames != nil && dm != nil {
			if dl, ok := dimnames.(*List); ok && len(dl.V) == len(dm) {
				good := true
				for i, e := range dl.V {
					if k := Length(e); k != dm[i] && k != 0 {
						good = false
					}
				}
				if good {
					obj = setAttr(obj, "dimnames", dimnames)
				}
			}
		}
		if names != nil && dm == nil && Length(names) == n {
			obj = s.setNames(c, obj, names)
		}
		return obj
	})
}

// the lagged differences of dates and of difftimes
func init() {
	lagDiffs := func(s *Session, c *callInfo, name string, units func(Value) string) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			a := s.matchFormals(c, "x", "lag", "differences", "...")
			x := a[0]
			if x == nil {
				s.errorf(c.call, "argument \"x\" is missing, with no default")
			}
			lag, diffs := 1, 1
			bad := func() { s.errorf(c.call, "'lag' and 'differences' must be integers >= 1") }
			for i, p := range []*int{&lag, &diffs} {
				if _, blank := a[i+1].(missingT); a[i+1] == nil || blank {
					continue
				}
				d := asDoubles(coerceTo(a[i+1], kDouble))
				if Length(a[i+1]) != 1 || len(d) != 1 || math.IsNaN(d[0]) || d[0] > math.MaxInt32 {
					bad()
				}
				*p = int(d[0])
			}
			if lag < 1 || diffs < 1 {
				bad()
			}
			d := append([]float64(nil), asDoubles(x)...)
			mk := func(v []float64) Value {
				return &Double{V: v, Attr: []Attr{{Name: "class", Value: mkString("difftime")}, {Name: "units", Value: mkString(units(x))}}}
			}
			if lag*diffs >= len(d) {
				return mk([]float64{})
			}
			for k := 0; k < diffs; k++ {
				out := make([]float64, len(d)-lag)
				for i := range out {
					out[i] = d[i+lag] - d[i]
				}
				d = out
			}
			return mk(d)
		})
	}
	lagDiffs(nil, nil, "diff.Date", func(Value) string { return "days" })
	lagDiffs(nil, nil, "diff.difftime", func(x Value) string {
		if u, ok := getAttr(x, "units").(*String); ok && len(u.V) == 1 {
			return u.V[0].S
		}
		return "days"
	})
}

// as.array.default is a Go native: a blind R draft converged on a 52-token run
// of R's as.R.
func init() {
	regBuiltin("as.array.default", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "...")
		x := a[0]
		if x == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		if getAttr(x, "dim") != nil {
			return x
		}
		if !isVector(x) {
			s.errorf(c.call, "'dim<-' applied to non-(vector or list)")
		}
		nm := namesOf(x)
		y := setAttr(x, "names", Null)
		y = setAttr(y, "dim", &Integer{V: []int32{int32(Length(x))}})
		if nm != nil && len(nm.V) > 0 {
			y = setAttr(y, "dimnames", &List{V: []Value{nm}})
		}
		return y
	})
}

// machineList is .Machine: the numerical characteristics of the IEEE double
// and the integer on a 64-bit LP64 platform, in the order R lists them.
func machineList() Value {
	names := []string{"double.eps", "double.neg.eps", "double.xmin", "double.xmax", "double.base",
		"double.digits", "double.rounding", "double.guard", "double.ulp.digits", "double.neg.ulp.digits",
		"double.exponent", "double.min.exp", "double.max.exp", "integer.max", "sizeof.long",
		"sizeof.longlong", "sizeof.longdouble", "sizeof.pointer", "sizeof.time_t"}
	vals := []Value{mkDouble(math.Pow(2, -52)), mkDouble(math.Pow(2, -53)), mkDouble(2.2250738585072014e-308),
		mkDouble(math.MaxFloat64), mkInt(2), mkInt(53), mkInt(5), mkInt(0), mkInt(-52), mkInt(-53),
		mkInt(11), mkInt(-1022), mkInt(1024), mkInt(math.MaxInt32), mkInt(8), mkInt(8), mkInt(8), mkInt(8), mkInt(8)}
	nm := make([]Str, len(names))
	for i, n := range names {
		nm[i] = Str{S: n}
	}
	return &List{V: vals, Attr: namesAttr(nm)}
}

// ..format_pval is the engine of format.pval, a Go native because a blind R
// draft converged on an 88-token normalised run of R's format.R; the R wrapper
// only evaluates the arguments.
func init() {
	regBuiltin("..format_pval", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "pv", "digits", "eps", "na.form")
		pv := asDoubles(coerceTo(a[0], kDouble))
		digits := int(asDoubles(coerceTo(a[1], kDouble))[0])
		eps := asDoubles(coerceTo(a[2], kDouble))[0]
		naForm := asStrings(a[3])[0].S
		fmtAll := func(v []float64, d int) []string {
			if len(v) == 0 {
				return nil
			}
			ci := &callInfo{call: c.call, env: c.env, args: []Value{&Double{V: v}, mkInt(int32(d))}, names: []string{"", "digits"}}
			res := builtins["format"].fn(s, ci)
			var out []string
			for _, e := range asStrings(res) {
				out = append(out, e.S)
			}
			return out
		}
		out := make([]Str, len(pv))
		var idx []int // the values that are not NA
		for i, v := range pv {
			if math.IsNaN(v) {
				out[i] = Str{S: naForm}
			} else {
				idx = append(idx, i)
			}
		}
		var big, plain, sci []int
		for _, i := range idx {
			if pv[i] >= eps {
				big = append(big, i)
				v := pv[i]
				floor := -4
				if digits < 2 {
					floor = -3
				}
				if v == 0 || math.Floor(math.Log10(v)) >= float64(floor) {
					plain = append(plain, i)
				} else {
					sci = append(sci, i)
				}
			}
		}
		pick := func(ix []int) []float64 {
			v := make([]float64, len(ix))
			for k, i := range ix {
				v[k] = pv[i]
			}
			return v
		}
		for _, grp := range [][]int{plain, sci} {
			for k, t := range fmtAll(pick(grp), digits) {
				out[grp[k]] = Str{S: t}
			}
		}
		w := 0
		for _, i := range big {
			if n := displayWidth(out[i].S); n > w {
				w = n
			}
		}
		d := digits - 2
		if d < 1 {
			d = 1
		}
		sep := " "
		if len(big) > 0 {
			if d > 1 && d+6 > w {
				d = w - 7
				if d < 1 {
					d = 1
				}
			}
			if d == 1 && w <= 6 {
				sep = ""
			}
		} else if d == 1 {
			sep = ""
		}
		var small string
		for _, i := range idx {
			if pv[i] < eps {
				if small == "" {
					small = "<" + sep + fmtAll([]float64{eps}, d)[0]
				}
				out[i] = Str{S: small}
			}
		}
		return &String{V: out}
	})
}

// as.character.condition and as.character.error are Go natives: blind R drafts
// converged on 43 exact (94 normalised) tokens of R's conditions.R.
func init() {
	text := func(name string, head func(x Value) string) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			a := s.matchFormals(c, "x", "...")
			x, ok := a[0].(*List)
			if !ok {
				s.errorf(c.call, "argument is not a condition")
			}
			msg := ""
			if m := listField(x, "message"); m != nil && Length(m) > 0 {
				if st, ok := coerceTo(m, kString).(*String); ok {
					msg = st.V[0].S
				}
			}
			call := listField(x, "call")
			if call == nil || call == Null {
				return mkString(head(x) + ": " + msg + "\n")
			}
			d := syntax.Deparse(s.toSyntax(call), nil)
			return mkString(head(x) + " in " + d[0] + ": " + msg + "\n")
		})
	}
	text("as.character.condition", func(x Value) string {
		if cl := classAttr(x); len(cl) > 0 {
			return cl[0]
		}
		return ""
	})
	text("as.character.error", func(Value) string { return "Error" })
}

// print.condition, conditionMessage.condition, conditionCall.condition and
// within.data.frame are Go natives: blind R drafts converged on R's conditions.R
// (71 normalised tokens) and eval.R (43 exact tokens).
func init() {
	regBuiltin("conditionMessage.condition", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "c")
		if l, ok := a[0].(*List); ok {
			if v := listField(l, "message"); v != nil {
				return v
			}
		}
		return Null
	})
	regBuiltin("conditionCall.condition", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "c")
		if l, ok := a[0].(*List); ok {
			if v := listField(l, "call"); v != nil {
				return v
			}
		}
		return Null
	})
	regBuiltin("print.condition", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "...")
		if a[0] == nil || len(classAttr(a[0])) == 0 {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		s.printCondition(a[0])
		s.visible = false
		return a[0]
	})
	regSpecial("within.data.frame", func(s *Session, c *callInfo) Value {
		var data, expr Value
		var positional []Value
		for _, a := range c.raw {
			switch {
			case a.HasName && a.Name == "data":
				data = a.Value
			case a.HasName && a.Name == "expr":
				expr = a.Value
			case !a.HasName:
				positional = append(positional, a.Value)
			}
		}
		for _, p := range positional {
			if data == nil {
				data = p
			} else if expr == nil {
				expr = p
			}
		}
		if data == nil {
			s.errorf(c.call, "argument \"data\" is missing, with no default")
		}
		if expr == nil {
			s.errorf(c.call, "argument \"expr\" is missing, with no default")
		}
		var df Value
		if p, ok := data.(*Promise); ok {
			df = s.force(p)
		} else {
			df = s.eval(data, c.env)
		}
		if p, ok := expr.(*Promise); ok {
			expr = p.Expr // the expression as written, not its value
		}
		cols, ok := df.(*List)
		if !ok {
			s.errorf(c.call, "invalid data frame")
		}
		e := NewEnv(c.env)
		nm := namesOf(cols)
		for i, v := range cols.V {
			if nm != nil && nm.V[i].S != "" {
				e.set(Intern(nm.V[i].S), v)
			}
		}
		s.eval(expr, e)
		var newNames []Str
		var newVals []Value
		have := map[string]bool{}
		for i := len(e.order) - 1; i >= 0; i-- {
			sym := e.order[i]
			b := e.lookup(sym)
			if b == nil || have[sym.Name] {
				continue
			}
			v := b.val
			if p, ok := v.(*Promise); ok {
				v = s.force(p)
			}
			if _, isNull := v.(nullT); isNull {
				continue
			}
			have[sym.Name] = true
			newNames = append(newNames, Str{S: sym.Name})
			newVals = append(newVals, v)
		}
		assign := builtins["[<-.data.frame"].fn
		out := df
		if len(newNames) > 0 {
			out = assign(s, &callInfo{call: c.call, env: c.env, args: []Value{out, &String{V: newNames}, &List{V: newVals, Attr: namesAttr(newNames)}}, names: []string{"", "", "value"}})
		}
		var gone []Str
		if nm != nil {
			for _, n := range nm.V {
				if !have[n.S] {
					gone = append(gone, n)
				}
			}
		}
		out = assign(s, &callInfo{call: c.call, env: c.env, args: []Value{out, &String{V: gone}, Null}, names: []string{"", "", "value"}})
		return out
	})
}

// callNamed calls the R function fname with evaluated, possibly named, arguments.
func (s *Session) callNamed(c *callInfo, fname string, args []Value, names []string) Value {
	f := s.findMethod(fname, c.env)
	if f == nil {
		s.errorf(c.call, "could not find function \"%s\"", fname)
	}
	return s.callMethod(f, fname, &Lang{Fun: Intern(fname)}, argsOf(args, names), c.env, nil)
}

// all.equal.default and all.equal.factor are Go natives: blind R drafts
// converged on 105 normalised tokens (default) and 41 exact tokens (factor) of
// R's all.equal.R. They choose a method and leave the comparing to the R ones.
func init() {
	regBuiltin("all.equal.default", func(s *Session, c *callInfo) Value {
		var target, current Value
		checkClass := true
		var rest []Value
		var restNames []string
		pos := 0
		for i, v := range c.args {
			switch {
			case c.names[i] == "check.class":
				checkClass = s.truthy(c, v, "check.class")
			case c.names[i] == "target" || (c.names[i] == "" && pos == 0 && target == nil):
				target = v
				if c.names[i] == "" {
					pos++
				}
			case c.names[i] == "current" || (c.names[i] == "" && current == nil):
				current = v
			default:
				rest = append(rest, v)
				restNames = append(restNames, c.names[i])
			}
		}
		if target == nil || current == nil {
			s.errorf(c.call, "argument \"current\" is missing, with no default")
		}
		forward := func(fname string) Value {
			return s.callNamed(c, fname, append([]Value{target, current}, rest...), append([]string{"", ""}, restNames...))
		}
		isFn := false
		switch target.(type) {
		case *Closure, *Builtin:
			isFn = true
		}
		_, isLang := target.(*Lang)
		_, isSym := target.(*Symbol)
		if isLang || isSym || isFn {
			return forward("all.equal.language")
		}
		_, e1 := target.(*Env)
		_, e2 := current.(*Env)
		if e1 || e2 {
			return forward("all.equal.environment")
		}
		if l, ok := target.(*List); ok && l != nil {
			return forward("all.equal.list")
		}
		switch target.(type) {
		case *Integer, *Double, *Complex:
			return forward("all.equal.numeric")
		case *String:
			return forward("all.equal.character")
		case *Logical, *Raw:
			return forward("all.equal.raw")
		}
		if checkClass {
			tk := s.callNamed(c, "data.class", []Value{target}, []string{""})
			ck := s.callNamed(c, "data.class", []Value{current}, []string{""})
			if a, b := asStrings(tk)[0].S, asStrings(ck)[0].S; a != b {
				return mkString("target is " + a + ", current is " + b)
			}
		}
		return &Logical{V: []int32{1}}
	})

	regBuiltin("all.equal.factor", func(s *Session, c *callInfo) Value {
		var target, current Value
		check := true
		var rest []Value
		var restNames []string
		for i, v := range c.args {
			switch {
			case c.names[i] == "check.attributes":
				check = s.truthy(c, v, "check.attributes")
			case target == nil && (c.names[i] == "" || c.names[i] == "target"):
				target = v
			case current == nil && (c.names[i] == "" || c.names[i] == "current"):
				current = v
			default:
				rest = append(rest, v)
				restNames = append(restNames, c.names[i])
			}
		}
		if target == nil || current == nil {
			s.errorf(c.call, "argument \"current\" is missing, with no default")
		}
		if !inheritsClass(target, "factor") {
			return mkString("'target' is not a factor")
		}
		if !inheritsClass(current, "factor") {
			return mkString("'current' is not a factor")
		}
		var notes []Str
		add := func(v Value) {
			if st, ok := v.(*String); ok {
				notes = append(notes, st.V...)
			}
		}
		if check {
			add(s.callNamed(c, "attr.all.equal", append([]Value{target, current}, rest...), append([]string{"", ""}, restNames...)))
		}
		tc := s.callNamed(c, "as.character", []Value{target}, []string{""})
		cc := s.callNamed(c, "as.character", []Value{current}, []string{""})
		args := append([]Value{tc, cc, mkLogical(check)}, rest...)
		names := append([]string{"", "", "check.attributes"}, restNames...)
		add(s.callNamed(c, "all.equal", args, names))
		if len(notes) == 0 {
			return &Logical{V: []int32{1}}
		}
		return &String{V: notes}
	})
}

// attr.all.equal and all.equal.list are Go natives too: blind R drafts
// converged on 105 normalised tokens (attr.all.equal) and 41 exact tokens
// (all.equal.list) of R's all.equal.R.
func init() {
	messages := func(v Value) []Str {
		if st, ok := v.(*String); ok {
			return st.V
		}
		return nil
	}
	isTrue := func(v Value) bool {
		l, ok := v.(*Logical)
		return ok && len(l.V) == 1 && l.V[0] == 1 && len(attrsOf(v)) == 0
	}
	stripLengths := func(notes []Str) []Str {
		var out []Str
		for _, n := range notes {
			if !strings.Contains(n.S, "Lengths") {
				out = append(out, n)
			}
		}
		return out
	}
	// parse reads (target, current, named options, the rest) from c's arguments
	parse := func(s *Session, c *callInfo, opts map[string]*Value) (target, current Value, rest []Value, restNames []string) {
		for i, v := range c.args {
			if p, ok := opts[c.names[i]]; ok && c.names[i] != "" {
				*p = v
				continue
			}
			switch {
			case target == nil && (c.names[i] == "" || c.names[i] == "target"):
				target = v
			case current == nil && (c.names[i] == "" || c.names[i] == "current"):
				current = v
			default:
				rest = append(rest, v)
				restNames = append(restNames, c.names[i])
			}
		}
		if target == nil || current == nil {
			s.errorf(c.call, "argument \"current\" is missing, with no default")
		}
		return
	}
	nameOfMode := func(s *Session, c *callInfo, v Value) string {
		return asStrings(s.callNamed(c, "mode", []Value{v}, []string{""}))[0].S
	}
	attrList := func(s *Session, c *callInfo, v Value) Value {
		l, _ := s.callNamed(c, "attributes", []Value{v}, []string{""}).(*List)
		if l == nil {
			return Null
		}
		out := &List{}
		var nm []Str
		if names := namesOf(l); names != nil {
			for i, n := range names.V {
				if n.S != "names" {
					out.V = append(out.V, l.V[i])
					nm = append(nm, n)
				}
			}
		}
		if len(out.V) > 0 {
			out.Attr = namesAttr(nm)
		}
		return out
	}
	sameStrings := func(a, b *String) bool {
		if len(a.V) != len(b.V) {
			return false
		}
		for i := range a.V {
			if a.V[i] != b.V[i] {
				return false
			}
		}
		return true
	}

	var attrEqual func(s *Session, c *callInfo, target, current Value, rest []Value, restNames []string, checkAttr, checkNames bool) []Str
	attrEqual = func(s *Session, c *callInfo, target, current Value, rest []Value, restNames []string, checkAttr, checkNames bool) []Str {
		var notes []Str
		tm, cm := nameOfMode(s, c, target), nameOfMode(s, c, current)
		if tm != cm {
			notes = append(notes, Str{S: "Modes: " + tm + ", " + cm})
		}
		if lt, lc := Length(target), Length(current); lt != lc {
			notes = append(notes, Str{S: "Lengths: " + strconv.Itoa(lt) + ", " + strconv.Itoa(lc)})
		}
		ta, ca := attrList(s, c, target), attrList(s, c, current)
		tn, cn := namesOf(target), namesOf(current)
		if checkNames {
			switch {
			case tn != nil && cn != nil && len(tn.V) > 0 && len(cn.V) > 0:
				if !sameStrings(tn, cn) {
					why := s.callNamed(c, "all.equal.character", []Value{tn, cn, mkLogical(checkAttr)}, []string{"", "", "check.attributes"})
					for _, m := range messages(why) {
						notes = append(notes, Str{S: "Names: " + m.S})
					}
				}
			case tn != nil && len(tn.V) > 0:
				notes = append(notes, Str{S: "names for target but not for current"})
			case cn != nil && len(cn.V) > 0:
				notes = append(notes, Str{S: "names for current but not for target"})
			}
		}
		if checkAttr && (Length(ta) > 0 || Length(ca) > 0) {
			sorted := func(v Value) Value {
				l, ok := v.(*List)
				if !ok || len(l.V) == 0 {
					return v
				}
				ord := asInts(coerceTo(s.callNamed(c, "order", []Value{namesOf(l)}, []string{""}), kInteger))
				out := &List{V: make([]Value, len(ord))}
				nm := make([]Str, len(ord))
				for k, o := range ord {
					out.V[k] = l.V[o-1]
					nm[k] = namesOf(l).V[o-1]
				}
				out.Attr = namesAttr(nm)
				return out
			}
			args := append([]Value{sorted(ta), sorted(ca)}, rest...)
			names := append([]string{"", ""}, restNames...)
			args = append(args, mkLogical(checkAttr))
			names = append(names, "check.attributes")
			why := s.callNamed(c, "all.equal", args, names)
			if !isTrue(why) {
				for _, m := range messages(why) {
					notes = append(notes, Str{S: "Attributes: < " + m.S + " >"})
				}
			}
		}
		return notes
	}

	regBuiltin("attr.all.equal", func(s *Session, c *callInfo) Value {
		var ca, cnm Value
		target, current, rest, restNames := parse(s, c, map[string]*Value{"check.attributes": &ca, "check.names": &cnm})
		checkAttr, checkNames := true, true
		for _, o := range []struct {
			v    Value
			name string
			dst  *bool
		}{{ca, "check.attributes", &checkAttr}, {cnm, "check.names", &checkNames}} {
			if o.v == nil {
				continue
			}
			l, ok := o.v.(*Logical)
			if !ok {
				s.errorf(c.call, "'%s' must be logical", o.name)
			}
			*o.dst = len(l.V) > 0 && l.V[0] == 1
		}
		notes := attrEqual(s, c, target, current, rest, restNames, checkAttr, checkNames)
		if len(notes) == 0 {
			return Null
		}
		return &String{V: notes}
	})

	regBuiltin("all.equal.list", func(s *Session, c *callInfo) Value {
		var ca, un Value
		target, current, rest, restNames := parse(s, c, map[string]*Value{"check.attributes": &ca, "use.names": &un})
		checkAttr := ca == nil || s.truthy(c, ca, "check.attributes")
		useNames := un == nil || s.truthy(c, un, "use.names")
		var notes []Str
		if checkAttr {
			notes = attrEqual(s, c, target, current, rest, restNames, true, true)
		}
		listLike := func(v Value) bool {
			if _, ok := v.(*List); ok {
				return true
			}
			if !isVector(v) {
				return false
			}
			for _, a := range attrsOf(v) {
				if a.Name != "names" && a.Name != "class" {
					return false
				}
			}
			return true
		}
		target, current = setAttr(target, "class", Null), setAttr(current, "class", Null)
		finish := func(extra ...string) Value {
			for _, e := range extra {
				notes = append(notes, Str{S: e})
			}
			if len(notes) == 0 {
				return &Logical{V: []int32{1}}
			}
			return &String{V: notes}
		}
		if !listLike(target) {
			return finish("target is not list-like")
		}
		if !listLike(current) {
			return finish("current is not list-like")
		}
		lt, lc := Length(target), Length(current)
		n := lt
		if lt != lc {
			notes = stripLengths(notes)
			if lc < n {
				n = lc
			}
			notes = append(notes, Str{S: "Length mismatch: comparison on first " + strconv.Itoa(n) + " components"})
		}
		tn, cn := namesOf(target), namesOf(current)
		for i := 1; i <= n; i++ {
			idx := mkInt(int32(i))
			a := s.callNamed(c, "[[", []Value{target, idx}, []string{"", ""})
			b := s.callNamed(c, "[[", []Value{current, idx}, []string{"", ""})
			args := append([]Value{a, b, mkLogical(checkAttr), mkLogical(useNames)}, rest...)
			names := append([]string{"", "", "check.attributes", "use.names"}, restNames...)
			why := s.callNamed(c, "all.equal", args, names)
			label := strconv.Itoa(i)
			// a component is named when both sides call it by the same name
			if useNames && tn != nil && cn != nil && i <= len(tn.V) && i <= len(cn.V) && tn.V[i-1].S != "" && tn.V[i-1] == cn.V[i-1] {
				label = "“" + tn.V[i-1].S + "”"
			}
			for _, m := range messages(why) {
				notes = append(notes, Str{S: "Component " + label + ": " + m.S})
			}
		}
		return finish()
	})
}

// all.equal.character, all.equal.raw and all.equal.numeric complete the family
// in Go: the R drafts of the first two converged on 79 normalised tokens of
// R's all.equal.R, and the R file as a whole passed the exact-containment limit.
func init() {
	sameKindMsg := func(s *Session, c *callInfo, a, b Value) string {
		x := asStrings(s.callNamed(c, "data.class", []Value{a}, []string{""}))[0].S
		y := asStrings(s.callNamed(c, "data.class", []Value{b}, []string{""}))[0].S
		return "target is " + x + ", current is " + y
	}
	common := func(s *Session, c *callInfo, fname string) (target, current Value, rest []Value, restNames []string, checkAttr bool) {
		checkAttr = true
		for i, v := range c.args {
			switch {
			case c.names[i] == "check.attributes":
				checkAttr = s.truthy(c, v, "check.attributes")
			case target == nil && (c.names[i] == "" || c.names[i] == "target"):
				target = v
			case current == nil && (c.names[i] == "" || c.names[i] == "current"):
				current = v
			default:
				rest = append(rest, v)
				restNames = append(restNames, c.names[i])
			}
		}
		if target == nil || current == nil {
			s.errorf(c.call, "argument \"current\" is missing, with no default")
		}
		return
	}
	attrNotes := func(s *Session, c *callInfo, target, current Value, rest []Value, restNames []string) []Str {
		v := s.callNamed(c, "attr.all.equal", append([]Value{target, current}, rest...), append([]string{"", ""}, restNames...))
		if st, ok := v.(*String); ok {
			return st.V
		}
		return nil
	}
	verdict := func(notes []Str) Value {
		if len(notes) == 0 {
			return &Logical{V: []int32{1}}
		}
		return &String{V: notes}
	}
	dropLengths := func(notes []Str) []Str {
		var out []Str
		for _, n := range notes {
			if !strings.Contains(n.S, "Lengths") {
				out = append(out, n)
			}
		}
		return out
	}
	mismatches := func(n int, what string) Str {
		if n > 1 {
			return Str{S: strconv.Itoa(n) + " " + what + " mismatches"}
		}
		return Str{S: strconv.Itoa(n) + " " + what + " mismatch"}
	}
	naPattern := func(tna, cna []bool) (Str, bool) {
		differ := false
		nt, nc := 0, 0
		for i := range tna {
			if tna[i] != cna[i] {
				differ = true
			}
			if tna[i] {
				nt++
			}
			if cna[i] {
				nc++
			}
		}
		return Str{S: "'is.NA' value mismatch: " + strconv.Itoa(nc) + " in current " + strconv.Itoa(nt) + " in target"}, differ
	}

	regBuiltin("all.equal.character", func(s *Session, c *callInfo) Value {
		target, current, rest, restNames, checkAttr := common(s, c, "all.equal.character")
		var notes []Str
		if checkAttr {
			notes = attrNotes(s, c, target, current, rest, restNames)
		}
		t, ok1 := target.(*String)
		cu, ok2 := current.(*String)
		if !ok1 || !ok2 {
			return verdict(append(notes, Str{S: sameKindMsg(s, c, target, current)}))
		}
		tv, cv := t.V, cu.V
		if len(tv) != len(cv) {
			n := len(tv)
			if len(cv) < n {
				n = len(cv)
			}
			notes = dropLengths(notes)
			notes = append(notes, Str{S: "Lengths (" + strconv.Itoa(len(tv)) + ", " + strconv.Itoa(len(cv)) + ") differ (string compare on first " + strconv.Itoa(n) + ")"})
			tv, cv = tv[:n], cv[:n]
		}
		tna, cna := make([]bool, len(tv)), make([]bool, len(cv))
		for i := range tv {
			tna[i], cna[i] = tv[i].NA, cv[i].NA
		}
		if m, differ := naPattern(tna, cna); differ {
			return verdict(append(notes, m))
		}
		bad := 0
		for i := range tv {
			if !tna[i] && tv[i].S != cv[i].S {
				bad++
			}
		}
		if bad > 0 {
			notes = append(notes, mismatches(bad, "string"))
		}
		return verdict(notes)
	})

	regBuiltin("all.equal.raw", func(s *Session, c *callInfo) Value {
		target, current, rest, restNames, checkAttr := common(s, c, "all.equal.raw")
		var notes []Str
		if checkAttr {
			notes = attrNotes(s, c, target, current, rest, restNames)
		}
		classDiffers := false
		if a, b := asStrings(s.callNamed(c, "data.class", []Value{target}, []string{""}))[0].S, asStrings(s.callNamed(c, "data.class", []Value{current}, []string{""}))[0].S; a != b {
			notes = append(notes, Str{S: "target is " + a + ", current is " + b})
			classDiffers = true
		}
		if Length(target) != Length(current) {
			if classDiffers {
				return verdict(notes) // the lengths are already in the notes
			}
			notes = dropLengths(notes)
			notes = append(notes, Str{S: "Lengths (" + strconv.Itoa(Length(target)) + ", " + strconv.Itoa(Length(current)) + ") differ"})
			return verdict(notes)
		}
		tv := asDoubles(coerceTo(target, kDouble))
		cv := asDoubles(coerceTo(current, kDouble))
		if len(tv) != Length(target) || len(cv) != Length(current) {
			// what cannot be read as numbers is compared by its text
			ts, cs := asStrings(target), asStrings(current)
			tv, cv = make([]float64, len(ts)), make([]float64, len(cs))
			for i := range ts {
				if ts[i].NA {
					tv[i] = math.NaN()
				} else if i < len(cs) && ts[i].S != cs[i].S {
					tv[i], cv[i] = 0, 1
				}
			}
			for i := range cs {
				if cs[i].NA {
					cv[i] = math.NaN()
				}
			}
		}
		if len(tv) != len(cv) {
			return verdict(notes)
		}
		tna, cna := make([]bool, len(tv)), make([]bool, len(cv))
		for i := range tv {
			tna[i], cna[i] = math.IsNaN(tv[i]), math.IsNaN(cv[i])
		}
		if m, differ := naPattern(tna, cna); differ {
			return verdict(append(notes, m))
		}
		bad := 0
		for i := range tv {
			if !tna[i] && tv[i] != cv[i] {
				bad++
			}
		}
		if bad > 0 {
			notes = append(notes, mismatches(bad, "element"))
		}
		return verdict(notes)
	})

	regBuiltin("all.equal.numeric", func(s *Session, c *callInfo) Value {
		var tolV, scaleV, countV, fmtV, ckAttr, ckClass Value
		opts := map[string]*Value{"tolerance": &tolV, "scale": &scaleV, "countEQ": &countV, "formatFUN": &fmtV,
			"check.attributes": &ckAttr, "check.class": &ckClass, "giveErr": new(Value)}
		var target, current Value
		var rest []Value
		var restNames []string
		for i, v := range c.args {
			if p, ok := opts[c.names[i]]; ok && c.names[i] != "" {
				*p = v
				continue
			}
			switch {
			case target == nil && (c.names[i] == "" || c.names[i] == "target"):
				target = v
			case current == nil && (c.names[i] == "" || c.names[i] == "current"):
				current = v
			default:
				rest = append(rest, v)
				restNames = append(restNames, c.names[i])
			}
		}
		if target == nil || current == nil {
			s.errorf(c.call, "argument \"current\" is missing, with no default")
		}
		tol := math.Sqrt(math.Pow(2, -52))
		if tolV != nil {
			switch tolV.(type) {
			case *Double, *Integer:
				tol = asDoubles(coerceTo(tolV, kDouble))[0]
			default:
				s.errorf(c.call, "'tolerance' should be numeric")
			}
		}
		var scale []float64
		if scaleV != nil {
			if _, isNull := scaleV.(nullT); !isNull {
				switch scaleV.(type) {
				case *Double, *Integer:
					scale = asDoubles(coerceTo(scaleV, kDouble))
				default:
					s.errorf(c.call, "'scale' should be numeric or NULL")
				}
			}
		}
		countEQ := countV != nil && s.truthy(c, countV, "countEQ")
		checkAttr := true
		if ckAttr != nil {
			l, ok := ckAttr.(*Logical)
			if !ok {
				s.errorf(c.call, "'check.attributes' must be logical")
			}
			checkAttr = len(l.V) > 0 && l.V[0] == 1
		}
		checkClass := ckClass == nil || s.truthy(c, ckClass, "check.class")
		var notes []Str
		if checkAttr {
			extra := append([]Value{mkDouble(tol)}, rest...)
			extraNames := append([]string{"tolerance"}, restNames...)
			if scaleV != nil {
				extra = append(extra, scaleV)
				extraNames = append(extraNames, "scale")
			}
			notes = attrNotes(s, c, target, current, extra, extraNames)
		}
		classDiffers := false
		if checkClass {
			if a, b := asStrings(s.callNamed(c, "data.class", []Value{target}, []string{""}))[0].S, asStrings(s.callNamed(c, "data.class", []Value{current}, []string{""}))[0].S; a != b {
				notes = append(notes, Str{S: "target is " + a + ", current is " + b})
				classDiffers = true
			}
		}
		_, cplx := target.(*Complex)
		switch current.(type) {
		case *Integer, *Double, *Complex:
		default:
			return verdict(notes)
		}
		if lt, lc := Length(target), Length(current); lt != lc {
			if classDiffers {
				return verdict(notes) // the lengths are already in the notes
			}
			notes = dropLengths(notes)
			kind := "Numeric"
			if cplx {
				kind = "Complex"
			}
			notes = append(notes, Str{S: kind + ": lengths (" + strconv.Itoa(lt) + ", " + strconv.Itoa(lc) + ") differ"})
			return verdict(notes)
		}
		n := Length(target)
		tz := make([]complex128, n)
		cz := make([]complex128, n)
		if cplx {
			copy(tz, target.(*Complex).V)
			if cc, ok := current.(*Complex); ok {
				copy(cz, cc.V)
			} else {
				for i, f := range asDoubles(coerceTo(current, kDouble)) {
					cz[i] = complex(f, 0)
				}
			}
		} else {
			for i, f := range asDoubles(coerceTo(target, kDouble)) {
				tz[i] = complex(f, 0)
			}
			for i, f := range asDoubles(coerceTo(current, kDouble)) {
				cz[i] = complex(f, 0)
			}
		}
		isNA := func(z complex128) bool { return math.IsNaN(real(z)) || math.IsNaN(imag(z)) }
		tna, cna := make([]bool, n), make([]bool, n)
		for i := range tz {
			tna[i], cna[i] = isNA(tz[i]), isNA(cz[i])
		}
		if m, differ := naPattern(tna, cna); differ {
			return verdict(append(notes, m))
		}
		same := make([]bool, n)
		allSame, nSame, nDiff := true, 0, 0
		var sumSame, sumT, sumD float64
		for i := range tz {
			same[i] = tna[i] || tz[i] == cz[i]
			if same[i] {
				nSame++
				if !tna[i] {
					sumSame += cmplx.Abs(tz[i])
				}
			} else {
				allSame = false
				nDiff++
				sumT += cmplx.Abs(tz[i])
				sumD += cmplx.Abs(tz[i] - cz[i])
			}
		}
		if allSame {
			return verdict(notes)
		}
		total := float64(nDiff)
		baseline := 0.0
		if countEQ && nSame > 0 {
			// the equal elements count as their mean size; a missing one makes it unknown
			baseline = sumSame / float64(nSame) * float64(nDiff)
			for i := range same {
				if same[i] && tna[i] {
					baseline = math.NaN()
				}
			}
		}
		err := sumD / total
		what := "absolute"
		if scale == nil {
			size := (baseline + sumT) / total
			if !math.IsInf(size, 0) && !math.IsNaN(size) && size > tol {
				err /= size
				what = "relative"
			}
		} else {
			for _, sc := range scale {
				if !(sc > 0) {
					s.errorf(c.call, "all(scale > 0) is not TRUE")
				}
			}
			err /= scale[0]
			if math.Abs(scale[0]-1) >= 1e-07 {
				what = "scaled"
			}
		}
		if cplx {
			what += " Mod"
		}
		if math.IsNaN(err) || err > tol {
			var shown string
			if f, ok := fmtV.(*Closure); ok {
				shown = asStrings(s.callFunction(f, []Value{mkDouble(err), mkString(what)}))[0].S
			} else {
				shown = asStrings(s.callNamed(c, "format", []Value{mkDouble(err)}, []string{""}))[0].S
			}
			notes = append(notes, Str{S: "Mean " + what + " difference: " + shown})
		}
		return verdict(notes)
	})
}

// all.equal.environment is a Go native: its R draft converged on 38 exact
// tokens of R's all.equal.R.
func init() {
	regBuiltin("all.equal.environment", func(s *Session, c *callInfo) Value {
		var target, current Value
		allNames := true
		var rest []Value
		var restNames []string
		for i, v := range c.args {
			switch {
			case c.names[i] == "all.names":
				allNames = s.truthy(c, v, "all.names")
			case c.names[i] == "evaluate":
			case target == nil && (c.names[i] == "" || c.names[i] == "target"):
				target = v
			case current == nil && (c.names[i] == "" || c.names[i] == "current"):
				current = v
			default:
				rest = append(rest, v)
				restNames = append(restNames, c.names[i])
			}
		}
		if target == nil || current == nil {
			s.errorf(c.call, "argument \"current\" is missing, with no default")
		}
		te, ok := target.(*Env)
		if !ok {
			return mkString("'target' is not an environment")
		}
		ce, ok := current.(*Env)
		if !ok {
			return mkString("'current' is not an environment")
		}
		if te == ce {
			return &Logical{V: []int32{1}}
		}
		list := func(e *Env) Value {
			return s.callNamed(c, "as.list.environment", []Value{e, mkLogical(allNames), mkLogical(true)}, []string{"", "all.names", "sorted"})
		}
		return s.callNamed(c, "all.equal.list", append([]Value{list(te), list(ce)}, rest...), append([]string{"", ""}, restNames...))
	})
}

// all.equal.language and all.equal.function are Go natives as well: their R
// drafts converged on 33 exact tokens of R's all.equal.R.
func init() {
	lang := func(s *Session, c *callInfo, target, current Value) Value {
		mode := func(v Value) string { return asStrings(s.callNamed(c, "mode", []Value{v}, []string{""}))[0].S }
		mt, mc := mode(target), mode(current)
		if mt == "expression" && mc == "expression" {
			return s.callNamed(c, "all.equal.list", []Value{target, current}, []string{"", ""})
		}
		var notes []Str
		if mt != mc {
			notes = append(notes, Str{S: "Modes of target, current: " + mt + ", " + mc})
		}
		dt := s.callNamed(c, "deparse", []Value{target}, []string{""})
		dc := s.callNamed(c, "deparse", []Value{current}, []string{""})
		a, b := asStrings(dt), asStrings(dc)
		same := len(a) == len(b)
		for i := 0; same && i < len(a); i++ {
			same = a[i] == b[i]
		}
		if !same {
			notes = append(notes, Str{S: "target, current do not match when deparsed"})
		}
		if len(notes) == 0 {
			return &Logical{V: []int32{1}}
		}
		return &String{V: notes}
	}
	pair := func(s *Session, c *callInfo) (Value, Value, bool) {
		var target, current Value
		checkEnv := true
		for i, v := range c.args {
			switch {
			case c.names[i] == "check.environment":
				checkEnv = s.truthy(c, v, "check.environment")
			case target == nil && (c.names[i] == "" || c.names[i] == "target"):
				target = v
			case current == nil && (c.names[i] == "" || c.names[i] == "current"):
				current = v
			}
		}
		if target == nil || current == nil {
			s.errorf(c.call, "argument \"current\" is missing, with no default")
		}
		return target, current, checkEnv
	}
	regBuiltin("all.equal.language", func(s *Session, c *callInfo) Value {
		target, current, _ := pair(s, c)
		return lang(s, c, target, current)
	})
	regBuiltin("all.equal.function", func(s *Session, c *callInfo) Value {
		target, current, checkEnv := pair(s, c)
		var notes []Str
		if r, ok := lang(s, c, target, current).(*String); ok {
			notes = append(notes, r.V...)
		}
		if checkEnv {
			env := func(v Value) Value { return s.callNamed(c, "environment", []Value{v}, []string{""}) }
			if e, ok := s.callNamed(c, "all.equal.environment", []Value{env(target), env(current)}, []string{"", ""}).(*String); ok {
				for _, m := range e.V {
					notes = append(notes, Str{S: "Environments: " + m.S})
				}
			}
		}
		if len(notes) == 0 {
			return &Logical{V: []int32{1}}
		}
		return &String{V: notes}
	})
}
