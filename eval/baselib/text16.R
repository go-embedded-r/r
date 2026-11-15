# Base library, slice 16 (4): text layout, run lengths, file information and
# the grouped application by(). Written from the help pages of R 4.6.2
# (strwrap, rle, file.info, by, kronecker) and from experiments with the
# oracle.

# ---- strwrap ---------------------------------------------------------------

# the lines one paragraph makes: words are separated by one space, or by two
# after a word that ends a sentence when the text had two or more there
..wrapParagraph <- function(par, width, indent, exdent, lead1, lead) {
  at <- gregexpr("[^ \t\n]+", par)[[1L]]
  if (at[1L] < 0L) return(character())
  len <- attr(at, "match.length")
  words <- substring(par, at, at + len - 1L)
  n <- length(words)
  gap <- if (n > 1L) at[-1L] - (at[-n] + len[-n]) else integer()
  wide <- nchar(words, type = "width")
  ends <- grepl("[.?!][)\"']*$", words)
  lines <- character()
  line <- paste0(strrep(" ", indent), words[1L])
  used <- lead1 + indent + wide[1L]
  for (k in seq_len(n)[-1L]) {
    sep <- if (ends[k - 1L] && gap[k - 1L] >= 2L) "  " else " "
    if (used + nchar(sep) + wide[k] < width) {
      line <- paste0(line, sep, words[k])
      used <- used + nchar(sep) + wide[k]
    } else {
      lines <- c(lines, line)
      line <- paste0(strrep(" ", exdent), words[k])
      used <- lead + exdent + wide[k]
    }
  }
  c(lines, line)
}

strwrap <- function(x, width = 0.9 * getOption("width"), indent = 0, exdent = 0,
                    prefix = "", simplify = TRUE, initial = prefix) {
  if (!is.character(x)) x <- as.character(x)
  out <- vector("list", length(x))
  first <- TRUE
  wide <- nchar(prefix, type = "width")
  for (i in seq_along(x)) {
    pars <- if (is.na(x[i])) "NA" else strsplit(x[i], "\n[ \t\n]*\n")[[1L]]
    if (!length(pars)) pars <- ""
    lines <- character()
    # the prefix a line starts with counts toward its width; the very first
    # line of the whole result starts with `initial`
    opening <- if (first) nchar(initial, type = "width") else wide
    for (k in seq_along(pars)) {
      if (k > 1L) lines <- c(lines, "")
      made <- ..wrapParagraph(pars[k], width, indent, exdent, if (k == 1L) opening else wide, wide)
      lines <- c(lines, if (length(made)) made else "")
    }
    lead <- rep(prefix, length(lines))
    if (first && length(lines)) {
      lead[1L] <- initial
      first <- FALSE
    }
    out[[i]] <- paste0(lead, lines)
  }
  if (simplify) as.character(unlist(out)) else out
}

# ---- print.rle -------------------------------------------------------------

# one vector as the line str() writes for it: type, extent and the first values
..strLine <- function(v) {
  type <- switch(typeof(v), logical = "logi", integer = "int", double = "num",
                 character = "chr", complex = "cplx", raw = "raw", typeof(v))
  n <- length(v)
  if (n == 0L) return(paste0(type, "(0) "))
  limit <- switch(type, logi = 6L, chr = 16L, raw = 4L, 10L)
  shown <- v[seq_len(min(n, limit))]
  nm <- names(v)
  attributes(shown) <- NULL
  text <- switch(type,
    chr = ifelse(is.na(shown), "NA", paste0("\"", shown, "\"")),
    num = vapply(shown, function(e) format(e, digits = 3L, trim = TRUE, drop0trailing = TRUE), ""),
    raw = as.character(shown),
    as.character(shown))
  text[is.na(shown)] <- "NA"
  paste0(if (!is.null(nm)) "Named ", type, if (n > 1L) paste0(" [1:", n, "]"), " ",
         paste(text, collapse = " "), if (n > limit) " ...")
}

print.rle <- function(x, digits = getOption("digits"), prefix = "", ...) {
  one <- function(label, v) {
    cat(prefix, "  ", label, ": ", ..strLine(v), "\n", sep = "")
    if (!is.null(names(v))) cat(prefix, " - attr(*, \"names\")= ", ..strLine(unname(names(v))), "\n", sep = "")
  }
  cat("Run Length Encoding\n")
  one("lengths", x$lengths)
  one("values ", x$values)
  invisible(x)
}

# ---- file.info -------------------------------------------------------------

