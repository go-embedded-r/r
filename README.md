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

### Clean-room process: the similarity gate

A rule that depends on what a contributor (human or model) remembers is not
enough: v0.12.0's `eval/baselib/datetime.R` was written from the help pages
but partly from memory of R's code, and came out close to R's text. It was
rewritten in v0.12.2 and `cmd/cleanroom` now checks every file of
`eval/baselib` against the R library code of R 4.6.2.

- **Method.** Our own lexer (package `syntax`) cuts both files into tokens
  (comments dropped); k-grams of tokens are hashed and thinned by winnowing
  (Schleimer, Wilkerson, Aiken, *Winnowing: Local Algorithms for Document
  Fingerprinting*, SIGMOD 2003: the minimum hash of each window, the rightmost
  on a tie; any shared run of `k+w-1` tokens shares a fingerprint, no run
  shorter than `k` can). Two views are compared: identifiers kept (*exact*) and
  every identifier replaced by one placeholder (*normalised*, which sees a
  renamed copy). The score is containment, the share of our fingerprints found
  in one R file. The longest common token run is reported with `file:line` on
  both sides. Runs count code tokens only (literals and commas do not, so a
  table of month names is not a copy). Parameters: `k = 8`, `w = 4`.
- **R's sources are read, never stored.** The checker reads
  `src/library/*/R/*.R` of an R 4.6.2 tree at check time (`RGO_RSOURCE`, or
  `--rsrc`) and prints only references (`file:line`), never R's text. CI
  downloads the tarball by SHA-256 in the oracle lane and runs the gate there
  (`RGO_REQUIRE_RSOURCE=1`: the gate must run, not skip).
- **Thresholds, and how they were chosen.** Measured on code known to be
  independent of R's library code, at `k = 8, w = 4`:

  | independent code (files of at least 200 tokens) | exact containment | normalised containment | longest exact run | longest normalised run |
  |---|---|---|---|---|
  | our own test inputs (`testdata/*.R`, `testdata/eval/*.R`; 25 files), max | 0.04 | 0.07 | 10 | 29 |
  | R's own tests and demos against its library code (144 files), p95 / p99 / max | 0.03 / 0.19 / 0.37 | 0.31 / 0.44 / 0.52 | 14 / 31 / 46 | 29 / 48 / 78 |

  Formal argument lists (`function(` up to the matching `)`, defaults
  included) are not fingerprinted in either view: a help page fixes the
  signature, so independent code must repeat it.

  The gate sits above the p99 of the second row: exact containment **0.20**,
  normalised containment **0.60** (these two only for files of at least 200
  tokens: a file of 100 tokens has a dozen fingerprints), longest exact run
  **40** code tokens, longest normalised run **60**. Files of 30 to 199 tokens
  are gated on an absolute count instead: **4** exact fingerprints shared with
  one R file (independent small files: at most 2 of 32). Limits of the method: a one-line function has one
  possible spelling, so a tiny file can be a copy and pass; the gate bounds
  accidents and careless restructuring, it does not replace the rule above.
- **The oracle does not show R's R code.** Printing a base function in the
  reference R shows its GPL body (this happened once, while checking a
  behaviour). `internal/oracle/oracle.R` therefore refuses, with the error
  `oracle: printing base closures is disallowed by the clean-room rule`, to
  print, deparse, `dput`, `body`, `as.list`, `format`, serialize or return any
  closure of base, stats, utils, methods and the other packages of R, however
  it is reached (auto-print, `:::`, `getAnywhere`, inside a list or an
  attribute, or a copy made outside the namespace with `environment<-`,
  `body<-`, `attr<-` or `as.function`, which is recognised by `identical()`
  against the bodies and formals of the functions of those packages, at any
  size); `.Internal` is refused in requests. Closures of our own print, but a
  user function that is textually a base function (`function(x) x` is
  `identity`) is refused too. Goldens that compare function printing use
  closures defined in the snippet. Limits: a base function copied with
  different formals is recognised only if its body is 100 serialized bytes or
  more; someone with an installed R can still read its sources, which is why
  the rule is also procedural:
- **Provenance in every PR.** `.github/pull_request_template.md` asks, for any
  change under `eval/baselib/`: the `.Rd` pages read, whether any R function
  body was seen and how (and who then rewrote the file blind), and the gate
  score of the first draft before any reshaping.
- **Run it.**

  ```sh
  export RGO_RSOURCE=/path/to/R-4.6.2          # an unpacked, checked R source tree
  go run ./cmd/cleanroom check                 # eval/baselib/*.R; exit 1 over the gate
  go run ./cmd/cleanroom check --list 20 f.R   # every shared run of 20+ tokens, as file:line
  go run ./cmd/cleanroom floor some/*.R        # the same scores for independent code
  go test ./internal/cleanroom ./cmd/cleanroom # the gate as a test; skips without R's sources
  ```

