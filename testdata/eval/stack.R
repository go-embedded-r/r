f <- function() sys.call(); f()
#---
f <- function() { x <- 1; g() }; g <- function() sys.frame(-1)$x; f()
#---
f <- function() sys.frame(0); identical(f(), globalenv())
#---
f <- function(n) if (n <= 1) 1 else n * Recall(n - 1); f(5); sapply(1:4, function(i) if (i <= 1) 1 else i * Recall(i - 1))
#---
f <- function(x) sys.function(); is.function(f(1))
#---
f <- function() on.exit(cat("bye\n")); f(); g <- function() { on.exit(cat("1\n")); on.exit(cat("2\n"), add = TRUE); sys.on.exit() }; g()
#---
tryCatch(stop("a"), error = function(e) conditionMessage(e)); tryCatch(stop("a", call. = FALSE), error = function(e) is.null(conditionCall(e)))
#---
f <- function() stop("inside"); tryCatch(f(), error = function(e) deparse(conditionCall(e)))
#---
warning("w1", call. = FALSE)
#---
f <- function() warning("w2"); f()
#---
withCallingHandlers({ warning("x"); "done" }, warning = function(w) { message("caught: ", conditionMessage(w)); invokeRestart("muffleWarning") })
#---
r <- withRestarts(invokeRestart("myR", 5), myR = function(v) v * 2); r
#---
withRestarts(computeRestarts()[[1]]$name, aa = function() 1)
#---
tryInvokeRestart("nonexistent"); 1
#---
invokeRestart("nonexistent")
#---
cond <- simpleCondition("custom message"); class(cond) <- c("custom", "condition"); withCallingHandlers(signalCondition(cond), custom = function(c) cat("got custom\n")); tryCatch(signalCondition(cond), condition = function(c) class(c))
#---
e <- errorCondition("bad", data = 42, class = "myError"); class(e); e$data; tryCatch(stop(e), myError = function(err) err$data)
#---
w <- warningCondition("careful", class = "myWarning"); class(w); tryCatch(warning(w), myWarning = function(x) "caught")
#---
tryCatch({ warning("first"); "not reached" }, warning = function(w) "W", error = function(e) "E")
#---
res <- try(stop("oops"), silent = TRUE); class(res); inherits(res, "try-error"); conditionMessage(attr(res, "condition"))
#---
f <- function() { warning("a"); warning("b"); 10 }; withCallingHandlers(f(), warning = function(w) { cat("W:", conditionMessage(w), "\n"); invokeRestart("muffleWarning") })
#---
tryCatch(as.integer("x"), warning = function(w) conditionMessage(w))
#---
suppressWarnings({ warning("x"); 5 }); suppressMessages({ message("m"); 6 })
#---
tryCatch(message("hello"), message = function(m) conditionMessage(m))
#---
conditionMessage(simpleError("E")); conditionCall(simpleError("E", quote(f(x)))); class(simpleWarning("W")); class(simpleMessage("M")); inherits(simpleCondition("c"), "condition")
#---
stop()
#---
stop("a", "b", 3)
#---
warning("a", "b", 3)
#---
f <- function() { tryCatch(return("early"), finally = cat("cleanup\n")); "late" }; f()
#---
tryCatch(stop("x"), error = function(e) "handled", finally = cat("finally\n"))
#---
nested <- function() tryCatch(stop("inner"), error = function(e) stop("outer")); tryCatch(nested(), error = function(e) conditionMessage(e))
#---
signalCondition(simpleCondition("ignored")); 1
#---
withCallingHandlers(
  { message("note"); 1 },
  message = function(m) { cat("msg:", conditionMessage(m)); invokeRestart("muffleMessage") })
#---
stopifnot(is.numeric(1)); tryCatch(stopifnot("x must be positive" = -1 > 0), error = function(e) conditionMessage(e))
#---
f <- function(x) { if (x > 0) stop("positive: ", x); "ok" }; sapply(c(-1, 0), f); tryCatch(sapply(c(-1, 2), f), error = function(e) conditionMessage(e))

#---
n0 <- sys.nframe(); f <- function(x) sys.nframe() - n0; f(1); g <- function() f(1); g()
#---
n0 <- length(sys.calls()); f <- function() g(); g <- function() { cs <- sys.calls(); list(length(cs) - n0, deparse(cs[[length(cs)]]), deparse(cs[[length(cs) - 1]])) }; f()
#---
n0 <- length(sys.frames()); f <- function() g(); g <- function() length(sys.frames()) - n0; f()

#---
f <- function(x, ...) deparse(match.call()); f(1, b = 2); f(x = 3); f(1, 2, 3)
#---
f <- function(x, y) deparse(match.call()); f(y = 1, 2); g <- function(...) f(...); g(1, 2)
#---
f <- function(a, bb) deparse(match.call()); f(b = 2, a = 1)
#---
f <- function(x, y = 2, ...) deparse(match.call()); f(1); f(); f(y = 3, 4, z = 5)
#---
withRestarts(unlist(lapply(computeRestarts()[1:2], function(r) r$name)), a = function() 1, b = function() 2)
#---
r <- withRestarts(findRestart("zz")$name, zz = function() 1); r; findRestart("nope_not_there")
#---
tryCatch(stop("x"), error = function(e) "handled", finally = cat("finally\n")); 2
#---
tryCatch(1, finally = cat("f\n"))
#---
res <- try(stop("zzz"), silent = TRUE); is.character(geterrmessage()); grepl("zzz", geterrmessage())

#---
f <- function(x, y) deparse(match.call(expand.dots = FALSE)); f(1, 2)
#---
f <- function(x, ...) { m <- match.call(expand.dots = FALSE); length(m$...) }; f(1, 2, 3)

#---
m <- quote(f(a = 1, bb = 2, 3)); m$a; m$bb; m$b; m$zz; m$c
