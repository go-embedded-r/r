# Base library, part 1: logical helpers and vector glue.
# Written from the help pages of R 4.6.2 and checked against the oracle
# (docs/m3-plan.md). BSD licence, like the rest of this repository.

isTRUE <- function(x) is.logical(x) && length(x) == 1L && !is.na(x) && x

isFALSE <- function(x) is.logical(x) && length(x) == 1L && !is.na(x) && !x

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
  lengx <- length(x)
  if (!after) c(values, x)
  else if (after >= lengx) c(x, values)
  else c(x[1L:after], values, x[(after + 1L):lengx])
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

ifelse <- function(test, yes, no) {
  if (is.atomic(test)) {
    if (typeof(test) != "logical") storage.mode(test) <- "logical"
  } else test <- as.logical(test)
  ans <- test
  len <- length(ans)
  ypos <- which(test)
  npos <- which(!test)
  if (length(ypos) > 0L) ans[ypos] <- rep(yes, length.out = len)[ypos]
  if (length(npos) > 0L) ans[npos] <- rep(no, length.out = len)[npos]
  ans
}

nrow <- function(x) dim(x)[1L]

ncol <- function(x) dim(x)[2L]

NROW <- function(x) if (length(d <- dim(x))) d[1L] else length(x)

NCOL <- function(x) if (length(d <- dim(x)) > 1L) d[2L] else 1L

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
