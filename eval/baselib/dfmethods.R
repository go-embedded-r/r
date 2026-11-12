# Base library, part 10: the methods of data frames. Written from the help
# pages of Extract.data.frame, format.data.frame, print.data.frame,
# rbind.data.frame, merge, split, subset, transform and duplicated in R 4.6.2
# and checked against the oracle (docs/m3-plan.md). BSD licence.

`[[.data.frame` <- function(x, ..., exact = TRUE) {
  na <- nargs() - !missing(exact)
  cols <- ..dfcols(x)
  if (na < 3L) return(cols[[..., exact = exact]])
  ij <- list(...)
  col <- cols[[ij[[2L]], exact = exact]]
  i <- ij[[1L]]
  if (is.character(i)) i <- match(i, attr(x, "row.names"))
  col[[i]]
}

`$<-.data.frame` <- function(x, name, value) {
  cl <- oldClass(x)
  class(x) <- NULL
  nrows <- .row_names_info(x, 2L)
  if (!is.null(value)) {
    n <- NROW(value)
    if (n > nrows || (n < nrows && (n == 0L || nrows %% n != 0L)))
      stop(sprintf(if (n == 1L) "replacement has %d row, data has %d" else
                   "replacement has %d rows, data has %d", n, nrows))
    if (n < nrows) value <- value[rep_len(seq_len(n), nrows)]
    if (is.atomic(value) && !is.null(names(value))) names(value) <- NULL
  }
  x[[name]] <- value
  class(x) <- cl
  x
}

`[[<-.data.frame` <- function(x, i, j, value) {
  if (nargs() < 4L) {
    if (is.numeric(i) && length(i) == 1L && i <= length(x)) {
      nm <- names(x)[i]
      if (is.null(value)) {
        cl <- oldClass(x)
        class(x) <- NULL
        x[[i]] <- NULL
        class(x) <- cl
        return(x)
      }
      i <- nm
    }
    cl <- oldClass(x)
    nrows <- .row_names_info(x, 2L)
    x <- unclass(x)
    if (!is.null(value)) {
      n <- NROW(value)
      if (n > nrows || (n < nrows && (n == 0L || nrows %% n != 0L)))
        stop(sprintf(if (n == 1L) "replacement has %d row, data has %d" else
                     "replacement has %d rows, data has %d", n, nrows))
      if (n < nrows) value <- value[rep_len(seq_len(n), nrows)]
    }
    x[[i]] <- value
    class(x) <- cl
    return(x)
  }
  col <- x[[j]]
  if (is.character(i)) i <- match(i, attr(x, "row.names"))
  col[[i]] <- value
  x[[j]] <- col
  x
}

# ---- format, print, matrix --------------------------------------------

format.data.frame <- function(x, ..., justify = "left") {
  cols <- lapply(..dfcols(x), function(col)
    if (is.character(col)) col else if (is.factor(col)) as.character(col) else format(col, ...))
  cols <- lapply(cols, function(col) {
    class(col) <- "AsIs"
    col
  })
  ..mkdf(unname(cols), names(x), row.names(x))
}

as.matrix.data.frame <- function(x, rownames.force = NA, ...) {
  nr <- .row_names_info(x, 2L)
  cols <- ..dfcols(x)
  rn <- if (identical(rownames.force, FALSE)) NULL else
    if (identical(rownames.force, TRUE) || .row_names_info(x) > 0L) row.names(x)
  if (any(vapply(cols, function(v) is.character(v) || is.factor(v), NA))) {
    cols <- lapply(cols, function(v) {
      if (is.character(v)) v else if (is.factor(v)) as.character(v)
      else if (is.logical(v)) as.character(v)
      else {
        gone <- is.na(v)
        out <- format(v)
        out[gone] <- NA
        out
      }
    })
  }
  if (length(cols) == 0L) return(matrix(NA, nr, 0L, dimnames = list(rn, NULL)))
  matrix(unlist(cols, use.names = FALSE), nr, length(cols), dimnames = list(rn, names(x)))
}

t.data.frame <- function(x) t(as.matrix(x))

