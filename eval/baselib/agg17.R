# Base library, slice 17 (2): what the formulas make possible. Written from the
# help pages of R 4.6.2 (aggregate, stack, unstack, xtabs, sort_by) and from
# experiments with the oracle.

# ---- aggregate ---------------------------------------------------------------

aggregate <- function(x, ...) UseMethod("aggregate")

registerS3method("aggregate", "default", function(x, ...) aggregate.data.frame(as.data.frame(x), ...))

aggregate.data.frame <- function(x, by, FUN, ..., simplify = TRUE, drop = TRUE) {
  if (!is.data.frame(x)) x <- as.data.frame(x)
  FUN <- match.fun(FUN)
  if (NROW(x) == 0L) stop("no rows to aggregate")
  if (!is.list(by)) stop("'by' must be a list")
  labels <- names(by)
  if (is.null(labels)) labels <- rep("", length(by))
  labels[!nzchar(labels)] <- paste0("Group.", seq_along(by))[!nzchar(labels)]
  if (any(lengths(by) != NROW(x))) stop("arguments must have same length")
  # one code per grouping variable and row; the first variable varies fastest
  kept <- vector("list", length(by))
  codes <- vector("list", length(by))
  for (j in seq_along(by)) {
    b <- by[[j]]
    if (is.factor(b)) {
      lev <- levels(b)
      code <- as.integer(b)
      if (drop) {
        used <- sort(unique(code[!is.na(code)]))
        code <- match(code, used)
        kept[[j]] <- factor(lev[used], levels = lev)
      } else {
        kept[[j]] <- factor(lev, levels = lev)
      }
    } else {
      lev <- sort(unique(b[!is.na(b)]))
      code <- match(b, lev)
      kept[[j]] <- lev
    }
    codes[[j]] <- code
  }
  sizes <- vapply(kept, length, 1L)
  if (prod(as.numeric(sizes)) > 1e7) stop("too many groups")
  id <- rep(1, NROW(x))
  stride <- 1
  for (j in seq_along(by)) {
    id <- id + (codes[[j]] - 1) * stride
    stride <- stride * sizes[j]
  }
  row_ok <- !is.na(id)
  groups <- if (drop) sort(unique(id[row_ok])) else seq_len(stride)
  out <- list()
  below <- 1
  for (j in seq_along(by)) {
    pos <- ((groups - 1) %/% below) %% sizes[j] + 1
    out[[labels[j]]] <- kept[[j]][pos]
    below <- below * sizes[j]
  }
  pieces <- factor(id[row_ok], levels = groups)
  for (nm in names(x)) {
    col <- x[[nm]][row_ok]
    got <- lapply(split(col, pieces), function(part) if (length(part)) FUN(part, ...))
    empty <- vapply(got, is.null, NA)
    if (any(empty)) got[empty] <- list(NA)
    one <- vapply(got, function(r) is.atomic(r) && length(r) == 1L, NA)
    if (simplify && all(one)) {
      out[[nm]] <- unlist(got, use.names = FALSE)
    } else if (simplify && length(unique(lengths(got))) == 1L && all(vapply(got, is.atomic, NA))) {
      out[[nm]] <- do.call(rbind, unname(got))
    } else {
      out[[nm]] <- unname(got)
    }
  }
  attr(out, "row.names") <- c(NA_integer_, -length(groups))
  class(out) <- "data.frame"
  out
}

registerS3method("aggregate", "formula", function(x, data, FUN, ..., subset, na.action = na.omit) {
  if (missing(x) || !inherits(x, "formula")) stop("argument 'x' is missing or not a formula")
  if (length(x) != 3L) stop("formula 'x' must have both left and right hand sides")
  lhs <- x[[2L]]
  rhs <- x[[3L]]
  sum_of <- function(nms) {
    out <- as.name(nms[1L])
    for (nm in nms[-1L]) out <- call("+", out, as.name(nm))
    out
  }
  if (identical(lhs, quote(.))) {
    lhs <- as.call(c(quote(cbind), lapply(setdiff(names(data), all.vars(rhs)), as.name)))
  } else if (identical(rhs, quote(.))) {
    rhs <- sum_of(setdiff(names(data), all.vars(lhs)))
  }
  made <- eval(call("~", lhs, rhs))
  environment(made) <- environment(x)
  ask <- list(as.name("model.frame"), made, quote(data), na.action = quote(na.action))
  if (!missing(subset)) ask$subset <- substitute(subset)
  frame <- eval(as.call(ask))
  vals <- frame[[1L]]
  vals <- if (is.matrix(vals)) as.data.frame(vals) else setNames(list(vals), names(frame)[1L])
  by <- as.list(frame[-1L])
  if (!length(by)) {
    fun <- match.fun(FUN)
    one <- lapply(vals, function(col) fun(col, ...))
    attr(one, "row.names") <- c(NA_integer_, -1L)
    class(one) <- "data.frame"
    return(one)
  }
  aggregate.data.frame(as.data.frame(vals, optional = TRUE), by, FUN, ...)
})

# repeating a factor repeats its codes and keeps its levels and class
rep.factor <- function(x, ...) {
  y <- NextMethod()
  structure(y, class = class(x), levels = levels(x))
}

# ---- stack and unstack --------------------------------------------------------

..by_levels <- function(values, groups) {
  got <- tapply(values, groups, as.vector)
  nm <- names(got)
  got <- unclass(got)
  dim(got) <- NULL
  names(got) <- nm
  got
}

stack <- function(x, ...) UseMethod("stack")

registerS3method("stack", "default", function(x, drop = FALSE, ...) {
  x <- as.list(x)
  keep <- vapply(x, is.vector, NA)
  if (!any(keep)) stop("at least one vector element is required")
  if (!all(keep)) warning("non-vector elements will be ignored")
  x <- x[keep]
  ind <- rep(factor(names(x), unique(names(x))), lengths(x))
  if (drop) ind <- droplevels(ind)
  data.frame(values = unlist(unname(x)), ind, stringsAsFactors = FALSE)
})

registerS3method("stack", "data.frame", function(x, select, drop = FALSE, ...) {
  if (!missing(select)) {
    pos <- as.list(seq_along(x))
    names(pos) <- names(x)
    x <- x[, eval(substitute(select), pos, parent.frame()), drop = FALSE]
  }
  keep <- vapply(x, is.vector, NA)
  if (!any(keep)) stop("at least one vector element is required")
  if (!all(keep)) warning("non-vector columns will be ignored")
  x <- x[keep]
  ind <- rep(factor(names(x), unique(names(x))), lengths(x))
  if (drop) ind <- droplevels(ind)
  data.frame(values = unlist(unname(x)), ind, stringsAsFactors = FALSE)
})

unstack <- function(x, ...) UseMethod("unstack")

registerS3method("unstack", "data.frame", function(x, form, ...) {
  form <- if (missing(form)) formula(x) else as.formula(form)
  if (length(form) != 3L) stop("'form' must be a two-sided formula")
  res <- ..by_levels(eval(form[[2L]], x), eval(form[[3L]], x))
  if (length(res) >= 2L && any(diff(lengths(res)) != 0L)) return(res)
  data.frame(res, stringsAsFactors = FALSE)
})

registerS3method("unstack", "default", function(x, form, ...) {
  x <- as.list(x)
  form <- as.formula(form)
  res <- ..by_levels(eval(form[[2L]], x), eval(form[[3L]], x))
  if (length(res) >= 2L && any(diff(lengths(res)) != 0L)) return(res)
  data.frame(res, stringsAsFactors = FALSE)
})
