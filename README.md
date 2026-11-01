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
| M1 parser | lexer, parser, srcrefs, `deparse` | planned |
| M2 – M8 | evaluator, base, stats, regex, graphics, notebook bridge, wasm playground, CRAN ports | planned |

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
