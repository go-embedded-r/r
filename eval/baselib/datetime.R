# Base library, part 5: dates and times. Written from the help pages of
# Dates, DateTimeClasses, as.Date, strptime, difftime, seq.Date and
# ISOdatetime in R 4.6.2 and from experiments against the oracle
# (docs/m3-plan.md). BSD licence. The conversions between seconds and
# calendar fields, the parsing and the formatting are Go (datetime.go, the
# .dt* helpers); what is here is the class machinery. Names that start with
# two dots are helpers of ours, not functions of base.

.Date <- function(xx, cl = "Date") `class<-`(xx, cl)

.POSIXct <- function(xx, tz = NULL, cl = c("POSIXct", "POSIXt")) {
  class(xx) <- cl
  attr(xx, "tzone") <- tz
  xx
}

.difftime <- function(xx, units, cl = "difftime") {
  class(xx) <- cl
  attr(xx, "units") <- units
  xx
}

Sys.time <- function() .POSIXct(.dtNow())

Sys.Date <- function() as.Date(as.POSIXlt(Sys.time()))

## ---- helpers --------------------------------------------------------------

..unitSeconds <- c(secs = 1, mins = 60, hours = 3600, days = 86400, weeks = 604800)

# a difftime as plain seconds
..inSeconds <- function(d) as.vector(d) * ..unitSeconds[[attr(d, "units")]]

# a difftime as a whole number of days, for date arithmetic
..inDays <- function(d) as.vector(round(..inSeconds(d) / 86400))

# the tzone attribute of a POSIXlt: the name alone for a UTC zone, else the
# name and the two abbreviations
..ltTzone <- function(tz) {
  if (!.dtIsUTC(tz)) .dtTzone(tz)
  else if (nzchar(tz)) tz
  else "UTC"
}

..mkLt <- function(fields, tzone) {
  class(fields) <- c("POSIXlt", "POSIXt")
  attr(fields, "tzone") <- tzone
  attr(fields, "balanced") <- TRUE
  fields
}

..wholeAsInteger <- function(x) if (all(is.na(x) | x == round(x))) as.integer(x) else x

..comparison <- function(op) op %in% c("<", ">", "==", "!=", "<=", ">=")

# the first of the formats that reads every given string, or NA
..firstFormat <- function(strings, formats, tz) {
  for (f in formats) {
    if (!any(is.na(strptime(strings, f, tz = tz)))) return(f)
  }
  NA_character_
}

## ---- POSIXlt and POSIXct ------------------------------------------------

as.POSIXlt <- function(x, tz = "", ...) UseMethod("as.POSIXlt")

as.POSIXlt.POSIXct <- function(x, tz = "", ...) {
  zone <- tz
  if ((missing(tz) || !nzchar(tz)) && !is.null(attr(x, "tzone"))) zone <- attr(x, "tzone")[1L]
  parts <- .dtToLt(unclass(x), zone)
  if (!is.null(nm <- names(x))) names(parts[[6L]]) <- nm
  ..mkLt(parts, ..ltTzone(zone))
}

as.POSIXlt.Date <- function(x, tz = "UTC", ...)
  ..mkLt(.dtToLt(floor(unclass(x)) * 86400, "UTC"), "UTC")

as.POSIXlt.numeric <- function(x, tz = "", origin, ...)
  as.POSIXlt(as.POSIXct(x, tz = tz, origin = origin), tz = tz)

as.POSIXlt.character <- function(x, tz = "", format,
                                 tryFormats = c("%Y-%m-%d %H:%M:%OS", "%Y/%m/%d %H:%M:%OS",
                                                "%Y-%m-%d %H:%M", "%Y/%m/%d %H:%M",
                                                "%Y-%m-%d", "%Y/%m/%d"),
                                 optional = FALSE, ...) {
  x <- unclass(x)
  retag <- function(lt) {
    if (nzchar(tz)) attr(lt, "tzone") <- tz
    lt
  }
  if (!missing(format)) return(retag(strptime(x, format, tz = tz)))
  present <- x[!is.na(x)]
  if (length(present) == 0L) return(retag(strptime(x, "%Y/%m/%d", tz = tz)))
  chosen <- ..firstFormat(present, tryFormats, tz)
  if (!is.na(chosen)) return(retag(strptime(x, chosen, tz = tz)))
  if (optional) return(as.POSIXlt.character(rep.int(NA_character_, length(x)), tz = tz))
  stop("character string is not in a standard unambiguous format")
}

