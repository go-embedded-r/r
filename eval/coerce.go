package eval

import (
	"math"
	"strconv"
	"strings"
)

// parseNumber is R's String2Real for the common spellings: blanks around, a
// decimal or hexadecimal number, Inf, NaN, NA.
func parseNumber(s Str) (float64, bool) {
	if s.NA {
		return naReal, true
	}
	t := strings.TrimSpace(s.S)
	if t == "" {
		return 0, false
	}
	neg := false
	body := t
	if body[0] == '+' || body[0] == '-' {
		neg = body[0] == '-'
		body = body[1:]
	}
	// R_strtod: NA, NaN and Inf are matched case-sensitively, infinity not
	signed := func(f float64) float64 {
		if neg {
			return -f
		}
		return f
	}
	switch {
	case body == "NA":
		return naReal, true
	case body == "NaN":
		return rNaN, true
	case body == "Inf" || strings.EqualFold(body, "infinity"):
		return signed(math.Inf(1)), true
	}
	var f float64
	if len(body) > 1 && body[0] == '0' && (body[1] == 'x' || body[1] == 'X') {
		v, err := strconv.ParseFloat(body, 64)
		if err != nil && !math.IsInf(v, 0) {
			// 0x without a p exponent: ParseFloat wants one
			u, err2 := strconv.ParseUint(body[2:], 16, 64)
			if err2 != nil {
				return 0, false
			}
			v = float64(u)
		}
		f = v
	} else {
		for _, c := range body {
			if !(c >= '0' && c <= '9') && c != '.' && c != 'e' && c != 'E' && c != '+' && c != '-' {
				return 0, false
			}
		}
		v, err := strconv.ParseFloat(body, 64)
		if err != nil && !math.IsInf(v, 0) {
			return 0, false
		}
		f = v
	}
	if neg {
		f = -f
	}
	return f, true
}

// firstArg is the data argument of as.integer() and friends, which are
// primitives that ignore any further argument.
func firstArg(c *callInfo) Value {
	if len(c.args) == 0 {
		return Null
	}
	return c.args[0]
}

func (s *Session) nargs1(c *callInfo, name string) Value {
	if len(c.args) != 1 {
		s.errorf(c.call, "%d arguments passed to '%s' which requires 1", len(c.args), name)
	}
	return c.args[0]
}

func (s *Session) asDouble(c *callInfo, name string) Value {
	x := firstArg(c)
	switch v := x.(type) {
	case nullT:
		return &Double{V: []float64{}}
	case *Double:
		return &Double{V: v.V}
	case *Logical, *Integer:
		return &Double{V: asDoubles(x)}
	case *String:
		out := make([]float64, len(v.V))
		warn := false
		for i, e := range v.V {
			f, ok := parseNumber(e)
			if !ok {
				f, warn = naReal, true
			}
			out[i] = f
		}
		if warn {
			s.warningf(c.call, "NAs introduced by coercion")
		}
		return &Double{V: out}
	}
	s.errorf(c.call, "cannot coerce type '%s' to vector of type 'double'", typeName(x))
	return nil
}

func (s *Session) asInteger(c *callInfo) Value {
	x := firstArg(c)
	switch v := x.(type) {
	case nullT:
		return &Integer{V: []int32{}}
	case *Logical, *Integer:
		return &Integer{V: asInts(x)}
	case *Double, *String:
		var d []float64
		warnCoerce := false
		if sv, ok := v.(*String); ok {
			d = make([]float64, len(sv.V))
			for i, e := range sv.V {
				f, ok := parseNumber(e)
				if !ok {
					f, warnCoerce = naReal, true
				}
				d[i] = f
			}
		} else {
			d = x.(*Double).V
		}
		out := make([]int32, len(d))
		warnRange := false
		for i, f := range d {
			switch {
			case math.IsNaN(f):
				out[i] = naInt
			case f >= 2147483648 || f <= -2147483649:
				out[i], warnRange = naInt, true
			default:
				out[i] = int32(f)
			}
		}
		if warnCoerce {
			s.warningf(c.call, "NAs introduced by coercion")
		}
		if warnRange {
			s.warningf(c.call, "NAs introduced by coercion to integer range")
		}
		return &Integer{V: out}
	}
	s.errorf(c.call, "cannot coerce type '%s' to vector of type 'integer'", typeName(x))
	return nil
}

