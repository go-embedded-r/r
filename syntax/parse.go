package syntax

import (
	"strings"
)

// Limits bound what the parser accepts; a zero field selects its default.
// Exceeding one is a *[LimitError] (which wraps [ErrLimit]).
type Limits struct {
	// MaxBytes bounds the source size (default 64 MiB).
	MaxBytes int
	// MaxDepth bounds the nesting of the syntax tree (default 5000).
	MaxDepth int
	// MaxNodes bounds the number of tree nodes (default 8 million).
	MaxNodes int
}

const (
	defaultMaxBytes = 64 << 20
	defaultMaxDepth = 5000
	defaultMaxNodes = 8 << 20
)

func (l Limits) withDefaults() Limits {
	if l.MaxBytes <= 0 {
		l.MaxBytes = defaultMaxBytes
	}
	if l.MaxDepth <= 0 {
		l.MaxDepth = defaultMaxDepth
	}
	if l.MaxNodes <= 0 {
		l.MaxNodes = defaultMaxNodes
	}
	return l
}

// Options configure a parse.
type Options struct {
	Limits
	// Filename is the name R shows in syntax errors (default "<text>", what
	// parse(text = ) uses).
	Filename string
	// KeepSource makes lexical errors use Filename too; without it they say
	// "<input>", as parse(keep.source = FALSE) does.
	KeepSource bool
}

// File is the result of a parse: the top-level expressions, the terminal
// tokens, and the source they came from.
type File struct {
	// Source is the text that was parsed: the input and the newline R adds.
	Source []byte
	// Exprs are the top-level expressions, in order.
	Exprs []Node
	// Tokens are the terminals getParseData() lists, in source order,
	// comments included.
	Tokens []Token
}

// ParseAll parses R source text as parse(text = src) would (a newline is
// appended to it), with the default limits. On a syntax error it returns the expressions before the
// error together with the *[Error].
func ParseAll(src string) (*File, error) { return ParseAllLimits(src, Limits{}) }

// ParseAllLimits is ParseAll with explicit limits.
func ParseAllLimits(src string, lim Limits) (*File, error) {
	return Parse(src, Options{Limits: lim})
}

// Parse parses src with the given options.
func Parse(src string, opt Options) (*File, error) {
	lim := opt.Limits.withDefaults()
	if len(src) > lim.MaxBytes {
		return nil, &LimitError{What: "bytes", Limit: lim.MaxBytes}
	}
	// like parse(text = ), the text is terminated by a newline
	b := append([]byte(src), '\n')
	if opt.Filename == "" {
		opt.Filename = "<text>"
	}
	lexFile := "<input>"
	if opt.KeepSource {
		lexFile = opt.Filename
	}
	p := &parser{lx: newLexer(b, lexFile), file: opt.Filename, lim: lim}
	f := &File{Source: b}
	err := p.run(f)
	f.Tokens = p.tokens()
	return f, err
}

type parser struct {
	lx   *lexer
	tok  lexTok
	file string
	lim  Limits

	depth int
	nodes int

	prevLast  Pos
	prevLastB int
}

func (p *parser) run(f *File) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch e := r.(type) {
			case *Error:
				err = e
			case *LimitError:
				err = e
			default:
				panic(r)
			}
		}
	}()
	p.advance()
	for {
		switch p.tok.kind {
		case tEOF:
			return nil
		case '\n':
			p.advance()
			continue
		}
		e := p.exprOrAssign(0)
		term := p.tok
		switch term.kind {
		case '\n', ';':
		default:
			p.unexpected()
		}
		if p.lx.havePH && hasPlaceholder(e) {
			pos := term.last
			if term.kind == '\n' {
				pos.Line--
			}
			panic(lexicalError(p.lx.file, "invalid use of pipe placeholder", pos))
		}
		f.Exprs = append(f.Exprs, e)
		p.advance()
	}
}

