package eval

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-embedded-r/r/syntax"
)

// cut.Date and cut.POSIXt (slice 16): Go natives, written from the help page
// of cut.Date and the oracle's answers. Their first R draft was too close to
// R's own file for the clean-room gate (72 exact and 81 normalised tokens
// against dates.R), so, as the rule says, it became native code.

var cutKinds = []string{"secs", "mins", "hours", "days", "weeks", "months", "years", "DSTdays", "quarters"}

const maxCutBreaks = 16777216

// levelsCall is the call the reference shows when the labels are fewer than
// the intervals.
const levelsCallText = "`levels<-.factor`(`*tmp*`, value = as.character(if (is.numeric(breaks)) x[!duplicated(res)] else breaks[-length(breaks)]))"

// parseCutSpec splits "2 months" into the unit it names and the count.
func (s *Session) parseCutSpec(c *callInfo, text string, allowed []string) (string, int) {
	bad := func() { s.errorf(c.call, "invalid specification of 'breaks'") }
	words := strings.Split(text, " ")
	for len(words) > 0 && words[len(words)-1] == "" {
		words = words[:len(words)-1]
	}
	if len(words) < 1 || len(words) > 2 {
		bad()
	}
	last := words[len(words)-1]
	kind, hits := "", 0
	for _, k := range allowed {
		if k == last {
			kind, hits = k, 1
			break
		}
		if last != "" && strings.HasPrefix(k, last) {
			kind, hits = k, hits+1
		}
	}
	if hits != 1 {
		bad()
	}
	mult := 1
	if len(words) == 2 {
		n, err := strconv.Atoi(words[0])
		if err != nil {
			bad()
		}
		mult = n
	}
	return kind, mult
}

func (s *Session) cutTimes(c *callInfo, isDate bool) Value {
	m, extra := s.splitArgs(c, []string{"x", "breaks", "labels", "start.on.monday", "right"})
	if m[0] == nil {
		s.errorf(c.call, "argument \"x\" is missing, with no default")
	}
	if m[1] == nil {
		s.errorf(c.call, "argument \"breaks\" is missing, with no default")
	}
	class, conv := "POSIXt", "as.POSIXct"
	if isDate {
		class, conv = "Date", "as.Date"
	}
	if !inheritsClass(m[0], class) {
		s.errorf(c.call, "'x' must be a date-time object")
	}
	call := func(name string, args ...Value) Value {
		return s.callNamed(c, name, args, make([]string, len(args)))
	}
	named := func(name string, args []Value, names []string) Value { return s.callNamed(c, name, args, names) }
	x := call(conv, m[0])
	monday := m[3] == nil || s.truthy(c, m[3], "start.on.monday")
	right := m[4] != nil && s.truthy(c, m[4], "right")
	var labels Value = Null
	if m[2] != nil {
		labels = m[2]
	}
	breaks := m[1]
	count := false
	switch b := breaks.(type) {
	case *String:
		if len(b.V) != 1 || b.V[0].NA {
			s.errorf(c.call, "invalid specification of 'breaks'")
		}
		if isDate {
			breaks = s.dateBreaks(c, x, b.V[0].S, monday)
		} else {
			breaks = s.timeBreaks(c, x, b.V[0].S, monday)
		}
	case *Double, *Integer:
		if inheritsClass(breaks, class) {
			breaks = call("sort", call(conv, breaks))
		} else if Length(b) == 1 && !isObject(breaks) {
			count = true
		} else {
			s.errorf(c.call, "invalid specification of 'breaks'")
		}
	default:
		s.errorf(c.call, "invalid specification of 'breaks'")
	}
	xnum := &Double{V: asDoubles(x)}
	var bnum Value = breaks
	if !count {
		bnum = &Double{V: asDoubles(breaks)}
	}
	args := []Value{xnum, bnum, labels, mkLogical(right)}
	names := []string{"", "", "labels", "right"}
	for _, e := range extra {
		args = append(args, e.Value)
		names = append(names, e.Name)
	}
	res := named("cut", args, names)
	if _, isNull := labels.(nullT); !isNull {
		return res
	}
	var value Value
	if count {
		value = call("as.character", call("[", x, call("!", call("duplicated", res))))
	} else {
		value = call("as.character", call("[", breaks, mkInt(int32(-Length(breaks)))))
	}
	if lev, ok := getAttr(res, "levels").(*String); ok && Length(value) < len(lev.V) {
		node, err := syntax.ParseAll(levelsCallText)
		if err == nil && len(node.Exprs) == 1 {
			s.errorf(s.fromSyntax(node.Exprs[0]).(*Lang), "number of levels differs")
		}
		s.errorf(c.call, "number of levels differs")
	}
	return named("levels<-", []Value{res, value}, []string{"", "value"})
}

