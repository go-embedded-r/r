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
| 9 | dates and times: `Sys.time`, `Date`, `POSIXct`, `POSIXlt`, `difftime`, `format`/`strptime`/`strftime` (TZ via Go's time database) | self-contained, many exports |
| 10 | connections and files: `file`, `readLines`, `writeLines`, `textConnection`, `readRDS`/`saveRDS` over `rds`, `scan`, `read.table` family in `utils` later | the notebook bridge needs them |
| 11 | linear algebra and numerics that base exports: `solve`, `qr`, `chol`, `svd`, `eigen`, `det` (the LAPACK routines are ported to Go by M4's numeric work; until then refused) | depends on go-numerics |
| 12 | the complex and raw types, expression vectors, formulas, S4 | each needs a value type; they come with the stats and methods packages |

A slice is tagged when its goldens, fuzz, `FuzzEval` and bomb tests are green and
`rparity basecov` has moved.

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
