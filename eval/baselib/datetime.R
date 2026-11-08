# Base library, part 9: dates, times and time differences. Written from the
# help pages of R 4.6.2 (Dates, DateTimeClasses, as.Date, as.POSIXlt,
# strptime, difftime, seq.Date, weekdays, ISOdatetime) and from experiments
# against the reference program, in the clean-room process of README.md. The
# calendar arithmetic is in Go (.dt* helpers, eval/datetime.go). BSD licence.

# ---- the three classes: constructors ---------------------------------------

.Date <- function(xx, cl = "Date") {
  class(xx) <- cl
  xx
}

.POSIXct <- function(xx, tz = NULL, cl = c("POSIXct", "POSIXt")) {
  class(xx) <- cl
  attr(xx, "tzone") <- tz
  xx
}

.difftime <- function(xx, units, cl = "difftime") {
  attr(xx, "units") <- NULL
  class(xx) <- cl
  attr(xx, "units") <- units
  xx
}

Sys.time <- function() .POSIXct(.dtNow())

Sys.Date <- function() as.Date(as.POSIXlt(Sys.time()))

# ---- small helpers (names start with a dot and are ours) -------------------

..units <- c("secs", "mins", "hours", "days", "weeks")
..unitSecs <- c(secs = 1, mins = 60, hours = 3600, days = 86400, weeks = 604800)

# seconds in a difftime or plain number, dropping attributes but the names
..asSecs <- function(z) {
  if (!inherits(z, "difftime")) return(z)
  v <- unclass(z) * ..unitSecs[[attr(z, "units")]]
  attr(v, "units") <- NULL
  v
}

# the time zone a date-time object carries ("" when it carries none)
..zoneOf <- function(x) {
  tz <- attr(x, "tzone")
  if (is.null(tz)) "" else tz[1L]
}

# the tzone attribute of a POSIXlt: one name for UTC and for a zone named by
# the caller of a parse, else the name with the two abbreviations
..tzAttr <- function(tz) {
  if (.dtIsUTC(tz)) return(if (nzchar(tz)) tz else .dtTzone(tz)[1L])
  .dtTzone(tz)
}

# wrap seconds since the epoch as a POSIXlt in zone tz
..mkLt <- function(secs, tz, nm = names(secs)) {
  lt <- .dtToLt(as.double(secs), tz)
  if (!is.null(nm)) names(lt$year) <- nm
  class(lt) <- c("POSIXlt", "POSIXt")
  attr(lt, "tzone") <- ..tzAttr(tz)
  attr(lt, "balanced") <- TRUE
  lt
}

# a difftime whose units attribute is set before its class
..mkDiff <- function(v, units) {
  attr(v, "units") <- units
  class(v) <- "difftime"
  v
}

..bareLt <- function(lt) {
  lt <- unclass(lt)
  attributes(lt) <- list(names = names(lt))
  lt
}

# ---- POSIXlt ---------------------------------------------------------------

as.POSIXlt <- function(x, tz = "", ...) UseMethod("as.POSIXlt")

as.POSIXlt.POSIXct <- function(x, tz = "", ...) {
  if (missing(tz)) tz <- ..zoneOf(x)
  ..mkLt(unclass(x), tz)
}

as.POSIXlt.Date <- function(x, tz = "UTC", ...) {
  ..mkLt(unclass(x) * 86400, "UTC")
}

as.POSIXlt.numeric <- function(x, tz = "", origin, ...) {
  as.POSIXlt(if (missing(origin)) as.POSIXct(x, tz = tz) else as.POSIXct(x, tz = tz, origin = origin), tz)
}

as.POSIXlt.character <- function(x, tz = "", format,
                                 tryFormats = c("%Y-%m-%d %H:%M:%OS", "%Y/%m/%d %H:%M:%OS",
                                                "%Y-%m-%d %H:%M", "%Y/%m/%d %H:%M",
                                                "%Y-%m-%d", "%Y/%m/%d"),
                                 optional = FALSE, ...) {
  x <- unclass(x)
  parse <- function(f) {
    lt <- strptime(x, f, tz = tz)
    if (nzchar(tz)) attr(lt, "tzone") <- tz
    lt
  }
  if (!missing(format)) return(parse(format))
  seen <- x[!is.na(x)]
  if (length(seen) == 0L) return(parse("%Y-%m-%d"))
  for (f in tryFormats)
    if (all(!is.na(strptime(seen, f, tz = tz)))) return(parse(f))
  if (optional) {
    x <- rep(NA_character_, length(x))
    return(parse("%Y-%m-%d"))
  }
  stop("character string is not in a standard unambiguous format")
}

