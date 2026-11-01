package rds_test

import (
	"bytes"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/go-embedded-r/r/internal/gen"
	"github.com/go-embedded-r/r/rds"
)

var typeofName = map[rds.Type]string{
	rds.NILSXP: "NULL", rds.SYMSXP: "symbol", rds.LISTSXP: "pairlist", rds.CLOSXP: "closure",
	rds.ENVSXP: "environment", rds.LANGSXP: "language", rds.SPECIALSXP: "special",
	rds.BUILTINSXP: "builtin", rds.LGLSXP: "logical", rds.INTSXP: "integer", rds.REALSXP: "double",
	rds.CPLXSXP: "complex", rds.STRSXP: "character", rds.VECSXP: "list", rds.EXPRSXP: "expression",
	rds.RAWSXP: "raw", rds.OBJSXP: "S4",
}

func length(v rds.Value) int {
	switch x := v.(type) {
	case *rds.Logical:
		return len(x.V)
	case *rds.Integer:
		return len(x.V)
	case *rds.Real:
		return len(x.V)
	case *rds.Complex:
		return len(x.V)
	case *rds.Character:
		return len(x.V)
	case *rds.Raw:
		return len(x.V)
	case *rds.List:
		return len(x.V)
	case *rds.Pairlist:
		return len(x.Items)
	}
	return -1
}

func goldens(t *testing.T) []gen.RDSRecord {
	t.Helper()
	recs, err := gen.ReadJSONL[gen.RDSRecord]("../testdata/golden/rds.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	if len(recs) < 40 {
		t.Fatalf("only %d golden records", len(recs))
	}
	return recs
}

// Every value R serialized decodes, agrees with R's own typeof/length, and
// encodes back to R's exact bytes: Go reads and writes what R does.
func TestGoldenRoundTrip(t *testing.T) {
	for _, rec := range goldens(t) {
		t.Run(rec.Src, func(t *testing.T) {
			raw, err := hex.DecodeString(rec.Hex)
			if err != nil {
				t.Fatal(err)
			}
			v, h, err := rds.Decode(raw, rds.Limits{})
			if err != nil {
				t.Fatal(err)
			}
			if h != rds.DefaultHeader {
				t.Errorf("header %+v, want %+v", h, rds.DefaultHeader)
			}
			if _, isAlt := v.(*rds.Altrep); isAlt {
				m, err := v.(*rds.Altrep).Materialize(1 << 20)
				if err != nil {
					t.Fatalf("materialize: %v", err)
				}
				if want := rec.Typeof; typeofName[m.Type()] != want || length(m) != rec.Length {
					t.Errorf("materialized %s of length %d, R says %s of length %d", typeofName[m.Type()], length(m), want, rec.Length)
				}
			} else if got := typeofName[v.Type()]; got != rec.Typeof {
				t.Errorf("typeof %q, R says %q", got, rec.Typeof)
			}
			if n := length(v); n >= 0 {
				if _, isAlt := v.(*rds.Altrep); !isAlt && n != rec.Length {
					t.Errorf("length %d, R says %d", n, rec.Length)
				}
			}
			out, err := rds.Encode(v)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(out, raw) {
				t.Errorf("re-encoded bytes differ:\n got %x\nwant %x", out, raw)
			}
		})
	}
}

func golden(t *testing.T, src string) []byte {
	t.Helper()
	for _, rec := range goldens(t) {
		if rec.Src == src {
			b, _ := hex.DecodeString(rec.Hex)
			return b
		}
	}
	t.Fatalf("no golden for %q", src)
	return nil
}

