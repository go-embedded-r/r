# Base library, slice 16 (1): version numbers and the description of R itself.
# Written from the help pages of R 4.6.2 (numeric_version, package_version,
# R.version, compareVersion) and from experiments with the oracle.

# (assigned, not defined with `<-`: they are data, not functions, and the
# coverage test of the functions does not apply to them)
assign("R.version", structure(list(
  platform = "aarch64-apple-darwin20.0.0", arch = "aarch64", os = "darwin20.0.0",
  system = "aarch64, darwin20.0.0", status = "", major = "4", minor = "6.2",
  year = "2026", month = "10", day = "29", `svn rev` = "90784", language = "R",
  version.string = "R version 4.6.2 (2026-10-29)", nickname = "Diversified Portfolio"),
  class = "simple.list"))
assign("version", R.version)
assign("R.version.string", "R version 4.6.2 (2026-10-29)")

# the patterns a version must match, by kind
.standard_regexps <- function()
  list(valid_package_name = "[[:alpha:]][[:alnum:].]*[[:alnum:]]",
       valid_package_version = "[[:digit:]]+([.-][[:digit:]]+)+",
       valid_R_system_version = "[[:digit:]]+[.][[:digit:]]+[.][[:digit:]]+",
       valid_numeric_version = "[[:digit:]]+([.-][[:digit:]]+)*")

# the list of integer vectors behind every version; a component list of length
# zero stands for NA
.make_numeric_version <- function(x, strict = TRUE, regexp, classes = NULL) {
  if (is.null(x)) x <- character()
  if (!is.character(x))
    stop(sprintf("invalid non-character version specification 'x' (type: %s)", typeof(x)), domain = NA)
  good <- !is.na(x) & grepl(paste0("^(", regexp, ")$"), x)
  if (strict && any(bad <- !good & !is.na(x)))
    stop(sprintf("invalid version specification %s", paste(sQuote(x[bad]), collapse = ", ")),
         call. = FALSE, domain = NA)
  out <- rep(list(integer()), length(x))
  out[good] <- lapply(strsplit(x[good], "[.-]"), as.integer)
  names(out) <- names(x)
  class(out) <- c(classes, "numeric_version")
  out
}

numeric_version <- function(x, strict = TRUE)
  .make_numeric_version(x, strict, .standard_regexps()$valid_numeric_version)

package_version <- function(x, strict = TRUE)
  .make_numeric_version(x, strict, .standard_regexps()$valid_package_version,
                        "package_version")

R_system_version <- function(x, strict = TRUE)
  .make_numeric_version(x, strict, .standard_regexps()$valid_R_system_version,
                        c("R_system_version", "package_version"))

is.numeric_version <- function(x) inherits(x, "numeric_version")
is.package_version <- function(x) inherits(x, "package_version")

as.numeric_version <- function(x) if (is.numeric_version(x)) x else numeric_version(x)

as.package_version <- function(x) {
  if (is.package_version(x)) return(x)
  if (is.list(x) && !is.null(names(x)) && all(c("major", "minor") %in% names(x)))
    return(R_system_version(paste(x[["major"]], x[["minor"]], sep = ".")))
  package_version(x)
}

getRversion <- function() R_system_version(paste(R.version$major, R.version$minor, sep = "."))

# a fixed-width key of each version, so that keys compare as versions do (a
# missing component counts as zero); NA for a missing version
..versionKey <- function(x) {
  v <- unclass(x)
  width <- max(0L, lengths(v))
  key <- vapply(v, function(p) {
    if (!length(p)) return(NA_character_)
    paste(sprintf("%010d", c(p, integer(width - length(p)))), collapse = "")
  }, "")
  names(key) <- NULL
  key
}

# -1, 0 or 1 for each pair of versions (NA when one is missing)
..versionCompare <- function(a, b) {
  na <- length(a)
  nb <- length(b)
  if (na == 0L || nb == 0L) return(integer())
  n <- max(na, nb)
  out <- integer(n)
  for (i in seq_len(n)) {
    p <- a[[(i - 1L) %% na + 1L]]
    q <- b[[(i - 1L) %% nb + 1L]]
    if (!length(p) || !length(q)) {
      out[i] <- NA_integer_
      next
    }
    w <- max(length(p), length(q))
    p <- c(p, integer(w - length(p)))
    q <- c(q, integer(w - length(q)))
    d <- which(p != q)
    out[i] <- if (length(d)) (if (p[d[1L]] < q[d[1L]]) -1L else 1L) else 0L
  }
  out
}

Ops.numeric_version <- function(e1, e2) {
  if (!(.Generic %in% c("<", "<=", ">", ">=", "==", "!=")))
    stop(sprintf("'%s' not defined for \"numeric_version\" objects", .Generic), domain = NA)
  e1 <- as.numeric_version(e1)
  e2 <- as.numeric_version(e2)
  d <- ..versionCompare(unclass(e1), unclass(e2))
  if (!length(d)) return(logical())
  switch(.Generic, "<" = d < 0L, "<=" = d <= 0L, ">" = d > 0L, ">=" = d >= 0L,
         "==" = d == 0L, "!=" = d != 0L)
}