as.POSIXlt.factor <- function(x, ...) {
  res <- as.POSIXlt(as.character(x), ...)
  attr(res, "balanced") <- NA
  res
}

as.POSIXlt.default <- function(x, tz = "", optional = FALSE, ...) {
  if (inherits(x, "POSIXlt")) {
    if (!missing(tz)) attr(x, "tzone") <- tz
    return(x)
  }
  if (is.null(x)) return(as.POSIXlt.character(character(), "UTC"))
  if (is.logical(x) && all(is.na(x))) return(as.POSIXlt.character(as.character(x), "UTC"))
  if (optional) return(as.POSIXlt.character(rep.int(NA_character_, length(x)), tz = tz))
  stop(gettextf("do not know how to convert '%s' to class %s",
                deparse1(substitute(x)), dQuote("POSIXlt")), domain = NA)
}

as.POSIXct <- function(x, tz = "", ...) UseMethod("as.POSIXct")

as.POSIXct.POSIXlt <- function(x, tz = "", ...) {
  zone <- tz
  if (missing(tz) && !is.null(attr(x, "tzone"))) zone <- attr(x, "tzone")[1L]
  secs <- .dtToCt(x, zone)
  names(secs) <- names(x$year)
  .POSIXct(secs, zone)
}

as.POSIXct.Date <- function(x, tz = "UTC", ...) .POSIXct(unclass(x) * 86400, "UTC")

as.POSIXct.numeric <- function(x, tz = "", origin, ...) {
  if (missing(origin)) return(.POSIXct(x, tz))
  start <- unclass(as.POSIXct(origin, tz = "GMT", ...))
  .POSIXct(start + x, tz)
}

as.POSIXct.default <- function(x, tz = "", ...) {
  if (inherits(x, "POSIXct")) {
    if (missing(tz)) return(x)
    secs <- unclass(x)
    attr(secs, "tzone") <- NULL
    return(.POSIXct(secs, tz))
  }
  if (is.null(x)) return(.POSIXct(integer(), tz))
  if (is.character(x) || is.factor(x)) return(as.POSIXct(as.POSIXlt(x, tz, ...), tz, ...))
  if (is.logical(x) && all(is.na(x))) return(.POSIXct(as.numeric(x), tz))
  stop(gettextf("do not know how to convert '%s' to class %s",
                deparse1(substitute(x)), dQuote("POSIXct")), domain = NA)
}

strptime <- function(x, format, tz = "") {
  parts <- .dtParse(as.character(x), format, tz)
  if (!is.null(nm <- names(x))) names(parts[[6L]]) <- nm
  ..mkLt(parts, ..ltTzone(tz))
}

format.POSIXlt <- function(x, format = "", usetz = FALSE,
                           digits = getOption("digits.secs"), ...) {
  if (!inherits(x, "POSIXlt")) stop("wrong class")
  # decimals shown for seconds: at most `digits`, as few as the data allow
  np <- 0L
  if (!is.null(digits)) {
    secs <- x$sec[is.finite(x$sec)]
    np <- min(6L, digits)
    for (k in seq_len(np) - 1L) {
      if (all(abs(secs - round(secs, k)) < 1e-06)) {
        np <- k
        break
      }
    }
  }
  blank <- format == ""
  if (any(blank)) {
    clock <- unclass(x)[c("hour", "min", "sec")]
    date_only <- all(unlist(lapply(clock, function(v) all(v[is.finite(v)] == 0))))
    format[blank] <- if (date_only) "%Y-%m-%d"
                     else if (np == 0L) "%Y-%m-%d %H:%M:%S"
                     else paste0("%Y-%m-%d %H:%M:%OS", np)
  }
  res <- .dtFormat(x, format, np)
  names(res) <- names(x$year)
  if (usetz) {
    zone <- x$zone
    if (is.null(zone)) zone <- rep("", length(res))
    if (!is.null(attr(x, "tzone")) && !any(nzchar(zone)))
      zone <- rep(attr(x, "tzone")[1L], length(res))
    tagged <- !is.na(res) & nzchar(zone) & is.finite(x$sec)
    res[tagged] <- paste(res[tagged], zone[tagged])
  }
  res
}

