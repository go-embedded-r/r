package eval

import (
	"math"
	"strconv"
	"strings"

	"github.com/go-embedded-r/r/syntax"
)

// ---- coercion to character, for comparison, c() and as.character ----------

func numToString(x float64) Str {
	if isNAReal(x) {
		return Str{NA: true}
	}
	return Str{S: syntax.FormatNumber(x)}
}

func asStrings(v Value) []Str {
	switch x := v.(type) {
	case *String:
		return x.V
	case *Logical:
		out := make([]Str, len(x.V))
		for i, e := range x.V {
			switch e {
			case naInt:
				out[i] = Str{NA: true}
			case 0:
				out[i] = Str{S: "FALSE"}
			default:
				out[i] = Str{S: "TRUE"}
			}
		}
		return out
	case *Integer:
		out := make([]Str, len(x.V))
		for i, e := range x.V {
			if e == naInt {
				out[i] = Str{NA: true}
			} else {
				out[i] = Str{S: strconv.Itoa(int(e))}
			}
		}
		return out
	case *Double:
		out := make([]Str, len(x.V))
		for i, e := range x.V {
			out[i] = numToString(e)
		}
		return out
	}
	return nil
}

// ---- comparison ------------------------------------------------------------

func (s *Session) compare(op string, c *callInfo) Value {
	if len(c.args) != 2 {
		s.errorf(c.call, "operator needs two arguments")
	}
	x, y := c.args[0], c.args[1]
	kx, ky := kindOf(x), kindOf(y)
	if kx > kString || ky > kString {
		s.errorf(c.call, "comparison (%s) is possible only for atomic and list types", op)
	}
	nx, ny := Length(x), Length(y)
	n := recycleLen(nx, ny)
	s.checkLen(c.call, n)
	if n > 0 && (n%nx != 0 || n%ny != 0) {
		s.warningf(c.call, "longer object length is not a multiple of shorter object length")
	}
	out := make([]int32, n)
	test := func(c int) int32 { // c is -1, 0, 1
		var r bool
		switch op {
		case "==":
			r = c == 0
		case "!=":
			r = c != 0
		case "<":
			r = c < 0
		case "<=":
			r = c <= 0
		case ">":
			r = c > 0
		case ">=":
			r = c >= 0
		}
		if r {
			return 1
		}
		return 0
	}
	if kx == kString || ky == kString {
		a, b := asStrings(x), asStrings(y)
		for i := range out {
			p, q := a[i%nx], b[i%ny]
			if p.NA || q.NA {
				out[i] = naInt
				continue
			}
			if op == "==" || op == "!=" {
				out[i] = test(strings.Compare(p.S, q.S))
			} else {
				out[i] = test(collate(p.S, q.S))
			}
		}
	} else {
		a, b := asDoubles(x), asDoubles(y)
		for i := range out {
			p, q := a[i%nx], b[i%ny]
			if math.IsNaN(p) || math.IsNaN(q) {
				out[i] = naInt
				continue
			}
			cmp := 0
			switch {
			case p < q:
				cmp = -1
			case p > q:
				cmp = 1
			}
			out[i] = test(cmp)
		}
	}
	return &Logical{V: out, Attr: copyAttrs(x, y, n)}
}

// ---- logic -----------------------------------------------------------------

func logicalOf(s *Session, c *callInfo, v Value, what string) []int32 {
	switch x := v.(type) {
	case nullT:
		return nil
	case *Logical:
		return x.V
	case *Integer:
		out := make([]int32, len(x.V))
		for i, e := range x.V {
			switch {
			case e == naInt:
				out[i] = naInt
			case e != 0:
				out[i] = 1
			}
		}
		return out
	case *Double:
		out := make([]int32, len(x.V))
		for i, e := range x.V {
			switch {
			case math.IsNaN(e):
				out[i] = naInt
			case e != 0:
				out[i] = 1
			}
		}
		return out
	}
	s.errorf(c.call, "%s", what)
	return nil
}

