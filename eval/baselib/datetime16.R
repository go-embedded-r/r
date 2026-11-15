# Base library, slice 16 (2): rounding, truncating, cutting and sequencing of
# dates and times, and their summaries. Written from the help pages of R 4.6.2
# (round.POSIXt, cut.Date, seq.POSIXt, summary.Date) and from experiments with
# the oracle; every layout below (which attributes a result has, and in which
# order) is what the oracle answered.

# a POSIXlt from its fields, class first and then the attributes the reference
# gives (tzone, and balanced when it has one)
..ltFrom <- function(fields, tzone, balanced) {
  class(fields) <- c("POSIXlt", "POSIXt")
  attr(fields, "tzone") <- tzone
  if (!is.null(balanced)) attr(fields, "balanced") <- balanced
  fields
}

# the start of the period of one unit that holds each time, as a POSIXlt
trunc.POSIXt <- function(x, units = c("secs", "mins", "hours", "days", "months", "years"), ...) {
  units <- match.arg(units)
  lt <- as.POSIXlt(x)
  f <- ..bareLt(lt)
  zone <- attr(lt, "tzone")
  if (units == "months" || units == "years") {
    live <- !is.na(f$year)
    f$sec[live] <- 0
    f$min[live] <- 0L
    f$hour[live] <- 0L
    f$mday[live] <- 1L
    if (units == "years") f$mon[live] <- 0L
    f$isdst[] <- -1L
    secs <- unclass(as.POSIXct(..ltFrom(f, zone, NULL)))
    attributes(secs) <- NULL
    return(..mkLt(secs, zone[1L], names(f$year)))
  }
  if (units == "secs") f$sec <- trunc(f$sec) else f$sec[] <- 0
  balanced <- NULL
  if (units == "secs" || units == "mins") {
    balanced <- NA
  } else {
    f$min[] <- 0L
    if (units == "days") {
      f$hour[] <- 0L
      f$isdst[] <- -1L
    }
  }
  if (units == "mins") f$sec[] <- 0
  ..ltFrom(f, zone, balanced)
}

# the nearest start of a period: half a unit later, then truncated; for a
# month or a year, whichever of the two bounding starts is nearer, the later
# one when they are equally near
round.POSIXt <- function(x, digits = c("secs", "mins", "hours", "days", "months", "years")) {
  units <- if (is.numeric(digits) && digits == 0) "secs" else match.arg(digits)
  ct <- as.POSIXct(x)
  zone <- attr(ct, "tzone")
  if (units == "months" || units == "years") {
    lo <- as.POSIXct(trunc.POSIXt(x, units))
    f <- ..bareLt(as.POSIXlt(lo))
    if (units == "months") f$mon <- f$mon + 1L else f$year <- f$year + 1L
    f$isdst[] <- -1L
    hi <- as.POSIXct(..ltFrom(f, attr(as.POSIXlt(lo), "tzone"), NULL))
    later <- unclass(ct) - unclass(lo) >= (unclass(hi) - unclass(lo)) / 2
    nearer <- ifelse(later, unclass(hi), unclass(lo))
    attributes(nearer) <- NULL
    out <- trunc.POSIXt(.POSIXct(nearer, zone), units)
    if (any(later, na.rm = TRUE)) {
      # when a time goes up the reference leaves the class after the other
      # attributes
      cl <- oldClass(out)
      class(out) <- NULL
      class(out) <- cl
    }
    return(out)
  }
  half <- c(secs = 0.5, mins = 30, hours = 1800, days = 43200)[[units]]
  trunc.POSIXt(.POSIXct(unclass(ct) + half, zone), units)
}

trunc.Date <- function(x, units = c("secs", "mins", "hours", "days", "months", "years"), ...) {
  units <- match.arg(units)
  if (units == "months" || units == "years") return(as.Date(trunc.POSIXt(as.POSIXlt(x), units)))
  v <- unclass(x)
  storage.mode(v) <- "double"
  .Date(floor(v), oldClass(x))
}