func (p *parser) tokens() []Token {
	out := make([]Token, len(p.lx.toks))
	for i, t := range p.lx.toks {
		out[i] = Token{Kind: t.kind, Text: t.text, Span: Span{Start: t.first, End: t.last, StartByte: t.fb, EndByte: t.lb}}
	}
	return out
}

// advance consumes the current token and reads the next.
func (p *parser) advance() {
	p.prevLast = p.tok.last
	p.prevLastB = p.tok.lb
	p.tok = p.lx.next()
}

func (p *parser) unexpected() {
	t := p.tok
	end := t.cend
	if end > len(p.lx.src) {
		end = len(p.lx.src)
	}
	panic(syntaxError(p.file, "unexpected "+t.kind.unexpected(), t.first, p.lx.src[:end]))
}

func (p *parser) expect(k Kind) {
	if p.tok.kind != k {
		p.unexpected()
	}
	p.advance()
}

func (p *parser) mark(t lexTok, name string) {
	if t.rec >= 0 {
		p.lx.toks[t.rec].kind = name
	}
}

func (p *parser) enter() {
	p.depth++
	if p.depth > p.lim.MaxDepth {
		panic(&LimitError{What: "depth", Limit: p.lim.MaxDepth})
	}
}

func (p *parser) leave() { p.depth-- }

func (p *parser) newNode() {
	p.nodes++
	if p.nodes > p.lim.MaxNodes {
		panic(&LimitError{What: "nodes", Limit: p.lim.MaxNodes})
	}
}

func (p *parser) spanFrom(t lexTok) Span {
	return Span{Start: t.first, End: p.prevLast, StartByte: t.fb, EndByte: p.prevLastB}
}

func (p *parser) spanOf(a, b Span) Span {
	return Span{Start: a.Start, End: b.End, StartByte: a.StartByte, EndByte: b.EndByte}
}

func (p *parser) sym(name string, s Span) *Sym {
	p.newNode()
	return &Sym{Pos: s, Name: name}
}

func (p *parser) call(fun string, funSpan Span, whole Span, args ...Node) *Call {
	p.newNode()
	c := &Call{Pos: whole, Fun: p.sym(fun, funSpan)}
	for _, a := range args {
		c.Args = append(c.Args, Arg{Value: a, Pos: spanOfNode(a)})
	}
	return c
}

func spanOfNode(n Node) Span {
	if n == nil {
		return Span{}
	}
	return n.Span()
}

func tokSpan(t lexTok) Span {
	return Span{Start: t.first, End: t.last, StartByte: t.fb, EndByte: t.lb}
}

// ---- expressions ----------------------------------------------------------

// Binding powers, low to high, following R's precedence table.
const (
	bpQuestion = 1
	bpEq       = 3
	bpLeftAsg  = 5
	bpRightAsg = 7
	bpTilde    = 9
	bpOr       = 11
	bpAnd      = 13
	bpNot      = 15
	bpCmp      = 17
	bpAdd      = 19
	bpMul      = 21
	bpSpecial  = 23
	bpColon    = 27
	bpUnary    = 29
	bpPow      = 31
	bpDollar   = 33
)

// exprOrAssign parses expr_or_assign_or_help: an expression that may use '='
// and '?'.
func (p *parser) exprOrAssign(minbp int) Node { return p.expr(minbp, true) }

// exprOrHelp parses expr_or_help: '?' allowed, '=' not.
func (p *parser) exprOrHelp() Node { return p.expr(0, false) }

func isCmp(k Kind) bool {
	switch k {
	case tLT, tLE, tGT, tGE, tEQ, tNE:
		return true
	}
	return false
}

// expr is the precedence-climbing core. eq says whether '=' may be used at
// this level (it is never accepted inside the operand of another operator).
func (p *parser) expr(minbp int, eq bool) Node {
	p.enter()
	defer p.leave()
	start := p.tok
	lhs := p.prefix()
	return p.infix(lhs, start, minbp, eq)
}