format.POSIXct <- function(x, format = "", tz = "", usetz = FALSE, ...) {
  zone <- tz
  if (missing(tz) && !is.null(attr(x, "tzone"))) zone <- attr(x, "tzone")[1L]
  structure(format.POSIXlt(as.POSIXlt(x, zone), format, usetz, ...), names = names(x))
}

strftime <- function(x, format = "", tz = "", usetz = FALSE, ...)
  format(as.POSIXlt(x, tz = tz), format = format, usetz = usetz, ...)

as.character.POSIXt <- function(x, digits = if (inherits(x, "POSIXlt")) 14L else 6L,
                                OutDec = ".", ...)
  format(x, digits = digits)

print.POSIXct <- function(x, tz = "", usetz = TRUE, max = NULL, ...) {
  if (length(x) == 0L) {
    cat("POSIXct of length 0\n")
  } else {
    text <- if (missing(tz)) format(x, usetz = usetz) else format(x, tz = tz, usetz = usetz)
    print(text, max = max, ...)
  }
  invisible(x)
}

print.POSIXlt <- function(x, tz = "", usetz = TRUE, max = NULL, ...) {
  if (length(x$year) == 0L) cat("POSIXlt of length 0\n")
  else print(format(x, usetz = usetz), max = max, ...)
  invisible(x)
}

length.POSIXlt <- function(x) length(unclass(x)[[1L]])

is.na.POSIXlt <- function(x) is.na(as.POSIXct(x))

ISOdatetime <- function(year, month, day, hour, min, sec, tz = "") {
  stamp <- paste(year, month, day, hour, min, sec, sep = "-")
  as.POSIXct(strptime(stamp, "%Y-%m-%d-%H-%M-%OS", tz = tz), tz = tz)
}

ISOdate <- function(year, month, day, hour = 12, min = 0, sec = 0, tz = "GMT")
  ISOdatetime(year, month, day, hour, min, sec, tz)

## ---- Date ------------------------------------------------------------------

as.Date <- function(x, ...) UseMethod("as.Date")

as.Date.default <- function(x, ...) {
  if (inherits(x, "Date")) return(x)
  if (is.null(x)) return(.Date(numeric()))
  if (is.logical(x) && all(is.na(x))) return(.Date(as.numeric(x)))
  stop(gettextf("do not know how to convert '%s' to class %s",
                deparse1(substitute(x)), dQuote("Date")), domain = NA)
}

as.Date.numeric <- function(x, origin, ...)
  if (missing(origin)) .Date(x) else as.Date(origin, ...) + x

as.Date.factor <- function(x, ...) as.Date(as.character(x), ...)

as.Date.character <- function(x, format, tryFormats = c("%Y-%m-%d", "%Y/%m/%d"),
                              optional = FALSE, ...) {
  if (!missing(format)) return(as.Date(strptime(x, format, tz = "GMT")))
  # the format is found on the first string that is not NA
  filled <- which(!is.na(x))
  if (length(filled) == 0L) return(as.Date(strptime(x, "%Y-%m-%d")))
  probe <- x[filled[1L]]
  chosen <- NA_character_
  for (f in tryFormats) {
    if (!is.na(strptime(probe, f, tz = "GMT"))) {
      chosen <- f
      break
    }
  }
  if (is.na(chosen)) {
    if (optional) return(as.Date.character(rep.int(NA_character_, length(x))))
    charToDate <- function(x) stop("character string is not in a standard unambiguous format")
    return(charToDate(x))
  }
  as.Date(strptime(x, chosen))
}