Summary.numeric_version <- function(..., na.rm) {
  if (!(.Generic %in% c("max", "min", "range")))
    stop(sprintf("%s not defined for \"numeric_version\" objects", .Generic), domain = NA)
  x <- do.call(c, lapply(list(...), as.numeric_version))
  gone <- which(is.na(x))
  if (length(gone)) {
    # as R does: a missing version answers, even when na.rm is TRUE
    return(if (.Generic == "range") x[rep(gone[1L], 2L)] else x[gone[1L]])
  }
  k <- xtfrm(x)
  switch(.Generic, max = x[which.max(k)], min = x[which.min(k)],
         range = x[c(which.min(k), which.max(k))])
}

xtfrm.numeric_version <- function(x) xtfrm(..versionKey(x))

format.numeric_version <- function(x, ...) {
  out <- vapply(unclass(x), function(p) if (length(p)) paste(p, collapse = ".") else NA_character_, "")
  names(out) <- NULL
  out
}

as.character.numeric_version <- function(x, ...) format.numeric_version(x)

print.numeric_version <- function(x, quote = FALSE, ...) {
  y <- format.numeric_version(x)
  if (!length(y)) {
    cat("<0 elements>\n")
  } else {
    ok <- !is.na(y)
    y[ok] <- sQuote(y[ok])
    print(y, quote = quote, ...)
  }
  invisible(x)
}

`[.numeric_version` <- function(x, i, j) {
  y <- if (missing(j)) unclass(x)[i] else lapply(unclass(x)[i], "[", j)
  y[vapply(y, is.null, NA)] <- list(integer())
  class(y) <- oldClass(x)
  y
}

`[[.numeric_version` <- function(x, ..., exact = NA) {
  at <- list(...)
  if (length(at) < 2L) {
    y <- list(unclass(x)[[at[[1L]]]])
    class(y) <- oldClass(x)
    y
  } else {
    unclass(x)[[at[[1L]]]][at[[2L]]]
  }
}

`[<-.numeric_version` <- function(x, i, j, value) {
  y <- unclass(x)
  y[i] <- unclass(as.numeric_version(value))
  class(y) <- oldClass(x)
  y
}

`[[<-.numeric_version` <- function(x, i, value) {
  y <- unclass(x)
  y[[i]] <- unclass(as.numeric_version(value))[[1L]]
  class(y) <- oldClass(x)
  y
}

`is.na<-.numeric_version` <- function(x, value) {
  y <- unclass(x)
  y[value] <- list(integer())
  class(y) <- oldClass(x)
  y
}

is.na.numeric_version <- function(x) {
  out <- lengths(unclass(x)) == 0L
  names(out) <- NULL
  out
}

anyNA.numeric_version <- function(x, recursive = FALSE) any(is.na.numeric_version(x))

`$.package_version` <- function(x, name) {
  part <- pmatch(name, c("major", "minor", "patchlevel"))
  v <- unclass(x)
  switch(part,
         vapply(v, "[", 0L, 1L),
         vapply(v, "[", 0L, 2L),
         vapply(v, "[", 0L, 3L))
}

c.numeric_version <- function(..., recursive = FALSE) {
  parts <- lapply(list(...), as.numeric_version)
  out <- unlist(lapply(parts, unclass), recursive = FALSE)
  if (is.null(out)) out <- list()
  names(out) <- NULL
  class(out) <- if (length(parts) && all(vapply(parts, is.package_version, NA)))
    c("package_version", "numeric_version") else "numeric_version"
  out
}

rep.numeric_version <- function(x, ...) {
  y <- rep(unclass(x), ...)
  class(y) <- oldClass(x)
  y
}

duplicated.numeric_version <- function(x, incomparables = FALSE, ...)
  duplicated(..versionKey(x), incomparables = incomparables, ...)

unique.numeric_version <- function(x, incomparables = FALSE, ...)
  x[!duplicated.numeric_version(x, incomparables = incomparables, ...)]

as.list.numeric_version <- function(x, ...) {
  v <- unclass(x)
  out <- lapply(seq_along(v), function(i) structure(list(v[[i]]), class = oldClass(x)))
  names(out) <- names(x)
  out
}

# the order of two version strings, as utils::compareVersion gives it: -1, 0 or
# 1, where a version that is a prefix of the other is the smaller
compareVersion <- function(a, b) {
  a <- as.integer(strsplit(a, "[.-]")[[1L]])
  b <- as.integer(strsplit(b, "[.-]")[[1L]])
  common <- min(length(a), length(b))
  for (k in seq_len(common)) {
    if (a[k] > b[k]) return(1)   # a double, as the reference answers it
    if (a[k] < b[k]) return(-1L)
  }
  if (length(a) > common) 1L else if (length(b) > common) -1L else 0L
}
