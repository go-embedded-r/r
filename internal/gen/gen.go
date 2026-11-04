// Package gen produces the golden files of testdata/golden from the reference
// R, and reads them back for the tests that run without R.
//
// The inputs are written for this repository (testdata/snippets.R,
// testdata/values.R); the goldens are what R 4.6.2 answers for them. No R
// source, test or output file belongs in this repository.
package gen

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/go-embedded-r/r/internal/oracle"
	"github.com/go-embedded-r/r/internal/parity"
	"github.com/go-embedded-r/r/rds"
)

// Oracle is the part of the reference R the generators use.
type Oracle interface {
	Deparse(src string) (*oracle.Reply, error)
	Call(op string, args ...any) (*oracle.Reply, error)
}

// Split cuts an input file into its snippets at the lines holding only
// "#---". A chunk whose first line starts with "## " is a header and dropped.
func Split(text string) []string {
	var out []string
	for _, chunk := range strings.Split(text, "\n#---\n") {
		chunk = strings.Trim(chunk, "\n")
		if chunk == "" || strings.HasPrefix(chunk, "## ") {
			continue
		}
		out = append(out, chunk)
	}
	return out
}

// ReadSnippets reads and splits a file under root.
func ReadSnippets(root, rel string) ([]string, error) {
	b, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		return nil, err
	}
	s := Split(string(b))
	if len(s) == 0 {
		return nil, fmt.Errorf("gen: %s holds no snippet", rel)
	}
	return s, nil
}

// DeparseRecord is one line of deparse.jsonl.
type DeparseRecord struct {
	Src string `json:"src"`
	// Deparse holds, for each top-level expression of Src, the lines deparse()
	// returns. Empty when Src does not parse.
	Deparse [][]string `json:"deparse,omitempty"`
	// Error is R's parse error message, when Src does not parse.
	Error string `json:"error,omitempty"`
}

// RDSRecord is one line of rds.jsonl.
type RDSRecord struct {
	Src    string   `json:"src"`
	Hex    string   `json:"hex"` // serialize(value, NULL, xdr = TRUE, version = 3)
	Typeof string   `json:"typeof"`
	Length int      `json:"length"`
	Class  []string `json:"class"`
}

