# Slice 21, part A: the row and column helpers, and the methods of factor,
# noquote, table, listof and simple.list. Written from the help pages of R
# 4.6.2 and the oracle's answers (docs/m3-plan.md); BSD licence.

..sumwhy <- function(x, m, n) {
  if (!is.numeric(x) && !is.complex(x) && !is.logical(x)) return("'x' must be numeric or complex")
  if (is.na(m) || is.na(n) || m < 0 || n < 0) return("invalid 'n' or 'm' argument")
  if (length(x) < m * n) return("'x' is too short")
  NULL
}

..sumshape <- function(x, m, n) matrix(x[seq_len(m * n)], m, n)

.rowSums <- function(x, m, n, na.rm = FALSE) {
  if (!is.null(why <- ..sumwhy(x, m, n))) stop(why)
  rowSums(..sumshape(x, m, n), na.rm = na.rm)
}
.colSums <- function(x, m, n, na.rm = FALSE) {
  if (!is.null(why <- ..sumwhy(x, m, n))) stop(why)
  colSums(..sumshape(x, m, n), na.rm = na.rm)
}
.rowMeans <- function(x, m, n, na.rm = FALSE) {
  if (!is.null(why <- ..sumwhy(x, m, n))) stop(why)
  rowMeans(..sumshape(x, m, n), na.rm = na.rm)
}
.colMeans <- function(x, m, n, na.rm = FALSE) {
  if (!is.null(why <- ..sumwhy(x, m, n))) stop(why)
  colMeans(..sumshape(x, m, n), na.rm = na.rm)
}

.row <- function(dim) {
  if (length(dim) != 2L) stop("a matrix-like object is required as argument to 'row'")
  matrix(rep.int(seq_len(dim[1L]), dim[2L]), dim[1L], dim[2L])
}

.col <- function(dim) {
  if (length(dim) != 2L) stop("a matrix-like object is required as argument to 'col'")
  matrix(rep(seq_len(dim[2L]), each = dim[1L]), dim[1L], dim[2L])
}

..refactor <- function(codes, x) {
  attr(codes, "levels") <- attr(x, "levels")
  ct <- attr(x, "contrasts")
  if (!is.null(ct)) attr(codes, "contrasts") <- ct
  class(codes) <- oldClass(x)
  codes
}

`[.factor` <- function(x, ..., drop = FALSE) {
  y <- ..refactor(unclass(x)[...], x)
  if (drop) factor(y, exclude = if (anyNA(levels(x))) NULL else NA) else y
}

`[[.factor` <- function(x, ...) ..refactor(unclass(x)[[...]], x)

print.factor <- function(x, quote = FALSE, max.levels = NULL, width = getOption("width"), ...) {
  ord <- is.ordered(x)
  if (length(x) == 0L) {
    cat(if (ord) "ordered" else "factor", "()\n", sep = "")
  } else {
    shown <- as.character(x)
    extra <- attributes(x)
    extra <- extra[!(names(extra) %in% c("levels", "class", "contrasts"))]
    attributes(shown) <- extra
    print(shown, quote = quote, ...)
  }
  if (is.logical(max.levels)) max.levels <- if (isFALSE(max.levels)) 0L else NULL
  if (!identical(max.levels, 0L) && !identical(max.levels, 0)) {
    lev <- levels(x)
    if (quote) lev <- paste0("\"", lev, "\"")
    lev[is.na(lev)] <- "<NA>"
    nl <- length(lev)
    sep <- if (ord) " < " else " "
    if (is.null(max.levels)) {
      room <- width - (nchar("Levels: ", "w") + 3L + 1L + 3L)
      acc <- cumsum(nchar(lev, "w") + nchar(sep, "w"))
      max.levels <- if (nl <= 1L || acc[nl] <= room) nl else max(1L, which.max(acc > room) - 1L)
    }
    cut <- nl > max.levels
    body <- if (cut) c(lev[seq_len(max(1L, max.levels - 1L))], "...", if (max.levels > 1L) lev[nl]) else lev
    cat(if (cut) paste(format(nl), ""), "Levels: ", paste(body, collapse = sep), "\n", sep = "")
  }
  invisible(x)
}

format.factor <- function(x, ...) {
  shown <- as.character(x)
  names(shown) <- names(x)
  dim(shown) <- dim(x)
  dimnames(shown) <- dimnames(x)
  format(shown, ...)
}

as.vector.factor <- function(x, mode = "any") {
  if (mode == "list") return(as.list(x))
  if (mode == "any" || mode == "character") return(as.vector(as.character(x), mode))
  as.vector(unclass(x), mode)
}

as.logical.factor <- function(x, ...) as.logical(levels(x))[x]

c.noquote <- function(..., recursive = FALSE) {
  r <- do.call(c, lapply(list(...), unclass))
  class(r) <- "noquote"
  r
}

`[.noquote` <- function(x, ...) {
  a <- attributes(x)
  r <- unclass(x)[...]
  attributes(r) <- c(attributes(r), a[!(names(a) %in% c("dim", "dimnames", "names"))])
  r
}

as.matrix.noquote <- function(x, ...) noquote(as.matrix(unclass(x), ...))

`[.table` <- function(x, ..., drop = TRUE) {
  plain <- x
  oldClass(plain) <- NULL
  r <- plain[..., drop = drop]
  if (length(dim(x)) == 1L && is.null(dim(r)) && !is.null(names(r))) {
    nm <- names(r)
    dim(r) <- length(r)
    dn <- names(dimnames(x))
    dimnames(r) <- if (length(dn) && nzchar(dn)) structure(list(nm), names = dn) else list(nm)
  }
  if (length(dim(r)) > 0L) class(r) <- oldClass(x)
  r
}

aperm.table <- function(a, perm = NULL, resize = TRUE, keep.class = TRUE, ...) {
  r <- aperm.default(a, perm, resize = resize)
  if (keep.class) class(r) <- oldClass(a)
  r
}

print.listof <- function(x, ...) {
  nn <- names(x)
  n <- length(x)
  if (length(nn) != n) nn <- paste("Component", seq_len(n))
  for (i in seq_len(n)) {
    cat(nn[i], ":\n")
    print(x[[i]], ...)
    cat("\n")
  }
  invisible(x)
}

`[.listof` <- function(x, ...) {
  r <- unclass(x)[...]
  class(r) <- oldClass(x)
  r
}

print.simple.list <- function(x, ...) print(noquote(cbind("_" = unlist(x))), ...)

`[.simple.list` <- function(x, ...) {
  r <- unclass(x)[...]
  class(r) <- oldClass(x)
  r
}
