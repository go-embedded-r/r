package syntax

import "strings"

// Names is what a notebook cell defines and uses at top level, worked out
// from the syntax alone under R's scoping rules.
type Names struct {
	// Assigned are the names the cell binds in the environment it runs in
	// (the top level), in order of first assignment: x <- v, x = v, v -> x,
	// "x" <- v, f(x) <- v and x$a <- v (which assign x), assign("x", v) with
	// a literal name, and the variable of a top-level for loop. Function
	// definitions are assignments of a function value like any other.
	Assigned []string
	// Read are the names the cell uses without having bound them first:
	// free variables of top-level code, and free variables of function bodies
	// that the cell does not itself assign anywhere at top level. Called
	// functions are included (f in f(x)), as are get("x") literals. Names
	// bound inside a function (its parameters and its local assignments) or
	// inside local() are not.
	Read []string
	// Dynamic is set when the set of names cannot be resolved statically: see
	// Reasons for which construct made it so.
	Dynamic bool
	// Reasons lists the constructs that set Dynamic, each once, in order of
	// appearance: "<<-", "eval", "assign with a computed name", ...
	Reasons []string
}

// TopLevelNames parses a cell and returns the names it assigns and reads.
// The error is the parse error (*Error or *LimitError), as from [ParseAll];
// the Names are valid only when it is nil.
func TopLevelNames(cell string) (*Names, error) {
	f, err := ParseAll(cell)
	if err != nil {
		return nil, err
	}
	return NamesOf(f.Exprs), nil
}

// NamesOf is TopLevelNames for expressions that are already parsed.
func NamesOf(exprs []Node) *Names {
	w := &nameWalker{assigned: map[string]bool{}, read: map[string]bool{}, reasons: map[string]bool{}}
	for _, e := range exprs {
		w.walk(e, nil)
	}
	out := &Names{Assigned: w.assignedOrder, Dynamic: len(w.reasonOrder) > 0, Reasons: w.reasonOrder}
	out.Read = append(out.Read, w.readOrder...)
	for _, n := range w.deferred {
		if !w.assigned[n] && !w.read[n] {
			w.read[n] = true
			out.Read = append(out.Read, n)
		}
	}
	return out
}

// dynamicCalls are the functions whose effect on the set of names cannot be
// known statically (or that run code that can do anything).
var dynamicCalls = map[string]bool{
	"eval": true, "evalq": true, "eval.parent": true, "parse": true, "source": true,
	"sys.source": true, "attach": true, "detach": true, "library": true, "require": true,
	"requireNamespace": true, "loadNamespace": true, "load": true, "sys.function": true,
	"list2env": true, "delayedAssign": true, "makeActiveBinding": true, "str2lang": true,
	"str2expression": true,
}

// noArgReads are the calls whose arguments are not variable reads.
var noArgReads = map[string]bool{
	"library": true, "require": true, "requireNamespace": true, "quote": true,
	"expression": true, "bquote": true, "substitute": true,
}

var keywordCalls = map[string]bool{
	"if": true, "for": true, "while": true, "repeat": true, "break": true, "next": true, "function": true,
}

type nscope struct {
	locals map[string]bool
	parent *nscope
	isFunc bool
}

func (s *nscope) has(name string) bool {
	for ; s != nil; s = s.parent {
		if s.locals[name] {
			return true
		}
	}
	return false
}

type nameWalker struct {
	assigned      map[string]bool
	assignedOrder []string
	read          map[string]bool
	readOrder     []string
	deferred      []string
	reasons       map[string]bool
	reasonOrder   []string
}

func (w *nameWalker) dynamic(reason string) {
	if !w.reasons[reason] {
		w.reasons[reason] = true
		w.reasonOrder = append(w.reasonOrder, reason)
	}
}

func (w *nameWalker) bind(name string, sc *nscope) {
	if sc != nil {
		sc.locals[name] = true
		return
	}
	if !w.assigned[name] {
		w.assigned[name] = true
		w.assignedOrder = append(w.assignedOrder, name)
	}
}

func (w *nameWalker) use(name string, sc *nscope, inFunc bool) {
	if name == "..." || (strings.HasPrefix(name, "..") && len(name) > 2 && isAllDigits(name[2:])) {
		return
	}
	if sc.has(name) {
		return
	}
	if inFunc {
		w.deferred = append(w.deferred, name)
		return
	}
	if w.assigned[name] || w.read[name] {
		return
	}
	w.read[name] = true
	w.readOrder = append(w.readOrder, name)
}

func isAllDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return s != ""
}

// inFunction reports whether sc is inside a function body (as opposed to the
// top level or only local() scopes).
func inFunction(sc *nscope) bool {
	for ; sc != nil; sc = sc.parent {
		if sc.isFunc {
			return true
		}
	}
	return false
}

func isCalleeName(name string) bool {
	if name == "" || keywordCalls[name] {
		return false
	}
	c := name[0]
	return c == '.' || c == '%' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c >= 0x80
}

func literalName(n Node) (string, bool) {
	switch x := n.(type) {
	case *Const:
		if x.Kind == String && !x.NA && !x.ph {
			return x.Str, true
		}
	case *Sym:
		return x.Name, true
	}
	return "", false
}

// walk visits n in scope sc (nil at top level).
func (w *nameWalker) walk(n Node, sc *nscope) {
	switch x := n.(type) {
	case nil:
	case *Const:
	case *Sym:
		w.use(x.Name, sc, inFunction(sc))
	case *Function:
		fs := &nscope{locals: map[string]bool{}, parent: sc, isFunc: true}
		for _, p := range x.Params {
			fs.locals[p.Name] = true
		}
		for _, p := range x.Params {
			if p.Default != nil {
				w.walk(p.Default, fs)
			}
		}
		w.walk(x.Body, fs)
	case *Call:
		w.call(x, sc)
	}
}

