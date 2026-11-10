x <- matrix(c(1,2,NA,4,5,6), 3, dimnames = list(letters[1:3], c("p", "q"))); scale(x)
#---
x <- matrix(c(1,2,NA,4,5,6), 3, dimnames = list(letters[1:3], c("p", "q"))); scale(x, FALSE, TRUE)
#---
x <- matrix(c(1,2,NA,4,5,6), 3, dimnames = list(letters[1:3], c("p", "q"))); scale(x, c(a = 1, b = 2), FALSE)
#---
x <- matrix(c(1,2,NA,4,5,6), 3); attributes(scale(x, FALSE, FALSE))
#---
scale(1:3)
#---
scale(matrix(1:6, 3), TRUE, c(2, 4))
#---
scale(matrix(1:6, 3), 1:3)
#---
scale(matrix(1:6, 3), TRUE, 1:3)
#---
unclass(rle(c(1, 1, NA, NA, 2)))
#---
unclass(rle(c(a = 1, b = 1, c = 2, d = 2, e = 2)))
#---
rle(c("a", "a", "b"))$lengths
#---
unclass(rle(character(0)))
#---
unclass(rle(c(TRUE, TRUE, NA)))
#---
rle(matrix(1:4, 2))
#---
rle(NULL)
#---
rle(factor(c("a", "a")))
#---
inverse.rle(rle(c(3L, 3L, 1L)))
#---
inverse.rle(list(lengths = 1:2, values = 1:2, foo = 1))
#---
inverse.rle(list(lengths = c(1.5, 2), values = 1:2))
#---
inverse.rle(list(lengths = c(-1L, 2L), values = 1:2))
#---
inverse.rle(list(lengths = c(2, 1), values = c(a = 1, b = 2)))
#---
inverse.rle(NULL)
#---
inverse.rle(list(lengths = 1:2, values = 1))
#---
inverse.rle(list(lengths = 1L))
#---
inverse.rle(list(lengths = NA_integer_, values = 1))
#---
inverse.rle(list(lengths = 1:2, values = list(1, 2)))
#---
inverse.rle(structure(list(lengths = 2L, values = factor("u")), class = "rle"))
#---
scale(matrix(list(1, 2), 1))
#---
scale(matrix(c("a", "b"), 1))
#---
scale(matrix(c("a", "b"), 1), FALSE, FALSE)
#---
scale(matrix(c("1", "2"), 1), c(1, 2), FALSE)