as.Date.POSIXct <- function(x, tz = "UTC", ...) {
  if (tz != "UTC") return(as.Date(as.POSIXlt(x, tz = tz)))
  days <- floor(unclass(x) / 86400)
  attr(days, "tzone") <- NULL
  .Date(days)
}

as.Date.POSIXlt <- function(x, ...) {
  days <- floor(.dtToCt(x, "UTC") / 86400)
  names(days) <- names(x$year)
  .Date(days)
}

format.Date <- function(x, ...) {
  text <- format(as.POSIXlt(x), ...)
  names(text) <- names(x)
  text
}

as.character.Date <- function(x, ...) format(x, ...)

print.Date <- function(x, max = NULL, ...) {
  if (length(x) == 0L) cat("Date of length 0\n")
  else print(format(x), max = max, ...)
  invisible(x)
}

is.numeric.Date <- function(x) FALSE
is.numeric.POSIXt <- function(x) FALSE
is.numeric.difftime <- function(x) FALSE

weekdays <- function(x, abbreviate) UseMethod("weekdays")
weekdays.Date <- function(x, abbreviate = FALSE) format(x, if (abbreviate) "%a" else "%A")
weekdays.POSIXt <- function(x, abbreviate = FALSE) format(x, if (abbreviate) "%a" else "%A")
months <- function(x, abbreviate) UseMethod("months")
months.Date <- function(x, abbreviate = FALSE) format(x, if (abbreviate) "%b" else "%B")
months.POSIXt <- function(x, abbreviate = FALSE) format(x, if (abbreviate) "%b" else "%B")
quarters <- function(x, ...) UseMethod("quarters")
quarters.Date <- function(x, ...) paste0("Q", as.POSIXlt(x)$mon %/% 3L + 1L)
quarters.POSIXt <- function(x, ...) paste0("Q", as.POSIXlt(x)$mon %/% 3L + 1L)

julian <- function(x, ...) UseMethod("julian")
julian.Date <- function(x, origin = as.Date("1970-01-01"), ...) {
  elapsed <- difftime(x, origin, units = "days")
  structure(unclass(elapsed)[seq_along(elapsed)], origin = origin)
}

`[.Date` <- function(x, ..., drop = TRUE) .Date(unclass(x)[...], oldClass(x))

`[[.Date` <- function(x, ..., drop = TRUE) .Date(unclass(x)[[...]], oldClass(x))

`[<-.Date` <- function(x, ..., value) {
  if (!length(value)) return(x)
  cls <- oldClass(x)
  bare <- unclass(x)
  bare[...] <- unclass(as.Date(value))
  .Date(bare, cls)
}

`[.POSIXct` <- function(x, ..., drop = TRUE)
  .POSIXct(unclass(x)[...], attr(x, "tzone"), oldClass(x))

`[[.POSIXct` <- function(x, ..., drop = TRUE)
  .POSIXct(unclass(x)[[...]], attr(x, "tzone"), oldClass(x))

`[.difftime` <- function(x, ..., drop = TRUE)
  .difftime(unclass(x)[...], attr(x, "units"), oldClass(x))

c.Date <- function(..., recursive = FALSE)
  .Date(c(unlist(lapply(list(...), function(e) unclass(as.Date(e))))))

c.POSIXct <- function(..., recursive = FALSE) {
  parts <- list(...)
  secs <- lapply(parts, function(e) unclass(as.POSIXct(e)))
  zones <- lapply(parts, attr, "tzone")
  same <- length(unique(zones)) == 1L
  .POSIXct(c(unlist(secs)), if (same) zones[[1L]] else NULL)
}