// Values built in Go encode to the bytes R writes for the same object.
func TestEncodeMatchesR(t *testing.T) {
	cases := []struct {
		src string
		v   rds.Value
	}{
		{"NULL", rds.Null},
		{"TRUE", &rds.Logical{V: []int32{1}}},
		{"c(TRUE, NA, FALSE)", &rds.Logical{V: []int32{1, rds.NAInteger, 0}}},
		{"42L", &rds.Integer{V: []int32{42}}},
		{"c(1L, NA, -3L)", &rds.Integer{V: []int32{1, rds.NAInteger, -3}}},
		{"1.5", &rds.Real{V: []float64{1.5}}},
		{`"a"`, rds.Strings("a")},
		{`c("héllo", NA, "", "日本語")`, &rds.Character{V: []rds.Str{rds.NewStr("héllo"), {NA: true}, rds.NewStr(""), rds.NewStr("日本語")}}},
		{"as.raw(c(0, 1, 255))", &rds.Raw{V: []byte{0, 1, 255}}},
		{"c(a = 1, b = 2)", &rds.Real{
			Meta: rds.Meta{Attrs: []rds.Attr{{Name: "names", Value: rds.Strings("a", "b")}}},
			V:    []float64{1, 2},
		}},
		{"globalenv()", rds.GlobalEnv},
		{"emptyenv()", rds.EmptyEnv},
		{"list(quote(x), quote(x), quote(y))", &rds.List{V: []rds.Value{
			&rds.Sym{Name: "x", Gp: 64}, &rds.Sym{Name: "x", Gp: 64}, &rds.Sym{Name: "y", Gp: 64},
		}}},
		{"quote(f(x, y = 1))", &rds.Pairlist{Kind: rds.LANGSXP, Items: []rds.Item{
			{Value: &rds.Sym{Name: "f", Gp: 64}},
			{Value: &rds.Sym{Name: "x", Gp: 64}},
			{Tag: &rds.Sym{Name: "y", Gp: 64}, Value: &rds.Real{V: []float64{1}}},
		}}},
	}
	for _, c := range cases {
		t.Run(c.src, func(t *testing.T) {
			got, err := rds.Encode(c.v)
			if err != nil {
				t.Fatal(err)
			}
			if want := golden(t, c.src); !bytes.Equal(got, want) {
				t.Errorf("got %x\nwant %x", got, want)
			}
		})
	}
}

// NA_real_ and NaN keep their payloads through a decode and an encode.
func TestNaNPayloads(t *testing.T) {
	raw := golden(t, "c(0.1, NA, NaN, Inf, -Inf, 1e300, -0)")
	v, _, err := rds.Decode(raw, rds.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	r := v.(*rds.Real)
	if got := f64bits(r.V[1]); got != rds.NAReal {
		t.Errorf("NA_real_ bits %#x, want %#x", got, rds.NAReal)
	}
	if got := f64bits(r.V[2]); got == rds.NAReal || got>>52&0x7ff != 0x7ff {
		t.Errorf("NaN bits %#x", got)
	}
	if got := f64bits(r.V[6]); got != 1<<63 {
		t.Errorf("-0 bits %#x", got)
	}
}

func TestAltrep(t *testing.T) {
	check := func(src string, want rds.Value) {
		t.Helper()
		v, _, err := rds.Decode(golden(t, src), rds.Limits{})
		if err != nil {
			t.Fatal(err)
		}
		a, ok := v.(*rds.Altrep)
		if !ok {
			t.Fatalf("%s: %T, want *Altrep", src, v)
		}
		got, err := a.Materialize(100)
		if err != nil {
			t.Fatal(err)
		}
		g, _ := rds.Encode(got)
		w, _ := rds.Encode(want)
		if !bytes.Equal(g, w) {
			t.Errorf("%s materialized to\n %x\nwant\n %x", src, g, w)
		}
	}
	check("1:10", &rds.Integer{V: []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}})
	check("3:1", &rds.Integer{V: []int32{3, 2, 1}})
	check("as.character(1:5)", &rds.Character{V: []rds.Str{rds.NewStr("1"), rds.NewStr("2"), rds.NewStr("3"), rds.NewStr("4"), rds.NewStr("5")}})

	a := mustAltrep(t, golden(t, "1:10"))
	if _, err := a.Materialize(5); err == nil {
		t.Error("a sequence longer than maxLen was materialized")
	}
	a.Class = &rds.Sym{Name: "mmap_integer"}
	if _, err := a.Materialize(100); !errors.Is(err, rds.ErrUnsupportedAltrep) {
		t.Errorf("unknown class: %v", err)
	}
}

func mustAltrep(t *testing.T, raw []byte) *rds.Altrep {
	t.Helper()
	v, _, err := rds.Decode(raw, rds.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	a, ok := v.(*rds.Altrep)
	if !ok {
		t.Fatalf("%T is not an *Altrep", v)
	}
	return a
}
