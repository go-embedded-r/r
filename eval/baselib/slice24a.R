# Base library, slice 24 (random streams). Written from the help pages of
# R 4.6.2 (?Random, ?RNGkind, ?set.seed, ?sample, ?jitter) and from
# experiments with the oracle. The generators are Go (..runif, ..rnorm,
# ..sample, ..sample2, ..set.seed, ..RNGkind).

runif <- function(n, min = 0, max = 1) ..runif(n, min, max)

rnorm <- function(n, mean = 0, sd = 1) ..rnorm(n, mean, sd)

set.seed <- function(seed, kind = NULL, normal.kind = NULL, sample.kind = NULL)
  invisible(..set.seed(seed, kind, normal.kind, sample.kind))

RNGkind <- function(kind = NULL, normal.kind = NULL, sample.kind = NULL) {
  r <- ..RNGkind(kind, normal.kind, sample.kind)
  if (is.null(kind) && is.null(normal.kind) && is.null(sample.kind)) r else invisible(r)
}

RNGversion <- function(vstr) {
  vnum <- as.numeric(strsplit(as.character(vstr), ".", fixed = TRUE)[[1L]])
  if (length(vnum) < 2L) stop("malformed version string")
  v <- vnum[1L] * 1e6 + vnum[2L] * 1e3 + (if (length(vnum) >= 3L) vnum[3L] else 0)
  if (v >= 3006000) RNGkind("Mersenne-Twister", "Inversion", "Rejection")
  else if (v >= 1007000) RNGkind("Mersenne-Twister", "Inversion", "Rounding")
  else RNGkind("Marsaglia-Multicarry", "Buggy Kinderman-Ramage", "Rounding")
}

sample <- function(x, size, replace = FALSE, prob = NULL) {
  if (length(x) != 1L || !is.numeric(x) || !is.finite(x) || x < 1) {
    # a set of things to shuffle: draw positions, then look them up
    if (missing(size)) size <- length(x)
    return(x[sample.int(length(x), size, replace, prob)])
  }
  # a single number n stands for 1:n
  if (missing(size)) size <- x
  sample.int(x, size, replace, prob)
}

sample.int <- function(n, size = n, replace = FALSE, prob = NULL,
                       useHash = n > 1e7 && !replace && is.null(prob) && !is.null(size) && size <= n/2) {
  stopifnot(length(n) == 1L)
  if (!useHash) return(..sample(n, size, replace, prob))
  stopifnot(is.null(prob), !replace)
  ..sample2(n, size)
}

jitter <- function(x, factor = 1, amount = NULL) {
  if (!length(x)) return(x)
  if (!is.numeric(x)) stop("'x' must be numeric")
  lims <- range(x, na.rm = TRUE)
  spread <- lims[2L] - lims[1L]
  if (spread == 0) spread <- abs(lims[1L])
  if (spread == 0) spread <- 1
  half <- if (is.null(amount)) {
    levels <- unique(sort.int(round(x, 3 - floor(log10(spread)))))
    gaps <- diff(levels)
    step <- if (length(gaps)) min(gaps) else if (levels != 0) levels / 10 else spread / 10
    factor / 5 * abs(step)
  } else if (amount == 0) {
    factor * (spread / 50)
  } else amount
  x + runif(length(x), -half, half)
}