func (s *Session) asLogical(c *callInfo) Value {
	x := firstArg(c)
	switch v := x.(type) {
	case nullT:
		return &Logical{V: []int32{}}
	case *Logical:
		return &Logical{V: v.V}
	case *Integer, *Double:
		return &Logical{V: logicalOf(s, c, x, "")}
	case *String:
		out := make([]int32, len(v.V))
		for i, e := range v.V {
			switch {
			case e.NA:
				out[i] = naInt
			case e.S == "T" || e.S == "TRUE" || e.S == "true" || e.S == "True":
				out[i] = 1
			case e.S == "F" || e.S == "FALSE" || e.S == "false" || e.S == "False":
				out[i] = 0
			default:
				out[i] = naInt
			}
		}
		return &Logical{V: out}
	}
	s.errorf(c.call, "cannot coerce type '%s' to vector of type 'logical'", typeName(x))
	return nil
}

func (s *Session) asCharacter(c *callInfo) Value {
	x := firstArg(c)
	switch v := x.(type) {
	case nullT:
		return &String{V: []Str{}}
	case *Symbol:
		return mkString(v.Name)
	case *Logical, *Integer, *Double, *String:
		return &String{V: asStrings(x)}
	}
	s.errorf(c.call, "cannot coerce type '%s' to vector of type 'character'", typeName(x))
	return nil
}

// ---- math and summaries -----------------------------------------------------

func (s *Session) math1(name string, f func(float64) float64, keepInt bool) func(s *Session, c *callInfo) Value {
	return func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, name)
		switch x.(type) {
		case *Logical, *Integer, *Double:
		default:
			s.errorf(c.call, "non-numeric argument to mathematical function")
		}
		if keepInt {
			if _, isDbl := x.(*Double); !isDbl {
				out := make([]int32, 0, Length(x))
				for _, e := range asInts(x) {
					switch {
					case e == naInt:
						out = append(out, naInt)
					default:
						out = append(out, int32(f(float64(e))))
					}
				}
				return &Integer{V: out, Attr: attrsOf(x)}
			}
		}
		d := asDoubles(x)
		out := make([]float64, len(d))
		nan := false
		for i, e := range d {
			if math.IsNaN(e) && name != "abs" { // abs is fabs: it clears the sign of a NaN too
				out[i] = e
				continue
			}
			out[i] = f(e)
			if math.IsNaN(out[i]) && !math.IsNaN(e) {
				nan = true
			}
		}
		if nan {
			s.warningf(c.call, "NaNs produced")
		}
		return &Double{V: out, Attr: attrsOf(x)}
	}
}

// narm separates the na.rm argument from the data arguments of a summary.
func (s *Session) narm(c *callInfo) (data []Value, rm bool) {
	for i, a := range c.args {
		if c.names[i] == "na.rm" {
			rm = len(logicalOf(s, c, a, "invalid 'na.rm' value")) > 0 && logicalOf(s, c, a, "")[0] == 1
			continue
		}
		data = append(data, a)
	}
	return
}

