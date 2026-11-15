package eval

import (
	"math"
	"strconv"
	"strings"
)

// scan (slice 18), written from ?scan and from what the oracle does with a
// series of probes; nothing of R's sources is read. The tokenizer works on the
// bytes of a connection (or of text=) and is bounded by what the file system
// holds (MaxFileBytes) and, for what it builds, by MaxVector.

// scanSpec is the dialect of one call.
type scanSpec struct {
	sep       byte // 0: white space
	quote     string
	dec       byte
	comment   byte // 0: none
	strip     []bool
	blankSkip bool
	fill      bool
	flush     bool
	multi     bool
	escapes   bool
	na        []string
}

// scanner tokenizes the unread bytes of a connection.
type scanner struct {
	d        []byte
	pos      int
	sp       *scanSpec
	unclosed bool // a quote was still open at the end of the data
	lastQ    byte // the quote character of the field being read
	// noQuote: with a separator, a field that is not read as text does not honour
	// quotes (a separator inside one still splits)
	noQuote bool
	// numWS: whitespace-separated, and the field is read as a number
	numWS bool
}

// shownText is the text of a field as R's conversion error shows it: the
// blanks removed, a quoted field between its quote characters.
func shownText(t scanTok) string {
	s := strings.Map(func(r rune) rune {
		if r == ' ' || r == '\t' {
			return -1
		}
		return r
	}, t.s)
	if t.quoted && t.qc != 0 {
		return string(t.qc) + s + string(t.qc)
	}
	return s
}

const (
	scBlank = iota // no field on this line (an empty line)
	scEnd          // end of the data, no field
	scField        // a field was read
)

type scanTok struct {
	s      string
	quoted bool
	qc     byte // the quote character of a quoted field
	kind   int
	eol    bool // the field was the last one of its line (kind scField), or the line ended (scBlank)
	eof    bool
}

// rc reads one byte: a line end is always "\n" (CR LF and a lone CR count).
func (sc *scanner) rc() (byte, bool) {
	if sc.pos >= len(sc.d) {
		return 0, false
	}
	c := sc.d[sc.pos]
	sc.pos++
	if c == '\r' {
		if sc.pos < len(sc.d) && sc.d[sc.pos] == '\n' {
			sc.pos++
		}
		return '\n', true
	}
	return c, true
}

func (sc *scanner) peek() (byte, bool) {
	save := sc.pos
	c, ok := sc.rc()
	sc.pos = save
	return c, ok
}

// rce is rc with the C escapes of allowEscapes processed; esc tells that the
// byte came from an escape.
func (sc *scanner) rce() (c byte, esc, ok bool) {
	c, ok = sc.rc()
	if !ok || c != '\\' || !sc.sp.escapes {
		return c, false, ok
	}
	n, ok2 := sc.rc()
	if !ok2 {
		return '\\', false, true
	}
	switch n {
	case 'n':
		return '\n', true, true
	case 't':
		return '\t', true, true
	case 'r':
		return '\r', true, true
	case 'a':
		return 7, true, true
	case 'b':
		return 8, true, true
	case 'f':
		return 12, true, true
	case 'v':
		return 11, true, true
	case '\\', '\'', '"', '`', ' ':
		return n, true, true
	case 'x':
		v, k := 0, 0
		for k < 2 {
			h, ok3 := sc.peek()
			if !ok3 || !isHex(h) {
				break
			}
			sc.rc()
			v = v*16 + hexVal(h)
			k++
		}
		if k == 0 {
			return 'x', true, true
		}
		return byte(v), true, true
	}
	if n >= '0' && n <= '7' {
		v := int(n - '0')
		for k := 1; k < 3; k++ {
			h, ok3 := sc.peek()
			if !ok3 || h < '0' || h > '7' {
				break
			}
			sc.rc()
			v = v*8 + int(h-'0')
		}
		return byte(v), true, true
	}
	// an unknown escape stays as it was written
	sc.pos--
	return '\\', false, true
}

