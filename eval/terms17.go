package eval

import (
	"sort"
	"strings"

	"github.com/go-embedded-r/r/syntax"
)

// terms.formula (slice 17), written from ?terms.formula, ?terms.object and
// ?formula and from what the oracle answers. A model formula is read as a set
// of terms, each a set of variables; the variables are the sub-expressions
// that are not formula operators, numbered in the order they are met.

const (
	maxTerms     = 1 << 16 // a formula that expands to more terms is refused
	maxFactorLen = 1 << 24 // cells of the factors matrix
)

type bitset []uint64

func (b bitset) key() string {
	var sb strings.Builder
	n := len(b)
	for n > 0 && b[n-1] == 0 {
		n--
	}
	for _, w := range b[:n] {
		for i := 0; i < 8; i++ {
			sb.WriteByte(byte(w >> (8 * i)))
		}
	}
	return sb.String()
}

func (b bitset) has(i int) bool { return i/64 < len(b) && b[i/64]>>(uint(i)%64)&1 == 1 }

func (b bitset) union(o bitset) bitset {
	n := len(b)
	if len(o) > n {
		n = len(o)
	}
	out := make(bitset, n)
	copy(out, b)
	for i, w := range o {
		out[i] |= w
	}
	return out
}

func (b bitset) subsetOf(o bitset) bool {
	for i, w := range b {
		var ow uint64
		if i < len(o) {
			ow = o[i]
		}
		if w&^ow != 0 {
			return false
		}
	}
	return true
}

func (b bitset) count() int {
	n := 0
	for i := 0; i < len(b)*64; i++ {
		if b.has(i) {
			n++
		}
	}
	return n
}

func (b bitset) with(i int) bitset {
	n := i/64 + 1
	if len(b) > n {
		n = len(b)
	}
	out := make(bitset, n)
	copy(out, b)
	out[i/64] |= 1 << (uint(i) % 64)
	return out
}

func (b bitset) without(i int) bitset {
	out := append(bitset(nil), b...)
	if i/64 < len(out) {
		out[i/64] &^= 1 << (uint(i) % 64)
	}
	return out
}

// termSet is a list of distinct terms in the order they first appeared.
type termSet struct {
	terms []bitset
	seen  map[string]bool
}

func (t *termSet) add(s *Session, call *Lang, b bitset) {
	k := b.key()
	if t.seen == nil {
		t.seen = map[string]bool{}
	}
	if t.seen[k] {
		return
	}
	if len(t.terms) >= maxTerms {
		s.errorf(call, "the formula expands to more than %d terms", maxTerms)
	}
	t.seen[k] = true
	t.terms = append(t.terms, b)
}

type formVars struct {
	exprs []Value
	index map[string]int
}

type termsCtx struct {
	s     *Session
	call  *Lang
	vars  formVars
	dots  []string // names that a `.` stands for
	hasDt bool
	allow bool // a `.` without data is a variable
	steps int
}

// termsResult is what a (sub)formula contributes: terms and an intercept
// decision (-1: none made, 0: removed, 1: present).
type termsResult struct {
	t    termSet
	icpt int
}

func (x *termsCtx) varKey(v Value) string {
	return strings.Join(syntax.Deparse(x.s.toSyntax(v), nil), "\n")
}

func (x *termsCtx) variable(v Value) int {
	k := x.varKey(v)
	if i, ok := x.vars.index[k]; ok {
		return i
	}
	if x.vars.index == nil {
		x.vars.index = map[string]int{}
	}
	if len(x.vars.exprs) >= maxFactorLen/16 {
		x.s.errorf(x.call, "the formula has too many variables")
	}
	x.vars.exprs = append(x.vars.exprs, v)
	i := len(x.vars.exprs) - 1
	x.vars.index[k] = i
	return i
}

func (x *termsCtx) bad() {
	x.s.errorf(x.call, "invalid model formula in ExtractVars")
}

// cross is the terms of a:b, every term of l joined with every term of r.
func (x *termsCtx) cross(l, r []bitset) termSet {
	var out termSet
	if len(l)*len(r) > maxTerms*4 {
		x.s.errorf(x.call, "the formula expands to more than %d terms", maxTerms)
	}
	for _, a := range l {
		for _, b := range r {
			out.add(x.s, x.call, a.union(b))
		}
	}
	return out
}