func (s *Session) logic(op string, c *callInfo) Value {
	if len(c.args) != 2 {
		s.errorf(c.call, "operator needs two arguments")
	}
	x, y := c.args[0], c.args[1]
	const msg = "operations are possible only for numeric, logical or complex types"
	a, b := logicalOf(s, c, x, msg), logicalOf(s, c, y, msg)
	nx, ny := len(a), len(b)
	n := recycleLen(nx, ny)
	s.checkLen(c.call, n)
	if n > 0 && (n%nx != 0 || n%ny != 0) {
		s.warningf(c.call, "longer object length is not a multiple of shorter object length")
	}
	out := make([]int32, n)
	for i := range out {
		p, q := a[i%nx], b[i%ny]
		if op == "&" {
			switch {
			case p == 0 || q == 0:
				out[i] = 0
			case p == naInt || q == naInt:
				out[i] = naInt
			default:
				out[i] = 1
			}
		} else {
			switch {
			case p == 1 || q == 1:
				out[i] = 1
			case p == naInt || q == naInt:
				out[i] = naInt
			default:
				out[i] = 0
			}
		}
	}
	return &Logical{V: out, Attr: copyAttrs(x, y, n)}
}

func (s *Session) not(c *callInfo) Value {
	if len(c.args) != 1 {
		s.errorf(c.call, "%d arguments passed to '!' which requires 1", len(c.args))
	}
	x := c.args[0]
	if _, isNull := x.(nullT); isNull {
		s.errorf(c.call, "invalid argument type")
	}
	a := logicalOf(s, c, x, "invalid argument type")
	out := make([]int32, len(a))
	for i, e := range a {
		switch e {
		case naInt:
			out[i] = naInt
		case 0:
			out[i] = 1
		}
	}
	return &Logical{V: out, Attr: attrsOf(x)}
}

// scalarLogical is the operand of && and ||: a length-one logical-like value,
// NA allowed; a zero-length operand is NA.
func (s *Session) scalarLogical(c *callInfo, v Value, which string) int32 {
	switch v.(type) {
	case *Logical, *Integer, *Double:
	default:
		s.errorf(c.call, "invalid '%s' type in 'x %s y'", which, c.call.Fun.(*Symbol).Name)
	}
	if n := Length(v); n != 1 {
		if n == 0 {
			return naInt
		}
		s.errorf(c.call, "'length = %d' in coercion to 'logical(1)'", n)
	}
	return logicalOf(s, c, v, "")[0]
}

// ---- the colon operator, c(), list(), length(), is.na() -------------------

func (s *Session) colon(c *callInfo) Value {
	if len(c.args) != 2 {
		s.errorf(c.call, "operator needs two arguments")
	}
	from, to := c.args[0], c.args[1]
	for _, v := range []Value{from, to} {
		if Length(v) == 0 {
			s.errorf(c.call, "argument of length 0")
		}
	}
	if Length(from) > 1 {
		s.warningf(c.call, "numerical expression has %d elements: only the first used", Length(from))
	}
	if Length(to) > 1 {
		s.warningf(c.call, "numerical expression has %d elements: only the first used", Length(to))
	}
	num := func(v Value) float64 {
		switch kindOf(v) {
		case kLogical, kInteger, kDouble:
			return asDoubles(v)[0]
		case kString:
			f, ok := parseNumber(asStrings(v)[0])
			if !ok {
				return rNaN
			}
			return f
		}
		s.errorf(c.call, "argument of length 0")
		return 0
	}
	a, b := num(from), num(to)
	if math.IsNaN(a) || math.IsNaN(b) {
		s.errorf(c.call, "NA/NaN argument")
	}
	nf := math.Floor(math.Abs(a-b)+1e-10) + 1
	if math.IsNaN(nf) { // Inf:Inf (Inf - Inf is NaN): R answers numeric(0)
		return &Double{V: []float64{}}
	}
	if nf > 1<<52 {
		s.errorf(c.call, "result would be too long a vector")
	}
	n := int(math.Min(nf, float64(s.opts.MaxVector)+1)) // checkLen refuses it
	s.checkLen(c.call, n)
	useInt := a == math.Floor(a) && a <= math.MaxInt32 && a > math.MinInt32
	if useInt {
		last := a + float64(n-1)
		if a > b {
			last = a - float64(n-1)
		}
		useInt = last <= math.MaxInt32 && last > math.MinInt32
	}
	if useInt {
		out := make([]int32, n)
		for i := range out {
			if a <= b {
				out[i] = int32(a) + int32(i)
			} else {
				out[i] = int32(a) - int32(i)
			}
		}
		return &Integer{V: out}
	}
	out := make([]float64, n)
	for i := range out {
		if a <= b {
			out[i] = a + float64(i)
		} else {
			out[i] = a - float64(i)
		}
	}
	return &Double{V: out}
}

