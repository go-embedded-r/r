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
				out[i] = complex(naReal, 0) // as R coerces: only the real part is NA
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
	if t == "i" {
		return complex(0, naReal), false // R reads the real part 0 and finds no number before the i
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
	// The C99 algorithm of Annex G (the divisor is scaled by a power of two
	// first), with the first product of each sum fused as the oracle's build
	// does; math.FMA makes that the same on every platform. Found by comparing
	// candidate algorithms with the oracle on 600 random quotients.
	ar, ai, c, d := real(a), imag(a), real(b), imag(b)
	logbw := math.Logb(math.Max(math.Abs(c), math.Abs(d)))
	ilogbw := 0
	if !math.IsInf(logbw, 0) && !math.IsNaN(logbw) {
		ilogbw = int(logbw)
		c = math.Ldexp(c, -ilogbw)
		d = math.Ldexp(d, -ilogbw)
	}
	denom := math.FMA(c, c, d*d)
	nx := math.FMA(ar, c, ai*d)
	ny := math.FMA(ai, c, -(ar * d))
	x, y := math.Ldexp(nx/denom, -ilogbw), math.Ldexp(ny/denom, -ilogbw)
	if math.IsNaN(x) && math.IsNaN(y) {
		a0, b0, c0, d0 := real(a), imag(a), real(b), imag(b)
		fin := func(v float64) bool { return !math.IsInf(v, 0) && !math.IsNaN(v) }
		ind := func(v float64) float64 {
			if math.IsInf(v, 0) {
				return 1
			}
			return 0
		}
		switch {
		case denom == 0 && (!math.IsNaN(a0) || !math.IsNaN(b0)):
			x = math.Copysign(math.Inf(1), c0) * a0
			y = math.Copysign(math.Inf(1), c0) * b0
		case (math.IsInf(a0, 0) || math.IsInf(b0, 0)) && fin(c0) && fin(d0):
			a1, b1 := math.Copysign(ind(a0), a0), math.Copysign(ind(b0), b0)
			x = math.Inf(1) * (a1*c0 + b1*d0)
			y = math.Inf(1) * (b1*c0 - a1*d0)
		case math.IsInf(logbw, 1) && fin(a0) && fin(b0):
			c1, d1 := math.Copysign(ind(c0), c0), math.Copysign(ind(d0), d0)
			x = 0 * (a0*c1 + b0*d1)
			y = 0 * (b0*c1 - a0*d1)
		}
	}
	return complex(x, y)
}

func cmul(a, b complex128) complex128 {
	ar, ai, br, bi := real(a), imag(a), real(b), imag(b)
	// float64(...) forbids the fused multiply-add a compiler may choose
	return complex(float64(ar*br)-float64(ai*bi), float64(ar*bi)+float64(ai*br))
}

// cpow is x^y: a small integer power by repeated multiplication, else through
// the exponential; 0 to a real power and the NA cases have their own answers.
func cpow(x, y complex128) complex128 {
	if y == 0 {
		return 1
	}
	if isNAComplex(x) || isNAComplex(y) {
		return naComplex
	}
	if imag(y) == 0 {
		k := real(y)
		if x == 0 {
			switch {
			case k > 0:
				return 0
			case k < 0:
				return complex(math.Inf(1), 0)
			}
		}
		if k > 0 && k == math.Trunc(k) && imag(x) == 0 && !math.IsNaN(real(x)) && !math.IsInf(real(x), 0) {
			// a real base stays real
			r := rpow(real(x), k)
			return complex(r, 0)
		}
		if k == math.Trunc(k) && math.Abs(k) <= 65536 {
			// square and multiply from the low bit, starting from 1 (found by
			// comparison with the oracle, which also squares a number directly)
			n := int(math.Abs(k))
			var res complex128 = 1
			base := x
			switch n {
			case 1:
				res = x
			case 2:
				res = cmul(x, x)
			default:
				for n > 0 {
					if n&1 == 1 {
						res = cmul(res, base)
					}
					base = cmul(base, base)
					n >>= 1
				}
			}
			if k < 0 {
				return cdiv(1, res)
			}
			return res
		}
	} else if x == 0 {
		return complex(rNaN, rNaN)
	}
	return cmplx.Exp(cmul(y, clog(x)))
}