as.POSIXlt.factor <- function(x, ...) {
  lt <- as.POSIXlt(as.character(x), ...)
  attr(lt, "balanced") <- NA
  lt
}

as.POSIXlt.default <- function(x, tz = "", optional = FALSE, ...) {
  if (inherits(x, "POSIXlt")) {
    if (!missing(tz)) attr(x, "tzone") <- tz
    return(x)
  }
  if (is.null(x)) return(as.POSIXlt(.POSIXct(numeric(), tz), tz))
  if (is.logical(x) && all(is.na(x))) return(as.POSIXlt(.POSIXct(as.numeric(x), tz), tz))
  if (optional) return(as.POSIXlt.character(rep.int(NA_character_, length(x)), tz = tz))
  stop(gettextf("do not know how to convert '%s' to class %s",
                deparse1(substitute(x)), dQuote("POSIXlt")), domain = NA)
}

strptime <- function(x, format, tz = "") {
  keep <- names(x)
  lt <- .dtParse(as.character(x), as.character(format), tz)
  if (!is.null(keep) && length(keep) == length(lt$year)) names(lt$year) <- keep
  class(lt) <- c("POSIXlt", "POSIXt")
  attr(lt, "tzone") <- ..tzAttr(tz)
  attr(lt, "balanced") <- TRUE
  lt
}

length.POSIXlt <- function(x) max(lengths(unclass(x)))

is.na.POSIXlt <- function(x) is.na(as.POSIXct(x))

# ---- POSIXct ---------------------------------------------------------------

as.POSIXct <- function(x, tz = "", ...) UseMethod("as.POSIXct")

as.POSIXct.POSIXlt <- function(x, tz = "", ...) {
  if (missing(tz)) tz <- ..zoneOf(x)
  secs <- .dtToCt(..bareLt(x), tz)
  nm <- names(unclass(x)$year)
  if (!is.null(nm)) names(secs) <- nm
  .POSIXct(secs, tz)
}

as.POSIXct.Date <- function(x, tz = "UTC", ...) .POSIXct(unclass(x) * 86400, "UTC")

as.POSIXct.numeric <- function(x, tz = "", origin, ...) {
  if (missing(origin)) return(.POSIXct(x, tz))
  .POSIXct(unclass(as.POSIXct(origin, tz = "GMT", ...)) + x, tz)
}

as.POSIXct.default <- function(x, tz = "", ...) {
  if (inherits(x, "POSIXct")) {
    if (!missing(tz)) attr(x, "tzone") <- tz
    return(x)
  }
  if (is.null(x)) return(.POSIXct(integer(), tz))
  if (is.logical(x) && all(is.na(x))) return(.POSIXct(as.numeric(x), tz))
  if (!(is.character(x) || is.factor(x)))
    stop(sprintf("do not know how to convert '%s' to class %s",
                 deparse1(substitute(x)), dQuote("POSIXct")))
  tz <- tz # now a value of this frame: not "missing" in the calls below
  as.POSIXct(as.POSIXlt(x, tz, ...), tz, ...)
}

ISOdatetime <- function(year, month, day, hour, min, sec, tz = "") {
  stamp <- paste(year, month, day, hour, min, sec, sep = "/")
  zone <- tz
  as.POSIXct(strptime(stamp, "%Y/%m/%d/%H/%M/%OS", tz = zone), tz = zone)
}

ISOdate <- function(year, month, day, hour = 12, min = 0, sec = 0, tz = "GMT")
  ISOdatetime(year, month, day, hour, min, sec, tz)

# ---- formatting and printing of date-times ---------------------------------

# the time-zone abbreviation of every element of a POSIXlt
..abbrev <- function(lt) {
  parts <- unclass(lt)
  z <- parts$zone
  tzn <- attr(lt, "tzone")
  n <- length(parts$year)
  if (is.null(z)) z <- character(n)
  z <- rep_len(z, n)
  blank <- !nzchar(z) | is.na(z)
  if (any(blank) && length(tzn) > 1L)
    z[blank] <- tzn[2L + (rep_len(parts$isdst, n)[blank] > 0)]
  else if (any(blank) && length(tzn) == 1L)
    z[blank] <- tzn
  z
}