func (s *Session) combine(c *callInfo) Value {
	kind := kNull
	total := 0
	for _, a := range c.args {
		k := kindOf(a)
		switch k {
		case kNull:
			continue
		case kOther:
			k = kList
		}
		if k > kind {
			kind = k
		}
		if kindOf(a) == kOther {
			total++
		} else {
			total += Length(a)
		}
	}
	s.checkLen(c.call, total)
	// names: the tag of an argument names its elements (tag for one element,
	// tag1, tag2, ... for several), alongside the names the argument has
	var names []Str
	anyNames := false
	for i, a := range c.args {
		if kindOf(a) == kNull {
			continue
		}
		if c.names[i] != "" || namesOf(a) != nil {
			anyNames = true
		}
	}
	addNames := func(i int, a Value, n int) {
		if !anyNames {
			return
		}
		tag := c.names[i]
		nm := namesOf(a)
		for j := 0; j < n; j++ {
			var inner string
			if nm != nil {
				inner = nm.V[j].S
			}
			switch {
			case tag != "" && inner != "":
				names = append(names, Str{S: tag + "." + inner})
			case tag != "" && n == 1:
				names = append(names, Str{S: tag})
			case tag != "":
				names = append(names, Str{S: tag + strconv.Itoa(j+1)})
			default:
				names = append(names, Str{S: inner})
			}
		}
	}
	var attr []Attr
	finish := func(v Value) Value {
		if anyNames {
			attr = namesAttr(names)
		}
		return withAttr(v, attr)
	}
	switch kind {
	case kNull:
		return Null
	case kLogical:
		out := make([]int32, 0, total)
		for i, a := range c.args {
			if kindOf(a) == kNull {
				continue
			}
			out = append(out, asInts(a)...)
			addNames(i, a, Length(a))
		}
		return finish(&Logical{V: out})
	case kInteger:
		out := make([]int32, 0, total)
		for i, a := range c.args {
			if kindOf(a) == kNull {
				continue
			}
			out = append(out, asInts(a)...)
			addNames(i, a, Length(a))
		}
		return finish(&Integer{V: out})
	case kDouble:
		out := make([]float64, 0, total)
		for i, a := range c.args {
			if kindOf(a) == kNull {
				continue
			}
			out = append(out, asDoubles(a)...)
			addNames(i, a, Length(a))
		}
		return finish(&Double{V: out})
	case kString:
		out := make([]Str, 0, total)
		for i, a := range c.args {
			if kindOf(a) == kNull {
				continue
			}
			out = append(out, asStrings(a)...)
			addNames(i, a, Length(a))
		}
		return finish(&String{V: out})
	}
	out := make([]Value, 0, total)
	for i, a := range c.args {
		switch x := a.(type) {
		case nullT:
		case *List:
			out = append(out, x.V...)
			addNames(i, a, len(x.V))
		case *Logical, *Integer, *Double, *String:
			for j := 0; j < Length(a); j++ {
				out = append(out, element(a, j))
			}
			addNames(i, a, Length(a))
		default:
			out = append(out, a)
			addNames(i, a, 1)
		}
	}
	return finish(&List{V: out})
}

func (s *Session) list(c *callInfo) Value {
	out := &List{V: append([]Value(nil), c.args...)}
	for _, n := range c.names {
		if n != "" {
			nm := make([]Str, len(c.names))
			for i, t := range c.names {
				nm[i] = Str{S: t}
			}
			out.Attr = namesAttr(nm)
			break
		}
	}
	return out
}