round.Date <- function(x, ...) {
  # the reference runs the next method, so it fails when called by name
  if (!exists(".Generic", inherits = FALSE)) stop("generic function not specified")
  v <- unclass(x)
  storage.mode(v) <- "double"
  .Date(round(v, ...), oldClass(x))
}

# ---- seq.POSIXt ------------------------------------------------------------

seq.POSIXt <- function(from, to, by, length.out = NULL, along.with = NULL, ...) {
  if (missing(from)) stop("'from' must be specified")
  if (!inherits(from, "POSIXt")) stop("'from' must be a \"POSIXt\" object")
  start <- as.POSIXct(from)
  if (length(start) != 1L) stop("'from' must be of length 1")
  zone <- attr(start, "tzone")
  if (!missing(to)) {
    if (!inherits(to, "POSIXt")) stop("'to' must be a \"POSIXt\" object")
    if (length(as.POSIXct(to)) != 1L) stop("'to' must be of length 1")
    stop_at <- as.POSIXct(to)
  }
  if (!missing(along.with)) length.out <- length(along.with)
  else if (!is.null(length.out)) {
    if (length(length.out) != 1L) stop("'length.out' must be of length 1")
    length.out <- ceiling(length.out)
  }
  if (sum(c(!missing(to), !missing(by), !is.null(length.out))) != 2L)
    stop("exactly three of 'to', 'from', 'by' and 'length.out' / 'along.with' must be specified")
  a <- as.numeric(unclass(start))
  if (missing(by))
    return(.POSIXct(..wholeInts(seq.int(a, as.numeric(unclass(stop_at)), length.out = length.out)), zone))
  if (length(by) != 1L) stop("'by' must be of length 1")
  mult <- 1
  kind <- "secs"
  if (inherits(by, "difftime")) {
    mult <- ..asSecs(by)
  } else if (is.character(by)) {
    words <- strsplit(by, " ", fixed = TRUE)[[1L]]
    if (length(words) > 2L || length(words) < 1L) stop("invalid 'by' string")
    kinds <- c("secs", "mins", "hours", "days", "weeks", "months", "years", "DSTdays", "quarters")
    k <- pmatch(words[length(words)], kinds)
    if (is.na(k)) stop("invalid string for 'by'")
    kind <- kinds[k]
    mult <- if (length(words) == 2L) as.integer(words[1L]) else 1L
    if (is.na(mult)) stop("invalid 'by' string")
  } else if (is.numeric(by)) {
    mult <- by
  } else stop("invalid mode for 'by'")
  if (kind %in% c("secs", "mins", "hours", "days", "weeks")) {
    step <- mult * c(secs = 1, mins = 60, hours = 3600, days = 86400, weeks = 604800)[[kind]]
    out <- if (!missing(to)) seq.int(a, as.numeric(unclass(stop_at)), by = step)
           else ..wholeInts(seq.int(a, by = step, length.out = length.out))
    return(.POSIXct(out, zone))
  }
  months <- kind %in% c("months", "years", "quarters")
  if (kind == "years") mult <- mult * 12L
  if (kind == "quarters") mult <- mult * 3L
  lt <- as.POSIXlt(start)
  if (!missing(to)) {
    end <- as.POSIXlt(stop_at)
    span <- if (months) abs((end$year - lt$year) * 12L + (end$mon - lt$mon))
            else abs(as.numeric(unclass(stop_at)) - a) / 86400
    count <- floor(span / abs(mult)) + 2L
  } else count <- length.out
  idx <- (seq_len(max(count, 0L)) - 1L) * mult
  f <- ..bareLt(lt)
  n <- length(idx)
  for (nm in c("sec", "min", "hour", "mday", "mon", "year")) f[[nm]] <- rep_len(f[[nm]], n)
  if (months) {
    total <- lt$mon + idx
    f$year <- lt$year + total %/% 12L
    f$mon <- total %% 12L
  } else {
    f$mday <- lt$mday + idx
  }
  f$isdst <- rep_len(-1L, n)
  f$wday <- rep_len(NA_integer_, n)
  f$yday <- rep_len(NA_integer_, n)
  f$zone <- rep_len("", n)
  f$gmtoff <- rep_len(NA_integer_, n)
  out <- as.numeric(unclass(as.POSIXct(..ltFrom(f, attr(lt, "tzone"), NULL))))
  if (!missing(to)) out <- if (mult > 0) out[out <= as.numeric(unclass(stop_at))] else out[out >= as.numeric(unclass(stop_at))]
  .POSIXct(out, zone)
}

