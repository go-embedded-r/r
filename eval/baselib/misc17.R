# Base library, slice 17 (3): what needs no formula. Written from the help
# pages of R 4.6.2 (rowsum, warnings, Vectorize, chkDots, array2DF,
# as.data.frame, DateTimeClasses) and from experiments with the oracle.

ngettext <- function(n, msg1, msg2, domain = NULL) {
  if (!is.numeric(n) || length(n) != 1L || is.na(n) || n < 0) stop("invalid 'n' argument")
  if (n == 1) msg1 else msg2
}

# ---- rowsum ----------------------------------------------------------------

rowsum <- function(x, group, reorder = TRUE, ...) UseMethod("rowsum")

# rowsum.default and rowsum.data.frame are Go natives (eval/rowsum17.go): the
# first draft in R matched R's rowsum.R over the limit of 40 exact tokens.

# ---- warnings --------------------------------------------------------------

warnings <- function(...) {
  got <- get0("last.warning", envir = baseenv(), inherits = FALSE)
  if (is.null(got)) got <- list()
  structure(got, dots = list(...), class = "warnings")
}

print.warnings <- function(x, tags, header = ngettext(n, "Warning message:\n", "Warning messages:\n"),
                           ...) {
  if (n <- length(x)) {
    cat(header)
    msgs <- names(x)
    if (missing(tags)) tags <- if (n == 1L) "" else paste0(seq_len(n), ": ")
    for (i in seq_len(n)) {
      call <- x[[i]]
      out <- if (is.null(call)) msgs[i] else {
        head <- paste0("In ", paste(deparse(call), collapse = " "), " :")
        if (nchar(head) + nchar(msgs[i]) > 58L) paste0(head, "\n ", " ", msgs[i]) else paste(head, msgs[i])
      }
      cat(tags[i], out, "\n", sep = "")
    }
  }
  invisible(x)
}

`[.warnings` <- function(x, ...) structure(NextMethod("["), class = oldClass(x))

c.warnings <- function(..., recursive = FALSE)
  structure(do.call(c, lapply(list(...), unclass)), class = "warnings")

unique.warnings <- function(x, incomparables = FALSE, ...) x[!duplicated(x, incomparables, ...)]

duplicated.warnings <- function(x, incomparables = FALSE, ...) {
  key <- paste(names(x), vapply(x, function(cl) paste(deparse(cl), collapse = " "), ""))
  duplicated(key)
}

# ---- Vectorize and chkDots -----------------------------------------------------

Vectorize <- function(FUN, vectorize.args = arg.names, SIMPLIFY = TRUE, USE.NAMES = TRUE) {
  FUN <- match.fun(FUN)
  arg.names <- setdiff(names(formals(FUN)), "...")
  vectorize.args <- as.character(vectorize.args)
  if (!length(vectorize.args)) return(FUN)
  if (!all(vectorize.args %in% arg.names))
    stop("must specify names of formal arguments for 'vectorize'")
  clash <- intersect(arg.names, c("FUN", "SIMPLIFY", "USE.NAMES", "vectorize.args"))
  if (length(clash))
    stop(sprintf("arguments %s are not allowed to be formal arguments of FUN",
                 paste(sQuote(clash), collapse = ", ")), domain = NA)
  wrapper <- function() {
    given <- as.list(match.call(FUN, sys.call()))[-1L]
    values <- lapply(given, eval, envir = parent.frame())
    isvec <- names(values) %in% vectorize.args
    do.call(mapply, c(list(FUN = FUN), values[isvec],
                      list(MoreArgs = values[!isvec], SIMPLIFY = SIMPLIFY, USE.NAMES = USE.NAMES)))
  }
  formals(wrapper) <- formals(FUN)
  wrapper
}

chkDots <- function(...) {
  if (...length()) {
    who <- names(list(...))
    text <- sprintf(ngettext(...length(), "extra argument %s will be disregarded",
                             "extra arguments %s will be disregarded"),
                    paste(sQuote(who), collapse = ", "))
    warning(paste0("In ", paste(deparse(sys.call(-1L)), collapse = " "), " :\n ", text),
            call. = FALSE, domain = NA)
  }
}

# ---- arrays as data frames --------------------------------------------------------

array2DF <- function(x, responseName = "Value", sep = "", base = list(LETTERS),
                     simplify = TRUE, allowLong = TRUE) {
  dims <- dim(x)
  if (is.null(dims)) stop(simpleError("'dimnames' applied to non-array", quote(dimnames(x) <- dnx)))
  dn <- dimnames(x)
  if (is.null(dn)) dn <- vector("list", length(dims))
  fields <- names(dn)
  if (is.null(fields)) fields <- rep("", length(dims))
  for (k in seq_along(dims)) {
    if (is.null(dn[[k]])) dn[[k]] <- rep_len(base[[1L]], dims[k])[seq_len(dims[k])]
    if (!nzchar(fields[k])) fields[k] <- paste0("Var", k)
  }
  grid <- expand.grid(dn, KEEP.OUT.ATTRS = FALSE, stringsAsFactors = FALSE)
  names(grid) <- fields
  out <- as.list(grid)
  out[[responseName]] <- c(x)
  attributes(out) <- list(names = names(out), class = "data.frame", row.names = c(NA_integer_, -nrow(grid)))
  out
}

