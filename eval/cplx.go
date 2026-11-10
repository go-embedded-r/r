package eval

import (
	"math"
	"math/cmplx"
	"strconv"
	"strings"

	"github.com/go-embedded-r/r/syntax"
)

// The complex type (CPLXSXP), written from ?complex, ?Arithmetic, ?Math and
// the oracle's behaviour; nothing is translated from R's sources.

var naComplex = complex(naReal, naReal)

// isNAComplex is R's test of a complex NA: either part is NA_real_.
func isNAComplex(z complex128) bool { return isNAReal(real(z)) || isNAReal(imag(z)) }

// asComplex converts a logical, integer, double, complex or raw vector.
func asComplex(v Value) []complex128 {
	switch x := v.(type) {
	case *Complex:
		return x.V
	case *Double:
		out := make([]complex128, len(x.V))
		for i, e := range x.V {
			out[i] = complex(e, 0)
		}
		return out
	case *Logical, *Integer:
		iv := asInts(v)
		out := make([]complex128, len(iv))
		for i, e := range iv {
			if e == naInt {
				out[i] = naComplex
			} else {
				out[i] = complex(float64(e), 0)
			}
		}
		return out
	case *Raw:
		out := make([]complex128, len(x.V))
		for i, e := range x.V {
			out[i] = complex(float64(e), 0)
		}
		return out
	}
	return nil
}

// complexStr is as.character of one complex number: each part with up to 15
// significant digits.
func complexStr(z complex128) Str {
	if isNAComplex(z) {
		return Str{NA: true}
	}
	re, im := real(z), imag(z)
	sign := "+"
	if im < 0 {
		sign = "-"
		im = -im
	}
	return Str{S: syntax.FormatNumber(re) + sign + syntax.FormatNumber(im) + "i"}
}

// parseComplex is String2Complex for the spellings R accepts: a real number,
// an imaginary one (2i), or both (1+2i).
func parseComplex(s Str) (complex128, bool) {
	if s.NA {
		return naComplex, true
	}
	t := strings.TrimSpace(s.S)
	if t == "" {
		return 0, false
	}
	if f, ok := parseNumber(Str{S: t}); ok {
		return complex(f, 0), true
	}
	if !strings.HasSuffix(t, "i") {
		return 0, false
	}
	body := t[:len(t)-1]
	// the split is the last sign that does not belong to an exponent
	split := -1
	for i := len(body) - 1; i > 0; i-- {
		if (body[i] == '+' || body[i] == '-') && body[i-1] != 'e' && body[i-1] != 'E' {
			split = i
			break
		}
	}
	if split < 0 {
		im, ok := parseNumber(Str{S: body})
		if !ok {
			return 0, false
		}
		return complex(0, im), true
	}
	re, ok1 := parseNumber(Str{S: body[:split]})
	im, ok2 := parseNumber(Str{S: body[split:]})
	if !ok1 || !ok2 {
		return 0, false
	}
	return complex(re, im), true
}

// ---- formatting ---------------------------------------------------------------

// complexFormat is the format shared by the elements of a complex vector: the
// real parts are formatted as one double vector and the absolute values of
// the imaginary parts as another.
type complexFormat struct{ re, im realFormat }

func formatComplexVec(x []complex128, digits int) complexFormat {
	var rs, is []float64
	for _, z := range x {
		if isNAComplex(z) {
			continue
		}
		rs = append(rs, real(z))
		is = append(is, math.Abs(imag(z)))
	}
	return complexFormat{formatRealVec(rs, digits), formatRealVec(is, digits)}
}

func (f complexFormat) width() int {
	w := f.re.w + f.im.w + 2
	if w < 2 {
		w = 2
	}
	return w
}

func (f complexFormat) encode(z complex128) string {
	if isNAComplex(z) {
		return padLeft("NA", f.width())
	}
	sign := "+"
	if imag(z) < 0 {
		sign = "-"
	}
	return f.re.encode(real(z)) + sign + f.im.encode(math.Abs(imag(z))) + "i"
}

// ---- arithmetic -----------------------------------------------------------------

func cdiv(a, b complex128) complex128 {
	ar, ai, br, bi := real(a), imag(a), real(b), imag(b)
	if math.IsNaN(ar) || math.IsNaN(ai) || math.IsNaN(br) || math.IsNaN(bi) || math.IsInf(br, 0) || math.IsInf(bi, 0) || math.IsInf(ar, 0) || math.IsInf(ai, 0) || (br == 0 && bi == 0) {
		return a / b
	}
	// Smith's algorithm
	if math.Abs(br) >= math.Abs(bi) {
		ratio := bi / br
		den := br + bi*ratio
		return complex((ar+ai*ratio)/den, (ai-ar*ratio)/den)
	}
	ratio := br / bi
	den := br*ratio + bi
	return complex((ar*ratio+ai)/den, (ai*ratio-ar)/den)
}

func cmul(a, b complex128) complex128 {
	ar, ai, br, bi := real(a), imag(a), real(b), imag(b)
	return complex(ar*br-ai*bi, ar*bi+ai*br)
}

