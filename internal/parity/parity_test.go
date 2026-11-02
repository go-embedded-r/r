package parity

import (
	"os"
	"testing"

	"github.com/go-embedded-r/r/internal/oracle"
)

func TestCompareLive(t *testing.T) {
	p, err := oracle.Start(oracle.Options{})
	if oracle.Unavailable(err) {
		if os.Getenv("RGO_REQUIRE_ORACLE") != "" {
			t.Fatal(err)
		}
		t.Skip(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	for _, src := range append([]string{"x <- f(1, b = 2)", "1 +", "a b", "0.1 + 1e300"}, Extras()...) {
		o, err := Compare(p, src)
		if err != nil {
			t.Fatal(err)
		}
		if !o.Tree || !o.Tokens || !o.Error {
			t.Errorf("%q: %+v", src, o)
		}
	}
}

func TestNumClose(t *testing.T) {
	if n, ok := NumClose("..dbl 3ff0000000000000", "..dbl 3ff0000000000002"); !ok || n != 1 {
		t.Errorf("%d %v", n, ok)
	}
	if _, ok := NumClose("..dbl 3ff0000000000000", "..dbl 3ff0000000001000"); ok {
		t.Error("far values are not close")
	}
	if _, ok := NumClose("..sym 61", "..sym 62"); ok {
		t.Error("symbols are not numbers")
	}
	if got := decodeLine("..sym 61"); got != `..sym "a"` {
		t.Error(got)
	}
}
