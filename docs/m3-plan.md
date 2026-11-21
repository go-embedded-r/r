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
| 13 (done in slice 20, v0.22.0: see row 20) | linear algebra that base exports: `solve`, `qr`, `chol`, `svd`, `eigen`, `det` | it waited for the extraction of go-ndarray's linear algebra into a shared library (go-embedded-ruby/ruby#918), not for go-numerics; go-ndarray v0.44.0 (`linalg`) is that library |
| 14 (v0.16.0, done: see README) | data frames and `table`, which slices 11 and 12 deferred: `data.frame`, `as.data.frame` (vector, list, matrix, table), `is.data.frame`, `[`/`[[`/`$` and their replacement forms, `dim`/`dimnames`/`row.names`, `rbind`/`cbind` of data frames and of vectors and matrices, `merge`, `split` of a data frame, `subset`, `transform`, `droplevels`, `print`/`format`/`as.matrix`/`t` of a data frame, `table`, `print.table`, `as.data.frame.table`, `expand.grid`, `interaction`, `cut`; goldens sharded per snippet file. **Deferred**: `summary` of a data frame or a table, the `Ops`/`Math`/`Summary` methods of data frames, `aggregate`, `by`, `stack`, `within`, `ftable`, `xtabs`, matrix and list columns; `head` and `tail` are `utils`, outside base | the notebook idioms (`data.frame`, `merge`, `table`, `cut`); 48.4 % of base's exports (659) |
| 15 (v0.17.0, done: see README) | the deferred base functions that need no shared engine: `summary` (default, factor, matrix, data frame) and `print.summaryDefault`, `format.pval`, the `Ops`/`Math`/`Summary` methods of factors and data frames, `within`, `all.equal` and `isSymmetric`, function objects (`formals`, `body<-`, `as.function`), `hexmode`/`octmode`, `AsIs`, the `diff`, `range`, `split`, `length<-`, `xtfrm` and `mtfrm` methods of dates and times, the `.default` methods of the generic natives, condition methods, `.Deprecated`/`.Defunct`. **Deferred**: `strwrap`, `pretty`, `kronecker`, `dput`/`dget`/`dump`, `print.rle`, `seq.POSIXt`, `round`/`trunc`/`cut` of dates, `summary` of dates, times, differences and tables (`pchisq`), `by`, `file.info`, `scan`, `numeric_version`, `Vectorize` (untestable through the oracle's guard) | the idioms around `summary(df)`, `all.equal`, `within`, `df * 2`; 57.4 % of base's exports (782) |
| 16 (v0.18.0, done: see README) | what slice 15 deferred, in the order of `rparity basecov`: `numeric_version`, `package_version`, `R_system_version`, `getRversion`, `R.version`, `compareVersion`; `seq.POSIXt`, `round`, `trunc` and `cut` of dates and times, `summary` of dates, times and differences; `strwrap`, `pretty`, `kronecker` and `%x%`, `print.rle`; `dput`, `dget`, `dump` and the `deparse` options (`width.cutoff`, `control`, `nlines`, `backtick`); `by`; `file.info`, `file.size`, `file.mtime`, `file.mode` over the session's file system; the marks of `format`, `formatC` and `prettyNum` (`big.interval`, `small.mark`, `decimal.mark`, `zero.print`, `drop0trailing`); `unique`, `duplicated` and `anyDuplicated` of matrices and arrays, `as.table`, `data.matrix`, `list2DF`, `sort_by`, `provideDimnames`, `dynGet`, `system.time`, `transform` and `subset` of other classes, `mode<-`, `logb`, `single`. **Deferred**: `scan`, `read.table` and the other readers of utils, `summary` of tables (`pchisq`), `sort_by` of a formula and every other use of `~` (formulas), `Vectorize` and the direct call `kronecker()` (the oracle refuses both, see the README), `chkDots`, `warnings`, `rowsum`, `jitter` and `sample` (random numbers, M4), the `[`, `[[` and `$<-` methods of `POSIXlt`, `str` (only the line `print.rle` needs) | the idioms around dates (`cut`, `trunc`, `seq` by month), version checks, `dput` and `strwrap`, `by`; 63.8 % of base's exports (869) |
| 17 (v0.19.0, done: see README) | formulas and what they allow: the `~` primitive (a call that carries `class` and `.Environment`), `print`, `format` and `deparse` of a formula, `terms.formula` (`+`, `-`, `:`, `*`, `/`, `^`, `%in%`, `.` against data, `I()`, `offset()`, `specials`, `keep.order`, `simplify`), `all.vars`, `all.names`, `update.formula`, `reformulate`, `as.formula`, `formula`, `model.frame.default` (`subset`, `na.action`, `drop.unused.levels`), `model.response`, `get_all_vars`, `na.omit` and its relatives, `complete.cases`; `aggregate` (default, data frame, formula), `stack`, `unstack`, `sort_by` of a formula; without formulas: `str`, `rowsum`, `warnings` and `print.warnings`, `Vectorize`, `chkDots`, `array2DF`, the `[`, `[[`, `[<-`, `[[<-` and `$<-` methods of `POSIXlt`, the remaining `as.data.frame` methods, `registerS3method`. **Deferred**: `scan`, `read.table` and its relatives, `xtabs`, `ftable`, `model.matrix` and the contrasts (slice 13), `summary` of tables, the statistics of stats (`median`, `var`, `quantile`), `serialize`, `summary.warnings` | the idioms around `y ~ x` and `aggregate(y ~ g, d, f)`, `str`, `stack`; 66.4 % of base's exports (905) |
| 18 (v0.20.0, done: see README) | tables in and out, and statistics: `scan`, `read.table`, `read.csv`, `read.csv2`, `read.delim`, `read.delim2`, `count.fields`, `type.convert`, `write.table`, `write.csv`, `write.csv2`; `xtabs`, `summary.table` (no chi-squared test); `median`, `var`, `sd`, `cov`, `cor` (three methods), `quantile` (nine types), `fivenum`, `IQR`, `mad`, `weighted.mean`. Written here, each with a stated error bound against math/big, before the binding to go-statistics/stats (slice 19) replaced the parts of them it could. **Deferred**: `ftable`, `ecdf`, `read.fwf`, `pchisq` and the independence test of `summary.table` (done in slice 19), the timing against R (go-embedded-r/r#48), `model.matrix` and the contrasts (slice 13) | the notebook idioms `read.csv`, `write.csv`, `summary(df$x)`, `quantile`, `cor`; 66.7 % of base's exports (908) |
| 19 (v0.21.0, done: see README) | the statistics move onto [go-statistics/stats](https://github.com/go-statistics/stats) v0.1.1 where it holds the stated bound of slice 18: the ranks of `cor(method = "spearman")` and the tau-b of `cor(method = "kendall")` (O(n log n), where the slice-18 kernel was O(n^2)); `pchisq` and `qchisq` (central and non-central, `lower.tail`, `log.p`) and the chi-squared test of independence of `summary.table`. The kernels go-statistics does not hold to the bound stay here, listed as gaps in the README. **Deferred**: `ftable`, `ecdf`, `read.fwf`, `model.matrix` and the contrasts (slice 13, done in slice 20), the other distribution functions | `summary(table(a, b))`, `cor(x, y, method = "kendall")` on large vectors |
| 20 (v0.22.0, done: see README; it is row 13) | linear algebra bound to [go-ndarray](https://github.com/go-ndarray/ndarray) v0.44.0 (`linalg`): `solve` (default, `qr`), `qr` (LINPACK `dqrdc2` and `LAPACK = TRUE`) with `qr.Q`, `qr.R`, `qr.X`, `qr.coef`, `qr.qy`, `qr.qty`, `qr.fitted`, `qr.resid`, `qr.solve`, `is.qr`, `as.qr`; `chol` (also `pivot = TRUE`), `chol2inv`, `det`, `determinant`, `svd`, `La.svd`, `eigen` (symmetric, general, complex results) and `print.eigen`, `backsolve`, `forwardsolve`, `kappa`, `rcond`, `norm`; `contr.treatment`, `contr.sum`, `contr.helmert`, `contr.poly`, `contr.SAS`, `contrasts`, `contrasts<-`, `model.matrix` on the `terms` of slice 17. Each states a bound against math/big; what go-ndarray lacks (pivoted and LINPACK QR, triangular solves, pivoted Cholesky, the LU factors for `rcond`) is written here and listed in the README as gaps. **Deferred**: `qr`, `svd` and `chol` of complex matrices, `poly` (so `model.matrix(~ poly(x, 2))`), the timing against R (go-embedded-r/r#48) | `solve`, `qr`, `lm`-style design matrices through `model.matrix`; 69.2 % of base's exports (943) |
| 21 (v0.23.0, done: see README) | the deferred stretch of slices 19 and 20: `ftable` (`ftable.default`, `ftable.formula`, `format`, `print`, `as.table`, `as.matrix`, `as.data.frame`), `ecdf` and its methods (`summary`, `quantile`, `knots`, `print`), `poly` and `predict.poly`, `read.fwf`, `formatDL`, the `srcfile` family, and the methods of `factor`, `noquote`, `table`, `listof` and `simple.list`; `.rowSums` and its kin; `regmatches(invert = NA)`, `is.pairlist` of a pairlist. Bounds: `ecdf` within u = 2^-53 and `poly` orthogonality within 64 n u, both against math/big. **Deferred**: package, namespace and autoload machinery, `sample`/`jitter`, `trace`/`untrace`, `plot`, the timing against R | `ecdf(x)(q)`, `poly(x, 3)`, `ftable(Titanic)`, `read.fwf`; 75.8 % of base's exports (1033) |
| 22 (v0.24.0, done: see README) | the small closures and natives that need no host: `charmatch`, `char.expand`, `findInterval`, `choose` (correctly rounded to k = 1000), `numToBits`, `numToInts`, `encodeString`, `format.info`, `sort.int`, `sort.list`, `rapply`, `asplit`, `drop`, `max.col`, `callCC`, `args`, `comment`, `.mapply`, `.makeMessage`, the assignment methods of `POSIXct` and `difftime`, and the factor, data.frame and ts conversions; the `try()` line wrap of R. Bound: `choose` within u = 2^-53 of math/big. **Deferred**: `lchoose`, `abbreviate`, `balancePOSIXlt`, `.invokeRestart`, the locale family, the debugger family, the timing against R | `choose(60, 30)`, `findInterval`, `encodeString`, `asplit`, `callCC`; 79.1 % of base's exports (1077) |

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
- **Statistics** (descriptive statistics, tests, the distributions): bind
  [go-statistics/stats](https://github.com/go-statistics/stats) (v0.1.1 since
  slice 19), and later [go-statistics/models](https://github.com/go-statistics/models)
  for `lm`, `glm` and `anova`; not go-numerics, which has no statistics kernel.
  The R layer keeps the glue (NA, `na.rm`, `use=`, `names`, recycling) and
  binds a function only where go-statistics meets the bound this repository
  states against math/big (the test is in `eval/stats18_bind_test.go`); where
  it does not, the kernel stays here and the gap is listed in the README, for
  the owner of go-statistics to close.
- **What R does better goes upstream.** When R's reference has an algorithm, a
  default or a test problem that go-numerics lacks, it is added to go-numerics
  for the whole fleet, not to the R layer. This applies only to what is
  available outside R's own GPL text: published papers and netlib's public
  domain or BSD code (QUADPACK, MINPACK, PORT). R's own C and R sources, nmath
  included, stay unread here and are never ported into a shared library.
- **Linear algebra** (slice 13, done in slice 20): binds go-ndarray's `linalg`
  (v0.44.0): `Solve`, `Inv`, `Det`, `Slogdet`, `Cholesky`, `SVD`, `Eigh`, `Eig`.
  What it lacks is written here (README, slice 20), for its owner to close.
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
