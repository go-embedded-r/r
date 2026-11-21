# Slice 22: small base closures.
# Written from the help pages of R 4.6.2 and the oracle's answers
# (docs/m3-plan.md); BSD licence.

charmatch <- function(x, table, nomatch = NA_integer_) {
  if (is.factor(x)) x <- as.character(x)
  if (is.factor(table)) table <- as.character(table)
  ..charmatchChecked(x, table, nomatch)
}
..charmatchChecked <- function(x, table, nomatch) {
  x <- as.character(x)
  table <- as.character(table)
  ..charmatch(x, table, nomatch)
}

char.expand <- function(input, target, nomatch = stop("no match")) {
  if (length(input) != 1L)
    stop("'input' must have length 1")
  if (!(is.character(input) && is.character(target)))
    stop("'input' and 'target' must be character vectors")
  y <- charmatch(input, target)
  if (anyNA(y)) {
    # R's eval() is a frame of its own, so the default stop() reports eval(nomatch)
    if (missing(nomatch)) stop(simpleError("no match", quote(eval(nomatch))))
    eval(nomatch)
  }
  target[y]
}

findInterval <- function(x, vec, rightmost.closed = FALSE, all.inside = FALSE,
                         left.open = FALSE) {
  if (!identical(FALSE, is.unsorted(vec)))
    stop("'vec' must be sorted non-decreasingly and not contain NAs")
  ..findInterval(as.double(x), as.double(vec), rightmost.closed, all.inside,
                 left.open)
}

.makeMessage <- function(..., domain = NULL, appendLF = FALSE) {
  args <- list(...)
  msg <- if (length(args)) {
    args <- lapply(list(...), as.character)
    paste(unlist(args), collapse = "")
  } else ""
  if (appendLF) paste0(msg, "\n") else msg
}

.mapply <- function(FUN, dots, MoreArgs) {
  FUN <- match.fun(FUN)
  n <- length(dots)
  if (n == 0L) return(list())
  lens <- vapply(dots, length, 1L)
  if (any(lens == 0L)) return(list())
  m <- max(lens)
  nms <- names(dots)
  res <- vector("list", m)
  for (i in seq_len(m)) {
    a <- lapply(dots, function(d) d[[(i - 1L) %% length(d) + 1L]])
    if (!is.null(nms)) names(a) <- nms
    res[i] <- list(do.call(FUN, c(a, MoreArgs)))
  }
  res
}

dontCheck <- function(x) x

default.stringsAsFactors <- function()
  .Defunct("`stringsAsFactors = FALSE`")

.valid.factor <- function(object) {
  levs <- levels(object)
  if (!is.character(levs))
    return("factor levels must be \"character\"")
  if (d <- anyDuplicated(levs))
    return(sprintf("duplicated level [%d] in factor", d))
  TRUE
}

.handleSimpleError <- function(h, msg, call) h(simpleError(msg, call))

as.expression.default <- function(x, ...) as.vector(x, "expression")

.S3method <- function(generic, class, method) {
  if (missing(method)) method <- paste(generic, class, sep = ".")
  registerS3method(generic, class, method, envir = parent.frame())
  invisible(NULL)
}

.POSIXlt <- function(xx, tz = NULL, cl = c("POSIXlt", "POSIXt")) {
  attr(xx, "tzone") <- tz
  class(xx) <- cl
  xx
}

remove <- rm

comment <- function(x) attr(x, "comment", exact = TRUE)
`comment<-` <- function(x, value) {
  if (!is.null(value) && !is.character(value))
    stop("attempt to set invalid 'comment' attribute")
  attr(x, "comment") <- value
  x
}

reg.finalizer <- function(e, f, onexit = FALSE) {
  if (!is.environment(e)) stop("first argument must be environment or external pointer")
  if (!is.function(f)) stop("second argument must be a function")
  invisible(NULL)
}

R.Version <- function() unclass(R.version)

as.character.factor <- function(x, ...) levels(x)[as.integer(x)]
as.character.default <- function(x, ...) as.vector(unclass(x), "character")

`[.Dlist` <- function(x, i) structure(NextMethod("["), class = oldClass(x))
print.Dlist <- function(x, ...) {
  cat(formatDL(x, ...), sep = "\n")
  invisible(x)
}