func (x *termsCtx) eval(e Value) termsResult {
	x.s.descend()
	defer x.s.ascend()
	x.steps++
	switch v := e.(type) {
	case *Double:
		if len(v.V) == 1 && (v.V[0] == 0 || v.V[0] == 1) {
			return termsResult{icpt: int(v.V[0])}
		}
		x.bad()
	case *Integer:
		if len(v.V) == 1 && (v.V[0] == 0 || v.V[0] == 1) {
			return termsResult{icpt: int(v.V[0])}
		}
		x.bad()
	case *Symbol:
		if v.Name == "." && x.hasDt {
			var out termsResult
			out.icpt = -1
			for _, n := range x.dots {
				out.t.add(x.s, x.call, bitset(nil).with(x.variable(Intern(n))))
			}
			return out
		}
		if v.Name == "." && !x.allow {
			x.s.errorf(x.call, "'.' in formula and no 'data' argument")
		}
		return x.single(e)
	case *Lang:
		op := ""
		if f, ok := v.Fun.(*Symbol); ok {
			op = f.Name
		}
		n := len(v.Args)
		arg := func(i int) Value { return v.Args[i].Value }
		switch {
		case op == "(" && n == 1:
			return x.eval(arg(0))
		case op == "+" && n == 1:
			return x.eval(arg(0))
		case op == "+" && n == 2:
			l, r := x.eval(arg(0)), x.eval(arg(1))
			out := termsResult{icpt: l.icpt}
			if r.icpt >= 0 {
				out.icpt = r.icpt
			}
			out.t = l.t
			out.t.terms = append([]bitset(nil), l.t.terms...)
			out.t.seen = map[string]bool{}
			for k := range l.t.seen {
				out.t.seen[k] = true
			}
			for _, b := range r.t.terms {
				out.t.add(x.s, x.call, b)
			}
			return out
		case op == "-" && n == 1:
			r := x.eval(arg(0))
			out := termsResult{icpt: -1}
			if r.icpt == 1 {
				out.icpt = 0
			} else if r.icpt == 0 {
				out.icpt = 1
			}
			return out
		case op == "-" && n == 2:
			l, r := x.eval(arg(0)), x.eval(arg(1))
			out := termsResult{icpt: l.icpt}
			switch r.icpt {
			case 1:
				out.icpt = 0
			case 0:
				out.icpt = 1
			}
			for _, b := range l.t.terms {
				if !r.t.seen[b.key()] {
					out.t.add(x.s, x.call, b)
				}
			}
			return out
		case op == "%in%" && n == 2:
			l, r := x.eval(arg(0)), x.eval(arg(1))
			out := termsResult{icpt: -1}
			var all bitset
			for _, b := range r.t.terms {
				all = all.union(b)
			}
			for _, b := range l.t.terms {
				out.t.add(x.s, x.call, b.union(all))
			}
			return out
		case op == ":" && n == 2:
			l, r := x.eval(arg(0)), x.eval(arg(1))
			return termsResult{t: x.cross(l.t.terms, r.t.terms), icpt: -1}
		case op == "*" && n == 2:
			l, r := x.eval(arg(0)), x.eval(arg(1))
			out := termsResult{icpt: -1}
			for _, b := range l.t.terms {
				out.t.add(x.s, x.call, b)
			}
			for _, b := range r.t.terms {
				out.t.add(x.s, x.call, b)
			}
			for _, b := range x.cross(l.t.terms, r.t.terms).terms {
				out.t.add(x.s, x.call, b)
			}
			return out
		case op == "/" && n == 2:
			l, r := x.eval(arg(0)), x.eval(arg(1))
			out := termsResult{icpt: -1}
			var all bitset
			for _, b := range l.t.terms {
				out.t.add(x.s, x.call, b)
				all = all.union(b)
			}
			for _, b := range r.t.terms {
				out.t.add(x.s, x.call, all.union(b))
			}
			return out
		case op == "^" && n == 2:
			l := x.eval(arg(0))
			pw := 0.0
			switch p := arg(1).(type) {
			case *Double:
				if len(p.V) == 1 {
					pw = p.V[0]
				}
			case *Integer:
				if len(p.V) == 1 {
					pw = float64(p.V[0])
				}
			}
			if pw < 1 || pw != float64(int(pw)) || pw > 64 {
				x.s.errorf(x.call, "invalid power in formula")
			}
			out := termsResult{icpt: -1, t: l.t}
			for i := 1; i < int(pw); i++ {
				var next termSet
				for _, b := range out.t.terms {
					next.add(x.s, x.call, b)
				}
				for _, b := range l.t.terms {
					next.add(x.s, x.call, b)
				}
				for _, b := range x.cross(out.t.terms, l.t.terms).terms {
					next.add(x.s, x.call, b)
				}
				out.t = next
			}
			return out
		}
		return x.single(e)
	}
	return x.single(e)
}

