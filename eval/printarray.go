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
func (s *Session) printMatrix(v Value, base, nr, nc int, rowNames, colNames *String, quote bool) {
	rw := len(strconv.Itoa(nr)) + 3
	if rowNames != nil {
		rw = 0
		for _, e := range rowNames.V {
			if n := utf8.RuneCountInString(labelOrNA(e)); n > rw {
				rw = n
			}
		}
	}
	rowLabel := func(i int) string {
		if rowNames != nil {
			return padRight(labelOrNA(rowNames.V[i]), rw)
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
			left = true
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
		s.out.WriteString(strings.Repeat(" ", rw))
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
	if len(dims) == 2 {
		s.printMatrix(v, 0, nr, nc, dimnamesOf(v, 0), dimnamesOf(v, 1), quote)
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
			if nm := dimnamesOf(v, i+2); nm != nil {
				label += labelOrNA(nm.V[c])
			} else {
				label += strconv.Itoa(c + 1)
			}
		}
		s.out.WriteString(label + "\n\n")
		s.printMatrix(v, k*slice, nr, nc, dimnamesOf(v, 0), dimnamesOf(v, 1), quote)
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
		s.out.WriteString("factor(0)\n")
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
		total += utf8.RuneCountInString(l.S) + 1
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
			parts = append(parts, l.S)
		}
	}
	prefix := ""
	if drop {
		prefix = strconv.Itoa(n) + " "
	}
	s.out.WriteString(prefix + "Levels: " + strings.Join(parts, " ") + "\n")
}
