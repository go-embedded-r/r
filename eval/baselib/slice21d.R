# Slice 21, part 3: ftable, ecdf, poly and read.fwf.  Written from the Rd
# pages and from oracle runs; the layout of format.ftable follows what the
# oracle printed for each method, justify and quote setting.

is.unsorted <- function(x, na.rm = FALSE, strictly = FALSE) {
  if (length(x) <= 1L) return(FALSE)
  if (!na.rm && anyNA(x)) return(NA)
  if (na.rm) x <- x[!is.na(x)]
  n <- length(x)
  if (n <= 1L) return(FALSE)
  a <- x[-n]
  b <- x[-1L]
  if (strictly) any(b <= a) else any(b < a)
}

## ---- ftable ----------------------------------------------------------

ftable <- function(x, ...) UseMethod("ftable")

ftable.default <- function(..., exclude = c(NA, NaN), row.vars = NULL, col.vars = NULL) {
  dots <- list(...)
  if (!length(dots)) stop("nothing to tabulate")
  first <- dots[[1L]]
  multi <- is.array(first) && length(dim(first)) > 1L
  x <- if (inherits(first, "ftable")) {
    as.table(first)
  } else if (is.list(first)) {
    table(first, exclude = exclude)
  } else if (multi && is.table(first)) {
    first
  } else if (multi) {
    as.table(first)
  } else {
    table(..., exclude = exclude)
  }
  dn <- dimnames(x)
  nx <- length(dim(x))
  if (!is.null(row.vars)) {
    if (is.character(row.vars)) {
      row.vars <- match(row.vars, names(dn))
      if (anyNA(row.vars)) stop("incorrect specification for 'row.vars'")
    } else if (!is.numeric(row.vars) || anyNA(row.vars) || any(row.vars < 1 | row.vars > nx))
      stop("incorrect specification for 'row.vars'")
    else row.vars <- as.integer(row.vars)
  }
  if (!is.null(col.vars)) {
    if (is.character(col.vars)) {
      col.vars <- match(col.vars, names(dn))
      if (anyNA(col.vars)) stop("incorrect specification for 'col.vars'")
    } else if (!is.numeric(col.vars) || anyNA(col.vars) || any(col.vars < 1 | col.vars > nx))
      stop("incorrect specification for 'col.vars'")
    else col.vars <- as.integer(col.vars)
  }
  keep <- c(row.vars, col.vars)
  marginal <- !is.null(row.vars) && !is.null(col.vars) && length(keep) < nx && !anyDuplicated(keep)
  if (marginal) {
    x <- marginSums(x, keep)
    row.vars <- match(row.vars, keep)
    col.vars <- match(col.vars, keep)
    dn <- dimnames(x)
    nx <- length(dim(x))
  }
  all <- seq_len(nx)
  if (is.null(row.vars) && is.null(col.vars)) {
    col.vars <- nx
    row.vars <- all[-nx]
  } else if (is.null(row.vars)) row.vars <- setdiff(all, col.vars)
  else if (is.null(col.vars)) col.vars <- setdiff(all, row.vars)
  y <- aperm(x, c(rev(row.vars), rev(col.vars)))
  d <- dim(x)
  dd <- c(prod(d[row.vars]), prod(d[col.vars]))
  if (marginal) {
    y <- as.vector(y)
    dim(y) <- dd
    attr(y, "row.vars") <- dn[row.vars]
    attr(y, "col.vars") <- dn[col.vars]
    class(y) <- "ftable"
    return(y)
  }
  structure(as.vector(y), dim = dd, class = "ftable", row.vars = dn[row.vars], col.vars = dn[col.vars])
}

ftable.formula <- function(formula, data = NULL, subset, na.action, ...) {
  if (length(formula) != 3L) stop("'formula' must have both left and right hand sides")
  dots <- list(...)
  exclude <- if (is.null(dots$exclude)) c(NA, NaN) else dots$exclude
  lv <- all.vars(formula[[2L]])
  rv <- all.vars(formula[[3L]])
  if (inherits(data, "table")) {
    dn <- names(dimnames(data))
    if ("." %in% lv) lv <- c(setdiff(lv, "."), setdiff(dn, c(rv, lv)))
    if ("." %in% rv) rv <- c(setdiff(rv, "."), setdiff(dn, c(rv, lv)))
    if (!all(rv %in% dn)) stop("incorrect variable names in rhs of formula")
    if (!all(lv %in% dn)) stop("incorrect variable names in lhs of formula")
    tab <- marginSums(data, match(c(rv, lv), dn))
    return(ftable.default(tab, row.vars = rv, col.vars = lv))
  }
  if ("." %in% c(lv, rv)) stop("cannot use dots in formula with given data")
  f2 <- eval(call("~", call("+", formula[[3L]], formula[[2L]])), environment(formula))
  environment(f2) <- environment(formula)
  mc <- match.call(expand.dots = FALSE)
  ml <- as.list(mc)
  extra <- ml[match(c("data", "subset", "na.action"), names(ml), 0L)]
  mf <- as.call(c(list(quote(model.frame), formula = f2), extra))
  mf <- eval(mf, parent.frame())
  tab <- table(mf, exclude = exclude)
  ftable.default(tab, row.vars = seq_along(rv), col.vars = length(rv) + seq_along(lv))
}