print.data.frame <- function(x, ..., digits = NULL, quote = FALSE, right = TRUE,
                             row.names = TRUE, max = NULL) {
  n <- .row_names_info(x, 2L)
  if (length(x) == 0L) {
    cat("data frame with 0 columns and ", n, if (n == 1L) " row" else " rows", "\n", sep = "")
  } else if (n == 0L) {
    print(names(x), quote = FALSE)
    cat("<0 rows> (or 0-length row.names)\n")
  } else {
    m <- as.matrix(format.data.frame(x, digits = digits))
    if (identical(row.names, FALSE)) dimnames(m)[[1L]] <- rep("", n)
    print(m, quote = quote, right = right)
  }
  invisible(x)
}

is.na.data.frame <- function(x) {
  m <- matrix(unlist(lapply(..dfcols(x), is.na), use.names = FALSE), .row_names_info(x, 2L), length(x))
  dimnames(m) <- list(if (.row_names_info(x) > 0L) row.names(x), names(x))
  m
}

duplicated.data.frame <- function(x, incomparables = FALSE, fromLast = FALSE, ...) {
  if (length(x) == 0L) return(logical(nrow(x)))
  duplicated(do.call(paste, c(unname(..dfcols(x)), sep = "\r")), fromLast = fromLast)
}

unique.data.frame <- function(x, incomparables = FALSE, fromLast = FALSE, ...)
  x[!duplicated.data.frame(x, fromLast = fromLast), , drop = FALSE]

# ---- binding ----------------------------------------------------------

cbind.data.frame <- function(..., deparse.level = 1)
  data.frame(..., check.names = FALSE)

rbind.data.frame <- function(..., deparse.level = 1, make.row.names = TRUE,
                             stringsAsFactors = FALSE) {
  match.names <- function(clabs, nmi) {
    if (!all(nmi %in% clabs)) stop("names do not match previous names")
  }
  args <- list(...)
  tags <- names(args)
  if (is.null(tags)) tags <- character(length(args))
  keep <- !vapply(args, is.null, NA)
  args <- args[keep]
  tags <- tags[keep]
  isdf <- vapply(args, is.data.frame, NA)
  if (!any(isdf)) stop("no data frame among the arguments")
  empty <- isdf & vapply(args, function(a) length(a) == 0L, NA)
  args <- args[!empty]
  tags <- tags[!empty]
  isdf <- isdf[!empty]
  if (length(args) == 0L) return(data.frame())
  nms <- names(args[[which(isdf)[1L]]])
  nc <- length(nms)
  pieces <- vector("list", nc)
  rns <- vector("list", length(args))
  auto <- logical(length(args))
  for (a in seq_along(args)) {
    ar <- args[[a]]
    if (isdf[a]) {
      if (length(ar) != nc) stop("numbers of columns of arguments do not match")
      xi <- ar
      clabs <- nms
      match.names(clabs, names(xi))
      n <- .row_names_info(ar, 2L)
      auto[a] <- .row_names_info(ar) < 0L
      rns[[a]] <- if (auto[a]) seq_len(n) else row.names(ar)
      for (k in seq_len(nc)) pieces[[k]][[a]] <- ar[[nms[k]]]
    } else {
      n <- 1L
      auto[a] <- TRUE
      rns[[a]] <- 1L
      ar <- rep_len(as.list(ar), nc)
      for (k in seq_len(nc)) pieces[[k]][[a]] <- ar[[k]]
    }
  }
  cols <- lapply(pieces, function(pc) {
    if (any(vapply(pc, is.factor, NA))) {
      pc <- lapply(pc, function(p) if (is.factor(p)) p else factor(p))
      do.call(c, pc)
    } else do.call(c, pc)
  })
  total <- sum(lengths(rns))
  if (!make.row.names || (all(auto) && !any(nzchar(tags)))) {
    rn <- .set_row_names(total)
  } else {
    rn <- unlist(lapply(seq_along(args), function(a) {
      if (!auto[a]) return(as.character(rns[[a]]))
      if (nzchar(tags[a])) {
        if (length(rns[[a]]) == 1L) tags[a] else paste(tags[a], rns[[a]], sep = ".")
      } else as.character(rns[[a]] + sum(lengths(rns)[seq_len(a - 1L)]))
    }))
    rn <- make.unique(rn)
  }
  ..mkdf(cols, nms, rn, TRUE)
}

