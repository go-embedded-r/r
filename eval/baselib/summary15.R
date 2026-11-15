# Base library, slice 15 (2): summary and its printing. Written from the help
# pages of R 4.6.2 (summary, quantile type 7, format.pval) and experiments with
# the oracle; quantile itself is in stats and is not here.

# the sample quantiles of type 7 of a vector without NA
..quantile7 <- function(x, probs) {
  n <- length(x)
  if (n == 0L) return(rep(NA_real_, length(probs)))
  x <- sort(as.numeric(x))
  at <- 1 + (n - 1) * probs
  lo <- floor(at)
  hi <- ceiling(at)
  frac <- at - lo
  below <- x[lo]
  above <- x[hi]
  ifelse(frac > 0 & below != above, (1 - frac) * below + frac * above, below)
}

..summaryClass <- function(x) {
  class(x) <- c("summaryDefault", "table")
  x
}

summary.default <- function(object, ..., digits, quantile.type = 7) {
  if (is.factor(object)) return(summary.factor(object, ...))
  if (is.numeric(object) || is.logical(object) && FALSE) {
    gone <- is.na(object)
    object <- object[!gone]
    q <- ..quantile7(object, c(0, 0.25, 0.5, 0.75, 1))
    q <- c(q[1L:3L], mean(object), q[4L:5L])
    if (!missing(digits)) q <- signif(q, digits)
    names(q) <- c("Min.", "1st Qu.", "Median", "Mean", "3rd Qu.", "Max.")
    if (any(gone)) q <- c(q, "NAs" = sum(gone))
    return(..summaryClass(q))
  }
  if (is.logical(object)) {
    tb <- table(object, exclude = NULL, useNA = "ifany")
    lab <- dimnames(tb)[[1L]]
    lab[is.na(lab)] <- "NAs"
    counts <- as.vector(unclass(tb))
    names(counts) <- lab
    return(..summaryClass(c(Mode = "logical", counts)))
  }
  if (is.character(object)) {
    gone <- is.na(object)
    seen <- object[!gone]
    v <- c(Length = length(object), N.unique = length(unique(seen)),
           N.blank = sum(seen == ""),
           Min.nchar = suppressWarnings(min(nchar(seen))),
           Max.nchar = suppressWarnings(max(nchar(seen))))
    if (any(gone)) v <- c(v, "NAs" = sum(gone))
    return(..summaryClass(v))
  }
  if (is.recursive(object) && !is.language(object) && (n <- length(object))) {
    cells <- array("", c(n, 3L), list(names(object), c("Length", "Class", "Mode")))
    for (i in seq_len(n)) {
      item <- object[[i]]
      cells[i, 1L] <- format(length(item))
      cl <- oldClass(item)
      cells[i, 2L] <- if (length(cl)) cl[1L] else "-none-"
      cells[i, 3L] <- mode(item)
    }
    return(..summaryClass(cells))
  }
  ..summaryClass(c(Length = length(object), Class = class(object), Mode = mode(object)))
}

summary.factor <- function(object, maxsum = 100L, ...) {
  gone <- is.na(object)
  lev <- levels(object)
  if (any(gone)) maxsum <- maxsum - 1L
  counts <- tabulate(as.integer(object), length(lev))
  if (length(lev)) names(counts) <- lev
  if (length(lev) > maxsum) {
    keep <- order(counts, decreasing = TRUE)
    rest <- maxsum:length(lev)
    counts <- c(counts[keep[-rest]], "(Other)" = sum(counts[keep[rest]]))
  }
  if (any(gone)) c(counts, "NAs" = sum(gone)) else counts
}

# the figures of a summary of dates or times, as text, the count of missing
# values last
..formatSummaryTime <- function(x, digits, ...) {
  v <- as.vector(unclass(x))
  names(v) <- names(x)
  shown <- if (inherits(x, "Date")) format(.Date(v), ...)
           else if (inherits(x, "POSIXct")) format(.POSIXct(v, attr(x, "tzone")), ...)
  names(shown) <- names(v)
  if (!is.null(attr(x, "NAs"))) shown <- c(shown, "NAs" = as.character(attr(x, "NAs")))
  shown
}

