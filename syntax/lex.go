package syntax

import (
	"math"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const eofByte = -1

type lexState struct{ off, line, col int }

// lexer reads R source the way R's does: byte by byte, counting columns in
// characters, with the context stack that decides which newlines end a
// statement.
type lexer struct {
	src  []byte
	off  int
	line int
	col  int
	hist [16]lexState
	hn   int

	file     string // name used in lexical errors
	ctx      []byte // ' ' (top), '(', '[', '{', 'i' (inside an if body)
	eatLines bool
	saved    *lexTok

	toks []recTok // every terminal, in source order (comments included)

	havePH bool // a pipe placeholder was seen
}

// recTok is a recorded terminal.
type recTok struct {
	kind  string
	text  string
	first Pos
	last  Pos
	fb    int // byte offsets
	lb    int
}

// lexTok is a token as the parser receives it.
type lexTok struct {
	kind  Kind
	text  string
	first Pos
	last  Pos
	fb    int // first byte
	lb    int // byte after the last
	cend  int // bytes consumed when it was returned (R's parse context)
	rec   int // index in lexer.toks, -1 if not recorded
	// values
	sym  string
	cst  *Const
	line int // line of the position after the token (for errors)
}

func newLexer(src []byte, file string) *lexer {
	return &lexer{src: src, line: 1, col: 0, file: file, ctx: []byte{' '}}
}

func (l *lexer) getc() int {
	l.hist[l.hn&15] = lexState{l.off, l.line, l.col}
	l.hn++
	if l.off >= len(l.src) {
		return eofByte
	}
	c := l.src[l.off]
	l.off++
	if c == '\n' {
		l.line++
		l.col = 0
	} else if c < 0x80 || c >= 0xC0 {
		l.col++
	}
	if c == '\t' {
		l.col = (l.col + 7) &^ 7
	}
	return int(c)
}

func (l *lexer) ungetc() {
	l.hn--
	s := l.hist[l.hn&15]
	l.off, l.line, l.col = s.off, s.line, s.col
}

func (l *lexer) nextchar(expect byte) bool {
	c := l.getc()
	if c == int(expect) {
		return true
	}
	l.ungetc()
	return false
}

func (l *lexer) pos() Pos { return Pos{l.line, l.col} }

// raise aborts the parse with a lexical error at the current position.
func (l *lexer) raise(format string, args ...any) {
	panic(lexicalError(l.file, format, l.pos(), args...))
}

func isDigit(c int) bool { return c >= '0' && c <= '9' }

// mbNext decodes the character whose first byte c was just read. The bytes
// that follow are not consumed. An invalid sequence is a lexical error, with
// the position R reports (after the bytes it tried to read).
func (l *lexer) mbNext(c int) (rune, int) {
	if c < 0x80 {
		return rune(c), 1
	}
	n := 1
	switch {
	case c >= 0xFC:
		n = 6
	case c >= 0xF8:
		n = 5
	case c >= 0xF0:
		n = 4
	case c >= 0xE0:
		n = 3
	case c >= 0xC0:
		n = 2
	}
	// l.off is after the lead byte
	start := l.off - 1
	if start+n > len(l.src) {
		// EOF while reading a multibyte character
		return 0, -1
	}
	r, size := utf8.DecodeRune(l.src[start : start+n])
	if r == utf8.RuneError || size != n {
		for i := 1; i < n; i++ {
			l.getc()
		}
		l.raise("invalid multibyte character in parser")
	}
	return r, n
}

func (l *lexer) skipSpace() int {
	for {
		c := l.getc()
		if c == ' ' || c == '\t' || c == '\f' {
			continue
		}
		if c >= 0x80 {
			r, n := l.mbNext(c)
			if n == -1 {
				l.ungetc()
				return eofByte
			}
			if isBlank(r) {
				for i := 1; i < n; i++ {
					l.getc()
				}
				continue
			}
		}
		return c
	}
}

// isBlank is iswblank for a UTF-8 locale: the space separators, except the
// no-break ones.
func isBlank(r rune) bool {
	if r == 0xA0 || r == 0x2007 || r == 0x202F {
		return false
	}
	return unicode.Is(unicode.Zs, r)
}

// skipComment reads a comment whose '#' has been read; it returns the
// character that ends it (a newline or EOF), and records the comment.
func (l *lexer) skipComment() int {
	startOff := l.off - 1
	firstLast := l.pos()
	first := firstLast
	var last Pos = l.pos()
	c := int('#')
	lastOff := l.off
	for c != '\n' && c != eofByte {
		last = l.pos()
		lastOff = l.off
		c = l.getc()
	}
	if c == '\n' {
		// the last character of the comment is the one before the newline
		lastOff = l.off - 1
	} else {
		lastOff = l.off
	}
	text := string(l.src[startOff:lastOff])
	kind := "COMMENT"
	if first.Col == 1 && isLineDirective(l.src[startOff:lastOff]) {
		kind = "LINE_DIRECTIVE"
	}
	l.toks = append(l.toks, recTok{kind: kind, text: text, first: first, last: last, fb: startOff, lb: lastOff})
	return c
}

func isLetter(r rune) bool { return unicode.IsLetter(r) }
func isAlnum(r rune) bool  { return unicode.IsLetter(r) || unicode.IsDigit(r) }

func (l *lexer) typeofnext() int {
	c := l.getc()
	k := 2
	if isDigit(c) {
		k = 1
	}
	l.ungetc()
	return k
}

// token returns the next token, ignoring the newline context rules (see
// yylex).
func (l *lexer) token() lexTok {
	if l.saved != nil {
		t := *l.saved
		l.saved = nil
		return t
	}
	c := l.skipSpace()
	if c == '#' {
		c = l.skipComment()
	}
	t := lexTok{rec: -1}
	t.first = l.pos()
	t.fb = l.off - 1
	if c == eofByte {
		t.kind = tEOF
		t.fb = len(l.src)
		t.lb = len(l.src)
		t.last = t.first
		t.cend = l.off
		return t
	}
	l.scan(c, &t)
	t.last = l.pos()
	t.lb = l.off
	t.cend = l.off
	if t.text == "" && t.kind != tERROR && t.kind != tINCOMPLETE && t.kind != '\n' {
		t.text = string(l.src[t.fb:t.lb])
	}
	return t
}

func (l *lexer) scan(c int, t *lexTok) {
	if c == '.' && l.typeofnext() >= 2 {
		l.symbol(c, t)
		return
	}
	if c == '.' || isDigit(c) {
		l.numeric(c, t)
		return
	}
	if c == 'r' || c == 'R' {
		if l.nextchar('"') {
			l.rawString('"', t)
			return
		} else if l.nextchar('\'') {
			l.rawString('\'', t)
			return
		}
	}
	if c == '"' || c == '\'' {
		l.str(c, false, t)
		return
	}
	if c == '%' {
		l.special(t)
		return
	}
	if c == '`' {
		l.str(c, true, t)
		return
	}
	if c == '_' {
		t.kind = tPLACEHOLDER
		t.text = "_"
		return
	}
	if c >= 0x80 {
		r, n := l.mbNext(c)
		if n == -1 {
			t.kind = tEOF
			return
		}
		if isLetter(r) {
			l.symbol(c, t)
			return
		}
		// not a letter: a multibyte character that is no token
		for i := 1; i < n; i++ {
			l.getc()
		}
		t.kind = tERROR
		return
	}
	if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
		l.symbol(c, t)
		return
	}
	switch c {
	case '<':
		switch {
		case l.nextchar('='):
			t.kind = tLE
		case l.nextchar('-'):
			t.kind = tLEFTASSIGN
		case l.nextchar('<'):
			if l.nextchar('-') {
				t.kind = tLEFTASSIGN
			} else {
				t.kind = tERROR
			}
		default:
			t.kind = tLT
		}
		if t.kind == tLEFTASSIGN {
			t.sym = "<-"
			if l.off-t.fb == 3 {
				t.sym = "<<-"
			}
		} else if t.kind == tLE {
			t.sym = "<="
		} else if t.kind == tLT {
			t.sym = "<"
		}
	case '-':
		if l.nextchar('>') {
			t.kind = tRIGHTASSIGN
			t.sym = "<-"
			if l.nextchar('>') {
				t.sym = "<<-"
			}
		} else {
			t.kind = '-'
			t.sym = "-"
		}
	case '>':
		if l.nextchar('=') {
			t.kind, t.sym = tGE, ">="
		} else {
			t.kind, t.sym = tGT, ">"
		}
	case '!':
		if l.nextchar('=') {
			t.kind, t.sym = tNE, "!="
		} else {
			t.kind, t.sym = '!', "!"
		}
	case '=':
		switch {
		case l.nextchar('='):
			t.kind, t.sym = tEQ, "=="
		case l.nextchar('>'):
			t.kind, t.sym = tPIPEBIND, "=>"
		default:
			t.kind, t.sym = tEQASSIGN, "="
		}
	case ':':
		switch {
		case l.nextchar(':'):
			if l.nextchar(':') {
				t.kind, t.sym = tNSGETINT, ":::"
			} else {
				t.kind, t.sym = tNSGET, "::"
			}
		case l.nextchar('='):
			t.kind, t.sym = tLEFTASSIGN, ":="
		default:
			t.kind, t.sym = ':', ":"
		}
	case '&':
		if l.nextchar('&') {
			t.kind, t.sym = tAND2, "&&"
		} else {
			t.kind, t.sym = tAND, "&"
		}
	case '|':
		switch {
		case l.nextchar('|'):
			t.kind, t.sym = tOR2, "||"
		case l.nextchar('>'):
			t.kind, t.sym = tPIPE, "|>"
		default:
			t.kind, t.sym = tOR, "|"
		}
	case '{', '}', '(', ')', ']', '?', '+', '/', '^', '~', '$', '@', '\\', '\n', ',', ';':
		t.kind = Kind(c)
		t.sym = string(rune(c))
	case '[':
		if l.nextchar('[') {
			t.kind, t.sym = tLBB, "[["
		} else {
			t.kind, t.sym = '[', "["
		}
	case '*':
		if l.nextchar('*') {
			t.kind, t.sym = '^', "^"
		} else {
			t.kind, t.sym = '*', "*"
		}
	default:
		t.kind = tUNDEF
	}
}

