# Base library, part 6: connections. Written from the help pages of
# connections, sink and capture.output in R 4.6.2 and checked against the
# oracle (docs/m3-plan.md). BSD licence. Connections, sinks and files are Go
# natives (eval/conn_fn.go, eval/files.go); this file holds the R-level glue.

summary <- function(object, ...) UseMethod("summary")

print.connection <- function(x, ...) {
  info <- unlist(summary.connection(x))
  cat("A connection with")
  print(matrix(info, ncol = 1L, dimnames = list(names(info), "")), ...)
  invisible(x)
}

capture.output <- function(..., file = NULL, append = FALSE,
                           type = c("output", "message"), split = FALSE) {
  type <- match.arg(type)
  ours <- is.null(file) || is.character(file)
  target <- file
  if (is.null(file)) {
    target <- textConnection("..captured", "w", local = TRUE)
  } else if (is.character(file)) {
    target <- file(file, if (isTRUE(append)) "a" else "w")
  }
  sink(target, type = type, split = split)
  finished <- FALSE
  finish <- function() {
    if (!finished) {
      finished <<- TRUE
      sink(type = type)
      if (ours) close(target)
    }
  }
  on.exit(finish())
  for (i in seq_len(...length())) {
    shown <- withVisible(...elt(i))
    if (shown$visible) print(shown$value)
  }
  finish()
  if (is.null(file)) ..captured else invisible(NULL)
}