# ---- cut -------------------------------------------------------------------

# the breaks written as "2 months": the unit it names and how many of it
..breakSpec <- function(spec, kinds) {
  words <- strsplit(spec, " ", fixed = TRUE)[[1L]]
  if (length(words) > 2L || length(words) < 1L) stop("invalid specification of 'breaks'")
  k <- pmatch(words[length(words)], kinds)
  if (is.na(k)) stop("invalid specification of 'breaks'")
  mult <- if (length(words) == 2L) as.integer(words[1L]) else 1L
  if (is.na(mult)) stop("invalid specification of 'breaks'")
  list(kind = kinds[k], mult = mult)
}

# the number of break points a specification asks for, refused when absurd
..checkBreaks <- function(n) {
  if (is.na(n) || n > 16777216) stop("'breaks' asks for too many intervals")
  as.integer(n)
}

# the factor of the intervals each time falls in, given the break values as
# numbers, the labels the caller gave and the labels to use otherwise
..cutFactor <- function(v, breaks, labels, right, fallback, ...) {
  codes <- cut(v, breaks, labels = FALSE, right = right, ...)
  if (is.logical(labels) && !labels) return(codes)
  nlev <- length(breaks) - 1L
  lev <- if (is.null(labels)) fallback(codes) else as.character(labels)
  if (length(lev) != nlev && is.null(labels) && !length(lev)) lev <- character(nlev)
  structure(codes, levels = lev, class = "factor")
}

cut.Date <- function(x, breaks, labels = NULL, start.on.monday = TRUE, right = FALSE, ...) {
  if (!inherits(x, "Date")) stop("'x' must be a date-time object")
  x <- as.Date(x)
  fallback <- NULL
  if (inherits(breaks, "Date")) {
    breaks <- sort(as.Date(breaks))
    fallback <- function(codes) as.character(breaks[-length(breaks)])
    num <- unclass(breaks)
  } else if (is.numeric(breaks) && length(breaks) == 1L) {
    num <- breaks
    fallback <- function(codes) as.character(x[!duplicated(codes)])
  } else if (is.character(breaks) && length(breaks) == 1L) {
    spec <- ..breakSpec(breaks, c("days", "weeks", "months", "years", "DSTdays", "quarters"))
    first <- min(x, na.rm = TRUE)
    last <- max(x, na.rm = TRUE)
    lt <- as.POSIXlt(first)
    if (spec$kind == "weeks") {
      first <- first - (lt$wday - if (start.on.monday) 1L else 0L) %% 7L
      step <- 7L * spec$mult
    } else if (spec$kind %in% c("days", "DSTdays")) {
      step <- spec$mult
    } else {
      mon <- switch(spec$kind, months = lt$mon, years = 0L, quarters = lt$mon - lt$mon %% 3L)
      first <- as.Date(ISOdate(lt$year + 1900L, mon + 1L, 1L, 0, tz = "GMT"))
      step <- spec$mult * switch(spec$kind, months = 1L, years = 12L, quarters = 3L)
    }
    n <- ..checkBreaks(ceiling(as.numeric(last - first) / (if (spec$kind %in% c("days", "DSTdays", "weeks")) step else step * 28)) + 2)
    pts <- if (spec$kind %in% c("days", "DSTdays", "weeks")) first + step * (seq_len(n) - 1L)
           else seq(first, by = paste(step, "months"), length.out = n)
    keep <- which(pts > last)[1L]
    breaks <- pts[seq_len(keep)]
    fallback <- function(codes) as.character(breaks[-length(breaks)])
    num <- unclass(breaks)
  } else stop("invalid specification of 'breaks'")
  ..cutFactor(unclass(x), num, labels, right, fallback, ...)
}