// infix continues an expression whose first operand lhs (starting at token
// start) has been parsed.
func (p *parser) infix(lhs Node, start lexTok, minbp int, eq bool) Node {
	lastCmp := false
	for {
		t := p.tok
		k := t.kind
		var lbp, rbp int
		wasCmp := lastCmp
		lastCmp = false
		switch {
		case k == '(' || k == '[' || k == tLBB:
			lhs = p.postfix(lhs, start)
			continue
		case k == '$' || k == '@':
			if bpDollar < minbp {
				return lhs
			}
			lhs = p.dollar(lhs, start)
			continue
		case k == '?':
			lbp, rbp = bpQuestion, bpQuestion+1
		case k == tEQASSIGN:
			if !eq {
				return lhs
			}
			lbp, rbp = bpEq, bpEq
		case k == tLEFTASSIGN:
			lbp, rbp = bpLeftAsg, bpLeftAsg
		case k == tRIGHTASSIGN:
			lbp, rbp = bpRightAsg, bpRightAsg+1
		case k == '~':
			lbp, rbp = bpTilde, bpTilde+1
		case k == tOR || k == tOR2:
			lbp, rbp = bpOr, bpOr+1
		case k == tAND || k == tAND2:
			lbp, rbp = bpAnd, bpAnd+1
		case isCmp(k):
			lbp, rbp = bpCmp, bpCmp+1
		case k == '+' || k == '-':
			lbp, rbp = bpAdd, bpAdd+1
		case k == '*' || k == '/':
			lbp, rbp = bpMul, bpMul+1
		case k == tSPECIAL || k == tPIPE:
			lbp, rbp = bpSpecial, bpSpecial+1
		case k == tPIPEBIND:
			lbp, rbp = bpSpecial+2, bpSpecial+3
		case k == ':':
			lbp, rbp = bpColon, bpColon+1
		case k == '^':
			lbp, rbp = bpPow, bpPow
		default:
			return lhs
		}
		if lbp < minbp {
			return lhs
		}
		if isCmp(k) && wasCmp {
			p.unexpected()
		}
		p.advance()
		opSpan := tokSpan(t)
		if k == tPIPEBIND {
			panic(lexicalError(p.lx.file, "'=>' is disabled; set '_R_USE_PIPEBIND_' envvar to a true value to enable it", t.first))
		}
		rstart := p.tok
		var rhs Node
		switch k {
		case '?':
			rhs = p.expr(rbp, eq)
		case tEQASSIGN:
			rhs = p.expr(rbp, true)
		default:
			rhs = p.expr(rbp, false)
		}
		whole := Span{Start: start.first, End: p.prevLast, StartByte: start.fb, EndByte: p.prevLastB}
		switch k {
		case tPIPE:
			lhs = p.pipe(lhs, rhs, rstart, whole)
		case tRIGHTASSIGN:
			lhs = p.call(t.sym, opSpan, whole, rhs, lhs)
		default:
			lhs = p.call(t.sym, opSpan, whole, lhs, rhs)
		}
		if isCmp(k) {
			lastCmp = true
		}
	}
}

