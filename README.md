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
| M2 evaluator | `eval`: values, environments, promises, closures, control flow, replacement functions, subsetting, attributes, S3 dispatch, conditions, printing, the `Eval` API with limits; 175 of the 210 primitives and 110 base closures, environments as values, the apply family, character and regular-expression functions, the math group, factors, matching | **slices 1 to 4 done**; M3 (base R, clean-room) is planned in [docs/m3-plan.md](docs/m3-plan.md), see [docs/m2-design.md](docs/m2-design.md) |
| M3 – M8 | base, stats, regex, graphics, notebook bridge, wasm playground, CRAN ports | planned |

### M2: the evaluator (package `eval`)

```go
s := eval.NewSession(eval.Options{}) // bounded by default, see "Limits" below
res, err := s.Eval("f <- function(x, y = 2) x * y; f(21)")            // one notebook cell
res, err = s.EvalContext(ctx, src)                                    // the same, cancellable
res.Value, res.Visible        // the last value, whether R would print it
res.Output, res.Stderr        // what print/cat and the auto-print wrote; message() output
res.Warnings                  // []eval.Warning{Message, Call}
// err is *eval.Error{Message, Call, Condition, Output, Stderr, Warnings} (R's own text),
// *syntax.Error or *syntax.LimitError; a timeout or a cancel unwraps to context.DeadlineExceeded / Canceled
s.Get("x"); s.Set("x", v)     // the global environment persists from one cell to the next
```

