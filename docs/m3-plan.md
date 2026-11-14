# M3 plan: base R, rewritten clean-room

Status: plan; the measuring tools (section 4) exist, the first slice is v0.7.0.

M2 gave an evaluator with 175 of R's 210 primitives. M3 is the rest of what
`library(base)` exports: 1 152 closures. The user's decision stands: base R's R
code is **rewritten under the BSD licence and never copied**. This document says
what that means in practice, in what order the work goes, and how progress is
measured.

## 1. The clean-room rule for base functions

- **Specification**: R's documentation (the help pages, the R language
  definition, R-ints, the "Writing R Extensions" manual) and the **observed
  behaviour of the oracle** (R 4.6.2 run through `internal/oracle`).
- **Never read, never copied**: the files of `src/library/base/R`, `src/main`,
  `src/nmath`, `tests/`. `CONTRIBUTING.md` says so for contributors; the tools
  below respect it: they work on *names and counts* obtained from the running
  oracle (`ls(baseenv())`, a walk over each closure's body that records only the
  names of the `.Internal` functions it calls), and keep nothing of the bodies.
- **How a function is written**: from its help page and from experiments. The
  author writes the experiments as snippets in `testdata/eval/*.R`, the oracle
  turns them into goldens, the Go (or R) code is made to agree. A function whose
  help page leaves a behaviour open is pinned by a snippet before it is written,
  never after.
- **Two homes for the code**:
  1. **Go natives** (`eval/*.go`), for anything on a hot path or that is
     naturally a loop over vectors (`paste`, `unique`, `sort`, `rev`, `which`,
     `pmax`, `cumsum`, the apply family, string functions).
  2. **Our own R** in `eval/baselib/*.R`, embedded and evaluated into the base
     environment when a session starts, for the closures that are glue
     (`Negate`, `isTRUE`, `Vectorize`, `match.arg`, `stopifnot`, `modifyList`,
     `mapply` wrappers, `do.call` helpers, argument checking). These files are
     written by us, carry the BSD header, and are subject to the same rule: the
     header of each states which help pages it was written from. A snippet of
     R code that matches R's own text by more than a signature is a defect and
     is rewritten.
- `.Internal(f(...))` itself is not provided (`rparity missing` lists it as
  refused): base closures are rewritten as Go or as our R and do not go through
  it. The 401 internals that base calls (`testdata/base/internals.tsv`) are
  therefore a *measure of work*, not an API to reproduce.

## 2. Order

The ranking is by what unlocks the most exports, not by call frequency: base
R's calls to `.Internal` are spread flat (the leaders, `as.vector` and `vector`,
have 8 and 7 sites; most internals one), so frequency ranks nothing.
`go run ./cmd/rparity basecov` prints, from `testdata/base/needs.tsv`, the
greedy order of the missing internals that completes the most closures, and the
count of closures that need none (a rewrite suffices: **701 of the 1 042
closures not yet implemented** at the time of writing). The slices:

