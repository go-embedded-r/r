# Base library, part 7: matrices and arrays. Written from the help pages of
# apply, sweep, scale, aperm, prop.table, marginSums, lower.tri, diag,
# rownames, kronecker and outer in R 4.6.2 and checked against the oracle
# (docs/m3-plan.md). BSD licence.

as.matrix <- function(x, ...) UseMethod("as.matrix")

as.matrix.default <- function(x, ...) {
  if (is.matrix(x)) return(x)
  nms <- names(x)
  array(x, c(length(x), 1L), if (is.null(nms)) NULL else list(nms, NULL))
}

mat.or.vec <- function(nr, nc) if (nc == 1L) numeric(nr) else matrix(0, nr, nc)

lower.tri <- function(x, diag = FALSE) {
  x <- as.matrix(x)
  if (diag) row(x) >= col(x) else row(x) > col(x)
}

upper.tri <- function(x, diag = FALSE) {
  x <- as.matrix(x)
  if (diag) row(x) <= col(x) else row(x) < col(x)
}

`diag<-` <- function(x, value) {
  dx <- dim(x)
  if (length(dx) != 2L) stop("only matrix diagonals can be replaced")
  n <- min(dx)
  if (length(value) != 1L && length(value) != n)
    stop("replacement diagonal has wrong length")
  if (n > 0L) {
    k <- seq_len(n)
    x[(k - 1L) * dx[1L] + k] <- value
  }
  x
}

rownames <- function(x, do.NULL = TRUE, prefix = "row") {
  dn <- dimnames(x)
  if (!is.null(dn[[1L]])) return(dn[[1L]])
  if (do.NULL) return(NULL)
  n <- NROW(x)
  if (n > 0L) paste0(prefix, seq_len(n)) else character()
}

colnames <- function(x, do.NULL = TRUE, prefix = "col") {
  dn <- dimnames(x)
  if (!is.null(dn[[2L]])) return(dn[[2L]])
  if (do.NULL) return(NULL)
  n <- NCOL(x)
  if (n > 0L) paste0(prefix, seq_len(n)) else character()
}

..setdimnames <- function(x, which, value) {
  dn <- dimnames(x)
  if (is.null(dn)) {
    if (is.null(value)) return(x)
    dn <- vector("list", length(dim(x)))
  }
  if (is.null(value)) dn[which] <- list(NULL) else dn[[which]] <- value
  dimnames(x) <- dn
  x
}

`rownames<-` <- function(x, value) {
  if (length(dim(x)) < 1L) stop("attempt to set 'rownames' on an object with no dimensions")
  ..setdimnames(x, 1L, value)
}

`colnames<-` <- function(x, value) {
  if (length(dim(x)) < 2L) stop("attempt to set 'colnames' on an object with less than two dimensions")
  ..setdimnames(x, 2L, value)
}

row <- function(x, as.factor = FALSE) {
  d <- dim(x)
  if (length(d) != 2L) stop("a matrix-like object is required as argument to 'row'")
  res <- array(rep(seq_len(d[1L]), d[2L]), d)
  if (as.factor) {
    res <- factor(res, levels = seq_len(d[1L]), labels = rownames(x, do.NULL = FALSE, prefix = ""))
    dim(res) <- d
  }
  res
}

col <- function(x, as.factor = FALSE) {
  d <- dim(x)
  if (length(d) != 2L) stop("a matrix-like object is required as argument to 'col'")
  res <- array(rep(seq_len(d[2L]), each = d[1L]), d)
  if (as.factor) {
    res <- factor(res, levels = seq_len(d[2L]), labels = colnames(x, do.NULL = FALSE, prefix = ""))
    dim(res) <- d
  }
  res
}

aperm <- function(a, perm = NULL, ...) UseMethod("aperm")

aperm.default <- function(a, perm = NULL, resize = TRUE, ...) {
  d <- dim(a)
  nd <- length(d)
  if (nd == 0L) stop("invalid first argument, must be an array")
  dn <- dimnames(a)
  if (is.null(perm)) perm <- rev(seq_len(nd))
  if (is.character(perm)) {
    if (is.null(names(dn))) stop("'a' does not have named dimnames")
    perm <- match(perm, names(dn))
  }
  if (length(perm) != nd) stop("'perm' is of wrong length ", length(perm), " (!= ", nd, ")")
  if (!identical(sort(as.integer(perm)), seq_len(nd))) stop("invalid 'perm' argument")
  perm <- as.integer(perm)
  nd2 <- d[perm]
  strides <- c(1, cumprod(d))[seq_len(nd)]
  pos <- seq_len(prod(d)) - 1
  src <- rep(0, length(pos))
  step <- 1
  for (k in seq_len(nd)) {
    src <- src + ((pos %/% step) %% nd2[k]) * strides[perm[k]]
    step <- step * nd2[k]
  }
  res <- a[src + 1]
  if (resize) {
    dim(res) <- nd2
    if (!is.null(dn)) dimnames(res) <- dn[perm]
  } else dim(res) <- d
  res
}

