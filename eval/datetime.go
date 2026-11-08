package eval

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata" // the time zone database travels with the evaluator
)

// Dates and times: the helpers behind baselib/datetime.R. Written from the
// help pages of strptime, as.POSIXlt and as.Date, and from what the oracle
// answers; the zone rules come from Go's time package (the IANA database).
// The helpers are hidden: they live in the base environment under dotted
// names and are not part of the functions the evaluator counts.

// hidden are native helpers installed in base but not listed by Builtins.
var hidden = map[string]*Builtin{}

func regHidden(name string, fn func(s *Session, c *callInfo) Value) {
	hidden[name] = &Builtin{Name: name, fn: fn}
}

var (
	monthNames = []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	dayNames   = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
)

// zoneNamed resolves a time zone name as R's tz argument gives it: "" is the
// session's zone (the TZ environment variable, else the system's), "UTC" and
// "GMT" are themselves. An unknown name is UTC, as R treats it.
func zoneNamed(name string) *time.Location {
	if name == "" {
		if tz := os.Getenv("TZ"); tz != "" {
			name = tz
		} else {
			return time.Local
		}
	}
	if loc, err := time.LoadLocation(name); err == nil {
		return loc
	}
	return time.UTC
}

// zoneAbbrevs are the standard and daylight abbreviations of a zone, as the
// tzone attribute of a POSIXlt lists them: found in the current year.
func zoneAbbrevs(loc *time.Location) (std, dst string) {
	// walk back through the zone's transitions from now: the first standard
	// and the first daylight abbreviation met
	t := time.Now().In(loc)
	for i := 0; i < 400 && t.Year() > 1900; i++ {
		name, _ := t.Zone()
		if t.IsDST() {
			if dst == "" {
				dst = name
			}
		} else if std == "" {
			std = name
		}
		if std != "" && dst != "" {
			break
		}
		start, _ := t.ZoneBounds()
		if start.IsZero() {
			break
		}
		t = start.Add(-time.Second)
	}
	if std == "" {
		std, _ = time.Now().In(loc).Zone()
	}
	if dst == "" {
		dst = std
	}
	return std, dst
}

func naInts(n int) []int32 {
	out := make([]int32, n)
	for i := range out {
		out[i] = naInt
	}
	return out
}

// ltFields are the names of the components of a POSIXlt, in order.
var ltFields = []string{"sec", "min", "hour", "mday", "mon", "year", "wday", "yday", "isdst", "zone", "gmtoff"}

// ltFromSeconds breaks seconds since the epoch into the components of a
// POSIXlt in the given zone.
func ltFromSeconds(secs []float64, loc *time.Location, utc bool) *List {
	n := len(secs)
	sec := make([]float64, n)
	ints := make([][]int32, 8) // min hour mday mon year wday yday isdst
	for i := range ints {
		ints[i] = make([]int32, n)
	}
	zone := make([]Str, n)
	off := make([]int32, n)
	for i, x := range secs {
		if math.IsNaN(x) || math.IsInf(x, 0) || math.Abs(x) > 1e15 {
			if math.IsInf(x, 0) || math.IsNaN(x) {
				sec[i] = x // Inf, NaN and NA stay what they were
			} else {
				sec[i] = naReal
			}
			for k := range ints {
				ints[k][i] = naInt
			}
			ints[7][i] = -1
			if utc {
				zn, _ := time.Now().In(loc).Zone()
				zone[i], off[i] = Str{S: zn}, 0
			} else {
				zone[i], off[i] = Str{S: ""}, naInt
			}
			continue
		}
		fl := math.Floor(x)
		t := time.Unix(int64(fl), 0).In(loc)
		sec[i] = float64(t.Second()) + (x - fl)
		ints[0][i] = int32(t.Minute())
		ints[1][i] = int32(t.Hour())
		ints[2][i] = int32(t.Day())
		ints[3][i] = int32(t.Month()) - 1
		ints[4][i] = int32(t.Year() - 1900)
		ints[5][i] = int32(t.Weekday())
		ints[6][i] = int32(t.YearDay() - 1)
		if t.IsDST() {
			ints[7][i] = 1
		}
		name, o := t.Zone()
		zone[i] = Str{S: name}
		off[i] = int32(o)
	}
	l := &List{V: []Value{&Double{V: sec}}}
	for _, v := range ints {
		l.V = append(l.V, &Integer{V: v})
	}
	l.V = append(l.V, &String{V: zone}, &Integer{V: off})
	l.Attr = namesAttr(strsOf(ltFields))
	return l
}

