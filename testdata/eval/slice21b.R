pairlist(a = 1, b = "x")
#---
typeof(pairlist(1, 2))
#---
pairlist()
#---
is.null(pairlist())
#---
x <- as.pairlist(list(a = 1, b = 2)); typeof(x)
#---
as.pairlist(1:2)
#---
length(pairlist(1, 2, 3))
#---
x <- c("a1b22", "c333", "d"); m <- regexpr("[0-9]+", x); regmatches(x, m) <- c("X", "Y"); x
#---
x <- c("a1b22", "c333", "d"); m <- gregexpr("[0-9]+", x); regmatches(x, m) <- list(c("<1>", "<2>"), "<3>", character()); x
#---
x <- c("a1b22", "c333"); m <- regexpr("[0-9]+", x); regmatches(x, m, invert = TRUE) <- list(c("p", "q"), c("r", "s")); x
#---
x <- c("ab", "cd"); m <- regexpr("z", x); regmatches(x, m) <- character(); x
#---
gregexec("([a-z])([0-9])", "a1b2")
#---
x <- "a1b2c3"; m <- gregexec("([a-z])([0-9])", x); regmatches(x, m)
#---
x <- c("a1", "zz"); m <- gregexec("([a-z])([0-9])", x); m[[2]]
#---
formatDL(c("a", "bb"), c("first", "second"))
#---
formatDL(c("a", "bb"), c("first", "second"), style = "list")
#---
formatDL(c(x = "term", y = "t2"), c("defn one", "defn two"), width = 20, indent = 6)
#---
formatDL(c("a", "bb", "ccc"))
#---
merge(1:2, c("a", "b"))
#---
merge.default(data.frame(k = 1:3, a = 4:6), data.frame(k = c(1, 3), b = c("u", "v")))
#---
d <- data.frame(a = 1:4, g = c(1, 2, 1, 2)); split(d, d$g) <- list(data.frame(a = 10:11, g = 1), data.frame(a = 20:21, g = 2)); d
#---
d <- data.frame(a = 1:4, g = c(1, 2, 1, 2)); `split<-.data.frame`(d, d$g, list(data.frame(a = 0L, g = 1), data.frame(a = 9L, g = 2)))
#---
anyNA(data.frame(a = c(1, NA), b = 1:2))
#---
anyNA(data.frame(a = c(1, 2), b = 1:2))
#---
anyNA(data.frame(a = c(1, NA), b = 1:2), recursive = TRUE)
#---
anyDuplicated(data.frame(a = c(1, 2, 1), b = c("x", "y", "x")))
#---
anyDuplicated(data.frame(a = c(1, 2, 1), b = c("x", "y", "x")), fromLast = TRUE)
#---
anyDuplicated(data.frame(a = 1:3))
#---
as.vector(data.frame(a = 1:2, b = c("x", "y")))
#---
as.vector(data.frame(a = 1:2, b = 3:4), "list")
#---
as.vector(data.frame(a = 1:2, b = 3:4), "character")
#---
all.equal(as.POSIXct("2024-01-02 03:04:05", tz = "UTC"), as.POSIXct("2024-01-02 03:04:06", tz = "UTC"))
#---
all.equal(as.POSIXct("2024-01-02 03:04:05", tz = "UTC"), as.POSIXct("2024-01-02 03:04:05.0001", tz = "UTC"))
#---
all.equal(as.POSIXct("2024-01-02", tz = "UTC"), as.POSIXct("2024-01-03", tz = "UTC"), tolerance = 0.01)
#---
all.equal(as.POSIXct("2024-01-02", tz = "UTC"), "2024")
#---
all.equal(as.POSIXlt("2024-01-02 03:04:05", tz = "UTC"), as.POSIXlt("2024-01-02 03:04:05", tz = "UTC"))
#---
all.equal(y ~ x, y ~ x)
#---
all.equal(y ~ x, y ~ z)
#---
all.equal(y ~ x, 3)
#---
all.equal(y ~ x + z, y ~ x + w)
#---
cond <- messageCondition("hello", call = NULL); class(cond)
#---
cond <- messageCondition("hello\n", call = quote(f(x))); conditionMessage(cond); conditionCall(cond)
#---
cond <- messageCondition("hi", class = "myMsg", extra = 1); class(cond); cond$extra
#---
tryCatch(packageStartupMessage("loading"), packageStartupMessage = function(m) class(m))
#---
packageStartupMessage("quiet ", "please")
#---
suppressPackageStartupMessages({ packageStartupMessage("x"); 5 })
#---
suppressPackageStartupMessages({ message("kept"); 6 })
#---
withRestarts(isRestart(computeRestarts()[[1]]), foo = function() 1)
#---
isRestart(1)
#---
withRestarts(restartDescription(computeRestarts()[[1]]), foo = list(handler = function(x, y) 1, description = "does foo"))
#---
withRestarts(restartFormals(computeRestarts()[[1]]), foo = function(x, y = 2) 1)
#---
withRestarts(restartDescription(computeRestarts()[[1]]), foo = function() 1)
#---
withAutoprint({ x <- 1:3; x * 2; y <- "a" })
#---
withAutoprint({ x <- 5; invisible(x); x + 1 }, echo = FALSE)
#---
withAutoprint(expression(1 + 1, "b"))
#---
withAutoprint({ a <- 1 + 1; a }, prompt.echo = "R> ", continue.echo = "R+ ")
#---
traceback()
#---
.traceback()
#---
x <- NULL; setHook("myhook", function(...) x <<- c(x, "one")); setHook("myhook", function(...) x <<- c(x, "two")); for (h in getHook("myhook")) h(); x
#---
setHook("h2", function(...) 1); setHook("h2", function(...) 2, "prepend"); length(getHook("h2"))
#---
setHook("h3", function(...) 1); setHook("h3", NULL, "replace"); length(getHook("h3"))
#---
length(getHook("nothing"))
#---
setHook("h4", "notafunction")
#---
id <- addTaskCallback(function(...) TRUE, name = "mine"); names(id); getTaskCallbackNames(); removeTaskCallback("mine"); getTaskCallbackNames()
#---
id <- addTaskCallback(function(...) TRUE); id2 <- addTaskCallback(function(...) TRUE); as.integer(id2) - as.integer(id); removeTaskCallback(id); removeTaskCallback(id2); length(getTaskCallbackNames())
#---
removeTaskCallback("none")
#---
id <- addTaskCallback(function(...) TRUE, name = "q"); removeTaskCallback("q"); removeTaskCallback("q")
#---
addTaskCallback("notafunction")
