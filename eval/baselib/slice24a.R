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
  if (length(x) == 1L && is.numeric(x) && is.finite(x) && x >= 1) {
    if (missing(size)) size <- x
    sample.int(x, size, replace, prob)
  } else {
    if (missing(size)) size <- length(x)
    x[sample.int(length(x), size, replace, prob)]
  }
}

sample.int <- function(n, size = n, replace = FALSE, prob = NULL,
                       useHash = (n > 1e7 && !replace && is.null(prob) && (!is.null(size)) && size <= n/2)) {
  stopifnot(length(n) == 1L)
  if (useHash) {
    stopifnot(is.null(prob), !replace)
    ..sample2(n, size)
  } else ..sample(n, size, replace, prob)
}

jitter <- function(x, factor = 1, amount = NULL) {
  if (length(x) == 0L) return(x)
  if (!is.numeric(x)) stop("'x' must be numeric")
  z <- diff(r <- range(x, na.rm = TRUE))
  if (z == 0) z <- abs(r[1L])
  if (z == 0) z <- 1
  if (is.null(amount)) {
    d <- diff(xx <- unique(sort.int(round(x, 3 - floor(log10(z))))))
    d <- if (length(d)) min(d) else if (xx != 0) xx / 10 else z / 10
    amount <- factor / 5 * abs(d)
  } else if (amount == 0) amount <- factor * (z / 50)
  x + runif(length(x), -amount, amount)
}
