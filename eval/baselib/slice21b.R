# Slice 21, part B: pairlists, regmatches<-, gregexec, formatDL, data frame
# methods, all.equal methods, conditions, restarts, withAutoprint, hooks, task
# callbacks, srcfiles. Written from the help pages of R 4.6.2 and the oracle's
# answers (docs/m3-plan.md); BSD licence.

as.pairlist <- function(x) if (is.null(x)) NULL else ..mkpair(as.list(x))

..rm_join <- function(a, b) paste0(c(rbind(a[seq_along(b)], b), a[length(b) + 1L]), collapse = "")

..rm_gaps <- function(s, st, ln) {
  en <- st + ln - 1L
  k <- length(st)
  gaps <- substring(s, c(1L, en + 1L), c(st - 1L, nchar(s)))
  list(gaps = gaps, hits = substring(s, st, en))
}

`regmatches<-` <- function(x, m, invert = FALSE, value) {
  n <- length(x)
  if (is.list(m)) {
    value <- rep_len(as.list(value), n)
    for (i in seq_len(n)) {
      st <- as.integer(m[[i]])
      ln <- attr(m[[i]], "match.length")
      none <- is.na(st[1L]) || st[1L] == -1L
      if (none) {
        if (invert && length(value[[i]])) x[i] <- as.character(value[[i]])[1L]
        next
      }
      pieces <- ..rm_gaps(x[i], st, ln)
      x[i] <- if (invert) ..rm_join(as.character(value[[i]]), pieces$hits) else ..rm_join(pieces$gaps, rep_len(as.character(value[[i]]), length(st)))
    }
    return(x)
  }
  ml <- attr(m, "match.length")
  hit <- which(!is.na(m) & m != -1L)
  if (invert) {
    value <- rep_len(as.list(value), n)
    for (i in hit) {
      pieces <- ..rm_gaps(x[i], m[i], ml[i])
      x[i] <- ..rm_join(as.character(value[[i]]), pieces$hits)
    }
    for (i in setdiff(seq_len(n), hit)) if (!is.na(m[i]) && length(value[[i]])) x[i] <- as.character(value[[i]])[1L]
    return(x)
  }
  value <- as.character(value)
  if (!length(hit)) return(x)
  if (!length(value)) stop("must have replacement values for matches")
  value <- rep_len(value, length(hit))
  for (j in seq_along(hit)) {
    i <- hit[j]
    pieces <- ..rm_gaps(x[i], m[i], ml[i])
    x[i] <- ..rm_join(pieces$gaps, value[j])
  }
  x
}

gregexec <- function(pattern, text, ignore.case = FALSE, perl = FALSE, fixed = FALSE, useBytes = FALSE) {
  if (!is.character(text)) text <- as.character(text)
  ascii <- !any(grepl("[^\001-\177]", text, useBytes = TRUE), na.rm = TRUE)
  whole <- gregexpr(pattern, text, ignore.case = ignore.case, perl = perl, fixed = fixed, useBytes = useBytes)
  lapply(seq_along(text), function(i) {
    st <- as.integer(whole[[i]])
    ln <- attr(whole[[i]], "match.length")
    if (is.na(st[1L]) || st[1L] == -1L) {
      return(structure(-1L, match.length = -1L, index.type = "chars", useBytes = ascii))
    }
    cols <- lapply(seq_along(st), function(j) {
      piece <- substring(text[i], st[j], st[j] + ln[j] - 1L)
      r <- regexec(pattern, piece, ignore.case = ignore.case, perl = perl, fixed = fixed, useBytes = useBytes)[[1L]]
      list(start = as.integer(r) + st[j] - 1L, len = attr(r, "match.length"))
    })
    start <- vapply(cols, function(cc) cc$start, integer(length(cols[[1L]]$start)))
    len <- vapply(cols, function(cc) cc$len, integer(length(cols[[1L]]$len)))
    if (is.null(dim(start))) {
      start <- matrix(start, ncol = length(st))
      len <- matrix(len, ncol = length(st))
    }
    shape <- dim(start)
    res <- as.vector(start)
    attr(res, "match.length") <- len
    attr(res, "useBytes") <- ascii
    attr(res, "index.type") <- "chars"
    dim(res) <- shape
    res
  })
}