func (l *lexer) symbol(c int, t *lexTok) {
	var buf []byte
	for {
		// c is the first byte of a character already known to be valid
		buf = append(buf, byte(c))
		if c >= 0x80 {
			_, n := l.mbNext(c)
			for i := 1; i < n; i++ {
				buf = append(buf, byte(l.getc()))
			}
		}
		c = l.getc()
		if c == eofByte {
			break
		}
		if c == '.' || c == '_' {
			continue
		}
		if c < 0x80 {
			if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || isDigit(c)) {
				break
			}
			continue
		}
		r, n := l.mbNext(c)
		if n == -1 || !isAlnum(r) {
			break
		}
	}
	l.ungetc()
	name := string(buf)
	t.sym = name
	t.text = name
	switch name {
	case "NULL":
		t.kind = tNULL
		t.cst = &Const{Kind: Null}
	case "NA":
		t.kind = tNUM
		t.cst = &Const{Kind: Logical, NA: true}
	case "TRUE":
		t.kind = tNUM
		t.cst = &Const{Kind: Logical, LogicalV: true}
	case "FALSE":
		t.kind = tNUM
		t.cst = &Const{Kind: Logical}
	case "Inf":
		t.kind = tNUM
		t.cst = &Const{Kind: Double, Double: math.Inf(1)}
	case "NaN":
		t.kind = tNUM
		t.cst = &Const{Kind: Double, Double: math.Float64frombits(0x7FF8000000000000)}
	case "NA_integer_":
		t.kind = tNUM
		t.cst = &Const{Kind: Integer, NA: true}
	case "NA_real_":
		t.kind = tNUM
		t.cst = &Const{Kind: Double, NA: true, Double: math.Float64frombits(0x7FF00000000007A2)}
	case "NA_character_":
		t.kind = tNUM
		t.cst = &Const{Kind: String, NA: true}
	case "NA_complex_":
		t.kind = tNUM
		na := math.Float64frombits(0x7FF00000000007A2)
		t.cst = &Const{Kind: Complex, NA: true, Double: na, Im: na}
	case "function":
		t.kind = tFUNCTION
	case "while":
		t.kind = tWHILE
	case "repeat":
		t.kind = tREPEAT
	case "for":
		t.kind = tFOR
	case "if":
		t.kind = tIF
	case "in":
		t.kind = tIN
	case "else":
		t.kind = tELSE
	case "next":
		t.kind = tNEXT
	case "break":
		t.kind = tBREAK
	default:
		t.kind = tSYMBOL
	}
}

