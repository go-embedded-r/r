## Row labels of a matrix without row names: [i,] right-aligned to the widest, from 9 to 1000 rows (issue #60). A line holding only #--- separates cells.
#---
print(matrix(1:18, 9))
#---
print(matrix(1:20, 10))
#---
print(matrix(1:22, 11))
#---
print(matrix(1:200, 100))
#---
print(matrix(1:198, 99))
#---
print(matrix(1:202, 101))
#---
print(matrix(1:2000, 1000, 2))
#---
print(matrix(c(1.5, 2:20), 10))
#---
print(matrix(letters[1:20], 10))
#---
print(matrix(c(TRUE, FALSE), 10, 2))
#---
print(matrix(1:20, 10, dimnames = list(NULL, c("a", "b"))))
#---
print(matrix(1:999, 999))
#---
print(matrix(1:9, 9, dimnames = list(NULL, "x")))
#---
print(table(rep(1:9, 2), rep(1:2, 9)))
#---
print(array(1:40, c(10, 2, 2)))
#---
print(matrix(list(1, 2), 10, 2))
#---
print(matrix(complex(real = 1:20), 10))
#---
print(matrix(as.raw(1:20), 10))
#---
print(matrix(1:20, 10), right = TRUE)
#---
matrix(1:30, 10)[8:10, ]
#---
matrix(1:30, 10)[c(1, 10), , drop = FALSE]
#---
t(matrix(1:20, 2))
#---
as.matrix(1:12)
#---
print(matrix(1:2000, 100))
#---
print(matrix(1:60, 10, dimnames = list(NULL, paste0("column", 1:6))))
#---
print(matrix(numeric(0), 10, 0))
#---
print(matrix(NA, 10, 1))