c.difftime <- function(..., recursive = FALSE) {
  parts <- list(...)
  if (length(parts) == 0L) return(.difftime(double(), "secs"))
  is_dt <- vapply(parts, inherits, NA, "difftime")
  units <- vapply(parts[is_dt], attr, "", "units")
  common <- if (all(units == units[1L])) units[1L] else "secs"
  parts[!is_dt] <- lapply(parts[!is_dt], as.difftime, units = common)
  if (common == "secs" && !all(units == "secs"))
    parts[is_dt] <- lapply(parts[is_dt], ..inSeconds)
  .difftime(unlist(parts), common)
}

rep.Date <- function(x, ...) {
  y <- NextMethod()
  .Date(y, oldClass(x))
}

rep.POSIXct <- function(x, ...) {
  y <- NextMethod()
  .POSIXct(y, attr(x, "tzone"), oldClass(x))
}

rep.difftime <- function(x, ...)
  .difftime(rep(unclass(x), ...), attr(x, "units"), oldClass(x))

mean.Date <- function(x, ...) .Date(mean(unclass(x), ...))

mean.POSIXct <- function(x, ...) .POSIXct(mean(unclass(x), ...), attr(x, "tzone"))

Summary.Date <- function(..., na.rm = FALSE) {
  if (!(.Generic %in% c("max", "min", "range")))
    stop(gettextf("%s not defined for \"Date\" objects", .Generic), domain = NA)
  val <- NextMethod(.Generic)
  class(val) <- oldClass(list(...)[[1L]])
  val
}

Summary.POSIXct <- function(..., na.rm = FALSE) {
  if (!(.Generic %in% c("max", "min", "range")))
    stop(gettextf("'%s' not defined for \"POSIXt\" objects", .Generic), domain = NA)
  zone <- attr(list(...)[[1L]], "tzone")
  .POSIXct(NextMethod(.Generic), zone)
}

# + and - on dates: the other operand is a number of days or a time difference
`+.Date` <- function(e1, e2) {
  if (nargs() == 1) return(e1)
  if (inherits(e1, "Date") && inherits(e2, "Date"))
    stop("binary + is not defined for \"Date\" objects")
  if (inherits(e1, "difftime")) e1 <- ..inDays(e1)
  if (inherits(e2, "difftime")) e2 <- ..inDays(e2)
  .Date(unclass(e1) + unclass(e2))
}

`-.Date` <- function(e1, e2) {
  if (!inherits(e1, "Date")) stop("can only subtract from \"Date\" objects")
  if (nargs() == 1) stop("unary - is not defined for \"Date\" objects")
  if (inherits(e2, "Date")) return(difftime(e1, e2, units = "days"))
  if (inherits(e2, "difftime")) e2 <- ..inDays(e2)
  if (!is.null(attr(e2, "class"))) stop("can only subtract numbers from \"Date\" objects")
  .Date(unclass(as.Date(e1)) - e2)
}

Ops.Date <- function(e1, e2) {
  op <- .Generic
  if (nargs() == 1) stop(gettextf("unary %s not defined for \"Date\" objects", op), domain = NA)
  if (!..comparison(op))
    stop(gettextf("%s not defined for \"Date\" objects", op), domain = NA)
  if (is.character(e1)) e1 <- as.Date(e1)
  if (is.character(e2)) e2 <- as.Date(e2)
  NextMethod(.Generic)
}

`+.POSIXt` <- function(e1, e2) {
  if (nargs() == 1) return(e1)
  if (inherits(e1, "POSIXt") && inherits(e2, "POSIXt"))
    stop("binary '+' is not defined for \"POSIXt\" objects")
  if (inherits(e1, "POSIXlt")) e1 <- as.POSIXct(e1)
  if (inherits(e2, "POSIXlt")) e2 <- as.POSIXct(e2)
  if (inherits(e1, "difftime")) e1 <- ..inSeconds(e1)
  if (inherits(e2, "difftime")) e2 <- ..inSeconds(e2)
  zone <- attr(if (inherits(e1, "POSIXt")) e1 else e2, "tzone")
  .POSIXct(unclass(e1) + unclass(e2), zone)
}