func (s *Session) arithComplex(op string, c *callInfo, x, y Value, n int, attrs []Attr) Value {
	nx, ny := Length(x), Length(y)
	switch op {
	case "+", "-", "*", "/", "^":
	default:
		s.errorf(nil, "unimplemented complex operation")
	}
	a, b := asComplex(x), asComplex(y)
	out := make([]complex128, n)
	for i := range out {
		p, q := a[i%nx], b[i%ny]
		switch op {
		case "+":
			out[i] = cadd(p, q)
		case "-":
			out[i] = csub(p, q)
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

// zmul is a product in which an exact zero factor gives zero, as the C99
// complex functions give it (sin 0 * cosh 800 is 0, not NaN).
func zmul(a, b float64) float64 {
	if a == 0 || b == 0 {
		return 0
	}
	return a * b
}

func csin(z complex128) complex128 {
	x, y := real(z), imag(z)
	return complex(zmul(math.Sin(x), math.Cosh(y)), zmul(math.Cos(x), math.Sinh(y)))
}

func ccos(z complex128) complex128 {
	x, y := real(z), imag(z)
	return complex(zmul(math.Cos(x), math.Cosh(y)), -zmul(math.Sin(x), math.Sinh(y)))
}

func csinh(z complex128) complex128 {
	x, y := real(z), imag(z)
	return complex(zmul(math.Sinh(x), math.Cos(y)), zmul(math.Cosh(x), math.Sin(y)))
}

func ccosh(z complex128) complex128 {
	x, y := real(z), imag(z)
	return complex(zmul(math.Cosh(x), math.Cos(y)), zmul(math.Sinh(x), math.Sin(y)))
}

func cexp(z complex128) complex128 {
	if imag(z) == 0 {
		return complex(math.Exp(real(z)), 0)
	}
	return cmplx.Exp(z)
}

// ctan is tan(z); where the library's answer overflows the limit is returned.
func ctan(z complex128) complex128 {
	x, y := real(z), imag(z)
	if !math.IsNaN(x) && !math.IsNaN(y) && !math.IsInf(x, 0) && math.Abs(y) >= 20 {
		// the oracle answers +-i, with a zero real part, from |Im z| = 20 on
		return complex(math.Copysign(0, math.Sin(2*x)), math.Copysign(1, y))
	}
	r := cmplx.Tan(z)
	if !math.IsNaN(real(r)) && !math.IsInf(real(r), 0) && !math.IsNaN(imag(r)) && !math.IsInf(imag(r), 0) {
		return r
	}
	if math.IsNaN(x) || math.IsNaN(y) || math.IsInf(x, 0) {
		return r
	}
	d := math.Cos(2*x) + math.Cosh(2*y)
	return complex(math.Sin(2*x)/d, math.Sinh(2*y)/d)
}

// ctanh is -i tan(iz).
func ctanh(z complex128) complex128 {
	w := ctan(complex(-imag(z), real(z)))
	return complex(imag(w), -real(w))
}

// catan is the arctangent as the oracle computes it: half the arctangent of
// (2x, 1-x^2-y^2), whose sign follows the sign of x (a signed zero too), and a
// quarter of the logarithm of a ratio of distances (the oracle's own last
// digits are less accurate than these: see the cplx-inv exclusion of the
// rawcplx fuzz mode). On +-i the real part is +-pi/4.
func catan(z complex128) complex128 {
	return catanXY(real(z), imag(z))
}

func catanXY(x, y float64) complex128 {
	var re float64
	if x == 0 && math.Abs(y) > 1 {
		re = math.Copysign(math.Pi/2, y) // beyond +-i the sign follows y
	} else if x == 0 && math.Abs(y) == 1 {
		re = math.Copysign(math.Pi/4, x) // on +-i it follows the sign of the zero x
	} else {
		if x == 0 {
			x = 0 // inside the unit disk -0 acts as +0 (atanh(0.5+0i) has a -0 imaginary part)
		}
		re = 0.5 * math.Atan2(2*x, 1-x*x-y*y)
	}
	den := x*x + (y-1)*(y-1)
	var im float64
	switch t := 4 * y / den; {
	case den == 0:
		im = math.Copysign(math.Inf(1), y)
	case t > -0.5:
		im = 0.25 * math.Log1p(t)
	default:
		im = 0.25 * math.Log((x*x+(y+1)*(y+1))/den)
	}
	return complex(re, im)
}

// hullAsin is the real and imaginary parts of asin and acos for x, y >= 0,
// by the algorithm of Hull, Fairgrieve and Tang ("Implementing the complex
// arcsine and arccosine functions using exception handling", 1997), which
// avoids the cancellation of the textbook formulas: it returns asin's real
// part, acos's real part and the common imaginary part (before signs).
func hullAsin(x, y float64) (asinRe, acosRe, im float64) {
	const aCross, bCross = 1.5, 0.6417
	r := chypot(x+1, y)
	s := chypot(x-1, y)
	a := 0.5 * (r + s)
	b := x / a
	// am1 is A-1 and amx is A-x, each computed without cancellation
	var am1, amx float64
	if x <= 1 {
		t := 0.0
		if y != 0 {
			t = y * y / (s + (1 - x))
		}
		am1 = 0.5 * (y*y/(r+x+1) + t)
		amx = 0.5 * (y*y/(r+x+1) + (s + (1 - x)))
	} else {
		am1 = 0.5 * (y*y/(r+x+1) + (s + (x - 1)))
		amx = 0.5 * (y*y/(r+x+1) + y*y/(s+(x-1)))
	}
	if b <= bCross {
		asinRe = math.Asin(b)
		acosRe = math.Acos(b)
	} else {
		t := math.Sqrt(amx * (a + x))
		asinRe = math.Atan(x / t)
		acosRe = math.Atan(t / x)
	}
	if a <= aCross {
		im = math.Log1p(am1 + math.Sqrt(am1*(a+1)))
	} else {
		im = math.Log(a + math.Sqrt(a*a-1))
	}
	return
}

// casin: on the real axis beyond +-1 the imaginary part is -sign(x) acosh|x|.
func casin(z complex128) complex128 {
	x, y := real(z), imag(z)
	if y == 0 && math.Abs(x) > 1 {
		return complex(math.Copysign(math.Pi/2, x), -math.Copysign(math.Acosh(math.Abs(x)), x))
	}
	if math.IsNaN(x) || math.IsNaN(y) || math.IsInf(x, 0) || math.IsInf(y, 0) {
		return cmplx.Asin(z)
	}
	re, _, im := hullAsin(math.Abs(x), math.Abs(y))
	return complex(math.Copysign(re, x), math.Copysign(im, y))
}

func cabs2(z complex128) float64 { return real(z)*real(z) + imag(z)*imag(z) }

func cacos(z complex128) complex128 {
	x, y := real(z), imag(z)
	if y == 0 && math.Abs(x) > 1 {
		a := casin(z)
		return complex(math.Pi/2-real(a), -imag(a))
	}
	if math.IsNaN(x) || math.IsNaN(y) || math.IsInf(x, 0) || math.IsInf(y, 0) {
		return cmplx.Acos(z)
	}
	_, re, im := hullAsin(math.Abs(x), math.Abs(y))
	if x < 0 {
		re = math.Pi - re
	}
	return complex(re, -math.Copysign(im, y))
}

// cacosh is i acos(z), the sign convention the oracle shows.
func cacosh(z complex128) complex128 {
	a := cacos(z)
	return complex(-imag(a), real(a))
}

// catanh is -i atan(iz); at +-1 it is +-Inf with a NaN imaginary part (nan
// tells the caller to warn).
func catanh(z complex128) (r complex128, nan bool) {
	if imag(z) == 0 && math.Abs(real(z)) == 1 {
		return complex(math.Copysign(math.Inf(1), real(z)), rNaN), true
	}
	w := catanXY(-imag(z), real(z))
	return complex(imag(w), -real(w)), false
}

var complexMath = map[string]func(complex128) complex128{
	"sqrt": cmplx.Sqrt, "exp": cexp,
	"sin": csin, "cos": ccos, "tan": ctan,
	"asin": casin, "acos": cacos, "atan": catan,
	"sinh": csinh, "cosh": ccosh, "tanh": ctanh,
	"asinh": casinh, "acosh": cacosh,
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
	if !ok && name != "atanh" {
		s.errorf(c.call, "unimplemented complex function")
	}
	out := make([]complex128, len(x.V))
	warn := false
	for i, z := range x.V {
		if isNAComplex(z) {
			out[i] = naComplex
			continue
		}
		if name == "atanh" {
			var w bool
			out[i], w = catanh(z)
			warn = warn || w
			continue
		}
		out[i] = f(z)
	}
	if warn {
		s.warningf(c.call, "NaNs produced in function \"atanh\"")
	}
	return &Complex{V: out, Attr: x.Attr}
}

func cabs(z complex128) float64 {
	if isNAComplex(z) {
		return naReal
	}
	return chypot(real(z), imag(z))
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
	l := clog(z)
	if base != nil {
		return cdiv(l, clog(*base))
	}
	return l
}

// roundComplex is round(z, digits) (each part rounded alike) or signif(z,
// digits) (both parts rounded to the decimal place the larger one needs).
func roundComplex(z complex128, digits int, signif bool, round func(float64, int) float64) complex128 {
	if isNAReal(real(z)) && isNAReal(imag(z)) {
		return naComplex
	}
	if !signif {
		return complex(round(real(z), digits), round(imag(z), digits))
	}
	m := 0.0
	if !math.IsNaN(real(z)) && !math.IsInf(real(z), 0) {
		m = math.Abs(real(z))
	}
	if !math.IsNaN(imag(z)) && !math.IsInf(imag(z), 0) && math.Abs(imag(z)) > m {
		m = math.Abs(imag(z))
	}
	if m == 0 {
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
				acc = cmulFused(acc, z)
			} else {
				acc = cumAdd(acc, z)
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
					out[i] = f(z)
				}
				return &Double{V: out, Attr: v.Attr}
			case *Logical, *Integer, *Double:
				z := asComplex(v)
				out := make([]float64, len(z))
				for i, e := range z {
					out[i] = f(e)
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
			switch t := v.(type) {
			case *Logical, *Integer, *Double:
				return asDoubles(v)
			case *String:
				out := make([]float64, len(t.V))
				warn := false
				for i, e := range t.V {
					f, ok := parseNumber(e)
					if !ok {
						f, warn = naReal, true
					}
					out[i] = f
				}
				if warn {
					s.warningf(c.call, "NAs introduced by coercion")
				}
				return out
			}
			s.errorf(c.call, "invalid argument type")
			return nil
		}
		n := 0
		if a[0] != nil {
			d := num(a[0], "length.out")
			if len(d) < 1 || math.IsNaN(d[0]) || d[0] < 0 {
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
	case *List:
		for _, e := range v.V {
			if Length(e) != 1 || !isVector(e) || kindOf(e) == kList {
				s.errorf(nil, "'list' object cannot be coerced to type 'complex'")
			}
		}
		comb := s.combine(&callInfo{call: c.call, env: c.env, args: v.V, names: make([]string, len(v.V))})
		return s.asComplexFn(&callInfo{call: c.call, env: c.env, args: []Value{withAttr(comb, nil)}, names: []string{""}})
	case *String:
		out := make([]complex128, len(v.V))
		warn := false
		for i, e := range v.V {
			z, ok := parseComplex(e)
			if !ok {
				warn = true
				if z == 0 {
					z = naComplex
				}
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

// formatComplexItems is format() of a complex vector: the real parts and the
// absolute values of the imaginary parts are formatted as two double vectors
// (nsmall, scientific and big.mark apply to each; width pads each), and
// trim drops the padding between the parts.
func (s *Session) formatComplexItems(x []complex128, digits, scipen, nsmall, width int, trim bool, bigMark Value) []string {
	var rs, is []float64
	for _, z := range x {
		if isNAComplex(z) {
			continue
		}
		rs = append(rs, real(z))
		is = append(is, math.Abs(imag(z)))
	}
	fr, fi := formatRealVecSci(rs, digits, scipen), formatRealVecSci(is, digits, scipen)
	for _, f := range []*realFormat{&fr, &fi} {
		if nsmall > 0 && !f.e && f.d < nsmall {
			f.w += nsmall - f.d
			if f.d == 0 {
				f.w++
			}
			f.d = nsmall
		}
		if width > f.w {
			f.w = width
		}
		if trim {
			f.w = 0
		}
	}
	mark := ""
	if st, ok := bigMark.(*String); ok && len(st.V) == 1 {
		mark = st.V[0].S
	}
	part := func(f realFormat, v float64) string {
		t := strings.TrimLeft(f.encode(v), " ")
		if mark != "" {
			t = insertBigMark(t, mark)
		}
		return padLeft(t, f.w)
	}
	items := make([]string, len(x))
	total := fr.w + fi.w + 2
	for i, z := range x {
		if isNAComplex(z) {
			items[i] = padLeft("NA", total)
			if trim {
				items[i] = "NA"
			}
			continue
		}
		sign := "+"
		if imag(z) < 0 {
			sign = "-"
		}
		items[i] = part(fr, real(z)) + sign + part(fi, math.Abs(imag(z))) + "i"
	}
	return items
}

// chypot is sqrt(x*x + y*y) rounded as the C library rounds it (correctly in
// all but the rarest cases): the sum of squares is carried as a double-double
// and the square root is corrected once. Infinities win over NaN.
func chypot(x, y float64) float64 {
	x, y = math.Abs(x), math.Abs(y)
	switch {
	case math.IsInf(x, 0) || math.IsInf(y, 0):
		return math.Inf(1)
	case math.IsNaN(x) || math.IsNaN(y):
		return x + y
	}
	if x < y {
		x, y = y, x
	}
	if x == 0 {
		return 0
	}
	// scale so the squares neither overflow nor underflow
	_, e := math.Frexp(x)
	sx, sy := math.Ldexp(x, -e), math.Ldexp(y, -e)
	p1 := sx * sx
	e1 := math.FMA(sx, sx, -p1)
	p2 := sy * sy
	e2 := math.FMA(sy, sy, -p2)
	hi := p1 + p2
	bb := hi - p1
	lo := (p1 - (hi - bb)) + (p2 - bb) + e1 + e2
	h := math.Sqrt(hi)
	r := math.FMA(-h, h, hi) + lo
	h += r / (2 * h)
	return math.Ldexp(h, e)
}

// cmulFused is the product prod() and cumprod() compute: the oracle's build
// fuses the first product of each part (the operator * does not), which
// math.FMA reproduces on every platform.
func cmulFused(a, b complex128) complex128 {
	ar, ai, br, bi := real(a), imag(a), real(b), imag(b)
	return complex(math.FMA(ar, br, -(ai*bi)), math.FMA(ar, bi, ai*br))
}

// cmulCum is the product cumprod() computes: fused like prod's, but with the
// two terms of the imaginary part the other way round (found by comparison).
func cmulCum(a, b complex128) complex128 {
	ar, ai, br, bi := real(a), imag(a), real(b), imag(b)
	return complex(math.FMA(ar, br, -(ai*bi)), math.FMA(ai, br, ar*bi))
}

// casinh is -i asin(iz).
func casinh(z complex128) complex128 {
	a := casin(complex(-imag(z), real(z)))
	return complex(imag(a), -real(a))
}

// clog is log z with the modulus as the C library rounds it.
func clog(z complex128) complex128 {
	x, y := real(z), imag(z)
	a, b := math.Abs(x), math.Abs(y)
	if a < b {
		a, b = b, a
	}
	var re float64
	switch {
	case math.IsInf(a, 0):
		re = math.Inf(1)
	case math.IsNaN(a) || math.IsNaN(b):
		re = math.NaN()
	case a == 0:
		re = math.Inf(-1)
	default:
		// log of the larger part plus half of log1p of the squared ratio: the
		// form that matches the oracle digit for digit
		t := b / a
		re = math.Log(a) + 0.5*math.Log1p(t*t)
	}
	return complex(re, math.Atan2(y, x))
}

// fadd and fsub are + and - with the NaN of the first operand winning, which
// keeps NA (a NaN with a payload) as R does whatever order the hardware takes.
func fadd(a, b float64) float64 {
	if a != a {
		return a
	}
	if b != b {
		return b
	}
	return a + b
}

func fsub(a, b float64) float64 {
	if a != a {
		return a
	}
	if b != b {
		return b
	}
	return a - b
}

func cadd(a, b complex128) complex128 { return complex(fadd(real(a), real(b)), fadd(imag(a), imag(b))) }
func csub(a, b complex128) complex128 { return complex(fsub(real(a), real(b)), fsub(imag(a), imag(b))) }

// cumAdd is the step of cumsum: the parts are added each on its own, and an NA
// in either operand of a part gives NA, a NaN gives NaN.
func cumAdd(acc, z complex128) complex128 {
	part := func(a, b float64) float64 {
		switch {
		case isNAReal(a) || isNAReal(b):
			return naReal
		case math.IsNaN(a) || math.IsNaN(b):
			return rNaN
		}
		return a + b
	}
	return complex(part(real(acc), real(z)), part(imag(acc), imag(z)))
}
