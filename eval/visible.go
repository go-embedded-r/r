package eval

// The functionals call R functions, and what those evaluated last (an
// assignment, invisible()) must not decide whether the result of the
// functional is printed: the result of sapply(1:3, function(i) x <- i) is
// visible. They are wrapped here, in a file that sorts after the ones that
// register them.
func init() {
	for _, name := range []string{"lapply", "sapply", "vapply", "Map", "mapply", "Reduce", "Filter", "Find", "Position"} {
		b := builtins[name]
		if b == nil {
			continue
		}
		inner := b.fn
		b.fn = func(s *Session, c *callInfo) Value {
			v := inner(s, c)
			s.visible = true
			return v
		}
	}
}
