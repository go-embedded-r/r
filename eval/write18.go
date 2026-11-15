package eval

import (
	"strconv"
	"strings"
)

// write.table, write.csv and write.csv2 (slice 18), written from the help page
// of write.table and from probes of the oracle.

const writeTableOpenCall = `file(file, ifelse(append, "a", "w"))`

var writeTableFormals = []string{"x", "file", "append", "quote", "sep", "eol", "na", "dec", "row.names", "col.names", "qmethod", "fileEncoding"}

func inheritsFrom(v Value, class string) bool {
	for _, c := range classAttr(v) {
		if c == class {
			return true
		}
	}
	return false
}

// rowNameStrings are the row names of a data frame as text.
func rowNameStrings(df *List, nrow int) []string {
	out := make([]string, nrow)
	rn := getAttr(df, "row.names")
	switch r := rn.(type) {
	case *String:
		for i := range out {
			if i < len(r.V) {
				out[i] = r.V[i].S
			}
		}
	case *Integer:
		if len(r.V) == 2 && r.V[0] == naInt {
			for i := range out {
				out[i] = strconv.Itoa(i + 1)
			}
		} else {
			for i := range out {
				if i < len(r.V) {
					out[i] = strconv.Itoa(int(r.V[i]))
				}
			}
		}
	default:
		for i := range out {
			out[i] = strconv.Itoa(i + 1)
		}
	}
	return out
}

func dfRows(df *List) int {
	rn := getAttr(df, "row.names")
	switch r := rn.(type) {
	case *Integer:
		if len(r.V) == 2 && r.V[0] == naInt {
			n := int(r.V[1])
			if n < 0 {
				n = -n
			}
			return n
		}
		return len(r.V)
	case *String:
		return len(r.V)
	}
	if len(df.V) > 0 {
		return Length(df.V[0])
	}
	return 0
}