format.POSIXlt <- function(x, format = "", usetz = FALSE,
                           digits = getOption("digits.secs"), ...) {
  parts <- ..bareLt(x)
  if (is.null(digits)) digits <- 0L
  digits <- min(6L, max(0L, digits))
  if (!is.character(format)) stop("invalid 'format' argument")
  nm <- names(unclass(x)$year)
  if (length(format) == 1L && !is.na(format) && format == "") {
    sec <- parts$sec
    live <- is.finite(sec)
    sec <- sec[live]
    np <- 0L
    while (np < digits && any(abs(sec - round(sec, np)) > 1e-6)) np <- np + 1L
    atMidnight <- all(c(parts$hour[live], parts$min[live], sec) == 0)
    format <- if (atMidnight) "%Y-%m-%d"
              else if (np > 0L) "%Y-%m-%d %H:%M:%OS"
              else "%Y-%m-%d %H:%M:%S"
    out <- .dtFormat(parts, format, np)
  } else out <- .dtFormat(parts, format, digits)
  if (usetz) {
    ab <- rep_len(..abbrev(x), length(out))
    real <- rep_len(is.finite(parts$sec), length(out))
    out <- ifelse(is.na(out) | !real, out, paste(out, ab))
  }
  if (!is.null(nm) && length(nm) == length(out)) names(out) <- nm
  out
}

format.POSIXct <- function(x, format = "", tz = "", usetz = FALSE, ...) {
  if (missing(tz)) tz <- ..zoneOf(x)
  out <- format.POSIXlt(as.POSIXlt(x, tz), format, usetz, ...)
  names(out) <- names(x)
  out
}

strftime <- function(x, format = "", tz = "", usetz = FALSE, ...) {
  if (!inherits(x, "POSIXlt")) x <- as.POSIXlt(x, tz = tz)
  format(x, format = format, usetz = usetz, ...)
}

# each element drops what is zero at its end: a midnight is just a date
as.character.POSIXt <- function(x, digits = if (inherits(x, "POSIXlt")) 14L else 6L,
                                OutDec = ".", ...) {
  lt <- if (inherits(x, "POSIXlt")) x else as.POSIXlt(x)
  parts <- ..bareLt(lt)
  n <- length(parts$year)
  out <- character(n)
  sec <- rep_len(parts$sec, n)
  for (i in seq_len(n)) {
    if (is.na(sec[i])) { out[i] <- NA_character_; next }
    s <- sec[i]
    one <- lapply(parts, function(col) rep_len(col, n)[i])
    whole <- abs(s - round(s)) < 1e-9
    f <- if (one$hour == 0 && one$min == 0 && s == 0) "%Y-%m-%d"
         else if (s == 0) "%Y-%m-%d %H:%M"
         else "%Y-%m-%d %H:%M:%S"
    txt <- .dtFormat(one, f, 0L)
    if (!whole && s != 0) {
      keep <- 0L
      for (k in seq_len(min(digits, 9L)))
        if (abs(s - round(s, k)) > 1e-9 * max(1, s)) keep <- k + 1L else { keep <- k; break }
      frac <- sub("0+$", "", sprintf(paste0("%.", min(max(keep, 1L), 9L), "f"), s - floor(s)))
      txt <- paste0(txt, substring(frac, 2L))
    }
    out[i] <- txt
  }
  nm <- names(unclass(lt)$year)
  if (!is.null(nm)) names(out) <- nm
  out
}

..printTimes <- function(x, tz, usetz, max, digits, ...) {
  if (is.null(max)) max <- getOption("max.print", 9999L)
  if (length(x) == 0L) {
    cat(class(x)[1L], " of length 0\n", sep = "")
  } else if (max < length(x)) {
    print(format(x[seq_len(max)], usetz = usetz), max = max + 1, ...)
    cat(" [ reached 'max' / getOption(\"max.print\") -- omitted",
        length(x) - max, "entries ]\n")
  } else print(format(x, usetz = usetz), max = max, ...)
  invisible(x)
}

