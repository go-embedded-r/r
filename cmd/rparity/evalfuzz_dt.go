package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// The counted exclusions of the datetime mode (see fuzzEnvExcluded for the
// others): the first line of a program is "# excl:<class>".
const (
	exclPre1902  = "pre1902"
	exclTzAbbrev = "tzabbrev"
	exclDstGap   = "dstgap"
	// the years after 2037, beyond the tabulated zone rules, where R and Go
	// extrapolate daylight-saving time their own ways
	exclPost2037 = "post2037"
)

// fuzzDatetime builds a random program over dates and times: formatting with
// every conversion, parsing, arithmetic, differences, sequences, time zone
// conversions, comparisons and the class-keeping methods. The first line is a
// comment that names the family.
func fuzzDatetime(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	zones := []string{"UTC", "GMT", "America/New_York", "Europe/Paris", "Asia/Tokyo", "Australia/Sydney", "Europe/London", "America/Los_Angeles"}
	if r.Intn(10) == 0 {
		return fuzzDatetimeExcluded(r)
	}
	// the years 1971 to 2037, where the zone rules are tabulated; later ones
	// are an exclusion of their own
	stamp := func(minYear, span int) string {
		s := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", minYear+r.Intn(span), 1+r.Intn(12), 1+r.Intn(28), (4+r.Intn(21))%24, r.Intn(60), r.Intn(60)) // never 01:00-03:59, where a zone's clock may skip or repeat (dstgap)
		if r.Intn(5) == 0 {
			s += fmt.Sprintf(".%03d", r.Intn(1000))
		}
		return s
	}
	day := func() string { return fmt.Sprintf("%04d-%02d-%02d", 1971+r.Intn(67), 1+r.Intn(12), 1+r.Intn(28)) }
	fmts := []string{
		"%Y-%m-%d %H:%M:%S", "%d/%m/%Y", "%A %d %B %Y", "%a %b %e", "%H:%M", "%I:%M %p", "%j", "%U %W %u %w",
		"%y%m%d", "%G-W%V-%u", "%D", "%F %T", "%R", "%c", "%x %X", "%Z %z", "%e|%k|%l", "%C %y", "%h %%", "%OS3", "%H:%M:%OS2", "%s", "%s %Z", "%z %s",
	}
	tz := func() string { return "tz = \"" + pick(zones) + "\"" }
	ct := func() string { return "as.POSIXct(\"" + stamp(1971, 67) + "\", " + tz() + ")" }
	cts := func(n int) string {
		parts := make([]string, n)
		for i := range parts {
			parts[i] = "\"" + stamp(1971, 67) + "\""
		}
		return "as.POSIXct(c(" + strings.Join(parts, ", ") + "), " + tz() + ")"
	}
	dt := func() string { return "as.Date(\"" + day() + "\")" }
	units := []string{"\"secs\"", "\"mins\"", "\"hours\"", "\"days\"", "\"weeks\"", "\"auto\""}

	var family, body string
	switch r.Intn(14) {
	case 0:
		family, body = "format-ct", "format("+ct()+", \""+pick(fmts)+"\""+pick([]string{"", ", usetz = TRUE", ", tz = \"" + pick(zones) + "\""})+")"
	case 1:
		family, body = "format-date", "format("+dt()+", \""+pick(fmts)+"\")"
	case 2:
		family, body = "lt-fields", "unclass(as.POSIXlt("+ct()+"))[c(\"sec\", \"min\", \"hour\", \"mday\", \"mon\", \"year\", \"wday\", \"yday\", \"isdst\", \"zone\", \"gmtoff\")]"
	case 3:
		// a round trip through text: what is formatted is parsed back
		f := pick([]string{"%Y-%m-%d %H:%M:%S", "%d/%m/%Y %H:%M", "%Y%m%d%H%M%S", "%Y-%j %H", "%b %d %Y %I:%M %p"})
		family, body = "roundtrip", "x <- "+ct()+"; s <- format(x, \""+f+"\"); list(s, format(strptime(s, \""+f+"\", "+tz()+"), \"%Y-%m-%d %H:%M:%S\"))"
	case 4:
		garbage := pick([]string{"2024-13-01", "2024-02-30", "24:61", "x", "", "2024-1-5 7:3", "  2024-03-10", "2024-03-10 trailing", "12/31/99", "99-12-31"})
		family, body = "parse", "strptime(\""+garbage+"\", \""+pick([]string{"%Y-%m-%d", "%m/%d/%y", "%y-%m-%d", "%H:%M", "%Y-%m-%d %H:%M", "%d.%m.%Y"})+"\", tz = \"UTC\")"
	case 5:
		family, body = "difftime", "difftime("+ct()+", "+ct()+", units = "+pick(units)+")"
	case 6:
		family, body = "diff-print", pick([]string{"print(", "format(", "units(", "as.numeric("})+ct()+" - "+ct()+")"
	case 7:
		family, body = "date-arith", pick([]string{dt() + " + " + fmt.Sprint(r.Intn(400)), dt() + " - " + fmt.Sprint(r.Intn(400)), dt() + " - " + dt(), "as.numeric(" + dt() + " - " + dt() + ")", dt() + " + as.difftime(" + fmt.Sprint(r.Intn(100)) + ", units = \"days\")"})
	case 8:
		family, body = "seq-date", "seq("+dt()+", by = \""+pick([]string{"day", "2 days", "week", "month", "2 months", "quarter", "year"})+"\", length.out = "+fmt.Sprint(1+r.Intn(6))+")"
	case 9:
		family, body = "convert", pick([]string{
			"as.numeric(" + ct() + ")", "as.Date(" + ct() + ")", "as.Date(" + ct() + ", tz = \"" + pick(zones) + "\")",
			"as.POSIXct(as.Date(\"" + day() + "\"))", "as.POSIXct(as.POSIXlt(" + ct() + "))", "as.POSIXct(" + fmt.Sprint(r.Intn(2000000000)) + ", tz = \"" + pick(zones) + "\")",
			"weekdays(" + dt() + ")", "months(" + ct() + ", TRUE)", "quarters(" + dt() + ")", "julian(" + dt() + ")",
		})
	case 10:
		family, body = "compare", pick([]string{ct() + " < " + ct(), dt() + " >= " + dt(), ct() + " == " + ct(), "range(" + cts(3) + ")", "rev(sort(" + cts(3) + "))", "diff(" + cts(3) + ")", "order(" + cts(4) + ")", "unique(c(" + dt() + ", " + dt() + "))"})
	case 11:
		family, body = "combine", pick([]string{"c(" + dt() + ", " + dt() + ")", "c(" + ct() + ", " + ct() + ")", "rep(" + dt() + ", 2)", "mean(" + cts(3) + ")", "(" + cts(3) + ")[2]", "max(c(" + dt() + ", " + dt() + "))"})
	case 12:
		family, body = "difftime-ops", pick([]string{
			"as.difftime(" + fmt.Sprint(1+r.Intn(50)) + ", units = " + pick(units[:5]) + ") + as.difftime(" + fmt.Sprint(1+r.Intn(50)) + ", units = " + pick(units[:5]) + ")",
			"as.difftime(" + fmt.Sprint(1+r.Intn(50)) + ", units = " + pick(units[:5]) + ") * " + fmt.Sprint(1+r.Intn(5)),
			"sum(as.difftime(c(1, 2), units = " + pick(units[:5]) + "), as.difftime(3, units = " + pick(units[:5]) + "))",
			"round(as.difftime(" + fmt.Sprint(r.Intn(1000)) + ".567, units = " + pick(units[:5]) + "), 1)",
			"format(as.difftime(" + fmt.Sprint(r.Intn(100)) + ", units = " + pick(units[:5]) + "))",
		})
	default:
		family, body = "tz-shift", "format("+ct()+", \"%Y-%m-%d %H:%M:%S %Z\", tz = \""+pick(zones)+"\")"
	}
	return "# " + family + "\n" + body
}

