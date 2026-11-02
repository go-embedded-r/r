// Package syntax is a lexer and parser for R 4.6.2, written from the R
// Language Definition and checked against the reference R (see cmd/rparity):
// it produces R's own language objects as an AST with a source map, the
// terminal tokens that utils::getParseData() reports, and R's parse errors
// word for word.
//
// The grammar and the context-sensitive handling of newlines follow R's
// documented behaviour; nothing is translated from R's sources.
package syntax

// Kind is a token kind. The single-character tokens use their character
// code; the others are the constants below.
type Kind int

// Token kinds with a name (the names of R's grammar).
const (
	tEOF Kind = 256 + iota
	tERROR
	tSTR
	tNUM
	tNULL
	tSYMBOL
	tFUNCTION
	tINCOMPLETE
	tLEFTASSIGN
	tEQASSIGN
	tRIGHTASSIGN
	tLBB
	tFOR
	tIN
	tIF
	tELSE
	tWHILE
	tNEXT
	tBREAK
	tREPEAT
	tGT
	tGE
	tLT
	tLE
	tEQ
	tNE
	tAND
	tOR
	tAND2
	tOR2
	tNSGET
	tNSGETINT
	tCOMMENT
	tSPECIAL
	tPIPE
	tPIPEBIND
	tPLACEHOLDER
	tUNDEF // a character that is no token
)

// Token names as getParseData() reports them.
var kindNames = map[Kind]string{
	tEOF: "END_OF_INPUT", tERROR: "ERROR", tSTR: "STR_CONST", tNUM: "NUM_CONST",
	tNULL: "NULL_CONST", tSYMBOL: "SYMBOL", tFUNCTION: "FUNCTION",
	tINCOMPLETE: "INCOMPLETE_STRING", tLEFTASSIGN: "LEFT_ASSIGN", tEQASSIGN: "EQ_ASSIGN",
	tRIGHTASSIGN: "RIGHT_ASSIGN", tLBB: "LBB", tFOR: "FOR", tIN: "IN", tIF: "IF",
	tELSE: "ELSE", tWHILE: "WHILE", tNEXT: "NEXT", tBREAK: "BREAK", tREPEAT: "REPEAT",
	tGT: "GT", tGE: "GE", tLT: "LT", tLE: "LE", tEQ: "EQ", tNE: "NE", tAND: "AND",
	tOR: "OR", tAND2: "AND2", tOR2: "OR2", tNSGET: "NS_GET", tNSGETINT: "NS_GET_INT",
	tCOMMENT: "COMMENT", tSPECIAL: "SPECIAL", tPIPE: "PIPE", tPIPEBIND: "PIPEBIND",
	tPLACEHOLDER: "PLACEHOLDER",
}

// name is the token name as getParseData() prints it: the grammar name, or
// the quoted character for a single-character token.
func (k Kind) name() string {
	if n, ok := kindNames[k]; ok {
		return n
	}
	switch k {
	case '\\':
		return `'\\'`
	case '\n':
		return `'\n'`
	}
	return "'" + string(rune(k)) + "'"
}

// unexpected is the description R puts after "unexpected " in a syntax error.
func (k Kind) unexpected() string {
	switch k {
	case tERROR, tPLACEHOLDER:
		return "input"
	case tUNDEF:
		return "invalid token"
	case tEOF:
		return "end of input"
	case tSTR:
		return "string constant"
	case tNUM:
		return "numeric constant"
	case tSYMBOL:
		return "symbol"
	case tLEFTASSIGN:
		return "assignment"
	case '\n':
		return "end of line"
	case tNULL:
		return "'NULL'"
	case tFUNCTION:
		return "'function'"
	case tEQASSIGN:
		return "'='"
	case tRIGHTASSIGN:
		return "'->'"
	case tLBB:
		return "'[['"
	case tFOR:
		return "'for'"
	case tIN:
		return "'in'"
	case tIF:
		return "'if'"
	case tELSE:
		return "'else'"
	case tWHILE:
		return "'while'"
	case tNEXT:
		return "'next'"
	case tBREAK:
		return "'break'"
	case tREPEAT:
		return "'repeat'"
	case tGT:
		return "'>'"
	case tGE:
		return "'>='"
	case tLT:
		return "'<'"
	case tLE:
		return "'<='"
	case tEQ:
		return "'=='"
	case tNE:
		return "'!='"
	case tAND:
		return "'&'"
	case tOR:
		return "'|'"
	case tAND2:
		return "'&&'"
	case tOR2:
		return "'||'"
	case tNSGET:
		return "'::'"
	case tNSGETINT:
		return "':::'"
	case tPIPE:
		return "'|>'"
	case tPIPEBIND:
		return "'=>'"
	case '\\':
		return `'\\'`
	}
	if n, ok := kindNames[k]; ok {
		return n
	}
	if k >= 256 {
		return "input"
	}
	if k < 0x20 || k >= 0x7f {
		return "input"
	}
	return "'" + string(rune(k)) + "'"
}

// Pos is a position as R counts it: the line (from 1) and the column (from
// 1; 0 for the position of a line end) in characters, with a tab advancing
// to the next multiple of 8.
type Pos struct{ Line, Col int }

// Span is a source range. Start and End are the positions of the first and
// the last character (End is inclusive, as getParseData's col2); StartByte
// and EndByte are byte offsets into the source, End exclusive.
type Span struct {
	Start, End         Pos
	StartByte, EndByte int
}

// Token is one terminal of the parse, as utils::getParseData() lists it.
type Token struct {
	// Kind is getParseData's token name: SYMBOL, SYMBOL_FUNCTION_CALL,
	// NUM_CONST, '+', LEFT_ASSIGN, COMMENT, ...
	Kind string
	Text string
	Span
}
