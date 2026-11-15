package eval

import (
	"strconv"
	"strings"
)

// The marks format, formatC and prettyNum can put in a number (slice 16):
// big.mark and big.interval, small.mark and small.interval, decimal.mark,
// zero.print and drop0trailing. Written from the help page of prettyNum and
// from the oracle's answers.

type numMarks struct {
	big           string
	bigInterval   int
	small         string
	smallInterval int
	decimal       string
	zero          *string
	drop0         bool
}

// active reports whether any option changes a number at all.
func (m numMarks) active() bool {
	return m.big != "" || m.small != "" || (m.decimal != "" && m.decimal != ".") || m.zero != nil || m.drop0
}

// finish writes the options that act on a number after it has been padded to
// its field: a zero printed as the text zero.print gives (as wide as the zero
// was), and the zeros that end a number dropped, with a zero exponent.
func (m numMarks) finish(item string) string {
	if m.zero == nil && !m.drop0 {
		return item
	}
	trimmed := strings.TrimSpace(item)
	if trimmed == "" {
		return item
	}
	if m.zero != nil {
		if f, err := strconv.ParseFloat(trimmed, 64); err == nil && f == 0 {
			z := *m.zero
			for len(z) < len(trimmed) {
				z += " "
			}
			return strings.Replace(item, trimmed, z, 1)
		}
	}
	if m.drop0 && strings.Contains(trimmed, ".") {
		lead := item[:strings.Index(item, trimmed)]
		if e := strings.IndexAny(trimmed, "eE"); e >= 0 {
			if trimmed[e+1:] == "+00" || trimmed[e+1:] == "-00" {
				return lead + trimmed[:e]
			}
			return item
		}
		return lead + strings.TrimSuffix(strings.TrimRight(trimmed, "0"), ".")
	}
	return item
}

// readMarks collects the options from matched arguments; a nil Value is an
// argument that was not given.
func (s *Session) readMarks(c *callInfo, big, bigInterval, small, smallInterval, decimal, zero, drop0 Value) numMarks {
	m := numMarks{bigInterval: 3, smallInterval: 5}
	if big != nil {
		m.big = s.stringArg(c, big, "big.mark")
	}
	if bigInterval != nil {
		m.bigInterval = s.intArg(c, bigInterval, "big.interval")
		if m.bigInterval < 1 {
			s.errorf(c.call, "invalid 'big.interval' argument")
		}
	}
	if small != nil {
		m.small = s.stringArg(c, small, "small.mark")
	}
	if smallInterval != nil {
		m.smallInterval = s.intArg(c, smallInterval, "small.interval")
		if m.smallInterval < 1 {
			s.errorf(c.call, "invalid 'small.interval' argument")
		}
	}
	if decimal != nil {
		m.decimal = s.stringArg(c, decimal, "decimal.mark")
	}
	if zero != nil {
		if _, isNull := zero.(nullT); !isNull {
			if st, ok := zero.(*String); ok && len(st.V) == 1 {
				z := st.V[0].S
				m.zero = &z
			} else if l, ok := zero.(*Logical); ok && len(l.V) == 1 && l.V[0] == 1 {
				z := "0"
				m.zero = &z
			}
		}
	}
	if drop0 != nil {
		m.drop0 = s.truthy(c, drop0, "drop0trailing")
	}
	return m
}

// apply writes one formatted number with the marks. Anything that is not a
// plain number (NA, a name, text) passes through.
func (m numMarks) apply(item string) string {
	if (m.big == "" && m.small == "" && (m.decimal == "" || m.decimal == ".")) || item == "" {
		return item
	}
	start := 0
	for start < len(item) && (item[start] == ' ' || item[start] == '-' || item[start] == '+') {
		start++
	}
	end := start
	for end < len(item) && item[end] >= '0' && item[end] <= '9' {
		end++
	}
	if end == start {
		return item
	}
	intPart := item[start:end]
	rest := item[end:]
	frac, tail := "", rest
	if strings.HasPrefix(rest, ".") {
		k := 1
		for k < len(rest) && rest[k] >= '0' && rest[k] <= '9' {
			k++
		}
		frac, tail = rest[1:k], rest[k:]
	}
	var b strings.Builder
	b.WriteString(item[:start])
	if m.big != "" {
		for i, ch := range intPart {
			if i > 0 && (len(intPart)-i)%m.bigInterval == 0 {
				b.WriteString(m.big)
			}
			b.WriteRune(ch)
		}
	} else {
		b.WriteString(intPart)
	}
	if frac != "" || rest != tail {
		dec := "."
		if m.decimal != "" {
			dec = m.decimal
		}
		b.WriteString(dec)
		if m.small != "" {
			for i, ch := range frac {
				if i > 0 && i%m.smallInterval == 0 {
					b.WriteString(m.small)
				}
				b.WriteRune(ch)
			}
		} else {
			b.WriteString(frac)
		}
	}
	b.WriteString(tail)
	return b.String()
}
