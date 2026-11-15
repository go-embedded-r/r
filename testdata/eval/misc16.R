as.table(matrix(1:4, 2))
#---
as.table(c(a = 1, b = 2))
#---
as.table(1:3)
#---
class(as.table(matrix(1:4, 2)))
#---
dimnames(as.table(matrix(1:4, 2)))
#---
as.table.default(matrix(1:4, 2, dimnames = list(x = c("p", "q"), y = c("r", "s"))))
#---
as.table(array(1:8, c(2, 2, 2)))
#---
tryCatch(as.table(list(1)), error = function(e) conditionMessage(e))
#---
simplify2array(list(1:2, 3:4))
#---
simplify2array(list(a = 1, b = 2))
#---
simplify2array(list(1:2, 3:4), higher = FALSE)
#---
simplify2array(list(1:2, 3:5))
#---
simplify2array(list(matrix(1:4, 2), matrix(5:8, 2)))
#---
simplify2array(list(matrix(1:4, 2), matrix(5:8, 2)), higher = FALSE)
#---
data.matrix(data.frame(a = 1:2, b = c("x", "y"), c = c(1.5, 2)))
#---
data.matrix(data.frame(a = 1:2, b = factor(c("x", "y"))))
#---
data.matrix(data.frame(a = c(TRUE, FALSE), b = c(1.5, NA)))
#---
data.matrix(data.frame(a = 1:2, row.names = c("r1", "r2")))
#---
list2DF(list(a = 1:2, b = c("x", "y")))
#---
list2DF(list(1:2, 3:4))
#---
list2DF()
#---
list2DF(list(a = 1:2), nrow = 2)
#---
tryCatch(list2DF(list(a = 1:3, b = 1), nrow = 3), error = function(e) conditionMessage(e))
#---
sort_by(data.frame(a = c(3, 1, 2), b = c("x", "y", "z")), data.frame(k = c(3, 1, 2)))
#---
sort_by(c(3, 1, 2), c(3, 1, 2))
#---
sort_by(c(b = 3, a = 1), c(2, 1))
#---
sort_by(data.frame(a = c(3, 1, 2), b = c("x", "y", "z")), list(c(3, 1, 2)), decreasing = TRUE)
#---
sort_by.default(c(3, 1, 2), c(3, 1, 2))
#---
sort_by.data.frame(data.frame(a = c(3, 1, 2)), list(c(3, 1, 2)))
#---
provideDimnames(matrix(1:4, 2))
#---
provideDimnames(matrix(1:4, 2, dimnames = list(c("a", "b"), NULL)))
#---
provideDimnames(matrix(1:4, 2), base = list(letters, LETTERS))
#---
provideDimnames(matrix(1:4, 2), unique = FALSE)
#---
logb(8, 2)
#---
logb(100, 10)
#---
logb(exp(2))
#---
logb(c(1, 10, 100), base = 10)
#---
x <- 1:3; mode(x) <- "character"; x
#---
x <- c("1", "2"); mode(x) <- "numeric"; x
#---
x <- 1:3; mode(x) <- "list"; x
#---
x <- c(1.5, 2); mode(x) <- "integer"; x
#---
x <- c(a = 1L, b = 2L); mode(x) <- "character"; x
#---
a1 <- 1; "a1" %in% objects()
#---
a1 <- 1; objects(pattern = "^a1$")
#---
g <- function() { zz <- 5; h() }; h <- function() dynGet("zz", ifnotfound = "none"); g()
#---
h <- function() dynGet("zz", ifnotfound = "none"); h()
#---
dynGet("nosuch", 3)
#---
cat(tryCatch(dynGet("nosuch"), error = function(e) conditionMessage(e)), "\n")
#---
unique(matrix(c(1, 1, 2, 3, 3, 4), 3, byrow = TRUE))
#---
unique(matrix(c(1, 2, 1, 2, 3, 4), 3, byrow = TRUE))
#---
duplicated(matrix(c(1, 2, 1, 2, 3, 4), 3, byrow = TRUE))
#---
duplicated(matrix(c(1, 1, 2, 2, 1, 1), 2, byrow = TRUE), MARGIN = 2)
#---
anyDuplicated(matrix(c(1, 2, 1, 2, 3, 4), 3, byrow = TRUE))
#---
unique(matrix(c(1, 2, 1, 2, 3, 4), 3, byrow = TRUE), MARGIN = 2)
#---
unique(array(c(1, 1, 2, 2), c(2, 2)), fromLast = TRUE)
#---
duplicated(array(1:8, c(2, 2, 2)), MARGIN = 3)
#---
unique(matrix(c("a", "a", "b", "b"), 2))
#---
anyDuplicated(matrix(1:4, 2))
#---
unique.matrix(matrix(c(1, 1, 2, 2), 2, byrow = TRUE))
#---
unique.array(array(c(1, 1, 2, 2), c(2, 2)))
#---
duplicated.matrix(matrix(c(1, 1, 2, 2), 2, byrow = TRUE))
#---
duplicated.array(array(c(1, 1, 2, 2), c(2, 2)))
#---
anyDuplicated.matrix(matrix(c(1, 1, 2, 2), 2, byrow = TRUE))
#---
anyDuplicated.array(array(c(1, 1, 2, 2), c(2, 2)))
#---
subset(matrix(1:6, 3, dimnames = list(NULL, c("a", "b"))), c(TRUE, FALSE, TRUE))
#---
subset(matrix(1:6, 3, dimnames = list(NULL, c("a", "b"))), c(TRUE, FALSE, TRUE), select = b)
#---
subset(matrix(1:6, 3, dimnames = list(NULL, c("a", "b"))), select = -a)
#---
subset.matrix(matrix(1:4, 2, dimnames = list(NULL, c("a", "b"))), c(TRUE, FALSE))
#---
transform(list(a = 1:2), b = a * 2)
#---
transform(c(a = 1), b = 2)
#---
transform.default(list(a = 1:2), b = a * 2)
#---
class(system.time(NULL))
#---
names(system.time(NULL))
#---
length(system.time(NULL))
#---
st <- system.time(for (i in 1:10) i); is.numeric(st)
#---
names(summary(system.time(NULL)))
#---
cat(tryCatch(unix.time(NULL), error = function(e) conditionMessage(e)), "\n")
#---
summary.proc_time(structure(c(1, 2, 3, 4, 5), class = "proc_time"))
#---
print.proc_time(structure(c(0, 0, 0, 0, 0), class = "proc_time"))
#---
structure(c(0, 0, 0, 0, 0), class = "proc_time")
#---
as.single(1.5)
#---
single(2)
#---
class(as.single(1))
#---
attributes(as.single(1.5))
#---
as.single.default(2)
#---
pretty(numeric(0))
#---
pretty(c(NA_real_, NA))
#---
pretty(5)
#---
pretty(0)
#---
pretty(c(0, 0))
#---
pretty(c(1, 10))
#---
pretty(c(0, 97))
#---
pretty(1:10, n = 3)
#---
pretty(c(-1, 1))
#---
pretty(c(0.001, 0.0099))
#---
pretty(1e6 * (1:5))
#---
pretty(c(123456, 654321), n = 10)
#---
pretty(c(1, 1.0000001))
#---
pretty(100)
#---
pretty(-3)
#---
pretty(c(3, 3), min.n = 0)
#---
pretty(c(2, 2), n = 1)
#---
pretty(1:3, n = 0)
#---
pretty(c(0, 1), n = 10, high.u.bias = 5)
#---
pretty(c(0, 1), n = 7, u5.bias = 0)
#---
pretty(c(-5.5, 7.2))
#---
pretty(c(1e10, 1e10 + 1))
#---
pretty(c(-1e-10, 1e-10))
#---
.pretty(c(0, 97))
#---
.pretty(c(0, 97), bounds = FALSE)
#---
names(.pretty(c(1, 5)))
#---
pretty(c(TRUE, FALSE))
#---
pretty(c(1, Inf))
#---
pretty(c(1, NA, 3))
#---
pretty(c(0.1, 0.35), n = 4)
#---
pretty(c(0, 1), eps.correct = 1)
#---
pretty(c(10, 11), n = 20)
#---
pretty(c(0, 1e300))
#---
pretty(12345.678)
#---
pretty(c(0, 12345.678))
#---
pretty(-1:1 / 3)
#---
pretty.default(c(0, 10))
#---
f <- factor(c("a", "b", "a")); levels(f) <- c("x", "y"); f
#---
f <- factor(c("a", "b", "c", "a")); levels(f) <- c("x", "y", "x"); f
#---
f <- factor(c("a", "b")); levels(f) <- c("x", "y", "z"); f
#---
f <- factor(c("a", "b", "c")); tryCatch(levels(f) <- c("x", "y"), error = function(e) conditionMessage(e))
#---
f <- factor(c("a", "b"), levels = c("b", "a")); levels(f) <- c("B", "A"); f
#---
f <- ordered(c("lo", "hi", "lo"), levels = c("lo", "hi")); levels(f) <- c("L", "H"); f
#---
f <- factor(c(u = "a", v = "b")); levels(f) <- c("x", "y"); f
#---
f <- factor(c("a", NA, "b")); levels(f) <- c("x", "y"); f
#---
f <- factor(c("a", "b")); `levels<-.factor`(f, c("p", "q"))
#---
tryCatch(pretty(c(-193.753, 0.0003), n = 1, min.n = 2), error = function(e) conditionMessage(e))
#---
tryCatch(pretty(1:3, n = 2, min.n = -1), error = function(e) conditionMessage(e))
#---
d <- as.Date(c("2024-03-12", "2024-01-05")); tryCatch(cut(d, 3), error = function(e) conditionMessage(e))
#---
formatC(c(42, 5338551, 0), format = "g", digits = 4, zero.print = "-")
#---
formatC(c(-4.39e+03, 0.00129147, 0), format = "g", digits = 5, big.mark = " ", big.interval = 2)
#---
formatC(c(76, 0), format = "g", digits = 5)
#---
formatC(c(7, 12345678), format = "fg", digits = 3, big.mark = ",")
#---
seq(0, 1, length.out = 11)[11] == 1
#---
seq(-8.4e-05, 5.5e7, length.out = 12)[12]
