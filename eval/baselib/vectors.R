# Base library, part 1: logical helpers and vector glue.
# Written from the help pages of R 4.6.2 and checked against the oracle
# (docs/m3-plan.md). BSD licence, like the rest of this repository.

# exactly one logical value, whatever attributes it carries
isTRUE <- function(x) is.logical(x) && length(x) == 1L && identical(as.vector(x), TRUE)

isFALSE <- function(x) is.logical(x) && length(x) == 1L && identical(as.vector(x), FALSE)

xor <- function(x, y) (x | y) & !(x & y)

Negate <- function(f) {
  f <- match.fun(f)
  function(...) !f(...)
}

setNames <- function(object = nm, nm) {
  names(object) <- nm
  object
}

is.element <- function(el, table) match(el, table, 0L) > 0L

union <- function(x, y) unique(c(as.vector(x), as.vector(y)))

intersect <- function(x, y) {
  if (is.null(x) || is.null(y)) return(NULL)
  u <- as.vector(x)
  v <- as.vector(y)
  c(u[!duplicated(unclass(u)) & (match(u, v, 0L) > 0L)], v[numeric()])
}

setdiff <- function(x, y) {
  x <- as.vector(x)
  y <- as.vector(y)
  unique(if (length(x) || length(y)) x[match(x, y, 0L) == 0L] else x)
}

append <- function(x, values, after = length(x)) {
  split <- min(max(after, 0L), length(x))
  head <- seq_len(split)
  tail <- seq_len(length(x) - split) + split
  c(x[head], values, x[tail])
}

head <- function(x, n = 6L, ...) {
  if (length(dim(x)) == 2L) {
    nr <- nrow(x)
    n <- if (n < 0L) max(nr + n, 0L) else min(n, nr)
    return(x[seq_len(n), , drop = FALSE])
  }
  len <- length(x)
  n <- if (n < 0L) max(len + n, 0L) else min(n, len)
  x[seq_len(n)]
}

tail <- function(x, n = 6L, ...) {
  if (length(dim(x)) == 2L) {
    nr <- nrow(x)
    n <- if (n < 0L) max(nr + n, 0L) else min(n, nr)
    idx <- seq.int(to = nr, length.out = n)
    r <- x[idx, , drop = FALSE]
    dn <- dimnames(x)
    if (is.null(dn) || is.null(dn[[1L]]))
      dimnames(r) <- list(paste0("[", idx, ",]"), if (is.null(dn)) NULL else dn[[2L]])
    return(r)
  }
  len <- length(x)
  n <- if (n < 0L) max(len + n, 0L) else min(n, len)
  x[seq.int(to = len, length.out = n)]
}

# yes is evaluated only when some element of test is true, no only when some
# is false (their promises are forced inside fill)
ifelse <- function(test, yes, no) {
  flag <- if (!is.atomic(test)) as.logical(test)
          else if (is.logical(test)) test
          else `storage.mode<-`(test, "logical")
  out <- flag
  size <- length(flag)
  fill <- function(where, from) {
    if (length(where) > 0L) out[where] <<- rep(from, length.out = size)[where]
  }
  fill(which(flag), yes)
  fill(which(!flag), no)
  out
}

nrow <- function(x) dim(x)[1L]

ncol <- function(x) dim(x)[2L]

NROW <- function(x) {
  extent <- dim(x)
  if (is.null(extent)) length(x) else extent[1L]
}

NCOL <- function(x) {
  extent <- dim(x)
  if (length(extent) < 2L) 1L else extent[2L]
}

lengths <- function(x, use.names = TRUE) {
  ans <- vapply(x, length, 1L, USE.NAMES = FALSE)
  if (use.names && !is.null(names(x))) names(ans) <- names(x)
  ans
}

rep_len <- function(x, length.out) rep(x, length.out = length.out)

diff <- function(x, lag = 1L, differences = 1L) {
  for (i in seq_len(differences)) {
    n <- length(x)
    if (lag >= n) return(x[0L])
    x <- x[(1L + lag):n] - x[1L:(n - lag)]
  }
  x
}

outer <- function(X, Y, FUN = "*", ...) {
  product <- is.character(FUN) && length(FUN) == 1L && FUN == "*"
  FUN <- match.fun(FUN)
  dX <- if (is.null(dim(X))) length(X) else dim(X)
  dY <- if (is.null(dim(Y))) length(Y) else dim(Y)
  nmX <- if (is.null(dim(X))) names(X) else NULL
  nmY <- if (is.null(dim(Y))) names(Y) else NULL
  if (product && ...length() == 0L) {
    robj <- as.vector(X) %*% t(as.vector(Y))
  } else {
    Xv <- rep(as.vector(X), times = length(Y))
    Yv <- rep(as.vector(Y), each = length(X))
    robj <- FUN(Xv, Yv, ...)
  }
  if (length(robj) != length(X) * length(Y)) stop("dims [product ", length(Xv), "] do not match the length of object [", length(robj), "]")
  dim(robj) <- c(dX, dY)
  if (!(is.null(nmX) && is.null(nmY))) dimnames(robj) <- list(nmX, nmY)
  robj
}