apply <- function(X, MARGIN, FUN, ..., simplify = TRUE) {
  FUN <- match.fun(FUN)
  d <- dim(X)
  if (!length(d)) stop("dim(X) must have a positive length")
  dn <- dimnames(X)
  if (is.character(MARGIN)) {
    if (is.null(names(dn))) stop("'X' must have named dimnames")
    MARGIN <- match(MARGIN, names(dn))
    if (anyNA(MARGIN)) stop("not all elements of 'MARGIN' are names of dimensions")
  }
  if (!all(MARGIN %in% seq_along(d))) stop("'MARGIN' does not match dim(X)")
  keep <- MARGIN
  rest <- seq_along(d)[-keep]
  cells <- prod(d[keep])
  size <- prod(d[rest])
  flat <- aperm(X, c(rest, keep))
  dim(flat) <- c(size, cells)
  restNames <- if (is.null(dn)) NULL else dn[rest]
  out <- vector("list", cells)
  for (j in seq_len(cells)) {
    piece <- flat[, j]
    if (length(rest) == 1L) {
      if (!is.null(restNames[[1L]])) names(piece) <- restNames[[1L]]
    } else if (length(rest) > 1L) {
      piece <- array(piece, d[rest], restNames)
    }
    out[j] <- list(FUN(piece, ...))
  }
  keepNames <- if (is.null(dn)) NULL else dn[keep]
  lens <- unique(lengths(out))
  if (simplify && length(out) > 0L && length(lens) == 1L && lens == 0L && !is.recursive(out[[1L]]))
    return(unlist(out))
  if (simplify && length(out) > 0L && length(lens) == 1L && lens >= 1L && !is.recursive(out[[1L]])) {
    vals <- unlist(out, use.names = FALSE)
    if (lens == 1L) {
      if (length(keep) == 1L) {
        names(vals) <- keepNames[[1L]]
        return(vals)
      }
      return(array(vals, d[keep], keepNames))
    }
    firstNames <- names(out[[1L]])
    allNames <- c(list(firstNames), keepNames)
    if (all(vapply(allNames, is.null, NA))) allNames <- NULL
    return(array(vals, c(lens, d[keep]), allNames))
  }
  if (length(keep) == 1L) {
    names(out) <- keepNames[[1L]]
    return(out)
  }
  dim(out) <- d[keep]
  if (!is.null(keepNames)) dimnames(out) <- keepNames
  out
}

sweep <- function(x, MARGIN, STATS, FUN = "-", check.margin = TRUE, ...) {
  FUN <- match.fun(FUN)
  d <- dim(x)
  if (is.character(MARGIN)) {
    MARGIN <- match(MARGIN, names(dimnames(x)))
  }
  if (check.margin) {
    extent <- prod(d[MARGIN])
    if (length(STATS) > extent)
      warning("STATS is longer than the extent of 'dim(x)[MARGIN]'")
    else if (extent %% length(STATS) != 0L)
      warning("STATS does not recycle exactly across MARGIN")
  }
  perm <- c(MARGIN, seq_along(d)[-MARGIN])
  FUN(x, aperm(array(STATS, d[perm]), order(perm)), ...)
}

scale <- function(x, center = TRUE, scale = TRUE) UseMethod("scale")

scale.default <- function(x, center = TRUE, scale = TRUE) {
  x <- as.matrix(x)
  if (!is.logical(center) && length(center) != ncol(x))
    stop("length of 'center' must equal the number of columns of 'x'")
  shift <- if (!is.logical(center)) center else if (isTRUE(center)) colMeans(x, na.rm = TRUE)
  if (!is.null(shift)) x <- x - rep(unname(shift), each = nrow(x))
  if (!is.logical(scale) && length(scale) != ncol(x))
    stop("length of 'scale' must equal the number of columns of 'x'")
  divisor <- if (!is.logical(scale)) scale else if (isTRUE(scale))
    apply(x, 2L, function(column) {
      column <- column[!is.na(column)]
      sqrt(sum(column^2) / max(1, length(column) - 1L))
    })
  if (!is.null(divisor)) x <- x / rep(unname(divisor), each = nrow(x))
  if (!is.null(shift)) attr(x, "scaled:center") <- shift
  if (!is.null(divisor)) attr(x, "scaled:scale") <- divisor
  x
}

marginSums <- function(x, margin = NULL) {
  if (!length(margin)) return(sum(x))
  res <- apply(x, margin, sum)
  dim(res) <- dim(x)[margin]
  dimnames(res) <- dimnames(x)[margin]
  res
}

margin.table <- marginSums

proportions <- function(x, margin = NULL) {
  if (length(margin)) sweep(x, margin, marginSums(x, margin), "/", check.margin = FALSE)
  else x / sum(x)
}

prop.table <- proportions

`%o%` <- function(X, Y) outer(X, Y)