formatDL <- function(x, y, style = c("table", "list"), width = 0.9 * getOption("width"), indent = NULL) {
  style <- match.arg(style)
  if (is.list(x)) {
    if (length(x) == 2L && diff(lengths(x)) == 0L) {
      y <- x[[2L]]
      x <- x[[1L]]
    } else stop("incorrect value for 'x'")
  } else if (is.matrix(x)) {
    y <- x[, 2L]
    x <- x[, 1L]
  } else {
    if (missing(y)) {
      if (is.null(names(x))) stop(simpleError("argument \"y\" is missing, with no default", sys.call()))
      y <- x
      x <- names(x)
    }
    if (length(x) != length(y)) stop("'x' and 'y' must have the same length")
  }
  x <- as.character(x)
  y <- as.character(y)
  if (is.null(indent)) indent <- if (style == "table") width / 3 else width / 9
  if (indent > 0.5 * width) stop("incorrect values of 'indent' and 'width'")
  pad <- strrep(" ", floor(indent))
  one <- function(i) {
    if (style == "list") {
      return(paste(strwrap(paste0(x[i], ": ", y[i]), width = width, exdent = indent), collapse = "\n"))
    }
    body <- strwrap(y[i], width = width - indent)
    if (!length(body)) body <- ""
    if (nchar(x[i], "w") >= indent - 2) {
      paste0(x[i], "\n", paste0(pad, body, collapse = "\n"))
    } else {
      paste0(formatC(x[i], width = -floor(indent)), paste(body, collapse = paste0("\n", pad)))
    }
  }
  vapply(seq_along(x), one, "")
}

merge.default <- function(x, y, ...) merge(as.data.frame(x), as.data.frame(y), ...)

`split<-.data.frame` <- function(x, f, drop = FALSE, ..., value) {
  n <- length(value)
  ix <- split(seq_len(nrow(x)), f, drop = drop, ...)
  j <- 0L
  for (i in ix) {
    j <- j %% n + 1L
    x[i, ] <- value[[j]]
  }
  x
}

anyNA.data.frame <- function(x, recursive = FALSE) {
  for (col in x) if (anyNA(col, recursive = recursive)) return(TRUE)
  FALSE
}

anyDuplicated.data.frame <- function(x, incomparables = FALSE, fromLast = FALSE, ...) {
  if (!length(x)) return(0L)
  key <- do.call(paste, c(lapply(x, function(col) ifelse(is.na(col), "\001NA", as.character(col))), sep = "\002"))
  anyDuplicated(key, fromLast = fromLast)
}

as.vector.data.frame <- function(x, mode = "any") {
  x <- as.list(x)
  if (mode == "any" || mode == "list") return(x)
  if (mode == "character") {
    return(unname(vapply(x, function(e) if (is.character(e) && length(e) == 1L) e else paste(deparse(e), collapse = " "), "")))
  }
  as.vector(x, mode)
}

all.equal.POSIXt <- function(target, current, ..., tolerance = 1e-3, scale) {
  if (!inherits(target, "POSIXt")) return("'target' is not a POSIXt")
  if (!inherits(current, "POSIXt")) return("'current' is not a POSIXt")
  tz1 <- attr(as.POSIXlt(target), "tzone")[1L]
  tz2 <- attr(as.POSIXlt(current), "tzone")[1L]
  msg <- if (!identical(tz1, tz2)) paste0("'tzone' attributes are inconsistent ('", tz1, "' and '", tz2, "')")
  r <- all.equal(as.numeric(as.POSIXct(target)), as.numeric(as.POSIXct(current)), ..., tolerance = tolerance, scale = 1)
  if (isTRUE(r)) r <- NULL
  if (is.null(msg) && is.null(r)) TRUE else c(msg, r)
}

