.rowSums(1:6, 2, 3)
#---
.colSums(c(1, 2, NA, 4, 5, 6), 2, 3, na.rm = TRUE)
#---
.rowMeans(c(1, 2, NA, 4, 5, 6), 2, 3, na.rm = TRUE)
#---
.colMeans(1:6, 3, 2)
#---
.rowSums(letters, 2, 13)
#---
.colSums(1:6, 4, 4)
#---
.row(c(2, 3))
#---
.col(c(2, 3))
#---
.row(3)
#---
f <- factor(c("a", "b", "c", "a"))
`[.factor`(f, 2:3)
#---
f <- factor(c("a", "b", "c", "a")); f[2:3, drop = TRUE]
#---
f <- factor(c("a", "b", NA, "a"), exclude = NULL); f[3]
#---
f <- factor(c("lo", "hi", "lo"), levels = c("lo", "hi"), ordered = TRUE); f[2:3]
#---
f <- factor(c("a", "b", "c")); f[[2]]
#---
f <- factor(c(x = "a", y = "b", z = "c")); f[[2]]
#---
f <- factor(c("a", "b", "c", "a")); print.factor(f)
#---
f <- factor(c("a", "b", "c", "a")); print(f, max.levels = 0)
#---
f <- factor(character(0)); print(f)
#---
f <- factor(c("a", "b", "c", "a"), levels = c("c", "b", "a"), ordered = TRUE); print(f)
#---
f <- factor(paste0("lev", 1:40)); print(f)
#---
f <- factor(c("a", NA, "b")); print(f)
#---
f <- factor(c(a = "x", b = "y")); print(f)
#---
f <- factor(c("a", "b", NA)); format(f)
#---
f <- factor(c("a", "bbb", NA)); format(f, width = 5)
#---
f <- factor(c("a", "bbb", NA)); format.factor(f, justify = "right")
#---
f <- factor(c("10", "20", "10")); as.vector(f)
#---
f <- factor(c("10", "20", "10")); as.vector(f, "integer")
#---
f <- factor(c("10", "20", "10")); as.vector(f, "list")
#---
f <- factor(c("TRUE", "FALSE", NA)); as.logical(f)
#---
f <- factor(c("T", "F", "yes")); as.logical(f)
#---
x <- noquote(c("a", "b")); c(x, "c")
#---
x <- noquote(c(a = "a", b = "b")); x[2]
#---
x <- noquote(matrix(c("a", "b", "c", "d"), 2)); x[1, ]
#---
x <- noquote(matrix(c("a", "b", "c", "d"), 2)); x[, 2, drop = FALSE]
#---
x <- noquote(c("a", "b")); as.matrix(x)
#---
t1 <- table(c("a", "b", "a", "c")); as.vector(t1[2:3]); class(t1[2:3])
#---
t1 <- table(c("a", "b", "a", "c")); t1[["a"]]
#---
t2 <- table(c(1, 2, 1, 2), c("x", "y", "y", "y")); t2[1, ]
#---
t2 <- table(c(1, 2, 1, 2), c("x", "y", "y", "y")); t2[, "y", drop = FALSE]
#---
t2 <- table(g = c(1, 2, 1, 2), h = c("x", "y", "y", "y")); aperm(t2)
#---
t2 <- table(g = c(1, 2, 1, 2), h = c("x", "y", "y", "y")); aperm.table(t2, c(2, 1))
#---
t3 <- table(a = c(1, 2, 1), b = c("u", "v", "v"), c = c(TRUE, FALSE, TRUE)); dim(aperm(t3, c(3, 1, 2)))
#---
x <- structure(list(a = 1:3, b = "z"), class = "listof"); print(x)
#---
x <- structure(list(a = 1:3, b = "z"), class = "listof"); x[1]
#---
x <- structure(list(1:3, "z"), class = "listof"); print(x)
#---
x <- structure(list(a = 1, b = "z", c = TRUE), class = "simple.list"); print(x)
#---
x <- structure(list(a = 1, b = "z", c = TRUE), class = "simple.list"); x[2:3]
#---
x <- structure(list(a = 1, b = "z"), class = "simple.list"); print.simple.list(x)
