p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(1, 1.1))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(1, 1))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(1, 1 + 1e-10))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(1, "a"))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(1:3, 1:4))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(c(a = 1), c(b = 1)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(c(a = 1), 1))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(1, c(a = 1)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal("a", "b"))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(c("a", "b"), c("a", "c")))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(c("a", "b"), c("x", "y")))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(c("a", "b"), c("a")))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(list(1, 2), list(1, 3)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(list(a = 1, b = 2), list(a = 1, c = 2)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(list(1), list(1, 2)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(TRUE, FALSE))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(c(TRUE, TRUE), c(TRUE, FALSE)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(NULL, NULL))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(NULL, 1))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(factor(c("a", "b")), factor(c("a", "c"))))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(factor("a"), "a"))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(1, 1.1, tolerance = 0.2))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(c(1, 2, 3), c(1, 2, 4)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(c(1, NA), c(1, 2)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(c(1, NA), c(1, NA)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(100, 101, scale = 1))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(0, 0.1))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(1L, 1.0))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(1L, 2L))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(1+2i, 1+3i))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(matrix(1:4, 2), matrix(1:4, 2)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(matrix(1:4, 2), matrix(1:6, 3)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(matrix(1:4, 2), 1:4))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(data.frame(a = 1:2), data.frame(a = c(1, 3))))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(1, 2, check.attributes = FALSE))
#---
isTRUE(all.equal(1, 1.5))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(quote(a + b), quote(a + c)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(quote(a), quote(a)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(environment(), globalenv()))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(as.raw(1:2), as.raw(c(1, 3))))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(1:2, c(1.5, 2), tolerance = 0.1))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(structure(1, class = "foo"), structure(1, class = "bar")))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(structure(1, a = 1), structure(1, a = 2)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(structure(1, a = 1), 1))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(as.Date("2024-01-01"), as.Date("2024-01-03")))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(1, 1.1, countEQ = TRUE))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(c(1, 1, 1, 2), c(1, 1, 1, 3)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(c(1, 1, 1, 2), c(1, 1, 1, 3), countEQ = TRUE))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(c("a", NA), c("a", "b")))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(NA, NA))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(1, NA))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(list(1, list(2, 3)), list(1, list(2, 4))))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(1:3, c(1, 2, 3)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(numeric(0), numeric(0)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(1e-20, 2e-20))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(c(a = 1, b = 2), c(a = 1, b = 2.5)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(character(0), 1))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal(TRUE, 1))
#---
isSymmetric(matrix(c(1, 2, 2, 1), 2))
#---
isSymmetric(matrix(1:6, 2))
#---
isSymmetric(matrix(c(1, 2, 2.0000001, 1), 2))
#---
isSymmetric.matrix(diag(3))
#---
isSymmetric(matrix(c(1, 2, 2, 1), 2, dimnames = list(c("a", "b"), c("a", "b"))))
#---
isSymmetric(matrix(c(1, 2, 2, 1), 2, dimnames = list(c("a", "b"), c("x", "y"))))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal.numeric(1, 1.5))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal.character("a", "b"))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal.list(list(1), list(2)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal.factor(factor("a"), factor("b")))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal.raw(as.raw(1), as.raw(2)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal.default(1, 2))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal.language(quote(a), quote(b)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(all.equal.environment(globalenv(), globalenv()))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); p(attr.all.equal(c(a = 1), c(b = 1)))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); f <- function(x) x + 1; g <- function(x) x + 2; p(all.equal(f, g))
#---
p <- function(x) cat(if (isTRUE(x)) "TRUE" else x, sep = "\n"); f <- function(x) x + 1; p(all.equal.function(f, f))