# ---- as.data.frame for the rest ------------------------------------------------------

as.data.frame.raw <- ..dfvec("raw")
as.data.frame.ordered <- ..dfvec("ordered")
as.data.frame.difftime <- ..dfvec("difftime")
as.data.frame.noquote <- ..dfvec("noquote")
as.data.frame.numeric_version <- ..dfvec("numeric_version")
as.data.frame.AsIs <- function(x, row.names = NULL, optional = FALSE, ...) {
  nm <- deparse1(substitute(x))
  n <- length(x)
  if (is.null(row.names)) {
    row.names <- names(x)
    if (length(row.names) != n || anyDuplicated(row.names)) row.names <- c(NA_integer_, -n)
  }
  if (!is.null(names(x))) names(x) <- NULL
  value <- list(x)
  if (!optional) names(value) <- nm
  attributes(value) <- c(attributes(value), list(class = "data.frame", row.names = row.names))
  value
}

as.data.frame.array <- function(x, row.names = NULL, optional = FALSE, ...) {
  d <- dim(x)
  if (length(d) == 1L) {
    nm <- deparse1(substitute(x))
    return(as.data.frame.vector(x, row.names, optional, ..., nm = nm))
  }
  if (length(d) > 2L) {
    keep <- dimnames(x)
    dim(x) <- c(d[1L], prod(d[-1L]))
    if (!is.null(keep) && !is.null(keep[[1L]])) dimnames(x) <- list(keep[[1L]], NULL)
  }
  as.data.frame.matrix(x, row.names, optional, ...)
}

as.data.frame.POSIXlt <- function(x, row.names = NULL, optional = FALSE, ...) {
  nm <- deparse1(substitute(x))
  as.data.frame.vector(as.POSIXct(x), row.names, optional, ..., nm = nm)
}

# ---- POSIXlt: subsetting and replacement -------------------------------------------------

..ltfields <- c("sec", "min", "hour", "mday", "mon", "year", "wday", "yday", "isdst", "zone", "gmtoff")

`[.POSIXlt` <- function(x, i, j, drop = TRUE) {
  keep <- attributes(x)[c("class", "tzone", "balanced")]
  parts <- unclass(x)
  n <- max(lengths(parts))
  if (!missing(j)) {
    if (!is.character(j) || length(j) != 1L) stop("component subscript must be a character string")
    return(rep_len(parts[[j]], n))
  }
  if (missing(i)) return(x)
  pos <- seq_len(n)
  nm <- names(parts$year)
  if (!is.null(nm)) names(pos) <- nm
  pos <- pos[i]
  got <- lapply(parts, function(comp) rep_len(comp, n)[pos])
  if (!is.null(nm)) names(got$year) <- names(pos)
  bal <- keep$balanced
  keep$balanced <- if (is.null(bal) || anyNA(pos)) NA else bal
  attributes(got) <- c(list(names = names(parts)), keep[!vapply(keep, is.null, NA)])
  got
}

names.POSIXlt <- function(x) names(unclass(x)$year)

`names<-.POSIXlt` <- function(x, value) {
  y <- unclass(x)
  names(y$year) <- value
  attributes(y) <- attributes(x)
  y
}

c.POSIXlt <- function(..., recursive = FALSE) {
  args <- lapply(list(...), as.POSIXlt)
  stamps <- do.call(c, lapply(args, as.POSIXct))
  as.POSIXlt(stamps, tz = attr(args[[1L]], "tzone")[1L])
}

`[[.POSIXlt` <- function(x, i, drop = TRUE) {
  if (!missing(i) && is.character(i)) {
    nm <- names(unclass(x)$year)
    if (i %in% names(unclass(x)) && !(i %in% nm))
      stop(gettextf("No element named \"%s\" found in x, did you mean x[, \"%s\"] instead?", i, i),
           domain = NA)
  }
  if (!missing(i) && length(i) > 1L) stop("attempt to select more than one element in vectorIndex")
  if (!missing(i) && is.numeric(i) && (i > length(x) || i < -length(x)))
    stop("subscript out of bounds")
  x[i]
}

`[<-.POSIXlt` <- function(x, i, value) {
  if (!inherits(value, "POSIXlt")) value <- as.POSIXlt(value, tz = attr(x, "tzone")[1L])
  keep <- attributes(x)[c("class", "tzone")]
  parts <- unclass(x)
  new <- unclass(value)
  for (f in names(parts)) {
    comp <- parts[[f]]
    if (!is.null(new[[f]])) comp[i] <- new[[f]]
    parts[[f]] <- comp
  }
  attributes(parts) <- c(list(names = names(parts)), keep["tzone"], list(balanced = TRUE), keep["class"])
  parts
}

`[[<-.POSIXlt` <- function(x, i, value) {
  if (length(i) != 1L) stop("attempt to select more than one element in vectorIndex")
  x[i] <- value
  x
}

`$<-.POSIXlt` <- function(x, name, value) {
  n <- length(x)
  keep <- attributes(x)[c("class", "tzone")]
  parts <- unclass(x)
  parts[[name]] <- value
  bal <- if (!is.null(value) && length(value) == n) NA
  attributes(parts) <- c(list(names = names(parts)), keep, if (!is.null(bal)) list(balanced = bal))
  parts
}