all.equal.formula <- function(target, current, ...) {
  if (!inherits(current, "formula")) current <- tryCatch(as.formula(current), error = function(e) current)
  if (length(target) != length(current)) {
    return(paste("target, current differ in having response: ", length(target) == 3L, ", ", length(current) == 3L, sep = ""))
  }
  if (!identical(deparse(target), deparse(current))) "formulas differ in contents" else TRUE
}

messageCondition <- function(message, ..., class = character(), call = NULL) {
  structure(class = c(class, "message", "condition"), list(message = as.character(message), call = call, ...))
}

packageStartupMessage <- function(..., domain = NULL, appendLF = TRUE) {
  msg <- paste0(..., if (appendLF) "\n")
  message(structure(class = c("packageStartupMessage", "simpleMessage", "message", "condition"), list(message = msg, call = sys.call())))
}

suppressPackageStartupMessages <- function(expr) {
  withCallingHandlers(expr, packageStartupMessage = function(c) invokeRestart("muffleMessage"))
}

isRestart <- function(x) inherits(x, "restart")
restartDescription <- function(r) r$description
restartFormals <- function(r) formals(r$handler)

withAutoprint <- function(exprs, evaluated = FALSE, local = parent.frame(), print. = TRUE, echo = TRUE,
                          max.deparse.length = Inf, width.cutoff = max(20, getOption("width")),
                          deparseCutoff = width.cutoff, skip.echo = 0, prompt.echo = getOption("prompt", "> "),
                          continue.echo = getOption("continue", "+ "), ...) {
  if (!evaluated) exprs <- substitute(exprs)
  if (is.call(exprs) && identical(exprs[[1L]], as.name("{"))) exprs <- as.list(exprs)[-1L]
  else if (is.expression(exprs)) exprs <- as.list(exprs)
  else exprs <- list(exprs)
  res <- list(value = NULL, visible = FALSE)
  for (e in exprs) {
    if (echo) {
      lines <- deparse(e, width.cutoff = deparseCutoff)
      if (length(lines) > 1L && is.finite(max.deparse.length)) lines <- lines[seq_len(max(1L, max.deparse.length))]
      cat(paste0(c(prompt.echo, rep(continue.echo, length(lines) - 1L)), lines), sep = "\n")
    }
    res <- withVisible(eval(e, local))
    if (print. && res$visible) print(res$value)
  }
  invisible(res)
}

.traceback <- function(x = NULL, max.lines = getOption("traceback.max.lines", getOption("deparse.max.lines", -1L))) {
  if (is.null(x) && exists(".Traceback", envir = baseenv(), inherits = FALSE)) x <- get(".Traceback", envir = baseenv())
  x
}

traceback <- function(x = NULL, max.lines = getOption("traceback.max.lines", getOption("deparse.max.lines", -1L))) {
  x <- .traceback(x, max.lines = max.lines)
  n <- length(x)
  if (n == 0L) {
    cat("No traceback available", "\n")
  } else {
    for (i in seq_len(n)) {
      label <- paste0(n - i + 1L, ": ")
      xi <- as.character(x[[i]])
      cat(paste0(c(label, rep(strrep(" ", nchar(label)), length(xi) - 1L)), xi), sep = "\n")
    }
  }
  invisible(x)
}

..userHooks <- new.env(parent = emptyenv())

getHook <- function(hookName) {
  if (exists(hookName, envir = ..userHooks, inherits = FALSE)) get(hookName, envir = ..userHooks, inherits = FALSE) else list()
}

setHook <- function(hookName, value, action = c("append", "prepend", "replace")) {
  action <- match.arg(action)
  old <- getHook(hookName)
  new <- switch(action, append = c(old, value), prepend = c(value, old), replace = value)
  if (length(new)) assign(hookName, new, envir = ..userHooks) else if (exists(hookName, envir = ..userHooks, inherits = FALSE)) rm(list = hookName, envir = ..userHooks)
  invisible()
}

..taskCallbacks <- new.env(parent = emptyenv())
..taskCallbacks$list <- list()

