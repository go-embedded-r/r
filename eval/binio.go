package eval

import (
	"encoding/binary"
	"math"
	"strings"
	"unicode/utf8"
)

// Binary input and output over the session's connections: rawConnection,
// readBin, writeBin, readChar, writeChar, seek. Written from the help pages of
// readBin, readChar, rawConnection and seek, and from what the oracle does;
// the bytes come from, and go to, the in-memory (or FileRoot-confined) file
// system and raw vectors of the session, nothing else.

// binReader is the data readBin and readChar read from: a raw vector (always
// read from its start), or the unread part of a connection.
type binReader struct {
	data []byte
	pos  int
	cn   *conn // nil for a raw vector
}

func (b *binReader) avail() int { return len(b.data) - b.pos }

// commit stores the read position back in the connection.
func (b *binReader) commit() {
	if b.cn != nil {
		b.cn.pos = b.pos
	}
}

// binSource resolves the con argument of readBin and readChar. what is the
// help page's name for the argument, textOK says whether a text-mode
// connection may be read (readChar), and the returned function releases a
// file opened for the call.
func (s *Session) binSource(c *callInfo, v Value, fname string, textOK bool) (*binReader, func()) {
	s.flushSinks()
	switch x := v.(type) {
	case *Raw:
		return &binReader{data: x.V}, func() {}
	case *String:
		if len(x.V) != 1 || x.V[0].NA {
			s.errorf(c.call, "invalid '%s' argument", "con")
		}
		cn := &conn{kind: kFile, class: "file", desc: x.V[0].S, path: s.abspath(x.V[0].S)}
		s.openConn(fileCall("con", "rb"), cn, "rb")
		return &binReader{data: cn.data, cn: cn}, func() { s.closeConn(cn, false) }
	}
	if !isConnValue(v) {
		if fname == "readChar" {
			s.errorf(c.call, "cannot read from this connection")
		}
		s.errorf(c.call, "can only read from a binary connection")
	}
	cn := s.getConn(c.call, v, "con")
	if fname == "readChar" && cn.kind == kFile {
		if !cn.open {
			s.openConn(c.call, cn, "rb")
			return &binReader{data: cn.data, cn: cn}, func() { s.closeConn(cn, false) }
		}
		if !cn.binary {
			s.warningf(c.call, "text connection used with readChar(), results may be incorrect")
		}
	}
	if cn.kind == kText && fname == "readChar" {
		s.warningf(c.call, "text connection used with readChar(), results may be incorrect")
		s.errorf(c.call, "'read' not enabled for this connection")
	}
	if !cn.open || (!cn.binary && !textOK) || cn.kind == kStd {
		if fname == "readChar" && cn.kind != kStd && !cn.open {
			s.errorf(c.call, "'read' not enabled for this connection")
		}
		s.errorf(c.call, "can only read from a binary connection")
	}
	if !cn.canRead {
		s.errorf(c.call, "cannot read from this connection")
	}
	return &binReader{data: cn.data, pos: cn.pos, cn: cn}, func() {}
}

// endianArg reads endian = c("little", "big", "swap"): true when the bytes
// are big-endian (little-endian is the machine's order).
func (s *Session) endianArg(c *callInfo, v Value) bool {
	if v == nil {
		return false
	}
	st, ok := v.(*String)
	if ok && len(st.V) > 1 {
		s.errorf(callOf(`if (!endian %in% c("big", "little", "swap")) stop("invalid 'endian' argument")`), "the condition has length > 1")
	}
	if !ok || len(st.V) < 1 || st.V[0].NA {
		s.errorf(c.call, "invalid '%s' argument", "endian")
	}
	switch st.V[0].S {
	case "little":
		return false
	case "big", "swap":
		return true
	}
	s.errorf(c.call, "invalid '%s' argument", "endian")
	return false
}

// binSize reads size = NA_integer_: def when it is NA or absent.
func (s *Session) binSize(c *callInfo, v Value, def int) (int, bool) {
	if v == nil {
		return def, false
	}
	var d []float64
	switch t := v.(type) {
	case *Logical, *Integer, *Double:
		d = asDoubles(v)
	case *String:
		for _, e := range t.V {
			f, ok := parseNumber(e)
			if !ok {
				s.warningf(c.call, "NAs introduced by coercion")
				f = naReal
			}
			d = append(d, f)
		}
	default:
		s.errorf(c.call, "invalid '%s' argument", "size")
	}
	if len(d) < 1 {
		s.errorf(c.call, "invalid '%s' argument", "size")
	}
	if math.IsNaN(d[0]) {
		return def, false
	}
	if d[0] > 1<<30 || d[0] < -(1<<30) {
		return int(d[0]/math.Abs(d[0])) << 30, true
	}
	return int(d[0]), true
}

