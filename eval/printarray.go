package eval

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

func labelOrNA(s Str) string {
	if s.NA {
		return "<NA>"
	}
	return s.S
}

// printMatrix prints a two-dimensional slice of v (elements at offset base,
// nr rows, nc columns) as R does.
//
// dnn is names(dimnames(v)), or nil: when it is set, the corner and a line over
// the columns carry the names of the two dimensions, and the row labels are
// offset (two columns, or as many as the row name needs).
func (s *Session) printMatrix(v Value, base, nr, nc int, rowNames, colNames *String, quote bool, dnn []string) {
	if nr > maxPrintRows || nc > maxPrintRows {
		s.errorf(nil, "matrix too large to print: more than %d rows or columns", maxPrintRows)
	}
	if nr == 0 && nc == 0 {
		s.out.WriteString("<0 x 0 matrix>\n")
		return
	}
	rw := len(strconv.Itoa(nr)) + 3
	if rowNames != nil {
		rw = 0
		for _, e := range rowNames.V {
			if n := utf8.RuneCountInString(labelOrNA(e)); n > rw {
				rw = n
			}
		}
	}
	lbloff := 0
	if dnn != nil {
		rnw := utf8.RuneCountInString(dnn[0])
		if rnw < rw+2 {
			lbloff = 2
		} else {
			lbloff = rnw - rw
		}
	}
	labw := rw // the width of the labels themselves
	rw += lbloff
	rowLabel := func(i int) string {
		if rowNames != nil {
			return strings.Repeat(" ", lbloff) + padRight(labelOrNA(rowNames.V[i]), labw)
		}
		return padLeft("["+strconv.Itoa(i+1)+",]", rw)
	}
	type column struct {
		label string
		items []string
		w     int
		left  bool
	}
	cols := make([]column, nc)
	for j := 0; j < nc; j++ {
		pos := make([]int, nr)
		for i := range pos {
			pos[i] = base + j*nr + i
		}
		col := takeElems(v, pos)
		items, left := formatElements(col, quote)
		if _, isStr := v.(*String); isStr {
			left = !s.printRight
		}
		label := "[," + strconv.Itoa(j+1) + "]"
		if colNames != nil {
			label = labelOrNA(colNames.V[j])
		}
		w := maxWidth(items)
		if n := utf8.RuneCountInString(label); n > w {
			w = n
		}
		cols[j] = column{label: label, items: items, w: w, left: left}
	}
	if nc == 0 {
		s.out.WriteString(strings.Repeat(" ", rw) + "\n")
		for i := 0; i < nr; i++ {
			s.out.WriteString(rowLabel(i) + "\n")
		}
		return
	}
	pad := func(c column, text string) string {
		if c.left {
			return " " + padRight(text, c.w)
		}
		return " " + padLeft(text, c.w)
	}
	for start := 0; start < nc; {
		end, width := start, rw
		for end < nc && (end == start || width+cols[end].w+1 <= printWidth) {
			width += cols[end].w + 1
			end++
		}
		if dnn != nil {
			s.out.WriteString(strings.Repeat(" ", rw) + dnn[1] + "\n")
			s.out.WriteString(padRight(dnn[0], rw))
		} else {
			s.out.WriteString(strings.Repeat(" ", rw))
		}
		for j := start; j < end; j++ {
			s.out.WriteString(pad(cols[j], cols[j].label))
		}
		s.out.WriteString("\n")
		for i := 0; i < nr; i++ {
			s.out.WriteString(rowLabel(i))
			for j := start; j < end; j++ {
				s.out.WriteString(pad(cols[j], cols[j].items[i]))
			}
			s.out.WriteString("\n")
		}
		start = end
	}
}

// maxPrintSlices bounds the matrix slices printArray emits; beyond it the
// print is refused (R would run for ever or exhaust memory).
const maxPrintSlices = 1 << 20

// maxPrintRows bounds the rows of a matrix printed with no column to show them.
const maxPrintRows = 1 << 22

// listCellText is what a list matrix shows for one cell: the value of a plain
// length-one vector, else the type and length ("numeric,3"); NULL for NULL.
func listCellText(e Value) string {
	switch v := e.(type) {
	case nullT:
		return "NULL"
	case *Logical, *Integer, *Double, *String:
		if Length(v) == 1 && !isObject(v) {
			items, _ := formatElements(v, true)
			return strings.TrimSpace(items[0])
		}
		name := typeName(v)
		if name == "double" {
			name = "numeric"
		}
		return name + "," + strconv.Itoa(Length(v))
	case *List:
		return "list," + strconv.Itoa(len(v.V))
	case *Closure, *Builtin:
		return "?"
	}
	return "?"
}

// dimnamesNames is names(dimnames(v)) as strings, or nil when there are none.
func dimnamesNames(v Value) []string {
	dn, ok := getAttr(v, "dimnames").(*List)
	if !ok {
		return nil
	}
	nm := namesOf(dn)
	if nm == nil {
		return nil
	}
	out := make([]string, len(nm.V))
	for i, e := range nm.V {
		if !e.NA {
			out[i] = e.S
		}
	}
	return out
}

