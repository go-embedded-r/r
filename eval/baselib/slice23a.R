# Base library, slice 23 (1). Written from the help pages of R 4.6.2
# (abbreviate, agrep, Sys.setlocale, Sys.which, warnings, conditions,
# DateTimeClasses) and from experiments with the oracle.

# ---- abbreviate ------------------------------------------------------------

abbreviate <- function(names.arg, minlength = 4L, use.classes = TRUE, dot = FALSE,
                       strict = FALSE, method = c("left.kept", "both.sides"), named = TRUE) {
  x <- sub(" +$", "", sub("^ +", "", as.character(names.arg)))
  if (minlength <= 0L) {
    r <- rep.int("", length(x))
    if (named) names(r) <- x
    return(r)
  }
  minlength <- as.integer(minlength)
  if (is.na(minlength)) stop("invalid 'minlength' argument")
  method <- match.arg(method)
  if (..nonascii(x)) warning("abbreviate used with non-ASCII chars")
  u <- unique(x[!is.na(x)])
  ab <- ..abbrev1(u, minlength, use.classes)
  if (!strict) {
    ml <- minlength
    top <- max(nchar(u), 0L)
    repeat {
      dup <- ab %in% ab[duplicated(ab)]
      if (!any(dup) || ml > top) break
      ml <- ml + 1L
      ab[dup] <- ..abbrev1(u[dup], ml, use.classes)
    }
  }
  r <- ab[match(x, u)]
  if (dot) {
    changed <- !is.na(r) & r != x
    r[changed] <- paste0(r[changed], ".")
  }
  if (named) names(r) <- x
  r
}

# ---- agrep -----------------------------------------------------------------

.amatch_bounds <- function(x = 0.1) {
  nm <- c("cost", "insertions", "deletions", "substitutions", "all")
  bad <- "match distance components must be non-negative"
  if (!is.list(x)) {
    if (!length(x)) stop(bad)
    if (is.na(x) || x < 0) stop(bad)
    v <- suppressWarnings(as.numeric(x))
    if (is.na(v)) stop(bad)
    return(c(v, NA, NA, NA, NA))
  }
  val <- rep(NA_real_, 5L)
  names(val) <- nm
  pos <- pmatch(names(x), nm)
  if (is.null(names(x))) pos <- rep(NA_integer_, length(x))
  if (anyNA(pos)) warning("unknown match distance components ignored")
  keep <- !is.na(pos)
  v <- suppressWarnings(as.numeric(unlist(x[keep])))
  if (anyNA(v) || any(v < 0)) stop(bad)
  val[pos[keep]] <- v
  if (is.na(val[["cost"]]) && is.na(val[["all"]])) val[["all"]] <- 0.1
  for (i in 2:4) if (is.na(val[[i]])) val[[i]] <- val[["all"]]
  val
}

.amatch_costs <- function(x = NULL) {
  nm <- c("insertions", "deletions", "substitutions")
  costs <- c(insertions = 1, deletions = 1, substitutions = 1)
  if (is.null(x)) return(costs)
  v <- suppressWarnings(as.numeric(unlist(x)))
  if (anyNA(v) || any(v < 0)) stop("cost components must be non-negative")
  xn <- names(x)
  if (is.null(xn)) return(costs)
  pos <- pmatch(xn, nm)
  if (any(is.na(pos) & nzchar(xn))) warning("unknown cost components ignored")
  keep <- !is.na(pos)
  costs[pos[keep]] <- v[keep]
  costs
}

