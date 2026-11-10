package syntax

import (
	"math"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// DeparseOptions are the options of deparse() that matter for syntax trees:
// the default width cutoff is 60 bytes, names are backquoted when they are
// not syntactic, and integers keep their L suffix (R's default "keepInteger",
// "keepNA").
type DeparseOptions struct {
	// Cutoff is the line width, in bytes, after which a line is broken at the
	// next break point (default 60).
	Cutoff int
	// NoBackquote turns off the backquoting of non-syntactic names, as
	// deparse(backtick = FALSE) does (R turns it off for a bare symbol).
	NoBackquote bool
	// NoKeepInteger drops the L suffix and NA_integer_/NA_real_ spellings.
	NoKeepInteger bool
}

// Deparse writes n as deparse() does, one string per line, byte for byte.
// The text is written from the documented behaviour of deparse() and
// checked against the oracle (cmd/rparity).
func Deparse(n Node, opt *DeparseOptions) []string {
	d := &deparser{cutoff: 60, backtick: true, keepInt: true, startline: true}
	if opt != nil {
		if opt.Cutoff > 0 {
			d.cutoff = opt.Cutoff
		}
		d.backtick = !opt.NoBackquote
		d.keepInt = !opt.NoKeepInteger
	}
	d.node(n)
	d.writeline()
	return d.lines
}

type ppKind int

const (
	ppFuncall ppKind = iota
	ppIf
	ppWhile
	ppFor
	ppRepeat
	ppBreak
	ppNext
	ppReturn
	ppFunction
	ppAssign
	ppAssign2
	ppCurly
	ppParen
	ppSubset
	ppDollar
	ppSubass
	ppBinary
	ppBinary2
	ppUnary
	ppForeign
)

// Operator precedences, low to high (R Internals, deparse).
const (
	precFn = iota
	precEq
	precLeft
	precRight
	precTilde
	precOr
	precAnd
	precNot
	precCompare
	precSum
	precProd
	precPercent
	precColon
	precSign
	precPower
	precSubset
	precDollar
	precNS
)

type ppInfo struct {
	kind  ppKind
	prec  int
	right bool
}

var ppTable = map[string]ppInfo{
	"if": {ppIf, precFn, true}, "while": {ppWhile, precFn, false}, "for": {ppFor, precFn, false},
	"repeat": {ppRepeat, precFn, false}, "break": {ppBreak, precFn, false}, "next": {ppNext, precFn, false},
	"return": {ppReturn, precFn, false}, "function": {ppFunction, precFn, false},
	"<-": {ppAssign, precLeft, true}, "=": {ppAssign, precEq, true}, "<<-": {ppAssign2, precLeft, true},
	"{": {ppCurly, precFn, false}, "(": {ppParen, precFn, false},
	"[": {ppSubset, precSubset, false}, "[[": {ppSubset, precSubset, false},
	"$": {ppDollar, precDollar, false}, "@": {ppDollar, precDollar, false},
	"[<-": {ppSubass, precLeft, true}, "[[<-": {ppSubass, precLeft, true}, "$<-": {ppSubass, precLeft, true}, "@<-": {ppSubass, precLeft, true},
	"+": {ppBinary, precSum, false}, "-": {ppBinary, precSum, false}, "*": {ppBinary, precProd, false},
	"/": {ppBinary2, precProd, false}, "^": {ppBinary2, precPower, true},
	"%%": {ppBinary2, precPercent, false}, "%/%": {ppBinary2, precPercent, false}, "%*%": {ppBinary, precPercent, false},
	"==": {ppBinary, precCompare, false}, "!=": {ppBinary, precCompare, false}, "<": {ppBinary, precCompare, false},
	"<=": {ppBinary, precCompare, false}, ">=": {ppBinary, precCompare, false}, ">": {ppBinary, precCompare, false},
	"&": {ppBinary, precAnd, false}, "|": {ppBinary, precOr, false}, "!": {ppUnary, precNot, false},
	"&&": {ppBinary, precAnd, false}, "||": {ppBinary, precOr, false},
	":": {ppBinary2, precColon, false}, "~": {ppBinary, precTilde, false},
	"::": {ppBinary2, precNS, false}, ":::": {ppBinary2, precNS, false},
	".C": {ppForeign, precFn, false}, ".Fortran": {ppForeign, precFn, false}, ".External": {ppForeign, precFn, false},
	".External2": {ppForeign, precFn, false}, ".Call": {ppForeign, precFn, false},
	".External.graphics": {ppForeign, precFn, false}, ".Call.graphics": {ppForeign, precFn, false},
	"is.loaded": {ppForeign, precFn, false}, "recordGraphics": {ppForeign, precFn, false},
}

func isUserBinop(name string) bool {
	return len(name) >= 2 && name[0] == '%' && name[len(name)-1] == '%'
}

// isBuiltin reports whether name is bound to a primitive in R's base
// environment (generated list, see prims_gen.go).
func isBuiltin(name string) bool {
	if _, ok := ppTable[name]; ok {
		return true
	}
	return primitiveNames[name]
}

type deparser struct {
	cutoff   int
	backtick bool
	keepInt  bool

	lines     []string
	buf       []byte
	indent    int
	incurly   int
	inlist    int
	startline bool
	left      int
	fnarg     bool
	elems     bool // the next args() are the elements of an expression vector
}

func (d *deparser) writeline() {
	d.lines = append(d.lines, string(d.buf))
	d.buf = d.buf[:0]
	d.startline = true
}

func (d *deparser) print(s string) {
	if d.startline {
		d.startline = false
		for i := 1; i <= d.indent; i++ {
			if i <= 4 {
				d.buf = append(d.buf, "    "...)
			} else {
				d.buf = append(d.buf, "  "...)
			}
		}
	}
	d.buf = append(d.buf, s...)
}

func (d *deparser) linebreak(lbreak *bool) {
	if len(d.buf) > d.cutoff {
		if !*lbreak {
			*lbreak = true
			d.indent++
		}
		d.writeline()
	}
}

// ---- names and strings ----------------------------------------------------

var reservedWords = map[string]bool{
	"NULL": true, "NA": true, "TRUE": true, "FALSE": true, "Inf": true, "NaN": true,
	"NA_integer_": true, "NA_real_": true, "NA_character_": true, "NA_complex_": true,
	"function": true, "while": true, "repeat": true, "for": true, "if": true, "in": true,
	"else": true, "next": true, "break": true,
}

// isValidName is R's isValidName: a syntactic name that is not reserved.
func isValidName(s string) bool {
	if s == "" {
		return false
	}
	r, n := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && n <= 1 {
		return false
	}
	if r != '.' && !unicode.IsLetter(r) {
		return false
	}
	if r == '.' && len(s) > 1 && s[1] >= '0' && s[1] <= '9' {
		return false
	}
	for _, c := range s[n:] {
		if !(unicode.IsLetter(c) || unicode.IsDigit(c) || c == '.' || c == '_') {
			return false
		}
	}
	if s == "..." {
		return true
	}
	return !reservedWords[s]
}

func (d *deparser) quotify(name string, quote byte) string {
	if isValidName(name) || name == "" {
		return name
	}
	return encodeString(name, quote)
}

// EncodeString writes s in the given quote with R's escapes, as print() and
// deparse() do.
func EncodeString(s string, quote byte) string { return encodeString(s, quote) }

// DeparseClosure writes a function the way print() shows one without source
// references: "function (x, y = 2) " and then the body on the next lines.
func DeparseClosure(params []Param, body Node) []string {
	d := &deparser{cutoff: 60, backtick: true, keepInt: true, startline: true}
	d.print("function (")
	d.formals(params)
	d.print(") ")
	d.writeline()
	d.node(body)
	d.writeline()
	return d.lines
}

// encodeString writes s in double (or other) quotes with R's escapes.
func encodeString(s string, quote byte) string {
	var b strings.Builder
	b.WriteByte(quote)
	for i := 0; i < len(s); {
		c := s[i]
		if c < 0x80 {
			i++
			switch c {
			case '\\':
				b.WriteString(`\\`)
			case '\a':
				b.WriteString(`\a`)
			case '\b':
				b.WriteString(`\b`)
			case '\f':
				b.WriteString(`\f`)
			case '\n':
				b.WriteString(`\n`)
			case '\r':
				b.WriteString(`\r`)
			case '\t':
				b.WriteString(`\t`)
			case '\v':
				b.WriteString(`\v`)
			case 0:
				b.WriteString(`\0`)
			default:
				switch {
				case c == quote:
					b.WriteByte('\\')
					b.WriteByte(c)
				case c < 0x20 || c == 0x7f:
					b.WriteString(`\` + octal3(c))
				default:
					b.WriteByte(c)
				}
			}
			continue
		}
		r, n := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && n == 1 {
			b.WriteString(`\x` + hex2(c))
			i++
			continue
		}
		if rPrintable(r) {
			b.WriteString(s[i : i+n])
		} else if r > 0xffff {
			b.WriteString(`\U` + hexN(uint32(r), 8))
		} else {
			b.WriteString(`\u` + hexN(uint32(r), 4))
		}
		i += n
	}
	b.WriteByte(quote)
	return b.String()
}

func octal3(c byte) string {
	return string([]byte{'0' + c>>6, '0' + (c>>3)&7, '0' + c&7})
}

func hex2(c byte) string { return hexN(uint32(c), 2) }

func hexN(v uint32, w int) string {
	s := strconv.FormatUint(uint64(v), 16)
	for len(s) < w {
		s = "0" + s
	}
	return s
}

// ---- numbers -------------------------------------------------------------

// pow10tbl are the powers of ten exact in a double.
var pow10tbl = [...]float64{1e00, 1e01, 1e02, 1e03, 1e04, 1e05, 1e06, 1e07, 1e08, 1e09, 1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19, 1e20, 1e21, 1e22}

// scientific finds the decimal exponent and the number of significant digits
// (at most 15) of |x| the way R's formatReal does with 15 digits, in double
// precision (R's computation, not the exact decimal expansion: they differ
// for some very large and very small values).
func scientific(x float64) (neg, kpower, nsig int, roundingWidens bool) {
	return scientificD(x, 15)
}

// Scientific is R's scientific() for the given number of significant digits
// (7 for print(), 15 for deparse()): the sign, the decimal exponent, the
// number of significant digits needed (at most digits), and whether rounding
// to digits makes the number one decimal place wider than it is.
func Scientific(x float64, digits int) (neg, kpower, nsig int, roundingWidens bool) {
	return scientificD(x, digits)
}

func scientificD(x float64, digits int) (neg, kpower, nsig int, roundingWidens bool) {
	if x == 0 {
		return 0, 0, 1, false
	}
	r := x
	if x < 0 {
		neg, r = 1, -x
	}
	kp := int(math.Floor(math.Log10(r))) - digits + 1
	rp := r
	switch {
	case kp >= -22 && kp <= 22:
		if kp >= 0 {
			rp /= pow10tbl[kp]
		} else {
			rp *= pow10tbl[-kp]
		}
	case kp <= -308:
		rp = (rp * 1e303) / rexp10(kp+303)
	default:
		rp /= rexp10(kp)
	}
	if rp < pow10tbl[digits-1] {
		rp *= 10
		kp--
	}
	alpha := math.RoundToEven(rp)
	nsig = digits
	for j := 1; j <= digits; j++ {
		alpha /= 10
		if alpha == math.Floor(alpha) {
			nsig--
		} else {
			break
		}
	}
	if nsig == 0 {
		nsig = 1
		kp++
	}
	kpower = kp + digits - 1
	rgt := digits - kpower
	if rgt < 0 {
		rgt = 0
	} else if rgt > 22 {
		rgt = 22
	}
	fuzz := 0.5 / pow10tbl[rgt]
	roundingWidens = kpower > 0 && kpower <= 22 && r < pow10tbl[kpower]-fuzz
	return
}

// formatReal15 is formatReal + EncodeReal for one value at 15 significant
// digits: fixed notation unless scientific is narrower.
func formatReal15(x float64) string {
	switch {
	case math.IsNaN(x):
		if math.Float64bits(x)&0xFFFFFFFF == 1954 {
			return "NA"
		}
		return "NaN"
	case math.IsInf(x, 1):
		return "Inf"
	case math.IsInf(x, -1):
		return "-Inf"
	}
	if x == 0 {
		x = 0 // -0 is written 0
	}
	neg, kpower, nsig, rw := scientific(x)
	left := kpower + 1
	if rw {
		left--
	}
	sleft := neg
	if left <= 0 {
		sleft++
	} else {
		sleft += left
	}
	rgt := nsig - left
	if rgt < 0 {
		rgt = 0
	}
	if left < 0 {
		sleft = 1 + neg
	}
	wF := sleft + rgt
	if rgt != 0 {
		wF++
	}
	e := 1
	if left > 100 || left <= -99 {
		e = 2
	}
	dsci := nsig - 1
	wE := neg + dsci + 4 + e
	if dsci > 0 {
		wE++
	}
	if wF <= wE {
		return strconv.FormatFloat(x, 'f', rgt, 64)
	}
	return strconv.FormatFloat(x, 'e', dsci, 64)
}

func (d *deparser) constant(c *Const) {
	switch c.Kind {
	case Null:
		d.print("NULL")
	case Logical:
		switch {
		case c.NA:
			d.print("NA")
		case c.LogicalV:
			d.print("TRUE")
		default:
			d.print("FALSE")
		}
	case Integer:
		switch {
		case c.NA && d.keepInt:
			d.print("NA_integer_")
		case c.NA:
			d.print("NA")
		case d.keepInt:
			d.print(strconv.Itoa(int(c.Int)) + "L")
		default:
			d.print(strconv.Itoa(int(c.Int)))
		}
	case Double:
		if c.NA && d.keepInt {
			d.print("NA_real_")
		} else {
			d.print(formatReal15(c.Double))
		}
	case Complex:
		d.print(formatComplex15(c))
	case Raw:
		d.print("0x" + strconv.FormatInt(int64(c.Int)|0x100, 16)[1:])
	case String:
		if c.NA {
			d.print("NA_character_")
		} else {
			d.print(encodeString(c.Str, '"'))
		}
	}
}

func formatComplex15(c *Const) string {
	if c.NA {
		return "NA_complex_"
	}
	re, im := c.Double, c.Im
	if math.IsNaN(re) || math.IsInf(im, 0) || math.IsNaN(im) {
		return "complex(real=" + formatReal15(re) + ", imaginary=" + formatReal15(im) + ")"
	}
	sign := "+"
	if im < 0 {
		sign = "-"
		im = -im
	}
	return formatReal15(re) + sign + formatReal15(im) + "i"
}

// ---- the recursive part --------------------------------------------------

func (d *deparser) node(n Node) {
	fnarg := d.fnarg
	d.fnarg = false
	prevLeft := d.left
	d.left = 0
	switch x := n.(type) {
	case nil:
		// R_MissingArg, the empty symbol: it prints nothing, but printing
		// it starts a pending line, which writes the indentation
		d.print("")
	case *Sym:
		if d.backtick {
			d.print(d.quotify(x.Name, '`'))
		} else {
			d.print(x.Name)
		}
	case *Const:
		d.constant(x)
	case *Function:
		d.print("function(")
		d.formals(x.Params)
		d.print(") ")
		d.node(x.Body)
	case *Call:
		d.call(x, fnarg, prevLeft)
	}
	d.left = prevLeft
}

func (d *deparser) formals(ps []Param) {
	lbreak := false
	for i, p := range ps {
		if p.Name == "..." {
			d.print(p.Name)
		} else if d.backtick {
			d.print(d.quotify(p.Name, '`'))
		} else {
			d.print(d.quotify(p.Name, '"'))
		}
		if p.Default != nil {
			d.print(" = ")
			d.fnarg = true
			d.node(p.Default)
		}
		if i < len(ps)-1 {
			d.print(", ")
			d.linebreak(&lbreak)
		}
	}
	if lbreak {
		d.indent--
	}
}

func (d *deparser) args(args []Arg) {
	elems := d.elems
	d.elems = false
	lbreak := false
	for i, a := range args {
		if a.HasName {
			if a.Name == "..." {
				d.print(a.Name)
			} else if d.backtick {
				d.print(d.quotify(a.Name, '`'))
			} else {
				d.print(d.quotify(a.Name, '"'))
			}
			d.print(" = ")
			if a.Value != nil {
				d.fnarg = !elems
				d.node(a.Value)
			}
		} else {
			d.fnarg = !elems
			d.node(a.Value)
		}
		if i < len(args)-1 {
			d.print(", ")
			d.linebreak(&lbreak)
		}
	}
	if lbreak {
		d.indent--
	}
}

func curlyAhead(n Node) bool {
	c, ok := n.(*Call)
	return ok && c.FunName() == "{"
}

// info is the printing info of the call's operator, or ok=false when the
// callee is not a builtin (or %op%) and prints as an ordinary call.
func callInfo(c *Call) (name string, fop ppInfo, user, ok bool) {
	s, isSym := c.Fun.(*Sym)
	if !isSym {
		return "", ppInfo{}, false, false
	}
	name = s.Name
	if info, found := ppTable[name]; found {
		return name, info, false, true
	}
	if isUserBinop(name) {
		for _, a := range c.Args {
			if a.HasName {
				return name, ppInfo{kind: ppFuncall}, true, true
			}
		}
		return name, ppInfo{ppBinary2, precPercent, false}, true, true
	}
	if primitiveNames[name] {
		return name, ppInfo{kind: ppFuncall}, false, true
	}
	return name, ppInfo{}, false, false
}

// needsParens decides whether arg, an operand of mainop, is bracketed.
func needsParens(mainop ppInfo, arg Node, left bool, deepLeft int) bool {
	if cc, ok := arg.(*Const); ok {
		if cc.Kind == Complex {
			return mainop.prec > precSum || (mainop.prec == precSum && left == mainop.right)
		}
		return false
	}
	c, ok := arg.(*Call)
	if !ok {
		if _, isFn := arg.(*Function); isFn {
			return false // PP_FUNCTION: default case
		}
		return false
	}
	name, ainfo, user, known := callInfo(c)
	if !known {
		return false
	}
	if user && ainfo.kind == ppFuncall {
		// %op% with named arguments is a function call; the isUserBinop
		// branch of needsparens only applies to symbols that are not builtins
		return false
	}
	if user {
		if mainop.prec > precPercent || (mainop.prec == precPercent && left == mainop.right) {
			return true
		}
		return false
	}
	_ = name
	switch ainfo.kind {
	case ppBinary, ppBinary2:
		switch len(c.Args) {
		case 1:
			if ainfo.prec == precSum {
				ainfo.prec = precSign
			}
			ainfo.kind = ppUnary
		case 2:
		default:
			return false
		}
	}
	cmp := func() bool {
		return mainop.prec > ainfo.prec || (mainop.prec == ainfo.prec && left == mainop.right)
	}
	switch ainfo.kind {
	case ppSubset:
		if mainop.kind == ppDollar || mainop.kind == ppSubset {
			if mainop.prec > ainfo.prec {
				return false
			}
		}
		fallthrough
	case ppBinary, ppBinary2:
		if mainop.prec == precCompare && ainfo.prec == precCompare {
			return true
		}
		return cmp()
	case ppAssign, ppAssign2, ppDollar:
		return cmp()
	case ppUnary:
		return (left && mainop.prec > ainfo.prec) || (deepLeft != 0 && deepLeft > ainfo.prec)
	case ppFor, ppIf, ppWhile, ppRepeat:
		return left || deepLeft != 0
	}
	return false
}

// parenthesizeCaller: does the function position of a call need brackets?
func parenthesizeCaller(op Node) bool {
	switch x := op.(type) {
	case *Function:
		return true
	case *Call:
		s, ok := x.Fun.(*Sym)
		if !ok {
			return true
		}
		if isUserBinop(s.Name) {
			return true
		}
		if info, found := ppTable[s.Name]; found {
			if info.prec >= precSubset || info.kind == ppFuncall || info.kind == ppParen || info.kind == ppCurly {
				return false
			}
			return true
		}
		return false
	}
	return false
}

func (d *deparser) operand(fop ppInfo, n Node, left bool, prevLeft int, ctx func()) {
	parens := needsParens(fop, n, left, prevLeft)
	if parens {
		d.print("(")
	}
	ctx()
	if parens {
		d.print(")")
	}
}

func (d *deparser) call(c *Call, fnarg bool, prevLeft int) {
	name, fop, user, builtin := callInfo(c)
	if !builtin {
		// an ordinary call, or one through a computed function
		if s, ok := c.Fun.(*Sym); ok {
			d.print(d.quotify(s.Name, '`'))
		} else {
			if parenthesizeCaller(c.Fun) {
				d.print("(")
				d.node(c.Fun)
				d.print(")")
			} else {
				d.node(c.Fun)
			}
		}
		d.print("(")
		d.elems = c.Elems
		d.args(c.Args)
		d.print(")")
		return
	}
	n := len(c.Args)
	argv := func(i int) Node {
		if i < n {
			return c.Args[i].Value
		}
		// CAR of the empty list is NULL
		return &Const{Kind: Null}
	}
	_ = user
	switch fop.kind {
	case ppBinary:
		switch n {
		case 1:
			fop.kind = ppUnary
			if fop.prec == precSum {
				fop.prec = precSign
			}
		case 2:
		default:
			fop.kind = ppFuncall
		}
	case ppBinary2:
		if n != 2 {
			fop.kind = ppFuncall
		} else if user {
			fop.kind = ppBinary
		}
	case ppDollar:
		if n != 2 {
			fop.kind = ppFuncall
			break
		}
		switch r := argv(1).(type) {
		case *Sym:
		case *Const:
			if !(r.Kind == String && !r.NA && r.Str != "") {
				fop.kind = ppFuncall
			}
		default:
			fop.kind = ppFuncall
		}
	}
	lookahead, lbreak := false, false
	switch fop.kind {
	case ppIf:
		d.print("if (")
		d.node(argv(0))
		d.print(") ")
		if d.incurly > 0 && d.inlist == 0 {
			lookahead = curlyAhead(argv(1))
			if !lookahead {
				d.writeline()
				d.indent++
			}
		}
		if n > 2 {
			d.node(argv(1))
			if d.incurly > 0 && d.inlist == 0 {
				d.writeline()
				if !lookahead {
					d.indent--
				}
			} else {
				d.print(" ")
			}
			d.print("else ")
			d.node(argv(2))
		} else {
			d.node(argv(1))
			if d.incurly > 0 && !lookahead && d.inlist == 0 {
				d.indent--
			}
		}
	case ppWhile:
		d.print("while (")
		d.node(argv(0))
		d.print(") ")
		d.node(argv(1))
	case ppFor:
		d.print("for (")
		d.node(argv(0))
		d.print(" in ")
		d.node(argv(1))
		d.print(") ")
		d.node(argv(2))
	case ppRepeat:
		d.print("repeat ")
		d.node(argv(0))
	case ppCurly:
		d.print("{")
		d.incurly++
		d.indent++
		d.writeline()
		for _, a := range c.Args {
			d.node(a.Value)
			d.writeline()
		}
		d.indent--
		d.print("}")
		d.incurly--
	case ppParen:
		d.print("(")
		d.node(argv(0))
		d.print(")")
	case ppSubset:
		d.operand(fop, argv(0), true, prevLeft, func() { d.node(argv(0)) })
		if name == "[" {
			d.print("[")
		} else {
			d.print("[[")
		}
		if n > 0 {
			d.args(c.Args[1:])
		}
		if name == "[" {
			d.print("]")
		} else {
			d.print("]]")
		}
	case ppFuncall, ppReturn:
		if d.backtick {
			d.print(d.quotify(name, '`'))
		} else {
			d.print(d.quotify(name, '"'))
		}
		d.print("(")
		d.inlist++
		d.elems = c.Elems
		d.args(c.Args)
		d.inlist--
		d.print(")")
	case ppForeign:
		d.print(name)
		d.print("(")
		d.inlist++
		d.args(c.Args)
		d.inlist--
		d.print(")")
	case ppFunction:
		// a call to `function` built by hand: formals are not representable
		d.print(name + "(")
		d.print(") ")
		d.node(argv(1))
	case ppAssign, ppAssign2:
		outer := fnarg && name == "="
		if outer {
			d.print("(")
		}
		parens := needsParens(fop, argv(0), true, prevLeft)
		if parens {
			d.print("(")
		}
		if parens {
			d.left = 0
		} else {
			d.left = fop.prec
		}
		d.node(argv(0))
		if parens {
			d.print(")")
		}
		d.print(" " + name + " ")
		parens = needsParens(fop, argv(1), false, prevLeft)
		if parens {
			d.print("(")
		}
		if parens {
			d.left = 0
		} else {
			d.left = prevLeft
		}
		d.node(argv(1))
		if parens {
			d.print(")")
		}
		if outer {
			d.print(")")
		}
		d.left = 0
	case ppDollar:
		parens := needsParens(fop, argv(0), true, prevLeft)
		if parens {
			d.print("(")
		}
		if parens {
			d.left = 0
		} else {
			d.left = fop.prec
		}
		d.node(argv(0))
		if parens {
			d.print(")")
		}
		d.print(name)
		if s, ok := argv(1).(*Const); ok && s.Kind == String && !s.NA && isValidName(s.Str) {
			d.print(s.Str)
		} else {
			parens = needsParens(fop, argv(1), false, prevLeft)
			if parens {
				d.print("(")
			}
			if parens {
				d.left = 0
			} else {
				d.left = prevLeft
			}
			d.node(argv(1))
			if parens {
				d.print(")")
			}
		}
		d.left = 0
	case ppBinary:
		parens := needsParens(fop, argv(0), true, prevLeft)
		if parens {
			d.print("(")
		}
		if parens {
			d.left = 0
		} else {
			d.left = fop.prec
		}
		d.node(argv(0))
		if parens {
			d.print(")")
		}
		d.print(" " + name + " ")
		d.linebreak(&lbreak)
		parens = needsParens(fop, argv(1), false, prevLeft)
		if parens {
			d.print("(")
		}
		if parens {
			d.left = 0
		} else {
			d.left = prevLeft
		}
		d.node(argv(1))
		if parens {
			d.print(")")
		}
		if lbreak {
			d.indent--
		}
		d.left = 0
	case ppBinary2:
		parens := needsParens(fop, argv(0), true, prevLeft)
		if parens {
			d.print("(")
		}
		if parens {
			d.left = 0
		} else {
			d.left = fop.prec
		}
		d.node(argv(0))
		if parens {
			d.print(")")
		}
		d.print(name)
		parens = needsParens(fop, argv(1), false, prevLeft)
		if parens {
			d.print("(")
		}
		if parens {
			d.left = 0
		} else {
			d.left = prevLeft
		}
		d.node(argv(1))
		if parens {
			d.print(")")
		}
		d.left = 0
	case ppUnary:
		d.print(name)
		parens := needsParens(fop, argv(0), false, prevLeft)
		if parens {
			d.print("(")
		}
		if parens {
			d.left = 0
		} else {
			d.left = prevLeft
		}
		d.node(argv(0))
		if parens {
			d.print(")")
		}
		d.left = 0
	case ppBreak:
		d.print("break")
	case ppNext:
		d.print("next")
	case ppSubass:
		d.print("`" + name + "`(")
		d.args(c.Args)
		d.print(")")
	}
}

// rexp10 is R's Rexp10: the power of ten as a correctly rounded double.
func rexp10(k int) float64 {
	f, _ := strconv.ParseFloat("1e"+strconv.Itoa(k), 64)
	return f
}

// rPrintable: R writes the character as itself, not as an escape.
func rPrintable(r rune) bool {
	lo, hi := 0, len(nonPrintable)
	for lo < hi {
		m := (lo + hi) / 2
		switch {
		case r < nonPrintable[m][0]:
			hi = m
		case r > nonPrintable[m][1]:
			lo = m + 1
		default:
			return false
		}
	}
	return true
}

// PrimitiveNames returns the names bound to primitives in R's base
// environment (a copy; see prims_gen.go).
func PrimitiveNames() map[string]bool {
	m := make(map[string]bool, len(primitiveNames))
	for k, v := range primitiveNames {
		m[k] = v
	}
	return m
}

// FormatNumber writes x as R does with 15 significant digits (the text of
// as.character(x) and of a deparsed double constant).
func FormatNumber(x float64) string { return formatReal15(x) }