func (s *Session) writeTable(c *callInfo, call *Lang, a []Value, csvNA bool) Value {
	cc := *c
	cc.call = call
	c = &cc
	x := a[0]
	if x == nil {
		s.errorf(c.call, "argument \"x\" is missing, with no default")
	}
	df, ok := x.(*List)
	if !ok || !inheritsFrom(x, "data.frame") {
		v := s.callNamed(c, "as.data.frame", []Value{x}, []string{""})
		df, ok = v.(*List)
		if !ok {
			s.errorf(c.call, "invalid argument")
		}
		switch x.(type) {
		case *Logical, *Integer, *Double, *String, *Complex:
			if getAttr(x, "dim") == nil && len(df.V) == 1 && !inheritsFrom(x, "table") {
				// a vector is a data frame with one column, x
				df = &List{V: df.V, Attr: append([]Attr(nil), df.Attr...)}
				for i, at := range df.Attr {
					if at.Name == "names" {
						df.Attr[i].Value = mkString("x")
					}
				}
			}
		}
	}
	nrow := dfRows(df)
	ncol := len(df.V)
	sep := s.optString(c, a[4], " ", "sep")
	eol := s.optString(c, a[5], "\n", "eol")
	na := s.optString(c, a[6], "NA", "na")
	dec := s.optString(c, a[7], ".", "dec")
	if len(dec) != 1 {
		s.errorf(c.call, "'dec' must be a single character")
	}
	qmethod := "escape"
	if a[10] != nil {
		q, ok := a[10].(*String)
		if !ok || len(q.V) == 0 {
			s.errorf(c.call, "invalid '%s' argument", "qmethod")
		}
		qmethod = q.V[0].S
		if qmethod != "escape" && qmethod != "double" {
			s.errorf(&Lang{Fun: Intern("match.arg"), Args: []Arg{{Value: Intern("qmethod")}}},
				"'arg' should be one of “escape”, “double”")
		}
	}
	appendMode := s.boolArg(c, a[2], false, "append")

	// which columns are quoted
	quoteAll, quoteHeader := true, true
	var quoteCols map[int]bool
	switch q := a[3].(type) {
	case nil:
	case *Logical:
		if len(q.V) != 1 || q.V[0] == naInt {
			s.errorf(c.call, "invalid '%s' argument", "quote")
		}
		quoteAll, quoteHeader = q.V[0] == 1, q.V[0] == 1
	case *Integer, *Double:
		quoteAll = false
		quoteCols = map[int]bool{}
		for _, d := range asDoubles(q) {
			if d != d || d < 1 || int(d) > ncol {
				s.errorf(c.call, "invalid numbers in 'quote'")
			}
			quoteCols[int(d)-1] = true
		}
	default:
		s.errorf(c.call, "invalid '%s' argument", "quote")
	}
	enc := func(t string) string {
		if qmethod == "double" {
			return `"` + strings.ReplaceAll(t, `"`, `""`) + `"`
		}
		return `"` + strings.ReplaceAll(t, `"`, `\"`) + `"`
	}
	quoteHdr := func(t string) string {
		if quoteHeader {
			return enc(t)
		}
		return t
	}

	// the row and column labels
	var rowLabels []string
	switch r := a[8].(type) {
	case nil:
		rowLabels = rowNameStrings(df, nrow)
	case *Logical:
		if len(r.V) != 1 || r.V[0] == naInt {
			s.errorf(c.call, "invalid 'row.names' specification")
		}
		if r.V[0] == 1 {
			rowLabels = rowNameStrings(df, nrow)
		}
	case *String:
		if len(r.V) != nrow {
			s.errorf(c.call, "invalid 'row.names' specification")
		}
		for _, e := range r.V {
			rowLabels = append(rowLabels, e.S)
		}
	default:
		s.errorf(c.call, "invalid 'row.names' specification")
	}
	haveRows := rowLabels != nil
	nm := namesOf(df)
	colLabels := make([]string, ncol)
	for j := range colLabels {
		if nm != nil && j < len(nm.V) {
			colLabels[j] = nm.V[j].S
		}
	}
	header, blankCorner := true, false
	switch cn := a[9].(type) {
	case nil:
		blankCorner = csvNA && haveRows
	case *Logical:
		if len(cn.V) != 1 {
			s.errorf(c.call, "invalid 'col.names' specification")
		}
		switch cn.V[0] {
		case naInt:
			if !haveRows {
				s.errorf(c.call, "'col.names = NA' makes no sense when 'row.names = FALSE'")
			}
			blankCorner = true
		case 0:
			header = false
		}
	case *String:
		if len(cn.V) != ncol {
			s.errorf(c.call, "invalid 'col.names' specification")
		}
		for j, e := range cn.V {
			colLabels[j] = e.S
		}
	default:
		s.errorf(c.call, "invalid 'col.names' specification")
	}

	// the columns as text
	cols := make([][]string, ncol)
	quoted := make([]bool, ncol)
	for j, col := range df.V {
		if Length(col) != nrow {
			s.errorf(c.call, "invalid 'x' argument")
		}
		cols[j] = make([]string, nrow)
		isFactor := inheritsFrom(col, "factor")
		st, isStr := col.(*String)
		isStr = isStr && !isFactor
		var txt []Str
		switch {
		case isStr:
			txt = st.V
		default:
			v := s.callNamed(c, "as.character", []Value{col}, []string{""})
			r, ok := v.(*String)
			if !ok || len(r.V) != nrow {
				s.errorf(c.call, "invalid 'x' argument")
			}
			txt = r.V
		}
		quoted[j] = (isStr || isFactor) && (quoteAll || quoteCols[j])
		dcol, isDbl := col.(*Double)
		for i, e := range txt {
			switch {
			case e.NA || (isDbl && dcol.V[i] != dcol.V[i]):
				cols[j][i] = na
			case quoted[j]:
				cols[j][i] = enc(e.S)
			case isDbl && !isFactor && dec != ".":
				cols[j][i] = strings.Replace(e.S, ".", dec, 1)
			default:
				cols[j][i] = e.S
			}
		}
	}

	var sb strings.Builder
	limit := s.opts.MaxFileBytes
	if header {
		first := true
		put := func(t string) {
			if !first {
				sb.WriteString(sep)
			}
			first = false
			sb.WriteString(t)
		}
		if blankCorner {
			put(quoteHdr(""))
		}
		for _, l := range colLabels {
			put(quoteHdr(l))
		}
		sb.WriteString(eol)
	}
	for i := 0; i < nrow; i++ {
		first := true
		if haveRows {
			sb.WriteString(quoteHdr(rowLabels[i]))
			first = false
		}
		for j := range cols {
			if !first {
				sb.WriteString(sep)
			}
			first = false
			sb.WriteString(cols[j][i])
		}
		sb.WriteString(eol)
		if sb.Len() > limit {
			s.errorf(c.call, "cannot write more than %d bytes", limit)
		}
	}
	var target Value = a[1]
	s.writeTo(c, target, "file", appendMode, sb.String(), callOf(writeTableOpenCall))
	return s.invisible(Null)
}

func init() {
	regBuiltin("write.table", func(s *Session, c *callInfo) Value {
		return s.writeTable(c, c.call, s.matchFormals(c, writeTableFormals...), false)
	})
	writeCSV := func(name, sep, dec string) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			ignored := map[string]bool{"append": true, "col.names": true, "sep": true, "dec": true, "qmethod": true}
			a := make([]Value, len(writeTableFormals))
			used := map[string]bool{}
			pos := 0
			for i, v := range c.args {
				n := c.names[i]
				if n == "" {
					for pos < len(writeTableFormals) && (used[writeTableFormals[pos]] || ignored[writeTableFormals[pos]]) {
						pos++
					}
					if pos >= len(writeTableFormals) {
						s.errorf(c.call, "unused argument (%s)", s.deparseValue(v)[0])
					}
					n = writeTableFormals[pos]
				}
				if ignored[n] {
					continue
				}
				k := -1
				for j, f := range writeTableFormals {
					if f == n {
						k = j
					}
				}
				if k < 0 {
					s.errorf(c.call, "unused argument (%s = %s)", n, s.deparseValue(v)[0])
				}
				used[n] = true
				a[k] = v
			}
			for _, n := range []string{"append", "col.names", "sep", "dec", "qmethod"} {
				for i := range c.args {
					if c.names[i] == n {
						s.warningf(c.call, "attempt to set '%s' ignored", n)
					}
				}
			}
			a[4], a[7], a[10] = mkString(sep), mkString(dec), mkString("double")
			a[2] = mkLogical(false)
			a[9] = nil
			return s.writeTable(c, c.call, a, true)
		})
	}
	writeCSV("write.csv", ",", ".")
	writeCSV("write.csv2", ";", ",")
}