func strsOf(ss []string) []Str {
	out := make([]Str, len(ss))
	for i, s := range ss {
		out[i] = Str{S: s}
	}
	return out
}

// ltColumns reads the numeric columns of a POSIXlt list; ok is false when it
// is not one. n is the length of the longest column.
func ltColumns(v Value) (cols [9][]float64, zone []Str, off []int32, n int, ok bool) {
	l, isList := v.(*List)
	if !isList || len(l.V) < 9 {
		return cols, nil, nil, 0, false
	}
	for i := 0; i < 9; i++ {
		cols[i] = asDoubles(l.V[i])
		if len(cols[i]) > n {
			n = len(cols[i])
		}
	}
	for i := 0; i < 9; i++ {
		if len(cols[i]) == 0 {
			n = 0
		}
	}
	if len(l.V) > 9 {
		if z, isStr := l.V[9].(*String); isStr {
			zone = z.V
		}
	}
	if len(l.V) > 10 {
		off = asInts(l.V[10])
	}
	return cols, zone, off, n, true
}

// ltTime is element i of the components as a time in loc, normalised the way
// mktime does (a day 32 is the first of the next month); ok is false for NA.
func ltTime(cols [9][]float64, i int, loc *time.Location) (t time.Time, frac float64, ok bool) {
	var f [9]float64
	for k := 0; k < 6; k++ { // sec min hour mday mon year
		c := cols[k]
		f[k] = c[i%len(c)]
		if math.IsNaN(f[k]) || math.IsInf(f[k], 0) {
			return t, 0, false
		}
	}
	sec := math.Floor(f[0])
	frac = f[0] - sec
	t = time.Date(int(f[5])+1900, time.Month(int(f[4])+1), int(f[3]), int(f[2]), int(f[1]), int(sec), 0, loc)
	return t, frac, true
}

// tzoneAttr is the tzone attribute of a POSIXlt made in a zone: the name
// alone for UTC, else the name and the two abbreviations.
func tzoneAttr(name string, loc *time.Location, third bool) Value {
	if name == "UTC" || name == "GMT" && !third {
		return mkString(name)
	}
	if !third {
		return mkString(name)
	}
	std, dst := zoneAbbrevs(loc)
	return strVec([]string{name, std, dst})
}

// isUTCZone says whether the zone a tz argument names is UTC: "UTC" or "GMT",
// or "" when the session's own zone is one of them.
func isUTCZone(tz string) bool {
	if tz == "" {
		tz = os.Getenv("TZ")
		if tz == "" && time.Local.String() == "UTC" {
			tz = "UTC"
		}
	}
	return tz == "UTC" || tz == "GMT"
}

