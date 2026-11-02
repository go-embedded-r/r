# The oracle: a request/reply loop around GNU R, written for this repository
# (BSD-3-Clause, no code from R). It is run by internal/oracle as
#
#     Rscript --vanilla -e <bootstrap>
#
# where the bootstrap (see oracle.go) reads this file from standard input, as a
# length-prefixed frame, and evaluates it: Rscript refuses a long -e argument.
#
# Framing. A request is a 4-byte big-endian length followed by that many bytes
# of serialize() output (XDR, written by the Go side): a list whose first
# element is the operation name (a character vector of length 1) and whose
# other elements are its arguments. A reply is one line of lowercase hex: the
# serialize(xdr = TRUE, version = 3) bytes of a list. The first reply is the
# greeting, sent before any request is read. Hex keeps the reply on a text
# stdout, which R cannot write raw bytes to.
#
# Every reply is list(ok = TRUE, ...) or list(ok = FALSE, class, message).
# Nothing a request evaluates can write to the reply stream: printed output is
# captured and returned as `output`.

local({
  stdin_con <- get0(".rgo_stdin", envir = globalenv(), ifnotfound = file("stdin", "rb"))
  # leave the global environment as R starts it: requests can observe it
  if (exists(".rgo_stdin", envir = globalenv(), inherits = FALSE)) rm(".rgo_stdin", envir = globalenv())
  reply_con <- stdout()

  send <- function(x) {
    bytes <- serialize(x, NULL, xdr = TRUE, version = 3L)
    cat(paste(as.character(bytes), collapse = ""), "\n", sep = "", file = reply_con)
    flush(reply_con)
  }

  failure <- function(cond) {
    list(ok = FALSE, class = class(cond), message = conditionMessage(cond))
  }

  # Run expr with output and warnings captured; never lets a condition out.
  guarded <- function(expr) {
    warnings <- character()
    out <- NULL
    res <- tryCatch(
      withCallingHandlers(
        {
          out <- utils::capture.output(value <- expr)
          list(ok = TRUE, value = value)
        },
        warning = function(w) {
          warnings <<- c(warnings, conditionMessage(w))
          invokeRestart("muffleWarning")
        }
      ),
      error = failure
    )
    res$output <- out
    res$warnings <- warnings
    res
  }

  hx <- function(s) {
    if (is.na(s)) "NA" else paste(as.character(charToRaw(s)), collapse = "")
  }
  dbits <- function(x) {
    paste(as.character(writeBin(x, raw(), endian = "big")), collapse = "")
  }

  # dump_tree(expression vector): "top N", then each expression, one node per
  # line, children indented by one "." per level. Strings and names are hex
  # (UTF-8 bytes), doubles are their IEEE bits, so that nothing depends on a
  # printing rule. A named argument is preceded by a "name" line.
  dump_tree <- function(exprs) {
    out <- character()
    emit <- function(d, s) out[[length(out) + 1L]] <<- paste0(strrep(".", d), s)
    node <- function(lst, i, d) {
      el <- lst[i]
      if (identical(el[[1L]], quote(expr = ))) return(emit(d, "missing"))
      x <- el[[1L]]
      if (is.call(x)) {
        xl <- as.list(x)
        if (identical(xl[[1L]], as.name("function"))) {
          fl <- as.list(xl[[2L]])
          emit(d, "fn")
          for (k in seq_along(fl)) {
            emit(d + 1L, paste0("formal ", hx(names(fl)[k])))
            if (identical(fl[k][[1L]], quote(expr = ))) emit(d + 2L, "nodefault")
            else node(fl, k, d + 2L)
          }
          node(xl, 3L, d + 1L)
        } else {
          emit(d, "call")
          nms <- names(xl)
          for (k in seq_along(xl)) {
            if (k > 1L && !is.null(nms) && nzchar(nms[k])) emit(d + 1L, paste0("name ", hx(nms[k])))
            node(xl, k, d + 1L)
          }
        }
      } else if (is.name(x)) emit(d, paste0("sym ", hx(as.character(x))))
      else if (is.null(x)) emit(d, "null")
      else if (is.logical(x)) emit(d, paste0("lgl ", if (is.na(x)) "NA" else if (x) "T" else "F"))
      else if (is.integer(x)) emit(d, paste0("int ", if (is.na(x)) "NA" else x))
      else if (is.double(x)) emit(d, paste0("dbl ", dbits(x)))
      else if (is.complex(x)) emit(d, paste0("cplx ", dbits(Re(x)), " ", dbits(Im(x))))
      else if (is.character(x)) emit(d, paste0("str ", hx(x)))
      else emit(d, paste0("other ", typeof(x)))
    }
    emit(0L, paste0("top ", length(exprs)))
    el <- as.list(exprs)
    for (i in seq_along(el)) node(el, i, 1L)
    out
  }

  ops <- list(
    # eval(source): evaluate every expression of the text in a fresh child
    # of the global environment; the value is that of the last one.
    eval = function(args) {
      guarded(eval(parse(text = args[[1L]], keep.source = FALSE),
                   new.env(parent = globalenv())))
    },
    # deparse(source): parse, then deparse each top-level expression; the
    # value is a list of character vectors, one per expression.
    deparse = function(args) {
      guarded(lapply(parse(text = args[[1L]], keep.source = FALSE), deparse))
    },
    # parsedata(source): getParseData() of the parse with source references.
    parsedata = function(args) {
      guarded({
        p <- parse(text = args[[1L]], keep.source = TRUE)
        utils::getParseData(p, includeText = NA)
      })
    },
    # describe(source): evaluate one expression; the value is a list holding
    # the bytes serialize() writes for the result and what R says about it.
    describe = function(args) {
      guarded({
        v <- eval(parse(text = args[[1L]], keep.source = FALSE),
                  new.env(parent = globalenv()))
        list(bytes = serialize(v, NULL, xdr = TRUE, version = 3L),
             typeof = typeof(v), length = length(v), class = class(v))
      })
    },
    # reserialize(bytes): unserialize and serialize again; the value is the
    # raw vector R writes for the object it read.
    reserialize = function(args) {
      guarded(serialize(unserialize(args[[1L]]), NULL, xdr = TRUE, version = 3L))
    },
    # tree(source): the parse as a canonical one-node-per-line dump (see
    # dump_tree); the value is a character vector of lines.
    tree = function(args) {
      guarded(dump_tree(parse(text = args[[1L]], keep.source = FALSE)))
    },
    # tokens(source): the terminal rows of getParseData(), one string per
    # token in source order: "line1:col1-line2:col2 TOKEN hex(text)".
    tokens = function(args) {
      guarded({
        p <- parse(text = args[[1L]], keep.source = TRUE)
        d <- utils::getParseData(p, includeText = NA)
        if (is.null(d)) {
          character()
        } else {
          d <- d[d$terminal, , drop = FALSE]
          d <- d[order(d$line1, d$col1), , drop = FALSE]
          sprintf("%d:%d-%d:%d %s %s", d$line1, d$col1, d$line2, d$col2,
                  d$token, vapply(d$text, hx, ""))
        }
      })
    },
    # parseerror(source): the message of the error parse() raises, or NULL.
    parseerror = function(args) {
      guarded(tryCatch({
        parse(text = args[[1L]], keep.source = FALSE)
        NULL
      }, error = function(e) conditionMessage(e)))
    },
    # benchparse(texts): seconds R spends in parse(text=, keep.source=FALSE) on
    # each text (a character vector), by elapsed time.
    benchparse = function(args) {
      guarded({
        texts <- args[[1L]]
        t0 <- proc.time()[["elapsed"]]
        for (s in texts) parse(text = s, keep.source = FALSE)
        proc.time()[["elapsed"]] - t0
      })
    },
    # sleep(seconds): for the timeout test.
    sleep = function(args) {
      guarded(Sys.sleep(as.numeric(args[[1L]])))
    }
  )

  send(list(
    ok = TRUE, ready = TRUE,
    version = paste(R.version$major, R.version$minor, sep = "."),
    svn = R.version$`svn rev`,
    platform = R.version$platform,
    banner = R.version.string
  ))

  repeat {
    n <- readBin(stdin_con, "integer", 1L, size = 4L, endian = "big")
    if (length(n) == 0L) break
    req <- unserialize(readBin(stdin_con, "raw", n))
    op <- req[[1L]]
    handler <- ops[[op]]
    if (is.null(handler)) {
      send(list(ok = FALSE, class = "oracleError",
                message = paste0("unknown operation '", op, "'")))
    } else {
      send(handler(req[-1L]))
    }
  }
})
