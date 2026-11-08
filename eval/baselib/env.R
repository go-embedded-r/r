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
  if (!USE.NAMES) names(res) <- NULL
  else if (length(res) == 0L) names(res) <- character(0)
  res
}

with <- function(data, expr, ...) UseMethod("with")

with.default <- function(data, expr, ...)
  eval(substitute(expr), data, enclos = parent.frame())