func init() {
	regHidden(".dtIsUTC", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "tz")
		tz := ""
		if st, ok := a[0].(*String); ok && len(st.V) > 0 && !st.V[0].NA {
			tz = st.V[0].S
		}
		return mkLogical(isUTCZone(tz))
	})
	regHidden(".dtNow", func(s *Session, c *callInfo) Value {
		t := time.Now()
		return mkDouble(float64(t.UnixNano()) / 1e9)
	})
	regHidden(".dtToLt", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "tz")
		tz := ""
		if st, ok := a[1].(*String); ok && len(st.V) > 0 && !st.V[0].NA {
			tz = st.V[0].S
		}
		return ltFromSeconds(asDoubles(a[0]), zoneNamed(tz), isUTCZone(tz))
	})
	regHidden(".dtToCt", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "lt", "tz")
		tz := ""
		if st, ok := a[1].(*String); ok && len(st.V) > 0 && !st.V[0].NA {
			tz = st.V[0].S
		}
		cols, _, _, n, ok := ltColumns(a[0])
		if !ok {
			s.errorf(c.call, "invalid 'x' argument")
		}
		loc := zoneNamed(tz)
		out := make([]float64, n)
		for i := range out {
			t, frac, ok := ltTime(cols, i, loc)
			if !ok {
				if sv := cols[0][i%len(cols[0])]; math.IsInf(sv, 0) || math.IsNaN(sv) {
					out[i] = sv
				} else {
					out[i] = naReal
				}
				continue
			}
			out[i] = float64(firstOccurrence(t, loc).Unix()) + frac
		}
		return &Double{V: out}
	})
	regHidden(".dtTzone", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "tz")
		tz := ""
		if st, ok := a[0].(*String); ok && len(st.V) > 0 && !st.V[0].NA {
			tz = st.V[0].S
		}
		loc := zoneNamed(tz)
		return tzoneAttr(tz, loc, true)
	})
	regHidden(".dtFormat", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "lt", "format", "np")
		np := 0
		if a[2] != nil {
			if d := asDoubles(a[2]); len(d) == 1 && d[0] >= 0 && d[0] <= 6 {
				np = int(d[0])
			}
		}
		cols, zone, off, n, ok := ltColumns(a[0])
		if !ok {
			s.errorf(c.call, "invalid 'x' argument")
		}
		fm := asStrings(a[1])
		if len(fm) == 0 {
			s.errorf(c.call, "invalid 'format' argument")
		}
		if n == 0 {
			return &String{V: []Str{}}
		}
		if len(fm) > n {
			n = len(fm)
		}
		out := make([]Str, n)
		for i := range out {
			f := fm[i%len(fm)]
			if f.NA {
				out[i] = Str{NA: true}
				continue
			}
			t, frac, ok := ltTime(cols, i, time.UTC)
			if !ok {
				// the non-finite times print as the numbers they are
				switch sv := cols[0][i%len(cols[0])]; {
				case math.IsInf(sv, 1):
					out[i] = Str{S: "Inf"}
				case math.IsInf(sv, -1):
					out[i] = Str{S: "-Inf"}
				case math.IsNaN(sv) && !isNAReal(sv):
					out[i] = Str{S: "NaN"}
				default:
					out[i] = Str{NA: true}
				}
				continue
			}
			var zn string
			var zo int32 = naInt
			if len(zone) > 0 {
				zn = zone[i%len(zone)].S
			}
			if len(off) > 0 {
				zo = off[i%len(off)]
			}
			if !strings.Contains(f.S, "%z") && !strings.Contains(f.S, "%Z") {
				// R's strftime knows the zone's offset only when the format
				// asks for the zone: %s is otherwise the fields read as UTC
				zo = naInt
			}
			txt, err := strftime(t, frac, expandOS(f.S, float64(t.Second())+frac, np), zn, zo)
			if err != "" {
				s.errorf(c.call, "%s", err)
			}
			out[i] = Str{S: txt}
		}
		return &String{V: out}
	})
	regHidden(".dtParse", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "format", "tz")
		xs, fm := asStrings(a[0]), asStrings(a[1])
		tz := ""
		if st, ok := a[2].(*String); ok && len(st.V) > 0 && !st.V[0].NA {
			tz = st.V[0].S
		}
		n := len(xs)
		if len(fm) == 0 {
			s.errorf(c.call, "invalid 'format' argument")
		}
		if len(xs) == 0 {
			n = 0
		} else if len(fm) > n {
			n = len(fm)
		}
		loc := zoneNamed(tz)
		now := time.Now().In(loc)
		sec := make([]float64, n)
		ints := make([][]int32, 8)
		for k := range ints {
			ints[k] = make([]int32, n)
		}
		zone := make([]Str, n)
		off := make([]int32, n)
		for i := 0; i < n; i++ {
			x, f := xs[i%len(xs)], fm[i%len(fm)]
			var p parsedTime
			ok := !x.NA && !f.NA && parseTime(x.S, f.S, &p)
			if ok {
				ok = p.complete(now)
			}
			if !ok {
				if p.badYday {
					s.warningf(s.currentCall(), "day-of-year %d in year %d is invalid", p.yday, p.year)
				}
				sec[i] = naReal
				for k := range ints {
					ints[k][i] = naInt
				}
				ints[7][i] = -1
				if isUTCZone(tz) {
					zn, _ := time.Now().In(loc).Zone()
					zone[i], off[i] = Str{S: zn}, 0
				} else {
					zone[i], off[i] = Str{S: ""}, naInt
				}
				continue
			}
			sec[i] = float64(p.sec) + p.frac
			ints[0][i], ints[1][i], ints[2][i], ints[3][i], ints[4][i] = int32(p.min), int32(p.hour), int32(p.mday), int32(p.mon), int32(p.year-1900)
			d := time.Date(p.year, time.Month(p.mon+1), p.mday, 12, 0, 0, 0, time.UTC)
			ints[5][i], ints[6][i] = int32(d.Weekday()), int32(d.YearDay()-1)
			t := time.Date(p.year, time.Month(p.mon+1), p.mday, p.hour, p.min, p.sec, 0, loc)
			name, o := t.Zone()
			if p.hasOffset {
				// %z: the time is as given, shifted to the zone asked for
				t = time.Date(p.year, time.Month(p.mon+1), p.mday, p.hour, p.min, p.sec, 0, time.UTC).Add(-time.Duration(p.offset) * time.Second).In(loc)
				name, o = t.Zone()
				ints[0][i], ints[1][i], ints[2][i], ints[3][i], ints[4][i] = int32(t.Minute()), int32(t.Hour()), int32(t.Day()), int32(t.Month())-1, int32(t.Year()-1900)
				sec[i] = float64(t.Second()) + p.frac
				ints[5][i], ints[6][i] = int32(t.Weekday()), int32(t.YearDay()-1)
			}
			if t.IsDST() {
				ints[7][i] = 1
			}
			zone[i] = Str{S: name}
			if isUTCZone(tz) {
				off[i] = 0
			} else {
				off[i] = naInt
			}
			_ = o
		}
		l := &List{V: []Value{&Double{V: sec}}}
		for _, v := range ints {
			l.V = append(l.V, &Integer{V: v})
		}
		l.V = append(l.V, &String{V: zone}, &Integer{V: off})
		l.Attr = namesAttr(strsOf(ltFields))
		return l
	})
}

