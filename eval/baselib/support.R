# Base library, part 6: small functions the others lean on. Written from
# the help pages of getOption, deparse, gettext, pmatch, rep and factor in
# R 4.6.2 and checked against the oracle (docs/m3-plan.md). BSD licence.

getOption <- function(x, default = NULL) {
  v <- options(x)[[1L]]
  if (is.null(v)) default else v
}

deparse1 <- function(expr, collapse = " ", width.cutoff = 500L, ...)
  paste(deparse(expr, width.cutoff, ...), collapse = collapse)

gettext <- function(..., domain = NULL, trim = TRUE) {
  args <- lapply(list(...), as.character)
  unlist(args)
}

gettextf <- function(fmt, ..., domain = NULL, trim = TRUE)
  sprintf(gettext(fmt, domain = domain, trim = trim), ...)

rep.int <- function(x, times) rep(x, times = times)

is.factor <- function(x) inherits(x, "factor")

pmatch <- function(x, table, nomatch = NA_integer_, duplicates.ok = FALSE) {
  x <- as.character(x)
  table <- as.character(table)
  res <- rep(as.integer(nomatch), length(x))
  done <- rep(FALSE, length(x))
  taken <- rep(FALSE, length(table))
  for (pass in 1:2) {
    for (i in seq_along(x)) {
      if (done[i]) next
      if (is.na(x[i]) || !nzchar(x[i])) next
      if (pass == 1L) {
        hit <- which(table == x[i] & (duplicates.ok | !taken))
        if (length(hit)) {
          res[i] <- hit[1L]
          done[i] <- TRUE
          taken[hit[1L]] <- TRUE
        }
      } else {
        hit <- which(startsWith(table, x[i]) & (duplicates.ok | !taken))
        if (length(hit) == 1L) {
          res[i] <- hit
          done[i] <- TRUE
          taken[hit] <- TRUE
        }
      }
    }
  }
  res
}

unname <- function(obj, force = FALSE) {
  if (!is.null(names(obj))) names(obj) <- NULL
  if (!is.null(dimnames(obj))) dimnames(obj) <- NULL
  obj
}