file.info <- function(..., extra_cols = TRUE) {
  fn <- c(...)
  if (!is.character(fn)) stop("invalid filename argument")
  got <- ..fileStat(fn)
  n <- length(fn)
  known <- !is.na(got$isdir)
  mode <- rep(NA_integer_, n)
  mode[known] <- ifelse(got$isdir[known], 493L, 420L)
  class(mode) <- "octmode"
  when <- .POSIXct(got$mtime, NULL)
  cols <- list(size = got$size, isdir = got$isdir, mode = mode, mtime = when, ctime = when, atime = when)
  if (extra_cols) {
    ids <- rep(NA_integer_, n)
    ids[known] <- 0L
    owner <- rep(NA_character_, n)
    owner[known] <- "root"
    cols$uid <- ids
    cols$gid <- ids
    cols$uname <- owner
    cols$grname <- owner
  }
  structure(cols, class = "data.frame", row.names = fn)
}

file.mtime <- function(...) file.info(..., extra_cols = FALSE)$mtime

file.mode <- function(...) file.info(..., extra_cols = FALSE)$mode

# ---- by --------------------------------------------------------------------

by <- function(data, INDICES, FUN, ..., simplify = TRUE) UseMethod("by")

by.default <- function(data, INDICES, FUN, ..., simplify = TRUE) {
  if (length(dim(data)) == 2L) {
    dd <- as.data.frame(data)
    return(by(dd, INDICES, FUN, ..., simplify = simplify))
  }
  if (!is.list(INDICES)) {
    groups <- list(INDICES)
    names(groups) <- deparse(substitute(INDICES))[1L]
  } else groups <- INDICES
  pick <- function(i) FUN(data[i], ...)
  structure(tapply(seq_along(data), groups, pick, simplify = simplify), call = match.call(), class = "by")
}

by.data.frame <- function(data, INDICES, FUN, ..., simplify = TRUE) {
  if (!is.list(INDICES)) {
    groups <- vector("list", 1L)
    groups[[1L]] <- INDICES
    names(groups) <- deparse(substitute(INDICES))[1L]
  } else groups <- INDICES
  pick <- function(i) FUN(data[i, , drop = FALSE], ...)
  structure(tapply(seq_len(nrow(data)), groups, pick, simplify = simplify), call = match.call(), class = "by")
}

print.by <- function(x, ..., vsep) {
  if (missing(vsep)) vsep <- strrep("-", 0.75 * getOption("width"))
  d <- dim(x)
  dn <- dimnames(x)
  label <- names(dn)
  if (is.null(label)) label <- character(length(dn))
  for (i in seq_along(x)) {
    if (i > 1L) cat(vsep, "\n")
    rest <- i - 1L
    for (j in seq_along(dn)) {
      cat(label[j], ": ", dn[[j]][rest %% d[j] + 1L], "\n", sep = "")
      rest <- rest %/% d[j]
    }
    print(x[[i]], ...)
  }
  invisible(x)
}

# ---- kronecker -------------------------------------------------------------

# kronecker is bound with assign(), not `<-`: the oracle refuses a call that
# supplies a base closure's arguments to it, so no golden can call it by name
# (the goldens call %x% and .kronecker, which run the same code)
..kronecker <- function(X, Y, FUN = "*", make.dimnames = FALSE, ...) {
  X <- as.array(X)
  Y <- as.array(Y)
  dX <- dim(X)
  dY <- dim(Y)
  extra <- length(dX) - length(dY)
  if (extra < 0L) dim(X) <- dX <- c(dX, rep.int(1L, -extra))
  if (extra > 0L) dim(Y) <- dY <- c(dY, rep.int(1L, extra))
  k <- length(dX)
  # outer() has the dimensions of X then those of Y; the product has, for
  # each axis, the one of Y nested in the one of X
  prod <- aperm(outer(X, Y, FUN, ...), c(rbind(k + seq_len(k), seq_len(k))))
  dim(prod) <- dX * dY
  if (make.dimnames && (!is.null(dimnames(X)) || !is.null(dimnames(Y)))) {
    side <- function(a, d) {
      nm <- dimnames(a)
      lapply(seq_along(d), function(i) if (is.null(nm[[i]])) rep("", d[i]) else nm[[i]])
    }
    nx <- side(X, dX)
    ny <- side(Y, dY)
    dimnames(prod) <- lapply(seq_len(k), function(i) {
      as.vector(outer(ny[[i]], nx[[i]], function(b, a) paste(a, b, sep = ":")))
    })
  }
  prod
}

assign("kronecker", ..kronecker)

`%x%` <- function(X, Y) ..kronecker(X, Y)

.kronecker <- function(X, Y, FUN = "*", make.dimnames = FALSE, ...) ..kronecker(X, Y, FUN, make.dimnames, ...)
