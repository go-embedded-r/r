# Base library, slice 16 (5): small closures that need no engine. Written from
# the help pages of R 4.6.2 (as.table, simplify2array, data.matrix, list2DF,
# sort_by, provideDimnames, mode, dynGet, duplicated, subset, transform,
# system.time, single) and from experiments with the oracle.

# ---- tables ----------------------------------------------------------------

as.table <- function(x, ...) UseMethod("as.table")

as.table.default <- function(x, ...) {
  if (is.table(x)) return(x)
  if (!(is.array(x) || is.numeric(x) || is.logical(x))) stop("cannot coerce to a table")
  if (is.null(dim(x))) {
    labels <- names(x)
    dim(x) <- length(x)
    if (is.null(labels)) labels <- make.unique(rep_len(LETTERS, length(x)))
    dimnames(x) <- list(labels)
  } else {
    given <- dimnames(x)
    if (is.null(given)) given <- vector("list", length(dim(x)))
    for (i in seq_along(given))
      if (is.null(given[[i]])) given[[i]] <- make.unique(rep_len(LETTERS, dim(x)[i]))
    dimnames(x) <- given
  }
  class(x) <- c("table", oldClass(x))
  x
}

provideDimnames <- function(x, sep = "", base = list(LETTERS), unique = TRUE) {
  d <- dim(x)
  given <- dimnames(x)
  if (is.null(given)) given <- vector("list", length(d))
  for (i in seq_along(d)) {
    if (!is.null(given[[i]])) next
    labels <- rep_len(base[[(i - 1L) %% length(base) + 1L]], d[i])
    given[[i]] <- if (unique) make.unique(labels, sep = sep) else labels
  }
  dimnames(x) <- given
  x
}

# ---- converting ------------------------------------------------------------

simplify2array <- function(x, higher = TRUE, except = c(0L, 1L))
  sapply(x, function(e) e, simplify = if (higher) "array" else TRUE)

data.matrix <- function(frame, rownames.force = NA) {
  if (!is.data.frame(frame)) return(as.matrix(frame))
  cols <- lapply(frame, function(v) {
    if (is.character(v)) v <- factor(v)
    if (is.factor(v) || is.logical(v)) as.integer(v) else v
  })
  out <- do.call(cbind, unname(cols))
  if (is.null(out)) out <- matrix(0L, nrow(frame), 0L)
  automatic <- is.integer(attr(frame, "row.names")) &&
    identical(attr(frame, "row.names"), seq_len(nrow(frame)))
  rn <- if (isTRUE(rownames.force) || (is.na(rownames.force) && !automatic)) row.names(frame)
  dimnames(out) <- list(rn, names(frame))
  out
}

list2DF <- function(x = list(), nrow = 0L) {
  stopifnot(is.list(x), is.null(nrow) || nrow >= 0L)
  n <- length(x)
  if (n) {
    sizes <- unique(lengths(x))
    if (length(sizes) > 1L || (!missing(nrow) && !is.null(nrow) && sizes != nrow))
      stop("all variables should have the same length")
    nrow <- sizes
  }
  if (is.null(names(x))) names(x) <- character(n)
  class(x) <- "data.frame"
  attr(x, "row.names") <- .set_row_names(nrow)
  x
}

`mode<-` <- function(x, value) {
  kept <- attributes(x)
  converter <- paste0("as.", value)
  x <- do.call(converter, list(x))
  attributes(x) <- kept
  x
}

logb <- function(x, base = exp(1)) if (missing(base)) log(x) else log(x, base)

as.single <- function(x, ...) UseMethod("as.single")
as.single.default <- function(x, ...) structure(as.double(x), Csingle = TRUE)
single <- function(length = 0L) structure(vector("double", length), Csingle = TRUE)

# ---- pretty ----------------------------------------------------------------

# the multiples of a unit (1, 2, 5 or 10 times a power of ten) around a range;
# ..prettyCore is in Go (eval/pretty16.go)
.pretty <- function(x, n = 5L, min.n = n %/% 3L, shrink.sml = 0.75, high.u.bias = 1.5,
                    u5.bias = 0.5 + 1.5 * high.u.bias, eps.correct = 0L, bounds = TRUE) {
  z <- ..prettyCore(x[1L], x[2L], n, min.n, shrink.sml, high.u.bias, u5.bias, eps.correct)
  if (bounds) list(l = min(z[5L], z[1L] * z[4L]), u = max(z[6L], z[2L] * z[4L]), n = as.integer(z[3L]))
  else list(ns = z[1L], nu = z[2L], n = as.integer(z[3L]), unit = z[4L])
}

pretty <- function(x, ...) UseMethod("pretty")