// prefix parses an operand: a constant, a name, a bracketed or braced
// expression, a unary operator, or a function/control-flow construct.
func (p *parser) prefix() Node {
	t := p.tok
	switch t.kind {
	case tNUM, tNULL:
		p.advance()
		return p.constant(t)
	case tSTR:
		p.advance()
		if p.tok.kind == tNSGET || p.tok.kind == tNSGETINT {
			return p.namespace(t)
		}
		return p.constant(t)
	case tSYMBOL:
		p.advance()
		if p.tok.kind == tNSGET || p.tok.kind == tNSGETINT {
			return p.namespace(t)
		}
		p.newNode()
		return &Sym{Pos: tokSpan(t), Name: t.sym}
	case tPLACEHOLDER:
		p.advance()
		p.lx.havePH = true
		p.newNode()
		return &Const{Pos: tokSpan(t), Kind: String, Str: "_", ph: true}
	case '{':
		return p.brace()
	case '(':
		p.advance()
		e := p.exprOrAssign(0)
		if p.tok.kind != ')' {
			p.unexpected()
		}
		p.advance()
		return p.call("(", tokSpan(t), p.spanFrom(t), e)
	case '-', '+':
		p.advance()
		e := p.expr(bpUnary+1, false)
		return p.call(t.sym, tokSpan(t), p.spanFrom(t), e)
	case '!':
		p.advance()
		e := p.expr(bpNot+1, false)
		return p.call("!", tokSpan(t), p.spanFrom(t), e)
	case '~':
		p.advance()
		e := p.expr(bpTilde+1, false)
		return p.call("~", tokSpan(t), p.spanFrom(t), e)
	case '?':
		p.advance()
		e := p.expr(bpQuestion+1, true)
		return p.call("?", tokSpan(t), p.spanFrom(t), e)
	case tFUNCTION, '\\':
		return p.function()
	case tIF:
		return p.ifExpr()
	case tFOR:
		return p.forExpr()
	case tWHILE:
		p.advance()
		cond := p.cond()
		body := p.exprOrAssign(0)
		return p.call("while", tokSpan(t), p.spanFrom(t), cond, body)
	case tREPEAT:
		p.advance()
		body := p.exprOrAssign(0)
		return p.call("repeat", tokSpan(t), p.spanFrom(t), body)
	case tNEXT, tBREAK:
		p.advance()
		name := "next"
		if t.kind == tBREAK {
			name = "break"
		}
		return p.call(name, tokSpan(t), p.spanFrom(t))
	}
	p.unexpected()
	return nil
}

func (p *parser) constant(t lexTok) Node {
	p.newNode()
	c := *t.cst
	c.Pos = tokSpan(t)
	return &c
}

// namespace parses pkg::name and pkg:::name; the left operand token t has
// been consumed and the operator is the current token.
func (p *parser) namespace(t lexTok) Node {
	op := p.tok
	p.advance()
	r := p.tok
	if r.kind != tSYMBOL && r.kind != tSTR {
		p.unexpected()
	}
	p.advance()
	if t.kind == tSYMBOL {
		p.mark(t, "SYMBOL_PACKAGE")
	}
	var l, rn Node
	if t.kind == tSYMBOL {
		p.newNode()
		l = &Sym{Pos: tokSpan(t), Name: t.sym}
	} else {
		l = p.constant(t)
	}
	if r.kind == tSYMBOL {
		p.newNode()
		rn = &Sym{Pos: tokSpan(r), Name: r.sym}
	} else {
		rn = p.constant(r)
	}
	return p.call(op.sym, tokSpan(op), p.spanFrom(t), l, rn)
}

func (p *parser) brace() Node {
	t := p.tok
	p.advance()
	var items []Node
	for {
		switch p.tok.kind {
		case ';', '\n':
			p.advance()
			continue
		case '}':
			p.lx.eatLines = false
			p.advance()
			c := p.call("{", tokSpan(t), p.spanFrom(t), items...)
			return c
		}
		e := p.exprOrAssign(0)
		items = append(items, e)
		switch p.tok.kind {
		case ';', '\n', '}':
		default:
			p.unexpected()
		}
	}
}

// cond parses '(' expr_or_help ')' and arranges for the newlines that follow
// to be ignored.
func (p *parser) cond() Node {
	if p.tok.kind != '(' {
		p.unexpected()
	}
	p.advance()
	e := p.exprOrHelp()
	if p.tok.kind != ')' {
		p.unexpected()
	}
	p.lx.eatLines = true
	p.advance()
	return e
}

