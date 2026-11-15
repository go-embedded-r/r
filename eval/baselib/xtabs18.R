# Base library, slice 18 (2): xtabs and summary.table. Written from the help
# pages of R 4.6.2 (xtabs, summary.table, table) and from experiments with the
# oracle. summary.table does not run the chi-squared test: it needs pchisq,
# which is not implemented (README, divergences).

xtabs <- function(formula = ~., data = parent.frame(), subset, sparse = FALSE, na.action, addNA = FALSE,
                  exclude = if (!addNA) c(NA, NaN), drop.unused.levels = FALSE) {
  if (missing(formula) && missing(data)) stop("must supply either 'formula' or 'data'")
  if (!inherits(formula, "formula")) stop("'formula' missing or incorrect")
  if (sparse) stop("sparse tables are not supported")
  cl <- match.call()
  rows <- if (missing(subset)) NULL else eval(substitute(subset), data, parent.frame())
  act <- if (!missing(na.action)) na.action else na.pass
  mf <- eval(bquote(model.frame(formula, data, subset = .(rows), na.action = .(act),
                                drop.unused.levels = .(drop.unused.levels))))
  has.y <- length(formula) == 3L
  vars <- if (has.y) mf[-1L] else mf
  y <- if (has.y) model.response(mf) else NULL
  if (length(vars) == 0L) stop("no factors in the formula")
  facs <- lapply(vars, function(v) factor(v, exclude = exclude))
  dims <- vapply(facs, nlevels, 1L)
  dn <- lapply(facs, function(f) if (nlevels(f)) levels(f))
  total <- prod(dims)
  if (total > 2^27) stop("table is too large")
  idx <- rep(1L, length(facs[[1L]]))
  stride <- 1L
  for (k in seq_along(facs)) {
    code <- as.integer(facs[[k]])
    idx <- ifelse(is.na(code) | is.na(idx), NA_integer_, idx + (code - 1L) * stride)
    stride <- stride * dims[k]
  }
  keep <- !is.na(idx)
  idx <- idx[keep]
  cells <- function(w) {
    out <- if (is.integer(w)) integer(total) else numeric(total)
    if (length(idx)) {
      r <- rowsum(w, idx)
      out[as.integer(rownames(r))] <- r
    }
    out
  }
  if (is.null(y)) {
    counts <- tabulate(idx, nbins = total)
    ans <- array(counts, dim = dims, dimnames = dn)
  } else if (is.matrix(y)) {
    cols <- lapply(seq_len(ncol(y)), function(j) cells(y[keep, j]))
    ans <- array(unlist(cols), dim = c(dims, ncol(y)), dimnames = c(dn, list(colnames(y))))
  } else {
    ans <- array(cells(if (is.factor(y)) as.integer(y)[keep] else y[keep]), dim = dims, dimnames = dn)
  }
  class(ans) <- c("xtabs", "table")
  attr(ans, "call") <- cl
  ans
}

print.xtabs <- function(x, na.print = "", ...) {
  ox <- x
  attr(x, "call") <- NULL
  print.table(x, na.print = na.print, ...)
  invisible(ox)
}

summary.table <- function(object, ...) {
  if (!is.table(object)) stop(gettextf("'object' must inherit from class %s", dQuote("table")), domain = NA)
  n.cases <- sum(object)
  n.vars <- length(dim(object))
  y <- list(n.vars = n.vars, n.cases = n.cases)
  if (n.vars > 1) y$call <- attr(object, "call")
  class(y) <- "summary.table"
  y
}

print.summary.table <- function(x, digits = max(1L, getOption("digits") - 3L), ...) {
  if (!is.null(x$call)) {
    cat("Call: ")
    print(x$call)
  }
  cat("Number of cases in table:", x$n.cases, "\n")
  cat("Number of factors:", x$n.vars, "\n")
  invisible(x)
}
