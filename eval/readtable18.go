package eval

import (
	"strconv"
	"strings"
)

// type.convert, count.fields and read.table (slice 18), written from the help
// pages of read.table, type.convert and count.fields and from probes of the
// oracle. Nothing of R's sources is read.

// ---- type.convert ---------------------------------------------------------

type tcOpts struct {
	na         []string
	asIs       bool
	dec        byte
	numerals   string
	tryLogical bool
}

func (o *tcOpts) isNA(e Str, trim bool) bool {
	if e.NA {
		return true
	}
	for _, n := range o.na {
		if e.S == n {
			return true
		}
	}
	return trim && strings.TrimSpace(e.S) == ""
}

// sigDigits counts the significant digits of a decimal constant.
func sigDigits(t string) int {
	n, started := 0, false
	for _, c := range t {
		switch {
		case c == 'e' || c == 'E':
			return n
		case c >= '1' && c <= '9':
			started = true
			n++
		case c == '0' && started:
			n++
		}
	}
	return n
}

// typeConvertStrings is type.convert on a character vector: the first of
// logical, integer, double, complex, character (or factor) that every element
// fits.
func (s *Session) typeConvertStrings(c *callInfo, x []Str, o *tcOpts) Value {
	allNA := true
	isLgl, isInt, isDbl, isCplx := o.tryLogical, true, true, true
	var lossy []string
	for _, e := range x {
		if o.isNA(e, true) {
			continue
		}
		allNA = false
		t := strings.TrimSpace(e.S)
		if isLgl {
			switch e.S {
			case "T", "F", "TRUE", "FALSE":
			default:
				isLgl = false
			}
		}
		if isInt {
			n, err := strconv.ParseInt(strings.TrimLeft(e.S, " \t"), 10, 64)
			if err != nil || n > 2147483647 || n <= -2147483648 {
				isInt = false
			}
		}
		if isDbl {
			if _, ok := scanNumber(t, o.dec); !ok {
				isDbl = false
			} else if sigDigits(t) > 15 && o.numerals != "allow.loss" {
				if o.numerals == "no.loss" {
					isDbl, isCplx = false, false
				} else {
					lossy = append(lossy, e.S)
				}
			}
		}
		if isCplx {
			u := t
			if o.dec != '.' {
				if strings.IndexByte(u, '.') >= 0 {
					isCplx = false
					continue
				}
				u = strings.Replace(u, string(o.dec), ".", 1)
			}
			if _, ok := parseComplex(Str{S: u}); !ok {
				isCplx = false
			}
		}
		if !isLgl && !isInt && !isDbl && !isCplx {
			break
		}
	}
	n := len(x)
	switch {
	case allNA && !o.tryLogical:
		out := make([]int32, n)
		for i := range out {
			out[i] = naInt
		}
		return &Integer{V: out}
	case allNA:
		out := make([]int32, n)
		for i := range out {
			out[i] = naInt
		}
		return &Logical{V: out}
	case isLgl:
		out := make([]int32, n)
		for i, e := range x {
			switch {
			case o.isNA(e, true):
				out[i] = naInt
			case e.S == "T" || e.S == "TRUE":
				out[i] = 1
			}
		}
		return &Logical{V: out}
	case isInt:
		out := make([]int32, n)
		for i, e := range x {
			if o.isNA(e, true) {
				out[i] = naInt
				continue
			}
			v, _ := strconv.ParseInt(strings.TrimLeft(e.S, " \t"), 10, 32)
			out[i] = int32(v)
		}
		return &Integer{V: out}
	case isDbl:
		for _, l := range lossy {
			s.warningf(c.call, "accuracy loss in conversion from \"%s\" to numeric", l)
		}
		out := make([]float64, n)
		for i, e := range x {
			if o.isNA(e, true) {
				out[i] = naReal
				continue
			}
			out[i], _ = scanNumber(strings.TrimSpace(e.S), o.dec)
		}
		return &Double{V: out}
	case isCplx:
		out := make([]complex128, n)
		for i, e := range x {
			if o.isNA(e, true) {
				out[i] = complex(naReal, naReal)
				continue
			}
			u := strings.TrimSpace(e.S)
			if o.dec != '.' {
				u = strings.Replace(u, string(o.dec), ".", 1)
			}
			out[i], _ = parseComplex(Str{S: u})
		}
		return &Complex{V: out}
	}
	out := make([]Str, n)
	for i, e := range x {
		if o.isNA(e, false) {
			out[i] = Str{NA: true}
		} else {
			out[i] = e
		}
	}
	res := &String{V: out}
	if o.asIs {
		return res
	}
	return s.callNamed(c, "factor", []Value{res}, []string{""})
}