func (x *termsCtx) single(e Value) termsResult {
	switch e.(type) {
	case *Symbol, *Lang:
	default:
		x.bad()
	}
	var out termsResult
	out.icpt = -1
	out.t.add(x.s, x.call, bitset(nil).with(x.variable(e)))
	return out
}

// dotReplace replaces `.` in e by the variables of the data, in parentheses
// unless it is the whole right-hand side.
func (s *Session) dotReplace(e Value, dots []string, parent string) Value {
	s.descend()
	defer s.ascend()
	switch v := e.(type) {
	case *Symbol:
		if v.Name != "." {
			return e
		}
		if len(dots) == 0 {
			return e
		}
		var sum Value = Intern(dots[0])
		for _, n := range dots[1:] {
			sum = &Lang{Fun: Intern("+"), Args: []Arg{{Value: sum}, {Value: Intern(n)}}}
		}
		switch parent {
		case "*", "-", "^", ":", "/":
			if len(dots) > 1 {
				sum = &Lang{Fun: Intern("("), Args: []Arg{{Value: sum}}}
			}
		}
		return sum
	case *Lang:
		f, ok := v.Fun.(*Symbol)
		if !ok {
			return e
		}
		switch f.Name {
		case "+", "-", "*", "/", "^", ":", "%in%", "(":
		default:
			return e // only the operators of a formula look for a dot
		}
		out := &Lang{Fun: v.Fun, Args: make([]Arg, len(v.Args)), Attr: v.Attr}
		for i, a := range v.Args {
			out.Args[i] = a
			out.Args[i].Value = s.dotReplace(a.Value, dots, f.Name)
		}
		return out
	}
	return e
}

func (s *Session) langHasDot(e Value) bool {
	s.descend()
	defer s.ascend()
	switch v := e.(type) {
	case *Symbol:
		return v.Name == "."
	case *Lang:
		for _, a := range v.Args {
			if s.langHasDot(a.Value) {
				return true
			}
		}
	}
	return false
}

func (s *Session) langVars(e Value, funcs bool, out *[]string) {
	s.descend()
	defer s.ascend()
	switch v := e.(type) {
	case *Symbol:
		if v.Name != "" {
			*out = append(*out, v.Name)
		}
	case *Lang:
		if funcs {
			if f, ok := v.Fun.(*Symbol); ok {
				*out = append(*out, f.Name)
			} else {
				s.langVars(v.Fun, funcs, out)
			}
		} else if _, ok := v.Fun.(*Symbol); !ok {
			s.langVars(v.Fun, funcs, out)
		}
		for _, a := range v.Args {
			s.langVars(a.Value, funcs, out)
		}
	case *List:
		if v.Expr {
			for _, e := range v.V {
				s.langVars(e, funcs, out)
			}
		}
	}
}