format.summaryDefault <- function(x, digits = max(3L, getOption("digits") - 3L), ...) {
  if (inherits(x, "Date") || inherits(x, "POSIXct")) return(..formatSummaryTime(x, digits, ...))
  if (inherits(x, "difftime")) {
    v <- as.vector(unclass(x))
    names(v) <- names(x)
    if (!is.null(attr(x, "NAs"))) v <- c(v, NAs = attr(x, "NAs"))
    x <- v
  }
  if (is.character(x)) {
    out <- format.default(as.vector(x), ...)
    names(out) <- names(x)
    return(out)
  }
  v <- unclass(x)
  attributes(v) <- NULL
  nm <- names(x)
  at <- match("NAs", nm, 0L)
  zap <- function(u) {
    ok <- is.finite(u)
    if (any(ok)) u[ok] <- zapsmall(u[ok], digits + 4L)
    u
  }
  out <- if (at) c(format(zap(v[-at]), digits = digits, ...), as.character(v[at]))
         else format(zap(v), digits = digits, ...)
  names(out) <- if (at) c(nm[-at], nm[at]) else nm
  out
}

print.summaryDefault <- function(x, digits = max(3L, getOption("digits") - 3L),
                                 quote = FALSE, ...) {
  if (inherits(x, "Date") || inherits(x, "POSIXct")) {
    # a summary of dates or times prints as quoted text unless asked not to
    print(..formatSummaryTime(x, digits), quote = if (missing(quote)) TRUE else quote, ...)
    return(invisible(x))
  }
  if (inherits(x, "difftime")) {
    cat("Time differences in ", attr(x, "units"), "\n", sep = "")
    v <- as.vector(unclass(x))
    names(v) <- names(x)
    if (!is.null(attr(x, "NAs"))) v <- c(v, NAs = attr(x, "NAs"))
    class(v) <- c("summaryDefault", "table")
    return(print.summaryDefault(v, digits = digits, quote = quote, ...))
  }
  if (length(dim(x)) == 2L) {
    print.table(x, quote = quote, ...)
  } else if (is.character(x)) {
    shown <- as.vector(x)
    names(shown) <- names(x)
    print.table(shown, quote = quote, right = TRUE, ...)
  } else {
    shown <- format(x, digits = digits)
    print.table(shown, digits = digits, quote = quote, right = TRUE, ...)
  }
  invisible(x)
}

summary.matrix <- function(object, ...) summary.data.frame(as.data.frame(object), ...)

summary.data.frame <- function(object, maxsum = 7L, digits, ...) {
  parts <- lapply(as.list(object), summary, maxsum = maxsum, digits = 12L, ...)
  nv <- length(parts)
  nm <- names(object)
  rows <- 0L
  cells <- vector("list", nv)
  labw <- integer(nv)
  fd <- if (missing(digits)) max(3L, getOption("digits") - 3L) else digits
  for (i in seq_len(nv)) {
    p <- parts[[i]]
    lab <- format(names(p))
    cells[[i]] <- paste0(lab, ":", format(p, digits = fd), "  ")
    labw[i] <- nchar(lab[1L], type = "w")
    rows <- max(rows, length(cells[[i]]))
  }
  out <- matrix(NA_character_, rows, nv)
  for (i in seq_len(nv)) out[seq_along(cells[[i]]), i] <- cells[[i]]
  head <- character(nv)
  for (i in seq_len(nv))
    head[i] <- paste0(strrep(" ", max(floor(labw[i] - nchar(nm[i]) / 2), 0L)), nm[i])
  if (nv > 0L) dimnames(out) <- list(rep.int("", rows), head)
  class(out) <- "table"
  out
}

format.pval <- function(pv, digits = max(1L, getOption("digits") - 2L),
                        eps = .Machine$double.eps, na.form = "NA", ...)
  ..format_pval(pv, digits, eps, na.form)