func (s *Session) dateBreaks(c *callInfo, x Value, text string, monday bool) Value {
	kind, mult := s.parseCutSpec(c, text, []string{"days", "weeks", "months", "years", "DSTdays", "quarters"})
	lo, hi := math.Inf(1), math.Inf(-1)
	for _, v := range asDoubles(x) {
		if !math.IsNaN(v) {
			lo, hi = math.Min(lo, v), math.Max(hi, v)
		}
	}
	if math.IsInf(lo, 0) || mult < 1 {
		s.errorf(c.call, "invalid specification of 'breaks'")
	}
	day := func(v float64) time.Time { return time.Unix(int64(math.Floor(v))*86400, 0).UTC() }
	days := func(t time.Time) float64 { return float64(t.Unix() / 86400) }
	byDays := kind == "days" || kind == "DSTdays" || kind == "weeks"
	step := mult
	first := lo
	t := day(lo)
	switch kind {
	case "weeks":
		back := int(t.Weekday())
		if monday {
			back--
		}
		back = ((back % 7) + 7) % 7
		first = lo - float64(back)
		step = 7 * mult
	case "months", "years", "quarters":
		mon := int(t.Month()) - 1
		if kind == "years" {
			mon = 0
		} else if kind == "quarters" {
			mon -= mon % 3
		}
		first = days(time.Date(t.Year(), time.Month(mon+1), 1, 0, 0, 0, 0, time.UTC))
		step = mult * map[string]int{"months": 1, "years": 12, "quarters": 3}[kind]
	}
	width := float64(step)
	if !byDays {
		width *= 28
	}
	if (hi-first)/width+2 > maxCutBreaks {
		s.errorf(c.call, "'breaks' asks for too many intervals")
	}
	out := []float64{}
	for i := 0; ; i++ {
		var v float64
		if byDays {
			v = first + float64(step*i)
		} else {
			f := day(first)
			v = days(time.Date(f.Year(), f.Month()+time.Month(step*i), 1, 0, 0, 0, 0, time.UTC))
		}
		out = append(out, v)
		if v > hi {
			break
		}
	}
	return setAttr(&Double{V: out}, "class", mkString("Date"))
}

func (s *Session) timeBreaks(c *callInfo, x Value, text string, monday bool) Value {
	kind, mult := s.parseCutSpec(c, text, cutKinds)
	if mult < 1 {
		s.errorf(c.call, "invalid specification of 'breaks'")
	}
	call := func(name string, args ...Value) Value {
		return s.callNamed(c, name, args, make([]string, len(args)))
	}
	rm := func(name string) Value {
		return s.callNamed(c, name, []Value{x, mkLogical(true)}, []string{"", "na.rm"})
	}
	early, late := rm("min"), rm("max")
	unit := map[string]string{"secs": "secs", "mins": "mins", "hours": "hours", "months": "months", "quarters": "months", "years": "years"}[kind]
	if unit == "" {
		unit = "days"
	}
	early = call("as.POSIXct", call("trunc.POSIXt", early, mkString(unit)))
	field := func(name string) int {
		lt, ok := call("as.POSIXlt", early).(*List)
		if !ok {
			return 0
		}
		v := asDoubles(listField(lt, name))
		if len(v) == 0 || math.IsNaN(v[0]) {
			return 0
		}
		return int(v[0])
	}
	switch kind {
	case "weeks":
		back := field("wday")
		if monday {
			back--
		}
		back = ((back % 7) + 7) % 7
		early = call("as.POSIXct", call("trunc.POSIXt", call("-", early, mkDouble(float64(back)*86400)), mkString("days")))
	case "quarters":
		if k := field("mon") % 3; k > 0 {
			shifted := s.callNamed(c, "seq.POSIXt", []Value{early, mkString(fmt.Sprintf("-%d months", k)), mkInt(2)}, []string{"", "by", "length.out"})
			early = call("[", shifted, mkInt(2))
		}
	}
	by := fmt.Sprintf("%d %s", mult, kind)
	if kind == "weeks" {
		by = fmt.Sprintf("%d days", mult*7)
	}
	width := float64(mult) * map[string]float64{"secs": 1, "mins": 60, "hours": 3600, "weeks": 604800, "months": 28 * 86400, "quarters": 84 * 86400, "years": 365 * 86400}[kind]
	if width == 0 { // days and DSTdays
		width = float64(mult) * 86400
	}
	lo, hi := asDoubles(early), asDoubles(late)
	if len(lo) != 1 || len(hi) != 1 || math.IsNaN(lo[0]) || math.IsNaN(hi[0]) {
		s.errorf(c.call, "invalid specification of 'breaks'")
	}
	n := math.Ceil((hi[0]-lo[0])/width) + 2
	if n > maxCutBreaks || n < 1 {
		s.errorf(c.call, "'breaks' asks for too many intervals")
	}
	pts := s.callNamed(c, "seq.POSIXt", []Value{early, mkString(by), mkInt(int32(n))}, []string{"", "by", "length.out"})
	vals := asDoubles(pts)
	keep := len(vals)
	for i, v := range vals {
		if v > hi[0] {
			keep = i + 1
			break
		}
	}
	idx := make([]int32, keep)
	for i := range idx {
		idx[i] = int32(i + 1)
	}
	return call("[", pts, &Integer{V: idx})
}

func init() {
	regBuiltin("cut.Date", func(s *Session, c *callInfo) Value { return s.cutTimes(c, true) })
	regBuiltin("cut.POSIXt", func(s *Session, c *callInfo) Value { return s.cutTimes(c, false) })
}
