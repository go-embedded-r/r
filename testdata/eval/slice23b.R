w <- structure(list(quote(f(x)), quote(f(x)), NULL, quote(g(y))), names = c("a", "a", "b", "a"), class = "warnings")
s <- summary(w)
print(class(s)); print(attr(s, "counts")); print(names(attributes(s))); print(length(s))
#---
w <- structure(list(quote(f(x)), quote(f(x)), NULL, quote(g(y))), names = c("a", "a", "b", "a"), class = "warnings")
print(summary(w))
#---
w <- structure(list(quote(f(x)), quote(f(x))), names = c("a", "a"), class = "warnings")
print(summary(w))
#---
w <- structure(list(NULL), names = "only", class = "warnings")
print(summary(w))
#---
w <- structure(list(NULL, NULL, NULL), names = c("m", "m", "m"), class = "warnings")
print(summary(w))
#---
w <- structure(list(quote(f(x))), names = "a", class = "warnings")
print(summary(w[0]))
#---
w <- structure(list(quote(h(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20)), quote(f(x))), names = c("a", "b"), class = "warnings")
print(summary(w))
#---
w <- structure(list(quote(f(x)), quote(f(x))), names = c("a long message that wraps past the line width of the console", "a long message that wraps past the line width of the console"), class = "warnings")
print(summary(w))
#---
w <- structure(list(quote(f(x)), quote(f(x)), quote(f(x))), names = c("one", "two", "one"), class = "warnings")
x <- print(summary(w)); print(class(x))
#---
print(.gt(c(3, 1, 2), 1, 2)); print(.gt(c(3, 1, 2), 2, 3)); print(.gt(c(1, 1), 1, 2)); print(.gt(c("b", "a"), 1, 2))
#---
print(tryCatch(.gt(c(NA, 1), 1, 2), error = conditionMessage))
#---
print(class(.gt(c(3, 1, 2), 1, 2)))
#---
print(.gtn(c(1, 2, 3), TRUE)); print(.gtn(c(1, 2, 2), TRUE)); print(.gtn(c(1, 2, 2), FALSE)); print(.gtn(c(3, 2, 1), FALSE))
#---
print(.gtn(c(1, NA, 3), FALSE)); print(.gtn(numeric(0), TRUE)); print(.gtn(1, TRUE)); print(.gtn(c("a", "b"), TRUE))
#---
g <- function(...) .rangeNum(..., isNumeric = is.numeric)
print(g(c(3, 1, 2), na.rm = FALSE, finite = FALSE)); print(g(c(3, NA, Inf), na.rm = TRUE, finite = FALSE))
print(g(c(3, NA, Inf), na.rm = FALSE, finite = TRUE)); print(g(c(3, NA), na.rm = FALSE, finite = FALSE))
#---
g <- function(...) .rangeNum(..., isNumeric = is.numeric)
print(g(1:3, 10, na.rm = FALSE, finite = FALSE)); print(g("a", na.rm = FALSE, finite = FALSE)); print(g(c(3, 5), "a", na.rm = FALSE, finite = FALSE))
print(g(TRUE, FALSE, na.rm = FALSE, finite = FALSE)); print(g(c(3, NaN), na.rm = FALSE, finite = FALSE)); print(g(c(3, NaN), na.rm = TRUE, finite = FALSE))
#---
print(.rangeNum(c("b", NA), na.rm = TRUE, finite = FALSE, isNumeric = function(x) FALSE))
print(.rangeNum(c("b", NA), na.rm = FALSE, finite = FALSE, isNumeric = function(x) FALSE))
print(.rangeNum(c("b", "a"), "c", na.rm = FALSE, finite = TRUE, isNumeric = function(x) FALSE))
#---
print(tryCatch(.rangeNum(numeric(0), na.rm = FALSE, finite = FALSE, isNumeric = is.numeric), warning = conditionMessage))
#---
print(.format.zeros(c("0.0", "1.5", "0.00"), "-"))
#---
print(.format.zeros(c("0.0", "1.5"), "-", replace = TRUE)); print(.format.zeros(c(a = "0.0", b = "1.5"), "-"))
#---
print(.format.zeros(c("0.0", " 0", " 0.00", "1"), "-")); print(.format.zeros(c("0.0", "-0.0", "NA"), "-"))
#---
print(.format.zeros(c("0e+00", "1e+00"), "-")); print(.format.zeros(character(0), "-")); print(.format.zeros(c("0", "10"), "-"))
#---
print(tryCatch(.format.zeros(c("0.0", "1.5"), "abcdef"), warning = conditionMessage))
#---
print(.format.zeros(c("0.0", "1.5"), "abcdef", warn.non.fitting = FALSE)); print(.format.zeros(c("0.0", "1.5"), "abcdef", replace = TRUE))
#---
print(.format.zeros(c("0.0", "1.5"), "-", nx = c(FALSE, TRUE))); print(.format.zeros(c("0.0", "1.5"), "-", nx = c(TRUE, TRUE)))
print(.format.zeros(c("0.0", "1.5"), "-", nx = c(FALSE, FALSE))); print(.format.zeros(c("0.0", "1.5"), "-", nx = c(NA, FALSE)))
#---
print(tryCatch(.format.zeros(c(0, 1.5), "-"), error = conditionMessage))
#---
print(tryCatch(.format.zeros(c("0.0", "1.5"), NA), error = conditionMessage))
#---
print(.format.zeros(matrix(c("0.0", "1.5", "0.0", "2"), 2), "-"))
#---
print(.doWrap(c(3, 1, 2), FALSE, TRUE, FALSE)); print(.doSortWrap(c(3, 1, NA), TRUE, TRUE, FALSE)); print(.doSortWrap(c(a = 1), FALSE, FALSE, FALSE))
#---
print(.tryResumeInterrupt())
#---
print(tryCatch(.signalSimpleWarning("hello", quote(f(x))), warning = conditionMessage))
#---
r <- withCallingHandlers(.signalSimpleWarning("hello", quote(f(x))),
  warning = function(w) { cat(class(w), "\n"); print(conditionCall(w)); invokeRestart("muffleWarning") })