func (l *lexer) special(t *lexTok) {
	for {
		c := l.getc()
		if c == eofByte || c == '%' {
			break
		}
		if c == '\n' {
			l.ungetc()
			t.kind = tERROR
			return
		}
	}
	t.kind = tSPECIAL
	t.sym = string(l.src[t.fb:l.off])
}

// numeric lexes a number whose first character c has been read.
func (l *lexer) numeric(c int, t *lexTok) {
	seendot := 0
	if c == '.' {
		seendot = 1
	}
	seenexp := 0
	last := c
	nd := 0
	asNumeric := false
	count := 1
	t.kind = tNUM
	for {
		c = l.getc()
		if !(isDigit(c) || c == '.' || c == 'e' || c == 'E' || c == 'x' || c == 'X' || c == 'L') {
			break
		}
		count++
		if c == 'L' {
			break
		}
		if c == 'x' || c == 'X' {
			if count > 2 || last != '0' {
				break
			}
			for {
				c = l.getc()
				if !(isDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') || c == '.') {
					break
				}
				if c == '.' {
					if seendot == 1 {
						t.kind = tERROR
						return
					}
					seendot = 1
				}
				nd++
			}
			if nd == 0 {
				t.kind = tERROR
				return
			}
			if c == 'p' || c == 'P' {
				seenexp = 1
				c = l.getc()
				if !isDigit(c) && c != '+' && c != '-' {
					t.kind = tERROR
					return
				}
				if c == '+' || c == '-' {
					c = l.getc()
				}
				nd = 0
				for isDigit(c) {
					c = l.getc()
					nd++
				}
				if nd == 0 {
					t.kind = tERROR
					return
				}
			}
			break
		}
		if c == 'E' || c == 'e' {
			if seenexp == 1 {
				break
			}
			seenexp = 1
			if seendot != 1 {
				seendot = 2
			}
			c = l.getc()
			if !isDigit(c) && c != '+' && c != '-' {
				t.kind = tERROR
				return
			}
			if c == '+' || c == '-' {
				c = l.getc()
				if !isDigit(c) {
					t.kind = tERROR
					return
				}
			}
		}
		if c == '.' {
			if seendot == 1 {
				break
			}
			seendot = 1
		}
		last = c
	}
	// c is the first character after the number's body, already read
	body := string(l.src[t.fb : l.off-boolToInt(c != eofByte)])
	if c == 'L' {
		body = string(l.src[t.fb : l.off-1])
	}
	hexText := body
	_ = hexText
	if c == 'L' {
		a := rAtof(body)
		b := int32(a)
		if a != float64(b) || math.IsNaN(a) {
			asNumeric = true
			seenexp = 1
		}
		if asNumeric {
			t.cst = &Const{Kind: Double, Double: a}
		} else {
			t.cst = &Const{Kind: Integer, Int: b}
		}
		return
	}
	if c == 'i' {
		t.cst = &Const{Kind: Complex, Double: 0, Im: rAtof(body)}
		return
	}
	if c != eofByte {
		l.ungetc()
	}
	_ = seenexp
	t.cst = &Const{Kind: Double, Double: rAtof(body)}
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// rAtof converts the text of a numeric literal to a double. Decimal
// literals are correctly rounded (strconv), so the result does not depend
// on the platform's long double; see the package documentation. Hexadecimal
// literals accumulate digit by digit in double precision.
func rAtof(s string) float64 {
	if len(s) > 1 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
		return hexAtof(s[2:])
	}
	return decAtof(s)
}

