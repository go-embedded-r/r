# The contrast matrices and model.matrix (slice 20), designed from the help
# pages of contrast.matrix, contrasts and model.matrix and from the oracle's
# answers, on the terms objects of slice 17.

..contr_levels <- function(n) {
  if (is.numeric(n) && length(n) == 1L) as.character(seq_len(n)) else as.character(n)
}

..contr_few <- function(n) is.numeric(n) && length(n) == 1L && n <= 1L

..contr_identity <- function(lev) {
  m <- diag(length(lev))
  dimnames(m) <- list(lev, lev)
  m
}

contr.treatment <- function(n, base = 1, contrasts = TRUE, sparse = FALSE) {
  if (..contr_few(n)) stop("not enough degrees of freedom to define contrasts")
  lev <- ..contr_levels(n)
  k <- length(lev)
  m <- ..contr_identity(lev)
  if (!contrasts) return(m)
  if (k < 2L) stop(gettextf("contrasts not defined for %d degrees of freedom", k - 1L), domain = NA)
  if (base < 1L || base > k) stop("baseline group number out of range")
  m[, -base, drop = FALSE]
}

contr.SAS <- function(n, contrasts = TRUE, sparse = FALSE) {
  contr.treatment(n, base = if (is.numeric(n) && length(n) == 1L) n else length(n), contrasts = contrasts, sparse = sparse)
}

contr.sum <- function(n, contrasts = TRUE, sparse = FALSE) {
  if (..contr_few(n)) stop("not enough degrees of freedom to define contrasts")
  lev <- ..contr_levels(n)
  k <- length(lev)
  if (!contrasts) return(..contr_identity(lev))
  if (k < 2L) stop(gettextf("contrasts not defined for %d degrees of freedom", k - 1L), domain = NA)
  m <- rbind(diag(k - 1L), -1)
  dimnames(m) <- list(lev, NULL)
  m
}

contr.helmert <- function(n, contrasts = TRUE, sparse = FALSE) {
  if (..contr_few(n)) stop("not enough degrees of freedom to define contrasts")
  lev <- ..contr_levels(n)
  k <- length(lev)
  if (!contrasts) return(..contr_identity(lev))
  if (k < 2L) stop(gettextf("contrasts not defined for %d degrees of freedom", k - 1L), domain = NA)
  m <- matrix(-1, k, k - 1L)
  below <- row(m) - col(m)
  m[below > 1L] <- 0
  m[below == 1L] <- seq_len(k - 1L)
  dimnames(m) <- list(lev, NULL)
  m
}

contr.poly <- function(n, scores = 1:n, contrasts = TRUE, sparse = FALSE) {
  if (length(n) > 1L) n <- length(n)
  if (n < 2L) stop(gettextf("contrasts not defined for %d degrees of freedom", n - 1L), domain = NA)
  if (length(scores) != n) stop("'scores' argument is of the wrong length")
  if (!is.numeric(scores) || anyDuplicated(scores)) stop("'scores' must all be different numbers")
  decomp <- qr(outer(scores, 0L:(n - 1L), "^"))
  z <- qr.Q(decomp)
  z <- sweep(z, 2L, sign(diag(qr.R(decomp))), "*")
  z <- sweep(z, 2L, sqrt(colSums(z^2)), "/")
  nm <- c(".L", ".Q", ".C", paste0("^", 4L:n))[seq_len(n - 1L)]
  if (contrasts) {
    z <- z[, -1L, drop = FALSE]
    dimnames(z) <- list(NULL, nm)
  } else {
    z[, 1L] <- 1
    dimnames(z) <- list(NULL, paste0("^", 0L:(n - 1L)))
  }
  z
}

contrasts <- function(x, contrasts = TRUE, sparse = FALSE) {
  if (is.logical(x)) x <- factor(x, levels = c(FALSE, TRUE))
  if (!is.factor(x)) stop("contrasts apply only to factors")
  lev <- levels(x)
  if (!contrasts) return(..contr_identity(lev))
  ctr <- attr(x, "contrasts")
  if (is.null(ctr)) {
    ctr <- getOption("contrasts")[[if (is.ordered(x)) 2L else 1L]]
  }
  if (is.character(ctr)) ctr <- get(ctr, mode = "function", envir = parent.frame())(lev)
  ctr
}

`contrasts<-` <- function(x, how.many = NULL, value) {
  if (is.logical(x)) x <- factor(x, levels = c(FALSE, TRUE))
  if (!is.factor(x)) stop("contrasts apply only to factors")
  nl <- nlevels(x)
  if (nl < 2L) stop("contrasts can be applied only to factors with 2 or more levels")
  if (is.function(value)) value <- value(nl)
  if (is.numeric(value)) {
    value <- as.matrix(value)
    if (nrow(value) != nl) stop("wrong number of contrast matrix rows")
    want <- if (is.null(how.many)) nl - 1L else how.many
    have <- ncol(value)
    rownames(value) <- levels(x)
    if (have < want) {
      fit <- qr(cbind(1, value))
      if (fit$rank != have + 1L) stop("singular contrast matrix")
      full <- qr.qy(fit, diag(nl))[, 2L:nl, drop = FALSE]
      full[, seq_len(have)] <- value
      dimnames(full) <- list(levels(x), NULL)
      cn <- colnames(value)
      if (!is.null(cn)) colnames(full) <- c(cn, rep.int("", want - have))
      value <- full
    } else value <- value[, seq_len(want), drop = FALSE]
  } else if (!is.character(value) && !is.null(value)) {
    stop("numeric contrasts or contrast name expected")
  }
  attr(x, "contrasts") <- value
  x
}