// ---- strftime ---------------------------------------------------------------

// isoWeek is the ISO 8601 week number and week-based year.
func isoWeek(t time.Time) (week, year int) { y, w := t.ISOWeek(); return w, y }

// strftime formats t (whose fields are those of the broken-down time, in
// UTC) with a C-locale format; zone and gmtoff are the zone's own.
func strftime(t time.Time, frac float64, f string, zone string, gmtoff int32) (string, string) {
	var b strings.Builder
	sec := float64(t.Second()) + frac
	for i := 0; i < len(f); i++ {
		ch := f[i]
		if ch != '%' {
			b.WriteByte(ch)
			continue
		}
		i++
		if i >= len(f) {
			b.WriteByte('%')
			break
		}
		c := f[i]
		if c == 'O' && i+1 < len(f) && f[i+1] == 'S' {
			i++
			digits := 0
			if i+1 < len(f) && f[i+1] >= '0' && f[i+1] <= '6' {
				digits = int(f[i+1] - '0')
				i++
			}
			if digits == 0 {
				fmt.Fprintf(&b, "%02d", int(math.Floor(sec)))
			} else {
				p := math.Pow(10, float64(digits))
				v := math.Floor(sec*p+1e-6) / p // truncated, not rounded
				fmt.Fprintf(&b, "%0*.*f", 3+digits, digits, v)
			}
			continue
		}
		if c == 'E' || c == 'O' { // alternative representations: the standard one
			if i+1 < len(f) {
				i++
				c = f[i]
			}
		}
		switch c {
		case '%':
			b.WriteByte('%')
		case 'a':
			b.WriteString(dayNames[t.Weekday()][:3])
		case 'A':
			b.WriteString(dayNames[t.Weekday()])
		case 'b', 'h':
			b.WriteString(monthNames[t.Month()-1][:3])
		case 'B':
			b.WriteString(monthNames[t.Month()-1])
		case 'c':
			s, _ := strftime(t, frac, "%a %b %e %H:%M:%S %Y", zone, gmtoff)
			b.WriteString(s)
		case 'C':
			fmt.Fprintf(&b, "%02d", floorDiv(t.Year(), 100))
		case 'd':
			fmt.Fprintf(&b, "%02d", t.Day())
		case 'D':
			s, _ := strftime(t, frac, "%m/%d/%y", zone, gmtoff)
			b.WriteString(s)
		case 'e':
			fmt.Fprintf(&b, "%2d", t.Day())
		case 'F':
			s, _ := strftime(t, frac, "%Y-%m-%d", zone, gmtoff)
			b.WriteString(s)
		case 'g':
			_, y := isoWeek(t)
			fmt.Fprintf(&b, "%02d", ((y%100)+100)%100)
		case 'G':
			_, y := isoWeek(t)
			b.WriteString(strconv.Itoa(y))
		case 'H':
			fmt.Fprintf(&b, "%02d", t.Hour())
		case 'I':
			h := t.Hour() % 12
			if h == 0 {
				h = 12
			}
			fmt.Fprintf(&b, "%02d", h)
		case 'j':
			fmt.Fprintf(&b, "%03d", t.YearDay())
		case 'k':
			fmt.Fprintf(&b, "%2d", t.Hour())
		case 'l':
			h := t.Hour() % 12
			if h == 0 {
				h = 12
			}
			fmt.Fprintf(&b, "%2d", h)
		case 'm':
			fmt.Fprintf(&b, "%02d", int(t.Month()))
		case 'M':
			fmt.Fprintf(&b, "%02d", t.Minute())
		case 'n':
			b.WriteByte('\n')
		case 'p':
			if t.Hour() < 12 {
				b.WriteString("AM")
			} else {
				b.WriteString("PM")
			}
		case 'P':
			if t.Hour() < 12 {
				b.WriteString("am")
			} else {
				b.WriteString("pm")
			}
		case 'r':
			s, _ := strftime(t, frac, "%I:%M:%S %p", zone, gmtoff)
			b.WriteString(s)
		case 'R':
			s, _ := strftime(t, frac, "%H:%M", zone, gmtoff)
			b.WriteString(s)
		case 's':
			// the fields read as UTC, less the zone's offset when it is known
			if gmtoff != naInt {
				fmt.Fprintf(&b, "%d", t.Unix()-int64(gmtoff))
			} else {
				fmt.Fprintf(&b, "%d", t.Unix())
			}
		case 'S':
			fmt.Fprintf(&b, "%02d", t.Second())
		case 't':
			b.WriteByte('\t')
		case 'T':
			s, _ := strftime(t, frac, "%H:%M:%S", zone, gmtoff)
			b.WriteString(s)
		case 'u':
			wd := int(t.Weekday())
			if wd == 0 {
				wd = 7
			}
			b.WriteString(strconv.Itoa(wd))
		case 'U':
			fmt.Fprintf(&b, "%02d", (t.YearDay()-1+7-int(t.Weekday()))/7)
		case 'V':
			w, _ := isoWeek(t)
			fmt.Fprintf(&b, "%02d", w)
		case 'w':
			b.WriteString(strconv.Itoa(int(t.Weekday())))
		case 'W':
			fmt.Fprintf(&b, "%02d", (t.YearDay()-1+7-(int(t.Weekday())+6)%7)/7)
		case 'x':
			s, _ := strftime(t, frac, "%m/%d/%y", zone, gmtoff)
			b.WriteString(s)
		case 'X':
			s, _ := strftime(t, frac, "%H:%M:%S", zone, gmtoff)
			b.WriteString(s)
		case 'y':
			fmt.Fprintf(&b, "%02d", ((t.Year()%100)+100)%100)
		case 'Y':
			if y := t.Year(); y >= 0 && y < 1000 {
				fmt.Fprintf(&b, "%04d", y)
			} else {
				b.WriteString(strconv.Itoa(y))
			}
		case 'z':
			if gmtoff == naInt {
				b.WriteString("+0000")
			} else {
				o := int(gmtoff)
				sign := '+'
				if o < 0 {
					sign, o = '-', -o
				}
				fmt.Fprintf(&b, "%c%02d%02d", sign, o/3600, (o%3600)/60)
			}
		case 'Z':
			b.WriteString(zone)
		case '+':
			s, _ := strftime(t, frac, "%a %b %e %H:%M:%S %Z %Y", zone, gmtoff)
			b.WriteString(s)
		default:
			b.WriteByte('%')
			b.WriteByte(c)
		}
	}
	if b.Len() > 2048 {
		return "", "output of format exceeds the limit of 2048 bytes"
	}
	return b.String(), ""
}

