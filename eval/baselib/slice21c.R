# Slice 21, part C: the methods of POSIXlt, through the POSIXct of the same
# instants. Written from the help pages of R 4.6.2 and the oracle's answers
# (docs/m3-plan.md); BSD licence.

mean.POSIXlt <- function(x, ...) as.POSIXlt(mean(as.POSIXct(x), ...))

sort.POSIXlt <- function(x, decreasing = FALSE, na.last = NA, ...) {
  x[order(as.POSIXct(x), na.last = na.last, decreasing = decreasing)]
}

duplicated.POSIXlt <- function(x, incomparables = FALSE, ...) duplicated(as.POSIXct(x), incomparables, ...)

unique.POSIXlt <- function(x, incomparables = FALSE, ...) x[!duplicated(as.POSIXct(x), incomparables, ...)]

anyNA.POSIXlt <- function(x, recursive = FALSE) anyNA(as.POSIXct(x))

rep.POSIXlt <- function(x, ...) as.POSIXlt(rep(as.POSIXct(x), ...))

is.finite.POSIXlt <- function(x) is.finite(as.POSIXct(x))
is.infinite.POSIXlt <- function(x) is.infinite(as.POSIXct(x))
is.nan.POSIXlt <- function(x) is.nan(as.POSIXct(x))

`length<-.POSIXlt` <- function(x, value) x[seq_len(value)]

as.double.POSIXlt <- function(x, ...) as.double(as.POSIXct(x))
xtfrm.POSIXlt <- function(x) as.double(as.POSIXct(x))

as.list.POSIXlt <- function(x, ...) {
  r <- lapply(seq_len(length(x)), function(i) {
    e <- x[i]
    attr(e, "balanced") <- NULL
    e
  })
  names(r) <- names(x)
  r
}

mtfrm.POSIXlt <- function(x) as.list(x)

as.vector.POSIXlt <- function(x, mode = "any") {
  parts <- as.list(x)
  if (mode == "any" || mode == "list") return(parts)
  if (mode == "character") {
    return(vapply(parts, function(e) {
      e <- unclass(e)
      attributes(e) <- list(names = names(e))
      e[] <- lapply(e, function(v) if (is.numeric(v)) as.double(v) else v)
      paste(deparse(e, width.cutoff = 500L), collapse = " ")
    }, "", USE.NAMES = FALSE))
  }
  as.vector(parts, mode)
}

as.matrix.POSIXlt <- function(x, ...) {
  u <- unclass(x)
  cols <- u[names(u) != "zone"]
  do.call(cbind, lapply(cols, as.double))
}

Summary.POSIXlt <- function(..., na.rm) {
  if (!(.Generic %in% c("max", "min", "range"))) {
    stop(gettextf("'%s' not defined for \"POSIXt\" objects", .Generic), domain = NA)
  }
  args <- lapply(list(...), as.POSIXct)
  as.POSIXlt(do.call(.Generic, c(args, na.rm = !missing(na.rm) && na.rm)))
}

julian.POSIXt <- function(x, origin = as.POSIXct("1970-01-01", tz = "GMT"), ...) {
  origin <- as.POSIXct(origin)
  if (length(origin) != 1L) stop("'origin' must be of length 1")
  structure(difftime(as.POSIXct(x), origin, units = "days"), origin = origin)
}
