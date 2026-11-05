package eval

import (
	"regexp"
	"strconv"
	"strings"
)

var posixClasses = map[string]bool{
	"alpha": true, "digit": true, "alnum": true, "upper": true, "lower": true,
	"space": true, "blank": true, "punct": true, "print": true, "graph": true,
	"cntrl": true, "xdigit": true,
}

var boundRE = regexp.MustCompile(`^([0-9]*)(,([0-9]*))?$`)

// translateTRE rewrites a TRE pattern as a Go regexp pattern. notbol makes ^
// match nothing (TRE's REG_NOTBOL). lazy reports non-greedy quantifiers,
// which TRE has and POSIX leftmost-longest does not. reason is TRE's own
// error text, or "backref" for a back-reference, which is not supported.
func translateTRE(p string, notbol bool) (out string, lazy bool, reason string) {
	var b strings.Builder
	rs := []rune(p)
	depth := 0
	atStart := true // nothing to repeat: start, after ( or |
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		switch r {
		case '\\':
			atStart = false
			if i+1 >= len(rs) {
				return "", false, "Trailing backslash"
			}
			i++
			n := rs[i]
			switch {
			case n >= '1' && n <= '9':
				return "", false, "backref"
			case n == '<' || n == '>':
				b.WriteString(`\b`)
			case n == 'e':
				b.WriteString(`\x1b`)
			case strings.ContainsRune("dDwWsSbBnrtfav", n):
				b.WriteRune('\\')
				b.WriteRune(n)
			case n < 128 && (n >= 'a' && n <= 'z' || n >= 'A' && n <= 'Z' || n == '0'):
				b.WriteRune(n)
			default:
				b.WriteString(regexp.QuoteMeta(string(n)))
			}
		case '[':
			atStart = false
			j := i + 1
			var cls strings.Builder
			cls.WriteByte('[')
			if j < len(rs) && rs[j] == '^' {
				cls.WriteByte('^')
				j++
			}
			first := true
			closed := false
			for ; j < len(rs); j++ {
				c := rs[j]
				if c == ']' && !first {
					closed = true
					break
				}
				first = false
				switch {
				case c == '[' && j+1 < len(rs) && rs[j+1] == ':':
					k := j + 2
					for k+1 < len(rs) && !(rs[k] == ':' && rs[k+1] == ']') {
						k++
					}
					if k+1 >= len(rs) {
						return "", false, "Missing ']'"
					}
					name := string(rs[j+2 : k])
					if !posixClasses[name] {
						return "", false, "Unknown character class name"
					}
					cls.WriteString("[:" + name + ":]")
					j = k + 1
				case c == '\\' || c == '[' || c == ']':
					cls.WriteByte('\\')
					cls.WriteRune(c)
				default:
					cls.WriteRune(c)
				}
			}
			if !closed {
				return "", false, "Missing ']'"
			}
			cls.WriteByte(']')
			b.WriteString(cls.String())
			i = j
		case '(':
			if i+1 < len(rs) && rs[i+1] == '?' {
				// (?:  and inline flags (?i) (?U); anything else (lookaround) is not TRE
				if i+2 < len(rs) && rs[i+2] == ':' {
					depth++
					atStart = true
					b.WriteString("(?:")
					i += 2
					continue
				}
				k := i + 2
				for k < len(rs) && strings.ContainsRune("imsxnrU-", rs[k]) {
					k++
				}
				if k < len(rs) && rs[k] == ')' && k > i+2 {
					var fl strings.Builder
					for _, f := range rs[i+2 : k] {
						if f == 'i' || f == 'U' || f == '-' {
							fl.WriteRune(f)
						}
					}
					if fl.Len() > 0 {
						b.WriteString("(?" + fl.String() + ")")
					}
					i = k
					continue
				}
				return "", false, "Invalid regexp"
			}
			depth++
			atStart = true
			b.WriteRune(r)
		case ')':
			if depth == 0 {
				atStart = false
				b.WriteString(`\)`)
				continue
			}
			depth--
			atStart = false
			b.WriteRune(r)
		case '|':
			atStart = true
			b.WriteRune(r)
		case '^':
			if notbol {
				b.WriteString(`[^\x00-\x{10FFFF}]`)
			} else {
				b.WriteRune('^')
			}
			// a repetition after ^ repeats the empty anchor
		case '*', '+', '?':
			if atStart {
				// nothing to repeat: TRE ignores the operator
				continue
			}
			b.WriteRune(r)
			if i+1 < len(rs) && rs[i+1] == '?' {
				lazy = true
			}
		case '{':
			end := -1
			for k := i + 1; k < len(rs); k++ {
				if rs[k] == '}' {
					end = k
					break
				}
			}
			if end < 0 {
				return "", false, "Missing '}'"
			}
			body := string(rs[i+1 : end])
			m := boundRE.FindStringSubmatch(body)
			if m == nil || body == "" || body == "," {
				return "", false, "Invalid contents of {}"
			}
			lo, hi := m[1], m[3]
			if lo != "" && hi != "" {
				l, _ := strconv.Atoi(lo)
				h, _ := strconv.Atoi(hi)
				if l > h {
					return "", false, "Invalid contents of {}"
				}
			}
			if atStart {
				i = end
				continue
			}
			if lo == "" {
				lo = "0"
			}
			if m[2] == "" {
				b.WriteString("{" + lo + "}")
			} else {
				b.WriteString("{" + lo + "," + hi + "}")
			}
			i = end
			if i+1 < len(rs) && rs[i+1] == '?' {
				lazy = true
			}
		default:
			atStart = false
			b.WriteRune(r)
		}
	}
	return b.String(), lazy, ""
}