pretty.default <- function(x, n = 5L, min.n = n %/% 3L, shrink.sml = 0.75, high.u.bias = 1.5,
                           u5.bias = 0.5 + 1.5 * high.u.bias, eps.correct = 0L, f.min = 2^-20, ...) {
  x <- as.numeric(x)
  x <- x[is.finite(x)]
  if (!length(x)) return(x)
  if (!is.numeric(n) || is.na(n) || n < 0) stop("invalid 'n' argument")
  if (!is.numeric(min.n) || is.na(min.n) || min.n < 0 || min.n > n) stop("invalid 'min.n' argument")
  z <- .pretty(range(x), n, min.n, shrink.sml, high.u.bias, u5.bias, eps.correct)
  s <- ..wholeInts(seq.int(z$l, z$u, length.out = z$n + 1L))
  if (!eps.correct && z$n) {
    delta <- (z$u - z$l) / z$n
    small <- abs(s) < 1e-14 * delta
    if (any(small)) s[small] <- 0
  }
  s
}

# ---- sorting by another vector ---------------------------------------------

sort_by <- function(x, y, ...) UseMethod("sort_by")

sort_by.default <- function(x, y, ...) x[order(y, ...)]

sort_by.data.frame <- function(x, y, ...) {
  if (inherits(y, "formula")) stop("a formula is not supported here")
  keys <- if (is.list(y)) unname(as.list(y)) else list(y)
  x[do.call(order, c(keys, list(...))), , drop = FALSE]
}

# ---- looking up the calling frames -----------------------------------------

dynGet <- function(x, ifnotfound = stop(gettextf("%s not found", sQuote(x)), domain = NA),
                   minframe = 1L, inherits = FALSE) {
  n <- sys.nframe() - 1L
  while (n > minframe) {
    n <- n - 1L
    env <- sys.frame(n)
    if (exists(x, envir = env, inherits = inherits)) return(get(x, envir = env, inherits = inherits))
  }
  ifnotfound
}

objects <- ls

# ---- arrays of rows: unique, duplicated, anyDuplicated ---------------------

# a text key for each slice of x along MARGIN
..sliceKeys <- function(x, MARGIN) apply(x, MARGIN, function(v) paste(v, collapse = "\r"))

duplicated.array <- function(x, incomparables = FALSE, MARGIN = 1L, fromLast = FALSE, ...) {
  if (!is.logical(incomparables) || incomparables) .NotYetUsed("incomparables != FALSE")
  res <- duplicated(as.vector(..sliceKeys(x, MARGIN)), fromLast = fromLast)
  dim(res) <- dim(x)[MARGIN]
  labels <- dimnames(x)[MARGIN]
  if (!is.null(labels) && !all(vapply(labels, is.null, NA))) dimnames(res) <- labels
  res
}

duplicated.matrix <- duplicated.array

anyDuplicated.array <- function(x, incomparables = FALSE, MARGIN = 1L, fromLast = FALSE, ...)
  anyDuplicated(as.vector(..sliceKeys(x, MARGIN)), fromLast = fromLast)

anyDuplicated.matrix <- anyDuplicated.array

unique.array <- function(x, incomparables = FALSE, MARGIN = 1, fromLast = FALSE, ...) {
  if (!is.logical(incomparables) || incomparables) .NotYetUsed("incomparables != FALSE")
  keep <- !as.vector(duplicated.array(x, MARGIN = MARGIN, fromLast = fromLast))
  index <- rep(list(TRUE), length(dim(x)))
  index[[MARGIN[1L]]] <- keep
  do.call(`[`, c(list(x), index, list(drop = FALSE)))
}

unique.matrix <- unique.array

# ---- subset and transform on other classes ---------------------------------

subset.matrix <- function(x, subset, select, drop = FALSE, ...) {
  rows <- if (missing(subset)) rep_len(TRUE, nrow(x)) else {
    r <- subset
    r & !is.na(r)
  }
  cols <- if (missing(select)) TRUE else {
    positions <- as.list(seq_len(ncol(x)))
    names(positions) <- colnames(x)
    eval(substitute(select), positions, parent.frame())
  }
  x[rows, cols, drop = drop]
}

transform.default <- function(`_data`, ...)
  transform.data.frame(data.frame(`_data`, check.names = FALSE), ...)

# ---- time spent ------------------------------------------------------------

system.time <- function(expr, gcFirst = TRUE) {
  begin <- proc.time()
  expr
  end <- proc.time()
  structure(unclass(end) - unclass(begin), class = "proc_time")
}

unix.time <- function(...) .Defunct("system.time")

summary.proc_time <- function(object, ...) {
  v <- unclass(object)
  c(user = v[[1L]] + v[[4L]], system = v[[2L]] + v[[5L]], elapsed = v[[3L]])
}

print.proc_time <- function(x, ...) {
  print(summary(x, ...), ...)
  invisible(x)
}
