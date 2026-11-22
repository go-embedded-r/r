# Base library, slice 23 (2): the small helpers of base, and the task callback
# manager. Written from the help pages of R 4.6.2 (taskCallback, range, sort)
# and from calls made to the oracle.

# the sign of the comparison of two elements: 1, 0 or -1
.gt <- function(x, i, j) {
  if (x[i] > x[j]) 1L else if (x[i] < x[j]) -1L else 0L
}

# is x out of order: does an element exceed (or, strictly, equal) its successor
.gtn <- function(x, strictly) {
  n <- length(x)
  if (n < 2L) return(FALSE)
  if (strictly) any(x[-n] >= x[-1L]) else any(x[-n] > x[-1L])
}

# .rangeNum is a Go native (eval/slice23.go)

# a vector that was sorted (or not) comes back as it is: the wrappers carry
# no information a session could use
.doWrap <- function(vec, decr, nalast, noNA) vec
.doSortWrap <- function(vec, decr, nalast, noNA) vec

.signalSimpleWarning <- function(msg, call)
  withRestarts({
    warning(simpleWarning(msg, call))
    NULL
  }, muffleWarning = function() NULL)

.tryResumeInterrupt <- function() {
  r <- findRestart("resume")
  if (!is.null(r)) invokeRestart(r)
  invisible(NULL)
}

.check_tzones <- function(...) {
  tzs <- unlist(lapply(list(...), function(x) {
    z <- attr(x, "tzone")
    if (is.null(z)) NULL else z[1L]
  }))
  tzs <- unique(tzs[nzchar(tzs)])
  if (length(tzs) > 1L) warning("'tzone' attributes are inconsistent")
  if (length(tzs)) tzs[1L] else NULL
}

# zero.print for formatted zeros: the digits of a zero are blanked and the
# print string stands where the first 0 was
.format.zeros <- function(x, zero.print, nx = suppressWarnings(as.numeric(x)),
                          replace = FALSE, warn.non.fitting = TRUE) {
  if (!is.character(x)) stop("replacing substrings in a non-character object")
  if (is.na(zero.print)) stop("missing value where TRUE/FALSE needed")
  zero <- !nx & !is.na(nx)
  if (!any(zero)) return(x)
  z <- which(zero)
  if (replace) {
    x[z] <- zero.print
    return(x)
  }
  for (i in z) {
    s <- x[i]
    at <- regexpr("0", s, fixed = TRUE)
    if (at < 0L) {
      x[i] <- gsub("[0-9.eE+]", " ", s)
      next
    }
    room <- nchar(s) - at + 1L
    if (nchar(zero.print) > room) {
      if (warn.non.fitting)
        warning("'zero.print' is truncated to fit into formatted zeros; consider 'replace=TRUE'")
      zp <- substr(zero.print, 1L, room)
    } else zp <- zero.print
    x[i] <- paste0(substr(s, 1L, at - 1L), formatC(zp, width = -room))
  }
  x
}

# ---- taskCallbackManager -----------------------------------------------------

taskCallbackManager <- function(handlers = list(), registered = FALSE,
                                verbose = getOption("verbose")) {
  suspended <- FALSE
  add <- function(f, data = NULL, name = NULL, register = TRUE) {
    if (is.null(name)) name <- as.character(length(handlers) + 1L)
    handlers[[name]] <<- list(f = f, data = data)
    if (register && !registered) register()
    name
  }
  evaluate <- function(expr, value, ok, visible) {
    if (suspended) return(TRUE)
    discard <- character()
    for (i in names(handlers)) {
      h <- handlers[[i]]
      val <- if (!is.null(h[["data"]])) h[["f"]](expr, value, ok, visible, h[["data"]])
             else h[["f"]](expr, value, ok, visible)
      if (!val) discard <- c(discard, i)
    }
    if (length(discard)) handlers[discard] <<- NULL
    TRUE
  }
  remove <- function(which) {
    if (is.character(which)) {
      if (is.na(match(which, names(handlers)))) stop(gettextf("no such element '%s'", which))
      handlers[[which]] <<- NULL
    } else {
      which <- as.integer(which)
      if (is.na(which) || which < 1L || which > length(handlers)) stop("invalid 'which' argument")
      handlers <<- handlers[-which]
    }
    TRUE
  }
  register <- function(name = "R-taskCallbackManager", userData = NULL) {
    id <- addTaskCallback(evaluate, name = name)
    registered <<- TRUE
    id
  }
  suspend <- function(status = TRUE) {
    suspended <<- status
    status
  }
  callbacks <- function() handlers
  list(add = add, evaluate = evaluate, remove = remove, register = register,
       suspend = suspend, callbacks = callbacks)
}
