# M4 plan: numerics and the random streams of R

Status: slice 26 done (v0.28.0, 82.1 % of base; see README); slices 27 to 31 planned. M3 ended at slice 23 (v0.25.0, 81.1 % of base).

M4 brings the numerical half of R: the random number streams, the special
functions of base, the d/p/q/r functions of `stats`, and the optimisers and
quadrature. The evaluator does not grow its own numerics: it binds the fleet
libraries (go-numerics, go-statistics, go-ndarray), and what they lack goes
**upstream as pull requests**, so that every binding (Ruby, Julia, C, R) gets
it. This document gives the order of the work, the contract each slice is held
to, and the gaps.

## 1. Rules

### Clean room (unchanged from M3)

- **Specification**: R's help pages (rendered from the `.Rd` files by the
  oracle) and the **observed behaviour of the oracle** (R 4.6.2, through
  `internal/oracle`). R's `.R` and `.c` files, `src/nmath` included, are never
  read.
- **Algorithms**: from the papers that the help pages cite, from netlib
  public-domain or BSD code, and from the APIs of the fleet libraries. A
  reference is cited only once it has been opened.
- The go-numerics roadmaps cite R source lines that their authors read to
  understand an algorithm. This plan cites those roadmaps as fleet documents
  only and never follows their R line numbers.

### Contract: a stated bound, not bit identity

The fleet rule (2026-11-11) applies: bit identity with R is a porting step,
never the goal. Each slice states an **error bound** and the **oracle** that
asserts it:

