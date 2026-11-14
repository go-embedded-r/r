package eval

import (
	"math"
	"strconv"
	"strings"
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
		if a[1] != nil {
			if Length(a[1]) != 1 {
				bad()
			}
			lag = int(asDoubles(coerceTo(a[1], kDouble))[0])
		}
		if a[2] != nil {
			if Length(a[2]) != 1 {
				bad()
			}
			diffs = int(asDoubles(coerceTo(a[2], kDouble))[0])
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
				if a[i+1] != nil {
					if Length(a[i+1]) != 1 {
						bad()
					}
					*p = int(asDoubles(coerceTo(a[i+1], kDouble))[0])
				}
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
