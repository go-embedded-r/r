.deparseOpts(c("keepNA", "keepInteger"))
#---
.deparseOpts("all")
#---
.deparseOpts("exact")
#---
.deparseOpts(NULL)
#---
.deparseOpts(c("keepNA", "keepInteger", "niceNames", "showAttributes"))
#---
.deparseOpts("digits17")
#---
.deparseOpts(c("quoteExpressions", "hexNumeric"))
#---
cat(tryCatch(.deparseOpts("bogus"), error = function(e) conditionMessage(e)), "\n")
#---
x <- c(a = 1L, b = NA); deparse(x, control = NULL)
#---
x <- c(a = 1L, b = NA); deparse(x, control = "keepInteger")
#---
x <- c(a = 1L, b = NA); deparse(x, control = c("keepInteger", "keepNA"))
#---
x <- c(a = 1L, b = NA); deparse(x, control = "niceNames")
#---
x <- c(a = 1L, b = NA); deparse(x, control = c("niceNames", "keepInteger"))
#---
y <- c(a = 1.5, b = NA); deparse(y, control = NULL)
#---
y <- c(a = 1.5, b = NA); deparse(y, control = "keepNA")
#---
y <- c(a = 1.5, b = NA); deparse(y, control = c("keepNA", "niceNames"))
#---
z <- structure(list(p = 1L, q = "s"), class = "foo", extra = c(u = 2)); deparse(z, control = NULL)
#---
z <- structure(list(p = 1L, q = "s"), class = "foo", extra = c(u = 2)); deparse(z, control = "showAttributes")
#---
z <- structure(list(p = 1L, q = "s"), class = "foo", extra = c(u = 2)); deparse(z, control = c("showAttributes", "niceNames"))
#---
f <- factor(c("a", "b")); deparse(f, control = NULL)
#---
f <- factor(c("a", "b")); deparse(f, control = "showAttributes")
#---
deparse(1:3, control = NULL)
#---
deparse(1:3, control = "keepInteger")
#---
deparse(c(1L, 3L), control = "keepInteger")
#---
deparse(c(1L, 3L), control = NULL)
#---
deparse(1L, control = NULL)
#---
deparse(1L)
#---
deparse(pi, control = "digits17")
#---
deparse(pi)
#---
deparse(0.1, control = "digits17")
#---
deparse(1/3, control = "digits17")
#---
deparse(1/3, control = "exact")
#---
deparse(0.5, control = "hexNumeric")
#---
deparse(quote(x), control = "quoteExpressions")
#---
deparse(quote(x + y), control = "quoteExpressions")
#---
deparse(quote(x + y), control = NULL)
#---
deparse(expression(a, b + 1), control = "quoteExpressions")
#---
deparse(letters[1:12], width.cutoff = 20)
#---
deparse(letters[1:12], width.cutoff = 500)
#---
deparse(letters[1:12], width.cutoff = 10)
#---
deparse(letters[1:12], width.cutoff = 60)
#---
deparse(letters[1:12], nlines = 2, width.cutoff = 20)
#---
deparse(c(a = 1, b = 2), control = "all")
#---
deparse(1:10, control = "all")
#---
deparse(c(1.5, NA), control = "all")
#---
deparse(NA_character_)
#---
deparse(NA_character_, control = NULL)
#---
deparse(c("a", NA), control = NULL)
#---
deparse(c("a", NA), control = "keepNA")
#---
deparse(NA, control = NULL)
#---
deparse(NA_integer_, control = NULL)
#---
deparse(NA_integer_, control = "keepNA")
#---
deparse(NA_real_, control = NULL)
#---
deparse(list(a = 1, 2), control = NULL)
#---
deparse(list(a = 1, 2), control = "niceNames")
#---
deparse(list(a = 1, 2), control = "showAttributes")
#---
deparse(c(a = 1), control = "showAttributes")
#---
deparse(c(a = 1), control = NULL)
#---
deparse(as.Date("2024-01-02"), control = NULL)
#---
deparse(as.Date("2024-01-02"))
#---
deparse(as.Date("2024-01-02"), control = "all")
#---
deparse(`names<-`(1:2, c("a", NA)))
#---
deparse(matrix(1:4, 2), control = NULL)
#---
deparse(matrix(1:4, 2))
#---
deparse(TRUE, control = "all")
#---
deparse(-1L, control = "all")
#---
deparse(c(1, 2, 3), control = "all")
#---
deparse(1e5, control = NULL)
#---
deparse(100000L)
#---
x <- c(a = 1L, b = NA); deparse(x, backtick = TRUE)
#---
deparse(quote(`a b` + 1), backtick = FALSE)
#---
deparse(quote(`a b` + 1))
#---
deparse(quote(`a b`), backtick = FALSE)
#---
deparse(quote(`a b`), backtick = TRUE)
#---
deparse(c(`a b` = 1))
#---
deparse(structure(1, `a b` = 2))
#---
tryCatch(deparse(1, width.cutoff = 5), warning = function(w) conditionMessage(w))
#---
tryCatch(deparse(1, width.cutoff = 600), warning = function(w) conditionMessage(w))
#---
cat(tryCatch(deparse(1, control = "bogus"), error = function(e) conditionMessage(e)), "\n")
#---
x <- c(1.5, 2, 3); y <- 1:3; dump(c("x", "y"), file = "")
#---
f <- function(a, b = 2) { a + b }; dump("f", file = "")
#---
ff <- factor(c("u", "v")); dump("ff", file = "")
#---
`a b` <- "s"; dump("a b", file = "")
#---
l <- list(p = 1L, q = "t"); dump("l", file = "")
#---
x <- c(1.5, 2, 3); r <- dump(c("x", "nosuch"), file = ""); r
#---
x <- c(1.5, 2, 3); suppressWarnings(dump(c("x", "nosuch"), file = ""))
#---
x <- c(1.5, 2, 3); tryCatch(dump(c("x", "nosuch"), file = ""), warning = function(w) conditionMessage(w))
#---
x <- c(1.5, 2, 3); dump(c("x"), file = "", control = NULL)
#---
y <- 1:3; dump(c("y"), file = "", control = "keepInteger")
#---
x <- c(1.5, 2, 3); y <- 1:3; dump(c("x", "y"), file = "", append = TRUE)
#---
f <- tempfile(); x <- c(1.5, 2, 3); dump("x", file = f); readLines(f)
#---
f <- tempfile(); x <- c(1.5, 2, 3); y <- 1:3; dump(c("x", "y"), f, append = FALSE); readLines(f)
#---
f <- tempfile(); x <- c(1.5, 2, 3); y <- 1:3; dump(c("x", "y"), f); dump("x", f, append = TRUE); readLines(f)
#---
f <- tempfile(); x <- c(1.5, 2, 3); dump("x", f); dget(f)
#---
x <- c(1.5, 2, 3); dput(x)
#---
f <- function(a, b = 2) { a + b }; dput(f)
#---
ff <- factor(c("u", "v")); dput(ff)
#---
l <- list(p = 1L, q = "t"); dput(l)
#---
dput(1:10)
#---
dput("a\"b")
#---
dput(NULL)
#---
dput(c(a = 1))
#---
x <- c(1.5, 2, 3); dput(x, control = NULL)
#---
y <- 1:3; dput(y, control = "all")
#---
x <- c(1.5, 2, 3); r2 <- dput(x); identical(r2, x)
#---
f <- tempfile(); x <- c(1.5, 2, 3); dput(x, file = f); readLines(f)
#---
f <- tempfile(); x <- c(1.5, 2, 3); dput(x, file = f); dget(f)
#---
dput(as.Date("2024-01-02"))
#---
dput(letters)
#---
dput(list(1, list(2, "a")))
#---
dput(quote(x + y))
#---
dput(expression(a, b))
#---
dput(1e-15)
#---
dput(c(0.1, 1/3))
#---
dput(TRUE)
#---
dput(NA_character_)
#---
dput(data.frame(a = 1:2, b = c("x", "y")))
#---
x <- c(1.5, 2, 3); dump("x", file = stdout())
#---
x <- c(1.5, 2, 3); y <- 1:3; dump(list = c("x", "y"), file = "")
#---
f <- function(a, b = 2) { a + b }; dump("f", file = "", control = NULL)
#---
cat(tryCatch(dump(1, file = ""), error = function(e) conditionMessage(e)), "\n")
#---
tryCatch(dget("nosuchfile.R"), error = function(e) "err", warning = function(w) "warn")
#---
dput(structure(list(sec = 45, min = 30L, hour = 10L, mday = 31L, mon = 6L, year = 124L, wday = 3L, yday = 212L, isdst = 1L, zone = "CEST", gmtoff = 7200L), class = c("POSIXlt", "POSIXt"), tzone = c("Europe/Paris", "CET", "CEST"), balanced = TRUE))
#---
deparse(structure(list(p = 1, q = "s"), class = "foo", extra = c(2, 3)))
#---
deparse(structure(list(p = 1, q = "s"), class = "foocccccc", extra = c(2, 3)))
#---
deparse(structure(list(p = 1, q = "s"), class = "foocccccc", extra = c(u = 2)))
#---
deparse(c(aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa = 1))
#---
deparse(c(aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa = 1, 2))
#---
deparse(c(aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa = 1, b = 2))
#---
deparse(list(aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa = 1, b = 2))
#---
deparse(c(1.5, 2, 3.25, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28))
#---
deparse(c(aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa = 1, bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb = 2, uuuu = 3))