format.ftable <- function(x, quote = TRUE, digits = getOption("digits"),
                          method = c("non.compact", "row.compact", "col.compact", "compact"),
                          lsep = " | ", justify = c("left", "right"), ...) {
  if (!inherits(x, "ftable")) stop("'x' must be an \"ftable\" object")
  method <- match.arg(method)
  justify <- rep_len(justify, 2L)
  rv <- attr(x, "row.vars")
  cv <- attr(x, "col.vars")
  nr <- length(rv)
  nc <- length(cv)
  rnames <- names(rv)
  if (is.null(rnames)) rnames <- character(nr)
  cnames <- names(cv)
  if (is.null(cnames)) cnames <- character(nc)
  q <- function(s) if (quote && length(s)) paste0("\"", s, "\"") else s
  # the label cells of a block of variables: one column per variable, a label
  # shown where its run begins
  labs <- function(l) {
    n <- length(l)
    lens <- lengths(l)
    N <- prod(lens)
    m <- matrix("", N, n)
    for (i in seq_len(n)) {
      stride <- prod(lens[seq_len(n)[-seq_len(i)]])
      at <- seq.int(0L, N - 1L)
      show <- at %% stride == 0L
      m[show, i] <- q(as.character(l[[i]])[(at[show] %/% stride) %% lens[i] + 1L])
    }
    m
  }
  nR <- prod(lengths(rv))
  nC <- prod(lengths(cv))
  RL <- labs(rv)
  CL <- t(labs(cv))
  dat <- matrix(format(as.vector(unclass(x)), digits = digits, ...), nR, nC)
  hr <- max(nc, 1L)
  merged <- method %in% c("col.compact", "compact")
  mcol <- max(nr, 1L)
  nlab <- if (merged) mcol else nr + 1L
  extra <- method %in% c("non.compact", "col.compact")
  nrows <- if (extra) nc + 1L + nR else hr + nR
  M <- matrix("", nrows, nlab + nC)
  top <- if (extra) nc + 1L else hr
  drows <- top + seq_len(nR)
  dcols <- nlab + seq_len(nC)
  if (nr > 0L) M[drows, seq_len(nr)] <- RL
  M[drows, dcols] <- dat
  if (nc > 0L) M[seq_len(nc), dcols] <- CL
  ccol <- nlab
  if (method == "compact") {
    if (nc > 1L) M[seq_len(nc - 1L), ccol] <- q(cnames[-nc])
    if (nr > 1L) M[hr, seq_len(nr - 1L)] <- q(rnames[-nr])
    M[hr, ccol] <- paste0(if (nr > 0L) rnames[nr] else "", lsep, if (nc > 0L) cnames[nc] else "")
  } else {
    if (nc > 0L) M[seq_len(nc), ccol] <- q(cnames)
    if (nr > 0L) M[top, seq_len(nr)] <- q(rnames)
  }
  for (j in seq_len(ncol(M))) M[, j] <- format(M[, j], justify = justify[if (j <= nlab) 1L else 2L])
  M
}

write.ftable <- function(x, file = "", quote = TRUE, append = FALSE,
                         digits = getOption("digits"), sep = " ", ...) {
  r <- format.ftable(x, quote = quote, digits = digits, ...)
  cat(t(r), file = file, append = append, sep = c(rep.int(sep, ncol(r) - 1L), "\n"))
  invisible(x)
}

print.ftable <- function(x, digits = getOption("digits"), ...) {
  write.ftable(x, quote = FALSE, digits = digits, ...)
  invisible(x)
}

as.table.ftable <- function(x, ...) {
  rv <- attr(x, "row.vars")
  cv <- attr(x, "col.vars")
  nr <- length(rv)
  nc <- length(cv)
  a <- array(as.vector(unclass(x)), c(rev(lengths(rv)), rev(lengths(cv))))
  a <- aperm(a, c(rev(seq_len(nr)), nr + rev(seq_len(nc))))
  dimnames(a) <- c(rv, cv)
  class(a) <- "table"
  a
}