func floorDiv(a, b int) int {
	q := a / b
	if a%b != 0 && (a < 0) != (b < 0) {
		q--
	}
	return q
}

// ---- strptime ---------------------------------------------------------------

type parsedTime struct {
	year, mon, mday      int // mon is 0-based
	hour, min, sec       int
	frac                 float64
	yday                 int
	hasYear, hasMon      bool
	hasMday, hasYday     bool
	hasHour, hasMin      bool
	hasSec               bool
	pm, hasPM            bool
	hour12               bool
	offset               int // seconds east of UTC, from %z
	hasOffset            bool
	century              int
	hasCentury, hasYear2 bool
	year2                int
	badYday              bool // %j named a day the year does not have
}

// complete fills what the input left out the way R does (an unspecified year,
// month or day is the current one) and checks the date is a real one.
func (p *parsedTime) complete(now time.Time) bool {
	if p.hasYear2 {
		if p.hasCentury {
			p.year, p.hasYear = p.century*100+p.year2, true
		} else if p.year2 < 69 {
			p.year, p.hasYear = 2000+p.year2, true
		} else {
			p.year, p.hasYear = 1900+p.year2, true
		}
	} else if !p.hasYear && p.hasCentury {
		p.year, p.hasYear = p.century*100, true
	}
	if p.hasPM {
		if !p.hour12 {
			// %p without %I is ignored
		} else if p.pm && p.hour < 12 {
			p.hour += 12
		} else if !p.pm && p.hour == 12 {
			p.hour = 0
		}
	}
	if !p.hasYear {
		p.year = now.Year()
	}
	if p.hasYday && !p.hasMon && !p.hasMday {
		d := time.Date(p.year, 1, 1, 12, 0, 0, 0, time.UTC).AddDate(0, 0, p.yday-1)
		if d.Year() != p.year {
			p.badYday = true
			return false
		}
		p.mon, p.mday = int(d.Month())-1, d.Day()
		return p.normalizeHour24()
	}
	if !p.hasMon && !p.hasMday {
		p.mon, p.mday = int(now.Month())-1, now.Day()
	} else if !p.hasMon {
		p.mon = int(now.Month()) - 1
	} else if !p.hasMday {
		return false
	}
	// the day must exist in the month
	d := time.Date(p.year, time.Month(p.mon+1), p.mday, 12, 0, 0, 0, time.UTC)
	if !(d.Day() == p.mday && int(d.Month()) == p.mon+1 && p.mday >= 1) {
		return false
	}
	return p.normalizeHour24()
}

