package eval

import (
	"sort"
	"strings"

	"github.com/go-embedded-r/tre"
)

// Regular expressions: R's default engine is TRE, a POSIX extended regular
// expression matcher. This file drives github.com/go-embedded-r/tre, a Go port
// of TRE built with its R-compatibility options, over R's own flags
// (REG_EXTENDED, plus REG_ICASE for ignore.case): lazy quantifiers, capture
// priority in repeated groups, back-references and \< \> behave as in R. The
// engine R does not use by default (perl = TRUE, PCRE) is refused with an
// error that says so. See the README, "Regular expressions".

type matcher struct {
	fixed  bool
	pat    string
	re     *tre.Regexp
	groups int
	s      *Session
	c      *callInfo

	// the subject converted to code units, kept across the calls that scan
	// one string repeatedly (gsub, gregexpr)
	lastStr string
	hasLast bool
	units   []rune
	offs    []int // byte offset of each character, plus the length; nil for ASCII
}

// compileRegex builds the matcher of a pattern, refusing what is not supported.
//
// R matches bytes (tre_regcompb: character classes are expanded over the 256
// byte values) when the pattern and every text are ASCII, and wide characters
// otherwise; the two modes differ in corners (a compilation error omits the
// pattern in wide mode, negated classes are lost in copies of a bounded
// repetition), so the choice is reproduced.
func (s *Session) compileRegex(c *callInfo, pattern string, texts []Str, fixed, icase, perl bool) *matcher {
	return s.compileRegexMsg(c, pattern, texts, fixed, icase, perl, false)
}

// compileRegexGrep is compileRegex for grep, grepl, regexpr and gregexpr, whose
// wide-character path reports a compilation error without the pattern.
func (s *Session) compileRegexGrep(c *callInfo, pattern string, texts []Str, fixed, icase, perl bool) *matcher {
	return s.compileRegexMsg(c, pattern, texts, fixed, icase, perl, true)
}

func (s *Session) compileRegexMsg(c *callInfo, pattern string, texts []Str, fixed, icase, perl, omitPatWide bool) *matcher {
	if perl {
		s.errorf(c.call, "perl = TRUE is not supported yet (R's PCRE engine); use the default extended regular expressions")
	}
	if fixed {
		if icase {
			s.warningf(c.call, "argument 'ignore.case = TRUE' will be ignored")
		}
		if pattern == "" {
			s.errorf(c.call, "zero-length pattern")
		}
		return &matcher{fixed: true, pat: pattern}
	}
	flags := tre.Extended
	if icase {
		flags |= tre.ICase
	}
	wide := !allASCII(pattern, texts)
	if !wide {
		flags |= tre.UseBytes
	}
	re, err := tre.CompileOptions(pattern, tre.Options{Flags: flags, RCompat: true,
		Limits: tre.Limits{MaxSteps: 1 << 26}})
	if err != nil {
		reason := err.Error()
		s.warningf(c.call, "TRE pattern compilation error '%s'", reason)
		if wide && omitPatWide {
			s.errorf(c.call, "invalid regular expression, reason '%s'", reason)
		}
		s.errorf(c.call, "invalid regular expression '%s', reason '%s'", pattern, reason)
	}
	return &matcher{pat: pattern, re: re, groups: re.NumSubexp(), s: s, c: c}
}

// convert fills units (and offs) for str.
func (m *matcher) convert(str string) {
	if m.hasLast && m.lastStr == str {
		return
	}
	m.lastStr, m.hasLast = str, true
	ascii := true
	for i := 0; i < len(str); i++ {
		if str[i] >= 0x80 {
			ascii = false
			break
		}
	}
	m.units = m.units[:0]
	if ascii {
		m.offs = nil
		for i := 0; i < len(str); i++ {
			m.units = append(m.units, rune(str[i]))
		}
		return
	}
	m.offs = m.offs[:0]
	for i, r := range str {
		m.units = append(m.units, r)
		m.offs = append(m.offs, i)
	}
	m.offs = append(m.offs, len(str))
}

// find returns the submatch byte offsets (absolute) of the first match at or
// after from, or nil. With from > 0 the text before is not seen, as with TRE's
// REG_NOTBOL on the rest of the string.
func (m *matcher) find(str string, from int) []int { return m.findBuf(nil, str, from) }

// findBuf is find writing the offsets into dst when it is large enough.
func (m *matcher) findBuf(dst []int, str string, from int) []int {
	if m.fixed {
		i := strings.Index(str[from:], m.pat)
		if i < 0 {
			return nil
		}
		return append(dst[:0], from+i, from+i+len(m.pat))
	}
	m.convert(str)
	start := from
	if m.offs != nil {
		start = sort.SearchInts(m.offs, from)
	}
	var ef tre.EFlags
	if from > 0 {
		ef = tre.NotBOL
	}
	pm, err := m.re.Exec(m.units[start:], m.groups+1, ef)
	if err != nil {
		if err == tre.ErrNoMatch {
			return nil
		}
		if err == tre.ErrSteps {
			m.s.errorf(m.c.call, "regular expression '%s' exceeds the match step limit on this text", m.pat)
		}
		m.s.errorf(m.c.call, "regular expression '%s' failed on the text: %s", m.pat, err.Error())
	}
	loc := dst[:0]
	if cap(loc) < 2*len(pm) {
		loc = make([]int, 0, 2*len(pm))
	}
	loc = loc[:2*len(pm)]
	for i, p := range pm {
		so, eo := p.So, p.Eo
		if so >= 0 {
			if m.offs != nil {
				so, eo = m.offs[so+start], m.offs[eo+start]
			} else {
				so, eo = so+start, eo+start
			}
		}
		loc[2*i], loc[2*i+1] = so, eo
	}
	return loc
}

// matches reports whether the pattern occurs in str (grepl): no submatches are
// computed, so the matcher stops at the first match.
func (m *matcher) matches(str string) bool {
	if m.fixed {
		return strings.Contains(str, m.pat)
	}
	m.convert(str)
	_, err := m.re.Exec(m.units, 0, 0)
	if err != nil {
		if err == tre.ErrNoMatch {
			return false
		}
		m.s.errorf(m.c.call, "regular expression '%s' exceeds the match step limit on this text", m.pat)
	}
	return true
}
