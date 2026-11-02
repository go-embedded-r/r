package syntax

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// ErrLimit is wrapped by the error returned when an input exceeds a
// configured limit (bytes, nesting depth or nodes).
var ErrLimit = errors.New("syntax: limit exceeded")

// Error is a parse failure. Its text is, byte for byte, the message R's
// parse() puts in the error condition it signals.
type Error struct {
	// Line and Col locate the error as R reports it (Col counts characters,
	// tab-expanded; 0 at the end of input).
	Line, Col int
	// Message is R's own description: "unexpected symbol", "unexpected end
	// of input", "'\q' is an unrecognized escape in character string", ...
	Message string
	// Lexical is set for errors R raises from its lexer or from the pipe
	// rewriting (their text carries "(<input>:line:col)" and no context).
	Lexical bool
	text    string
}

// Error returns the complete message, including the source context lines and
// the caret for a syntax error.
func (e *Error) Error() string { return e.text }

// LimitError reports a configured limit being exceeded; it wraps [ErrLimit].
type LimitError struct {
	What  string // "bytes", "depth" or "nodes"
	Limit int
}

func (e *LimitError) Error() string {
	return fmt.Sprintf("syntax: %s limit %d exceeded", e.What, e.Limit)
}

// Unwrap lets errors.Is(err, ErrLimit) work.
func (e *LimitError) Unwrap() error { return ErrLimit }

// syntaxError builds R's message for a syntax error detected at the token
// that starts at pos; consumed is the source up to the end of that token
// (what R's lexer had read), with a final newline counted as in R.
func syntaxError(filename, msg string, pos Pos, consumed []byte) *Error {
	ctx := consumed
	if len(ctx) > 256 {
		ctx = ctx[len(ctx)-256:]
	}
	// a multibyte character cut short by the end of the input was pushed
	// back by R's lexer, and does not appear in the context
	for len(ctx) > 0 {
		r, n := utf8.DecodeLastRune(ctx)
		if r != utf8.RuneError || n != 1 {
			break
		}
		ctx = ctx[:len(ctx)-1]
		consumed = consumed[:len(consumed)-1]
	}
	lines := strings.Split(string(ctx), "\n")
	// R_ParseContextLine is the line of the last character read; a trailing
	// empty line (the context ends with a newline) is dropped and the line
	// number decremented.
	ctxLine := 1
	for _, b := range consumed {
		if b == '\n' {
			ctxLine++
		}
	}
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
		ctxLine--
	}
	for i := range lines {
		lines[i] = tabExpand(lines[i])
	}
	var text string
	head := fmt.Sprintf("%s:%d:%d: %s", filename, pos.Line, pos.Col, msg)
	switch n := len(lines); {
	case n == 0:
		text = head
	case n == 1:
		width := len(fmt.Sprintf("%d: ", ctxLine))
		text = fmt.Sprintf("%s\n%d: %s\n%*s", head, ctxLine, lines[0], width+pos.Col+1, "^")
	default:
		width := len(fmt.Sprintf("%d:", ctxLine))
		text = fmt.Sprintf("%s\n%d: %s\n%d: %s\n%*s", head, ctxLine-1, lines[n-2], ctxLine, lines[n-1], width+pos.Col+1, "^")
	}
	return &Error{Line: pos.Line, Col: pos.Col, Message: msg, text: text}
}

// tabExpand replaces each tab by spaces up to the next multiple of 8, in
// bytes, as R does when it prints the context (the buffer is 200 bytes).
func tabExpand(s string) string {
	var b []byte
	for i := 0; i < len(s) && len(b) < 192; i++ {
		if s[i] == '\t' {
			for {
				b = append(b, ' ')
				if len(b)&7 == 0 {
					break
				}
			}
		} else {
			b = append(b, s[i])
		}
	}
	return string(b)
}

// lexicalError builds the message of an error raised by the lexer or by the
// pipe rewriting: R's text, a space, and "(file:line:col)".
func lexicalError(filename, format string, pos Pos, args ...any) *Error {
	msg := fmt.Sprintf(format, args...)
	return &Error{Line: pos.Line, Col: pos.Col, Message: msg, Lexical: true,
		text: fmt.Sprintf("%s (%s:%d:%d)", msg, filename, pos.Line, pos.Col)}
}