// normalizeHour24 turns 24:00:00, which ISO 8601 allows, into midnight of the
// next day; any other time in hour 24 is invalid.
func (p *parsedTime) normalizeHour24() bool {
	if p.hour != 24 {
		return true
	}
	if p.min != 0 || p.sec != 0 {
		return false
	}
	d := time.Date(p.year, time.Month(p.mon+1), p.mday, 12, 0, 0, 0, time.UTC).AddDate(0, 0, 1)
	p.year, p.mon, p.mday, p.hour = d.Year(), int(d.Month())-1, d.Day(), 0
	return true
}

func skipSpace(s string, i int) int {
	for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r' || s[i] == '\f' || s[i] == '\v') {
		i++
	}
	return i
}

// number reads up to maxDigits decimal digits (leading spaces allowed).
func number(s string, i, maxDigits, lo, hi int) (int, int, bool) {
	i = skipSpace(s, i)
	j := i
	v := 0
	for j < len(s) && j-i < maxDigits && s[j] >= '0' && s[j] <= '9' {
		v = v*10 + int(s[j]-'0')
		j++
	}
	if j == i || v < lo || v > hi {
		return 0, i, false
	}
	return v, j, true
}

// name matches one of the names, in full or abbreviated to three letters,
// case-insensitively; it returns the index.
func name(s string, i int, names []string) (int, int, bool) {
	rest := strings.ToLower(s[i:])
	for k, n := range names {
		ln := strings.ToLower(n)
		if strings.HasPrefix(rest, ln) {
			return k, i + len(ln), true
		}
	}
	for k, n := range names {
		if strings.HasPrefix(rest, strings.ToLower(n[:3])) {
			return k, i + 3, true
		}
	}
	return 0, i, false
}

