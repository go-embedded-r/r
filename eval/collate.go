package eval

import (
	"strings"
	"unicode"
)

// String collation. R 4.6.2 run with LC_ALL=C.UTF-8 (the oracle's locale)
// still collates with ICU, not with code points: punctuation before digits
// before letters, letters compared without regard to accent or case first,
// then accents, then case with lowercase first. This file implements that
// ordering for Latin text, enough for sort(), order(), ls() and the
// comparison operators on strings; characters outside Latin script compare by
// code point after the Latin letters. A session whose Options.CCollation is
// set compares by code point instead (R's method = "radix" and LC_COLLATE=C).

const asciiPunctOrder = " _-,;:!?.'\"()[]{}@*/\\&#%`^+<=>|~$"

var accentBase = map[rune]rune{}

func init() {
	for base, set := range map[rune]string{
		'a': "àáâãäåāăą", 'c': "çćĉċč", 'd': "ďđ", 'e': "èéêëēĕėęě", 'g': "ĝğġģ", 'h': "ĥħ",
		'i': "ìíîïĩīĭįı", 'j': "ĵ", 'k': "ķ", 'l': "ĺļľŀł", 'n': "ñńņňŉ", 'o': "òóôõöøōŏő",
		'r': "ŕŗř", 's': "śŝşš", 't': "ţťŧ", 'u': "ùúûüũūŭůűų", 'w': "ŵ", 'y': "ýÿŷ", 'z': "źżž",
	} {
		for _, r := range set {
			accentBase[r] = base
			accentBase[unicode.ToUpper(r)] = unicode.ToUpper(base)
		}
	}
}

// collKey is the three-level key of one character.
type collKey struct {
	primary   int
	secondary rune
	upper     bool
}

func keyOf(r rune) collKey {
	if i := strings.IndexRune(asciiPunctOrder, r); i >= 0 {
		return collKey{primary: i}
	}
	base, accented := r, rune(0)
	if b, ok := accentBase[r]; ok {
		base, accented = b, r
	}
	lower := unicode.ToLower(base)
	up := lower != base
	switch {
	case unicode.IsDigit(base):
		return collKey{primary: 1000 + int(base), secondary: accented}
	case unicode.IsLetter(base) && base < 0x250:
		return collKey{primary: 2000 + int(lower), secondary: accented, upper: up}
	case unicode.IsLetter(base):
		return collKey{primary: 100000 + int(lower), secondary: accented, upper: up}
	}
	// other symbols and punctuation: after the ASCII ones, before digits
	return collKey{primary: 100 + int(base)%800}
}

// collate compares two strings: -1, 0 or 1.
func collate(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	n := min(len(ra), len(rb))
	ka, kb := make([]collKey, n), make([]collKey, n)
	for i := 0; i < n; i++ {
		ka[i], kb[i] = keyOf(ra[i]), keyOf(rb[i])
		if ka[i].primary != kb[i].primary {
			if ka[i].primary < kb[i].primary {
				return -1
			}
			return 1
		}
	}
	if len(ra) != len(rb) {
		if len(ra) < len(rb) {
			return -1
		}
		return 1
	}
	for i := 0; i < n; i++ {
		if ka[i].secondary != kb[i].secondary {
			if ka[i].secondary < kb[i].secondary {
				return -1
			}
			return 1
		}
	}
	for i := 0; i < n; i++ {
		if ka[i].upper != kb[i].upper {
			if !ka[i].upper {
				return -1
			}
			return 1
		}
	}
	return strings.Compare(a, b)
}
