package evalcheck

import "os"

// The oracle runs R with TZ=UTC (internal/oracle), so the session's own time
// zone, the one a date-time without a tz argument is read in, is UTC there.
// The Go evaluator reads TZ the same way, and is given the same value, so
// that a comparison does not depend on the machine's zone.
func init() { os.Setenv("TZ", "UTC") }