as.data.frame.ftable <- function(x, row.names = NULL, optional = FALSE, ...)
  as.data.frame(as.table(x), row.names = row.names, optional = optional)

as.matrix.ftable <- function(x, sep = "_", ...) {
  rv <- attr(x, "row.vars")
  cv <- attr(x, "col.vars")
  joint <- function(l) {
    lens <- lengths(l)
    N <- prod(lens)
    out <- character(N)
    for (i in seq_along(l)) {
      stride <- prod(lens[seq_along(l)[-seq_len(i)]])
      lab <- as.character(l[[i]])[(seq.int(0L, N - 1L) %/% stride) %% lens[i] + 1L]
      out <- if (i == 1L) lab else paste(out, lab, sep = sep)
    }
    out
  }
  nm <- function(l) {
    n <- names(l)
    if (is.null(n)) NULL else paste(n, collapse = sep)
  }
  m <- matrix(as.vector(unclass(x)), nrow(x), ncol(x))
  rn <- nm(rv)
  cn <- nm(cv)
  dn <- list(joint(rv), joint(cv))
  if (!is.null(rn) || !is.null(cn)) names(dn) <- c(if (is.null(rn)) "" else rn, if (is.null(cn)) "" else cn)
  dimnames(m) <- dn
  m
}

## ---- approxfun, ecdf -------------------------------------------------

approxfun <- function(x, y = NULL, method = "linear", yleft, yright, rule = 1, f = 0,
                      ties = mean, na.rm = TRUE) {
  if (is.null(y)) {
    y <- x
    x <- seq_along(y)
  }
  method <- pmatch(method, c("linear", "constant"))
  if (is.na(method)) stop("invalid interpolation method")
  x <- as.numeric(x)
  y <- as.numeric(y)
  if (length(x) != length(y)) stop("'x' and 'y' lengths differ")
  ok <- !(is.na(x) | is.na(y))
  if (na.rm) {
    x <- x[ok]
    y <- y[ok]
  } else if (any(!ok)) stop("approx(): attempted to interpolate NA values")
  if (length(x) == 0L) stop("zero non-NA points")
  if (!identical(ties, "ordered")) {
    if (is.unsorted(x)) {
      o <- order(x)
      x <- x[o]
      y <- y[o]
    }
    if (anyDuplicated(x)) {
      if (is.character(ties) && length(ties) == 1L) ties <- match.fun(ties[[1L]])
      else if (is.list(ties) && length(ties) == 2L && is.function(ties[[2L]])) ties <- ties[[2L]]
      y <- as.vector(tapply(y, match(x, x), ties))
      x <- unique(x)
    }
  } else if (is.unsorted(x)) stop("x must be sorted when ties = \"ordered\"")
  if (length(x) <= 1L && method == 1L) stop("need at least two non-NA values to interpolate")
  rule <- rep_len(rule, 2L)
  # outside the data a rule of 1 gives NA, a rule of 2 the nearest end value
  ends <- c(y[1L], y[length(y)])
  edge <- function(i) if (rule[i] == 1) NA else ends[i]
  if (missing(yleft)) yleft <- edge(1L)
  if (missing(yright)) yright <- edge(2L)
  force(f)
  env <- new.env(parent = environment(approxfun))
  assign("x", x, envir = env)
  assign("y", y, envir = env)
  assign("method", method, envir = env)
  assign("yleft", yleft, envir = env)
  assign("yright", yright, envir = env)
  assign("f", f, envir = env)
  assign("na.rm", na.rm, envir = env)
  rval <- function(v) ..approxeval(x, y, v, method, yleft, yright, f)
  environment(rval) <- env
  rval
}

ecdf <- function(x) {
  if ((is.numeric(x) || is.logical(x)) && all(is.na(x))) stop("'x' must have 1 or more non-missing values")
  if (!is.numeric(x) && !is.logical(x)) {
    # R fails inside approxfun on the unique values, and the message names that call
    shown <- as.call(list(as.name("approxfun"), quote(vals), quote(cumsum(tabulate(match(x, vals)))/n),
                          method = "constant", yleft = 0, yright = 1, f = 0, ties = "ordered"))
    stop(simpleError("zero non-NA points", shown))
  }
  k <- ..ecdfknots(x)
  fn <- approxfun(k$x, k$y, method = "constant", yleft = 0, yright = 1, f = 0, ties = "ordered")
  class(fn) <- c("ecdf", "stepfun", class(fn))
  assign("nobs", k$n, envir = environment(fn))
  attr(fn, "call") <- sys.call()
  fn
}