| slice | contents | why first |
|---|---|---|
| 4 (v0.7.0) | the rest of M2: 65 primitives (math group, cumulative functions, `range`, predicates, `switch`, `substitute`, `call`, `%*%`), `factor`, `match`, `unique`, collation | the primitives are what every closure is made of |
| 5 (v0.8.0, done: 22.2 %) | vectors: `as.vector`, `vector`, `rev`, `order`, `rank`, `pmax`, `pmin`, `which.max`, `tabulate`, `table`, `rowSums` family, `outer`, `diag`, `t`, `apply`, `sweep`, `rowsum`, set operations | the largest group of pure closures; `as.vector` alone unlocks 6 |
| 6 (v0.9.0, part done: 23.3 %) | strings and text: `paste` family, `format` family (`formatC`, `prettyNum`, `format.Date`), `sprintf` completion, `toString`, `sprintf`-based `message`, `strtoi`, `chartr`, `casefold`, `iconv`, `shQuote`, `dQuote`, `noquote`, `cat` options | `format` and `paste` are called by most of the rest |
| 7 (v0.10.0, done: 24.4 %) | conditions and the call stack: `.addCondHands` family as Go, `tryCatch` and `withCallingHandlers` over it, `warning()` and `stop()` with `call.`, `conditionCall`, `simpleCondition`, restarts (`withRestarts`, `computeRestarts`, `invokeRestart`), `sys.*`, `match.call`, `match.arg`, `Recall`, `on.exit` options, `trace`-less `browser` refusal | the condition internals are the `.Internal` family that base uses most as a group |
| 8 (v0.11.0, done: see README) | environments and namespaces: `lockBinding` family, `lockEnvironment`, `makeActiveBinding`, `delayedAssign`, `as.list.environment`, `eapply`, `with`/`within`, `mget`/`get`/`get0`/`exists`/`rm` with modes and `inherits`, a namespace model of one namespace (base: `asNamespace`, `getNamespace`, `isNamespace`, `::`, `:::`, `topenv`, `requireNamespace`), attributes on environments. **Deferred to the slices that need them**: `sys.source` (connections, slice 10), `library`/`require`/`loadNamespace` and a table of namespaces for ported packages (there is no package to load yet), `attach`/`detach`/`search` (the search path) | needed to host ported packages |
| 9 (v0.12.0, done: see README) | dates and times: `Sys.time`, `Sys.Date`, `Date`, `POSIXct`, `POSIXlt`, `difftime` with their arithmetic, comparison, `c`, `[`, `rep`, `seq.Date`, `format`/`strptime`/`strftime`/`as.Date`/`as.POSIXct`/`as.POSIXlt` (zones from Go's embedded IANA database), `ISOdate`, `weekdays`/`months`/`quarters`/`julian`, `options`. **Deferred**: `round`/`trunc`/`cut` of dates, `seq.POSIXt`, `Sys.setenv`/`Sys.timezone`/`OlsonNames` (with the system calls of slice 10), `summary.Date` | self-contained, many exports (every method counts as an export) |
| 10 (v0.13.0, done: see README) | connections and files: `file`, `textConnection`, `readLines`, `writeLines`, `cat` to files, `sink`, `capture.output`, the file-system functions over a private in-memory file system (`Options.FileRoot` for a confined host directory), `readRDS`/`saveRDS` over `rds`, `source`/`sys.source`, `Sys.getenv`/`Sys.setenv` over the session's own table, `Sys.sleep`. **Deferred**: `scan`, `readline` input, `readBin`/`writeBin`/`readChar`/`writeChar`/`rawConnection` (need the raw type, slice 12), `gzfile`/`bzfile`/`url`/`pipe`/sockets, `file.info`/`file.mtime`, `parse` (expression vectors, slice 12), `dput`/`dget`/`dump`, `save`/`load`, `read.table` family in `utils` | the notebook bridge needs them |
| 11 (v0.14.0, done: see README) | matrices, arrays and grouping: `apply`, `sweep`, `scale`, `aperm`, `prop.table`, `margin.table`, `rownames`/`colnames`, `lower.tri`, `diag<-`, `split`, `unsplit`, `tapply`, `rle`, `sequence`, `replace`, `replicate`, `bquote`, `is.vector` and the printing they need (named dimnames, one-dimensional arrays, list matrices). Taken ahead of linear algebra because 530 closures needed no missing internal and the numerics wait for go-numerics. **Deferred**: `table` and everything with data frames, `kronecker`, `Vectorize`, `alist`, `interaction`, `cut`, `strwrap`, `pretty`, `print.rle` (needs `str`) | 40 % of base's exports; the notebook idioms (`apply`, `tapply`, `split`) |
| 12 (v0.15.0, done: see README) | the complex and raw types, expression vectors: `raw`, `as.raw`, `rawToChar`, `charToRaw`, `rawToBits`, `intToBits`, `packBits`, `rawShift`, `bitwAnd` and the other `bitw` functions, `complex`, `as.complex`, `Re`, `Im`, `Mod`, `Arg`, `Conj` and the math group on complex numbers, `expression`, `as.expression`, `parse`, `str2lang`, `str2expression`, `eval` of an expression vector; `rawConnection`, `rawConnectionValue`, `readBin`, `writeBin`, `readChar`, `writeChar`, `seek`, `isSeekable` over `eval/conn.go`. **Deferred**: `polyroot` (its last bits are those of a Jenkins-Traub implementation), formulas (`~`) and S4 (they need attributes on calls and the methods package), `as.hexmode`, `dput`, `serialize`, `gzfile`/`memCompress` (no zlib is written here: the language-neutral engines are shared, see go-embedded-ruby/ruby#918) | the value types that the numeric and stats work and the slice-10 deferrals needed |
| 13 | linear algebra that base exports: `solve`, `qr`, `chol`, `svd`, `eigen`, `det` | depends on the extraction of go-ndarray's linear algebra into a shared library (go-embedded-ruby/ruby#918), not on go-numerics; until then refused |
| 14 (v0.16.0, done: see README) | data frames and `table`, which slices 11 and 12 deferred: `data.frame`, `as.data.frame` (vector, list, matrix, table), `is.data.frame`, `[`/`[[`/`$` and their replacement forms, `dim`/`dimnames`/`row.names`, `rbind`/`cbind` of data frames and of vectors and matrices, `merge`, `split` of a data frame, `subset`, `transform`, `droplevels`, `print`/`format`/`as.matrix`/`t` of a data frame, `table`, `print.table`, `as.data.frame.table`, `expand.grid`, `interaction`, `cut`; goldens sharded per snippet file. **Deferred**: `summary` of a data frame or a table, the `Ops`/`Math`/`Summary` methods of data frames, `aggregate`, `by`, `stack`, `within`, `ftable`, `xtabs`, matrix and list columns; `head` and `tail` are `utils`, outside base | the notebook idioms (`data.frame`, `merge`, `table`, `cut`); 48.4 % of base's exports (659) |

A slice is tagged when its goldens, fuzz, `FuzzEval` and bomb tests are green and
`rparity basecov` has moved.

### Shared engines (decided 2026-11-10 and 2026-11-11)

go-embedded-r does not write its own copy of an engine that several languages
of the fleet need (go-embedded-ruby/ruby#918). It binds the shared library and
keeps only a thin layer that is faithful to GNU R: call shapes, argument
defaults, the result lists, the messages.

- **Numerics** (M4, stats): `optim`, `optimize`, `uniroot`, `nlm`, `nlminb`,
  `nls`, `integrate`, ODE solvers, and the special functions (`gamma`, `beta`,
  `besselJ` and the rest) bind
  [go-numerics](https://github.com/go-numerics). Where GNU R's algorithm differs
  from go-numerics' default, the R layer selects R's (for example `optim`'s
  Nelder–Mead constants, `uniroot`'s tolerance rule, `integrate`'s QUADPACK
  routine) and records in its divergences when a result is not bit-identical to
  GNU R's. go-numerics now states an accuracy contract: a bound on the error
  against an extended-precision oracle (math/big or mpmath), see go-numerics
  optimize#61 and special#69. When a result is not bit-identical to GNU R's,
  the R layer quotes that bound as the guarantee instead of calling the
  difference "last-bit noise"; bit-exact arithmetic stays an option.
- **What R does better goes upstream.** When R's reference has an algorithm, a
  default or a test problem that go-numerics lacks, it is added to go-numerics
  for the whole fleet, not to the R layer. This applies only to what is
  available outside R's own GPL text: published papers and netlib's public
  domain or BSD code (QUADPACK, MINPACK, PORT). R's own C and R sources, nmath
  included, stay unread here and are never ported into a shared library.
- **Linear algebra**: slice 13, on go-ndarray's linear algebra once it is a
  shared library.
- **Regex `perl = TRUE`**: go-regexp/pcre2. **Compression** (`gzfile`,
  `memCompress`): the shared zlib. **Random numbers** (M4): the shared bit
  generators, as long as their state reads back as R's `.Random.seed`.
- **Files, environment, clock**: go-embedded-r has the no-host mode (the
  default in-memory tree) and the confined-root mode (`Options.FileRoot`) since
  v0.13.0. It adopts the shared host and sandbox layer
  (go-embedded-julia/julia#32) when that layer has both modes.

## 3. Gates, per slice

The gates of M2 apply unchanged:

1. goldens from our own snippets, compared with the oracle byte for byte (value,
   visibility, output, stderr, warnings, error message and call);
2. a differential fuzz mode for the new area, with **every exclusion counted**
   in the output (`evalfuzz` prints how many cases were compared through a
   weaker check, for example libm results through `round(x, 6)`);
3. `FuzzEval` for the whole evaluator, with each crasher kept in the corpus;
4. bomb tests that assert a refusal under `GOMEMLIMIT`, and the safe defaults of
   `eval.Options` (v0.4.0);
5. performance against the reference, single core, once a function is on a hot
   path (`rparity bench` style tables in the README).

## 4. Measuring coverage

- `testdata/base/exports.tsv`: the 1 362 functions exported by base, with kind
  (primitive or closure). Generated by the oracle, regenerated when R is bumped.
- `testdata/base/needs.tsv`: for each closure, the `.Internal` names it calls.
  `internals.tsv`: per internal, the number of call sites and of closures.
- **Coverage** = the fraction of base exports that the evaluator **implements
  and for which a passing golden exercises the function** (a name that appears
  as a call in `testdata/eval`, whose snippets all agree with R). It is printed
  by `go run ./cmd/rparity basecov` and recorded in the README at each tag.
  Starting point at v0.6.0: 14.4 % of base's exports agree with the oracle
  (196 of 1 362: 102 of the 210 primitives and 94 closures); 19.8 % (269)
  at v0.7.0. A function counted once is not
  uncounted when more cases are added; a refusal (`eval.Unsupported`) counts as
  not implemented.
- The metric is deliberately a floor: one passing snippet per function is the
  entry ticket, the differential fuzz modes and the base R test suite
  (run through the oracle's own snippet generator, never copied) are what make a
  function trustworthy.

## 5. Collation (decided in slice 4)

R 4.6.2 started with `LC_ALL=C.UTF-8` (the oracle's locale) still orders
strings with ICU, not by code point: `sort(c("b", "A", "a", "B"))` is
`"a" "A" "b" "B"`, punctuation sorts before digits before letters, and accents
come after the plain letter. The evaluator reproduces this for Latin text
(`eval/collate.go`: three-level comparison, ICU's order of the ASCII
punctuation), for `sort`, `order`, `ls`, `<`/`>` on strings, and `min`/`max` of
strings; `==` is exact. Strings outside Latin script compare by code point after
the Latin letters (a documented divergence). `method = "radix"` and
`LC_COLLATE=C` will select code-point order when `Sys.setlocale` exists (slice 6).