func (s *Session) summary(name string) func(s *Session, c *callInfo) Value {
	return func(s *Session, c *callInfo) Value {
		data, rm := s.narm(c)
		allInt := true
		if name == "min" || name == "max" {
			for _, a := range data {
				if _, isStr := a.(*String); isStr {
					return s.strExtreme(c, name, data, rm)
				}
			}
		}
		for _, a := range data {
			switch a.(type) {
			case nullT, *Logical, *Integer:
			case *Double:
				allInt = false
			default:
				s.errorf(c.call, "invalid 'type' (%s) of argument", typeName(a))
			}
		}
		if allInt && name != "prod" {
			var acc int64
			sawAny := false
			var best int32
			na := false
			for _, a := range data {
				for _, e := range asInts(a) {
					if e == naInt {
						if !rm {
							na = true
						}
						continue
					}
					switch name {
					case "sum":
						acc += int64(e)
					case "min":
						if !sawAny || e < best {
							best = e
						}
					case "max":
						if !sawAny || e > best {
							best = e
						}
					}
					sawAny = true
				}
			}
			switch name {
			case "sum":
				if na {
					return &Integer{V: []int32{naInt}}
				}
				if acc > math.MaxInt32 || acc < -math.MaxInt32 {
					return mkDouble(float64(acc)) // R widens instead of overflowing
				}
				return mkInt(int32(acc))
			default:
				if na {
					return &Integer{V: []int32{naInt}}
				}
				if !sawAny {
					s.warningf(c.call, "no non-missing arguments to %s; returning %s", name, map[string]string{"min": "Inf", "max": "-Inf"}[name])
					if name == "min" {
						return mkDouble(math.Inf(1))
					}
					return mkDouble(math.Inf(-1))
				}
				return mkInt(best)
			}
		}
		if name == "sum" || name == "prod" {
			// R accumulates each argument on its own, then combines the partial
			// results: sum(a, b) is (sum a) + (sum b), prod(0, x) is 0 * (prod x)
			unit := 0.0
			if name == "prod" {
				unit = 1
			}
			total := unit
			totNaN := false
			var nanv float64
			for _, a := range data {
				_, isDbl := a.(*Double)
				part := unit
				partNaN := false
				var pn float64
				for _, e := range asDoubles(a) {
					if math.IsNaN(e) {
						if rm {
							continue
						}
						if name == "sum" && !isDbl {
							return mkDouble(naReal) // an integer NA ends the sum, exactly
						}
						if !partNaN || (isNAReal(e) && !isNAReal(pn)) {
							pn = e
						}
						partNaN = true
						continue
					}
					if name == "sum" {
						part += e
					} else {
						part *= e
					}
				}
				if partNaN {
					cand := quiet(pn)
					if !totNaN || (isNAReal(cand) && !isNAReal(nanv)) {
						nanv = cand
					}
					totNaN = true
				} else if name == "sum" {
					total += part
				} else {
					total *= part
				}
			}
			if totNaN {
				return mkDouble(nanv)
			}
			return mkDouble(total)
		}
		acc := math.Inf(1)
		if name == "max" {
			acc = math.Inf(-1)
		}
		sawAny, sawNaN := false, false
		var naVal float64
		for _, a := range data {
			for _, e := range asDoubles(a) {
				if math.IsNaN(e) {
					if rm {
						continue
					}
					if !sawNaN || (isNAReal(e) && !isNAReal(naVal)) {
						naVal = e // any NA trumps all the NaNs
					}
					sawNaN = true
					continue
				}
				sawAny = true
				if name == "min" && e < acc { // not math.Min: min(0, -0) is 0
					acc = e
				} else if name == "max" && e > acc {
					acc = e
				}
			}
		}
		if sawNaN {
			return mkDouble(naVal)
		}
		if !sawAny && (name == "min" || name == "max") {
			s.warningf(c.call, "no non-missing arguments to %s; returning %s", name, map[string]string{"min": "Inf", "max": "-Inf"}[name])
		}
		return mkDouble(acc)
	}
}

func (s *Session) anyAll(isAny bool) func(s *Session, c *callInfo) Value {
	return func(s *Session, c *callInfo) Value {
		data, rm := s.narm(c)
		sawNA := false
		for _, a := range data {
			var vals []int32
			switch x := a.(type) {
			case nullT, *Logical, *Integer:
				vals = logicalOf(s, c, a, "")
			case *Double:
				s.warningf(c.call, "coercing argument of type 'double' to logical")
				vals = logicalOf(s, c, a, "")
			case *String:
				s.warningf(c.call, "coercing argument of type 'character' to logical")
				vals = s.asLogical(&callInfo{call: c.call, args: []Value{x}}).(*Logical).V
			default:
				s.errorf(c.call, "invalid 'type' (%s) of argument", typeName(a))
			}
			for _, e := range vals {
				switch {
				case e == naInt:
					if !rm {
						sawNA = true
					}
				case isAny && e == 1:
					return mkLogical(true)
				case !isAny && e == 0:
					return mkLogical(false)
				}
			}
		}
		if sawNA {
			return &Logical{V: []int32{naInt}}
		}
		return mkLogical(!isAny)
	}
}