- **mpmath** at 192 and 320 bits, stored as a double-double, for the special
  and distribution functions. This is the contract of go-numerics/special
  (`docs/accuracy.md`) and of go-statistics/stats (`docs/accuracy.md`); M4
  adopts those budgets rather than inventing its own. special states them in
  its README "Accuracy contract" (go-numerics/special#69): per-region budgets
  of 1.5 times the measured maximum, and four `BetaIncInv`/`BetaIncCInv`
  regions (a or b < 1e-3) unbounded today (go-numerics/special#11).
- **`math/big`** for roots and minimisers (the root or minimum is recomputed at
  extended precision; the bound is on |x − x\*| and on f). This is
  go-numerics/optimize's contract (go-numerics/optimize#61: exact optima and a
  300-bit `math/big` oracle, `TestAccuracyBound…` per solver family).
- **mpmath `quad`** at 50+ digits for quadrature (true error within the
  requested tolerance, and the reported error bounds the true error).

R's own value is always recorded beside ours in the golden. Where R is
farther from the oracle than we are, the golden says so, and the
difference is counted as an exclusion in the fuzz output, never hidden.

The one exception is the **random streams** (section 3): a seeded stream is a
definition, not an approximation, so state equality is exact.

## 2. Order

| slice | R functions | binding | R layer writes | bound and oracle | why here |
|---|---|---|---|---|---|
| 24 (v0.26.0, done: see README) | `.Random.seed`, `set.seed`, `RNGkind`, `RNGversion`, `runif`, `sample`, `sample.int`, `rnorm` (Inversion), `jitter` | go-ndarray/ndarray v0.46.0 `random.MT19937` (`State`, `SetState`, `Uint32`) | `set.seed`'s seed scrambling, the 32-bit-to-double map, `.Random.seed[1]` kind code, the Rejection and Rounding samplers, the Walker alias method (?sample: replace, more than 200 probable values; Ripley 1987), the two-draw path for n ≥ 2^31 | `.Random.seed` equal to the oracle's after every call (exact); `runif` exact; `rnorm` within 1 ulp of mpmath's Φ⁻¹ at the uniform actually drawn | every later r\* function and every simulation in a test needs it, and its state is fully observable (section 3) |
| 25 (v0.27.0, done: see README) | `gamma`, `lgamma`, `beta`, `lbeta`, `digamma`, `trigamma`, `psigamma`, `choose`/`lchoose` (completing slice 23), `besselI`/`J`/`K`/`Y` with `expon.scaled`, the last bits of `exp` and `^` (README "Known divergences") | go-numerics/special v0.20.2 (`Gamma`, `LogGamma`, `LogGammaSign`, `Beta`, `LogBeta`, `Digamma`, `Polygamma`, `Iv`, `Ive`, `Jv`, `Kv`, `Kve`, `Yv`) | recycling, `NA`/`NaN` and warning rules, `nu < 0` reflection for `besselI`/`besselK` as ?Bessel states, `lgamma`'s sign attribute | special's measured worst cases (`docs/accuracy-table.md`): gamma 0.50 ulp, loggamma 1.00, beta 1.00, polygamma 0.50, iv 2.76, jv 3.77, kv 1.89, yv 2.66; logbeta (19 ulp near log B = 0), digamma (39 ulp near negative roots) and jv at half-integer order (10 210 ulp) are listed as divergences until the upstream fix | needed by the distributions of 26–27; one binding, small R layer; closes #14 (`lgamma`, `digamma`, `trigamma` less accurate than R's, found by `evalfuzz --mode math`) |
| 26 (v0.28.0, done: see README) | d/p/q of the continuous families: `norm`, `unif`, `exp`, `gamma`, `beta`, `t`, `f`, `chisq`, `lnorm`, `logis`, `cauchy`, `weibull`; arguments `log`, `lower.tail`, `log.p` | go-statistics/stats v0.1.1 (`PDF`, `LogPDF`, `CDF`, `LogCDF`, `SF`, `LogSF`, `PPF`, `ISF`) | parameter recycling, the `lower.tail`/`log.p` dispatch onto CDF/SF/LogCDF/LogSF, quantiles at log probabilities, R's parameterisations (`rate` vs `scale`, `ncp` argument presence) and boundary values (`sd = 0` point mass from ?Normal, `min == max` from ?Uniform) | stats' budgets (`docs/accuracy.md`: normal ≤ 0.55 ulp, gamma family ≤ 2.1, beta cdf/sf ≤ 2.6, t ≤ 2.4 and quantiles ≤ 9.4 below 1e-20, F ≤ 2.3); the log-probability quantiles are measured against mpmath before they are claimed, because stats has no `PPF` at a log probability yet | the largest group of `stats` exports, and the most called |
| 27 | d/p/q of the discrete and rank families: `binom`, `pois`, `geom`, `nbinom` (`size`/`prob` and `mu`), `hyper`; noncentral `chisq`, `t`, `beta`, `f`; `ptukey`/`qtukey`; `wilcox`, `signrank` | go-statistics/stats v0.1.1 (Binomial, Poisson, Geometric, NegativeBinomial, Hypergeometric, NoncentralChiSquared) | the `mu` parameterisation of `nbinom`, integer-valued checks and warnings, `ncp` dispatch (noncentral χ² with `df = 0` is go-statistics/stats#26, open) | binomial/negative binomial cdf/sf ≤ 3.5 ulp; noncentral χ² pdf ≤ 9.8, cdf ≤ 116, sf ≤ 151 ulp (stats `docs/accuracy.md`); noncentral t/beta/F, Tukey, Wilcoxon and signed-rank have no fleet implementation and wait for the upstream PRs (section 4); until then they are `eval.Unsupported`, not approximated | needs 26's dispatch; the missing families are the main upstream work of M4 |
| 28 | `rbinom`, `rpois`, `rgamma`, `rbeta`, `rhyper`, `rexp`, `rweibull`, `rgeom`, `rnbinom`, `rchisq`, `rt`, `rf`, `rcauchy`, `rlogis`, `rlnorm`, `rmultinom`, `rsignrank`, `rwilcox` | the samplers of section 4 (go-statistics), driven by slice 24's uniform source | argument recycling and the draw order across a vector of parameters | streams: the `.Random.seed` after each call equals the oracle's (exact: same draws consumed); values within the bound of the transform they apply (for example, rweibull by inversion within 1 ulp of the exact inverse at the drawn uniform) | needs 24 and 26; R's streams follow R's algorithms, not NumPy's, so stats' `Rvs` cannot be reused |
| 29 | `optimize`/`optimise`, `uniroot` (with `extendInt`), `polyroot` | go-numerics/optimize v0.9.0 (`Fminbound`, `Brentq`) | R's tolerance conventions (`tol = .Machine$double.eps^0.25`), the result lists (`root`, `f.root`, `iter`, `init.it`, `estim.prec`), `extendInt` until upstream D2 lands | root: the final bracket contains a `math/big` root, and its width is within R's documented `tol`; minimiser: within eps·abs(x) + tol/3 of the `math/big` minimiser (?optimize); `polyroot`: each root within a stated relative bound of the `math/big` (or mpmath `polyroots`) root, the bound set after measurement | Brent (1973) in both (?optimize, ?uniroot cite it; ?uniroot names netlib `c/brent.shar` `zeroin.c`); small and much used inside other R code |
| 30 | `optim` (Nelder-Mead, BFGS, CG, L-BFGS-B, SANN, Brent), `optimHess`, `constrOptim`, `nlm`, `nlminb` | go-numerics/optimize v0.9.0 (`Minimize` with `NelderMead`, `BFGS`, `CG`, `L-BFGS-B`; `ApproxFprime`) | `control` (`fnscale`, `parscale`, `ndeps`, `maxit`, `abstol`, `reltol`, `REPORT`, `trace`), the result list (`par`, `value`, `counts`, `convergence`, `message`, `hessian`), SANN's cooling `temp / log(((t-1) %/% tmax)*tmax + exp(1))` (?optim) | minimiser within ‖x − x\*‖ ≤ 1e-8·max(1, ‖x\*‖) and f within 1e-12 relative of the known optimum refined in double-double, on the MGH problems (the bound of optimize's roadmap); R's `counts` recorded, not asserted, where the method differs (section 4) | uses 29 for `method = "Brent"`; the optimisers in optimize today follow SciPy, so this slice is where most divergence from R's iterates is expected |
| 31 | `integrate` | go-numerics/integrate v0.3.0 (`Quad`) | `rel.tol` default `.Machine$double.eps^0.25`, the check `rel.tol ≥ max(50·eps, 0.5e-28)` (?integrate), `subdivisions`, the `integrate` object and its `print` | true error within the requested tolerance against mpmath `quad`, and `abs.error` bounds it (integrate's roadmap oracle rule) | ?integrate: QUADPACK `dqags`/`dqagi` (Piessens et al. 1983), which `Quad` provides; R calls `f` with a **vector** of nodes, see the upstream batch integrand |

`nls`, `deriv`, `splinefun`, `approxfun`, `fft` and `convolve` were read in
their help pages but belong with the modelling work of M5 (`nls` uses PORT
through `nlminb` and variable projection, optimize roadmap item 8) or need a
library not inventoried here (`fft`).

Timing against R, single core, is part of each slice's gate and closes
go-embedded-r/r#48 for the functions of M4.

## 3. `.Random.seed` compatibility

R 3.6.0 and later default to `RNGkind("Mersenne-Twister", "Inversion",
"Rejection")` (?RNGkind). The facts below come from ?RNGkind, ?sample and
?Normal; everything the pages leave open is pinned by an oracle experiment
before the code is written.

- **Mersenne-Twister**: Matsumoto & Nishimura (1998), "code updated in 2002".
  The tempering and the 624-word recurrence are the published algorithm, which
  go-ndarray's `MT19937.Uint32` already implements (`random/mt19937.go:108-136`).
- **The seed vector**: `.Random.seed` is 626 integers: element 1 codes the kinds
  (units and tens: the uniform kind; hundreds: the normal kind; ten-thousands:
  the sample kind), element 2 is the position in the state, elements 3–626 are
  the 624 words, read as unsigned (so they can be negative or `NA` in R). Our
  `.Random.seed` is a plain integer vector in the global environment; reading
  it calls `State()`, assigning it calls `SetState` after validation.
- **`set.seed`**: R uses "its own initialization method due to B. D. Ripley"
  (?RNGkind), not Matsumoto–Nishimura's `init_genrand` and not NumPy's legacy
  seeding (`SeedLegacy`). The method is pinned by experiment: `set.seed(s)` for
  a few hundred `s` (0, 1, negative, `.Machine$integer.max`, `NA`), with
  `.Random.seed` dumped after each, is the golden; the scrambling is then
  written to reproduce the 626 integers exactly.
- **The uniform map**: ?RNGkind says the generators "return 32-bit integer
  values that are converted to doubles", and ?Uniform that `runif` never
  returns 0 or 1. The exact map (scaling constant, the guard at the ends) is
  pinned by drawing after `set.seed` and comparing each double with the 32-bit
  word read from `.Random.seed`.
- **Inversion**: `rnorm` applies the normal quantile to uniform draws (?RNGkind
  refers to ?qnorm: Wichura 1988, AS 241). How many 32-bit words each normal
  consumes is read from the change of `.Random.seed[2]` across `rnorm(1)`,
  not assumed. The quantile is go-statistics' `Normal.PPF` (≤ 0.55 ulp);
  where R's AS 241 is farther from mpmath than ours, the value differs in the
  last bit and the golden counts it.
- **Rejection**: `sample.kind = "Rejection"` (R 3.6.0, PR#17494) replaced
  "Rounding". The number of words drawn per index, and the rejection rule, are
  pinned the same way, from `.Random.seed[2]`, for n just below and above
  powers of two and for n ≥ 2^31 ("two random numbers are used", ?sample).
  "Rounding" is kept for `RNGversion("3.5.0")`.
- **Walker alias** and sequential sampling without replacement (?sample),
  from Ripley (1987).
- **Other kinds**: Wichmann-Hill, Marsaglia-Multicarry, Super-Duper,
  Knuth-TAOCP(-2002), L'Ecuyer-CMRG and Box-Muller / Kinderman-Ramage /
  Ahrens-Dieter are refused with `eval.Unsupported` in slice 24 and added
  later from their papers, in demand order; ?RNGkind's description of each
  (for example the Wichmann-Hill multipliers and moduli) is the starting
  point.
- **Gates**: the goldens compare the whole `.Random.seed` after each call,
  not only the values; a differential fuzz mode runs random sequences of
  `set.seed`, `runif`, `sample`, `rnorm` and `RNGkind` against the oracle and
  compares the state after each.

## 4. Upstream gaps: pull requests, coordinated with the owners

Each gap goes to its library as a **pull request**, agreed beforehand with
that organisation's owners. Nothing is written into their trackers from this
repository. Until a gap lands, the R function it serves is
`eval.Unsupported` (never a local approximation that would need replacing).

**go-ndarray/ndarray** (`random`):

- `MT19937.SetState` stores `Pos` without checking it
  (`random/mt19937.go:159-162`); a position above 624 makes the next
  `Uint32` index past the state. R lets a user assign any `.Random.seed`, so
  either the setter validates, or the R layer does. The R layer validates in
  slice 24 in any case; the PR adds the check where every binding benefits.

**go-numerics/special** (its `docs/roadmap.md`, 2026-11-11, lists these):

- the log-space set (`Log1mExp`, `LogAddExp`, `LogSubExp`, `Log1pmx`),
  `LogGamma1p`, `SinPi`/`CosPi`/`TanPi` (items 1, 3, 4): the `log.p` and
  `lower.tail` paths of slices 26–27;
- `Stirlerr`, `Bd0` and the raw binomial/Poisson densities, today private in
  go-statistics (item 2): the saddle-point densities (Loader 2000, cited by
  ?Binomial, ?Poisson, ?GammaDist);
- log-probability inverses for `GammaPInv` and `BetaIncInv` (item A1): `qgamma`,
  `qbeta` with `log.p = TRUE`; and a bound for `BetaIncInv` at a or b < 1e-3
  (go-numerics/special#11), which `qbeta` reaches;
- noncentral beta, t and F, Owen's T (items 5, 10): Lenth (1987) AS 226 and
  Frick (1990), Lenth (1989) AS 243, as ?Beta, ?TDist and ?FDist cite them;
- Tukey studentized range and the exact Wilcoxon and signed-rank
  distributions (item 11; the roadmap leaves open whether they belong in
  go-statistics): Copenhaver & Holland (1988) and Odeh & Evans AS 70 (?Tukey);
- the accuracy outliers that slice 25 lists as divergences: logbeta near
  log B = 0, digamma near its negative roots, `jv` at half-integer order.

**go-statistics/stats**:

- d/p/q for the families of slice 27 that it lacks (noncentral t, beta and F;
  Tukey; Wilcoxon; signed rank), on top of special's items above;
- **samplers over a uniform source**, separate from the NumPy-compatible
  `Rvs`, implementing the algorithms R documents: Kachitvichyanukul &
  Schmeiser (1988) for the binomial, Ahrens & Dieter (1982, TOMS 8(2)) for
  Poisson, Ahrens & Dieter (1982, CACM 25(1)) and (1974) for gamma, Cheng
  (1978) for beta, Kachitvichyanukul & Schmeiser (1985) for the
  hypergeometric, Ahrens & Dieter (1972) for the exponential, Devroye (1986)
  for the negative binomial and geometric;
- quantiles at a log probability (`PPF` from `LogCDF`/`LogSF`);
- its dependencies are old (special v0.5.3, optimize v0.6.0 in its `go.mod`):
  a bump to special ≥ v0.20 and optimize ≥ v0.9 comes first, so the r module
  does not carry two versions' accuracy.

**go-numerics/optimize** (its `docs/roadmap.md` already plans the first two):

- PORT `nlminb` from netlib (roadmap item 7; netlib licence notice to be
  checked first, as the roadmap says);
- `extendInt`-style bracket extension (roadmap D2);
- R's `optim` variants, which are not SciPy's: Nelder-Mead with
  alpha 1.0, beta 0.5, gamma 2.0, and BFGS and CG "based originally on Pascal
  code in Nash (1990)" (?optim), CG `type` 1–3 (Fletcher-Reeves,
  Polak-Ribière, Beale-Sorenson); implemented from Nash's book, the code's
  licence checked before any is read;
- SANN (Bélisle 1992, ?optim);
- `nlm`'s method: Dennis & Schnabel (1983) and Schnabel, Koontz & Weiss
  (1985), cited by ?nlm;
- `polyroot`: no polynomial root finder exists in go-numerics or go-ndarray
  (searched in the module cache); Jenkins & Traub (1972) Algorithm 419
  (?polyroot). The repository is agreed with the owners.

**go-numerics/integrate**:

- a batched (vector) integrand for `Quad`: today `Quad` takes
  `func(float64) float64`, while R calls `f` once per rule application with all
  its nodes (?integrate: "`f` must accept a vector of inputs"); the number and
  order of `f`'s calls is observable in R. This is part of integrate's roadmap
  item 4 (batched calls).

## 5. Gates, per slice

The five gates of M3 (goldens, differential fuzz with exclusions counted,
`FuzzEval`, bomb tests under `GOMEMLIMIT`, performance) apply, plus:

1. every numeric golden records our error and R's error against the oracle
   (mpmath, `math/big`), and fails when ours exceeds the slice's bound;
2. bomb tests for the new resource paths: `sample.int(1e15, 1e9)`,
   `rnorm(2^31)`, `integrate` with `subdivisions = 1e9`, `optim` with
   `maxit = .Machine$integer.max` all refuse under `eval.Options` limits;
3. the coverage metric of M3 extends to `stats` (`rparity` gains a
   `statscov` mode on the oracle's `ls("package:stats")`).
