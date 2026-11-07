package rds_test

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"
	"testing"

	"github.com/go-embedded-r/r/rds"
)

// stream builds a version-3 XDR stream from 32-bit words and raw bytes.
func stream(parts ...any) []byte {
	b := []byte("X\n")
	b = binary.BigEndian.AppendUint32(b, 3)
	b = binary.BigEndian.AppendUint32(b, 0x00040602)
	b = binary.BigEndian.AppendUint32(b, 0x00030500)
	b = binary.BigEndian.AppendUint32(b, 5)
	b = append(b, "UTF-8"...)
	for _, p := range parts {
		switch x := p.(type) {
		case uint32:
			b = binary.BigEndian.AppendUint32(b, x)
		case int:
			b = binary.BigEndian.AppendUint32(b, uint32(x))
		case string:
			b = append(b, x...)
		}
	}
	return b
}

func TestDecodeRefusals(t *testing.T) {
	const nilv = 254
	cases := []struct {
		name string
		in   []byte
		want string
	}{
		{"empty", nil, "unexpected end"},
		{"not R", []byte("hello world"), "no XDR header"},
		{"ascii format", []byte("A\n3\n"), "not supported"},
		{"native format", []byte("B\n"), "not supported"},
		{"version 4", append([]byte("X\n"), 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0), "version 4"},
		{"truncated header", []byte("X\n\x00\x00"), "unexpected end"},
		{"bad encoding length", append([]byte("X\n"), 0, 0, 0, 3, 0, 4, 6, 2, 0, 3, 5, 0, 0, 0, 0, 99), "encoding-name"},
		{"trailing data", stream(nilv, 0), "trailing"},
		{"unknown type", stream(77), "unknown type"},
		{"bytecode", stream(21), "byte-code"},
		{"promise", stream(5), "promises"},
		{"persistent", stream(247), "persistent"},
		{"ref out of range", stream(1<<8 | 255), "out of range"},
		{"ref zero index out of range", stream(255, 0), "out of range"},
		{"negative length", stream(13, -5), "negative"},
		{"length beyond input", stream(14, 1<<30), "does not fit"},
		{"huge length", stream(14, -1, 1<<30, 0), "does not fit"},
		{"overflowing length", stream(14, -1, uint32(1<<31), uint32(0)), "out of range"},
		{"truncated length", stream(13, -1, 0), "unexpected end"},
		{"string vector too long", stream(16, 1000), "does not fit"},
		{"string element not CHARSXP", stream(16, 1, 13, 0), "expected a CHARSXP"},
		{"string element negative", stream(16, 1, 9, -7), "bad string length"},
		{"symbol without name", stream(1, 13), "expected a CHARSXP"},
		{"builtin name length", stream(8, 999), "builtin-name"},
		{"tag not symbol", stream(2|0x400, nilv, nilv, nilv), "not a symbol"},
		{"attributes not pairlist", stream(13|0x200, 1, 5, 13, 1, 5), "attributes"},
		{"untagged attribute", stream(13|0x200, 1, 5, 2, nilv, nilv), "untagged"},
		{"interior attributes", stream(2, nilv, 2|0x200, nilv), "interior"},
		{"namespace without zero", stream(249, 1), "zero is required"},
		{"altrep malformed info", stream(238, nilv, nilv, nilv), "ALTREP"},
		{"altrep bad class", stream(238, 2, nilv, nilv, nilv), "ALTREP"},
		{"empty after type", stream(13), "unexpected end"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, _, err := rds.Decode(c.in, rds.Limits{})
			if err == nil {
				t.Fatal("accepted")
			}
			var re *rds.Error
			if !errors.As(err, &re) {
				t.Fatalf("error %T is not an *rds.Error", err)
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Errorf("error %q lacks %q", err, c.want)
			}
		})
	}
}

func TestDepthBomb(t *testing.T) {
	// a list of one list of one list ... : 100000 levels in 800 KB
	parts := make([]any, 0, 200000)
	for i := 0; i < 100000; i++ {
		parts = append(parts, 19, 1)
	}
	parts = append(parts, 254)
	_, _, err := rds.Decode(stream(parts...), rds.Limits{})
	if err == nil || !strings.Contains(err.Error(), "nesting") {
		t.Fatalf("depth bomb: %v", err)
	}
	// within the limit the same shape is fine
	if _, _, err := rds.Decode(stream(append(parts[:2*100], 254)...), rds.Limits{MaxDepth: 150}); err != nil {
		t.Errorf("100 levels with MaxDepth 150: %v", err)
	}
	if _, _, err := rds.Decode(stream(append(parts[:2*100], 254)...), rds.Limits{MaxDepth: 50}); err == nil {
		t.Error("100 levels with MaxDepth 50 accepted")
	}
	// the encoder refuses what the decoder would
	var v rds.Value = rds.Null
	for i := 0; i < 5000; i++ {
		v = &rds.List{V: []rds.Value{v}}
	}
	if _, err := rds.Encode(v); err == nil {
		t.Error("Encode accepted 5000 levels")
	}
}