func (s *Session) unknownSize(c *callInfo, size int) {
	s.errorf(c.call, "size %d is unknown on this machine", size)
}

func order(big bool) binary.ByteOrder {
	if big {
		return binary.BigEndian
	}
	return binary.LittleEndian
}

// whatMode is the mode readBin reads: the argument is one of the names, or
// stands for its type.
func (s *Session) whatMode(c *callInfo, v Value) string {
	if st, ok := v.(*String); ok {
		if len(st.V) > 1 {
			s.errorf(nil, "'length = %d' in coercion to 'logical(1)'", len(st.V))
		}
		if len(st.V) == 1 && !st.V[0].NA {
			switch st.V[0].S {
			case "numeric", "double":
				return "double"
			case "integer", "int":
				return "integer"
			case "logical", "complex", "character", "raw":
				return st.V[0].S
			}
		}
		return "character"
	}
	switch v.(type) {
	case *Integer:
		return "integer"
	case *Double:
		return "double"
	case *Logical:
		return "logical"
	case *Complex:
		return "complex"
	case *Raw:
		return "raw"
	}
	s.errorf(c.call, "invalid '%s' argument", "what")
	return ""
}

func init() {
	regBuiltin("readBin", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "con", "what", "n", "size", "signed", "endian")
		if a[0] == nil {
			s.errorf(c.call, "argument \"con\" is missing, with no default")
		}
		if a[1] == nil {
			s.errorf(c.call, "argument \"what\" is missing, with no default")
		}
		br, release := s.binSource(c, a[0], "readBin", false)
		defer release()
		mode := s.whatMode(c, a[1])
		n := 1
		if a[2] != nil {
			var d []float64
			switch t := a[2].(type) {
			case *Integer, *Double:
				d = asDoubles(a[2])
			case *String:
				for _, e := range t.V {
					f, _ := parseNumber(e)
					d = append(d, f)
				}
			default:
				s.errorf(c.call, "invalid '%s' argument", "n")
			}
			if len(d) < 1 || math.IsNaN(d[0]) || d[0] < 0 {
				s.errorf(c.call, "invalid '%s' argument", "n")
			}
			if d[0] > float64(s.opts.MaxVector) {
				s.errorf(c.call, "cannot allocate vector of length %.0f", d[0])
			}
			n = int(d[0])
		}
		padStrings := br.cn == nil
		signed := true
		if a[4] != nil {
			signed = s.boolArg(c, a[4], true, "signed")
		}
		big := s.endianArg(c, a[5])
		bo := order(big)
		defer br.commit()
		switch mode {
		case "integer":
			size, _ := s.binSize(c, a[3], 4)
			switch size {
			case 1, 2, 4, 8:
			default:
				s.unknownSize(c, size)
			}
			if !signed && size != 1 && size != 2 {
				s.warningf(c.call, "'signed = FALSE' is only valid for integers of sizes 1 and 2")
			}
			cnt := min(n, br.avail()/size)
			s.checkLen(c.call, cnt)
			out := make([]int32, cnt)
			for i := range out {
				b := br.data[br.pos : br.pos+size]
				br.pos += size
				switch size {
				case 1:
					if signed {
						out[i] = int32(int8(b[0]))
					} else {
						out[i] = int32(b[0])
					}
				case 2:
					v := bo.Uint16(b)
					if signed {
						out[i] = int32(int16(v))
					} else {
						out[i] = int32(v)
					}
				case 4:
					out[i] = int32(bo.Uint32(b))
				default:
					out[i] = int32(bo.Uint64(b))
				}
			}
			return &Integer{V: out}
		case "logical":
			size, _ := s.binSize(c, a[3], 4)
			switch size {
			case 1, 2, 4, 8:
			default:
				s.unknownSize(c, size)
			}
			cnt := min(n, br.avail()/size)
			s.checkLen(c.call, cnt)
			out := make([]int32, cnt)
			for i := range out {
				b := br.data[br.pos : br.pos+size]
				br.pos += size
				var v int32
				switch size {
				case 1:
					v = int32(int8(b[0]))
				case 2:
					v = int32(int16(bo.Uint16(b)))
				case 4:
					v = int32(bo.Uint32(b))
				default:
					v = int32(bo.Uint64(b))
				}
				out[i] = v // R keeps the integer it read (2 stays 2); only INT_MIN is NA
			}
			return &Logical{V: out}
		case "double":
			size, _ := s.binSize(c, a[3], 8)
			if size != 4 && size != 8 {
				s.unknownSize(c, size)
			}
			cnt := min(n, br.avail()/size)
			s.checkLen(c.call, cnt)
			out := make([]float64, cnt)
			for i := range out {
				b := br.data[br.pos : br.pos+size]
				br.pos += size
				if size == 4 {
					out[i] = float64(math.Float32frombits(bo.Uint32(b)))
				} else {
					out[i] = math.Float64frombits(bo.Uint64(b))
				}
			}
			return &Double{V: out}
		case "complex":
			if size, set := s.binSize(c, a[3], 16); set && size != 16 {
				s.errorf(c.call, "size changing is not supported for complex vectors")
			}
			cnt := min(n, br.avail()/16)
			s.checkLen(c.call, cnt)
			out := make([]complex128, cnt)
			for i := range out {
				b := br.data[br.pos : br.pos+16]
				br.pos += 16
				out[i] = complex(math.Float64frombits(bo.Uint64(b)), math.Float64frombits(bo.Uint64(b[8:])))
			}
			return &Complex{V: out}
		case "raw":
			if size, set := s.binSize(c, a[3], 1); set && size != 1 {
				s.errorf(c.call, "raw is always of size 1")
			}
			cnt := min(n, br.avail())
			s.checkLen(c.call, cnt)
			out := append([]byte(nil), br.data[br.pos:br.pos+cnt]...)
			br.pos += cnt
			return &Raw{V: out}
		}
		// character: NUL-terminated strings
		out := []Str{}
		for len(out) < n && (br.avail() > 0 || padStrings) {
			if br.avail() == 0 {
				// a raw vector is read as if it went on with empty strings
				s.checkLen(c.call, n)
				out = append(out, make([]Str, n-len(out))...)
				break
			}
			rest := br.data[br.pos:]
			i := 0
			for i < len(rest) && rest[i] != 0 {
				i++
			}
			if i == len(rest) {
				if br.cn != nil {
					s.warningf(c.call, "incomplete string at end of file has been discarded")
					br.pos = len(br.data)
					break
				}
			}
			out = append(out, Str{S: string(rest[:i])})
			br.pos += min(i+1, len(rest))
			s.checkLen(c.call, len(out))
		}
		return &String{V: out}
	})

	regBuiltin("writeBin", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "object", "con", "size", "endian", "useBytes")
		if a[0] == nil {
			s.errorf(c.call, "argument \"object\" is missing, with no default")
		}
		if a[1] == nil {
			s.errorf(c.call, "argument \"con\" is missing, with no default")
		}
		obj := a[0]
		switch obj.(type) {
		case *Logical, *Integer, *Double, *Complex, *String, *Raw:
			if isObject(obj) {
				s.errorf(c.call, "can only write vector objects")
			}
		default:
			s.errorf(c.call, "can only write vector objects")
		}
		big := s.endianArg(c, a[3])
		bo := order(big)
		size, sizeSet := s.binSize(c, a[2], 0)
		var buf []byte
		le := func(n int) []byte { buf = append(buf, make([]byte, n)...); return buf[len(buf)-n:] }
		switch x := obj.(type) {
		case *Logical, *Integer:
			if size == 0 && sizeSet {
				s.unknownSize(c, size)
			}
			if size == 0 {
				size = 4
			}
			iv := asInts(x)
			switch size {
			case 1, 2, 4, 8:
			default:
				s.unknownSize(c, size)
			}
			s.checkLen(c.call, len(iv)*size/8+1)
			buf = make([]byte, 0, len(iv)*size)
			for _, v := range iv {
				b := le(size)
				switch size {
				case 1:
					b[0] = byte(v)
				case 2:
					bo.PutUint16(b, uint16(v))
				case 4:
					bo.PutUint32(b, uint32(v))
				default:
					bo.PutUint64(b, uint64(int64(v)))
				}
			}
		case *Double:
			if size == 0 && sizeSet {
				s.unknownSize(c, size)
			}
			if size == 0 {
				size = 8
			}
			if size != 4 && size != 8 {
				s.unknownSize(c, size)
			}
			s.checkLen(c.call, len(x.V)*size/8+1)
			buf = make([]byte, 0, len(x.V)*size)
			for _, v := range x.V {
				b := le(size)
				if size == 4 {
					bo.PutUint32(b, math.Float32bits(float32(v)))
				} else {
					bo.PutUint64(b, math.Float64bits(v))
				}
			}
		case *Complex:
			if size != 0 && size != 16 {
				s.errorf(c.call, "size changing is not supported for complex vectors")
			}
			s.checkLen(c.call, len(x.V)*2)
			buf = make([]byte, 0, len(x.V)*16)
			for _, v := range x.V {
				b := le(16)
				bo.PutUint64(b, math.Float64bits(real(v)))
				bo.PutUint64(b[8:], math.Float64bits(imag(v)))
			}
		case *Raw:
			if size != 0 && size != 1 {
				s.errorf(c.call, "raw is always of size 1")
			}
			buf = append([]byte(nil), x.V...)
		case *String:
			total := 0
			for _, e := range x.V {
				total += len(e.S) + 1
			}
			s.checkLen(c.call, total/8+1)
			buf = make([]byte, 0, total)
			for _, e := range x.V {
				if e.NA {
					buf = append(buf, 'N', 'A', 0)
				} else {
					buf = append(buf, e.S...)
					buf = append(buf, 0)
				}
			}
		}
		return s.binSink(c, a[1], buf, "writeBin")
	})

	regBuiltin("writeChar", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "object", "con", "nchars", "eos", "useBytes")
		if a[0] == nil {
			s.errorf(c.call, "argument \"object\" is missing, with no default")
		}
		if a[1] == nil {
			s.errorf(c.call, "argument \"con\" is missing, with no default")
		}
		st, ok := a[0].(*String)
		if !ok {
			s.errorf(c.call, "can only write character objects")
		}
		useBytes := s.boolArg(c, a[4], false, "useBytes")
		var eos []byte
		hasEos := true
		if a[3] != nil {
			switch e := a[3].(type) {
			case nullT:
				hasEos = false
			case *String:
				if len(e.V) != 1 {
					s.errorf(c.call, "invalid '%s' argument", "sep")
				}
				eos = []byte(e.V[0].S)
				if e.V[0].NA {
					eos = []byte("NA")
				}
			default:
				s.errorf(c.call, "invalid '%s' argument", "sep")
			}
		}
		nch := make([]int, len(st.V))
		for i, e := range st.V {
			if useBytes {
				nch[i] = len(e.S)
			} else {
				nch[i] = utf8.RuneCountInString(e.S)
			}
		}
		nstr := len(st.V)
		if a[2] == nil {
			for _, e := range st.V {
				if e.NA {
					s.errorf(c.call, "invalid '%s' argument", "nchars")
				}
			}
		} else {
			var d []float64
			switch t := a[2].(type) {
			case *Logical, *Integer, *Double:
				d = asDoubles(a[2])
			case *String:
				for _, e := range t.V {
					f, ok := parseNumber(e)
					if !ok {
						s.warningf(c.call, "NAs introduced by coercion")
						f = naReal
					}
					d = append(d, f)
				}
			default:
				s.errorf(c.call, "invalid '%s' argument", "nchars")
			}
			if len(d) > len(st.V) {
				s.errorf(c.call, "'object' is too short")
			}
			nstr = len(d)
			nch = nch[:nstr]
			for i := range nch {
				if math.IsNaN(d[i]) || d[i] < 0 {
					s.errorf(c.call, "invalid '%s' argument", "nchars")
				}
				if d[i] > float64(s.opts.MaxVector) {
					s.errorf(c.call, "cannot allocate vector of length %.0f", d[i])
				}
				nch[i] = int(d[i])
			}
		}
		var buf []byte
		for i, e := range st.V[:nstr] {
			str := e.S
			if e.NA {
				str = "NA"
			}
			have := utf8.RuneCountInString(str)
			if useBytes {
				have = len(str)
			}
			want := nch[i]
			switch {
			case want <= have:
				if useBytes {
					buf = append(buf, str[:want]...)
				} else {
					k := 0
					for want > 0 {
						_, w := utf8.DecodeRuneInString(str[k:])
						k += w
						want--
					}
					buf = append(buf, str[:k]...)
				}
			default:
				s.warningf(c.call, "writeChar: more characters requested than are in the string - will zero-pad")
				buf = append(buf, str...)
				pad := want - have
				s.chargeBytes(pad)
				buf = append(buf, make([]byte, pad)...)
			}
			if hasEos {
				buf = append(buf, eos...)
				buf = append(buf, 0)
			}
			s.chargeBytes(len(buf))
		}
		return s.binSink(c, a[1], buf, "writeChar")
	})

	regBuiltin("readChar", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "con", "nchars", "useBytes")
		if a[0] == nil {
			s.errorf(c.call, "argument \"con\" is missing, with no default")
		}
		if a[1] == nil {
			s.errorf(c.call, "argument \"nchars\" is missing, with no default")
		}
		useBytes := s.boolArg(c, a[2], false, "useBytes")
		br, release := s.binSource(c, a[0], "readChar", true)
		defer release()
		defer br.commit()
		var counts []int
		switch x := a[1].(type) {
		case *Logical, *Integer, *Double:
			for _, f := range asDoubles(x) {
				switch {
				case math.IsNaN(f):
					s.errorf(c.call, "invalid '%s' argument", "nchars")
				case f >= 2147483648 || f <= -2147483649:
					s.warningf(c.call, "NAs introduced by coercion to integer range")
					s.errorf(c.call, "invalid '%s' argument", "nchars")
				case f < 0:
					s.errorf(c.call, "invalid '%s' argument", "nchars")
				}
				counts = append(counts, int(f))
			}
		case *String:
			for _, e := range x.V {
				f, ok := parseNumber(e)
				if !ok || math.IsNaN(f) {
					s.warningf(c.call, "NAs introduced by coercion")
					s.errorf(c.call, "invalid '%s' argument", "nchars")
				}
				if f < 0 {
					s.errorf(c.call, "invalid '%s' argument", "nchars")
				}
				counts = append(counts, int(f))
			}
		default:
			s.errorf(c.call, "invalid '%s' argument", "nchars")
		}
		out := []Str{}
		for _, n := range counts {
			if br.avail() <= 0 && n > 0 {
				break
			}
			rest := br.data[br.pos:]
			k := 0
			if useBytes {
				k = min(n, len(rest))
			} else {
				for i := 0; i < n && k < len(rest); i++ {
					_, w := utf8.DecodeRune(rest[k:])
					k += w
				}
			}
			chunk := rest[:k]
			br.pos += k
			if !useBytes && !utf8.Valid(chunk) {
				s.errorf(c.call, "invalid UTF-8 input in readChar()")
			}
			trim := chunk
			if br.cn != nil {
				// a connection's string stops at its first nul, with a warning
				if i := strings.IndexByte(string(chunk), 0); i >= 0 {
					s.warningf(c.call, "truncating string with embedded nuls")
					trim = chunk[:i]
				}
			} else {
				for len(trim) > 0 && trim[len(trim)-1] == 0 {
					trim = trim[:len(trim)-1]
				}
				if strings.IndexByte(string(trim), 0) >= 0 {
					s.errorf(c.call, "embedded nul in string: '%s'", nulEscape(trim))
				}
			}
			out = append(out, Str{S: string(trim)})
			s.checkLen(c.call, len(out))
			s.chargeBytes(len(trim))
		}
		return &String{V: out}
	})

	regBuiltin("rawConnection", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "object", "open")
		if a[0] == nil {
			s.errorf(c.call, "argument \"object\" is missing, with no default")
		}
		r, ok := a[0].(*Raw)
		if !ok {
			s.errorf(c.call, "invalid '%s' argument", "raw")
		}
		open := "r"
		if a[1] != nil {
			st, ok := a[1].(*String)
			if !ok || len(st.V) != 1 || st.V[0].NA {
				s.errorf(c.call, "invalid '%s' argument", "open")
			}
			open = st.V[0].S
		}
		cn := &conn{kind: kRawConn, class: "rawConnection", desc: descriptionOf(c, "object", a[0]), mode: open, open: true, binary: true}
		if strings.HasPrefix(open, "r") {
			cn.canRead = true
			cn.canWrite = strings.Contains(open, "+")
			cn.data = append([]byte(nil), r.V...)
			if cn.canWrite {
				cn.buf = append([]byte(nil), r.V...)
			}
		}
		if strings.HasPrefix(open, "w") || strings.HasPrefix(open, "a") {
			// the initial bytes are kept, whatever the mode
			cn.canWrite = true
			cn.buf = append([]byte(nil), r.V...)
			if strings.HasPrefix(open, "a") {
				cn.wpos = len(cn.buf)
			}
		}
		return connValue(s.newConn(cn))
	})
	regBuiltin("rawConnectionValue", func(s *Session, c *callInfo) Value {
		v := s.nargs1(c, "rawConnectionValue")
		if !isConnValue(v) {
			s.errorf(c.call, "'con' is not a rawConnection")
		}
		cn := s.getConn(c.call, v, "con")
		if cn.kind != kRawConn {
			s.errorf(c.call, "'con' is not a rawConnection")
		}
		if !cn.canWrite {
			s.errorf(c.call, "'con' is not an output rawConnection")
		}
		s.flushSinks()
		return &Raw{V: append([]byte(nil), cn.buf...)}
	})
	regBuiltin("isSeekable", func(s *Session, c *callInfo) Value {
		cn := s.getConn(c.call, s.nargs1(c, "isSeekable"), "con")
		return mkLogical(cn.kind == kRawConn || cn.kind == kFile)
	})
	connGeneric("seek", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "con", "where", "origin", "rw")
		cn := s.getConn(c.call, a[0], "con")
		if !cn.open {
			s.errorf(c.call, "connection is not open")
		}
		if cn.kind != kRawConn && cn.kind != kFile {
			s.errorf(c.call, "seek not enabled for this connection")
		}
		rw := s.optString(c, a[3], "", "rw")
		if !cn.canRead || (rw != "" && strings.HasPrefix("write", rw)) {
			if a[1] != nil && !math.IsNaN(asDoublesOrNil(a[1])[0]) {
				s.errorf(c.call, "seek on a write connection is not supported")
			}
			return mkDouble(float64(len(cn.buf)))
		}
		old := float64(cn.pos)
		if a[1] != nil {
			d := asDoublesOrNil(a[1])
			if len(d) != 1 {
				s.errorf(c.call, "invalid '%s' argument", "where")
			}
			if !math.IsNaN(d[0]) {
				origin := s.optString(c, a[2], "start", "origin")
				base := 0.0
				switch {
				case strings.HasPrefix("start", origin):
				case strings.HasPrefix("current", origin):
					base = float64(cn.pos)
				case strings.HasPrefix("end", origin):
					base = float64(len(cn.data))
				default:
					s.errorf(c.call, "invalid '%s' argument", "origin")
				}
				np := base + d[0]
				if np < 0 || np > float64(len(cn.data)) {
					s.errorf(c.call, "seek position out of range")
				}
				cn.pos = int(np)
			}
		}
		return mkDouble(old)
	})
}