print.POSIXct <- function(x, tz = "", usetz = TRUE, max = NULL,
                          digits = getOption("digits.secs"), ...) {
  shown <- if (missing(tz)) x else .POSIXct(unclass(x), tz)
  ..printTimes(shown, tz, usetz, max, digits, ...)
  invisible(x)
}

print.POSIXlt <- function(x, tz = "", usetz = TRUE, max = NULL, ...) {
  shown <- if (missing(tz)) x else as.POSIXlt(as.POSIXct(x), tz)
  ..printTimes(shown, tz, usetz, max, NULL, ...)
  invisible(x)
}

# ---- Date ------------------------------------------------------------------

as.Date <- function(x, ...) UseMethod("as.Date")

as.Date.default <- function(x, ...) {
  if (inherits(x, "Date")) return(x)
  if (is.null(x)) return(.Date(numeric()))
  if (is.logical(x) && all(is.na(x))) return(.Date(as.numeric(x)))
  stop(gettextf("do not know how to convert '%s' to class %s",
                deparse1(substitute(x)), dQuote("Date")), domain = NA)
}

as.Date.numeric <- function(x, origin, ...) {
  if (missing(origin)) .Date(x) else as.Date(origin, ...) + x
}

as.Date.factor <- function(x, ...) as.Date(as.character(x), ...)

as.Date.character <- function(x, format, tryFormats = c("%Y-%m-%d", "%Y/%m/%d"),
                              optional = FALSE, ...) {
  # settles the format from the first element that is not NA
  charToDate <- function(x) {
    seen <- x[!is.na(x)]
    if (length(seen) == 0L) return(NA_character_)
    for (f in tryFormats)
      if (!is.na(strptime(seen[1L], f, tz = "GMT"))) return(f)
    if (optional) return(NA_character_)
    stop("character string is not in a standard unambiguous format")
  }
  if (missing(format)) {
    format <- charToDate(x)
    if (is.na(format)) return(.Date(rep(NA_real_, length(x))))
  }
  as.Date(strptime(x, format, tz = "GMT"))
}

as.Date.POSIXlt <- function(x, ...) {
  days <- floor(.dtToCt(..bareLt(x), "UTC") / 86400)
  nm <- names(unclass(x)$year)
  if (!is.null(nm)) names(days) <- nm
  .Date(days)
}

as.Date.POSIXct <- function(x, tz = "UTC", ...) {
  if (tz == "UTC") {
    out <- .Date(floor(unclass(x) / 86400))
    attr(out, "tzone") <- NULL
    out
  } else as.Date(as.POSIXlt(x, tz = tz))
}

format.Date <- function(x, format = "%Y-%m-%d", ...) {
  out <- format.POSIXlt(as.POSIXlt(.Date(floor(unclass(x)))), format, ...)
  names(out) <- names(x)
  out
}

as.character.Date <- function(x, ...) {
  out <- format.Date(x, ...)
  out[is.na(x)] <- NA_character_
  out
}

print.Date <- function(x, max = NULL, ...) {
  if (is.null(max)) max <- getOption("max.print", 9999L)
  if (length(x) == 0L) cat(class(x)[1L], " of length 0\n", sep = "")
  else if (max < length(x)) {
    print(format(x[seq_len(max)]), max = max + 1, ...)
    cat(" [ reached 'max' / getOption(\"max.print\") -- omitted",
        length(x) - max, "entries ]\n")
  } else print(format(x), max = max, ...)
  invisible(x)
}

is.numeric.Date <- function(x) FALSE
is.numeric.POSIXt <- function(x) FALSE
is.numeric.difftime <- function(x) FALSE

# ---- weekdays, months, quarters, julian ------------------------------------

weekdays <- function(x, abbreviate) UseMethod("weekdays")
weekdays.Date <- function(x, abbreviate = FALSE)
  format(x, if (abbreviate) "%a" else "%A")
weekdays.POSIXt <- function(x, abbreviate = FALSE)
  format(x, if (abbreviate) "%a" else "%A")

months <- function(x, abbreviate) UseMethod("months")
months.Date <- function(x, abbreviate = FALSE)
  format(x, if (abbreviate) "%b" else "%B")
months.POSIXt <- function(x, abbreviate = FALSE)
  format(x, if (abbreviate) "%b" else "%B")