func (w *nameWalker) walkArgs(args []Arg, sc *nscope) {
	for _, a := range args {
		w.walk(a.Value, sc)
	}
}

func (w *nameWalker) call(c *Call, sc *nscope) {
	name := c.FunName()
	// strip base:: and friends: base::eval is eval
	if q, ok := c.Fun.(*Call); ok && (q.FunName() == "::" || q.FunName() == ":::") && len(q.Args) == 2 {
		if s, ok := q.Args[1].Value.(*Sym); ok {
			name = s.Name
		}
	}
	switch name {
	case "<-", "=":
		if len(c.Args) == 2 {
			w.walk(c.Args[1].Value, sc)
			w.assign(c.Args[0].Value, sc)
			return
		}
	case "<<-":
		w.dynamic("<<-")
		if len(c.Args) == 2 {
			w.walk(c.Args[1].Value, sc)
			w.targetReads(c.Args[0].Value, sc)
			return
		}
	case "$", "@":
		if len(c.Args) >= 1 {
			w.walk(c.Args[0].Value, sc)
			return
		}
	case "::", ":::":
		return
	case "for":
		if len(c.Args) == 3 {
			w.walk(c.Args[1].Value, sc)
			if s, ok := c.Args[0].Value.(*Sym); ok {
				w.bind(s.Name, sc)
			}
			w.walk(c.Args[2].Value, sc)
			return
		}
	case "local":
		if len(c.Args) >= 1 {
			ls := &nscope{locals: map[string]bool{}, parent: sc, isFunc: inFunction(sc)}
			w.walk(c.Args[0].Value, ls)
			for _, a := range c.Args[1:] {
				w.walk(a.Value, sc)
			}
			return
		}
	case "assign":
		w.assignCall(c, sc)
		return
	case "get", "get0", "exists":
		if len(c.Args) >= 1 {
			if lit, ok := stringLiteral(c.Args[0].Value); ok {
				w.use(lit, sc, inFunction(sc))
			} else {
				w.dynamic(name + " with a computed name")
			}
			w.walkArgs(c.Args[1:], sc)
			return
		}
	case "environment":
	}
	if dynamicCalls[name] {
		w.dynamic(name)
	}
	if noArgReads[name] {
		return
	}
	if s, ok := c.Fun.(*Sym); ok {
		if isCalleeName(s.Name) {
			w.use(s.Name, sc, inFunction(sc))
		}
	} else if _, isPkg := c.Fun.(*Call); !isPkg || name == "" {
		w.walk(c.Fun, sc)
	}
	w.walkArgs(c.Args, sc)
}

func stringLiteral(n Node) (string, bool) {
	if c, ok := n.(*Const); ok && c.Kind == String && !c.NA && !c.ph {
		return c.Str, true
	}
	return "", false
}

func (w *nameWalker) assignCall(c *Call, sc *nscope) {
	// assign(x, value, pos, envir, inherits, immediate): match by name, then
	// by position
	formals := []string{"x", "value", "pos", "envir", "inherits", "immediate"}
	got := map[string]Node{}
	used := map[string]bool{}
	for _, a := range c.Args {
		if a.HasName {
			for _, f := range formals {
				if f == a.Name {
					got[f], used[f] = a.Value, true
				}
			}
		}
	}
	next := 0
	for _, a := range c.Args {
		if a.HasName {
			continue
		}
		for next < len(formals) && used[formals[next]] {
			next++
		}
		if next < len(formals) {
			got[formals[next]], used[formals[next]] = a.Value, true
			next++
		}
	}
	for _, f := range formals[1:] {
		if used[f] {
			w.walk(got[f], sc)
		}
	}
	lit, ok := stringLiteral(got["x"])
	switch {
	case !ok:
		w.walk(got["x"], sc)
		w.dynamic("assign with a computed name")
	case used["pos"] || used["envir"]:
		w.dynamic("assign into another environment")
	default:
		w.bind(lit, sc)
	}
}

// assign handles the target of x <- v.
func (w *nameWalker) assign(target Node, sc *nscope) {
	switch t := target.(type) {
	case *Sym:
		w.bind(t.Name, sc)
	case *Const:
		if t.Kind == String && !t.NA && !t.ph {
			w.bind(t.Str, sc)
		}
	case *Call:
		// f(x) <- v, x$a <- v, x[i] <- v: x is read and assigned
		root, fnames := w.targetRoot(t, sc)
		for _, f := range fnames {
			if f == "environment" {
				w.dynamic("environment<-")
			}
		}
		if root != "" {
			w.use(root, sc, inFunction(sc))
			w.bind(root, sc)
		}
	}
}

// targetReads walks a <<- target without binding.
func (w *nameWalker) targetReads(target Node, sc *nscope) {
	if t, ok := target.(*Call); ok {
		w.targetRoot(t, sc)
	}
}

// targetRoot descends the first arguments of a replacement target to the
// variable being modified, walks the other arguments as reads, and returns
// the root name and the replacement functions met on the way.
func (w *nameWalker) targetRoot(t *Call, sc *nscope) (root string, funs []string) {
	var cur Node = t
	for {
		c, ok := cur.(*Call)
		if !ok || len(c.Args) == 0 {
			break
		}
		funs = append(funs, c.FunName())
		start := 1
		if c.FunName() == "$" || c.FunName() == "@" {
			start = len(c.Args)
		}
		w.walkArgs(c.Args[start:], sc)
		cur = c.Args[0].Value
	}
	switch r := cur.(type) {
	case *Sym:
		return r.Name, funs
	case *Const:
		if r.Kind == String && !r.NA && !r.ph {
			return r.Str, funs
		}
	}
	return "", funs
}