// cpow is x^y: a small integer power by repeated multiplication, else through
// the exponential.
func cpow(x, y complex128) complex128 {
	if imag(y) == 0 {
		k := real(y)
		if k == 0 {
			return 1
		}
		if k == math.Trunc(k) && math.Abs(k) <= 65536 {
			n := int(math.Abs(k))
			var res complex128 = 1
			base := x
			for n > 0 {
				if n&1 == 1 {
					res = cmul(res, base)
				}
				base = cmul(base, base)
				n >>= 1
			}
			if k < 0 {
				return cdiv(1, res)
			}
			return res
		}
	}
	if x == 0 {
		if real(y) > 0 {
			return 0
		}
	}
	return cmplx.Exp(cmul(y, cmplx.Log(x)))
}

func (s *Session) arithComplex(op string, c *callInfo, x, y Value, n int, attrs []Attr) Value {
	nx, ny := Length(x), Length(y)
	switch op {
	case "+", "-", "*", "/", "^":
	default:
		s.errorf(c.call, "invalid operation on complex numbers")
	}
	a, b := asComplex(x), asComplex(y)
	out := make([]complex128, n)
	for i := range out {
		p, q := a[i%nx], b[i%ny]
		switch op {
		case "+":
			out[i] = p + q
		case "-":
			out[i] = p - q
		case "*":
			out[i] = cmul(p, q)
		case "/":
			out[i] = cdiv(p, q)
		case "^":
			out[i] = cpow(p, q)
		}
	}
	return &Complex{V: out, Attr: attrs}
}

func (s *Session) compareComplex(op string, c *callInfo, x, y Value, n int) Value {
	if op != "==" && op != "!=" {
		s.errorf(c.call, "invalid comparison with complex values")
	}
	nx, ny := Length(x), Length(y)
	a, b := asComplex(x), asComplex(y)
	out := make([]int32, n)
	for i := range out {
		p, q := a[i%nx], b[i%ny]
		if math.IsNaN(real(p)) || math.IsNaN(imag(p)) || math.IsNaN(real(q)) || math.IsNaN(imag(q)) {
			out[i] = naInt
			continue
		}
		eq := p == q
		if (op == "==") == eq {
			out[i] = 1
		}
	}
	var keep []Attr
	for _, at := range copyAttrs(x, y, n) {
		if at.Name == "names" || at.Name == "dim" || at.Name == "dimnames" {
			keep = append(keep, at)
		}
	}
	return &Logical{V: out, Attr: keep}
}

// ---- the math group ---------------------------------------------------------------

var complexMath = map[string]func(complex128) complex128{
	"sqrt": cmplx.Sqrt, "exp": cmplx.Exp,
	"sin": cmplx.Sin, "cos": cmplx.Cos, "tan": cmplx.Tan,
	"asin": cmplx.Asin, "acos": cmplx.Acos, "atan": cmplx.Atan,
	"sinh": cmplx.Sinh, "cosh": cmplx.Cosh, "tanh": cmplx.Tanh,
	"asinh": cmplx.Asinh, "acosh": cmplx.Acosh, "atanh": cmplx.Atanh,
}

// complexUnary is the math group on a complex vector.
func (s *Session) complexUnary(c *callInfo, name string, x *Complex) Value {
	switch name {
	case "abs":
		out := make([]float64, len(x.V))
		for i, z := range x.V {
			out[i] = cabs(z)
		}
		return &Double{V: out, Attr: x.Attr}
	}
	f, ok := complexMath[name]
	if !ok {
		s.errorf(c.call, "unimplemented complex function")
	}
	out := make([]complex128, len(x.V))
	for i, z := range x.V {
		if isNAComplex(z) {
			out[i] = naComplex
			continue
		}
		out[i] = f(z)
	}
	return &Complex{V: out, Attr: x.Attr}
}

func cabs(z complex128) float64 {
	if isNAComplex(z) {
		return naReal
	}
	return math.Hypot(real(z), imag(z))
}

func carg(z complex128) float64 {
	if isNAComplex(z) {
		return naReal
	}
	return math.Atan2(imag(z), real(z))
}

// complexLog is log(z) and log(z, base).
func complexLog(z complex128, base *complex128) complex128 {
	if isNAComplex(z) {
		return naComplex
	}
	l := cmplx.Log(z)
	if base != nil {
		return cdiv(l, cmplx.Log(*base))
	}
	return l
}

// roundComplex is round(z, digits) (each part rounded alike) or signif(z,
// digits) (both parts rounded to the decimal place the larger one needs).
func roundComplex(z complex128, digits int, signif bool, round func(float64, int) float64) complex128 {
	if isNAComplex(z) {
		return naComplex
	}
	if !signif {
		return complex(round(real(z), digits), round(imag(z), digits))
	}
	m := math.Max(math.Abs(real(z)), math.Abs(imag(z)))
	if m == 0 || math.IsInf(m, 0) || math.IsNaN(m) {
		return z
	}
	if digits < 1 {
		digits = 1
	}
	dig := digits - int(math.Floor(math.Log10(m))) - 1
	return complex(round(real(z), dig), round(imag(z), dig))
}

