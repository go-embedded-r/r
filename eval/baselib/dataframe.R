# Base library, part 9: data frames. Written from the help pages of
# data.frame, as.data.frame, row.names, Extract.data.frame, rbind.data.frame,
# merge, subset, transform, droplevels and the factor methods in R 4.6.2 and
# checked against the oracle (docs/m3-plan.md). BSD licence.

..mkdf <- function(cols, nms, rn, rnsecond = FALSE) {
  if (is.na(rnsecond)) {
    attr(cols, "row.names") <- rn
    names(cols) <- nms
    attr(cols, "class") <- "data.frame"
    return(cols)
  }
  names(cols) <- nms
  if (rnsecond) {
    attr(cols, "row.names") <- rn
    attr(cols, "class") <- "data.frame"
  } else {
    attr(cols, "class") <- "data.frame"
    attr(cols, "row.names") <- rn
  }
  cols
}

..dfcols <- function(x) {
  nm <- names(x)
  x <- unclass(x)
  attr(x, "row.names") <- NULL
  names(x) <- nm
  x
}

is.data.frame <- function(x) inherits(x, "data.frame")

# the store step of `[<-.data.frame` (a Go native, eval/dfassign.go)
..dfsetcol <- function(col, ii, v, has.i, newn) {
  if (is.null(v)) return(col[seq_len(newn)])
  if (has.i) col[ii] <- v else col[seq_len(newn)] <- v
  col
}

dim.data.frame <- function(x) c(.row_names_info(x, 2L), length(x))

dimnames.data.frame <- function(x) list(row.names(x), names(x))

`dimnames<-.data.frame` <- function(x, value) {
  if (!is.list(value) || length(value) != 2L)
    stop("invalid 'dimnames' given for data frame")
  row.names(x) <- value[[1L]]
  names(x) <- as.character(value[[2L]])
  x
}

row.names <- function(x) UseMethod("row.names")
row.names.data.frame <- function(x) as.character(attr(x, "row.names"))
row.names.default <- function(x) if (!is.null(dim(x))) rownames(x)

`row.names<-` <- function(x, value) UseMethod("row.names<-")
`row.names<-.default` <- function(x, value) {
  if (!is.null(dim(x))) dimnames(x)[[1L]] <- value else attr(x, "row.names") <- value
  x
}
`row.names<-.data.frame` <- function(x, value) `.rowNamesDF<-`(x, value = value)

`.rowNamesDF<-` <- function(x, make.names = FALSE, value) {
  n <- .row_names_info(x, 2L)
  if (is.null(value)) {
    attr(x, "row.names") <- .set_row_names(n)
    return(x)
  }
  if (length(value) != n) {
    if (n == 0L) stop("invalid 'row.names' length")
    stop("invalid 'row.names' length")
  }
  if (is.object(value) || !is.integer(value)) value <- as.character(value)
  if (anyDuplicated(value)) {
    dup <- value[duplicated(value)]
    warning(sprintf("non-unique value when setting 'row.names': %s",
                    sQuote(dup[1L])), call. = FALSE)
    stop("duplicate 'row.names' are not allowed")
  }
  if (anyNA(value)) stop("missing values in 'row.names' are not allowed")
  attr(x, "row.names") <- value
  x
}

as.list.data.frame <- function(x, ...) {
  nm <- names(x)
  x <- unclass(x)
  attr(x, "row.names") <- NULL
  names(x) <- nm
  x
}

# ---- factors ----------------------------------------------------------

c.factor <- function(..., recursive = TRUE) {
  x <- list(...)
  if (!all(vapply(x, is.factor, NA))) return(unlist(lapply(x, unclass)))
  lev <- unique(unlist(lapply(x, levels)))
  codes <- unlist(lapply(x, function(f) match(levels(f)[f], lev)))
  structure(codes, levels = lev, class = "factor")
}

droplevels <- function(x, ...) UseMethod("droplevels")
droplevels.factor <- function(x, exclude = if (anyNA(levels(x))) NULL else NA, ...)
  factor(x, exclude = exclude)
droplevels.data.frame <- function(x, except = NULL, exclude, ...) {
  ix <- !(seq_along(x) %in% except)
  for (i in which(ix)) if (is.factor(x[[i]])) x[[i]] <- droplevels(x[[i]])
  x
}

# ---- building ---------------------------------------------------------