// parseTime reads s as f says; trailing input is ignored. It reports whether
// every conversion matched.
func parseTime(s, f string, p *parsedTime) bool {
	i := 0
	return parseTimeFrom(s, &i, f, p)
}

// parseTimeFrom is parseTime from position *ip of s; it leaves *ip after what
// it read.
func parseTimeFrom(s string, ip *int, f string, p *parsedTime) bool {
	i := *ip
	defer func() { *ip = i }()
	var ok bool
	for k := 0; k < len(f); k++ {
		c := f[k]
		switch {
		case c == ' ' || c == '\t' || c == '\n':
			i = skipSpace(s, i)
			continue
		case c != '%':
			if i >= len(s) || s[i] != c {
				return false
			}
			i++
			continue
		}
		k++
		if k >= len(f) {
			return false
		}
		c = f[k]
		if (c == 'E' || c == 'O') && k+1 < len(f) && c == 'O' && f[k+1] == 'S' {
			k++
			c = 'Q' // %OS
		} else if c == 'E' || c == 'O' {
			k++
			if k >= len(f) {
				return false
			}
			c = f[k]
		}
		var v int
		switch c {
		case '%':
			if i >= len(s) || s[i] != '%' {
				return false
			}
			i++
		case 'a', 'A':
			v, i, ok = name(s, i, dayNames)
			if !ok {
				return false
			}
		case 'b', 'B', 'h':
			v, i, ok = name(s, i, monthNames)
			if !ok {
				return false
			}
			p.mon, p.hasMon = v, true
		case 'c':
			if !parseTime2(s, &i, "%a %b %e %H:%M:%S %Y", p) {
				return false
			}
		case 'C':
			if v, i, ok = number(s, i, 2, 0, 99); !ok {
				return false
			}
			p.century, p.hasCentury = v, true
		case 'd', 'e':
			if v, i, ok = number(s, i, 2, 1, 31); !ok {
				return false
			}
			p.mday, p.hasMday = v, true
		case 'D':
			// as R reads it on input: year first, unlike the output
			if !parseTime2(s, &i, "%y/%m/%d", p) {
				return false
			}
		case 'F':
			if !parseTime2(s, &i, "%Y-%m-%d", p) {
				return false
			}
		case 'H', 'k':
			if v, i, ok = number(s, i, 2, 0, 24); !ok {
				return false
			}
			p.hour, p.hasHour = v, true
		case 'I', 'l':
			if v, i, ok = number(s, i, 2, 1, 12); !ok {
				return false
			}
			p.hour, p.hasHour, p.hour12 = v, true, true
		case 'j':
			if v, i, ok = number(s, i, 3, 1, 366); !ok {
				return false
			}
			p.yday, p.hasYday = v, true
		case 'm':
			if v, i, ok = number(s, i, 2, 1, 12); !ok {
				return false
			}
			p.mon, p.hasMon = v-1, true
		case 'M':
			if v, i, ok = number(s, i, 2, 0, 59); !ok {
				return false
			}
			p.min, p.hasMin = v, true
		case 'n', 't':
			i = skipSpace(s, i)
		case 'p':
			rest := strings.ToLower(s[min(i, len(s)):])
			switch {
			case strings.HasPrefix(rest, "am"):
				p.pm, p.hasPM = false, true
				i += 2
			case strings.HasPrefix(rest, "pm"):
				p.pm, p.hasPM = true, true
				i += 2
			default:
				return false
			}
		case 'r':
			if !parseTime2(s, &i, "%I:%M:%S %p", p) {
				return false
			}
		case 'R':
			if !parseTime2(s, &i, "%H:%M", p) {
				return false
			}
		case 'S':
			if v, i, ok = number(s, i, 2, 0, 60); !ok {
				return false
			}
			p.sec, p.hasSec = v, true
		case 'Q': // %OS: seconds with a fractional part
			if v, i, ok = number(s, i, 2, 0, 60); !ok {
				return false
			}
			p.sec, p.hasSec = v, true
			if i < len(s) && s[i] == '.' {
				j := i + 1
				for j < len(s) && s[j] >= '0' && s[j] <= '9' {
					j++
				}
				if j > i+1 {
					fr, _ := strconv.ParseFloat("0"+s[i:j], 64)
					p.frac = fr
					i = j
				}
			}
		case 'T', 'X':
			if !parseTime2(s, &i, "%H:%M:%S", p) {
				return false
			}
		case 'u':
			if v, i, ok = number(s, i, 1, 1, 7); !ok {
				return false
			}
		case 'w':
			if v, i, ok = number(s, i, 1, 0, 6); !ok {
				return false
			}
		case 'U', 'W', 'V':
			if v, i, ok = number(s, i, 2, 0, 53); !ok {
				return false
			}
		case 'g':
			if v, i, ok = number(s, i, 2, 0, 99); !ok {
				return false
			}
		case 'G':
			if v, i, ok = number(s, i, 4, 0, 9999); !ok {
				return false
			}
		case 'x':
			if !parseTime2(s, &i, "%y/%m/%d", p) {
				return false
			}
		case 'y':
			if v, i, ok = number(s, i, 2, 0, 99); !ok {
				return false
			}
			p.year2, p.hasYear2 = v, true
		case 'Y':
			if v, i, ok = number(s, i, 4, 0, 9999); !ok {
				return false
			}
			p.year, p.hasYear = v, true
		case 'z':
			i = skipSpace(s, i)
			if i >= len(s) || (s[i] != '+' && s[i] != '-') {
				return false
			}
			sign := 1
			if s[i] == '-' {
				sign = -1
			}
			i++
			hh, j, ok1 := number(s, i, 2, 0, 14)
			if !ok1 || j-i != 2 {
				return false
			}
			i = j
			mm, j, ok2 := number(s, i, 2, 0, 59)
			if !ok2 || j-i != 2 {
				return false
			}
			i = j
			p.offset, p.hasOffset = sign*(hh*3600+mm*60), true
		default:
			return false
		}
	}
	return true
}

