# Base library, part 5: dates and times. Written from the help pages of
# Dates, DateTimeClasses, as.Date, strptime, difftime, seq.Date and
# ISOdatetime in R 4.6.2 and checked against the oracle (docs/m3-plan.md).
# BSD licence. The conversions between seconds and calendar fields, the
# parsing and the formatting are Go (datetime.go, the .dt* helpers); what is
# here is the class machinery.

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

## ---- POSIXlt and POSIXct ------------------------------------------------

..ltTzone <- function(tz) {
  if (identical(tz, "UTC") || identical(tz, "GMT")) tz else .dtTzone(tz)
}

..mkLt <- function(fields, tzone) {
  class(fields) <- c("POSIXlt", "POSIXt")
  attr(fields, "tzone") <- tzone
  attr(fields, "balanced") <- TRUE
  fields
}

as.POSIXlt <- function(x, tz = "", ...) UseMethod("as.POSIXlt")

as.POSIXlt.POSIXct <- function(x, tz = "", ...) {
  if (missing(tz) && !is.null(tzone <- attr(x, "tzone"))) tz <- tzone[1L]
  fields <- .dtToLt(unclass(x), tz)
  if (!is.null(nm <- names(x))) names(fields[[6L]]) <- nm
  ..mkLt(fields, ..ltTzone(tz))
}

as.POSIXlt.Date <- function(x, tz = "UTC", ...) {
  ..mkLt(.dtToLt(floor(unclass(x)) * 86400, "UTC"), "UTC")
}

as.POSIXlt.numeric <- function(x, tz = "", origin, ...) {
  as.POSIXlt(as.POSIXct(x, tz = tz, origin = origin), tz = tz)
}

as.POSIXlt.character <- function(x, tz = "", format,
                                 tryFormats = c("%Y-%m-%d %H:%M:%OS", "%Y/%m/%d %H:%M:%OS",
                                                "%Y-%m-%d %H:%M", "%Y/%m/%d %H:%M",
                                                "%Y-%m-%d", "%Y/%m/%d"),
                                 optional = FALSE, ...) {
  x <- unclass(x)
  if (!missing(format)) {
    res <- strptime(x, format, tz = tz)
    if (nzchar(tz)) attr(res, "tzone") <- tz
    return(res)
  }
  xx <- x[!is.na(x)]
  if (!length(xx)) {
    res <- strptime(x, "%Y/%m/%d", tz = tz)
    if (nzchar(tz)) attr(res, "tzone") <- tz
    return(res)
  }
  for (ff in tryFormats) {
    if (all(!is.na(strptime(xx, ff, tz = tz)))) {
      res <- strptime(x, ff, tz = tz)
      if (nzchar(tz)) attr(res, "tzone") <- tz
      return(res)
    }
  }
  if (optional) as.POSIXlt.character(rep.int(NA_character_, length(x)), tz = tz)
  else stop("character string is not in a standard unambiguous format")
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
  if (optional) as.POSIXlt.character(rep.int(NA_character_, length(x)), tz = tz)
  else stop(gettextf("do not know how to convert '%s' to class %s",
                     deparse1(substitute(x)), dQuote("POSIXlt")), domain = NA)
}

as.POSIXct <- function(x, tz = "", ...) UseMethod("as.POSIXct")

as.POSIXct.POSIXlt <- function(x, tz = "", ...) {
  tzone <- attr(x, "tzone")
  if (missing(tz) && !is.null(tzone)) tz <- tzone[1L]
  secs <- .dtToCt(x, tz)
  names(secs) <- names(x$year)
  .POSIXct(secs, tz)
}

as.POSIXct.Date <- function(x, tz = "UTC", ...) .POSIXct(unclass(x) * 86400, "UTC")

as.POSIXct.numeric <- function(x, tz = "", origin, ...) {
  if (missing(origin)) .POSIXct(x, tz)
  else .POSIXct(unclass(as.POSIXct(origin, tz = "GMT", ...)) + x, tz)
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
  res <- .dtParse(as.character(x), format, tz)
  if (!is.null(nm <- names(x))) names(res[[6L]]) <- nm
  ..mkLt(res, ..ltTzone(tz))
}