func init() {
	pred := func(name string, f func(Value) bool) {
		regBuiltin(name, func(s *Session, c *callInfo) Value { return mkLogical(f(s.nargs1(c, name))) })
	}
	pred("is.null", func(v Value) bool { _, ok := v.(nullT); return ok })
	pred("is.logical", func(v Value) bool { _, ok := v.(*Logical); return ok })
	pred("is.integer", func(v Value) bool { _, ok := v.(*Integer); return ok })
	pred("is.double", func(v Value) bool { _, ok := v.(*Double); return ok })
	pred("is.numeric", func(v Value) bool {
		switch v.(type) {
		case *Integer, *Double:
			return getAttr(v, "class") == nil
		}
		return false
	})
	pred("is.character", func(v Value) bool { _, ok := v.(*String); return ok })
	pred("is.list", func(v Value) bool { _, ok := v.(*List); return ok })
	pred("is.function", func(v Value) bool {
		switch v.(type) {
		case *Closure, *Builtin:
			return true
		}
		return false
	})
	pred("is.symbol", func(v Value) bool { _, ok := v.(*Symbol); return ok })
	pred("is.name", func(v Value) bool { _, ok := v.(*Symbol); return ok })
	pred("is.call", func(v Value) bool { _, ok := v.(*Lang); return ok })
	pred("is.environment", func(v Value) bool { _, ok := v.(*Env); return ok })
	pred("is.atomic", func(v Value) bool {
		switch v.(type) {
		case *Logical, *Integer, *Double, *String:
			return true
		}
		return false
	})
	regBuiltin("as.integer", func(s *Session, c *callInfo) Value { return s.asInteger(c) })
	regBuiltin("as.double", func(s *Session, c *callInfo) Value { return s.asDouble(c, "as.double") })
	regBuiltin("as.numeric", func(s *Session, c *callInfo) Value { return s.asDouble(c, "as.numeric") })
	regBuiltin("as.logical", func(s *Session, c *callInfo) Value { return s.asLogical(c) })
	regBuiltin("as.character", func(s *Session, c *callInfo) Value { return s.asCharacter(c) })
	regBuiltin("abs", (*Session)(nil).math1("abs", math.Abs, true))
	regBuiltin("sqrt", (*Session)(nil).math1("sqrt", math.Sqrt, false))
	regBuiltin("exp", (*Session)(nil).math1("exp", math.Exp, false))
	regBuiltin("floor", (*Session)(nil).math1("floor", math.Floor, false))
	regBuiltin("ceiling", (*Session)(nil).math1("ceiling", math.Ceil, false))
	for _, n := range []string{"sum", "prod", "min", "max"} {
		regBuiltin(n, (*Session)(nil).summary(n))
	}
	regBuiltin("any", (*Session)(nil).anyAll(true))
	regBuiltin("all", (*Session)(nil).anyAll(false))
	regBuiltin("seq_len", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "seq_len")
		if Length(x) == 0 {
			s.errorf(c.call, "argument of length 0")
		}
		if Length(x) > 1 {
			s.warningf(c.call, "first element used of 'length.out' argument")
		}
		var d []float64
		switch kindOf(x) {
		case kLogical, kInteger, kDouble:
			d = asDoubles(x)
		case kString:
			if f, ok := parseNumber(asStrings(x)[0]); ok {
				d = []float64{f}
			}
		}
		if len(d) == 0 || math.IsNaN(d[0]) || math.IsInf(d[0], 0) || d[0] < 0 {
			s.errorf(c.call, "argument must be coercible to non-negative integer")
		}
		if d[0] > 1<<52 {
			s.errorf(c.call, "result would be too long a vector")
		}
		n := int(d[0])
		s.checkLen(c.call, n)
		out := make([]int32, n)
		for i := range out {
			out[i] = int32(i + 1)
		}
		return &Integer{V: out}
	})
	regBuiltin("seq_along", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "seq_along")
		out := make([]int32, Length(x))
		for i := range out {
			out[i] = int32(i + 1)
		}
		return &Integer{V: out}
	})
	regBuiltin("names", func(s *Session, c *callInfo) Value {
		if n := namesOf(s.nargs1(c, "names")); n != nil {
			return n
		}
		return Null
	})
}

// quiet sets the quiet bit of a NaN, as an arithmetic operation on it does.
func quiet(x float64) float64 {
	return math.Float64frombits(math.Float64bits(x) | 0x0008000000000000)
}

// strExtreme is min or max when an argument is a character vector: all the
// arguments are compared as strings.
func (s *Session) strExtreme(c *callInfo, name string, data []Value, rm bool) Value {
	var best Str
	have, na := false, false
	for _, a := range data {
		switch a.(type) {
		case nullT, *Logical, *Integer, *Double, *String:
		default:
			s.errorf(c.call, "invalid 'type' (%s) of argument", typeName(a))
		}
		for _, e := range asStrings(a) {
			if e.NA {
				if !rm {
					na = true
				}
				continue
			}
			if !have || (name == "min" && e.S < best.S) || (name == "max" && e.S > best.S) {
				best, have = e, true
			}
		}
	}
	if na {
		return &String{V: []Str{{NA: true}}}
	}
	if !have {
		inf := map[string]string{"min": "Inf", "max": "-Inf"}[name]
		s.errorf(c.call, "no non-missing arguments to %s; returning %s", name, inf)
	}
	return &String{V: []Str{best}}
}