func (p *parser) ifExpr() Node {
	t := p.tok
	p.advance()
	cond := p.cond()
	yes := p.exprOrAssign(0)
	if p.tok.kind == tELSE {
		p.advance()
		no := p.exprOrAssign(0)
		return p.call("if", tokSpan(t), p.spanFrom(t), cond, yes, no)
	}
	return p.call("if", tokSpan(t), p.spanFrom(t), cond, yes)
}

func (p *parser) forExpr() Node {
	t := p.tok
	p.advance()
	if p.tok.kind != '(' {
		p.unexpected()
	}
	p.advance()
	v := p.tok
	if v.kind != tSYMBOL {
		p.unexpected()
	}
	p.advance()
	p.expect(tIN)
	seq := p.exprOrHelp()
	if p.tok.kind != ')' {
		p.unexpected()
	}
	p.lx.eatLines = true
	p.advance()
	body := p.exprOrAssign(0)
	p.newNode()
	vs := &Sym{Pos: tokSpan(v), Name: v.sym}
	return p.call("for", tokSpan(t), p.spanFrom(t), vs, seq, body)
}

func (p *parser) function() Node {
	t := p.tok
	lambda := t.kind == '\\'
	p.advance()
	if p.tok.kind != '(' {
		p.unexpected()
	}
	p.advance()
	var params []Param
	badValue := p.tok.kind == ','
	for p.tok.kind != ')' {
		if badValue {
			// R accepts a leading comma in the grammar and fails when it
			// builds the formals
			p.advance()
		}
		s := p.tok
		if s.kind != tSYMBOL {
			p.unexpected()
		}
		p.advance()
		p.mark(s, "SYMBOL_FORMALS")
		for _, q := range params {
			if q.Name == s.sym {
				panic(lexicalError(p.lx.file, "repeated formal argument '%s'", s.first, s.sym))
			}
		}
		prm := Param{Name: s.sym, Pos: tokSpan(s)}
		if p.tok.kind == tEQASSIGN {
			p.mark(p.tok, "EQ_FORMALS")
			p.advance()
			prm.Default = p.exprOrHelp()
			prm.Pos = p.spanFrom(s)
		}
		if badValue {
			panic(&Error{Message: "bad value", Lexical: true, text: "bad value"})
		}
		params = append(params, prm)
		if p.tok.kind == ',' {
			p.advance()
			if p.tok.kind == ')' {
				p.unexpected()
			}
			continue
		}
		if p.tok.kind != ')' {
			p.unexpected()
		}
	}
	p.lx.eatLines = true
	p.advance()
	body := p.exprOrAssign(0)
	p.newNode()
	return &Function{Pos: p.spanFrom(t), Params: params, Body: body, Lambda: lambda}
}

// postfix parses a call or a subscript on lhs (starting at token start).
func (p *parser) postfix(lhs Node, start lexTok) Node {
	t := p.tok
	p.advance()
	switch t.kind {
	case '(':
		args := p.sublist(')')
		p.expect(')')
		// the callee's last SYMBOL token becomes SYMBOL_FUNCTION_CALL
		p.markCallee(lhs)
		if len(args) == 1 && args[0].Value == nil && !args[0].HasName {
			args = nil
		}
		fun := lhs
		if c, ok := lhs.(*Const); ok && c.Kind == String && !c.NA && !c.ph {
			p.newNode()
			fun = &Sym{Pos: c.Pos, Name: c.Str}
		}
		p.newNode()
		return &Call{Pos: p.spanFrom(start), Fun: fun, Args: args}
	case '[':
		args := p.sublist(']')
		p.expect(']')
		return p.subscript("[", t, lhs, args, start)
	default: // [[
		args := p.sublist(']')
		p.expect(']')
		p.expect(']')
		return p.subscript("[[", t, lhs, args, start)
	}
}