addTaskCallback <- function(f, data = NULL, name = character()) {
  if (!is.function(f)) stop("handler must be a function")
  cur <- ..taskCallbacks$list
  id <- length(cur) + 1L
  nm <- if (length(name)) as.character(name)[1L] else as.character(id)
  cur[[length(cur) + 1L]] <- list(f = f, data = data, name = nm)
  ..taskCallbacks$list <- cur
  structure(id, names = nm)
}

removeTaskCallback <- function(id) {
  cur <- ..taskCallbacks$list
  nms <- vapply(cur, function(e) e$name, "")
  pos <- if (is.character(id)) match(id, nms) else as.integer(id)
  if (length(pos) != 1L || is.na(pos) || pos < 1L || pos > length(cur)) return(FALSE)
  ..taskCallbacks$list <- cur[-pos]
  TRUE
}

getTaskCallbackNames <- function() vapply(..taskCallbacks$list, function(e) e$name, "")

print.function <- function(x, useSource = TRUE, ...) {
  print.default(x, useSource = useSource, ...)
  invisible(x)
}

write.dcf <- function(x, file = "", append = FALSE, useBytes = FALSE, indent = 0.1 * getOption("width"),
                      width = 0.9 * getOption("width"), keep.white = NULL) {
  if (is.atomic(x) && !is.matrix(x)) x <- as.list(x)
  x <- as.data.frame(x, stringsAsFactors = FALSE)
  x[] <- lapply(x, as.character)
  recs <- lapply(seq_len(nrow(x)), function(i) {
    out <- character()
    for (tag in names(x)) {
      val <- x[[tag]][i]
      if (is.na(val)) next
      if (tag %in% keep.white) {
        out <- c(out, strsplit(paste0(tag, ": ", gsub("\n([^[:blank:]])", "\n \\1", val)), "\n", fixed = TRUE)[[1L]])
      } else {
        out <- c(out, strwrap(paste0(tag, ": ", val), width = width, exdent = indent))
      }
    }
    out
  })
  lines <- character()
  for (i in seq_along(recs)) lines <- c(lines, if (i > 1L) "", recs[[i]])
  if (is.character(file) && !nzchar(file)) cat(lines, sep = "\n")
  else if (is.character(file)) {
    con <- file(file, if (append) "a" else "w")
    on.exit(close(con))
    writeLines(lines, con)
  } else writeLines(lines, file)
  invisible()
}

sys.status <- function() list(sys.calls = sys.calls(), sys.parents = sys.parents(), sys.frames = sys.frames())

truncate <- function(con, ...) UseMethod("truncate")

chooseOpsMethod <- function(x, y, mx, my, cl, reverse) UseMethod("chooseOpsMethod")
chooseOpsMethod.default <- function(x, y, mx, my, cl, reverse) FALSE

nameOfClass <- function(x) UseMethod("nameOfClass")
nameOfClass.default <- function(x) NULL

srcfile <- function(filename, encoding = getOption("encoding"), Enc = "unknown") {
  e <- new.env(hash = TRUE, parent = emptyenv())
  e$filename <- filename
  e$encoding <- encoding
  e$Enc <- Enc
  e$isFile <- file.exists(filename)
  e$timestamp <- if (e$isFile) file.mtime(filename) else Sys.time()
  e$wd <- getwd()
  class(e) <- "srcfile"
  e
}

srcfilecopy <- function(filename, lines, timestamp = Sys.time(), isFile = FALSE) {
  e <- new.env(hash = TRUE, parent = emptyenv())
  e$filename <- filename
  e$wd <- getwd()
  e$isFile <- isFile
  e$lines <- as.character(lines)
  e$timestamp <- timestamp
  e$Enc <- "unknown"
  class(e) <- c("srcfilecopy", "srcfile")
  e
}

srcfilealias <- function(filename, srcfile) {
  e <- new.env(hash = TRUE, parent = emptyenv())
  e$filename <- filename
  e$original <- srcfile
  class(e) <- c("srcfilealias", "srcfile")
  e
}

getSrcLines <- function(srcfile, first, last) {
  if (first > last) return(character())
  while (inherits(srcfile, "srcfilealias")) srcfile <- srcfile$original
  lines <- if (inherits(srcfile, "srcfilecopy")) srcfile$lines else readLines(srcfile$filename, warn = FALSE)
  lines[seq.int(first, min(last, length(lines)))]
}