// decAtof accumulates the digits and scales by a power of ten, all in double
// precision: the arithmetic of R on a platform whose long double is a double
// (macOS arm64), where the compiler fuses 10*ans+digit into one rounding. The
// fusion is written as math.FMA so that the result is the same on every
// platform, whatever the Go compiler does for its target.
func decAtof(s string) float64 {
	var ans float64
	expn := 0
	ndigits := 0
	i := 0
	for ; i < len(s) && isDigit(int(s[i])); i++ {
		ans = math.FMA(10, ans, float64(s[i]-'0'))
		ndigits++
	}
	if i < len(s) && s[i] == '.' {
		i++
		for ; i < len(s) && isDigit(int(s[i])); i++ {
			ans = math.FMA(10, ans, float64(s[i]-'0'))
			ndigits++
			expn--
		}
	}
	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		i++
		sign := 1
		if i < len(s) && (s[i] == '+' || s[i] == '-') {
			if s[i] == '-' {
				sign = -1
			}
			i++
		}
		n := 0
		for ; i < len(s) && isDigit(int(s[i])); i++ {
			if n < 9999 {
				n = n*10 + int(s[i]-'0')
			}
		}
		expn += sign * n
	}
	p10, fac := 10.0, 1.0
	if expn < 0 {
		if expn+ndigits < -300 {
			// scale the digits down first, so that the power of ten below
			// cannot overflow
			for n := 0; n < ndigits; n++ {
				ans /= 10
			}
			expn += ndigits
		}
		if expn < -307 {
			// use underflow, not overflow
			for n := -expn; n != 0; n >>= 1 {
				if n&1 != 0 {
					fac /= p10
				}
				p10 *= p10
			}
			ans *= fac
		} else {
			for n := -expn; n != 0; n >>= 1 {
				if n&1 != 0 {
					fac *= p10
				}
				p10 *= p10
			}
			ans /= fac
		}
	} else if ans != 0 {
		for n := expn; n != 0; n >>= 1 {
			if n&1 != 0 {
				fac *= p10
			}
			p10 *= p10
		}
		ans *= fac
	}
	return ans
}

