# Base library, part 4: environments. Written from the help pages of
# as.list, eapply and with in R 4.6.2 and checked against the oracle
# (docs/m3-plan.md). BSD licence. The methods as.list.default,
# as.list.environment and as.list.function are Go natives.

as.list <- function(x, ...) UseMethod("as.list")

as.list.factor <- function(x, ...) {
  res <- vector("list", length(x))
  for (i in seq_along(x)) {
    one <- x[i]
    names(one) <- NULL
    res[[i]] <- one
  }
  names(res) <- names(x)
  res
}

eapply <- function(env, FUN, ..., all.names = FALSE, USE.NAMES = TRUE) {
  FUN <- match.fun(FUN)
  res <- lapply(as.list.environment(env, all.names = all.names), FUN, ...)
  if (!isTRUE(as.logical(USE.NAMES))) names(res) <- NULL
  else if (length(res) == 0L) names(res) <- character(0)
  res
}

with <- function(data, expr, ...) UseMethod("with")

with.default <- function(data, expr, ...)
  eval(substitute(expr), data, enclos = parent.frame())

within <- function(data, expr, ...) UseMethod("within")

within.list <- function(data, expr, ...) {
  env <- list2env(data, parent = parent.frame())
  eval(substitute(expr), env)
  now <- as.list.environment(env, all.names = TRUE)
  have <- names(now)
  gone <- names(data)[!(names(data) %in% have)]
  for (n in names(data)[names(data) %in% have]) data[n] <- list(now[[n]])
  fresh <- rev(have[!(have %in% names(data))])
  for (n in fresh) data[n] <- list(now[[n]])
  for (n in gone) data[[n]] <- NULL
  data
}
