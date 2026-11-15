# Base library, slice 16 (3): writing and reading back R objects as text:
# dput, dget, dump and the deparse options. Written from the help pages of
# R 4.6.2 (dput, dget, dump, .deparseOpts) and from experiments with the oracle;
# deparse itself, and the options it takes, are in Go (eval/deparse16.go).

# the integer that stands for a set of deparse options
.deparseOpts <- function(control) {
  if (!length(control)) return(0)
  known <- c("all", "keepInteger", "quoteExpressions", "showAttributes", "useSource",
             "warnIncomplete", "delayPromises", "keepNA", "S_compatible", "hexNumeric",
             "digits17", "niceNames", "exact")
  bit <- c(all = 1631L, keepInteger = 1L, quoteExpressions = 2L, showAttributes = 4L, useSource = 8L,
           warnIncomplete = 16L, delayPromises = 32L, keepNA = 64L, S_compatible = 128L,
           hexNumeric = 256L, digits17 = 512L, niceNames = 1024L, exact = 1375L)
  hit <- pmatch(control, known)
  if (anyNA(hit))
    stop(sprintf("deparse option %s is not recognized", sQuote(control[is.na(hit)][1L])), domain = NA)
  as.numeric(Reduce(bitwOr, bit[known[hit]], 0L))
}

dput <- function(x, file = "", control = c("keepNA", "keepInteger", "niceNames", "showAttributes")) {
  if (is.character(file)) {
    if (nzchar(file)) {
      file <- file(file, "wt")
      on.exit(close(file))
    } else file <- stdout()
  }
  writeLines(deparse(x, control = control, backtick = TRUE), file)
  invisible(x)
}

dget <- function(file, keep.source = FALSE)
  eval(parse(file = file, keep.source = keep.source))

dump <- function(list, file = "dumpdata.R", append = FALSE, control = "all",
                 envir = parent.frame(), evaluate = TRUE) {
  if (!is.character(list)) stop("invalid first argument")
  if (is.character(file)) {
    if (nzchar(file)) {
      file <- file(file, if (append) "at" else "wt")
      on.exit(close(file))
    } else file <- stdout()
  }
  present <- vapply(list, function(name) exists(name, envir = envir), NA)
  for (name in list[!present]) warning(sprintf("object '%s' not found", name))
  found <- list[present]
  for (name in found) {
    head <- if (identical(name, make.names(name))) name else paste0("`", name, "`")
    writeLines(c(paste(head, "<-"), deparse(get(name, envir = envir), control = control, backtick = TRUE)), file)
  }
  invisible(found)
}