cut.POSIXt <- function(x, breaks, labels = NULL, start.on.monday = TRUE, right = FALSE, ...) {
  if (!inherits(x, "POSIXt")) stop("'x' must be a date-time object")
  x <- as.POSIXct(x)
  zone <- attr(x, "tzone")
  fallback <- NULL
  if (inherits(breaks, "POSIXt")) {
    breaks <- sort(as.POSIXct(breaks))
    fallback <- function(codes) as.character(breaks[-length(breaks)])
    num <- unclass(breaks)
  } else if (is.numeric(breaks) && length(breaks) == 1L) {
    num <- breaks
    fallback <- function(codes) as.character(x[!duplicated(codes)])
  } else if (is.character(breaks) && length(breaks) == 1L) {
    spec <- ..breakSpec(breaks, c("secs", "mins", "hours", "days", "weeks", "months", "years", "DSTdays", "quarters"))
    first <- min(x, na.rm = TRUE)
    last <- max(x, na.rm = TRUE)
    unit <- switch(spec$kind, secs = "secs", mins = "mins", hours = "hours", "days")
    if (spec$kind %in% c("months", "years", "quarters")) unit <- "months"
    if (spec$kind == "years") unit <- "years"
    first <- as.POSIXct(trunc.POSIXt(first, unit))
    lt <- as.POSIXlt(first)
    if (spec$kind == "weeks") {
      first <- as.POSIXct(trunc.POSIXt(first - ((lt$wday - if (start.on.monday) 1L else 0L) %% 7L) * 86400, "days"))
    } else if (spec$kind == "quarters") {
      f <- ..bareLt(lt)
      f$mon <- f$mon - f$mon %% 3L
      f$isdst[] <- -1L
      first <- as.POSIXct(..ltFrom(f, attr(lt, "tzone"), NULL))
    }
    by <- switch(spec$kind, weeks = paste(spec$mult * 7L, "days"), DSTdays = paste(spec$mult, "DSTdays"),
                 paste(spec$mult, spec$kind))
    step <- switch(spec$kind, secs = 1, mins = 60, hours = 3600, days = 86400, weeks = 604800, DSTdays = 86400,
                   months = 28 * 86400, quarters = 84 * 86400, years = 365 * 86400) * spec$mult
    n <- ..checkBreaks(ceiling((as.numeric(unclass(last)) - as.numeric(unclass(first))) / step) + 2)
    pts <- seq.POSIXt(first, by = by, length.out = n)
    keep <- which(pts > last)[1L]
    breaks <- pts[seq_len(keep)]
    fallback <- function(codes) as.character(breaks[-length(breaks)])
    num <- unclass(breaks)
    attributes(num) <- NULL
  } else stop("invalid specification of 'breaks'")
  xv <- unclass(x)
  attributes(xv) <- NULL
  ..cutFactor(xv, num, labels, right, fallback, ...)
}

# ---- summaries of dates and times ------------------------------------------

# the six figures of summary.default over the numbers of a date-time object,
# with the count of missing values kept apart and the class given
..summaryTime <- function(object, digits, cls, ...) {
  s <- summary.default(object, digits = digits, ...)
  v <- unclass(s)
  at <- match("NAs", names(v), 0L)
  nas <- if (at) as.integer(v[at])
  if (at) v <- v[-at]
  attributes(v) <- list(names = names(v))
  class(v) <- cls
  if (at) attr(v, "NAs") <- nas
  v
}

summary.Date <- function(object, digits = 12L, ...)
  ..summaryTime(as.vector(unclass(object)) + 0, digits, c("summaryDefault", "Date"), ...)

summary.POSIXct <- function(object, digits = 15L, ...) {
  v <- ..summaryTime(as.vector(unclass(object)), digits, c("summaryDefault", "POSIXct", "POSIXt"), ...)
  if (!is.null(attr(object, "tzone"))) attr(v, "tzone") <- attr(object, "tzone")
  v
}

summary.POSIXlt <- function(object, digits = 15L, ...) summary.POSIXct(as.POSIXct(object), digits, ...)

summary.difftime <- function(object, ...) {
  v <- ..summaryTime(as.vector(unclass(object)), 7L, c("summaryDefault", "difftime"), ...)
  attr(v, "units") <- attr(object, "units")
  v
}