func correctlyRounded(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func hexAtof(s string) float64 {
	var ans float64
	i := 0
	for ; i < len(s) && isHex(s[i]); i++ {
		ans = math.FMA(16, ans, float64(hexVal(s[i])))
	}
	exph := 0
	if i < len(s) && s[i] == '.' {
		i++
		for ; i < len(s) && isHex(s[i]); i++ {
			ans = math.FMA(16, ans, float64(hexVal(s[i])))
			exph += 4
		}
	}
	expn := 0
	if i < len(s) && (s[i] == 'p' || s[i] == 'P') {
		i++
		neg := false
		if i < len(s) && (s[i] == '+' || s[i] == '-') {
			neg = s[i] == '-'
			i++
		}
		for ; i < len(s) && isDigit(int(s[i])); i++ {
			if expn < 100000 {
				expn = expn*10 + int(s[i]-'0')
			}
		}
		if neg {
			expn = -expn
		}
	}
	return math.Ldexp(ans, expn-exph)
}

func isHex(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func hexVal(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c-'a') + 10
	}
	return int(c-'A') + 10
}

// str lexes a quoted string or, with forSymbol, a backquoted name. The
// opening quote c has been read.
func (l *lexer) str(quote int, forSymbol bool, t *lexTok) {
	var out []byte
	nwide := 0
	oct, wide := false, false
	endOK := false
loop:
	for {
		c := l.getc()
		if c == eofByte {
			break
		}
		if c == quote {
			endOK = true
			break
		}
		if c == '\n' {
			l.ungetc()
			c = '\\'
		}
		if c == '\\' {
			c = l.getc()
			if c == eofByte {
				break loop
			}
			switch {
			case c >= '0' && c <= '7':
				octal := c - '0'
				c = l.getc()
				if c == eofByte {
					break loop
				}
				if c >= '0' && c <= '7' {
					octal = 8*octal + c - '0'
					c = l.getc()
					if c == eofByte {
						break loop
					}
					if c >= '0' && c <= '7' {
						octal = 8*octal + c - '0'
					} else {
						l.ungetc()
					}
				} else {
					l.ungetc()
				}
				if octal == 0 {
					l.raise("nul character not allowed")
				}
				if octal > 0xff {
					l.raise("\\%o exceeds maximum allowed octal value \\377", octal)
				}
				c = octal
				oct = true
			case c == 'x':
				val := 0
				for i := 0; i < 2; i++ {
					c = l.getc()
					if c == eofByte {
						break loop
					}
					ext := hexDigit(c)
					if ext < 0 {
						l.ungetc()
						if i == 0 {
							l.raise("'\\x' used without hex digits in character string")
						}
						break
					}
					val = 16*val + ext
				}
				if val == 0 {
					l.raise("nul character not allowed")
				}
				c = val
				oct = true
			case c == 'u' || c == 'U':
				maxd := 4
				if c == 'U' {
					maxd = 8
				}
				esc := byte(c)
				if forSymbol {
					if esc == 'U' {
						l.raise("\\Uxxxxxxxx sequences not supported inside backticks")
					}
					l.raise("\\uxxxx sequences not supported inside backticks")
				}
				c = l.getc()
				if c == eofByte {
					break loop
				}
				delim := false
				if c == '{' {
					delim = true
				} else {
					l.ungetc()
				}
				var val uint32
				for i := 0; i < maxd; i++ {
					c = l.getc()
					if c == eofByte {
						break loop
					}
					ext := hexDigit(c)
					if ext < 0 {
						l.ungetc()
						if i == 0 {
							l.raise("'\\%c' used without hex digits in character string", esc)
						}
						break
					}
					val = 16*val + uint32(ext)
				}
				if delim {
					c = l.getc()
					if c == eofByte {
						break loop
					}
					if c != '}' {
						if esc == 'u' {
							// R reports this one by line only
							panic(&Error{Line: l.line, Col: l.col, Lexical: true,
								Message: "invalid \\u{xxxx} sequence", text: "invalid \\u{xxxx} sequence (line " + strconv.Itoa(l.line) + ")"})
						}
						l.raise("invalid \\U{xxxxxxxx} sequence")
					}
				}
				if val == 0 {
					l.raise("nul character not allowed")
				}
				if val > 0x10FFFF {
					if delim {
						l.raise("invalid \\U{xxxxxxxx} value %6x", val)
					}
					l.raise("invalid \\Uxxxxxxxx value %6x", val)
				}
				out = utf8.AppendRune(out, rune(val))
				nwide++
				wide = true
				continue
			default:
				switch c {
				case 'a':
					c = '\a'
				case 'b':
					c = '\b'
				case 'f':
					c = '\f'
				case 'n':
					c = '\n'
				case 'r':
					c = '\r'
				case 't':
					c = '\t'
				case 'v':
					c = '\v'
				case '\\', '"', '\'', '`', ' ', '\n':
				default:
					l.raise("'\\%s' is an unrecognized escape in character string", string(rune(c)))
				}
			}
		} else if c >= 0x80 {
			r, n := l.mbNext(c)
			if n == -1 {
				break loop
			}
			if r >= 0x202A && r <= 0x2069 && !(r > 0x202E && r < 0x2066) {
				l.raise("bidi formatting not allowed, use escapes instead (\\u%04x)", r)
			}
			out = append(out, byte(c))
			for i := 0; i < n-1; i++ {
				c = l.getc()
				if c == eofByte {
					break loop
				}
				if c == '\n' {
					l.ungetc()
					c = '\\'
				}
				out = append(out, byte(c))
			}
			nwide++
			continue
		} else {
			nwide++
		}
		out = append(out, byte(c))
	}
	if !endOK {
		t.kind = tINCOMPLETE
		return
	}
	text := string(l.src[t.fb:l.off])
	if len(text) > 1000 {
		if forSymbol || !wide {
			text = "[" + strconv.Itoa(len(out)) + " chars quoted with '" + string(rune(quote)) + "']"
		} else {
			text = "[" + strconv.Itoa(nwide+1) + " wide chars quoted with '" + string(rune(quote)) + "']"
		}
	}
	t.text = text
	if forSymbol {
		if len(out) == 0 {
			panic(&Error{Message: "attempt to use zero-length variable name", Lexical: true, text: "attempt to use zero-length variable name"})
		}
		t.kind = tSYMBOL
		t.sym = string(out)
		return
	}
	if wide && oct {
		l.raise("mixing Unicode and octal/hex escapes in a string is not allowed")
	}
	if wide && nwide >= 10000 {
		l.raise("string containing Unicode escapes not in this locale\nis too long (max 10000 chars)")
	}
	t.kind = tSTR
	t.cst = &Const{Kind: String, Str: string(out)}
}

