package syntax_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/go-embedded-r/r/syntax"
)

func TestTopLevelNames(t *testing.T) {
	cases := []struct {
		name, src string
		assigned  []string
		read      []string
		dynamic   string // the expected reason, "" for static
	}{
		{"arrow", "x <- 1", []string{"x"}, nil, ""},
		{"equals", "x = 1", []string{"x"}, nil, ""},
		{"right arrow", "1 -> x", []string{"x"}, nil, ""},
		{"string target", `"x" <- 1`, []string{"x"}, nil, ""},
		{"chain", "a <- b <- 1", []string{"b", "a"}, nil, ""},
		{"reads", "y <- x + z", []string{"y"}, []string{"x", "z"}, ""},
		{"read after assign is internal", "x <- 1; y <- x", []string{"x", "y"}, nil, ""},
		{"read before assign is external", "y <- x; x <- 1", []string{"y", "x"}, []string{"x"}, ""},
		{"self update", "x <- x + 1", []string{"x"}, []string{"x"}, ""},
		{"function definition", "f <- function(a, b = 2) a + b + k", []string{"f"}, []string{"k"}, ""},
		{"lambda", "g <- \\(x) x + w", []string{"g"}, []string{"w"}, ""},
		{"function locals", "f <- function() { v <- 1; v + u }", []string{"f"}, []string{"u"}, ""},
		{"function sees later top-level definition", "f <- function() g(); g <- function() 1", []string{"f", "g"}, nil, ""},
		{"calls are reads", "print(x)", nil, []string{"print", "x"}, ""},
		{"namespaced call", "stats::sd(x)", nil, []string{"x"}, ""},
		{"complex assignment", "names(x)[2] <- 'a'", []string{"x"}, []string{"x"}, ""},
		{"dollar assignment", "df$col <- v", []string{"df"}, []string{"v", "df"}, ""},
		{"dollar read ignores the field", "df$col", nil, []string{"df"}, ""},
		{"slot", "obj@slot <- 1", []string{"obj"}, []string{"obj"}, ""},
		{"for variable", "for (i in xs) total <- total + i", []string{"i", "total"}, []string{"xs", "total"}, ""},
		{"assign literal", `assign("made", 1)`, []string{"made"}, []string{"assign"}[:0], ""},
		{"assign literal with value read", `assign("made", y)`, []string{"made"}, []string{"y"}, ""},
		{"get literal", `get("there")`, nil, []string{"there"}, ""},
		{"local is local", "local({ hidden <- 1; hidden + seen })", nil, []string{"seen"}, ""},
		{"pipe", "r <- xs |> sort() |> head(n = 3)", []string{"r"}, []string{"xs", "sort", "head"}, ""},
		{"named arguments are not reads", "f(a = 1)", nil, []string{"f"}, ""},
		{"formula", "m <- lm(y ~ x, data = d)", []string{"m"}, []string{"lm", "y", "x", "d"}, ""},
		{"quote", "e <- quote(a + b)", []string{"e"}, nil, ""},
		{"dots", "f <- function(...) list(...)", []string{"f"}, []string{"list"}, ""},
		{"parameters shadow", "x <- 1; f <- function(x) x", []string{"x", "f"}, nil, ""},
		{"backquoted", "`my var` <- 2", []string{"my var"}, nil, ""},
		{"if", "if (cond) a <- 1 else b <- 2", []string{"a", "b"}, []string{"cond"}, ""},
		{"empty", "", nil, nil, ""},

		{"superassign", "x <<- 1", nil, nil, "<<-"},
		{"superassign in function", "f <- function() count <<- count + 1", []string{"f"}, []string{"count"}, "<<-"},
		{"right superassign", "1 ->> y", nil, nil, "<<-"},
		{"assign computed", `assign(nm, 1)`, nil, []string{"nm"}, "assign with a computed name"},
		{"assign paste", `assign(paste0("v", i), 1)`, nil, []string{"paste0", "i"}, "assign with a computed name"},
		{"assign other env", `assign("x", 1, envir = e)`, nil, []string{"e"}, "assign into another environment"},
		{"assign pos", `assign("x", 1, 1)`, nil, nil, "assign into another environment"},
		{"get computed", "get(nm)", nil, nil, "get with a computed name"},
		{"exists computed", "exists(paste('a', 'b'))", nil, nil, "exists with a computed name"},
		{"eval", "eval(expr)", nil, []string{"eval", "expr"}, "eval"},
		{"base::eval", "base::eval(expr)", nil, []string{"expr"}, "eval"},
		{"parse", "parse(text = s)", nil, []string{"parse", "s"}, "parse"},
		{"source", `source("f.R")`, nil, []string{"source"}, "source"},
		{"attach", "attach(df)", nil, []string{"attach", "df"}, "attach"},
		{"library", "library(dplyr)", nil, nil, "library"},
		{"require", `require("dplyr")`, nil, nil, "require"},
		{"sys.function", "f <- function() sys.function()", []string{"f"}, []string{"sys.function"}, "sys.function"},
		{"environment<-", "environment(f) <- e", []string{"f"}, []string{"e", "f"}, "environment<-"},
		{"list2env", "list2env(l, globalenv())", nil, []string{"list2env", "l", "globalenv"}, "list2env"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			n, err := syntax.TopLevelNames(c.src)
			if err != nil {
				t.Fatal(err)
			}
			if !sameSet(n.Assigned, c.assigned) {
				t.Errorf("assigned %v, want %v", n.Assigned, c.assigned)
			}
			if !sameSet(n.Read, c.read) {
				t.Errorf("read %v, want %v", n.Read, c.read)
			}
			if n.Dynamic != (c.dynamic != "") {
				t.Errorf("dynamic = %v (%v), want reason %q", n.Dynamic, n.Reasons, c.dynamic)
			}
			if c.dynamic != "" && !contains(n.Reasons, c.dynamic) {
				t.Errorf("reasons %v lack %q", n.Reasons, c.dynamic)
			}
		})
	}
}

func contains(s []string, x string) bool {
	for _, y := range s {
		if y == x {
			return true
		}
	}
	return false
}

func sameSet(a, b []string) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	x, y := map[string]bool{}, map[string]bool{}
	for _, s := range a {
		x[s] = true
	}
	for _, s := range b {
		y[s] = true
	}
	return reflect.DeepEqual(x, y)
}

func TestTopLevelNamesOrderAndErrors(t *testing.T) {
	n, err := syntax.TopLevelNames("b <- 1\na <- 2\nb <- 3")
	if err != nil || strings.Join(n.Assigned, ",") != "b,a" {
		t.Errorf("%v %v", n, err)
	}
	n, _ = syntax.TopLevelNames("eval(a); eval(b); x <<- 1")
	if strings.Join(n.Reasons, ",") != "eval,<<-" {
		t.Errorf("reasons %v", n.Reasons)
	}
	if _, err := syntax.TopLevelNames("x <- "); err == nil {
		t.Error("a syntax error must be reported")
	}
}
