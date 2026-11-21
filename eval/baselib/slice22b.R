# Slice 22: base closures that need a little structure.
# Written from the help pages of R 4.6.2 and the oracle's answers
# (docs/m3-plan.md); BSD licence.

callCC <- function(f) {
  tag <- new.env()
  live <- TRUE
  k <- function(v) {
    if (!live) stop("no function to return from, jumping to top level")
    signalCondition(structure(class = c("..callCCjump", "condition"),
                              list(message = "", call = NULL, tag = tag, value = v)))
    stop("no function to return from, jumping to top level")
  }
  on.exit(live <- FALSE)
  tryCatch(f(k), ..callCCjump = function(j) {
    if (identical(j$tag, tag)) j$value else stop(j)
  })
}

`[<-.POSIXct` <- function(x, ..., value) {
  if (!length(value)) return(x)
  value <- unclass(as.POSIXct(value))
  cl <- oldClass(x)
  tz <- attr(x, "tzone")
  attr(x, "tzone") <- NULL
  x <- unclass(x)
  x[...] <- value
  class(x) <- cl
  attr(x, "tzone") <- tz
  x
}

`[<-.difftime` <- function(x, ..., value) {
  if (!length(value)) return(x)
  u <- attr(x, "units")
  if (inherits(value, "difftime")) {
    secs <- c(secs = 1, mins = 60, hours = 3600, days = 86400, weeks = 604800)
    value <- as.vector(unclass(value)) * secs[[attr(value, "units")]] / secs[[u]]
  } else value <- as.vector(unclass(value))
  cl <- oldClass(x)
  attr(x, "units") <- NULL
  x <- unclass(x)
  x[...] <- value
  class(x) <- cl
  attr(x, "units") <- u
  x
}

as.data.frame.ts <- function(x, row.names = NULL, optional = FALSE, ...) {
  if (is.matrix(x)) as.data.frame.matrix(x, row.names, optional, ...)
  else as.data.frame.vector(x, row.names, optional, ...)
}

as.data.frame.model.matrix <- function(x, row.names = NULL, optional = FALSE,
                                       make.names = TRUE, ...) {
  d <- dim(x)
  nrows <- d[[1L]]
  dn <- dimnames(x)
  if (is.null(row.names)) row.names <- dn[[1L]]
  value <- list(x)
  if (!optional) names(value) <- deparse(substitute(x))[[1L]]
  if (is.null(row.names)) row.names <- .set_row_names(nrows)
  attr(value, "row.names") <- row.names
  class(value) <- "data.frame"
  value
}

sort.int <- function(x, partial = NULL, na.last = NA, decreasing = FALSE,
                     method = c("auto", "shell", "quick", "radix"),
                     index.return = FALSE) {
  if (!is.null(partial)) stop("'partial' is not supported")
  method <- match.arg(method)
  if (is.na(na.last)) x <- x[!is.na(x)]
  o <- order(x, na.last = na.last, decreasing = decreasing)
  y <- x[o]
  if (index.return) list(x = y, ix = o) else y
}

sort.list <- function(x, partial = NULL, na.last = TRUE, decreasing = FALSE,
                      method = c("auto", "shell", "quick", "radix")) {
  if (!is.null(partial)) stop("'partial' is not supported")
  method <- match.arg(method)
  if (!is.atomic(x))
    stop("'x' must be atomic for 'sort.list', method \"shell\" and \"quick\"\nHave you called 'sort' on a list?")
  if (is.na(na.last)) x <- x[!is.na(x)]
  order(x, na.last = na.last, decreasing = decreasing)
}

rapply <- function(object, f, classes = "ANY", deflt = NULL,
                   how = c("unlist", "replace", "list"), ...) {
  how <- match.arg(how)
  if (typeof(object) != "list" && !is.expression(object))
    stop("'object' must be a list or expression")
  f <- match.fun(f)
  want <- function(x) identical(classes, "ANY") || any(class(x) %in% classes) ||
    any(inherits(x, classes))
  rec <- function(x) {
    if (typeof(x) == "list" || is.expression(x)) {
      y <- x
      for (i in seq_along(x)) y[i] <- list(rec(x[[i]]))
      if (how != "replace") {
        keep <- names(x)
        attributes(y) <- NULL
        names(y) <- keep
      }
      y
    } else if (want(x)) f(x, ...)
    else if (how == "replace") x
    else deflt
  }
  res <- rec(object)
  if (how == "unlist") unlist(res, recursive = TRUE) else res
}

asplit <- function(x, MARGIN) {
  dl <- length(dim(x))
  if (dl == 0L) stop("dim(x) must have a positive length")
  if (is.object(x)) x <- if (dl == 2L) as.matrix(x) else as.array(x)
  d <- dim(x)
  dn <- dimnames(x)
  if (is.character(MARGIN)) {
    labels <- names(dn)
    if (is.null(labels)) stop("'x' must have named dimnames")
    pos <- match(MARGIN, labels)
    if (anyNA(pos)) stop("not all elements of 'MARGIN' are names of dimensions")
    MARGIN <- pos
  }
  rest <- setdiff(seq_len(dl), MARGIN)
  # one cell of the result per combination of the margin indices; a cell is
  # the sub-array that fixes those indices, kept with the other dimensions
  grid <- as.matrix(do.call(expand.grid, lapply(d[MARGIN], seq_len)))
  cells <- vector("list", nrow(grid))
  for (k in seq_along(cells)) {
    idx <- rep(list(quote(expr = )), dl)
    for (j in seq_along(MARGIN)) idx[[MARGIN[j]]] <- grid[k, j]
    piece <- do.call(`[`, c(list(x), idx, list(drop = FALSE)))
    dim(piece) <- d[rest]
    if (!is.null(dn) && !all(vapply(dn[rest], is.null, NA))) dimnames(piece) <- dn[rest]
    cells[[k]] <- piece
  }
  dim(cells) <- d[MARGIN]
  if (!is.null(dn) && !all(vapply(dn[MARGIN], is.null, NA))) dimnames(cells) <- dn[MARGIN]
  cells
}

drop <- function(x) {
  d <- dim(x)
  if (is.null(d) || !any(d == 1L)) return(x)
  keep <- d != 1L
  dn <- dimnames(x)
  y <- x
  if (sum(keep) <= 1L) {
    dim(y) <- NULL
    if (!is.null(dn) && sum(keep) == 1L && !is.null(nm <- dn[[which(keep)]]))
      names(y) <- nm
  } else {
    dim(y) <- d[keep]
    if (!is.null(dn)) {
      dn2 <- dn[keep]
      if (!all(vapply(dn2, is.null, NA))) dimnames(y) <- dn2
    }
  }
  y
}

args <- function(name) {
  if (is.character(name)) name <- get(name, mode = "function", envir = parent.frame())
  if (typeof(name) == "closure") {
    f <- name
    body(f) <- NULL
    environment(f) <- globalenv()
    return(f)
  }
  NULL
}

max.col <- function(m, ties.method = c("random", "first", "last")) {
  ties.method <- match.arg(ties.method)
  m <- as.matrix(m)
  res <- integer(nrow(m))
  for (i in seq_len(nrow(m))) {
    r <- m[i, ]
    if (anyNA(r)) { res[i] <- NA_integer_; next }
    w <- which(r == max(r))
    res[i] <- switch(ties.method,
      first = w[1L],
      last = w[length(w)],
      random = if (length(w) == 1L) w else
        stop("ties.method = \"random\" needs the random number generator, which this evaluator does not draw from"))
  }
  res
}

t.default <- function(x) t(unclass(x))