// parseTime2 parses a sub-format at *i of s.
func parseTime2(s string, i *int, f string, p *parsedTime) bool {
	return parseTimeFrom(s, i, f, p)
}

// expandOS settles the %OS conversions of a format for seconds sec: the first
// one is written out with the digits it names (or np when it names none), the
// following ones are plain seconds, as R handles only one per format.
func expandOS(f string, sec float64, np int) string {
	if !strings.Contains(f, "%OS") {
		return f
	}
	var b strings.Builder
	seen := false
	for i := 0; i < len(f); i++ {
		if f[i] != '%' || i+1 >= len(f) {
			b.WriteByte(f[i])
			continue
		}
		if f[i+1] == '%' {
			b.WriteString("%%")
			i++
			continue
		}
		if strings.HasPrefix(f[i:], "%OS") {
			i += 2
			if seen {
				b.WriteString("%S")
				continue
			}
			seen = true
			digits := np
			if i+1 < len(f) && f[i+1] >= '0' && f[i+1] <= '6' {
				digits = int(f[i+1] - '0')
				i++
			}
			if digits == 0 {
				fmt.Fprintf(&b, "%02d", int(math.Floor(sec)))
			} else {
				p := math.Pow(10, float64(digits))
				v := math.Floor(sec*p+1e-6) / p // truncated, not rounded
				fmt.Fprintf(&b, "%0*.*f", 3+digits, digits, v)
			}
			continue
		}
		b.WriteByte('%')
	}
	return b.String()
}

// firstOccurrence is t, except that a wall-clock time which occurs twice (the
// hour a daylight-saving change repeats) is read as its first occurrence, as
// R's mktime does; Go leaves that choice open.
func firstOccurrence(t time.Time, loc *time.Location) time.Time {
	y, mo, d := t.Date()
	h, mi, s := t.Clock()
	best := t
	for _, delta := range []int{-3 * 3600, 3 * 3600} {
		_, off := t.Add(time.Duration(delta) * time.Second).Zone()
		_, offHere := t.Zone()
		if off == offHere {
			continue
		}
		// the same wall clock read with the neighbouring offset
		naive := time.Date(y, mo, d, h, mi, s, 0, time.UTC).Unix()
		cand := time.Unix(naive-int64(off), 0).In(loc)
		cy, cmo, cd := cand.Date()
		ch, cmi, cs := cand.Clock()
		if cy == y && cmo == mo && cd == d && ch == h && cmi == mi && cs == s && cand.Before(best) {
			best = cand
		}
	}
	return best
}
