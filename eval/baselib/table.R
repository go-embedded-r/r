# Base library, part 11: contingency tables. Written from the help pages of
# table, as.data.frame.table, expand.grid and interaction in R 4.6.2 and
# checked against the oracle (docs/m3-plan.md). BSD licence.

table <- function(..., exclude = if (useNA == "no") c(NA, NaN), useNA = c("no", "ifany", "always"),
                  dnn = NULL, deparse.level = 1) {
  useNA <- match.arg(useNA)
  vars <- list(...)
  if (length(vars) == 0L) stop("nothing to tabulate")
  labs <- names(vars)
  if (length(vars) == 1L && is.list(vars[[1L]])) {
    vars <- vars[[1L]]
    labs <- names(vars)
  } else if (deparse.level > 0L) {
    ex <- as.list(substitute(list(...)))[-1L]
    if (is.null(labs)) labs <- character(length(vars))
    for (k in seq_along(vars))
      if (!nzchar(labs[k]) && is.symbol(ex[[k]])) labs[k] <- as.character(ex[[k]])
  }
  if (is.null(labs)) labs <- character(length(vars))
  if (!is.null(dnn)) labs <- dnn
  if (length(labs) != length(vars)) stop("'dnn' must be a character vector with one name per factor")
  size <- 1
  lev <- vector("list", length(vars))
  cod <- vector("list", length(vars))
  for (k in seq_along(vars)) {
    v <- vars[[k]]
    if (is.factor(v)) {
      ll <- levels(v)
      v <- factor(v, levels = ll[!(ll %in% exclude)], exclude = NULL)
    } else v <- factor(v, exclude = exclude)
    if (useNA == "always") v <- addNA(v) else if (useNA == "ifany") v <- addNA(v, ifany = TRUE)
    lev[[k]] <- levels(v)
    cod[[k]] <- as.integer(v)
    size <- size * length(lev[[k]])
  }
  if (size > 2^27) stop("attempt to make a table that is too large")
  if (length(unique(lengths(cod))) > 1L) stop("all arguments must have the same length")
  cell <- 1L
  stride <- 1L
  for (k in seq_along(cod)) {
    cell <- cell + stride * (cod[[k]] - 1L)
    stride <- stride * length(lev[[k]])
  }
  lev <- lapply(lev, function(l) if (length(l)) l)
  names(lev) <- labs
  structure(array(tabulate(cell, size), vapply(lev, length, 1L), lev), class = "table")
}

cut <- function(x, ...) UseMethod("cut")

print.table <- function(x, digits = getOption("digits"), quote = FALSE, na.print = "",
                        zero.print = "0", right = is.numeric(x) || is.complex(x),
                        justify = "none", ...) {
  d <- dim(x)
  if (any(d == 0L)) {
    cat("< table of extent", paste(d, collapse = " x "), ">\n")
  } else print(unclass(x), quote = quote)
  invisible(x)
}

expand.grid <- function(..., KEEP.OUT.ATTRS = TRUE, stringsAsFactors = TRUE) {
  vars <- list(...)
  if (length(vars) == 1L && is.list(vars[[1L]])) vars <- vars[[1L]]
  k <- length(vars)
  if (k == 0L) return(data.frame())
  fallback <- paste0("Var", seq_len(k))
  nm <- names(vars)
  if (is.null(nm)) nm <- fallback else nm[!nzchar(nm)] <- fallback[!nzchar(nm)]
  d <- lengths(vars)
  total <- prod(d)
  if (total > 2^27) stop("the grid is too large")
  cols <- vector("list", k)
  inner <- 1
  for (i in seq_len(k)) {
    v <- vars[[i]]
    n <- length(v)
    pos <- rep.int(rep.int(seq_len(n), rep.int(inner, n)), total / (inner * n))
    col <- v[pos]
    if (stringsAsFactors && is.character(v)) col <- factor(col, levels = unique(v))
    cols[[i]] <- col
    inner <- inner * n
  }
  names(cols) <- nm
  if (KEEP.OUT.ATTRS) {
    dn <- vector("list", k)
    for (i in seq_len(k)) {
      v <- vars[[i]]
      dn[[i]] <- paste0(nm[i], "=", if (is.numeric(v)) format(v) else v)
    }
    names(dn) <- nm
    names(d) <- nm
    attr(cols, "out.attrs") <- list(dim = d, dimnames = dn)
  }
  attr(cols, "class") <- "data.frame"
  attr(cols, "row.names") <- .set_row_names(as.integer(total))
  cols
}

as.data.frame.table <- function(x, row.names = NULL, ..., responseName = "Freq",
                                stringsAsFactors = TRUE, sep = "", base = list(LETTERS)) {
  dn <- dimnames(x)
  nms <- names(dn)
  if (is.null(nms)) nms <- character(length(dn))
  nms[!nzchar(nms)] <- paste0("Var", seq_along(dn))[!nzchar(nms)]
  names(dn) <- nms
  grid <- expand.grid(dn, KEEP.OUT.ATTRS = FALSE, stringsAsFactors = stringsAsFactors)
  counts <- list(as.vector(unclass(x)))
  names(counts) <- responseName
  do.call(data.frame, c(list(grid), counts, list(check.names = FALSE)))
}

interaction <- function(..., drop = FALSE, sep = ".", lex.order = FALSE) {
  vars <- list(...)
  if (length(vars) == 1L && is.list(vars[[1L]])) vars <- vars[[1L]]
  vars <- lapply(vars, as.factor)
  lv <- lapply(vars, levels)
  nl <- lengths(lv)
  code <- as.integer(vars[[1L]])
  labs <- lv[[1L]]
  for (k in seq_along(vars)[-1L]) {
    if (lex.order) {
      code <- (code - 1L) * nl[k] + as.integer(vars[[k]])
      labs <- paste(rep(labs, each = nl[k]), rep(lv[[k]], times = length(labs)), sep = sep)
    } else {
      code <- code + length(labs) * (as.integer(vars[[k]]) - 1L)
      labs <- paste(rep(labs, times = nl[k]), rep(lv[[k]], each = length(labs)), sep = sep)
    }
  }
  f <- structure(code, levels = labs, class = "factor")
  if (drop) droplevels(f) else f
}
