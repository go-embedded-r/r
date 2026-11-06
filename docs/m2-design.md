# M2 design: the evaluator core

Status: the first three slices (§7, §12, §13) are implemented in package `eval` (README, "M2"). Scope: milestone M2 of
[`bib/r.md`](https://github.com/go-scikit/roadmap/blob/main/bib/r.md) (§2.3, §2.4, §5):
SEXPs, environments, promises, argument matching, closures, contexts and
conditions, S3, and the 210 primitives, in package `eval`.

## 0. Rules

- **Clean room.** R's C code is a specification, never a source: no function
  is translated, no table copied. The design below is written from the R
  Language Definition (R-lang), *R Internals* (R-ints), the help pages, and
  what the oracle (GNU R 4.6.2) does. Where this document says "R does X", X
  was either read in R-lang/R-ints or observed on the oracle; each slice is
  checked by the oracle lane before it is merged.
- **BSD-3-Clause**, and base R's R code is rewritten (decision of 2026-11-01).
  M2 needs no base R code: the first slices run on primitives alone.
- **The oracle judges.** Goldens are R's answers to inputs written in this
  repository; the evaluator's answers are compared byte for byte (see §9).
- **Bounded.** Recursion depth (R's own `expressions` limit, 5000), vector
  allocation size, and evaluation steps are limits of the session, refused
  with R's own error, never a Go crash. Evaluation of untrusted cells is a
  design constraint, not an afterthought.

## 1. Values

One Go interface, one concrete type per SEXPTYPE that user code can build:

```go
type Value interface{ Type() Type }   // Type is rds.Type: the SEXPTYPE codes

*Null                         // the singleton NULL
*Symbol                       // interned: one *Symbol per name, compared by pointer
*Logical  { v []int32;  attr }  // TRUE 1, FALSE 0, NA = math.MinInt32 (as R)
*Integer  { v []int32;  attr }  // NA = math.MinInt32
*Double   { v []float64; attr } // NA_real_ is the NaN with low word 1954 (payload kept)
*Complex  { v []complex128; attr }
*String   { v []Str;   attr }   // Str{S string; NA bool}; strings are byte strings with a UTF-8 flag
*Raw      { v []byte;  attr }
*List     { v []Value; attr }   // VECSXP; EXPRSXP is a *List with the expression flag
*Lang     { fun Value; args []Arg; attr } // a call: a function and tagged arguments
*Closure  { formals []Formal; body Value; env *Env; attr }
*Builtin  { name string; fn func(*Call) Value; kind special|builtin } // evaluated or not
*Env      // see §2
*Promise  { expr Value; env *Env; val Value; forced bool }
*Missing  // the empty argument (R_MissingArg), distinct from NULL
```

Choices, each with its reason:

1. **Typed Go slices, not `[]interface{}`.** Arithmetic loops run over `[]float64`
   and `[]int32` (the fleet's go-simd kernels can take them later); NA is a value
   inside the slice (the INT_MIN and NaN-payload conventions of R, which the
   `rds` package already round-trips), never a side bitmap.
2. **Calls are `*Lang` with an argument slice, not pairlists.** R's pairlist is a
   C linked list; R code observes it only through `as.list`, `[[`, `length`,
   `as.call`, `match.call`, `sys.call`, all of which a slice supports with the
   same results. `*Lang` is what `syntax` already produces (`syntax.Call`);
   conversion is one pass (`eval.FromSyntax`). A real pairlist type exists only
   for `formals`/`attributes`/`pairlist()` and is a separate `*Pairlist`.
3. **Attributes are an ordered slice of (symbol, value)**, nil when empty (the
   common case costs one pointer); order is observable (`attributes(x)`) and
   preserved, as in `rds`. `names`, `dim`, `class` have typed accessors.
4. **Symbols are interned in a process-wide table** (sync.Map), so lookup keys
   are pointers and `quote(x) == quote(x)` is pointer equality.
5. **`*Missing` is its own type**, because `f(x, )` and `x[1, ]` need an argument
   that is neither NULL nor a symbol, and `missing()` tests for it.
6. **Byte-code and the `compiler` package are out of scope** (bib/r.md §2.3): the
   AST walker is the semantics; a faster tier is a later, measured decision.

### Copy semantics

R copies on modification, observably (`y <- x; y[1] <- 0` leaves `x`). Go has a
GC, not reference counts. Design: every vector carries a **shared** counter
(`named`, 0..2, saturating, as R's historic NAMED): binding a value to a
variable or passing it as an argument raises it; a replacement function
(`[<-`, `names<-`, ...) updates in place when `named <= 1` and copies
otherwise. Values returned from primitives start at 0. The first slice has no
replacement functions, so every value is effectively immutable; the counter is
added with `[<-` and tested by `tracemem`-free observable behaviour only.

## 2. Environments

```go
type Env struct {
    parent *Env
    vars   map[*Symbol]*binding   // frames of a function call: a small slice first
    order  []*Symbol              // insertion order, for ls(sorted = FALSE)
    locked bool; attr *Attrs
}
type binding struct{ val Value; locked, active bool }
```

- A **function call frame** has few variables: it uses a linear slice of
  bindings up to 8 entries and switches to the map above that. Lookup walks
  `parent` pointers (`find`, `findFun` which skips non-functions, as R does for
  a call's function position).
- Distinguished environments: `Empty`, `Base` (primitives live here), `Global`
  (parent: the search path, initially just Base), and per-session. There is no
  namespace machinery in M2 (that is M3: `::` over a registry of packages).
- Locking and active bindings exist in the struct from the start (`lockBinding`,
  `makeActiveBinding` need them; the first slice only reads `locked`).

## 3. Evaluation, promises, laziness

`eval(v Value, env *Env) Value` is a switch on the type of `v`:

| `v` | result |
|---|---|
| constant (vectors, NULL) | itself |
| `*Symbol` | the binding's value: a forced promise's value, a forced-on-demand promise, error `object 'x' not found` when unbound, `argument "x" is missing, with no default` for `*Missing` |
| `*Promise` | forced once (value cached, environment dropped) |
| `*Lang` | resolve the function (`findFun` for a symbol), then by kind below |

- **Specials** (`quote`, `if`, `for`, `function`, `<-`, `{`, `&&`, `switch`, ...)
  get the unevaluated call and the environment and evaluate what they need.
- **Builtins** (`+`, `length`, `c`, ...) get their arguments evaluated, left to
  right, in the caller's environment, with `...` expanded.
- **Closures**: arguments become **promises** `{expr, env}`; the call frame is a
  new `*Env` whose parent is the closure's environment; matched formals are
  bound to promises (or to their default expressions as promises in the new
  frame); the body is evaluated in the frame. A promise is forced the first time
  its symbol is read; forcing a promise that is being forced is the error
  `promise already under evaluation: recursive default argument reference or
  earlier problems?`. `substitute`, `missing`, `force`, `delayedAssign` follow
  from the promise's `expr` and `forced` fields.
- **Visibility** is a per-session flag set by each builtin (`invisible`, `<-`,
  `for`, ... clear it; every other call sets it), read by the top level to
  decide whether to print.

### Argument matching (R-lang §Functions / Evaluation)

Three passes over the supplied named arguments and the formals, as documented:
exact name, then unique partial name (only for formals before `...`), then
positional for what is left; remaining arguments go to `...` if present, else
the error `unused argument (...)`. `...` is a `*Dots` value holding the
promises; `..1`, `...length()`, `...elt()`, and passing `...` on are slices of
that. Errors use R's own wording (`formal argument "x" matched by multiple
actual arguments`, `argument 2 matches multiple formal arguments`).

## 4. Control flow

Lexical control flow is **not** panic-based on the hot path:

- Statements (`{`, `if`, loops) are evaluated by `evalStmt`, which returns the
  value plus a `flow` code (`normal`, `break`, `next`, `return`). A function
  body, loop body and `{` propagate the code up to the loop or call that owns
  it; there is no allocation per `break`.
- `return(x)`, `break`, `next` reached **through a promise or a non-statement
  position** (`f(return(1))`, `lapply(x, function(i) break)`) and every
  non-lexical exit (errors, `invokeRestart`, `stop`, `tryCatch` unwinding) use a
  Go panic with a typed value, recovered at the owning context (function call,
  loop, `tryCatch`, `withRestarts`) and nowhere else. Panics never cross a
  wasm callback boundary (bib/r.md §2.3): the session entry point recovers.
- **Contexts**: a stack of frames `{call *Lang, fn *Closure, env *Env, promargs,
  onexit []Value, sysparent *Env}` backs `sys.call`, `sys.function`,
  `parent.frame`, `match.call`, `on.exit`, `Recall`, `nargs`, `UseMethod`
  (dispatch looks at the caller frame), and the depth limit.

## 5. Conditions

R's system (R-lang §Exception handling), in Go:

- A **handler stack** (`[]handlerEntry{classes, fn, calling bool, target}`) per
  session, pushed by `withCallingHandlers` (calling) and `tryCatch` (exiting).
- `signalCondition`/`stop`/`warning`/`message` build a condition object (a
  list of class `c("simpleError","error","condition")` with `message` and
  `call`) and walk the stack from the top: a **calling** handler is run in place,
  on a stack truncated below it; an **exiting** handler panics with a `unwind`
  value naming its `tryCatch` target, recovered there, which then runs the
  handler. With no handler `stop` reaches the top level: the session returns an
  `*eval.Error` carrying R's message `Error in f(x) : msg` formatting, the call,
  and the condition object.
- **Restarts** (`withRestarts`, `invokeRestart("muffleWarning")`) are entries on
  a second stack; `warning()` establishes `muffleWarning`, `stop()` none.
- Deferred **warnings** accumulate per top-level call and are returned with the
  result (`Result.Warnings`), printed by the notebook, as R prints them after
  the call completes.
- Errors raised by the runtime itself (type errors, `object not found`,
  `non-numeric argument to binary operator`) are conditions with R's message
  text and call, so `tryCatch` sees what it sees in R.

## 6. S3 and the rest of dispatch (slices after the first)

`UseMethod(generic, object)`: find `generic.class` for each element of the
implicit class (`class(x)`, then `matrix`/`array`/`integer`/`numeric`/`function`
defaults per R-lang), then `generic.default`; the method runs with the generic's
promises (not re-evaluated), `.Generic`, `.Class`, and `NextMethod` state. Group
generics `Ops`/`Math`/`Summary`/`Complex` dispatch from the arithmetic
primitives when an operand has a class attribute. S4 (`standardGeneric`) is M4.

## 7. The first slice

Constants and symbols; `*Env` with `Global`/`Base`/`Empty`; `<-`, `=`, `<<-`
(top level only), `->` (a parse-time rewrite already), `function`, closure calls
with full argument matching and promises, `{`, `(`, `if`, `for`, `while`,
`repeat`, `break`, `next`, `return`, `quote`, `invisible`, `missing`,
`identity`-like closures written by the tests; the arithmetic primitives `+ - *
/ ^ %% %/%` (unary and binary), comparison `== != < <= > >=`, logic `! & | && ||`,
`:`, `c`, `length`, `is.na`, and the type predicates, with R's **recycling**
(length of the result is the longer; zero-length gives zero-length; a warning
`longer object length is not a multiple of shorter object length`), **NA rules**
(NA propagation, `NA & FALSE` is `FALSE`, `NA_integer_` arithmetic, integer
overflow gives `NA` with the warning `NAs produced by integer overflow`), and
**type promotion** (logical < integer < double < complex, `/` and `^` always
double, `%/%` and `%%` on integers stay integer). Results are compared with the
oracle on goldens generated from our own snippets (§9).

Implemented in the first slice, 71 primitives: `quote function { ( invisible if while repeat for break next return <- = <<- missing + - * / ^ %% %/% == != < <= > >= & | ! && || : c list length is.na is.null is.logical is.integer is.double is.numeric is.character is.list is.function is.symbol is.name is.call is.environment is.atomic as.integer as.double as.numeric as.logical as.character abs sqrt exp floor ceiling sum prod min max any all seq_len seq_along names`.

## 8. Order of the primitives

Counted by `go run ./cmd/rparity primcount <dir>` (call sites of each primitive
in the R files of a tree, plus references as a value, plus files; it parses with
`syntax` and keeps nothing of the sources). Over the 163 files of
`src/library/base/R` of R 4.6.2 (23 000 lines), the 210 primitives are used as
follows; the first column is the milestone slice.

| slice | primitives (call sites in base R) |
|---|---|
| 1 (this PR series) | `<-` 6257, `if` 3195, `{` 2130, `(` 312, `return` 305, `for` 156, `while` 8, `repeat` 3, `break` 8, `next` 17, `function`, `quote` 17, `missing` 200, `invisible` 80, `!` 1112, `&&` 593, `\|\|` 325, `==` 528, `!=` 216, `>` 287, `<` 138, `<=` 48, `>=` 36, `-` 418, `+` 214, `*` 144, `/` 60, `^` 9, `%%` 37, `%/%` 13, `&` 32, `\|` 19, `:` 59, `c` 543, `length` 885, `is.na` 262, `is.null` 407, `is.character` 174, `is.logical`, `is.numeric`, `is.integer`, `is.double`, `is.function`, `as.integer`/`as.character`/`as.numeric`/`as.double`/`as.logical`, `seq_len` 136, `seq_along` 75, `any` 186, `all` 84, `sum` 54, `min`/`max`/`prod`/`range`, `abs`, `sqrt`, `exp`, `floor`, `ceiling`, `round` |
| 2 | subsetting and its replacement forms: `[` 1426, `[[` 520, `$` 572, `@`, `[<-`, `[[<-`, `$<-`, `<<-` (35), `names` 402, `names<-`, `attr` 216, `attr<-`, `attributes` 47, `class` 98, `class<-`, `oldClass` 73, `unclass` 87, `inherits`-support, `dim` 135, `dimnames` 76, `list` 210, `unlist`-support, `rep` 34, `seq.int` 33, `is.list`, `is.matrix`, `is.array`, `levels<-`, `length<-`, `environment`, `environment<-` |
| 3 | calling and context: `UseMethod` 97, `NextMethod`, `standardGeneric`, `switch` 55, `substitute` 51, `on.exit` 63, `nargs` 24, `...length`, `...elt`, `...names`, `interactive`, `sys.call`-family (closures over `.Internal`), `forceAndCall`, `Recall`, `.Internal` 519 (the dispatch to the 563 internals), `::`/`:::` 58/23, `emptyenv`, `baseenv`, `as.environment`, `isS4`, `.isMethodsDispatchOn` |
| 4 | math and summaries: `cumsum`, `cumprod`, `cummax`, `cummin`, `log`, `log10`, `log2`, `exp`, `sin`..., `trunc`, `sign`, `Re`/`Im`/`Mod`/`Arg`/`Conj`, `is.finite`/`is.infinite`/`is.nan`, `anyNA` 65, `nzchar` 88, `max`/`min` on character, `xtfrm`, bit ops |
| 5 | the foreign and low-level rest: `.Call`, `.Fortran`, `.External*`, `.C` (refused or routed to Go ports), `lazyLoadDBfetch`, `proc.time`, `gc.time`, `Sys.*`-support |

`.Internal` is the larger half: base R's closures call **412 distinct internals
at 1 to 7 sites each** (most once; leaders `as.vector` 7, `vector` 6, `unlist`
4, `qsort` 4, `.signalCondition` 4), so frequency does not rank them; they are
ordered by the closures that matter first (M3): `vector`, `as.vector`, `unlist`,
`paste`, `sprintf`/`format`-family, `identical`, `match`, `unique`, the
condition internals (`.addCondHands`, `.signalCondition`, `.invokeRestart`,
`.dfltStop`), `environment` and `sys.*` internals, `lapply`/`vapply`/`mapply`,
`rev`/`sort`/`order`, string internals, then the rest as the base test suite
reaches them.

The full table (all 210 primitives, with references and files) and the
`.Internal` table are the output of the command above; they are regenerated, not
maintained by hand.

## 9. Testing: the oracle's `eval` lane

`internal/oracle`'s `eval` operation evaluates source in a fresh child of R's
global environment and returns the value (serialized), the captured output, the
warnings, and the condition class and message on error. For M2:

- `testdata/eval/*.R` are snippets written here (one expression or a small
  program per `#---` block). `roracle eval-golden` writes
  `testdata/golden/eval.jsonl`: for each snippet the value as `rds` bytes (hex)
  with `typeof`, `length`, `class`, the visibility, the warnings, and the error
  message when it fails.
- `eval` package tests replay the goldens **without R**: evaluate each snippet,
  serialize the value with `rds.Encode`, and compare bytes, warnings and error
  text. With R, the same comparison runs live (CI oracle lane), and the goldens
  are regenerated and must not change (the platform rule of the README: inputs
  must not depend on `long double` or FMA).
- Equality is exact for logical, integer, character, list structure, attributes
  and their order; doubles are compared by bits where R computes in `double`
  (arithmetic, comparison), by the stated tolerance elsewhere (math functions
  through go-numerics, M4).
- Differential fuzz: random expressions over a small grammar are evaluated by R
  and by Go (the oracle's `eval` over a batch), values and error messages
  compared; this is what finds recycling and NA corner cases.
- Unit tests of the runtime itself (environments, matching, promises,
  conditions) are plain Go tests with goldens only where R's wording matters.

## 10. API for go-notebook

```go
s := eval.NewSession(eval.Options{MaxDepth: 5000, MaxVector: 1 << 28, MaxSteps: ...})
res, err := s.Eval(src)       // parse (package syntax) + evaluate each expression in order
res.Value     eval.Value      // the last value
res.Visible   bool            // whether R would print it
res.Output    string          // text written by print/cat/message while running
res.Warnings  []string
// err is *eval.Error (R's message, call, condition) or *syntax.Error / *syntax.LimitError
s.Get("x"), s.Set("x", v)     // names assigned by one cell are seen by the next
```

One `Session` is one persistent global environment (the notebook's shared
interpreter); `TopLevelNames` (package `syntax`) tells the notebook which names
a cell defines. The session is not safe for concurrent use; cells run one at a
time. `Value` converts to and from Go (`float64`, `[]float64`, `string`, ...) in
`eval/bridge.go`, following the value bridge of bib/r.md §5 (M6).

## 12. Slice 2 (implemented)

Replacement functions with always-copy semantics (the ownership optimisation of §1 stays deferred), subsetting, attributes, S3 dispatch, the condition system, value printing and the `Eval` API with limits (`Timeout`, `MaxAlloc`, `MaxSteps`, `MaxDepth`, `MaxVector`; a cancelled context). Handlers are a stack on the session: exiting handlers unwind with a typed panic caught by their `tryCatch`, calling handlers run with the stack truncated below them, restarts are found by name. Errors are signalled as conditions before they end the evaluation, so `tryCatch(error=)` and `withCallingHandlers(error=)` see them. The oracle's `evalcase` captures printed output and standard error with `sink`, auto-prints visible values, and expands ALTREP (compact sequences) by a version-2 serialize round trip, nested lists included. The differential fuzz has two modes: `arith` and `subset`.

## 13. Slice 3 (implemented)

Environments are values (`*Env`; `Session.Empty` is the root above base); `parent.frame(n)` and `sys.call` read the frame stack by environment, not by position, so they work from a builtin called inside any closure. The apply family is native Go that calls `callWith`, the one place that applies a function to already evaluated, possibly named arguments; `sapply` and `mapply` simplify as `simplify2array` does for vectors and matrices. Regular expressions are described in the README: the matcher is github.com/go-embedded-r/tre, a Go port of TRE with R's modifications (parallel matcher, backtracking matcher for back-references, resource limits and a step limit that refuse hostile patterns), with R's own loops for `gsub`, `gregexpr` and `strsplit` written on top. The differential fuzz has five modes: arith, subset, regex, strings, apply.

## 11. Risks

- The **slice sizes**: the first slice is small on purpose; every later slice is
  one PR with oracle goldens.
- **Printing** (`print`, `format`, auto-print) is a large, exact piece of R
  (`print.default`, `formatReal`); it comes with M3 and is judged by the same
  oracle; the first slice returns values, and the oracle lane compares values,
  not printed text.
- **R's numeric printing digits** (`scientific()`) is already ported for
  `deparse` (`syntax/deparse.go`) and is reused.
- **Performance** is measured against R single-core per slice (fleet rule) once
  there are loops to measure; the first number is for `for` + arithmetic.
