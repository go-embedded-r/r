## Subsetting with [, [[ and $. A line holding only #--- separates cells.
#---
x <- c(10, 20, 30, 40); x[2]; x[c(1, 3)]; x[-1]; x[-c(1, 2)]; x[c(TRUE, FALSE)]; x[5]; x[0]; x[c(1, NA)]; x[10:11]
#---
x <- c(10, 20, 30); x[c(-1, 2)]
#---
x <- c(a = 1, b = 2, c = 3); x["b"]; x[c("a", "c")]; x["z"]; x[c(TRUE, FALSE, TRUE)]; x[2:3]; x[-2]
#---
x <- c(a = 1, b = 2); x[NA]; x[c(NA, 1)]; x[NA_character_]
#---
x <- 1:5; x[x > 2]; x[x %% 2 == 0]; x[c(TRUE, NA, FALSE, TRUE, FALSE)]
#---
x <- 1:3; x[]; x[TRUE]; x[FALSE]; x[1.9]; x[-1.9]; x[c(1.2, 2.7)]
#---
x <- 1:3; x[1, 2]
#---
x <- 1:3; x["a"]
#---
x <- c("a", "b"); x[3]; x[c(1, 3)]
#---
x <- c(TRUE, NA); x[3]; x[0]
#---
x <- list(1, "a", TRUE); x[2]; x[-1]; x[4]; x[c(1, 3)]
#---
x <- list(a = 1, b = "z"); x["a"]; x["q"]; x[c("b", "a")]
#---
NULL[1]; NULL[[1]]; NULL$a; NULL["a"]
#---
x <- 1:10; x[c(-1, -2, -20)]
#---
x <- 1:3; x[c(TRUE, FALSE, TRUE, TRUE)]
#---
x <- 1:3; x[list(1)]
#---
x <- 1:3; x[quote(a)]
#---
sum[1]
#---
f <- function(x) x; f[1]
#---
x <- list(a = 1, b = list(c = 2, d = "e")); x[["a"]]; x[["b"]]; x[["b"]][["d"]]; x[[c(2, 1)]]; x[[c("b", "d")]]
#---
x <- list(1, 2); x[[3]]
#---
x <- list(1, 2); x[[0]]
#---
x <- list(1, 2); x[[-1]]
#---
x <- list(1, 2, 3); x[[-1]]
#---
x <- list(a = 1); x[["b"]]; x[["a"]]
#---
x <- c(a = 1, b = 2); x[["b"]]; x[[1]]; x[["a", exact = TRUE]]
#---
x <- c(a = 1, b = 2); x[["z"]]
#---
x <- c(1, 2); x[[3]]
#---
x <- c(1, 2); x[[1:2]]
#---
x <- list(abc = 1); x[["a"]]; x[["a", exact = FALSE]]
#---
x <- 1:3; x[[2]]; x[[NA]]
#---
x <- list(1); x[[NA]]; x[[NA_integer_]]
#---
x <- list(1, 2); x[[integer(0)]]
#---
x <- list(1, 2); x[[c(1, 2)]]
#---
x <- c(a = 5); x[[1]]
#---
x <- list(a = 1, bb = 2, bcd = 3); x$a; x$b; x$bc; x$bb; x$z; x[["b"]]
#---
x <- list(abc = 1, abd = 2); x$ab; x$abc
#---
x <- c(a = 1); x$a
#---
f <- function(x) x; f$a
#---
x <- list(`a b` = 1); x$`a b`; x$"a b"
#---
x <- 1:6; dim(x) <- c(2, 3); x[1, 2]; x[2, ]; x[, 3]; x[, 2:3]; x[1, , drop = FALSE]; x[, 1, drop = FALSE]; x[-1, ]; x[c(TRUE, FALSE), ]
#---
x <- 1:6; dim(x) <- c(2, 3); x[3, 1]
#---
x <- 1:6; dim(x) <- c(2, 3); x[1, 2, 3]
#---
x <- 1:6; dim(x) <- c(2, 3); x[5]; x[c(1, 6)]; x[[2, 3]]; x[[5]]
#---
x <- 1:6; dim(x) <- c(2, 3); x[NA, 1]; x[c(1, NA), 2]
#---
x <- 1:24; dim(x) <- c(2, 3, 4); x[2, 3, 4]; x[, , 1]; x[1, , ]; x[, 2, 3]; x[2, 2:3, 1:2]
#---
x <- 1:6; dim(x) <- c(2, 3); dimnames(x) <- list(c("a", "b"), c("p", "q", "r")); x["a", "q"]; x["b", ]; x[, "r"]; x[, c("p", "r")]; x["a", "q", drop = FALSE]
#---
x <- 1:6; dim(x) <- c(2, 3); dimnames(x) <- list(c("a", "b"), NULL); x["a", ]; x[, 1]; x[2, 2, drop = FALSE]
#---
x <- 1:6; dim(x) <- c(2, 3); x["a", ]
#---
x <- 1:6; dim(x) <- c(2, 3); x[0, ]; x[, 0]
#---
x <- 1:4; names(x) <- c("a", "b", "c", "d"); x[c("a", "a", "d")]
#---
x <- c(1.5, 2.5); x[2L]; x[2.5]; x[-0.5]
#---
x <- structure(1:3, class = "foo", extra = 1); x[1]
#---
x <- structure(c(a = 1, b = 2), extra = 1); x[1]; x[[1]]
#---
f <- structure(1:3, levels = c("a", "b", "c"), class = "factor"); f[2]; f[c(1, 3)]

#---
matrix(1:6, 2); matrix(1:6, ncol = 2); matrix(1:6, 2, byrow = TRUE); matrix(0, 2, 2)
#---
matrix(1:4, 2, dimnames = list(c("a", "b"), c("x", "y")))
#---
matrix(c("a", "b", "c", "d"), 2)
#---
matrix(numeric(0), 2, 2)
#---
matrix(1:6, 2, 3, 4, 5, 6, 7)
#---
array(1:24, c(2, 3, 4))
#---
array(0, c(2, 2)); array(1:3); array(1:4, c(2, 2), list(c("a", "b"), NULL))
#---
m <- matrix(1:6, 2, dimnames = list(c("r1", "r2"), c("A", "B", "C"))); m[0, c(1, 2)]; m[integer(0), 0]; m[, "B"]; m["r2", ]
#---
x <- c(a = 1, b = 2); x[NULL]; x[[-1.9]]; x[-1]
#---
x <- c(TRUE, FALSE, NA); x[[NULL]]
#---
x <- c("p", "q", "r"); x[[c("a", "c")]]
#---
x <- c(x = "p", y = "q"); x[[FALSE]]
#---
letters[1:3]; LETTERS[24:26]
#---
x <- 1:5; names(x) <- letters; x