// tcArgs reads the options of type.convert from the matched formals.
func (s *Session) tcArgs(c *callInfo, na, asIs, dec, numerals, tryLogical Value) *tcOpts {
	o := &tcOpts{na: []string{"NA"}, dec: '.', numerals: "allow.loss", tryLogical: true}
	if na != nil {
		st, ok := na.(*String)
		if !ok {
			s.errorf(c.call, "invalid '%s' argument", "na.strings")
		}
		o.na = o.na[:0]
		for _, e := range st.V {
			if !e.NA {
				o.na = append(o.na, e.S)
			}
		}
	}
	if asIs == nil {
		s.warningf(c.call, "'as.is' should be specified by the caller; using TRUE")
		o.asIs = true
	} else {
		o.asIs = s.boolArg(c, asIs, true, "as.is")
	}
	o.dec = s.scanByteArg(c, dec, '.', "dec", "invalid decimal separator: must be one byte", false)
	if numerals != nil {
		st, ok := numerals.(*String)
		if !ok || len(st.V) == 0 {
			s.errorf(c.call, "invalid '%s' argument", "numerals")
		}
		found := false
		for _, m := range []string{"allow.loss", "warn.loss", "no.loss"} {
			if st.V[0].S == m {
				found = true
			}
		}
		if !found {
			s.errorf(&Lang{Fun: Intern("match.arg"), Args: []Arg{{Value: Intern("numerals")}}},
				"'arg' should be one of “allow.loss”, “warn.loss”, “no.loss”")
		}
		o.numerals = st.V[0].S
	}
	o.tryLogical = s.boolArg(c, tryLogical, true, "tryLogical")
	return o
}

func init() {
	regBuiltin("..type.convert.default", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "na.strings", "as.is", "dec", "numerals", "tryLogical", "...")
		o := s.tcArgs(c, a[1], a[2], a[3], a[4], a[5])
		x := a[0]
		if x == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		var strs []Str
		if st, ok := x.(*String); ok && classAttr(x) == nil {
			strs = st.V
		} else {
			cv := s.callNamed(c, "as.character", []Value{x}, []string{""})
			st, ok := cv.(*String)
			if !ok {
				s.errorf(c.call, "invalid argument")
			}
			strs = st.V
		}
		res := s.typeConvertStrings(c, strs, o)
		// the dimensions stay, the names do not
		var keep []Attr
		for _, at := range attrsOf(x) {
			if at.Name == "dim" || at.Name == "dimnames" {
				keep = append(keep, at)
			}
		}
		if keep != nil {
			return withAttr(res, append(attrsOf(res), keep...))
		}
		return res
	})
	typeConvertEach := func(name string) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			if len(c.args) == 0 {
				s.errorf(c.call, "argument \"x\" is missing, with no default")
			}
			x, ok := c.args[0].(*List)
			if !ok {
				s.errorf(c.call, "invalid argument")
			}
			out := &List{V: make([]Value, len(x.V)), Attr: x.Attr}
			if name == "..type.convert.data.frame" {
				// the class comes last, after the row names
				out.Attr = nil
				var cls []Attr
				for _, at := range x.Attr {
					if at.Name == "class" {
						cls = append(cls, at)
					} else {
						out.Attr = append(out.Attr, at)
					}
				}
				out.Attr = append(out.Attr, cls...)
			}
			for i, e := range x.V {
				args := append([]Value{e}, c.args[1:]...)
				names := append([]string{""}, c.names[1:]...)
				out.V[i] = s.callNamed(c, "type.convert", args, names)
			}
			return out
		})
	}
	typeConvertEach("..type.convert.list")
	typeConvertEach("..type.convert.data.frame")
}

// ---- count.fields ---------------------------------------------------------