format.POSIXlt <- function(x, format = "", usetz = FALSE,
                           digits = getOption("digits.secs"), ...) {
  if (!inherits(x, "POSIXlt")) stop("wrong class")
  np <- 0L
  if (!is.null(digits)) {
    secs <- x$sec[is.finite(x$sec)]
    np <- min(6L, digits)
    for (i in seq_len(np) - 1L)
      if (all(abs(secs - round(secs, i)) < 1e-06)) {
        np <- i
        break
      }
  }
  if (any(f0 <- format == "")) {
    midnight <- all(unlist(lapply(unclass(x)[c("hour", "min", "sec")],
                                  function(v) all(v[is.finite(v)] == 0))))
    format[f0] <-
      if (midnight) "%Y-%m-%d"
      else if (np == 0L) "%Y-%m-%d %H:%M:%S"
      else paste0("%Y-%m-%d %H:%M:%OS", np)
  }
  res <- .dtFormat(x, format, np)
  names(res) <- names(x$year)
  if (usetz) {
    zone <- if (!is.null(x$zone)) x$zone else rep("", length(res))
    tzn <- attr(x, "tzone")
    if (!is.null(tzn) && all(!nzchar(zone))) zone <- rep(tzn[1L], length(res))
    ok <- !is.na(res) & nzchar(zone)
    res[ok] <- paste(res[ok], zone[ok])
  }
  res
}

format.POSIXct <- function(x, format = "", tz = "", usetz = FALSE, ...) {
  if (missing(tz) && !is.null(tzone <- attr(x, "tzone"))) tz <- tzone[1L]
  structure(format.POSIXlt(as.POSIXlt(x, tz), format, usetz, ...), names = names(x))
}

strftime <- function(x, format = "", tz = "", usetz = FALSE, ...)
  format(as.POSIXlt(x, tz = tz), format = format, usetz = usetz, ...)

as.character.POSIXt <- function(x, digits = if (inherits(x, "POSIXlt")) 14L else 6L,
                                OutDec = ".", ...) {
  format(x, digits = digits)
}

