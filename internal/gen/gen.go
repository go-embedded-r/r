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
	"os"
	"path/filepath"
	"strings"

	"github.com/go-embedded-r/r/internal/oracle"
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
	for _, chunk := range strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n#---\n") {
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