quarters <- function(x, ...) UseMethod("quarters")
quarters.Date <- function(x, ...)
  paste0("Q", 1L + as.POSIXlt(x)$mon %/% 3L)
quarters.POSIXt <- function(x, ...)
  paste0("Q", 1L + as.POSIXlt(x)$mon %/% 3L)

julian <- function(x, ...) UseMethod("julian")
julian.Date <- function(x, origin = as.Date("1970-01-01"), ...) {
  days <- unclass(x) - unclass(origin)
  attr(days, "origin") <- origin
  days
}

# ---- subsetting, combining, replicating ------------------------------------

`[.Date` <- function(x, ..., drop = TRUE) .Date(unclass(x)[...], oldClass(x))
`[[.Date` <- function(x, ..., drop = TRUE) .Date(unclass(x)[[...]], oldClass(x))
`[<-.Date` <- function(x, ..., value) {
  if (length(value) == 0L) return(x)
  days <- unclass(x)
  days[...] <- unclass(as.Date(value))
  .Date(days, oldClass(x))
}

`[.POSIXct` <- function(x, ..., drop = TRUE)
  .POSIXct(unclass(x)[...], attr(x, "tzone"), oldClass(x))
`[[.POSIXct` <- function(x, ..., drop = TRUE)
  .POSIXct(unclass(x)[[...]], attr(x, "tzone"), oldClass(x))

`[.difftime` <- function(x, ..., drop = TRUE)
  .difftime(unclass(x)[...], attr(x, "units"), oldClass(x))

c.Date <- function(..., recursive = FALSE)
  .Date(unlist(lapply(list(...), function(e) unclass(as.Date(e)))))

c.POSIXct <- function(..., recursive = FALSE) {
  parts <- list(...)
  zones <- lapply(parts, attr, "tzone")
  tz <- if (length(unique(zones)) == 1L) zones[[1L]] else NULL
  .POSIXct(unlist(lapply(parts, function(e) unclass(as.POSIXct(e)))), tz)
}

c.difftime <- function(..., recursive = FALSE) {
  parts <- list(...)
  if (length(parts) == 0L) return(.difftime(numeric(), "secs"))
  made <- vapply(parts, inherits, NA, "difftime")
  us <- unique(vapply(parts[made], attr, "", "units"))
  if (length(us) == 1L)
    return(.difftime(unlist(lapply(parts, as.vector)), us))
  .difftime(unlist(lapply(parts, function(e) as.vector(..asSecs(e)))), "secs")
}

rep.Date <- function(x, ...) {
  y <- NextMethod()
  .Date(y, oldClass(x))
}
rep.POSIXct <- function(x, ...) {
  y <- NextMethod()
  .POSIXct(y, attr(x, "tzone"), oldClass(x))
}
rep.difftime <- function(x, ...) .difftime(rep(as.vector(unclass(x)), ...), attr(x, "units"), oldClass(x))

as.list.Date <- function(x, ...) lapply(unclass(x), .Date, oldClass(x))
as.list.POSIXct <- function(x, ...) lapply(unclass(x), .POSIXct, attr(x, "tzone"), oldClass(x))
as.list.difftime <- function(x, ...) lapply(unclass(x), .difftime, attr(x, "units"), oldClass(x))

mean.Date <- function(x, ...) .Date(mean(unclass(x), ...))
mean.POSIXct <- function(x, ...) .POSIXct(mean(unclass(x), ...), attr(x, "tzone"))

Summary.Date <- function(..., na.rm = FALSE) {
  if (!(.Generic %in% c("min", "max", "range")))
    stop(gettextf("%s not defined for \"Date\" objects", .Generic), domain = NA)
  .Date(do.call(.Generic, c(lapply(list(...), function(e) unclass(as.Date(e))), na.rm = na.rm)))
}

Summary.POSIXct <- function(..., na.rm = FALSE) {
  if (!(.Generic %in% c("min", "max", "range")))
    stop(gettextf("'%s' not defined for \"POSIXt\" objects", .Generic), domain = NA)
  args <- list(...)
  v <- do.call(.Generic, c(lapply(args, function(e) unclass(as.POSIXct(e))), na.rm = na.rm))
  .POSIXct(v, attr(args[[1L]], "tzone"))
}