func (p *parser) subscript(name string, t lexTok, lhs Node, args []Arg, start lexTok) Node {
	p.newNode()
	c := &Call{Pos: p.spanFrom(start), Fun: p.sym(name, tokSpan(t))}
	c.Args = append(c.Args, Arg{Value: lhs, Pos: spanOfNode(lhs)})
	c.Args = append(c.Args, args...)
	return c
}

// markCallee turns the last SYMBOL token of the callee expression into
// SYMBOL_FUNCTION_CALL: the callee itself, or the name after $ or ::.
func (p *parser) markCallee(callee Node) {
	var sp Span
	switch c := callee.(type) {
	case *Sym:
		sp = c.Pos
	case *Call:
		switch c.FunName() {
		case "$", "::", ":::":
			if len(c.Args) == 2 {
				if s, ok := c.Args[1].Value.(*Sym); ok {
					sp = s.Pos
				}
			}
		}
	}
	if sp.EndByte == 0 {
		return
	}
	// find the recorded token at that byte offset, searching backwards
	toks := p.lx.toks
	for i := len(toks) - 1; i >= 0; i-- {
		if toks[i].fb == sp.StartByte {
			if toks[i].kind == "SYMBOL" {
				toks[i].kind = "SYMBOL_FUNCTION_CALL"
			}
			return
		}
		if toks[i].fb < sp.StartByte {
			return
		}
	}
}

func (p *parser) dollar(lhs Node, start lexTok) Node {
	op := p.tok
	p.advance()
	r := p.tok
	if r.kind != tSYMBOL && r.kind != tSTR {
		p.unexpected()
	}
	p.advance()
	var rn Node
	if r.kind == tSYMBOL {
		if op.kind == '@' {
			p.mark(r, "SLOT")
		}
		p.newNode()
		rn = &Sym{Pos: tokSpan(r), Name: r.sym}
	} else {
		rn = p.constant(r)
	}
	return p.call(op.sym, tokSpan(op), p.spanFrom(start), lhs, rn)
}

// sublist parses the arguments of a call or subscript up to (not including)
// the closing token.
func (p *parser) sublist(closer Kind) []Arg {
	var args []Arg
	for {
		args = append(args, p.sub(closer))
		if p.tok.kind == ',' {
			p.lx.eatLines = true
			p.advance()
			continue
		}
		return args
	}
}

func (p *parser) sub(closer Kind) Arg {
	t := p.tok
	if t.kind == ',' || t.kind == closer {
		return Arg{}
	}
	switch t.kind {
	case tSYMBOL, tSTR, tNULL:
		p.advance()
		if p.tok.kind == tEQASSIGN {
			eq := p.tok
			p.advance()
			name := t.sym
			if t.kind == tSTR {
				name = t.cst.Str
			}
			if t.kind == tNULL {
				name = "NULL"
			}
			p.mark(eq, "EQ_SUB")
			if t.kind == tSYMBOL {
				p.mark(t, "SYMBOL_SUB")
			}
			a := Arg{Name: name, HasName: true}
			if p.tok.kind == ',' || p.tok.kind == closer {
				a.Pos = p.spanFrom(t)
				return a
			}
			a.Value = p.exprOrHelp()
			a.Pos = p.spanFrom(t)
			return a
		}
		// not a name: t is the first operand of an expression
		var first Node
		switch t.kind {
		case tSYMBOL:
			if p.tok.kind == tNSGET || p.tok.kind == tNSGETINT {
				first = p.namespace(t)
			} else {
				p.newNode()
				first = &Sym{Pos: tokSpan(t), Name: t.sym}
			}
		case tSTR:
			if p.tok.kind == tNSGET || p.tok.kind == tNSGETINT {
				first = p.namespace(t)
			} else {
				first = p.constant(t)
			}
		default:
			first = p.constant(t)
		}
		p.enter()
		v := p.infix(first, t, 0, false)
		p.leave()
		return Arg{Value: v, Pos: spanOfNode(v)}
	}
	v := p.exprOrHelp()
	return Arg{Value: v, Pos: spanOfNode(v)}
}