func init() {
	regBuiltin("count.fields", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "file", "sep", "quote", "skip", "blank.lines.skip", "comment.char")
		sp := &scanSpec{dec: '.', blankSkip: true, multi: true}
		sp.sep = s.scanByteArg(c, a[1], 0, "sep", "invalid 'sep' value: must be one byte", true)
		sp.quote = "'\""
		if sp.sep == '\n' {
			sp.quote = ""
		}
		if a[2] != nil {
			q, ok := a[2].(*String)
			if !ok || len(q.V) != 1 {
				s.errorf(c.call, "invalid quote symbol set")
			}
			sp.quote = q.V[0].S
		}
		skip := s.scanLimit(c, a[3], 0, "skip")
		sp.blankSkip = s.boolArg(c, a[4], true, "blank.lines.skip")
		sp.comment = '#'
		if a[5] != nil {
			sp.comment = s.scanByteArg(c, a[5], '#', "comment.char", "invalid 'comment.char' argument", true)
		}
		sp.strip = []bool{false}
		if a[0] == nil {
			s.errorf(c.call, "argument \"file\" is missing, with no default")
		}
		cn, ours := s.openInput(c, a[0], nil, "r", c.call)
		if !cn.canRead {
			s.errorf(c.call, "cannot read from this connection")
		}
		s.flushSinks()
		defer func() {
			if ours {
				s.closeConn(cn, false)
			}
		}()
		sc := &scanner{d: cn.data, pos: cn.pos, sp: sp}
		sc.skipLines(skip)
		var out []int32
		count := 0
		lines := 0 // physical lines consumed since the start of the record
		line := 1
		for {
			start := sc.pos
			t := sc.next(false)
			nl := newlinesIn(sc.d[start:sc.pos])
			if t.kind == scEnd {
				break
			}
			if sc.unclosed {
				s.errorf(c.call, "quoted string on line %d terminated by EOF", line+newlinesIn(sc.d[start:firstNonBlank(sc.d, start)]))
			}
			if t.kind == scBlank {
				line++
				if !sp.blankSkip {
					out = append(out, 0)
				}
				continue
			}
			count++
			if t.eol {
				// the terminating line end is not inside the field
				extra := nl
				if !t.eof {
					extra--
				}
				for i := 0; i < extra; i++ {
					out = append(out, naInt)
				}
				out = append(out, int32(count))
				if len(out) > s.opts.MaxVector {
					s.errorf(c.call, "cannot allocate vector of length %d", len(out))
				}
				line += lines + extra + 1
				count, lines = 0, 0
			} else {
				lines += nl
			}
		}
		if out == nil {
			return Null
		}
		return &Integer{V: out}
	})
}

func newlinesIn(b []byte) int {
	n := 0
	for i, c := range b {
		if c == '\n' || (c == '\r' && (i+1 >= len(b) || b[i+1] != '\n')) {
			n++
		}
	}
	return n
}

func firstNonBlank(d []byte, i int) int {
	for i < len(d) && (d[i] == ' ' || d[i] == '\t' || d[i] == '\n' || d[i] == '\r') {
		i++
	}
	return i
}

// ---- read.table -----------------------------------------------------------

const readTableScanCall = "scan(file = file, what = what, sep = sep, quote = quote, dec = dec, nmax = nrows, skip = 0, na.strings = na.strings, quiet = TRUE, fill = fill, strip.white = strip.white, blank.lines.skip = blank.lines.skip, multi.line = FALSE, comment.char = comment.char, allowEscapes = allowEscapes, flush = flush, encoding = encoding, skipNul = skipNul)"

var readTableScanLang *Lang

var readTableFormals = []string{"file", "header", "sep", "quote", "dec", "fill", "row.names", "col.names", "as.is",
	"tryLogical", "na.strings", "colClasses", "nrows", "skip", "check.names", "strip.white", "blank.lines.skip",
	"comment.char", "allowEscapes", "flush", "stringsAsFactors", "fileEncoding", "encoding", "text", "skipNul"}

const readCSVCall = "read.table(file = file, header = header, sep = sep, quote = quote, dec = dec, fill = fill, comment.char = comment.char, ...)"