# ---- arithmetic and comparison ---------------------------------------------

# operators that compare; everything else on dates and times is refused
..comparison <- c("<", ">", "==", "!=", "<=", ">=")

# a day count from a number of days, or from a difftime of any units
..days <- function(z) if (inherits(z, "difftime")) ..asSecs(z) / 86400 else z

`+.Date` <- function(e1, e2) {
  if (nargs() < 2L) return(e1)
  if (inherits(e1, "Date") && inherits(e2, "Date"))
    stop("binary + is not defined for \"Date\" objects")
  .Date(unclass(..days(e1)) + unclass(..days(e2)))
}

`-.Date` <- function(e1, e2) {
  if (!inherits(e1, "Date")) stop("can only subtract from \"Date\" objects")
  if (nargs() < 2L) stop("unary - is not defined for \"Date\" objects")
  if (inherits(e2, "Date")) return(difftime(e1, e2, units = "days"))
  e2 <- ..days(e2)
  if (is.object(e2)) stop("can only subtract numbers from \"Date\" objects")
  .Date(unclass(as.Date(e1)) - e2)
}

Ops.Date <- function(e1, e2) {
  if (nargs() < 2L) stop(gettextf("unary %s not defined for \"Date\" objects", .Generic), domain = NA)
  if (!(.Generic %in% ..comparison))
    stop(gettextf("%s not defined for \"Date\" objects", .Generic), domain = NA)
  toDay <- function(z) unclass(if (is.character(z)) as.Date(z) else z)
  op <- get(.Generic, mode = "function")
  op(toDay(e1), toDay(e2))
}

`+.POSIXt` <- function(e1, e2) {
  if (nargs() < 2L) return(e1)
  if (inherits(e1, "POSIXt") && inherits(e2, "POSIXt"))
    stop("binary '+' is not defined for \"POSIXt\" objects")
  zoned <- if (inherits(e1, "POSIXt")) e1 else e2
  total <- unclass(..asSecs(if (inherits(e1, "POSIXlt")) as.POSIXct(e1) else e1)) +
    unclass(..asSecs(if (inherits(e2, "POSIXlt")) as.POSIXct(e2) else e2))
  .POSIXct(total, attr(if (inherits(zoned, "POSIXlt")) as.POSIXct(zoned) else zoned, "tzone"))
}

`-.POSIXt` <- function(e1, e2) {
  if (nargs() < 2L) stop("unary '-' is not defined for \"POSIXt\" objects")
  if (!inherits(e1, "POSIXt")) stop("can only subtract from \"POSIXt\" objects")
  if (inherits(e2, "POSIXt")) return(difftime(e1, e2))
  e2 <- ..asSecs(e2)
  if (is.object(e2)) stop("can only subtract numbers from \"POSIXt\" objects")
  e1 <- as.POSIXct(e1)
  shifted <- unclass(e1) - e2
  .POSIXct(shifted, attr(e1, "tzone"))
}

Ops.POSIXt <- function(e1, e2) {
  if (nargs() < 2L) stop(gettextf("unary '%s' not defined for \"POSIXt\" objects", .Generic), domain = NA)
  if (!(.Generic %in% ..comparison))
    stop(gettextf("'%s' not defined for \"POSIXt\" objects", .Generic), domain = NA)
  as.instant <- function(z) if (inherits(z, "POSIXlt") || is.character(z)) as.POSIXct(z) else z
  e1 <- as.instant(e1)
  e2 <- as.instant(e2)
  named <- unique(c(..zoneOf(e1), ..zoneOf(e2)))
  if (sum(nzchar(named)) > 1L) warning("'tzone' attributes are inconsistent")
  get(.Generic, mode = "function")(unclass(e1), unclass(e2))
}

# ---- difftime --------------------------------------------------------------

difftime <- function(time1, time2, tz,
                     units = c("auto", "secs", "mins", "hours", "days", "weeks")) {
  zoned <- !missing(tz)
  asCt <- function(t) if (zoned) as.POSIXct(t, tz = tz) else as.POSIXct(t)
  gap <- unclass(asCt(time1)) - unclass(asCt(time2))
  attr(gap, "tzone") <- NULL
  units <- match.arg(units)
  if (units == "auto") {
    least <- if (all(is.na(gap))) NA_real_ else min(abs(gap), na.rm = TRUE)
    units <- if (!is.finite(least) || least < 60) "secs"
             else if (least < 3600) "mins"
             else if (least < 86400) "hours"
             else "days"
  }
  .difftime(gap / ..unitSecs[[units]], units)
}