func hexDigit(c int) int {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	}
	return -1
}

// rawString lexes r"(...)" and its variants; r and the opening quote have
// been read.
func (l *lexer) rawString(quote byte, t *lexTok) {
	ndash := 0
	for l.nextchar('-') {
		ndash++
	}
	c := l.getc()
	var delim int
	switch c {
	case '(':
		delim = ')'
	case '[':
		delim = ']'
	case '{':
		delim = '}'
	case '|':
		delim = '|'
	default:
		l.raise("malformed raw string literal")
	}
	var out []byte
	wide := false
	nwide := 0
	endOK := false
	for {
		c = l.getc()
		if c == eofByte {
			break
		}
		if c == delim {
			nd := 0
			for nd < ndash && l.nextchar('-') {
				nd++
			}
			if nd == ndash && l.nextchar(quote) {
				endOK = true
				break
			}
			out = append(out, byte(delim))
			for i := 0; i < nd; i++ {
				out = append(out, '-')
			}
			nwide += 1 + nd
			continue
		}
		if c >= 0x80 {
			r, n := l.mbNext(c)
			if n == -1 {
				break
			}
			if r >= 0x202A && r <= 0x2069 && !(r > 0x202E && r < 0x2066) {
				l.raise("bidi formatting not allowed, use escapes instead (\\u%04x)", r)
			}
			out = append(out, byte(c))
			eof := false
			for i := 0; i < n-1; i++ {
				c = l.getc()
				if c == eofByte {
					eof = true
					break
				}
				out = append(out, byte(c))
			}
			if eof {
				break
			}
			nwide++
			continue
		}
		out = append(out, byte(c))
		nwide++
	}
	if !endOK {
		t.kind = tINCOMPLETE
		return
	}
	text := string(l.src[t.fb:l.off])
	if len(text) > 1000 {
		text = "[" + strconv.Itoa(len(out)) + " chars quoted with '" + string(rune(quote)) + "']"
	}
	_ = wide
	t.text = text
	t.kind = tSTR
	t.cst = &Const{Kind: String, Str: string(out)}
}