func init() {
	regBuiltin("read.table", func(s *Session, c *callInfo) Value {
		return s.readTable(c, c.call, s.matchFormals(c, readTableFormals...))
	})
	// read.csv and its relatives: read.table with other defaults, and the call
	// R's errors show
	readCSV := func(name, sep, dec string) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			fixed := []string{"file", "header", "sep", "quote", "dec", "fill", "comment.char"}
			a := make([]Value, len(readTableFormals))
			idx := func(n string) int {
				for i, f := range readTableFormals {
					if f == n {
						return i
					}
				}
				return -1
			}
			a[1], a[2], a[3], a[4], a[5], a[17] = mkLogical(true), mkString(sep), mkString("\""), mkString(dec), mkLogical(true), mkString("")
			given := map[int]bool{}
			var unused []string
			next := 0
			for i, v := range c.args {
				k := -1
				if n := c.names[i]; n != "" {
					k = idx(n)
				} else {
					for next < len(fixed) && given[idx(fixed[next])] {
						next++
					}
					if next < len(fixed) {
						k = idx(fixed[next])
					}
				}
				if k < 0 || given[k] {
					unused = append(unused, s.deparseValue(v)[0])
					continue
				}
				given[k], a[k] = true, v
			}
			if len(unused) > 0 {
				s.errorf(c.call, "unused argument (%s)", strings.Join(unused, ", "))
			}
			return s.readTable(c, callOf(readCSVCall), a)
		})
	}
	readCSV("read.csv", ",", ".")
	readCSV("read.csv2", ";", ",")
	readCSV("read.delim", "\t", ".")
	readCSV("read.delim2", "\t", ",")
}

