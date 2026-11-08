package main

import (
	"math/rand"
	"regexp"
	"strings"
)

// fuzzEnv builds a random program over environments: bindings, locks, active
// bindings, delayed assignment, lookups with modes and defaults, the
// namespace functions, with/within/eapply. Every program ends in a value that
// is not an environment (an environment's address is not a value), and every
// listing is sorted (R's order for a hashed environment is not defined). The
// first line is a comment that names the family, for the per-family count.
func fuzzEnv(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	if r.Intn(8) == 0 {
		return fuzzEnvExcluded(r)
	}
	names := []string{"a", "b", "c", ".h", "zz"}
	// no closures among the values: R prints a closure defined here with its
	// environment's address, which is not a value to compare
	vals := []string{"1", "\"s\"", "TRUE", "NULL", "c(1, 2)", "list(1)", "2L", "NA", "\"t\""}
	modes := []string{"\"any\"", "\"numeric\"", "\"character\"", "\"function\"", "\"logical\"", "\"list\"", "\"NULL\"", "c(\"numeric\", \"character\")"}
	nm := func() string { return "\"" + pick(names) + "\"" }
	flag := func() string { return pick([]string{"TRUE", "FALSE", "NA", "1", "0"}) }

	// the environment under test: some bindings, maybe locks
	var setup strings.Builder
	setup.WriteString("e <- new.env(); ")
	for _, n := range []string{"a", "b", "c", ".h"} {
		if r.Intn(2) == 0 {
			setup.WriteString("assign(\"" + n + "\", " + pick(vals) + ", e); ")
		}
	}
	if r.Intn(3) == 0 {
		setup.WriteString("tryCatch(lockBinding(" + nm() + ", e), error = function(err) NULL); ")
	}
	if r.Intn(4) == 0 {
		setup.WriteString("lockEnvironment(e" + pick([]string{"", ", TRUE", ", FALSE"}) + "); ")
	}
	probe := func(body string) string {
		// the outcome of body, whatever it does, then the visible state of e
		return setup.String() + body + "; sort(ls(e, all.names = TRUE))"
	}
	// the packages R installs besides base are left out: their namespaces
	// resolve through our base environment, not through their own exports
	nsArgs := []string{"\"base\"", "\"foo\"", "1", "NULL", "baseenv()", "globalenv()", "asNamespace(\"base\")", "c(\"base\", \"foo\")", "quote(base)"}
	nsStrArgs := []string{"\"base\"", "\"foo\"", "1", "quote(base)", "c(\"base\", \"foo\")"}
	nsNames := []string{"sum", "pi", "paste", "nope", "T", "letters", "lapply", "`+`"}

	family, body := "", ""
	switch r.Intn(16) {
	case 0:
		family, body = "get", "get("+nm()+", envir = e, mode = "+pick(modes)+", inherits = "+flag()+")"
	case 1:
		family, body = "exists", "exists("+nm()+", envir = e, mode = "+pick(modes)+", inherits = "+flag()+")"
	case 2:
		family, body = "get0", "get0("+nm()+", envir = e, mode = "+pick(modes)+", inherits = "+flag()+", ifnotfound = "+pick(vals)+")"
	case 3:
		family = "mget"
		body = "mget(c(" + nm() + ", " + nm() + "), envir = e, mode = " + pick(modes) + ", inherits = " + flag()
		if r.Intn(2) == 0 {
			body += ", ifnotfound = " + pick([]string{"NA", "list(NA)", "list(0, 1)", "list(function(x) x)", "list(1, 2, 3)", "\"none\""})
		}
		body += ")"
		body = "r <- tryCatch(" + body + ", error = function(err) conditionMessage(err)); if (is.list(r)) lapply(r[order(names(r))], function(v) if (is.function(v)) \"<function>\" else v) else r"
		if r.Intn(2) == 0 {
			body = "tryCatch(mget(c(" + nm() + ", " + nm() + "), envir = e), error = function(err) conditionMessage(err))"
		}
	case 4:
		family, body = "assign", "assign("+nm()+", "+pick(vals)+", envir = e)"
	case 5:
		family = "subassign"
		body = pick([]string{"e$" + pick([]string{"a", "b", "new"}) + " <- " + pick(vals), "e[[" + nm() + "]] <- " + pick(vals), "local(" + pick([]string{"a", "b", "zz"}) + " <- " + pick(vals) + ", e)", "evalq(" + pick([]string{"a", "b"}) + " <- " + pick(vals) + ", e)"})
	case 6:
		family = "rm"
		body = "rm(list = c(" + nm() + ", " + nm() + "), envir = e" + pick([]string{"", ", inherits = TRUE", ", inherits = FALSE"}) + ")"
	case 7:
		family = "locks"
		body = pick([]string{
			"lockBinding(" + nm() + ", e)", "unlockBinding(" + nm() + ", e)", "bindingIsLocked(" + nm() + ", e)",
			"environmentIsLocked(e)", "lockEnvironment(e, " + flag() + ")", "bindingIsLocked(quote(" + pick(names[:3]) + "), e)",
			"lockBinding(" + pick([]string{"1", "NULL", "quote(a)", "c(\"a\", \"b\")"}) + ", " + pick([]string{"e", "5", "globalenv()"}) + ")",
		})
	case 8:
		family = "active"
		body = pick([]string{
			"makeActiveBinding(" + nm() + ", function(v) if (missing(v)) 42 else cat(\"set\", v, \"\\n\"), e); e[[" + nm() + "]]",
			"makeActiveBinding(" + nm() + ", " + pick([]string{"function() 1", "5", "sum", "NULL"}) + ", e); bindingIsActive(" + nm() + ", e)",
			"bindingIsActive(" + nm() + ", e)", "activeBindingFunction(" + nm() + ", e)",
			"makeActiveBinding(\"q\", function(v) if (missing(v)) 7 else cat(\"w\", v, \"\\n\"), e); e$q <- 1; get(\"q\", e)",
		})
	case 9:
		family = "delayed"
		body = pick([]string{
			"delayedAssign(" + nm() + ", {cat(\"forced\\n\"); 9}, assign.env = e); " + pick([]string{"exists(\"a\", e)", "get(\"a\", e)", "e$a", "mget(\"a\", envir = e, ifnotfound = NA)"}),
			"delayedAssign(" + pick([]string{"1", "NULL", "c(\"a\", \"b\")", "\"x\""}) + ", 1, assign.env = e)",
			"delayedAssign(\"w\", 1 + 1, eval.env = e, assign.env = e); w <- 1; get(\"w\", e)",
		})
	case 10:
		family = "aslist"
		body = pick([]string{
			"as.list(e, all.names = " + flag() + ", sorted = TRUE)",
			"names(as.list(e, all.names = " + flag() + ", sorted = " + flag() + "))[0]",
			"length(as.list(e, all.names = " + flag() + "))",
			"as.list(" + pick([]string{"1:3", "c(a = 1)", "list(1, b = 2)", "\"x\"", "NULL", "TRUE", "quote(f(x, y = 2))", "factor(c(\"u\", \"v\"))", "function(x, y = 2) x + y", "NA", "new.env()[0]"}) + ")",
		})
	case 11:
		family = "eapply"
		body = pick([]string{
			"r <- eapply(e, function(v) length(v), all.names = " + flag() + "); r[order(names(r))]",
			"r <- eapply(e, class, USE.NAMES = " + flag() + "); length(r)",
		})
	case 12:
		family = "with"
		body = pick([]string{
			"with(list(a = 1, b = 2), " + pick([]string{"a + b", "a * zz", "c(a, b)", "names(environment())[0]"}) + ")",
			"with(e, " + pick([]string{"a", "b", "exists(\"c\")", "zz"}) + ")",
			"within(list(a = 1, b = 2), " + pick([]string{"c <- a + b", "{ d <- 1; e2 <- 2 }", "rm(a)", "a <- NULL", "{ b <- b + 1; z <- 0 }", "NULL"}) + ")",
		})
	case 13:
		family = "namespace"
		body = pick([]string{
			"asNamespace(" + pick(nsArgs) + ")", "getNamespace(" + pick(nsStrArgs) + ")", "isNamespace(" + pick(nsArgs) + ")",
			"isBaseNamespace(" + pick(nsArgs) + ")", "getNamespaceName(" + pick(nsArgs) + ")",
			"isNamespaceLoaded(" + pick([]string{"\"base\"", "\"foo\"", "quote(base)"}) + ")",
			"requireNamespace(" + pick([]string{"\"base\"", "\"foo\""}) + ", quietly = TRUE)",
			"environmentName(" + pick(nsArgs) + ")", "identical(" + pick([]string{"\"base\"", "baseenv()", "asNamespace(\"base\")"}) + ", baseenv())",
		})
	case 14:
		family = "colon"
		pkg := pick([]string{"base", "base", "foo", "\"base\""})
		op := pick([]string{"::", ":::"})
		name := pick(nsNames)
		body = pkg + op + name
		switch name {
		case "pi", "T", "letters":
		default:
			// a closure of R is byte-compiled, which a serialization cannot carry
			body = "is.function(" + body + ")"
		}
	default:
		family = "topenv"
		body = pick([]string{
			"environmentName(topenv(" + pick([]string{"e", "globalenv()", "baseenv()", "emptyenv()", "asNamespace(\"base\")", "new.env(parent = baseenv())", "new.env(parent = emptyenv())", "parent.env(e)"}) + "))",
			"environmentName(parent.env(" + pick([]string{"asNamespace(\"base\")", "baseenv()", "e"}) + "))",
			"environmentName(" + pick([]string{"e", "5", "NULL", "\"base\"", "emptyenv()", "baseenv()", "asNamespace(\"base\")"}) + ")",
			"{ attr(e, \"name\") <- " + pick([]string{"\"nm\"", "1", "NULL"}) + "; environmentName(e) }",
		})
	}
	if family == "get" || family == "get0" {
		// a function found in base prints differently (its formals); only
		// the fact that one was found is compared
		body = "x <- " + body + "; if (is.function(x)) \"<function>\" else x"
	}
	prog := body
	if family != "namespace" && family != "colon" && family != "topenv" && !strings.Contains(body, "sort(ls") {
		prog = probe(body)
	} else if family == "topenv" || family == "namespace" || family == "colon" {
		prog = setup.String() + body
	}
	return "# " + family + "\n" + prog
}