// A pairlist is read iteratively, so a very long chain is cheap and not
// subject to the depth limit.
func TestLongPairlist(t *testing.T) {
	parts := make([]any, 0, 3*50000)
	for i := 0; i < 50000; i++ {
		parts = append(parts, 2, 254)
	}
	parts = append(parts, 254)
	v, _, err := rds.Decode(stream(parts...), rds.Limits{MaxDepth: 10})
	if err != nil {
		t.Fatal(err)
	}
	if n := len(v.(*rds.Pairlist).Items); n != 50000 {
		t.Errorf("%d cells", n)
	}
}

func TestCompressedInput(t *testing.T) {
	plain, _ := hex.DecodeString("580a000000030004060200030500000000055554462d3800000" + "00d00000002000000010000000" + "2")
	var gz bytes.Buffer
	w := gzip.NewWriter(&gz)
	w.Write(plain)
	w.Close()
	for name, in := range map[string][]byte{"plain": plain, "gzip": gz.Bytes()} {
		v, _, err := rds.Decode(in, rds.Limits{})
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if got := v.(*rds.Integer).V; len(got) != 2 || got[0] != 1 || got[1] != 2 {
			t.Errorf("%s: %v", name, got)
		}
	}

	// bzip2, as R's memCompress(serialize(c(1L, 2L), NULL), "bzip2") writes it
	bz, _ := hex.DecodeString("425a6839314159265359c4cabfbf000013de007f1200020040010006402000314c0000d4f49918988c76eaeb246837500820218c4aabcc9f8bb9229c284862655fdf80")
	if v, _, err := rds.Decode(bz, rds.Limits{}); err != nil {
		t.Errorf("bzip2: %v", err)
	} else if len(v.(*rds.Integer).V) != 2 {
		t.Errorf("bzip2: %v", v)
	}
	if _, _, err := rds.Decode(bz[:20], rds.Limits{}); err == nil {
		t.Error("truncated bzip2 accepted")
	}

	if _, _, err := rds.Decode([]byte("\xfd7zXZ\x00rest"), rds.Limits{}); err == nil || !strings.Contains(err.Error(), "xz") {
		t.Errorf("xz: %v", err)
	}
	if _, _, err := rds.Decode([]byte{0x1f, 0x8b, 0, 0}, rds.Limits{}); err == nil {
		t.Error("corrupt gzip accepted")
	}
	if _, _, err := rds.Decode(gz.Bytes()[:gz.Len()-6], rds.Limits{}); err == nil {
		t.Error("truncated gzip accepted")
	}
}

// A decompression bomb is refused under the byte limit, not inflated.
func TestDecompressionBomb(t *testing.T) {
	var gz bytes.Buffer
	w := gzip.NewWriter(&gz)
	zeros := make([]byte, 1<<20)
	for i := 0; i < 64; i++ { // 64 MiB of zeros in about 64 KiB
		w.Write(zeros)
	}
	w.Close()
	if gz.Len() > 1<<18 {
		t.Fatalf("the bomb is %d bytes", gz.Len())
	}
	_, _, err := rds.Decode(gz.Bytes(), rds.Limits{MaxBytes: 1 << 20})
	if !errors.Is(err, rds.ErrTooLarge) {
		t.Fatalf("bomb: %v", err)
	}
	if _, _, err := rds.Decode(make([]byte, 2<<20), rds.Limits{MaxBytes: 1 << 20}); !errors.Is(err, rds.ErrTooLarge) {
		t.Errorf("large plain input: %v", err)
	}
	if _, _, err := rds.DecodeReader(bytes.NewReader(make([]byte, 2<<20)), rds.Limits{MaxBytes: 1 << 20}); !errors.Is(err, rds.ErrTooLarge) {
		t.Errorf("large reader input: %v", err)
	}
}

