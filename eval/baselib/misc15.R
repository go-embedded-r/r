# Base library, slice 15 (1): small closures of base that were deferred.
# Written from the help pages of R 4.6.2 and experiments with the oracle.

alist <- function(...) as.list(sys.call())[-1L]

as.function <- function(x, ...) UseMethod("as.function")

enquote <- function(cl) as.call(list(quote(base::quote), cl))

.Deprecated <- function(new, package = NULL, msg,
                        old = as.character(sys.call(sys.parent()))[1L]) {
  if (missing(msg)) {
    msg <- gettextf("'%s' is deprecated.\n", old)
    if (!missing(new)) msg <- c(msg, gettextf("Use '%s' instead.\n", new))
    msg <- c(msg, if (is.null(package)) gettext("See help(\"Deprecated\")")
                  else gettextf("See help(\"Deprecated\") and help(\"%s-deprecated\").", package))
  } else msg <- as.character(msg)
  if (missing(new)) new <- NULL
  warning(warningCondition(paste(msg, collapse = ""), old = old, new = new, package = package,
                           class = "deprecatedWarning", call = sys.call(sys.parent())))
}

.Defunct <- function(new, package = NULL, msg,
                     old = as.character(sys.call(sys.parent()))[1L]) {
  if (missing(msg)) {
    msg <- gettextf("'%s' is defunct.\n", old)
    if (!missing(new)) msg <- c(msg, gettextf("Use '%s' instead.\n", new))
    msg <- c(msg, if (is.null(package)) gettext("See help(\"Defunct\")")
                  else gettextf("See help(\"Defunct\") and help(\"%s-defunct\").", package))
  } else msg <- as.character(msg)
  if (missing(new)) new <- NULL
  stop(errorCondition(paste(msg, collapse = ""), old = old, new = new, package = package,
                      class = "defunctError", call = sys.call(sys.parent())))
}

.NotYetImplemented <- function()
  stop(gettextf("'%s' is not implemented yet", as.character(sys.call(sys.parent()))[1L]),
       call. = FALSE)

.NotYetUsed <- function(arg, error = TRUE) {
  msg <- gettextf("argument '%s' is not used (yet)", arg)
  if (error) stop(msg, domain = NA, call. = FALSE) else warning(msg, domain = NA, call. = FALSE)
}

is.R <- function() .Defunct(msg = "'is.R' is defunct.\nSee help(\"Defunct\")")

as.null <- function(x, ...) UseMethod("as.null")

as.null.default <- function(x, ...) NULL

isa <- function(x, what) all(match(what, class(x), 0L) > 0L)

is.table <- function(x) inherits(x, "table")

is.ordered <- function(x) inherits(x, "ordered")

ordered <- function(x, ...) factor(x, ..., ordered = TRUE)

as.ordered <- function(x) if (is.ordered(x)) x else ordered(x)

gl <- function(n, k, length = n * k, labels = seq_len(n), ordered = FALSE)
  factor(rep_len(rep(seq_len(n), each = k), length), levels = seq_len(n),
         labels = labels, ordered = ordered)

levels.default <- function(x) attr(x, "levels")

labels <- function(object, ...) UseMethod("labels")

# labels.default, mostattributes<- are Go natives (eval/funcs15.go)

data.class <- function(x) {
  cl <- oldClass(x)
  if (length(cl)) return(cl[1L])
  nd <- length(dim(x))
  if (nd == 2L) "matrix" else if (nd) "array" else mode(x)
}

I <- function(x) {
  class(x) <- unique.default(c("AsIs", oldClass(x)))
  x
}

print.AsIs <- function(x, ...) {
  cl <- oldClass(x)
  oldClass(x) <- cl[cl != "AsIs"]
  NextMethod("print")
  invisible(x)
}

`[.AsIs` <- function(x, i, ...) I(NextMethod("["))

format.AsIs <- function(x, width = 12, ...) {
  if (is.character(x) && !is.list(x)) return(format.default(x, ...))
  s <- character(length(x))
  for (i in seq_along(x)) {
    e <- x[[i]]
    cl <- oldClass(e)
    oldClass(e) <- cl[cl != "AsIs"]
    s[i] <- toString(e, width = width, ...)
  }
  dim(s) <- dim(x)
  dimnames(s) <- dimnames(x)
  format.default(s, justify = "right")
}

toString <- function(x, ...) UseMethod("toString")

toString.default <- function(x, width = NULL, ...) {
  s <- paste(x, collapse = ", ")
  if (is.null(width) || width == 0) return(s)
  if (width < 0L) stop("'width' must be positive")
  if (nchar(s, type = "w") > width) {
    width <- max(6, width)
    s <- paste0(strtrim(s, width - 4), "....")
  }
  s
}

`is.na<-` <- function(x, value) UseMethod("is.na<-")

`is.na<-.default` <- function(x, value) {
  x[value] <- NA
  x
}

`is.na<-.factor` <- function(x, value) {
  lv <- levels(x)
  x[value] <- NA
  structure(x, levels = lv, class = oldClass(x))
}

`substring<-` <- function(text, first, last = 1000000L, value) `substr<-`(text, first, last, value)