The design is in [docs/m2-design.md](docs/m2-design.md). **Slice 2** adds: replacement functions with copy-on-modify (`x[i] <- v`, `x[[i]] <- v`, `x$a <- v`, `names(x) <- v`, `attr(x, "a") <- v`, `f(g(x)) <- v`, user `"f<-"` functions); `[`, `[[` and `$` with R's drop, NA, out-of-bounds, negative, partial-matching and matrix/array rules; attributes (`attr`, `attributes`, `class`, `dim`, `dimnames`, `levels`, `structure`, factors keep levels); S3 (`UseMethod`, `NextMethod`, internal dispatch on `[`, `[[`, `$`, their assignment forms, `length`, `print`, and the group generics `Ops`, `Math`, `Summary`); the condition system (`stop`, `warning`, `message`, `signalCondition`, `tryCatch` with `finally`, `withCallingHandlers`, `try`, `on.exit`, restarts `muffleWarning`/`muffleMessage`, `suppressWarnings`/`suppressMessages`, custom conditions); and `print`/auto-print that reproduces R's text for atomic vectors, named vectors, lists, matrices and arrays, factors, conditions, closures and attributes (R's `formatReal` with 7 digits, `[k]` labels at width 80). `matrix`, `array`, `letters` and `LETTERS` are there for the tests. The first slice: constants and symbols, environments, `<-`/`=`/`<<-`, `function` and closures with R's three-pass argument matching, promises (lazy arguments, `missing`, `...`), `{`, `if`, `for`, `while`, `repeat`, `break`, `next`, `return`, `quote`, and arithmetic, comparison and logic with R's recycling, NA rules, type promotion and integer overflow; `c`, `list`, `length`, `is.*`, `as.*`, `sum`/`prod`/`min`/`max`/`any`/`all`, `abs`/`sqrt`/`floor`/`ceiling`/`exp`, `seq_len`, `seq_along`. **94 of R's 210 primitives** plus 32 functions that are closures in R (`print`, `cat`, `paste`, `tryCatch`, `structure`, `inherits`, ...), each exercised by the goldens; the rest of the primitives raise `could not find function` until their slice. Errors and warnings are R's own: `Error in f(1) : object 'y' not found`.

**Limits (safe by default since v0.4.0).** The zero `eval.Options` is bounded: `MaxSteps` 30 000 000 calls and loop iterations, `MaxAlloc` 16 GiB of vector storage allocated in total by one Eval (copies included, so it is a budget, not live memory), `Timeout` 10 s, `MaxDepth` 5000 nested evaluations, `MaxVector` 2^28 elements. The first three values were measured, not guessed (`RGO_MEASURE=1 go test -run TestMeasureLimits -v ./eval`): the 1 154 goldens need at most 8 879 steps, 0.8 MB and 3 ms, and the heaviest cells a notebook runs on purpose (a loop of 10^6 iterations: 3 million steps, 0.17 s; a list grown by 20 000 assignments: 1.6 GB of copies, 0.63 s) set the scale; each default is about ten times that. A limit that is hit ends the evaluation with an R error that names it and leaves the session usable. To remove a limit you must say so: `MaxSteps: eval.Unbounded` (also for `MaxAlloc` and `Timeout`); a zero never means unbounded and any other negative value panics. A context given to `EvalContext` still stops an unbounded session. The bomb tests (`while (TRUE) 0`, `numeric(1e12)`, `1:1e12`, unbounded recursion, an allocation loop) assert a refusal under a Go memory limit, and in CI under `ulimit -v`.

Judged by the oracle, never by expectations written by hand:

| check | result |
|---|---|
| goldens: 1 154 snippets of `testdata/eval` (math group, language primitives, factors and matching, environments, apply, strings, regular expressions, arithmetic, comparison, sequences, coercion, summaries, closures, matching, laziness, loops, subsetting, assignment, attributes, S3, conditions, printing, errors), each compared with R 4.6.2's serialized value (byte for byte), visibility, printed output, standard error, warnings, and error message and call; replayed by `go test` without R | **1 154 / 1 154** |
| differential fuzz, arithmetic/logic/summary (`go run ./cmd/rparity evalfuzz`; seeds 3, 6, 7, 8) | **189 998 / 190 000** |
| differential fuzz, subsetting and assignment on vectors, lists, matrices and arrays (`evalfuzz --mode subset`; seeds 7, 11, 21) | **35 963 / 36 000** (see below) |
| differential fuzz, regular expressions: random POSIX extended patterns, with lazy quantifiers, back-references, `\<` `\>` `\b` and non-ASCII text, through `grepl`, `sub`, `gsub`, `regexpr`, `gregexpr`, `strsplit`, `regmatches` (`--mode regex`; seeds 1, 7, 21; each function reported on its own line; nothing excluded) | **60 000 / 60 000** |
| differential fuzz, character functions: `nchar`, `substr`, `substring`, `toupper`, `trimws`, `strsplit`, `sprintf`, `format`, `chartr` (`--mode strings`; seeds 1, 5, 21) | **20 000 / 20 000** |
| differential fuzz, `lapply`, `sapply`, `vapply`, `Map`, `mapply`, `Reduce`, `Filter`, `unlist`, `rep`, `seq`, `rev` (`--mode apply`; seeds 3, 21) | **18 000 / 18 000** |
| differential fuzz, the math group (`--mode math`; seeds 1, 4, 9). **Exclusions, counted:** 5 979 of the last run's 20 000 expressions are libm functions (trigonometric, hyperbolic, logarithmic, gamma family, `sinpi` family, `log(x, base)`) compared through `round(x, 6)`, because their last bit depends on the C library (macOS and glibc differ from each other and from Go); `digamma`, `trigamma`, `gamma` and `lgamma` at 1e-10 are left out (a different summation order, 1e-16 relative) | **20 000 / 20 000** |
| `FuzzEval`: random source never panics, hangs or allocates without bound (session limits: depth, steps, vector length; the corpus keeps the 9 crashers it found) | about 70 million executions in the last runs, clean; the corpus keeps the 13 crashers it found |

**Slice 3** adds: environments as values (`environment`, `environment<-`, `new.env`, `globalenv`, `emptyenv`, `baseenv`, `parent.env`, `parent.frame`, `sys.call`, `sys.function`, `assign`, `get`, `get0`, `exists`, `mget`, `ls`, `rm`, `local`, `eval`, `evalq`, `list2env`, `as.environment`); the apply family (`lapply`, `sapply`, `vapply`, `Map`, `mapply`, `Reduce`, `Filter`, `Find`, `Position`, `do.call`, `match.fun`); character functions (`nchar`, `substr`, `substring`, `strsplit`, `sprintf`, `format`, `toupper`, `tolower`, `casefold`, `chartr`, `trimws`, `startsWith`, `endsWith`, `Encoding`, `enc2utf8`, `utf8ToInt`, `intToUtf8`, `validUTF8`); and regular expressions (`grepl`, `grep`, `sub`, `gsub`, `regexpr`, `gregexpr`, `regmatches`) with POSIX extended syntax and leftmost-longest matching as R's default engine TRE has. **Regular expressions are TRE.** The matcher is [github.com/go-embedded-r/tre](https://github.com/go-embedded-r/tre), a Go port of TRE (BSD-2, the upstream notice kept there) with the options that reproduce R 4.6.2's modified copy (`{,n}`, bound overflow, the backtracking restart test, shared bracket positions, negated classes in copies of a repetition); no R source is copied. Lazy quantifiers, capture priority in repeated groups, back-references (`(a)\\1`, run by TRE's backtracking matcher, refused by a step limit when exponential) and `\<` `\>` behave as in R, and R's choice of byte mode for all-ASCII input and wide mode otherwise is reproduced, with its quirks (an error text without the pattern for `grepl`/`regexpr` in wide mode, empty-match tokens of `strsplit`). `perl = TRUE` (PCRE) still stops with an error that says it is not supported; TRE's approximate-matching syntax (`{~1}`) is not ported and is an error; classes of non-ASCII characters follow Go's Unicode tables, not a C library's. Widths of `sprintf` and `format` above 2^24 are refused; `sprintf("%05d", NA)` zero-pads the `NA` as macOS does, where glibc pads with spaces (the C library decides in R, so the golden leaves it out). `format` of a factor and `sort` of non-atomic values are not implemented; `sort` and `ls` order strings by code point, not by locale.

**Slice 4** adds the rest of the math group (`log`, `log2`, `log10`, `log1p`, `expm1`, trigonometric and hyperbolic functions and their inverses, `sinpi` family, `atan2`, `gamma`, `lgamma`, `digamma`, `trigamma`, `trunc`, `sign`, `round`, `signif` with R 4's rounding algorithm), `cumsum`/`cumprod`/`cummax`/`cummin`, `range`, `is.finite`/`is.infinite`/`is.nan`, `anyNA`, `nzchar`, `is.matrix`/`is.array`/`is.object`/`is.recursive`/`is.language`, `length<-`, `storage.mode<-`, `switch`, `substitute`, `nargs`, `...length`/`...elt`/`...names`, `call`, `as.call`, `forceAndCall`, `.Primitive`, `.subset`/`.subset2`, `::`/`:::` (resolved in the base and global environments), `%*%`/`crossprod`/`tcrossprod`, `seq.int`, `proc.time`; and `typeof`, `mode`, `deparse`, `as.name`, `unique`, `duplicated`, `match`, `%in%`, `factor`, `levels`, `format` of factors. **Collation**: strings sort and compare as R 4.6.2 does under the oracle's `C.UTF-8` (ICU: punctuation, digits, letters; accents after the plain letter; lowercase before uppercase), see [docs/m3-plan.md](docs/m3-plan.md) section 5; non-Latin scripts compare by code point after the Latin letters. **Not supported, refused with an error that says why** (`eval.Unsupported`, listed by `go run ./cmd/rparity missing`, not counted as implemented): the complex and raw types (`Re`, `Im`, `Mod`, `Arg`, `Conj`, `as.complex`, `as.raw`), S4 (`@`, `standardGeneric`, `.isMethodsDispatchOn`), formulas (`~`) and expression vectors (`expression`), foreign calls (`.C`, `.Call`, `.External*`, `.Fortran`, `lazyLoadDBfetch`), `.Internal`, and the tracing and debugging primitives. Coverage of base's 1 362 exports (`go run ./cmd/rparity basecov`): 269 agree with the oracle, 19.8 %, at v0.7.0 (196, 14.4 %, at v0.6.0).

Known divergences of slice 2: `[[<-` with a vector subscript on a list that goes through an NA, a name that is absent or a logical (`x[[c(1, NA)]] <- v`, `x[[c("a", "c")]] <- v` on an unnamed list: R's message comes from its recursive indexing, ours differs, 31 of 26 000 fuzz cases); a few R messages for `[[<-` on atomic vectors with 0, negative or NA subscripts follow the oracle for the common cases only; `print` of a function does not show R's `<environment: 0x...>` address (a placeholder), `.Primitive` formals are not shown, and printing an array with more than 2^20 slices or 2^22 rows is refused instead of running for ever; `print` of list elements dispatches only to the methods R registers itself (factor, condition), as R does. Known divergences, all from the platform rather than the semantics: the bits of a NaN (its sign, whether it is quiet) are whatever the hardware made of it, so goldens and the fuzz compare NA against NaN, not their payloads; `exp` and the power with a fractional or negative exponent can differ from the C library's `pow`/`exp` in the last bit (the fuzz keeps the exponents integral and leaves `exp` out; they come with go-numerics in M4); the 2 of 180 000 above are one of each (a `pow` last bit, and `NA %% 3` added to a NaN); `min`/`max` of a character vector mixed with numbers is not collated like R's. The modulus family is written with explicit fused multiply-adds, as R is compiled on arm64, so it gives the same result everywhere; R on x86-64 Linux does not fuse, so its `%%` of huge doubles can differ there (the CI lane reports the evaluator fuzz on Linux, it does not gate it).

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

The parser (M1) is released as [v0.1.0](https://github.com/go-embedded-r/r/releases/tag/v0.1.0).
The organisation landing page is <https://go-embedded-r.github.io/> and the
documentation site is <https://go-embedded-r.github.io/docs/>.
