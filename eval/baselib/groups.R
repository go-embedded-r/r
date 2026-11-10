# Base library, part 8: grouping and small helpers. Written from the help
# pages of split, unsplit, tapply, rle, sequence, replace, replicate, setequal,
# getElement, eval.parent, bquote, zapsmall, factorial and the match operators
# in R 4.6.2 and checked against the oracle (docs/m3-plan.md). BSD licence.

force <- function(x) x

`%||%` <- function(x, y) if (is.null(x)) y else x

`%notin%` <- function(x, table) is.na(match(x, table))

grepv <- function(pattern, x, ignore.case = FALSE, perl = FALSE, value = TRUE,
                  fixed = FALSE, useBytes = FALSE, invert = FALSE)
  grep(pattern, x, ignore.case = ignore.case, perl = perl, value = value,
       fixed = fixed, useBytes = useBytes, invert = invert)

eval.parent <- function(expr, n = 1) eval(expr, parent.frame(n + 1))

replace <- function(x, list, values) {
  x[list] <- values
  x
}

replicate <- function(n, expr, simplify = "array") {
  body <- eval.parent(substitute(function(...) expr))
  sapply(integer(n), body, simplify = simplify)
}

setequal <- function(x, y) {
  x <- as.vector(x)
  y <- as.vector(y)
  !anyNA(match(x, y)) && !anyNA(match(y, x))
}

getElement <- function(object, name) object[[name, exact = TRUE]]

factorial <- function(x) gamma(x + 1)

lfactorial <- function(x) lgamma(x + 1)

zapsmall <- function(x, digits = getOption("digits"),
                     mFUN = function(x, ina) max(abs(x[!ina])), min.d = 0L) {
  if (length(digits) == 0L) stop("invalid 'digits'")
  ina <- is.na(x)
  if (all(ina)) return(x)
  top <- mFUN(x, ina)
  if (top > 0) digits <- max(min.d, digits - as.numeric(ceiling(log10(top))))
  round(x, digits = digits)
}

sequence <- function(nvec, ...) UseMethod("sequence")

sequence.default <- function(nvec, from = 1L, by = 1L, ...) {
  n <- length(nvec)
  longest <- max(length(from), length(by))
  if (n > 0L && n < longest)
    warning(gettextf("length(nvec) %d < %d = max(length(from), length(by)) -- future R`s default 'recycle = TRUE' will recycle 'nvec'", n, longest), call. = FALSE)
  if (n == 0L) return(integer())
  from <- rep_len(from, n)
  by <- rep_len(by, n)
  out <- vector("list", n)
  for (i in seq_len(n)) out[[i]] <- from[i] + by[i] * (seq_len(nvec[i]) - 1L)
  as.integer(unlist(out))
}





split <- function(x, f, drop = FALSE, ...) UseMethod("split")

..dropUnused <- function(f) {
  used <- tabulate(as.integer(f), nlevels(f)) > 0L
  factor(as.character(f), levels = levels(f)[used])
}

..crossFactors <- function(fs, sep, lex.order) {
  fs <- lapply(fs, as.factor)
  lev <- levels(fs[[1L]])
  code <- as.integer(fs[[1L]]) - 1L
  for (k in seq_along(fs)[-1L]) {
    lk <- levels(fs[[k]])
    ck <- as.integer(fs[[k]]) - 1L
    if (lex.order) {
      code <- code * length(lk) + ck
      lev <- as.vector(t(outer(lev, lk, paste, sep = sep)))
    } else {
      code <- code + length(lev) * ck
      lev <- as.vector(outer(lev, lk, paste, sep = sep))
    }
  }
  factor(code + 1L, levels = seq_along(lev), labels = lev)
}