..dfpiece <- function(xi, lab, named, saf) {
  if (is.list(xi) && !is.data.frame(xi))
    xi <- do.call(data.frame, c(xi, list(check.names = FALSE, stringsAsFactors = saf)))
  if (is.data.frame(xi)) {
    cl <- ..dfcols(xi)
    nm <- names(cl)
    if (named && length(cl) > 1L) nm <- paste(lab, nm, sep = ".")
    rn <- if (.row_names_info(xi) > 0L) attr(xi, "row.names")
    return(list(cols = cl, names = nm, rn = rn, nr = .row_names_info(xi, 2L)))
  }
  if (is.matrix(xi)) {
    nc <- ncol(xi)
    cn <- colnames(xi)
    cl <- lapply(seq_len(nc), function(k) {
      v <- xi[, k]
      names(v) <- NULL
      v
    })
    nm <- if (!is.null(cn)) cn else if (named) as.character(seq_len(nc)) else
      paste0("X", seq_len(nc))
    if (named && (!is.null(cn) || nc > 1L)) nm <- paste(lab, nm, sep = ".")
    if (named && nc == 1L) nm <- lab
    if (saf) cl <- lapply(cl, function(v) if (is.character(v)) factor(v) else v)
    return(list(cols = cl, names = nm, rn = rownames(xi), nr = nrow(xi)))
  }
  rn <- names(xi)
  if (!is.null(rn) && !is.object(xi)) names(xi) <- NULL
  if (saf && is.character(xi)) xi <- factor(xi)
  list(cols = list(xi), names = lab, rn = rn, nr = NROW(xi))
}

data.frame <- function(..., row.names = NULL, check.rows = FALSE, check.names = TRUE,
                       fix.empty.names = TRUE, stringsAsFactors = FALSE) {
  x <- list(...)
  nx <- length(x)
  labs <- names(x)
  if (is.null(labs)) labs <- character(nx)
  exprs <- as.list(substitute(list(...)))[-1L]
  cols <- list()
  nms <- character()
  rn <- NULL
  nrs <- integer()
  argnr <- integer()
  for (i in seq_len(nx)) {
    xi <- x[[i]]
    named <- !is.na(labs[i]) && nzchar(labs[i])
    if (is.null(xi)) {
      if (named) argnr <- c(argnr, 0L)
      next
    }
    lab <- if (named) labs[i] else deparse(exprs[[i]])[1L]
    p <- ..dfpiece(xi, lab, named, stringsAsFactors)
    cols <- c(cols, p$cols)
    nms <- c(nms, p$names)
    nrs <- c(nrs, rep(p$nr, length(p$cols)))
    argnr <- c(argnr, p$nr)
    if (is.null(rn) && !is.null(p$rn) && length(p$rn) == p$nr && !anyDuplicated(p$rn))
      rn <- p$rn
  }
  if (length(cols) == 0L) {
    if (!is.null(row.names)) return(..mkdf(list(), character(), row.names, TRUE))
    return(..mkdf(list(), character(), integer(), TRUE))
  }
  nr <- max(nrs)
  if (any(nrs != nr & (nrs == 0L | nr %% nrs != 0L)) || (nr > 0L && any(argnr == 0L)))
    stop(gettextf("arguments imply differing number of rows: %s",
                  paste(unique(argnr), collapse = ", ")), domain = NA)
  for (k in seq_along(cols))
    if (nrs[k] != nr) cols[[k]] <- cols[[k]][rep_len(seq_len(nrs[k]), nr)]
  if (!is.null(row.names)) {
    if (length(row.names) == 1L && nr != 1L) {
      k <- if (is.character(row.names)) match(row.names, nms) else row.names
      if (is.na(k) || k < 1L || k > length(cols)) stop("'row.names' should specify one of the variables")
      rn <- cols[[k]]
      cols <- cols[-k]
      nms <- nms[-k]
    } else {
      if (length(row.names) != nr) stop("row names supplied are of the wrong length")
      rn <- row.names
    }
  }
  if (check.names) nms <- make.names(nms, unique = TRUE)
  if (is.null(rn)) rn <- .set_row_names(nr)
  else {
    if (anyDuplicated(rn)) stop(gettextf("duplicate row.names: %s",
                                         paste(unique(rn[duplicated(rn)]), collapse = ", ")), domain = NA)
    if (!is.integer(rn)) rn <- as.character(rn)
  }
  ..mkdf(cols, nms, rn)
}

as.data.frame <- function(x, row.names = NULL, optional = FALSE, ...)
  UseMethod("as.data.frame")

as.data.frame.data.frame <- function(x, row.names = NULL, ...) {
  if (!is.null(row.names)) row.names(x) <- row.names
  x
}