print(r)
#---
u <- as.POSIXct("2020-01-01", tz = "UTC"); p <- as.POSIXct("2020-01-01", tz = "Europe/Paris")
print(.check_tzones(u)); print(.check_tzones(u, u)); print(.check_tzones(1)); print(.check_tzones()); print(.check_tzones(u, NULL))
#---
u <- as.POSIXct("2020-01-01", tz = "UTC"); p <- as.POSIXct("2020-01-01", tz = "Europe/Paris")
print(tryCatch(.check_tzones(u, p), warning = conditionMessage))
print(tryCatch(.check_tzones(as.POSIXlt(u), as.POSIXlt(p)), warning = conditionMessage))
print(tryCatch(.check_tzones(u, u, p), warning = conditionMessage))
#---
u <- as.POSIXct("2020-01-01", tz = "UTC")
print(.check_tzones(u, as.Date("2020-01-01"))); print(.check_tzones(as.POSIXct("2020-01-01", tz = ""), u))
#---
m <- taskCallbackManager()
print(class(m)); print(names(m))
#---
m <- taskCallbackManager()
print(m$add(function(...) TRUE, name = "first")); print(m$add(function(...) TRUE)); print(names(m$callbacks()))
#---
m <- taskCallbackManager()
m$add(function(expr, value, ok, visible, data) { cat("saw", deparse(expr), "data", data, "\n"); TRUE }, data = 7, name = "h")
print(m$evaluate(quote(1 + 1), 2, TRUE, TRUE))
#---
m <- taskCallbackManager()
m$add(function(expr, value, ok, visible) FALSE, name = "gone")
m$add(function(expr, value, ok, visible) TRUE, name = "stays")
print(m$evaluate(quote(x), 1, TRUE, TRUE)); print(names(m$callbacks()))
#---
m <- taskCallbackManager()
print(tryCatch(m$remove("zzz"), error = conditionMessage))
print(tryCatch(m$remove(9), error = conditionMessage))
#---
m <- taskCallbackManager()
m$add(function(...) TRUE, name = "a"); m$add(function(...) TRUE, name = "b"); m$add(function(...) TRUE, name = "c")
print(m$remove("b")); print(names(m$callbacks())); print(m$remove(1)); print(names(m$callbacks()))
#---
m <- taskCallbackManager()
m$add(function(...) { cat("called\n"); TRUE }, name = "x")
print(m$suspend()); print(m$evaluate(quote(1), 1, TRUE, TRUE)); print(m$suspend(FALSE)); print(m$evaluate(quote(1), 1, TRUE, TRUE))
#---
m <- taskCallbackManager()
print(tryCatch({ m$add(function(...) stop("boom"), name = "bad"); m$evaluate(quote(1), 1, TRUE, TRUE) }, error = conditionMessage))
#---
m <- taskCallbackManager()
m$add(function(...) TRUE, name = "k")
print("R-taskCallbackManager" %in% getTaskCallbackNames())
print(removeTaskCallback("R-taskCallbackManager"))
#---
m <- taskCallbackManager()
invisible(m$register(name = "zz")); print("zz" %in% getTaskCallbackNames())
invisible(removeTaskCallback("zz"))
#---
tc <- textConnection("a"); f <- file(tempfile(), "w")
m <- showConnections()
m <- m[m[, "class"] %in% c("textConnection", "file") & m[, "description"] != "stdin", , drop = FALSE]
print(dim(m)); print(colnames(m)); print(unname(m[, "class"])); print(unname(m[, "mode"])); print(unname(m[, "isopen"]))
print(unname(m[, "can read"])); print(unname(m[, "can write"])); print(unname(m[, "text"])); print(unname(m[1, "description"]))
close(tc); close(f)
#---
tc <- textConnection("a")
m <- showConnections(all = TRUE)
print(rownames(m)[1:3]); print(m["0", ]); print(m["1", "mode"]); print(m["2", "description"])
close(tc)
#---
tc <- textConnection("a"); print(typeof(showConnections())); print(is.matrix(showConnections())); close(tc)
#---
tc <- textConnection("a"); cc <- file(tempfile(), "w")
k <- function() { m <- showConnections(); sum(m[, "class"] %in% c("textConnection", "file") & m[, "description"] != "stdin") }
print(k()); close(tc); print(k()); close(cc); print(k())
#---
x <- textConnection("zz"); y <- textConnection("yy"); close(x)
m <- showConnections(); print(sum(m[, "description"] %in% c("\"zz\"", "\"yy\"")))
close(y)
#---
print(c(1, list(2, list(3, "a")), recursive = TRUE)); print(c(list(a = 1), b = 2, recursive = TRUE))
#---
print(c(a = 1, b = 2, use.names = FALSE)); print(c(list(a = 1, b = list(c = 2)), recursive = TRUE, use.names = FALSE))
#---
print(c(list(1, 2), recursive = FALSE)); print(c(NULL, recursive = TRUE))