..agrepl <- function(pattern, x, max.distance, costs, ignore.case, fixed) {
  pattern <- as.character(pattern)
  if (!length(pattern)) stop("invalid 'pattern' argument")
  if (length(pattern) > 1L)
    warning("argument 'pattern' has length > 1 and only the first element will be used")
  pattern <- pattern[1L]
  if (!is.character(x)) x <- as.character(x)
  costs <- .amatch_costs(costs)
  bounds <- .amatch_bounds(max.distance)
  if (is.na(pattern)) return(rep.int(NA, length(x)))
  if (!nzchar(pattern)) stop("'pattern' must be a non-empty character string")
  if (any(costs != 1)) stop("non-default 'costs' are not supported")
  if (!fixed && grepl("[][\\\\.^$*+?(){}|]", pattern))
    stop("approximate matching of regular expressions is not supported; use fixed = TRUE")
  len <- nchar(pattern)
  bounds <- ifelse(is.na(bounds), NA_real_, ifelse(bounds < 1, ceiling(bounds * len), trunc(bounds)))
  ..agrepGo(pattern, x, ignore.case, bounds)
}

agrep <- function(pattern, x, max.distance = 0.1, costs = NULL, ignore.case = FALSE,
                  value = FALSE, fixed = TRUE, useBytes = FALSE) {
  hit <- ..agrepl(pattern, x, max.distance, costs, ignore.case, fixed)
  if (!is.character(x)) x <- as.character(x)
  if (anyNA(hit) && length(hit))
    return(if (value) rep(NA_character_, length(hit)) else rep(NA_integer_, length(hit)))
  if (value) x[hit] else which(hit)
}

agrepl <- function(pattern, x, max.distance = 0.1, costs = NULL, ignore.case = FALSE,
                   fixed = TRUE, useBytes = FALSE) {
  ..agrepl(pattern, x, max.distance, costs, ignore.case, fixed)
}

# ---- locale ----------------------------------------------------------------

..locale <- new.env()
..localeCats <- c("LC_COLLATE", "LC_CTYPE", "LC_MONETARY", "LC_NUMERIC", "LC_TIME",
                  "LC_MESSAGES", "LC_PAPER", "LC_MEASUREMENT")
..localeReset <- function() {
  for (k in ..localeCats) assign(k, "C.UTF-8", envir = ..locale)
  assign("LC_NUMERIC", "C", envir = ..locale)
  assign("LC_PAPER", "", envir = ..locale)
  assign("LC_MEASUREMENT", "", envir = ..locale)
}
..localeReset()

Sys.getlocale <- function(category = "LC_ALL") {
  category <- match(category, c("LC_ALL", ..localeCats))
  if (is.na(category)) stop("invalid 'category' argument")
  if (category > 1L) return(get(..localeCats[category - 1L], envir = ..locale))
  six <- c("LC_CTYPE", "LC_NUMERIC", "LC_TIME", "LC_COLLATE", "LC_MONETARY", "LC_MESSAGES")
  v <- vapply(six, get, "", envir = ..locale)
  if (all(v == v[1L])) return(unname(v[1L]))
  paste0(c(paste0(six, "=", v), paste0("LC_PAPER=", get("LC_PAPER", envir = ..locale)),
           "LC_NAME=C", "LC_ADDRESS=C", "LC_TELEPHONE=C",
           paste0("LC_MEASUREMENT=", get("LC_MEASUREMENT", envir = ..locale)),
           "LC_IDENTIFICATION=C"), collapse = ";")
}

Sys.setlocale <- function(category = "LC_ALL", locale = "") {
  cat <- match(category, c("LC_ALL", ..localeCats))
  if (is.na(cat)) stop("invalid 'category' argument")
  if (!is.character(locale) || length(locale) != 1L || is.na(locale))
    stop("invalid 'locale' argument")
  if (locale != "" && !locale %in% c("C", "POSIX", "C.UTF-8", "C.utf8")) {
    warning(gettextf("OS reports request to set locale to \"%s\" cannot be honored", locale),
            call. = FALSE)
    return("")
  }
  cats <- if (cat == 1L) ..localeCats else ..localeCats[cat - 1L]
  if (locale == "") {
    saved <- as.list(..locale)[setdiff(..localeCats, cats)]
    ..localeReset()
    for (k in names(saved)) assign(k, saved[[k]], envir = ..locale)
  } else for (k in cats) assign(k, if (k == "LC_NUMERIC" && cat == 1L) "C" else locale, envir = ..locale)
  Sys.getlocale(category)
}