as.data.frame.list <- function(x, row.names = NULL, optional = FALSE, ...,
                               cut.names = FALSE, col.names = names(x),
                               fix.empty.names = TRUE, check.names = !optional,
                               stringsAsFactors = FALSE) {
  if (is.null(col.names)) col.names <- paste0("V", seq_along(x))
  names(x) <- col.names
  y <- do.call(data.frame, c(x, list(check.names = check.names,
                                      stringsAsFactors = stringsAsFactors)))
  if (!is.null(row.names)) row.names(y) <- row.names
  y
}

as.data.frame.matrix <- function(x, row.names = NULL, optional = FALSE, make.names = TRUE,
                                 ..., stringsAsFactors = FALSE) {
  nr <- nrow(x)
  nc <- ncol(x)
  rn <- rownames(x)
  cn <- colnames(x)
  if (is.null(cn)) cn <- paste0("V", seq_len(nc))
  cols <- lapply(seq_len(nc), function(k) {
    v <- x[, k]
    names(v) <- NULL
    if (stringsAsFactors && is.character(v)) factor(v) else v
  })
  if (is.null(row.names)) row.names <- rn
  if (is.null(row.names) || anyDuplicated(row.names)) row.names <- .set_row_names(nr)
  ..mkdf(cols, cn, row.names)
}

as.data.frame.vector <- function(x, row.names = NULL, optional = FALSE, ..., nm = deparse1(substitute(x))) {
  force(nm)
  n <- length(x)
  if (is.null(row.names)) {
    row.names <- names(x)
    if (length(row.names) != n || anyDuplicated(row.names)) row.names <- .set_row_names(n)
  }
  if (!is.null(names(x))) names(x) <- NULL
  value <- list(x)
  if (!optional) names(value) <- nm
  attr(value, "row.names") <- row.names
  class(value) <- "data.frame"
  value
}
..dfvec <- function(cls) {
  cls
  function(x, row.names = NULL, optional = FALSE, ..., nm = deparse1(substitute(x))) {
    force(nm)
    if (!exists(".Generic", inherits = FALSE))
      warning(gettextf("Direct call of '%s()' is deprecated.  Use '%s()' or '%s()' instead",
                       paste0("as.data.frame.", cls), "as.data.frame.vector", "as.data.frame"),
              call. = FALSE, domain = NA)
    as.data.frame.vector(x, row.names, optional, ..., nm = nm)
  }
}
as.data.frame.integer <- ..dfvec("integer")
as.data.frame.numeric <- ..dfvec("numeric")
as.data.frame.character <- as.data.frame.vector
as.data.frame.logical <- ..dfvec("logical")
as.data.frame.factor <- ..dfvec("factor")
as.data.frame.complex <- ..dfvec("complex")
as.data.frame.Date <- ..dfvec("Date")
as.data.frame.POSIXct <- ..dfvec("POSIXct")
as.data.frame.default <- function(x, ...)
  stop(gettextf("cannot coerce class %s to a data.frame", sQuote(deparse(class(x))[1L])),
       domain = NA)

# ---- extraction -------------------------------------------------------

`[.data.frame` <- function(x, i, j, drop = if (missing(i)) TRUE else length(cols) == 1) {
  na <- nargs() - !missing(drop)
  cols <- ..dfcols(x)
  nr <- .row_names_info(x, 2L)
  rnx <- attr(x, "row.names")
  auto <- .row_names_info(x) < 0L
  if (na < 3L) {
    if (missing(i)) return(x)
    cols <- cols[i]
    if (anyNA(names(cols))) stop("undefined columns selected")
    if (anyDuplicated(names(cols))) names(cols) <- make.unique(names(cols))
    return(..mkdf(unname(cols), names(cols), .row_names_info(x, 0L), TRUE))
  }
  if (!missing(j)) {
    cols <- cols[j]
    if (anyNA(names(cols))) stop("undefined columns selected")
    if (anyDuplicated(names(cols))) names(cols) <- make.unique(names(cols))
  }
  if (!missing(i)) {
    if (is.character(i)) i <- pmatch(i, as.character(rnx), duplicates.ok = TRUE)
    ix <- seq_len(nr)[i]
    cols <- lapply(cols, function(col) col[ix])
    rn <- rnx[ix]
    if (anyNA(rn) || anyDuplicated(rn)) {
      rn <- as.character(rn)
      rn[is.na(rn)] <- "NA"
      rn <- make.unique(rn)
    }
    nr <- length(ix)
  } else rn <- .row_names_info(x, 0L)
  if (drop) {
    if (length(cols) == 1L) return(cols[[1L]])
    if (nr == 1L && !missing(drop)) return(cols)
  }
  ..mkdf(unname(cols), names(cols), rn, if (missing(i)) FALSE else TRUE)
}
