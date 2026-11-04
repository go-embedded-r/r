lapply(1:3, function(i) i^2)
#---
lapply(c(a = 1, b = 2), function(x) x + 1)
#---
lapply(list(1, "a", TRUE), class)
#---
lapply(1:2, function(x, y) x + y, y = 10)
#---
lapply(NULL, identity)
#---
sapply(1:5, function(i) i * 2)
#---
sapply(c("a", "bb", "ccc"), nchar)
#---
sapply(c(a = 1, b = 2), function(x) x * 10)
#---
sapply(1:3, function(i) c(i, i^2))
#---
sapply(1:3, function(i) c(x = i, y = i^2))
#---
sapply(1:3, function(i) if (i > 1) 1:i else 1)
#---
sapply(list(), length)
#---
sapply(1:3, function(i) letters[i])
#---
sapply(1:3, function(i) i, simplify = FALSE)
#---
sapply(c("x", "y"), toupper, USE.NAMES = FALSE)
#---
sapply(1:3, function(i) list(i))
#---
vapply(1:3, function(i) i * 2, numeric(1))
#---
vapply(c("a", "bb"), nchar, integer(1))
#---
vapply(1:2, function(i) c(i, i), numeric(2))
#---
vapply(1:3, function(i) letters[i], character(1))
#---
vapply(1:3, function(i) "a", numeric(1))
#---
vapply(1:3, function(i) 1:2, numeric(1))
#---
vapply(integer(0), function(i) i, numeric(1))
#---
Map(function(x, y) x + y, 1:3, 4:6)
#---
Map(function(n, v) paste(n, v), c("a", "b"), c("x", "y"))
#---
Map(`+`, c(a = 1, b = 2), 10)
#---
mapply(function(x, y) x * y, 1:3, 4:6)
#---
mapply(function(x, y) x + y, 1:2, c(a = 3, b = 4))
#---
mapply(rep, 1:3, 3:1)
#---
mapply(function(x, p) x^p, 2:3, MoreArgs = list(p = 2))
#---
mapply(function(x, y) x + y, 1:3, 1:3, SIMPLIFY = FALSE)
#---
Reduce(`+`, 1:5)
#---
Reduce(`+`, 1:4, accumulate = TRUE)
#---
Reduce(function(a, b) paste0(a, b), c("a", "b", "c"), accumulate = TRUE, right = TRUE)
#---
Reduce(`+`, list(), 0)
#---
Reduce(`+`, list())
#---
Reduce(function(a, b) a * b, 1:5, 100)
#---
Filter(function(x) x > 2, 1:5)
#---
Filter(function(x) x > 1, c(a = 1, b = 2, c = 3))
#---
Filter(Negate_not_defined <- function(x) TRUE, list(1, 2))
#---
Find(function(x) x > 2, 1:5); Find(function(x) x > 10, 1:5)
#---
Position(function(x) x > 2, c(1, 3, 5)); Position(function(x) x > 2, c(1, 3, 5), right = TRUE); Position(function(x) x > 9, 1:3)
#---
do.call(sum, list(1, 2, 3))
#---
do.call("paste", list("a", "b", sep = "-"))
#---
f <- function(x, y) x - y; do.call(f, list(y = 1, x = 10))
#---
do.call(function(...) length(list(...)), list(1, 2, 3))
#---
do.call("nonexistent_fn", list())
#---
do.call(1, list())
#---
match.fun("sum")(1, 2)
#---
match.fun(5)
#---
unlist(list(1, 2, 3))
#---
unlist(list(a = 1, b = 2))
#---
unlist(list(a = 1:2, b = 3))
#---
unlist(list(a = list(b = 1, c = 2), d = 3))
#---
unlist(list(a = c(x = 1, y = 2), b = 3))
#---
unlist(list(1, "a", TRUE))
#---
unlist(list(1:2, list(3, 4)), recursive = FALSE)
#---
unlist(list(a = 1, b = 2), use.names = FALSE)
#---
unlist(list())
#---
unlist(5)
#---
rev(1:5); rev(c(a = 1, b = 2)); rev(list(1, "a")); rev(NULL)
#---
which(c(TRUE, FALSE, TRUE, NA)); which(c(a = TRUE, b = FALSE, c = TRUE))
#---
which(1:3)
#---
rep(1:2, 3); rep(1:2, each = 2); rep("a", 0); rep(c(1, 2), times = c(2, 3)); rep(1:3, length.out = 7)
#---
rep(c(a = 1, b = 2), 2)
#---
rep(NULL, 3)
#---
rep(1, -1)
#---
rep(list(1, "a"), 2)
#---
seq(5); seq(2, 8); seq(10, 1); seq(0, 1, by = 0.25); seq(1, 10, length.out = 4); seq(5, 1, by = -2)
#---
seq_len(3); seq(2, 11, 3)
#---
seq(1, 10, by = -1)
#---
seq(c(5, 6, 7))