knots <- function(Fn, ...) UseMethod("knots")
knots.stepfun <- function(Fn, ...) environment(Fn)$x

print.ecdf <- function(x, digits = getOption("digits") - 2L, ...) {
  cat("Empirical CDF \nCall: ")
  print(attr(x, "call"), ...)
  cat(..ecdfline(environment(x)$x, digits))
  invisible(x)
}

summary.ecdf <- function(object, ...) {
  n <- length(environment(object)$x)
  header <- paste("Empirical CDF:\t ", n, "unique values with summary\n")
  structure(summary(knots(object), ...), header = header, class = "summary.ecdf")
}

print.summary.ecdf <- function(x, ...) {
  cat(attr(x, "header"))
  y <- x
  attr(y, "header") <- NULL
  class(y) <- c("summaryDefault", "table")
  print(y, ...)
  invisible(x)
}

quantile.ecdf <- function(x, ...) {
  e <- environment(x)
  quantile(rep.int(e$x, diff(c(0, round(e$nobs * e$y)))), ...)
}

## ---- poly -------------------------------------------------------------

predict <- function(object, ...) UseMethod("predict")

# the three-term recurrence of the monic polynomials with the centres alpha and
# the squared norms norm2 (norm2[1] = 1, norm2[2] = n), scaled to unit norm
..polyrec <- function(x, degree, alpha, norm2) {
  Z <- matrix(1, length(x), degree + 1L)
  if (degree > 0L) Z[, 2L] <- x - alpha[1L]
  if (degree > 1L)
    for (i in 2:degree) Z[, i + 1L] <- (x - alpha[i]) * Z[, i] - (norm2[i + 1L]/norm2[i]) * Z[, i - 1L]
  Z <- Z/rep(sqrt(norm2[-1L]), each = length(x))
  Z[, -1L, drop = FALSE]
}

poly <- function(x, ..., degree = 1, coefs = NULL, raw = FALSE, simple = FALSE) {
  dots <- list(...)
  if (length(dots)) {
    if (length(dots) == 1L && length(dots[[1L]]) == 1L) degree <- dots[[1L]]
    else return(polym(x, ..., degree = degree, coefs = coefs, raw = raw))
  }
  if (is.matrix(x)) {
    if (ncol(x) > 1L) {
      m <- lapply(seq_len(ncol(x)), function(j) x[, j])
      return(do.call(polym, c(m, list(degree = degree, coefs = coefs, raw = raw))))
    }
    x <- x[, 1L]
  }
  if (degree < 1) stop("'degree' must be at least 1")
  if (!raw && anyNA(x)) stop("missing values are not allowed in 'poly'")
  nms <- names(x)
  rnames <- nms
  if (raw) {
    Z <- outer(x, 1L:degree, `^`)
    colnames(Z) <- 1L:degree
    if (!is.null(rnames)) rownames(Z) <- rnames
    if (simple) return(Z)
    return(structure(Z, degree = 1L:degree, class = c("poly", "matrix")))
  }
  if (is.null(coefs)) {
    if (degree >= length(unique(x))) stop("'degree' must be less than number of unique points")
    xbar <- mean(x)
    xc <- x - xbar
    norm2 <- numeric(degree + 2L)
    alpha <- numeric(degree)
    norm2[1L] <- 1
    norm2[2L] <- length(x)
    alpha[1L] <- xbar
    P0 <- rep(1, length(x))
    P1 <- xc
    norm2[3L] <- sum(P1^2)
    if (degree > 1L) {
      alpha[2L] <- xbar + sum(xc * P1^2)/norm2[3L]
      Pm <- P0
      Pc <- P1
      for (i in 2:degree) {
        Pn <- (x - alpha[i]) * Pc - (norm2[i + 1L]/norm2[i]) * Pm
        norm2[i + 2L] <- sum(Pn^2)
        if (i < degree) alpha[i + 1L] <- xbar + sum(xc * Pn^2)/norm2[i + 2L]
        Pm <- Pc
        Pc <- Pn
      }
    }
    coefs <- list(alpha = alpha, norm2 = norm2)
  }
  Z <- ..polyrec(x, degree, coefs$alpha, coefs$norm2)
  colnames(Z) <- 1L:degree
  if (!is.null(rnames)) rownames(Z) <- rnames
  if (simple) return(Z)
  structure(Z, coefs = coefs, degree = 1L:degree, class = c("poly", "matrix"))
}

