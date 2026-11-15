package eval

import (
	"strconv"
	"strings"

	"github.com/go-embedded-r/r/syntax"
)

// model.frame.default (slice 17) is a Go native. Its first draft in R scored
// an exact run of 80 tokens and a normalised run of 80 against the model.frame
// of R's stats (limits 40 and 60); under the rule it was not reworked until it
// passed but written again from the help page and the oracle. The closure in
// formula17.R only forwards what R's lazy evaluation needs (the unevaluated
// subset and extra arguments); this file does the rest.

// mfClass is the class of a model-frame variable, as dataClasses records it.
func mfClass(v Value) string {
	switch {
	case kindOf(v) == kLogical:
		return "logical"
	case inheritsClass(v, "ordered"):
		return "ordered"
	case inheritsClass(v, "factor"):
		return "factor"
	case kindOf(v) == kString:
		return "character"
	}
	if d := dimsOf(v); len(d) == 2 && (kindOf(v) == kInteger || kindOf(v) == kDouble) {
		return "nmatrix." + strconv.Itoa(d[1])
	}
	if k := kindOf(v); k == kInteger || k == kDouble {
		return "numeric"
	}
	return "other"
}

func init() {
	// ..getAllVars(missFormula, missData, formula, data, extras, callerEnv)
	regBuiltin("..getAllVars", func(s *Session, c *callInfo) Value {
		flag := func(i int) bool { l := asInts(c.args[i]); return len(l) > 0 && l[0] == 1 }
		missFormula, missData := flag(0), flag(1)
		formula, data, extrasExpr := c.args[2], c.args[3], c.args[4]
		caller, _ := c.args[5].(*Env)
		if caller == nil {
			caller = c.env
		}
		callNamed := func(name string, args []Value, names []string) Value { return s.callNamed(c, name, args, names) }
		isFrame := func(v Value) bool { return inheritsClass(v, "data.frame") }
		switch {
		case missFormula:
			if !missData && isFrame(data) && getAttr(data, "terms") != nil {
				return data
			}
			formula = callNamed("as.formula", []Value{data}, []string{""})
		case missData && isFrame(formula):
			data = formula
			formula = callNamed("as.formula", []Value{data}, []string{""})
		}
		formula = callNamed("as.formula", []Value{formula}, []string{""})
		env, _ := getAttr(formula, ".Environment").(*Env)
		if env == nil {
			env = caller
		}
		var dataEnv Value = data
		if _, isNull := data.(nullT); isNull {
			dataEnv = env
		}
		e := s.evalEnvArg(c, dataEnv, env)
		var vars []string
		s.langVars(formula, false, &vars)
		seen := map[string]bool{}
		out := &List{}
		var names []Str
		for _, v := range vars {
			if seen[v] {
				continue
			}
			seen[v] = true
			val := s.underEval(callOf("FUN(X[[i]], ...)"), func() Value { return s.eval(Intern(v), e) })
			if _, isNull := val.(nullT); isNull {
				continue
			}
			out.V = append(out.V, val)
			names = append(names, Str{S: v})
		}
		if extrasExpr != nil {
			if ex, ok := s.underEval(callOf("eval(extras, data, env)"), func() Value { return s.eval(extrasExpr, e) }).(*List); ok {
				nm := namesOf(ex)
				for i, v := range ex.V {
					if _, isNull := v.(nullT); isNull {
						continue
					}
					out.V = append(out.V, v)
					if nm != nil && nm.V[i].S != "" {
						names = append(names, nm.V[i])
					} else {
						names = append(names, Str{NA: true})
					}
				}
			}
		}
		if len(out.V) == 0 {
			return callNamed("data.frame", nil, nil)
		}
		out.Attr = namesAttr(names)
		return callNamed("as.data.frame", []Value{out, mkLogical(false), mkLogical(true)}, []string{"", "stringsAsFactors", "optional"})
	})
	// model.response(data, type = "any"): written from the help page and the
	// oracle (its first draft in R was an exact run of 59 tokens of R's)
	regBuiltin("model.response", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "data", "type")
		if a[0] == nil {
			s.errorf(c.call, "argument \"data\" is missing, with no default")
		}
		data := a[0]
		typ := "any"
		if a[1] != nil {
			typ = s.stringArg(c, a[1], "type")
		}
		resp := getAttr(getAttr(data, "terms"), "response")
		if resp == nil || Length(resp) == 0 {
			s.errorf(c.call, "argument is of length zero")
		}
		if d := asDoubles(resp); len(d) > 0 && d[0] <= 0 {
			return Null
		}
		cols, ok := data.(*List)
		if !ok || len(cols.V) == 0 {
			s.errorf(c.call, "invalid data")
		}
		v := cols.V[0]
		switch {
		case typ == "numeric" && inheritsClass(v, "factor"):
			s.warningf(c.call, "using type = \"numeric\" with a factor response will be ignored")
		case typ == "numeric" || typ == "double":
			if k := kindOf(v); k == kLogical || k == kInteger {
				v = withAttr(&Double{V: asDoubles(v)}, attrsOf(v))
			}
		case typ != "any":
			s.errorf(c.call, "invalid response type")
		}
		rows := getAttr(data, "row.names")
		if rows != nil {
			rows = expandRowNames(rows)
		}
		d := dimsOf(v)
		switch {
		case len(d) == 2 && d[1] == 1:
			var keep []Attr
			if rn := dimnamesOf(v, 0); rn != nil {
				keep = namesAttr(rn.V)
			}
			return withAttr(v, keep)
		case len(d) == 2:
			if dimnamesOf(v, 0) == nil && rows != nil && Length(rows) == d[0] {
				var cn Value = Null
				if c2 := dimnamesOf(v, 1); c2 != nil {
					cn = c2
				}
				return setAttr(v, "dimnames", &List{V: []Value{&String{V: asStrings(rows)}, cn}})
			}
			return v
		}
		if rows != nil && Length(rows) == Length(v) {
			return setAttr(v, "names", &String{V: asStrings(rows)})
		}
		return v
	})
	regBuiltin(".MFclass", func(s *Session, c *callInfo) Value {
		if len(c.args) != 1 {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		return mkString(mfClass(c.args[0]))
	})
	// ..modelFrame(missFormula, missData, missNA, formula, data, subset, na.action,
	//              drop.unused.levels, xlev, extras, callerEnv)
	regBuiltin("..modelFrame", func(s *Session, c *callInfo) Value {
		flag := func(i int) bool { l := asInts(c.args[i]); return len(l) > 0 && l[0] == 1 }
		missFormula, missData, missNA := flag(0), flag(1), flag(2)
		formula, data := c.args[3], c.args[4]
		subsetExpr, naAction, dropUnused, xlev, extrasExpr := c.args[5], c.args[6], flag(7), c.args[8], c.args[9]
		caller, _ := c.args[10].(*Env)
		if caller == nil {
			caller = c.env
		}
		if missNA {
			found := getAttr(data, "na.action")
			if found != nil && kindOf(found) != kInteger && kindOf(found) != kDouble {
				naAction = found
			} else if o := s.options["na.action"]; o != nil {
				naAction = o
			} else {
				naAction = mkString("na.omit")
			}
		}
		callNamed := func(name string, args []Value, names []string) Value { return s.callNamed(c, name, args, names) }
		isFrame := func(v Value) bool { return inheritsClass(v, "data.frame") }
		switch {
		case missFormula:
			if !missData && isFrame(data) && getAttr(data, "terms") != nil {
				return data
			}
			formula = callNamed("as.formula", []Value{data}, []string{""})
		case missData && isFrame(formula):
			if getAttr(formula, "terms") != nil {
				return formula
			}
			data = formula
			formula = callNamed("as.formula", []Value{data}, []string{""})
		}
		if !inheritsClass(formula, "terms") {
			f := callNamed("as.formula", []Value{formula}, []string{""})
			formula = callNamed("terms", []Value{f, data}, []string{"", "data"})
		}
		tt := formula
		env, _ := getAttr(tt, ".Environment").(*Env)
		if env == nil {
			env = caller
		}
		var dataEnv Value = data
		switch d := data.(type) {
		case nullT:
			dataEnv = env
		case *Env:
		default:
			if len(dimsOf(d)) >= 2 {
				s.errorf(c.call, "'data' must be a data.frame, not a matrix or an array")
			}
			if isObject(d) && !isFrame(d) {
				dataEnv = callNamed("as.data.frame", []Value{d}, []string{""})
			}
		}
		evalIn := func(expr Value, shown string) Value {
			e := s.evalEnvArg(c, dataEnv, env)
			return s.underEval(callOf(shown), func() Value { return s.eval(expr, e) })
		}
		vars, _ := getAttr(tt, "variables").(*Lang)
		if vars == nil {
			s.errorf(c.call, "invalid terms object")
		}
		cols, _ := evalIn(vars, "eval(predvars, data, env)").(*List)
		if cols == nil {
			s.errorf(c.call, "invalid variables")
		}
		var vals []Value
		var names []Str
		for i := range vars.Args {
			vals = append(vals, cols.V[i])
			names = append(names, Str{S: strings.Join(syntax.Deparse(s.toSyntax(vars.Args[i].Value), &syntax.DeparseOptions{Cutoff: 500}), " ")})
		}
		if extrasExpr != nil {
			if ex, ok := evalIn(extrasExpr, "eval(extras, data, env)").(*List); ok {
				nm := namesOf(ex)
				for i, e := range ex.V {
					if _, isNull := e.(nullT); isNull || nm == nil || nm.V[i].S == "" {
						continue
					}
					vals = append(vals, e)
					names = append(names, Str{S: "(" + nm.V[i].S + ")"})
				}
			}
		}
		rowsOf := func(v Value) int {
			if d := dimsOf(v); len(d) >= 1 {
				return d[0]
			}
			return Length(v)
		}
		for i, v := range vals {
			switch k := kindOf(v); {
			case k == kLogical || k == kInteger || k == kDouble || k == kString || k == kComplex || k == kRaw:
			default:
				s.errorf(c.call, "invalid type (%s) for variable '%s'", typeName(v), names[i].S)
			}
		}
		nr := 0
		if len(vals) > 0 {
			nr = rowsOf(vals[0])
		} else if isFrame(data) {
			if rn := expandRowNames(getAttr(data, "row.names")); rn != nil {
				nr = Length(rn)
			}
		}
		for i, v := range vals {
			if rowsOf(v) != nr {
				s.errorf(c.call, "variable lengths differ (found for '%s')", names[i].S)
			}
		}
		// the row names: the data frame's, or the names of the first variable
		var rn Value
		if isFrame(data) {
			if r := expandRowNames(getAttr(data, "row.names")); r != nil && Length(r) == nr {
				rn = r
			}
		}
		if rn == nil {
			if len(vals) > 0 {
				var lab *String
				if d := dimsOf(vals[0]); len(d) == 2 {
					lab = dimnamesOf(vals[0], 0)
				} else {
					lab = namesOf(vals[0])
				}
				if lab != nil && len(lab.V) == nr {
					seen := map[string]bool{}
					dup := false
					for _, e := range lab.V {
						dup = dup || seen[e.S]
						seen[e.S] = true
					}
					if !dup {
						rn = lab
					}
				}
			}
			if rn == nil {
				seq := make([]int32, nr)
				for i := range seq {
					seq[i] = int32(i + 1)
				}
				rn = &Integer{V: seq}
			}
		}
		keep := make([]int32, nr)
		for i := range keep {
			keep[i] = int32(i + 1)
		}
		if _, isNull := subsetExpr.(nullT); !isNull && subsetExpr != nil {
			sub := evalIn(subsetExpr, "eval(substitute(subset), data, env)")
			if _, isNull := sub.(nullT); !isNull {
				k := callNamed("[", []Value{&Integer{V: keep}, sub}, []string{"", ""})
				var kk []int32
				for _, v := range asInts(k) {
					if v != naInt {
						kk = append(kk, v)
					}
				}
				keep = kk
			}
		}
		whole := len(keep) == nr
		for i, v := range keep {
			whole = whole && int(v) == i+1
		}
		// the levels the model was fitted with
		if xl, ok := xlev.(*List); ok {
			xn := namesOf(xl)
			for i, lv := range xl.V {
				if xn == nil {
					break
				}
				for j := range names {
					if names[j].S != xn.V[i].S || !(inheritsClass(vals[j], "factor") || kindOf(vals[j]) == kString) {
						continue
					}
					var sv []Str
					if inheritsClass(vals[j], "factor") {
						sv = s.factorStrings(vals[j])
					} else {
						sv = asStrings(vals[j])
					}
					known := map[string]bool{}
					for _, e := range asStrings(lv) {
						known[e.S] = true
					}
					var bad []string
					for _, e := range sv {
						if !e.NA && !known[e.S] {
							bad = append(bad, e.S)
						}
					}
					if len(bad) > 0 {
						s.errorf(c.call, "factor %s has new levels %s", xn.V[i].S, strings.Join(uniqueStrings(bad), ", "))
					}
					vals[j] = callNamed("factor", []Value{&String{V: sv}, lv}, []string{"", "levels"})
				}
			}
		}
		out := &List{V: make([]Value, len(vals))}
		classes := make([]Str, len(vals))
		for j, v := range vals {
			classes[j] = Str{S: mfClass(v)}
			idx := &Integer{V: keep}
			if len(dimsOf(v)) == 2 {
				out.V[j] = callNamed("[", []Value{v, idx, Missing, mkLogical(false)}, []string{"", "", "", "drop"})
			} else {
				out.V[j] = callNamed("[", []Value{v, idx}, []string{"", ""})
			}
		}
		out.Attr = namesAttr(names)
		rnKept := callNamed("[", []Value{rn, &Integer{V: keep}}, []string{"", ""})
		auto := false
		if iv, ok := rn.(*Integer); ok && whole && namesOf(rn) == nil {
			auto = true
			for i, e := range iv.V {
				auto = auto && int(e) == i+1
			}
		}
		// dataClasses and predvars are added to the terms object
		tt = setAttr(tt, "predvars", getAttr(tt, "variables"))
		cls := &String{V: classes}
		cls.Attr = namesAttr(names)
		tt = setAttr(tt, "dataClasses", cls)
		rnAttr := func() Value {
			if auto {
				return &Integer{V: []int32{naInt, int32(len(keep))}}
			}
			return rnKept
		}
		setRN := func(v Value, rn Value) Value {
			return &List{V: v.(*List).V, Attr: append(append([]Attr(nil), attrsOf(v)...), Attr{Name: "row.names", Value: compactRowNames(rn)})}
		}
		var frame Value = out
		if whole {
			frame = &List{V: out.V, Attr: append(append([]Attr(nil), out.Attr...), Attr{Name: "class", Value: mkString("data.frame")}, Attr{Name: "row.names", Value: rnAttr()})}
		} else {
			frame = setRN(out, rnKept)
			frame = setAttr(frame, "class", mkString("data.frame"))
		}
		frame = setAttr(frame, "terms", tt)
		if _, isNull := naAction.(nullT); !isNull && naAction != nil {
			fn := naAction
			if st, ok := fn.(*String); ok && len(st.V) == 1 {
				fn = s.findMethod(st.V[0].S, c.env)
				if fn == nil {
					s.errorf(c.call, "object '%s' of mode 'function' was not found", st.V[0].S)
				}
			}
			kept := frame
			frame = s.callFunction(fn, []Value{frame})
			if !s.identical(frame, kept) {
				// what the method answered gets the terms again, in R's order
				act := getAttr(frame, "na.action")
				rows := getAttr(frame, "row.names")
				fl, _ := frame.(*List)
				if fl == nil {
					s.errorf(c.call, "invalid result of na.action")
				}
				nattr := []Attr{{Name: "names", Value: getAttr(frame, "names")}, {Name: "terms", Value: tt}}
				if rows != nil {
					nattr = append(nattr, Attr{Name: "row.names", Value: rows})
				}
				nattr = append(nattr, Attr{Name: "class", Value: mkString("data.frame")})
				if act != nil {
					nattr = append(nattr, Attr{Name: "na.action", Value: act})
				}
				frame = &List{V: fl.V, Attr: nattr}
			}
		}
		if dropUnused {
			fl := frame.(*List)
			nl := &List{V: append([]Value(nil), fl.V...), Attr: fl.Attr}
			for i, col := range nl.V {
				if !inheritsClass(col, "factor") {
					continue
				}
				lev, _ := getAttr(col, "levels").(*String)
				if lev == nil {
					continue
				}
				used := make([]bool, len(lev.V))
				for _, code := range asInts(col) {
					if code != naInt && int(code) >= 1 && int(code) <= len(lev.V) {
						used[code-1] = true
					}
				}
				remap := make([]int32, len(lev.V))
				var nlev []Str
				for j, u := range used {
					if u {
						nlev = append(nlev, lev.V[j])
						remap[j] = int32(len(nlev))
					}
				}
				if len(nlev) == len(lev.V) {
					continue
				}
				codes := make([]int32, Length(col))
				for k, code := range asInts(col) {
					if code == naInt {
						codes[k] = naInt
					} else {
						codes[k] = remap[code-1]
					}
				}
				var attrs []Attr
				for _, a := range attrsOf(col) {
					if a.Name == "levels" {
						a.Value = &String{V: nlev}
					}
					attrs = append(attrs, a)
				}
				nl.V[i] = &Integer{V: codes, Attr: attrs}
			}
			frame = nl
		}
		return frame
	})
}

func uniqueStrings(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, e := range in {
		if !seen[e] {
			seen[e] = true
			out = append(out, e)
		}
	}
	return out
}
