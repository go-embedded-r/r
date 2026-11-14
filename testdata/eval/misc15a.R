x <- as.hexmode(c(10, 255, 16)); unclass(x)
#---
x <- as.hexmode(c(10, 255, 16)); class(x)
#---
x <- as.hexmode(c(10, 255, 16)); x
#---
x <- as.hexmode(c(10, 255, 16)); format(x)
#---
x <- as.hexmode(c(10, 255, 16)); format(x, width = 4, upper.case = TRUE)
#---
x <- as.hexmode(c(10, 255, 16)); as.character(x)
#---
x <- as.hexmode(c(10, 255, 16)); x[2]
#---
x <- as.hexmode(c(10, 255, 16)); x & as.hexmode(15)
#---
x <- as.hexmode(c(10, 255, 16)); x | 1L
#---
x <- as.hexmode(c(10, 255, 16)); !x
#---
as.octmode(c(8, 64))
#---
format(as.octmode(c(8, 64)), width = 5)
#---
as.character(as.octmode(c(8, 64, NA)))
#---
format(as.hexmode(c(1, NA, 300)))
#---
as.hexmode("ff")
#---
as.hexmode(c("a", "zz"))
#---
as.hexmode(1.5)
#---
as.octmode("17")
#---
as.octmode(as.octmode(9L))
#---
as.octmode(7L) | as.octmode(8L)
#---
as.octmode(15L) & 9L
#---
!as.octmode(5L)
#---
as.octmode(1:10)[3:4]
#---
isa(1L, "integer")
#---
isa(structure(1, class = c("a", "b")), c("a", "b"))
#---
isa(structure(1, class = c("a", "b")), c("b", "a"))
#---
isa(1, "integer")
#---
isa(matrix(1), "matrix")
#---
gl(2, 3)
#---
gl(2, 1, 5, labels = c("x", "y"), ordered = TRUE)
#---
gl(3, 2, 4)
#---
gl(2, 2, 3, labels = c("a", "b"))
#---
y <- I(1:3); class(y)
#---
I(1:3)
#---
I(c(a = 1, b = 2))
#---
format(I(c(1.5, 2)))
#---
y <- I(1:3); y[2:3]
#---
class(I(factor("a")))
#---
format(I(list(1, "a", 1:3)))
#---
is.table(table(1))
#---
is.table(1)
#---
data.class(1:3)
#---
data.class(matrix(1:4, 2))
#---
data.class(array(1:8, c(2, 2, 2)))
#---
data.class(factor("a"))
#---
data.class(list(1))
#---
data.class("a")
#---
labels(1:3)
#---
labels(c(a = 1, b = 2))
#---
labels(matrix(1:4, 2))
#---
labels(matrix(1:4, 2, dimnames = list(c("a", "b"), NULL)))
#---
ordered(c("b", "a"))
#---
ordered(c("b", "a"), levels = c("b", "a"))
#---
is.ordered(ordered(1))
#---
is.ordered(factor(1))
#---
as.ordered(c(3, 1, 2))
#---
as.ordered(ordered(c("x", "y")))
#---
x <- 1:5; is.na(x) <- 2; x
#---
x <- c(a = 1, b = 2, c = 3); is.na(x) <- "b"; x
#---
f <- factor(c("a", "b")); is.na(f) <- 1; f
#---
x <- list(1, 2); is.na(x) <- 1; x
#---
s <- "abcdef"; substring(s, 2, 3) <- "ZZZ"; s
#---
s <- c("abcdef", "ghijkl"); substring(s, 3) <- c("XY", "Z"); s
#---
x <- 1:6; mostattributes(x) <- list(dim = c(2, 3), foo = "bar", names = letters[1:6]); attributes(x)
#---
x <- 1:3; mostattributes(x) <- list(names = c("a", "b", "c"), dim = 4, class = "zz"); attributes(x)
#---
x <- 1:4; mostattributes(x) <- list(dim = c(2, 2), dimnames = list(c("a", "b"), c("c", "d"))); x
#---
x <- 1:4; mostattributes(x) <- 3
#---
arrayInd(5, c(2, 3))
#---
arrayInd(c(2, 5), c(2, 3), .dimnames = list(letters[1:2], NULL), useNames = TRUE)
#---
arrayInd(c(1, 7, 8), c(2, 2, 2))
#---
arrayInd(3, c(2, 2), useNames = TRUE)
#---
slice.index(array(1:8, c(2, 2, 2)), 3)
#---
slice.index(matrix(0, 2, 3), 1)
#---
slice.index(matrix(0, 2, 3), 2)
#---
slice.index(matrix(0, 2, 3), 1:2)
#---
slice.index(1:3, 1)
#---
as.array(1:3)
#---
as.array(c(a = 1, b = 2))
#---
as.array(matrix(1:4, 2))
#---
as.array(list(1, 2))
#---
as.null(1)
#---
as.null.default(1:3)
#---
toString(1:5, collapse = "; ")
#---
toString(c("a", "b"), width = 3)
#---
toString(1:20, width = 12)
#---
toString(c("x", "y"), sep = "-")
#---
toString(1:3, width = -1)
#---
toString(NULL)
#---
toString.default(letters[1:4], width = 0)
#---
alist(a = , b = 1 + 2, 3)
#---
length(alist())
#---
enquote(quote(a + b))
#---
enquote(c(1, 2))
#---
f <- function() .Deprecated("g"); f()
#---
f <- function() .Deprecated(); tryCatch(f(), warning = function(w) class(w))
#---
f <- function() .Deprecated("g", package = "pk"); f()
#---
f <- function() .Deprecated(msg = "custom"); f()
#---
f <- function() .Defunct("g"); f()
#---
f <- function() .Defunct(); r <- tryCatch(f(), error = function(e) e); list(class(r), conditionCall(r))
#---
f <- function() .Defunct(msg = "custom"); f()
#---
f <- function() .NotYetImplemented(); f()
#---
f <- function(a) .NotYetUsed("a"); f(1)
#---
f <- function(a) .NotYetUsed("a", error = FALSE); f(1)
#---
x <- 1.5; class(x) <- "factor"
#---
class(structure(1:3, class = "factor"))
#---
x <- 1:3; class(x) <- "zz"; names(x) <- letters[1:3]; names(attributes(x))
#---
x <- 1:3; names(x) <- letters[1:3]; class(x) <- "zz"; names(attributes(x))
#---
strtrim("abcdef", 3)
#---
strtrim(c(a = "abcdef", b = "xy", c = NA), c(1, 5, 2))
#---
strtrim("abc", -1)
#---
strtrim(1, 1)
#---
strtrim(c("abc", "def"), 1:3)
#---
strtrim("日本語です", 5)
#---
strtrim(character(0), 1)
#---
strtrim(matrix(c("abc", "defg"), 1), 2)
#---
mean(c(1, 2, 60, 100, 3), trim = 0.2)
#---
mean(c(1, 2, 60, 100, 3), trim = 0.5)
#---
mean(c(1, NA, 3), trim = 0.1)
#---
mean(c(1, NA, 3, 4, 5, 9), trim = 0.1, na.rm = TRUE)
#---
mean(1:10, trim = 0.6)
#---
mean(c(1, 5, 2, 3, 10, 40, 41), trim = 0.3)
#---
is.primitive(sum)
#---
is.primitive(mean)
#---
is.primitive(`if`)
#---
is.primitive(1)
#---
is.R()
#---
f <- function(a, b = 2, ...) NULL; formals(f)$b
#---
f <- function(a, b = 2, ...) NULL; length(formals(f))
#---
f <- function(a, b = 2, ...) NULL; formals("f")$b
#---
f <- function(a, b = 2, ...) NULL; names(formals(f))
#---
formals(function() 1)
#---
formals(sum)
#---
f <- function(a) a * 2; body(f) <- quote(a + 100); f(1)
#---
f <- function(a) a * 2; body(f) <- quote(a + 100); body(f)
#---
f <- function(a, b) a * 2; formals(f) <- alist(x = , y = 7); f(3)
#---
f <- function(a, b) a * 2; formals(f) <- alist(x = , y = 7); names(formals(f))
#---
f <- function(a, b) a * 2; formals(f)$b <- 10; formals(f)$b
#---
g <- as.function(alist(a = , b = 2, a * b)); g(3)
#---
g <- as.function(alist(a = , b = 2, a * b)); names(formals(g))
#---
as.function(alist(1))()
#---
as.function(list())
#---
as.function(1)
#---
diff(c(1, 4, 9, 16))
#---
diff(c(1, 4, 9, 16), lag = 2)
#---
diff(c(1, 4, 9, 16), differences = 2)
#---
diff(1:3, differences = 3)
#---
diff(matrix(c(1, 4, 9, 16, 25, 36), 3))
#---
diff(1:5, lag = 0)
#---
diff(c(a = 1, b = 3, c = 8))
#---
unique.default(c(1, 2, 1))
#---
unique.default(factor(c("a", "b", "a")))
#---
duplicated.default(c(1, 2, 1))
#---
anyDuplicated.default(c(1, 2, 1))
#---
rev.default(1:3)
#---
rev.default(c(a = 1, b = 2))
#---
sort.default(c(3, 1, 2), decreasing = TRUE)
#---
xtfrm.default(c(3, 1, 2))
#---
range.default(c(3, 1, 2))
#---
range.default(c(3, NA, 2), na.rm = TRUE)
#---
mean.default(c(1, 2, 6))
#---
mean.default(c(1, 2, 60), trim = 0.4)
#---
seq.default(1, 10, length.out = 4)
#---
seq.default(5)
#---
format.default(c(1.5, 10))
#---
f <- tempfile(); write(1:7, f); readLines(f)
#---
f <- tempfile(); write(c("a", "b"), f); readLines(f)
#---
f <- tempfile(); write(1:4, f, ncolumns = 2, sep = ","); readLines(f)
#---
f <- tempfile(); write(1:2, f); write(3:4, f, append = TRUE); readLines(f)
#---
x <- as.hexmode(c(10, 255)); `!.hexmode`(x)
#---
`!.octmode`(as.octmode(8L))
#---
`&.hexmode`(as.hexmode(12L), 10L)
#---
`&.octmode`(as.octmode(12L), 10L)
#---
`|.hexmode`(as.hexmode(12L), 10L)
#---
`|.octmode`(as.octmode(12L), 10L)
#---
`[.hexmode`(as.hexmode(1:5 * 40L), 2:3)
#---
`[.octmode`(as.octmode(1:5 * 40L), 2:3)
#---
`[.AsIs`(I(c(a = 1, b = 2)), 2)
#---
as.character.hexmode(as.hexmode(c(255L, 16L)))
#---
as.character.octmode(as.octmode(c(255L, 16L)))
#---
format.hexmode(as.hexmode(255L), upper.case = TRUE)
#---
format.octmode(as.octmode(c(8L, 100L)))
#---
format.AsIs(I(c("a", "bb")))
#---
format.AsIs(I(list(1:3, "x")), width = 4)
#---
print.hexmode(as.hexmode(255L))
#---
print.octmode(as.octmode(255L))
#---
x <- I(1:2); tryCatch(print.AsIs(x), error = function(e) "err")
#---
`is.na<-.default`(1:3, 2)
#---
`is.na<-.factor`(factor(c("a", "b")), 2)
#---
labels.default(c(x = 1))
#---
levels.default(factor(c("b", "a")))
#---
levels.default(1:3)
#---
as.array.default(1:3)
#---
as.function.default(alist(x = , x + 1))(2)
#---
f <- function() sys.nframe() - sys.parent(); f()
#---
g <- function() f(); f <- function() sys.nframe() - sys.parent(); g()
#---
f <- function() sys.nframe() - sys.parent(2); g <- function() f(); h <- function() g(); h()
#---
f <- function() sys.nframe() - sys.parent(1); g <- function() f(); h <- function() g(); h()
#---
tryCatch(print.AsIs(I("a")), error = function(e) conditionMessage(e))
#---
diff.default(c(1, 4, 9), 2)
#---
diff.Date(as.Date(c("2024-01-05", "2024-01-08", "2024-02-01")))
#---
diff.Date(as.Date(c("2024-01-05", "2024-01-08", "2024-02-01")), differences = 2)
#---
diff.Date(as.Date("2024-01-05"))
#---
diff.POSIXt(as.POSIXct(c("2024-03-10 00:00:00", "2024-03-10 00:00:30"), tz = "UTC"))
#---
diff.POSIXt(as.POSIXct(c("2024-03-10 00:00:00", "2024-03-10 00:10:30", "2024-03-10 01:00:00"), tz = "UTC"), differences = 2)
#---
diff.POSIXt(as.POSIXct(c("2024-03-10", "2024-03-14", "2024-03-20"), tz = "UTC"), lag = 2)
#---
diff.POSIXt(as.POSIXlt(c("2024-03-10", "2024-03-14"), tz = "UTC"))
#---
diff.POSIXt(as.POSIXct("2024-03-10", tz = "UTC"))
#---
diff.difftime(as.difftime(c(1, 4, 9), units = "mins"))
#---
diff.difftime(as.difftime(c(1, 4, 9), units = "hours"), differences = 2)
#---
diff(as.Date(c("2024-01-05", "2024-01-08")), lag = 0)
#---
diff(as.POSIXct(c("2024-03-10 00:00", "2024-03-10 03:00"), tz = "UTC"), differences = 0)