`-.POSIXt` <- function(e1, e2) {
  if (!inherits(e1, "POSIXt")) stop("can only subtract from \"POSIXt\" objects")
  if (nargs() == 1) stop("unary '-' is not defined for \"POSIXt\" objects")
  if (inherits(e2, "POSIXt")) return(difftime(e1, e2))
  if (inherits(e2, "difftime")) e2 <- ..inSeconds(e2)
  if (!is.null(attr(e2, "class"))) stop("can only subtract numbers from \"POSIXt\" objects")
  e1 <- as.POSIXct(e1)
  .POSIXct(unclass(e1) - e2, attr(e1, "tzone"))
}

Ops.POSIXt <- function(e1, e2) {
  op <- .Generic
  if (nargs() == 1) stop(gettextf("unary '%s' not defined for \"POSIXt\" objects", op), domain = NA)
  if (!..comparison(op))
    stop(gettextf("'%s' not defined for \"POSIXt\" objects", op), domain = NA)
  zoneOf <- function(other) if (is.null(z <- attr(other, "tzone"))) "" else z[1L]
  if (is.character(e1)) e1 <- as.POSIXct(e1, tz = zoneOf(e2))
  else if (inherits(e1, "POSIXlt")) e1 <- as.POSIXct(e1)
  if (is.character(e2)) e2 <- as.POSIXct(e2, tz = zoneOf(e1))
  else if (inherits(e2, "POSIXlt")) e2 <- as.POSIXct(e2)
  zones <- unique(c(zoneOf(e1), zoneOf(e2)))
  if (length(zones[nzchar(zones)]) > 1L) warning("'tzone' attributes are inconsistent")
  NextMethod(.Generic)
}

## ---- difftime --------------------------------------------------------------

difftime <- function(time1, time2, tz, units = c("auto", "secs", "mins", "hours", "days", "weeks")) {
  if (missing(tz)) {
    t1 <- as.POSIXct(time1)
    t2 <- as.POSIXct(time2)
  } else {
    t1 <- as.POSIXct(time1, tz = tz)
    t2 <- as.POSIXct(time2, tz = tz)
  }
  secs <- unclass(t1) - unclass(t2)
  attr(secs, "tzone") <- NULL
  units <- match.arg(units)
  if (units == "auto") {
    # the largest unit (weeks excepted) in which the smallest gap is at least one
    smallest <- if (all(is.na(secs))) 0 else min(abs(secs), na.rm = TRUE)
    units <- if (!is.finite(smallest) || smallest < 60) "secs"
             else if (smallest < 3600) "mins"
             else if (smallest < 86400) "hours"
             else "days"
  }
  .difftime(secs / ..unitSeconds[[units]], units = units)
}

as.difftime <- function(tim, format = "%X", units = "auto", tz = "UTC") {
  if (inherits(tim, "difftime")) return(tim)
  if (is.character(tim)) {
    midnight <- strptime("0:0:0", format = "%X", tz = tz)
    return(difftime(strptime(tim, format = format, tz = tz), midnight, units = units))
  }
  if (!is.numeric(tim)) stop("'tim' is not character or numeric")
  if (units == "auto") stop("need explicit units for numeric conversion")
  if (!(units %in% names(..unitSeconds))) stop("invalid units specified")
  .difftime(tim, units = units)
}

units <- function(x) UseMethod("units")

units.difftime <- function(x) attr(x, "units")

`units<-` <- function(x, value) UseMethod("units<-")

`units<-.difftime` <- function(x, value) {
  from <- units(x)
  if (from == value) return(x)
  if (!(value %in% names(..unitSeconds))) stop("invalid units specified")
  .difftime(unclass(x) * (..unitSeconds[[from]] / ..unitSeconds[[value]]), value)
}

as.double.difftime <- function(x, units = "auto", ...) {
  if (units != "auto") units(x) <- units
  as.vector(unclass(x), "double")
}