func isHex(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
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

func (sc *scanner) isQuote(c byte) bool {
	return c != 0 && strings.IndexByte(sc.sp.quote, c) >= 0
}

// skipLines drops n lines, quotes ignored.
func (sc *scanner) skipLines(n int) {
	for ; n > 0; n-- {
		for {
			c, ok := sc.rc()
			if !ok {
				return
			}
			if c == '\n' {
				break
			}
		}
	}
}

// skipToEOL drops the rest of the line, the line end included.
func (sc *scanner) skipToEOL() {
	for {
		c, ok := sc.rc()
		if !ok || c == '\n' {
			return
		}
	}
}

// next reads one field. strip is the strip.white flag of its column.
func (sc *scanner) next(strip bool) scanTok {
	if sc.sp.sep == 0 {
		return sc.nextWhite()
	}
	return sc.nextSep(strip)
}

func (sc *scanner) nextWhite() scanTok {
	sp := sc.sp
	sc.lastQ = 0
	var c byte
	var esc, ok bool
	for {
		c, esc, ok = sc.rce()
		if !ok {
			return scanTok{kind: scEnd, eof: true}
		}
		if c != ' ' && c != '\t' {
			if !esc && c == sp.comment && c != 0 {
				sc.skipToEOL()
				return scanTok{kind: scBlank, eol: true}
			}
			if c == '\n' {
				return scanTok{kind: scBlank, eol: true}
			}
			break
		}
	}
	var b []byte
	quoted := false
	if !esc && sc.isQuote(c) {
		q := c
		sc.lastQ = q
		quoted = true
		for {
			c, esc, ok = sc.rce()
			if !ok {
				sc.unclosed = true
				break
			}
			if !esc && c == q {
				break
			}
			if !esc && c == '\\' {
				n, ok2 := sc.peek()
				if ok2 && (n == q) {
					sc.rc()
					b = append(b, n)
					continue
				}
			}
			b = append(b, c)
		}
		if sc.numWS {
			// a field read as a number does not end at its closing quote
			var rest []byte
			for {
				save := sc.pos
				d, e2, ok2 := sc.rce()
				if !ok2 {
					break
				}
				if !e2 && (d == ' ' || d == '\t' || d == '\n') {
					sc.pos = save
					break
				}
				rest = append(rest, d)
			}
			if len(rest) > 0 {
				b = append(append(append([]byte{q}, b...), q), rest...)
				quoted = false
			}
		}
	} else {
		b = append(b, c)
		for {
			save := sc.pos
			c, esc, ok = sc.rce()
			if !ok {
				break
			}
			if c == ' ' || c == '\t' || c == '\n' || (!esc && c == sp.comment && c != 0) {
				// the terminator is not part of the field: finishWhite reads it
				sc.pos = save
				break
			}
			b = append(b, c)
		}
	}
	return sc.finishWhite(string(b), quoted)
}

// finishWhite looks past the blanks after a field for the end of the line.
func (sc *scanner) finishWhite(s string, quoted bool) scanTok {
	for {
		save := sc.pos
		c, esc, ok := sc.rce()
		if !ok {
			return scanTok{s: s, quoted: quoted, qc: sc.lastQ, kind: scField, eol: true, eof: true}
		}
		switch {
		case c == ' ' || c == '\t':
		case c == '\n':
			return scanTok{s: s, quoted: quoted, qc: sc.lastQ, kind: scField, eol: true}
		case !esc && c == sc.sp.comment && c != 0:
			sc.skipToEOL()
			return scanTok{s: s, quoted: quoted, qc: sc.lastQ, kind: scField, eol: true}
		default:
			sc.pos = save
			return scanTok{s: s, quoted: quoted, qc: sc.lastQ, kind: scField}
		}
	}
}

func (sc *scanner) nextSep(strip bool) scanTok {
	sp := sc.sp
	sc.lastQ = 0
	var b []byte
	var inq byte
	quoted := false
	protect := 0
	started := false
	lineStart := sc.pos == 0 || sc.d[sc.pos-1] == '\n' || sc.d[sc.pos-1] == '\r'
	afterSep := sc.pos > 0 && sc.d[sc.pos-1] == sp.sep && sp.sep != '\n'
	nread := 0
	for {
		c, esc, ok := sc.rce()
		if !ok {
			if inq != 0 {
				sc.unclosed = true
			}
			if !started && !afterSep {
				return scanTok{kind: scEnd, eof: true}
			}
			return sc.endSep(b, protect, strip, quoted, true, true)
		}
		started = true
		nread++
		if inq != 0 {
			if !esc && c == inq {
				if n, ok2 := sc.peek(); ok2 && n == inq {
					sc.rc()
					b = append(b, c)
					continue
				}
				inq = 0
				protect = len(b)
				continue
			}
			b = append(b, c)
			continue
		}
		if !esc {
			switch {
			case c == '\n':
				if lineStart && nread == 1 {
					// an empty line
					return scanTok{kind: scBlank, eol: true}
				}
				return sc.endSep(b, protect, strip, quoted, true, false)
			case c == sp.sep:
				return sc.endSep(b, protect, strip, quoted, false, false)
			case sp.comment != 0 && c == sp.comment:
				sc.skipToEOL()
				if lineStart && nread == 1 {
					return scanTok{kind: scBlank, eol: true}
				}
				return sc.endSep(b, protect, strip, quoted, true, false)
			case !sc.noQuote && sc.isQuote(c):
				inq, quoted = c, true
				if sc.lastQ == 0 {
					sc.lastQ = c
				}
				continue
			}
		}
		if strip && len(b) == 0 && !quoted && (c == ' ' || c == '\t') {
			continue
		}
		b = append(b, c)
	}
}

func (sc *scanner) endSep(b []byte, protect int, strip, quoted, eol, eof bool) scanTok {
	if strip {
		for len(b) > protect && (b[len(b)-1] == ' ' || b[len(b)-1] == '\t') {
			b = b[:len(b)-1]
		}
	}
	return scanTok{s: string(b), quoted: quoted, qc: sc.lastQ, kind: scField, eol: eol, eof: eof}
}

// ---- conversion -----------------------------------------------------------

type scanKind int

const (
	skNull scanKind = iota
	skLogical
	skInteger
	skDouble
	skComplex
	skString
)

func scanKindOf(v Value) (scanKind, bool) {
	switch v.(type) {
	case nullT:
		return skNull, true
	case *Logical:
		return skLogical, true
	case *Integer:
		return skInteger, true
	case *Double:
		return skDouble, true
	case *Complex:
		return skComplex, true
	case *String:
		return skString, true
	}
	return 0, false
}

type scanCol struct {
	kind scanKind
	l    []int32
	d    []float64
	c    []complex128
	s    []Str
}

func (c *scanCol) len() int {
	switch c.kind {
	case skLogical, skInteger:
		return len(c.l)
	case skDouble:
		return len(c.d)
	case skComplex:
		return len(c.c)
	case skString:
		return len(c.s)
	}
	return 0
}

func (c *scanCol) pad() { c.add(scanTok{}, nil, nil, nil, nil) }

var scanExpected = map[scanKind]string{skLogical: "a logical", skInteger: "an integer", skDouble: "a real", skComplex: "a complex"}

// add converts one field and appends it; bad is called with the text when it
// does not convert.
func (c *scanCol) add(t scanTok, sp *scanSpec, dec *byte, bad func(kind scanKind, text string), _ any) {
	if c.kind == skNull {
		return
	}
	isNA := false
	txt := t.s
	if sp != nil {
		for _, n := range sp.na {
			if txt == n {
				isNA = true
				break
			}
		}
	}
	if c.kind != skString {
		txt = strings.TrimSpace(txt)
		if txt == "" {
			isNA = true
		}
	}
	switch c.kind {
	case skString:
		if isNA {
			c.s = append(c.s, Str{NA: true})
		} else {
			c.s = append(c.s, Str{S: txt})
		}
	case skLogical:
		v := int32(naInt)
		if !isNA {
			switch txt {
			case "T", "TRUE", "true", "True":
				v = 1
			case "F", "FALSE", "false", "False":
				v = 0
			default:
				bad(c.kind, shownText(t))
			}
		}
		c.l = append(c.l, v)
	case skInteger:
		v := int32(naInt)
		if !isNA {
			n, err := strconv.ParseInt(txt, 10, 64)
			if err != nil || n > math.MaxInt32 || n <= math.MinInt32 {
				bad(c.kind, shownText(t))
			} else {
				v = int32(n)
			}
		}
		c.l = append(c.l, v)
	case skDouble:
		v := naReal
		if !isNA {
			f, ok := scanNumber(txt, *dec)
			if !ok {
				bad(c.kind, shownText(t))
			}
			v = f
		}
		c.d = append(c.d, v)
	case skComplex:
		v := complex(naReal, naReal)
		if !isNA {
			if *dec != '.' {
				if strings.IndexByte(txt, '.') >= 0 {
					bad(c.kind, shownText(t))
				}
				txt = strings.Replace(txt, string(*dec), ".", 1)
			}
			z, ok := parseComplex(Str{S: txt})
			if !ok {
				bad(c.kind, shownText(t))
			}
			v = z
		}
		c.c = append(c.c, v)
	}
}

// scanNumber parses a real with the decimal mark dec.
func scanNumber(txt string, dec byte) (float64, bool) {
	if dec != '.' {
		if strings.IndexByte(txt, '.') >= 0 {
			return 0, false
		}
		txt = strings.Replace(txt, string(dec), ".", 1)
	}
	return parseNumber(Str{S: txt})
}

func (c *scanCol) value() Value {
	switch c.kind {
	case skLogical:
		return &Logical{V: c.l}
	case skInteger:
		return &Integer{V: c.l}
	case skDouble:
		return &Double{V: c.d}
	case skComplex:
		return &Complex{V: c.c}
	case skString:
		return &String{V: c.s}
	}
	return Null
}

func newScanCol(k scanKind) *scanCol { return &scanCol{kind: k} }

// ---- the builtin ----------------------------------------------------------

// scanInput finds the bytes scan reads: text=, a file name or a connection. It
// answers the connection and whether scan opened it (and so closes it).
func (s *Session) scanInput(c *callInfo, file, text Value) (cn *conn, ours bool) {
	return s.openInput(c, file, text, "r")
}

// openInput is scanInput with the mode of the file call R shows on a failure.
func (s *Session) openInput(c *callInfo, file, text Value, mode string, openCall ...*Lang) (cn *conn, ours bool) {
	if text != nil {
		ts, ok := text.(*String)
		if !ok {
			s.errorf(c.call, "invalid 'text' argument")
		}
		var sb strings.Builder
		for _, e := range ts.V {
			if e.NA {
				sb.WriteString("NA")
			} else {
				sb.WriteString(e.S)
			}
			sb.WriteByte('\n')
		}
		return &conn{kind: kText, class: "textConnection", desc: "text", open: true, canRead: true, data: []byte(sb.String())}, true
	}
	switch x := file.(type) {
	case nil:
		return s.io.conns[0], false
	case *String:
		if len(x.V) != 1 || x.V[0].NA {
			s.errorf(c.call, "invalid '%s' argument", "file")
		}
		if x.V[0].S == "" {
			return s.io.conns[0], false
		}
		cn = &conn{kind: kFile, class: "file", desc: x.V[0].S, path: s.abspath(x.V[0].S)}
		oc := fileCall("file", mode)
		if len(openCall) > 0 {
			oc = openCall[0]
		}
		s.openConn(oc, cn, "r")
		return cn, true
	}
	cn = s.getConn(c.call, file, "file")
	if !cn.open {
		s.openConn(c.call, cn, "rt")
		return cn, true
	}
	return cn, false
}

// scanSep reads a one-byte argument (sep, dec, comment.char).
func (s *Session) scanByteArg(c *callInfo, v Value, def byte, what, msg string, allowEmpty bool) byte {
	if v == nil {
		return def
	}
	st, ok := v.(*String)
	if !ok || len(st.V) != 1 || st.V[0].NA {
		s.errorf(c.call, "invalid '%s' argument", what)
	}
	switch len(st.V[0].S) {
	case 0:
		if allowEmpty {
			return 0
		}
	case 1:
		return st.V[0].S[0]
	}
	s.errorf(c.call, "%s", msg)
	return 0
}

func (s *Session) scanLimit(c *callInfo, v Value, def int, what string) int {
	if v == nil {
		return def
	}
	d := asDoublesOrNil(v)
	if len(d) < 1 || d[0] != d[0] {
		return def
	}
	return int(max(min(d[0], math.MaxInt32), -1))
}

// scanRun reads the records of the columns cols from sc, until the data ends or
// maxRec records (nlines lines) were read; call is the call an error shows.
func (s *Session) scanRun(call *Lang, sc *scanner, cols []*scanCol, maxRec, nlines int) (nrec int) {
	sp := sc.sp
	nc := len(cols)
	bad := func(kind scanKind, text string) {
		s.errorf(call, "scan() expected '%s', got '%s'", scanExpected[kind], text)
	}
	dec := sp.dec
	colsread, linesRead, skipped := 0, 0, 0
	natural := true
loop:
	for maxRec == 0 || nrec < maxRec {
		col := cols[colsread]
		sc.noQuote = sp.sep != 0 && col.kind != skString
		sc.numWS = sp.sep == 0 && col.kind != skString
		startPos := sc.pos
		t := sc.next(sp.strip[colsread%len(sp.strip)])
		if nc > 1 && colsread == 0 && t.kind == scField && t.eol && t.s == "" && !t.quoted && sp.sep != 0 && startPos > 0 && sc.d[startPos-1] == sp.sep {
			// a separator that ends a line leaves no field at the start of a record
			linesRead++
			if nlines > 0 && linesRead >= nlines {
				natural = false
				break loop
			}
			continue
		}
		switch t.kind {
		case scEnd:
			break loop
		case scBlank:
			// in the middle of a record an empty line is an empty field
			if sp.blankSkip && colsread == 0 {
				// it counts toward nlines, not toward the line number of an error
				linesRead++
				skipped++
				if nlines > 0 && linesRead >= nlines {
					natural = false
					break loop
				}
				continue
			}
			t = scanTok{kind: scField, eol: true}
		}
		col.add(t, sp, &dec, bad, nil)
		if n := col.len(); n > s.opts.MaxVector {
			s.errorf(call, "cannot allocate vector of length %d", n)
		}
		colsread++
		if colsread == nc {
			nrec++
			colsread = 0
			if sp.flush && !t.eol {
				for {
					u := sc.next(false)
					if u.kind == scEnd || u.eol {
						t.eol = u.kind != scEnd
						break
					}
				}
			}
		}
		if t.eol {
			linesRead++
			if colsread > 0 && (sp.fill || !sp.multi) {
				if !sp.fill {
					s.errorf(call, "line %d did not have %d elements", linesRead-skipped, nc)
				}
				for ; colsread < nc; colsread++ {
					cols[colsread].pad()
				}
				nrec++
				colsread = 0
			}
			if nlines > 0 && linesRead >= nlines {
				natural = false
				break loop
			}
		}
	}
	if maxRec > 0 && nrec >= maxRec {
		natural = false
	}
	if sc.unclosed {
		s.warningf(call, "EOF within quoted string")
	}
	if colsread > 0 {
		missing := nc - colsread
		if !natural || missing >= 2 {
			s.warningf(call, "number of items read is not a multiple of the number of columns")
		}
		for ; colsread < nc; colsread++ {
			cols[colsread].pad()
		}
		nrec++
	}
	return nrec
}

func init() {
	regBuiltin("scan", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "file", "what", "nmax", "n", "sep", "quote", "dec", "skip", "nlines", "na.strings",
			"flush", "fill", "strip.white", "quiet", "blank.lines.skip", "multi.line", "comment.char", "allowEscapes",
			"fileEncoding", "encoding", "text", "skipNul")
		sp := &scanSpec{dec: '.', blankSkip: true, multi: true, na: []string{"NA"}}
		sp.sep = s.scanByteArg(c, a[4], 0, "sep", "invalid 'sep' value: must be one byte", true)
		sp.quote = "'\""
		if sp.sep == '\n' {
			sp.quote = ""
		}
		if a[5] != nil {
			q, ok := a[5].(*String)
			if !ok || len(q.V) != 1 {
				s.errorf(c.call, "invalid quote symbol set")
			}
			sp.quote = q.V[0].S
		}
		sp.dec = s.scanByteArg(c, a[6], '.', "dec", "invalid decimal separator: must be one byte", false)
		sp.comment = s.scanByteArg(c, a[16], 0, "comment.char", "invalid 'comment.char' argument", true)
		sp.flush = s.boolArg(c, a[10], false, "flush")
		sp.fill = s.boolArg(c, a[11], false, "fill")
		quiet := s.boolArg(c, a[13], false, "quiet")
		sp.blankSkip = s.boolArg(c, a[14], true, "blank.lines.skip")
		sp.multi = s.boolArg(c, a[15], true, "multi.line")
		sp.escapes = s.boolArg(c, a[17], false, "allowEscapes")
		if a[9] != nil {
			switch x := a[9].(type) {
			case nullT:
				sp.na = nil
			case *String:
				sp.na = sp.na[:0]
				for _, e := range x.V {
					if !e.NA {
						sp.na = append(sp.na, e.S)
					}
				}
			default:
				s.errorf(c.call, "invalid '%s' argument", "na.strings")
			}
		}
		if a[2] != nil && a[3] != nil {
			s.errorf(c.call, "either specify 'nmax' or 'n', but not both.")
		}
		nmax := s.scanLimit(c, a[2], 0, "nmax")
		nitems := s.scanLimit(c, a[3], 0, "n")
		skip := s.scanLimit(c, a[7], 0, "skip")
		nlines := s.scanLimit(c, a[8], 0, "nlines")

		// the columns
		var cols []*scanCol
		var names *String
		isList := false
		switch w := a[1].(type) {
		case nil:
			cols = []*scanCol{newScanCol(skDouble)}
		case *List:
			isList = true
			names = namesOf(w)
			for _, e := range w.V {
				k, ok := scanKindOf(e)
				if !ok {
					s.errorf(c.call, "unimplemented type '%s' in 'extractItem'\n", typeName(e))
				}
				cols = append(cols, newScanCol(k))
			}
			if len(cols) == 0 {
				s.errorf(c.call, "invalid '%s' argument", "what")
			}
		default:
			k, ok := scanKindOf(w)
			if !ok || k == skNull {
				s.errorf(c.call, "invalid '%s' argument", "what")
			}
			cols = []*scanCol{newScanCol(k)}
		}
		nc := len(cols)
		sp.strip = []bool{false}
		if a[12] != nil {
			l, ok := a[12].(*Logical)
			if !ok || len(l.V) == 0 || (len(l.V) != 1 && len(l.V) != nc) {
				s.errorf(c.call, "invalid '%s' length", "strip.white")
			}
			sp.strip = make([]bool, len(l.V))
			for i, v := range l.V {
				sp.strip[i] = v == 1
			}
		}

		cn, ours := s.scanInput(c, a[0], a[20])
		if !cn.canRead {
			s.errorf(c.call, "cannot read from this connection")
		}
		s.flushSinks()
		sc := &scanner{d: cn.data, pos: cn.pos, sp: sp}
		defer func() {
			if ours {
				s.closeConn(cn, false)
			}
		}()
		sc.skipLines(skip)

		maxRec := 0
		switch {
		case isList && nmax > 0:
			maxRec = nmax
		case !isList && nmax > 0:
			maxRec = nmax
		case !isList && nitems > 0:
			maxRec = nitems
		}
		nrec := s.scanRun(c.call, sc, cols, maxRec, nlines)
		if !ours {
			cn.pos = sc.pos
		}
		if !quiet {
			unit := "item"
			if isList {
				unit = "record"
			}
			if nrec != 1 {
				unit += "s"
			}
			s.stderr.WriteString("Read " + strconv.Itoa(nrec) + " " + unit + "\n")
		}
		if !isList {
			return cols[0].value()
		}
		out := &List{V: make([]Value, nc)}
		for i, cl := range cols {
			out.V[i] = cl.value()
		}
		if names != nil {
			out.Attr = namesAttr(names.V)
		}
		return out
	})
}
