# r — an R interpreter in pure Go

`github.com/go-embedded-r/r` is the start of an R implementation written in Go
(`CGO_ENABLED=0`), meant to be embedded in Go programs and to run in the
browser through WebAssembly. It will back the ` ```r ` cells of
[go-notebook](https://github.com/go-notebook) the way
[julia](https://github.com/go-embedded-julia/julia) backs its Julia cells and
[rbgo](https://github.com/go-embedded-ruby/ruby) its Ruby cells.

The plan is [go-scikit/roadmap `bib/r.md`](https://github.com/go-scikit/roadmap/blob/main/bib/r.md).
Every stage is checked against the reference implementation, the pinned
**GNU R 4.6.2** (CRAN tarball, SVN revision 90784), never against expectations
written by hand.

## Licence, and what "clean room" means here

This repository is **BSD-3-Clause** (see [LICENSE](LICENSE)). The sources of R
are GPL-2-or-later, so:

- **R is an oracle, never a source.** R is run as a black box and what it
  prints is recorded. No R source file, no R test file, no R help page, no
  `.Rout.save` transcript and no text of R's manuals is copied into this
  repository, and none is downloaded at run time. The golden files are the
  answers of R 4.6.2 to inputs written here.
- **Base R's R code is rewritten, not copied.** The R-language parts of the
  base packages (`base`, `stats`, `utils`, `methods`, …) will be re-implemented
  under the BSD licence from the documented behaviour and from what the oracle
  does. Fetching the GPL code on first use, and a GPL port, were considered
  and rejected (decision of 2026-11-01, recorded in the roadmap).
- **Contributors work from specifications**: the *R Language Definition*, the
  *R Internals* manual, the help pages, and the oracle's behaviour. Do not open
  `src/main`, `src/nmath` or the R code of the base packages while writing the
  counterpart. See [CONTRIBUTING.md](CONTRIBUTING.md) and [NOTICE](NOTICE).

## Status

| milestone | contents | state |
|---|---|---|
| M0 oracle harness | `internal/oracle`, `cmd/roracle`, `rds` (serialization reader/writer), golden files | **done** |
| M1 parser | `syntax`: lexer, parser, AST with source map, tokens, R's parse errors, `deparse`, `TopLevelNames` | **done** |
| M2 – M8 | evaluator, base, stats, regex, graphics, notebook bridge, wasm playground, CRAN ports | planned |

### M1: the parser (package `syntax`)

```go
f, err := syntax.ParseAll(src)                    // or ParseAllLimits(src, syntax.Limits{MaxBytes, MaxDepth, MaxNodes})
names, err := syntax.TopLevelNames(cell)          // names a cell assigns and reads
lines := syntax.Deparse(f.Exprs[0], nil)          // deparse(), one string per line
```

- **Tree.** R's own shape: operators, `if`/`for`/`{`/`[` and the rest are calls to a
  symbol; `Function` carries the formals; `|>` is rewritten as R does; `->` is
  `<-` with the sides swapped; a string callee is a symbol. Every node has a
  `Span` (line, column and byte range). `File.Tokens` are the terminals of
  `getParseData()` (`SYMBOL_FUNCTION_CALL`, `SYMBOL_SUB`, `EQ_FORMALS`, `COMMENT`, `LINE_DIRECTIVE`...).
- **Errors.** `*syntax.Error` has `Line`, `Col`, `Message` and its text is the
  message `parse()` signals, byte for byte, source context and caret included;
  limits give a `*syntax.LimitError` (wraps `ErrLimit`). The statements before
  an error are kept in the returned `File`.
- **Bounds.** Bytes, tree depth (also for chains such as `a+a+a+...`, built by a
  loop) and node count; R's own 50-level context stack (`contextstack overflow`).
  Bomb tests (200 000-deep chains of every shape) assert a refusal under a
  memory limit.
- **`TopLevelNames`.** `Assigned`: `<-`, `=`, `->`, `"x" <-`, `f(x) <-`, `x$a <-`
  (assign x), `assign("lit", ...)`, `for` variables. `Read`: free variables and
  called names not bound earlier in the cell, under R scoping (function
  parameters and locals, `local()`; a function body may use anything the cell
  assigns anywhere). `Dynamic` (with `Reasons`) when names cannot be resolved:
  `<<-`, assign/get with a computed name or another environment, eval, parse, source, attach, library, require, sys.function, `environment<-`, load, list2env...
- **Numbers.** A decimal literal is computed in double precision, digit by
  digit with an explicit fused multiply-add (`math.FMA`), then scaled by a power
  of ten: what R does on macOS arm64, where `long double` is `double` and the
  compiler fuses `10*ans + digit`. 150 000 random literals are identical to
  that R's, bit for bit, and the result is the same on every platform. R on
  x86-64 Linux accumulates in 80-bit extended precision and can differ by one
  unit in the last place on literals of more than 15 digits (the CI oracle lane
  saw it on `1e300`), and by up to 6 ulp on literals near the subnormal range
  (measured on R's own `tests/`: 290 literals differ, at most by 6 ulp). The parser is deterministic and matches the 64-bit
  reference; the platform-dependent literals are in `testdata/numbers.R`,
  compared within 16 ulp and kept out of the goldens. One literal differs by more: `1.7976931348623157e308` (the largest double) is `Inf` in 64-bit-double R and in the Go parser, and is kept by 80-bit R. Hexadecimal literals
  follow R's rule (`0x1.8p1`).

#### Parity with R 4.6.2

Judged by the oracle, never by expectations written by hand: `go run ./cmd/rparity ...`
parses the same text with R and with Go and compares (i) the terminal tokens
and positions with `getParseData()`, (ii) the tree with a dump of the language
objects of `parse(keep.source = FALSE)` (names, constants as IEEE bits, empty
arguments, formals), (iii) the error text with `parse()`'s. R's own sources are
read where they are unpacked and never committed.

| corpus | inputs | tree | tokens | deparse (`deparse()` and `print()`) | error text, byte for byte |
|---|---|---|---|---|---|
| R `src/library`: every `.R` file of the base packages (7.8 MB) | 844 (842 parse, 2 rejected) | **842 / 842** | **842 / 842** | **842 / 842** | **2 / 2** |
| R `tests/`: 95 files (1 is not UTF-8 and is skipped) | 94 (90 parse, 4 rejected) | **90 / 90** | **90 / 90** | **90 / 90** | **4 / 4** |
| our snippets (`testdata/snippets.R`) | 42 | **42 / 42** | **42 / 42** | **42 / 42** | n/a |
| our error cases (`testdata/errors.R`) | 52 (1 parses) | **1 / 1** | **1 / 1** | **1 / 1** | **51 / 51** |
| our edge cases: Unicode, raw strings, `\|>`, `\(x)`, `->>`, hex/`L`/`i` numbers (`testdata/edge.R`) | 88 (3 rejected) | **85 / 85** | **85 / 85** | **85 / 85** | **3 / 3** |
| our deparse cases: parentheses, precedence, line breaking, indentation, constants, names (`testdata/deparse.R`) | 65 (2 rejected) | **63 / 63** | **63 / 63** | **63 / 63** | **2 / 2** |
| odd input: CRLF, tabs, form feeds, Unicode spaces, bidi controls, invalid UTF-8, `=>`, 50 and 5000 parentheses (`parity.Extras`, `values.R`) | 120 | all | all | all | all |
| differential fuzz: 610 000 random edits of lines of R's sources (seeds 21, 31, 41, 42) | 610 000 (about two thirds rejected) | **all** | **all** | **all** | **all** |
| decimal literals, random (up to 20 digits, exponents -330..330) | 100 000 | **100 000 identical doubles** | | 99 995 identical, 5 not (below) | |

Deparse is written from the documented behaviour of `deparse()` (precedence and parenthesisation, the 60-byte line cutoff with its break points and indentation, `if`/`else` layout inside braces, `structure`-free constants, 15-significant-digit numbers, string escapes with R's own table of printable characters, generated by the oracle into `syntax/print_gen.go`; the list of primitives that deparse as function calls into `syntax/prims_gen.go`). Not reachable byte for byte: numeric constants with a decimal exponent beyond about 200 in magnitude, where R takes a power of ten from the C library's `exp10`, which on macOS is not always the correctly rounded value: 5 of 100 000 random literals of that kind (`1e229`, `1e303`, `1e-244`) deparse with one more digit in Go (`1.00000000000000e+229`). No such constant occurs in R's sources or tests. On x86-64 Linux, where R parses and formats numbers in 80-bit `long double`, 6 of the 90 files of R's `tests/` deparse with a different last digit in some constant (`1e-270` against `9.99999999999999e-271`, `-5152.07935578853` against `...52`): the CI lane reports that verdict and gates the other three columns and all of `src/library`. R on x86-64 Linux formats with `long double` and its own digits routine; the CI lane compares the same corpora there.

Speed, one core each, interleaved rounds, median of 11 after a warm-up round, 1-minute load 7.8 on an Apple M-series laptop: the 842 files of
`src/library` that both parse (7.8 MB) take **0.214 s** in Go (`syntax.ParseAll`, which also builds the source map and the token list) and
**0.269 s** in R's `parse(text =, keep.source = FALSE)`: **Go/R = 0.80**. (`go run ./cmd/rparity bench DIR`.)

Known gaps: the effect of `#line` directives on the line numbers of later errors;
`parse(file=)`
encoding options; an invalid UTF-8 input is parsed as R's lexer does, and
for such text R cannot produce parse data, so tokens are not compared.

### M0: what exists

- **`internal/oracle`** starts one `Rscript --vanilla` process (the pinned
  release; any other version is refused at the greeting) and talks to it over
  stdin/stdout. Requests are R's own serialization format, written by the Go
  `rds` package; replies are R serializations, one hex line each, so no JSON
  package is needed on the R side. Operations: `eval`, `deparse`, `parsedata`
  (`utils::getParseData`), `describe`, `reserialize`. Output, warnings and the
  condition class and message are captured and returned; each request gets a
  fresh environment; a per-call timeout kills R; `R_MAX_VSIZE` bounds its
  vector heap. The test suite starts a fake R from its own test binary, so the
  protocol and every failure mode are tested without R.
- **`rds`** reads and writes R's serialization format (XDR, versions 2 and 3,
  gzip or bzip2 compressed): vectors of every atomic type, lists, pairlists,
  calls, symbols, closures, environments (with sharing and cycles), namespace
  and package references, built-ins, external pointers, S4 objects and ALTREP
  objects (kept in serialized form; `Altrep.Materialize` expands the compact
  sequences, wrappers and deferred strings). A decoded value encodes back to
  the **same bytes**: attributes in order, the object and general-purpose
  bits, string encoding flags, NaN payloads (`NA_real_` is a NaN). Byte-code
  and promises are refused. The decoder is built for untrusted input: size
  limit after decompression (decompression bombs are refused), nesting limit,
  every claimed length checked against the bytes left before allocating. It
  is fuzzed (`FuzzDecode`, corpus in `rds/testdata/fuzz`).
- **`cmd/roracle`** regenerates the golden files (`roracle regen`), prints the
  identity of the R that answers (`roracle version`) and evaluates a snippet
  (`roracle eval 'expr'`).
- **Goldens** (`testdata/golden`): `deparse.jsonl` is `parse` + `deparse` of
  the snippets in `testdata/snippets.R` (the input of the M1 parser), and
  `rds.jsonl` is `serialize()` of the expressions of `testdata/values.R`. Both
  inputs are written for this repository. The tests replay the goldens without
  R: every value R serialized decodes, agrees with R's `typeof` and `length`,
  and re-encodes to R's exact bytes; values built in Go encode to the bytes R
  writes. With R, the live tests check that R reads Go's bytes back unchanged
  and that regenerating the goldens reproduces the committed files.

## Getting the oracle

The tests that need R skip when `Rscript` is not found. To run them, install
the pinned release, which conda-forge publishes for macOS arm64 and Linux
x86-64, and point `RGO_RSCRIPT` at it:

```
micromamba create -y -p ~/r-4.6.2 -c conda-forge --override-channels r-base=4.6.2
export RGO_RSCRIPT=~/r-4.6.2/bin/Rscript
RGO_REQUIRE_ORACLE=1 go test ./...      # a missing R is a failure, not a skip
go run ./cmd/roracle regen              # rewrite testdata/golden
```

`RGO_REQUIRE_ORACLE=1` turns the skips into failures (CI sets it).

## Layout

| path | role |
|---|---|
| `rds/` | R serialization reader and writer |
| `internal/oracle/` | the driver of the reference R, and the R side of the protocol (`oracle.R`) |
| `internal/gen/` | golden-file generators and readers |
| `cmd/roracle/` | the command that regenerates the goldens |
| `testdata/` | inputs written for this repository, and `golden/` |

There is no release yet: tags will follow once the parser (M1) lands. The
organisation landing page and the documentation site are not set up at M0; this
README is the documentation.