# products of univariate polynomials of the variables, total degree <= degree
polym <- function(..., degree = 1, coefs = NULL, raw = FALSE) {
  v <- list(...)
  nv <- length(v)
  if (nv == 0L) stop("must supply one or more vectors")
  if (nv == 1L && is.list(v[[1L]]) && length(v[[1L]]) > 1L) {
    v <- v[[1L]]
    nv <- length(v)
  }
  if (nv == 1L) return(poly(v[[1L]], degree, coefs = coefs, raw = raw))
  if (any(lengths(v) != length(v[[1L]]))) stop("arguments must have the same length")
  if (degree < 1) stop("'degree' must be at least 1")
  uni <- vector("list", nv)
  for (j in seq_len(nv))
    uni[[j]] <- poly(v[[j]], degree, coefs = if (is.null(coefs)) NULL else coefs[[j]], raw = raw, simple = FALSE)
  grid <- as.matrix(do.call(expand.grid, rep(list(0:degree), nv)))
  grid <- grid[rowSums(grid) <= degree & rowSums(grid) > 0, , drop = FALSE]
  n <- length(v[[1L]])
  Z <- matrix(1, n, nrow(grid))
  for (k in seq_len(nrow(grid)))
    for (j in seq_len(nv))
      if (grid[k, j] > 0L) Z[, k] <- Z[, k] * uni[[j]][, grid[k, j]]
  idx <- 1L + as.vector(grid %*% (degree + 1L)^(seq_len(nv) - 1L))
  cn <- apply(grid, 1L, paste, collapse = ".")
  dimnames(Z) <- list(NULL, structure(cn, names = idx))
  tot <- rowSums(grid)
  if (raw) return(structure(Z, degree = tot, class = c("poly", "matrix")))
  structure(Z, degree = tot, coefs = lapply(uni, attr, "coefs"), class = c("poly", "matrix"))
}

predict.poly <- function(object, newdata, ...) {
  if (missing(newdata)) return(object)
  if (is.null(attr(object, "coefs"))) {
    d <- attr(object, "degree")
    Z <- outer(newdata, d, `^`)
    colnames(Z) <- d
    return(Z)
  }
  d <- max(attr(object, "degree"))
  coefs <- attr(object, "coefs")
  Z <- ..polyrec(newdata, d, coefs$alpha, coefs$norm2)
  colnames(Z) <- 1L:d
  Z
}

## ---- read.fwf ---------------------------------------------------------

read.fwf <- function(file, widths, header = FALSE, sep = "\t", skip = 0L, row.names, col.names,
                     n = -1L, buffersize = 2000, fileEncoding = "", ...) {
  if (is.list(widths)) stop("'widths' as a list (records over several lines) is not supported")
  if (is.character(file)) {
    file <- if (nzchar(fileEncoding)) file(file, "rt", encoding = fileEncoding) else file(file, "rt")
    on.exit(close(file))
  } else if (!isOpen(file)) {
    open(file, "rt")
    on.exit(close(file))
  }
  if (skip > 0L) readLines(file, skip)
  headerline <- if (header) readLines(file, 1L) else character()
  lines <- readLines(file, n)
  w <- abs(widths)
  ends <- cumsum(w)
  keep <- widths > 0
  first <- (ends - w + 1)[keep]
  last <- ends[keep]
  seps <- rep_len(sep, length(first) - 1L)
  nc <- nchar(lines)
  # a field that starts beyond the end of its line reads as NA
  piece <- function(k) ifelse(first[k] <= nc, substring(lines, first[k], last[k]), "NA")
  out <- piece(1L)
  for (k in seq_along(first)[-1L]) out <- paste0(out, seps[k - 1L], piece(k))
  if (!length(lines)) out <- character()
  if (length(c(headerline, out)) == 0L) {
    stop(simpleError("no lines available in input", quote(read.table(file = FILE, header = header,
      sep = sep, row.names = row.names, col.names = col.names, quote = "", ...))))
  }
  FILE <- textConnection(c(headerline, out))
  on.exit(close(FILE), add = TRUE)
  # the call is built from symbols so an error shows it as R's read.fwf does
  args <- list(file = quote(FILE), header = quote(header), sep = quote(sep))
  if (!missing(row.names)) args$row.names <- quote(row.names)
  if (!missing(col.names)) args$col.names <- quote(col.names)
  args$quote <- ""
  args[[length(args) + 1L]] <- quote(...)
  tryCatch(eval(as.call(c(list(quote(read.table)), args))), error = function(e) {
    cl <- conditionCall(e)
    if (!is.call(cl) || !identical(cl[[1L]], quote(read.table))) stop(e)
    # R's read.fwf passes every argument on, so its errors show them all
    stop(simpleError(conditionMessage(e), quote(read.table(file = FILE, header = header, sep = sep,
      row.names = row.names, col.names = col.names, quote = "", ...))))
  })
}