func init() {
	// ..termsFormula(call, x, specials, dotnames, keep.order, simplify, allowDot)
	regBuiltin("..termsFormula", func(s *Session, c *callInfo) Value {
		call, _ := c.args[0].(*Lang)
		if call == nil {
			call = c.call
		}
		f, ok := c.args[1].(*Lang)
		if !ok || len(f.Args) < 1 || len(f.Args) > 2 {
			s.errorf(call, "argument is not a valid model")
		}
		if sym, ok := f.Fun.(*Symbol); !ok || sym.Name != "~" {
			s.errorf(call, "argument is not a valid model")
		}
		var specials []string
		if sp, ok := c.args[2].(*String); ok {
			for _, e := range sp.V {
				specials = append(specials, e.S)
			}
		}
		x := &termsCtx{s: s, call: call}
		if dn, ok := c.args[3].(*String); ok {
			x.hasDt = true
			for _, e := range dn.V {
				x.dots = append(x.dots, e.S)
			}
		}
		keepOrder := len(logicalOf(s, c, c.args[4], "invalid 'keep.order'")) > 0 && logicalOf(s, c, c.args[4], "")[0] == 1
		simplify := len(logicalOf(s, c, c.args[5], "invalid 'simplify'")) > 0 && logicalOf(s, c, c.args[5], "")[0] == 1
		x.allow = len(logicalOf(s, c, c.args[6], "invalid 'allowDotAsName'")) > 0 && logicalOf(s, c, c.args[6], "")[0] == 1

		resp := 0
		var lhs, rhs Value
		if len(f.Args) == 2 {
			lhs, rhs = f.Args[0].Value, f.Args[1].Value
		} else {
			rhs = f.Args[0].Value
		}
		if lhs != nil {
			resp = 1
			x.variable(lhs)
		}
		// the variables of the left-hand side are not what a `.` stands for
		if x.hasDt && lhs != nil {
			var lv []string
			s.langVars(lhs, false, &lv)
			drop := map[string]bool{}
			for _, n := range lv {
				drop[n] = true
			}
			var keep []string
			for _, n := range x.dots {
				if !drop[n] {
					keep = append(keep, n)
				}
			}
			x.dots = keep
		}
		out := f
		if x.hasDt && s.langHasDot(rhs) {
			nf := &Lang{Fun: f.Fun, Args: append([]Arg(nil), f.Args...)}
			nf.Args[len(nf.Args)-1].Value = s.dotReplace(rhs, x.dots, "")
			out = nf
		}
		res := x.eval(rhs)
		icpt := 1
		if res.icpt >= 0 {
			icpt = res.icpt
		}
		// offsets are variables whose call is offset(...)
		isOff := func(i int) bool {
			l, ok := x.vars.exprs[i].(*Lang)
			if !ok {
				return false
			}
			f, ok := l.Fun.(*Symbol)
			return ok && f.Name == "offset"
		}
		var offs []int32
		for i := range x.vars.exprs {
			if isOff(i) && i >= resp {
				offs = append(offs, int32(i+1))
			}
		}
		var terms []bitset
	next:
		for _, b := range res.t.terms {
			for i := range x.vars.exprs {
				if isOff(i) && b.has(i) {
					continue next
				}
			}
			terms = append(terms, b)
		}
		if !keepOrder {
			sort.SliceStable(terms, func(i, j int) bool { return terms[i].count() < terms[j].count() })
		}
		nv := len(x.vars.exprs)
		if nv*len(terms) > maxFactorLen {
			s.errorf(call, "the formula has too many terms")
		}
		vnames := make([]Str, nv)
		for i, e := range x.vars.exprs {
			vnames[i] = Str{S: x.varKey(e)}
		}
		labels := make([]Str, len(terms))
		order := make([]int32, len(terms))
		var factors Value = &Integer{V: []int32{}}
		if len(terms) > 0 {
			fm := make([]int32, nv*len(terms))
			for j, b := range terms {
				var parts []string
				for i := 0; i < nv; i++ {
					if !b.has(i) {
						continue
					}
					parts = append(parts, vnames[i].S)
					order[j]++
					// contrasts (1) when what is left of the term once the variable is
					// taken out is nothing, or lies inside a term that comes before it
					code := int32(1)
					rest := b.without(i)
					if rest.count() > 0 {
						code = 2
						for _, earlier := range terms[:j] {
							if rest.subsetOf(earlier) {
								code = 1
								break
							}
						}
					}
					fm[j*nv+i] = code
				}
				labels[j] = Str{S: strings.Join(parts, ":")}
			}
			factors = &Integer{V: fm, Attr: []Attr{
				{Name: "dim", Value: &Integer{V: []int32{int32(nv), int32(len(terms))}}},
				{Name: "dimnames", Value: &List{V: []Value{&String{V: vnames}, &String{V: labels}}}},
			}}
		}
		varCall := &Lang{Fun: Intern("list")}
		for _, e := range x.vars.exprs {
			varCall.Args = append(varCall.Args, Arg{Value: e})
		}
		attrs := []Attr{{Name: "variables", Value: varCall}}
		if len(offs) > 0 {
			attrs = append(attrs, Attr{Name: "offset", Value: &Integer{V: offs}})
		}
		attrs = append(attrs, Attr{Name: "factors", Value: factors}, Attr{Name: "term.labels", Value: &String{V: labels}})
		if sp, ok := c.args[2].(*String); ok && len(sp.V) > 0 {
			sl := &List{Pair: true}
			var nm []Str
			for _, name := range specials {
				var idx []int32
				for i, e := range x.vars.exprs {
					if l, ok := e.(*Lang); ok {
						if hf, ok := l.Fun.(*Symbol); ok && hf.Name == name {
							idx = append(idx, int32(i+1))
						}
					}
				}
				if len(idx) > 0 {
					sl.V = append(sl.V, &Integer{V: idx})
				} else {
					sl.V = append(sl.V, Null)
				}
				nm = append(nm, Str{S: name})
			}
			sl.Attr = namesAttr(nm)
			attrs = append(attrs, Attr{Name: "specials", Value: sl})
		}
		attrs = append(attrs,
			Attr{Name: "order", Value: &Integer{V: order}},
			Attr{Name: "intercept", Value: mkInt(int32(icpt))},
			Attr{Name: "response", Value: mkInt(int32(resp))},
			Attr{Name: "class", Value: strVec([]string{"terms", "formula"})})
		var env Value = Null
		if e := getAttr(f, ".Environment"); e != nil {
			env = e
		}
		attrs = append(attrs, Attr{Name: ".Environment", Value: env})
		if simplify {
			var rhsV Value
			add := func(v Value) {
				if rhsV == nil {
					rhsV = v
				} else {
					rhsV = &Lang{Fun: Intern("+"), Args: []Arg{{Value: rhsV}, {Value: v}}}
				}
			}
			for j := range terms {
				var tv Value
				for i := 0; i < nv; i++ {
					if terms[j].has(i) {
						if tv == nil {
							tv = x.vars.exprs[i]
						} else {
							tv = &Lang{Fun: Intern(":"), Args: []Arg{{Value: tv}, {Value: x.vars.exprs[i]}}}
						}
					}
				}
				add(tv)
			}
			for _, o := range offs {
				add(x.vars.exprs[o-1])
			}
			if rhsV == nil {
				rhsV = mkDouble(1)
				if icpt == 0 {
					rhsV = mkDouble(0)
				}
			} else if icpt == 0 {
				rhsV = &Lang{Fun: Intern("-"), Args: []Arg{{Value: rhsV}, {Value: mkDouble(1)}}}
			}
			nf := &Lang{Fun: f.Fun}
			if lhs != nil {
				nf.Args = append(nf.Args, Arg{Value: lhs})
			}
			nf.Args = append(nf.Args, Arg{Value: rhsV})
			out = nf
		}
		return &Lang{Fun: out.Fun, Args: out.Args, Attr: attrs}
	})

	// all.names(expr, functions, max.names, unique) and all.vars
	names := func(isVars bool) func(*Session, *callInfo) Value {
		return func(s *Session, c *callInfo) Value {
			a := s.matchFormals(c, "expr", "functions", "max.names", "unique")
			if a[0] == nil {
				s.errorf(c.call, "argument \"expr\" is missing, with no default")
			}
			flag := func(v Value, def bool) bool {
				if v == nil {
					return def
				}
				l := logicalOf(s, c, v, "invalid argument")
				return len(l) > 0 && l[0] == 1
			}
			funcs := flag(a[1], !isVars)
			uniq := flag(a[3], isVars)
			maxN := -1
			if a[2] != nil {
				if d := asDoubles(a[2]); len(d) == 1 && d[0] >= 0 {
					maxN = int(d[0])
				}
			}
			var out []string
			s.langVars(a[0], funcs, &out)
			if uniq {
				seen := map[string]bool{}
				var u []string
				for _, n := range out {
					if !seen[n] {
						seen[n] = true
						u = append(u, n)
					}
				}
				out = u
			}
			if maxN >= 0 && len(out) > maxN {
				out = out[:maxN]
			}
			return strVec(out)
		}
	}
	regBuiltin("all.names", names(false))
	regBuiltin("all.vars", names(true))
}