format.difftime <- function(x, ...) paste(format(unclass(x), ...), units(x))

print.difftime <- function(x, digits = getOption("digits"), ...) {
  if (length(x) == 0L) {
    cat("difftime of length 0\n")
  } else if (is.array(x) || length(x) > 1L) {
    cat("Time differences in ", attr(x, "units"), "\n", sep = "")
    plain <- unclass(x)
    attr(plain, "units") <- NULL
    print(plain, digits = digits, ...)
  } else {
    cat("Time difference of ", format(unclass(x), digits = digits), " ",
        attr(x, "units"), "\n", sep = "")
  }
  invisible(x)
}

`*.difftime` <- function(e1, e2) {
  if (inherits(e1, "difftime") && inherits(e2, "difftime"))
    stop("both arguments of * cannot be \"difftime\" objects")
  if (inherits(e2, "difftime")) {
    held <- e1
    e1 <- e2
    e2 <- held
  }
  .difftime(e2 * unclass(e1), attr(e1, "units"))
}

`/.difftime` <- function(e1, e2) {
  if (inherits(e2, "difftime")) stop("second argument of / cannot be a \"difftime\" object")
  .difftime(unclass(e1) / e2, attr(e1, "units"))
}

Ops.difftime <- function(e1, e2) {
  op <- .Generic
  if (nargs() == 1) {
    if (op == "+") return(e1)
    if (op == "-") {
      e1[] <- -unclass(e1)
      return(e1)
    }
    stop(gettextf("unary '%s' not defined for \"difftime\" objects", op), domain = NA)
  }
  if (..comparison(op)) {
    if (inherits(e1, "difftime") && inherits(e2, "difftime")) {
      e1 <- ..inSeconds(e1)
      e2 <- ..inSeconds(e2)
    }
    return(NextMethod(.Generic))
  }
  if (!(op %in% c("+", "-")))
    stop(gettextf("'%s' not defined for \"difftime\" objects", op), domain = NA)
  d1 <- inherits(e1, "difftime")
  d2 <- inherits(e2, "difftime")
  if (d1 && !d2) return(.difftime(NextMethod(.Generic), units = attr(e1, "units")))
  if (d2 && !d1) return(.difftime(NextMethod(.Generic), units = attr(e2, "units")))
  if (attr(e1, "units") == attr(e2, "units"))
    return(.difftime(NextMethod(.Generic), units = attr(e1, "units")))
  # the attributes of the result are laid out as R's arithmetic lays them out:
  # those of the second operand first, and an operand already in seconds is
  # stored class first
  secs_second <- attr(e2, "units") == "secs"
  e1 <- ..inSeconds(e1)
  e2 <- ..inSeconds(e2)
  res <- as.vector(NextMethod(.Generic))
  if (secs_second) {
    class(res) <- "difftime"
    attr(res, "units") <- "secs"
  } else {
    attr(res, "units") <- "secs"
    class(res) <- "difftime"
  }
  res
}

Math.difftime <- function(x, ...) {
  if (!(.Generic %in% c("abs", "sign", "floor", "ceiling", "trunc", "round", "signif")))
    stop(gettextf("'%s' not defined for \"difftime\" objects", .Generic), domain = NA)
  .difftime(NextMethod(), attr(x, "units"))
}

Summary.difftime <- function(..., na.rm) {
  if (!(.Generic %in% c("max", "min", "sum", "range")))
    stop(gettextf("'%s' not defined for \"difftime\" objects", .Generic), domain = NA)
  parts <- list(...)
  units <- vapply(parts, attr, "", "units")
  if (all(units == units[1L])) {
    out <- do.call(.Generic, c(lapply(parts, as.vector), na.rm = na.rm))
    .difftime(out, units[1L])
  } else {
    out <- do.call(.Generic, c(lapply(parts, ..inSeconds), na.rm = na.rm))
    .difftime(out, "secs")
  }
}

mean.difftime <- function(x, ...) .difftime(mean(unclass(x), ...), attr(x, "units"))

