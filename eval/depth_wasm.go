//go:build js || wasip1

package eval

// DefaultMaxDepth is the default of Options.MaxDepth: R's option
// "expressions" (5000) on the other platforms, 300 here. A WebAssembly engine
// runs Go calls on its own native stack, which no Go code can grow, and a
// stack overflow there is not a Go panic: it kills the instance (a JS
// RangeError). Measured with node 26 and V8's default stack of 984 KB, the
// largest MaxDepth that survives is 4265 for f <- function() f(), 2721 for a
// function calling itself through sapply, 1961 for all.equal on nested lists
// and 832 for the same once the first evaluation left garbage behind (the
// allocator and the collector run on top of the deep stack). Chromium's
// frames are larger by about a third. 300 keeps a margin of 2.7 over the
// worst of those in node and about 1.7 in Chromium. A host that knows its
// stack is larger (wazero's is far larger: 80359 for the simplest function)
// sets Options.MaxDepth.
const DefaultMaxDepth = 300
