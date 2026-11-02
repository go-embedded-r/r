package syntax_test

import (
	"errors"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/go-embedded-r/r/syntax"
)

func TestErrorCarriesPosition(t *testing.T) {
	_, err := syntax.ParseAll("x <- 1\ny <- 2 3")
	var e *syntax.Error
	if !errors.As(err, &e) {
		t.Fatalf("%T %v", err, err)
	}
	if e.Line != 2 || e.Col != 8 || e.Message != "unexpected numeric constant" || e.Lexical {
		t.Errorf("%+v", e)
	}
	if !strings.HasPrefix(e.Error(), "<text>:2:8: unexpected numeric constant\n1: x <- 1\n2: y <- 2 3\n") {
		t.Errorf("%q", e.Error())
	}
	_, err = syntax.ParseAll(`"\q"`)
	if !errors.As(err, &e) || !e.Lexical || e.Line != 1 || e.Col != 3 || e.Message != `'\q' is an unrecognized escape in character string` {
		t.Errorf("%v %+v", err, e)
	}
	f, err := syntax.ParseAll("a <- 1\nb <- 2 3")
	if err == nil || f == nil || len(f.Exprs) != 1 {
		t.Errorf("expressions before the error are kept: %v %v", f, err)
	}
}

func TestLimits(t *testing.T) {
	_, err := syntax.ParseAllLimits(strings.Repeat("x\n", 100), syntax.Limits{MaxBytes: 50})
	var le *syntax.LimitError
	if !errors.Is(err, syntax.ErrLimit) || !errors.As(err, &le) || le.What != "bytes" || le.Limit != 50 || !strings.Contains(err.Error(), "bytes") {
		t.Errorf("bytes: %v", err)
	}
	_, err = syntax.ParseAllLimits(strings.Repeat("x\n", 100), syntax.Limits{MaxNodes: 10})
	if !errors.As(err, &le) || le.What != "nodes" {
		t.Errorf("nodes: %v", err)
	}
	_, err = syntax.ParseAllLimits("f(g(h(i(1))))", syntax.Limits{MaxDepth: 3})
	if !errors.As(err, &le) || le.What != "depth" {
		t.Errorf("depth: %v", err)
	}
	if _, err := syntax.ParseAllLimits("f(g(h(i(1))))", syntax.Limits{MaxDepth: 20}); err != nil {
		t.Error(err)
	}
}

// The bombs: inputs that are small and would build a huge or very deep tree
// must be refused, under a memory limit, without crashing.
func TestBombs(t *testing.T) {
	old := debug.SetMemoryLimit(512 << 20)
	defer debug.SetMemoryLimit(old)
	n := 200000
	bombs := map[string]string{
		"left chain":     strings.Repeat("a+", n) + "a",
		"unary chain":    strings.Repeat("-", n) + "a",
		"not chain":      strings.Repeat("!", n) + "a",
		"call chain":     "f" + strings.Repeat("(1)", n),
		"subscript":      "x" + strings.Repeat("[1]", n),
		"dollar chain":   "x" + strings.Repeat("$a", n),
		"right assoc":    strings.Repeat("a^", n) + "a",
		"assign chain":   strings.Repeat("a<-", n) + "a",
		"parens":         strings.Repeat("(", n) + "a" + strings.Repeat(")", n),
		"braces":         strings.Repeat("{", n),
		"brackets":       "x" + strings.Repeat("[[", n),
		"if chain":       strings.Repeat("if(a)b else ", n) + "c",
		"function chain": strings.Repeat("function(x)", n) + "x",
		"lambda chain":   strings.Repeat(`\(x)`, n) + "x",
		"pipe chain":     "x" + strings.Repeat(" |> f()", n),
		"while chain":    strings.Repeat("while(a)", n) + "b",
		"formula chain":  strings.Repeat("~", n) + "a",
		"nested ifs":     strings.Repeat("{if(a)", n),
		"call args":      strings.Repeat("f(", n),
	}
	for name, src := range bombs {
		f, err := syntax.ParseAll(src)
		if err == nil {
			// a chain is accepted only when it stays within the limits
			t.Errorf("%s: accepted (%d expressions)", name, len(f.Exprs))
			continue
		}
		var se *syntax.Error
		if !errors.Is(err, syntax.ErrLimit) && !errors.As(err, &se) {
			t.Errorf("%s: %T %v", name, err, err)
		}
	}
	// the context stack of R itself: 50 levels of brackets
	_, err := syntax.ParseAll(strings.Repeat("(", 60) + "a" + strings.Repeat(")", 60))
	if err == nil || !strings.Contains(err.Error(), "contextstack overflow") {
		t.Errorf("60 parentheses: %v", err)
	}
	if _, err := syntax.ParseAll(strings.Repeat("(", 40) + "a" + strings.Repeat(")", 40)); err != nil {
		t.Errorf("40 parentheses: %v", err)
	}
	_, err = syntax.ParseAll(strings.Repeat("x[[", 30))
	if err == nil || !strings.Contains(err.Error(), "contextstack overflow") {
		t.Errorf("30 [[: %v", err)
	}
	// a long but flat input is fine
	if _, err := syntax.ParseAllLimits(strings.Repeat("x <- 1\n", 100000), syntax.Limits{}); err != nil {
		t.Error(err)
	}
	// a large string literal
	if _, err := syntax.ParseAll(`"` + strings.Repeat("a", 5<<20) + `"`); err != nil {
		t.Error(err)
	}
}

