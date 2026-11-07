## Conditions: stop, warning, message, tryCatch, withCallingHandlers, on.exit. A line holding only #--- separates cells.
#---
stop("boom")
#---
f <- function() stop("inside f"); f()
#---
f <- function() stop("a", "b", 3); f()
#---
f <- function() stop("no call", call. = FALSE); f()
#---
tryCatch(stop("x"), error = function(e) conditionMessage(e))
#---
tryCatch(stop("x"), error = function(e) "caught")
#---
tryCatch(1 + 1, error = function(e) "caught")
#---
f <- function() stop("deep"); tryCatch(f(), error = function(e) conditionCall(e))
#---
tryCatch(stop("x"), error = function(e) class(e))
#---
tryCatch(warning("w"), warning = function(w) class(w))
#---
tryCatch(warning("w"), warning = function(w) conditionMessage(w))
#---
tryCatch(message("m"), message = function(m) conditionMessage(m))
#---
tryCatch(signalCondition(simpleCondition("sc")), condition = function(c) "got it")
#---
tryCatch(1 + "a", error = function(e) conditionMessage(e))
#---
tryCatch(undefined_var, error = function(e) conditionMessage(e))
#---
tryCatch(undefined_fn(), error = function(e) conditionMessage(e))
#---
tryCatch(stop("x"), error = function(e) "first", error = function(e) "second")
#---
tryCatch(stop("x"), condition = function(c) "condition first", error = function(e) "error second")
#---
tryCatch({ cat("a\n"); stop("x"); cat("b\n") }, error = function(e) cat("handler\n"))
#---
tryCatch(tryCatch(stop("inner"), error = function(e) stop("outer")), error = function(e) conditionMessage(e))
#---
tryCatch(tryCatch(stop("inner"), warning = function(w) "w"), error = function(e) "outer caught")
#---
tryCatch(stop("x"), error = function(e) tryCatch(stop("y"), error = function(e2) paste(conditionMessage(e), conditionMessage(e2))))
#---
f <- function() { for (i in 1:3) tryCatch({ if (i == 2) next; cat(i, "\n") }, error = function(e) NULL) }; f()
#---
r <- tryCatch(stop("x"), error = function(e) NULL); is.null(r)
#---
res <- withCallingHandlers(1 + 1, error = function(e) "never"); res
#---
withCallingHandlers({ warning("w1"); "value" }, warning = function(w) { cat("saw", conditionMessage(w), "\n"); invokeRestart("muffleWarning") })
#---
withCallingHandlers({ warning("w1"); "value" }, warning = function(w) cat("saw but not muffled\n"))
#---
withCallingHandlers({ message("m1"); "value" }, message = function(m) { cat("got", conditionMessage(m)); invokeRestart("muffleMessage") })
#---
tryCatch(withCallingHandlers(stop("e"), error = function(e) cat("first\n")), error = function(e) cat("then exiting\n"))
#---
withCallingHandlers({ warning("a"); warning("b"); "done" }, warning = function(w) { cat(conditionMessage(w), "\n"); invokeRestart("muffleWarning") })
#---
f <- function() { warning("fw"); "returned" }; f()
#---
f <- function() { warning("fw"); "returned" }; suppressWarnings(f())
#---
suppressWarnings({ warning("a"); 1 + 1 })
#---
suppressMessages({ message("hidden"); "shown" })
#---
message("to stderr")
#---
message("a", "b"); message("c")
#---
f <- function() message("from f"); f()
#---
message("no newline", appendLF = FALSE)
#---
tryCatch(message("x"), message = function(m) class(m))
#---
x <- warning("returned"); x
#---
warning("top-level warning")
#---
f <- function() warning("in f"); f()
#---
f <- function() warning("no call", call. = FALSE); f()
#---
1:3 + 1:2
#---
sqrt(-1)
#---
as.integer("x")
#---
withCallingHandlers(as.integer("x"), warning = function(w) { cat("W:", conditionMessage(w), "\n"); invokeRestart("muffleWarning") })
#---
tryCatch(as.integer("x"), warning = function(w) "warned")
#---
tryCatch(sqrt(-1), warning = function(w) conditionMessage(w))
#---
cnd <- simpleCondition("custom msg"); class(cnd) <- c("myCondition", "condition"); tryCatch(signalCondition(cnd), myCondition = function(c) "custom caught")
#---
cnd <- structure(class = c("myError", "error", "condition"), list(message = "my error", call = NULL)); tryCatch(stop(cnd), myError = function(e) conditionMessage(e))
#---
cnd <- structure(class = c("myError", "error", "condition"), list(message = "my error", call = NULL)); stop(cnd)
#---
cnd <- structure(class = c("myWarning", "warning", "condition"), list(message = "my warning", call = NULL)); tryCatch(warning(cnd), myWarning = function(w) "custom warning caught")
#---
cnd <- structure(class = c("myWarning", "warning", "condition"), list(message = "my warning", call = NULL)); warning(cnd)
#---
e <- simpleError("E"); class(e); conditionMessage(e); is.null(conditionCall(e))
#---
e <- simpleError("E", quote(f(x))); conditionCall(e)
#---
w <- simpleWarning("W"); class(w)
#---
m <- simpleMessage("M\n"); class(m)
#---
c1 <- simpleCondition("sc"); class(c1)
#---
tryCatch(stop(simpleError("from object")), error = function(e) conditionMessage(e))
#---
signalCondition(simpleCondition("nobody listens"))
#---
f <- function() { on.exit(cat("bye\n")); cat("hi\n"); invisible(1) }; f()
#---
f <- function() { on.exit(cat("one\n")); on.exit(cat("two\n")); 1 }; f()
#---
f <- function() { on.exit(cat("one\n")); on.exit(cat("two\n"), add = TRUE); 1 }; f()
#---
f <- function() { on.exit(cat("one\n")); on.exit(cat("two\n"), add = TRUE, after = FALSE); 1 }; f()
#---
f <- function() { on.exit(cat("cleanup\n")); stop("fail") }; tryCatch(f(), error = function(e) "caught")
#---
f <- function() { on.exit(cat("cleanup\n")); stop("fail") }; f()
#---
f <- function() { x <- 1; on.exit(x <- 2); x }; f()
#---
f <- function() { on.exit(cat("a\n")); on.exit(); 1 }; f()
#---
f <- function() { g <- function() { on.exit(cat("g exit\n")); stop("g fails") }; tryCatch(g(), error = function(e) cat("caught in f\n")); on.exit(cat("f exit\n")); 1 }; f()
#---
f <- function() { on.exit(cat("exit\n")); return("value") }; f()
#---
try(stop("tried"), silent = TRUE); "continued"
#---
x <- try(1 + 1); x
#---
inherits(try(stop("x"), silent = TRUE), "try-error")
#---
stop()
#---
f <- function() stop(); f()
#---
tryCatch(stop(), error = function(e) conditionMessage(e))
#---
warning()
#---
tryCatch(warning(), warning = function(w) conditionMessage(w))
#---
f <- function(x) { if (x > 0) stop("positive") ; "ok" }; sapply_like <- function(v) { r <- character(0); for (i in v) r <- c(r, tryCatch(f(i), error = function(e) conditionMessage(e))); r }; sapply_like(c(-1, 1))
#---
h <- function(e) "H"; tryCatch(stop("x"), error = h)
#---
tryCatch(stop("x"), error = function(e) invisible("quiet"))
#---
tryCatch(invisible(5))
#---
tryCatch(NULL)
#---
tryCatch(expr = 1 + 2)
#---
tryCatch()
#---
withCallingHandlers()
#---
withCallingHandlers(expr = "e")
#---
res <- NULL; withCallingHandlers({ warning("W"); res <- "continued" }, warning = function(w) invokeRestart("muffleWarning")); res
#---
invokeRestart("muffleWarning")
#---
invokeRestart("nosuchrestart")
#---
f <- function() { warning("first"); warning("second"); 3 }; r <- withCallingHandlers(f(), warning = function(w) invokeRestart("muffleWarning")); r
#---
tryCatch({ warning("w"); "not reached" }, warning = function(w) "stopped at warning")
#---
e <- tryCatch(stop("E"), error = identity_fn <- function(e) e); conditionMessage(e); inherits(e, "error")
#---
cat("a", "b", sep = ""); cat("\n"); cat("x", 1, TRUE, NULL, NA, "\n"); cat(1.5, 1e5, 100000L, 0.1 + 0.2, 1/3, "\n"); cat("a\n", "b\n"); cat(c(1, 2, 3), sep = ", "); cat("\n")
#---
cat(sum)
#---
cat("tab\there\n"); cat('quote"s\n')
#---
print("x"); print(1:3); print(NULL); invisible(print("p"))
#---
x <- print("returned"); x

#---
f <- function() { on.exit(cat("a\n"), TRUE); on.exit(cat("b\n"), add = TRUE, after = FALSE); 1 }; f()
#---
f <- function() { on.exit(, cat()) }; f()