func (s *Session) readTable(c0 *callInfo, call *Lang, a []Value) Value {
	cc := *c0
	cc.call = call
	c := &cc
	{
		if a[0] == nil && a[23] == nil {
			s.errorf(c.call, "argument \"file\" is missing, with no default")
		}
		if readTableScanLang == nil {
			readTableScanLang = callOf(readTableScanCall)
		}
		sp := &scanSpec{dec: '.', blankSkip: true, multi: false, na: []string{"NA"}}
		sp.sep = s.scanByteArg(c, a[2], 0, "sep", "invalid 'sep' value: must be one byte", true)
		sp.quote = "\"'"
		if a[3] != nil {
			q, ok := a[3].(*String)
			if !ok || len(q.V) != 1 {
				s.errorf(c.call, "invalid quote symbol set")
			}
			sp.quote = q.V[0].S
		}
		sp.dec = s.scanByteArg(c, a[4], '.', "dec", "invalid decimal separator: must be one byte", false)
		sp.blankSkip = s.boolArg(c, a[16], true, "blank.lines.skip")
		sp.fill = s.boolArg(c, a[5], !sp.blankSkip, "fill")
		sp.comment = s.scanByteArg(c, a[17], '#', "comment.char", "invalid 'comment.char' argument", true)
		sp.escapes = s.boolArg(c, a[18], false, "allowEscapes")
		sp.flush = s.boolArg(c, a[19], false, "flush")
		stringsAsFactors := s.boolArg(c, a[20], false, "stringsAsFactors")
		checkNames := s.boolArg(c, a[14], true, "check.names")
		tryLogical := s.boolArg(c, a[9], true, "tryLogical")
		sp.strip = []bool{s.boolArg(c, a[15], false, "strip.white")}
		if a[10] != nil {
			st, ok := a[10].(*String)
			if !ok {
				s.errorf(c.call, "invalid '%s' argument", "na.strings")
			}
			sp.na = sp.na[:0]
			for _, e := range st.V {
				if !e.NA {
					sp.na = append(sp.na, e.S)
				}
			}
		}
		nrows := s.scanLimit(c, a[12], -1, "nrows")
		skip := s.scanLimit(c, a[13], 0, "skip")

		cn, ours := s.openInput(c, a[0], a[23], "rt")
		if !cn.canRead {
			s.errorf(c.call, "cannot read from this connection")
		}
		s.flushSinks()
		defer func() {
			if ours {
				s.closeConn(cn, false)
			}
		}()
		sc := &scanner{d: cn.data, pos: cn.pos, sp: sp}
		sc.skipLines(skip)

		// the first five lines give the number of columns and the header
		pre := &scanner{d: sc.d, pos: sc.pos, sp: sp}
		var first []scanTok
		var afterFirst int
		counts := []int{}
		cur := 0
	preLoop:
		for len(counts) < 5 {
			t := pre.next(sp.strip[0])
			switch t.kind {
			case scEnd:
				if cur > 0 {
					counts = append(counts, cur)
				}
				break preLoop
			case scBlank:
				if sp.blankSkip {
					continue
				}
				t = scanTok{kind: scField, eol: true}
			}
			if len(counts) == 0 {
				first = append(first, t)
			}
			cur++
			if cur > s.opts.MaxVector {
				s.errorf(c.call, "cannot allocate vector of length %d", cur)
			}
			if t.eol {
				if len(counts) == 0 {
					afterFirst = pre.pos
				}
				counts = append(counts, cur)
				cur = 0
			}
		}
		if len(counts) == 0 {
			s.errorf(c.call, "no lines available in input")
		}
		col1, cols := counts[0], 0
		for _, n := range counts {
			cols = max(cols, n)
		}
		var colNames []string
		userNames := false
		if a[7] != nil {
			st, ok := a[7].(*String)
			if !ok {
				s.errorf(c.call, "invalid '%s' argument", "col.names")
			}
			for _, e := range st.V {
				colNames = append(colNames, e.S)
			}
			userNames = true
			cols = max(cols, len(colNames))
		}
		header := col1 == cols-1
		if a[1] != nil {
			header = s.boolArg(c, a[1], false, "header")
		}
		rlabp := cols-col1 == 1 && header
		if header {
			sc.pos = afterFirst
		}
		if !userNames {
			if header {
				for _, t := range first {
					colNames = append(colNames, t.s)
				}
			} else {
				for i := 1; i <= cols; i++ {
					colNames = append(colNames, "V"+strconv.Itoa(i))
				}
			}
			if rlabp {
				colNames = append([]string{"row.names"}, colNames...)
			}
		}
		if len(colNames) < cols {
			s.errorf(c.call, "more columns than column names")
		}
		if checkNames {
			mk := s.callNamed(c, "make.names", []Value{strVec(colNames), mkLogical(true)}, []string{"", "unique"})
			colNames = colNames[:0]
			for _, e := range mk.(*String).V {
				colNames = append(colNames, e.S)
			}
		}
		colNames = colNames[:cols]

		// the classes of the columns
		classes := make([]string, cols)
		if a[11] != nil {
			cc, ok := a[11].(*String)
			if !ok {
				s.errorf(c.call, "invalid '%s' argument", "colClasses")
			}
			if nm := namesOf(cc); nm != nil {
				for i, e := range cc.V {
					for j, n := range colNames {
						if n == nm.V[i].S && !e.NA {
							classes[j] = e.S
						}
					}
				}
			} else if len(cc.V) > 0 {
				for j := range classes {
					if e := cc.V[j%len(cc.V)]; !e.NA {
						classes[j] = e.S
					}
				}
			}
		}
		// as.is, per column
		asIs := make([]bool, cols)
		switch x := a[8].(type) {
		case nil:
			for j := range asIs {
				asIs[j] = !stringsAsFactors
			}
		case *Logical:
			if len(x.V) == 0 {
				s.errorf(c.call, "invalid '%s' argument", "as.is")
			}
			for j := range asIs {
				asIs[j] = x.V[j%len(x.V)] == 1
			}
		case *Integer, *Double:
			for j := range asIs {
				asIs[j] = false
			}
			for _, d := range asDoubles(x) {
				if int(d) >= 1 && int(d) <= cols {
					asIs[int(d)-1] = true
				}
			}
		case *String:
			for j := range asIs {
				for _, e := range x.V {
					if e.S == colNames[j] {
						asIs[j] = true
					}
				}
			}
		default:
			s.errorf(c.call, "invalid '%s' argument", "as.is")
		}
		scanCols := make([]*scanCol, cols)
		for j := range scanCols {
			switch classes[j] {
			case "NULL":
				scanCols[j] = newScanCol(skNull)
			case "logical":
				scanCols[j] = newScanCol(skLogical)
			case "integer":
				scanCols[j] = newScanCol(skInteger)
			case "numeric", "double":
				scanCols[j] = newScanCol(skDouble)
			case "complex":
				scanCols[j] = newScanCol(skComplex)
			default:
				scanCols[j] = newScanCol(skString)
			}
		}
		maxRec := 0
		if nrows > 0 {
			maxRec = nrows
		}
		nrec := s.scanRun(readTableScanLang, sc, scanCols, maxRec, 0)
		if !ours {
			cn.pos = sc.pos
		}

		// the columns
		vals := make([]Value, cols)
		for j, col := range scanCols {
			switch {
			case col.kind == skNull:
			case col.kind != skString:
				vals[j] = col.value()
			case classes[j] == "character":
				vals[j] = col.value()
			case classes[j] == "factor":
				vals[j] = s.callNamed(c, "factor", []Value{col.value()}, []string{""})
			case classes[j] == "Date":
				vals[j] = s.callNamed(c, "as.Date", []Value{col.value()}, []string{""})
			case classes[j] == "POSIXct":
				vals[j] = s.callNamed(c, "as.POSIXct", []Value{col.value()}, []string{""})
			case classes[j] != "":
				s.errorf(c.call, "no method for coercing this S4 class to a vector")
			default:
				o := &tcOpts{asIs: asIs[j], dec: sp.dec, numerals: "allow.loss", tryLogical: tryLogical}
				vals[j] = s.typeConvertStrings(c, col.s, o)
			}
		}
		// row names
		var rowNames Value
		drop := make([]bool, cols)
		if rlabp && a[6] == nil {
			rowNames = scanCols[0].value()
			drop[0] = true
		}
		switch x := a[6].(type) {
		case nil:
		case nullT:
		default:
			n := Length(x)
			isStr := false
			if _, ok := x.(*String); ok {
				isStr = true
			}
			k := -1
			if n == 1 {
				if isStr {
					for j, nm := range colNames {
						if nm == x.(*String).V[0].S {
							k = j
						}
					}
					if k < 0 && nrec != 1 {
						s.errorf(c.call, "invalid 'row.names' specification")
					}
				} else if d := asDoublesOrNil(x); len(d) == 1 && d[0] >= 1 && int(d[0]) <= cols {
					k = int(d[0]) - 1
				} else if nrec != 1 {
					s.errorf(c.call, "invalid 'row.names' specification")
				}
			}
			if k >= 0 {
				if vals[k] == nil {
					s.errorf(c.call, "invalid 'row.names' specification")
				}
				rowNames = vals[k]
				drop[k] = true
			} else {
				if n != nrec {
					s.errorf(c.call, "invalid 'row.names' length")
				}
				rowNames = x
			}
		}
		var outV []Value
		var outN []Str
		for j := range vals {
			if vals[j] == nil || drop[j] {
				continue
			}
			outV = append(outV, vals[j])
			outN = append(outN, Str{S: colNames[j]})
		}
		if len(outV) > s.opts.MaxVector {
			s.errorf(c.call, "cannot allocate vector of length %d", len(outV))
		}
		var rn Value
		if rowNames != nil {
			switch r := rowNames.(type) {
			case *String:
				seen := make(map[string]bool, len(r.V))
				hasNA := false
				for _, e := range r.V {
					key := e.S
					if e.NA {
						hasNA, key = true, "\x00NA"
					}
					if seen[key] {
						s.errorf(c.call, "duplicate 'row.names' are not allowed")
					}
					seen[key] = true
				}
				if hasNA {
					s.errorf(c.call, "missing values in 'row.names' are not allowed")
				}
				rn = &String{V: r.V}
			case *Integer:
				seen := make(map[int32]bool, len(r.V))
				for _, e := range r.V {
					if e == naInt {
						s.errorf(c.call, "missing values in 'row.names' are not allowed")
					}
					if seen[e] {
						s.errorf(c.call, "duplicate 'row.names' are not allowed")
					}
					seen[e] = true
				}
				rn = &Integer{V: r.V}
			default:
				cv := s.callNamed(c, "as.character", []Value{rowNames}, []string{""})
				st := cv.(*String)
				seen := make(map[string]bool, len(st.V))
				for _, e := range st.V {
					if e.NA || seen[e.S] {
						s.errorf(c.call, "duplicate 'row.names' are not allowed")
					}
					seen[e.S] = true
				}
				rn = st
			}
		} else if nrec == 0 {
			rn = &Integer{V: []int32{}}
		} else {
			rn = &Integer{V: []int32{naInt, int32(-nrec)}}
		}
		out := &List{V: outV}
		if out.V == nil {
			out.V = []Value{}
		}
		out.Attr = []Attr{
			{Name: "names", Value: &String{V: outN}},
			{Name: "class", Value: mkString("data.frame")},
			{Name: "row.names", Value: rn},
		}
		return out
	}
}