func TestSpansAndTokens(t *testing.T) {
	f, err := syntax.ParseAll("x <- f(a, 'é')  # c\ny")
	if err != nil {
		t.Fatal(err)
	}
	asg := f.Exprs[0].(*syntax.Call)
	if asg.Pos.Start != (syntax.Pos{Line: 1, Col: 1}) || asg.Pos.End != (syntax.Pos{Line: 1, Col: 14}) {
		t.Errorf("assignment span %+v", asg.Pos)
	}
	if got := string(f.Source[asg.Pos.StartByte:asg.Pos.EndByte]); got != "x <- f(a, 'é')" {
		t.Errorf("source of the assignment: %q", got)
	}
	call := asg.Args[1].Value.(*syntax.Call)
	if call.FunName() != "f" || len(call.Args) != 2 {
		t.Errorf("%+v", call)
	}
	var kinds []string
	for _, tk := range f.Tokens {
		kinds = append(kinds, tk.Kind)
	}
	want := "SYMBOL LEFT_ASSIGN SYMBOL_FUNCTION_CALL '(' SYMBOL ',' STR_CONST ')' COMMENT SYMBOL"
	if strings.Join(kinds, " ") != want {
		t.Errorf("tokens %v", kinds)
	}
	if f.Tokens[0].Text != "x" || f.Tokens[6].Text != "'é'" {
		t.Errorf("%+v", f.Tokens)
	}
}

func TestTreeShape(t *testing.T) {
	f, err := syntax.ParseAll("function(a, b = 2) a\n\\(z) z\nf(x = 1, , 'y' = 2)\n`if`(a)\n'g'(1)\nx |> f(y = 2)")
	if err != nil {
		t.Fatal(err)
	}
	fn := f.Exprs[0].(*syntax.Function)
	if len(fn.Params) != 2 || fn.Params[1].Default == nil || fn.Params[0].Default != nil || fn.Lambda {
		t.Errorf("%+v", fn)
	}
	if !f.Exprs[1].(*syntax.Function).Lambda {
		t.Error("lambda")
	}
	c := f.Exprs[2].(*syntax.Call)
	if len(c.Args) != 3 || !c.Args[0].HasName || c.Args[1].Value != nil || c.Args[1].HasName || c.Args[2].Name != "y" {
		t.Errorf("%+v", c.Args)
	}
	if f.Exprs[4].(*syntax.Call).FunName() != "g" {
		t.Error("a string callee is a symbol")
	}
	p := f.Exprs[5].(*syntax.Call)
	if p.FunName() != "f" || len(p.Args) != 2 || p.Args[0].Value.(*syntax.Sym).Name != "x" {
		t.Errorf("pipe: %+v", p)
	}
}

func TestConstants(t *testing.T) {
	f, err := syntax.ParseAll("1L; 1.5; 0x10L; 2i; 1e3L; 1.5L; TRUE; NA; NA_integer_; NA_character_; NULL; Inf; 'a'; 100000L")
	if err != nil {
		t.Fatal(err)
	}
	kinds := []syntax.ConstKind{syntax.Integer, syntax.Double, syntax.Integer, syntax.Complex, syntax.Integer,
		syntax.Double, syntax.Logical, syntax.Logical, syntax.Integer, syntax.String, syntax.Null, syntax.Double, syntax.String, syntax.Integer}
	for i, e := range f.Exprs {
		if e.(*syntax.Const).Kind != kinds[i] {
			t.Errorf("#%d kind %v want %v", i, e.(*syntax.Const).Kind, kinds[i])
		}
	}
	if f.Exprs[2].(*syntax.Const).Int != 16 || f.Exprs[0].(*syntax.Const).Int != 1 {
		t.Error("integer values")
	}
}

func TestEmptyAndComments(t *testing.T) {
	for _, src := range []string{"", "\n\n", "# only a comment", ";", "  \t "} {
		f, err := syntax.ParseAll(src)
		if src == ";" {
			if err == nil {
				t.Errorf("%q accepted", src)
			}
			continue
		}
		if err != nil || len(f.Exprs) != 0 {
			t.Errorf("%q: %v %v", src, f, err)
		}
	}
}