// DeparseGolden asks R to parse and deparse each snippet.
func DeparseGolden(o Oracle, snippets []string) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	for _, s := range snippets {
		rec := DeparseRecord{Src: s}
		r, err := o.Deparse(s)
		switch {
		case err == nil:
			l, ok := r.Value.(*rds.List)
			if !ok {
				return nil, fmt.Errorf("gen: deparse of %q returned %T", s, r.Value)
			}
			for _, e := range l.V {
				c, ok := e.(*rds.Character)
				if !ok {
					return nil, fmt.Errorf("gen: deparse of %q: element is %T", s, e)
				}
				lines := make([]string, len(c.V))
				for i, x := range c.V {
					lines[i] = x.S
				}
				rec.Deparse = append(rec.Deparse, lines)
			}
		case errors.Is(err, oracle.ErrR):
			rec.Error = r.Message
		default:
			return nil, err
		}
		if err := enc.Encode(rec); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// RDSGolden asks R to evaluate each expression and serialize the result.
func RDSGolden(o Oracle, exprs []string) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	for _, s := range exprs {
		r, err := o.Call("describe", s)
		if err != nil {
			return nil, fmt.Errorf("gen: %q: %w", s, err)
		}
		l, ok := r.Value.(*rds.List)
		if !ok || len(l.V) != 4 {
			return nil, fmt.Errorf("gen: describe of %q returned %T", s, r.Value)
		}
		raw, ok1 := l.V[0].(*rds.Raw)
		typ, ok2 := l.V[1].(*rds.Character)
		n, ok3 := l.V[2].(*rds.Integer)
		cls, ok4 := l.V[3].(*rds.Character)
		if !ok1 || !ok2 || !ok3 || !ok4 || len(typ.V) != 1 || len(n.V) != 1 {
			return nil, fmt.Errorf("gen: describe of %q has an unexpected shape", s)
		}
		rec := RDSRecord{Src: s, Hex: hex.EncodeToString(raw.V), Typeof: typ.V[0].S, Length: int(n.V[0])}
		for _, c := range cls.V {
			rec.Class = append(rec.Class, c.S)
		}
		if err := enc.Encode(rec); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// ParseRecord is one line of parse.jsonl: what R makes of a source text.
type ParseRecord struct {
	Src string `json:"src"`
	// Tree is the canonical dump of R's parse (see syntax.Dump), when it parses.
	Tree []string `json:"tree,omitempty"`
	// Tokens are the terminals of getParseData, "l1:c1-l2:c2 KIND hex(text)".
	Tokens []string `json:"tokens,omitempty"`
	// Deparse and Print are, for each top-level expression, the lines of
	// deparse() and of print() (see the oracle's deparse2 operation).
	Deparse [][]string `json:"deparse,omitempty"`
	Print   [][]string `json:"print,omitempty"`
	// Error is the message of parse(), when it does not parse.
	Error string `json:"error,omitempty"`
}

// ParseGolden asks R to parse each source text (texts that are not valid
// UTF-8 are left out: JSON cannot carry them).
func ParseGolden(o parity.Oracle, srcs []string) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	for _, s := range srcs {
		if !utf8.ValidString(s) {
			continue
		}
		rec := ParseRecord{Src: s}
		tree, rErr, err := parity.TreeLines(o, s)
		if err != nil {
			return nil, err
		}
		if rErr != "" {
			rec.Error = rErr
		} else {
			rec.Tree = tree
			r, err := o.Call("tokens", s)
			if err != nil {
				return nil, err
			}
			ch, _ := r.Value.(*rds.Character)
			rec.Tokens = []string{}
			if ch != nil {
				for _, t := range ch.V {
					rec.Tokens = append(rec.Tokens, t.S)
				}
			}
			for _, mode := range []string{"default", "print"} {
				dr, err := o.Call("deparse2", s, mode)
				if err != nil {
					return nil, err
				}
				l, _ := dr.Value.(*rds.List)
				var all [][]string
				if l != nil {
					for _, e := range l.V {
						lines := []string{}
						if c, ok := e.(*rds.Character); ok {
							for _, x := range c.V {
								lines = append(lines, x.S)
							}
						}
						all = append(all, lines)
					}
				}
				if mode == "default" {
					rec.Deparse = all
				} else {
					rec.Print = all
				}
			}
		}
		if err := enc.Encode(rec); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// ParseInputs are the sources of parse.jsonl: the snippet files and the
// extra cases of package parity.
func ParseInputs(root string) ([]string, error) {
	var all []string
	for _, f := range []string{"testdata/snippets.R", "testdata/errors.R", "testdata/edge.R", "testdata/deparse.R", "testdata/values.R"} {
		s, err := ReadSnippets(root, f)
		if err != nil {
			return nil, err
		}
		all = append(all, s...)
	}
	return append(all, parity.Extras()...), nil
}

// ReadJSONL decodes a golden file into records of type T.
func ReadJSONL[T any](path string) ([]T, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var out []T
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	for dec.More() {
		var t T
		if err := dec.Decode(&t); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		out = append(out, t)
	}
	return out, nil
}

// EvalRecord is one line of eval.jsonl: what R's evaluator does with a snippet.
type EvalRecord struct {
	Src      string   `json:"src"`
	OK       bool     `json:"ok"`
	Hex      string   `json:"hex,omitempty"` // serialize() of the last value
	Typeof   string   `json:"typeof,omitempty"`
	Visible  bool     `json:"visible,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
	Output   string   `json:"output,omitempty"` // what the cell printed on standard output
	Stderr   string   `json:"stderr,omitempty"` // message() output
	Error    string   `json:"error,omitempty"`  // conditionMessage
	Call     string   `json:"call,omitempty"`   // deparsed conditionCall
}

// EvalGolden asks R to evaluate each snippet as a cell.
func EvalGolden(o parity.Oracle, snippets []string) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	for _, s := range snippets {
		r, err := o.Call("evalcase", s)
		if err != nil {
			return nil, fmt.Errorf("gen: %q: %w", s, err)
		}
		l, ok := r.Value.(*rds.List)
		if !ok {
			return nil, fmt.Errorf("gen: evalcase of %q returned %T", s, r.Value)
		}
		names, _ := l.Attr("names").(*rds.Character)
		if names == nil || len(names.V) != len(l.V) {
			return nil, fmt.Errorf("gen: evalcase of %q returned an unnamed list", s)
		}
		f := map[string]rds.Value{}
		for i, n := range names.V {
			f[n.S] = l.V[i]
		}
		str := func(k string) string {
			if c, ok := f[k].(*rds.Character); ok && len(c.V) > 0 {
				return c.V[0].S
			}
			return ""
		}
		rec := EvalRecord{Src: s}
		if lg, ok := f["ok"].(*rds.Logical); ok && len(lg.V) == 1 && lg.V[0] == 1 {
			rec.OK = true
			if raw, ok := f["bytes"].(*rds.Raw); ok {
				v, _, err := rds.Decode(raw.V, rds.Limits{})
				if err != nil {
					return nil, fmt.Errorf("gen: value of %q: %w", s, err)
				}
				norm, err := rds.Encode(NormalizeNaN(v))
				if err != nil {
					return nil, err
				}
				rec.Hex = hex.EncodeToString(norm)
			}
			rec.Typeof = str("typeof")
			if v, ok := f["visible"].(*rds.Logical); ok && len(v.V) == 1 {
				rec.Visible = v.V[0] == 1
			}
		} else {
			rec.Error, rec.Call = str("message"), str("call")
		}
		rec.Output, rec.Stderr = str("output"), str("stderr")
		if w, ok := f["warnings"].(*rds.Character); ok {
			for _, x := range w.V {
				rec.Warnings = append(rec.Warnings, x.S)
			}
		}
		if err := enc.Encode(rec); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// EvalInputs reads the snippets under testdata/eval.
func EvalInputs(root string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(root, "testdata/eval"))
	if err != nil {
		return nil, err
	}
	var all []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".R") {
			s, err := ReadSnippets(root, "testdata/eval/"+e.Name())
			if err != nil {
				return nil, err
			}
			all = append(all, s...)
		}
	}
	return all, nil
}

// ParseJSONL decodes golden data already in memory.
func ParseJSONL[T any](b []byte) ([]T, error) {
	var out []T
	dec := json.NewDecoder(bytes.NewReader(b))
	for dec.More() {
		var t T
		if err := dec.Decode(&t); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

// NormalizeNaN rewrites, in place, the NaNs of every double vector of v to
// one NaN and one NA: the bits of a NaN (its sign, whether it is quiet) are
// whatever the hardware made of it (the sign bit of 0/0 differs between
// x86-64 and arm64), so goldens keep only whether an element is NA or NaN.
func NormalizeNaN(v rds.Value) rds.Value {
	switch x := v.(type) {
	case *rds.Real:
		for i, f := range x.V {
			if f != f {
				if math.Float64bits(f)&0xFFFFFFFF == 1954 {
					x.V[i] = math.Float64frombits(rds.NAReal)
				} else {
					x.V[i] = math.Float64frombits(0x7FF8000000000000)
				}
			}
		}
		normAttrs(x.Meta)
	case *rds.List:
		for _, e := range x.V {
			NormalizeNaN(e)
		}
		normAttrs(x.Meta)
	case *rds.Integer:
		normAttrs(x.Meta)
	case *rds.Logical:
		normAttrs(x.Meta)
	case *rds.Character:
		normAttrs(x.Meta)
	}
	return v
}

func normAttrs(m rds.Meta) {
	for _, a := range m.Attrs {
		NormalizeNaN(a.Value)
	}
}