func (s *Session) isNA(c *callInfo) Value {
	if len(c.args) != 1 {
		s.errorf(c.call, "%d arguments passed to 'is.na' which requires 1", len(c.args))
	}
	x := c.args[0]
	var out []int32
	b := func(v bool) int32 {
		if v {
			return 1
		}
		return 0
	}
	switch v := x.(type) {
	case *Logical:
		for _, e := range v.V {
			out = append(out, b(e == naInt))
		}
	case *Integer:
		for _, e := range v.V {
			out = append(out, b(e == naInt))
		}
	case *Double:
		for _, e := range v.V {
			out = append(out, b(math.IsNaN(e)))
		}
	case *String:
		for _, e := range v.V {
			out = append(out, b(e.NA))
		}
	case *List:
		for _, e := range v.V {
			r := int32(0)
			if Length(e) == 1 {
				switch e.(type) {
				case *Logical, *Integer, *Double, *String:
					r = s.isNA(&callInfo{call: c.call, args: []Value{e}}).(*Logical).V[0]
				}
			}
			out = append(out, r)
		}
	case nullT:
		return &Logical{V: []int32{}}
	default:
		s.warningf(c.call, "is.na() applied to non-(list or vector) of type '%s'", typeName(x))
		return mkLogical(false)
	}
	var attr []Attr
	for _, a := range attrsOf(x) {
		if a.Name == "names" || a.Name == "dim" || a.Name == "dimnames" {
			attr = append(attr, a)
		}
	}
	return &Logical{V: out, Attr: attr}
}

func typeName(v Value) string {
	switch v.(type) {
	case nullT:
		return "NULL"
	case *Symbol:
		return "symbol"
	case *Logical:
		return "logical"
	case *Integer:
		return "integer"
	case *Double:
		return "double"
	case *String:
		return "character"
	case *List:
		return "list"
	case *Lang:
		return "language"
	case *Closure:
		return "closure"
	case *Builtin:
		if v.(*Builtin).Special {
			return "special"
		}
		return "builtin"
	case *Env:
		return "environment"
	}
	return "unknown"
}

func init() {
	for _, op := range []string{"==", "!=", "<", "<=", ">", ">="} {
		op := op
		regBuiltin(op, func(s *Session, c *callInfo) Value { return s.compare(op, c) })
	}
	regBuiltin("&", func(s *Session, c *callInfo) Value { return s.logic("&", c) })
	regBuiltin("|", func(s *Session, c *callInfo) Value { return s.logic("|", c) })
	regBuiltin("!", func(s *Session, c *callInfo) Value { return s.not(c) })
	regSpecial("&&", func(s *Session, c *callInfo) Value {
		if len(c.raw) != 2 {
			s.errorf(nil, "'&&' operator requires 2 arguments")
		}
		x := s.scalarLogical(c, s.eval(c.raw[0].Value, c.env), "x")
		if x == 0 {
			return mkLogical(false)
		}
		y := s.scalarLogical(c, s.eval(c.raw[1].Value, c.env), "y")
		switch {
		case y == 0:
			return mkLogical(false)
		case x == naInt || y == naInt:
			return &Logical{V: []int32{naInt}}
		}
		return mkLogical(true)
	})
	regSpecial("||", func(s *Session, c *callInfo) Value {
		if len(c.raw) != 2 {
			s.errorf(nil, "'||' operator requires 2 arguments")
		}
		x := s.scalarLogical(c, s.eval(c.raw[0].Value, c.env), "x")
		if x == 1 {
			return mkLogical(true)
		}
		y := s.scalarLogical(c, s.eval(c.raw[1].Value, c.env), "y")
		switch {
		case y == 1:
			return mkLogical(true)
		case x == naInt || y == naInt:
			return &Logical{V: []int32{naInt}}
		}
		return mkLogical(false)
	})
	regBuiltin(":", func(s *Session, c *callInfo) Value { return s.colon(c) })
	regBuiltin("c", func(s *Session, c *callInfo) Value { return s.combine(c) })
	regBuiltin("list", func(s *Session, c *callInfo) Value { return s.list(c) })
	regBuiltin("length", func(s *Session, c *callInfo) Value {
		if len(c.args) != 1 {
			s.errorf(c.call, "%d arguments passed to 'length' which requires 1", len(c.args))
		}
		return mkInt(int32(Length(c.args[0])))
	})
	regBuiltin("is.na", func(s *Session, c *callInfo) Value { return s.isNA(c) })
}