// ---- the native pipe ------------------------------------------------------

var specialSymbols = map[string]bool{}

func init() {
	for _, s := range strings.Fields("if while repeat for break next return function ( { + - * / ^ %% %/% %*% : == != < > <= >= & | && || ! <- <<- = $ [ [[ $<- [<- [[<-") {
		specialSymbols[s] = true
	}
}

func asCall(n Node) (fun Node, args []Arg, ok bool) {
	switch c := n.(type) {
	case *Call:
		return c.Fun, c.Args, true
	case *Function:
		return &Sym{Name: "function"}, nil, true
	}
	return nil, nil, false
}

func (p *parser) pipe(lhs, rhs Node, rstart lexTok, whole Span) Node {
	fail := func(format string, a ...any) {
		panic(lexicalError(p.lx.file, format, rstart.first, a...))
	}
	call, isCall := rhs.(*Call)
	if !isCall {
		if _, isFn := rhs.(*Function); isFn {
			fail("function '%s' not supported in RHS call of a pipe", "function")
		}
		fail("The pipe operator requires a function call as RHS")
	}
	if hasPlaceholder(call.Fun) {
		fail("pipe placeholder cannot be used in the RHS function")
	}
	// _$a[1]$b and the like
	if slot, found := p.extractorChain(call, call, rstart); found {
		slot.Value = lhs
		slot.Pos = spanOfNode(lhs)
		return rhs
	}
	for i := range call.Args {
		if isPlaceholder(call.Args[i].Value) {
			if !call.Args[i].HasName {
				fail("pipe placeholder can only be used as a named argument")
			}
			for _, r := range call.Args[i+1:] {
				if isPlaceholder(r.Value) {
					fail("pipe placeholder may only appear once")
				}
			}
			call.Args[i].Value = lhs
			call.Pos = whole
			return rhs
		}
	}
	if s, ok := call.Fun.(*Sym); ok && specialSymbols[s.Name] {
		fail("function '%s' not supported in RHS call of a pipe", s.Name)
	}
	args := append([]Arg{{Value: lhs, Pos: spanOfNode(lhs)}}, call.Args...)
	call.Args = args
	call.Pos = whole
	return call
}

func isPlaceholder(n Node) bool {
	c, ok := n.(*Const)
	return ok && c.ph
}

// hasPlaceholder reports whether the placeholder occurs in the expression,
// descending into calls (and the body of a function) as R does.
func hasPlaceholder(n Node) bool {
	switch x := n.(type) {
	case *Const:
		return x.ph
	case *Call:
		if hasPlaceholder(x.Fun) {
			return true
		}
		for _, a := range x.Args {
			if a.Value != nil && hasPlaceholder(a.Value) {
				return true
			}
		}
	case *Function:
		return hasPlaceholder(x.Body)
	}
	return false
}

// extractorChain looks for the placeholder as the head of a chain of
// extractions ($, @, [, [[) in expr and returns the argument slot holding it.
func (p *parser) extractorChain(rhs, expr *Call, rstart lexTok) (*Arg, bool) {
	switch expr.FunName() {
	case "[", "[[", "$", "@":
	default:
		return nil, false
	}
	if len(expr.Args) == 0 {
		return nil, false
	}
	var slot *Arg
	var found bool
	first := expr.Args[0].Value
	if isPlaceholder(first) {
		slot, found = &expr.Args[0], true
	} else if c, ok := first.(*Call); ok {
		slot, found = p.extractorChain(rhs, c, rstart)
	}
	if found {
		for _, a := range expr.Args[1:] {
			if a.Value != nil && hasPlaceholder(a.Value) {
				panic(lexicalError(p.lx.file, "pipe placeholder may only appear once", rstart.first))
			}
		}
	}
	return slot, found
}
