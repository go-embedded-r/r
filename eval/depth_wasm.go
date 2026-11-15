//go:build js || wasip1

package eval

// DefaultMaxDepth is the default of Options.MaxDepth: R's option
// "expressions" (5000) on the other platforms, 1000 here. A WebAssembly
// engine runs Go calls on its own native stack, which no Go code can grow:
// V8's default is 984 KB, and one nesting costs up to about 360 bytes of it
// (an R function calling itself through sapply, measured with node 26 as a
// JS RangeError at a MaxDepth of 2717). 1000 keeps a margin of 2.7 over that
// in node and about 1.7 in Chromium, whose frames are larger. A host that
// knows its stack is larger sets Options.MaxDepth.
const DefaultMaxDepth = 1000
