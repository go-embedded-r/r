//go:build !js && !wasip1

package eval

// DefaultMaxDepth is the default of Options.MaxDepth: R's option
// "expressions". See depth_wasm.go for the value under js/wasm and wasip1.
const DefaultMaxDepth = 5000