Sys.localeconv <- function() {
  r <- c(".", "", "", "", "", "", "", "", "", "", rep("127", 8L))
  names(r) <- c("decimal_point", "thousands_sep", "grouping", "int_curr_symbol",
                "currency_symbol", "mon_decimal_point", "mon_thousands_sep", "mon_grouping",
                "positive_sign", "negative_sign", "int_frac_digits", "frac_digits",
                "p_cs_precedes", "p_sep_by_space", "n_cs_precedes", "n_sep_by_space",
                "p_sign_posn", "n_sign_posn")
  r
}

# a session has no host programs to find: every name is reported as not found
Sys.which <- function(names) {
  x <- as.character(names)
  r <- character(length(x))
  r[is.na(x)] <- NA_character_
  names(r) <- x
  r
}

# ---- warnings, restarts ----------------------------------------------------

print.restart <- function(x, ...) {
  cat(paste("<restart:", x[[1L]], ">\n"))
  invisible(x)
}

invokeRestartInteractively <- function(r) {
  if (!interactive()) stop("not an interactive session")
  invokeRestart(r)
}

summary.warnings <- function(object, ...) {
  key <- paste(names(object), vapply(object, function(cl) paste(deparse(cl), collapse = " "), ""))
  uk <- unique(key)
  r <- object[!duplicated(key)]
  class(r) <- "summary.warnings"
  attr(r, "counts") <- tabulate(match(key, uk), length(uk))
  r
}

print.summary.warnings <- function(x, ...) {
  n <- length(x)
  counts <- attr(x, "counts")
  if (!n) {
    cat("No warnings\n")
    return(invisible(x))
  }
  if (n == 1L) {
    cat(counts, " identical warnings:\n", sep = "")
    tags <- ""
  } else {
    cat("Summary of (a total of ", sum(counts), ") warning messages:\n", sep = "")
    tags <- paste0(counts, "x : ")
  }
  msgs <- names(x)
  for (i in seq_len(n)) {
    call <- unclass(x)[[i]]
    out <- if (is.null(call)) msgs[i] else {
      dp <- deparse(call, width.cutoff = 50L)
      head <- paste0("In ", dp[1L], if (length(dp) > 1L) " ..." else "", " :")
      if (nchar(dp[1L]) + nchar(msgs[i]) > 75L) paste0(head, "\n  ", msgs[i])
      else paste(head, msgs[i])
    }
    cat(tags[i], out, "\n", sep = "")
  }
  invisible(x)
}

# ---- balancePOSIXlt ----------------------------------------------------------

balancePOSIXlt <- function(x, fill.only = FALSE, classed = TRUE) {
  if (!is.list(x) || length(unclass(x)) < 9L)
    stop("a valid \"POSIXlt\" object is a list of at least 9 elements")
  if (is.na(fill.only)) stop("invalid 'fill.only' argument")
  if (is.na(classed)) stop("invalid 'classed' argument")
  tz <- attr(x, "tzone")
  if (fill.only) {
    n <- max(lengths(unclass(x)))
    was <- attr(x, "balanced")
    r <- lapply(unclass(x), function(e) if (length(e) == n || !length(e)) e else rep_len(e, n))
    attributes(r) <- list(names = names(unclass(x)))
    if (classed) class(r) <- c("POSIXlt", "POSIXt")
    attr(r, "tzone") <- tz
    attr(r, "balanced") <- if (isTRUE(was)) TRUE else NA
    return(r)
  }
  r <- as.POSIXlt(as.POSIXct(x), tz = if (is.null(tz)) "" else tz[1L])
  if (!is.null(tz) && !tz[1L] %in% c("UTC", "GMT") && !is.null(r$gmtoff)) r$gmtoff <- rep(NA_integer_, length(r$sec))
  attr(r, "tzone") <- tz
  attr(r, "balanced") <- TRUE
  if (!classed) {
    a <- attributes(r)
    r <- unclass(r)
    attributes(r) <- a[c("names", "tzone", "balanced")]
  }
  r
}