// complexSummary is sum, prod and mean of complex data.
func complexSum(data []Value, prod bool, rm bool) complex128 {
	var acc complex128
	if prod {
		acc = 1
	}
	for _, a := range data {
		for _, z := range asComplex(a) {
			if rm && (math.IsNaN(real(z)) || math.IsNaN(imag(z))) {
				continue
			}
			if prod {
				acc = cmul(acc, z)
			} else {
				acc += z
			}
		}
	}
	return acc
}

func rawHex(b byte) string {
	const hexd = "0123456789abcdef"
	return string([]byte{hexd[b>>4], hexd[b&15]})
}

func itoaInt(i int) string { return strconv.Itoa(i) }

// ---- the builtins -------------------------------------------------------------------

func init() {
	regBuiltin("is.complex", func(s *Session, c *callInfo) Value {
		_, ok := s.nargs1(c, "is.complex").(*Complex)
		return mkLogical(ok)
	})
	part := func(name string, f func(complex128) float64) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			x := s.nargs1(c, name)
			switch v := x.(type) {
			case *Complex:
				out := make([]float64, len(v.V))
				for i, z := range v.V {
					if isNAComplex(z) {
						out[i] = naReal
					} else {
						out[i] = f(z)
					}
				}
				return &Double{V: out, Attr: v.Attr}
			case *Logical, *Integer, *Double:
				z := asComplex(v)
				out := make([]float64, len(z))
				for i, e := range z {
					if isNAComplex(e) {
						out[i] = naReal
					} else {
						out[i] = f(e)
					}
				}
				return &Double{V: out, Attr: attrsOf(x)}
			}
			s.errorf(c.call, "non-numeric argument to function")
			return nil
		})
	}
	part("Re", func(z complex128) float64 { return real(z) })
	part("Im", func(z complex128) float64 { return imag(z) })
	part("Mod", cabs)
	part("Arg", carg)
	regBuiltin("Conj", func(s *Session, c *callInfo) Value {
		x := s.nargs1(c, "Conj")
		switch v := x.(type) {
		case *Complex:
			out := make([]complex128, len(v.V))
			for i, z := range v.V {
				out[i] = complex(real(z), -imag(z))
			}
			return &Complex{V: out, Attr: v.Attr}
		case *Logical, *Integer:
			return x
		case *Double:
			return x
		}
		s.errorf(c.call, "non-numeric argument to function")
		return nil
	})
	regBuiltin("as.complex", func(s *Session, c *callInfo) Value { return s.asComplexFn(c) })
	regBuiltin("complex", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "length.out", "real", "imaginary", "modulus", "argument")
		num := func(v Value, what string) []float64 {
			if v == nil {
				return nil
			}
			switch v.(type) {
			case *Logical, *Integer, *Double:
				return asDoubles(v)
			}
			s.errorf(c.call, "invalid argument type")
			return nil
		}
		n := 0
		if a[0] != nil {
			d := num(a[0], "length.out")
			if len(d) != 1 || math.IsNaN(d[0]) || d[0] < 0 {
				s.errorf(c.call, "invalid length")
			}
			if d[0] > float64(s.opts.MaxVector) {
				s.checkLen(c.call, s.opts.MaxVector+1)
			}
			n = int(d[0])
		}
		re, im, mo, ar := num(a[1], "real"), num(a[2], "imaginary"), num(a[3], "modulus"), num(a[4], "argument")
		polar := a[3] != nil || a[4] != nil
		for _, l := range []int{len(re), len(im), len(mo), len(ar)} {
			if l > n {
				n = l
			}
		}
		s.checkLen(c.call, n)
		out := make([]complex128, n)
		pick := func(v []float64, i int, def float64) float64 {
			if len(v) == 0 {
				return def
			}
			return v[i%len(v)]
		}
		for i := range out {
			if polar {
				m, g := pick(mo, i, 1), pick(ar, i, 0)
				out[i] = complex(m*math.Cos(g), m*math.Sin(g))
			} else {
				out[i] = complex(pick(re, i, 0), pick(im, i, 0))
			}
		}
		return &Complex{V: out}
	})
}

func (s *Session) asComplexFn(c *callInfo) Value {
	x := firstArg(c)
	switch v := x.(type) {
	case nullT:
		return &Complex{V: []complex128{}}
	case *Complex:
		return &Complex{V: v.V}
	case *Logical, *Integer, *Double, *Raw:
		return &Complex{V: asComplex(x)}
	case *String:
		out := make([]complex128, len(v.V))
		warn := false
		for i, e := range v.V {
			z, ok := parseComplex(e)
			if !ok {
				z, warn = naComplex, true
			}
			out[i] = z
		}
		if warn {
			s.warningf(c.call, "NAs introduced by coercion")
		}
		return &Complex{V: out}
	}
	s.errorf(c.call, "cannot coerce type '%s' to vector of type 'complex'", typeName(x))
	return nil
}

// cless orders complex numbers by real part, then imaginary part.
func cless(a, b complex128) bool {
	if real(a) != real(b) {
		return real(a) < real(b)
	}
	return imag(a) < imag(b)
}