func TestDecodeReader(t *testing.T) {
	in := stream(13, 1, 7)
	v, _, err := rds.DecodeReader(bytes.NewReader(in), rds.Limits{})
	if err != nil || v.(*rds.Integer).V[0] != 7 {
		t.Fatalf("%v %v", v, err)
	}
	if _, _, err := rds.DecodeReader(errReader{}, rds.Limits{}); err == nil {
		t.Error("reader error swallowed")
	}
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, errors.New("disk on fire") }

// Memory stays proportional to the input: a stream that claims 2^31 elements
// is refused before any allocation.
func TestClaimedLengthsAreChecked(t *testing.T) {
	for _, typ := range []int{10, 13, 14, 15, 16, 19, 20, 24} {
		in := stream(typ, 0x7fffffff)
		if _, _, err := rds.Decode(in, rds.Limits{}); err == nil || !strings.Contains(err.Error(), "does not fit") {
			t.Errorf("type %d: %v", typ, err)
		}
	}
}

func TestNamespaceAndPackageRefs(t *testing.T) {
	// namespace:stats (ref 1) twice, in a list
	ns := []any{249, 0, 2, 9 | 64<<12, 5, "stats", 9 | 64<<12, 5, "4.6.2"}
	in := stream(append(append([]any{19, 2}, ns...), 1<<8|255)...)
	v, _, err := rds.Decode(in, rds.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	l := v.(*rds.List)
	if l.V[0] != l.V[1] || l.V[0].(*rds.NamespaceEnv).Info[0] != "stats" {
		t.Errorf("namespace reference not shared: %#v", l.V)
	}
	out, err := rds.Encode(v)
	if err != nil || !bytes.Equal(out, in) {
		t.Errorf("round trip: %v\n%x\n%x", err, out, in)
	}
}

func TestPackedAndLongReferences(t *testing.T) {
	// 2^23 distinct symbols would be needed to reach the unpacked form; the
	// writer's rule is exercised by a direct check of both encodings.
	var env rds.Value = &rds.Env{Enclos: rds.EmptyEnv, Frame: rds.Null, HashTab: rds.Null}
	l := &rds.List{V: []rds.Value{env, env}}
	b, err := rds.Encode(l)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasSuffix(b, []byte{0, 0, 1, 255}) {
		t.Errorf("packed reference missing: %x", b)
	}
	// an unpacked reference (flag 255 then the index) decodes the same
	long := stream(19, 2, 4, 0, 242, 254, 254, 254, 255, 1)
	v, _, err := rds.Decode(long, rds.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	if lv := v.(*rds.List); lv.V[0] != lv.V[1] {
		t.Error("unpacked reference not shared")
	}
}

func TestEncodeRefusals(t *testing.T) {
	cases := map[string]rds.Value{
		"nil":            nil,
		"empty pairlist": &rds.Pairlist{Kind: rds.LISTSXP},
		"bad kind":       &rds.Pairlist{Kind: rds.INTSXP, Items: []rds.Item{{Value: rds.Null}}},
		"nil child":      &rds.List{V: []rds.Value{nil}},
		"foreign":        fakeValue{},
	}
	for name, v := range cases {
		if _, err := rds.Encode(v); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
	if _, err := rds.EncodeWith(rds.Null, rds.Header{Version: 9}, 0); err == nil {
		t.Error("version 9 accepted")
	}
	alt := &rds.Altrep{Class: &rds.Sym{Name: "compact_intseq"}, Pkg: &rds.Sym{Name: "base"}, Of: rds.INTSXP, State: rds.Null}
	if _, err := rds.EncodeWith(alt, rds.Header{Version: 2}, 0); err == nil {
		t.Error("ALTREP in a version-2 stream accepted")
	}
}

type fakeValue struct{}

func (fakeValue) Type() rds.Type { return rds.NILSXP }

func TestVersion2(t *testing.T) {
	h := rds.Header{Version: 2, Writer: 0x00030400, MinReader: 0x00020300}
	b, err := rds.EncodeWith(&rds.Integer{V: []int32{9}}, h, 0)
	if err != nil {
		t.Fatal(err)
	}
	v, got, err := rds.Decode(b, rds.Limits{})
	if err != nil || got != h || v.(*rds.Integer).V[0] != 9 {
		t.Errorf("%v %+v %v", v, got, err)
	}
}

// v2Stream wraps the body of an item, taken from a version-3 encoding, in a
// version-2 header (which has no native-encoding field).
func v2Stream(t *testing.T, v rds.Value) []byte {
	t.Helper()
	b3, err := rds.Encode(v)
	if err != nil {
		t.Fatal(err)
	}
	body := b3[2+4+4+4+4+len("UTF-8"):]
	b := []byte("X\n")
	b = binary.BigEndian.AppendUint32(b, 2)
	b = binary.BigEndian.AppendUint32(b, 0x00040602)
	b = binary.BigEndian.AppendUint32(b, 0x00020300)
	return append(b, body...)
}

// R's ReadItem accepts ALTREP_SXP in a version-2 stream (serialize.c has no
// version test there) but only writes it in version 3. The decoder must hand
// back what the version-2 encoder can write: the expanded vector.
func TestAltrepInVersion2(t *testing.T) {
	alt := func(class string, state rds.Value) *rds.Altrep {
		return &rds.Altrep{Class: &rds.Sym{Name: class}, Pkg: &rds.Sym{Name: "base"}, Of: rds.INTSXP, State: state}
	}
	seq := func(n float64) rds.Value { return &rds.Real{V: []float64{n, 1, 1}} }

	v, h, err := rds.Decode(v2Stream(t, alt("compact_intseq", seq(3))), rds.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := v.(*rds.Integer); !ok || len(got.V) != 3 || got.V[2] != 3 {
		t.Fatalf("not expanded: %#v", v)
	}
	out, err := rds.EncodeWith(v, h, 0)
	if err != nil {
		t.Fatalf("decoded but cannot encode: %v", err)
	}
	if _, _, err := rds.Decode(out, rds.Limits{}); err != nil {
		t.Errorf("own output: %v", err)
	}

	if _, _, err := rds.Decode(v2Stream(t, alt("no_such_class", rds.Null)), rds.Limits{}); err == nil {
		t.Error("unknown ALTREP class accepted in a version-2 stream")
	}
	lim := rds.Limits{MaxBytes: 1 << 10}
	if _, _, err := rds.Decode(v2Stream(t, alt("compact_intseq", seq(1e9))), lim); err == nil {
		t.Error("expansion beyond the limit accepted")
	}
	// the budget is shared: two sequences that fit alone do not fit together
	two := &rds.List{V: []rds.Value{alt("compact_intseq", seq(100)), alt("compact_intseq", seq(100))}}
	if _, _, err := rds.Decode(v2Stream(t, two), lim); err == nil {
		t.Error("cumulative expansion beyond the limit accepted")
	}
}

func TestOtherNodes(t *testing.T) {
	// everything the golden set does not reach, through an encode/decode/encode cycle
	sym := &rds.Sym{Name: "x", Gp: 64}
	env := &rds.Env{
		Locked: 1, Enclos: rds.GlobalEnv,
		Frame:   &rds.Pairlist{Kind: rds.LISTSXP, Items: []rds.Item{{Tag: sym, Value: rds.Strings("v")}}},
		HashTab: rds.Null,
		Meta:    rds.Meta{Attrs: []rds.Attr{{Name: "class", Value: rds.Strings("e")}}},
	}
	v := &rds.List{V: []rds.Value{
		&rds.Complex{V: []complex128{complex(1, 2)}},
		&rds.List{Expr: true, V: []rds.Value{sym}},
		&rds.Closure{Env: env, Formals: &rds.Pairlist{Kind: rds.LISTSXP, Items: []rds.Item{{Tag: sym, Value: rds.MissingArg}}}, Body: sym},
		&rds.Closure{Formals: rds.Null, Body: rds.Null, Meta: rds.Meta{Attrs: []rds.Attr{{Name: "a", Value: rds.Null}}}},
		&rds.Builtin{Kind: rds.BUILTINSXP, Name: "sum"},
		&rds.PackageEnv{Info: []string{"package:stats"}},
		&rds.ExtPtr{Prot: rds.Null, Tag: sym},
		&rds.WeakRef{Meta: rds.Meta{Attrs: []rds.Attr{{Name: "k", Value: rds.Null}}}},
		&rds.S4{Meta: rds.Meta{Obj: true, Gp: 16, Attrs: []rds.Attr{{Name: "class", Value: rds.Strings("C")}}}},
		&rds.Pairlist{Kind: rds.DOTSXP, Items: []rds.Item{{Value: rds.UnboundValue}}, Tail: rds.BaseNamespace},
		&rds.Pairlist{Kind: rds.LISTSXP, Items: []rds.Item{{Value: rds.BaseEnv}}, Meta: rds.Meta{Attrs: []rds.Attr{{Name: "n", Value: rds.Null}}}},
		&rds.Altrep{
			Class: &rds.Sym{Name: "wrap_real"}, Pkg: &rds.Sym{Name: "base"}, Of: rds.REALSXP,
			State: &rds.List{V: []rds.Value{&rds.Real{V: []float64{1}}, &rds.Integer{V: []int32{0, 0}}}},
			Meta:  rds.Meta{Attrs: []rds.Attr{{Name: "names", Value: rds.Strings("a")}}},
		},
	}}
	first, err := rds.Encode(v)
	if err != nil {
		t.Fatal(err)
	}
	back, _, err := rds.Decode(first, rds.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	second, err := rds.Encode(back)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Errorf("not stable:\n%x\n%x", first, second)
	}
	wrapped := back.(*rds.List).V[len(v.V)-1].(*rds.Altrep)
	m, err := wrapped.Materialize(10)
	if err != nil {
		t.Fatal(err)
	}
	if r := m.(*rds.Real); r.V[0] != 1 || r.Attr("names") == nil {
		t.Errorf("wrapper: %+v", r)
	}
	if v := back.(*rds.List).V[2].(*rds.Closure); v.Env.(*rds.Env).Locked != 1 || v.Env.(*rds.Env).Attr("class") == nil {
		t.Errorf("closure env: %+v", v.Env)
	}
	for _, s := range []*rds.Special{rds.GlobalEnv, rds.MissingArg} {
		if s.String() == "" {
			t.Error("empty name")
		}
	}
	if rds.MissingArg.Type() != rds.SYMSXP || rds.GlobalEnv.Type() != rds.ENVSXP {
		t.Error("Special.Type")
	}
}

func TestMaterializeErrors(t *testing.T) {
	sym := func(n string) *rds.Sym { return &rds.Sym{Name: n} }
	mk := func(class string, state rds.Value) *rds.Altrep {
		return &rds.Altrep{Class: sym(class), Pkg: sym("base"), Of: rds.INTSXP, State: state}
	}
	real3 := func(n, start, step float64) *rds.Real { return &rds.Real{V: []float64{n, start, step}} }
	bad := []*rds.Altrep{
		mk("compact_intseq", rds.Null),
		mk("compact_intseq", real3(-1, 1, 1)),
		mk("compact_intseq", real3(5, 2147483647, 1)),
		mk("compact_intseq", real3(5, -2147483648, 1)),
		mk("wrap_real", rds.Null),
		mk("wrap_real", &rds.List{V: []rds.Value{rds.Null, rds.Null}}),
		mk("deferred_string", rds.Null),
		mk("deferred_string", &rds.Pairlist{Kind: rds.LISTSXP, Items: []rds.Item{{Value: &rds.Real{V: []float64{1}}}}}),
		mk("deferred_string", &rds.Pairlist{Kind: rds.LISTSXP, Items: []rds.Item{{Value: &rds.Integer{V: make([]int32, 20)}}}}),
		mk("deferred_string", &rds.Pairlist{Kind: rds.LISTSXP, Items: []rds.Item{{Value: mk("compact_intseq", rds.Null)}}}),
	}
	for i, a := range bad {
		if _, err := a.Materialize(10); err == nil {
			t.Errorf("case %d accepted", i)
		}
	}
	// realseq and the remaining wrapper kinds
	for _, v := range []rds.Value{&rds.Logical{}, &rds.Integer{}, &rds.Complex{}, &rds.Raw{}, &rds.Character{}, &rds.List{}} {
		a := mk("wrap_x", nil)
		a.Class = sym("wrap_list")
		a.State = &rds.List{V: []rds.Value{v, rds.Null}}
		if _, err := a.Materialize(10); err != nil {
			t.Errorf("wrapper of %T: %v", v, err)
		}
	}
	rs, err := mk("compact_realseq", real3(3, 0.5, 2)).Materialize(10)
	if err != nil || rs.(*rds.Real).V[2] != 4.5 {
		t.Errorf("realseq: %v %v", rs, err)
	}
	ds, err := mk("deferred_string", &rds.Pairlist{Kind: rds.LISTSXP, Items: []rds.Item{{Value: &rds.Integer{V: []int32{rds.NAInteger, 7}}}}}).Materialize(10)
	if err != nil || !ds.(*rds.Character).V[0].NA || ds.(*rds.Character).V[1].S != "7" {
		t.Errorf("deferred_string: %v %v", ds, err)
	}
}