print.POSIXct <- function(x, tz = "", usetz = TRUE, max = NULL, ...) {
  if (length(x) == 0L) cat("POSIXct of length 0\n")
  else {
    shown <- if (missing(tz)) format(x, usetz = usetz) else format(x, tz = tz, usetz = usetz)
    print(shown, max = max, ...)
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
  x <- paste(year, month, day, hour, min, sec, sep = "-")
  as.POSIXct(strptime(x, "%Y-%m-%d-%H-%M-%OS", tz = tz), tz = tz)
}

ISOdate <- function(year, month, day, hour = 12, min = 0, sec = 0, tz = "GMT")
  ISOdatetime(year, month, day, hour, min, sec, tz)

## ---- Date ------------------------------------------------------------------

as.Date <- function(x, ...) UseMethod("as.Date")

as.Date.default <- function(x, ...) {
  if (inherits(x, "Date")) x
  else if (is.null(x)) .Date(numeric())
  else if (is.logical(x) && all(is.na(x))) .Date(as.numeric(x))
  else stop(gettextf("do not know how to convert '%s' to class %s",
                     deparse1(substitute(x)), dQuote("Date")), domain = NA)
}

as.Date.numeric <- function(x, origin, ...) {
  if (missing(origin)) .Date(x) else as.Date(origin, ...) + x
}

as.Date.factor <- function(x, ...) as.Date(as.character(x), ...)

as.Date.character <- function(x, format, tryFormats = c("%Y-%m-%d", "%Y/%m/%d"),
                              optional = FALSE, ...) {
  charToDate <- function(x) {
    xx <- x[1L]
    if (is.na(xx)) {
      j <- 1L
      while (is.na(xx) && (j <- j + 1L) <= length(x)) xx <- x[j]
      if (is.na(xx)) f <- "%Y-%m-%d"
    }
    if (is.na(xx)) strptime(x, f)
    else {
      for (ff in tryFormats)
        if (!is.na(strptime(xx, ff, tz = "GMT"))) return(strptime(x, ff))
      if (optional) as.Date.character(rep.int(NA_character_, length(x)), tz = tz)
      else stop("character string is not in a standard unambiguous format")
    }
  }
  res <- if (missing(format)) charToDate(x) else strptime(x, format, tz = "GMT")
  as.Date(res)
}

as.Date.POSIXct <- function(x, tz = "UTC", ...) {
  if (tz == "UTC") {
    z <- floor(unclass(x) / 86400)
    attr(z, "tzone") <- NULL
    .Date(z)
  } else as.Date(as.POSIXlt(x, tz = tz))
}

as.Date.POSIXlt <- function(x, ...) {
  days <- floor(.dtToCt(x, "UTC") / 86400)
  names(days) <- names(x$year)
  .Date(days)
}

format.Date <- function(x, ...) {
  xx <- format(as.POSIXlt(x), ...)
  names(xx) <- names(x)
  xx
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
weekdays.Date <- function(x, abbreviate = FALSE) format(x, ifelse(abbreviate, "%a", "%A"))
weekdays.POSIXt <- function(x, abbreviate = FALSE) format(x, ifelse(abbreviate, "%a", "%A"))
months <- function(x, abbreviate) UseMethod("months")
months.Date <- function(x, abbreviate = FALSE) format(x, ifelse(abbreviate, "%b", "%B"))
months.POSIXt <- function(x, abbreviate = FALSE) format(x, ifelse(abbreviate, "%b", "%B"))
quarters <- function(x, ...) UseMethod("quarters")
quarters.Date <- function(x, ...) paste0("Q", as.POSIXlt(x)$mon %/% 3L + 1L)
quarters.POSIXt <- function(x, ...) paste0("Q", as.POSIXlt(x)$mon %/% 3L + 1L)

julian <- function(x, ...) UseMethod("julian")
julian.Date <- function(x, origin = as.Date("1970-01-01"), ...) {
  res <- difftime(x, origin, units = "days")
  structure(unclass(res)[seq_along(res)], origin = origin)
}

`[.Date` <- function(x, ..., drop = TRUE) {
  cl <- oldClass(x)
  .Date(unclass(x)[...], cl)
}
`[[.Date` <- function(x, ..., drop = TRUE) {
  cl <- oldClass(x)
  .Date(unclass(x)[[...]], cl)
}
`[<-.Date` <- function(x, ..., value) {
  if (!length(value)) return(x)
  value <- unclass(as.Date(value))
  cl <- oldClass(x)
  x <- unclass(x)
  x[...] <- value
  .Date(x, cl)
}
`[.POSIXct` <- function(x, ..., drop = TRUE) {
  cl <- oldClass(x)
  tz <- attr(x, "tzone")
  .POSIXct(unclass(x)[...], tz, cl)
}
`[[.POSIXct` <- function(x, ..., drop = TRUE) {
  cl <- oldClass(x)
  tz <- attr(x, "tzone")
  .POSIXct(unclass(x)[[...]], tz, cl)
}
`[.difftime` <- function(x, ..., drop = TRUE) {
  cl <- oldClass(x)
  u <- attr(x, "units")
  .difftime(unclass(x)[...], u, cl)
}

c.Date <- function(..., recursive = FALSE)
  .Date(c(unlist(lapply(list(...), function(e) unclass(as.Date(e))))))

c.POSIXct <- function(..., recursive = FALSE) {
  x <- lapply(list(...), function(e) unclass(as.POSIXct(e)))
  tzones <- lapply(list(...), attr, "tzone")
  tz <- if (length(unique(tzones)) == 1L) tzones[[1L]] else NULL
  .POSIXct(c(unlist(x)), tz)
}

c.difftime <- function(..., recursive = FALSE) {
  args <- list(...)
  if (!length(args)) return(.difftime(double(), "secs"))
  ind <- sapply(args, inherits, "difftime")
  pos <- which(!ind)
  units <- sapply(args[ind], attr, "units")
  if (all(units == (un1 <- units[1L]))) {
    if (length(pos)) args[pos] <- lapply(args[pos], as.difftime, units = un1)
    .difftime(unlist(args), un1)
  } else {
    if (length(pos)) args[pos] <- lapply(args[pos], as.difftime, units = "secs")
    args[ind] <- lapply(args[ind], function(x) as.vector(..coerceTimeUnit(x)))
    .difftime(unlist(args), "secs")
  }
}

rep.Date <- function(x, ...) {
  y <- NextMethod()
  .Date(y, oldClass(x))
}
rep.POSIXct <- function(x, ...) {
  y <- NextMethod()
  .POSIXct(y, attr(x, "tzone"), oldClass(x))
}
rep.difftime <- function(x, ...) .difftime(rep(unclass(x), ...), attr(x, "units"), oldClass(x))

mean.Date <- function(x, ...) .Date(mean(unclass(x), ...))
mean.POSIXct <- function(x, ...) .POSIXct(mean(unclass(x), ...), attr(x, "tzone"))

Summary.Date <- function(..., na.rm = FALSE) {
  ok <- switch(.Generic, max = , min = , range = TRUE, FALSE)
  if (!ok) stop(gettextf("%s not defined for \"Date\" objects", .Generic), domain = NA)
  val <- NextMethod(.Generic)
  class(val) <- oldClass(list(...)[[1L]])
  val
}

Summary.POSIXct <- function(..., na.rm = FALSE) {
  ok <- switch(.Generic, max = , min = , range = TRUE, FALSE)
  if (!ok) stop(gettextf("'%s' not defined for \"POSIXt\" objects", .Generic), domain = NA)
  args <- list(...)
  tz <- attr(args[[1L]], "tzone")
  .POSIXct(NextMethod(.Generic), tz)
}

`+.Date` <- function(e1, e2) {
  coerceTimeUnit <- function(x) as.vector(round(switch(attr(x, "units"),
    secs = x/86400, mins = x/1440, hours = x/24, days = x, weeks = 7 * x)))
  if (nargs() == 1) return(e1)
  if (inherits(e1, "Date") && inherits(e2, "Date"))
    stop("binary + is not defined for \"Date\" objects")
  if (inherits(e1, "difftime")) e1 <- coerceTimeUnit(e1)
  if (inherits(e2, "difftime")) e2 <- coerceTimeUnit(e2)
  .Date(unclass(e1) + unclass(e2))
}

`-.Date` <- function(e1, e2) {
  coerceTimeUnit <- function(x) as.vector(round(switch(attr(x, "units"),
    secs = x/86400, mins = x/1440, hours = x/24, days = x, weeks = 7 * x)))
  if (!inherits(e1, "Date")) stop("can only subtract from \"Date\" objects")
  if (nargs() == 1) stop("unary - is not defined for \"Date\" objects")
  if (inherits(e2, "Date")) return(difftime(e1, e2, units = "days"))
  if (inherits(e2, "difftime")) e2 <- coerceTimeUnit(e2)
  if (!is.null(attr(e2, "class"))) stop("can only subtract numbers from \"Date\" objects")
  .Date(unclass(as.Date(e1)) - e2)
}

Ops.Date <- function(e1, e2) {
  if (nargs() == 1) stop(gettextf("unary %s not defined for \"Date\" objects", .Generic), domain = NA)
  boolean <- switch(.Generic, `<` = , `>` = , `==` = , `!=` = , `<=` = , `>=` = TRUE, FALSE)
  if (!boolean) stop(gettextf("%s not defined for \"Date\" objects", .Generic), domain = NA)
  if (is.character(e1)) e1 <- as.Date(e1)
  if (is.character(e2)) e2 <- as.Date(e2)
  NextMethod(.Generic)
}

`+.POSIXt` <- function(e1, e2) {
  coerceTimeUnit <- function(x) as.vector(switch(attr(x, "units"),
    secs = x, mins = 60 * x, hours = 3600 * x, days = 86400 * x, weeks = 7 * 86400 * x))
  if (nargs() == 1) return(e1)
  if (inherits(e1, "POSIXt") && inherits(e2, "POSIXt"))
    stop("binary '+' is not defined for \"POSIXt\" objects")
  if (inherits(e1, "POSIXlt")) e1 <- as.POSIXct(e1)
  if (inherits(e2, "POSIXlt")) e2 <- as.POSIXct(e2)
  if (inherits(e1, "difftime")) e1 <- coerceTimeUnit(e1)
  if (inherits(e2, "difftime")) e2 <- coerceTimeUnit(e2)
  .POSIXct(unclass(e1) + unclass(e2), ..checkTz(e1, e2))
}

..checkTz <- function(e1, e2) {
  tz <- if (inherits(e1, "POSIXt")) attr(e1, "tzone") else attr(e2, "tzone")
  if (is.null(tz)) NULL else tz
}

`-.POSIXt` <- function(e1, e2) {
  coerceTimeUnit <- function(x) as.vector(switch(attr(x, "units"),
    secs = x, mins = 60 * x, hours = 3600 * x, days = 86400 * x, weeks = 7 * 86400 * x))
  if (!inherits(e1, "POSIXt")) stop("can only subtract from \"POSIXt\" objects")
  if (nargs() == 1) stop("unary '-' is not defined for \"POSIXt\" objects")
  if (inherits(e2, "POSIXt")) return(difftime(e1, e2))
  if (inherits(e2, "difftime")) e2 <- coerceTimeUnit(e2)
  if (!is.null(attr(e2, "class"))) stop("can only subtract numbers from \"POSIXt\" objects")
  e1 <- as.POSIXct(e1)
  .POSIXct(unclass(e1) - e2, attr(e1, "tzone"))
}

Ops.POSIXt <- function(e1, e2) {
  if (nargs() == 1) stop(gettextf("unary '%s' not defined for \"POSIXt\" objects", .Generic), domain = NA)
  boolean <- switch(.Generic, `<` = , `>` = , `==` = , `!=` = , `<=` = , `>=` = TRUE, FALSE)
  if (!boolean) stop(gettextf("'%s' not defined for \"POSIXt\" objects", .Generic), domain = NA)
  zoneOf <- function(other) if (is.null(z <- attr(other, "tzone"))) "" else z[1L]
  if (is.character(e1)) e1 <- as.POSIXct(e1, tz = zoneOf(e2))
  else if (inherits(e1, "POSIXlt")) e1 <- as.POSIXct(e1)
  if (is.character(e2)) e2 <- as.POSIXct(e2, tz = zoneOf(e1))
  else if (inherits(e2, "POSIXlt")) e2 <- as.POSIXct(e2)
  NextMethod(.Generic)
}

## ---- difftime --------------------------------------------------------------

difftime <- function(time1, time2, tz, units = c("auto", "secs", "mins", "hours", "days", "weeks")) {
  if (missing(tz)) {
    time1 <- as.POSIXct(time1)
    time2 <- as.POSIXct(time2)
  } else {
    time1 <- as.POSIXct(time1, tz = tz)
    time2 <- as.POSIXct(time2, tz = tz)
  }
  z <- unclass(time1) - unclass(time2)
  attr(z, "tzone") <- NULL
  units <- match.arg(units)
  if (units == "auto")
    units <- if (all(is.na(z))) "secs"
             else {
               zz <- min(abs(z), na.rm = TRUE)
               if (!is.finite(zz) || zz < 60) "secs"
               else if (zz < 3600) "mins"
               else if (zz < 86400) "hours"
               else "days"
             }
  switch(units,
    secs = .difftime(z, units = "secs"),
    mins = .difftime(z/60, units = "mins"),
    hours = .difftime(z/3600, units = "hours"),
    days = .difftime(z/86400, units = "days"),
    weeks = .difftime(z/(7 * 86400), units = "weeks"))
}

as.difftime <- function(tim, format = "%X", units = "auto", tz = "UTC") {
  if (inherits(tim, "difftime")) return(tim)
  if (is.character(tim)) {
    difftime(strptime(tim, format = format, tz = tz),
             strptime("0:0:0", format = "%X", tz = tz), units = units)
  } else {
    if (!is.numeric(tim)) stop("'tim' is not character or numeric")
    if (units == "auto") stop("need explicit units for numeric conversion")
    if (!(units %in% c("secs", "mins", "hours", "days", "weeks"))) stop("invalid units specified")
    .difftime(tim, units = units)
  }
}

units <- function(x) UseMethod("units")
units.difftime <- function(x) attr(x, "units")
`units<-` <- function(x, value) UseMethod("units<-")
`units<-.difftime` <- function(x, value) {
  from <- units(x)
  if (from == value) return(x)
  if (!(value %in% c("secs", "mins", "hours", "days", "weeks"))) stop("invalid units specified")
  sc <- cumprod(c(secs = 1, mins = 60, hours = 60, days = 24, weeks = 7))
  newx <- unclass(x) * as.vector(sc[from]/sc[value])
  .difftime(newx, value)
}

..coerceTimeUnit <- function(x)
  as.vector(switch(attr(x, "units"),
    secs = x, mins = 60 * x, hours = 3600 * x, days = 86400 * x, weeks = 7 * 86400 * x))

as.double.difftime <- function(x, units = "auto", ...) {
  if (units != "auto") units(x) <- units
  as.vector(unclass(x), "double")
}

format.difftime <- function(x, ...) paste(format(unclass(x), ...), units(x))

print.difftime <- function(x, digits = getOption("digits"), ...) {
  if (length(x) == 0L) cat("difftime of length 0\n")
  else if (is.array(x) || length(x) > 1L) {
    cat("Time differences in ", attr(x, "units"), "\n", sep = "")
    y <- unclass(x)
    attr(y, "units") <- NULL
    print(y, digits = digits, ...)
  } else cat("Time difference of ", format(unclass(x), digits = digits), " ",
             attr(x, "units"), "\n", sep = "")
  invisible(x)
}

`*.difftime` <- function(e1, e2) {
  if (inherits(e1, "difftime") && inherits(e2, "difftime"))
    stop("both arguments of * cannot be \"difftime\" objects")
  if (inherits(e2, "difftime")) {
    tmp <- e1
    e1 <- e2
    e2 <- tmp
  }
  .difftime(e2 * unclass(e1), attr(e1, "units"))
}

`/.difftime` <- function(e1, e2) {
  if (inherits(e2, "difftime")) stop("second argument of / cannot be a \"difftime\" object")
  .difftime(unclass(e1) / e2, attr(e1, "units"))
}

Ops.difftime <- function(e1, e2) {
  if (nargs() == 1) {
    switch(.Generic,
      `+` = {},
      `-` = e1[] <- -unclass(e1),
      stop(gettextf("unary '%s' not defined for \"difftime\" objects", .Generic), domain = NA))
    return(e1)
  }
  boolean <- switch(.Generic, `<` = , `>` = , `==` = , `!=` = , `<=` = , `>=` = TRUE, FALSE)
  if (boolean) {
    if (inherits(e1, "difftime") && inherits(e2, "difftime")) {
      e1 <- ..coerceTimeUnit(e1)
      e2 <- ..coerceTimeUnit(e2)
    }
    NextMethod(.Generic)
  } else if (.Generic == "+" || .Generic == "-") {
    if (inherits(e1, "difftime") && !inherits(e2, "difftime"))
      return(.difftime(NextMethod(.Generic), units = attr(e1, "units")))
    if (!inherits(e1, "difftime") && inherits(e2, "difftime"))
      return(.difftime(NextMethod(.Generic), units = attr(e2, "units")))
    u1 <- attr(e1, "units")
    if (attr(e2, "units") == u1) {
      .difftime(NextMethod(.Generic), units = u1)
    } else {
      e1 <- ..coerceTimeUnit(e1)
      e2 <- ..coerceTimeUnit(e2)
      res <- as.vector(NextMethod(.Generic))
      attr(res, "units") <- "secs"
      class(res) <- "difftime"
      res
    }
  } else {
    stop(gettextf("'%s' not defined for \"difftime\" objects", .Generic), domain = NA)
  }
}

Math.difftime <- function(x, ...) {
  switch(.Generic,
    abs = , sign = , floor = , ceiling = , trunc = , round = , signif = {
      units <- attr(x, "units")
      .difftime(NextMethod(), units)
    },
    stop(gettextf("'%s' not defined for \"difftime\" objects", .Generic), domain = NA))
}

Summary.difftime <- function(..., na.rm) {
  ok <- switch(.Generic, max = , min = , sum = , range = TRUE, FALSE)
  if (!ok) stop(gettextf("'%s' not defined for \"difftime\" objects", .Generic), domain = NA)
  x <- list(...)
  units <- sapply(x, attr, "units")
  if (all(units == (u1 <- units[1L]))) {
    args <- c(lapply(x, as.vector), na.rm = na.rm)
    .difftime(do.call(.Generic, args), u1)
  } else {
    args <- c(lapply(x, ..coerceTimeUnit), na.rm = na.rm)
    .difftime(do.call(.Generic, args), "secs")
  }
}

mean.difftime <- function(x, ...) .difftime(mean(unclass(x), ...), attr(x, "units"))

## ---- seq.Date --------------------------------------------------------------

seq.Date <- function(from, to, by, length.out = NULL, along.with = NULL, ...) {
  if (missing(from)) stop("'from' must be specified")
  if (!inherits(from, "Date")) stop("'from' must be a \"Date\" object")
  if (length(as.Date(from)) != 1L) stop("'from' must be of length 1")
  if (!missing(to)) {
    if (!inherits(to, "Date")) stop("'to' must be a \"Date\" object")
    if (length(as.Date(to)) != 1L) stop("'to' must be of length 1")
  }
  if (!missing(along.with)) length.out <- length(along.with)
  else if (!is.null(length.out)) {
    if (length(length.out) != 1L) stop("'length.out' must be of length 1")
    length.out <- ceiling(length.out)
  }
  if (missing(by)) {
    if (missing(to) && is.null(length.out))
      stop("without 'by', when one of 'to', 'from' is missing, 'length.out' / 'along.with' must be specified")
    from <- unclass(as.Date(from))
    to <- if (missing(to)) NULL else unclass(as.Date(to))
    res <- if (is.null(length.out)) seq.int(from, to) else seq.int(from, to, length.out = length.out)
    return(.Date(..wholeAsInteger(res)))
  }
  if (missing(to) + is.null(length.out) != 1L)
    stop("given 'by', exactly two of 'to', 'from' and 'length.out' / 'along.with' must be specified")
  if (length(by) != 1L) stop("'by' must be of length 1")
  valid <- 0L
  if (inherits(by, "difftime")) {
    by <- switch(attr(by, "units"), secs = 1/86400, mins = 1/1440, hours = 1/24,
                 days = 1, weeks = 7) * unclass(by)
  } else if (is.character(by)) {
    by2 <- strsplit(by, " ", fixed = TRUE)[[1L]]
    if (length(by2) > 2L || length(by2) < 1L) stop("invalid 'by' string")
    valid <- pmatch(by2[length(by2)], c("days", "weeks", "months", "quarters", "years", "DSTdays"))
    if (is.na(valid)) stop("invalid string for 'by'")
    if (valid <= 2L) {
      by <- c(1, 7)[valid]
      if (length(by2) == 2L) by <- by * as.integer(by2[1L])
    } else by <- if (length(by2) == 2L) as.integer(by2[1L]) else 1
  } else if (!is.numeric(by)) stop("invalid mode for 'by'")
  if (is.na(by)) stop("'by' is NA")
  if (valid <= 2L || valid == 6L) {
    from <- unclass(as.Date(from))
    if (!is.null(length.out)) res <- ..wholeAsInteger(seq.int(from, by = by, length.out = length.out))
    else {
      to0 <- unclass(as.Date(to))
      res <- seq.int(0, to0 - from, by) + from
    }
    res <- .Date(res)
  } else {
    r1 <- as.POSIXlt(from)
    if (valid == 5L) {
      if (missing(to)) {
        yr <- seq.int(r1$year, by = by, length.out = length.out)
      } else {
        to0 <- as.POSIXlt(to)
        yr <- seq.int(r1$year, to0$year, by)
      }
      r1$year <- yr
      res <- as.Date(r1)
    } else {
      if (valid == 4L) by <- by * 3
      if (missing(to)) {
        mon <- seq.int(r1$mon, by = by, length.out = length.out)
      } else {
        to0 <- as.POSIXlt(to)
        mon <- seq.int(r1$mon, 12 * (to0$year - r1$year) + to0$mon, by)
      }
      r1$mon <- mon
      res <- as.Date(r1)
    }
  }
  if (!missing(to)) {
    to <- as.Date(to)
    res <- if (by > 0) res[res <= to] else res[res >= to]
  }
  res
}

..wholeAsInteger <- function(x) if (all(is.na(x) | x == round(x))) as.integer(x) else x

as.list.Date <- function(x, ...) lapply(unclass(x), .Date, oldClass(x))

as.list.POSIXct <- function(x, ...) {
  nms <- names(x)
  tz <- attr(x, "tzone")
  cl <- oldClass(x)
  out <- lapply(unname(unclass(x)), .POSIXct, tz, cl)
  names(out) <- nms
  out
}

as.list.difftime <- function(x, ...) {
  units <- attr(x, "units")
  lapply(unclass(x), .difftime, units, oldClass(x))
}