as.difftime <- function(tim, format = "%X", units = "auto", tz = "UTC") {
  if (inherits(tim, "difftime")) return(tim)
  if (is.character(tim)) {
    midnight <- strptime("0:0:0", format = "%X", tz = tz)
    return(difftime(strptime(tim, format = format, tz = tz), midnight, units = units))
  }
  if (!is.numeric(tim)) stop("'tim' is not character or numeric")
  if (units == "auto") stop("need explicit units for numeric conversion")
  if (!(units %in% ..units)) stop("invalid units specified")
  .difftime(tim, units = units)
}

units <- function(x) UseMethod("units")
units.difftime <- function(x) attr(x, "units")

`units<-` <- function(x, value) UseMethod("units<-")
`units<-.difftime` <- function(x, value) {
  from <- attr(x, "units")
  if (identical(from, value)) return(x)
  if (!(value %in% ..units)) stop("invalid units specified")
  y <- unclass(x) * (..unitSecs[[from]] / ..unitSecs[[value]])
  attr(y, "units") <- value
  class(y) <- "difftime"
  y
}

as.double.difftime <- function(x, units = "auto", ...) {
  if (units != "auto") units(x) <- units
  as.vector(unclass(x), "double")
}

format.difftime <- function(x, ..., with.units = TRUE) {
  body <- format(as.vector(unclass(x), "double"), ...)
  if (!is.null(names(x))) names(body) <- names(x)
  if (with.units) paste(body, attr(x, "units")) else body
}

print.difftime <- function(x, digits = getOption("digits"), ...) {
  u <- attr(x, "units")
  n <- length(x)
  if (n == 0L) cat("difftime of length 0\n")
  else if (n == 1L && !is.array(x))
    cat(sprintf("Time difference of %s %s\n", format(unclass(x), digits = digits), u))
  else {
    cat(sprintf("Time differences in %s\n", u))
    print(structure(unclass(x), units = NULL), digits = digits, ...)
  }
  invisible(x)
}

`*.difftime` <- function(e1, e2) {
  firstIs <- inherits(e1, "difftime")
  if (firstIs && inherits(e2, "difftime"))
    stop("both arguments of * cannot be \"difftime\" objects")
  span <- if (firstIs) e1 else e2
  factor <- if (firstIs) e2 else e1
  ..mkDiff(as.vector(unclass(span)) * factor, attr(span, "units"))
}

`/.difftime` <- function(e1, e2) {
  if (inherits(e2, "difftime"))
    stop("second argument of / cannot be a \"difftime\" object")
  ..mkDiff(as.vector(unclass(e1)) / e2, attr(e1, "units"))
}

Ops.difftime <- function(e1, e2) {
  cmp <- .Generic %in% c("<", ">", "==", "!=", "<=", ">=")
  if (nargs() == 1L) {
    if (!(.Generic %in% c("+", "-")))
      stop(gettextf("unary '%s' not defined for \"difftime\" objects", .Generic), domain = NA)
    return(.difftime(get(.Generic, mode = "function")(unclass(e1)), attr(e1, "units")))
  }
  if (!cmp && !(.Generic %in% c("+", "-")))
    stop(gettextf("'%s' not defined for \"difftime\" objects", .Generic), domain = NA)
  both <- inherits(e1, "difftime") && inherits(e2, "difftime")
  common <- attr(if (inherits(e1, "difftime")) e1 else e2, "units")
  mixed <- both && !identical(attr(e1, "units"), attr(e2, "units"))
  secondInSecs <- inherits(e2, "difftime") && identical(attr(e2, "units"), "secs")
  if (mixed) {
    e1 <- ..asSecs(e1)
    e2 <- ..asSecs(e2)
    common <- "secs"
  }
  r <- get(.Generic, mode = "function")(unclass(e1), unclass(e2))
  if (cmp) { attr(r, "units") <- NULL; r }
  else if (mixed && !secondInSecs) ..mkDiff(r, common)
  else .difftime(r, common)
}