// ---- context-sensitive newlines (the part of R's lexer that sits between
// the characters and the grammar) -------------------------------------

func (l *lexer) top() byte { return l.ctx[len(l.ctx)-1] }

// contextStackSize is R's CONTEXTSTACK_SIZE: brackets, braces, parentheses and
// the if bodies inside them nest at most this deep.
const contextStackSize = 50

func (l *lexer) push(b byte) {
	if len(l.ctx)-1 >= contextStackSize-1 {
		l.raise("contextstack overflow")
	}
	l.ctx = append(l.ctx, b)
}

func (l *lexer) pop() {
	if len(l.ctx) > 1 {
		l.ctx = l.ctx[:len(l.ctx)-1]
	}
}

func (l *lexer) ifPush() {
	switch l.top() {
	case '{', '[', '(', 'i':
		l.push('i')
	}
}

func (l *lexer) ifPop() {
	if l.top() == 'i' {
		l.pop()
	}
}

// next returns the next token as the parser sees it, applying the rules for
// newlines, and records terminals.
func (l *lexer) next() lexTok {
again:
	t := l.tokenRecorded()
	if t.kind == '\n' {
		if l.eatLines || l.top() == '[' || l.top() == '(' {
			goto again
		}
		if l.top() == 'i' {
			for t.kind == '\n' {
				t = l.tokenRecorded()
			}
			if t.kind == '}' || t.kind == ')' || t.kind == ']' {
				for l.top() == 'i' {
					l.ifPop()
				}
				l.pop()
				return t
			}
			if t.kind == ',' {
				l.ifPop()
				return t
			}
			if t.kind == tELSE {
				l.eatLines = true
				l.ifPop()
				return t
			}
			l.ifPop()
			// push the token back; the newline ends the statement
			saved := t
			l.saved = &saved
			if t.rec >= 0 {
				// it will be recorded again when it is returned
				l.toks = l.toks[:len(l.toks)-1]
			}
			nl := lexTok{kind: '\n', cend: t.cend, rec: -1}
			nl.first = t.first
			return nl
		}
		return t
	}
	switch t.kind {
	case '+', '-', '*', '/', '^', tLT, tLE, tGE, tGT, tEQ, tNE, tOR, tAND, tOR2, tAND2,
		tPIPE, tPIPEBIND, tSPECIAL, tFUNCTION, tWHILE, tREPEAT, tFOR, tIN, '?', '!', ':',
		'~', '$', '@', tLEFTASSIGN, tRIGHTASSIGN, tEQASSIGN:
		l.eatLines = true
	case tIF:
		l.ifPush()
		l.eatLines = true
	case tELSE:
		l.ifPop()
		l.eatLines = true
	case ';', ',':
		l.ifPop()
	case tSYMBOL, tPLACEHOLDER, tSTR, tNUM, tNULL, tNEXT, tBREAK:
		l.eatLines = false
	case tLBB:
		if len(l.ctx)-1 >= contextStackSize-2 {
			l.raise("contextstack overflow")
		}
		l.push('[')
		l.push('[')
	case '[':
		l.push('[')
	case '{':
		l.push('{')
		l.eatLines = true
	case '(':
		l.push('(')
	case ']':
		for l.top() == 'i' {
			l.ifPop()
		}
		l.pop()
		l.eatLines = false
	case '}':
		for l.top() == 'i' {
			l.ifPop()
		}
		l.pop()
	case ')':
		for l.top() == 'i' {
			l.ifPop()
		}
		l.pop()
		l.eatLines = false
	}
	return t
}

// tokenRecorded lexes one token and records it as a terminal.
func (l *lexer) tokenRecorded() lexTok {
	t := l.token()
	if t.kind != '\n' && t.kind != tEOF && t.text != "" && t.kind != tERROR && t.kind != tINCOMPLETE {
		t.rec = len(l.toks)
		l.toks = append(l.toks, recTok{kind: t.kind.name(), text: t.text, first: t.first, last: t.last, fb: t.fb, lb: t.lb})
	}
	return t
}

func trimNL(s string) string { return strings.TrimRight(s, "\n") }

// isLineDirective reports whether a comment at the start of a line is a
// "#line N" directive: R's lexer records those under their own token name.
// (The directive's effect on the line numbers of later errors is not
// implemented.)
func isLineDirective(c []byte) bool {
	if len(c) < 6 || string(c[:5]) != "#line" {
		return false
	}
	i := 5
	for i < len(c) && (c[i] == ' ' || c[i] == '\t' || c[i] == '\f') {
		i++
	}
	return i < len(c) && isDigit(int(c[i]))
}