as.list.Date <- function(x, ...) lapply(unclass(x), .Date, oldClass(x))

as.list.POSIXct <- function(x, ...) {
  nms <- names(x)
  zone <- attr(x, "tzone")
  cls <- oldClass(x)
  out <- lapply(unname(unclass(x)), .POSIXct, zone, cls)
  names(out) <- nms
  out
}

as.list.difftime <- function(x, ...) lapply(unclass(x), .difftime, attr(x, "units"), oldClass(x))

## ---- seq.Date --------------------------------------------------------------

seq.Date <- function(from, to, by, length.out = NULL, along.with = NULL, ...) {
  if (missing(from)) stop("'from' must be specified")
  if (!inherits(from, "Date")) stop("'from' must be a \"Date\" object")
  if (length(as.Date(from)) != 1L) stop("'from' must be of length 1")
  if (!missing(to)) {
    if (!inherits(to, "Date")) stop("'to' must be a \"Date\" object")
    if (length(as.Date(to)) != 1L) stop("'to' must be of length 1")
  }
  if (!missing(along.with)) {
    length.out <- length(along.with)
  } else if (!is.null(length.out)) {
    if (length(length.out) != 1L) stop("'length.out' must be of length 1")
    length.out <- ceiling(length.out)
  }
  start <- unclass(as.Date(from))
  if (missing(by)) {
    if (missing(to) && is.null(length.out))
      stop("without 'by', when one of 'to', 'from' is missing, 'length.out' / 'along.with' must be specified")
    end <- if (missing(to)) NULL else unclass(as.Date(to))
    days <- if (is.null(length.out)) seq.int(start, end) else seq.int(start, end, length.out = length.out)
    return(.Date(..wholeAsInteger(days)))
  }
  if (missing(to) + is.null(length.out) != 1L)
    stop("given 'by', exactly two of 'to', 'from' and 'length.out' / 'along.with' must be specified")
  if (length(by) != 1L) stop("'by' must be of length 1")
  unit <- 0L # 0: a number of days; else the position in the list below
  if (inherits(by, "difftime")) {
    by <- unclass(by) * c(secs = 1/86400, mins = 1/1440, hours = 1/24, days = 1, weeks = 7)[[attr(by, "units")]]
  } else if (is.character(by)) {
    words <- strsplit(by, " ", fixed = TRUE)[[1L]]
    if (length(words) < 1L || length(words) > 2L) stop("invalid 'by' string")
    unit <- pmatch(words[length(words)], c("days", "weeks", "months", "quarters", "years", "DSTdays"))
    if (is.na(unit)) stop("invalid string for 'by'")
    count <- if (length(words) == 2L) as.integer(words[1L]) else 1
    by <- if (unit <= 2L) count * c(1, 7)[unit] else count
  } else if (!is.numeric(by)) {
    stop("invalid mode for 'by'")
  }
  if (is.na(by)) stop("'by' is NA")
  if (unit <= 2L || unit == 6L) {
    if (!is.null(length.out)) {
      days <- ..wholeAsInteger(seq.int(start, by = by, length.out = length.out))
    } else {
      days <- seq.int(0, unclass(as.Date(to)) - start, by) + start
    }
    res <- .Date(days)
  } else {
    lt <- as.POSIXlt(from)
    if (unit == 5L) {
      lt$year <- if (missing(to)) seq.int(lt$year, by = by, length.out = length.out)
                 else seq.int(lt$year, as.POSIXlt(to)$year, by)
    } else {
      if (unit == 4L) by <- by * 3
      lt$mon <- if (missing(to)) seq.int(lt$mon, by = by, length.out = length.out)
                else seq.int(lt$mon, 12 * (as.POSIXlt(to)$year - lt$year) + as.POSIXlt(to)$mon, by)
    }
    res <- as.Date(lt)
  }
  if (!missing(to)) {
    limit <- as.Date(to)
    res <- if (by > 0) res[res <= limit] else res[res >= limit]
  }
  res
}