split.default <- function(x, f, drop = FALSE, sep = ".", lex.order = FALSE, ...) {
  if (is.list(f)) f <- ..crossFactors(f, sep, lex.order)
  if (!is.factor(f)) f <- as.factor(f)
  if (drop) f <- ..dropUnused(f)
  if (length(f) > 0L && length(x) %% length(f) != 0L)
    warning("data length is not a multiple of split variable")
  codes <- rep_len(as.integer(f), length(x))
  lev <- levels(f)
  out <- lapply(seq_along(lev), function(i) x[which(codes == i)])
  names(out) <- lev
  out
}

`split<-` <- function(x, f, drop = FALSE, ..., value) UseMethod("split<-")

`split<-.default` <- function(x, f, drop = FALSE, ..., value) {
  groups <- split(seq_along(x), f, drop = drop, ...)
  k <- 0L
  for (idx in groups) {
    k <- k %% length(value) + 1L
    x[idx] <- value[[k]]
  }
  x
}

unsplit <- function(value, f, drop = FALSE) {
  n <- length(if (is.list(f)) f[[1L]] else f)
  x <- value[[1L]][rep(NA_integer_, n)]
  split(x, f, drop = drop) <- value
  x
}

tapply <- function(X, INDEX, FUN = NULL, ..., default = NA, simplify = TRUE) {
  if (!is.null(FUN)) FUN <- match.fun(FUN)
  if (!is.list(INDEX)) INDEX <- list(INDEX)
  INDEX <- lapply(INDEX, as.factor)
  if (!all(lengths(INDEX) == length(X))) stop("arguments must have same length")
  levs <- lapply(INDEX, levels)
  extent <- lengths(levs)
  code <- rep(1, length(X))
  stride <- 1
  for (k in seq_along(INDEX)) {
    code <- code + (as.integer(INDEX[[k]]) - 1) * stride
    stride <- stride * extent[k]
  }
  if (is.null(FUN)) return(as.integer(code))
  cells <- split(X, factor(code, levels = seq_len(prod(extent))))
  seen <- seq_along(cells) %in% code
  res <- vector("list", length(cells))
  for (i in which(seen)) res[i] <- list(FUN(cells[[i]], ...))
  shape <- function(v) {
    dim(v) <- extent
    dimnames(v) <- levs
    v
  }
  if (simplify && all(lengths(res[seen]) == 1L)) {
    vals <- unlist(res[seen], use.names = FALSE)
    out <- if (is.na(default) && is.atomic(vals) && !is.null(vals)) rep(vals[NA_integer_], length(cells)) else rep(default, length(cells))
    out[seen] <- vals
    return(shape(out))
  }
  shape(res)
}

bquote <- function(expr, where = parent.frame(), splice = FALSE) {
  rewrite <- function(e) {
    if (!is.call(e)) return(e)
    if (identical(e[[1L]], quote(.)) && length(e) == 2L) return(eval(e[[2L]], where))
    parts <- as.list(e)
    labels <- names(parts)
    if (is.null(labels)) labels <- character(length(parts))
    out <- list()
    outLabels <- character()
    for (i in seq_along(parts)) {
      p <- parts[[i]]
      if (splice && is.call(p) && length(p) == 2L && identical(p[[1L]], quote(..))) {
        vals <- as.list(eval(p[[2L]], where))
        valLabels <- names(vals)
        if (is.null(valLabels)) valLabels <- character(length(vals))
        out <- c(out, vals)
        outLabels <- c(outLabels, valLabels)
      } else {
        out <- c(out, list(rewrite(p)))
        outLabels <- c(outLabels, labels[i])
      }
    }
    if (any(outLabels != "")) names(out) <- outLabels
    as.call(out)
  }
  rewrite(substitute(expr))
}

is.vector <- function(x, mode = "any") {
  a <- attributes(x)
  if (length(a) > 0L && !identical(names(a), "names")) return(FALSE)
  switch(mode,
         any = (is.atomic(x) || is.list(x)) && !is.null(x),
         numeric = is.numeric(x),
         name = is.symbol(x),
         symbol = is.symbol(x),
         identical(typeof(x), mode))
}