Math.difftime <- function(x, ...) {
  if (!(.Generic %in% c("abs", "sign", "floor", "ceiling", "trunc", "round", "signif")))
    stop(gettextf("'%s' not defined for \"difftime\" objects", .Generic), domain = NA)
  .difftime(get(.Generic, mode = "function")(unclass(x), ...), attr(x, "units"))
}

Summary.difftime <- function(..., na.rm = FALSE) {
  if (!(.Generic %in% c("max", "min", "sum", "range")))
    stop(gettextf("'%s' not defined for \"difftime\" objects", .Generic), domain = NA)
  args <- list(...)
  us <- vapply(args, function(e) if (inherits(e, "difftime")) attr(e, "units") else NA_character_, "")
  if (length(unique(us)) == 1L) {
    v <- do.call(.Generic, c(lapply(args, unclass), na.rm = na.rm))
    return(.difftime(as.vector(v), us[1L]))
  }
  v <- do.call(.Generic, c(lapply(args, ..asSecs), na.rm = na.rm))
  .difftime(as.vector(v), "secs")
}

mean.difftime <- function(x, ...) .difftime(mean(unclass(x), ...), attr(x, "units"))

# ---- seq.Date --------------------------------------------------------------

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
  counted <- !is.null(length.out)
  if (missing(by)) {
    if (missing(to))
      stop("without 'by', when one of 'to', 'from' is missing, 'length.out' / 'along.with' must be specified")
    if (!counted) return(.Date(..wholeInts(seq.int(unclass(from), unclass(to), by = 1))))
    return(.Date(..wholeInts(seq.int(unclass(from), unclass(to), length.out = length.out))))
  }
  if (missing(to) == !counted)
    stop("given 'by', exactly two of 'to', 'from' and 'length.out' / 'along.with' must be specified")
  if (length(by) != 1L) stop("'by' must be of length 1")
  step <- by
  if (inherits(by, "difftime")) {
    step <- ..asSecs(by) / 86400
  } else if (is.character(by)) {
    words <- strsplit(by, " ", fixed = TRUE)[[1L]]
    if (length(words) > 2L || length(words) < 1L) stop("invalid 'by' string")
    kinds <- c("days", "weeks", "months", "quarters", "years", "DSTdays")
    k <- pmatch(words[length(words)], kinds)
    if (is.na(k)) stop("invalid string for 'by'")
    kind <- kinds[k]
    mult <- if (length(words) == 2L) as.integer(words[1L]) else 1L
    if (is.na(mult)) stop("invalid 'by' string")
    if (kind %in% c("days", "weeks", "DSTdays")) {
      step <- mult * c(days = 1, weeks = 7, DSTdays = 1)[[kind]]
    } else {
      if (kind == "quarters") { mult <- mult * 3L; kind <- "months" }
      lt <- as.POSIXlt(from)
      if (kind == "years") mult <- mult * 12L
      if (!missing(to)) {
        lo <- as.POSIXlt(to)
        total <- (lo$year - lt$year) * 12L + (lo$mon - lt$mon)
        count <- total %/% mult + 1L
        if (mult < 0L) count <- (-total) %/% (-mult) + 1L
      } else count <- length.out
      idx <- (seq_len(max(count, 0L)) - 1L) * mult
      first <- lt$mon + idx
      lt$year <- lt$year + first %/% 12L
      lt$mon <- first %% 12L
      lt$mday <- rep_len(lt$mday, length(idx))
      lt$sec <- lt$min <- lt$hour <- rep_len(0, length(idx))
      lt$isdst <- rep_len(0L, length(idx))
      lt$wday <- lt$yday <- rep_len(NA_integer_, length(idx))
      out <- as.Date(lt)
      if (!missing(to)) out <- if (mult > 0L) out[out <= to] else out[out >= to]
      return(out)
    }
  } else if (!is.numeric(by)) stop("'by' must be a number, a difftime or a character string")
  if (!missing(to)) {
    return(.Date(seq.int(unclass(from), unclass(to), by = step)))
  }
  .Date(..wholeInts(seq.int(unclass(from), by = step, length.out = length.out)))
}

# a run of whole numbers that fits an integer is stored as integers
..wholeInts <- function(v) {
  if (!anyNA(v) && all(v == floor(v)) && all(abs(v) < 2147483647)) as.integer(v) else v
}
