package eval

import "github.com/go-embedded-r/r/rds"

// Env is an environment: a frame of bindings and an enclosure.
type Env struct {
	parent *Env
	small  []binding            // the frame of a call: a few bindings, searched linearly
	big    map[*Symbol]*binding // above smallMax bindings
	order  []*Symbol
	// missing records, for a function frame, the formals that were not
	// supplied (they may still hold a default).
	missing map[*Symbol]bool
	fn      *Closure // the function whose frame this is, nil otherwise
	locked  bool
	name    string
	// special is set once the frame holds a locked or active binding, or is
	// locked itself: only then does an assignment take the slow path.
	special bool
	// allLocked locks every binding of the frame at once (the base frame).
	allLocked bool
	// alias makes this environment another name for the frame of alias (the
	// base namespace shares the frame of the base environment).
	alias *Env
	// ns says the environment is a namespace.
	ns    bool
	attrs []Attr
}

type binding struct {
	sym    *Symbol
	val    Value
	locked bool
	active Value // the function of an active binding, nil for a regular one
}

const smallMax = 8

func (e *Env) Type() Type { return rds.ENVSXP }

// NewEnv makes an environment enclosed by parent.
func NewEnv(parent *Env) *Env { return &Env{parent: parent} }

// Parent is the enclosure.
func (e *Env) Parent() *Env { return e.parent }

func (e *Env) lookup(sym *Symbol) *binding {
	if e.alias != nil {
		return e.alias.lookup(sym)
	}
	if e.big != nil {
		return e.big[sym]
	}
	for i := range e.small {
		if e.small[i].sym == sym {
			return &e.small[i]
		}
	}
	return nil
}

// Get returns the value bound to name in this frame only.
func (e *Env) Get(name string) (Value, bool) {
	if b := e.lookup(Intern(name)); b != nil {
		return b.val, true
	}
	return nil, false
}

// set binds sym in this frame.
func (e *Env) set(sym *Symbol, v Value) {
	if e.alias != nil {
		e.alias.set(sym, v)
		return
	}
	if b := e.lookup(sym); b != nil {
		b.val = v
		delete(e.missing, sym) // assigning to a formal makes it not missing
		return
	}
	e.order = append(e.order, sym)
	if e.big != nil {
		e.big[sym] = &binding{sym: sym, val: v}
		return
	}
	if len(e.small) >= smallMax {
		e.big = make(map[*Symbol]*binding, 2*smallMax)
		for i := range e.small {
			b := e.small[i]
			e.big[b.sym] = &b
		}
		e.small = nil
		e.big[sym] = &binding{sym: sym, val: v}
		return
	}
	e.small = append(e.small, binding{sym: sym, val: v})
}

// Set binds name in this frame.
func (e *Env) Set(name string, v Value) { e.set(Intern(name), v) }

// Names lists the bound names in order of creation (ls(sorted = FALSE)).
func (e *Env) Names() []string {
	if e.alias != nil {
		return e.alias.Names()
	}
	out := make([]string, len(e.order))
	for i, s := range e.order {
		out[i] = s.Name
	}
	return out
}

// find looks sym up through the enclosures.
func (e *Env) find(sym *Symbol) (*binding, *Env) {
	for env := e; env != nil; env = env.parent {
		if b := env.lookup(sym); b != nil {
			return b, env
		}
	}
	return nil, nil
}