func fuzzDatetimeExcluded(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	var class, prog string
	switch r.Intn(4) {
	case 3:
		class = exclPost2037
		prog = fmt.Sprintf("format(as.POSIXct(\"%04d-%02d-%02d %02d:%02d:00\", tz = \"%s\"), \"%%s|%%Z|%%H|%%z\")",
			2038+r.Intn(60), 1+r.Intn(12), 1+r.Intn(28), r.Intn(24), r.Intn(60),
			pick([]string{"Europe/Paris", "America/New_York", "Australia/Sydney", "America/Los_Angeles"}))
	case 0:
		class = exclPre1902
		y := 1800 + r.Intn(102)
		prog = pick([]string{
			fmt.Sprintf("format(as.POSIXct(\"%04d-%02d-%02d %02d:%02d:00\", tz = \"UTC\"), \"%%Y-%%m-%%d %%H:%%M %%j %%u\")", y, 1+r.Intn(12), 1+r.Intn(28), r.Intn(24), r.Intn(60)),
			fmt.Sprintf("as.Date(\"%04d-%02d-%02d\") + 1000", y, 1+r.Intn(12), 1+r.Intn(28)),
			fmt.Sprintf("weekdays(as.Date(\"%04d-%02d-%02d\"))", y, 1+r.Intn(12), 1+r.Intn(28)),
			fmt.Sprintf("unclass(as.POSIXct(\"%04d-01-01\", tz = \"UTC\"))", y),
		})
	case 1:
		class = exclTzAbbrev
		prog = "attributes(as.POSIXlt(as.POSIXct(\"2024-06-01 12:00:00\", tz = \"" +
			pick([]string{"Asia/Kolkata", "America/Sao_Paulo", "Asia/Dubai", "Africa/Cairo", "Pacific/Auckland", "Asia/Tokyo", "Europe/Moscow", "America/Phoenix", "Australia/Perth"}) + "\")))$tzone"
	default:
		class = exclDstGap
		prog = pick([]string{
			"as.POSIXct(\"2024-03-31 02:30:00\", tz = \"Europe/Paris\")", "as.POSIXct(\"2024-10-27 02:30:00\", tz = \"Europe/Paris\")",
			"as.POSIXct(\"2024-03-10 02:30:00\", tz = \"America/New_York\")", "as.POSIXct(\"2024-11-03 01:30:00\", tz = \"America/New_York\")",
			"as.POSIXct(\"2024-10-06 02:30:00\", tz = \"Australia/Sydney\")", "as.POSIXct(\"2024-04-07 02:30:00\", tz = \"Australia/Sydney\")",
			"as.POSIXct(\"2024-03-31 01:30:00\", tz = \"Europe/London\")", "as.POSIXct(\"2024-10-27 01:30:00\", tz = \"Europe/London\")",
		})
	}
	return "# excl:" + class + "\n" + prog
}