// binSink delivers the bytes writeBin and writeChar produce: a raw vector
// target gets them back as the value; a file name or a connection receives
// them and the value is an invisible NULL.
func (s *Session) binSink(c *callInfo, target Value, buf []byte, fname string) Value {
	s.chargeBytes(len(buf))
	switch x := target.(type) {
	case *Raw:
		return &Raw{V: buf}
	case *String:
		if len(x.V) != 1 || x.V[0].NA {
			s.errorf(c.call, "invalid '%s' argument", "con")
		}
		cn := &conn{kind: kFile, class: "file", desc: x.V[0].S, path: s.abspath(x.V[0].S)}
		s.openConn(fileCall("con", "wb"), cn, "wb")
		s.connWriteBytes(c.call, cn, buf)
		s.closeConn(cn, false)
		return s.invisible(Null)
	}
	if !isConnValue(target) {
		if fname == "writeChar" {
			s.errorf(c.call, "'write' not enabled for this connection")
		}
		s.errorf(c.call, "can only write to a binary connection")
	}
	cn := s.getConn(c.call, target, "con")
	s.flushSinks()
	if fname == "writeChar" && cn.kind == kFile && !cn.open {
		// a file connection that is not open is opened for the call
		s.openConn(c.call, cn, "wb")
		defer s.closeConn(cn, false)
	}
	switch {
	case cn.kind == kStd || cn.kind == kText:
		if fname == "writeChar" {
			s.errorf(c.call, "'write' not enabled for this connection")
		}
		s.errorf(c.call, "can only write to a binary connection")
	case !cn.open || (!cn.binary && fname == "writeBin"):
		s.errorf(c.call, "can only write to a binary connection")
	case !cn.canWrite:
		s.errorf(c.call, "cannot write to this connection")
	}
	if len(buf) == 0 && cn.kind == kRawConn {
		s.warningf(c.call, "problem writing to connection")
	}
	s.connWriteBytes(c.call, cn, buf)
	return s.invisible(Null)
}

func (s *Session) connWriteBytes(call *Lang, cn *conn, b []byte) {
	if len(b) == 0 {
		return
	}
	s.connWrite(call, cn, string(b))
}