arrayInd <- function(ind, .dim, .dimnames = NULL, useNames = FALSE) {
  m <- length(ind)
  rank <- length(.dim)
  ind <- ind - 1L
  res <- matrix(0L, m, rank)
  for (k in seq_len(rank)) {
    res[, k] <- as.integer(ind %% .dim[k] + 1L)
    ind <- ind %/% .dim[k]
  }
  if (useNames) {
    rn <- if (!is.null(.dimnames) && !is.null(.dimnames[[1L]])) .dimnames[[1L]][res[, 1L]]
    cn <- if (rank == 2L) c("row", "col") else paste0("dim", seq_len(rank))
    dimnames(res) <- list(rn, cn)
  }
  res
}

slice.index <- function(x, MARGIN) {
  d <- dim(x)
  if (is.null(d)) d <- length(x)
  res <- 0L
  mult <- 1L
  for (m in MARGIN) {
    inner <- prod(d[seq_len(m - 1L)])
    one <- rep(rep(seq_len(d[m]), each = inner), length.out = prod(d))
    res <- res + (one - 1L) * mult
    mult <- mult * d[m]
  }
  array(as.integer(res + 1L), d)
}

as.array <- function(x, ...) UseMethod("as.array")

# as.array.default is a Go native (eval/funcs15.go)

write <- function(x, file = "data", ncolumns = if (is.character(x)) 1 else 5,
                  append = FALSE, sep = " ")
  cat(x, file = file, sep = c(rep.int(sep, ncolumns - 1), "\n"), append = append)

# bit operations on the classes "hexmode" and "octmode"
# as.hexmode and as.octmode are Go natives (eval/funcs15.go)

..digitsMode <- function(x, fmt, width) {
  s <- sprintf(fmt, unclass(x))
  s[is.na(x)] <- NA_character_
  attributes(s) <- attributes(unclass(x))
  if (is.null(width) && length(s) > 0L) width <- max(nchar(s), na.rm = TRUE)
  if (!is.null(width)) {
    pad <- pmax(width - nchar(s), 0L)
    s <- ifelse(is.na(s), NA_character_, paste0(strrep("0", pad), s))
  }
  s
}

format.hexmode <- function(x, width = NULL, upper.case = FALSE, ...)
  ..digitsMode(x, if (upper.case) "%X" else "%x", width)

format.octmode <- function(x, width = NULL, ...) ..digitsMode(x, "%o", width)

as.character.hexmode <- function(x, ...) {
  s <- sprintf("%x", unclass(x))
  s[is.na(x)] <- NA_character_
  s
}

as.character.octmode <- function(x, ...) {
  s <- sprintf("%o", unclass(x))
  s[is.na(x)] <- NA_character_
  s
}

print.hexmode <- function(x, ...) {
  print(format(x), ...)
  invisible(x)
}

print.octmode <- function(x, ...) {
  print(format(x), ...)
  invisible(x)
}

`[.hexmode` <- function(x, i) {
  cl <- oldClass(x)
  y <- NextMethod("[")
  oldClass(y) <- cl
  y
}

`[.octmode` <- function(x, i) {
  cl <- oldClass(x)
  y <- NextMethod("[")
  oldClass(y) <- cl
  y
}

`!.hexmode` <- function(a) as.hexmode(bitwNot(unclass(a)))
`&.hexmode` <- function(a, b) as.hexmode(bitwAnd(unclass(as.hexmode(a)), unclass(as.hexmode(b))))
`|.hexmode` <- function(a, b) as.hexmode(bitwOr(unclass(as.hexmode(a)), unclass(as.hexmode(b))))
`!.octmode` <- function(a) as.octmode(bitwNot(unclass(a)))
`&.octmode` <- function(a, b) as.octmode(bitwAnd(unclass(as.octmode(a)), unclass(as.octmode(b))))
`|.octmode` <- function(a, b) as.octmode(bitwOr(unclass(as.octmode(a)), unclass(as.octmode(b))))

diff <- function(x, ...) UseMethod("diff")

diff.default <- function(x, lag = 1L, differences = 1L, ...) {
  ismat <- is.matrix(x)
  n <- if (ismat) nrow(x) else length(x)
  if (length(lag) != 1L || length(differences) > 1L || lag < 1L || differences < 1L)
    stop("'lag' and 'differences' must be integers >= 1")
  if (lag * differences >= n) return(x[0L])
  r <- unclass(x)
  for (i in seq_len(differences)) {
    n <- if (ismat) nrow(r) else length(r)
    keep <- seq_len(n - lag)
    later <- keep + lag
    r <- if (ismat) r[later, , drop = FALSE] - r[keep, , drop = FALSE] else r[later] - r[keep]
  }
  class(r) <- oldClass(x)
  r
}

# all.equal and its methods are Go natives (eval/funcs15.go); the generics are here
all.equal <- function(target, current, ...) UseMethod("all.equal")

isSymmetric <- function(object, ...) UseMethod("isSymmetric")

isSymmetric.matrix <- function(object, tol = 100 * .Machine$double.eps, tol1 = 8 * tol, ...) {
  d <- dim(object)
  if (d[1L] != d[2L]) return(FALSE)
  if (!identical(rownames(object), colnames(object))) return(FALSE)
  isTRUE(all.equal(object, t(object), tolerance = tol, ...))
}
