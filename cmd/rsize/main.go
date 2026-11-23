// Command rsize is the minimal program the size budget measures: one session
// and one evaluation, and nothing else. Its size under GOOS=js and
// GOOS=wasip1 is what a host pays for the interpreter (go-embedded-r/r#62).
package main

import (
	"fmt"
	"os"

	"github.com/go-embedded-r/r/eval"
)

func main() {
	s := eval.NewSession(eval.Options{})
	res, err := s.Eval("1 + 1")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	_ = res
}