model.matrix <- function(object, ...) UseMethod("model.matrix")

model.matrix.default <- function(object, data = environment(object), contrasts.arg = NULL, xlev = NULL, ...) {
  tt <- terms(object)
  if (is.null(attr(data, "terms"))) data <- model.frame(object, data, xlev = xlev)
  fac <- attr(tt, "factors")
  has_int <- attr(tt, "intercept") == 1L
  n <- nrow(data)
  vnames <- if (length(fac)) rownames(fac) else character()
  if (!is.null(contrasts.arg)) {
    if (!is.list(contrasts.arg)) warning("'contrasts.arg' is not a list")
    for (nm in names(contrasts.arg)) {
      if (!nm %in% names(data)) {
        warning(gettextf("variable '%s' is absent, its contrast will be ignored", nm), domain = NA)
      } else {
        col <- data[[nm]]
        if (is.character(col)) col <- factor(col)
        contrasts(col) <- contrasts.arg[[nm]]
        data[[nm]] <- col
      }
    }
  }
  isfac <- logical(length(vnames))
  for (i in seq_along(vnames)) {
    v <- data[[vnames[i]]]
    if (is.character(v)) {
      v <- factor(v)
      data[[vnames[i]]] <- v
    }
    isfac[i] <- is.factor(v) || is.logical(v)
  }
  contr.funs <- as.character(getOption("contrasts"))
  isOF <- vapply(data, is.ordered, NA)
  for (nn in vnames[isfac]) {
    if (is.null(attr(data[[nn]], "contrasts"))) contrasts(data[[nn]]) <- contr.funs[1 + isOF[nn]]
  }
  # without an intercept the first factor met is coded by indicators
  code <- fac
  if (length(code) && !has_int) {
    done <- FALSE
    for (j in seq_len(ncol(code))) {
      for (i in seq_len(nrow(code))) {
        if (isfac[i] && code[i, j] > 0L) {
          code[i, j] <- 2L
          done <- TRUE
          break
        }
      }
      if (done) break
    }
  }
  blocks <- list()
  cmat <- list()
  for (i in which(isfac)) {
    v <- data[[vnames[i]]]
    if (is.logical(v)) v <- factor(v, levels = c(FALSE, TRUE))
    if (any(code[i, ] > 0L)) cmat[[vnames[i]]] <- contrasts(v)
  }
  cols <- list()
  names_out <- character()
  assign <- integer()
  if (has_int) {
    cols[[1L]] <- rep(1, n)
    names_out <- "(Intercept)"
    assign <- 0L
  }
  for (j in seq_len(if (length(fac)) ncol(fac) else 0L)) {
    cur <- matrix(1, n, 1L)
    curnames <- ""
    for (i in which(code[, j] > 0L)) {
      v <- data[[vnames[i]]]
      if (isfac[i]) {
        if (is.logical(v)) v <- factor(v, levels = c(FALSE, TRUE))
        if (code[i, j] == 2L) {
          blk <- outer(as.integer(v), seq_len(nlevels(v)), "==") + 0
          bn <- paste0(vnames[i], levels(v))
        } else {
          cm <- cmat[[vnames[i]]]
          blk <- cm[as.integer(v), , drop = FALSE]
          bn <- paste0(vnames[i], if (is.null(colnames(cm))) seq_len(ncol(cm)) else colnames(cm))
          dimnames(blk) <- NULL
        }
      } else {
        blk <- as.matrix(v)
        bn <- if (ncol(blk) == 1L && is.null(colnames(blk))) vnames[i] else paste0(vnames[i], if (is.null(colnames(blk))) seq_len(ncol(blk)) else colnames(blk))
        dimnames(blk) <- NULL
      }
      nxt <- matrix(0, n, ncol(cur) * ncol(blk))
      nxtnames <- character(ncol(nxt))
      k <- 0L
      for (b in seq_len(ncol(blk))) {
        for (a in seq_len(ncol(cur))) {
          k <- k + 1L
          nxt[, k] <- cur[, a] * blk[, b]
          nxtnames[k] <- if (nzchar(curnames[a])) paste(curnames[a], bn[b], sep = ":") else bn[b]
        }
      }
      cur <- nxt
      curnames <- nxtnames
    }
    for (a in seq_len(ncol(cur))) cols[[length(cols) + 1L]] <- cur[, a]
    names_out <- c(names_out, curnames)
    assign <- c(assign, rep.int(j, ncol(cur)))
  }
  ans <- matrix(0, n, length(cols))
  for (k in seq_along(cols)) ans[, k] <- cols[[k]]
  dimnames(ans) <- list(row.names(data), names_out)
  attr(ans, "assign") <- assign
  if (length(cmat)) {
    ctr <- list()
    for (nm in names(cmat)) {
      v <- data[[nm]]
      how <- if (is.factor(v)) attr(v, "contrasts") else NULL
      if (is.null(how)) how <- getOption("contrasts")[[if (is.ordered(v)) 2L else 1L]]
      ctr[[nm]] <- how
    }
    attr(ans, "contrasts") <- ctr
  }
  ans
}