- **Audit of `eval/baselib`** (v0.12.3; k = 8, w = 4). "Best R file" is the one
  with the highest normalised containment (against R's largest files that is
  generic structure, which is why the normalised column alone is not the gate).
  "Before" is the first draft, scored with formal argument lists still
  fingerprinted; "after" is the file in the repository, scored as the gate now
  does. Runs are in code tokens, as `file:line` of ours against R's.

  | File | Best R file | Exact | Normalised | Longest exact run | Rewritten |
  |---|---|---|---|---|---|
  | `datetime.R` | `base/R/datetime.R` | 0.11 (before 0.23) | 0.34 (0.47) | 36, ours 374-377 / `base/R/dates.R` 125-128 (88) | yes, from the Rd pages and the oracle, without the old file |
  | `vectors.R` | `tools/R/QC.R` | 0.00 (0.00) | 0.40 (0.41) | 22, ours 87-93 / `base/R/matrix.R` 27-32 (74, `ifelse`) | `ifelse`, `append`, `isTRUE`, `isFALSE`, `NROW`, `NCOL` |
  | `conditions.R` | `methods/R/RMethodUtils.R` | 0.00 (0.27) | 0.43 (0.43) | 6, ours 13 / `base/R/message.R` 82 (13) | yes, blind, by a second agent (7 shared fingerprints before, all signatures; 5 after the first rewrite) |
  | `env.R` | `tools/R/QC.R` | 0.00 (0.01) | 0.52 (0.55) | 22, ours 8-10 / `base/R/factor.R` 148-151 (26) | no |
  | `support.R` | `tools/R/QC.R` | 0.01 (0.01) | 0.41 (0.47) | 30, ours 55-57 / `base/R/unname.R` 19-23 (35) | no |
  | `strings.R` | `base/R/namespace.R` | 0.00 (0.03) | 0.23 (0.28) | 26, ours 25-29 / `base/R/print.R` 59-65 (31) | no |
  | `arrays.R` (v0.14.0) | `tools/R/QC.R` | 0.00 (0.00) | 0.38 (0.37) | 26, ours 55-58 / `base/R/matrix.R` 49-52 (60, ours 214-223 / `base/R/scale.R` 36-45: **the first draft failed the gate**) | `scale.default` redesigned (it was close to R's text, a sign it had been recalled rather than derived from the help page): per-column vectors and one subtraction and division, instead of two `sweep` calls and a nested `if`; the rest as first written |
  | `groups.R` (v0.14.0) | `tools/R/QC.R` | 0.00 (0.01) | 0.43 (0.43) | 21, ours 142-143 / `base/R/tapply.R` 27-28 (51, ours 66-69 / `base/R/rle.R` 19-26: **the first draft failed the gate**) | `rle` redesigned (the empty case was a special return that R also has): a condition that needs no special case, run starts by `cumsum`/`tabulate`; the rest as first written |
  | `connections.R` (v0.13.0) | `tools/R/check.R` | 0.00 (0.00) | 0.42 (0.43) | 13, ours 34-35 / `utils/R/capture.output.R` 36-38 (38, ours 34-39 / the same file 36-45) | no: written from the `connections`, `sink` and `capture.output` help pages and the oracle; the first draft is the file, plus `summary`, one line |

  The gate is exact containment 0.20, normalised 0.60, exact run 40, normalised
  run 60 (4 shared fingerprints under 200 tokens); every file is under it. The
  longest runs left (22 to 36 tokens) are a documented one-liner or an error
  message the oracle fixes.

## Status

| milestone | contents | state |
|---|---|---|
| M0 oracle harness | `internal/oracle`, `cmd/roracle`, `rds` (serialization reader/writer), golden files | **done** |
| M1 parser | `syntax`: lexer, parser, AST with source map, tokens, R's parse errors, `deparse`, `TopLevelNames` | **done** |
| M2 evaluator | `eval`: values, environments, promises, closures, control flow, replacement functions, subsetting, attributes, S3, conditions, printing, the `Eval` API with safe limits, apply family, strings, regular expressions (over `go-embedded-r/tre`), math group, factors; 175 of the 210 primitives implemented and the other 35 refused with an error that says why (complex, raw, S4, formulas, expression vectors, foreign calls, `.Internal`) | **done** (v0.7.0), see [docs/m2-design.md](docs/m2-design.md) |
| M3 base R, clean-room | rewrite of base's 1 152 closures under BSD from the help pages and the oracle, never from R's sources; plan, order and coverage metric in [docs/m3-plan.md](docs/m3-plan.md); `go run ./cmd/rparity basecov` | **in progress**: 40.8 % of base's 1 362 exports agree with the oracle (556; v0.14.0; 37.7 % at v0.13.0, 33.9 % at v0.12.0, 26.4 % at v0.11.0, 24.4 % at v0.10.0, 23.3 % at v0.9.0, 22.2 % at v0.8.0, 19.8 % at v0.7.0, 14.4 % at v0.6.0) |
| M4 – M8 | numerics, stats, graphics, notebook bridge, wasm playground, CRAN ports | planned |

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

**Limits (safe by default since v0.4.0).** The zero `eval.Options` is bounded: `MaxSteps` 30 000 000 calls and loop iterations, `MaxAlloc` 16 GiB of vector storage allocated in total by one Eval (copies included, so it is a budget, not live memory), `Timeout` 10 s, `MaxDepth` 5000 nested evaluations, `MaxVector` 2^28 elements. The first three values were measured, not guessed (`RGO_MEASURE=1 go test -run TestMeasureLimits -v ./eval`): the 1 154 goldens need at most 8 879 steps, 0.8 MB and 3 ms, and the heaviest cells a notebook runs on purpose (a loop of 10^6 iterations: 3 million steps, 0.17 s; a list grown by 20 000 assignments: 1.6 GB of copies, 0.63 s) set the scale; each default is about ten times that. A limit that is hit ends the evaluation with an R error that names it and leaves the session usable. To remove a limit you must say so: `MaxSteps: eval.Unbounded` (also for `MaxAlloc` and `Timeout`); a zero never means unbounded and any other negative value panics. A context given to `EvalContext` still stops an unbounded session. The bomb tests (`while (TRUE) 0`, `numeric(1e12)`, `1:1e12`, unbounded recursion, an allocation loop) assert a refusal under a Go memory limit, and in CI under `ulimit -v`. Since v0.13.0 the files are bounded too: `MaxFileBytes` (default `DefaultMaxFileBytes`, 256 MiB; `eval.Unbounded` removes it) limits the in-memory file system, any one file read or written under `FileRoot`, and what a text connection holds; files and the environment variables of a session are its own (see slice 10).

Judged by the oracle, never by expectations written by hand:

| check | result |
|---|---|
| goldens: 2 707 snippets (9 of them recorded for macOS only, `testdata/eval/platform_darwin.R`, and compared only there) of `testdata/eval` (matrices, arrays and grouping, connections and files, dates and times, environments and namespaces, vectors, math group, language primitives, factors and matching, environments, apply, strings, regular expressions, arithmetic, comparison, sequences, coercion, summaries, closures, matching, laziness, loops, subsetting, assignment, attributes, S3, conditions, printing, errors), each compared with R 4.6.2's serialized value (byte for byte), visibility, printed output, standard error, warnings, and error message and call; replayed by `go test` without R | **2 707 / 2 707** |
| differential fuzz, environments (`--mode env`; seeds 1, 2, 3, 5 000 programs each): `get`, `exists`, `get0`, `mget` with modes, `inherits` and `ifnotfound`; `assign`, `$<-`, `[[<-`, `local`, `evalq` into locked bindings and environments; `rm`; the lock functions; active bindings; `delayedAssign`; `as.list`, `eapply`, `with`, `within`; the namespace functions; `::`, `:::`; `topenv`, `environmentName`. **Exclusions, generated on purpose and counted (1 858 of the 15 000 programs, not in the figure)**: 473 print an environment's address (compared with the address replaced, 0 differ), 458 print a builtin (our `print` shows no formals; compared as `.Primitive("name")`, 0 differ), 485 use a package R installs besides base (compared exactly: 363 differ, the divergence of slice 8), 442 list a hashed environment unsorted (compared exactly: 210 differ, R's order is undefined) | **13 142 / 13 142** compared, plus the exclusions above |
| differential fuzz, dates and times (`--mode datetime`; seeds 1, 2, 3, 6 000 programs each): formatting with 40 conversions (`%s` included), parsing, `POSIXlt` fields, round trips through text, `difftime` in every unit, arithmetic and comparison on `Date`/`POSIXct`/`difftime`, `seq` by day, week, month, quarter, year, `c`, `rep`, `rev`, `sort`, `unique`, time zone conversion in 8 zones; years 1971 to 2037, never 01:00 to 03:59 local. **Exclusions, generated on purpose and counted (1 848 of the 18 000 programs, not in the figure)**: 462 in years before 1902 (R warns once per session, so the warning is not compared; the rest is, 0 differ), 463 in years after 2037 in zones with daylight saving (compared exactly, 0 differ), 445 `tzone` attributes of zones with unusual daylight abbreviations (compared exactly, 56 differ), 478 local times that a daylight-saving change skips or repeats (compared exactly, 64 differ) | **16 152 / 16 152** compared, plus the exclusions above |
| differential fuzz, matrices, arrays and grouping (`--mode arrays`; seeds 1, 2, 3, 6 000 programs each, 14 families: `apply` on matrices and 3-d arrays, `sweep`, `scale`, `aperm`, `prop.table`/`marginSums`, `rownames`/`colnames` and their replacements, `lower.tri`, `diag<-`, `split`/`unsplit`/`split<-`, `tapply` (one and two factors, `default`), `rle`/`inverse.rle`, `sequence`, `replace`/`setequal`/`zapsmall`/`factorial`, `replicate` and `sapply(simplify = "array")`, `bquote` and subsetting of calls; the values are integers and halves, so that no C-library noise is in play; nothing excluded). The first run found 11 differences in 2 000 programs (the names `apply` gives the dimnames of its result, a list matrix printed as a list), fixed over three more runs | **18 000 / 18 000** |
| differential fuzz, connections and files (`--mode io`; seeds 1, 2, 3, 6 000 programs each, 11 families: text connections read and written, `cat` separators, files by `cat`/`writeLines`/`readLines` and by connection, the file-system functions in a temporary directory, `sink` and `capture.output`, `readRDS`/`saveRDS` round trips, `Sys.setenv`/`Sys.getenv`, `source`; a path never reaches the output). **One exclusion, generated on purpose and counted (1 706 of the 18 000 programs, not in the figure)**: non-ASCII text that goes through a connection (R leaves it unmarked, ours is marked UTF-8: compared exactly except for the mark, 0 differ) | **16 294 / 16 294** compared (the first run found 215 differences in 1 000 programs: `cat`'s separators, `textConnection(open = "a")`'s warning, a tag on a nested object; fixed) |
| differential fuzz, arithmetic/logic/summary (`go run ./cmd/rparity evalfuzz`; seeds 3, 6, 7, 8) | **189 998 / 190 000** |
| differential fuzz, subsetting and assignment on vectors, lists, matrices and arrays (`evalfuzz --mode subset`; seeds 7, 11, 21) | **35 963 / 36 000** (see below) |
| differential fuzz, regular expressions: random POSIX extended patterns, with lazy quantifiers, back-references, `\<` `\>` `\b` and non-ASCII text, through `grepl`, `sub`, `gsub`, `regexpr`, `gregexpr`, `strsplit`, `regmatches` (`--mode regex`; seeds 1, 7, 21; each function reported on its own line; nothing excluded) | **60 000 / 60 000** |
| differential fuzz, character functions: `nchar`, `substr`, `substring`, `toupper`, `trimws`, `strsplit`, `sprintf`, `format`, `chartr` (`--mode strings`; seeds 1, 5, 21) | **20 000 / 20 000** |
| differential fuzz, `lapply`, `sapply`, `vapply`, `Map`, `mapply`, `Reduce`, `Filter`, `unlist`, `rep`, `seq`, `rev` (`--mode apply`; seeds 3, 21) | **18 000 / 18 000** |
| differential fuzz, strings, part 2: `strrep`, `substr<-`, `strtoi`, `make.unique`, `basename`, `dirname`, `formatC`, `prettyNum`, `toString`, `sQuote`/`dQuote`/`shQuote`, `noquote`, `file.path` (`--mode strings2`; seeds 1, 3; no exclusions) | **16 000 / 16 000** (after the first run's 255 differences were fixed) |
| differential fuzz, vectors: `as.vector`, `order`, `rank`, `pmax`, `pmin`, `which.max`, `tabulate`, row/column sums and means, `t`, `diag`, `outer`, `mean`, set operations, `head`, `tail`, `ifelse`, `diff`, `append` (`--mode vec`; seeds 1, 5; no exclusions) | **18 000 / 18 000** |
| differential fuzz, the math group (`--mode math`; seed 9, 20 000 expressions). **Comparison classes, counted:** exact for everything except 5 644 libm cases (trigonometric, hyperbolic, logarithmic, `gamma`, `expm1`, `log1p`, `log(x, base)`), compared within 4 ulp, where the largest difference seen is 2 ulp (`log`, `log10`) and 1 ulp for the others: last-bit noise of the C library, which differs between macOS and glibc; and 335 gamma-family cases (`lgamma`, `digamma`, `trigamma`) compared through `signif(x, 12)` because our routines are less accurate than R's, an open defect (issue filed). `sinpi`, `cospi`, `tanpi` agree exactly | **20 000 / 20 000** |
| `FuzzEval`: random source never panics, hangs or allocates without bound (session limits: depth, steps, vector length; the corpus keeps the 9 crashers it found) | about 70 million executions in the last runs, clean; the corpus keeps the 16 crashers it found (the newest: `identical` of a call and a symbol panicked on a type assertion, found through `bquote(0()())`; before it, `unique` of a list of 44 000 elements took 6 s, pairwise, and is hashed now) |

**Slice 11 (M3, v0.14.0)** adds matrices, arrays and grouping, with 351 new goldens (`testdata/eval/arrays.R`) and a fuzz mode (`--mode arrays`); basecov goes from 513 to 556 of 1 362 (37.7 % to 40.8 %). **Functions** (`eval/baselib/arrays.R` and `groups.R`, written from the help pages and the oracle): `apply` (any number of margins, names as margins, results of any shape: a vector, an array with the dimnames R gives it, or a list), `sweep` (with its two `check.margin` warnings), `scale` (a generic, with `scale.default`), `aperm` (generic, `perm` by number or by dimnames name, `resize`), `prop.table`/`proportions`, `margin.table`/`marginSums`, `%o%`, `as.matrix` (generic; the default method), `mat.or.vec`, `lower.tri`, `upper.tri`, `row`, `col` (`as.factor`), `diag<-`, `rownames`, `colnames` (`do.NULL`, `prefix`) and their replacement forms, `split` (generic; one factor, or a list of factors crossed with `sep` and `lex.order`; `drop`; a grouping shorter than the data is recycled with R's warning), `split<-`, `unsplit`, `tapply` (one or several factors, `default`, `simplify`, array and list results), `rle` and `inverse.rle`, `sequence` (with its recycling warning), `replace`, `replicate`, `setequal`, `getElement`, `force`, `eval.parent`, `%||%`, `%notin%`, `grepv`, `factorial`, `lfactorial`, `zapsmall`, `bquote` (with `splice = TRUE`), `is.vector`. **In the evaluator**, because the above needs them: a matrix or array whose dimnames are named prints R's header lines (the row-label offset rule, the column line, `, , c = u` for slices); a one-dimensional array prints as a named vector with the name of its dimension over it, and `names()`, `[` and `[[` on it use its dimnames; a list matrix prints as a grid of one-line summaries (`numeric,3`); `[[` and `[` work on calls (a call is the list of its function and arguments, and `q[-1]` is a call again); `sapply(simplify = "array")`; `array(NULL)` is R's error; the functionals (`lapply`, `sapply`, `vapply`, `Map`, `mapply`, `Reduce`, `Filter`, `Find`, `Position`) return a visible value whatever their function evaluated last. **Divergences**: there are no data frames, so none of the `.data.frame` methods and no `table`; `rle` objects print as plain lists (R's `print.rle` needs `str`); `kronecker` and `%x%` are not in this release (the oracle cannot evaluate them, so they could not be checked); `Vectorize` and `alist` are not either; `split`'s grouping is by factor levels only (no `interaction` function); `tapply` takes factors or what `as.factor` can make of them; `aperm` is written in R, with index arithmetic over the whole array (memory of the order of the array, bounded by the vector limit); `is.primitive` is not in this release (the evaluator's builtins are not R's primitives); the error call of a failure inside `as.matrix.default` and `colnames` shows our code where R shows its own (two goldens that would have copied the text of R's call were left out); `factorial`, `lfactorial` and `zapsmall` inherit the last-bit differences of `gamma`, `lgamma` and `round` (the goldens compare them through `signif`). **Found by the fuzz and fixed**: the names `apply` gives the first dimension of its result (only when the function's results are named and as long as the dimension it iterates over); a result list from `apply` is not unlisted unless the first result is a plain vector; `tapply` of vectors of different lengths is an error; `split(drop = TRUE)` keeps the order of the levels.

**Slice 10 (M3, v0.13.0)** adds connections and files, with 328 new goldens (`testdata/eval/connections.R`) and a fuzz mode (`--mode io`); basecov goes from 462 to 513 of 1 362 (33.9 % to 37.7 %). **The file system is private by default**: a session has an in-memory tree (bounded by `Options.MaxFileBytes`, default 256 MiB) that the host never sees and that no other session shares; `Options.FileRoot` gives the session a directory of the host as its root instead, confined with `os.Root` (`..` is clamped at the root, symbolic links cannot leave it; `Session.Close` releases it). `Sys.getenv` reads a table of the session's own (`Options.Env`, empty by default), never the host's environment. **Connections**: `file` (modes `r`, `w`, `a`, with `b`, `t` and `+`), `textConnection` (read; write and append to a variable, `local`, anonymous with `textConnectionValue`; the variable is updated line by line), `stdin`, `stdout`, `stderr`, `open`, `close`, `flush`, `isOpen`, `isIncomplete`, `summary` (generic; `summary.connection` and `print.connection`), `readLines` (`n`, `ok`, `warn`; LF, CRLF and CR; the incomplete-line warning), `writeLines` (`sep`), `readline` (non-interactive: writes the prompt, returns `""`), `nullfile`, `isatty`. `cat` takes `file` (a path or a connection), `append`, and places its separators as R does (before every argument but the first, `NULL` excepted; a separator with a newline in it ends the output with one). `sink` (a stack of 20, `append`, `split`, `type = "message"`) and `capture.output` divert the evaluator's own output, so what a cell prints reaches a connection and never the `Result`; a sink left open stays across cells. **Files**: `file.exists`, `dir.exists`, `file.create`, `file.remove`, `file.rename`, `file.copy`, `file.append`, `file.size`, `unlink` (wildcards, `recursive`), `dir.create` (`recursive`), `list.files`/`dir` (`pattern`, `all.files`, `full.names`, `recursive`, `ignore.case`, `include.dirs`, `no..`, `fixed`), `list.dirs`, `tempfile`, `tempdir`, `getwd`, `setwd`, `normalizePath`, `path.expand`, with R's warnings and messages. **Serialization and code**: `readRDS` and `saveRDS` over package `rds` (gzip, `compress = FALSE`, `version = 2`; vectors, lists, attributes, factors, dates, calls), `source` (`local`, `print.eval`; value and visibility returned as R does) and `sys.source`. **Environment variables and time**: `Sys.getenv`, `Sys.setenv`, `Sys.unsetenv`, `Sys.sleep` (it obeys the time limit and the context). **Bounds**: a file or text connection past `MaxFileBytes` is refused (`the file size limit of this session was reached`), 128 connections and 20 sinks at most, 65 536 entries in the in-memory tree, paths of at most 4 096 bytes and 256 levels (v0.13.1), `readRDS` refuses a document that expands past the limit; the bomb tests assert these under a memory limit. **Divergences**: a connection is an integer with its class and no `conn_id` attribute, numbered from 3 (R's are numbered from what is open); `stdout()` is always connection 1 (R's is the current sink); strings read from a connection are not marked unknown as R marks them, so `Encoding()` and the serialized bytes of non-ASCII text differ (the fuzz counts these, 1 706 programs); the variable of an output text connection is not locked; `file("")` (the anonymous file), `gzfile`, `bzfile`, `url`, `pipe`, `fifo`, sockets, `rawConnection`, `readBin`, `writeBin`, `readChar`, `scan`, `seek`, `pushBack`, `showConnections` and `closeAllConnections` do not exist yet; `saveRDS` refuses closures, environments (other than the four standard ones), `ascii = TRUE`, and `bzip2`/`xz`, which R writes; `readRDS` reads a pairlist as a list and refuses environments, closures and S4 objects; `source` refuses `echo`, `verbose` and `exprs` (no expression vectors), and `parse` does not exist; `Sys.getenv()` with no argument returns a plain named vector (R's has class `Dlist`), `Sys.setenv(TZ = )` does not change the zone the date functions use (they still read the process's `TZ`); `path.expand` leaves `~` alone; `file.size` of a directory is 4 096; there is no `file.info`, `file.mtime` or `Sys.glob`; `summary` has no default method, so `summary(1:3)` says there is no applicable method where it said it could not find the function; a sink's output is moved to its connection when the evaluator next checks its limits or runs a connection function (so it is not copied while it piles up), and what a cell diverted is never in its `Result`. **Fixed on the way** (found by the fuzz): `cat` separated by position instead of by argument; an element of a list that is an object printed under its `[[i]]` tag where R's default method prints none.

**Slice 9 (M3, v0.12.0)** adds dates and times, with 452 new goldens and a fuzz mode (`--mode datetime`). **Classes**: `Date`, `POSIXct`, `POSIXlt` (a list of the 11 components with `tzone` and `balanced`), `difftime`. **Functions**: `Sys.time`, `Sys.Date`, `as.Date`, `as.POSIXct`, `as.POSIXlt` (character with `format` and `tryFormats`, numeric with `origin`, `Date`, between classes, factors), `strptime` (every conversion R's help page lists, `%OS`, `%z`, 24:00, names matched without case, trailing text ignored), `format`/`strftime`/`as.character` (`usetz`, `digits`, `%OSn`, `options(digits.secs)`), `ISOdate`, `ISOdatetime`, `weekdays`, `months`, `quarters`, `julian`, `difftime`, `as.difftime`, `units`, `units<-`, `format`/`print` of differences, `seq.Date` (day, week, month, quarter, year, a number, a difftime), `c`, `[`, `[[`, `[<-`, `rep`, `rev`, `sort`, `unique`, `mean`, `range`, `diff` on all three classes, `+`, `-`, `*`, `/`, the comparisons and `Math`/`Summary` of `difftime` with R's rules (a time difference added to a date wins over the `difftime` group method; mixed units go to seconds; comparing date-times of different zones warns). Non-finite times (`Inf`, `NaN`, `NA`) are kept apart as R 4.2 does. **Time zones** come from Go's embedded IANA database (`time/tzdata`, about 450 KB in the binary): `tz = ""` is the session's zone, the `TZ` variable else the system's (the oracle runs with `TZ=UTC` and so do the goldens); a repeated local hour reads as its first occurrence, as R's does. **Support added on the way**: `options`, `getOption`, `deparse1`, `gettext`, `gettextf`, `rep.int`, `is.factor`, `pmatch`, `unname`; `as.character` and `paste` of a factor give its labels (they gave the codes); a comparison keeps names and shape but no longer the class; `print` forwards its arguments to the method; a list element that is an object is printed by its method (looked up from the global environment, as R's console does); `deparse` writes attributes as `structure(...)`; `as.numeric` finds `as.double.<class>` methods; `lapply` over an object goes through its `as.list` method; `rev`, `sort` and `unique` keep the class of what they are given. **Divergences**: the C locale only (English names, `%x` is `%m/%d/%y`, `%c` the C form; `%E` and `%O` modifiers read as the plain conversion); R's warning `datetimes before 1902 may not be accurate` (once per session) is never given; the daylight abbreviation in the `tzone` attribute of a zone with unusual or no daylight time is found by scanning Go's transitions and differs from R's for some zones (56 of 445 fuzz cases: the attribute has three elements, the first is always right); of the local times that a daylight-saving change skips or repeats, ours read as R's in most cases, and 64 of 478 fuzz cases differ. `%s` is not the epoch time on macOS: R gives the fields read as UTC unless the format also has `%z` or `%Z`, in which case it subtracts the zone's offset; that is reproduced exactly (the 9 goldens of `platform_darwin.R`, the fuzz), and whether glibc does the same is not known, so those goldens are not compared on Linux. The call shown in an error from a long `structure(...)` breaks its lines where R does not; `Date` and `POSIXct` have no `round`, `trunc`, `cut`, `summary`, `seq.POSIXt`, `Sys.setenv`, `Sys.timezone`, `OlsonNames`, and there are no data frames yet, so no date columns. `NewSession` now costs 56 µs and 1 091 allocations (626 at v0.11.0), measured at 1-minute load 5. The growth is not the time zone tables (Go loads a zone on first use): `datetime.R` alone is 420 allocations, 110 closures defined eagerly at session start; installing a library file on first use would remove them.

**Slice 8 (M3, v0.11.0)** adds environments and a one-namespace model, with 268 new goldens and a fuzz mode (`--mode env`). **Bindings**: `lockBinding`, `unlockBinding`, `bindingIsLocked`, `lockEnvironment(bindings =)`, `environmentIsLocked`, `makeActiveBinding`, `bindingIsActive`, `activeBindingFunction`, `delayedAssign`; every assignment (`<-`, `<<-`, `assign`, `$<-`, `[[<-`, `list2env`, replacement forms) honours locks and runs active bindings, with R's messages and calls; the base environment and all its bindings are locked as in R (`assign("sum", 1, baseenv())` is refused). `rm` takes `inherits`; `get`/`get0`/`exists`/`mget` take `mode` (several modes are read as R does), `inherits`, `ifnotfound` (values or functions) and give R's messages, curly quotes included; `exists` does not force a promise. `as.list` is generic, with `as.list.environment` (`all.names`, `sorted`), `.default` (vectors, calls), `.function` and `.factor`; `eapply`, `with`, `within` (lists), `withVisible`, `names(env)`. Environments carry attributes (`attr<-`, `class<-`, so S3 dispatch works on them; `environmentName` honours a `name` attribute). **Namespaces**: `asNamespace`, `getNamespace`, `isNamespace`, `isBaseNamespace`, `getNamespaceName`, `getNamespaceExports`, `isNamespaceLoaded`, `requireNamespace`, `.BaseNamespaceEnv`, `.GlobalEnv`, `topenv` with its real rule, `::` and `:::`. The namespace of base is an environment of its own (printed `<environment: namespace:base>`, enclosed by the global environment) that shares the base frame. Code run by `eval`, `evalq`, `local` and `with` now reports R's call for errors raised directly in it (`eval(quote(x), new.env())`). Divergences: there is no package loading, so no `library`, `require`, `loadNamespace`, `attach`, `search`, `sys.source` (planned with the slices that need them); the installed packages other than base (`stats`, `utils`, ...) resolve through our base environment (`utils::head` works, `stats::sd` is "'sd' is not an exported object") and any other name is `there is no package called ‘x’`; `parent.env(globalenv())` is base, R's is the last attached package; `as.list.environment` and `names(env)` list in creation order where R's hashed order is unspecified; `sys.call()` and `parent.frame()` inside `eval`/`local` do not see R's eval frame; the text `requireNamespace("foo")` writes without `quiet` was read from the oracle, not goldened (its harness adds a connection warning); `print` of an environment does not list its attributes; `getNamespace(<environment>)` raises R's message with another call; `within` has no data frame method.

**Slice 3** adds: environments as values (`environment`, `environment<-`, `new.env`, `globalenv`, `emptyenv`, `baseenv`, `parent.env`, `parent.frame`, `sys.call`, `sys.function`, `assign`, `get`, `get0`, `exists`, `mget`, `ls`, `rm`, `local`, `eval`, `evalq`, `list2env`, `as.environment`); the apply family (`lapply`, `sapply`, `vapply`, `Map`, `mapply`, `Reduce`, `Filter`, `Find`, `Position`, `do.call`, `match.fun`); character functions (`nchar`, `substr`, `substring`, `strsplit`, `sprintf`, `format`, `toupper`, `tolower`, `casefold`, `chartr`, `trimws`, `startsWith`, `endsWith`, `Encoding`, `enc2utf8`, `utf8ToInt`, `intToUtf8`, `validUTF8`); and regular expressions (`grepl`, `grep`, `sub`, `gsub`, `regexpr`, `gregexpr`, `regmatches`) with POSIX extended syntax and leftmost-longest matching as R's default engine TRE has. **Regular expressions are TRE.** The matcher is [github.com/go-embedded-r/tre](https://github.com/go-embedded-r/tre), a Go port of TRE (BSD-2, the upstream notice kept there) with the options that reproduce R 4.6.2's modified copy (`{,n}`, bound overflow, the backtracking restart test, shared bracket positions, negated classes in copies of a repetition); no R source is copied. Lazy quantifiers, capture priority in repeated groups, back-references (`(a)\\1`, run by TRE's backtracking matcher, refused by a step limit when exponential) and `\<` `\>` behave as in R, and R's choice of byte mode for all-ASCII input and wide mode otherwise is reproduced, with its quirks (an error text without the pattern for `grepl`/`regexpr` in wide mode, empty-match tokens of `strsplit`). `perl = TRUE` (PCRE) still stops with an error that says it is not supported; TRE's approximate-matching syntax (`{~1}`) is not ported and is an error; classes of non-ASCII characters follow Go's Unicode tables, not a C library's. Patterns beyond tre v0.1.3's compile limits (nesting thousands deep, huge repetition expansion, more than about 2^17 nodes) stop with `regular expression '...' is too large for this evaluator`; no golden or fuzz case reaches them, R compiles some of them. Widths of `sprintf` and `format` above 2^24 are refused; `sprintf("%05d", NA)` zero-pads the `NA` as macOS does, where glibc pads with spaces (the C library decides in R, so the golden leaves it out). `format` of a factor and `sort` of non-atomic values are not implemented; `sort` and `ls` order strings by code point, not by locale.

**Slice 7 (M3, v0.10.0)** adds restarts (`withRestarts` with functions and arguments, `computeRestarts`, `findRestart`, `invokeRestart(r, ...)`, `tryInvokeRestart`), the call stack (`sys.nframe`, `sys.calls`, `sys.frames`, `sys.frame`, `sys.on.exit`), `Recall`, `match.call` (positional, partial and exact matching, `...` expanded from the caller, `expand.dots`), `errorCondition`, `warningCondition`, `geterrmessage`, and `$` on a call. Divergences, all because R runs `tryCatch` and `try` through closures we do not have: the `conditionCall` of a `stop()` made directly inside `try()` or `tryCatch()` at top level is empty here and R's internal call there (`doTryCatch(...)`, `try(...)`); `match.call(expand.dots = FALSE)` keeps `...` as a list where R keeps a pairlist (`deparse` shows `list(` for `pairlist(`); absolute frame numbers differ from the oracle's, which runs each snippet inside its own harness frames, so the goldens compare numbers of frames relative to a starting point; `computeRestarts` has no top-level `abort` restart.

**Slice 6 (M3, v0.9.0)** adds string functions: `strrep`, `substr<-`, `strtoi`, `make.unique`, `basename`, `dirname`, `formatC` (formats `d`, `f`, `e`, `g`, `s`; flags `-`, `0`, `+`, space; `big.mark`; negative width), `prettyNum`, `toString`, `sQuote`/`dQuote` (curly quotes, as the oracle's UTF-8 locale gives), `shQuote` (type `sh`), `noquote` and `print(quote = FALSE)`, `file.path`. Part 2 adds `format(big.mark =)`, `formatC(mode =, format = "fg")`, `shQuote(type = "csh"/"cmd"/"cmd2")` and `match.arg` (native special: defaults from the calling function's formals, partial matching, R's error text), and `sprintf("%#x", 0)` prints no prefix, as C does.

**Slice 5 (M3, v0.8.0)** starts the base library: `eval/baselib/*.R` is base R's glue written by us (`isTRUE`, `isFALSE`, `xor`, `Negate`, `setNames`, `union`, `intersect`, `setdiff`, `is.element`, `append`, `head`, `tail`, `ifelse`, `nrow`, `ncol`, `NROW`, `NCOL`, `lengths`, `rep_len`, `diff`, `outer`), embedded, parsed once and evaluated into the base environment of each session (43 µs per `NewSession`), under the clean-room rule of [docs/m3-plan.md](docs/m3-plan.md); Go natives for `as.vector`, `order`, `rank`, `pmax`, `pmin`, `which.max`, `which.min`, `tabulate`, `rowSums`/`colSums`/`rowMeans`/`colMeans`, `t`, `diag`, `mean`, `stopifnot`. `mean` sums in `float64` with R's second pass where R uses `long double`: last-bit differences are possible on long sums. `rank` supports `average`, `first`, `min`, `max`.

**Slice 4** adds the rest of the math group (`log`, `log2`, `log10`, `log1p`, `expm1`, trigonometric and hyperbolic functions and their inverses, `sinpi` family, `atan2`, `gamma`, `lgamma`, `digamma`, `trigamma`, `trunc`, `sign`, `round`, `signif` with R 4's rounding algorithm), `cumsum`/`cumprod`/`cummax`/`cummin`, `range`, `is.finite`/`is.infinite`/`is.nan`, `anyNA`, `nzchar`, `is.matrix`/`is.array`/`is.object`/`is.recursive`/`is.language`, `length<-`, `storage.mode<-`, `switch`, `substitute`, `nargs`, `...length`/`...elt`/`...names`, `call`, `as.call`, `forceAndCall`, `.Primitive`, `.subset`/`.subset2`, `::`/`:::` (resolved in the base and global environments), `%*%`/`crossprod`/`tcrossprod`, `seq.int`, `proc.time`; and `typeof`, `mode`, `deparse`, `as.name`, `unique`, `duplicated`, `match`, `%in%`, `factor`, `levels`, `format` of factors. **Collation**: strings sort and compare as R 4.6.2 does under the oracle's `C.UTF-8` (ICU: punctuation, digits, letters; accents after the plain letter; lowercase before uppercase), see [docs/m3-plan.md](docs/m3-plan.md) section 5; non-Latin scripts compare by code point after the Latin letters. **Not supported, refused with an error that says why** (`eval.Unsupported`, listed by `go run ./cmd/rparity missing`, not counted as implemented): the complex and raw types (`Re`, `Im`, `Mod`, `Arg`, `Conj`, `as.complex`, `as.raw`), S4 (`@`, `standardGeneric`, `.isMethodsDispatchOn`), formulas (`~`) and expression vectors (`expression`), foreign calls (`.C`, `.Call`, `.External*`, `.Fortran`, `lazyLoadDBfetch`), `.Internal`, and the tracing and debugging primitives. Coverage of base's 1 362 exports (`go run ./cmd/rparity basecov`): 269 agree with the oracle, 19.8 %, at v0.7.0 (196, 14.4 %, at v0.6.0).

Divergence of slice 5, fixed in v0.14.0: a matrix whose `dimnames` carry names (`names(dimnames(x))`, as `table()` makes) printed without the variable-name header lines R adds.

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