// The classes of case that cannot be compared exactly, each generated on
// purpose so that the report counts them (the first line of the program is
// "# excl:<class>"):
//
//	addr       R prints an environment's address; compared after the address is
//	           replaced by a placeholder on both sides
//	builtin    R prints a builtin with its formals, ours does not; compared
//	           after both are reduced to .Primitive("name")
//	installed  a package R installs besides base resolves through our base
//	           environment; compared exactly, a difference is a known divergence
//	hash       the order of an unsorted listing of a hashed environment is not
//	           defined by R; compared exactly, a difference is a known divergence
const (
	exclAddr      = "addr"
	exclBuiltin   = "builtin"
	exclInstalled = "installed"
	exclHash      = "hash"
)

func fuzzEnvExcluded(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	var class, prog string
	switch r.Intn(4) {
	case 0:
		class = exclAddr
		prog = "f <- function(x) x; e <- new.env(); assign(\"f\", f, e); print(" +
			pick([]string{"get(\"f\", e)", "mget(\"f\", envir = e)", "as.list(e)", "f", "eapply(e, identity)", "list(f, 1)"}) + "); 0"
	case 1:
		class = exclBuiltin
		prog = "print(" + pick([]string{
			"get(\"sum\", mode = \"function\")", "mget(\"c\", envir = baseenv())", "sum", "`+`", "list(length, 1)",
			"get0(\"max\", mode = \"function\")", "base::sum", "asNamespace(\"base\")$min",
		}) + "); 0"
	case 2:
		class = exclInstalled
		prog = pick([]string{
			"stats::paste", "is.function(utils::head)", "is.function(stats::sd)", "utils::sum(1, 2)",
			"environmentName(asNamespace(\"stats\"))", "isNamespace(asNamespace(\"utils\"))", "stats:::nope", "isNamespaceLoaded(\"stats\")",
			"requireNamespace(\"stats\", quietly = TRUE)",
		})
	default:
		class = exclHash
		prog = "e <- new.env(); " + pick([]string{"assign(\"b\", 1, e); assign(\"a\", 2, e); assign(\"c\", 3, e); ", "e$x <- 1; e$y <- 2; e$z <- 3; "}) +
			pick([]string{"names(as.list(e))", "names(e)", "names(eapply(e, identity))", "unlist(as.list(e))"})
	}
	return "# excl:" + class + "\n" + prog
}

var (
	addrRE    = regexp.MustCompile(`<(environment|bytecode): 0x[0-9a-f]+>`)
	builtinRE = regexp.MustCompile(`function \([^)]*\) +\.Primitive\("([^"]*)"\)`)
)

// normAddr replaces the address R prints for an environment by a placeholder.
func normAddr(s string) string { return addrRE.ReplaceAllString(s, "<$1>") }

// normBuiltin reduces the printed form of a builtin to its name.
func normBuiltin(s string) string { return builtinRE.ReplaceAllString(s, ".Primitive(\"$1\")") }