// oneDView is a one-dimensional array with dimnames as the named vector R
// prints it as: the dimension's labels become names; dnn is the name of the
// dimension when the dimnames are named.
func oneDView(v Value) (view Value, dnn string, named bool) {
	lab := dimnamesOf(v, 0)
	if lab == nil {
		return v, "", false
	}
	if n := dimnamesNames(v); n != nil {
		dnn, named = n[0], true
	}
	var attr []Attr
	for _, a := range attrsOf(v) {
		if a.Name != "dim" && a.Name != "dimnames" && a.Name != "names" {
			attr = append(attr, a)
		}
	}
	attr = append(attr, Attr{Name: "names", Value: lab})
	return withAttr(v, attr), dnn, named
}

func dimnamesOf(v Value, d int) *String {
	if dn, ok := getAttr(v, "dimnames").(*List); ok && d < len(dn.V) {
		if nm, ok := dn.V[d].(*String); ok {
			return nm
		}
	}
	return nil
}

// printArray prints an array of two or more dimensions: a matrix, or one
// matrix for each slice of the higher dimensions.
func (s *Session) printArray(v Value, dims []int, quote bool) {
	nr, nc := dims[0], dims[1]
	dnn := dimnamesNames(v)
	var mdnn []string
	if len(dnn) >= 2 {
		mdnn = dnn[:2]
	}
	if len(dims) == 2 {
		s.printMatrix(v, 0, nr, nc, dimnamesOf(v, 0), dimnamesOf(v, 1), quote, mdnn)
		return
	}
	slice := nr * nc
	nslices := 1
	for _, d := range dims[2:] {
		nslices *= d
		if nslices > maxPrintSlices || nslices < 0 {
			// an array with a zero extent and astronomically many slices
			s.errorf(nil, "array too large to print: more than %d slices", maxPrintSlices)
		}
	}
	counter := make([]int, len(dims)-2)
	for k := 0; k < nslices; k++ {
		s.checkInterrupt()
		label := ", , "
		for i, c := range counter {
			if i > 0 {
				label += ", "
			}
			one := strconv.Itoa(c + 1)
			if nm := dimnamesOf(v, i+2); nm != nil {
				one = labelOrNA(nm.V[c])
			}
			if i+2 < len(dnn) && dnn[i+2] != "" {
				one = dnn[i+2] + " = " + one
			}
			label += one
		}
		s.out.WriteString(label + "\n\n")
		s.printMatrix(v, k*slice, nr, nc, dimnamesOf(v, 0), dimnamesOf(v, 1), quote, mdnn)
		s.out.WriteString("\n")
		for i := range counter {
			counter[i]++
			if counter[i] < dims[i+2] {
				break
			}
			counter[i] = 0
		}
	}
}

func (s *Session) printFactor(v Value) {
	lev, _ := getAttr(v, "levels").(*String)
	codes := asInts(v)
	if lev == nil {
		lev = &String{}
	}
	labels := make([]Str, len(codes))
	for i, c := range codes {
		if c == naInt || int(c) < 1 || int(c) > len(lev.V) {
			labels[i] = Str{NA: true}
		} else {
			labels[i] = lev.V[c-1]
		}
	}
	if len(codes) == 0 {
		s.out.WriteString("factor()\n")
	} else {
		var attr []Attr
		if nm := namesOf(v); nm != nil {
			attr = namesAttr(nm.V)
		}
		s.printVector(&String{V: labels, Attr: attr}, false)
	}
	// the levels line, truncated as print.factor does
	n := len(lev.V)
	width := printWidth - (len("Levels: ") + 3 + 1 + 3)
	cum := make([]int, n)
	total := 0
	for i, l := range lev.V {
		total += utf8.RuneCountInString(labelOrNA(l)) + 1
		cum[i] = total
	}
	maxl := n
	if n > 1 && cum[n-1] > width {
		first := 0
		for i, c := range cum {
			if c > width {
				first = i + 1
				break
			}
		}
		maxl = first - 1
		if maxl < 1 {
			maxl = 1
		}
	}
	drop := n > maxl
	var parts []string
	if drop {
		k := maxl - 1
		if k < 1 {
			k = 1
		}
		for _, l := range lev.V[:k] {
			parts = append(parts, l.S)
		}
		parts = append(parts, "...")
		if maxl > 1 {
			parts = append(parts, lev.V[n-1].S)
		}
	} else {
		for _, l := range lev.V {
			parts = append(parts, labelOrNA(l))
		}
	}
	prefix := ""
	if drop {
		prefix = strconv.Itoa(n) + " "
	}
	sep := " "
	if hasString(classAttr(v), "ordered") {
		sep = " < "
	}
	s.out.WriteString(prefix + "Levels: " + strings.Join(parts, sep) + "\n")
}
