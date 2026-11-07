# Base library, part 2: strings. Written from the help pages of R 4.6.2 and
# checked against the oracle (docs/m3-plan.md). BSD licence.

toString <- function(x, ...) paste(x, collapse = ", ")

sQuote <- function(x, q = TRUE) {
  if (!length(x)) return(character())
  if (isFALSE(q)) paste0("'", x, "'") else paste0("‘", x, "’")
}

dQuote <- function(x, q = TRUE) {
  if (!length(x)) return(character())
  if (isFALSE(q)) paste0("\"", x, "\"") else paste0("“", x, "”")
}

shQuote <- function(string, type = c("sh", "csh", "cmd", "cmd2")) {
  type <- match.arg(type)
  if (type == "cmd") return(paste0("\"", gsub("\"", "\\\\\"", string), "\""))
  if (type == "cmd2") return(gsub("([()%!^\"<>&|])", "^\\1", string))
  if (type == "csh") {
    has_single_csh <- grepl("'", string, fixed = TRUE)
    return(ifelse(has_single_csh, paste0("\"", gsub("([\"$`\\\\!])", "\\\\\\1", string), "\""), paste0("'", string, "'")))
  }
  has_single <- grepl("'", string, fixed = TRUE)
  ifelse(has_single, paste0("\"", gsub("([\"$`\\\\])", "\\\\\\1", string), "\""), paste0("'", string, "'"))
}

noquote <- function(obj, right = FALSE) {
  if (!inherits(obj, "noquote")) class(obj) <- c(attr(obj, "class"), "noquote")
  obj
}

print.noquote <- function(x, ...) {
  y <- x
  cl <- oldClass(y)
  keep <- cl[cl != "noquote"]
  oldClass(y) <- if (length(keep)) keep else NULL
  print(y, quote = FALSE)
  invisible(x)
}

file.path <- function(..., fsep = "/") {
  args <- list(...)
  if (length(args) == 0L || any(lengths(args) == 0L)) return(character())
  paste(..., sep = fsep)
}