# ---- merge, split, subset, transform ----------------------------------

merge <- function(x, y, ...) UseMethod("merge")

merge.data.frame <- function(x, y, by = intersect(names(x), names(y)), by.x = by,
                             by.y = by, all = FALSE, all.x = all, all.y = all,
                             sort = TRUE, suffixes = c(".x", ".y"), ...) {
  nx <- nrow(x)
  ny <- nrow(y)
  if (is.numeric(by.x)) by.x <- names(x)[by.x]
  if (is.numeric(by.y)) by.y <- names(y)[by.y]
  if (length(by.x) != length(by.y)) stop("'by.x' and 'by.y' specify different numbers of columns")
  if (length(by.x) == 0L) {
    if (as.numeric(nx) * ny > 2^27) stop("the cross product is too large")
    ix <- rep(seq_len(nx), ny)
    iy <- rep(seq_len(ny), each = nx)
    ux <- uy <- integer()
  } else {
    if (!all(by.x %in% names(x)) || !all(by.y %in% names(y)))
      stop("'by' must specify a uniquely valid column")
    kx <- do.call(paste, c(unname(as.list(x[by.x])), sep = "\r"))
    ky <- do.call(paste, c(unname(as.list(y[by.y])), sep = "\r"))
    hits <- split(seq_len(ny), ky)[kx]
    lens <- vapply(hits, length, 1L)
    if (sum(as.numeric(lens)) > 2^27) stop("the result of 'merge' is too large")
    ix <- rep(seq_len(nx), lens)
    iy <- unlist(hits, use.names = FALSE)
    if (is.null(iy)) iy <- integer()
    ux <- if (all.x) which(lens == 0L) else integer()
    uy <- if (all.y) setdiff(seq_len(ny), iy) else integer()
  }
  rx <- setdiff(names(x), by.x)
  ry <- setdiff(names(y), by.y)
  both <- intersect(rx, ry)
  nmx <- ifelse(rx %in% both, paste0(rx, suffixes[1L]), rx)
  nmy <- ifelse(ry %in% both, paste0(ry, suffixes[2L]), ry)
  rowsx <- c(ix, ux, rep(NA_integer_, length(uy)))
  rowsy <- c(iy, rep(NA_integer_, length(ux)), uy)
  cols <- list()
  for (k in seq_along(by.x)) {
    v <- x[[by.x[k]]][c(ix, ux)]
    if (length(uy)) v <- c(v, y[[by.y[k]]][uy])
    cols[[by.x[k]]] <- v
  }
  for (k in seq_along(rx)) cols[[nmx[k]]] <- x[[rx[k]]][rowsx]
  for (k in seq_along(ry)) cols[[nmy[k]]] <- y[[ry[k]]][rowsy]
  if (sort && length(by.x) > 0L && length(rowsx) > 1L) {
    o <- do.call(order, unname(cols[seq_along(by.x)]))
    cols <- lapply(cols, function(v) v[o])
  }
  ..mkdf(unname(cols), names(cols), .set_row_names(length(rowsx)), length(by.x) > 0L)
}

split.data.frame <- function(x, f, drop = FALSE, ...) {
  idx <- split(seq_len(nrow(x)), f, drop = drop, ...)
  lapply(idx, function(rows) x[rows, , drop = FALSE])
}

subset <- function(x, ...) UseMethod("subset")
subset.default <- function(x, subset, ...) x[subset & !is.na(subset)]
subset.data.frame <- function(x, subset, select, drop = FALSE, ...)
  ..subsetDF(x, if (missing(subset)) NULL else eval(substitute(subset), x, parent.frame()),
            if (missing(select)) NULL else
              eval(substitute(select), as.list(structure(seq_along(x), names = names(x))), parent.frame()),
            drop)

transform <- function(`_data`, ...) UseMethod("transform")
transform.data.frame <- function(`_data`, ...) {
  e <- eval(substitute(list(...)), `_data`, parent.frame())
  r <- ..transformDF(`_data`, e)
  if (length(r[[2L]])) do.call("data.frame", c(list(r[[1L]]), r[[2L]])) else r[[1L]]
}
