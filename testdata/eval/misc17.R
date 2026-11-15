rowsum(c(1, 2, 3, 4), c("a", "b", "a", "b"))
#---
rowsum(matrix(1:6, 3), c(2, 1, 2))
#---
rowsum(data.frame(x = 1:4, y = c(2.5, 1, 1, 1)), c("b", "a", "b", "a"))
#---
rowsum(c(1, NA, 3), c("a", "a", "b"))
#---
rowsum(c(1, NA, 3), c("a", "a", "b"), na.rm = TRUE)
#---
rowsum(c(1, 2, 3), c("b", "a", "b"), reorder = FALSE)
#---
rowsum(1:4, factor(c("x", "y", "x", "y"), levels = c("y", "x")))
#---
rowsum(matrix(1:4, 2, dimnames = list(NULL, c("p", "q"))), c("g", "g"))
#---
rowsum(c(1, 2), c(NA, "a"))
#---
tryCatch(rowsum(1:3, 1:2), error = function(e) conditionMessage(e))
#---
tryCatch(rowsum("a", "b"), error = function(e) conditionMessage(e))
#---
rowsum(c(TRUE, FALSE, TRUE), c(1, 1, 2))
#---
rowsum(1:4, c(1.5, 2.5, 1.5, 2.5))
#---
rowsum.default(1:4, c(1, 1, 2, 2))
#---
rowsum.data.frame(data.frame(x = 1:4), c(1, 1, 2, 2))
#---
f <- function(x, ...) chkDots(...); tryCatch(f(1, a = 2), warning = function(w) conditionMessage(w))
#---
w <- structure(list("first" = quote(f(x)), "second" = NULL), class = "warnings"); w
#---
w <- structure(list("only" = quote(g(1, 2))), class = "warnings"); w
#---
w <- structure(list("only" = NULL), class = "warnings"); w
#---
w <- structure(list(), class = "warnings"); w
#---
print(structure(list(a = quote(f())), class = "warnings"))
#---
w <- structure(list("a very long warning message that goes on and on and on and on and on and on" = quote(some_function_name(argument_one, argument_two, argument_three))), class = "warnings"); w
#---
w <- structure(list("m1" = quote(f(x)), "m2" = quote(g(y))), class = "warnings"); print(w, tags = c("A", "B"))
#---
w <- structure(list("m1" = quote(f(x))), class = "warnings"); c(w, w)
#---
rev(structure(list("m1" = quote(f(x)), "m2" = NULL), class = "warnings"))
#---
unique(structure(list("m1" = quote(f(x)), "m1" = quote(f(x))), class = "warnings"))
#---
as.character(structure(list("m1" = quote(f(x))), class = "warnings"))
#---
array2DF(array(1:4, c(2, 2), dimnames = list(a = c("x", "y"), b = c("p", "q"))))
#---
array2DF(array(1:4, c(2, 2)))
#---
array2DF(array(1:4, c(2, 2), dimnames = list(c("x", "y"), c("p", "q"))), responseName = "val")
#---
array2DF(array(1:3, 3, dimnames = list(k = c("a", "b", "c"))))
#---
array2DF(table(c("a", "b", "a")))
#---
array2DF(array(1:8, c(2, 2, 2), dimnames = list(a = c("x", "y"), b = c("p", "q"), c = c("u", "v"))))
#---
array2DF(list(a = 1, b = 2))
#---
array2DF(matrix(1:4, 2), simplify = FALSE)
#---
array2DF(array(list(1, "a", TRUE, NULL), c(2, 2)))
#---
array2DF(array(1:4, c(2, 2), dimnames = list(a = c("x", "y"), b = c("p", "q"))), allowLong = FALSE)
#---
as.data.frame(1:3)
#---
as.data.frame(c(a = 1, b = 2))
#---
as.data.frame(letters[1:2])
#---
as.data.frame(list(a = 1:2, b = 3:4))
#---
as.data.frame(matrix(1:4, 2))
#---
as.data.frame(matrix(1:4, 2, dimnames = list(NULL, c("u", "v"))))
#---
as.data.frame(factor(c("a", "b")))
#---
as.data.frame(as.Date("2024-01-01") + 0:1)
#---
as.data.frame(table(c("a", "b", "a")))
#---
as.data.frame(array(1:8, c(2, 2, 2)))
#---
as.data.frame(NULL)
#---
as.data.frame(TRUE)
#---
as.data.frame(1:2, nm = "z")
#---
as.data.frame(1:2, optional = TRUE)
#---
as.data.frame(as.POSIXct("2024-01-01 10:00:00", tz = "UTC"))
#---
as.data.frame(difftime(as.Date("2024-01-02"), as.Date("2024-01-01")))
#---
as.data.frame(list(1, 2))
#---
as.data.frame(list(a = 1, b = "x"), stringsAsFactors = TRUE)
#---
as.data.frame(ordered(c("a", "b")))
#---
as.data.frame(as.numeric_version("1.2.3"))
#---
as.data.frame(1i)
#---
as.data.frame(raw(2))
#---
as.data.frame(quote(x))
#---
as.data.frame(structure(list(a = 1:2), class = "foo"))
#---
as.data.frame(I(c("a", "b")))
#---
as.data.frame.vector(1:2)
#---
as.data.frame.integer(1:2)
#---
as.data.frame.numeric(c(1.5, 2))
#---
as.data.frame.character(c("a", "b"))
#---
as.data.frame.logical(c(TRUE, FALSE))
#---
as.data.frame.factor(factor("a"))
#---
as.data.frame.Date(as.Date("2024-01-01"))
#---
as.data.frame.list(list(a = 1))
#---
as.data.frame.matrix(matrix(1:2, 1))
#---
as.data.frame.table(table(c(1, 1)))
#---
as.data.frame.AsIs(I(1:2))
#---
as.data.frame.complex(1i)
#---
as.data.frame.raw(raw(1))
#---
as.data.frame.ordered(ordered("a"))
#---
as.data.frame.POSIXct(as.POSIXct("2024-01-01", tz = "UTC"))
#---
as.data.frame.difftime(as.difftime(1, units = "mins"))
#---
as.data.frame.array(array(1:4, c(2, 2)))
#---
as.data.frame.noquote(noquote("a"))
#---
as.data.frame.numeric_version(as.numeric_version("1.2"))
#---
as.data.frame.POSIXlt(as.POSIXlt("2024-01-01", tz = "UTC"))
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
x[2]
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
x[2:3]
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
x[-1]
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
x[c(TRUE, FALSE, TRUE)]
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
x[5]
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
x[0]
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
class(x[2])
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
length(x[2:3])
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
x[[2]]
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
x[["year"]]
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
x[[1]] <- as.POSIXlt("2000-01-01 00:00:00", tz = "UTC"); x
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
y <- x; y$year <- y$year + 1L; y
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
y <- x; y$hour <- 5L; format(y)
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
y <- x; y$zone <- NULL; names(unclass(y))
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
x$year
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
x$mon
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
names(unclass(x))
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
unlist(unclass(x[1]))
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
attributes(x[2])$tzone
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
attr(x[1], "balanced")
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
y <- x; y[2] <- as.POSIXlt("1999-12-31 12:00:00", tz = "UTC"); y
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
y <- x; y[[2]] <- as.POSIXlt("1999-12-31 12:00:00", tz = "UTC"); y
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
rev(x)
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
x[order(c(3, 1, 2))]
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
as.character(x[1])
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
y <- x; y[4] <- as.POSIXlt("2001-01-01", tz = "UTC"); length(y)
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
format(x[2], "%Y")
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
is.na(x[1])
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
l <- as.POSIXlt("2024-06-01 12:00:00", tz = "Europe/Paris"); l[1]; attributes(l[1])$tzone
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
names(as.POSIXlt(c(a = "2024-01-01"), tz = "UTC"))
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30", "2024-03-25 23:59:59"), tz = "UTC")
x2 <- as.POSIXlt(c(a = "2024-01-01", b = "2024-01-02"), tz = "UTC"); x2["b"]
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$year <- y$year + 1L; attributes(y)
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$hour <- 5L; attr(y, "balanced")
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$zone <- NULL; attr(y, "balanced")
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$foo <- 1; names(unclass(y))
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$sec <- c(1, 2, 3); length(y)
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$year <- 1L; y
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
attributes(x[1:2])
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
attributes(x[2:1])
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
x[NA]
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
x[c(1, NA)]
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
x[-5]
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
tryCatch(x["a"], error = function(e) conditionMessage(e))
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
tryCatch(x[1, 2], error = function(e) conditionMessage(e))
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
x[]
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
x[TRUE]
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
x[integer(0)]
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
length(x[3])
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
x[[1]]
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y[[3]] <- x[[1]]; y
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y[3] <- x[1]; length(y)
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y[1] <- NA; y
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y[1] <- as.POSIXct("2001-01-01 00:00:00", tz = "UTC"); y
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y[1] <- "2001-05-05"; y
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y[c(1, 2)] <- x[2]; y
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y[0] <- x[1]; length(y)
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
tryCatch({y <- x; y[1] <- 5}, error = function(e) conditionMessage(e))
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- as.POSIXlt("2024-06-01 12:00:00", tz = "Europe/Paris"); y[1]
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
attributes(y[1])
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$year <- 100L; attr(y, "balanced")
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$year <- c(100L, 101L); attr(y, "balanced")
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$year <- x$year + 1L; attr(y, "balanced")
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$hour <- c(5L, 6L); attr(y, "balanced")
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$hour <- 5L; attr(y, "balanced")
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$min <- x$min + 1L; attr(y, "balanced")
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$mon <- 3L; attr(y, "balanced")
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$sec <- 1.5; attr(y, "balanced")
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$mday <- x$mday - 1L; attr(y, "balanced")
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$foo <- 1:2; attr(y, "balanced")
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$zone <- c("A", "B"); attr(y, "balanced")
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y[1] <- x[2]; attr(y, "balanced")
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y[[1]] <- x[2]; attr(y, "balanced")
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y[3] <- x[2]; attr(y, "balanced")
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$year <- 100L; y$hour <- 3L; attr(y, "balanced")
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
attr(as.POSIXlt("2024-01-01", tz = "UTC"), "balanced")
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$year <- 100L; c(y, y)
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$hour <- 25L; format(y)
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$hour <- 25L; attributes(y)$balanced
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$hour <- 25L; y[1]
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC")
y <- x; y$hour <- 25L; attr(y[1], "balanced")
#---
ngettext(1, "one", "many")
#---
ngettext(2, "one", "many")
#---
tryCatch(ngettext(-1, "a", "b"), error = function(e) conditionMessage(e))
#---
f <- function(x, ...) chkDots(...); tryCatch(f(1, a = 2), warning = function(w) conditionMessage(w))
#---
f <- function(x, ...) chkDots(...); f(1)
#---
f <- function(x, ...) chkDots(...); tryCatch(f(1, a = 2, b = 3), warning = function(w) conditionMessage(w))
#---
rowsum(c(1, 2, 3), c("a", "b", "a"), reorder = TRUE)
#---
warnings()
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC"); `[.POSIXlt`(x, 2)
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC"); `[[.POSIXlt`(x, 2)
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC"); `[<-.POSIXlt`(x, 1, value = x[2])
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC"); `[[<-.POSIXlt`(x, 1, value = x[2])
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC"); `$<-.POSIXlt`(x, "hour", 3L)
#---
x <- as.POSIXlt(c("2024-01-15 10:30:00", "2024-02-20 11:45:30"), tz = "UTC"); c.POSIXlt(x, x)
#---
x <- as.POSIXlt(c(a = "2024-01-15 10:30:00", b = "2024-02-20 11:45:30"), tz = "UTC"); names(x) <- c("p", "q"); names(x)
#---
w <- structure(list(a = quote(f()), b = quote(g())), class = "warnings"); `[.warnings`(w, 2)
#---
w <- structure(list(a = quote(f()), b = quote(g())), class = "warnings"); c.warnings(w, w)
#---
w <- structure(list(a = quote(f()), a = quote(f())), class = "warnings"); unique.warnings(w)
#---
w <- structure(list(a = quote(f()), a = quote(f())), class = "warnings"); duplicated.warnings(w)
#---
w <- structure(list(a = quote(f())), class = "warnings"); print.warnings(w)
#---
x <- as.POSIXlt(c(a = "2024-01-15 10:30:00", b = "2024-02-20 11:45:30"), tz = "UTC"); names.POSIXlt(x)
#---
x <- as.POSIXlt(c(a = "2024-01-15 10:30:00", b = "2024-02-20 11:45:30"), tz = "UTC"); names(`names<-.POSIXlt`(x, value = c("p", "q")))
#---
tryCatch(Vectorize(function(x) x, "zz"), error = function(e) conditionMessage(e))
#---
class(Vectorize(function() 1))