srcref <- function(srcfile, lloc) {
  stopifnot(inherits(srcfile, "srcfile"))
  lloc <- as.integer(lloc)
  if (length(lloc) == 4L) lloc <- c(lloc, lloc[2L], lloc[4L], lloc[1L], lloc[3L])
  else if (length(lloc) == 6L) lloc <- c(lloc, lloc[1L], lloc[3L])
  structure(lloc, srcfile = srcfile, class = "srcref")
}

as.character.srcref <- function(x, useSource = TRUE, to = x, ...) {
  loc <- unclass(x)
  last <- unclass(to)
  sf <- attr(x, "srcfile")
  first_line <- loc[1L]
  last_line <- last[3L]
  first_col <- if (length(loc) >= 6L) loc[5L] else loc[2L]
  last_col <- if (length(last) >= 6L) last[6L] else last[4L]
  lines <- getSrcLines(sf, first_line, last_line)
  if (length(lines)) {
    nl <- length(lines)
    lines[nl] <- substring(lines[nl], 1L, last_col)
    lines[1L] <- substring(lines[1L], first_col)
  }
  lines
}

print.srcref <- function(x, useSource = TRUE, ...) {
  if (useSource) {
    cat(as.character(x, useSource = TRUE), sep = "\n")
  } else {
    loc <- unclass(x)
    c1 <- if (length(loc) >= 6L) loc[5L] else loc[2L]
    c2 <- if (length(loc) >= 6L) loc[6L] else loc[4L]
    cat("<srcref: file \"", attr(x, "srcfile")$filename, "\" chars ", loc[1L], ":", c1, " to ", loc[3L], ":", c2, ">\n", sep = "")
  }
  invisible(x)
}

summary.srcref <- function(object, useSource = FALSE, ...) {
  print(object, useSource = useSource, ...)
  invisible(object)
}

print.srcfile <- function(x, ...) {
  cat(x$filename, "\n")
  invisible(x)
}

summary.srcfile <- function(object, ...) {
  cat(object$filename, "\n")
  if (inherits(object$timestamp, "POSIXt")) cat("Timestamp: ", format(object$timestamp, usetz = TRUE), "\n", sep = "")
  enc <- object$encoding
  cat("Encoding: \"", enc, "\"", sep = "")
  if (!is.null(object$Enc) && object$Enc != "unknown" && !identical(object$Enc, enc)) {
    cat(", re-encoded to \"", object$Enc, "\"", sep = "")
  }
  cat("\n")
  invisible(object)
}

# the connection of a source file stays open between calls; `line` is the line
# to read up to, and is required
..srcfile_position <- function(con, line, make, ...) {
  if (!is.null(con$line) && con$line > line) close(con)
  conn <- con$conn
  if (is.null(conn)) {
    con$conn <- conn <- make()
    con$line <- 1L
  } else if (!isOpen(conn)) {
    open(conn, "r")
    con$line <- 1L
  }
  if (con$line < line) {
    readLines(conn, line - con$line, warn = FALSE)
    con$line <- line
  }
  invisible(conn)
}

open.srcfile <- function(con, line, ...) {
  mt <- file.mtime(con$filename)
  if (is.na(mt) || mt != con$timestamp) {
    warning(gettextf("Timestamp of %s has changed", sQuote(con$filename)), call. = FALSE)
  }
  ..srcfile_position(con, line, function() file(con$filename, "r", encoding = con$encoding))
}

open.srcfilecopy <- function(con, line, ...) {
  ..srcfile_position(con, line, function() textConnection(con$lines, "r"))
}

open.srcfilealias <- function(con, line, ...) open(con$original, line, ...)

close.srcfile <- function(con, ...) {
  # observed: a source file with an open connection answers invisibly, one without answers NULL visibly
  if (is.null(con$conn)) return(NULL)
  close(con$conn)
  rm("conn", envir = con)
  invisible(NULL)
}

close.srcfilealias <- function(con, ...) close(con$original, ...)
