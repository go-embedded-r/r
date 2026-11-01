# The oracle: a request/reply loop around GNU R, written for this repository
# (BSD-3-Clause, no code from R). It is run by internal/oracle as
#
#     Rscript --vanilla -e <this file>
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
  stdin_con <- file("stdin", "rb")
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
