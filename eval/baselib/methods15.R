# Base library, slice 15 (3): S3 methods of conditions, factors, dates and data
# frames. Written from the help pages of R 4.6.2 and experiments with the oracle.

# conditionMessage.condition, conditionCall.condition, print.condition,
# as.character.condition, as.character.error and within.data.frame are Go
# natives (eval/funcs15.go)

pmax.int <- function(..., na.rm = FALSE) pmax(..., na.rm = na.rm)

pmin.int <- function(..., na.rm = FALSE) pmin(..., na.rm = na.rm)

# matching keys: the dates and times match by their numbers
mtfrm <- function(x) UseMethod("mtfrm")

mtfrm.default <- function(x) {
  key <- as.character(x)
  if (length(key) != length(x)) stop("cannot mtfrm")
  key
}

mtfrm.Date <- function(x) as.numeric(x)

mtfrm.POSIXct <- function(x) as.numeric(x)

xtfrm.Date <- function(x) as.numeric(x)

xtfrm.POSIXct <- function(x) as.numeric(x)

xtfrm.difftime <- function(x) as.numeric(x)

xtfrm.factor <- function(x) as.integer(x)

xtfrm.AsIs <- function(x) {
  cl <- oldClass(x)
  oldClass(x) <- cl[cl != "AsIs"]
  NextMethod("xtfrm")
}

xtfrm.data.frame <- function(x) stop("cannot xtfrm data frames")

range.Date <- function(..., na.rm = FALSE, finite = FALSE)
  .Date(range(unclass(c(...)), na.rm = na.rm, finite = finite))

range.POSIXct <- function(..., na.rm = FALSE, finite = FALSE) {
  x <- c(...)
  .POSIXct(range(unclass(x), na.rm = na.rm, finite = finite), attr(x, "tzone"))
}

split.Date <- function(x, f, drop = FALSE, ...)
  lapply(split.default(seq_along(x), f, drop = drop, ...), function(ind) x[ind])

split.POSIXct <- function(x, f, drop = FALSE, ...)
  lapply(split.default(seq_along(x), f, drop = drop, ...), function(ind) x[ind])

`length<-.Date` <- function(x, value) .Date(NextMethod(), oldClass(x))

`length<-.POSIXct` <- function(x, value) .POSIXct(NextMethod(), attr(x, "tzone"), oldClass(x))

`length<-.difftime` <- function(x, value) .difftime(NextMethod(), attr(x, "units"), oldClass(x))

`length<-.factor` <- function(x, value) {
  lev <- levels(x)
  cl <- oldClass(x)
  y <- NextMethod()
  attr(y, "levels") <- lev
  oldClass(y) <- cl
  y
}

Math.Date <- function(x, ...) stop(gettextf("%s not defined for \"Date\" objects", .Generic), domain = NA)

Math.POSIXt <- function(x, ...) stop(gettextf("'%s' not defined for \"POSIXt\" objects", .Generic), domain = NA)

Math.factor <- function(x, ...) stop(gettextf("%s not meaningful for factors", sQuote(.Generic)))

Summary.factor <- function(..., na.rm) stop(gettextf("%s not meaningful for factors", sQuote(.Generic)))

Summary.ordered <- function(..., na.rm) {
  if (!(.Generic %in% c("max", "min", "range")))
    stop(gettextf("'%s' not defined for ordered factors", .Generic))
  parts <- list(...)
  lev <- levels(parts[[1L]])
  same <- vapply(parts, function(p) is.ordered(p) && identical(levels(p), lev), NA)
  if (!all(same)) stop("The level sets of factors are different")
  pick <- do.call(.Generic, c(lapply(parts, as.integer), na.rm = na.rm))
  factor(lev[pick], levels = lev, ordered = TRUE)
}

Ops.factor <- function(e1, e2) {
  if (!(.Generic %in% c("==", "!="))) {
    warning(gettextf("%s not meaningful for factors", sQuote(.Generic)))
    return(rep.int(NA, max(length(e1), if (!missing(e2)) length(e2))))
  }
  if (nzchar(.Method[1L])) e1 <- as.character(e1)
  if (nzchar(.Method[2L])) e2 <- as.character(e2)
  get(.Generic, mode = "function")(e1, e2)
}

Ops.ordered <- function(e1, e2) {
  if (.Generic %in% c("==", "!=")) return(NextMethod(.Generic))
  if (!(.Generic %in% c("<", ">", "<=", ">="))) {
    warning(gettextf("'%s' is not meaningful for ordered factors", .Generic))
    return(rep.int(NA, max(length(e1), if (!missing(e2)) length(e2))))
  }
  if (is.ordered(e1) && is.ordered(e2)) {
    if (!identical(levels(e1), levels(e2))) stop("level sets of factors are different")
  } else {
    lev <- levels(if (is.ordered(e1)) e1 else e2)
    if (!is.ordered(e1)) e1 <- match(e1, lev)
    if (!is.ordered(e2)) e2 <- match(e2, lev)
  }
  get(.Generic, mode = "function")(as.integer(e1), as.integer(e2))
}

Math.data.frame <- function(x, ...) {
  ok <- vapply(x, function(col) is.numeric(col) || is.logical(col) || is.complex(col), NA)
  if (!all(ok)) {
    nm <- names(x)
    if (is.null(nm)) nm <- seq_along(x)
    stop("non-numeric-alike variable(s) in data frame: ", paste(nm[!ok], collapse = ", "))
  }
  x[] <- lapply(x, .Generic, ...)
  x
}

Summary.data.frame <- function(..., na.rm) {
  parts <- lapply(list(...), function(d) {
    m <- as.matrix(d)
    if (!is.numeric(m) && !is.logical(m) && !is.complex(m))
      stop("only defined on a data frame with all numeric-alike variables")
    m
  })
  do.call(.Generic, c(parts, na.rm = na.rm))
}

Ops.data.frame <- function(e1, e2) {
  unary <- nargs() == 1L
  left <- nzchar(.Method[1L])
  right <- !unary && nzchar(.Method[2L])
  df <- if (left) e1 else e2
  nr <- nrow(df)
  nc <- length(df)
  logic <- .Generic %in% c("==", "!=", "<", ">", "<=", ">=", "&", "|", "!")
  if (unary) {
    cols <- lapply(e1, .Generic)
  } else {
    if (left && right) {
      if (!identical(dim(e1), dim(e2)))
        stop(gettextf("%s only defined for equally-sized data frames", sQuote(.Generic)))
      pa <- as.list(e1)
      pb <- as.list(e2)
    } else {
      vec <- rep_len(as.vector(if (left) e2 else e1), nr * nc)
      pieces <- split(vec, rep.int(seq_len(nc), rep.int(nr, nc)))
      names(pieces) <- NULL
      if (left) { pa <- as.list(e1); pb <- pieces }
      else { pa <- pieces; pb <- as.list(e2) }
    }
    f <- get(.Generic, mode = "function")
    cols <- vector("list", nc)
    for (i in seq_len(nc)) cols[[i]] <- f(pa[[i]], pb[[i]])
  }
  if (logic) {
    out <- matrix(unlist(cols, use.names = FALSE), nr, nc)
    dimnames(out) <- list(if (.row_names_info(df) > 0L) row.names(df), names(df))
    out
  } else {
    names(cols) <- names(df)
    attr(cols, "class") <- oldClass(df)
    attr(cols, "row.names") <- .row_names_info(df, 0L)
    cols
  }
}
